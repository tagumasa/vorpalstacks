package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	s3store "vorpalstacks/internal/store/aws/s3"
	"vorpalstacks/internal/utils/timeutils"
)

// Object metadata, delete and admin-copy cores for both transport planes.
// The streaming put/get/copy engines live in object_stream_core.go.

// AdminHeadObjectInput carries the fields needed for HeadObject. The
// conditional fields mirror GetObjectStreamInput's so both read operations
// evaluate preconditions in Core; the admin plane leaves them unset.
type AdminHeadObjectInput struct {
	Bucket            string
	Key               string
	VersionID         string
	IfMatch           string
	IfNoneMatch       string
	IfModifiedSince   *time.Time
	IfUnmodifiedSince *time.Time
}

// AdminGetObjectInput carries the fields needed for GetObject.
type AdminGetObjectInput struct {
	Bucket    string
	Key       string
	VersionID string
}

// AdminPutObjectInput carries the fields needed for PutObject.
type AdminPutObjectInput struct {
	Bucket      string
	Key         string
	Body        []byte
	ContentType string
	Metadata    map[string]string
}

// AdminDeleteObjectInput carries the fields needed for DeleteObject.
type AdminDeleteObjectInput struct {
	Bucket    string
	Key       string
	VersionID string
}

// AdminObjectIdentifier identifies a single object for bulk delete.
type AdminObjectIdentifier struct {
	Key       string
	VersionID string
}

// AdminDeleteObjectsInput carries the fields needed for DeleteObjects.
type AdminDeleteObjectsInput struct {
	Bucket  string
	Objects []AdminObjectIdentifier
}

// AdminCopyObjectInput carries the fields needed for CopyObject.
type AdminCopyObjectInput struct {
	Bucket       string
	Key          string
	CopySource   string
	ContentType  string
	StorageClass string
}

// ---------------------------------------------------------------------------
// Streaming input/result structs — used by both HTTP and admin handlers
// for operations involving io.Reader bodies.
// ---------------------------------------------------------------------------

// AdminHeadObjectResult holds the transport-agnostic result of HeadObject.
type AdminHeadObjectResult struct {
	Object     *s3store.Object
	Expiration string
}

// AdminGetObjectResult holds the transport-agnostic result of GetObject.
type AdminGetObjectResult struct {
	Object *s3store.Object
	Body   []byte
}

// AdminPutObjectResult holds the transport-agnostic result of PutObject.
type AdminPutObjectResult struct {
	ETag      string
	VersionID string
	Size      int64
	KMSKeyID  string
}

// AdminDeleteObjectResult holds the transport-agnostic result of DeleteObject.
type AdminDeleteObjectResult struct {
	VersionID      string
	IsDeleteMarker bool
}

// AdminDeletedObject holds info about a single successfully deleted object.
type AdminDeletedObject struct {
	Key                   string
	VersionID             string
	DeleteMarker          bool
	DeleteMarkerVersionID string
}

// AdminDeleteError holds info about a single failed deletion.
type AdminDeleteError struct {
	Key     string
	Code    string
	Message string
}

// AdminDeleteObjectsResult holds the transport-agnostic result of DeleteObjects.
type AdminDeleteObjectsResult struct {
	Deleted []AdminDeletedObject
	Errors  []AdminDeleteError
}

// AdminCopyObjectResult holds the transport-agnostic result of CopyObject.
type AdminCopyObjectResult struct {
	ETag         string
	LastModified string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path shared by admin
// handler only.  The HTTP API layer has its own operations in
// bucket_operations.go / object_*.go.
// ---------------------------------------------------------------------------

// headObjectCore retrieves metadata for an object without returning the body.
// Conditional requests are evaluated here, as GET's Core does, so both read
// operations share one precondition path.
func (s *S3Service) headObjectCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminHeadObjectInput) (*AdminHeadObjectResult, error) {
	obj, err := objectStore.HeadWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		return nil, mapVersionReadError(err, obj, in.Key, in.VersionID)
	}
	if in.IfMatch != "" || in.IfNoneMatch != "" || in.IfModifiedSince != nil || in.IfUnmodifiedSince != nil {
		if err := checkObjectPreconditions(obj, in.IfMatch, in.IfNoneMatch, in.IfModifiedSince, in.IfUnmodifiedSince); err != nil {
			return nil, err
		}
	}
	return &AdminHeadObjectResult{Object: obj, Expiration: s.objectExpirationHeaderFor(in.Bucket, obj)}, nil
}

