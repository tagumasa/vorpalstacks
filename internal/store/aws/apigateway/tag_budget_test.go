package apigateway

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
)

// The inline REST API tag merge honours API Gateway's documented fifty on
// the merged set: two sequential under-cap writes that together pass the
// cap leave the second rejected with the service's BadRequestException
// identity, and overwriting an existing key at the cap still lands.
func TestRestApiTagBudgetBoundsMergedSet(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { st.Close() })
	store := NewRestApiStore(st, "123456789012", "us-east-1")

	api, err := store.Create(&RestApi{Name: "budget-api"})
	require.NoError(t, err)

	at := func(prefix string, n int) map[string]string {
		m := make(map[string]string, n)
		for i := 0; i < n; i++ {
			m[fmt.Sprintf("%s-%02d", prefix, i)] = "v"
		}
		return m
	}

	require.NoError(t, store.Tag(api.Id, at("a", 30)))
	require.NoError(t, store.Tag(api.Id, at("b", 20)))

	err = store.Tag(api.Id, map[string]string{"extra": "v"})
	require.Error(t, err)
	apiErr, ok := err.(*awserrors.AWSError)
	require.True(t, ok, "overflow must carry the wire identity, got %T", err)
	require.Equal(t, "BadRequestException", apiErr.Code)
	require.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)

	require.NoError(t, store.Tag(api.Id, map[string]string{"a-00": "rewritten"}))
	fetched, err := store.Get(api.Id)
	require.NoError(t, err)
	require.Len(t, fetched.Tags, 50)
}

// The documented fifty-tag cap applies to every taggable API Gateway
// resource, not the REST API alone: the stage, domain name, API key and
// usage plan merges each bound the merged set with the service's
// BadRequestException identity.
func TestTaggableFamilyTagBudgetBoundsMergedSet(t *testing.T) {
	st, err := storage.Open(t.TempDir())
	require.NoError(t, err)
	t.Cleanup(func() { st.Close() })

	at := func(prefix string, n int) map[string]string {
		m := make(map[string]string, n)
		for i := 0; i < n; i++ {
			m[fmt.Sprintf("%s-%02d", prefix, i)] = "v"
		}
		return m
	}
	expectOverflow := func(t *testing.T, err error) {
		t.Helper()
		require.Error(t, err)
		apiErr, ok := err.(*awserrors.AWSError)
		require.True(t, ok, "overflow must carry the wire identity, got %T", err)
		require.Equal(t, "BadRequestException", apiErr.Code)
		require.Equal(t, http.StatusBadRequest, apiErr.HTTPStatus)
	}

	restApi := NewRestApiStore(st, "123456789012", "us-east-1")
	api, err := restApi.Create(&RestApi{Name: "family-api"})
	require.NoError(t, err)
	_, err = restApi.CreateStage(api.Id, &Stage{StageName: "prod"})
	require.NoError(t, err)
	require.NoError(t, restApi.TagStage(api.Id, "prod", at("s", 30)))
	require.NoError(t, restApi.TagStage(api.Id, "prod", at("t", 20)))
	expectOverflow(t, restApi.TagStage(api.Id, "prod", map[string]string{"extra": "v"}))

	domains := NewDomainStore(st, "123456789012", "us-east-1")
	_, err = domains.CreateDomainName(&DomainName{DomainName: "budget.example.com"})
	require.NoError(t, err)
	require.NoError(t, domains.TagDomainName("budget.example.com", at("d", 30)))
	require.NoError(t, domains.TagDomainName("budget.example.com", at("e", 20)))
	expectOverflow(t, domains.TagDomainName("budget.example.com", map[string]string{"extra": "v"}))

	usage := NewUsageStore(st, "123456789012", "us-east-1")
	key, err := usage.CreateApiKey(&ApiKey{Name: "budget-key"})
	require.NoError(t, err)
	require.NoError(t, usage.TagApiKey(key.Id, at("k", 30)))
	require.NoError(t, usage.TagApiKey(key.Id, at("l", 20)))
	expectOverflow(t, usage.TagApiKey(key.Id, map[string]string{"extra": "v"}))

	plan, err := usage.CreateUsagePlan(&UsagePlan{Name: "budget-plan"})
	require.NoError(t, err)
	require.NoError(t, usage.TagUsagePlan(plan.Id, at("p", 30)))
	require.NoError(t, usage.TagUsagePlan(plan.Id, at("q", 20)))
	expectOverflow(t, usage.TagUsagePlan(plan.Id, map[string]string{"extra": "v"}))
}
