package route53

import (
	"context"
	"strings"
	"time"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	route53store "vorpalstacks/internal/store/aws/route53"
)

func normalizeRecordName(name string) string {
	name = strings.ToLower(name)
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	return name
}

// ChangeResourceRecordSets creates, updates, or deletes resource record sets in a hosted zone.
func (s *Route53Service) ChangeResourceRecordSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	hostedZoneId, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := changeResourceRecordSetsCore(st, ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZoneId,
		ChangeBatch:  request.GetMapParam(req.Parameters, "ChangeBatch"),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + result.ChangeId,
			"Status":      result.Status,
			"SubmittedAt": result.SubmittedAt.Format(time.RFC3339),
			"Comment":     result.Comment,
		},
	}, nil
}

// ListResourceRecordSets returns resource record sets in a hosted zone with pagination and filtering support.
func (s *Route53Service) ListResourceRecordSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	hostedZoneId, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	startRecordName := request.GetStringParam(req.Parameters, "StartRecordName")
	if startRecordName == "" {
		startRecordName = request.GetStringParam(req.Parameters, "name")
	}
	startRecordType := request.GetStringParam(req.Parameters, "StartRecordType")
	if startRecordType == "" {
		startRecordType = request.GetStringParam(req.Parameters, "type")
	}
	startRecordIdentifier := request.GetStringParam(req.Parameters, "StartRecordIdentifier")
	maxItems := int(request.GetIntParam(req.Parameters, "MaxItems"))
	if maxItems == 0 {
		maxItems = int(request.GetIntParam(req.Parameters, "maxitems"))
	}
	if maxItems <= 0 {
		maxItems = 300
	}

	if startRecordName != "" {
		startRecordName = normalizeRecordName(startRecordName)
	}
	startRecordType = strings.ToUpper(startRecordType)

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := listResourceRecordSetsCore(st, ListResourceRecordSetsInput{
		HostedZoneId:          hostedZoneId,
		StartRecordName:       startRecordName,
		StartRecordType:       startRecordType,
		StartRecordIdentifier: startRecordIdentifier,
		MaxItems:              maxItems,
	})
	if err != nil {
		return nil, err
	}

	records := make([]interface{}, len(result.RecordSets))
	for i, rs := range result.RecordSets {
		records[i] = s.recordSetToResponse(rs)
	}

	resp := map[string]interface{}{
		"ResourceRecordSets": protocol.XMLElements{ElementName: "ResourceRecordSet", Items: records},
		"IsTruncated":        result.IsTruncated,
		"MaxItems":           result.MaxItems,
	}

	if result.IsTruncated {
		resp["NextRecordName"] = result.NextRecordName
		resp["NextRecordType"] = result.NextRecordType
		if result.NextRecordIdentifier != "" {
			resp["NextRecordIdentifier"] = result.NextRecordIdentifier
		}
	}

	return resp, nil
}

// GetChange returns the status of a change batch request.
func (s *Route53Service) GetChange(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id, err := extractChangeId(req.Parameters, "Id")
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	change, err := getChangeCore(st, id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + change.ID,
			"Status":      change.Status,
			"SubmittedAt": change.SubmittedAt.Format(time.RFC3339),
			"Comment":     change.Comment,
		},
	}, nil
}

func (s *Route53Service) recordSetToResponse(rs *route53store.ResourceRecordSet) map[string]interface{} {
	result := map[string]interface{}{
		"Name": rs.Name,
		"Type": rs.Type,
	}

	if rs.AliasTarget == nil {
		result["TTL"] = rs.TTL
	}

	if len(rs.ResourceRecords) > 0 {
		records := make([]interface{}, len(rs.ResourceRecords))
		for i, r := range rs.ResourceRecords {
			records[i] = map[string]interface{}{"Value": r.Value}
		}
		result["ResourceRecords"] = protocol.XMLElements{ElementName: "ResourceRecord", Items: records}
	}

	if rs.AliasTarget != nil {
		result["AliasTarget"] = map[string]interface{}{
			"HostedZoneId":         rs.AliasTarget.HostedZoneID,
			"DNSName":              rs.AliasTarget.DNSName,
			"EvaluateTargetHealth": rs.AliasTarget.EvaluateTargetHealth,
		}
	}

	if rs.SetIdentifier != "" {
		result["SetIdentifier"] = rs.SetIdentifier
	}

	if rs.SetIdentifier != "" && rs.Region == "" && rs.Failover == "" && rs.GeoLocation == nil && !rs.MultiValueAnswer {
		result["Weight"] = rs.Weight
	}

	if rs.Region != "" {
		result["Region"] = rs.Region
	}

	if rs.Failover != "" {
		result["Failover"] = rs.Failover
	}

	if rs.MultiValueAnswer {
		result["MultiValueAnswer"] = true
	}

	if rs.HealthCheckID != "" {
		result["HealthCheckId"] = rs.HealthCheckID
	}

	if rs.TrafficPolicyInstanceID != "" {
		result["TrafficPolicyInstanceId"] = rs.TrafficPolicyInstanceID
	}

	if rs.GeoLocation != nil {
		result["GeoLocation"] = map[string]interface{}{
			"ContinentCode":   rs.GeoLocation.ContinentCode,
			"CountryCode":     rs.GeoLocation.CountryCode,
			"SubdivisionCode": rs.GeoLocation.SubdivisionCode,
		}
	}

	if rs.CidrRoutingConfig != nil {
		result["CidrRoutingConfig"] = map[string]interface{}{
			"CollectionId": rs.CidrRoutingConfig.CollectionId,
			"LocationName": rs.CidrRoutingConfig.LocationName,
		}
	}

	if rs.GeoProximityLocation != nil {
		gp := map[string]interface{}{}
		if rs.GeoProximityLocation.AWSRegion != "" {
			gp["AWSRegion"] = rs.GeoProximityLocation.AWSRegion
		}
		if rs.GeoProximityLocation.LocalZoneGroup != "" {
			gp["LocalZoneGroup"] = rs.GeoProximityLocation.LocalZoneGroup
		}
		if rs.GeoProximityLocation.Coordinates != nil {
			gp["Coordinates"] = map[string]interface{}{
				"Latitude":  rs.GeoProximityLocation.Coordinates.Latitude,
				"Longitude": rs.GeoProximityLocation.Coordinates.Longitude,
			}
		}
		if rs.GeoProximityLocation.Bias != 0 {
			gp["Bias"] = rs.GeoProximityLocation.Bias
		}
		result["GeoProximityLocation"] = gp
	}

	return result
}
