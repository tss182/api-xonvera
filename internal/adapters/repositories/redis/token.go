package repositoriesRedis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"app/xonvera-core/internal/core/domain"
	"app/xonvera-core/internal/core/ports/repository"

	"github.com/redis/go-redis/v9"
)

// TokenRedisRepository handles token storage using Redis
type TokenRepository struct {
	client *redis.Client
}

// NewTokenRepository creates a new token Redis repository
func NewTokenRepository(client *redis.Client) portRepository.TokenRepository {
	return &TokenRepository{client: client}
}

// Create saves a token to Redis with expiration
func (r *TokenRepository) Create(ctx context.Context, token *domain.Token) error {
	// Serialize token to JSON
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Store by access token
	accessKey := fmt.Sprintf("token:access:%s", token.AccessToken)
	if err := r.client.Set(ctx, accessKey, data, time.Until(token.ExpiresAt)).Err(); err != nil {
		return fmt.Errorf("failed to store access token: %w", err)
	}

	// Store by refresh token
	refreshKey := fmt.Sprintf("token:refresh:%s", token.RefreshToken)
	if err := r.client.Set(ctx, refreshKey, data, time.Until(token.RefreshExpiresAt)).Err(); err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Store user's token reference
	userKey := fmt.Sprintf("token:user:%d", token.UserID)
	if err := r.client.Set(ctx, userKey, token.AccessToken, time.Until(token.ExpiresAt)).Err(); err != nil {
		return fmt.Errorf("failed to store user token reference: %w", err)
	}

	return nil
}

// FindByAccessToken retrieves a token by access token from Redis
func (r *TokenRepository) FindByAccessToken(ctx context.Context, accessToken string) (*domain.Token, error) {
	key := fmt.Sprintf("token:access:%s", accessToken)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("token not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var token domain.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// FindByRefreshToken retrieves a token by refresh token from Redis
func (r *TokenRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	key := fmt.Sprintf("token:refresh:%s", refreshToken)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, errors.New("token not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	var token domain.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// Update updates an existing token in Redis
func (r *TokenRepository) Update(ctx context.Context, token *domain.Token) error {
	// Delete old tokens first (we'll recreate with new values)
	// This is simpler than trying to update in place

	// Serialize new token data
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Update access token
	accessKey := fmt.Sprintf("token:access:%s", token.AccessToken)
	if err := r.client.Set(ctx, accessKey, data, time.Until(token.ExpiresAt)).Err(); err != nil {
		return fmt.Errorf("failed to update access token: %w", err)
	}

	// Update refresh token
	refreshKey := fmt.Sprintf("token:refresh:%s", token.RefreshToken)
	if err := r.client.Set(ctx, refreshKey, data, time.Until(token.RefreshExpiresAt)).Err(); err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	return nil
}

// Delete removes a token from Redis by finding it first via access token
func (r *TokenRepository) Delete(ctx context.Context, id uint) error {
	// Note: This implementation is constrained by the interface which uses ID
	// In practice, the Logout flow should call this after finding the token
	// The actual deletion is handled in the Logout method which has the token
	// This is a no-op as we'll delete via DeleteByAccessToken helper
	return nil
}

// DeleteByAccessToken is a helper to delete a token by its access token
func (r *TokenRepository) DeleteByAccessToken(ctx context.Context, accessToken string) error {
	// Get the token first to find the refresh token
	token, err := r.FindByAccessToken(ctx, accessToken)
	if err != nil {
		return err
	}

	// Delete access token
	accessKey := fmt.Sprintf("token:access:%s", token.AccessToken)
	if err := r.client.Del(ctx, accessKey).Err(); err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}

	// Delete refresh token
	refreshKey := fmt.Sprintf("token:refresh:%s", token.RefreshToken)
	if err := r.client.Del(ctx, refreshKey).Err(); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	// Delete user reference
	userKey := fmt.Sprintf("token:user:%d", token.UserID)
	if err := r.client.Del(ctx, userKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user token reference: %w", err)
	}

	return nil
}

// DeleteByUserID removes all tokens for a user from Redis
func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	// Get user's token reference
	userKey := fmt.Sprintf("token:user:%d", userID)
	accessToken, err := r.client.Get(ctx, userKey).Result()
	if err == redis.Nil {
		return nil // No tokens found
	}
	if err != nil {
		return fmt.Errorf("failed to get user token: %w", err)
	}

	// Get the token to find refresh token
	token, err := r.FindByAccessToken(ctx, accessToken)
	if err != nil {
		return err
	}

	// Delete access token
	accessKey := fmt.Sprintf("token:access:%s", token.AccessToken)
	if err := r.client.Del(ctx, accessKey).Err(); err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}

	// Delete refresh token
	refreshKey := fmt.Sprintf("token:refresh:%s", token.RefreshToken)
	if err := r.client.Del(ctx, refreshKey).Err(); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	// Delete user reference
	if err := r.client.Del(ctx, userKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user token reference: %w", err)
	}

	return nil
}

// DeleteExpiredTokens is a no-op for Redis as Redis automatically expires keys
func (r *TokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	// Redis automatically handles TTL expiration, so no action needed
	return nil
}
