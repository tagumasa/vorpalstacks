package appsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The API list cores validate pagination through listOptionsFromParams, so
// a whitespace-bearing token and a negative maxResults are rejected instead
// of silently normalised. The store is never reached on these paths.
func TestListApisCoreRejectsInvalidPagination(t *testing.T) {
	s := &AppSyncService{}

	_, _, err := s.listApisCore(nil, 25, "token with space")
	assert.Error(t, err)
	ae, ok := err.(*AppSyncError)
	assert.True(t, ok)
	assert.Equal(t, "BadRequestException", ae.Code)

	_, _, err = s.listApisCore(nil, -1, "")
	assert.Error(t, err, "negative maxResults must be rejected, not mapped to the default")
	ae, ok = err.(*AppSyncError)
	assert.True(t, ok)
	assert.Equal(t, "BadRequestException", ae.Code)
}
