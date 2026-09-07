package dynamodb

import (
	"encoding/base64"
)

// This file is the single wire codec for AttributeValue: the conversion
// between the typed attribute values the store works with and the JSON wire
// shape ({"S": …}, {"N": …}, base64 "B", …) the API responses and the
// DynamoDB Streams event format carry. The service layer's response
// builders delegate here, so the wire shape has exactly one definition.

// BuildAttributeValueWire renders one typed attribute value in the wire
// shape. A nil value renders as nil.
func BuildAttributeValueWire(av *AttributeValue) map[string]interface{} {
	if av == nil {
		return nil
	}

	result := make(map[string]interface{})

	if av.S != nil {
		result["S"] = *av.S
	}
	if av.N != nil {
		result["N"] = *av.N
	}
	if av.B != nil {
		result["B"] = base64.StdEncoding.EncodeToString(av.B)
	}
	if av.BOOL != nil {
		result["BOOL"] = *av.BOOL
	}
	if av.NULL != nil && *av.NULL {
		result["NULL"] = true
	}
	if av.SS != nil {
		result["SS"] = av.SS
	}
	if av.NS != nil {
		result["NS"] = av.NS
	}
	if av.BS != nil {
		var encodedBS []string
		for _, b := range av.BS {
			encodedBS = append(encodedBS, base64.StdEncoding.EncodeToString(b))
		}
		result["BS"] = encodedBS
	}
	if av.M != nil {
		m := make(map[string]interface{})
		for k, v := range av.M {
			m[k] = BuildAttributeValueWire(v)
		}
		result["M"] = m
	}
	if av.L != nil {
		var l []interface{}
		for _, item := range av.L {
			l = append(l, BuildAttributeValueWire(item))
		}
		result["L"] = l
	}

	return result
}

// BuildItemWire renders a typed attribute map in the wire shape. A nil map
// renders as an empty map, the response convention for items.
func BuildItemWire(attrs map[string]*AttributeValue) map[string]interface{} {
	if attrs == nil {
		return map[string]interface{}{}
	}

	result := make(map[string]interface{}, len(attrs))
	for k, v := range attrs {
		result[k] = BuildAttributeValueWire(v)
	}
	return result
}

// ParseAttributeValueWire decodes one wire-shaped attribute value back into
// its typed form; values BuildAttributeValueWire produced round-trip
// exactly. Unknown shapes decode as nil.
func ParseAttributeValueWire(m map[string]interface{}) *AttributeValue {
	if m == nil {
		return nil
	}
	av := &AttributeValue{}
	if v, ok := m["S"].(string); ok {
		av.S = &v
	}
	if v, ok := m["N"].(string); ok {
		av.N = &v
	}
	if v, ok := m["B"].(string); ok {
		if raw, err := base64.StdEncoding.DecodeString(v); err == nil {
			av.B = raw
		}
	}
	if v, ok := m["BOOL"].(bool); ok {
		av.BOOL = &v
	}
	if _, ok := m["NULL"]; ok {
		flag := true
		av.NULL = &flag
	}
	if v, ok := m["SS"].([]string); ok {
		av.SS = v
	}
	if v, ok := m["NS"].([]string); ok {
		av.NS = v
	}
	if v, ok := m["BS"].([]interface{}); ok {
		set := make([][]byte, 0, len(v))
		for _, entry := range v {
			if encoded, isStr := entry.(string); isStr {
				if raw, err := base64.StdEncoding.DecodeString(encoded); err == nil {
					set = append(set, raw)
				}
			}
		}
		av.BS = set
	} else if v, ok := m["BS"].([]string); ok {
		set := make([][]byte, 0, len(v))
		for _, encoded := range v {
			if raw, err := base64.StdEncoding.DecodeString(encoded); err == nil {
				set = append(set, raw)
			}
		}
		av.BS = set
	}
	if v, ok := m["M"].(map[string]interface{}); ok {
		nested := make(map[string]*AttributeValue, len(v))
		for k, entry := range v {
			if entryMap, isMap := entry.(map[string]interface{}); isMap {
				nested[k] = ParseAttributeValueWire(entryMap)
			}
		}
		av.M = nested
	}
	if v, ok := m["L"].([]interface{}); ok {
		list := make([]*AttributeValue, 0, len(v))
		for _, entry := range v {
			if entryMap, isMap := entry.(map[string]interface{}); isMap {
				list = append(list, ParseAttributeValueWire(entryMap))
			}
		}
		av.L = list
	}
	return av
}

// ParseItemWire decodes a wire-shaped attribute map back into typed
// attribute values. A nil map decodes as nil.
func ParseItemWire(m map[string]interface{}) map[string]*AttributeValue {
	if m == nil {
		return nil
	}
	result := make(map[string]*AttributeValue, len(m))
	for k, v := range m {
		if entryMap, ok := v.(map[string]interface{}); ok {
			result[k] = ParseAttributeValueWire(entryMap)
		}
	}
	return result
}

// BuildStreamImages renders the Keys member and the stream-view-filtered
// NewImage/OldImage members of a stream record in the wire shape. It is the
// single image builder: the service's stream capture and the TTL worker's
// expiry records both render through it, so a service-initiated record is
// byte-identical to a client-driven one.
func BuildStreamImages(streamViewType StreamViewType, keys, newImage, oldImage map[string]*AttributeValue) (keysResp, newImageResp, oldImageResp map[string]interface{}) {
	keysResp = BuildItemWire(keys)

	switch streamViewType {
	case StreamViewTypeNewImage:
		if newImage != nil {
			newImageResp = BuildItemWire(newImage)
		}
	case StreamViewTypeOldImage:
		if oldImage != nil {
			oldImageResp = BuildItemWire(oldImage)
		}
	case StreamViewTypeNewAndOldImages:
		if newImage != nil {
			newImageResp = BuildItemWire(newImage)
		}
		if oldImage != nil {
			oldImageResp = BuildItemWire(oldImage)
		}
	case StreamViewTypeKeysOnly:
		// Only keys are included.
	}

	return keysResp, newImageResp, oldImageResp
}
