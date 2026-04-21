package vtl

import (
	"strings"
	"testing"
)

func TestCtxArgs(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Args: map[string]interface{}{
			"id":    "post-1",
			"title": "Hello World",
		},
	}

	tests := []struct {
		template string
		expected string
	}{
		{`$ctx.args.id`, "post-1"},
		{`$ctx.args.title`, "Hello World"},
		{`$ctx.args.missing`, "null"},
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

func TestCtxSource(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Source: map[string]interface{}{
			"id":     "post-1",
			"title":  "Test Post",
			"author": map[string]interface{}{"name": "Alice", "email": "alice@test.com"},
		},
	}

	tests := []struct {
		template string
		expected string
	}{
		{`$ctx.source.id`, "post-1"},
		{`$ctx.source.title`, "Test Post"},
		{`$ctx.source.author.name`, "Alice"},
		{`$ctx.source.missing`, "null"},
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

func TestCtxSourceNull(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Source: nil,
	}

	result, err := engine.Transform(`$ctx.source`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "null" {
		t.Errorf("Transform = %q, want %q", result, "null")
	}
}

func TestCtxInfo(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Info: &AppSyncFieldInfo{
			FieldName:           "getPost",
			ParentTypeName:      "Query",
			SelectionSetGraphQL: "{ id title author { name } }",
			SelectionSetList:    []string{"id", "title", "author"},
			RootTypeName:        "Query",
		},
	}

	tests := []struct {
		template string
		expected string
	}{
		{`$ctx.info.fieldName`, "getPost"},
		{`$ctx.info.parentTypeName`, "Query"},
		{`$ctx.info.selectionSetGraphQL`, "{ id title author { name } }"},
		{`$ctx.info.selectionSetList`, `["id","title","author"]`},
		{`$ctx.info.rootTypeName`, "Query"},
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

func TestCtxInfoVariables(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Info: &AppSyncFieldInfo{
			Variables: map[string]interface{}{"id": "123"},
		},
	}

	result, err := engine.Transform(`$ctx.info.variables`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != `{"id":"123"}` {
		t.Errorf("Transform = %q, want %q", result, `{"id":"123"}`)
	}
}

func TestCtxStash(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Stash: map[string]interface{}{
			"authorId": "a1",
		},
	}

	tests := []struct {
		template string
		expected string
	}{
		{`$ctx.stash.authorId`, "a1"},
		{`$ctx.stash.missing`, "null"},
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

func TestCtxIdentity(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Identity: map[string]interface{}{
			"sub":            "user-123",
			"username":       "alice",
			"cognito:groups": []string{"admin", "editor"},
		},
	}

	result, err := engine.Transform(`$ctx.identity.sub`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "user-123" {
		t.Errorf("Transform = %q, want %q", result, "user-123")
	}
}

func TestCtxResult(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Result: map[string]interface{}{"id": "post-1", "title": "Test"},
	}

	result, err := engine.Transform(`$ctx.result`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	expected := `{"id":"post-1","title":"Test"}`
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestCtxResultNull(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Result: nil,
	}

	result, err := engine.Transform(`$ctx.result`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "null" {
		t.Errorf("Transform = %q, want %q", result, "null")
	}
}

func TestCtxRequestHeaders(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Request: map[string]interface{}{
			"headers": map[string]interface{}{
				"x-api-key":     "da2-test",
				"Authorization": "Bearer token123",
			},
			"domainName": "api.example.com",
		},
	}

	result, err := engine.Transform(`$ctx.request.headers.x-api-key`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "da2-test" {
		t.Errorf("Transform = %q, want %q", result, "da2-test")
	}
}

func TestCtxPrev(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Prev: map[string]interface{}{
			"result": map[string]interface{}{"id": "post-1"},
		},
	}

	result, err := engine.Transform(`$ctx.prev.result.id`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "post-1" {
		t.Errorf("Transform = %q, want %q", result, "post-1")
	}
}

func TestCtxTopLevelRefs(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Args:   map[string]interface{}{"id": "1"},
		Source: map[string]interface{}{"title": "Test"},
		Stash:  map[string]interface{}{"key": "val"},
	}

	tests := []struct {
		template string
		expected string
	}{
		{`$ctx.args`, `{"id":"1"}`},
		{`$ctx.source`, `{"title":"Test"}`},
		{`$ctx.stash`, `{"key":"val"}`},
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

func TestCtxNilAppSyncCtx(t *testing.T) {
	engine := NewEngine()

	result, err := engine.Transform(`$ctx.args.id`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != `$ctx.args.id` {
		t.Errorf("Transform should not modify $ctx when AppSyncCtx is nil, got %q", result)
	}
}

func TestAppSyncRequestTemplate(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Args: map[string]interface{}{"id": "post-1"},
		Info: &AppSyncFieldInfo{
			FieldName:      "getPost",
			ParentTypeName: "Query",
		},
	}

	template := `{"version":"2018-05-29","operation":"GetItem","key":{"id":{"S":"$ctx.args.id"}}}`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	expected := `{"version":"2018-05-29","operation":"GetItem","key":{"id":{"S":"post-1"}}}`
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestAppSyncResponseTemplate(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Result: map[string]interface{}{
			"id":    "post-1",
			"title": "Test Post",
		},
	}

	template := `{"id":$ctx.result.id,"title":$ctx.result.title}`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	expected := `{"id":post-1,"title":Test Post}`
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestUtilIsNull(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	tests := []struct {
		template string
		expected string
	}{
		{`$util.isNull(null)`, "true"},
		{`$util.isNull("hello")`, "false"},
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

func TestUtilIsNullOrEmpty(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	tests := []struct {
		template string
		expected string
	}{
		{`$util.isNullOrEmpty(null)`, "true"},
		{`$util.isNullOrEmpty("")`, "true"},
		{`$util.isNullOrEmpty('')`, "true"},
		{`$util.isNullOrEmpty("hello")`, "false"},
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

func TestUtilDefaultIfNull(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	tests := []struct {
		template string
		expected string
	}{
		{`$util.defaultIfNull(null, "default")`, "default"},
		{`$util.defaultIfNull("value", "default")`, "value"},
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

func TestUtilAutoId(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	result, err := engine.Transform(`$util.autoId()`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(result) != 36 {
		t.Errorf("autoId should return UUID (36 chars), got %q (%d chars)", result, len(result))
	}
	if result[8] != '-' || result[13] != '-' || result[18] != '-' || result[23] != '-' {
		t.Errorf("autoId should return UUID format, got %q", result)
	}
}

func TestUtilAutoIdUnique(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.autoId()_$util.autoId()`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	parts := strings.SplitN(result, "_", 2)
	if len(parts) != 2 {
		t.Fatalf("Expected two UUIDs separated by _, got %q", result)
	}
	if parts[0] == parts[1] {
		t.Errorf("Two autoId calls should return different UUIDs, got same: %q", parts[0])
	}
}

func TestUtilAutoIdInSet(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `#set($postId = $util.autoId())$postId`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(result) != 36 {
		t.Errorf("autoId inside #set should return UUID (36 chars), got %q (%d chars)", result, len(result))
	}
	if result[8] != '-' || result[13] != '-' || result[18] != '-' || result[23] != '-' {
		t.Errorf("autoId inside #set should return UUID format, got %q", result)
	}
}

func TestUtilAutoIdInSetUsedInPayload(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Args: map[string]interface{}{
			"title": "Test Post",
		},
	}

	template := `#set($postId = $util.autoId()){"key":{"id":{"S":"$postId"}},"attributeValues":{"id":{"S":"$postId"},"title":{"S":"$ctx.args.title"}}}`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if !strings.Contains(result, `"S":"`) {
		t.Errorf("Expected DynamoDB JSON format, got %q", result)
	}
	if strings.Contains(result, "$postId") {
		t.Errorf("Variable $postId should be resolved, got %q", result)
	}
	if !strings.Contains(result, "Test Post") {
		t.Errorf("Expected $ctx.args.title resolved, got %q", result)
	}
}

func TestUtilTimeNowInSet(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `#set($now = $util.time.nowISO8601())$now`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(result) != 24 {
		t.Errorf("time.nowISO8601 inside #set should return 24 chars, got %q (%d)", result, len(result))
	}
}

func TestUtilEpochSecondsInSet(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `#set($ts = $util.time.nowEpochSeconds())$ts`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(result) < 9 {
		t.Errorf("time.nowEpochSeconds inside #set should return ~10 digits, got %q", result)
	}
}

func TestUtilError(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.error("Not found", "NotFound")`
	_, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(engine.AppSyncCtx.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(engine.AppSyncCtx.Errors))
	}
	appErr := engine.AppSyncCtx.Errors[0]
	if appErr.Message != "Not found" {
		t.Errorf("Error message = %q, want %q", appErr.Message, "Not found")
	}
	if appErr.ErrorType != "NotFound" {
		t.Errorf("Error type = %q, want %q", appErr.ErrorType, "NotFound")
	}
}

func TestUtilAppendError(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.appendError("Warning", "Warning")rest of template`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(engine.AppSyncCtx.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(engine.AppSyncCtx.Errors))
	}
	appErr := engine.AppSyncCtx.Errors[0]
	if appErr.Message != "Warning" {
		t.Errorf("Error message = %q, want %q", appErr.Message, "Warning")
	}
	if result != "rest of template" {
		t.Errorf("Transform = %q, want %q", result, "rest of template")
	}
}

func TestUtilErrorWithData(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.appendError("Invalid input", "BadRequest", {"field":"title"})`
	_, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(engine.AppSyncCtx.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(engine.AppSyncCtx.Errors))
	}
	appErr := engine.AppSyncCtx.Errors[0]
	if appErr.Data == nil {
		t.Fatal("Expected error data, got nil")
	}
	data, ok := appErr.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map data, got %T", appErr.Data)
	}
	if data["field"] != "title" {
		t.Errorf("Error data field = %v, want %q", data["field"], "title")
	}
}

func TestUtilUnauthorized(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.unauthorized()`
	_, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(engine.AppSyncCtx.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(engine.AppSyncCtx.Errors))
	}
	appErr := engine.AppSyncCtx.Errors[0]
	if appErr.Message != "Unauthorized" {
		t.Errorf("Error message = %q, want %q", appErr.Message, "Unauthorized")
	}
}

func TestUtilTimeNowISO8601(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	result, err := engine.Transform(`$util.time.nowISO8601()`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(result) != 24 {
		t.Errorf("ISO8601 timestamp should be 24 chars (e.g. 2006-01-02T15:04:05.000Z), got %q (%d)", result, len(result))
	}
	if !strings.HasSuffix(result, "Z") {
		t.Errorf("ISO8601 timestamp should end with Z, got %q", result)
	}
}

func TestUtilTimeNowEpochSeconds(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	result, err := engine.Transform(`$util.time.nowEpochSeconds()`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(result) < 9 {
		t.Errorf("Epoch seconds should be ~10 digits, got %q", result)
	}
}

func TestUtilTimeNowEpochMilliSeconds(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	result, err := engine.Transform(`$util.time.nowEpochMilliSeconds()`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(result) < 12 {
		t.Errorf("Epoch millis should be ~13 digits, got %q", result)
	}
}

func TestUtilQr(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.qr($myVar = "test")$myVar`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "test" {
		t.Errorf("Transform = %q, want %q", result, "test")
	}
}

func TestUtilQuietIf(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.quietIf(true)rest`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "rest" {
		t.Errorf("Transform = %q, want %q", result, "rest")
	}
}

func TestUtilValidatePass(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.validate(true, "Must be true")`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "" {
		t.Errorf("Transform = %q, want empty", result)
	}
	if len(engine.AppSyncCtx.Errors) != 0 {
		t.Errorf("Expected no errors on valid condition, got %d", len(engine.AppSyncCtx.Errors))
	}
}

func TestUtilValidateFail(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.validate(false, "Title is required", "Validation")`
	_, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if len(engine.AppSyncCtx.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(engine.AppSyncCtx.Errors))
	}
	appErr := engine.AppSyncCtx.Errors[0]
	if appErr.Message != "Title is required" {
		t.Errorf("Error message = %q, want %q", appErr.Message, "Title is required")
	}
	if appErr.ErrorType != "Validation" {
		t.Errorf("Error type = %q, want %q", appErr.ErrorType, "Validation")
	}
}

func TestUtilIsString(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	tests := []struct {
		template string
		expected string
	}{
		{`$util.isString("hello")`, "true"},
		{`$util.isString(null)`, "false"},
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

func TestUtilIsList(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	tests := []struct {
		template string
		expected string
	}{
		{`$util.isList([1,2,3])`, "true"},
		{`$util.isList("hello")`, "false"},
		{`$util.isList(null)`, "false"},
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

func TestUtilDynamoDBToMapValues(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	result, err := engine.Transform(`$util.dynamodb.toMapValues($ctx.result)`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "null" {
		t.Errorf("Transform = %q, want %q", result, "null")
	}
}

func TestAppSyncControlFlowIntegration(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Args: map[string]interface{}{
			"id": "post-1",
		},
		Source: map[string]interface{}{
			"title": "Existing Post",
		},
	}

	template := `#if($ctx.args.id != null){"id":"$ctx.args.id"}#else{"id":"unknown"}#end`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	expected := `{"id":"post-1"}`
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestAppSyncForeachWithCtx(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Args: map[string]interface{}{
			"ids": []interface{}{"a", "b", "c"},
		},
	}

	engine.context.Context["ids"] = engine.AppSyncCtx.Args["ids"]

	template := `#foreach($id in $ids)"$id"#if($foreach.hasNext),#end#end`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	expected := `"a","b","c"`
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestAppSyncFullResolverPipeline(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{
		Args: map[string]interface{}{
			"id": "post-1",
		},
		Result: map[string]interface{}{
			"id":     "post-1",
			"title":  "Hello World",
			"author": map[string]interface{}{"id": "a1", "name": "Alice"},
		},
		Info: &AppSyncFieldInfo{
			FieldName:      "getPost",
			ParentTypeName: "Query",
		},
	}

	template := `{"id":"$ctx.result.id","title":"$ctx.result.title"}`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	expected := `{"id":"post-1","title":"Hello World"}`
	if result != expected {
		t.Errorf("Transform = %q, want %q", result, expected)
	}
}

func TestAppSyncLogFunctions(t *testing.T) {
	engine := NewEngine()
	engine.AppSyncCtx = &AppSyncContext{}

	template := `$util.log.info("processing")$util.log.error("something failed")$util.log.warn("careful")done`
	result, err := engine.Transform(template)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if result != "done" {
		t.Errorf("Transform = %q, want %q", result, "done")
	}
}

func TestAppSyncCtxIgnoredWithoutAppSyncCtx(t *testing.T) {
	engine := NewEngine()
	engine.context.Context = map[string]interface{}{
		"requestId": "req-123",
	}

	result, err := engine.Transform(`$context.requestId and $ctx.args.id`)
	if err != nil {
		t.Fatalf("Transform error: %v", err)
	}
	if !strings.Contains(result, "req-123") {
		t.Errorf("Context variables should still work, got %q", result)
	}
	if !strings.Contains(result, "$ctx.args.id") {
		t.Errorf("$ctx should be preserved without AppSyncCtx, got %q", result)
	}
}
