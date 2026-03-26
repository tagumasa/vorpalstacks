package route53

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/route53"
	route53connect "vorpalstacks/internal/pb/aws/route53/route53connect"
)

// AdminHandler provides Route 53 service administration functionality.
// It implements the Route53ServiceHandler interface for gRPC-Web communication.
type AdminHandler struct{}

var _ route53connect.Route53ServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Route 53 AdminHandler.
func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// ListHostedZones lists hosted zones in Route 53.
// Returns an empty list as this operation is not fully implemented.
func (h *AdminHandler) ListHostedZones(ctx context.Context, req *connect.Request[pb.ListHostedZonesRequest]) (*connect.Response[pb.ListHostedZonesResponse], error) {
	return connect.NewResponse(&pb.ListHostedZonesResponse{
		Hostedzones: []*pb.HostedZone{},
	}), nil
}

// ActivateKeySigningKey activates a key signing key for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ActivateKeySigningKey(ctx context.Context, req *connect.Request[pb.ActivateKeySigningKeyRequest]) (*connect.Response[pb.ActivateKeySigningKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ActivateKeySigningKey is not implemented"))
}

// AssociateVPCWithHostedZone associates an Amazon VPC with a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) AssociateVPCWithHostedZone(ctx context.Context, req *connect.Request[pb.AssociateVPCWithHostedZoneRequest]) (*connect.Response[pb.AssociateVPCWithHostedZoneResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.AssociateVPCWithHostedZone is not implemented"))
}

// ChangeCidrCollection changes the CIDR collection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ChangeCidrCollection(ctx context.Context, req *connect.Request[pb.ChangeCidrCollectionRequest]) (*connect.Response[pb.ChangeCidrCollectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ChangeCidrCollection is not implemented"))
}

// ChangeResourceRecordSets creates, changes, or deletes DNS records.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ChangeResourceRecordSets(ctx context.Context, req *connect.Request[pb.ChangeResourceRecordSetsRequest]) (*connect.Response[pb.ChangeResourceRecordSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ChangeResourceRecordSets is not implemented"))
}

// ChangeTagsForResource adds, edits, or deletes tags for a Route 53 resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ChangeTagsForResource(ctx context.Context, req *connect.Request[pb.ChangeTagsForResourceRequest]) (*connect.Response[pb.ChangeTagsForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ChangeTagsForResource is not implemented"))
}

// CreateCidrCollection creates a CIDR collection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateCidrCollection(ctx context.Context, req *connect.Request[pb.CreateCidrCollectionRequest]) (*connect.Response[pb.CreateCidrCollectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateCidrCollection is not implemented"))
}

// CreateHealthCheck creates a health check.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateHealthCheck(ctx context.Context, req *connect.Request[pb.CreateHealthCheckRequest]) (*connect.Response[pb.CreateHealthCheckResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateHealthCheck is not implemented"))
}

// CreateHostedZone creates a new hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateHostedZone(ctx context.Context, req *connect.Request[pb.CreateHostedZoneRequest]) (*connect.Response[pb.CreateHostedZoneResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateHostedZone is not implemented"))
}

// CreateKeySigningKey creates a key signing key for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateKeySigningKey(ctx context.Context, req *connect.Request[pb.CreateKeySigningKeyRequest]) (*connect.Response[pb.CreateKeySigningKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateKeySigningKey is not implemented"))
}

// CreateQueryLoggingConfig creates a configuration for DNS query logging.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateQueryLoggingConfig(ctx context.Context, req *connect.Request[pb.CreateQueryLoggingConfigRequest]) (*connect.Response[pb.CreateQueryLoggingConfigResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateQueryLoggingConfig is not implemented"))
}

// CreateReusableDelegationSet creates a reusable delegation set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateReusableDelegationSet(ctx context.Context, req *connect.Request[pb.CreateReusableDelegationSetRequest]) (*connect.Response[pb.CreateReusableDelegationSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateReusableDelegationSet is not implemented"))
}

// CreateTrafficPolicy creates a traffic policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTrafficPolicy(ctx context.Context, req *connect.Request[pb.CreateTrafficPolicyRequest]) (*connect.Response[pb.CreateTrafficPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateTrafficPolicy is not implemented"))
}

