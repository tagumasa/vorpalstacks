package chunk

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/klauspost/compress/zstd"
)

// fileSequence uniquifies chunk file names across concurrent writers. Its
// value only ever grows: consumers whose chunk files persist across
// restarts raise the floor through SeedFileSequence so a replayed write
// cannot reuse a name an earlier process issued.
var fileSequence uint64

// SeedFileSequence raises the file-name sequence floor to at least the
// given value. File names are unique only within one process lifetime by
// default; a consumer restarting over durable chunk files seeds the
// sequence from the highest sequence already present on disk, so a
// backfill replayed at the same timestamps can never truncate an earlier
// file by recreating its name.
func SeedFileSequence(floor uint64) {
	for {
		cur := atomic.LoadUint64(&fileSequence)
		if floor <= cur || atomic.CompareAndSwapUint64(&fileSequence, cur, floor) {
			return
		}
	}
}

var (
	// ErrEmptyEntries is returned when attempting to write a chunk with no entries.
	ErrEmptyEntries = errors.New("cannot write empty entries")
	// ErrChunkTooLarge is returned when a chunk exceeds the maximum allowed size.
	ErrChunkTooLarge = errors.New("chunk exceeds maximum size")
	// ErrMessageTooLarge is returned when a single message exceeds the maximum allowed size.
	ErrMessageTooLarge = errors.New("message exceeds maximum size")
	// ErrWriteFailed is returned when writing a chunk to disk fails.
	ErrWriteFailed = errors.New("failed to write chunk")
)

// Entry represents a timestamped log entry.
type Entry interface {
	Timestamp() int64
	Message() []byte
}

// Ingestible is an optional interface implemented by entries that carry
// an ingestion timestamp.  The writer checks for this interface and, when
// present, serialises the ingestion timestamp into V2 chunk files.
type Ingestible interface {
	IngestionTimeUnixMilli() int64
}

// SimpleEntry implements Entry for basic log entries.
type SimpleEntry struct {
	Ts          int64
	IngestionTs int64
	Msg         []byte
}

// Timestamp returns the Unix timestamp in nanoseconds for this entry.
func (e SimpleEntry) Timestamp() int64 {
	return e.Ts
}

// Message returns the raw byte content of this entry.
func (e SimpleEntry) Message() []byte {
	return e.Msg
}

// IngestionTimeUnixMilli returns the ingestion timestamp in Unix milliseconds.
func (e SimpleEntry) IngestionTimeUnixMilli() int64 {
	return e.IngestionTs
}

// WriterOptions contains configuration options for creating a Writer.
type WriterOptions struct {
	ChunksDir   string
	ChunkSize   int
	Encoding    Encoding
	SyncOnWrite bool
}

// DefaultWriterOptions returns a WriterOptions with sensible defaults.
func DefaultWriterOptions(chunksDir string) *WriterOptions {
	return &WriterOptions{
		ChunksDir:   chunksDir,
		ChunkSize:   DefaultChunkSize,
		Encoding:    EncodingZstd,
		SyncOnWrite: false,
	}
}

// Writer provides functionality for writing log entries to chunk files.
type Writer struct {
	opts      *WriterOptions
	buffer    []Entry
	chunkPath string
}

// NewWriter creates a new Writer with specified options.
func NewWriter(opts *WriterOptions) *Writer {
	if opts == nil {
		opts = &WriterOptions{}
	}
	if opts.ChunksDir == "" {
		opts.ChunksDir = "."
	}
	if opts.ChunkSize <= 0 {
		opts.ChunkSize = DefaultChunkSize
	}
	if opts.Encoding == 0 {
		opts.Encoding = EncodingZstd
	}

	return &Writer{
		opts:      opts,
		buffer:    make([]Entry, 0, opts.ChunkSize),
		chunkPath: "",
	}
}

// Write writes a single entry to buffer, flushing if buffer is full.
func (w *Writer) Write(entry Entry) error {
	if len(entry.Message()) > MaxMessageSize {
		return fmt.Errorf("%w: message length %d", ErrMessageTooLarge, len(entry.Message()))
	}

	w.buffer = append(w.buffer, entry)

	if len(w.buffer) >= w.opts.ChunkSize {
		if _, err := w.Flush(); err != nil {
			return err
		}
	}

	return nil
}

// WriteBatch writes multiple entries to buffer, flushing if buffer is full.
func (w *Writer) WriteBatch(entries []Entry) error {
	for _, e := range entries {
		if len(e.Message()) > MaxMessageSize {
			return fmt.Errorf("%w: message length %d", ErrMessageTooLarge, len(e.Message()))
		}
	}

	w.buffer = append(w.buffer, entries...)

	if len(w.buffer) >= w.opts.ChunkSize {
		if _, err := w.Flush(); err != nil {
			return err
		}
	}

	return nil
}

