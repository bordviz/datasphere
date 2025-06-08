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
	"github.com/bordviz/datasphere/internal/lib/nullable"
)

func (r *ChunkRequests) CreateChunk(ctx context.Context, tx *sql.Tx, model *dto.Chunk, requestID string) (int64, error) {
	const op = "storage.requests.chunk.CreateChunk"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("create new chunk", slog.Any("model", *model))

	q := `
		INSERT INTO chunk
		(file_id, chunk_number, file_key)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var id int64
	if err := tx.QueryRowContext(
		ctx,
		q,
		nullable.IsNullable(model.FileID),
		nullable.IsNullable(model.ChunkNumber),
		nullable.IsNullable(model.FileKey),
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

func (r *ChunkRequests) getFileChunksCount(ctx context.Context, tx *sql.Tx, fileKey int64, requestID string) (int64, error) {
	const op = "storage.requests.chunk.getFileChunksCount"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("get file chunks count", slog.Int64("file_key", fileKey))

	q := `
		SELECT COUNT(*)
		FROM chunk
		WHERE file_key = $1;
	`

	var count int64
	if err := tx.QueryRowContext(ctx, q, fileKey).Scan(&count); err != nil {
		log.Error("failed to get file chunks count", sl.Err(err))
		return 0, customerror.NewCustomError(fmt.Sprintf("failed to get file chunks count: %s", err.Error()), 400)
	}

	if count == 0 {
		log.Error("file chunks not found")
		return 0, customerror.NewCustomError("file chunks not found", 404)
	}

	log.Debug("file chunks count fetched successfully", slog.Int64("count", count))
	return count, nil
}

func (r *ChunkRequests) GetFileChunks(ctx context.Context, tx *sql.Tx, fileKey, limit, offset int64, requestID string) (*models.FileChunks, error) {
	const op = "storage.requests.chunk.GetFileChunks"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("get file chunks", slog.Int64("file_key", fileKey), slog.Int64("limit", limit), slog.Int64("offset", offset))

	q := `
		SELECT id, file_id, chunk_number, file_key
		FROM chunk
		WHERE file_key = $1
		ORDER BY chunk_number
		LIMIT $2 OFFSET $3;
	`

	var res models.FileChunks
	count, err := r.getFileChunksCount(ctx, tx, fileKey, requestID)
	if err != nil {
		return nil, err
	}
	res.Count = count

	rows, err := tx.QueryContext(ctx, q, fileKey, limit, offset)
	if err != nil {
		log.Error("failed to get file chunks", sl.Err(err))
		return nil, customerror.NewCustomError(fmt.Sprintf("failed to get file chunks: %s", err.Error()), 400)
	}
	defer rows.Close()

	for rows.Next() {
		var model models.Chunk

		if err := rows.Scan(
			&model.ID,
			&model.FileID,
			&model.ChunkNumber,
			&model.FileKey,
		); err != nil {
			log.Error("failed to get scan file chunk", sl.Err(err))
			return nil, customerror.NewCustomError(fmt.Sprintf("failed to get scan file chunk: %s", err.Error()), 400)
		}
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
