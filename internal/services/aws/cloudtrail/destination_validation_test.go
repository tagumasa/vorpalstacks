package cloudtrail

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vorpalstacks/internal/common/invokers"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// Pins for the CreateTrail/UpdateTrail destination validation: the S3
// bucket's existence and policy sufficiency, the SNS topic's resolution
// (name or ARN) and policy sufficiency, and the CloudWatch Logs stream
// naming the delivery uses.

// newTrailTestService builds a service whose registry answers every S3
// bucket as an existing, sufficiently-privileged destination — for Core
// tests whose subject is not the destination itself.
func newTrailTestService(store *cloudtrailstore.CloudTrailStore) *CloudTrailService {
	svc := NewCloudTrailService(store.GetAccountID(), store.GetRegion())
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: newPermissiveS3Invoker()})
	return svc
}

// fakeSNSInvoker is an in-memory SNSInvoker: topics map ARN to policy
// document, and publishes are recorded for assertions.
type fakeSNSInvoker struct {
	topics    map[string]string
	published []fakeSNSPublish
}

type fakeSNSPublish struct {
	topicARN string
	message  string
	subject  string
}

func (f *fakeSNSInvoker) GetTopic(_ context.Context, topicARN string) (string, error) {
	if _, ok := f.topics[topicARN]; ok {
		return topicARN, nil
	}
	return "", fmt.Errorf("topic %s not found", topicARN)
}

func (f *fakeSNSInvoker) GetTopicPolicy(_ context.Context, topicARN string) (string, error) {
	if _, ok := f.topics[topicARN]; ok {
		return f.topics[topicARN], nil
	}
	return "", fmt.Errorf("topic %s not found", topicARN)
}

func (f *fakeSNSInvoker) ListSubscriptionsByTopic(_ context.Context, _ string) ([]invokers.SubscriptionInfo, error) {
	return nil, nil
}

func (f *fakeSNSInvoker) PublishToTopic(_ context.Context, topicARN, message, subject string, _ map[string]invokers.SQSMessageAttribute) (string, error) {
	if _, ok := f.topics[topicARN]; !ok {
		return "", fmt.Errorf("topic %s not found", topicARN)
	}
	f.published = append(f.published, fakeSNSPublish{topicARN: topicARN, message: message, subject: subject})
	return "msg-1", nil
}

