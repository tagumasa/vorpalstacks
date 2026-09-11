package iam

import (
	"errors"
	"testing"
)

// The catalogue is versioned data: a template ARN can carry several minor
// versions, and Get resolves the requested one among them — the default
// selector (minorVersion 0) picks the template's default minor version.
// The fixtures append a second minor version to the seeded Example
// template for the duration of each test and restore the catalogue after.
func TestRoleTemplateStoreResolvesMinorVersions(t *testing.T) {
	restore := withExtraMinorVersion(t, 2, 1)
	defer restore()

	store := NewRoleTemplateStore()
	const arn = "arn:aws:iam::aws:role-template/awsserviceprincipal/Example:1"

	first, err := store.Get(arn, 1)
	if err != nil {
		t.Fatalf("resolve minor 1: %v", err)
	}
	if first.MinorVersion != 1 {
		t.Fatalf("minor 1 lookup: got minor %d", first.MinorVersion)
	}

	second, err := store.Get(arn, 2)
	if err != nil {
		t.Fatalf("resolve minor 2: %v", err)
	}
	if second.MinorVersion != 2 {
		t.Fatalf("minor 2 lookup: got minor %d", second.MinorVersion)
	}

	def, err := store.Get(arn, 0)
	if err != nil {
		t.Fatalf("resolve default minor: %v", err)
	}
	if def.MinorVersion != 1 {
		t.Fatalf("default minor must stay 1, got %d", def.MinorVersion)
	}

	if _, err := store.Get(arn, 99); !errors.Is(err, ErrRoleTemplateVersionNotFound) {
		t.Fatalf("unknown minor: got %v, want ErrRoleTemplateVersionNotFound", err)
	}
	if _, err := store.Get("arn:aws:iam::aws:role-template/awsserviceprincipal/Absent:1", 0); !errors.Is(err, ErrRoleTemplateNotFound) {
		t.Fatalf("unknown template: got %v, want ErrRoleTemplateNotFound", err)
	}
}

// withExtraMinorVersion appends a copy of the seeded Example template at
// the given minor version and returns a function restoring the catalogue.
func withExtraMinorVersion(t *testing.T, minorVersion, defaultMinor int) func() {
	t.Helper()
	const arn = "arn:aws:iam::aws:role-template/awsserviceprincipal/Example:1"
	for i := range roleTemplateCatalogue {
		if roleTemplateCatalogue[i].TemplateArn == arn {
			extra := roleTemplateCatalogue[i]
			extra.MinorVersion = minorVersion
			extra.DefaultMinorVersion = defaultMinor
			roleTemplateCatalogue = append(roleTemplateCatalogue, extra)
			original := roleTemplateCatalogue
			return func() { roleTemplateCatalogue = original[:len(original)-1] }
		}
	}
	t.Fatal("seeded Example template not found")
	return nil
}
