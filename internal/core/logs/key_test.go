package logs

import (
	"testing"
)

func TestNewKey(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		tenantID  string
		region    string
		id        string
		wantNS    string
		wantID    string
	}{
		{
			name:      "basic key",
			namespace: "logs",
			tenantID:  "tenant1",
			region:    "us-east-1",
			id:        "log123",
			wantNS:    "logs",
			wantID:    "log123",
		},
		{
			name:      "empty parts",
			namespace: "",
			tenantID:  "",
			region:    "",
			id:        "",
			wantNS:    "",
			wantID:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := NewKey(tt.namespace, tt.tenantID, tt.region, tt.id)
			if key.Namespace != tt.wantNS {
				t.Errorf("Namespace = %v, want %v", key.Namespace, tt.wantNS)
			}
			if key.ID != tt.wantID {
				t.Errorf("ID = %v, want %v", key.ID, tt.wantID)
			}
		})
	}
}

func TestNewLogKey(t *testing.T) {
	key := NewLogKey("tenant1", "us-east-1", "log123")
	if key.Namespace != "logs" {
		t.Errorf("Namespace = %v, want logs", key.Namespace)
	}
	if key.TenantID != "tenant1" {
		t.Errorf("TenantID = %v, want tenant1", key.TenantID)
	}
	if key.Region != "us-east-1" {
		t.Errorf("Region = %v, want us-east-1", key.Region)
	}
	if key.ID != "log123" {
		t.Errorf("ID = %v, want log123", key.ID)
	}
}

func TestKey_Encode(t *testing.T) {
	key := &Key{
		Namespace: "logs",
		TenantID:  "tenant1",
		Region:    "us-east-1",
		ID:        "log123",
	}
	want := "logs:tenant1:us-east-1:log123"
	got := key.Encode()
	if got != want {
		t.Errorf("Encode() = %v, want %v", got, want)
	}
}

func TestKey_EncodePrefix(t *testing.T) {
	key := &Key{
		Namespace: "logs",
		TenantID:  "tenant1",
		Region:    "us-east-1",
		ID:        "log123",
	}
	want := "logs:tenant1:us-east-1"
	got := key.EncodePrefix()
	if got != want {
		t.Errorf("EncodePrefix() = %v, want %v", got, want)
	}
}

func TestDecodeKey(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNil bool
		wantNS  string
		wantID  string
	}{
		{
			name:    "valid key",
			input:   "logs:tenant1:us-east-1:log123",
			wantNil: false,
			wantNS:  "logs",
			wantID:  "log123",
		},
		{
			name:    "invalid - not enough parts",
			input:   "logs:tenant1:us-east-1",
			wantNil: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := DecodeKey(tt.input)
			if (key == nil) != tt.wantNil {
				t.Errorf("DecodeKey() = nil %v, want nil %v", key == nil, tt.wantNil)
				return
			}
			if !tt.wantNil {
				if key.Namespace != tt.wantNS {
					t.Errorf("Namespace = %v, want %v", key.Namespace, tt.wantNS)
				}
				if key.ID != tt.wantID {
					t.Errorf("ID = %v, want %v", key.ID, tt.wantID)
				}
			}
		})
	}
}

func TestIndexKey_Encode(t *testing.T) {
	tests := []struct {
		name       string
		indexType  IndexType
		tenantID   string
		region     string
		segment1   string
		segment2   string
		entryID    string
		wantPrefix string
	}{
		{
			name:       "time index",
			indexType:  IndexByTime,
			tenantID:   "tenant1",
			region:     "us-east-1",
			segment1:   "20240101",
			segment2:   "",
			entryID:    "entry1",
			wantPrefix: "idx_time",
		},
		{
			name:       "level index",
			indexType:  IndexByLevel,
			tenantID:   "tenant1",
			region:     "us-east-1",
			segment1:   "error",
			segment2:   "1234567890",
			entryID:    "entry1",
			wantPrefix: "idx_level",
		},
		{
			name:       "source index",
			indexType:  IndexBySource,
			tenantID:   "tenant1",
			region:     "us-east-1",
			segment1:   "myapp",
			segment2:   "1234567890",
			entryID:    "entry1",
			wantPrefix: "idx_source",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := &IndexKey{
				IndexType: tt.indexType,
				TenantID:  tt.tenantID,
				Region:    tt.region,
				Segment1:  tt.segment1,
				Segment2:  tt.segment2,
				EntryID:   tt.entryID,
			}
			got := key.Encode()
			if len(got) == 0 {
				t.Error("Encode() returned empty string")
			}
		})
	}
}

func TestIndexKey_EncodePrefix(t *testing.T) {
	tests := []struct {
		name      string
		indexType IndexType
		tenantID  string
		region    string
		segment1  string
		segment2  string
	}{
		{
			name:      "with segment2",
			indexType: IndexByLevel,
			tenantID:  "tenant1",
			region:    "us-east-1",
			segment1:  "error",
			segment2:  "1234567890",
		},
		{
			name:      "with only segment1",
			indexType: IndexByTime,
			tenantID:  "tenant1",
			region:    "us-east-1",
			segment1:  "20240101",
			segment2:  "",
		},
		{
			name:      "no segments",
			indexType: IndexBySource,
			tenantID:  "tenant1",
			region:    "us-east-1",
			segment1:  "",
			segment2:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := &IndexKey{
				IndexType: tt.indexType,
				TenantID:  tt.tenantID,
				Region:    tt.region,
				Segment1:  tt.segment1,
				Segment2:  tt.segment2,
			}
			got := key.EncodePrefix()
			if len(got) == 0 {
				t.Error("EncodePrefix() returned empty string")
			}
		})
	}
}

func TestNewTimeIndexKey(t *testing.T) {
	key := NewTimeIndexKey("tenant1", "us-east-1", "20240101", "entry1")
	if key.IndexType != IndexByTime {
		t.Errorf("IndexType = %v, want IndexByTime", key.IndexType)
	}
	if key.TenantID != "tenant1" {
		t.Errorf("TenantID = %v, want tenant1", key.TenantID)
	}
	if key.Segment1 != "20240101" {
		t.Errorf("Segment1 = %v, want 20240101", key.Segment1)
	}
}

func TestNewLevelIndexKey(t *testing.T) {
	key := NewLevelIndexKey("tenant1", "us-east-1", "error", 1234567890, "entry1")
	if key.IndexType != IndexByLevel {
		t.Errorf("IndexType = %v, want IndexByLevel", key.IndexType)
	}
	if key.Segment1 != "error" {
		t.Errorf("Segment1 = %v, want error", key.Segment1)
	}
	if key.Segment2 != "1234567890" {
		t.Errorf("Segment2 = %v, want 1234567890", key.Segment2)
	}
}

func TestNewSourceIndexKey(t *testing.T) {
	key := NewSourceIndexKey("tenant1", "us-east-1", "myapp", 1234567890, "entry1")
	if key.IndexType != IndexBySource {
		t.Errorf("IndexType = %v, want IndexBySource", key.IndexType)
	}
	if key.Segment1 != "myapp" {
		t.Errorf("Segment1 = %v, want myapp", key.Segment1)
	}
}
