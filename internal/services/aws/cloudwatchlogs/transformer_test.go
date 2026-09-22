package cloudwatchlogs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"vorpalstacks/internal/common/request"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The E2 pins: the recipe engine against the processor pages' own
// documented input→output examples, the CRUD validation rows, and the
// application seam (group-level override, account-policy prefix
// selection, the original-vs-transformed read split, and the
// applyOnTransformedLogs filters).

// transformOf runs one recipe over one message and parses the JSON
// output for example comparison.
func transformOf(t *testing.T, config string, message string) map[string]interface{} {
	t.Helper()
	var parsed []map[string]interface{}
	if err := json.Unmarshal([]byte(config), &parsed); err != nil {
		t.Fatalf("recipe does not parse: %v", err)
	}
	out, _ := transformEvent(parsed, message, TransformContext{
		LogGroup:  "myLogGroupName",
		LogStream: "myLogStreamName",
		AccountID: "012345678912",
		Region:    "us-east-1",
	})
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(out), &record); err != nil {
		t.Fatalf("output is not JSON (%q): %v", out, err)
	}
	return record
}

func TestTransformEngineDocumentedExamples(t *testing.T) {
	t.Run("parseJSON destination", func(t *testing.T) {
		got := transformOf(t, `[{"parseJSON":{"destination":"new_key"}}]`,
			`{"outer_key": {"inner_key": "inner_value"}}`)
		inner, _ := got["new_key"].(map[string]interface{})
		if inner == nil || inner["outer_key"] == nil {
			t.Fatalf("parseJSON destination example: %v", got)
		}
	})

	t.Run("grok extracts unstructured fields", func(t *testing.T) {
		got := transformOf(t, `[{"grok":{"match":"%{NUMBER:version} %{HOSTNAME:hostname} %{NOTSPACE:status} %{QUOTEDSTRING:logMsg}"}}]`,
			`293750 server-01.internal-network.local OK "[Thread-000] token generated"`)
		if got["version"] != "293750" || got["hostname"] != "server-01.internal-network.local" ||
			got["status"] != "OK" || got["logMsg"] != "[Thread-000] token generated" {
			t.Fatalf("grok example 1: %v", got)
		}
	})

	t.Run("grok dotted field names nest", func(t *testing.T) {
		got := transformOf(t, `[{"grok":{"match":"%{IP:client.ip} %{WORD:method} %{URIPATHPARAM:request.uri} %{NUMBER:response.status} %{NUMBER:response.bytes}"}}]`,
			`192.168.1.1 GET /index.html?param=value 200 1234`)
		client, _ := got["client"].(map[string]interface{})
		resp, _ := got["response"].(map[string]interface{})
		reqPath, _ := got["request"].(map[string]interface{})
		if client == nil || client["ip"] != "192.168.1.1" || got["method"] != "GET" ||
			reqPath == nil || reqPath["uri"] != "/index.html?param=value" ||
			resp == nil || resp["status"] != "200" || resp["bytes"] != "1234" {
			t.Fatalf("grok dotted example: %v", got)
		}
	})

	t.Run("grok apache common log typed output", func(t *testing.T) {
		got := transformOf(t, `[{"grok":{"match":"%{APACHE_ACCESS_LOG}"}}]`,
			`127.0.0.1 - - [03/Aug/2023:12:34:56 +0000] "GET /page.html HTTP/1.1" 200 1234`)
		if got["request"] != "/page.html" || got["http_method"] != "GET" ||
			got["http_version"] != "1.1" || got["remote_host"] != "127.0.0.1" {
			t.Fatalf("apache example: %v", got)
		}
		if got["status_code"] != float64(200) || got["response_size"] != float64(1234) {
			t.Fatalf("apache typed fields: %v", got)
		}
		// The documented output's timestamp field converts the captured
		// httpdate to ISO-8601 UTC ("timestamp": "2023-08-03T12:34:56Z")
		// — no raw httpdate field remains.
		if got["timestamp"] != "2023-08-03T12:34:56Z" {
			t.Fatalf("apache timestamp: %v", got)
		}
		if _, still := got["httpdate"]; still {
			t.Fatalf("apache raw httpdate field leaked: %v", got)
		}
	})

	t.Run("grok nginx common log", func(t *testing.T) {
		got := transformOf(t, `[{"grok":{"match":"%{NGINX_ACCESS_LOG}"}}]`,
			`192.168.1.100 - Foo [03/Aug/2023:12:34:56 +0000] "GET /account/login.html HTTP/1.1" 200 42 "https://www.amazon.com/" "Mozilla/5.0 (Windows NT 10.0)"`)
		if got["auth_user"] != "Foo" || got["request"] != "/account/login.html" ||
			got["status_code"] != float64(200) || got["response_size"] != float64(42) ||
			got["referrer"] != "https://www.amazon.com/" {
			t.Fatalf("nginx example: %v", got)
		}
		if got["timestamp"] != "2023-08-03T12:34:56Z" {
			t.Fatalf("nginx timestamp: %v", got)
		}
	})

	t.Run("grok common log non-UTC offset normalises", func(t *testing.T) {
		got := transformOf(t, `[{"grok":{"match":"%{APACHE_ACCESS_LOG}"}}]`,
			`127.0.0.1 - - [03/Aug/2023:12:34:56 +0640] "GET /page.html HTTP/1.1" 200 1234`)
		if got["timestamp"] != "2023-08-03T05:54:56Z" {
			t.Fatalf("non-UTC offset normalisation: %v", got)
		}
	})

	t.Run("grok httpdate one-digit day converts", func(t *testing.T) {
		// The HTTPDATE grammar admits a one-digit day (`0?[1-9]`); the
		// conversion layout must accept what the pattern matches.
		got := transformOf(t, `[{"grok":{"match":"%{APACHE_ACCESS_LOG}"}}]`,
			`127.0.0.1 - - [3/Aug/2023:10:00:00 +0000] "GET /page.html HTTP/1.1" 200 1234`)
		if got["timestamp"] != "2023-08-03T10:00:00Z" {
			t.Fatalf("one-digit day timestamp: %v", got)
		}
	})

	t.Run("grok syslog5424", func(t *testing.T) {
		got := transformOf(t, `[{"grok":{"match":"%{SYSLOG5424}"}}]`,
			`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource= "Application" eventID="1011"][examplePriority@32473 class="high"]`)
		if got["pri"] != float64(165) || got["version"] != float64(1) ||
			got["hostname"] != "mymachine.example.com" || got["app"] != "evntslog" ||
			got["msg_id"] != "ID47" || got["message"] != "[examplePriority@32473 class=\"high\"]" {
			t.Fatalf("syslog5424 example: %v", got)
		}
		if !strings.Contains(got["structured_data"].(string), "exampleSDID@32473") {
			t.Fatalf("syslog5424 structured data: %v", got["structured_data"])
		}
		// The sample's PROCID is RFC 5424's nil token and the documented
		// output omits procid entirely.
		if _, exists := got["procid"]; exists {
			t.Fatalf("nil PROCID emitted: %v", got)
		}
		// A present PROCID is captured.
		got = transformOf(t, `[{"grok":{"match":"%{SYSLOG5424}"}}]`,
			`<34>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog 8710 ID47 -`)
		if got["procid"] != "8710" {
			t.Fatalf("present PROCID: %v", got)
		}
	})

	t.Run("grok after parseJSON on a field", func(t *testing.T) {
		got := transformOf(t, `[{"parseJSON":{}},{"grok":{"source":"logMsg","match":"%{WORD:http_method} %{NOTSPACE:request} HTTP/%{NUMBER:http_version}"}}]`,
			`{"timestamp": "2024-11-23T16:03:12Z", "level": "ERROR", "logMsg": "GET /page.html HTTP/1.1"}`)
		if got["http_method"] != "GET" || got["request"] != "/page.html" || got["http_version"] != "1.1" {
			t.Fatalf("grok+parseJSON example: %v", got)
		}
	})

	t.Run("csv with delimiter and quotes", func(t *testing.T) {
		got := transformOf(t, `[{"csv":{"delimiter":":","quoteCharacter":"'"}}]`,
			`'Akua Mansa':28:'New York: USA'`)
		if got["column_1"] != "Akua Mansa" || got["column_2"] != "28" || got["column_3"] != "New York: USA" {
			t.Fatalf("csv example 1: %v", got)
		}
	})

	t.Run("csv columns and destination", func(t *testing.T) {
		got := transformOf(t, `[{"parseJSON":{}},{"csv":{"source":"logMsg","delimiter":":","quoteCharacter":"'","columns":["name","age","location"],"destination":"msg"}}]`,
			`{"timestamp": "2024-11-23T16:03:12Z", "type": "user_data", "logMsg": "'Akua Mansa':28:'New York: USA'"}`)
		msg, _ := got["msg"].(map[string]interface{})
		if msg == nil || msg["name"] != "Akua Mansa" || msg["age"] != "28" || msg["location"] != "New York: USA" {
			t.Fatalf("csv example 2: %v", got)
		}
	})

	t.Run("parseKeyValue with non-match and prefix", func(t *testing.T) {
		got := transformOf(t, `[{"parseKeyValue":{"destination":"new_key","fieldDelimiter":"!","keyValueDelimiter":":","nonMatchValue":"defaultValue","keyPrefix":"parsed_"}}]`,
			`key1:value1!key2:value2!key3:value3!key4`)
		nk, _ := got["new_key"].(map[string]interface{})
		if nk == nil || nk["parsed_key1"] != "value1" || nk["parsed_key4"] != "defaultValue" {
			t.Fatalf("parseKeyValue example: %v", got)
		}
	})

	t.Run("json mutate examples", func(t *testing.T) {
		base := `[{"parseJSON":{}}`
		got := transformOf(t, base+`,{"addKeys":{"entries":[{"key":"outer_key.new_key","value":"new_value"}]}}]`,
			`{"outer_key": {"inner_key": "inner_value"}}`)
		outer, _ := got["outer_key"].(map[string]interface{})
		if outer == nil || outer["new_key"] != "new_value" {
			t.Fatalf("addKeys example: %v", got)
		}

		got = transformOf(t, base+`,{"deleteKeys":{"withKeys":["outer_key.inner_key"]}}]`,
			`{"outer_key": {"inner_key": "inner_value"}}`)
		outer, _ = got["outer_key"].(map[string]interface{})
		if outer == nil || len(outer) != 0 {
			t.Fatalf("deleteKeys example: %v", got)
		}

		got = transformOf(t, base+`,{"moveKeys":{"entries":[{"source":"outer_key1.inner_key1","target":"outer_key2"}]}}]`,
			`{"outer_key1": {"inner_key1": "inner_value1"}, "outer_key2": {"inner_key2": "inner_value2"}}`)
		outer2, _ := got["outer_key2"].(map[string]interface{})
		outer1, _ := got["outer_key1"].(map[string]interface{})
		if outer2 == nil || outer2["inner_key1"] != "inner_value1" || len(outer1) != 0 {
			t.Fatalf("moveKeys example: %v", got)
		}

		got = transformOf(t, base+`,{"renameKeys":{"entries":[{"key":"outer_key","renameTo":"new_key"}]}}]`,
			`{"outer_key": {"inner_key": "inner_value"}}`)
		if got["outer_key"] != nil || got["new_key"] == nil {
			t.Fatalf("renameKeys example: %v", got)
		}

		// A move or rename that cannot overwrite (the destination
		// exists, overwriteIfExists false by default) keeps the value
		// under its original key — the write is blocked, so the source
		// is never deleted around it.
		got = transformOf(t, base+`,{"moveKeys":{"entries":[{"source":"outer_key1.inner_key1","target":"outer_key2"}]}}]`,
			`{"outer_key1": {"inner_key1": "inner_value1"}, "outer_key2": {"inner_key1": "occupied"}}`)
		outer1, _ = got["outer_key1"].(map[string]interface{})
		outer2, _ = got["outer_key2"].(map[string]interface{})
		if outer1 == nil || outer1["inner_key1"] != "inner_value1" ||
			outer2 == nil || outer2["inner_key1"] != "occupied" {
			t.Fatalf("blocked moveKeys keeps the source: %v", got)
		}
		got = transformOf(t, base+`,{"renameKeys":{"entries":[{"key":"old_key","renameTo":"new_key"}]}}]`,
			`{"old_key": "old_value", "new_key": "occupied"}`)
		if got["old_key"] != "old_value" || got["new_key"] != "occupied" {
			t.Fatalf("blocked renameKeys keeps the source: %v", got)
		}
		// With the overwrite flag the blocked move succeeds.
		got = transformOf(t, base+`,{"moveKeys":{"entries":[{"source":"outer_key1.inner_key1","target":"outer_key2","overwriteIfExists":true}]}}]`,
			`{"outer_key1": {"inner_key1": "inner_value1"}, "outer_key2": {"inner_key1": "occupied"}}`)
		outer1, _ = got["outer_key1"].(map[string]interface{})
		outer2, _ = got["outer_key2"].(map[string]interface{})
		if outer1 == nil || len(outer1) != 0 || outer2 == nil || outer2["inner_key1"] != "inner_value1" {
			t.Fatalf("overwriting moveKeys relocates the value: %v", got)
		}

		got = transformOf(t, base+`,{"copyValue":{"entries":[{"source":"outer_key.inner_key","target":"new_key"}]}}]`,
			`{"outer_key": {"inner_key": "inner_value"}}`)
		if got["new_key"] != "inner_value" {
			t.Fatalf("copyValue example: %v", got)
		}

		// The documented metadata form: the copyValue sources
		// @logGroupName/@logGroupStream/@accountId/@regionName read the
		// event's platform identity instead of the record. The doc
		// section's first entry spells its source outer_key.new_key while
		// its own input carries outer_key.inner_key and its output shows
		// new_key = inner_value, so the pin follows the output-consistent
		// spelling of that record-path entry.
		got = transformOf(t, base+`,{"copyValue":{"entries":[{"source":"outer_key.inner_key","target":"new_key"},{"source":"@logGroupName","target":"log_group_name"},{"source":"@logGroupStream","target":"log_group_stream"},{"source":"@accountId","target":"account_id"},{"source":"@regionName","target":"region_name"}]}}]`,
			`{"outer_key": {"inner_key": "inner_value"}}`)
		if got["new_key"] != "inner_value" ||
			got["log_group_name"] != "myLogGroupName" || got["log_group_stream"] != "myLogStreamName" ||
			got["account_id"] != "012345678912" || got["region_name"] != "us-east-1" {
			t.Fatalf("copyValue metadata example: %v", got)
		}
	})

	t.Run("listToMap flatten forms", func(t *testing.T) {
		source := `[{"parseJSON":{}},{"listToMap":{"source":"outer_key","key":"inner_key","valueKey":"inner_value","flatten":false}}]`
		msg := `{"outer_key": [{"inner_key": "a", "inner_value": "val-a"}, {"inner_key": "b", "inner_value": "val-b1"}, {"inner_key": "b", "inner_value": "val-b2"}, {"inner_key": "c", "inner_value": "val-c"}]}`
		got := transformOf(t, source, msg)
		a, _ := got["a"].([]interface{})
		b, _ := got["b"].([]interface{})
		if a == nil || a[0] != "val-a" || b == nil || len(b) != 2 || b[1] != "val-b2" {
			t.Fatalf("listToMap list form: %v", got)
		}
		got = transformOf(t, `[{"parseJSON":{}},{"listToMap":{"source":"outer_key","key":"inner_key","valueKey":"inner_value","flatten":true,"flattenedElement":"first"}}]`, msg)
		if got["a"] != "val-a" || got["b"] != "val-b1" || got["c"] != "val-c" {
			t.Fatalf("listToMap flatten first: %v", got)
		}
		got = transformOf(t, `[{"parseJSON":{}},{"listToMap":{"source":"outer_key","key":"inner_key","valueKey":"inner_value","flatten":true,"flattenedElement":"last"}}]`, msg)
		if got["b"] != "val-b2" {
			t.Fatalf("listToMap flatten last: %v", got)
		}
	})

	t.Run("string mutate examples", func(t *testing.T) {
		base := `[{"parseJSON":{}}`
		got := transformOf(t, base+`,{"lowerCaseString":{"withKeys":["outer_key.inner_key"]}}]`,
			`{"outer_key": {"inner_key": "INNER_VALUE"}}`)
		outer, _ := got["outer_key"].(map[string]interface{})
		if outer["inner_key"] != "inner_value" {
			t.Fatalf("lowerCaseString example: %v", got)
		}

		got = transformOf(t, base+`,{"trimString":{"withKeys":["outer_key.inner_key"]}}]`,
			`{"outer_key": {"inner_key": "   inner_value  "}}`)
		outer, _ = got["outer_key"].(map[string]interface{})
		if outer["inner_key"] != "inner_value" {
			t.Fatalf("trimString example: %v", got)
		}

		got = transformOf(t, base+`,{"splitString":{"entries":[{"source":"outer_key.inner_key","delimiter":"_"}]}}]`,
			`{"outer_key": {"inner_key": "inner_value"}}`)
		outer, _ = got["outer_key"].(map[string]interface{})
		parts, _ := outer["inner_key"].([]interface{})
		if len(parts) != 2 || parts[0] != "inner" || parts[1] != "value" {
			t.Fatalf("splitString example: %v", got)
		}

		got = transformOf(t, base+`,{"substituteString":{"entries":[
			{"source":"outer_key.inner_key1","from":"\\[\\]","to":"value1"},
			{"source":"outer_key.inner_key2","from":"[0-9]{3}-[0-9]{3}-[0-9]{3}","to":"xxx-xxx-xxx"},
			{"source":"outer_key.inner_key3","from":"cat","to":"dog"}]}}]`,
			`{"outer_key": {"inner_key1": "[]", "inner_key2": "123-345-567", "inner_key3": "A cat takes a catnap."}}`)
		outer, _ = got["outer_key"].(map[string]interface{})
		if outer["inner_key1"] != "value1" || outer["inner_key2"] != "xxx-xxx-xxx" || outer["inner_key3"] != "A dog takes a dognap." {
			t.Fatalf("substituteString example: %v", got)
		}

		got = transformOf(t, base+`,{"substituteString":{"entries":[
			{"source":"outer_key.inner_key1","from":"(\\w+), (\\w+), and (\\w+)","to":"$1 and $3"},
			{"source":"outer_key.inner_key2","from":"^arn:aws:sts::(?P<account_id>\\d{12}):assumed-role/(?P<role_name>[\\w+=,.@-]+)/(?P<role_session_name>[\\w+=,.@-]+)$","to":"${account_id}:${role_name}:${role_session_name}"}]}}]`,
			`{"outer_key": {"inner_key1": "Tom, Dick, and Harry", "inner_key2": "arn:aws:sts::123456789012:assumed-role/MyImportantRole/MySession"}}`)
		outer, _ = got["outer_key"].(map[string]interface{})
		if outer["inner_key1"] != "Tom and Harry" ||
			outer["inner_key2"] != "123456789012:MyImportantRole:MySession" {
			t.Fatalf("substituteString backreference example: %v", got)
		}
	})

	t.Run("typeConverter", func(t *testing.T) {
		got := transformOf(t, `[{"parseJSON":{}},{"typeConverter":{"entries":[{"key":"status","type":"integer"}]}}]`,
			`{"name": "value", "status": "200"}`)
		if got["status"] != float64(200) {
			t.Fatalf("typeConverter example: %v", got)
		}

		// A numeric capture the grok composites write is int64 inside
		// the pipeline (JSON only flattens it at the seam): the boolean
		// cast must convert it, not silently skip it.
		got = transformOf(t, `[{"grok":{"match":"%{SYSLOG5424}"}},{"typeConverter":{"entries":[{"key":"pri","type":"boolean"}]}}]`,
			`<165>1 2003-10-11T22:14:15.003Z mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3"]`)
		if got["pri"] != true {
			t.Fatalf("int64 numeric capture to boolean: %v", got)
		}
	})

	t.Run("dateTimeConverter", func(t *testing.T) {
		got := transformOf(t, `[{"parseJSON":{}},{"dateTimeConverter":{"source":"german_datetime","target":"target_1","locale":"de","matchPatterns":["EEEE dd. MMMM yyyy HH:mm:ss"],"sourceTimezone":"Europe/Berlin","targetTimezone":"America/New_York","targetFormat":"yyyy-MM-dd'T'HH:mm:ss z"}}]`,
			`{"german_datetime": "Samstag 05. Dezember 1998 11:00:00"}`)
		rendered, _ := got["target_1"].(string)
		if !strings.HasPrefix(rendered, "1998-12-05T") {
			t.Fatalf("dateTimeConverter example: %v", got)
		}

		// The documented default targetFormat — yyyy-MM-dd'T'HH:mm:ss.SSS'Z'
		// — renders its quoted literals verbatim: no quote characters leak
		// and the trailing 'Z' stays the literal Z, never the zone token.
		got = transformOf(t, `[{"parseJSON":{}},{"dateTimeConverter":{"source":"iso_datetime","target":"target_2","matchPatterns":["yyyy-MM-dd'T'HH:mm:ss.SSS'Z'"]}}]`,
			`{"iso_datetime": "2024-06-01T12:30:45.123Z"}`)
		rendered, _ = got["target_2"].(string)
		if rendered != "2024-06-01T12:30:45.123Z" {
			t.Fatalf("dateTimeConverter default targetFormat render: %q", rendered)
		}
	})
}

