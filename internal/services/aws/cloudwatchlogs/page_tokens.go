package cloudwatchlogs

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"vorpalstacks/internal/common/pagination"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The service-plane scoped page tokens: stream listings and query
// listings page through one vocabulary each, carrying the request
// identity (group, filters, ordering) plus a value cursor, so a token
// replayed against a different request is rejected instead of silently
// repositioning the walk. The rejection is two-layered: the version
// namespace below keeps a token of one vocabulary foreign to every
// other vocabulary's decode, and within a vocabulary the
// request-identity members must match.

// The page-token version bytes are one namespace: the wire form is
// base64-JSON whose decode ignores unknown keys, so a token minted by
// one vocabulary structurally decodes into every other and the version
// byte is the only cross-vocabulary discriminator — two vocabularies
// never share a value, and the shared decode below rejects a foreign
// version as an invalid token instead of repositioning that walk.
const (
	streamPageTokenVersion       byte = 1
	queryPageTokenVersion        byte = 2
	queryResultsPageTokenVersion byte = 3
	// listingPageTokenVersion versions the generic name-cursor listing
	// vocabulary (scopedListingToken below).
	listingPageTokenVersion byte = 4
)

// scopedPageToken is the contract each page vocabulary's token fulfils
// for the shared codec pair below: stamp sets the two fields every
// encode must carry (the vocabulary version and the mint clock — the
// generic cannot reach the concrete fields), and minted reads the mint
// back for the expiry check.
type scopedPageToken[T any] interface {
	*T
	pagination.ScopedToken
	stamp(version byte)
	minted() int64
}

// encodeScopedPageToken stamps the vocabulary's version and mint onto
// the token and seals it — one body behind every page vocabulary's
// encoder, so no vocabulary can mint an unstamped or mis-stamped token.
func encodeScopedPageToken[T any, PT scopedPageToken[T]](token PT, version byte) (string, error) {
	token.stamp(version)
	return pagination.EncodeScopedToken(*token)
}

// decodeScopedPageToken unwraps one vocabulary's page token: the wire
// form must decode, carry the vocabulary's version, and sit inside the
// vocabulary's TTL — the shared half of every decode; member checks
// beyond that discipline stay at the vocabulary's own wrapper.
func decodeScopedPageToken[T any, PT scopedPageToken[T]](s string, version byte, ttl time.Duration) (T, error) {
	var token T
	ptr := PT(&token)
	if err := pagination.DecodeScopedToken(s, ptr); err != nil || ptr.TokenVersion() != version ||
		pageTokenExpired(ptr.minted(), ttl) {
		return token, errInvalidPageToken
	}
	return token, nil
}

// errInvalidPageToken maps a foreign or stale token to the operations'
// documented invalid-parameter error.
var errInvalidPageToken = NewLogsError("InvalidParameterException", "Invalid nextToken", 400)

// The documented token expiries: GetQueryResults' nextToken "expires
// after 1 hour" while the describe listings' tokens "expire after 24
// hours" — the shared 24-hour value the store plane defines
// (logsstore.TokenExpiry24h) so both planes reference one definition. A
// token past its expiry rejects as invalid rather than repositioning the
// walk.
const (
	resultsPageTokenTTL = time.Hour
	listPageTokenTTL    = logsstore.TokenExpiry24h
)

// pageTokenNowMs is the clock the tokens mint from; tests backdate a
// mint by overriding it.
var pageTokenNowMs = func() int64 { return time.Now().UnixMilli() }

// pageTokenExpired reports whether a mint timestamp has passed its
// listing's TTL; a zero mint (a token from a vocabulary predating the
// field) is always expired.
func pageTokenExpired(mintedMs int64, ttl time.Duration) bool {
	if mintedMs <= 0 {
		return true
	}
	return time.Now().UnixMilli()-mintedMs > ttl.Milliseconds()
}

