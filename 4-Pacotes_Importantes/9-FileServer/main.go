package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	fileserver := http.FileServer(http.Dir("./public"))
	mux.Handle("/", fileserver)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
