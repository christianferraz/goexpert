package entity

import (
	"errors"
)

type Order struct {
	ID         string
	Price      float64
	Tax        float64
	FinalPrice float64
}

func NewOrder(id string, price, tax float64) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}

	if err := order.IsValid(); err != nil {
		return nil, err
	}
	return order, nil
}

func (order *Order) IsValid() error {
	if order.ID == "" {
		return errors.New("invalid price")
	}
	if order.Price <= 0 {
		return errors.New("order price cannot be zero or negative")
	}
	if order.Tax <= 0 {
		return errors.New("order tax cannot be zero or negative")
	}
	return nil
}

func (order *Order) CalculateFinalPrice() error {
	order.FinalPrice = order.Price + order.Tax
	if err := order.IsValid(); err != nil {
		return err
	}
	return nil
}
