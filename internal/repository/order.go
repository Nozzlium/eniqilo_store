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
	query := `
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
	_, err := repo.db.Exec(
		ctx,
		query,
		order.ID,
		order.CustomerID,
		order.TotalPrice,
		order.PaymentAmount,
		order.Change,
	)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}
