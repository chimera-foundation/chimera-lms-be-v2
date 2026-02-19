package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	Upload(ctx context.Context, path string, file io.Reader) (url string, err error)
	Delete(ctx context.Context, path string) error
}
