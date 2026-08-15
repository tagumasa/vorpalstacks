package iot

import "testing"

// restJson1 serialises Timestamp members as JSON numbers; the suppression
// handlers must read the raw member instead of a string-only accessor
// that silently turns numbers into zero.
func TestTimestampMemberParam(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		want   int64
	}{
		{"json number", map[string]interface{}{"expirationDate": float64(1735689600)}, 1735689600},
		{"int", map[string]interface{}{"expirationDate": 1735689600}, 1735689600},
		{"epoch string", map[string]interface{}{"expirationDate": "1735689600"}, 1735689600},
		{"iso string", map[string]interface{}{"expirationDate": "2025-01-01T00:00:00Z"}, 1735689600},
		{"absent", map[string]interface{}{}, 0},
		{"garbage string", map[string]interface{}{"expirationDate": "soon"}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := timestampMemberParam(tt.params, "expirationDate"); got != tt.want {
				t.Fatalf("timestampMemberParam = %d, want %d", got, tt.want)
			}
		})
	}
}
