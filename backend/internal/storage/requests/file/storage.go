package file

import (
	"log/slog"
)

type FileRequests struct {
	log *slog.Logger
}

func New(log *slog.Logger) *FileRequests {
	return &FileRequests{log: log}
}
