package main

import (
	"net/http"
	"testing"
	"time"
)

func TestServerRate(t *testing.T) {
	limit := time.Tick(1111 * time.Millisecond)

	for i := 0; i < 10; i++ {
		<-limit

		res, err := http.Get("http://localhost:8080/")
		if res.StatusCode != http.StatusOK {
			t.Errorf("Status code esperado: %d, recebido: %d", http.StatusOK, res.StatusCode)
		}
		if err != nil {
			t.Errorf("Erro na requisição: %v", err)
		}

	}

}
