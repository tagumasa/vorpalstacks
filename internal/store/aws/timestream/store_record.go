package timestream

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/chunk"
)

// RecordStore manages Timestream records.
type RecordStore struct {
	tableStore *TableStore
	region     string
	dataPath   string
	chunkSize  int
	index      *chunk.PebbleIndex
	useIndex   bool
	chunkMu    sync.Map
	buffers    map[string]*chunkBuffer
	bufferMu   sync.Mutex
	bufferSize int
}

// NewRecordStore creates a new record store.
func NewRecordStore(store storage.BasicStorage, tableStore *TableStore, region, dataPath string) *RecordStore {
	s := &RecordStore{
		tableStore: tableStore,
		region:     region,
		dataPath:   dataPath,
		chunkSize:  16 * 1024 * 1024,
		useIndex:   false,
		buffers:    make(map[string]*chunkBuffer),
		bufferSize: 1000,
	}
	go s.cleanupChunkLocks()
	return s
}

// NewRecordStoreWithIndex creates a new record store with an index.
func NewRecordStoreWithIndex(store storage.BasicStorage, tableStore *TableStore, region, dataPath string) (*RecordStore, error) {
	indexPath := filepath.Join(dataPath, region, "timestream_chunks", "index.db")
	idx, err := chunk.GetDefaultIndexManager().GetIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	s := &RecordStore{
		tableStore: tableStore,
		region:     region,
		dataPath:   dataPath,
		chunkSize:  16 * 1024 * 1024,
		index:      idx,
		useIndex:   true,
		buffers:    make(map[string]*chunkBuffer),
		bufferSize: 1000,
	}
	go s.cleanupChunkLocks()
	return s, nil
}

func convertDimensions(dimensions []Dimension) []chunk.Dimension {
	result := make([]chunk.Dimension, len(dimensions))
	for i, d := range dimensions {
		result[i] = chunk.Dimension{Name: d.Name, Value: d.Value}
	}
	return result
}

func convertMeasureValues(measureValues []MeasureValue) []chunk.MeasureValue {
	result := make([]chunk.MeasureValue, len(measureValues))
	for i, mv := range measureValues {
		result[i] = chunk.MeasureValue{Name: mv.Name, Value: mv.Value, Type: string(mv.Type)}
	}
	return result
}

func sanitizePathComponent(name string) string {
	safe := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, name)
	if safe == "" {
		safe = "unnamed"
	}
	return safe
}

func (s *RecordStore) getChunkPath(databaseName, tableName string, ts time.Time) string {
	hour := ts.Format("2006-01-02-15")
	return fmt.Sprintf("%s/timestream_chunks/%s/%s/%s/%s.chunk", s.dataPath, s.region, sanitizePathComponent(databaseName), sanitizePathComponent(tableName), hour)
}

func (s *RecordStore) getChunkLock(chunkPath string) *sync.Mutex {
	mu, _ := s.chunkMu.LoadOrStore(chunkPath, &sync.Mutex{})
	if typed, ok := mu.(*sync.Mutex); ok {
		return typed
	}
	return &sync.Mutex{}
}

func (s *RecordStore) cleanupChunkLocks() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for range ticker.C {
		cutoff := time.Now().Add(-2 * time.Hour).Format("2006-01-02-15")
		s.chunkMu.Range(func(key, _ interface{}) bool {
			path := key.(string)
			idx := strings.LastIndex(path, "/")
			if idx >= 0 {
				hour := strings.TrimSuffix(path[idx+1:], ".chunk")
				if hour < cutoff {
					s.chunkMu.Delete(key)
				}
			}
			return true
		})
	}
}

