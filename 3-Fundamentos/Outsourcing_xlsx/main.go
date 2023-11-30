package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func main() {
	diretorio := "./xlsx"
	arquivosExcel, err := filepath.Glob(filepath.Join(diretorio, "*.xlsx"))
	if err != nil {
		log.Fatal(err)
	}

	somaTotal := 0
	somaTotalPaginas := 0
	somaTotalPaginasExcedentes := 0
	somaTotalPaginasColoridas := 0
	somaTotalPaginasColoridasExcedentes := 0
	const (
		cotaPaginas                    = 11408659
		cotaPaginasExcedentes          = 7605773
		cotaPaginasColoridas           = 633600
		cotaPaginasColoridasExcedentes = 422400
		cotaTotal                      = 20070432
	)
	for _, arquivo := range arquivosExcel {
		caminhoArquivo, err := filepath.Abs(arquivo)
		if err != nil {
			log.Fatal(err)
		}

		xlFile, err := xlsx.OpenFile(caminhoArquivo)
		if err != nil {
			log.Printf("Erro ao abrir o arquivo %s: %v\n", arquivo, err)
			continue
		}

		// Acesse a planilha 'CONSOLIDADO'
		planilha := xlFile.Sheet["CONSOLIDADO"]

		// Obtenha os valores das células desejadas
		valorPaginas, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(12, 1).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas: %v\n", err)
			continue
		}

		valorPaginasExcedentes, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(14, 1).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas excedentes: %v\n", err)
			continue
		}

		valorPaginasColoridas, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(12, 3).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas coloridas: %v\n", err)
			continue
		}

		valorPaginasColoridasExcedentes, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(14, 3).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas coloridas excedentes: %v\n", err)
			continue
		}

		somaTotal += valorPaginas + valorPaginasExcedentes + valorPaginasColoridas + valorPaginasColoridasExcedentes
		somaTotalPaginas += valorPaginas
		somaTotalPaginasExcedentes += valorPaginasExcedentes
		somaTotalPaginasColoridas += valorPaginasColoridas
		somaTotalPaginasColoridasExcedentes += valorPaginasColoridasExcedentes
	}

	// Imprimir os resultados
	fmt.Printf("O valor total de páginas impressas p&b é %s faltam %s\n", formatNumber(somaTotalPaginas), formatNumber(cotaPaginas-somaTotalPaginas))
	fmt.Printf("O valor total de páginas excedentes impressas p&b é %s faltam %s\n", formatNumber(somaTotalPaginasExcedentes), formatNumber(cotaPaginasExcedentes-somaTotalPaginasExcedentes))
	fmt.Printf("O valor total de páginas coloridas impressas é %s faltam %s\n", formatNumber(somaTotalPaginasColoridas), formatNumber(cotaPaginasColoridas-somaTotalPaginasColoridas))
	fmt.Printf("O valor total de páginas coloridas impressas é %s faltam %s\n", formatNumber(somaTotalPaginasColoridasExcedentes), formatNumber(cotaPaginasColoridasExcedentes-somaTotalPaginasColoridasExcedentes))
	fmt.Printf("O valor total impressas é %s faltam %s\n", formatNumber(somaTotal), formatNumber(cotaTotal-somaTotal))
}

// Função auxiliar para formatar números com separador de milhares
func formatNumber(numero int) string {
	// Convert numero to string
	numeroStr := strconv.Itoa(numero)

	// Determine the length of the number
	numeroLen := len(numeroStr)

	// Check if a thousands separator is needed
	separatorNeeded := numeroLen % 3
	if separatorNeeded == 0 {
		separatorNeeded = 3
	}

	// Iterate over the digits of the number, adding a thousands separator every three digits
	var stringNumero string
	for i := 0; i < numeroLen; i++ {
		if i > 0 && i%3 == separatorNeeded {
			stringNumero += "."
		}
		stringNumero += string(numeroStr[i])
	}

	// Return the formatted number
	return stringNumero
}
