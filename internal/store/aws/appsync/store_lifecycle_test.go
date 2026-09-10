package appsync

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"vorpalstacks/internal/core/storage"
)

func newLifecycleTestStore(t *testing.T) *AppSyncStore {
	t.Helper()
	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })
	return NewAppSyncStore(ps, "123456789012", "us-east-1")
}

// Deleting an Event API must cascade to its API keys: they are keyed under
// apiId/ in the shared api-keys bucket and would otherwise outlive the API.
func TestDeleteApiByIdCascadesApiKeys(t *testing.T) {
	store := newLifecycleTestStore(t)

	api, err := store.CreateApi(&Api{
		Name: "cascade-event",
		EventConfig: &EventConfig{
			ConnectionAuthModes:       []AuthMode{{AuthType: "API_KEY"}},
			DefaultPublishAuthModes:   []AuthMode{{AuthType: "API_KEY"}},
			DefaultSubscribeAuthModes: []AuthMode{{AuthType: "API_KEY"}},
		},
	})
	if err != nil {
		t.Fatalf("CreateApi: %v", err)
	}
	if err := store.CreateApiKey(api.ApiId, &ApiKey{Id: "k1", Expires: 2000000000}); err != nil {
		t.Fatalf("CreateApiKey: %v", err)
	}

	if err := store.DeleteApiById(api.ApiId); err != nil {
		t.Fatalf("DeleteApiById: %v", err)
	}

	_, err = store.GetApiKey(api.ApiId, "k1")
	assert.Equal(t, ErrApiKeyNotFound, err, "API keys must not outlive the Event API")
	_, err = store.GetApiById(api.ApiId)
	assert.Equal(t, ErrApiNotFound, err)
}

// Deleting a GraphQL API must cascade to its resolver cache entries and API
// keys, both keyed under apiId/ prefixes.
func TestDeleteGraphqlApiByIdCascadesResolverCacheAndApiKeys(t *testing.T) {
	store := newLifecycleTestStore(t)

	api, err := store.CreateGraphqlApi(&GraphqlApi{Name: "cascade-gql", AuthenticationType: "API_KEY"})
	if err != nil {
		t.Fatalf("CreateGraphqlApi: %v", err)
	}
	if err := store.CreateApiKey(api.ApiId, &ApiKey{Id: "k1", Expires: 2000000000}); err != nil {
		t.Fatalf("CreateApiKey: %v", err)
	}
	if err := store.PutResolverCacheEntry(api.ApiId, "ck", &ResolverCacheEntry{
		Result:   json.RawMessage(`{}`),
		CachedAt: 1,
		TTL:      60,
	}); err != nil {
		t.Fatalf("PutResolverCacheEntry: %v", err)
	}

	if err := store.DeleteGraphqlApiById(api.ApiId); err != nil {
		t.Fatalf("DeleteGraphqlApiById: %v", err)
	}

	_, err = store.GetResolverCacheEntry(api.ApiId, "ck")
	assert.Error(t, err, "resolver cache entries must not outlive the GraphQL API")
	_, err = store.GetApiKey(api.ApiId, "k1")
	assert.Equal(t, ErrApiKeyNotFound, err)
	_, err = store.GetGraphqlApiById(api.ApiId)
	assert.Equal(t, ErrGraphqlApiNotFound, err)
}