// The numeric captures the grok composites write, and a prior integer
// conversion's own output, are int64 inside the pipeline: every numeric
// cast must accept an int64 source (the double and boolean switches once
// missed it and silently left the value unconverted).
func TestConvertTypeInt64Sources(t *testing.T) {
	if v, ok := convertType(int64(443), "double"); !ok || v != float64(443) {
		t.Fatalf("int64 to double: %v %v", v, ok)
	}
	if v, ok := convertType(int64(1), "boolean"); !ok || v != true {
		t.Fatalf("int64 one to boolean: %v %v", v, ok)
	}
	if v, ok := convertType(int64(0), "boolean"); !ok || v != false {
		t.Fatalf("int64 zero to boolean: %v %v", v, ok)
	}
	if v, ok := convertType(int64(7), "integer"); !ok || v != int64(7) {
		t.Fatalf("int64 to integer identity: %v %v", v, ok)
	}
}

// The CRUD rows: structural rules, per-processor member rows, the
// Standard-class gate, and the round trip.
func TestTransformerCRUDRows(t *testing.T) {
	svc, store := newReadTestService(t, "xform-group")

	valid := []map[string]interface{}{
		{"parseJSON": map[string]interface{}{}},
		{"addKeys": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"key": "env", "value": "test"},
		}}},
	}
	if _, err := svc.putTransformerCore("xform-group", valid, "us-east-1"); err != nil {
		t.Fatalf("put transformer: %v", err)
	}
	saved, err := svc.getTransformerCore("xform-group", "us-east-1")
	if err != nil {
		t.Fatalf("get transformer: %v", err)
	}
	if len(saved.Config) != 2 || saved.CreationTime == 0 || saved.LastModifiedTime == 0 {
		t.Fatalf("saved transformer: %+v", saved)
	}
	// An update keeps the creation stamp.
	time.Sleep(2 * time.Millisecond)
	if _, err := svc.putTransformerCore("arn:aws:logs:us-east-1:000000000000:log-group:xform-group", valid, "us-east-1"); err != nil {
		t.Fatalf("put transformer by ARN: %v", err)
	}
	updated, _ := svc.getTransformerCore("xform-group", "us-east-1")
	if updated.CreationTime != saved.CreationTime || updated.LastModifiedTime == saved.LastModifiedTime {
		t.Fatalf("update stamps: creation %d→%d, modified %d→%d",
			saved.CreationTime, updated.CreationTime, saved.LastModifiedTime, updated.LastModifiedTime)
	}

	if _, err := svc.getTransformerCore("no-such-group", "us-east-1"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("get missing transformer: %v", err)
	}
	if err := svc.deleteTransformerCore("xform-group", "us-east-1"); err != nil {
		t.Fatalf("delete transformer: %v", err)
	}
	if err := svc.deleteTransformerCore("xform-group", "us-east-1"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("delete missing transformer: %v", err)
	}

	cases := []struct {
		name   string
		config []map[string]interface{}
	}{
		{"empty config", nil},
		{"first processor not a parser", []map[string]interface{}{
			{"addKeys": map[string]interface{}{"entries": []interface{}{map[string]interface{}{"key": "k", "value": "v"}}}},
			{"parseJSON": map[string]interface{}{}},
		}},
		{"no parser", []map[string]interface{}{
			{"addKeys": map[string]interface{}{"entries": []interface{}{map[string]interface{}{"key": "k", "value": "v"}}}},
		}},
		{"two groks", []map[string]interface{}{
			{"grok": map[string]interface{}{"match": "%{WORD:w}"}},
			{"grok": map[string]interface{}{"match": "%{WORD:w2}"}},
		}},
		{"unknown processor", []map[string]interface{}{{"parseXML": map[string]interface{}{}}}},
		{"unpublished built-in parser", []map[string]interface{}{{"parseToOCSF": map[string]interface{}{
			"eventSource": "VPC_FLOW", "ocsfVersion": "V1_5",
		}}}},
		{"built-in parser with a non-message source", []map[string]interface{}{{"parseVPC": map[string]interface{}{"source": "payload"}}}},
		{"grok without match", []map[string]interface{}{{"grok": map[string]interface{}{}}}},
		{"addKeys without entries", []map[string]interface{}{{"parseJSON": map[string]interface{}{}}, {"addKeys": map[string]interface{}{}}}},
		{"addKeys entry without value", []map[string]interface{}{{"parseJSON": map[string]interface{}{}}, {"addKeys": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"key": "k"},
		}}}}},
		{"splitString entry without delimiter", []map[string]interface{}{{"parseJSON": map[string]interface{}{}}, {"splitString": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"source": "s"},
		}}}}},
		{"substituteString entry without from", []map[string]interface{}{{"parseJSON": map[string]interface{}{}}, {"substituteString": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"source": "s", "to": "t"},
		}}}}},
		{"substituteString entry without to", []map[string]interface{}{{"parseJSON": map[string]interface{}{}}, {"substituteString": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"source": "s", "from": "f"},
		}}}}},
		{"addKeys entry ceiling", []map[string]interface{}{{"parseJSON": map[string]interface{}{}}, {"addKeys": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"key": "k1", "value": "v"}, map[string]interface{}{"key": "k2", "value": "v"},
			map[string]interface{}{"key": "k3", "value": "v"}, map[string]interface{}{"key": "k4", "value": "v"},
			map[string]interface{}{"key": "k5", "value": "v"}, map[string]interface{}{"key": "k6", "value": "v"},
		}}}}},
	}
	for _, tc := range cases {
		_, err := svc.putTransformerCore("xform-group", tc.config, "us-east-1")
		if code := logsErrorCode(err); code != "InvalidParameterException" {
			t.Fatalf("%s: code=%q err=%v", tc.name, code, err)
		}
	}

	// A missing group rejects; an Infrequent-Access group rejects.
	if _, err := svc.putTransformerCore("no-such-group", valid, "us-east-1"); logsErrorCode(err) != "ResourceNotFoundException" {
		t.Fatalf("missing group: %v", err)
	}
	// The identifier carries the shared LogGroupIdentifier shape's
	// element traits: an identifier outside its alphabet rejects at
	// validation instead of answering as an unknown group.
	if _, err := svc.putTransformerCore("bad identifier!", valid, "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("invalid identifier: %v", err)
	}
	if _, err := svc.getTransformerCore("bad identifier!", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("invalid identifier (get): %v", err)
	}
	if err := svc.deleteTransformerCore("bad identifier!", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("invalid identifier (delete): %v", err)
	}
	_ = store
}