// sufficientTopicPolicy is the AWS-documented SNS topic policy for
// CloudTrail notifications.
func sufficientTopicPolicy(topicARN string) string {
	return fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Sid":"AWSCloudTrailSNSPolicy20131101","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"%s"}]}`, topicARN)
}

// fakeLogsInvoker is an in-memory LogsInvoker: streams exist under groups
// the caller created, and PutLogEvents records the entries per stream. A
// stream under a missing group fails, mirroring the real service;
// failPuts injects that many PutLogEvents failures before recording
// resumes.
type fakeLogsInvoker struct {
	groups     map[string]bool
	streams    map[string]bool
	events     map[string][]invokers.LogsLogEntry
	batchSizes map[string][]int
	failPuts   int
}

func newFakeLogsInvoker(groups ...string) *fakeLogsInvoker {
	f := &fakeLogsInvoker{
		groups:     make(map[string]bool),
		streams:    make(map[string]bool),
		events:     make(map[string][]invokers.LogsLogEntry),
		batchSizes: make(map[string][]int),
	}
	for _, g := range groups {
		f.groups[g] = true
	}
	return f
}

func (f *fakeLogsInvoker) EnsureLogGroup(_ context.Context, _ string, group, _ string) error {
	f.groups[group] = true
	return nil
}

func (f *fakeLogsInvoker) EnsureLogStream(_ context.Context, _ string, group, stream string) error {
	if !f.groups[group] {
		return fmt.Errorf("log group %s does not exist", group)
	}
	f.streams[group+"/"+stream] = true
	return nil
}

func (f *fakeLogsInvoker) PutLogEvents(_ context.Context, _ string, group, stream string, entries []invokers.LogsLogEntry) error {
	if !f.streams[group+"/"+stream] {
		return fmt.Errorf("log stream %s does not exist", stream)
	}
	if f.failPuts > 0 {
		f.failPuts--
		return fmt.Errorf("injected stream failure")
	}
	f.events[group+"/"+stream] = append(f.events[group+"/"+stream], entries...)
	f.batchSizes[group+"/"+stream] = append(f.batchSizes[group+"/"+stream], len(entries))
	return nil
}

// documentedBucketPolicy renders the AWS-documented bucket policy for a
// trail's delivery prefix.
func documentedBucketPolicy(bucket, prefix, accountID, region, trailName string) string {
	trailARN := fmt.Sprintf("arn:aws:cloudtrail:%s:%s:trail/%s", region, accountID, trailName)
	return fmt.Sprintf(`{"Version":"2012-10-17","Statement":[`+
		`{"Sid":"AWSCloudTrailAclCheck20150319","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::%s"},`+
		`{"Sid":"AWSCloudTrailWrite20150319","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::%s/%sAWSLogs/%s/*","Condition":{"StringEquals":{"s3:x-amz-acl":"bucket-owner-full-control","aws:SourceArn":"%s"}}}]}`,
		bucket, bucket, prefix, accountID, trailARN)
}

func TestS3BucketPolicySufficiency(t *testing.T) {
	const (
		bucket  = "amzn-s3-demo-bucket"
		account = "123456789012"
		region  = "us-east-2"
		trail   = "arn:aws:cloudtrail:us-east-2:123456789012:trail/Trail1"
	)
	cases := []struct {
		name   string
		policy string
		prefix string
		want   bool
	}{
		{"documented policy with prefix and conditions", documentedBucketPolicy(bucket, "pref/", account, region, "Trail1"), "pref/", true},
		{"documented policy without prefix", documentedBucketPolicy(bucket, "", account, region, "Trail1"), "", true},
		{"no policy", "", "", false},
		{"not json", "not json", "", false},
		{"put object only", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::amzn-s3-demo-bucket/*"}]}`, "", false},
		{"get acl only", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::amzn-s3-demo-bucket"}]}`, "", false},
		{"wrong principal", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]}]}`, "", false},
		{"wildcard principal", `{"Statement":[{"Effect":"Allow","Principal":"*","Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]}]}`, "", true},
		{"service wildcard actions", `{"Statement":[{"Effect":"Allow","Principal":{"Service":["cloudtrail.amazonaws.com"]},"Action":"s3:*","Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]}]}`, "", true},
		{"partial wildcard actions", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucket*","s3:Put*"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]}]}`, "", true},
		{"mixed-case partial wildcard action", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["S3:GETBUCKETACL","s3:Put*"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]}]}`, "", true},
		{"deny effect", `{"Statement":[{"Effect":"Deny","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]}]}`, "", false},
		{"bucket arn does not cover objects", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":"arn:aws:s3:::amzn-s3-demo-bucket"}]}`, "", false},
		{"source arn condition matching this trail", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"StringEquals":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/Trail1"}}}]}`, "", true},
		{"source arn condition naming another trail", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"StringEquals":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/OtherTrail"}}}]}`, "", false},
		{"source arn array covering this trail", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"StringEquals":{"aws:SourceArn":["arn:aws:cloudtrail:us-east-1:123456789012:trail/A","arn:aws:cloudtrail:us-east-2:123456789012:trail/Trail1"]}}}]}`, "", true},
		{"source account condition matching", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"StringEquals":{"aws:SourceAccount":"123456789012"}}}]}`, "", true},
		{"StringLike source account pattern", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"StringLike":{"aws:SourceAccount":"123456789*"}}}]}`, "", true},
		{"StringEquals source account mismatch fails", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"StringEquals":{"aws:SourceAccount":"999999999999"}}}]}`, "", false},
		{"secure transport demanded on the allow", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"Bool":{"aws:SecureTransport":"true"}}}]}`, "", true},
		{"secure transport demanded false fails closed", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"Bool":{"aws:SecureTransport":"false"}}}]}`, "", false},
		{"Bool on a key outside the delivery context fails closed", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"Bool":{"aws:MultiFactorAuthPresent":"true"}}}]}`, "", false},
		{"unevaluable condition key fails closed", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"IpAddress":{"aws:SourceIp":"10.0.0.0/8"}}}]}`, "", false},
		{"prefix policy does not cover a deeper prefix", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/AWSLogs/*"]}]}`, "pref", false},
		{"ArnEquals condition is a pattern match", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"ArnEquals":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/*"}}}]}`, "", true},
		{"ArnLike condition with single-char wildcard", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"],"Condition":{"ArnLike":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/Trail?"}}}]}`, "", true},
		{"single-char wildcard in a resource pattern", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/?WSLogs/*"]}]}`, "", true},
		{"deny statement overrides the allow", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]},{"Effect":"Deny","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::amzn-s3-demo-bucket/*"}]}`, "", false},
		{"deny on an unrelated action does not veto", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]},{"Effect":"Deny","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:DeleteObject","Resource":"arn:aws:s3:::amzn-s3-demo-bucket/*"}]}`, "", true},
		{"deny with an unsatisfied condition does not veto", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]},{"Effect":"Deny","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::amzn-s3-demo-bucket/*","Condition":{"StringEquals":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/Other"}}}]}`, "", true},
		{"deny with a numeric source-account value vetoes", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]},{"Effect":"Deny","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::amzn-s3-demo-bucket/*","Condition":{"StringEquals":{"aws:SourceAccount":123456789012}}}]}`, "", false},
		{"deny with a boolean condition parses and stays unmatched", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":["s3:GetBucketAcl","s3:PutObject"],"Resource":["arn:aws:s3:::amzn-s3-demo-bucket","arn:aws:s3:::amzn-s3-demo-bucket/*"]},{"Effect":"Deny","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:PutObject","Resource":"arn:aws:s3:::amzn-s3-demo-bucket/*","Condition":{"Bool":{"aws:SecureTransport":false}}}]}`, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, s3BucketPolicySufficient(tc.policy, bucket, tc.prefix, account, region, trail))
		})
	}
}

// TestJsonStringOrListScalars pins the policy-value decoder: strings and
// string lists pass through, unquoted Boolean and numeric scalars decode
// as their JSON literal (the grammar permits dropping the quotation
// marks), and a value form the grammar does not cover still errors so its
// statement is reported unparseable rather than silently reshaped.
func TestJsonStringOrListScalars(t *testing.T) {
	cases := []struct {
		raw  string
		want []string
	}{
		{`"abc"`, []string{"abc"}},
		{`["a","b"]`, []string{"a", "b"}},
		{`false`, []string{"false"}},
		{`true`, []string{"true"}},
		{`443`, []string{"443"}},
		{`1.5`, []string{"1.5"}},
		{`["a", 443, false]`, []string{"a", "443", "false"}},
	}
	for _, tc := range cases {
		got, err := jsonStringOrListDecode([]byte(tc.raw))
		require.NoError(t, err, tc.raw)
		assert.Equal(t, tc.want, got, tc.raw)
	}
	for _, raw := range []string{`{"k":1}`, `[{"k":1}]`} {
		_, err := jsonStringOrListDecode([]byte(raw))
		assert.Error(t, err, raw)
	}
}

func TestSnsTopicPolicySufficiency(t *testing.T) {
	const (
		topic   = "arn:aws:sns:us-east-2:123456789012:myTopic"
		trail   = "arn:aws:cloudtrail:us-east-2:123456789012:trail/Trail1"
		account = "123456789012"
	)
	cases := []struct {
		name   string
		policy string
		want   bool
	}{
		{"documented policy", sufficientTopicPolicy(topic), true},
		{"no policy", "", false},
		{"wrong action", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Subscribe","Resource":"arn:aws:sns:us-east-2:123456789012:myTopic"}]}`, false},
		{"wrong resource", sufficientTopicPolicy("arn:aws:sns:us-east-2:123456789012:other"), false},
		{"wildcard resource", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"sns:Publish","Resource":"*"}]}`, true},
		{"uppercase action from the documented example", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"arn:aws:sns:us-east-2:123456789012:myTopic"}]}`, true},
		{"source arn condition naming this trail", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"arn:aws:sns:us-east-2:123456789012:myTopic","Condition":{"StringEquals":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/Trail1"}}}]}`, true},
		{"source arn condition naming another trail", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"arn:aws:sns:us-east-2:123456789012:myTopic","Condition":{"StringEquals":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/Other"}}}]}`, false},
		{"deny statement overrides the allow", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"sns:Publish","Resource":"arn:aws:sns:us-east-2:123456789012:myTopic"},{"Effect":"Deny","Principal":"*","Action":"sns:Publish","Resource":"arn:aws:sns:us-east-2:123456789012:myTopic"}]}`, false},
		{"ArnLike condition is a pattern match", `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"SNS:Publish","Resource":"arn:aws:sns:us-east-2:123456789012:myTopic","Condition":{"ArnLike":{"aws:SourceArn":"arn:aws:cloudtrail:us-east-2:123456789012:trail/*"}}}]}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, snsTopicPolicySufficient(tc.policy, topic, trail, account))
		})
	}
}

func TestCreateTrailDestinationValidation(t *testing.T) {
	store := newErrorTestServiceStore(t)
	ctx := context.Background()

	s3 := newFakeS3Invoker("dest-bucket")
	svc := NewCloudTrailService(store.GetAccountID(), store.GetRegion())
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3})

	t.Run("missing bucket is rejected with S3BucketDoesNotExistException", func(t *testing.T) {
		_, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "no-bucket-trail", S3BucketName: "ghost-bucket", Region: store.GetRegion()})
		requireAWSCode(t, err, "S3BucketDoesNotExistException", 404)
	})

	t.Run("insufficient policy is rejected with InsufficientS3BucketPolicyException", func(t *testing.T) {
		s3.policies["dest-bucket"] = `{"Statement":[{"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Action":"s3:GetBucketAcl","Resource":"arn:aws:s3:::dest-bucket"}]}`
		_, err := svc.createTrailCore(ctx, store, CreateTrailInput{Name: "bad-policy-trail", S3BucketName: "dest-bucket", Region: store.GetRegion()})
		requireAWSCode(t, err, "InsufficientS3BucketPolicyException", 403)
		s3.policies["dest-bucket"] = sufficientAnyBucketPolicy("dest-bucket")
	})

	t.Run("missing sns topic is rejected with InsufficientSnsTopicPolicyException", func(t *testing.T) {
		_, err := svc.createTrailCore(ctx, store, CreateTrailInput{
			Name: "no-topic-trail", S3BucketName: "dest-bucket", SnsTopicName: "ghost-topic", Region: store.GetRegion(),
		})
		requireAWSCode(t, err, "InsufficientSnsTopicPolicyException", 403)
	})

	t.Run("topic without policy is rejected", func(t *testing.T) {
		sns := &fakeSNSInvoker{topics: map[string]string{
			"arn:aws:sns:us-east-1:acc123:topic-no-policy": "",
		}}
		svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3, sns: sns})
		_, err := svc.createTrailCore(ctx, store, CreateTrailInput{
			Name: "no-policy-topic-trail", S3BucketName: "dest-bucket", SnsTopicName: "topic-no-policy", Region: store.GetRegion(),
		})
		requireAWSCode(t, err, "InsufficientSnsTopicPolicyException", 403)
	})

	t.Run("resolvable topic by name stores the ARN", func(t *testing.T) {
		sns := &fakeSNSInvoker{topics: map[string]string{}}
		topicARN := "arn:aws:sns:us-east-1:acc123:notify"
		sns.topics[topicARN] = sufficientTopicPolicy(topicARN)
		svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3, sns: sns})
		created, err := svc.createTrailCore(ctx, store, CreateTrailInput{
			Name: "topic-trail", S3BucketName: "dest-bucket", SnsTopicName: "notify", Region: store.GetRegion(),
		})
		require.NoError(t, err)
		assert.Equal(t, topicARN, created.SnsTopicARN)
	})

	t.Run("topic by ARN is accepted verbatim", func(t *testing.T) {
		sns := &fakeSNSInvoker{topics: map[string]string{}}
		topicARN := "arn:aws:sns:us-east-1:acc123:by-arn"
		sns.topics[topicARN] = sufficientTopicPolicy(topicARN)
		svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3, sns: sns})
		created, err := svc.createTrailCore(ctx, store, CreateTrailInput{
			Name: "topic-arn-trail", S3BucketName: "dest-bucket", SnsTopicName: topicARN, Region: store.GetRegion(),
		})
		require.NoError(t, err)
		assert.Equal(t, topicARN, created.SnsTopicARN)
	})

	t.Run("non-sns ARN reference is an invalid topic name", func(t *testing.T) {
		_, err := svc.createTrailCore(ctx, store, CreateTrailInput{
			Name: "bad-arn-trail", S3BucketName: "dest-bucket",
			SnsTopicName: "arn:aws:sqs:us-east-1:acc123:queue", Region: store.GetRegion(),
		})
		requireAWSCode(t, err, "InvalidSnsTopicNameException", 400)
	})
}

func TestUpdateTrailDestinationValidation(t *testing.T) {
	store := newErrorTestServiceStore(t)
	ctx := context.Background()

	s3 := newFakeS3Invoker("upd-bucket")
	svc := NewCloudTrailService(store.GetAccountID(), store.GetRegion())
	svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3})

	created, err := svc.createTrailCore(ctx, store, CreateTrailInput{
		Name: "upd-dest-trail", S3BucketName: "upd-bucket", Region: store.GetRegion(),
	})
	require.NoError(t, err)

	t.Run("bucket change to a missing bucket is rejected", func(t *testing.T) {
		_, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{
			Name:   created.Name,
			Params: map[string]interface{}{"S3BucketName": "ghost-bucket"},
		})
		requireAWSCode(t, err, "S3BucketDoesNotExistException", 404)
	})

	t.Run("malformed bucket name is rejected as a name, not a destination", func(t *testing.T) {
		_, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{
			Name:   created.Name,
			Params: map[string]interface{}{"S3BucketName": "..invalid.."},
		})
		requireAWSCode(t, err, "InvalidS3BucketNameException", 400)
	})

	t.Run("bucket change to an insufficiently privileged bucket is rejected", func(t *testing.T) {
		s3.buckets["plain-bucket"] = map[string][]byte{}
		_, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{
			Name:   created.Name,
			Params: map[string]interface{}{"S3BucketName": "plain-bucket"},
		})
		requireAWSCode(t, err, "InsufficientS3BucketPolicyException", 403)
	})

	t.Run("prefix change covered by the policy succeeds", func(t *testing.T) {
		s3.policies["upd-bucket"] = sufficientAnyBucketPolicy("upd-bucket")
		_, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{
			Name:   created.Name,
			Params: map[string]interface{}{"S3KeyPrefix": "logs"},
		})
		require.NoError(t, err)
	})

	t.Run("sns topic update resolves the ARN and clear empties it", func(t *testing.T) {
		sns := &fakeSNSInvoker{topics: map[string]string{}}
		topicARN := "arn:aws:sns:us-east-1:acc123:upd-notify"
		sns.topics[topicARN] = sufficientTopicPolicy(topicARN)
		svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3, sns: sns})

		updated, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{
			Name:   created.Name,
			Params: map[string]interface{}{"SnsTopicName": "upd-notify"},
		})
		require.NoError(t, err)
		assert.Equal(t, topicARN, updated.SnsTopicARN)

		cleared, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{
			Name:   created.Name,
			Params: map[string]interface{}{"SnsTopicName": ""},
		})
		require.NoError(t, err)
		assert.Empty(t, cleared.SnsTopicName)
		assert.Empty(t, cleared.SnsTopicARN)
	})

	t.Run("recorded topic without an ARN is healed on update", func(t *testing.T) {
		sns := &fakeSNSInvoker{topics: map[string]string{}}
		topicARN := "arn:aws:sns:us-east-1:acc123:heal-notify"
		sns.topics[topicARN] = sufficientTopicPolicy(topicARN)
		svc.SetInvokerRegistry(fakeInvokerRegistry{s3: s3, sns: sns})

		// Simulate a pre-resolution record: a topic name with no ARN.
		_, err := store.MutateTrail(created.Name, func(tr *cloudtrailstore.Trail) error {
			tr.SnsTopicName = "heal-notify"
			tr.SnsTopicARN = ""
			return nil
		})
		require.NoError(t, err)

		updated, err := svc.updateTrailCore(ctx, store, UpdateTrailInput{
			Name:   created.Name,
			Params: map[string]interface{}{},
		})
		require.NoError(t, err)
		assert.Equal(t, topicARN, updated.SnsTopicARN)
	})
}

func TestCloudWatchLogsStreamNaming(t *testing.T) {
	assert.Equal(t, "123456789012_CloudTrail_us-east-2",
		fmt.Sprintf(cloudTrailLogStreamFormat, "123456789012", "us-east-2"))
	group, region := cloudWatchLogsDestination(
		"arn:aws:logs:us-east-2:123456789012:log-group:CloudTrail/logs", "us-east-1")
	assert.Equal(t, "CloudTrail/logs", group)
	assert.Equal(t, "us-east-2", region)
	// An unparseable ARN keeps the textual extraction and addresses the
	// trail's home region.
	group, region = cloudWatchLogsDestination("log-group:CloudTrail/logs", "us-east-1")
	assert.Equal(t, "CloudTrail/logs", group)
	assert.Equal(t, "us-east-1", region)
}
