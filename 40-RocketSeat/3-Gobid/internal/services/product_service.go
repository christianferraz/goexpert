package services

import (
	"context"
	"errors"
	"time"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewProductService(pool *pgxpool.Pool) ProductService {
	return ProductService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

func (ps *ProductService) CreateProduct(ctx context.Context, sellerId uuid.UUID, productName, description string, baseprice float64, auctionEnd time.Time) (uuid.UUID, error) {
	idv7, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, err
	}
	args := pgstore.CreateProductsParams{
		ID:          idv7,
		SellerID:    sellerId,
		ProductName: productName,
		Description: description,
		Baseprice:   baseprice,
		AuctionEnd:  auctionEnd,
	}
	id, err := ps.queries.CreateProducts(ctx, args)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil

}

var ErrProductNotFound = errors.New("product not found")

func (ps *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (pgstore.Product, error) {
	product, err := ps.queries.GetProductsById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Product{}, ErrProductNotFound
		}
		return pgstore.Product{}, err
	}
	return product, nil
}
