package main

import "net/http"

func main() {
	http.HandleFunc("/", FuncHandler)
	http.ListenAndServe(":8080", nil)
}

func FuncHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
