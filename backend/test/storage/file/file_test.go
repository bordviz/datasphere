package file_test

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
	st, err = suite.New(t, 3, true)
	require.NoError(t, err)
}

var fakeFiles []*models.File = []*models.File{
	{
		ID:          1,
		Filename:    "fake_music.mp3",
		Size:        1024,
		Name:        "Music file",
		ChunksCount: 3,
	},
	{
		ID:          2,
		Filename:    "fake_video.mp4",
		Size:        2048,
		Name:        "Video file",
		ChunksCount: 7,
	},
	{
		ID:          3,
		Filename:    "fake_txt.txt",
		Size:        512,
		Name:        "Text file",
		ChunksCount: 1,
	},
	{
		ID:          4,
		Filename:    "fake_zip.zip",
		Size:        4096,
		Name:        "ZIP file",
		ChunksCount: 33,
	},
}

func TestCreateFile(t *testing.T) {
	tests := []struct {
		name     string
		model    *models.File
		expected int64
		err      error
	}{
		{
			name:     "create valid misic file",
			model:    fakeFiles[0],
			expected: 1,
		},
		{
			name:     "create valid video file",
			model:    fakeFiles[1],
			expected: 2,
		},
		{
			name:     "create valid text file",
			model:    fakeFiles[2],
			expected: 3,
		},
		{
			name:     "create valid ZIP file",
			model:    fakeFiles[3],
			expected: 4,
		},
		{
			name:  "create file with empty filename",
			model: &models.File{Size: 1024, Name: "Error file", ChunksCount: 10},
			err:   customerror.NewCustomError("failed to create new file: NOT NULL constraint failed: file.filename", 400),
		},
		{
			name:  "create file with empty name",
			model: &models.File{Size: 1024, Filename: "Error file", ChunksCount: 10},
			err:   customerror.NewCustomError("failed to create new file: NOT NULL constraint failed: file.name", 400),
		},
		{
			name:  "create file with empty chunks count",
			model: &models.File{Filename: "Error file", Size: 1024, Name: "Error file"},
			err:   customerror.NewCustomError("failed to create new file: NOT NULL constraint failed: file.chunks_count", 400),
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

			model := &dto.File{
				Filename:    tt.model.Filename,
				Size:        tt.model.Size,
				Name:        tt.model.Name,
				ChunksCount: tt.model.ChunksCount,
			}

			id, err := st.Storage.CreateFile(ctx, tx, model, tt.name)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, id)

			if err == nil {
				err = tx.Commit()
				require.NoError(t, err)
			}
		})
	}
}

func TestGetFileByID(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		expected *models.File
		err      error
	}{
		{
			name:     "get valid music file",
			id:       1,
			expected: fakeFiles[0],
		},
		{
			name:     "get valid video file",
			id:       2,
			expected: fakeFiles[1],
		},
		{
			name:     "get valid text file",
			id:       3,
			expected: fakeFiles[2],
		},
		{
			name:     "get valid ZIP file",
			id:       4,
			expected: fakeFiles[3],
		},
		{
			name: "get file with not exists id",
			id:   100,
			err:  customerror.NewCustomError("file not found", 404),
		},
		{
			name: "get file with negative id",
			id:   -100,
			err:  customerror.NewCustomError("file not found", 404),
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

			model, err := st.Storage.GetFileByID(ctx, tx, tt.id, tt.name)
			require.Equal(t, tt.expected, model)
			require.Equal(t, tt.err, err)
		})
	}
}

func TestDeleteFile(t *testing.T) {
	tests := []struct {
		name     string
		id       int64
		expected int64
		err      error
	}{
		{
			name:     "delete valid text file",
			id:       3,
			expected: 3,
		},
		{
			name: "delete not exists file",
			id:   100,
			err:  customerror.NewCustomError("file not found", 404),
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

			id, err := st.Storage.DeleteFile(ctx, tx, tt.id, tt.name)
			require.Equal(t, tt.err, err)

			if err == nil {
				require.Equal(t, tt.expected, id)
				err = tx.Commit()
				require.NoError(t, err)
			}
		})
	}
}

func TestSearchFiles(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []*models.File
		err      error
	}{
		{
			name:  "get files by valid query: fake",
			query: "fake",
			expected: []*models.File{
				fakeFiles[0],
				fakeFiles[1],
				fakeFiles[3],
			},
		},
		{
			name:  "get files by valid query: ideo",
			query: "ideo",
			expected: []*models.File{
				fakeFiles[1],
			},
		},
		{
			name:  "get files by valid query: .zip",
			query: ".zip",
			expected: []*models.File{
				fakeFiles[3],
			},
		},
		{
			name:  "get files by valid query: Music",
			query: "Music",
			expected: []*models.File{
				fakeFiles[0],
			},
		},
		{
			name:  "get files by not valid query: Film",
			query: "Film",
			err:   customerror.NewCustomError("files not found", 404),
		},
		{
			name:  "get files by not valid query: text",
			query: "text",
			err:   customerror.NewCustomError("files not found", 404),
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

			res, err := st.Storage.SearchFiles(ctx, tx, tt.query, tt.name)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.expected, res)

		})
	}
}
