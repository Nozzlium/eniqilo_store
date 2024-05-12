package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/util"
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
		ProductDetails: productDetails,
		CreatedAt: util.ToISO8601(
			order.CreatedAt,
		),
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
	return body.Paid > 0 &&
		body.Change >= 0
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
	CreatedAt  string    `query:"createdAt"`
	Limit      int       `query:"limit"`
	Offset     int       `query:"offset"`
	CustomerID uuid.UUID `query:"customerId"`
}

func (soq SearchOrderQuery) BuildWhereClauseAndParams() ([]string, []interface{}) {
	if soq.CustomerID == uuid.Nil {
		return []string{}, []interface{}{}
	}
	return []string{"o.customer_id = $%d"}, []interface{}{soq.CustomerID}
}

func (soq SearchOrderQuery) BuildPagination() (string, []interface{}) {
	return util.DefaultPaginationBuilder(soq.Limit, soq.Offset)
}

func (soq SearchOrderQuery) BuildOrderByClause() []string {
	var sqlClause []string

	if soq.CreatedAt != "" ||
		OrderBy(
			soq.CreatedAt,
		).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf(
				"o.created_at %s",
				soq.CreatedAt,
			),
		)
	} else {
		sqlClause = append(
			sqlClause,
			"o.created_at desc",
		)
	}

	return sqlClause
}

type SearchOrderDetailQuery struct {
	CreatedAt string      `query:"createdAt"`
	IDs       []uuid.UUID `query:"customerId"`
}

func (sodq SearchOrderDetailQuery) BuildWhereClauseAndParams() ([]string, []interface{}) {
	return []string{"o.id = any($%d)"}, []interface{}{sodq.IDs}
}

func (sodq SearchOrderDetailQuery) BuildOrderByClause() []string {
	var sqlClause []string

	if sodq.CreatedAt != "" ||
		OrderBy(
			sodq.CreatedAt,
		).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf(
				"o.created_at %s",
				sodq.CreatedAt,
			),
		)
	} else {
		sqlClause = append(
			sqlClause,
			"o.created_at desc",
		)
	}

	return sqlClause
}
