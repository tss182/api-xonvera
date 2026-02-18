package domain

import (
	"time"
)

type Token struct {
	ID               uint
	UserID           uint
	AccessToken      string
	RefreshToken     string
	ExpiresAt        time.Time
	RefreshExpiresAt time.Time
	Timestamp

	// Relations
	User *User
}

func (Token) TableName() string {
	return "auth.tokens"
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
