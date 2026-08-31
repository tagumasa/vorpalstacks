package route53

import (
	"context"
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	tagutil "vorpalstacks/internal/common/tags"
	route53store "vorpalstacks/internal/store/aws/route53"
)

type delegationSetResponse struct {
	ID          string               `json:"id"`
	NameServers protocol.XMLElements `json:"nameServers"`
}

// CreateHostedZone creates a new hosted zone in Route 53.
func (s *Route53Service) CreateHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")

	callerRef := request.GetStringParam(req.Parameters, "CallerReference")

	var comment string
	privateZone := request.GetBoolParam(req.Parameters, "PrivateZone")
	if hzConfig := request.GetMapParam(req.Parameters, "HostedZoneConfig"); hzConfig != nil {
		comment = request.GetStringParam(hzConfig, "Comment")
		privateZone = request.GetBoolParam(hzConfig, "PrivateZone")
	}

	vpcID, vpcRegion := "", ""
	vpcMap := request.GetMapParam(req.Parameters, "VPC")
	if vpcMap != nil {
		if vpc := parseVPC(vpcMap); vpc != nil {
			vpcID = vpc.VPCID
			vpcRegion = vpc.VPCRegion
			if err := s.validateVPC(ctx, reqCtx.GetRegion(), vpc.VPCID, vpc.VPCRegion); err != nil {
				return nil, awserrors.NewAWSError("InvalidVPCId",
					fmt.Sprintf("The VPC %s in region %s does not exist", vpc.VPCID, vpc.VPCRegion), 400)
			}
		}
	}

	delegationSetID := ""
	if dsID := request.GetStringParam(req.Parameters, "DelegationSetId"); dsID != "" {
		delegationSetID = strings.TrimPrefix(dsID, "/delegationset/")
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.createHostedZoneCore(st, CreateHostedZoneInput{
		Name:            name,
		CallerReference: callerRef,
		Comment:         comment,
		PrivateZone:     privateZone,
		VPCID:           vpcID,
		VPCRegion:       vpcRegion,
		DelegationSetID: delegationSetID,
		Region:          reqCtx.GetRegion(),
		Tags:            tagutil.ParseTags(req.Parameters, "HostedZoneTags"),
	})
	if err != nil {
		return nil, err
	}

	zone := result.Zone
	resp := map[string]interface{}{
		"HostedZone":    s.hostedZoneToResponse(zone),
		"DelegationSet": buildDelegationSetResponse(zone.NameServers, zone.DelegationSetID),
	}
	if zone.ChangeID != "" {
		resp["ChangeInfo"] = map[string]interface{}{
			"Id":          "/change/" + zone.ChangeID,
			"Status":      "INSYNC",
			"SubmittedAt": zone.CreatedAt.Format(time.RFC3339),
		}
	}
	return resp, nil
}

// GetHostedZone returns details about a hosted zone by its ID.
func (s *Route53Service) GetHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id, err := extractHostedZoneId(req.Parameters, "Id")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	zone, err := getHostedZoneCore(st, id)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"HostedZone": s.hostedZoneToResponse(zone),
	}

	if len(zone.NameServers) > 0 {
		result["DelegationSet"] = buildDelegationSetResponse(zone.NameServers, zone.DelegationSetID)
	}

	if len(zone.VPCs) > 0 {
		vpcs := make([]interface{}, len(zone.VPCs))
		for i, vpc := range zone.VPCs {
			vpcs[i] = map[string]interface{}{
				"VPCRegion": vpc.VPCRegion,
				"VPCId":     vpc.VPCID,
			}
		}
		result["VPCs"] = protocol.XMLElements{ElementName: "VPC", Items: vpcs}
	}

	return result, nil
}

