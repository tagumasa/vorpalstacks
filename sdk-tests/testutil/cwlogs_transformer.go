package testutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// The log transformer pins (E2): the TestTransformer dry run over the
// documented pipeline examples, the CRUD lifecycle with its validation
// rows, the ingestion application (originals on the API read plane,
// transformed fields through Insights), and the account-level
// TRANSFORMER_POLICY prefix selection.

func (tc *cwlogsTestCtx) transformerTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "TestTransformer_DocumentedPipeline", func() error {
		resp, err := tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{
				{ParseJSON: &types.ParseJSON{}},
				{Grok: &types.Grok{
					Source: aws.String("logMsg"),
					Match:  aws.String(`%{WORD:http_method} %{NOTSPACE:request} HTTP/%{NUMBER:http_version}`),
				}},
			},
			LogEventMessages: []string{
				`{"timestamp": "2024-11-23T16:03:12Z", "level": "ERROR", "logMsg": "GET /page.html HTTP/1.1"}`,
			},
		})
		if err != nil {
			return fmt.Errorf("test transformer: %v", err)
		}
		if len(resp.TransformedLogs) != 1 {
			return fmt.Errorf("transformed logs: %d", len(resp.TransformedLogs))
		}
		record := resp.TransformedLogs[0]
		if aws.ToString(record.EventMessage) == "" || record.EventNumber != 1 {
			return fmt.Errorf("record identity: %+v", record)
		}
		transformed := aws.ToString(record.TransformedEventMessage)
		for _, want := range []string{`"http_method":"GET"`, `"request":"/page.html"`, `"http_version":"1.1"`} {
			if !strings.Contains(transformed, want) {
				return fmt.Errorf("transformed log %q missing %s", transformed, want)
			}
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "TransformerCRUD_Lifecycle", func() error {
		group, cleanupGroup, err := tc.newLogGroupFixture("xform")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		groupArn, err := tc.findLogGroupARN(group)
		if err != nil {
			return err
		}

		config := []types.Processor{
			{ParseJSON: &types.ParseJSON{}},
			{AddKeys: &types.AddKeys{Entries: []types.AddKeyEntry{{
				Key: aws.String("enriched"), Value: aws.String("yes"),
			}}}},
		}
		if _, err := tc.client.PutTransformer(tc.ctx, &cloudwatchlogs.PutTransformerInput{
			LogGroupIdentifier: groupArn,
			TransformerConfig:  config,
		}); err != nil {
			return fmt.Errorf("put transformer: %v", err)
		}

		got, err := tc.client.GetTransformer(tc.ctx, &cloudwatchlogs.GetTransformerInput{
			LogGroupIdentifier: groupArn,
		})
		if err != nil {
			return fmt.Errorf("get transformer: %v", err)
		}
		if len(got.TransformerConfig) != 2 {
			return fmt.Errorf("config round trip: %d processors", len(got.TransformerConfig))
		}
		if got.CreationTime == nil || got.LastModifiedTime == nil {
			return fmt.Errorf("stamps missing: %v %v", got.CreationTime, got.LastModifiedTime)
		}
		if !strings.Contains(aws.ToString(got.LogGroupIdentifier), group) {
			return fmt.Errorf("identifier echo: %v", got.LogGroupIdentifier)
		}

		// The validation rows.
		var ipe *types.InvalidParameterException
		_, err = tc.client.PutTransformer(tc.ctx, &cloudwatchlogs.PutTransformerInput{
			LogGroupIdentifier: groupArn,
			TransformerConfig: []types.Processor{
				{AddKeys: &types.AddKeys{Entries: []types.AddKeyEntry{{Key: aws.String("k"), Value: aws.String("v")}}}},
			},
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("first-not-parser: %v", err)
		}
		_, err = tc.client.PutTransformer(tc.ctx, &cloudwatchlogs.PutTransformerInput{
			LogGroupIdentifier: groupArn,
			TransformerConfig: []types.Processor{{ParseToOCSF: &types.ParseToOCSF{
				EventSource: types.EventSourceVpcFlow, OcsfVersion: types.OCSFVersionV15,
			}}},
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("unpublished built-in parser: %v", err)
		}

		if _, err := tc.client.DeleteTransformer(tc.ctx, &cloudwatchlogs.DeleteTransformerInput{
			LogGroupIdentifier: groupArn,
		}); err != nil {
			return fmt.Errorf("delete transformer: %v", err)
		}
		var rnf *types.ResourceNotFoundException
		_, err = tc.client.GetTransformer(tc.ctx, &cloudwatchlogs.GetTransformerInput{
			LogGroupIdentifier: groupArn,
		})
		if !errors.As(err, &rnf) {
			return fmt.Errorf("deleted transformer: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "Transformer_GrokFieldNameValidation", func() error {
		var ipe *types.InvalidParameterException
		_, err := tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{
				{Grok: &types.Grok{Match: aws.String(`%{IP:my-field}`)}},
			},
			LogEventMessages: []string{"203.0.113.9"},
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("dash field name on the dry run: %v", err)
		}
		_, err = tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{
				{Grok: &types.Grok{Match: aws.String(`%{IP:client..ip}`)}},
			},
			LogEventMessages: []string{"203.0.113.9"},
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("empty path segment on the dry run: %v", err)
		}

		group, cleanupGroup, err := tc.newLogGroupFixture("grok-name")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		groupArn, err := tc.findLogGroupARN(group)
		if err != nil {
			return err
		}
		_, err = tc.client.PutTransformer(tc.ctx, &cloudwatchlogs.PutTransformerInput{
			LogGroupIdentifier: groupArn,
			TransformerConfig: []types.Processor{
				{Grok: &types.Grok{Match: aws.String(`%{IP:my-field}`)}},
			},
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("dash field name at Put: %v", err)
		}

		// The dotted form is the documented JSON-path notation: it must
		// survive the new check and extract a nested field.
		if _, err := tc.client.PutTransformer(tc.ctx, &cloudwatchlogs.PutTransformerInput{
			LogGroupIdentifier: groupArn,
			TransformerConfig: []types.Processor{
				{Grok: &types.Grok{Match: aws.String(`%{IP:client.ip}`)}},
			},
		}); err != nil {
			return fmt.Errorf("dotted field name at Put: %v", err)
		}
		resp, err := tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{
				{Grok: &types.Grok{Match: aws.String(`%{IP:client.ip}`)}},
			},
			LogEventMessages: []string{"203.0.113.9"},
		})
		if err != nil {
			return fmt.Errorf("dotted field name dry run: %v", err)
		}
		if len(resp.TransformedLogs) != 1 ||
			!strings.Contains(aws.ToString(resp.TransformedLogs[0].TransformedEventMessage), `"client":{"ip":"203.0.113.9"}`) {
			return fmt.Errorf("dotted extraction: %+v", resp.TransformedLogs)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "Transformer_AppliedAtIngestion", func() error {
		group, cleanupGroup, err := tc.newGroupStreamFixture("applied", "xform-stream")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		const stream = "xform-stream"
		groupArn, err := tc.findLogGroupARN(group)

		if _, err := tc.client.PutTransformer(tc.ctx, &cloudwatchlogs.PutTransformerInput{
			LogGroupIdentifier: groupArn,
			TransformerConfig: []types.Processor{
				{ParseJSON: &types.ParseJSON{}},
				{AddKeys: &types.AddKeys{Entries: []types.AddKeyEntry{{
					Key: aws.String("marker"), Value: aws.String("transformed"),
				}}}},
			},
		}); err != nil {
			return fmt.Errorf("put transformer: %v", err)
		}

		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(group, stream, `{"level":"ERROR","msg":"boom"}`, now); err != nil {
			return err
		}

		// The API read plane keeps the original message.
		events, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName: aws.String(group), LogStreamName: aws.String(stream),
		})
		if err != nil {
			return fmt.Errorf("get log events: %v", err)
		}
		if len(events.Events) != 1 || aws.ToString(events.Events[0].Message) != `{"level":"ERROR","msg":"boom"}` {
			return fmt.Errorf("original message replaced: %v", events.Events)
		}

		// The query plane reads the transformed form.
		startResp, err := tc.client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:           aws.Int64(now/1000 - 60),
			EndTime:             aws.Int64(now/1000 + 60),
			LogGroupIdentifiers: []string{aws.ToString(groupArn)},
			QueryString:         aws.String(`filter marker = "transformed" | fields @message`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		resResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		if len(resResp.Results) == 0 {
			return fmt.Errorf("query saw no transformed record")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "AccountPolicy_TransformerPrefix", func() error {
		prefix := "/" + tc.uniquePrefix("acct-xform")
		inner := prefix + "/inner"
		outer := "/" + tc.uniquePrefix("outside")

		for _, g := range []string{inner, outer} {
			if err := tc.createLogGroup(g); err != nil {
				return err
			}
			defer tc.deleteLogGroup(g)
			if err := tc.createLogStream(g, "s1"); err != nil {
				return err
			}
		}
		policyName := tc.uniquePrefix("xform-policy")
		doc := `[{"parseJSON":{}},{"addKeys":{"entries":[{"key":"scope","value":"account"}]}}]`
		if _, err := tc.client.PutAccountPolicy(tc.ctx, &cloudwatchlogs.PutAccountPolicyInput{
			PolicyName:        aws.String(policyName),
			PolicyType:        types.PolicyTypeTransformerPolicy,
			PolicyDocument:    aws.String(doc),
			SelectionCriteria: aws.String("LogGroupNamePrefix = " + prefix),
		}); err != nil {
			return fmt.Errorf("put account policy: %v", err)
		}
		defer tc.client.DeleteAccountPolicy(tc.ctx, &cloudwatchlogs.DeleteAccountPolicyInput{
			PolicyName: aws.String(policyName), PolicyType: types.PolicyTypeTransformerPolicy,
		})

		// An overlapping prefix rejects.
		var ipe *types.InvalidParameterException
		_, err := tc.client.PutAccountPolicy(tc.ctx, &cloudwatchlogs.PutAccountPolicyInput{
			PolicyName:        aws.String(tc.uniquePrefix("xform-policy-2")),
			PolicyType:        types.PolicyTypeTransformerPolicy,
			PolicyDocument:    aws.String(doc),
			SelectionCriteria: aws.String("LogGroupNamePrefix = " + prefix[:len(prefix)-2]),
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("overlapping prefix: %v", err)
		}

		// A group inside the prefix transforms; a group outside does not.
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(inner, "s1", `{"a":1}`, now); err != nil {
			return err
		}
		if err := tc.putLogEvent(outer, "s1", `{"a":1}`, now); err != nil {
			return err
		}

		innerArn, err := tc.findLogGroupARN(inner)
		if err != nil {
			return err
		}
		startResp, err := tc.client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime: aws.Int64(now/1000 - 60), EndTime: aws.Int64(now/1000 + 60),
			LogGroupIdentifiers: []string{aws.ToString(innerArn)},
			QueryString:         aws.String(`filter scope = "account"`),
		})
		if err != nil {
			return fmt.Errorf("start inner query: %v", err)
		}
		innerResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		if len(innerResp.Results) == 0 {
			return fmt.Errorf("account policy did not transform the prefixed group")
		}

		outerArn, err := tc.findLogGroupARN(outer)
		if err != nil {
			return err
		}
		startResp, err = tc.client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime: aws.Int64(now/1000 - 60), EndTime: aws.Int64(now/1000 + 60),
			LogGroupIdentifiers: []string{aws.ToString(outerArn)},
			QueryString:         aws.String(`filter scope = "account"`),
		})
		if err != nil {
			return fmt.Errorf("start outer query: %v", err)
		}
		outerResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		if len(outerResp.Results) != 0 {
			return fmt.Errorf("account policy transformed a group outside its prefix")
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "Transformer_BuiltinVendedParsers", func() error {
		// parseVPC: the documented example over the dry-run seam.
		resp, err := tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{{ParseVPC: &types.ParseVPC{}}},
			LogEventMessages: []string{
				"2 123456789010 eni-abc123de 192.0.2.0 192.0.2.24 20641 22 6 20 4249 1418530010 1418530070 ACCEPT OK",
			},
		})
		if err != nil {
			return fmt.Errorf("test transformer parseVPC: %v", err)
		}
		vpc, err := transformedRecordJSON(resp.TransformedLogs[0])
		if err != nil {
			return err
		}
		if vpc["accountId"] != "123456789010" {
			return fmt.Errorf("parseVPC accountId: %v (%T)", vpc["accountId"], vpc["accountId"])
		}
		if vpc["srcPort"] != float64(20641) || vpc["version"] != float64(2) {
			return fmt.Errorf("parseVPC numeric typing: %v %v", vpc["srcPort"], vpc["version"])
		}
		if vpc["action"] != "ACCEPT" || vpc["logStatus"] != "OK" {
			return fmt.Errorf("parseVPC string fields: %v %v", vpc["action"], vpc["logStatus"])
		}

		// parseCloudfront: URL decoding plus integer/double typing.
		resp, err = tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{{ParseCloudfront: &types.ParseCloudfront{}}},
			LogEventMessages: []string{
				"2019-12-04  21:02:31   LAX1   392    192.0.2.24    GET    d111111abcdef8.cloudfront.net  /index.html    200    -  Mozilla/5.0%20(Windows%20NT%2010.0;%20Win64;%20x64)%20AppleWebKit/537.36%20(KHTML,%20like%20Gecko)%20Chrome/78.0.3904.108%20Safari/537.36  -  -  Hit    SOX4xwn4XV6Q4rgb7XiVGOHms_BGlTAC4KyHmureZmBNrjGdRLiNIQ==   d111111abcdef8.cloudfront.net  https  23 0.001  -  TLSv1.2    ECDHE-RSA-AES128-GCM-SHA256    Hit    HTTP/2.0   -  -  11040  0.001  Hit    text/html  78 -  -",
			},
		})
		if err != nil {
			return fmt.Errorf("test transformer parseCloudfront: %v", err)
		}
		cf, err := transformedRecordJSON(resp.TransformedLogs[0])
		if err != nil {
			return err
		}
		if ua, _ := cf["cs(User-Agent)"].(string); !strings.Contains(ua, "(Windows NT 10.0; Win64; x64)") {
			return fmt.Errorf("parseCloudfront user agent decoding: %v", cf["cs(User-Agent)"])
		}
		if cf["sc-bytes"] != float64(392) || cf["time-taken"] != float64(0.001) {
			return fmt.Errorf("parseCloudfront numeric typing: %v %v", cf["sc-bytes"], cf["time-taken"])
		}

		// parseWAF: headers and labels arrays become JSON objects.
		resp, err = tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{{ParseWAF: &types.ParseWAF{}}},
			LogEventMessages: []string{
				`{"timestamp":1576280412771,"formatVersion":1,"action":"BLOCK","httpRequest":{"clientIp":"1.1.1.1","headers":[{"name":"Host","value":"localhost:1989"},{"name":"User-Agent","value":"curl/7.61.1"}],"uri":"/myUri","httpVersion":"HTTP/1.1","httpMethod":"GET"},"labels":[{"name":"value"}]}`,
			},
		})
		if err != nil {
			return fmt.Errorf("test transformer parseWAF: %v", err)
		}
		waf, err := transformedRecordJSON(resp.TransformedLogs[0])
		if err != nil {
			return err
		}
		httpRequest, _ := waf["httpRequest"].(map[string]interface{})
		headers, _ := httpRequest["headers"].(map[string]interface{})
		if headers["Host"] != "localhost:1989" || headers["User-Agent"] != "curl/7.61.1" {
			return fmt.Errorf("parseWAF headers object: %v", httpRequest["headers"])
		}
		labels, _ := waf["labels"].(map[string]interface{})
		if labels["name"] != "value" {
			return fmt.Errorf("parseWAF labels object: %v", waf["labels"])
		}

		// parseRoute53: the version types as a double, the rest stay
		// strings.
		resp, err = tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{{ParseRoute53: &types.ParseRoute53{}}},
			LogEventMessages: []string{
				"1.0 2017-12-13T08:15:50.235Z Z123412341234 example.com AAAA NOERROR TCP IAD12 192.0.2.0 198.51.100.0/24",
			},
		})
		if err != nil {
			return fmt.Errorf("test transformer parseRoute53: %v", err)
		}
		r53, err := transformedRecordJSON(resp.TransformedLogs[0])
		if err != nil {
			return err
		}
		if r53["version"] != float64(1) || r53["hostZoneId"] != "Z123412341234" || r53["ednsClientSubnet"] != "198.51.100.0/24" {
			return fmt.Errorf("parseRoute53 fields: %v", r53)
		}

		// parsePostgres: the six documented fields.
		resp, err = tc.client.TestTransformer(tc.ctx, &cloudwatchlogs.TestTransformerInput{
			TransformerConfig: []types.Processor{{ParsePostgres: &types.ParsePostgres{}}},
			LogEventMessages: []string{
				`2019-03-10 03:54:59 UTC:10.0.0.123(52834):postgres@logtestdb:[20175]:ERROR: column "wrong_column_name" does not exist at character 8`,
			},
		})
		if err != nil {
			return fmt.Errorf("test transformer parsePostgres: %v", err)
		}
		pg, err := transformedRecordJSON(resp.TransformedLogs[0])
		if err != nil {
			return err
		}
		if pg["logTime"] != "2019-03-10 03:54:59 UTC" || pg["srcIp"] != "10.0.0.123(52834)" ||
			pg["userName"] != "postgres" || pg["dbName"] != "logtestdb" ||
			pg["processId"] != "20175" || pg["logLevel"] != "ERROR" {
			return fmt.Errorf("parsePostgres fields: %v", pg)
		}

		// The applied-at-ingestion path: a group-level parseVPC
		// transformer substitutes the parsed fields for the query plane.
		group, cleanupGroup, err := tc.newLogGroupFixture("xformbuiltin")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		groupArn, err := tc.findLogGroupARN(group)
		if err != nil {
			return err
		}
		if _, err := tc.client.PutTransformer(tc.ctx, &cloudwatchlogs.PutTransformerInput{
			LogGroupIdentifier: groupArn,
			TransformerConfig:  []types.Processor{{ParseVPC: &types.ParseVPC{}}},
		}); err != nil {
			return fmt.Errorf("put transformer: %v", err)
		}
		const stream = "s1"
		if err := tc.createLogStream(group, stream); err != nil {
			return err
		}
		if err := tc.putLogEvent(group, stream,
			"2 123456789010 eni-abc123de 192.0.2.0 192.0.2.24 20641 22 6 20 4249 1418530010 1418530070 ACCEPT OK",
			time.Now().UnixMilli()); err != nil {
			return err
		}
		nowMs := time.Now().UnixMilli()
		start, err := tc.client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:           aws.Int64(nowMs/1000 - 3600),
			EndTime:             aws.Int64(nowMs/1000 + 60),
			LogGroupIdentifiers: []string{aws.ToString(groupArn)},
			QueryString:         aws.String(`filter ispresent(interfaceId) | stats count() as flows`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		queryResp, err := tc.waitForQuery(start.QueryId, 20)
		if err != nil {
			return err
		}
		if len(queryResp.Results) == 0 || aws.ToString(queryResp.Results[0][0].Value) != "1" {
			var dump string
			for _, row := range queryResp.Results {
				for _, f := range row {
					dump += aws.ToString(f.Field) + "=" + aws.ToString(f.Value) + ";"
				}
			}
			return fmt.Errorf("parsed flow not visible to the query plane: rows=%d [%s]", len(queryResp.Results), dump)
		}
		return nil
	}))

	return results
}

// transformedRecordJSON decodes one TransformedLogRecord's transformed
// message as the JSON object the built-in parser assertions read.
func transformedRecordJSON(record types.TransformedLogRecord) (map[string]interface{}, error) {
	transformed := aws.ToString(record.TransformedEventMessage)
	if transformed == "" {
		return nil, fmt.Errorf("record carries no transformed message")
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(transformed), &parsed); err != nil {
		return nil, fmt.Errorf("transformed message %q is not a JSON object: %v", transformed, err)
	}
	return parsed, nil
}
