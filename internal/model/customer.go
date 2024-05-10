package model

import (
	"regexp"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID
	PhoneNumber string
	Name        string
}

type CustomerRegisterBody struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
}

func (body CustomerRegisterBody) IsValid(
	phoneNumberRegex *regexp.Regexp,
) bool {
	if nameLen := len(body.Name); nameLen < 5 ||
		nameLen > 50 {
		return false
	}

	if phoneLen := len(body.PhoneNumber); phoneLen < 10 ||
		phoneLen > 16 {
		return false
	}

	if !phoneNumberRegex.MatchString(
		body.PhoneNumber,
	) {
		return false
	}

	return true
}
