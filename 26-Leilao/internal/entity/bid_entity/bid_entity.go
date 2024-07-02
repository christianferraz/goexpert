package bid_entity

import "time"

type Bid struct {
	Id        string    `bson:"_id,omitempty"`
	UserId    string    `bson:"user_id"`
	AuctionId string    `bson:"auction_id"`
	Amount    float64   `bson:"amount"`
	Timestamp time.Time `bson:"timestamp"`
}
