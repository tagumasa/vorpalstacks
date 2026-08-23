// Package cognito provides storage layer for AWS Cognito service entities
// including user pools, users, groups, and tokens.
package cognitoidentityprovider

import (
	"sync"

	"vorpalstacks/internal/core/storage"
	"vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// CognitoStore provides Cognito storage operations.
type CognitoStore struct {
	*common.BaseStore
	usersStore             *common.BaseStore
	groupsStore            *common.BaseStore
	clientsStore           *common.BaseStore
	refreshTokensStore     *common.BaseStore
	idTokensStore          *common.BaseStore
	accessTokensStore      *common.BaseStore
	challengeSessionsStore *common.BaseStore
	devicesStore           *common.BaseStore
	authEventsStore        *common.BaseStore
	userImportJobsStore    *common.BaseStore
	webauthnStore          *common.BaseStore
	*common.TagStore
	arnBuilder  *svcarn.ARNBuilder
	accountID   string
	region      string
	groupMu     sync.Mutex
	createMu    sync.Mutex
	importJobMu sync.Mutex
}

// NewCognitoStore creates a new Cognito identity provider store.
func NewCognitoStore(store storage.BasicStorage, accountID, region string) *CognitoStore {
	return &CognitoStore{
		BaseStore:              common.NewBaseStore(store.Bucket(userPoolBucketName(region)), "cognito-userpools"),
		usersStore:             common.NewBaseStore(store.Bucket(userBucketName(region)), "cognito-users"),
		groupsStore:            common.NewBaseStore(store.Bucket(groupBucketName(region)), "cognito-groups"),
		clientsStore:           common.NewBaseStore(store.Bucket(clientBucketName(region)), "cognito-clients"),
		refreshTokensStore:     common.NewBaseStore(store.Bucket(refreshTokenBucketName(region)), "cognito-refreshtokens"),
		idTokensStore:          common.NewBaseStore(store.Bucket(idTokenBucketName(region)), "cognito-idtokens"),
		accessTokensStore:      common.NewBaseStore(store.Bucket(accessTokenBucketName(region)), "cognito-accesstokens"),
		challengeSessionsStore: common.NewBaseStore(store.Bucket(challengeSessionBucketName(region)), "cognito-challengesessions"),
		devicesStore:           common.NewBaseStore(store.Bucket(deviceBucketName(region)), "cognito-devices"),
		authEventsStore:        common.NewBaseStore(store.Bucket(authEventBucketName(region)), "cognito-authevents"),
		userImportJobsStore:    common.NewBaseStore(store.Bucket(userImportJobBucketName(region)), "cognito-userimportjobs"),
		webauthnStore:          common.NewBaseStore(store.Bucket(webauthnCredentialBucketName(region)), "cognito-webauthn"),
		TagStore:               common.NewTagStoreWithRegion(store, "cognito", region),
		arnBuilder:             svcarn.NewARNBuilder(accountID, region),
		accountID:              accountID,
		region:                 region,
	}
}

func (s *CognitoStore) buildUserPoolArn(userPoolID string) string {
	return s.arnBuilder.Cognito().UserPool(userPoolID)
}
