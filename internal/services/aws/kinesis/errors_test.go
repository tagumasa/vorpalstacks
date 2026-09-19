package kinesis

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	kinesisstore "vorpalstacks/internal/store/aws/kinesis"
)

func TestKinesisErrors(t *testing.T) {
	// The HTTP statuses follow the AWS Kinesis API reference's per-
	// operation error tables: awsJson1_1 carries every modelled client
	// error at 400, identity travelling in __type.
	t.Run("ErrResourceNotFound", func(t *testing.T) {
		assert.Equal(t, "ResourceNotFoundException: Requested resource not found.", ErrResourceNotFound.Error())
		assert.Equal(t, http.StatusBadRequest, ErrResourceNotFound.GetHTTPStatusCode())
	})

	t.Run("ErrResourceInUse", func(t *testing.T) {
		assert.Equal(t, "ResourceInUseException: The resource is in use.", ErrResourceInUse.Error())
		assert.Equal(t, http.StatusBadRequest, ErrResourceInUse.GetHTTPStatusCode())
	})

	t.Run("ErrInvalidArgument", func(t *testing.T) {
		assert.Equal(t, "InvalidArgumentException: Invalid argument.", ErrInvalidArgument.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidArgument.GetHTTPStatusCode())
	})

	t.Run("ErrLimitExceeded", func(t *testing.T) {
		assert.Equal(t, "LimitExceededException: The requested resource exceeds the maximum number allowed.", ErrLimitExceeded.Error())
		assert.Equal(t, http.StatusBadRequest, ErrLimitExceeded.GetHTTPStatusCode())
	})

	t.Run("ErrProvisionedThroughputExceeded", func(t *testing.T) {
		assert.Equal(t, "ProvisionedThroughputExceededException: Rate exceeded for this shard.", ErrProvisionedThroughputExceeded.Error())
		assert.Equal(t, http.StatusBadRequest, ErrProvisionedThroughputExceeded.GetHTTPStatusCode())
	})

	t.Run("ErrExpiredIterator", func(t *testing.T) {
		assert.Equal(t, "ExpiredIteratorException: Iterator expired.", ErrExpiredIterator.Error())
		assert.Equal(t, http.StatusBadRequest, ErrExpiredIterator.GetHTTPStatusCode())
	})

	t.Run("ErrInvalidIterator", func(t *testing.T) {
		assert.Equal(t, "InvalidArgumentException: Invalid iterator.", ErrInvalidIterator.Error())
		assert.Equal(t, http.StatusBadRequest, ErrInvalidIterator.GetHTTPStatusCode())
	})

	t.Run("ErrShardClosed", func(t *testing.T) {
		assert.Equal(t, "InvalidArgumentException: Shard is closed.", ErrShardClosed.Error())
		assert.Equal(t, http.StatusBadRequest, ErrShardClosed.GetHTTPStatusCode())
	})

	t.Run("ErrExpiredNextToken", func(t *testing.T) {
		assert.Equal(t, "ExpiredNextTokenException: The pagination token passed to the operation is expired.", ErrExpiredNextToken.Error())
		assert.Equal(t, http.StatusBadRequest, ErrExpiredNextToken.GetHTTPStatusCode())
	})

	t.Run("ErrInternalFailure", func(t *testing.T) {
		assert.Equal(t, "InternalFailureException: The processing of the request failed because of an unknown error, exception, or failure.", ErrInternalFailure.Error())
		assert.Equal(t, http.StatusInternalServerError, ErrInternalFailure.GetHTTPStatusCode())
	})
}

// TestMapStoreErrorFallback pins the three outcomes of the single
// store-error mapping site: a table-mapped sentinel answers its declared
// AWS identity (a duplicate consumer is ResourceInUseException, the shape
// RegisterStreamConsumer declares), an already wire-shaped error passes
// through unchanged, and an unmapped wrapped storage failure answers
// InternalFailureException — never a client-fault 400.
func TestMapStoreErrorFallback(t *testing.T) {
	svc := NewKinesisService("000000000000")

	if err := svc.mapStoreError(kinesisstore.ErrConsumerAlreadyExists); !errors.Is(err, ErrResourceInUse) {
		t.Fatalf("duplicate consumer: want ResourceInUseException, got: %v", err)
	}

	if err := svc.mapStoreError(ErrResourceNotFound); !errors.Is(err, ErrResourceNotFound) {
		t.Fatalf("wire-shaped error must pass through, got: %v", err)
	}

	wrapped := fmt.Errorf("failed to delete records for shard %s: %w", "shardId-000000000000", errors.New("pebble: closed"))
	if err := svc.mapStoreError(wrapped); !errors.Is(err, ErrInternalFailure) {
		t.Fatalf("unmapped infrastructure error: want InternalFailureException, got: %v", err)
	}
}
