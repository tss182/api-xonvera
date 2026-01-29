package repositoriesSql

import (
	portRepository "app/xonvera-core/internal/core/ports/repository"

	"gorm.io/gorm"
)

func txDb(tx, db *gorm.DB) *gorm.DB {
	if tx == nil {
		return db
	}
	return tx
}

type txRepository struct {
	db *gorm.DB
}

func NewTxRepository(db *gorm.DB) portRepository.TxRepository {
	return &txRepository{db: db}
}

func (r *txRepository) Begin() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
