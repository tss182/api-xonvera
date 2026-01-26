package services

import (
	"testing"
	"time"

	"app/xonvera-core/internal/infrastructure/config"
)

func TestTokenService_GenerateTokenPair_Success(t *testing.T) {
	cfg := &config.TokenConfig{
		SecretKey:      "a_very_secret_key_that_is_long_enough",
		Expired:        time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewTokenService(cfg)

	tokenPair, err := service.GenerateTokenPair(1)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if tokenPair == nil {
		t.Error("expected token pair, got nil")
	}
	if tokenPair.AccessToken == "" {
		t.Error("expected access token, got empty string")
	}
	if tokenPair.RefreshToken == "" {
		t.Error("expected refresh token, got empty string")
	}
	if tokenPair.ExpiresAt == 0 {
		t.Error("expected expires_at, got 0")
	}
}

func TestTokenService_ValidateToken_Success(t *testing.T) {
	cfg := &config.TokenConfig{
		SecretKey:      "a_very_secret_key_that_is_long_enough",
		Expired:        time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewTokenService(cfg)

	// Generate a token
	tokenPair, err := service.GenerateTokenPair(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Validate the token
	userID, err := service.ValidateToken(tokenPair.AccessToken)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if userID != 1 {
		t.Errorf("expected user ID 1, got %d", userID)
	}
}

func TestTokenService_ValidateToken_InvalidToken(t *testing.T) {
	cfg := &config.TokenConfig{
		SecretKey:      "a_very_secret_key_that_is_long_enough",
		Expired:        time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewTokenService(cfg)

	// Try to validate invalid token
	_, err := service.ValidateToken("invalid_token")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTokenService_ValidateRefreshToken_Success(t *testing.T) {
	cfg := &config.TokenConfig{
		SecretKey:      "a_very_secret_key_that_is_long_enough",
		Expired:        time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewTokenService(cfg)

	// Generate a token
	tokenPair, err := service.GenerateTokenPair(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Validate the refresh token
	userID, err := service.ValidateRefreshToken(tokenPair.RefreshToken)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if userID != 1 {
		t.Errorf("expected user ID 1, got %d", userID)
	}
}

func TestTokenService_ValidateRefreshToken_InvalidToken(t *testing.T) {
	cfg := &config.TokenConfig{
		SecretKey:      "a_very_secret_key_that_is_long_enough",
		Expired:        time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewTokenService(cfg)

	// Try to validate invalid refresh token
	_, err := service.ValidateRefreshToken("invalid_refresh_token")

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestTokenService_TokenExpiration(t *testing.T) {
	cfg := &config.TokenConfig{
		SecretKey:      "a_very_secret_key_that_is_long_enough",
		Expired:        time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewTokenService(cfg)

	tokenPair, err := service.GenerateTokenPair(1)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	// Check that expires_at is in the future
	now := time.Now().Unix()
	if tokenPair.ExpiresAt <= now {
		t.Errorf("expected expires_at to be in future, got %d (now: %d)", tokenPair.ExpiresAt, now)
	}

	// Check that it's approximately 60 minutes in the future
	expectedTime := now + (60 * 60) // 60 minutes
	tolerance := int64(5)           // 5 second tolerance
	if tokenPair.ExpiresAt < expectedTime-tolerance || tokenPair.ExpiresAt > expectedTime+tolerance {
		t.Errorf("expected expires_at around %d, got %d", expectedTime, tokenPair.ExpiresAt)
	}
}

func TestTokenService_DifferentUsersGetDifferentTokens(t *testing.T) {
	cfg := &config.TokenConfig{
		SecretKey:      "a_very_secret_key_that_is_long_enough",
		Expired:        time.Hour * 1,
		RefreshExpired: time.Hour * 24 * 7,
	}

	service := NewTokenService(cfg)

	// Generate tokens for different users
	tokenPair1, err := service.GenerateTokenPair(1)
	if err != nil {
		t.Fatalf("failed to generate token for user 1: %v", err)
	}

	tokenPair2, err := service.GenerateTokenPair(2)
	if err != nil {
		t.Fatalf("failed to generate token for user 2: %v", err)
	}

	// Tokens should be different
	if tokenPair1.AccessToken == tokenPair2.AccessToken {
		t.Error("tokens for different users should be different")
	}
	if tokenPair1.RefreshToken == tokenPair2.RefreshToken {
		t.Error("refresh tokens for different users should be different")
	}

	// But validate should return correct user IDs
	userID1, _ := service.ValidateToken(tokenPair1.AccessToken)
	userID2, _ := service.ValidateToken(tokenPair2.AccessToken)

	if userID1 != 1 {
		t.Errorf("expected user ID 1, got %d", userID1)
	}
	if userID2 != 2 {
		t.Errorf("expected user ID 2, got %d", userID2)
	}
}
