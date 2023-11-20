package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

func main() {
	c := http.Client{}
	// precisa antes, bufferizar os dados para ser lido pelo Reader
	jsonVar := bytes.NewBuffer([]byte(`{"nome":"jefferson"}`))
	resp, err := c.Post("https://www.google.com.br", "application/json", jsonVar)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.CopyBuffer(os.Stdout, resp.Body, nil)
}
