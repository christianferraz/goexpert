package entity

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type Product struct {
	ID    string
	Name  string
	Price float64
	DB    *sql.DB
}

func NewProduct(db *sql.DB) *Product {
	return &Product{
		DB: db,
	}
}

func (p *Product) InsertProduct(name string, price float64) (*Product, error) {
	id := uuid.New().String()
	stmt, err := p.DB.Prepare("INSERT INTO products(id, name, price) VALUES(?,?,?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id, name, price)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:    id,
		Name:  name,
		Price: price,
		DB:    p.DB,
	}, nil

}

func (p *Product) UpdateProduct(id string, name string, price float64) error {
	stmt, err := p.DB.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")
	if err != nil {
		log.Fatal("eu")
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(name, price, id)
	if err != nil {
		return err
	}
	return nil
}
