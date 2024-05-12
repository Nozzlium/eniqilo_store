package service

import (
	"context"

	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
)

type OrderService struct {
	orderRepository    *repository.OrderRepository
	productRepository  *repository.ProductRepository
	customerRepository *repository.CustomerRepository
}

func NewOrderService(
	orderRepository *repository.OrderRepository,
	productRepository *repository.ProductRepository,
	customerRepository *repository.CustomerRepository,
) *OrderService {
	return &OrderService{
		orderRepository:    orderRepository,
		productRepository:  productRepository,
		customerRepository: customerRepository,
	}
}

func (service *OrderService) Create(
	ctx context.Context,
	order model.Order,
) (model.Order, error) {
	_, err := service.customerRepository.FindByID(ctx, order.CustomerID)
	if err != nil {
		return model.Order{}, err
	}

	stringIds := make(
		[]string,
		0,
		len(order.ProductOrders),
	)
	for _, orderProduct := range order.ProductOrders {
		stringIds = append(
			stringIds,
			orderProduct.ProductID.String(),
		)
	}
	products, err := service.productRepository.FindByIds(
		ctx,
		stringIds,
	)
	if err != nil {
		return model.Order{}, err
	}

	var actualTotal float64 = 0
	updatedProducts := make(
		[]model.Product,
		0,
		len(order.ProductOrders),
	)

	for i, orderProduct := range order.ProductOrders {
		tempProd, ok := products[orderProduct.ProductID]
		if !ok {
			return model.Order{}, constant.ErrNotFound
		}

		if tempProd.Stock == 0 ||
			tempProd.Stock < orderProduct.Quantity {
			return model.Order{}, constant.ErrInsufficientStock
		}

		itemTotal := float64(
			orderProduct.Quantity,
		) * tempProd.Price

		tempProd.Stock = tempProd.Stock - orderProduct.Quantity
		updatedProducts = append(
			updatedProducts,
			tempProd,
		)

		orderProduct.Price = tempProd.Price
		orderProduct.TotalPrice = itemTotal
		order.ProductOrders[i] = orderProduct

		actualTotal += itemTotal
	}

	if order.PaymentAmount < actualTotal {
		return model.Order{}, constant.ErrInsufficientFund
	}

	if order.PaymentAmount-actualTotal != order.Change {
		return model.Order{}, constant.ErrInvalidChange
	}

	result, err := service.orderRepository.Save(
		ctx,
		order,
		updatedProducts,
	)
	if err != nil {
		return model.Order{}, err
	}

	return result, nil
}

func (service *OrderService) Search(ctx context.Context, query model.SearchOrderQuery) ([]model.Order, error) {
	return nil, nil
}
