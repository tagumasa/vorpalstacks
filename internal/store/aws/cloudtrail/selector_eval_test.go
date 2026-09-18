package cloudtrail

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAdvancedSelectorMatchesDocumentationMatrix pins the user guide's
// "Example showing multiple conditions for the resources.ARN field" as the
// oracle: one selector carrying all six operator lists on resources.ARN,
// evaluated against the documentation's six objects. object1 matches
// Equals, object2 StartsWith and object3 EndsWith (the SELECT group OR'd);
// object4, object5 and object6 are each excluded by the single DESELECT
// condition they meet despite satisfying StartsWith (the DESELECT group
// AND'd with it).
func TestAdvancedSelectorMatchesDocumentationMatrix(t *testing.T) {
	sel := AdvancedEventSelector{
		Name: "S3Select",
		FieldSelectors: []AdvancedFieldSelector{
			{Field: "eventCategory", Equals: []string{"Data"}},
			{Field: "resources.type", Equals: []string{"AWS::S3::Object"}},
			{
				Field:         "resources.ARN",
				Equals:        []string{"arn:aws:s3:::amzn-s3-demo-bucket/object1"},
				StartsWith:    []string{"arn:aws:s3:::amzn-s3-demo-bucket/"},
				EndsWith:      []string{"object3"},
				NotStartsWith: []string{"arn:aws:s3:::amzn-s3-demo-bucket/deselect"},
				NotEndsWith:   []string{"object5"},
				NotEquals:     []string{"arn:aws:s3:::amzn-s3-demo-bucket/object6"},
			},
		},
	}
	cases := []struct {
		arn  string
		want bool
	}{
		{"arn:aws:s3:::amzn-s3-demo-bucket/object1", true},
		{"arn:aws:s3:::amzn-s3-demo-bucket/object2", true},
		{"arn:aws:s3:::amzn-s3-demo-bucket1/object3", true},
		{"arn:aws:s3:::amzn-s3-demo-bucket/deselectObject4", false},
		{"arn:aws:s3:::amzn-s3-demo-bucket/object5", false},
		{"arn:aws:s3:::amzn-s3-demo-bucket/object6", false},
	}
	for _, tc := range cases {
		e := &Event{
			EventCategory: "Data",
			Resources:     []Resource{{ResourceType: "AWS::S3::Object", ResourceName: tc.arn}},
		}
		assert.Equal(t, tc.want, AdvancedSelectorMatches(sel, e), "arn %s", tc.arn)
	}
}

// TestAdvancedSelectorGroupSemantics pins the group semantics around the
// documentation matrix: the vacuous SELECT group of a deselect-only
// selector, the SELECT group's OR, a contradictory selector never holding,
// and the fail-closed rulings (no operator values; a field the event does
// not carry).
func TestAdvancedSelectorGroupSemantics(t *testing.T) {
	event := &Event{
		EventName:     "CreateUser",
		EventSource:   "iam.amazonaws.com",
		EventType:     "AwsApiCall",
		ReadOnly:      "false",
		EventCategory: "Management",
		UserIdentity:  &UserIdentity{Type: "IAMUser", UserName: "alice", ARN: "arn:aws:iam::acc123:user/alice"},
	}
	holds := func(fs AdvancedFieldSelector) bool {
		return AdvancedSelectorMatches(AdvancedEventSelector{
			FieldSelectors: []AdvancedFieldSelector{fs},
		}, event)
	}

	// Deselect-only: the absent SELECT group holds vacuously, so a value
	// no DESELECT condition hits qualifies.
	assert.True(t, holds(AdvancedFieldSelector{
		Field:     "eventSource",
		NotEquals: []string{"s3.amazonaws.com"},
	}))
	assert.False(t, holds(AdvancedFieldSelector{
		Field:     "eventSource",
		NotEquals: []string{"iam.amazonaws.com"},
	}))

	// The SELECT group is OR'd: either value alone qualifies.
	assert.True(t, holds(AdvancedFieldSelector{
		Field:      "eventName",
		Equals:     []string{"DeleteUser"},
		StartsWith: []string{"Create"},
	}))

	// Both groups are AND'd: a met SELECT condition does not rescue a met
	// DESELECT condition.
	assert.False(t, holds(AdvancedFieldSelector{
		Field:      "eventName",
		StartsWith: []string{"Create"},
		NotEquals:  []string{"CreateUser"},
	}))

	// A selector that contradicts itself never holds.
	assert.False(t, holds(AdvancedFieldSelector{
		Field:     "readOnly",
		Equals:    []string{"false"},
		NotEquals: []string{"false"},
	}))

	// Fail-closed: no operator values at all never holds; a field the
	// event does not carry never holds (errorCode is unset here, and no
	// event this platform records carries vpcEndpointId).
	assert.False(t, holds(AdvancedFieldSelector{Field: "eventName"}))
	assert.False(t, holds(AdvancedFieldSelector{Field: "errorCode", Equals: []string{"AccessDenied"}}))
	assert.False(t, holds(AdvancedFieldSelector{Field: "vpcEndpointId", Equals: []string{"vpce-1"}}))

	// Field comparison folds case, on both the resolver and the rules.
	assert.True(t, holds(AdvancedFieldSelector{
		Field:  "USERIDENTITY.ARN",
		Equals: []string{"arn:aws:iam::acc123:user/alice"},
	}))

	// Multi-valued fields: one qualifying resource value qualifies the
	// event; none qualifying excludes it.
	multi := &Event{
		EventCategory: "Data",
		Resources: []Resource{
			{ResourceType: "AWS::S3::Object", ResourceName: "arn:aws:s3:::b/object9"},
			{ResourceType: "AWS::DynamoDB::Table", ResourceName: "arn:aws:dynamodb:us-east-1:acc123:table/t"},
		},
	}
	assert.True(t, AdvancedSelectorMatches(AdvancedEventSelector{
		FieldSelectors: []AdvancedFieldSelector{
			{Field: "resources.type", Equals: []string{"AWS::DynamoDB::Table"}},
		},
	}, multi))
	assert.False(t, AdvancedSelectorMatches(AdvancedEventSelector{
		FieldSelectors: []AdvancedFieldSelector{
			{Field: "resources.type", Equals: []string{"AWS::Lambda::Function"}},
		},
	}, multi))
}

