package jsonutils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/validator"
)

func EncodeJSON[T any](w http.ResponseWriter, r *http.Request, statuscode int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statuscode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	return nil
}

func DecodeValidJSON[T validator.Validator](r *http.Request) (T, map[string]string, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, nil, fmt.Errorf("failed to decode JSON: %w", err)
	}
	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}
	return data, nil, nil
}

func DecodeJSON[T any](r *http.Request) (T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return data, fmt.Errorf("failed to decode JSON: %w", err)
	}
	return data, nil
}
