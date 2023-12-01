package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
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
	// file := xlsx.NewFile()
	// sheet, err := file.AddSheet("Sheet1")
	// if err != nil {
	// 	fmt.Printf(err.Error())
	// }

	// i := 0
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

		// sheet.Cell(i, 0).SetValue(arquivo)
		// sheet.Cell(i, 1).SetValue(valorPaginas)
		// sheet.Cell(i, 2).SetValue(valorPaginasExcedentes)
		// sheet.Cell(i, 3).SetValue(valorPaginasColoridas)
		// sheet.Cell(i, 4).SetValue(valorPaginasColoridasExcedentes)
		// i++
		somaTotal += valorPaginas + valorPaginasExcedentes + valorPaginasColoridas + valorPaginasColoridasExcedentes
		somaTotalPaginas += valorPaginas
		somaTotalPaginasExcedentes += valorPaginasExcedentes
		somaTotalPaginasColoridas += valorPaginasColoridas
		somaTotalPaginasColoridasExcedentes += valorPaginasColoridasExcedentes
	}
	// sheet.File.Save("teste.xlsx")
	// Imprimir os resultados
	fmt.Printf("O valor total de páginas impressas p&b é %s faltam %s\n", formatNumber(somaTotalPaginas), formatNumber(cotaPaginas-somaTotalPaginas))
	fmt.Printf("O valor total de páginas excedentes impressas p&b é %s faltam %s\n", formatNumber(somaTotalPaginasExcedentes), formatNumber(cotaPaginasExcedentes-somaTotalPaginasExcedentes))
	fmt.Printf("O valor total de páginas coloridas impressas é %s faltam %s\n", formatNumber(somaTotalPaginasColoridas), formatNumber(cotaPaginasColoridas-somaTotalPaginasColoridas))
	fmt.Printf("O valor total de páginas coloridas impressas é %s faltam %s\n", formatNumber(somaTotalPaginasColoridasExcedentes), formatNumber(cotaPaginasColoridasExcedentes-somaTotalPaginasColoridasExcedentes))
	fmt.Printf("O valor total impressas é %s faltam %s\n", formatNumber(somaTotal), formatNumber(cotaTotal-somaTotal))
}

// Função auxiliar para formatar números com separador de milhares
func formatNumber(value int) string {
	// Converter o valor para uma string
	valueStr := strconv.Itoa(value)

	// Remover caracteres não numéricos
	valueStr = regexp.MustCompile(`\D`).ReplaceAllString(valueStr, "")

	// Remover zeros à esquerda
	valueStr = regexp.MustCompile(`^[0]+`).ReplaceAllString(valueStr, "")

	// Adicionar pontos como separadores de milhares
	var result string
	for i := len(valueStr) - 1; i >= 0; i-- {
		if (len(valueStr)-i-1)%3 == 0 && i != len(valueStr)-1 {
			result = "." + result
		}
		result = string(valueStr[i]) + result
	}

	return result
}
