package lambda

import "testing"

// TestCanonicalRuntime pins the canonical-form contract of the runtime
// single source: case variants normalise to the lowercase model value and
// identifiers outside the current set — including unmodelled and
// deprecated values — are rejected.
func TestCanonicalRuntime(t *testing.T) {
	cases := []struct {
		in   string
		want string
		ok   bool
	}{
		{"python3.12", "python3.12", true},
		{"Python3.12", "python3.12", true},
		{"NODEJS22.X", "nodejs22.x", true},
		{"Java17.AL2023", "java17.al2023", true},
		{"dotnet9", "", false},    // unmodelled
		{"nodejs20.x", "", false}, // deprecated, not creatable
		{"nodejs26.x", "", false}, // non-public feature tag
		{"bogus", "", false},
		{"", "", false},
	}
	for _, tc := range cases {
		got, ok := CanonicalRuntime(tc.in)
		if ok != tc.ok || string(got) != tc.want {
			t.Errorf("CanonicalRuntime(%q) = (%q, %v), want (%q, %v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}

// TestCurrentRuntimesHaveImages pins the invariant the invocation plane
// relies on: every runtime accepted for new functions has a mapped
// container image, so validation can never admit a runtime that would
// fail the image lookup.
func TestCurrentRuntimesHaveImages(t *testing.T) {
	for _, r := range CurrentRuntimes {
		if _, ok := RuntimeImageMapping[r]; !ok {
			t.Errorf("current runtime %q has no image mapping", r)
		}
	}
}

// TestGetImageForRuntimeNoFallback pins that an unmapped runtime is
// reported as a miss: the silent provided:al2 substitution executed the
// wrong runtime with no error anywhere, so the lookup must make the miss
// explicit for the caller to surface.
func TestGetImageForRuntimeNoFallback(t *testing.T) {
	if _, ok := GetImageForRuntime("Python3.12"); ok {
		t.Fatal("mixed-case runtime must be an unmapped miss, not a fallback image")
	}
	if _, ok := GetImageForRuntime("bogus"); ok {
		t.Fatal("unknown runtime must be an unmapped miss, not a fallback image")
	}
	if image, ok := GetImageForRuntime(RuntimePython312); !ok || image == "" {
		t.Fatalf("mapped runtime must resolve, got (%q, %v)", image, ok)
	}
}
