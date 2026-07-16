package chunk

import "encoding/binary"

const (
	// MagicVLog is the magic bytes identifier for VLog chunk files.
	// This is written at the beginning of each chunk file for identification.
	MagicVLog = "VLOG"
	// VersionV1 is the version number for the current chunk file format.
	VersionV1 = 1
	// EncodingNone indicates no compression is applied to chunk data.
	EncodingNone = 0
	// EncodingGzip indicates GZIP compression is applied to chunk data.
	EncodingGzip = 1
	// EncodingZstd indicates Zstandard compression is applied to chunk data.
	EncodingZstd = 2

	// DefaultChunkSize is the default number of entries allowed in a single chunk.
	DefaultChunkSize = 10000
	// MaxChunkSize is the maximum number of entries allowed in a single chunk.
	MaxChunkSize = 100000
	// MaxMessageSize is the maximum size in bytes for a single message entry.
	MaxMessageSize = 1024 * 1024
)

// Header represents the header metadata for a chunk file.
type Header struct {
	Magic      [4]byte
	Version    uint8
	Encoding   uint8
	EntryCount uint32
	MinTs      int64
	MaxTs      int64
}

// HeaderSize returns the size in bytes of a Header.
func HeaderSize() int {
	return binary.Size(Header{})
}

// ValidateHeader checks whether a Header has valid magic bytes,
// version number, and encoding type.
func ValidateHeader(h *Header) bool {
	return string(h.Magic[:]) == MagicVLog &&
		h.Version == VersionV1 &&
		(h.Encoding == EncodingGzip || h.Encoding == EncodingZstd || h.Encoding == EncodingNone)
}

// Encoding represents the compression encoding used for chunk data.
type Encoding int

// Uint8 returns the uint8 value of the encoding type.
func (e Encoding) Uint8() uint8 {
	return uint8(e)
}

// String returns the string representation of encoding.
func (e Encoding) String() string {
	switch e {
	case EncodingGzip:
		return "gzip"
	case EncodingZstd:
		return "zstd"
	case EncodingNone:
		return "none"
	default:
		return "unknown"
	}
}
