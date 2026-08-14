package ssm

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return NewStore(st, "123456789012", "us-east-1")
}

// TestGetParameterHistory_OrderingAcrossPages verifies that history pages
// are strictly newest-version-first across page boundaries. History keys
// encode the version in decimal, so lexicographic key order (v10 before v2)
// must never leak into the output.
func TestGetParameterHistory_OrderingAcrossPages(t *testing.T) {
	store := newTestStore(t)

	const totalVersions = 12
	param := NewParameter("order-param", "v0", ParameterTypeString)
	for i := 1; i <= totalVersions; i++ {
		p := *param
		p.Value = fmt.Sprintf("value-%d", i)
		if _, err := store.PutParameter(&p, true); err != nil {
			t.Fatalf("put version %d: %v", i, err)
		}
	}

	var got []int64
	marker := ""
	pages := 0
	for {
		items, nextMarker, truncated, err := store.GetParameterHistory("order-param", 5, marker)
		if err != nil {
			t.Fatalf("get history page: %v", err)
		}
		for _, it := range items {
			got = append(got, it.Version)
		}
		pages++
		if !truncated {
			break
		}
		marker = nextMarker
		if pages > 10 {
			t.Fatal("pagination did not terminate")
		}
	}

	if len(got) != totalVersions {
		t.Fatalf("got %d versions, want %d", len(got), totalVersions)
	}
	want := int64(totalVersions)
	for _, v := range got {
		if v != want {
			t.Fatalf("versions not strictly descending: got sequence %v", got)
		}
		want--
	}
}

// TestGetParameterHistory_InvalidMarker verifies that a marker which does
// not parse as a positive version number is rejected.
func TestGetParameterHistory_InvalidMarker(t *testing.T) {
	store := newTestStore(t)
	param := NewParameter("marker-param", "v1", ParameterTypeString)
	if _, err := store.PutParameter(param, false); err != nil {
		t.Fatalf("put: %v", err)
	}

	if _, _, _, err := store.GetParameterHistory("marker-param", 10, "not-a-version"); !errors.Is(err, ErrInvalidNextToken) {
		t.Fatalf("err = %v, want ErrInvalidNextToken", err)
	}
	if _, _, _, err := store.GetParameterHistory("marker-param", 10, "0"); !errors.Is(err, ErrInvalidNextToken) {
		t.Fatalf("err = %v, want ErrInvalidNextToken for zero marker", err)
	}
}

// TestGetParameterHistory_PoliciesInHistory verifies that the Policies
// document is copied into each history entry so GetParameterHistory can
// return it (the Smithy ParameterHistory shape includes Policies).
func TestGetParameterHistory_PoliciesInHistory(t *testing.T) {
	store := newTestStore(t)
	param := NewParameter("policy-param", "v1", ParameterTypeString)
	param.Policies = `[{"Type":"Expiration","Version":"1"}]`
	if _, err := store.PutParameter(param, false); err != nil {
		t.Fatalf("put: %v", err)
	}

	items, _, _, err := store.GetParameterHistory("policy-param", 10, "")
	if err != nil {
		t.Fatalf("get history: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}
	if items[0].Policies != param.Policies {
		t.Fatalf("history Policies = %q, want %q", items[0].Policies, param.Policies)
	}
}

// TestValidateParameterName_UserSpecifiableLength verifies the AWS
// effective limit: the documented 2048-character maximum includes 1037
// characters reserved for internal use, so callers may specify at most 1011.
func TestValidateParameterName_UserSpecifiableLength(t *testing.T) {
	if err := ValidateParameterName(strings.Repeat("a", MaxParameterNameLength)); err != nil {
		t.Fatalf("1011-char name rejected: %v", err)
	}
	if err := ValidateParameterName(strings.Repeat("a", MaxParameterNameLength+1)); !errors.Is(err, ErrInvalidParameterName) {
		t.Fatalf("1012-char name err = %v, want ErrInvalidParameterName", err)
	}
}

// TestValidateParameterName_HierarchyDepth verifies the fifteen-level
// hierarchy depth cap ("/a/b/c" is three levels).
func TestValidateParameterName_HierarchyDepth(t *testing.T) {
	ok := "/" + strings.Join(repeatN("l", MaxHierarchyDepth), "/")
	if err := ValidateParameterName(ok); err != nil {
		t.Fatalf("15-level hierarchy rejected: %v", err)
	}
	tooDeep := "/" + strings.Join(repeatN("l", MaxHierarchyDepth+1), "/")
	if err := ValidateParameterName(tooDeep); !errors.Is(err, ErrHierarchyLevelLimitExceeded) {
		t.Fatalf("16-level hierarchy err = %v, want ErrHierarchyLevelLimitExceeded", err)
	}
}

func repeatN(s string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = s
	}
	return out
}
