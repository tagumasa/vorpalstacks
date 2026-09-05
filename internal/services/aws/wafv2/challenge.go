package wafv2

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	waf "vorpalstacks/internal/common/invokers/waf"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/wafv2/inspection"
	wafstore "vorpalstacks/internal/store/aws/waf"
)

const (
	// challengeIDBytes is the entropy of an issued challenge
	// identifier; the identifier is unguessable so a client cannot
	// precompute a solution before the interstitial is served.
	challengeIDBytes = 16
	// challengeTTL bounds how long an issued challenge stays
	// redeemable; the interstitial solves in milliseconds, so a window
	// of minutes only accommodates slow clients.
	challengeTTL = 5 * time.Minute
	// challengeRegistryMax bounds the outstanding-challenge registry;
	// the TTL sweep runs whenever the bound is reached.
	challengeRegistryMax = 10000
)

// outstandingChallenge is one interstitial challenge issued by an
// interrupting Captcha or Challenge response.
type outstandingChallenge struct {
	kind     string
	issuedAt time.Time
}

// challengeRegistry tracks issued challenges so the token exchange
// endpoint accepts solutions only for fresh, single-use challenges.
type challengeRegistry struct {
	mu      sync.Mutex
	entries map[string]outstandingChallenge
}

func newChallengeRegistry() *challengeRegistry {
	return &challengeRegistry{entries: map[string]outstandingChallenge{}}
}

// issue registers a fresh challenge of the given kind and returns its
// identifier. An entropy failure returns an empty identifier, and the
// caller serves the interruption without an interstitial body.
func (r *challengeRegistry) issue(kind string, now time.Time) string {
	buf := make([]byte, challengeIDBytes)
	if _, err := rand.Read(buf); err != nil {
		return ""
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.entries) >= challengeRegistryMax {
		for id, entry := range r.entries {
			if now.Sub(entry.issuedAt) > challengeTTL {
				delete(r.entries, id)
			}
		}
	}
	id := hex.EncodeToString(buf)
	r.entries[id] = outstandingChallenge{kind: kind, issuedAt: now}
	return id
}

// redeem consumes a challenge identifier and returns its kind when the
// challenge is still within its validity window; every identifier is
// single-use.
func (r *challengeRegistry) redeem(id string, now time.Time) (string, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.entries[id]
	if !ok {
		return "", false
	}
	delete(r.entries, id)
	if now.Sub(entry.issuedAt) > challengeTTL {
		return "", false
	}
	return entry.kind, true
}

// wafv2TokenStoreKeyType keys the singleton token store in the service's
// store cache; the signing secret is shared by every region because tokens
// scope to protected hosts, not regions. The named type keeps the key
// distinct from the other empty-struct keys that share the map — two
// unnamed struct{} values are the same interface key.
type wafv2TokenStoreKeyType struct{}

var wafv2TokenStoreKey wafv2TokenStoreKeyType

// tokenStoreLoad returns the singleton token store over the global
// storage.
func (s *WAFv2Service) tokenStoreLoad() (*wafstore.TokenStore, error) {
	if cached, ok := s.stores.Load(wafv2TokenStoreKey); ok {
		return cached.(*wafstore.TokenStore), nil
	}
	if s.storageManager == nil {
		return nil, invalidParamError("wafv2 storage manager not initialised")
	}
	globalStorage, err := s.storageManager.GetGlobalStorage()
	if err != nil {
		return nil, err
	}
	store := wafstore.NewTokenStore(globalStorage)
	actual, _ := s.stores.LoadOrStore(wafv2TokenStoreKey, store)
	return actual.(*wafstore.TokenStore), nil
}

// tokenValidator adapts the persistent signing secret to the
// evaluator's token-validator contract.
type serviceTokenValidator struct {
	secret []byte
}

func (v serviceTokenValidator) ValidateToken(value string) (inspection.ChallengeToken, bool) {
	return inspection.ParseToken(v.secret, value)
}

// challengeInterstitial registers a challenge of the interrupting kind
// and renders the interstitial page that solves it. An empty return
// means the challenge could not be issued and the interruption is
// served without a body.
func (s *WAFv2Service) challengeInterstitial(kind string) string {
	id := s.challenges.issue(strings.ToLower(kind), time.Now())
	if id == "" {
		return ""
	}
	return inspection.ChallengeInterstitialHTML(kind, id, waf.ChallengeEndpointPath, inspection.ChallengeSolutionHexZeros)
}

// ExchangeWAFToken verifies an interstitial solution and answers with
// the aws-waf-token cookie. The exchanged token carries the solve
// timestamp of the challenge kind the interstitial presented and
// preserves the other kind's timestamp from a token the request already
// carries — the token characteristics page describes the two timestamps
// as independent state of one token. The return value always reports
// the request as served; a failed verification is a 400, never a
// pass-through to the protected resource.
func (s *WAFv2Service) ExchangeWAFToken(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {
	var submission struct {
		ChallengeID string `json:"challengeId"`
		Counter     string `json:"counter"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&submission); err != nil {
		http.Error(w, "malformed challenge submission", http.StatusBadRequest)
		return true
	}
	kind, ok := s.challenges.redeem(submission.ChallengeID, time.Now())
	if !ok || !inspection.VerifyChallengeSolution(submission.ChallengeID, submission.Counter) {
		http.Error(w, "challenge verification failed", http.StatusBadRequest)
		return true
	}
	secret, err := s.tokenStoreLoad()
	if err == nil {
		_, err = secret.SigningKey()
	}
	if err != nil {
		logs.Warn("wafv2 token signing key unavailable, challenge exchange failed", logs.Err(err))
		http.Error(w, "token issuance unavailable", http.StatusInternalServerError)
		return true
	}
	key, _ := secret.SigningKey()
	now := time.Now()
	token := inspection.ChallengeToken{IssuedAt: now.Unix(), Domains: []string{r.Host}}
	if existing, ok := inspection.ParseToken(key, existingTokenValue(r)); ok {
		token.ChallengeSolvedAt = existing.ChallengeSolvedAt
		token.CaptchaSolvedAt = existing.CaptchaSolvedAt
	}
	if kind == strings.ToLower(inspection.ActionCaptcha) {
		token.CaptchaSolvedAt = now.Unix()
	} else {
		token.ChallengeSolvedAt = now.Unix()
	}
	http.SetCookie(w, &http.Cookie{
		Name:     inspection.TokenCookieName,
		Value:    inspection.SignToken(key, token),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"token":true}`))
	return true
}

// existingTokenValue reads the aws-waf-token cookie from the exchange
// request, if the client already carries one.
func existingTokenValue(r *http.Request) string {
	cookie, err := r.Cookie(inspection.TokenCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}
