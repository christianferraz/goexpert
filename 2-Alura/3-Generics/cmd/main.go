package main

import (
	"fmt"

	"github.com/christianferraz/goexpert/2-Alura/3-Generics/internal/models"
	"github.com/christianferraz/goexpert/2-Alura/3-Generics/internal/services"
)

func main() {
	fmt.Println("Sistema de Estoque")
	estoque := services.NewEstoque()
	itens := []models.Item{
		{ID: 1, Name: "Fone", Quantity: 10, Price: 100},
		{ID: 2, Name: "Camiseta", Quantity: 1, Price: 55.99},
		{ID: 3, Name: "Mouse", Quantity: 2, Price: 12.99},
	}

	for _, item := range itens {
		err := estoque.AddItem(item)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	for _, item := range estoque.ListItems() {
		fmt.Printf("\nID: %d | Item: %s | Quantidade: %d | Preço: %2.f", item.ID, item.Name, item.Quantity, item.Price)
	}
	fmt.Println("\n\nCusto total do estoque:", estoque.CalculateTotalCost())
	fmt.Println("\n\nLogs de auditoria:")
	logs := estoque.ViewAuditLogs()
	for _, log := range logs {
		fmt.Printf("\n%s: %s - %d", log.Timestamp.Format("02/01/2006 15:04:05"), log.Action, log.ItemID)
	}
	itemParaBuscar, err := services.FindBy(estoque.ListItems(), func(item models.Item) bool {
		return item.Name == "Camiseta"
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("\n\nItem encontrado: ", itemParaBuscar)
	alura := services.Fornecedor{
		CNPJ:    "123456789",
		Contato: "Fulano",
		Cidade:  "São Paulo",
	}
	fmt.Println(alura.GetInfo())
	if alura.VerificarDisponibilidade(10, 15) {
		fmt.Println("Disponível")
	} else {
		fmt.Println("Indisponível")
	}
}
