package cloudtrail

import (
	"reflect"
	"testing"
)

// TestApplyTagsAcceptsBothWireForms pins the shared tag applier: the
// JSON-protocol list form and the JSON string holding the same list must both
// populate the destination map, a nil map must be initialised, and entries
// with an empty key must be ignored. Every tag-bearing CloudTrail resource
// routes through this single applier.
func TestApplyTagsAcceptsBothWireForms(t *testing.T) {
	listForm := []interface{}{
		map[string]interface{}{"Key": "env", "Value": "prod"},
		map[string]interface{}{"Key": "", "Value": "ignored"},
	}
	stringForm := `[{"Key":"env","Value":"prod"}]`

	var fromList map[string]string
	applyTags(&fromList, listForm)
	if !reflect.DeepEqual(fromList, map[string]string{"env": "prod"}) {
		t.Fatalf("applyTags(list) = %v, want map[env:prod]", fromList)
	}

	var fromString map[string]string
	applyTags(&fromString, stringForm)
	if !reflect.DeepEqual(fromString, map[string]string{"env": "prod"}) {
		t.Fatalf("applyTags(string) = %v, want map[env:prod]", fromString)
	}

	// An existing map is extended, not replaced.
	existing := map[string]string{"keep": "me"}
	applyTags(&existing, listForm)
	if !reflect.DeepEqual(existing, map[string]string{"keep": "me", "env": "prod"}) {
		t.Fatalf("applyTags(existing) = %v, want keep+env", existing)
	}

	untouched := map[string]string{"keep": "me"}
	applyTags(&untouched, 42)
	if !reflect.DeepEqual(untouched, map[string]string{"keep": "me"}) {
		t.Fatalf("applyTags(unsupported) = %v, want unchanged", untouched)
	}
}

// TestFormatTagsListRendersWireShape pins the TagsList rendering: a list of
// {Key, Value} maps, one per tag.
func TestFormatTagsListRendersWireShape(t *testing.T) {
	got := formatTagsList(map[string]string{"env": "prod"})
	want := []interface{}{map[string]interface{}{"Key": "env", "Value": "prod"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("formatTagsList = %v, want %v", got, want)
	}
	if got := formatTagsList(nil); len(got) != 0 {
		t.Fatalf("formatTagsList(nil) = %v, want empty", got)
	}
}