// CreateTrafficPolicyInstance creates a traffic policy instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTrafficPolicyInstance(ctx context.Context, req *connect.Request[pb.CreateTrafficPolicyInstanceRequest]) (*connect.Response[pb.CreateTrafficPolicyInstanceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateTrafficPolicyInstance is not implemented"))
}

// CreateTrafficPolicyVersion creates a new version of an existing traffic policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateTrafficPolicyVersion(ctx context.Context, req *connect.Request[pb.CreateTrafficPolicyVersionRequest]) (*connect.Response[pb.CreateTrafficPolicyVersionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateTrafficPolicyVersion is not implemented"))
}

// CreateVPCAssociationAuthorization authorises associating a VPC with a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) CreateVPCAssociationAuthorization(ctx context.Context, req *connect.Request[pb.CreateVPCAssociationAuthorizationRequest]) (*connect.Response[pb.CreateVPCAssociationAuthorizationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.CreateVPCAssociationAuthorization is not implemented"))
}

// DeactivateKeySigningKey deactivates a key signing key for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeactivateKeySigningKey(ctx context.Context, req *connect.Request[pb.DeactivateKeySigningKeyRequest]) (*connect.Response[pb.DeactivateKeySigningKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeactivateKeySigningKey is not implemented"))
}

// DeleteCidrCollection deletes a CIDR collection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteCidrCollection(ctx context.Context, req *connect.Request[pb.DeleteCidrCollectionRequest]) (*connect.Response[pb.DeleteCidrCollectionResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteCidrCollection is not implemented"))
}

// DeleteHealthCheck deletes a health check.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteHealthCheck(ctx context.Context, req *connect.Request[pb.DeleteHealthCheckRequest]) (*connect.Response[pb.DeleteHealthCheckResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteHealthCheck is not implemented"))
}

// DeleteHostedZone deletes a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteHostedZone(ctx context.Context, req *connect.Request[pb.DeleteHostedZoneRequest]) (*connect.Response[pb.DeleteHostedZoneResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteHostedZone is not implemented"))
}

// DeleteKeySigningKey deletes a key signing key.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteKeySigningKey(ctx context.Context, req *connect.Request[pb.DeleteKeySigningKeyRequest]) (*connect.Response[pb.DeleteKeySigningKeyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteKeySigningKey is not implemented"))
}

// DeleteQueryLoggingConfig deletes a DNS query logging configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteQueryLoggingConfig(ctx context.Context, req *connect.Request[pb.DeleteQueryLoggingConfigRequest]) (*connect.Response[pb.DeleteQueryLoggingConfigResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteQueryLoggingConfig is not implemented"))
}

// DeleteReusableDelegationSet deletes a reusable delegation set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteReusableDelegationSet(ctx context.Context, req *connect.Request[pb.DeleteReusableDelegationSetRequest]) (*connect.Response[pb.DeleteReusableDelegationSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteReusableDelegationSet is not implemented"))
}

// DeleteTrafficPolicy deletes a traffic policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteTrafficPolicy(ctx context.Context, req *connect.Request[pb.DeleteTrafficPolicyRequest]) (*connect.Response[pb.DeleteTrafficPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteTrafficPolicy is not implemented"))
}

// DeleteTrafficPolicyInstance deletes a traffic policy instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteTrafficPolicyInstance(ctx context.Context, req *connect.Request[pb.DeleteTrafficPolicyInstanceRequest]) (*connect.Response[pb.DeleteTrafficPolicyInstanceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteTrafficPolicyInstance is not implemented"))
}

// DeleteVPCAssociationAuthorization removes authorisation for associating a VPC with a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DeleteVPCAssociationAuthorization(ctx context.Context, req *connect.Request[pb.DeleteVPCAssociationAuthorizationRequest]) (*connect.Response[pb.DeleteVPCAssociationAuthorizationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DeleteVPCAssociationAuthorization is not implemented"))
}

// DisableHostedZoneDNSSEC disables DNSSEC for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisableHostedZoneDNSSEC(ctx context.Context, req *connect.Request[pb.DisableHostedZoneDNSSECRequest]) (*connect.Response[pb.DisableHostedZoneDNSSECResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DisableHostedZoneDNSSEC is not implemented"))
}

