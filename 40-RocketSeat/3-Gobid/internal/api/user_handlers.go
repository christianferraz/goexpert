package api

import (
	"net/http"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/jsonutils"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/usecase/user"
)

func (a *Api) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusBadRequest, problems)

	}
}

func (a *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {

	panic("not implemented")
}

func (a *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {

	panic("not implemented")
}
