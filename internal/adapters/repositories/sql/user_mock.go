package repositoriesSql

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

// MockUserRepository is a mock implementation of portRepository.UserRepository for testing
type MockUserRepository struct {
	CreateFunc             func(ctx context.Context, user *domain.User) error
	FindByEmailFunc        func(ctx context.Context, email string) (*domain.User, error)
	FindByPhoneFunc        func(ctx context.Context, phone string) (*domain.User, error)
	FindByEmailOrPhoneFunc func(ctx context.Context, username string) (*domain.User, error)
	FindByIDFunc           func(ctx context.Context, id uint) (*domain.User, error)
	ExistsByEmailFunc      func(ctx context.Context, email string) (bool, error)
	ExistsByPhoneFunc      func(ctx context.Context, phone string) (bool, error)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	if m.FindByPhoneFunc != nil {
		return m.FindByPhoneFunc(ctx, phone)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByEmailOrPhone(ctx context.Context, username string) (*domain.User, error) {
	if m.FindByEmailOrPhoneFunc != nil {
		return m.FindByEmailOrPhoneFunc(ctx, username)
	}
	return nil, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	if m.ExistsByEmailFunc != nil {
		return m.ExistsByEmailFunc(ctx, email)
	}
	return false, nil
}

func (m *MockUserRepository) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	if m.ExistsByPhoneFunc != nil {
		return m.ExistsByPhoneFunc(ctx, phone)
	}
	return false, nil
}
