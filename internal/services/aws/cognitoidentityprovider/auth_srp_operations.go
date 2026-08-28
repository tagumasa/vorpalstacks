package cognitoidentityprovider

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// computeSrpVerifier derives a fresh random salt and the matching SRP verifier
// for the supplied password. It must be invoked at every site that stores a
// password hash so that the user can later authenticate via USER_SRP_AUTH.
//
// userPoolID is the full Cognito pool ID (e.g. "us-east-1_abc123"); the part
// after the underscore (poolName) is required by Cognito's SRP variant inner
// hash. The returned saltHex and verifierHex are lowercase hex strings suitable
// for direct assignment to User.SrpSalt and User.SrpVerifier.
func computeSrpVerifier(userPoolID, username, password string) (saltHex, verifierHex string, err error) {
	idx := strings.Index(userPoolID, "_")
	if idx < 0 || idx == len(userPoolID)-1 {
		return "", "", fmt.Errorf("invalid user pool ID %q: missing region prefix", userPoolID)
	}
	poolName := userPoolID[idx+1:]
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}
	saltHex = hex.EncodeToString(salt)
	v := ComputeVerifier(saltHex, poolName, username, password)
	return saltHex, hex.EncodeToString(v.Bytes()), nil
}

// poolNameFromID extracts the portion of a Cognito user pool ID after the
// underscore (e.g. "us-east-1_abc123" => "abc123"). The pool name is used as
// part of the Cognito SRP inner hash and the claim message. The boolean is
// false when the ID does not contain a valid region/name separator.
func poolNameFromID(userPoolID string) (string, bool) {
	idx := strings.Index(userPoolID, "_")
	if idx < 0 || idx == len(userPoolID)-1 {
		return "", false
	}
	return userPoolID[idx+1:], true
}
