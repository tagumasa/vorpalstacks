package eventbridge

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	eventsstore "vorpalstacks/internal/store/aws/eventbridge"
)

// recordedRequest is one invocation captured by the test endpoint.
type recordedRequest struct {
	method  string
	path    string
	query   string
	headers http.Header
	body    string
}

// apiDestinationTestServer is a locally served invocation endpoint that
// records every request and answers from a scripted status sequence.
type apiDestinationTestServer struct {
	mu       sync.Mutex
	requests []recordedRequest
	statuses []int
	headers  map[string]string
	server   *httptest.Server
}

func newAPIDestinationTestServer(statuses ...int) *apiDestinationTestServer {
	s := &apiDestinationTestServer{statuses: statuses}
	s.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		s.mu.Lock()
		s.requests = append(s.requests, recordedRequest{
			method:  r.Method,
			path:    r.URL.Path,
			query:   r.URL.RawQuery,
			headers: r.Header.Clone(),
			body:    string(body),
		})
		status := http.StatusOK
		if len(s.statuses) > 0 {
			status = s.statuses[0]
			s.statuses = s.statuses[1:]
		}
		for key, value := range s.headers {
			w.Header().Set(key, value)
		}
		s.mu.Unlock()
		w.WriteHeader(status)
	}))
	return s
}

func (s *apiDestinationTestServer) url() string { return s.server.URL }

func (s *apiDestinationTestServer) close() { s.server.Close() }

func (s *apiDestinationTestServer) recorded() []recordedRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]recordedRequest(nil), s.requests...)
}

// seedApiDestinationFixture registers a service+store pair (fake secrets
// invoker wired) and seeds a BASIC connection whose credentials live in
// the fake secret store.
func seedApiDestinationFixture(t *testing.T) (*EventsService, *eventsstore.EventsStore, *fakeSecretsInvoker) {
	t.Helper()
	svc, store, fake := newConnectionTestService(t)
	svc.SetEventsStore("us-east-1", store)
	ctx := context.Background()
	if _, err := svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "dest-conn",
		AuthorizationType: "BASIC",
		AuthParameters: &eventsstore.AuthParameters{
			BasicAuthParameters: &eventsstore.BasicAuthParameters{Username: "alice", Password: "wonderland"},
			InvocationHttpParameters: &eventsstore.ConnectionHttpParameters{
				HeaderParameters: []eventsstore.ConnectionHeaderParameter{
					{Key: "X-Conn-Header", Value: "conn-value", IsValueSecret: false},
					{Key: "X-Conn-Secret", Value: "conn-secret", IsValueSecret: true},
				},
				QueryStringParameters: []eventsstore.ConnectionQueryStringParameter{
					{Key: "src", Value: "connection"},
				},
				BodyParameters: []eventsstore.ConnectionBodyParameter{
					{Key: "injectedBy", Value: "connection"},
				},
			},
		},
	}); err != nil {
		t.Fatal(err)
	}
	return svc, store, fake
}

// seedDestinationRecord writes an API destination record directly with a
// local endpoint: the delivery path exercises the invocation contract, the
// create-path endpoint validation is pinned separately.
func seedDestinationRecord(t *testing.T, store *eventsstore.EventsStore, connectionARN, endpoint, method string, rate int32) *eventsstore.ApiDestination {
	t.Helper()
	dest := &eventsstore.ApiDestination{
		Name:                         "dest-1",
		ConnectionARN:                connectionARN,
		HttpMethod:                   method,
		InvocationEndpoint:           endpoint,
		InvocationRateLimitPerSecond: rate,
	}
	if err := store.CreateApiDestination(context.Background(), dest); err != nil {
		t.Fatal(err)
	}
	return dest
}

func apiDestJob(dest *eventsstore.ApiDestination, target eventsstore.Target) *deliveryJob {
	return &deliveryJob{
		region:  "us-east-1",
		event:   &eventsstore.Event{ID: "evt-1"},
		target:  target,
		payload: []byte(`{"orderId":7}`),
	}
}

