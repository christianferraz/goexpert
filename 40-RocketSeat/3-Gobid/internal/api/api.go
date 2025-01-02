package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Api struct {
	router *chi.Mux
	ctx    context.Context
}

func NewApi(ctx context.Context) *Api {
	return &Api{
		router: chi.NewRouter(),
		ctx:    ctx,
	}
}

func (api *Api) Start() {
	server := &http.Server{

		Addr:    ":8080",
		Handler: api.router,
	}
	server.ListenAndServe()
}
