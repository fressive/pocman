package http

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fressive/pocman/server/internal/conf"
	"github.com/fressive/pocman/server/internal/data"
	"github.com/fressive/pocman/server/internal/model/dto"
	"github.com/fressive/pocman/server/internal/server/http/response"
	"github.com/gin-gonic/gin"
)

const (
	uploadFieldName = "file"
	maxUploadBytes  = 128 << 20 // 128 MiB per upload keeps memory bounded
)

type FileHandler struct{}

// NewFileHandler instantiates a FileHandler.
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// FileUpload receives a multipart file, sanitizes the filename, persists it inside storage.path,
// and returns the stored path so callers can reference the artifact.
func (h *FileHandler) FileUpload(c *gin.Context) {
	storageRoot, err := prepareStorageRoot()
	if err != nil {
		response.Error(c, 11001, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadBytes)
	upload, err := c.FormFile(uploadFieldName)
	if err != nil {
		response.Error(c, 11002, fmt.Sprintf("failed to read upload: %v", err))
		return
	}

	originalName := sanitizeFilename(upload.Filename)
	if originalName == "" {
		originalName = "file"
	}

	vulnIDParam := c.PostForm("vuln_id")
	if vulnIDParam == "" {
		response.Error(c, 11012, "vuln_id is required")
		return
	}
	vulnID, err := strconv.ParseUint(vulnIDParam, 10, 64)
	if err != nil || vulnID == 0 {
		response.Error(c, 11013, fmt.Sprintf("invalid vuln_id: %v", err))
		return
	}

	targetDir := filepath.Join(storageRoot, fmt.Sprintf("%d", vulnID))
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		response.Error(c, 11003, fmt.Sprintf("failed to prepare storage directory: %v", err))
		return
	}

	now := time.Now()
	initialName := fmt.Sprintf("%s_%s", now.Format("20060102_150405"), originalName)
	target := filepath.Join(targetDir, initialName)
	if err := c.SaveUploadedFile(upload, target); err != nil {
		response.Error(c, 11004, fmt.Sprintf("failed to persist file: %v", err))
		return
	}

	sha256sum, err := computeSHA256(target)
	if err != nil {
		response.Error(c, 11005, fmt.Sprintf("failed to compute file hash: %v", err))
		return
	}

	// We will NOT detect replicates for some reference concerns when removing files

	hashExt := filepath.Ext(originalName)
	hashedName := sha256sum
	if hashExt != "" {
		hashedName = fmt.Sprintf("%s%s", sha256sum, hashExt)
	}
	hashTarget := filepath.Join(targetDir, hashedName)
	if hashTarget != target {
		if err := os.Rename(target, hashTarget); err != nil {
			response.Error(c, 11007, fmt.Sprintf("failed to rename file: %v", err))
			return
		}
		target = hashTarget
	}

	relPath, err := filepath.Rel(storageRoot, target)
	if err != nil {
		relPath = hashedName
	}
	relativePath := filepath.ToSlash(relPath)

	record := dto.UploadedFile{
		OriginalName: originalName,
		StoredName:   hashedName,
		RelativePath: relativePath,
		Size:         upload.Size,
		SHA256:       sha256sum,
		VulnID:       uint(vulnID),
	}
	if err := data.DB.Create(&record).Error; err != nil {
		response.Error(c, 11006, fmt.Sprintf("failed to persist metadata: %v", err))
		return
	}

	response.Success(c, gin.H{
		"id":            record.ID,
		"created_at":    record.CreatedAt,
		"original_name": originalName,
		"stored_name":   hashedName,
		"relative_path": relativePath,
		"size":          upload.Size,
		"sha256":        sha256sum,
	})
}

