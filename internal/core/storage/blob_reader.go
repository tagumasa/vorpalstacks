package storage

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type memoryReader struct {
	*bytes.Reader
	meta *BlobMetadata
}

func newMemoryReader(data []byte, meta *BlobMetadata) *memoryReader {
	return &memoryReader{
		Reader: bytes.NewReader(data),
		meta:   meta,
	}
}

// Size returns the total size of the in-memory blob in bytes.
func (r *memoryReader) Size() int64 { return r.Reader.Size() }

// ETag returns the entity tag for the in-memory blob.
func (r *memoryReader) ETag() string { return r.meta.ETag }

// Close releases resources held by the memory reader (no-op for in-memory blobs).
func (r *memoryReader) Close() error { return nil }

type fileReader struct {
	io.ReadCloser
	size int64
	meta *BlobMetadata
}

func newFileReader(rc io.ReadCloser, size int64, meta *BlobMetadata) *fileReader {
	return &fileReader{ReadCloser: rc, size: size, meta: meta}
}

// Size returns the total size of the file-backed blob in bytes.
func (r *fileReader) Size() int64 { return r.size }

// ETag returns the entity tag for the file-backed blob.
func (r *fileReader) ETag() string { return r.meta.ETag }

type sectionFileReader struct {
	*io.SectionReader
	file *os.File
	meta *BlobMetadata
}

func newSectionFileReader(f *os.File, offset, length int64, meta *BlobMetadata) *sectionFileReader {
	return &sectionFileReader{
		SectionReader: io.NewSectionReader(f, offset, length),
		file:          f,
		meta:          meta,
	}
}

// Size returns the size of the section being read.
func (r *sectionFileReader) Size() int64 { return r.SectionReader.Size() }

// ETag returns the entity tag for the section's parent blob.
func (r *sectionFileReader) ETag() string { return r.meta.ETag }

// Close releases the underlying file handle for the section reader.
func (r *sectionFileReader) Close() error {
	return r.file.Close()
}

// openAndStat opens a file and retrieves its FileInfo, returning a descriptive error if the file does not exist.
func openAndStat(path, notFoundMsg string) (*os.File, os.FileInfo, error) {
	// #nosec G304
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, fmt.Errorf("object not found: %s", notFoundMsg)
		}
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, fmt.Errorf("failed to stat file: %w", err)
	}
	return f, info, nil
}

// clampRange constrains a byte range [offset, offset+length) to lie within [0, maxSize].
func clampRange(offset, length, maxSize int64) (start, end int64) {
	start = offset
	if start < 0 {
		start = 0
	}
	if start > maxSize {
		start = maxSize
	}
	end = start + length
	if end > maxSize {
		end = maxSize
	}
	if end < start {
		end = start
	}
	return start, end
}
