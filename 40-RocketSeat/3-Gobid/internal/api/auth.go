package api

import (
	"net/http"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/jsonutils"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !api.Sessions.Exists(r.Context(), "AuthenticatedUserId") {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
