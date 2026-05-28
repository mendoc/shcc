package api

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	PublicKey string    `json:"public_key"`
}

type Share struct {
	ID          uuid.UUID `json:"id"`
	Owner       string    `json:"owner"`
	To          string    `json:"to"`
	Credentials string    `json:"credentials"` // JSON chiffré et encodé
	CreatedAt   time.Time `json:"created_at"`
	ExpiredAt   time.Time `json:"expired_at"`
}
