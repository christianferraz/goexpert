package internal

import "fmt"

// Listener 1: Registra a criação de um usuário
func UserCreatedListener(event Event) {
	if user, ok := event.Data.(map[string]string); ok {
		fmt.Printf("Usuário criado: %s (Email: %s)\n", user["name"], user["email"])
	}
}

// Listener 2: Loga a criação de um usuário
func LogUserCreated(event Event) {
	if user, ok := event.Data.(map[string]string); ok {
		fmt.Printf("Log: Um novo usuário foi criado: %s\n", user["name"])
	}
}
