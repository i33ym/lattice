package lattice

import (
	"context"
	"io"
)

// BlobStore is the interface for binary object storage.
type BlobStore interface {
	Put(ctx context.Context, key string, r io.Reader, meta BlobMeta) error
	Get(ctx context.Context, key string) (io.ReadCloser, BlobMeta, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}
