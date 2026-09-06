package apigateway

import (
	"errors"
	"testing"

	"vorpalstacks/internal/core/storage"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// fakeBucket implements storage.Bucket with programmable Get and Put
// behaviour; Has defaults to true so write paths pass the store's existence
// check. Every other method is unreachable in these tests and panics if
// called.
type fakeBucket struct {
	get func(key []byte) ([]byte, error)
	put func(key, value []byte) error
	has func(key []byte) bool
}

func (b *fakeBucket) Get(key []byte) ([]byte, error) { return b.get(key) }
func (b *fakeBucket) Put(key, value []byte) error {
	if b.put != nil {
		return b.put(key, value)
	}
	panic("unreachable")
}
func (b *fakeBucket) Has(key []byte) bool {
	if b.has != nil {
		return b.has(key)
	}
	return true
}
func (b *fakeBucket) Delete(key []byte) error { panic("unreachable") }
func (b *fakeBucket) ForEach(fn func(k, v []byte) error) error {
	panic("unreachable")
}
func (b *fakeBucket) ScanPrefix(prefix []byte) storage.Iterator { panic("unreachable") }
func (b *fakeBucket) ScanRange(start, end []byte) storage.Iterator {
	panic("unreachable")
}
func (b *fakeBucket) Count() int { panic("unreachable") }

// fakeStorage implements storage.BasicStorage serving one bucket for every
// name, which is sufficient because each test targets a single store.
type fakeStorage struct {
	bucket storage.Bucket
}

func (s fakeStorage) Close() error                 { return nil }
func (s fakeStorage) Bucket(string) storage.Bucket { return s.bucket }
func (s fakeStorage) CreateBucket(string) error    { return nil }
func (s fakeStorage) DeleteBucket(string) error    { return nil }
func (s fakeStorage) ListBuckets() []string        { return nil }

// TestDeleteResourceClassification pins the delete classification: the store
// sentinels for child-resources and root-resource deletions map to
// BadRequestException, never falling through to the 500 fallback.
func TestDeleteResourceClassification(t *testing.T) {
	apiJSON := []byte(`{"id":"api1","name":"probe","resources":{
		"res-root":{"id":"res-root","parent_id":"","path":"/","path_part":"/","resource_methods":{}},
		"res-a":{"id":"res-a","parent_id":"res-root","path":"/a","path_part":"a","resource_methods":{}},
		"res-b":{"id":"res-b","parent_id":"res-a","path":"/a/b","path_part":"b","resource_methods":{}}
	}}`)
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	fs := fakeStorage{bucket: &fakeBucket{get: func(key []byte) ([]byte, error) { return apiJSON, nil }}}
	stores := &apiGatewayStores{
		restApis: apigatewaystore.NewRestApiStore(fs, "123456789012", "us-east-1"),
	}

	tests := []struct {
		name       string
		resourceId string
		wantMsg    string
	}{
		{name: "resource with children is a BadRequest", resourceId: "res-a", wantMsg: "Resource has child resources"},
		{name: "root resource deletion is a BadRequest", resourceId: "res-root", wantMsg: "Cannot delete the root resource"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.deleteResourceCore(stores, "api1", tt.resourceId)
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			apiErr, ok := err.(*ApiGatewayError)
			if !ok {
				t.Fatalf("expected *ApiGatewayError, got %T: %v", err, err)
			}
			if apiErr.Code != "BadRequestException" {
				t.Errorf("code = %q, want BadRequestException", apiErr.Code)
			}
			if apiErr.HTTPStatus != 400 {
				t.Errorf("http status = %d, want 400", apiErr.HTTPStatus)
			}
			if apiErr.Message != tt.wantMsg {
				t.Errorf("message = %q, want %q", apiErr.Message, tt.wantMsg)
			}
		})
	}
}

