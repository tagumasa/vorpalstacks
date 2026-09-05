package s3

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/csv"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/parquet-go/parquet-go"
	"github.com/scritchley/orc"

	"vorpalstacks/internal/core/logs"
	s3store "vorpalstacks/internal/store/aws/s3"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// InventoryReportWorker periodically scans buckets carrying inventory
// configurations and publishes the reports to their destination buckets.
//
// Schedule semantics (production): a report is due at the next UTC midnight
// (Daily) or the next Sunday 00:00 UTC (Weekly) strictly after the
// configuration's schedule anchor — its last delivery, or the configuration's
// creation time before the first delivery. In TEST_MODE the cadence is
// compressed (45s Daily / 90s Weekly) so end-to-end tests can observe a
// delivery within one poll window; the compression is a test-runner-only
// mechanism.
type InventoryReportWorker struct {
	svc       *S3Service
	interval  time.Duration
	testMode  bool
	now       func() time.Time
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	startOnce sync.Once

	failMu   sync.Mutex
	failures map[string]*deliveryFailure
}

// deliveryFailure tracks the consecutive delivery failures of one
// configuration so repeated failures back off instead of rescanning the
// source bucket on every tick.
type deliveryFailure struct {
	count int
	last  time.Time
}

// NewInventoryReportWorker creates an InventoryReportWorker with a 5-minute
// default interval, tightened to 10 seconds in TEST_MODE.
func NewInventoryReportWorker(svc *S3Service) *InventoryReportWorker {
	testMode := os.Getenv("TEST_MODE") == "true"
	interval := 5 * time.Minute
	if testMode {
		interval = 10 * time.Second
	}
	return &InventoryReportWorker{
		svc:      svc,
		interval: interval,
		testMode: testMode,
		now:      time.Now,
		failures: map[string]*deliveryFailure{},
	}
}

// Start launches the inventory report goroutine.
func (w *InventoryReportWorker) Start() {
	w.startOnce.Do(func() {
		w.ctx, w.cancel = context.WithCancel(context.Background())
		w.wg.Add(1)
		go w.run()
		logs.Info("s3: inventory report worker started", logs.Any("interval", w.interval))
	})
}

// Close gracefully stops the inventory report worker.
func (w *InventoryReportWorker) Close() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
	logs.Info("s3: inventory report worker stopped")
}

func (w *InventoryReportWorker) run() {
	defer w.wg.Done()

	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-timer.C:
			w.deliverDueReports()
			timer.Reset(w.interval)
		}
	}
}

// deliverDueReports walks every region's buckets and produces the reports
// whose schedule boundary has passed. The scan goes through the same
// per-region stores the request plane uses; an ad-hoc store keyed by the
// account id would read a different, empty keyspace.
func (w *InventoryReportWorker) deliverDueReports() {
	seen := map[string]bool{}
	w.svc.s3RegionStores(func(region string, bucketStore *s3store.BucketStore, objectStore *s3store.ObjectStore) {
		buckets, err := bucketStore.List()
		if err != nil {
			logs.Warn("s3: inventory worker failed to list buckets", logs.String("region", region), logs.Err(err))
			return
		}
		for _, bucket := range buckets {
			for _, config := range bucket.InventoryConfigurations {
				if !inventoryReportDue(config, w.now(), w.testMode) {
					continue
				}
				key := w.failureKey(region, bucket.Name, config.ID)
				seen[key] = true
				if w.backedOff(key, w.now()) {
					continue
				}
				if err := w.deliverReport(bucketStore, objectStore, bucket.Name, config); err != nil {
					w.noteFailure(key, w.now())
					logs.Warn("s3: inventory report delivery failed",
						logs.String("bucket", bucket.Name),
						logs.String("configId", config.ID),
						logs.Err(err))
					continue
				}
				w.noteSuccess(key)
			}
		}
	})
	w.pruneFailures(seen)
}

// failureKey identifies one configuration's failure state across regions.
func (w *InventoryReportWorker) failureKey(region, bucket, id string) string {
	return region + "/" + bucket + "/" + id
}

