package cloudwatch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/core/storage/chunk"
	"vorpalstacks/internal/store/aws/common"
)

func metricBucketName(region string) string {
	return "cw_metrics-" + region
}

const metricBufferSize = 1000

type metricBuffer struct {
	entries []chunk.CloudWatchMetricEntry
}

// MetricChunkStore provides CloudWatch metric chunk storage operations.
type MetricChunkStore struct {
	*common.BaseStore
	storage    storage.BasicStorage
	region     string
	dataPath   string
	index      *chunk.PebbleIndex
	useIndex   bool
	keyLocker  common.KeyLocker
	chunkPaths sync.Map
	buffers    map[string]*metricBuffer
	bufferMu   sync.Mutex
	stopCh     chan struct{}
	closeOnce  sync.Once
}

// NewMetricChunkStore creates a new CloudWatch metric chunk store.
func NewMetricChunkStore(store storage.BasicStorage, region, dataPath string) *MetricChunkStore {
	s := &MetricChunkStore{
		BaseStore: common.NewBaseStore(store.Bucket(metricBucketName(region)), "cloudwatch"),
		storage:   store,
		region:    region,
		dataPath:  dataPath,
		useIndex:  false,
		buffers:   make(map[string]*metricBuffer),
		stopCh:    make(chan struct{}),
	}
	go s.cleanupChunkLocks()
	return s
}

// NewMetricChunkStoreWithIndex creates a new CloudWatch metric chunk store with an index.
func NewMetricChunkStoreWithIndex(store storage.BasicStorage, region, dataPath string) (*MetricChunkStore, error) {
	indexPath := filepath.Join(dataPath, region, "cw_metric_index", "index.db")
	idx, err := chunk.GetDefaultIndexManager().GetIndex(indexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric index: %w", err)
	}

	s := &MetricChunkStore{
		BaseStore: common.NewBaseStore(store.Bucket(metricBucketName(region)), "cloudwatch"),
		storage:   store,
		region:    region,
		dataPath:  dataPath,
		index:     idx,
		useIndex:  true,
		buffers:   make(map[string]*metricBuffer),
		stopCh:    make(chan struct{}),
	}
	go s.cleanupChunkLocks()
	return s, nil
}

// Close stops the background cleanup goroutine.
func (s *MetricChunkStore) Close() {
	s.closeOnce.Do(func() { close(s.stopCh) })
}

func (s *MetricChunkStore) getChunkPath(namespace, metricName string, ts time.Time) string {
	safeNS := strings.ReplaceAll(namespace, "/", "_")
	safeMetric := strings.ReplaceAll(metricName, "/", "_")
	hour := ts.Format("2006-01-02-15")
	return fmt.Sprintf("%s/%s/cw_metric_chunks/%s/%s/%s.chunk", s.dataPath, s.region, safeNS, safeMetric, hour)
}

func (s *MetricChunkStore) cleanupChunkLocks() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
		}
		cutoff := time.Now().Add(-2 * time.Hour).Format("2006-01-02-15")
		s.keyLocker.Range(func(path string) bool {
			idx := strings.LastIndex(path, "/")
			if idx >= 0 {
				hour := strings.TrimSuffix(path[idx+1:], ".chunk")
				if hour < cutoff {
					s.keyLocker.Delete(path)
				}
			}
			return true
		})
		s.chunkPaths.Range(func(key, _ interface{}) bool {
			path := key.(string)
			idx := strings.LastIndex(path, "/")
			if idx >= 0 {
				hour := strings.TrimSuffix(path[idx+1:], ".chunk")
				if hour < cutoff {
					s.chunkPaths.Delete(key)
				}
			}
			return true
		})
	}
}

func (s *MetricChunkStore) writeChunk(chunkPath string, entry *chunk.CloudWatchMetricEntry) error {
	s.bufferMu.Lock()
	buf, ok := s.buffers[chunkPath]
	if !ok {
		buf = &metricBuffer{entries: make([]chunk.CloudWatchMetricEntry, 0, metricBufferSize)}
		s.buffers[chunkPath] = buf
	}
	buf.entries = append(buf.entries, *entry)
	shouldFlush := len(buf.entries) >= metricBufferSize
	s.bufferMu.Unlock()

	if shouldFlush {
		return s.flushBuffer(chunkPath)
	}
	return nil
}

