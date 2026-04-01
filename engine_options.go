package lattice

import "log/slog"

// Option configures an Engine.
type Option func(*Engine)

// WithStore sets the relational store for the engine.
func WithStore(s Store) Option {
	return func(e *Engine) { e.store = s }
}

// WithBlobStore sets the blob store for the engine.
func WithBlobStore(bs BlobStore) Option {
	return func(e *Engine) { e.blobStore = bs }
}

// WithVectorStore sets the vector store for the engine.
func WithVectorStore(vs VectorStore) Option {
	return func(e *Engine) { e.vectorStore = vs }
}

// WithLogger sets the logger for the engine.
func WithLogger(l *slog.Logger) Option {
	return func(e *Engine) { e.logger = l }
}
