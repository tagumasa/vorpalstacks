package s3

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"time"
)

const (
	eventContentTypeHeader      = ":content-type"
	eventMessageTypeHeader      = ":message-type"
	eventEventTypeHeader        = ":event-type"
	eventContentTypeOctetStream = "application/octet-stream"
	eventMessageTypeEvent       = "event"
)

var eventTypeNames = map[string]string{
	"records":      "Records",
	"stats":        "Stats",
	"progress":     "Progress",
	"continuation": "Cont",
	"end":          "End",
}

type eventStreamEncoder struct {
	w io.Writer
}

func newEventStreamEncoder(w io.Writer) *eventStreamEncoder {
	return &eventStreamEncoder{w: w}
}

type eventStreamHeader struct {
	name  string
	value string
}

func (e *eventStreamEncoder) encodeEvent(eventType string, payload []byte, headers []eventStreamHeader) error {
	var headerBuf bytes.Buffer
	for _, h := range headers {
		if err := e.encodeHeader(&headerBuf, h.name, h.value); err != nil {
			return err
		}
	}

	headersLen := headerBuf.Len()
	payloadLen := len(payload)

	totalLen := 4 + 4 + 4 + headersLen + payloadLen + 4

	prelude := make([]byte, 8)
	binary.BigEndian.PutUint32(prelude[0:4], uint32(totalLen))
	binary.BigEndian.PutUint32(prelude[4:8], uint32(headersLen))
	preludeCrc := crc32.ChecksumIEEE(prelude)

	var message bytes.Buffer
	if err := binary.Write(&message, binary.BigEndian, uint32(totalLen)); err != nil {
		return fmt.Errorf("failed to write total length: %w", err)
	}
	if err := binary.Write(&message, binary.BigEndian, uint32(headersLen)); err != nil {
		return fmt.Errorf("failed to write headers length: %w", err)
	}
	if err := binary.Write(&message, binary.BigEndian, preludeCrc); err != nil {
		return fmt.Errorf("failed to write prelude CRC: %w", err)
	}
	message.Write(headerBuf.Bytes())
	message.Write(payload)

	messageCrc := crc32.ChecksumIEEE(message.Bytes())
	if err := binary.Write(&message, binary.BigEndian, messageCrc); err != nil {
		return fmt.Errorf("failed to write message CRC: %w", err)
	}

	_, err := e.w.Write(message.Bytes())
	return err
}

func (e *eventStreamEncoder) encodeHeader(w io.Writer, name, value string) error {
	nameBytes := []byte(name)
	valueBytes := []byte(value)

	if _, err := w.Write([]byte{byte(len(nameBytes))}); err != nil {
		return fmt.Errorf("failed to write header name length: %w", err)
	}
	if _, err := w.Write(nameBytes); err != nil {
		return fmt.Errorf("failed to write header name: %w", err)
	}

	if _, err := w.Write([]byte{0x07}); err != nil {
		return fmt.Errorf("failed to write header type: %w", err)
	}

	if err := binary.Write(w, binary.BigEndian, uint16(len(valueBytes))); err != nil {
		return fmt.Errorf("failed to write header value length: %w", err)
	}
	if _, err := w.Write(valueBytes); err != nil {
		return fmt.Errorf("failed to write header value: %w", err)
	}

	return nil
}

func (e *eventStreamEncoder) writeRecordsEvent(payload []byte) error {
	headers := []eventStreamHeader{
		{eventMessageTypeHeader, eventMessageTypeEvent},
		{eventEventTypeHeader, eventTypeNames["records"]},
		{eventContentTypeHeader, eventContentTypeOctetStream},
	}
	return e.encodeEvent("records", payload, headers)
}

func (e *eventStreamEncoder) writeStatsEvent(stats *Stats) error {
	statsXML := fmt.Sprintf(`<Stats><BytesScanned>%d</BytesScanned><BytesProcessed>%d</BytesProcessed><BytesReturned>%d</BytesReturned></Stats>`,
		stats.BytesScanned, stats.BytesProcessed, stats.BytesReturned)

	headers := []eventStreamHeader{
		{eventMessageTypeHeader, eventMessageTypeEvent},
		{eventEventTypeHeader, eventTypeNames["stats"]},
		{eventContentTypeHeader, "text/xml"},
	}
	return e.encodeEvent("stats", []byte(statsXML), headers)
}