// ListHostedZones returns a list of hosted zones with pagination support.
func (s *Route53Service) ListHostedZones(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	delegationSetId := request.GetStringParam(req.Parameters, "DelegationSetId")
	if delegationSetId != "" {
		delegationSetId = strings.TrimPrefix(delegationSetId, "/delegationset/")
	}
	hostedZoneType := request.GetStringParam(req.Parameters, "HostedZoneType")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.listHostedZonesCore(st, ListHostedZonesInput{
		Marker:          marker,
		MaxItems:        maxItems,
		DelegationSetID: delegationSetId,
		HostedZoneType:  hostedZoneType,
	})
	if err != nil {
		return nil, err
	}

	return s.buildHostedZonesListResponse(result.HostedZones, result.IsTruncated, marker, result.NextMarker, maxItems), nil
}

// ListHostedZonesByName returns hosted zones sorted by name with optional DNS name filter.
func (s *Route53Service) ListHostedZonesByName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	dnsName := request.GetStringParam(req.Parameters, "DNSName")
	hostedZoneIdMarker := request.GetStringParam(req.Parameters, "HostedZoneId")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")
	if maxItems <= 0 {
		maxItems = 100
	}
	if dnsName != "" {
		dnsName = strings.ToLower(dnsName)
		if !strings.HasSuffix(dnsName, ".") {
			dnsName = dnsName + "."
		}
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := listHostedZonesByNameCore(st, ListHostedZonesByNameInput{
		DNSName:            dnsName,
		HostedZoneIdMarker: hostedZoneIdMarker,
		MaxItems:           maxItems,
	})
	if err != nil {
		return nil, err
	}
	filtered := result.HostedZones
	isTruncated := result.IsTruncated

	resp := map[string]interface{}{
		"HostedZones": protocol.XMLElements{ElementName: "HostedZone", Items: func() []interface{} {
			items := make([]interface{}, len(filtered))
			for i, z := range filtered {
				items[i] = s.hostedZoneToResponse(z)
			}
			return items
		}()},
		"IsTruncated": isTruncated,
		"MaxItems":    maxItems,
	}

	if dnsName != "" {
		resp["DNSName"] = dnsName
	}
	if hostedZoneIdMarker != "" {
		resp["HostedZoneId"] = hostedZoneIdMarker
	}

	if isTruncated && len(filtered) > 0 {
		lastZone := filtered[len(filtered)-1]
		resp["NextDNSName"] = lastZone.Name
		resp["NextHostedZoneId"] = lastZone.ID
	}

	return resp, nil
}

// DeleteHostedZone deletes a hosted zone by its ID. The zone must be empty.
func (s *Route53Service) DeleteHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id, err := extractHostedZoneId(req.Parameters, "Id")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.deleteHostedZoneCore(st, DeleteHostedZoneInput{Id: id})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + result.ChangeId,
			"Status":      result.Status,
			"SubmittedAt": result.SubmittedAt.Format(time.RFC3339),
		},
	}, nil
}

// UpdateHostedZoneComment updates the comment for a hosted zone.
func (s *Route53Service) UpdateHostedZoneComment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id, err := extractHostedZoneId(req.Parameters, "Id")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	zone, err := updateHostedZoneCommentCore(st, UpdateHostedZoneCommentInput{
		Id:      id,
		Comment: request.GetStringParam(req.Parameters, "Comment"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"HostedZone": s.hostedZoneToResponse(zone),
	}, nil
}

func (s *Route53Service) hostedZoneToResponse(zone *route53store.HostedZone) map[string]interface{} {
	result := map[string]interface{}{
		"Id":                     "/hostedzone/" + zone.ID,
		"Name":                   zone.Name,
		"CallerReference":        zone.CallerReference,
		"ResourceRecordSetCount": zone.ResourceRecordSetCount,
	}

	if zone.Config != nil {
		config := map[string]interface{}{
			"PrivateZone": zone.Config.PrivateZone,
		}
		if zone.Config.Comment != "" {
			config["Comment"] = zone.Config.Comment
		}
		result["Config"] = config
	}

	return result
}

// GetDNSSEC retrieves the DNSSEC signing status and configuration for a hosted zone.
func (s *Route53Service) GetDNSSEC(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Status": map[string]interface{}{
			"ServeSignature": "NOT_SIGNING",
			"StatusMessage":  "",
		},
		"KeySigningKeys": []interface{}{},
	}, nil
}
