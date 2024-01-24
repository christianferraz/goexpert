package main

import (
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestServerRate(t *testing.T) {
	var wg sync.WaitGroup
	requestsPerSecond := 81
	testDuration := 10 // duração do teste em segundos
	totalRequests := requestsPerSecond * testDuration

	wg.Add(totalRequests)
	for i := 0; i < totalRequests; i++ {
		go func() {
			defer wg.Done()
			time.Sleep(time.Second / time.Duration(requestsPerSecond))
			res, err := http.Get("http://localhost:8080/")
			if res.StatusCode != http.StatusOK {
				t.Errorf("Status code esperado: %d, recebido: %d", http.StatusOK, res.StatusCode)
			}
			if err != nil {
				t.Errorf("Erro na requisição: %v", err)
			}
		}()
	}

	wg.Wait() // Aguarda todas as goroutines terminarem
}
