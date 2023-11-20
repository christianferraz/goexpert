package main

import (
	"html/template"
	"os"
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
	template := template.Must(template.New("template.html").ParseFiles("template.html"))
	err := template.Execute(os.Stdout, curso)
	if err != nil {
		panic(err)
	}
}
