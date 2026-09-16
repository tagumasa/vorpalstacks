package eventbridge

import (
	"fmt"
	"strings"
	"testing"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

func inputConfig(target eventsstore.Target) *eventsstore.Target {
	return &target
}

// The input configuration is validated at PutTargets: mutual exclusion of
// the three members, JSON validity of Input, and the modelled length,
// entry-count and pattern bounds of the transformer.
func TestValidateTargetInputConfiguration(t *testing.T) {
	longInput := `{"pad":"` + strings.Repeat("x", eventsstore.TargetInputMaxLength) + `"}`
	longPath := "$." + strings.Repeat("a", eventsstore.TargetInputPathMaxLength+1)
	overHundred := make(map[string]string, eventsstore.InputPathsMapMaxEntries+1)
	for i := 0; i <= eventsstore.InputPathsMapMaxEntries; i++ {
		overHundred[fmt.Sprintf("k%d", i)] = "$.detail"
	}

	cases := []struct {
		name    string
		target  *eventsstore.Target
		wantErr string
	}{
		{
			name:    "no input configuration passes",
			target:  inputConfig(eventsstore.Target{}),
			wantErr: "",
		},
		{
			name:    "Input alone passes",
			target:  inputConfig(eventsstore.Target{Input: `{"action":"test"}`}),
			wantErr: "",
		},
		{
			name:    "string constant Input passes",
			target:  inputConfig(eventsstore.Target{Input: `"HelloWorld!"`}),
			wantErr: "",
		},
		{
			name:    "InputPath alone passes",
			target:  inputConfig(eventsstore.Target{InputPath: "$.detail"}),
			wantErr: "",
		},
		{
			name:    "InputTransformer alone passes",
			target:  inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{InputTemplate: `"t <v>"`, InputPathsMap: map[string]string{"v": "$.detail.v"}}}),
			wantErr: "",
		},
		{
			name:    "Input with InputPath rejected",
			target:  inputConfig(eventsstore.Target{Input: `{}`, InputPath: "$.detail"}),
			wantErr: "mutually exclusive",
		},
		{
			name:    "Input with InputTransformer rejected",
			target:  inputConfig(eventsstore.Target{Input: `{}`, InputTransformer: &eventsstore.InputTransformer{InputTemplate: "t"}}),
			wantErr: "mutually exclusive",
		},
		{
			name:    "InputPath with InputTransformer rejected",
			target:  inputConfig(eventsstore.Target{InputPath: "$.detail", InputTransformer: &eventsstore.InputTransformer{InputTemplate: "t"}}),
			wantErr: "mutually exclusive",
		},
		{
			name:    "Input must be valid JSON",
			target:  inputConfig(eventsstore.Target{Input: `{not json`}),
			wantErr: "valid JSON",
		},
		{
			name:    "Input length bound",
			target:  inputConfig(eventsstore.Target{Input: longInput}),
			wantErr: "at most 8192",
		},
		{
			// @length counts Unicode scalar values, so a multibyte Input
			// under the bound passes even though its byte size exceeds it.
			name:    "multibyte Input under the bound passes",
			target:  inputConfig(eventsstore.Target{Input: `{"pad":"` + strings.Repeat("あ", 5000) + `"}`}),
			wantErr: "",
		},
		{
			name:    "multibyte Input over the bound rejected",
			target:  inputConfig(eventsstore.Target{Input: `{"pad":"` + strings.Repeat("あ", eventsstore.TargetInputMaxLength) + `"}`}),
			wantErr: "at most 8192",
		},
		{
			name:    "InputPath length bound",
			target:  inputConfig(eventsstore.Target{InputPath: longPath}),
			wantErr: "at most 256",
		},
		{
			name:    "empty InputTemplate rejected",
			target:  inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{}}),
			wantErr: "between 1 and 8192",
		},
		{
			name: "InputTemplate length bound",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: strings.Repeat("a", eventsstore.InputTemplateMaxLength+1),
			}}),
			wantErr: "between 1 and 8192",
		},
		{
			// The template bound is scalar values too: 5,000 multibyte
			// runes are 15,000 bytes yet under the 8,192-character bound.
			name: "multibyte InputTemplate under the bound passes",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: strings.Repeat("あ", 5000),
			}}),
			wantErr: "",
		},
		{
			name: "InputPathsMap entry-count bound",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: "t",
				InputPathsMap: overHundred,
			}}),
			wantErr: "at most 100 entries",
		},
		{
			name: "InputPathsMap key pattern rejects dots",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: "t",
				InputPathsMap: map[string]string{"aws.events.rule-arn": "$.detail"},
			}}),
			wantErr: "must match [A-Za-z0-9_-]",
		},
		{
			name: "InputPathsMap key pattern rejects spaces",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: "t",
				InputPathsMap: map[string]string{"my key": "$.detail"},
			}}),
			wantErr: "must match [A-Za-z0-9_-]",
		},
		{
			name: "InputPathsMap key length bound",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: "t",
				InputPathsMap: map[string]string{strings.Repeat("k", eventsstore.InputPathsMapKeyMaxLength+1): "$.detail"},
			}}),
			wantErr: "at most 256",
		},
		{
			name: "InputPathsMap value length bound",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: "t",
				InputPathsMap: map[string]string{"k": longPath},
			}}),
			wantErr: "at most 256",
		},
		{
			name: "placeholder as object key rejected",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: `{"<key>": "value"}`,
				InputPathsMap: map[string]string{"key": "$.detail"},
			}}),
			wantErr: "object key",
		},
		{
			name: "placeholder as object value passes",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: `{"key": <value>}`,
				InputPathsMap: map[string]string{"value": "$.detail"},
			}}),
			wantErr: "",
		},
		{
			name: "quoted placeholder in a string template passes",
			target: inputConfig(eventsstore.Target{InputTransformer: &eventsstore.InputTransformer{
				InputTemplate: `"<key> matters"`,
				InputPathsMap: map[string]string{"key": "$.detail"},
			}}),
			wantErr: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTargetInputConfiguration(tc.target)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("expected acceptance, got: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got: %v", tc.wantErr, err)
			}
		})
	}
}
