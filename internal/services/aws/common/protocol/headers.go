// Package protocol provides AWS protocol encoding utilities for vorpalstacks.
package protocol

import "net/http"

// SetProtocolHeaders sets the protocol-specific HTTP headers for AWS API responses.
func SetProtocolHeaders(w http.ResponseWriter, contentType, xAmzTarget, smithyProtocol string, xAmznQueryMode bool) {
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	if xAmzTarget != "" {
		w.Header().Set("X-Amz-Target", xAmzTarget)
	}

	if smithyProtocol != "" {
		w.Header().Set("smithy-protocol", smithyProtocol)
	}

	if xAmznQueryMode {
		w.Header().Set("X-Amzn-Query-Mode", "true")
	}
}
