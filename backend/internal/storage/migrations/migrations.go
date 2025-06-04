package migrations

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MigrationsHandler struct {
	migrator *migrate.Migrate
}

func NewMigrationHandler(storagePath, migrationsPath string) (*MigrationsHandler, error) {
	dsn := fmt.Sprintf("sqlite3://%s", storagePath)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new migrations handler: %s", err.Error())
	}
	return &MigrationsHandler{
		migrator: m,
	}, nil
}

func (m *MigrationsHandler) Up() error {
	if err := m.migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations error: %s", err.Error())
	}
	return nil
}

func (m *MigrationsHandler) Down() error {
	if err := m.migrator.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("revert migrations error: %s", err.Error())
	}
	return nil
}
