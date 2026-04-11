package protocol

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

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
	case []float64:
		result := make([]interface{}, len(val))
		for i, f := range val {
			if math.IsNaN(f) || math.IsInf(f, 0) {
				result[i] = nanMarker
			} else {
				result[i] = f
			}
		}
		return result
	case []interface{}:
		for i, vv := range val {
			val[i] = sanitizeFloats(vv)
		}
		return val
	default:
		return v
	}
}

// EncodeJSONResponse serialises the response as JSON and writes it to the HTTP response writer with status 200.
func EncodeJSONResponse(w http.ResponseWriter, response interface{}) error {
	sanitized := sanitizeFloats(response)
	converted := ConvertTimestampsToSeconds(sanitized)

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(converted); err != nil {
		return err
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(buf.Bytes())
	return err
}

// EncodeJSONResponseWithStatus serialises the response as JSON and writes it with the specified HTTP status code.
func EncodeJSONResponseWithStatus(w http.ResponseWriter, statusCode int, response interface{}) error {
	sanitized := sanitizeFloats(response)
	converted := ConvertTimestampsToSeconds(sanitized)

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(converted); err != nil {
		return err
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(statusCode)
	_, err := w.Write(buf.Bytes())
	return err
}

var jsonTimestampKeys = map[string]bool{
	"Timestamp":                          true,
	"StateUpdatedTimestamp":              true,
	"StateTransitionedTimestamp":         true,
	"AlarmConfigurationUpdatedTimestamp": true,
	"CreationDate":                       true,
	"LastUpdated":                        true,
}

// ConvertTimestampsToSeconds converts millisecond-precision timestamp values to seconds for JSON serialisation.
// It also converts time.Time and []time.Time values to epoch seconds for services like CloudWatch
// that include Timestamps as []time.Time in JSON responses.
func ConvertTimestampsToSeconds(v interface{}) interface{} {
	switch val := v.(type) {
	case time.Time:
		return float64(val.UnixMilli()) / 1000
	case []time.Time:
		result := make([]interface{}, len(val))
		for i, t := range val {
			result[i] = float64(t.UnixMilli()) / 1000
		}
		return result
	case map[string]interface{}:
		for k, vv := range val {
			if jsonTimestampKeys[k] {
				if ts, ok := vv.(int64); ok {
					if ts > 1e12 {
						val[k] = float64(ts) / 1000
					} else {
						val[k] = float64(ts)
					}
					continue
				}
				if ts, ok := vv.(float64); ok && ts > 1e12 {
					val[k] = ts / 1000
					continue
				}
			}
			val[k] = ConvertTimestampsToSeconds(vv)
		}
		return val
	case []interface{}:
		for i, vv := range val {
			val[i] = ConvertTimestampsToSeconds(vv)
		}
		return val
	default:
		if m, ok := toGenericSlice(v); ok {
			result := make([]interface{}, len(m))
			for i, vv := range m {
				result[i] = ConvertTimestampsToSeconds(vv)
			}
			return result
		}
		return v
	}
}
