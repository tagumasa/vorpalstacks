package acm

import (
	"testing"
)

// TestValidateX509Filter pins the validation/evaluation symmetry for every
// X509AttributeFilter union member: each member the evaluator honours must
// be validated here, so malformed filters fail with ValidationException
// instead of silently matching nothing.
func TestValidateX509Filter(t *testing.T) {
	tests := []struct {
		name    string
		filter  map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid Subject.CommonName filter",
			filter: map[string]interface{}{
				"Subject": map[string]interface{}{
					"CommonName": map[string]interface{}{
						"Value":              "example.com",
						"ComparisonOperator": "EQUALS",
					},
				},
			},
		},
		{
			name: "valid SubjectAlternativeName.DnsName filter",
			filter: map[string]interface{}{
				"SubjectAlternativeName": map[string]interface{}{
					"DnsName": map[string]interface{}{
						"Value":              "example.com",
						"ComparisonOperator": "CONTAINS",
					},
				},
			},
		},
		{
			name:    "valid KeyAlgorithm filter",
			filter:  map[string]interface{}{"KeyAlgorithm": "RSA_2048"},
			wantErr: false,
		},
		{
			name:    "invalid KeyAlgorithm enum value",
			filter:  map[string]interface{}{"KeyAlgorithm": "RSA_1024_LEGACY"},
			wantErr: true,
		},
		{
			name:    "valid KeyUsage filter",
			filter:  map[string]interface{}{"KeyUsage": "DIGITAL_SIGNATURE"},
			wantErr: false,
		},
		{
			name:    "invalid KeyUsage enum value",
			filter:  map[string]interface{}{"KeyUsage": "NOT_A_KEY_USAGE"},
			wantErr: true,
		},
		{
			name:    "wrong-type KeyUsage",
			filter:  map[string]interface{}{"KeyUsage": 42},
			wantErr: true,
		},
		{
			name:    "valid ExtendedKeyUsage filter",
			filter:  map[string]interface{}{"ExtendedKeyUsage": "TLS_WEB_SERVER_AUTHENTICATION"},
			wantErr: false,
		},
		{
			name:    "invalid ExtendedKeyUsage enum value",
			filter:  map[string]interface{}{"ExtendedKeyUsage": "SERVER_AUTH"},
			wantErr: true,
		},
		{
			name:    "valid SerialNumber filter",
			filter:  map[string]interface{}{"SerialNumber": "e5:87:ef:34:7a:4a:0f:de"},
			wantErr: false,
		},
		{
			name:    "wrong-type SerialNumber",
			filter:  map[string]interface{}{"SerialNumber": 123},
			wantErr: true,
		},
		{
			name:    "SerialNumber violating the hex-pair pattern",
			filter:  map[string]interface{}{"SerialNumber": "not-hex"},
			wantErr: true,
		},
		{
			name: "valid NotAfter TimestampRange",
			filter: map[string]interface{}{
				"NotAfter": map[string]interface{}{"Start": float64(1700000000)},
			},
			wantErr: false,
		},
		{
			name:    "wrong-type NotAfter",
			filter:  map[string]interface{}{"NotAfter": "2026-01-01T00:00:00Z"},
			wantErr: true,
		},
		{
			name: "NotAfter with non-timestamp Start",
			filter: map[string]interface{}{
				"NotAfter": map[string]interface{}{"Start": "2026-01-01T00:00:00Z"},
			},
			wantErr: true,
		},
		{
			name: "valid NotBefore TimestampRange",
			filter: map[string]interface{}{
				"NotBefore": map[string]interface{}{"End": float64(1800000000)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		err := validateX509Filter(tt.filter)
		if tt.wantErr && err == nil {
			t.Errorf("%s: expected validation error, got nil", tt.name)
		}
		if !tt.wantErr && err != nil {
			t.Errorf("%s: unexpected validation error: %v", tt.name, err)
		}
	}
}
