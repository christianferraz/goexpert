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

const bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTQ0OTY4NzgsInJvbGUiOiIxYTFhMGJiZjA2ZjMxYWI2ZTA2ZTU1NmFiZTNhMjhhZmJhOTUxNWQ4ZGNjODQwZGZlM2NjN2E1NDliZTg2MTY2Iiwic3ViIjoiYjVjNmU3YzIxODMwYjY0MDljZGU4YjhkN2FkMzVkYjNiODExYWRiYzg1NDY2ZjlhMmRlNzNmNDE5YjliNTJlMDM1NWFkMDg1NzNkYjg1NDEzMjg1MGI0NTUyNTU4ODI3MDIzNTRhZWE3MWRmYWUyZWY3ZGQyOGMzMDNkMDE1YzUifQ.EtNC4xjUq6ZvAwSJgd5tDIS0R8nlV0AXlZ6Mo7oxdog"

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

		}
	}
	fmt.Printf("Itens nao encontrados %v\n", i)
}
