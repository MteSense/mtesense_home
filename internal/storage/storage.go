package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
)

type Service struct {
	dir string
}

type UploadedFile struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
}

func NewService(dir string) *Service {
	return &Service{dir: dir}
}

func (s *Service) Ensure() error {
	return os.MkdirAll(s.dir, 0o755)
}

func (s *Service) Save(file multipart.File, header *multipart.FileHeader) (UploadedFile, error) {
	if header.Size > 5*1024*1024 {
		return UploadedFile{}, fmt.Errorf("file exceeds 5MB")
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true, ".svg": true}
	if !allowed[ext] {
		return UploadedFile{}, fmt.Errorf("unsupported file type")
	}
	name, err := randomName(ext)
	if err != nil {
		return UploadedFile{}, err
	}
	path := filepath.Join(s.dir, name)
	output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return UploadedFile{}, fmt.Errorf("create upload: %w", err)
	}
	defer output.Close()
	if _, err := io.Copy(output, file); err != nil {
		return UploadedFile{}, fmt.Errorf("write upload: %w", err)
	}
	return UploadedFile{URL: "/uploads/" + name, Filename: name}, nil
}

func randomName(ext string) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate filename: %w", err)
	}
	return hex.EncodeToString(bytes) + ext, nil
}
