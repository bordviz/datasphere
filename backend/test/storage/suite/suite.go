package suite

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/bordviz/datasphere/internal/config"
	"github.com/bordviz/datasphere/internal/logger"
	"github.com/bordviz/datasphere/internal/storage"
	"github.com/bordviz/datasphere/internal/storage/migrations"
	"github.com/bordviz/datasphere/internal/storage/sqlite"
)

type Suite struct {
	DB      *sql.DB
	Storage *storage.Storage
}

const configPath = "config/test.yml"

func New(t *testing.T, upLevel int) (*Suite, error) {
	t.Helper()

	cfg, err := config.LoadConfigFromPath(strings.Repeat("../", upLevel) + configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create new config: %s", err.Error())
	}

	log, err := logger.New(cfg.Env)
	if err != nil {
		return nil, fmt.Errorf("failed to create new logger: %s", err.Error())
	}

	db, err := sqlite.New(strings.Repeat("../", upLevel) + cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err.Error())
	}

	migrationsHandler, err := migrations.NewMigrationHandler(strings.Repeat("../", upLevel)+cfg.StoragePath, strings.Repeat("../", upLevel)+cfg.MigrationsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrations handler: %s", err.Error())
	}

	if err := migrationsHandler.Down(); err != nil {
		return nil, fmt.Errorf("failed to revert migrations: %s", err.Error())
	}

	if err := migrationsHandler.Up(); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %s", err.Error())
	}

	storage := storage.New(log)

	return &Suite{
		DB:      db,
		Storage: storage,
	}, nil
}
