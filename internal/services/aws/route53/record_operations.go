package route53

import (
	"context"
	"log"
	"strings"
	"time"

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
		log.Printf("CreateChange failed: %v", err)
		return nil, NewAPIError("CreateChange", "Failed to create change. See server logs for details.", 500)
	}

	var appliedChanges []*route53store.ResourceRecordSet

	for _, c := range changesList {
		changeMap, ok := c.(map[string]interface{})
		if !ok {
			continue
		}

		action, _ := changeMap["Action"].(string)
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
			if err := st.RecordSets().Create(hostedZoneId, rrs); err != nil {
				for _, ac := range appliedChanges {
					if delErr := st.RecordSets().Delete(hostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						log.Printf("Failed to rollback record %s: %v", ac.Name, delErr)
					}
				}
				log.Printf("CREATE record failed: %v", err)
				return nil, NewAPIError("InvalidChangeBatch", "Failed to create resource record set. See server logs for details.", 400)
			}
			appliedChanges = append(appliedChanges, rrs)
		case "UPSERT":
			oldRRS, _ := st.RecordSets().Get(hostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier)
			if err := st.RecordSets().Upsert(hostedZoneId, rrs); err != nil {
				for _, ac := range appliedChanges {
					if delErr := st.RecordSets().Delete(hostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						log.Printf("Failed to rollback record %s: %v", ac.Name, delErr)
					}
				}
				if oldRRS != nil {
					if createErr := st.RecordSets().Create(hostedZoneId, oldRRS); createErr != nil {
						log.Printf("Failed to restore record %s: %v", oldRRS.Name, createErr)
					}
				}
				log.Printf("UPSERT record failed: %v", err)
				return nil, NewAPIError("InvalidChangeBatch", "Failed to upsert resource record set. See server logs for details.", 400)
			}
			appliedChanges = append(appliedChanges, rrs)
		case "DELETE":
			deletedRRS, _ := st.RecordSets().Get(hostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier)
			if err := st.RecordSets().Delete(hostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier); err != nil {
				for _, ac := range appliedChanges {
					if delErr := st.RecordSets().Delete(hostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						log.Printf("Failed to rollback record %s: %v", ac.Name, delErr)
					}
				}
				if deletedRRS != nil {
					if createErr := st.RecordSets().Create(hostedZoneId, deletedRRS); createErr != nil {
						log.Printf("Failed to restore record %s: %v", deletedRRS.Name, createErr)
					}
				}
				log.Printf("DELETE record failed: %v", err)
				return nil, NewAPIError("InvalidChangeBatch", "Failed to delete resource record set. See server logs for details.", 400)
			}
		}
	}

	if err := st.Changes().UpdateStatus(changeId, "INSYNC"); err != nil {
		log.Printf("UpdateChange status failed: %v", err)
		return nil, NewAPIError("UpdateChange", "Failed to update change status. See server logs for details.", 500)
	}
	change.Status = "INSYNC"

	recordSets, err := st.RecordSets().List(hostedZoneId)
	if err != nil {
		log.Printf("ListRecordSets failed: %v", err)
		return nil, NewAPIError("ListRecordSets", "Failed to list record sets. See server logs for details.", 500)
	}
	zone.ResourceRecordSetCount = len(recordSets)
	if err := st.HostedZones().Update(zone); err != nil {
		log.Printf("UpdateHostedZone failed: %v", err)
		return nil, NewAPIError("UpdateHostedZone", "Failed to update hosted zone. See server logs for details.", 500)
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
	startRecordType := request.GetStringParam(req.Parameters, "StartRecordType")
	maxItems := int(request.GetIntParam(req.Parameters, "MaxItems"))

	var filtered []*route53store.ResourceRecordSet
	started := startRecordName == "" && startRecordType == ""

	for _, rs := range recordSets {
		if started {
			filtered = append(filtered, rs)
		} else if strings.EqualFold(rs.Name, startRecordName) && strings.EqualFold(rs.Type, startRecordType) {
			started = true
		}
	}

	if maxItems > 0 && len(filtered) > maxItems {
		filtered = filtered[:maxItems]
	}

	records := make([]interface{}, len(filtered))
	for i, rs := range filtered {
		records[i] = s.recordSetToResponse(rs)
	}

	isTruncated := maxItems > 0 && len(filtered) == maxItems && (startRecordName != "" || startRecordType != "" || len(records) < len(recordSets))

	return map[string]interface{}{
		"ResourceRecordSets": protocol.XMLElements{ElementName: "ResourceRecordSet", Items: records},
		"IsTruncated":        isTruncated,
	}, nil
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
		"TTL":  rs.TTL,
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