// backedOff reports whether a configuration with consecutive delivery
// failures is still within its backoff window, which doubles per failure
// up to eight poll intervals; a configuration without failures is never
// backed off.
func (w *InventoryReportWorker) backedOff(key string, now time.Time) bool {
	w.failMu.Lock()
	defer w.failMu.Unlock()
	f := w.failures[key]
	if f == nil {
		return false
	}
	window := w.interval << (f.count - 1)
	if window > 8*w.interval {
		window = 8 * w.interval
	}
	return now.Sub(f.last) < window
}

func (w *InventoryReportWorker) noteFailure(key string, now time.Time) {
	w.failMu.Lock()
	defer w.failMu.Unlock()
	f := w.failures[key]
	if f == nil {
		f = &deliveryFailure{}
		w.failures[key] = f
	}
	f.count++
	f.last = now
}

func (w *InventoryReportWorker) noteSuccess(key string) {
	w.failMu.Lock()
	defer w.failMu.Unlock()
	delete(w.failures, key)
}

// pruneFailures drops the failure state of configurations that no longer
// exist.
func (w *InventoryReportWorker) pruneFailures(seen map[string]bool) {
	w.failMu.Lock()
	defer w.failMu.Unlock()
	for key := range w.failures {
		if !seen[key] {
			delete(w.failures, key)
		}
	}
}

// deliverReport produces one report for one configuration and, on success,
// advances the configuration's schedule anchor. The destination is resolved
// before the source bucket is scanned so an unusable destination never pays
// for the scan.
func (w *InventoryReportWorker) deliverReport(bucketStore *s3store.BucketStore, sourceObjects *s3store.ObjectStore, sourceBucket string, config *s3store.InventoryConfiguration) error {
	scanStart := w.now().UTC()

	dest := config.Destination.S3BucketDestination
	destBucket, err := bucketNameFromARN(dest.Bucket)
	if err != nil {
		return err
	}
	// The destination bucket may live in a different region than the source;
	// resolve its region through the same cross-region lookup the request
	// plane uses.
	_, destRegion := w.svc.s3FindBucket(destBucket)
	destObjects := w.svc.s3Objects(destRegion)
	if destObjects == nil {
		return fmt.Errorf("inventory destination bucket %q does not exist", destBucket)
	}
	destPrefix := strings.Trim(dest.Prefix, "/")

	// The configuration's InventoryEncryption choice applies to every object
	// of the delivery; an unset choice leaves the objects unencrypted. When
	// both members are present the more specific SSE-KMS choice wins.
	var encType EncryptionType
	var kmsKeyID string
	if dest.Encryption != nil {
		if dest.Encryption.SSES3 {
			encType = EncryptionTypeSSE_S3
		}
		if dest.Encryption.SSEKMS != nil {
			encType = EncryptionTypeSSE_KMS
			kmsKeyID = dest.Encryption.SSEKMS.KeyID
		}
	}

	rows, err := collectInventoryRows(sourceObjects, sourceBucket, config)
	if err != nil {
		return err
	}

	var data []byte
	var extension, contentType, fileSchema string
	switch dest.Format {
	case "CSV":
		data, fileSchema, err = buildCSVReport(sourceBucket, config, rows)
		extension, contentType = "csv.gz", "application/gzip"
	case "Parquet":
		data, fileSchema, err = buildParquetReport(sourceBucket, config, rows)
		extension, contentType = "parquet", "application/octet-stream"
	case "ORC":
		data, fileSchema, err = buildORCReport(sourceBucket, config, rows)
		extension, contentType = "orc", "application/octet-stream"
	default:
		return fmt.Errorf("unsupported inventory report format %q", dest.Format)
	}
	if err != nil {
		return err
	}

	// Every delivery key derives from the schedule window's anchor stamp, so
	// a partial failure rewrites the same keys on retry instead of piling up
	// orphan data files; the manifest is the commit point and the anchor
	// advances only after every object of the delivery landed. AWS documents
	// unpredictable data file names — a window-stable name is this
	// platform's delivery-key choice.
	windowStamp := config.LastDelivery.UTC().Format("2006-01-02T15-04-05Z")
	dataKey := inventoryDataKey(destPrefix, sourceBucket, config.ID, windowStamp, extension)
	manifestKey := inventoryManifestKey(destPrefix, sourceBucket, config.ID, windowStamp)
	checksumKey := manifestKey + ".checksum"
	symlinkKey := inventorySymlinkKey(destPrefix, sourceBucket, config.ID, windowStamp)

	if err := w.putReportObject(destObjects, destBucket, dataKey, bytes.NewReader(data), contentType, encType, kmsKeyID); err != nil {
		return err
	}

	manifest := buildInventoryManifest(sourceBucket, dest.Bucket, scanStart, dest.Format, fileSchema,
		[]inventoryManifestFile{{
			Key:         dataKey,
			Size:        int64(len(data)),
			MD5Checksum: md5Hex(data),
		}})
	manifestJSON, err := json.MarshalIndent(manifest, "", "    ")
	if err != nil {
		return err
	}
	if err := w.putReportObject(destObjects, destBucket, manifestKey, bytes.NewReader(manifestJSON), "application/json", encType, kmsKeyID); err != nil {
		return err
	}
	if err := w.putReportObject(destObjects, destBucket, checksumKey,
		strings.NewReader(md5Hex(manifestJSON)), "text/plain", encType, kmsKeyID); err != nil {
		return err
	}
	symlink := fmt.Sprintf("s3://%s/%s\n", destBucket, dataKey)
	if err := w.putReportObject(destObjects, destBucket, symlinkKey, strings.NewReader(symlink), "text/plain", encType, kmsKeyID); err != nil {
		return err
	}

	// The anchor advances only after every object of the delivery landed, so
	// a failed delivery is retried on the next tick instead of being lost.
	config.LastDelivery = scanStart
	return bucketStore.SetInventoryConfiguration(sourceBucket, config.ID, config)
}

