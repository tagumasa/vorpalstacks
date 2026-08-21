// Package endpoint builds the platform-local URL forms that API Gateway
// integrations reference when they address other services on this
// deployment instead of AWS regional endpoints.
package endpoint

import (
	"fmt"

	appconfig "vorpalstacks/internal/config"
)

// SQSQueueURL constructs the URL of an SQS queue on this platform for use
// in API Gateway integration request URIs.
func SQSQueueURL(accountID, queueName string) string {
	return fmt.Sprintf("%s/%s/%s", appconfig.BaseURL(), accountID, queueName)
}