// TestAdvancedSelectorFieldRules pins the vocabulary and the per-field
// operator restrictions the service-layer validator enforces.
func TestAdvancedSelectorFieldRules(t *testing.T) {
	for _, field := range []string{
		"eventCategory", "eventSource", "eventName", "eventType", "readOnly",
		"sessionCredentialFromConsole", "userIdentity.arn", "resources.type",
		"resources.ARN", "errorCode", "vpcEndpointId", "username",
	} {
		_, known := AdvancedSelectorFieldRules(field)
		assert.True(t, known, "field %s must be known", field)
	}
	for _, field := range []string{"", "notAField", "resource.type", "eventcategory "} {
		_, known := AdvancedSelectorFieldRules(field)
		assert.False(t, known, "field %q must be unknown", field)
	}
	restrictedTo, known := AdvancedSelectorFieldRules("readOnly")
	require.True(t, known)
	assert.Equal(t, map[string]bool{"Equals": true}, restrictedTo)
	restrictedTo, known = AdvancedSelectorFieldRules("SessionCredentialFromConsole")
	require.True(t, known)
	assert.Equal(t, map[string]bool{"Equals": true, "NotEquals": true}, restrictedTo)
	restrictedTo, known = AdvancedSelectorFieldRules("errorCode")
	require.True(t, known)
	assert.Equal(t, map[string]bool{"Equals": true}, restrictedTo)
	unrestricted, known := AdvancedSelectorFieldRules("userIdentity.arn")
	require.True(t, known)
	assert.Nil(t, unrestricted)
}

// TestAdvancedSelectorSessionCredentialFromConsole pins the console
// field's value resolution: no event this platform writes originates
// from a console session — the record member is absent unless true — so
// the field's value is false for every event, the documented
// Equals ["false"] / NotEquals ["true"] forms select everything, and
// their complements select nothing.
func TestAdvancedSelectorSessionCredentialFromConsole(t *testing.T) {
	e := &Event{EventID: "e-1", EventName: "PutObject", EventSource: "s3.amazonaws.com"}

	equalsFalse := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "sessionCredentialFromConsole", Equals: []string{"false"}},
	}}
	assert.True(t, AdvancedSelectorMatches(equalsFalse, e))

	notEqualsTrue := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "sessionCredentialFromConsole", NotEquals: []string{"true"}},
	}}
	assert.True(t, AdvancedSelectorMatches(notEqualsTrue, e))

	equalsTrue := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "sessionCredentialFromConsole", Equals: []string{"true"}},
	}}
	assert.False(t, AdvancedSelectorMatches(equalsTrue, e))

	notEqualsFalse := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "sessionCredentialFromConsole", NotEquals: []string{"false"}},
	}}
	assert.False(t, AdvancedSelectorMatches(notEqualsFalse, e))
}

// TestAdvancedSelectorSelectsNetworkActivity pins the category detection
// the selector-level eventSource rule keys on: the eventCategory field
// with an Equals value of NetworkActivity, case-folded, and nothing else.
func TestAdvancedSelectorSelectsNetworkActivity(t *testing.T) {
	na := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "eventCategory", Equals: []string{"networkactivity"}},
	}}
	assert.True(t, AdvancedSelectorSelectsNetworkActivity(na))
	management := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "eventCategory", Equals: []string{"Management"}},
	}}
	assert.False(t, AdvancedSelectorSelectsNetworkActivity(management))
	deselected := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "eventCategory", NotEquals: []string{"NetworkActivity"}},
	}}
	assert.False(t, AdvancedSelectorSelectsNetworkActivity(deselected))
	other := AdvancedEventSelector{FieldSelectors: []AdvancedFieldSelector{
		{Field: "eventName", Equals: []string{"NetworkActivity"}},
	}}
	assert.False(t, AdvancedSelectorSelectsNetworkActivity(other))
}
