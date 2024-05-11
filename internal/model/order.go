package model

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID
	CustomerID    uuid.UUID
	TotalPrice    big.Float
	PaymentAmount big.Float
	Change        big.Float
	ProductOrders []ProductOrder
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type ProductOrder struct {
	OrderID    uuid.UUID
	ProductID  uuid.UUID
	Quantity   uint64
	Price      big.Float
	TotalPrice big.Float
}
