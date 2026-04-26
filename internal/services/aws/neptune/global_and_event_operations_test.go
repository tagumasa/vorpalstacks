package neptune

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClusterIDFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{"standard ARN", "arn:aws:neptune:us-east-1:123456789:cluster:my-cluster", "my-cluster"},
		{"cluster with hyphens", "arn:aws:neptune:eu-west-1:999:cluster:my-test-cluster-1", "my-test-cluster-1"},
		{"not an ARN returns input", "just-a-string", "just-a-string"},
		{"ARN without cluster segment", "arn:aws:neptune:us-east-1:123456789:snapshot:my-snap", "arn:aws:neptune:us-east-1:123456789:snapshot:my-snap"},
		{"empty string", "", ""},
		{"cluster at end", "arn:aws:neptune:us-east-1:123:cluster:", ""},
		{"multiple colons in ID", "arn:aws:neptune:us-east-1:123:cluster:my:complex:id", "my"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clusterIDFromARN(tt.arn)
			assert.Equal(t, tt.want, got)
		})
	}
}
