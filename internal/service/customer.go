package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
)

type CustomerService struct {
	CustomerRepository *repository.CustomerRepository
	Conn               *pgx.Conn
	PhoneNumberRegex   *regexp.Regexp
}

type CustomerData struct {
	UserID      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

func (service *CustomerService) Create(
	ctx context.Context,
	customer model.Customer,
) (model.Customer, error) {
	newID, err := uuid.NewV7()
	if err != nil {
		return model.Customer{}, err
	}

	customer.ID = newID
	saved, err := service.CustomerRepository.Save(
		ctx,
		service.Conn,
		customer,
	)
	if err != nil {
		return model.Customer{}, err
	}

	return saved, nil
}

func (service *CustomerService) FindCustomers(
	ctx context.Context,
	customer model.Customer,
) ([]CustomerData, error) {
	customer.PhoneNumber = validatePhoneNumberWildCard(
		customer.PhoneNumber,
	)

	customerData, err := findAllCustomers(
		ctx,
		service.Conn,
		customer,
	)
	if err != nil {
		return nil, err
	}

	return customerData, nil
}

func validatePhoneNumberWildCard(
	phoneNumber string,
) string {
	fmt.Println(phoneNumber)
	if phoneNumber == "" {
		return ""
	}
	prefix := "+62"
	length := 3
	if len(phoneNumber) < 3 {
		length = len(phoneNumber)
	}
	for i := 0; i < length; i++ {
		if prefix[i] != phoneNumber[i] {
			return ""
		}
	}
	return phoneNumber
}

func findAllCustomers(
	ctx context.Context,
	conn *pgx.Conn,
	customer model.Customer,
) ([]CustomerData, error) {
	query, params := buildQuery(
		customer,
	)

	rows, err := conn.Query(
		ctx,
		query,
		params...)
	if err != nil {
		return []CustomerData{}, err
	}
	defer rows.Close()

	res := make([]CustomerData, 0, 20)
	for rows.Next() {
		tempCust := CustomerData{}
		err := rows.Scan(
			&tempCust.UserID,
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
			customer.PhoneNumber,
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
			"%s;",
			base,
		), params
	}

	return fmt.Sprintf(
		"%s where %s;",
		base,
		strings.Join(
			paramQueries,
			" and ",
		),
	), params
}
