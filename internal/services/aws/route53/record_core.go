package route53

import (
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// ChangeResourceRecordSetsInput carries the parsed ChangeBatch for
// ChangeResourceRecordSets. The nested member maps keep their wire shape —
// Core interprets them with the same member names as the wire protocol.
type ChangeResourceRecordSetsInput struct {
	HostedZoneId string
	ChangeBatch  map[string]interface{}
}

// ChangeResourceRecordSetsResult holds the outcome of the change batch.
type ChangeResourceRecordSetsResult struct {
	ChangeId    string
	Status      string
	Comment     string
	SubmittedAt time.Time
}

// ListResourceRecordSetsInput carries the canonical pagination start point
// for ListResourceRecordSets: the start record name is already normalised
// (lowercase, trailing dot) and the type already upper-cased by the caller.
type ListResourceRecordSetsInput struct {
	HostedZoneId          string
	StartRecordName       string
	StartRecordType       string
	StartRecordIdentifier string
	MaxItems              int
}

// ListResourceRecordSetsResult holds one page of record sets together with
// the resume point when the page is truncated.
type ListResourceRecordSetsResult struct {
	RecordSets           []*route53store.ResourceRecordSet
	IsTruncated          bool
	MaxItems             int
	NextRecordName       string
	NextRecordType       string
	NextRecordIdentifier string
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getChangeCore looks up a change by ID, mapping store misses to the
// NoSuchChange AWS error.
func getChangeCore(stores *route53store.Route53Stores, id string) (*route53store.ChangeInfo, error) {
	change, err := stores.Changes().Get(id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchChangeError(id)
		}
		return nil, mapStoreError(err)
	}
	return change, nil
}

// changeResourceRecordSetsCore is the single entry point applying a change
// batch to a hosted zone. It pre-validates the batch shape before creating
// any persistent state, applies each change with rollback of already-applied
// changes on failure, and finally flips the ChangeInfo to INSYNC.
func changeResourceRecordSetsCore(stores *route53store.Route53Stores, input ChangeResourceRecordSetsInput) (*ChangeResourceRecordSetsResult, error) {
	zone, err := getHostedZoneCore(stores, input.HostedZoneId)
	if err != nil {
		return nil, err
	}

	changeBatch := input.ChangeBatch
	if changeBatch == nil {
		return nil, awserrors.NewAWSError("InvalidInput", "ChangeBatch is required", 400)
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
		return nil, awserrors.NewAWSError("InvalidInput", "Changes must be an array", 400)
	}

	if len(changesList) == 0 {
		return nil, awserrors.NewAWSError("InvalidInput", "Changes are required", 400)
	}

	// Pre-validate all changes before creating any persistent state so
	// that an invalid change batch does not leave an orphaned PENDING
	// ChangeInfo record.
	type parsedChange struct {
		action string
		rrsRaw map[string]interface{}
	}
	parsed := make([]parsedChange, 0, len(changesList))
	for _, c := range changesList {
		changeMap, ok := c.(map[string]interface{})
		if !ok {
			return nil, awserrors.NewAWSError("InvalidChangeBatch", "Each element in Changes must be a map", 400)
		}
		action, _ := changeMap["Action"].(string)
		if !validateChangeAction(action) {
			return nil, awserrors.NewAWSError("InvalidInput", "Invalid action: "+action+". Must be CREATE, UPSERT, or DELETE", 400)
		}
		rrsRaw, _ := changeMap["ResourceRecordSet"].(map[string]interface{})
		if rrsRaw == nil {
			return nil, awserrors.NewAWSError("InvalidChangeBatch", "ResourceRecordSet is required for each change", 400)
		}
		parsed = append(parsed, parsedChange{action: action, rrsRaw: rrsRaw})
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

	if err := stores.Changes().Create(change); err != nil {
		return nil, awserrors.NewAWSError("CreateChange", "Failed to create change: "+err.Error(), 500)
	}

	var appliedChanges []*route53store.ResourceRecordSet

	for _, pc := range parsed {
		action := pc.action
		rrsRaw := pc.rrsRaw

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
		if rrs.Type == "" {
			return nil, awserrors.NewAWSError("InvalidInput", "Type is required for resource record set", 400)
		}
		if !validateRecordType(rrs.Type) {
			return nil, awserrors.NewAWSError("InvalidInput", fmt.Sprintf("Invalid record type: %q. Must be a valid RRType enum value", rrs.Type), 400)
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
					value, present := resourceRecordValueFromMap(rMap)
					if err := validateResourceRecordValue(present, value); err != nil {
						return nil, err
					}
					rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{
						Value: value,
					})
				}
			}
		} else if recordsMap, ok := rrsRaw["ResourceRecords"].(map[string]interface{}); ok {
			if rrStr, ok := recordsMap["ResourceRecord"].(string); ok {
				if err := validateResourceRecordValue(true, rrStr); err != nil {
					return nil, err
				}
				rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{Value: rrStr})
			} else if rrMap, ok := recordsMap["ResourceRecord"].(map[string]interface{}); ok {
				value, present := resourceRecordValueFromMap(rrMap)
				if err := validateResourceRecordValue(present, value); err != nil {
					return nil, err
				}
				rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{
					Value: value,
				})
			} else if rrArr, ok := recordsMap["ResourceRecord"].([]interface{}); ok {
				for _, r := range rrArr {
					if rMap, ok := r.(map[string]interface{}); ok {
						value, present := resourceRecordValueFromMap(rMap)
						if err := validateResourceRecordValue(present, value); err != nil {
							return nil, err
						}
						rrs.ResourceRecords = append(rrs.ResourceRecords, &route53store.ResourceRecord{
							Value: value,
						})
					} else if rStr, ok := r.(string); ok {
						if err := validateResourceRecordValue(true, rStr); err != nil {
							return nil, err
						}
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

		if cidrRaw, ok := rrsRaw["CidrRoutingConfig"].(map[string]interface{}); ok {
			rrs.CidrRoutingConfig = &route53store.CidrRoutingConfig{
				CollectionId: request.GetStringParam(cidrRaw, "CollectionId"),
				LocationName: request.GetStringParam(cidrRaw, "LocationName"),
			}
		}

		if gpRaw, ok := rrsRaw["GeoProximityLocation"].(map[string]interface{}); ok {
			gp := &route53store.GeoProximityLocation{
				AWSRegion:      request.GetStringParam(gpRaw, "AWSRegion"),
				LocalZoneGroup: request.GetStringParam(gpRaw, "LocalZoneGroup"),
				Bias:           int64(request.GetIntParam(gpRaw, "Bias")),
			}
			if coordsRaw, ok := gpRaw["Coordinates"].(map[string]interface{}); ok {
				gp.Coordinates = &route53store.Coordinates{
					Latitude:  request.GetFloatParam(coordsRaw, "Latitude"),
					Longitude: request.GetFloatParam(coordsRaw, "Longitude"),
				}
			}
			rrs.GeoProximityLocation = gp
		}

		// Validate TTL and HealthCheckId only for CREATE/UPSERT.
		// DELETE only requires Name + Type (+ SetIdentifier); AWS does
		// not enforce TTL or HealthCheckId existence on deletion.
		if action != "DELETE" {
			// Non-alias records must have TTL > 0.
			if rrs.AliasTarget == nil && rrs.TTL <= 0 {
				return nil, awserrors.NewAWSError("InvalidInput", "TTL is required and must be greater than 0 for non-alias resource record sets", 400)
			}

			// ResourceRecords carries @length(min 1) in the API model: a
			// non-alias record set must hold at least one value, so a
			// present-but-empty list is an invalid change batch.
			if rrs.AliasTarget == nil && len(rrs.ResourceRecords) == 0 {
				return nil, awserrors.NewAWSError("InvalidChangeBatch",
					"Invalid Resource Record: the record set must contain at least one resource record", 400)
			}

			// Verify that the referenced HealthCheckId exists.
			if rrs.HealthCheckID != "" {
				if !stores.HealthChecks().Exists(rrs.HealthCheckID) {
					return nil, awserrors.NewAWSError("InvalidChangeBatch",
						"No health check found with id: "+rrs.HealthCheckID, 400)
				}
			}
		}

		switch action {
		case "CREATE":
			if err := stores.RecordSets().Create(input.HostedZoneId, rrs); err != nil {
				for _, ac := range appliedChanges {
					if delErr := stores.RecordSets().Delete(input.HostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						logs.Error("Failed to rollback record", logs.String("name", ac.Name), logs.Err(delErr))
					}
				}
				return nil, awserrors.NewAWSError("InvalidChangeBatch", "Failed to create resource record set "+rrs.Name+": "+err.Error(), 400)
			}
			appliedChanges = append(appliedChanges, rrs)
		case "UPSERT":
			oldRRS, _ := stores.RecordSets().Get(input.HostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier)
			if err := stores.RecordSets().Upsert(input.HostedZoneId, rrs); err != nil {
				for _, ac := range appliedChanges {
					if delErr := stores.RecordSets().Delete(input.HostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						logs.Error("Failed to rollback record", logs.String("name", ac.Name), logs.Err(delErr))
					}
				}
				if oldRRS != nil {
					if restoreErr := stores.RecordSets().Upsert(input.HostedZoneId, oldRRS); restoreErr != nil {
						logs.Error("Failed to restore record", logs.String("name", oldRRS.Name), logs.Err(restoreErr))
					}
				}
				logs.Error("UPSERT record failed", logs.Err(err))
				return nil, awserrors.NewAWSError("InvalidChangeBatch", "Failed to upsert resource record set "+rrs.Name+": "+err.Error(), 400)
			}
			appliedChanges = append(appliedChanges, rrs)
		case "DELETE":
			deletedRRS, _ := stores.RecordSets().Get(input.HostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier)
			if err := stores.RecordSets().Delete(input.HostedZoneId, rrs.Name, rrs.Type, rrs.SetIdentifier); err != nil {
				for _, ac := range appliedChanges {
					if delErr := stores.RecordSets().Delete(input.HostedZoneId, ac.Name, ac.Type, ac.SetIdentifier); delErr != nil {
						logs.Error("Failed to rollback record", logs.String("name", ac.Name), logs.Err(delErr))
					}
				}
				if deletedRRS != nil {
					if createErr := stores.RecordSets().Create(input.HostedZoneId, deletedRRS); createErr != nil {
						logs.Error("Failed to restore record", logs.String("name", deletedRRS.Name), logs.Err(createErr))
					}
				}
				logs.Error("DELETE record failed", logs.Err(err))
				return nil, awserrors.NewAWSError("InvalidChangeBatch", "Failed to delete resource record set "+rrs.Name+": "+err.Error(), 400)
			}
		}
	}

	if err := stores.Changes().UpdateStatus(changeId, "INSYNC"); err != nil {
		return nil, awserrors.NewAWSError("UpdateChange", "Failed to update change status: "+err.Error(), 500)
	}
	change.Status = "INSYNC"

	recordSets, err := stores.RecordSets().List(input.HostedZoneId)
	if err != nil {
		return nil, awserrors.NewAWSError("ListRecordSets", "Failed to list record sets: "+err.Error(), 500)
	}
	zone.ResourceRecordSetCount = len(recordSets)
	if err := stores.HostedZones().Update(zone); err != nil {
		return nil, awserrors.NewAWSError("UpdateHostedZone", "Failed to update hosted zone: "+err.Error(), 500)
	}

	return &ChangeResourceRecordSetsResult{
		ChangeId:    change.ID,
		Status:      change.Status,
		Comment:     change.Comment,
		SubmittedAt: change.SubmittedAt,
	}, nil
}

// resourceRecordValueFromMap reads the Value member out of a raw
// ResourceRecord map, mirroring the case fallback of GetStringParam. The
// second return reports whether the member is present at all; a missing
// member is distinguishable from an explicitly empty value because the
// wire shape requires the member while the RData length minimum is 0.
func resourceRecordValueFromMap(raw map[string]interface{}) (string, bool) {
	for _, key := range []string{"Value", "value"} {
		if v, ok := raw[key]; ok {
			if s, ok := v.(string); ok {
				return s, true
			}
			return "", true
		}
	}
	return "", false
}

// listResourceRecordSetsCore is the single entry point for listing a hosted
// zone's record sets. It walks the zone's records from the requested start
// point (name, then type, then set identifier) and truncates the page at
// MaxItems, reporting the resume record for truncated pages.
func listResourceRecordSetsCore(stores *route53store.Route53Stores, input ListResourceRecordSetsInput) (*ListResourceRecordSetsResult, error) {
	if _, err := getHostedZoneCore(stores, input.HostedZoneId); err != nil {
		return nil, err
	}

	recordSets, err := stores.RecordSets().List(input.HostedZoneId)
	if err != nil {
		return nil, mapStoreError(err)
	}

	var filtered []*route53store.ResourceRecordSet
	started := input.StartRecordName == ""

	for _, rs := range recordSets {
		if !started {
			// Use lexicographic comparison so that if the start
			// record was deleted, pagination continues from the
			// next record instead of returning nothing.
			rsName := strings.ToLower(rs.Name)
			nameCmp := strings.Compare(rsName, input.StartRecordName)
			if nameCmp > 0 {
				started = true
			} else if nameCmp == 0 {
				typeCmp := strings.Compare(rs.Type, input.StartRecordType)
				if typeCmp > 0 {
					started = true
				} else if typeCmp == 0 {
					if input.StartRecordIdentifier == "" || rs.SetIdentifier >= input.StartRecordIdentifier {
						started = true
					}
				}
			}
		}
		if started {
			filtered = append(filtered, rs)
		}
	}

	allRecords := filtered

	totalFiltered := len(filtered)
	if input.MaxItems > 0 && totalFiltered > input.MaxItems {
		filtered = filtered[:input.MaxItems]
	}

	isTruncated := input.MaxItems > 0 && totalFiltered > input.MaxItems

	result := &ListResourceRecordSetsResult{
		RecordSets:  filtered,
		IsTruncated: isTruncated,
		MaxItems:    input.MaxItems,
	}
	if isTruncated && len(allRecords) > input.MaxItems {
		nextRecord := allRecords[input.MaxItems]
		result.NextRecordName = nextRecord.Name
		result.NextRecordType = nextRecord.Type
		result.NextRecordIdentifier = nextRecord.SetIdentifier
	}
	return result, nil
}
