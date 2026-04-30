package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
)

type vtlTest struct {
	name         string
	template     string
	body         string
	validate     func(body string) error
	respTmpl     string
}

func (r *TestRunner) invokeMock(ctx context.Context, client *apigateway.Client, apiID, pathPart, template, body string, respTmpl string) (string, error) {
	resResp, err := client.CreateResource(ctx, &apigateway.CreateResourceInput{
		RestApiId: aws.String(apiID),
		ParentId:  aws.String(apiID),
		PathPart:  aws.String(pathPart),
	})
	if err != nil {
		return "", fmt.Errorf("create resource: %w", err)
	}

	_, err = client.PutMethod(ctx, &apigateway.PutMethodInput{
		RestApiId:         aws.String(apiID),
		ResourceId:        resResp.Id,
		HttpMethod:        aws.String("POST"),
		AuthorizationType: aws.String("NONE"),
	})
	if err != nil {
		return "", fmt.Errorf("put method: %w", err)
	}

	putIntInput := &apigateway.PutIntegrationInput{
		RestApiId:  aws.String(apiID),
		ResourceId: resResp.Id,
		HttpMethod: aws.String("POST"),
		Type:       types.IntegrationTypeMock,
		RequestTemplates: map[string]string{
			"application/json": template,
		},
	}
	_, err = client.PutIntegration(ctx, putIntInput)
	if err != nil {
		return "", fmt.Errorf("put integration: %w", err)
	}

	if respTmpl != "" {
		_, err = client.PutIntegrationResponse(ctx, &apigateway.PutIntegrationResponseInput{
			RestApiId:         aws.String(apiID),
			ResourceId:        resResp.Id,
			HttpMethod:        aws.String("POST"),
			StatusCode:        aws.String("200"),
			ResponseTemplates: map[string]string{"application/json": respTmpl},
		})
		if err != nil {
			return "", fmt.Errorf("put integration response: %w", err)
		}
	}

	testBody := ""
	if body != "" {
		testBody = body
	}
	resp, err := client.TestInvokeMethod(ctx, &apigateway.TestInvokeMethodInput{
		RestApiId:  aws.String(apiID),
		ResourceId: resResp.Id,
		HttpMethod: aws.String("POST"),
		Body:       aws.String(testBody),
	})
	if err != nil {
		return "", fmt.Errorf("test invoke: %w", err)
	}
	if resp.Status != 200 {
		return "", fmt.Errorf("status %d", resp.Status)
	}
	if resp.Body == nil {
		return "", fmt.Errorf("body is nil")
	}
	return *resp.Body, nil
}

