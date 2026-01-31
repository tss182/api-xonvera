package repositoriesSql

import (
	portRepository "app/xonvera-core/internal/core/ports/repository"

	"gorm.io/gorm"
)

// gormTransaction wraps *gorm.DB to implement the Transaction interface
type gormTransaction struct {
	db *gorm.DB
}

// Commit commits the transaction
func (t *gormTransaction) Commit() error {
	return t.db.Commit().Error
}

// Rollback rolls back the transaction
func (t *gormTransaction) Rollback() error {
	return t.db.Rollback().Error
}

// GetDB returns the underlying *gorm.DB for internal repository use
func (t *gormTransaction) GetDB() *gorm.DB {
	return t.db
}

// txDb is a helper that returns the appropriate DB handle
// If tx is a gormTransaction, it returns the wrapped DB
// Otherwise, it returns the default db
func txDb(tx portRepository.Transaction, db *gorm.DB) *gorm.DB {
	if tx == nil {
		return db
	}
	if gormTx, ok := tx.(*gormTransaction); ok {
		return gormTx.db
	}
	return db
}

type txRepository struct {
	db *gorm.DB
}

func NewTxRepository(db *gorm.DB) portRepository.TxRepository {
	return &txRepository{db: db}
}

func (r *txRepository) Begin() (portRepository.Transaction, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &gormTransaction{db: tx}, nil
}
