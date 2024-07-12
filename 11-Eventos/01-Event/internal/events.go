package internal

import (
	"sync"
)

// Event representa um evento genérico
type Event struct {
	Name string
	Data interface{}
}

// EventHandler é uma função que lida com eventos
type EventHandler func(event Event)

// EventDispatcher gerencia eventos e seus ouvintes
type EventDispatcher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewEventDispatcher cria uma nova instância de EventDispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// Register adiciona um ouvinte para um evento específico
func (d *EventDispatcher) Register(eventName string, handler EventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.handlers[eventName] = append(d.handlers[eventName], handler)
}

// Dispatch envia um evento para todos os ouvintes registrados
func (d *EventDispatcher) Dispatch(event Event) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if handlers, found := d.handlers[event.Name]; found {
		for _, handler := range handlers {
			go handler(event) // Executa o handler em uma goroutine
		}
	}
}