// FileDownload streams a previously uploaded file using the record `id` or `sha256`.
func (h *FileHandler) FileDownload(c *gin.Context) {
	storageRoot, err := prepareStorageRoot()
	if err != nil {
		response.Error(c, 11001, err.Error())
		return
	}

	vulnParam := c.Query("vuln_id")
	if vulnParam != "" {
		vulnID, err := strconv.ParseUint(vulnParam, 10, 64)
		if err != nil || vulnID == 0 {
			response.Error(c, 11014, fmt.Sprintf("invalid vuln_id: %v", err))
			return
		}

		var records []dto.UploadedFile
		if err := data.DB.Where("vuln_id = ?", vulnID).Find(&records).Error; err != nil {
			response.Error(c, 11010, fmt.Sprintf("file metadata not found: %v", err))
			return
		}
		if len(records) == 0 {
			response.Error(c, 11010, "file metadata not found")
			return
		}
		if len(records) == 1 {
			serveFile(c, storageRoot, records[0])
			return
		}
		streamVulnFiles(c, storageRoot, uint(vulnID), records)
		return
	}

	idParam := c.Query("id")
	shaParam := c.Query("sha256")
	if idParam == "" && shaParam == "" {
		response.Error(c, 11008, "must provide id or sha256 to download")
		return
	}

	var record dto.UploadedFile
	query := data.DB
	if idParam != "" {
		id, err := strconv.ParseUint(idParam, 10, 64)
		if err != nil {
			response.Error(c, 11009, fmt.Sprintf("invalid id: %v", err))
			return
		}
		query = query.Where("id = ?", id)
	} else {
		query = query.Where("sha256 = ?", shaParam)
	}

	if err := query.First(&record).Error; err != nil {
		response.Error(c, 11010, fmt.Sprintf("file metadata not found: %v", err))
		return
	}

	serveFile(c, storageRoot, record)
}

func prepareStorageRoot() (string, error) {
	if conf.ServerConfig.Data == nil || conf.ServerConfig.Data.Storage == nil {
		return "", fmt.Errorf("storage configuration is missing")
	}

	storagePath := strings.TrimSpace(conf.ServerConfig.Data.Storage.Path)
	if storagePath == "" {
		return "", fmt.Errorf("storage.path is empty")
	}

	abs, err := filepath.Abs(storagePath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve storage path: %w", err)
	}

	if err := os.MkdirAll(abs, 0o755); err != nil {
		return "", fmt.Errorf("failed to ensure storage directory: %w", err)
	}

	return abs, nil
}

func sanitizeFilename(name string) string {
	clean := filepath.Base(name)
	if clean == "" {
		return ""
	}

	clean = strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '.', r == '-', r == '_':
			return r
		default:
			return '_'
		}
	}, clean)

	clean = strings.Trim(clean, "._-")
	if clean == "" {
		return ""
	}

	return clean
}

func computeSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file for hashing: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash file contents: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func serveFile(c *gin.Context, storageRoot string, record dto.UploadedFile) {
	filePath := filepath.Join(storageRoot, record.RelativePath)
	if _, err := os.Stat(filePath); err != nil {
		response.Error(c, 11011, fmt.Sprintf("file missing: %v", err))
		return
	}
	c.FileAttachment(filePath, record.OriginalName)
}

func streamVulnFiles(c *gin.Context, storageRoot string, vulnID uint, records []dto.UploadedFile) {
	paths := make([]string, len(records))
	for i, record := range records {
		filePath := filepath.Join(storageRoot, record.RelativePath)
		if _, err := os.Stat(filePath); err != nil {
			response.Error(c, 11011, fmt.Sprintf("file missing: %v", err))
			return
		}
		paths[i] = filePath
	}

	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"vuln_%d_files.zip\"", vulnID))
	zw := zip.NewWriter(c.Writer)
	defer zw.Close()

	for i, record := range records {
		if err := addFileToZip(zw, paths[i], record); err != nil {
			response.Error(c, 11015, fmt.Sprintf("failed to add file to archive: %v", err))
			return
		}
	}
}

func addFileToZip(zw *zip.Writer, filePath string, record dto.UploadedFile) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	entryName := record.StoredName
	if entryName == "" {
		entryName = record.OriginalName
	}

	w, err := zw.Create(entryName)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, f)
	return err
}
