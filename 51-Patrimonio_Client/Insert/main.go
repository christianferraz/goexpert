package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type InventoryItem struct {
	Uorg            int     `json:"id_lotacao"`
	CondicaoDeUso   string  `json:"bem_condicao"`
	FichaSiads      string  `json:"ficha_siads"`
	FichaEbserh     string  `json:"ficha_ebserh"`
	DescricaoDoBem  string  `json:"bem_descricao"`
	Marca           string  `json:"marca"`
	Modelo          string  `json:"modelo"`
	Serie           string  `json:"serie"`
	DataAquisicao   string  `json:"data_tombamento"`
	ValorDepreciado float64 `json:"bem_valor_depreciado"`
	ValorAtual      float64 `json:"bem_valor_atual"`
	FichaUFMS       string  `json:"ficha_ufms"`
}

func converterParaFloat(valorStr string) (float64, error) {
	// Remover pontos usados para separar milhares
	valorSemPontos := strings.Replace(valorStr, ".", "", -1)
	// Substituir vírgula por ponto para adequar ao formato de número flutuante
	valorCorrigido := strings.Replace(valorSemPontos, ",", ".", -1)

	// Converter a string corrigida para float64
	valorFloat, err := strconv.ParseFloat(valorCorrigido, 64)
	if err != nil {
		return 0, err
	}
	return valorFloat, nil
}

func ReadExcel(fileName string) ([]InventoryItem, error) {
	var items []InventoryItem

	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			if rowIndex == 0 { // Pula o cabeçalho
				continue
			}
			dataAquisicaoStr := row.Cells[8].String()
			if dataAquisicaoStr == "" {
				dataAquisicaoStr = "01-01-70"
			}
			dataAquisicao, err := time.Parse("01-02-06", dataAquisicaoStr) // MM/DD/YYYY
			if err != nil {
				log.Fatalf("Erro ao analisar a data: %v no item %v", err, row.Cells[11].String())
			}
			valorDep, err := row.Cells[9].Float()
			if err != nil {
				log.Fatalf("Erro ao analisar o valor depreciado: %v no item %v", err, row.Cells[11].String())
			}
			valorAt, err := row.Cells[10].Float()
			if err != nil {
				log.Fatalf("Erro ao analisar o valor depreciado: %v no item %v", err, row.Cells[11].String())
			}
			uorgTmp, err := row.Cells[0].Int()
			if err != nil {
				log.Fatalf("Erro ao analisar a unidade de origem: %v no item %v", err, row.Cells[11].String())
			}

			item := InventoryItem{
				FichaSiads:     row.Cells[2].String(),
				FichaEbserh:    row.Cells[3].String(),
				DescricaoDoBem: row.Cells[4].String(),
				Marca:          row.Cells[5].String(),
				Modelo:         row.Cells[6].String(),
				Serie:          row.Cells[7].String(),
				// FichaUFMS:      row.Cells[11].String(),
				DataAquisicao: dataAquisicao.Format("02/01/2006"),
				// DataAquisicao:   dataAquisicaoStr,
				ValorDepreciado: valorDep,
				ValorAtual:      valorAt,
				CondicaoDeUso:   row.Cells[1].String(),
				Uorg:            uorgTmp,
			}
			items = append(items, item)
		}
	}

	return items, nil
}

const bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTQ0OTE5OTQsInJvbGUiOiI3YTI1MDU0YTViYmVkY2ZhZTViNDNiN2ZlNzdhZDBjMjM1ZjcwNTMyNjJkNDYxMDBmYjA3YTljMDI2YWIzODJmIiwic3ViIjoiNTA5MzQ1YTc5YjJjZTJlZDVkNjg4NjkyNjc1Mjk0MGM0Nzk3ODdjOTAzNTZmY2YwZjg5OGI1ZWEzYzUyYWMwMmIzMjYwNjk3YmJmMzc4ZGRjODg3YTY5OGU5MmIwMTFhNjlmZWM3ZmRkNjBhOWMwOTU1N2M2OWZhY2RhMTM0ZDYifQ.kCxVJL-GcmghZeNSaUnFAH2dIuikJKHeJhQUX4m3oi4"

func main() {
	excelFileName := "ufms.xlsx"
	const apiURL = "http://localhost:8080/patrimonio"

	items, err := ReadExcel(excelFileName)
	if err != nil {
		panic(err)
	}

	for _, item := range items {
		req, err := http.NewRequest("GET", "http://localhost:8080/patrimonio/"+item.FichaSiads, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Authorization", "Bearer "+bearerToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			// handle error
			fmt.Printf("Erro ao fazer a solicitação: %v\n", err)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// handle error
			fmt.Printf("Erro ao ler a resposta: %v\n", err)
			return
		}
		if string(bytes.TrimSpace(body)) == "null" {
			err := SendToAPI(item, apiURL)
			if err != nil {
				fmt.Println(item.FichaEbserh, err)
				continue
			}
			fmt.Printf("Item %s enviado para a API\n", item.FichaEbserh)
			continue
		}
	}
}

func SendToAPI(item InventoryItem, apiURL string) error {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Erro ao enviar a solicitação: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro ao ler a resposta: %v", err)
	}

	// Faça algo com a resposta
	fmt.Printf("Resposta do servidor: %s\n", body)

	return nil
}
