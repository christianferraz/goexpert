package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"
)

func main() {
	itemsJSON := `[{"item_carrinho_id": "c4b7a01d-494f-4806-89b9-683ffafab2b9", "item_nome": "Mini Curso", "item_quantidade": 1, "item_preco": 0.00, "item_tipo": "mini_curso"}, {"item_carrinho_id": "16e9c5ce-7e03-41fb-bfeb-7184dd79903e", "item_nome": "Ortopedistas não Sócios", "item_quantidade": 1, "item_preco": 1000.00, "item_tipo": "ingresso"}]`

	var items []ItemCarrinho
	err := json.Unmarshal([]byte(itemsJSON), &items)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	emailContent := TemplateEmail("dfsfaf3432432", items, "teste@teste.com")
	fmt.Println(emailContent)

}

type PaymentData struct {
	Number      string
	Total       float64
	Date        string
	Items       []ItemCarrinho
	EmitterName string
}

type ItemCarrinho struct {
	ItemCarrinhoID string  `json:"item_carrinho_id"`
	ItemNome       string  `json:"item_nome"`
	ItemQuantidade int     `json:"item_quantidade"`
	ItemPreco      float64 `json:"item_preco"`
	ItemTipo       string  `json:"item_tipo"`
	ItemTotal      float64 // Adiciona o total para cada item
}

func TemplateEmail(txid string, items []ItemCarrinho, email string) string {
	// Calculate total
	var total float64
	for i, item := range items {
		itemTotal := item.ItemPreco * float64(item.ItemQuantidade)
		items[i].ItemTotal = itemTotal // Armazena o total para o item
		total += itemTotal
	}

	data := PaymentData{
		Number:      txid,
		Total:       total,
		Date:        time.Now().Format("02/01/2006"),
		EmitterName: "SOCIEDADE BRASILEIRA DE ORTOPEDIA PEDIATRICA - SBOP",
		Items:       items,
	}

	tmpl := `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>Confirmação de Pagamento</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { width: 100%; max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { color: #4CAF50; }
        .details { margin-bottom: 20px; }
        table { width: 100%; border-collapse: collapse; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Olá!</h1>
        <div class="header">
            <h2>✓ O pagamento da cobrança Nº {{.Number}} foi aprovado.</h2>
        </div>
        <div class="details">
            <p><strong>Total:</strong> R$ {{printf "%.2f" .Total}}</p>
            <p><strong>Data da Cobrança:</strong> {{.Date}}</p>
            <p><strong>Nº da cobrança:</strong> {{.Number}}</p>
        </div>
        <h3>Dados da Cobrança</h3>
        <table>
            <thead>
                <tr>
                    <th>Itens</th>
                    <th>Quantidade</th>
                    <th>Valor</th>
                    <th>Total</th>
                </tr>
            </thead>
            <tbody>
                {{range .Items}}
                <tr>
                    <td>{{.ItemNome}}</td>
                    <td>{{.ItemQuantidade}}</td>
                    <td>R$ {{printf "%.2f" .ItemPreco}}</td>
                    <td>R$ {{printf "%.2f" .ItemTotal}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        <h3>Dados do emissor</h3>
        <p>{{.EmitterName}}</p>
    </div>
</body>
</html>
`

	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return ""
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return ""
	}

	return buf.String()
}
