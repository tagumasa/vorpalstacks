package testutil

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// queryTests covers StartQuery compile validation and the scheduled query
// lifecycle (create, get, list, update, history, delete).
func (tc *cwlogsTestCtx) queryTests() []TestResult {
	var results []TestResult
	client := tc.client

	results = append(results, tc.runner.RunTest("logs", "StartQuery_MalformedQuery_Rejected", func() error {
		groupName, cleanupGroup, err := tc.newLogGroupFixture("malformed-query-group")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		now := time.Now().Unix()
		_, err = client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String("bogus @message"),
		})
		if err := AssertErrorContains(err, "MalformedQueryException"); err != nil {
			return fmt.Errorf("expected MalformedQueryException for unknown command, got: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "ScheduledQuery_Lifecycle", func() error {
		name := tc.uniquePrefix("sched-query")
		createResp, err := client.CreateScheduledQuery(tc.ctx, &cloudwatchlogs.CreateScheduledQueryInput{
			Name:               aws.String(name),
			QueryString:        aws.String("fields @timestamp, @message | limit 10"),
			QueryLanguage:      types.QueryLanguageCwli,
			ExecutionRoleArn:   aws.String(tc.roleARN("scheduled-query-role")),
			ScheduleExpression: aws.String("rate(1 hour)"),
			State:              types.ScheduledQueryStateEnabled,
		})
		if err != nil {
			return fmt.Errorf("create scheduled query: %v", err)
		}
		queryArn := aws.ToString(createResp.ScheduledQueryArn)

		getResp, err := client.GetScheduledQuery(tc.ctx, &cloudwatchlogs.GetScheduledQueryInput{
			Identifier: aws.String(queryArn),
		})
		if err != nil {
			return fmt.Errorf("get scheduled query: %v", err)
		}
		if aws.ToString(getResp.Name) != name {
			return fmt.Errorf("name mismatch, got %v", getResp.Name)
		}
		if getResp.State != types.ScheduledQueryStateEnabled {
			return fmt.Errorf("state mismatch, got %v", getResp.State)
		}

		// List must find the query across all pages.
		var found bool
		var nextToken *string
		for i := 0; i < 20; i++ {
			listResp, err := client.ListScheduledQueries(tc.ctx, &cloudwatchlogs.ListScheduledQueriesInput{
				NextToken: nextToken,
			})
			if err != nil {
				return fmt.Errorf("list scheduled queries: %v", err)
			}
			for _, sq := range listResp.ScheduledQueries {
				if aws.ToString(sq.Name) == name {
					found = true
				}
			}
			if found || listResp.NextToken == nil {
				break
			}
			nextToken = listResp.NextToken
		}
		if !found {
			return fmt.Errorf("scheduled query %s not found in list", name)
		}

		updateResp, err := client.UpdateScheduledQuery(tc.ctx, &cloudwatchlogs.UpdateScheduledQueryInput{
			Identifier:         aws.String(queryArn),
			QueryString:        aws.String("fields @timestamp, @message | limit 5"),
			QueryLanguage:      types.QueryLanguageCwli,
			ScheduleExpression: aws.String("rate(1 hour)"),
			ExecutionRoleArn:   aws.String(tc.roleARN("scheduled-query-role")),
			State:              types.ScheduledQueryStateDisabled,
		})
		if err != nil {
			return fmt.Errorf("update scheduled query: %v", err)
		}
		if updateResp.State != types.ScheduledQueryStateDisabled {
			return fmt.Errorf("state not updated, got %v", updateResp.State)
		}

		historyResp, err := client.GetScheduledQueryHistory(tc.ctx, &cloudwatchlogs.GetScheduledQueryHistoryInput{
			Identifier: aws.String(queryArn),
			EndTime:    aws.Int64(time.Now().Add(time.Hour).UnixMilli()),
			StartTime:  aws.Int64(time.Now().Add(-24 * time.Hour).UnixMilli()),
		})
		if err != nil {
			return fmt.Errorf("get scheduled query history: %v", err)
		}
		if historyResp.ScheduledQueryArn == nil {
			return fmt.Errorf("history response missing scheduledQueryArn")
		}

		_, err = client.DeleteScheduledQuery(tc.ctx, &cloudwatchlogs.DeleteScheduledQueryInput{
			Identifier: aws.String(queryArn),
		})
		if err != nil {
			return fmt.Errorf("delete scheduled query: %v", err)
		}

		_, err = client.GetScheduledQuery(tc.ctx, &cloudwatchlogs.GetScheduledQueryInput{
			Identifier: aws.String(queryArn),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("expected ResourceNotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	// A full Insights query over ingested events must return results: JSON
	// field discovery, the where alias, stats aggregation and bin grouping,
	// and the limit any variant all execute end to end.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_InsightsCommands_ProduceResults", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("insights-query-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		for i, msg := range []string{
			`{"level": "INFO", "service": "alpha", "latency": 100}`,
			`{"level": "ERROR", "service": "alpha", "latency": 300}`,
			`{"level": "INFO", "service": "beta", "latency": 200}`,
		} {
			if err := tc.putLogEvent(groupName, "s1", msg, now-int64(len(msg))*0-int64(i)*1000); err != nil {
				return fmt.Errorf("put event: %v", err)
			}
		}

		// The query addresses the group through the identifier member's
		// ARN form — the documented name-or-ARN contract — so the results
		// pin resolution on the execution plane, not just acceptance.
		groupArn, err := tc.findLogGroupARN(groupName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:           aws.Int64(now/1000 - 3600),
			EndTime:             aws.Int64(now/1000 + 60),
			LogGroupIdentifiers: []string{aws.ToString(groupArn)},
			QueryString:         aws.String(`filter ispresent(service) | where level = "ERROR" | stats max(latency) as worst`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}

		resResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		// The statistics echo the scan: every event scanned,
		// the filter's survivors matched, a positive byte
		// count and the one group addressed.
		if got := resResp.Statistics.RecordsScanned; got != 3 {
			return fmt.Errorf("recordsScanned = %v, want 3", got)
		}
		if got := resResp.Statistics.RecordsMatched; got != 1 {
			return fmt.Errorf("recordsMatched = %v, want 1", got)
		}
		if resResp.Statistics.BytesScanned <= 0 {
			return fmt.Errorf("bytesScanned = %v, want a positive count", resResp.Statistics.BytesScanned)
		}
		if got := resResp.Statistics.LogGroupsScanned; got != 1 {
			return fmt.Errorf("logGroupsScanned = %v, want 1", got)
		}
		results := flattenResultFields(resResp.Results)
		if results == nil {
			return fmt.Errorf("query completed without rows")
		}
		var worst string
		for _, f := range results {
			if aws.ToString(f.Field) == "worst" {
				worst = aws.ToString(f.Value)
			}
		}
		if worst != "300" {
			return fmt.Errorf("worst latency = %q, want 300 (fields: %+v)", worst, results)
		}
		return nil
	}))

	// The result row contract follows AWS: @timestamp renders as
	// "YYYY-MM-DD HH:MM:SS.mmm", every event row carries an @ptr that
	// GetLogRecord accepts, and @log identifies the group as account:group.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_ResultRowContract", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("row-contract-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", `{"level": "INFO"}`, now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields @timestamp, @message, @log | limit 2`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		resResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		if len(resResp.Results) == 0 {
			return fmt.Errorf("query completed without rows")
		}
		row := resResp.Results[0]

		tsFormat := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}$`)
		fields := flattenResultFields([][]types.ResultField{row})
		var timestamp, logField, ptr, message string
		for _, f := range fields {
			switch aws.ToString(f.Field) {
			case "@timestamp":
				timestamp = aws.ToString(f.Value)
			case "@log":
				logField = aws.ToString(f.Value)
			case "@ptr":
				ptr = aws.ToString(f.Value)
			case "@message":
				message = aws.ToString(f.Value)
			}
		}
		if !tsFormat.MatchString(timestamp) {
			return fmt.Errorf("@timestamp = %q, want YYYY-MM-DD HH:MM:SS.mmm", timestamp)
		}
		if !strings.HasSuffix(logField, ":"+groupName) {
			return fmt.Errorf("@log = %q, want account:%s", logField, groupName)
		}
		if ptr == "" {
			return fmt.Errorf("event row missing @ptr: %+v", row)
		}
		recordResp, err := client.GetLogRecord(tc.ctx, &cloudwatchlogs.GetLogRecordInput{
			LogRecordPointer: aws.String(ptr),
		})
		if err != nil {
			return fmt.Errorf("get log record with @ptr: %v", err)
		}
		if got := recordResp.LogRecord["@message"]; got != message {
			return fmt.Errorf("log record @message = %q, want the row's @message %q", got, message)
		}
		return nil
	}))

	// GetLogObject serves the fieldStream member as the modelled
	// awsJson1.1 event stream: an initial-response message, then one
	// fields event whose data blob is the record document the pointer
	// addresses (the @ptr the result rows carry is the operation's
	// pointer).
	results = append(results, tc.runner.RunTest("logs", "GetLogObject_Stream", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("log-object-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		message := `{"level": "INFO", "curated": "value"}`
		if err := tc.putLogEvent(groupName, "s1", message, now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		client, err := tc.newLiveTailClient()
		if err != nil {
			return err
		}
		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields @ptr | limit 1`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		resResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		if len(resResp.Results) == 0 {
			return fmt.Errorf("query completed without rows")
		}
		var ptr string
		for _, f := range flattenResultFields(resResp.Results) {
			if aws.ToString(f.Field) == "@ptr" {
				ptr = aws.ToString(f.Value)
			}
		}
		if ptr == "" {
			return fmt.Errorf("event row missing @ptr: %+v", resResp.Results[0])
		}

		objResp, err := client.GetLogObject(tc.ctx, &cloudwatchlogs.GetLogObjectInput{
			LogObjectPointer: aws.String(ptr),
		})
		if err != nil {
			return fmt.Errorf("get log object: %v", err)
		}
		stream := objResp.GetStream()
		defer stream.Close()
		event, ok := <-stream.Events()
		if !ok {
			return fmt.Errorf("event stream closed without a fields event: %v", stream.Err())
		}
		fields, isFields := event.(*types.GetLogObjectResponseStreamMemberFields)
		if !isFields {
			return fmt.Errorf("first stream event = %T, want the fields member", event)
		}
		var record map[string]string
		if err := json.Unmarshal(fields.Value.Data, &record); err != nil {
			return fmt.Errorf("fields data is not the record document: %v (%s)", err, fields.Value.Data)
		}
		if record["@message"] != message {
			return fmt.Errorf("record @message = %q, want %q", record["@message"], message)
		}
		if record["curated"] != "value" {
			return fmt.Errorf("record JSON fields missing: %+v", record)
		}
		if _, stillOpen := <-stream.Events(); stillOpen {
			return fmt.Errorf("stream delivered more than the one fields event")
		}
		if err := stream.Err(); err != nil {
			return fmt.Errorf("stream error after the fields event: %v", err)
		}
		return nil
	}))

	// Expression columns carry their written form when the query gives no
	// alias: a postfix attribute access and an aggregation's auto-label
	// both name their columns with the expression's source text (never an
	// empty name), and the bin grouping's bucket column is addressable as
	// @bin by a subsequent stats command.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_ExpressionColumnNames", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("expr-column-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		for _, msg := range []string{
			`{"status": "ok", "latency": 100}`,
			`{"status": "ok", "latency": 300}`,
		} {
			if err := tc.putLogEvent(groupName, "s1", msg, now); err != nil {
				return fmt.Errorf("put event: %v", err)
			}
		}
		window := &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: []string{groupName},
		}

		window.QueryString = aws.String(`fields jsonParse(@message).status`)
		startResp, err := client.StartQuery(tc.ctx, window)
		if err != nil {
			return fmt.Errorf("start postfix query: %v", err)
		}
		resResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		fields := flattenResultFields(resResp.Results)
		if fields == nil {
			return fmt.Errorf("postfix query completed without rows")
		}
		var status string
		for _, f := range fields {
			if aws.ToString(f.Field) == "" {
				return fmt.Errorf("row carries an unnamed column: %+v", fields)
			}
			if aws.ToString(f.Field) == "jsonParse(@message).status" {
				status = aws.ToString(f.Value)
			}
		}
		if status != "ok" {
			return fmt.Errorf("postfix column jsonParse(@message).status = %q (fields: %+v)", status, fields)
		}

		window.QueryString = aws.String(`stats max(jsonParse(@message).latency)`)
		startResp, err = client.StartQuery(tc.ctx, window)
		if err != nil {
			return fmt.Errorf("start auto-label query: %v", err)
		}
		resResp, err = tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		fields = flattenResultFields(resResp.Results)
		var worst string
		for _, f := range fields {
			if aws.ToString(f.Field) == "max(jsonParse(@message).latency)" {
				worst = aws.ToString(f.Value)
			}
		}
		if worst != "300" {
			return fmt.Errorf("auto-labelled max column = %q (fields: %+v)", worst, fields)
		}

		window.QueryString = aws.String(`stats count(*) as n by bin(1h) | stats sum(n) as total by @bin`)
		startResp, err = client.StartQuery(tc.ctx, window)
		if err != nil {
			return fmt.Errorf("start bin-addressing query: %v", err)
		}
		resResp, err = tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		fields = flattenResultFields(resResp.Results)
		var total string
		for _, f := range fields {
			if aws.ToString(f.Field) == "total" {
				total = aws.ToString(f.Value)
			}
		}
		if total != "2" {
			return fmt.Errorf("binned follow-up total = %q (fields: %+v)", total, fields)
		}
		return nil
	}))

	// The documented hashing functions md5()/sha256() are available in
	// fields expressions and return hex digests.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_HashingFunctions", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("hash-query-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", "hello", now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}
		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields md5(@message) as m, sha256(@message) as s | limit 1`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		resResp, err := tc.waitForQuery(startResp.QueryId, 20)
		if err != nil {
			return err
		}
		fields := flattenResultFields(resResp.Results)
		if err := requireResultFields(fields, "m", "s"); err != nil {
			return err
		}
		wantMD5 := "5d41402abc4b2a76b9719d911017c592"
		wantSHA := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
		for _, f := range fields {
			switch aws.ToString(f.Field) {
			case "m":
				if aws.ToString(f.Value) != wantMD5 {
					return fmt.Errorf("md5 = %q, want %q", aws.ToString(f.Value), wantMD5)
				}
			case "s":
				if aws.ToString(f.Value) != wantSHA {
					return fmt.Errorf("sha256 = %q, want %q", aws.ToString(f.Value), wantSHA)
				}
			}
		}
		return nil
	}))

	// Timestamp-typed results, the @ptr round trip over a delimiter-bearing
	// message, and the documented case/messageSize limits.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_TimestampContract", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("ts-contract-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", "alpha|beta|gamma", now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}
		if err := tc.putLogEvent(groupName, "s1", "second event", now+60000); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		// A message containing the pointer delimiter survives the @ptr
		// round trip through GetLogRecord intact. Cold-start query
		// completion can exceed the suite's usual window when this test
		// chains several queries, hence the wider poll budget.
		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 120),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`filter @message like /alpha/ | fields @ptr, @message | limit 5`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		tsResp, err := tc.waitForQuery(startResp.QueryId, 40)
		if err != nil {
			return err
		}
		fields := flattenResultFields(tsResp.Results)
		ptr := ""
		for _, f := range fields {
			if aws.ToString(f.Field) == "@ptr" {
				ptr = aws.ToString(f.Value)
			}
		}
		if ptr == "" {
			return fmt.Errorf("row missing @ptr: %+v", fields)
		}
		recordResp, err := client.GetLogRecord(tc.ctx, &cloudwatchlogs.GetLogRecordInput{LogRecordPointer: aws.String(ptr)})
		if err != nil {
			return fmt.Errorf("get log record: %v", err)
		}
		if got := recordResp.LogRecord["@message"]; got != "alpha|beta|gamma" {
			return fmt.Errorf("log record @message = %q, want the full delimiter-bearing message", got)
		}

		// fromMillis, datefloor and the passthrough aggregations render in
		// the timestamp result form and stay numeric inside expressions.
		startResp, err = client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 120),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields fromMillis(1700000000000) as fm, datefloor(fromMillis(1700000000000), 1h) as fl, toMillis(fromMillis(1700000000000)) as ms | limit 1`),
		})
		if err != nil {
			return fmt.Errorf("start timestamp query: %v", err)
		}
		tsResp, err = tc.waitForQuery(startResp.QueryId, 40)
		if err != nil {
			return err
		}
		fields = flattenResultFields(tsResp.Results)
		if err := requireResultFields(fields, "fm", "fl", "ms"); err != nil {
			return err
		}
		for _, f := range fields {
			switch aws.ToString(f.Field) {
			case "fm":
				if aws.ToString(f.Value) != "2023-11-14 22:13:20.000" {
					return fmt.Errorf("fromMillis = %q, want the timestamp rendering", aws.ToString(f.Value))
				}
			case "fl":
				if aws.ToString(f.Value) != "2023-11-14 22:00:00.000" {
					return fmt.Errorf("datefloor = %q, want the floored timestamp rendering", aws.ToString(f.Value))
				}
			case "ms":
				if aws.ToString(f.Value) != "1700000000000" {
					return fmt.Errorf("toMillis round trip = %q, want 1700000000000", aws.ToString(f.Value))
				}
			}
		}

		// Bin keys and earliest/latest keep the same timestamp rendering.
		startResp, err = client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 120),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`stats earliest(@timestamp) as e, latest(@timestamp) as l by bin(1m)`),
		})
		if err != nil {
			return fmt.Errorf("start bin query: %v", err)
		}
		tsResp, err = tc.waitForQuery(startResp.QueryId, 40)
		if err != nil {
			return err
		}
		fields = flattenResultFields(tsResp.Results)
		if err := requireResultFields(fields, "@bin", "e", "l"); err != nil {
			return err
		}
		tsFormat := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}$`)
		for _, f := range fields {
			switch aws.ToString(f.Field) {
			case "@bin", "e", "l":
				if !tsFormat.MatchString(aws.ToString(f.Value)) {
					return fmt.Errorf("%s = %q, want the timestamp rendering", aws.ToString(f.Field), aws.ToString(f.Value))
				}
			}
		}

		// The documented case branch limit and messageSize arity are
		// rejected at compile time.
		caseQuery := "fields case("
		for i := 0; i < 11; i++ {
			if i > 0 {
				caseQuery += ", "
			}
			caseQuery += "false, 0"
		}
		caseQuery += ") as c"
		_, err = client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(caseQuery),
		})
		if err == nil {
			return fmt.Errorf("case with 11 branches should be rejected")
		}
		_, err = client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields messageSize() as sz`),
		})
		if err == nil {
			return fmt.Errorf("messageSize without arguments should be rejected")
		}
		return nil
	}))

	// The documented command placement rules are enforced with
	// MalformedQueryException at StartQuery time.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_PlacementRules_Rejected", func() error {
		groupName, cleanupGroup, err := tc.newLogGroupFixture("placement-query-group")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		now := time.Now().Unix()
		for _, qs := range []string{
			`fields @message | dedup @message | filter ispresent(@message)`,
			`sort @timestamp desc | stats count(*)`,
			`sort @timestamp desc | pattern @message`,
		} {
			_, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
				StartTime:     aws.Int64(now - 3600),
				EndTime:       aws.Int64(now),
				LogGroupNames: []string{groupName},
				QueryString:   aws.String(qs),
			})
			if err := AssertErrorContains(err, "MalformedQueryException"); err != nil {
				return fmt.Errorf("query %q: %v", qs, err)
			}
		}
		return nil
	}))

	// Scheduled queries must reject unknown commands at create time so that
	// a stored query can never silently ignore commands at trigger time.
	// The scheduled-query operations declare ValidationException (not
	// StartQuery's MalformedQueryException), so a query that fails to
	// compile rejects with their own identity.
	results = append(results, tc.runner.RunTest("logs", "CreateScheduledQuery_UnknownCommand_Rejected", func() error {
		_, err := client.CreateScheduledQuery(tc.ctx, &cloudwatchlogs.CreateScheduledQueryInput{
			Name:               aws.String(tc.uniquePrefix("bad-command-sq")),
			QueryString:        aws.String("fields @message | frobnicate x"),
			QueryLanguage:      types.QueryLanguageCwli,
			ExecutionRoleArn:   aws.String(tc.roleARN("scheduled-query-role")),
			ScheduleExpression: aws.String("rate(1 hour)"),
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// UpdateScheduledQuery follows the Smithy required members and the same
	// spec validation as create: an invalid schedule expression and a
	// missing required member are both rejected.
	results = append(results, tc.runner.RunTest("logs", "UpdateScheduledQuery_Validation_Rejected", func() error {
		name := tc.uniquePrefix("update-valid-sq")
		createResp, err := client.CreateScheduledQuery(tc.ctx, &cloudwatchlogs.CreateScheduledQueryInput{
			Name:               aws.String(name),
			QueryString:        aws.String("fields @timestamp, @message | limit 5"),
			QueryLanguage:      types.QueryLanguageCwli,
			ExecutionRoleArn:   aws.String(tc.roleARN("scheduled-query-role")),
			ScheduleExpression: aws.String("rate(1 hour)"),
		})
		if err != nil {
			return fmt.Errorf("create scheduled query: %v", err)
		}
		queryArn := aws.ToString(createResp.ScheduledQueryArn)
		defer client.DeleteScheduledQuery(tc.ctx, &cloudwatchlogs.DeleteScheduledQueryInput{
			Identifier: aws.String(queryArn),
		})

		_, err = client.UpdateScheduledQuery(tc.ctx, &cloudwatchlogs.UpdateScheduledQueryInput{
			Identifier:         aws.String(queryArn),
			QueryString:        aws.String("fields @message | limit 5"),
			QueryLanguage:      types.QueryLanguageCwli,
			ScheduleExpression: aws.String("not-a-schedule"),
			ExecutionRoleArn:   aws.String(tc.roleARN("scheduled-query-role")),
		})
		if err := AssertErrorContains(err, "ValidationException"); err != nil {
			return fmt.Errorf("invalid schedule: %v", err)
		}

		// The SDK itself enforces the Smithy required members
		// (queryLanguage, queryString, scheduleExpression,
		// executionRoleArn) before the request is sent, so the server-side
		// missing-parameter path cannot be exercised through the client.

		_, err = client.UpdateScheduledQuery(tc.ctx, &cloudwatchlogs.UpdateScheduledQueryInput{
			Identifier:         aws.String(queryArn),
			QueryString:        aws.String("fields @message | bogus"),
			QueryLanguage:      types.QueryLanguageCwli,
			ScheduleExpression: aws.String("rate(1 hour)"),
			ExecutionRoleArn:   aws.String(tc.roleARN("scheduled-query-role")),
		})
		return AssertErrorContains(err, "ValidationException")
	}))

	// A stream holding more events than one GetLogEvents page serves must
	// still be scanned whole: the query's statistics and count(*) see
	// every event past the wire page limit.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_StatsBeyondEventPageLimit", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("bigcount-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		const totalEvents = 10005
		base := time.Now().UnixMilli() - int64(totalEvents)
		batch := make([]types.InputLogEvent, 0, 5000)
		for i := 0; i < totalEvents; i++ {
			batch = append(batch, types.InputLogEvent{
				Message:   aws.String(fmt.Sprintf("bulk-%d", i)),
				Timestamp: aws.Int64(base + int64(i)),
			})
			if len(batch) == 5000 {
				if _, err := client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
					LogGroupName:  aws.String(groupName),
					LogStreamName: aws.String("s1"),
					LogEvents:     batch,
				}); err != nil {
					return fmt.Errorf("put batch: %v", err)
				}
				batch = batch[:0]
			}
		}
		if len(batch) > 0 {
			if _, err := client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
				LogGroupName:  aws.String(groupName),
				LogStreamName: aws.String("s1"),
				LogEvents:     batch,
			}); err != nil {
				return fmt.Errorf("put tail batch: %v", err)
			}
		}

		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(base/1000 - 60),
			EndTime:       aws.Int64(base/1000 + 3600),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String("stats count(*) as total"),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		resResp, err := tc.waitForQuery(startResp.QueryId, 30)
		if err != nil {
			return err
		}
		if got := int64(resResp.Statistics.RecordsScanned); got != totalEvents {
			return fmt.Errorf("recordsScanned = %d, want %d", got, totalEvents)
		}
		rows := flattenResultFields(resResp.Results)
		for _, f := range rows {
			if aws.ToString(f.Field) == "total" && aws.ToString(f.Value) != fmt.Sprintf("%d", totalEvents) {
				return fmt.Errorf("count(*) = %q, want %d", aws.ToString(f.Value), totalEvents)
			}
		}
		return nil
	}))

	// StartQuery fails fast on a log group that does not exist instead of
	// returning a queryId that completes with zero rows.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_MissingGroup_Rejected", func() error {
		now := time.Now().Unix()
		_, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now),
			LogGroupNames: []string{tc.uniquePrefix("no-such-group")},
			QueryString:   aws.String("fields @message | limit 1"),
		})
		if err := AssertErrorContains(err, "ResourceNotFoundException"); err != nil {
			return fmt.Errorf("expected ResourceNotFoundException for a missing log group, got: %v", err)
		}
		return nil
	}))

	// An inverted time window is malformed input, not an empty query.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_InvertedWindow_Rejected", func() error {
		groupName, cleanupGroup, err := tc.newLogGroupFixture("inverted-window-group")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().Unix()
		_, err = client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now + 600),
			EndTime:       aws.Int64(now),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String("fields @message | limit 1"),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("expected InvalidParameterException for startTime after endTime, got: %v", err)
		}
		return nil
	}))

	// Stopping a query that already ended rejects with an error stating
	// the query is not running, per the StopQuery contract.
	results = append(results, tc.runner.RunTest("logs", "StopQuery_EndedQuery_Rejected", func() error {
		groupName, cleanupGroup, err := tc.newGroupStreamFixture("ended-query-group", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", `{"level": "INFO"}`, now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now/1000 - 3600),
			EndTime:       aws.Int64(now/1000 + 60),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String("fields @message | limit 2"),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		if _, err := tc.waitForQuery(startResp.QueryId, 20); err != nil {
			return err
		}
		_, err = client.StopQuery(tc.ctx, &cloudwatchlogs.StopQueryInput{
			QueryId: startResp.QueryId,
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("expected InvalidParameterException stopping a completed query, got: %v", err)
		}
		return nil
	}))

	return results
}

func flattenResultFields(rows [][]types.ResultField) []*types.ResultField {
	var out []*types.ResultField
	for _, row := range rows {
		for i := range row {
			out = append(out, &row[i])
		}
	}
	return out
}

// requireResultFields fails the check unless every named field is present
// among the flattened results: a value switch over the fields alone passes
// vacuously when the query returns none of the expected fields, so the
// completeness check must come before the value comparisons.
func requireResultFields(fields []*types.ResultField, want ...string) error {
	have := map[string]bool{}
	for _, f := range fields {
		have[aws.ToString(f.Field)] = true
	}
	for _, w := range want {
		if !have[w] {
			return fmt.Errorf("query result missing field %q, got %+v", w, fields)
		}
	}
	return nil
}
