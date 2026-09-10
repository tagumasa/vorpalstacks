package s3

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// fakeBlobStore is an in-memory BlobStore standing in for the file-tier
// blob backend in store-level tests. Only the multipart surface has real
// behaviour; the remaining interface methods are inert stubs.
type fakeBlobStore struct {
	mu     sync.Mutex
	nextID int
}

func (f *fakeBlobStore) CreateMultipartUpload(ctx context.Context, bucket, key string, metadata *storage.BlobMetadata) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.nextID++
	return fmt.Sprintf("upload-%03d", f.nextID), nil
}

func (f *fakeBlobStore) UploadPart(ctx context.Context, bucket, key, uploadID string, partNumber int, reader io.Reader) (string, error) {
	return fmt.Sprintf("etag-%s-%d", uploadID, partNumber), nil
}

func (f *fakeBlobStore) CompleteMultipartUpload(ctx context.Context, bucket, key, versionID, uploadID string, parts []storage.PartInfo) (*storage.BlobMetadata, error) {
	return &storage.BlobMetadata{Key: key, Size: int64(len(parts)), ETag: "completed-etag", LastModified: time.Now().UTC()}, nil
}

func (f *fakeBlobStore) AbortMultipartUpload(ctx context.Context, bucket, key, uploadID string) error {
	return nil
}

func (f *fakeBlobStore) ListParts(ctx context.Context, bucket, key, uploadID string) ([]storage.PartInfo, error) {
	return nil, nil
}

func (f *fakeBlobStore) Put(ctx context.Context, bucket, key string, reader io.Reader, metadata *storage.BlobMetadata) (*storage.BlobMetadata, error) {
	return metadata, nil
}

// emptyReadCloser satisfies storage.BlobReader for blob reads that carry
// no payload in tests.
type emptyReadCloser struct{}

func (emptyReadCloser) Read([]byte) (int, error) { return 0, io.EOF }
func (emptyReadCloser) Close() error             { return nil }
func (emptyReadCloser) Size() int64              { return 0 }
func (emptyReadCloser) ETag() string             { return "" }

func (f *fakeBlobStore) Get(ctx context.Context, bucket, key string) (storage.BlobReader, *storage.BlobMetadata, error) {
	return emptyReadCloser{}, &storage.BlobMetadata{Key: key}, nil
}

func (f *fakeBlobStore) GetRange(ctx context.Context, bucket, key string, offset, length int64) (storage.BlobReader, *storage.BlobMetadata, error) {
	return emptyReadCloser{}, &storage.BlobMetadata{Key: key}, nil
}

func (f *fakeBlobStore) Delete(ctx context.Context, bucket, key string) error {
	return nil
}

func (f *fakeBlobStore) Exists(ctx context.Context, bucket, key string) (bool, error) {
	return false, nil
}

func (f *fakeBlobStore) Head(ctx context.Context, bucket, key string) (*storage.BlobMetadata, error) {
	return &storage.BlobMetadata{Key: key}, nil
}

func (f *fakeBlobStore) PutWithVersion(ctx context.Context, bucket, key, versionId string, reader io.Reader, metadata *storage.BlobMetadata) (*storage.BlobMetadata, error) {
	return metadata, nil
}

func (f *fakeBlobStore) GetWithVersion(ctx context.Context, bucket, key, versionId string) (storage.BlobReader, *storage.BlobMetadata, error) {
	return emptyReadCloser{}, &storage.BlobMetadata{Key: key}, nil
}

func (f *fakeBlobStore) GetRangeWithVersion(ctx context.Context, bucket, key, versionId string, offset, length int64) (storage.BlobReader, *storage.BlobMetadata, error) {
	return emptyReadCloser{}, &storage.BlobMetadata{Key: key}, nil
}

func (f *fakeBlobStore) DeleteWithVersion(ctx context.Context, bucket, key, versionId string) error {
	return nil
}

func (f *fakeBlobStore) CopyWithVersion(ctx context.Context, srcBucket, srcKey, srcVersionId, dstBucket, dstKey string) (*storage.BlobMetadata, error) {
	return nil, nil
}

