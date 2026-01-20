package portService

import (
	"context"

	"app/xonvera-core/internal/core/domain"
)

// PackageService defines the interface for package operations exposed to handlers.
// Handlers should depend on this service rather than accessing repositories directly.
type PackageService interface {
	GetPackages(ctx context.Context) ([]domain.Package, error)
	GetPackageByID(ctx context.Context, id string) (*domain.Package, error)
}
