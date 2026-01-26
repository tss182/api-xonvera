package repositoriesRedis

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

// MockTokenRepository is a mock implementation of portRepository.TokenRepository for testing
type MockTokenRepository struct {
	CreateFunc              func(ctx context.Context, token *domain.Token) error
	FindByAccessTokenFunc   func(ctx context.Context, accessToken string) (*domain.Token, error)
	FindByRefreshTokenFunc  func(ctx context.Context, refreshToken string) (*domain.Token, error)
	UpdateFunc              func(ctx context.Context, token *domain.Token) error
	DeleteFunc              func(ctx context.Context, id uint) error
	DeleteByUserIDFunc      func(ctx context.Context, userID uint) error
	DeleteExpiredTokensFunc func(ctx context.Context) error
}

func (m *MockTokenRepository) Create(ctx context.Context, token *domain.Token) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, token)
	}
	return nil
}

func (m *MockTokenRepository) FindByAccessToken(ctx context.Context, accessToken string) (*domain.Token, error) {
	if m.FindByAccessTokenFunc != nil {
		return m.FindByAccessTokenFunc(ctx, accessToken)
	}
	return nil, nil
}

func (m *MockTokenRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	if m.FindByRefreshTokenFunc != nil {
		return m.FindByRefreshTokenFunc(ctx, refreshToken)
	}
	return nil, nil
}

func (m *MockTokenRepository) Update(ctx context.Context, token *domain.Token) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, token)
	}
	return nil
}

func (m *MockTokenRepository) Delete(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockTokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	if m.DeleteByUserIDFunc != nil {
		return m.DeleteByUserIDFunc(ctx, userID)
	}
	return nil
}

func (m *MockTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	if m.DeleteExpiredTokensFunc != nil {
		return m.DeleteExpiredTokensFunc(ctx)
	}
	return nil
}
