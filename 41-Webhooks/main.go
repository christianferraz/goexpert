package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Mapa para manter o controle dos clientes conectados
var clients = make(map[*websocket.Conn]bool)

// Canal para transmitir mensagens a todos os clientes
var broadcast = make(chan []byte)

// Configuração do Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Permitir todas as conexões - em produção, implemente verificação adequada
		return true
	},
}

func main() {
	// Servir arquivos estáticos do diretório "./public"
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Rota para lidar com conexões WebSocket
	http.HandleFunc("/ws", handleConnections)

	// Iniciar uma goroutine para lidar com mensagens
	go handleMessages()

	// Iniciar o servidor na porta 8080
	fmt.Println("Servidor iniciado na porta 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Atualizar a conexão HTTP para WebSocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao atualizar conexão:", err)
		return
	}
	defer ws.Close()

	// Registrar novo cliente
	clients[ws] = true
	fmt.Println("Cliente conectado")

	for {
		// Ler mensagem do cliente
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler mensagem:", err)
			delete(clients, ws)
			break
		}
		fmt.Printf("Mensagem recebida: %s\n", message)

		// Enviar mensagem para o canal broadcast
		broadcast <- message
	}
}

func handleMessages() {
	for {
		// Receber próxima mensagem do canal broadcast
		message := <-broadcast
		// Enviar mensagem para todos os clientes conectados
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Erro ao enviar mensagem:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
