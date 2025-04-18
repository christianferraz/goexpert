package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger, api.Sessions.LoadAndSave)
	// csrfMiddleware := csrf.Protect([]byte(os.Getenv("GOBID_CSRF_KEY")), csrf.Secure(false))
	// api.Router.Use(csrfMiddleware)
	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/users", func(r chi.Router) {
				r.Post("/signup", api.handleCreateUser)
				r.Post("/signin", api.handleLoginUser)
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/logout", api.handleLogoutUser)
				})
			})
			r.Route("/products", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/", api.handleCreateProduct)
					r.Get("/ws/subscribe/{product_id}", api.handleSubscribeUserToAuction)
				})
			})
		})
	})
}
