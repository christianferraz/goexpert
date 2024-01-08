package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func ConectaComBancoDeDados() *sql.DB {
	conexao := "user=postgres dbname=alura_loja password=Postgres2024! host=localhost sslmode=disable"
	db, err := sql.Open("postgres", conexao)
	if err != nil {
		log.Fatal(err)
	}
	// Verifica a conexão com o banco de dados
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}