func (s *MetricChunkStore) flushBuffer(chunkPath string) error {
	s.bufferMu.Lock()
	buf := s.buffers[chunkPath]
	if buf == nil || len(buf.entries) == 0 {
		s.bufferMu.Unlock()
		return nil
	}
	entries := buf.entries
	buf.entries = make([]chunk.CloudWatchMetricEntry, 0, metricBufferSize)
	s.bufferMu.Unlock()

	s.keyLocker.Lock(chunkPath)
	defer s.keyLocker.Unlock(chunkPath)

	chunkDir := filepath.Dir(chunkPath)
	if err := os.MkdirAll(chunkDir, 0755); err != nil {
		return err
	}

	var existingEntries []chunk.CloudWatchMetricEntry
	if v, ok := s.chunkPaths.Load(chunkPath); ok {
		if readEntries, err := s.readChunkFile(v.(string)); err == nil {
			existingEntries = readEntries
		}
	}

	allEntries := make([]chunk.CloudWatchMetricEntry, 0, len(existingEntries)+len(entries))
	allEntries = append(allEntries, existingEntries...)
	allEntries = append(allEntries, entries...)

	chunkEntries := make([]chunk.Entry, len(allEntries))
	for i := range allEntries {
		chunkEntries[i] = &allEntries[i]
	}

	writer := chunk.NewWriter(&chunk.WriterOptions{
		ChunksDir:   chunkDir,
		Encoding:    chunk.EncodingZstd,
		SyncOnWrite: true,
	})

	if err := writer.WriteBatch(chunkEntries); err != nil {
		return err
	}

	actualChunkPath, err := writer.Flush()
	if err != nil {
		return err
	}

	s.chunkPaths.Store(chunkPath, actualChunkPath)

	if existingPath, ok := s.chunkPaths.Load(chunkPath); ok {
		if existingStr, ok := existingPath.(string); ok && existingStr != "" && existingStr != actualChunkPath {
			os.Remove(existingStr)
		}
	}

	if s.useIndex && s.index != nil {
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

		chunkID := filepath.Base(actualChunkPath)
		meta := chunk.Meta{
			ChunkID:    chunkID,
			MinTs:      minTs,
			MaxTs:      maxTs,
			EntryCount: len(allEntries),
			ChunkPath:  actualChunkPath,
			Tags: map[string]string{
				"namespace":   entries[0].Namespace,
				"metric_name": entries[0].MetricName,
				"dimensions":  entries[0].BuildDimensionsKey(),
			},
		}

		if err := s.index.Update(meta); err != nil {
			return fmt.Errorf("failed to update chunk meta in index: %w", err)
		}
	}

	return nil
}

// FlushAllBuffers forces all buffered metric entries to be written to disk.
func (s *MetricChunkStore) FlushAllBuffers() error {
	s.bufferMu.Lock()
	paths := make([]string, 0, len(s.buffers))
	for p := range s.buffers {
		paths = append(paths, p)
	}
	s.bufferMu.Unlock()

	for _, p := range paths {
		if err := s.flushBuffer(p); err != nil {
			return err
		}
	}
	return nil
}

func (s *MetricChunkStore) readChunk(chunkPath string) ([]chunk.CloudWatchMetricEntry, error) {
	if v, ok := s.chunkPaths.Load(chunkPath); ok {
		return s.readChunkFile(v.(string))
	}
	return s.readChunkFile(chunkPath)
}

func (s *MetricChunkStore) readChunkFile(path string) ([]chunk.CloudWatchMetricEntry, error) {
	entries, err := chunk.Read(path)
	if err != nil {
		return nil, err
	}

	result := make([]chunk.CloudWatchMetricEntry, len(entries))
	for i, e := range entries {
		var ce chunk.CloudWatchMetricEntry
		if err := json.Unmarshal(e.Message(), &ce); err != nil {
			return nil, err
		}
		result[i] = ce
	}

	return result, nil
}