// getObjectCore retrieves metadata and body for an object.
func (s *S3Service) getObjectCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminGetObjectInput) (*AdminGetObjectResult, error) {
	streamResult, err := s.getObjectStreamCore(ctx, objectStore, GetObjectStreamInput{
		Bucket:    in.Bucket,
		Key:       in.Key,
		VersionID: in.VersionID,
	})
	if err != nil {
		return nil, err
	}
	defer streamResult.Body.Close()

	data, err := io.ReadAll(streamResult.Body)
	if err != nil {
		return nil, err
	}

	return &AdminGetObjectResult{
		Object: &s3store.Object{
			ETag:               streamResult.ETag,
			LastModified:       streamResult.LastModified,
			ContentType:        streamResult.ContentType,
			Metadata:           streamResult.Metadata,
			Size:               streamResult.ContentLength,
			StorageClass:       s3store.ObjectStorageClass(streamResult.StorageClass),
			VersionID:          streamResult.VersionID,
			ContentEncoding:    streamResult.ContentEncoding,
			ContentLanguage:    streamResult.ContentLanguage,
			ContentDisposition: streamResult.ContentDisposition,
			CacheControl:       streamResult.CacheControl,
			ReplicationStatus:  streamResult.ReplicationStatus,
			SSEMetadata:        streamResult.SSEMetadata,
			Tags:               streamResult.Tags,
		},
		Body: data,
	}, nil
}

// putObjectCore validates the upload, determines encryption settings, and
// stores the object.
func (s *S3Service) putObjectCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminPutObjectInput) (*AdminPutObjectResult, error) {
	if in.Bucket == "" {
		return nil, NewInvalidArgumentError("bucket is required")
	}
	if in.Key == "" {
		return nil, NewInvalidArgumentError("key is required")
	}
	streamResult, err := s.putObjectStreamCore(ctx, bucketStore, objectStore, PutObjectStreamInput{
		Body:          bytes.NewReader(in.Body),
		ContentLength: int64(len(in.Body)),
		Bucket:        in.Bucket,
		Key:           in.Key,
		ContentType:   in.ContentType,
		Metadata:      in.Metadata,
	})
	if err != nil {
		return nil, err
	}
	obj := streamResult.Object
	result := &AdminPutObjectResult{
		ETag:      formatETag(obj.ETag),
		VersionID: obj.VersionID,
		Size:      obj.Size,
	}
	if obj.SSEMetadata != nil {
		result.KMSKeyID = obj.SSEMetadata.KMSKeyID
	}
	return result, nil
}

// deleteObjectCore deletes a single object. In a versioned bucket, deleting
// without a specific VersionID creates a delete marker.
func (s *S3Service) deleteObjectCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminDeleteObjectInput) (*AdminDeleteObjectResult, error) {
	if in.Bucket == "" {
		return nil, NewInvalidArgumentError("bucket is required")
	}
	if in.Key == "" {
		return nil, NewInvalidArgumentError("key is required")
	}
	result, err := objectStore.DeleteWithVersion(ctx, in.Bucket, in.Key, in.VersionID)
	if err != nil {
		// A non-versioned bucket only ever holds the null version, so a
		// version-addressed delete there references a version that does
		// not exist — the AWS-documented NoSuchVersion case.
		if errors.Is(err, s3store.ErrVersioningNotEnabled) {
			return nil, ErrNoSuchVersion
		}
		return nil, err
	}
	if result == nil {
		return &AdminDeleteObjectResult{}, nil
	}
	return &AdminDeleteObjectResult{
		VersionID:      result.VersionID,
		IsDeleteMarker: result.IsDeleteMarker,
	}, nil
}

