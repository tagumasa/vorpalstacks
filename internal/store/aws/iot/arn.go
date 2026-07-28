package iot

import (
	"strings"

	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ThingNameFromARN extracts the thing name from an IoT thing ARN
// (arn:aws:iot:region:account:thing/<name>). Returns an empty string for
// non-thing ARNs or malformed input.
func ThingNameFromARN(arn string) string {
	if idx := strings.Index(arn, ":thing/"); idx >= 0 {
		return arn[idx+len(":thing/"):]
	}
	return ""
}

func BuildThingARN(accountID, region, thingName string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Thing(thingName)
}

func BuildCertificateARN(accountID, region, certID string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Certificate(certID)
}

func BuildPolicyARN(accountID, region, policyName string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Policy(policyName)
}

func BuildRuleARN(accountID, region, ruleName string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Rule(ruleName)
}

func BuildJobARN(accountID, region, jobID string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Job(jobID)
}

func BuildThingTypeARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().ThingType(name)
}

func BuildThingGroupARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().ThingGroup(name)
}

func BuildRoleAliasARN(accountID, region, alias string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().RoleAlias(alias)
}

func BuildBillingGroupARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().BillingGroup(name)
}

func BuildAuthorizerARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Authorizer(name)
}

func BuildProvisioningTemplateARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().ProvisioningTemplate(name)
}

func BuildSecurityProfileARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().SecurityProfile(name)
}

func BuildStreamARN(accountID, region, streamID string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Stream(streamID)
}

func BuildPackageARN(accountID, region, packageName string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Package(packageName)
}

func BuildCommandARN(accountID, region, commandID string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Command(commandID)
}

func BuildCertificateProviderARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().CertificateProvider(name)
}

func BuildDomainConfigurationARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().DomainConfiguration(name)
}

func BuildCACertificateARN(accountID, region, certID string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().CACertificate(certID)
}

func BuildCustomMetricARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().CustomMetric(name)
}

func BuildDimensionARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().Dimension(name)
}

func BuildMitigationActionARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().MitigationAction(name)
}

func BuildScheduledAuditARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().ScheduledAudit(name)
}

func BuildFleetMetricARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().FleetMetric(name)
}

func BuildTopicRuleDestinationARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().TopicRuleDestination(name)
}

func BuildJobTemplateARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().JobTemplate(name)
}

func BuildOTAUpdateARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).IoT().OTAUpdate(name)
}
