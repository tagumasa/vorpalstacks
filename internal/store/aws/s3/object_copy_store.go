package s3

import (
	"context"
	"io"
)

// CopyOverrides selects what an object copy replaces from the source. The
// zero value copies the source's current version with its content type,
// metadata, content-* fields and storage class intact.
type CopyOverrides struct {
	// VersionId addresses a specific source version instead of the
	// current one.
	VersionId string
	// ContentType replaces the source content type when non-empty.
	ContentType string
	// The content-* fields replace the source values when metadata is
	// replaced (they only take effect under ReplaceMetadata, mirroring the
	// metadata directive contract).
	ContentEncoding    string
	ContentDisposition string
	ContentLanguage    string
	CacheControl       string
	// ReplaceMetadata swaps Metadata in for the source metadata map.
	ReplaceMetadata bool
	Metadata        map[string]string
	// StorageClass replaces the source storage class when non-empty.
	StorageClass ObjectStorageClass
}

// CopyObject copies an object from one location to another. It is the
// single path for every copy flavour: it opens the source (a specific
// version when the overrides address one), defaults the content type and
// storage class from the source object, applies the overrides, and
// persists the copy through PutWithVersioning.
func (s *ObjectStore) CopyObject(ctx context.Context, srcBucket, srcKey, dstBucket, dstKey string, overrides *CopyOverrides) (*Object, error) {
	if overrides == nil {
		overrides = &CopyOverrides{}
	}

	reader, srcObj, err := s.openObjectReader(ctx, srcBucket, srcKey, overrides.VersionId)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	contentType := overrides.ContentType
	if contentType == "" {
		contentType = srcObj.ContentType
	}
	metadata := srcObj.Metadata
	if overrides.ReplaceMetadata {
		metadata = overrides.Metadata
	}
	// The content-* fields follow the same directive as the metadata map:
	// a plain copy carries them from the source, a metadata replacement
	// installs the overrides' values.
	sysMeta := &SystemMetadata{
		ContentEncoding:    srcObj.ContentEncoding,
		ContentLanguage:    srcObj.ContentLanguage,
		ContentDisposition: srcObj.ContentDisposition,
		CacheControl:       srcObj.CacheControl,
	}
	if overrides.ReplaceMetadata {
		sysMeta.ContentEncoding = overrides.ContentEncoding
		sysMeta.ContentLanguage = overrides.ContentLanguage
		sysMeta.ContentDisposition = overrides.ContentDisposition
		sysMeta.CacheControl = overrides.CacheControl
	}
	storageClass := overrides.StorageClass
	if storageClass == "" {
		storageClass = srcObj.StorageClass
	}
	if storageClass == "" {
		storageClass = StorageClassStandard
	}

	return s.PutWithVersioning(ctx, dstBucket, dstKey, reader, contentType, metadata, false, storageClass, sysMeta)
}

// openObjectReader returns a reader and source object metadata for the given
// bucket/key. An explicit versionId addresses that version by record layout
// — the record and blob tiers resolve it regardless of the bucket's current
// versioning status, so a suspended bucket still copies its pre-suspension
// versions — while an empty versionId reads the current version.
func (s *ObjectStore) openObjectReader(ctx context.Context, bucket, key, versionId string) (io.ReadCloser, *Object, error) {
	if versionId != "" {
		return s.GetWithVersion(ctx, bucket, key, versionId)
	}
	return s.Get(ctx, bucket, key)
}
