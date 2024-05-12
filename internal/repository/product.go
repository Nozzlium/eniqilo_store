package repository

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/constant"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/util"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(
	db *pgx.Conn,
) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Search(
	ctx context.Context,
	searchQuery model.SearchProductQuery,
) ([]model.Product, error) {
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
    where deleted_at is null`)

	queryString, params := util.BuildQueryStringAndParams(
		&query,
		searchQuery.BuildWhereClauseAndParams,
		searchQuery.BuildPagination,
		searchQuery.BuildOrderByClause,
	)

	var products []model.Product
	rows, err := r.db.Query(
		ctx,
		queryString,
		params...)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
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

		p.Category = p.Category.FromDBEnumType(
			category,
		)
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) Save(
	ctx context.Context,
	product model.Product,
) error {
	query := `
  insert into products (
    id,
    name,
    sku,
    price,
    stock,
    notes,
    category,
    image_url,
    is_available,
    location,
    created_at,
    updated_at,
    created_by
  ) values 
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`

	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.Name,
		product.SKU,
		product.Price,
		product.Stock,
		product.Notes,
		product.Category.ToDBEnumType(),
		product.ImageURL,
		product.IsAvailable,
		product.Location,
		product.CreatedAt,
		product.UpdatedAt,
		product.CreatedBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Update(
	ctx context.Context,
	product model.Product,
) error {
	query := `
  update products set
    name = $1,
    sku = $2,
    price = $3,
    stock = $4,
    notes = $5,
    category = $6,
    image_url = $7,
    is_available = $8,
    location = $9,
    updated_at = $10,
    updated_by = $11
  where id = $12`

	_, err := r.db.Exec(ctx, query,
		product.Name,
		product.SKU,
		product.Price,
		product.Stock,
		product.Notes,
		product.Category.ToDBEnumType(),
		product.ImageURL,
		product.IsAvailable,
		product.Location,
		product.UpdatedAt,
		product.UpdatedBy,
		product.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) Delete(
	ctx context.Context,
	id, deletedBy uuid.UUID,
	deletedAt time.Time,
) error {
	query := `
  update products set
    deleted_at = $1,
    deleted_by = $2
  where id = $3`

	_, err := r.db.Exec(ctx, query,
		deletedAt,
		deletedBy,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *ProductRepository) FindBySKU(
	ctx context.Context,
	sku string,
) (model.Product, error) {
	var (
		p        model.Product
		category string
	)

	query := `select
			id,
			name,
			sku,
			category,
			image_url,
			stock,
			notes,
			price,
			location, 
			is_available
    from products p 
    where sku = $1 and deleted_at is null`

	row := r.db.QueryRow(
		ctx,
		query,
		sku,
	)
	err := row.Scan(
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
	)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return model.Product{}, constant.ErrNotFound
		}
		return p, err
	}

	p.Category = p.Category.FromDBEnumType(
		category,
	)

	return p, nil
}

func (r *ProductRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (model.Product, error) {
	var (
		p        model.Product
		category string
	)

	query := `select
			name,
			sku,
			category,
			image_url,
			stock,
			notes,
			price,
			location, 
			is_available
    from products p 
    where id = $1 and deleted_at is null`

	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(
		&p.Name,
		&p.SKU,
		&category,
		&p.ImageURL,
		&p.Stock,
		&p.Notes,
		&p.Price,
		&p.Location,
		&p.IsAvailable,
	)
	if err != nil {
		return p, err
	}

	p.Category = p.Category.FromDBEnumType(
		category,
	)

	return p, nil
}

func (r *ProductRepository) FindByIds(
	ctx context.Context,
	ids []uuid.UUID,
) (map[uuid.UUID]model.Product, error) {
	query := `
    select 
      id, 
      name, 
      stock, 
      price, 
      is_available 
    from products
    where id = any($1::uuid[])
    and deleted_at is null
  `
	rows, err := r.db.Query(
		ctx,
		query,
		ids,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make(
		map[uuid.UUID]model.Product,
	)
	for rows.Next() {
		temp := model.Product{}
		err = rows.Scan(
			&temp.ID,
			&temp.Name,
			&temp.Stock,
			&temp.Price,
			&temp.IsAvailable,
		)
		if err != nil {
			return nil, err
		}

		res[temp.ID] = temp
	}

	return res, nil
}