// putReportObject writes one delivered report object. When the
// configuration carries an InventoryEncryption choice, the content is
// encrypted with the same chunked AES-GCM scheme the object-put path uses
// and stored together with its SSE metadata.
func (w *InventoryReportWorker) putReportObject(objects s3store.ObjectStoreInterface, bucket, key string, content io.Reader, contentType string, encType EncryptionType, kmsKeyID string) error {
	if encType == "" {
		_, err := objects.Put(w.ctx, bucket, key, content, contentType, nil)
		return err
	}
	res, err := w.svc.encryptionManager.EncryptStream(content, encType, nil, bucket, key, kmsKeyID, nil)
	if err != nil {
		return err
	}
	_, err = objects.PutEncrypted(w.ctx, bucket, key, res.EncryptedData, contentType, nil, res.SSEMetadata, s3store.StorageClassStandard, nil)
	return err
}

// collectInventoryRows gathers the object metadata rows the report lists:
// every version (delete markers included) for IncludedObjectVersions=All, or
// the current version of each key for Current.
func collectInventoryRows(objectStore *s3store.ObjectStore, sourceBucket string, config *s3store.InventoryConfiguration) ([]*s3store.Object, error) {
	prefix := ""
	if config.Filter != nil {
		prefix = config.Filter.Prefix
	}

	var rows []*s3store.Object
	if config.IncludedObjectVersions == "All" {
		keyMarker, versionMarker := "", ""
		for {
			result, err := objectStore.ListObjectVersions(sourceBucket, prefix, "", keyMarker, versionMarker, 1000)
			if err != nil {
				return nil, err
			}
			rows = append(rows, result.Objects...)
			if !result.IsTruncated {
				break
			}
			keyMarker, versionMarker = result.NextVersionKeyMarker, result.NextVersionIDMarker
		}
		return rows, nil
	}

	marker := ""
	for {
		result, err := objectStore.List(sourceBucket, prefix, "", marker, 1000)
		if err != nil {
			return nil, err
		}
		for _, obj := range result.Objects {
			if obj.IsDeleteMarker {
				continue
			}
			rows = append(rows, obj)
		}
		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}
	return rows, nil
}