// DisassociateVPCFromHostedZone disassociates an Amazon VPC from a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) DisassociateVPCFromHostedZone(ctx context.Context, req *connect.Request[pb.DisassociateVPCFromHostedZoneRequest]) (*connect.Response[pb.DisassociateVPCFromHostedZoneResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.DisassociateVPCFromHostedZone is not implemented"))
}

// EnableHostedZoneDNSSEC enables DNSSEC for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) EnableHostedZoneDNSSEC(ctx context.Context, req *connect.Request[pb.EnableHostedZoneDNSSECRequest]) (*connect.Response[pb.EnableHostedZoneDNSSECResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.EnableHostedZoneDNSSEC is not implemented"))
}

// GetAccountLimit returns the current limit for a Route 53 resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetAccountLimit(ctx context.Context, req *connect.Request[pb.GetAccountLimitRequest]) (*connect.Response[pb.GetAccountLimitResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetAccountLimit is not implemented"))
}

// GetChange returns the status of a change batch request.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetChange(ctx context.Context, req *connect.Request[pb.GetChangeRequest]) (*connect.Response[pb.GetChangeResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetChange is not implemented"))
}

// GetCheckerIpRanges returns the IP ranges used by Route 53 health checkers.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetCheckerIpRanges(ctx context.Context, req *connect.Request[pb.GetCheckerIpRangesRequest]) (*connect.Response[pb.GetCheckerIpRangesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetCheckerIpRanges is not implemented"))
}

// GetDNSSEC returns DNSSEC configuration information for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetDNSSEC(ctx context.Context, req *connect.Request[pb.GetDNSSECRequest]) (*connect.Response[pb.GetDNSSECResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetDNSSEC is not implemented"))
}

// GetGeoLocation returns information about a geographic location.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetGeoLocation(ctx context.Context, req *connect.Request[pb.GetGeoLocationRequest]) (*connect.Response[pb.GetGeoLocationResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetGeoLocation is not implemented"))
}

// GetHealthCheck returns information about a health check.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHealthCheck(ctx context.Context, req *connect.Request[pb.GetHealthCheckRequest]) (*connect.Response[pb.GetHealthCheckResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetHealthCheck is not implemented"))
}

// GetHealthCheckCount returns the number of health checks.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHealthCheckCount(ctx context.Context, req *connect.Request[pb.GetHealthCheckCountRequest]) (*connect.Response[pb.GetHealthCheckCountResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetHealthCheckCount is not implemented"))
}

// GetHealthCheckLastFailureReason returns the reason for the last failure of a health check.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHealthCheckLastFailureReason(ctx context.Context, req *connect.Request[pb.GetHealthCheckLastFailureReasonRequest]) (*connect.Response[pb.GetHealthCheckLastFailureReasonResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetHealthCheckLastFailureReason is not implemented"))
}

// GetHealthCheckStatus returns the status of a health check.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHealthCheckStatus(ctx context.Context, req *connect.Request[pb.GetHealthCheckStatusRequest]) (*connect.Response[pb.GetHealthCheckStatusResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetHealthCheckStatus is not implemented"))
}

// GetHostedZone returns information about a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHostedZone(ctx context.Context, req *connect.Request[pb.GetHostedZoneRequest]) (*connect.Response[pb.GetHostedZoneResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetHostedZone is not implemented"))
}

// GetHostedZoneCount returns the number of hosted zones.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHostedZoneCount(ctx context.Context, req *connect.Request[pb.GetHostedZoneCountRequest]) (*connect.Response[pb.GetHostedZoneCountResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetHostedZoneCount is not implemented"))
}

// GetHostedZoneLimit returns the limit for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetHostedZoneLimit(ctx context.Context, req *connect.Request[pb.GetHostedZoneLimitRequest]) (*connect.Response[pb.GetHostedZoneLimitResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetHostedZoneLimit is not implemented"))
}

// GetQueryLoggingConfig returns information about a DNS query logging configuration.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetQueryLoggingConfig(ctx context.Context, req *connect.Request[pb.GetQueryLoggingConfigRequest]) (*connect.Response[pb.GetQueryLoggingConfigResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetQueryLoggingConfig is not implemented"))
}

// GetReusableDelegationSet returns information about a reusable delegation set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetReusableDelegationSet(ctx context.Context, req *connect.Request[pb.GetReusableDelegationSetRequest]) (*connect.Response[pb.GetReusableDelegationSetResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetReusableDelegationSet is not implemented"))
}

