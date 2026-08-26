package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"

	"vorpalstacks-sdk-tests/config"
)

type appsyncResources struct {
	ctx    context.Context
	client *appsync.Client
	uid    int64

	apiId      string
	tagsApiId  string
	ownerApiId string

	nsName string

	taggedApiId  string
	taggedApiArn string

	nsArn     string
	tagNsName string

	gqlApiId     string
	gqlTagsApiId string

	functionId string

	apiKeyId     string
	descApiKeyId string

	domainName    string
	tagDomainName string

	mergedApiId   string
	sourceApiId2  string
	associationId string

	mergedAssocId string
}

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

	res := &appsyncResources{
		ctx:    context.Background(),
		client: appsync.NewFromConfig(cfg),
		uid:    time.Now().UnixNano(),
	}

	results = append(results, r.runAppSyncEventApiTests(res)...)
	results = append(results, r.runAppSyncChannelTests(res)...)
	results = append(results, r.runAppSyncTagTests(res)...)
	results = append(results, r.runAppSyncGraphqlApiTests(res)...)
	results = append(results, r.runAppSyncDataSourceTests(res)...)
	results = append(results, r.runAppSyncTypeTests(res)...)
	results = append(results, r.runAppSyncResolverTests(res)...)
	results = append(results, r.runAppSyncVTLTests(res)...)
	results = append(results, r.runAppSyncApiKeyCacheTests(res)...)
	results = append(results, r.runAppSyncGraphQLAuthTests(res)...)
	results = append(results, r.runAppSyncDomainTests(res)...)
	results = append(results, r.runAppSyncCleanupTests(res)...)

	return results
}

func minEventConfig() *types.EventConfig {
	return &types.EventConfig{
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
}

// allApis walks every page of ListApis and returns all APIs.
func (res *appsyncResources) allApis() ([]types.Api, error) {
	return paginate(func(next *string) ([]types.Api, *string, error) {
		resp, err := res.client.ListApis(res.ctx, &appsync.ListApisInput{
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Apis, resp.NextToken, nil
	})
}

// allGraphqlApis walks every page of ListGraphqlApis and returns all APIs.
func (res *appsyncResources) allGraphqlApis() ([]types.GraphqlApi, error) {
	return paginate(func(next *string) ([]types.GraphqlApi, *string, error) {
		resp, err := res.client.ListGraphqlApis(res.ctx, &appsync.ListGraphqlApisInput{
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.GraphqlApis, resp.NextToken, nil
	})
}

// allDataSources walks every page of ListDataSources for one API.
func (res *appsyncResources) allDataSources(apiID string) ([]types.DataSource, error) {
	return paginate(func(next *string) ([]types.DataSource, *string, error) {
		resp, err := res.client.ListDataSources(res.ctx, &appsync.ListDataSourcesInput{
			ApiId:     aws.String(apiID),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DataSources, resp.NextToken, nil
	})
}

// allResolvers walks every page of ListResolvers for one API type.
func (res *appsyncResources) allResolvers(apiID, typeName string) ([]types.Resolver, error) {
	return paginate(func(next *string) ([]types.Resolver, *string, error) {
		resp, err := res.client.ListResolvers(res.ctx, &appsync.ListResolversInput{
			ApiId:     aws.String(apiID),
			TypeName:  aws.String(typeName),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Resolvers, resp.NextToken, nil
	})
}

// allFunctions walks every page of ListFunctions for one API.
func (res *appsyncResources) allFunctions(apiID string) ([]types.FunctionConfiguration, error) {
	return paginate(func(next *string) ([]types.FunctionConfiguration, *string, error) {
		resp, err := res.client.ListFunctions(res.ctx, &appsync.ListFunctionsInput{
			ApiId:     aws.String(apiID),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Functions, resp.NextToken, nil
	})
}

// allResolversByFunction walks every page of ListResolversByFunction.
func (res *appsyncResources) allResolversByFunction(apiID, functionID string) ([]types.Resolver, error) {
	return paginate(func(next *string) ([]types.Resolver, *string, error) {
		resp, err := res.client.ListResolversByFunction(res.ctx, &appsync.ListResolversByFunctionInput{
			ApiId:      aws.String(apiID),
			FunctionId: aws.String(functionID),
			NextToken:  next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Resolvers, resp.NextToken, nil
	})
}

// allApiKeys walks every page of ListApiKeys for one API.
func (res *appsyncResources) allApiKeys(apiID string) ([]types.ApiKey, error) {
	return paginate(func(next *string) ([]types.ApiKey, *string, error) {
		resp, err := res.client.ListApiKeys(res.ctx, &appsync.ListApiKeysInput{
			ApiId:     aws.String(apiID),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.ApiKeys, resp.NextToken, nil
	})
}

// allChannelNamespaces walks every page of ListChannelNamespaces for one API.
func (res *appsyncResources) allChannelNamespaces(apiID string) ([]types.ChannelNamespace, error) {
	return paginate(func(next *string) ([]types.ChannelNamespace, *string, error) {
		resp, err := res.client.ListChannelNamespaces(res.ctx, &appsync.ListChannelNamespacesInput{
			ApiId:     aws.String(apiID),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.ChannelNamespaces, resp.NextToken, nil
	})
}

// allDomainNames walks every page of ListDomainNames.
func (res *appsyncResources) allDomainNames() ([]types.DomainNameConfig, error) {
	return paginate(func(next *string) ([]types.DomainNameConfig, *string, error) {
		resp, err := res.client.ListDomainNames(res.ctx, &appsync.ListDomainNamesInput{
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.DomainNameConfigs, resp.NextToken, nil
	})
}

// allTypes walks every page of ListTypes for one API and definition format.
func (res *appsyncResources) allTypes(apiID string, format types.TypeDefinitionFormat) ([]types.Type, error) {
	return paginate(func(next *string) ([]types.Type, *string, error) {
		resp, err := res.client.ListTypes(res.ctx, &appsync.ListTypesInput{
			ApiId:     aws.String(apiID),
			Format:    format,
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Types, resp.NextToken, nil
	})
}

// allSourceApiAssociations walks every page of ListSourceApiAssociations for one merged API.
func (res *appsyncResources) allSourceApiAssociations(mergedApiID string) ([]types.SourceApiAssociationSummary, error) {
	return paginate(func(next *string) ([]types.SourceApiAssociationSummary, *string, error) {
		resp, err := res.client.ListSourceApiAssociations(res.ctx, &appsync.ListSourceApiAssociationsInput{
			ApiId:     aws.String(mergedApiID),
			NextToken: next,
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.SourceApiAssociationSummaries, resp.NextToken, nil
	})
}
