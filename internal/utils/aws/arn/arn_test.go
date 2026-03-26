package arn

import (
	"testing"
)

func TestParseARN(t *testing.T) {
	tests := []struct {
		name    string
		arn     string
		want    *ParsedARN
		wantErr bool
	}{
		{
			name:    "valid S3 bucket ARN",
			arn:     "arn:aws:s3:::my-bucket",
			want:    &ParsedARN{Partition: "aws", Service: "s3", Region: "", AccountID: "", Resource: "my-bucket"},
			wantErr: false,
		},
		{
			name:    "valid S3 bucket ARN with region",
			arn:     "arn:aws:s3:us-east-1:123456789012:my-bucket",
			want:    &ParsedARN{Partition: "aws", Service: "s3", Region: "us-east-1", AccountID: "123456789012", Resource: "my-bucket"},
			wantErr: false,
		},
		{
			name:    "valid IAM role ARN",
			arn:     "arn:aws:iam::123456789012:role/MyRole",
			want:    &ParsedARN{Partition: "aws", Service: "iam", Region: "", AccountID: "123456789012", Resource: "role/MyRole", RoleName: "MyRole"},
			wantErr: false,
		},
		{
			name:    "valid IAM role ARN with path",
			arn:     "arn:aws:iam::123456789012:role/path/to/MyRole",
			want:    &ParsedARN{Partition: "aws", Service: "iam", Region: "", AccountID: "123456789012", Resource: "role/path/to/MyRole", RoleName: "MyRole"},
			wantErr: false,
		},
		{
			name:    "valid Lambda function ARN",
			arn:     "arn:aws:lambda:us-east-1:123456789012:function:my-function",
			want:    &ParsedARN{Partition: "aws", Service: "lambda", Region: "us-east-1", AccountID: "123456789012", Resource: "function:my-function"},
			wantErr: false,
		},
		{
			name:    "valid DynamoDB table ARN",
			arn:     "arn:aws:dynamodb:us-east-1:123456789012:table/my-table",
			want:    &ParsedARN{Partition: "aws", Service: "dynamodb", Region: "us-east-1", AccountID: "123456789012", Resource: "table/my-table"},
			wantErr: false,
		},
		{
			name:    "China partition ARN",
			arn:     "arn:aws-cn:iam::123456789012:role/MyRole",
			want:    &ParsedARN{Partition: "aws-cn", Service: "iam", Region: "", AccountID: "123456789012", Resource: "role/MyRole", RoleName: "MyRole"},
			wantErr: false,
		},
		{
			name:    "invalid ARN - no arn prefix",
			arn:     "s3:::my-bucket",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "S3 ARN with empty resource",
			arn:     "arn:aws:s3:::",
			want:    &ParsedARN{Partition: "aws", Service: "s3", Region: "", AccountID: "", Resource: ""},
			wantErr: false,
		},
		{
			name:    "empty ARN",
			arn:     "",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseARN(tt.arn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseARN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Partition != tt.want.Partition {
				t.Errorf("ParseARN().Partition = %v, want %v", got.Partition, tt.want.Partition)
			}
			if got.Service != tt.want.Service {
				t.Errorf("ParseARN().Service = %v, want %v", got.Service, tt.want.Service)
			}
			if got.Region != tt.want.Region {
				t.Errorf("ParseARN().Region = %v, want %v", got.Region, tt.want.Region)
			}
			if got.AccountID != tt.want.AccountID {
				t.Errorf("ParseARN().AccountID = %v, want %v", got.AccountID, tt.want.AccountID)
			}
			if got.Resource != tt.want.Resource {
				t.Errorf("ParseARN().Resource = %v, want %v", got.Resource, tt.want.Resource)
			}
			if got.RoleName != tt.want.RoleName {
				t.Errorf("ParseARN().RoleName = %v, want %v", got.RoleName, tt.want.RoleName)
			}
		})
	}
}