func (f *fakeBlobStore) Copy(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string) (*storage.BlobMetadata, error) {
	return nil, nil
}

func (f *fakeBlobStore) CreateBucket(ctx context.Context, name string) error {
	return nil
}

func (f *fakeBlobStore) DeleteBucket(ctx context.Context, name string) error {
	return nil
}

// newMultipartTestStore builds an ObjectStore over a throwaway Pebble
// instance and the in-memory blob fake.
func newMultipartTestStore(t *testing.T) (*ObjectStore, *BucketStore) {
	t.Helper()
	st, err := storage.NewPebbleStorage(&storage.Config{Path: t.TempDir()})
	if err != nil {
		t.Fatalf("pebble storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	bucketStore := NewBucketStore(st, "123456789012", "test-region")
	objectStore, err := NewObjectStore(st, &fakeBlobStore{}, bucketStore, "123456789012", "test-region")
	if err != nil {
		t.Fatalf("object store: %v", err)
	}
	return objectStore, bucketStore
}

func createTestUpload(t *testing.T, store *ObjectStore, bucket, key string) *MultipartUpload {
	t.Helper()
	upload, err := store.CreateMultipartUpload(context.Background(), bucket, key, "", nil, "", "", "", nil, nil, "", nil)
	if err != nil {
		t.Fatalf("CreateMultipartUpload(%s): %v", key, err)
	}
	return upload
}

// Concurrent part uploads on one upload ID must all survive: the upload
// record is a read-modify-write cycle, so a broken per-key lock (the
// unlock-then-delete pattern) loses the loser's part. The AWS SDK upload
// managers upload parts concurrently by default.
func TestUploadPartConcurrentPartsAllSurvive(t *testing.T) {
	store, _ := newMultipartTestStore(t)
	upload := createTestUpload(t, store, "bkt", "key.txt")

	const workers = 4
	const iters = 25
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := 0; i < iters; i++ {
				partNumber := w*iters + i + 1
				if _, err := store.UploadPart(context.Background(), "bkt", "key.txt", upload.UploadID, partNumber, strings.NewReader("x"), 0, 1, nil, nil); err != nil {
					t.Errorf("UploadPart(%d): %v", partNumber, err)
					return
				}
			}
		}(w)
	}
	wg.Wait()

	got, err := store.GetMultipartUpload(upload.UploadID)
	if err != nil {
		t.Fatalf("GetMultipartUpload: %v", err)
	}
	if len(got.Parts) != workers*iters {
		t.Fatalf("parts = %d, want %d — concurrent uploads lost parts", len(got.Parts), workers*iters)
	}
}

// AbortMultipartUpload addresses an upload by the bucket/key/uploadId
// triple. A valid uploadId under the wrong bucket or key addresses no
// upload at that address — AWS reports NoSuchUpload — and must leave the
// real upload untouched instead of destroying it.
func TestAbortMultipartUploadAddressMismatch(t *testing.T) {
	store, _ := newMultipartTestStore(t)
	upload := createTestUpload(t, store, "bkt", "key.txt")
	if _, err := store.UploadPart(context.Background(), "bkt", "key.txt", upload.UploadID, 1, strings.NewReader("x"), 0, 1, nil, nil); err != nil {
		t.Fatalf("UploadPart: %v", err)
	}

	for _, addr := range [][2]string{
		{"bkt", "other.txt"},
		{"other-bkt", "key.txt"},
	} {
		err := store.AbortMultipartUpload(context.Background(), addr[0], addr[1], upload.UploadID)
		if err == nil || !errors.Is(err, ErrUploadNotFound) {
			t.Fatalf("AbortMultipartUpload(%s, %s) = %v, want ErrUploadNotFound", addr[0], addr[1], err)
		}

		got, getErr := store.GetMultipartUpload(upload.UploadID)
		if getErr != nil {
			t.Fatalf("upload was destroyed by a mismatched abort: %v", getErr)
		}
		if len(got.Parts) != 1 {
			t.Fatalf("parts = %d after mismatched abort, want 1", len(got.Parts))
		}
	}

	if err := store.AbortMultipartUpload(context.Background(), "bkt", "key.txt", upload.UploadID); err != nil {
		t.Fatalf("AbortMultipartUpload(correct address): %v", err)
	}
	if _, err := store.GetMultipartUpload(upload.UploadID); !errors.Is(err, ErrUploadNotFound) {
		t.Fatalf("GetMultipartUpload after correct abort = %v, want ErrUploadNotFound", err)
	}
}

