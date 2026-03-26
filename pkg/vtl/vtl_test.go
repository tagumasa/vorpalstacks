package vtl

import (
	"encoding/base64"
	"net/url"
	"testing"
)

func TestInputPath(t *testing.T) {
	engine := NewEngine()
	engine.context.JSONBody = map[string]interface{}{
		"name": "test",
		"nested": map[string]interface{}{
			"value": 123,
		},
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$input.path('$.name')", "test"},
		{"$input.path('$.nested.value')", "123"},
		{"$input.path('$.missing')", ""},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestInputJson(t *testing.T) {
	engine := NewEngine()
	engine.context.JSONBody = map[string]interface{}{
		"name":  "test",
		"items": []interface{}{1, 2, 3},
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$input.json('$.name')", `"test"`},
		{"$input.json('$')", `{"items":[1,2,3],"name":"test"}`},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestInputBody(t *testing.T) {
	engine := NewEngine()
	engine.context.Body = `{"key":"value"}`

	result, err := engine.Transform("$input.body")
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != `{"key":"value"}` {
		t.Errorf("Transform = %q, want %q", result, `{"key":"value"}`)
	}
}

func TestInputParams(t *testing.T) {
	engine := NewEngine()
	engine.context.Params = map[string]string{
		"id":   "123",
		"name": "test",
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$input.params('id')", "123"},
		{"$input.params('name')", "test"},
		{"$input.params('missing')", ""},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestInputParamsAll(t *testing.T) {
	engine := NewEngine()
	engine.context.Params = map[string]string{
		"id": "123",
	}

	result, err := engine.Transform("$input.params()")
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != `{"id":"123"}` {
		t.Errorf("Transform = %q, want %q", result, `{"id":"123"}`)
	}
}

func TestUtilUrlEncode(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		template string
		expected string
	}{
		{"$util.urlEncode('hello world')", url.QueryEscape("hello world")},
		{"$util.urlEncode('test=value')", url.QueryEscape("test=value")},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestUtilUrlDecode(t *testing.T) {
	engine := NewEngine()

	encoded := url.QueryEscape("hello world")
	template := "$util.urlDecode('" + encoded + "')"

	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "hello world" {
		t.Errorf("Transform = %q, want %q", result, "hello world")
	}
}

func TestUtilBase64Encode(t *testing.T) {
	engine := NewEngine()

	result, err := engine.Transform("$util.base64Encode('hello')")
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	expected := base64.StdEncoding.EncodeToString([]byte("hello"))
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestUtilBase64Decode(t *testing.T) {
	engine := NewEngine()

	encoded := base64.StdEncoding.EncodeToString([]byte("hello"))
	template := "$util.base64Decode('" + encoded + "')"

	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "hello" {
		t.Errorf("Transform = %q, want %q", result, "hello")
	}
}

func TestUtilParseJson(t *testing.T) {
	engine := NewEngine()

	tests := []struct {
		template string
		expected string
	}{
		{`$util.parseJson('{"key":"value"}')`, `{"key":"value"}`},
		{`$util.parseJson('[1,2,3]')`, `[1,2,3]`},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestContextVariables(t *testing.T) {
	engine := NewEngine()
	engine.context.Context = map[string]interface{}{
		"requestId":  "req-123",
		"apiId":      "api-456",
		"stage":      "prod",
		"httpMethod": "GET",
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$context.requestId", "req-123"},
		{"$context.apiId", "api-456"},
		{"$context.stage", "prod"},
		{"$context.httpMethod", "GET"},
		{"$context.missing", ""},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestContextAuthorizer(t *testing.T) {
	engine := NewEngine()
	engine.context.Authorizer = map[string]interface{}{
		"principalId": "user-123",
		"claims":      "admin",
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$context.authorizer.principalId", "user-123"},
		{"$context.authorizer.claims", "admin"},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestContextIdentity(t *testing.T) {
	engine := NewEngine()
	engine.context.Identity = map[string]interface{}{
		"apiKey":   "key-123",
		"sourceIp": "192.168.1.1",
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$context.identity.apiKey", "key-123"},
		{"$context.identity.sourceIp", "192.168.1.1"},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestStageVariables(t *testing.T) {
	engine := NewEngine()
	engine.context.StageVars = map[string]string{
		"dbHost": "localhost",
		"dbPort": "5432",
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$stageVariables.dbHost", "localhost"},
		{"$stageVariables.dbPort", "5432"},
		{"$stageVariables['dbHost']", "localhost"},
		{"${stageVariables['dbPort']}", "5432"},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestTransformRequest(t *testing.T) {
	engine := NewEngine()
	engine.context.Context = map[string]interface{}{
		"stage": "prod",
	}

	templates := map[string]string{
		"application/json": `{"stage":"$context.stage","body":$input.json('$')}`,
	}

	body := []byte(`{"name":"test"}`)
	result, err := engine.TransformRequest(templates, "application/json", body)
	if err != nil {
		t.Fatalf("TransformRequest error: %v", err)
	}

	expected := `{"stage":"prod","body":{"name":"test"}}`
	if string(result) != expected {
		t.Errorf("TransformRequest = %q, want %q", string(result), expected)
	}
}

func TestTransformResponse(t *testing.T) {
	engine := NewEngine()
	engine.context.Context = map[string]interface{}{
		"requestId": "req-123",
	}

	templates := map[string]string{
		"application/json": `{"requestId":"$context.requestId","data":$input.json('$')}`,
	}

	body := []byte(`{"result":"success"}`)
	result, err := engine.TransformResponse(templates, "application/json", body)
	if err != nil {
		t.Fatalf("TransformResponse error: %v", err)
	}

	expected := `{"requestId":"req-123","data":{"result":"success"}}`
	if string(result) != expected {
		t.Errorf("TransformResponse = %q, want %q", string(result), expected)
	}
}

func TestComplexTemplate(t *testing.T) {
	engine := NewEngine()
	engine.context.JSONBody = map[string]interface{}{
		"userId": 123,
		"action": "update",
	}
	engine.context.Context = map[string]interface{}{
		"requestId":  "req-456",
		"apiId":      "api-789",
		"httpMethod": "POST",
	}
	engine.context.StageVars = map[string]string{
		"version": "v1",
	}

	template := `{"requestId":"$context.requestId","method":"$context.httpMethod","userId":$input.path('$.userId'),"version":"$stageVariables.version"}`

	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}

	expected := `{"requestId":"req-456","method":"POST","userId":123,"version":"v1"}`
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestArrayAccess(t *testing.T) {
	engine := NewEngine()
	engine.context.JSONBody = map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"name": "first", "value": 1},
			map[string]interface{}{"name": "second", "value": 2},
		},
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$input.path('$.items[0].name')", "first"},
		{"$input.path('$.items[1].value')", "2"},
		{"$input.json('$.items[0]')", `{"name":"first","value":1}`},
		{"$input.json('$.items[1].name')", `"second"`},
		{"$input.path('$.items[5]')", ""},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestWebSocketVariables(t *testing.T) {
	engine := NewEngine()
	engine.context.Context = map[string]interface{}{
		"connectionId": "conn-123",
		"connectedAt":  int64(1234567890),
		"eventType":    "MESSAGE",
		"routeKey":     "$connect",
		"messageId":    "msg-456",
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$context.connectionId", "conn-123"},
		{"$context.eventType", "MESSAGE"},
		{"$context.routeKey", "$connect"},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestRequestResponseOverride(t *testing.T) {
	engine := NewEngine()
	engine.context.RequestOverride = map[string]string{
		"header.Content-Type": "application/json",
	}
	engine.context.ResponseOverride = map[string]string{
		"header.X-Custom": "value",
	}

	tests := []struct {
		template string
		expected string
	}{
		{"$context.requestOverride.header.Content-Type", "application/json"},
		{"$context.responseOverride.header.X-Custom", "value"},
		{"$context.requestOverride.missing", ""},
	}

	for _, tt := range tests {
		result, err := engine.Transform(tt.template)
		if err != nil {
			t.Errorf("Transform(%q) error: %v", tt.template, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("Transform(%q) = %q, want %q", tt.template, result, tt.expected)
		}
	}
}

func TestEscapeJavaScriptSingleQuote(t *testing.T) {
	engine := NewEngine()

	result, err := engine.Transform(`$util.escapeJavaScript("it's")`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != `it\'s` {
		t.Errorf("Transform = %q, want %q", result, `it\'s`)
	}
}

func TestIfTrue(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["enabled"] = true

	result, err := engine.Transform(`#if($enabled)yes#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "yes" {
		t.Errorf("Transform = %q, want %q", result, "yes")
	}
}

func TestIfFalse(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["enabled"] = false

	result, err := engine.Transform(`#if($enabled)yes#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "" {
		t.Errorf("Transform = %q, want empty", result)
	}
}

func TestIfElse(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["enabled"] = false

	result, err := engine.Transform(`#if($enabled)yes#else no#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != " no" {
		t.Errorf("Transform = %q, want %q", result, " no")
	}
}

func TestIfElseifElse(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["value"] = 2

	result, err := engine.Transform(`#if($value == 1)one#elseif($value == 2)two#else other#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "two" {
		t.Errorf("Transform = %q, want %q", result, "two")
	}
}

func TestNestedIf(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["a"] = true
	engine.context.Context["b"] = true

	result, err := engine.Transform(`#if($a)#if($b)both#end#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "both" {
		t.Errorf("Transform = %q, want %q", result, "both")
	}
}

func TestForeach(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["items"] = []interface{}{"a", "b", "c"}

	result, err := engine.Transform(`#foreach($item in $items)$item#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "abc" {
		t.Errorf("Transform = %q, want %q", result, "abc")
	}
}

func TestForeachWithIndex(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["items"] = []interface{}{"a", "b"}

	result, err := engine.Transform(`#foreach($item in $items)$foreach.index#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "12" {
		t.Errorf("Transform = %q, want %q", result, "12")
	}
}

func TestForeachWithFirst(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["items"] = []interface{}{"a", "b", "c"}

	result, err := engine.Transform(`#foreach($item in $items)#if($foreach.first)F#end$item#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "Fabc" {
		t.Errorf("Transform = %q, want %q", result, "Fabc")
	}
}

func TestSetVariable(t *testing.T) {
	engine := NewEngine()

	result, err := engine.Transform(`#set($name = "test")$name`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "test" {
		t.Errorf("Transform = %q, want %q", result, "test")
	}
}

func TestSetFromVariable(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["original"] = "value"

	result, err := engine.Transform(`#set($copy = $original)$copy`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "value" {
		t.Errorf("Transform = %q, want %q", result, "value")
	}
}

func TestComparisonEquals(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["value"] = "test"

	result, err := engine.Transform(`#if($value == "test")match#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "match" {
		t.Errorf("Transform = %q, want %q", result, "match")
	}
}

func TestComparisonNotEquals(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["value"] = "other"

	result, err := engine.Transform(`#if($value != "test")no_match#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "no_match" {
		t.Errorf("Transform = %q, want %q", result, "no_match")
	}
}

func TestLogicalAnd(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["a"] = true
	engine.context.Context["b"] = true

	result, err := engine.Transform(`#if($a && $b)both#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "both" {
		t.Errorf("Transform = %q, want %q", result, "both")
	}
}

func TestLogicalOr(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["a"] = false
	engine.context.Context["b"] = true

	result, err := engine.Transform(`#if($a || $b)one#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "one" {
		t.Errorf("Transform = %q, want %q", result, "one")
	}
}

func TestLogicalNot(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["a"] = false

	result, err := engine.Transform(`#if(!$a)not_a#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "not_a" {
		t.Errorf("Transform = %q, want %q", result, "not_a")
	}
}

func TestNestedForeach(t *testing.T) {
	engine := NewEngine()
	engine.context.Context["outer"] = []interface{}{
		[]interface{}{1, 2},
		[]interface{}{3, 4},
	}

	result, err := engine.Transform(`#foreach($o in $outer)#foreach($i in $o)$i#end#end`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "1234" {
		t.Errorf("Transform = %q, want %q", result, "1234")
	}
}

func TestCommentBlock(t *testing.T) {
	engine := NewEngine()

	result, err := engine.Transform(`hello#* this is a comment *#world`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "helloworld" {
		t.Errorf("Transform = %q, want %q", result, "helloworld")
	}
}
