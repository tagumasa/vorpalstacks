package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
)

func (tc *athenaTestCtx) testEdgeCases() []TestResult {
	var results []TestResult
	client := tc.client
	ctx := tc.ctx

	results = append(results, tc.runner.RunTest("athena", "GetWorkGroup_NonExistent", func() error {
		_, err := client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
			WorkGroup: aws.String("nonexistent_wg_xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "DeleteWorkGroup_NonExistent", func() error {
		_, err := client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{
			WorkGroup: aws.String("nonexistent_wg_xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetNamedQuery_NonExistent", func() error {
		_, err := client.GetNamedQuery(ctx, &athena.GetNamedQueryInput{
			NamedQueryId: aws.String("00000000-0000-0000-0000-000000000000"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "GetDataCatalog_NonExistent", func() error {
		_, err := client.GetDataCatalog(ctx, &athena.GetDataCatalogInput{
			Name: aws.String("nonexistent_catalog_xyz"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "CreateWorkGroup_Duplicate", func() error {
		dupWGName := fmt.Sprintf("dup-wg-%d", time.Now().UnixNano()%1000000000)
		_, err := client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(dupWGName),
		})
		if err != nil {
			return fmt.Errorf("first create: %v", err)
		}
		defer client.DeleteWorkGroup(ctx, &athena.DeleteWorkGroupInput{WorkGroup: aws.String(dupWGName)})

		_, err = client.CreateWorkGroup(ctx, &athena.CreateWorkGroupInput{
			Name: aws.String(dupWGName),
		})
		if err == nil {
			return fmt.Errorf("expected error for duplicate work group")
		}
		if err := AssertErrorContains(err, "ResourceAlreadyExistsException"); err != nil {
			return fmt.Errorf("expected ResourceAlreadyExistsException for duplicate, got: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("athena", "ListEngineVersions", func() error {
		resp, err := client.ListEngineVersions(ctx, &athena.ListEngineVersionsInput{})
		if err != nil {
			return err
		}
		if len(resp.EngineVersions) == 0 {
			return fmt.Errorf("expected at least 1 engine version")
		}
		for _, ev := range resp.EngineVersions {
			if aws.ToString(ev.SelectedEngineVersion) == "" {
				return fmt.Errorf("engine version with empty selected version string found")
			}
		}
		return nil
	}))

	return results
}
