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
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			_ = jsonutils.EncodeJSON(w, r, http.StatusConflict, map[string]string{"error": err.Error()})
		}
	}
	_ = jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]string{"id": id.String()})
}

func (a *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problem, err := jsonutils.DecodeValidJSON[user.LoginUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problem)
	}
	id, err := a.UserService.AuthenticateUser(r.Context(), data.Email, data.PasswordHash)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			_ = jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]any{"error": "Invalid credentials"})
			return
		}
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	err = a.Sessions.RenewToken(r.Context())
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	a.Sessions.Put(r.Context(), "AuthenticatedUserId", id.String())
	_ = jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{"message": "logged in, sucessfully"})
}

func (a *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := a.Sessions.RenewToken(r.Context())
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	a.Sessions.Remove(r.Context(), "AuthenticatedUserId")
	_ = jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{"message": "logged out, sucessfully"})
}
