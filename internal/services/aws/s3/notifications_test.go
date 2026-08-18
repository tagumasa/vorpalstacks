package s3

import "testing"

// The notification event families must line up with the S3ObjectOp values
// published by the object operations, so a configuration subscribing to a
// documented event type receives the corresponding operation.
func TestMatchS3Event(t *testing.T) {
	tests := []struct {
		name       string
		configured []string
		eventName  string
		want       bool
	}{
		{"exact tagging put", []string{"s3:ObjectTagging:Put"}, "s3:ObjectTagging:Put", true},
		{"exact tagging delete", []string{"s3:ObjectTagging:Delete"}, "s3:ObjectTagging:Delete", true},
		{"exact restore post", []string{"s3:ObjectRestore:Post"}, "s3:ObjectRestore:Post", true},
		{"tagging wildcard matches put", []string{"s3:ObjectTagging:*"}, "s3:ObjectTagging:Put", true},
		{"tagging wildcard matches delete", []string{"s3:ObjectTagging:*"}, "s3:ObjectTagging:Delete", true},
		{"restore wildcard matches post", []string{"s3:ObjectRestore:*"}, "s3:ObjectRestore:Post", true},
		{"created wildcard matches put", []string{"s3:ObjectCreated:*"}, "s3:ObjectCreated:Put", true},
		{"created wildcard misses tagging", []string{"s3:ObjectCreated:*"}, "s3:ObjectTagging:Put", false},
		{"tagging config misses created", []string{"s3:ObjectTagging:Put"}, "s3:ObjectCreated:Put", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchS3Event(tt.configured, tt.eventName); got != tt.want {
				t.Errorf("matchS3Event(%v, %q) = %v, want %v", tt.configured, tt.eventName, got, tt.want)
			}
		})
	}
}