// inventoryReportDue reports whether a configuration's schedule boundary has
// passed. The anchor is the last delivery (or creation time before the
// first one); in TEST_MODE the cadence compresses to fixed periods.
func inventoryReportDue(config *s3store.InventoryConfiguration, now time.Time, testMode bool) bool {
	if config == nil || !config.IsEnabled || config.Schedule == nil {
		return false
	}
	anchor := config.LastDelivery
	if anchor.IsZero() {
		return false
	}
	if testMode {
		if config.Schedule.Frequency == "Weekly" {
			return now.Sub(anchor) >= 90*time.Second
		}
		return now.Sub(anchor) >= 45*time.Second
	}
	return now.UTC().After(nextScheduleBoundary(anchor, config.Schedule.Frequency))
}

// nextScheduleBoundary returns the first Daily (UTC midnight) or Weekly
// (Sunday 00:00 UTC) boundary strictly after the anchor.
func nextScheduleBoundary(anchor time.Time, frequency string) time.Time {
	utc := anchor.UTC()
	day := time.Date(utc.Year(), utc.Month(), utc.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, 1)
	if frequency == "Weekly" {
		for weekday := day.Weekday(); weekday != time.Sunday; weekday = day.Weekday() {
			day = day.AddDate(0, 0, 1)
		}
	}
	return day
}

// inventoryColumnOrder is the canonical report column order: bucket and key
// always, the version columns when IncludedObjectVersions=All, then the
// optional fields in the order AWS lists them.
var inventoryColumnOrder = []string{
	"Size", "LastModifiedDate", "ETag", "StorageClass", "IsMultipartUploaded",
	"ReplicationStatus", "EncryptionStatus", "ObjectLockRetainUntilDate",
	"ObjectLockMode", "ObjectLockLegalHoldStatus", "IntelligentTieringAccessTier",
	"BucketKeyStatus", "ChecksumAlgorithm", "ObjectAccessControlList",
	"ObjectOwner", "LifecycleExpirationDate",
}

// inventoryColumns returns the report's column names in canonical order.
func inventoryColumns(config *s3store.InventoryConfiguration) []string {
	columns := []string{"Bucket", "Key"}
	if config.IncludedObjectVersions == "All" {
		columns = append(columns, "VersionId", "IsLatest", "IsDeleteMarker")
	}
	selected := make(map[string]bool, len(config.OptionalFields))
	for _, field := range config.OptionalFields {
		selected[field] = true
	}
	for _, field := range inventoryColumnOrder {
		if selected[field] {
			columns = append(columns, field)
		}
	}
	return columns
}

// reportColumnValue renders one object's value for one column. CSV booleans
// are the uppercase words AWS's examples show, empty values stay empty, and
// the CSV key column is percent-encoded (slashes literal) as AWS documents.
func reportColumnValue(column string, bucket string, obj *s3store.Object) string {
	switch column {
	case "Bucket":
		return bucket
	case "Key":
		return obj.Key
	case "VersionId":
		if obj.VersionID == "" {
			// Inventory reports carry the literal string "null" for objects
			// without a version id.
			return "null"
		}
		return obj.VersionID
	case "IsLatest":
		return boolWord(obj.IsLatest)
	case "IsDeleteMarker":
		return boolWord(obj.IsDeleteMarker)
	case "Size":
		return fmt.Sprintf("%d", obj.Size)
	case "LastModifiedDate":
		return obj.LastModified.UTC().Format("2006-01-02T15:04:05.000Z")
	case "ETag":
		return strings.Trim(obj.ETag, `"`)
	case "StorageClass":
		if obj.StorageClass == "" {
			return "STANDARD"
		}
		return string(obj.StorageClass)
	case "IsMultipartUploaded":
		return boolWord(len(obj.Parts) > 0)
	case "ReplicationStatus":
		return obj.ReplicationStatus
	case "EncryptionStatus":
		return objectEncryptionStatus(obj)
	case "ObjectLockRetainUntilDate":
		if obj.ObjectLockRetention == nil || obj.ObjectLockRetention.RetainUntilDate == nil {
			return ""
		}
		return obj.ObjectLockRetention.RetainUntilDate.UTC().Format("2006-01-02T15:04:05.000Z")
	case "ObjectLockMode":
		if obj.ObjectLockRetention == nil {
			return ""
		}
		return string(obj.ObjectLockRetention.Mode)
	case "ObjectLockLegalHoldStatus":
		if obj.ObjectLockLegalHold == nil {
			return ""
		}
		return string(obj.ObjectLockLegalHold.Status)
	case "IntelligentTieringAccessTier", "ChecksumAlgorithm", "LifecycleExpirationDate":
		// No substrate on this single-tier platform: the columns exist but
		// carry no values.
		return ""
	case "BucketKeyStatus":
		// The column reports S3 Bucket Key usage for SSE-KMS objects. The
		// platform has no bucket keys, so every SSE-KMS object is DISABLED,
		// and objects without SSE-KMS carry no bucket-key status.
		switch obj.ServerSideEncryption {
		case "aws:kms", "aws:kms+dbz":
			return "DISABLED"
		}
		return ""
	case "ObjectAccessControlList":
		if obj.ACL == nil {
			return ""
		}
		return base64.StdEncoding.EncodeToString([]byte(objectACLReportJSON(obj)))
	case "ObjectOwner":
		if obj.Owner == nil {
			return ""
		}
		return obj.Owner.ID
	}
	return ""
}

