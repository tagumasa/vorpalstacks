package testutil

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaLayerTests(tc *lambdaTestContext) []TestResult {
	var results []TestResult

	layerName := tc.unique("TestLayer")
	layerZipContent := base64.StdEncoding.EncodeToString([]byte("exports.handler = async (event) => { return 1; };"))

	results = append(results, tc.r.RunTest("lambda", "PublishLayerVersion", func() error {
		resp, err := tc.client.PublishLayerVersion(tc.ctx, &lambda.PublishLayerVersionInput{
			LayerName: aws.String(layerName),
			Content: &types.LayerVersionContentInput{
				ZipFile: []byte(layerZipContent),
			},
			Description:        aws.String("Test layer version"),
			CompatibleRuntimes: []types.Runtime{types.RuntimeNodejs22x},
		})
		if err != nil {
			return err
		}
		if resp.LayerArn == nil {
			return fmt.Errorf("LayerArn is nil")
		}
		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}
		if resp.Content == nil || resp.Content.CodeSha256 == nil {
			return fmt.Errorf("CodeSha256 is nil")
		}
		// Two more versions so pagination has multiple pages.
		for i := 2; i <= 3; i++ {
			if _, err := tc.client.PublishLayerVersion(tc.ctx, &lambda.PublishLayerVersionInput{
				LayerName: aws.String(layerName),
				Content:   &types.LayerVersionContentInput{ZipFile: []byte(layerZipContent)},
			}); err != nil {
				return fmt.Errorf("publish version %d: %v", i, err)
			}
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListLayerVersions_Pagination", func() error {
		first, err := tc.client.ListLayerVersions(tc.ctx, &lambda.ListLayerVersionsInput{
			LayerName: aws.String(layerName),
			MaxItems:  aws.Int32(2),
		})
		if err != nil {
			return err
		}
		if len(first.LayerVersions) != 2 {
			return fmt.Errorf("first page should carry 2 versions, got %d", len(first.LayerVersions))
		}
		if first.NextMarker == nil || *first.NextMarker == "" {
			return fmt.Errorf("first page should be truncated with a NextMarker")
		}
		second, err := tc.client.ListLayerVersions(tc.ctx, &lambda.ListLayerVersionsInput{
			LayerName: aws.String(layerName),
			Marker:    first.NextMarker,
		})
		if err != nil {
			return err
		}
		if len(second.LayerVersions) == 0 {
			return fmt.Errorf("second page should carry the remaining version")
		}
		for _, v := range second.LayerVersions {
			for _, seen := range first.LayerVersions {
				if seen.Version == v.Version {
					return fmt.Errorf("version %d repeated across pages", v.Version)
				}
			}
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetLayerVersion", func() error {
		resp, err := tc.client.GetLayerVersion(tc.ctx, &lambda.GetLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(1),
		})
		if err != nil {
			return err
		}
		if resp.Content == nil || resp.Content.CodeSha256 == nil {
			return fmt.Errorf("CodeSha256 is nil")
		}
		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "ListLayers", func() error {
		var nextMarker *string
		for page := 0; page < 100; page++ {
			resp, err := tc.client.ListLayers(tc.ctx, &lambda.ListLayersInput{
				Marker: nextMarker,
			})
			if err != nil {
				return err
			}
			for _, l := range resp.Layers {
				if l.LayerName != nil && *l.LayerName == layerName {
					return nil
				}
			}
			if resp.NextMarker == nil {
				break
			}
			nextMarker = resp.NextMarker
		}
		return fmt.Errorf("layer %s not found in ListLayers", layerName)
	}))

	results = append(results, tc.r.RunTest("lambda", "ListLayerVersions", func() error {
		resp, err := tc.client.ListLayerVersions(tc.ctx, &lambda.ListLayerVersionsInput{
			LayerName: aws.String(layerName),
		})
		if err != nil {
			return err
		}
		if resp.LayerVersions == nil {
			return fmt.Errorf("layer versions list is nil")
		}
		if len(resp.LayerVersions) == 0 {
			return fmt.Errorf("expected at least 1 layer version")
		}
		if resp.LayerVersions[0].Version != 1 {
			return fmt.Errorf("expected first version 1, got %d", resp.LayerVersions[0].Version)
		}
		if resp.LayerVersions[0].CompatibleRuntimes == nil || len(resp.LayerVersions[0].CompatibleRuntimes) == 0 {
			return fmt.Errorf("CompatibleRuntimes is nil or empty")
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "DeleteLayerVersion", func() error {
		_, err := tc.client.DeleteLayerVersion(tc.ctx, &lambda.DeleteLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(1),
		})
		if err != nil {
			return err
		}
		_, err = tc.client.GetLayerVersion(tc.ctx, &lambda.GetLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(1),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "GetLayerVersion_NonExistent", func() error {
		_, err := tc.client.GetLayerVersion(tc.ctx, &lambda.GetLayerVersionInput{
			LayerName:     aws.String("nonexistent-layer-xyz"),
			VersionNumber: aws.Int64(999),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, tc.r.RunTest("lambda", "LayerVersionPermission_OlderVersion", func() error {
		// Version 3 is the latest; version 2 is an older published version.
		_, err := tc.client.AddLayerVersionPermission(tc.ctx, &lambda.AddLayerVersionPermissionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(2),
			StatementId:   aws.String("cross-account"),
			Action:        aws.String("lambda:GetLayerVersion"),
			Principal:     aws.String("*"),
		})
		if err != nil {
			return fmt.Errorf("add permission on older version: %v", err)
		}

		policyResp, err := tc.client.GetLayerVersionPolicy(tc.ctx, &lambda.GetLayerVersionPolicyInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(2),
		})
		if err != nil {
			return err
		}
		if policyResp.Policy == nil || !strings.Contains(*policyResp.Policy, "cross-account") {
			return fmt.Errorf("policy should carry the added statement")
		}
		if policyResp.RevisionId == nil || *policyResp.RevisionId == "" {
			return fmt.Errorf("RevisionId is nil or empty")
		}
		again, err := tc.client.GetLayerVersionPolicy(tc.ctx, &lambda.GetLayerVersionPolicyInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(2),
		})
		if err != nil {
			return err
		}
		if again.RevisionId == nil || policyResp.RevisionId == nil || *again.RevisionId != *policyResp.RevisionId {
			return fmt.Errorf("RevisionId should be stable across reads")
		}

		if _, err := tc.client.RemoveLayerVersionPermission(tc.ctx, &lambda.RemoveLayerVersionPermissionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(2),
			StatementId:   aws.String("cross-account"),
		}); err != nil {
			return fmt.Errorf("remove permission: %v", err)
		}
		return nil
	}))

	return results
}
