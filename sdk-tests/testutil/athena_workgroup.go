package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestCtx) testWorkGroups() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("athena", "ListWorkGroups", func() error {
		resp, err := client.ListWorkGroups(ctx, &athena.ListWorkGroupsInput{
			MaxResults: aws.Int32(10),
		})
		if err != nil {
			return err
		}
		if resp.WorkGroups == nil {
			return fmt.Errorf("work groups list is nil")
		}
		var foundPrimary bool
		for _, wg := range resp.WorkGroups {
			if aws.ToString(wg.Name) == "primary" {
				foundPrimary = true
				if wg.State != types.WorkGroupStateEnabled {
					return fmt.Errorf("primary work group state: expected ENABLED, got %s", wg.State)
				}
			}
		}
		if !foundPrimary {
			return fmt.Errorf("primary work group not found in list")
		}
		return nil
	}))

	workGroupName := fmt.Sprintf("test-wg-%d", time.Now().UnixNano()%1000000000)
	results = append(results, tc.runner.RunTest("athena", "CreateWorkGroup", func() error {
		resp, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(workGroupName),
			Configuration: &types.WorkGroupConfiguration{
				ResultConfiguration: &types.ResultConfiguration{
					OutputLocation: aws.String("s3://test-bucket/athena/"),
				},
			},
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetWorkGroup", func() error {
		resp, err := client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
			WorkGroup: aws.String(workGroupName),
		})
		if err != nil {
			return err
		}
		wg := resp.WorkGroup
		if wg == nil {
			return fmt.Errorf("work group is nil")
		}
		if aws.ToString(wg.Name) != workGroupName {
			return fmt.Errorf("expected name %q, got %q", workGroupName, aws.ToString(wg.Name))
		}
		if wg.State != types.WorkGroupStateEnabled {
			return fmt.Errorf("expected state ENABLED, got %s", wg.State)
		}
		if wg.Configuration == nil || wg.Configuration.ResultConfiguration == nil {
			return fmt.Errorf("configuration or result configuration is nil")
		}
		if aws.ToString(wg.Configuration.ResultConfiguration.OutputLocation) != "s3://test-bucket/athena/" {
			return fmt.Errorf("expected output location 's3://test-bucket/athena/', got %q", aws.ToString(wg.Configuration.ResultConfiguration.OutputLocation))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "UpdateWorkGroup", func() error {
		_, err := client.UpdateWorkGroup(ctx, &athena.UpdateWorkGroupInput{
			WorkGroup:   aws.String(workGroupName),
			Description: aws.String("Updated work group"),
		})
		if err != nil {
			return err
		}
		verifyResp, err := client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
			WorkGroup: aws.String(workGroupName),
		})
		if err != nil {
			return err
		}
		if aws.ToString(verifyResp.WorkGroup.Description) != "Updated work group" {
			return fmt.Errorf("expected description 'Updated work group', got %q", aws.ToString(verifyResp.WorkGroup.Description))
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteWorkGroup", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup:             aws.String(workGroupName),
			RecursiveDeleteOption: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		_, err = client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
			WorkGroup: aws.String(workGroupName),
		})
		if err == nil {
			return fmt.Errorf("work group should be deleted but GetWorkGroup succeeded")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListWorkGroups_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano()%1000000000)
		var pgWorkGroups []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("PagWG-%s-%d", pgTs, i)
			_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
				Name: aws.String(name),
			})
			if err != nil {
				for _, wn := range pgWorkGroups {
					client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wn)})
				}
				return fmt.Errorf("create work group %s: %v", name, err)
			}
			pgWorkGroups = append(pgWorkGroups, name)
		}

		var allWorkGroups []string
		var nextToken *string
		for {
			resp, err := client.ListWorkGroups(ctx, &athena.ListWorkGroupsInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, wn := range pgWorkGroups {
					client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wn)})
				}
				return fmt.Errorf("list work groups page: %v", err)
			}
			for _, wg := range resp.WorkGroups {
				if wg.Name != nil && strings.Contains(*wg.Name, "PagWG-"+pgTs) {
					allWorkGroups = append(allWorkGroups, *wg.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, wn := range pgWorkGroups {
			client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(wn)})
		}
		if len(allWorkGroups) != 5 {
			return fmt.Errorf("expected 5 paginated work groups, got %d", len(allWorkGroups))
		}
		return nil
	}))

	return results
}
