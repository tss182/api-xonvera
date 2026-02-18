package portService

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshToken(ctx context.Context, req *domain.RefreshTokenRequest) (*domain.AuthResponse, error)
	Logout(ctx context.Context, accessToken string) error
	ValidateAccessToken(ctx context.Context, accessToken string) (uint, error)
}
