package main

import (
	"fmt"
	"os"

	"github.com/bordviz/datasphere/internal/config"
	"github.com/bordviz/datasphere/internal/logger"
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
}
