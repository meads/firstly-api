package domain

import (
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	UserID       int64     `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	IsRevoked    bool      `json:"isRevoked"`
	ExpiresAt    time.Time `json:"expiresAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

type CreateSessionParams struct {
	ID           string    `json:"id"`
	UserID       int64     `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	IsRevoked    bool      `json:"isRevoked"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
