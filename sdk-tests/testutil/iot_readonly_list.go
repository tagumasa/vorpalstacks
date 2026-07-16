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
		out, err := tc.client.ListIndices(tc.ctx, &iot.ListIndicesInput{})
		if err != nil {
			return err
		}
		_ = out.IndexNames
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_FleetMetrics", func() error {
		out, err := tc.client.ListFleetMetrics(tc.ctx, &iot.ListFleetMetricsInput{})
		if err != nil {
			return err
		}
		_ = out.FleetMetrics
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_CustomMetrics", func() error {
		out, err := tc.client.ListCustomMetrics(tc.ctx, &iot.ListCustomMetricsInput{})
		if err != nil {
			return err
		}
		_ = out.MetricNames
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_Dimensions", func() error {
		out, err := tc.client.ListDimensions(tc.ctx, &iot.ListDimensionsInput{})
		if err != nil {
			return err
		}
		_ = out.DimensionNames
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_MitigationActions", func() error {
		out, err := tc.client.ListMitigationActions(tc.ctx, &iot.ListMitigationActionsInput{})
		if err != nil {
			return err
		}
		_ = out.ActionIdentifiers
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ScheduledAudits", func() error {
		out, err := tc.client.ListScheduledAudits(tc.ctx, &iot.ListScheduledAuditsInput{})
		if err != nil {
			return err
		}
		_ = out.ScheduledAudits
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_RoleAliases", func() error {
		out, err := tc.client.ListRoleAliases(tc.ctx, &iot.ListRoleAliasesInput{})
		if err != nil {
			return err
		}
		_ = out.RoleAliases
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
		out, err := tc.client.ListTopicRuleDestinations(tc.ctx, &iot.ListTopicRuleDestinationsInput{})
		if err != nil {
			return err
		}
		_ = out.DestinationSummaries
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditTasks", func() error {
		now := time.Now()
		start := now.Add(-24 * time.Hour)
		out, err := tc.client.ListAuditTasks(tc.ctx, &iot.ListAuditTasksInput{
			StartTime: &start,
			EndTime:   &now,
		})
		if err != nil {
			return err
		}
		_ = out.Tasks
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditFindings", func() error {
		out, err := tc.client.ListAuditFindings(tc.ctx, &iot.ListAuditFindingsInput{})
		if err != nil {
			return err
		}
		_ = out.Findings
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_AuditSuppressions", func() error {
		out, err := tc.client.ListAuditSuppressions(tc.ctx, &iot.ListAuditSuppressionsInput{})
		if err != nil {
			return err
		}
		_ = out.Suppressions
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_OutgoingCertificates", func() error {
		out, err := tc.client.ListOutgoingCertificates(tc.ctx, &iot.ListOutgoingCertificatesInput{})
		if err != nil {
			return err
		}
		_ = out.OutgoingCertificates
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ManagedJobTemplates", func() error {
		out, err := tc.client.ListManagedJobTemplates(tc.ctx, &iot.ListManagedJobTemplatesInput{})
		if err != nil {
			return err
		}
		_ = out.ManagedJobTemplates
		return nil
	}))

	results = append(results, r.RunTest("iot", "List_ThingRegistrationTasks", func() error {
		out, err := tc.client.ListThingRegistrationTasks(tc.ctx, &iot.ListThingRegistrationTasksInput{})
		if err != nil {
			return err
		}
		_ = out.TaskIds
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
		_ = out.EncryptionType
		return nil
	}))

	results = append(results, r.RunTest("iot", "Describe_AccountAuditConfiguration", func() error {
		out, err := tc.client.DescribeAccountAuditConfiguration(tc.ctx, &iot.DescribeAccountAuditConfigurationInput{})
		if err != nil {
			return err
		}
		_ = out.AuditCheckConfigurations
		return nil
	}))

	results = append(results, r.RunTest("iot", "Get_IndexingConfiguration", func() error {
		out, err := tc.client.GetIndexingConfiguration(tc.ctx, &iot.GetIndexingConfigurationInput{})
		if err != nil {
			return err
		}
		_ = out.ThingIndexingConfiguration
		return nil
	}))

	results = append(results, r.RunTest("iot", "Get_V2LoggingOptions", func() error {
		out, err := tc.client.GetV2LoggingOptions(tc.ctx, &iot.GetV2LoggingOptionsInput{})
		if err != nil {
			return err
		}
		_ = out.DefaultLogLevel
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
