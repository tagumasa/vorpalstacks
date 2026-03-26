package route53

import (
	"context"
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	tagutil "vorpalstacks/internal/services/aws/common/tags"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// CreateHostedZone creates a new hosted zone in Route 53.
func (s *Route53Service) CreateHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, NewAPIError("InvalidInput", "Name is required", 400)
	}

	callerRef := request.GetStringParam(req.Parameters, "CallerReference")
	if callerRef == "" {
		callerRef = fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s-%d", name, time.Now().UnixNano()))))
	}

	var comment string
	privateZone := request.GetBoolParam(req.Parameters, "PrivateZone")
	if hzConfig := request.GetMapParam(req.Parameters, "HostedZoneConfig"); hzConfig != nil {
		comment = request.GetStringParam(hzConfig, "Comment")
		privateZone = request.GetBoolParam(hzConfig, "PrivateZone")
	}

	vpcMap := request.GetMapParam(req.Parameters, "VPC")
	var vpcs []*route53store.VPC
	if vpcMap != nil {
		if vpc := parseVPC(vpcMap); vpc != nil {
			vpcs = append(vpcs, vpc)
		}
	}

	zone := &route53store.HostedZone{
		ID:                     generateHostedZoneId(),
		Name:                   normalizeZoneNameForCreate(name),
		CallerReference:        callerRef,
		Config:                 &route53store.HostedZoneConfig{Comment: comment, PrivateZone: privateZone},
		ResourceRecordSetCount: 0,
		Private:                privateZone,
		VPCs:                   vpcs,
		Region:                 reqCtx.GetRegion(),
		AccountID:              s.accountID,
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.HostedZones().Create(zone); err != nil {
		if route53store.IsAlreadyExists(err) {
			return nil, NewAPIError("HostedZoneAlreadyExists", fmt.Sprintf("Hosted zone already exists: %s", name), 400)
		}
		return nil, mapStoreError(err)
	}

	if tags := tagutil.ParseTags(req.Parameters, "HostedZoneTags"); len(tags) > 0 {
		resourceKey := "hostedzone/" + zone.ID
		if err := st.Tags().TagResource(resourceKey, tags); err != nil {
			return nil, NewAPIError("TagResource", err.Error(), 500)
		}
	}

	changeId := generateChangeId()
	now := time.Now()
	if err := st.Changes().Create(&route53store.ChangeInfo{
		ID:          changeId,
		Status:      "INSYNC",
		SubmittedAt: now,
		Comment:     comment,
	}); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"HostedZone": s.hostedZoneToResponse(zone),
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + changeId,
			"Status":      "INSYNC",
			"SubmittedAt": now.Format(time.RFC3339),
		},
		"DelegationSet": map[string]interface{}{
			"Id":              "/delegationset/" + generateDelegationSetId(),
			"CallerReference": callerRef,
			"NameServers": protocol.XMLElements{
				ElementName: "NameServer",
				Items: []interface{}{
					"ns-1.awsdns.com.",
					"ns-2.awsdns.net.",
					"ns-3.awsdns.org.",
					"ns-4.awsdns.co.uk.",
				},
			},
		},
	}, nil
}

// GetHostedZone returns details about a hosted zone by its ID.
func (s *Route53Service) GetHostedZone(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id, err := extractHostedZoneId(req.Parameters, "Id")
	if err != nil {
		return nil, err
	}

	zone, err := s.getHostedZoneById(reqCtx, id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"HostedZone": s.hostedZoneToResponse(zone),
		"DelegationSet": map[string]interface{}{
			"NameServers": protocol.XMLElements{
				ElementName: "NameServer",
				Items: []interface{}{
					"ns-1.awsdns.com.",
					"ns-2.awsdns.net.",
					"ns-3.awsdns.org.",
					"ns-4.awsdns.co.uk.",
				},
			},
		},
		"VPCs": []interface{}{},
	}, nil
}

// ListHostedZones returns a list of hosted zones with pagination support.
func (s *Route53Service) ListHostedZones(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	result, err := st.HostedZones().List(marker, maxItems)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return s.buildHostedZonesListResponse(result.HostedZones, result.IsTruncated, result.Marker, maxItems), nil
}

// ListHostedZonesByName returns hosted zones sorted by name with optional DNS name filter.
func (s *Route53Service) ListHostedZonesByName(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	dnsName := request.GetStringParam(req.Parameters, "DNSName")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := request.GetIntParam(req.Parameters, "MaxItems")

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if dnsName != "" {
		zone, err := st.HostedZones().GetByName(dnsName)
		if err != nil {
			if route53store.IsNotFound(err) {
				return s.buildHostedZonesListResponse(nil, false, "", maxItems), nil
			}
			return nil, mapStoreError(err)
		}
		return s.buildHostedZonesListResponse([]*route53store.HostedZone{zone}, false, "", maxItems), nil
	}

	result, err := st.HostedZones().List(marker, maxItems)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return s.buildHostedZonesListResponse(result.HostedZones, result.IsTruncated, result.Marker, maxItems), nil
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
	recordSets, err := st.RecordSets().List(id)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if len(recordSets) > 0 {
		return nil, NewAPIError("HostedZoneNotEmpty", "The hosted zone must be empty before it can be deleted.", 400)
	}

	if err := st.HostedZones().Delete(id); err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchHostedZoneError(id)
		}
		return nil, mapStoreError(err)
	}

	changeId := generateChangeId()
	now := time.Now()
	if err := st.Changes().Create(&route53store.ChangeInfo{
		ID:          changeId,
		Status:      "INSYNC",
		SubmittedAt: now,
	}); err != nil {
		return nil, mapStoreError(err)
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + changeId,
			"Status":      "INSYNC",
			"SubmittedAt": now.Format(time.RFC3339),
		},
	}, nil
}

// UpdateHostedZoneComment updates the comment for a hosted zone.
func (s *Route53Service) UpdateHostedZoneComment(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id, err := extractHostedZoneId(req.Parameters, "Id")
	if err != nil {
		return nil, err
	}

	zone, err := s.getHostedZoneById(reqCtx, id)
	if err != nil {
		return nil, err
	}

	comment := request.GetStringParam(req.Parameters, "Comment")
	if zone.Config == nil {
		zone.Config = &route53store.HostedZoneConfig{}
	}
	zone.Config.Comment = comment

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.HostedZones().Update(zone); err != nil {
		return nil, mapStoreError(err)
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

	if len(zone.VPCs) > 0 {
		vpcs := make([]map[string]string, len(zone.VPCs))
		for i, vpc := range zone.VPCs {
			vpcs[i] = map[string]string{
				"VPCRegion": vpc.VPCRegion,
				"VPCId":     vpc.VPCID,
			}
		}
		result["VPCs"] = vpcs
	}

	return result
}

func normalizeZoneNameForCreate(name string) string {
	name = strings.ToLower(name)
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	return name
}
