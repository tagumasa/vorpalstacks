package sns

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// SNS message and batch size limits (AWS documented constants).
const (
	maxMessageSize    = 256 * 1024 // 256 KB
	maxSubjectLength  = 100
	maxBatchTotalSize = 256 * 1024 // 256 KB total for all entries
	maxBatchEntries   = 10
)

// validMessageAttributeDataTypes is the complete set of DataType values
// accepted by SNS message attributes.
var validMessageAttributeDataTypes = map[string]bool{
	"String":       true,
	"Number":       true,
	"String.Array": true,
	"Binary":       true,
}

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

// validatePublishParams validates all top-level Publish parameters that have
// AWS-documented constraints. Returns an InvalidParameterException when any
// constraint is violated.
func validatePublishParams(topic *snsstore.Topic, message, subject, messageStructure, messageGroupId, messageDeduplicationId string) error {
	if len(message) > maxMessageSize {
		return awserrors.NewAWSError("InvalidParameter", fmt.Sprintf("Message too long: %d bytes (maximum %d)", len(message), maxMessageSize), 400)
	}

	if len(subject) > maxSubjectLength {
		return awserrors.NewInvalidParameterException(fmt.Sprintf("Subject too long: %d characters (maximum %d)", len(subject), maxSubjectLength))
	}

	if messageStructure != "" && messageStructure != "json" {
		return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid MessageStructure: %s. Valid value: json", messageStructure))
	}

	if messageStructure == "json" {
		var msgMap map[string]interface{}
		if err := json.Unmarshal([]byte(message), &msgMap); err != nil {
			return awserrors.NewInvalidParameterException(fmt.Sprintf("Invalid parameter: MessageStructure is json but message body is not valid JSON: %s", err.Error()))
		}
		if _, ok := msgMap["default"]; !ok {
			return awserrors.NewInvalidParameterException("Invalid parameter: MessageStructure is json but message body does not contain a 'default' key")
		}
	}

	if !topic.IsFifoTopic() {
		if messageGroupId != "" {
			return awserrors.NewInvalidParameterException("MessageGroupId is only valid for FIFO topics")
		}
		if messageDeduplicationId != "" {
			return awserrors.NewInvalidParameterException("MessageDeduplicationId is only valid for FIFO topics")
		}
	}

	return nil
}

// parseMessageAttributes extracts SNS message attributes from a request params
// map and populates the Message's MessageAttributes field. Returns an error
// (fail-closed) when an attribute entry is malformed or has an invalid
// DataType instead of silently skipping it.
func parseMessageAttributes(params map[string]interface{}, msg *snsstore.Message) error {
	var attrs map[string]interface{}
	for _, key := range []string{"MessageAttributes", "messageAttributes"} {
		if m, ok := params[key].(map[string]interface{}); ok {
			attrs = m
			break
		}
	}
	if attrs == nil {
		return nil
	}

	msg.MessageAttributes = make(map[string]*snsstore.MessageAttribute, len(attrs))
	for k, v := range attrs {
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
