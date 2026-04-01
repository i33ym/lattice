package lattice

import "context"

// VectorRecord holds a vector with its associated ID and metadata.
type VectorRecord struct {
	ID       string
	Vector   []float32
	Metadata map[string]any
}

// VectorResult is a search result with a similarity score.
type VectorResult struct {
	ID       string
	Score    float32
	Metadata map[string]any
}

// VectorStore is the interface for vector similarity storage.
type VectorStore interface {
	CreateCollection(ctx context.Context, name string, dim int) error
	DropCollection(ctx context.Context, name string) error
	Upsert(ctx context.Context, collection string, records []VectorRecord) error
	Search(ctx context.Context, collection string, vector []float32, k int) ([]VectorResult, error)
	Delete(ctx context.Context, collection string, ids []string) error
}
