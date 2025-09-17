package model

import "github.com/google/uuid"

type Location struct {
	UserId      uuid.UUID `json:"userId"`
	DisplayName string    `json:"display_name"`
	Latitude    float64   `json:"lat"`
	Longitude   float64   `json:"lon"`
}
