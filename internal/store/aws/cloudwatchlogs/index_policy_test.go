package cloudwatchlogs

import (
	"errors"
	"testing"
)

// A read failure under the policy replacement is not policy-absence: the
// replacement aborts with the error instead of building from a phantom
// nil — the build that would overwrite the unreadable record's policy
// and punch a hole in its INACTIVE trail.
func TestReplaceIndexPolicyPropagatesReadFailure(t *testing.T) {
	s := newLogsTestStore(t)
	if err := s.putIndexPolicyStamped(&IndexPolicy{
		LogGroupName: "read-fail-group", PolicyDocument: `{"Fields":["requestId"]}`,
	}); err != nil {
		t.Fatal(err)
	}
	orig := readIndexPolicyRecord
	readIndexPolicyRecord = func(*Store, string) (*IndexPolicy, error) {
		return nil, errors.New("backend unavailable")
	}
	defer func() { readIndexPolicyRecord = orig }()

	err := s.ReplaceIndexPolicy("read-fail-group", func(current *IndexPolicy, trail *IndexInactiveTrail) (*IndexPolicy, *IndexInactiveTrail, error) {
		if current != nil {
			t.Fatal("a failed read must not surface as policy-absence (current=nil)")
		}
		return &IndexPolicy{LogGroupName: "read-fail-group", PolicyDocument: `{"Fields":["other"]}`}, trail, nil
	})
	if err == nil {
		t.Fatal("ReplaceIndexPolicy must propagate the read failure")
	}
	readIndexPolicyRecord = orig
	policy, err := s.GetIndexPolicy("read-fail-group")
	if err != nil {
		t.Fatal(err)
	}
	if policy.PolicyDocument != `{"Fields":["requestId"]}` {
		t.Fatalf("the existing policy must be untouched, got %s", policy.PolicyDocument)
	}
}
