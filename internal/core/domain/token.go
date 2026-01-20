package domain

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	UserID           uint           `json:"user_id" gorm:"index;not null"`
	AccessToken      string         `json:"-" gorm:"type:text;index;not null"`
	RefreshToken     string         `json:"-" gorm:"type:text;uniqueIndex;not null"`
	ExpiresAt        time.Time      `json:"expires_at" gorm:"not null"`
	RefreshExpiresAt time.Time      `json:"refresh_expires_at" gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relations
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Token) TableName() string {
	return "auth.tokens"
}
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
