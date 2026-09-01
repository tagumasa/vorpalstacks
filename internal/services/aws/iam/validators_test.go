package iam

import (
	"strings"
	"testing"

	tagutil "vorpalstacks/internal/common/tags"
)

// TestValidateTagEntriesUnicodeLengths pins that tag key and value limits
// are counted in Unicode characters, matching the Smithy tagKeyType and
// tagValueType @length traits.
func TestValidateTagEntriesUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if err := validateTagEntries([]tagutil.Tag{{Key: strings.Repeat(cjk, 100), Value: strings.Repeat(cjk, 200)}}); err != nil {
		t.Errorf("100-character CJK key with 200-character value rejected: %v", err)
	}
	if err := validateTagEntries([]tagutil.Tag{{Key: strings.Repeat(cjk, 129), Value: "v"}}); err == nil {
		t.Error("129-character CJK key accepted, want rejection")
	}
	if err := validateTagEntries([]tagutil.Tag{{Key: "k", Value: strings.Repeat(cjk, 257)}}); err == nil {
		t.Error("257-character CJK value accepted, want rejection")
	}
}

// TestValidateRoleDescriptionLatin1 pins that roleDescriptionType lengths
// are counted in Unicode characters: the shape's pattern admits Latin-1
// supplement characters (2 bytes each in UTF-8), so 600 such characters
// (1200 bytes) fall within the 1000-character bound.
func TestValidateRoleDescriptionLatin1(t *testing.T) {
	latin1 := "é" // U+00E9, inside the pattern's \u00A1-\u00FF range

	if !validateRoleDescription(strings.Repeat(latin1, 600)) {
		t.Error("600-character Latin-1 role description rejected")
	}
	if validateRoleDescription(strings.Repeat(latin1, 1001)) {
		t.Error("1001-character Latin-1 role description accepted")
	}
	if !validateRoleDescription(strings.Repeat("a", 1000)) {
		t.Error("1000-character ASCII role description rejected")
	}
	if validateRoleDescription("✓") {
		t.Error("character outside the Latin-1 pattern accepted")
	}
}

// TestValidateClientIDUnicodeLengths pins that clientIDType @length(1,255)
// is counted in Unicode characters; the shape carries no pattern, so
// rune-legal multibyte client IDs must not be rejected on byte length.
func TestValidateClientIDUnicodeLengths(t *testing.T) {
	cjk := "\u65e5" // one CJK character, 3 bytes

	if !validateClientID(strings.Repeat(cjk, 255)) {
		t.Error("255-character CJK client ID rejected")
	}
	if validateClientID(strings.Repeat(cjk, 256)) {
		t.Error("256-character CJK client ID accepted")
	}
	if validateClientID("") {
		t.Error("empty client ID accepted")
	}
}

// TestGetEntityCoresRejectEmptyName pins the shared empty-identifier
// rejection in the user/role/policy/group get cores: an omitted member is
// a client error on both planes, reported before any store access.
func TestGetEntityCoresRejectEmptyName(t *testing.T) {
	svc := &IAMService{}

	cases := []struct {
		name string
		call func() error
		want string
	}{
		{"user", func() error { _, err := svc.getUserCore(nil, ""); return err }, "Required parameter UserName is missing."},
		{"role", func() error { _, err := svc.getRoleCore(nil, ""); return err }, "Required parameter RoleName is missing."},
		{"policy", func() error { _, err := svc.getPolicyCore(nil, ""); return err }, "Required parameter PolicyArn is missing."},
		{"group", func() error { _, err := svc.getGroupCore(nil, ""); return err }, "Required parameter GroupName is missing."},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Fatal("expected a validation error, got nil")
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error %q does not contain %q", err.Error(), tc.want)
			}
		})
	}
}

// TestPutPermissionsBoundaryCoresValidationOrder pins the member precedence
// of the permissions-boundary operation cores: when a request omits both the
// entity name and the boundary ARN, the name is reported first, and a
// present name with an omitted boundary reports the boundary. Both checks
// fire before any store access.
func TestPutPermissionsBoundaryCoresValidationOrder(t *testing.T) {
	svc := &IAMService{}

	t.Run("role both omitted reports RoleName", func(t *testing.T) {
		err := svc.putRolePermissionsBoundaryCore(nil, "", "")
		if err == nil || !strings.Contains(err.Error(), "Required parameter RoleName is missing.") {
			t.Fatalf("expected the RoleName validation error, got %v", err)
		}
	})
	t.Run("role boundary omitted", func(t *testing.T) {
		err := svc.putRolePermissionsBoundaryCore(nil, "ops", "")
		if err == nil || !strings.Contains(err.Error(), "Required parameter PermissionsBoundary is missing.") {
			t.Fatalf("expected the PermissionsBoundary validation error, got %v", err)
		}
	})
	t.Run("user both omitted reports UserName", func(t *testing.T) {
		err := svc.putUserPermissionsBoundaryCore(nil, "", "")
		if err == nil || !strings.Contains(err.Error(), "Required parameter UserName is missing.") {
			t.Fatalf("expected the UserName validation error, got %v", err)
		}
	})
	t.Run("user boundary omitted", func(t *testing.T) {
		err := svc.putUserPermissionsBoundaryCore(nil, "alice", "")
		if err == nil || !strings.Contains(err.Error(), "Required parameter PermissionsBoundary is missing.") {
			t.Fatalf("expected the PermissionsBoundary validation error, got %v", err)
		}
	})
}

// TestDeletePermissionsBoundaryCoresRejectEmptyName pins the empty-name
// rejection in the permissions-boundary delete cores: it fires as a client
// error before any store access.
func TestDeletePermissionsBoundaryCoresRejectEmptyName(t *testing.T) {
	svc := &IAMService{}

	if err := svc.deleteRolePermissionsBoundaryCore(nil, ""); err == nil ||
		!strings.Contains(err.Error(), "Required parameter RoleName is missing.") {
		t.Fatalf("expected the RoleName validation error, got %v", err)
	}
	if err := svc.deleteUserPermissionsBoundaryCore(nil, ""); err == nil ||
		!strings.Contains(err.Error(), "Required parameter UserName is missing.") {
		t.Fatalf("expected the UserName validation error, got %v", err)
	}
}
