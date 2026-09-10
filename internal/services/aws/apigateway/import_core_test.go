package apigateway

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"testing"

	"vorpalstacks/internal/core/storage"
	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// memBucket is an in-memory storage.Bucket: writes persist so a store
// built on it round-trips records through the real BaseStore logic.
type memBucket struct {
	data map[string][]byte
}

func newMemBucket() *memBucket { return &memBucket{data: make(map[string][]byte)} }

func (b *memBucket) Get(key []byte) ([]byte, error) {
	if v, ok := b.data[string(key)]; ok {
		return v, nil
	}
	return nil, nil
}

func (b *memBucket) Put(key, value []byte) error {
	encoded := make([]byte, len(value))
	copy(encoded, value)
	b.data[string(key)] = encoded
	return nil
}

func (b *memBucket) Delete(key []byte) error {
	delete(b.data, string(key))
	return nil
}

func (b *memBucket) Has(key []byte) bool {
	_, ok := b.data[string(key)]
	return ok
}

func (b *memBucket) ForEach(fn func(k, v []byte) error) error {
	keys := make([]string, 0, len(b.data))
	for k := range b.data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if err := fn([]byte(k), b.data[k]); err != nil {
			return err
		}
	}
	return nil
}

func (b *memBucket) ScanPrefix(prefix []byte) storage.Iterator {
	keys := make([]string, 0, len(b.data))
	for k := range b.data {
		if strings.HasPrefix(k, string(prefix)) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return &memIterator{keys: keys, idx: -1, bucket: b}
}

func (b *memBucket) ScanPrefixReverse(prefix, before []byte) storage.Iterator {
	keys := make([]string, 0, len(b.data))
	for k := range b.data {
		if strings.HasPrefix(k, string(prefix)) && (before == nil || k < string(before)) {
			keys = append(keys, k)
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return &memIterator{keys: keys, idx: -1, bucket: b}
}

func (b *memBucket) ScanRange(start, end []byte) storage.Iterator {
	keys := make([]string, 0, len(b.data))
	for k := range b.data {
		if k >= string(start) && k < string(end) {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	return &memIterator{keys: keys, idx: -1, bucket: b}
}

func (b *memBucket) Count() int { return len(b.data) }

type memIterator struct {
	keys   []string
	idx    int
	bucket *memBucket
}

func (it *memIterator) Next() bool {
	it.idx++
	return it.idx < len(it.keys)
}

func (it *memIterator) Key() []byte   { return []byte(it.keys[it.idx]) }
func (it *memIterator) Value() []byte { return it.bucket.data[it.keys[it.idx]] }
func (it *memIterator) Error() error  { return nil }
func (it *memIterator) Close()        {}

// memStorage serves one in-memory bucket for every name.
type memStorage struct{ bucket *memBucket }

func (s memStorage) Close() error                 { return nil }
func (s memStorage) Bucket(string) storage.Bucket { return s.bucket }
func (s memStorage) CreateBucket(string) error    { return nil }
func (s memStorage) DeleteBucket(string) error    { return nil }
func (s memStorage) ListBuckets() []string        { return nil }

// failingForEachBucket wraps memBucket with a switch that makes the
// record listing — the read path GetApiKeyByValue runs — fail
// deterministically once enabled.
type failingForEachBucket struct {
	*memBucket
	fail bool
}

func (b *failingForEachBucket) ForEach(fn func(k, v []byte) error) error {
	if b.fail {
		return errors.New("simulated storage read failure")
	}
	return b.memBucket.ForEach(fn)
}

// failingStorage serves the failing bucket for every name.
type failingStorage struct{ bucket *failingForEachBucket }

func (s failingStorage) Close() error                 { return nil }
func (s failingStorage) Bucket(string) storage.Bucket { return s.bucket }
func (s failingStorage) CreateBucket(string) error    { return nil }
func (s failingStorage) DeleteBucket(string) error    { return nil }
func (s failingStorage) ListBuckets() []string        { return nil }

func newImportStores(t *testing.T) *apiGatewayStores {
	t.Helper()
	fs := memStorage{bucket: newMemBucket()}
	return &apiGatewayStores{
		restApis: apigatewaystore.NewRestApiStore(fs, "123456789012", "us-east-1"),
		usage:    apigatewaystore.NewUsageStore(fs, "123456789012", "us-east-1"),
		account:  apigatewaystore.NewAccountStore(fs, "123456789012", "us-east-1"),
	}
}

// swaggerDefinition is a Swagger 2.0 import document exercising the
// metadata, paths, definitions, security and integration surfaces.
const swaggerDefinition = `{
  "swagger": "2.0",
  "info": {"title": "imported-api", "version": "1.0.0", "description": "an imported API"},
  "basePath": "/team/v1",
  "securityDefinitions": {
    "api_key": {"type": "apiKey", "name": "x-api-key", "in": "header"},
    "my_auth": {
      "type": "apiKey", "name": "Authorization", "in": "header",
      "x-amazon-apigateway-authorizer": {
        "type": "token",
        "authorizerUri": "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:fn/invocations",
        "identitySource": "method.request.header.Authorization",
        "authorizerResultTtlInSeconds": 30
      }
    },
    "pools": {
      "type": "apiKey", "name": "Authorization", "in": "header",
      "x-amazon-apigateway-authorizer": {
        "type": "cognito_user_pools",
        "providerARNs": ["arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_ABC"]
      }
    }
  },
  "paths": {
    "/": {"get": {"operationId": "ListRoot", "responses": {"200": {"description": "ok"}}}},
    "/orders/{id}": {
      "get": {
        "operationId": "GetOrder",
        "security": [{"my_auth": []}],
        "x-amazon-apigateway-api-key-required": true,
        "responses": {"200": {"description": "ok"}, "404": {"description": "missing"}},
        "x-amazon-apigateway-integration": {"type": "AWS_PROXY", "httpMethod": "POST", "uri": "arn:aws:apigateway:us-east-1:lambda:path/2015-03-31/functions/arn:order/invocations", "passthroughBehavior": "when_no_match"}
      },
      "delete": {"operationId": "DeleteOrder", "responses": {"204": {"description": "gone"}}},
      "post": {"operationId": "CreateOrder", "security": [{"pools": []}], "responses": {"201": {"description": "created"}}}
    }
  },
  "definitions": {"Order": {"type": "object", "properties": {"id": {"type": "string"}}}}
}`

func TestParseApiDefinition(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")

	doc, err := svc.parseApiDefinition([]byte(swaggerDefinition))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if doc.Name != "imported-api" || doc.Version != "1.0.0" || doc.Description != "an imported API" {
		t.Errorf("metadata mismatch: %+v", doc)
	}
	if doc.BasePath != "/team/v1" {
		t.Errorf("basePath = %q, want /team/v1", doc.BasePath)
	}
	if len(doc.Warnings) != 0 {
		t.Errorf("unexpected warnings: %v", doc.Warnings)
	}

	if _, ok := doc.Authorizers["my_auth"]; !ok {
		t.Fatalf("authorizer not parsed: %v", doc.Authorizers)
	}
	auth := doc.Authorizers["my_auth"]
	if auth.Type != "TOKEN" || auth.AuthorizerResultTtlInSeconds != 30 {
		t.Errorf("authorizer mismatch: %+v", auth)
	}
	pools := doc.Authorizers["pools"]
	if pools == nil || pools.Type != "COGNITO_USER_POOLS" {
		t.Fatalf("cognito authorizer not parsed: %+v", doc.Authorizers)
	}
	if len(pools.ProviderArns) != 1 || pools.ProviderArns[0] != "arn:aws:cognito-idp:us-east-1:123456789012:userpool/us-east-1_ABC" {
		t.Errorf("providerARNs not parsed as an array: %+v", pools.ProviderArns)
	}

	orderPost := findImportedOperation(t, doc, "/orders/{id}", "POST")
	if orderPost.Method.AuthorizationType != "COGNITO_USER_POOLS" {
		t.Errorf("cognito authorizer binding missing: %+v", orderPost.Method)
	}

	orderGet := findImportedOperation(t, doc, "/orders/{id}", "GET")
	if orderGet.Method.OperationName != "GetOrder" || !orderGet.Method.ApiKeyRequired {
		t.Errorf("operation mismatch: %+v", orderGet)
	}
	if orderGet.AuthorizerName != "my_auth" || orderGet.Method.AuthorizationType != "TOKEN" {
		t.Errorf("authorizer binding missing: %+v", orderGet)
	}
	if orderGet.Method.MethodIntegration == nil || orderGet.Method.MethodIntegration.Type != "AWS_PROXY" {
		t.Errorf("integration not parsed: %+v", orderGet.Method.MethodIntegration)
	}
	if len(orderGet.Method.MethodResponses) != 2 {
		t.Errorf("method responses = %v, want 2", orderGet.Method.MethodResponses)
	}

	model, ok := doc.Models["Order"]
	if !ok || !strings.Contains(model.Schema, `"id"`) {
		t.Errorf("definitions not parsed: %+v", doc.Models)
	}
}

// TestParseApiDefinitionDocumentSecurity pins that the document-level
// security declaration is ignored on import: an operation without its
// own security stays unprotected, while an operation-level declaration —
// including an explicit empty list — is what binds.
func TestParseApiDefinitionDocumentSecurity(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	doc, err := svc.parseApiDefinition([]byte(`{
  "swagger": "2.0",
  "info": {"title": "rootsec", "version": "1.0.0"},
  "securityDefinitions": {"api_key": {"type": "apiKey", "name": "x-api-key", "in": "header"}},
  "security": [{"api_key": []}],
  "paths": {
    "/inherit": {"get": {"operationId": "Inherit", "responses": {"200": {"description": "ok"}}}},
    "/optout": {"get": {"operationId": "OptOut", "security": [], "responses": {"200": {"description": "ok"}}}},
    "/own": {"get": {"operationId": "Own", "security": [{"api_key": []}], "responses": {"200": {"description": "ok"}}}}
  }
}`))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if inherit := findImportedOperation(t, doc, "/inherit", "GET"); inherit.Method.ApiKeyRequired {
		t.Errorf("document default security was applied: %+v", inherit.Method)
	}
	if optout := findImportedOperation(t, doc, "/optout", "GET"); optout.Method.ApiKeyRequired {
		t.Errorf("operation-level empty security bound a requirement: %+v", optout.Method)
	}
	if own := findImportedOperation(t, doc, "/own", "GET"); !own.Method.ApiKeyRequired {
		t.Errorf("operation-level security did not bind: %+v", own.Method)
	}
}

func TestParseApiDefinitionYAMLOpenAPI30(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	doc, err := svc.parseApiDefinition([]byte(`
openapi: 3.0.1
info:
  title: yaml-api
  version: 2.0.0
servers:
  - url: https://api.example.com/prefix
paths:
  /items:
    get:
      operationId: ListItems
      responses:
        "200":
          description: ok
components:
  schemas:
    Item:
      type: object
`))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if doc.Name != "yaml-api" {
		t.Errorf("title = %q, want yaml-api", doc.Name)
	}
	if doc.BasePath != "/prefix" {
		t.Errorf("server-derived basePath = %q, want /prefix", doc.BasePath)
	}
	if _, ok := doc.Models["Item"]; !ok {
		t.Errorf("components.schemas not parsed: %v", doc.Models)
	}
	if _, ok := doc.Paths["/items"]; !ok {
		t.Errorf("paths not parsed: %v", doc.Paths)
	}
}

func TestParseApiDefinitionRejections(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	for name, body := range map[string]string{
		"empty body":        ``,
		"not a document":    `just a string`,
		"missing family":    `{"info": {"title": "x"}}`,
		"missing title":     `{"swagger": "2.0"}`,
		"yaml missing info": "openapi: 3.0.1\npaths: {}\n",
		"authorizer on a non-apiKey scheme": `{"swagger": "2.0", "info": {"title": "x"},
			"securityDefinitions": {"bad": {"type": "oauth2",
			"x-amazon-apigateway-authorizer": {"type": "token"}}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := svc.parseApiDefinition([]byte(body))
			apiErr, ok := err.(*ApiGatewayError)
			if !ok {
				t.Fatalf("expected *ApiGatewayError, got %T: %v", err, err)
			}
			if apiErr.Code != "BadRequestException" {
				t.Errorf("code = %s, want BadRequestException", apiErr.Code)
			}
		})
	}
}

func findImportedOperation(t *testing.T, doc *apiDefinitionDocument, path, verb string) *apigatewaystore.ImportedMethod {
	t.Helper()
	for _, imported := range doc.Paths[path] {
		if imported.Verb == verb {
			return imported
		}
	}
	t.Fatalf("operation %s %s not found in %v", verb, path, doc.Paths)
	return nil
}

// TestApplyBasePathMode pins the documented basePath semantics: with
// basePath /a/b/c and paths /e and /f, prepend yields /a/b/c/e and
// /a/b/c/f, split yields /b/c/e and /b/c/f, and ignore keeps the declared
// paths.
func TestApplyBasePathMode(t *testing.T) {
	doc := &apiDefinitionDocument{BasePath: "/a/b/c", Paths: map[string][]*apigatewaystore.ImportedMethod{
		"/e": nil, "/f": nil,
	}}
	cases := []struct {
		mode string
		want []string
	}{
		{mode: "", want: []string{"/e", "/f"}},
		{mode: "ignore", want: []string{"/e", "/f"}},
		{mode: "prepend", want: []string{"/a/b/c/e", "/a/b/c/f"}},
		{mode: "split", want: []string{"/b/c/e", "/b/c/f"}},
	}
	for _, tc := range cases {
		t.Run("mode="+tc.mode, func(t *testing.T) {
			got := applyBasePathMode(doc, tc.mode)
			keys := make([]string, 0, len(got))
			for k := range got {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			if strings.Join(keys, ",") != strings.Join(tc.want, ",") {
				t.Errorf("paths = %v, want %v", keys, tc.want)
			}
		})
	}
}

// TestImportRestApiCoreRoundTrip pins the fresh import: the document's
// metadata names the API, declared paths become a resource tree with
// intermediate resources, operations carry their methods, definitions
// become models, and authorizers resolve onto the methods that reference
// them.
func TestImportRestApiCoreRoundTrip(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	stores := newImportStores(t)

	params, err := parseImportParameters(map[string]interface{}{})
	if err != nil {
		t.Fatalf("params: %v", err)
	}
	created, warnings, err := svc.importRestApiCore(stores, []byte(swaggerDefinition), params, false)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(warnings) != 0 {
		t.Errorf("unexpected warnings: %v", warnings)
	}
	if created.Name != "imported-api" || created.Version != "1.0.0" {
		t.Errorf("metadata mismatch: %+v", created)
	}

	paths := resourcePaths(created)
	for _, want := range []string{"/", "/orders", "/orders/{id}"} {
		if _, ok := paths[want]; !ok {
			t.Errorf("path %s missing from import: %v", want, paths)
		}
	}
	orderResource := paths["/orders/{id}"]
	if orderResource.ParentId != paths["/orders"].Id {
		t.Errorf("intermediate parent not linked: %s vs %s", orderResource.ParentId, paths["/orders"].Id)
	}
	get := orderResource.ResourceMethods["GET"]
	if get == nil || get.OperationName != "GetOrder" || !get.ApiKeyRequired || get.AuthorizationType != "TOKEN" {
		t.Errorf("GET /orders/{id} mismatch: %+v", get)
	}
	if get.AuthorizerId == "" {
		t.Error("authorizer reference not resolved to an id")
	}
	if get.MethodIntegration == nil || get.MethodIntegration.IntegrationHttpMethod != "POST" {
		t.Errorf("integration mismatch: %+v", get.MethodIntegration)
	}
	if _, ok := created.Models["Order"]; !ok {
		t.Errorf("model Order missing: %v", created.Models)
	}
}

// resourcePaths indexes a REST API's resources by path.
func resourcePaths(api *apigatewaystore.RestApi) map[string]*apigatewaystore.Resource {
	byPath := make(map[string]*apigatewaystore.Resource, len(api.Resources))
	for _, r := range api.Resources {
		byPath[r.Path] = r
	}
	return byPath
}

// TestPutRestApiCoreModes pins the merge/overwrite contract: merge
// installs the document on top of the existing API (unmentioned resources
// survive, mentioned paths gain exactly the imported verbs), overwrite
// replaces the resource tree, models and authorizers wholesale while the
// API identity survives both modes.
func TestPutRestApiCoreModes(t *testing.T) {
	seedDoc := `{"swagger":"2.0","info":{"title":"seed"},"paths":{
		"/keep":{"get":{"operationId":"KeepGet","responses":{"200":{"description":"ok"}}}}}}`
	otherDoc := `{"swagger":"2.0","info":{"title":"other"},"paths":{
		"/new":{"post":{"operationId":"NewPost","responses":{"200":{"description":"ok"}}}}},"definitions":{"M":{"type":"object"}}}`

	params, err := parseImportParameters(map[string]interface{}{})
	if err != nil {
		t.Fatalf("params: %v", err)
	}

	t.Run("merge keeps unmentioned resources", func(t *testing.T) {
		svc := NewAPIGatewayService("123456789012", "us-east-1")
		stores := newImportStores(t)
		seeded, _, err := svc.importRestApiCore(stores, []byte(seedDoc), params, false)
		if err != nil {
			t.Fatalf("seed import failed: %v", err)
		}

		merged, _, err := svc.putRestApiCore(stores, seeded.Id, "merge", []byte(otherDoc), params, false)
		if err != nil {
			t.Fatalf("merge failed: %v", err)
		}
		paths := resourcePaths(merged)
		if _, ok := paths["/keep"]; !ok {
			t.Errorf("merge dropped unmentioned /keep: %v", paths)
		}
		if _, ok := paths["/new"]; !ok {
			t.Errorf("merge did not install /new: %v", paths)
		}
		if _, ok := merged.Models["M"]; !ok {
			t.Errorf("merge did not install model M: %v", merged.Models)
		}
	})

	t.Run("overwrite replaces the definition", func(t *testing.T) {
		svc := NewAPIGatewayService("123456789012", "us-east-1")
		stores := newImportStores(t)
		seeded, _, err := svc.importRestApiCore(stores, []byte(seedDoc), params, false)
		if err != nil {
			t.Fatalf("seed import failed: %v", err)
		}

		replaced, _, err := svc.putRestApiCore(stores, seeded.Id, "overwrite", []byte(otherDoc), params, false)
		if err != nil {
			t.Fatalf("overwrite failed: %v", err)
		}
		if replaced.Id != seeded.Id {
			t.Errorf("overwrite changed the API id: %s vs %s", replaced.Id, seeded.Id)
		}
		paths := resourcePaths(replaced)
		if _, ok := paths["/keep"]; ok {
			t.Errorf("overwrite kept /keep: %v", paths)
		}
		if _, ok := paths["/new"]; !ok {
			t.Errorf("overwrite did not install /new: %v", paths)
		}
	})

	t.Run("unknown api is NotFound", func(t *testing.T) {
		svc := NewAPIGatewayService("123456789012", "us-east-1")
		stores := newImportStores(t)
		_, _, err := svc.putRestApiCore(stores, "apiX", "merge", []byte(seedDoc), params, false)
		apiErr, ok := err.(*ApiGatewayError)
		if !ok || apiErr.Code != "NotFoundException" {
			t.Fatalf("expected NotFoundException, got %v", err)
		}
	})

	t.Run("bogus mode is BadRequest", func(t *testing.T) {
		svc := NewAPIGatewayService("123456789012", "us-east-1")
		stores := newImportStores(t)
		seeded, _, err := svc.importRestApiCore(stores, []byte(seedDoc), params, false)
		if err != nil {
			t.Fatalf("seed import failed: %v", err)
		}
		_, _, err = svc.putRestApiCore(stores, seeded.Id, "bogus", []byte(seedDoc), params, false)
		apiErr, ok := err.(*ApiGatewayError)
		if !ok || apiErr.Code != "BadRequestException" {
			t.Fatalf("expected BadRequestException, got %v", err)
		}
	})
}

// TestImportRestApiFailOnWarnings pins the rollback contract: with
// failonwarnings set, an operation referencing an unknown security scheme
// produces a warning that rejects the import before anything is created.
func TestImportRestApiFailOnWarnings(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	stores := newImportStores(t)
	doc := `{"swagger":"2.0","info":{"title":"warny"},"paths":{
		"/a":{"get":{"operationId":"A","security":[{"ghost":[]}],"responses":{"200":{"description":"ok"}}}}}}`

	params, _ := parseImportParameters(map[string]interface{}{})
	_, _, err := svc.importRestApiCore(stores, []byte(doc), params, true)
	apiErr, ok := err.(*ApiGatewayError)
	if !ok || apiErr.Code != "BadRequestException" {
		t.Fatalf("expected BadRequestException with failonwarnings, got %v", err)
	}
	if !strings.Contains(apiErr.Message, "ghost") {
		t.Errorf("warning not carried into the error: %s", apiErr.Message)
	}

	// Without failonwarnings the import succeeds and the warning is
	// reported back with the created API.
	created, warnings, err := svc.importRestApiCore(stores, []byte(doc), params, false)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if len(warnings) != 1 || !strings.Contains(warnings[0], "ghost") {
		t.Errorf("warnings = %v, want the ghost scheme warning", warnings)
	}
	if len(created.Warnings) != 1 {
		t.Errorf("created API does not carry the warning: %+v", created.Warnings)
	}
}

func TestParseApiKeyCsv(t *testing.T) {
	rows, err := parseApiKeyCsv([]byte(
		"Name,key,description,Enabled,usageplanIds\n" +
			"First,apikey1234abcdefghij0123456789,An imported key,TRUE,\"p1, p2\"\n" +
			"Second,apikey5678abcdefghij0123456789,,false,\n"))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(rows))
	}
	if rows[0].name != "First" || rows[0].description != "An imported key" || !rows[0].enabled {
		t.Errorf("row 0 mismatch: %+v", rows[0])
	}
	if len(rows[0].usagePlanIds) != 2 || rows[0].usagePlanIds[0] != "p1" || rows[0].usagePlanIds[1] != "p2" {
		t.Errorf("row 0 usagePlanIds = %v, want [p1 p2]", rows[0].usagePlanIds)
	}
	if rows[1].enabled {
		t.Errorf("row 1 enabled = true, want false")
	}

	for name, body := range map[string]string{
		"missing name column": "key\napikey1234abcdefghij0123456789\n",
		"missing key column":  "name\nFirst\n",
		"no header":           "",
		"broken csv":          "name,key\n\"unterminated\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseApiKeyCsv([]byte(body)); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

// TestImportApiKeysCore pins the key import: rows import (upserting by
// key value), a row with an invalid value warns and is skipped, the
// format parameter rejects anything but csv, failonwarnings turns any
// warning into a rollback, and usage plan associations are created with
// duplicates reported as warnings.
func TestImportApiKeysCore(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")

	newTestStores := func(t *testing.T) *apiGatewayStores {
		stores := newImportStores(t)
		if _, err := stores.usage.CreateUsagePlan(&apigatewaystore.UsagePlan{Name: "plan-a"}); err != nil {
			t.Fatalf("seed usage plan: %v", err)
		}
		return stores
	}
	planId := func(t *testing.T, stores *apiGatewayStores) string {
		plans, err := stores.usage.ListUsagePlans(storecommon.ListOptions{})
		if err != nil || len(plans.Items) == 0 {
			t.Fatalf("list usage plans: %v", err)
		}
		return plans.Items[0].Id
	}

	t.Run("imports rows and associates plans", func(t *testing.T) {
		stores := newTestStores(t)
		plan := planId(t, stores)
		csv := "name,key,usageplanIds\n" +
			"First,apikey1234abcdefghij0123456789," + plan + "\n" +
			"Bad,short,\"" + plan + "\"\n"
		ids, warnings, err := svc.importApiKeysCore(stores, []byte(csv), "csv", false)
		if err != nil {
			t.Fatalf("import failed: %v", err)
		}
		if len(ids) != 1 {
			t.Errorf("ids = %v, want one", ids)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], "Bad") {
			t.Errorf("warnings = %v, want the short-value row", warnings)
		}

		keys, err := stores.usage.ListApiKeys(storecommon.ListOptions{})
		if err != nil {
			t.Fatalf("list keys: %v", err)
		}
		if len(keys.Items) != 1 || keys.Items[0].Name != "First" {
			t.Errorf("stored keys = %+v", keys.Items)
		}
	})

	t.Run("long name row warns and is skipped", func(t *testing.T) {
		stores := newTestStores(t)
		csv := "name,key\n" +
			strings.Repeat("n", apigatewaystore.ApiKeyNameMaxLength+1) + ",apikey1234abcdefghij0123456789\n" +
			"Good,apikeyabcdefghij0123456789xyz\n"
		ids, warnings, err := svc.importApiKeysCore(stores, []byte(csv), "csv", false)
		if err != nil {
			t.Fatalf("import failed: %v", err)
		}
		if len(ids) != 1 {
			t.Errorf("ids = %v, want one", ids)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], "Invalid key name") {
			t.Errorf("warnings = %v, want the long-name row", warnings)
		}
		keys, err := stores.usage.ListApiKeys(storecommon.ListOptions{})
		if err != nil {
			t.Fatalf("list keys: %v", err)
		}
		if len(keys.Items) != 1 || keys.Items[0].Name != "Good" {
			t.Errorf("stored keys = %+v", keys.Items)
		}
	})

	t.Run("reimport overwrites by value", func(t *testing.T) {
		stores := newTestStores(t)
		first := "name,key,description\nFirst,apikey1234abcdefghij0123456789,one\n"
		if _, _, err := svc.importApiKeysCore(stores, []byte(first), "csv", false); err != nil {
			t.Fatalf("first import: %v", err)
		}
		second := "name,key,description\nRenamed,apikey1234abcdefghij0123456789,two\n"
		if _, _, err := svc.importApiKeysCore(stores, []byte(second), "csv", false); err != nil {
			t.Fatalf("second import: %v", err)
		}
		keys, err := stores.usage.ListApiKeys(storecommon.ListOptions{})
		if err != nil {
			t.Fatalf("list keys: %v", err)
		}
		if len(keys.Items) != 1 || keys.Items[0].Name != "Renamed" || keys.Items[0].Description != "two" {
			t.Errorf("reimport did not overwrite: %+v", keys.Items)
		}
	})

	t.Run("stored duplicate association warns and rolls back on demand", func(t *testing.T) {
		stores := newTestStores(t)
		plan := planId(t, stores)
		csv := "name,key,usageplanIds\nFirst,apikey1234abcdefghij0123456789," + plan + "\n"
		if _, _, err := svc.importApiKeysCore(stores, []byte(csv), "csv", false); err != nil {
			t.Fatalf("first import: %v", err)
		}

		_, warnings, err := svc.importApiKeysCore(stores, []byte(csv), "csv", false)
		if err != nil {
			t.Fatalf("re-import: %v", err)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], plan) {
			t.Fatalf("warnings = %v, want the duplicate association with %s", warnings, plan)
		}

		_, _, err = svc.importApiKeysCore(stores, []byte(csv), "csv", true)
		apiErr, ok := err.(*ApiGatewayError)
		if !ok || apiErr.Code != "BadRequestException" {
			t.Fatalf("expected BadRequestException with failonwarnings, got %v", err)
		}
	})

	t.Run("stored duplicate with a distinct key id resolves by value", func(t *testing.T) {
		stores := newTestStores(t)
		plan := planId(t, stores)
		// A key created through the default path has a server-generated
		// identifier distinct from its value; usage-plan-key records live
		// under that identifier, so the import must resolve the row's
		// value to it — using the value itself would miss the duplicate.
		seeded, err := stores.usage.CreateApiKey(&apigatewaystore.ApiKey{
			Name:  "Seeded",
			Value: "apikey1234abcdefghij0123456789",
		})
		if err != nil {
			t.Fatalf("seed key: %v", err)
		}
		if seeded.Id == seeded.Value {
			t.Fatalf("seed produced Id==Value (%s); the fixture needs a distinct id", seeded.Id)
		}
		if _, err := stores.usage.CreateUsagePlanKey(plan, &apigatewaystore.UsagePlanKey{Id: seeded.Id}); err != nil {
			t.Fatalf("seed association: %v", err)
		}

		csv := "name,key,usageplanIds\nRenamed,apikey1234abcdefghij0123456789," + plan + "\n"
		ids, warnings, err := svc.importApiKeysCore(stores, []byte(csv), "csv", false)
		if err != nil {
			t.Fatalf("import: %v", err)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], plan) {
			t.Fatalf("warnings = %v, want the duplicate association held under the seeded key id", warnings)
		}
		if len(ids) != 1 || ids[0] != seeded.Id {
			t.Fatalf("ids = %v, want the seeded key id %s", ids, seeded.Id)
		}
		keys, err := stores.usage.ListApiKeys(storecommon.ListOptions{})
		if err != nil {
			t.Fatalf("list keys: %v", err)
		}
		if len(keys.Items) != 1 || keys.Items[0].Id != seeded.Id || keys.Items[0].Name != "Renamed" {
			t.Fatalf("stored keys = %+v", keys.Items)
		}
	})

	t.Run("in-file duplicate association warns", func(t *testing.T) {
		stores := newTestStores(t)
		plan := planId(t, stores)
		csv := "name,key,usageplanIds\n" +
			"First,apikey1234abcdefghij0123456789," + plan + "\n" +
			"Second,apikey1234abcdefghij0123456789," + plan + "\n"
		_, warnings, err := svc.importApiKeysCore(stores, []byte(csv), "csv", false)
		if err != nil {
			t.Fatalf("import: %v", err)
		}
		if len(warnings) != 1 || !strings.Contains(warnings[0], "apikey1234abcdefghij0123456789") {
			t.Fatalf("warnings = %v, want the in-file duplicate association", warnings)
		}
	})

	t.Run("failonwarnings rolls back", func(t *testing.T) {
		stores := newTestStores(t)
		csv := "name,key\nGood,apikey1234abcdefghij0123456789\nBad,short\n"
		_, _, err := svc.importApiKeysCore(stores, []byte(csv), "csv", true)
		apiErr, ok := err.(*ApiGatewayError)
		if !ok || apiErr.Code != "BadRequestException" {
			t.Fatalf("expected BadRequestException, got %v", err)
		}
		keys, err := stores.usage.ListApiKeys(storecommon.ListOptions{})
		if err != nil {
			t.Fatalf("list keys: %v", err)
		}
		if len(keys.Items) != 0 {
			t.Errorf("rollback left %d keys behind", len(keys.Items))
		}
	})

	t.Run("storage read failure fails the import", func(t *testing.T) {
		bucket := &failingForEachBucket{memBucket: newMemBucket()}
		stores := &apiGatewayStores{
			usage: apigatewaystore.NewUsageStore(failingStorage{bucket: bucket}, "123456789012", "us-east-1"),
		}
		if _, err := stores.usage.CreateApiKey(&apigatewaystore.ApiKey{
			Name:  "Seeded",
			Value: "apikey1234abcdefghij0123456789",
		}); err != nil {
			t.Fatalf("seed key: %v", err)
		}
		// The value→id lookup lists the key records; making that listing
		// fail must fail the import rather than silently drop the
		// duplicate-association warning.
		bucket.fail = true
		csv := "name,key\nRenamed,apikey1234abcdefghij0123456789\n"
		if _, _, err := svc.importApiKeysCore(stores, []byte(csv), "csv", false); err == nil {
			t.Fatal("expected the storage read failure to fail the import, got nil")
		}
	})

	t.Run("format validation", func(t *testing.T) {
		stores := newImportStores(t)
		if _, _, err := svc.importApiKeysCore(stores, []byte("x"), "", false); err == nil {
			t.Error("empty format accepted")
		}
		if _, _, err := svc.importApiKeysCore(stores, []byte("x"), "json", false); err == nil {
			t.Error("json format accepted")
		}
	})
}

// TestImportParametersValidation pins the parameters map: the documented
// values pass, anything else rejects.
func TestImportParametersValidation(t *testing.T) {
	valid := map[string]interface{}{"ignore": "documentation", "endpointConfigurationTypes": "REGIONAL", "basepath": "prepend"}
	if _, err := parseImportParameters(valid); err != nil {
		t.Fatalf("valid parameters rejected: %v", err)
	}
	for name, params := range map[string]map[string]interface{}{
		"bogus ignore":   {"ignore": "models"},
		"bogus endpoint": {"endpointConfigurationTypes": "FANCY"},
		"bogus basepath": {"basepath": "replace"},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseImportParameters(params); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}

// TestImportedDefinitionJSONShape documents the store shape of an imported
// definition as it serialises into the RestApi record.
func TestImportedDefinitionJSONShape(t *testing.T) {
	def := &apigatewaystore.ImportedDefinition{
		Resources: map[string][]*apigatewaystore.ImportedMethod{
			"/a": {{Verb: "GET", Method: &apigatewaystore.Method{OperationName: "A"}}},
		},
	}
	encoded, err := json.Marshal(def.Resources)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !strings.Contains(string(encoded), "A") {
		t.Errorf("operation lost in serialisation: %s", encoded)
	}
}