// ListMultipartUploads must honour the documented marker contract: without
// an upload-id-marker only keys greater than the key-marker are included;
// with one, equal-key uploads with greater upload IDs are included too; and
// a marker key with no active uploads must not truncate the listing.
func TestListMultipartUploadsMarkerSemantics(t *testing.T) {
	store, _ := newMultipartTestStore(t)

	createTestUpload(t, store, "bkt", "a")
	bFirst := createTestUpload(t, store, "bkt", "b")
	bSecond := createTestUpload(t, store, "bkt", "b")
	cUpload := createTestUpload(t, store, "bkt", "c")

	if bSecond.UploadID <= bFirst.UploadID {
		t.Fatalf("test requires b's upload IDs in lexicographic order, got %s then %s", bFirst.UploadID, bSecond.UploadID)
	}

	keys := func(resp *MultipartUploadListResult) []string {
		var out []string
		for _, u := range resp.Uploads {
			out = append(out, u.Key)
		}
		return out
	}
	entries := func(resp *MultipartUploadListResult) []string {
		var out []string
		for _, u := range resp.Uploads {
			out = append(out, u.Key+"/"+u.UploadID)
		}
		return out
	}

	// Pair marker: equal-key uploads with upload IDs greater than the
	// marker must be included alongside later keys.
	resp, err := store.ListMultipartUploads("bkt", "", "b", bFirst.UploadID, 100)
	if err != nil {
		t.Fatalf("ListMultipartUploads(pair): %v", err)
	}
	got := strings.Join(entries(resp), ",")
	want := "b/" + bSecond.UploadID + ",c/" + cUpload.UploadID
	if got != want {
		t.Fatalf("pair-marker listing = %q, want %q", got, want)
	}

	// Key marker only: equal-key uploads are excluded, later keys included.
	resp, err = store.ListMultipartUploads("bkt", "", "b", "", 100)
	if err != nil {
		t.Fatalf("ListMultipartUploads(key only): %v", err)
	}
	if got := strings.Join(keys(resp), ","); got != "c" {
		t.Fatalf("key-only listing = %q, want c", got)
	}

	// Marker key with no active uploads must not truncate the listing.
	if err := store.AbortMultipartUpload(context.Background(), "bkt", "b", bFirst.UploadID); err != nil {
		t.Fatalf("AbortMultipartUpload(b/first): %v", err)
	}
	if err := store.AbortMultipartUpload(context.Background(), "bkt", "b", bSecond.UploadID); err != nil {
		t.Fatalf("AbortMultipartUpload(b/second): %v", err)
	}
	resp, err = store.ListMultipartUploads("bkt", "", "b", "", 100)
	if err != nil {
		t.Fatalf("ListMultipartUploads(deleted marker key): %v", err)
	}
	if got := strings.Join(keys(resp), ","); got != "c" {
		t.Fatalf("deleted-marker-key listing = %q, want c (the marker key having no active uploads must not truncate the listing)", got)
	}

	// A pair marker whose own entry has been aborted must still include the
	// equal-key uploads that sort after the upload-id-marker.
	bThird := createTestUpload(t, store, "bkt", "b")
	resp, err = store.ListMultipartUploads("bkt", "", "b", bFirst.UploadID, 100)
	if err != nil {
		t.Fatalf("ListMultipartUploads(deleted pair-marker entry): %v", err)
	}
	got = strings.Join(entries(resp), ",")
	want = "b/" + bThird.UploadID + ",c/" + cUpload.UploadID
	if got != want {
		t.Fatalf("deleted-pair-marker listing = %q, want %q", got, want)
	}
}

