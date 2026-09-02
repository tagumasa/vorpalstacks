package cloudfront

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/core/storage"
	cloudfrontstore "vorpalstacks/internal/store/aws/cloudfront"
)

// newCDPTestStores builds a service store bundle over a temporary storage.
func newCDPTestStores(t *testing.T) *cloudfrontStores {
	t.Helper()
	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	service := NewCloudFrontService("123456789012")
	return service.createStores(st)
}

// createStagingDistribution inserts an enabled staging distribution and
// returns it.
func createStagingDistribution(t *testing.T, stores *cloudfrontStores, callerRef string) *cloudfrontstore.Distribution {
	t.Helper()
	dist, err := stores.distributions.Create(callerRef, &cloudfrontstore.DistributionConfig{
		Enabled: true,
		Staging: true,
		Origins: cloudfrontstore.Origins{Quantity: 1, Items: []*cloudfrontstore.Origin{{
			ID:                 "origin-1",
			DomainName:         "127.0.0.1:9",
			CustomOriginConfig: &cloudfrontstore.CustomOriginConfig{OriginProtocolPolicy: "http-only"},
		}}},
		DefaultCacheBehavior: &cloudfrontstore.CacheBehavior{
			TargetOriginId:       "origin-1",
			ViewerProtocolPolicy: "allow-all",
			ForwardedValues:      &cloudfrontstore.ForwardedValues{QueryString: false},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return dist
}

func singleWeightConfig(weight float64) *cloudfrontstore.TrafficConfig {
	return &cloudfrontstore.TrafficConfig{
		Type:               "SingleWeight",
		SingleWeightConfig: &cloudfrontstore.ContinuousDeploymentSingleWeightConfig{Weight: weight},
	}
}

func cdpConfig(staging *cloudfrontstore.Distribution, traffic *cloudfrontstore.TrafficConfig) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
	return &cloudfrontstore.ContinuousDeploymentPolicyConfig{
		StagingDistributionDnsNames: &cloudfrontstore.StagingDistributionDnsNames{
			Quantity: 1,
			Items:    []string{staging.DomainName},
		},
		Enabled:       true,
		TrafficConfig: traffic,
	}
}

func TestContinuousDeploymentPolicyCRUD(t *testing.T) {
	service := NewCloudFrontService("123456789012")
	stores := newCDPTestStores(t)
	staging := createStagingDistribution(t, stores, "staging-1")

	created, err := service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
		Config: cdpConfig(staging, singleWeightConfig(0.10)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID == "" || created.ETag == "" {
		t.Fatalf("created policy lacks an identity: %+v", created)
	}

	got, err := service.getContinuousDeploymentPolicyCore(stores, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ContinuousDeploymentPolicyConfig.TrafficConfig.SingleWeightConfig.Weight != 0.10 {
		t.Fatalf("round-tripped weight = %v", got.ContinuousDeploymentPolicyConfig.TrafficConfig.SingleWeightConfig.Weight)
	}

	updated, err := service.updateContinuousDeploymentPolicyCore(stores, UpdateContinuousDeploymentPolicyInput{
		Id:      created.ID,
		IfMatch: created.ETag,
		Config:  cdpConfig(staging, singleWeightConfig(0.05)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.ETag == created.ETag {
		t.Fatal("update must rotate the ETag")
	}

	// A policy referenced by a distribution cannot be deleted.
	primary, err := stores.distributions.Create("primary-1", &cloudfrontstore.DistributionConfig{
		Enabled: true,
		Origins: cloudfrontstore.Origins{Quantity: 1, Items: []*cloudfrontstore.Origin{{
			ID:                 "origin-1",
			DomainName:         "127.0.0.1:9",
			CustomOriginConfig: &cloudfrontstore.CustomOriginConfig{OriginProtocolPolicy: "http-only"},
		}}},
		DefaultCacheBehavior: &cloudfrontstore.CacheBehavior{
			TargetOriginId:       "origin-1",
			ViewerProtocolPolicy: "allow-all",
			ForwardedValues:      &cloudfrontstore.ForwardedValues{QueryString: false},
		},
		ContinuousDeploymentPolicyId: created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.deleteContinuousDeploymentPolicyCore(stores, created.ID, updated.ETag); err == nil {
		t.Fatal("delete must fail while a distribution references the policy")
	}

	primary.DistributionConfig.ContinuousDeploymentPolicyId = ""
	if _, err := stores.distributions.Update(primary.ID, primary.DistributionConfig); err != nil {
		t.Fatal(err)
	}
	if err := service.deleteContinuousDeploymentPolicyCore(stores, created.ID, updated.ETag); err != nil {
		t.Fatalf("delete after detach failed: %v", err)
	}
	if _, err := service.getContinuousDeploymentPolicyCore(stores, created.ID); err == nil {
		t.Fatal("deleted policy must be gone")
	}
}

func TestContinuousDeploymentPolicyValidation(t *testing.T) {
	service := NewCloudFrontService("123456789012")
	// Each case builds its configuration against its own isolated
	// storage, so one case's outcome cannot influence another's.
	cases := []struct {
		name string
		cfg  func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig
		code string
	}{
		{
			name: "unknown staging distribution",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(&cloudfrontstore.Distribution{DomainName: "d99999999.cloudfront.net"}, singleWeightConfig(0.1))
			},
			code: "InvalidArgument",
		},
		{
			name: "staging name of a non-staging distribution",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				normal, err := stores.distributions.Create("normal-"+t.Name(), &cloudfrontstore.DistributionConfig{Enabled: true})
				if err != nil {
					t.Fatal(err)
				}
				return cdpConfig(&cloudfrontstore.Distribution{DomainName: normal.DomainName}, singleWeightConfig(0.1))
			},
			code: "InvalidArgument",
		},
		{
			name: "quantity mismatch",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				cfg := cdpConfig(createStagingDistribution(t, stores, "s-qty"), singleWeightConfig(0.1))
				cfg.StagingDistributionDnsNames.Quantity = 2
				return cfg
			},
			code: "InconsistentQuantities",
		},
		{
			name: "weight above the quota",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-wmax"), singleWeightConfig(0.2))
			},
			code: "InvalidArgument",
		},
		{
			name: "negative weight",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-wmin"), singleWeightConfig(-0.1))
			},
			code: "InvalidArgument",
		},
		{
			name: "stickiness idle below the floor",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-idlelo"), &cloudfrontstore.TrafficConfig{
					Type: "SingleWeight",
					SingleWeightConfig: &cloudfrontstore.ContinuousDeploymentSingleWeightConfig{
						Weight:                  0.1,
						SessionStickinessConfig: &cloudfrontstore.SessionStickinessConfig{IdleTTL: 60, MaximumTTL: 600},
					},
				})
			},
			code: "InvalidArgument",
		},
		{
			name: "stickiness maximum above the ceiling",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-ttlmax"), &cloudfrontstore.TrafficConfig{
					Type: "SingleWeight",
					SingleWeightConfig: &cloudfrontstore.ContinuousDeploymentSingleWeightConfig{
						Weight:                  0.1,
						SessionStickinessConfig: &cloudfrontstore.SessionStickinessConfig{IdleTTL: 3600, MaximumTTL: 4000},
					},
				})
			},
			code: "InvalidArgument",
		},
		{
			name: "stickiness idle above maximum TTL",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-order"), &cloudfrontstore.TrafficConfig{
					Type: "SingleWeight",
					SingleWeightConfig: &cloudfrontstore.ContinuousDeploymentSingleWeightConfig{
						Weight:                  0.1,
						SessionStickinessConfig: &cloudfrontstore.SessionStickinessConfig{IdleTTL: 1800, MaximumTTL: 600},
					},
				})
			},
			code: "InvalidArgument",
		},
		{
			name: "header without the aws-cf-cd- prefix",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-prefix"), &cloudfrontstore.TrafficConfig{
					Type: "SingleHeader",
					SingleHeaderConfig: &cloudfrontstore.ContinuousDeploymentSingleHeaderConfig{
						Header: "x-test",
						Value:  "1",
					},
				})
			},
			code: "InvalidArgument",
		},
		{
			name: "traffic type without its configuration",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-missing"), &cloudfrontstore.TrafficConfig{Type: "SingleHeader"})
			},
			code: "InvalidArgument",
		},
		{
			name: "both traffic configurations present",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-both"), &cloudfrontstore.TrafficConfig{
					Type:               "SingleWeight",
					SingleWeightConfig: &cloudfrontstore.ContinuousDeploymentSingleWeightConfig{Weight: 0.1},
					SingleHeaderConfig: &cloudfrontstore.ContinuousDeploymentSingleHeaderConfig{
						Header: "aws-cf-cd-test",
						Value:  "1",
					},
				})
			},
			code: "InvalidArgument",
		},
		{
			name: "unknown traffic type",
			cfg: func(t *testing.T, stores *cloudfrontStores) *cloudfrontstore.ContinuousDeploymentPolicyConfig {
				return cdpConfig(createStagingDistribution(t, stores, "s-type"), &cloudfrontstore.TrafficConfig{Type: "RoundRobin"})
			},
			code: "InvalidArgument",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			isolated := newCDPTestStores(t)
			_, err := service.createContinuousDeploymentPolicyCore(isolated,
				CreateContinuousDeploymentPolicyInput{Config: tc.cfg(t, isolated)})
			if err == nil {
				t.Fatal("validation must reject the configuration")
			}
			if !strings.Contains(err.Error(), tc.code) {
				t.Fatalf("error = %v, want code %s", err, tc.code)
			}
		})
	}
}

