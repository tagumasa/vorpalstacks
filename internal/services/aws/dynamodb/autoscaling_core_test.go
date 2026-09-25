package dynamodb

import (
	"context"
	"errors"
	"sync"
	"testing"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// TestParseAutoScalingSettingsTargetValueRange pins the documented
// TargetTrackingScalingPolicyConfiguration.TargetValue metric range on the
// shared auto-scaling settings parser: the member is required and its value
// must fall within the documented Base-10 bounds.
func TestParseAutoScalingSettingsTargetValueRange(t *testing.T) {
	settings := func(target float64) map[string]interface{} {
		return map[string]interface{}{
			"MinimumUnits": float64(5),
			"MaximumUnits": float64(50),
			"ScalingPolicyUpdate": map[string]interface{}{
				"PolicyName": "tracking-policy",
				"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
					"TargetValue": target,
				},
			},
		}
	}

	for _, valid := range []float64{50, autoScalingTargetValueMin, autoScalingTargetValueMax} {
		if _, err := parseAutoScalingSettings(settings(valid)); err != nil {
			t.Fatalf("TargetValue %v: unexpected error: %v", valid, err)
		}
	}
	for _, invalid := range []float64{0, -1, 1e-200, 2e108} {
		if _, err := parseAutoScalingSettings(settings(invalid)); err == nil {
			t.Fatalf("TargetValue %v: expected rejection", invalid)
		}
	}

	// The member is required: a target tracking configuration without a
	// TargetValue is rejected.
	noTarget := map[string]interface{}{
		"ScalingPolicyUpdate": map[string]interface{}{
			"PolicyName": "tracking-policy",
			"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
				"DisableScaleIn": true,
			},
		},
	}
	if _, err := parseAutoScalingSettings(noTarget); err == nil {
		t.Fatal("missing TargetValue: expected rejection")
	}
}

