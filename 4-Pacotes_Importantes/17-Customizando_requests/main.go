package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	c := http.Client{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	// precisa antes, bufferizar os dados para ser lido pelo Reader
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com.br", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Print(string(body))
}
