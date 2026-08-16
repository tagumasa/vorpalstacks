package testutil

import (
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
		groupName := tc.uniquePrefix("malformed-query-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)

		now := time.Now().Unix()
		_, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
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
			ExecutionRoleArn:   aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
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
			ExecutionRoleArn:   aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
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
		groupName := tc.uniquePrefix("insights-query-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)
		if err := tc.createLogStream(groupName, "s1"); err != nil {
			return fmt.Errorf("create log stream: %v", err)
		}
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

		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 60000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`filter ispresent(service) | where level = "ERROR" | stats max(latency) as worst`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}

		var results []*types.ResultField
		for i := 0; i < 20; i++ {
			resResp, err := client.GetQueryResults(tc.ctx, &cloudwatchlogs.GetQueryResultsInput{
				QueryId: startResp.QueryId,
			})
			if err != nil {
				return fmt.Errorf("get results: %v", err)
			}
			if resResp.Status == types.QueryStatusComplete {
				results = flattenResultFields(resResp.Results)
				break
			}
			if resResp.Status == types.QueryStatusFailed || resResp.Status == types.QueryStatusCancelled {
				return fmt.Errorf("query finished with status %s", resResp.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		if results == nil {
			return fmt.Errorf("query did not complete in time")
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
		groupName := tc.uniquePrefix("row-contract-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)
		if err := tc.createLogStream(groupName, "s1"); err != nil {
			return fmt.Errorf("create log stream: %v", err)
		}
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", `{"level": "INFO"}`, now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 60000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields @timestamp, @message, @log | limit 2`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		var row []types.ResultField
		for i := 0; i < 20; i++ {
			resResp, err := client.GetQueryResults(tc.ctx, &cloudwatchlogs.GetQueryResultsInput{
				QueryId: startResp.QueryId,
			})
			if err != nil {
				return fmt.Errorf("get results: %v", err)
			}
			if resResp.Status == types.QueryStatusComplete {
				if len(resResp.Results) == 0 {
					return fmt.Errorf("query completed without rows")
				}
				row = resResp.Results[0]
				break
			}
			if resResp.Status == types.QueryStatusFailed || resResp.Status == types.QueryStatusCancelled {
				return fmt.Errorf("query finished with status %s", resResp.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		if row == nil {
			return fmt.Errorf("query did not complete in time")
		}

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

	// The documented hashing functions md5()/sha256() are available in
	// fields expressions and return hex digests.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_HashingFunctions", func() error {
		groupName := tc.uniquePrefix("hash-query-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)
		if err := tc.createLogStream(groupName, "s1"); err != nil {
			return fmt.Errorf("create log stream: %v", err)
		}
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", "hello", now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}
		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 60000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields md5(@message) as m, sha256(@message) as s | limit 1`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		var fields []*types.ResultField
		for i := 0; i < 20; i++ {
			resResp, err := client.GetQueryResults(tc.ctx, &cloudwatchlogs.GetQueryResultsInput{
				QueryId: startResp.QueryId,
			})
			if err != nil {
				return fmt.Errorf("get results: %v", err)
			}
			if resResp.Status == types.QueryStatusComplete {
				fields = flattenResultFields(resResp.Results)
				break
			}
			if resResp.Status == types.QueryStatusFailed || resResp.Status == types.QueryStatusCancelled {
				return fmt.Errorf("query finished with status %s", resResp.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		if fields == nil {
			return fmt.Errorf("query did not complete in time")
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
		groupName := tc.uniquePrefix("ts-contract-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)
		if err := tc.createLogStream(groupName, "s1"); err != nil {
			return fmt.Errorf("create log stream: %v", err)
		}
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", "alpha|beta|gamma", now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}
		if err := tc.putLogEvent(groupName, "s1", "second event", now+60000); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		waitComplete := func(queryId *string) ([]*types.ResultField, error) {
			// Cold-start query completion can exceed the suite's usual
			// window when this test chains several queries.
			for i := 0; i < 40; i++ {
				resResp, err := client.GetQueryResults(tc.ctx, &cloudwatchlogs.GetQueryResultsInput{QueryId: queryId})
				if err != nil {
					return nil, fmt.Errorf("get results: %v", err)
				}
				if resResp.Status == types.QueryStatusComplete {
					return flattenResultFields(resResp.Results), nil
				}
				if resResp.Status == types.QueryStatusFailed || resResp.Status == types.QueryStatusCancelled {
					return nil, fmt.Errorf("query finished with status %s", resResp.Status)
				}
				time.Sleep(250 * time.Millisecond)
			}
			return nil, fmt.Errorf("query did not complete in time")
		}

		// A message containing the pointer delimiter survives the @ptr
		// round trip through GetLogRecord intact.
		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 120000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`filter @message like /alpha/ | fields @ptr, @message | limit 5`),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		fields, err := waitComplete(startResp.QueryId)
		if err != nil {
			return err
		}
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
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 120000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`fields fromMillis(1700000000000) as fm, datefloor(fromMillis(1700000000000), 1h) as fl, toMillis(fromMillis(1700000000000)) as ms | limit 1`),
		})
		if err != nil {
			return fmt.Errorf("start timestamp query: %v", err)
		}
		fields, err = waitComplete(startResp.QueryId)
		if err != nil {
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
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 120000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(`stats earliest(@timestamp) as e, latest(@timestamp) as l by bin(1m)`),
		})
		if err != nil {
			return fmt.Errorf("start bin query: %v", err)
		}
		fields, err = waitComplete(startResp.QueryId)
		if err != nil {
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
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(caseQuery),
		})
		if err == nil {
			return fmt.Errorf("case with 11 branches should be rejected")
		}
		_, err = client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now),
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
		groupName := tc.uniquePrefix("placement-query-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)

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
	results = append(results, tc.runner.RunTest("logs", "CreateScheduledQuery_UnknownCommand_Rejected", func() error {
		_, err := client.CreateScheduledQuery(tc.ctx, &cloudwatchlogs.CreateScheduledQueryInput{
			Name:               aws.String(tc.uniquePrefix("bad-command-sq")),
			QueryString:        aws.String("fields @message | frobnicate x"),
			QueryLanguage:      types.QueryLanguageCwli,
			ExecutionRoleArn:   aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
			ScheduleExpression: aws.String("rate(1 hour)"),
		})
		return AssertErrorContains(err, "MalformedQueryException")
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
			ExecutionRoleArn:   aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
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
			ExecutionRoleArn:   aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
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
			ExecutionRoleArn:   aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
		})
		return AssertErrorContains(err, "MalformedQueryException")
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
