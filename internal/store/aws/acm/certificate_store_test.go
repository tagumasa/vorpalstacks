package acm

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
)

// fakeBucket implements storage.Bucket with programmable Get behaviour.
// Every other method is unreachable in these tests and panics if called.
type fakeBucket struct {
	get func(key []byte) ([]byte, error)
}

func (b *fakeBucket) Get(key []byte) ([]byte, error) { return b.get(key) }
func (b *fakeBucket) Put(key, value []byte) error    { panic("unreachable") }
func (b *fakeBucket) Delete(key []byte) error        { panic("unreachable") }
func (b *fakeBucket) Has(key []byte) bool            { panic("unreachable") }
func (b *fakeBucket) ForEach(fn func(k, v []byte) error) error {
	panic("unreachable")
}
func (b *fakeBucket) ScanPrefix(prefix []byte) storage.Iterator { panic("unreachable") }
func (b *fakeBucket) ScanRange(start, end []byte) storage.Iterator {
	panic("unreachable")
}
func (b *fakeBucket) Count() int { panic("unreachable") }

// fakeStorage implements storage.BasicStorage serving one bucket for every
// name; the certificate and config buckets are indistinguishable here,
// which is sufficient because each test targets a single method.
type fakeStorage struct {
	bucket storage.Bucket
}

func (s fakeStorage) Close() error                 { return nil }
func (s fakeStorage) Bucket(string) storage.Bucket { return s.bucket }
func (s fakeStorage) CreateBucket(string) error    { return nil }
func (s fakeStorage) DeleteBucket(string) error    { return nil }
func (s fakeStorage) ListBuckets() []string        { return nil }

func TestGetAccountConfigurationErrorPropagation(t *testing.T) {
	tests := []struct {
		name    string
		get     func(key []byte) ([]byte, error)
		wantErr bool
		want    int
	}{
		{
			// BaseStore.Get turns a nil-value read into ErrNotFound; a
			// never-configured account must receive the AWS default.
			name:    "missing key returns the AWS default with nil error",
			get:     func([]byte) ([]byte, error) { return nil, nil },
			wantErr: false,
			want:    45,
		},
		{
			name:    "stored configuration returned as-is",
			get:     func([]byte) ([]byte, error) { return []byte(`{"ExpiryEvents":{"DaysBeforeExpiry":30}}`), nil },
			wantErr: false,
			want:    30,
		},
		{
			name:    "storage failure propagates instead of the default",
			get:     func([]byte) ([]byte, error) { return nil, errors.New("pebble: i/o error") },
			wantErr: true,
		},
	}
	for _, tt := range tests {
		store := NewCertificateStore(fakeStorage{bucket: &fakeBucket{get: tt.get}}, "123456789012", "us-east-1")
		config, err := store.GetAccountConfiguration("123456789012", "us-east-1")
		if tt.wantErr {
			if err == nil {
				t.Errorf("%s: expected error, got config %+v", tt.name, config)
			}
			continue
		}
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.name, err)
		}
		if config.ExpiryEvents.DaysBeforeExpiry != tt.want {
			t.Errorf("%s: DaysBeforeExpiry = %d, want %d", tt.name, config.ExpiryEvents.DaysBeforeExpiry, tt.want)
		}
	}
}

func TestInUseByErrorSemantics(t *testing.T) {
	certArn := "arn:aws:acm:us-east-1:123456789012:certificate/0000000000000001"
	resourceArn := "arn:aws:cloudfront::123456789012:distribution/example"
	ioFailure := errors.New("pebble: i/o error")

	ops := map[string]func(store *CertificateStore) error{
		"AddInUseBy":    func(store *CertificateStore) error { return store.AddInUseBy(certArn, resourceArn) },
		"RemoveInUseBy": func(store *CertificateStore) error { return store.RemoveInUseBy(certArn, resourceArn) },
	}

	for op, call := range ops {
		missing := NewCertificateStore(fakeStorage{bucket: &fakeBucket{get: func([]byte) ([]byte, error) { return nil, nil }}}, "123456789012", "us-east-1")
		if err := call(missing); !IsNotFound(err) {
			t.Errorf("%s on absent certificate: error should classify as not-found, got %v", op, err)
		}

		failing := NewCertificateStore(fakeStorage{bucket: &fakeBucket{get: func([]byte) ([]byte, error) { return nil, ioFailure }}}, "123456789012", "us-east-1")
		err := call(failing)
		if err == nil {
			t.Fatalf("%s on storage failure: expected error", op)
		}
		if IsNotFound(err) {
			t.Errorf("%s on storage failure: mislabelled as not-found: %v", op, err)
		}
		if !errors.Is(err, ioFailure) {
			t.Errorf("%s on storage failure: underlying cause lost: %v", op, err)
		}
	}
}
