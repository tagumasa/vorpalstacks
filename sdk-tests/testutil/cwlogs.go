package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"vorpalstacks-sdk-tests/config"
)

type cwlogsTestCtx struct {
	client    *cloudwatchlogs.Client
	iamClient *iam.Client
	ctx       context.Context
	runner    *TestRunner
	region    string
}

func (r *TestRunner) RunCloudWatchLogsTests() []TestResult {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{
			Service:  "logs",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		}}
	}

	tc := &cwlogsTestCtx{
		client:    cloudwatchlogs.NewFromConfig(cfg),
		iamClient: iam.NewFromConfig(cfg),
		ctx:       context.Background(),
		runner:    r,
		region:    r.region,
	}

	return tc.runAll()
}

func (tc *cwlogsTestCtx) runAll() []TestResult {
	var results []TestResult

	results = append(results, tc.basicTests()...)
	results = append(results, tc.eventTests()...)
	results = append(results, tc.metricFilterTests()...)
	results = append(results, tc.subscriptionTests()...)
	results = append(results, tc.destinationTests()...)
	results = append(results, tc.retentionTests()...)
	results = append(results, tc.tagTests()...)
	results = append(results, tc.paginationTests()...)
	results = append(results, tc.legacyTests()...)

	return results
}

func (tc *cwlogsTestCtx) uniquePrefix(base string) string {
	return fmt.Sprintf("%s-%d", base, time.Now().UnixNano())
}

func (tc *cwlogsTestCtx) createLogGroup(name string) error {
	_, err := tc.client.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(name),
	})
	return err
}

func (tc *cwlogsTestCtx) deleteLogGroup(name string) {
	tc.client.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
		LogGroupName: aws.String(name),
	})
}

func (tc *cwlogsTestCtx) createLogStream(group, stream string) error {
	_, err := tc.client.CreateLogStream(tc.ctx, &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
	})
	return err
}

func (tc *cwlogsTestCtx) putLogEvent(group, stream, message string, ts int64) error {
	_, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  aws.String(group),
		LogStreamName: aws.String(stream),
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(ts),
			},
		},
	})
	return err
}

func (tc *cwlogsTestCtx) findLogGroupARN(prefix string) (*string, error) {
	resp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(prefix),
	})
	if err != nil {
		return nil, fmt.Errorf("describe: %v", err)
	}
	if len(resp.LogGroups) == 0 || resp.LogGroups[0].Arn == nil {
		return nil, fmt.Errorf("log group ARN not found for prefix %s", prefix)
	}
	return resp.LogGroups[0].Arn, nil
}

func (tc *cwlogsTestCtx) createSubscriptionRole() (string, func(), error) {
	roleName := tc.uniquePrefix("test-sub-role")
	_, err := tc.iamClient.CreateRole(tc.ctx, &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"logs.amazonaws.com"},"Action":"sts:AssumeRole"}]}`),
	})
	if err != nil {
		return "", nil, fmt.Errorf("create role: %v", err)
	}
	cleanup := func() {
		tc.iamClient.DeleteRole(tc.ctx, &iam.DeleteRoleInput{RoleName: aws.String(roleName)})
	}
	return fmt.Sprintf("arn:aws:iam::000000000000:role/%s", roleName), cleanup, nil
}

func (tc *cwlogsTestCtx) putMetricFilter(group, filterName, pattern, metricName, metricNS string) error {
	_, err := tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
		LogGroupName:  aws.String(group),
		FilterName:    aws.String(filterName),
		FilterPattern: aws.String(pattern),
		MetricTransformations: []types.MetricTransformation{
			{
				MetricName:      aws.String(metricName),
				MetricNamespace: aws.String(metricNS),
				MetricValue:     aws.String("1"),
			},
		},
	})
	return err
}

func (tc *cwlogsTestCtx) putSubscriptionFilter(group, filterName, pattern, destARN, roleARN string) error {
	_, err := tc.client.PutSubscriptionFilter(tc.ctx, &cloudwatchlogs.PutSubscriptionFilterInput{
		LogGroupName:   aws.String(group),
		FilterName:     aws.String(filterName),
		FilterPattern:  aws.String(pattern),
		DestinationArn: aws.String(destARN),
		RoleArn:        aws.String(roleARN),
		Distribution:   types.DistributionByLogStream,
	})
	return err
}

func (tc *cwlogsTestCtx) collectAllLogGroups(prefix string) ([]types.LogGroup, error) {
	var all []types.LogGroup
	var nextToken *string
	for {
		resp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String(prefix),
			NextToken:          nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("describe page: %v", err)
		}
		all = append(all, resp.LogGroups...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *cwlogsTestCtx) collectAllLogStreams(group string) ([]types.LogStream, error) {
	var all []types.LogStream
	var nextToken *string
	for {
		resp, err := tc.client.DescribeLogStreams(tc.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(group),
			NextToken:    nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("describe page: %v", err)
		}
		all = append(all, resp.LogStreams...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}
