package testutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	iottypes "github.com/aws/aws-sdk-go-v2/service/iot/types"
)

// runIoTReadonlyListTests covers List/Describe operations that require no
// prior resource setup. Each List test paginates through ALL pages
// (pagination is mandatory for list operations) and verifies the response
// structure. Describe tests verify the response contains the expected
// top-level field.
func (r *TestRunner) runIoTReadonlyListTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "List_BillingGroups", func() error {
		groups, err := paginate(func(next *string) ([]iottypes.GroupNameAndArn, *string, error) {
			out, err := tc.client.ListBillingGroups(tc.ctx, &iot.ListBillingGroupsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.BillingGroups, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, g := range groups {
			if g.GroupName == nil || *g.GroupName == "" {
				return fmt.Errorf("ListBillingGroups returned group with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ThingGroups", func() error {
		groups, err := paginate(func(next *string) ([]iottypes.GroupNameAndArn, *string, error) {
			out, err := tc.client.ListThingGroups(tc.ctx, &iot.ListThingGroupsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.ThingGroups, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, g := range groups {
			if g.GroupName == nil || *g.GroupName == "" {
				return fmt.Errorf("ListThingGroups returned group with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ThingTypes", func() error {
		types, err := paginate(func(next *string) ([]iottypes.ThingTypeDefinition, *string, error) {
			out, err := tc.client.ListThingTypes(tc.ctx, &iot.ListThingTypesInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.ThingTypes, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, t := range types {
			if t.ThingTypeName == nil || *t.ThingTypeName == "" {
				return fmt.Errorf("ListThingTypes returned type with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_Indices", func() error {
		names, err := paginate(func(next *string) ([]string, *string, error) {
			out, err := tc.client.ListIndices(tc.ctx, &iot.ListIndicesInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.IndexNames, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, name := range names {
			if name == "" {
				return fmt.Errorf("ListIndices returned index with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_FleetMetrics", func() error {
		metrics, err := paginate(func(next *string) ([]iottypes.FleetMetricNameAndArn, *string, error) {
			out, err := tc.client.ListFleetMetrics(tc.ctx, &iot.ListFleetMetricsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.FleetMetrics, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, m := range metrics {
			if m.MetricName == nil || *m.MetricName == "" {
				return fmt.Errorf("ListFleetMetrics returned metric with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_CustomMetrics", func() error {
		names, err := paginate(func(next *string) ([]string, *string, error) {
			out, err := tc.client.ListCustomMetrics(tc.ctx, &iot.ListCustomMetricsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.MetricNames, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, name := range names {
			if name == "" {
				return fmt.Errorf("ListCustomMetrics returned metric with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_Dimensions", func() error {
		names, err := paginate(func(next *string) ([]string, *string, error) {
			out, err := tc.client.ListDimensions(tc.ctx, &iot.ListDimensionsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.DimensionNames, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, name := range names {
			if name == "" {
				return fmt.Errorf("ListDimensions returned dimension with empty name")
			}
		}
		return nil
	}))

	// A nextToken the server could never have issued (non-numeric, or an
	// offset beyond the result set) must be rejected, not silently treated
	// as page zero.
	results = append(results, r.RunTest("iot", "ListDimensions_InvalidNextTokenRejected", func() error {
		_, err := tc.client.ListDimensions(tc.ctx, &iot.ListDimensionsInput{NextToken: aws.String("abc")})
		if vErr := expectAWSErrorCode(err, "InvalidRequestException"); vErr != nil {
			return vErr
		}
		_, err = tc.client.ListDimensions(tc.ctx, &iot.ListDimensionsInput{NextToken: aws.String("999999")})
		return expectAWSErrorCode(err, "InvalidRequestException")
	}))

	results = append(results, r.RunTest("iot", "List_MitigationActions", func() error {
		actions, err := paginate(func(next *string) ([]iottypes.MitigationActionIdentifier, *string, error) {
			out, err := tc.client.ListMitigationActions(tc.ctx, &iot.ListMitigationActionsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.ActionIdentifiers, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, a := range actions {
			if a.ActionName == nil || *a.ActionName == "" {
				return fmt.Errorf("ListMitigationActions returned action with empty name")
			}
			if a.ActionArn == nil || *a.ActionArn == "" {
				return fmt.Errorf("ListMitigationActions returned action with empty ARN")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ScheduledAudits", func() error {
		audits, err := paginate(func(next *string) ([]iottypes.ScheduledAuditMetadata, *string, error) {
			out, err := tc.client.ListScheduledAudits(tc.ctx, &iot.ListScheduledAuditsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.ScheduledAudits, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, a := range audits {
			if a.ScheduledAuditName == nil || *a.ScheduledAuditName == "" {
				return fmt.Errorf("ListScheduledAudits returned audit with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_RoleAliases", func() error {
		aliases, err := paginate(func(next *string) ([]string, *string, error) {
			out, err := tc.client.ListRoleAliases(tc.ctx, &iot.ListRoleAliasesInput{Marker: next})
			if err != nil {
				return nil, nil, err
			}
			return out.RoleAliases, out.NextMarker, nil
		})
		if err != nil {
			return err
		}
		for _, alias := range aliases {
			if alias == "" {
				return fmt.Errorf("ListRoleAliases returned empty role alias")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_TopicRules", func() error {
		rules, err := paginate(func(next *string) ([]iottypes.TopicRuleListItem, *string, error) {
			out, err := tc.client.ListTopicRules(tc.ctx, &iot.ListTopicRulesInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.Rules, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, rule := range rules {
			if rule.RuleName == nil || *rule.RuleName == "" {
				return fmt.Errorf("ListTopicRules returned rule with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_TopicRuleDestinations", func() error {
		dests, err := paginate(func(next *string) ([]iottypes.TopicRuleDestinationSummary, *string, error) {
			out, err := tc.client.ListTopicRuleDestinations(tc.ctx, &iot.ListTopicRuleDestinationsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.DestinationSummaries, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, d := range dests {
			if d.Status == "" {
				return fmt.Errorf("ListTopicRuleDestinations returned destination with empty status")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditTasks", func() error {
		now := time.Now()
		start := now.Add(-24 * time.Hour)
		tasks, err := paginate(func(next *string) ([]iottypes.AuditTaskMetadata, *string, error) {
			out, err := tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
				StartTime: &start,
				EndTime:   &now,
				NextToken: next,
			})
			if err != nil {
				return nil, nil, err
			}
			return out.Tasks, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, t := range tasks {
			if t.TaskId == nil || *t.TaskId == "" {
				return fmt.Errorf("ListAuditTasks returned task with empty taskId")
			}
			if t.TaskStatus == "" {
				return fmt.Errorf("ListAuditTasks returned task with empty status")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditFindings", func() error {
		findings, err := paginate(func(next *string) ([]iottypes.AuditFinding, *string, error) {
			out, err := tc.client.ListAuditFindings(tc.ctx, &iot.ListAuditFindingsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.Findings, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, f := range findings {
			if f.CheckName == nil || *f.CheckName == "" {
				return fmt.Errorf("ListAuditFindings returned finding with empty checkName")
			}
			if f.TaskId == nil || *f.TaskId == "" {
				return fmt.Errorf("ListAuditFindings returned finding with empty taskId")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditSuppressions", func() error {
		suppressions, err := paginate(func(next *string) ([]iottypes.AuditSuppression, *string, error) {
			out, err := tc.client.ListAuditSuppressions(tc.ctx, &iot.ListAuditSuppressionsInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.Suppressions, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, s := range suppressions {
			if s.CheckName == nil || *s.CheckName == "" {
				return fmt.Errorf("ListAuditSuppressions returned suppression with empty checkName")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_OutgoingCertificates", func() error {
		certs, err := paginate(func(next *string) ([]iottypes.OutgoingCertificate, *string, error) {
			out, err := tc.client.ListOutgoingCertificates(tc.ctx, &iot.ListOutgoingCertificatesInput{Marker: next})
			if err != nil {
				return nil, nil, err
			}
			return out.OutgoingCertificates, out.NextMarker, nil
		})
		if err != nil {
			return err
		}
		for _, c := range certs {
			if c.CertificateId == nil || *c.CertificateId == "" {
				return fmt.Errorf("ListOutgoingCertificates returned certificate with empty ID")
			}
			if c.TransferDate == nil {
				return fmt.Errorf("ListOutgoingCertificates returned certificate with nil transferDate")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ManagedJobTemplates", func() error {
		templates, err := paginate(func(next *string) ([]iottypes.ManagedJobTemplateSummary, *string, error) {
			out, err := tc.client.ListManagedJobTemplates(tc.ctx, &iot.ListManagedJobTemplatesInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.ManagedJobTemplates, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, t := range templates {
			if t.TemplateName == nil || *t.TemplateName == "" {
				return fmt.Errorf("ListManagedJobTemplates returned template with empty name")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ThingRegistrationTasks", func() error {
		ids, err := paginate(func(next *string) ([]string, *string, error) {
			out, err := tc.client.ListThingRegistrationTasks(tc.ctx, &iot.ListThingRegistrationTasksInput{NextToken: next})
			if err != nil {
				return nil, nil, err
			}
			return out.TaskIds, out.NextToken, nil
		})
		if err != nil {
			return err
		}
		for _, id := range ids {
			if id == "" {
				return fmt.Errorf("ListThingRegistrationTasks returned empty task ID")
			}
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Describe_Endpoint", func() error {
		resp, err := tc.client.DescribeEndpoint(tc.ctx, &iot.DescribeEndpointInput{})
		if err != nil {
			return err
		}
		return validateEndpointAddress(resp.EndpointAddress, true)
	}))

	for _, ep := range []string{"iot:Data", "iot:Data-ATS", "iot:Data-ALPN", "iot:Jobs"} {
		et := ep
		results = append(results, r.RunTest("iot", "Describe_Endpoint_"+strings.ReplaceAll(et, ":", "_"), func() error {
			resp, err := tc.client.DescribeEndpoint(tc.ctx, &iot.DescribeEndpointInput{
				EndpointType: aws.String(et),
			})
			if err != nil {
				return err
			}
			return validateEndpointAddress(resp.EndpointAddress, true)
		}))
	}

	results = append(results, r.RunTest("iot", "Describe_Endpoint_CredentialProvider", func() error {
		resp, err := tc.client.DescribeEndpoint(tc.ctx, &iot.DescribeEndpointInput{
			EndpointType: aws.String("iot:CredentialProvider"),
		})
		if err != nil {
			return err
		}
		return validateEndpointAddress(resp.EndpointAddress, false)
	}))

	results = append(results, r.RunTest("iot", "Describe_EventConfigurations", func() error {
		out, err := tc.client.DescribeEventConfigurations(tc.ctx, &iot.DescribeEventConfigurationsInput{})
		if err != nil {
			return err
		}
		if out.EventConfigurations == nil {
			return fmt.Errorf("DescribeEventConfigurations returned nil EventConfigurations")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Describe_EncryptionConfiguration", func() error {
		out, err := tc.client.DescribeEncryptionConfiguration(tc.ctx, &iot.DescribeEncryptionConfigurationInput{})
		if err != nil {
			return err
		}
		if out.EncryptionType == "" {
			return fmt.Errorf("DescribeEncryptionConfiguration returned empty encryptionType")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Describe_AccountAuditConfiguration", func() error {
		out, err := tc.client.DescribeAccountAuditConfiguration(tc.ctx, &iot.DescribeAccountAuditConfigurationInput{})
		if err != nil {
			return err
		}
		if out.AuditCheckConfigurations == nil {
			return fmt.Errorf("DescribeAccountAuditConfiguration returned nil auditCheckConfigurations")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Get_IndexingConfiguration", func() error {
		out, err := tc.client.GetIndexingConfiguration(tc.ctx, &iot.GetIndexingConfigurationInput{})
		if err != nil {
			return err
		}
		if out.ThingIndexingConfiguration == nil {
			return fmt.Errorf("GetIndexingConfiguration returned nil thingIndexingConfiguration")
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "Get_V2LoggingOptions", func() error {
		out, err := tc.client.GetV2LoggingOptions(tc.ctx, &iot.GetV2LoggingOptionsInput{})
		if err != nil {
			return err
		}
		if out.DefaultLogLevel == "" {
			return fmt.Errorf("GetV2LoggingOptions returned empty defaultLogLevel")
		}
		return nil
	}))

	return results
}

// IoT MQTT brokers now use per-region dynamic port allocation
// (mirrors internal/common/serviceports.DynamicRangeStart..End). MQTT
// endpoint types must return a port within this range, or the legacy
// fixed port (serviceports.IotMQTT = 50107) as a fallback, so that IoT
// clients can connect to the correct regional broker.
const (
	iotMQTTPortRangeStart = 50200
	iotMQTTPortRangeEnd   = 50400
	iotMQTTLegacyPort     = 50107
)

// isMQTTBrokerPort returns true if the port is in the dynamic allocation
// range or equals the legacy fixed broker port.
func isMQTTBrokerPort(port int) bool {
	return port == iotMQTTLegacyPort || (port >= iotMQTTPortRangeStart && port <= iotMQTTPortRangeEnd)
}

// validateEndpointAddress checks that a DescribeEndpoint response contains
// a connectable localhost:{port} address. When expectMQTT is true, the port
// must be a valid regional MQTT broker port (dynamic range or legacy
// fallback).
func validateEndpointAddress(addr *string, expectMQTT bool) error {
	if addr == nil || *addr == "" {
		return fmt.Errorf("DescribeEndpoint returned empty endpointAddress")
	}
	if !strings.HasPrefix(*addr, "localhost:") {
		return fmt.Errorf("endpointAddress should be localhost:{port}, got %q", *addr)
	}
	if expectMQTT {
		portStr := strings.TrimPrefix(*addr, "localhost:")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return fmt.Errorf("MQTT endpoint port not numeric: %q", *addr)
		}
		if !isMQTTBrokerPort(port) {
			return fmt.Errorf("MQTT endpoint port %d is not a valid broker port (dynamic range %d-%d or legacy %d)", port, iotMQTTPortRangeStart, iotMQTTPortRangeEnd, iotMQTTLegacyPort)
		}
	}
	return nil
}
