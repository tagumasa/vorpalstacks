// Package protocol provides AWS protocol encoding utilities for vorpalstacks.
package protocol

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// XMLElements wraps a slice of items with an XML element name for serialisation.
type XMLElements struct {
	ElementName string
	Items       []interface{}
}

// EncodeXMLResponse writes the response as XML to the HTTP response using standard AWS XML encoding.
func EncodeXMLResponse(w http.ResponseWriter, operationName string, response interface{}) error {
	// Treat nil or an empty map as an empty result.
	if response == nil {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Result></Result>`)); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}
	if m, ok := response.(map[string]interface{}); ok && len(m) == 0 {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Result></Result>`)); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}

	rootName := operationName + "Response"
	resultName := operationName + "Result"

	var xmlBuilder strings.Builder
	xmlBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlBuilder.WriteString("<" + rootName + ">")
	xmlBuilder.WriteString("<" + resultName + ">")
	encodeValueToXML(&xmlBuilder, response)
	xmlBuilder.WriteString("</" + resultName + ">")
	xmlBuilder.WriteString("</" + rootName + ">")

	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xmlBuilder.String())); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

// EncodeRestXMLResponse writes the response as XML using REST-XML encoding (without Result wrapper).
func EncodeRestXMLResponse(w http.ResponseWriter, operationName string, response interface{}) error {
	if response == nil {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Result></Result>`)); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}
	if m, ok := response.(map[string]interface{}); ok && len(m) == 0 {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Result></Result>`)); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}

	rootName := operationName + "Response"

	var xmlBuilder strings.Builder
	xmlBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlBuilder.WriteString("<" + rootName + ">")
	encodeValueToXML(&xmlBuilder, response)
	xmlBuilder.WriteString("</" + rootName + ">")

	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xmlBuilder.String())); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

// EncodeRestXMLResponseWithNamespace writes the response as XML with the specified namespace, using REST-XML encoding.
func EncodeRestXMLResponseWithNamespace(w http.ResponseWriter, operationName string, xmlNamespace string, response interface{}) error {
	if response == nil {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Result xmlns="%s"></Result>`, xmlNamespace))); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}
	if m, ok := response.(map[string]interface{}); ok && len(m) == 0 {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?><Result xmlns="%s"></Result>`, xmlNamespace))); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}

	rootName := operationName + "Response"

	var xmlBuilder strings.Builder
	xmlBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlBuilder.WriteString(fmt.Sprintf(`<%s xmlns="%s">`, rootName, xmlNamespace))
	encodeValueToXML(&xmlBuilder, response)
	xmlBuilder.WriteString("</" + rootName + ">")

	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xmlBuilder.String())); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

// EncodeRestXMLPayloadResponse writes the response as XML using REST-XML payload encoding.
func EncodeRestXMLPayloadResponse(w http.ResponseWriter, payloadRoot string, response interface{}) error {
	if response == nil {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Result></Result>`)); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}
	if m, ok := response.(map[string]interface{}); ok && len(m) == 0 {
		w.Header().Set("Content-Type", "application/xml")
		if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?><Result></Result>`)); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		return nil
	}

	payload := response
	if m, ok := response.(map[string]interface{}); ok {
		if inner, exists := m[payloadRoot]; exists {
			payload = inner
		}
	}

	var xmlBuilder strings.Builder
	xmlBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlBuilder.WriteString("<" + payloadRoot + ">")
	encodeValueToXML(&xmlBuilder, payload)
	xmlBuilder.WriteString("</" + payloadRoot + ">")

	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xmlBuilder.String())); err != nil {
		return fmt.Errorf("write response: %w", err)
	}

	return nil
}

