package iam

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/storage"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// faultBucket lets a test choose the failure mode of every storage read:
// a non-nil ioErr simulates an infrastructure fault, while a nil ioErr
// (returning no data) makes BaseStore.Get signal the common not-found
// sentinel.
type faultBucket struct {
	ioErr error
}

func (b *faultBucket) Get(key []byte) ([]byte, error)           { return nil, b.ioErr }
func (b *faultBucket) Put(key, value []byte) error              { return b.ioErr }
func (b *faultBucket) Delete(key []byte) error                  { return b.ioErr }
func (b *faultBucket) Has(key []byte) bool                      { return false }
func (b *faultBucket) ForEach(fn func(k, v []byte) error) error { return b.ioErr }
func (b *faultBucket) ScanPrefix(prefix []byte) storage.Iterator {
	return emptyIterator{}
}
func (b *faultBucket) ScanPrefixReverse(prefix, before []byte) storage.Iterator {
	return emptyIterator{}
}
func (b *faultBucket) ScanRange(start, end []byte) storage.Iterator {
	return emptyIterator{}
}
func (b *faultBucket) Count() int { return 0 }

type emptyIterator struct{}

func (emptyIterator) Next() bool    { return false }
func (emptyIterator) Key() []byte   { return nil }
func (emptyIterator) Value() []byte { return nil }
func (emptyIterator) Error() error  { return nil }
func (emptyIterator) Close()        {}

type faultStorage struct{ bucket *faultBucket }

func (s faultStorage) Close() error                      { return nil }
func (s faultStorage) Bucket(name string) storage.Bucket { return s.bucket }
func (s faultStorage) CreateBucket(name string) error    { return nil }
func (s faultStorage) DeleteBucket(name string) error    { return nil }

// A storage-layer failure must surface as a 5xx InternalFailure carrying
// the cause, never as NoSuchEntity; an entity that is genuinely absent
// still maps to the site's NoSuchEntity constructor, through both the
// common sentinel and the per-family store sentinels.
func TestStoreReadFailureIsNotNoSuchEntity(t *testing.T) {
	s := &IAMService{}

	t.Run("infrastructure fault surfaces as 5xx", func(t *testing.T) {
		store := iamstore.NewIAMStore(faultStorage{&faultBucket{
			ioErr: errors.New("simulated storage I/O failure"),
		}}, "123456789012")

		_, err := s.getUserCore(store, "someuser")
		if err == nil {
			t.Fatal("expected an error from a failing storage layer")
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("expected *awserrors.AWSError, got %T", err)
		}
		assert.Equal(t, http.StatusInternalServerError, awsErr.GetHTTPStatusCode())
		assert.Contains(t, awsErr.Error(), "InternalFailure")
		assert.Contains(t, awsErr.Error(), "simulated storage I/O failure")
		assert.NotContains(t, awsErr.Error(), "NoSuchEntity")
	})

	t.Run("absent entity still maps to NoSuchEntity", func(t *testing.T) {
		store := iamstore.NewIAMStore(faultStorage{&faultBucket{}}, "123456789012")

		_, err := s.getUserCore(store, "someuser")
		if err == nil {
			t.Fatal("expected NoSuchEntity for an absent user")
		}
		awsErr, ok := err.(*awserrors.AWSError)
		if !ok {
			t.Fatalf("expected *awserrors.AWSError, got %T", err)
		}
		assert.Equal(t, http.StatusNotFound, awsErr.GetHTTPStatusCode())
		assert.Equal(t, "NoSuchEntity: The user with name someuser cannot be found.", awsErr.Error())
	})

	t.Run("family sentinel maps to the site's not-found constructor", func(t *testing.T) {
		wrapped := iamstore.NewStoreError("get_policy", iamstore.ErrPolicyNotFound)
		mapped := storeReadError(wrapped, iamstore.ErrPolicyNotFound,
			NewNoSuchPolicyError("arn:aws:iam::123456789012:policy/test"))
		assert.Equal(t, http.StatusNotFound, mapped.GetHTTPStatusCode())
		assert.Equal(t, "NoSuchEntity: The policy with name arn:aws:iam::123456789012:policy/test cannot be found.", mapped.Error())
	})
}
