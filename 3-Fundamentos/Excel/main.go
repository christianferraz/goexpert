package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/tealeg/xlsx"
)

func main() {
	// Mapa para armazenar a contagem de cada item
	contagem := make(map[string]int)
	var i = 0
	// Abra o arquivo Excel (substitua "caminho/do/arquivo.xlsx" pelo caminho real do seu arquivo)
	arquivo, err := xlsx.OpenFile("query_results-2023-11-27_84628.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	// Iterar sobre as planilhas e células
	for _, planilha := range arquivo.Sheets {
		for _, linha := range planilha.Rows {
			// Certifique-se de que há células suficientes na linha
			if len(linha.Cells) > 0 {
				// Obtenha o valor da célula na primeira coluna
				celula := linha.Cells[5]
				item := strings.TrimSpace(strings.ToLower(celula.String()))

				// Incrementar a contagem para o item no mapa
				contagem[item]++
				i++
			}
		}
	}

	// Imprimir a contagem de cada item
	fmt.Printf("Contagem de itens na coluna: %v\n", i)
	for item, quantidade := range contagem {
		fmt.Printf("%s: %d\n", item, quantidade)
	}
}