// The application seam: ingestion persists transformed copies, the API
// read plane keeps originals, the query plane reads transformed, and
// the group-level transformer overrides the account policy.
func TestTransformerApplicationAtIngestion(t *testing.T) {
	svc, store := newReadTestService(t, "app-group")

	recipe := []map[string]interface{}{
		{"parseJSON": map[string]interface{}{}},
		{"addKeys": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"key": "enriched", "value": "yes"},
		}}},
	}
	if _, err := svc.putTransformerCore("app-group", recipe, "us-east-1"); err != nil {
		t.Fatal(err)
	}
	const stream = "s1"
	if err := store.CreateLogStream(logsstore.NewLogStream(stream, "app-group")); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName: "app-group", LogStreamName: stream,
		Events: []PutLogEvent{
			{LogEntry: logsstore.LogEntry{Timestamp: now, Message: `{"level":"ERROR","msg":"boom"}`}, TimestampSet: true},
		},
		Region: "us-east-1",
	}); err != nil {
		t.Fatal(err)
	}

	// The API read plane keeps the original.
	events, err := fetchAllLogEvents(store, "app-group", stream, 0, 0)
	if err != nil || len(events) != 1 {
		t.Fatalf("read events: %v %d", err, len(events))
	}
	if events[0].Message != `{"level":"ERROR","msg":"boom"}` {
		t.Fatalf("original message replaced: %q", events[0].Message)
	}

	// The query plane reads the transformed form.
	transformed, err := store.TransformedMessageMap("app-group")
	if err != nil {
		t.Fatal(err)
	}
	digest := logsstore.TransformedMessageDigest(now, stream, events[0].Message)
	out, ok := transformed[digest]
	if !ok {
		t.Fatalf("transformed copy missing (digest %s, map %v)", digest, transformed)
	}
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(out), &record); err != nil {
		t.Fatal(err)
	}
	if record["level"] != "ERROR" || record["enriched"] != "yes" {
		t.Fatalf("transformed record: %v", record)
	}

	// A deleted transformer stops transformation (the cache invalidates
	// immediately).
	if err := svc.deleteTransformerCore("app-group", "us-east-1"); err != nil {
		t.Fatal(err)
	}
	if recipe := svc.effectiveTransformer(store, "us-east-1", "app-group"); recipe.config != nil {
		t.Fatalf("transformer survived deletion: %v", recipe.config)
	}
}