// TestParseAutoScalingSettingsTypedMembers pins the typed parse: present
// members land in the description, absent members stay nil, and capacity
// members typed as Long reject fractional or negative values.
func TestParseAutoScalingSettingsTypedMembers(t *testing.T) {
	roleArn := "arn:aws:iam::123456789012:role/auto-scaling"
	desc, err := parseAutoScalingSettings(map[string]interface{}{
		"MinimumUnits":        float64(5),
		"MaximumUnits":        float64(50),
		"AutoScalingDisabled": true,
		"AutoScalingRoleArn":  roleArn,
		"ScalingPolicyUpdate": map[string]interface{}{
			"PolicyName": "tracking-policy",
			"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{
				"TargetValue":      70.0,
				"DisableScaleIn":   true,
				"ScaleInCooldown":  float64(300),
				"ScaleOutCooldown": float64(60),
			},
		},
	})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if desc.MinimumUnits == nil || *desc.MinimumUnits != 5 {
		t.Fatalf("MinimumUnits = %v, want 5", desc.MinimumUnits)
	}
	if desc.MaximumUnits == nil || *desc.MaximumUnits != 50 {
		t.Fatalf("MaximumUnits = %v, want 50", desc.MaximumUnits)
	}
	if desc.AutoScalingDisabled == nil || !*desc.AutoScalingDisabled {
		t.Fatalf("AutoScalingDisabled = %v, want true", desc.AutoScalingDisabled)
	}
	if desc.AutoScalingRoleArn == nil || *desc.AutoScalingRoleArn != roleArn {
		t.Fatalf("AutoScalingRoleArn = %v", desc.AutoScalingRoleArn)
	}
	if len(desc.ScalingPolicies) != 1 {
		t.Fatalf("ScalingPolicies = %v, want one policy", desc.ScalingPolicies)
	}
	policy := desc.ScalingPolicies[0]
	if policy.PolicyName == nil || *policy.PolicyName != "tracking-policy" {
		t.Fatalf("PolicyName = %v", policy.PolicyName)
	}
	tt := policy.TargetTrackingScalingPolicyConfiguration
	if tt == nil || tt.TargetValue != 70 || tt.DisableScaleIn == nil || !*tt.DisableScaleIn {
		t.Fatalf("target tracking = %+v", tt)
	}
	if tt.ScaleInCooldown == nil || *tt.ScaleInCooldown != 300 || tt.ScaleOutCooldown == nil || *tt.ScaleOutCooldown != 60 {
		t.Fatalf("cooldowns = %v/%v", tt.ScaleInCooldown, tt.ScaleOutCooldown)
	}

	empty, err := parseAutoScalingSettings(map[string]interface{}{})
	if err != nil {
		t.Fatalf("empty parse: %v", err)
	}
	if empty.MinimumUnits != nil || empty.MaximumUnits != nil || empty.AutoScalingDisabled != nil ||
		empty.AutoScalingRoleArn != nil || empty.ScalingPolicies != nil {
		t.Fatalf("absent members must stay nil, got %+v", empty)
	}

	for _, invalid := range []interface{}{float64(1.5), float64(-1), "5"} {
		if _, err := parseAutoScalingSettings(map[string]interface{}{"MinimumUnits": invalid}); err == nil {
			t.Fatalf("MinimumUnits %v: expected rejection", invalid)
		}
	}

	// The cooldown members are modelled as IntegerObject (no range trait):
	// a non-integral number is rejected rather than truncated, and a value
	// outside the storage type's int32 bounds is rejected with it.
	settingsWith := func(ttOverride map[string]interface{}) map[string]interface{} {
		tt := map[string]interface{}{
			"TargetValue":      70.0,
			"ScaleInCooldown":  float64(300),
			"ScaleOutCooldown": float64(60),
		}
		for k, v := range ttOverride {
			tt[k] = v
		}
		return map[string]interface{}{
			"ScalingPolicyUpdate": map[string]interface{}{
				"PolicyName": "tracking-policy",
				"TargetTrackingScalingPolicyConfiguration": tt,
			},
		}
	}
	for _, invalid := range []interface{}{3.7, -0.5, "300", true, float64(3) * float64(1<<30)} {
		for _, member := range []string{"ScaleInCooldown", "ScaleOutCooldown"} {
			if _, err := parseAutoScalingSettings(settingsWith(map[string]interface{}{member: invalid})); err == nil {
				t.Fatalf("%s %v: expected rejection", member, invalid)
			}
		}
	}

	// A present value of the wrong type on the Boolean and String members
	// is the invalid-parameter error, never a silently skipped update.
	if _, err := parseAutoScalingSettings(map[string]interface{}{"AutoScalingDisabled": "yes"}); err == nil {
		t.Fatal("non-bool AutoScalingDisabled: expected rejection")
	}
	if _, err := parseAutoScalingSettings(map[string]interface{}{"AutoScalingRoleArn": 42}); err == nil {
		t.Fatal("non-string AutoScalingRoleArn: expected rejection")
	}
	if _, err := parseAutoScalingSettings(settingsWith(map[string]interface{}{"DisableScaleIn": "yes"})); err == nil {
		t.Fatal("non-bool DisableScaleIn: expected rejection")
	}

	// An explicit null member is absent — the protocol omits unset
	// members rather than carrying null — so no pointer is set and no
	// validation runs against a zero value the caller never sent.
	nulls, err := parseAutoScalingSettings(map[string]interface{}{
		"AutoScalingDisabled": nil,
		"AutoScalingRoleArn":  nil,
	})
	if err != nil {
		t.Fatalf("null members must parse as absent: %v", err)
	}
	if nulls.AutoScalingDisabled != nil || nulls.AutoScalingRoleArn != nil {
		t.Fatalf("null members must stay nil: %+v", nulls)
	}
	ttNulls, err := parseAutoScalingSettings(settingsWith(map[string]interface{}{
		"DisableScaleIn":   nil,
		"ScaleInCooldown":  nil,
		"ScaleOutCooldown": nil,
	}))
	if err != nil {
		t.Fatalf("null target-tracking members must parse as absent: %v", err)
	}
	policyTT := ttNulls.ScalingPolicies[0].TargetTrackingScalingPolicyConfiguration
	if policyTT.DisableScaleIn != nil || policyTT.ScaleInCooldown != nil || policyTT.ScaleOutCooldown != nil {
		t.Fatalf("null target-tracking members must stay nil: %+v", policyTT)
	}

	// The policy name rides the same typed-member family: a present
	// non-string value is the invalid-parameter error, never a silently
	// skipped policy name.
	if _, err := parseAutoScalingSettings(map[string]interface{}{
		"ScalingPolicyUpdate": map[string]interface{}{
			"PolicyName": 42,
			"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{"TargetValue": 70.0},
		},
	}); err == nil {
		t.Fatal("non-string PolicyName: expected rejection")
	}

	// The structure members take the same discipline: a present value of
	// another type is the invalid-parameter error, never a silently
	// skipped specification.
	if _, err := parseAutoScalingSettings(map[string]interface{}{"ScalingPolicyUpdate": "not-a-structure"}); err == nil {
		t.Fatal("non-structure ScalingPolicyUpdate: expected rejection")
	}
	if _, err := parseAutoScalingSettings(map[string]interface{}{
		"ScalingPolicyUpdate": map[string]interface{}{"TargetTrackingScalingPolicyConfiguration": 42},
	}); err == nil {
		t.Fatal("non-structure TargetTrackingScalingPolicyConfiguration: expected rejection")
	}
	if _, err := parseOptionalAutoScalingSettings(map[string]interface{}{"ProvisionedWriteCapacityAutoScalingUpdate": "not-a-structure"}, "ProvisionedWriteCapacityAutoScalingUpdate"); err == nil {
		t.Fatal("non-structure AutoScalingSettingsUpdate member: expected rejection")
	}
}

