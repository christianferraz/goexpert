package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	c := http.Client{Timeout: time.Second * 5}
	req, err := c.Get("https://www.google.com.br")
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}
