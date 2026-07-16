package neptunedata

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	"vorpalstacks/internal/core/storage/graphengine"
)

// GraphSON v3 type namespace constants used by TinkerPop Gremlin Server.
// Each value is serialised as the "@type" field in typed JSON messages.
const (
	graphsonVertex   = "g:Vertex"
	graphsonEdge     = "g:Edge"
	graphsonVProp    = "g:VertexProperty"
	graphsonProperty = "g:Property"
	graphsonPath     = "g:Path"
	graphsonList     = "g:List"
	graphsonSet      = "g:Set"
	graphsonMap      = "g:Map"
	graphsonInt32    = "g:Int32"
	graphsonInt64    = "g:Int64"
	graphsonFloat    = "g:Float"
	graphsonDouble   = "g:Double"
	graphsonDate     = "g:Date"
	graphsonUUID     = "g:UUID"
	graphsonBulkSet  = "g:BulkSet"
)

// EncodeGraphSON3 converts a Go value produced by gremlinparser into a
// GraphSON v3-compatible value. Graph engine types (*graphengine.Node,
// *graphengine.Edge) are wrapped in typed envelopes; primitives and
// collections are mapped recursively.
func EncodeGraphSON3(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case *graphengine.Node:
		return encodeNode(val)
	case *graphengine.Edge:
		return encodeEdge(val)
	case *graphengine.TraversalResult:
		return encodeTraversalResult(val)
	case *graphengine.PathResult:
		return encodePathResult(val)
	case map[string]interface{}:
		return encodeMap(val)
	case []interface{}:
		return encodeSlice(val)
	case string:
		return val
	case bool:
		return val
	case int:
		return float64(val)
	case int32:
		return wrapTyped(graphsonInt32, float64(val))
	case int64:
		return wrapTyped(graphsonInt64, float64(val))
	case uint64:
		if val > uint64(math.MaxInt64) {
			return wrapTyped(graphsonInt64, float64(val))
		}
		return float64(val)
	case float32:
		return wrapTyped(graphsonFloat, float64(val))
	case float64:
		return val
	case time.Time:
		return wrapTyped(graphsonDate, val.Format(time.RFC3339Nano))
	case graphengine.NodeID:
		return float64(val)
	case graphengine.EdgeID:
		return float64(val)
	default:
		return encodeByReflection(v)
	}
}

// EncodeGraphSON3Result is a convenience wrapper that takes the raw result of
// ExecuteQuery and wraps it in a GraphSON v3 typed list envelope suitable for
// the "result.data" field of a TinkerPop response.
func EncodeGraphSON3Result(data interface{}) interface{} {
	if data == nil {
		return wrapTyped(graphsonList, []interface{}{})
	}
	encoded := EncodeGraphSON3(data)
	return wrapTyped(graphsonList, asSlice(encoded))
}

// asSlice ensures the value is represented as a JSON array. A single non-slice
// value is wrapped into a one-element slice.
func asSlice(v interface{}) []interface{} {
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return []interface{}{v}
}

// wrapTyped wraps a value in a GraphSON v3 typed envelope:
//
//	{"@type": typeName, "@value": value}
func wrapTyped(typeName string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"@type":  typeName,
		"@value": value,
	}
}

// encodeNode converts a graph engine Node into a GraphSON v3 Vertex envelope.
// Vertex format: {"@type":"g:Vertex","@value":{"id":ID,"label":"lbl","properties":{...}}}
func encodeNode(n *graphengine.Node) map[string]interface{} {
	props := make(map[string]interface{})
	for k, raw := range n.Props {
		vp := wrapTyped(graphsonVProp, map[string]interface{}{
			"id":    float64(n.ID),
			"label": k,
			"value": EncodeGraphSON3(raw),
		})
		props[k] = []interface{}{vp}
	}
	return wrapTyped(graphsonVertex, map[string]interface{}{
		"id":         float64(n.ID),
		"label":      firstLabel(n.Labels),
		"properties": props,
	})
}

// encodeEdge converts a graph engine Edge into a GraphSON v3 Edge envelope.
// Edge format: {"@type":"g:Edge","@value":{id,label,inVLabel,outVLabel,inV,outV,properties}}
func encodeEdge(e *graphengine.Edge) map[string]interface{} {
	props := make(map[string]interface{})
	for k, v := range e.Props {
		props[k] = wrapTyped(graphsonProperty, map[string]interface{}{
			"key":   k,
			"value": EncodeGraphSON3(v),
		})
	}
	return wrapTyped(graphsonEdge, map[string]interface{}{
		"id":         float64(e.ID),
		"label":      e.Label,
		"inVLabel":   "",
		"outVLabel":  "",
		"inV":        float64(e.To),
		"outV":       float64(e.From),
		"properties": props,
	})
}

// encodeTraversalResult converts a TraversalResult into a GraphSON v3 map.
// TraversalResult carries a Node, Depth, and Path from BFS operations.
func encodeTraversalResult(tr *graphengine.TraversalResult) interface{} {
	pathEdges := make([]interface{}, len(tr.Path))
	for i, e := range tr.Path {
		pathEdges[i] = EncodeGraphSON3(e)
	}
	return wrapTyped(graphsonMap, []interface{}{
		"node", EncodeGraphSON3(tr.Node),
		"depth", EncodeGraphSON3(tr.Depth),
		"path", wrapTyped(graphsonList, pathEdges),
	})
}

