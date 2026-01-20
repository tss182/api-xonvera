package services

import (
	"context"

	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"
	portService "app/xonvera-core/internal/core/ports/service"
)

// PackageService is the concrete implementation of portService.PackageService.
// It orchestrates package-related operations and encapsulates repository access.
type PackageService struct {
	repo portRepository.PackageRepository
}

// NewPackageService builds a PackageService.
func NewPackageService(repo portRepository.PackageRepository) portService.PackageService {
	return &PackageService{repo: repo}
}

// GetPackages retrieves all packages.
func (s *PackageService) GetPackages(ctx context.Context) ([]domain.Package, error) {
	return s.repo.GetAll(ctx)
}

// GetPackageByID retrieves a package by its ID.
func (s *PackageService) GetPackageByID(ctx context.Context, id string) (*domain.Package, error) {
	return s.repo.GetByID(ctx, id)
}
