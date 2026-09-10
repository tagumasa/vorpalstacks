package appsync

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	appsyncstore "vorpalstacks/internal/store/aws/appsync"
)

func TestMapStoreErrorEPassesTypedErrorsThrough(t *testing.T) {
	// Validation helpers return service-level AppSyncError values; the
	// mapping must not degrade them to the 500 fallback.
	ae, ok := mapStoreErrorE(NewNotFoundException("Event API")).(*AppSyncError)
	assert.True(t, ok, "typed AppSyncError must survive the mapping")
	assert.Equal(t, "NotFoundException", ae.Code)
	assert.Equal(t, http.StatusNotFound, ae.HTTPStatus)
	assert.Equal(t, "Event API not found.", ae.Message)
}

func TestMapStoreErrorEMapsStoreSentinels(t *testing.T) {
	ae, ok := mapStoreErrorE(appsyncstore.ErrGraphqlApiNotFound).(*AppSyncError)
	assert.True(t, ok)
	assert.Equal(t, "NotFoundException", ae.Code)
	assert.Equal(t, http.StatusNotFound, ae.HTTPStatus)
}

func TestMapStoreErrorEUnknownFallsBackToInternalFailure(t *testing.T) {
	ae, ok := mapStoreErrorE(errors.New("boom")).(*AppSyncError)
	assert.True(t, ok)
	assert.Equal(t, "InternalFailureException", ae.Code)
	assert.Equal(t, http.StatusInternalServerError, ae.HTTPStatus)
}

func TestGraphQLWireErrorFromStoreError(t *testing.T) {
	wireErr := graphqlWireErrorFromStoreError(appsyncstore.ErrGraphqlApiNotFound)
	assert.Equal(t, http.StatusNotFound, wireErr.HTTPStatus)
	assert.Equal(t, "NotFoundException", wireErr.ErrorType)

	internal := graphqlWireErrorFromStoreError(errors.New("boom"))
	assert.Equal(t, http.StatusInternalServerError, internal.HTTPStatus)
	assert.Equal(t, "InternalFailureException", internal.ErrorType)
}