func convertToChunkDimensions(dims []Dimension) []chunk.CloudWatchDimension {
	result := make([]chunk.CloudWatchDimension, len(dims))
	for i, d := range dims {
		result[i] = chunk.CloudWatchDimension{Name: d.Name, Value: d.Value}
	}
	return result
}

func convertToChunkStatisticSet(ss *StatisticSet) *chunk.CloudWatchStatisticSet {
	if ss == nil {
		return nil
	}
	return &chunk.CloudWatchStatisticSet{
		Maximum:     ss.Maximum,
		Minimum:     ss.Minimum,
		SampleCount: ss.SampleCount,
		Sum:         ss.Sum,
	}
}

func convertFromChunkDimensions(dims []chunk.CloudWatchDimension) []Dimension {
	result := make([]Dimension, len(dims))
	for i, d := range dims {
		result[i] = Dimension{Name: d.Name, Value: d.Value}
	}
	return result
}

func convertFromChunkStatisticSet(ss *chunk.CloudWatchStatisticSet) *StatisticSet {
	if ss == nil {
		return nil
	}
	return &StatisticSet{
		Maximum:     ss.Maximum,
		Minimum:     ss.Minimum,
		SampleCount: ss.SampleCount,
		Sum:         ss.Sum,
	}
}

// PutMetricData stores metric data in the CloudWatch store.
func (s *MetricChunkStore) PutMetricData(namespace string, metricData []MetricDatum) error {
	for _, datum := range metricData {
		ts := datum.Timestamp
		if ts.IsZero() {
			ts = time.Now().UTC()
		}

		entry := &chunk.CloudWatchMetricEntry{
			Namespace:         namespace,
			MetricName:        datum.MetricName,
			Dimensions:        convertToChunkDimensions(datum.Dimensions),
			Time:              ts,
			Value:             datum.Value,
			Values:            datum.Values,
			Counts:            datum.Counts,
			Unit:              string(datum.Unit),
			StorageResolution: datum.StorageResolution,
			StatisticValues:   convertToChunkStatisticSet(datum.StatisticValues),
		}

		chunkPath := s.getChunkPath(namespace, datum.MetricName, ts)
		if err := s.writeChunk(chunkPath, entry); err != nil {
			return err
		}
	}
	return s.FlushAllBuffers()
}

// MetricDataPoint represents a CloudWatch metric data point.
type MetricDataPoint struct {
	MetricName        string
	Namespace         string
	Dimensions        []Dimension
	Timestamp         time.Time
	Value             float64
	Values            []float64
	Counts            []float64
	Unit              StandardUnit
	StorageResolution int32
	StatisticValues   *StatisticSet
}

// MetricQuery represents a query for CloudWatch metric data.
type MetricQuery struct {
	Namespace  string
	MetricName string
	Dimensions []Dimension
	StartTime  time.Time
	EndTime    time.Time
	Period     int32
	Statistics []string
}

