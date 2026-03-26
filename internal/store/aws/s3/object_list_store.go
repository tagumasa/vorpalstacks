package s3

import (
	"strings"

	pb "vorpalstacks/internal/pb/storage/storage_s3"

	"google.golang.org/protobuf/proto"
)

// List lists objects in a bucket with optional prefix and delimiter.
func (s *ObjectStore) List(bucket, prefix, delimiter, marker string, maxKeys int) (*ObjectListResult, error) {
	if maxKeys <= 0 {
		maxKeys = 1000
	}

	var objects []*Object
	var commonPrefixes []string
	prefixSet := make(map[string]bool)
	count := 0
	started := marker == ""
	skipUntilNewKey := ""
	hasMore := false
	wasVersioned := s.wasVersioned(bucket)

	var searchPrefix string
	if wasVersioned {
		searchPrefix = bucket + "#" + prefix
	} else {
		searchPrefix = s.storageKey(bucket, prefix)
	}

	err := s.ScanPrefix(searchPrefix, func(key string, value []byte) error {
		var pbObj pb.Object
		if err := proto.Unmarshal(value, &pbObj); err != nil {
			return err
		}
		obj := ProtoToObject(&pbObj)

		if !started {
			if obj.Key > marker {
				started = true
			} else {
				if obj.Key == marker {
					started = true
					skipUntilNewKey = marker
				}
				return nil
			}
		}

		if skipUntilNewKey != "" && obj.Key == skipUntilNewKey {
			return nil
		}

		if wasVersioned && strings.HasSuffix(key, "#_latest") {
			return nil
		}

		if wasVersioned && !obj.IsLatest {
			return nil
		}

		if obj.IsDeleteMarker {
			return nil
		}

		if delimiter != "" {
			trimmed := strings.TrimPrefix(obj.Key, prefix)
			idx := strings.Index(trimmed, delimiter)
			if idx >= 0 {
				commonPrefix := obj.Key[:len(prefix)+idx+1]
				if !prefixSet[commonPrefix] {
					prefixSet[commonPrefix] = true
					if count < maxKeys {
						commonPrefixes = append(commonPrefixes, commonPrefix)
						count++
					} else {
						hasMore = true
					}
				}
				return nil
			}
		}

		if count < maxKeys {
			objects = append(objects, obj)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	result := &ObjectListResult{
		Objects:        objects,
		CommonPrefixes: commonPrefixes,
		IsTruncated:    hasMore,
	}

	if len(objects) > 0 {
		result.NextMarker = objects[len(objects)-1].Key
	} else if len(commonPrefixes) > 0 {
		result.NextMarker = commonPrefixes[len(commonPrefixes)-1]
	}

	return result, nil
}

// ListObjectVersions lists object versions in a bucket.
func (s *ObjectStore) ListObjectVersions(bucket, prefix, delimiter, keyMarker, versionIdMarker string, maxKeys int) (*ObjectListResult, error) {
	if maxKeys <= 0 {
		maxKeys = 1000
	}

	var versions []*Object
	var commonPrefixes []string
	prefixSet := make(map[string]bool)
	count := 0
	started := keyMarker == ""
	hasMore := false

	searchPrefix := bucket + "#" + prefix

	err := s.ScanPrefix(searchPrefix, func(key string, value []byte) error {
		var pbObj pb.Object
		if err := proto.Unmarshal(value, &pbObj); err != nil {
			return err
		}
		obj := ProtoToObject(&pbObj)

		if strings.HasSuffix(key, "#_latest") {
			return nil
		}

		if !started {
			if obj.Key == keyMarker && (versionIdMarker == "" || obj.VersionID == versionIdMarker) {
				started = true
				return nil
			} else if obj.Key > keyMarker {
				started = true
			} else {
				return nil
			}
		}

		if prefix != "" && delimiter != "" {
			idx := strings.Index(strings.TrimPrefix(obj.Key, prefix), delimiter)
			if idx >= 0 {
				commonPrefix := prefix + obj.Key[len(prefix):len(prefix)+idx+1]
				if !prefixSet[commonPrefix] {
					prefixSet[commonPrefix] = true
					if count < maxKeys {
						commonPrefixes = append(commonPrefixes, commonPrefix)
						count++
					} else {
						hasMore = true
					}
				}
				return nil
			}
		}

		if count < maxKeys {
			versions = append(versions, obj)
			count++
		} else {
			hasMore = true
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	result := &ObjectListResult{
		Objects:        versions,
		CommonPrefixes: commonPrefixes,
		IsTruncated:    hasMore,
	}

	if len(versions) > 0 {
		result.NextVersionKeyMarker = versions[len(versions)-1].Key
		result.NextVersionIDMarker = versions[len(versions)-1].VersionID
	}

	return result, nil
}

// CountByBucket returns the number of objects in a bucket.
func (s *ObjectStore) CountByBucket(bucket string) (int, error) {
	count := 0
	isVersioned := s.isVersioningEnabled(bucket)

	if isVersioned {
		latestPrefix := bucket + "#"
		seenKeys := make(map[string]bool)
		err := s.ScanPrefix(latestPrefix, func(key string, value []byte) error {
			var pbObj pb.Object
			if err := proto.Unmarshal(value, &pbObj); err != nil {
				return err
			}
			obj := ProtoToObject(&pbObj)

			if strings.HasSuffix(key, "#_latest") {
				if !obj.IsDeleteMarker && !seenKeys[obj.Key] {
					seenKeys[obj.Key] = true
					count++
				}
			}
			return nil
		})
		return count, err
	}

	prefix := s.storageKey(bucket, "")
	err := s.ScanPrefix(prefix, func(key string, value []byte) error {
		var pbObj pb.Object
		if err := proto.Unmarshal(value, &pbObj); err != nil {
			return err
		}
		obj := ProtoToObject(&pbObj)
		if !obj.IsDeleteMarker {
			count++
		}
		return nil
	})
	return count, err
}

// CountMultipartUploadsByBucket returns the number of in-progress multipart uploads in a bucket.
func (s *ObjectStore) CountMultipartUploadsByBucket(bucket string) (int, error) {
	count := 0
	err := s.storage.Bucket(multipartBucketName(s.region)).ForEach(func(k, v []byte) error {
		var upload pb.MultipartUpload
		if err := proto.Unmarshal(v, &upload); err != nil {
			return err
		}
		if upload.BucketName == bucket {
			count++
		}
		return nil
	})
	return count, err
}