// The account-level TRANSFORMER_POLICY applies by LogGroupNamePrefix,
// and a group-level transformer overrides it.
func TestTransformerAccountPolicySelection(t *testing.T) {
	svc, store := newReadTestService(t, "/aws/lambda/group")

	policyDoc := `[{"parseJSON":{}},{"addKeys":{"entries":[{"key":"source","value":"account"}]}}]`
	if _, err := svc.putAccountPolicyCore("acct-xform", policyDoc, "TRANSFORMER_POLICY", "",
		"LogGroupNamePrefix = /aws/lambda", "us-east-1"); err != nil {
		t.Fatalf("put account policy: %v", err)
	}

	recipe := svc.effectiveTransformer(store, "us-east-1", "/aws/lambda/group")
	if recipe.config == nil {
		t.Fatal("account policy did not apply to the matching group")
	}
	if recipe := svc.effectiveTransformer(store, "us-east-1", "/aws/other"); recipe.config != nil {
		t.Fatalf("account policy applied outside its prefix: %v", recipe.config)
	}

	// An overlapping prefix rejects.
	_, err := svc.putAccountPolicyCore("acct-xform-2", policyDoc, "TRANSFORMER_POLICY", "",
		"LogGroupNamePrefix = /aws", "us-east-1")
	if logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("overlapping prefix: %v", err)
	}

	// The account-wide policy and prefix-scoped policies never coexist,
	// in either creation order: one policy for all log groups, or as
	// many as twenty scoped to subsets — never both forms at once.
	if _, err := svc.putAccountPolicyCore("acct-xform-wide", policyDoc, "TRANSFORMER_POLICY", "",
		"", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("account-wide over a prefix policy: %v", err)
	}

	// The group-level transformer overrides the account form.
	groupRecipe := []map[string]interface{}{
		{"parseJSON": map[string]interface{}{}},
		{"addKeys": map[string]interface{}{"entries": []interface{}{
			map[string]interface{}{"key": "source", "value": "group"},
		}}},
	}
	if _, err := svc.putTransformerCore("/aws/lambda/group", groupRecipe, "us-east-1"); err != nil {
		t.Fatal(err)
	}
	recipe = svc.effectiveTransformer(store, "us-east-1", "/aws/lambda/group")
	if recipe.config == nil || len(recipe.config) != 2 {
		t.Fatalf("group override did not take effect: %v", recipe.config)
	}

	// An invalid criteria rejects.
	if _, err := svc.putAccountPolicyCore("acct-xform-3", policyDoc, "TRANSFORMER_POLICY", "",
		"LogGroupName NOT IN []", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("subscription-form criteria on transformer policy: %v", err)
	}
}

