package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
)

type vtlEvalTest struct {
	name     string
	context  string
	template string
	validate func(out *appsync.EvaluateMappingTemplateOutput) error
}

func vtlResult(out *appsync.EvaluateMappingTemplateOutput) string {
	if out.EvaluationResult != nil {
		return *out.EvaluationResult
	}
	return ""
}

func vtlStash(out *appsync.EvaluateMappingTemplateOutput) string {
	if out.Stash != nil {
		return *out.Stash
	}
	return ""
}

func (r *TestRunner) runAppSyncVTLTests(res *appsyncResources) []TestResult {
	var results []TestResult
	ctx := res.ctx
	client := res.client

	tests := []vtlEvalTest{
		{
			name:     "VTL_CtxArgs_FieldAccess",
			context:  `{"args":{"id":"post-123","title":"Hello"}}`,
			template: `$ctx.args.id`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "post-123") {
					return fmt.Errorf("expected post-123, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxArgs_TopLevelJSON",
			context:  `{"args":{"id":"post-123","title":"Hello"}}`,
			template: `$ctx.args`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				if m["id"] != "post-123" {
					return fmt.Errorf("expected id=post-123, got %v", m["id"])
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxArgs_MultipleFields",
			context:  `{"args":{"id":"abc","title":"My Post","status":"draft"}}`,
			template: `{"id":"$ctx.args.id","title":"$ctx.args.title"}`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				r := vtlResult(out)
				if !strings.Contains(r, `"abc"`) || !strings.Contains(r, `"My Post"`) {
					return fmt.Errorf("expected both fields resolved, got %s", r)
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxSource_NestedAccess",
			context:  `{"source":{"author":{"name":"Alice","id":"u1"}}}`,
			template: `$ctx.source.author.name`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "Alice") {
					return fmt.Errorf("expected Alice, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_CtxStash_ReadWrite",
			context: `{"stash":{"existingKey":"existingVal"}}`,
			template: `#set($ctx.stash.newKey = "newVal")
$ctx.stash.existingKey`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "existingVal") {
					return fmt.Errorf("expected existingVal in result, got %s", vtlResult(out))
				}
				var s map[string]interface{}
				if e := json.Unmarshal([]byte(vtlStash(out)), &s); e != nil {
					return fmt.Errorf("stash is not valid JSON: %s", vtlStash(out))
				}
				if s["newKey"] != "newVal" {
					return fmt.Errorf("expected stash.newKey=newVal, got %v", s["newKey"])
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxIdentity",
			context:  `{"identity":{"sub":"user-789","username":"bob","claims":{"email":"bob@test.com"}}}`,
			template: `$ctx.identity.username`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "bob") {
					return fmt.Errorf("expected bob, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxInfo_FieldName",
			context:  `{"info":{"fieldName":"getPost","parentTypeName":"Query","selectionSetGraphQL":"{ id title }","selectionSetList":["id","title"]}}`,
			template: `$ctx.info.fieldName`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "getPost") {
					return fmt.Errorf("expected getPost, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxInfo_ParentTypeName",
			context:  `{"info":{"fieldName":"getPost","parentTypeName":"Query"}}`,
			template: `$ctx.info.parentTypeName`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "Query") {
					return fmt.Errorf("expected Query, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxResult_FieldAccess",
			context:  `{"result":{"id":"post-1","title":"Existing Post","views":42}}`,
			template: `$ctx.result.title`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "Existing Post") {
					return fmt.Errorf("expected 'Existing Post', got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxResult_TopLevelJSON",
			context:  `{"result":{"id":"post-1","title":"Test"}}`,
			template: `$ctx.result`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				if m["id"] != "post-1" {
					return fmt.Errorf("expected id=post-1, got %v", m["id"])
				}
				return nil
			},
		},
		{
			name:     "VTL_CtxRequest_Headers",
			context:  `{"request":{"headers":{"x-api-key":"da2-testkey","authorization":"Bearer tok"}}}`,
			template: `$ctx.request.headers.x-api-key`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "da2-testkey") {
					return fmt.Errorf("expected da2-testkey, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilAutoId",
			context:  `{}`,
			template: `$util.autoId()`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				trimmed := strings.TrimSpace(vtlResult(out))
				if len(trimmed) != 36 {
					return fmt.Errorf("expected UUID (36 chars), got %d chars: %s", len(trimmed), trimmed)
				}
				if strings.Count(trimmed, "-") != 4 {
					return fmt.Errorf("expected UUID with 4 dashes, got %s", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilToJson_Object",
			context:  `{"args":{"id":"post-1","title":"Test"}}`,
			template: `$util.toJson($ctx.args)`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				if m["id"] != "post-1" {
					return fmt.Errorf("expected id=post-1, got %v", m["id"])
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilToJson_MapInSet",
			context: `{"args":{"name":"test"}}`,
			template: `#set($payload = {"version":"2018-05-29","operation":"GetItem","key":{"id":{"S":"$ctx.args.name"}}})
$util.toJson($payload)`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				if m["version"] != "2018-05-29" {
					return fmt.Errorf("expected version=2018-05-29, got %v", m["version"])
				}
				if m["operation"] != "GetItem" {
					return fmt.Errorf("expected operation=GetItem, got %v", m["operation"])
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilQr_QuietReference",
			context: `{}`,
			template: `#set($x = $util.qr("ignored"))
visible`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				trimmed := strings.TrimSpace(vtlResult(out))
				if trimmed != "visible" {
					return fmt.Errorf("expected 'visible', got '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilError",
			context:  `{}`,
			template: `$util.error("something went wrong", "CustomError")`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if out.Error == nil || out.Error.Message == nil {
					return fmt.Errorf("expected error in response, got nil")
				}
				if !strings.Contains(*out.Error.Message, "something went wrong") {
					return fmt.Errorf("expected error message 'something went wrong', got %s", *out.Error.Message)
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilUnauthorized",
			context:  `{}`,
			template: `$util.unauthorized()`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if out.Error == nil || out.Error.Message == nil {
					return fmt.Errorf("expected unauthorized error, got nil")
				}
				if !strings.Contains(*out.Error.Message, "Unauthorized") {
					return fmt.Errorf("expected Unauthorized error, got %s", *out.Error.Message)
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilIsNull",
			context: `{"args":{"val":null}}`,
			template: `#if($util.isNull($ctx.args.val))
null-checked
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "null-checked") {
					return fmt.Errorf("expected null-checked branch, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilIsNullOrEmpty_EmptyString",
			context: `{"args":{"val":""}}`,
			template: `#if($util.isNullOrEmpty($ctx.args.val))
empty
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "empty") {
					return fmt.Errorf("expected empty branch, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilDefaultIfNull",
			context:  `{"args":{"val":null}}`,
			template: `$util.defaultIfNull($ctx.args.val, "fallback")`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "fallback") {
					return fmt.Errorf("expected fallback, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilIsString",
			context: `{"args":{"val":"hello"}}`,
			template: `#if($util.isString($ctx.args.val))
yes-string
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "yes-string") {
					return fmt.Errorf("expected yes-string branch, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilIsList",
			context: `{"args":{"val":["a","b"]}}`,
			template: `#if($util.isList($ctx.args.val))
yes-list
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "yes-list") {
					return fmt.Errorf("expected yes-list branch, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilValidate_Pass",
			context: `{"args":{"email":"test@example.com"}}`,
			template: `$util.validate($ctx.args.email != "", "email is required")
valid`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "valid") {
					return fmt.Errorf("expected valid output, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilValidate_Fail",
			context:  `{"args":{"email":""}}`,
			template: `$util.validate($ctx.args.email != "", "email is required")`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if out.Error == nil || out.Error.Message == nil {
					return fmt.Errorf("expected validation error, got nil")
				}
				if !strings.Contains(*out.Error.Message, "email is required") {
					return fmt.Errorf("expected 'email is required', got %s", *out.Error.Message)
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilDynamoDBToMapValues",
			context:  `{"args":{"id":"post-1","count":5}}`,
			template: `$util.toJson($util.dynamodb.toMapValues($ctx.args))`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				idVal, ok := m["id"]
				if !ok {
					return fmt.Errorf("expected id key in result, got %v", m)
				}
				idMap, ok := idVal.(map[string]interface{})
				if !ok {
					return fmt.Errorf("expected id to be a DynamoDB attribute, got %T", idVal)
				}
				if s, ok := idMap["S"].(string); !ok || s != "post-1" {
					return fmt.Errorf("expected id.S=post-1, got %v", idMap)
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilTimeNowISO8601",
			context:  `{}`,
			template: `$util.time.nowISO8601()`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				trimmed := strings.TrimSpace(vtlResult(out))
				if !strings.Contains(trimmed, "T") || !strings.Contains(trimmed, "Z") {
					return fmt.Errorf("expected ISO 8601 format (contains T and Z), got %s", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_UtilTimeNowEpochSeconds",
			context:  `{}`,
			template: `$util.time.nowEpochSeconds()`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				trimmed := strings.TrimSpace(vtlResult(out))
				if len(trimmed) < 10 {
					return fmt.Errorf("expected epoch seconds (>=10 digits), got %s", trimmed)
				}
				return nil
			},
		},
		{
			name:    "VTL_SetVariable",
			context: `{}`,
			template: `#set($greeting = "Hello World")
$greeting`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "Hello World") {
					return fmt.Errorf("expected Hello World, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_SetFromCtxArgs",
			context: `{"args":{"name":"Alice"}}`,
			template: `#set($name = $ctx.args.name)
Hello $name`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "Hello Alice") {
					return fmt.Errorf("expected 'Hello Alice', got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_If_TrueBranch",
			context: `{"args":{"active":true}}`,
			template: `#if($ctx.args.active)
active-yes
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "active-yes") {
					return fmt.Errorf("expected active-yes, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_IfElse_TrueBranch",
			context: `{"args":{"status":"active"}}`,
			template: `#if($ctx.args.status == "active")
yes-active
#else
not-active
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				r := vtlResult(out)
				if !strings.Contains(r, "yes-active") {
					return fmt.Errorf("expected yes-active, got %s", r)
				}
				if strings.Contains(r, "not-active") {
					return fmt.Errorf("else branch should not appear")
				}
				return nil
			},
		},
		{
			name:    "VTL_IfElse_FalseBranch",
			context: `{"args":{"status":"inactive"}}`,
			template: `#if($ctx.args.status == "active")
yes-active
#else
not-active
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				r := vtlResult(out)
				if !strings.Contains(r, "not-active") {
					return fmt.Errorf("expected not-active, got %s", r)
				}
				if strings.Contains(r, "yes-active") {
					return fmt.Errorf("true branch should not appear")
				}
				return nil
			},
		},
		{
			name:    "VTL_IfElseIfElse_MiddleBranch",
			context: `{"args":{"level":2}}`,
			template: `#if($ctx.args.level == 1)
one
#elseif($ctx.args.level == 2)
two
#else
other
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				r := vtlResult(out)
				if !strings.Contains(r, "two") {
					return fmt.Errorf("expected two, got %s", r)
				}
				if strings.Contains(r, "one") || strings.Contains(r, "other") {
					return fmt.Errorf("wrong branch appeared")
				}
				return nil
			},
		},
		{
			name:    "VTL_Foreach_Basic",
			context: `{"args":{"items":["a","b","c"]}}`,
			template: `#foreach($item in $ctx.args.items)
$item
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				r := vtlResult(out)
				if !strings.Contains(r, "a") || !strings.Contains(r, "b") || !strings.Contains(r, "c") {
					return fmt.Errorf("expected all items, got %s", r)
				}
				return nil
			},
		},
		{
			name:    "VTL_Foreach_Count",
			context: `{"args":{"items":["x","y","z"]}}`,
			template: `#foreach($item in $ctx.args.items)
$foreach.count:$item
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				r := vtlResult(out)
				if !strings.Contains(r, "1:") || !strings.Contains(r, "2:") || !strings.Contains(r, "3:") {
					return fmt.Errorf("expected foreach count 1-3, got %s", r)
				}
				return nil
			},
		},
		{
			name:    "VTL_NestedIf",
			context: `{"args":{"role":"admin","active":true}}`,
			template: `#if($ctx.args.role == "admin")
#if($ctx.args.active)
admin-active
#end
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "admin-active") {
					return fmt.Errorf("expected admin-active, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_Composite_DynamoDBGetItem",
			context:  `{"args":{"id":"post-1"}}`,
			template: `{"version":"2018-05-29","operation":"GetItem","key":{"id":{"S":"$ctx.args.id"}}}`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				if m["version"] != "2018-05-29" {
					return fmt.Errorf("expected version=2018-05-29, got %v", m["version"])
				}
				if m["operation"] != "GetItem" {
					return fmt.Errorf("expected operation=GetItem, got %v", m["operation"])
				}
				return nil
			},
		},
		{
			name:    "VTL_PipelineStash_Passing",
			context: `{"args":{"postId":"p100"},"stash":{}}`,
			template: `#set($ctx.stash.postId = $ctx.args.postId)
{"version":"2018-05-29","operation":"GetItem","key":{"id":{"S":"$ctx.args.postId"}}}`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var s map[string]interface{}
				if e := json.Unmarshal([]byte(vtlStash(out)), &s); e != nil {
					return fmt.Errorf("stash is not valid JSON: %s", vtlStash(out))
				}
				if s["postId"] != "p100" {
					return fmt.Errorf("expected stash.postId=p100, got %v", s["postId"])
				}
				if !strings.Contains(vtlResult(out), "p100") {
					return fmt.Errorf("expected p100 in result, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilMathEquals",
			context: `{"args":{"a":5,"b":5}}`,
			template: `#if($util.math.equals($ctx.args.a, $ctx.args.b))
equal
#end`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "equal") {
					return fmt.Errorf("expected equal, got %s", vtlResult(out))
				}
				return nil
			},
		},
		{
			name:     "VTL_ResponseTemplate_UseResult",
			context:  `{"result":{"id":"post-1","title":"Test Post","author":"Alice"}}`,
			template: `{"id":"$ctx.result.id","title":"$ctx.result.title","author":"$ctx.result.author"}`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				if m["id"] != "post-1" {
					return fmt.Errorf("expected id=post-1, got %v", m["id"])
				}
				if m["author"] != "Alice" {
					return fmt.Errorf("expected author=Alice, got %v", m["author"])
				}
				return nil
			},
		},
		{
			name:    "VTL_UtilAppendError_NonFatal",
			context: `{}`,
			template: `$util.appendError("warning message", "WarningType")
still-running`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				if !strings.Contains(vtlResult(out), "still-running") {
					return fmt.Errorf("template should continue after appendError, got %s", vtlResult(out))
				}
				if out.OutErrors == nil || !strings.Contains(*out.OutErrors, "warning message") {
					return fmt.Errorf("expected outErrors to contain 'warning message', got %v", out.OutErrors)
				}
				return nil
			},
		},
		{
			name:    "VTL_MultipleSet_Composite",
			context: `{"args":{"firstName":"John","lastName":"Doe"}}`,
			template: `#set($first = $ctx.args.firstName)
#set($last = $ctx.args.lastName)
#set($full = "$first $last")
{"fullName":"$full"}`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var m map[string]interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &m); e != nil {
					return fmt.Errorf("result is not valid JSON: %s", vtlResult(out))
				}
				if m["fullName"] != "John Doe" {
					return fmt.Errorf("expected fullName=John Doe, got %v", m["fullName"])
				}
				return nil
			},
		},
		{
			name:    "VTL_Foreach_BuildJSONArray",
			context: `{"args":{"ids":["id1","id2","id3"]}}`,
			template: `#set($output = [])
#foreach($id in $ctx.args.ids)
#set($added = $output.add($id))
#end
$util.toJson($output)`,
			validate: func(out *appsync.EvaluateMappingTemplateOutput) error {
				var arr []interface{}
				if e := json.Unmarshal([]byte(vtlResult(out)), &arr); e != nil {
					return fmt.Errorf("result is not valid JSON array: %s", vtlResult(out))
				}
				if len(arr) != 3 {
					return fmt.Errorf("expected 3 items, got %d", len(arr))
				}
				return nil
			},
		},
	}

	for _, t := range tests {
		tc := t
		results = append(results, r.RunTest("appsync", tc.name, func() error {
			out, evalErr := evalVTL(ctx, client, tc.context, tc.template)
			if evalErr != nil {
				return evalErr
			}
			return tc.validate(out)
		}))
	}

	return results
}

func evalVTL(ctx context.Context, client *appsync.Client, contextJSON, template string) (*appsync.EvaluateMappingTemplateOutput, error) {
	out, err := client.EvaluateMappingTemplate(ctx, &appsync.EvaluateMappingTemplateInput{
		Context:  aws.String(contextJSON),
		Template: aws.String(template),
	})
	if err != nil {
		return nil, fmt.Errorf("EvaluateMappingTemplate API error: %w", err)
	}
	return out, nil
}