// TestApiDestinationInvocationContract pins the documented invocation
// contract end to end: the destination's method and endpoint, the target's
// path values filling the endpoint wildcard, the query merge with the
// connection's parameters taking precedence, the fixed header set, the
// Basic authorisation header, the connection's invocation headers, and the
// body-parameter merge into the JSON payload.
func TestApiDestinationInvocationContract(t *testing.T) {
	ctx := context.Background()
	svc, store, _ := seedApiDestinationFixture(t)
	endpoint := newAPIDestinationTestServer(http.StatusOK)
	defer endpoint.close()

	connection, err := store.GetConnection(ctx, "dest-conn")
	if err != nil {
		t.Fatal(err)
	}
	dest := seedDestinationRecord(t, store, connection.ARN, endpoint.url()+"/widgets/*", "PUT", 0)

	target := eventsstore.Target{
		ARN: dest.ARN,
		HttpParameters: &eventsstore.HttpParameters{
			PathParameterValues: []string{"order-7"},
			HeaderParameters:    map[string]string{"X-Target-Header": "target-value", "X-Conn-Header": "target-loses"},
			QueryStringParameters: map[string]string{
				"src":   "target-loses",
				"extra": "target-wins",
			},
		},
	}
	if err := svc.deliverToApiDestination(ctx, apiDestJob(dest, target)); err != nil {
		t.Fatalf("invocation must succeed: %v", err)
	}

	recorded := endpoint.recorded()
	if len(recorded) != 1 {
		t.Fatalf("exactly one invocation expected, got %d", len(recorded))
	}
	req := recorded[0]
	if req.method != "PUT" {
		t.Fatalf("method = %s, want PUT", req.method)
	}
	if req.path != "/widgets/order-7" {
		t.Fatalf("path wildcard must be filled by the target's value, got %s", req.path)
	}
	if !strings.Contains(req.query, "src=connection") || !strings.Contains(req.query, "extra=target-wins") {
		t.Fatalf("query merge: connection values take precedence, got %q", req.query)
	}

	expectedBasic := "Basic " + base64.StdEncoding.EncodeToString([]byte("alice:wonderland"))
	if got := req.headers.Get("Authorization"); got != expectedBasic {
		t.Fatalf("Authorization = %q, want %q", got, expectedBasic)
	}
	if got := req.headers.Get("X-Conn-Header"); got != "conn-value" {
		t.Fatalf("connection header must take precedence, got %q", got)
	}
	if got := req.headers.Get("X-Conn-Secret"); got != "conn-secret" {
		t.Fatalf("secret-marked connection header value must be sent, got %q", got)
	}
	if got := req.headers.Get("X-Target-Header"); got != "target-value" {
		t.Fatalf("target header must be sent, got %q", got)
	}
	if got := req.headers.Get("User-Agent"); got != apiDestinationUserAgent {
		t.Fatalf("User-Agent = %q", got)
	}
	if got := req.headers.Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("Content-Type = %q", got)
	}
	if got := req.headers.Get("Range"); got != "bytes=0-1048575" {
		t.Fatalf("Range = %q", got)
	}

	var body map[string]interface{}
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatalf("body must stay JSON: %v (%s)", err, req.body)
	}
	if body["orderId"] != float64(7) || body["injectedBy"] != "connection" {
		t.Fatalf("connection body parameters merge into the payload, got %s", req.body)
	}
}

// TestApiDestinationStatusClassification pins the retry classification:
// 2xx succeeds; 401/407/409/429/5xx are retried; 1xx/3xx and the other 4xx
// are permanent.
func TestApiDestinationStatusClassification(t *testing.T) {
	cases := []struct {
		status    int
		permanent bool
	}{
		{http.StatusOK, false},
		{http.StatusCreated, false},
		{http.StatusUnauthorized, false},
		{http.StatusProxyAuthRequired, false},
		{http.StatusConflict, false},
		{http.StatusTooManyRequests, false},
		{http.StatusInternalServerError, false},
		{http.StatusServiceUnavailable, false},
		{http.StatusMovedPermanently, true},
		{http.StatusNotFound, true},
		{http.StatusBadRequest, true},
		{http.StatusForbidden, true},
	}
	for _, tc := range cases {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			err := classifyAPIDestinationStatus(&apiDestinationError{statusCode: tc.status})
			if tc.status >= 200 && tc.status < 300 {
				if err != nil {
					t.Fatalf("2xx must succeed, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("non-2xx must fail")
			}
			if got := errors.Is(err, errPermanentDelivery); got != tc.permanent {
				t.Fatalf("permanent = %v, want %v (err %v)", got, tc.permanent, err)
			}
		})
	}
}

