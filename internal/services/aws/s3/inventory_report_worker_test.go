package s3

import (
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/parquet-go/parquet-go"
	"github.com/scritchley/orc"

	s3store "vorpalstacks/internal/store/aws/s3"
)

func TestInventoryReportDueFollowsScheduleBoundaries(t *testing.T) {
	base := time.Date(2026, 9, 3, 10, 30, 0, 0, time.UTC) // Thursday

	daily := &s3store.InventoryConfiguration{
		IsEnabled:    true,
		Schedule:     &s3store.InventorySchedule{Frequency: "Daily"},
		LastDelivery: base,
	}
	if inventoryReportDue(daily, base.Add(time.Hour), false) {
		t.Fatal("daily report must not be due before the next UTC midnight")
	}
	if !inventoryReportDue(daily, time.Date(2026, 9, 4, 0, 0, 1, 0, time.UTC), false) {
		t.Fatal("daily report must be due after the next UTC midnight")
	}

	weekly := &s3store.InventoryConfiguration{
		IsEnabled:    true,
		Schedule:     &s3store.InventorySchedule{Frequency: "Weekly"},
		LastDelivery: base,
	}
	if inventoryReportDue(weekly, time.Date(2026, 9, 5, 0, 0, 1, 0, time.UTC), false) {
		t.Fatal("weekly report must not be due before Sunday UTC")
	}
	if !inventoryReportDue(weekly, time.Date(2026, 9, 6, 0, 0, 1, 0, time.UTC), false) {
		t.Fatal("weekly report must be due after Sunday 00:00 UTC")
	}

	disabled := *daily
	disabled.IsEnabled = false
	if inventoryReportDue(&disabled, base.AddDate(0, 0, 3), false) {
		t.Fatal("a disabled configuration never comes due")
	}
}

func TestInventoryReportDueTestModeCompression(t *testing.T) {
	base := time.Now()
	daily := &s3store.InventoryConfiguration{
		IsEnabled:    true,
		Schedule:     &s3store.InventorySchedule{Frequency: "Daily"},
		LastDelivery: base,
	}
	if inventoryReportDue(daily, base.Add(30*time.Second), true) {
		t.Fatal("compressed daily cadence must not be due before its period")
	}
	if !inventoryReportDue(daily, base.Add(46*time.Second), true) {
		t.Fatal("compressed daily cadence must be due after 45s")
	}
}

func TestBuildCSVReportColumnsAndEncoding(t *testing.T) {
	config := &s3store.InventoryConfiguration{
		IncludedObjectVersions: "All",
		OptionalFields:         []string{"Size", "LastModifiedDate", "ETag", "StorageClass"},
	}
	modified := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	rows := []*s3store.Object{
		{Key: "plain/key.txt", VersionID: "v2", IsLatest: true, Size: 1500, LastModified: modified, ETag: `"abc123"`, StorageClass: "STANDARD"},
		{Key: "spaced key +&.csv", IsDeleteMarker: true},
	}

	data, schema, err := buildCSVReport("src-bucket", config, rows)
	if err != nil {
		t.Fatalf("buildCSVReport: %v", err)
	}
	wantSchema := "Bucket, Key, VersionId, IsLatest, IsDeleteMarker, Size, LastModifiedDate, ETag, StorageClass"
	if schema != wantSchema {
		t.Fatalf("schema = %q, want %q", schema, wantSchema)
	}

	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	records, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		t.Fatalf("csv read: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("records = %d, want 2", len(records))
	}
	wantFirst := []string{"src-bucket", "plain/key.txt", "v2", "TRUE", "FALSE", "1500", "2026-09-01T12:00:00.000Z", "abc123", "STANDARD"}
	if strings.Join(records[0], "|") != strings.Join(wantFirst, "|") {
		t.Fatalf("first record = %v, want %v", records[0], wantFirst)
	}
	second := records[1]
	if second[1] != "spaced%20key%20%2B%26.csv" {
		t.Fatalf("key column = %q, want percent-encoded form", second[1])
	}
	if second[2] != "null" {
		t.Fatalf("empty version id = %q, want the literal string null", second[2])
	}
	if second[5] != "0" || second[8] != "STANDARD" {
		t.Fatalf("delete-marker row = %v, want zero size and STANDARD class", second)
	}
}

