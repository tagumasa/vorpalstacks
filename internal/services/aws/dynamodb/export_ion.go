package dynamodb

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// ---------------------------------------------------------------------------
// Ion export serialisation — the ExportFormat=ION data path.
//
// The developer guide ("DynamoDB table export output format", Amazon Ion)
// documents the contract: the output is Ion text; the DynamoDB datatypes
// map to Ion datatypes (String → string, Boolean → bool, Number → decimal,
// Binary → blob, List → list, Map → struct) with the set types
// disambiguated by Ion type annotations ($dynamodb_SS, $dynamodb_NS,
// $dynamodb_BS); "Items in an Ion export are delimited by newlines. Each
// line begins with an Ion version marker, followed by an item in Ion
// format."
// ---------------------------------------------------------------------------

// buildIonItemLine renders one item as its Ion export line: the Ion
// version marker, then the item wrapped in the Item struct the format's
// example shows. The newline delimiter is the writer's, shared with the
// JSON-lines path.
func buildIonItemLine(item *dbstore.Item) []byte {
	var b strings.Builder
	b.WriteString("$ion_1_0 ")
	b.WriteString(`{"Item":`)
	writeIonAttributes(&b, item.Attributes)
	b.WriteString("}")
	return []byte(b.String())
}

// buildIonIncrementalRecordLine renders one incremental export record as
// its Ion export line. The members mirror the DynamoDB JSON record's —
// Metadata, Keys, and the images the view type selected — with the
// attribute maps rendered through the shared Ion attribute writer.
func buildIonIncrementalRecordLine(record *incrementalExportRecord) []byte {
	var b strings.Builder
	b.WriteString("$ion_1_0 ")
	b.WriteString("{")
	writeIonString(&b, "Metadata")
	b.WriteString(":{")
	writeIonString(&b, "WriteTimestampMicros")
	b.WriteString(":{")
	writeIonString(&b, "N")
	b.WriteString(":")
	b.WriteString(ionDecimal(strconv.FormatInt(record.WriteTimestampMicros, 10)))
	b.WriteString("}}")
	b.WriteString(",")
	writeIonString(&b, "Keys")
	b.WriteString(":")
	writeIonAttributes(&b, record.Keys)
	if record.NewImage != nil {
		b.WriteString(",")
		writeIonString(&b, "NewImage")
		b.WriteString(":")
		writeIonAttributes(&b, record.NewImage)
	}
	if record.OldImage != nil {
		b.WriteString(",")
		writeIonString(&b, "OldImage")
		b.WriteString(":")
		writeIonAttributes(&b, record.OldImage)
	}
	b.WriteString("}")
	return []byte(b.String())
}

// writeIonAttributes writes an attribute map as an Ion struct. Field names
// are quoted strings (legal Ion text field names, and the only safe form
// for arbitrary DynamoDB attribute names); the names are sorted so the
// output is deterministic, matching the key order the JSON path's
// encoding/json marshal produces.
func writeIonAttributes(b *strings.Builder, attrs map[string]*dbstore.AttributeValue) {
	names := make([]string, 0, len(attrs))
	for name := range attrs {
		names = append(names, name)
	}
	sort.Strings(names)
	b.WriteString("{")
	for i, name := range names {
		if i > 0 {
			b.WriteString(",")
		}
		writeIonString(b, name)
		b.WriteString(":")
		writeIonAttributeValue(b, attrs[name])
	}
	b.WriteString("}")
}

// writeIonAttributeValue writes one DynamoDB attribute value in its
// documented Ion mapping. An empty value (every member nil) writes the Ion
// null, the same value the NULL type maps to.
func writeIonAttributeValue(b *strings.Builder, v *dbstore.AttributeValue) {
	switch {
	case v.S != nil:
		writeIonString(b, *v.S)
	case v.N != nil:
		b.WriteString(ionDecimal(*v.N))
	case v.B != nil:
		b.WriteString("{{")
		b.WriteString(base64.StdEncoding.EncodeToString(v.B))
		b.WriteString("}}")
	case v.BOOL != nil:
		if *v.BOOL {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case v.SS != nil:
		b.WriteString("$dynamodb_SS::[")
		for i, s := range v.SS {
			if i > 0 {
				b.WriteString(",")
			}
			writeIonString(b, s)
		}
		b.WriteString("]")
	case v.NS != nil:
		b.WriteString("$dynamodb_NS::[")
		for i, n := range v.NS {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(ionDecimal(n))
		}
		b.WriteString("]")
	case v.BS != nil:
		b.WriteString("$dynamodb_BS::[")
		for i, bin := range v.BS {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString("{{")
			b.WriteString(base64.StdEncoding.EncodeToString(bin))
			b.WriteString("}}")
		}
		b.WriteString("]")
	case v.L != nil:
		b.WriteString("[")
		for i, e := range v.L {
			if i > 0 {
				b.WriteString(",")
			}
			writeIonAttributeValue(b, e)
		}
		b.WriteString("]")
	case v.M != nil:
		writeIonAttributes(b, v.M)
	default:
		b.WriteString("null")
	}
}

// writeIonString writes a quoted Ion text string. Quote, backslash and the
// control characters take the escapes Ion shares with JSON; UTF-8 content
// passes through unescaped (Ion text is UTF-8).
func writeIonString(b *strings.Builder, s string) {
	b.WriteString(`"`)
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(b, `\u%04X`, r)
			} else {
				b.WriteRune(r)
			}
		}
	}
	b.WriteString(`"`)
}

// ionDecimal renders a DynamoDB number string as an Ion decimal literal:
// a bare digit run would parse as an Ion int, so an integer-shaped number
// gains the trailing dot, and an exponent marker becomes the Ion decimal
// exponent (d), which likewise denotes a decimal rather than an int.
func ionDecimal(n string) string {
	if i := strings.IndexAny(n, "eE"); i >= 0 {
		return n[:i] + "d" + n[i+1:]
	}
	if strings.Contains(n, ".") {
		return n
	}
	return n + "."
}