// TestUpdateTableReplicaAutoScalingRejectsWrongTypedMembers pins the
// typed-member presence contract on the replica face: the ReplicaUpdates
// carrier and the per-replica read-capacity structure member reject a
// present wrong-typed value — the invalid-parameter error — instead of
// silently treating it as an absent update.
func TestUpdateTableReplicaAutoScalingRejectsWrongTypedMembers(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "AsTypeFaceTable",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	update := func(replicaUpdates interface{}) error {
		_, err := svc.UpdateTableReplicaAutoScaling(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":      "AsTypeFaceTable",
			"ReplicaUpdates": replicaUpdates,
		}})
		return err
	}

	if err := update("not-a-list"); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-list ReplicaUpdates: expected the invalid-parameter rejection, got %v", err)
	}
	if err := update([]interface{}{map[string]interface{}{
		"RegionName": reqCtx.GetRegion(),
		"ReplicaProvisionedReadCapacityAutoScalingUpdate": "not-a-structure",
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Errorf("non-structure read-capacity member: expected the invalid-parameter rejection, got %v", err)
	}

	// Positive control: the well-typed update still lands.
	if err := update([]interface{}{map[string]interface{}{
		"RegionName": reqCtx.GetRegion(),
		"ReplicaProvisionedReadCapacityAutoScalingUpdate": map[string]interface{}{"MinimumUnits": 1.0, "MaximumUnits": 10.0},
	}}); err != nil {
		t.Fatalf("well-typed update: %v", err)
	}
}

// The per-index updates list holds the presence contract at every level:
// a present member of the wrong wire type — the list itself, an entry,
// IndexName, or the capacity-side member — is the invalid-parameter error,
// while absence (the list, IndexName, or the capacity member) still
// contributes nothing.
func TestParseGSIAutoScalingListMemberTypes(t *testing.T) {
	const capacity = "ProvisionedWriteCapacityAutoScalingUpdate"
	if _, err := parseGSIAutoScalingList("not-a-list", capacity); err == nil {
		t.Fatal("non-list updates member: expected rejection")
	}
	if _, err := parseGSIAutoScalingList([]interface{}{"not-a-map"}, capacity); err == nil {
		t.Fatal("non-map entry: expected rejection")
	}
	if _, err := parseGSIAutoScalingList([]interface{}{map[string]interface{}{"IndexName": 5}}, capacity); err == nil {
		t.Fatal("non-string IndexName: expected rejection")
	}
	if _, err := parseGSIAutoScalingList([]interface{}{
		map[string]interface{}{"IndexName": "gsi1", capacity: "not-a-structure"},
	}, capacity); err == nil {
		t.Fatal("non-structure capacity member: expected rejection")
	}

	absent, err := parseGSIAutoScalingList(nil, capacity)
	if err != nil || len(absent) != 0 {
		t.Fatalf("absent updates member must contribute nothing: %v, %v", absent, err)
	}
	skips, err := parseGSIAutoScalingList([]interface{}{
		nil,
		map[string]interface{}{},
		map[string]interface{}{"IndexName": "gsi1"},
	}, capacity)
	if err != nil || len(skips) != 0 {
		t.Fatalf("entries lacking IndexName or the capacity member contribute nothing: %v, %v", skips, err)
	}
}

func replicaDesc(region string) dbstore.ReplicaAutoScalingDescription {
	return dbstore.ReplicaAutoScalingDescription{RegionName: region}
}

func TestMergeReplicaDescriptionsUpsertsPerRegion(t *testing.T) {
	existing := []dbstore.ReplicaAutoScalingDescription{
		replicaDesc("us-east-1"),
		replicaDesc("eu-west-1"),
	}
	// An update for one region and a new region: us-east-1 is replaced,
	// eu-west-1 is preserved untouched, ap-southeast-1 is appended.
	units := int64(1)
	updatedEast := dbstore.ReplicaAutoScalingDescription{
		RegionName: "us-east-1",
		Read:       &dbstore.AutoScalingSettingsDescription{MinimumUnits: &units},
	}
	merged := mergeReplicaDescriptions(existing, []dbstore.ReplicaAutoScalingDescription{updatedEast, replicaDesc("ap-southeast-1")})

	if len(merged) != 3 {
		t.Fatalf("merged = %v, want 3 replicas", merged)
	}
	if merged[0].Read == nil {
		t.Fatalf("us-east-1 must carry the updated settings, got %v", merged[0])
	}
	if merged[1].RegionName != "eu-west-1" {
		t.Fatalf("unmentioned region must be preserved in place, got %v", merged[1])
	}
	if merged[2].RegionName != "ap-southeast-1" {
		t.Fatalf("new region must be appended, got %v", merged[2])
	}

	// An empty update leaves the stored list unchanged — an update without
	// ReplicaUpdates must not erase the stored replicas.
	same := mergeReplicaDescriptions(existing, nil)
	if len(same) != 2 || same[0].RegionName != "us-east-1" || same[1].RegionName != "eu-west-1" {
		t.Fatalf("empty update must preserve the stored replicas, got %v", same)
	}

	// The input list is never mutated.
	if existing[0].Read != nil {
		t.Fatal("merge must not mutate the stored descriptions")
	}
}

func TestReplicaDescriptionsFromSettingsNeverNil(t *testing.T) {
	if got := replicaDescriptionsFromSettings(nil); got == nil || len(got) != 0 {
		t.Fatalf("nil settings must yield an empty non-nil list, got %v", got)
	}
	got := replicaDescriptionsFromSettings(&dbstore.TableReplicaAutoScalingSettings{
		Replicas: []dbstore.ReplicaAutoScalingDescription{replicaDesc("us-east-1")},
	})
	if len(got) != 1 || got[0].RegionName != "us-east-1" {
		t.Fatalf("stored replicas must pass through, got %v", got)
	}
}

// TestAutoScalingSettingsWireRoundTrip pins the wire rendering: only set
// members appear, and the render covers the full member set.
func TestAutoScalingSettingsWireRoundTrip(t *testing.T) {
	if got := autoScalingSettingsToWire(nil); got != nil {
		t.Fatalf("nil description must render nil, got %v", got)
	}
	units := int64(5)
	wire := autoScalingSettingsToWire(&dbstore.AutoScalingSettingsDescription{
		MinimumUnits: &units,
		ScalingPolicies: []dbstore.AutoScalingPolicyDescription{{
			TargetTrackingScalingPolicyConfiguration: &dbstore.TargetTrackingScalingPolicyConfiguration{TargetValue: 70},
		}},
	})
	if wire["MinimumUnits"] != int64(5) {
		t.Fatalf("MinimumUnits = %v", wire["MinimumUnits"])
	}
	if _, ok := wire["MaximumUnits"]; ok {
		t.Fatal("absent members must not render")
	}
	policies, ok := wire["ScalingPolicies"].([]interface{})
	if !ok || len(policies) != 1 {
		t.Fatalf("ScalingPolicies = %v", wire["ScalingPolicies"])
	}
	policy := policies[0].(map[string]interface{})
	if _, ok := policy["PolicyName"]; ok {
		t.Fatal("absent PolicyName must not render")
	}
	tt, ok := policy["TargetTrackingScalingPolicyConfiguration"].(map[string]interface{})
	if !ok || tt["TargetValue"] != float64(70) {
		t.Fatalf("target tracking = %v", policy)
	}
}

// TestConcurrentReplicaAutoScalingUpdatesKeepEveryRegion pins the locked
// read-modify-write behind the settings write: concurrent
// UpdateTableReplicaAutoScaling calls, each updating a different region's
// replica description, must all land — an unlocked read-merge-write would
// have every caller start from the same stored state and the last write
// would silently discard the others' replicas.
func TestConcurrentReplicaAutoScalingUpdatesKeepEveryRegion(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":            "AsRaceTable",
		"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
		"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
		"BillingMode":          "PAY_PER_REQUEST",
	}}); err != nil {
		t.Fatalf("create table: %v", err)
	}

	regions := []string{"us-east-1", "eu-west-1", "us-west-2", "ap-northeast-1", "af-south-1", "eu-central-1", "sa-east-1", "me-south-1"}
	var wg sync.WaitGroup
	for _, region := range regions {
		wg.Add(1)
		go func(region string) {
			defer wg.Done()
			_, err := svc.UpdateTableReplicaAutoScaling(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
				"TableName": "AsRaceTable",
				"ReplicaUpdates": []interface{}{map[string]interface{}{
					"RegionName": region,
					"ReplicaProvisionedReadCapacityAutoScalingUpdate": map[string]interface{}{
						"MinimumUnits": 1.0,
						"MaximumUnits": 10.0,
					},
				}},
			}})
			if err != nil {
				t.Errorf("update %s: %v", region, err)
			}
		}(region)
	}
	wg.Wait()

	resp, err := svc.DescribeTableReplicaAutoScaling(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "AsRaceTable"}})
	if err != nil {
		t.Fatalf("describe: %v", err)
	}
	replicas := resp.(map[string]interface{})["TableAutoScalingDescription"].(map[string]interface{})["Replicas"].([]map[string]interface{})
	if len(replicas) != len(regions) {
		t.Fatalf("stored %d replica descriptions, want %d — a concurrent update was lost: %+v", len(replicas), len(regions), replicas)
	}
	seen := make(map[string]bool, len(replicas))
	for _, replica := range replicas {
		region := replica["RegionName"].(string)
		seen[region] = true
		read := replica["ReplicaProvisionedReadCapacityAutoScalingSettings"].(map[string]interface{})
		if read["MinimumUnits"].(int64) != 1 || read["MaximumUnits"].(int64) != 10 {
			t.Fatalf("region %s settings: %+v", region, read)
		}
	}
	for _, region := range regions {
		if !seen[region] {
			t.Fatalf("region %s missing from the stored settings: %+v", region, replicas)
		}
	}
}