// The applyOnTransformedLogs filters evaluate (and deliver) the
// transformed form.
func TestApplyOnTransformedLogsFilters(t *testing.T) {
	svc, store := newReadTestService(t, "xform-filter-group")
	const stream = "app"

	recipe := []map[string]interface{}{
		{"grok": map[string]interface{}{"match": `%{WORD:level} %{GREEDYDATA:rest}`}},
	}
	if _, err := svc.putTransformerCore("xform-filter-group", recipe, "us-east-1"); err != nil {
		t.Fatal(err)
	}
	// The pattern matches the TRANSFORMED form ("ERROR message" → level
	// extraction) but not the original JSON.
	if err := svc.putMetricFilterCore("xform-filter-group", "xform", "{ $.level = ERROR }", true,
		[]logsstore.MetricTransformation{{
			MetricName: "XformErrors", MetricNamespace: "Xform", MetricValue: "1",
			MetricNameSet: true, MetricNamespaceSet: true, MetricValueSet: true,
		}},
		true, "", nil, "us-east-1"); err != nil {
		t.Fatal(err)
	}

	metrics := &valueRecordingMetricInvoker{}
	svc.SetCloudWatchMetricInvoker(metrics)

	if err := store.CreateLogStream(logsstore.NewLogStream(stream, "xform-filter-group")); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UnixMilli()
	if _, err := svc.putLogEventsCore(PutLogEventsInput{
		LogGroupName: "xform-filter-group", LogStreamName: stream,
		Events: []PutLogEvent{
			{LogEntry: logsstore.LogEntry{Timestamp: now, Message: `ERROR something broke`}, TimestampSet: true},
		},
		Region: "us-east-1",
	}); err != nil {
		t.Fatal(err)
	}

	// The fan-out is asynchronous; wait for the recording invoker to
	// observe the emission.
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if calls, _ := metrics.recorded(); len(calls) > 0 {
			for _, c := range calls {
				if c.metric == "XformErrors" {
					return
				}
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("transformed-form metric never fired")
}

// TestTransformerSDKFacingVocabulary pins the transformer operations'
// response shapes through the handlers.
func TestTransformerSDKFacingVocabulary(t *testing.T) {
	svc, _ := newReadTestService(t, "vocab-xform-group")
	reqCtx := request.NewRequestContext(context.Background(), nil, "000000000000", "us-east-1")
	ctx := context.Background()
	req := func(params map[string]interface{}) *request.ParsedRequest {
		return &request.ParsedRequest{Operation: "Vocab", Parameters: params}
	}

	putResp, err := svc.PutTransformer(ctx, reqCtx, req(map[string]interface{}{
		"logGroupIdentifier": "vocab-xform-group",
		"transformerConfig": []interface{}{
			map[string]interface{}{"parseJSON": map[string]interface{}{}},
		},
	}))
	if err != nil {
		t.Fatalf("put transformer: %v", err)
	}
	if len(putResp.(map[string]interface{})) != 0 {
		t.Fatalf("PutTransformer response is not empty: %v", putResp)
	}

	getResp, err := svc.GetTransformer(ctx, reqCtx, req(map[string]interface{}{
		"logGroupIdentifier": "vocab-xform-group",
	}))
	if err != nil {
		t.Fatalf("get transformer: %v", err)
	}
	get := getResp.(map[string]interface{})
	for k := range get {
		switch k {
		case "logGroupIdentifier", "creationTime", "lastModifiedTime", "transformerConfig":
		default:
			t.Fatalf("GetTransformer carries the unmodelled member %q", k)
		}
	}

	testResp, err := svc.TestTransformer(ctx, reqCtx, req(map[string]interface{}{
		"transformerConfig": []interface{}{
			map[string]interface{}{"parseJSON": map[string]interface{}{}},
		},
		"logEventMessages": []interface{}{`{"a": 1}`},
	}))
	if err != nil {
		t.Fatalf("test transformer: %v", err)
	}
	if _, ok := testResp.(map[string]interface{})["transformedLogs"]; !ok {
		t.Fatalf("TestTransformer response: %v", testResp)
	}

	delResp, err := svc.DeleteTransformer(ctx, reqCtx, req(map[string]interface{}{
		"logGroupIdentifier": "vocab-xform-group",
	}))
	if err != nil {
		t.Fatalf("delete transformer: %v", err)
	}
	if len(delResp.(map[string]interface{})) != 0 {
		t.Fatalf("DeleteTransformer response is not empty: %v", delResp)
	}
}

// A dateTimeConverter without matchPatterns rejects with the processor's
// validation error — the absent member must not panic the validator's
// type assertion.
func TestDateTimeConverterWithoutMatchPatternsRejected(t *testing.T) {
	err := validateTransformerConfig([]map[string]interface{}{
		{"parseJSON": map[string]interface{}{}},
		{"dateTimeConverter": map[string]interface{}{
			"source": "t", "target": "u", "locale": "en",
		}},
	})
	if err == nil {
		t.Fatal("dateTimeConverter without matchPatterns must reject")
	}
	if !strings.Contains(err.Error(), "matchPatterns") {
		t.Fatalf("error must name the missing member: %v", err)
	}
}

// An empty selectionCriteria is the catch-all and a one-character
// LogGroupNamePrefix is a legal prefix: both must be selectable — a
// length-1 sentinel in the best-prefix comparison barred them both.
func TestTransformerPolicyEmptyAndOneCharPrefixApply(t *testing.T) {
	svc, store := newReadTestService(t, "catchall-group")

	policyDoc := `[{"parseJSON":{}},{"addKeys":{"entries":[{"key":"source","value":"account"}]}}]`
	if _, err := svc.putAccountPolicyCore("acct-catchall", policyDoc, "TRANSFORMER_POLICY", "",
		"", "us-east-1"); err != nil {
		t.Fatalf("put catch-all policy: %v", err)
	}
	if recipe := svc.effectiveTransformer(store, "us-east-1", "catchall-group"); recipe.config == nil {
		t.Fatal("the empty selectionCriteria (catch-all) never applied")
	}

	// The reverse coexistence order rejects too: with the account-wide
	// policy standing, a prefix-scoped policy cannot join it.
	if _, err := svc.putAccountPolicyCore("acct-scoped", policyDoc, "TRANSFORMER_POLICY", "",
		"LogGroupNamePrefix = /aws", "us-east-1"); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("prefix policy over the account-wide policy: %v", err)
	}

	// A one-character LogGroupNamePrefix is a legal prefix on its own
	// store (no overlap with the catch-all above).
	svc2, store2 := newReadTestService(t, "xgroup")
	if _, err := svc2.putAccountPolicyCore("acct-onechar", policyDoc, "TRANSFORMER_POLICY", "",
		"LogGroupNamePrefix = x", "us-east-1"); err != nil {
		t.Fatalf("put one-char-prefix policy: %v", err)
	}
	if recipe := svc2.effectiveTransformer(store2, "us-east-1", "xgroup"); recipe.config == nil {
		t.Fatal("a one-character LogGroupNamePrefix never applied")
	}
	if recipe := svc2.effectiveTransformer(store2, "us-east-1", "other"); recipe.config != nil {
		t.Fatalf("one-char prefix applied outside its scope: %v", recipe.config)
	}
}

// A deleted group's memoised transformer recipe dies with the group: a
// same-named group recreated right after the deletion starts with no
// recipe, not with the deleted one until the cache TTL lapses.
func TestTransformerCacheInvalidatedOnGroupDeletion(t *testing.T) {
	svc, store := newReadTestService(t, "recycled-group")
	recipe := []map[string]interface{}{
		{"parseJSON": map[string]interface{}{}},
	}
	if _, err := svc.putTransformerCore("recycled-group", recipe, "us-east-1"); err != nil {
		t.Fatal(err)
	}
	if got := svc.effectiveTransformer(store, "us-east-1", "recycled-group"); got.config == nil {
		t.Fatal("transformer did not apply before the deletion")
	}

	if err := svc.deleteLogGroupCore(DeleteLogGroupInput{LogGroupName: "recycled-group", Region: "us-east-1"}); err != nil {
		t.Fatal(err)
	}
	if err := store.CreateLogGroup(logsstore.NewLogGroup("recycled-group", "us-east-1", "000000000000")); err != nil {
		t.Fatal(err)
	}
	if got := svc.effectiveTransformer(store, "us-east-1", "recycled-group"); got.config != nil {
		t.Fatalf("recreated group inherited the deleted recipe: %v", got.config)
	}
}

// The compiled grok cache reuses one compiled form per match string —
// the ingestion hot path asks per event — and caches a failed compile
// as a stable miss.
func TestGrokCompileCacheReusesCompiledForm(t *testing.T) {
	first, ok := compiledGrokFor("%{IP:client.ip}")
	if !ok {
		t.Fatal("dotted match failed to compile")
	}
	second, ok := compiledGrokFor("%{IP:client.ip}")
	if !ok || first != second {
		t.Fatalf("compiled grok not reused: %p %p", first, second)
	}
	if _, ok := compiledGrokFor("%{IP:my-field}"); ok {
		t.Fatal("non-compiling match reported ok")
	}
	if _, ok := compiledGrokFor("%{IP:my-field}"); ok {
		t.Fatal("cached compile failure flip-flopped")
	}
}

// grokMatches reports whether a single-pattern match expression matches
// the whole text.
func grokMatches(match, text string) bool {
	compiled, ok := compileGrok(match)
	if !ok {
		return false
	}
	return compiled.re.MatchString(text)
}

// The pattern-library rows against their documented examples: UNIXPATH
// carries query parameters, MonthName admits full names, the seconds
// families include the leap second, DATE/DATESTAMP accept both date
// orders, the ISO-8601 offset hours are one-or-two digits, and the ARN
// row refuses ARNs missing information between colons.
func TestGrokPatternLibraryDocumentedRows(t *testing.T) {
	rows := []struct {
		name  string
		match string
		text  string
		want  bool
	}{
		{"UNIXPATH query parameter", `%{UNIXPATH:p}`, "/search?q=regex", true},
		{"UNIXPATH plain path", `%{UNIXPATH:p}`, "/category/sub", true},
		{"MONTH abbreviated", `%{MONTH:m}`, "Nov", true},
		{"MONTH full name", `%{MONTH:m}`, "November", true},
		{"HTTPDATE full month", `%{HTTPDATE:d}`, "23/November/2024:14:30:00 +0640", true},
		{"SYSLOGTIMESTAMP full month", `%{SYSLOGTIMESTAMP:d}`, "November 29 14:30:00", true},
		{"TIME leap second", `%{TIME:t}`, "09:45:60", true},
		{"ISO8601_SECOND leap second", `%{ISO8601_SECOND:s}`, "60", true},
		{"DATESTAMP documented EU row", `%{DATESTAMP:d}`, "29-11-2024 14:30:00", true},
		{"DATESTAMP US order", `%{DATESTAMP:d}`, "11/29/2024 14:30:00", true},
		{"ISO8601_TIMEZONE -530", `%{ISO8601_TIMEZONE:tz}`, "-530", true},
		{"ISO8601_TIMEZONE +5:30", `%{ISO8601_TIMEZONE:tz}`, "+5:30", true},
		{"TIMESTAMP_ISO8601 one-digit offset hour", `%{TIMESTAMP_ISO8601:ts}`, "23-5-1T1:25+5:30", true},
		{"TIMESTAMP_ISO8601 documented form", `%{TIMESTAMP_ISO8601:ts}`, "2023-05-15T14:30:00+05:30", true},
		{"ARN documented example", `%{ARN:a}`, "arn:aws:iam:us-east-1:123456789012:user/johndoe", true},
		{"ARN empty segments refuse", `%{ARN:a}`, "arn:aws:iam::123456789012:role/x", false},
		{"ARN empty region refuse", `%{ARN:a}`, "arn:aws:s3:::bucket", false},
	}
	for _, row := range rows {
		if got := grokMatches(row.match, row.text); got != row.want {
			t.Errorf("%s: match=%v want=%v", row.name, got, row.want)
		}
	}

	// A plain %{HTTPDATE:timestamp} capture stays raw — only the
	// common-log composites convert (documented grok example).
	got := transformOf(t, `[{"grok":{"match":"%{HTTPDATE:timestamp} %{IPORHOST:clientip} %{NUMBER:response_status}"}}]`,
		"23/Nov/2024:10:25:15 -0900 172.16.0.1 200")
	if got["timestamp"] != "23/Nov/2024:10:25:15 -0900" {
		t.Fatalf("plain HTTPDATE capture must stay raw: %v", got)
	}

	// The URIHOST alias's port carries its colon: the documented worked
	// example matches example.com:443 as the whole host instead of
	// skipping forward onto the bare port digits. The port is captured
	// from the first, port-bearing use and the second, portless use
	// contributes nothing — the documented output carries one port key.
	got = transformOf(t, `[{"grok":{"match":"%{URIHOST:host} %{URIHOST:ip}"}}]`,
		"example.com:443 10.0.0.1")
	if got["host"] != "example.com:443" || got["ip"] != "10.0.0.1" || got["port"] != "443" {
		t.Fatalf("URIHOST documented example: %v", got)
	}
	// An optional capture that did not participate is omitted, never
	// stored as an empty string.
	got = transformOf(t, `[{"grok":{"match":"%{URIHOST:host}"}}]`, "example.com")
	if got["host"] != "example.com" {
		t.Fatalf("portless URIHOST host: %v", got)
	}
	if _, exists := got["port"]; exists {
		t.Fatalf("portless URIHOST emitted a port: %v", got)
	}
}

// The grok match expression's structural rules: the twenty-pattern
// ceiling, the five-times budgets, GREEDYDATA_MULTILINE's single use,
// and the common-log composites' first-and-alone rule with their
// documented valid and invalid examples.
func TestValidateGrokMatchStructuralRules(t *testing.T) {
	valid := []string{
		"%{NGINX_ACCESS_LOG} %{DATA}",
		"%{SYSLOG5424}%{DATA:logMsg}",
		"%{APACHE_ACCESS_LOG} %{GREEDYDATA:logMsg}",
		"%{NUMBER:timestamp} [%{NUMBER:db} %{IP:client_ip}:%{NUMBER:client_port}] %{GREEDYDATA:data}",
		// The dotted form is the documented JSON-path notation.
		"%{IP:client.ip} %{NUMBER:client.port}",
		strings.Repeat("%{NUMBER:n} ", 19) + "%{NUMBER:n}",
		// Five combined uses of the shared-budget patterns are legal.
		strings.Repeat("%{DATA:d} ", 4) + "%{GREEDYDATA:g}",
	}
	for _, match := range valid {
		if err := validateGrokMatch(match); err != nil {
			t.Errorf("valid match %q rejected: %v", match, err)
		}
	}
	invalid := []struct{ match, because string }{
		{strings.Repeat("%{NUMBER:n} ", 20) + "%{NUMBER:n}", "twenty-one patterns"},
		{strings.Repeat("%{DATA:d} ", 5) + "%{DATA:d}", "six DATA patterns"},
		{strings.Repeat("%{DATA:d} ", 5) + "%{GREEDYDATA:g}", "six combined DATA and GREEDYDATA"},
		{strings.Repeat("%{GREEDYDATA_MULTILINE:m} ", 1) + "%{GREEDYDATA_MULTILINE:m}", "two GREEDYDATA_MULTILINE"},
		{strings.Repeat("%{NOTSPACE:n} ", 5) + "%{NOTSPACE:n}", "six NOTSPACE patterns"},
		{"%{GREEDYDATA:logMsg} %{SYSLOG5424}", "pattern precedes the common log pattern"},
		{"%{APACHE_ACCESS_LOG} %{SYSLOG5424}", "multiple common log patterns"},
		{"%{NGINX_ACCESS_LOG} %{NUMBER:num}", "non-DATA follower"},
		{"%{NGINX_ACCESS_LOG} %{DATA} %{DATA}", "two followers"},
		{"%{NOT_A_PATTERN:x}", "unknown pattern"},
		// A field name outside the word-character vocabulary cannot
		// become a named capture group: Put rejects it instead of
		// accepting a processor whose compilation silently fails.
		{"%{IP:my-field}", "dash in the field name"},
		{"%{IP:client..ip}", "empty path segment"},
		{"%{IP:.client}", "leading path dot"},
		{"%{IP:名前}", "non-ASCII field name"},
	}
	for _, row := range invalid {
		if err := validateGrokMatch(row.match); err == nil {
			t.Errorf("invalid match (%s) accepted: %q", row.because, row.match)
		}
	}
}

// configFromJSON parses a processor-array configuration for the
// validation rows.
func configFromJSON(t *testing.T, config string) []map[string]interface{} {
	t.Helper()
	var parsed []map[string]interface{}
	if err := json.Unmarshal([]byte(config), &parsed); err != nil {
		t.Fatalf("config does not parse: %v", err)
	}
	return parsed
}

// The configuration's per-member rows: the model's entry ceilings
// (deleteKeys and typeConverter at five, the case/trim siblings' ten
// preserved), the path length and nested-depth bounds, the csv and
// dateTimeConverter member rows, the union's single-member Processor
// elements, and the vended parsers' first-and-alone rule.
func TestTransformerMemberConstraintRows(t *testing.T) {
	reject := []string{
		`[{"parseJSON":{}},{"deleteKeys":{"withKeys":["a","b","c","d","e","f"]}}]`,
		`[{"parseJSON":{}},{"typeConverter":{"entries":[` +
			`{"key":"a","type":"string"},{"key":"b","type":"string"},{"key":"c","type":"string"},` +
			`{"key":"d","type":"string"},{"key":"e","type":"string"},{"key":"f","type":"string"}]}}]`,
		`[{"parseJSON":{}},{"lowerCaseString":{"withKeys":[` +
			`"a","b","c","d","e","f","g","h","i","j","k"]}}]`,
		`[{"grok":{"match":"%{NUMBER:n}"},"addKeys":{"entries":[]}}]`,
		`[{"parseWAF":{}},{"parseVPC":{}}]`,
		`[{"parseJSON":{}},{"parseWAF":{}}]`,
		`[{"parseJSON":{}},{"addKeys":{"entries":[{"key":"a.b.c.d","value":"v"}]}}]`,
		`[{"parseJSON":{}},{"addKeys":{"entries":[{"key":"k","value":"` +
			strings.Repeat("v", 257) + `"}]}}]`,
		`[{"parseJSON":{"source":"a.b.c.d"}}]`,
		`[{"csv":{"delimiter":"ab"}}]`,
		`[{"csv":{"quoteCharacter":"ab"}}]`,
		`[{"csv":{"columns":[` + func() string {
			names := make([]string, 101)
			for i := range names {
				names[i] = fmt.Sprintf("c%03d", i)
			}
			return `"` + strings.Join(names, `","`) + `"`
		}() + `]}}]`,
		`[{"parseJSON":{}},{"dateTimeConverter":{"source":"t","target":"u","locale":"en","matchPatterns":[` +
			`"yyyy","yyyy","yyyy","yyyy","yyyy","yyyy"]}}]`,
		`[{"parseJSON":{}},{"dateTimeConverter":{"source":"t","target":"u","locale":"en",` +
			`"matchPatterns":["yyyy"],"targetFormat":"` + strings.Repeat("y", 65) + `"}}]`,
		`[{"parseJSON":{}},{"typeConverter":{"entries":[{"key":"k","type":"float"}]}}]`,
		`[{"parseKeyValue":{"fieldDelimiter":""}}]`,
		`[{"parseJSON":{}},{"listToMap":{"source":"s","key":"k","flatten":true}}]`,
		`[{"parseJSON":{}},{"listToMap":{"source":"s","key":"k","flatten":true,"flattenedElement":"middle"}}]`,
	}
	for _, config := range reject {
		if err := validateTransformerConfig(configFromJSON(t, config)); err == nil {
			t.Errorf("invalid configuration accepted: %s", config)
		}
	}
	accept := []string{
		`[{"parseJSON":{}},{"deleteKeys":{"withKeys":["a","b","c","d","e"]}}]`,
		`[{"parseJSON":{}},{"typeConverter":{"entries":[` +
			`{"key":"a","type":"string"},{"key":"b","type":"string"},{"key":"c","type":"string"},` +
			`{"key":"d","type":"string"},{"key":"e","type":"string"}]}}]`,
		`[{"parseJSON":{}},{"lowerCaseString":{"withKeys":[` +
			`"a","b","c","d","e","f","g","h","i","j"]}}]`,
		`[{"csv":{"delimiter":"\\t"}}]`,
		`[{"parseJSON":{}},{"addKeys":{"entries":[{"key":"a.b.c","value":"` +
			strings.Repeat("v", 256) + `"}]}}]`,
		`[{"parseJSON":{}},{"listToMap":{"source":"s","key":"k","flatten":true,"flattenedElement":"first"}}]`,
	}
	for _, config := range accept {
		if err := validateTransformerConfig(configFromJSON(t, config)); err != nil {
			t.Errorf("valid configuration rejected: %s: %v", config, err)
		}
	}
}

// The TRANSFORMER_POLICY document contract at Put: "A transformer
// policy must include one JSON block with the array of processors and
// their configurations" — the bare array form alone parses, and the
// recipe it carries must be a valid transformer.
func TestTransformerPolicyDocumentValidation(t *testing.T) {
	svc, _ := newReadTestService(t, "poldoc-group")
	valid := `[{"parseJSON":{}}]`

	for _, row := range []struct{ name, document string }{
		{"wrapped object form", `{"transformerConfig":[{"parseJSON":{}}]}`},
		{"not JSON", `not-json`},
		{"empty array", `[]`},
		{"invalid recipe", `[{"parseJSON":{}},{"deleteKeys":{"withKeys":["a","b","c","d","e","f"]}}]`},
	} {
		if _, err := svc.putAccountPolicyCore("pol-"+row.name, row.document, "TRANSFORMER_POLICY",
			"", "LogGroupNamePrefix = /pol", "us-east-1"); err == nil {
			t.Errorf("%s document accepted", row.name)
		}
	}
	if _, err := svc.putAccountPolicyCore("pol-valid", valid, "TRANSFORMER_POLICY",
		"", "LogGroupNamePrefix = /pol", "us-east-1"); err != nil {
		t.Fatalf("valid bare-array document rejected: %v", err)
	}
}

// "You can create as many as 20 account-level transformer policies ...
// in a Region" — the 21st distinct policy rejects with
// LimitExceededException and overwriting an existing name stays legal.
func TestTransformerPolicyTwentyCap(t *testing.T) {
	svc, _ := newReadTestService(t, "cap-group")
	doc := `[{"parseJSON":{}}]`
	for i := 0; i < 20; i++ {
		name := fmt.Sprintf("tp%02d", i)
		prefix := fmt.Sprintf("cap%02d", i)
		if _, err := svc.putAccountPolicyCore(name, doc, "TRANSFORMER_POLICY",
			"", "LogGroupNamePrefix = "+prefix, "us-east-1"); err != nil {
			t.Fatalf("policy %s: %v", name, err)
		}
	}
	_, err := svc.putAccountPolicyCore("tp21", doc, "TRANSFORMER_POLICY",
		"", "LogGroupNamePrefix = capxyz", "us-east-1")
	if err == nil || !strings.Contains(err.Error(), "LimitExceeded") {
		t.Fatalf("21st policy: %v", err)
	}
	if _, err := svc.putAccountPolicyCore("tp00", doc, "TRANSFORMER_POLICY",
		"", "LogGroupNamePrefix = cap00", "us-east-1"); err != nil {
		t.Fatalf("overwrite at the cap: %v", err)
	}
}

// "Log transformation and enrichment is supported only for log groups
// in the Standard log class" — an account-level policy matches no
// non-Standard group, mirroring the group-level form's Put-time check.
func TestTransformerPolicySkipsNonStandardGroups(t *testing.T) {
	svc, store := newReadTestService(t, "std-policy-group")
	ia := logsstore.NewLogGroup("ia-policy-group", "us-east-1", "000000000000")
	ia.LogGroupClass = "INFREQUENT_ACCESS"
	if err := store.CreateLogGroup(ia); err != nil {
		t.Fatal(err)
	}
	doc := `[{"parseJSON":{}}]`
	if _, err := svc.putAccountPolicyCore("pol-class", doc, "TRANSFORMER_POLICY",
		"", "", "us-east-1"); err != nil {
		t.Fatalf("put catch-all policy: %v", err)
	}
	if recipe := svc.effectiveTransformer(store, "us-east-1", "ia-policy-group"); recipe.config != nil {
		t.Fatalf("account policy transformed a non-Standard group: %v", recipe.config)
	}
	if recipe := svc.effectiveTransformer(store, "us-east-1", "std-policy-group"); recipe.config == nil {
		t.Fatal("account policy skipped a Standard group")
	}
}

// The documented transformation failure surface: "Whenever CloudWatch
// Logs tries and fails to transform a log event, it adds a
// @transformationError system field to that log event" — the 512KB
// bound and the 200-field extraction ceiling persist the marker record,
// and the batch's TransformedLogEvents / TransformedBytes /
// TransformationErrors metrics publish with the level's dimension.
func TestTransformationErrorMarkerAndMetrics(t *testing.T) {
	parseJSONRecipe := `[{"parseJSON":{}}]`

	t.Run("oversize event carries the marker", func(t *testing.T) {
		out, failed := transformEvent(configFromJSON(t, parseJSONRecipe),
			`{"a":"`+strings.Repeat("x", maxTransformMessageBytes)+`"}`, TransformContext{})
		if !failed || out == "" {
			t.Fatalf("oversize event must fail with a marker record: %q failed=%v", out, failed)
		}
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(out), &record); err != nil {
			t.Fatal(err)
		}
		if _, ok := record["@transformationError"]; !ok {
			t.Fatalf("marker record missing the system field: %v", record)
		}
	})

	t.Run("field ceiling carries the marker", func(t *testing.T) {
		var b strings.Builder
		b.WriteString("{")
		for i := 0; i <= maxTransformedFieldCount; i++ {
			if i > 0 {
				b.WriteString(",")
			}
			fmt.Fprintf(&b, `"k%03d":%d`, i, i)
		}
		b.WriteString("}")
		out, failed := transformEvent(configFromJSON(t, parseJSONRecipe), b.String(), TransformContext{})
		if !failed {
			t.Fatalf("201 fields must fail: %q", out)
		}
		var record map[string]interface{}
		if err := json.Unmarshal([]byte(out), &record); err != nil {
			t.Fatal(err)
		}
		if _, ok := record["@transformationError"]; !ok {
			t.Fatalf("marker record missing the system field: %v", record)
		}
	})

	t.Run("marker surfaces as the single-@ system field", func(t *testing.T) {
		marker := `{"@transformationError":"reason"}`
		transformedRow := queryResultRow{}
		discoverJSONFields(marker, &transformedRow, true)
		if _, ok := transformedRow.fields["@transformationError"]; !ok {
			t.Fatalf("transformed copy must surface the system field: %v", transformedRow.fields)
		}
		originalRow := queryResultRow{}
		discoverJSONFields(marker, &originalRow, false)
		if _, ok := originalRow.fields["@@transformationError"]; !ok {
			t.Fatalf("an original message's @-key must render doubled: %v", originalRow.fields)
		}
	})

	t.Run("metrics publish with the level's dimension", func(t *testing.T) {
		svc, store := newReadTestService(t, "metric-group")
		metrics := &valueRecordingMetricInvoker{}
		svc.SetCloudWatchMetricInvoker(metrics)
		if _, err := svc.putTransformerCore("metric-group", configFromJSON(t, parseJSONRecipe), "us-east-1"); err != nil {
			t.Fatal(err)
		}
		svc.transformIngestedBatch(store, "us-east-1", "metric-group", "s1", []logsstore.LogEntry{
			{Timestamp: 1000, Message: `{"a":1}`},
			{Timestamp: 2000, Message: `{"b":"` + strings.Repeat("x", maxTransformMessageBytes) + `"}`},
		})
		var events, bytes, failures float64
		metricCalls, metricDims := metrics.recorded()
		for i, call := range metricCalls {
			if call.namespace != "AWS/Logs" {
				t.Fatalf("namespace %q", call.namespace)
			}
			if metricDims[i]["LogGroupname"] != "metric-group" {
				t.Fatalf("dimension %v", metricDims[i])
			}
			switch call.metric {
			case "TransformedLogEvents":
				events = call.value
			case "TransformedBytes":
				bytes = call.value
			case "TransformationErrors":
				failures = call.value
			}
		}
		if events != 1 || failures != 1 || bytes != float64(len(`{"a":1}`)) {
			t.Fatalf("events=%v bytes=%v failures=%v", events, bytes, failures)
		}
	})

	t.Run("account-level metrics carry PolicyLevel", func(t *testing.T) {
		svc, store := newReadTestService(t, "acct-metric-group")
		metrics := &valueRecordingMetricInvoker{}
		svc.SetCloudWatchMetricInvoker(metrics)
		if _, err := svc.putAccountPolicyCore("acct-metrics", parseJSONRecipe, "TRANSFORMER_POLICY",
			"", "LogGroupNamePrefix = acct-metric", "us-east-1"); err != nil {
			t.Fatal(err)
		}
		svc.transformIngestedBatch(store, "us-east-1", "acct-metric-group", "s1", []logsstore.LogEntry{
			{Timestamp: 1000, Message: `{"a":1}`},
		})
		acctCalls, acctDims := metrics.recorded()
		if len(acctCalls) == 0 {
			t.Fatal("no transformation metrics emitted")
		}
		for i := range acctCalls {
			if acctDims[i]["PolicyLevel"] != "AccountPolicy" {
				t.Fatalf("dimension %v", acctDims[i])
			}
			if acctDims[i]["LogGroupname"] != "" {
				t.Fatalf("account-level metric carries the group dimension: %v", acctDims[i])
			}
		}
	})
}