func boolWord(value bool) string {
	if value {
		return "TRUE"
	}
	return "FALSE"
}

func objectEncryptionStatus(obj *s3store.Object) string {
	switch obj.ServerSideEncryption {
	case "AES256":
		return "SSE-S3"
	case "aws:kms", "aws:kms+dbz":
		return "SSE-KMS"
	case "":
		if obj.SSEMetadata != nil && obj.SSEMetadata.EncryptionType == s3store.SSETypeCustomer {
			return "SSE-C"
		}
		return "NOT-SSE"
	default:
		return "NOT-SSE"
	}
}

// objectACLReportJSON renders the object ACL in the JSON form AWS embeds
// (base64-encoded) in the ObjectAccessControlList column.
func objectACLReportJSON(obj *s3store.Object) string {
	type reportGrant struct {
		CanonicalID string `json:"canonicalId"`
		Type        string `json:"type"`
		Permission  string `json:"permission"`
	}
	grants := make([]reportGrant, 0, len(obj.ACL.Grants))
	for _, grant := range obj.ACL.Grants {
		if grant.Grantee == nil {
			continue
		}
		grants = append(grants, reportGrant{
			CanonicalID: grant.Grantee.ID,
			Type:        string(grant.Grantee.Type),
			Permission:  string(grant.Permission),
		})
	}
	payload := struct {
		Version string        `json:"version"`
		Status  string        `json:"status"`
		Grants  []reportGrant `json:"grants"`
	}{Version: "2022-11-10", Status: "AVAILABLE", Grants: grants}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return ""
	}
	return string(encoded)
}

// csvEncodeKey percent-encodes a key for the CSV report column: AWS
// URL-encodes key names in CSV inventory files; slashes stay literal so the
// keys remain path-addressable.
func csvEncodeKey(key string) string {
	var b strings.Builder
	for _, ch := range key {
		if ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9' ||
			ch == '-' || ch == '_' || ch == '.' || ch == '~' || ch == '/' {
			b.WriteRune(ch)
			continue
		}
		for _, byteValue := range []byte(string(ch)) {
			fmt.Fprintf(&b, "%%%02X", byteValue)
		}
	}
	return b.String()
}

