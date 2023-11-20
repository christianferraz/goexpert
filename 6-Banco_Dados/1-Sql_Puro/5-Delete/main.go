package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = DeleteProduct(db, "82f1aced-24f2-4847-8f84-4deb00ec4c83")
	if err != nil {
		panic(err)
	}
}

func DeleteProduct(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM products WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	a, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	println(a.RowsAffected())
	return nil
}
