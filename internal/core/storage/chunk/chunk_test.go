package chunk

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestWriterAndReader(t *testing.T) {
	tempDir := t.TempDir()
	chunksDir := filepath.Join(tempDir, "chunks")

	entries := []Entry{
		SimpleEntry{Ts: 1000, Msg: []byte("message 1")},
		SimpleEntry{Ts: 2000, Msg: []byte("message 2")},
		SimpleEntry{Ts: 3000, Msg: []byte("message 3")},
	}

	t.Run("Zstd encoding", func(t *testing.T) {
		opts := &WriterOptions{
			ChunksDir: chunksDir,
			Encoding:  EncodingZstd,
			ChunkSize: 10,
		}

		w := NewWriter(opts)
		if err := w.WriteBatch(entries); err != nil {
			t.Fatalf("WriteBatch failed: %v", err)
		}

		chunkPath, err := w.Flush()
		if err != nil {
			t.Fatalf("Flush failed: %v", err)
		}

		if chunkPath == "" {
			t.Fatal("chunkPath is empty")
		}

		header, err := ReadHeader(chunkPath)
		if err != nil {
			t.Fatalf("ReadHeader failed: %v", err)
		}

		if header.EntryCount != 3 {
			t.Errorf("Expected 3 entries, got %d", header.EntryCount)
		}

		if header.MinTs != 1000 {
			t.Errorf("Expected MinTs 1000, got %d", header.MinTs)
		}

		if header.MaxTs != 3000 {
			t.Errorf("Expected MaxTs 3000, got %d", header.MaxTs)
		}

		if header.Encoding != EncodingZstd {
			t.Errorf("Expected EncodingZstd, got %d", header.Encoding)
		}

		readEntries, err := Read(chunkPath)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}

		if len(readEntries) != 3 {
			t.Fatalf("Expected 3 entries, got %d", len(readEntries))
		}

		for i, entry := range readEntries {
			if entry.Timestamp() != entries[i].Timestamp() {
				t.Errorf("Entry %d: expected timestamp %d, got %d", i, entries[i].Timestamp(), entry.Timestamp())
			}
			if string(entry.Message()) != string(entries[i].Message()) {
				t.Errorf("Entry %d: expected message %q, got %q", i, entries[i].Message(), entry.Message())
			}
		}
	})

	t.Run("Gzip encoding", func(t *testing.T) {
		opts := &WriterOptions{
			ChunksDir:   filepath.Join(chunksDir, "gzip"),
			Encoding:    EncodingGzip,
			ChunkSize:   10,
			SyncOnWrite: true,
		}

		w := NewWriter(opts)
		if err := w.WriteBatch(entries); err != nil {
			t.Fatalf("WriteBatch failed: %v", err)
		}

		chunkPath, err := w.Flush()
		if err != nil {
			t.Fatalf("Flush failed: %v", err)
		}

		if chunkPath == "" {
			t.Fatal("Expected chunk path to be returned")
		}

		data, err := os.ReadFile(chunkPath)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		if len(data) < 2 || data[0] != 0x1f || data[1] != 0x8b {
			t.Errorf("Expected gzip magic bytes at start of file")
		}

		// Full round-trip: ReadHeader + Read should return correct entries.
		header, err := ReadHeader(chunkPath)
		if err != nil {
			t.Fatalf("ReadHeader failed: %v", err)
		}
		if header.EntryCount != 3 {
			t.Errorf("Expected 3 entries, got %d", header.EntryCount)
		}
		if header.Encoding != EncodingGzip {
			t.Errorf("Expected EncodingGzip, got %d", header.Encoding)
		}

		readEntries, err := Read(chunkPath)
		if err != nil {
			t.Fatalf("Read failed: %v", err)
		}
		if len(readEntries) != 3 {
			t.Fatalf("Expected 3 entries, got %d", len(readEntries))
		}
		for i, entry := range readEntries {
			if entry.Timestamp() != entries[i].Timestamp() {
				t.Errorf("Entry %d: expected timestamp %d, got %d", i, entries[i].Timestamp(), entry.Timestamp())
			}
			if string(entry.Message()) != string(entries[i].Message()) {
				t.Errorf("Entry %d: expected message %q, got %q", i, entries[i].Message(), entry.Message())
			}
		}
	})

	t.Run("Auto flush on buffer full", func(t *testing.T) {
		opts := &WriterOptions{
			ChunksDir: chunksDir,
			Encoding:  EncodingZstd,
			ChunkSize: 2,
		}

		w := NewWriter(opts)

		if err := w.Write(entries[0]); err != nil {
			t.Fatalf("Write first failed: %v", err)
		}
		if err := w.Write(entries[1]); err != nil {
			t.Fatalf("Write second failed: %v", err)
		}

		chunkPath1 := w.GetChunkPath()
		if chunkPath1 == "" {
			t.Fatal("Expected chunk to be flushed after 2 entries")
		}

		if err := w.Write(entries[2]); err != nil {
			t.Fatalf("Write third failed: %v", err)
		}

		chunkPath2 := w.GetChunkPath()
		if chunkPath2 != chunkPath1 {
			t.Fatal("Expected same chunk path - third entry hasn't triggered flush yet")
		}

		finalPath, err := w.Flush()
		if err != nil {
			t.Fatalf("Final flush failed: %v", err)
		}
		if finalPath == chunkPath1 {
			t.Fatal("Expected new chunk path after final flush")
		}
	})
}

