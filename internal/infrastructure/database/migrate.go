package database

import (
	"errors"

	"app/xonvera-core/internal/infrastructure/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

type Migrator struct {
	migrate *migrate.Migrate
}

func NewMigrator(databaseURL, migrationsPath string) (*Migrator, error) {
	m, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return nil, err
	}
	return &Migrator{migrate: m}, nil
}

func (m *Migrator) Up() error {
	logger.Info("Running database migrations up...")
	err := m.migrate.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if errors.Is(err, migrate.ErrNoChange) {
		logger.Info("No new migrations to apply")
	} else {
		logger.Info("Database migrations applied successfully")
	}
	return nil
}

func (m *Migrator) Down() error {
	logger.Info("Running database migrations down...")
	err := m.migrate.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	logger.Info("Database migrations rolled back successfully")
	return nil
}

func (m *Migrator) Steps(n int) error {
	logger.Info("Running database migration steps", zap.Int("steps", n))
	return m.migrate.Steps(n)
}

func (m *Migrator) Version() (uint, bool, error) {
	return m.migrate.Version()
}

func (m *Migrator) Close() error {
	sourceErr, dbErr := m.migrate.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}
