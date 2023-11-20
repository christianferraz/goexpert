package main

import (
	"html/template"
	"os"
)

type Curso struct {
	Nome         string
	CargaHoraria int
}

func main() {
	curso := Curso{
		"Curso go",
		340,
	}
	template := template.New("CursoTemplate")
	template, _ = template.Parse("Curso: {{.Nome}} e Carga horária: {{.CargaHoraria}}\n")
	err := template.Execute(os.Stdout, curso)
	if err != nil {
		panic(err)
	}
}
