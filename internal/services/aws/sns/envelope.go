package sns

// envelope.go is the notification-envelope contract: the sole construction
// and signing of the envelopes every delivery protocol receives (the SNS
// Notification, the SubscriptionConfirmation) and the protocol-message
// extraction the delivery paths consume. Delivery routing lives in
// delivery.go.

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

// snsServiceURL returns the base of the AWS service endpoint URLs placed
// in notification and confirmation envelopes. The platform cannot serve
// these AWS-account-owned domains (recorded platform exclusion); the URLs
// exist for envelope shape parity, and this helper is their single
// definition.
func snsServiceURL(region string) string {
	return fmt.Sprintf("https://sns.%s.amazonaws.com", region)
}

// envelopeTimestamp renders the Timestamp member every envelope carries:
// UTC at millisecond precision — the format AWS sample envelopes use
// (e.g. 2019-01-31T04:37:04.321Z). This is the single definition; no
// envelope construction may format its own timestamp.
func envelopeTimestamp(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z07:00")
}

// buildNotificationEnvelope constructs the SNS Notification envelope every
// delivery protocol receives in non-raw mode — SQS, HTTP(S) and Lambda
// deliveries serialise this same object (the Lambda record nests it under
// Sns), so member inclusion and the timestamp format cannot drift per
// protocol. signatureVersion is the topic's SignatureVersion attribute
// ("By default, SignatureVersion is set to 1"): the envelope carries it as
// a member and the signing step selects the hash it names. The envelope is
// not yet signed: the caller applies signNotificationEnvelope before
// serialising.
func (s *SNSService) buildNotificationEnvelope(msg *snsstore.Message, sub *snsstore.Subscription, region string, message string, signatureVersion string) map[string]interface{} {
	if region == "" {
		region = s.defaultRegion
	}
	payload := map[string]interface{}{
		"Type":             "Notification",
		"MessageId":        msg.MessageId,
		"TopicArn":         msg.TopicArn,
		"Message":          message,
		"Timestamp":        envelopeTimestamp(msg.PublishedTimestamp),
		"SignatureVersion": signatureVersion,
		"UnsubscribeURL":   snsServiceURL(region) + fmt.Sprintf("/?Action=Unsubscribe&SubscriptionArn=%s", sub.SubscriptionArn),
	}

	if msg.Subject != "" {
		payload["Subject"] = msg.Subject
	}

	if len(msg.MessageAttributes) > 0 {
		attrs := make(map[string]interface{}, len(msg.MessageAttributes))
		for k, v := range msg.MessageAttributes {
			attrs[k] = map[string]interface{}{
				"Type":  v.Type,
				"Value": messageAttributeValue(v),
			}
		}
		payload["MessageAttributes"] = attrs
	}

	return payload
}

// buildConfirmationEnvelope constructs the SubscriptionConfirmation message
// a pending HTTP(S) subscription receives: the same envelope vocabulary as
// a notification (Type names the message kind, the timestamp format is
// shared, the signature version the topic selects) plus the Token and
// SubscribeURL members the subscriber confirms with. The caller applies
// signNotificationEnvelope before serialising.
func (s *SNSService) buildConfirmationEnvelope(sub *snsstore.Subscription, region string, signatureVersion string) map[string]interface{} {
	if region == "" {
		region = s.defaultRegion
	}
	return map[string]interface{}{
		"Type":             "SubscriptionConfirmation",
		"MessageId":        uuid.New().String(),
		"Token":            sub.ConfirmationToken,
		"TopicArn":         sub.TopicArn,
		"Message":          fmt.Sprintf("You have chosen to subscribe to the topic %s.\nTo confirm the subscription, visit the SubscribeURL included in this message.", sub.TopicArn),
		"SubscribeURL":     snsServiceURL(region) + fmt.Sprintf("/?Action=ConfirmSubscription&TopicArn=%s&Token=%s", sub.TopicArn, sub.ConfirmationToken),
		"Timestamp":        envelopeTimestamp(time.Now()),
		"SignatureVersion": signatureVersion,
	}
}

// signNotificationEnvelope adds the signature fields to a constructed
// envelope (notification or confirmation) in place; every delivery
// protocol receives the signed object.
func (s *SNSService) signNotificationEnvelope(payload map[string]interface{}, region string, signatureVersion string) {
	if signature, certURL := s.signEnvelope(payload, region, signatureVersion); signature != "" {
		payload["Signature"] = signature
		payload["SigningCertURL"] = certURL
	}
}

// signEnvelope signs an envelope with the service's signing key.
// SignatureVersion selects the hash — version 1 is "an SHA1 hash of the
// message", version 2 "an SHA256 hash of the message" (verifying the
// signatures of Amazon SNS messages) — and the signature is RSA over that
// hash. The StringToSign is the envelope's string members as name-value
// pairs in byte-sorted name order, each pair ending with a newline:
// notifications carry Message, MessageId, Subject (when present),
// Timestamp, TopicArn and Type; SubscriptionConfirmation (and
// UnsubscribeConfirmation) carry SubscribeURL and Token instead of Subject.
func (s *SNSService) signEnvelope(payload map[string]interface{}, region string, signatureVersion string) (string, string) {
	s.initSigningKey()
	if s.signingKey == nil {
		logs.Warn("Failed to get cached RSA key for signing")
		return "", ""
	}

	certURL := snsServiceURL(region) + fmt.Sprintf("/SimpleNotificationService-%x.pem", sha256.Sum256(s.signingCertPEM))

	var strToSign strings.Builder
	for _, name := range []string{"Message", "MessageId", "Subject", "SubscribeURL", "Timestamp", "Token", "TopicArn", "Type"} {
		value, ok := payload[name].(string)
		if !ok || value == "" {
			continue
		}
		strToSign.WriteString(name)
		strToSign.WriteString("\n")
		strToSign.WriteString(value)
		strToSign.WriteString("\n")
	}

	hash := crypto.SHA1
	if signatureVersion == "2" {
		hash = crypto.SHA256
	}
	digest := hash.New()
	digest.Write([]byte(strToSign.String()))
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.signingKey, hash, digest.Sum(nil))
	if err != nil {
		logs.Warn("Failed to sign payload", logs.String("error", err.Error()))
		return "", ""
	}

	return base64.StdEncoding.EncodeToString(signature), certURL
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
