package testutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func runLambdaLayerTests(
	r *TestRunner,
	ctx context.Context,
	client *lambda.Client,
) []TestResult {
	var results []TestResult

	layerName := fmt.Sprintf("TestLayer-%d", time.Now().UnixNano())
	layerZipContent := base64.StdEncoding.EncodeToString([]byte("exports.handler = async (event) => { return 1; };"))

	results = append(results, r.RunTest("lambda", "PublishLayerVersion", func() error {
		resp, err := client.PublishLayerVersion(ctx, &lambda.PublishLayerVersionInput{
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
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetLayerVersion", func() error {
		resp, err := client.GetLayerVersion(ctx, &lambda.GetLayerVersionInput{
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

	results = append(results, r.RunTest("lambda", "ListLayers", func() error {
		resp, err := client.ListLayers(ctx, &lambda.ListLayersInput{})
		if err != nil {
			return err
		}
		if resp.Layers == nil {
			return fmt.Errorf("layers list is nil")
		}
		found := false
		for _, l := range resp.Layers {
			if l.LayerName != nil && *l.LayerName == layerName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("layer %s not found in ListLayers", layerName)
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "ListLayerVersions", func() error {
		resp, err := client.ListLayerVersions(ctx, &lambda.ListLayerVersionsInput{
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

	results = append(results, r.RunTest("lambda", "DeleteLayerVersion", func() error {
		_, err := client.DeleteLayerVersion(ctx, &lambda.DeleteLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(1),
		})
		if err != nil {
			return err
		}
		_, err = client.GetLayerVersion(ctx, &lambda.GetLayerVersionInput{
			LayerName:     aws.String(layerName),
			VersionNumber: aws.Int64(1),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("lambda", "GetLayerVersion_NonExistent", func() error {
		_, err := client.GetLayerVersion(ctx, &lambda.GetLayerVersionInput{
			LayerName:     aws.String("nonexistent-layer-xyz"),
			VersionNumber: aws.Int64(999),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
