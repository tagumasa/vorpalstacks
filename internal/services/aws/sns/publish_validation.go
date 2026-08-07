package sns

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// messageEntrySize estimates the serialised wire size of a single Publish or
// PublishBatch entry: message body + subject + all message attributes
// (name + DataType + StringValue or BinaryValue).  AWS counts the full
// serialised request toward the 256 KB batch limit, so excluding attributes
// under-estimates and allows oversized batches through.
func messageEntrySize(message, subject string, attrs map[string]*snsstore.MessageAttribute) int {
	size := len(message) + len(subject)
	for name, attr := range attrs {
		size += len(name)
		size += len(attr.Type)
		size += len(attr.StringValue)
		if attr.BinaryValue != nil {
			size += base64.StdEncoding.EncodedLen(len(attr.BinaryValue))
		}
	}
	return size
}

// parseMessageAttributes extracts SNS message attributes from a request params
// map and populates the Message's MessageAttributes field. Returns an error
// (fail-closed) when an attribute entry is malformed or has an invalid
// DataType instead of silently skipping it.
//
// M7: enforces max 10 attributes, String value max 256 chars, Binary value
// max 256 bytes.
// M8: enforces attribute name format [a-zA-Z0-9_.-]{1,256}.
func parseMessageAttributes(params map[string]interface{}, msg *snsstore.Message) error {
	var attrs map[string]interface{}
	for _, key := range []string{"MessageAttributes", "messageAttributes"} {
		if m, ok := params[key].(map[string]interface{}); ok {
			attrs = m
			break
		}
	}

	// Also handle the AWS Query API flat-key format. The SDK sends message
	// attributes as:
	//   MessageAttributes.entry.N.Name=<name>
	//   MessageAttributes.entry.N.Value.DataType=<type>
	//   MessageAttributes.entry.N.Value.StringValue=<value>
	// The HTTP query parser stores these as flat string keys, so we parse
	// them manually.
	if attrs == nil {
		attrs = make(map[string]interface{})
		for i := 1; i <= 30; i++ {
			name := request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Name", i))
			if name == "" {
				break
			}
			attrs[name] = map[string]interface{}{
				"DataType":    request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Value.DataType", i)),
				"StringValue": request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Value.StringValue", i)),
				"BinaryValue": request.GetStringParam(params, fmt.Sprintf("MessageAttributes.entry.%d.Value.BinaryValue", i)),
			}
		}
		if len(attrs) == 0 {
			attrs = nil
		}
	}

	if attrs == nil {
		return nil
	}

	// M7: maximum 10 message attributes.
	if len(attrs) > maxMessageAttributes {
		return awserrors.NewInvalidParameterException(fmt.Sprintf("Too many message attributes: %d (maximum %d)", len(attrs), maxMessageAttributes))
	}

	msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute, len(attrs))
	for k, v := range attrs {
		// M8: validate attribute name format.
		if err := validateMessageAttributeName(k); err != nil {
			return err
		}

		attrMap, ok := v.(map[string]interface{})
		if !ok {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: value must be a map", k))
		}

		dataType := firstString(attrMap, "DataType", "dataType")
		if dataType == "" {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: DataType is required", k))
		}
		if !validMessageAttributeDataTypes[dataType] {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: DataType %q is not valid (String, Number, String.Array, Binary)", k, dataType))
		}

		attr := &snsstore.MessageAttribute{
			Type:        dataType,
			StringValue: firstString(attrMap, "StringValue", "stringValue"),
		}
		if raw := firstString(attrMap, "BinaryValue", "binaryValue"); raw != "" {
			decoded, err := base64.StdEncoding.DecodeString(raw)
			if err != nil {
				return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid message attribute %q: BinaryValue is not valid base64", k))
			}
			attr.BinaryValue = decoded
		}

		// M7: validate attribute value sizes.
		if err := validateMessageAttributeLimits(k, attr.StringValue, attr.BinaryValue); err != nil {
			return err
		}

		msg.MessageAttributes[k] = attr
	}

	return nil
}

// extractProtocolMessage returns the protocol-specific message from a
// structured (MessageStructure=json) message. When MessageStructure is not
// "json", the raw message is returned unchanged.
//
// Fail-closed: when MessageStructure is "json" but the message body cannot be
// parsed as valid JSON, or when neither the protocol-specific key nor a
// "default" key exists, an error is returned. AWS rejects these cases with
// InvalidParameter.
func extractProtocolMessage(msg *snsstore.Message, protocol string) (string, error) {
	if msg.MessageStructure != "json" {
		return msg.Message, nil
	}

	var msgMap map[string]string
	if err := json.Unmarshal([]byte(msg.Message), &msgMap); err != nil {
		return "", awserrors.NewInvalidParameterException(fmt.Sprintf("MessageStructure is json but message body is not valid JSON: %s", err.Error()))
	}

	if protocolMsg, ok := msgMap[protocol]; ok {
		return protocolMsg, nil
	}
	if defaultMsg, ok := msgMap["default"]; ok {
		return defaultMsg, nil
	}

	return "", awserrors.NewInvalidParameterException("MessageStructure is json but neither protocol-specific key nor 'default' key found in message body")
}

// firstString returns the first non-empty string value found for any of the
// given keys in the map.
func firstString(m map[string]interface{}, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k].(string); ok && v != "" {
			return v
		}
	}
	return ""
}

// messageAttributeValue returns the serialisable value for an SNS message
// attribute. String/Number/String.Array types return the string value;
// Binary types return the base64-encoded representation.
func messageAttributeValue(attr *snsstore.MessageAttribute) string {
	if len(attr.BinaryValue) > 0 {
		return base64.StdEncoding.EncodeToString(attr.BinaryValue)
	}
	return attr.StringValue
}
