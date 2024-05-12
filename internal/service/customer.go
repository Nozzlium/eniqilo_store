package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
)

type CustomerService struct {
	CustomerRepository *repository.CustomerRepository
}

func NewCustomerService(
	customerRepository *repository.CustomerRepository,
) *CustomerService {
	return &CustomerService{
		CustomerRepository: customerRepository,
	}
}

func (service *CustomerService) Create(
	ctx context.Context,
	customer model.Customer,
) (model.CustomerData, error) {
	savedCustomer, err := service.CustomerRepository.FindByPhoneNumber(
		ctx,
		customer.PhoneNumber,
	)
	if err != nil {
		if !errors.Is(
			err,
			constant.ErrNotFound,
		) {
			return model.CustomerData{}, err
		}
	}

	if savedCustomer.PhoneNumber == customer.PhoneNumber {
		return model.CustomerData{}, constant.ErrConflict
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return model.CustomerData{}, err
	}

	customer.ID = newID
	saved, err := service.CustomerRepository.Save(
		ctx,
		customer,
	)
	if err != nil {
		return model.CustomerData{}, err
	}

	return model.CustomerData{
		UserID:      saved.ID.String(),
		Name:        saved.Name,
		PhoneNumber: saved.PhoneNumber,
	}, nil
}

func (service *CustomerService) FindCustomers(
	ctx context.Context,
	customer model.Customer,
) ([]model.CustomerData, error) {
	customers, err := service.CustomerRepository.FindAllCustomers(
		ctx,
		customer,
	)
	if err != nil {
		return nil, err
	}

	res := make(
		[]model.CustomerData,
		0,
		len(customers),
	)
	for _, customer := range customers {
		res = append(
			res,
			model.CustomerData{
				UserID:      customer.ID.String(),
				PhoneNumber: customer.PhoneNumber,
				Name:        customer.Name,
			},
		)
	}
	return res, nil
}
