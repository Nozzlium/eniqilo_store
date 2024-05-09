package model

import "github.com/google/uuid"

type Customer struct {
	ID          uuid.UUID
	PhoneNumber string
	Name        string
}
