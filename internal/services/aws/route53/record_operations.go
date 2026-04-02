package route53

import (
	"context"
	"fmt"
	"strings"
	"time"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
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

	zone, err := s.getHostedZoneById(reqCtx, hostedZoneId)
	if err != nil {
		return nil, err
	}

	changeBatch := request.GetMapParam(req.Parameters, "ChangeBatch")
	if changeBatch == nil {
		return nil, NewAPIError("InvalidInput", "ChangeBatch is required", 400)
	}

	var changesList []interface{}
	switch c := changeBatch["Changes"].(type) {
	case []interface{}:
		changesList = c
	case map[string]interface{}:
		if changeArr, ok := c["Change"].([]interface{}); ok {
			changesList = changeArr
		} else if changeMap, ok := c["Change"].(map[string]interface{}); ok {
			changesList = []interface{}{changeMap}
		}
	default:
		return nil, NewAPIError("InvalidInput", "Changes must be an array", 400)
	}

	if len(changesList) == 0 {
		return nil, NewAPIError("InvalidInput", "Changes are required", 400)
	}

	changeId := generateChangeId()
	change := &route53store.ChangeInfo{
		ID:          changeId,
		Status:      "PENDING",
		SubmittedAt: time.Now(),
	}
	if cmt, ok := changeBatch["Comment"].(string); ok {
		change.Comment = cmt
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := st.Changes().Create(change); err != nil {
		return nil, NewAPIError("CreateChange", fmt.Sprintf("Failed to create change: %v", err), 500)
	}

	var appliedChanges []*route53store.ResourceRecordSet

	for _, c := range changesList {
		changeMap, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		action, _ := changeMap["Action"].(string)
		if action != "CREATE" && action != "UPSERT" && action != "DELETE" {
			return nil, NewAPIError("InvalidInput", fmt.Sprintf("Invalid action: %s. Must be CREATE, UPSERT, or DELETE", action), 400)
		}
		rrsRaw, _ := changeMap["ResourceRecordSet"].(map[string]interface{})
		if rrsRaw == nil {
			continue
		}

		name := request.GetStringParam(rrsRaw, "Name")
		if name == "" {
			name = request.GetStringParam(rrsRaw, "name")
		}
		rrs := &route53store.ResourceRecordSet{
			Name: normalizeRecordName(name),
			Type: strings.ToUpper(request.GetStringParam(rrsRaw, "Type")),
			TTL:  int64(request.GetIntParam(rrsRaw, "TTL")),
		}
		if rrs.Type == "" {
			rrs.Type = strings.ToUpper(request.GetStringParam(rrsRaw, "type"))
		}
		if rrs.TTL == 0 {
			rrs.TTL = int64(request.GetIntParam(rrsRaw, "ttl"))
		}
		rrs.SetIdentifier = request.GetStringParam(rrsRaw, "SetIdentifier")
		rrs.Weight = int64(request.GetIntParam(rrsRaw, "Weight"))
		rrs.Region = request.GetStringParam(rrsRaw, "Region")
		rrs.Failover = request.GetStringParam(rrsRaw, "Failover")
		rrs.HealthCheckID = request.GetStringParam(rrsRaw, "HealthCheckId")
		rrs.MultiValueAnswer = request.GetBoolParam(rrsRaw, "MultiValueAnswer")
		rrs.TrafficPolicyInstanceID = request.GetStringParam(rrsRaw, "TrafficPolicyInstanceId")

		if geoRaw, ok := rrsRaw["GeoLocation"].(map[string]interface{}); ok {
			rrs.GeoLocation = &route53store.GeoLocation{
				ContinentCode:   request.GetStringParam(geoRaw, "ContinentCode"),
				CountryCode:     request.GetStringParam(geoRaw, "CountryCode"),
				SubdivisionCode: request.GetStringParam(geoRaw, "SubdivisionCode"),
			}
		}

		if recordsRaw, ok := rrsRaw["ResourceRecords"].([]interface{}); ok {
			for _, r := range recordsRaw {
				if rMap, ok := r.(map[string]interface{}); ok {
					rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{
						Value: request.GetStringParam(rMap, "Value"),
					})
				}
			}
		} else if recordsMap, ok := rrsRaw["ResourceRecords"].(map[string]interface{}); ok {
			if rrStr, ok := recordsMap["ResourceRecord"].(string); ok {
				rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{Value: rrStr})
			} else if rrMap, ok := recordsMap["ResourceRecord"].(map[string]interface{}); ok {
				rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{
					Value: request.GetStringParam(rrMap, "Value"),
				})
			} else if rrArr, ok := recordsMap["ResourceRecord"].([]interface{}); ok {
				for _, r := range rrArr {
					if rMap, ok := r.(map[string]interface{}); ok {
						rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{
							Value: request.GetStringParam(rMap, "Value"),
						})
					} else if rStr, ok := r.(string); ok {
						rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{Value: rStr})
					}
				}
			}
		}

		if aliasRaw, ok := rrsRaw["AliasTarget"].(map[string]interface{}); ok {
			rrs.AliasTarget = &route53store.AliasTarget{
				HostedZoneID:         request.GetStringParam(aliasRaw, "HostedZoneId"),
				DNSName:              request.GetStringParam(aliasRaw, "DNSName"),
				EvaluateTargetHealth: request.GetBoolParam(aliasRaw, "EvaluateTargetHealth"),
			}
		}

		switch action {
		case "CREATE":
			if err := st.RecordSets().Upsert(hostedZoneId, rrs); err != nil {
				for _, ac := range appliedChanges {
					if delErr := st.RecordSets().Delete(hostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						logs.Error("Failed to rollback record", logs.String("name", ac.Name), logs.Err(delErr))
					}
				}
				return nil, NewAPIError("InvalidChangeBatch", fmt.Sprintf("Failed to create resource record set %s: %v", rrs.Name, err), 400)
			}
			appliedChanges = append(appliedChanges, rrs)
		case "UPSERT":
			oldRRS, _ := st.RecordSets().Get(hostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier)
			if err := st.RecordSets().Upsert(hostedZoneId, rrs); err != nil {
				for _, ac := range appliedChanges {
					if delErr := st.RecordSets().Delete(hostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						logs.Error("Failed to rollback record", logs.String("name", ac.Name), logs.Err(delErr))
					}
				}
				if oldRRS != nil {
					if restoreErr := st.RecordSets().Upsert(hostedZoneId, oldRRS); restoreErr != nil {
						logs.Error("Failed to restore record", logs.String("name", oldRRS.Name), logs.Err(restoreErr))
					}
				}
				logs.Error("UPSERT record failed", logs.Err(err))
				return nil, NewAPIError("InvalidChangeBatch", fmt.Sprintf("Failed to upsert resource record set %s: %v", rrs.Name, err), 400)
			}
			appliedChanges = append(appliedChanges, rrs)
		case "DELETE":
			deletedRRS, _ := st.RecordSets().Get(hostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier)
			if err := st.RecordSets().Delete(hostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier); err != nil {
				for _, ac := range appliedChanges {
					if delErr := st.RecordSets().Delete(hostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						logs.Error("Failed to rollback record", logs.String("name", ac.Name), logs.Err(delErr))
					}
				}
				if deletedRRS != nil {
					if createErr := st.RecordSets().Create(hostedZoneId, deletedRRS); createErr != nil {
						logs.Error("Failed to restore record", logs.String("name", deletedRRS.Name), logs.Err(createErr))
					}
				}
				logs.Error("DELETE record failed", logs.Err(err))
				return nil, NewAPIError("InvalidChangeBatch", fmt.Sprintf("Failed to delete resource record set %s: %v", rrs.Name, err), 400)
			}
		}
	}

	if err := st.Changes().UpdateStatus(changeId, "INSYNC"); err != nil {
		return nil, NewAPIError("UpdateChange", fmt.Sprintf("Failed to update change status: %v", err), 500)
	}
	change.Status = "INSYNC"

	recordSets, err := st.RecordSets().List(hostedZoneId)
	if err != nil {
		return nil, NewAPIError("ListRecordSets", fmt.Sprintf("Failed to list record sets: %v", err), 500)
	}
	zone.ResourceRecordSetCount = len(recordSets)
	if err := st.HostedZones().Update(zone); err != nil {
		return nil, NewAPIError("UpdateHostedZone", fmt.Sprintf("Failed to update hosted zone: %v", err), 500)
	}

	return map[string]interface{}{
		"ChangeInfo": map[string]interface{}{
			"Id":          "/change/" + changeId,
			"Status":      change.Status,
			"SubmittedAt": change.SubmittedAt.Format(time.RFC3339),
			"Comment":     change.Comment,
		},
	}, nil
}

