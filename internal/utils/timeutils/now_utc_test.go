package timeutils

import (
	"testing"
	"time"
)

func TestNowUTC(t *testing.T) {
	before := time.Now().UTC()
	result := NowUTC()
	after := time.Now().UTC()

	if result.Before(before) || result.After(after) {
		t.Errorf("NowUTC() = %v, should be within test execution window", result)
	}

	if result.Location() != time.UTC {
		t.Errorf("NowUTC() location = %v, want UTC", result.Location())
	}
}

func TestNowString(t *testing.T) {
	result := NowString()

	if result == "" {
		t.Error("NowString() should not return empty string")
	}

	_, err := time.Parse(ISO8601UTCFormat, result)
	if err != nil {
		t.Errorf("NowString() format = %q, should match ISO8601UTCFormat", result)
	}
}

func TestParseTime(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		format    string
		wantErr   bool
		checkFunc func(time.Time) bool
	}{
		{
			name:    "with explicit format",
			input:   "2024-01-01T12:00:00Z",
			format:  time.RFC3339,
			wantErr: false,
			checkFunc: func(t time.Time) bool {
				return t.Year() == 2024 && t.Month() == time.January && t.Day() == 1
			},
		},
		{
			name:    "empty string with auto parse",
			input:   "",
			format:  "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "2024-01-01",
			format:  "invalid-format",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.checkFunc != nil {
				if !tt.checkFunc(result) {
					t.Errorf("ParseTime() result = %v, did not pass check", result)
				}
			}
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name   string
		input  time.Time
		format string
		want   string
	}{
		{
			name:   "RFC3339 format",
			input:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			format: time.RFC3339,
			want:   "2024-01-01T12:00:00Z",
		},
		{
			name:   "custom format",
			input:  time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
			format: "2006-01-02",
			want:   "2024-01-01",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(tt.input, tt.format)
			if result != tt.want {
				t.Errorf("FormatTime() = %q, want %q", result, tt.want)
			}
		})
	}
}