// TestApiDestinationDeauthorizedConnectionIsPermanent pins that invoking
// through a deauthorized connection fails permanently: the credentials are
// gone and no retry can recover them.
func TestApiDestinationDeauthorizedConnectionIsPermanent(t *testing.T) {
	ctx := context.Background()
	svc, store, _ := seedApiDestinationFixture(t)
	endpoint := newAPIDestinationTestServer(http.StatusOK)
	defer endpoint.close()

	connection, err := store.GetConnection(ctx, "dest-conn")
	if err != nil {
		t.Fatal(err)
	}
	dest := seedDestinationRecord(t, store, connection.ARN, endpoint.url(), "POST", 0)
	if _, err := svc.deauthorizeConnectionCore(ctx, store, "dest-conn"); err != nil {
		t.Fatal(err)
	}

	err = svc.deliverToApiDestination(ctx, apiDestJob(dest, eventsstore.Target{ARN: dest.ARN}))
	if err == nil || !errors.Is(err, errPermanentDelivery) {
		t.Fatalf("deauthorized connection must fail permanently, got %v", err)
	}
	if len(endpoint.recorded()) != 0 {
		t.Fatal("no request may reach the endpoint without credentials")
	}
}

// TestApiDestinationOAuthTokenExchange pins the OAuth flow: the invocation
// carries the exchanged bearer token, the token is cached across
// invocations, and a 401 forces a synchronous refresh whose new token
// rides the retry.
func TestApiDestinationOAuthTokenExchange(t *testing.T) {
	ctx := context.Background()
	svc, store, fake := newConnectionTestService(t)
	svc.SetEventsStore("us-east-1", store)

	var tokenMu sync.Mutex
	tokenIssues := 0
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "grant_type=client_credentials") {
			t.Errorf("token request must carry the client-credentials grant, got %s", body)
		}
		if !strings.Contains(string(body), "client_id=oauth-id") {
			t.Errorf("token request must carry the client ID, got %s", body)
		}
		tokenMu.Lock()
		tokenIssues++
		issue := tokenIssues
		tokenMu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "tok-" + string(rune('0'+issue)),
			"token_type":   "Bearer",
			"expires_in":   3600,
		})
	}))
	defer tokenServer.Close()

	var endpointMu sync.Mutex
	var authorizations []string
	statuses := []int{http.StatusOK, http.StatusOK, http.StatusUnauthorized, http.StatusOK}
	endpoint := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		endpointMu.Lock()
		authorizations = append(authorizations, r.Header.Get("Authorization"))
		status := statuses[0]
		statuses = statuses[1:]
		endpointMu.Unlock()
		w.WriteHeader(status)
	}))
	defer endpoint.Close()

	// The OAuth connection's credentials go through the fake secret store.
	_, cleanup, err := createOAuthConnectionForTest(ctx, svc, store, tokenServer.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	_ = fake

	connection, err := store.GetConnection(ctx, "oauth-conn")
	if err != nil {
		t.Fatal(err)
	}
	dest := seedDestinationRecord(t, store, connection.ARN, endpoint.URL, "POST", 0)
	job := apiDestJob(dest, eventsstore.Target{ARN: dest.ARN})

	// First invocation: token exchange + authorised request.
	if err := svc.deliverToApiDestination(ctx, job); err != nil {
		t.Fatal(err)
	}
	// Second invocation: served from the token cache.
	if err := svc.deliverToApiDestination(ctx, job); err != nil {
		t.Fatal(err)
	}
	// Third invocation: endpoint answers 401 → token refresh + one retry.
	if err := svc.deliverToApiDestination(ctx, job); err != nil {
		t.Fatal(err)
	}

	if len(authorizations) != 4 {
		t.Fatalf("expected 3 invocations plus the 401 retry, got %d", len(authorizations))
	}
	if authorizations[0] != "Bearer tok-1" || authorizations[1] != "Bearer tok-1" {
		t.Fatalf("cached token must serve the first two invocations, got %v", authorizations[:2])
	}
	if authorizations[2] != "Bearer tok-1" || authorizations[3] != "Bearer tok-2" {
		t.Fatalf("the 401 must force a refresh and the retry must carry the new token, got %v", authorizations[2:])
	}
	tokenMu.Lock()
	issues := tokenIssues
	tokenMu.Unlock()
	if issues != 2 {
		t.Fatalf("exactly two token exchanges expected (initial + 401 refresh), got %d", issues)
	}
}

