package portRepository

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

// UserRepository defines the interface for user data persistence
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByPhone(ctx context.Context, phone string) (*domain.User, error)
	FindByEmailOrPhone(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id uint) (*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}
