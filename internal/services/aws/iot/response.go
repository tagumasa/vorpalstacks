package iot

import (
	"fmt"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

func thingResponse(t *iotstore.Thing) map[string]interface{} {
	return map[string]interface{}{
		"thingName":        t.ThingName,
		"thingArn":         t.ThingARN,
		"thingId":          t.ThingID,
		"thingTypeName":    t.ThingTypeName,
		"attributes":       t.Attributes,
		"attributeNames":   t.AttributeNames,
		"version":          t.Version,
		"creationDate":     t.CreationDate.Unix(),
		"lastModifiedDate": t.LastModifiedDate.Unix(),
	}
}

func thingDescribeResponse(t *iotstore.Thing) map[string]interface{} {
	return map[string]interface{}{
		"thingName":        t.ThingName,
		"thingArn":         t.ThingARN,
		"thingId":          t.ThingID,
		"thingTypeName":    t.ThingTypeName,
		"attributes":       ensureMap(t.Attributes),
		"attributeNames":   t.AttributeNames,
		"version":          t.Version,
		"creationDate":     t.CreationDate.Unix(),
		"lastModifiedDate": t.LastModifiedDate.Unix(),
	}
}

func thingTypeResponse(tt *iotstore.ThingType) map[string]interface{} {
	props := map[string]interface{}{"thingTypeDescription": tt.Description}
	if len(tt.SearchableAttributes) > 0 {
		props["searchableAttributes"] = tt.SearchableAttributes
	}
	return map[string]interface{}{
		"thingTypeName":       tt.ThingTypeName,
		"thingTypeArn":        tt.ThingTypeARN,
		"thingTypeId":         tt.ThingTypeID,
		"description":         tt.Description,
		"thingTypeProperties": props,
		"version":             tt.Version,
		"creationDate":        tt.CreationDate.Unix(),
		"lastModifiedDate":    tt.LastModifiedDate.Unix(),
	}
}

func thingTypeDescribeResponse(tt *iotstore.ThingType) map[string]interface{} {
	props := map[string]interface{}{"thingTypeDescription": tt.Description}
	if len(tt.SearchableAttributes) > 0 {
		props["searchableAttributes"] = tt.SearchableAttributes
	}
	metadata := map[string]interface{}{
		"deprecated":   tt.Deprecated,
		"creationDate": tt.CreationDate.Unix(),
	}
	if !tt.DeprecationDate.IsZero() {
		metadata["deprecationDate"] = tt.DeprecationDate.Unix()
	}
	return map[string]interface{}{
		"thingTypeName":       tt.ThingTypeName,
		"thingTypeArn":        tt.ThingTypeARN,
		"thingTypeId":         tt.ThingTypeID,
		"description":         tt.Description,
		"thingTypeProperties": props,
		"thingTypeMetadata":   metadata,
	}
}

func thingGroupResponse(tg *iotstore.ThingGroup) map[string]interface{} {
	return map[string]interface{}{
		"thingGroupName":   tg.GroupName,
		"thingGroupArn":    tg.GroupARN,
		"thingGroupId":     tg.GroupID,
		"parentGroupName":  tg.ParentGroupName,
		"description":      tg.Description,
		"attributes":       tg.Attributes,
		"creationDate":     tg.CreationDate.Unix(),
		"lastModifiedDate": tg.LastModifiedDate.Unix(),
	}
}

// groupNameAndArnResponse renders the GroupNameAndArn summary shape used by
// ListThingGroups and ListBillingGroups (AWS returns groupName/groupArn, not
// the full thingGroupName/thingGroupArn detail shape).
func groupNameAndArnResponse(name, arn string) map[string]interface{} {
	return map[string]interface{}{
		"groupName": name,
		"groupArn":  arn,
	}
}

func thingGroupDescribeResponse(tg *iotstore.ThingGroup) map[string]interface{} {
	return map[string]interface{}{
		"thingGroupName":  tg.GroupName,
		"thingGroupArn":   tg.GroupARN,
		"thingGroupId":    tg.GroupID,
		"parentGroupName": tg.ParentGroupName,
		"description":     tg.Description,
		"thingGroupProperties": map[string]interface{}{
			"thingGroupDescription": tg.Description,
			"attributePayload": map[string]interface{}{
				"attributes": ensureMap(tg.Attributes),
			},
		},
		"attributes":       tg.Attributes,
		"creationDate":     tg.CreationDate.Unix(),
		"lastModifiedDate": tg.LastModifiedDate.Unix(),
	}
}

func billingGroupResponse(bg *iotstore.BillingGroup) map[string]interface{} {
	return map[string]interface{}{
		"billingGroupName": bg.GroupName,
		"billingGroupArn":  bg.GroupARN,
		"billingGroupId":   bg.GroupID,
		"description":      bg.Description,
		"creationDate":     bg.CreationDate.Unix(),
		"lastModifiedDate": bg.LastModifiedDate.Unix(),
	}
}

func billingGroupDescribeResponse(bg *iotstore.BillingGroup) map[string]interface{} {
	return map[string]interface{}{
		"billingGroupName":       bg.GroupName,
		"billingGroupArn":        bg.GroupARN,
		"billingGroupId":         bg.GroupID,
		"version":                bg.Version,
		"description":            bg.Description,
		"billingGroupProperties": map[string]interface{}{"billingGroupDescription": bg.Description},
		"attributes":             bg.Attributes,
		"creationDate":           bg.CreationDate.Unix(),
		"lastModifiedDate":       bg.LastModifiedDate.Unix(),
	}
}

func certificateResponse(c *iotstore.Certificate) map[string]interface{} {
	return map[string]interface{}{
		"certificateArn":  c.CertificateARN,
		"certificateId":   c.CertificateID,
		"status":          c.Status,
		"certificateMode": c.CertificateMode,
		"creationDate":    c.CreationDate.Unix(),
	}
}

func certificateDetailResponse(c *iotstore.Certificate) map[string]interface{} {
	return map[string]interface{}{
		"certificateDescription": map[string]interface{}{
			"certificateArn":   c.CertificateARN,
			"certificateId":    c.CertificateID,
			"status":           c.Status,
			"certificatePem":   c.CertificatePEM,
			"lastModifiedDate": c.LastModifiedDate.Unix(),
			"creationDate":     c.CreationDate.Unix(),
		},
	}
}

func authorizerResponse(a *iotstore.Authorizer) map[string]interface{} {
	status := "ACTIVE"
	if !a.Status {
		status = "INACTIVE"
	}
	return map[string]interface{}{
		"authorizerName":        a.AuthorizerName,
		"authorizerArn":         a.AuthorizerARN,
		"authorizerFunctionArn": a.AuthorizerFunctionARN,
		"tokenKeyName":          a.TokenName,
		"tokenSignature":        a.TokenSignature,
		"status":                status,
		"enableCachingForHttp":  a.EnableCachingForHTTP,
		"creationDate":          a.CreationDate.Unix(),
		"lastModifiedDate":      a.LastModifiedDate.Unix(),
	}
}

func jobResponse(j *iotstore.Job) map[string]interface{} {
	return map[string]interface{}{
		"jobId":           j.JobID,
		"jobArn":          j.JobARN,
		"description":     j.Description,
		"status":          j.Status,
		"targetSelection": j.TargetSelection,
		"targets":         j.Targets,
		"createdAt":       j.CreatedAt.Unix(),
		"lastUpdatedAt":   j.LastUpdatedAt.Unix(),
	}
}

func policyResponse(p *iotstore.Policy) map[string]interface{} {
	return map[string]interface{}{
		"policyName":       p.PolicyName,
		"policyArn":        p.PolicyARN,
		"policyDocument":   p.PolicyDocument,
		"creationDate":     p.CreationDate.Unix(),
		"lastModifiedDate": p.LastModifiedDate.Unix(),
		"defaultVersionId": fmt.Sprintf("%d", p.Version),
	}
}

func roleAliasResponse(ra *iotstore.RoleAlias) map[string]interface{} {
	return map[string]interface{}{
		"roleAliasArn":              ra.RoleAliasARN,
		"roleAlias":                 ra.RoleAlias,
		"credentialDurationSeconds": ra.CredentialDurationSeconds,
		"roleArn":                   ra.RoleARN,
		"owner":                     ra.Owner,
		"creationDate":              ra.CreationDate.Unix(),
		"lastModifiedDate":          ra.LastModifiedDate.Unix(),
	}
}

func provisioningTemplateResponse(t *iotstore.ProvisioningTemplate) map[string]interface{} {
	resp := map[string]interface{}{
		"templateName":        t.TemplateName,
		"templateArn":         t.TemplateARN,
		"description":         t.Description,
		"enabled":             t.Enabled,
		"provisioningRoleArn": t.ProvisioningRoleARN,
		"creationDate":        t.CreationDate.Unix(),
		"lastModifiedDate":    t.LastModifiedDate.Unix(),
		"type":                t.Type,
	}
	if t.TemplateBody != "" {
		resp["templateBody"] = t.TemplateBody
	}
	return resp
}

// securityProfileIdentifierResponse renders the SecurityProfileIdentifier
// summary shape used by ListSecurityProfiles (AWS returns name/arn, not the
// full securityProfileName/securityProfileArn detail shape).
func securityProfileIdentifierResponse(name, arn string) map[string]interface{} {
	return map[string]interface{}{
		"name": name,
		"arn":  arn,
	}
}

func domainConfigResponse(dc *iotstore.DomainConfiguration) map[string]interface{} {
	return map[string]interface{}{
		"domainConfigurationName": dc.DomainConfigurationName,
		"domainConfigurationArn":  dc.DomainConfigurationARN,
		"serviceType":             dc.ServiceType,
	}
}

func domainConfigDetailResponse(dc *iotstore.DomainConfiguration) map[string]interface{} {
	resp := map[string]interface{}{
		"domainConfigurationName":   dc.DomainConfigurationName,
		"domainConfigurationArn":    dc.DomainConfigurationARN,
		"domainName":                dc.DomainName,
		"serverCertificateArns":     dc.ServerCertificateARNs,
		"serviceType":               dc.ServiceType,
		"domainConfigurationStatus": dc.DomainConfigurationStatus,
		"authenticationType":        dc.AuthenticationType,
		"applicationProtocol":       dc.ApplicationProtocol,
		"creationDate":              dc.CreationDate.Unix(),
		"lastModifiedDate":          dc.LastModifiedDate.Unix(),
	}
	if dc.AuthorizerConfig != "" {
		resp["authorizerConfig"] = dc.AuthorizerConfig
	}
	return resp
}

func violationEventResponse(e *iotstore.ViolationEvent) map[string]interface{} {
	resp := map[string]interface{}{
		"violationId":         e.ViolationID,
		"thingName":           e.ThingName,
		"securityProfileName": e.SecurityProfileName,
		"violationEventType":  e.ViolationEventType,
		"verificationState":   e.VerificationState,
		"violationEventTime":  e.ViolationEventTime.Unix(),
	}
	if e.VerificationStateDescription != "" {
		resp["verificationStateDescription"] = e.VerificationStateDescription
	}
	if e.Behavior != nil {
		resp["behavior"] = behaviorResponse(e.Behavior)
	}
	if e.MetricValue != nil {
		resp["metricValue"] = metricValueResponse(e.MetricValue)
	}
	return resp
}

func behaviorResponse(b *iotstore.Behavior) map[string]interface{} {
	resp := map[string]interface{}{
		"name": b.Name,
	}
	if b.Metric != "" {
		resp["metric"] = b.Metric
	}
	if b.MetricDimension != "" {
		resp["metricDimension"] = b.MetricDimension
	}
	if b.SuppressAlerts {
		resp["suppressAlerts"] = b.SuppressAlerts
	}
	if b.ExportMetric {
		resp["exportMetric"] = b.ExportMetric
	}
	return resp
}

func metricValueResponse(m *iotstore.MetricValue) map[string]interface{} {
	resp := map[string]interface{}{}
	if m.Count != 0 {
		resp["count"] = m.Count
	}
	if m.Number != 0 {
		resp["number"] = m.Number
	}
	if len(m.Cidrs) > 0 {
		resp["cidrs"] = m.Cidrs
	}
	if len(m.Ports) > 0 {
		resp["ports"] = m.Ports
	}
	if len(m.Numbers) > 0 {
		resp["numbers"] = m.Numbers
	}
	if len(m.Strings) > 0 {
		resp["strings"] = m.Strings
	}
	return resp
}
