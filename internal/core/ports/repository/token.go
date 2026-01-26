package portRepository

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)



// TokenRepository defines the interface for token data persistence
type TokenRepository interface {
	Create(ctx context.Context, token *domain.Token) error
	FindByAccessToken(ctx context.Context, accessToken string) (*domain.Token, error)
	FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error)
	Update(ctx context.Context, token *domain.Token) error
	Delete(ctx context.Context, id uint) error
	DeleteByUserID(ctx context.Context, userID uint) error
	DeleteExpiredTokens(ctx context.Context) error
}
