package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/christianferraz/goexpert/13-Graphql/graph"
	"github.com/christianferraz/goexpert/13-Graphql/internal/database"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// inserir o db para conexão com o banco
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	categoryDB := database.NewCategory(db)
	courseDB := database.NewCourse(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CategoryDB: categoryDB,
		CoursesDB:  courseDB,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
