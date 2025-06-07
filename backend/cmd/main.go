package main

import (
	"fmt"
	"os"

	"github.com/bordviz/datasphere/internal/config"
	"github.com/bordviz/datasphere/internal/lib/logger/sl"
	"github.com/bordviz/datasphere/internal/logger"
	"github.com/bordviz/datasphere/internal/storage"
	"github.com/bordviz/datasphere/internal/storage/migrations"
	"github.com/bordviz/datasphere/internal/storage/sqlite"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.Env)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Debug("Debug messages are available")
	log.Info("Info messages are available")
	log.Warn("Warn messages are available")
	log.Error("Error messages are available")

	db, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to connect to database", sl.Err(err))
		os.Exit(1)
	}

	migrationsHandler, err := migrations.NewMigrationHandler(cfg.StoragePath, cfg.MigrationsPath)
	if err != nil {
		log.Error("failed to create migrations handler", sl.Err(err))
		os.Exit(1)
	}

	if err := migrationsHandler.Up(); err != nil {
		log.Error("failed to apply migrations", sl.Err(err))
		os.Exit(1)
	}

	storage := storage.New(log)
	_ = db
	_ = storage

}
