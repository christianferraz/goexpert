package user

import (
	"context"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/validator"
)

type LoginUserReq struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

func (req LoginUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(validator.NotBlank(req.Email), "email", "This field is not empty")
	eval.CheckField(validator.NotBlank(req.PasswordHash), "password_hash", "This field is not empty")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "This field must be a valid email")
	return eval
}
