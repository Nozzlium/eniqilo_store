package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
)

type OrderService struct {
	OrderRepository   repository.OrderRepository
	ProductRepository repository.ProductRepository
}

func (service *OrderService) Create(
	ctx context.Context,
	order model.Order,
) (model.Order, error) {
	tx, err := service.OrderRepository.DB.BeginTx(
		ctx,
		pgx.TxOptions{},
	)
	if err != nil {
		return model.Order{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// TODO: Still needs to check on user validity, code not merged

	queryOrderProduct := `
    insert into
      orders (
        order_id,
        product_id,
        quantity,
        price,
        total_price
    ) values (
      $1, $2, $3, $4, $5
    )
  `
	orderProductBatch := &pgx.Batch{}
	updatedProds := make(
		[]model.Product,
		0,
		len(order.ProductOrders),
	)
	var actualTotalPrice float64 = 0
	for _, orderProduct := range order.ProductOrders {
		temp, err := service.ProductRepository.FindById(
			ctx,
			orderProduct.ProductID,
		)
		if err != nil {
			return model.Order{}, err
		}

		if temp.Stock < int(
			orderProduct.Quantity,
		) {
			return model.Order{}, model.ErrInsufficientStock
		}

		if !temp.IsAvailable {
			return model.Order{}, model.ErrInsufficientStock
		}

		actualPrice := float64(
			orderProduct.Quantity,
		) * temp.Price

		orderProduct.Price = temp.Price
		orderProduct.TotalPrice = actualPrice
		actualTotalPrice += actualPrice

		orderProductBatch.Queue(
			queryOrderProduct,
			orderProduct.OrderID,
			orderProduct.ProductID,
			orderProduct.Quantity,
			orderProduct.Price,
			orderProduct.TotalPrice,
		)
		updatedProds = append(
			updatedProds,
			temp,
		)
	}

	if order.PaymentAmount < actualTotalPrice {
		return model.Order{}, model.ErrInsufficientFund
	}

	if actualChange := actualTotalPrice - order.PaymentAmount; actualChange != order.Change {
		return model.Order{}, model.ErrInvalidChange
	}

	order.TotalPrice = actualTotalPrice
	_, err = service.OrderRepository.SaveTx(
		ctx,
		tx,
		order,
	)
	if err != nil {
		return model.Order{}, err
	}

	batchRes := tx.SendBatch(
		ctx,
		orderProductBatch,
	)
	err = batchRes.Close()
	if err != nil {
		return model.Order{}, err
	}

	for _, updatedProduct := range updatedProds {
		_, err = service.ProductRepository.SaveTx(
			ctx,
			tx,
			updatedProduct,
		)
		if err != nil {
			return model.Order{}, err
		}
	}

	return order, nil
}
