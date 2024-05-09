package model

import "github.com/google/uuid"

type LoginResponse struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

type User struct {
	Name        string
	Password    string
	PhoneNumber string
	ID          uuid.UUID
}
