package iot

import (
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"google.golang.org/protobuf/proto"

	"vorpalstacks/internal/pb/aws/iot"
	iotstore "vorpalstacks/internal/store/aws/iot"
	"vorpalstacks/internal/utils/timeutils"
)

// This file is the sole location permitted to import the iotstore package
// directly. It converts between store-level structs and proto response
// messages.

// getStoreFromHeaders resolves the region from gRPC-Web request headers and
// returns the IoT store for that region.
func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (iotstore.IotStoreInterface, error) {
	region := defaults.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

func toPbThingAttribute(t *iotstore.Thing) *iot.ThingAttribute {
	ta := &iot.ThingAttribute{
		Attributes:    t.Attributes,
		Thingarn:      t.ThingARN,
		Thingname:     t.ThingName,
		Thingtypename: t.ThingTypeName,
	}
	if t.Version != 0 {
		v := t.Version
		ta.Version = &v
	}
	return ta
}

func toPbDescribeThingResponse(t *iotstore.Thing) *iot.DescribeThingResponse {
	resp := &iot.DescribeThingResponse{
		Attributes:       t.Attributes,
		Thingarn:         t.ThingARN,
		Thingid:          t.ThingID,
		Thingname:        t.ThingName,
		Thingtypename:    t.ThingTypeName,
		Defaultclientid:  t.DefaultClientId,
		Billinggroupname: t.BillingGroupName,
	}
	if t.Version != 0 {
		v := t.Version
		resp.Version = &v
	}
	return resp
}

func toPbCreateThingResponse(t *iotstore.Thing) *iot.CreateThingResponse {
	return &iot.CreateThingResponse{
		Thingname: t.ThingName,
		Thingarn:  t.ThingARN,
		Thingid:   t.ThingID,
	}
}

func toPbPolicy(p *iotstore.Policy) *iot.Policy {
	return &iot.Policy{
		Policyname: p.PolicyName,
		Policyarn:  p.PolicyARN,
	}
}

func toPbCreatePolicyResponse(p *iotstore.Policy) *iot.CreatePolicyResponse {
	return &iot.CreatePolicyResponse{
		Policyname:      p.PolicyName,
		Policyarn:       p.PolicyARN,
		Policydocument:  p.PolicyDocument,
		Policyversionid: "1",
	}
}

func toPbGetPolicyResponse(p *iotstore.Policy) *iot.GetPolicyResponse {
	resp := &iot.GetPolicyResponse{
		Policyname:       p.PolicyName,
		Policyarn:        p.PolicyARN,
		Policydocument:   p.PolicyDocument,
		Defaultversionid: "1",
		Generationid:     "1",
	}
	if !p.CreationDate.IsZero() {
		resp.Creationdate = p.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !p.LastModifiedDate.IsZero() {
		resp.Lastmodifieddate = p.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}
	return resp
}

func certStatusToProto(status string) iot.CertificateStatus {
	switch status {
	case "ACTIVE":
		return iot.CertificateStatus_CERTIFICATE_STATUS_ACTIVE
	case "INACTIVE":
		return iot.CertificateStatus_CERTIFICATE_STATUS_INACTIVE
	case "REVOKED":
		return iot.CertificateStatus_CERTIFICATE_STATUS_REVOKED
	case "PENDING_ACTIVATION":
		return iot.CertificateStatus_CERTIFICATE_STATUS_PENDING_ACTIVATION
	case "PENDING_TRANSFER":
		return iot.CertificateStatus_CERTIFICATE_STATUS_PENDING_TRANSFER
	case "REGISTER_INACTIVE":
		return iot.CertificateStatus_CERTIFICATE_STATUS_REGISTER_INACTIVE
	}
	return iot.CertificateStatus_CERTIFICATE_STATUS_INACTIVE
}

func certStatusFromProto(s iot.CertificateStatus) string {
	switch s {
	case iot.CertificateStatus_CERTIFICATE_STATUS_ACTIVE:
		return "ACTIVE"
	case iot.CertificateStatus_CERTIFICATE_STATUS_INACTIVE:
		return "INACTIVE"
	case iot.CertificateStatus_CERTIFICATE_STATUS_REVOKED:
		return "REVOKED"
	case iot.CertificateStatus_CERTIFICATE_STATUS_PENDING_ACTIVATION:
		return "PENDING_ACTIVATION"
	case iot.CertificateStatus_CERTIFICATE_STATUS_PENDING_TRANSFER:
		return "PENDING_TRANSFER"
	case iot.CertificateStatus_CERTIFICATE_STATUS_REGISTER_INACTIVE:
		return "REGISTER_INACTIVE"
	}
	return "INACTIVE"
}

func certModeToProto(mode string) iot.CertificateMode {
	if mode == "SNI_ONLY" {
		return iot.CertificateMode_CERTIFICATE_MODE_SNI_ONLY
	}
	return iot.CertificateMode_CERTIFICATE_MODE_DEFAULT
}

func toPbCertificate(c *iotstore.Certificate) *iot.Certificate {
	cert := &iot.Certificate{
		Certificatearn:  c.CertificateARN,
		Certificateid:   c.CertificateID,
		Certificatemode: certModeToProto(c.CertificateMode),
		Status:          certStatusToProto(c.Status),
	}
	if !c.CreationDate.IsZero() {
		cert.Creationdate = c.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	return cert
}

func toPbCertificateDescription(c *iotstore.Certificate) *iot.CertificateDescription {
	desc := &iot.CertificateDescription{
		Certificatearn:  c.CertificateARN,
		Certificateid:   c.CertificateID,
		Certificatepem:  c.CertificatePEM,
		Cacertificateid: c.CaCertificateID,
		Status:          certStatusToProto(c.Status),
		Certificatemode: certModeToProto(c.CertificateMode),
		Generationid:    "1",
	}
	cv := int32(1)
	desc.Customerversion = &cv
	if !c.CreationDate.IsZero() {
		desc.Creationdate = c.CreationDate.Format(timeutils.ISO8601UTCFormat)
	}
	if !c.LastModifiedDate.IsZero() {
		desc.Lastmodifieddate = c.LastModifiedDate.Format(timeutils.ISO8601UTCFormat)
	}
	return desc
}

func toPbTopicRuleListItem(r *iotstore.TopicRule) *iot.TopicRuleListItem {
	item := &iot.TopicRuleListItem{
		Rulearn:      r.ARN,
		Rulename:     r.RuleName,
		Topicpattern: r.TopicPattern,
		Createdat:    r.CreatedAt,
	}
	if r.RuleDisabled {
		item.Ruledisabled = proto.Bool(true)
	}
	return item
}
