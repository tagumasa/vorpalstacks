package dynamodb

import (
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
)

// newTablePersistStore opens a store holding one hash-only table.
func newTablePersistStore(t *testing.T) *DynamoDBStore {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatalf("open storage: %v", err)
	}
	t.Cleanup(func() { st.Close() })
	store := NewDynamoDBStore(st, "123456789012", "us-east-1")
	if _, err := store.Tables().Create(CreateTableParams{
		Name:                 "PersistTbl",
		KeySchema:            []*KeySchemaElement{{AttributeName: "id", KeyType: KeyTypeHash}},
		AttributeDefinitions: []*AttributeDefinition{{AttributeName: "id", AttributeType: ScalarAttributeTypeS}},
		BillingMode:          BillingModePayPerRequest,
	}); err != nil {
		t.Fatalf("create table: %v", err)
	}
	return store
}

// Resource-policy revisions must survive persistence: two SetResourcePolicy
// calls advance the stored revision to 1 then 2, re-read from disk.
func TestResourcePolicyRevisionPersists(t *testing.T) {
	tables := newTablePersistStore(t).Tables()

	if err := tables.SetResourcePolicy("PersistTbl", `{"Version":"2012-10-17"}`); err != nil {
		t.Fatalf("first put: %v", err)
	}
	rev, err := tables.GetResourcePolicyRevisionId("PersistTbl")
	if err != nil || rev != 1 {
		t.Fatalf("after first put: rev=%d err=%v, want 1", rev, err)
	}
	if err := tables.SetResourcePolicy("PersistTbl", `{"Version":"2012-10-17","x":1}`); err != nil {
		t.Fatalf("second put: %v", err)
	}
	rev, err = tables.GetResourcePolicyRevisionId("PersistTbl")
	if err != nil || rev != 2 {
		t.Fatalf("after second put: rev=%d err=%v, want 2", rev, err)
	}
}

// ContributorInsightsUpdatedAt is written by SetContributorInsights and must
// round-trip through the persisted table record.
func TestContributorInsightsUpdatedAtPersists(t *testing.T) {
	tables := newTablePersistStore(t).Tables()
	if err := tables.SetContributorInsights("PersistTbl", true, ""); err != nil {
		t.Fatalf("set contributor insights: %v", err)
	}
	table, err := tables.Get("PersistTbl")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if table.ContributorInsightsUpdatedAt.IsZero() {
		t.Fatal("ContributorInsightsUpdatedAt lost on persist")
	}
	if time.Since(table.ContributorInsightsUpdatedAt) > time.Minute {
		t.Fatalf("ContributorInsightsUpdatedAt not current: %v", table.ContributorInsightsUpdatedAt)
	}
}

// TableClass persists in the table record alone: a class set through Update
// survives, and clearing it sticks — no side bucket exists to resurrect it.
func TestTableClassRoundTripsInTableRecord(t *testing.T) {
	tables := newTablePersistStore(t).Tables()

	if _, err := tables.Update("PersistTbl", func(table *Table) error {
		table.TableClass = "STANDARD"
		return nil
	}); err != nil {
		t.Fatalf("set class: %v", err)
	}
	table, err := tables.Get("PersistTbl")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if table.TableClass != "STANDARD" {
		t.Fatalf("class lost: %q", table.TableClass)
	}

	if _, err := tables.Update("PersistTbl", func(table *Table) error {
		table.TableClass = ""
		return nil
	}); err != nil {
		t.Fatalf("clear class: %v", err)
	}
	table, err = tables.Get("PersistTbl")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if table.TableClass != "" {
		t.Fatalf("cleared class resurrected: %q", table.TableClass)
	}
}

// GlobalTableSourceArn set on a table persists in the table record.
func TestGlobalTableSourceArnPersists(t *testing.T) {
	tables := newTablePersistStore(t).Tables()
	const sourceArn = "arn:aws:dynamodb:us-east-1:123456789012:table/SourceTbl"
	if _, err := tables.Update("PersistTbl", func(table *Table) error {
		table.GlobalTableSourceArn = sourceArn
		return nil
	}); err != nil {
		t.Fatalf("set source arn: %v", err)
	}
	table, err := tables.Get("PersistTbl")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if table.GlobalTableSourceArn != sourceArn {
		t.Fatalf("source arn lost: %q", table.GlobalTableSourceArn)
	}
}

