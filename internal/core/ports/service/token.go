package portService

import (
	"app/xonvera-core/internal/core/domain"
)

// TokenService defines the interface for token operations
type TokenService interface {
	GenerateTokenPair(userID uint) (*domain.TokenPair, error)
	ValidateToken(token string) (uint, error)
	ValidateRefreshToken(token string) (uint, error)
}
