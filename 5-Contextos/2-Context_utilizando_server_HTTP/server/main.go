package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("request chegando")
	defer log.Println("request finalizado")
	select {
	case <-ctx.Done():
		log.Println("request cancelada pelo cliente")
	case <-time.After(time.Second * 5):
		log.Println("processamento finalizado")
		w.Write([]byte("Deu tudo ok"))
	}
}
