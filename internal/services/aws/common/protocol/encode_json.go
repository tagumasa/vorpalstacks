package protocol

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"vorpalstacks/internal/utils/buffer"
)

type nanSentinel struct{}

var nanMarker = nanSentinel{}

func sanitizeFloats(v interface{}) interface{} {
	switch val := v.(type) {
	case float64:
		if math.IsNaN(val) || math.IsInf(val, 0) {
			return nanMarker
		}
		return val
	case map[string]interface{}:
		for k, vv := range val {
			sanitized := sanitizeFloats(vv)
			if _, isNan := sanitized.(nanSentinel); isNan {
				delete(val, k)
			} else {
				val[k] = sanitized
			}
		}
		return val
	case []interface{}:
		for i, vv := range val {
			val[i] = sanitizeFloats(vv)
		}
		return val
	default:
		return v
	}
}

// EncodeJSONResponse writes the response as JSON to the HTTP response with appropriate headers.
func EncodeJSONResponse(w http.ResponseWriter, response interface{}) error {
	sanitized := sanitizeFloats(response)

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(sanitized); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(buf.Bytes())
	return err
}

// EncodeJSONResponseWithStatus writes the response as JSON with the specified HTTP status code.
func EncodeJSONResponseWithStatus(w http.ResponseWriter, statusCode int, response interface{}) error {
	sanitized := sanitizeFloats(response)

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(sanitized); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(statusCode)
	_, err := w.Write(buf.Bytes())
	return err
}
