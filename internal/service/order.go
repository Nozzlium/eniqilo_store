package service

import (
	"context"

	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
)

type OrderService struct {
	OrderRepository repository.OrderRepository
}

func (service *OrderService) Create(
	ctx context.Context,
	order model.Order,
) (model.Order, error) {
  :
}
