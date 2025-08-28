package jsonutils

import (
	"fmt"
	"net/http"

	"github.com/christianferraz/goexpert/44-Validator/internal/validator"
)

func DecodeValidJSON[T validator.Validator](r *http.Request) (T, map[string]string, error) {
	var data T

	if problems := data.Valid(r.Context()); len(problems) > 0 {
		return data, problems, fmt.Errorf("invalid %T: %d problems", data, len(problems))
	}

	return data, nil, nil
}
