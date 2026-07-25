package testutil

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

func (r *TestRunner) runCloudTrailImportTests(tc *cloudTrailTestContext) []TestResult {
	var results []TestResult

	// Create a dedicated EDS for import tests.
	var edsARN string
	var edsID string
	results = append(results, r.RunTest("cloudtrail", "Import_Setup_EDS", func() error {
		resp, err := tc.client.CreateEventDataStore(tc.ctx, &cloudtrail.CreateEventDataStoreInput{
			Name:                         aws.String(fmt.Sprintf("import-test-%d", time.Now().UnixNano())),
			TerminationProtectionEnabled: aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("create EDS for import: %w", err)
		}
		if resp.EventDataStoreArn == nil {
			return fmt.Errorf("EDS ARN is nil")
		}
		edsARN = *resp.EventDataStoreArn
		edsID = *resp.EventDataStoreArn
		return nil
	}))
	defer func() {
		if edsARN != "" {
			_, _ = tc.client.DeleteEventDataStore(tc.ctx, &cloudtrail.DeleteEventDataStoreInput{
				EventDataStore: aws.String(edsID),
			})
		}
	}()

	// StartImport.
	var importID string
	results = append(results, r.RunTest("cloudtrail", "StartImport_Success", func() error {
		resp, err := tc.client.StartImport(tc.ctx, &cloudtrail.StartImportInput{
			Destinations: []string{edsARN},
			ImportSource: &types.ImportSource{
				S3: &types.S3ImportSource{
					S3LocationUri:         aws.String("s3://test-bucket/CloudTrail/"),
					S3BucketRegion:        aws.String("us-east-1"),
					S3BucketAccessRoleArn: aws.String("arn:aws:iam::123456789012:role/CloudTrailImport"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("StartImport failed: %w", err)
		}
		if resp.ImportId == nil {
			return fmt.Errorf("ImportId is nil")
		}
		importID = *resp.ImportId
		return nil
	}))

	// GetImport.
	results = append(results, r.RunTest("cloudtrail", "GetImport_Success", func() error {
		resp, err := tc.client.GetImport(tc.ctx, &cloudtrail.GetImportInput{
			ImportId: aws.String(importID),
		})
		if err != nil {
			return fmt.Errorf("GetImport failed: %w", err)
		}
		if resp.ImportId == nil || *resp.ImportId != importID {
			return fmt.Errorf("ImportId mismatch")
		}
		if resp.ImportStatus == "" {
			return fmt.Errorf("ImportStatus is empty")
		}
		return nil
	}))

	// GetImport_NotFound.
	results = append(results, r.RunTest("cloudtrail", "GetImport_NotFound", func() error {
		_, err := tc.client.GetImport(tc.ctx, &cloudtrail.GetImportInput{
			ImportId: aws.String("nonexistent-import-id"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent import")
		}
		return nil
	}))

	// ListImports.
	results = append(results, r.RunTest("cloudtrail", "ListImports_Success", func() error {
		resp, err := tc.client.ListImports(tc.ctx, &cloudtrail.ListImportsInput{
			Destination: aws.String(edsARN),
		})
		if err != nil {
			return fmt.Errorf("ListImports failed: %w", err)
		}
		found := false
		for _, imp := range resp.Imports {
			if imp.ImportId != nil && *imp.ImportId == importID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("import %s not found in list", importID)
		}
		return nil
	}))

	// ListImportFailures.
	results = append(results, r.RunTest("cloudtrail", "ListImportFailures_Success", func() error {
		_, err := tc.client.ListImportFailures(tc.ctx, &cloudtrail.ListImportFailuresInput{
			ImportId: aws.String(importID),
		})
		if err != nil {
			return fmt.Errorf("ListImportFailures failed: %w", err)
		}
		return nil
	}))

	// StartImport_MissingDestinations.
	results = append(results, r.RunTest("cloudtrail", "StartImport_MissingDestinations", func() error {
		_, err := tc.client.StartImport(tc.ctx, &cloudtrail.StartImportInput{
			ImportSource: &types.ImportSource{
				S3: &types.S3ImportSource{
					S3LocationUri:         aws.String("s3://test-bucket/CloudTrail/"),
					S3BucketRegion:        aws.String("us-east-1"),
					S3BucketAccessRoleArn: aws.String("arn:aws:iam::123456789012:role/CloudTrailImport"),
				},
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for missing Destinations")
		}
		return nil
	}))

	return results
}
