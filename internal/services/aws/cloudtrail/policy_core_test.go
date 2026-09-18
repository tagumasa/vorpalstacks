package cloudtrail

import (
	"errors"
	"strings"
	"testing"

	awserrors "vorpalstacks/internal/common/errors"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// The pins below drive the resource-policy dispatch with AWS-format ARNs:
// every resource the store creates carries its documented resource word
// (trail/, channel/, eventdatastore/). The documented ResourceArn set is the
// event data store, dashboard, and channel — trails are not policy targets,
// so a trail ARN is rejected with ResourceTypeNotSupportedException, and a
// hyphenated event data store word falls through to the same
// unsupported-type shape. The stored policy document must be a valid
// resource-based policy carrying Statement entries with principals.

func TestResourcePolicyRoundTripWithAWSFormatARNs(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	eds, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("eds-policy", "acc123", "us-east-1"))
	if err != nil {
		t.Fatalf("failed to create event data store: %v", err)
	}
	channel, err := store.CreateChannel(cloudtrailstore.NewChannel("channel-policy", "arn:aws:cloudtrail:acc123:channel-source", "acc123", "us-east-1"))
	if err != nil {
		t.Fatalf("failed to create channel: %v", err)
	}

	if !strings.Contains(eds.EventDataStoreARN, ":eventdatastore/") {
		t.Fatalf("created ARN %q lacks the AWS eventdatastore/ resource word", eds.EventDataStoreARN)
	}

	policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"cloudtrail:GetQueryResults","Resource":"*"}]}`
	for _, arn := range []string{eds.EventDataStoreARN, channel.ChannelARN} {
		if _, err := svc.putResourcePolicyCore(store, PutResourcePolicyInput{ResourceARN: arn, Policy: policy}); err != nil {
			t.Fatalf("PutResourcePolicy(%s) failed: %v", arn, err)
		}
		got, err := svc.getResourcePolicyCore(store, ResourcePolicyInput{ResourceARN: arn})
		if err != nil {
			t.Fatalf("GetResourcePolicy(%s) failed: %v", arn, err)
		}
		if got["ResourcePolicy"] != policy {
			t.Fatalf("round-tripped policy for %s = %v", arn, got["ResourcePolicy"])
		}
		if err := svc.deleteResourcePolicyCore(store, ResourcePolicyInput{ResourceARN: arn}); err != nil {
			t.Fatalf("DeleteResourcePolicy(%s) failed: %v", arn, err)
		}
	}
}

func TestResourcePolicyRejectsUnsupportedResourceTypes(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	trail, err := store.CreateTrail(cloudtrailstore.NewTrail("trail-policy", "policy-bucket", "us-east-1"))
	if err != nil {
		t.Fatalf("failed to create trail: %v", err)
	}

	// Trails are not in the documented ResourceArn set ("event data store,
	// dashboard, or channel"); the well-formed trail ARN is rejected as an
	// unsupported resource type, on every policy operation.
	for _, call := range []func() error{
		func() error {
			_, err := svc.putResourcePolicyCore(store, PutResourcePolicyInput{ResourceARN: trail.TrailARN, Policy: `{"Version":"2012-10-17","Statement":[]}`})
			return err
		},
		func() error {
			_, err := svc.getResourcePolicyCore(store, ResourcePolicyInput{ResourceARN: trail.TrailARN})
			return err
		},
		func() error {
			return svc.deleteResourcePolicyCore(store, ResourcePolicyInput{ResourceARN: trail.TrailARN})
		},
	} {
		err := call()
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.GetCode() != "ResourceTypeNotSupportedException" {
			t.Fatalf("trail ARN error = %v, want ResourceTypeNotSupportedException", err)
		}
	}
}

func TestResourcePolicyRejectsHyphenatedResourceWord(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	_, err := svc.getResourcePolicyCore(store, ResourcePolicyInput{
		ResourceARN: "arn:aws:cloudtrail:us-east-1:acc123:eventdata-store/EXAMPLE-f852-4e8f-8bd1-bcf6cEXAMPLE",
	})
	var awsErr *awserrors.AWSError
	if !errors.As(err, &awsErr) || awsErr.GetCode() != "ResourceTypeNotSupportedException" {
		t.Fatalf("hyphenated resource word error = %v, want ResourceTypeNotSupportedException", err)
	}
}

func TestPutResourcePolicyDocumentValidity(t *testing.T) {
	svc := &CloudTrailService{}
	store := newErrorTestServiceStore(t)

	eds, err := store.CreateEventDataStore(cloudtrailstore.NewEventDataStore("eds-policy-doc", "acc123", "us-east-1"))
	if err != nil {
		t.Fatalf("failed to create event data store: %v", err)
	}

	cases := []struct {
		name   string
		policy string
	}{
		{"syntax error", `{"Version":"2012-10-17"`},
		{"no statement", `{"Version":"2012-10-17"}`},
		{"statement without principal", `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":"cloudtrail:GetQueryResults","Resource":"*"}]}`},
	}
	for _, tc := range cases {
		_, err := svc.putResourcePolicyCore(store, PutResourcePolicyInput{ResourceARN: eds.EventDataStoreARN, Policy: tc.policy})
		var awsErr *awserrors.AWSError
		if !errors.As(err, &awsErr) || awsErr.GetCode() != "ResourcePolicyNotValidException" {
			t.Fatalf("%s: error = %v, want ResourcePolicyNotValidException", tc.name, err)
		}
	}

	// A single statement object and a NotPrincipal statement are both
	// valid documents.
	for _, policy := range []string{
		`{"Version":"2012-10-17","Statement":{"Effect":"Allow","Principal":"*","Action":"cloudtrail:GetQueryResults","Resource":"*"}}`,
		`{"Version":"2012-10-17","Statement":[{"Effect":"Deny","NotPrincipal":{"AWS":"*"},"Action":"cloudtrail:GetQueryResults","Resource":"*"}]}`,
	} {
		if _, err := svc.putResourcePolicyCore(store, PutResourcePolicyInput{ResourceARN: eds.EventDataStoreARN, Policy: policy}); err != nil {
			t.Fatalf("valid policy rejected: %v", err)
		}
	}
}
