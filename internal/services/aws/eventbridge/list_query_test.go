package eventbridge

import (
	"strings"
	"testing"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// TestNormaliseListLimit pins the shared list Limit window: every list
// operation's Limit member targets the LimitMax100 shape (@range 1-100,
// vendored model); an omitted Limit (the int32 zero value) defaults to the
// maximum, never to a smaller silent default.
func TestNormaliseListLimit(t *testing.T) {
	cases := []struct {
		name    string
		limit   int32
		want    int32
		wantErr bool
	}{
		{"omitted defaults to the maximum", 0, 100, false},
		{"minimum accepted", 1, 1, false},
		{"maximum accepted", 100, 100, false},
		{"negative rejected", -1, 0, true},
		{"above maximum rejected", 101, 0, true},
	}
	for _, tc := range cases {
		got, err := normaliseListLimit(tc.limit)
		if tc.wantErr {
			if err == nil {
				t.Errorf("%s: accepted, want rejection", tc.name)
			} else if !strings.Contains(err.Error(), "Limit must be between 1 and 100") {
				t.Errorf("%s: error = %v, want the 1-100 window message", tc.name, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tc.name, err)
			continue
		}
		if got != tc.want {
			t.Errorf("%s: limit = %d, want %d", tc.name, got, tc.want)
		}
	}
}

// TestListNamePrefixValidation pins the NamePrefix contract across every
// family: the member targets the family's name shape (jq over the vendored
// model), so a prefix outside the charset or length window is a
// ValidationException rather than a silent empty match.
func TestListNamePrefixValidation(t *testing.T) {
	valid := map[string]string{
		"event-bus":       "bus.",
		"rule":            "rule-",
		"archive":         "arch_",
		"connection":      "conn",
		"api-destination": "dest",
		"replay":          "rep",
	}
	for kind, prefix := range valid {
		if err := validateListNamePrefix(prefix, kind); err != nil {
			t.Errorf("%s: valid prefix %q rejected: %v", kind, prefix, err)
		}
	}

	invalid := map[string]string{
		"event-bus":       "bad prefix!", // space and '!' outside the charset
		"rule":            "bad prefix!",
		"archive":         "bad prefix!",
		"connection":      "bad prefix!",
		"api-destination": "bad prefix!",
		"replay":          "bad prefix!",
	}
	for kind, prefix := range invalid {
		if err := validateListNamePrefix(prefix, kind); err == nil {
			t.Errorf("%s: invalid prefix %q accepted, want rejection", kind, prefix)
		}
	}

	// The event-bus charset admits '/' (its name shape does); the others
	// do not.
	if err := validateListNamePrefix("team/", "event-bus"); err != nil {
		t.Errorf("event-bus prefix with slash rejected: %v", err)
	}
	if err := validateListNamePrefix("team/", "rule"); err == nil {
		t.Error("rule prefix with slash accepted, want rejection")
	}

	// Length windows follow the family's name shape: 49 characters is a
	// valid archive prefix, 49 'a's as a prefix exceed nothing (the
	// archive window is 48) — the over-window prefix is rejected.
	if err := validateListNamePrefix(strings.Repeat("a", 48), "archive"); err != nil {
		t.Errorf("48-character archive prefix rejected: %v", err)
	}
	if err := validateListNamePrefix(strings.Repeat("a", 49), "archive"); err == nil {
		t.Error("49-character archive prefix accepted, want rejection (ArchiveName @length max 48)")
	}
}

// TestStateFilterVocabularies pins the list-filter enum validation against
// the complete model vocabularies: every model member is accepted, a
// non-member is a ValidationException naming the vocabulary, and the
// fabricated DELETING archive state (not a model member) is rejected.
func TestStateFilterVocabularies(t *testing.T) {
	store := newEventBusCoreTestStore(t)
	svc := newEventBusCoreTestService()
	ctx := t.Context()

	for _, state := range eventsstore.ArchiveStateVocabulary() {
		if _, err := svc.listArchivesCore(ctx, store, ListArchivesInput{State: state}); err != nil {
			t.Errorf("ListArchives State=%s rejected: %v", state, err)
		}
	}
	if _, err := svc.listArchivesCore(ctx, store, ListArchivesInput{State: "DELETING"}); err == nil {
		t.Error("ListArchives State=DELETING accepted; it is not an ArchiveState model member")
	} else if !strings.Contains(err.Error(), "ArchiveState enum") {
		t.Errorf("ListArchives State=DELETING error = %v, want the enum message", err)
	}

	for _, state := range eventsstore.ConnectionStateVocabulary() {
		if _, err := svc.listConnectionsCore(ctx, store, ListConnectionsInput{State: state}); err != nil {
			t.Errorf("ListConnections State=%s rejected: %v", state, err)
		}
	}
	if _, err := svc.listConnectionsCore(ctx, store, ListConnectionsInput{State: "BOGUS"}); err == nil {
		t.Error("ListConnections State=BOGUS accepted, want rejection")
	}

	for _, state := range eventsstore.ReplayStateVocabulary() {
		if _, err := svc.listReplaysCore(ctx, store, ListReplaysInput{State: eventsstore.ReplayState(state)}); err != nil {
			t.Errorf("ListReplays State=%s rejected: %v", state, err)
		}
	}
	if _, err := svc.listReplaysCore(ctx, store, ListReplaysInput{State: eventsstore.ReplayState("BOGUS")}); err == nil {
		t.Error("ListReplays State=BOGUS accepted, want rejection")
	}
}

// TestStateVocabulariesMatchModel pins the store's state vocabularies
// against the vendored Smithy enums (2015-10-07), regenerated at fix time:
// ArchiveState has six members, ConnectionState nine, ReplayState six.
func TestStateVocabulariesMatchModel(t *testing.T) {
	want := map[string][]string{
		"ArchiveState":    []string{"CREATE_FAILED", "CREATING", "DISABLED", "ENABLED", "UPDATE_FAILED", "UPDATING"},
		"ConnectionState": []string{"ACTIVE", "AUTHORIZED", "AUTHORIZING", "CREATING", "DEAUTHORIZED", "DEAUTHORIZING", "DELETING", "FAILED_CONNECTIVITY", "UPDATING"},
		"ReplayState":     []string{"CANCELLED", "CANCELLING", "COMPLETED", "FAILED", "RUNNING", "STARTING"},
	}
	got := map[string][]string{
		"ArchiveState":    eventsstore.ArchiveStateVocabulary(),
		"ConnectionState": eventsstore.ConnectionStateVocabulary(),
		"ReplayState":     eventsstore.ReplayStateVocabulary(),
	}
	for enum, wantValues := range want {
		if strings.Join(got[enum], ",") != strings.Join(wantValues, ",") {
			t.Errorf("%s vocabulary = %v, want %v (Smithy model 2015-10-07)", enum, got[enum], wantValues)
		}
	}
}
