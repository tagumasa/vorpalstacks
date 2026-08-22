package testutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// lookupTableTests covers the lookup table API family and the lookup /
// cidrlookup query commands that enrich events with table data.
func (tc *cwlogsTestCtx) lookupTableTests() []TestResult {
	var results []TestResult
	client := tc.client

	// lookupTableName allows only alphanumerics and underscores, so the
	// hyphenated uniquePrefix cannot be used here.
	uniqueTableName := func(base string) string {
		return fmt.Sprintf("%s_%d", base, time.Now().UnixNano())
	}

	results = append(results, tc.runner.RunTest("logs", "LookupTable_Lifecycle", func() error {
		name := uniqueTableName("users_lifecycle")
		body := "id,name,department\nu1,Alice,Engineering\nu2,Bob,Support\n"
		createResp, err := client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String(name),
			TableBody:       aws.String(body),
			Description:     aws.String("user reference data"),
		})
		if err != nil {
			return fmt.Errorf("create lookup table: %v", err)
		}
		arn := aws.ToString(createResp.LookupTableArn)
		if arn == "" || createResp.CreatedAt == nil {
			return fmt.Errorf("create response missing arn or createdAt: %+v", createResp)
		}
		defer client.DeleteLookupTable(tc.ctx, &cloudwatchlogs.DeleteLookupTableInput{
			LookupTableArn: aws.String(arn),
		})

		getResp, err := client.GetLookupTable(tc.ctx, &cloudwatchlogs.GetLookupTableInput{
			LookupTableArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("get lookup table: %v", err)
		}
		if aws.ToString(getResp.TableBody) != body {
			return fmt.Errorf("tableBody round trip: %q", aws.ToString(getResp.TableBody))
		}
		if aws.ToString(getResp.LookupTableName) != name {
			return fmt.Errorf("lookupTableName = %q, want %q", aws.ToString(getResp.LookupTableName), name)
		}

		listResp, err := client.DescribeLookupTables(tc.ctx, &cloudwatchlogs.DescribeLookupTablesInput{
			LookupTableNamePrefix: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe lookup tables: %v", err)
		}
		if len(listResp.LookupTables) != 1 {
			return fmt.Errorf("describe returned %d tables, want 1", len(listResp.LookupTables))
		}
		sum := listResp.LookupTables[0]
		if aws.ToInt64(sum.RecordsCount) != 2 {
			return fmt.Errorf("recordsCount = %d, want 2", aws.ToInt64(sum.RecordsCount))
		}
		if len(sum.TableFields) != 3 || sum.TableFields[0] != "id" {
			return fmt.Errorf("tableFields = %v", sum.TableFields)
		}

		updateResp, err := client.UpdateLookupTable(tc.ctx, &cloudwatchlogs.UpdateLookupTableInput{
			LookupTableArn: aws.String(arn),
			TableBody:      aws.String("id,name\ng1,Carol\n"),
		})
		if err != nil {
			return fmt.Errorf("update lookup table: %v", err)
		}
		if updateResp.LastUpdatedTime == nil {
			return fmt.Errorf("update response missing lastUpdatedTime")
		}
		listResp, err = client.DescribeLookupTables(tc.ctx, &cloudwatchlogs.DescribeLookupTablesInput{
			LookupTableNamePrefix: aws.String(name),
		})
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		if len(listResp.LookupTables) != 1 || aws.ToInt64(listResp.LookupTables[0].RecordsCount) != 1 {
			return fmt.Errorf("recordsCount after replacement = %+v, want 1", listResp.LookupTables)
		}

		_, err = client.DeleteLookupTable(tc.ctx, &cloudwatchlogs.DeleteLookupTableInput{
			LookupTableArn: aws.String(arn),
		})
		if err != nil {
			return fmt.Errorf("delete lookup table: %v", err)
		}
		_, err = client.GetLookupTable(tc.ctx, &cloudwatchlogs.GetLookupTableInput{
			LookupTableArn: aws.String(arn),
		})
		if err2 := AssertErrorContains(err, "ResourceNotFoundException"); err2 != nil {
			return fmt.Errorf("expected ResourceNotFoundException after delete, got: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "LookupTable_CreateValidation", func() error {
		_, err := client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String("bad name!"),
			TableBody:       aws.String("a\n1\n"),
		})
		if err2 := AssertErrorContains(err, "InvalidParameterException"); err2 != nil {
			return fmt.Errorf("invalid name: %v", err)
		}

		_, err = client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String("bad_csv"),
			TableBody:       aws.String("a,b\n\"unterminated,1\n"),
		})
		if err2 := AssertErrorContains(err, "InvalidParameterException"); err2 != nil {
			return fmt.Errorf("malformed CSV: %v", err)
		}

		name := uniqueTableName("dup_lookup")
		body := "a\n1\n"
		createResp, err := client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String(name),
			TableBody:       aws.String(body),
		})
		if err != nil {
			return fmt.Errorf("create: %v", err)
		}
		defer client.DeleteLookupTable(tc.ctx, &cloudwatchlogs.DeleteLookupTableInput{
			LookupTableArn: createResp.LookupTableArn,
		})
		_, err = client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String(name),
			TableBody:       aws.String(body),
		})
		if err2 := AssertErrorContains(err, "ResourceAlreadyExistsException"); err2 != nil {
			return fmt.Errorf("duplicate name: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "LookupTable_Quota", func() error {
		// The documented quota allows 100 lookup tables per account per
		// Region; the 101st creation is rejected.
		prefix := uniqueTableName("quota_l")
		var arns []string
		defer func() {
			for _, arn := range arns {
				client.DeleteLookupTable(tc.ctx, &cloudwatchlogs.DeleteLookupTableInput{
					LookupTableArn: aws.String(arn),
				})
			}
		}()
		existing, err := client.DescribeLookupTables(tc.ctx, &cloudwatchlogs.DescribeLookupTablesInput{})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		budget := 100 - len(existing.LookupTables)
		for i := 0; i < budget; i++ {
			resp, err := client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
				LookupTableName: aws.String(fmt.Sprintf("%s_%d", prefix, i)),
				TableBody:       aws.String("a\n1\n"),
			})
			if err != nil {
				return fmt.Errorf("fill table %d: %v", i, err)
			}
			arns = append(arns, aws.ToString(resp.LookupTableArn))
		}
		_, err = client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String(prefix + "_overflow"),
			TableBody:       aws.String("a\n1\n"),
		})
		if err2 := AssertErrorContains(err, "LimitExceededException"); err2 != nil {
			return fmt.Errorf("101st table: %v", err)
		}
		return nil
	}))

	// lookup enriches events with reference data matched on field values.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_LookupCommand", func() error {
		tableName := uniqueTableName("users_q")
		createResp, err := client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String(tableName),
			TableBody:       aws.String("id,name,department\nu1,Alice,Engineering\nu2,Bob,Support\n"),
		})
		if err != nil {
			return fmt.Errorf("create lookup table: %v", err)
		}
		defer client.DeleteLookupTable(tc.ctx, &cloudwatchlogs.DeleteLookupTableInput{
			LookupTableArn: createResp.LookupTableArn,
		})

		groupName := tc.uniquePrefix("lookup-query-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)
		if err := tc.createLogStream(groupName, "s1"); err != nil {
			return fmt.Errorf("create log stream: %v", err)
		}
		now := time.Now().UnixMilli()
		for i, msg := range []string{`{"user_id": "u1"}`, `{"user_id": "u3"}`} {
			if err := tc.putLogEvent(groupName, "s1", msg, now-int64(i)*1000); err != nil {
				return fmt.Errorf("put event: %v", err)
			}
		}

		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 60000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(fmt.Sprintf(`lookup %s id as user_id OUTPUT name, department | fields user_id, name, department | sort user_id asc`, tableName)),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		var all [][]types.ResultField
		for i := 0; i < 20; i++ {
			resResp, err := client.GetQueryResults(tc.ctx, &cloudwatchlogs.GetQueryResultsInput{
				QueryId: startResp.QueryId,
			})
			if err != nil {
				return fmt.Errorf("get results: %v", err)
			}
			if resResp.Status == types.QueryStatusComplete {
				all = resResp.Results
				break
			}
			if resResp.Status == types.QueryStatusFailed || resResp.Status == types.QueryStatusCancelled {
				return fmt.Errorf("query finished with status %s", resResp.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		if len(all) != 2 {
			return fmt.Errorf("rows = %d, want 2", len(all))
		}
		for _, row := range all {
			var uid, name, dept string
			for _, f := range row {
				switch aws.ToString(f.Field) {
				case "user_id":
					uid = aws.ToString(f.Value)
				case "name":
					name = aws.ToString(f.Value)
				case "department":
					dept = aws.ToString(f.Value)
				}
			}
			if uid == "u1" && (name != "Alice" || dept != "Engineering") {
				return fmt.Errorf("u1 enrichment: name=%q dept=%q", name, dept)
			}
			if uid == "u3" && name != "" {
				return fmt.Errorf("unmatched u3 should carry null name, got %q", name)
			}
		}
		return nil
	}))

	// cidrlookup enriches events by matching an IP field against CIDR
	// ranges in a table column.
	results = append(results, tc.runner.RunTest("logs", "StartQuery_CidrLookupCommand", func() error {
		tableName := uniqueTableName("nets_q")
		tableArn := ""
		createResp, err := client.CreateLookupTable(tc.ctx, &cloudwatchlogs.CreateLookupTableInput{
			LookupTableName: aws.String(tableName),
			TableBody:       aws.String("cidr,region,owner\n10.0.0.0/8,corp,it\n192.168.1.0/24,office,netops\n"),
		})
		if err != nil {
			return fmt.Errorf("create lookup table: %v", err)
		}
		tableArn = aws.ToString(createResp.LookupTableArn)
		defer client.DeleteLookupTable(tc.ctx, &cloudwatchlogs.DeleteLookupTableInput{
			LookupTableArn: aws.String(tableArn),
		})

		groupName := tc.uniquePrefix("cidr-query-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)
		if err := tc.createLogStream(groupName, "s1"); err != nil {
			return fmt.Errorf("create log stream: %v", err)
		}
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", `{"ip": "10.1.2.3"}`, now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}
		if err := tc.putLogEvent(groupName, "s1", `{"ip": "203.0.113.9"}`, now-1000); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		startResp, err := client.StartQuery(tc.ctx, &cloudwatchlogs.StartQueryInput{
			StartTime:     aws.Int64(now - 3600),
			EndTime:       aws.Int64(now + 60000),
			LogGroupNames: []string{groupName},
			QueryString:   aws.String(fmt.Sprintf(`cidrlookup %s ip as cidr OUTPUT region, owner | fields ip, region | sort ip desc`, tableName)),
		})
		if err != nil {
			return fmt.Errorf("start query: %v", err)
		}
		var all [][]types.ResultField
		for i := 0; i < 20; i++ {
			resResp, err := client.GetQueryResults(tc.ctx, &cloudwatchlogs.GetQueryResultsInput{
				QueryId: startResp.QueryId,
			})
			if err != nil {
				return fmt.Errorf("get results: %v", err)
			}
			if resResp.Status == types.QueryStatusComplete {
				all = resResp.Results
				break
			}
			if resResp.Status == types.QueryStatusFailed || resResp.Status == types.QueryStatusCancelled {
				return fmt.Errorf("query finished with status %s", resResp.Status)
			}
			time.Sleep(200 * time.Millisecond)
		}
		if len(all) != 2 {
			return fmt.Errorf("rows = %d, want 2", len(all))
		}
		for _, row := range all {
			var ip, region string
			for _, f := range row {
				switch aws.ToString(f.Field) {
				case "ip":
					ip = aws.ToString(f.Value)
				case "region":
					region = aws.ToString(f.Value)
				}
			}
			if strings.HasPrefix(ip, "10.") && region != "corp" {
				return fmt.Errorf("10.x enrichment: region=%q", region)
			}
			if strings.HasPrefix(ip, "203.") && region != "" {
				return fmt.Errorf("unmatched IP should carry null region, got %q", region)
			}
		}
		return nil
	}))

	// A scheduled query whose destination configuration violates the
	// documented member constraints is rejected at create time.
	results = append(results, tc.runner.RunTest("logs", "ScheduledQuery_DestinationValidation", func() error {
		base := func(dest *types.DestinationConfiguration) *cloudwatchlogs.CreateScheduledQueryInput {
			return &cloudwatchlogs.CreateScheduledQueryInput{
				Name:                     aws.String(tc.uniquePrefix("sched-dest")),
				QueryString:              aws.String("fields @message"),
				QueryLanguage:            types.QueryLanguageCwli,
				ExecutionRoleArn:         aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
				ScheduleExpression:       aws.String("rate(1 hour)"),
				State:                    types.ScheduledQueryStateDisabled,
				DestinationConfiguration: dest,
			}
		}
		role := aws.String("arn:aws:iam::123456789012:role/deliver")

		// The SDK enforces required members client-side, so the server-side
		// requirement check is exercised through a non-IAM role ARN instead.
		_, err := client.CreateScheduledQuery(tc.ctx, base(&types.DestinationConfiguration{
			LookupTableConfiguration: &types.LookupTableConfiguration{
				TableName: aws.String("bad_role"),
				RoleArn:   aws.String("arn:aws:s3:::not-a-role"),
			},
		}))
		if err2 := AssertErrorContains(err, "InvalidParameterException"); err2 != nil {
			return fmt.Errorf("non-IAM roleArn: %v", err)
		}

		_, err = client.CreateScheduledQuery(tc.ctx, base(&types.DestinationConfiguration{
			LookupTableConfiguration: &types.LookupTableConfiguration{
				TableName: aws.String("bad-name"),
				RoleArn:   role,
			},
		}))
		if err2 := AssertErrorContains(err, "InvalidParameterException"); err2 != nil {
			return fmt.Errorf("invalid tableName: %v", err)
		}

		_, err = client.CreateScheduledQuery(tc.ctx, base(&types.DestinationConfiguration{
			S3Configuration: &types.S3Configuration{
				DestinationIdentifier: aws.String("https://not-s3/prefix"),
				RoleArn:               role,
			},
		}))
		if err2 := AssertErrorContains(err, "InvalidParameterException"); err2 != nil {
			return fmt.Errorf("invalid s3 uri: %v", err)
		}

		createResp, err := client.CreateScheduledQuery(tc.ctx, base(&types.DestinationConfiguration{
			LookupTableConfiguration: &types.LookupTableConfiguration{
				TableName: aws.String(uniqueTableName("dest_valid")),
				RoleArn:   role,
			},
		}))
		if err != nil {
			return fmt.Errorf("valid destination rejected: %v", err)
		}
		_, err = client.DeleteScheduledQuery(tc.ctx, &cloudwatchlogs.DeleteScheduledQueryInput{
			Identifier: createResp.ScheduledQueryArn,
		})
		if err != nil {
			return fmt.Errorf("delete scheduled query: %v", err)
		}
		return nil
	}))

	// A scheduled query with a lookup table destination populates the table
	// with its result rows on the first scheduled execution and reports the
	// delivery through GetScheduledQueryHistory.
	results = append(results, tc.runner.RunTest("logs", "ScheduledQuery_LookupTableDestination", func() error {
		tableName := uniqueTableName("dest_refresh")
		groupName := tc.uniquePrefix("dest-group")
		if err := tc.createLogGroup(groupName); err != nil {
			return fmt.Errorf("create log group: %v", err)
		}
		defer tc.deleteLogGroup(groupName)
		if err := tc.createLogStream(groupName, "s1"); err != nil {
			return fmt.Errorf("create log stream: %v", err)
		}
		now := time.Now().UnixMilli()
		if err := tc.putLogEvent(groupName, "s1", `{"user_id":"u7","code":200}`, now); err != nil {
			return fmt.Errorf("put event: %v", err)
		}

		createResp, err := client.CreateScheduledQuery(tc.ctx, &cloudwatchlogs.CreateScheduledQueryInput{
			Name:                aws.String(tc.uniquePrefix("sched-lookup")),
			QueryString:         aws.String("fields user_id, code"),
			QueryLanguage:       types.QueryLanguageCwli,
			LogGroupIdentifiers: []string{groupName},
			ExecutionRoleArn:    aws.String("arn:aws:iam::123456789012:role/scheduled-query-role"),
			ScheduleExpression:  aws.String("rate(1 minute)"),
			State:               types.ScheduledQueryStateEnabled,
			DestinationConfiguration: &types.DestinationConfiguration{
				LookupTableConfiguration: &types.LookupTableConfiguration{
					TableName:   aws.String(tableName),
					RoleArn:     aws.String("arn:aws:iam::123456789012:role/deliver"),
					Description: aws.String("populated by scheduled query"),
				},
			},
		})
		if err != nil {
			return fmt.Errorf("create scheduled query: %v", err)
		}
		queryArn := aws.ToString(createResp.ScheduledQueryArn)
		defer client.DeleteScheduledQuery(tc.ctx, &cloudwatchlogs.DeleteScheduledQueryInput{
			Identifier: aws.String(queryArn),
		})
		tableArn := ""
		defer func() {
			if tableArn == "" {
				return
			}
			client.DeleteLookupTable(tc.ctx, &cloudwatchlogs.DeleteLookupTableInput{
				LookupTableArn: aws.String(tableArn),
			})
		}()

		// The AWS rate() contract runs the first query one full interval
		// after creation, so the first delivery arrives just past the
		// one-minute mark.
		deadline := time.Now().Add(150 * time.Second)
		var table *cloudwatchlogs.GetLookupTableOutput
		for time.Now().Before(deadline) {
			getResp, err := client.GetLookupTable(tc.ctx, &cloudwatchlogs.GetLookupTableInput{
				LookupTableArn: aws.String(tableName),
			})
			if err == nil {
				table = getResp
				break
			}
			time.Sleep(3 * time.Second)
		}
		if table == nil {
			return fmt.Errorf("lookup table %s was not populated by the scheduled query", tableName)
		}
		tableArn = aws.ToString(table.LookupTableArn)
		if !strings.Contains(aws.ToString(table.TableBody), "u7") || !strings.Contains(aws.ToString(table.TableBody), "200") {
			return fmt.Errorf("delivered body missing query results: %q", aws.ToString(table.TableBody))
		}
		if aws.ToString(table.Description) != "populated by scheduled query" {
			return fmt.Errorf("description = %q", aws.ToString(table.Description))
		}

		historyResp, err := client.GetScheduledQueryHistory(tc.ctx, &cloudwatchlogs.GetScheduledQueryHistoryInput{
			Identifier: aws.String(queryArn),
			StartTime:  aws.Int64(now - 60000),
			EndTime:    aws.Int64(time.Now().Add(time.Hour).UnixMilli()),
		})
		if err != nil {
			return fmt.Errorf("get scheduled query history: %v", err)
		}
		var delivered bool
		for _, trig := range historyResp.TriggerHistory {
			for _, dest := range trig.Destinations {
				if dest.DestinationType == types.ScheduledQueryDestinationTypeLookupTable &&
					dest.Status == types.ActionStatusComplete {
					delivered = true
				}
			}
		}
		if !delivered {
			return fmt.Errorf("history does not report a completed LOOKUP_TABLE delivery: %+v", historyResp.TriggerHistory)
		}
		return nil
	}))

	return results
}
