package eventbridge

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// apiDestinationClientTimeout bounds every request to an API destination
// endpoint: "EventBridge requests to an API destination endpoint must have
// a maximum client execution timeout of 5 seconds. If the target endpoint
// takes longer than 5 seconds to respond, EventBridge times out the
// request" (user guide, API destinations).
const apiDestinationClientTimeout = 5 * time.Second

// apiDestinationUserAgent is the fixed User-Agent of API destination
// requests; AWS documents it in the "Headers included in requests to API
// destinations" table, and it cannot be overridden by connection or target
// header parameters.
const apiDestinationUserAgent = "Amazon/EventBridge/ApiDestinations"

// errRetryAfterHint is implemented by API destination delivery errors that
// carry the endpoint's Retry-After delay: "EventBridge API destinations
// read the standard HTTP response header Retry-After to find out how long
// to wait before making a follow-up request ... EventBridge chooses the
// more conservative value between the defined retry policy and the
// Retry-After header" (user guide). A negative Retry-After — "EventBridge
// stops retrying delivery for that event" — is a permanent failure.
type errRetryAfterHint interface {
	error
	RetryAfterDelay() time.Duration
}

// apiDestinationError is one API destination invocation outcome: the HTTP
// status (0 for transport failures), whether the endpoint's Retry-After
// header asked the delivery to stop, and the retry delay it requested.
type apiDestinationError struct {
	targetARN  string
	statusCode int
	retryAfter time.Duration
	stopRetry  bool
	cause      error
}

func (e *apiDestinationError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("api destination %s invocation failed: %v", e.targetARN, e.cause)
	}
	return fmt.Sprintf("api destination %s returned status %d", e.targetARN, e.statusCode)
}

func (e *apiDestinationError) Unwrap() error { return e.cause }

// RetryAfterDelay reports the endpoint's requested follow-up delay.
func (e *apiDestinationError) RetryAfterDelay() time.Duration { return e.retryAfter }

// classifyAPIDestinationStatus maps an HTTP status to the retry contract:
// "Retries events associated with error codes 401, 407, 409, 429, and 5xx.
// Does not retry events associated with error codes 1xx, 2xx, 3xx, and 4xx
// (other than those noted above)" (user guide). 2xx answers nil.
func classifyAPIDestinationStatus(e *apiDestinationError) error {
	switch {
	case e.statusCode >= 200 && e.statusCode < 300:
		return nil
	case e.statusCode == 401, e.statusCode == 407, e.statusCode == 409,
		e.statusCode == 429, e.statusCode >= 500:
		return e
	default:
		return fmt.Errorf("%w: %v", errPermanentDelivery, e)
	}
}

// apiDestinationPacer enforces an API destination's invocation rate
// ("The maximum number of requests per second to send to the HTTP
// invocation endpoint"): each invocation reserves the next slot of a
// strictly spaced timeline, so bursts queue instead of exceeding the rate.
type apiDestinationPacer struct {
	mu       sync.Mutex
	interval time.Duration
	next     time.Time
}