func TestBuildCSVReportCurrentOmitsVersionColumns(t *testing.T) {
	config := &s3store.InventoryConfiguration{
		IncludedObjectVersions: "Current",
		OptionalFields:         []string{"Size"},
	}
	data, schema, err := buildCSVReport("b", config, nil)
	if err != nil {
		t.Fatalf("buildCSVReport: %v", err)
	}
	if schema != "Bucket, Key, Size" {
		t.Fatalf("schema = %q, want Bucket, Key, Size", schema)
	}
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip reader: %v", err)
	}
	body, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(body) != 0 {
		t.Fatalf("empty report body = %q, want empty", body)
	}
}

func TestBuildParquetReportRoundTrip(t *testing.T) {
	config := &s3store.InventoryConfiguration{
		IncludedObjectVersions: "All",
		OptionalFields:         []string{"Size", "LastModifiedDate", "StorageClass", "IsMultipartUploaded"},
	}
	modified := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	rows := []*s3store.Object{
		{Key: "data/one.bin", VersionID: "v1", IsLatest: true, Size: 42, LastModified: modified, StorageClass: "STANDARD"},
		{Key: "data/two.bin", VersionID: "v2", Size: 7, LastModified: modified, Parts: []s3store.ObjectPartBoundary{{PartNumber: 1}}},
	}

	data, schema, err := buildParquetReport("src-bucket", config, rows)
	if err != nil {
		t.Fatalf("buildParquetReport: %v", err)
	}
	if !strings.Contains(schema, "required binary bucket (UTF8)") ||
		!strings.Contains(schema, "optional int64 last_modified_date (TIMESTAMP_MILLIS)") ||
		!strings.Contains(schema, "optional boolean is_multipart_uploaded") {
		t.Fatalf("schema = %q, want the AWS column forms", schema)
	}

	type row struct {
		Bucket              string  `parquet:"bucket"`
		Key                 string  `parquet:"key"`
		VersionID           *string `parquet:"version_id,optional"`
		IsLatest            *bool   `parquet:"is_latest,optional"`
		IsDeleteMarker      *bool   `parquet:"is_delete_marker,optional"`
		Size                *int64  `parquet:"size,optional"`
		LastModifiedDate    *int64  `parquet:"last_modified_date,optional"`
		StorageClass        *string `parquet:"storage_class,optional"`
		IsMultipartUploaded *bool   `parquet:"is_multipart_uploaded,optional"`
	}
	rowsOut, err := parquet.Read[row](bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("parquet read: %v", err)
	}
	got := rowsOut
	if len(got) != 2 {
		t.Fatalf("rows = %d, want 2", len(got))
	}
	if got[0].Bucket != "src-bucket" || got[0].Key != "data/one.bin" || got[0].Size == nil || *got[0].Size != 42 {
		t.Fatalf("first row = %+v", got[0])
	}
	if got[0].LastModifiedDate == nil || *got[0].LastModifiedDate != modified.UnixMilli() {
		t.Fatalf("first row timestamp = %v, want %d", got[0].LastModifiedDate, modified.UnixMilli())
	}
	if got[1].IsMultipartUploaded == nil || !*got[1].IsMultipartUploaded {
		t.Fatalf("second row multipart flag = %v, want true", got[1].IsMultipartUploaded)
	}
}

