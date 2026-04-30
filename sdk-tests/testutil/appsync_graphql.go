package testutil

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
)

func (r *TestRunner) runAppSyncGraphqlApiTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client
	uid := res.uid

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
			return fmt.Errorf("arn is nil")
		}
		if !strings.HasPrefix(*resp.GraphqlApi.Arn, "arn:aws:appsync:") {
			return fmt.Errorf("invalid ARN format: %s", *resp.GraphqlApi.Arn)
		}
		res.gqlApiId = *resp.GraphqlApi.ApiId
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
			return fmt.Errorf("tags not returned")
		}
		if resp.GraphqlApi.Tags["env"] != "dev" {
			return fmt.Errorf("tag not persisted: %v", resp.GraphqlApi.Tags)
		}
		res.gqlTagsApiId = *resp.GraphqlApi.ApiId
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetGraphqlApi", func() error {
		resp, err := client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String(res.gqlApiId)})
		if err != nil {
			return err
		}
		if resp.GraphqlApi == nil {
			return fmt.Errorf("GraphqlApi is nil")
		}
		if *resp.GraphqlApi.ApiId != res.gqlApiId {
			return fmt.Errorf("expected ApiId %s, got %s", res.gqlApiId, *resp.GraphqlApi.ApiId)
		}
		if resp.GraphqlApi.AuthenticationType != types.AuthenticationTypeApiKey {
			return fmt.Errorf("expected API_KEY auth type, got %s", resp.GraphqlApi.AuthenticationType)
		}
		return nil
	}))

	results = append(results, r.RunTest("appsync", "GetGraphqlApi_NonExistent", func() error {
		_, err := client.GetGraphqlApi(ctx, &appsync.GetGraphqlApiInput{ApiId: aws.String("does-not-exist")})
		return AssertErrorContains(err, "NotFoundException")
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
			if *api.ApiId == res.gqlApiId {
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
			ApiId:              aws.String(res.gqlApiId),
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
		return AssertErrorContains(err, "NotFoundException")
	}))

	results = append(results, r.runAppSyncSchemaTests(res)...)
	results = append(results, r.runAppSyncEnvironmentVariableTests(res)...)

	return results
}

func (r *TestRunner) runAppSyncSchemaTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "StartSchemaCreation", func() error {
		schemaSDL := []byte("type Query { hello: String } type Mutation { addPost(title: String!): Post } type Post { id: ID! title: String! }")
		resp, err := client.StartSchemaCreation(ctx, &appsync.StartSchemaCreationInput{
			ApiId:      aws.String(res.gqlApiId),
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
			ApiId: aws.String(res.gqlApiId),
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
			ApiId:  aws.String(res.gqlApiId),
			Format: types.OutputTypeSdl,
		})
		if err != nil {
			return err
		}
		schemaStr := string(resp.Schema)
		if len(schemaStr) == 0 {
			return fmt.Errorf("schema is empty")
		}
		for _, expected := range []string{"type Query", "type Mutation", "type Post"} {
			if !strings.Contains(schemaStr, expected) {
				return fmt.Errorf("schema missing %q", expected)
			}
		}
		return nil
	}))

	return results
}

func (r *TestRunner) runAppSyncEnvironmentVariableTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	results = append(results, r.RunTest("appsync", "PutGraphqlApiEnvironmentVariables", func() error {
		resp, err := client.PutGraphqlApiEnvironmentVariables(ctx, &appsync.PutGraphqlApiEnvironmentVariablesInput{
			ApiId: aws.String(res.gqlApiId),
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
			ApiId: aws.String(res.gqlApiId),
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
			ApiId: aws.String(res.gqlApiId),
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

	return results
}
