package cognitoidentityprovider

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

func findTokenByValue[T any](store *common.BaseStore, tokenValue string, getToken func(*T) string, getExpires func(*T) time.Time) (*T, error) {
	var found *T
	var foundKey string
	err := store.ForEach(func(key string, value []byte) error {
		var t T
		if err := json.Unmarshal(value, &t); err != nil {
			return err
		}
		if getToken(&t) == tokenValue {
			cp := t
			found = &cp
			foundKey = key
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if found == nil {
		return nil, ErrTokenNotFound
	}
	if time.Now().After(getExpires(found)) {
		_ = store.Delete(foundKey)
		return nil, ErrTokenExpired
	}
	return found, nil
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

func userIndexKey(userID string) string {
	return "useridx:" + userID
}

func providerIndexKey(userPoolID, providerName, providerAttrValue string) string {
	return "providx:" + userPoolID + "#" + providerName + "#" + providerAttrValue
}

func clientIndexKey(clientID string) string {
	return "clientidx:" + clientID
}

func encodePrivateKeyToPEM(key *rsa.PrivateKey) string {
	der := x509.MarshalPKCS1PrivateKey(key)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: der,
	}
	return string(pem.EncodeToMemory(block))
}

func encodePublicKeyToPEM(key *rsa.PublicKey) string {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return ""
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}
	return string(pem.EncodeToMemory(block))
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

func domainKey(domain string) string {
	return "domain:" + domain
}