// Flush forces any buffered entries to be written to disk.
// Returns the path to the written chunk file if any entries were flushed.
func (w *Writer) Flush() (string, error) {
	if len(w.buffer) == 0 {
		return "", nil
	}

	entries := w.buffer
	w.buffer = w.buffer[:0]

	// The file name must be unique: concurrent writers flushing batches
	// whose first entries share a millisecond timestamp and whose entry
	// counts match would otherwise write to the same path, and every chunk
	// index pointing at that path would read the one surviving file's
	// events. The per-process sequence guarantees uniqueness.
	chunkID := fmt.Sprintf("%d-%d-%d", entries[0].Timestamp(), len(entries), atomic.AddUint64(&fileSequence, 1))
	chunkPath := filepath.Join(w.opts.ChunksDir, chunkID+".vlog")

	absPath, err := filepath.Abs(chunkPath)
	if err != nil {
		absPath = chunkPath
	}
	chunkPath = absPath

	_, err = writeEncodedChunk(w, chunkPath, entries)

	if err != nil {
		// The partial file must not survive the failure: nothing will
		// index it, and no sweep knows about unindexed files, so keeping
		// it leaks storage for ever. The removal error is secondary to
		// the write failure the caller must see.
		os.Remove(chunkPath)
		w.buffer = append(entries, w.buffer...)
		return "", fmt.Errorf("%w: %w", ErrWriteFailed, err)
	}

	w.chunkPath = chunkPath
	return chunkPath, nil
}

// GetChunkPath returns the path of the most recently written chunk.
func (w *Writer) GetChunkPath() string {
	return w.chunkPath
}

// syncDir fsyncs a directory so the creation of a file inside it survives
// a crash: on filesystems where fsyncing the file alone does not persist
// the directory entry, an index committed after the file fsync could point
// at a file the crash never made visible.
func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

// encodeChunk dispatches one buffered batch to its encoding writer. It is
// indirected through the writeEncodedChunk var so tests can exercise the
// write-failure path of Flush (a partial file on disk must not survive
// the error).
func encodeChunk(w *Writer, path string, entries []Entry) (*Header, error) {
	switch w.opts.Encoding {
	case EncodingGzip:
		return w.writeGzipChunk(path, entries)
	case EncodingZstd:
		return w.writeZstdChunk(path, entries)
	default:
		return w.writeRawChunk(path, entries)
	}
}

var writeEncodedChunk = encodeChunk

func (w *Writer) writeGzipChunk(path string, entries []Entry) (*Header, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gzWriter := gzip.NewWriter(f)

	header, err := w.writeEntries(gzWriter, entries)
	if err != nil {
		gzWriter.Close()
		return nil, err
	}

	if err := gzWriter.Close(); err != nil {
		return nil, fmt.Errorf("gzip close failed: %w", err)
	}

	if w.opts.SyncOnWrite {
		if err := f.Sync(); err != nil {
			return nil, err
		}
		if err := syncDir(filepath.Dir(path)); err != nil {
			return nil, err
		}
	}

	return header, nil
}

func (w *Writer) writeZstdChunk(path string, entries []Entry) (*Header, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	header, err := w.writeEntries(&buf, entries)
	if err != nil {
		return nil, err
	}

	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(zstd.EncoderLevelFromZstd(3)))
	if err != nil {
		return nil, err
	}
	defer encoder.Close()
	compressed := encoder.EncodeAll(buf.Bytes(), nil)

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := f.Write(compressed); err != nil {
		return nil, err
	}

	if w.opts.SyncOnWrite {
		if err := f.Sync(); err != nil {
			return nil, err
		}
		if err := syncDir(filepath.Dir(path)); err != nil {
			return nil, err
		}
	}

	return header, nil
}

func (w *Writer) writeRawChunk(path string, entries []Entry) (*Header, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	header, err := w.writeEntries(f, entries)
	if err != nil {
		return nil, err
	}

	if w.opts.SyncOnWrite {
		if err := f.Sync(); err != nil {
			return nil, err
		}
		if err := syncDir(filepath.Dir(path)); err != nil {
			return nil, err
		}
	}

	return header, nil
}

func (w *Writer) writeEntries(writer io.Writer, entries []Entry) (*Header, error) {
	if len(entries) == 0 {
		return nil, ErrEmptyEntries
	}

	magic := [4]byte{'V', 'L', 'O', 'G'}
	header := &Header{
		Magic:      magic,
		Version:    VersionV2,
		Encoding:   w.opts.Encoding.Uint8(),
		EntryCount: uint32(len(entries)),
		MinTs:      entries[0].Timestamp(),
		MaxTs:      entries[0].Timestamp(),
	}

	for _, e := range entries {
		ts := e.Timestamp()
		if ts < header.MinTs {
			header.MinTs = ts
		}
		if ts > header.MaxTs {
			header.MaxTs = ts
		}
	}

	if err := binary.Write(writer, binary.BigEndian, header); err != nil {
		return nil, err
	}

	for _, e := range entries {
		if err := binary.Write(writer, binary.BigEndian, e.Timestamp()); err != nil {
			return nil, err
		}

		var ingestionTs int64
		if ig, ok := e.(Ingestible); ok {
			ingestionTs = ig.IngestionTimeUnixMilli()
		}
		if err := binary.Write(writer, binary.BigEndian, ingestionTs); err != nil {
			return nil, err
		}

		msgBytes := e.Message()
		if err := binary.Write(writer, binary.BigEndian, uint32(len(msgBytes))); err != nil {
			return nil, err
		}
		if _, err := writer.Write(msgBytes); err != nil {
			return nil, err
		}
	}

	return header, nil
}