func TestContinuousDeploymentPolicyStagingExclusivityAndQuota(t *testing.T) {
	service := NewCloudFrontService("123456789012")
	stores := newCDPTestStores(t)
	staging := createStagingDistribution(t, stores, "staging-excl")

	if _, err := service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
		Config: cdpConfig(staging, singleWeightConfig(0.05)),
	}); err != nil {
		t.Fatal(err)
	}
	_, err := service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
		Config: cdpConfig(staging, singleWeightConfig(0.05)),
	})
	if err == nil || !strings.Contains(err.Error(), "StagingDistributionInUse") {
		t.Fatalf("second policy on one staging distribution = %v, want StagingDistributionInUse", err)
	}

	// The account quota is pinned by filling up to the limit.
	for i := 1; i < cloudfrontstore.MaxContinuousDeploymentPolicies; i++ {
		dist := createStagingDistribution(t, stores, fmt.Sprintf("staging-quota-%d", i))
		if _, err := service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
			Config: cdpConfig(dist, singleWeightConfig(0.05)),
		}); err != nil {
			t.Fatalf("policy %d: %v", i, err)
		}
	}
	overflow := createStagingDistribution(t, stores, "staging-overflow")
	_, err = service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
		Config: cdpConfig(overflow, singleWeightConfig(0.05)),
	})
	if err == nil || !strings.Contains(err.Error(), "TooManyContinuousDeploymentPolicies") {
		t.Fatalf("policy over the quota = %v, want TooManyContinuousDeploymentPolicies", err)
	}
}