// Deleting a resource must remove its tag rows with it — both on the
// per-resource deletes and on the API-level cascades, whose bulk prefix
// deletes sweep the child ARNs first.
func TestDeletesRemoveTagRows(t *testing.T) {
	store := newLifecycleTestStore(t)

	gqlApi, err := store.CreateGraphqlApi(&GraphqlApi{Name: "tag-sweep-gql", AuthenticationType: "API_KEY"})
	if err != nil {
		t.Fatalf("CreateGraphqlApi: %v", err)
	}
	if _, err := store.CreateDataSource(&DataSource{ApiId: gqlApi.ApiId, Name: "ds1", Type: "NONE"}); err != nil {
		t.Fatalf("CreateDataSource: %v", err)
	}
	if _, err := store.CreateResolver(&Resolver{ApiId: gqlApi.ApiId, TypeName: "T", FieldName: "f", Kind: "UNIT"}); err != nil {
		t.Fatalf("CreateResolver: %v", err)
	}
	if err := store.CreateApiKey(gqlApi.ApiId, &ApiKey{Id: "k1", Expires: 2000000000}); err != nil {
		t.Fatalf("CreateApiKey: %v", err)
	}

	dsArn := store.BuildDataSourceARN(gqlApi.ApiId, "ds1")
	if err := store.TagStore.Tag(dsArn, map[string]string{"tier": "ds"}); err != nil {
		t.Fatalf("Tag datasource: %v", err)
	}
	resolverArn := store.BuildResolverARN(gqlApi.ApiId, "T", "f")
	if err := store.TagStore.Tag(resolverArn, map[string]string{"tier": "resolver"}); err != nil {
		t.Fatalf("Tag resolver: %v", err)
	}
	keyArn := store.BuildApiKeyARN(gqlApi.ApiId, "k1")
	if err := store.TagStore.Tag(keyArn, map[string]string{"tier": "key"}); err != nil {
		t.Fatalf("Tag api key: %v", err)
	}
	if err := store.TagStore.Tag(gqlApi.Arn, map[string]string{"tier": "api"}); err != nil {
		t.Fatalf("Tag api: %v", err)
	}

	// Per-resource delete removes the record's tag row.
	if err := store.DeleteDataSource(gqlApi.ApiId, "ds1"); err != nil {
		t.Fatalf("DeleteDataSource: %v", err)
	}
	tags, err := store.TagStore.List(dsArn)
	if err != nil {
		t.Fatalf("List datasource tags: %v", err)
	}
	assert.Empty(t, tags, "datasource tag rows must not outlive the data source")

	// API-level cascade removes every remaining child tag row.
	if err := store.DeleteGraphqlApiById(gqlApi.ApiId); err != nil {
		t.Fatalf("DeleteGraphqlApiById: %v", err)
	}
	for arn, label := range map[string]string{
		resolverArn: "resolver",
		keyArn:      "api key",
		gqlApi.Arn:  "api",
	} {
		tags, err := store.TagStore.List(arn)
		if err != nil {
			t.Fatalf("List %s tags: %v", label, err)
		}
		assert.Empty(t, tags, "%s tag rows must not outlive the API", label)
	}
}

// Deleting an Event API must sweep the channel namespaces' tag rows along
// with the bulk prefix delete of the namespace records.
func TestDeleteApiByIdSweepsChannelNamespaceTags(t *testing.T) {
	store := newLifecycleTestStore(t)

	api, err := store.CreateApi(&Api{
		Name: "tag-sweep-event",
		EventConfig: &EventConfig{
			ConnectionAuthModes:       []AuthMode{{AuthType: "API_KEY"}},
			DefaultPublishAuthModes:   []AuthMode{{AuthType: "API_KEY"}},
			DefaultSubscribeAuthModes: []AuthMode{{AuthType: "API_KEY"}},
		},
	})
	if err != nil {
		t.Fatalf("CreateApi: %v", err)
	}
	ns, err := store.CreateChannelNamespace(&ChannelNamespace{ApiId: api.ApiId, Name: "default"})
	if err != nil {
		t.Fatalf("CreateChannelNamespace: %v", err)
	}
	if err := store.TagStore.Tag(ns.ChannelNamespaceArn, map[string]string{"tier": "ns"}); err != nil {
		t.Fatalf("Tag namespace: %v", err)
	}

	if err := store.DeleteApiById(api.ApiId); err != nil {
		t.Fatalf("DeleteApiById: %v", err)
	}

	tags, err := store.TagStore.List(ns.ChannelNamespaceArn)
	if err != nil {
		t.Fatalf("List namespace tags: %v", err)
	}
	assert.Empty(t, tags, "namespace tag rows must not outlive the Event API")
}

// An explicitly supplied zero limit is a valid value and must survive the
// update merge; an omitted limit must leave the stored value untouched.
func TestExplicitZeroLimitsRoundTrip(t *testing.T) {
	store := newLifecycleTestStore(t)

	api, err := store.CreateGraphqlApi(&GraphqlApi{
		Name:               "limits",
		AuthenticationType: "API_KEY",
		QueryDepthLimit:    8,
		QueryDepthLimitSet: true,
	})
	if err != nil {
		t.Fatalf("CreateGraphqlApi: %v", err)
	}

	updated, err := store.UpdateGraphqlApiById(api.ApiId, &GraphqlApi{
		QueryDepthLimit:    0,
		QueryDepthLimitSet: true,
	})
	if err != nil {
		t.Fatalf("UpdateGraphqlApiById explicit zero: %v", err)
	}
	assert.Equal(t, int32(0), updated.QueryDepthLimit, "explicit zero must overwrite the stored limit")

	updated, err = store.UpdateGraphqlApiById(api.ApiId, &GraphqlApi{})
	if err != nil {
		t.Fatalf("UpdateGraphqlApiById omitted: %v", err)
	}
	assert.Equal(t, int32(0), updated.QueryDepthLimit, "omitted limit must leave the stored value untouched")

	// Resolver MaxBatchSize follows the same presence semantics.
	if _, err := store.CreateResolver(&Resolver{ApiId: api.ApiId, TypeName: "T", FieldName: "f", Kind: "UNIT"}); err != nil {
		t.Fatalf("CreateResolver: %v", err)
	}
	res, err := store.UpdateResolver(&Resolver{
		ApiId:           api.ApiId,
		TypeName:        "T",
		FieldName:       "f",
		MaxBatchSize:    0,
		MaxBatchSizeSet: true,
	})
	if err != nil {
		t.Fatalf("UpdateResolver explicit zero: %v", err)
	}
	assert.Equal(t, int32(0), res.MaxBatchSize, "explicit zero must overwrite the stored batch size")
}

