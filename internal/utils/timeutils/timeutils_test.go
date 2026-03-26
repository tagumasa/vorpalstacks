package timeutils

import (
	"testing"
	"time"
)

func TestFormatTimeCustom(t *testing.T) {
	tests := []struct {
		name    string
		t       time.Time
		format  string
		want    string
		wantErr bool
	}{
		{
			name:    "with custom format",
			t:       time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
			format:  "2006-01-02",
			want:    "2024-01-15",
			wantErr: false,
		},
		{
			name:    "default format",
			t:       time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
			format:  "",
			want:    "2024-01-15T10:30:45.000Z",
			wantErr: false,
		},
		{
			name:    "ISO8601 format",
			t:       time.Date(2024, 1, 15, 10, 30, 45, 123000000, time.UTC),
			format:  ISO8601UTCFormat,
			want:    "2024-01-15T10:30:45.123Z",
			wantErr: false,
		},
		{
			name:    "RFC3339 format",
			t:       time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC),
			format:  RFC3339Format,
			want:    "2024-01-15T10:30:45Z",
			wantErr: false,
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

func TestDurationString(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
		want string
	}{
		{"zero", 0, "0s"},
		{"seconds only", 45 * time.Second, "45s"},
		{"minutes and seconds", 90 * time.Second, "1m30s"},
		{"hours minutes seconds", 2*time.Hour + 5*time.Minute + 30*time.Second, "2h5m30s"},
		{"hours only", 3 * time.Hour, "3h"},
		{"minutes only", 15 * time.Minute, "15m"},
		{"sub-second", 500 * time.Millisecond, "500ms"},
		{"complex", 1*time.Hour + 1*time.Minute + 1*time.Second, "1h1m1s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DurationString(tt.d); got != tt.want {
				t.Errorf("DurationString(%v) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}

func TestDurationString_LessThanSecond(t *testing.T) {
	d := 500 * time.Millisecond
	got := DurationString(d)
	if got != "500ms" {
		t.Errorf("DurationString() = %v, want 500ms", got)
	}
}
