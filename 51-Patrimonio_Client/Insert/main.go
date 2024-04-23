package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tealeg/xlsx"
)

type InventoryItem struct {
	FichaEbserh    string `json:"ficha_ebserh"`
	FichaUFMS      string `json:"ficha_ufms"`
	DescricaoDoBem string `json:"bem_descricao"`
	Setor          int    `json:"id_lotacao"`
	Marca          string `json:"marca"`
	Modelo         string `json:"modelo"`
	Serie          string `json:"serie"`
	DataAquisicao  string `json:"data_tombamento"`
	CondicaoDeUso  string `json:"bem_condicao"`
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
			dataAquisicaoStr := row.Cells[6].String()
			if dataAquisicaoStr == "" {
				dataAquisicaoStr = "01-01-70"
			}
			dataAquisicao, err := time.Parse("01-02-06", dataAquisicaoStr) // MM/DD/YYYY
			if err != nil {
				log.Fatalf("Erro ao analisar a data: %v", err)
			}
			row.Cells[6].SetFormat("02/01/2006")
			item := InventoryItem{
				FichaEbserh:    row.Cells[0].String(),
				FichaUFMS:      row.Cells[1].String(),
				DescricaoDoBem: row.Cells[2].String(),
				Marca:          row.Cells[3].String(),
				Modelo:         row.Cells[4].String(),
				Serie:          row.Cells[5].String(),
				DataAquisicao:  dataAquisicao.Format("02/01/2006"),
				CondicaoDeUso:  row.Cells[7].String(),
				Setor:          1,
			}
			items = append(items, item)
		}
	}

	return items, nil
}

const bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM4ODg2ODIsInJvbGUiOiI5OThhZjk1NWU1YWFmY2FlYTYzNTIxYzRiNmU2ZTczOWNhZGM0MWZjN2RmYzY1ZmM3YzYxZTgyOWNhODI3MzY0Iiwic3ViIjoiZmNmMDQ4YTg2M2U0M2YzZGEwMDEzOWQ2YWE5ODg1NzRlNTgxNTMxOGNkOTRiZjRiYzkzMWUyNTFhN2ZkM2U0OWE3OWRiOGY3ODhlYjY0OTg1NDdhYzRmYTJjNTZjNGZkMDA5NGNkNzNiNmU1MjQ1YjBkY2U2ODQwMjZiNmU2NjIifQ.aWtY9Ez_GAG-LKwqC7oVd9QeqYzx6Iwpy6lft0NV6sU"

func main() {

	excelFileName := "patrimonio.xlsx"
	const apiURL = "http://localhost:8080/patrimonio"

	items, err := ReadExcel(excelFileName)
	if err != nil {
		panic(err)
	}

	for _, item := range items {

		req, err := http.NewRequest("GET", "http://localhost:8080/patrimonio/"+item.FichaEbserh, nil)
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
