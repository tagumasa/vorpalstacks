// Scoped pagination tokens: request-identity-carrying continuation
// tokens shared by services whose page cursors must survive re-reads of
// the data they walk. A scoped token serialises to an opaque base64
// string; internally it is the JSON form of a struct that names the
// request that produced it (resource, filters, window) plus a value
// cursor, so a token replayed against a different request is rejected
// instead of silently repositioning the walk.

package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// ErrInvalidScopedToken reports a continuation token that is not a
// decodable scoped token of the caller's vocabulary. Consumers surface it
// as their operation's documented invalid-token error.
var ErrInvalidScopedToken = errors.New("invalid scoped pagination token")

// ScopedToken is the contract a scoped-token struct fulfils: it reports
// the vocabulary version it was encoded with. A struct whose Version
// field did not survive decoding (a foreign or legacy token) reports
// zero and fails validation.
type ScopedToken interface {
	TokenVersion() byte
}

// EncodeScopedToken serialises a scoped token to its opaque client-facing
// form: base64 (raw URL alphabet) over the JSON encoding of the struct.
func EncodeScopedToken(token any) (string, error) {
	encoded, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(encoded), nil
}

// DecodeScopedToken reverses EncodeScopedToken into dst. Anything that is
// not the base64 JSON form of a versioned token — a token from another
// vocabulary, a legacy positional token, or plain garbage — yields
// ErrInvalidScopedToken.
func DecodeScopedToken(encoded string, dst ScopedToken) error {
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return ErrInvalidScopedToken
	}
	if err := json.Unmarshal(decoded, dst); err != nil {
		return ErrInvalidScopedToken
	}
	if dst.TokenVersion() == 0 {
		return ErrInvalidScopedToken
	}
	return nil
}
