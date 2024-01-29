package tax

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

	totalRequests := 100 // Número de requisições menor que o limite
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()
			resp, err := http.Get("http://localhost:8080/") // Substitua pela URL do seu servidor
			if err != nil || resp.StatusCode != http.StatusOK {
				t.Errorf("Requisição falhou ou foi inesperadamente limitada")
			}
		}()
	}

	wg.Wait()
}

// TestRateLimiterOverLimit testa se o rate limiter está bloqueando requisições acima do limite
func TestRateLimiterOverLimit(t *testing.T) {

	totalRequests := 15 // Número de requisições maior que o limite
	var wg sync.WaitGroup
	wg.Add(totalRequests)

	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(1000 * time.Millisecond)           // Espaça as requisições
			resp, _ := http.Get("http://localhost:8080/") // Substitua pela URL do seu servidor
			if resp.StatusCode == http.StatusTooManyRequests {
				t.Log("Requisição corretamente limitada")
			} else {
				t.Errorf("Requisição passou inesperadamente")
			}
		}()
	}

	wg.Wait()
}
