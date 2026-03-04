package dto

import (
	"gorm.io/gorm"
)

// UploadedFile holds metadata persisted for every file uploaded via HTTP.
type UploadedFile struct {
	gorm.Model
	OriginalName string `json:"original_name" gorm:"size:255;not null"`
	StoredName   string `json:"stored_name" gorm:"size:255;not null"`
	RelativePath string `json:"relative_path" gorm:"size:1024;not null"`
	Size         int64  `json:"size"`
	SHA256       string `json:"sha256" gorm:"size:64;not null"`
	VulnID       uint   `json:"vuln_id" gorm:"not null;index"`
}
