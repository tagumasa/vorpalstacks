// PutAuditEvents — the cloudtrail-data service's ingestion endpoint. The
// operation carries its scalar members as URI query parameters (httpQuery
// traits in the vendored model) and the auditEvents list as the REST-JSON
// body; events flow through a channel into the channel's destination event
// data stores, bypassing the stores' ingestion selectors — the channel's
// destinations are the routing.
package cloudtrail

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
)

// externalIDPattern compiles the model's ExternalId pattern.
var externalIDPattern = regexp.MustCompile(cloudtrailstore.ExternalIDPattern)

// PutAuditEventsInput carries the cloudtrail-data PutAuditEvents members.
type PutAuditEventsInput struct {
	ChannelARN    string
	ExternalID    string
	ExternalIDSet bool
	AuditEvents   []AuditEventInput
}

// AuditEventInput is one payload entry: the source event's ID, the JSON
// record, and the optional base64-SHA256 checksum over that record.
type AuditEventInput struct {
	ID                string
	EventData         string
	EventDataChecksum string
	ChecksumSet       bool
}

// PutAuditEvents ingests external application events into CloudTrail Lake
// through a channel (the cloudtrail-data service).
func (s *CloudTrailService) PutAuditEvents(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	in := PutAuditEventsInput{
		ChannelARN: request.GetStringParam(req.Parameters, "channelArn"),
		ExternalID: request.GetStringParam(req.Parameters, "externalId"),
	}
	// Presence, not nullness: an explicit empty externalId is a present
	// member that fails the model's length bound, while a JSON null is the
	// absent member — the same reading the auditEvents member takes one
	// block below.
	if v, ok := req.Parameters["externalId"]; ok && v != nil {
		in.ExternalIDSet = true
	}
	raw, ok := req.Parameters["auditEvents"]
	if !ok || raw == nil {
		return nil, newChannelUnsupportedSchema("auditEvents is required")
	}
	list, ok := raw.([]interface{})
	if !ok {
		return nil, newChannelUnsupportedSchema("auditEvents must be a list of event entries")
	}
	for _, item := range list {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, newChannelUnsupportedSchema("auditEvents must be a list of event entries")
		}
		ev := AuditEventInput{
			ID:        wireString(m, "id"),
			EventData: wireString(m, "eventData"),
		}
		ev.EventDataChecksum = wireString(m, "eventDataChecksum")
		ev.ChecksumSet = ev.EventDataChecksum != ""
		in.AuditEvents = append(in.AuditEvents, ev)
	}

	return s.putAuditEventsCore(store, in)
}

// putAuditEventsCore is the single validation and ingestion path for
// PutAuditEvents. Entries are validated individually: an id or eventData
// the model marks required (FieldNotFound), a checksum that does not match
// the record (InvalidChecksum — the base64-SHA256 algorithm the reference
// documents), or a record that does not parse (InvalidData) fails that
// entry alone; the response carries the survivors and the failures
// separately, both members always present. Duplicate non-empty ids across
// the request are refused wholesale (DuplicatedAuditEventId is a
// request-level error); an empty id is not a duplicate identity — it fails
// its own entry.
func (s *CloudTrailService) putAuditEventsCore(store cloudtrailstore.CloudTrailStoreInterface, in PutAuditEventsInput) (map[string]interface{}, error) {
	if in.ChannelARN == "" {
		return nil, newInvalidChannelARN("channelArn is required")
	}
	ch, err := resolveDataChannel(store, in.ChannelARN)
	if err != nil {
		return nil, err
	}

	// externalId is "conditionally required when the channel's resource
	// policy includes an external ID"; platform channels carry no
	// external-ID policy conditions, so the member is accepted with
	// nothing to match against.
	destEDSs, err := resolveChannelDestinationEDSs(store, ch.Destinations)
	if err != nil {
		// The management-plane destination errors have no counterpart in
		// this service's declared vocabulary; a destination that no longer
		// resolves surfaces as the channel being unusable.
		var wireErr interface{ GetCode() string }
		if errors.As(err, &wireErr) && wireErr.GetCode() == "EventDataStoreNotFoundException" {
			return nil, newChannelNotFound("The channel's destination event data store does not exist")
		}
		return nil, err
	}

	// One request carries at least one and at most 100 events: "You can add
	// up to 100 of these events (or up to 1 MB) per PutAuditEvents request"
	// (Quotas in AWS CloudTrail); the cloudtrail-data model marks the
	// AuditEvents list @length 1-100 (the typed SDK rejects an out-of-range
	// list client-side, so both guards answer the raw-HTTP path). The
	// service declares no generic validation shape — ChannelUnsupportedSchema
	// is the auditEvents member's validation shape.
	if len(in.AuditEvents) == 0 {
		return nil, newChannelUnsupportedSchema("auditEvents must contain at least one entry")
	}
	if len(in.AuditEvents) > cloudtrailstore.MaxAuditEventsPerRequest {
		return nil, newChannelUnsupportedSchema(fmt.Sprintf(
			"auditEvents must contain at most %d entries", cloudtrailstore.MaxAuditEventsPerRequest))
	}

	// externalId is optional, but a present value must satisfy the model's
	// ExternalId shape — length 2-1224 and the printable pattern. The
	// service declares no dedicated shape for the member, so
	// ChannelUnsupportedSchema serves as it does for the list bound above.
	if in.ExternalIDSet {
		if len(in.ExternalID) < cloudtrailstore.MinExternalIDLength ||
			len(in.ExternalID) > cloudtrailstore.MaxExternalIDLength ||
			!externalIDPattern.MatchString(in.ExternalID) {
			return nil, newChannelUnsupportedSchema(
				"externalId must be 2-1224 characters drawn from [A-Za-z0-9_+=,.@:/-]")
		}
	}

	// Duplicate non-empty ids are refused wholesale (DuplicatedAuditEventId
	// is a request-level error). An empty id is not a duplicate identity:
	// it fails its own entry below, per the model's required id member.
	seen := make(map[string]bool, len(in.AuditEvents))
	for _, ev := range in.AuditEvents {
		if ev.ID == "" {
			continue
		}
		if seen[ev.ID] {
			return nil, newDuplicatedAuditEventId(
				fmt.Sprintf("Two or more entries in the request have the same event ID: %s", ev.ID))
		}
		seen[ev.ID] = true
	}

	successful := make([]map[string]interface{}, 0, len(in.AuditEvents))
	failed := make([]map[string]interface{}, 0)
	for _, ev := range in.AuditEvents {
		eventID, failCode, failMessage := ingestAuditEvent(store, destEDSs, ev)
		if failCode == "" {
			successful = append(successful, map[string]interface{}{
				"id":      ev.ID,
				"eventID": eventID,
			})
			continue
		}
		failed = append(failed, map[string]interface{}{
			"id":           ev.ID,
			"errorCode":    failCode,
			"errorMessage": failMessage,
		})
	}

	return map[string]interface{}{
		"successful": successful,
		"failed":     failed,
	}, nil
}

