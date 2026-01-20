package dto

import (
	"app/xonvera-core/internal/core/domain"
)

// ToUserResponse converts domain.User to UserResponse
func ToUserResponse(user *domain.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToAuthResponse converts domain.AuthResponse to dto.AuthResponse
func ToAuthResponse(domainResp *domain.AuthResponse) *AuthResponse {
	if domainResp == nil {
		return nil
	}
	return &AuthResponse{
		User:         ToUserResponse(domainResp.User),
		AccessToken:  domainResp.AccessToken,
		RefreshToken: domainResp.RefreshToken,
		ExpiresAt:    domainResp.ExpiresAt,
	}
}

// ToRegisterRequest converts dto.RegisterRequest to domain.RegisterRequest
func ToRegisterRequest(req *RegisterRequest) *domain.RegisterRequest {
	if req == nil {
		return nil
	}
	return &domain.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	}
}

// ToLoginRequest converts dto.LoginRequest to domain.LoginRequest
func ToLoginRequest(req *LoginRequest) *domain.LoginRequest {
	if req == nil {
		return nil
	}
	return &domain.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

// ToRefreshTokenRequest converts dto.RefreshTokenRequest to domain.RefreshTokenRequest
func ToRefreshTokenRequest(req *RefreshTokenRequest) *domain.RefreshTokenRequest {
	if req == nil {
		return nil
	}
	return &domain.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}
}
