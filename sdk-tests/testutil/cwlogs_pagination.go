package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func (tc *cwlogsTestCtx) paginationTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "DescribeLogGroups_Pagination", func() error {
		pgPrefix := fmt.Sprintf("PagGroup-%d-", time.Now().UnixNano())
		groupNames := []string{
			pgPrefix + "alpha",
			pgPrefix + "beta",
			pgPrefix + "gamma",
		}
		for _, name := range groupNames {
			if err := tc.createLogGroup(name); err != nil {
				return fmt.Errorf("create %s: %v", name, err)
			}
			defer tc.deleteLogGroup(name)
		}

		var allGroups []types.LogGroup
		var nextToken *string
		for {
			resp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
				LogGroupNamePrefix: aws.String(pgPrefix),
				Limit:              aws.Int32(2),
				NextToken:          nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe page: %v", err)
			}
			allGroups = append(allGroups, resp.LogGroups...)
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}

		if len(allGroups) != 3 {
			return fmt.Errorf("expected 3 groups across pages, got %d", len(allGroups))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeLogStreams_Pagination", func() error {
		psName, cleanupGroup, err := tc.newLogGroupFixture("PagStreamGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		streamNames := []string{psName + "-s1", psName + "-s2", psName + "-s3"}
		for _, sn := range streamNames {
			if err := tc.createLogStream(psName, sn); err != nil {
				return fmt.Errorf("create stream %s: %v", sn, err)
			}
		}

		var allStreams []types.LogStream
		var nextToken *string
		for {
			resp, err := tc.client.DescribeLogStreams(tc.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
				LogGroupName: aws.String(psName),
				Limit:        aws.Int32(2),
				NextToken:    nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe page: %v", err)
			}
			allStreams = append(allStreams, resp.LogStreams...)
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}

		if len(allStreams) != 3 {
			return fmt.Errorf("expected 3 streams across pages, got %d", len(allStreams))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeLogStreams_NamePrefix", func() error {
		npName, cleanupGroup, err := tc.newLogGroupFixture("NpStreamGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		for _, sn := range []string{"app-server-1", "app-server-2", "db-server-1"} {
			if err := tc.createLogStream(npName, sn); err != nil {
				return fmt.Errorf("create stream %s: %v", sn, err)
			}
		}

		resp, err := tc.client.DescribeLogStreams(tc.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(npName),
			LogStreamNamePrefix: aws.String("app-"),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.LogStreams) != 2 {
			return fmt.Errorf("expected 2 streams with prefix 'app-', got %d", len(resp.LogStreams))
		}
		for _, ls := range resp.LogStreams {
			if !strings.HasPrefix(*ls.LogStreamName, "app-") {
				return fmt.Errorf("unexpected stream: %q", *ls.LogStreamName)
			}
		}

		// The identifier member addresses the same listing by ARN.
		npArn, err := tc.findLogGroupARN(npName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		arnResp, err := tc.client.DescribeLogStreams(tc.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupIdentifier:  npArn,
			LogStreamNamePrefix: aws.String("app-"),
		})
		if err != nil {
			return fmt.Errorf("describe by identifier ARN: %v", err)
		}
		if len(arnResp.LogStreams) != 2 {
			return fmt.Errorf("identifier-ARN describe expected 2 streams, got %d", len(arnResp.LogStreams))
		}

		// A prefix carrying '/' addresses streams whose names carry the
		// same separator.
		for _, sn := range []string{"app/alpha", "app/beta", "db/one"} {
			if err := tc.createLogStream(npName, sn); err != nil {
				return fmt.Errorf("create stream %s: %v", sn, err)
			}
		}
		slashResp, err := tc.client.DescribeLogStreams(tc.ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(npName),
			LogStreamNamePrefix: aws.String("app/"),
		})
		if err != nil {
			return fmt.Errorf("describe with slash prefix: %v", err)
		}
		if len(slashResp.LogStreams) != 2 {
			return fmt.Errorf("expected 2 streams with prefix 'app/', got %d", len(slashResp.LogStreams))
		}
		for _, ls := range slashResp.LogStreams {
			if !strings.HasPrefix(*ls.LogStreamName, "app/") {
				return fmt.Errorf("unexpected stream under slash prefix: %q", *ls.LogStreamName)
			}
		}
		return nil
	}))

	return results
}
