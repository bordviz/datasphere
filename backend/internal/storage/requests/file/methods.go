package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/bordviz/datasphere/internal/domain/dto"
	"github.com/bordviz/datasphere/internal/domain/models"
	"github.com/bordviz/datasphere/internal/lib/customerror"
	"github.com/bordviz/datasphere/internal/lib/logger/sl"
	"github.com/bordviz/datasphere/internal/lib/logger/with"
	"github.com/bordviz/datasphere/internal/lib/nullable"
)

func (r *FileRequests) CreateFile(ctx context.Context, tx *sql.Tx, model *dto.File, requestID string) (int64, error) {
	const op = "storage.requests.file.CreateFile"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("create new file", slog.Any("model", *model))

	q := `
		INSERT INTO file
		(filename, size, name, chunks_count)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`

	var id int64
	if err := tx.QueryRowContext(
		ctx,
		q,
		nullable.IsNullable(model.Filename),
		nullable.IsNullable(model.Size),
		nullable.IsNullable(model.Name),
		nullable.IsNullable(model.ChunksCount),
	).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			log.Error("no rows returned when creating new file", sl.Err(err))
			return 0, customerror.NewCustomError("no rows returned when creating new file", 500)
		}
		log.Error("failed to create new file", sl.Err(err))
		return 0, customerror.NewCustomError(fmt.Sprintf("failed to create new file: %s", err.Error()), 400)
	}

	log.Debug("new file successfully created", slog.Int64("id", id))
	return id, nil
}

func (r *FileRequests) GetFileByID(ctx context.Context, tx *sql.Tx, id int64, requestID string) (*models.File, error) {
	const op = "storage.requests.file.GetFileByID"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("get file by id", slog.Int64("id", id))

	q := `
		SELECT id, filename, size, name, chunks_count 
		FROM file
		WHERE id = $1;
	`

	var model models.File
	if err := tx.QueryRowContext(
		ctx,
		q,
		id,
	).Scan(
		&model.ID,
		&model.Filename,
		&model.Size,
		&model.Name,
		&model.ChunksCount,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Error("file not found")
			return nil, customerror.NewCustomError("file not found", 404)
		}

		log.Error("failed to get file by id", sl.Err(err))
		return nil, customerror.NewCustomError(fmt.Sprintf("failed to get file by id: %s", err.Error()), 400)
	}

	log.Debug("file by id fetched successfully", slog.Any("model", model))
	return &model, nil
}

func (r *FileRequests) SearchFiles(ctx context.Context, tx *sql.Tx, query, requestID string) ([]*models.File, error) {
	const op = "storage.requests.file.SearchFile"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("search file", slog.String("query", query))

	q := `
		SELECT id, filename, size, name, chunks_count
		FROM file
		WHERE LOWER(filename) LIKE $1
		OR LOWER(name) LIKE $1;
	`

	rows, err := tx.QueryContext(ctx, q, "%"+strings.ToLower(query)+"%")
	if err != nil {
		log.Error("failed to get files", sl.Err(err))
		return nil, customerror.NewCustomError(fmt.Sprintf("failed to get files: %s", err.Error()), 400)
	}
	defer rows.Close()

	var res []*models.File
	for rows.Next() {
		var model models.File
		if err := rows.Scan(
			&model.ID,
			&model.Filename,
			&model.Size,
			&model.Name,
			&model.ChunksCount,
		); err != nil {
			log.Error("failed to scan model", sl.Err(err))
			return nil, customerror.NewCustomError(fmt.Sprintf("failed to scan file: %s", err.Error()), 500)
		}

		res = append(res, &model)
	}

	if err := rows.Err(); err != nil {
		log.Error("rows error", sl.Err(err))
		return nil, customerror.NewCustomError(fmt.Sprintf("rows error: %s", err.Error()), 500)
	}

	if len(res) == 0 {
		log.Error("files not found")
		return nil, customerror.NewCustomError("files not found", 404)
	}

	log.Debug("files fetched successfully", slog.Int("count", len(res)), slog.Any("models", res))
	return res, nil
}

func (r *FileRequests) DeleteFile(ctx context.Context, tx *sql.Tx, id int64, requestID string) (int64, error) {
	const op = "storage.requests.file.DeleteFile"

	log := with.WithOpAndRequestID(r.log, op, requestID)

	log.Debug("delete file", slog.Int64("id", id))

	q := `
		DELETE FROM file
		WHERE id = $1
		RETURNING id;
	`

	var fileID int64
	if err := tx.QueryRowContext(ctx, q, id).Scan(&fileID); err != nil {
		if err == sql.ErrNoRows {
			log.Error("file not found")
			return 0, customerror.NewCustomError("file not found", 404)
		}

		log.Error("failed to delete file", sl.Err(err))
		return 0, customerror.NewCustomError(fmt.Sprintf("failed to delete file: %s", err.Error()), 400)
	}

	log.Debug("file successfully deleted", slog.Int64("id", fileID))
	return fileID, nil
}
