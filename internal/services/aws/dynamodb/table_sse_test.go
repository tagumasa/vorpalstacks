// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"testing"

	"vorpalstacks/internal/common/request"
)

// TestSSESpecificationFollowsTheModel pins the specification-side semantics
// the SSESpecification member documentation states: Enabled=true alone
// yields the KMS-with-managed-key state (the description echoes the
// documented default key's alias ARN — SSEDescription.KMSMasterKeyArn is
// "The KMS key ARN used for the KMS encryption"), an explicit
// KMSMasterKeyId settles to its ARN form (alias identifiers to alias ARNs,
// key identifiers to key ARNs, a full ARN verbatim, and the wired KMS
// resolver's answer when it is attached), AES256 inside a specification is
// refused (KMS is the only supported specification value), a specification
// without Enabled that carries SSEType or KMSMasterKeyId contradicts the
// owned-key default the absence states, and UpdateTable refuses an
// explicit disable — encryption at rest has no off state, only key types.
func TestSSESpecificationFollowsTheModel(t *testing.T) {
	svc, reqCtx := billingModePlaneFixture(t)
	ctx := context.Background()

	create := func(name string, sse interface{}) error {
		params := map[string]interface{}{
			"TableName":            name,
			"KeySchema":            []interface{}{map[string]interface{}{"AttributeName": "pk", "KeyType": "HASH"}},
			"AttributeDefinitions": []interface{}{map[string]interface{}{"AttributeName": "pk", "AttributeType": "S"}},
			"BillingMode":          "PAY_PER_REQUEST",
		}
		if sse != nil {
			params["SSESpecification"] = sse
		}
		_, err := svc.CreateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}
	describeSSE := func(name string) map[string]interface{} {
		t.Helper()
		resp, err := svc.DescribeTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{"TableName": name}})
		if err != nil {
			t.Fatalf("describe %s: %v", name, err)
		}
		sse, _ := resp.(map[string]interface{})["Table"].(map[string]interface{})["SSEDescription"].(map[string]interface{})
		return sse
	}

	if err := create("SseDefaultTable", map[string]interface{}{"Enabled": true}); err != nil {
		t.Fatalf("create with Enabled alone: %v", err)
	}
	sse := describeSSE("SseDefaultTable")
	if sse == nil {
		t.Fatal("enabled table carries no SSEDescription")
	}
	if sse["Status"] != "ENABLED" || sse["SSEType"] != "KMS" {
		t.Fatalf("default specification state: %+v", sse)
	}
	if sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:"+defaultSSEKeyAlias {
		t.Fatalf("default key alias ARN: %v", sse["KMSMasterKeyArn"])
	}

	if err := create("SseKeyTable", map[string]interface{}{"Enabled": true, "KMSMasterKeyId": "alias/my-key"}); err != nil {
		t.Fatalf("create with an explicit key: %v", err)
	}
	if sse := describeSSE("SseKeyTable"); sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:alias/my-key" {
		t.Fatalf("explicit alias key: %+v", sse)
	}

	if err := create("SseKeyIdTable", map[string]interface{}{"Enabled": true, "KMSMasterKeyId": "1234abcd-12ab-34cd-56ef-1234567890ab"}); err != nil {
		t.Fatalf("create with a key-id identifier: %v", err)
	}
	if sse := describeSSE("SseKeyIdTable"); sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab" {
		t.Fatalf("explicit key-id identifier: %+v", sse)
	}

	// The identifier's form is validated where it enters: an ARN that is
	// no key identifier — another service's ARN, or a KMS ARN whose
	// resource names no key or alias — is rejected before any resolver is
	// consulted, while a full KMS key ARN is stored verbatim.
	for _, bad := range []string{
		"arn:aws:sqs:us-east-1:123456789012:queue/leak",
		"arn:aws:kms:us-east-1:123456789012:policy/p1",
	} {
		if err := create("SseBadArnTable", map[string]interface{}{"Enabled": true, "KMSMasterKeyId": bad}); !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("KMSMasterKeyId %q: expected ErrInvalidParameter, got %v", bad, err)
		}
	}
	if err := create("SseKeyArnTable", map[string]interface{}{"Enabled": true, "KMSMasterKeyId": "arn:aws:kms:us-east-1:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab"}); err != nil {
		t.Fatalf("create with a full key ARN: %v", err)
	}
	if sse := describeSSE("SseKeyArnTable"); sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:key/1234abcd-12ab-34cd-56ef-1234567890ab" {
		t.Fatalf("full key ARN verbatim: %+v", sse)
	}

	// The wired resolver settles identifiers to the underlying key's ARN
	// and rejects an identifier the KMS store cannot resolve — the
	// platform extension the SQS and Kinesis key members carry.
	svc.SetKMSResolver(stubKMSResolver{arn: "arn:aws:kms:us-east-1:123456789012:key/resolved-1"})
	if err := create("SseResolvedTable", map[string]interface{}{"Enabled": true, "KMSMasterKeyId": "alias/my-key"}); err != nil {
		t.Fatalf("create with a resolvable key: %v", err)
	}
	if sse := describeSSE("SseResolvedTable"); sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:key/resolved-1" {
		t.Fatalf("resolved key ARN: %+v", sse)
	}
	// The managed default never consults the resolver: its reserved alias
	// is not one the platform's KMS hosts, so the deterministic alias ARN
	// applies even under a resolver that answers everything.
	if err := create("SseDefaultResolvedTable", map[string]interface{}{"Enabled": true}); err != nil {
		t.Fatalf("create with Enabled alone under a resolver: %v", err)
	}
	if sse := describeSSE("SseDefaultResolvedTable"); sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:"+defaultSSEKeyAlias {
		t.Fatalf("default key under a resolver: %+v", sse)
	}
	// Naming the default's reserved alias explicitly selects the same
	// effective key as the omission — identical outcome, no resolver
	// consultation, never the rejection an unresolvable identifier takes.
	if err := create("SseExplicitDefaultTable", map[string]interface{}{"Enabled": true, "KMSMasterKeyId": defaultSSEKeyAlias}); err != nil {
		t.Fatalf("create naming the default alias explicitly under a resolver: %v", err)
	}
	if sse := describeSSE("SseExplicitDefaultTable"); sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:"+defaultSSEKeyAlias {
		t.Fatalf("explicitly named default key: %+v", sse)
	}
	svc.SetKMSResolver(stubKMSResolver{err: errors.New("unresolvable")})
	if err := create("SseMissingKeyTable", map[string]interface{}{"Enabled": true, "KMSMasterKeyId": "alias/gone"}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("unresolvable key: expected ErrInvalidParameter, got %v", err)
	}
	// An erroring resolver still never sees the managed default.
	if err := create("SseDefaultMissingTable", map[string]interface{}{"Enabled": true}); err != nil {
		t.Fatalf("create with Enabled alone under an erroring resolver: %v", err)
	}
	svc.SetKMSResolver(nil)

	if err := create("SseAesTable", map[string]interface{}{"Enabled": true, "SSEType": "AES256"}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("AES256 inside a specification: expected ErrInvalidParameter, got %v", err)
	}

	if err := create("SseContradictionTable", map[string]interface{}{"SSEType": "KMS"}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("SSEType without Enabled: expected ErrInvalidParameter, got %v", err)
	}
	if err := create("SseContradictionTable", map[string]interface{}{"KMSMasterKeyId": "alias/x"}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("KMSMasterKeyId without Enabled: expected ErrInvalidParameter, got %v", err)
	}

	// A present member of the wrong wire type is a request error, never a
	// silently skipped specification: the structure member itself and the
	// Boolean Enabled both reject mistyped values.
	if err := create("SseMistypedSpecTable", "not-a-structure"); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("non-structure SSESpecification: expected ErrInvalidParameter, got %v", err)
	}
	if err := create("SseMistypedEnabledTable", map[string]interface{}{"Enabled": "yes"}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("non-bool Enabled: expected ErrInvalidParameter, got %v", err)
	}

	// An explicit disable at creation states the owned-key default a fresh
	// table already carries: accepted, no description.
	if err := create("SseDisabledCreateTable", map[string]interface{}{"Enabled": false}); err != nil {
		t.Fatalf("create with an explicit disable: %v", err)
	}
	if sse := describeSSE("SseDisabledCreateTable"); sse != nil {
		t.Fatalf("disabled-at-create table carries a description: %+v", sse)
	}

	// The update plane enables SSE through the same default form, and
	// refuses an explicit disable: encryption at rest has no off state.
	if err := create("SseUpdateTable", nil); err != nil {
		t.Fatalf("create plain table: %v", err)
	}
	resp, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":        "SseUpdateTable",
		"SSESpecification": map[string]interface{}{"Enabled": true},
	}})
	if err != nil {
		t.Fatalf("update SSE: %v", err)
	}
	sse = resp.(map[string]interface{})["TableDescription"].(map[string]interface{})["SSEDescription"].(map[string]interface{})
	if sse["SSEType"] != "KMS" || sse["KMSMasterKeyArn"] != "arn:aws:kms:us-east-1:123456789012:"+defaultSSEKeyAlias {
		t.Fatalf("update-plane description: %+v", sse)
	}
	if _, err := svc.UpdateTable(ctx, reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":        "SseUpdateTable",
		"SSESpecification": map[string]interface{}{"Enabled": false},
	}}); !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("explicit disable on update: expected ErrInvalidParameter, got %v", err)
	}
}

// stubKMSResolver pins the SSE key-identifier resolution contract without
// a KMS plane: it answers a fixed ARN or a fixed error.
type stubKMSResolver struct {
	arn string
	err error
}

func (r stubKMSResolver) ResolveKeyArn(ctx context.Context, region, keyID string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.arn, nil
}