func (e *eventStreamEncoder) writeProgressEvent(progress *Progress) error {
	progressXML := fmt.Sprintf(`<Progress><BytesScanned>%d</BytesScanned><BytesProcessed>%d</BytesProcessed><BytesReturned>%d</BytesReturned></Progress>`,
		progress.BytesScanned, progress.BytesProcessed, progress.BytesReturned)

	headers := []eventStreamHeader{
		{eventMessageTypeHeader, eventMessageTypeEvent},
		{eventEventTypeHeader, eventTypeNames["progress"]},
		{eventContentTypeHeader, "text/xml"},
	}
	return e.encodeEvent("progress", []byte(progressXML), headers)
}

func (e *eventStreamEncoder) writeEndEvent() error {
	headers := []eventStreamHeader{
		{eventMessageTypeHeader, eventMessageTypeEvent},
		{eventEventTypeHeader, eventTypeNames["end"]},
	}
	return e.encodeEvent("end", nil, headers)
}

// EventStreamWriter defines the interface for writing S3 Select event streams.
type EventStreamWriter interface {
	WriteRecords(data []byte) error
	AddBytesScanned(n int64)
	AddBytesProcessed(n int64)
	MaybeWriteProgress() error
	WriteStats() error
	WriteEnd() error
	Flush() error
}

// SelectEventStreamWriter writes S3 Select event streams.
type SelectEventStreamWriter struct {
	encoder      *eventStreamEncoder
	stats        Stats
	flush        func() error
	progress     *RequestProgress
	lastProgress time.Time
}

// NewSelectEventStreamWriter creates a new SelectEventStreamWriter.
func NewSelectEventStreamWriter(w io.Writer, progress *RequestProgress) *SelectEventStreamWriter {
	return &SelectEventStreamWriter{
		encoder:  newEventStreamEncoder(w),
		progress: progress,
	}
}

// SetFlush sets the flush function for the event stream writer.
func (s *SelectEventStreamWriter) SetFlush(flush func() error) {
	s.flush = flush
}

// WriteRecords writes records to the event stream.
func (s *SelectEventStreamWriter) WriteRecords(data []byte) error {
	s.stats.BytesReturned += int64(len(data))
	return s.encoder.writeRecordsEvent(data)
}

// AddBytesScanned adds to the bytes scanned counter.
func (s *SelectEventStreamWriter) AddBytesScanned(n int64) {
	s.stats.BytesScanned += n
}

// AddBytesProcessed adds to the bytes processed counter.
func (s *SelectEventStreamWriter) AddBytesProcessed(n int64) {
	s.stats.BytesProcessed += n
}

// MaybeWriteProgress writes a progress event if enough time has elapsed.
func (s *SelectEventStreamWriter) MaybeWriteProgress() error {
	if s.progress == nil || !s.progress.Enabled {
		return nil
	}

	now := time.Now()
	if now.Sub(s.lastProgress) < time.Second {
		return nil
	}
	s.lastProgress = now

	return s.encoder.writeProgressEvent(&Progress{
		BytesScanned:   s.stats.BytesScanned,
		BytesProcessed: s.stats.BytesProcessed,
		BytesReturned:  s.stats.BytesReturned,
	})
}

// WriteStats writes a stats event to the event stream.
func (s *SelectEventStreamWriter) WriteStats() error {
	return s.encoder.writeStatsEvent(&s.stats)
}

// WriteEnd writes an end event to the event stream.
func (s *SelectEventStreamWriter) WriteEnd() error {
	return s.encoder.writeEndEvent()
}

// Flush flushes any buffered data.
func (s *SelectEventStreamWriter) Flush() error {
	if s.flush != nil {
		return s.flush()
	}
	return nil
}

func marshalJSONRecord(record map[string]interface{}, delimiter string) ([]byte, error) {
	data, err := json.Marshal(record)
	if err != nil {
		return nil, err
	}
	return append(data, []byte(delimiter)...), nil
}

func marshalCSVRecord(fields []string, delimiter, fieldDelimiter string) ([]byte, error) {
	var buf bytes.Buffer
	for i, field := range fields {
		if i > 0 {
			buf.WriteString(fieldDelimiter)
		}
		needsQuote := stringsContainsAny(field, fieldDelimiter+"\n\"")
		if needsQuote {
			buf.WriteByte('"')
			for _, c := range field {
				if c == '"' {
					buf.WriteString(`""`)
				} else {
					buf.WriteRune(c)
				}
			}
			buf.WriteByte('"')
		} else {
			buf.WriteString(field)
		}
	}
	buf.WriteString(delimiter)
	return buf.Bytes(), nil
}

func stringsContainsAny(s, chars string) bool {
	for _, c := range s {
		for _, m := range chars {
			if c == m {
				return true
			}
		}
	}
	return false
}
