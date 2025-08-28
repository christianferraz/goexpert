package main

import (
	"context"
	"fmt"

	"github.com/christianferraz/goexpert/44-Validator/internal/validator"
)

// Definir struct com tags
type Usuario struct {
	Nome  string `validate:"required,min=2,max=50,alpha"`
	Email string `validate:"required,email"`
	CPF   string `validate:"required,cpf"`
	Site  string `validate:"omitempty,url"`
	Idade int    `validate:"required,min=18,max=100"`
}

func main() {
	ctx := context.Background()
	erros := validator.ValidateStruct(ctx, Usuario{Nome: "Christian", Email: "invalid-email", CPF: "123.456.789-09", Site: "https://example.com"})
	if len(erros) == 0 {
		// Válido!
	} else {
		// Processar erros
		for campo, erro := range erros {
			fmt.Printf("%s: %s\n", campo, erro)
		}
	}
}

// Validar
