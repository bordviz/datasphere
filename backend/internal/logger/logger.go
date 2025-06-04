package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bordviz/datasphere/internal/lib/logger/slogpretty"
)

func New(env string) (*slog.Logger, error) {
	log := new(slog.Logger)

	switch env {
	case "test, disable":
		log = slog.New(slog.DiscardHandler)
	case "local":
		log = setupPrettyLogger()
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	default:
		return nil, fmt.Errorf("failed to create new logger, invalid environment: %s, expected: disable, local, dev, prod", env)
	}

	return log, nil
}

func setupPrettyLogger() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOptions: &slog.HandlerOptions{Level: slog.LevelDebug},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
