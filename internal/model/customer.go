package model

import (
	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/util"
)

type Customer struct {
	Name        string
	PhoneNumber string
	ID          uuid.UUID
}

type CustomerData struct {
	UserID      string `json:"userId"`
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

type CustomerRegisterBody struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

func (body CustomerRegisterBody) IsValid() bool {
	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return false
	}

	if phoneLen := len(body.PhoneNumber); phoneLen < 10 ||
		phoneLen > 16 {
		return false
	}

	if !util.ValidatePhoneNumber(
		body.PhoneNumber,
	) {
		return false
	}

	return true
}
