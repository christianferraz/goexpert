package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"regexp"
)

func main() {
	// Abrindo o arquivo de texto
	file, err := os.Open("dados.txt")
	if err != nil {
		log.Fatalf("Não foi possível abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Preparando para escrever no arquivo CSV
	csvFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("Não foi possível criar o arquivo CSV: %v", err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Escrevendo o cabeçalho do CSV
	writer.Write([]string{"Descrição", "Quantidade", "Valor Unitário"})

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`Qtde\.:(\d+(?:,\d+)?)\s+UN:\s+[A-Z]{2}\s+Vl.\s+Unit\.\:\s+(\d+(?:,\d+)?)`)

	// Lendo o arquivo linha por linha
	for scanner.Scan() {
		descricao := scanner.Text()

		// Lendo a próxima linha para obter a quantidade e o valor unitário
		if scanner.Scan() {
			linhaQuantidade := scanner.Text()
			matches := re.FindStringSubmatch(linhaQuantidade)

			if len(matches) > 1 {
				// Escrevendo no CSV
				writer.Write([]string{descricao, matches[1], matches[2]})
			}

			// Lendo a terceira linha para ignorá-la (valor total)
			if !scanner.Scan() {
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Erro ao ler o arquivo: %v", err)
	}
}
