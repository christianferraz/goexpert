package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/tealeg/xlsx"
)

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

const bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM4ODg2ODIsInJvbGUiOiI5OThhZjk1NWU1YWFmY2FlYTYzNTIxYzRiNmU2ZTczOWNhZGM0MWZjN2RmYzY1ZmM3YzYxZTgyOWNhODI3MzY0Iiwic3ViIjoiZmNmMDQ4YTg2M2U0M2YzZGEwMDEzOWQ2YWE5ODg1NzRlNTgxNTMxOGNkOTRiZjRiYzkzMWUyNTFhN2ZkM2U0OWE3OWRiOGY3ODhlYjY0OTg1NDdhYzRmYTJjNTZjNGZkMDA5NGNkNzNiNmU1MjQ1YjBkY2U2ODQwMjZiNmU2NjIifQ.aWtY9Ez_GAG-LKwqC7oVd9QeqYzx6Iwpy6lft0NV6sU"

func main() {

	excelFileName := "a.xlsx"

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
			continue

		}
	}
}
