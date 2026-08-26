package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

func (tc *athenaTestContext) testWorkGroups() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("athena", "ListWorkGroups", func() error {
		resp, err := tc.client.ListWorkGroups(tc.ctx, &athena.ListWorkGroupsInput{
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

	workGroupName := tc.uniqueName("test-wg")
	// Created here and kept alive for the Get/Update/Delete scenarios
	// below; the DeleteWorkGroup scenario performs the actual deletion.
	results = append(results, tc.runner.RunTest("athena", "CreateWorkGroup", func() error {
		return tc.createWorkGroup(workGroupName, &types.WorkGroupConfiguration{
			ResultConfiguration: &types.ResultConfiguration{
				OutputLocation: aws.String("s3://test-bucket/athena/"),
			},
		})
	}))

	results = append(results, tc.runner.RunTest("athena", "GetWorkGroup", func() error {
		resp, err := tc.client.GetWorkGroup(tc.ctx, &athena.GetWorkGroupInput{
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
		_, err := tc.client.UpdateWorkGroup(tc.ctx, &athena.UpdateWorkGroupInput{
			WorkGroup:   aws.String(workGroupName),
			Description: aws.String("Updated work group"),
		})
		if err != nil {
			return err
		}
		verifyResp, err := tc.client.GetWorkGroup(tc.ctx, &athena.GetWorkGroupInput{
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
		_, err := tc.client.DeleteWorkGroup(tc.ctx, &athena.DeleteWorkGroupInput{
			WorkGroup:             aws.String(workGroupName),
			RecursiveDeleteOption: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetWorkGroup(tc.ctx, &athena.GetWorkGroupInput{
			WorkGroup: aws.String(workGroupName),
		})
		if err == nil {
			return fmt.Errorf("work group should be deleted but GetWorkGroup succeeded")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListWorkGroups_Pagination", func() error {
		prefix := tc.uniqueName("PagWG")
		var created []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("%s-%d", prefix, i)
			if err := tc.createWorkGroup(name, nil); err != nil {
				for _, wn := range created {
					tc.deleteWorkGroup(wn)
				}
				return fmt.Errorf("create work group %s: %v", name, err)
			}
			created = append(created, name)
		}
		defer func() {
			for _, wn := range created {
				tc.deleteWorkGroup(wn)
			}
		}()

		allWGs, err := tc.allWorkGroups()
		if err != nil {
			return err
		}
		var found []string
		for _, wg := range allWGs {
			if wg.Name != nil && strings.HasPrefix(aws.ToString(wg.Name), prefix+"-") {
				found = append(found, aws.ToString(wg.Name))
			}
		}
		if len(found) != len(created) {
			return fmt.Errorf("expected %d paginated work groups across pages, found %d", len(created), len(found))
		}
		for _, want := range created {
			foundIt := false
			for _, got := range found {
				if got == want {
					foundIt = true
					break
				}
			}
			if !foundIt {
				return fmt.Errorf("created work group %s not found across pages", want)
			}
		}
		return nil
	}))

	return results
}
