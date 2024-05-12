package model

import (
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
