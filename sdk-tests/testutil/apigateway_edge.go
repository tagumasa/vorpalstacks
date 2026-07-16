package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

func (r *TestRunner) runAPIGatewayEdgeTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("apigateway", "TagResource_UntagResource_ListTags", func() error {
		tagAPI := fmt.Sprintf("TagAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tagAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		arn := fmt.Sprintf("arn:aws:apigateway:%s::/restapis/%s", r.region, *createResp.Id)

		_, err = client.TagResource(ctx, &apigateway.TagResourceInput{
			ResourceArn: aws.String(arn),
			Tags: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		})
		if err != nil {
			return fmt.Errorf("tag: %v", err)
		}

		tagResp, err := client.GetTags(ctx, &apigateway.GetTagsInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("get tags: %v", err)
		}
		if tagResp.Tags == nil || tagResp.Tags["key1"] != "value1" {
			return fmt.Errorf("tags mismatch, got %v", tagResp.Tags)
		}

		_, err = client.UntagResource(ctx, &apigateway.UntagResourceInput{
			ResourceArn: aws.String(arn),
			TagKeys:     []string{"key2"},
		})
		if err != nil {
			return fmt.Errorf("untag: %v", err)
		}

		tagResp2, err := client.GetTags(ctx, &apigateway.GetTagsInput{
			ResourceArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("get tags after untag: %v", err)
		}
		if _, exists := tagResp2.Tags["key2"]; exists {
			return fmt.Errorf("key2 should have been removed")
		}
		if tagResp2.Tags["key1"] != "value1" {
			return fmt.Errorf("key1 should still exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetRestApi_NonExistent", func() error {
		_, err := client.GetRestApi(ctx, &apigateway.GetRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "DeleteRestApi_NonExistent", func() error {
		_, err := client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{
			RestApiId: aws.String("nonexistent_xyz"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("apigateway", "GetStage_NonExistent", func() error {
		tmpAPI := fmt.Sprintf("TmpAPI-%d", time.Now().UnixNano())
		createResp, err := client.CreateRestApi(ctx, &apigateway.CreateRestApiInput{
			Name: aws.String(tmpAPI),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteRestApi(ctx, &apigateway.DeleteRestApiInput{RestApiId: createResp.Id})

		_, err = client.GetStage(ctx, &apigateway.GetStageInput{
			RestApiId: createResp.Id,
			StageName: aws.String("nonexistent_stage"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	return results
}
