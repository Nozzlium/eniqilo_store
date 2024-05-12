package repository

import (
	"bytes"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/nozzlium/eniqilo_store/internal/model"
	"github.com/nozzlium/eniqilo_store/internal/util"
)

type OrderRepository struct {
	db *pgx.Conn
}

func NewOrderRepository(
	db *pgx.Conn,
) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) Save(
	ctx context.Context,
	order model.Order,
	products []model.Product,
) (model.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return model.Order{}, err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	// create order entity
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
	batch.Queue(
		queryOrder,
		order.ID,
		order.CustomerID,
		order.TotalPrice,
		order.PaymentAmount,
		order.Change,
	)

	// create order_product entity
	// total price not included for insersion
	// column setup to generated always stored
	queryOrderProduct := `
    insert into
      order_product (
        order_id,
        product_id,
        quantity,
        price
    ) values (
      $1, $2, $3, $4
    )
  `
	for _, orderProduct := range order.ProductOrders {
		batch.Queue(
			queryOrderProduct,
			orderProduct.OrderID,
			orderProduct.ProductID,
			orderProduct.Quantity,
			orderProduct.Price,
		)
	}

	// update product stock
	queryUpdateStock := `
    update products 
    set stock = $1 
    where id = $2;
  `
	for _, product := range products {
		batch.Queue(
			queryUpdateStock,
			product.Stock,
			product.ID,
		)
	}

	batchRes := tx.SendBatch(
		ctx,
		batch,
	)
	if err := batchRes.Close(); err != nil {
		return model.Order{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) Search(
	ctx context.Context,
	searchQuery model.SearchOrderQuery,
) (map[uuid.UUID]model.Order, []uuid.UUID, error) {
	var query bytes.Buffer
	query.WriteString(`
    select
      o.id,
      o.customer_id,
      o.payment_amount,
      o.change,
      o.created_at,
      op.product_id,
      op.quantity,
    from orders o
    join order_product op on o.id = op.order_id
    where 1=1`)

	queryString, params := util.BuildQueryStringAndParams(
		&query,
		searchQuery.BuildWhereClauseAndParams,
		searchQuery.BuildPagination,
		searchQuery.BuildOrderByClause,
	)

	var (
		orders   []uuid.UUID
		orderMap = make(
			map[uuid.UUID]model.Order,
		)
	)
	rows, err := r.db.Query(
		ctx,
		queryString,
		params...)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return orderMap, orders, nil
		}
		return orderMap, orders, err
	}

	for rows.Next() {
		var (
			o         model.Order
			quantity  int
			productID uuid.UUID
		)

		err := rows.Scan(
			&o.ID,
			&o.CustomerID,
			&o.PaymentAmount,
			&o.Change,
			&o.CreatedAt,
			&productID,
			&quantity,
		)
		if err != nil {
			return orderMap, orders, err
		}

		om, ok := orderMap[o.ID]
		if !ok {
			orders = append(
				orders,
				o.ID,
			)
			o.ProductOrders = make(
				[]model.ProductOrder,
				0,
			)
		} else {
			o = om
		}

		o.ProductOrders = append(
			o.ProductOrders,
			model.ProductOrder{
				OrderID:   o.ID,
				ProductID: productID,
				Quantity:  quantity,
			},
		)

		orderMap[o.ID] = o
	}

	return orderMap, orders, nil
}
