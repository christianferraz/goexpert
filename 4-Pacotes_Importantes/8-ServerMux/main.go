package main

import (
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Servidor 1
	wg.Add(1)
	go func() {
		defer wg.Done()
		mux := http.NewServeMux()
		mux.HandleFunc("/", HomeHandler)
		mux.Handle("/blog", &blog{"Blog"})
		http.ListenAndServe(":8080", mux)
	}()

	// Servidor 2
	wg.Add(1)
	go func() {
		defer wg.Done()
		mux2 := http.NewServeMux()
		mux2.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World 2"))
		})
		http.ListenAndServe(":8081", mux2)
	}()

	wg.Wait()
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

type blog struct {
	title string
}

func (b *blog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(b.title))
}
