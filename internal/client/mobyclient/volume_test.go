package mobyclient

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "valid RFC3339 time",
			input:    "2024-01-15T10:30:00Z",
			expected: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "valid RFC3339 time with timezone",
			input:    "2024-01-15T10:30:00+09:00",
			expected: time.Date(2024, 1, 15, 1, 30, 0, 0, time.UTC),
		},
		{
			name:     "empty string",
			input:    "",
			expected: time.Time{},
		},
		{
			name:     "invalid time format",
			input:    "invalid",
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTime(tt.input)
			if tt.input == "" || tt.input == "invalid" {
				assert.True(t, result.IsZero())
			} else {
				assert.Equal(t, tt.expected.Unix(), result.Unix())
			}
		})
	}
}