// wait reserves the next rate slot, blocking until it arrives. It returns
// the context error when the delivery is cancelled while queued.
func (p *apiDestinationPacer) wait(ctx context.Context) error {
	p.mu.Lock()
	now := time.Now()
	if p.next.Before(now) {
		p.next = now
	}
	slot := p.next
	p.next = p.next.Add(p.interval)
	p.mu.Unlock()

	delay := time.Until(slot)
	if delay <= 0 {
		return nil
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

// oauthToken is a cached OAuth access token with its expiry.
type oauthToken struct {
	accessToken string
	tokenType   string
	expiresAt   time.Time
}

// oauthRefreshLead is how long before expiry a token is proactively
// refreshed on the invocation path: "Proactively during an HTTPS
// invocation, if the token expires within 60 seconds" (user guide).
const oauthRefreshLead = 60 * time.Second

// deliverToApiDestination invokes an API destination target's HTTPS
// endpoint with the event payload, authorised by the destination's
// connection. The connection's credentials are read from its credential
// secret (the record keeps only the non-credential halves).
func (s *EventsService) deliverToApiDestination(ctx context.Context, job *deliveryJob) error {
	targetARN := job.target.ARN
	destName := svcarn.ExtractApiDestinationNameFromARN(targetARN)
	if destName == "" {
		return permanentDeliveryError("api destination target ARN %s does not name an api-destination resource", targetARN)
	}
	_, _, arnRegion, _, _ := svcarn.SplitARN(targetARN)
	if arnRegion == "" {
		arnRegion = job.region
	}
	store, err := s.GetStoreForRegion(arnRegion)
	if err != nil {
		return fmt.Errorf("failed to get store for region %s: %w", arnRegion, err)
	}
	dest, err := store.GetApiDestination(ctx, destName)
	if err != nil {
		return permanentDeliveryError("api destination %s does not exist", destName)
	}

	connection, err := store.GetConnection(ctx, svcarn.ExtractConnectionNameFromARN(dest.ConnectionARN))
	if err != nil {
		return permanentDeliveryError("connection for api destination %s does not exist (%s)", destName, dest.ConnectionARN)
	}
	if connection.State == eventsstore.ConnectionStateDeauthorized || connection.AuthParameters == nil {
		return permanentDeliveryError("connection %s is deauthorized; api destination invocations have no credentials", connection.Name)
	}
	credentials, err := s.connectionCredentials(ctx, store, connection)
	if err != nil {
		return fmt.Errorf("load credentials for connection %s: %w", connection.Name, err)
	}

	if dest.InvocationRateLimitPerSecond > 0 {
		if err := s.paceApiDestination(ctx, dest.Name, dest.InvocationRateLimitPerSecond); err != nil {
			return err
		}
	}

	err = s.invokeApiDestination(ctx, job, connection, credentials, dest)
	if err == nil {
		return nil
	}
	// A 401/407 from an OAuth connection refreshes the token and retries
	// the invocation once ("EventBridge refreshes OAuth tokens: When a 401
	// or 407 response is returned", user guide).
	var apiErr *apiDestinationError
	if errors.As(err, &apiErr) &&
		(apiErr.statusCode == http.StatusUnauthorized || apiErr.statusCode == http.StatusProxyAuthRequired) &&
		connection.AuthorizationType == "OAUTH_CLIENT_CREDENTIALS" {
		s.dropOAuthToken(connection.ARN)
		return s.invokeApiDestination(ctx, job, connection, credentials, dest)
	}
	return err
}

// invokeApiDestination performs one authorised HTTP request to the
// destination's endpoint and classifies the outcome.
func (s *EventsService) invokeApiDestination(
	ctx context.Context,
	job *deliveryJob,
	connection *eventsstore.Connection,
	credentials *eventsstore.AuthParameters,
	dest *eventsstore.ApiDestination,
) error {
	endpoint, headers, body, err := s.buildApiDestinationRequest(ctx, job, connection, credentials, dest)
	if err != nil {
		return fmt.Errorf("build api destination request: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, apiDestinationClientTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, dest.HttpMethod, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build api destination request for %s: %w", endpoint, err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// "Connection: close" is part of the documented fixed header set.
	req.Close = true

	resp, err := s.apiDestHTTP.Do(req)
	if err != nil {
		// Timeouts and transport faults are retried like a 5xx: the
		// endpoint may recover.
		return &apiDestinationError{targetARN: job.target.ARN, cause: err}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(resp.Body, 1<<16))

	apiErr := &apiDestinationError{targetARN: job.target.ARN, statusCode: resp.StatusCode}
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		if seconds, parseErr := strconv.ParseFloat(retryAfter, 64); parseErr == nil {
			if seconds < 0 {
				apiErr.stopRetry = true
			} else {
				apiErr.retryAfter = time.Duration(seconds * float64(time.Second))
			}
		}
	}
	if apiErr.stopRetry {
		return fmt.Errorf("%w: endpoint %s sent a negative Retry-After", errPermanentDelivery, endpoint)
	}
	return classifyAPIDestinationStatus(apiErr)
}

// buildApiDestinationRequest assembles the invocation URL, headers and
// body: the fixed header set, the connection's authorisation, the
// connection's invocation parameters merged with the target's HTTP
// parameters (the connection wins conflicts — "In case of any conflicting
// keys, values from the Connection take precedence", Smithy model), and
// the target's path values filling the endpoint's "*" wildcards.
func (s *EventsService) buildApiDestinationRequest(
	ctx context.Context,
	job *deliveryJob,
	connection *eventsstore.Connection,
	credentials *eventsstore.AuthParameters,
	dest *eventsstore.ApiDestination,
) (string, map[string]string, []byte, error) {
	headers := map[string]string{
		"User-Agent":      apiDestinationUserAgent,
		"Content-Type":    "application/json; charset=utf-8",
		"Range":           "bytes=0-1048575",
		"Accept-Encoding": "gzip,deflate",
	}

	body := job.payload
	var connectionParams *eventsstore.ConnectionHttpParameters
	if credentials != nil {
		connectionParams = credentials.InvocationHttpParameters
	}
	if connectionParams != nil {
		for _, h := range connectionParams.HeaderParameters {
			headers[h.Key] = h.Value
		}
	}
	if job.target.HttpParameters != nil {
		for key, value := range job.target.HttpParameters.HeaderParameters {
			if _, ok := headers[key]; !ok {
				headers[key] = value
			}
		}
	}

	// Authorisation headers per the connection's type.
	authHeader, err := s.apiDestinationAuthHeader(ctx, connection, credentials)
	if err != nil {
		return "", nil, nil, err
	}
	if authHeader != "" {
		headers["Authorization"] = authHeader
	}
	if credentials != nil && credentials.ApiKeyAuthParameters != nil {
		headers[credentials.ApiKeyAuthParameters.ApiKeyName] = credentials.ApiKeyAuthParameters.ApiKeyValue
	}

	// Body parameters merge into a JSON-object payload; the payload of an
	// input transformer is the event document the endpoint consumes.
	if connectionParams != nil && len(connectionParams.BodyParameters) > 0 {
		var object map[string]interface{}
		if json.Unmarshal(body, &object) == nil && object != nil {
			for _, b := range connectionParams.BodyParameters {
				object[b.Key] = b.Value
			}
			if merged, err := json.Marshal(object); err == nil {
				body = merged
			}
		}
	}

	endpoint, err := buildApiDestinationURL(dest.InvocationEndpoint, connectionParams, job.target)
	if err != nil {
		return "", nil, nil, err
	}
	return endpoint, headers, body, nil
}

// buildApiDestinationURL applies the target's path values to the endpoint
// wildcards and merges the connection's and target's query parameters
// (the connection wins conflicts).
func buildApiDestinationURL(endpoint string, connectionParams *eventsstore.ConnectionHttpParameters, target eventsstore.Target) (string, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse invocation endpoint %s: %w", endpoint, err)
	}

	if target.HttpParameters != nil && len(target.HttpParameters.PathParameterValues) > 0 {
		path := parsed.Path
		for _, value := range target.HttpParameters.PathParameterValues {
			if !strings.Contains(path, "*") {
				break
			}
			path = strings.Replace(path, "*", value, 1)
		}
		parsed.Path = path
	}

	query := parsed.Query()
	if target.HttpParameters != nil {
		for key, value := range target.HttpParameters.QueryStringParameters {
			if _, conflict := query[key]; !conflict {
				query.Set(key, value)
			}
		}
	}
	if connectionParams != nil {
		// Connection values take precedence over target values.
		for _, q := range connectionParams.QueryStringParameters {
			query.Set(q.Key, q.Value)
		}
	}
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

// apiDestinationAuthHeader resolves the authorisation an invocation
// carries: Basic credentials encode into the Authorization header, an API
// key connection returns no Authorization (the key rides its own header),
// and an OAuth connection exchanges its client credentials for a managed
// access token.
func (s *EventsService) apiDestinationAuthHeader(ctx context.Context, connection *eventsstore.Connection, credentials *eventsstore.AuthParameters) (string, error) {
	switch connection.AuthorizationType {
	case "BASIC":
		if credentials == nil || credentials.BasicAuthParameters == nil {
			return "", nil
		}
		basic := credentials.BasicAuthParameters.Username + ":" + credentials.BasicAuthParameters.Password
		return "Basic " + base64.StdEncoding.EncodeToString([]byte(basic)), nil
	case "OAUTH_CLIENT_CREDENTIALS":
		token, err := s.oauthAccessToken(ctx, connection, credentials)
		if err != nil {
			return "", err
		}
		return token.tokenType + " " + token.accessToken, nil
	default:
		return "", nil
	}
}

// connectionCredentials loads the full authorization parameters of a
// connection from its credential secret.
func (s *EventsService) connectionCredentials(ctx context.Context, store *eventsstore.EventsStore, connection *eventsstore.Connection) (*eventsstore.AuthParameters, error) {
	if connection.SecretArn == "" {
		return nil, fmt.Errorf("connection %s has no credential secret", connection.Name)
	}
	invoker := s.secretsInvoker()
	if invoker == nil {
		return nil, errSecretsManagerUnavailable()
	}
	payload, err := invoker.GetServiceSecretString(ctx, store.Region(), connection.SecretArn)
	if err != nil {
		return nil, err
	}
	var params eventsstore.AuthParameters
	if err := json.Unmarshal([]byte(payload), &params); err != nil {
		return nil, fmt.Errorf("decode credential secret: %w", err)
	}
	return &params, nil
}

// paceApiDestination waits for the destination's next rate slot. The pacer
// is created on first use and replaced whenever the record's rate changes,
// so a rate update takes effect on the next invocation.
func (s *EventsService) paceApiDestination(ctx context.Context, name string, ratePerSecond int32) error {
	interval := time.Second / time.Duration(ratePerSecond)
	s.apiDestPacersMu.Lock()
	pacer, ok := s.apiDestPacers[name]
	if !ok || pacer.interval != interval {
		pacer = &apiDestinationPacer{interval: interval}
		s.apiDestPacers[name] = pacer
	}
	s.apiDestPacersMu.Unlock()
	return pacer.wait(ctx)
}

// newAPIDestinationHTTPClient builds the shared API destination HTTP client
// at service construction, so no delivery path ever writes the field. In
// TEST_MODE the client accepts otherwise-untrusted certificates so the test
// suite can point destinations at locally served TLS endpoints; the gate
// server runs the suite with TEST_MODE=true, which is constant for the
// process lifetime and therefore readable once here.
func newAPIDestinationHTTPClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if os.Getenv("TEST_MODE") == "true" {
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
	return &http.Client{Transport: transport}
}

// oauthAccessToken returns a usable OAuth access token for the connection:
// a cached token that is not about to expire, or a fresh client-credentials
// exchange ("For OAuth authorization, EventBridge also exchanges your
// client ID and secret for an access token and then manages it securely",
// user guide). AWS does not document the token request's wire format, so
// the exchange follows RFC 6749 §4.4 (grant_type=client_credentials with
// the client credentials in the form body); the connection's
// OAuthHttpParameters ride the token request in all three families —
// headers on the request, query-string members appended to the
// authorization endpoint, and body members joining the form.
func (s *EventsService) oauthAccessToken(ctx context.Context, connection *eventsstore.Connection, credentials *eventsstore.AuthParameters) (*oauthToken, error) {
	if credentials == nil || credentials.OAuthParameters == nil {
		return nil, fmt.Errorf("connection %s carries no OAuth parameters", connection.Name)
	}
	oauth := credentials.OAuthParameters

	s.apiDestTokensMu.Lock()
	cached, ok := s.apiDestTokens[connection.ARN]
	s.apiDestTokensMu.Unlock()
	if ok && (cached.expiresAt.IsZero() || time.Now().Add(oauthRefreshLead).Before(cached.expiresAt)) {
		return cached, nil
	}

	token, err := s.exchangeOAuthToken(ctx, connection.ARN, oauth)
	if err != nil {
		return nil, err
	}
	s.apiDestTokensMu.Lock()
	s.apiDestTokens[connection.ARN] = token
	s.apiDestTokensMu.Unlock()
	return token, nil
}

// dropOAuthToken discards the connection's cached token (a 401/407 forces
// the next invocation into a synchronous refresh).
func (s *EventsService) dropOAuthToken(connectionARN string) {
	s.apiDestTokensMu.Lock()
	delete(s.apiDestTokens, connectionARN)
	s.apiDestTokensMu.Unlock()
}

// exchangeOAuthToken performs the client-credentials token exchange
// against the connection's authorization endpoint.
func (s *EventsService) exchangeOAuthToken(ctx context.Context, connectionARN string, oauth *eventsstore.OAuthParameters) (*oauthToken, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {oauth.ClientParameters.ClientID},
		"client_secret": {oauth.ClientParameters.ClientSecret},
	}
	endpoint := oauth.AuthorizationEndpoint
	if params := oauth.OAuthHttpParameters; params != nil {
		// Body members join the client-credentials form; query-string
		// members are appended to the authorization endpoint's query.
		// Body keys carry no modelled pattern, so an empty key is skipped
		// rather than emitted as a valueless form field.
		for _, b := range params.BodyParameters {
			if b.Key != "" {
				form.Set(b.Key, b.Value)
			}
		}
		if len(params.QueryStringParameters) > 0 {
			endpointURL, perr := url.Parse(endpoint)
			if perr != nil {
				return nil, fmt.Errorf("build OAuth token request: parse authorization endpoint: %w", perr)
			}
			query := endpointURL.Query()
			for _, q := range params.QueryStringParameters {
				query.Add(q.Key, q.Value)
			}
			endpointURL.RawQuery = query.Encode()
			endpoint = endpointURL.String()
		}
	}

	reqCtx, cancel := context.WithTimeout(ctx, apiDestinationClientTimeout)
	defer cancel()
	method := oauth.HttpMethod
	if method == "" {
		method = http.MethodPost
	}
	req, err := http.NewRequestWithContext(reqCtx, method, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build OAuth token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if oauth.OAuthHttpParameters != nil {
		for _, h := range oauth.OAuthHttpParameters.HeaderParameters {
			req.Header.Set(h.Key, h.Value)
		}
	}
	req.Close = true

	resp, err := s.apiDestHTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OAuth token exchange for connection %s: %w", connectionARN, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read OAuth token response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("OAuth token exchange for connection %s returned status %d", connectionARN, resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, fmt.Errorf("decode OAuth token response: %w", err)
	}
	if tokenResponse.AccessToken == "" {
		return nil, fmt.Errorf("OAuth token response for connection %s carries no access_token", connectionARN)
	}
	tokenType := tokenResponse.TokenType
	if tokenType == "" {
		tokenType = "Bearer"
	}
	token := &oauthToken{
		accessToken: tokenResponse.AccessToken,
		tokenType:   tokenType,
	}
	if tokenResponse.ExpiresIn > 0 {
		token.expiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second)
	}
	logs.Debug("EventBridge: exchanged OAuth token for API destination connection",
		logs.String("connectionArn", connectionARN))
	return token, nil
}
