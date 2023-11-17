package main

import "net/http"

func main() {
	http.HandleFunc("/", FuncHandler)
	http.ListenAndServe(":8080", nil)
}

func FuncHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`
		<html>
		<head>
		<meta charset="utf-8">
		<title>Olá Mundo!</title>
		</head>
		
		<h1>Olá Mundo!</h1>
		</html>`))
}
