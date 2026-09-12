package cognitoidentityprovider

import (
	"strings"

	"github.com/google/uuid"
)

// The store's identifier generators, one set for the whole package. Every
// identifier is a UUID: record keys use single UUIDs, the two long opaque
// values the wire surfaces treat as secrets (client secrets and token
// values) concatenate two UUIDs, and pool IDs carry the region prefix the
// service's addressing relies on.

func generateClientID() string {
	return uuid.New().String()
}

// generateClientSecret produces a value in the Smithy ClientSecretType
// shape: at most 64 characters drawn from [\w+]. Two hyphen-free UUIDs
// concatenate to exactly 64 hex characters.
func generateClientSecret() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "") + strings.ReplaceAll(uuid.New().String(), "-", "")
}

// keyIDPrefixLength is the UUID prefix length used for pool key IDs and
// region-stripped pool IDs alike.
const keyIDPrefixLength = 8

func generateUserPoolID(region string) string {
	id := uuid.New().String()
	return region + "_" + id[:keyIDPrefixLength]
}

func generateID() string {
	return uuid.New().String()
}

func generateToken() string {
	return uuid.New().String() + uuid.New().String()
}
