package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptunegraph"
)

func (r *TestRunner) runNeptunegraphTagTests(tc *neptunegraphContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("neptunegraph", "TagResource", func() error {
		if err := tc.requireGraphWithARN(); err != nil {
			return err
		}
		_, err := tc.client.TagResource(tc.ctx, &neptunegraph.TagResourceInput{
			ResourceArn: aws.String(tc.graphARN),
			Tags: map[string]string{
				"ExtraTag":  "extra-value",
				"CreatedBy": "sdk-tests",
			},
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "TagResource_ReservedPrefixRejected", func() error {
		if err := tc.requireGraphARN(); err != nil {
			return err
		}
		_, err := tc.client.TagResource(tc.ctx, &neptunegraph.TagResourceInput{
			ResourceArn: aws.String(tc.graphARN),
			Tags: map[string]string{
				"aws:reserved": "v",
			},
		})
		if err == nil {
			return fmt.Errorf("expected error for aws:-prefixed tag key")
		}
		return nil
	}))

	// Tag values are limited to 256 Unicode characters, not bytes.
	results = append(results, r.RunTest("neptunegraph", "TagResource_MultibyteTagValueAccepted", func() error {
		if err := tc.requireGraphARN(); err != nil {
			return err
		}
		value := strings.Repeat("\u65e5", 200)
		_, err := tc.client.TagResource(tc.ctx, &neptunegraph.TagResourceInput{
			ResourceArn: aws.String(tc.graphARN),
			Tags: map[string]string{
				"MultibyteValue": value,
			},
		})
		if err != nil {
			return fmt.Errorf("200-character multibyte tag value rejected: %w", err)
		}
		resp, err := tc.client.ListTagsForResource(tc.ctx, &neptunegraph.ListTagsForResourceInput{
			ResourceArn: aws.String(tc.graphARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags["MultibyteValue"] != value {
			return fmt.Errorf("multibyte tag value not stored faithfully: %q", resp.Tags["MultibyteValue"])
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "ListTagsForResource", func() error {
		if err := tc.requireGraphWithARN(); err != nil {
			return err
		}
		resp, err := tc.client.ListTagsForResource(tc.ctx, &neptunegraph.ListTagsForResourceInput{
			ResourceArn: aws.String(tc.graphARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("expected non-nil Tags map")
		}
		if resp.Tags["Environment"] != "test" {
			return fmt.Errorf("expected tag Environment=test, got %v", resp.Tags["Environment"])
		}
		if resp.Tags["Owner"] != "sdk-test" {
			return fmt.Errorf("expected tag Owner=sdk-test, got %v", resp.Tags["Owner"])
		}
		if resp.Tags["ExtraTag"] != "extra-value" {
			return fmt.Errorf("expected tag ExtraTag=extra-value, got %v", resp.Tags["ExtraTag"])
		}
		return nil
	}))

	results = append(results, r.RunTest("neptunegraph", "UntagResource", func() error {
		if err := tc.requireGraphWithARN(); err != nil {
			return err
		}
		_, err := tc.client.UntagResource(tc.ctx, &neptunegraph.UntagResourceInput{
			ResourceArn: aws.String(tc.graphARN),
			TagKeys:     []string{"ExtraTag"},
		})
		return err
	}))

	results = append(results, r.RunTest("neptunegraph", "ListTagsForResource_AfterUntag", func() error {
		if err := tc.requireGraphWithARN(); err != nil {
			return err
		}
		resp, err := tc.client.ListTagsForResource(tc.ctx, &neptunegraph.ListTagsForResourceInput{
			ResourceArn: aws.String(tc.graphARN),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("expected non-nil Tags map")
		}
		if _, ok := resp.Tags["ExtraTag"]; ok {
			return fmt.Errorf("expected ExtraTag to be removed")
		}
		if resp.Tags["Environment"] != "test" {
			return fmt.Errorf("expected tag Environment=test still present, got %v", resp.Tags["Environment"])
		}
		return nil
	}))

	return results
}