// TestCoreErrorMapping pins the error contract of the get cores and the
// resource-delete fallback: a missing resource must surface as
// NotFoundException (404) through the store sentinels, while a genuine
// storage failure must surface as InternalFailure (500) — never as a
// blanket 404 or an ad-hoc 500 code.
func TestCoreErrorMapping(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")

	apiJSON := []byte(`{"id":"api1","name":"probe","resources":{
		"res-root":{"id":"res-root","parent_id":"","path":"/","path_part":"/","resource_methods":{}},
		"res-a":{"id":"res-a","parent_id":"res-root","path":"/a","path_part":"a","resource_methods":{}},
		"res-b":{"id":"res-b","parent_id":"res-a","path":"/a/b","path_part":"b","resource_methods":{}}
	}}`)

	newStores := func(get func(key []byte) ([]byte, error), put func(key, value []byte) error) *apiGatewayStores {
		fs := fakeStorage{bucket: &fakeBucket{get: get, put: put}}
		return &apiGatewayStores{
			restApis: apigatewaystore.NewRestApiStore(fs, "123456789012", "us-east-1"),
			usage:    apigatewaystore.NewUsageStore(fs, "123456789012", "us-east-1"),
			domains:  apigatewaystore.NewDomainStore(fs, "123456789012", "us-east-1"),
		}
	}

	tests := []struct {
		name      string
		get       func(key []byte) ([]byte, error)
		put       func(key, value []byte) error
		invoke    func(stores *apiGatewayStores) error
		wantCode  string
		wantHTTP  int
		wantMsgIn string
	}{
		{
			name: "missing usage plan key maps to NotFoundException",
			get:  func([]byte) ([]byte, error) { return nil, nil },
			invoke: func(stores *apiGatewayStores) error {
				_, err := svc.getUsagePlanKeyCore(stores, "plan1", "key1")
				return err
			},
			wantCode: "NotFoundException",
			wantHTTP: 404,
		},
		{
			name: "storage failure on usage plan key read maps to InternalFailure",
			get:  func([]byte) ([]byte, error) { return nil, errors.New("pebble: corrupted log") },
			invoke: func(stores *apiGatewayStores) error {
				_, err := svc.getUsagePlanKeyCore(stores, "plan1", "key1")
				return err
			},
			wantCode: "InternalFailure",
			wantHTTP: 500,
		},
		{
			name: "missing domain name maps to NotFoundException",
			get:  func([]byte) ([]byte, error) { return nil, nil },
			invoke: func(stores *apiGatewayStores) error {
				_, err := svc.getDomainNameCore(stores, "api.example.com", "")
				return err
			},
			wantCode: "NotFoundException",
			wantHTTP: 404,
		},
		{
			name: "storage failure on domain read maps to InternalFailure",
			get:  func([]byte) ([]byte, error) { return nil, errors.New("pebble: io error") },
			invoke: func(stores *apiGatewayStores) error {
				_, err := svc.getDomainNameCore(stores, "api.example.com", "")
				return err
			},
			wantCode: "InternalFailure",
			wantHTTP: 500,
		},
		{
			name: "storage failure on resource delete write maps to InternalFailure",
			get:  func([]byte) ([]byte, error) { return apiJSON, nil },
			put:  func(key, value []byte) error { return errors.New("pebble: sync failed") },
			invoke: func(stores *apiGatewayStores) error {
				return svc.deleteResourceCore(stores, "api1", "res-b")
			},
			wantCode: "InternalFailure",
			wantHTTP: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.invoke(newStores(tt.get, tt.put))
			if err == nil {
				t.Fatal("expected an error, got nil")
			}
			apiErr, ok := err.(*ApiGatewayError)
			if !ok {
				t.Fatalf("expected *ApiGatewayError, got %T: %v", err, err)
			}
			if apiErr.Code != tt.wantCode {
				t.Errorf("code = %q, want %q", apiErr.Code, tt.wantCode)
			}
			if apiErr.HTTPStatus != tt.wantHTTP {
				t.Errorf("http status = %d, want %d", apiErr.HTTPStatus, tt.wantHTTP)
			}
		})
	}
}
