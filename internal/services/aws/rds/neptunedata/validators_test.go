package neptunedata

import "testing"

func TestValidateLoaderFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"csv", "csv", true},
		{"opencypher", "opencypher", true},
		{"ntriples", "ntriples", true},
		{"nquads", "nquads", true},
		{"rdfxml", "rdfxml", true},
		{"turtle", "turtle", true},
		{"uppercase CSV", "CSV", false},
		{"unknown", "json", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateLoaderFormat(tt.input); got != tt.want {
				t.Errorf("validateLoaderFormat(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateLoaderMode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"RESUME", "RESUME", true},
		{"NEW", "NEW", true},
		{"AUTO", "AUTO", true},
		{"empty (optional)", "", true},
		{"lowercase", "resume", false},
		{"unknown", "REPLACE", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateLoaderMode(tt.input); got != tt.want {
				t.Errorf("validateLoaderMode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateLoaderParallelism(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"LOW", "LOW", true},
		{"MEDIUM", "MEDIUM", true},
		{"HIGH", "HIGH", true},
		{"OVERSUBSCRIBE", "OVERSUBSCRIBE", true},
		{"empty (optional)", "", true},
		{"lowercase", "low", false},
		{"unknown", "EXTREME", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateLoaderParallelism(tt.input); got != tt.want {
				t.Errorf("validateLoaderParallelism(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateExplainMode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"static", "static", true},
		{"details", "details", true},
		{"dynamic", "dynamic", true},
		{"empty (optional)", "", true},
		{"uppercase", "STATIC", false},
		{"unknown", "verbose", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateExplainMode(tt.input); got != tt.want {
				t.Errorf("validateExplainMode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateIteratorType(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"AT_SEQUENCE_NUMBER", "AT_SEQUENCE_NUMBER", true},
		{"AFTER_SEQUENCE_NUMBER", "AFTER_SEQUENCE_NUMBER", true},
		{"LATEST", "LATEST", true},
		{"TRIM_HORIZON", "TRIM_HORIZON", true},
		{"empty (optional)", "", true},
		{"lowercase", "latest", false},
		{"unknown", "EARLIEST", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateIteratorType(tt.input); got != tt.want {
				t.Errorf("validateIteratorType(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateManageStatsMode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"disableAutoCompute", "disableAutoCompute", true},
		{"enableAutoCompute", "enableAutoCompute", true},
		{"refresh", "refresh", true},
		{"empty", "", false},
		{"uppercase", "REFRESH", false},
		{"unknown", "reset", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateManageStatsMode(tt.input); got != tt.want {
				t.Errorf("validateManageStatsMode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateGraphSummaryMode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"basic", "basic", true},
		{"detailed", "detailed", true},
		{"empty (optional)", "", true},
		{"uppercase", "DETAILED", false},
		{"unknown", "summary", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateGraphSummaryMode(tt.input); got != tt.want {
				t.Errorf("validateGraphSummaryMode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateIamRoleArn(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"valid role", "arn:aws:iam::123456789012:role/NeptuneLoadRole", true},
		{"valid path role", "arn:aws:iam::123456789012:role/service-role/MyRole", true},
		{"empty account id", "arn:aws:iam:::role/NeptuneLoadRole", true},
		{"s3 arn", "arn:aws:s3:::mybucket", false},
		{"kms arn", "arn:aws:kms:us-east-1:123456789012:key/abc", false},
		{"empty", "", false},
		{"plain string", "not-an-arn", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateIamRoleArn(tt.input); got != tt.want {
				t.Errorf("validateIamRoleArn(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateS3SourceURI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"s3 uri", "s3://mybucket/data/", true},
		{"https s3", "https://s3.us-east-1.amazonaws.com/bucket/key", true},
		{"file uri", "file:///tmp/data.csv", false},
		{"empty", "", false},
		{"plain path", "/tmp/data.csv", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateS3SourceURI(tt.input); got != tt.want {
				t.Errorf("validateS3SourceURI(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
