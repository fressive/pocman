package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fressive/pocman/common/pkg/model"
)

type uploadResponse struct {
	ID         uint32         `json:"id"`
	StoredName string         `json:"stored_name"`
	Original   string         `json:"original_name"`
	Sha256     string         `json:"sha256"`
	Type       model.FileType `json:"type"`
	Extension  string         `json:"extension"`
	VulnID     uint           `json:"vuln_id"`
}

// UploadFile sends a multipart request to attach a local file to a vulnerability.
func (c *Client) UploadFile(ctx context.Context, path string, vulnID uint64, fileType model.FileType) (uploadResponse, error) {
	var out uploadResponse
	file, err := os.Open(path)
	if err != nil {
		return out, fmt.Errorf("open file %s: %w", path, err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if err := writer.WriteField("vuln_id", strconv.FormatUint(vulnID, 10)); err != nil {
		return out, err
	}
	if err := writer.WriteField("file_type", strconv.Itoa(int(fileType))); err != nil {
		return out, err
	}
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return out, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return out, fmt.Errorf("copy file contents: %w", err)
	}
	if err := writer.Close(); err != nil {
		return out, err
	}

	err = c.DoMultipart(ctx, http.MethodPost, "/api/v1/file/upload", writer.FormDataContentType(), body, &out)
	if err != nil {
		return out, err
	}

	return out, nil
}
