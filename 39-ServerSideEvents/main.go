package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// Payload é a estrutura que representa os dados que o webhook receberá
type Payload struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

var (
	clients   = make(map[chan string]bool)
	clientsMu sync.Mutex
)

func main() {
	// Define o manipulador para a rota do webhook
	http.HandleFunc("/webhook", webhookHandler)

	// Define o manipulador para a rota de SSE
	http.HandleFunc("/events", sseHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Inicia o servidor HTTP na porta 8080
	log.Println("Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// webhookHandler é a função que lida com as solicitações do webhook
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Não foi possível ler o corpo da solicitação", http.StatusBadRequest)
		return
	}

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Formato de JSON inválido", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("Evento recebido: %s, Dados recebidos: %s", payload.Event, payload.Data)
	broadcastMessage(message)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Webhook recebido com sucesso")
}

// sseHandler é a função que lida com as conexões SSE
func sseHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Servidor não suporta streaming", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan string)
	clientsMu.Lock()
	clients[messageChan] = true
	clientsMu.Unlock()

	defer func() {
		clientsMu.Lock()
		delete(clients, messageChan)
		clientsMu.Unlock()
		close(messageChan)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case message := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// broadcastMessage envia uma mensagem para todos os clientes SSE conectados
func broadcastMessage(message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	for client := range clients {
		client <- message
	}
}
