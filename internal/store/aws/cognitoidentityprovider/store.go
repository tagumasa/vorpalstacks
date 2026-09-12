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
	arnBuilder *svcarn.ARNBuilder
	accountID  string
	region     string
	// recordMu serialises every user- and group-record write — creation,
	// attribute updates, deletion, group-membership maintenance and app
	// client creation. One record family, one guard: user records are
	// rewritten whole by writers entering from several operations, and two
	// writers on different mutexes would silently overwrite each other's
	// changes. The lock order against the pool lock is fixed:
	// poolKeyLocker → recordMu (pool deletion takes both, in that order).
	recordMu    sync.Mutex
	importJobMu sync.Mutex
	// domainMu serialises user-pool domain mutations. The domain binding
	// rules (a domain string belongs to at most one pool, a pool owns at
	// most one domain) are check-then-write sequences over the shared
	// "domain:" key space, so concurrent creates, updates and deletes must
	// not interleave.
	domainMu sync.Mutex
	// poolKeyLocker serialises mutations of a single user-pool record
	// (updates, MFA configuration, schema additions, deletion) so that
	// read-modify-write cycles cannot lose each other's changes and a
	// deleted pool cannot be resurrected by a racing writer.
	poolKeyLocker common.KeyLocker
	// usernameCaseCache memoises each pool's username-case sensitivity —
	// a bool consulted before every username-keyed store call (user reads,
	// writes, deletes, group membership). The full pool record it derives
	// from is orders of magnitude costlier to decode than the flag is to
	// consult, so the flag is cached and kept exact by the pool-record write
	// paths: an update stores the fresh value, a deletion evicts the entry.
	usernameCaseCache sync.Map
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