// deleteObjectsCore deletes multiple objects, collecting per-object results.
func (s *S3Service) deleteObjectsCore(ctx context.Context, objectStore s3store.ObjectStoreInterface, in AdminDeleteObjectsInput) (*AdminDeleteObjectsResult, error) {
	if in.Bucket == "" {
		return nil, NewInvalidArgumentError("bucket is required")
	}
	if len(in.Objects) == 0 {
		return nil, NewInvalidArgumentError("no objects specified for deletion")
	}
	result := &AdminDeleteObjectsResult{}

	for _, obj := range in.Objects {
		if obj.Key == "" {
			result.Errors = append(result.Errors, AdminDeleteError{
				Key:     obj.Key,
				Code:    "InvalidArgument",
				Message: "object key is required",
			})
			continue
		}

		if obj.VersionID != "" {
			delResult, err := objectStore.DeleteWithVersion(ctx, in.Bucket, obj.Key, obj.VersionID)
			if err != nil {
				code := "InternalError"
				if errors.Is(err, s3store.ErrVersioningNotEnabled) {
					code = "NoSuchVersion"
				}
				result.Errors = append(result.Errors, AdminDeleteError{
					Key:     obj.Key,
					Code:    code,
					Message: err.Error(),
				})
				continue
			}
			deletedObj := AdminDeletedObject{
				Key:       obj.Key,
				VersionID: obj.VersionID,
			}
			// A version-addressed delete reports DeleteMarker only when the
			// removed version WAS a delete marker; removing a regular
			// version reports neither field.
			if delResult != nil && delResult.IsDeleteMarker {
				deletedObj.DeleteMarker = true
				deletedObj.DeleteMarkerVersionID = delResult.VersionID
			}
			result.Deleted = append(result.Deleted, deletedObj)
		} else {
			if err := objectStore.Delete(ctx, in.Bucket, obj.Key); err != nil {
				result.Errors = append(result.Errors, AdminDeleteError{
					Key:     obj.Key,
					Code:    "InternalError",
					Message: err.Error(),
				})
				continue
			}
			result.Deleted = append(result.Deleted, AdminDeletedObject{
				Key: obj.Key,
			})
		}
	}

	return result, nil
}

// copyObjectCore copies an object, handling encryption for the destination.
func (s *S3Service) copyObjectCore(ctx context.Context, bucketStore s3store.BucketStoreInterface, objectStore s3store.ObjectStoreInterface, in AdminCopyObjectInput) (*AdminCopyObjectResult, error) {
	if in.Bucket == "" {
		return nil, NewInvalidArgumentError("bucket is required")
	}
	if in.Key == "" {
		return nil, NewInvalidArgumentError("key is required")
	}
	if in.CopySource == "" {
		return nil, NewInvalidArgumentError("copy source is required")
	}
	streamResult, err := s.copyObjectStreamCore(ctx, bucketStore, objectStore, CopyObjectStreamInput{
		Bucket:       in.Bucket,
		Key:          in.Key,
		CopySource:   in.CopySource,
		ContentType:  in.ContentType,
		StorageClass: in.StorageClass,
	})
	if err != nil {
		return nil, err
	}
	obj := streamResult.Object
	return &AdminCopyObjectResult{
		ETag:         formatETag(obj.ETag),
		LastModified: obj.LastModified.Format(timeutils.ISO8601UTCFormat),
	}, nil
}