// EncodeQueryXMLResponse writes the response as XML using Query API XML encoding with ResponseMetadata.
func EncodeQueryXMLResponse(w http.ResponseWriter, operationName string, response interface{}) error {
	rootName := operationName + "Response"
	resultName := operationName + "Result"

	var xmlBuilder strings.Builder
	xmlBuilder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlBuilder.WriteString("<" + rootName + ">")
	xmlBuilder.WriteString("<" + resultName + ">")

	if response != nil {
		encodeValueToXML(&xmlBuilder, response)
	}

	xmlBuilder.WriteString("</" + resultName + ">")
	xmlBuilder.WriteString("<ResponseMetadata><RequestId>example-request-id</RequestId></ResponseMetadata>")
	xmlBuilder.WriteString("</" + rootName + ">")

	w.Header().Set("Content-Type", "application/xml")
	if _, err := w.Write([]byte(xmlBuilder.String())); err != nil {
		return fmt.Errorf("write response: %w", err)
	}
	return nil
}

// encodeValueToXML encodes any value to XML. Common types are handled via type
// assertions (zero-allocation fast path); structs fall through to the
// reflection-based path so that implemented handlers returning structs are
// encoded correctly without an intermediate JSON round-trip.
func encodeValueToXML(builder *strings.Builder, v interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		encodeMapToXML(builder, val)
	case map[string]string:
		for k, s := range val {
			builder.WriteString("<entry>")
			builder.WriteString("<key>" + escapeXML(k) + "</key>")
			builder.WriteString("<value>" + escapeXML(s) + "</value>")
			builder.WriteString("</entry>")
		}
	case XMLElements:
		for _, item := range val.Items {
			builder.WriteString("<" + val.ElementName + ">")
			encodeValueToXML(builder, item)
			builder.WriteString("</" + val.ElementName + ">")
		}
	case []interface{}:
		for _, item := range val {
			builder.WriteString("<member>")
			encodeValueToXML(builder, item)
			builder.WriteString("</member>")
		}
	case []map[string]interface{}:
		for _, itemMap := range val {
			builder.WriteString("<member>")
			encodeMapToXML(builder, itemMap)
			builder.WriteString("</member>")
		}
	case []string:
		for _, item := range val {
			builder.WriteString("<member>")
			builder.WriteString(escapeXML(item))
			builder.WriteString("</member>")
		}
	case time.Time:
		builder.WriteString(escapeXML(val.UTC().Format(time.RFC3339)))
	default:
		rv := reflect.ValueOf(v)
		for rv.Kind() == reflect.Ptr {
			if rv.IsNil() {
				return
			}
			rv = rv.Elem()
		}
		switch rv.Kind() {
		case reflect.Struct:
			encodeStructToXML(builder, rv)
		case reflect.Map:
			encodeReflectMapToXML(builder, rv)
		case reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				builder.WriteString("<member>")
				encodeValueToXML(builder, rv.Index(i).Interface())
				builder.WriteString("</member>")
			}
		default:
			builder.WriteString(anyToString(v))
		}
	}
}

// encodeStructToXML serialises a struct to XML using reflection. Field names
// are taken from the json struct tag when present, matching the convention used
// by the rest of the codebase.
func encodeStructToXML(builder *strings.Builder, rv reflect.Value) {
	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}
		fv := rv.Field(i)
		// Dereference pointer fields.
		for fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				break
			}
			fv = fv.Elem()
		}
		if fv.Kind() == reflect.Ptr && fv.IsNil() {
			continue
		}
		key, skip, omitempty := parseJSONTag(field.Tag.Get("json"), field.Name)
		if skip {
			continue
		}
		if omitempty && fv.IsZero() {
			continue
		}
		xmlKey := toTitleCase(key)
		builder.WriteString("<" + xmlKey + ">")
		encodeValueToXML(builder, fv.Interface())
		builder.WriteString("</" + xmlKey + ">")
	}
}

