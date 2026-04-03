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

func (r *memoryReader) Size() int64 { return r.Reader.Size() }

func (r *memoryReader) ETag() string { return r.meta.ETag }

func (r *memoryReader) Close() error { return nil }

type fileReader struct {
	io.ReadCloser
	size int64
	meta *BlobMetadata
}

func newFileReader(rc io.ReadCloser, size int64, meta *BlobMetadata) *fileReader {
	return &fileReader{ReadCloser: rc, size: size, meta: meta}
}

func (r *fileReader) Size() int64 { return r.size }

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

func (r *sectionFileReader) Size() int64 { return r.SectionReader.Size() }

func (r *sectionFileReader) ETag() string { return r.meta.ETag }

func (r *sectionFileReader) Close() error {
	return r.file.Close()
}

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
