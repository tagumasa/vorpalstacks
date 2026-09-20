package sns

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/common/errors"
	snsstore "vorpalstacks/internal/store/aws/sns"
)

func TestSNSErrors(t *testing.T) {
	t.Run("predefined errors", func(t *testing.T) {
		assert.Equal(t, "AuthorizationError: Authorization error", ErrAuthorizationError.Error())
		assert.Equal(t, 403, ErrAuthorizationError.GetHTTPStatusCode())

		assert.Equal(t, "FilterPolicyLimitExceeded: Filter policy limit exceeded", ErrFilterLimitExceeded.Error())
		assert.Equal(t, 403, ErrFilterLimitExceeded.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Topic does not exist", ErrTopicNotFound.Error())
		assert.Equal(t, 404, ErrTopicNotFound.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Subscription does not exist", ErrSubscriptionNotFound.Error())
		assert.Equal(t, 404, ErrSubscriptionNotFound.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Platform application does not exist", ErrPlatformAppNotFound.Error())
		assert.Equal(t, 404, ErrPlatformAppNotFound.GetHTTPStatusCode())

		assert.Equal(t, "NotFound: Endpoint does not exist", ErrEndpointNotFound.Error())
		assert.Equal(t, 404, ErrEndpointNotFound.GetHTTPStatusCode())

		// The tag family reports the ResourceNotFoundException wire code,
		// distinct from the NotFoundException "NotFound" code above.
		assert.Equal(t, "ResourceNotFound: Topic does not exist", ErrResourceNotFound.Error())
		assert.Equal(t, 404, ErrResourceNotFound.GetHTTPStatusCode())

		// The message is the model's own TagLimitExceededException documentation.
		assert.Equal(t, "TagLimitExceeded: Can't add more than 50 tags to a topic.", ErrTagLimitExceeded.Error())
		assert.Equal(t, 400, ErrTagLimitExceeded.GetHTTPStatusCode())

		assert.Equal(t, "BatchEntryIdsNotDistinct: Two or more batch entries have the same ID", ErrBatchEntryIdsNotDistinct.Error())
		assert.Equal(t, 400, ErrBatchEntryIdsNotDistinct.GetHTTPStatusCode())

		assert.Equal(t, "TooManyEntriesInBatchRequest: Maximum number of entries per request are 10", ErrTooManyEntriesInBatch.Error())
		assert.Equal(t, 400, ErrTooManyEntriesInBatch.GetHTTPStatusCode())
	})

	t.Run("NewInvalidParameter carries the model's query wire code", func(t *testing.T) {
		err := NewInvalidParameter("invalid parameter value")
		assert.Equal(t, "InvalidParameterException: invalid parameter value", err.Error())
		assert.Equal(t, 400, err.GetHTTPStatusCode())
		// The query protocols report the awsQueryError trait code on the
		// wire ("InvalidParameter"), not the shape name.
		assert.Equal(t, "InvalidParameter", err.QueryErrorCode)
	})

	t.Run("mapStoreError translates every store sentinel", func(t *testing.T) {
		assert.Nil(t, mapStoreError(nil))

		cases := []struct {
			storeErr   error
			wireCode   string
			queryCode  string
			message    string
			httpStatus int
		}{
			{snsstore.ErrTopicNotFound, "NotFound", "", "Topic does not exist", 404},
			{snsstore.ErrSubscriptionNotFound, "NotFound", "", "Subscription does not exist", 404},
			{snsstore.ErrPlatformApplicationNotFound, "NotFound", "", "Platform application does not exist", 404},
			{snsstore.ErrEndpointNotFound, "NotFound", "", "Endpoint does not exist", 404},
			{snsstore.ErrPlatformApplicationAlreadyExists, "InvalidParameterException", "InvalidParameter", "Platform application already exists with the same name", 400},
			{snsstore.ErrInvalidToken, "InvalidParameterException", "InvalidParameter", "Invalid parameter: Token", 400},
		}
		for _, tc := range cases {
			mapped, ok := mapStoreError(tc.storeErr).(*awserrors.AWSError)
			assert.True(t, ok, "store error %v must map to an AWSError", tc.storeErr)
			assert.Equal(t, tc.wireCode, mapped.GetCode(), "code for %v", tc.storeErr)
			assert.Equal(t, tc.queryCode, mapped.QueryErrorCode, "query wire code for %v", tc.storeErr)
			assert.Equal(t, tc.message, mapped.GetMessage(), "message for %v", tc.storeErr)
			assert.Equal(t, tc.httpStatus, mapped.GetHTTPStatusCode(), "status for %v", tc.storeErr)
		}

		// Wrapped sentinels still translate; unknown errors pass through
		// unchanged so storage failures surface as themselves.
		assert.Equal(t, ErrTopicNotFound, mapStoreError(fmt.Errorf("get topic: %w", snsstore.ErrTopicNotFound)))
		foreign := errors.New("storage failure")
		assert.Equal(t, foreign, mapStoreError(foreign))
	})
}

// errTokenScanStorageFailure stands in for a failing token scan (a record
// the scanner cannot decode); the stub answers nothing else, so the test
// fails loudly if the core calls any other store method first.
type errTokenScanStorageFailure struct{}

func (errTokenScanStorageFailure) Error() string { return "stub: token scan storage failure" }

type tokenScanFailingStore struct {
	snsstore.SNSStoreInterface
}

func (t *tokenScanFailingStore) FindSubscriptionByToken(topicArn, token string) (*snsstore.Subscription, error) {
	return nil, errTokenScanStorageFailure{}
}

// TestConfirmSubscriptionTokenScanStorageFailureIsInternal pins the error
// plane of a failing token lookup: the miss alone answers InvalidParameter
// — a storage failure must not wear that client-error identity but pass
// through the seam and surface as the internal error it is.
func TestConfirmSubscriptionTokenScanStorageFailureIsInternal(t *testing.T) {
	svc := NewSNSService(nil, "123456789012", "us-east-1")
	want := errTokenScanStorageFailure{}
	_, err := svc.confirmSubscriptionCore(&tokenScanFailingStore{}, ConfirmSubscriptionInput{
		TopicArn: "arn:aws:sns:us-east-1:123456789012:token-scan",
		Token:    "token-x",
	})
	if !errors.Is(err, want) {
		t.Fatalf("confirmSubscriptionCore over a failing scan: err = %v, want the passthrough storage failure %v", err, want)
	}
}
