package rdbengine

import (
	"encoding/binary"
	"math"
)

// Key encoding for Pebble row storage. Keys must sort correctly in byte order:
// - Strings/bytes: 0x00 and 0x01 are escaped as [0x01,0x00] and [0x01,0x01]
//   so that the 0x00 separator between composite PK components is unambiguous.
// - Integers: big-endian fixed-width with sign-bit flip (^0x80 on MSB)
//   so signed values sort correctly: -3 < -1 < 0 < 5
//   Without this, negative values (0xFF..) sort after all positives.

// EncodePK encodes a composite primary key from column values.
// Each component is encoded according to its type and concatenated with 0x00 separator.
func EncodePK(schema *TableSchema, row Row) []byte {
	var buf []byte
	first := true
	for _, col := range schema.PrimaryKeyColumns() {
		if !first {
			buf = append(buf, 0x00)
		}
		first = false
		cv, ok := row[col.Name]
		if !ok || cv.Value == nil {
			buf = append(buf, 0x00)
			continue
		}
		buf = append(buf, encodeValue(cv)...)
	}
	return buf
}

// encodeValue encodes a single ColumnValue for use as a Pebble key component.
func encodeValue(cv ColumnValue) []byte {
	switch v := cv.Value.(type) {
	case int:
		return encodeInt64(int64(v))
	case int8:
		return encodeInt64(int64(v))
	case int16:
		return encodeInt64(int64(v))
	case int32:
		return encodeInt32(v)
	case int64:
		return encodeInt64(v)
	case uint:
		return encodeUint64(uint64(v))
	case uint8:
		return encodeUint64(uint64(v))
	case uint16:
		return encodeUint64(uint64(v))
	case uint32:
		return encodeUint64(uint64(v))
	case uint64:
		return encodeUint64(v)
	case float32:
		return encodeFloat64(float64(v))
	case float64:
		return encodeFloat64(v)
	case bool:
		if v {
			return []byte{0x01}
		}
		return []byte{0x00}
	case string:
		return escapeKeyBytes([]byte(v))
	case []byte:
		return escapeKeyBytes(v)
	default:
		return []byte{}
	}
}

// encodeInt32 encodes a signed int32 with sign-bit flip for correct byte-order sorting.
func encodeInt32(v int32) []byte {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v)^0x80000000)
	return buf[:]
}

// encodeInt64 encodes a signed int64 with sign-bit flip for correct byte-order sorting.
func encodeInt64(v int64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(v)^0x8000000000000000)
	return buf[:]
}

// encodeUint64 encodes an unsigned uint64 in big-endian byte order.
func encodeUint64(v uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], v)
	return buf[:]
}

// encodeFloat64 encodes a float64 for correct byte-order sorting.
// Uses IEEE 754 representation with sign-bit manipulation.
func encodeFloat64(v float64) []byte {
	var buf [8]byte
	bits := math.Float64bits(v)
	if v >= 0 {
		bits |= 0x8000000000000000
	} else {
		bits = ^bits
	}
	binary.BigEndian.PutUint64(buf[:], bits)
	return buf[:]
}

// decodeInt32 decodes a sign-bit-flipped int32 from Pebble key bytes.
func decodeInt32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b) ^ 0x80000000)
}

// decodeInt64 decodes a sign-bit-flipped int64 from Pebble key bytes.
func decodeInt64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b) ^ 0x8000000000000000)
}

// decodeUint64 decodes a big-endian uint64 from Pebble key bytes.
func decodeUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// decodeFloat64 decodes a float64 from Pebble key bytes.
func decodeFloat64(b []byte) float64 {
	bits := binary.BigEndian.Uint64(b)
	if bits&0x8000000000000000 != 0 {
		bits &^= 0x8000000000000000
	} else {
		bits = ^bits
	}
	return math.Float64frombits(bits)
}

// escapeKeyBytes escapes 0x00 and 0x01 bytes in string/binary key components so
// that the 0x00 composite-PK separator is unambiguous. The encoding preserves
// lexicographic sort order: escaped("a\x00b") < escaped("a\x01b").
func escapeKeyBytes(b []byte) []byte {
	extra := 0
	for _, c := range b {
		if c <= 0x01 {
			extra++
		}
	}
	if extra == 0 {
		return b
	}
	out := make([]byte, len(b)+extra)
	j := 0
	for _, c := range b {
		if c <= 0x01 {
			out[j] = 0x01
			out[j+1] = c
			j += 2
		} else {
			out[j] = c
			j++
		}
	}
	return out[:j]
}

