package domain

type (
	RegisterRequest struct {
		Name     string `json:"name" validate:"required,min=2,max=100"`
		Phone    string `json:"phone" validate:"required,min=10,max=15"`
		Password string `json:"password" validate:"required,min=6,max=100"`
	}

	LoginRequest struct {
		Username string `json:"username" validate:"required"` // email or phone
		Password string `json:"password" validate:"required"`
	}

	UserResponse struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		CreatedAt string `json:"created_at"`
	}

	AuthResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
	}
)
