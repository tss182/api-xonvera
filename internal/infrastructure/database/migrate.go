package database

import (
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
	log.Println("Running migrations up...")
	err := m.migrate.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No new migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}
	return nil
}

func (m *Migrator) Down() error {
	log.Println("Running migrations down...")
	err := m.migrate.Down()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	log.Println("Migrations rolled back successfully")
	return nil
}

func (m *Migrator) Steps(n int) error {
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