func TestValidateHeader(t *testing.T) {
	magic := [4]byte{'V', 'L', 'O', 'G'}

	t.Run("Valid header", func(t *testing.T) {
		header := &Header{
			Magic:      magic,
			Version:    VersionV1,
			Encoding:   EncodingZstd,
			EntryCount: 10,
			MinTs:      1000,
			MaxTs:      2000,
		}

		if !ValidateHeader(header) {
			t.Error("Valid header should pass validation")
		}
	})

	t.Run("Invalid magic", func(t *testing.T) {
		header := &Header{
			Magic:      [4]byte{'X', 'X', 'X', 'X'},
			Version:    VersionV1,
			Encoding:   EncodingZstd,
			EntryCount: 10,
			MinTs:      1000,
			MaxTs:      2000,
		}

		if ValidateHeader(header) {
			t.Error("Invalid magic should fail validation")
		}
	})

	t.Run("Invalid version", func(t *testing.T) {
		header := &Header{
			Magic:      magic,
			Version:    99,
			Encoding:   EncodingZstd,
			EntryCount: 10,
			MinTs:      1000,
			MaxTs:      2000,
		}

		if ValidateHeader(header) {
			t.Error("Invalid version should fail validation")
		}
	})

	t.Run("Invalid encoding", func(t *testing.T) {
		header := &Header{
			Magic:      magic,
			Version:    VersionV1,
			Encoding:   99,
			EntryCount: 10,
			MinTs:      1000,
			MaxTs:      2000,
		}

		if ValidateHeader(header) {
			t.Error("Invalid encoding should fail validation")
		}
	})
}

func TestMaxMessageSize(t *testing.T) {
	tempDir := t.TempDir()
	opts := &WriterOptions{
		ChunksDir: tempDir,
		Encoding:  EncodingZstd,
		ChunkSize: 10,
	}

	w := NewWriter(opts)

	largeMsg := make([]byte, MaxMessageSize+1)
	entry := SimpleEntry{Ts: 1000, Msg: largeMsg}

	err := w.Write(entry)
	if err == nil {
		t.Error("Expected error for message exceeding MaxMessageSize")
	}
}

func TestEmptyWrite(t *testing.T) {
	tempDir := t.TempDir()
	opts := &WriterOptions{
		ChunksDir: tempDir,
		Encoding:  EncodingZstd,
		ChunkSize: 10,
	}

	w := NewWriter(opts)

	chunkPath, err := w.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	if chunkPath != "" {
		t.Error("Expected empty chunkPath for empty buffer")
	}
}

