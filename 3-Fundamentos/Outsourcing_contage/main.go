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

	file, err := xlsx.OpenFile("teste.xlsx")
	if err != nil {
		fmt.Println(err.Error())
	}
	sheet := file.Sheet["Sheet"]
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	j := 0
	for i, arquivo := range arquivosExcel {

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
		valorPaginas, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(13, 1).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas: %v\n", err)
			continue
		}
		qtdPaginas, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(12, 1).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas: %v\n", err)
			continue
		}

		valorPaginasExcedentes, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(15, 1).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas excedentes: %v\n", err)
			continue
		}

		qtdPaginasExcedentes, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(14, 1).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas excedentes: %v\n", err)
			continue
		}

		valorPaginasColoridas, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(12, 3).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas coloridas: %v\n", err)
			continue
		}
		qtdPaginasColoridas, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(13, 3).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas coloridas: %v\n", err)
			continue
		}

		valorPaginasColoridasExcedentes, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(15, 3).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas coloridas excedentes: %v\n", err)
			continue
		}
		qtdPaginasColoridasExcedentes, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(14, 3).String()))
		if err != nil {
			log.Printf("Erro ao converter o valor de páginas coloridas excedentes: %v\n", err)
			continue
		}
		// valor_total, err := strconv.Atoi(strings.TrimSpace(planilha.Cell(22, 2).String()))
		// if err != nil {
		// 	log.Printf("Erro ao converter o valor de páginas coloridas excedentes: %v\n", err)
		// 	continue
		// }
		// linha x coluna
		i++
		i = i + j
		sheet.Cell(i, 8).SetValue(arquivo)
		sheet.Cell(i, 4).SetValue(qtdPaginas)
		sheet.Cell(i, 5).SetValue(valorPaginas)
		sheet.Cell(i, 4).SetValue(qtdPaginasExcedentes)
		sheet.Cell(i, 5).SetValue(valorPaginasExcedentes)
		sheet.Cell(i, 4).SetValue(qtdPaginasColoridas)
		sheet.Cell(i, 5).SetValue(valorPaginasColoridas)
		sheet.Cell(i, 4).SetValue(qtdPaginasColoridasExcedentes)
		sheet.Cell(i, 5).SetValue(valorPaginasColoridasExcedentes)
		j = j + 4

	}
	// sheet.File.Save("teste.xlsx")
	// Imprimir os resultados

}
