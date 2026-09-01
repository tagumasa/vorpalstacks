package neptune

import (
	"errors"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	rdstypes "vorpalstacks/internal/store/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
)

// newTagTargetTestStore opens a throwaway Neptune store seeded with one
// cluster and one DB instance.
func newTagTargetTestStore(t *testing.T) neptunestore.NeptuneStoreInterface {
	t.Helper()
	raw, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { raw.Close() })
	store := neptunestore.NewNeptuneStore(raw)

	if err := store.CreateCluster(&rdstypes.DBCluster{DBClusterIdentifier: "pin-cluster", Engine: "neptune"}); err != nil {
		t.Fatalf("seed cluster: %v", err)
	}
	if err := store.CreateInstance(&rdstypes.DBInstance{DBInstanceIdentifier: "pin-instance", Engine: "neptune"}); err != nil {
		t.Fatalf("seed instance: %v", err)
	}
	return store
}

// TestValidateNeptuneTagTarget pins the tag-target resolution behind the
// Neptune tag trio: a seeded cluster and instance resolve, a missing
// resource of a hosted kind fails with the resource-specific NotFoundFault
// the model attaches to the tag operations, and an ARN of any other kind
// is invalid input.
func TestValidateNeptuneTagTarget(t *testing.T) {
	store := newTagTargetTestStore(t)
	arn := func(kind, name string) string {
		return "arn:aws:rds:us-east-1:000000000000:" + kind + ":" + name
	}

	t.Run("seeded resources resolve", func(t *testing.T) {
		for _, resource := range []string{
			arn("cluster", "pin-cluster"),
			arn("db", "pin-instance"),
		} {
			if err := validateNeptuneTagTarget(store, resource); err != nil {
				t.Errorf("validateNeptuneTagTarget(%q) = %v, want nil", resource, err)
			}
		}
	})

	t.Run("missing resource is the kind-specific NotFoundFault", func(t *testing.T) {
		tests := []struct{ arn, wantCode string }{
			{arn("cluster", "absent"), "DBClusterNotFoundFault"},
			{arn("db", "absent"), "DBInstanceNotFoundFault"},
			{arn("cluster-snapshot", "absent"), "DBClusterSnapshotNotFoundFault"},
			{arn("snapshot", "absent"), "DBSnapshotNotFoundFault"},
			{arn("subgrp", "absent"), "DBSubnetGroupNotFoundFault"},
		}
		for _, tt := range tests {
			err := validateNeptuneTagTarget(store, tt.arn)
			if err == nil {
				t.Fatalf("validateNeptuneTagTarget(%q) = nil, want %s", tt.arn, tt.wantCode)
			}
			var awsErr *awserrors.AWSError
			if !errors.As(err, &awsErr) {
				t.Fatalf("validateNeptuneTagTarget(%q) = %T, want *awserrors.AWSError", tt.arn, err)
			}
			if awsErr.Code != tt.wantCode {
				t.Errorf("%s: code = %q, want %q", tt.arn, awsErr.Code, tt.wantCode)
			}
		}
	})

	t.Run("unrecognised ARN kind is invalid input", func(t *testing.T) {
		for _, resource := range []string{
			arn("something-else", "whatever"),
			"arn:aws:rds:us-east-1:000000000000:cluster",
			"not-an-arn",
		} {
			err := validateNeptuneTagTarget(store, resource)
			if err == nil {
				t.Fatalf("validateNeptuneTagTarget(%q) = nil, want invalid-parameter", resource)
			}
			var awsErr *awserrors.AWSError
			if !errors.As(err, &awsErr) {
				t.Fatalf("validateNeptuneTagTarget(%q) = %T, want *awserrors.AWSError", resource, err)
			}
			if awsErr.Code != "InvalidParameterValue" {
				t.Errorf("%s: code = %q, want InvalidParameterValue", resource, awsErr.Code)
			}
		}
	})
}
