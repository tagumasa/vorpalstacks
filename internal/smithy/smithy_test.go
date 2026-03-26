package smithy

import "testing"

func TestShapeSkipTraits(t *testing.T) {
	tests := []struct {
		trait string
		want  bool
	}{
		{TraitDocs, true},
		{TraitDocumentation, true},
		{TraitSmithyAPIInput, true},
		{TraitSmithyAPIOutput, true},
		{TraitSmithyAPIErrors, true},
		{TraitSmithyAPIHttp, true},
		{TraitSmithyAPIReadonly, true},
		{TraitSmithyAPIError, true},
		{TraitSmithyAPIHttpError, true},
		{TraitSmithyAPIHttpCode, true},
		{TraitAwsAPIErrorCode, true},
		{TraitAwsAPIService, true},
		{TraitAwsAPIVersion, true},
		{"unknown.trait", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.trait, func(t *testing.T) {
			if got := ShapeSkipTraits[tt.trait]; got != tt.want {
				t.Errorf("ShapeSkipTraits[%q] = %v, want %v", tt.trait, got, tt.want)
			}
		})
	}
}

func TestMemberSkipTraits(t *testing.T) {
	tests := []struct {
		trait string
		want  bool
	}{
		{TraitDocs, true},
		{TraitDocumentation, true},
		{TraitSmithyAPIHttpLabel, true},
		{TraitSmithyAPIHttpQuery, true},
		{TraitSmithyAPIHttpHeader, true},
		{TraitSmithyAPIHttpPayload, true},
		{TraitSmithyAPIJsonName, true},
		{TraitSmithyAPIRequired, true},
		{TraitSmithyAPIHttpPrefix, true},
		{TraitSmithyHttpDefaultPayloadMediaType, true},
		{TraitSmithyAPIHttpLocationName, true},
		{"unknown.trait", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.trait, func(t *testing.T) {
			if got := MemberSkipTraits[tt.trait]; got != tt.want {
				t.Errorf("MemberSkipTraits[%q] = %v, want %v", tt.trait, got, tt.want)
			}
		})
	}
}

func TestUsesApplicationJSON(t *testing.T) {
	tests := []struct {
		service string
		want    bool
	}{
		{"s3", true},
		{"lambda", true},
		{"dynamodb", true},
		{"ec2", true},
		{"iam", true},
		{"sns", true},
		{"sqs", true},
		{"unknown", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.service, func(t *testing.T) {
			if got := UsesApplicationJSON(tt.service); got != tt.want {
				t.Errorf("UsesApplicationJSON(%q) = %v, want %v", tt.service, got, tt.want)
			}
		})
	}
}
