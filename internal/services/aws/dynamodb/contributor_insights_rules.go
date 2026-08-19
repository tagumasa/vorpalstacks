package dynamodb

import (
	"fmt"
	"strconv"
	"strings"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// contributorInsightRulePrefix is the prefix of the CloudWatch Contributor
// Insights rules that DynamoDB creates on a contributor insights-enabled
// table. The middle segment selects the tracked key layout and the kind of
// traffic: PKC/SKC track accesses, PKT/SKT track throttles.
const contributorInsightRulePrefix = "DynamoDBContributorInsights"

// ContributorInsightsRuleNames derives the CloudWatch rule names a
// contributor insights-enabled table exposes: accessed-key rules in
// accessed-and-throttled mode plus the throttled-key rules, scoped to the
// tracked key layouts of the table's key schema. The names are derived
// from the table state, so disabling insights removes them without any
// stored rule list.
func ContributorInsightsRuleNames(table *dbstore.Table) []string {
	if table == nil || !table.ContributorInsightsEnabled {
		return nil
	}
	mode := table.ContributorInsightsMode
	if mode == "" {
		mode = "ACCESSED_AND_THROTTLED_KEYS"
	}
	stamp := table.ContributorInsightsUpdatedAt.UnixMilli()
	if stamp <= 0 {
		stamp = table.CreationDateTime.UnixMilli()
	}

	hasSortKey := false
	for _, ks := range table.KeySchema {
		if ks.KeyType == dbstore.KeyTypeRange {
			hasSortKey = true
			break
		}
	}

	var rules []string
	if mode == "ACCESSED_AND_THROTTLED_KEYS" {
		rules = append(rules, contributorInsightRuleName("PKC", table.Name, stamp))
		if hasSortKey {
			rules = append(rules, contributorInsightRuleName("SKC", table.Name, stamp))
		}
	}
	rules = append(rules, contributorInsightRuleName("PKT", table.Name, stamp))
	if hasSortKey {
		rules = append(rules, contributorInsightRuleName("SKT", table.Name, stamp))
	}
	return rules
}

func contributorInsightRuleName(layout, resource string, stamp int64) string {
	return fmt.Sprintf("%s-%s-%s-%d", contributorInsightRulePrefix, layout, resource, stamp)
}

// ParseContributorInsightsRule extracts the table name and tracked key
// layout from a DynamoDB contributor insights rule name. Table names may
// contain hyphens, so the layout is the first segment, the creation
// timestamp the last, and everything between them the table name.
func ParseContributorInsightsRule(ruleName string) (tableName, layout string, ok bool) {
	rest, found := strings.CutPrefix(ruleName, contributorInsightRulePrefix+"-")
	if !found {
		return "", "", false
	}
	parts := strings.Split(rest, "-")
	if len(parts) < 3 {
		return "", "", false
	}
	switch parts[0] {
	case "PKC", "SKC", "PKT", "SKT":
	default:
		return "", "", false
	}
	if _, err := strconv.ParseInt(parts[len(parts)-1], 10, 64); err != nil {
		return "", "", false
	}
	return strings.Join(parts[1:len(parts)-1], "-"), parts[0], true
}
