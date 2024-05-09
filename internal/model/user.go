package model

import (
	"github.com/google/uuid"
	"github.com/nozzlium/eniqilo_store/internal/util"
)

type RegisterRespose struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

type LoginResponse struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type RegisterBody struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

func (b RegisterBody) IsValid() bool {
	if nameLength := len(b.Name); nameLength < 5 || nameLength > 50 {
		return false
	}

	if phoneLength := len(b.PhoneNumber); phoneLength < 10 || phoneLength > 16 {
		return false
	}

	if !util.ValidatePhoneNumber(b.PhoneNumber) {
		return false
	}

	if passLength := len(b.Password); passLength < 5 || passLength > 15 {
		return false
	}

	return true
}

type LoginBody struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}

func (b LoginBody) IsValid() bool {
	if phoneLength := len(b.PhoneNumber); phoneLength < 10 || phoneLength > 16 {
		return false
	}

	if !util.ValidatePhoneNumber(b.PhoneNumber) {
		return false
	}

	if passLength := len(b.Password); passLength < 5 || passLength > 15 {
		return false
	}

	return true
}

type User struct {
	Name        string
	Password    string
	PhoneNumber string
	ID          uuid.UUID
}