// GetMetricStatistics retrieves statistics for a metric query.
func (s *MetricChunkStore) GetMetricStatistics(query MetricQuery) ([]*MetricStatistics, error) {
	if query.StartTime.IsZero() || query.EndTime.IsZero() {
		return nil, ErrInvalidParameter
	}

	var dataPoints []*MetricDataPoint

	if s.useIndex && s.index != nil {
		minTs := query.StartTime.UnixNano()
		maxTs := query.EndTime.UnixNano()

		metas, err := s.index.QueryByTimeRange(minTs, maxTs)
		if err != nil {
			return nil, fmt.Errorf("failed to query chunks by time range: %w", err)
		}

		dimsKey := s.buildDimensionsKey(query.Dimensions)
		for _, meta := range metas {
			if meta.Tags["namespace"] != query.Namespace {
				continue
			}
			if query.MetricName != "" && meta.Tags["metric_name"] != query.MetricName {
				continue
			}
			if len(query.Dimensions) > 0 && meta.Tags["dimensions"] != dimsKey {
				continue
			}

			entries, err := s.readChunkFile(meta.ChunkPath)
			if err != nil {
				continue
			}

			for _, entry := range entries {
				if (entry.Time.Equal(query.StartTime) || entry.Time.After(query.StartTime)) &&
					(entry.Time.Equal(query.EndTime) || entry.Time.Before(query.EndTime)) {
					dataPoints = append(dataPoints, &MetricDataPoint{
						MetricName:        entry.MetricName,
						Namespace:         entry.Namespace,
						Dimensions:        convertFromChunkDimensions(entry.Dimensions),
						Timestamp:         entry.Time,
						Value:             entry.Value,
						Values:            entry.Values,
						Counts:            entry.Counts,
						Unit:              StandardUnit(entry.Unit),
						StorageResolution: entry.StorageResolution,
						StatisticValues:   convertFromChunkStatisticSet(entry.StatisticValues),
					})
				}
			}
		}
	} else {
		chunkDir := fmt.Sprintf("%s/%s/cw_metric_chunks", s.dataPath, s.region)
		safeNS := strings.ReplaceAll(query.Namespace, "/", "_")
		nsDir := filepath.Join(chunkDir, safeNS)

		if query.MetricName != "" {
			safeMetric := strings.ReplaceAll(query.MetricName, "/", "_")
			metricDir := filepath.Join(nsDir, safeMetric)
			s.scanMetricChunks(metricDir, query, &dataPoints)
		} else {
			entries, err := os.ReadDir(nsDir)
			if err != nil {
				return nil, nil
			}
			for _, entry := range entries {
				if entry.IsDir() {
					metricDir := filepath.Join(nsDir, entry.Name())
					s.scanMetricChunks(metricDir, query, &dataPoints)
				}
			}
		}
	}

	if len(dataPoints) == 0 {
		return nil, nil
	}

	return s.aggregateStatistics(dataPoints, query), nil
}

func (s *MetricChunkStore) scanMetricChunks(metricDir string, query MetricQuery, dataPoints *[]*MetricDataPoint) {
	entries, err := os.ReadDir(metricDir)
	if err != nil {
		return
	}

	dimsKey := s.buildDimensionsKey(query.Dimensions)
	for _, entry := range entries {
		if !entry.IsDir() {
			chunkPath := filepath.Join(metricDir, entry.Name())
			chunkEntries, err := s.readChunk(chunkPath)
			if err != nil {
				continue
			}

			for _, ce := range chunkEntries {
				if ce.Namespace != query.Namespace {
					continue
				}
				if query.MetricName != "" && ce.MetricName != query.MetricName {
					continue
				}
				if len(query.Dimensions) > 0 && ce.BuildDimensionsKey() != dimsKey {
					continue
				}
				if (ce.Time.Equal(query.StartTime) || ce.Time.After(query.StartTime)) &&
					(ce.Time.Equal(query.EndTime) || ce.Time.Before(query.EndTime)) {
					*dataPoints = append(*dataPoints, &MetricDataPoint{
						MetricName:        ce.MetricName,
						Namespace:         ce.Namespace,
						Dimensions:        convertFromChunkDimensions(ce.Dimensions),
						Timestamp:         ce.Time,
						Value:             ce.Value,
						Values:            ce.Values,
						Counts:            ce.Counts,
						Unit:              StandardUnit(ce.Unit),
						StorageResolution: ce.StorageResolution,
						StatisticValues:   convertFromChunkStatisticSet(ce.StatisticValues),
					})
				}
			}
		}
	}
}

func (s *MetricChunkStore) aggregateStatistics(dataPoints []*MetricDataPoint, query MetricQuery) []*MetricStatistics {
	periodDuration := time.Duration(query.Period) * time.Second
	periodGroups := make(map[int64][]*MetricDataPoint)

	for _, dp := range dataPoints {
		periodStart := dp.Timestamp.Truncate(periodDuration)
		periodKey := periodStart.Unix()
		periodGroups[periodKey] = append(periodGroups[periodKey], dp)
	}

	var results []*MetricStatistics
	for periodKey, points := range periodGroups {
		stats := &MetricStatistics{
			Timestamp: time.Unix(periodKey, 0),
		}

		var sum, min, max float64
		var count float64
		first := true

		for _, dp := range points {
			if dp.StatisticValues != nil {
				sv := dp.StatisticValues
				if first || sv.Minimum < min {
					min = sv.Minimum
				}
				if first || sv.Maximum > max {
					max = sv.Maximum
				}
				sum += sv.Sum
				count += sv.SampleCount
				first = false
			} else {
				values := dp.Values
				counts := dp.Counts

				if len(values) == 0 {
					values = []float64{dp.Value}
				}

				for i, v := range values {
					c := float64(1)
					if i < len(counts) {
						c = counts[i]
					}
					if first || v < min {
						min = v
					}
					if first || v > max {
						max = v
					}
					sum += v * c
					count += c
					first = false
				}
			}

			if dp.Unit != "" {
				stats.Unit = dp.Unit
			}
		}

		for _, stat := range query.Statistics {
			switch strings.ToLower(stat) {
			case "sum":
				stats.Sum = sum
			case "average":
				if count > 0 {
					stats.Average = sum / count
				}
			case "minimum":
				stats.Minimum = min
			case "maximum":
				stats.Maximum = max
			case "samplecount":
				stats.SampleCount = count
			}
		}

		results = append(results, stats)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.Before(results[j].Timestamp)
	})

	return results
}

