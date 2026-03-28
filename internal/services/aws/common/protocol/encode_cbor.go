package protocol

import (
	"net/http"
	"reflect"
	"time"

	"github.com/fxamacker/cbor/v2"
)

var cborEncMode cbor.EncMode

var timestampKeys = map[string]bool{
	"Timestamp":                          true,
	"StateUpdatedTimestamp":              true,
	"StateTransitionedTimestamp":         true,
	"AlarmConfigurationUpdatedTimestamp": true,
	"CreationDate":                       true,
	"LastUpdated":                        true,
}

func init() {
	opts := cbor.EncOptions{
		Time:    cbor.TimeUnixDynamic,
		TimeTag: cbor.EncTagRequired,
	}
	var err error
	cborEncMode, err = opts.EncMode()
	if err != nil {
		panic("failed to initialize CBOR encoder: " + err.Error())
	}
}

func convertInt64Timestamps(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, vv := range val {
			if timestampKeys[k] {
				if ts, ok := vv.(int64); ok {
					val[k] = time.UnixMilli(ts)
					continue
				}
				if ts, ok := vv.(float64); ok {
					val[k] = time.UnixMilli(int64(ts))
					continue
				}
			}
			val[k] = convertInt64Timestamps(vv)
		}
		return val
	case []interface{}:
		for i, vv := range val {
			val[i] = convertInt64Timestamps(vv)
		}
		return val
	default:
		if m, ok := toGenericSlice(v); ok {
			result := make([]interface{}, len(m))
			for i, vv := range m {
				result[i] = convertInt64Timestamps(vv)
			}
			return result
		}
		return v
	}
}

func toGenericSlice(v interface{}) ([]interface{}, bool) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice {
		return nil, false
	}
	if rv.IsNil() {
		return nil, false
	}
	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return result, true
}

func EncodeCBORResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/cbor")
	w.Header().Set("smithy-protocol", "rpc-v2-cbor")
	sanitized := sanitizeFloats(response)
	converted := convertInt64Timestamps(sanitized)
	data, err := cborEncMode.Marshal(converted)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