// GetReusableDelegationSetLimit returns the limit for a reusable delegation set.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetReusableDelegationSetLimit(ctx context.Context, req *connect.Request[pb.GetReusableDelegationSetLimitRequest]) (*connect.Response[pb.GetReusableDelegationSetLimitResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetReusableDelegationSetLimit is not implemented"))
}

// GetTrafficPolicy returns information about a traffic policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetTrafficPolicy(ctx context.Context, req *connect.Request[pb.GetTrafficPolicyRequest]) (*connect.Response[pb.GetTrafficPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetTrafficPolicy is not implemented"))
}

// GetTrafficPolicyInstance returns information about a traffic policy instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetTrafficPolicyInstance(ctx context.Context, req *connect.Request[pb.GetTrafficPolicyInstanceRequest]) (*connect.Response[pb.GetTrafficPolicyInstanceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetTrafficPolicyInstance is not implemented"))
}

// GetTrafficPolicyInstanceCount returns the number of traffic policy instances.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) GetTrafficPolicyInstanceCount(ctx context.Context, req *connect.Request[pb.GetTrafficPolicyInstanceCountRequest]) (*connect.Response[pb.GetTrafficPolicyInstanceCountResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.GetTrafficPolicyInstanceCount is not implemented"))
}

// ListCidrBlocks lists the CIDR blocks in a CIDR collection.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListCidrBlocks(ctx context.Context, req *connect.Request[pb.ListCidrBlocksRequest]) (*connect.Response[pb.ListCidrBlocksResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListCidrBlocks is not implemented"))
}

// ListCidrCollections lists CIDR collections.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListCidrCollections(ctx context.Context, req *connect.Request[pb.ListCidrCollectionsRequest]) (*connect.Response[pb.ListCidrCollectionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListCidrCollections is not implemented"))
}

// ListCidrLocations lists locations for CIDR collections.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListCidrLocations(ctx context.Context, req *connect.Request[pb.ListCidrLocationsRequest]) (*connect.Response[pb.ListCidrLocationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListCidrLocations is not implemented"))
}

// ListGeoLocations lists geographic locations available for use with geolocation routing.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListGeoLocations(ctx context.Context, req *connect.Request[pb.ListGeoLocationsRequest]) (*connect.Response[pb.ListGeoLocationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListGeoLocations is not implemented"))
}

// ListHealthChecks lists health checks.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListHealthChecks(ctx context.Context, req *connect.Request[pb.ListHealthChecksRequest]) (*connect.Response[pb.ListHealthChecksResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListHealthChecks is not implemented"))
}

// ListHostedZonesByName lists hosted zones by name.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListHostedZonesByName(ctx context.Context, req *connect.Request[pb.ListHostedZonesByNameRequest]) (*connect.Response[pb.ListHostedZonesByNameResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListHostedZonesByName is not implemented"))
}

// ListHostedZonesByVPC lists hosted zones associated with a VPC.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListHostedZonesByVPC(ctx context.Context, req *connect.Request[pb.ListHostedZonesByVPCRequest]) (*connect.Response[pb.ListHostedZonesByVPCResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListHostedZonesByVPC is not implemented"))
}

// ListQueryLoggingConfigs lists DNS query logging configurations.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListQueryLoggingConfigs(ctx context.Context, req *connect.Request[pb.ListQueryLoggingConfigsRequest]) (*connect.Response[pb.ListQueryLoggingConfigsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListQueryLoggingConfigs is not implemented"))
}

// ListResourceRecordSets lists DNS records in a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListResourceRecordSets(ctx context.Context, req *connect.Request[pb.ListResourceRecordSetsRequest]) (*connect.Response[pb.ListResourceRecordSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListResourceRecordSets is not implemented"))
}

// ListReusableDelegationSets lists reusable delegation sets.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListReusableDelegationSets(ctx context.Context, req *connect.Request[pb.ListReusableDelegationSetsRequest]) (*connect.Response[pb.ListReusableDelegationSetsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListReusableDelegationSets is not implemented"))
}

// ListTagsForResource lists tags for a Route 53 resource.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResource(ctx context.Context, req *connect.Request[pb.ListTagsForResourceRequest]) (*connect.Response[pb.ListTagsForResourceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListTagsForResource is not implemented"))
}

