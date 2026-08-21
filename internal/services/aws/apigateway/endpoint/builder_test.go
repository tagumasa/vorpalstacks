package endpoint

import (
	"strings"
	"testing"
)

func TestSQSQueueURL(t *testing.T) {
	got := SQSQueueURL("123456789012", "my-queue")
	if !strings.HasSuffix(got, "/123456789012/my-queue") {
		t.Errorf("SQSQueueURL = %q, want suffix /123456789012/my-queue", got)
	}
}
