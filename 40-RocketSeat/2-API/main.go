package main

/* Encurtamento de URL */

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/christianferraz/goexpert/40-RocketSeat/2-API/api"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fail to log", "error", err.Error())
		panic(err)
	}
	slog.Info("done")
}

func run() error {
	db := make(map[string]string)
	handler := api.NewHandler(db)
	s := http.Server{
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
		Addr:         ":8080",
		Handler:      handler,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
