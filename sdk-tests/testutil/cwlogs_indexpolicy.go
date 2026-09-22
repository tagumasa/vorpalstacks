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

// The field index pins (E3): the log-group-level policy lifecycle with
// its precedence over the account-level FIELD_INDEX_POLICY, the
// DescribeFieldIndexes derivation (DEFAULT listing, CUSTOM fields with
// event-derived bounds, the INACTIVE trail) and the documented
// validation rows.

// collectAllFieldIndexes pages DescribeFieldIndexes to exhaustion.
func (tc *cwlogsTestCtx) collectAllFieldIndexes(groups []string, categories []types.IndexCategory) ([]types.FieldIndex, error) {
	var all []types.FieldIndex
	token := ""
	for {
		resp, err := tc.client.DescribeFieldIndexes(tc.ctx, &cloudwatchlogs.DescribeFieldIndexesInput{
			LogGroupIdentifiers: groups,
			IndexCategories:     categories,
			NextToken:           condString(token),
		})
		if err != nil {
			return nil, err
		}
		all = append(all, resp.FieldIndexes...)
		if resp.NextToken == nil || aws.ToString(resp.NextToken) == "" {
			return all, nil
		}
		token = aws.ToString(resp.NextToken)
	}
}

func condString(v string) *string {
	if v == "" {
		return nil
	}
	return aws.String(v)
}

// fieldIndexRow finds one (name, category) row of a group's listing.
func fieldIndexRow(rows []types.FieldIndex, name string, category types.IndexCategory) *types.FieldIndex {
	for i := range rows {
		if aws.ToString(rows[i].FieldIndexName) == name && rows[i].IndexCategory == category {
			return &rows[i]
		}
	}
	return nil
}

