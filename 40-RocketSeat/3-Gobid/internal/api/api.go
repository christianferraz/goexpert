package api

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Api struct {
	Router         *chi.Mux
	UserService    services.UserService
	ProductService services.ProductService
	BidsService    services.BidsService
	Sessions       *scs.SessionManager
	WSupgrader     websocket.Upgrader
	AuctionLobby   services.AuctionLobby
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
