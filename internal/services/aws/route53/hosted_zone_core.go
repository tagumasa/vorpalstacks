package route53

import (
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// CreateHostedZoneInput carries every field that CreateHostedZone needs,
// independent of the wire protocol. Both the HTTP API handler and the
// admin gRPC handler build this struct and delegate to createHostedZoneCore.
type CreateHostedZoneInput struct {
	Name            string
	CallerReference string
	Comment         string
	PrivateZone     bool
	VPCID           string
	VPCRegion       string
}

// CreateHostedZoneResult holds the outcome of createHostedZoneCore.
type CreateHostedZoneResult struct {
	Zone *route53store.HostedZone
}

// ListHostedZonesInput carries the parameters for ListHostedZones.
type ListHostedZonesInput struct {
	Marker   string
	MaxItems int
}

// ListHostedZonesResult holds the outcome of listHostedZonesCore.
type ListHostedZonesResult struct {
	HostedZones []*route53store.HostedZone
	IsTruncated bool
	Marker      string
	NextMarker  string
}

// DeleteHostedZoneInput carries the parameters for DeleteHostedZone.
type DeleteHostedZoneInput struct {
	Id string
}

// DeleteHostedZoneResult holds the outcome of deleteHostedZoneCore.
type DeleteHostedZoneResult struct {
	ChangeId    string
	Status      string
	SubmittedAt time.Time
}

// ---------------------------------------------------------------------------
// Core methods — shared by both HTTP API and admin gRPC handlers
// ---------------------------------------------------------------------------

// createHostedZoneCore contains the shared creation logic for a hosted zone.
// It validates input, generates default NS/SOA records, and persists the zone.
func (s *Route53Service) createHostedZoneCore(stores *route53store.Route53Stores, input CreateHostedZoneInput) (*CreateHostedZoneResult, error) {
	name := strings.ToLower(input.Name)
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	if err := validateDomainName(name); err != nil {
		return nil, err
	}

	callerRef := input.CallerReference
	if callerRef == "" {
		callerRef = fmt.Sprintf("%x", fmt.Sprintf("%s-%d", name, time.Now().UnixNano()))
	}

	privateZone := input.PrivateZone
	if input.VPCID != "" {
		privateZone = true
	}

	if privateZone && input.VPCID == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "A VPC is required when creating a private hosted zone", 400)
	}

	nameServers := generateNameServers(4)

	zone := &route53store.HostedZone{
		ID:                     generateHostedZoneId(),
		Name:                   route53store.NormalizeZoneName(name),
		CallerReference:        callerRef,
		Config:                 &route53store.HostedZoneConfig{Comment: input.Comment, PrivateZone: privateZone},
		ResourceRecordSetCount: 0,
		Private:                privateZone,
		DelegationSetID:        "",
		NameServers:            nameServers,
		AccountID:              s.accountID,
		CreatedAt:              time.Now(),
	}

	if input.VPCID != "" {
		zone.VPCs = []*route53store.VPC{{
			VPCID:     input.VPCID,
			VPCRegion: input.VPCRegion,
		}}
	}

	if existing, err := stores.HostedZones().GetByCallerReference(callerRef); err == nil && existing != nil {
		return &CreateHostedZoneResult{Zone: existing}, nil
	}

	if err := stores.HostedZones().Create(zone); err != nil {
		if route53store.IsAlreadyExists(err) {
			return nil, awserrors.NewAWSError("HostedZoneAlreadyExists", fmt.Sprintf("Hosted zone already exists: %s", name), 400)
		}
		return nil, mapStoreError(err)
	}

	nsRecords := make([]*route53store.ResourceRecord, len(nameServers))
	for i, ns := range nameServers {
		nsRecords[i] = &route53store.ResourceRecord{Value: ns}
	}
	if err := stores.RecordSets().Create(zone.ID, &route53store.ResourceRecordSet{
		Name:            zone.Name,
		Type:            "NS",
		TTL:             172800,
		ResourceRecords: nsRecords,
	}); err != nil {
		return nil, mapStoreError(err)
	}

	if err := stores.RecordSets().Create(zone.ID, &route53store.ResourceRecordSet{
		Name: zone.Name,
		Type: "SOA",
		TTL:  900,
		ResourceRecords: []*route53store.ResourceRecord{
			{Value: fmt.Sprintf("%s %s 1 7200 900 1209600 86400", zone.Name, nameServers[0])},
		},
	}); err != nil {
		return nil, mapStoreError(err)
	}

	zone.ResourceRecordSetCount = 2
	zone.ChangeID = generateChangeId()
	if err := stores.HostedZones().Update(zone); err != nil {
		return nil, mapStoreError(err)
	}

	if err := stores.Changes().Create(&route53store.ChangeInfo{
		ID:          zone.ChangeID,
		Status:      "INSYNC",
		SubmittedAt: zone.CreatedAt,
		Comment:     input.Comment,
	}); err != nil {
		return nil, mapStoreError(err)
	}

	return &CreateHostedZoneResult{Zone: zone}, nil
}

// listHostedZonesCore contains the shared listing logic for hosted zones.
func (s *Route53Service) listHostedZonesCore(stores *route53store.Route53Stores, input ListHostedZonesInput) (*ListHostedZonesResult, error) {
	maxItems := input.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}

	result, err := stores.HostedZones().List(input.Marker, maxItems)
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &ListHostedZonesResult{
		HostedZones: result.HostedZones,
		IsTruncated: result.IsTruncated,
		Marker:      input.Marker,
		NextMarker:  result.Marker,
	}, nil
}

// deleteHostedZoneCore contains the shared deletion logic for a hosted zone.
// It verifies the zone is empty, cleans up NS/SOA records and tags, deletes
// the zone, and creates a ChangeInfo record.
func (s *Route53Service) deleteHostedZoneCore(stores *route53store.Route53Stores, input DeleteHostedZoneInput) (*DeleteHostedZoneResult, error) {
	id := strings.TrimPrefix(input.Id, "/hostedzone/")

	recordSets, err := stores.RecordSets().List(id)
	if err != nil {
		return nil, mapStoreError(err)
	}

	if len(recordSets) > 0 {
		userRecords := 0
		for _, rs := range recordSets {
			if rs.Type != "NS" && rs.Type != "SOA" {
				userRecords++
			}
		}
		if userRecords > 0 {
			return nil, awserrors.NewAWSError("HostedZoneNotEmpty", "The hosted zone must be empty before it can be deleted.", 400)
		}
	}

	// Clean up default NS and SOA records.
	for _, rs := range recordSets {
		if rs.Type == "NS" || rs.Type == "SOA" {
			_ = stores.RecordSets().Delete(id, rs.Name, rs.Type, rs.SetIdentifier)
		}
	}

	if err := stores.Tags().Raw().Delete("hostedzone/" + id); err != nil {
		return nil, awserrors.NewAWSError("InvalidInput", fmt.Sprintf("Failed to delete tags: %v", err), 500)
	}

	if err := stores.HostedZones().Delete(id); err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchHostedZoneError(id)
		}
		return nil, mapStoreError(err)
	}

	changeId := generateChangeId()
	now := time.Now()
	if err := stores.Changes().Create(&route53store.ChangeInfo{
		ID:          changeId,
		Status:      "INSYNC",
		SubmittedAt: now,
	}); err != nil {
		return nil, mapStoreError(err)
	}

	return &DeleteHostedZoneResult{
		ChangeId:    changeId,
		Status:      "INSYNC",
		SubmittedAt: now,
	}, nil
}
