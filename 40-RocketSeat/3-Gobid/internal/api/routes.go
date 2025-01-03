package api

import "github.com/go-chi/chi/v5"

func (api *Api) BindRoutes() {
	api.router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleCreateUser)
				r.Post("/login", api.handleLoginUser)
				r.Post("/logout", api.handleLogoutUser)
			})
		})
	})
}
