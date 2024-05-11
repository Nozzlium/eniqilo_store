package repository

import (
	"bytes"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/util"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(db *pgx.Conn) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Search(ctx context.Context, searchQuery model.SearchProductQuery) ([]model.Product, error) {
	var query bytes.Buffer
	query.WriteString(`
    select
			id,
			name,
			sku,
			category,
			image_url,
			stock,
			notes,
			price,
			location, 
			is_available,
			created_at
    from products p 
    where 1=1`)

	queryString, params := util.BuildQueryStringAndParams(
		&query,
		searchQuery.BuildWhereClauseAndParams,
		searchQuery.BuildPagination,
		searchQuery.BuildOrderByClause,
	)

	var products []model.Product
	rows, err := r.db.Query(ctx, queryString, params...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return products, nil
		}
		return nil, err
	}

	for rows.Next() {
		var (
			p        model.Product
			category string
		)

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.SKU,
			&category,
			&p.ImageURL,
			&p.Stock,
			&p.Notes,
			&p.Price,
			&p.Location,
			&p.IsAvailable,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		p.Category = p.Category.FromDBEnumType(category)
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) Save(ctx context.Context, product model.Product) error {
	query := `
  insert into products (
    id,
    name,
    price,
    stock,
    notes,
    category,
    image_url
  ) values 
  ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
		product.Notes,
		product.Category.ToDBEnumType(),
		product.ImageURL,
	)

	if err != nil {
		return err
	}

	return nil
}
