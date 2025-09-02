package model

import "github.com/google/uuid"

type User struct {
	UserId    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
}
