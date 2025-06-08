package chunk

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/bordviz/datasphere/internal/domain/dto"
	"github.com/bordviz/datasphere/internal/domain/models"
	"github.com/bordviz/datasphere/internal/lib/customerror"
	"github.com/bordviz/datasphere/test/storage/suite"
	"github.com/stretchr/testify/require"
)

var st *suite.Suite

func TestInitSuite(t *testing.T) {
	var err error
	st, err = suite.New(t, 3, false)
	require.NoError(t, err)
}

var fakeChunks []*models.Chunk = []*models.Chunk{
	{
		ID:          1,
		FileID:      "testid1",
		ChunkNumber: 1,
		FileKey:     1,
	},
	{
		ID:          2,
		FileID:      "testid2",
		ChunkNumber: 2,
		FileKey:     1,
	},
	{
		ID:          3,
		FileID:      "testid3",
		ChunkNumber: 3,
		FileKey:     1,
	},
	{
		ID:          4,
		FileID:      "testid1",
		ChunkNumber: 1,
		FileKey:     2,
	},
	{
		ID:          5,
		FileID:      "testid2",
		ChunkNumber: 2,
		FileKey:     2,
	},
}

func TestCreateChunk(t *testing.T) {
	tests := []struct {
		name     string
		model    *models.Chunk
		expected int64
		err      error
	}{
		{
			name:     "create chunk 1 for music file",
			model:    fakeChunks[0],
			expected: 1,
		},
		{
			name:     "create chunk 2 for music file",
			model:    fakeChunks[1],
			expected: 2,
		},
		{
			name:     "create chunk 3 for music file",
			model:    fakeChunks[2],
			expected: 3,
		},
		{
			name:     "create chunk 1 for video file",
			model:    fakeChunks[3],
			expected: 4,
		},
		{
			name:     "create chunk 2 for video file",
			model:    fakeChunks[4],
			expected: 5,
		},
		{
			name:  "create invalid chunk with empty file id",
			model: &models.Chunk{ChunkNumber: 1, FileKey: 1},
			err:   customerror.NewCustomError("failed to create new chunk: NOT NULL constraint failed: chunk.file_id", 400),
		},
		{
			name:  "create invalid chunk with empty chunk nember",
			model: &models.Chunk{FileID: "testid", FileKey: 1},
			err:   customerror.NewCustomError("failed to create new chunk: NOT NULL constraint failed: chunk.chunk_number", 400),
		},
		{
			name:  "create invalid chunk with empty file key",
			model: &models.Chunk{FileID: "testid", ChunkNumber: 1},
			err:   customerror.NewCustomError("failed to create new chunk: NOT NULL constraint failed: chunk.file_key", 400),
		},
		{
			name:  "create invalid chunk with not exists file key",
			model: &models.Chunk{FileID: "testid", ChunkNumber: 1, FileKey: 100},
			err:   customerror.NewCustomError("failed to create new chunk: FOREIGN KEY constraint failed", 400),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			tx, err := st.DB.BeginTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  false,
			})
			require.NoError(t, err)
			defer tx.Rollback()

			model := &dto.Chunk{
				FileID:      tt.model.FileID,
				ChunkNumber: tt.model.ChunkNumber,
				FileKey:     tt.model.FileKey,
			}

			id, err := st.Storage.CreateChunk(ctx, tx, model, tt.name)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, id)

			if err == nil {
				err = tx.Commit()
				require.NoError(t, err)
			}
		})
	}
}

func TestGetFileChunks(t *testing.T) {
	tests := []struct {
		name     string
		fileKey  int64
		limit    int64
		offset   int64
		expected *models.FileChunks
		err      error
	}{
		{
			name:    "get valid music file chunks",
			fileKey: 1,
			limit:   10,
			offset:  0,
			expected: &models.FileChunks{
				Count:  3,
				Chunks: fakeChunks[:3],
			},
		},
		{
			name:    "get valid video file valid chunks",
			fileKey: 2,
			limit:   10,
			offset:  0,
			expected: &models.FileChunks{
				Count:  2,
				Chunks: fakeChunks[3:],
			},
		},
		{
			name:    "get valid music file chunks with limit",
			fileKey: 1,
			limit:   2,
			offset:  0,
			expected: &models.FileChunks{
				Count:  3,
				Chunks: fakeChunks[:2],
			},
		},
		{
			name:    "get valid music file chunks with offset",
			fileKey: 1,
			limit:   10,
			offset:  1,
			expected: &models.FileChunks{
				Count:  3,
				Chunks: fakeChunks[1:3],
			},
		},
		{
			name:    "get file chunks for not exists file",
			fileKey: 100,
			limit:   100,
			offset:  0,
			err:     customerror.NewCustomError("file chunks not found", 404),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			tx, err := st.DB.BeginTx(ctx, &sql.TxOptions{
				Isolation: sql.LevelReadCommitted,
				ReadOnly:  true,
			})
			require.NoError(t, err)
			defer tx.Rollback()

			res, err := st.Storage.GetFileChunks(ctx, tx, tt.fileKey, tt.limit, tt.offset, tt.name)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, res)
		})
	}
}
