package timeutils

import (
	"testing"
	"time"
)

func TestFormatTimeCustom(t *testing.T) {
	tests := []struct {
		name   string
		t      time.Time
		format string
		want   string
	}{
		{
			name:   "with custom format",
			t:      time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
			format: "2006-01-02",
			want:   "2024-01-15",
		},
		{
			name:   "default format",
			t:      time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
			format: "",
			want:   "2024-01-15T10:30:45.000Z",
		},
		{
			name:   "ISO8601 format",
			t:      time.Date(2024, 1, 15, 10, 30, 45, 123000000, time.UTC),
			format: ISO8601UTCFormat,
			want:   "2024-01-15T10:30:45.123Z",
		},
		{
			name:   "RFC3339 format",
			t:      time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
			format: RFC3339Format,
			want:   "2024-01-15T10:30:45Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatTime(tt.t, tt.format)
			if got != tt.want {
				t.Errorf("FormatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
