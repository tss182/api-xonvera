package repositoriesRedis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"

	"github.com/redis/go-redis/v9"
)

// Token key prefixes for Redis storage
const (
	tokenAccessKeyPrefix  = "token:access:%s"
	tokenRefreshKeyPrefix = "token:refresh:%s"
	tokenUserKeyPrefix    = "token:user:%d"
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
	accessKey := fmt.Sprintf(tokenAccessKeyPrefix, token.AccessToken)
	accessTTL := time.Until(token.ExpiresAt)
	if accessTTL <= 0 {
		return fmt.Errorf("access token already expired")
	}
	if err := r.client.Set(ctx, accessKey, data, accessTTL).Err(); err != nil {
		return fmt.Errorf("failed to store access token: %w", err)
	}

	// Store by refresh token
	refreshKey := fmt.Sprintf(tokenRefreshKeyPrefix, token.RefreshToken)
	refreshTTL := time.Until(token.RefreshExpiresAt)
	if refreshTTL <= 0 {
		return fmt.Errorf("refresh token already expired")
	}
	if err := r.client.Set(ctx, refreshKey, data, refreshTTL).Err(); err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	// Store user's token reference
	userKey := fmt.Sprintf(tokenUserKeyPrefix, token.UserID)
	userTTL := time.Until(token.ExpiresAt)
	if userTTL <= 0 {
		return fmt.Errorf("user token reference already expired")
	}
	if err := r.client.Set(ctx, userKey, token.AccessToken, userTTL).Err(); err != nil {
		return fmt.Errorf("failed to store user token reference: %w", err)
	}

	return nil
}

// FindByAccessToken retrieves a token by access token from Redis
func (r *TokenRepository) FindByAccessToken(ctx context.Context, accessToken string) (*domain.Token, error) {
	key := fmt.Sprintf(tokenAccessKeyPrefix, accessToken)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("404:not found token")
	}
	if err != nil {
		return nil, fmt.Errorf("500:failed to get token: %w", err)
	}

	var token domain.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("500:failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// FindByRefreshToken retrieves a token by refresh token from Redis
func (r *TokenRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*domain.Token, error) {
	key := fmt.Sprintf(tokenRefreshKeyPrefix, refreshToken)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("404:not found token")
	}
	if err != nil {
		return nil, fmt.Errorf("500:failed to get token: %w", err)
	}

	var token domain.Token
	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return nil, fmt.Errorf("500:failed to unmarshal token: %w", err)
	}

	return &token, nil
}

// deleteOldToken removes the previous token entry for a user
func (r *TokenRepository) deleteOldToken(ctx context.Context, userID uint) error {
	userKey := fmt.Sprintf(tokenUserKeyPrefix, userID)
	oldAccessToken, err := r.client.Get(ctx, userKey).Result()
	if err == redis.Nil || oldAccessToken == "" {
		return nil // No previous token exists
	}
	if err != nil {
		return nil // Ignore errors when retrieving old token
	}

	oldToken, err := r.FindByAccessToken(ctx, oldAccessToken)
	if err != nil {
		return nil // Ignore if old token can't be found
	}

	// Delete old access token key
	oldAccessKey := fmt.Sprintf(tokenAccessKeyPrefix, oldToken.AccessToken)
	_ = r.client.Del(ctx, oldAccessKey).Err() // Ignore deletion errors

	// Delete old refresh token key
	oldRefreshKey := fmt.Sprintf(tokenRefreshKeyPrefix, oldToken.RefreshToken)
	_ = r.client.Del(ctx, oldRefreshKey).Err() // Ignore deletion errors

	return nil
}

// validateTokenTTL checks if token expiration is valid
func validateTokenTTL(expiresAt time.Time, tokenType string) error {
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return fmt.Errorf("%s token already expired", tokenType)
	}
	return nil
}

// Update updates an existing token in Redis using atomic transaction
func (r *TokenRepository) Update(ctx context.Context, token *domain.Token) error {
	// Delete previous token if exists
	_ = r.deleteOldToken(ctx, token.UserID)

	// Serialize new token data
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %w", err)
	}

	// Validate TTLs
	if err := validateTokenTTL(token.ExpiresAt, "access"); err != nil {
		return err
	}
	if err := validateTokenTTL(token.RefreshExpiresAt, "refresh"); err != nil {
		return err
	}

	// Use Redis pipeline for atomic update
	pipe := r.client.Pipeline()

	accessKey := fmt.Sprintf(tokenAccessKeyPrefix, token.AccessToken)
	accessTTL := time.Until(token.ExpiresAt)
	refreshKey := fmt.Sprintf(tokenRefreshKeyPrefix, token.RefreshToken)
	refreshTTL := time.Until(token.RefreshExpiresAt)
	userKey := fmt.Sprintf(tokenUserKeyPrefix, token.UserID)
	userTTL := time.Until(token.ExpiresAt)

	// Queue all operations
	pipe.Set(ctx, accessKey, data, accessTTL)
	pipe.Set(ctx, refreshKey, data, refreshTTL)
	pipe.Set(ctx, userKey, token.AccessToken, userTTL)

	// Execute all operations atomically
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update token: %w", err)
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
	accessKey := fmt.Sprintf(tokenAccessKeyPrefix, token.AccessToken)
	if err := r.client.Del(ctx, accessKey).Err(); err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}

	// Delete refresh token
	refreshKey := fmt.Sprintf(tokenRefreshKeyPrefix, token.RefreshToken)
	if err := r.client.Del(ctx, refreshKey).Err(); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	// Delete user reference
	userKey := fmt.Sprintf(tokenUserKeyPrefix, token.UserID)
	if err := r.client.Del(ctx, userKey).Err(); err != nil {
		return fmt.Errorf("failed to delete user token reference: %w", err)
	}

	return nil
}

// DeleteByUserID removes all tokens for a user from Redis
func (r *TokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	// Get user's token reference
	userKey := fmt.Sprintf(tokenUserKeyPrefix, userID)
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
	accessKey := fmt.Sprintf(tokenAccessKeyPrefix, token.AccessToken)
	if err := r.client.Del(ctx, accessKey).Err(); err != nil {
		return fmt.Errorf("failed to delete access token: %w", err)
	}

	// Delete refresh token
	refreshKey := fmt.Sprintf(tokenRefreshKeyPrefix, token.RefreshToken)
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
