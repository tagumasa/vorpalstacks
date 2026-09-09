package s3

import (
	"strings"
	"testing"
	"time"

	s3store "vorpalstacks/internal/store/aws/s3"
)

func lcPtrInt32(v int32) *int32 { return &v }

// The lifecycle validation must enforce the AWS contracts for transition
// destinations, ExpiredObjectDeleteMarker exclusivity, and the
// NewerNoncurrentVersions range and Filter requirement.
func TestValidateLifecycleRulesTransitionContract(t *testing.T) {
	filter := &LifecycleRuleFilterInput{Prefix: "docs/"}
	valid := []LifecycleRuleInput{{
		ID:     "valid",
		Status: "Enabled",
		Filter: filter,
		Transitions: []LifecycleTransitionInput{
			{Days: lcPtrInt32(1), StorageClass: string(s3store.StorageClassGlacierIR)},
		},
		NoncurrentVersionTransitions: []NoncurrentVersionTransitionInput{
			{NoncurrentDays: lcPtrInt32(1), NewerNoncurrentVersions: lcPtrInt32(1), StorageClass: string(s3store.StorageClassStandardIA)},
		},
	}}
	if err := validateLifecycleRules(valid); err != nil {
		t.Fatalf("valid transition configuration rejected: %v", err)
	}

	cases := []struct {
		name  string
		rules []LifecycleRuleInput
		want  string // error-code substring
	}{
		{
			name: "STANDARD is not a transition destination",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				Transitions: []LifecycleTransitionInput{{Days: lcPtrInt32(1), StorageClass: "STANDARD"}},
			}},
			want: "InvalidArgument",
		},
		{
			name: "REDUCED_REDUNDANCY is not a transition destination",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				Transitions: []LifecycleTransitionInput{{Days: lcPtrInt32(1), StorageClass: "REDUCED_REDUNDANCY"}},
			}},
			want: "InvalidArgument",
		},
		{
			name: "Transition StorageClass is required",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				Transitions: []LifecycleTransitionInput{{Days: lcPtrInt32(1)}},
			}},
			want: "InvalidArgument",
		},
		{
			name: "NoncurrentVersionTransition StorageClass must be a transition destination",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				NoncurrentVersionTransitions: []NoncurrentVersionTransitionInput{{NoncurrentDays: lcPtrInt32(1), StorageClass: "STANDARD"}},
			}},
			want: "InvalidArgument",
		},
		{
			name: "ExpiredObjectDeleteMarker excludes Days",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				Expiration: &LifecycleExpirationInput{Days: lcPtrInt32(1), ExpiredObjectDeleteMarker: lcPtrBool(true)},
			}},
			want: "InvalidArgument",
		},
		{
			name: "ExpiredObjectDeleteMarker excludes Date",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				Expiration: &LifecycleExpirationInput{Date: &time.Time{}, ExpiredObjectDeleteMarker: lcPtrBool(true)},
			}},
			want: "InvalidArgument",
		},
		{
			name: "ExpiredObjectDeleteMarker excludes a tag filter",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled",
				Filter:     &LifecycleRuleFilterInput{Tag: &Tag{Key: "k", Value: "v"}},
				Expiration: &LifecycleExpirationInput{ExpiredObjectDeleteMarker: lcPtrBool(true)},
			}},
			want: "InvalidArgument",
		},
		{
			name: "AbortIncompleteMultipartUpload excludes a tag filter",
			rules: []LifecycleRuleInput{{
				ID:                             "t",
				Status:                         "Enabled",
				Filter:                         &LifecycleRuleFilterInput{And: &LifecycleRuleAndOperatorInput{Tags: []Tag{{Key: "k", Value: "v"}}}},
				AbortIncompleteMultipartUpload: &AbortIncompleteUploadInput{DaysAfterInitiation: lcPtrInt32(1)},
			}},
			want: "InvalidArgument",
		},
		{
			name: "NoncurrentVersionExpiration NewerNoncurrentVersions requires a Filter",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled",
				NoncurrentVersionExpiration: &NoncurrentVersionExpirationInput{NoncurrentDays: lcPtrInt32(1), NewerNoncurrentVersions: lcPtrInt32(1)},
			}},
			want: "InvalidRequest",
		},
		{
			name: "NoncurrentVersionTransition NewerNoncurrentVersions requires a Filter",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled",
				NoncurrentVersionTransitions: []NoncurrentVersionTransitionInput{{NoncurrentDays: lcPtrInt32(1), NewerNoncurrentVersions: lcPtrInt32(1), StorageClass: "GLACIER"}},
			}},
			want: "InvalidRequest",
		},
		{
			name: "NewerNoncurrentVersions upper bound",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				NoncurrentVersionExpiration: &NoncurrentVersionExpirationInput{NoncurrentDays: lcPtrInt32(1), NewerNoncurrentVersions: lcPtrInt32(s3store.MaxNewerNoncurrentVersions + 1)},
			}},
			want: "InvalidArgument",
		},
		{
			name: "NewerNoncurrentVersions lower bound",
			rules: []LifecycleRuleInput{{
				ID: "t", Status: "Enabled", Filter: filter,
				NoncurrentVersionTransitions: []NoncurrentVersionTransitionInput{{NoncurrentDays: lcPtrInt32(1), NewerNoncurrentVersions: lcPtrInt32(0), StorageClass: "GLACIER"}},
			}},
			want: "InvalidArgument",
		},
	}
	for _, c := range cases {
		err := validateLifecycleRules(c.rules)
		if err == nil {
			t.Fatalf("%s: expected rejection", c.name)
		}
		if !strings.Contains(err.Error(), c.want) {
			t.Fatalf("%s: error %q does not carry code %s", c.name, err.Error(), c.want)
		}
	}
}

func lcPtrBool(v bool) *bool { return &v }