// TestBuildORCReportReadBack pins the ORC report format: the manifest
// schema matches the struct form AWS documents, and the file decodes back
// through the ORC reader with the documented column types and values.
func TestBuildORCReportReadBack(t *testing.T) {
	config := &s3store.InventoryConfiguration{
		IncludedObjectVersions: "All",
		OptionalFields:         []string{"Size", "LastModifiedDate", "BucketKeyStatus"},
	}
	modified := time.Date(2026, 9, 1, 12, 0, 0, 0, time.UTC)
	rows := []*s3store.Object{
		{Key: "kept/kms.txt", VersionID: "v1", IsLatest: true, Size: 42, LastModified: modified, ServerSideEncryption: "aws:kms"},
		{Key: "kept/plain.txt", VersionID: "v2", IsLatest: false, Size: 7, LastModified: modified},
	}

	data, schema, err := buildORCReport("src-bucket", config, rows)
	if err != nil {
		t.Fatalf("buildORCReport: %v", err)
	}
	wantSchema := "struct<bucket:string,key:string,version_id:string,is_latest:boolean,is_delete_marker:boolean,size:bigint,last_modified_date:timestamp,bucket_key_status:string>"
	if schema != wantSchema {
		t.Fatalf("schema = %q, want %q", schema, wantSchema)
	}

	reader, err := orc.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("orc reader: %v", err)
	}
	cursor := reader.Select("bucket", "key", "size", "bucket_key_status")
	seen := map[string]string{}
	for cursor.Stripes() {
		for cursor.Next() {
			row := cursor.Row()
			seen[row[1].(string)] = fmt.Sprintf("%d|%s", row[2].(int64), row[3].(string))
		}
	}
	if err := cursor.Err(); err != nil {
		t.Fatalf("cursor: %v", err)
	}
	if len(seen) != 2 {
		t.Fatalf("read back %d rows, want 2", len(seen))
	}
	if got := seen["kept/kms.txt"]; got != "42|DISABLED" {
		t.Fatalf("kms row = %q, want 42|DISABLED", got)
	}
	if got := seen["kept/plain.txt"]; got != "7|" {
		t.Fatalf("plain row = %q, want 7| (empty bucket key status)", got)
	}
}

func TestBuildInventoryManifest(t *testing.T) {
	scanStart := time.Date(2026, 9, 5, 8, 30, 0, 0, time.UTC)
	files := []inventoryManifestFile{{Key: "src/id/data/x.csv.gz", Size: 12, MD5Checksum: "abc"}}
	manifest := buildInventoryManifest("src", "arn:aws:s3:::dest", scanStart, "CSV", "Bucket, Key", files)
	if manifest.Version != "2016-11-30" {
		t.Fatalf("version = %q", manifest.Version)
	}
	if manifest.CreationTimestamp != "1788597000000" {
		t.Fatalf("creationTimestamp = %q, want epoch millis string", manifest.CreationTimestamp)
	}
	if manifest.DestinationBucket != "arn:aws:s3:::dest" || manifest.Files[0].Key != "src/id/data/x.csv.gz" {
		t.Fatalf("manifest = %+v", manifest)
	}
	encoded, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	for _, field := range []string{"sourceBucket", "destinationBucket", "creationTimestamp", "fileFormat", "fileSchema", "files", "MD5checksum"} {
		if !strings.Contains(string(encoded), `"`+field+`"`) {
			t.Fatalf("manifest json missing field %q: %s", field, encoded)
		}
	}
}

func TestInventoryReportKeyLayout(t *testing.T) {
	windowStamp := "2026-11-06T21-32-10Z"

	data := inventoryDataKey("", "src-bucket", "report1", windowStamp, "csv.gz")
	if data != "src-bucket/report1/data/delivery-2026-11-06T21-32-10Z.csv.gz" {
		t.Fatalf("data key = %q", data)
	}
	manifest := inventoryManifestKey("prefix1", "src-bucket", "report1", windowStamp)
	if manifest != "prefix1/src-bucket/report1/2026-11-06T21-32-10Z/manifest.json" {
		t.Fatalf("manifest key = %q", manifest)
	}
	symlink := inventorySymlinkKey("prefix1", "src-bucket", "report1", windowStamp)
	if symlink != "prefix1/src-bucket/report1/hive/dt=2026-11-06T21-32-10Z/symlink.txt" {
		t.Fatalf("symlink key = %q", symlink)
	}
}

