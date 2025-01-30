package services

import (
	"context"
	"errors"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/store/pgstore"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BidsService struct {
	pool    *pgxpool.Pool
	queries *pgstore.Queries
}

func NewBidsService(pool *pgxpool.Pool) BidsService {
	return BidsService{
		pool:    pool,
		queries: pgstore.New(pool),
	}
}

var ErrBidIsTooLow = errors.New("the bid value is too low")

func (bs *BidsService) Placebid(ctx context.Context, productID, BidderID uuid.UUID, amount float64) (pgstore.Bid, error) {
	product, err := bs.queries.GetProductsById(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}
	higuestBid, err := bs.queries.GetHighestBidByProductID(ctx, productID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return pgstore.Bid{}, err
		}
	}
	if product.Baseprice >= amount || higuestBid.BidAmount >= amount {
		return pgstore.Bid{}, ErrBidIsTooLow
	}
	higuestBid, err = bs.queries.CreateBid(ctx, pgstore.CreateBidParams{
		ProductID: productID,
		BidderID:  BidderID,
		BidAmount: amount,
	})
	if err != nil {
		return pgstore.Bid{}, err
	}
	return higuestBid, nil
}

func (bs *BidsService) CreateBid(ctx context.Context, productID, bidderID uuid.UUID, bidAmount float64) error {
	args := pgstore.CreateBidParams{
		ProductID: productID,
		BidderID:  bidderID,
		BidAmount: bidAmount,
	}
	_, err := bs.queries.CreateBid(ctx, args)
	if err != nil {
		return err
	}
	return nil
}