// buildCSVReport renders the gzip-compressed CSV report and its schema
// string (the comma-separated column names).
func buildCSVReport(bucket string, config *s3store.InventoryConfiguration, rows []*s3store.Object) ([]byte, string, error) {
	columns := inventoryColumns(config)

	var raw bytes.Buffer
	writer := csv.NewWriter(&raw)
	for _, obj := range rows {
		record := make([]string, len(columns))
		for i, column := range columns {
			value := reportColumnValue(column, bucket, obj)
			if column == "Key" {
				value = csvEncodeKey(value)
			}
			record[i] = value
		}
		if err := writer.Write(record); err != nil {
			return nil, "", err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", err
	}

	var out bytes.Buffer
	gzipWriter := gzip.NewWriter(&out)
	if _, err := gzipWriter.Write(raw.Bytes()); err != nil {
		return nil, "", err
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, "", err
	}
	return out.Bytes(), strings.Join(columns, ", "), nil
}

// parquetColumnType maps one report column to its Parquet field: the Go
// value kind drives the physical type, names follow the AWS schema.
func buildParquetReport(bucket string, config *s3store.InventoryConfiguration, rows []*s3store.Object) ([]byte, string, error) {
	columns := inventoryColumns(config)

	type columnSpec struct {
		name     string
		optional bool
		kind     parquet.Kind
	}
	specs := make([]columnSpec, 0, len(columns))
	for _, column := range columns {
		spec := columnSpec{name: parquetColumnName(column)}
		switch column {
		case "Bucket", "Key":
			spec.kind = parquet.ByteArray
		case "IsLatest", "IsDeleteMarker", "IsMultipartUploaded":
			spec.kind = parquet.Boolean
			spec.optional = true
		case "Size", "LastModifiedDate", "ObjectLockRetainUntilDate":
			spec.kind = parquet.Int64
			spec.optional = true
		default:
			spec.kind = parquet.ByteArray
			spec.optional = true
		}
		specs = append(specs, spec)
	}

	fields := make([]reflect.StructField, len(specs))
	for i, spec := range specs {
		tag := reflect.StructTag(fmt.Sprintf(`parquet:"%s"`, spec.name))
		if spec.optional {
			tag = reflect.StructTag(fmt.Sprintf(`parquet:"%s,optional"`, spec.name))
		}
		fields[i] = reflect.StructField{
			Name: fmt.Sprintf("Field%d", i),
			Type: reflect.TypeOf(""),
			Tag:  tag,
		}
		// Boolean and int64 columns need matching Go kinds.
		switch spec.kind {
		case parquet.Boolean:
			fields[i].Type = reflect.TypeOf(false)
		case parquet.Int64:
			fields[i].Type = reflect.TypeOf(int64(0))
		}
	}
	rowType := reflect.StructOf(fields)

	values := make([]any, 0, len(rows))
	for _, obj := range rows {
		row := reflect.New(rowType).Elem()
		for i, column := range columns {
			switch fields[i].Type.Kind() {
			case reflect.Bool:
				row.Field(i).SetBool(reportColumnBool(column, obj))
			case reflect.Int64:
				row.Field(i).SetInt(reportColumnInt64(column, obj))
			default:
				row.Field(i).SetString(reportColumnValue(column, bucket, obj))
			}
		}
		values = append(values, row.Interface())
	}

	var out bytes.Buffer
	writer := parquet.NewWriter(&out, parquet.Compression(&parquet.Snappy))
	for _, value := range values {
		if err := writer.Write(value); err != nil {
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return out.Bytes(), parquetFileSchema(config), nil
}

// orcColumnCategory maps one report column to its ORC type per the AWS
// inventory schema.
func orcColumnCategory(column string) string {
	switch column {
	case "IsLatest", "IsDeleteMarker", "IsMultipartUploaded":
		return "boolean"
	case "Size":
		return "bigint"
	case "LastModifiedDate", "ObjectLockRetainUntilDate":
		return "timestamp"
	default:
		return "string"
	}
}

// orcFileSchema renders the manifest's fileSchema string for ORC reports in
// the struct form AWS documents.
func orcFileSchema(config *s3store.InventoryConfiguration) string {
	var b strings.Builder
	b.WriteString("struct<bucket:string,key:string")
	if config.IncludedObjectVersions == "All" {
		b.WriteString(",version_id:string,is_latest:boolean,is_delete_marker:boolean")
	}
	for _, field := range inventoryColumnOrder {
		if !optionalFieldSelected(config, field) {
			continue
		}
		b.WriteString("," + parquetColumnName(field) + ":" + orcColumnCategory(field))
	}
	b.WriteString(">")
	return b.String()
}

// buildORCReport renders the ORC report and its schema string, compressed
// with ZLIB as the AWS ORC deliveries are.
func buildORCReport(bucket string, config *s3store.InventoryConfiguration, rows []*s3store.Object) ([]byte, string, error) {
	columns := inventoryColumns(config)
	schema, err := orc.ParseSchema(orcFileSchema(config))
	if err != nil {
		return nil, "", fmt.Errorf("orc schema: %w", err)
	}

	var out bytes.Buffer
	writer, err := orc.NewWriter(&out, orc.SetSchema(schema), orc.SetCompression(orc.CompressionZlib{}))
	if err != nil {
		return nil, "", err
	}
	for _, obj := range rows {
		values := make([]interface{}, 0, len(columns))
		for _, column := range columns {
			switch orcColumnCategory(column) {
			case "boolean":
				values = append(values, reportColumnBool(column, obj))
			case "bigint":
				values = append(values, reportColumnInt64(column, obj))
			case "timestamp":
				if column == "ObjectLockRetainUntilDate" && obj.ObjectLockRetention != nil &&
					obj.ObjectLockRetention.RetainUntilDate != nil {
					values = append(values, *obj.ObjectLockRetention.RetainUntilDate)
					continue
				}
				values = append(values, obj.LastModified)
			default:
				values = append(values, reportColumnValue(column, bucket, obj))
			}
		}
		if err := writer.Write(values...); err != nil {
			return nil, "", err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return out.Bytes(), orcFileSchema(config), nil
}

// reportColumnInt64 renders the int64 columns: sizes in bytes, timestamps
// as epoch milliseconds (the AWS Parquet schema's TIMESTAMP_MILLIS).
func reportColumnInt64(column string, obj *s3store.Object) int64 {
	switch column {
	case "Size":
		return obj.Size
	case "LastModifiedDate":
		return obj.LastModified.UnixMilli()
	case "ObjectLockRetainUntilDate":
		if obj.ObjectLockRetention == nil || obj.ObjectLockRetention.RetainUntilDate == nil {
			return 0
		}
		return obj.ObjectLockRetention.RetainUntilDate.UnixMilli()
	}
	return 0
}

func reportColumnBool(column string, obj *s3store.Object) bool {
	switch column {
	case "IsLatest":
		return obj.IsLatest
	case "IsDeleteMarker":
		return obj.IsDeleteMarker
	case "IsMultipartUploaded":
		return len(obj.Parts) > 0
	}
	return false
}

// parquetColumnName maps the AWS column names to the snake_case names of
// the AWS Parquet inventory schema.
func parquetColumnName(column string) string {
	switch column {
	case "VersionId":
		return "version_id"
	case "StorageClass":
		return "storage_class"
	case "IsLatest":
		return "is_latest"
	case "IsDeleteMarker":
		return "is_delete_marker"
	case "LastModifiedDate":
		return "last_modified_date"
	case "ETag":
		return "e_tag"
	case "IsMultipartUploaded":
		return "is_multipart_uploaded"
	case "ReplicationStatus":
		return "replication_status"
	case "EncryptionStatus":
		return "encryption_status"
	case "ObjectLockRetainUntilDate":
		return "object_lock_retain_until_date"
	case "ObjectLockMode":
		return "object_lock_mode"
	case "ObjectLockLegalHoldStatus":
		return "object_lock_legal_hold_status"
	case "IntelligentTieringAccessTier":
		return "intelligent_tiering_access_tier"
	case "BucketKeyStatus":
		return "bucket_key_status"
	case "ChecksumAlgorithm":
		return "checksum_algorithm"
	case "ObjectAccessControlList":
		return "object_access_control_list"
	case "ObjectOwner":
		return "object_owner"
	case "LifecycleExpirationDate":
		return "lifecycle_expiration_date"
	default:
		return strings.ToLower(column)
	}
}

// parquetFileSchema renders the manifest's fileSchema string for Parquet
// reports in the AWS schema form.
func parquetFileSchema(config *s3store.InventoryConfiguration) string {
	var b strings.Builder
	b.WriteString("message s3.inventory { ")
	b.WriteString("required binary bucket (UTF8); required binary key (UTF8);")
	if config.IncludedObjectVersions == "All" {
		b.WriteString(" optional binary version_id (UTF8); optional boolean is_latest; optional boolean is_delete_marker;")
	}
	for _, field := range inventoryColumnOrder {
		if !optionalFieldSelected(config, field) {
			continue
		}
		name := parquetColumnName(field)
		switch field {
		case "Size":
			b.WriteString(" optional int64 size;")
		case "IsMultipartUploaded":
			b.WriteString(" optional boolean is_multipart_uploaded;")
		case "LastModifiedDate", "ObjectLockRetainUntilDate":
			b.WriteString(" optional int64 " + name + " (TIMESTAMP_MILLIS);")
		default:
			b.WriteString(" optional binary " + name + " (UTF8);")
		}
	}
	b.WriteString(" }")
	return b.String()
}

func optionalFieldSelected(config *s3store.InventoryConfiguration, field string) bool {
	for _, candidate := range config.OptionalFields {
		if candidate == field {
			return true
		}
	}
	return false
}

// inventoryManifestFile describes one data file of a delivery.
type inventoryManifestFile struct {
	Key         string `json:"key"`
	Size        int64  `json:"size"`
	MD5Checksum string `json:"MD5checksum"`
}

// inventoryManifest is the manifest.json AWS publishes with every delivery.
type inventoryManifest struct {
	SourceBucket      string                  `json:"sourceBucket"`
	DestinationBucket string                  `json:"destinationBucket"`
	Version           string                  `json:"version"`
	CreationTimestamp string                  `json:"creationTimestamp"`
	FileFormat        string                  `json:"fileFormat"`
	FileSchema        string                  `json:"fileSchema"`
	Files             []inventoryManifestFile `json:"files"`
}

func buildInventoryManifest(sourceBucket, destinationBucket string, scanStart time.Time, format, fileSchema string, files []inventoryManifestFile) inventoryManifest {
	return inventoryManifest{
		SourceBucket:      sourceBucket,
		DestinationBucket: destinationBucket,
		Version:           "2016-11-30",
		CreationTimestamp: fmt.Sprintf("%d", scanStart.UnixMilli()),
		FileFormat:        format,
		FileSchema:        fileSchema,
		Files:             files,
	}
}

func md5Hex(data []byte) string {
	sum := md5.Sum(data)
	return hex.EncodeToString(sum[:])
}

// bucketNameFromARN extracts the bucket name from the destination ARN the
// configuration carries (arn:aws:s3:::bucket-name).
func bucketNameFromARN(arn string) (string, error) {
	parsed, err := svcarn.ParseARN(arn)
	if err != nil {
		return "", fmt.Errorf("parse inventory destination %q: %w", arn, err)
	}
	return parsed.Resource, nil
}

// inventoryDataKey builds the data file's delivery key. The file name is
// stable per schedule window (delivery-<window stamp>) so retries of a
// partially failed delivery converge on the same object.
func inventoryDataKey(destPrefix, sourceBucket, configID, windowStamp, extension string) string {
	name := fmt.Sprintf("delivery-%s.%s", windowStamp, extension)
	parts := []string{sourceBucket, configID, "data", name}
	if destPrefix != "" {
		return destPrefix + "/" + strings.Join(parts, "/")
	}
	return strings.Join(parts, "/")
}

// inventoryManifestKey builds the manifest's delivery key under the window's
// stamp folder.
func inventoryManifestKey(destPrefix, sourceBucket, configID, windowStamp string) string {
	parts := []string{sourceBucket, configID, windowStamp, "manifest.json"}
	if destPrefix != "" {
		return destPrefix + "/" + strings.Join(parts, "/")
	}
	return strings.Join(parts, "/")
}

// inventorySymlinkKey builds the Hive-style symlink's delivery key under the
// window's dt= partition folder.
func inventorySymlinkKey(destPrefix, sourceBucket, configID, windowStamp string) string {
	parts := []string{sourceBucket, configID, "hive", "dt=" + windowStamp, "symlink.txt"}
	if destPrefix != "" {
		return destPrefix + "/" + strings.Join(parts, "/")
	}
	return strings.Join(parts, "/")
}
