package api

import (
	"errors"
	"net/http"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/jsonutils"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProduct := chi.URLParam(r, "product_id")
	productID, err := uuid.Parse(rawProduct)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}
	_, err = api.ProductService.GetProductByID(r.Context(), productID)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJSON(w, r, http.StatusNotFound, map[string]string{"error": "product not found"})
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	userId, ok := api.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	conn, err := api.WSupgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{"error": "could not upgrade connection to a websocket protocol"})
		return
	}
}