// TestApiDestinationOAuthTokenRequestParameters pins that all three
// connection HTTP-parameter families reach the token exchange: headers on
// the request, query-string members appended to the authorization
// endpoint, and body members joining the client-credentials form.
func TestApiDestinationOAuthTokenRequestParameters(t *testing.T) {
	ctx := context.Background()
	svc, _, _ := newConnectionTestService(t)

	var mu sync.Mutex
	var gotQuery, gotHeader, gotForm string
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		gotQuery = r.URL.Query().Get("audience")
		gotHeader = r.Header.Get("X-Tenant")
		gotForm = string(body)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token": "tok-params",
			"token_type":   "Bearer",
			"expires_in":   60,
		})
	}))
	defer tokenServer.Close()

	oauth := &eventsstore.OAuthParameters{
		AuthorizationEndpoint: tokenServer.URL,
		HttpMethod:            "POST",
		ClientParameters:      &eventsstore.OAuthClientParameters{ClientID: "id", ClientSecret: "secret"},
		OAuthHttpParameters: &eventsstore.ConnectionHttpParameters{
			HeaderParameters:      []eventsstore.ConnectionHeaderParameter{{Key: "X-Tenant", Value: "acme"}},
			QueryStringParameters: []eventsstore.ConnectionQueryStringParameter{{Key: "audience", Value: "api://orders"}},
			BodyParameters:        []eventsstore.ConnectionBodyParameter{{Key: "scope", Value: "orders.read"}},
		},
	}
	if _, err := svc.exchangeOAuthToken(ctx, "arn:aws:events:us-east-1:000000000000:connection/params/id", oauth); err != nil {
		t.Fatal(err)
	}

	mu.Lock()
	defer mu.Unlock()
	if gotQuery != "api://orders" {
		t.Errorf("token request query = %q, want the connection's query-string parameter", gotQuery)
	}
	if gotHeader != "acme" {
		t.Errorf("token request header = %q, want the connection's header parameter", gotHeader)
	}
	if !strings.Contains(gotForm, "scope=orders.read") {
		t.Errorf("token request form = %q, want the connection's body parameter", gotForm)
	}
	if !strings.Contains(gotForm, "grant_type=client_credentials") {
		t.Errorf("token request form = %q, want the grant alongside the body parameters", gotForm)
	}
}

func createOAuthConnectionForTest(ctx context.Context, svc *EventsService, store *eventsstore.EventsStore, tokenEndpoint string) (*eventsstore.Connection, func(), error) {
	connection, err := svc.createConnectionCore(ctx, store, CreateConnectionInput{
		Name:              "oauth-conn",
		AuthorizationType: "OAUTH_CLIENT_CREDENTIALS",
		AuthParameters: &eventsstore.AuthParameters{
			OAuthParameters: &eventsstore.OAuthParameters{
				AuthorizationEndpoint: tokenEndpoint,
				HttpMethod:            "POST",
				ClientParameters:      &eventsstore.OAuthClientParameters{ClientID: "oauth-id", ClientSecret: "oauth-secret"},
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}
	return connection, func() {}, nil
}

// TestApiDestinationPacerSpacing pins the rate pacer's arithmetic: the
// first slot is immediate, later slots are spaced one interval apart.
func TestApiDestinationPacerSpacing(t *testing.T) {
	pacer := &apiDestinationPacer{interval: 40 * time.Millisecond}
	ctx := context.Background()

	start := time.Now()
	if err := pacer.wait(ctx); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed > 20*time.Millisecond {
		t.Fatalf("first slot must be immediate, took %v", elapsed)
	}

	if err := pacer.wait(ctx); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed < 30*time.Millisecond {
		t.Fatalf("second slot must be spaced by the interval, took %v", elapsed)
	}
}