func TestReadNonExistentFile(t *testing.T) {
	_, err := Read("/non/existent/file.vlog")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// SeedFileSequence raises the file-name sequence floor: the next flushed
// file must carry a sequence above it, so a process restarting over
// durable files never recreates an earlier file's name (os.Create would
// truncate it).
func TestSeedFileSequenceRaisesFloor(t *testing.T) {
	const floor = 1 << 20
	SeedFileSequence(floor)

	w := NewWriter(&WriterOptions{ChunksDir: t.TempDir(), Encoding: EncodingZstd, ChunkSize: 10})
	if err := w.WriteBatch([]Entry{SimpleEntry{Ts: 1000, Msg: []byte("seeded")}}); err != nil {
		t.Fatal(err)
	}
	path, err := w.Flush()
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(strings.TrimSuffix(filepath.Base(path), ".vlog"), "-")
	if len(parts) != 3 {
		t.Fatalf("chunk file name %q is not the three-component form", filepath.Base(path))
	}
	seq, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	if seq <= floor {
		t.Fatalf("file sequence %d does not exceed the seeded floor %d", seq, floor)
	}
}

// A write that fails partway through must not leak the partial file:
// nothing will index it, and no sweep knows about unindexed files.
func TestFlushRemovesPartialFileOnWriteFailure(t *testing.T) {
	dir := t.TempDir()
	orig := writeEncodedChunk
	writeEncodedChunk = func(w *Writer, path string, entries []Entry) (*Header, error) {
		// Simulate a write that created the file and failed partway
		// through: the partial file is on disk, the caller sees an error.
		if err := os.WriteFile(path, []byte("partial"), 0644); err != nil {
			return nil, err
		}
		return nil, errors.New("injected mid-write failure")
	}
	t.Cleanup(func() { writeEncodedChunk = orig })

	w := NewWriter(&WriterOptions{ChunksDir: dir, Encoding: EncodingZstd, ChunkSize: 10})
	if err := w.WriteBatch([]Entry{SimpleEntry{Ts: 1000, Msg: []byte("doomed")}}); err != nil {
		t.Fatal(err)
	}
	if _, err := w.Flush(); err == nil {
		t.Fatal("Flush reported success despite the injected write failure")
	}
	remaining, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range remaining {
		if !f.IsDir() {
			t.Fatalf("partial chunk file leaked: %s", filepath.Join(dir, f.Name()))
		}
	}
	// The buffered entries survive the failure for a retry.
	if err := w.WriteBatch([]Entry{SimpleEntry{Ts: 1000, Msg: []byte("doomed")}}); err != nil {
		t.Fatal(err)
	}
}

// zeroReader serves an endless stream of zero bytes without holding the
// expansion in memory.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = 0
	}
	return len(p), nil
}

// The entry count and every message length come from the file itself: a
// corrupt or planted header must fail as corruption, not allocate the
// requested gigabytes until the process dies.
func TestDecodeRejectsHostileAllocationRequests(t *testing.T) {
	magic := [4]byte{'V', 'L', 'O', 'G'}

	t.Run("entry count beyond the stream", func(t *testing.T) {
		var buf bytes.Buffer
		header := &Header{Magic: magic, Version: VersionV2, Encoding: EncodingNone, EntryCount: 0xFFFFFFFF}
		if err := binary.Write(&buf, binary.BigEndian, header); err != nil {
			t.Fatal(err)
		}
		if _, err := Decode(buf.Bytes()); err == nil {
			t.Fatal("Decode accepted an entry count far beyond the bytes the file holds")
		}
	})

	t.Run("message length beyond the stream", func(t *testing.T) {
		var buf bytes.Buffer
		header := &Header{Magic: magic, Version: VersionV2, Encoding: EncodingNone, EntryCount: 1}
		if err := binary.Write(&buf, binary.BigEndian, header); err != nil {
			t.Fatal(err)
		}
		if err := binary.Write(&buf, binary.BigEndian, int64(1000)); err != nil {
			t.Fatal(err)
		}
		if err := binary.Write(&buf, binary.BigEndian, int64(2000)); err != nil {
			t.Fatal(err)
		}
		if err := binary.Write(&buf, binary.BigEndian, uint32(0xFFFFFFFF)); err != nil {
			t.Fatal(err)
		}
		if _, err := Decode(buf.Bytes()); err == nil {
			t.Fatal("Decode accepted a message length beyond the bytes the file holds")
		}
	})

	t.Run("decompression bomb", func(t *testing.T) {
		var compressed bytes.Buffer
		gz := gzip.NewWriter(&compressed)
		// Just past the bound: the limited read stops at the bound and
		// reports corruption instead of materialising the expansion.
		if _, err := io.CopyN(gz, zeroReader{}, MaxDecompressedChunkBytes+1); err != nil {
			t.Fatal(err)
		}
		if err := gz.Close(); err != nil {
			t.Fatal(err)
		}
		if _, err := Decode(compressed.Bytes()); err == nil {
			t.Fatal("Decode accepted a payload expanding beyond the decompressed-size bound")
		}
	})
}

