package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	req, err := http.Get("http://www.google.com.br")
	if err != nil {
		fmt.Println("Erro ao abrir a página do Google:", err)
		return
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Erro ao ler a página do Google:", err)
		return
	}
	fmt.Println(string(res))
}