func (r *TestRunner) vtlTests(ctx context.Context, client *apigateway.Client, apiID string) []TestResult {
	var results []TestResult

	tests := []vtlTest{
		{
			name:     "VTL_InputJson_Root",
			template: `$input.json('$')`,
			body:     `{"name": "Alice", "age": 30}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("not valid JSON: %s", body)
				}
				if parsed["name"] != "Alice" {
					return fmt.Errorf("$.name should be Alice, got %v", parsed["name"])
				}
				if parsed["age"] != float64(30) {
					return fmt.Errorf("$.age should be 30, got %v", parsed["age"])
				}
				return nil
			},
		},
		{
			name:     "VTL_InputJson_Path",
			template: `$input.json('$.name')`,
			body:     `{"name": "Bob", "age": 25}`,
			validate: func(body string) error {
				var parsed interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				str, ok := parsed.(string)
				if !ok || str != "Bob" {
					return fmt.Errorf("expected JSON string \"Bob\", got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_InputPath_Field",
			template: `$input.path('$.city')`,
			body:     `{"name": "Carol", "city": "Tokyo"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "Tokyo" {
					return fmt.Errorf("expected exactly 'Tokyo', got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_InputBody",
			template: `$input.body`,
			body:     `{"raw": "payload"}`,
			validate: func(body string) error {
				if !strings.Contains(body, `"raw": "payload"`) && !strings.Contains(body, `"raw":"payload"`) {
					return fmt.Errorf("expected raw body payload, got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_InputJson_Nested",
			template: `$input.json('$.address.city')`,
			body:     `{"name": "Dave", "address": {"city": "Osaka", "zip": "530"}}`,
			validate: func(body string) error {
				var parsed interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				str, ok := parsed.(string)
				if !ok || str != "Osaka" {
					return fmt.Errorf("expected JSON string \"Osaka\", got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_InputJson_ArrayElement",
			template: `$input.json('$.items[0]')`,
			body:     `{"items": ["alpha", "beta", "gamma"]}`,
			validate: func(body string) error {
				var parsed interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				str, ok := parsed.(string)
				if !ok || str != "alpha" {
					return fmt.Errorf("expected JSON string \"alpha\", got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_Context_ApiId",
			template: `{"apiId": "$context.apiId"}`,
			body:     `{}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				got, ok := parsed["apiId"].(string)
				if !ok {
					return fmt.Errorf("apiId is not a string: %v", parsed["apiId"])
				}
				if got != apiID {
					return fmt.Errorf("expected apiId %s, got: %s", apiID, got)
				}
				return nil
			},
		},
		{
			name:     "VTL_Context_Method",
			template: `{"method": "$context.method"}`,
			body:     `{}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				got, ok := parsed["method"].(string)
				if !ok {
					return fmt.Errorf("method is not a string: %v", parsed["method"])
				}
				if got != "POST" {
					return fmt.Errorf("expected method POST, got: %s", got)
				}
				return nil
			},
		},
		{
			name:     "VTL_Context_Stage",
			template: `{"stage": "$context.stage"}`,
			body:     `{}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				stage, ok := parsed["stage"].(string)
				if !ok {
					return fmt.Errorf("stage is not a string: %v", parsed["stage"])
				}
				if stage == "$context.stage" {
					return fmt.Errorf("$context.stage was not resolved at all, raw template left: %s", stage)
				}
				return nil
			},
		},
		{
			name:     "VTL_Set_Variable",
			template: `#set($greeting = "Hello World"){"msg":"$greeting"}`,
			body:     `{}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				got, ok := parsed["msg"].(string)
				if !ok {
					return fmt.Errorf("msg is not a string: %v", parsed["msg"])
				}
				if got != "Hello World" {
					return fmt.Errorf("expected msg \"Hello World\", got: %s", got)
				}
				return nil
			},
		},
		{
			name:     "VTL_Set_FromInput",
			template: `#set($name = $input.path('$.name')){"greeting":"Hello $name"}`,
			body:     `{"name": "Eve"}`,
			validate: func(body string) error {
				if !strings.Contains(body, "Eve") {
					return fmt.Errorf("expected Eve in output, got: %s", body)
				}
				if !strings.Contains(body, "Hello") {
					return fmt.Errorf("expected Hello in output, got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_If_TrueBranch",
			template: `#set($val = $input.path('$.active'))#if($val == 'true')ACTIVE#end`,
			body:     `{"active": "true"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "ACTIVE" {
					return fmt.Errorf("expected exactly 'ACTIVE', got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_If_FalseOmitted",
			template: `#set($val = $input.path('$.active'))#if($val == 'true')ISACTIVE#endDONE`,
			body:     `{"active": "false"}`,
			validate: func(body string) error {
				if strings.Contains(body, "ISACTIVE") {
					return fmt.Errorf("ISACTIVE should not appear when condition false, got: %s", body)
				}
				if !strings.Contains(body, "DONE") {
					return fmt.Errorf("expected DONE outside #if, got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_IfElse_TrueBranch",
			template: `#set($role = $input.path('$.role'))#if($role == 'admin')ADMIN#elseUSER#end`,
			body:     `{"role": "admin"}`,
			validate: func(body string) error {
				if !strings.Contains(body, "ADMIN") {
					return fmt.Errorf("expected ADMIN from #if true, got: %s", body)
				}
				if strings.Contains(body, "USER") {
					return fmt.Errorf("USER should not appear in true branch, got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_IfElse_FalseBranch",
			template: `#set($role = $input.path('$.role'))#if($role == 'admin')ADMIN#elseUSER#end`,
			body:     `{"role": "viewer"}`,
			validate: func(body string) error {
				if !strings.Contains(body, "USER") {
					return fmt.Errorf("expected USER from #else, got: %s", body)
				}
				if strings.Contains(body, "ADMIN") {
					return fmt.Errorf("ADMIN should not appear in false branch, got: %s", body)
				}
				return nil
			},
		},
		{
			name:     "VTL_Foreach_Iteration",
			template: `#set($items = $input.json('$.items'))#foreach($item in $items)$item #end`,
			body:     `{"items": ["a", "b", "c"]}`,
			validate: func(body string) error {
				parts := strings.Fields(strings.TrimSpace(body))
				if len(parts) != 3 {
					return fmt.Errorf("expected exactly 3 space-separated items, got %d: '%s'", len(parts), body)
				}
				if parts[0] != "a" || parts[1] != "b" || parts[2] != "c" {
					return fmt.Errorf("expected items [a b c] in order, got: %v", parts)
				}
				return nil
			},
		},
		{
			name:     "VTL_Util_UrlEncode",
			template: `$util.urlEncode("hello world")`,
			body:     `{}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "hello+world" && trimmed != "hello%20world" {
					return fmt.Errorf("expected URL-encoded 'hello world', got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_Util_Base64Encode",
			template: `$util.base64Encode("test")`,
			body:     `{}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "dGVzdA==" {
					return fmt.Errorf("expected base64 'dGVzdA==', got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_Composite_Template",
			template: `{"name":"$input.path('$.name')","api":"$context.apiId","version":"2"}`,
			body:     `{"name": "Frank"}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				name, _ := parsed["name"].(string)
				if name != "Frank" {
					return fmt.Errorf("expected name \"Frank\", got: %s", name)
				}
				api, _ := parsed["api"].(string)
				if api != apiID {
					return fmt.Errorf("expected api \"%s\", got: %s", apiID, api)
				}
				ver, _ := parsed["version"].(string)
				if ver != "2" {
					return fmt.Errorf("expected version \"2\", got: %s", ver)
				}
				return nil
			},
		},
		{
			name:     "VTL_ResponseTemplate",
			template: `{"name":"$input.path('$.name')"}`,
			body:     `{"name": "Grace"}`,
			respTmpl: `{"wrapped": $input.json('$'), "status": "ok"}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				status, _ := parsed["status"].(string)
				if status != "ok" {
					return fmt.Errorf("expected status \"ok\", got: %v", parsed["status"])
				}
				wrapped, ok := parsed["wrapped"].(map[string]interface{})
				if !ok {
					return fmt.Errorf("expected wrapped to be an object, got: %v", parsed["wrapped"])
				}
				if wrapped["name"] != "Grace" {
					return fmt.Errorf("expected wrapped.name \"Grace\", got: %v", wrapped["name"])
				}
				return nil
			},
		},
		{
			name:     "VTL_IfElseIfElse_AllBranches",
			template: `#set($level = $input.path('$.level'))#if($level == 'high')HIGH#elseif($level == 'mid')MID#elseLOW#end`,
			body:     `{"level": "mid"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "MID" {
					return fmt.Errorf("expected 'MID' from #elseif, got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_IfElseIfElse_DefaultBranch",
			template: `#set($level = $input.path('$.level'))#if($level == 'high')HIGH#elseif($level == 'mid')MID#elseLOW#end`,
			body:     `{"level": "unknown"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "LOW" {
					return fmt.Errorf("expected 'LOW' from #else (default), got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_NestedIf",
			template: `#set($a = $input.path('$.a'))#set($b = $input.path('$.b'))#if($a == '1')A#if($b == '2')AND_B#end#end`,
			body:     `{"a": "1", "b": "2"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "AAND_B" {
					return fmt.Errorf("expected 'AAND_B' from nested #if, got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_NestedIf_OuterTrueInnerFalse",
			template: `#set($a = $input.path('$.a'))#set($b = $input.path('$.b'))#if($a == '1')A#if($b == '2')AND_B#end#end`,
			body:     `{"a": "1", "b": "9"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "A" {
					return fmt.Errorf("expected 'A' only (inner false), got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_Foreach_Count",
			template: `#set($items = $input.json('$.items'))#foreach($item in $items)[$foreach.count]$item#end`,
			body:     `{"items": ["x", "y"]}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if !strings.Contains(trimmed, "[1]x") || !strings.Contains(trimmed, "[2]y") {
					return fmt.Errorf("expected foreach count [1]x [2]y, got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_MultipleSet",
			template: `#set($a = $input.path('$.x'))#set($b = $input.path('$.y')){"sum":"$a+$b"}`,
			body:     `{"x": "10", "y": "20"}`,
			validate: func(body string) error {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(body), &parsed); err != nil {
					return fmt.Errorf("output is not valid JSON: %s", body)
				}
				sum, ok := parsed["sum"].(string)
				if !ok {
					return fmt.Errorf("sum is not a string: %v", parsed["sum"])
				}
				if sum != "10+20" {
					return fmt.Errorf("expected sum \"10+20\", got: %s", sum)
				}
				return nil
			},
		},
		{
			name:     "VTL_EmptyBody",
			template: `$input.body`,
			body:     ``,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "" {
					return fmt.Errorf("expected empty output for empty body, got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_Set_BooleanCondition",
			template: `#set($active = $input.path('$.active'))#if($active == 'true')ON#end`,
			body:     `{"active": "true"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "ON" {
					return fmt.Errorf("expected 'ON' for boolean string compare, got: '%s'", trimmed)
				}
				return nil
			},
		},
		{
			name:     "VTL_InputJson_MissingPath",
			template: `$input.json('$.nonexistent')`,
			body:     `{"name": "test"}`,
			validate: func(body string) error {
				trimmed := strings.TrimSpace(body)
				if trimmed != "" && trimmed != "null" {
					return fmt.Errorf("expected empty or null for missing path, got: '%s'", trimmed)
				}
				return nil
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		testName := "TestInvokeMethod_" + tc.name
		results = append(results, r.RunTest("apigateway", testName, func() error {
			if apiID == "" {
				return fmt.Errorf("API ID not available")
			}
			pathPart := fmt.Sprintf("vtl-%d", time.Now().UnixNano())
			body, err := r.invokeMock(ctx, client, apiID, pathPart, tc.template, tc.body, tc.respTmpl)
			if err != nil {
				return err
			}
			return tc.validate(body)
		}))
	}

	return results
}