func TestContinuousDeploymentPolicyIfMatch(t *testing.T) {
	service := NewCloudFrontService("123456789012")
	stores := newCDPTestStores(t)
	staging := createStagingDistribution(t, stores, "staging-ifmatch")
	created, err := service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
		Config: cdpConfig(staging, singleWeightConfig(0.05)),
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.updateContinuousDeploymentPolicyCore(stores, UpdateContinuousDeploymentPolicyInput{
		Id: created.ID, Config: cdpConfig(staging, singleWeightConfig(0.1)),
	}); err == nil || !strings.Contains(err.Error(), "InvalidIfMatchVersion") {
		t.Fatalf("missing If-Match = %v, want InvalidIfMatchVersion", err)
	}
	if _, err := service.updateContinuousDeploymentPolicyCore(stores, UpdateContinuousDeploymentPolicyInput{
		Id: created.ID, IfMatch: "stale", Config: cdpConfig(staging, singleWeightConfig(0.1)),
	}); err == nil || !strings.Contains(err.Error(), "PreconditionFailed") {
		t.Fatalf("stale If-Match = %v, want PreconditionFailed", err)
	}
	if err := service.deleteContinuousDeploymentPolicyCore(stores, "missing-id", "*"); err == nil ||
		!strings.Contains(err.Error(), "NoSuchContinuousDeploymentPolicy") {
		t.Fatalf("missing policy delete = %v, want NoSuchContinuousDeploymentPolicy", err)
	}
}

