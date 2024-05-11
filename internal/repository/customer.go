package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
)

type CustomerRepository struct {
	DB *pgx.Conn
}

func NewCustomerRepository(
	db *pgx.Conn,
) (*CustomerRepository, error) {
	if db == nil {
		return &CustomerRepository{}, errors.New(
			"cannot init, db is nil",
		)
	}

	return &CustomerRepository{
		DB: db,
	}, nil
}

func (repo *CustomerRepository) Save(
	ctx context.Context,
	customer model.Customer,
) (model.Customer, error) {
	query := `
    insert into customers
      (id, phone_number, name)
    values
      ($1, $2, $3);
  `
	_, err := repo.DB.Exec(
		ctx,
		query,
		customer.ID,
		customer.PhoneNumber,
		customer.Name,
	)
	if err != nil {
		return model.Customer{}, err
	}

	return customer, nil
}

func (repo *CustomerRepository) FindByPhoneNumber(
	ctx context.Context,
	phoneNumber string,
) (model.Customer, error) {
	query := `
    select 
      id, phone_number, name
    from customers
      where phone_number = $1
  `

	customer := model.Customer{}
	err := repo.DB.QueryRow(
		ctx,
		query,
		phoneNumber,
	).Scan(&customer.ID, &customer.PhoneNumber, &customer.Name)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return customer, model.ErrNotFound
		}
		return customer, err
	}

	return customer, nil
}
