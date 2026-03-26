// This is a generated file - do not edit.
//
// Generated from route53.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:async' as $async;
import 'dart:core' as $core;

import 'package:grpc/service_api.dart' as $grpc;
import 'package:protobuf/protobuf.dart' as $pb;

import 'route53.pb.dart' as $0;

export 'route53.pb.dart';

/// Route53Service provides route53 API operations.
@$pb.GrpcServiceName('route53.Route53Service')
class Route53ServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  Route53ServiceClient(super.channel, {super.options, super.interceptors});

  /// Activates a key-signing key (KSK) so that it can be used for signing by DNSSEC. This operation changes the KSK status to ACTIVE.
  /// HTTP: POST /2013-04-01/keysigningkey/{HostedZoneId}/{Name}/activate
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ActivateKeySigningKeyResponse> activateKeySigningKey(
    $0.ActivateKeySigningKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$activateKeySigningKey, request, options: options);
  }

  /// Associates an Amazon VPC with a private hosted zone. To perform the association, the VPC and the private hosted zone must already exist. You can't convert a public hosted zone into a private hosted...
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/associatevpc
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.AssociateVPCWithHostedZoneResponse>
      associateVPCWithHostedZone(
    $0.AssociateVPCWithHostedZoneRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateVPCWithHostedZone, request,
        options: options);
  }

  /// Creates, changes, or deletes CIDR blocks within a collection. Contains authoritative IP information mapping blocks to one or multiple locations. A change request can update multiple locations in a ...
  /// HTTP: POST /2013-04-01/cidrcollection/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ChangeCidrCollectionResponse> changeCidrCollection(
    $0.ChangeCidrCollectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changeCidrCollection, request, options: options);
  }

  /// Creates, changes, or deletes a resource record set, which contains authoritative DNS information for a specified domain name or subdomain name. For example, you can use ChangeResourceRecordSets to ...
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/rrset
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ChangeResourceRecordSetsResponse>
      changeResourceRecordSets(
    $0.ChangeResourceRecordSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changeResourceRecordSets, request,
        options: options);
  }

  /// Adds, edits, or deletes tags for a health check or a hosted zone. For information about using tags for cost allocation, see Using Cost Allocation Tags in the Billing and Cost Management User Guide.
  /// HTTP: POST /2013-04-01/tags/{ResourceType}/{ResourceId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ChangeTagsForResourceResponse> changeTagsForResource(
    $0.ChangeTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changeTagsForResource, request, options: options);
  }

  /// Creates a CIDR collection in the current Amazon Web Services account.
  /// HTTP: POST /2013-04-01/cidrcollection
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateCidrCollectionResponse> createCidrCollection(
    $0.CreateCidrCollectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCidrCollection, request, options: options);
  }

  /// Creates a new health check. For information about adding health checks to resource record sets, see HealthCheckId in ChangeResourceRecordSets. ELB Load Balancers If you're registering EC2 instances...
  /// HTTP: POST /2013-04-01/healthcheck
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateHealthCheckResponse> createHealthCheck(
    $0.CreateHealthCheckRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createHealthCheck, request, options: options);
  }

  /// Creates a new public or private hosted zone. You create records in a public hosted zone to define how you want to route traffic on the internet for a domain, such as example.com, and its subdomains...
  /// HTTP: POST /2013-04-01/hostedzone
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateHostedZoneResponse> createHostedZone(
    $0.CreateHostedZoneRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createHostedZone, request, options: options);
  }

  /// Creates a new key-signing key (KSK) associated with a hosted zone. You can only have two KSKs per hosted zone.
  /// HTTP: POST /2013-04-01/keysigningkey
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateKeySigningKeyResponse> createKeySigningKey(
    $0.CreateKeySigningKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createKeySigningKey, request, options: options);
  }

  /// Creates a configuration for DNS query logging. After you create a query logging configuration, Amazon Route 53 begins to publish log data to an Amazon CloudWatch Logs log group. DNS query logs cont...
  /// HTTP: POST /2013-04-01/queryloggingconfig
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateQueryLoggingConfigResponse>
      createQueryLoggingConfig(
    $0.CreateQueryLoggingConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createQueryLoggingConfig, request,
        options: options);
  }

  /// Creates a delegation set (a group of four name servers) that can be reused by multiple hosted zones that were created by the same Amazon Web Services account. You can also create a reusable delegat...
  /// HTTP: POST /2013-04-01/delegationset
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateReusableDelegationSetResponse>
      createReusableDelegationSet(
    $0.CreateReusableDelegationSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createReusableDelegationSet, request,
        options: options);
  }

  /// Creates a traffic policy, which you use to create multiple DNS resource record sets for one domain name (such as example.com) or one subdomain name (such as www.example.com).
  /// HTTP: POST /2013-04-01/trafficpolicy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateTrafficPolicyResponse> createTrafficPolicy(
    $0.CreateTrafficPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTrafficPolicy, request, options: options);
  }

  /// Creates resource record sets in a specified hosted zone based on the settings in a specified traffic policy version. In addition, CreateTrafficPolicyInstance associates the resource record sets wit...
  /// HTTP: POST /2013-04-01/trafficpolicyinstance
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateTrafficPolicyInstanceResponse>
      createTrafficPolicyInstance(
    $0.CreateTrafficPolicyInstanceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTrafficPolicyInstance, request,
        options: options);
  }

  /// Creates a new version of an existing traffic policy. When you create a new version of a traffic policy, you specify the ID of the traffic policy that you want to update and a JSON-formatted documen...
  /// HTTP: POST /2013-04-01/trafficpolicy/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateTrafficPolicyVersionResponse>
      createTrafficPolicyVersion(
    $0.CreateTrafficPolicyVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTrafficPolicyVersion, request,
        options: options);
  }

  /// Authorizes the Amazon Web Services account that created a specified VPC to submit an AssociateVPCWithHostedZone request to associate the VPC with a specified hosted zone that was created by a diffe...
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/authorizevpcassociation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.CreateVPCAssociationAuthorizationResponse>
      createVPCAssociationAuthorization(
    $0.CreateVPCAssociationAuthorizationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createVPCAssociationAuthorization, request,
        options: options);
  }

  /// Deactivates a key-signing key (KSK) so that it will not be used for signing by DNSSEC. This operation changes the KSK status to INACTIVE.
  /// HTTP: POST /2013-04-01/keysigningkey/{HostedZoneId}/{Name}/deactivate
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeactivateKeySigningKeyResponse>
      deactivateKeySigningKey(
    $0.DeactivateKeySigningKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deactivateKeySigningKey, request,
        options: options);
  }

  /// Deletes a CIDR collection in the current Amazon Web Services account. The collection must be empty before it can be deleted.
  /// HTTP: DELETE /2013-04-01/cidrcollection/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteCidrCollectionResponse> deleteCidrCollection(
    $0.DeleteCidrCollectionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCidrCollection, request, options: options);
  }

  /// Deletes a health check. Amazon Route 53 does not prevent you from deleting a health check even if the health check is associated with one or more resource record sets. If you delete a health check ...
  /// HTTP: DELETE /2013-04-01/healthcheck/{HealthCheckId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteHealthCheckResponse> deleteHealthCheck(
    $0.DeleteHealthCheckRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteHealthCheck, request, options: options);
  }

  /// Deletes a hosted zone. If the hosted zone was created by another service, such as Cloud Map, see Deleting Public Hosted Zones That Were Created by Another Service in the Amazon Route 53 Developer ...
  /// HTTP: DELETE /2013-04-01/hostedzone/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteHostedZoneResponse> deleteHostedZone(
    $0.DeleteHostedZoneRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteHostedZone, request, options: options);
  }

  /// Deletes a key-signing key (KSK). Before you can delete a KSK, you must deactivate it. The KSK must be deactivated before you can delete it regardless of whether the hosted zone is enabled for DNSSE...
  /// HTTP: DELETE /2013-04-01/keysigningkey/{HostedZoneId}/{Name}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteKeySigningKeyResponse> deleteKeySigningKey(
    $0.DeleteKeySigningKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteKeySigningKey, request, options: options);
  }

  /// Deletes a configuration for DNS query logging. If you delete a configuration, Amazon Route 53 stops sending query logs to CloudWatch Logs. Route 53 doesn't delete any logs that are already in Cloud...
  /// HTTP: DELETE /2013-04-01/queryloggingconfig/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteQueryLoggingConfigResponse>
      deleteQueryLoggingConfig(
    $0.DeleteQueryLoggingConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteQueryLoggingConfig, request,
        options: options);
  }

  /// Deletes a reusable delegation set. You can delete a reusable delegation set only if it isn't associated with any hosted zones. To verify that the reusable delegation set is not associated with any ...
  /// HTTP: DELETE /2013-04-01/delegationset/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteReusableDelegationSetResponse>
      deleteReusableDelegationSet(
    $0.DeleteReusableDelegationSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteReusableDelegationSet, request,
        options: options);
  }

  /// Deletes a traffic policy. When you delete a traffic policy, Route 53 sets a flag on the policy to indicate that it has been deleted. However, Route 53 never fully deletes the traffic policy. Note t...
  /// HTTP: DELETE /2013-04-01/trafficpolicy/{Id}/{Version}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteTrafficPolicyResponse> deleteTrafficPolicy(
    $0.DeleteTrafficPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTrafficPolicy, request, options: options);
  }

  /// Deletes a traffic policy instance and all of the resource record sets that Amazon Route 53 created when you created the instance. In the Route 53 console, traffic policy instances are known as poli...
  /// HTTP: DELETE /2013-04-01/trafficpolicyinstance/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteTrafficPolicyInstanceResponse>
      deleteTrafficPolicyInstance(
    $0.DeleteTrafficPolicyInstanceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTrafficPolicyInstance, request,
        options: options);
  }

  /// Removes authorization to submit an AssociateVPCWithHostedZone request to associate a specified VPC with a hosted zone that was created by a different account. You must use the account that created ...
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/deauthorizevpcassociation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DeleteVPCAssociationAuthorizationResponse>
      deleteVPCAssociationAuthorization(
    $0.DeleteVPCAssociationAuthorizationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteVPCAssociationAuthorization, request,
        options: options);
  }

  /// Disables DNSSEC signing in a specific hosted zone. This action does not deactivate any key-signing keys (KSKs) that are active in the hosted zone.
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/disable-dnssec
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DisableHostedZoneDNSSECResponse>
      disableHostedZoneDNSSEC(
    $0.DisableHostedZoneDNSSECRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableHostedZoneDNSSEC, request,
        options: options);
  }

  /// Disassociates an Amazon Virtual Private Cloud (Amazon VPC) from an Amazon Route 53 private hosted zone. Note the following: You can't disassociate the last Amazon VPC from a private hosted zone. Yo...
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/disassociatevpc
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.DisassociateVPCFromHostedZoneResponse>
      disassociateVPCFromHostedZone(
    $0.DisassociateVPCFromHostedZoneRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disassociateVPCFromHostedZone, request,
        options: options);
  }

  /// Enables DNSSEC signing in a specific hosted zone.
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/enable-dnssec
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.EnableHostedZoneDNSSECResponse>
      enableHostedZoneDNSSEC(
    $0.EnableHostedZoneDNSSECRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableHostedZoneDNSSEC, request,
        options: options);
  }

  /// Gets the specified limit for the current account, for example, the maximum number of health checks that you can create using the account. For the default limit, see Limits in the Amazon Route 53 De...
  /// HTTP: GET /2013-04-01/accountlimit/{Type}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetAccountLimitResponse> getAccountLimit(
    $0.GetAccountLimitRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccountLimit, request, options: options);
  }

  /// Returns the current status of a change batch request. The status is one of the following values: PENDING indicates that the changes in this request have not propagated to all Amazon Route 53 DNS se...
  /// HTTP: GET /2013-04-01/change/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetChangeResponse> getChange(
    $0.GetChangeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getChange, request, options: options);
  }

  /// Route 53 does not perform authorization for this API because it retrieves information that is already available to the public. GetCheckerIpRanges still works, but we recommend that you download ip-...
  /// HTTP: GET /2013-04-01/checkeripranges
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetCheckerIpRangesResponse> getCheckerIpRanges(
    $0.GetCheckerIpRangesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCheckerIpRanges, request, options: options);
  }

  /// Returns information about DNSSEC for a specific hosted zone, including the key-signing keys (KSKs) in the hosted zone.
  /// HTTP: GET /2013-04-01/hostedzone/{HostedZoneId}/dnssec
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetDNSSECResponse> getDNSSEC(
    $0.GetDNSSECRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDNSSEC, request, options: options);
  }

  /// Gets information about whether a specified geographic location is supported for Amazon Route 53 geolocation resource record sets. Route 53 does not perform authorization for this API because it ret...
  /// HTTP: GET /2013-04-01/geolocation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetGeoLocationResponse> getGeoLocation(
    $0.GetGeoLocationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGeoLocation, request, options: options);
  }

  /// Gets information about a specified health check.
  /// HTTP: GET /2013-04-01/healthcheck/{HealthCheckId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetHealthCheckResponse> getHealthCheck(
    $0.GetHealthCheckRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHealthCheck, request, options: options);
  }

  /// Retrieves the number of health checks that are associated with the current Amazon Web Services account.
  /// HTTP: GET /2013-04-01/healthcheckcount
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetHealthCheckCountResponse> getHealthCheckCount(
    $0.GetHealthCheckCountRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHealthCheckCount, request, options: options);
  }

  /// Gets the reason that a specified health check failed most recently.
  /// HTTP: GET /2013-04-01/healthcheck/{HealthCheckId}/lastfailurereason
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetHealthCheckLastFailureReasonResponse>
      getHealthCheckLastFailureReason(
    $0.GetHealthCheckLastFailureReasonRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHealthCheckLastFailureReason, request,
        options: options);
  }

  /// Gets status of a specified health check. This API is intended for use during development to diagnose behavior. It doesn’t support production use-cases with high query rates that require immediate...
  /// HTTP: GET /2013-04-01/healthcheck/{HealthCheckId}/status
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetHealthCheckStatusResponse> getHealthCheckStatus(
    $0.GetHealthCheckStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHealthCheckStatus, request, options: options);
  }

  /// Gets information about a specified hosted zone including the four name servers assigned to the hosted zone. returns the VPCs associated with the specified hosted zone and does not reflect the VPC a...
  /// HTTP: GET /2013-04-01/hostedzone/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetHostedZoneResponse> getHostedZone(
    $0.GetHostedZoneRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHostedZone, request, options: options);
  }

  /// Retrieves the number of hosted zones that are associated with the current Amazon Web Services account.
  /// HTTP: GET /2013-04-01/hostedzonecount
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetHostedZoneCountResponse> getHostedZoneCount(
    $0.GetHostedZoneCountRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHostedZoneCount, request, options: options);
  }

  /// Gets the specified limit for a specified hosted zone, for example, the maximum number of records that you can create in the hosted zone. For the default limit, see Limits in the Amazon Route 53 Dev...
  /// HTTP: GET /2013-04-01/hostedzonelimit/{HostedZoneId}/{Type}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetHostedZoneLimitResponse> getHostedZoneLimit(
    $0.GetHostedZoneLimitRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHostedZoneLimit, request, options: options);
  }

  /// Gets information about a specified configuration for DNS query logging. For more information about DNS query logs, see CreateQueryLoggingConfig and Logging DNS Queries.
  /// HTTP: GET /2013-04-01/queryloggingconfig/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetQueryLoggingConfigResponse> getQueryLoggingConfig(
    $0.GetQueryLoggingConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueryLoggingConfig, request, options: options);
  }

  /// Retrieves information about a specified reusable delegation set, including the four name servers that are assigned to the delegation set.
  /// HTTP: GET /2013-04-01/delegationset/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetReusableDelegationSetResponse>
      getReusableDelegationSet(
    $0.GetReusableDelegationSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getReusableDelegationSet, request,
        options: options);
  }

  /// Gets the maximum number of hosted zones that you can associate with the specified reusable delegation set. For the default limit, see Limits in the Amazon Route 53 Developer Guide. To request a hig...
  /// HTTP: GET /2013-04-01/reusabledelegationsetlimit/{DelegationSetId}/{Type}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetReusableDelegationSetLimitResponse>
      getReusableDelegationSetLimit(
    $0.GetReusableDelegationSetLimitRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getReusableDelegationSetLimit, request,
        options: options);
  }

  /// Gets information about a specific traffic policy version. For information about how of deleting a traffic policy affects the response from GetTrafficPolicy, see DeleteTrafficPolicy.
  /// HTTP: GET /2013-04-01/trafficpolicy/{Id}/{Version}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetTrafficPolicyResponse> getTrafficPolicy(
    $0.GetTrafficPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTrafficPolicy, request, options: options);
  }

  /// Gets information about a specified traffic policy instance. Use GetTrafficPolicyInstance with the id of new traffic policy instance to confirm that the CreateTrafficPolicyInstance or an UpdateTraff...
  /// HTTP: GET /2013-04-01/trafficpolicyinstance/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetTrafficPolicyInstanceResponse>
      getTrafficPolicyInstance(
    $0.GetTrafficPolicyInstanceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTrafficPolicyInstance, request,
        options: options);
  }

  /// Gets the number of traffic policy instances that are associated with the current Amazon Web Services account.
  /// HTTP: GET /2013-04-01/trafficpolicyinstancecount
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.GetTrafficPolicyInstanceCountResponse>
      getTrafficPolicyInstanceCount(
    $0.GetTrafficPolicyInstanceCountRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTrafficPolicyInstanceCount, request,
        options: options);
  }

  /// Returns a paginated list of location objects and their CIDR blocks.
  /// HTTP: GET /2013-04-01/cidrcollection/{CollectionId}/cidrblocks
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListCidrBlocksResponse> listCidrBlocks(
    $0.ListCidrBlocksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCidrBlocks, request, options: options);
  }

  /// Returns a paginated list of CIDR collections in the Amazon Web Services account (metadata only).
  /// HTTP: GET /2013-04-01/cidrcollection
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListCidrCollectionsResponse> listCidrCollections(
    $0.ListCidrCollectionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCidrCollections, request, options: options);
  }

  /// Returns a paginated list of CIDR locations for the given collection (metadata only, does not include CIDR blocks).
  /// HTTP: GET /2013-04-01/cidrcollection/{CollectionId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListCidrLocationsResponse> listCidrLocations(
    $0.ListCidrLocationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCidrLocations, request, options: options);
  }

  /// Retrieves a list of supported geographic locations. Countries are listed first, and continents are listed last. If Amazon Route 53 supports subdivisions for a country (for example, states or provin...
  /// HTTP: GET /2013-04-01/geolocations
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListGeoLocationsResponse> listGeoLocations(
    $0.ListGeoLocationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGeoLocations, request, options: options);
  }

  /// Retrieve a list of the health checks that are associated with the current Amazon Web Services account.
  /// HTTP: GET /2013-04-01/healthcheck
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListHealthChecksResponse> listHealthChecks(
    $0.ListHealthChecksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listHealthChecks, request, options: options);
  }

  /// Retrieves a list of the public and private hosted zones that are associated with the current Amazon Web Services account. The response includes a HostedZones child element for each hosted zone. Ama...
  /// HTTP: GET /2013-04-01/hostedzone
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListHostedZonesResponse> listHostedZones(
    $0.ListHostedZonesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listHostedZones, request, options: options);
  }

  /// Retrieves a list of your hosted zones in lexicographic order. The response includes a HostedZones child element for each hosted zone created by the current Amazon Web Services account. ListHostedZo...
  /// HTTP: GET /2013-04-01/hostedzonesbyname
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListHostedZonesByNameResponse> listHostedZonesByName(
    $0.ListHostedZonesByNameRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listHostedZonesByName, request, options: options);
  }

  /// Lists all the private hosted zones that a specified VPC is associated with, regardless of which Amazon Web Services account or Amazon Web Services service owns the hosted zones. The HostedZoneOwner...
  /// HTTP: GET /2013-04-01/hostedzonesbyvpc
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListHostedZonesByVPCResponse> listHostedZonesByVPC(
    $0.ListHostedZonesByVPCRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listHostedZonesByVPC, request, options: options);
  }

  /// Lists the configurations for DNS query logging that are associated with the current Amazon Web Services account or the configuration that is associated with a specified hosted zone. For more inform...
  /// HTTP: GET /2013-04-01/queryloggingconfig
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListQueryLoggingConfigsResponse>
      listQueryLoggingConfigs(
    $0.ListQueryLoggingConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listQueryLoggingConfigs, request,
        options: options);
  }

  /// Lists the resource record sets in a specified hosted zone. ListResourceRecordSets returns up to 300 resource record sets at a time in ASCII order, beginning at a position specified by the name and ...
  /// HTTP: GET /2013-04-01/hostedzone/{HostedZoneId}/rrset
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListResourceRecordSetsResponse>
      listResourceRecordSets(
    $0.ListResourceRecordSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResourceRecordSets, request,
        options: options);
  }

  /// Retrieves a list of the reusable delegation sets that are associated with the current Amazon Web Services account.
  /// HTTP: GET /2013-04-01/delegationset
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListReusableDelegationSetsResponse>
      listReusableDelegationSets(
    $0.ListReusableDelegationSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listReusableDelegationSets, request,
        options: options);
  }

  /// Lists tags for one health check or hosted zone. For information about using tags for cost allocation, see Using Cost Allocation Tags in the Billing and Cost Management User Guide.
  /// HTTP: GET /2013-04-01/tags/{ResourceType}/{ResourceId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Lists tags for up to 10 health checks or hosted zones. For information about using tags for cost allocation, see Using Cost Allocation Tags in the Billing and Cost Management User Guide.
  /// HTTP: POST /2013-04-01/tags/{ResourceType}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTagsForResourcesResponse> listTagsForResources(
    $0.ListTagsForResourcesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResources, request, options: options);
  }

  /// Gets information about the latest version for every traffic policy that is associated with the current Amazon Web Services account. Policies are listed in the order that they were created in. For i...
  /// HTTP: GET /2013-04-01/trafficpolicies
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTrafficPoliciesResponse> listTrafficPolicies(
    $0.ListTrafficPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTrafficPolicies, request, options: options);
  }

  /// Gets information about the traffic policy instances that you created by using the current Amazon Web Services account. After you submit an UpdateTrafficPolicyInstance request, there's a brief delay...
  /// HTTP: GET /2013-04-01/trafficpolicyinstances
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTrafficPolicyInstancesResponse>
      listTrafficPolicyInstances(
    $0.ListTrafficPolicyInstancesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTrafficPolicyInstances, request,
        options: options);
  }

  /// Gets information about the traffic policy instances that you created in a specified hosted zone. After you submit a CreateTrafficPolicyInstance or an UpdateTrafficPolicyInstance request, there's a ...
  /// HTTP: GET /2013-04-01/trafficpolicyinstances/hostedzone
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTrafficPolicyInstancesByHostedZoneResponse>
      listTrafficPolicyInstancesByHostedZone(
    $0.ListTrafficPolicyInstancesByHostedZoneRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTrafficPolicyInstancesByHostedZone, request,
        options: options);
  }

  /// Gets information about the traffic policy instances that you created by using a specify traffic policy version. After you submit a CreateTrafficPolicyInstance or an UpdateTrafficPolicyInstance requ...
  /// HTTP: GET /2013-04-01/trafficpolicyinstances/trafficpolicy
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTrafficPolicyInstancesByPolicyResponse>
      listTrafficPolicyInstancesByPolicy(
    $0.ListTrafficPolicyInstancesByPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTrafficPolicyInstancesByPolicy, request,
        options: options);
  }

  /// Gets information about all of the versions for a specified traffic policy. Traffic policy versions are listed in numerical order by VersionNumber.
  /// HTTP: GET /2013-04-01/trafficpolicies/{Id}/versions
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListTrafficPolicyVersionsResponse>
      listTrafficPolicyVersions(
    $0.ListTrafficPolicyVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTrafficPolicyVersions, request,
        options: options);
  }

  /// Gets a list of the VPCs that were created by other accounts and that can be associated with a specified hosted zone because you've submitted one or more CreateVPCAssociationAuthorization requests. ...
  /// HTTP: GET /2013-04-01/hostedzone/{HostedZoneId}/authorizevpcassociation
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.ListVPCAssociationAuthorizationsResponse>
      listVPCAssociationAuthorizations(
    $0.ListVPCAssociationAuthorizationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listVPCAssociationAuthorizations, request,
        options: options);
  }

  /// Gets the value that Amazon Route 53 returns in response to a DNS request for a specified record name and type. You can optionally specify the IP address of a DNS resolver, an EDNS0 client subnet IP...
  /// HTTP: GET /2013-04-01/testdnsanswer
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.TestDNSAnswerResponse> testDNSAnswer(
    $0.TestDNSAnswerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testDNSAnswer, request, options: options);
  }

  /// Updates an existing health check. Note that some values can't be updated. For more information about updating health checks, see Creating, Updating, and Deleting Health Checks in the Amazon Route 5...
  /// HTTP: POST /2013-04-01/healthcheck/{HealthCheckId}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateHealthCheckResponse> updateHealthCheck(
    $0.UpdateHealthCheckRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateHealthCheck, request, options: options);
  }

  /// Updates the comment for a specified hosted zone.
  /// HTTP: POST /2013-04-01/hostedzone/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateHostedZoneCommentResponse>
      updateHostedZoneComment(
    $0.UpdateHostedZoneCommentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateHostedZoneComment, request,
        options: options);
  }

  /// Updates the features configuration for a hosted zone. This operation allows you to enable or disable specific features for your hosted zone, such as accelerated recovery. Accelerated recovery enabl...
  /// HTTP: POST /2013-04-01/hostedzone/{HostedZoneId}/features
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateHostedZoneFeaturesResponse>
      updateHostedZoneFeatures(
    $0.UpdateHostedZoneFeaturesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateHostedZoneFeatures, request,
        options: options);
  }

  /// Updates the comment for a specified traffic policy version.
  /// HTTP: POST /2013-04-01/trafficpolicy/{Id}/{Version}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateTrafficPolicyCommentResponse>
      updateTrafficPolicyComment(
    $0.UpdateTrafficPolicyCommentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTrafficPolicyComment, request,
        options: options);
  }

  /// After you submit a UpdateTrafficPolicyInstance request, there's a brief delay while Route 53 creates the resource record sets that are specified in the traffic policy definition. Use GetTrafficPol...
  /// HTTP: POST /2013-04-01/trafficpolicyinstance/{Id}
  /// Protocol: restXml
  $grpc.ResponseFuture<$0.UpdateTrafficPolicyInstanceResponse>
      updateTrafficPolicyInstance(
    $0.UpdateTrafficPolicyInstanceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTrafficPolicyInstance, request,
        options: options);
  }

  // method descriptors

  static final _$activateKeySigningKey = $grpc.ClientMethod<
          $0.ActivateKeySigningKeyRequest, $0.ActivateKeySigningKeyResponse>(
      '/route53.Route53Service/ActivateKeySigningKey',
      ($0.ActivateKeySigningKeyRequest value) => value.writeToBuffer(),
      $0.ActivateKeySigningKeyResponse.fromBuffer);
  static final _$associateVPCWithHostedZone = $grpc.ClientMethod<
          $0.AssociateVPCWithHostedZoneRequest,
          $0.AssociateVPCWithHostedZoneResponse>(
      '/route53.Route53Service/AssociateVPCWithHostedZone',
      ($0.AssociateVPCWithHostedZoneRequest value) => value.writeToBuffer(),
      $0.AssociateVPCWithHostedZoneResponse.fromBuffer);
  static final _$changeCidrCollection = $grpc.ClientMethod<
          $0.ChangeCidrCollectionRequest, $0.ChangeCidrCollectionResponse>(
      '/route53.Route53Service/ChangeCidrCollection',
      ($0.ChangeCidrCollectionRequest value) => value.writeToBuffer(),
      $0.ChangeCidrCollectionResponse.fromBuffer);
  static final _$changeResourceRecordSets = $grpc.ClientMethod<
          $0.ChangeResourceRecordSetsRequest,
          $0.ChangeResourceRecordSetsResponse>(
      '/route53.Route53Service/ChangeResourceRecordSets',
      ($0.ChangeResourceRecordSetsRequest value) => value.writeToBuffer(),
      $0.ChangeResourceRecordSetsResponse.fromBuffer);
  static final _$changeTagsForResource = $grpc.ClientMethod<
          $0.ChangeTagsForResourceRequest, $0.ChangeTagsForResourceResponse>(
      '/route53.Route53Service/ChangeTagsForResource',
      ($0.ChangeTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ChangeTagsForResourceResponse.fromBuffer);
  static final _$createCidrCollection = $grpc.ClientMethod<
          $0.CreateCidrCollectionRequest, $0.CreateCidrCollectionResponse>(
      '/route53.Route53Service/CreateCidrCollection',
      ($0.CreateCidrCollectionRequest value) => value.writeToBuffer(),
      $0.CreateCidrCollectionResponse.fromBuffer);
  static final _$createHealthCheck = $grpc.ClientMethod<
          $0.CreateHealthCheckRequest, $0.CreateHealthCheckResponse>(
      '/route53.Route53Service/CreateHealthCheck',
      ($0.CreateHealthCheckRequest value) => value.writeToBuffer(),
      $0.CreateHealthCheckResponse.fromBuffer);
  static final _$createHostedZone = $grpc.ClientMethod<
          $0.CreateHostedZoneRequest, $0.CreateHostedZoneResponse>(
      '/route53.Route53Service/CreateHostedZone',
      ($0.CreateHostedZoneRequest value) => value.writeToBuffer(),
      $0.CreateHostedZoneResponse.fromBuffer);
  static final _$createKeySigningKey = $grpc.ClientMethod<
          $0.CreateKeySigningKeyRequest, $0.CreateKeySigningKeyResponse>(
      '/route53.Route53Service/CreateKeySigningKey',
      ($0.CreateKeySigningKeyRequest value) => value.writeToBuffer(),
      $0.CreateKeySigningKeyResponse.fromBuffer);
  static final _$createQueryLoggingConfig = $grpc.ClientMethod<
          $0.CreateQueryLoggingConfigRequest,
          $0.CreateQueryLoggingConfigResponse>(
      '/route53.Route53Service/CreateQueryLoggingConfig',
      ($0.CreateQueryLoggingConfigRequest value) => value.writeToBuffer(),
      $0.CreateQueryLoggingConfigResponse.fromBuffer);
  static final _$createReusableDelegationSet = $grpc.ClientMethod<
          $0.CreateReusableDelegationSetRequest,
          $0.CreateReusableDelegationSetResponse>(
      '/route53.Route53Service/CreateReusableDelegationSet',
      ($0.CreateReusableDelegationSetRequest value) => value.writeToBuffer(),
      $0.CreateReusableDelegationSetResponse.fromBuffer);
  static final _$createTrafficPolicy = $grpc.ClientMethod<
          $0.CreateTrafficPolicyRequest, $0.CreateTrafficPolicyResponse>(
      '/route53.Route53Service/CreateTrafficPolicy',
      ($0.CreateTrafficPolicyRequest value) => value.writeToBuffer(),
      $0.CreateTrafficPolicyResponse.fromBuffer);
  static final _$createTrafficPolicyInstance = $grpc.ClientMethod<
          $0.CreateTrafficPolicyInstanceRequest,
          $0.CreateTrafficPolicyInstanceResponse>(
      '/route53.Route53Service/CreateTrafficPolicyInstance',
      ($0.CreateTrafficPolicyInstanceRequest value) => value.writeToBuffer(),
      $0.CreateTrafficPolicyInstanceResponse.fromBuffer);
  static final _$createTrafficPolicyVersion = $grpc.ClientMethod<
          $0.CreateTrafficPolicyVersionRequest,
          $0.CreateTrafficPolicyVersionResponse>(
      '/route53.Route53Service/CreateTrafficPolicyVersion',
      ($0.CreateTrafficPolicyVersionRequest value) => value.writeToBuffer(),
      $0.CreateTrafficPolicyVersionResponse.fromBuffer);
  static final _$createVPCAssociationAuthorization = $grpc.ClientMethod<
          $0.CreateVPCAssociationAuthorizationRequest,
          $0.CreateVPCAssociationAuthorizationResponse>(
      '/route53.Route53Service/CreateVPCAssociationAuthorization',
      ($0.CreateVPCAssociationAuthorizationRequest value) =>
          value.writeToBuffer(),
      $0.CreateVPCAssociationAuthorizationResponse.fromBuffer);
  static final _$deactivateKeySigningKey = $grpc.ClientMethod<
          $0.DeactivateKeySigningKeyRequest,
          $0.DeactivateKeySigningKeyResponse>(
      '/route53.Route53Service/DeactivateKeySigningKey',
      ($0.DeactivateKeySigningKeyRequest value) => value.writeToBuffer(),
      $0.DeactivateKeySigningKeyResponse.fromBuffer);
  static final _$deleteCidrCollection = $grpc.ClientMethod<
          $0.DeleteCidrCollectionRequest, $0.DeleteCidrCollectionResponse>(
      '/route53.Route53Service/DeleteCidrCollection',
      ($0.DeleteCidrCollectionRequest value) => value.writeToBuffer(),
      $0.DeleteCidrCollectionResponse.fromBuffer);
  static final _$deleteHealthCheck = $grpc.ClientMethod<
          $0.DeleteHealthCheckRequest, $0.DeleteHealthCheckResponse>(
      '/route53.Route53Service/DeleteHealthCheck',
      ($0.DeleteHealthCheckRequest value) => value.writeToBuffer(),
      $0.DeleteHealthCheckResponse.fromBuffer);
  static final _$deleteHostedZone = $grpc.ClientMethod<
          $0.DeleteHostedZoneRequest, $0.DeleteHostedZoneResponse>(
      '/route53.Route53Service/DeleteHostedZone',
      ($0.DeleteHostedZoneRequest value) => value.writeToBuffer(),
      $0.DeleteHostedZoneResponse.fromBuffer);
  static final _$deleteKeySigningKey = $grpc.ClientMethod<
          $0.DeleteKeySigningKeyRequest, $0.DeleteKeySigningKeyResponse>(
      '/route53.Route53Service/DeleteKeySigningKey',
      ($0.DeleteKeySigningKeyRequest value) => value.writeToBuffer(),
      $0.DeleteKeySigningKeyResponse.fromBuffer);
  static final _$deleteQueryLoggingConfig = $grpc.ClientMethod<
          $0.DeleteQueryLoggingConfigRequest,
          $0.DeleteQueryLoggingConfigResponse>(
      '/route53.Route53Service/DeleteQueryLoggingConfig',
      ($0.DeleteQueryLoggingConfigRequest value) => value.writeToBuffer(),
      $0.DeleteQueryLoggingConfigResponse.fromBuffer);
  static final _$deleteReusableDelegationSet = $grpc.ClientMethod<
          $0.DeleteReusableDelegationSetRequest,
          $0.DeleteReusableDelegationSetResponse>(
      '/route53.Route53Service/DeleteReusableDelegationSet',
      ($0.DeleteReusableDelegationSetRequest value) => value.writeToBuffer(),
      $0.DeleteReusableDelegationSetResponse.fromBuffer);
  static final _$deleteTrafficPolicy = $grpc.ClientMethod<
          $0.DeleteTrafficPolicyRequest, $0.DeleteTrafficPolicyResponse>(
      '/route53.Route53Service/DeleteTrafficPolicy',
      ($0.DeleteTrafficPolicyRequest value) => value.writeToBuffer(),
      $0.DeleteTrafficPolicyResponse.fromBuffer);
  static final _$deleteTrafficPolicyInstance = $grpc.ClientMethod<
          $0.DeleteTrafficPolicyInstanceRequest,
          $0.DeleteTrafficPolicyInstanceResponse>(
      '/route53.Route53Service/DeleteTrafficPolicyInstance',
      ($0.DeleteTrafficPolicyInstanceRequest value) => value.writeToBuffer(),
      $0.DeleteTrafficPolicyInstanceResponse.fromBuffer);
  static final _$deleteVPCAssociationAuthorization = $grpc.ClientMethod<
          $0.DeleteVPCAssociationAuthorizationRequest,
          $0.DeleteVPCAssociationAuthorizationResponse>(
      '/route53.Route53Service/DeleteVPCAssociationAuthorization',
      ($0.DeleteVPCAssociationAuthorizationRequest value) =>
          value.writeToBuffer(),
      $0.DeleteVPCAssociationAuthorizationResponse.fromBuffer);
  static final _$disableHostedZoneDNSSEC = $grpc.ClientMethod<
          $0.DisableHostedZoneDNSSECRequest,
          $0.DisableHostedZoneDNSSECResponse>(
      '/route53.Route53Service/DisableHostedZoneDNSSEC',
      ($0.DisableHostedZoneDNSSECRequest value) => value.writeToBuffer(),
      $0.DisableHostedZoneDNSSECResponse.fromBuffer);
  static final _$disassociateVPCFromHostedZone = $grpc.ClientMethod<
          $0.DisassociateVPCFromHostedZoneRequest,
          $0.DisassociateVPCFromHostedZoneResponse>(
      '/route53.Route53Service/DisassociateVPCFromHostedZone',
      ($0.DisassociateVPCFromHostedZoneRequest value) => value.writeToBuffer(),
      $0.DisassociateVPCFromHostedZoneResponse.fromBuffer);
  static final _$enableHostedZoneDNSSEC = $grpc.ClientMethod<
          $0.EnableHostedZoneDNSSECRequest, $0.EnableHostedZoneDNSSECResponse>(
      '/route53.Route53Service/EnableHostedZoneDNSSEC',
      ($0.EnableHostedZoneDNSSECRequest value) => value.writeToBuffer(),
      $0.EnableHostedZoneDNSSECResponse.fromBuffer);
  static final _$getAccountLimit =
      $grpc.ClientMethod<$0.GetAccountLimitRequest, $0.GetAccountLimitResponse>(
          '/route53.Route53Service/GetAccountLimit',
          ($0.GetAccountLimitRequest value) => value.writeToBuffer(),
          $0.GetAccountLimitResponse.fromBuffer);
  static final _$getChange =
      $grpc.ClientMethod<$0.GetChangeRequest, $0.GetChangeResponse>(
          '/route53.Route53Service/GetChange',
          ($0.GetChangeRequest value) => value.writeToBuffer(),
          $0.GetChangeResponse.fromBuffer);
  static final _$getCheckerIpRanges = $grpc.ClientMethod<
          $0.GetCheckerIpRangesRequest, $0.GetCheckerIpRangesResponse>(
      '/route53.Route53Service/GetCheckerIpRanges',
      ($0.GetCheckerIpRangesRequest value) => value.writeToBuffer(),
      $0.GetCheckerIpRangesResponse.fromBuffer);
  static final _$getDNSSEC =
      $grpc.ClientMethod<$0.GetDNSSECRequest, $0.GetDNSSECResponse>(
          '/route53.Route53Service/GetDNSSEC',
          ($0.GetDNSSECRequest value) => value.writeToBuffer(),
          $0.GetDNSSECResponse.fromBuffer);
  static final _$getGeoLocation =
      $grpc.ClientMethod<$0.GetGeoLocationRequest, $0.GetGeoLocationResponse>(
          '/route53.Route53Service/GetGeoLocation',
          ($0.GetGeoLocationRequest value) => value.writeToBuffer(),
          $0.GetGeoLocationResponse.fromBuffer);
  static final _$getHealthCheck =
      $grpc.ClientMethod<$0.GetHealthCheckRequest, $0.GetHealthCheckResponse>(
          '/route53.Route53Service/GetHealthCheck',
          ($0.GetHealthCheckRequest value) => value.writeToBuffer(),
          $0.GetHealthCheckResponse.fromBuffer);
  static final _$getHealthCheckCount = $grpc.ClientMethod<
          $0.GetHealthCheckCountRequest, $0.GetHealthCheckCountResponse>(
      '/route53.Route53Service/GetHealthCheckCount',
      ($0.GetHealthCheckCountRequest value) => value.writeToBuffer(),
      $0.GetHealthCheckCountResponse.fromBuffer);
  static final _$getHealthCheckLastFailureReason = $grpc.ClientMethod<
          $0.GetHealthCheckLastFailureReasonRequest,
          $0.GetHealthCheckLastFailureReasonResponse>(
      '/route53.Route53Service/GetHealthCheckLastFailureReason',
      ($0.GetHealthCheckLastFailureReasonRequest value) =>
          value.writeToBuffer(),
      $0.GetHealthCheckLastFailureReasonResponse.fromBuffer);
  static final _$getHealthCheckStatus = $grpc.ClientMethod<
          $0.GetHealthCheckStatusRequest, $0.GetHealthCheckStatusResponse>(
      '/route53.Route53Service/GetHealthCheckStatus',
      ($0.GetHealthCheckStatusRequest value) => value.writeToBuffer(),
      $0.GetHealthCheckStatusResponse.fromBuffer);
  static final _$getHostedZone =
      $grpc.ClientMethod<$0.GetHostedZoneRequest, $0.GetHostedZoneResponse>(
          '/route53.Route53Service/GetHostedZone',
          ($0.GetHostedZoneRequest value) => value.writeToBuffer(),
          $0.GetHostedZoneResponse.fromBuffer);
  static final _$getHostedZoneCount = $grpc.ClientMethod<
          $0.GetHostedZoneCountRequest, $0.GetHostedZoneCountResponse>(
      '/route53.Route53Service/GetHostedZoneCount',
      ($0.GetHostedZoneCountRequest value) => value.writeToBuffer(),
      $0.GetHostedZoneCountResponse.fromBuffer);
  static final _$getHostedZoneLimit = $grpc.ClientMethod<
          $0.GetHostedZoneLimitRequest, $0.GetHostedZoneLimitResponse>(
      '/route53.Route53Service/GetHostedZoneLimit',
      ($0.GetHostedZoneLimitRequest value) => value.writeToBuffer(),
      $0.GetHostedZoneLimitResponse.fromBuffer);
  static final _$getQueryLoggingConfig = $grpc.ClientMethod<
          $0.GetQueryLoggingConfigRequest, $0.GetQueryLoggingConfigResponse>(
      '/route53.Route53Service/GetQueryLoggingConfig',
      ($0.GetQueryLoggingConfigRequest value) => value.writeToBuffer(),
      $0.GetQueryLoggingConfigResponse.fromBuffer);
  static final _$getReusableDelegationSet = $grpc.ClientMethod<
          $0.GetReusableDelegationSetRequest,
          $0.GetReusableDelegationSetResponse>(
      '/route53.Route53Service/GetReusableDelegationSet',
      ($0.GetReusableDelegationSetRequest value) => value.writeToBuffer(),
      $0.GetReusableDelegationSetResponse.fromBuffer);
  static final _$getReusableDelegationSetLimit = $grpc.ClientMethod<
          $0.GetReusableDelegationSetLimitRequest,
          $0.GetReusableDelegationSetLimitResponse>(
      '/route53.Route53Service/GetReusableDelegationSetLimit',
      ($0.GetReusableDelegationSetLimitRequest value) => value.writeToBuffer(),
      $0.GetReusableDelegationSetLimitResponse.fromBuffer);
  static final _$getTrafficPolicy = $grpc.ClientMethod<
          $0.GetTrafficPolicyRequest, $0.GetTrafficPolicyResponse>(
      '/route53.Route53Service/GetTrafficPolicy',
      ($0.GetTrafficPolicyRequest value) => value.writeToBuffer(),
      $0.GetTrafficPolicyResponse.fromBuffer);
  static final _$getTrafficPolicyInstance = $grpc.ClientMethod<
          $0.GetTrafficPolicyInstanceRequest,
          $0.GetTrafficPolicyInstanceResponse>(
      '/route53.Route53Service/GetTrafficPolicyInstance',
      ($0.GetTrafficPolicyInstanceRequest value) => value.writeToBuffer(),
      $0.GetTrafficPolicyInstanceResponse.fromBuffer);
  static final _$getTrafficPolicyInstanceCount = $grpc.ClientMethod<
          $0.GetTrafficPolicyInstanceCountRequest,
          $0.GetTrafficPolicyInstanceCountResponse>(
      '/route53.Route53Service/GetTrafficPolicyInstanceCount',
      ($0.GetTrafficPolicyInstanceCountRequest value) => value.writeToBuffer(),
      $0.GetTrafficPolicyInstanceCountResponse.fromBuffer);
  static final _$listCidrBlocks =
      $grpc.ClientMethod<$0.ListCidrBlocksRequest, $0.ListCidrBlocksResponse>(
          '/route53.Route53Service/ListCidrBlocks',
          ($0.ListCidrBlocksRequest value) => value.writeToBuffer(),
          $0.ListCidrBlocksResponse.fromBuffer);
  static final _$listCidrCollections = $grpc.ClientMethod<
          $0.ListCidrCollectionsRequest, $0.ListCidrCollectionsResponse>(
      '/route53.Route53Service/ListCidrCollections',
      ($0.ListCidrCollectionsRequest value) => value.writeToBuffer(),
      $0.ListCidrCollectionsResponse.fromBuffer);
  static final _$listCidrLocations = $grpc.ClientMethod<
          $0.ListCidrLocationsRequest, $0.ListCidrLocationsResponse>(
      '/route53.Route53Service/ListCidrLocations',
      ($0.ListCidrLocationsRequest value) => value.writeToBuffer(),
      $0.ListCidrLocationsResponse.fromBuffer);
  static final _$listGeoLocations = $grpc.ClientMethod<
          $0.ListGeoLocationsRequest, $0.ListGeoLocationsResponse>(
      '/route53.Route53Service/ListGeoLocations',
      ($0.ListGeoLocationsRequest value) => value.writeToBuffer(),
      $0.ListGeoLocationsResponse.fromBuffer);
  static final _$listHealthChecks = $grpc.ClientMethod<
          $0.ListHealthChecksRequest, $0.ListHealthChecksResponse>(
      '/route53.Route53Service/ListHealthChecks',
      ($0.ListHealthChecksRequest value) => value.writeToBuffer(),
      $0.ListHealthChecksResponse.fromBuffer);
  static final _$listHostedZones =
      $grpc.ClientMethod<$0.ListHostedZonesRequest, $0.ListHostedZonesResponse>(
          '/route53.Route53Service/ListHostedZones',
          ($0.ListHostedZonesRequest value) => value.writeToBuffer(),
          $0.ListHostedZonesResponse.fromBuffer);
  static final _$listHostedZonesByName = $grpc.ClientMethod<
          $0.ListHostedZonesByNameRequest, $0.ListHostedZonesByNameResponse>(
      '/route53.Route53Service/ListHostedZonesByName',
      ($0.ListHostedZonesByNameRequest value) => value.writeToBuffer(),
      $0.ListHostedZonesByNameResponse.fromBuffer);
  static final _$listHostedZonesByVPC = $grpc.ClientMethod<
          $0.ListHostedZonesByVPCRequest, $0.ListHostedZonesByVPCResponse>(
      '/route53.Route53Service/ListHostedZonesByVPC',
      ($0.ListHostedZonesByVPCRequest value) => value.writeToBuffer(),
      $0.ListHostedZonesByVPCResponse.fromBuffer);
  static final _$listQueryLoggingConfigs = $grpc.ClientMethod<
          $0.ListQueryLoggingConfigsRequest,
          $0.ListQueryLoggingConfigsResponse>(
      '/route53.Route53Service/ListQueryLoggingConfigs',
      ($0.ListQueryLoggingConfigsRequest value) => value.writeToBuffer(),
      $0.ListQueryLoggingConfigsResponse.fromBuffer);
  static final _$listResourceRecordSets = $grpc.ClientMethod<
          $0.ListResourceRecordSetsRequest, $0.ListResourceRecordSetsResponse>(
      '/route53.Route53Service/ListResourceRecordSets',
      ($0.ListResourceRecordSetsRequest value) => value.writeToBuffer(),
      $0.ListResourceRecordSetsResponse.fromBuffer);
  static final _$listReusableDelegationSets = $grpc.ClientMethod<
          $0.ListReusableDelegationSetsRequest,
          $0.ListReusableDelegationSetsResponse>(
      '/route53.Route53Service/ListReusableDelegationSets',
      ($0.ListReusableDelegationSetsRequest value) => value.writeToBuffer(),
      $0.ListReusableDelegationSetsResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/route53.Route53Service/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listTagsForResources = $grpc.ClientMethod<
          $0.ListTagsForResourcesRequest, $0.ListTagsForResourcesResponse>(
      '/route53.Route53Service/ListTagsForResources',
      ($0.ListTagsForResourcesRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourcesResponse.fromBuffer);
  static final _$listTrafficPolicies = $grpc.ClientMethod<
          $0.ListTrafficPoliciesRequest, $0.ListTrafficPoliciesResponse>(
      '/route53.Route53Service/ListTrafficPolicies',
      ($0.ListTrafficPoliciesRequest value) => value.writeToBuffer(),
      $0.ListTrafficPoliciesResponse.fromBuffer);
  static final _$listTrafficPolicyInstances = $grpc.ClientMethod<
          $0.ListTrafficPolicyInstancesRequest,
          $0.ListTrafficPolicyInstancesResponse>(
      '/route53.Route53Service/ListTrafficPolicyInstances',
      ($0.ListTrafficPolicyInstancesRequest value) => value.writeToBuffer(),
      $0.ListTrafficPolicyInstancesResponse.fromBuffer);
  static final _$listTrafficPolicyInstancesByHostedZone = $grpc.ClientMethod<
          $0.ListTrafficPolicyInstancesByHostedZoneRequest,
          $0.ListTrafficPolicyInstancesByHostedZoneResponse>(
      '/route53.Route53Service/ListTrafficPolicyInstancesByHostedZone',
      ($0.ListTrafficPolicyInstancesByHostedZoneRequest value) =>
          value.writeToBuffer(),
      $0.ListTrafficPolicyInstancesByHostedZoneResponse.fromBuffer);
  static final _$listTrafficPolicyInstancesByPolicy = $grpc.ClientMethod<
          $0.ListTrafficPolicyInstancesByPolicyRequest,
          $0.ListTrafficPolicyInstancesByPolicyResponse>(
      '/route53.Route53Service/ListTrafficPolicyInstancesByPolicy',
      ($0.ListTrafficPolicyInstancesByPolicyRequest value) =>
          value.writeToBuffer(),
      $0.ListTrafficPolicyInstancesByPolicyResponse.fromBuffer);
  static final _$listTrafficPolicyVersions = $grpc.ClientMethod<
          $0.ListTrafficPolicyVersionsRequest,
          $0.ListTrafficPolicyVersionsResponse>(
      '/route53.Route53Service/ListTrafficPolicyVersions',
      ($0.ListTrafficPolicyVersionsRequest value) => value.writeToBuffer(),
      $0.ListTrafficPolicyVersionsResponse.fromBuffer);
  static final _$listVPCAssociationAuthorizations = $grpc.ClientMethod<
          $0.ListVPCAssociationAuthorizationsRequest,
          $0.ListVPCAssociationAuthorizationsResponse>(
      '/route53.Route53Service/ListVPCAssociationAuthorizations',
      ($0.ListVPCAssociationAuthorizationsRequest value) =>
          value.writeToBuffer(),
      $0.ListVPCAssociationAuthorizationsResponse.fromBuffer);
  static final _$testDNSAnswer =
      $grpc.ClientMethod<$0.TestDNSAnswerRequest, $0.TestDNSAnswerResponse>(
          '/route53.Route53Service/TestDNSAnswer',
          ($0.TestDNSAnswerRequest value) => value.writeToBuffer(),
          $0.TestDNSAnswerResponse.fromBuffer);
  static final _$updateHealthCheck = $grpc.ClientMethod<
          $0.UpdateHealthCheckRequest, $0.UpdateHealthCheckResponse>(
      '/route53.Route53Service/UpdateHealthCheck',
      ($0.UpdateHealthCheckRequest value) => value.writeToBuffer(),
      $0.UpdateHealthCheckResponse.fromBuffer);
  static final _$updateHostedZoneComment = $grpc.ClientMethod<
          $0.UpdateHostedZoneCommentRequest,
          $0.UpdateHostedZoneCommentResponse>(
      '/route53.Route53Service/UpdateHostedZoneComment',
      ($0.UpdateHostedZoneCommentRequest value) => value.writeToBuffer(),
      $0.UpdateHostedZoneCommentResponse.fromBuffer);
  static final _$updateHostedZoneFeatures = $grpc.ClientMethod<
          $0.UpdateHostedZoneFeaturesRequest,
          $0.UpdateHostedZoneFeaturesResponse>(
      '/route53.Route53Service/UpdateHostedZoneFeatures',
      ($0.UpdateHostedZoneFeaturesRequest value) => value.writeToBuffer(),
      $0.UpdateHostedZoneFeaturesResponse.fromBuffer);
  static final _$updateTrafficPolicyComment = $grpc.ClientMethod<
          $0.UpdateTrafficPolicyCommentRequest,
          $0.UpdateTrafficPolicyCommentResponse>(
      '/route53.Route53Service/UpdateTrafficPolicyComment',
      ($0.UpdateTrafficPolicyCommentRequest value) => value.writeToBuffer(),
      $0.UpdateTrafficPolicyCommentResponse.fromBuffer);
  static final _$updateTrafficPolicyInstance = $grpc.ClientMethod<
          $0.UpdateTrafficPolicyInstanceRequest,
          $0.UpdateTrafficPolicyInstanceResponse>(
      '/route53.Route53Service/UpdateTrafficPolicyInstance',
      ($0.UpdateTrafficPolicyInstanceRequest value) => value.writeToBuffer(),
      $0.UpdateTrafficPolicyInstanceResponse.fromBuffer);
}

@$pb.GrpcServiceName('route53.Route53Service')
abstract class Route53ServiceBase extends $grpc.Service {
  $core.String get $name => 'route53.Route53Service';

  Route53ServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.ActivateKeySigningKeyRequest,
            $0.ActivateKeySigningKeyResponse>(
        'ActivateKeySigningKey',
        activateKeySigningKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ActivateKeySigningKeyRequest.fromBuffer(value),
        ($0.ActivateKeySigningKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssociateVPCWithHostedZoneRequest,
            $0.AssociateVPCWithHostedZoneResponse>(
        'AssociateVPCWithHostedZone',
        associateVPCWithHostedZone_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateVPCWithHostedZoneRequest.fromBuffer(value),
        ($0.AssociateVPCWithHostedZoneResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangeCidrCollectionRequest,
            $0.ChangeCidrCollectionResponse>(
        'ChangeCidrCollection',
        changeCidrCollection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangeCidrCollectionRequest.fromBuffer(value),
        ($0.ChangeCidrCollectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangeResourceRecordSetsRequest,
            $0.ChangeResourceRecordSetsResponse>(
        'ChangeResourceRecordSets',
        changeResourceRecordSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangeResourceRecordSetsRequest.fromBuffer(value),
        ($0.ChangeResourceRecordSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangeTagsForResourceRequest,
            $0.ChangeTagsForResourceResponse>(
        'ChangeTagsForResource',
        changeTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangeTagsForResourceRequest.fromBuffer(value),
        ($0.ChangeTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateCidrCollectionRequest,
            $0.CreateCidrCollectionResponse>(
        'CreateCidrCollection',
        createCidrCollection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCidrCollectionRequest.fromBuffer(value),
        ($0.CreateCidrCollectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateHealthCheckRequest,
            $0.CreateHealthCheckResponse>(
        'CreateHealthCheck',
        createHealthCheck_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateHealthCheckRequest.fromBuffer(value),
        ($0.CreateHealthCheckResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateHostedZoneRequest,
            $0.CreateHostedZoneResponse>(
        'CreateHostedZone',
        createHostedZone_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateHostedZoneRequest.fromBuffer(value),
        ($0.CreateHostedZoneResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateKeySigningKeyRequest,
            $0.CreateKeySigningKeyResponse>(
        'CreateKeySigningKey',
        createKeySigningKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateKeySigningKeyRequest.fromBuffer(value),
        ($0.CreateKeySigningKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateQueryLoggingConfigRequest,
            $0.CreateQueryLoggingConfigResponse>(
        'CreateQueryLoggingConfig',
        createQueryLoggingConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateQueryLoggingConfigRequest.fromBuffer(value),
        ($0.CreateQueryLoggingConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateReusableDelegationSetRequest,
            $0.CreateReusableDelegationSetResponse>(
        'CreateReusableDelegationSet',
        createReusableDelegationSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateReusableDelegationSetRequest.fromBuffer(value),
        ($0.CreateReusableDelegationSetResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateTrafficPolicyRequest,
            $0.CreateTrafficPolicyResponse>(
        'CreateTrafficPolicy',
        createTrafficPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateTrafficPolicyRequest.fromBuffer(value),
        ($0.CreateTrafficPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateTrafficPolicyInstanceRequest,
            $0.CreateTrafficPolicyInstanceResponse>(
        'CreateTrafficPolicyInstance',
        createTrafficPolicyInstance_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateTrafficPolicyInstanceRequest.fromBuffer(value),
        ($0.CreateTrafficPolicyInstanceResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateTrafficPolicyVersionRequest,
            $0.CreateTrafficPolicyVersionResponse>(
        'CreateTrafficPolicyVersion',
        createTrafficPolicyVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateTrafficPolicyVersionRequest.fromBuffer(value),
        ($0.CreateTrafficPolicyVersionResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateVPCAssociationAuthorizationRequest,
            $0.CreateVPCAssociationAuthorizationResponse>(
        'CreateVPCAssociationAuthorization',
        createVPCAssociationAuthorization_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateVPCAssociationAuthorizationRequest.fromBuffer(value),
        ($0.CreateVPCAssociationAuthorizationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeactivateKeySigningKeyRequest,
            $0.DeactivateKeySigningKeyResponse>(
        'DeactivateKeySigningKey',
        deactivateKeySigningKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeactivateKeySigningKeyRequest.fromBuffer(value),
        ($0.DeactivateKeySigningKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteCidrCollectionRequest,
            $0.DeleteCidrCollectionResponse>(
        'DeleteCidrCollection',
        deleteCidrCollection_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCidrCollectionRequest.fromBuffer(value),
        ($0.DeleteCidrCollectionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteHealthCheckRequest,
            $0.DeleteHealthCheckResponse>(
        'DeleteHealthCheck',
        deleteHealthCheck_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteHealthCheckRequest.fromBuffer(value),
        ($0.DeleteHealthCheckResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteHostedZoneRequest,
            $0.DeleteHostedZoneResponse>(
        'DeleteHostedZone',
        deleteHostedZone_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteHostedZoneRequest.fromBuffer(value),
        ($0.DeleteHostedZoneResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteKeySigningKeyRequest,
            $0.DeleteKeySigningKeyResponse>(
        'DeleteKeySigningKey',
        deleteKeySigningKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteKeySigningKeyRequest.fromBuffer(value),
        ($0.DeleteKeySigningKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteQueryLoggingConfigRequest,
            $0.DeleteQueryLoggingConfigResponse>(
        'DeleteQueryLoggingConfig',
        deleteQueryLoggingConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteQueryLoggingConfigRequest.fromBuffer(value),
        ($0.DeleteQueryLoggingConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteReusableDelegationSetRequest,
            $0.DeleteReusableDelegationSetResponse>(
        'DeleteReusableDelegationSet',
        deleteReusableDelegationSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteReusableDelegationSetRequest.fromBuffer(value),
        ($0.DeleteReusableDelegationSetResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTrafficPolicyRequest,
            $0.DeleteTrafficPolicyResponse>(
        'DeleteTrafficPolicy',
        deleteTrafficPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteTrafficPolicyRequest.fromBuffer(value),
        ($0.DeleteTrafficPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTrafficPolicyInstanceRequest,
            $0.DeleteTrafficPolicyInstanceResponse>(
        'DeleteTrafficPolicyInstance',
        deleteTrafficPolicyInstance_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteTrafficPolicyInstanceRequest.fromBuffer(value),
        ($0.DeleteTrafficPolicyInstanceResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteVPCAssociationAuthorizationRequest,
            $0.DeleteVPCAssociationAuthorizationResponse>(
        'DeleteVPCAssociationAuthorization',
        deleteVPCAssociationAuthorization_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteVPCAssociationAuthorizationRequest.fromBuffer(value),
        ($0.DeleteVPCAssociationAuthorizationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableHostedZoneDNSSECRequest,
            $0.DisableHostedZoneDNSSECResponse>(
        'DisableHostedZoneDNSSEC',
        disableHostedZoneDNSSEC_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableHostedZoneDNSSECRequest.fromBuffer(value),
        ($0.DisableHostedZoneDNSSECResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisassociateVPCFromHostedZoneRequest,
            $0.DisassociateVPCFromHostedZoneResponse>(
        'DisassociateVPCFromHostedZone',
        disassociateVPCFromHostedZone_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisassociateVPCFromHostedZoneRequest.fromBuffer(value),
        ($0.DisassociateVPCFromHostedZoneResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableHostedZoneDNSSECRequest,
            $0.EnableHostedZoneDNSSECResponse>(
        'EnableHostedZoneDNSSEC',
        enableHostedZoneDNSSEC_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableHostedZoneDNSSECRequest.fromBuffer(value),
        ($0.EnableHostedZoneDNSSECResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccountLimitRequest,
            $0.GetAccountLimitResponse>(
        'GetAccountLimit',
        getAccountLimit_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAccountLimitRequest.fromBuffer(value),
        ($0.GetAccountLimitResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetChangeRequest, $0.GetChangeResponse>(
        'GetChange',
        getChange_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetChangeRequest.fromBuffer(value),
        ($0.GetChangeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCheckerIpRangesRequest,
            $0.GetCheckerIpRangesResponse>(
        'GetCheckerIpRanges',
        getCheckerIpRanges_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCheckerIpRangesRequest.fromBuffer(value),
        ($0.GetCheckerIpRangesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDNSSECRequest, $0.GetDNSSECResponse>(
        'GetDNSSEC',
        getDNSSEC_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetDNSSECRequest.fromBuffer(value),
        ($0.GetDNSSECResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetGeoLocationRequest,
            $0.GetGeoLocationResponse>(
        'GetGeoLocation',
        getGeoLocation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetGeoLocationRequest.fromBuffer(value),
        ($0.GetGeoLocationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHealthCheckRequest,
            $0.GetHealthCheckResponse>(
        'GetHealthCheck',
        getHealthCheck_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHealthCheckRequest.fromBuffer(value),
        ($0.GetHealthCheckResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHealthCheckCountRequest,
            $0.GetHealthCheckCountResponse>(
        'GetHealthCheckCount',
        getHealthCheckCount_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHealthCheckCountRequest.fromBuffer(value),
        ($0.GetHealthCheckCountResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHealthCheckLastFailureReasonRequest,
            $0.GetHealthCheckLastFailureReasonResponse>(
        'GetHealthCheckLastFailureReason',
        getHealthCheckLastFailureReason_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHealthCheckLastFailureReasonRequest.fromBuffer(value),
        ($0.GetHealthCheckLastFailureReasonResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHealthCheckStatusRequest,
            $0.GetHealthCheckStatusResponse>(
        'GetHealthCheckStatus',
        getHealthCheckStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHealthCheckStatusRequest.fromBuffer(value),
        ($0.GetHealthCheckStatusResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetHostedZoneRequest, $0.GetHostedZoneResponse>(
            'GetHostedZone',
            getHostedZone_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetHostedZoneRequest.fromBuffer(value),
            ($0.GetHostedZoneResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHostedZoneCountRequest,
            $0.GetHostedZoneCountResponse>(
        'GetHostedZoneCount',
        getHostedZoneCount_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHostedZoneCountRequest.fromBuffer(value),
        ($0.GetHostedZoneCountResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHostedZoneLimitRequest,
            $0.GetHostedZoneLimitResponse>(
        'GetHostedZoneLimit',
        getHostedZoneLimit_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHostedZoneLimitRequest.fromBuffer(value),
        ($0.GetHostedZoneLimitResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetQueryLoggingConfigRequest,
            $0.GetQueryLoggingConfigResponse>(
        'GetQueryLoggingConfig',
        getQueryLoggingConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetQueryLoggingConfigRequest.fromBuffer(value),
        ($0.GetQueryLoggingConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetReusableDelegationSetRequest,
            $0.GetReusableDelegationSetResponse>(
        'GetReusableDelegationSet',
        getReusableDelegationSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetReusableDelegationSetRequest.fromBuffer(value),
        ($0.GetReusableDelegationSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetReusableDelegationSetLimitRequest,
            $0.GetReusableDelegationSetLimitResponse>(
        'GetReusableDelegationSetLimit',
        getReusableDelegationSetLimit_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetReusableDelegationSetLimitRequest.fromBuffer(value),
        ($0.GetReusableDelegationSetLimitResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTrafficPolicyRequest,
            $0.GetTrafficPolicyResponse>(
        'GetTrafficPolicy',
        getTrafficPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTrafficPolicyRequest.fromBuffer(value),
        ($0.GetTrafficPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTrafficPolicyInstanceRequest,
            $0.GetTrafficPolicyInstanceResponse>(
        'GetTrafficPolicyInstance',
        getTrafficPolicyInstance_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTrafficPolicyInstanceRequest.fromBuffer(value),
        ($0.GetTrafficPolicyInstanceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTrafficPolicyInstanceCountRequest,
            $0.GetTrafficPolicyInstanceCountResponse>(
        'GetTrafficPolicyInstanceCount',
        getTrafficPolicyInstanceCount_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTrafficPolicyInstanceCountRequest.fromBuffer(value),
        ($0.GetTrafficPolicyInstanceCountResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCidrBlocksRequest,
            $0.ListCidrBlocksResponse>(
        'ListCidrBlocks',
        listCidrBlocks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCidrBlocksRequest.fromBuffer(value),
        ($0.ListCidrBlocksResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCidrCollectionsRequest,
            $0.ListCidrCollectionsResponse>(
        'ListCidrCollections',
        listCidrCollections_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCidrCollectionsRequest.fromBuffer(value),
        ($0.ListCidrCollectionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCidrLocationsRequest,
            $0.ListCidrLocationsResponse>(
        'ListCidrLocations',
        listCidrLocations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCidrLocationsRequest.fromBuffer(value),
        ($0.ListCidrLocationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGeoLocationsRequest,
            $0.ListGeoLocationsResponse>(
        'ListGeoLocations',
        listGeoLocations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListGeoLocationsRequest.fromBuffer(value),
        ($0.ListGeoLocationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListHealthChecksRequest,
            $0.ListHealthChecksResponse>(
        'ListHealthChecks',
        listHealthChecks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListHealthChecksRequest.fromBuffer(value),
        ($0.ListHealthChecksResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListHostedZonesRequest,
            $0.ListHostedZonesResponse>(
        'ListHostedZones',
        listHostedZones_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListHostedZonesRequest.fromBuffer(value),
        ($0.ListHostedZonesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListHostedZonesByNameRequest,
            $0.ListHostedZonesByNameResponse>(
        'ListHostedZonesByName',
        listHostedZonesByName_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListHostedZonesByNameRequest.fromBuffer(value),
        ($0.ListHostedZonesByNameResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListHostedZonesByVPCRequest,
            $0.ListHostedZonesByVPCResponse>(
        'ListHostedZonesByVPC',
        listHostedZonesByVPC_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListHostedZonesByVPCRequest.fromBuffer(value),
        ($0.ListHostedZonesByVPCResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListQueryLoggingConfigsRequest,
            $0.ListQueryLoggingConfigsResponse>(
        'ListQueryLoggingConfigs',
        listQueryLoggingConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListQueryLoggingConfigsRequest.fromBuffer(value),
        ($0.ListQueryLoggingConfigsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResourceRecordSetsRequest,
            $0.ListResourceRecordSetsResponse>(
        'ListResourceRecordSets',
        listResourceRecordSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResourceRecordSetsRequest.fromBuffer(value),
        ($0.ListResourceRecordSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListReusableDelegationSetsRequest,
            $0.ListReusableDelegationSetsResponse>(
        'ListReusableDelegationSets',
        listReusableDelegationSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListReusableDelegationSetsRequest.fromBuffer(value),
        ($0.ListReusableDelegationSetsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourcesRequest,
            $0.ListTagsForResourcesResponse>(
        'ListTagsForResources',
        listTagsForResources_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourcesRequest.fromBuffer(value),
        ($0.ListTagsForResourcesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTrafficPoliciesRequest,
            $0.ListTrafficPoliciesResponse>(
        'ListTrafficPolicies',
        listTrafficPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTrafficPoliciesRequest.fromBuffer(value),
        ($0.ListTrafficPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTrafficPolicyInstancesRequest,
            $0.ListTrafficPolicyInstancesResponse>(
        'ListTrafficPolicyInstances',
        listTrafficPolicyInstances_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTrafficPolicyInstancesRequest.fromBuffer(value),
        ($0.ListTrafficPolicyInstancesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListTrafficPolicyInstancesByHostedZoneRequest,
            $0.ListTrafficPolicyInstancesByHostedZoneResponse>(
        'ListTrafficPolicyInstancesByHostedZone',
        listTrafficPolicyInstancesByHostedZone_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTrafficPolicyInstancesByHostedZoneRequest.fromBuffer(value),
        ($0.ListTrafficPolicyInstancesByHostedZoneResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTrafficPolicyInstancesByPolicyRequest,
            $0.ListTrafficPolicyInstancesByPolicyResponse>(
        'ListTrafficPolicyInstancesByPolicy',
        listTrafficPolicyInstancesByPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTrafficPolicyInstancesByPolicyRequest.fromBuffer(value),
        ($0.ListTrafficPolicyInstancesByPolicyResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTrafficPolicyVersionsRequest,
            $0.ListTrafficPolicyVersionsResponse>(
        'ListTrafficPolicyVersions',
        listTrafficPolicyVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTrafficPolicyVersionsRequest.fromBuffer(value),
        ($0.ListTrafficPolicyVersionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListVPCAssociationAuthorizationsRequest,
            $0.ListVPCAssociationAuthorizationsResponse>(
        'ListVPCAssociationAuthorizations',
        listVPCAssociationAuthorizations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListVPCAssociationAuthorizationsRequest.fromBuffer(value),
        ($0.ListVPCAssociationAuthorizationsResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TestDNSAnswerRequest, $0.TestDNSAnswerResponse>(
            'TestDNSAnswer',
            testDNSAnswer_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TestDNSAnswerRequest.fromBuffer(value),
            ($0.TestDNSAnswerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateHealthCheckRequest,
            $0.UpdateHealthCheckResponse>(
        'UpdateHealthCheck',
        updateHealthCheck_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateHealthCheckRequest.fromBuffer(value),
        ($0.UpdateHealthCheckResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateHostedZoneCommentRequest,
            $0.UpdateHostedZoneCommentResponse>(
        'UpdateHostedZoneComment',
        updateHostedZoneComment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateHostedZoneCommentRequest.fromBuffer(value),
        ($0.UpdateHostedZoneCommentResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateHostedZoneFeaturesRequest,
            $0.UpdateHostedZoneFeaturesResponse>(
        'UpdateHostedZoneFeatures',
        updateHostedZoneFeatures_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateHostedZoneFeaturesRequest.fromBuffer(value),
        ($0.UpdateHostedZoneFeaturesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateTrafficPolicyCommentRequest,
            $0.UpdateTrafficPolicyCommentResponse>(
        'UpdateTrafficPolicyComment',
        updateTrafficPolicyComment_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateTrafficPolicyCommentRequest.fromBuffer(value),
        ($0.UpdateTrafficPolicyCommentResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateTrafficPolicyInstanceRequest,
            $0.UpdateTrafficPolicyInstanceResponse>(
        'UpdateTrafficPolicyInstance',
        updateTrafficPolicyInstance_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateTrafficPolicyInstanceRequest.fromBuffer(value),
        ($0.UpdateTrafficPolicyInstanceResponse value) =>
            value.writeToBuffer()));
  }

  $async.Future<$0.ActivateKeySigningKeyResponse> activateKeySigningKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ActivateKeySigningKeyRequest> $request) async {
    return activateKeySigningKey($call, await $request);
  }

  $async.Future<$0.ActivateKeySigningKeyResponse> activateKeySigningKey(
      $grpc.ServiceCall call, $0.ActivateKeySigningKeyRequest request);

  $async.Future<$0.AssociateVPCWithHostedZoneResponse>
      associateVPCWithHostedZone_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AssociateVPCWithHostedZoneRequest> $request) async {
    return associateVPCWithHostedZone($call, await $request);
  }

  $async.Future<$0.AssociateVPCWithHostedZoneResponse>
      associateVPCWithHostedZone(
          $grpc.ServiceCall call, $0.AssociateVPCWithHostedZoneRequest request);

  $async.Future<$0.ChangeCidrCollectionResponse> changeCidrCollection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ChangeCidrCollectionRequest> $request) async {
    return changeCidrCollection($call, await $request);
  }

  $async.Future<$0.ChangeCidrCollectionResponse> changeCidrCollection(
      $grpc.ServiceCall call, $0.ChangeCidrCollectionRequest request);

  $async.Future<$0.ChangeResourceRecordSetsResponse>
      changeResourceRecordSets_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ChangeResourceRecordSetsRequest> $request) async {
    return changeResourceRecordSets($call, await $request);
  }

  $async.Future<$0.ChangeResourceRecordSetsResponse> changeResourceRecordSets(
      $grpc.ServiceCall call, $0.ChangeResourceRecordSetsRequest request);

  $async.Future<$0.ChangeTagsForResourceResponse> changeTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ChangeTagsForResourceRequest> $request) async {
    return changeTagsForResource($call, await $request);
  }

  $async.Future<$0.ChangeTagsForResourceResponse> changeTagsForResource(
      $grpc.ServiceCall call, $0.ChangeTagsForResourceRequest request);

  $async.Future<$0.CreateCidrCollectionResponse> createCidrCollection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateCidrCollectionRequest> $request) async {
    return createCidrCollection($call, await $request);
  }

  $async.Future<$0.CreateCidrCollectionResponse> createCidrCollection(
      $grpc.ServiceCall call, $0.CreateCidrCollectionRequest request);

  $async.Future<$0.CreateHealthCheckResponse> createHealthCheck_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateHealthCheckRequest> $request) async {
    return createHealthCheck($call, await $request);
  }

  $async.Future<$0.CreateHealthCheckResponse> createHealthCheck(
      $grpc.ServiceCall call, $0.CreateHealthCheckRequest request);

  $async.Future<$0.CreateHostedZoneResponse> createHostedZone_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateHostedZoneRequest> $request) async {
    return createHostedZone($call, await $request);
  }

  $async.Future<$0.CreateHostedZoneResponse> createHostedZone(
      $grpc.ServiceCall call, $0.CreateHostedZoneRequest request);

  $async.Future<$0.CreateKeySigningKeyResponse> createKeySigningKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateKeySigningKeyRequest> $request) async {
    return createKeySigningKey($call, await $request);
  }

  $async.Future<$0.CreateKeySigningKeyResponse> createKeySigningKey(
      $grpc.ServiceCall call, $0.CreateKeySigningKeyRequest request);

  $async.Future<$0.CreateQueryLoggingConfigResponse>
      createQueryLoggingConfig_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateQueryLoggingConfigRequest> $request) async {
    return createQueryLoggingConfig($call, await $request);
  }

  $async.Future<$0.CreateQueryLoggingConfigResponse> createQueryLoggingConfig(
      $grpc.ServiceCall call, $0.CreateQueryLoggingConfigRequest request);

  $async.Future<$0.CreateReusableDelegationSetResponse>
      createReusableDelegationSet_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateReusableDelegationSetRequest> $request) async {
    return createReusableDelegationSet($call, await $request);
  }

  $async.Future<$0.CreateReusableDelegationSetResponse>
      createReusableDelegationSet($grpc.ServiceCall call,
          $0.CreateReusableDelegationSetRequest request);

  $async.Future<$0.CreateTrafficPolicyResponse> createTrafficPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateTrafficPolicyRequest> $request) async {
    return createTrafficPolicy($call, await $request);
  }

  $async.Future<$0.CreateTrafficPolicyResponse> createTrafficPolicy(
      $grpc.ServiceCall call, $0.CreateTrafficPolicyRequest request);

  $async.Future<$0.CreateTrafficPolicyInstanceResponse>
      createTrafficPolicyInstance_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateTrafficPolicyInstanceRequest> $request) async {
    return createTrafficPolicyInstance($call, await $request);
  }

  $async.Future<$0.CreateTrafficPolicyInstanceResponse>
      createTrafficPolicyInstance($grpc.ServiceCall call,
          $0.CreateTrafficPolicyInstanceRequest request);

  $async.Future<$0.CreateTrafficPolicyVersionResponse>
      createTrafficPolicyVersion_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateTrafficPolicyVersionRequest> $request) async {
    return createTrafficPolicyVersion($call, await $request);
  }

  $async.Future<$0.CreateTrafficPolicyVersionResponse>
      createTrafficPolicyVersion(
          $grpc.ServiceCall call, $0.CreateTrafficPolicyVersionRequest request);

  $async.Future<$0.CreateVPCAssociationAuthorizationResponse>
      createVPCAssociationAuthorization_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateVPCAssociationAuthorizationRequest>
              $request) async {
    return createVPCAssociationAuthorization($call, await $request);
  }

  $async.Future<$0.CreateVPCAssociationAuthorizationResponse>
      createVPCAssociationAuthorization($grpc.ServiceCall call,
          $0.CreateVPCAssociationAuthorizationRequest request);

  $async.Future<$0.DeactivateKeySigningKeyResponse> deactivateKeySigningKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeactivateKeySigningKeyRequest> $request) async {
    return deactivateKeySigningKey($call, await $request);
  }

  $async.Future<$0.DeactivateKeySigningKeyResponse> deactivateKeySigningKey(
      $grpc.ServiceCall call, $0.DeactivateKeySigningKeyRequest request);

  $async.Future<$0.DeleteCidrCollectionResponse> deleteCidrCollection_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteCidrCollectionRequest> $request) async {
    return deleteCidrCollection($call, await $request);
  }

  $async.Future<$0.DeleteCidrCollectionResponse> deleteCidrCollection(
      $grpc.ServiceCall call, $0.DeleteCidrCollectionRequest request);

  $async.Future<$0.DeleteHealthCheckResponse> deleteHealthCheck_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteHealthCheckRequest> $request) async {
    return deleteHealthCheck($call, await $request);
  }

  $async.Future<$0.DeleteHealthCheckResponse> deleteHealthCheck(
      $grpc.ServiceCall call, $0.DeleteHealthCheckRequest request);

  $async.Future<$0.DeleteHostedZoneResponse> deleteHostedZone_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteHostedZoneRequest> $request) async {
    return deleteHostedZone($call, await $request);
  }

  $async.Future<$0.DeleteHostedZoneResponse> deleteHostedZone(
      $grpc.ServiceCall call, $0.DeleteHostedZoneRequest request);

  $async.Future<$0.DeleteKeySigningKeyResponse> deleteKeySigningKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteKeySigningKeyRequest> $request) async {
    return deleteKeySigningKey($call, await $request);
  }

  $async.Future<$0.DeleteKeySigningKeyResponse> deleteKeySigningKey(
      $grpc.ServiceCall call, $0.DeleteKeySigningKeyRequest request);

  $async.Future<$0.DeleteQueryLoggingConfigResponse>
      deleteQueryLoggingConfig_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteQueryLoggingConfigRequest> $request) async {
    return deleteQueryLoggingConfig($call, await $request);
  }

  $async.Future<$0.DeleteQueryLoggingConfigResponse> deleteQueryLoggingConfig(
      $grpc.ServiceCall call, $0.DeleteQueryLoggingConfigRequest request);

  $async.Future<$0.DeleteReusableDelegationSetResponse>
      deleteReusableDelegationSet_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteReusableDelegationSetRequest> $request) async {
    return deleteReusableDelegationSet($call, await $request);
  }

  $async.Future<$0.DeleteReusableDelegationSetResponse>
      deleteReusableDelegationSet($grpc.ServiceCall call,
          $0.DeleteReusableDelegationSetRequest request);

  $async.Future<$0.DeleteTrafficPolicyResponse> deleteTrafficPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteTrafficPolicyRequest> $request) async {
    return deleteTrafficPolicy($call, await $request);
  }

  $async.Future<$0.DeleteTrafficPolicyResponse> deleteTrafficPolicy(
      $grpc.ServiceCall call, $0.DeleteTrafficPolicyRequest request);

  $async.Future<$0.DeleteTrafficPolicyInstanceResponse>
      deleteTrafficPolicyInstance_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteTrafficPolicyInstanceRequest> $request) async {
    return deleteTrafficPolicyInstance($call, await $request);
  }

  $async.Future<$0.DeleteTrafficPolicyInstanceResponse>
      deleteTrafficPolicyInstance($grpc.ServiceCall call,
          $0.DeleteTrafficPolicyInstanceRequest request);

  $async.Future<$0.DeleteVPCAssociationAuthorizationResponse>
      deleteVPCAssociationAuthorization_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeleteVPCAssociationAuthorizationRequest>
              $request) async {
    return deleteVPCAssociationAuthorization($call, await $request);
  }

  $async.Future<$0.DeleteVPCAssociationAuthorizationResponse>
      deleteVPCAssociationAuthorization($grpc.ServiceCall call,
          $0.DeleteVPCAssociationAuthorizationRequest request);

  $async.Future<$0.DisableHostedZoneDNSSECResponse> disableHostedZoneDNSSEC_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DisableHostedZoneDNSSECRequest> $request) async {
    return disableHostedZoneDNSSEC($call, await $request);
  }

  $async.Future<$0.DisableHostedZoneDNSSECResponse> disableHostedZoneDNSSEC(
      $grpc.ServiceCall call, $0.DisableHostedZoneDNSSECRequest request);

  $async.Future<$0.DisassociateVPCFromHostedZoneResponse>
      disassociateVPCFromHostedZone_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DisassociateVPCFromHostedZoneRequest>
              $request) async {
    return disassociateVPCFromHostedZone($call, await $request);
  }

  $async.Future<$0.DisassociateVPCFromHostedZoneResponse>
      disassociateVPCFromHostedZone($grpc.ServiceCall call,
          $0.DisassociateVPCFromHostedZoneRequest request);

  $async.Future<$0.EnableHostedZoneDNSSECResponse> enableHostedZoneDNSSEC_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.EnableHostedZoneDNSSECRequest> $request) async {
    return enableHostedZoneDNSSEC($call, await $request);
  }

  $async.Future<$0.EnableHostedZoneDNSSECResponse> enableHostedZoneDNSSEC(
      $grpc.ServiceCall call, $0.EnableHostedZoneDNSSECRequest request);

  $async.Future<$0.GetAccountLimitResponse> getAccountLimit_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAccountLimitRequest> $request) async {
    return getAccountLimit($call, await $request);
  }

  $async.Future<$0.GetAccountLimitResponse> getAccountLimit(
      $grpc.ServiceCall call, $0.GetAccountLimitRequest request);

  $async.Future<$0.GetChangeResponse> getChange_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetChangeRequest> $request) async {
    return getChange($call, await $request);
  }

  $async.Future<$0.GetChangeResponse> getChange(
      $grpc.ServiceCall call, $0.GetChangeRequest request);

  $async.Future<$0.GetCheckerIpRangesResponse> getCheckerIpRanges_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCheckerIpRangesRequest> $request) async {
    return getCheckerIpRanges($call, await $request);
  }

  $async.Future<$0.GetCheckerIpRangesResponse> getCheckerIpRanges(
      $grpc.ServiceCall call, $0.GetCheckerIpRangesRequest request);

  $async.Future<$0.GetDNSSECResponse> getDNSSEC_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDNSSECRequest> $request) async {
    return getDNSSEC($call, await $request);
  }

  $async.Future<$0.GetDNSSECResponse> getDNSSEC(
      $grpc.ServiceCall call, $0.GetDNSSECRequest request);

  $async.Future<$0.GetGeoLocationResponse> getGeoLocation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetGeoLocationRequest> $request) async {
    return getGeoLocation($call, await $request);
  }

  $async.Future<$0.GetGeoLocationResponse> getGeoLocation(
      $grpc.ServiceCall call, $0.GetGeoLocationRequest request);

  $async.Future<$0.GetHealthCheckResponse> getHealthCheck_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHealthCheckRequest> $request) async {
    return getHealthCheck($call, await $request);
  }

  $async.Future<$0.GetHealthCheckResponse> getHealthCheck(
      $grpc.ServiceCall call, $0.GetHealthCheckRequest request);

  $async.Future<$0.GetHealthCheckCountResponse> getHealthCheckCount_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHealthCheckCountRequest> $request) async {
    return getHealthCheckCount($call, await $request);
  }

  $async.Future<$0.GetHealthCheckCountResponse> getHealthCheckCount(
      $grpc.ServiceCall call, $0.GetHealthCheckCountRequest request);

  $async.Future<$0.GetHealthCheckLastFailureReasonResponse>
      getHealthCheckLastFailureReason_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetHealthCheckLastFailureReasonRequest>
              $request) async {
    return getHealthCheckLastFailureReason($call, await $request);
  }

  $async.Future<$0.GetHealthCheckLastFailureReasonResponse>
      getHealthCheckLastFailureReason($grpc.ServiceCall call,
          $0.GetHealthCheckLastFailureReasonRequest request);

  $async.Future<$0.GetHealthCheckStatusResponse> getHealthCheckStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHealthCheckStatusRequest> $request) async {
    return getHealthCheckStatus($call, await $request);
  }

  $async.Future<$0.GetHealthCheckStatusResponse> getHealthCheckStatus(
      $grpc.ServiceCall call, $0.GetHealthCheckStatusRequest request);

  $async.Future<$0.GetHostedZoneResponse> getHostedZone_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHostedZoneRequest> $request) async {
    return getHostedZone($call, await $request);
  }

  $async.Future<$0.GetHostedZoneResponse> getHostedZone(
      $grpc.ServiceCall call, $0.GetHostedZoneRequest request);

  $async.Future<$0.GetHostedZoneCountResponse> getHostedZoneCount_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHostedZoneCountRequest> $request) async {
    return getHostedZoneCount($call, await $request);
  }

  $async.Future<$0.GetHostedZoneCountResponse> getHostedZoneCount(
      $grpc.ServiceCall call, $0.GetHostedZoneCountRequest request);

  $async.Future<$0.GetHostedZoneLimitResponse> getHostedZoneLimit_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHostedZoneLimitRequest> $request) async {
    return getHostedZoneLimit($call, await $request);
  }

  $async.Future<$0.GetHostedZoneLimitResponse> getHostedZoneLimit(
      $grpc.ServiceCall call, $0.GetHostedZoneLimitRequest request);

  $async.Future<$0.GetQueryLoggingConfigResponse> getQueryLoggingConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetQueryLoggingConfigRequest> $request) async {
    return getQueryLoggingConfig($call, await $request);
  }

  $async.Future<$0.GetQueryLoggingConfigResponse> getQueryLoggingConfig(
      $grpc.ServiceCall call, $0.GetQueryLoggingConfigRequest request);

  $async.Future<$0.GetReusableDelegationSetResponse>
      getReusableDelegationSet_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetReusableDelegationSetRequest> $request) async {
    return getReusableDelegationSet($call, await $request);
  }

  $async.Future<$0.GetReusableDelegationSetResponse> getReusableDelegationSet(
      $grpc.ServiceCall call, $0.GetReusableDelegationSetRequest request);

  $async.Future<$0.GetReusableDelegationSetLimitResponse>
      getReusableDelegationSetLimit_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetReusableDelegationSetLimitRequest>
              $request) async {
    return getReusableDelegationSetLimit($call, await $request);
  }

  $async.Future<$0.GetReusableDelegationSetLimitResponse>
      getReusableDelegationSetLimit($grpc.ServiceCall call,
          $0.GetReusableDelegationSetLimitRequest request);

  $async.Future<$0.GetTrafficPolicyResponse> getTrafficPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTrafficPolicyRequest> $request) async {
    return getTrafficPolicy($call, await $request);
  }

  $async.Future<$0.GetTrafficPolicyResponse> getTrafficPolicy(
      $grpc.ServiceCall call, $0.GetTrafficPolicyRequest request);

  $async.Future<$0.GetTrafficPolicyInstanceResponse>
      getTrafficPolicyInstance_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetTrafficPolicyInstanceRequest> $request) async {
    return getTrafficPolicyInstance($call, await $request);
  }

  $async.Future<$0.GetTrafficPolicyInstanceResponse> getTrafficPolicyInstance(
      $grpc.ServiceCall call, $0.GetTrafficPolicyInstanceRequest request);

  $async.Future<$0.GetTrafficPolicyInstanceCountResponse>
      getTrafficPolicyInstanceCount_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetTrafficPolicyInstanceCountRequest>
              $request) async {
    return getTrafficPolicyInstanceCount($call, await $request);
  }

  $async.Future<$0.GetTrafficPolicyInstanceCountResponse>
      getTrafficPolicyInstanceCount($grpc.ServiceCall call,
          $0.GetTrafficPolicyInstanceCountRequest request);

  $async.Future<$0.ListCidrBlocksResponse> listCidrBlocks_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCidrBlocksRequest> $request) async {
    return listCidrBlocks($call, await $request);
  }

  $async.Future<$0.ListCidrBlocksResponse> listCidrBlocks(
      $grpc.ServiceCall call, $0.ListCidrBlocksRequest request);

  $async.Future<$0.ListCidrCollectionsResponse> listCidrCollections_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCidrCollectionsRequest> $request) async {
    return listCidrCollections($call, await $request);
  }

  $async.Future<$0.ListCidrCollectionsResponse> listCidrCollections(
      $grpc.ServiceCall call, $0.ListCidrCollectionsRequest request);

  $async.Future<$0.ListCidrLocationsResponse> listCidrLocations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCidrLocationsRequest> $request) async {
    return listCidrLocations($call, await $request);
  }

  $async.Future<$0.ListCidrLocationsResponse> listCidrLocations(
      $grpc.ServiceCall call, $0.ListCidrLocationsRequest request);

  $async.Future<$0.ListGeoLocationsResponse> listGeoLocations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListGeoLocationsRequest> $request) async {
    return listGeoLocations($call, await $request);
  }

  $async.Future<$0.ListGeoLocationsResponse> listGeoLocations(
      $grpc.ServiceCall call, $0.ListGeoLocationsRequest request);

  $async.Future<$0.ListHealthChecksResponse> listHealthChecks_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListHealthChecksRequest> $request) async {
    return listHealthChecks($call, await $request);
  }

  $async.Future<$0.ListHealthChecksResponse> listHealthChecks(
      $grpc.ServiceCall call, $0.ListHealthChecksRequest request);

  $async.Future<$0.ListHostedZonesResponse> listHostedZones_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListHostedZonesRequest> $request) async {
    return listHostedZones($call, await $request);
  }

  $async.Future<$0.ListHostedZonesResponse> listHostedZones(
      $grpc.ServiceCall call, $0.ListHostedZonesRequest request);

  $async.Future<$0.ListHostedZonesByNameResponse> listHostedZonesByName_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListHostedZonesByNameRequest> $request) async {
    return listHostedZonesByName($call, await $request);
  }

  $async.Future<$0.ListHostedZonesByNameResponse> listHostedZonesByName(
      $grpc.ServiceCall call, $0.ListHostedZonesByNameRequest request);

  $async.Future<$0.ListHostedZonesByVPCResponse> listHostedZonesByVPC_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListHostedZonesByVPCRequest> $request) async {
    return listHostedZonesByVPC($call, await $request);
  }

  $async.Future<$0.ListHostedZonesByVPCResponse> listHostedZonesByVPC(
      $grpc.ServiceCall call, $0.ListHostedZonesByVPCRequest request);

  $async.Future<$0.ListQueryLoggingConfigsResponse> listQueryLoggingConfigs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListQueryLoggingConfigsRequest> $request) async {
    return listQueryLoggingConfigs($call, await $request);
  }

  $async.Future<$0.ListQueryLoggingConfigsResponse> listQueryLoggingConfigs(
      $grpc.ServiceCall call, $0.ListQueryLoggingConfigsRequest request);

  $async.Future<$0.ListResourceRecordSetsResponse> listResourceRecordSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListResourceRecordSetsRequest> $request) async {
    return listResourceRecordSets($call, await $request);
  }

  $async.Future<$0.ListResourceRecordSetsResponse> listResourceRecordSets(
      $grpc.ServiceCall call, $0.ListResourceRecordSetsRequest request);

  $async.Future<$0.ListReusableDelegationSetsResponse>
      listReusableDelegationSets_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListReusableDelegationSetsRequest> $request) async {
    return listReusableDelegationSets($call, await $request);
  }

  $async.Future<$0.ListReusableDelegationSetsResponse>
      listReusableDelegationSets(
          $grpc.ServiceCall call, $0.ListReusableDelegationSetsRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListTagsForResourcesResponse> listTagsForResources_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourcesRequest> $request) async {
    return listTagsForResources($call, await $request);
  }

  $async.Future<$0.ListTagsForResourcesResponse> listTagsForResources(
      $grpc.ServiceCall call, $0.ListTagsForResourcesRequest request);

  $async.Future<$0.ListTrafficPoliciesResponse> listTrafficPolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTrafficPoliciesRequest> $request) async {
    return listTrafficPolicies($call, await $request);
  }

  $async.Future<$0.ListTrafficPoliciesResponse> listTrafficPolicies(
      $grpc.ServiceCall call, $0.ListTrafficPoliciesRequest request);

  $async.Future<$0.ListTrafficPolicyInstancesResponse>
      listTrafficPolicyInstances_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListTrafficPolicyInstancesRequest> $request) async {
    return listTrafficPolicyInstances($call, await $request);
  }

  $async.Future<$0.ListTrafficPolicyInstancesResponse>
      listTrafficPolicyInstances(
          $grpc.ServiceCall call, $0.ListTrafficPolicyInstancesRequest request);

  $async.Future<$0.ListTrafficPolicyInstancesByHostedZoneResponse>
      listTrafficPolicyInstancesByHostedZone_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListTrafficPolicyInstancesByHostedZoneRequest>
              $request) async {
    return listTrafficPolicyInstancesByHostedZone($call, await $request);
  }

  $async.Future<$0.ListTrafficPolicyInstancesByHostedZoneResponse>
      listTrafficPolicyInstancesByHostedZone($grpc.ServiceCall call,
          $0.ListTrafficPolicyInstancesByHostedZoneRequest request);

  $async.Future<$0.ListTrafficPolicyInstancesByPolicyResponse>
      listTrafficPolicyInstancesByPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListTrafficPolicyInstancesByPolicyRequest>
              $request) async {
    return listTrafficPolicyInstancesByPolicy($call, await $request);
  }

  $async.Future<$0.ListTrafficPolicyInstancesByPolicyResponse>
      listTrafficPolicyInstancesByPolicy($grpc.ServiceCall call,
          $0.ListTrafficPolicyInstancesByPolicyRequest request);

  $async.Future<$0.ListTrafficPolicyVersionsResponse>
      listTrafficPolicyVersions_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListTrafficPolicyVersionsRequest> $request) async {
    return listTrafficPolicyVersions($call, await $request);
  }

  $async.Future<$0.ListTrafficPolicyVersionsResponse> listTrafficPolicyVersions(
      $grpc.ServiceCall call, $0.ListTrafficPolicyVersionsRequest request);

  $async.Future<$0.ListVPCAssociationAuthorizationsResponse>
      listVPCAssociationAuthorizations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListVPCAssociationAuthorizationsRequest>
              $request) async {
    return listVPCAssociationAuthorizations($call, await $request);
  }

  $async.Future<$0.ListVPCAssociationAuthorizationsResponse>
      listVPCAssociationAuthorizations($grpc.ServiceCall call,
          $0.ListVPCAssociationAuthorizationsRequest request);

  $async.Future<$0.TestDNSAnswerResponse> testDNSAnswer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestDNSAnswerRequest> $request) async {
    return testDNSAnswer($call, await $request);
  }

  $async.Future<$0.TestDNSAnswerResponse> testDNSAnswer(
      $grpc.ServiceCall call, $0.TestDNSAnswerRequest request);

  $async.Future<$0.UpdateHealthCheckResponse> updateHealthCheck_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateHealthCheckRequest> $request) async {
    return updateHealthCheck($call, await $request);
  }

  $async.Future<$0.UpdateHealthCheckResponse> updateHealthCheck(
      $grpc.ServiceCall call, $0.UpdateHealthCheckRequest request);

  $async.Future<$0.UpdateHostedZoneCommentResponse> updateHostedZoneComment_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateHostedZoneCommentRequest> $request) async {
    return updateHostedZoneComment($call, await $request);
  }

  $async.Future<$0.UpdateHostedZoneCommentResponse> updateHostedZoneComment(
      $grpc.ServiceCall call, $0.UpdateHostedZoneCommentRequest request);

  $async.Future<$0.UpdateHostedZoneFeaturesResponse>
      updateHostedZoneFeatures_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateHostedZoneFeaturesRequest> $request) async {
    return updateHostedZoneFeatures($call, await $request);
  }

  $async.Future<$0.UpdateHostedZoneFeaturesResponse> updateHostedZoneFeatures(
      $grpc.ServiceCall call, $0.UpdateHostedZoneFeaturesRequest request);

  $async.Future<$0.UpdateTrafficPolicyCommentResponse>
      updateTrafficPolicyComment_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateTrafficPolicyCommentRequest> $request) async {
    return updateTrafficPolicyComment($call, await $request);
  }

  $async.Future<$0.UpdateTrafficPolicyCommentResponse>
      updateTrafficPolicyComment(
          $grpc.ServiceCall call, $0.UpdateTrafficPolicyCommentRequest request);

  $async.Future<$0.UpdateTrafficPolicyInstanceResponse>
      updateTrafficPolicyInstance_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateTrafficPolicyInstanceRequest> $request) async {
    return updateTrafficPolicyInstance($call, await $request);
  }

  $async.Future<$0.UpdateTrafficPolicyInstanceResponse>
      updateTrafficPolicyInstance($grpc.ServiceCall call,
          $0.UpdateTrafficPolicyInstanceRequest request);
}
