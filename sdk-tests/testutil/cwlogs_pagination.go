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
			tc.createLogGroup(name)
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
		psName := tc.uniquePrefix("PagStreamGroup")
		if err := tc.createLogGroup(psName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(psName)

		streamNames := []string{psName + "-s1", psName + "-s2", psName + "-s3"}
		for _, sn := range streamNames {
			tc.createLogStream(psName, sn)
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
		npName := tc.uniquePrefix("NpStreamGroup")
		if err := tc.createLogGroup(npName); err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer tc.deleteLogGroup(npName)

		tc.createLogStream(npName, "app-server-1")
		tc.createLogStream(npName, "app-server-2")
		tc.createLogStream(npName, "db-server-1")

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
		return nil
	}))

	return results
}
