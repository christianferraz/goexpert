package main

import (
	"html/template"
	"net/http"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func main() {
	curso := Cursos{
		{
			"Curso go",
			340,
		},
		{
			"Curso de Python",
			10,
		},
		{
			"Curso de PHP",
			30,
		},
	}
	templates := []string{
		"header.html",
		"content.html",
		"footer.html",
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.New("content.html").ParseFiles(templates...))
		err := template.Execute(w, curso)
		if err != nil {
			panic(err)
		}

	})
	http.ListenAndServe(":8080", nil)
}