// TestDeleteTableRemovesAutoScalingSettings pins the cascade sweep: the
// auto-scaling settings record is per-table state, so DeleteTable removes
// it with the table, and a same-name re-created table starts from an empty
// description instead of inheriting the dead table's replicas.
func TestDeleteTableRemovesAutoScalingSettings(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	create := func() {
		t.Helper()
		if _, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":            "AsLeakTable",
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		}}); err != nil {
			t.Fatalf("create table: %v", err)
		}
	}
	create()

	if _, err := svc.UpdateTableReplicaAutoScaling(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "AsLeakTable",
		"ReplicaUpdates": []interface{}{map[string]interface{}{
			"RegionName": reqCtx.GetRegion(),
			"ReplicaProvisionedReadCapacityAutoScalingUpdate": map[string]interface{}{
				"MinimumUnits": 5.0,
				"MaximumUnits": 50.0,
			},
		}},
	}}); err != nil {
		t.Fatalf("update auto-scaling: %v", err)
	}

	describe := func() []map[string]interface{} {
		t.Helper()
		resp, err := svc.DescribeTableReplicaAutoScaling(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "AsLeakTable"}})
		if err != nil {
			t.Fatalf("describe auto-scaling: %v", err)
		}
		return resp.(map[string]interface{})["TableAutoScalingDescription"].(map[string]interface{})["Replicas"].([]map[string]interface{})
	}

	// Positive control: the stored settings are visible before the delete,
	// so their absence afterwards is the sweep's work, not a seed gap.
	if replicas := describe(); len(replicas) != 1 || replicas[0]["RegionName"] != reqCtx.GetRegion() {
		t.Fatalf("pre-delete describe: %+v", replicas)
	}

	if _, err := svc.DeleteTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": "AsLeakTable"}}); err != nil {
		t.Fatalf("delete table: %v", err)
	}
	create()

	if replicas := describe(); len(replicas) != 0 {
		t.Fatalf("re-created table inherited the dead table's settings: %+v", replicas)
	}
}

