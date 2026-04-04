package testutil

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"

	"vorpalstacks-sdk-tests/config"
)

func (r *TestRunner) RunAppSyncTests() []TestResult {
	var results []TestResult

	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return append(results, TestResult{
			Service:  "appsync",
			TestName: "Setup",
			Status:   "FAIL",
			Error:    fmt.Sprintf("Failed to load config: %v", err),
		})
	}

	client := appsync.NewFromConfig(cfg)
	ctx := context.Background()
	uid := time.Now().UnixNano()

	minEventConfig := &types.EventConfig{
		AuthProviders: []types.AuthProvider{
			{AuthType: types.AuthenticationTypeApiKey},
		},
		ConnectionAuthModes: []types.AuthMode{
			{AuthType: types.AuthenticationTypeApiKey},
		},
		DefaultPublishAuthModes: []types.AuthMode{
			{AuthType: types.AuthenticationTypeApiKey},
		},
		DefaultSubscribeAuthModes: []types.AuthMode{
			{AuthType: types.AuthenticationTypeApiKey},
		},
	}

	// ================================================================
	// Event API (v2) CRUD
	// ================================================================
	var apiId string
	var tagsApiId string
	var ownerApiId string

	results = append(results, r.RunTest("appsync", "CreateApi", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:        aws.String(fmt.Sprintf("test-api-%d", uid)),
			EventConfig: minEventConfig,
		})
		if err != nil {
			return err
		}
		if resp.Api == nil {
			return fmt.Errorf("Api is nil")
		}
		if resp.Api.ApiId == nil || *resp.Api.ApiId == "" {
			return fmt.Errorf("ApiId is empty")
		}
		apiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateApi_WithTags", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:        aws.String(fmt.Sprintf("test-api-tags-%d", uid)),
			EventConfig: minEventConfig,
			Tags: map[string]string{
				"env":  "test",
				"team": "platform",
			},
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || resp.Api.ApiId == nil {
			return fmt.Errorf("invalid response")
		}
		if len(resp.Api.Tags) != 2 || resp.Api.Tags["env"] != "test" || resp.Api.Tags["team"] != "platform" {
			return fmt.Errorf("tags not persisted: %v", resp.Api.Tags)
		}
		tagsApiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateApi_WithOwnerContact", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:         aws.String(fmt.Sprintf("test-api-owner-%d", uid)),
			EventConfig:  minEventConfig,
			OwnerContact: aws.String("test@example.com"),
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || resp.Api.OwnerContact == nil || *resp.Api.OwnerContact != "test@example.com" {
			return fmt.Errorf("ownerContact not set: %v", resp.Api.OwnerContact)
		}
		ownerApiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApi", func() error {
		resp, err := client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(apiId)})
		if err != nil {
			return err
		}
		if resp.Api == nil {
			return fmt.Errorf("Api is nil")
		}
		if *resp.Api.ApiId != apiId {
			return fmt.Errorf("expected ApiId %s, got %s", apiId, *resp.Api.ApiId)
		}
		if resp.Api.EventConfig == nil {
			return fmt.Errorf("EventConfig is nil")
		}
		if len(resp.Api.EventConfig.AuthProviders) != 1 {
			return fmt.Errorf("expected 1 auth provider, got %d", len(resp.Api.EventConfig.AuthProviders))
		}
		if resp.Api.ApiArn == nil {
			return fmt.Errorf("ApiArn is nil")
		}
		if !strings.HasPrefix(*resp.Api.ApiArn, "arn:aws:appsync:") {
			return fmt.Errorf("invalid ARN format: %s", *resp.Api.ApiArn)
		}
		if resp.Api.Created == nil {
			return fmt.Errorf("Created timestamp is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApi_NonExistent", func() error {
		_, err := client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String("does-not-exist")})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListApis", func() error {
		resp, err := client.ListApis(ctx, &appsync.ListApisInput{})
		if err != nil {
			return err
		}
		if len(resp.Apis) < 3 {
			return fmt.Errorf("expected at least 3 APIs, got %d", len(resp.Apis))
		}
		found := false
		for _, api := range resp.Apis {
			if *api.ApiId == apiId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created API not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListApis_WithPagination", func() error {
		resp, err := client.ListApis(ctx, &appsync.ListApisInput{MaxResults: 1})
		if err != nil {
			return err
		}
		if len(resp.Apis) != 1 {
			return fmt.Errorf("expected 1 API with maxResults=1, got %d", len(resp.Apis))
		}
		if resp.NextToken == nil {
			return fmt.Errorf("expected NextToken when more results exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApi", func() error {
		newName := fmt.Sprintf("updated-api-%d", uid)
		resp, err := client.UpdateApi(ctx, &appsync.UpdateApiInput{
			ApiId:       aws.String(apiId),
			Name:        aws.String(newName),
			EventConfig: minEventConfig,
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || *resp.Api.Name != newName {
			return fmt.Errorf("expected name %s, got %v", newName, resp.Api.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApi_NonExistent", func() error {
		_, err := client.UpdateApi(ctx, &appsync.UpdateApiInput{
			ApiId:       aws.String("does-not-exist"),
			Name:        aws.String("noop"),
			EventConfig: minEventConfig,
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// ================================================================
	// Channel Namespace CRUD
	// ================================================================
	var nsName = fmt.Sprintf("testns-%d", uid)

	results = append(results, r.RunTest("appsync", "CreateChannelNamespace", func() error {
		resp, err := client.CreateChannelNamespace(ctx, &appsync.CreateChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(nsName),
			Tags: map[string]string{
				"type": "test",
			},
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil {
			return fmt.Errorf("ChannelNamespace is nil")
		}
		if *resp.ChannelNamespace.Name != nsName {
			return fmt.Errorf("expected name %s, got %s", nsName, *resp.ChannelNamespace.Name)
		}
		if resp.ChannelNamespace.ApiId == nil || *resp.ChannelNamespace.ApiId != apiId {
			return fmt.Errorf("ApiId mismatch: %v", resp.ChannelNamespace.ApiId)
		}
		if resp.ChannelNamespace.ChannelNamespaceArn == nil {
			return fmt.Errorf("ChannelNamespaceArn is nil")
		}
		if resp.ChannelNamespace.Created == nil {
			return fmt.Errorf("Created timestamp is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetChannelNamespace", func() error {
		resp, err := client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(nsName),
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil || *resp.ChannelNamespace.Name != nsName {
			return fmt.Errorf("expected namespace %s, got %v", nsName, resp.ChannelNamespace)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetChannelNamespace_NonExistent", func() error {
		_, err := client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String("does-not-exist"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListChannelNamespaces", func() error {
		resp, err := client.ListChannelNamespaces(ctx, &appsync.ListChannelNamespacesInput{
			ApiId: aws.String(apiId),
		})
		if err != nil {
			return err
		}
		if len(resp.ChannelNamespaces) < 1 {
			return fmt.Errorf("expected at least 1 namespace, got %d", len(resp.ChannelNamespaces))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateChannelNamespace", func() error {
		resp, err := client.UpdateChannelNamespace(ctx, &appsync.UpdateChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(nsName),
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil || *resp.ChannelNamespace.Name != nsName {
			return fmt.Errorf("expected namespace %s, got %v", nsName, resp.ChannelNamespace)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateChannelNamespace_NonExistent", func() error {
		_, err := client.UpdateChannelNamespace(ctx, &appsync.UpdateChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String("does-not-exist"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// ================================================================
	// Tag operations
	// ================================================================
	taggedApiArn := ""
	var taggedApiId string

	results = append(results, r.RunTest("appsync", "CreateApi_ForTagging", func() error {
		resp, err := client.CreateApi(ctx, &appsync.CreateApiInput{
			Name:        aws.String(fmt.Sprintf("test-api-tag-%d", uid)),
			EventConfig: minEventConfig,
			Tags: map[string]string{
				"key1": "value1",
			},
		})
		if err != nil {
			return err
		}
		if resp.Api == nil || resp.Api.ApiArn == nil {
			return fmt.Errorf("invalid response")
		}
		taggedApiArn = *resp.Api.ApiArn
		taggedApiId = *resp.Api.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListTagsForResource", func() error {
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(taggedApiArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags == nil {
			return fmt.Errorf("Tags is nil")
		}
		if resp.Tags["key1"] != "value1" {
			return fmt.Errorf("expected key1=value1, got: %v", resp.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "TagResource", func() error {
		_, err := client.TagResource(ctx, &appsync.TagResourceInput{
			ResourceArn: aws.String(taggedApiArn),
			Tags: map[string]string{
				"key2": "value2",
				"key3": "value3",
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(taggedApiArn),
		})
		if err != nil {
			return err
		}
		if len(resp.Tags) != 3 {
			return fmt.Errorf("expected 3 tags, got %d", len(resp.Tags))
		}
		if resp.Tags["key2"] != "value2" || resp.Tags["key3"] != "value3" {
			return fmt.Errorf("new tags not found: %v", resp.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UntagResource", func() error {
		_, err := client.UntagResource(ctx, &appsync.UntagResourceInput{
			ResourceArn: aws.String(taggedApiArn),
			TagKeys:     []string{"key2"},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(taggedApiArn),
		})
		if err != nil {
			return err
		}
		if _, exists := resp.Tags["key2"]; exists {
			return fmt.Errorf("key2 should have been removed: %v", resp.Tags)
		}
		if resp.Tags["key1"] != "value1" || resp.Tags["key3"] != "value3" {
			return fmt.Errorf("remaining tags incorrect: %v", resp.Tags)
		}
		return nil
	}))

	// Tag operations on ChannelNamespace
	var nsArn string
	var tagNsName string

	results = append(results, r.RunTest("appsync", "CreateChannelNamespace_ForTagging", func() error {
		tagNsName = fmt.Sprintf("tag-ns-%d", uid)
		resp, err := client.CreateChannelNamespace(ctx, &appsync.CreateChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(tagNsName),
			Tags: map[string]string{
				"nsKey": "nsValue",
			},
		})
		if err != nil {
			return err
		}
		if resp.ChannelNamespace == nil || resp.ChannelNamespace.ChannelNamespaceArn == nil {
			return fmt.Errorf("invalid response")
		}
		nsArn = *resp.ChannelNamespace.ChannelNamespaceArn
		return nil
	}))

	results = append(results, r.RunTest("appsync", "TagResource_ChannelNamespace", func() error {
		_, err := client.TagResource(ctx, &appsync.TagResourceInput{
			ResourceArn: aws.String(nsArn),
			Tags: map[string]string{
				"added": "yes",
			},
		})
		if err != nil {
			return err
		}
		resp, err := client.ListTagsForResource(ctx, &appsync.ListTagsForResourceInput{
			ResourceArn: aws.String(nsArn),
		})
		if err != nil {
			return err
		}
		if resp.Tags["nsKey"] != "nsValue" || resp.Tags["added"] != "yes" {
			return fmt.Errorf("namespace tags incorrect: %v", resp.Tags)
		}
		return nil
	}))

	// ================================================================
	// GraphQL API (v1) CRUD
	// ================================================================
	var gqlApiId string
	var gqlTagsApiId string

	results = append(results, r.RunTest("appsync", "CreateGraphqlApi", func() error {
		resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
			Name:               aws.String(fmt.Sprintf("test-gql-%d", uid)),
			AuthenticationType: types.AuthenticationTypeApiKey,
		})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil {
			return fmt.Errorf("GraphqlApi is nil")
		}
		if resp.GraphqlApi.ApiId == nil || *resp.GraphqlApi.ApiId == "" {
			return fmt.Errorf("ApiId is empty")
		}
		if resp.GraphqlApi.AuthenticationType != types.AuthenticationTypeApiKey {
			return fmt.Errorf("expected API_KEY auth type, got %s", resp.GraphqlApi.AuthenticationType)
		}
		if resp.GraphqlApi.Arn == nil {
			return fmt.Errorf("Arn is nil")
		}
		if !strings.HasPrefix(*resp.GraphqlApi.Arn, "arn:aws:appsync:") {
			return fmt.Errorf("invalid ARN format: %s", *resp.GraphqlApi.Arn)
		}
		gqlApiId = *resp.GraphqlApi.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateGraphqlApi_WithTags", func() error {
		resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
			Name:               aws.String(fmt.Sprintf("test-gql-tags-%d", uid)),
			AuthenticationType: types.AuthenticationTypeApiKey,
			Tags:               map[string]string{"env": "dev"},
		})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil || resp.GraphqlApi.Tags == nil {
			return fmt.Errorf("Tags not returned")
		}
		if resp.GraphqlApi.Tags["env"] != "dev" {
			return fmt.Errorf("tag not persisted: %v", resp.GraphqlApi.Tags)
		}
		gqlTagsApiId = *resp.GraphqlApi.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetGraphqlApi", func() error {
		resp, err := client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(gqlApiId)})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil {
			return fmt.Errorf("GraphqlApi is nil")
		}
		if *resp.GraphqlApi.ApiId != gqlApiId {
			return fmt.Errorf("expected ApiId %s, got %s", gqlApiId, *resp.GraphqlApi.ApiId)
		}
		if resp.GraphqlApi.AuthenticationType != types.AuthenticationTypeApiKey {
			return fmt.Errorf("expected API_KEY auth type, got %s", resp.GraphqlApi.AuthenticationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetGraphqlApi_NonExistent", func() error {
		_, err := client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String("does-not-exist")})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListGraphqlApis", func() error {
		resp, err := client.ListGraphqlApis(ctx, &appsync.ListGraphqlApisInput{})
		if err != nil {
			return err
		}
		if len(resp.GraphqlApis) < 2 {
			return fmt.Errorf("expected at least 2 GraphQL APIs, got %d", len(resp.GraphqlApis))
		}
		found := false
		for _, api := range resp.GraphqlApis {
			if *api.ApiId == gqlApiId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("created GraphQL API not found in list")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListGraphqlApis_WithPagination", func() error {
		resp, err := client.ListGraphqlApis(ctx, &appsync.ListGraphqlApisInput{MaxResults: 1})
		if err != nil {
			return err
		}
		if len(resp.GraphqlApis) != 1 {
			return fmt.Errorf("expected 1 API with maxResults=1, got %d", len(resp.GraphqlApis))
		}
		if resp.NextToken == nil {
			return fmt.Errorf("expected NextToken when more results exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateGraphqlApi", func() error {
		newName := fmt.Sprintf("updated-gql-%d", uid)
		resp, err := client.UpdateGraphqlApi(ctx, &appsync.UpdateGraphqlApiInput{
			ApiId:              aws.String(gqlApiId),
			Name:               aws.String(newName),
			AuthenticationType: types.AuthenticationTypeOpenidConnect,
		})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil || *resp.GraphqlApi.Name != newName {
			return fmt.Errorf("expected name %s, got %v", newName, resp.GraphqlApi.Name)
		}
		if resp.GraphqlApi.AuthenticationType != types.AuthenticationTypeOpenidConnect {
			return fmt.Errorf("expected OPENID_CONNECT auth type, got %s", resp.GraphqlApi.AuthenticationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateGraphqlApi_NonExistent", func() error {
		_, err := client.UpdateGraphqlApi(ctx, &appsync.UpdateGraphqlApiInput{
			ApiId:              aws.String("does-not-exist"),
			Name:               aws.String("noop"),
			AuthenticationType: types.AuthenticationTypeApiKey,
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent GraphQL API")
		}
		return nil
	}))

	// ================================================================
	// DataSource CRUD
	// ================================================================
	results = append(results, r.RunTest("appsync", "CreateDataSource", func() error {
		resp, err := client.CreateDataSource(ctx, &appsync.CreateDataSourceInput{
			ApiId: aws.String(gqlApiId),
			Name:  aws.String("testDS"),
			Type:  types.DataSourceTypeNone,
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil {
			return fmt.Errorf("DataSource is nil")
		}
		if *resp.DataSource.Name != "testDS" {
			return fmt.Errorf("expected name testDS, got %s", *resp.DataSource.Name)
		}
		if resp.DataSource.Type != types.DataSourceTypeNone {
			return fmt.Errorf("expected NONE type, got %s", resp.DataSource.Type)
		}
		if resp.DataSource.DataSourceArn == nil {
			return fmt.Errorf("DataSourceArn is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateDataSource_WithDescription", func() error {
		resp, err := client.CreateDataSource(ctx, &appsync.CreateDataSourceInput{
			ApiId:       aws.String(gqlApiId),
			Name:        aws.String("testDS2"),
			Type:        types.DataSourceTypeNone,
			Description: aws.String("test description"),
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil {
			return fmt.Errorf("DataSource is nil")
		}
		if resp.DataSource.Description == nil || *resp.DataSource.Description != "test description" {
			return fmt.Errorf("description not set: %v", resp.DataSource.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDataSource", func() error {
		resp, err := client.GetDataSource(ctx, &appsync.GetDataSourceInput{
			ApiId: aws.String(gqlApiId),
			Name:  aws.String("testDS"),
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil || *resp.DataSource.Name != "testDS" {
			return fmt.Errorf("expected testDS, got %v", resp.DataSource)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDataSource_NonExistent", func() error {
		_, err := client.GetDataSource(ctx, &appsync.GetDataSourceInput{
			ApiId: aws.String(gqlApiId),
			Name:  aws.String("does-not-exist"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListDataSources", func() error {
		resp, err := client.ListDataSources(ctx, &appsync.ListDataSourcesInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if len(resp.DataSources) < 2 {
			return fmt.Errorf("expected at least 2 data sources, got %d", len(resp.DataSources))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateDataSource", func() error {
		resp, err := client.UpdateDataSource(ctx, &appsync.UpdateDataSourceInput{
			ApiId:       aws.String(gqlApiId),
			Name:        aws.String("testDS"),
			Type:        types.DataSourceTypeNone,
			Description: aws.String("updated description"),
		})
		if err != nil {
			return err
		}
		if resp.DataSource == nil {
			return fmt.Errorf("DataSource is nil")
		}
		if resp.DataSource.Description == nil || *resp.DataSource.Description != "updated description" {
			return fmt.Errorf("description not updated: %v", resp.DataSource.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteDataSource", func() error {
		_, err := client.DeleteDataSource(ctx, &appsync.DeleteDataSourceInput{
			ApiId: aws.String(gqlApiId),
			Name:  aws.String("testDS2"),
		})
		if err != nil {
			return err
		}
		_, err = client.GetDataSource(ctx, &appsync.GetDataSourceInput{
			ApiId: aws.String(gqlApiId),
			Name:  aws.String("testDS2"),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting data source")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteDataSource_NonExistent", func() error {
		_, err := client.DeleteDataSource(ctx, &appsync.DeleteDataSourceInput{
			ApiId: aws.String(gqlApiId),
			Name:  aws.String("already-deleted"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent data source")
		}
		return nil
	}))

	// ================================================================
	// Type CRUD
	// ================================================================
	sdl := `type Post { id: ID! title: String! }`

	results = append(results, r.RunTest("appsync", "CreateType", func() error {
		resp, err := client.CreateType(ctx, &appsync.CreateTypeInput{
			ApiId:      aws.String(gqlApiId),
			Definition: aws.String(sdl),
			Format:     types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if resp.Type == nil {
			return fmt.Errorf("Type is nil")
		}
		if resp.Type.Name == nil || *resp.Type.Name != "Post" {
			return fmt.Errorf("expected name Post, got %v", resp.Type.Name)
		}
		if resp.Type.Format != types.TypeDefinitionFormatSdl {
			return fmt.Errorf("expected SDL format, got %s", resp.Type.Format)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetType", func() error {
		resp, err := client.GetType(ctx, &appsync.GetTypeInput{
			ApiId:    aws.String(gqlApiId),
			TypeName: aws.String("Post"),
			Format:   types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if resp.Type == nil || *resp.Type.Name != "Post" {
			return fmt.Errorf("expected Post type, got %v", resp.Type)
		}
		if resp.Type.Definition == nil || *resp.Type.Definition == "" {
			return fmt.Errorf("definition is empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetType_NonExistent", func() error {
		_, err := client.GetType(ctx, &appsync.GetTypeInput{
			ApiId:    aws.String(gqlApiId),
			TypeName: aws.String("DoesNotExist"),
			Format:   types.TypeDefinitionFormatSdl,
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListTypes", func() error {
		resp, err := client.ListTypes(ctx, &appsync.ListTypesInput{
			ApiId:  aws.String(gqlApiId),
			Format: types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if len(resp.Types) < 1 {
			return fmt.Errorf("expected at least 1 type, got %d", len(resp.Types))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateType", func() error {
		updatedSDL := `type Post { id: ID! title: String! content: String }`
		resp, err := client.UpdateType(ctx, &appsync.UpdateTypeInput{
			ApiId:      aws.String(gqlApiId),
			TypeName:   aws.String("Post"),
			Definition: aws.String(updatedSDL),
			Format:     types.TypeDefinitionFormatSdl,
		})
		if err != nil {
			return err
		}
		if resp.Type == nil {
			return fmt.Errorf("Type is nil")
		}
		if resp.Type.Definition == nil || !strings.Contains(*resp.Type.Definition, "content") {
			return fmt.Errorf("definition not updated: %v", resp.Type.Definition)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteType", func() error {
		_, err := client.DeleteType(ctx, &appsync.DeleteTypeInput{
			ApiId:    aws.String(gqlApiId),
			TypeName: aws.String("Post"),
		})
		if err != nil {
			return err
		}
		_, err = client.GetType(ctx, &appsync.GetTypeInput{
			ApiId:    aws.String(gqlApiId),
			TypeName: aws.String("Post"),
			Format:   types.TypeDefinitionFormatSdl,
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting type")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteType_NonExistent", func() error {
		_, err := client.DeleteType(ctx, &appsync.DeleteTypeInput{
			ApiId:    aws.String(gqlApiId),
			TypeName: aws.String("already-deleted"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent type")
		}
		return nil
	}))

	// ================================================================
	// Resolver CRUD
	// ================================================================
	results = append(results, r.RunTest("appsync", "CreateResolver", func() error {
		resp, err := client.CreateResolver(ctx, &appsync.CreateResolverInput{
			ApiId:          aws.String(gqlApiId),
			TypeName:       aws.String("Query"),
			FieldName:      aws.String("getPost"),
			DataSourceName: aws.String("testDS"),
		})
		if err != nil {
			return err
		}
		if resp.Resolver == nil {
			return fmt.Errorf("Resolver is nil")
		}
		if *resp.Resolver.FieldName != "getPost" {
			return fmt.Errorf("expected fieldName getPost, got %s", *resp.Resolver.FieldName)
		}
		if *resp.Resolver.TypeName != "Query" {
			return fmt.Errorf("expected typeName Query, got %s", *resp.Resolver.TypeName)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetResolver", func() error {
		resp, err := client.GetResolver(ctx, &appsync.GetResolverInput{
			ApiId:     aws.String(gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("getPost"),
		})
		if err != nil {
			return err
		}
		if resp.Resolver == nil || *resp.Resolver.FieldName != "getPost" {
			return fmt.Errorf("expected getPost resolver, got %v", resp.Resolver)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetResolver_NonExistent", func() error {
		_, err := client.GetResolver(ctx, &appsync.GetResolverInput{
			ApiId:     aws.String(gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("doesNotExist"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListResolvers", func() error {
		resp, err := client.ListResolvers(ctx, &appsync.ListResolversInput{
			ApiId:    aws.String(gqlApiId),
			TypeName: aws.String("Query"),
		})
		if err != nil {
			return err
		}
		if len(resp.Resolvers) < 1 {
			return fmt.Errorf("expected at least 1 resolver, got %d", len(resp.Resolvers))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateResolver", func() error {
		resp, err := client.UpdateResolver(ctx, &appsync.UpdateResolverInput{
			ApiId:                  aws.String(gqlApiId),
			TypeName:               aws.String("Query"),
			FieldName:              aws.String("getPost"),
			RequestMappingTemplate: aws.String("{\"version\": \"2017-02-28\", \"payload\": {}}"),
		})
		if err != nil {
			return err
		}
		if resp.Resolver == nil {
			return fmt.Errorf("Resolver is nil")
		}
		if resp.Resolver.RequestMappingTemplate == nil || *resp.Resolver.RequestMappingTemplate == "" {
			return fmt.Errorf("requestMappingTemplate not updated")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteResolver", func() error {
		_, err := client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId:     aws.String(gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("getPost"),
		})
		if err != nil {
			return err
		}
		_, err = client.GetResolver(ctx, &appsync.GetResolverInput{
			ApiId:     aws.String(gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("getPost"),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting resolver")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteResolver_NonExistent", func() error {
		_, err := client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId:     aws.String(gqlApiId),
			TypeName:  aws.String("Query"),
			FieldName: aws.String("already-deleted"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent resolver")
		}
		return nil
	}))

	// ================================================================
	// Function (AppSync) CRUD
	// ================================================================
	var functionId string

	results = append(results, r.RunTest("appsync", "CreateFunction", func() error {
		resp, err := client.CreateFunction(ctx, &appsync.CreateFunctionInput{
			ApiId:           aws.String(gqlApiId),
			Name:            aws.String("testFn"),
			DataSourceName:  aws.String("testDS"),
			FunctionVersion: aws.String("2018-05-29"),
		})
		if err != nil {
			return err
		}
		if resp.FunctionConfiguration == nil {
			return fmt.Errorf("FunctionConfiguration is nil")
		}
		if resp.FunctionConfiguration.FunctionId == nil || *resp.FunctionConfiguration.FunctionId == "" {
			return fmt.Errorf("FunctionId is empty")
		}
		if *resp.FunctionConfiguration.Name != "testFn" {
			return fmt.Errorf("expected name testFn, got %s", *resp.FunctionConfiguration.Name)
		}
		functionId = *resp.FunctionConfiguration.FunctionId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetFunction", func() error {
		resp, err := client.GetFunction(ctx, &appsync.GetFunctionInput{
			ApiId:      aws.String(gqlApiId),
			FunctionId: aws.String(functionId),
		})
		if err != nil {
			return err
		}
		if resp.FunctionConfiguration == nil || *resp.FunctionConfiguration.FunctionId != functionId {
			return fmt.Errorf("expected function %s, got %v", functionId, resp.FunctionConfiguration)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetFunction_NonExistent", func() error {
		_, err := client.GetFunction(ctx, &appsync.GetFunctionInput{
			ApiId:      aws.String(gqlApiId),
			FunctionId: aws.String("does-not-exist"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListFunctions", func() error {
		resp, err := client.ListFunctions(ctx, &appsync.ListFunctionsInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if len(resp.Functions) < 1 {
			return fmt.Errorf("expected at least 1 function, got %d", len(resp.Functions))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateFunction", func() error {
		resp, err := client.UpdateFunction(ctx, &appsync.UpdateFunctionInput{
			ApiId:           aws.String(gqlApiId),
			FunctionId:      aws.String(functionId),
			Name:            aws.String("updatedFn"),
			DataSourceName:  aws.String("testDS"),
			FunctionVersion: aws.String("2018-05-29"),
		})
		if err != nil {
			return err
		}
		if resp.FunctionConfiguration == nil || *resp.FunctionConfiguration.Name != "updatedFn" {
			return fmt.Errorf("expected name updatedFn, got %v", resp.FunctionConfiguration.Name)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteFunction", func() error {
		_, err := client.DeleteFunction(ctx, &appsync.DeleteFunctionInput{
			ApiId:      aws.String(gqlApiId),
			FunctionId: aws.String(functionId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetFunction(ctx, &appsync.GetFunctionInput{
			ApiId:      aws.String(gqlApiId),
			FunctionId: aws.String(functionId),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting function")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteFunction_NonExistent", func() error {
		_, err := client.DeleteFunction(ctx, &appsync.DeleteFunctionInput{
			ApiId:      aws.String(gqlApiId),
			FunctionId: aws.String("already-deleted"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// ================================================================
	// ListResolversByFunction
	// ================================================================
	results = append(results, r.RunTest("appsync", "ListResolversByFunction", func() error {
		createFnResp, err := client.CreateFunction(ctx, &appsync.CreateFunctionInput{
			ApiId:           aws.String(gqlApiId),
			Name:            aws.String("listByFn"),
			DataSourceName:  aws.String("testDS"),
			FunctionVersion: aws.String("2018-05-29"),
		})
		if err != nil {
			return err
		}
		listFnId := *createFnResp.FunctionConfiguration.FunctionId

		_, err = client.CreateResolver(ctx, &appsync.CreateResolverInput{
			ApiId:          aws.String(gqlApiId),
			TypeName:       aws.String("Query"),
			FieldName:      aws.String("fnTest"),
			DataSourceName: aws.String("testDS"),
		})
		if err != nil {
			return err
		}

		pipelineResp, err := client.CreateResolver(ctx, &appsync.CreateResolverInput{
			ApiId:          aws.String(gqlApiId),
			TypeName:       aws.String("Query"),
			FieldName:      aws.String("pipelineTest"),
			Kind:           types.ResolverKindPipeline,
			PipelineConfig: &types.PipelineConfig{Functions: []string{listFnId}},
		})
		if err != nil {
			return err
		}

		listResp, err := client.ListResolversByFunction(ctx, &appsync.ListResolversByFunctionInput{
			ApiId:      aws.String(gqlApiId),
			FunctionId: aws.String(listFnId),
		})
		if err != nil {
			return err
		}
		if len(listResp.Resolvers) < 1 {
			return fmt.Errorf("expected at least 1 resolver, got %d", len(listResp.Resolvers))
		}
		found := false
		for _, r := range listResp.Resolvers {
			if *r.FieldName == *pipelineResp.Resolver.FieldName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("pipeline resolver not found in ListResolversByFunction")
		}

		client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId: aws.String(gqlApiId), TypeName: aws.String("Query"), FieldName: aws.String("fnTest"),
		})
		client.DeleteResolver(ctx, &appsync.DeleteResolverInput{
			ApiId: aws.String(gqlApiId), TypeName: aws.String("Query"), FieldName: aws.String("pipelineTest"),
		})
		client.DeleteFunction(ctx, &appsync.DeleteFunctionInput{
			ApiId: aws.String(gqlApiId), FunctionId: aws.String(listFnId),
		})
		return nil
	}))

	// ================================================================
	// Schema operations
	// ================================================================
	results = append(results, r.RunTest("appsync", "StartSchemaCreation", func() error {
		schemaSDL := []byte("type Query { hello: String } type Mutation { addPost(title: String!): Post } type Post { id: ID! title: String! }")
		resp, err := client.StartSchemaCreation(ctx, &appsync.StartSchemaCreationInput{
			ApiId:      aws.String(gqlApiId),
			Definition: schemaSDL,
		})
		if err != nil {
			return err
		}
		if resp.Status != "PROCESSING" && resp.Status != "SUCCESS" {
			return fmt.Errorf("expected PROCESSING or SUCCESS status, got %s", resp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetSchemaCreationStatus", func() error {
		resp, err := client.GetSchemaCreationStatus(ctx, &appsync.GetSchemaCreationStatusInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.Status != "SUCCESS" && resp.Status != "PROCESSING" {
			return fmt.Errorf("expected SUCCESS or PROCESSING status, got %s", resp.Status)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetIntrospectionSchema", func() error {
		resp, err := client.GetIntrospectionSchema(ctx, &appsync.GetIntrospectionSchemaInput{
			ApiId:  aws.String(gqlApiId),
			Format: types.OutputTypeSdl,
		})
		if err != nil {
			return err
		}
		schemaStr := string(resp.Schema)
		if len(schemaStr) == 0 {
			return fmt.Errorf("Schema is empty")
		}
		for _, expected := range []string{"type Query", "type Mutation", "type Post"} {
			if !strings.Contains(schemaStr, expected) {
				return fmt.Errorf("schema missing %q", expected)
			}
		}
		return nil
	}))

	// ================================================================
	// Environment Variables
	// ================================================================
	results = append(results, r.RunTest("appsync", "PutGraphqlApiEnvironmentVariables", func() error {
		resp, err := client.PutGraphqlApiEnvironmentVariables(ctx, &appsync.PutGraphqlApiEnvironmentVariablesInput{
			ApiId: aws.String(gqlApiId),
			EnvironmentVariables: map[string]string{
				"ENV1": "value1",
				"ENV2": "value2",
			},
		})
		if err != nil {
			return err
		}
		if resp.EnvironmentVariables == nil {
			return fmt.Errorf("EnvironmentVariables is nil")
		}
		if resp.EnvironmentVariables["ENV1"] != "value1" || resp.EnvironmentVariables["ENV2"] != "value2" {
			return fmt.Errorf("env vars not persisted: %v", resp.EnvironmentVariables)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetGraphqlApiEnvironmentVariables", func() error {
		resp, err := client.GetGraphqlApiEnvironmentVariables(ctx, &appsync.GetGraphqlApiEnvironmentVariablesInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.EnvironmentVariables == nil {
			return fmt.Errorf("EnvironmentVariables is nil")
		}
		if resp.EnvironmentVariables["ENV1"] != "value1" || resp.EnvironmentVariables["ENV2"] != "value2" {
			return fmt.Errorf("env vars incorrect: %v", resp.EnvironmentVariables)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "PutGraphqlApiEnvironmentVariables_Replace", func() error {
		resp, err := client.PutGraphqlApiEnvironmentVariables(ctx, &appsync.PutGraphqlApiEnvironmentVariablesInput{
			ApiId: aws.String(gqlApiId),
			EnvironmentVariables: map[string]string{
				"ENV3": "value3",
			},
		})
		if err != nil {
			return err
		}
		if len(resp.EnvironmentVariables) != 1 {
			return fmt.Errorf("expected 1 env var after replace, got %d", len(resp.EnvironmentVariables))
		}
		if resp.EnvironmentVariables["ENV3"] != "value3" {
			return fmt.Errorf("expected ENV3=value3, got %v", resp.EnvironmentVariables)
		}
		if _, exists := resp.EnvironmentVariables["ENV1"]; exists {
			return fmt.Errorf("ENV1 should have been replaced")
		}
		return nil
	}))

	// ================================================================
	// API Key CRUD
	// ================================================================
	var apiKeyId string
	var descApiKeyId string

	results = append(results, r.RunTest("appsync", "CreateApiKey", func() error {
		resp, err := client.CreateApiKey(ctx, &appsync.CreateApiKeyInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.ApiKey == nil {
			return fmt.Errorf("ApiKey is nil")
		}
		if resp.ApiKey.Id == nil || *resp.ApiKey.Id == "" {
			return fmt.Errorf("Id is empty")
		}
		if resp.ApiKey.Expires == 0 {
			return fmt.Errorf("Expires should be set (default 365 days)")
		}
		if resp.ApiKey.Deletes == 0 {
			return fmt.Errorf("Deletes should be set (same as Expires)")
		}
		apiKeyId = *resp.ApiKey.Id
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateApiKey_WithDescription", func() error {
		resp, err := client.CreateApiKey(ctx, &appsync.CreateApiKeyInput{
			ApiId:       aws.String(gqlApiId),
			Description: aws.String("test key"),
		})
		if err != nil {
			return err
		}
		if resp.ApiKey == nil {
			return fmt.Errorf("ApiKey is nil")
		}
		if resp.ApiKey.Description == nil || *resp.ApiKey.Description != "test key" {
			return fmt.Errorf("description not set: %v", resp.ApiKey.Description)
		}
		descApiKeyId = *resp.ApiKey.Id
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListApiKeys", func() error {
		resp, err := client.ListApiKeys(ctx, &appsync.ListApiKeysInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if len(resp.ApiKeys) < 2 {
			return fmt.Errorf("expected at least 2 API keys, got %d", len(resp.ApiKeys))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListApiKeys_WithPagination", func() error {
		resp, err := client.ListApiKeys(ctx, &appsync.ListApiKeysInput{
			ApiId:      aws.String(gqlApiId),
			MaxResults: 1,
		})
		if err != nil {
			return err
		}
		if len(resp.ApiKeys) != 1 {
			return fmt.Errorf("expected 1 API key with maxResults=1, got %d", len(resp.ApiKeys))
		}
		if resp.NextToken == nil {
			return fmt.Errorf("expected NextToken when more results exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApiKey", func() error {
		newExpiry := time.Now().Add(180 * 24 * time.Hour).Unix()
		resp, err := client.UpdateApiKey(ctx, &appsync.UpdateApiKeyInput{
			ApiId:       aws.String(gqlApiId),
			Id:          aws.String(apiKeyId),
			Description: aws.String("updated key"),
			Expires:     newExpiry,
		})
		if err != nil {
			return err
		}
		if resp.ApiKey == nil {
			return fmt.Errorf("ApiKey is nil")
		}
		if resp.ApiKey.Description == nil || *resp.ApiKey.Description != "updated key" {
			return fmt.Errorf("description not updated: %v", resp.ApiKey.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApiKey_NonExistent", func() error {
		_, err := client.UpdateApiKey(ctx, &appsync.UpdateApiKeyInput{
			ApiId:       aws.String(gqlApiId),
			Id:          aws.String("does-not-exist"),
			Description: aws.String("noop"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent API key")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiKey", func() error {
		_, err := client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
			ApiId: aws.String(gqlApiId),
			Id:    aws.String(apiKeyId),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiKey_NonExistent", func() error {
		_, err := client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
			ApiId: aws.String(gqlApiId),
			Id:    aws.String("already-deleted"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent API key")
		}
		return nil
	}))

	// ================================================================
	// Cache operations
	// ================================================================
	results = append(results, r.RunTest("appsync", "CreateApiCache", func() error {
		resp, err := client.CreateApiCache(ctx, &appsync.CreateApiCacheInput{
			ApiId:              aws.String(gqlApiId),
			Type:               types.ApiCacheTypeSmall,
			Ttl:                300,
			ApiCachingBehavior: types.ApiCachingBehaviorFullRequestCaching,
		})
		if err != nil {
			return err
		}
		if resp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil")
		}
		if resp.ApiCache.Type != types.ApiCacheTypeSmall {
			return fmt.Errorf("expected SMALL type, got %s", resp.ApiCache.Type)
		}
		if resp.ApiCache.Ttl != 300 {
			return fmt.Errorf("expected TTL 300, got %d", resp.ApiCache.Ttl)
		}
		if resp.ApiCache.ApiCachingBehavior != types.ApiCachingBehaviorFullRequestCaching {
			return fmt.Errorf("expected FULL_REQUEST_CACHING, got %s", resp.ApiCache.ApiCachingBehavior)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApiCache", func() error {
		resp, err := client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil")
		}
		if resp.ApiCache.Ttl != 300 {
			return fmt.Errorf("expected TTL 300, got %d", resp.ApiCache.Ttl)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApiCache_NonExistent", func() error {
		_, err := client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String("does-not-exist"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateApiCache", func() error {
		resp, err := client.UpdateApiCache(ctx, &appsync.UpdateApiCacheInput{
			ApiId:              aws.String(gqlApiId),
			Type:               types.ApiCacheTypeMedium,
			Ttl:                600,
			ApiCachingBehavior: types.ApiCachingBehaviorPerResolverCaching,
		})
		if err != nil {
			return err
		}
		if resp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil")
		}
		if resp.ApiCache.Type != types.ApiCacheTypeMedium {
			return fmt.Errorf("expected MEDIUM type, got %s", resp.ApiCache.Type)
		}
		if resp.ApiCache.Ttl != 600 {
			return fmt.Errorf("expected TTL 600, got %d", resp.ApiCache.Ttl)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "FlushApiCache", func() error {
		_, err := client.FlushApiCache(ctx, &appsync.FlushApiCacheInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		verifyResp, err := client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return fmt.Errorf("cache should still exist after flush: %w", err)
		}
		if verifyResp.ApiCache == nil {
			return fmt.Errorf("ApiCache is nil after flush")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "FlushApiCache_NonExistent", func() error {
		_, err := client.FlushApiCache(ctx, &appsync.FlushApiCacheInput{
			ApiId: aws.String("does-not-exist"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent cache")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiCache", func() error {
		_, err := client.DeleteApiCache(ctx, &appsync.DeleteApiCacheInput{
			ApiId: aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApiCache(ctx, &appsync.GetApiCacheInput{
			ApiId: aws.String(gqlApiId),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting cache")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApiCache_NonExistent", func() error {
		_, err := client.DeleteApiCache(ctx, &appsync.DeleteApiCacheInput{
			ApiId: aws.String(gqlApiId),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent cache")
		}
		return nil
	}))

	// ================================================================
	// Domain Name CRUD
	// ================================================================
	domainName := fmt.Sprintf("test-domain-%d.example.com", uid)
	var tagDomainName string

	results = append(results, r.RunTest("appsync", "CreateDomainName", func() error {
		resp, err := client.CreateDomainName(ctx, &appsync.CreateDomainNameInput{
			DomainName:     aws.String(domainName),
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/test-cert"),
			Description:    aws.String("test domain"),
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil {
			return fmt.Errorf("DomainNameConfig is nil")
		}
		if *resp.DomainNameConfig.DomainName != domainName {
			return fmt.Errorf("expected domain %s, got %s", domainName, *resp.DomainNameConfig.DomainName)
		}
		if resp.DomainNameConfig.DomainNameArn == nil || *resp.DomainNameConfig.DomainNameArn == "" {
			return fmt.Errorf("DomainNameArn is nil")
		}
		if resp.DomainNameConfig.AppsyncDomainName == nil || *resp.DomainNameConfig.AppsyncDomainName == "" {
			return fmt.Errorf("AppsyncDomainName is nil")
		}
		if resp.DomainNameConfig.HostedZoneId == nil || *resp.DomainNameConfig.HostedZoneId == "" {
			return fmt.Errorf("HostedZoneId is nil")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateDomainName_WithTags", func() error {
		tagDomainName = fmt.Sprintf("tag-domain-%d.example.com", uid)
		resp, err := client.CreateDomainName(ctx, &appsync.CreateDomainNameInput{
			DomainName:     aws.String(tagDomainName),
			CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/tag-cert"),
			Tags:           map[string]string{"env": "prod"},
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil || resp.DomainNameConfig.Tags == nil {
			return fmt.Errorf("Tags not returned")
		}
		if resp.DomainNameConfig.Tags["env"] != "prod" {
			return fmt.Errorf("tag not persisted: %v", resp.DomainNameConfig.Tags)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDomainName", func() error {
		resp, err := client.GetDomainName(ctx, &appsync.GetDomainNameInput{
			DomainName: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil || *resp.DomainNameConfig.DomainName != domainName {
			return fmt.Errorf("expected domain %s, got %v", domainName, resp.DomainNameConfig)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetDomainName_NonExistent", func() error {
		_, err := client.GetDomainName(ctx, &appsync.GetDomainNameInput{
			DomainName: aws.String("does-not-exist.example.com"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListDomainNames", func() error {
		resp, err := client.ListDomainNames(ctx, &appsync.ListDomainNamesInput{})
		if err != nil {
			return err
		}
		if len(resp.DomainNameConfigs) < 2 {
			return fmt.Errorf("expected at least 2 domain names, got %d", len(resp.DomainNameConfigs))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListDomainNames_WithPagination", func() error {
		resp, err := client.ListDomainNames(ctx, &appsync.ListDomainNamesInput{MaxResults: 1})
		if err != nil {
			return err
		}
		if len(resp.DomainNameConfigs) != 1 {
			return fmt.Errorf("expected 1 domain name with maxResults=1, got %d", len(resp.DomainNameConfigs))
		}
		if resp.NextToken == nil {
			return fmt.Errorf("expected NextToken when more results exist")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateDomainName", func() error {
		resp, err := client.UpdateDomainName(ctx, &appsync.UpdateDomainNameInput{
			DomainName:  aws.String(domainName),
			Description: aws.String("updated description"),
		})
		if err != nil {
			return err
		}
		if resp.DomainNameConfig == nil {
			return fmt.Errorf("DomainNameConfig is nil")
		}
		if resp.DomainNameConfig.Description == nil || *resp.DomainNameConfig.Description != "updated description" {
			return fmt.Errorf("description not updated: %v", resp.DomainNameConfig.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateDomainName_NonExistent", func() error {
		_, err := client.UpdateDomainName(ctx, &appsync.UpdateDomainNameInput{
			DomainName: aws.String("does-not-exist.example.com"),
		})
		if err == nil {
			return fmt.Errorf("expected error for non-existent domain name")
		}
		return nil
	}))

	// ================================================================
	// Domain Association
	// ================================================================
	results = append(results, r.RunTest("appsync", "AssociateApi", func() error {
		resp, err := client.AssociateApi(ctx, &appsync.AssociateApiInput{
			DomainName: aws.String(domainName),
			ApiId:      aws.String(gqlApiId),
		})
		if err != nil {
			return err
		}
		if resp.ApiAssociation == nil {
			return fmt.Errorf("ApiAssociation is nil")
		}
		if resp.ApiAssociation.ApiId == nil || *resp.ApiAssociation.ApiId != gqlApiId {
			return fmt.Errorf("expected apiId %s, got %v", gqlApiId, resp.ApiAssociation.ApiId)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetApiAssociation", func() error {
		resp, err := client.GetApiAssociation(ctx, &appsync.GetApiAssociationInput{
			DomainName: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		if resp.ApiAssociation == nil {
			return fmt.Errorf("ApiAssociation is nil")
		}
		if resp.ApiAssociation.ApiId == nil || *resp.ApiAssociation.ApiId != gqlApiId {
			return fmt.Errorf("expected apiId %s, got %v", gqlApiId, resp.ApiAssociation.ApiId)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateApi", func() error {
		_, err := client.DisassociateApi(ctx, &appsync.DisassociateApiInput{
			DomainName: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApiAssociation(ctx, &appsync.GetApiAssociationInput{
			DomainName: aws.String(domainName),
		})
		if err == nil {
			return fmt.Errorf("expected error after disassociating API")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateApi_NotAssociated", func() error {
		_, err := client.DisassociateApi(ctx, &appsync.DisassociateApiInput{
			DomainName: aws.String(domainName),
		})
		if err == nil {
			return fmt.Errorf("expected error for disassociating non-associated domain")
		}
		return nil
	}))

	// ================================================================
	// Merged API operations
	// ================================================================
	var mergedApiId string
	var sourceApiId2 string
	var associationId string

	results = append(results, r.RunTest("appsync", "CreateGraphqlApi_ForMerged", func() error {
		resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
			Name:               aws.String(fmt.Sprintf("merged-api-%d", uid)),
			AuthenticationType: types.AuthenticationTypeApiKey,
			ApiType:            types.GraphQLApiTypeMerged,
		})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil || resp.GraphqlApi.ApiId == nil {
			return fmt.Errorf("invalid response")
		}
		if resp.GraphqlApi.ApiType != types.GraphQLApiTypeMerged {
			return fmt.Errorf("expected MERGED api type, got %s", resp.GraphqlApi.ApiType)
		}
		mergedApiId = *resp.GraphqlApi.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "CreateGraphqlApi_ForSource", func() error {
		resp, err := client.CreateGraphqlApi(ctx, &appsync.CreateGraphqlApiInput{
			Name:               aws.String(fmt.Sprintf("source-api-%d", uid)),
			AuthenticationType: types.AuthenticationTypeApiKey,
		})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil || resp.GraphqlApi.ApiId == nil {
			return fmt.Errorf("invalid response")
		}
		sourceApiId2 = *resp.GraphqlApi.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "AssociateSourceGraphqlApi", func() error {
		resp, err := client.AssociateSourceGraphqlApi(ctx, &appsync.AssociateSourceGraphqlApiInput{
			MergedApiIdentifier: aws.String(mergedApiId),
			SourceApiIdentifier: aws.String(sourceApiId2),
			Description:         aws.String("test association"),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.AssociationId == nil || *resp.SourceApiAssociation.AssociationId == "" {
			return fmt.Errorf("AssociationId is empty")
		}
		if resp.SourceApiAssociation.SourceApiAssociationStatus != types.SourceApiAssociationStatusMergeScheduled {
			return fmt.Errorf("expected MERGE_SCHEDULED, got %s", resp.SourceApiAssociation.SourceApiAssociationStatus)
		}
		if resp.SourceApiAssociation.Description == nil || *resp.SourceApiAssociation.Description != "test association" {
			return fmt.Errorf("description not set: %v", resp.SourceApiAssociation.Description)
		}
		associationId = *resp.SourceApiAssociation.AssociationId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetSourceApiAssociation", func() error {
		resp, err := client.GetSourceApiAssociation(ctx, &appsync.GetSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(mergedApiId),
			AssociationId:       aws.String(associationId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil || *resp.SourceApiAssociation.AssociationId != associationId {
			return fmt.Errorf("expected association %s, got %v", associationId, resp.SourceApiAssociation)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetSourceApiAssociation_NonExistent", func() error {
		_, err := client.GetSourceApiAssociation(ctx, &appsync.GetSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(mergedApiId),
			AssociationId:       aws.String("does-not-exist"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "UpdateSourceApiAssociation", func() error {
		resp, err := client.UpdateSourceApiAssociation(ctx, &appsync.UpdateSourceApiAssociationInput{
			MergedApiIdentifier: aws.String(mergedApiId),
			AssociationId:       aws.String(associationId),
			Description:         aws.String("updated association"),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.Description == nil || *resp.SourceApiAssociation.Description != "updated association" {
			return fmt.Errorf("description not updated: %v", resp.SourceApiAssociation.Description)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "StartSchemaMerge", func() error {
		resp, err := client.StartSchemaMerge(ctx, &appsync.StartSchemaMergeInput{
			MergedApiIdentifier: aws.String(mergedApiId),
			AssociationId:       aws.String(associationId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociationStatus != types.SourceApiAssociationStatusMergeSuccess {
			return fmt.Errorf("expected MERGE_SUCCESS, got %s", resp.SourceApiAssociationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "ListSourceApiAssociations", func() error {
		resp, err := client.ListSourceApiAssociations(ctx, &appsync.ListSourceApiAssociationsInput{
			ApiId: aws.String(sourceApiId2),
		})
		if err != nil {
			return err
		}
		if len(resp.SourceApiAssociationSummaries) < 1 {
			return fmt.Errorf("expected at least 1 association, got %d", len(resp.SourceApiAssociationSummaries))
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateSourceGraphqlApi", func() error {
		resp, err := client.DisassociateSourceGraphqlApi(ctx, &appsync.DisassociateSourceGraphqlApiInput{
			MergedApiIdentifier: aws.String(mergedApiId),
			AssociationId:       aws.String(associationId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociationStatus != types.SourceApiAssociationStatusMergeSuccess {
			return fmt.Errorf("expected MERGE_SUCCESS status, got %s", resp.SourceApiAssociationStatus)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateSourceGraphqlApi_NonExistent", func() error {
		_, err := client.DisassociateSourceGraphqlApi(ctx, &appsync.DisassociateSourceGraphqlApiInput{
			MergedApiIdentifier: aws.String(mergedApiId),
			AssociationId:       aws.String("already-deleted"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// ================================================================
	// AssociateMergedGraphqlApi (source API side)
	// ================================================================
	var mergedAssocId string

	results = append(results, r.RunTest("appsync", "AssociateMergedGraphqlApi", func() error {
		resp, err := client.AssociateMergedGraphqlApi(ctx, &appsync.AssociateMergedGraphqlApiInput{
			SourceApiIdentifier: aws.String(sourceApiId2),
			MergedApiIdentifier: aws.String(mergedApiId),
			Description:         aws.String("merged from source side"),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociation == nil {
			return fmt.Errorf("SourceApiAssociation is nil")
		}
		if resp.SourceApiAssociation.AssociationId == nil || *resp.SourceApiAssociation.AssociationId == "" {
			return fmt.Errorf("AssociationId is empty")
		}
		mergedAssocId = *resp.SourceApiAssociation.AssociationId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateMergedGraphqlApi", func() error {
		resp, err := client.DisassociateMergedGraphqlApi(ctx, &appsync.DisassociateMergedGraphqlApiInput{
			SourceApiIdentifier: aws.String(sourceApiId2),
			AssociationId:       aws.String(mergedAssocId),
		})
		if err != nil {
			return err
		}
		if resp.SourceApiAssociationStatus == "" {
			return fmt.Errorf("expected non-empty status, got empty")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DisassociateMergedGraphqlApi_NonExistent", func() error {
		_, err := client.DisassociateMergedGraphqlApi(ctx, &appsync.DisassociateMergedGraphqlApiInput{
			SourceApiIdentifier: aws.String(sourceApiId2),
			AssociationId:       aws.String("already-deleted"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	// ================================================================
	// Cleanup: Delete domain names
	// ================================================================
	results = append(results, r.RunTest("appsync", "DeleteDomainName", func() error {
		_, err := client.DeleteDomainName(ctx, &appsync.DeleteDomainNameInput{
			DomainName: aws.String(domainName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetDomainName(ctx, &appsync.GetDomainNameInput{
			DomainName: aws.String(domainName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting domain name")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteDomainName_NonExistent", func() error {
		_, err := client.DeleteDomainName(ctx, &appsync.DeleteDomainNameInput{
			DomainName: aws.String("already-deleted.example.com"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteDomainName_WithTags", func() error {
		_, err := client.DeleteDomainName(ctx, &appsync.DeleteDomainNameInput{
			DomainName: aws.String(tagDomainName),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	// ================================================================
	// Cleanup: Delete GraphQL APIs
	// ================================================================
	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi", func() error {
		if descApiKeyId != "" {
			client.DeleteApiKey(ctx, &appsync.DeleteApiKeyInput{
				ApiId: aws.String(gqlApiId),
				Id:    aws.String(descApiKeyId),
			})
		}
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(gqlApiId)})
		if err != nil {
			return err
		}
		_, err = client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(gqlApiId)})
		if err == nil {
			return fmt.Errorf("expected error after deleting GraphQL API")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_NonExistent", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String("already-deleted")})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_WithTags", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(gqlTagsApiId)})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_Merged", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(mergedApiId)})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteGraphqlApi_Source", func() error {
		_, err := client.DeleteGraphqlApi(ctx, &appsync.DeleteGraphqlApiInput{ApiId: aws.String(sourceApiId2)})
		if err != nil {
			return err
		}
		return nil
	}))

	// ================================================================
	// Cleanup: Delete Event API (v2) resources
	// ================================================================
	results = append(results, r.RunTest("appsync", "DeleteChannelNamespace", func() error {
		_, err := client.DeleteChannelNamespace(ctx, &appsync.DeleteChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(nsName),
		})
		if err != nil {
			return err
		}
		_, err = client.GetChannelNamespace(ctx, &appsync.GetChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(nsName),
		})
		if err == nil {
			return fmt.Errorf("expected error after deleting namespace")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteChannelNamespace_NonExistent", func() error {
		_, err := client.DeleteChannelNamespace(ctx, &appsync.DeleteChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String("already-deleted"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteChannelNamespace_ForTagging", func() error {
		_, err := client.DeleteChannelNamespace(ctx, &appsync.DeleteChannelNamespaceInput{
			ApiId: aws.String(apiId),
			Name:  aws.String(tagNsName),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	// === PAGINATION FOLLOW-UP TESTS ===

	results = append(results, r.RunTest("appsync", "ListApis_NextTokenFollowUp", func() error {
		var allApis []string
		nextToken := ""
		pageCount := 0
		for {
			input := &appsync.ListApisInput{MaxResults: 1}
			if nextToken != "" {
				input.NextToken = aws.String(nextToken)
			}
			resp, err := client.ListApis(ctx, input)
			if err != nil {
				return fmt.Errorf("list apis page: %v", err)
			}
			pageCount++
			for _, api := range resp.Apis {
				if api.Name != nil {
					allApis = append(allApis, *api.Name)
				}
			}
			if resp.NextToken != nil && *resp.NextToken != "" {
				nextToken = *resp.NextToken
			} else {
				break
			}
		}
		if pageCount < 2 {
			return fmt.Errorf("expected at least 2 pages for ListApis with MaxResults=1, got %d", pageCount)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(apiId),
		})
		if err != nil {
			return err
		}
		_, err = client.GetApi(ctx, &appsync.GetApiInput{ApiId: aws.String(apiId)})
		if err == nil {
			return fmt.Errorf("expected error after deleting API")
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_NonExistent", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String("already-deleted"),
		})
		if err := AssertErrorContains(err, "NotFoundException"); err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_WithTags", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(tagsApiId),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_ForTagging", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(taggedApiId),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "DeleteApi_WithOwnerContact", func() error {
		_, err := client.DeleteApi(ctx, &appsync.DeleteApiInput{
			ApiId: aws.String(ownerApiId),
		})
		if err != nil {
			return err
		}
		return nil
	}))

	return results
}
