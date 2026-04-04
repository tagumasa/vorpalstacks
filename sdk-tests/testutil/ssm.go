package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunSSMTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "ssm",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := ssm.NewFromConfig(cfg)
	ctx := context.Background()

	ts := fmt.Sprintf("%d", time.Now().UnixNano())

	results = append(results, r.RunTest("ssm", "PutParameter", func() error {
		resp, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String("/test/param-" + ts),
			Value: aws.String("test-value"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter", func() error {
		resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String("/test/param-" + ts),
		})
		if err != nil {
			return err
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameters", func() error {
		resp, err := client.GetParameters(ctx, &ssm.GetParametersInput{
			Names: []string{"/test/param-" + ts},
		})
		if err != nil {
			return err
		}
		if resp.Parameters == nil {
			return fmt.Errorf("parameters list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParametersByPath", func() error {
		resp, err := client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:      aws.String("/test"),
			Recursive: aws.Bool(true),
		})
		if err != nil {
			return err
		}
		if resp.Parameters == nil {
			return fmt.Errorf("parameters list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters", func() error {
		resp, err := client.DescribeParameters(ctx, &ssm.DescribeParametersInput{
			Filters: []types.ParametersFilter{
				{
					Key:    types.ParametersFilterKeyName,
					Values: []string{"/test/param-" + ts},
				},
			},
		})
		if err != nil {
			return err
		}
		if resp.Parameters == nil {
			return fmt.Errorf("parameters list is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameter", func() error {
		resp, err := client.DeleteParameter(ctx, &ssm.DeleteParameterInput{
			Name: aws.String("/test/param-" + ts),
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return fmt.Errorf("response is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter_NonExistent", func() error {
		_, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String("/nonexistent/param-xyz"),
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameter_NonExistent", func() error {
		_, err := client.DeleteParameter(ctx, &ssm.DeleteParameterInput{
			Name: aws.String("/nonexistent/param-xyz"),
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_GetParameter_Roundtrip", func() error {
		rtName := fmt.Sprintf("/rt/param-%d", time.Now().UnixNano())
		rtValue := "roundtrip-value-12345"
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(rtName),
			Value: aws.String(rtValue),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(rtName)})

		resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String(rtName),
		})
		if err != nil {
			return fmt.Errorf("get: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != rtValue {
			return fmt.Errorf("value mismatch: got %q, want %q", aws.ToString(resp.Parameter.Value), rtValue)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameters_InvalidNames", func() error {
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String("/valid/param-test"),
			Value: aws.String("valid"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String("/valid/param-test")})

		resp, err := client.GetParameters(ctx, &ssm.GetParametersInput{
			Names: []string{"/valid/param-test", "/nonexistent/param-xyz"},
		})
		if err != nil {
			return fmt.Errorf("get params: %v", err)
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 valid parameter, got %d", len(resp.Parameters))
		}
		if len(resp.InvalidParameters) != 1 {
			return fmt.Errorf("expected 1 invalid parameter, got %d", len(resp.InvalidParameters))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters_ContainsCreated", func() error {
		dpName := fmt.Sprintf("/dp/param-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:        aws.String(dpName),
			Value:       aws.String("desc-test"),
			Type:        types.ParameterTypeString,
			Description: aws.String("Test description for search"),
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(dpName)})

		resp, err := client.DescribeParameters(ctx, &ssm.DescribeParametersInput{
			Filters: []types.ParametersFilter{
				{
					Key:    types.ParametersFilterKeyName,
					Values: []string{dpName},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 parameter, got %d", len(resp.Parameters))
		}
		if resp.Parameters[0].Description == nil || *resp.Parameters[0].Description != "Test description for search" {
			return fmt.Errorf("description mismatch")
		}
		return nil
	}))

	// === TAG OPERATIONS ===

	results = append(results, r.RunTest("ssm", "AddTagsToResource_ListTagsForResource", func() error {
		name := fmt.Sprintf("/tag/param-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("tagged"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.AddTagsToResource(ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("env"), Value: aws.String("test")},
				{Key: aws.String("team"), Value: aws.String("platform")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		resp, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(resp.TagList) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(resp.TagList))
		}
		tagMap := make(map[string]string)
		for _, t := range resp.TagList {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if tagMap["env"] != "test" || tagMap["team"] != "platform" {
			return fmt.Errorf("tag values mismatch: %v", tagMap)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "RemoveTagsFromResource", func() error {
		name := fmt.Sprintf("/tag/rm/param-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("tagged"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.AddTagsToResource(ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("keep"), Value: aws.String("yes")},
				{Key: aws.String("remove"), Value: aws.String("me")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags: %v", err)
		}

		_, err = client.RemoveTagsFromResource(ctx, &ssm.RemoveTagsFromResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			TagKeys:      []string{"remove"},
		})
		if err != nil {
			return fmt.Errorf("remove tags: %v", err)
		}

		resp, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(resp.TagList) != 1 {
			return fmt.Errorf("expected 1 tag after removal, got %d", len(resp.TagList))
		}
		if aws.ToString(resp.TagList[0].Key) != "keep" {
			return fmt.Errorf("expected 'keep' tag, got %q", aws.ToString(resp.TagList[0].Key))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "AddTagsToResource_Merge", func() error {
		name := fmt.Sprintf("/tag/merge/param-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("tagged"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.AddTagsToResource(ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("a"), Value: aws.String("1")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags 1: %v", err)
		}

		_, err = client.AddTagsToResource(ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
			Tags: []types.Tag{
				{Key: aws.String("b"), Value: aws.String("2")},
				{Key: aws.String("a"), Value: aws.String("updated")},
			},
		})
		if err != nil {
			return fmt.Errorf("add tags 2: %v", err)
		}

		resp, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		tagMap := make(map[string]string)
		for _, t := range resp.TagList {
			tagMap[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
		if len(tagMap) != 2 {
			return fmt.Errorf("expected 2 tags, got %d", len(tagMap))
		}
		if tagMap["a"] != "updated" {
			return fmt.Errorf("expected 'a'='updated', got %q", tagMap["a"])
		}
		if tagMap["b"] != "2" {
			return fmt.Errorf("expected 'b'='2', got %q", tagMap["b"])
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "ListTagsForResource_Empty", func() error {
		name := fmt.Sprintf("/tag/empty/param-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("notags"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		resp, err := client.ListTagsForResource(ctx, &ssm.ListTagsForResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("list tags: %v", err)
		}
		if len(resp.TagList) != 0 {
			return fmt.Errorf("expected 0 tags, got %d", len(resp.TagList))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "AddTagsToResource_NonExistent", func() error {
		_, err := client.AddTagsToResource(ctx, &ssm.AddTagsToResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String("/nonexistent/param-xyz"),
			Tags: []types.Tag{
				{Key: aws.String("k"), Value: aws.String("v")},
			},
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "RemoveTagsFromResource_NonExistent", func() error {
		_, err := client.RemoveTagsFromResource(ctx, &ssm.RemoveTagsFromResourceInput{
			ResourceType: types.ResourceTypeForTaggingParameter,
			ResourceId:   aws.String("/nonexistent/param-xyz"),
			TagKeys:      []string{"k"},
		})
		if err := AssertErrorContains(err, "ParameterNotFound"); err != nil {
			return err
		}
		return nil
	}))

	// === VERSION OPERATIONS ===

	results = append(results, r.RunTest("ssm", "PutParameter_Overwrite_IncrementsVersion", func() error {
		name := fmt.Sprintf("/ver/overwrite-%d", time.Now().UnixNano())
		resp1, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("v1"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		if resp1.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp1.Version)
		}

		resp2, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String("v2"),
			Type:      types.ParameterTypeString,
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}
		if resp2.Version != 2 {
			return fmt.Errorf("expected version 2, got %d", resp2.Version)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_Duplicate_NoOverwrite", func() error {
		name := fmt.Sprintf("/ver/dup-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("first"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put first: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("second"),
			Type:  types.ParameterTypeString,
		})
		if err == nil {
			return fmt.Errorf("expected error when creating duplicate without Overwrite")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameterHistory_TwoVersions", func() error {
		name := fmt.Sprintf("/ver/hist-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("v1"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String("v2"),
			Type:      types.ParameterTypeString,
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		resp, err := client.GetParameterHistory(ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		if len(resp.Parameters) != 2 {
			return fmt.Errorf("expected 2 history entries, got %d", len(resp.Parameters))
		}
		if resp.Parameters[0].Version != 2 {
			return fmt.Errorf("expected first entry version 2 (newest first), got %d", resp.Parameters[0].Version)
		}
		if resp.Parameters[1].Version != 1 {
			return fmt.Errorf("expected second entry version 1, got %d", resp.Parameters[1].Version)
		}
		if resp.Parameters[0].Value == nil || *resp.Parameters[0].Value != "v2" {
			return fmt.Errorf("expected v2 value, got %q", aws.ToString(resp.Parameters[0].Value))
		}
		if resp.Parameters[1].Value == nil || *resp.Parameters[1].Value != "v1" {
			return fmt.Errorf("expected v1 value, got %q", aws.ToString(resp.Parameters[1].Value))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameterHistory_ContainsLabels", func() error {
		name := fmt.Sprintf("/ver/histlabel-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("v1"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.LabelParameterVersion(ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"golden"},
		})
		if err != nil {
			return fmt.Errorf("label: %v", err)
		}

		resp, err := client.GetParameterHistory(ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		if len(resp.Parameters) != 1 {
			return fmt.Errorf("expected 1 history entry, got %d", len(resp.Parameters))
		}
		if len(resp.Parameters[0].Labels) != 1 || resp.Parameters[0].Labels[0] != "golden" {
			return fmt.Errorf("expected label 'golden', got %v", resp.Parameters[0].Labels)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "LabelParameterVersion_GetByLabel", func() error {
		name := fmt.Sprintf("/ver/getlabel-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("original"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.LabelParameterVersion(ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"mylabel"},
		})
		if err != nil {
			return fmt.Errorf("label: %v", err)
		}

		selector := name + ":mylabel"
		resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String(selector),
		})
		if err != nil {
			return fmt.Errorf("get by label: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != "original" {
			return fmt.Errorf("value mismatch: got %q, want 'original'", aws.ToString(resp.Parameter.Value))
		}
		if resp.Parameter.Selector == nil || *resp.Parameter.Selector != selector {
			return fmt.Errorf("selector mismatch: got %q, want %q", aws.ToString(resp.Parameter.Selector), selector)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "LabelParameterVersion_MovesLabel", func() error {
		name := fmt.Sprintf("/ver/movelabel-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("v1"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put v1: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String("v2"),
			Type:      types.ParameterTypeString,
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		_, err = client.LabelParameterVersion(ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(1),
			Labels:           []string{"latest"},
		})
		if err != nil {
			return fmt.Errorf("label v1: %v", err)
		}

		_, err = client.LabelParameterVersion(ctx, &ssm.LabelParameterVersionInput{
			Name:             aws.String(name),
			ParameterVersion: aws.Int64(2),
			Labels:           []string{"latest"},
		})
		if err != nil {
			return fmt.Errorf("label v2: %v", err)
		}

		hist, err := client.GetParameterHistory(ctx, &ssm.GetParameterHistoryInput{
			Name: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("history: %v", err)
		}
		for _, p := range hist.Parameters {
			if p.Version == 1 {
				if contains(p.Labels, "latest") {
					return fmt.Errorf("v1 should not have 'latest' label after move")
				}
			}
			if p.Version == 2 {
				if !contains(p.Labels, "latest") {
					return fmt.Errorf("v2 should have 'latest' label")
				}
			}
		}
		return nil
	}))

	// === BATCH OPERATIONS ===

	results = append(results, r.RunTest("ssm", "DeleteParameters_Success", func() error {
		n1 := fmt.Sprintf("/batch/del1-%d", time.Now().UnixNano())
		n2 := fmt.Sprintf("/batch/del2-%d", time.Now().UnixNano()+1)
		for _, n := range []string{n1, n2} {
			_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
				Name:  aws.String(n),
				Value: aws.String("batch"),
				Type:  types.ParameterTypeString,
			})
			if err != nil {
				return fmt.Errorf("put %s: %v", n, err)
			}
		}

		resp, err := client.DeleteParameters(ctx, &ssm.DeleteParametersInput{
			Names: []string{n1, n2},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		if len(resp.DeletedParameters) != 2 {
			return fmt.Errorf("expected 2 deleted, got %d", len(resp.DeletedParameters))
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DeleteParameters_MixedValidInvalid", func() error {
		n1 := fmt.Sprintf("/batch/mixed1-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(n1),
			Value: aws.String("batch"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}

		resp, err := client.DeleteParameters(ctx, &ssm.DeleteParametersInput{
			Names: []string{n1, "/nonexistent/batch-xyz"},
		})
		if err != nil {
			return fmt.Errorf("delete: %v", err)
		}
		if len(resp.DeletedParameters) != 1 {
			return fmt.Errorf("expected 1 deleted, got %d", len(resp.DeletedParameters))
		}
		if len(resp.InvalidParameters) != 1 {
			return fmt.Errorf("expected 1 invalid, got %d", len(resp.InvalidParameters))
		}
		return nil
	}))

	// === EDGE CASES ===

	results = append(results, r.RunTest("ssm", "GetParametersByPath_NonRecursive", func() error {
		base := fmt.Sprintf("/nr/param-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(base),
			Value: aws.String("direct"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put direct: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(base)})

		_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(base + "/child"),
			Value: aws.String("child"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put child: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(base + "/child")})

		resp, err := client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:      aws.String("/nr"),
			Recursive: aws.Bool(false),
		})
		if err != nil {
			return fmt.Errorf("get by path: %v", err)
		}
		for _, p := range resp.Parameters {
			if strings.Contains(aws.ToString(p.Name), "/child") {
				return fmt.Errorf("non-recursive should not return children, got %q", aws.ToString(p.Name))
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "GetParameter_WithVersionSelector", func() error {
		name := fmt.Sprintf("/ver/sel-%d", time.Now().UnixNano())
		_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("version-one"),
			Type:  types.ParameterTypeString,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		_, err = client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:      aws.String(name),
			Value:     aws.String("version-two"),
			Type:      types.ParameterTypeString,
			Overwrite: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("put v2: %v", err)
		}

		resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{
			Name: aws.String(name + ":1"),
		})
		if err != nil {
			return fmt.Errorf("get v1: %v", err)
		}
		if resp.Parameter == nil {
			return fmt.Errorf("parameter is nil")
		}
		if resp.Parameter.Value == nil || *resp.Parameter.Value != "version-one" {
			return fmt.Errorf("expected 'version-one', got %q", aws.ToString(resp.Parameter.Value))
		}
		if resp.Parameter.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Parameter.Version)
		}
		if resp.Parameter.ARN == nil || *resp.Parameter.ARN == "" {
			return fmt.Errorf("ARN should not be empty for version selector")
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters_TypeFilter", func() error {
		strName := fmt.Sprintf("/tf/string-%d", time.Now().UnixNano())
		slName := fmt.Sprintf("/tf/stringlist-%d", time.Now().UnixNano()+1)
		for _, input := range []*ssm.PutParameterInput{
			{Name: aws.String(strName), Value: aws.String("x"), Type: types.ParameterTypeString},
			{Name: aws.String(slName), Value: aws.String("a,b,c"), Type: types.ParameterTypeStringList},
		} {
			_, err := client.PutParameter(ctx, input)
			if err != nil {
				return fmt.Errorf("put: %v", err)
			}
			defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: input.Name})
		}

		resp, err := client.DescribeParameters(ctx, &ssm.DescribeParametersInput{
			Filters: []types.ParametersFilter{
				{
					Key:    types.ParametersFilterKeyType,
					Values: []string{"String"},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		for _, p := range resp.Parameters {
			if p.Type != types.ParameterTypeString {
				return fmt.Errorf("type filter returned non-String type: %v", p.Type)
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "PutParameter_ReturnsVersionAndTier", func() error {
		name := fmt.Sprintf("/resp/param-%d", time.Now().UnixNano())
		resp, err := client.PutParameter(ctx, &ssm.PutParameterInput{
			Name:  aws.String(name),
			Value: aws.String("check"),
			Type:  types.ParameterTypeString,
			Tier:  types.ParameterTierStandard,
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		defer client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(name)})

		if resp.Version != 1 {
			return fmt.Errorf("expected version 1, got %d", resp.Version)
		}
		if resp.Tier != types.ParameterTierStandard {
			return fmt.Errorf("expected Standard tier, got %v", resp.Tier)
		}
		return nil
	}))

	results = append(results, r.RunTest("ssm", "DescribeParameters_Pagination", func() error {
		pgTs := fmt.Sprintf("%d", time.Now().UnixNano())
		var pgParams []string
		for i := 0; i < 5; i++ {
			name := fmt.Sprintf("/pag/param-%s-%d", pgTs, i)
			_, err := client.PutParameter(ctx, &ssm.PutParameterInput{
				Name:  aws.String(name),
				Value: aws.String("pagval"),
				Type:  types.ParameterTypeString,
			})
			if err != nil {
				for _, pn := range pgParams {
					client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(pn)})
				}
				return fmt.Errorf("put parameter %s: %v", name, err)
			}
			pgParams = append(pgParams, name)
		}

		var allParams []string
		var nextToken *string
		for {
			resp, err := client.DescribeParameters(ctx, &ssm.DescribeParametersInput{
				MaxResults: aws.Int32(2),
				NextToken:  nextToken,
			})
			if err != nil {
				for _, pn := range pgParams {
					client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(pn)})
				}
				return fmt.Errorf("describe parameters page: %v", err)
			}
			for _, p := range resp.Parameters {
				if p.Name != nil && strings.Contains(*p.Name, "/pag/param-"+pgTs) {
					allParams = append(allParams, *p.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = resp.NextToken
			} else {
				break
			}
		}

		for _, pn := range pgParams {
			client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(pn)})
		}
		if len(allParams) != 5 {
			return fmt.Errorf("expected 5 paginated parameters, got %d", len(allParams))
		}
		return nil
	}))

	return results
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