func TestDecodeInvalidData(t *testing.T) {
	invalidData := []byte("invalid chunk data")
	_, err := Decode(invalidData)
	if err == nil {
		t.Error("Expected error for invalid data")
	}
}

func TestWriteToExistingDirectory(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := filepath.Join(tempDir, "existing")
	if err := os.MkdirAll(existingDir, 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	entries := []Entry{
		SimpleEntry{Ts: 1000, Msg: []byte("test")},
	}

	opts := &WriterOptions{
		ChunksDir: existingDir,
		Encoding:  EncodingZstd,
		ChunkSize: 10,
	}

	w := NewWriter(opts)
	if err := w.WriteBatch(entries); err != nil {
		t.Fatalf("WriteBatch failed: %v", err)
	}

	if _, err := w.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}
}

func TestIngestionTimeRoundTrip(t *testing.T) {
	tempDir := t.TempDir()

	entries := []Entry{
		SimpleEntry{Ts: 1000, IngestionTs: 5000, Msg: []byte("message 1")},
		SimpleEntry{Ts: 2000, IngestionTs: 6000, Msg: []byte("message 2")},
		SimpleEntry{Ts: 3000, IngestionTs: 7000, Msg: []byte("message 3")},
	}

	for _, enc := range []Encoding{EncodingZstd, EncodingGzip} {
		t.Run(enc.String(), func(t *testing.T) {
			opts := &WriterOptions{
				ChunksDir: filepath.Join(tempDir, enc.String()),
				Encoding:  enc,
				ChunkSize: 10,
			}

			w := NewWriter(opts)
			if err := w.WriteBatch(entries); err != nil {
				t.Fatalf("WriteBatch failed: %v", err)
			}

			chunkPath, err := w.Flush()
			if err != nil {
				t.Fatalf("Flush failed: %v", err)
			}
			if chunkPath == "" {
				t.Fatal("chunkPath is empty")
			}

			header, err := ReadHeader(chunkPath)
			if err != nil {
				t.Fatalf("ReadHeader failed: %v", err)
			}

			if header.Version != VersionV2 {
				t.Errorf("Expected VersionV2 for Ingestible entries, got version %d", header.Version)
			}

			readEntries, err := Read(chunkPath)
			if err != nil {
				t.Fatalf("Read failed: %v", err)
			}

			if len(readEntries) != 3 {
				t.Fatalf("Expected 3 entries, got %d", len(readEntries))
			}

			for i, entry := range readEntries {
				if entry.Timestamp() != entries[i].Timestamp() {
					t.Errorf("Entry %d: expected timestamp %d, got %d", i, entries[i].Timestamp(), entry.Timestamp())
				}
				if string(entry.Message()) != string(entries[i].Message()) {
					t.Errorf("Entry %d: expected message %q, got %q", i, entries[i].Message(), entry.Message())
				}
				ingestible, ok := entry.(Ingestible)
				if !ok {
					t.Fatalf("Entry %d: expected Ingestible interface", i)
				}
				if ingestible.IngestionTimeUnixMilli() != entries[i].(SimpleEntry).IngestionTs {
					t.Errorf("Entry %d: expected ingestionTs %d, got %d",
						i, entries[i].(SimpleEntry).IngestionTs, ingestible.IngestionTimeUnixMilli())
				}
			}
		})
	}
}

func TestZeroIngestionTimeRoundTrip(t *testing.T) {
	tempDir := t.TempDir()

	entries := []Entry{
		SimpleEntry{Ts: 1000, Msg: []byte("zero ingestion")},
	}

	opts := &WriterOptions{
		ChunksDir: tempDir,
		Encoding:  EncodingZstd,
		ChunkSize: 10,
	}

	w := NewWriter(opts)
	if err := w.WriteBatch(entries); err != nil {
		t.Fatalf("WriteBatch failed: %v", err)
	}

	chunkPath, err := w.Flush()
	if err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	readEntries, err := Read(chunkPath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if len(readEntries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(readEntries))
	}

	ingestible, ok := readEntries[0].(Ingestible)
	if !ok {
		t.Fatal("Expected Ingestible interface even for zero IngestionTs")
	}
	if ingestible.IngestionTimeUnixMilli() != 0 {
		t.Errorf("Expected 0 ingestionTs, got %d", ingestible.IngestionTimeUnixMilli())
	}
}

