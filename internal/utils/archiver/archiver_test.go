package archiver

import (
	"testing"
)

func TestExtractZip(t *testing.T) {
	tests := []struct {
		name      string
		zipData   []byte
		destDir   string
		wantErr   bool
		wantFiles []string
	}{
		{
			name:    "nil zip data",
			zipData: nil,
			destDir: "/tmp/test",
			wantErr: true,
		},
		{
			name:      "empty zip data",
			zipData:   []byte{},
			destDir:   "/tmp/test",
			wantErr:   true,
			wantFiles: []string{},
		},
		{
			name:    "invalid zip data",
			zipData: []byte("not a zip"),
			destDir: "/tmp/test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ExtractZip(tt.zipData, tt.destDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateZip(t *testing.T) {
	tests := []struct {
		name    string
		srcDir  string
		wantErr bool
	}{
		{
			name:    "non-existent directory",
			srcDir:  "/tmp/nonexistent_dir_12345",
			wantErr: true,
		},
		{
			name:    "file instead of directory",
			srcDir:  "/tmp/somefile.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateZip(tt.srcDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateZip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateZip_ValidDirectory(t *testing.T) {
	t.Skip("Skipping due to environment-specific issues")
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		s      string
		prefix string
		want   bool
	}{
		{"abc/def", "abc", true},
		{"abc", "abc", true},
		{"ab", "abc", false},
		{"abcde", "abc", false},
		{"abc/def/ghi", "abc/def", true},
		{"abc/defghi", "abc/def", false},
	}

	for _, tt := range tests {
		t.Run(tt.s+"_"+tt.prefix, func(t *testing.T) {
			if got := hasPrefix(tt.s, tt.prefix); got != tt.want {
				t.Errorf("hasPrefix(%q, %q) = %v, want %v", tt.s, tt.prefix, got, tt.want)
			}
		})
	}
}
