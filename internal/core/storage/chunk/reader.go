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
	// ErrInvalidChunkFile is returned when a chunk file does not have a valid format or header.
	ErrInvalidChunkFile = errors.New("invalid chunk file")
	// ErrCorruptChunk is returned when chunk data cannot be decoded or decompressed.
	ErrCorruptChunk = errors.New("corrupt chunk data")
)

// ReaderOptions contains configuration options for creating a Reader.
type ReaderOptions struct {
	ChunksDir string
}

// DefaultReaderOptions returns a ReaderOptions with sensible defaults.
func DefaultReaderOptions(chunksDir string) *ReaderOptions {
	if chunksDir == "" {
		chunksDir = "."
	}
	return &ReaderOptions{
		ChunksDir: chunksDir,
	}
}

// Reader provides functionality for reading log entries from chunk files.
type Reader struct {
	opts *ReaderOptions
}

// NewReader creates a new Reader with specified options.
func NewReader(opts *ReaderOptions) *Reader {
	if opts == nil {
		opts = DefaultReaderOptions("")
	}
	if opts.ChunksDir == "" {
		opts.ChunksDir = "."
	}

	return &Reader{
		opts: opts,
	}
}

// ReadHeader reads only the header from a chunk file.
func ReadHeader(chunkPath string) (*Header, error) {
	data, err := os.ReadFile(chunkPath)
	if err != nil {
		return nil, err
	}

	if isZstdCompressed(data) {
		decompressed, err := decompressZstd(data)
		if err != nil {
			return nil, fmt.Errorf("%w: zstd decompression failed: %v", ErrCorruptChunk, err)
		}
		data = decompressed
	}

	return readHeader(data)
}

// Read reads all entries from a chunk file.
func Read(chunkPath string) ([]Entry, error) {
	r := NewReader(nil)
	return r.Read(chunkPath)
}

// Read reads all entries from a chunk file using this reader.
func (r *Reader) Read(chunkPath string) ([]Entry, error) {
	fullPath := chunkPath
	if r.opts.ChunksDir != "" && !filepath.IsAbs(chunkPath) {
		fullPath = filepath.Join(r.opts.ChunksDir, chunkPath)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	entries, err := Decode(data)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorruptChunk, err)
	}

	return entries, nil
}

// Decode decompresses and parses chunk data.
func Decode(data []byte) ([]Entry, error) {
	fileFullyCompressed := isZstdCompressed(data)
	if fileFullyCompressed {
		decompressed, err := decompressZstd(data)
		if err != nil {
			return nil, fmt.Errorf("%w: zstd decompression failed: %v", ErrCorruptChunk, err)
		}
		data = decompressed
	}

	header, err := readHeader(data)
	if err != nil {
		return nil, err
	}

	var decompressed []byte

	if fileFullyCompressed {
		decompressed = data[HeaderSize():]
	} else {
		switch header.Encoding {
		case EncodingGzip:
			decompressed, err = decompressGzip(data[HeaderSize():])
		case EncodingZstd:
			decompressed, err = decompressZstd(data[HeaderSize():])
		case EncodingNone:
			decompressed = data[HeaderSize():]
		default:
			return nil, fmt.Errorf("%w: unknown encoding %d", ErrInvalidChunkFile, header.Encoding)
		}

		if err != nil {
			return nil, fmt.Errorf("%w: decompression failed: %v", ErrCorruptChunk, err)
		}
	}

	entries, err := parseEntries(bytes.NewReader(decompressed), header.EntryCount)
	if err != nil {
		return nil, fmt.Errorf("%w: parse failed: %v", ErrCorruptChunk, err)
	}

	return entries, nil
}

func isZstdCompressed(data []byte) bool {
	if len(data) < 4 {
		return false
	}
	return data[0] == 0x28 && data[1] == 0xB5 && data[2] == 0x2F && data[3] == 0xFD
}

func readHeader(data []byte) (*Header, error) {
	if len(data) < HeaderSize() {
		return nil, fmt.Errorf("%w: too short for header", ErrInvalidChunkFile)
	}

	header := &Header{}
	buf := bytes.NewReader(data[:HeaderSize()])
	if err := binary.Read(buf, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("%w: failed to read header: %v", ErrInvalidChunkFile, err)
	}

	if !ValidateHeader(header) {
		return nil, fmt.Errorf("%w: invalid header", ErrInvalidChunkFile)
	}

	return header, nil
}

func readHeaderFromReader(r io.Reader) (*Header, error) {
	header := &Header{}
	if err := binary.Read(r, binary.BigEndian, header); err != nil {
		return nil, fmt.Errorf("%w: failed to read header: %v", ErrInvalidChunkFile, err)
	}

	if !ValidateHeader(header) {
		return nil, fmt.Errorf("%w: invalid header", ErrInvalidChunkFile)
	}

	return header, nil
}

func parseEntries(r io.Reader, entryCount uint32) ([]Entry, error) {
	entries := make([]Entry, 0, entryCount)

	for i := uint32(0); i < entryCount; i++ {
		var ts int64
		if err := binary.Read(r, binary.BigEndian, &ts); err != nil {
			return nil, fmt.Errorf("failed to read timestamp: %w", err)
		}

		var msgLen uint32
		if err := binary.Read(r, binary.BigEndian, &msgLen); err != nil {
			return nil, fmt.Errorf("failed to read message length: %w", err)
		}

		msgBytes := make([]byte, msgLen)
		if _, err := io.ReadFull(r, msgBytes); err != nil {
			return nil, fmt.Errorf("failed to read message: %w", err)
		}

		entries = append(entries, SimpleEntry{
			Ts:  ts,
			Msg: msgBytes,
		})
	}

	return entries, nil
}

func decompressGzip(data []byte) ([]byte, error) {
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer gzReader.Close()

	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}

func decompressZstd(data []byte) ([]byte, error) {
	decoder, err := zstd.NewReader(nil)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	decompressed, err := decoder.DecodeAll(data, nil)
	if err != nil {
		return nil, err
	}

	return decompressed, nil
}
