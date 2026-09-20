package sns

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"regexp"
	"strings"
	"testing"
	"time"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

// envelopeTimestampPattern pins the envelope Timestamp format: UTC at
// millisecond precision (AWS sample envelopes carry e.g.
// 2019-01-31T04:37:04.321Z).
var envelopeTimestampPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$`)

// TestNotificationEnvelopeVocabulary pins the member set of the one
// notification envelope every protocol receives: the base members, the
// conditional members (Subject only when set, MessageAttributes only when
// present), and the signature fields the signing step adds — the same
// object reaches SQS, HTTP(S), Lambda and DLQ destinations, so this is the
// whole protocol family's vocabulary in one pin.
func TestNotificationEnvelopeVocabulary(t *testing.T) {
	svc, _ := newFanoutTestService(t)
	msg := &snsstore.Message{
		MessageId:          "mid-vocab",
		TopicArn:           "arn:aws:sns:us-east-1:123456789012:vocab-topic",
		Subject:            "vocab subject",
		Message:            "vocab body",
		PublishedTimestamp: time.Date(2026, 9, 19, 12, 34, 56, 789000000, time.UTC),
		MessageAttributes: map[string]*snsstore.MessageAttribute{
			"rank": {Type: "Number", StringValue: "3"},
		},
	}
	sub := &snsstore.Subscription{
		SubscriptionArn: "arn:aws:sns:us-east-1:123456789012:vocab-topic:00000000-0000-0000-0000-000000000000",
	}

	envelope := svc.buildNotificationEnvelope(msg, sub, "us-east-1", "vocab body", "2")
	svc.signNotificationEnvelope(envelope, "us-east-1", "2")

	for name, want := range map[string]interface{}{
		"Type":             "Notification",
		"MessageId":        "mid-vocab",
		"TopicArn":         msg.TopicArn,
		"Message":          "vocab body",
		"Subject":          "vocab subject",
		"SignatureVersion": "2",
	} {
		if envelope[name] != want {
			t.Errorf("envelope[%s] = %v, want %v", name, envelope[name], want)
		}
	}
	if ts, _ := envelope["Timestamp"].(string); !envelopeTimestampPattern.MatchString(ts) {
		t.Errorf("envelope Timestamp = %q, want millisecond-precision UTC", ts)
	}
	if unsub, _ := envelope["UnsubscribeURL"].(string); unsub == "" {
		t.Error("envelope carries no UnsubscribeURL")
	} else if want := "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=" + sub.SubscriptionArn; unsub != want {
		t.Errorf("envelope UnsubscribeURL = %q, want %q", unsub, want)
	}
	if sig, _ := envelope["Signature"].(string); sig == "" {
		t.Error("the signed envelope carries no Signature")
	}
	if cert, _ := envelope["SigningCertURL"].(string); cert == "" {
		t.Error("the signed envelope carries no SigningCertURL")
	}
	attrs, ok := envelope["MessageAttributes"].(map[string]interface{})
	if !ok {
		t.Fatalf("envelope MessageAttributes = %#v, want the attribute map", envelope["MessageAttributes"])
	}
	rank, ok := attrs["rank"].(map[string]interface{})
	if !ok || rank["Type"] != "Number" || rank["Value"] != "3" {
		t.Errorf("envelope attribute rank = %#v, want Type Number Value 3", attrs["rank"])
	}
}

// TestNotificationEnvelopeOmitsAbsentMembers pins the conditional members:
// a message without Subject and without message attributes produces an
// envelope without the Subject and MessageAttributes keys — never explicit
// empty members.
func TestNotificationEnvelopeOmitsAbsentMembers(t *testing.T) {
	svc, _ := newFanoutTestService(t)
	msg := &snsstore.Message{
		MessageId:          "mid-bare",
		TopicArn:           "arn:aws:sns:us-east-1:123456789012:bare-topic",
		Message:            "bare body",
		PublishedTimestamp: time.Date(2026, 9, 19, 12, 34, 56, 0, time.UTC),
	}
	sub := &snsstore.Subscription{SubscriptionArn: "arn:aws:sns:us-east-1:123456789012:bare-topic:sub"}

	envelope := svc.buildNotificationEnvelope(msg, sub, "us-east-1", "bare body", "1")

	if _, present := envelope["Subject"]; present {
		t.Error("envelope carries Subject for a message without one")
	}
	if _, present := envelope["MessageAttributes"]; present {
		t.Error("envelope carries MessageAttributes for a message without any")
	}
	if envelope["SignatureVersion"] != "1" {
		t.Errorf("envelope SignatureVersion = %v, want the passed version 1", envelope["SignatureVersion"])
	}
}

// TestSignatureVersionSelectsSigningHash pins the honoured SignatureVersion:
// "SignatureVersion1 – Uses an SHA1 hash of the message. SignatureVersion2 –
// Uses an SHA256 hash of the message" — each envelope's signature verifies
// against the hash its SignatureVersion member names, and only against that
// hash.
func TestSignatureVersionSelectsSigningHash(t *testing.T) {
	svc, _ := newFanoutTestService(t)
	svc.initSigningKey()
	if svc.signingKey == nil {
		t.Fatal("signing key unavailable")
	}
	msg := &snsstore.Message{
		MessageId:          "mid-sigver",
		TopicArn:           "arn:aws:sns:us-east-1:123456789012:sigver-topic",
		Message:            "signature version body",
		PublishedTimestamp: time.Date(2026, 9, 19, 12, 34, 56, 0, time.UTC),
	}
	sub := &snsstore.Subscription{SubscriptionArn: "arn:aws:sns:us-east-1:123456789012:sigver-topic:sub"}

	for _, version := range []string{"1", "2"} {
		envelope := svc.buildNotificationEnvelope(msg, sub, "us-east-1", msg.Message, version)
		svc.signNotificationEnvelope(envelope, "us-east-1", version)

		sigB64, _ := envelope["Signature"].(string)
		if sigB64 == "" {
			t.Fatalf("version %s envelope carries no Signature", version)
		}
		sig, err := base64.StdEncoding.DecodeString(sigB64)
		if err != nil {
			t.Fatalf("version %s signature decode: %v", version, err)
		}

		var strToSign strings.Builder
		for _, name := range []string{"Message", "MessageId", "Subject", "SubscribeURL", "Timestamp", "Token", "TopicArn", "Type"} {
			if value, ok := envelope[name].(string); ok && value != "" {
				strToSign.WriteString(name)
				strToSign.WriteString("\n")
				strToSign.WriteString(value)
				strToSign.WriteString("\n")
			}
		}

		wantHash := crypto.SHA1
		notWantHash := crypto.SHA256
		if version == "2" {
			wantHash, notWantHash = crypto.SHA256, crypto.SHA1
		}
		digest := wantHash.New()
		digest.Write([]byte(strToSign.String()))
		if err := rsa.VerifyPKCS1v15(&svc.signingKey.PublicKey, wantHash, digest.Sum(nil), sig); err != nil {
			t.Errorf("version %s signature does not verify against its documented hash: %v", version, err)
		}
		wrongDigest := notWantHash.New()
		wrongDigest.Write([]byte(strToSign.String()))
		if err := rsa.VerifyPKCS1v15(&svc.signingKey.PublicKey, notWantHash, wrongDigest.Sum(nil), sig); err == nil {
			t.Errorf("version %s signature verifies against the other version's hash — the version is not honoured", version)
		}
	}
}

// TestConfirmationEnvelopeVocabulary pins the SubscriptionConfirmation
// message: the notification vocabulary's shared members plus Token and
// SubscribeURL, and the signature fields — AWS signs confirmation messages,
// and the platform signs them with the same signing step as notifications.
func TestConfirmationEnvelopeVocabulary(t *testing.T) {
	svc, _ := newFanoutTestService(t)
	sub := &snsstore.Subscription{
		SubscriptionArn:   "arn:aws:sns:us-east-1:123456789012:conf-topic:sub",
		TopicArn:          "arn:aws:sns:us-east-1:123456789012:conf-topic",
		ConfirmationToken: "conf-token",
	}

	envelope := svc.buildConfirmationEnvelope(sub, "us-east-1", "1")
	svc.signNotificationEnvelope(envelope, "us-east-1", "1")

	if envelope["Type"] != "SubscriptionConfirmation" {
		t.Errorf("envelope Type = %v, want SubscriptionConfirmation", envelope["Type"])
	}
	if envelope["Token"] != "conf-token" {
		t.Errorf("envelope Token = %v, want the subscription's confirmation token", envelope["Token"])
	}
	if id, _ := envelope["MessageId"].(string); id == "" {
		t.Error("confirmation envelope carries no MessageId")
	}
	if url, _ := envelope["SubscribeURL"].(string); url == "" {
		t.Error("confirmation envelope carries no SubscribeURL")
	} else if want := "https://sns.us-east-1.amazonaws.com/?Action=ConfirmSubscription&TopicArn=" + sub.TopicArn + "&Token=conf-token"; url != want {
		t.Errorf("confirmation SubscribeURL = %q, want %q", url, want)
	}
	if ts, _ := envelope["Timestamp"].(string); !envelopeTimestampPattern.MatchString(ts) {
		t.Errorf("confirmation Timestamp = %q, want millisecond-precision UTC", ts)
	}
	if sig, _ := envelope["Signature"].(string); sig == "" {
		t.Error("the signed confirmation envelope carries no Signature")
	}
	if cert, _ := envelope["SigningCertURL"].(string); cert == "" {
		t.Error("the signed confirmation envelope carries no SigningCertURL")
	}
}
