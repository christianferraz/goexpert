package main

import (
	"database/sql"
	"fmt"

	"github.com/christianferraz/goexpert/6-Banco_Dados/1-Sql_Puro/4-Select_Multiple/internal/entity"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	product, err := SelectAllProduct(db)
	if err != nil {
		panic(err)
	}
	for _, p := range product {
		fmt.Printf("ID: %s, Name: %s, Price: %f\n", p.ID, p.Name, p.Price)
	}
}

func SelectAllProduct(db *sql.DB) ([]entity.Product, error) {
	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []entity.Product
	for rows.Next() {
		var product entity.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}
