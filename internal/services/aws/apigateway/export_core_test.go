package apigateway

import (
	"strings"
	"testing"

	apigatewaystore "vorpalstacks/internal/store/aws/apigateway"
)

// exportFixture serves a REST API record whose single deployment snapshots
// one resource with one GET method, so the export document has a path, an
// operation and a model definition to render.
const exportFixture = `{"id":"api1","name":"probe","version":"1.0",
 "resources":{"res-root":{"id":"res-root","path":"/","path_part":"/",
   "resource_methods":{"GET":{"http_method":"GET","authorization_type":"NONE","operation_name":"ProbeGet"}}}},
 "models":{"ProbeModel":{"name":"ProbeModel","schema":"{\"type\":\"object\"}"}},
 "deployments":{"dep1":{"id":"dep1","created_date":"2026-01-02T03:04:05Z",
   "snapshot":{"resources":{"res-root":{"id":"res-root","path":"/","path_part":"/",
     "resource_methods":{"GET":{"http_method":"GET","authorization_type":"NONE","operation_name":"ProbeGet"}}}},
     "models":{"ProbeModel":{"name":"ProbeModel","schema":"{\"type\":\"object\"}"}}}}},
 "stages":{"prod":{"stage_name":"prod","rest_api_id":"api1","deployment_id":"dep1"}}}`

// TestGetExportCoreYAMLAccept pins the export document families and the
// Accept-driven content negotiation. The YAML branch and the trimming of
// media-type parameters are unreachable through the Go SDK client (its
// build-stack middleware overwrites Accept with application/json on every
// API Gateway operation), so the server contract is pinned here.
func TestGetExportCoreYAMLAccept(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	fs := fakeStorage{bucket: &fakeBucket{get: func(key []byte) ([]byte, error) {
		return []byte(exportFixture), nil
	}}}
	stores := &apiGatewayStores{
		restApis: apigatewaystore.NewRestApiStore(fs, "123456789012", "us-east-1"),
	}

	tests := []struct {
		name          string
		exportType    string
		accepts       string
		wantContent   string
		wantFamilyKey string
		wantFamilyVer string
	}{
		{name: "swagger json", exportType: "swagger", accepts: "", wantContent: "application/json", wantFamilyKey: "swagger", wantFamilyVer: "2.0"},
		{name: "oas30 json explicit", exportType: "oas30", accepts: "application/json", wantContent: "application/json", wantFamilyKey: "openapi", wantFamilyVer: "3.0.1"},
		{name: "oas30 yaml", exportType: "oas30", accepts: "application/yaml", wantContent: "application/yaml", wantFamilyKey: "openapi", wantFamilyVer: "3.0.1"},
		{name: "yaml with media-type parameter", exportType: "oas30", accepts: "application/yaml; charset=utf-8", wantContent: "application/yaml", wantFamilyKey: "openapi", wantFamilyVer: "3.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, contentType, err := svc.getExportCore(stores, "api1", "prod", tt.exportType, tt.accepts)
			if err != nil {
				t.Fatalf("getExportCore failed: %v", err)
			}
			if contentType != tt.wantContent {
				t.Errorf("contentType = %q, want %q", contentType, tt.wantContent)
			}
			text := string(body)
			if tt.wantContent == "application/json" {
				if !strings.HasPrefix(text, "{") {
					t.Errorf("json export does not start with '{': %.40s", text)
				}
				if !strings.Contains(text, `"`+tt.wantFamilyKey+`": "`+tt.wantFamilyVer+`"`) {
					t.Errorf("json export missing %q: %q marker", tt.wantFamilyKey, tt.wantFamilyVer)
				}
			} else {
				if strings.HasPrefix(text, "{") {
					t.Errorf("yaml export looks like json: %.40s", text)
				}
				if !strings.Contains(text, tt.wantFamilyKey+": "+tt.wantFamilyVer) {
					t.Errorf("yaml export missing %q: %q marker", tt.wantFamilyKey, tt.wantFamilyVer)
				}
			}
			if !strings.Contains(text, "ProbeGet") {
				t.Errorf("export omits the snapshot's method: %.80s", text)
			}
			if !strings.Contains(text, "ProbeModel") {
				t.Errorf("export omits the snapshot's model definition: %.80s", text)
			}
		})
	}
}

// TestGetModelTemplateCoreJSON pins the generated sample mapping template:
// one quoted entry per object property, comma-separated, with no trailing
// comma — the rendered skeleton must stay valid JSON.
func TestGetModelTemplateCoreJSON(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	fixture := `{"id":"api1","name":"probe","version":"1.0",
 "models":{"ProbeModel":{"name":"ProbeModel",
   "schema":"{\"type\":\"object\",\"properties\":{\"id\":{\"type\":\"string\"},\"name\":{\"type\":\"string\"}}}"}}}`
	fs := fakeStorage{bucket: &fakeBucket{get: func(key []byte) ([]byte, error) {
		return []byte(fixture), nil
	}}}
	stores := &apiGatewayStores{
		restApis: apigatewaystore.NewRestApiStore(fs, "123456789012", "us-east-1"),
	}

	tmpl, err := svc.getModelTemplateCore(stores, "api1", "ProbeModel")
	if err != nil {
		t.Fatalf("getModelTemplateCore: %v", err)
	}
	want := "#set($inputRoot = $input.path('$'))\n{\n  \"id\": \"$inputRoot.id\",\n  \"name\": \"$inputRoot.name\"\n}"
	if tmpl != want {
		t.Errorf("template = %q, want %q", tmpl, want)
	}
}

// TestGetExportCoreRejections pins the core's input validation: an unknown
// export family and an unsupported Accept media type are BadRequestException,
// never a 500 or a silent fallback.
func TestGetExportCoreRejections(t *testing.T) {
	svc := NewAPIGatewayService("123456789012", "us-east-1")
	fs := fakeStorage{bucket: &fakeBucket{get: func(key []byte) ([]byte, error) {
		return []byte(exportFixture), nil
	}}}
	stores := &apiGatewayStores{
		restApis: apigatewaystore.NewRestApiStore(fs, "123456789012", "us-east-1"),
	}

	for name, tc := range map[string]struct {
		exportType, accepts string
	}{
		"bogus export type": {exportType: "bogus", accepts: ""},
		"bogus accepts":     {exportType: "swagger", accepts: "text/html"},
	} {
		t.Run(name, func(t *testing.T) {
			_, _, err := svc.getExportCore(stores, "api1", "prod", tc.exportType, tc.accepts)
			apiErr, ok := err.(*ApiGatewayError)
			if !ok {
				t.Fatalf("expected *ApiGatewayError, got %T: %v", err, err)
			}
			if apiErr.Code != "BadRequestException" || apiErr.HTTPStatus != 400 {
				t.Errorf("code/status = %s/%d, want BadRequestException/400", apiErr.Code, apiErr.HTTPStatus)
			}
		})
	}
}
