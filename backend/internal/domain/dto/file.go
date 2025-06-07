package dto

type File struct {
	Filename    string `json:"filename"`
	Size        int64  `json:"size"`
	Name        string `json:"name"`
	ChunksCount int    `json:"chunks_count"`
}