func TestSplitARN(t *testing.T) {
	tests := []struct {
		name     string
		arn      string
		wantPart string
		wantServ string
		wantReg  string
		wantAcct string
		wantRes  string
	}{
		{
			name:     "valid S3 ARN",
			arn:      "arn:aws:s3:us-east-1:123456789012:my-bucket",
			wantPart: "aws",
			wantServ: "s3",
			wantReg:  "us-east-1",
			wantAcct: "123456789012",
			wantRes:  "my-bucket",
		},
		{
			name:     "valid Lambda ARN",
			arn:      "arn:aws:lambda:us-east-1:123456789012:function:my-function",
			wantPart: "aws",
			wantServ: "lambda",
			wantReg:  "us-east-1",
			wantAcct: "123456789012",
			wantRes:  "function:my-function",
		},
		{
			name:     "invalid ARN",
			arn:      "not-an-arn",
			wantPart: "",
			wantServ: "",
			wantReg:  "",
			wantAcct: "",
			wantRes:  "",
		},
		{
			name:     "empty string",
			arn:      "",
			wantPart: "",
			wantServ: "",
			wantReg:  "",
			wantAcct: "",
			wantRes:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			part, serv, reg, acct, res := SplitARN(tt.arn)
			if part != tt.wantPart {
				t.Errorf("SplitARN() partition = %v, want %v", part, tt.wantPart)
			}
			if serv != tt.wantServ {
				t.Errorf("SplitARN() service = %v, want %v", serv, tt.wantServ)
			}
			if reg != tt.wantReg {
				t.Errorf("SplitARN() region = %v, want %v", reg, tt.wantReg)
			}
			if acct != tt.wantAcct {
				t.Errorf("SplitARN() account = %v, want %v", acct, tt.wantAcct)
			}
			if res != tt.wantRes {
				t.Errorf("SplitARN() resource = %v, want %v", res, tt.wantRes)
			}
		})
	}
}

func TestGetServiceFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "S3 service",
			arn:  "arn:aws:s3:::my-bucket",
			want: "s3",
		},
		{
			name: "Lambda service",
			arn:  "arn:aws:lambda:us-east-1:123456789012:function:my-function",
			want: "lambda",
		},
		{
			name: "IAM service",
			arn:  "arn:aws:iam::123456789012:role/MyRole",
			want: "iam",
		},
		{
			name: "empty ARN",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetServiceFromARN(tt.arn); got != tt.want {
				t.Errorf("GetServiceFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidRoleARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want bool
	}{
		{
			name: "valid role ARN",
			arn:  "arn:aws:iam::123456789012:role/MyRole",
			want: true,
		},
		{
			name: "valid role ARN with path",
			arn:  "arn:aws:iam::123456789012:role/path/to/MyRole",
			want: true,
		},
		{
			name: "invalid - not IAM service",
			arn:  "arn:aws:s3:::my-bucket",
			want: false,
		},
		{
			name: "invalid - not a role",
			arn:  "arn:aws:iam::123456789012:user/myuser",
			want: false,
		},
		{
			name: "invalid - empty",
			arn:  "",
			want: false,
		},
		{
			name: "invalid - malformed",
			arn:  "not-an-arn",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidRoleARN(tt.arn); got != tt.want {
				t.Errorf("IsValidRoleARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParsedARN_IsCrossAccount(t *testing.T) {
	tests := []struct {
		name         string
		arn          *ParsedARN
		localAccount string
		want         bool
	}{
		{
			name:         "cross account",
			arn:          &ParsedARN{AccountID: "123456789012"},
			localAccount: "999999999999",
			want:         true,
		},
		{
			name:         "same account",
			arn:          &ParsedARN{AccountID: "123456789012"},
			localAccount: "123456789012",
			want:         false,
		},
		{
			name:         "empty account",
			arn:          &ParsedARN{AccountID: ""},
			localAccount: "123456789012",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arn.IsCrossAccount(tt.localAccount); got != tt.want {
				t.Errorf("IsCrossAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractResourceFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "S3 bucket",
			arn:  "arn:aws:s3:us-east-1:123456789012:my-bucket",
			want: "my-bucket",
		},
		{
			name: "Lambda function",
			arn:  "arn:aws:lambda:us-east-1:123456789012:function:my-function",
			want: "function:my-function",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractResourceFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractResourceFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractFunctionNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "Lambda function ARN",
			arn:  "arn:aws:lambda:us-east-1:123456789012:function:my-function",
			want: "my-function",
		},
		{
			name: "Lambda function ARN with qualifier",
			arn:  "arn:aws:lambda:us-east-1:123456789012:function:my-function:$LATEST",
			want: "my-function",
		},
		{
			name: "non-Lambda ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractFunctionNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractFunctionNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRoleNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "IAM role ARN",
			arn:  "arn:aws:iam::123456789012:role/MyRole",
			want: "MyRole",
		},
		{
			name: "IAM role ARN with path",
			arn:  "arn:aws:iam::123456789012:role/path/to/MyRole",
			want: "MyRole",
		},
		{
			name: "non-role ARN returns resource",
			arn:  "arn:aws:s3:::my-bucket",
			want: "my-bucket",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractRoleNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractRoleNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractQueueNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "SQS queue ARN",
			arn:  "arn:aws:sqs:us-east-1:123456789012:my-queue",
			want: "my-queue",
		},
		{
			name: "SQS queue ARN with FIFO",
			arn:  "arn:aws:sqs:us-east-1:123456789012:my-queue.fifo",
			want: "my-queue.fifo",
		},
		{
			name: "non-SQS ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "my-bucket",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractQueueNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractQueueNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractQueueNameFromURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "SQS queue URL",
			url:  "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue",
			want: "my-queue",
		},
		{
			name: "SQS FIFO queue URL",
			url:  "https://sqs.us-east-1.amazonaws.com/123456789012/my-queue.fifo",
			want: "my-queue.fifo",
		},
		{
			name: "empty",
			url:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractQueueNameFromURL(tt.url); got != tt.want {
				t.Errorf("ExtractQueueNameFromURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractTopicNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "SNS topic ARN",
			arn:  "arn:aws:sns:us-east-1:123456789012:my-topic",
			want: "my-topic",
		},
		{
			name: "SNS FIFO topic ARN",
			arn:  "arn:aws:sns:us-east-1:123456789012:my-topic.fifo",
			want: "my-topic.fifo",
		},
		{
			name: "non-SNS ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractTopicNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractTopicNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractLogGroupNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "Logs log group ARN",
			arn:  "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/my-function",
			want: "/aws/lambda/my-function",
		},
		{
			name: "Logs log group ARN with stream",
			arn:  "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/my-function:log-stream:my-stream",
			want: "/aws/lambda/my-function",
		},
		{
			name: "non-Logs ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractLogGroupNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractLogGroupNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractLogStreamNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "Logs log stream ARN",
			arn:  "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/my-function:log-stream:my-stream",
			want: "my-stream",
		},
		{
			name: "Logs log group without stream",
			arn:  "arn:aws:logs:us-east-1:123456789012:log-group:/aws/lambda/my-function",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractLogStreamNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractLogStreamNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractStreamNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "Kinesis stream ARN",
			arn:  "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream",
			want: "my-stream",
		},
		{
			name: "Kinesis stream ARN with partition",
			arn:  "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream/consumer/my-consumer",
			want: "my-stream",
		},
		{
			name: "non-Kinesis ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractStreamNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractStreamNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractEventBusNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "EventBridge event bus ARN",
			arn:  "arn:aws:events:us-east-1:123456789012:event-bus/my-bus",
			want: "my-bus",
		},
		{
			name: "EventBridge rule ARN",
			arn:  "arn:aws:events:us-east-1:123456789012:rule/my-bus/my-rule",
			want: "my-bus",
		},
		{
			name: "non-EventBridge ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractEventBusNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractEventBusNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRuleNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "EventBridge rule ARN with bus",
			arn:  "arn:aws:events:us-east-1:123456789012:rule/my-bus/my-rule",
			want: "my-rule",
		},
		{
			name: "EventBridge rule ARN standalone",
			arn:  "arn:aws:events:us-east-1:123456789012:rule/my-rule",
			want: "my-rule",
		},
		{
			name: "non-EventBridge ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractRuleNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractRuleNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractArchiveNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "EventBridge archive ARN",
			arn:  "arn:aws:events:us-east-1:123456789012:archive/my-archive",
			want: "my-archive",
		},
		{
			name: "non-archive ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractArchiveNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractArchiveNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractStateMachineNameFromARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want string
	}{
		{
			name: "StepFunctions state machine ARN",
			arn:  "arn:aws:states:us-east-1:123456789012:stateMachine:my-state-machine",
			want: "my-state-machine",
		},
		{
			name: "StepFunctions execution ARN",
			arn:  "arn:aws:states:us-east-1:123456789012:execution:my-state-machine:my-execution",
			want: "my-state-machine",
		},
		{
			name: "non-StepFunctions ARN",
			arn:  "arn:aws:s3:::my-bucket",
			want: "",
		},
		{
			name: "empty",
			arn:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractStateMachineNameFromARN(tt.arn); got != tt.want {
				t.Errorf("ExtractStateMachineNameFromARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsLambdaARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want bool
	}{
		{"Lambda ARN", "arn:aws:lambda:us-east-1:123456789012:function:my-function", true},
		{"S3 ARN", "arn:aws:s3:::my-bucket", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLambdaARN(tt.arn); got != tt.want {
				t.Errorf("IsLambdaARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsKinesisARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want bool
	}{
		{"Kinesis ARN", "arn:aws:kinesis:us-east-1:123456789012:stream/my-stream", true},
		{"S3 ARN", "arn:aws:s3:::my-bucket", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsKinesisARN(tt.arn); got != tt.want {
				t.Errorf("IsKinesisARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSQSARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want bool
	}{
		{"SQS ARN", "arn:aws:sqs:us-east-1:123456789012:my-queue", true},
		{"S3 ARN", "arn:aws:s3:::my-bucket", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSQSARN(tt.arn); got != tt.want {
				t.Errorf("IsSQSARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSNSARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want bool
	}{
		{"SNS ARN", "arn:aws:sns:us-east-1:123456789012:my-topic", true},
		{"S3 ARN", "arn:aws:s3:::my-bucket", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSNSARN(tt.arn); got != tt.want {
				t.Errorf("IsSNSARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEventBridgeARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want bool
	}{
		{"EventBridge ARN", "arn:aws:events:us-east-1:123456789012:event-bus/my-bus", true},
		{"S3 ARN", "arn:aws:s3:::my-bucket", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEventBridgeARN(tt.arn); got != tt.want {
				t.Errorf("IsEventBridgeARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsStateMachineARN(t *testing.T) {
	tests := []struct {
		name string
		arn  string
		want bool
	}{
		{"StateMachine ARN", "arn:aws:states:us-east-1:123456789012:stateMachine:my-sm", true},
		{"S3 ARN", "arn:aws:s3:::my-bucket", false},
		{"empty", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStateMachineARN(tt.arn); got != tt.want {
				t.Errorf("IsStateMachineARN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewARNBuilder(t *testing.T) {
	tests := []struct {
		name        string
		accountID   string
		region      string
		wantPart    string
		wantRegion  string
		wantAccount string
	}{
		{
			name:        "standard region",
			accountID:   "123456789012",
			region:      "us-east-1",
			wantPart:    PartitionAWS,
			wantRegion:  "us-east-1",
			wantAccount: "123456789012",
		},
		{
			name:        "China region",
			accountID:   "123456789012",
			region:      "cn-north-1",
			wantPart:    PartitionAWSCN,
			wantRegion:  "cn-north-1",
			wantAccount: "123456789012",
		},
		{
			name:        "empty region",
			accountID:   "123456789012",
			region:      "",
			wantPart:    PartitionAWS,
			wantRegion:  "",
			wantAccount: "123456789012",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewARNBuilder(tt.accountID, tt.region)
			if b.Partition() != tt.wantPart {
				t.Errorf("NewARNBuilder().Partition() = %v, want %v", b.Partition(), tt.wantPart)
			}
			if b.Region() != tt.wantRegion {
				t.Errorf("NewARNBuilder().Region() = %v, want %v", b.Region(), tt.wantRegion)
			}
			if b.AccountId() != tt.wantAccount {
				t.Errorf("NewARNBuilder().AccountId() = %v, want %v", b.AccountId(), tt.wantAccount)
			}
		})
	}
}

func TestARNBuilder_Build(t *testing.T) {
	b := NewARNBuilder("123456789012", "us-east-1")
	want := "arn:aws:s3:us-east-1:123456789012:my-bucket"
	got := b.Build("s3", "my-bucket")
	if got != want {
		t.Errorf("Build() = %v, want %v", got, want)
	}
}

func TestARNBuilder_BuildGlobal(t *testing.T) {
	b := NewARNBuilder("123456789012", "us-east-1")
	want := "arn:aws:iam:::role/MyRole"
	got := b.BuildGlobal("iam", "role/MyRole")
	if got != want {
		t.Errorf("BuildGlobal() = %v, want %v", got, want)
	}
}

func TestARNBuilder_BuildChinaPartition(t *testing.T) {
	b := NewARNBuilder("123456789012", "cn-north-1")
	want := "arn:aws-cn:s3:cn-north-1:123456789012:my-bucket"
	got := b.Build("s3", "my-bucket")
	if got != want {
		t.Errorf("Build() = %v, want %v", got, want)
	}
}
