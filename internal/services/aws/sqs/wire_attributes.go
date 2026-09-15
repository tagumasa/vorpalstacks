package sqs

import (
	"strconv"
	"strings"

	"vorpalstacks/internal/common/request"
	sqsstore "vorpalstacks/internal/store/aws/sqs"
)

// ---------------------------------------------------------------------------
// Wire→store message-attribute parsing — the single implementation serving
// both wire formats and every former copy site (single send, both send-batch
// arms, single and batch system attributes). The parse rules live here once:
//
//   - empty-value rule: "Name, type, value and the message body must not be
//     empty or null" (MessageAttributeValue documentation trait,
//     sqs-2012-11-05 model) — an attribute carrying no non-empty StringValue
//     and no non-empty BinaryValue is rejected with InvalidParameterValue;
//   - BinaryValue rule: a present non-empty BinaryValue is base64-decoded on
//     both wire formats, and a decode failure is InvalidParameterValue;
//   - malformed-attribute policy: an attribute whose DataType member is
//     missing is carried with an empty DataType so the shared attribute
//     validator owns that rejection, and a slot with value members but no
//     Name, or a gap with indexed entries beyond it, rejects the request
//     outright — parsing never silently truncates the remaining attribute
//     list.
// ---------------------------------------------------------------------------

// parseRequestMessageAttributes parses the MessageAttributes member of a
// single-send request: the JSON map when present, otherwise the flattened
// query keys (MessageAttribute.N.Name / MessageAttribute.N.Value.*).
func parseRequestMessageAttributes(params map[string]interface{}) (map[string]*sqsstore.MessageAttributeValue, error) {
	if jsonAttrs, ok := params["MessageAttributes"].(map[string]interface{}); ok && len(jsonAttrs) > 0 {
		return parseJSONMessageAttributes(jsonAttrs)
	}
	return parseQueryMessageAttributes(params, "MessageAttribute.")
}

// parseJSONMessageAttributes converts a JSON-protocol attribute map
// (attribute name → value members) into the store representation. A value
// that is not an object violates the MessageAttributeValue shape and
// rejects the request — parsing never silently drops a member of the
// attribute map.
func parseJSONMessageAttributes(jsonAttrs map[string]interface{}) (map[string]*sqsstore.MessageAttributeValue, error) {
	result := make(map[string]*sqsstore.MessageAttributeValue, len(jsonAttrs))
	for name, val := range jsonAttrs {
		attrMap, ok := val.(map[string]interface{})
		if !ok {
			return nil, ErrSerializationException
		}
		attr, err := parseJSONMessageAttributeValue(attrMap)
		if err != nil {
			return nil, err
		}
		result[name] = attr
	}
	return result, nil
}

// parseJSONMessageAttributeValue converts one JSON-protocol attribute value
// into the store representation.
func parseJSONMessageAttributeValue(attrMap map[string]interface{}) (*sqsstore.MessageAttributeValue, error) {
	sv, _ := attrMap["StringValue"].(string)
	bv, _ := attrMap["BinaryValue"].(string)
	return buildMessageAttributeValue(
		request.GetStringParam(attrMap, "DataType"), sv, bv)
}

// parseQueryMessageAttributes reads a flattened-query attribute list into
// the store representation: prefix+N+".Name" carries the attribute name and
// prefix+N+".Value.*" its members. The prefix is the indexed list's stem —
// "MessageAttribute." / "MessageSystemAttribute." for single sends, or the
// batch-entry-qualified variant.
func parseQueryMessageAttributes(params map[string]interface{}, prefix string) (map[string]*sqsstore.MessageAttributeValue, error) {
	result := make(map[string]*sqsstore.MessageAttributeValue)
	for i := 1; ; i++ {
		name := request.GetParamCaseInsensitive(params, prefix+strconv.Itoa(i)+".Name")
		if name == "" {
			// A slot with value members but no name, or a gap with indexed
			// entries of the same stem beyond it, is a malformed list that
			// rejects the request — parsing never silently truncates the
			// remaining attribute list.
			if hasQueryEntryMembers(params, prefix+strconv.Itoa(i)+".Value.") ||
				hasQueryListEntriesBeyond(params, strings.TrimSuffix(prefix, "."), i) {
				return nil, ErrInvalidParameterValue
			}
			break
		}
		attr, err := parseQueryMessageAttributeValue(params, prefix+strconv.Itoa(i)+".Value.")
		if err != nil {
			return nil, err
		}
		result[name] = attr
	}
	return result, nil
}

// parseQueryMessageAttributeValue reads one flattened-query attribute value;
// valuePrefix points at the "...N.Value." key stem.
func parseQueryMessageAttributeValue(params map[string]interface{}, valuePrefix string) (*sqsstore.MessageAttributeValue, error) {
	return buildMessageAttributeValue(
		request.GetParamCaseInsensitive(params, valuePrefix+"DataType"),
		request.GetParamCaseInsensitive(params, valuePrefix+"StringValue"),
		request.GetParamCaseInsensitive(params, valuePrefix+"BinaryValue"))
}

// buildMessageAttributeValue is the single value-coercion rule both wire
// arms feed: the empty-value rule ("Name, type, value and the message body
// must not be empty or null", MessageAttributeValue documentation trait)
// and the base64 BinaryValue decode with its failure policy live here once.
func buildMessageAttributeValue(dataType, stringValue, binaryValue string) (*sqsstore.MessageAttributeValue, error) {
	attr := &sqsstore.MessageAttributeValue{
		DataType: dataType,
	}
	if stringValue != "" {
		attr.StringValue = &stringValue
	}
	if binaryValue != "" {
		decoded, err := sqsstore.DecodeBinaryValue(binaryValue)
		if err != nil {
			return nil, ErrInvalidParameterValue
		}
		attr.BinaryValue = decoded
	}
	if attr.StringValue == nil && attr.BinaryValue == nil {
		return nil, ErrInvalidParameterValue
	}
	return attr, nil
}

// parseSystemAttributes parses the MessageSystemAttributes member of a send
// request from the JSON member map (request level for SendMessage, entry
// level for SendMessageBatch) or, when absent there, from the flattened
// query keys under queryPrefix. An empty queryPrefix disables the query arm
// for containers that hold no request-parameter keys (the JSON batch entry
// map). Only AWSTraceHeader is valid; any other system attribute name is
// rejected with InvalidParameterValue.
func parseSystemAttributes(container map[string]interface{}, queryPrefix string) (map[string]*sqsstore.MessageAttributeValue, error) {
	var attrs map[string]*sqsstore.MessageAttributeValue
	var err error
	if jsonAttrs, ok := container["MessageSystemAttributes"].(map[string]interface{}); ok {
		attrs, err = parseJSONMessageAttributes(jsonAttrs)
	} else if queryPrefix != "" {
		attrs, err = parseQueryMessageAttributes(container, queryPrefix)
	} else {
		attrs = make(map[string]*sqsstore.MessageAttributeValue)
	}
	if err != nil {
		return nil, err
	}
	if err := rejectNonTraceSystemAttributes(attrs); err != nil {
		return nil, err
	}
	return attrs, nil
}

// rejectNonTraceSystemAttributes rejects any system-attribute name other
// than AWSTraceHeader — the only system attribute valid on sends.
func rejectNonTraceSystemAttributes(attrs map[string]*sqsstore.MessageAttributeValue) error {
	for name := range attrs {
		if name != "AWSTraceHeader" {
			return ErrInvalidParameterValue
		}
	}
	return nil
}
