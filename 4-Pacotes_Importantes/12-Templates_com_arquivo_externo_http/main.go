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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.New("template.html").ParseFiles("template.html"))
		err := template.Execute(w, curso)
		if err != nil {
			panic(err)
		}

	})
	http.ListenAndServe(":8080", nil)
}
