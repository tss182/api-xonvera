package portRepository

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

type PackageRepository interface {
	GetAll(ctx context.Context) ([]domain.Package, error)
	GetByID(ctx context.Context, id string) (*domain.Package, error)
	Create(ctx context.Context, pkg *domain.Package) error
	Update(ctx context.Context, pkg *domain.Package) error
	Delete(ctx context.Context, id string) error
}
