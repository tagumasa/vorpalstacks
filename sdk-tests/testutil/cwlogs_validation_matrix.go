package testutil

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// The C2 validation-matrix SDK pins: the member-trait rows that are
// reachable through the typed SDK (the enum members are client-side
// validated by the SDK and stay unit-pinned on the server side).
func (tc *cwlogsTestCtx) validationMatrixTests() []TestResult {
	var results []TestResult

	// The Tags map's length trait is 1-50 and the operation declares
	// TooManyTagsException: the 51st entry rejects with it.
	results = append(results, tc.runner.RunTest("logs", "TagResource_TooManyTags", func() error {
		tgName, cleanupGroup, err := tc.newLogGroupFixture("TooManyTagsGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		arn, err := tc.findLogGroupARN(tgName)
		if err != nil {
			return err
		}
		tags := make(map[string]string, 51)
		for i := 0; i < 51; i++ {
			tags[fmt.Sprintf("Key%02d", i)] = "v"
		}
		_, err = tc.client.TagResource(tc.ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: arn,
			Tags:        tags,
		})
		if err := AssertErrorContains(err, "TooManyTagsException"); err != nil {
			return fmt.Errorf("51 tags: %v", err)
		}
		// The 50-tag ceiling itself is accepted.
		delete(tags, "Key50")
		if _, err := tc.client.TagResource(tc.ctx, &cloudwatchlogs.TagResourceInput{
			ResourceArn: arn,
			Tags:        tags,
		}); err != nil {
			return fmt.Errorf("50 tags: %v", err)
		}
		return nil
	}))

	// metricTransformations carries @length 1..1 — exactly one
	// transformation per metric filter.
	results = append(results, tc.runner.RunTest("logs", "PutMetricFilter_TwoTransformations_Rejected", func() error {
		mfName, cleanupGroup, err := tc.newLogGroupFixture("TwoXformGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		_, err = tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:  aws.String(mfName),
			FilterName:    aws.String("two"),
			FilterPattern: aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{
				{MetricName: aws.String("One"), MetricNamespace: aws.String("App"), MetricValue: aws.String("1")},
				{MetricName: aws.String("Two"), MetricNamespace: aws.String("App"), MetricValue: aws.String("1")},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("two transformations: %v", err)
		}
		return nil
	}))

	// logGroupNamePattern is a case-sensitive substring match; the
	// response is restricted to arn, creationTime and logGroupName while
	// it is in use, and it is mutually exclusive with logGroupNamePrefix.
	results = append(results, tc.runner.RunTest("logs", "DescribeLogGroups_Pattern_SubstringAndRestrictedResponse", func() error {
		base := tc.uniquePrefix("pattern")
		groups := []string{base + "DataLogs", base + "aws/DataLogs", base + "GroupDataLogs", base + "datalogs", base + "Groupdata"}
		for _, g := range groups {
			if err := tc.createLogGroup(g); err != nil {
				return fmt.Errorf("create %s: %v", g, err)
			}
			defer tc.deleteLogGroup(g)
		}

		// The pattern is an unanchored substring, so residue groups from
		// a crashed earlier run of this test match it just as
		// legitimately as the fixtures, and the fixtures cannot be
		// assumed to sit on the first page: the pinned contract is that
		// every returned name contains the substring (a case-insensitive
		// matcher returns the lowercase fixtures and fails here), the
		// restricted member set holds for every returned group, and all
		// three fixtures appear across the full traversal.
		want := map[string]bool{groups[0]: true, groups[1]: true, groups[2]: true}
		matched := 0
		seen := 0
		var nextToken *string
		for {
			resp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
				LogGroupNamePattern: aws.String("DataLogs"),
				NextToken:           nextToken,
			})
			if err != nil {
				return fmt.Errorf("describe with pattern: %v", err)
			}
			for _, lg := range resp.LogGroups {
				name := aws.ToString(lg.LogGroupName)
				if !strings.Contains(name, "DataLogs") {
					return fmt.Errorf("pattern response carried group %q, which does not contain the substring", name)
				}
				if lg.Arn == nil || lg.CreationTime == nil {
					return fmt.Errorf("restricted response must carry logGroupName, arn and creationTime")
				}
				if lg.StoredBytes != nil || lg.MetricFilterCount != nil {
					return fmt.Errorf("restricted response must not carry storedBytes or metricFilterCount")
				}
				seen++
				if want[name] {
					matched++
				}
			}
			if resp.NextToken == nil || *resp.NextToken == "" {
				break
			}
			nextToken = resp.NextToken
		}
		if matched != 3 {
			return fmt.Errorf("substring matches: want the 3 case-sensitive fixtures, got %d (%d groups returned)", matched, seen)
		}

		_, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePattern: aws.String(base),
			LogGroupNamePrefix:  aws.String(base),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("pattern+prefix: %v", err)
		}
		return nil
	}))

	// The documented PutLogEvents byte limits: each event no larger than
	// 1 MB, and the batch — the sum of all messages in UTF-8 plus 26
	// bytes per event — no larger than 1,048,576 bytes.
	results = append(results, tc.runner.RunTest("logs", "PutLogEvents_EventAndBatchSizeLimits", func() error {
		evName, cleanupGroup, err := tc.newGroupStreamFixture("SizeLimitGroup", "s")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		_, err = tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(evName),
			LogStreamName: aws.String("s"),
			LogEvents: []types.InputLogEvent{{
				Message:   aws.String(strings.Repeat("x", 1048577)),
				Timestamp: aws.Int64(time.Now().UnixMilli()),
			}},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("event over 1 MB: %v", err)
		}

		_, err = tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
			LogGroupName:  aws.String(evName),
			LogStreamName: aws.String("s"),
			LogEvents: []types.InputLogEvent{
				{Message: aws.String(strings.Repeat("a", 600000)), Timestamp: aws.Int64(time.Now().UnixMilli())},
				{Message: aws.String(strings.Repeat("b", 600000)), Timestamp: aws.Int64(time.Now().UnixMilli() + 1)},
			},
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("batch over 1 MB: %v", err)
		}
		return nil
	}))

	// The documented GetLogEvents default page budget: as many events as
	// fit in a response of 1 MB (up to 10,000) — pagination still serves
	// the whole stream.
	results = append(results, tc.runner.RunTest("logs", "GetLogEvents_ResponseSizeCap", func() error {
		rcName, cleanupGroup, err := tc.newGroupStreamFixture("ResponseCapGroup", "s")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		ts := time.Now().UnixMilli()
		// Four 270,000-byte events exceed the 1 MB response budget in
		// total while each two-event write stays inside the batch
		// ceiling.
		for i := 0; i < 2; i++ {
			batch := make([]types.InputLogEvent, 2)
			for j := range batch {
				batch[j] = types.InputLogEvent{
					Message:   aws.String(strings.Repeat("m", 270000)),
					Timestamp: aws.Int64(ts + int64(i*2+j)),
				}
			}
			if _, err := tc.client.PutLogEvents(tc.ctx, &cloudwatchlogs.PutLogEventsInput{
				LogGroupName: aws.String(rcName), LogStreamName: aws.String("s"), LogEvents: batch,
			}); err != nil {
				return fmt.Errorf("seed batch: %v", err)
			}
		}

		first, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
			LogGroupName: aws.String(rcName), LogStreamName: aws.String("s"), StartFromHead: aws.Bool(true),
		})
		if err != nil {
			return fmt.Errorf("first page: %v", err)
		}
		if len(first.Events) == 0 || len(first.Events) >= 4 {
			return fmt.Errorf("first page should hold a strict subset of the four events, got %d", len(first.Events))
		}

		served := len(first.Events)
		token := first.NextForwardToken
		for token != nil && *token != "" {
			page, err := tc.client.GetLogEvents(tc.ctx, &cloudwatchlogs.GetLogEventsInput{
				LogGroupName: aws.String(rcName), LogStreamName: aws.String("s"),
				StartFromHead: aws.Bool(true), NextToken: token,
			})
			if err != nil {
				return fmt.Errorf("page: %v", err)
			}
			if page.NextForwardToken == nil || *page.NextForwardToken == *token {
				served += len(page.Events)
				break
			}
			served += len(page.Events)
			token = page.NextForwardToken
		}
		if served != 4 {
			return fmt.Errorf("pagination should serve all four events, served %d", served)
		}
		return nil
	}))

	// The resource-scoped revision identity: "The expected revision ID of
	// the resource policy. Required when resourceArn is provided... Use
	// null when creating a resource policy for the first time" and the
	// response's revisionId is "Only returned for resource-scoped
	// policies".
	results = append(results, tc.runner.RunTest("logs", "PutResourcePolicy_RevisionId", func() error {
		policyName := tc.uniquePrefix("rev-policy")
		groupName, cleanupGroup, err := tc.newLogGroupFixture("RevPolicyGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		// A resource-scoped delete must quote the current revision; the
		// deferred cleanup re-reads the last issued revision at run time.
		var lastRevision *string
		defer func() {
			del := &cloudwatchlogs.DeleteResourcePolicyInput{PolicyName: aws.String(policyName)}
			if lastRevision != nil {
				del.ExpectedRevisionId = lastRevision
			}
			_, _ = tc.client.DeleteResourcePolicy(tc.ctx, del)
		}()
		groupArn, err := tc.findLogGroupARN(groupName)
		if err != nil {
			return fmt.Errorf("find ARN: %v", err)
		}
		doc := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"route53.amazonaws.com"},"Action":"logs:PutLogEvents","Resource":"*"}]}`

		// The account-wide form returns no revisionId at all.
		wide, err := tc.client.PutResourcePolicy(tc.ctx, &cloudwatchlogs.PutResourcePolicyInput{
			PolicyName: aws.String(policyName + "-wide"), PolicyDocument: aws.String(doc),
		})
		defer tc.client.DeleteResourcePolicy(tc.ctx, &cloudwatchlogs.DeleteResourcePolicyInput{PolicyName: aws.String(policyName + "-wide")})
		if err != nil {
			return fmt.Errorf("account-wide put: %v", err)
		}
		if wide.RevisionId != nil {
			return fmt.Errorf("account-wide put returned a revisionId, which only resource-scoped policies carry")
		}

		// The first-time null creates and issues a revision; an update
		// must quote it; a stale quote rejects.
		first, err := tc.client.PutResourcePolicy(tc.ctx, &cloudwatchlogs.PutResourcePolicyInput{
			PolicyName: aws.String(policyName), PolicyDocument: aws.String(doc), ResourceArn: groupArn,
		})
		if err != nil {
			return fmt.Errorf("first put: %v", err)
		}
		if first.RevisionId == nil || *first.RevisionId == "" {
			return fmt.Errorf("first put issued no revisionId")
		}
		if _, noQuoteErr := tc.client.PutResourcePolicy(tc.ctx, &cloudwatchlogs.PutResourcePolicyInput{
			PolicyName: aws.String(policyName), PolicyDocument: aws.String(doc), ResourceArn: groupArn,
		}); noQuoteErr == nil {
			return fmt.Errorf("update without a revision quote must reject")
		}
		_, err = tc.client.PutResourcePolicy(tc.ctx, &cloudwatchlogs.PutResourcePolicyInput{
			PolicyName: aws.String(policyName), PolicyDocument: aws.String(doc), ResourceArn: groupArn, ExpectedRevisionId: aws.String("stale"),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("stale revision: %v", err)
		}
		second, err := tc.client.PutResourcePolicy(tc.ctx, &cloudwatchlogs.PutResourcePolicyInput{
			PolicyName: aws.String(policyName), PolicyDocument: aws.String(doc), ResourceArn: groupArn, ExpectedRevisionId: first.RevisionId,
		})
		if err != nil {
			return fmt.Errorf("current revision: %v", err)
		}
		if second.RevisionId == nil || *second.RevisionId == *first.RevisionId {
			return fmt.Errorf("update must mint a fresh revisionId")
		}
		lastRevision = second.RevisionId
		// "One per LogGroup resourceARN": a second policy on the same ARN
		// rejects with the quota identity. The no-retry client carries the
		// call: LimitExceededException is in the SDK retryer's throttling
		// class, so the standard client would wait out its retry ladder
		// before this assertion runs.
		_, err = tc.noRetryClient.PutResourcePolicy(tc.ctx, &cloudwatchlogs.PutResourcePolicyInput{
			PolicyName: aws.String(policyName + "-second"), PolicyDocument: aws.String(doc), ResourceArn: groupArn,
		})
		if err := AssertErrorContains(err, "LimitExceededException"); err != nil {
			return fmt.Errorf("second policy on one resourceArn: %v", err)
		}

		// The delete contract: a resource-scoped delete without a revision
		// quote rejects, and a stale quote rejects; the deferred cleanup
		// holds the current quote.
		_, err = tc.client.DeleteResourcePolicy(tc.ctx, &cloudwatchlogs.DeleteResourcePolicyInput{PolicyName: aws.String(policyName)})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("resource-scoped delete without a revision quote: %v", err)
		}
		_, err = tc.client.DeleteResourcePolicy(tc.ctx, &cloudwatchlogs.DeleteResourcePolicyInput{
			PolicyName: aws.String(policyName), ExpectedRevisionId: aws.String("stale"),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("resource-scoped delete with a stale quote: %v", err)
		}

		// The resourceArn form: "Currently only supports LogGroup ARN."
		_, err = tc.client.PutResourcePolicy(tc.ctx, &cloudwatchlogs.PutResourcePolicyInput{
			PolicyName: aws.String(policyName + "-badarn"), PolicyDocument: aws.String(doc),
			ResourceArn: aws.String(fmt.Sprintf("arn:aws:logs:%s:%s:destination:not-a-group", tc.region, tc.runner.AccountID())),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("non-log-group resourceArn: %v", err)
		}

		// The listing scope: absent defaults to ACCOUNT, RESOURCE selects
		// the resource-scoped form, and a foreign value rejects.
		defaultScope, err := tc.client.DescribeResourcePolicies(tc.ctx, &cloudwatchlogs.DescribeResourcePoliciesInput{})
		if err != nil {
			return fmt.Errorf("default-scope describe: %v", err)
		}
		for _, p := range defaultScope.ResourcePolicies {
			if p.PolicyScope == types.PolicyScopeResource {
				return fmt.Errorf("default scope listed a RESOURCE policy named %s", aws.ToString(p.PolicyName))
			}
		}
		scopedList, err := tc.client.DescribeResourcePolicies(tc.ctx, &cloudwatchlogs.DescribeResourcePoliciesInput{PolicyScope: types.PolicyScopeResource})
		if err != nil {
			return fmt.Errorf("RESOURCE-scope describe: %v", err)
		}
		if len(scopedList.ResourcePolicies) != 1 || aws.ToString(scopedList.ResourcePolicies[0].PolicyName) != policyName {
			return fmt.Errorf("RESOURCE scope returned %d policies, want the scoped one alone", len(scopedList.ResourcePolicies))
		}
		_, err = tc.client.DescribeResourcePolicies(tc.ctx, &cloudwatchlogs.DescribeResourcePoliciesInput{PolicyScope: types.PolicyScope("BOTH")})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("foreign policyScope: %v", err)
		}
		return nil
	}))

	// AssociateKmsKey's account-level query-result ARN form is accepted
	// and recorded; the target members are mutually exclusive.
	results = append(results, tc.runner.RunTest("logs", "AssociateKmsKey_QueryResultArn", func() error {
		kmsName, cleanupGroup, err := tc.newLogGroupFixture("QueryResultKmsGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		queryResultArn := fmt.Sprintf("arn:aws:logs:%s:%s:query-result:*", tc.region, tc.runner.AccountID())
		keyArn := fmt.Sprintf("arn:aws:kms:%s:%s:key/1234abcd-12ab-34cd-56ef-1234567890ab", tc.region, tc.runner.AccountID())
		if _, err := tc.client.AssociateKmsKey(tc.ctx, &cloudwatchlogs.AssociateKmsKeyInput{
			ResourceIdentifier: aws.String(queryResultArn), KmsKeyId: aws.String(keyArn),
		}); err != nil {
			return fmt.Errorf("query-result association: %v", err)
		}
		if _, err := tc.client.DisassociateKmsKey(tc.ctx, &cloudwatchlogs.DisassociateKmsKeyInput{
			ResourceIdentifier: aws.String(queryResultArn),
		}); err != nil {
			return fmt.Errorf("query-result disassociation: %v", err)
		}
		_, err = tc.client.AssociateKmsKey(tc.ctx, &cloudwatchlogs.AssociateKmsKeyInput{
			LogGroupName: aws.String(kmsName), ResourceIdentifier: aws.String(queryResultArn), KmsKeyId: aws.String(keyArn),
		})
		if err := AssertErrorContains(err, "InvalidParameterException"); err != nil {
			return fmt.Errorf("both target members: %v", err)
		}
		return nil
	}))

	// The per-log-group metric filter quota is 100 (CloudWatch Logs
	// quotas page, not adjustable): the hundred-and-first distinct filter
	// name rejects with LimitExceededException — the quota-shaped error
	// PutMetricFilter declares — while an update of an existing name at
	// the cap still succeeds.
	results = append(results, tc.runner.RunTest("logs", "PutMetricFilter_PerGroupQuota_100", func() error {
		qName, cleanupGroup, err := tc.newLogGroupFixture("QuotaGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		xform := types.MetricTransformation{
			MetricName: aws.String("QuotaCount"), MetricNamespace: aws.String("Quota"), MetricValue: aws.String("1"),
		}
		for i := 0; i < 100; i++ {
			if _, err := tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
				LogGroupName:          aws.String(qName),
				FilterName:            aws.String(fmt.Sprintf("f-%03d", i)),
				FilterPattern:         aws.String("ERROR"),
				MetricTransformations: []types.MetricTransformation{xform},
			}); err != nil {
				return fmt.Errorf("filter %d within the quota: %v", i, err)
			}
		}

		// The 101st filter rejects with the quota identity; the no-retry
		// client carries the call so the SDK's throttling-class retry
		// ladder adds no wait before the assertion.
		_, err = tc.noRetryClient.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:          aws.String(qName),
			FilterName:            aws.String("f-over"),
			FilterPattern:         aws.String("ERROR"),
			MetricTransformations: []types.MetricTransformation{xform},
		})
		if err := AssertErrorContains(err, "LimitExceededException"); err != nil {
			return fmt.Errorf("filter beyond the quota: %v", err)
		}

		if _, err := tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:          aws.String(qName),
			FilterName:            aws.String("f-000"),
			FilterPattern:         aws.String("WARN"),
			MetricTransformations: []types.MetricTransformation{xform},
		}); err != nil {
			return fmt.Errorf("update at the cap: %v", err)
		}
		return nil
	}))

	// An update keeps the stored creation time: DescribeMetricFilters
	// reports the filter's original creation, not the last update.
	results = append(results, tc.runner.RunTest("logs", "PutMetricFilter_CreationTimePreservedOnUpdate", func() error {
		tsName, cleanupGroup, err := tc.newLogGroupFixture("CreationTimeGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		xform := types.MetricTransformation{
			MetricName: aws.String("TsCount"), MetricNamespace: aws.String("Ts"), MetricValue: aws.String("1"),
		}
		put := func(name, pattern string) error {
			_, err := tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
				LogGroupName:          aws.String(tsName),
				FilterName:            aws.String(name),
				FilterPattern:         aws.String(pattern),
				MetricTransformations: []types.MetricTransformation{xform},
			})
			return err
		}
		if err := put("keeper", "ERROR"); err != nil {
			return fmt.Errorf("first put: %v", err)
		}
		first, err := tc.collectAllMetricFilters(tsName, "keeper")
		if err != nil {
			return fmt.Errorf("describe after create: %v", err)
		}
		if len(first) != 1 || first[0].CreationTime == nil {
			return fmt.Errorf("expected one filter with a creationTime after create")
		}
		created := *first[0].CreationTime

		time.Sleep(20 * time.Millisecond)
		if err := put("keeper", "WARN"); err != nil {
			return fmt.Errorf("second put: %v", err)
		}
		second, err := tc.collectAllMetricFilters(tsName, "keeper")
		if err != nil {
			return fmt.Errorf("describe after update: %v", err)
		}
		if len(second) != 1 {
			return fmt.Errorf("expected one filter after update")
		}
		if *second[0].CreationTime != created {
			return fmt.Errorf("update must preserve creationTime: first %d, second %d",
				created, *second[0].CreationTime)
		}
		if second[0].FilterPattern == nil || *second[0].FilterPattern != "WARN" {
			return fmt.Errorf("update must still change the pattern")
		}
		return nil
	}))

	// TestMetricFilter's eventNumber is the matched event's position in
	// logEventMessages, not a match ordinal.
	results = append(results, tc.runner.RunTest("logs", "TestMetricFilter_EventNumberIsInputPosition", func() error {
		resp, err := tc.client.TestMetricFilter(tc.ctx, &cloudwatchlogs.TestMetricFilterInput{
			FilterPattern: aws.String("ERROR"),
			LogEventMessages: []string{
				"nothing here",
				"ERROR at two",
				"still nothing",
				"warn only",
				"ERROR at five",
			},
		})
		if err != nil {
			return fmt.Errorf("test metric filter: %v", err)
		}
		if len(resp.Matches) != 2 {
			return fmt.Errorf("two messages match, got %d records", len(resp.Matches))
		}
		for i, want := range []int64{2, 5} {
			if resp.Matches[i].EventNumber != want {
				return fmt.Errorf("match %d must report input position %d, got %d",
					i, want, resp.Matches[i].EventNumber)
			}
		}
		return nil
	}))

	// DeleteLogGroup tears down the group's data protection policy with
	// it: the policy record dies with the group, and a same-named
	// recreation starts policy-free.
	results = append(results, tc.runner.RunTest("logs", "DeleteLogGroup_TearsDownDataProtectionPolicy", func() error {
		dppName := tc.uniquePrefix("DppTeardown")
		if err := tc.createLogGroup(dppName); err != nil {
			return fmt.Errorf("create group: %v", err)
		}
		if _, err := tc.client.PutDataProtectionPolicy(tc.ctx, &cloudwatchlogs.PutDataProtectionPolicyInput{
			LogGroupIdentifier: aws.String(dppName),
			PolicyDocument: aws.String(`{"Name":"data-protection-policy","Version":"2021-06","Statement":[` +
				`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Audit":{"FindingsDestination":{}}}},` +
				`{"DataIdentifer":["arn:aws:dataprotection::aws:DataIdentifier/EmailAddress"],"Operation":{"Deidentify":{"MaskConfig":{}}}}]}`),
		}); err != nil {
			return fmt.Errorf("put policy: %v", err)
		}
		if _, err := tc.client.DeleteLogGroup(tc.ctx, &cloudwatchlogs.DeleteLogGroupInput{
			LogGroupName: aws.String(dppName),
		}); err != nil {
			return fmt.Errorf("delete group: %v", err)
		}
		var rnf *types.ResourceNotFoundException
		if _, err := tc.client.GetDataProtectionPolicy(tc.ctx, &cloudwatchlogs.GetDataProtectionPolicyInput{
			LogGroupIdentifier: aws.String(dppName),
		}); err == nil {
			return fmt.Errorf("the policy must die with the group")
		} else if !errors.As(err, &rnf) {
			return fmt.Errorf("a policy for a deleted group must answer ResourceNotFoundException, got %v", err)
		}

		if err := tc.createLogGroup(dppName); err != nil {
			return fmt.Errorf("recreate group: %v", err)
		}
		defer tc.deleteLogGroup(dppName)
		if _, err := tc.client.GetDataProtectionPolicy(tc.ctx, &cloudwatchlogs.GetDataProtectionPolicyInput{
			LogGroupIdentifier: aws.String(dppName),
		}); err == nil {
			return fmt.Errorf("a recreated group must start policy-free")
		} else if !errors.As(err, &rnf) {
			return fmt.Errorf("a recreated group's absent policy must answer ResourceNotFoundException, got %v", err)
		}
		return nil
	}))

	// Deleting a log stream releases its bytes from the group's
	// storedBytes on the same basis ingestion added them.
	results = append(results, tc.runner.RunTest("logs", "DeleteLogStream_ReleasesStoredBytes", func() error {
		group, cleanupGroup, err := tc.newLogGroupFixture("StoredBytes")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		for _, stream := range []string{"sa", "sb"} {
			if err := tc.createLogStream(group, stream); err != nil {
				return fmt.Errorf("create stream %s: %v", stream, err)
			}
		}
		msgA := strings.Repeat("a", 300)
		msgB := strings.Repeat("b", 170)
		ts := time.Now().UnixMilli()
		if err := tc.putLogEvent(group, "sa", msgA, ts); err != nil {
			return fmt.Errorf("put stream sa: %v", err)
		}
		if err := tc.putLogEvent(group, "sb", msgB, ts); err != nil {
			return fmt.Errorf("put stream sb: %v", err)
		}
		storedBytes := func() (int64, error) {
			resp, err := tc.client.DescribeLogGroups(tc.ctx, &cloudwatchlogs.DescribeLogGroupsInput{
				LogGroupNamePrefix: aws.String(group),
			})
			if err != nil {
				return 0, err
			}
			for _, lg := range resp.LogGroups {
				if aws.ToString(lg.LogGroupName) == group {
					return aws.ToInt64(lg.StoredBytes), nil
				}
			}
			return 0, fmt.Errorf("group %s not found in describe", group)
		}
		if got, err := storedBytes(); err != nil {
			return fmt.Errorf("describe after ingestion: %v", err)
		} else if want := int64(len(msgA) + len(msgB)); got != want {
			return fmt.Errorf("storedBytes after ingestion = %d, want the ingested message bytes %d", got, want)
		}
		if _, err := tc.client.DeleteLogStream(tc.ctx, &cloudwatchlogs.DeleteLogStreamInput{
			LogGroupName:  aws.String(group),
			LogStreamName: aws.String("sa"),
		}); err != nil {
			return fmt.Errorf("delete stream sa: %v", err)
		}
		if got, err := storedBytes(); err != nil {
			return fmt.Errorf("describe after stream delete: %v", err)
		} else if want := int64(len(msgB)); got != want {
			return fmt.Errorf("storedBytes after deleting stream sa = %d, want the surviving stream's bytes %d", got, want)
		}
		return nil
	}))

	// The transform-era metric-filter members survive the full
	// Put→Describe round-trip (persistence carries them).
	results = append(results, tc.runner.RunTest("logs", "PutMetricFilter_TransformSelectionMembers_Roundtrip", func() error {
		mfName, cleanupGroup, err := tc.newLogGroupFixture("MfSelectionGroup")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		_, err = tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:              aws.String(mfName),
			FilterName:                aws.String("select"),
			FilterPattern:             aws.String("{ $.latency > 0 }"),
			MetricTransformations:     []types.MetricTransformation{{MetricName: aws.String("Lat"), MetricNamespace: aws.String("Sel"), MetricValue: aws.String("$.latency")}},
			ApplyOnTransformedLogs:    true,
			FieldSelectionCriteria:    aws.String(`@aws.region = "us-east-1" OR @aws.account IN ["000000000000"]`),
			EmitSystemFieldDimensions: []string{"@aws.account", "@aws.region"},
		})
		if err != nil {
			return fmt.Errorf("put: %v", err)
		}
		resp, err := tc.client.DescribeMetricFilters(tc.ctx, &cloudwatchlogs.DescribeMetricFiltersInput{
			LogGroupName: aws.String(mfName), FilterNamePrefix: aws.String("select"),
		})
		if err != nil {
			return fmt.Errorf("describe: %v", err)
		}
		if len(resp.MetricFilters) != 1 {
			return fmt.Errorf("expected one filter, got %d", len(resp.MetricFilters))
		}
		f := resp.MetricFilters[0]
		if f.ApplyOnTransformedLogs != true {
			return fmt.Errorf("applyOnTransformedLogs must survive the round-trip, got %t", f.ApplyOnTransformedLogs)
		}
		if aws.ToString(f.FieldSelectionCriteria) != `@aws.region = "us-east-1" OR @aws.account IN ["000000000000"]` {
			return fmt.Errorf("fieldSelectionCriteria must survive the round-trip, got %q", aws.ToString(f.FieldSelectionCriteria))
		}
		if len(f.EmitSystemFieldDimensions) != 2 || f.EmitSystemFieldDimensions[0] != "@aws.account" || f.EmitSystemFieldDimensions[1] != "@aws.region" {
			return fmt.Errorf("emitSystemFieldDimensions must survive the round-trip, got %v", f.EmitSystemFieldDimensions)
		}
		// The member's value set is documented: "Valid values are
		// @aws.account and @aws.region" — a foreign value rejects.
		if _, err := tc.client.PutMetricFilter(tc.ctx, &cloudwatchlogs.PutMetricFilterInput{
			LogGroupName:              aws.String(mfName),
			FilterName:                aws.String("select-bad"),
			FilterPattern:             aws.String("ERROR"),
			MetricTransformations:     []types.MetricTransformation{{MetricName: aws.String("Lat2"), MetricNamespace: aws.String("Sel"), MetricValue: aws.String("1")}},
			EmitSystemFieldDimensions: []string{"$.host"},
		}); err == nil {
			return fmt.Errorf("a foreign emitSystemFieldDimensions value must reject")
		}
		return nil
	}))

	return results
}