// ListTagsForResources lists tags for multiple Route 53 resources.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTagsForResources(ctx context.Context, req *connect.Request[pb.ListTagsForResourcesRequest]) (*connect.Response[pb.ListTagsForResourcesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListTagsForResources is not implemented"))
}

// ListTrafficPolicies lists traffic policies.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTrafficPolicies(ctx context.Context, req *connect.Request[pb.ListTrafficPoliciesRequest]) (*connect.Response[pb.ListTrafficPoliciesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListTrafficPolicies is not implemented"))
}

// ListTrafficPolicyInstances lists traffic policy instances.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTrafficPolicyInstances(ctx context.Context, req *connect.Request[pb.ListTrafficPolicyInstancesRequest]) (*connect.Response[pb.ListTrafficPolicyInstancesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListTrafficPolicyInstances is not implemented"))
}

// ListTrafficPolicyInstancesByHostedZone lists traffic policy instances in a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTrafficPolicyInstancesByHostedZone(ctx context.Context, req *connect.Request[pb.ListTrafficPolicyInstancesByHostedZoneRequest]) (*connect.Response[pb.ListTrafficPolicyInstancesByHostedZoneResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListTrafficPolicyInstancesByHostedZone is not implemented"))
}

// ListTrafficPolicyInstancesByPolicy lists traffic policy instances for a traffic policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTrafficPolicyInstancesByPolicy(ctx context.Context, req *connect.Request[pb.ListTrafficPolicyInstancesByPolicyRequest]) (*connect.Response[pb.ListTrafficPolicyInstancesByPolicyResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListTrafficPolicyInstancesByPolicy is not implemented"))
}

// ListTrafficPolicyVersions lists versions of a traffic policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListTrafficPolicyVersions(ctx context.Context, req *connect.Request[pb.ListTrafficPolicyVersionsRequest]) (*connect.Response[pb.ListTrafficPolicyVersionsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListTrafficPolicyVersions is not implemented"))
}

// ListVPCAssociationAuthorizations lists authorisations for VPC association.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) ListVPCAssociationAuthorizations(ctx context.Context, req *connect.Request[pb.ListVPCAssociationAuthorizationsRequest]) (*connect.Response[pb.ListVPCAssociationAuthorizationsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.ListVPCAssociationAuthorizations is not implemented"))
}

// TestDNSAnswer tests DNS resolution for a domain name.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) TestDNSAnswer(ctx context.Context, req *connect.Request[pb.TestDNSAnswerRequest]) (*connect.Response[pb.TestDNSAnswerResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.TestDNSAnswer is not implemented"))
}

// UpdateHealthCheck updates a health check.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateHealthCheck(ctx context.Context, req *connect.Request[pb.UpdateHealthCheckRequest]) (*connect.Response[pb.UpdateHealthCheckResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.UpdateHealthCheck is not implemented"))
}

// UpdateHostedZoneComment updates the comment for a hosted zone.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateHostedZoneComment(ctx context.Context, req *connect.Request[pb.UpdateHostedZoneCommentRequest]) (*connect.Response[pb.UpdateHostedZoneCommentResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.UpdateHostedZoneComment is not implemented"))
}

// UpdateHostedZoneFeatures enables or disables hosted zone features.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateHostedZoneFeatures(ctx context.Context, req *connect.Request[pb.UpdateHostedZoneFeaturesRequest]) (*connect.Response[pb.UpdateHostedZoneFeaturesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.UpdateHostedZoneFeatures is not implemented"))
}

// UpdateTrafficPolicyComment updates the comment for a traffic policy.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateTrafficPolicyComment(ctx context.Context, req *connect.Request[pb.UpdateTrafficPolicyCommentRequest]) (*connect.Response[pb.UpdateTrafficPolicyCommentResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.UpdateTrafficPolicyComment is not implemented"))
}

// UpdateTrafficPolicyInstance updates a traffic policy instance.
// Use the AWS REST API for this operation as gRPC-Web does not support it.
func (h *AdminHandler) UpdateTrafficPolicyInstance(ctx context.Context, req *connect.Request[pb.UpdateTrafficPolicyInstanceRequest]) (*connect.Response[pb.UpdateTrafficPolicyInstanceResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("route53.Route53Service.UpdateTrafficPolicyInstance is not implemented"))
}
