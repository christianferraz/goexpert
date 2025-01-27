package api

import (
	"errors"
	"net/http"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/jsonutils"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/services"
	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/usecase/user"
)

func (a *Api) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusBadRequest, problems)
	}
	id, err := a.UserService.CreateUser(r.Context(),
		data.UserName,
		data.Email,
		data.PasswordHash,
		data.Bio,
	)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrPassword) {
			_ = jsonutils.EncodeJSON(w, r, http.StatusConflict, map[string]string{"error": err.Error()})
		}
	}
	_ = jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]string{"id": id.String()})
}

func (a *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {

	panic("not implemented")
}

func (a *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {

	panic("not implemented")
}
