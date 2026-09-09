package storage

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// A versioned write and the versioned reads must agree on the storage
// locations byte for byte. The historical bug: versioned writes munged
// "key#versionId" into the plain put path, so the file-tier path was
// sanitised with the version embedded in the final segment, while reads
// sanitised the bare key and appended the version afterwards — large
// versioned objects under keys whose final segment rewrites (trailing
// "/", "/.", "/..") were written but unreachable.
func TestVersionedBlobAddressesRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	bs, err := Open(tmpDir)
	require.NoError(t, err)
	defer bs.Close()
	s, err := NewHybridBlobStore(bs, tmpDir)
	require.NoError(t, err)
	ctx := context.Background()

	// Above the small-object threshold so every write lands in the file
	// tier, where the divergent paths lived.
	large := bytes.Repeat([]byte("x"), int(SmallObjectThreshold)+1)

	for _, key := range []string{"plain/key.txt", "dir/", "dots/.", "dots/..", "double//"} {
		_, err = s.PutWithVersion(ctx, "addr", key, "v1", bytes.NewReader(large), nil)
		require.NoError(t, err, "PutWithVersion(%q)", key)

		reader, meta, err := s.GetWithVersion(ctx, "addr", key, "v1")
		require.NoError(t, err, "GetWithVersion(%q)", key)
		got, readErr := io.ReadAll(reader)
		require.NoError(t, reader.Close())
		require.NoError(t, readErr)
		require.Len(t, got, len(large))
		require.EqualValues(t, len(large), meta.Size)

		rng, _, err := s.GetRangeWithVersion(ctx, "addr", key, "v1", 1, 10)
		require.NoError(t, err, "GetRangeWithVersion(%q)", key)
		_, readErr = io.ReadAll(rng)
		require.NoError(t, rng.Close())
		require.NoError(t, readErr)

		require.NoError(t, s.DeleteWithVersion(ctx, "addr", key, "v1"))
	}
}

// A key containing "#" must be its own object: the composed Pebble keys and
// file paths escape the separator, so the plain key "a#b" cannot alias the
// versioned address of version "b" of object "a" in either tier.
func TestHashKeyDoesNotAliasVersionedAddress(t *testing.T) {
	tmpDir := t.TempDir()
	bs, err := Open(tmpDir)
	require.NoError(t, err)
	defer bs.Close()
	s, err := NewHybridBlobStore(bs, tmpDir)
	require.NoError(t, err)
	ctx := context.Background()

	// Above the small-object threshold so both objects exercise the file
	// tier as well; the small tier is covered by the second round below.
	large := bytes.Repeat([]byte("x"), int(SmallObjectThreshold)+1)

	for _, round := range []struct{ plainBody, versionBody string }{
		{"small-plain", "small-versioned"},
		{string(large) + "p", string(large) + "v"},
	} {
		_, err = s.Put(ctx, "hash", "a#b", strings.NewReader(round.plainBody), nil)
		require.NoError(t, err)
		_, err = s.PutWithVersion(ctx, "hash", "a", "b", strings.NewReader(round.versionBody), nil)
		require.NoError(t, err)

		got, err := readAll(s.Get(ctx, "hash", "a#b"))
		require.NoError(t, err)
		require.Equal(t, round.plainBody, got, "plain hash key read")

		got, err = readAll(s.GetWithVersion(ctx, "hash", "a", "b"))
		require.NoError(t, err)
		require.Equal(t, round.versionBody, got, "versioned read must not see the hash key's bytes")

		require.NoError(t, s.Delete(ctx, "hash", "a#b"))
		got, err = readAll(s.GetWithVersion(ctx, "hash", "a", "b"))
		require.NoError(t, err, "deleting the hash key must not delete the versioned object")
		require.Equal(t, round.versionBody, got)

		require.NoError(t, s.DeleteWithVersion(ctx, "hash", "a", "b"))
	}
}

func readAll(r BlobReader, m *BlobMetadata, err error) (string, error) {
	if err != nil {
		return "", err
	}
	defer r.Close()
	data, readErr := io.ReadAll(r)
	if readErr != nil {
		return "", readErr
	}
	return string(data), nil
}

// ListParts returns ascending part-number order regardless of the lexical
// directory order the parts were written in (part.10 precedes part.2
// lexically), so pagination and truncation markers stay numeric.
func TestBlobListPartsAscendingOrder(t *testing.T) {
	tmpDir := t.TempDir()
	bs, err := Open(tmpDir)
	require.NoError(t, err)
	defer bs.Close()
	s, err := NewHybridBlobStore(bs, tmpDir)
	require.NoError(t, err)
	ctx := context.Background()

	uploadID, err := s.CreateMultipartUpload(ctx, "sort", "sorted.txt", nil)
	require.NoError(t, err)
	for partNum := 12; partNum >= 1; partNum-- {
		_, err := s.UploadPart(ctx, "sort", "sorted.txt", uploadID, partNum, strings.NewReader("x"))
		require.NoError(t, err)
	}

	parts, err := s.ListParts(ctx, "sort", "sorted.txt", uploadID)
	require.NoError(t, err)
	require.Len(t, parts, 12)
	for i, p := range parts {
		require.Equal(t, i+1, p.PartNumber, "parts must ascend by part number")
	}
}

// The multipart complete must store the assembled object at the address the
// versioned reads resolve to — including keys whose final segment rewrites.
// A null version addresses the plain copy, which the "null"-version read
// fallback resolves.
func TestMultipartCompleteVersionedAddress(t *testing.T) {
	tmpDir := t.TempDir()
	bs, err := Open(tmpDir)
	require.NoError(t, err)
	defer bs.Close()
	s, err := NewHybridBlobStore(bs, tmpDir)
	require.NoError(t, err)
	ctx := context.Background()

	part := bytes.Repeat([]byte("y"), int(SmallObjectThreshold)+1)

	for _, tc := range []struct {
		key       string
		versionId string
	}{
		{"plain/key.txt", "v9"},
		{"dir/", "v9"},
		{"dir/", ""}, // null version: the plain address
	} {
		uploadID, err := s.CreateMultipartUpload(ctx, "mpu", tc.key, nil)
		require.NoError(t, err)
		_, err = s.UploadPart(ctx, "mpu", tc.key, uploadID, 1, bytes.NewReader(part))
		require.NoError(t, err)

		_, err = s.CompleteMultipartUpload(ctx, "mpu", tc.key, tc.versionId, uploadID, []PartInfo{{PartNumber: 1}})
		require.NoError(t, err, "CompleteMultipartUpload(%q, %q)", tc.key, tc.versionId)

		version := tc.versionId
		if version == "" {
			version = "null"
		}
		reader, meta, err := s.GetWithVersion(ctx, "mpu", tc.key, version)
		require.NoError(t, err, "GetWithVersion(%q, %q)", tc.key, version)
		got, readErr := io.ReadAll(reader)
		require.NoError(t, reader.Close())
		require.NoError(t, readErr)
		require.Len(t, got, len(part))
		require.EqualValues(t, len(part), meta.Size)
	}
}
