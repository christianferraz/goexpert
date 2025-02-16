package api

import (
	"github.com/christianferraz/goexpert/40-RocketSeat/4-Taskify/internal/services"
	"github.com/go-chi/chi/v5"
)

type Application struct {
	Router      *chi.Mux
	TaskService services.TaskService
}
