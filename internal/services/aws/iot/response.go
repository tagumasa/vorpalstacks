package iot

import (
	"encoding/json"
	"fmt"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

func thingResponse(t *iotstore.Thing) map[string]interface{} {
	return map[string]interface{}{
		"thingName":     t.ThingName,
		"thingArn":      t.ThingARN,
		"thingId":       t.ThingID,
		"thingTypeName": t.ThingTypeName,
		"attributes":    ensureMap(t.Attributes),
		"version":       t.Version,
	}
}

func thingDescribeResponse(t *iotstore.Thing) map[string]interface{} {
	resp := map[string]interface{}{
		"thingName":     t.ThingName,
		"thingArn":      t.ThingARN,
		"thingId":       t.ThingID,
		"thingTypeName": t.ThingTypeName,
		"attributes":    ensureMap(t.Attributes),
		"version":       t.Version,
	}
	if t.DefaultClientId != "" {
		resp["defaultClientId"] = t.DefaultClientId
	} else {
		resp["defaultClientId"] = t.ThingName
	}
	return resp
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
	resp := map[string]interface{}{
		"authorizerName":        a.AuthorizerName,
		"authorizerArn":         a.AuthorizerARN,
		"authorizerFunctionArn": a.AuthorizerFunctionARN,
		"tokenKeyName":          a.TokenName,
		"status":                status,
		"enableCachingForHttp":  a.EnableCachingForHTTP,
		"signingDisabled":       a.SigningDisabled,
		"creationDate":          a.CreationDate.Unix(),
		"lastModifiedDate":      a.LastModifiedDate.Unix(),
	}
	if len(a.TokenSigningPublicKeys) > 0 {
		resp["tokenSigningPublicKeys"] = a.TokenSigningPublicKeys
	}
	return resp
}

func jobResponse(j *iotstore.Job) map[string]interface{} {
	resp := map[string]interface{}{
		"jobId":           j.JobID,
		"jobArn":          j.JobARN,
		"description":     j.Description,
		"status":          j.Status,
		"targetSelection": j.TargetSelection,
		"targets":         j.Targets,
		"createdAt":       j.CreatedAt.Unix(),
		"lastUpdatedAt":   j.LastUpdatedAt.Unix(),
		"forceCanceled":   j.ForceCanceledFlag,
	}
	if j.ReasonCode != "" {
		resp["reasonCode"] = j.ReasonCode
	}
	if j.Comment != "" {
		resp["comment"] = j.Comment
	}
	if j.NamespaceID != "" {
		resp["namespaceId"] = j.NamespaceID
	}
	if j.CompletedAt != "" {
		resp["completedAt"] = j.CompletedAt
	}
	if j.JobTemplateARN != "" {
		resp["jobTemplateArn"] = j.JobTemplateARN
	}
	if j.IsConcurrent {
		resp["isConcurrent"] = j.IsConcurrent
	}
	if v := rawJSONOrOmit(j.PresignedUrlConfig); v != nil {
		resp["presignedUrlConfig"] = v
	}
	if v := rawJSONOrOmit(j.JobExecutionsRolloutConfig); v != nil {
		resp["jobExecutionsRolloutConfig"] = v
	}
	if v := rawJSONOrOmit(j.AbortConfig); v != nil {
		resp["abortConfig"] = v
	}
	if v := rawJSONOrOmit(j.TimeoutConfig); v != nil {
		resp["timeoutConfig"] = v
	}
	if v := rawJSONOrOmit(j.JobExecutionsRetryConfig); v != nil {
		resp["jobExecutionsRetryConfig"] = v
	}
	if v := rawJSONOrOmit(j.DocumentParameters); v != nil {
		resp["documentParameters"] = v
	}
	if v := rawJSONOrOmit(j.SchedulingConfig); v != nil {
		resp["schedulingConfig"] = v
	}
	if v := rawJSONOrOmit(j.ScheduledJobRollouts); v != nil {
		resp["scheduledJobRollouts"] = v
	}
	if v := rawJSONOrOmit(j.DestinationPackageVersions); v != nil {
		resp["destinationPackageVersions"] = v
	}
	return resp
}

// rawJSONOrOmit validates that s contains syntactically valid JSON before
// returning it as json.RawMessage. If s is empty or invalid JSON, it returns
// nil so the caller can omit the field from the response map entirely.
// Without this guard, a non-JSON string stored in a config field would
// corrupt the entire response serialisation.
func rawJSONOrOmit(s string) interface{} {
	if s == "" {
		return nil
	}
	if !json.Valid([]byte(s)) {
		return nil
	}
	return json.RawMessage(s)
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

// provisioningTemplateSummaryResponse renders the
// ProvisioningTemplateSummary shape used by ListProvisioningTemplates.
func provisioningTemplateSummaryResponse(t *iotstore.ProvisioningTemplate) map[string]interface{} {
	return map[string]interface{}{
		"templateArn":      t.TemplateARN,
		"templateName":     t.TemplateName,
		"description":      t.Description,
		"creationDate":     t.CreationDate.Unix(),
		"lastModifiedDate": t.LastModifiedDate.Unix(),
		"enabled":          t.Enabled,
		"type":             t.Type,
	}
}

// provisioningTemplateDetailResponse renders the
// DescribeProvisioningTemplate output: the model's eleven members including
// defaultVersionId and the preProvisioningHook structure.
func provisioningTemplateDetailResponse(t *iotstore.ProvisioningTemplate) map[string]interface{} {
	resp := map[string]interface{}{
		"templateArn":         t.TemplateARN,
		"templateName":        t.TemplateName,
		"description":         t.Description,
		"creationDate":        t.CreationDate.Unix(),
		"lastModifiedDate":    t.LastModifiedDate.Unix(),
		"defaultVersionId":    t.DefaultVersionID,
		"enabled":             t.Enabled,
		"provisioningRoleArn": t.ProvisioningRoleARN,
		"type":                t.Type,
	}
	if t.TemplateBody != "" {
		resp["templateBody"] = t.TemplateBody
	}
	if t.PreProvisioningHook != "" {
		var hook interface{}
		if err := json.Unmarshal([]byte(t.PreProvisioningHook), &hook); err == nil {
			resp["preProvisioningHook"] = hook
		} else {
			resp["preProvisioningHook"] = t.PreProvisioningHook
		}
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

// domainConfigDetailResponse renders the DescribeDomainConfiguration
// output: exactly the API model's member names (serverCertificates
// summaries, the authorizerConfig structure, lastStatusChangeDate).
func domainConfigDetailResponse(dc *iotstore.DomainConfiguration) map[string]interface{} {
	resp := map[string]interface{}{
		"domainConfigurationName":   dc.DomainConfigurationName,
		"domainConfigurationArn":    dc.DomainConfigurationARN,
		"domainName":                dc.DomainName,
		"serviceType":               dc.ServiceType,
		"domainConfigurationStatus": dc.DomainConfigurationStatus,
		"lastStatusChangeDate":      dc.LastModifiedDate.Unix(),
	}
	if dc.AuthenticationType != "" {
		resp["authenticationType"] = dc.AuthenticationType
	}
	if dc.ApplicationProtocol != "" {
		resp["applicationProtocol"] = dc.ApplicationProtocol
	}
	if len(dc.ServerCertificateARNs) > 0 {
		summaries := make([]interface{}, 0, len(dc.ServerCertificateARNs))
		for _, arn := range dc.ServerCertificateARNs {
			summaries = append(summaries, map[string]interface{}{
				"serverCertificateArn": arn,
			})
		}
		resp["serverCertificates"] = summaries
	}
	if cfg := domainConfigAuthorizerConfigValue(dc.AuthorizerConfig); cfg != nil {
		resp["authorizerConfig"] = cfg
	}
	return resp
}

// activeViolationResponse renders the ActiveViolation shape used by
// ListActiveViolations: the most-recent violation members, without the
// violationEventType/violationEventTime/metricValue members that belong to
// the ViolationEvent shape.
func activeViolationResponse(e *iotstore.ViolationEvent) map[string]interface{} {
	resp := map[string]interface{}{
		"violationId":         e.ViolationID,
		"thingName":           e.ThingName,
		"securityProfileName": e.SecurityProfileName,
		"lastViolationTime":   e.ViolationEventTime.Unix(),
		"violationStartTime":  e.ViolationEventTime.Unix(),
	}
	if e.VerificationState != "" {
		resp["verificationState"] = e.VerificationState
	}
	if e.VerificationStateDescription != "" {
		resp["verificationStateDescription"] = e.VerificationStateDescription
	}
	if e.Behavior != nil {
		resp["behavior"] = behaviorResponse(e.Behavior)
	}
	if e.MetricValue != nil {
		resp["lastViolationValue"] = metricValueResponse(e.MetricValue)
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