// Typed auto-scaling settings round-trip through the persisted record with
// pointer members (absence) preserved.
func TestAutoScalingSettingsTypedRoundTrip(t *testing.T) {
	tables := newTablePersistStore(t).Tables()

	if settings, err := tables.GetAutoScalingSettings("PersistTbl"); err != nil || settings != nil {
		t.Fatalf("absent settings: got %v, %v", settings, err)
	}

	minUnits, maxUnits := int64(5), int64(50)
	disabled := true
	roleArn := "arn:aws:iam::123456789012:role/auto-scaling"
	policyName := "tracking-policy"
	cooldown := int32(300)
	stored := &TableReplicaAutoScalingSettings{
		Replicas: []ReplicaAutoScalingDescription{{
			RegionName: "us-east-1",
			Read: &AutoScalingSettingsDescription{
				MinimumUnits:        &minUnits,
				MaximumUnits:        &maxUnits,
				AutoScalingDisabled: &disabled,
				AutoScalingRoleArn:  &roleArn,
				ScalingPolicies: []AutoScalingPolicyDescription{{
					PolicyName: &policyName,
					TargetTrackingScalingPolicyConfiguration: &TargetTrackingScalingPolicyConfiguration{
						DisableScaleIn:   &disabled,
						ScaleInCooldown:  &cooldown,
						ScaleOutCooldown: &cooldown,
						TargetValue:      70.5,
					},
				}},
			},
			GlobalSecondaryIndexes: []IndexAutoScalingSettings{{
				IndexName:                     "gsi-1",
				ProvisionedWriteCapacityUnits: &maxUnits,
				Write: &AutoScalingSettingsDescription{
					MinimumUnits: &minUnits,
				},
			}},
		}},
	}
	if err := tables.SetAutoScalingSettings("PersistTbl", stored); err != nil {
		t.Fatalf("set auto-scaling settings: %v", err)
	}

	loaded, err := tables.GetAutoScalingSettings("PersistTbl")
	if err != nil {
		t.Fatalf("get auto-scaling settings: %v", err)
	}
	if len(loaded.Replicas) != 1 {
		t.Fatalf("replicas = %v", loaded.Replicas)
	}
	replica := loaded.Replicas[0]
	if replica.RegionName != "us-east-1" || replica.Write != nil {
		t.Fatalf("replica = %+v", replica)
	}
	read := replica.Read
	if read == nil || read.MinimumUnits == nil || *read.MinimumUnits != 5 ||
		read.MaximumUnits == nil || *read.MaximumUnits != 50 ||
		read.AutoScalingDisabled == nil || !*read.AutoScalingDisabled ||
		read.AutoScalingRoleArn == nil || *read.AutoScalingRoleArn != roleArn {
		t.Fatalf("read settings = %+v", read)
	}
	if len(read.ScalingPolicies) != 1 || read.ScalingPolicies[0].PolicyName == nil ||
		*read.ScalingPolicies[0].PolicyName != policyName {
		t.Fatalf("scaling policies = %+v", read.ScalingPolicies)
	}
	tt := read.ScalingPolicies[0].TargetTrackingScalingPolicyConfiguration
	if tt == nil || tt.TargetValue != 70.5 || tt.ScaleInCooldown == nil || *tt.ScaleInCooldown != 300 {
		t.Fatalf("target tracking = %+v", tt)
	}
	if len(replica.GlobalSecondaryIndexes) != 1 ||
		replica.GlobalSecondaryIndexes[0].IndexName != "gsi-1" ||
		replica.GlobalSecondaryIndexes[0].ProvisionedWriteCapacityUnits == nil ||
		*replica.GlobalSecondaryIndexes[0].ProvisionedWriteCapacityUnits != 50 ||
		replica.GlobalSecondaryIndexes[0].Read != nil {
		t.Fatalf("index settings = %+v", replica.GlobalSecondaryIndexes)
	}
}