// checkObjectPreconditions evaluates the conditional request headers
// (If-Match, If-None-Match, If-Modified-Since, If-Unmodified-Since) against
// object metadata in the order RFC 7232 section 6 fixes — the order the S3
// API reference points to on every conditional header: If-Match, then
// If-Unmodified-Since only when If-Match is absent, then If-None-Match,
// then If-Modified-Since only when If-None-Match is absent. This yields the
// pairwise outcomes the S3 API reference documents: a true If-Match with a
// false If-Unmodified-Since returns 200 OK, and a false If-None-Match with
// a true If-Modified-Since returns 304 Not Modified. A failed If-Match or
// If-Unmodified-Since yields 412 PreconditionFailed, while a matching
// If-None-Match or an unmodified If-Modified-Since yields 304 NotModified.
// GET and HEAD share this evaluation.
func checkObjectPreconditions(obj *s3store.Object, ifMatch, ifNoneMatch string, ifModifiedSince, ifUnmodifiedSince *time.Time) error {
	if ifMatch != "" {
		if ifMatch == "*" {
			// Wildcard: object must exist; the caller resolved the object.
		} else if strings.Trim(obj.ETag, "\"") != strings.Trim(ifMatch, "\"") {
			return ErrPreconditionFailed
		}
	} else if ifUnmodifiedSince != nil && obj.LastModified.After(*ifUnmodifiedSince) {
		return ErrPreconditionFailed
	}
	if ifNoneMatch != "" {
		if ifNoneMatch == "*" {
			return ErrNotModified
		}
		if strings.Trim(obj.ETag, "\"") == strings.Trim(ifNoneMatch, "\"") {
			return ErrNotModified
		}
	} else if ifModifiedSince != nil && !obj.LastModified.After(*ifModifiedSince) {
		return ErrNotModified
	}
	return nil
}

// checkCopySourcePreconditions evaluates the x-amz-copy-source-if-* request
// headers against the copy source object in the RFC 7232 section 6 order
// the read preconditions use. A failed source precondition fails the whole
// copy with 412 PreconditionFailed: unlike a read, a copy cannot complete
// meaningfully with a 304 response, so the If-None-Match and
// If-Modified-Since failures yield 412 here as they would for any method
// other than GET and HEAD.
func checkCopySourcePreconditions(obj *s3store.Object, ifMatch, ifNoneMatch string, ifModifiedSince, ifUnmodifiedSince *time.Time) error {
	if ifMatch != "" {
		if ifMatch == "*" {
			// Wildcard: source object must exist; the caller resolved it.
		} else if strings.Trim(obj.ETag, "\"") != strings.Trim(ifMatch, "\"") {
			return ErrPreconditionFailed
		}
	} else if ifUnmodifiedSince != nil && obj.LastModified.After(*ifUnmodifiedSince) {
		return ErrPreconditionFailed
	}
	if ifNoneMatch != "" {
		if ifNoneMatch == "*" || strings.Trim(obj.ETag, "\"") == strings.Trim(ifNoneMatch, "\"") {
			return ErrPreconditionFailed
		}
	} else if ifModifiedSince != nil && !obj.LastModified.After(*ifModifiedSince) {
		return ErrPreconditionFailed
	}
	return nil
}

// resolvePartWindow maps a partNumber read to the part's byte window within
// the object.  Objects completed without persisted boundaries (plain uploads
// and multipart objects written before boundaries were kept) are treated as
// a single implicit part.  An unsatisfiable part number yields the
// documented 416 InvalidRange error, consistent with the API reference
// describing partNumber as "effectively performing a 'ranged' GET request".
func resolvePartWindow(parts []s3store.ObjectPartBoundary, partNumber int, totalSize int64) (start, length int64, partsCount int32, err error) {
	if len(parts) == 0 {
		if partNumber == 1 {
			// A plain object is a single implicit part; the parts-count
			// header is only documented for multipart-uploaded objects,
			// so no count is reported.
			return 0, totalSize, 0, nil
		}
		return 0, 0, 0, ErrInvalidRange
	}
	var offset int64
	for _, p := range parts {
		if p.PartNumber == partNumber {
			return offset, p.Size, int32(len(parts)), nil
		}
		offset += p.Size
	}
	return 0, 0, 0, ErrInvalidRange
}

