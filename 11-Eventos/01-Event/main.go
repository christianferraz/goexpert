package main

import (
	"time"

	"github.com/christianferraz/goexpert/11-Eventos/01-Event/internal"
)

func main() {
	dispatcher := internal.NewEventDispatcher()

	// Registra os listeners para o evento "user_created"
	dispatcher.Register("user_created", internal.UserCreatedListener)
	dispatcher.Register("user_created", internal.LogUserCreated)

	// Cria um evento de usuário criado
	userData := map[string]string{
		"name":  "John Doe",
		"email": "john.doe@example.com",
	}
	event := internal.Event{Name: "user_created", Data: userData}

	// Despacha o evento
	dispatcher.Dispatch(event)

	// Aguarda um momento para garantir que os listeners processem o evento
	time.Sleep(1 * time.Second)
}