// The dateTimeConverter's z letter renders the zone NAME — the
// documented German example targetFormat "yyyy-MM-dd'T'HH:mm:ss z"
// produces "MEZ" for Europe/Berlin.
func TestDateTimeConverterZoneLetter(t *testing.T) {
	rows := []struct {
		source, want string
	}{
		{"05. Dezember 1998 11:00:00", "1998-12-05T11:00:00 MEZ"},
		{"05. Juni 1998 11:00:00", "1998-06-05T11:00:00 MESZ"},
		// The locale tables' non-ASCII month names are whole words too:
		// März, février and août localise and parse like their ASCII
		// siblings.
		{"01. März 2024 09:30:00", "2024-03-01T09:30:00 MEZ"},
		{"01 février 2024 08:15:00", "2024-02-01T08:15:00 CET"},
		{"15 août 2023 14:00:00", "2023-08-15T14:00:00 CEST"},
	}
	for _, row := range rows {
		recipe := `[{"parseJSON":{}},{"dateTimeConverter":{"source":"german_datetime","target":"target_1",` +
			`"locale":"de","matchPatterns":["dd. MMMM yyyy HH:mm:ss"],` +
			`"sourceTimezone":"Europe/Berlin","targetTimezone":"Europe/Berlin",` +
			`"targetFormat":"yyyy-MM-dd'T'HH:mm:ss z"}}]`
		if strings.Contains(row.source, "février") || strings.Contains(row.source, "août") {
			recipe = `[{"parseJSON":{}},{"dateTimeConverter":{"source":"french_datetime","target":"target_1",` +
				`"locale":"fr","matchPatterns":["dd MMMM yyyy HH:mm:ss"],` +
				`"sourceTimezone":"Europe/Paris","targetTimezone":"Europe/Berlin",` +
				`"targetFormat":"yyyy-MM-dd'T'HH:mm:ss z"}}]`
		}
		got := transformOf(t, recipe, `{"german_datetime": "`+row.source+`", "french_datetime": "`+row.source+`"}`)
		if got["target_1"] != row.want {
			t.Errorf("%s: got %v want %s", row.source, got["target_1"], row.want)
		}
	}
}

// The transformerConfig member is a list of Processor structures: an
// element that is not an object is a malformed request and rejects, never
// a silently dropped entry — the stored pipeline must be the pipeline the
// caller sent.
func TestTransformerConfigNonObjectElementsReject(t *testing.T) {
	bad := map[string]interface{}{
		"transformerConfig": []interface{}{
			"garbage",
			map[string]interface{}{"parseJSON": map[string]interface{}{}},
		},
	}
	if _, err := transformerConfigFromParams(bad); logsErrorCode(err) != "InvalidParameterException" {
		t.Fatalf("non-object element: %v", err)
	}
	good := map[string]interface{}{
		"transformerConfig": []interface{}{
			map[string]interface{}{"parseJSON": map[string]interface{}{}},
		},
	}
	config, err := transformerConfigFromParams(good)
	if err != nil {
		t.Fatalf("object elements: %v", err)
	}
	if len(config) != 1 {
		t.Fatalf("object elements: want 1 processor, got %d", len(config))
	}
}
