package sns

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

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
		return "", NewInvalidParameter(fmt.Sprintf("MessageStructure is json but message body is not valid JSON: %s", err.Error()))
	}

	if protocolMsg, ok := msgMap[protocol]; ok {
		return protocolMsg, nil
	}
	if defaultMsg, ok := msgMap["default"]; ok {
		return defaultMsg, nil
	}

	return "", NewInvalidParameter("MessageStructure is json but neither protocol-specific key nor 'default' key found in message body")
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