// TestGSIAutoScalingWriteReadMemberExtraction pins both per-index list
// parsers: each side reads only its own capacity member by name, entries
// that are absent (null) or without an index name or without the member
// are skipped, and a malformed settings payload propagates the parse
// error. A present entry of the wrong wire type is the invalid-parameter
// error instead (TestParseGSIAutoScalingListMemberTypes).
func TestGSIAutoScalingWriteReadMemberExtraction(t *testing.T) {
	w, err := parseGSIAutoScalingWrite([]interface{}{
		map[string]interface{}{
			"IndexName": "gsi",
			"ProvisionedWriteCapacityAutoScalingUpdate": map[string]interface{}{"MinimumUnits": 2.0, "MaximumUnits": 20.0},
		},
		map[string]interface{}{"IndexName": "memberless"},
		map[string]interface{}{"ProvisionedWriteCapacityAutoScalingUpdate": map[string]interface{}{}},
		nil,
	})
	if err != nil {
		t.Fatalf("write list: %v", err)
	}
	if len(w) != 1 || w["gsi"] == nil || w["gsi"].MinimumUnits == nil || *w["gsi"].MinimumUnits != 2 {
		t.Fatalf("write-side extraction = %+v", w)
	}

	r, err := parseGSIAutoScalingRead([]interface{}{
		map[string]interface{}{
			"IndexName": "gsi",
			"ProvisionedReadCapacityAutoScalingUpdate": map[string]interface{}{"MinimumUnits": 1.0, "MaximumUnits": 10.0},
		},
	})
	if err != nil {
		t.Fatalf("read list: %v", err)
	}
	if len(r) != 1 || r["gsi"] == nil || r["gsi"].MaximumUnits == nil || *r["gsi"].MaximumUnits != 10 {
		t.Fatalf("read-side extraction = %+v", r)
	}

	// A malformed settings payload inside an entry fails the whole parse.
	if _, err := parseGSIAutoScalingWrite([]interface{}{map[string]interface{}{
		"IndexName": "gsi",
		"ProvisionedWriteCapacityAutoScalingUpdate": map[string]interface{}{
			"ScalingPolicyUpdate": map[string]interface{}{
				"TargetTrackingScalingPolicyConfiguration": map[string]interface{}{},
			},
		},
	}}); err == nil {
		t.Error("expected a settings payload without the required TargetValue to fail the parse")
	}
}
