package classifier

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"testing"
)

func newClassifier(t *testing.T) *Classifier {
	t.Helper()
	return &Classifier{index: nil}
}

func classify(t *testing.T, c *Classifier, method, path, contentType string, headers map[string]string, body string) *ClassifiedRequest {
	t.Helper()
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, path, bodyReader)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	cr, err := c.Classify(req)
	if err != nil {
		t.Fatalf("Classify returned error: %v", err)
	}
	return cr
}

func TestProtocolDetection(t *testing.T) {
	c := newClassifier(t)

	tests := []struct {
		name        string
		method      string
		path        string
		contentType string
		headers     map[string]string
		want        Protocol
	}{
		{"JSON-RPC via X-Amz-Target", "POST", "/", "application/json", map[string]string{"X-Amz-Target": "DynamoDB_20120810.CreateTable"}, ProtocolJSONRPC},
		{"Query via URL Action", "GET", "/?Action=GetUser", "", nil, ProtocolQuery},
		{"Query via form body", "POST", "/", "application/x-www-form-urlencoded", nil, ProtocolQuery},
		{"CBOR", "POST", "/", "application/cbor", nil, ProtocolCBOR},
		{"Neptune gremlin", "POST", "/gremlin", "application/json", nil, ProtocolNeptune},
		{"REST XML Route53", "POST", "/2013-04-01/hostedzone", "text/xml", nil, ProtocolRESTXML},
		{"REST JSON SESv2", "POST", "/v2/email/configuration-sets", "application/json", nil, ProtocolRESTJSON},
		{"REST JSON Lambda", "GET", "/2015-03-31/functions", "application/json", nil, ProtocolRESTJSON},
		{"Unknown", "GET", "/unknown", "text/plain", nil, ProtocolUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := classify(t, c, tt.method, tt.path, tt.contentType, tt.headers, "")
			if cr.Protocol != tt.want {
				t.Errorf("protocol = %v, want %v", cr.Protocol, tt.want)
			}
		})
	}
}

