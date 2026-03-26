// This is a generated file - do not edit.
//
// Generated from email.proto.

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

import 'email.pb.dart' as $0;

export 'email.pb.dart';

/// SESv2Service provides email API operations.
@$pb.GrpcServiceName('email.SESv2Service')
class SESv2ServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SESv2ServiceClient(super.channel, {super.options, super.interceptors});

  /// Retrieves batches of metric data collected based on your sending activity. You can execute this operation no more than 16 times per second, and with at most 160 queries from the batches per second ...
  /// HTTP: POST /v2/email/metrics/batch
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.BatchGetMetricDataResponse> batchGetMetricData(
    $0.BatchGetMetricDataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchGetMetricData, request, options: options);
  }

  /// Cancels an export job.
  /// HTTP: PUT /v2/email/export-jobs/{JobId}/cancel
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CancelExportJobResponse> cancelExportJob(
    $0.CancelExportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelExportJob, request, options: options);
  }

  /// Create a configuration set. Configuration sets are groups of rules that you can apply to the emails that you send. You apply a configuration set to an email by specifying the name of the configurat...
  /// HTTP: POST /v2/email/configuration-sets
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateConfigurationSetResponse>
      createConfigurationSet(
    $0.CreateConfigurationSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createConfigurationSet, request,
        options: options);
  }

  /// Create an event destination. Events include message sends, deliveries, opens, clicks, bounces, and complaints. Event destinations are places that you can send information about these events to. For...
  /// HTTP: POST /v2/email/configuration-sets/{ConfigurationSetName}/event-destinations
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateConfigurationSetEventDestinationResponse>
      createConfigurationSetEventDestination(
    $0.CreateConfigurationSetEventDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createConfigurationSetEventDestination, request,
        options: options);
  }

  /// Creates a contact, which is an end-user who is receiving the email, and adds them to a contact list.
  /// HTTP: POST /v2/email/contact-lists/{ContactListName}/contacts
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateContactResponse> createContact(
    $0.CreateContactRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createContact, request, options: options);
  }

  /// Creates a contact list.
  /// HTTP: POST /v2/email/contact-lists
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateContactListResponse> createContactList(
    $0.CreateContactListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createContactList, request, options: options);
  }

  /// Creates a new custom verification email template. For more information about custom verification email templates, see Using custom verification email templates in the Amazon SES Developer Guide. Yo...
  /// HTTP: POST /v2/email/custom-verification-email-templates
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateCustomVerificationEmailTemplateResponse>
      createCustomVerificationEmailTemplate(
    $0.CreateCustomVerificationEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCustomVerificationEmailTemplate, request,
        options: options);
  }

  /// Create a new pool of dedicated IP addresses. A pool can include one or more dedicated IP addresses that are associated with your Amazon Web Services account. You can associate a pool with a configu...
  /// HTTP: POST /v2/email/dedicated-ip-pools
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateDedicatedIpPoolResponse> createDedicatedIpPool(
    $0.CreateDedicatedIpPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDedicatedIpPool, request, options: options);
  }

  /// Create a new predictive inbox placement test. Predictive inbox placement tests can help you predict how your messages will be handled by various email providers around the world. When you perform a...
  /// HTTP: POST /v2/email/deliverability-dashboard/test
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateDeliverabilityTestReportResponse>
      createDeliverabilityTestReport(
    $0.CreateDeliverabilityTestReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDeliverabilityTestReport, request,
        options: options);
  }

  /// Starts the process of verifying an email identity. An identity is an email address or domain that you use when you send email. Before you can use an identity to send email, you first have to verify...
  /// HTTP: POST /v2/email/identities
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateEmailIdentityResponse> createEmailIdentity(
    $0.CreateEmailIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createEmailIdentity, request, options: options);
  }

  /// Creates the specified sending authorization policy for the given identity (an email address or a domain). This API is for the identity owner only. If you have not verified the identity, this API wi...
  /// HTTP: POST /v2/email/identities/{EmailIdentity}/policies/{PolicyName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateEmailIdentityPolicyResponse>
      createEmailIdentityPolicy(
    $0.CreateEmailIdentityPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createEmailIdentityPolicy, request,
        options: options);
  }

  /// Creates an email template. Email templates enable you to send personalized email to one or more destinations in a single API operation. For more information, see the Amazon SES Developer Guide. You...
  /// HTTP: POST /v2/email/templates
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateEmailTemplateResponse> createEmailTemplate(
    $0.CreateEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createEmailTemplate, request, options: options);
  }

  /// Creates an export job for a data source and destination. You can execute this operation no more than once per second.
  /// HTTP: POST /v2/email/export-jobs
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateExportJobResponse> createExportJob(
    $0.CreateExportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createExportJob, request, options: options);
  }

  /// Creates an import job for a data destination.
  /// HTTP: POST /v2/email/import-jobs
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateImportJobResponse> createImportJob(
    $0.CreateImportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createImportJob, request, options: options);
  }

  /// Creates a multi-region endpoint (global-endpoint). The primary region is going to be the AWS-Region where the operation is executed. The secondary region has to be provided in request's parameters....
  /// HTTP: POST /v2/email/multi-region-endpoints
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateMultiRegionEndpointResponse>
      createMultiRegionEndpoint(
    $0.CreateMultiRegionEndpointRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createMultiRegionEndpoint, request,
        options: options);
  }

  /// Create a tenant. Tenants are logical containers that group related SES resources together. Each tenant can have its own set of resources like email identities, configuration sets, and templates, al...
  /// HTTP: POST /v2/email/tenants
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateTenantResponse> createTenant(
    $0.CreateTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTenant, request, options: options);
  }

  /// Associate a resource with a tenant. Resources can be email identities, configuration sets, or email templates. When you associate a resource with a tenant, you can use that resource when sending em...
  /// HTTP: POST /v2/email/tenants/resources
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateTenantResourceAssociationResponse>
      createTenantResourceAssociation(
    $0.CreateTenantResourceAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTenantResourceAssociation, request,
        options: options);
  }

  /// Delete an existing configuration set. Configuration sets are groups of rules that you can apply to the emails you send. You apply a configuration set to an email by including a reference to the con...
  /// HTTP: DELETE /v2/email/configuration-sets/{ConfigurationSetName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteConfigurationSetResponse>
      deleteConfigurationSet(
    $0.DeleteConfigurationSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteConfigurationSet, request,
        options: options);
  }

  /// Delete an event destination. Events include message sends, deliveries, opens, clicks, bounces, and complaints. Event destinations are places that you can send information about these events to. For...
  /// HTTP: DELETE /v2/email/configuration-sets/{ConfigurationSetName}/event-destinations/{EventDestinationName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteConfigurationSetEventDestinationResponse>
      deleteConfigurationSetEventDestination(
    $0.DeleteConfigurationSetEventDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteConfigurationSetEventDestination, request,
        options: options);
  }

  /// Removes a contact from a contact list.
  /// HTTP: DELETE /v2/email/contact-lists/{ContactListName}/contacts/{EmailAddress}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteContactResponse> deleteContact(
    $0.DeleteContactRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteContact, request, options: options);
  }

  /// Deletes a contact list and all of the contacts on that list.
  /// HTTP: DELETE /v2/email/contact-lists/{ContactListName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteContactListResponse> deleteContactList(
    $0.DeleteContactListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteContactList, request, options: options);
  }

  /// Deletes an existing custom verification email template. For more information about custom verification email templates, see Using custom verification email templates in the Amazon SES Developer Gui...
  /// HTTP: DELETE /v2/email/custom-verification-email-templates/{TemplateName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteCustomVerificationEmailTemplateResponse>
      deleteCustomVerificationEmailTemplate(
    $0.DeleteCustomVerificationEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCustomVerificationEmailTemplate, request,
        options: options);
  }

  /// Delete a dedicated IP pool.
  /// HTTP: DELETE /v2/email/dedicated-ip-pools/{PoolName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteDedicatedIpPoolResponse> deleteDedicatedIpPool(
    $0.DeleteDedicatedIpPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDedicatedIpPool, request, options: options);
  }

  /// Deletes an email identity. An identity can be either an email address or a domain name.
  /// HTTP: DELETE /v2/email/identities/{EmailIdentity}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteEmailIdentityResponse> deleteEmailIdentity(
    $0.DeleteEmailIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteEmailIdentity, request, options: options);
  }

  /// Deletes the specified sending authorization policy for the given identity (an email address or a domain). This API returns successfully even if a policy with the specified name does not exist. This...
  /// HTTP: DELETE /v2/email/identities/{EmailIdentity}/policies/{PolicyName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteEmailIdentityPolicyResponse>
      deleteEmailIdentityPolicy(
    $0.DeleteEmailIdentityPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteEmailIdentityPolicy, request,
        options: options);
  }

  /// Deletes an email template. You can execute this operation no more than once per second.
  /// HTTP: DELETE /v2/email/templates/{TemplateName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteEmailTemplateResponse> deleteEmailTemplate(
    $0.DeleteEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteEmailTemplate, request, options: options);
  }

  /// Deletes a multi-region endpoint (global-endpoint). Only multi-region endpoints (global-endpoints) whose primary region is the AWS-Region where operation is executed can be deleted.
  /// HTTP: DELETE /v2/email/multi-region-endpoints/{EndpointName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteMultiRegionEndpointResponse>
      deleteMultiRegionEndpoint(
    $0.DeleteMultiRegionEndpointRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMultiRegionEndpoint, request,
        options: options);
  }

  /// Removes an email address from the suppression list for your account.
  /// HTTP: DELETE /v2/email/suppression/addresses/{EmailAddress}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteSuppressedDestinationResponse>
      deleteSuppressedDestination(
    $0.DeleteSuppressedDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSuppressedDestination, request,
        options: options);
  }

  /// Delete an existing tenant. When you delete a tenant, its associations with resources are removed, but the resources themselves are not deleted.
  /// HTTP: POST /v2/email/tenants/delete
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteTenantResponse> deleteTenant(
    $0.DeleteTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTenant, request, options: options);
  }

  /// Delete an association between a tenant and a resource. When you delete a tenant-resource association, the resource itself is not deleted, only its association with the specific tenant is removed. A...
  /// HTTP: POST /v2/email/tenants/resources/delete
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteTenantResourceAssociationResponse>
      deleteTenantResourceAssociation(
    $0.DeleteTenantResourceAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTenantResourceAssociation, request,
        options: options);
  }

  /// Obtain information about the email-sending status and capabilities of your Amazon SES account in the current Amazon Web Services Region.
  /// HTTP: GET /v2/email/account
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetAccountResponse> getAccount(
    $0.GetAccountRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccount, request, options: options);
  }

  /// Retrieve a list of the blacklists that your dedicated IP addresses appear on.
  /// HTTP: GET /v2/email/deliverability-dashboard/blacklist-report
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetBlacklistReportsResponse> getBlacklistReports(
    $0.GetBlacklistReportsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getBlacklistReports, request, options: options);
  }

  /// Get information about an existing configuration set, including the dedicated IP pool that it's associated with, whether or not it's enabled for sending email, and more. Configuration sets are group...
  /// HTTP: GET /v2/email/configuration-sets/{ConfigurationSetName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetConfigurationSetResponse> getConfigurationSet(
    $0.GetConfigurationSetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConfigurationSet, request, options: options);
  }

  /// Retrieve a list of event destinations that are associated with a configuration set. Events include message sends, deliveries, opens, clicks, bounces, and complaints. Event destinations are places t...
  /// HTTP: GET /v2/email/configuration-sets/{ConfigurationSetName}/event-destinations
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetConfigurationSetEventDestinationsResponse>
      getConfigurationSetEventDestinations(
    $0.GetConfigurationSetEventDestinationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConfigurationSetEventDestinations, request,
        options: options);
  }

  /// Returns a contact from a contact list.
  /// HTTP: GET /v2/email/contact-lists/{ContactListName}/contacts/{EmailAddress}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetContactResponse> getContact(
    $0.GetContactRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContact, request, options: options);
  }

  /// Returns contact list metadata. It does not return any information about the contacts present in the list.
  /// HTTP: GET /v2/email/contact-lists/{ContactListName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetContactListResponse> getContactList(
    $0.GetContactListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContactList, request, options: options);
  }

  /// Returns the custom email verification template for the template name you specify. For more information about custom verification email templates, see Using custom verification email templates in th...
  /// HTTP: GET /v2/email/custom-verification-email-templates/{TemplateName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetCustomVerificationEmailTemplateResponse>
      getCustomVerificationEmailTemplate(
    $0.GetCustomVerificationEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCustomVerificationEmailTemplate, request,
        options: options);
  }

  /// Get information about a dedicated IP address, including the name of the dedicated IP pool that it's associated with, as well information about the automatic warm-up process for the address.
  /// HTTP: GET /v2/email/dedicated-ips/{Ip}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDedicatedIpResponse> getDedicatedIp(
    $0.GetDedicatedIpRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDedicatedIp, request, options: options);
  }

  /// Retrieve information about the dedicated pool.
  /// HTTP: GET /v2/email/dedicated-ip-pools/{PoolName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDedicatedIpPoolResponse> getDedicatedIpPool(
    $0.GetDedicatedIpPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDedicatedIpPool, request, options: options);
  }

  /// List the dedicated IP addresses that are associated with your Amazon Web Services account.
  /// HTTP: GET /v2/email/dedicated-ips
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDedicatedIpsResponse> getDedicatedIps(
    $0.GetDedicatedIpsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDedicatedIps, request, options: options);
  }

  /// Retrieve information about the status of the Deliverability dashboard for your account. When the Deliverability dashboard is enabled, you gain access to reputation, deliverability, and other metric...
  /// HTTP: GET /v2/email/deliverability-dashboard
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDeliverabilityDashboardOptionsResponse>
      getDeliverabilityDashboardOptions(
    $0.GetDeliverabilityDashboardOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeliverabilityDashboardOptions, request,
        options: options);
  }

  /// Retrieve the results of a predictive inbox placement test.
  /// HTTP: GET /v2/email/deliverability-dashboard/test-reports/{ReportId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDeliverabilityTestReportResponse>
      getDeliverabilityTestReport(
    $0.GetDeliverabilityTestReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeliverabilityTestReport, request,
        options: options);
  }

  /// Retrieve all the deliverability data for a specific campaign. This data is available for a campaign only if the campaign sent email by using a domain that the Deliverability dashboard is enabled for.
  /// HTTP: GET /v2/email/deliverability-dashboard/campaigns/{CampaignId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDomainDeliverabilityCampaignResponse>
      getDomainDeliverabilityCampaign(
    $0.GetDomainDeliverabilityCampaignRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDomainDeliverabilityCampaign, request,
        options: options);
  }

  /// Retrieve inbox placement and engagement rates for the domains that you use to send email.
  /// HTTP: GET /v2/email/deliverability-dashboard/statistics-report/{Domain}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDomainStatisticsReportResponse>
      getDomainStatisticsReport(
    $0.GetDomainStatisticsReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDomainStatisticsReport, request,
        options: options);
  }

  /// Provides validation insights about a specific email address, including syntax validation, DNS record checks, mailbox existence, and other deliverability factors.
  /// HTTP: POST /v2/email/email-address-insights
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetEmailAddressInsightsResponse>
      getEmailAddressInsights(
    $0.GetEmailAddressInsightsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEmailAddressInsights, request,
        options: options);
  }

  /// Provides information about a specific identity, including the identity's verification status, sending authorization policies, its DKIM authentication status, and its custom Mail-From settings.
  /// HTTP: GET /v2/email/identities/{EmailIdentity}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetEmailIdentityResponse> getEmailIdentity(
    $0.GetEmailIdentityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEmailIdentity, request, options: options);
  }

  /// Returns the requested sending authorization policies for the given identity (an email address or a domain). The policies are returned as a map of policy names to policy contents. You can retrieve a...
  /// HTTP: GET /v2/email/identities/{EmailIdentity}/policies
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetEmailIdentityPoliciesResponse>
      getEmailIdentityPolicies(
    $0.GetEmailIdentityPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEmailIdentityPolicies, request,
        options: options);
  }

  /// Displays the template object (which includes the subject line, HTML part and text part) for the template you specify. You can execute this operation no more than once per second.
  /// HTTP: GET /v2/email/templates/{TemplateName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetEmailTemplateResponse> getEmailTemplate(
    $0.GetEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEmailTemplate, request, options: options);
  }

  /// Provides information about an export job.
  /// HTTP: GET /v2/email/export-jobs/{JobId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetExportJobResponse> getExportJob(
    $0.GetExportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getExportJob, request, options: options);
  }

  /// Provides information about an import job.
  /// HTTP: GET /v2/email/import-jobs/{JobId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetImportJobResponse> getImportJob(
    $0.GetImportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getImportJob, request, options: options);
  }

  /// Provides information about a specific message, including the from address, the subject, the recipient address, email tags, as well as events associated with the message. You can execute this operat...
  /// HTTP: GET /v2/email/insights/{MessageId}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetMessageInsightsResponse> getMessageInsights(
    $0.GetMessageInsightsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMessageInsights, request, options: options);
  }

  /// Displays the multi-region endpoint (global-endpoint) configuration. Only multi-region endpoints (global-endpoints) whose primary region is the AWS-Region where operation is executed can be displayed.
  /// HTTP: GET /v2/email/multi-region-endpoints/{EndpointName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetMultiRegionEndpointResponse>
      getMultiRegionEndpoint(
    $0.GetMultiRegionEndpointRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMultiRegionEndpoint, request,
        options: options);
  }

  /// Retrieve information about a specific reputation entity, including its reputation management policy, customer-managed status, Amazon Web Services Amazon SES-managed status, and aggregate sending st...
  /// HTTP: GET /v2/email/reputation/entities/{ReputationEntityType}/{ReputationEntityReference}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetReputationEntityResponse> getReputationEntity(
    $0.GetReputationEntityRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getReputationEntity, request, options: options);
  }

  /// Retrieves information about a specific email address that's on the suppression list for your account.
  /// HTTP: GET /v2/email/suppression/addresses/{EmailAddress}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetSuppressedDestinationResponse>
      getSuppressedDestination(
    $0.GetSuppressedDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSuppressedDestination, request,
        options: options);
  }

  /// Get information about a specific tenant, including the tenant's name, ID, ARN, creation timestamp, tags, and sending status.
  /// HTTP: POST /v2/email/tenants/get
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetTenantResponse> getTenant(
    $0.GetTenantRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTenant, request, options: options);
  }

  /// List all of the configuration sets associated with your account in the current region. Configuration sets are groups of rules that you can apply to the emails you send. You apply a configuration se...
  /// HTTP: GET /v2/email/configuration-sets
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListConfigurationSetsResponse> listConfigurationSets(
    $0.ListConfigurationSetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listConfigurationSets, request, options: options);
  }

  /// Lists all of the contact lists available. If your output includes a "NextToken" field with a string value, this indicates there may be additional contacts on the filtered list - regardless of the n...
  /// HTTP: GET /v2/email/contact-lists
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListContactListsResponse> listContactLists(
    $0.ListContactListsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listContactLists, request, options: options);
  }

  /// Lists the contacts present in a specific contact list.
  /// HTTP: POST /v2/email/contact-lists/{ContactListName}/contacts/list
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListContactsResponse> listContacts(
    $0.ListContactsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listContacts, request, options: options);
  }

  /// Lists the existing custom verification email templates for your account in the current Amazon Web Services Region. For more information about custom verification email templates, see Using custom v...
  /// HTTP: GET /v2/email/custom-verification-email-templates
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListCustomVerificationEmailTemplatesResponse>
      listCustomVerificationEmailTemplates(
    $0.ListCustomVerificationEmailTemplatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCustomVerificationEmailTemplates, request,
        options: options);
  }

  /// List all of the dedicated IP pools that exist in your Amazon Web Services account in the current Region.
  /// HTTP: GET /v2/email/dedicated-ip-pools
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListDedicatedIpPoolsResponse> listDedicatedIpPools(
    $0.ListDedicatedIpPoolsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDedicatedIpPools, request, options: options);
  }

  /// Show a list of the predictive inbox placement tests that you've performed, regardless of their statuses. For predictive inbox placement tests that are complete, you can use the GetDeliverabilityTes...
  /// HTTP: GET /v2/email/deliverability-dashboard/test-reports
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListDeliverabilityTestReportsResponse>
      listDeliverabilityTestReports(
    $0.ListDeliverabilityTestReportsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDeliverabilityTestReports, request,
        options: options);
  }

  /// Retrieve deliverability data for all the campaigns that used a specific domain to send email during a specified time range. This data is available for a domain only if you enabled the Deliverabilit...
  /// HTTP: GET /v2/email/deliverability-dashboard/domains/{SubscribedDomain}/campaigns
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListDomainDeliverabilityCampaignsResponse>
      listDomainDeliverabilityCampaigns(
    $0.ListDomainDeliverabilityCampaignsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDomainDeliverabilityCampaigns, request,
        options: options);
  }

  /// Returns a list of all of the email identities that are associated with your Amazon Web Services account. An identity can be either an email address or a domain. This operation returns identities th...
  /// HTTP: GET /v2/email/identities
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListEmailIdentitiesResponse> listEmailIdentities(
    $0.ListEmailIdentitiesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEmailIdentities, request, options: options);
  }

  /// Lists the email templates present in your Amazon SES account in the current Amazon Web Services Region. You can execute this operation no more than once per second.
  /// HTTP: GET /v2/email/templates
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListEmailTemplatesResponse> listEmailTemplates(
    $0.ListEmailTemplatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEmailTemplates, request, options: options);
  }

  /// Lists all of the export jobs.
  /// HTTP: POST /v2/email/list-export-jobs
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListExportJobsResponse> listExportJobs(
    $0.ListExportJobsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listExportJobs, request, options: options);
  }

  /// Lists all of the import jobs.
  /// HTTP: POST /v2/email/import-jobs/list
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListImportJobsResponse> listImportJobs(
    $0.ListImportJobsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listImportJobs, request, options: options);
  }

  /// List the multi-region endpoints (global-endpoints). Only multi-region endpoints (global-endpoints) whose primary region is the AWS-Region where operation is executed will be listed.
  /// HTTP: GET /v2/email/multi-region-endpoints
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListMultiRegionEndpointsResponse>
      listMultiRegionEndpoints(
    $0.ListMultiRegionEndpointsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMultiRegionEndpoints, request,
        options: options);
  }

  /// Lists the recommendations present in your Amazon SES account in the current Amazon Web Services Region. You can execute this operation no more than once per second.
  /// HTTP: POST /v2/email/vdm/recommendations
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListRecommendationsResponse> listRecommendations(
    $0.ListRecommendationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRecommendations, request, options: options);
  }

  /// List reputation entities in your Amazon SES account in the current Amazon Web Services Region. You can filter the results by entity type, reputation impact, sending status, or entity reference pref...
  /// HTTP: POST /v2/email/reputation/entities
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListReputationEntitiesResponse>
      listReputationEntities(
    $0.ListReputationEntitiesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listReputationEntities, request,
        options: options);
  }

  /// List all tenants associated with a specific resource. This operation returns a list of tenants that are associated with the specified resource. This is useful for understanding which tenants are cu...
  /// HTTP: POST /v2/email/resources/tenants/list
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListResourceTenantsResponse> listResourceTenants(
    $0.ListResourceTenantsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResourceTenants, request, options: options);
  }

  /// Retrieves a list of email addresses that are on the suppression list for your account.
  /// HTTP: GET /v2/email/suppression/addresses
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListSuppressedDestinationsResponse>
      listSuppressedDestinations(
    $0.ListSuppressedDestinationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSuppressedDestinations, request,
        options: options);
  }

  /// Retrieve a list of the tags (keys and values) that are associated with a specified resource. A tag is a label that you optionally define and associate with a resource. Each tag consists of a requ...
  /// HTTP: GET /v2/email/tags
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// List all resources associated with a specific tenant. This operation returns a list of resources (email identities, configuration sets, or email templates) that are associated with the specified te...
  /// HTTP: POST /v2/email/tenants/resources/list
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListTenantResourcesResponse> listTenantResources(
    $0.ListTenantResourcesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTenantResources, request, options: options);
  }

  /// List all tenants associated with your account in the current Amazon Web Services Region. This operation returns basic information about each tenant, such as tenant name, ID, ARN, and creation times...
  /// HTTP: POST /v2/email/tenants/list
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListTenantsResponse> listTenants(
    $0.ListTenantsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTenants, request, options: options);
  }

  /// Enable or disable the automatic warm-up feature for dedicated IP addresses.
  /// HTTP: PUT /v2/email/account/dedicated-ips/warmup
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutAccountDedicatedIpWarmupAttributesResponse>
      putAccountDedicatedIpWarmupAttributes(
    $0.PutAccountDedicatedIpWarmupAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAccountDedicatedIpWarmupAttributes, request,
        options: options);
  }

  /// Update your Amazon SES account details.
  /// HTTP: POST /v2/email/account/details
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutAccountDetailsResponse> putAccountDetails(
    $0.PutAccountDetailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAccountDetails, request, options: options);
  }

  /// Enable or disable the ability of your account to send email.
  /// HTTP: PUT /v2/email/account/sending
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutAccountSendingAttributesResponse>
      putAccountSendingAttributes(
    $0.PutAccountSendingAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAccountSendingAttributes, request,
        options: options);
  }

  /// Change the settings for the account-level suppression list.
  /// HTTP: PUT /v2/email/account/suppression
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutAccountSuppressionAttributesResponse>
      putAccountSuppressionAttributes(
    $0.PutAccountSuppressionAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAccountSuppressionAttributes, request,
        options: options);
  }

  /// Update your Amazon SES account VDM attributes. You can execute this operation no more than once per second.
  /// HTTP: PUT /v2/email/account/vdm
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutAccountVdmAttributesResponse>
      putAccountVdmAttributes(
    $0.PutAccountVdmAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAccountVdmAttributes, request,
        options: options);
  }

  /// Associate the configuration set with a MailManager archive. When you send email using the SendEmail or SendBulkEmail operations the message as it will be given to the receiving SMTP server will be ...
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/archiving-options
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutConfigurationSetArchivingOptionsResponse>
      putConfigurationSetArchivingOptions(
    $0.PutConfigurationSetArchivingOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putConfigurationSetArchivingOptions, request,
        options: options);
  }

  /// Associate a configuration set with a dedicated IP pool. You can use dedicated IP pools to create groups of dedicated IP addresses for sending specific types of email.
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/delivery-options
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutConfigurationSetDeliveryOptionsResponse>
      putConfigurationSetDeliveryOptions(
    $0.PutConfigurationSetDeliveryOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putConfigurationSetDeliveryOptions, request,
        options: options);
  }

  /// Enable or disable collection of reputation metrics for emails that you send using a particular configuration set in a specific Amazon Web Services Region.
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/reputation-options
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutConfigurationSetReputationOptionsResponse>
      putConfigurationSetReputationOptions(
    $0.PutConfigurationSetReputationOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putConfigurationSetReputationOptions, request,
        options: options);
  }

  /// Enable or disable email sending for messages that use a particular configuration set in a specific Amazon Web Services Region.
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/sending
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutConfigurationSetSendingOptionsResponse>
      putConfigurationSetSendingOptions(
    $0.PutConfigurationSetSendingOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putConfigurationSetSendingOptions, request,
        options: options);
  }

  /// Specify the account suppression list preferences for a configuration set.
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/suppression-options
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutConfigurationSetSuppressionOptionsResponse>
      putConfigurationSetSuppressionOptions(
    $0.PutConfigurationSetSuppressionOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putConfigurationSetSuppressionOptions, request,
        options: options);
  }

  /// Specify a custom domain to use for open and click tracking elements in email that you send.
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/tracking-options
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutConfigurationSetTrackingOptionsResponse>
      putConfigurationSetTrackingOptions(
    $0.PutConfigurationSetTrackingOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putConfigurationSetTrackingOptions, request,
        options: options);
  }

  /// Specify VDM preferences for email that you send using the configuration set. You can execute this operation no more than once per second.
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/vdm-options
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutConfigurationSetVdmOptionsResponse>
      putConfigurationSetVdmOptions(
    $0.PutConfigurationSetVdmOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putConfigurationSetVdmOptions, request,
        options: options);
  }

  /// Move a dedicated IP address to an existing dedicated IP pool. The dedicated IP address that you specify must already exist, and must be associated with your Amazon Web Services account. The dedicat...
  /// HTTP: PUT /v2/email/dedicated-ips/{Ip}/pool
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutDedicatedIpInPoolResponse> putDedicatedIpInPool(
    $0.PutDedicatedIpInPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDedicatedIpInPool, request, options: options);
  }

  /// Used to convert a dedicated IP pool to a different scaling mode. MANAGED pools cannot be converted to STANDARD scaling mode.
  /// HTTP: PUT /v2/email/dedicated-ip-pools/{PoolName}/scaling
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutDedicatedIpPoolScalingAttributesResponse>
      putDedicatedIpPoolScalingAttributes(
    $0.PutDedicatedIpPoolScalingAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDedicatedIpPoolScalingAttributes, request,
        options: options);
  }

  ///
  /// HTTP: PUT /v2/email/dedicated-ips/{Ip}/warmup
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutDedicatedIpWarmupAttributesResponse>
      putDedicatedIpWarmupAttributes(
    $0.PutDedicatedIpWarmupAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDedicatedIpWarmupAttributes, request,
        options: options);
  }

  /// Enable or disable the Deliverability dashboard. When you enable the Deliverability dashboard, you gain access to reputation, deliverability, and other metrics for the domains that you use to send e...
  /// HTTP: PUT /v2/email/deliverability-dashboard
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutDeliverabilityDashboardOptionResponse>
      putDeliverabilityDashboardOption(
    $0.PutDeliverabilityDashboardOptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putDeliverabilityDashboardOption, request,
        options: options);
  }

  /// Used to associate a configuration set with an email identity.
  /// HTTP: PUT /v2/email/identities/{EmailIdentity}/configuration-set
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutEmailIdentityConfigurationSetAttributesResponse>
      putEmailIdentityConfigurationSetAttributes(
    $0.PutEmailIdentityConfigurationSetAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$putEmailIdentityConfigurationSetAttributes, request,
        options: options);
  }

  /// Used to enable or disable DKIM authentication for an email identity.
  /// HTTP: PUT /v2/email/identities/{EmailIdentity}/dkim
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutEmailIdentityDkimAttributesResponse>
      putEmailIdentityDkimAttributes(
    $0.PutEmailIdentityDkimAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putEmailIdentityDkimAttributes, request,
        options: options);
  }

  /// Used to configure or change the DKIM authentication settings for an email domain identity. You can use this operation to do any of the following: Update the signing attributes for an identity that ...
  /// HTTP: PUT /v2/email/identities/{EmailIdentity}/dkim/signing
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutEmailIdentityDkimSigningAttributesResponse>
      putEmailIdentityDkimSigningAttributes(
    $0.PutEmailIdentityDkimSigningAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putEmailIdentityDkimSigningAttributes, request,
        options: options);
  }

  /// Used to enable or disable feedback forwarding for an identity. This setting determines what happens when an identity is used to send an email that results in a bounce or complaint event. If the val...
  /// HTTP: PUT /v2/email/identities/{EmailIdentity}/feedback
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutEmailIdentityFeedbackAttributesResponse>
      putEmailIdentityFeedbackAttributes(
    $0.PutEmailIdentityFeedbackAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putEmailIdentityFeedbackAttributes, request,
        options: options);
  }

  /// Used to enable or disable the custom Mail-From domain configuration for an email identity.
  /// HTTP: PUT /v2/email/identities/{EmailIdentity}/mail-from
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutEmailIdentityMailFromAttributesResponse>
      putEmailIdentityMailFromAttributes(
    $0.PutEmailIdentityMailFromAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putEmailIdentityMailFromAttributes, request,
        options: options);
  }

  /// Adds an email address to the suppression list for your account.
  /// HTTP: PUT /v2/email/suppression/addresses
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutSuppressedDestinationResponse>
      putSuppressedDestination(
    $0.PutSuppressedDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putSuppressedDestination, request,
        options: options);
  }

  /// Composes an email message to multiple destinations.
  /// HTTP: POST /v2/email/outbound-bulk-emails
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SendBulkEmailResponse> sendBulkEmail(
    $0.SendBulkEmailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendBulkEmail, request, options: options);
  }

  /// Adds an email address to the list of identities for your Amazon SES account in the current Amazon Web Services Region and attempts to verify it. As a result of executing this operation, a customize...
  /// HTTP: POST /v2/email/outbound-custom-verification-emails
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SendCustomVerificationEmailResponse>
      sendCustomVerificationEmail(
    $0.SendCustomVerificationEmailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendCustomVerificationEmail, request,
        options: options);
  }

  /// Sends an email message. You can use the Amazon SES API v2 to send the following types of messages: Simple – A standard email message. When you create this type of message, you specify the sender,...
  /// HTTP: POST /v2/email/outbound-emails
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SendEmailResponse> sendEmail(
    $0.SendEmailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendEmail, request, options: options);
  }

  /// Add one or more tags (keys and values) to a specified resource. A tag is a label that you optionally define and associate with a resource. Tags can help you categorize and manage resources in diff...
  /// HTTP: POST /v2/email/tags
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Creates a preview of the MIME content of an email when provided with a template and a set of replacement data. You can execute this operation no more than once per second.
  /// HTTP: POST /v2/email/templates/{TemplateName}/render
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.TestRenderEmailTemplateResponse>
      testRenderEmailTemplate(
    $0.TestRenderEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$testRenderEmailTemplate, request,
        options: options);
  }

  /// Remove one or more tags (keys and values) from a specified resource.
  /// HTTP: DELETE /v2/email/tags
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Update the configuration of an event destination for a configuration set. Events include message sends, deliveries, opens, clicks, bounces, and complaints. Event destinations are places that you ca...
  /// HTTP: PUT /v2/email/configuration-sets/{ConfigurationSetName}/event-destinations/{EventDestinationName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateConfigurationSetEventDestinationResponse>
      updateConfigurationSetEventDestination(
    $0.UpdateConfigurationSetEventDestinationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateConfigurationSetEventDestination, request,
        options: options);
  }

  /// Updates a contact's preferences for a list. You must specify all existing topic preferences in the TopicPreferences object, not just the ones that need updating; otherwise, all your existing prefer...
  /// HTTP: PUT /v2/email/contact-lists/{ContactListName}/contacts/{EmailAddress}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateContactResponse> updateContact(
    $0.UpdateContactRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateContact, request, options: options);
  }

  /// Updates contact list metadata. This operation does a complete replacement.
  /// HTTP: PUT /v2/email/contact-lists/{ContactListName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateContactListResponse> updateContactList(
    $0.UpdateContactListRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateContactList, request, options: options);
  }

  /// Updates an existing custom verification email template. For more information about custom verification email templates, see Using custom verification email templates in the Amazon SES Developer Gui...
  /// HTTP: PUT /v2/email/custom-verification-email-templates/{TemplateName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateCustomVerificationEmailTemplateResponse>
      updateCustomVerificationEmailTemplate(
    $0.UpdateCustomVerificationEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateCustomVerificationEmailTemplate, request,
        options: options);
  }

  /// Updates the specified sending authorization policy for the given identity (an email address or a domain). This API returns successfully even if a policy with the specified name does not exist. This...
  /// HTTP: PUT /v2/email/identities/{EmailIdentity}/policies/{PolicyName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateEmailIdentityPolicyResponse>
      updateEmailIdentityPolicy(
    $0.UpdateEmailIdentityPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateEmailIdentityPolicy, request,
        options: options);
  }

  /// Updates an email template. Email templates enable you to send personalized email to one or more destinations in a single API operation. For more information, see the Amazon SES Developer Guide. You...
  /// HTTP: PUT /v2/email/templates/{TemplateName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateEmailTemplateResponse> updateEmailTemplate(
    $0.UpdateEmailTemplateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateEmailTemplate, request, options: options);
  }

  /// Update the customer-managed sending status for a reputation entity. This allows you to enable, disable, or reinstate sending for the entity. The customer-managed status works in conjunction with th...
  /// HTTP: PUT /v2/email/reputation/entities/{ReputationEntityType}/{ReputationEntityReference}/customer-managed-status
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateReputationEntityCustomerManagedStatusResponse>
      updateReputationEntityCustomerManagedStatus(
    $0.UpdateReputationEntityCustomerManagedStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$updateReputationEntityCustomerManagedStatus, request,
        options: options);
  }

  /// Update the reputation management policy for a reputation entity. The policy determines how the entity responds to reputation findings, such as automatically pausing sending when certain thresholds ...
  /// HTTP: PUT /v2/email/reputation/entities/{ReputationEntityType}/{ReputationEntityReference}/policy
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateReputationEntityPolicyResponse>
      updateReputationEntityPolicy(
    $0.UpdateReputationEntityPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateReputationEntityPolicy, request,
        options: options);
  }

  // method descriptors

  static final _$batchGetMetricData = $grpc.ClientMethod<
          $0.BatchGetMetricDataRequest, $0.BatchGetMetricDataResponse>(
      '/email.SESv2Service/BatchGetMetricData',
      ($0.BatchGetMetricDataRequest value) => value.writeToBuffer(),
      $0.BatchGetMetricDataResponse.fromBuffer);
  static final _$cancelExportJob =
      $grpc.ClientMethod<$0.CancelExportJobRequest, $0.CancelExportJobResponse>(
          '/email.SESv2Service/CancelExportJob',
          ($0.CancelExportJobRequest value) => value.writeToBuffer(),
          $0.CancelExportJobResponse.fromBuffer);
  static final _$createConfigurationSet = $grpc.ClientMethod<
          $0.CreateConfigurationSetRequest, $0.CreateConfigurationSetResponse>(
      '/email.SESv2Service/CreateConfigurationSet',
      ($0.CreateConfigurationSetRequest value) => value.writeToBuffer(),
      $0.CreateConfigurationSetResponse.fromBuffer);
  static final _$createConfigurationSetEventDestination = $grpc.ClientMethod<
          $0.CreateConfigurationSetEventDestinationRequest,
          $0.CreateConfigurationSetEventDestinationResponse>(
      '/email.SESv2Service/CreateConfigurationSetEventDestination',
      ($0.CreateConfigurationSetEventDestinationRequest value) =>
          value.writeToBuffer(),
      $0.CreateConfigurationSetEventDestinationResponse.fromBuffer);
  static final _$createContact =
      $grpc.ClientMethod<$0.CreateContactRequest, $0.CreateContactResponse>(
          '/email.SESv2Service/CreateContact',
          ($0.CreateContactRequest value) => value.writeToBuffer(),
          $0.CreateContactResponse.fromBuffer);
  static final _$createContactList = $grpc.ClientMethod<
          $0.CreateContactListRequest, $0.CreateContactListResponse>(
      '/email.SESv2Service/CreateContactList',
      ($0.CreateContactListRequest value) => value.writeToBuffer(),
      $0.CreateContactListResponse.fromBuffer);
  static final _$createCustomVerificationEmailTemplate = $grpc.ClientMethod<
          $0.CreateCustomVerificationEmailTemplateRequest,
          $0.CreateCustomVerificationEmailTemplateResponse>(
      '/email.SESv2Service/CreateCustomVerificationEmailTemplate',
      ($0.CreateCustomVerificationEmailTemplateRequest value) =>
          value.writeToBuffer(),
      $0.CreateCustomVerificationEmailTemplateResponse.fromBuffer);
  static final _$createDedicatedIpPool = $grpc.ClientMethod<
          $0.CreateDedicatedIpPoolRequest, $0.CreateDedicatedIpPoolResponse>(
      '/email.SESv2Service/CreateDedicatedIpPool',
      ($0.CreateDedicatedIpPoolRequest value) => value.writeToBuffer(),
      $0.CreateDedicatedIpPoolResponse.fromBuffer);
  static final _$createDeliverabilityTestReport = $grpc.ClientMethod<
          $0.CreateDeliverabilityTestReportRequest,
          $0.CreateDeliverabilityTestReportResponse>(
      '/email.SESv2Service/CreateDeliverabilityTestReport',
      ($0.CreateDeliverabilityTestReportRequest value) => value.writeToBuffer(),
      $0.CreateDeliverabilityTestReportResponse.fromBuffer);
  static final _$createEmailIdentity = $grpc.ClientMethod<
          $0.CreateEmailIdentityRequest, $0.CreateEmailIdentityResponse>(
      '/email.SESv2Service/CreateEmailIdentity',
      ($0.CreateEmailIdentityRequest value) => value.writeToBuffer(),
      $0.CreateEmailIdentityResponse.fromBuffer);
  static final _$createEmailIdentityPolicy = $grpc.ClientMethod<
          $0.CreateEmailIdentityPolicyRequest,
          $0.CreateEmailIdentityPolicyResponse>(
      '/email.SESv2Service/CreateEmailIdentityPolicy',
      ($0.CreateEmailIdentityPolicyRequest value) => value.writeToBuffer(),
      $0.CreateEmailIdentityPolicyResponse.fromBuffer);
  static final _$createEmailTemplate = $grpc.ClientMethod<
          $0.CreateEmailTemplateRequest, $0.CreateEmailTemplateResponse>(
      '/email.SESv2Service/CreateEmailTemplate',
      ($0.CreateEmailTemplateRequest value) => value.writeToBuffer(),
      $0.CreateEmailTemplateResponse.fromBuffer);
  static final _$createExportJob =
      $grpc.ClientMethod<$0.CreateExportJobRequest, $0.CreateExportJobResponse>(
          '/email.SESv2Service/CreateExportJob',
          ($0.CreateExportJobRequest value) => value.writeToBuffer(),
          $0.CreateExportJobResponse.fromBuffer);
  static final _$createImportJob =
      $grpc.ClientMethod<$0.CreateImportJobRequest, $0.CreateImportJobResponse>(
          '/email.SESv2Service/CreateImportJob',
          ($0.CreateImportJobRequest value) => value.writeToBuffer(),
          $0.CreateImportJobResponse.fromBuffer);
  static final _$createMultiRegionEndpoint = $grpc.ClientMethod<
          $0.CreateMultiRegionEndpointRequest,
          $0.CreateMultiRegionEndpointResponse>(
      '/email.SESv2Service/CreateMultiRegionEndpoint',
      ($0.CreateMultiRegionEndpointRequest value) => value.writeToBuffer(),
      $0.CreateMultiRegionEndpointResponse.fromBuffer);
  static final _$createTenant =
      $grpc.ClientMethod<$0.CreateTenantRequest, $0.CreateTenantResponse>(
          '/email.SESv2Service/CreateTenant',
          ($0.CreateTenantRequest value) => value.writeToBuffer(),
          $0.CreateTenantResponse.fromBuffer);
  static final _$createTenantResourceAssociation = $grpc.ClientMethod<
          $0.CreateTenantResourceAssociationRequest,
          $0.CreateTenantResourceAssociationResponse>(
      '/email.SESv2Service/CreateTenantResourceAssociation',
      ($0.CreateTenantResourceAssociationRequest value) =>
          value.writeToBuffer(),
      $0.CreateTenantResourceAssociationResponse.fromBuffer);
  static final _$deleteConfigurationSet = $grpc.ClientMethod<
          $0.DeleteConfigurationSetRequest, $0.DeleteConfigurationSetResponse>(
      '/email.SESv2Service/DeleteConfigurationSet',
      ($0.DeleteConfigurationSetRequest value) => value.writeToBuffer(),
      $0.DeleteConfigurationSetResponse.fromBuffer);
  static final _$deleteConfigurationSetEventDestination = $grpc.ClientMethod<
          $0.DeleteConfigurationSetEventDestinationRequest,
          $0.DeleteConfigurationSetEventDestinationResponse>(
      '/email.SESv2Service/DeleteConfigurationSetEventDestination',
      ($0.DeleteConfigurationSetEventDestinationRequest value) =>
          value.writeToBuffer(),
      $0.DeleteConfigurationSetEventDestinationResponse.fromBuffer);
  static final _$deleteContact =
      $grpc.ClientMethod<$0.DeleteContactRequest, $0.DeleteContactResponse>(
          '/email.SESv2Service/DeleteContact',
          ($0.DeleteContactRequest value) => value.writeToBuffer(),
          $0.DeleteContactResponse.fromBuffer);
  static final _$deleteContactList = $grpc.ClientMethod<
          $0.DeleteContactListRequest, $0.DeleteContactListResponse>(
      '/email.SESv2Service/DeleteContactList',
      ($0.DeleteContactListRequest value) => value.writeToBuffer(),
      $0.DeleteContactListResponse.fromBuffer);
  static final _$deleteCustomVerificationEmailTemplate = $grpc.ClientMethod<
          $0.DeleteCustomVerificationEmailTemplateRequest,
          $0.DeleteCustomVerificationEmailTemplateResponse>(
      '/email.SESv2Service/DeleteCustomVerificationEmailTemplate',
      ($0.DeleteCustomVerificationEmailTemplateRequest value) =>
          value.writeToBuffer(),
      $0.DeleteCustomVerificationEmailTemplateResponse.fromBuffer);
  static final _$deleteDedicatedIpPool = $grpc.ClientMethod<
          $0.DeleteDedicatedIpPoolRequest, $0.DeleteDedicatedIpPoolResponse>(
      '/email.SESv2Service/DeleteDedicatedIpPool',
      ($0.DeleteDedicatedIpPoolRequest value) => value.writeToBuffer(),
      $0.DeleteDedicatedIpPoolResponse.fromBuffer);
  static final _$deleteEmailIdentity = $grpc.ClientMethod<
          $0.DeleteEmailIdentityRequest, $0.DeleteEmailIdentityResponse>(
      '/email.SESv2Service/DeleteEmailIdentity',
      ($0.DeleteEmailIdentityRequest value) => value.writeToBuffer(),
      $0.DeleteEmailIdentityResponse.fromBuffer);
  static final _$deleteEmailIdentityPolicy = $grpc.ClientMethod<
          $0.DeleteEmailIdentityPolicyRequest,
          $0.DeleteEmailIdentityPolicyResponse>(
      '/email.SESv2Service/DeleteEmailIdentityPolicy',
      ($0.DeleteEmailIdentityPolicyRequest value) => value.writeToBuffer(),
      $0.DeleteEmailIdentityPolicyResponse.fromBuffer);
  static final _$deleteEmailTemplate = $grpc.ClientMethod<
          $0.DeleteEmailTemplateRequest, $0.DeleteEmailTemplateResponse>(
      '/email.SESv2Service/DeleteEmailTemplate',
      ($0.DeleteEmailTemplateRequest value) => value.writeToBuffer(),
      $0.DeleteEmailTemplateResponse.fromBuffer);
  static final _$deleteMultiRegionEndpoint = $grpc.ClientMethod<
          $0.DeleteMultiRegionEndpointRequest,
          $0.DeleteMultiRegionEndpointResponse>(
      '/email.SESv2Service/DeleteMultiRegionEndpoint',
      ($0.DeleteMultiRegionEndpointRequest value) => value.writeToBuffer(),
      $0.DeleteMultiRegionEndpointResponse.fromBuffer);
  static final _$deleteSuppressedDestination = $grpc.ClientMethod<
          $0.DeleteSuppressedDestinationRequest,
          $0.DeleteSuppressedDestinationResponse>(
      '/email.SESv2Service/DeleteSuppressedDestination',
      ($0.DeleteSuppressedDestinationRequest value) => value.writeToBuffer(),
      $0.DeleteSuppressedDestinationResponse.fromBuffer);
  static final _$deleteTenant =
      $grpc.ClientMethod<$0.DeleteTenantRequest, $0.DeleteTenantResponse>(
          '/email.SESv2Service/DeleteTenant',
          ($0.DeleteTenantRequest value) => value.writeToBuffer(),
          $0.DeleteTenantResponse.fromBuffer);
  static final _$deleteTenantResourceAssociation = $grpc.ClientMethod<
          $0.DeleteTenantResourceAssociationRequest,
          $0.DeleteTenantResourceAssociationResponse>(
      '/email.SESv2Service/DeleteTenantResourceAssociation',
      ($0.DeleteTenantResourceAssociationRequest value) =>
          value.writeToBuffer(),
      $0.DeleteTenantResourceAssociationResponse.fromBuffer);
  static final _$getAccount =
      $grpc.ClientMethod<$0.GetAccountRequest, $0.GetAccountResponse>(
          '/email.SESv2Service/GetAccount',
          ($0.GetAccountRequest value) => value.writeToBuffer(),
          $0.GetAccountResponse.fromBuffer);
  static final _$getBlacklistReports = $grpc.ClientMethod<
          $0.GetBlacklistReportsRequest, $0.GetBlacklistReportsResponse>(
      '/email.SESv2Service/GetBlacklistReports',
      ($0.GetBlacklistReportsRequest value) => value.writeToBuffer(),
      $0.GetBlacklistReportsResponse.fromBuffer);
  static final _$getConfigurationSet = $grpc.ClientMethod<
          $0.GetConfigurationSetRequest, $0.GetConfigurationSetResponse>(
      '/email.SESv2Service/GetConfigurationSet',
      ($0.GetConfigurationSetRequest value) => value.writeToBuffer(),
      $0.GetConfigurationSetResponse.fromBuffer);
  static final _$getConfigurationSetEventDestinations = $grpc.ClientMethod<
          $0.GetConfigurationSetEventDestinationsRequest,
          $0.GetConfigurationSetEventDestinationsResponse>(
      '/email.SESv2Service/GetConfigurationSetEventDestinations',
      ($0.GetConfigurationSetEventDestinationsRequest value) =>
          value.writeToBuffer(),
      $0.GetConfigurationSetEventDestinationsResponse.fromBuffer);
  static final _$getContact =
      $grpc.ClientMethod<$0.GetContactRequest, $0.GetContactResponse>(
          '/email.SESv2Service/GetContact',
          ($0.GetContactRequest value) => value.writeToBuffer(),
          $0.GetContactResponse.fromBuffer);
  static final _$getContactList =
      $grpc.ClientMethod<$0.GetContactListRequest, $0.GetContactListResponse>(
          '/email.SESv2Service/GetContactList',
          ($0.GetContactListRequest value) => value.writeToBuffer(),
          $0.GetContactListResponse.fromBuffer);
  static final _$getCustomVerificationEmailTemplate = $grpc.ClientMethod<
          $0.GetCustomVerificationEmailTemplateRequest,
          $0.GetCustomVerificationEmailTemplateResponse>(
      '/email.SESv2Service/GetCustomVerificationEmailTemplate',
      ($0.GetCustomVerificationEmailTemplateRequest value) =>
          value.writeToBuffer(),
      $0.GetCustomVerificationEmailTemplateResponse.fromBuffer);
  static final _$getDedicatedIp =
      $grpc.ClientMethod<$0.GetDedicatedIpRequest, $0.GetDedicatedIpResponse>(
          '/email.SESv2Service/GetDedicatedIp',
          ($0.GetDedicatedIpRequest value) => value.writeToBuffer(),
          $0.GetDedicatedIpResponse.fromBuffer);
  static final _$getDedicatedIpPool = $grpc.ClientMethod<
          $0.GetDedicatedIpPoolRequest, $0.GetDedicatedIpPoolResponse>(
      '/email.SESv2Service/GetDedicatedIpPool',
      ($0.GetDedicatedIpPoolRequest value) => value.writeToBuffer(),
      $0.GetDedicatedIpPoolResponse.fromBuffer);
  static final _$getDedicatedIps =
      $grpc.ClientMethod<$0.GetDedicatedIpsRequest, $0.GetDedicatedIpsResponse>(
          '/email.SESv2Service/GetDedicatedIps',
          ($0.GetDedicatedIpsRequest value) => value.writeToBuffer(),
          $0.GetDedicatedIpsResponse.fromBuffer);
  static final _$getDeliverabilityDashboardOptions = $grpc.ClientMethod<
          $0.GetDeliverabilityDashboardOptionsRequest,
          $0.GetDeliverabilityDashboardOptionsResponse>(
      '/email.SESv2Service/GetDeliverabilityDashboardOptions',
      ($0.GetDeliverabilityDashboardOptionsRequest value) =>
          value.writeToBuffer(),
      $0.GetDeliverabilityDashboardOptionsResponse.fromBuffer);
  static final _$getDeliverabilityTestReport = $grpc.ClientMethod<
          $0.GetDeliverabilityTestReportRequest,
          $0.GetDeliverabilityTestReportResponse>(
      '/email.SESv2Service/GetDeliverabilityTestReport',
      ($0.GetDeliverabilityTestReportRequest value) => value.writeToBuffer(),
      $0.GetDeliverabilityTestReportResponse.fromBuffer);
  static final _$getDomainDeliverabilityCampaign = $grpc.ClientMethod<
          $0.GetDomainDeliverabilityCampaignRequest,
          $0.GetDomainDeliverabilityCampaignResponse>(
      '/email.SESv2Service/GetDomainDeliverabilityCampaign',
      ($0.GetDomainDeliverabilityCampaignRequest value) =>
          value.writeToBuffer(),
      $0.GetDomainDeliverabilityCampaignResponse.fromBuffer);
  static final _$getDomainStatisticsReport = $grpc.ClientMethod<
          $0.GetDomainStatisticsReportRequest,
          $0.GetDomainStatisticsReportResponse>(
      '/email.SESv2Service/GetDomainStatisticsReport',
      ($0.GetDomainStatisticsReportRequest value) => value.writeToBuffer(),
      $0.GetDomainStatisticsReportResponse.fromBuffer);
  static final _$getEmailAddressInsights = $grpc.ClientMethod<
          $0.GetEmailAddressInsightsRequest,
          $0.GetEmailAddressInsightsResponse>(
      '/email.SESv2Service/GetEmailAddressInsights',
      ($0.GetEmailAddressInsightsRequest value) => value.writeToBuffer(),
      $0.GetEmailAddressInsightsResponse.fromBuffer);
  static final _$getEmailIdentity = $grpc.ClientMethod<
          $0.GetEmailIdentityRequest, $0.GetEmailIdentityResponse>(
      '/email.SESv2Service/GetEmailIdentity',
      ($0.GetEmailIdentityRequest value) => value.writeToBuffer(),
      $0.GetEmailIdentityResponse.fromBuffer);
  static final _$getEmailIdentityPolicies = $grpc.ClientMethod<
          $0.GetEmailIdentityPoliciesRequest,
          $0.GetEmailIdentityPoliciesResponse>(
      '/email.SESv2Service/GetEmailIdentityPolicies',
      ($0.GetEmailIdentityPoliciesRequest value) => value.writeToBuffer(),
      $0.GetEmailIdentityPoliciesResponse.fromBuffer);
  static final _$getEmailTemplate = $grpc.ClientMethod<
          $0.GetEmailTemplateRequest, $0.GetEmailTemplateResponse>(
      '/email.SESv2Service/GetEmailTemplate',
      ($0.GetEmailTemplateRequest value) => value.writeToBuffer(),
      $0.GetEmailTemplateResponse.fromBuffer);
  static final _$getExportJob =
      $grpc.ClientMethod<$0.GetExportJobRequest, $0.GetExportJobResponse>(
          '/email.SESv2Service/GetExportJob',
          ($0.GetExportJobRequest value) => value.writeToBuffer(),
          $0.GetExportJobResponse.fromBuffer);
  static final _$getImportJob =
      $grpc.ClientMethod<$0.GetImportJobRequest, $0.GetImportJobResponse>(
          '/email.SESv2Service/GetImportJob',
          ($0.GetImportJobRequest value) => value.writeToBuffer(),
          $0.GetImportJobResponse.fromBuffer);
  static final _$getMessageInsights = $grpc.ClientMethod<
          $0.GetMessageInsightsRequest, $0.GetMessageInsightsResponse>(
      '/email.SESv2Service/GetMessageInsights',
      ($0.GetMessageInsightsRequest value) => value.writeToBuffer(),
      $0.GetMessageInsightsResponse.fromBuffer);
  static final _$getMultiRegionEndpoint = $grpc.ClientMethod<
          $0.GetMultiRegionEndpointRequest, $0.GetMultiRegionEndpointResponse>(
      '/email.SESv2Service/GetMultiRegionEndpoint',
      ($0.GetMultiRegionEndpointRequest value) => value.writeToBuffer(),
      $0.GetMultiRegionEndpointResponse.fromBuffer);
  static final _$getReputationEntity = $grpc.ClientMethod<
          $0.GetReputationEntityRequest, $0.GetReputationEntityResponse>(
      '/email.SESv2Service/GetReputationEntity',
      ($0.GetReputationEntityRequest value) => value.writeToBuffer(),
      $0.GetReputationEntityResponse.fromBuffer);
  static final _$getSuppressedDestination = $grpc.ClientMethod<
          $0.GetSuppressedDestinationRequest,
          $0.GetSuppressedDestinationResponse>(
      '/email.SESv2Service/GetSuppressedDestination',
      ($0.GetSuppressedDestinationRequest value) => value.writeToBuffer(),
      $0.GetSuppressedDestinationResponse.fromBuffer);
  static final _$getTenant =
      $grpc.ClientMethod<$0.GetTenantRequest, $0.GetTenantResponse>(
          '/email.SESv2Service/GetTenant',
          ($0.GetTenantRequest value) => value.writeToBuffer(),
          $0.GetTenantResponse.fromBuffer);
  static final _$listConfigurationSets = $grpc.ClientMethod<
          $0.ListConfigurationSetsRequest, $0.ListConfigurationSetsResponse>(
      '/email.SESv2Service/ListConfigurationSets',
      ($0.ListConfigurationSetsRequest value) => value.writeToBuffer(),
      $0.ListConfigurationSetsResponse.fromBuffer);
  static final _$listContactLists = $grpc.ClientMethod<
          $0.ListContactListsRequest, $0.ListContactListsResponse>(
      '/email.SESv2Service/ListContactLists',
      ($0.ListContactListsRequest value) => value.writeToBuffer(),
      $0.ListContactListsResponse.fromBuffer);
  static final _$listContacts =
      $grpc.ClientMethod<$0.ListContactsRequest, $0.ListContactsResponse>(
          '/email.SESv2Service/ListContacts',
          ($0.ListContactsRequest value) => value.writeToBuffer(),
          $0.ListContactsResponse.fromBuffer);
  static final _$listCustomVerificationEmailTemplates = $grpc.ClientMethod<
          $0.ListCustomVerificationEmailTemplatesRequest,
          $0.ListCustomVerificationEmailTemplatesResponse>(
      '/email.SESv2Service/ListCustomVerificationEmailTemplates',
      ($0.ListCustomVerificationEmailTemplatesRequest value) =>
          value.writeToBuffer(),
      $0.ListCustomVerificationEmailTemplatesResponse.fromBuffer);
  static final _$listDedicatedIpPools = $grpc.ClientMethod<
          $0.ListDedicatedIpPoolsRequest, $0.ListDedicatedIpPoolsResponse>(
      '/email.SESv2Service/ListDedicatedIpPools',
      ($0.ListDedicatedIpPoolsRequest value) => value.writeToBuffer(),
      $0.ListDedicatedIpPoolsResponse.fromBuffer);
  static final _$listDeliverabilityTestReports = $grpc.ClientMethod<
          $0.ListDeliverabilityTestReportsRequest,
          $0.ListDeliverabilityTestReportsResponse>(
      '/email.SESv2Service/ListDeliverabilityTestReports',
      ($0.ListDeliverabilityTestReportsRequest value) => value.writeToBuffer(),
      $0.ListDeliverabilityTestReportsResponse.fromBuffer);
  static final _$listDomainDeliverabilityCampaigns = $grpc.ClientMethod<
          $0.ListDomainDeliverabilityCampaignsRequest,
          $0.ListDomainDeliverabilityCampaignsResponse>(
      '/email.SESv2Service/ListDomainDeliverabilityCampaigns',
      ($0.ListDomainDeliverabilityCampaignsRequest value) =>
          value.writeToBuffer(),
      $0.ListDomainDeliverabilityCampaignsResponse.fromBuffer);
  static final _$listEmailIdentities = $grpc.ClientMethod<
          $0.ListEmailIdentitiesRequest, $0.ListEmailIdentitiesResponse>(
      '/email.SESv2Service/ListEmailIdentities',
      ($0.ListEmailIdentitiesRequest value) => value.writeToBuffer(),
      $0.ListEmailIdentitiesResponse.fromBuffer);
  static final _$listEmailTemplates = $grpc.ClientMethod<
          $0.ListEmailTemplatesRequest, $0.ListEmailTemplatesResponse>(
      '/email.SESv2Service/ListEmailTemplates',
      ($0.ListEmailTemplatesRequest value) => value.writeToBuffer(),
      $0.ListEmailTemplatesResponse.fromBuffer);
  static final _$listExportJobs =
      $grpc.ClientMethod<$0.ListExportJobsRequest, $0.ListExportJobsResponse>(
          '/email.SESv2Service/ListExportJobs',
          ($0.ListExportJobsRequest value) => value.writeToBuffer(),
          $0.ListExportJobsResponse.fromBuffer);
  static final _$listImportJobs =
      $grpc.ClientMethod<$0.ListImportJobsRequest, $0.ListImportJobsResponse>(
          '/email.SESv2Service/ListImportJobs',
          ($0.ListImportJobsRequest value) => value.writeToBuffer(),
          $0.ListImportJobsResponse.fromBuffer);
  static final _$listMultiRegionEndpoints = $grpc.ClientMethod<
          $0.ListMultiRegionEndpointsRequest,
          $0.ListMultiRegionEndpointsResponse>(
      '/email.SESv2Service/ListMultiRegionEndpoints',
      ($0.ListMultiRegionEndpointsRequest value) => value.writeToBuffer(),
      $0.ListMultiRegionEndpointsResponse.fromBuffer);
  static final _$listRecommendations = $grpc.ClientMethod<
          $0.ListRecommendationsRequest, $0.ListRecommendationsResponse>(
      '/email.SESv2Service/ListRecommendations',
      ($0.ListRecommendationsRequest value) => value.writeToBuffer(),
      $0.ListRecommendationsResponse.fromBuffer);
  static final _$listReputationEntities = $grpc.ClientMethod<
          $0.ListReputationEntitiesRequest, $0.ListReputationEntitiesResponse>(
      '/email.SESv2Service/ListReputationEntities',
      ($0.ListReputationEntitiesRequest value) => value.writeToBuffer(),
      $0.ListReputationEntitiesResponse.fromBuffer);
  static final _$listResourceTenants = $grpc.ClientMethod<
          $0.ListResourceTenantsRequest, $0.ListResourceTenantsResponse>(
      '/email.SESv2Service/ListResourceTenants',
      ($0.ListResourceTenantsRequest value) => value.writeToBuffer(),
      $0.ListResourceTenantsResponse.fromBuffer);
  static final _$listSuppressedDestinations = $grpc.ClientMethod<
          $0.ListSuppressedDestinationsRequest,
          $0.ListSuppressedDestinationsResponse>(
      '/email.SESv2Service/ListSuppressedDestinations',
      ($0.ListSuppressedDestinationsRequest value) => value.writeToBuffer(),
      $0.ListSuppressedDestinationsResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/email.SESv2Service/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listTenantResources = $grpc.ClientMethod<
          $0.ListTenantResourcesRequest, $0.ListTenantResourcesResponse>(
      '/email.SESv2Service/ListTenantResources',
      ($0.ListTenantResourcesRequest value) => value.writeToBuffer(),
      $0.ListTenantResourcesResponse.fromBuffer);
  static final _$listTenants =
      $grpc.ClientMethod<$0.ListTenantsRequest, $0.ListTenantsResponse>(
          '/email.SESv2Service/ListTenants',
          ($0.ListTenantsRequest value) => value.writeToBuffer(),
          $0.ListTenantsResponse.fromBuffer);
  static final _$putAccountDedicatedIpWarmupAttributes = $grpc.ClientMethod<
          $0.PutAccountDedicatedIpWarmupAttributesRequest,
          $0.PutAccountDedicatedIpWarmupAttributesResponse>(
      '/email.SESv2Service/PutAccountDedicatedIpWarmupAttributes',
      ($0.PutAccountDedicatedIpWarmupAttributesRequest value) =>
          value.writeToBuffer(),
      $0.PutAccountDedicatedIpWarmupAttributesResponse.fromBuffer);
  static final _$putAccountDetails = $grpc.ClientMethod<
          $0.PutAccountDetailsRequest, $0.PutAccountDetailsResponse>(
      '/email.SESv2Service/PutAccountDetails',
      ($0.PutAccountDetailsRequest value) => value.writeToBuffer(),
      $0.PutAccountDetailsResponse.fromBuffer);
  static final _$putAccountSendingAttributes = $grpc.ClientMethod<
          $0.PutAccountSendingAttributesRequest,
          $0.PutAccountSendingAttributesResponse>(
      '/email.SESv2Service/PutAccountSendingAttributes',
      ($0.PutAccountSendingAttributesRequest value) => value.writeToBuffer(),
      $0.PutAccountSendingAttributesResponse.fromBuffer);
  static final _$putAccountSuppressionAttributes = $grpc.ClientMethod<
          $0.PutAccountSuppressionAttributesRequest,
          $0.PutAccountSuppressionAttributesResponse>(
      '/email.SESv2Service/PutAccountSuppressionAttributes',
      ($0.PutAccountSuppressionAttributesRequest value) =>
          value.writeToBuffer(),
      $0.PutAccountSuppressionAttributesResponse.fromBuffer);
  static final _$putAccountVdmAttributes = $grpc.ClientMethod<
          $0.PutAccountVdmAttributesRequest,
          $0.PutAccountVdmAttributesResponse>(
      '/email.SESv2Service/PutAccountVdmAttributes',
      ($0.PutAccountVdmAttributesRequest value) => value.writeToBuffer(),
      $0.PutAccountVdmAttributesResponse.fromBuffer);
  static final _$putConfigurationSetArchivingOptions = $grpc.ClientMethod<
          $0.PutConfigurationSetArchivingOptionsRequest,
          $0.PutConfigurationSetArchivingOptionsResponse>(
      '/email.SESv2Service/PutConfigurationSetArchivingOptions',
      ($0.PutConfigurationSetArchivingOptionsRequest value) =>
          value.writeToBuffer(),
      $0.PutConfigurationSetArchivingOptionsResponse.fromBuffer);
  static final _$putConfigurationSetDeliveryOptions = $grpc.ClientMethod<
          $0.PutConfigurationSetDeliveryOptionsRequest,
          $0.PutConfigurationSetDeliveryOptionsResponse>(
      '/email.SESv2Service/PutConfigurationSetDeliveryOptions',
      ($0.PutConfigurationSetDeliveryOptionsRequest value) =>
          value.writeToBuffer(),
      $0.PutConfigurationSetDeliveryOptionsResponse.fromBuffer);
  static final _$putConfigurationSetReputationOptions = $grpc.ClientMethod<
          $0.PutConfigurationSetReputationOptionsRequest,
          $0.PutConfigurationSetReputationOptionsResponse>(
      '/email.SESv2Service/PutConfigurationSetReputationOptions',
      ($0.PutConfigurationSetReputationOptionsRequest value) =>
          value.writeToBuffer(),
      $0.PutConfigurationSetReputationOptionsResponse.fromBuffer);
  static final _$putConfigurationSetSendingOptions = $grpc.ClientMethod<
          $0.PutConfigurationSetSendingOptionsRequest,
          $0.PutConfigurationSetSendingOptionsResponse>(
      '/email.SESv2Service/PutConfigurationSetSendingOptions',
      ($0.PutConfigurationSetSendingOptionsRequest value) =>
          value.writeToBuffer(),
      $0.PutConfigurationSetSendingOptionsResponse.fromBuffer);
  static final _$putConfigurationSetSuppressionOptions = $grpc.ClientMethod<
          $0.PutConfigurationSetSuppressionOptionsRequest,
          $0.PutConfigurationSetSuppressionOptionsResponse>(
      '/email.SESv2Service/PutConfigurationSetSuppressionOptions',
      ($0.PutConfigurationSetSuppressionOptionsRequest value) =>
          value.writeToBuffer(),
      $0.PutConfigurationSetSuppressionOptionsResponse.fromBuffer);
  static final _$putConfigurationSetTrackingOptions = $grpc.ClientMethod<
          $0.PutConfigurationSetTrackingOptionsRequest,
          $0.PutConfigurationSetTrackingOptionsResponse>(
      '/email.SESv2Service/PutConfigurationSetTrackingOptions',
      ($0.PutConfigurationSetTrackingOptionsRequest value) =>
          value.writeToBuffer(),
      $0.PutConfigurationSetTrackingOptionsResponse.fromBuffer);
  static final _$putConfigurationSetVdmOptions = $grpc.ClientMethod<
          $0.PutConfigurationSetVdmOptionsRequest,
          $0.PutConfigurationSetVdmOptionsResponse>(
      '/email.SESv2Service/PutConfigurationSetVdmOptions',
      ($0.PutConfigurationSetVdmOptionsRequest value) => value.writeToBuffer(),
      $0.PutConfigurationSetVdmOptionsResponse.fromBuffer);
  static final _$putDedicatedIpInPool = $grpc.ClientMethod<
          $0.PutDedicatedIpInPoolRequest, $0.PutDedicatedIpInPoolResponse>(
      '/email.SESv2Service/PutDedicatedIpInPool',
      ($0.PutDedicatedIpInPoolRequest value) => value.writeToBuffer(),
      $0.PutDedicatedIpInPoolResponse.fromBuffer);
  static final _$putDedicatedIpPoolScalingAttributes = $grpc.ClientMethod<
          $0.PutDedicatedIpPoolScalingAttributesRequest,
          $0.PutDedicatedIpPoolScalingAttributesResponse>(
      '/email.SESv2Service/PutDedicatedIpPoolScalingAttributes',
      ($0.PutDedicatedIpPoolScalingAttributesRequest value) =>
          value.writeToBuffer(),
      $0.PutDedicatedIpPoolScalingAttributesResponse.fromBuffer);
  static final _$putDedicatedIpWarmupAttributes = $grpc.ClientMethod<
          $0.PutDedicatedIpWarmupAttributesRequest,
          $0.PutDedicatedIpWarmupAttributesResponse>(
      '/email.SESv2Service/PutDedicatedIpWarmupAttributes',
      ($0.PutDedicatedIpWarmupAttributesRequest value) => value.writeToBuffer(),
      $0.PutDedicatedIpWarmupAttributesResponse.fromBuffer);
  static final _$putDeliverabilityDashboardOption = $grpc.ClientMethod<
          $0.PutDeliverabilityDashboardOptionRequest,
          $0.PutDeliverabilityDashboardOptionResponse>(
      '/email.SESv2Service/PutDeliverabilityDashboardOption',
      ($0.PutDeliverabilityDashboardOptionRequest value) =>
          value.writeToBuffer(),
      $0.PutDeliverabilityDashboardOptionResponse.fromBuffer);
  static final _$putEmailIdentityConfigurationSetAttributes =
      $grpc.ClientMethod<$0.PutEmailIdentityConfigurationSetAttributesRequest,
              $0.PutEmailIdentityConfigurationSetAttributesResponse>(
          '/email.SESv2Service/PutEmailIdentityConfigurationSetAttributes',
          ($0.PutEmailIdentityConfigurationSetAttributesRequest value) =>
              value.writeToBuffer(),
          $0.PutEmailIdentityConfigurationSetAttributesResponse.fromBuffer);
  static final _$putEmailIdentityDkimAttributes = $grpc.ClientMethod<
          $0.PutEmailIdentityDkimAttributesRequest,
          $0.PutEmailIdentityDkimAttributesResponse>(
      '/email.SESv2Service/PutEmailIdentityDkimAttributes',
      ($0.PutEmailIdentityDkimAttributesRequest value) => value.writeToBuffer(),
      $0.PutEmailIdentityDkimAttributesResponse.fromBuffer);
  static final _$putEmailIdentityDkimSigningAttributes = $grpc.ClientMethod<
          $0.PutEmailIdentityDkimSigningAttributesRequest,
          $0.PutEmailIdentityDkimSigningAttributesResponse>(
      '/email.SESv2Service/PutEmailIdentityDkimSigningAttributes',
      ($0.PutEmailIdentityDkimSigningAttributesRequest value) =>
          value.writeToBuffer(),
      $0.PutEmailIdentityDkimSigningAttributesResponse.fromBuffer);
  static final _$putEmailIdentityFeedbackAttributes = $grpc.ClientMethod<
          $0.PutEmailIdentityFeedbackAttributesRequest,
          $0.PutEmailIdentityFeedbackAttributesResponse>(
      '/email.SESv2Service/PutEmailIdentityFeedbackAttributes',
      ($0.PutEmailIdentityFeedbackAttributesRequest value) =>
          value.writeToBuffer(),
      $0.PutEmailIdentityFeedbackAttributesResponse.fromBuffer);
  static final _$putEmailIdentityMailFromAttributes = $grpc.ClientMethod<
          $0.PutEmailIdentityMailFromAttributesRequest,
          $0.PutEmailIdentityMailFromAttributesResponse>(
      '/email.SESv2Service/PutEmailIdentityMailFromAttributes',
      ($0.PutEmailIdentityMailFromAttributesRequest value) =>
          value.writeToBuffer(),
      $0.PutEmailIdentityMailFromAttributesResponse.fromBuffer);
  static final _$putSuppressedDestination = $grpc.ClientMethod<
          $0.PutSuppressedDestinationRequest,
          $0.PutSuppressedDestinationResponse>(
      '/email.SESv2Service/PutSuppressedDestination',
      ($0.PutSuppressedDestinationRequest value) => value.writeToBuffer(),
      $0.PutSuppressedDestinationResponse.fromBuffer);
  static final _$sendBulkEmail =
      $grpc.ClientMethod<$0.SendBulkEmailRequest, $0.SendBulkEmailResponse>(
          '/email.SESv2Service/SendBulkEmail',
          ($0.SendBulkEmailRequest value) => value.writeToBuffer(),
          $0.SendBulkEmailResponse.fromBuffer);
  static final _$sendCustomVerificationEmail = $grpc.ClientMethod<
          $0.SendCustomVerificationEmailRequest,
          $0.SendCustomVerificationEmailResponse>(
      '/email.SESv2Service/SendCustomVerificationEmail',
      ($0.SendCustomVerificationEmailRequest value) => value.writeToBuffer(),
      $0.SendCustomVerificationEmailResponse.fromBuffer);
  static final _$sendEmail =
      $grpc.ClientMethod<$0.SendEmailRequest, $0.SendEmailResponse>(
          '/email.SESv2Service/SendEmail',
          ($0.SendEmailRequest value) => value.writeToBuffer(),
          $0.SendEmailResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/email.SESv2Service/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$testRenderEmailTemplate = $grpc.ClientMethod<
          $0.TestRenderEmailTemplateRequest,
          $0.TestRenderEmailTemplateResponse>(
      '/email.SESv2Service/TestRenderEmailTemplate',
      ($0.TestRenderEmailTemplateRequest value) => value.writeToBuffer(),
      $0.TestRenderEmailTemplateResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/email.SESv2Service/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateConfigurationSetEventDestination = $grpc.ClientMethod<
          $0.UpdateConfigurationSetEventDestinationRequest,
          $0.UpdateConfigurationSetEventDestinationResponse>(
      '/email.SESv2Service/UpdateConfigurationSetEventDestination',
      ($0.UpdateConfigurationSetEventDestinationRequest value) =>
          value.writeToBuffer(),
      $0.UpdateConfigurationSetEventDestinationResponse.fromBuffer);
  static final _$updateContact =
      $grpc.ClientMethod<$0.UpdateContactRequest, $0.UpdateContactResponse>(
          '/email.SESv2Service/UpdateContact',
          ($0.UpdateContactRequest value) => value.writeToBuffer(),
          $0.UpdateContactResponse.fromBuffer);
  static final _$updateContactList = $grpc.ClientMethod<
          $0.UpdateContactListRequest, $0.UpdateContactListResponse>(
      '/email.SESv2Service/UpdateContactList',
      ($0.UpdateContactListRequest value) => value.writeToBuffer(),
      $0.UpdateContactListResponse.fromBuffer);
  static final _$updateCustomVerificationEmailTemplate = $grpc.ClientMethod<
          $0.UpdateCustomVerificationEmailTemplateRequest,
          $0.UpdateCustomVerificationEmailTemplateResponse>(
      '/email.SESv2Service/UpdateCustomVerificationEmailTemplate',
      ($0.UpdateCustomVerificationEmailTemplateRequest value) =>
          value.writeToBuffer(),
      $0.UpdateCustomVerificationEmailTemplateResponse.fromBuffer);
  static final _$updateEmailIdentityPolicy = $grpc.ClientMethod<
          $0.UpdateEmailIdentityPolicyRequest,
          $0.UpdateEmailIdentityPolicyResponse>(
      '/email.SESv2Service/UpdateEmailIdentityPolicy',
      ($0.UpdateEmailIdentityPolicyRequest value) => value.writeToBuffer(),
      $0.UpdateEmailIdentityPolicyResponse.fromBuffer);
  static final _$updateEmailTemplate = $grpc.ClientMethod<
          $0.UpdateEmailTemplateRequest, $0.UpdateEmailTemplateResponse>(
      '/email.SESv2Service/UpdateEmailTemplate',
      ($0.UpdateEmailTemplateRequest value) => value.writeToBuffer(),
      $0.UpdateEmailTemplateResponse.fromBuffer);
  static final _$updateReputationEntityCustomerManagedStatus =
      $grpc.ClientMethod<$0.UpdateReputationEntityCustomerManagedStatusRequest,
              $0.UpdateReputationEntityCustomerManagedStatusResponse>(
          '/email.SESv2Service/UpdateReputationEntityCustomerManagedStatus',
          ($0.UpdateReputationEntityCustomerManagedStatusRequest value) =>
              value.writeToBuffer(),
          $0.UpdateReputationEntityCustomerManagedStatusResponse.fromBuffer);
  static final _$updateReputationEntityPolicy = $grpc.ClientMethod<
          $0.UpdateReputationEntityPolicyRequest,
          $0.UpdateReputationEntityPolicyResponse>(
      '/email.SESv2Service/UpdateReputationEntityPolicy',
      ($0.UpdateReputationEntityPolicyRequest value) => value.writeToBuffer(),
      $0.UpdateReputationEntityPolicyResponse.fromBuffer);
}

@$pb.GrpcServiceName('email.SESv2Service')
abstract class SESv2ServiceBase extends $grpc.Service {
  $core.String get $name => 'email.SESv2Service';

  SESv2ServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.BatchGetMetricDataRequest,
            $0.BatchGetMetricDataResponse>(
        'BatchGetMetricData',
        batchGetMetricData_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.BatchGetMetricDataRequest.fromBuffer(value),
        ($0.BatchGetMetricDataResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelExportJobRequest,
            $0.CancelExportJobResponse>(
        'CancelExportJob',
        cancelExportJob_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelExportJobRequest.fromBuffer(value),
        ($0.CancelExportJobResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateConfigurationSetRequest,
            $0.CreateConfigurationSetResponse>(
        'CreateConfigurationSet',
        createConfigurationSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateConfigurationSetRequest.fromBuffer(value),
        ($0.CreateConfigurationSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.CreateConfigurationSetEventDestinationRequest,
            $0.CreateConfigurationSetEventDestinationResponse>(
        'CreateConfigurationSetEventDestination',
        createConfigurationSetEventDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateConfigurationSetEventDestinationRequest.fromBuffer(value),
        ($0.CreateConfigurationSetEventDestinationResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateContactRequest, $0.CreateContactResponse>(
            'CreateContact',
            createContact_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateContactRequest.fromBuffer(value),
            ($0.CreateContactResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateContactListRequest,
            $0.CreateContactListResponse>(
        'CreateContactList',
        createContactList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateContactListRequest.fromBuffer(value),
        ($0.CreateContactListResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.CreateCustomVerificationEmailTemplateRequest,
            $0.CreateCustomVerificationEmailTemplateResponse>(
        'CreateCustomVerificationEmailTemplate',
        createCustomVerificationEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCustomVerificationEmailTemplateRequest.fromBuffer(value),
        ($0.CreateCustomVerificationEmailTemplateResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDedicatedIpPoolRequest,
            $0.CreateDedicatedIpPoolResponse>(
        'CreateDedicatedIpPool',
        createDedicatedIpPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDedicatedIpPoolRequest.fromBuffer(value),
        ($0.CreateDedicatedIpPoolResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDeliverabilityTestReportRequest,
            $0.CreateDeliverabilityTestReportResponse>(
        'CreateDeliverabilityTestReport',
        createDeliverabilityTestReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDeliverabilityTestReportRequest.fromBuffer(value),
        ($0.CreateDeliverabilityTestReportResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateEmailIdentityRequest,
            $0.CreateEmailIdentityResponse>(
        'CreateEmailIdentity',
        createEmailIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateEmailIdentityRequest.fromBuffer(value),
        ($0.CreateEmailIdentityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateEmailIdentityPolicyRequest,
            $0.CreateEmailIdentityPolicyResponse>(
        'CreateEmailIdentityPolicy',
        createEmailIdentityPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateEmailIdentityPolicyRequest.fromBuffer(value),
        ($0.CreateEmailIdentityPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateEmailTemplateRequest,
            $0.CreateEmailTemplateResponse>(
        'CreateEmailTemplate',
        createEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateEmailTemplateRequest.fromBuffer(value),
        ($0.CreateEmailTemplateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateExportJobRequest,
            $0.CreateExportJobResponse>(
        'CreateExportJob',
        createExportJob_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateExportJobRequest.fromBuffer(value),
        ($0.CreateExportJobResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateImportJobRequest,
            $0.CreateImportJobResponse>(
        'CreateImportJob',
        createImportJob_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateImportJobRequest.fromBuffer(value),
        ($0.CreateImportJobResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateMultiRegionEndpointRequest,
            $0.CreateMultiRegionEndpointResponse>(
        'CreateMultiRegionEndpoint',
        createMultiRegionEndpoint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateMultiRegionEndpointRequest.fromBuffer(value),
        ($0.CreateMultiRegionEndpointResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateTenantRequest, $0.CreateTenantResponse>(
            'CreateTenant',
            createTenant_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateTenantRequest.fromBuffer(value),
            ($0.CreateTenantResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateTenantResourceAssociationRequest,
            $0.CreateTenantResourceAssociationResponse>(
        'CreateTenantResourceAssociation',
        createTenantResourceAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateTenantResourceAssociationRequest.fromBuffer(value),
        ($0.CreateTenantResourceAssociationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteConfigurationSetRequest,
            $0.DeleteConfigurationSetResponse>(
        'DeleteConfigurationSet',
        deleteConfigurationSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteConfigurationSetRequest.fromBuffer(value),
        ($0.DeleteConfigurationSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeleteConfigurationSetEventDestinationRequest,
            $0.DeleteConfigurationSetEventDestinationResponse>(
        'DeleteConfigurationSetEventDestination',
        deleteConfigurationSetEventDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteConfigurationSetEventDestinationRequest.fromBuffer(value),
        ($0.DeleteConfigurationSetEventDestinationResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteContactRequest, $0.DeleteContactResponse>(
            'DeleteContact',
            deleteContact_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteContactRequest.fromBuffer(value),
            ($0.DeleteContactResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteContactListRequest,
            $0.DeleteContactListResponse>(
        'DeleteContactList',
        deleteContactList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteContactListRequest.fromBuffer(value),
        ($0.DeleteContactListResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeleteCustomVerificationEmailTemplateRequest,
            $0.DeleteCustomVerificationEmailTemplateResponse>(
        'DeleteCustomVerificationEmailTemplate',
        deleteCustomVerificationEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCustomVerificationEmailTemplateRequest.fromBuffer(value),
        ($0.DeleteCustomVerificationEmailTemplateResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDedicatedIpPoolRequest,
            $0.DeleteDedicatedIpPoolResponse>(
        'DeleteDedicatedIpPool',
        deleteDedicatedIpPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDedicatedIpPoolRequest.fromBuffer(value),
        ($0.DeleteDedicatedIpPoolResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteEmailIdentityRequest,
            $0.DeleteEmailIdentityResponse>(
        'DeleteEmailIdentity',
        deleteEmailIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteEmailIdentityRequest.fromBuffer(value),
        ($0.DeleteEmailIdentityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteEmailIdentityPolicyRequest,
            $0.DeleteEmailIdentityPolicyResponse>(
        'DeleteEmailIdentityPolicy',
        deleteEmailIdentityPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteEmailIdentityPolicyRequest.fromBuffer(value),
        ($0.DeleteEmailIdentityPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteEmailTemplateRequest,
            $0.DeleteEmailTemplateResponse>(
        'DeleteEmailTemplate',
        deleteEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteEmailTemplateRequest.fromBuffer(value),
        ($0.DeleteEmailTemplateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMultiRegionEndpointRequest,
            $0.DeleteMultiRegionEndpointResponse>(
        'DeleteMultiRegionEndpoint',
        deleteMultiRegionEndpoint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMultiRegionEndpointRequest.fromBuffer(value),
        ($0.DeleteMultiRegionEndpointResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteSuppressedDestinationRequest,
            $0.DeleteSuppressedDestinationResponse>(
        'DeleteSuppressedDestination',
        deleteSuppressedDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteSuppressedDestinationRequest.fromBuffer(value),
        ($0.DeleteSuppressedDestinationResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteTenantRequest, $0.DeleteTenantResponse>(
            'DeleteTenant',
            deleteTenant_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteTenantRequest.fromBuffer(value),
            ($0.DeleteTenantResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTenantResourceAssociationRequest,
            $0.DeleteTenantResourceAssociationResponse>(
        'DeleteTenantResourceAssociation',
        deleteTenantResourceAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteTenantResourceAssociationRequest.fromBuffer(value),
        ($0.DeleteTenantResourceAssociationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccountRequest, $0.GetAccountResponse>(
        'GetAccount',
        getAccount_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetAccountRequest.fromBuffer(value),
        ($0.GetAccountResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetBlacklistReportsRequest,
            $0.GetBlacklistReportsResponse>(
        'GetBlacklistReports',
        getBlacklistReports_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetBlacklistReportsRequest.fromBuffer(value),
        ($0.GetBlacklistReportsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConfigurationSetRequest,
            $0.GetConfigurationSetResponse>(
        'GetConfigurationSet',
        getConfigurationSet_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConfigurationSetRequest.fromBuffer(value),
        ($0.GetConfigurationSetResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetConfigurationSetEventDestinationsRequest,
            $0.GetConfigurationSetEventDestinationsResponse>(
        'GetConfigurationSetEventDestinations',
        getConfigurationSetEventDestinations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConfigurationSetEventDestinationsRequest.fromBuffer(value),
        ($0.GetConfigurationSetEventDestinationsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetContactRequest, $0.GetContactResponse>(
        'GetContact',
        getContact_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetContactRequest.fromBuffer(value),
        ($0.GetContactResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetContactListRequest,
            $0.GetContactListResponse>(
        'GetContactList',
        getContactList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetContactListRequest.fromBuffer(value),
        ($0.GetContactListResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCustomVerificationEmailTemplateRequest,
            $0.GetCustomVerificationEmailTemplateResponse>(
        'GetCustomVerificationEmailTemplate',
        getCustomVerificationEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCustomVerificationEmailTemplateRequest.fromBuffer(value),
        ($0.GetCustomVerificationEmailTemplateResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDedicatedIpRequest,
            $0.GetDedicatedIpResponse>(
        'GetDedicatedIp',
        getDedicatedIp_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDedicatedIpRequest.fromBuffer(value),
        ($0.GetDedicatedIpResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDedicatedIpPoolRequest,
            $0.GetDedicatedIpPoolResponse>(
        'GetDedicatedIpPool',
        getDedicatedIpPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDedicatedIpPoolRequest.fromBuffer(value),
        ($0.GetDedicatedIpPoolResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDedicatedIpsRequest,
            $0.GetDedicatedIpsResponse>(
        'GetDedicatedIps',
        getDedicatedIps_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDedicatedIpsRequest.fromBuffer(value),
        ($0.GetDedicatedIpsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeliverabilityDashboardOptionsRequest,
            $0.GetDeliverabilityDashboardOptionsResponse>(
        'GetDeliverabilityDashboardOptions',
        getDeliverabilityDashboardOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeliverabilityDashboardOptionsRequest.fromBuffer(value),
        ($0.GetDeliverabilityDashboardOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeliverabilityTestReportRequest,
            $0.GetDeliverabilityTestReportResponse>(
        'GetDeliverabilityTestReport',
        getDeliverabilityTestReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeliverabilityTestReportRequest.fromBuffer(value),
        ($0.GetDeliverabilityTestReportResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDomainDeliverabilityCampaignRequest,
            $0.GetDomainDeliverabilityCampaignResponse>(
        'GetDomainDeliverabilityCampaign',
        getDomainDeliverabilityCampaign_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDomainDeliverabilityCampaignRequest.fromBuffer(value),
        ($0.GetDomainDeliverabilityCampaignResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDomainStatisticsReportRequest,
            $0.GetDomainStatisticsReportResponse>(
        'GetDomainStatisticsReport',
        getDomainStatisticsReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDomainStatisticsReportRequest.fromBuffer(value),
        ($0.GetDomainStatisticsReportResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEmailAddressInsightsRequest,
            $0.GetEmailAddressInsightsResponse>(
        'GetEmailAddressInsights',
        getEmailAddressInsights_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEmailAddressInsightsRequest.fromBuffer(value),
        ($0.GetEmailAddressInsightsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEmailIdentityRequest,
            $0.GetEmailIdentityResponse>(
        'GetEmailIdentity',
        getEmailIdentity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEmailIdentityRequest.fromBuffer(value),
        ($0.GetEmailIdentityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEmailIdentityPoliciesRequest,
            $0.GetEmailIdentityPoliciesResponse>(
        'GetEmailIdentityPolicies',
        getEmailIdentityPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEmailIdentityPoliciesRequest.fromBuffer(value),
        ($0.GetEmailIdentityPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEmailTemplateRequest,
            $0.GetEmailTemplateResponse>(
        'GetEmailTemplate',
        getEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEmailTemplateRequest.fromBuffer(value),
        ($0.GetEmailTemplateResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetExportJobRequest, $0.GetExportJobResponse>(
            'GetExportJob',
            getExportJob_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetExportJobRequest.fromBuffer(value),
            ($0.GetExportJobResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetImportJobRequest, $0.GetImportJobResponse>(
            'GetImportJob',
            getImportJob_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetImportJobRequest.fromBuffer(value),
            ($0.GetImportJobResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMessageInsightsRequest,
            $0.GetMessageInsightsResponse>(
        'GetMessageInsights',
        getMessageInsights_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMessageInsightsRequest.fromBuffer(value),
        ($0.GetMessageInsightsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMultiRegionEndpointRequest,
            $0.GetMultiRegionEndpointResponse>(
        'GetMultiRegionEndpoint',
        getMultiRegionEndpoint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMultiRegionEndpointRequest.fromBuffer(value),
        ($0.GetMultiRegionEndpointResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetReputationEntityRequest,
            $0.GetReputationEntityResponse>(
        'GetReputationEntity',
        getReputationEntity_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetReputationEntityRequest.fromBuffer(value),
        ($0.GetReputationEntityResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSuppressedDestinationRequest,
            $0.GetSuppressedDestinationResponse>(
        'GetSuppressedDestination',
        getSuppressedDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSuppressedDestinationRequest.fromBuffer(value),
        ($0.GetSuppressedDestinationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTenantRequest, $0.GetTenantResponse>(
        'GetTenant',
        getTenant_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetTenantRequest.fromBuffer(value),
        ($0.GetTenantResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListConfigurationSetsRequest,
            $0.ListConfigurationSetsResponse>(
        'ListConfigurationSets',
        listConfigurationSets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListConfigurationSetsRequest.fromBuffer(value),
        ($0.ListConfigurationSetsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListContactListsRequest,
            $0.ListContactListsResponse>(
        'ListContactLists',
        listContactLists_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListContactListsRequest.fromBuffer(value),
        ($0.ListContactListsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListContactsRequest, $0.ListContactsResponse>(
            'ListContacts',
            listContacts_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListContactsRequest.fromBuffer(value),
            ($0.ListContactsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListCustomVerificationEmailTemplatesRequest,
            $0.ListCustomVerificationEmailTemplatesResponse>(
        'ListCustomVerificationEmailTemplates',
        listCustomVerificationEmailTemplates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCustomVerificationEmailTemplatesRequest.fromBuffer(value),
        ($0.ListCustomVerificationEmailTemplatesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDedicatedIpPoolsRequest,
            $0.ListDedicatedIpPoolsResponse>(
        'ListDedicatedIpPools',
        listDedicatedIpPools_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDedicatedIpPoolsRequest.fromBuffer(value),
        ($0.ListDedicatedIpPoolsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDeliverabilityTestReportsRequest,
            $0.ListDeliverabilityTestReportsResponse>(
        'ListDeliverabilityTestReports',
        listDeliverabilityTestReports_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDeliverabilityTestReportsRequest.fromBuffer(value),
        ($0.ListDeliverabilityTestReportsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDomainDeliverabilityCampaignsRequest,
            $0.ListDomainDeliverabilityCampaignsResponse>(
        'ListDomainDeliverabilityCampaigns',
        listDomainDeliverabilityCampaigns_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDomainDeliverabilityCampaignsRequest.fromBuffer(value),
        ($0.ListDomainDeliverabilityCampaignsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEmailIdentitiesRequest,
            $0.ListEmailIdentitiesResponse>(
        'ListEmailIdentities',
        listEmailIdentities_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEmailIdentitiesRequest.fromBuffer(value),
        ($0.ListEmailIdentitiesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEmailTemplatesRequest,
            $0.ListEmailTemplatesResponse>(
        'ListEmailTemplates',
        listEmailTemplates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEmailTemplatesRequest.fromBuffer(value),
        ($0.ListEmailTemplatesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListExportJobsRequest,
            $0.ListExportJobsResponse>(
        'ListExportJobs',
        listExportJobs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListExportJobsRequest.fromBuffer(value),
        ($0.ListExportJobsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListImportJobsRequest,
            $0.ListImportJobsResponse>(
        'ListImportJobs',
        listImportJobs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListImportJobsRequest.fromBuffer(value),
        ($0.ListImportJobsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMultiRegionEndpointsRequest,
            $0.ListMultiRegionEndpointsResponse>(
        'ListMultiRegionEndpoints',
        listMultiRegionEndpoints_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListMultiRegionEndpointsRequest.fromBuffer(value),
        ($0.ListMultiRegionEndpointsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRecommendationsRequest,
            $0.ListRecommendationsResponse>(
        'ListRecommendations',
        listRecommendations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRecommendationsRequest.fromBuffer(value),
        ($0.ListRecommendationsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListReputationEntitiesRequest,
            $0.ListReputationEntitiesResponse>(
        'ListReputationEntities',
        listReputationEntities_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListReputationEntitiesRequest.fromBuffer(value),
        ($0.ListReputationEntitiesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResourceTenantsRequest,
            $0.ListResourceTenantsResponse>(
        'ListResourceTenants',
        listResourceTenants_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResourceTenantsRequest.fromBuffer(value),
        ($0.ListResourceTenantsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSuppressedDestinationsRequest,
            $0.ListSuppressedDestinationsResponse>(
        'ListSuppressedDestinations',
        listSuppressedDestinations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSuppressedDestinationsRequest.fromBuffer(value),
        ($0.ListSuppressedDestinationsResponse value) =>
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
    $addMethod($grpc.ServiceMethod<$0.ListTenantResourcesRequest,
            $0.ListTenantResourcesResponse>(
        'ListTenantResources',
        listTenantResources_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTenantResourcesRequest.fromBuffer(value),
        ($0.ListTenantResourcesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListTenantsRequest, $0.ListTenantsResponse>(
            'ListTenants',
            listTenants_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListTenantsRequest.fromBuffer(value),
            ($0.ListTenantsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutAccountDedicatedIpWarmupAttributesRequest,
            $0.PutAccountDedicatedIpWarmupAttributesResponse>(
        'PutAccountDedicatedIpWarmupAttributes',
        putAccountDedicatedIpWarmupAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAccountDedicatedIpWarmupAttributesRequest.fromBuffer(value),
        ($0.PutAccountDedicatedIpWarmupAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAccountDetailsRequest,
            $0.PutAccountDetailsResponse>(
        'PutAccountDetails',
        putAccountDetails_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAccountDetailsRequest.fromBuffer(value),
        ($0.PutAccountDetailsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAccountSendingAttributesRequest,
            $0.PutAccountSendingAttributesResponse>(
        'PutAccountSendingAttributes',
        putAccountSendingAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAccountSendingAttributesRequest.fromBuffer(value),
        ($0.PutAccountSendingAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAccountSuppressionAttributesRequest,
            $0.PutAccountSuppressionAttributesResponse>(
        'PutAccountSuppressionAttributes',
        putAccountSuppressionAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAccountSuppressionAttributesRequest.fromBuffer(value),
        ($0.PutAccountSuppressionAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAccountVdmAttributesRequest,
            $0.PutAccountVdmAttributesResponse>(
        'PutAccountVdmAttributes',
        putAccountVdmAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAccountVdmAttributesRequest.fromBuffer(value),
        ($0.PutAccountVdmAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutConfigurationSetArchivingOptionsRequest,
            $0.PutConfigurationSetArchivingOptionsResponse>(
        'PutConfigurationSetArchivingOptions',
        putConfigurationSetArchivingOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutConfigurationSetArchivingOptionsRequest.fromBuffer(value),
        ($0.PutConfigurationSetArchivingOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutConfigurationSetDeliveryOptionsRequest,
            $0.PutConfigurationSetDeliveryOptionsResponse>(
        'PutConfigurationSetDeliveryOptions',
        putConfigurationSetDeliveryOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutConfigurationSetDeliveryOptionsRequest.fromBuffer(value),
        ($0.PutConfigurationSetDeliveryOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutConfigurationSetReputationOptionsRequest,
            $0.PutConfigurationSetReputationOptionsResponse>(
        'PutConfigurationSetReputationOptions',
        putConfigurationSetReputationOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutConfigurationSetReputationOptionsRequest.fromBuffer(value),
        ($0.PutConfigurationSetReputationOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutConfigurationSetSendingOptionsRequest,
            $0.PutConfigurationSetSendingOptionsResponse>(
        'PutConfigurationSetSendingOptions',
        putConfigurationSetSendingOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutConfigurationSetSendingOptionsRequest.fromBuffer(value),
        ($0.PutConfigurationSetSendingOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutConfigurationSetSuppressionOptionsRequest,
            $0.PutConfigurationSetSuppressionOptionsResponse>(
        'PutConfigurationSetSuppressionOptions',
        putConfigurationSetSuppressionOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutConfigurationSetSuppressionOptionsRequest.fromBuffer(value),
        ($0.PutConfigurationSetSuppressionOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutConfigurationSetTrackingOptionsRequest,
            $0.PutConfigurationSetTrackingOptionsResponse>(
        'PutConfigurationSetTrackingOptions',
        putConfigurationSetTrackingOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutConfigurationSetTrackingOptionsRequest.fromBuffer(value),
        ($0.PutConfigurationSetTrackingOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutConfigurationSetVdmOptionsRequest,
            $0.PutConfigurationSetVdmOptionsResponse>(
        'PutConfigurationSetVdmOptions',
        putConfigurationSetVdmOptions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutConfigurationSetVdmOptionsRequest.fromBuffer(value),
        ($0.PutConfigurationSetVdmOptionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDedicatedIpInPoolRequest,
            $0.PutDedicatedIpInPoolResponse>(
        'PutDedicatedIpInPool',
        putDedicatedIpInPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDedicatedIpInPoolRequest.fromBuffer(value),
        ($0.PutDedicatedIpInPoolResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutDedicatedIpPoolScalingAttributesRequest,
            $0.PutDedicatedIpPoolScalingAttributesResponse>(
        'PutDedicatedIpPoolScalingAttributes',
        putDedicatedIpPoolScalingAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDedicatedIpPoolScalingAttributesRequest.fromBuffer(value),
        ($0.PutDedicatedIpPoolScalingAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDedicatedIpWarmupAttributesRequest,
            $0.PutDedicatedIpWarmupAttributesResponse>(
        'PutDedicatedIpWarmupAttributes',
        putDedicatedIpWarmupAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDedicatedIpWarmupAttributesRequest.fromBuffer(value),
        ($0.PutDedicatedIpWarmupAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutDeliverabilityDashboardOptionRequest,
            $0.PutDeliverabilityDashboardOptionResponse>(
        'PutDeliverabilityDashboardOption',
        putDeliverabilityDashboardOption_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutDeliverabilityDashboardOptionRequest.fromBuffer(value),
        ($0.PutDeliverabilityDashboardOptionResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutEmailIdentityConfigurationSetAttributesRequest,
            $0.PutEmailIdentityConfigurationSetAttributesResponse>(
        'PutEmailIdentityConfigurationSetAttributes',
        putEmailIdentityConfigurationSetAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutEmailIdentityConfigurationSetAttributesRequest.fromBuffer(
                value),
        ($0.PutEmailIdentityConfigurationSetAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutEmailIdentityDkimAttributesRequest,
            $0.PutEmailIdentityDkimAttributesResponse>(
        'PutEmailIdentityDkimAttributes',
        putEmailIdentityDkimAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutEmailIdentityDkimAttributesRequest.fromBuffer(value),
        ($0.PutEmailIdentityDkimAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.PutEmailIdentityDkimSigningAttributesRequest,
            $0.PutEmailIdentityDkimSigningAttributesResponse>(
        'PutEmailIdentityDkimSigningAttributes',
        putEmailIdentityDkimSigningAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutEmailIdentityDkimSigningAttributesRequest.fromBuffer(value),
        ($0.PutEmailIdentityDkimSigningAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutEmailIdentityFeedbackAttributesRequest,
            $0.PutEmailIdentityFeedbackAttributesResponse>(
        'PutEmailIdentityFeedbackAttributes',
        putEmailIdentityFeedbackAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutEmailIdentityFeedbackAttributesRequest.fromBuffer(value),
        ($0.PutEmailIdentityFeedbackAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutEmailIdentityMailFromAttributesRequest,
            $0.PutEmailIdentityMailFromAttributesResponse>(
        'PutEmailIdentityMailFromAttributes',
        putEmailIdentityMailFromAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutEmailIdentityMailFromAttributesRequest.fromBuffer(value),
        ($0.PutEmailIdentityMailFromAttributesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutSuppressedDestinationRequest,
            $0.PutSuppressedDestinationResponse>(
        'PutSuppressedDestination',
        putSuppressedDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutSuppressedDestinationRequest.fromBuffer(value),
        ($0.PutSuppressedDestinationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.SendBulkEmailRequest, $0.SendBulkEmailResponse>(
            'SendBulkEmail',
            sendBulkEmail_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.SendBulkEmailRequest.fromBuffer(value),
            ($0.SendBulkEmailResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendCustomVerificationEmailRequest,
            $0.SendCustomVerificationEmailResponse>(
        'SendCustomVerificationEmail',
        sendCustomVerificationEmail_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendCustomVerificationEmailRequest.fromBuffer(value),
        ($0.SendCustomVerificationEmailResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendEmailRequest, $0.SendEmailResponse>(
        'SendEmail',
        sendEmail_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SendEmailRequest.fromBuffer(value),
        ($0.SendEmailResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
            'TagResource',
            tagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TagResourceRequest.fromBuffer(value),
            ($0.TagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TestRenderEmailTemplateRequest,
            $0.TestRenderEmailTemplateResponse>(
        'TestRenderEmailTemplate',
        testRenderEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TestRenderEmailTemplateRequest.fromBuffer(value),
        ($0.TestRenderEmailTemplateResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceRequest.fromBuffer(value),
            ($0.UntagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateConfigurationSetEventDestinationRequest,
            $0.UpdateConfigurationSetEventDestinationResponse>(
        'UpdateConfigurationSetEventDestination',
        updateConfigurationSetEventDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateConfigurationSetEventDestinationRequest.fromBuffer(value),
        ($0.UpdateConfigurationSetEventDestinationResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateContactRequest, $0.UpdateContactResponse>(
            'UpdateContact',
            updateContact_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateContactRequest.fromBuffer(value),
            ($0.UpdateContactResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateContactListRequest,
            $0.UpdateContactListResponse>(
        'UpdateContactList',
        updateContactList_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateContactListRequest.fromBuffer(value),
        ($0.UpdateContactListResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateCustomVerificationEmailTemplateRequest,
            $0.UpdateCustomVerificationEmailTemplateResponse>(
        'UpdateCustomVerificationEmailTemplate',
        updateCustomVerificationEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateCustomVerificationEmailTemplateRequest.fromBuffer(value),
        ($0.UpdateCustomVerificationEmailTemplateResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateEmailIdentityPolicyRequest,
            $0.UpdateEmailIdentityPolicyResponse>(
        'UpdateEmailIdentityPolicy',
        updateEmailIdentityPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateEmailIdentityPolicyRequest.fromBuffer(value),
        ($0.UpdateEmailIdentityPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateEmailTemplateRequest,
            $0.UpdateEmailTemplateResponse>(
        'UpdateEmailTemplate',
        updateEmailTemplate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateEmailTemplateRequest.fromBuffer(value),
        ($0.UpdateEmailTemplateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateReputationEntityCustomerManagedStatusRequest,
            $0.UpdateReputationEntityCustomerManagedStatusResponse>(
        'UpdateReputationEntityCustomerManagedStatus',
        updateReputationEntityCustomerManagedStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateReputationEntityCustomerManagedStatusRequest.fromBuffer(
                value),
        ($0.UpdateReputationEntityCustomerManagedStatusResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateReputationEntityPolicyRequest,
            $0.UpdateReputationEntityPolicyResponse>(
        'UpdateReputationEntityPolicy',
        updateReputationEntityPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateReputationEntityPolicyRequest.fromBuffer(value),
        ($0.UpdateReputationEntityPolicyResponse value) =>
            value.writeToBuffer()));
  }

  $async.Future<$0.BatchGetMetricDataResponse> batchGetMetricData_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.BatchGetMetricDataRequest> $request) async {
    return batchGetMetricData($call, await $request);
  }

  $async.Future<$0.BatchGetMetricDataResponse> batchGetMetricData(
      $grpc.ServiceCall call, $0.BatchGetMetricDataRequest request);

  $async.Future<$0.CancelExportJobResponse> cancelExportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelExportJobRequest> $request) async {
    return cancelExportJob($call, await $request);
  }

  $async.Future<$0.CancelExportJobResponse> cancelExportJob(
      $grpc.ServiceCall call, $0.CancelExportJobRequest request);

  $async.Future<$0.CreateConfigurationSetResponse> createConfigurationSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateConfigurationSetRequest> $request) async {
    return createConfigurationSet($call, await $request);
  }

  $async.Future<$0.CreateConfigurationSetResponse> createConfigurationSet(
      $grpc.ServiceCall call, $0.CreateConfigurationSetRequest request);

  $async.Future<$0.CreateConfigurationSetEventDestinationResponse>
      createConfigurationSetEventDestination_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateConfigurationSetEventDestinationRequest>
              $request) async {
    return createConfigurationSetEventDestination($call, await $request);
  }

  $async.Future<$0.CreateConfigurationSetEventDestinationResponse>
      createConfigurationSetEventDestination($grpc.ServiceCall call,
          $0.CreateConfigurationSetEventDestinationRequest request);

  $async.Future<$0.CreateContactResponse> createContact_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateContactRequest> $request) async {
    return createContact($call, await $request);
  }

  $async.Future<$0.CreateContactResponse> createContact(
      $grpc.ServiceCall call, $0.CreateContactRequest request);

  $async.Future<$0.CreateContactListResponse> createContactList_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateContactListRequest> $request) async {
    return createContactList($call, await $request);
  }

  $async.Future<$0.CreateContactListResponse> createContactList(
      $grpc.ServiceCall call, $0.CreateContactListRequest request);

  $async.Future<$0.CreateCustomVerificationEmailTemplateResponse>
      createCustomVerificationEmailTemplate_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateCustomVerificationEmailTemplateRequest>
              $request) async {
    return createCustomVerificationEmailTemplate($call, await $request);
  }

  $async.Future<$0.CreateCustomVerificationEmailTemplateResponse>
      createCustomVerificationEmailTemplate($grpc.ServiceCall call,
          $0.CreateCustomVerificationEmailTemplateRequest request);

  $async.Future<$0.CreateDedicatedIpPoolResponse> createDedicatedIpPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDedicatedIpPoolRequest> $request) async {
    return createDedicatedIpPool($call, await $request);
  }

  $async.Future<$0.CreateDedicatedIpPoolResponse> createDedicatedIpPool(
      $grpc.ServiceCall call, $0.CreateDedicatedIpPoolRequest request);

  $async.Future<$0.CreateDeliverabilityTestReportResponse>
      createDeliverabilityTestReport_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateDeliverabilityTestReportRequest>
              $request) async {
    return createDeliverabilityTestReport($call, await $request);
  }

  $async.Future<$0.CreateDeliverabilityTestReportResponse>
      createDeliverabilityTestReport($grpc.ServiceCall call,
          $0.CreateDeliverabilityTestReportRequest request);

  $async.Future<$0.CreateEmailIdentityResponse> createEmailIdentity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateEmailIdentityRequest> $request) async {
    return createEmailIdentity($call, await $request);
  }

  $async.Future<$0.CreateEmailIdentityResponse> createEmailIdentity(
      $grpc.ServiceCall call, $0.CreateEmailIdentityRequest request);

  $async.Future<$0.CreateEmailIdentityPolicyResponse>
      createEmailIdentityPolicy_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateEmailIdentityPolicyRequest> $request) async {
    return createEmailIdentityPolicy($call, await $request);
  }

  $async.Future<$0.CreateEmailIdentityPolicyResponse> createEmailIdentityPolicy(
      $grpc.ServiceCall call, $0.CreateEmailIdentityPolicyRequest request);

  $async.Future<$0.CreateEmailTemplateResponse> createEmailTemplate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateEmailTemplateRequest> $request) async {
    return createEmailTemplate($call, await $request);
  }

  $async.Future<$0.CreateEmailTemplateResponse> createEmailTemplate(
      $grpc.ServiceCall call, $0.CreateEmailTemplateRequest request);

  $async.Future<$0.CreateExportJobResponse> createExportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateExportJobRequest> $request) async {
    return createExportJob($call, await $request);
  }

  $async.Future<$0.CreateExportJobResponse> createExportJob(
      $grpc.ServiceCall call, $0.CreateExportJobRequest request);

  $async.Future<$0.CreateImportJobResponse> createImportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateImportJobRequest> $request) async {
    return createImportJob($call, await $request);
  }

  $async.Future<$0.CreateImportJobResponse> createImportJob(
      $grpc.ServiceCall call, $0.CreateImportJobRequest request);

  $async.Future<$0.CreateMultiRegionEndpointResponse>
      createMultiRegionEndpoint_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateMultiRegionEndpointRequest> $request) async {
    return createMultiRegionEndpoint($call, await $request);
  }

  $async.Future<$0.CreateMultiRegionEndpointResponse> createMultiRegionEndpoint(
      $grpc.ServiceCall call, $0.CreateMultiRegionEndpointRequest request);

  $async.Future<$0.CreateTenantResponse> createTenant_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateTenantRequest> $request) async {
    return createTenant($call, await $request);
  }

  $async.Future<$0.CreateTenantResponse> createTenant(
      $grpc.ServiceCall call, $0.CreateTenantRequest request);

  $async.Future<$0.CreateTenantResourceAssociationResponse>
      createTenantResourceAssociation_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateTenantResourceAssociationRequest>
              $request) async {
    return createTenantResourceAssociation($call, await $request);
  }

  $async.Future<$0.CreateTenantResourceAssociationResponse>
      createTenantResourceAssociation($grpc.ServiceCall call,
          $0.CreateTenantResourceAssociationRequest request);

  $async.Future<$0.DeleteConfigurationSetResponse> deleteConfigurationSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteConfigurationSetRequest> $request) async {
    return deleteConfigurationSet($call, await $request);
  }

  $async.Future<$0.DeleteConfigurationSetResponse> deleteConfigurationSet(
      $grpc.ServiceCall call, $0.DeleteConfigurationSetRequest request);

  $async.Future<$0.DeleteConfigurationSetEventDestinationResponse>
      deleteConfigurationSetEventDestination_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeleteConfigurationSetEventDestinationRequest>
              $request) async {
    return deleteConfigurationSetEventDestination($call, await $request);
  }

  $async.Future<$0.DeleteConfigurationSetEventDestinationResponse>
      deleteConfigurationSetEventDestination($grpc.ServiceCall call,
          $0.DeleteConfigurationSetEventDestinationRequest request);

  $async.Future<$0.DeleteContactResponse> deleteContact_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteContactRequest> $request) async {
    return deleteContact($call, await $request);
  }

  $async.Future<$0.DeleteContactResponse> deleteContact(
      $grpc.ServiceCall call, $0.DeleteContactRequest request);

  $async.Future<$0.DeleteContactListResponse> deleteContactList_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteContactListRequest> $request) async {
    return deleteContactList($call, await $request);
  }

  $async.Future<$0.DeleteContactListResponse> deleteContactList(
      $grpc.ServiceCall call, $0.DeleteContactListRequest request);

  $async.Future<$0.DeleteCustomVerificationEmailTemplateResponse>
      deleteCustomVerificationEmailTemplate_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeleteCustomVerificationEmailTemplateRequest>
              $request) async {
    return deleteCustomVerificationEmailTemplate($call, await $request);
  }

  $async.Future<$0.DeleteCustomVerificationEmailTemplateResponse>
      deleteCustomVerificationEmailTemplate($grpc.ServiceCall call,
          $0.DeleteCustomVerificationEmailTemplateRequest request);

  $async.Future<$0.DeleteDedicatedIpPoolResponse> deleteDedicatedIpPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDedicatedIpPoolRequest> $request) async {
    return deleteDedicatedIpPool($call, await $request);
  }

  $async.Future<$0.DeleteDedicatedIpPoolResponse> deleteDedicatedIpPool(
      $grpc.ServiceCall call, $0.DeleteDedicatedIpPoolRequest request);

  $async.Future<$0.DeleteEmailIdentityResponse> deleteEmailIdentity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteEmailIdentityRequest> $request) async {
    return deleteEmailIdentity($call, await $request);
  }

  $async.Future<$0.DeleteEmailIdentityResponse> deleteEmailIdentity(
      $grpc.ServiceCall call, $0.DeleteEmailIdentityRequest request);

  $async.Future<$0.DeleteEmailIdentityPolicyResponse>
      deleteEmailIdentityPolicy_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteEmailIdentityPolicyRequest> $request) async {
    return deleteEmailIdentityPolicy($call, await $request);
  }

  $async.Future<$0.DeleteEmailIdentityPolicyResponse> deleteEmailIdentityPolicy(
      $grpc.ServiceCall call, $0.DeleteEmailIdentityPolicyRequest request);

  $async.Future<$0.DeleteEmailTemplateResponse> deleteEmailTemplate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteEmailTemplateRequest> $request) async {
    return deleteEmailTemplate($call, await $request);
  }

  $async.Future<$0.DeleteEmailTemplateResponse> deleteEmailTemplate(
      $grpc.ServiceCall call, $0.DeleteEmailTemplateRequest request);

  $async.Future<$0.DeleteMultiRegionEndpointResponse>
      deleteMultiRegionEndpoint_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteMultiRegionEndpointRequest> $request) async {
    return deleteMultiRegionEndpoint($call, await $request);
  }

  $async.Future<$0.DeleteMultiRegionEndpointResponse> deleteMultiRegionEndpoint(
      $grpc.ServiceCall call, $0.DeleteMultiRegionEndpointRequest request);

  $async.Future<$0.DeleteSuppressedDestinationResponse>
      deleteSuppressedDestination_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteSuppressedDestinationRequest> $request) async {
    return deleteSuppressedDestination($call, await $request);
  }

  $async.Future<$0.DeleteSuppressedDestinationResponse>
      deleteSuppressedDestination($grpc.ServiceCall call,
          $0.DeleteSuppressedDestinationRequest request);

  $async.Future<$0.DeleteTenantResponse> deleteTenant_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteTenantRequest> $request) async {
    return deleteTenant($call, await $request);
  }

  $async.Future<$0.DeleteTenantResponse> deleteTenant(
      $grpc.ServiceCall call, $0.DeleteTenantRequest request);

  $async.Future<$0.DeleteTenantResourceAssociationResponse>
      deleteTenantResourceAssociation_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeleteTenantResourceAssociationRequest>
              $request) async {
    return deleteTenantResourceAssociation($call, await $request);
  }

  $async.Future<$0.DeleteTenantResourceAssociationResponse>
      deleteTenantResourceAssociation($grpc.ServiceCall call,
          $0.DeleteTenantResourceAssociationRequest request);

  $async.Future<$0.GetAccountResponse> getAccount_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAccountRequest> $request) async {
    return getAccount($call, await $request);
  }

  $async.Future<$0.GetAccountResponse> getAccount(
      $grpc.ServiceCall call, $0.GetAccountRequest request);

  $async.Future<$0.GetBlacklistReportsResponse> getBlacklistReports_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetBlacklistReportsRequest> $request) async {
    return getBlacklistReports($call, await $request);
  }

  $async.Future<$0.GetBlacklistReportsResponse> getBlacklistReports(
      $grpc.ServiceCall call, $0.GetBlacklistReportsRequest request);

  $async.Future<$0.GetConfigurationSetResponse> getConfigurationSet_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetConfigurationSetRequest> $request) async {
    return getConfigurationSet($call, await $request);
  }

  $async.Future<$0.GetConfigurationSetResponse> getConfigurationSet(
      $grpc.ServiceCall call, $0.GetConfigurationSetRequest request);

  $async.Future<$0.GetConfigurationSetEventDestinationsResponse>
      getConfigurationSetEventDestinations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetConfigurationSetEventDestinationsRequest>
              $request) async {
    return getConfigurationSetEventDestinations($call, await $request);
  }

  $async.Future<$0.GetConfigurationSetEventDestinationsResponse>
      getConfigurationSetEventDestinations($grpc.ServiceCall call,
          $0.GetConfigurationSetEventDestinationsRequest request);

  $async.Future<$0.GetContactResponse> getContact_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetContactRequest> $request) async {
    return getContact($call, await $request);
  }

  $async.Future<$0.GetContactResponse> getContact(
      $grpc.ServiceCall call, $0.GetContactRequest request);

  $async.Future<$0.GetContactListResponse> getContactList_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetContactListRequest> $request) async {
    return getContactList($call, await $request);
  }

  $async.Future<$0.GetContactListResponse> getContactList(
      $grpc.ServiceCall call, $0.GetContactListRequest request);

  $async.Future<$0.GetCustomVerificationEmailTemplateResponse>
      getCustomVerificationEmailTemplate_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetCustomVerificationEmailTemplateRequest>
              $request) async {
    return getCustomVerificationEmailTemplate($call, await $request);
  }

  $async.Future<$0.GetCustomVerificationEmailTemplateResponse>
      getCustomVerificationEmailTemplate($grpc.ServiceCall call,
          $0.GetCustomVerificationEmailTemplateRequest request);

  $async.Future<$0.GetDedicatedIpResponse> getDedicatedIp_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDedicatedIpRequest> $request) async {
    return getDedicatedIp($call, await $request);
  }

  $async.Future<$0.GetDedicatedIpResponse> getDedicatedIp(
      $grpc.ServiceCall call, $0.GetDedicatedIpRequest request);

  $async.Future<$0.GetDedicatedIpPoolResponse> getDedicatedIpPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDedicatedIpPoolRequest> $request) async {
    return getDedicatedIpPool($call, await $request);
  }

  $async.Future<$0.GetDedicatedIpPoolResponse> getDedicatedIpPool(
      $grpc.ServiceCall call, $0.GetDedicatedIpPoolRequest request);

  $async.Future<$0.GetDedicatedIpsResponse> getDedicatedIps_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDedicatedIpsRequest> $request) async {
    return getDedicatedIps($call, await $request);
  }

  $async.Future<$0.GetDedicatedIpsResponse> getDedicatedIps(
      $grpc.ServiceCall call, $0.GetDedicatedIpsRequest request);

  $async.Future<$0.GetDeliverabilityDashboardOptionsResponse>
      getDeliverabilityDashboardOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetDeliverabilityDashboardOptionsRequest>
              $request) async {
    return getDeliverabilityDashboardOptions($call, await $request);
  }

  $async.Future<$0.GetDeliverabilityDashboardOptionsResponse>
      getDeliverabilityDashboardOptions($grpc.ServiceCall call,
          $0.GetDeliverabilityDashboardOptionsRequest request);

  $async.Future<$0.GetDeliverabilityTestReportResponse>
      getDeliverabilityTestReport_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetDeliverabilityTestReportRequest> $request) async {
    return getDeliverabilityTestReport($call, await $request);
  }

  $async.Future<$0.GetDeliverabilityTestReportResponse>
      getDeliverabilityTestReport($grpc.ServiceCall call,
          $0.GetDeliverabilityTestReportRequest request);

  $async.Future<$0.GetDomainDeliverabilityCampaignResponse>
      getDomainDeliverabilityCampaign_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetDomainDeliverabilityCampaignRequest>
              $request) async {
    return getDomainDeliverabilityCampaign($call, await $request);
  }

  $async.Future<$0.GetDomainDeliverabilityCampaignResponse>
      getDomainDeliverabilityCampaign($grpc.ServiceCall call,
          $0.GetDomainDeliverabilityCampaignRequest request);

  $async.Future<$0.GetDomainStatisticsReportResponse>
      getDomainStatisticsReport_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetDomainStatisticsReportRequest> $request) async {
    return getDomainStatisticsReport($call, await $request);
  }

  $async.Future<$0.GetDomainStatisticsReportResponse> getDomainStatisticsReport(
      $grpc.ServiceCall call, $0.GetDomainStatisticsReportRequest request);

  $async.Future<$0.GetEmailAddressInsightsResponse> getEmailAddressInsights_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEmailAddressInsightsRequest> $request) async {
    return getEmailAddressInsights($call, await $request);
  }

  $async.Future<$0.GetEmailAddressInsightsResponse> getEmailAddressInsights(
      $grpc.ServiceCall call, $0.GetEmailAddressInsightsRequest request);

  $async.Future<$0.GetEmailIdentityResponse> getEmailIdentity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEmailIdentityRequest> $request) async {
    return getEmailIdentity($call, await $request);
  }

  $async.Future<$0.GetEmailIdentityResponse> getEmailIdentity(
      $grpc.ServiceCall call, $0.GetEmailIdentityRequest request);

  $async.Future<$0.GetEmailIdentityPoliciesResponse>
      getEmailIdentityPolicies_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetEmailIdentityPoliciesRequest> $request) async {
    return getEmailIdentityPolicies($call, await $request);
  }

  $async.Future<$0.GetEmailIdentityPoliciesResponse> getEmailIdentityPolicies(
      $grpc.ServiceCall call, $0.GetEmailIdentityPoliciesRequest request);

  $async.Future<$0.GetEmailTemplateResponse> getEmailTemplate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEmailTemplateRequest> $request) async {
    return getEmailTemplate($call, await $request);
  }

  $async.Future<$0.GetEmailTemplateResponse> getEmailTemplate(
      $grpc.ServiceCall call, $0.GetEmailTemplateRequest request);

  $async.Future<$0.GetExportJobResponse> getExportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetExportJobRequest> $request) async {
    return getExportJob($call, await $request);
  }

  $async.Future<$0.GetExportJobResponse> getExportJob(
      $grpc.ServiceCall call, $0.GetExportJobRequest request);

  $async.Future<$0.GetImportJobResponse> getImportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetImportJobRequest> $request) async {
    return getImportJob($call, await $request);
  }

  $async.Future<$0.GetImportJobResponse> getImportJob(
      $grpc.ServiceCall call, $0.GetImportJobRequest request);

  $async.Future<$0.GetMessageInsightsResponse> getMessageInsights_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMessageInsightsRequest> $request) async {
    return getMessageInsights($call, await $request);
  }

  $async.Future<$0.GetMessageInsightsResponse> getMessageInsights(
      $grpc.ServiceCall call, $0.GetMessageInsightsRequest request);

  $async.Future<$0.GetMultiRegionEndpointResponse> getMultiRegionEndpoint_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMultiRegionEndpointRequest> $request) async {
    return getMultiRegionEndpoint($call, await $request);
  }

  $async.Future<$0.GetMultiRegionEndpointResponse> getMultiRegionEndpoint(
      $grpc.ServiceCall call, $0.GetMultiRegionEndpointRequest request);

  $async.Future<$0.GetReputationEntityResponse> getReputationEntity_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetReputationEntityRequest> $request) async {
    return getReputationEntity($call, await $request);
  }

  $async.Future<$0.GetReputationEntityResponse> getReputationEntity(
      $grpc.ServiceCall call, $0.GetReputationEntityRequest request);

  $async.Future<$0.GetSuppressedDestinationResponse>
      getSuppressedDestination_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetSuppressedDestinationRequest> $request) async {
    return getSuppressedDestination($call, await $request);
  }

  $async.Future<$0.GetSuppressedDestinationResponse> getSuppressedDestination(
      $grpc.ServiceCall call, $0.GetSuppressedDestinationRequest request);

  $async.Future<$0.GetTenantResponse> getTenant_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetTenantRequest> $request) async {
    return getTenant($call, await $request);
  }

  $async.Future<$0.GetTenantResponse> getTenant(
      $grpc.ServiceCall call, $0.GetTenantRequest request);

  $async.Future<$0.ListConfigurationSetsResponse> listConfigurationSets_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListConfigurationSetsRequest> $request) async {
    return listConfigurationSets($call, await $request);
  }

  $async.Future<$0.ListConfigurationSetsResponse> listConfigurationSets(
      $grpc.ServiceCall call, $0.ListConfigurationSetsRequest request);

  $async.Future<$0.ListContactListsResponse> listContactLists_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListContactListsRequest> $request) async {
    return listContactLists($call, await $request);
  }

  $async.Future<$0.ListContactListsResponse> listContactLists(
      $grpc.ServiceCall call, $0.ListContactListsRequest request);

  $async.Future<$0.ListContactsResponse> listContacts_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListContactsRequest> $request) async {
    return listContacts($call, await $request);
  }

  $async.Future<$0.ListContactsResponse> listContacts(
      $grpc.ServiceCall call, $0.ListContactsRequest request);

  $async.Future<$0.ListCustomVerificationEmailTemplatesResponse>
      listCustomVerificationEmailTemplates_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListCustomVerificationEmailTemplatesRequest>
              $request) async {
    return listCustomVerificationEmailTemplates($call, await $request);
  }

  $async.Future<$0.ListCustomVerificationEmailTemplatesResponse>
      listCustomVerificationEmailTemplates($grpc.ServiceCall call,
          $0.ListCustomVerificationEmailTemplatesRequest request);

  $async.Future<$0.ListDedicatedIpPoolsResponse> listDedicatedIpPools_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDedicatedIpPoolsRequest> $request) async {
    return listDedicatedIpPools($call, await $request);
  }

  $async.Future<$0.ListDedicatedIpPoolsResponse> listDedicatedIpPools(
      $grpc.ServiceCall call, $0.ListDedicatedIpPoolsRequest request);

  $async.Future<$0.ListDeliverabilityTestReportsResponse>
      listDeliverabilityTestReports_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDeliverabilityTestReportsRequest>
              $request) async {
    return listDeliverabilityTestReports($call, await $request);
  }

  $async.Future<$0.ListDeliverabilityTestReportsResponse>
      listDeliverabilityTestReports($grpc.ServiceCall call,
          $0.ListDeliverabilityTestReportsRequest request);

  $async.Future<$0.ListDomainDeliverabilityCampaignsResponse>
      listDomainDeliverabilityCampaigns_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDomainDeliverabilityCampaignsRequest>
              $request) async {
    return listDomainDeliverabilityCampaigns($call, await $request);
  }

  $async.Future<$0.ListDomainDeliverabilityCampaignsResponse>
      listDomainDeliverabilityCampaigns($grpc.ServiceCall call,
          $0.ListDomainDeliverabilityCampaignsRequest request);

  $async.Future<$0.ListEmailIdentitiesResponse> listEmailIdentities_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEmailIdentitiesRequest> $request) async {
    return listEmailIdentities($call, await $request);
  }

  $async.Future<$0.ListEmailIdentitiesResponse> listEmailIdentities(
      $grpc.ServiceCall call, $0.ListEmailIdentitiesRequest request);

  $async.Future<$0.ListEmailTemplatesResponse> listEmailTemplates_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEmailTemplatesRequest> $request) async {
    return listEmailTemplates($call, await $request);
  }

  $async.Future<$0.ListEmailTemplatesResponse> listEmailTemplates(
      $grpc.ServiceCall call, $0.ListEmailTemplatesRequest request);

  $async.Future<$0.ListExportJobsResponse> listExportJobs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListExportJobsRequest> $request) async {
    return listExportJobs($call, await $request);
  }

  $async.Future<$0.ListExportJobsResponse> listExportJobs(
      $grpc.ServiceCall call, $0.ListExportJobsRequest request);

  $async.Future<$0.ListImportJobsResponse> listImportJobs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListImportJobsRequest> $request) async {
    return listImportJobs($call, await $request);
  }

  $async.Future<$0.ListImportJobsResponse> listImportJobs(
      $grpc.ServiceCall call, $0.ListImportJobsRequest request);

  $async.Future<$0.ListMultiRegionEndpointsResponse>
      listMultiRegionEndpoints_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListMultiRegionEndpointsRequest> $request) async {
    return listMultiRegionEndpoints($call, await $request);
  }

  $async.Future<$0.ListMultiRegionEndpointsResponse> listMultiRegionEndpoints(
      $grpc.ServiceCall call, $0.ListMultiRegionEndpointsRequest request);

  $async.Future<$0.ListRecommendationsResponse> listRecommendations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRecommendationsRequest> $request) async {
    return listRecommendations($call, await $request);
  }

  $async.Future<$0.ListRecommendationsResponse> listRecommendations(
      $grpc.ServiceCall call, $0.ListRecommendationsRequest request);

  $async.Future<$0.ListReputationEntitiesResponse> listReputationEntities_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListReputationEntitiesRequest> $request) async {
    return listReputationEntities($call, await $request);
  }

  $async.Future<$0.ListReputationEntitiesResponse> listReputationEntities(
      $grpc.ServiceCall call, $0.ListReputationEntitiesRequest request);

  $async.Future<$0.ListResourceTenantsResponse> listResourceTenants_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListResourceTenantsRequest> $request) async {
    return listResourceTenants($call, await $request);
  }

  $async.Future<$0.ListResourceTenantsResponse> listResourceTenants(
      $grpc.ServiceCall call, $0.ListResourceTenantsRequest request);

  $async.Future<$0.ListSuppressedDestinationsResponse>
      listSuppressedDestinations_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListSuppressedDestinationsRequest> $request) async {
    return listSuppressedDestinations($call, await $request);
  }

  $async.Future<$0.ListSuppressedDestinationsResponse>
      listSuppressedDestinations(
          $grpc.ServiceCall call, $0.ListSuppressedDestinationsRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListTenantResourcesResponse> listTenantResources_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTenantResourcesRequest> $request) async {
    return listTenantResources($call, await $request);
  }

  $async.Future<$0.ListTenantResourcesResponse> listTenantResources(
      $grpc.ServiceCall call, $0.ListTenantResourcesRequest request);

  $async.Future<$0.ListTenantsResponse> listTenants_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTenantsRequest> $request) async {
    return listTenants($call, await $request);
  }

  $async.Future<$0.ListTenantsResponse> listTenants(
      $grpc.ServiceCall call, $0.ListTenantsRequest request);

  $async.Future<$0.PutAccountDedicatedIpWarmupAttributesResponse>
      putAccountDedicatedIpWarmupAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutAccountDedicatedIpWarmupAttributesRequest>
              $request) async {
    return putAccountDedicatedIpWarmupAttributes($call, await $request);
  }

  $async.Future<$0.PutAccountDedicatedIpWarmupAttributesResponse>
      putAccountDedicatedIpWarmupAttributes($grpc.ServiceCall call,
          $0.PutAccountDedicatedIpWarmupAttributesRequest request);

  $async.Future<$0.PutAccountDetailsResponse> putAccountDetails_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutAccountDetailsRequest> $request) async {
    return putAccountDetails($call, await $request);
  }

  $async.Future<$0.PutAccountDetailsResponse> putAccountDetails(
      $grpc.ServiceCall call, $0.PutAccountDetailsRequest request);

  $async.Future<$0.PutAccountSendingAttributesResponse>
      putAccountSendingAttributes_Pre($grpc.ServiceCall $call,
          $async.Future<$0.PutAccountSendingAttributesRequest> $request) async {
    return putAccountSendingAttributes($call, await $request);
  }

  $async.Future<$0.PutAccountSendingAttributesResponse>
      putAccountSendingAttributes($grpc.ServiceCall call,
          $0.PutAccountSendingAttributesRequest request);

  $async.Future<$0.PutAccountSuppressionAttributesResponse>
      putAccountSuppressionAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutAccountSuppressionAttributesRequest>
              $request) async {
    return putAccountSuppressionAttributes($call, await $request);
  }

  $async.Future<$0.PutAccountSuppressionAttributesResponse>
      putAccountSuppressionAttributes($grpc.ServiceCall call,
          $0.PutAccountSuppressionAttributesRequest request);

  $async.Future<$0.PutAccountVdmAttributesResponse> putAccountVdmAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutAccountVdmAttributesRequest> $request) async {
    return putAccountVdmAttributes($call, await $request);
  }

  $async.Future<$0.PutAccountVdmAttributesResponse> putAccountVdmAttributes(
      $grpc.ServiceCall call, $0.PutAccountVdmAttributesRequest request);

  $async.Future<$0.PutConfigurationSetArchivingOptionsResponse>
      putConfigurationSetArchivingOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutConfigurationSetArchivingOptionsRequest>
              $request) async {
    return putConfigurationSetArchivingOptions($call, await $request);
  }

  $async.Future<$0.PutConfigurationSetArchivingOptionsResponse>
      putConfigurationSetArchivingOptions($grpc.ServiceCall call,
          $0.PutConfigurationSetArchivingOptionsRequest request);

  $async.Future<$0.PutConfigurationSetDeliveryOptionsResponse>
      putConfigurationSetDeliveryOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutConfigurationSetDeliveryOptionsRequest>
              $request) async {
    return putConfigurationSetDeliveryOptions($call, await $request);
  }

  $async.Future<$0.PutConfigurationSetDeliveryOptionsResponse>
      putConfigurationSetDeliveryOptions($grpc.ServiceCall call,
          $0.PutConfigurationSetDeliveryOptionsRequest request);

  $async.Future<$0.PutConfigurationSetReputationOptionsResponse>
      putConfigurationSetReputationOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutConfigurationSetReputationOptionsRequest>
              $request) async {
    return putConfigurationSetReputationOptions($call, await $request);
  }

  $async.Future<$0.PutConfigurationSetReputationOptionsResponse>
      putConfigurationSetReputationOptions($grpc.ServiceCall call,
          $0.PutConfigurationSetReputationOptionsRequest request);

  $async.Future<$0.PutConfigurationSetSendingOptionsResponse>
      putConfigurationSetSendingOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutConfigurationSetSendingOptionsRequest>
              $request) async {
    return putConfigurationSetSendingOptions($call, await $request);
  }

  $async.Future<$0.PutConfigurationSetSendingOptionsResponse>
      putConfigurationSetSendingOptions($grpc.ServiceCall call,
          $0.PutConfigurationSetSendingOptionsRequest request);

  $async.Future<$0.PutConfigurationSetSuppressionOptionsResponse>
      putConfigurationSetSuppressionOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutConfigurationSetSuppressionOptionsRequest>
              $request) async {
    return putConfigurationSetSuppressionOptions($call, await $request);
  }

  $async.Future<$0.PutConfigurationSetSuppressionOptionsResponse>
      putConfigurationSetSuppressionOptions($grpc.ServiceCall call,
          $0.PutConfigurationSetSuppressionOptionsRequest request);

  $async.Future<$0.PutConfigurationSetTrackingOptionsResponse>
      putConfigurationSetTrackingOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutConfigurationSetTrackingOptionsRequest>
              $request) async {
    return putConfigurationSetTrackingOptions($call, await $request);
  }

  $async.Future<$0.PutConfigurationSetTrackingOptionsResponse>
      putConfigurationSetTrackingOptions($grpc.ServiceCall call,
          $0.PutConfigurationSetTrackingOptionsRequest request);

  $async.Future<$0.PutConfigurationSetVdmOptionsResponse>
      putConfigurationSetVdmOptions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutConfigurationSetVdmOptionsRequest>
              $request) async {
    return putConfigurationSetVdmOptions($call, await $request);
  }

  $async.Future<$0.PutConfigurationSetVdmOptionsResponse>
      putConfigurationSetVdmOptions($grpc.ServiceCall call,
          $0.PutConfigurationSetVdmOptionsRequest request);

  $async.Future<$0.PutDedicatedIpInPoolResponse> putDedicatedIpInPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutDedicatedIpInPoolRequest> $request) async {
    return putDedicatedIpInPool($call, await $request);
  }

  $async.Future<$0.PutDedicatedIpInPoolResponse> putDedicatedIpInPool(
      $grpc.ServiceCall call, $0.PutDedicatedIpInPoolRequest request);

  $async.Future<$0.PutDedicatedIpPoolScalingAttributesResponse>
      putDedicatedIpPoolScalingAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutDedicatedIpPoolScalingAttributesRequest>
              $request) async {
    return putDedicatedIpPoolScalingAttributes($call, await $request);
  }

  $async.Future<$0.PutDedicatedIpPoolScalingAttributesResponse>
      putDedicatedIpPoolScalingAttributes($grpc.ServiceCall call,
          $0.PutDedicatedIpPoolScalingAttributesRequest request);

  $async.Future<$0.PutDedicatedIpWarmupAttributesResponse>
      putDedicatedIpWarmupAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutDedicatedIpWarmupAttributesRequest>
              $request) async {
    return putDedicatedIpWarmupAttributes($call, await $request);
  }

  $async.Future<$0.PutDedicatedIpWarmupAttributesResponse>
      putDedicatedIpWarmupAttributes($grpc.ServiceCall call,
          $0.PutDedicatedIpWarmupAttributesRequest request);

  $async.Future<$0.PutDeliverabilityDashboardOptionResponse>
      putDeliverabilityDashboardOption_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutDeliverabilityDashboardOptionRequest>
              $request) async {
    return putDeliverabilityDashboardOption($call, await $request);
  }

  $async.Future<$0.PutDeliverabilityDashboardOptionResponse>
      putDeliverabilityDashboardOption($grpc.ServiceCall call,
          $0.PutDeliverabilityDashboardOptionRequest request);

  $async.Future<$0.PutEmailIdentityConfigurationSetAttributesResponse>
      putEmailIdentityConfigurationSetAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutEmailIdentityConfigurationSetAttributesRequest>
              $request) async {
    return putEmailIdentityConfigurationSetAttributes($call, await $request);
  }

  $async.Future<$0.PutEmailIdentityConfigurationSetAttributesResponse>
      putEmailIdentityConfigurationSetAttributes($grpc.ServiceCall call,
          $0.PutEmailIdentityConfigurationSetAttributesRequest request);

  $async.Future<$0.PutEmailIdentityDkimAttributesResponse>
      putEmailIdentityDkimAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutEmailIdentityDkimAttributesRequest>
              $request) async {
    return putEmailIdentityDkimAttributes($call, await $request);
  }

  $async.Future<$0.PutEmailIdentityDkimAttributesResponse>
      putEmailIdentityDkimAttributes($grpc.ServiceCall call,
          $0.PutEmailIdentityDkimAttributesRequest request);

  $async.Future<$0.PutEmailIdentityDkimSigningAttributesResponse>
      putEmailIdentityDkimSigningAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutEmailIdentityDkimSigningAttributesRequest>
              $request) async {
    return putEmailIdentityDkimSigningAttributes($call, await $request);
  }

  $async.Future<$0.PutEmailIdentityDkimSigningAttributesResponse>
      putEmailIdentityDkimSigningAttributes($grpc.ServiceCall call,
          $0.PutEmailIdentityDkimSigningAttributesRequest request);

  $async.Future<$0.PutEmailIdentityFeedbackAttributesResponse>
      putEmailIdentityFeedbackAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutEmailIdentityFeedbackAttributesRequest>
              $request) async {
    return putEmailIdentityFeedbackAttributes($call, await $request);
  }

  $async.Future<$0.PutEmailIdentityFeedbackAttributesResponse>
      putEmailIdentityFeedbackAttributes($grpc.ServiceCall call,
          $0.PutEmailIdentityFeedbackAttributesRequest request);

  $async.Future<$0.PutEmailIdentityMailFromAttributesResponse>
      putEmailIdentityMailFromAttributes_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutEmailIdentityMailFromAttributesRequest>
              $request) async {
    return putEmailIdentityMailFromAttributes($call, await $request);
  }

  $async.Future<$0.PutEmailIdentityMailFromAttributesResponse>
      putEmailIdentityMailFromAttributes($grpc.ServiceCall call,
          $0.PutEmailIdentityMailFromAttributesRequest request);

  $async.Future<$0.PutSuppressedDestinationResponse>
      putSuppressedDestination_Pre($grpc.ServiceCall $call,
          $async.Future<$0.PutSuppressedDestinationRequest> $request) async {
    return putSuppressedDestination($call, await $request);
  }

  $async.Future<$0.PutSuppressedDestinationResponse> putSuppressedDestination(
      $grpc.ServiceCall call, $0.PutSuppressedDestinationRequest request);

  $async.Future<$0.SendBulkEmailResponse> sendBulkEmail_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendBulkEmailRequest> $request) async {
    return sendBulkEmail($call, await $request);
  }

  $async.Future<$0.SendBulkEmailResponse> sendBulkEmail(
      $grpc.ServiceCall call, $0.SendBulkEmailRequest request);

  $async.Future<$0.SendCustomVerificationEmailResponse>
      sendCustomVerificationEmail_Pre($grpc.ServiceCall $call,
          $async.Future<$0.SendCustomVerificationEmailRequest> $request) async {
    return sendCustomVerificationEmail($call, await $request);
  }

  $async.Future<$0.SendCustomVerificationEmailResponse>
      sendCustomVerificationEmail($grpc.ServiceCall call,
          $0.SendCustomVerificationEmailRequest request);

  $async.Future<$0.SendEmailResponse> sendEmail_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendEmailRequest> $request) async {
    return sendEmail($call, await $request);
  }

  $async.Future<$0.SendEmailResponse> sendEmail(
      $grpc.ServiceCall call, $0.SendEmailRequest request);

  $async.Future<$0.TagResourceResponse> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceResponse> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.TestRenderEmailTemplateResponse> testRenderEmailTemplate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TestRenderEmailTemplateRequest> $request) async {
    return testRenderEmailTemplate($call, await $request);
  }

  $async.Future<$0.TestRenderEmailTemplateResponse> testRenderEmailTemplate(
      $grpc.ServiceCall call, $0.TestRenderEmailTemplateRequest request);

  $async.Future<$0.UntagResourceResponse> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceResponse> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.UpdateConfigurationSetEventDestinationResponse>
      updateConfigurationSetEventDestination_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateConfigurationSetEventDestinationRequest>
              $request) async {
    return updateConfigurationSetEventDestination($call, await $request);
  }

  $async.Future<$0.UpdateConfigurationSetEventDestinationResponse>
      updateConfigurationSetEventDestination($grpc.ServiceCall call,
          $0.UpdateConfigurationSetEventDestinationRequest request);

  $async.Future<$0.UpdateContactResponse> updateContact_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateContactRequest> $request) async {
    return updateContact($call, await $request);
  }

  $async.Future<$0.UpdateContactResponse> updateContact(
      $grpc.ServiceCall call, $0.UpdateContactRequest request);

  $async.Future<$0.UpdateContactListResponse> updateContactList_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateContactListRequest> $request) async {
    return updateContactList($call, await $request);
  }

  $async.Future<$0.UpdateContactListResponse> updateContactList(
      $grpc.ServiceCall call, $0.UpdateContactListRequest request);

  $async.Future<$0.UpdateCustomVerificationEmailTemplateResponse>
      updateCustomVerificationEmailTemplate_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateCustomVerificationEmailTemplateRequest>
              $request) async {
    return updateCustomVerificationEmailTemplate($call, await $request);
  }

  $async.Future<$0.UpdateCustomVerificationEmailTemplateResponse>
      updateCustomVerificationEmailTemplate($grpc.ServiceCall call,
          $0.UpdateCustomVerificationEmailTemplateRequest request);

  $async.Future<$0.UpdateEmailIdentityPolicyResponse>
      updateEmailIdentityPolicy_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateEmailIdentityPolicyRequest> $request) async {
    return updateEmailIdentityPolicy($call, await $request);
  }

  $async.Future<$0.UpdateEmailIdentityPolicyResponse> updateEmailIdentityPolicy(
      $grpc.ServiceCall call, $0.UpdateEmailIdentityPolicyRequest request);

  $async.Future<$0.UpdateEmailTemplateResponse> updateEmailTemplate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateEmailTemplateRequest> $request) async {
    return updateEmailTemplate($call, await $request);
  }

  $async.Future<$0.UpdateEmailTemplateResponse> updateEmailTemplate(
      $grpc.ServiceCall call, $0.UpdateEmailTemplateRequest request);

  $async.Future<$0.UpdateReputationEntityCustomerManagedStatusResponse>
      updateReputationEntityCustomerManagedStatus_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateReputationEntityCustomerManagedStatusRequest>
              $request) async {
    return updateReputationEntityCustomerManagedStatus($call, await $request);
  }

  $async.Future<$0.UpdateReputationEntityCustomerManagedStatusResponse>
      updateReputationEntityCustomerManagedStatus($grpc.ServiceCall call,
          $0.UpdateReputationEntityCustomerManagedStatusRequest request);

  $async.Future<$0.UpdateReputationEntityPolicyResponse>
      updateReputationEntityPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateReputationEntityPolicyRequest>
              $request) async {
    return updateReputationEntityPolicy($call, await $request);
  }

  $async.Future<$0.UpdateReputationEntityPolicyResponse>
      updateReputationEntityPolicy($grpc.ServiceCall call,
          $0.UpdateReputationEntityPolicyRequest request);
}