func (tc *cwlogsTestCtx) indexPolicyTests() []TestResult {
	var results []TestResult

	results = append(results, tc.runner.RunTest("logs", "IndexPolicy_PutDescribeDeleteLifecycle", func() error {
		group, cleanupGroup, err := tc.newLogGroupFixture("idx")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		put, err := tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"Fields":["RequestId"],"FieldsV2":{"APIName":{"type":"FACET"}}}`),
		})
		if err != nil {
			return fmt.Errorf("put index policy: %v", err)
		}
		policy := put.IndexPolicy
		if policy == nil {
			return fmt.Errorf("put response carries no indexPolicy")
		}
		if policy.Source != types.IndexSourceLogGroup {
			return fmt.Errorf("put source: %v", policy.Source)
		}
		if policy.PolicyName != nil {
			return fmt.Errorf("group-level policy carries policyName %q", aws.ToString(policy.PolicyName))
		}
		if policy.LastUpdateTime == nil || !strings.Contains(aws.ToString(policy.LogGroupIdentifier), group) {
			return fmt.Errorf("put echo: %+v", policy)
		}

		desc, err := tc.client.DescribeIndexPolicies(tc.ctx, &cloudwatchlogs.DescribeIndexPoliciesInput{
			LogGroupIdentifiers: []string{group},
		})
		if err != nil {
			return fmt.Errorf("describe index policies: %v", err)
		}
		if len(desc.IndexPolicies) != 1 || desc.IndexPolicies[0].Source != types.IndexSourceLogGroup {
			return fmt.Errorf("describe rows: %+v", desc.IndexPolicies)
		}
		if aws.ToString(desc.IndexPolicies[0].PolicyDocument) != aws.ToString(policy.PolicyDocument) {
			return fmt.Errorf("document round trip: %v", desc.IndexPolicies[0].PolicyDocument)
		}

		// The listing reports the facet configuration's type.
		rows, err := tc.collectAllFieldIndexes([]string{group}, []types.IndexCategory{types.IndexCategoryCustom})
		if err != nil {
			return fmt.Errorf("describe field indexes: %v", err)
		}
		if row := fieldIndexRow(rows, "RequestId", types.IndexCategoryCustom); row == nil || row.Type != types.IndexTypeFieldIndex {
			return fmt.Errorf("RequestId row: %+v", row)
		}
		if row := fieldIndexRow(rows, "APIName", types.IndexCategoryCustom); row == nil || row.Type != types.IndexTypeFacet {
			return fmt.Errorf("APIName row: %+v", row)
		}

		if _, err := tc.client.DeleteIndexPolicy(tc.ctx, &cloudwatchlogs.DeleteIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
		}); err != nil {
			return fmt.Errorf("delete index policy: %v", err)
		}
		desc, err = tc.client.DescribeIndexPolicies(tc.ctx, &cloudwatchlogs.DescribeIndexPoliciesInput{
			LogGroupIdentifiers: []string{group},
		})
		if err != nil {
			return fmt.Errorf("describe after delete: %v", err)
		}
		if len(desc.IndexPolicies) != 0 {
			return fmt.Errorf("policy survived deletion: %+v", desc.IndexPolicies)
		}
		var rnf *types.ResourceNotFoundException
		_, err = tc.client.DeleteIndexPolicy(tc.ctx, &cloudwatchlogs.DeleteIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
		})
		if !errors.As(err, &rnf) {
			return fmt.Errorf("second delete: %v", err)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "IndexPolicy_AccountPolicyPrecedence", func() error {
		prefix := tc.uniquePrefix("idxacct")
		group := prefix + "-app"
		policyName := prefix + "-policy"
		if err := tc.createLogGroup(group); err != nil {
			return err
		}
		defer tc.deleteLogGroup(group)
		defer tc.client.DeleteAccountPolicy(tc.ctx, &cloudwatchlogs.DeleteAccountPolicyInput{
			PolicyName: aws.String(policyName), PolicyType: types.PolicyTypeFieldIndexPolicy,
		})

		if _, err := tc.client.PutAccountPolicy(tc.ctx, &cloudwatchlogs.PutAccountPolicyInput{
			PolicyName:        aws.String(policyName),
			PolicyDocument:    aws.String(`{"Fields":["SessionId"]}`),
			PolicyType:        types.PolicyTypeFieldIndexPolicy,
			SelectionCriteria: aws.String("LogGroupNamePrefix = " + prefix),
		}); err != nil {
			return fmt.Errorf("put account policy: %v", err)
		}

		// The account policy applies while the group carries no
		// group-level policy.
		desc, err := tc.client.DescribeIndexPolicies(tc.ctx, &cloudwatchlogs.DescribeIndexPoliciesInput{
			LogGroupIdentifiers: []string{group},
		})
		if err != nil {
			return fmt.Errorf("describe account policy: %v", err)
		}
		if len(desc.IndexPolicies) != 1 {
			return fmt.Errorf("account policy rows: %+v", desc.IndexPolicies)
		}
		accountRow := desc.IndexPolicies[0]
		if accountRow.Source != types.IndexSourceAccount || aws.ToString(accountRow.PolicyName) != policyName {
			return fmt.Errorf("account policy row: %+v", accountRow)
		}

		// The group-level policy overrides it.
		if _, err := tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"Fields":["RequestId"]}`),
		}); err != nil {
			return fmt.Errorf("put group policy: %v", err)
		}
		desc, err = tc.client.DescribeIndexPolicies(tc.ctx, &cloudwatchlogs.DescribeIndexPoliciesInput{
			LogGroupIdentifiers: []string{group},
		})
		if err != nil {
			return fmt.Errorf("describe group policy: %v", err)
		}
		if len(desc.IndexPolicies) != 1 || desc.IndexPolicies[0].Source != types.IndexSourceLogGroup {
			return fmt.Errorf("group policy does not override: %+v", desc.IndexPolicies)
		}

		// Deleting the group policy falls back to the account policy,
		// whose fields the group then indexes.
		if _, err := tc.client.DeleteIndexPolicy(tc.ctx, &cloudwatchlogs.DeleteIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
		}); err != nil {
			return fmt.Errorf("delete group policy: %v", err)
		}
		rows, err := tc.collectAllFieldIndexes([]string{group}, []types.IndexCategory{types.IndexCategoryCustom})
		if err != nil {
			return fmt.Errorf("field indexes after fallback: %v", err)
		}
		if row := fieldIndexRow(rows, "SessionId", types.IndexCategoryCustom); row == nil {
			return fmt.Errorf("account policy field not indexed after fallback: %+v", rows)
		}
		if row := fieldIndexRow(rows, "RequestId", types.IndexCategoryCustom); row != nil {
			return fmt.Errorf("deleted group policy still indexes: %+v", row)
		}
		inactiveRows, err := tc.collectAllFieldIndexes([]string{group}, []types.IndexCategory{types.IndexCategoryInactive})
		if err != nil {
			return fmt.Errorf("inactive field indexes after fallback: %v", err)
		}
		if row := fieldIndexRow(inactiveRows, "RequestId", types.IndexCategoryInactive); row == nil {
			return fmt.Errorf("deleted group policy's field missing from INACTIVE: %+v", inactiveRows)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "DescribeFieldIndexes_DefaultCustomInactive", func() error {
		group, cleanupGroup, err := tc.newGroupStreamFixture("idxscan", "s1")
		if err != nil {
			return err
		}
		defer cleanupGroup()
		const stream = "s1"

		preTs := time.Now().UnixMilli() - 120000
		if err := tc.putLogEvent(group, stream, `{"TransactionId":"t-old","level":"INFO"}`, preTs); err != nil {
			return err
		}
		if _, err := tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"Fields":["TransactionId","userIdentity.accessKeyId"]}`),
		}); err != nil {
			return fmt.Errorf("put index policy: %v", err)
		}
		postTs := time.Now().UnixMilli() + 60000
		if err := tc.putLogEvent(group, stream, `{"TransactionId":"t-new","userIdentity":{"accessKeyId":"AKIA"}}`, postTs); err != nil {
			return err
		}

		rows, err := tc.collectAllFieldIndexes([]string{group}, nil)
		if err != nil {
			return fmt.Errorf("describe field indexes: %v", err)
		}
		// The DEFAULT listing: every documented default field, the
		// identity ones bounded by the group's event history.
		for _, name := range []string{
			"@logStream", "@aws.region", "@aws.account", "@source.log",
			"@data_source_name", "@data_source_type", "@data_format",
			"traceId", "severityText", "attributes.session.id",
		} {
			if row := fieldIndexRow(rows, name, types.IndexCategoryDefault); row == nil {
				return fmt.Errorf("DEFAULT row %s missing: %d rows", name, len(rows))
			}
		}
		if row := fieldIndexRow(rows, "@logStream", types.IndexCategoryDefault); row == nil || aws.ToInt64(row.FirstEventTime) != preTs {
			return fmt.Errorf("@logStream DEFAULT bounds: %+v", row)
		}
		// CUSTOM bounds count only events at or after the policy's
		// update stamp: the pre-policy match does not count.
		custom := fieldIndexRow(rows, "TransactionId", types.IndexCategoryCustom)
		if custom == nil {
			return fmt.Errorf("CUSTOM TransactionId missing")
		}
		if aws.ToInt64(custom.FirstEventTime) != postTs || aws.ToInt64(custom.LastEventTime) != postTs {
			return fmt.Errorf("CUSTOM bounds: first=%d last=%d want %d",
				aws.ToInt64(custom.FirstEventTime), aws.ToInt64(custom.LastEventTime), postTs)
		}
		if custom.LogGroupIdentifier == nil || !strings.Contains(aws.ToString(custom.LogGroupIdentifier), group) {
			return fmt.Errorf("single-group CUSTOM row omits its group ARN: %+v", custom)
		}
		if row := fieldIndexRow(rows, "userIdentity.accessKeyId", types.IndexCategoryCustom); row == nil {
			return fmt.Errorf("dotted JSON path missing from CUSTOM")
		}

		// Dropping the field moves it to INACTIVE; its bounds freeze at
		// the range the ingestion scans accumulated while the policy
		// watched the field ("the earliest log event that matches this
		// field index, after the index policy that contains it was
		// created" — the pre-policy event never carried into the index).
		if _, err := tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"Fields":["userIdentity.accessKeyId"]}`),
		}); err != nil {
			return fmt.Errorf("replace index policy: %v", err)
		}
		rows, err = tc.collectAllFieldIndexes([]string{group}, []types.IndexCategory{types.IndexCategoryInactive})
		if err != nil {
			return fmt.Errorf("describe inactive: %v", err)
		}
		inactive := fieldIndexRow(rows, "TransactionId", types.IndexCategoryInactive)
		if inactive == nil || aws.ToInt64(inactive.FirstEventTime) != postTs {
			return fmt.Errorf("INACTIVE TransactionId bounds: %+v", inactive)
		}

		autoRows, err := tc.collectAllFieldIndexes([]string{group}, []types.IndexCategory{types.IndexCategoryAuto})
		if err != nil {
			return fmt.Errorf("describe auto: %v", err)
		}
		if len(autoRows) != 0 {
			return fmt.Errorf("AUTO category carries rows: %+v", autoRows)
		}
		return nil
	}))

	results = append(results, tc.runner.RunTest("logs", "IndexPolicy_RejectionRows", func() error {
		group, cleanupGroup, err := tc.newLogGroupFixture("idxrej")
		if err != nil {
			return err
		}
		defer cleanupGroup()

		var ipe *types.InvalidParameterException
		many := make([]string, 21)
		for i := range many {
			many[i] = fmt.Sprintf("F%02d", i)
		}
		_, err = tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"Fields": ["` + strings.Join(many, `","`) + `"]}`),
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("twenty-first field: %v", err)
		}
		_, err = tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"Fields":["A"],"FieldsV2":{"A":{"type":"FACET"}}}`),
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("overlapping members: %v", err)
		}
		_, err = tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"Fields":["@userId"]}`),
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("single-at custom field: %v", err)
		}
		_, err = tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{"FieldsV2":{"A":{"type":"PARTIAL"}}}`),
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("unsupported type: %v", err)
		}
		_, err = tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group),
			PolicyDocument:     aws.String(`{}`),
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("empty document: %v", err)
		}
		var rnf *types.ResourceNotFoundException
		_, err = tc.client.PutIndexPolicy(tc.ctx, &cloudwatchlogs.PutIndexPolicyInput{
			LogGroupIdentifier: aws.String(group + "-absent"),
			PolicyDocument:     aws.String(`{"Fields":["A"]}`),
		})
		if !errors.As(err, &rnf) {
			return fmt.Errorf("missing group: %v", err)
		}
		_, err = tc.client.DescribeIndexPolicies(tc.ctx, &cloudwatchlogs.DescribeIndexPoliciesInput{
			LogGroupIdentifiers: []string{group, group + "-second"},
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("two identifiers: %v", err)
		}
		_, err = tc.client.DescribeFieldIndexes(tc.ctx, &cloudwatchlogs.DescribeFieldIndexesInput{
			LogGroupIdentifiers: []string{group},
			IndexCategories:     []types.IndexCategory{types.IndexCategory("NOT_A_CATEGORY")},
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("invalid category: %v", err)
		}
		// The data-source criteria form rejects: the platform carries
		// no vended-logs data sources.
		_, err = tc.client.PutAccountPolicy(tc.ctx, &cloudwatchlogs.PutAccountPolicyInput{
			PolicyName:        aws.String(group + "-ds"),
			PolicyDocument:    aws.String(`{"Fields":["A"]}`),
			PolicyType:        types.PolicyTypeFieldIndexPolicy,
			SelectionCriteria: aws.String("DataSourceName = amazon_vpc AND DataSourceType = flow"),
		})
		if !errors.As(err, &ipe) {
			return fmt.Errorf("data-source criteria: %v", err)
		}
		return nil
	}))

	return results
}