// rowKeyPrefix returns the key prefix for all rows in a given database table.
// Format: {engine}/{db}/{table}/
func rowKeyPrefix(engine, db, table string) []byte {
	return []byte(engine + "/" + db + "/" + table + "/")
}

// rowKey returns the full key for a specific row.
// Format: {engine}/{db}/{table}/{pk_encoded}
func rowKey(engine, db, table string, pk []byte) []byte {
	prefix := rowKeyPrefix(engine, db, table)
	key := make([]byte, len(prefix)+len(pk))
	copy(key, prefix)
	copy(key[len(prefix):], pk)
	return key
}

// rowEndKey returns the exclusive upper bound for all rows in a table.
func rowEndKey(engine, db, table string) []byte {
	prefix := rowKeyPrefix(engine, db, table)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF
	return end
}

// indexKeyPrefix returns the key prefix for all entries of a secondary index.
// Format: idx/{engine}/{db}/{table}/{idx_name}/
func indexKeyPrefix(engine, db, table, idxName string) []byte {
	return []byte("idx/" + engine + "/" + db + "/" + table + "/" + idxName + "/")
}

// indexKey returns a full key for a secondary index entry.
// Format: idx/{engine}/{db}/{table}/{idx_name}/{col_val_encoded}/{pk_encoded}
func indexKey(engine, db, table, idxName string, colVal []byte, pk []byte) []byte {
	prefix := indexKeyPrefix(engine, db, table, idxName)
	key := make([]byte, len(prefix)+len(colVal)+1+len(pk))
	copy(key, prefix)
	copy(key[len(prefix):], colVal)
	key[len(prefix)+len(colVal)] = 0x00
	copy(key[len(prefix)+len(colVal)+1:], pk)
	return key
}

// indexEndKey returns the exclusive upper bound for an index.
func indexEndKey(engine, db, table, idxName string) []byte {
	prefix := indexKeyPrefix(engine, db, table, idxName)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF
	return end
}

// uniqueKeyPrefix returns the key prefix for unique constraint entries.
// Format: uniq/{engine}/{db}/{table}/{idx_name}/
func uniqueKeyPrefix(engine, db, table, idxName string) []byte {
	return []byte("uniq/" + engine + "/" + db + "/" + table + "/" + idxName + "/")
}

// uniqueKey returns the full key for a unique constraint entry.
// Format: uniq/{engine}/{db}/{table}/{idx_name}/{col_val_encoded}
func uniqueKey(engine, db, table, idxName string, colVal []byte) []byte {
	prefix := uniqueKeyPrefix(engine, db, table, idxName)
	key := make([]byte, len(prefix)+len(colVal))
	copy(key, prefix)
	copy(key[len(prefix):], colVal)
	return key
}

// uniqueEndKey returns the exclusive upper bound for a unique constraint.
func uniqueEndKey(engine, db, table, idxName string) []byte {
	prefix := uniqueKeyPrefix(engine, db, table, idxName)
	end := make([]byte, len(prefix)+1)
	copy(end, prefix)
	end[len(prefix)] = 0xFF
	return end
}

// catalogDBKey returns the key for a database metadata entry.
// Format: catalog/{engine}/db/{db_name}
func catalogDBKey(engine, db string) []byte {
	return []byte("catalog/" + engine + "/db/" + db)
}

// catalogTableKey returns the key for a table schema entry.
// Format: catalog/{engine}/db/{db_name}/table/{tbl_name}
func catalogTableKey(engine, db, table string) []byte {
	return []byte("catalog/" + engine + "/db/" + db + "/table/" + table)
}

// catalogTablePrefix returns the key prefix for all tables in a database.
func catalogTablePrefix(engine, db string) []byte {
	return []byte("catalog/" + engine + "/db/" + db + "/table/")
}

// catalogIndexKey returns the key for an index definition entry.
// Format: catalog/{engine}/db/{db_name}/table/{tbl_name}/idx/{idx_name}
func catalogIndexKey(engine, db, table, idxName string) []byte {
	return []byte("catalog/" + engine + "/db/" + db + "/table/" + table + "/idx/" + idxName)
}

// catalogIndexPrefix returns the key prefix for all indexes on a table.
func catalogIndexPrefix(engine, db, table string) []byte {
	return []byte("catalog/" + engine + "/db/" + db + "/table/" + table + "/idx/")
}

// catalogDBPrefix returns the key prefix for all databases of an engine.
func catalogDBPrefix(engine string) []byte {
	return []byte("catalog/" + engine + "/db/")
}

// autoincKey returns the key for a table's auto-increment counter.
// Format: autoinc/{engine}/{db}/{table}
func autoincKey(engine, db, table string) []byte {
	return []byte("autoinc/" + engine + "/" + db + "/" + table)
}