func TestDistributionPolicyReferenceValidation(t *testing.T) {
	service := NewCloudFrontService("123456789012")
	stores := newCDPTestStores(t)
	staging := createStagingDistribution(t, stores, "staging-ref")
	policy, err := service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
		Config: cdpConfig(staging, singleWeightConfig(0.05)),
	})
	if err != nil {
		t.Fatal(err)
	}

	baseConfig := func() *cloudfrontstore.DistributionConfig {
		return &cloudfrontstore.DistributionConfig{
			Enabled: true,
			Origins: cloudfrontstore.Origins{Quantity: 1, Items: []*cloudfrontstore.Origin{{
				ID:                 "o",
				DomainName:         "127.0.0.1:9",
				CustomOriginConfig: &cloudfrontstore.CustomOriginConfig{OriginProtocolPolicy: "http-only"},
			}}},
			DefaultCacheBehavior: &cloudfrontstore.CacheBehavior{
				TargetOriginId:       "o",
				ViewerProtocolPolicy: "allow-all",
				ForwardedValues:      &cloudfrontstore.ForwardedValues{QueryString: false},
			},
		}
	}

	unknown := baseConfig()
	unknown.ContinuousDeploymentPolicyId = "does-not-exist"
	if err := service.validateDistributionPolicyReference(stores, unknown); err == nil ||
		!strings.Contains(err.Error(), "NoSuchContinuousDeploymentPolicy") {
		t.Fatalf("unknown policy reference = %v", err)
	}

	valid := baseConfig()
	valid.ContinuousDeploymentPolicyId = policy.ID
	if err := service.validateDistributionPolicyReference(stores, valid); err != nil {
		t.Fatalf("valid reference rejected: %v", err)
	}

	stagingRef := baseConfig()
	stagingRef.Staging = true
	stagingRef.ContinuousDeploymentPolicyId = policy.ID
	if err := service.validateDistributionPolicyReference(stores, stagingRef); err == nil {
		t.Fatal("a staging distribution must not carry a policy")
	}
}

// --- runtime traffic routing ---

// newCDPRoutingServer builds a distribution server whose primary and
// staging distributions proxy to distinct test origins; each origin
// replies with a marker naming the serving distribution, so the routing
// decision is observable in the response body.
func newCDPRoutingServer(t *testing.T, traffic *cloudfrontstore.TrafficConfig) (*DistributionServer, *cloudfrontstore.Distribution, *cloudfrontStores) {
	t.Helper()
	primaryOrigin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("served-by-primary"))
	}))
	t.Cleanup(primaryOrigin.Close)
	stagingOrigin := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("served-by-staging"))
	}))
	t.Cleanup(stagingOrigin.Close)

	st, err := storage.Open(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	service := NewCloudFrontService("123456789012")
	stores := service.createStores(st)

	buildConfig := func(staging bool, originDomain string) *cloudfrontstore.DistributionConfig {
		return &cloudfrontstore.DistributionConfig{
			Enabled: true,
			Staging: staging,
			Origins: cloudfrontstore.Origins{Quantity: 1, Items: []*cloudfrontstore.Origin{{
				ID:                 "origin-1",
				DomainName:         originDomain,
				CustomOriginConfig: &cloudfrontstore.CustomOriginConfig{OriginProtocolPolicy: "http-only"},
			}}},
			DefaultCacheBehavior: &cloudfrontstore.CacheBehavior{
				TargetOriginId:       "origin-1",
				ViewerProtocolPolicy: "allow-all",
				ForwardedValues:      &cloudfrontstore.ForwardedValues{QueryString: false},
				MinTTL:               0,
				DefaultTTL:           0,
				MaxTTL:               0,
			},
		}
	}

	primary, err := stores.distributions.Create("cdp-primary",
		buildConfig(false, strings.TrimPrefix(primaryOrigin.URL, "http://")))
	if err != nil {
		t.Fatal(err)
	}
	stagingDist, err := stores.distributions.Create("cdp-staging",
		buildConfig(true, strings.TrimPrefix(stagingOrigin.URL, "http://")))
	if err != nil {
		t.Fatal(err)
	}

	policyStore := cloudfrontstore.NewContinuousDeploymentPolicyStore(st)
	policy, err := policyStore.Create(cdpConfig(stagingDist, traffic))
	if err != nil {
		t.Fatal(err)
	}
	primary.DistributionConfig.ContinuousDeploymentPolicyId = policy.ID
	if _, err := stores.distributions.Update(primary.ID, primary.DistributionConfig); err != nil {
		t.Fatal(err)
	}

	server := &DistributionServer{
		accountID: "123456789012",
		client:    &http.Client{Timeout: 5 * time.Second},
		cache:     newResponseCache(),
	}
	server.SetDistributionStore(stores.distributions)
	server.deploymentPolicies = policyStore
	return server, primary, stores
}

