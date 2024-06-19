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
	"sync"
	"time"

	"github.com/tealeg/xlsx"
)

type InventoryItem struct {
	IDMaterial      string  `json:"id_material"`
	Uorg            int     `json:"id_lotacao"`
	CondicaoDeUso   string  `json:"bem_condicao"`
	FichaSiads      string  `json:"ficha_siads"`
	DataAquisicao   string  `json:"data_tombamento"`
	ValorDepreciado float64 `json:"valor_depreciado"`
	ValorAtual      float64 `json:"valor_original"`
	ContaContabil   int     `json:"conta_contabil"`
}

const apiURL = "http://localhost:8080/patrimonio/"
const bearerToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc4ODc2NTksInJvbGUiOiJhZGM4YmM1Zjg0OGExNzZiMThlNTBjNTIwNjM4ZTk0OWI2NjYxNjM0NThlMDMzODUwYThkNjBmYzNiMzYzODRiIiwic3ViIjoiZjQ2NjE3YWFkMjAyZjQ2MDBhMTY3MmY1ZDM4OThiMTEwNjY5NDI1ZTBmOGRlNTFjYzJjZDVhYmNhNWJjMDdiY2E3ZTlkMDE4MzhhMDY5Y2JjMzBiZGVjN2UxYWYyMjY0MThhMDE1NWFlOTZhNjRlYzQ3YWQxZTNmZTQxMTkzYWQifQ.l9r4Dx1tb52Xv8CIm9_wR08Mf0PEhxXzjMWcDyNv6LI"

func ReadExcel(fileName string) error {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, sheet := range xlFile.Sheets {
		for rowIndex, row := range sheet.Rows {
			if rowIndex == 0 { // Pula o cabeçalho
				continue
			}

			wg.Add(1)
			go func(row *xlsx.Row) {
				defer wg.Done()
				processRow(row)
			}(row)
		}
	}
	wg.Wait()
	return nil
}

func processRow(row *xlsx.Row) {
	dataAquisicaoStr := row.Cells[9].String()
	if dataAquisicaoStr == "" {
		dataAquisicaoStr = "01/01/1970"
	}
	dataAquisicao, err := time.Parse("02/01/2006", dataAquisicaoStr) // MM/DD/YYYY
	if err != nil {
		log.Printf("Erro ao analisar a data: %v no item %v", err, row.Cells[0].String())
		return
	}
	valorDep, err := row.Cells[11].Float()
	if err != nil {
		log.Printf("Erro ao analisar o valor depreciado: %v no item %v", err, row.Cells[0].String())
		return
	}
	valorAt, err := row.Cells[10].Float()
	if err != nil {
		log.Printf("Erro ao analisar o valor atual: %v no item %v", err, row.Cells[0].String())
		return
	}
	uorgTmp, err := row.Cells[5].Int()
	if err != nil {
		log.Printf("Erro ao analisar a unidade de origem: %v no item %v", err, row.Cells[0].String())
		return
	}
	contacontabil, err := row.Cells[3].Int()
	if err != nil {
		log.Printf("Erro ao analisar a conta contábil: %v no item %v", err, row.Cells[3].String())
		return
	}

	item := InventoryItem{
		FichaSiads:      row.Cells[0].String(),
		ContaContabil:   contacontabil,
		DataAquisicao:   dataAquisicao.Format("02/01/2006"),
		ValorDepreciado: valorDep,
		ValorAtual:      valorAt,
		CondicaoDeUso:   row.Cells[8].String(),
		Uorg:            uorgTmp,
	}

	if verificaSeExiste(&item) {
		atualizaItem(&item)
	} else {
		insereItem(&item)
	}
}

func insereItem(item *InventoryItem) bool {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return false
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar a solicitação: %v", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta: %v", err)
		return false
	}

	fmt.Printf("Resposta do servidor: %s\n", body)
	return true
}

func atualizaItem(item *InventoryItem) bool {
	jsonData, err := json.Marshal(item)
	if err != nil {
		return false
	}

	req, err := http.NewRequest("PATCH", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Erro ao enviar a solicitação: %v", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler a resposta: %v", err)
		return false
	}

	fmt.Printf("Resposta do servidor: %s\n", body)
	return true
}

func verificaSeExiste(item *InventoryItem) bool {
	req, err := http.NewRequest("GET", apiURL+item.FichaSiads, nil)
	if err != nil {
		fmt.Printf("Erro ao criar a requisição: %v\n", err)
		return false
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer a solicitação: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler a resposta: %v\n", err)
		return false
	}

	if string(bytes.TrimSpace(body)) == "null" {
		fmt.Printf("Item %s não encontrado\n", item.FichaSiads)
		return false
	} else {
		var items []InventoryItem
		err := json.Unmarshal(body, &items)
		if err != nil {
			fmt.Printf("Erro ao desserializar a resposta: %v\n", err)
			return false
		}

		if len(items) == 0 {
			fmt.Printf("Item %s não encontrado\n", item.FichaSiads)
			return false
		}

		// Atualiza o item original com o primeiro item encontrado
		item.IDMaterial = items[0].IDMaterial
		return true
	}
}

func converterParaFloat(valorStr string) (float64, error) {
	valorSemPontos := strings.Replace(valorStr, ".", "", -1)
	valorCorrigido := strings.Replace(valorSemPontos, ",", ".", -1)
	valorFloat, err := strconv.ParseFloat(valorCorrigido, 64)
	if err != nil {
		return 0, err
	}
	return valorFloat, nil
}

func main() {
	excelFileName := "rmb.xlsx"

	err := ReadExcel(excelFileName)
	if err != nil {
		panic(err)
	}
}
