package auth

import "testing"

// TestExtractFunctionRefFromURI pins that the authorizer's function
// reference reaches the invoker as parsed: a qualified ARN reference keeps
// its function name and alias instead of collapsing to the last colon
// segment, and a URI outside the invocation grammar is rejected.
func TestExtractFunctionRefFromURI(t *testing.T) {
	cases := []struct {
		name string
		uri  string
		ref  string
		ok   bool
	}{
		{
			name: "bare function name",
			uri:  "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/my-fn/invocations",
			ref:  "my-fn",
			ok:   true,
		},
		{
			name: "qualified function ARN keeps name and alias",
			uri:  "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:my-fn:PROD/invocations",
			ref:  "arn:aws:lambda:us-east-1:123456789012:function:my-fn:PROD",
			ok:   true,
		},
		{
			name: "URI without the invocations suffix",
			uri:  "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/my-fn",
		},
		{
			name: "not a Lambda URI",
			uri:  "arn:aws:apigateway:us-east-1:sqs:path/123456789012/my-queue",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ref, err := extractFunctionRefFromURI(c.uri)
			if c.ok {
				if err != nil || ref != c.ref {
					t.Errorf("extractFunctionRefFromURI(%q) = (%q, %v), want (%q, nil)", c.uri, ref, err, c.ref)
				}
			} else if err == nil {
				t.Errorf("extractFunctionRefFromURI(%q) accepted the URI as %q", c.uri, ref)
			}
		})
	}
}