// cdpRequest issues a request through the server; the origin markers
// name the distribution that served it.
func cdpRequest(t *testing.T, server *DistributionServer, dist *cloudfrontstore.Distribution, path string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, "http://edge.example.test"+path, nil)
	r.Host = dist.ID + ".cloudfront.net"
	for name, value := range headers {
		r.Header.Set(name, value)
	}
	recorder := httptest.NewRecorder()
	server.HandleRequest(recorder, r)
	return recorder
}

func TestContinuousDeploymentStagingDirectAccessBlocked(t *testing.T) {
	server, _, stores := newCDPRoutingServer(t, singleWeightConfig(0.0))
	stagingDist, err := stores.distributions.GetByCallerReference("cdp-staging")
	if err != nil {
		t.Fatal(err)
	}
	response := cdpRequest(t, server, stagingDist, "/obj", nil)
	if response.Code != http.StatusForbidden {
		t.Fatalf("direct staging access = %d, want 403", response.Code)
	}
}

func TestContinuousDeploymentSingleHeaderRouting(t *testing.T) {
	server, primary, _ := newCDPRoutingServer(t, &cloudfrontstore.TrafficConfig{
		Type: "SingleHeader",
		SingleHeaderConfig: &cloudfrontstore.ContinuousDeploymentSingleHeaderConfig{
			Header: "aws-cf-cd-test",
			Value:  "staging-please",
		},
	})

	if direct := cdpRequest(t, server, primary, "/plain", nil); direct.Body.String() != "served-by-primary" {
		t.Fatalf("plain request = %q, want the primary origin", direct.Body.String())
	}
	if staged := cdpRequest(t, server, primary, "/staged", map[string]string{"aws-cf-cd-test": "staging-please"}); staged.Body.String() != "served-by-staging" {
		t.Fatalf("header-routed request = %q, want the staging origin", staged.Body.String())
	}
	if wrong := cdpRequest(t, server, primary, "/wrong", map[string]string{"aws-cf-cd-test": "other"}); wrong.Body.String() != "served-by-primary" {
		t.Fatalf("mismatched header value = %q, want the primary origin", wrong.Body.String())
	}
}

func TestContinuousDeploymentSingleWeightRouting(t *testing.T) {
	server, primary, _ := newCDPRoutingServer(t, singleWeightConfig(0.0))
	if response := cdpRequest(t, server, primary, "/obj", nil); response.Body.String() != "served-by-primary" {
		t.Fatalf("weight 0 must keep every request on the primary, got %q", response.Body.String())
	}

	full, primaryDist, _ := newCDPRoutingServer(t, singleWeightConfig(1.0))
	if response := cdpRequest(t, full, primaryDist, "/obj", nil); response.Body.String() != "served-by-staging" {
		t.Fatalf("weight 1 must route to the staging distribution, got %q", response.Body.String())
	}
}

func TestContinuousDeploymentStickinessCookie(t *testing.T) {
	server, primary, _ := newCDPRoutingServer(t, &cloudfrontstore.TrafficConfig{
		Type: "SingleWeight",
		SingleWeightConfig: &cloudfrontstore.ContinuousDeploymentSingleWeightConfig{
			Weight: 1.0,
			SessionStickinessConfig: &cloudfrontstore.SessionStickinessConfig{
				IdleTTL:    300,
				MaximumTTL: 600,
			},
		},
	})

	response := cdpRequest(t, server, primary, "/sticky", nil)
	if response.Code != http.StatusOK {
		t.Fatalf("sticky request = %d", response.Code)
	}
	cookie := response.Header().Get("Set-Cookie")
	if !strings.Contains(cookie, stickyCDPCookieName) {
		t.Fatalf("weight policy with stickiness must set the session cookie, got %q", cookie)
	}
	if !strings.Contains(cookie, "staging") && !strings.Contains(cookie, "primary") {
		t.Fatalf("session cookie carries the assignment, got %q", cookie)
	}
}

