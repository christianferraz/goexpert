package auction

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/christianferraz/goexpert/26-Leilao/configuration/logger"
	"github.com/christianferraz/goexpert/26-Leilao/internal/entity/auction_entity"
	"github.com/christianferraz/goexpert/26-Leilao/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
)

func (ar *AuctionRepository) FindAuctionById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}
	var auction AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auction); err != nil {
		logger.Error("Error trying to find auction by id", err)
		return nil, internal_error.InternalServerError(fmt.Sprintf("Error trying to find auction by id %s", id))
	}
	return &auction_entity.Auction{
		Id:          auction.Id,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   time.Unix(auction.Timestamp, 0),
	}, nil
}

func (ar *AuctionRepository) FindAuctions(ctx context.Context, status auction_entity.ProductCondition, category, productName string) ([]auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{}
	if status != 0 {
		filter["status"] = status
	}
	if category != "" {
		filter["category"] = strings.ToLower(category)
	}
	if productName != "" {
		filter["productName"] = strings.ToLower(productName)
	}
	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error trying to find auctions", err)
		return nil, internal_error.InternalServerError("Error trying to find auctions")
	}
	defer cursor.Close(ctx)
	var auctionsEntityMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsEntityMongo); err != nil {
		logger.Error("Error trying to find auctions", err)
		return nil, internal_error.InternalServerError("Error trying to find auctions")
	}
	var auctions []auction_entity.Auction
	for _, auction := range auctionsEntityMongo {
		auctions = append(auctions, auction_entity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   auction.Condition,
			Status:      auction.Status,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	}
	return auctions, nil
}
