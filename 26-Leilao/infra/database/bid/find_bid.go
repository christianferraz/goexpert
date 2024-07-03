package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/christianferraz/goexpert/26-Leilao/configuration/logger"
	"github.com/christianferraz/goexpert/26-Leilao/internal/entity/bid_entity"
	"github.com/christianferraz/goexpert/26-Leilao/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (bd *BidRepository) FindBidByAuctionId(ctx context.Context, auctionId string) ([]bid_entity.Bid, error) {
	cursor, err := bd.Collection.Find(ctx, bson.M{"auctionId": auctionId})
	if err != nil {
		logger.Error(fmt.Sprintf("Error trying to find bids by auctionid %s", auctionId), err)
		return nil, internal_error.InternalServerError(fmt.Sprintf("Error trying to find bids by auctionid %s", err))
	}
	defer cursor.Close(ctx)
	var bidEntitiesMongo []BidEntityMongo
	if err := cursor.All(ctx, &bidEntitiesMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to decode bids %s", auctionId), err)
		return nil, internal_error.InternalServerError(fmt.Sprintf("Error trying to decode bids %s", auctionId))
	}
	var bidEntities []bid_entity.Bid
	for _, bidEntityMongo := range bidEntitiesMongo {
		bidEntities = append(bidEntities, bid_entity.Bid{
			Id:        bidEntityMongo.Id,
			AuctionId: bidEntityMongo.AuctionId,
			UserId:    bidEntityMongo.UserId,
			Amount:    bidEntityMongo.Amount,
			Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
		})
	}
	return bidEntities, nil
}

func (bd *BidRepository) FindWinnerBidByAuctionId(ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auctionId": auctionId}
	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})
	var bidEntityMongo BidEntityMongo
	if err := bd.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find winner bid by auctionid %s", auctionId), err)
		return nil, internal_error.InternalServerError(fmt.Sprintf("Error trying to find winner bid by auctionid %s", auctionId))
	}
	return &bid_entity.Bid{
		Id:        bidEntityMongo.Id,
		AuctionId: bidEntityMongo.AuctionId,
		UserId:    bidEntityMongo.UserId,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