// encodeReflectMapToXML encodes a reflect.Value of kind Map whose key type is
// string (or string-convertible) to XML.
func encodeReflectMapToXML(builder *strings.Builder, rv reflect.Value) {
	for _, k := range rv.MapKeys() {
		var rawKey string
		if k.Kind() == reflect.String {
			rawKey = k.String()
		} else {
			rawKey = anyToString(k.Interface())
		}
		if !isValidXMLName(rawKey) {
			continue
		}
		xmlKey := toTitleCase(rawKey)
		builder.WriteString("<" + xmlKey + ">")
		encodeValueToXML(builder, rv.MapIndex(k).Interface())
		builder.WriteString("</" + xmlKey + ">")
	}
}

// parseJSONTag extracts the XML element name and options from a json struct tag.
func parseJSONTag(tag, fieldName string) (key string, skip bool, omitempty bool) {
	if tag == "" {
		return fieldName, false, false
	}
	parts := strings.SplitN(tag, ",", 2)
	name := parts[0]
	if name == "-" {
		return "", true, false
	}
	if name == "" {
		name = fieldName
	}
	if len(parts) > 1 && strings.Contains(parts[1], "omitempty") {
		omitempty = true
	}
	return name, false, omitempty
}

// encodeMapToXML recursively encodes a map[string]interface{} to XML. The
// indent parameter is kept for future formatted output but is not currently
// used.
func encodeMapToXML(builder *strings.Builder, data map[string]interface{}) {
	for key, value := range data {
		if !isValidXMLName(key) {
			continue
		}
		xmlKey := toTitleCase(key)
		builder.WriteString("<" + xmlKey + ">")

		switch v := value.(type) {
		case map[string]interface{}:
			if len(v) > 0 {
				encodeMapToXML(builder, v)
			}
		case map[string]string:
			for k, val := range v {
				builder.WriteString("<entry>")
				builder.WriteString("<key>" + escapeXML(k) + "</key>")
				builder.WriteString("<value>" + escapeXML(val) + "</value>")
				builder.WriteString("</entry>")
			}
		case XMLElements:
			for _, item := range v.Items {
				builder.WriteString("<" + v.ElementName + ">")
				encodeValueToXML(builder, item)
				builder.WriteString("</" + v.ElementName + ">")
			}
		case []interface{}:
			for _, item := range v {
				builder.WriteString("<member>")
				if itemMap, ok := item.(map[string]interface{}); ok {
					encodeMapToXML(builder, itemMap)
				} else {
					encodeValueToXML(builder, item)
				}
				builder.WriteString("</member>")
			}
		case []map[string]interface{}:
			for _, itemMap := range v {
				builder.WriteString("<member>")
				encodeMapToXML(builder, itemMap)
				builder.WriteString("</member>")
			}
		case []string:
			for _, item := range v {
				builder.WriteString("<member>")
				builder.WriteString(escapeXML(item))
				builder.WriteString("</member>")
			}
		default:
			encodeValueToXML(builder, v)
		}

		builder.WriteString("</" + xmlKey + ">")
	}
}

func anyToString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return escapeXML(val)
	case time.Time:
		return escapeXML(val.UTC().Format(time.RFC3339))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return numberToString(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return ""
	}
}

func escapeXML(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for _, r := range s {
		switch r {
		case '&':
			result.WriteString("&amp;")
		case '<':
			result.WriteString("&lt;")
		case '>':
			result.WriteString("&gt;")
		case '"':
			result.WriteString("&quot;")
		case '\'':
			result.WriteString("&apos;")
		default:
			if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
				result.WriteString("&#x" + strconv.FormatInt(int64(r), 16) + ";")
			} else {
				result.WriteRune(r)
			}
		}
	}
	return result.String()
}

func isValidXMLName(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if r == '&' || r == '<' || r == '>' || r == '"' || r == '\'' || r < 0x20 {
			return false
		}
		if i == 0 {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || r == ':') {
				return false
			}
		}
	}
	return true
}

func toTitleCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = []rune(strings.ToUpper(string(runes[0])))[0]
	return string(runes)
}

func numberToString(v interface{}) string {
	switch val := v.(type) {
	case int:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return ""
	}
}

func intToString(n int64) string {
	return strconv.FormatInt(n, 10)
}