// ListParts returns the parts and the upload record they were read from in
// one result, so a caller renders Initiator, Owner and StorageClass without
// a second lookup that an interleaved abort could fail after the page was
// already assembled. A valid uploadId under the wrong key addresses no
// upload and reports ErrUploadNotFound.
func TestListPartsReturnsUploadRecordAndParts(t *testing.T) {
	store, _ := newMultipartTestStore(t)
	upload := createTestUpload(t, store, "bkt", "key.txt")
	for _, pn := range []int{1, 2} {
		if _, err := store.UploadPart(context.Background(), "bkt", "key.txt", upload.UploadID, pn, strings.NewReader("x"), 0, 1, nil, nil); err != nil {
			t.Fatalf("UploadPart(%d): %v", pn, err)
		}
	}

	first, err := store.ListParts(context.Background(), "bkt", "key.txt", upload.UploadID, 0, 1)
	if err != nil {
		t.Fatalf("ListParts(first page): %v", err)
	}
	if first.Upload == nil {
		t.Fatal("ListParts returned no upload record with the parts")
	}
	if first.Upload.UploadID != upload.UploadID || first.Upload.BucketName != "bkt" || first.Upload.Key != "key.txt" {
		t.Fatalf("record identity = %s/%s/%s, want %s/bkt/key.txt", first.Upload.UploadID, first.Upload.BucketName, first.Upload.Key, upload.UploadID)
	}
	if len(first.Parts) != 1 || first.Parts[0].PartNumber != 1 {
		t.Fatalf("first page parts = %v, want part 1 only", first.Parts)
	}
	if !first.IsTruncated || first.NextPartNumberMarker != 1 {
		t.Fatalf("first page truncated = %v marker = %d, want true/1", first.IsTruncated, first.NextPartNumberMarker)
	}

	second, err := store.ListParts(context.Background(), "bkt", "key.txt", upload.UploadID, first.NextPartNumberMarker, 1)
	if err != nil {
		t.Fatalf("ListParts(second page): %v", err)
	}
	if second.Upload == nil || second.Upload.UploadID != upload.UploadID {
		t.Fatal("second page carries no upload record")
	}
	if len(second.Parts) != 1 || second.Parts[0].PartNumber != 2 || second.IsTruncated {
		t.Fatalf("second page = %v truncated = %v, want part 2 only, not truncated", second.Parts, second.IsTruncated)
	}

	if _, err := store.ListParts(context.Background(), "bkt", "other.txt", upload.UploadID, 0, 100); !errors.Is(err, ErrUploadNotFound) {
		t.Fatalf("ListParts(wrong key) = %v, want ErrUploadNotFound", err)
	}
}

// A missing upload is ErrUploadNotFound whatever the page size: the record
// is resolved before the empty-page branch, so the not-found outcome does
// not depend on the caller making a second lookup.
func TestListPartsZeroMaxPartsMissingUploadIsNoSuchUpload(t *testing.T) {
	store, _ := newMultipartTestStore(t)

	if _, err := store.ListParts(context.Background(), "bkt", "key.txt", "no-such-upload", 0, 0); !errors.Is(err, ErrUploadNotFound) {
		t.Fatalf("ListParts(missing upload, zero page) = %v, want ErrUploadNotFound", err)
	}

	upload := createTestUpload(t, store, "bkt", "key.txt")
	empty, err := store.ListParts(context.Background(), "bkt", "key.txt", upload.UploadID, 0, 0)
	if err != nil {
		t.Fatalf("ListParts(existing upload, zero page): %v", err)
	}
	if len(empty.Parts) != 0 || empty.IsTruncated {
		t.Fatalf("zero-page result = %v truncated = %v, want an empty non-truncated page", empty.Parts, empty.IsTruncated)
	}
	if empty.Upload == nil || empty.Upload.UploadID != upload.UploadID {
		t.Fatal("zero-page result carries no upload record")
	}
}
