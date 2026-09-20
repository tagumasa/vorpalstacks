package sns

import (
	"fmt"
	"maps"
	"net/url"
	"slices"
	"strings"

	snsstore "vorpalstacks/internal/store/aws/sns"
)

// protocolRegistry is the single table of subscription-protocol facts:
// which protocols Subscribe accepts, which are auto-confirmed at Subscribe
// time without a token round-trip, which have a delivery handler in the
// delivery engine, which pending protocols receive a confirmation POST,
// and how each protocol's endpoint is validated. The validate-time set,
// the auto-confirm set, the delivery dispatch and the confirmation-POST
// gate all read this table, so they can never drift apart silently.
type protocolInfo struct {
	// autoConfirm: the subscription becomes immediately active without a
	// confirmation round-trip (AWS auto-confirms sqs, lambda and
	// application subscriptions).
	autoConfirm bool
	// deliveryHandler: the delivery engine has a handler for the
	// protocol. A protocol without one fails every delivery at fan-out
	// time (routed to the DLQ when configured).
	deliveryHandler bool
	// confirmationPOST: a pending subscription of this protocol receives
	// a SubscriptionConfirmation POST at Subscribe time — the protocols
	// whose confirmation rides an endpoint the platform can reach.
	confirmationPOST bool
	// rawDelivery: the protocol is in the documented RawMessageDelivery
	// set — "enables raw message delivery to Amazon SQS or HTTP/S
	// endpoints" (SetSubscriptionAttributes, Subscribe).
	rawDelivery bool
	// fifoCapable: the protocol may serve a FIFO topic. "SNS FIFO topics
	// can't deliver messages to customer managed endpoints ... Attempts
	// to subscribe customer managed endpoints to SNS FIFO topics result
	// in errors" (FIFO message delivery) — the customer-managed set is
	// email addresses, mobile apps, phone numbers (SMS) and HTTP(S)
	// endpoints; sqs, lambda and firehose are AWS-managed.
	fifoCapable bool
	// subscribeRefusal: non-empty when Subscribe rejects the protocol
	// outright on this platform, carrying the refusal's reason. Empty
	// means the protocol is accepted.
	subscribeRefusal string
	// validateEndpoint checks the endpoint format at Subscribe time;
	// nil means the protocol's endpoint has no format check.
	validateEndpoint func(endpoint string) error
}

var protocolRegistry = map[string]protocolInfo{
	"http":       {autoConfirm: false, deliveryHandler: true, confirmationPOST: true, rawDelivery: true, validateEndpoint: validateHTTPSEndpoint("http")},
	"https":      {autoConfirm: false, deliveryHandler: true, confirmationPOST: true, rawDelivery: true, validateEndpoint: validateHTTPSEndpoint("https")},
	"email":      {autoConfirm: false, deliveryHandler: false, validateEndpoint: validateEmailEndpoint("email")},
	"email-json": {autoConfirm: false, deliveryHandler: false, validateEndpoint: validateEmailEndpoint("email-json")},
	"sms":        {autoConfirm: false, deliveryHandler: false, validateEndpoint: nil},
	"sqs":        {autoConfirm: true, deliveryHandler: true, rawDelivery: true, fifoCapable: true, validateEndpoint: validateSQSEndpoint},
	"application": {
		autoConfirm: true, deliveryHandler: false,
		// Mobile push delivery requires the external push services (APNs,
		// FCM and siblings) this platform never implements; an
		// auto-confirmed subscription would silently drop every publish,
		// so Subscribe refuses the protocol instead of selling a working
		// subscription. The exclusion is recorded in docs/services.md.
		subscribeRefusal: "mobile push (application protocol) delivery is not available on this platform",
		validateEndpoint: validateApplicationEndpoint,
	},
	"lambda": {autoConfirm: true, deliveryHandler: true, fifoCapable: true, validateEndpoint: validateLambdaEndpoint},
	"firehose": {
		autoConfirm: false, deliveryHandler: false, fifoCapable: true,
		// Firehose delivery lands with the Firehose service; until then
		// subscriptions stay pending and every publish to them routes to
		// the DLQ path — the accept-and-pend carryover recorded beside
		// the Cognito and EventBridge Firehose destinations in
		// docs/services.md.
		validateEndpoint: validateFirehoseEndpoint,
	},
}

// sortedVocabulary renders a registry's keys as the sorted, comma-separated
// list an acceptance error quotes, so the vocabulary an error names can
// never disagree with the check that rejected the input. The protocol,
// platform and action acceptance checks all render through this helper.
func sortedVocabulary[V any](m map[string]V) string {
	return strings.Join(slices.Sorted(maps.Keys(m)), ", ")
}

