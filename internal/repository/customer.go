package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
)

type CustomerRepository struct {
	db *pgx.Conn
}

func NewCustomerRepository(
	db *pgx.Conn,
) *CustomerRepository {
	return &CustomerRepository{
		db: db,
	}
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
	_, err := repo.db.Exec(
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
	err := repo.db.QueryRow(
		ctx,
		query,
		phoneNumber,
	).Scan(&customer.ID, &customer.PhoneNumber, &customer.Name)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return customer, constant.ErrNotFound
		}
		return customer, err
	}

	return customer, nil
}

func (r *CustomerRepository) FindAllCustomers(
	ctx context.Context,
	customer model.Customer,
) ([]model.Customer, error) {
	query, params := buildQuery(
		customer,
	)

	rows, err := r.db.Query(
		ctx,
		query,
		params...)
	if err != nil {
		return []model.Customer{}, err
	}
	defer rows.Close()

	res := make([]model.Customer, 0, 20)
	for rows.Next() {
		tempCust := model.Customer{}
		err := rows.Scan(
			&tempCust.ID,
			&tempCust.PhoneNumber,
			&tempCust.Name,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, tempCust)
	}

	return res, nil
}

func buildQuery(
	customer model.Customer,
) (string, []any) {
	paramCounter := 0
	paramQueries := make([]string, 0, 2)
	params := make([]any, 0, 2)
	base := `
    select 
      id, 
      phone_number, 
      name
    from customers
  `

	if customer.PhoneNumber != "" {
		paramCounter++
		paramQueries = append(
			paramQueries,
			fmt.Sprintf(
				`phone_number like $%d || '%%'`,
				paramCounter,
			),
		)
		params = append(
			params,
			"+"+customer.PhoneNumber,
		)
	}

	if customer.Name != "" {
		paramCounter++
		paramQueries = append(
			paramQueries,
			fmt.Sprintf(
				"name ilike '%%' || $%d || '%%'",
				paramCounter,
			),
		)
		params = append(
			params,
			customer.Name,
		)
	}

	if paramCounter == 0 {
		return fmt.Sprintf(
			"%s order by created_at desc;",
			base,
		), params
	}

	return fmt.Sprintf(
		"%s where %s order by created_at desc;",
		base,
		strings.Join(
			paramQueries,
			" and ",
		),
	), params
}