func TestServiceNameBySigningService(t *testing.T) {
	c := newClassifier(t)

	tests := []struct {
		name    string
		headers map[string]string
		body    string
		want    string
	}{
		{"timestream query via target", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/timestream/aws4_request", "X-Amz-Target": "Timestream_20181101.Query"}, "", "timestream-query"},
		{"timestream write via target", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/timestream/aws4_request", "X-Amz-Target": "Timestream_20181101.WriteRecords"}, "", "timestream-write"},
		{"timestream TagResource", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/timestream/aws4_request", "X-Amz-Target": "Timestream_20181101.TagResource"}, "", "timestream-query"},
		{"email", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/email/aws4_request"}, "", "email"},
		{"ses", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/ses/aws4_request"}, "", "email"},
		{"logs", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/logs/aws4_request"}, "", "logs"},
		{"states", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/states/aws4_request"}, "", "states"},
		{"runtime.sagemaker", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/runtime.sagemaker/aws4_request"}, "", "lambda"},
		{"wafv2", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/wafv2/aws4_request"}, "", "wafv2"},
		{"events", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/events/aws4_request"}, "", "eventbridge"},
		{"rds", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/rds/aws4_request"}, "", "neptune"},
		{"neptune-db", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/neptune-db/aws4_request"}, "", "neptunedata"},
		{"dynamodb identity", map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/dynamodb/aws4_request"}, "", "dynamodb"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := classify(t, c, "POST", "/", "application/json", tt.headers, tt.body)
			if cr.ServiceName != tt.want {
				t.Errorf("serviceName = %q, want %q", cr.ServiceName, tt.want)
			}
		})
	}
}

func TestServiceNameByTarget(t *testing.T) {
	c := newClassifier(t)

	tests := []struct {
		name   string
		target string
		want   string
	}{
		{"AmazonSNS.Publish", "AmazonSNS.Publish", "sns"},
		{"AmazonSQS.SendMessage", "AmazonSQS.SendMessage", "sqs"},
		{"AmazonSSM.GetParameter", "AmazonSSM.GetParameter", "ssm"},
		{"AWSCognitoIdentityProviderService.SignUp", "AWSCognitoIdentityProviderService.SignUp", "cognito-idp"},
		{"AWSCognitoIdentityService.GetId", "AWSCognitoIdentityService.GetId", "cognito-identity"},
		{"AWSEvents.PutRule", "AWSEvents.PutRule", "eventbridge"},
		{"Logs_20140328.CreateLogGroup", "Logs_20140328.CreateLogGroup", "logs"},
		{"DynamoDB_20120810.CreateTable", "DynamoDB_20120810.CreateTable", "dynamodb"},
		{"Kinesis_20131202.PutRecord", "Kinesis_20131202.PutRecord", "kinesis"},
		{"KMS_20141211.Encrypt", "KMS_20141211.Encrypt", "kms"},
		{"TrentService.Encrypt", "TrentService.Encrypt", "kms"},
		{"secretsmanager.GetSecretValue", "secretsmanager.GetSecretValue", "secretsmanager"},
		{"AWSStepFunctions.StartExecution", "AWSStepFunctions.StartExecution", "states"},
		{"Timestream_20181101.Query", "Timestream_20181101.Query", "timestream-query"},
		{"Timestream_20181101.WriteRecords", "Timestream_20181101.WriteRecords", "timestream-write"},
		{"GraniteServiceVersion20100801.GetMetricData", "GraniteServiceVersion20100801.GetMetricData", "monitoring"},
		{"AmazonAthena.StartQueryExecution", "AmazonAthena.StartQueryExecution", "athena"},
		{"CertificateManager.RequestCertificate", "CertificateManager.RequestCertificate", "acm"},
		{"AWSWAF_20190729.ListWebACLs", "AWSWAF_20190729.ListWebACLs", "wafv2"},
		{"AWSWAF_20150824.ListWebACLs", "AWSWAF_20150824.ListWebACLs", "waf"},
		{"CloudTrail_20131101.CreateTrail", "CloudTrail_20131101.CreateTrail", "cloudtrail"},
		{"com.amazonaws.cloudtrail.CreateTrail", "com.amazonaws.cloudtrail.CreateTrail", "cloudtrail"},
		{"AWSSTS.GetCallerIdentity", "AWSSTS.GetCallerIdentity", "sts"},
		{"amazonlambda.Invoke lowercase", "amazonlambda.Invoke", "lambda"},
		{"amazoncloudwatch.GetMetricStatistics lowercase", "amazoncloudwatch.GetMetricStatistics", "cloudwatch"},
		{"timestream_20181101.WriteRecords lowercase", "timestream_20181101.WriteRecords", "timestream-write"},
		{"timestreamquery_20181101.Query lowercase", "timestreamquery_20181101.Query", "timestream-query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := classify(t, c, "POST", "/", "application/x-amz-json-1.0", map[string]string{"X-Amz-Target": tt.target}, "")
			if cr.ServiceName != tt.want {
				t.Errorf("serviceName = %q, want %q", cr.ServiceName, tt.want)
			}
		})
	}
}

func TestServiceNameByPath(t *testing.T) {
	c := newClassifier(t)

	tests := []struct {
		name string
		path string
		want string
	}{
		{"lambda 2015-03-31", "/2015-03-31/functions", "lambda"},
		{"apigateway restapis", "/restapis/abc123", "apigateway"},
		{"apigateway apikeys", "/apikeys/abc123", "apigateway"},
		{"apigateway tags", "/tags/arn:aws:apigateway:us-east-1:restapi", "apigateway"},
		{"apigateway runtime excluded", "/restapis/abc123/stage/_user_request_/foo", ""},
		{"scheduler schedule-groups", "/schedule-groups/default", "scheduler"},
		{"scheduler schedules", "/schedules/my-schedule", "scheduler"},
		{"scheduler tags", "/tags/arn:aws:scheduler:region:account:schedule/name", "scheduler"},
		{"email", "/v2/email/configuration-sets", "email"},
		{"route53", "/2013-04-01/hostedzone/123", "route53"},
		{"cloudfront", "/2020-05-31/distribution", "cloudfront"},
		{"neptunedata gremlin", "/gremlin", "neptunedata"},
		{"neptunedata sparql", "/sparql/statistics", "neptunedata"},
		{"monitoring operation", "/service/GraniteServiceVersion20100801/operation/GetMetricData", "monitoring"},
		{"monitoring prefix", "/service/GraniteServiceVersion20100801/metrics", "monitoring"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := classify(t, c, "GET", tt.path, "", nil, "")
			if cr.ServiceName != tt.want {
				t.Errorf("serviceName = %q, want %q", cr.ServiceName, tt.want)
			}
		})
	}
}

func TestServiceNameByAction(t *testing.T) {
	c := newClassifier(t)

	tests := []struct {
		name   string
		action string
		want   string
	}{
		{"CreateUser", "CreateUser", "iam"},
		{"CreateTopic", "CreateTopic", "sns"},
		{"SendMessage", "SendMessage", "sqs"},
		{"GetCallerIdentity", "GetCallerIdentity", "sts"},
		{"PutEvents", "PutEvents", "events"},
		{"CreateStateMachine", "CreateStateMachine", "states"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/?Action="+tt.action, nil)
			if err != nil {
				t.Fatal(err)
			}
			cr, err := c.Classify(req)
			if err != nil {
				t.Fatal(err)
			}
			if cr.ServiceName != tt.want {
				t.Errorf("serviceName = %q, want %q", cr.ServiceName, tt.want)
			}
		})
	}
}

func TestAuthExtraction(t *testing.T) {
	c := newClassifier(t)

	cr := classify(t, c, "POST", "/", "application/json",
		map[string]string{"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/ap-northeast-1/dynamodb/aws4_request"}, "")

	if cr.Region != "ap-northeast-1" {
		t.Errorf("region = %q, want %q", cr.Region, "ap-northeast-1")
	}
	if cr.AccessKeyID != "AKIAIOSFODNN7EXAMPLE" {
		t.Errorf("accessKeyID = %q, want %q", cr.AccessKeyID, "AKIAIOSFODNN7EXAMPLE")
	}
}

func TestBodyParsing(t *testing.T) {
	c := newClassifier(t)

	t.Run("JSON body", func(t *testing.T) {
		cr := classify(t, c, "POST", "/", "application/json", nil, `{"TableName":"test","Key":{"id":{"S":"1"}}}`)
		if cr.Parameters["TableName"] != "test" {
			t.Errorf("Parameters[TableName] = %v, want %q", cr.Parameters["TableName"], "test")
		}
	})

	t.Run("form body", func(t *testing.T) {
		cr := classify(t, c, "POST", "/", "application/x-www-form-urlencoded", nil, "Action=GetUser&UserName=alice")
		if cr.Parameters["Action"] != "GetUser" {
			t.Errorf("Parameters[Action] = %v, want %q", cr.Parameters["Action"], "GetUser")
		}
	})

	t.Run("body restore", func(t *testing.T) {
		body := `{"test":true}`
		req, _ := http.NewRequest("POST", "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		cr, err := c.Classify(req)
		if err != nil {
			t.Fatal(err)
		}
		if cr == nil {
			t.Fatal("expected non-nil ClassifiedRequest")
		}
		readBack, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(readBack, []byte(body)) {
			t.Errorf("body not restored: got %q, want %q", string(readBack), body)
		}
	})
}

func TestOperationExtraction(t *testing.T) {
	c := newClassifier(t)

	tests := []struct {
		name    string
		method  string
		path    string
		headers map[string]string
		body    string
		want    string
	}{
		{"X-Amz-Target", "POST", "/", map[string]string{"X-Amz-Target": "DynamoDB_20120810.CreateTable"}, "", "CreateTable"},
		{"service/operation path", "POST", "/service/GraniteServiceVersion20100801/operation/GetMetricData", nil, "", "GetMetricData"},
		{"URL Action", "GET", "/?Action=CreateUser", nil, "", "CreateUser"},
		{"Lambda GetFunction", "GET", "/2015-03-31/functions/my-func", nil, "", "GetFunction"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := classify(t, c, tt.method, tt.path, "application/json", tt.headers, tt.body)
			if cr.Operation != tt.want {
				t.Errorf("operation = %q, want %q", cr.Operation, tt.want)
			}
		})
	}
}

func TestQueryParamsMerge(t *testing.T) {
	c := newClassifier(t)

	req, _ := http.NewRequest("GET", "/?MaxItems=10&Prefix=photos/", nil)
	cr, err := c.Classify(req)
	if err != nil {
		t.Fatal(err)
	}
	if cr.Parameters["MaxItems"] != "10" {
		t.Errorf("Parameters[MaxItems] = %v, want %q", cr.Parameters["MaxItems"], "10")
	}
	if cr.Parameters["Prefix"] != "photos/" {
		t.Errorf("Parameters[Prefix] = %v, want %q", cr.Parameters["Prefix"], "photos/")
	}
}

func TestLookupServiceByPath(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/2015-03-31/functions", "lambda"},
		{"/restapis/abc", "apigateway"},
		{"/restapis/abc/stage/_user_request_/foo", ""},
		{"/schedule-groups/x", "scheduler"},
		{"/v2/email/config", "email"},
		{"/2013-04-01/hostedzone", "route53"},
		{"/2020-05-31/distribution", "cloudfront"},
		{"/gremlin", "neptunedata"},
		{"/service/GraniteServiceVersion20100801/op", "monitoring"},
		{"/unknown/path", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := lookupServiceByPath(tt.path)
			if got != tt.want {
				t.Errorf("lookupServiceByPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestLookupServiceByTarget(t *testing.T) {
	tests := []struct {
		target string
		want   string
	}{
		{"AmazonSNS", "sns"},
		{"AmazonSQS", "sqs"},
		{"AWSEvents", "eventbridge"},
		{"DynamoDB_20120810", "dynamodb"},
		{"CertificateManager", "acm"},
		{"CloudTrail_20131101", "cloudtrail"},
		{"com.amazonaws.cloudtrail", "cloudtrail"},
		{"amazonlambda", "lambda"},
		{"amazoncloudwatch", "cloudwatch"},
		{"unknownservice", "unknownservice"},
	}

	for _, tt := range tests {
		t.Run(tt.target, func(t *testing.T) {
			got := lookupServiceByTarget(tt.target)
			if got != tt.want {
				t.Errorf("lookupServiceByTarget(%q) = %q, want %q", tt.target, got, tt.want)
			}
		})
	}
}

func TestIsApiGatewayPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/restapis/abc123", true},
		{"/apikeys/abc", true},
		{"/usageplans/1", true},
		{"/domainnames/example.com", true},
		{"/vpclinks/1", true},
		{"/apis/1", true},
		{"/authorizers/1", true},
		{"/tags/arn:aws:apigateway:region:restapi", true},
		{"/restapis/abc/stage/_user_request_/foo", false},
		{"/unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isApiGatewayPath(tt.path)
			if got != tt.want {
				t.Errorf("isApiGatewayPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNeptunedataPath(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/gremlin", true},
		{"/gremlin/status", true},
		{"/opencypher", true},
		{"/opencypher/status", true},
		{"/status", true},
		{"/system", true},
		{"/loader", true},
		{"/loader/job-1", true},
		{"/propertygraph", true},
		{"/propertygraph/statistics", true},
		{"/sparql/statistics", true},
		{"/rdf/statistics/summary", true},
		{"/ml/endpoints", true},
		{"/statusbar", false},
		{"/status-page", false},
		{"/systeminfo", false},
		{"/loaderx", false},
		{"/ml", false},
		{"/mls", false},
		{"/unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := isNeptunedataPath(tt.path)
			if got != tt.want {
				t.Errorf("isNeptunedataPath(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestFormBodyActionExtraction(t *testing.T) {
	c := newClassifier(t)

	cr := classify(t, c, "POST", "/", "application/x-www-form-urlencoded", nil, "Action=CreateUser&UserName=alice")
	if cr.ServiceName != "iam" {
		t.Errorf("serviceName = %q, want %q", cr.ServiceName, "iam")
	}
	if cr.Operation != "CreateUser" {
		t.Errorf("operation = %q, want %q", cr.Operation, "CreateUser")
	}
	if cr.Parameters["UserName"] != "alice" {
		t.Errorf("Parameters[UserName] = %v, want %q", cr.Parameters["UserName"], "alice")
	}
}

func TestSigningServiceOverridesPath(t *testing.T) {
	c := newClassifier(t)

	cr := classify(t, c, "POST", "/2015-03-31/functions", "application/json",
		map[string]string{
			"Authorization": "AWS4-HMAC-SHA256 Credential=AKIAIOSFODNN7EXAMPLE/20260101/us-east-1/dynamodb/aws4_request",
			"X-Amz-Target":  "DynamoDB_20120810.CreateTable",
		}, "")

	if cr.ServiceName != "dynamodb" {
		t.Errorf("serviceName = %q, want %q (SigningService should take priority over path)", cr.ServiceName, "dynamodb")
	}
}

func TestDefaultRegion(t *testing.T) {
	c := newClassifier(t)

	cr := classify(t, c, "GET", "/", "", nil, "")
	if cr.Region != "us-east-1" {
		t.Errorf("region = %q, want default %q", cr.Region, "us-east-1")
	}
}

func TestClassifyWithEncodedURL(t *testing.T) {
	c := newClassifier(t)

	req, _ := http.NewRequest("GET", "/v2/email/configuration-sets?ConfigurationSetName=MySet", nil)
	cr, err := c.Classify(req)
	if err != nil {
		t.Fatal(err)
	}
	if cr.ServiceName != "email" {
		t.Errorf("serviceName = %q, want %q", cr.ServiceName, "email")
	}
	if cr.QueryParams.Get("ConfigurationSetName") != "MySet" {
		t.Errorf("QueryParams[ConfigurationSetName] = %q, want %q", cr.QueryParams.Get("ConfigurationSetName"), "MySet")
	}
}
