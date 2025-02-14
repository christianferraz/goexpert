package api

import "net/http"

func (api *Application) handleCreateTAsk(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Task List"))
}
