package chunk

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/bordviz/datasphere/internal/domain/dto"
	"github.com/bordviz/datasphere/internal/domain/models"
	"github.com/bordviz/datasphere/internal/lib/customerror"
	"github.com/bordviz/datasphere/internal/lib/logger/sl"
	"github.com/bordviz/datasphere/internal/lib/logger/with"
)

func (r *ChunkRequests) CreateChunk(ctx context.Context, tx *sql.Tx, model *dto.Chunk, requestID string) (int64, error) {
	const op = "storage.requests.chunk.CreateChunk"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("create new chunk", slog.Any("model", *model))

	q := `
		INSERT INTO chunk
		(file_id, chunk_number, file)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id int64
	if err := tx.QueryRowContext(
		ctx,
		q,
		model.FileID,
		model.ChunkNumber,
		model.File,
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			log.Error("no rows returned when creating new chunk", sl.Err(err))
			return 0, customerror.NewCustomError("no rows returned when creating new chunk", 500)
		}
		log.Error("failed to create new chunk", sl.Err(err))
		return 0, customerror.NewCustomError(fmt.Sprintf("failed to create new chunk: %s", err.Error()), 400)
	}

	log.Debug("new chunk successfully created", slog.Int64("id", id))
	return id, nil
}

func (r *ChunkRequests) GetFileChunks(ctx context.Context, tx *sql.Tx, file, limit, offset int64, requestID string) (*models.FileChunks, error) {
	const op = "storage.requests.chunk.GetFileChunks"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	q := `
		SELECT COUNT(*) AS count, id, file_id, chunk_number, file
		FROM chunk
		WHERE file = $1
		ORDER BY chunk_number
		LIMIT $2 OFFSET $3;
	`

	var res models.FileChunks
	rows, err := tx.QueryContext(ctx, q, file, limit, offset)
	if err != nil {
		log.Error("failed to get file chunks", sl.Err(err))
		return nil, customerror.NewCustomError(fmt.Sprintf("failed to get file chunks: %s", err.Error()), 400)
	}
	defer rows.Close()

	for rows.Next() {
		var model models.Chunk
		var count int64

		if err := rows.Scan(
			&count,
			&model.FileID,
			&model.ChunkNumber,
			&model.File,
		); err != nil {
			log.Error("failed to get scan file chunk", sl.Err(err))
			return nil, customerror.NewCustomError(fmt.Sprintf("failed to get scan file chunk: %s", err.Error()), 400)
		}

		res.Count = count
		res.Chunks = append(res.Chunks, &model)
	}

	if err := rows.Err(); err != nil {
		log.Error("scan rows error", sl.Err(err))
		return nil, customerror.NewCustomError(fmt.Sprintf("scan rows error: %s", err.Error()), 500)
	}

	if len(res.Chunks) == 0 {
		log.Error("file chunks not found")
		return nil, customerror.NewCustomError("file chunks not found", 404)
	}

	log.Debug("file chunk fetched successfully", slog.Int("count", len(res.Chunks)), slog.Any("chunks", res.Chunks))
	return &res, nil
}
