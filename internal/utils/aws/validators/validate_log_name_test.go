package validators

import (
	"strings"
	"testing"
)

func TestValidateLogStreamName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple name", "log-stream-001", false},
		{"valid with underscore", "log_stream_001", false},
		{"valid with slash", "app/service/logs", false},
		{"invalid with colon", "region:zone:logs", true},
		{"valid max length", strings.Repeat("a", 512), false},
		{"valid alphanumeric", "LogStream123", false},
		{"empty string", "", true},
		{"too long", strings.Repeat("a", 513), true},
		{"invalid character - space", "log stream", true},
		{"invalid character - period", "log.stream", true},
		{"invalid character - asterisk", "log*stream", true},
		{"invalid character - parenthesis", "log(stream)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogStreamName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogStreamName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLogGroupName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple name", "/aws/lambda/my-function", false},
		{"valid with hyphen", "log-group-001", false},
		{"valid with underscore", "log_group_001", false},
		{"valid with slash", "app/service/logs", false},
		{"valid with colon", "region:zone:logs", false},
		{"valid max length", strings.Repeat("a", 512), false},
		{"valid path format", "/aws/lambda/function-name", false},
		{"empty string", "", true},
		{"too long", strings.Repeat("a", 513), true},
		{"invalid character - space", "log group", true},
		{"invalid character - period", "log.group", true},
		{"invalid character - asterisk", "log*group", true},
		{"invalid character - at sign", "log@group", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLogGroupName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateLogGroupName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
