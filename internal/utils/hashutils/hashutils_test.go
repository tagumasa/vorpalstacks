package hashutils

import (
	"bytes"
	"testing"
)

func TestCalculateETag(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantLen  int
		wantSame bool
	}{
		{
			name:     "empty data",
			data:     []byte{},
			wantLen:  64,
			wantSame: true,
		},
		{
			name:     "hello world",
			data:     []byte("hello world"),
			wantLen:  64,
			wantSame: true,
		},
		{
			name:     "test data",
			data:     []byte("test data for etag"),
			wantLen:  64,
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateETag(tt.data)
			if len(got) != tt.wantLen {
				t.Errorf("CalculateETag() len = %d, want %d", len(got), tt.wantLen)
			}
			if tt.wantSame && CalculateETag(tt.data) != got {
				t.Errorf("CalculateETag() = %v, want same", got)
			}
		})
	}
}

func TestCalculateETagFromReader(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantLen int
		wantErr bool
	}{
		{
			name:    "valid reader",
			data:    []byte("hello world"),
			wantLen: 64,
			wantErr: false,
		},
		{
			name:    "empty reader",
			data:    []byte{},
			wantLen: 64,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.data)
			got, err := CalculateETagFromReader(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateETagFromReader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.wantLen {
				t.Errorf("CalculateETagFromReader() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestCalculateETagFromReader_NilReader(t *testing.T) {
	_, err := CalculateETagFromReader(nil)
	if err == nil {
		t.Error("CalculateETagFromReader() expected error for nil reader")
	}
}
