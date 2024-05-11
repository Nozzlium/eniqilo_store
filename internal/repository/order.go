package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
)

type OrderRepository struct {
	db *pgx.Conn
}

func (repo *OrderRepository) Save(
	ctx context.Context,
	order model.Order,
) (model.Order, error) {
	tx, err := repo.db.BeginTx(
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

	queryOrder := `
    insert into
      orders (
        id,
        customer_id,
        total_price,
        payment_amount,
        change
    ) values (
      $1, $2, $3, $4, $5
    );
  `
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
	_, err = tx.Exec(
		ctx,
		queryOrder,
		order.ID,
		order.CustomerID,
		order.TotalPrice,
		order.PaymentAmount,
		order.Change,
	)
	if err != nil {
		return model.Order{}, err
	}

	orderProductBatch := &pgx.Batch{}
	for _, orderProduct := range order.ProductOrders {
		orderProductBatch.Queue(
			queryOrderProduct,
			orderProduct.OrderID,
			orderProduct.ProductID,
			orderProduct.Quantity,
			orderProduct.Price,
			orderProduct.TotalPrice,
		)
	}
	batchRes := tx.SendBatch(
		ctx,
		orderProductBatch,
	)
	err = batchRes.Close()
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
