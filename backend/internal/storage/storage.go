package storage

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/bordviz/datasphere/internal/domain/dto"
	"github.com/bordviz/datasphere/internal/domain/models"
	"github.com/bordviz/datasphere/internal/storage/requests/chunk"
	"github.com/bordviz/datasphere/internal/storage/requests/file"
)

type Storage struct {
	log *slog.Logger
	FileStorage
	ChunkStorage
}

type FileStorage interface {
	CreateFile(ctx context.Context, tx *sql.Tx, model *dto.File, requestID string) (int64, error)
	GetFileByID(ctx context.Context, tx *sql.Tx, id int64, requestID string) (*models.File, error)
	SearchFiles(ctx context.Context, tx *sql.Tx, query, requestID string) ([]*models.File, error)
	DeleteFile(ctx context.Context, tx *sql.Tx, id int64, requestID string) (int64, error)
}

type ChunkStorage interface {
	CreateChunk(ctx context.Context, tx *sql.Tx, model *dto.Chunk, requestID string) (int64, error)
	GetFileChunks(ctx context.Context, tx *sql.Tx, file, limit, offset int64, requestID string) (*models.FileChunks, error)
}

func New(log *slog.Logger) *Storage {
	return &Storage{
		log:          log,
		FileStorage:  file.New(log),
		ChunkStorage: chunk.New(log),
	}
}