// ingestAuditEvent validates and ingests one audit event entry. It returns
// the CloudTrail-assigned event ID on success, or the model's errorCode
// vocabulary entry and message on failure.
func ingestAuditEvent(store cloudtrailstore.CloudTrailStoreInterface, destEDSs []*cloudtrailstore.EventDataStore, ev AuditEventInput) (string, string, string) {
	if ev.ID == "" {
		return "", "FieldNotFound", "The event entry is missing its id"
	}
	if ev.EventData == "" {
		return "", "FieldNotFound", "The event entry is missing its eventData"
	}
	if ev.ChecksumSet {
		digest := sha256.Sum256([]byte(ev.EventData))
		if base64.StdEncoding.EncodeToString(digest[:]) != ev.EventDataChecksum {
			return "", "InvalidChecksum", "The checksum does not match the event data"
		}
	}

	event, err := cloudtrailstore.RecordEventFromPayload(ev.EventData, uuid.NewString(), "ActivityAuditLog")
	if err != nil {
		return "", "InvalidData", "The event data is not a valid JSON event record"
	}

	// Every destination is verified accepting before any write happens,
	// and the copies join one store transaction: a refusal or failure
	// discovered at a later destination must not leave an earlier one
	// holding the event, whose retry (a fresh event record) would
	// duplicate the surviving copies.
	for _, eds := range destEDSs {
		// A destination with ingestion stopped accepts nothing; the
		// mapping of this condition to InvalidRecipient is inferred — the
		// errorCode vocabulary is documented, this trigger is not.
		if !eds.IngestionEnabled || eds.Status == "PENDING_DELETION" {
			return "", "InvalidRecipient", fmt.Sprintf(
				"The destination event data store %s is not accepting events", eds.EventDataStoreARN)
		}
	}
	edsIDs := make([]string, len(destEDSs))
	for i, eds := range destEDSs {
		edsIDs[i] = eds.EventDataStoreID
	}
	if err := store.PutEventIntoEDSs(edsIDs, event); err != nil {
		if errors.Is(err, cloudtrailstore.ErrEventDataStoreNotFound) {
			// A destination record vanished between the request's
			// resolution and the write (the restore window's hard
			// delete): the same not-accepting condition the ingestion
			// guard reports, naming the destination that no longer
			// resolves.
			for _, eds := range destEDSs {
				if _, err := store.GetEventDataStore(eds.EventDataStoreID); errors.Is(err, cloudtrailstore.ErrEventDataStoreNotFound) {
					return "", "InvalidRecipient", fmt.Sprintf(
						"The destination event data store %s is not accepting events", eds.EventDataStoreARN)
				}
			}
			return "", "InvalidRecipient", "The destination event data store is not accepting events"
		}
		return "", "InternalFailure", err.Error()
	}
	return event.EventID, "", ""
}

// resolveDataChannel resolves the channelArn query member — "The ARN or ID
// (the ARN suffix) of a channel" — to the stored channel over the shared
// ARN-or-UUID resolution. A value in ARN form that is not a CloudTrail
// channel ARN answers the data model's InvalidChannelARN; an ARN or suffix
// that names no stored channel answers ChannelNotFound.
func resolveDataChannel(store cloudtrailstore.CloudTrailStoreInterface, channelArn string) (*cloudtrailstore.Channel, error) {
	ch, err := resolveChannelByARNOrUUID(store, channelArn, func(v string) error {
		return newInvalidChannelARN(
			fmt.Sprintf("The specified channel ARN is not a valid channel ARN: %s", v))
	})
	if err != nil {
		if errors.Is(err, cloudtrailstore.ErrChannelNotFound) {
			return nil, newChannelNotFound("The channel could not be found")
		}
		return nil, err
	}
	return ch, nil
}

// wireString reads a string member from a decoded JSON object; a member
// that is absent or not a string reads as empty.
func wireString(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
