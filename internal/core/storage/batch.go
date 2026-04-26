package storage

// Batch provides atomic multi-key write operations within a Bucket.
// Operations are buffered until Commit is called. Close discards
// uncommitted operations and releases resources.
type Batch interface {
	Put(key, value []byte) error
	Delete(key []byte) error
	DeleteRange(start, end []byte) error
	Commit() error
	Close() error
}

// BatchBucket extends Bucket with atomic batch writes and range deletion.
// Services that require multi-key atomicity (e.g. graph databases) use
// this interface instead of plain Bucket.
type BatchBucket interface {
	Bucket
	NewBatch() Batch
	DeleteRange(start, end []byte) error
}
