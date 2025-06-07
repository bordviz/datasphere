package models

type Chunk struct {
	ID          int64  `json:"id"`
	FileID      string `json:"file_id"`
	ChunkNumber int64  `json:"chunk_number"`
	File        int64  `json:"file"`
}

type FileChunks struct {
	Count  int64    `json:"count"`
	Chunks []*Chunk `json:"chunks"`
}
