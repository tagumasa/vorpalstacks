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

	"github.com/klauspost/compress/zstd"
)

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

// SimpleEntry implements Entry for basic log entries.
type SimpleEntry struct {
	Ts  int64
	Msg []byte
}

// Timestamp returns the Unix timestamp in nanoseconds for this entry.
func (e SimpleEntry) Timestamp() int64 {
	return e.Ts
}

// Message returns the raw byte content of this entry.
func (e SimpleEntry) Message() []byte {
	return e.Msg
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

	chunkID := fmt.Sprintf("%d-%d", entries[0].Timestamp(), len(entries))
	chunkPath := filepath.Join(w.opts.ChunksDir, chunkID+".vlog")

	absPath, err := filepath.Abs(chunkPath)
	if err != nil {
		absPath = chunkPath
	}
	chunkPath = absPath

	switch w.opts.Encoding {
	case EncodingGzip:
		_, err = w.writeGzipChunk(chunkPath, entries)
	case EncodingZstd:
		_, err = w.writeZstdChunk(chunkPath, entries)
	default:
		_, err = w.writeRawChunk(chunkPath, entries)
	}

	if err != nil {
		w.buffer = append(entries, w.buffer...)
		return "", fmt.Errorf("%w: %v", ErrWriteFailed, err)
	}

	w.chunkPath = chunkPath
	return chunkPath, nil
}

// GetChunkPath returns the path of the most recently written chunk.
func (w *Writer) GetChunkPath() string {
	return w.chunkPath
}

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
		Version:    VersionV1,
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
