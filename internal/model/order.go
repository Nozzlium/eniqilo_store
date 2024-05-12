package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	ProductOrders []ProductOrder
	ID            uuid.UUID
	CustomerID    uuid.UUID
	TotalPrice    float64
	PaymentAmount float64
	Change        float64
}

func (order Order) ToResponseBody() OrderResponseBody {
	productDetails := make(
		[]ProductDetailBody,
		0,
		len(order.ProductOrders),
	)
	for _, product := range order.ProductOrders {
		productDetails = append(
			productDetails,
			ProductDetailBody{
				ProductID: product.ProductID.String(),
				Quantity: int(
					product.Quantity,
				),
			},
		)
	}

	return OrderResponseBody{
		TransactionId:  order.ID.String(),
		CustomerID:     order.CustomerID.String(),
		Paid:           order.PaymentAmount,
		Change:         order.Change,
		CreatedAt:      order.CreatedAt.String(),
		ProductDetails: productDetails,
	}
}

type ProductOrder struct {
	OrderID    uuid.UUID
	ProductID  uuid.UUID
	Quantity   int
	Price      float64
	TotalPrice float64
}

type OrderRequestBody struct {
	CustomerID     string              `json:"customerId"`
	ProductDetails []ProductDetailBody `json:"productDetails"`
	Paid           float64             `json:"paid"`
	Change         float64             `json:"change"`
}

func (body OrderRequestBody) IsValid() bool {
	return body.Paid > 0 && body.Change >= 0
}

type ProductDetailBody struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (body ProductDetailBody) IsValid() bool {
	return body.Quantity > 0
}

type OrderResponseBody struct {
	CreatedAt      string              `json:"createdAt"`
	TransactionId  string              `json:"transactionId"`
	CustomerID     string              `json:"customerId"`
	ProductDetails []ProductDetailBody `json:"productDetails"`
	Paid           float64             `json:"paid"`
	Change         float64             `json:"change"`
}

type SearchOrderQuery struct {
	CustomerID string `query:"customerId"`
	CreatedAt  string `query:"createdAt"`
	Limit      int    `query:"limit"`
	Offset     int    `query:"offset"`
}

func (soq SearchOrderQuery) BuildWhereClauseAndParams() ([]string, []interface{}) {
	var (
		sqlClause []string
		params    []interface{}
	)

	if soq.CustomerID != "" {
		params = append(params, soq.CustomerID)
		sqlClause = append(
			sqlClause,
			"o.customer_id = $%d",
		)
	}

	return sqlClause, params
}

func (soq SearchOrderQuery) BuildPagination() (string, []interface{}) {
	var params []interface{}

	limit := 5
	offset := 0
	if soq.Limit > 0 {
		limit = soq.Limit
	}
	if soq.Offset > 0 {
		offset = soq.Offset
	}
	params = append(
		params,
		limit,
		offset,
	)

	return "limit $%d offset $%d", params
}

func (soq SearchOrderQuery) BuildOrderByClause() []string {
	var sqlClause []string

	if soq.CreatedAt != "" ||
		OrderBy(
			soq.CreatedAt,
		).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf("o.created_at %s", soq.CreatedAt),
		)
	} else {
		sqlClause = append(
			sqlClause,
			"o.created_at desc",
		)
	}

	return sqlClause
}
