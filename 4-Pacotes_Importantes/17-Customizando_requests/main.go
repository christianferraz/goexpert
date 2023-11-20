package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	c := http.Client{}
	// precisa antes, bufferizar os dados para ser lido pelo Reader
	req, err := http.NewRequest(http.MethodGet, "https://www.google.com.br", nil)
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
