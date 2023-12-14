package usecase

import (
	"github.com/christianferraz/goexpert/20-CleanArch/internal/entity"
	"github.com/christianferraz/goexpert/20-CleanArch/pkg/events"
)

type OrderOutputDTO struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"final_price"`
}

type ListOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderListed     events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewListOrderUseCase(
	OrderRepository entity.OrderRepositoryInterface,
	OrderCreated events.EventInterface,
	EventDispatcher events.EventDispatcherInterface,
) *ListOrderUseCase {
	return &ListOrderUseCase{
		OrderRepository: OrderRepository,
		OrderListed:     OrderCreated,
		EventDispatcher: EventDispatcher,
	}
}

func (n *ListOrderUseCase) Execute() ([]OrderOutputDTO, error) {
	orders, err := n.OrderRepository.GetOrders()

	if err != nil {
		return []OrderOutputDTO{}, err
	}
	var dtos []OrderOutputDTO
	for _, order := range orders {
		order.CalculateFinalPrice()
		dtos = append(dtos, OrderOutputDTO{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.Price + order.Tax,
		})
	}
	n.OrderListed.SetPayload(dtos)
	n.EventDispatcher.Dispatch(n.OrderListed)

	return dtos, nil
}
