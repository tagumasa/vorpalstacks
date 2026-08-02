package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/timestreamquery"
	tqtypes "github.com/aws/aws-sdk-go-v2/service/timestreamquery/types"
)

func (r *TestRunner) runTimestreamScheduledTests(tc *tsTestContext) []TestResult {
	var results []TestResult
	sqName := tc.uniqueName("sq")
	sqDBName := tc.uniqueName("sq-db")
	sqTableName := tc.uniqueName("sq-table")
	sqRoleName := tc.uniqueName("ts-sq-role")
	var sqARN string

	defer func() {
		tc.deleteTable(sqDBName, sqTableName)
		tc.deleteDatabase(sqDBName)
		IAMDeleteRole(tc.iamClient, sqRoleName)
	}()

	createIAMRole := func(roleName string) error {
		return IAMCreateRole(tc.iamClient, roleName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"timestream.amazonaws.com"},"Action":"sts:AssumeRole"}]}`)
	}

	results = append(results, r.RunTest("timestream", "ScheduledQuery_Setup", func() error {
		if err := createIAMRole(sqRoleName); err != nil {
			return fmt.Errorf("create role: %v", err)
		}
		if err := tc.createDatabase(sqDBName); err != nil {
			return fmt.Errorf("create db: %v", err)
		}
		return tc.createTable(sqDBName, sqTableName)
	}))

	sqInput := &timestreamquery.CreateScheduledQueryInput{
		Name:        aws.String(sqName),
		QueryString: aws.String(fmt.Sprintf(`SELECT * FROM "%s"."%s"`, sqDBName, sqTableName)),
		ScheduleConfiguration: &tqtypes.ScheduleConfiguration{
			ScheduleExpression: aws.String("cron(0 0 * * ? *)"),
		},
		NotificationConfiguration: &tqtypes.NotificationConfiguration{
			SnsConfiguration: &tqtypes.SnsConfiguration{
				TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:test-topic", tc.region, tc.accountID)),
			},
		},
		ErrorReportConfiguration: &tqtypes.ErrorReportConfiguration{
			S3Configuration: &tqtypes.S3Configuration{
				BucketName: aws.String("error-report-bucket"),
			},
		},
		ScheduledQueryExecutionRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, sqRoleName)),
		Tags: []tqtypes.Tag{
			{Key: aws.String("env"), Value: aws.String("scheduled-test")},
		},
	}

	results = append(results, r.RunTest("timestream", "CreateScheduledQuery_Verify", func() error {
		resp, err := tc.queryClient.CreateScheduledQuery(tc.ctx, sqInput)
		if err != nil {
			return err
		}
		if resp.Arn == nil || *resp.Arn == "" {
			return fmt.Errorf("expected Arn to be set")
		}
		sqARN = *resp.Arn
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DescribeScheduledQuery_Verify", func() error {
		resp, err := tc.queryClient.DescribeScheduledQuery(tc.ctx, &timestreamquery.DescribeScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
		})
		if err != nil {
			return err
		}
		if resp.ScheduledQuery == nil {
			return fmt.Errorf("expected ScheduledQuery")
		}
		if resp.ScheduledQuery.Name == nil || *resp.ScheduledQuery.Name != sqName {
			return fmt.Errorf("expected Name=%s, got %v", sqName, resp.ScheduledQuery.Name)
		}
		if resp.ScheduledQuery.CreationTime == nil {
			return fmt.Errorf("expected CreationTime to be set")
		}
		if resp.ScheduledQuery.State == "" {
			return fmt.Errorf("expected State to be set")
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ListScheduledQueries_Verify", func() error {
		resp, err := tc.queryClient.ListScheduledQueries(tc.ctx, &timestreamquery.ListScheduledQueriesInput{})
		if err != nil {
			return err
		}
		if len(resp.ScheduledQueries) == 0 {
			return fmt.Errorf("expected at least 1 scheduled query")
		}
		found := false
		for _, sq := range resp.ScheduledQueries {
			if sq.Name != nil && *sq.Name == sqName {
				found = true
				if sq.Arn == nil || *sq.Arn == "" {
					return fmt.Errorf("listed scheduled query has nil/empty Arn")
				}
				break
			}
		}
		if !found {
			return fmt.Errorf("scheduled query %s not found in list", sqName)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UpdateScheduledQuery_Disable", func() error {
		_, err := tc.queryClient.UpdateScheduledQuery(tc.ctx, &timestreamquery.UpdateScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
			State:             tqtypes.ScheduledQueryStateDisabled,
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "UpdateScheduledQuery_Verify", func() error {
		resp, err := tc.queryClient.DescribeScheduledQuery(tc.ctx, &timestreamquery.DescribeScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
		})
		if err != nil {
			return err
		}
		if resp.ScheduledQuery.State != tqtypes.ScheduledQueryStateDisabled {
			return fmt.Errorf("expected state DISABLED, got %s", resp.ScheduledQuery.State)
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "ExecuteScheduledQuery", func() error {
		_, err := tc.queryClient.ExecuteScheduledQuery(tc.ctx, &timestreamquery.ExecuteScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
			InvocationTime:    aws.Time(time.Now().UTC()),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "TagResource_ScheduledQuery", func() error {
		_, err := tc.queryClient.TagResource(tc.ctx, &timestreamquery.TagResourceInput{
			ResourceARN: aws.String(sqARN),
			Tags: []tqtypes.Tag{
				{Key: aws.String("extra"), Value: aws.String("tag")},
			},
		})
		if err != nil {
			return err
		}
		listResp, err := tc.queryClient.ListTagsForResource(tc.ctx, &timestreamquery.ListTagsForResourceInput{
			ResourceARN: aws.String(sqARN),
		})
		if err != nil {
			return err
		}
		if len(listResp.Tags) < 2 {
			return fmt.Errorf("expected at least 2 tags after tagging, got %d", len(listResp.Tags))
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "UntagResource_ScheduledQuery", func() error {
		_, err := tc.queryClient.UntagResource(tc.ctx, &timestreamquery.UntagResourceInput{
			ResourceARN: aws.String(sqARN),
			TagKeys:     []string{"extra"},
		})
		if err != nil {
			return err
		}
		listResp, err := tc.queryClient.ListTagsForResource(tc.ctx, &timestreamquery.ListTagsForResourceInput{
			ResourceARN: aws.String(sqARN),
		})
		if err != nil {
			return err
		}
		for _, t := range listResp.Tags {
			if t.Key != nil && *t.Key == "extra" {
				return fmt.Errorf("tag 'extra' should have been removed")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("timestream", "DeleteScheduledQuery", func() error {
		_, err := tc.queryClient.DeleteScheduledQuery(tc.ctx, &timestreamquery.DeleteScheduledQueryInput{
			ScheduledQueryArn: aws.String(sqARN),
		})
		return err
	}))

	results = append(results, r.RunTest("timestream", "DescribeScheduledQuery_NonExistent", func() error {
		fakeARN := fmt.Sprintf("arn:aws:timestream:%s:%s:scheduled-query/nonexistent", tc.region, tc.accountID)
		_, err := tc.queryClient.DescribeScheduledQuery(tc.ctx, &timestreamquery.DescribeScheduledQueryInput{
			ScheduledQueryArn: aws.String(fakeARN),
		})
		return AssertErrorContains(err, "ResourceNotFoundException")
	}))

	results = append(results, r.RunTest("timestream", "CreateScheduledQuery_Duplicate", func() error {
		dupSQName := tc.uniqueName("dup-sq")
		dupInput := &timestreamquery.CreateScheduledQueryInput{
			Name:        aws.String(dupSQName),
			QueryString: aws.String("SELECT 1"),
			ScheduleConfiguration: &tqtypes.ScheduleConfiguration{
				ScheduleExpression: aws.String("cron(0 0 * * ? *)"),
			},
			NotificationConfiguration: &tqtypes.NotificationConfiguration{
				SnsConfiguration: &tqtypes.SnsConfiguration{
					TopicArn: aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:test-topic", tc.region, tc.accountID)),
				},
			},
			ErrorReportConfiguration: &tqtypes.ErrorReportConfiguration{
				S3Configuration: &tqtypes.S3Configuration{
					BucketName: aws.String("error-report-bucket"),
				},
			},
			ScheduledQueryExecutionRoleArn: aws.String(fmt.Sprintf("arn:aws:iam::%s:role/%s", tc.accountID, sqRoleName)),
		}
		_, err := tc.queryClient.CreateScheduledQuery(tc.ctx, dupInput)
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer tc.queryClient.DeleteScheduledQuery(tc.ctx, &timestreamquery.DeleteScheduledQueryInput{
			ScheduledQueryArn: aws.String(fmt.Sprintf("arn:aws:timestream:%s:%s:scheduled-query/%s", tc.region, tc.accountID, dupSQName)),
		})
		_, err = tc.queryClient.CreateScheduledQuery(tc.ctx, dupInput)
		return AssertErrorContains(err, "ConflictException")
	}))

	return results
}
