package chunk

import (
	"os"
	"path/filepath"
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
			ChunksDir: chunksDir,
			Encoding:  EncodingGzip,
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
			t.Fatal("Expected chunk path to be returned")
		}

		data, err := os.ReadFile(chunkPath)
		if err != nil {
			t.Fatalf("ReadFile failed: %v", err)
		}

		if len(data) < 2 || data[0] != 0x1f || data[1] != 0x8b {
			t.Errorf("Expected gzip magic bytes at start of file")
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
