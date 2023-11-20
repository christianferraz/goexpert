package main

import (
	"html/template"
	"net/http"
	"strings"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

type Cursos []Curso

func ToUpper(input string) string {
	return strings.ToUpper(input)
}

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
		t := template.New("content.html")
		t.Funcs(template.FuncMap{"ToUpper": ToUpper})
		template := template.Must(t.ParseFiles(templates...))
		err := template.Execute(w, curso)
		if err != nil {
			panic(err)
		}

	})
	http.ListenAndServe(":8080", nil)
}