// protocolAutoConfirms reports whether the protocol's subscriptions are
// auto-confirmed at Subscribe time.
func protocolAutoConfirms(protocol string) bool {
	return protocolRegistry[protocol].autoConfirm
}

// protocolHasDeliveryHandler reports whether the delivery engine owns a
// handler for the protocol — the dispatchability fact the fan-out engine
// consults before dispatching.
func protocolHasDeliveryHandler(protocol string) bool {
	return protocolRegistry[protocol].deliveryHandler
}

// protocolSendsConfirmationPost reports whether a pending subscription of
// the protocol receives a SubscriptionConfirmation POST at Subscribe time.
func protocolSendsConfirmationPost(protocol string) bool {
	return protocolRegistry[protocol].confirmationPOST
}

// protocolSupportsRawDelivery reports whether the protocol is in the
// documented RawMessageDelivery set.
func protocolSupportsRawDelivery(protocol string) bool {
	return protocolRegistry[protocol].rawDelivery
}

// validateProtocolEndpoint dispatches the Subscribe-time endpoint format
// check to the protocol's registry entry. Unknown protocols and protocols
// without a check pass (protocol acceptance itself is validateProtocol's
// job).
func validateProtocolEndpoint(protocol, endpoint string) error {
	info, ok := protocolRegistry[protocol]
	if !ok || info.validateEndpoint == nil {
		return nil
	}
	return info.validateEndpoint(endpoint)
}

// validateHTTPSEndpoint returns the endpoint validator for the http/https
// schemes: the URL must carry the scheme's prefix and stay within the
// endpoint URL length cap.
func validateHTTPSEndpoint(scheme string) func(string) error {
	return func(endpoint string) error {
		if !strings.HasPrefix(endpoint, scheme+"://") {
			return NewInvalidParameter(fmt.Sprintf("Endpoint must be a valid %s URL starting with %s://", strings.ToUpper(scheme), scheme))
		}
		if len(endpoint) > snsstore.MaxEndpointURLLength {
			return NewInvalidParameter(fmt.Sprintf("Endpoint URL too long: %d characters (maximum %d)", len(endpoint), snsstore.MaxEndpointURLLength))
		}
		if _, err := url.Parse(endpoint); err != nil {
			return NewInvalidParameter(fmt.Sprintf("Invalid endpoint URL: %s", err.Error()))
		}
		return nil
	}
}

// validateEmailEndpoint requires an @-containing address for the email
// protocols (delivery itself is a recorded platform exclusion).
func validateEmailEndpoint(protocol string) func(string) error {
	return func(endpoint string) error {
		if !strings.Contains(endpoint, "@") {
			return NewInvalidParameter("Endpoint must be a valid email address for protocol " + protocol)
		}
		return nil
	}
}

// validateSQSEndpoint accepts an SQS queue URL (http/https) or queue ARN.
func validateSQSEndpoint(endpoint string) error {
	if !strings.HasPrefix(endpoint, "http") && !strings.HasPrefix(endpoint, "arn:") {
		return NewInvalidParameter("Endpoint must be a valid SQS queue URL or ARN for protocol sqs")
	}
	return nil
}

// validateApplicationEndpoint requires a platform endpoint ARN.
func validateApplicationEndpoint(endpoint string) error {
	if !strings.HasPrefix(endpoint, "arn:") {
		return NewInvalidParameter("Endpoint must be a valid platform endpoint ARN for protocol application")
	}
	return nil
}

// validateLambdaEndpoint accepts a function ARN or a function name over
// the Lambda name charset (letters, digits, hyphens, underscores).
func validateLambdaEndpoint(endpoint string) error {
	if strings.HasPrefix(endpoint, "arn:") {
		return nil
	}
	if endpoint == "" {
		return NewInvalidParameter("Endpoint must be a valid Lambda function ARN or name")
	}
	for _, c := range endpoint {
		if !isAlnumHyphenUnderscore(c) {
			return NewInvalidParameter("Endpoint must be a valid Lambda function ARN or name")
		}
	}
	return nil
}

// validateFirehoseEndpoint requires a Firehose delivery stream ARN.
func validateFirehoseEndpoint(endpoint string) error {
	if !strings.HasPrefix(endpoint, "arn:") {
		return NewInvalidParameter("Endpoint must be a valid Firehose delivery stream ARN for protocol firehose")
	}
	return nil
}
