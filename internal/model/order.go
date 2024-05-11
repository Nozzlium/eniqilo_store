package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID
	CustomerID    uuid.UUID
	TotalPrice    float64
	PaymentAmount float64
	Change        float64
	ProductOrders []ProductOrder
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
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
	Quantity   uint64
	Price      float64
	TotalPrice float64
}

type OrderRequestBody struct {
	CustomerID     string              `json:"customerId"`
	Paid           float64             `json:"paid"`
	Change         float64             `json:"change"`
	ProductDetails []ProductDetailBody `json:"productDetails"`
}

func (body OrderRequestBody) IsValid() bool {
	if body.Paid < 1 {
		return false
	}

	for _, product := range body.ProductDetails {
		if !product.isValid() {
			return false
		}
	}

	return true
}

type ProductDetailBody struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

func (body ProductDetailBody) isValid() bool {
	if body.Quantity < 1 {
		return false
	}

	return true
}

type OrderResponseBody struct {
	TransactionId  string              `json:"transactionId"`
	CustomerID     string              `json:"customerId"`
	Paid           float64             `json:"paid"`
	Change         float64             `json:"change"`
	CreatedAt      string              `json:"createdAt"`
	ProductDetails []ProductDetailBody `json:"productDetails"`
}
