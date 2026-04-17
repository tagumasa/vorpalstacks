package cognitoidentityprovider

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

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

func userPoolUserKey(userPoolID, username string) string {
	return userPoolID + "#" + username
}

func userPoolGroupKey(userPoolID, groupName string) string {
	return userPoolID + "#" + groupName
}

func userPoolClientKey(userPoolID, clientID string) string {
	return userPoolID + "#" + clientID
}

func refreshTokenKey(userPoolID, userID, token string) string {
	return userPoolID + "#" + userID + "#" + token
}

func idTokenKey(userPoolID, userID, token string) string {
	return userPoolID + "#" + userID + "#" + token
}

func accessTokenKey(userPoolID, userID, token string) string {
	return userPoolID + "#" + userID + "#" + token
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
