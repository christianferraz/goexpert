package user

import (
	"context"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/validator"
)

type CreateUserReq struct {
	UserName     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Bio          string `json:"bio"`
}

func (req CreateUserReq) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(validator.NotBlank(req.UserName), "username", "This field is not empty")
	eval.CheckField(validator.NotBlank(req.Email), "email", "This field is not empty")
	eval.CheckField(validator.NotBlank(req.Bio), "bio", "This field is not empty")
	eval.CheckField(validator.MinChars(req.UserName, 10) && validator.MaxChars(req.UserName, 255), "username", "This field must have a length between 10 and 255 characters")
	eval.CheckField(validator.MinChars(req.PasswordHash, 8), "password_hash", "This field must have a length greater than 8 characters")
	eval.CheckField(validator.Matches(req.Email, validator.EmailRX), "email", "This field must be a valid email")
	return eval
}
