package main

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestRateLimiterUnderLimitWithToken(t *testing.T) {
	totalRequests := 99 // Número de requisições menor que o limite
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func(i int) {
			defer wg.Done()
			req, err := http.NewRequest("GET", "http://localhost:8080/", nil)
			if err != nil {
				t.Errorf("Erro ao criar requisição %d: %v", i, err)
				return
			}
			req.Header.Set("API_KEY", "token1")
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("Erro ao fazer requisição %d: %v", i, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Requisição %d falhou ou foi inesperadamente limitada. Status: %v", i, resp.StatusCode)
			}
		}(i)
	}

	wg.Wait()
	t.Logf("Todas as requisições foram feitas com sucesso")
}

func TestRateLimiterUnderLimit(t *testing.T) {
	// ...

	totalRequests := 99 // Número de requisições menor que o limite
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func(i int) {
			defer wg.Done()
			resp, err := http.Get("http://localhost:8080/")

			if err != nil || resp.StatusCode != http.StatusOK {
				t.Errorf("Requisição falhou ou foi inesperadamente limitada %v", resp.StatusCode)
			}
			defer resp.Body.Close()
		}(i)
	}

	wg.Wait()
	t.Logf("Todas as requisições foram feitas com sucesso")
}

func TestRateLimiterOverLimit(t *testing.T) {
	totalRequests := 9 // Número de requisições maior que o limite
	bl := false
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	// Canal para coletar os resultados das requisições
	results := make(chan bool, totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(1000 * time.Millisecond) // Espaça as requisições
			resp, err := http.Get("http://localhost:8080/")
			if err != nil {
				t.Errorf("Erro ao fazer requisição: %v", err)
				results <- false
				return
			}
			defer resp.Body.Close()

			// Verifica se o código de status é 429 Too Many Requests
			results <- resp.StatusCode == http.StatusTooManyRequests
		}()
	}

	wg.Wait()
	close(results)

	// Verifica se todas as requisições excederam o limite
	for result := range results {
		if result {
			t.Logf("O rate limiter  bloqueou  as requisições")
			bl = true
			break
		}
	}
	if !bl {
		t.Errorf("O rate limiter não bloqueou as requisições")
	}
}
