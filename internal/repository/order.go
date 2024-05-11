package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
)

type OrderRepository struct {
	DB *pgx.Conn
}

func (r *OrderRepository) Save(
	ctx context.Context,
	order model.Order,
	products []model.Product,
) (model.Order, error) {
	tx, err := r.DB.Begin(
		ctx,
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
	batch := &pgx.Batch{}
	batch.Queue(
		queryOrder,
		order.ID,
		order.CustomerID,
		order.TotalPrice,
		order.PaymentAmount,
		order.Change,
	)

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
	for _, orderProduct := range order.ProductOrders {
		batch.Queue(
			queryOrderProduct,
			orderProduct.OrderID,
			orderProduct.ProductID,
			orderProduct.Quantity,
			orderProduct.Price,
			orderProduct.TotalPrice,
		)
	}

	queryUpdateStock := `
    update products 
    set quantity = $1 
    where id = $2;
  `
	for _, product := range products {
		batch.Queue(
			queryUpdateStock,
			product.Stock,
			product.ID,
		)
	}

	batchRes := tx.SendBatch(ctx, batch)
	if err := batchRes.Close(); err != nil {
		return model.Order{}, err
	}

	return order, nil
}