// streamPageToken scopes a DescribeLogStreams page. The cursor is the
// last served stream's name (name ordering) or its (LastEventTime, name)
// pair (event-time ordering); both orderings qualify as total orders
// because stream names are unique within a group.
type streamPageToken struct {
	Version     byte   `json:"v"`
	Group       string `json:"g"`
	Prefix      string `json:"p,omitempty"`
	OrderBy     string `json:"o,omitempty"`
	Descending  bool   `json:"de,omitempty"`
	LastName    string `json:"n,omitempty"`
	LastEventTs int64  `json:"t,omitempty"`
	MintedMs    int64  `json:"m"`
}

// TokenVersion implements pagination.ScopedToken.
func (t streamPageToken) TokenVersion() byte { return t.Version }

// stamp sets the shared codec's two encode fields; minted reads the
// mint its expiry check ages from.
func (t *streamPageToken) stamp(version byte) {
	t.Version = version
	t.MintedMs = pageTokenNowMs()
}

func (t *streamPageToken) minted() int64 { return t.MintedMs }

func encodeStreamPageToken(t streamPageToken) (string, error) {
	return encodeScopedPageToken(&t, streamPageTokenVersion)
}

func decodeStreamPageToken(s string) (streamPageToken, error) {
	return decodeScopedPageToken[streamPageToken](s, streamPageTokenVersion, listPageTokenTTL)
}

// queryPageToken scopes a DescribeQueries page. Queries list newest
// first; the cursor is the (creation time, query id) pair of the last
// served query.
type queryPageToken struct {
	Version        byte   `json:"v"`
	StatusFilter   string `json:"sf,omitempty"`
	LogGroupName   string `json:"lg,omitempty"`
	QueryLanguage  string `json:"ql,omitempty"`
	CursorCreateNs int64  `json:"c"`
	CursorQueryId  string `json:"q,omitempty"`
	MintedMs       int64  `json:"m"`
}

// TokenVersion implements pagination.ScopedToken.
func (t queryPageToken) TokenVersion() byte { return t.Version }

// stamp sets the shared codec's two encode fields; minted reads the
// mint its expiry check ages from.
func (t *queryPageToken) stamp(version byte) {
	t.Version = version
	t.MintedMs = pageTokenNowMs()
}

func (t *queryPageToken) minted() int64 { return t.MintedMs }

func encodeQueryPageToken(t queryPageToken) (string, error) {
	return encodeScopedPageToken(&t, queryPageTokenVersion)
}

func decodeQueryPageToken(s string) (queryPageToken, error) {
	return decodeScopedPageToken[queryPageToken](s, queryPageTokenVersion, listPageTokenTTL)
}

// queryResultsPageToken scopes a GetQueryResults page to the query whose
// results it pages ("The token expires after 1 hour"): the cursor is the
// offset of the next unserved row, valid only against the minting query.
type queryResultsPageToken struct {
	Version  byte   `json:"v"`
	QueryId  string `json:"q"`
	Offset   int    `json:"o"`
	MintedMs int64  `json:"m"`
}

// TokenVersion implements pagination.ScopedToken.
func (t queryResultsPageToken) TokenVersion() byte { return t.Version }

// stamp sets the shared codec's two encode fields; minted reads the
// mint its expiry check ages from.
func (t *queryResultsPageToken) stamp(version byte) {
	t.Version = version
	t.MintedMs = pageTokenNowMs()
}

func (t *queryResultsPageToken) minted() int64 { return t.MintedMs }

func encodeQueryResultsPageToken(t queryResultsPageToken) (string, error) {
	return encodeScopedPageToken(&t, queryResultsPageTokenVersion)
}

func decodeQueryResultsPageToken(s string) (queryResultsPageToken, error) {
	t, err := decodeScopedPageToken[queryResultsPageToken](s, queryResultsPageTokenVersion, resultsPageTokenTTL)
	if err != nil || t.Offset < 0 {
		return t, errInvalidPageToken
	}
	return t, nil
}

// scopedListingToken scopes the listings whose previous tokens were bare
// names or raw store keys (log groups, metric filters, resource policies,
// subscription filters, field indexes, a query's log groups): the token
// carries a digest of the request identity plus the opaque value cursor,
// so a token minted by one request repositions no other, and a typed or
// stale token rejects instead of silently ending the walk.
type scopedListingToken struct {
	Version  byte   `json:"v"`
	Scope    string `json:"s"`
	Cursor   string `json:"c,omitempty"`
	MintedMs int64  `json:"m"`
}