// --- fault-injection seam for storage-failure paths ---

// faultBucket wraps a storage bucket, failing Put for keys with the given
// prefix and Get for the exact key, so tests can exercise storage failures
// on the apiId-index paths while record writes succeed.
type faultBucket struct {
	storage.Bucket
	failPutPrefix string
	failGetKey    string
}

func (b *faultBucket) Put(key, value []byte) error {
	if b.failPutPrefix != "" && strings.HasPrefix(string(key), b.failPutPrefix) {
		return errors.New("injected put failure")
	}
	return b.Bucket.Put(key, value)
}

func (b *faultBucket) Get(key []byte) ([]byte, error) {
	if b.failGetKey != "" && string(key) == b.failGetKey {
		return nil, errors.New("injected get failure")
	}
	return b.Bucket.Get(key)
}

// faultStorage wraps a BasicStorage so the named bucket is fault-injecting.
type faultStorage struct {
	storage.BasicStorage
	bucketName    string
	failPutPrefix string
	failGetKey    string
}

func (s *faultStorage) Bucket(name string) storage.Bucket {
	b := s.BasicStorage.Bucket(name)
	if name == s.bucketName {
		return &faultBucket{Bucket: b, failPutPrefix: s.failPutPrefix, failGetKey: s.failGetKey}
	}
	return b
}

func newFaultLifecycleTestStore(t *testing.T, bucketName, failPutPrefix, failGetKey string) *AppSyncStore {
	t.Helper()
	ps, err := storage.NewPebbleStorage(&storage.Config{
		Path:       t.TempDir(),
		TTLEnabled: false,
	})
	if err != nil {
		t.Fatalf("NewPebbleStorage: %v", err)
	}
	t.Cleanup(func() { ps.Close() })
	return NewAppSyncStore(&faultStorage{
		BasicStorage:  ps,
		bucketName:    bucketName,
		failPutPrefix: failPutPrefix,
		failGetKey:    failGetKey,
	}, "123456789012", "us-east-1")
}

// A failed apiId-index write must fail the create and roll the record
// back: the index is the only ID-lookup path, so a create that succeeded
// without its index would leave a record unreachable by ID.
func TestCreateApiFailsWhenIndexWriteFails(t *testing.T) {
	store := newFaultLifecycleTestStore(t, apiBucketName("us-east-1"), "#id:", "")

	_, err := store.CreateApi(&Api{
		Name: "orphan-guard",
		EventConfig: &EventConfig{
			ConnectionAuthModes:       []AuthMode{{AuthType: "API_KEY"}},
			DefaultPublishAuthModes:   []AuthMode{{AuthType: "API_KEY"}},
			DefaultSubscribeAuthModes: []AuthMode{{AuthType: "API_KEY"}},
		},
	})
	if err == nil {
		t.Fatal("expected create to fail when the index write fails")
	}
	if _, err := store.GetApi("orphan-guard"); err == nil {
		t.Fatal("record must be rolled back after the failed index write")
	}
}

// The GraphQL-plane creator shares the index-write obligation.
func TestCreateGraphqlApiFailsWhenIndexWriteFails(t *testing.T) {
	store := newFaultLifecycleTestStore(t, graphqlApiBucketName("us-east-1"), "#id:", "")

	_, err := store.CreateGraphqlApi(&GraphqlApi{Name: "orphan-guard", AuthenticationType: "API_KEY"})
	if err == nil {
		t.Fatal("expected create to fail when the index write fails")
	}
	if _, err := store.GetGraphqlApi("orphan-guard"); err == nil {
		t.Fatal("record must be rolled back after the failed index write")
	}
}

// A storage failure on the index read must surface as an error, not
// collapse into NotFound: only a genuinely absent index entry means the
// API does not exist.
func TestGetApiByIdDistinguishesIndexReadFailureFromAbsence(t *testing.T) {
	healthy := newLifecycleTestStore(t)
	if _, err := healthy.GetApiById("no-such-api"); err != ErrApiNotFound {
		t.Fatalf("absent API must report ErrApiNotFound, got %v", err)
	}

	store := newFaultLifecycleTestStore(t, apiBucketName("us-east-1"), "", apiIdIndexKey("no-such-api"))
	_, err := store.GetApiById("no-such-api")
	if err == nil || err == ErrApiNotFound {
		t.Fatalf("storage failure must not surface as NotFound, got %v", err)
	}
}
