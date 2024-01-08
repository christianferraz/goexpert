package main

import (
	"log"
	"net/http"

	"github.com/christianferraz/goexpert/2-Alura/2-LojaVirtual/routes"
)

func main() {

	// Serve static files

	routes.CarregaRotas()
	log.Fatal(http.ListenAndServe(":8000", nil))
}