func (s *MetricChunkStore) buildDimensionsKey(dimensions []Dimension) string {
	if len(dimensions) == 0 {
		return ""
	}
	parts := make([]string, len(dimensions))
	for i, d := range dimensions {
		parts[i] = fmt.Sprintf("%s=%s", d.Name, d.Value)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}

// MetricStatistics represents CloudWatch metric statistics.
type MetricStatistics struct {
	Timestamp   time.Time
	Sum         float64
	Average     float64
	Minimum     float64
	Maximum     float64
	SampleCount float64
	Unit        StandardUnit
}

// ListMetrics returns a list of metrics matching the specified criteria.
func (s *MetricChunkStore) ListMetrics(namespace, metricName string, dimensions []Dimension) ([]MetricDatum, error) {
	var metrics []MetricDatum
	seen := make(map[string]bool)

	chunkDir := fmt.Sprintf("%s/%s/cw_metric_chunks", s.dataPath, s.region)
	s.listMetricsFromDir(chunkDir, namespace, metricName, dimensions, seen, &metrics)

	return metrics, nil
}

func (s *MetricChunkStore) listMetricsFromDir(chunkDir, namespace, metricName string, dimensions []Dimension, seen map[string]bool, metrics *[]MetricDatum) {
	nsEntries, err := os.ReadDir(chunkDir)
	if err != nil {
		return
	}

	for _, nsEntry := range nsEntries {
		if !nsEntry.IsDir() {
			continue
		}

		metricEntries, err := os.ReadDir(filepath.Join(chunkDir, nsEntry.Name()))
		if err != nil {
			continue
		}

		for _, metricEntry := range metricEntries {
			if !metricEntry.IsDir() {
				continue
			}

			hourEntries, err := os.ReadDir(filepath.Join(chunkDir, nsEntry.Name(), metricEntry.Name()))
			if err != nil {
				continue
			}

			for _, hourEntry := range hourEntries {
				if hourEntry.IsDir() {
					continue
				}

				chunkPath := filepath.Join(chunkDir, nsEntry.Name(), metricEntry.Name(), hourEntry.Name())
				entries, err := s.readChunk(chunkPath)
				if err != nil {
					continue
				}

				for _, entry := range entries {
					if namespace != "" && entry.Namespace != namespace {
						continue
					}
					if metricName != "" && entry.MetricName != metricName {
						continue
					}
					if len(dimensions) > 0 && !matchDimensions(convertFromChunkDimensions(entry.Dimensions), dimensions) {
						continue
					}

					metricKey := entry.Namespace + "#" + entry.MetricName + "#" + entry.BuildDimensionsKey()
					if !seen[metricKey] {
						seen[metricKey] = true
						*metrics = append(*metrics, MetricDatum{
							Namespace:  entry.Namespace,
							MetricName: entry.MetricName,
							Dimensions: convertFromChunkDimensions(entry.Dimensions),
						})
					}
				}
			}
		}
	}
}

func matchDimensions(dataDims, queryDims []Dimension) bool {
	for _, qd := range queryDims {
		found := false
		for _, dd := range dataDims {
			if dd.Name == qd.Name && (qd.Value == "" || dd.Value == qd.Value) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