// ListResourceRecordSets returns resource record sets in a hosted zone with pagination and filtering support.
func (s *Route53Service) ListResourceRecordSets(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	hostedZoneId, err := extractHostedZoneId(req.Parameters, "HostedZoneId")
	if err != nil {
		return nil, err
	}

	_, err = s.getHostedZoneById(reqCtx, hostedZoneId)
	if err != nil {
		return nil, err
	}

	st, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	recordSets, err := st.RecordSets().List(hostedZoneId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	startRecordName := request.GetStringParam(req.Parameters, "StartRecordName")
	if startRecordName == "" {
		startRecordName = request.GetStringParam(req.Parameters, "name")
	}
	startRecordType := request.GetStringParam(req.Parameters, "StartRecordType")
	if startRecordType == "" {
		startRecordType = request.GetStringParam(req.Parameters, "type")
	}
	maxItems := int(request.GetIntParam(req.Parameters, "MaxItems"))
	if maxItems == 0 {
		maxItems = int(request.GetIntParam(req.Parameters, "maxitems"))
	}

	var filtered []*route53store.ResourceRecordSet
	started := startRecordName == "" && startRecordType == ""

	for _, rs := range recordSets {
		if started {
			filtered = append(filtered, rs)
		} else if strings.EqualFold(rs.Name, startRecordName) && strings.EqualFold(rs.Type, startRecordType) {
			started = true
			filtered = append(filtered, rs)
		}
	}

	totalFiltered := len(filtered)
	if maxItems > 0 && totalFiltered > maxItems {
		filtered = filtered[:maxItems]
	}

	records := make([]interface{}, len(filtered))
	for i, rs := range filtered {
		records[i] = s.recordSetToResponse(rs)
	}

	isTruncated := maxItems > 0 && totalFiltered > maxItems

	result := map[string]interface{}{
		"ResourceRecordSets": protocol.XMLElements{ElementName: "ResourceRecordSet", Items: records},
		"IsTruncated":        isTruncated,
		"MaxItems":           maxItems,
	}

	if isTruncated && len(filtered) > 0 {
		lastRecord := filtered[len(filtered)-1]
		result["NextRecordName"] = lastRecord.Name
		result["NextRecordType"] = lastRecord.Type
	}

	return result, nil
}

// GetChange returns the status of a change batch request.
func (s *Route53Service) GetChange(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	id, err := extractChangeId(req.Parameters, "Id")
	if err != nil {
		return nil, err
	}

	change, err := s.getChangeById(reqCtx, id)
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

	if rs.Weight > 0 {
		result["Weight"] = rs.Weight
	}

	if rs.Region != "" {
		result["Region"] = rs.Region
	}

	if rs.Failover != "" {
		result["Failover"] = rs.Failover
	}

	if rs.HealthCheckID != "" {
		result["HealthCheckId"] = rs.HealthCheckID
	}

	if rs.GeoLocation != nil {
		result["GeoLocation"] = map[string]interface{}{
			"ContinentCode":   rs.GeoLocation.ContinentCode,
			"CountryCode":     rs.GeoLocation.CountryCode,
			"SubdivisionCode": rs.GeoLocation.SubdivisionCode,
		}
	}

	return result
}
