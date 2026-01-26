package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"app/xonvera-core/internal/core/domain"
	"app/xonvera-core/internal/core/ports/service"
	"app/xonvera-core/internal/infrastructure/config"

	"github.com/o1egl/paseto"
)

type tokenService struct {
	paseto          *paseto.V2
	secretKey       []byte
	expireIn        time.Duration
	refreshExpireIn time.Duration
}

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

type TokenPayload struct {
	UserID    uint      `json:"user_id"`
	Type      string    `json:"type"` // "access" or "refresh"
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func NewTokenService(cfg *config.TokenConfig) portService.TokenService {
	return &tokenService{
		paseto:          paseto.NewV2(),
		secretKey:       []byte(cfg.SecretKey),
		expireIn:        cfg.Expired,
		refreshExpireIn: cfg.RefreshExpired,
	}
}

func (s *tokenService) GenerateTokenPair(userID uint) (*domain.TokenPair, error) {
	now := time.Now()
	expiresAt := now.Add(s.expireIn)

	// Generate access token
	accessPayload := TokenPayload{
		UserID:    userID,
		Type:      tokenTypeAccess,
		IssuedAt:  now,
		ExpiredAt: expiresAt,
	}

	symmetricKey := make([]byte, 32)
	copy(symmetricKey, s.secretKey)

	accessToken, err := s.paseto.Encrypt(symmetricKey, accessPayload, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token (random string + paseto for validation)
	refreshPayload := TokenPayload{
		UserID:    userID,
		Type:      tokenTypeRefresh,
		IssuedAt:  now,
		ExpiredAt: now.Add(s.refreshExpireIn),
	}

	refreshToken, err := s.paseto.Encrypt(symmetricKey, refreshPayload, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Add random suffix to make refresh token unique
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	refreshToken = refreshToken + "." + hex.EncodeToString(randomBytes)

	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt.Unix(),
	}, nil
}

func (s *tokenService) ValidateToken(token string) (uint, error) {
	var payload TokenPayload

	symmetricKey := make([]byte, 32)
	copy(symmetricKey, s.secretKey)

	err := s.paseto.Decrypt(token, symmetricKey, &payload, nil)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	if payload.Type != tokenTypeAccess {
		return 0, fmt.Errorf("invalid token type")
	}

	if time.Now().After(payload.ExpiredAt) {
		return 0, fmt.Errorf("token has expired")
	}

	return payload.UserID, nil
}

func (s *tokenService) ValidateRefreshToken(token string) (uint, error) {
	// Remove random suffix
	var pasetoToken string
	for i := len(token) - 1; i >= 0; i-- {
		if token[i] == '.' {
			pasetoToken = token[:i]
			break
		}
	}
	if pasetoToken == "" {
		return 0, fmt.Errorf("invalid refresh token format")
	}

	var payload TokenPayload

	symmetricKey := make([]byte, 32)
	copy(symmetricKey, s.secretKey)

	err := s.paseto.Decrypt(pasetoToken, symmetricKey, &payload, nil)
	if err != nil {
		return 0, fmt.Errorf("invalid refresh token: %w", err)
	}

	if payload.Type != tokenTypeRefresh {
		return 0, fmt.Errorf("invalid token type")
	}

	if time.Now().After(payload.ExpiredAt) {
		return 0, fmt.Errorf("refresh token has expired")
	}

	return payload.UserID, nil
}
