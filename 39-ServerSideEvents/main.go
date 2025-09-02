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
	clientsMu sync.RWMutex // Usando RWMutex para melhor performance em leituras
)

func main() {
	// Define o manipulador para a rota do webhook
	http.HandleFunc("/webhook", webhookHandler)

	// Define o manipulador para a rota de SSE
	http.HandleFunc("/events", sseHandler)

	// Serve arquivo estático
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

	// Adiciona limite de tamanho para o corpo da requisição
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limite

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Não foi possível ler o corpo da solicitação", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload Payload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Formato de JSON inválido", http.StatusBadRequest)
		return
	}

	// Validação básica dos campos
	if payload.Event == "" {
		http.Error(w, "Campo 'event' é obrigatório", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("Evento recebido: %s, Dados recebidos: %s", payload.Event, payload.Data)
	broadcastMessage(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Retorna resposta em JSON
	response := map[string]string{"status": "success", "message": "Webhook recebido com sucesso"}
	json.NewEncoder(w).Encode(response)
}

// sseHandler é a função que lida com as conexões SSE
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Verifica se o servidor suporta flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Servidor não suporta streaming", http.StatusInternalServerError)
		return
	}

	// Cria canal para mensagens com buffer
	messageChan := make(chan string, 10)

	// Adiciona cliente à lista
	clientsMu.Lock()
	clients[messageChan] = true
	clientsMu.Unlock()

	// Remove cliente quando a conexão for fechada
	defer func() {
		clientsMu.Lock()
		delete(clients, messageChan)
		clientsMu.Unlock()
		close(messageChan)
		log.Println("Cliente SSE desconectado")
	}()

	// Define headers para SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Cache-Control")

	// Envia mensagem inicial de conexão
	fmt.Fprintf(w, "data: Conectado ao stream de eventos\n\n")
	flusher.Flush()

	log.Println("Novo cliente SSE conectado")

	// Loop principal do SSE
	for {
		select {
		case message, ok := <-messageChan:
			if !ok {
				return // Canal foi fechado
			}
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			log.Println("Contexto da requisição cancelado")
			return
		}
	}
}

// broadcastMessage envia uma mensagem para todos os clientes SSE conectados
func broadcastMessage(message string) {
	clientsMu.RLock()
	defer clientsMu.RUnlock()

	log.Printf("Broadcasting mensagem para %d clientes: %s", len(clients), message)

	for client := range clients {
		select {
		case client <- message:
			// Mensagem enviada com sucesso
		default:
			// Canal está cheio ou bloqueado, pula este cliente
			log.Println("Aviso: Cliente com canal cheio, pulando...")
		}
	}
}