func TestPebbleIndex(t *testing.T) {
	tempDir := t.TempDir()
	indexPath := filepath.Join(tempDir, "index.db")

	idx, err := NewPebbleIndex(&IndexOptions{
		Path: indexPath,
	})
	if err != nil {
		t.Fatalf("Failed to create index: %v", err)
	}
	defer idx.Close()

	t.Run("Add and Get", func(t *testing.T) {
		meta := Meta{
			ChunkID:    "test-chunk-1",
			MinTs:      1000,
			MaxTs:      2000,
			EntryCount: 10,
			ChunkPath:  "/path/to/chunk1.chunk",
			Tags: map[string]string{
				"database": "test-db",
				"table":    "test-table",
			},
		}

		if err := idx.Add(meta); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		retrieved, err := idx.Get("test-chunk-1")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		if retrieved.ChunkID != meta.ChunkID {
			t.Errorf("Expected ChunkID %q, got %q", meta.ChunkID, retrieved.ChunkID)
		}
		if retrieved.MinTs != meta.MinTs {
			t.Errorf("Expected MinTs %d, got %d", meta.MinTs, retrieved.MinTs)
		}
		if retrieved.MaxTs != meta.MaxTs {
			t.Errorf("Expected MaxTs %d, got %d", meta.MaxTs, retrieved.MaxTs)
		}
		if retrieved.EntryCount != meta.EntryCount {
			t.Errorf("Expected EntryCount %d, got %d", meta.EntryCount, retrieved.EntryCount)
		}
		if retrieved.ChunkPath != meta.ChunkPath {
			t.Errorf("Expected ChunkPath %q, got %q", meta.ChunkPath, retrieved.ChunkPath)
		}
		if retrieved.Tags["database"] != meta.Tags["database"] {
			t.Errorf("Expected database tag %q, got %q", meta.Tags["database"], retrieved.Tags["database"])
		}
	})

	t.Run("AddBatch", func(t *testing.T) {
		metas := []Meta{
			{
				ChunkID:    "test-chunk-2",
				MinTs:      3000,
				MaxTs:      4000,
				EntryCount: 20,
				ChunkPath:  "/path/to/chunk2.chunk",
			},
			{
				ChunkID:    "test-chunk-3",
				MinTs:      5000,
				MaxTs:      6000,
				EntryCount: 30,
				ChunkPath:  "/path/to/chunk3.chunk",
			},
		}

		if err := idx.AddBatch(metas); err != nil {
			t.Fatalf("AddBatch failed: %v", err)
		}

		for _, meta := range metas {
			retrieved, err := idx.Get(meta.ChunkID)
			if err != nil {
				t.Fatalf("Get failed for %s: %v", meta.ChunkID, err)
			}
			if retrieved.ChunkID != meta.ChunkID {
				t.Errorf("Expected ChunkID %q, got %q", meta.ChunkID, retrieved.ChunkID)
			}
		}
	})

	t.Run("Get non-existent", func(t *testing.T) {
		_, err := idx.Get("non-existent")
		if err != ErrChunkNotFound {
			t.Errorf("Expected ErrChunkNotFound, got %v", err)
		}
	})

	t.Run("QueryByTimeRange", func(t *testing.T) {
		metas, err := idx.QueryByTimeRange(2500, 5500)
		if err != nil {
			t.Fatalf("QueryByTimeRange failed: %v", err)
		}

		if len(metas) != 2 {
			t.Errorf("Expected 2 chunks, got %d", len(metas))
		}

		for _, meta := range metas {
			if meta.MaxTs < 2500 || meta.MinTs > 5500 {
				t.Errorf("Chunk %s does not overlap with query range: MinTs=%d, MaxTs=%d", meta.ChunkID, meta.MinTs, meta.MaxTs)
			}
		}
	})

	t.Run("QueryByTimeRange with MinTs outside range", func(t *testing.T) {
		wideChunk := Meta{
			ChunkID:    "wide-chunk",
			MinTs:      100,
			MaxTs:      10000,
			EntryCount: 50,
			ChunkPath:  "/path/to/wide.chunk",
		}
		if err := idx.Add(wideChunk); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		metas, err := idx.QueryByTimeRange(5000, 6000)
		if err != nil {
			t.Fatalf("QueryByTimeRange failed: %v", err)
		}

		found := false
		for _, meta := range metas {
			if meta.ChunkID == "wide-chunk" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected wide-chunk (MinTs=100, MaxTs=10000) to be found in query range 5000-6000")
		}

		if err := idx.Delete("wide-chunk"); err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}
	})

	t.Run("Count", func(t *testing.T) {
		count, err := idx.Count()
		if err != nil {
			t.Fatalf("Count failed: %v", err)
		}

		if count != 3 {
			t.Errorf("Expected 3 entries, got %d", count)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		if err := idx.Delete("test-chunk-1"); err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		_, err := idx.Get("test-chunk-1")
		if err != ErrChunkNotFound {
			t.Errorf("Expected ErrChunkNotFound after delete, got %v", err)
		}
	})

	t.Run("DeleteBefore", func(t *testing.T) {
		if count, err := idx.Count(); err != nil {
			t.Fatalf("Count failed: %v", err)
		} else if count != 2 {
			t.Fatalf("Expected 2 chunks before DeleteBefore, got %d", count)
		}

		metas, err := idx.QueryByTimeRange(0, 10000)
		if err != nil {
			t.Fatalf("QueryByTimeRange failed: %v", err)
		}

		if len(metas) != 2 {
			t.Fatalf("Expected 2 chunks before DeleteBefore, got %d", len(metas))
		}

		count, err := idx.DeleteBefore(5000)
		if err != nil {
			t.Fatalf("DeleteBefore failed: %v", err)
		}

		if count != 1 {
			t.Errorf("Expected to delete 1 entry, got %d", count)
		}

		remainingCount, err := idx.Count()
		if err != nil {
			t.Fatalf("Count failed: %v", err)
		} else if remainingCount != 1 {
			t.Errorf("Expected 1 remaining entry after DeleteBefore, got %d", remainingCount)
		}

		remaining, err := idx.QueryByTimeRange(0, 10000)
		if err != nil {
			t.Fatalf("QueryByTimeRange failed: %v", err)
		}

		if len(remaining) != 1 {
			t.Errorf("Expected 1 remaining entry, got %d", len(remaining))
		}

		if remaining[0].ChunkID != "test-chunk-3" {
			t.Errorf("Expected test-chunk-3, got %s", remaining[0].ChunkID)
		}
	})

	t.Run("Update", func(t *testing.T) {
		updateChunk := Meta{
			ChunkID:    "update-chunk",
			MinTs:      100,
			MaxTs:      200,
			EntryCount: 5,
			ChunkPath:  "/path/to/update.chunk",
			Tags: map[string]string{
				"key": "value1",
			},
		}
		if err := idx.Add(updateChunk); err != nil {
			t.Fatalf("Add failed: %v", err)
		}

		countBefore, _ := idx.Count()

		updatedChunk := Meta{
			ChunkID:    "update-chunk",
			MinTs:      50,
			MaxTs:      300,
			EntryCount: 10,
			ChunkPath:  "/path/to/update.chunk",
			Tags: map[string]string{
				"key": "value2",
			},
		}
		if err := idx.Update(updatedChunk); err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		countAfter, _ := idx.Count()
		if countAfter != countBefore {
			t.Errorf("Update should not change entry count: before=%d, after=%d", countBefore, countAfter)
		}

		retrieved, err := idx.Get("update-chunk")
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}
		if retrieved.MinTs != 50 {
			t.Errorf("Expected MinTs 50, got %d", retrieved.MinTs)
		}
		if retrieved.MaxTs != 300 {
			t.Errorf("Expected MaxTs 300, got %d", retrieved.MaxTs)
		}
		if retrieved.Tags["key"] != "value2" {
			t.Errorf("Expected tag key=value2, got %s", retrieved.Tags["key"])
		}

		metas, err := idx.QueryByTimeRange(60, 70)
		if err != nil {
			t.Fatalf("QueryByTimeRange failed: %v", err)
		}
		found := false
		for _, m := range metas {
			if m.ChunkID == "update-chunk" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Update should update time index: expected to find chunk in range 60-70")
		}

		if err := idx.Delete("update-chunk"); err != nil {
			t.Fatalf("Cleanup failed: %v", err)
		}
	})
}
