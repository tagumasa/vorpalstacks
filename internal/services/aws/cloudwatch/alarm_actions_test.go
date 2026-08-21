package cloudwatch

import "testing"

// TestExtractAlarmNameFromARN pins alarm-name extraction from the
// resource field of a CloudWatch alarm ARN
// (arn:<partition>:cloudwatch:<region>:<account>:alarm:<name>).
func TestExtractAlarmNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{"standard alarm ARN", "arn:aws:cloudwatch:us-east-1:123456789012:alarm:MyAlarm", "MyAlarm"},
		{"name with dots", "arn:aws:cloudwatch:us-east-1:123456789012:alarm:my.team.alarm", "my.team.alarm"},
		{"non-alarm resource", "arn:aws:cloudwatch:us-east-1:123456789012:insight-rule/MyRule", ""},
		{"bare name", "MyAlarm", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractAlarmNameFromARN(tt.arn); got != tt.want {
				t.Errorf("extractAlarmNameFromARN(%q) = %q, want %q", tt.arn, got, tt.want)
			}
		})
	}
}
