package iot

import (
	"fmt"

	svcarn "vorpalstacks/internal/utils/aws/arn"
)

func BuildThingARN(accountID, region, thingName string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("thing/%s", thingName))
}

func BuildCertificateARN(accountID, region, certID string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("cert/%s", certID))
}

func BuildPolicyARN(accountID, region, policyName string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("policy/%s", policyName))
}

func BuildRuleARN(accountID, region, ruleName string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("rule/%s", ruleName))
}

func BuildJobARN(accountID, region, jobID string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("job/%s", jobID))
}

func BuildThingTypeARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("thingtype/%s", name))
}

func BuildThingGroupARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("thinggroup/%s", name))
}

func BuildRoleAliasARN(accountID, region, alias string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("rolealias/%s", alias))
}

func BuildBillingGroupARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("billinggroup/%s", name))
}

func BuildAuthorizerARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("authorizer/%s", name))
}

func BuildProvisioningTemplateARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("provisioningtemplate/%s", name))
}

func BuildDetectorModelARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iotevents", fmt.Sprintf("detectorModel/%s", name))
}

func BuildInputARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iotevents", fmt.Sprintf("input/%s", name))
}

func BuildSecurityProfileARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("securityprofile/%s", name))
}

func BuildDomainConfigurationARN(accountID, region, name string) string {
	return svcarn.NewARNBuilder(accountID, region).Build("iot", fmt.Sprintf("domainconfiguration/%s", name))
}
