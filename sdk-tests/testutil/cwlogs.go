package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"vorpalstacks-sdk-tests/config"
)

type cwlogsTestCtx struct {
	client    *cloudwatchlogs.Client
	iamClient *iam.Client
	kmsClient *kms.Client
	// noRetryClient serves the calls that deliberately provoke a quota
	// rejection: LimitExceededException sits in the SDK standard
	// retryer's throttling class, so an ordinary client pays the full
	// retry ladder (three attempts with exponential backoff) before the
	// assertion sees the error it is waiting for.
	noRetryClient *cloudwatchlogs.Client
	ctx           context.Context
	runner        *TestRunner
	region        string
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
		kmsClient: kms.NewFromConfig(cfg),
		noRetryClient: cloudwatchlogs.NewFromConfig(cfg, func(o *cloudwatchlogs.Options) {
			o.RetryMaxAttempts = 1
		}),
		ctx:    context.Background(),
		runner: r,
		region: r.region,
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
	results = append(results, tc.validationMatrixTests()...)
	results = append(results, tc.paginationTests()...)
	results = append(results, tc.legacyTests()...)
	results = append(results, tc.queryTests()...)
	results = append(results, tc.lookupTableTests()...)
	results = append(results, tc.deliveryTests()...)
	results = append(results, tc.transformerTests()...)
	results = append(results, tc.indexPolicyTests()...)
	results = append(results, tc.e4Tests()...)
	results = append(results, tc.liveTailTests()...)
	results = append(results, tc.httpIngestionTests()...)

	return results
}

func (tc *cwlogsTestCtx) uniquePrefix(base string) string {
	return fmt.Sprintf("%s-%d", base, time.Now().UnixNano())
}

// roleARN builds an IAM role ARN against the runner's configured account,
// so fixtures survive account-configuration changes.
func (tc *cwlogsTestCtx) roleARN(name string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.runner.AccountID(), name)
}

// logGroupARN addresses a log group in the ARN form the members that
// mandate it take ("Specify each log group by its ARN", StartLiveTail).
func (tc *cwlogsTestCtx) logGroupARN(name string) string {
	return fmt.Sprintf("arn:aws:logs:%s:%s:log-group:%s", tc.region, tc.runner.AccountID(), name)
}

func (tc *cwlogsTestCtx) createLogGroup(name string) error {
	_, err := tc.client.CreateLogGroup(tc.ctx, &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(name),
	})
	return err
}

// newLogGroupFixture creates one isolated log group for a single test and
// returns its name with the cleanup that deletes it: the family keeps
// per-test isolation (one test's failure cannot cascade into its
// siblings), with the scaffold stated once instead of per test.
func (tc *cwlogsTestCtx) newLogGroupFixture(base string) (string, func(), error) {
	name := tc.uniquePrefix(base)
	cleanup := func() { tc.deleteLogGroup(name) }
	if err := tc.createLogGroup(name); err != nil {
		return name, cleanup, err
	}
	return name, cleanup, nil
}

// newGroupStreamFixture additionally creates one named log stream inside
// the isolated group, for the tests that read events back.
func (tc *cwlogsTestCtx) newGroupStreamFixture(base, stream string) (string, func(), error) {
	name, cleanup, err := tc.newLogGroupFixture(base)
	if err != nil {
		return name, cleanup, err
	}
	if err := tc.createLogStream(name, stream); err != nil {
		return name, cleanup, err
	}
	return name, cleanup, nil
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

// findLogGroupARN resolves the group's ARN in identifier form: the
// listing's ARN carries the AWS log-group suffix ":*", and the
// identifier members' documentation says "Don't include an * at the
// end" — every API this helper feeds takes the identifier form.
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
	return aws.String(strings.TrimSuffix(aws.ToString(resp.LogGroups[0].Arn), ":*")), nil
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
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.runner.AccountID(), roleName), cleanup, nil
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

// The collectAll* walkers below exist for the suite-wide full-traversal
// rule: during full regression other services create resources in
// parallel, so a membership assertion can never rely on the first page.

func (tc *cwlogsTestCtx) collectAllMetricFilters(group, prefix string) ([]types.MetricFilter, error) {
	var all []types.MetricFilter
	var nextToken *string
	for {
		resp, err := tc.client.DescribeMetricFilters(tc.ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName:     aws.String(group),
			FilterNamePrefix: aws.String(prefix),
			NextToken:        nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("describe page: %v", err)
		}
		all = append(all, resp.MetricFilters...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *cwlogsTestCtx) collectAllSubscriptionFilters(group, prefix string) ([]types.SubscriptionFilter, error) {
	var all []types.SubscriptionFilter
	var nextToken *string
	for {
		resp, err := tc.client.DescribeSubscriptionFilters(tc.ctx, &cloudwatchlogs.DescribeSubscriptionFiltersInput{
			LogGroupName:     aws.String(group),
			FilterNamePrefix: aws.String(prefix),
			NextToken:        nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("describe page: %v", err)
		}
		all = append(all, resp.SubscriptionFilters...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

func (tc *cwlogsTestCtx) collectAllDestinations(prefix string) ([]types.Destination, error) {
	var all []types.Destination
	var nextToken *string
	for {
		resp, err := tc.client.DescribeDestinations(tc.ctx, &cloudwatchlogs.DescribeDestinationsInput{
			DestinationNamePrefix: aws.String(prefix),
			NextToken:             nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("describe page: %v", err)
		}
		all = append(all, resp.Destinations...)
		if resp.NextToken == nil || *resp.NextToken == "" {
			break
		}
		nextToken = resp.NextToken
	}
	return all, nil
}

// waitForQuery polls GetQueryResults until the query reports Complete and
// returns that response (its Statistics and Results describe the finished
// query); Failed and Cancelled surface as errors. attempts caps the poll —
// the interval is the suite's usual 200ms, so attempts sizes the wait
// window a caller is willing to accept.
func (tc *cwlogsTestCtx) waitForQuery(queryID *string, attempts int) (*cloudwatchlogs.GetQueryResultsOutput, error) {
	for i := 0; i < attempts; i++ {
		resp, err := tc.client.GetQueryResults(tc.ctx, &cloudwatchlogs.GetQueryResultsInput{
			QueryId: queryID,
		})
		if err != nil {
			return nil, fmt.Errorf("get results: %v", err)
		}
		if resp.Status == types.QueryStatusComplete {
			return resp, nil
		}
		if resp.Status == types.QueryStatusFailed || resp.Status == types.QueryStatusCancelled {
			return nil, fmt.Errorf("query finished with status %s", resp.Status)
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil, fmt.Errorf("query did not complete in time")
}
