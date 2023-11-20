package main

import (
	"database/sql"
	"fmt"

	"github.com/christianferraz/goexpert/6-Banco_Dados/1-Sql_Puro/3-Select_one/internal/entity"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	product, err := SelectProduct(db, "e1a8d935-11fc-4da3-ba1e-db5d42eca180")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Product ID: %s\nProduct Name: %s\nProduct Price: %0.2f\n", product.ID, product.Name, product.Price)
}

func SelectProduct(db *sql.DB, id string) (*entity.Product, error) {
	stmt, err := db.Prepare("SELECT * FROM products WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var product entity.Product
	row := stmt.QueryRow(id)
	err = row.Scan(&product.ID, &product.Name, &product.Price)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