// TestStickyAssignmentSessionBounds pins the two-clock cookie form: the
// idle TTL runs from the last refresh while the maximum TTL runs from the
// session start, so a viewer that keeps returning within the idle window
// still leaves the session once the maximum session time passes.
func TestStickyAssignmentSessionBounds(t *testing.T) {
	cfg := &cloudfrontstore.SessionStickinessConfig{IdleTTL: 300, MaximumTTL: 600}
	now := time.Now().Truncate(time.Second)
	cookieValue := func(staging bool, start, last time.Time) string {
		assignment := "primary"
		if staging {
			assignment = "staging"
		}
		return fmt.Sprintf("%s.%d.%d", assignment, start.Unix(), last.Unix())
	}
	requestWithCookie := func(value string) *http.Request {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: stickyCDPCookieName, Value: value})
		return req
	}

	fresh := httptest.NewRecorder()
	if got := stickyAssignment(fresh, requestWithCookie(cookieValue(true, now.Add(-100*time.Second), now.Add(-50*time.Second))), false, cfg, now); !got {
		t.Fatal("a session inside both TTLs lost its assignment")
	}
	reissued := fresh.Header().Get("Set-Cookie")
	if !strings.Contains(reissued, fmt.Sprintf(".%d.%d", now.Add(-100*time.Second).Unix(), now.Unix())) {
		t.Fatalf("cookie refresh must keep the session start and update the last refresh, got %q", reissued)
	}

	expiredMax := httptest.NewRecorder()
	if got := stickyAssignment(expiredMax, requestWithCookie(cookieValue(true, now.Add(-601*time.Second), now.Add(-10*time.Second))), false, cfg, now); got {
		t.Fatal("a session past the maximum TTL kept its assignment despite refreshing within the idle window")
	}

	expiredIdle := httptest.NewRecorder()
	if got := stickyAssignment(expiredIdle, requestWithCookie(cookieValue(true, now.Add(-400*time.Second), now.Add(-301*time.Second))), true, cfg, now); !got {
		t.Fatal("an idle-expired session must fall back to the fresh weight assignment (true here)")
	}
	if issued := expiredIdle.Header().Get("Set-Cookie"); !strings.Contains(issued, fmt.Sprintf(".%d.%d", now.Unix(), now.Unix())) {
		t.Fatalf("a new session must carry the current time as both clocks, got %q", issued)
	}
}

// TestAssociateDistributionWebACLWithPolicy pins the continuous-deployment
// association rule: while a policy is attached, associating a new web ACL
// (or changing to a different one) is rejected, but re-associating the ACL
// already in place succeeds.
func TestAssociateDistributionWebACLWithPolicy(t *testing.T) {
	service := NewCloudFrontService("123456789012")
	stores := newCDPTestStores(t)
	staging := createStagingDistribution(t, stores, "staging-waf")
	policy, err := service.createContinuousDeploymentPolicyCore(stores, CreateContinuousDeploymentPolicyInput{
		Config: cdpConfig(staging, singleWeightConfig(0.05)),
	})
	if err != nil {
		t.Fatal(err)
	}
	config := &cloudfrontstore.DistributionConfig{
		Enabled:                      true,
		ContinuousDeploymentPolicyId: policy.ID,
		WebACLId:                     "arn:aws:wafv2:us-east-1:123456789012:webacl/pinned-acl/1",
		Origins: cloudfrontstore.Origins{Quantity: 1, Items: []*cloudfrontstore.Origin{{
			ID:                 "o",
			DomainName:         "127.0.0.1:9",
			CustomOriginConfig: &cloudfrontstore.CustomOriginConfig{OriginProtocolPolicy: "http-only"},
		}}},
		DefaultCacheBehavior: &cloudfrontstore.CacheBehavior{
			TargetOriginId:       "o",
			ViewerProtocolPolicy: "allow-all",
			ForwardedValues:      &cloudfrontstore.ForwardedValues{QueryString: false},
		},
	}
	dist, err := stores.distributions.Create("waf-assoc-ref", config)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.associateDistributionWebACLCore(context.Background(), stores, AssociateDistributionWebACLInput{
		Id:        dist.ID,
		WebACLArn: config.WebACLId,
	}); err != nil {
		t.Fatalf("re-associating the pinned web ACL must stay allowed: %v", err)
	}
	if _, err := service.associateDistributionWebACLCore(context.Background(), stores, AssociateDistributionWebACLInput{
		Id:        dist.ID,
		WebACLArn: "arn:aws:wafv2:us-east-1:123456789012:webacl/other-acl/1",
	}); err == nil || !strings.Contains(err.Error(), "InvalidArgument") {
		t.Fatalf("changing the web ACL under an attached policy must be rejected, got %v", err)
	}
}
