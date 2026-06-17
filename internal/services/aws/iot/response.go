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
	return map[string]interface{}{
		"thingTypeName":       tt.ThingTypeName,
		"thingTypeArn":        tt.ThingTypeARN,
		"thingTypeId":         tt.ThingTypeID,
		"description":         tt.Description,
		"thingTypeProperties": map[string]interface{}{"thingTypeDescription": tt.Description},
		"version":             tt.Version,
		"creationDate":        tt.CreationDate.Unix(),
		"lastModifiedDate":    tt.LastModifiedDate.Unix(),
	}
}

func thingTypeDescribeResponse(tt *iotstore.ThingType) map[string]interface{} {
	props := map[string]interface{}{"thingTypeDescription": tt.Description}
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
	return map[string]interface{}{
		"templateName":        t.TemplateName,
		"templateArn":         t.TemplateARN,
		"description":         t.Description,
		"enabled":             t.Enabled,
		"provisioningRoleArn": t.ProvisioningRoleARN,
		"creationDate":        t.CreationDate.Unix(),
		"lastModifiedDate":    t.LastModifiedDate.Unix(),
	}
}

func securityProfileResponse(sp *iotstore.SecurityProfile) map[string]interface{} {
	return map[string]interface{}{
		"securityProfileName":        sp.SecurityProfileName,
		"securityProfileArn":         sp.SecurityProfileARN,
		"securityProfileDescription": sp.SecurityProfileDescription,
	}
}

func detectorModelResponse(m *iotstore.DetectorModel) map[string]interface{} {
	return map[string]interface{}{
		"detectorModelName": m.DetectorModelName,
		"detectorModelArn":  m.DetectorModelARN,
		"roleArn":           m.RoleARN,
	}
}

func detectorModelDetailResponse(m *iotstore.DetectorModel) map[string]interface{} {
	return map[string]interface{}{
		"detectorModelName":        m.DetectorModelName,
		"detectorModelArn":         m.DetectorModelARN,
		"detectorModelDescription": m.DetectorModelDescription,
		"roleArn":                  m.RoleARN,
		"detectorModelDefinition":  m.DetectorModelDefinition,
		"status":                   m.Status,
		"key":                      m.Key,
		"evaluationMethod":         m.EvaluationMethod,
		"detectorModelVersion":     m.DetectorModelVersion,
		"creationDate":             m.CreationDate.Unix(),
		"lastModifiedDate":         m.LastModifiedDate.Unix(),
	}
}

func inputResponse(inp *iotstore.Input) map[string]interface{} {
	return map[string]interface{}{
		"inputName": inp.InputName,
		"inputArn":  inp.InputARN,
	}
}

func inputDetailResponse(inp *iotstore.Input) map[string]interface{} {
	return map[string]interface{}{
		"inputName":        inp.InputName,
		"inputArn":         inp.InputARN,
		"inputDescription": inp.InputDescription,
		"inputDefinition":  inp.InputDefinition,
		"status":           inp.Status,
		"creationDate":     inp.CreationDate.Unix(),
		"lastModifiedDate": inp.LastModifiedDate.Unix(),
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
