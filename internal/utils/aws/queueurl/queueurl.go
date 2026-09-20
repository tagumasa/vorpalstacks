// Package queueurl resolves the region an SQS queue URL embeds.
package queueurl

import (
	"net/url"
	"strings"
)

// RegionFromQueueURL extracts the region an AWS-form SQS queue URL embeds
// in its host (https://sqs.<region>.amazonaws.com/<account>/<name>).
// Platform queue URLs are host-local and embed none; the caller then falls
// back to its own region context. A host of any other shape names no
// region — extracting a middle label from a foreign host would invent one.
func RegionFromQueueURL(queueURL string) string {
	u, err := url.Parse(queueURL)
	if err != nil {
		return ""
	}
	host := u.Hostname()
	if !strings.HasPrefix(host, "sqs.") {
		return ""
	}
	rest := strings.TrimSuffix(strings.TrimPrefix(host, "sqs."), ".amazonaws.com")
	if rest == "" || strings.Contains(rest, ".") {
		return ""
	}
	return rest
}
