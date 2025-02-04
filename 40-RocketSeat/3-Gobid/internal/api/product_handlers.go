package api

import (
	"context"
	"net/http"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/jsonutils"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/services"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/usecase/product"
	"github.com/google/uuid"
)

func (a *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[product.ProductUseCase](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}
	userId, ok := a.Sessions.Get(r.Context(), "AuthenticatedUserId").(uuid.UUID)
	if !ok {
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": "Not authenticated"})
		return
	}
	productId, err := a.ProductService.CreateProduct(r.Context(),
		userId,
		data.ProductName,
		data.Description,
		data.Baseprice,
		data.AuctionEnd,
	)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	// nao passar o contexto da requisição (r.Context()) para o serviço, assim que a requisição terminar o contexto é cancelado
	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)

	auctionRoom := services.NewAuctionRoom(ctx, productId, &a.BidsService)
	a.AuctionLobby.Lock()
	a.AuctionLobby.Rooms[productId] = auctionRoom
	a.AuctionLobby.Unlock()
	_ = jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"message":    "Auction has started with success",
		"product_id": productId.String()})
}
