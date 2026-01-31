package portService

import (
	"context"

	"app/xonvera-core/internal/adapters/dto"
)

// AuthService defines the interface for authentication operations
type AuthService interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, accessToken string) error
	ValidateAccessToken(ctx context.Context, accessToken string) (uint, error)
}
