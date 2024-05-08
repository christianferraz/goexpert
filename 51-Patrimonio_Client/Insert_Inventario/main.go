package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tealeg/xlsx"
)

type MaterialInputAddDTO struct {
	IDMaterial  string `json:"id_material"`
	FichaEBSERH string `json:"ficha_ebserh"`
	FichaUFMS   string `json:"ficha_ufms"`
	FichaSIADS  string `json:"ficha_siads"`
}

type MaterialInventarioAddDTO struct {
	IDInventario string `json:"id_inventario"`
	IDMaterial   string `json:"id_material"`
	Status       string `json:"status"`
	Observacoes  string `json:"observacoes"`
}

func ReadExcel(fileName string) ([]string, error) {
	var items []string

	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			if rowIndex == 0 { // Pula o cabeçalho
				continue
			}

			items = append(items, row.Cells[0].String())
		}
	}

	return items, nil
}

const bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUxMDU3NDcsInJvbGUiOiIyNjc2NjZjZTkyOThmNTk5MGYzMzE4ZGRiMDgzNGNiMDYxNDhiYjcyZGI5MzMzNmIzYzk5ZTUxZjUyOTY3YTNlIiwic3ViIjoiNzkwYTIyN2ZhNjk1NjA3ODNkOWMwZGU0MzQzMzQ1ODc4ODEyMjEwY2ZlMDVlOTQzNDRhMTlmNmQzMjM4ZWE5MjBjNDZlYjNhOTE1NTFjNWNkNTQxZDJiZDYyMWViZWFmY2Y5NmU0NmRhNWIzZWVlZDJlMDA0MjFiZDhmNmUwMmEifQ.ynaf0sXHPv85DvkRWLBnt0AMCNCxEUxQAd3q9xptbiY"

func main() {

	excelFileName := "a.xlsx"
	i := 0
	items, err := ReadExcel(excelFileName)
	if err != nil {
		panic(err)
	}

	for _, item := range items {

		req, err := http.NewRequest("GET", "http://localhost:8080/patrimonio/"+item, nil)
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

			fmt.Printf("Item %s nao encontrado\n", item)
			i++
			continue

		} else {
			var material MaterialInputAddDTO
			err := json.Unmarshal(body, &material)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Item %v encontrado\n", material.IDMaterial)
			SendToAPI(&MaterialInventarioAddDTO{
				IDInventario: "c2299d9c-32dc-4646-a013-8e71e5764392",
				IDMaterial:   material.IDMaterial,
				Status:       "Disponivel",
				Observacoes:  "",
			}, "http://localhost:8080/inventario/material")
		}
	}
	fmt.Printf("Itens nao encontrados %v\n", i)
}

func SendToAPI(item *MaterialInventarioAddDTO, apiURL string) error {
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