// encodePathResult converts a PathResult into a GraphSON v3 Path envelope.
// Path format: {"@type":"g:Path","@value":{"labels":[[...]],"objects":[...]}}
func encodePathResult(pr *graphengine.PathResult) interface{} {
	labels := make([]interface{}, len(pr.Nodes)+1)
	for i := range labels {
		labels[i] = []interface{}{}
	}
	objects := make([]interface{}, 0, len(pr.Nodes)+len(pr.Edges))
	for _, n := range pr.Nodes {
		objects = append(objects, EncodeGraphSON3(n))
	}
	return wrapTyped(graphsonPath, map[string]interface{}{
		"labels":  wrapTyped(graphsonList, labels),
		"objects": wrapTyped(graphsonList, objects),
	})
}

// encodeMap recursively encodes a Go map into a GraphSON v3 Map envelope.
// GraphSON v3 maps are serialised as alternating key-value arrays:
//
//	{"@type":"g:Map","@value":["k1",v1,"k2",v2,...]}
func encodeMap(m map[string]interface{}) interface{} {
	pairs := make([]interface{}, 0, len(m)*2)
	for k, v := range m {
		pairs = append(pairs, k)
		pairs = append(pairs, EncodeGraphSON3(v))
	}
	return wrapTyped(graphsonMap, pairs)
}

// encodeSlice recursively encodes a Go slice into a GraphSON v3 List envelope.
func encodeSlice(s []interface{}) interface{} {
	items := make([]interface{}, len(s))
	for i, v := range s {
		items[i] = EncodeGraphSON3(v)
	}
	return wrapTyped(graphsonList, items)
}

// encodeByReflection handles values not covered by the type switch, using
// reflection to identify underlying types (e.g. named types wrapping int).
func encodeByReflection(v interface{}) interface{} {
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return float64(rv.Int())
	case reflect.Int64:
		return wrapTyped(graphsonInt64, float64(rv.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return float64(rv.Uint())
	case reflect.Uint64:
		return wrapTyped(graphsonInt64, float64(rv.Uint()))
	case reflect.Float32:
		return wrapTyped(graphsonFloat, float64(rv.Float()))
	case reflect.Float64:
		return rv.Float()
	case reflect.Bool:
		return rv.Bool()
	case reflect.String:
		return rv.String()
	case reflect.Slice, reflect.Array:
		items := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			items[i] = EncodeGraphSON3(rv.Index(i).Interface())
		}
		return wrapTyped(graphsonList, items)
	case reflect.Map:
		pairs := make([]interface{}, 0, rv.Len()*2)
		iter := rv.MapRange()
		for iter.Next() {
			pairs = append(pairs, fmt.Sprintf("%v", iter.Key().Interface()))
			pairs = append(pairs, EncodeGraphSON3(iter.Value().Interface()))
		}
		return wrapTyped(graphsonMap, pairs)
	case reflect.Struct:
		if _, ok := v.(fmt.Stringer); ok {
			return fmt.Sprintf("%v", v)
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// firstLabel returns the first label of a node, or "vertex" if none exist.
func firstLabel(labels []string) string {
	if len(labels) > 0 {
		return labels[0]
	}
	return "vertex"
}

// DecodeGraphSON3Args decodes a GraphSON v3 value back into a plain Go value,
// stripping @type/@value envelopes. Used when parsing incoming WebSocket
// message arguments (bindings, etc.).
func DecodeGraphSON3Args(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	m, ok := v.(map[string]interface{})
	if !ok {
		return v
	}

	typeName, hasType := m["@type"]
	raw, hasValue := m["@value"]
	if !hasType || !hasValue {
		return v
	}

	ts, ok := typeName.(string)
	if !ok {
		return v
	}

	switch ts {
	case graphsonList, graphsonSet, graphsonBulkSet:
		if arr, ok := raw.([]interface{}); ok {
			result := make([]interface{}, len(arr))
			for i, elem := range arr {
				result[i] = DecodeGraphSON3Args(elem)
			}
			return result
		}
		return raw

	case graphsonMap:
		if arr, ok := raw.([]interface{}); ok {
			result := make(map[string]interface{})
			for i := 0; i+1 < len(arr); i += 2 {
				key := fmt.Sprintf("%v", DecodeGraphSON3Args(arr[i]))
				result[key] = DecodeGraphSON3Args(arr[i+1])
			}
			return result
		}
		return raw

	case graphsonInt32, graphsonInt64, graphsonFloat, graphsonDouble:
		return raw

	case graphsonDate:
		if s, ok := raw.(string); ok {
			if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
				return t
			}
		}
		return raw

	case graphsonUUID:
		if s, ok := raw.(string); ok {
			return strings.ToLower(s)
		}
		return raw

	case graphsonVertex, graphsonEdge, graphsonVProp, graphsonProperty, graphsonPath:
		return raw

	default:
		return raw
	}
}
