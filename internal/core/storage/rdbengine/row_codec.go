package rdbengine

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// encodeRow serialises a Row into Pebble value bytes using a compact binary
// format: uint32 column count, then for each column: uint16 name length, name
// bytes, uint8 type tag, uint32 value length, value bytes.
func encodeRow(row Row) ([]byte, error) {
	cols := make([]string, 0, len(row))
	for k := range row {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	estimated := 4 + len(cols)*20
	buf := make([]byte, 0, estimated)
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(cols)))

	for _, name := range cols {
		cv := row[name]
		nb := []byte(name)
		buf = binary.BigEndian.AppendUint16(buf, uint16(len(nb)))
		buf = append(buf, nb...)
		buf = append(buf, byte(cv.Type))

		vb, err := encodeColumnValue(cv)
		if err != nil {
			return nil, err
		}
		buf = binary.BigEndian.AppendUint32(buf, uint32(len(vb)))
		buf = append(buf, vb...)
	}
	return buf, nil
}

// decodeRow deserialises a Row from Pebble value bytes.
func decodeRow(data []byte) (Row, error) {
	if len(data) < 4 {
		return nil, fmt.Errorf("rdbengine: row data too short (%d bytes)", len(data))
	}
	colCount := int(binary.BigEndian.Uint32(data[0:4]))
	off := 4
	row := make(Row, colCount)

	for i := 0; i < colCount; i++ {
		if off+2 > len(data) {
			return nil, fmt.Errorf("rdbengine: truncated column name at column %d", i)
		}
		nameLen := int(binary.BigEndian.Uint16(data[off : off+2]))
		off += 2

		if off+nameLen > len(data) {
			return nil, fmt.Errorf("rdbengine: truncated column name data at column %d", i)
		}
		name := string(data[off : off+nameLen])
		off += nameLen

		if off+1 > len(data) {
			return nil, fmt.Errorf("rdbengine: missing type tag at column %d", i)
		}
		colType := ColumnType(data[off])
		off++

		if off+4 > len(data) {
			return nil, fmt.Errorf("rdbengine: truncated value length at column %d", i)
		}
		valLen := int(binary.BigEndian.Uint32(data[off : off+4]))
		off += 4

		if off+valLen > len(data) {
			return nil, fmt.Errorf("rdbengine: truncated value data at column %d", i)
		}
		valData := data[off : off+valLen]
		off += valLen

		val, err := decodeColumnValue(colType, valData)
		if err != nil {
			return nil, fmt.Errorf("rdbengine: column %q: %w", name, err)
		}
		row[name] = ColumnValue{Type: colType, Value: val}
	}
	return row, nil
}

// encodeColumnValue serialises a single typed column value.
func encodeColumnValue(cv ColumnValue) ([]byte, error) {
	if cv.Value == nil {
		return nil, nil
	}
	switch v := cv.Value.(type) {
	case bool:
		if v {
			return []byte{0x01}, nil
		}
		return []byte{0x00}, nil
	case int8:
		return encodeInt64(int64(v)), nil
	case int16:
		return encodeInt64(int64(v)), nil
	case int32:
		return encodeInt64(int64(v)), nil
	case int64:
		return encodeInt64(v), nil
	case int:
		return encodeInt64(int64(v)), nil
	case uint8:
		return encodeUint64(uint64(v)), nil
	case uint16:
		return encodeUint64(uint64(v)), nil
	case uint32:
		return encodeUint64(uint64(v)), nil
	case uint64:
		return encodeUint64(v), nil
	case uint:
		return encodeUint64(uint64(v)), nil
	case float32:
		return encodeFloat64(float64(v)), nil
	case float64:
		return encodeFloat64(v), nil
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	case time.Time:
		return v.MarshalBinary()
	case json.RawMessage:
		return v, nil
	default:
		return json.Marshal(v)
	}
}

// decodeColumnValue deserialises a single typed column value.
func decodeColumnValue(colType ColumnType, data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, nil
	}
	switch colType {
	case ColumnTypeBool:
		return data[0] == 0x01, nil
	case ColumnTypeInt8, ColumnTypeInt16, ColumnTypeInt32, ColumnTypeInt64,
		ColumnTypeSerial, ColumnTypeBigSerial:
		if len(data) < 8 {
			return nil, fmt.Errorf("integer value too short: %d bytes", len(data))
		}
		return decodeInt64(data[:8]), nil
	case ColumnTypeFloat32, ColumnTypeFloat64:
		if len(data) < 8 {
			return nil, fmt.Errorf("float value too short: %d bytes", len(data))
		}
		return decodeFloat64(data[:8]), nil
	case ColumnTypeString, ColumnTypeText:
		return string(data), nil
	case ColumnTypeBytes, ColumnTypeBlob:
		cp := make([]byte, len(data))
		copy(cp, data)
		return cp, nil
	case ColumnTypeDate, ColumnTypeTimestamp, ColumnTypeTimestampTZ:
		var t time.Time
		if err := t.UnmarshalBinary(data); err != nil {
			return nil, err
		}
		return t, nil
	case ColumnTypeJSON:
		return json.RawMessage(data), nil
	case ColumnTypeDecimal:
		return string(data), nil
	case ColumnTypeUUID:
		return string(data), nil
	default:
		return string(data), nil
	}
}
