package api

import "net/http"

func (a *Api) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Create User"))
}

func (a *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {

	panic("not implemented")
}

func (a *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {

	panic("not implemented")
}
