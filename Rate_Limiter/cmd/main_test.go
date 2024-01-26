package main

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

// TestRateLimiterUnderLimit testa se as requisições estão passando quando estão abaixo do limite
func TestRateLimiterUnderLimit(t *testing.T) {
	// Substitua pela lógica de inicialização do seu servidor e rate limiter
	// ...

	totalRequests := 9 // Número de requisições menor que o limite
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func(i int) {
			defer wg.Done()
			resp, err := http.Get("http://localhost:8080/") // Substitua pela URL do seu servidor
			if err != nil || resp.StatusCode != http.StatusOK {
				t.Errorf("Requisição falhou ou foi inesperadamente limitada %v", resp.StatusCode)
			}
			defer resp.Body.Close()
		}(i)
	}

	wg.Wait()
	t.Logf("Todas as requisições foram feitas com sucesso")
}

// TestRateLimiterOverLimit testa se o rate limiter está bloqueando requisições acima do limite
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
			time.Sleep(1000 * time.Millisecond)             // Espaça as requisições
			resp, err := http.Get("http://localhost:8080/") // Substitua pela URL do seu servidor
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
