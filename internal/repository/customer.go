package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
)

type CustomerRepository struct{}

func (repo *CustomerRepository) Save(
	ctx context.Context,
	conn *pgx.Conn,
	customer model.Customer,
) (model.Customer, error) {
	query := `
    insert into customers
      (id, phone_number, name);
    values
      ($1, $2, $3);
  `
	err := conn.QueryRow(
		ctx,
		query,
		customer.ID,
		customer.PhoneNumber,
		customer.Name,
	).Scan()
	if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}
