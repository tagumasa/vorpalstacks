package lambda

import (
	"errors"
	"testing"
)

// TestUpdateAliasAtomicallyChecksRevisionPrecondition pins the alias-side
// optimistic-locking precondition: a stale expected revision fails with
// ErrRevisionMismatch inside the store lock while a matching one lets the
// update through.
func TestUpdateAliasAtomicallyChecksRevisionPrecondition(t *testing.T) {
	s := resourcePolicyStore(t)

	created, err := s.CreateAliasAtomically("fn", func(fn *Function) (*Alias, error) {
		return &Alias{Name: "live", FunctionVersion: "$LATEST"}, nil
	})
	if err != nil {
		t.Fatalf("create alias: %v", err)
	}
	if created.RevisionId == "" {
		t.Fatalf("created alias carries no revision")
	}

	if _, err := s.UpdateAliasAtomically("fn", "live", "stale-revision", func(fn *Function, existing *Alias) error {
		existing.Description = "must not apply"
		return nil
	}); !errors.Is(err, ErrRevisionMismatch) {
		t.Fatalf("stale revision: expected ErrRevisionMismatch, got %v", err)
	}

	alias, err := s.GetAlias("fn", "live")
	if err != nil {
		t.Fatalf("reload alias: %v", err)
	}
	if alias.Description == "must not apply" {
		t.Fatalf("stale update leaked through the precondition")
	}

	updated, err := s.UpdateAliasAtomically("fn", "live", alias.RevisionId, func(fn *Function, existing *Alias) error {
		existing.Description = "applied"
		return nil
	})
	if err != nil {
		t.Fatalf("update with matching revision: %v", err)
	}
	if updated.Description != "applied" {
		t.Fatalf("matching update did not apply: %+v", updated)
	}
}
