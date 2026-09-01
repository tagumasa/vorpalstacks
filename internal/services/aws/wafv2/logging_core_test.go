package wafv2

import (
	"strings"
	"testing"
)

// TestPutLoggingConfigurationCoreMemberErrors pins the member-specific
// required-member errors: a missing LoggingConfiguration container member
// and a present container whose ResourceArn is empty produce distinct
// messages, both before any store access.
func TestPutLoggingConfigurationCoreMemberErrors(t *testing.T) {
	svc := &WAFv2Service{}

	t.Run("missing LoggingConfiguration member", func(t *testing.T) {
		_, err := svc.putLoggingConfigurationCore(nil, LoggingConfigInput{})
		if err == nil {
			t.Fatal("expected a validation error, got nil")
		}
		if !strings.Contains(err.Error(), "LoggingConfiguration is required") {
			t.Fatalf("error %q does not contain %q", err.Error(), "LoggingConfiguration is required")
		}
	})

	t.Run("present LoggingConfiguration without ResourceArn", func(t *testing.T) {
		_, err := svc.putLoggingConfigurationCore(nil, LoggingConfigInput{LoggingConfigurationPresent: true})
		if err == nil {
			t.Fatal("expected a validation error, got nil")
		}
		if !strings.Contains(err.Error(), "ResourceArn is required") {
			t.Fatalf("error %q does not contain %q", err.Error(), "ResourceArn is required")
		}
	})
}
