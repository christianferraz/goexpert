package api

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/services"
	"github.com/go-chi/chi/v5"
)

type Api struct {
	Router      *chi.Mux
	UserService services.UserService
	Sessions    *scs.SessionManager
}

func NewApi(ctx context.Context) *Api {
	return &Api{
		Router: chi.NewRouter(),
	}
}

func (api *Api) Start() {
	server := &http.Server{

		Addr:    ":8080",
		Handler: api.Router,
	}
	server.ListenAndServe()
}
