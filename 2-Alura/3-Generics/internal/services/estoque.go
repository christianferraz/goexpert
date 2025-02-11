package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/christianferraz/goexpert/2-Alura/3-Generics/internal/models"
)

type Estoque struct {
	items map[string]models.Item
	logs  []models.Log
}

func NewEstoque() *Estoque {
	return &Estoque{
		items: make(map[string]models.Item),
		logs:  []models.Log{},
	}
}

func (e *Estoque) AddItem(item models.Item) error {
	if item.Quantity <= 0 {
		return fmt.Errorf("erro ao adicionar item: [ID:%d] possui uma quantidade inválida (zero ou negativa)", item.ID)
	}
	existingItem, exists := e.items[strconv.Itoa(item.ID)]
	if exists {
		item.Quantity += existingItem.Quantity
	}
	e.items[strconv.Itoa(item.ID)] = item
	e.logs = append(e.logs, models.Log{Timestamp: time.Now(), Action: "add", ItemID: item.ID})
	return nil
}

func (e *Estoque) ListItems() []models.Item {
	var itemList []models.Item
	for _, item := range e.items {
		itemList = append(itemList, item)
	}
	return itemList
}

func (e *Estoque) ViewAuditLogs() []models.Log {
	return e.logs
}

func (e *Estoque) CalculateTotalCost() float64 {
	var totalCost float64
	for _, item := range e.items {
		totalCost += float64(item.Quantity) * item.Price
	}
	return totalCost
}

func (e *Estoque) RemoveItem(itemID int, quantity int, user string) error {
	item, exists := e.items[strconv.Itoa(itemID)]
	if !exists {
		return fmt.Errorf("item não encontrado")
	}
	if quantity <= 0 {
		return fmt.Errorf("quantidade inválida")
	}
	if item.Quantity < quantity {
		return fmt.Errorf("quantidade insuficiente")
	}
	item.Quantity -= quantity
	if item.Quantity == 0 {
		delete(e.items, strconv.Itoa(itemID))
	}
	e.logs = append(e.logs, models.Log{Timestamp: time.Now(), Action: "remove", ItemID: itemID, User: user})
	return nil
}

func FindBy[T any](data []T, comparator func(T) bool) ([]T, error) {
	var result []T
	for _, v := range data {
		if comparator(v) {
			result = append(result, v)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("nenhum item foi encontrado")
	}
	return result, nil
}