func (s *RecordStore) writeChunk(chunkPath string, entry *chunk.TimestreamEntry) error {
	mu := s.getChunkLock(chunkPath)
	mu.Lock()
	defer mu.Unlock()

	chunkDir := filepath.Dir(chunkPath)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return err
	}

	entries, err := s.readChunk(chunkPath)
	if err != nil {
		entries = []chunk.TimestreamEntry{*entry}
	} else {
		entries = append(entries, *entry)
	}

	chunkEntries := make([]chunk.Entry, len(entries))
	for i := range entries {
		chunkEntries[i] = &entries[i]
	}

	writer := chunk.NewWriter(&chunk.WriterOptions{
		ChunksDir:   chunkDir,
		ChunkSize:   s.chunkSize,
		Encoding:    chunk.EncodingZstd,
		SyncOnWrite: true,
	})

	if err := writer.WriteBatch(chunkEntries); err != nil {
		return err
	}

	if _, err := writer.Flush(); err != nil {
		return err
	}

	if s.useIndex && s.index != nil {
		minTs := entries[0].Time.UnixNano()
		maxTs := entries[0].Time.UnixNano()
		for _, e := range entries {
			ts := e.Time.UnixNano()
			if ts < minTs {
				minTs = ts
			}
			if ts > maxTs {
				maxTs = ts
			}
		}

		chunkID := filepath.Base(chunkPath)
		meta := chunk.Meta{
			ChunkID:    chunkID,
			MinTs:      minTs,
			MaxTs:      maxTs,
			EntryCount: len(entries),
			ChunkPath:  chunkPath,
			Tags: map[string]string{
				"database": entry.DatabaseName,
				"table":    entry.TableName,
			},
		}

		if err := s.index.Update(meta); err != nil {
			return fmt.Errorf("failed to update chunk meta in index: %w", err)
		}
	}

	return nil
}

func (s *RecordStore) readChunk(chunkPath string) ([]chunk.TimestreamEntry, error) {
	entries, err := chunk.Read(chunkPath)
	if err != nil {
		return nil, err
	}

	result := make([]chunk.TimestreamEntry, len(entries))
	for i, e := range entries {
		var te chunk.TimestreamEntry
		if err := json.Unmarshal(e.Message(), &te); err != nil {
			return nil, err
		}
		result[i] = te
	}

	return result, nil
}

func (s *RecordStore) parseTimestamp(timeStr string, timeUnit TimeUnit) (time.Time, error) {
	if timeStr == "" {
		return time.Now().UTC(), nil
	}

	var ts int64
	if _, err := fmt.Sscanf(timeStr, "%d", &ts); err != nil {
		return time.Time{}, ErrInvalidParameter
	}

	switch timeUnit {
	case TimeUnitSeconds:
		return time.Unix(ts, 0).UTC(), nil
	case TimeUnitMilliseconds:
		return time.UnixMilli(ts).UTC(), nil
	case TimeUnitMicroseconds:
		return time.UnixMicro(ts).UTC(), nil
	case TimeUnitNanoseconds:
		return time.Unix(0, ts).UTC(), nil
	default:
		return time.UnixMilli(ts).UTC(), nil
	}
}

// StoredRecord represents a stored Timestream record.
type StoredRecord struct {
	DatabaseName     string           `json:"databaseName"`
	TableName        string           `json:"tableName"`
	Dimensions       []Dimension      `json:"dimensions,omitempty"`
	MeasureName      string           `json:"measureName"`
	MeasureValue     string           `json:"measureValue,omitempty"`
	MeasureValueType MeasureValueType `json:"measureValueType,omitempty"`
	MeasureValues    []MeasureValue   `json:"measureValues,omitempty"`
	Timestamp        time.Time        `json:"timestamp"`
	Version          int64            `json:"version,omitempty"`
}

