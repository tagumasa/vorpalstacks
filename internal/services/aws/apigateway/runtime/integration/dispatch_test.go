package integration

import "testing"

// TestAWSServiceFromURI pins the structural dispatch of AWS-type
// integrations: the backend service comes from the resource segment of the
// parsed integration ARN, never from substring probes over the whole string —
// a URI merely containing a probe like "lambda:" must not dispatch.
func TestAWSServiceFromURI(t *testing.T) {
	tests := []struct {
		name string
		uri  string
		want string
	}{
		{
			name: "lambda path-style invocation URI",
			uri:  "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:aws:lambda:us-east-1:123456789012:function:my-fn/invocations",
			want: "lambda",
		},
		{
			name: "sqs path-style URI",
			uri:  "arn:aws:apigateway:us-east-1:sqs:path/123456789012/my-queue",
			want: "sqs",
		},
		{
			name: "sns path-style URI",
			uri:  "arn:aws:apigateway:eu-west-1:sns:path/my-topic",
			want: "sns",
		},
		{
			name: "dynamodb action-style URI",
			uri:  "arn:aws:apigateway:us-east-1:dynamodb:action/PutItem",
			want: "dynamodb",
		},
		{
			name: "kinesis action-style URI",
			uri:  "arn:aws:apigateway:us-east-1:kinesis:action/PutRecord",
			want: "kinesis",
		},
		{
			name: "step functions action-style URI",
			uri:  "arn:aws:apigateway:us-east-1:states:action/StartExecution",
			want: "states",
		},
		{
			name: "https URI carrying the old lambda probe text must not dispatch",
			uri:  "https://example.com/redirect?to=lambda:path/foo",
			want: "",
		},
		{
			name: "https URI carrying the old sqs probe text must not dispatch",
			uri:  "https://example.com/a:sqs:b",
			want: "",
		},
		{
			name: "non-apigateway ARN does not dispatch",
			uri:  "arn:aws:sqs:us-east-1:123456789012:my-queue",
			want: "",
		},
		{
			name: "empty URI",
			uri:  "",
			want: "",
		},
		{
			name: "URI without the path/action marker does not dispatch",
			uri:  "arn:aws:apigateway:us-east-1:lambda",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := awsServiceFromURI(tt.uri); got != tt.want {
				t.Errorf("awsServiceFromURI(%q) = %q, want %q", tt.uri, got, tt.want)
			}
		})
	}
}
