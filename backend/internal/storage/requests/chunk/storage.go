package chunk

import "log/slog"

type ChunkRequests struct {
	log *slog.Logger
}

func New(log *slog.Logger) *ChunkRequests {
	return &ChunkRequests{log: log}
}
