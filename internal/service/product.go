package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/repository"
	"github.com/nozzlium/eniqilo_store/internal/util"
)

type ProductService struct {
	repository *repository.ProductRepository
}

func NewProductService(repository *repository.ProductRepository) *ProductService {
	return &ProductService{repository: repository}
}

func (s ProductService) Search(ctx context.Context, query model.SearchProductQuery) ([]model.Product, error) {
	var products []model.Product
	products, err := s.repository.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	return products, nil
}

/*
*
  - "id": "", // should be string
    "name": "",
    "sku": "",
    "category": "",
    "imageUrl": "",
    "stock": 1,
    "notes":"",
    "price":1,
    "location": "",
    "isAvailable": true,
    "createdAt": "" // should in ISO 8601 format
*/
func (s ProductService) Save(ctx context.Context, product model.Product) (string, string, error) {
	now := util.Now()
	id, err := uuid.NewV7()
	if err != nil {
		return "", "", err
	}

	// Check if product already exists from its SKU
	// ignored for now
	// existingProduct, err := s.repository.FindBySKU(ctx, product.SKU)
	// if err != nil && !errors.Is(err, constant.ErrNotFound) {
	// 	return "", "", fmt.Errorf("failed to find product: %v", err)
	// }
	//
	// if existingProduct.SKU == product.SKU {
	// 	return "", "", constant.ErrProductExists
	// }

	product.ID = id
	product.CreatedAt = now
	product.UpdatedAt = now
	product.CreatedBy = uuid.MustParse(ctx.Value("userID").(string))
	err = s.repository.Save(ctx, product)
	if err != nil {
		return "", "", err
	}

	return id.String(), util.ToISO8601(now), nil
}

func (s ProductService) Update(ctx context.Context, id string, product model.Product) error {
	now := util.Now()

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return constant.ErrNotFound
	}

	existingProduct, err := s.repository.FindByID(ctx, uuidID)
	if err != nil {
		if errors.Is(err, constant.ErrNotFound) {
			return err
		}
		return fmt.Errorf("failed to find product: %v", err)
	}

	err = existingProduct.CompareAndUpdate(product)
	if err != nil {
		return err
	}

	product.UpdatedAt = now
	product.UpdatedBy = uuid.MustParse(ctx.Value("userID").(string))
	err = s.repository.Update(ctx, product)
	if err != nil {
		return err
	}

	return nil
}

func (s ProductService) Delete(ctx context.Context, id string) error {
	now := util.Now()

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return constant.ErrNotFound
	}

	deletedAt := now
	deletedBy := uuid.MustParse(ctx.Value("userID").(string))
	err = s.repository.Delete(ctx, uuidID, deletedBy, deletedAt)
	if err != nil {
		return err
	}

	return nil
}
