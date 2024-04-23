package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /books/{id}", GetBookHandler)
	mux.HandleFunc("GET /books/{d...}", BooksPathHandler)
	// caminho exato
	mux.HandleFunc("GET /books/test/{$}", BooksHandler)
	if err := http.ListenAndServe(":9000", mux); err != nil {
		fmt.Println(err)
	}
}

func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("Book " + id))
}

func BooksPathHandler(w http.ResponseWriter, r *http.Request) {
	dirpath := r.PathValue("d")
	fmt.Fprintf(w, "Books in %s\n", dirpath)
}

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All books"))
}