// WriteRecords writes records to the store.
func (s *RecordStore) WriteRecords(databaseName, tableName string, records []Record) ([]RejectedRecord, error) {
	if _, err := s.tableStore.GetTable(databaseName, tableName); err != nil {
		return nil, ErrTableNotFound
	}

	var rejectedRecords []RejectedRecord
	chunkEntries := make(map[string][]indexedEntry)

	for i, record := range records {
		if record.MeasureName == "" {
			rejectedRecords = append(rejectedRecords, RejectedRecord{
				RecordIndex: int64(i),
				Reason:      "MeasureName is required",
			})
			continue
		}

		if record.MeasureValueType == MeasureValueTypeMulti {
			if len(record.MeasureValues) == 0 {
				rejectedRecords = append(rejectedRecords, RejectedRecord{
					RecordIndex: int64(i),
					Reason:      "MeasureValues is required when MeasureValueType is MULTI",
				})
				continue
			}
		} else {
			if record.MeasureValue == "" {
				rejectedRecords = append(rejectedRecords, RejectedRecord{
					RecordIndex: int64(i),
					Reason:      "MeasureValue is required",
				})
				continue
			}
		}

		timeUnit := record.TimeUnit
		if timeUnit == "" {
			timeUnit = TimeUnitMilliseconds
		}

		ts, err := s.parseTimestamp(record.Time, timeUnit)
		if err != nil {
			rejectedRecords = append(rejectedRecords, RejectedRecord{
				RecordIndex: int64(i),
				Reason:      "Invalid timestamp",
			})
			continue
		}

		version := record.Version
		if version == 0 {
			version = 1
		}

		entry := &chunk.TimestreamEntry{
			DatabaseName:     databaseName,
			TableName:        tableName,
			Dimensions:       convertDimensions(record.Dimensions),
			MeasureName:      record.MeasureName,
			MeasureValue:     record.MeasureValue,
			MeasureValueType: string(record.MeasureValueType),
			MeasureValues:    convertMeasureValues(record.MeasureValues),
			Time:             ts,
			Version:          version,
		}

		chunkPath := s.getChunkPath(databaseName, tableName, ts)
		chunkEntries[chunkPath] = append(chunkEntries[chunkPath], indexedEntry{entry: entry, originalIdx: i})
	}

	for chunkPath, indexedEntries := range chunkEntries {
		entries := make([]*chunk.TimestreamEntry, len(indexedEntries))
		for j, ie := range indexedEntries {
			entries[j] = ie.entry
		}
		if err := s.writeChunkBatch(chunkPath, entries); err != nil {
			for _, ie := range indexedEntries {
				rejectedRecords = append(rejectedRecords, RejectedRecord{
					RecordIndex: int64(ie.originalIdx),
					Reason:      err.Error(),
				})
			}
		}
	}

	return rejectedRecords, nil
}

func (s *RecordStore) writeChunkBatch(chunkPath string, newEntries []*chunk.TimestreamEntry) error {
	mu := s.getChunkLock(chunkPath)
	mu.Lock()
	defer mu.Unlock()

	s.bufferMu.Lock()
	buf, ok := s.buffers[chunkPath]
	if !ok {
		buf = &chunkBuffer{
			entries:  make([]chunk.TimestreamEntry, 0, s.bufferSize),
			seenKeys: make(map[string]int64),
		}
		s.buffers[chunkPath] = buf
	}
	s.bufferMu.Unlock()

	for _, entry := range newEntries {
		key := fmt.Sprintf("%s|%s|%s|%d",
			entry.DatabaseName,
			entry.BuildDimensionsKey(),
			entry.MeasureName,
			entry.Time.UnixNano())

		if existingVersion, exists := buf.seenKeys[key]; exists {
			if existingVersion != 0 && existingVersion != entry.Version {
				continue
			}
		}
		buf.entries = append(buf.entries, *entry)
		buf.seenKeys[key] = entry.Version
	}

	if len(buf.entries) >= s.bufferSize {
		return s.flushBuffer(chunkPath, buf)
	}

	return nil
}

func (s *RecordStore) flushBuffer(chunkPath string, buf *chunkBuffer) error {
	if len(buf.entries) == 0 {
		return nil
	}

	chunkDir := filepath.Dir(chunkPath)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return err
	}

	existingEntries, _ := s.readChunk(chunkPath)
	allEntries := make([]chunk.TimestreamEntry, 0, len(existingEntries)+len(buf.entries))
	allEntries = append(allEntries, existingEntries...)
	allEntries = append(allEntries, buf.entries...)

	chunkEntries := make([]chunk.Entry, len(allEntries))
	for i := range allEntries {
		chunkEntries[i] = &allEntries[i]
	}

	writer := chunk.NewWriter(&chunk.WriterOptions{
		ChunksDir:   chunkDir,
		ChunkSize:   s.chunkSize,
		Encoding:    chunk.EncodingZstd,
		SyncOnWrite: false,
	})

	if err := writer.WriteBatch(chunkEntries); err != nil {
		return err
	}

	if _, err := writer.Flush(); err != nil {
		return err
	}

	if s.useIndex && s.index != nil && len(allEntries) > 0 {
		minTs := allEntries[0].Time.UnixNano()
		maxTs := allEntries[0].Time.UnixNano()
		for _, e := range allEntries {
			ts := e.Time.UnixNano()
			if ts < minTs {
				minTs = ts
			}
			if ts > maxTs {
				maxTs = ts
			}
		}

		chunkID := filepath.Base(chunkPath)
		meta := chunk.Meta{
			ChunkID:    chunkID,
			MinTs:      minTs,
			MaxTs:      maxTs,
			EntryCount: len(allEntries),
			ChunkPath:  chunkPath,
			Tags: map[string]string{
				"database": allEntries[0].DatabaseName,
				"table":    allEntries[0].TableName,
			},
		}

		if err := s.index.Update(meta); err != nil {
			return fmt.Errorf("failed to update chunk meta in index: %w", err)
		}
	}

	buf.entries = buf.entries[:0]
	for k := range buf.seenKeys {
		delete(buf.seenKeys, k)
	}

	return nil
}

