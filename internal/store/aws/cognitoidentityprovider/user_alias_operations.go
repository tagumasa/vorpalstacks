package cognitoidentityprovider

import (
	"strings"
	"time"

	"vorpalstacks/internal/store/aws/common"
)

// aliasAttributeOrder is the fixed resolution order of the alias attributes,
// so that an identifier claimed under more than one attribute always resolves
// to the same username. It mirrors the AliasAttributeType enum values.
var aliasAttributeOrder = []string{"email", "phone_number", "preferred_username"}

// aliasAttributeSet returns the attributes whose values are sign-in
// identifiers with pool-wide uniqueness: the union of the pool's
// AliasAttributes and UsernameAttributes. Both configurations accept only
// email, phone_number and preferred_username (the AliasAttributeType and
// UsernameAttributeType enums), enforced when the pool is created or
// updated.
func aliasAttributeSet(pool *UserPool) map[string]bool {
	set := make(map[string]bool, len(pool.AliasAttributes)+len(pool.UsernameAttributes))
	for _, a := range pool.AliasAttributes {
		set[a] = true
	}
	for _, a := range pool.UsernameAttributes {
		set[a] = true
	}
	return set
}

// userAliasClaims returns the alias values the user actively claims, as
// attribute name → value. An email or phone_number claims its value only
// while verified: Amazon Cognito lets unverified duplicates coexist so the
// public sign-up API cannot be used to enumerate accounts, and the conflict
// surfaces when the second holder verifies (ConfirmSignUp,
// VerifyUserAttribute) or a write activates the verified flag.
// preferred_username has no verification round and claims on presence.
func userAliasClaims(pool *UserPool, attrs map[string]string) map[string]string {
	claims := make(map[string]string)
	for attr := range aliasAttributeSet(pool) {
		value := attrs[attr]
		if value == "" {
			continue
		}
		if attr == "email" || attr == "phone_number" {
			if attrs[attr+"_verified"] == "true" {
				claims[attr] = value
			}
			continue
		}
		claims[attr] = value
	}
	return claims
}

// aliasIndexKey is the alias uniqueness index: it maps a claimed attribute
// value to the claiming username. Values are lowercased so that alias
// sign-in matches case-insensitively, following the case-insensitive
// username index precedent.
func aliasIndexKey(userPoolID, attribute, value string) string {
	return "aliasidx:" + userPoolID + "#" + attribute + "#" + strings.ToLower(value)
}

// claimUserAliases writes the index entries for the user's alias claims. A
// value another user already claims rejects the write with ErrAliasExists,
// or, when migrate is set, is taken over from the previous holder — the
// ForceAliasCreation semantics the model documents for AdminCreateUser: "if
// this parameter is set to True and the phone number or email address
// specified in the UserAttributes parameter already exists as an alias with
// a different user, this request migrates the alias from the previous user
// to the newly-created user. The previous user will no longer be able to
// log in using that alias." The caller must hold the user-write lock.
func (s *CognitoStore) claimUserAliases(pool *UserPool, user *User, migrate bool) error {
	for attr, value := range userAliasClaims(pool, user.Attributes) {
		key := aliasIndexKey(pool.ID, attr, value)
		var holder string
		if err := s.usersStore.Get(key, &holder); err != nil && !common.IsNotFound(err) {
			// The uniqueness gate fails closed: an unreadable entry must
			// not read as an unclaimed value.
			return err
		}
		if holder != "" && holder != user.Username {
			if !migrate {
				return ErrAliasExists
			}
			previous, err := s.GetUser(pool.ID, holder)
			if err == nil {
				// Releasing the claim un-verifies a contact attribute — the
				// minimal change the documented outcome requires within the
				// verified-claim model — and removes an attribute with no
				// verification round outright.
				if attr == "email" || attr == "phone_number" {
					previous.Attributes[attr+"_verified"] = "false"
				} else {
					delete(previous.Attributes, attr)
				}
				previous.LastModifiedDate = time.Now().UTC()
				if err := s.usersStore.Put(userPoolUserKey(pool.ID, previous.Username), previous); err != nil {
					return err
				}
			}
		}
		// A failed claim write would leave the alias unrecorded, so a later
		// conflicting claim is not rejected — the write propagates.
		if err := s.usersStore.Put(key, user.Username); err != nil {
			return err
		}
	}
	return nil
}

// releaseAliasClaims deletes the index entries for the supplied claims,
// except those in keep. Each deletion is holder-checked so that an entry
// already taken over by another user is left alone. The caller must hold
// the user-write lock.
func (s *CognitoStore) releaseAliasClaims(userPoolID, username string, claims, keep map[string]string) {
	for attr, value := range claims {
		if keep[attr] == value {
			continue
		}
		key := aliasIndexKey(userPoolID, attr, value)
		var holder string
		if err := s.usersStore.Get(key, &holder); err == nil && holder == username {
			_ = s.usersStore.Delete(key)
		}
	}
}

// resolveAliasUsername resolves a client-supplied identifier to the username
// of the user actively claiming it as an alias value; an empty string means
// no alias matches. The model's standard Username-parameter documentation
// ("The value of this parameter is typically your user's username, but it
// can be any of their alias attributes") makes this the general resolution
// for user-addressing operations when the literal username does not match.
func (s *CognitoStore) resolveAliasUsername(userPoolID, identifier string) string {
	pool, err := s.GetUserPool(userPoolID)
	if err != nil {
		return ""
	}
	configured := aliasAttributeSet(pool)
	for _, attr := range aliasAttributeOrder {
		if !configured[attr] {
			continue
		}
		var username string
		if err := s.usersStore.Get(aliasIndexKey(userPoolID, attr, identifier), &username); err == nil && username != "" {
			return username
		}
	}
	return ""
}
