package s3

import (
	"strings"
	"testing"
	"time"
)

// Every list-response timestamp carries milliseconds — the format the S3
// API reference examples render (2009-10-12T17:50:30.000Z) across object
// contents, version entries, delete markers, parts, and upload initiations.
func TestListXMLTimestampsCarryMilliseconds(t *testing.T) {
	stamp := time.Date(2026, 9, 9, 12, 0, 0, 500000000, time.UTC)

	contents := &ListObjectsV2Output{
		Contents: []*ObjectContent{{Key: "a.txt", LastModified: stamp, ETag: `"e"`, Size: 1, StorageClass: "STANDARD"}},
	}
	if xml := contents.ToXML(); !strings.Contains(xml, "2026-09-09T12:00:00.500Z") {
		t.Fatalf("Contents LastModified lacks the millisecond format: %s", xml)
	}

	versions := &ListObjectVersionsOutput{
		Versions:      []*ObjectVersion{{Key: "a.txt", LastModified: stamp, VersionId: "v1"}},
		DeleteMarkers: []*DeleteMarkerEntry{{Key: "b.txt", LastModified: stamp, VersionId: "m1"}},
	}
	if xml := versions.ToXML(); !strings.Contains(xml, "2026-09-09T12:00:00.500Z") {
		t.Fatalf("Version/DeleteMarker LastModified lacks the millisecond format: %s", xml)
	}

	parts := &ListPartsOutput{
		Parts: []*Part{{PartNumber: 1, LastModified: stamp, ETag: `"e"`, Size: 1}},
	}
	if xml := parts.ToXML(); !strings.Contains(xml, "2026-09-09T12:00:00.500Z") {
		t.Fatalf("Part LastModified lacks the millisecond format: %s", xml)
	}

	uploads := &ListMultipartUploadsOutput{
		Uploads: []*Upload{{Key: "a.txt", UploadId: "u1", StorageClass: "STANDARD", Initiated: stamp}},
	}
	if xml := uploads.ToXML(); !strings.Contains(xml, "2026-09-09T12:00:00.500Z") {
		t.Fatalf("Upload Initiated lacks the millisecond format: %s", xml)
	}
}
