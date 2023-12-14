package web

import (
	"encoding/json"
	"net/http"

	"github.com/christianferraz/goexpert/20-CleanArch/internal/entity"
	"github.com/christianferraz/goexpert/20-CleanArch/internal/usecase"
	"github.com/christianferraz/goexpert/20-CleanArch/pkg/events"
)

type WebOrderHandler struct {
	EventDispatcher   events.EventDispatcherInterface
	OrderRepository   entity.OrderRepositoryInterface
	OrderCreatedEvent events.EventInterface
}

func NewWebOrderHandler(
	EventDispatcher events.EventDispatcherInterface,
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreatedEvent events.EventInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		EventDispatcher:   EventDispatcher,
		OrderRepository:   OrderRepository,
		OrderCreatedEvent: OrderCreatedEvent,
	}
}

func (h *WebOrderHandler) Orders(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		var dto usecase.OrderInputDTO
		err := json.NewDecoder(r.Body).Decode(&dto)
		if err != nil {
			println(" carai ", r.Method)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
		output, err := createOrder.Execute(dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "GET" {
		listOrders := usecase.NewListOrderUseCase(h.OrderRepository, h.OrderCreatedEvent, h.EventDispatcher)
		dto, err := listOrders.Execute()
		err = json.NewEncoder(w).Encode(dto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