// resolvePartRange maps a partNumber request to an absolute byte window
// within the object.  An optional Range header is evaluated within the
// selected part.  The returned end offset is inclusive.
func resolvePartRange(parts []s3store.ObjectPartBoundary, partNumber int, totalSize int64, rangeHeader string) (start, end int64, partsCount int32, err error) {
	partStart, partLength, count, err := resolvePartWindow(parts, partNumber, totalSize)
	if err != nil {
		return 0, 0, 0, err
	}
	start = partStart
	end = partStart + partLength - 1
	partsCount = count

	if rangeHeader != "" {
		ranges, rangeErr := parseRangeHeader(rangeHeader)
		if rangeErr != nil {
			return 0, 0, 0, rangeErr
		}
		r := ranges[0]
		var relStart, relLength int64
		if r.Start == -1 {
			relLength = r.Length
			relStart = partLength - relLength
			if relStart < 0 {
				relStart = 0
				relLength = partLength
			}
		} else {
			relStart = r.Start
			if r.Length == -1 {
				relLength = partLength - relStart
			} else {
				relLength = r.Length
			}
			if relLength < 0 {
				relLength = 0
			}
		}
		if relStart >= partLength && partLength > 0 {
			return 0, 0, 0, ErrInvalidRange
		}
		relEnd := relStart + relLength - 1
		if relEnd > partLength-1 {
			relEnd = partLength - 1
		}
		start += relStart
		end = start + (relEnd - relStart)
	}
	return start, end, partsCount, nil
}

// isArchiveClass reports whether a storage class places objects in an
// archive tier that must be restored before the object data can be read.
// GLACIER_IR is excluded because it offers real-time retrieval.
func isArchiveClass(cls s3store.ObjectStorageClass) bool {
	return cls == s3store.StorageClassGlacier || cls == s3store.StorageClassDeepArchive
}

// objectRestored reports whether an archived object currently has a
// temporary restored copy available for reads. The storage class is
// unchanged by a restore, so the expiry timestamp alone carries the state.
func objectRestored(obj *s3store.Object, now time.Time) bool {
	return obj.RestoreExpiry != nil && now.Before(*obj.RestoreExpiry)
}

// restoreHeaderValue renders the x-amz-restore response header for an
// object with an active temporary copy, e.g.
// ongoing-request="false", expiry-date="Wed, 12 Aug 2020 00:00:00 GMT".
// It returns the empty string when no restored copy is active.
func restoreHeaderValue(obj *s3store.Object, now time.Time) string {
	if !objectRestored(obj, now) {
		return ""
	}
	return fmt.Sprintf(`ongoing-request="false", expiry-date=%q`, obj.RestoreExpiry.UTC().Format(http.TimeFormat))
}

// nextRestoreExpiry computes when a temporary restored copy expires: the
// requested number of days is added to the completion time and the result
// is rounded up to the following midnight UTC, as documented by the
// restore API (a copy restored at 10:30 for 3 days expires at 00:00 on
// the day after restore-time + 3 days).
func nextRestoreExpiry(now time.Time, days int) time.Time {
	end := now.UTC().AddDate(0, 0, days)
	midnight := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
	if end.Equal(midnight) {
		return midnight
	}
	return midnight.AddDate(0, 0, 1)
}

// validateStorageClass checks that the value is either empty (defaults to
// STANDARD) or a storage class this platform persists. The acceptance set
// lives in the store package beside the storage-class constants, so the
// accepted values and the persisted enum cannot drift apart. It reads the
// store acceptance set directly, so it lives with the Core functions rather
// than in the transport-side validators file.
func validateStorageClass(sc string) error {
	if sc == "" {
		return nil
	}
	if !s3store.IsValidStorageClass(s3store.ObjectStorageClass(sc)) {
		return NewInvalidArgumentError(fmt.Sprintf("invalid StorageClass: %s", sc))
	}
	return nil
}
