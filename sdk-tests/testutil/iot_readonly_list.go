package testutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// runIoTReadonlyListTests covers List/Describe operations that require no
// prior resource setup. Each List test paginates through ALL pages (per
// AGENTS.md requirement) and verifies the response structure. Describe tests
// verify the response contains the expected top-level field.
func (r *TestRunner) runIoTReadonlyListTests(tc *iotTestContext) []TestResult {
	var results []TestResult

	results = append(results, r.RunTest("iot", "List_BillingGroups", func() error {
		var token *string
		for {
			out, err := tc.client.ListBillingGroups(tc.ctx, &iot.ListBillingGroupsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, g := range out.BillingGroups {
				if g.GroupName == nil || *g.GroupName == "" {
					return fmt.Errorf("ListBillingGroups returned group with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ThingGroups", func() error {
		var token *string
		for {
			out, err := tc.client.ListThingGroups(tc.ctx, &iot.ListThingGroupsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, g := range out.ThingGroups {
				if g.GroupName == nil || *g.GroupName == "" {
					return fmt.Errorf("ListThingGroups returned group with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ThingTypes", func() error {
		var token *string
		for {
			out, err := tc.client.ListThingTypes(tc.ctx, &iot.ListThingTypesInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, t := range out.ThingTypes {
				if t.ThingTypeName == nil || *t.ThingTypeName == "" {
					return fmt.Errorf("ListThingTypes returned type with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_Indices", func() error {
		var token *string
		for {
			out, err := tc.client.ListIndices(tc.ctx, &iot.ListIndicesInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, name := range out.IndexNames {
				if name == "" {
					return fmt.Errorf("ListIndices returned index with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_FleetMetrics", func() error {
		var token *string
		for {
			out, err := tc.client.ListFleetMetrics(tc.ctx, &iot.ListFleetMetricsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, m := range out.FleetMetrics {
				if m.MetricName == nil || *m.MetricName == "" {
					return fmt.Errorf("ListFleetMetrics returned metric with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_CustomMetrics", func() error {
		var token *string
		for {
			out, err := tc.client.ListCustomMetrics(tc.ctx, &iot.ListCustomMetricsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, name := range out.MetricNames {
				if name == "" {
					return fmt.Errorf("ListCustomMetrics returned metric with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_Dimensions", func() error {
		var token *string
		for {
			out, err := tc.client.ListDimensions(tc.ctx, &iot.ListDimensionsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, name := range out.DimensionNames {
				if name == "" {
					return fmt.Errorf("ListDimensions returned dimension with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_MitigationActions", func() error {
		var token *string
		for {
			out, err := tc.client.ListMitigationActions(tc.ctx, &iot.ListMitigationActionsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, a := range out.ActionIdentifiers {
				if a.ActionName == nil || *a.ActionName == "" {
					return fmt.Errorf("ListMitigationActions returned action with empty name")
				}
				if a.ActionArn == nil || *a.ActionArn == "" {
					return fmt.Errorf("ListMitigationActions returned action with empty ARN")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ScheduledAudits", func() error {
		var token *string
		for {
			out, err := tc.client.ListScheduledAudits(tc.ctx, &iot.ListScheduledAuditsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, a := range out.ScheduledAudits {
				if a.ScheduledAuditName == nil || *a.ScheduledAuditName == "" {
					return fmt.Errorf("ListScheduledAudits returned audit with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_RoleAliases", func() error {
		var marker *string
		for {
			out, err := tc.client.ListRoleAliases(tc.ctx, &iot.ListRoleAliasesInput{Marker: marker})
			if err != nil {
				return err
			}
			for _, alias := range out.RoleAliases {
				if alias == "" {
					return fmt.Errorf("ListRoleAliases returned empty role alias")
				}
			}
			if out.NextMarker == nil || *out.NextMarker == "" {
				break
			}
			marker = out.NextMarker
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_TopicRules", func() error {
		var token *string
		for {
			out, err := tc.client.ListTopicRules(tc.ctx, &iot.ListTopicRulesInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, rule := range out.Rules {
				if rule.RuleName == nil || *rule.RuleName == "" {
					return fmt.Errorf("ListTopicRules returned rule with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_TopicRuleDestinations", func() error {
		var token *string
		for {
			out, err := tc.client.ListTopicRuleDestinations(tc.ctx, &iot.ListTopicRuleDestinationsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, d := range out.DestinationSummaries {
				if d.Status == "" {
					return fmt.Errorf("ListTopicRuleDestinations returned destination with empty status")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditTasks", func() error {
		now := time.Now()
		start := now.Add(-24 * time.Hour)
		var token *string
		for {
			out, err := tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
				StartTime: &start,
				EndTime:   &now,
				NextToken: token,
			})
			if err != nil {
				return err
			}
			for _, t := range out.Tasks {
				if t.TaskId == nil || *t.TaskId == "" {
					return fmt.Errorf("ListAuditTasks returned task with empty taskId")
				}
				if t.TaskStatus == "" {
					return fmt.Errorf("ListAuditTasks returned task with empty status")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditFindings", func() error {
		var token *string
		for {
			out, err := tc.client.ListAuditFindings(tc.ctx, &iot.ListAuditFindingsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, f := range out.Findings {
				if f.CheckName == nil || *f.CheckName == "" {
					return fmt.Errorf("ListAuditFindings returned finding with empty checkName")
				}
				if f.TaskId == nil || *f.TaskId == "" {
					return fmt.Errorf("ListAuditFindings returned finding with empty taskId")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditSuppressions", func() error {
		var token *string
		for {
			out, err := tc.client.ListAuditSuppressions(tc.ctx, &iot.ListAuditSuppressionsInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, s := range out.Suppressions {
				if s.CheckName == nil || *s.CheckName == "" {
					return fmt.Errorf("ListAuditSuppressions returned suppression with empty checkName")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_OutgoingCertificates", func() error {
		var marker *string
		for {
			out, err := tc.client.ListOutgoingCertificates(tc.ctx, &iot.ListOutgoingCertificatesInput{Marker: marker})
			if err != nil {
				return err
			}
			for _, c := range out.OutgoingCertificates {
				if c.CertificateId == nil || *c.CertificateId == "" {
					return fmt.Errorf("ListOutgoingCertificates returned certificate with empty ID")
				}
				if c.TransferDate == nil {
					return fmt.Errorf("ListOutgoingCertificates returned certificate with nil transferDate")
				}
			}
			if out.NextMarker == nil || *out.NextMarker == "" {
				break
			}
			marker = out.NextMarker
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ManagedJobTemplates", func() error {
		var token *string
		for {
			out, err := tc.client.ListManagedJobTemplates(tc.ctx, &iot.ListManagedJobTemplatesInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, t := range out.ManagedJobTemplates {
				if t.TemplateName == nil || *t.TemplateName == "" {
					return fmt.Errorf("ListManagedJobTemplates returned template with empty name")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
		}
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ThingRegistrationTasks", func() error {
		var token *string
		for {
			out, err := tc.client.ListThingRegistrationTasks(tc.ctx, &iot.ListThingRegistrationTasksInput{NextToken: token})
			if err != nil {
				return err
			}
			for _, id := range out.TaskIds {
				if id == "" {
					return fmt.Errorf("ListThingRegistrationTasks returned empty task ID")
				}
			}
			if out.NextToken == nil || *out.NextToken == "" {
				break
			}
			token = out.NextToken
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
