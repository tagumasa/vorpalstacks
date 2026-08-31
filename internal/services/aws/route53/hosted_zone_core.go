package route53

import (
	"context"
	"fmt"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
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
	DelegationSetID string
	Region          string
	Tags            []types.Tag
}

// CreateHostedZoneResult holds the outcome of createHostedZoneCore.
type CreateHostedZoneResult struct {
	Zone       *route53store.HostedZone
	Idempotent bool
}

// ListHostedZonesInput carries the parameters for ListHostedZones.
type ListHostedZonesInput struct {
	Marker          string
	MaxItems        int
	DelegationSetID string
	HostedZoneType  string
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

// ListHostedZonesByNameInput carries the parameters for ListHostedZonesByName.
type ListHostedZonesByNameInput struct {
	DNSName            string
	HostedZoneIdMarker string
	MaxItems           int
}

// ListHostedZonesByNameResult holds the outcome of listHostedZonesByNameCore.
type ListHostedZonesByNameResult struct {
	HostedZones []*route53store.HostedZone
	IsTruncated bool
}

// UpdateHostedZoneCommentInput carries the parameters for
// UpdateHostedZoneComment.
type UpdateHostedZoneCommentInput struct {
	Id      string
	Comment string
}

// AssociateVPCWithHostedZoneInput carries the parameters for
// AssociateVPCWithHostedZone. A nil VPC means the request carried no VPC
// member at all.
type AssociateVPCWithHostedZoneInput struct {
	HostedZoneId string
	VPC          *route53store.VPC
	Comment      string
	Region       string
}

// DisassociateVPCFromHostedZoneInput carries the parameters for
// DisassociateVPCFromHostedZone. A nil VPC means the request carried no VPC
// member at all.
type DisassociateVPCFromHostedZoneInput struct {
	HostedZoneId string
	VPC          *route53store.VPC
	Comment      string
}

// VPCAssociationResult holds the ChangeInfo outcome of a VPC association or
// disassociation.
type VPCAssociationResult struct {
	ChangeId    string
	Status      string
	Comment     string
	SubmittedAt time.Time
}

// ---------------------------------------------------------------------------
// Core methods — shared by both HTTP API and admin gRPC handlers
// ---------------------------------------------------------------------------

// createHostedZoneCore is the single entry point for creating a hosted zone,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *Route53Service) createHostedZoneCore(stores *route53store.Route53Stores, input CreateHostedZoneInput) (*CreateHostedZoneResult, error) {
	name := strings.ToLower(input.Name)
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "Name is required", 400)
	}
	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	if err := validateDomainName(name); err != nil {
		return nil, err
	}

	// CallerReference is the caller's idempotency token: it must arrive with
	// the request (Nonce, 1 to 128 characters) and is never synthesised
	// server-side, because a fresh token per retry would defeat the
	// execute-once semantics the member provides.
	if err := validateHostedZoneCallerReference(input.CallerReference); err != nil {
		return nil, err
	}
	callerRef := input.CallerReference

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
		DelegationSetID:        input.DelegationSetID,
		NameServers:            nameServers,
		Region:                 input.Region,
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
		return &CreateHostedZoneResult{Zone: existing, Idempotent: true}, nil
	}

	if err := stores.HostedZones().Create(zone); err != nil {
		if route53store.IsAlreadyExists(err) {
			return nil, NewHostedZoneAlreadyExistsError(name)
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

	if len(input.Tags) > 0 {
		resourceKey := "hostedzone/" + zone.ID
		if err := stores.Tags().Tag(resourceKey, input.Tags); err != nil {
			return nil, mapStoreError(err)
		}
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

// listHostedZonesCore is the single entry point for listing hosted zones,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *Route53Service) listHostedZonesCore(stores *route53store.Route53Stores, input ListHostedZonesInput) (*ListHostedZonesResult, error) {
	maxItems := input.MaxItems
	if maxItems <= 0 {
		maxItems = 100
	}

	if input.DelegationSetID != "" || input.HostedZoneType != "" {
		allZones, err := stores.HostedZones().ListByName()
		if err != nil {
			return nil, mapStoreError(err)
		}
		var filtered []*route53store.HostedZone
		for _, z := range allZones {
			if input.DelegationSetID != "" && z.DelegationSetID != input.DelegationSetID {
				continue
			}
			if isPrivateHostedZoneFilter(input.HostedZoneType) && !z.Private {
				continue
			}
			filtered = append(filtered, z)
		}

		isTruncated := len(filtered) > maxItems
		nextMarker := ""
		if isTruncated {
			filtered = filtered[:maxItems]
			if len(filtered) > 0 {
				nextMarker = filtered[maxItems-1].ID
			}
		}
		return &ListHostedZonesResult{
			HostedZones: filtered,
			IsTruncated: isTruncated,
			Marker:      input.Marker,
			NextMarker:  nextMarker,
		}, nil
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

// deleteHostedZoneCore is the single entry point for deleting a hosted zone,
// shared by the HTTP API and the admin gRPC-Web handler.
func (s *Route53Service) deleteHostedZoneCore(stores *route53store.Route53Stores, input DeleteHostedZoneInput) (*DeleteHostedZoneResult, error) {
	if input.Id == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "Id is required", 400)
	}

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

// getHostedZoneCore looks up a hosted zone by ID, mapping store misses to
// the NoSuchHostedZone AWS error.
func getHostedZoneCore(stores *route53store.Route53Stores, id string) (*route53store.HostedZone, error) {
	zone, err := stores.HostedZones().Get(id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchHostedZoneError(id)
		}
		return nil, mapStoreError(err)
	}
	return zone, nil
}

// listHostedZonesByNameCore is the single entry point for the name-ordered
// hosted zone listing. When DNSName is set the result starts at that name
// (or after the HostedZoneIdMarker zone when both are given) and is
// truncated at MaxItems.
func listHostedZonesByNameCore(stores *route53store.Route53Stores, input ListHostedZonesByNameInput) (*ListHostedZonesByNameResult, error) {
	allZones, err := stores.HostedZones().ListByName()
	if err != nil {
		return nil, mapStoreError(err)
	}

	var filtered []*route53store.HostedZone
	if input.DNSName != "" {
		started := false
		if input.HostedZoneIdMarker != "" {
			for _, z := range allZones {
				if strings.EqualFold(z.Name, input.DNSName) && z.ID == input.HostedZoneIdMarker {
					started = true
					continue
				}
				if started {
					filtered = append(filtered, z)
				}
			}
		} else {
			for _, z := range allZones {
				if !started {
					if strings.Compare(z.Name, input.DNSName) >= 0 {
						started = true
					}
				}
				if started {
					filtered = append(filtered, z)
				}
			}
		}
	} else {
		filtered = allZones
	}

	isTruncated := len(filtered) > input.MaxItems
	if isTruncated {
		filtered = filtered[:input.MaxItems]
	}

	return &ListHostedZonesByNameResult{
		HostedZones: filtered,
		IsTruncated: isTruncated,
	}, nil
}

// updateHostedZoneCommentCore is the single entry point updating a hosted
// zone's comment. It validates the comment length and persists the zone.
func updateHostedZoneCommentCore(stores *route53store.Route53Stores, input UpdateHostedZoneCommentInput) (*route53store.HostedZone, error) {
	zone, err := getHostedZoneCore(stores, input.Id)
	if err != nil {
		return nil, err
	}

	if err := validateComment(input.Comment); err != nil {
		return nil, err
	}
	if zone.Config == nil {
		zone.Config = &route53store.HostedZoneConfig{}
	}
	zone.Config.Comment = input.Comment

	if err := stores.HostedZones().Update(zone); err != nil {
		return nil, mapStoreError(err)
	}
	return zone, nil
}

// associateVPCWithHostedZoneCore is the single entry point associating a
// VPC with a private hosted zone. It verifies the VPC exists in EC2 (best
// effort via the event bus), rejects duplicate associations, records the
// ChangeInfo, and persists the zone.
func (s *Route53Service) associateVPCWithHostedZoneCore(ctx context.Context, stores *route53store.Route53Stores, input AssociateVPCWithHostedZoneInput) (*VPCAssociationResult, error) {
	if input.VPC == nil {
		return nil, awserrors.NewAWSError("InvalidInput", "VPC is required", 400)
	}

	zone, err := getHostedZoneCore(stores, input.HostedZoneId)
	if err != nil {
		return nil, err
	}

	if !zone.Private {
		return nil, awserrors.NewAWSError("InvalidInput", "Cannot associate VPC with a public hosted zone", 400)
	}

	if err := s.validateVPC(ctx, input.Region, input.VPC.VPCID, input.VPC.VPCRegion); err != nil {
		return nil, awserrors.NewAWSError("InvalidVPCId",
			fmt.Sprintf("The VPC %s in region %s does not exist", input.VPC.VPCID, input.VPC.VPCRegion), 400)
	}
	for _, existing := range zone.VPCs {
		if existing.VPCID == input.VPC.VPCID && existing.VPCRegion == input.VPC.VPCRegion {
			return nil, awserrors.NewAWSError("VPCAlreadyAssociated", "VPC is already associated with the hosted zone", 400)
		}
	}
	zone.VPCs = append(zone.VPCs, input.VPC)

	if err := stores.HostedZones().Update(zone); err != nil {
		return nil, mapStoreError(err)
	}

	// Validate the ChangeInfo comment after persisting the association,
	// matching the original operation order.
	if err := validateComment(input.Comment); err != nil {
		return nil, err
	}
	changeId := generateChangeId()
	now := time.Now()
	if err := stores.Changes().Create(&route53store.ChangeInfo{
		ID:          changeId,
		Status:      "INSYNC",
		SubmittedAt: now,
		Comment:     input.Comment,
	}); err != nil {
		return nil, mapStoreError(err)
	}

	return &VPCAssociationResult{
		ChangeId:    changeId,
		Status:      "INSYNC",
		Comment:     input.Comment,
		SubmittedAt: now,
	}, nil
}

// disassociateVPCFromHostedZoneCore is the single entry point removing a VPC
// association from a private hosted zone. It rejects the removal when the
// VPC is not associated or when it is the last association, records the
// ChangeInfo, and persists the zone.
func disassociateVPCFromHostedZoneCore(stores *route53store.Route53Stores, input DisassociateVPCFromHostedZoneInput) (*VPCAssociationResult, error) {
	if input.VPC == nil {
		return nil, awserrors.NewAWSError("InvalidInput", "VPC is required", 400)
	}

	zone, err := getHostedZoneCore(stores, input.HostedZoneId)
	if err != nil {
		return nil, err
	}

	var newVPCs []*route53store.VPC
	for _, v := range zone.VPCs {
		if v.VPCID != input.VPC.VPCID || v.VPCRegion != input.VPC.VPCRegion {
			newVPCs = append(newVPCs, v)
		}
	}

	if len(newVPCs) == len(zone.VPCs) {
		// Use VPCAssociationNotFound instead of InvalidInput.
		return nil, NewVPCAssociationNotFoundError(input.VPC.VPCID)
	}

	// Prevent removing the last VPC association from a private zone.
	if len(newVPCs) == 0 {
		return nil, NewLastVPCAssociationError()
	}

	zone.VPCs = newVPCs
	if err := stores.HostedZones().Update(zone); err != nil {
		return nil, mapStoreError(err)
	}

	// Validate the ChangeInfo comment after persisting the
	// disassociation, matching the original operation order.
	if err := validateComment(input.Comment); err != nil {
		return nil, err
	}
	changeId := generateChangeId()
	now := time.Now()
	if err := stores.Changes().Create(&route53store.ChangeInfo{
		ID:          changeId,
		Status:      "INSYNC",
		SubmittedAt: now,
		Comment:     input.Comment,
	}); err != nil {
		return nil, mapStoreError(err)
	}

	return &VPCAssociationResult{
		ChangeId:    changeId,
		Status:      "INSYNC",
		Comment:     input.Comment,
		SubmittedAt: now,
	}, nil
}
