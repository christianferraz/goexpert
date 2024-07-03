package user_entity

import (
	"context"

	"github.com/christianferraz/goexpert/26-Leilao/internal/internal_error"
)

type User struct {
	Id   string
	Name string
}

type UserRepositoryInterface interface {
	FindUserById(ctx context.Context, useId string) (*User, *internal_error.InternalError)
}
