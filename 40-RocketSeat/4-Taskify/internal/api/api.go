package api

import "github.com/go-chi/chi"

type Application struct {
	Router *chi.Mux
}
