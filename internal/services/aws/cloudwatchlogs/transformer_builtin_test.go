package cloudwatchlogs

import (
	"encoding/json"
	"reflect"
	"testing"
)

// The built-in vended-format parser pins: every documented example on
// the Built-in processors page transforms to exactly the documented
// output, and the structural rules (first processor, @message-only
// input, the parseToOCSF rejection) hold.

// transformOfJSON runs one processor list over a message and returns the
// output as the JSON-decoded object the pin compares against the
// documented example.
func transformOfJSON(t *testing.T, config []map[string]interface{}, message string) map[string]interface{} {
	t.Helper()
	out, _ := transformEvent(config, message, TransformContext{})
	if out == "" {
		t.Fatalf("no output for %q", message)
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("output %q is not a JSON object: %v", out, err)
	}
	return parsed
}

func TestTransformEngineBuiltinParsersDocumentedExamples(t *testing.T) {
	cases := []struct {
		name    string
		config  []map[string]interface{}
		message string
		want    string
	}{
		{
			name:    "parseWAF documented example",
			config:  []map[string]interface{}{{"parseWAF": map[string]interface{}{}}},
			message: `{"timestamp":1576280412771,"formatVersion":1,"webaclId":"arn:aws:wafv2:ap-southeast-2:111122223333:regional/webacl/STMTest/1EXAMPLE-2ARN-3ARN-4ARN-123456EXAMPLE","terminatingRuleId":"STMTest_SQLi_XSS","terminatingRuleType":"REGULAR","action":"BLOCK","terminatingRuleMatchDetails":[{"conditionType":"SQL_INJECTION","sensitivityLevel":"HIGH","location":"HEADER","matchedData":["10","AND","1"]}],"httpSourceName":"-","httpSourceId":"-","ruleGroupList":[],"rateBasedRuleList":[],"nonTerminatingMatchingRules":[],"httpRequest":{"clientIp":"1.1.1.1","country":"AU","headers":[{"name":"Host","value":"localhost:1989"},{"name":"User-Agent","value":"curl/7.61.1"},{"name":"Accept","value":"*/*"},{"name":"x-stm-test","value":"10 AND 1=1"}],"uri":"/myUri","args":"","httpVersion":"HTTP/1.1","httpMethod":"GET","requestId":"rid"},"labels":[{"name":"value"}]}`,
			want:    `{"httpRequest":{"headers":{"Host":"localhost:1989","User-Agent":"curl/7.61.1","Accept":"*/*","x-stm-test":"10 AND 1=1"},"clientIp":"1.1.1.1","country":"AU","uri":"/myUri","args":"","httpVersion":"HTTP/1.1","httpMethod":"GET","requestId":"rid"},"labels":{"name":"value"},"timestamp":1576280412771,"formatVersion":1,"webaclId":"arn:aws:wafv2:ap-southeast-2:111122223333:regional/webacl/STMTest/1EXAMPLE-2ARN-3ARN-4ARN-123456EXAMPLE","terminatingRuleId":"STMTest_SQLi_XSS","terminatingRuleType":"REGULAR","action":"BLOCK","terminatingRuleMatchDetails":[{"conditionType":"SQL_INJECTION","sensitivityLevel":"HIGH","location":"HEADER","matchedData":["10","AND","1"]}],"httpSourceName":"-","httpSourceId":"-","ruleGroupList":[],"rateBasedRuleList":[],"nonTerminatingMatchingRules":[]}`,
		},
		{
			name:    "parsePostgres documented example",
			config:  []map[string]interface{}{{"parsePostgres": map[string]interface{}{}}},
			message: `2019-03-10 03:54:59 UTC:10.0.0.123(52834):postgres@logtestdb:[20175]:ERROR: column "wrong_column_name" does not exist at character 8`,
			want:    `{"logTime":"2019-03-10 03:54:59 UTC","srcIp":"10.0.0.123(52834)","userName":"postgres","dbName":"logtestdb","processId":"20175","logLevel":"ERROR"}`,
		},
		{
			name:    "parseCloudfront documented example",
			config:  []map[string]interface{}{{"parseCloudfront": map[string]interface{}{}}},
			message: "2019-12-04  21:02:31   LAX1   392    192.0.2.24    GET    d111111abcdef8.cloudfront.net  /index.html    200    -  Mozilla/5.0%20(Windows%20NT%2010.0;%20Win64;%20x64)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/78.0.3904.108%20Safari/537.36  -  -  Hit    SOX4xwn4XV6Q4rgb7XiVGOHms_BGlTAC4KyHmureZmBNrjGdRLiNIQ==   d111111abcdef8.cloudfront.net  https  23 0.001  -  TLSv1.2    ECDHE-RSA-AES128-GCM-SHA256    Hit    HTTP/2.0   -  -  11040  0.001  Hit    text/html  78 -  -",
			want:    `{"date":"2019-12-04","time":"21:02:31","x-edge-location":"LAX1","sc-bytes":392,"c-ip":"192.0.2.24","cs-method":"GET","cs(Host)":"d111111abcdef8.cloudfront.net","cs-uri-stem":"/index.html","sc-status":200,"cs(Referer)":"-","cs(User-Agent)":"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36","cs-uri-query":"-","cs(Cookie)":"-","x-edge-result-type":"Hit","x-edge-request-id":"SOX4xwn4XV6Q4rgb7XiVGOHms_BGlTAC4KyHmureZmBNrjGdRLiNIQ==","x-host-header":"d111111abcdef8.cloudfront.net","cs-protocol":"https","cs-bytes":23,"time-taken":0.001,"x-forwarded-for":"-","ssl-protocol":"TLSv1.2","ssl-cipher":"ECDHE-RSA-AES128-GCM-SHA256","x-edge-response-result-type":"Hit","cs-protocol-version":"HTTP/2.0","fle-status":"-","fle-encrypted-fields":"-","c-port":11040,"time-to-first-byte":0.001,"x-edge-detailed-result-type":"Hit","sc-content-type":"text/html","sc-content-len":78,"sc-range-start":"-","sc-range-end":"-"}`,
		},
		{
			name:    "parseRoute53 documented example",
			config:  []map[string]interface{}{{"parseRoute53": map[string]interface{}{}}},
			message: `1.0 2017-12-13T08:15:50.235Z Z123412341234 example.com AAAA NOERROR TCP IAD12 192.0.2.0 198.51.100.0/24`,
			want:    `{"version":1.0,"queryTimestamp":"2017-12-13T08:15:50.235Z","hostZoneId":"Z123412341234","queryName":"example.com","queryType":"AAAA","responseCode":"NOERROR","protocol":"TCP","edgeLocation":"IAD12","resolverIp":"192.0.2.0","ednsClientSubnet":"198.51.100.0/24"}`,
		},
		{
			name:    "parseVPC documented example",
			config:  []map[string]interface{}{{"parseVPC": map[string]interface{}{}}},
			message: `2 123456789010 eni-abc123de 192.0.2.0 192.0.2.24 20641 22 6 20 4249 1418530010 1418530070 ACCEPT OK`,
			want:    `{"version":2,"accountId":"123456789010","interfaceId":"eni-abc123de","srcAddr":"192.0.2.0","dstAddr":"192.0.2.24","srcPort":20641,"dstPort":22,"protocol":6,"packets":20,"bytes":4249,"start":1418530010,"end":1418530070,"action":"ACCEPT","logStatus":"OK"}`,
		},
	}
	for _, tc := range cases {
		got := transformOfJSON(t, tc.config, tc.message)
		var want map[string]interface{}
		if err := json.Unmarshal([]byte(tc.want), &want); err != nil {
			t.Fatalf("%s: expected output is not JSON: %v", tc.name, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("%s:\n got %#v\nwant %#v", tc.name, got, want)
		}
	}

	// An unparsable message leaves the event untransformed (the original
	// still ingests).
	if out, _ := transformEvent([]map[string]interface{}{{"parseVPC": map[string]interface{}{}}},
		"not a flow log at all", TransformContext{}); out != "" {
		t.Fatalf("unparsable message transformed: %q", out)
	}
	if out, _ := transformEvent([]map[string]interface{}{{"parsePostgres": map[string]interface{}{}}},
		"not a postgres line", TransformContext{}); out != "" {
		t.Fatalf("unparsable postgres line transformed: %q", out)
	}
	if out, _ := transformEvent([]map[string]interface{}{{"parseWAF": map[string]interface{}{}}},
		"not json", TransformContext{}); out != "" {
		t.Fatalf("unparsable WAF log transformed: %q", out)
	}

	// The built-in parsers compose: a later mutator runs over the parsed
	// fields.
	composed := transformOfJSON(t, []map[string]interface{}{
		{"parseVPC": map[string]interface{}{}},
		{"upperCaseString": map[string]interface{}{"withKeys": []interface{}{"action"}}},
	}, `2 123456789010 eni-abc123de 192.0.2.0 192.0.2.24 20641 22 6 20 4249 1418530010 1418530070 accept OK`)
	if composed["action"] != "ACCEPT" {
		t.Fatalf("composed action: %v", composed["action"])
	}
}
