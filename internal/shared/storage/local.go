package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
	serveURL string
}

func NewLocalStorage(basePath, serveURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory %s: %w", basePath, err)
	}
	return &LocalStorage{
		basePath: basePath,
		serveURL: serveURL,
	}, nil
}

func (s *LocalStorage) Upload(_ context.Context, path string, file io.Reader) (string, error) {
	fullPath := filepath.Join(s.basePath, path)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", fullPath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("failed to write file %s: %w", fullPath, err)
	}

	url := s.serveURL + "/" + path
	return url, nil
}

func (s *LocalStorage) Delete(_ context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file %s: %w", fullPath, err)
	}
	return nil
}
