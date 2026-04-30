package testutil

import (
	"context"
	"fmt"
	"time"

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
