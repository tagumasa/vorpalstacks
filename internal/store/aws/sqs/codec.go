package sqs

import (
	"bytes"
	"encoding/base64"
	"sort"
)

// CalculateMessageAttributesMD5 computes the MD5 digest of message attributes
// per the AWS SQS specification. Exported for cross-package use.
func CalculateMessageAttributesMD5(attrs map[string]*MessageAttributeValue) string {
	if len(attrs) == 0 {
		return calculateMD5("")
	}

	keys := make([]string, 0, len(attrs))
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		v := attrs[k]
		if v == nil {
			continue
		}

		buf.Write(uint32ToBytes(uint32(len(k))))
		buf.WriteString(k)

		buf.Write(uint32ToBytes(uint32(len(v.DataType))))
		buf.WriteString(v.DataType)

		if v.StringValue != nil {
			buf.WriteByte(1)
			buf.Write(uint32ToBytes(uint32(len(*v.StringValue))))
			buf.WriteString(*v.StringValue)
		} else if v.BinaryValue != nil {
			buf.WriteByte(2)
			buf.Write(uint32ToBytes(uint32(len(v.BinaryValue))))
			buf.Write(v.BinaryValue)
		} else if len(v.StringListValues) > 0 {
			buf.WriteByte(3)
			buf.Write(uint32ToBytes(uint32(len(v.StringListValues))))
			for _, sv := range v.StringListValues {
				buf.Write(uint32ToBytes(uint32(len(sv))))
				buf.WriteString(sv)
			}
		} else if len(v.BinaryListValues) > 0 {
			buf.WriteByte(4)
			buf.Write(uint32ToBytes(uint32(len(v.BinaryListValues))))
			for _, bv := range v.BinaryListValues {
				buf.Write(uint32ToBytes(uint32(len(bv))))
				buf.Write(bv)
			}
		} else {
			buf.WriteByte(1)
			buf.Write(uint32ToBytes(0))
		}
	}

	return calculateMD5(buf.String())
}

func uint32ToBytes(n uint32) []byte {
	return []byte{
		byte(n >> 24),
		byte(n >> 16),
		byte(n >> 8),
		byte(n),
	}
}

// DecodeBinaryValue decodes a base64-encoded string into a byte slice.
func DecodeBinaryValue(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// EncodeBinaryValue encodes a byte slice into a base64-encoded string.
func EncodeBinaryValue(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
