package user_usecase

import (
	"context"

	"github.com/christianferraz/goexpert/26-Leilao/internal/entity/user_entity"
	"github.com/christianferraz/goexpert/26-Leilao/internal/internal_error"
)

type UserUseCase struct {
	UserRepository user_entity.UserRepositoryInterface
}

type UserOutputDTO struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	FindUserById(ctx context.Context, id string) (UserOutputDTO, *internal_error.InternalError)
}

func NewUserUseCase() *UserUseCase {
	return &UserUseCase{}

}

func (u *UserUseCase) FindUserById(ctx context.Context, id string) (*UserOutputDTO, error) {
	userEntity, err := u.UserRepository.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserOutputDTO{
		Id:   userEntity.Id,
		Name: userEntity.Name,
	}, nil
}
