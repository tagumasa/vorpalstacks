package eventbus

import (
	"context"
	"testing"
)

// TestLambdaResourcePolicyFnEffects pins the Lambda policy conversion: a
// Deny entry becomes a pattern-less Deny statement (its empty patterns
// match everything, so it overrides any Allow), and a well-formed Allow
// entry keeps matching exactly its principal/action/resource.
func TestLambdaResourcePolicyFnEffects(t *testing.T) {
	fn := LambdaResourcePolicyFn(func(ctx context.Context, functionARN string) ([]LambdaPolicyEntry, error) {
		return []LambdaPolicyEntry{
			{Effect: "Allow", Principal: "sqs.amazonaws.com", Action: "lambda:InvokeFunction", Resource: "arn:aws:lambda:us-east-1:123456789012:function:fn"},
			{Effect: "Deny"},
		}, nil
	})
	doc, err := fn(context.Background(), "arn:aws:lambda:us-east-1:123456789012:function:fn")
	if err != nil {
		t.Fatalf("lookup failed: %v", err)
	}
	if len(doc.Statement) != 2 {
		t.Fatalf("expected 2 converted statements, got %d", len(doc.Statement))
	}
	if doc.Statement[0].Effect != "Allow" || doc.Statement[1].Effect != "Deny" {
		t.Fatalf("unexpected statement effects: %+v", doc.Statement)
	}

	ev := NewSimplePolicyEvaluator()
	// The Deny entry denies a delivery the Allow entry would grant...
	ok, err := ev.Evaluate(context.Background(), doc, "sqs.amazonaws.com", "lambda:InvokeFunction", "arn:aws:lambda:us-east-1:123456789012:function:fn")
	if err != nil || ok {
		t.Fatalf("Deny entry must override the matching Allow: ok=%v err=%v", ok, err)
	}
	// ...and one it would not.
	ok, err = ev.Evaluate(context.Background(), doc, "events.amazonaws.com", "lambda:InvokeFunction", "arn:aws:lambda:us-east-1:123456789012:function:fn")
	if err != nil || ok {
		t.Fatalf("Deny entry must deny unmatched principals too: ok=%v err=%v", ok, err)
	}
}

// TestLambdaResourcePolicyFnSkipsUnevaluatableAllow verifies that an Allow
// entry whose action or principal is not a single string (array-shaped or
// absent) is dropped rather than converted into a wildcard grant.
func TestLambdaResourcePolicyFnSkipsUnevaluatableAllow(t *testing.T) {
	fn := LambdaResourcePolicyFn(func(ctx context.Context, functionARN string) ([]LambdaPolicyEntry, error) {
		return []LambdaPolicyEntry{
			{Effect: "Allow", Principal: "sns.amazonaws.com"},
			{Effect: "Allow", Action: "lambda:InvokeFunction"},
		}, nil
	})
	doc, err := fn(context.Background(), "arn:aws:lambda:us-east-1:123456789012:function:fn")
	if err != nil {
		t.Fatalf("lookup failed: %v", err)
	}
	if len(doc.Statement) != 0 {
		t.Fatalf("unevaluatable Allow entries must be skipped, got %+v", doc.Statement)
	}
}
