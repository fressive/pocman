package model

type FileType int

const (
	Document FileType = iota
	Resource
)

type File struct {
	ID        uint32   `json:"id"`
	Type      FileType `json:"type"`
	Checksum  string   `json:"checksum"`
	Extension string   `json:"extension"`
}
