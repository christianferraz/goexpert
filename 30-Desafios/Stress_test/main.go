/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"flag"
	"fmt"

	"github.com/christianferraz/goexpert/30-Desafios/Stress_test/tests"
)

func main() {
	var (
		url         = flag.String("url", "", "URL do serviço a ser testado")
		requests    = flag.Int("requests", 100, "Número total de requests")
		concurrency = flag.Int("concurrency", 10, "Número de chamadas simultâneas")
	)
	flag.Parse()
	if *url == "" {
		fmt.Println("URL é um parâmetro obrigatório")
		return
	}
	fmt.Printf("Iniciando testes para %s com %d requests e concorrência de %d\n", *url, *requests, *concurrency)
	report := tests.RunLoadTest(*url, *requests, *concurrency)
	report.Print()
}
