package dto

type Chunk struct {
	FileID      string `json:"file_id"`
	ChunkNumber int64  `json:"chunk_number"`
	File        int64  `json:"file"`
}
