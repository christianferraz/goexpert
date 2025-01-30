package product

import (
	"context"
	"time"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/validator"
	"github.com/google/uuid"
)

type ProductUseCase struct {
	SellerID    uuid.UUID `json:"seller_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	Baseprice   float64   `json:"baseprice"`
	AuctionEnd  time.Time `json:"auction_end"`
}

const minAucttionDuration = time.Hour * 2

func (p ProductUseCase) Valid(ctx context.Context) validator.Evaluator {
	var eval validator.Evaluator
	eval.CheckField(validator.NotBlank(p.ProductName), "product_name", "This field is not empty")
	eval.CheckField(validator.NotBlank(p.Description), "description", "This field is not empty")
	eval.CheckField(validator.MinChars(p.ProductName, 10) && validator.MaxChars(p.ProductName, 255), "product_name", "This field must have a length between 10 and 255 characters")
	eval.CheckField(p.Baseprice > 0, "baseprice", "This field must be greater than 0")
	eval.CheckField(p.AuctionEnd.Sub(time.Now()) >= minAucttionDuration, "auction_end", "This field must be greater than the current date")
	return eval
}
