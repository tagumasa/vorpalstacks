package cognitoidentityprovider

import (
	"encoding/json"
	"strings"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// findTokenByValue resolves a token value through the token-value index and
// loads the record it points at: one index read plus one record read, with
// no bucket scan. An index miss means the value identifies no token — the
// index is written together with the record, so a record without an index
// entry is not resurrected by scanning. An index entry whose record is gone
// is dropped. An expired record takes its index entry with it.
func findTokenByValue[T any](store *common.BaseStore, tokenValue string, getExpires func(*T) time.Time) (*T, error) {
	var key string
	if err := store.Get(tokenIndexKey(tokenValue), &key); err != nil {
		return nil, ErrTokenNotFound
	}
	var t T
	if err := store.Get(key, &t); err != nil {
		_ = store.Delete(tokenIndexKey(tokenValue))
		return nil, ErrTokenNotFound
	}
	if time.Now().After(getExpires(&t)) {
		_ = store.Delete(key)
		_ = store.Delete(tokenIndexKey(tokenValue))
		return nil, ErrTokenExpired
	}
	return &t, nil
}

// deleteTokenEntry removes a token record and its value-index entry in one
// step. Scan-driven deletions (per-user token sweeps and the pool cascade)
// receive the record's bytes, from which the index key derives; a value
// that cannot be parsed still loses its record key.
func deleteTokenEntry(store *common.BaseStore, key string, value []byte) error {
	var rec struct {
		Token string
	}
	if err := json.Unmarshal(value, &rec); err == nil && rec.Token != "" {
		_ = store.Delete(tokenIndexKey(rec.Token))
	}
	return store.Delete(key)
}

func userPoolBucketName(region string) string {
	return "cognito-userpools-" + region
}

func userBucketName(region string) string {
	return "cognito-users-" + region
}

func groupBucketName(region string) string {
	return "cognito-groups-" + region
}

func clientBucketName(region string) string {
	return "cognito-clients-" + region
}

func refreshTokenBucketName(region string) string {
	return "cognito-refreshtokens-" + region
}

func idTokenBucketName(region string) string {
	return "cognito-idtokens-" + region
}

func accessTokenBucketName(region string) string {
	return "cognito-accesstokens-" + region
}

func challengeSessionBucketName(region string) string {
	return "cognito-challengesessions-" + region
}

func deviceBucketName(region string) string {
	return "cognito-devices-" + region
}

func deviceKey(userPoolID, userID, deviceKey string) string {
	return userPoolID + "#" + userID + "#" + deviceKey
}

func devicePrefix(userPoolID, userID string) string {
	return userPoolID + "#" + userID + "#"
}

func authEventBucketName(region string) string {
	return "cognito-authevents-" + region
}

func authEventKey(userPoolID, userID, eventID string) string {
	return userPoolID + "#" + userID + "#" + eventID
}

func authEventPrefix(userPoolID, userID string) string {
	return userPoolID + "#" + userID + "#"
}

func logDeliveryKey(userPoolID string) string {
	return "logdelivery:" + userPoolID
}

func riskConfigKey(userPoolID, clientID string) string {
	if clientID == "" {
		return "riskconfig:" + userPoolID
	}
	return "riskconfig:" + userPoolID + "#" + clientID
}

func uiCustomizationKey(userPoolID, clientID string) string {
	if clientID == "" {
		return "uicustomization:" + userPoolID
	}
	return "uicustomization:" + userPoolID + "#" + clientID
}

func userImportJobBucketName(region string) string {
	return "cognito-userimportjobs-" + region
}

func userImportJobKey(userPoolID, jobID string) string {
	return userPoolID + "#" + jobID
}

func userImportJobPrefix(userPoolID string) string {
	return userPoolID + "#"
}

func webauthnCredentialBucketName(region string) string {
	return "cognito-webauthn-" + region
}

func webauthnKey(userPoolID, userID, credID string) string {
	return userPoolID + "#" + userID + "#" + credID
}

func webauthnPrefix(userPoolID, userID string) string {
	return userPoolID + "#" + userID + "#"
}

func managedLoginBrandingKey(userPoolID, brandingID string) string {
	return "managedlogin:" + userPoolID + "#" + brandingID
}

func managedLoginBrandingPrefix(userPoolID string) string {
	return "managedlogin:" + userPoolID + "#"
}

func termsKey(userPoolID, termsID string) string {
	return "terms:" + userPoolID + "#" + termsID
}

func termsPrefix(userPoolID string) string {
	return "terms:" + userPoolID + "#"
}

func userPoolReplicaKey(userPoolID, regionName string) string {
	return "replica:" + userPoolID + "#" + regionName
}

func userPoolReplicaPrefix(userPoolID string) string {
	return "replica:" + userPoolID + "#"
}

func userPoolUserKey(userPoolID, username string) string {
	return userPoolID + "#" + username
}

func userPoolGroupKey(userPoolID, groupName string) string {
	return userPoolID + "#" + groupName
}

func userPoolClientKey(userPoolID, clientID string) string {
	return userPoolID + "#" + clientID
}

func tokenKey(userPoolID, userID, token string) string {
	return userPoolID + "#" + userID + "#" + token
}

// tokenIndexKey is the token-value index: it maps a token's opaque value to
// the record's primary key, so value lookups (every authenticated operation
// resolves its access token this way) cost one read instead of a full-bucket
// scan. One index lives in each token bucket, scoped to that token family
// exactly like the records it points at.
func tokenIndexKey(tokenValue string) string {
	return "tokenidx:" + tokenValue
}

// activeImportJobKey is the region's active-import marker: it holds the
// primary key ("<poolID>#<jobID>") of the import job that currently owns
// the single active slot, letting the one-active-job start guard answer
// from one read instead of walking every job in the bucket. The referenced
// job record remains the authority: a marker whose job is gone or terminal
// is stale and is taken over by the next start.
const activeImportJobKey = "importjob-active"

func userIndexKey(userID string) string {
	return "useridx:" + userID
}

func providerIndexKey(userPoolID, providerName, providerAttrValue string) string {
	return "providx:" + userPoolID + "#" + providerName + "#" + providerAttrValue
}

func clientIndexKey(clientID string) string {
	return "clientidx:" + clientID
}

func resourceServerKey(userPoolID, identifier string) string {
	return "resourceserver:" + userPoolID + "#" + identifier
}

func resourceServerPrefix(userPoolID string) string {
	return "resourceserver:" + userPoolID + "#"
}

func identityProviderKey(userPoolID, providerName string) string {
	return "identityprovider:" + userPoolID + "#" + providerName
}

func identityProviderPrefix(userPoolID string) string {
	return "identityprovider:" + userPoolID + "#"
}

// domainKey folds the domain to lowercase: the production router
// lowercases the Host header before domain extraction, so storage, lookup
// and deletion all fold the same way and a mixed-case domain cannot
// resolve differently per mount path.
func domainKey(domain string) string {
	return "domain:" + strings.ToLower(domain)
}

func provisionedLimitKey(category string) string {
	return "provisionedlimit:" + category
}
