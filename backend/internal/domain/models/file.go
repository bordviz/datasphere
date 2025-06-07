package models

type File struct {
	ID          int64  `json:"id"`
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	Name        string `json:"name"`
	ChunksCount int    `json:"chunks_count"`
}
