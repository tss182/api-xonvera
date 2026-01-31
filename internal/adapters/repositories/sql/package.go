package repositoriesSql

import (
	"context"

	"app/xonvera-core/internal/core/domain"
	portRepository "app/xonvera-core/internal/core/ports/repository"

	"gorm.io/gorm"
)

type packageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) portRepository.PackageRepository {
	return &packageRepository{db: db}
}

func (r *packageRepository) GetAll(ctx context.Context) ([]domain.Package, error) {
	var packages []domain.Package
	if err := r.db.WithContext(ctx).Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

func (r *packageRepository) GetByID(ctx context.Context, id string) (*domain.Package, error) {
	var pkg domain.Package
	if err := r.db.WithContext(ctx).First(&pkg, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &pkg, nil
}

func (r *packageRepository) Create(ctx context.Context, pkg *domain.Package) error {
	return r.db.WithContext(ctx).Create(pkg).Error
}

func (r *packageRepository) Update(ctx context.Context, pkg *domain.Package) error {
	return r.db.WithContext(ctx).Save(pkg).Error
}

func (r *packageRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Package{}, "id = ?", id).Error
}