// FlushAllBuffers flushes all buffered records to storage.
func (s *RecordStore) FlushAllBuffers() error {
	s.bufferMu.Lock()
	defer s.bufferMu.Unlock()

	var lastErr error
	for chunkPath, buf := range s.buffers {
		if err := s.flushBuffer(chunkPath, buf); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// QueryRecords retrieves records within a time range.
func (s *RecordStore) QueryRecords(databaseName, tableName string, startTime, endTime time.Time) ([]*StoredRecord, error) {
	var records []*StoredRecord

	if s.useIndex && s.index != nil {
		minTs := startTime.UnixNano()
		maxTs := endTime.UnixNano()

		metas, err := s.index.QueryByTimeRange(minTs, maxTs)
		if err != nil {
			return nil, fmt.Errorf("failed to query chunks by time range: %w", err)
		}

		for _, meta := range metas {
			if meta.Tags["database"] != databaseName || meta.Tags["table"] != tableName {
				continue
			}

			entries, err := s.readChunk(meta.ChunkPath)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if !startTime.IsZero() && entry.Time.Before(startTime) {
					continue
				}
				if !endTime.IsZero() && entry.Time.After(endTime) {
					continue
				}

				record := &StoredRecord{
					DatabaseName:     entry.DatabaseName,
					TableName:        entry.TableName,
					Dimensions:       convertDimensionsBack(entry.Dimensions),
					MeasureName:      entry.MeasureName,
					MeasureValue:     entry.MeasureValue,
					MeasureValueType: MeasureValueType(entry.MeasureValueType),
					MeasureValues:    convertMeasureValuesBack(entry.MeasureValues),
					Timestamp:        entry.Time,
					Version:          entry.Version,
				}
				records = append(records, record)
			}
		}
	} else {
		chunkDir := fmt.Sprintf("%s/timestream_chunks/%s/%s/%s", s.dataPath, s.region, sanitizePathComponent(databaseName), sanitizePathComponent(tableName))

		startHour := startTime.Format("2006-01-02-15")
		endHour := endTime.Add(1 * time.Hour).Format("2006-01-02-15")

		for h := startTime.Truncate(time.Hour); h.Before(endTime.Add(1 * time.Hour)); h = h.Add(1 * time.Hour) {
			hourStr := h.Format("2006-01-02-15")
			if hourStr < startHour || hourStr > endHour {
				continue
			}

			chunkPath := fmt.Sprintf("%s/%s.chunk", chunkDir, hourStr)
			entries, err := s.readChunk(chunkPath)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if entry.DatabaseName != databaseName || entry.TableName != tableName {
					continue
				}

				if !startTime.IsZero() && entry.Time.Before(startTime) {
					continue
				}
				if !endTime.IsZero() && entry.Time.After(endTime) {
					continue
				}

				record := &StoredRecord{
					DatabaseName:     entry.DatabaseName,
					TableName:        entry.TableName,
					Dimensions:       convertDimensionsBack(entry.Dimensions),
					MeasureName:      entry.MeasureName,
					MeasureValue:     entry.MeasureValue,
					MeasureValueType: MeasureValueType(entry.MeasureValueType),
					MeasureValues:    convertMeasureValuesBack(entry.MeasureValues),
					Timestamp:        entry.Time,
					Version:          entry.Version,
				}
				records = append(records, record)
			}
		}
	}

	return records, nil
}

func convertDimensionsBack(dimensions []chunk.Dimension) []Dimension {
	result := make([]Dimension, len(dimensions))
	for i, d := range dimensions {
		result[i] = Dimension{Name: d.Name, Value: d.Value}
	}
	return result
}

func convertMeasureValuesBack(measureValues []chunk.MeasureValue) []MeasureValue {
	result := make([]MeasureValue, len(measureValues))
	for i, mv := range measureValues {
		result[i] = MeasureValue{Name: mv.Name, Value: mv.Value, Type: MeasureValueType(mv.Type)}
	}
	return result
}
