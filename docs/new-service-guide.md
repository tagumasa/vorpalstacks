# Adding a New AWS Service

This guide explains how to add a new AWS service to Vorpalstacks. It covers every file that must be created or modified, with concrete examples drawn from existing services.

---

## Table of Contents

1. [Overview](#overview)
2. [Step 1 — Store Layer](#step-1--store-layer)
3. [Step 2 — Service Layer](#step-2--service-layer)
4. [Step 3 — Store Provider Registration](#step-3--store-provider-registration)
5. [Step 4 — Protocol Routing](#step-4--protocol-routing)
6. [Step 5 — Server Integration (main.go)](#step-5--server-integration-maingo)
7. [Step 6 — gRPC-Web Admin Handler](#step-6--grpc-web-admin-handler)
8. [Step 7 — SDK Conformance Tests](#step-7--sdk-conformance-tests)
9. [Common Patterns](#common-patterns)
10. [Checklist](#checklist)

---

## Overview

Vorpalstacks uses a two-layer architecture:

| Layer | Directory | Responsibility |
|-------|-----------|----------------|
| Store | `internal/store/aws/{service}/` | Data persistence, CRUD, ARN generation |
| Service | `internal/services/aws/{service}/` | Business logic, validation, response formatting |

Adding a new service requires changes in roughly 13 files across 6 directories. The exact set depends on the service's protocol (JSON RPC, Query, REST) and whether it has cross-service dependencies.

### File Change Summary

| # | File | Purpose |
|---|------|---------|
| 1 | `internal/store/aws/{svc}/types.go` | Data types with JSON tags |
| 2 | `internal/store/aws/{svc}/interface.go` | Store interface (for StoreProvider) |
| 3 | `internal/store/aws/{svc}/store.go` | Store implementation (embeds BaseStore) |
| 4 | `internal/store/aws/{svc}/errors.go` | Store-level errors |
| 5 | `internal/services/aws/{svc}/errors.go` | Service-level errors (embed AWSError) |
| 6 | `internal/services/aws/{svc}/helpers.go` | Parameter extraction helpers |
| 7 | `internal/services/aws/{svc}/service.go` | Service struct, RegisterHandlers |
| 8 | `internal/services/aws/{svc}/*_operations.go` | Operation handlers |
| 9 | `internal/services/aws/common/interfaces/store_provider.go` | Add GetXxxStore() |
| 10 | `internal/services/aws/common/request/context.go` | Add GetXxxStore() accessor |
| 11 | `internal/server/actionregistry/registry.go` | Query protocol action registration |
| 12 | `internal/server/http/request_extraction.go` | X-Amz-Target or path-based routing |
| 13 | `main.go` | Service instantiation and registration |

Optional files:

| # | File | Purpose |
|---|------|---------|
| 14 | `internal/services/aws/{svc}/admin_handler.go` | gRPC-Web admin handler |
| 15 | `internal/server/grpcweb/register.go` | Admin handler registration |
| 16 | `sdk-tests/testutil/{svc}.go` | Go SDK conformance tests |
| 17 | `sdk-tests/testutil/runner.go` | Test runner registration |
| 18 | `internal/config/bootstrap.go` | Feature flag for the new service |

---

## Step 1 — Store Layer

### Types

`internal/store/aws/{svc}/types.go`:

```go
package {svc}

import (
    "time"
    storecommon "vorpalstacks/internal/store/aws/common"
)

type Thing struct {
    Name       string            `json:"name"`
    Arn        string            `json:"arn"`
    Status     string            `json:"status"`
    Attributes map[string]string `json:"attributes,omitempty"`
    Tags       []storecommon.Tag `json:"tags,omitempty"`
    CreatedAt  time.Time         `json:"createdAt"`
}

func NewThing(name, arn string) *Thing {
    return &Thing{
        Name:       name,
        Arn:        arn,
        Attributes: make(map[string]string),
        CreatedAt:  time.Now().UTC(),
    }
}
```

### Interface

`internal/store/aws/{svc}/interface.go`:

Define an interface for the store. This is required for the StoreProvider pattern and allows mocking in tests.

```go
package {svc}

import storecommon "vorpalstacks/internal/store/aws/common"

type ThingStoreInterface interface {
    Get(arn string) (*Thing, error)
    Create(thing *Thing) error
    Update(thing *Thing) error
    Delete(arn string) error
    Exists(arn string) bool
    List(marker string, maxItems int) (*ThingListResult, error)
    GetTags(arn string) ([]storecommon.Tag, error)
    AddTags(arn string, tags []storecommon.Tag) error
    RemoveTags(arn string, tagKeys []string) error
}

type ThingListResult struct {
    Things      []*Thing
    NextMarker  string
    IsTruncated bool
}
```

### Store Implementation

`internal/store/aws/{svc}/store.go`:

Embed `*storecommon.BaseStore` and implement the interface. Use `ARNBuilder` for ARN generation.

```go
package {svc}

import (
    "sync"
    storecommon "vorpalstacks/internal/store/aws/common"
    "vorpalstacks/internal/core/storage"
)

type ThingStore struct {
    *storecommon.BaseStore
    tagStore   *storecommon.TagStore
    arnBuilder *storecommon.ARNBuilder
    accountID  string
    region     string
    mu         sync.RWMutex
}

func NewThingStore(store storage.BasicStorage, accountID, region string) *ThingStore {
    return &ThingStore{
        BaseStore:  storecommon.NewBaseStore(store.Bucket("things-"+region), "{svc}"),
        tagStore:   storecommon.NewTagStore(store, "{svc}"),
        arnBuilder: storecommon.NewARNBuilder(accountID, region),
        accountID:  accountID,
        region:     region,
    }
}

var _ ThingStoreInterface = (*ThingStore)(nil)

func (s *ThingStore) buildArn(name string) string {
    return s.arnBuilder.Build("{svc}", "thing/"+name)
}

func (s *ThingStore) Create(thing *Thing) error {
    thing.Arn = s.buildArn(thing.Name)
    if s.Exists(thing.Arn) {
        return ErrThingAlreadyExists
    }
    return s.Put(thing.Arn, thing)
}

func (s *ThingStore) Get(arn string) (*Thing, error) {
    var thing Thing
    if err := s.Get(arn, &thing); err != nil {
        return nil, ErrThingNotFound
    }
    return &thing, nil
}

func (s *ThingStore) Delete(arn string) error {
    s.tagStore.UntagResource(arn, nil)
    return s.BaseStore.Delete(arn)
}

func (s *ThingStore) GetTags(arn string) ([]storecommon.Tag, error) {
    return s.tagStore.ListTagsAsSlice(arn)
}

func (s *ThingStore) AddTags(arn string, tags []storecommon.Tag) error {
    tagMap := make(map[string]string)
    for _, t := range tags {
        tagMap[t.Key] = t.Value
    }
    return s.tagStore.TagResource(arn, tagMap)
}

func (s *ThingStore) RemoveTags(arn string, tagKeys []string) error {
    return s.tagStore.UntagResource(arn, tagKeys)
}
```

### Store Errors

`internal/store/aws/{svc}/errors.go`:

```go
package {svc}

import "errors"

var (
    ErrThingNotFound      = errors.New("thing not found")
    ErrThingAlreadyExists = errors.New("thing already exists")
)
```

### BaseStore Methods Reference

| Method | Signature | Description |
|--------|-----------|-------------|
| `Get` | `(key string, dest interface{}) error` | JSON unmarshal |
| `Put` | `(key string, data interface{}) error` | JSON marshal |
| `Delete` | `(key string) error` | Delete key |
| `Exists` | `(key string) bool` | Check existence |
| `ForEach` | `(fn func(key string, value []byte) error) error` | Iterate all items |
| `ScanPrefix` | `(prefix string, fn func(key string, value []byte) error) error` | Iterate by prefix |
| `DeleteByPrefix` | `(prefix string) error` | Delete all with prefix |
| `ScanRange` | `(start, end string, fn func(key, value []byte) error) error` | Iterate by key range |
| `Count` | `() int` | Count items |

Use `ScanPrefix` over `ForEach` when possible — it uses the underlying storage engine's prefix scan and avoids iterating unrelated keys. The generic `ListProto` helper uses `ForEach` internally, so for large datasets prefer a direct `ScanPrefix` loop.

---

## Step 2 — Service Layer

### Errors

`internal/services/aws/{svc}/errors.go`:

Service errors must embed `*awserrors.AWSError` so the dispatcher can extract the correct HTTP status code and error format. Implement the `GetAWSError() *awserrors.AWSError` method to satisfy the `awsErrorHolder` interface.

```go
package {svc}

import (
    "fmt"
    "vorpalstacks/internal/services/aws/common/errors"
)

type SvcError struct {
    *errors.AWSError
}

func (e *SvcError) GetAWSError() *errors.AWSError {
    return e.AWSError
}

var (
    ErrInvalidParameter = &SvcError{errors.NewValidationException("invalid parameter")}
    ErrThingNotFound    = &SvcError{errors.NewResourceNotFoundException("Thing", "")}
)

func NewThingNotFoundException(name string) *SvcError {
    return &SvcError{errors.NewResourceNotFoundException("Thing", name)}
}

func NewInvalidParameterException(msg string) *SvcError {
    return &SvcError{errors.NewValidationException(msg)}
}
```

**Important**: The `awsErrorHolder` interface is defined inline in the dispatcher's error extraction chain (`internal/server/dispatcher/response.go`). If an error does not implement `GetAWSError()`, the dispatcher falls back to a generic 500 Internal Server Error, which AWS SDKs treat as retryable — potentially causing infinite retry loops.

### Helpers

`internal/services/aws/{svc}/helpers.go`:

```go
package {svc}

import "vorpalstacks/internal/services/aws/common/request"

func getParam(req *request.ParsedRequest, key string) string {
    return request.GetParamLowerFirst(req.Parameters, key)
}

func getIntParam(req *request.ParsedRequest, key string) int {
    return request.GetIntParam(req.Parameters, key)
}

func getBoolParam(req *request.ParsedRequest, key string) bool {
    return request.GetBoolParam(req.Parameters, key)
}
```

For Query protocol services (SQS, IAM), use `request.GetParamCaseInsensitive` instead of `GetParamLowerFirst` — Query protocol uses PascalCase parameters.

### Service

`internal/services/aws/{svc}/service.go`:

```go
package {svc}

import (
    "vorpalstacks/internal/core/storage"
    "vorpalstacks/internal/server/dispatcher"
    "vorpalstacks/internal/services/aws/common/request"
    svcstore "vorpalstacks/internal/store/aws/{svc}"
)

type thingStores struct {
    things     svcstore.ThingStoreInterface
    arnBuilder *storecommon.ARNBuilder
}

type SvcService struct {
    accountID string
}

func NewSvcService(store storage.BasicStorage, accountID, region string) *SvcService {
    return &SvcService{accountID: accountID}
}

func (s *SvcService) store(reqCtx *request.RequestContext) (*thingStores, error) {
    if thingStore := reqCtx.GetThingStore(); thingStore != nil {
        return &thingStores{
            things:     thingStore,
            arnBuilder: svcstore.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()),
        }, nil
    }
    storage, err := reqCtx.GetStorage()
    if err != nil {
        return nil, err
    }
    return &thingStores{
        things:     svcstore.NewThingStore(storage, reqCtx.GetAccountID(), reqCtx.GetRegion()),
        arnBuilder: svcstore.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()),
    }, nil
}

func (s *SvcService) RegisterHandlers(d *dispatcher.Dispatcher) {
    d.RegisterHandlerForService("{svc}", "CreateThing", s.CreateThing)
    d.RegisterHandlerForService("{svc}", "GetThing", s.GetThing)
    d.RegisterHandlerForService("{svc}", "DeleteThing", s.DeleteThing)
    d.RegisterHandlerForService("{svc}", "ListThings", s.ListThings)
    d.RegisterHandlerForService("{svc}", "TagResource", s.TagResource)
    d.RegisterHandlerForService("{svc}", "UntagResource", s.UntagResource)
    d.RegisterHandlerForService("{svc}", "ListTagsForResource", s.ListTagsForResource)
}
```

**Key points**:

- Always use `RegisterHandlerForService` (not `RegisterHandler`) to avoid operation name collisions across services.
- The `store(reqCtx)` method follows the StoreProvider-first pattern: it checks `reqCtx.GetXxxStore()` first, then falls back to creating a store from raw storage.
- For regional services that need per-region store caching, use `sync.Map` + `storecommon.GetOrCreateStoreE` (see Scheduler pattern in main.go).

### Operations

`internal/services/aws/{svc}/thing_operations.go`:

```go
package {svc}

import (
    "vorpalstacks/internal/services/aws/common/pagination"
    "vorpalstacks/internal/services/aws/common/request"
)

func (s *SvcService) CreateThing(req *request.ParsedRequest) (interface{}, error) {
    name := getParam(req, "Name")
    if name == "" {
        return nil, NewInvalidParameterException("Name is required")
    }

    stores, err := s.store(req.Ctx)
    if err != nil {
        return nil, err
    }

    thing := svcstore.NewThing(name, "")
    if err := stores.things.Create(thing); err != nil {
        if errors.Is(err, svcstore.ErrThingAlreadyExists) {
            return nil, NewThingAlreadyExistsException(name)
        }
        return nil, err
    }

    return map[string]interface{}{
        "ThingArn":  thing.Arn,
        "ThingName": thing.Name,
    }, nil
}

func (s *SvcService) ListThings(req *request.ParsedRequest) (interface{}, error) {
    stores, err := s.store(req.Ctx)
    if err != nil {
        return nil, err
    }

    maxItems := pagination.GetMaxItems(req.Parameters, 100)
    marker := pagination.GetMarker(req.Parameters)

    result, err := stores.things.List(marker, maxItems)
    if err != nil {
        return nil, err
    }

    things := make([]map[string]interface{}, 0, len(result.Things))
    for _, t := range result.Things {
        things = append(things, map[string]interface{}{
            "ThingArn":  t.Arn,
            "ThingName": t.Name,
        })
    }

    return pagination.BuildListResponse("Things", things, result.NextMarker), nil
}
```

### IAM Role Validation

If the service needs to validate IAM roles (e.g., Lambda, Step Functions, EventBridge, Scheduler, CloudTrail), use `ValidateRoleForServiceWithErrors` with service-specific error factories:

```go
validator := reqCtx.GetIAMValidator()
if validator == nil {
    return nil, NewAccessDeniedException("IAM not available")
}

errFactories := &iam.RoleErrorFactories{
    InvalidArnError:         func(arn string) error { return NewInvalidArnException(arn) },
    RoleNotFoundError:       func(arn string) error { return NewRoleNotFoundException(arn) },
    RoleCannotBeAssumedError: func(arn string) error { return NewAccessDeniedException("role cannot be assumed") },
}

if err := validator.ValidateRoleForServiceWithErrors(
    req.Ctx, roleArn, iam.ServicePrincipal{Service: "{svc}.amazonaws.com"}, errFactories,
); err != nil {
    return nil, err
}
```

**Do not** use the default `ValidateRoleForService` — it uses Lambda-specific error messages.

---

## Step 3 — Store Provider Registration

The StoreProvider is a central registry that gives the RequestContext typed access to all service stores. Every new service must register here.

### Interface

Add to `internal/services/aws/common/interfaces/store_provider.go`:

```go
// In the StoreProvider interface:
GetXxxStore() xxxstore.XxxStoreInterface

// In StoreProviderImpl struct:
xxxStore xxxstore.XxxStoreInterface

// In NewStoreProvider parameters:
xxxStore xxxstore.XxxStoreInterface,

// In NewStoreProvider return:
xxxStore: xxxStore,

// Getter method:
func (s *StoreProviderImpl) GetXxxStore() xxxstore.XxxStoreInterface { return s.xxxStore }
```

### RequestContext

Add to `internal/services/aws/common/request/context.go`:

```go
func (c *RequestContext) GetXxxStore() xxxstore.XxxStoreInterface {
    if c.storeProvider == nil {
        return nil
    }
    return c.storeProvider.GetXxxStore()
}
```

### Server Construction

The StoreProvider is constructed in `internal/server/http/server.go` (~line 152). Add the new store instance to the `NewStoreProvider(...)` call. The store instance is typically created in `main.go` and passed to the server constructor.

---

## Step 4 — Protocol Routing

Vorpalstacks auto-detects the protocol from the request. The routing priority is:

1. Path-based matching (REST services with custom paths)
2. `X-Amz-Target` header → JSON RPC protocol
3. `Action` parameter → Query protocol
4. Hostname subdomain

### JSON RPC Protocol (most common)

For services that use `X-Amz-Target` (Lambda, DynamoDB, SNS, EventBridge, Step Functions, ACM, KMS, etc.), add the target prefix mapping in `internal/server/http/request_extraction.go` in the `extractServiceFromTarget` function.

The `X-Amz-Target` format is `{ServicePrefix}.{OperationName}`. For example:
- `AWSLambda.CreateFunction` → Lambda
- `DynamoDB_20120810.CreateTable` → DynamoDB
- `AWSCertificateManagerRequestCertificate` → ACM

### Query Protocol

For services using Query protocol (SQS, IAM, SNS, STS), register all action names in `internal/server/actionregistry/registry.go`:

```go
// In initDefaults():
r.Register("{svc}", []string{
    "CreateThing",
    "GetThing",
    "DeleteThing",
    "ListThings",
    "TagResource",
    "UntagResource",
    "ListTagsForResource",
})
```

### REST Path-Based Services

For services with explicit HTTP paths (S3, Route53, CloudFront, API Gateway, Scheduler, SESv2), add a path detection function in `internal/server/http/request_extraction.go`:

```go
func is{Svc}Path(path string) bool {
    return strings.HasPrefix(path, "/v1/things/") ||
           strings.HasPrefix(path, "/2024-01-01/things/")
}
```

Then add the check in the service extraction chain.

---

## Step 5 — Server Integration (main.go)

### Feature Flag

Add to `internal/config/bootstrap.go`:

```go
// In BootstrapConfig struct:
Xxx bool

// In LoadBootstrapConfig():
Xxx: envBool("XXX_ENABLED", true),
```

### Service Registration

Add to `main.go`:

```go
import svcxxx "vorpalstacks/internal/services/aws/{svc}"

// Pattern A: Always-on service (no feature flag)
xxxService := svcxxx.NewSvcService(server.Storage(), cfg.AccountID, cfg.Region)
xxxService.RegisterHandlers(server.Dispatcher())

// Pattern B: Feature-flagged service
if cfg.Xxx {
    xxxService := svcxxx.NewSvcService(server.Storage(), cfg.AccountID, cfg.Region)
    xxxService.RegisterHandlers(server.Dispatcher())
}

// Pattern C: With cross-service dependencies
if cfg.Xxx {
    xxxService := svcxxx.NewSvcService(server.StorageManager(), cfg.AccountID)
    if lambdaService != nil {
        xxxService.SetLambdaInvoker(lambdaService)
    }
    if sqsStoreInstance != nil {
        xxxService.SetSQSStore(sqsStoreInstance)
    }
    xxxService.RegisterHandlers(server.Dispatcher())
}

// Pattern D: With background goroutine
if cfg.Xxx {
    xxxService := svcxxx.NewSvcService(server.StorageManager(), cfg.AccountID)
    xxxService.RegisterHandlers(server.Dispatcher())
    xxxService.Start()
    server.RegisterShutdownHook(func(ctx context.Context) {
        xxxService.Stop()
    })
}
```

---

## Step 6 — gRPC-Web Admin Handler

Create `internal/services/aws/{svc}/admin_handler.go`:

```go
package {svc}

import (
    "errors"
    "connectrpc.com/connect"
    pb "vorpalstacks/proto/gen/go/proto/{svc}"
    pbconnect "vorpalstacks/proto/gen/go/proto/{svc}/{svc}connect"
)

type AdminHandler struct {
    pbconnect.Unimplemented{Svc}ServiceHandler
}

var _ pbconnect.{Svc}ServiceHandler = (*AdminHandler)(nil)

func NewAdminHandler() *AdminHandler {
    return &AdminHandler{}
}

func (h *AdminHandler) CreateThing(ctx context.Context, req *connect.Request[pb.CreateThingRequest]) (*connect.Response[pb.CreateThingResponse], error) {
    return nil, connect.NewError(connect.CodeUnimplemented, errors.New("not implemented"))
}
```

Register in `internal/server/grpcweb/register.go`:

```go
s.Handle(pbconnect.New{Svc}ServiceHandler(svcxxx.NewAdminHandler()))
```

For admin handlers that need storage access, accept a `*storage.RegionStorageManager` in the constructor and extract the region from the request header (`req.Header().Get("X-Aws-Region")`).

---

## Step 7 — SDK Conformance Tests

Create `sdk-tests/testutil/{svc}.go`:

```go
package testutil

import (
    "context"
    "fmt"
    svc "vorpalstacks/internal/services/aws/{svc}/client"
)

func (r *TestRunner) Run{Svc}Tests() []TestResult {
    var results []TestResult
    cfg, _ := config.LoadDefaultAWSConfig(config.AWSConfig{Endpoint: r.endpoint, Region: r.region})
    client := svc.NewFromConfig(cfg)
    ctx := context.Background()

    results = append(results, r.RunTest("{svc}", "CreateThing", func() error {
        resp, err := client.CreateThing(ctx, &svc.CreateThingInput{Name: aws.String("test-thing")})
        if err != nil { return err }
        if resp.ThingArn == nil { return fmt.Errorf("nil ThingArn") }
        return nil
    }))

    results = append(results, r.RunTest("{svc}", "GetThing", func() error {
        _, err := client.GetThing(ctx, &svc.GetThingInput{ThingName: aws.String("test-thing")})
        return err
    }))

    return results
}
```

Register in `sdk-tests/testutil/runner.go`:

```go
// In GetAllServices():
"{svc}",

// In RunServiceTests() switch:
case "{svc}":
    results = append(results, r.Run{Svc}Tests()...)
```

---

## Common Patterns

### Pagination

```go
import "vorpalstacks/internal/services/aws/common/pagination"

maxItems := pagination.GetMaxItems(req.Parameters, 100)
marker := pagination.GetMarker(req.Parameters)

result, err := store.List(marker, maxItems)
// ...

nextToken := ""
if result.IsTruncated {
    nextToken = result.NextMarker
}
return pagination.BuildListResponse("Things", things, nextToken), nil
```

Constants: `pagination.DefaultMaxItems = 100`, `pagination.AbsoluteMaxItems = 1000`.

### Tag Operations

```go
import "vorpalstacks/internal/services/aws/common/request"

tags := request.ParseTags(req.Parameters, "Tags")

func (s *SvcService) TagResource(req *request.ParsedRequest) (interface{}, error) {
    arn := getParam(req, "ResourceArn")
    tags := request.ParseTags(req.Parameters, "Tags")
    tagSlice := make([]storecommon.Tag, 0, len(tags))
    for k, v := range tags {
        tagSlice = append(tagSlice, storecommon.Tag{Key: k, Value: v})
    }
    if err := store.AddTags(arn, tagSlice); err != nil {
        return nil, err
    }
    return map[string]interface{}{}, nil
}
```

### Response Formatting

Service handlers return `(interface{}, error)`. The response can be:
- A `map[string]interface{}` for JSON protocol
- A typed struct (will be serialised via JSON)
- `nil` for operations with no response body

### Timestamp Handling

The `ConvertTimestampsToSeconds` function in `internal/services/aws/common/protocol/encode_json.go` automatically converts millisecond-precision epoch values (> 1e12) to seconds for known timestamp keys. If your service has custom timestamp fields, add them to the `jsonTimestampKeys` map.

### Content-Type Preservation

`EncodeJSONResponse` sets `Content-Type: application/json` only if no Content-Type has already been set. Services that need protocol-specific content types (e.g., `application/x-amz-json-1.0`) should set it via `SetProtocolHeaders` before calling the response writer.

### Error Response Formats

The dispatcher auto-selects the error format based on the request protocol:
- JSON RPC → JSON error body + `X-Amzn-ErrorType` header
- Query → XML error body
- REST-XML → XML error body
- CBOR → CBOR error body

Service errors only need to embed `*errors.AWSError` with the correct code and HTTP status. The dispatcher handles serialisation.

---

## Checklist

### Store Layer
- [ ] `types.go` — Data types with JSON tags
- [ ] `interface.go` — Store interface + compile-time check (`var _ Interface = (*Store)(nil)`)
- [ ] `store.go` — Embed BaseStore, use ARNBuilder
- [ ] `errors.go` — Store-level errors
- [ ] Use `ScanPrefix` over `ForEach` for prefix queries
- [ ] Initialise nil maps before use

### Service Layer
- [ ] `errors.go` — Embed `*errors.AWSError`, implement `GetAWSError()`
- [ ] `helpers.go` — Parameter extraction (JSON: `GetParamLowerFirst`, Query: `GetParamCaseInsensitive`)
- [ ] `service.go` — Service struct, lazy `store(reqCtx)`, `RegisterHandlers`
- [ ] `*_operations.go` — Operation handlers
- [ ] Use `RegisterHandlerForService` (never bare `RegisterHandler`)
- [ ] Use `ValidateRoleForServiceWithErrors` with service-specific error factories

### Integration
- [ ] `store_provider.go` — Interface, impl struct, constructor param, getter
- [ ] `request/context.go` — `GetXxxStore()` accessor
- [ ] `actionregistry/registry.go` — Query actions (if Query protocol)
- [ ] `request_extraction.go` — X-Amz-Target or path mapping
- [ ] `config/bootstrap.go` — Feature flag + PrintStartupBanner
- [ ] `main.go` — Import, construct, RegisterHandlers, shutdown hook

### Admin Handler
- [ ] `admin_handler.go` — Embed Unimplemented handler, compile-time check
- [ ] `grpcweb/register.go` — Register handler

### Testing
- [ ] `sdk-tests/testutil/{svc}.go` — Conformance tests
- [ ] `sdk-tests/testutil/runner.go` — Register in GetAllServices + switch

### Code Standards
- [ ] British English for all comments and documentation
- [ ] No unnecessary comments
- [ ] Use common utilities from `internal/services/aws/common/`