// TokenVersion implements pagination.ScopedToken.
func (t scopedListingToken) TokenVersion() byte { return t.Version }

// stamp sets the shared codec's two encode fields; minted reads the
// mint its expiry check ages from.
func (t *scopedListingToken) stamp(version byte) {
	t.Version = version
	t.MintedMs = pageTokenNowMs()
}

func (t *scopedListingToken) minted() int64 { return t.MintedMs }

func encodeListingToken(scope, cursor string) (string, error) {
	return encodeScopedPageToken(&scopedListingToken{Scope: scope, Cursor: cursor}, listingPageTokenVersion)
}

// decodeListingToken returns the cursor a token carries, rejecting a
// token minted for any other scope (or one that is not a token at all).
func decodeListingToken(token, scope string) (string, error) {
	t, err := decodeScopedPageToken[scopedListingToken](token, listingPageTokenVersion, listPageTokenTTL)
	if err != nil || t.Scope != scope {
		return "", errInvalidPageToken
	}
	return t.Cursor, nil
}

// listingScope digests a listing's request identity into the opaque
// scope its tokens carry: the parts join with a separator outside the
// members' charset so no two identities digest alike.
func listingScope(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return hex.EncodeToString(sum[:12])
}

// paginateScopedListing runs the scoped-listing round-trip every
// paginated CWL listing shares: the incoming token unwraps under the
// scope's vocabulary, the listing pages by its natural key, and the
// outgoing page marker wraps under the same vocabulary — one body, so
// the marker always comes from the pagination result and never from
// per-site knowledge about where the walk stopped. The scope arrives
// pre-digested (listingScope): the request-identity parts stay the
// caller's knowledge.
func paginateScopedListing[T any](scope, nextToken string, items []T, limit int, keyOf func(T) string) (pagination.SliceResult[T], error) {
	cursor := ""
	if nextToken != "" {
		var err error
		cursor, err = decodeListingToken(nextToken, scope)
		if err != nil {
			return pagination.SliceResult[T]{}, err
		}
	}
	result := pagination.PaginateSlice(items, cursor, limit, keyOf)
	if err := rewrapResultMarker(scope, &result); err != nil {
		return pagination.SliceResult[T]{}, err
	}
	return result, nil
}

// paginateScopedListingByPosition is the positional-cursor twin of
// paginateScopedListing for the listings whose page marker is a
// position rather than a natural key (a batch index inside one import
// task), where no unique key exists to walk.
func paginateScopedListingByPosition[T any](scope, nextToken string, items []T, limit int) (pagination.SliceResult[T], error) {
	cursor := ""
	if nextToken != "" {
		var err error
		cursor, err = decodeListingToken(nextToken, scope)
		if err != nil {
			return pagination.SliceResult[T]{}, err
		}
	}
	result := pagination.PaginateSliceByPosition(items, cursor, limit)
	if err := rewrapResultMarker(scope, &result); err != nil {
		return pagination.SliceResult[T]{}, err
	}
	return result, nil
}

// wrapScopedListingMarker wraps a marker a lower plane produced (the
// store's own walk, a hand-filtered selection) in the scoped listing
// vocabulary — the wrap half of paginateScopedListing for the listings
// whose pagination happens below this plane. An empty marker stays
// empty: no token is minted for an exhausted walk.
func wrapScopedListingMarker(scope, marker string) (string, error) {
	if marker == "" {
		return "", nil
	}
	return encodeListingToken(scope, marker)
}

// rewrapResultMarker replaces a page result's bare marker with its
// scoped-token form in place.
func rewrapResultMarker[T any](scope string, result *pagination.SliceResult[T]) error {
	if result.NextMarker == "" {
		return nil
	}
	wrapped, err := encodeListingToken(scope, result.NextMarker)
	if err != nil {
		return err
	}
	result.NextMarker = wrapped
	return nil
}
