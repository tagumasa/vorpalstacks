package s3

import (
	"encoding/xml"
	"io"
)

// ExpressionType defines the type of expression used in S3 Select.
type ExpressionType string

// Expression type constants.
const (
	// ExpressionTypeSQL indicates a SQL expression.
	ExpressionTypeSQL ExpressionType = "SQL"
)

// FileHeaderInfo defines how header information is handled in CSV files.
type FileHeaderInfo string

// File header info constants.
const (
	// FileHeaderInfoUse indicates the header should be used.
	FileHeaderInfoUse FileHeaderInfo = "USE"
	// FileHeaderInfoIgnore indicates the header should be ignored.
	FileHeaderInfoIgnore FileHeaderInfo = "IGNORE"
	// FileHeaderInfoNone indicates there is no header.
	FileHeaderInfoNone FileHeaderInfo = "NONE"
)

// QuoteFields defines how fields are quoted in CSV output.
type QuoteFields string

// Quote fields constants.
const (
	// QuoteFieldsAlways indicates all fields should be quoted.
	QuoteFieldsAlways QuoteFields = "ALWAYS"
	// QuoteFieldsAsNeeded indicates fields should be quoted only when necessary.
	QuoteFieldsAsNeeded QuoteFields = "ASNEEDED"
)

// CompressionType defines the compression type for S3 Select.
type CompressionType string

// Compression type constants.
const (
	// CompressionTypeNone indicates no compression.
	CompressionTypeNone CompressionType = "NONE"
	// CompressionTypeGZIP indicates GZIP compression.
	CompressionTypeGZIP CompressionType = "GZIP"
	// CompressionTypeBZIP2 indicates BZIP2 compression.
	CompressionTypeBZIP2 CompressionType = "BZIP2"
)

// InputSerialization defines the input serialization format for S3 Select.
type InputSerialization struct {
	CSV             *CSVInput       `xml:"CSV,omitempty"`
	JSON            *JSONInput      `xml:"JSON,omitempty"`
	Parquet         *ParquetInput   `xml:"Parquet,omitempty"`
	CompressionType CompressionType `xml:"CompressionType,omitempty"`
}

// CSVInput defines CSV input format settings.
type CSVInput struct {
	FileHeaderInfo             FileHeaderInfo `xml:"FileHeaderInfo,omitempty"`
	Comments                   string         `xml:"Comments,omitempty"`
	QuoteEscapeCharacter       string         `xml:"QuoteEscapeCharacter,omitempty"`
	RecordDelimiter            string         `xml:"RecordDelimiter,omitempty"`
	FieldDelimiter             string         `xml:"FieldDelimiter,omitempty"`
	QuoteCharacter             string         `xml:"QuoteCharacter,omitempty"`
	AllowQuotedRecordDelimiter bool           `xml:"AllowQuotedRecordDelimiter,omitempty"`
}

// JSONInput defines JSON input format settings.
type JSONInput struct {
	Type string `xml:"Type,omitempty"`
}

// ParquetInput defines Parquet input format settings.
type ParquetInput struct{}

// OutputSerialization defines the output serialization format for S3 Select.
type OutputSerialization struct {
	CSV  *CSVOutput  `xml:"CSV,omitempty"`
	JSON *JSONOutput `xml:"JSON,omitempty"`
}

// CSVOutput defines CSV output format settings.
type CSVOutput struct {
	QuoteFields          QuoteFields `xml:"QuoteFields,omitempty"`
	QuoteEscapeCharacter string      `xml:"QuoteEscapeCharacter,omitempty"`
	RecordDelimiter      string      `xml:"RecordDelimiter,omitempty"`
	FieldDelimiter       string      `xml:"FieldDelimiter,omitempty"`
	QuoteCharacter       string      `xml:"QuoteCharacter,omitempty"`
}

// JSONOutput defines JSON output format settings.
type JSONOutput struct {
	RecordDelimiter string `xml:"RecordDelimiter,omitempty"`
}

// RequestProgress defines request progress settings.
type RequestProgress struct {
	Enabled bool `xml:"Enabled,omitempty"`
}

// ScanRange defines the byte range for scanning.
type ScanRange struct {
	Start int64 `xml:"Start,omitempty"`
	End   int64 `xml:"End,omitempty"`
}

// SelectObjectContentInput represents the input for SelectObjectContent operation.
type SelectObjectContentInput struct {
	XMLName             xml.Name             `xml:"SelectObjectContentRequest"`
	Expression          string               `xml:"Expression"`
	ExpressionType      ExpressionType       `xml:"ExpressionType"`
	InputSerialization  *InputSerialization  `xml:"InputSerialization"`
	OutputSerialization *OutputSerialization `xml:"OutputSerialization"`
	RequestProgress     *RequestProgress     `xml:"RequestProgress,omitempty"`
	ScanRange           *ScanRange           `xml:"ScanRange,omitempty"`

	Bucket               string `xml:"-"`
	Key                  string `xml:"-"`
	SSECustomerAlgorithm string `xml:"-"`
	SSECustomerKey       string `xml:"-"`
	SSECustomerKeyMD5    string `xml:"-"`
}

// SelectObjectContentOutput represents the output from SelectObjectContent operation.
type SelectObjectContentOutput struct {
	Payload io.ReadCloser `xml:"-"`
}

// Stats holds statistics about the S3 Select query execution.
type Stats struct {
	BytesScanned   int64 `xml:"BytesScanned"`
	BytesProcessed int64 `xml:"BytesProcessed"`
	BytesReturned  int64 `xml:"BytesReturned"`
}

// Progress holds progress information about S3 Select query execution.
type Progress struct {
	BytesScanned   int64 `xml:"BytesScanned"`
	BytesProcessed int64 `xml:"BytesProcessed"`
	BytesReturned  int64 `xml:"BytesReturned"`
}

// SelectObjectContentEventStream represents the event stream for SelectObjectContent.
type SelectObjectContentEventStream struct {
	Records  *RecordsEvent      `xml:"Records,omitempty"`
	Stats    *StatsEvent        `xml:"Stats,omitempty"`
	Progress *ProgressEvent     `xml:"Progress,omitempty"`
	Cont     *ContinuationEvent `xml:"Cont,omitempty"`
	End      *EndEvent          `xml:"End,omitempty"`
}

// RecordsEvent contains the data records returned by S3 Select.
type RecordsEvent struct {
	Payload []byte `xml:",innerxml"`
}

// StatsEvent contains statistics about the query execution.
type StatsEvent struct {
	Details Stats `xml:"Stats"`
}

// ProgressEvent contains progress information during query execution.
type ProgressEvent struct {
	Details Progress `xml:"Progress"`
}

// ContinuationEvent indicates that more data is available.
type ContinuationEvent struct{}

// EndEvent indicates the end of the S3 Select query response.
type EndEvent struct{}
