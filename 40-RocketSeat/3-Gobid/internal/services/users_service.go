package services

import (
	"context"
	"errors"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

var (
	ErrDuplicatedEmailOrUsername = errors.New("duplicated or email already exists")
	ErrInvalidCredentials        = errors.New("invalid credentials")
)

func NewUserService(pool *pgxpool.Pool) UserService {
	return UserService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (us *UserService) CreateUser(ctx context.Context, userName, email, password, bio string) (uuid.UUID, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return uuid.Nil, err
	}
	idv7, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}
	args := pgstore.CreateUserParams{
		ID:           idv7,
		UserName:     userName,
		Email:        email,
		PasswordHash: hash,
		Bio:          bio,
	}
	id, err := us.queries.CreateUser(ctx, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return uuid.Nil, ErrDuplicatedEmailOrUsername
		}
		return uuid.Nil, err
	}
	return id, nil
}

func (us *UserService) AuthenticateUser(ctx context.Context, email, password string) (uuid.UUID, error) {
	user, err := us.queries.GetUserByUsername(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, ErrInvalidCredentials
		}
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return uuid.UUID{}, err
	}
	return user.ID, nil
}
