package wafv2

import "testing"

// TestMatchesResourceType pins the ResourceType-to-service-namespace
// mapping used by ListResourcesForWebACL filtering. The namespaces follow
// the AWS ARN service namespaces reference; App Runner services live in
// the apprunner namespace and AgentCore gateways in bedrock-agentcore.
func TestMatchesResourceType(t *testing.T) {
	tests := []struct {
		name         string
		resourceArn  string
		resourceType string
		want         bool
	}{
		{"ALB", "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-alb/abc", "APPLICATION_LOAD_BALANCER", true},
		{"API Gateway", "arn:aws:apigateway:us-east-1::/restapis/api123/stages/prod", "API_GATEWAY", true},
		{"AppSync", "arn:aws:appsync:us-east-1:123456789012:apis/api123", "APPSYNC", true},
		{"Cognito user pool", "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_abc", "COGNITO_USER_POOL", true},
		{"App Runner service", "arn:aws:apprunner:us-east-1:123456789012:service/my-service/0123456789abcdef", "APP_RUNNER_SERVICE", true},
		{"Verified Access instance", "arn:aws:ec2:us-east-1:123456789012:verified-access-instance/vai-0123456789abcdef0", "VERIFIED_ACCESS_INSTANCE", true},
		{"Amplify app", "arn:aws:amplify:us-east-1:123456789012:apps/abc123", "AMPLIFY", true},
		{"AgentCore gateway", "arn:aws:bedrock-agentcore:us-west-2:123456789012:gateway/gateway-0123456789", "AGENTCORE_GATEWAY", true},

		{"wrong service", "arn:aws:sqs:us-east-1:123456789012:queue-name", "APPLICATION_LOAD_BALANCER", false},
		{"ALB type against SQS ARN", "arn:aws:sqs:us-east-1:123456789012:queue-name", "APP_RUNNER_SERVICE", false},
		{"plain ec2 ARN is not Verified Access", "arn:aws:ec2:us-east-1:123456789012:instance/i-0123456789abcdef0", "VERIFIED_ACCESS_INSTANCE", false},
		{"runner substring is not App Runner", "arn:aws:foo:us-east-1:123456789012:runner/x", "APP_RUNNER_SERVICE", false},
		{"agentcore substring is not AgentCore", "arn:aws:foo:us-east-1:123456789012:agentcore/x", "AGENTCORE_GATEWAY", false},
		{"not an ARN", "queue-name", "API_GATEWAY", false},
		{"unknown type", "arn:aws:apprunner:us-east-1:123456789012:service/my-service/0123", "S3_BUCKET", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesResourceType(tt.resourceArn, tt.resourceType); got != tt.want {
				t.Errorf("matchesResourceType(%q, %q) = %v, want %v", tt.resourceArn, tt.resourceType, got, tt.want)
			}
		})
	}
}