// TestInventoryDeliveryKeysAreStablePerWindow pins the retry-convergence
// contract: every delivery key derives from the schedule window's anchor
// stamp, so a partially failed delivery rewrites the same keys on retry
// instead of piling up orphans, and a different window produces different
// keys.
func TestInventoryDeliveryKeysAreStablePerWindow(t *testing.T) {
	anchor := time.Date(2026, 11, 6, 21, 32, 10, 0, time.UTC)
	next := anchor.Add(45 * time.Second)
	stampOf := func(lastDelivery time.Time) string {
		return lastDelivery.UTC().Format("2006-01-02T15-04-05Z")
	}

	build := func(stamp string) (string, string, string) {
		return inventoryDataKey("prefix1", "src-bucket", "report1", stamp, "csv.gz"),
			inventoryManifestKey("prefix1", "src-bucket", "report1", stamp),
			inventorySymlinkKey("prefix1", "src-bucket", "report1", stamp)
	}

	dataA, manifestA, symlinkA := build(stampOf(anchor))
	dataA2, manifestA2, symlinkA2 := build(stampOf(anchor))
	if dataA != dataA2 || manifestA != manifestA2 || symlinkA != symlinkA2 {
		t.Fatalf("keys must be stable within one window: %q %q %q vs %q %q %q", dataA, manifestA, symlinkA, dataA2, manifestA2, symlinkA2)
	}

	dataB, manifestB, symlinkB := build(stampOf(next))
	if dataA == dataB || manifestA == manifestB || symlinkA == symlinkB {
		t.Fatalf("a different window must produce different keys: %q %q %q vs %q %q %q", dataA, manifestA, symlinkA, dataB, manifestB, symlinkB)
	}
}

func TestNextScheduleBoundary(t *testing.T) {
	// Thursday 2026-09-03 10:30 UTC: the next daily boundary is Friday
	// midnight; the next weekly boundary is Sunday midnight.
	anchor := time.Date(2026, 9, 3, 10, 30, 0, 0, time.UTC)
	if got := nextScheduleBoundary(anchor, "Daily"); !got.Equal(time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("daily boundary = %v", got)
	}
	if got := nextScheduleBoundary(anchor, "Weekly"); !got.Equal(time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("weekly boundary = %v", got)
	}
	// Exactly at midnight the boundary must move to the next day, not the
	// same instant (strictly after the anchor).
	midnight := time.Date(2026, 9, 4, 0, 0, 0, 0, time.UTC)
	if got := nextScheduleBoundary(midnight, "Daily"); !got.Equal(time.Date(2026, 9, 5, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("boundary at midnight = %v", got)
	}
}

// TestDeliveryFailureBackoffWindows pins the consecutive-failure backoff of
// the report worker: the window doubles per failure up to eight poll
// intervals, success clears the state, and pruning drops deleted
// configurations.
func TestDeliveryFailureBackoffWindows(t *testing.T) {
	w := &InventoryReportWorker{interval: 5 * time.Minute, failures: map[string]*deliveryFailure{}}
	base := time.Date(2026, 9, 5, 12, 0, 0, 0, time.UTC)
	key := w.failureKey("us-east-1", "bucket", "cfg")
	other := w.failureKey("us-east-1", "bucket", "other")

	if w.backedOff(key, base) {
		t.Fatal("a configuration without failures must never be backed off")
	}

	w.noteFailure(key, base)
	if !w.backedOff(key, base.Add(4*time.Minute)) {
		t.Fatal("first failure must back off for one interval")
	}
	if w.backedOff(key, base.Add(5*time.Minute)) {
		t.Fatal("first failure backoff must expire after one interval")
	}

	w.noteFailure(key, base.Add(6*time.Minute))
	if !w.backedOff(key, base.Add(15*time.Minute)) {
		t.Fatal("second failure must back off for two intervals")
	}
	if w.backedOff(key, base.Add(16*time.Minute)) {
		t.Fatal("second failure backoff must expire after two intervals")
	}

	for i := 0; i < 4; i++ {
		w.noteFailure(key, base)
	}
	// Six failures: 5m << 5 = 160m, capped at eight intervals (40m).
	if !w.backedOff(key, base.Add(39*time.Minute)) {
		t.Fatal("repeated failures must stay within the capped backoff window")
	}
	if w.backedOff(key, base.Add(40*time.Minute)) {
		t.Fatal("backoff window must be capped at eight intervals")
	}

	w.noteSuccess(key)
	if w.backedOff(key, base.Add(time.Second)) {
		t.Fatal("success must clear the failure state")
	}

	w.noteFailure(other, base)
	w.pruneFailures(map[string]bool{key: true})
	if w.backedOff(other, base) {
		t.Fatal("pruning must drop the failure state of configurations that no longer exist")
	}
}
