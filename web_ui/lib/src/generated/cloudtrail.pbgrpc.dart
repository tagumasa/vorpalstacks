// This is a generated file - do not edit.
//
// Generated from cloudtrail.proto.

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

import 'cloudtrail.pb.dart' as $0;

export 'cloudtrail.pb.dart';

/// CloudTrailService provides cloudtrail API operations.
@$pb.GrpcServiceName('cloudtrail.CloudTrailService')
class CloudTrailServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CloudTrailServiceClient(super.channel, {super.options, super.interceptors});

  /// Adds one or more tags to a trail, event data store, dashboard, or channel, up to a limit of 50. Overwrites an existing tag's value when a new value is specified for an existing tag key. Tag key nam...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AddTagsResponse> addTags(
    $0.AddTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addTags, request, options: options);
  }

  /// Cancels a query if the query is not in a terminated state, such as CANCELLED, FAILED, TIMED_OUT, or FINISHED. You must specify an ARN value for EventDataStore. The ID of the query that you want to ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelQueryResponse> cancelQuery(
    $0.CancelQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelQuery, request, options: options);
  }

  /// Creates a channel for CloudTrail to ingest events from a partner or external source. After you create a channel, a CloudTrail Lake event data store can log events from the partner or source that yo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateChannelResponse> createChannel(
    $0.CreateChannelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createChannel, request, options: options);
  }

  /// Creates a custom dashboard or the Highlights dashboard. Custom dashboards - Custom dashboards allow you to query events in any event data store type. You can add up to 10 widgets to a custom dashbo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateDashboardResponse> createDashboard(
    $0.CreateDashboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDashboard, request, options: options);
  }

  /// Creates a new event data store.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateEventDataStoreResponse> createEventDataStore(
    $0.CreateEventDataStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createEventDataStore, request, options: options);
  }

  /// Creates a trail that specifies the settings for delivery of log data to an Amazon S3 bucket.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateTrailResponse> createTrail(
    $0.CreateTrailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTrail, request, options: options);
  }

  /// Deletes a channel.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteChannelResponse> deleteChannel(
    $0.DeleteChannelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteChannel, request, options: options);
  }

  /// Deletes the specified dashboard. You cannot delete a dashboard that has termination protection enabled.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteDashboardResponse> deleteDashboard(
    $0.DeleteDashboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDashboard, request, options: options);
  }

  /// Disables the event data store specified by EventDataStore, which accepts an event data store ARN. After you run DeleteEventDataStore, the event data store enters a PENDING_DELETION state, and is au...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteEventDataStoreResponse> deleteEventDataStore(
    $0.DeleteEventDataStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteEventDataStore, request, options: options);
  }

  /// Deletes the resource-based policy attached to the CloudTrail event data store, dashboard, or channel.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteResourcePolicyResponse> deleteResourcePolicy(
    $0.DeleteResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourcePolicy, request, options: options);
  }

  /// Deletes a trail. This operation must be called from the Region in which the trail was created. DeleteTrail cannot be called on the shadow trails (replicated trails in other Regions) of a trail that...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteTrailResponse> deleteTrail(
    $0.DeleteTrailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTrail, request, options: options);
  }

  /// Removes CloudTrail delegated administrator permissions from a member account in an organization.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeregisterOrganizationDelegatedAdminResponse>
      deregisterOrganizationDelegatedAdmin(
    $0.DeregisterOrganizationDelegatedAdminRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deregisterOrganizationDelegatedAdmin, request,
        options: options);
  }

  /// Returns metadata about a query, including query run time in milliseconds, number of events scanned and matched, and query status. If the query results were delivered to an S3 bucket, the response a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeQueryResponse> describeQuery(
    $0.DescribeQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeQuery, request, options: options);
  }

  /// Retrieves settings for one or more trails associated with the current Region for your account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeTrailsResponse> describeTrails(
    $0.DescribeTrailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeTrails, request, options: options);
  }

  /// Disables Lake query federation on the specified event data store. When you disable federation, CloudTrail disables the integration with Glue, Lake Formation, and Amazon Athena. After disabling Lake...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DisableFederationResponse> disableFederation(
    $0.DisableFederationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableFederation, request, options: options);
  }

  /// Enables Lake query federation on the specified event data store. Federating an event data store lets you view the metadata associated with the event data store in the Glue Data Catalog and run SQL ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.EnableFederationResponse> enableFederation(
    $0.EnableFederationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableFederation, request, options: options);
  }

  /// Generates a query from a natural language prompt. This operation uses generative artificial intelligence (generative AI) to produce a ready-to-use SQL query from the prompt. The prompt can be a que...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GenerateQueryResponse> generateQuery(
    $0.GenerateQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateQuery, request, options: options);
  }

  /// Returns information about a specific channel.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetChannelResponse> getChannel(
    $0.GetChannelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getChannel, request, options: options);
  }

  /// Returns the specified dashboard.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDashboardResponse> getDashboard(
    $0.GetDashboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDashboard, request, options: options);
  }

  /// Retrieves the current event configuration settings for the specified event data store or trail. The response includes maximum event size configuration, the context key selectors configured for the ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetEventConfigurationResponse> getEventConfiguration(
    $0.GetEventConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEventConfiguration, request, options: options);
  }

  /// Returns information about an event data store specified as either an ARN or the ID portion of the ARN.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetEventDataStoreResponse> getEventDataStore(
    $0.GetEventDataStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEventDataStore, request, options: options);
  }

  /// Describes the settings for the event selectors that you configured for your trail. The information returned for your event selectors includes the following: If your event selector includes read-onl...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetEventSelectorsResponse> getEventSelectors(
    $0.GetEventSelectorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEventSelectors, request, options: options);
  }

  /// Returns information about a specific import.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetImportResponse> getImport(
    $0.GetImportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getImport, request, options: options);
  }

  /// Describes the settings for the Insights event selectors that you configured for your trail or event data store. GetInsightSelectors shows if CloudTrail Insights logging is enabled and which Insight...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetInsightSelectorsResponse> getInsightSelectors(
    $0.GetInsightSelectorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInsightSelectors, request, options: options);
  }

  /// Gets event data results of a query. You must specify the QueryID value returned by the StartQuery operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetQueryResultsResponse> getQueryResults(
    $0.GetQueryResultsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueryResults, request, options: options);
  }

  /// Retrieves the JSON text of the resource-based policy document attached to the CloudTrail event data store, dashboard, or channel.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetResourcePolicyResponse> getResourcePolicy(
    $0.GetResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourcePolicy, request, options: options);
  }

  /// Returns settings information for a specified trail.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetTrailResponse> getTrail(
    $0.GetTrailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTrail, request, options: options);
  }

  /// Returns a JSON-formatted list of information about the specified trail. Fields include information on delivery errors, Amazon SNS and Amazon S3 errors, and start and stop logging times for each tra...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetTrailStatusResponse> getTrailStatus(
    $0.GetTrailStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTrailStatus, request, options: options);
  }

  /// Lists the channels in the current account, and their source names.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListChannelsResponse> listChannels(
    $0.ListChannelsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listChannels, request, options: options);
  }

  /// Returns information about all dashboards in the account, in the current Region.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListDashboardsResponse> listDashboards(
    $0.ListDashboardsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDashboards, request, options: options);
  }

  /// Returns information about all event data stores in the account, in the current Region.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListEventDataStoresResponse> listEventDataStores(
    $0.ListEventDataStoresRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEventDataStores, request, options: options);
  }

  /// Returns a list of failures for the specified import.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListImportFailuresResponse> listImportFailures(
    $0.ListImportFailuresRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listImportFailures, request, options: options);
  }

  /// Returns information on all imports, or a select set of imports by ImportStatus or Destination.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListImportsResponse> listImports(
    $0.ListImportsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listImports, request, options: options);
  }

  /// Returns Insights events generated on a trail that logs data events. You can list Insights events that occurred in a Region within the last 90 days. ListInsightsData supports the following Dimension...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListInsightsDataResponse> listInsightsData(
    $0.ListInsightsDataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInsightsData, request, options: options);
  }

  /// Returns Insights metrics data for trails that have enabled Insights. The request must include the EventSource, EventName, and InsightType parameters. If the InsightType is set to ApiErrorRateInsigh...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListInsightsMetricDataResponse>
      listInsightsMetricData(
    $0.ListInsightsMetricDataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInsightsMetricData, request,
        options: options);
  }

  /// Returns all public keys whose private keys were used to sign the digest files within the specified time range. The public key is needed to validate digest files that were signed with its correspond...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListPublicKeysResponse> listPublicKeys(
    $0.ListPublicKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPublicKeys, request, options: options);
  }

  /// Returns a list of queries and query statuses for the past seven days. You must specify an ARN value for EventDataStore. Optionally, to shorten the list of results, you can specify a time range, for...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListQueriesResponse> listQueries(
    $0.ListQueriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listQueries, request, options: options);
  }

  /// Lists the tags for the specified trails, event data stores, dashboards, or channels in the current Region.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsResponse> listTags(
    $0.ListTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTags, request, options: options);
  }

  /// Lists trails that are in the current account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTrailsResponse> listTrails(
    $0.ListTrailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTrails, request, options: options);
  }

  /// Looks up management events or CloudTrail Insights events that are captured by CloudTrail. You can look up events that occurred in a Region within the last 90 days. LookupEvents returns recent Insig...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.LookupEventsResponse> lookupEvents(
    $0.LookupEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$lookupEvents, request, options: options);
  }

  /// Updates the event configuration settings for the specified event data store or trail. This operation supports updating the maximum event size, adding or modifying context key selectors for event da...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutEventConfigurationResponse> putEventConfiguration(
    $0.PutEventConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putEventConfiguration, request, options: options);
  }

  /// Configures event selectors (also referred to as basic event selectors) or advanced event selectors for your trail. You can use either AdvancedEventSelectors or EventSelectors, but not both. If you ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutEventSelectorsResponse> putEventSelectors(
    $0.PutEventSelectorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putEventSelectors, request, options: options);
  }

  /// Lets you enable Insights event logging on specific event categories by specifying the Insights selectors that you want to enable on an existing trail or event data store. You also use PutInsightSel...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutInsightSelectorsResponse> putInsightSelectors(
    $0.PutInsightSelectorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putInsightSelectors, request, options: options);
  }

  /// Attaches a resource-based permission policy to a CloudTrail event data store, dashboard, or channel. For more information about resource-based policies, see CloudTrail resource-based policy example...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutResourcePolicyResponse> putResourcePolicy(
    $0.PutResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putResourcePolicy, request, options: options);
  }

  /// Registers an organization’s member account as the CloudTrail delegated administrator.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RegisterOrganizationDelegatedAdminResponse>
      registerOrganizationDelegatedAdmin(
    $0.RegisterOrganizationDelegatedAdminRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerOrganizationDelegatedAdmin, request,
        options: options);
  }

  /// Removes the specified tags from a trail, event data store, dashboard, or channel.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RemoveTagsResponse> removeTags(
    $0.RemoveTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeTags, request, options: options);
  }

  /// Restores a deleted event data store specified by EventDataStore, which accepts an event data store ARN. You can only restore a deleted event data store within the seven-day wait period after deleti...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RestoreEventDataStoreResponse> restoreEventDataStore(
    $0.RestoreEventDataStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$restoreEventDataStore, request, options: options);
  }

  /// Searches sample queries and returns a list of sample queries that are sorted by relevance. To search for sample queries, provide a natural language SearchPhrase in English.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SearchSampleQueriesResponse> searchSampleQueries(
    $0.SearchSampleQueriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$searchSampleQueries, request, options: options);
  }

  /// Starts a refresh of the specified dashboard. Each time a dashboard is refreshed, CloudTrail runs queries to populate the dashboard's widgets. CloudTrail must be granted permissions to run the Start...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartDashboardRefreshResponse> startDashboardRefresh(
    $0.StartDashboardRefreshRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startDashboardRefresh, request, options: options);
  }

  /// Starts the ingestion of live events on an event data store specified as either an ARN or the ID portion of the ARN. To start ingestion, the event data store Status must be STOPPED_INGESTION and the...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartEventDataStoreIngestionResponse>
      startEventDataStoreIngestion(
    $0.StartEventDataStoreIngestionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startEventDataStoreIngestion, request,
        options: options);
  }

  /// Starts an import of logged trail events from a source S3 bucket to a destination event data store. By default, CloudTrail only imports events contained in the S3 bucket's CloudTrail prefix and the ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartImportResponse> startImport(
    $0.StartImportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startImport, request, options: options);
  }

  /// Starts the recording of Amazon Web Services API calls and log file delivery for a trail. For a trail that is enabled in all Regions, this operation must be called from the Region in which the trail...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartLoggingResponse> startLogging(
    $0.StartLoggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startLogging, request, options: options);
  }

  /// Starts a CloudTrail Lake query. Use the QueryStatement parameter to provide your SQL query, enclosed in single quotation marks. Use the optional DeliveryS3Uri parameter to deliver the query results...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartQueryResponse> startQuery(
    $0.StartQueryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startQuery, request, options: options);
  }

  /// Stops the ingestion of live events on an event data store specified as either an ARN or the ID portion of the ARN. To stop ingestion, the event data store Status must be ENABLED and the eventCatego...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopEventDataStoreIngestionResponse>
      stopEventDataStoreIngestion(
    $0.StopEventDataStoreIngestionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopEventDataStoreIngestion, request,
        options: options);
  }

  /// Stops a specified import.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopImportResponse> stopImport(
    $0.StopImportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopImport, request, options: options);
  }

  /// Suspends the recording of Amazon Web Services API calls and log file delivery for the specified trail. Under most circumstances, there is no need to use this action. You can update a trail without ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopLoggingResponse> stopLogging(
    $0.StopLoggingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopLogging, request, options: options);
  }

  /// Updates a channel specified by a required channel ARN or UUID.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateChannelResponse> updateChannel(
    $0.UpdateChannelRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateChannel, request, options: options);
  }

  /// Updates the specified dashboard. To set a refresh schedule, CloudTrail must be granted permissions to run the StartDashboardRefresh operation to refresh the dashboard on your behalf. To provide per...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateDashboardResponse> updateDashboard(
    $0.UpdateDashboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDashboard, request, options: options);
  }

  /// Updates an event data store. The required EventDataStore value is an ARN or the ID portion of the ARN. Other parameters are optional, but at least one optional parameter must be specified, or Cloud...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateEventDataStoreResponse> updateEventDataStore(
    $0.UpdateEventDataStoreRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateEventDataStore, request, options: options);
  }

  /// Updates trail settings that control what events you are logging, and how to handle log files. Changes to a trail do not require stopping the CloudTrail service. Use this action to designate an exis...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateTrailResponse> updateTrail(
    $0.UpdateTrailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTrail, request, options: options);
  }

  // method descriptors

  static final _$addTags =
      $grpc.ClientMethod<$0.AddTagsRequest, $0.AddTagsResponse>(
          '/cloudtrail.CloudTrailService/AddTags',
          ($0.AddTagsRequest value) => value.writeToBuffer(),
          $0.AddTagsResponse.fromBuffer);
  static final _$cancelQuery =
      $grpc.ClientMethod<$0.CancelQueryRequest, $0.CancelQueryResponse>(
          '/cloudtrail.CloudTrailService/CancelQuery',
          ($0.CancelQueryRequest value) => value.writeToBuffer(),
          $0.CancelQueryResponse.fromBuffer);
  static final _$createChannel =
      $grpc.ClientMethod<$0.CreateChannelRequest, $0.CreateChannelResponse>(
          '/cloudtrail.CloudTrailService/CreateChannel',
          ($0.CreateChannelRequest value) => value.writeToBuffer(),
          $0.CreateChannelResponse.fromBuffer);
  static final _$createDashboard =
      $grpc.ClientMethod<$0.CreateDashboardRequest, $0.CreateDashboardResponse>(
          '/cloudtrail.CloudTrailService/CreateDashboard',
          ($0.CreateDashboardRequest value) => value.writeToBuffer(),
          $0.CreateDashboardResponse.fromBuffer);
  static final _$createEventDataStore = $grpc.ClientMethod<
          $0.CreateEventDataStoreRequest, $0.CreateEventDataStoreResponse>(
      '/cloudtrail.CloudTrailService/CreateEventDataStore',
      ($0.CreateEventDataStoreRequest value) => value.writeToBuffer(),
      $0.CreateEventDataStoreResponse.fromBuffer);
  static final _$createTrail =
      $grpc.ClientMethod<$0.CreateTrailRequest, $0.CreateTrailResponse>(
          '/cloudtrail.CloudTrailService/CreateTrail',
          ($0.CreateTrailRequest value) => value.writeToBuffer(),
          $0.CreateTrailResponse.fromBuffer);
  static final _$deleteChannel =
      $grpc.ClientMethod<$0.DeleteChannelRequest, $0.DeleteChannelResponse>(
          '/cloudtrail.CloudTrailService/DeleteChannel',
          ($0.DeleteChannelRequest value) => value.writeToBuffer(),
          $0.DeleteChannelResponse.fromBuffer);
  static final _$deleteDashboard =
      $grpc.ClientMethod<$0.DeleteDashboardRequest, $0.DeleteDashboardResponse>(
          '/cloudtrail.CloudTrailService/DeleteDashboard',
          ($0.DeleteDashboardRequest value) => value.writeToBuffer(),
          $0.DeleteDashboardResponse.fromBuffer);
  static final _$deleteEventDataStore = $grpc.ClientMethod<
          $0.DeleteEventDataStoreRequest, $0.DeleteEventDataStoreResponse>(
      '/cloudtrail.CloudTrailService/DeleteEventDataStore',
      ($0.DeleteEventDataStoreRequest value) => value.writeToBuffer(),
      $0.DeleteEventDataStoreResponse.fromBuffer);
  static final _$deleteResourcePolicy = $grpc.ClientMethod<
          $0.DeleteResourcePolicyRequest, $0.DeleteResourcePolicyResponse>(
      '/cloudtrail.CloudTrailService/DeleteResourcePolicy',
      ($0.DeleteResourcePolicyRequest value) => value.writeToBuffer(),
      $0.DeleteResourcePolicyResponse.fromBuffer);
  static final _$deleteTrail =
      $grpc.ClientMethod<$0.DeleteTrailRequest, $0.DeleteTrailResponse>(
          '/cloudtrail.CloudTrailService/DeleteTrail',
          ($0.DeleteTrailRequest value) => value.writeToBuffer(),
          $0.DeleteTrailResponse.fromBuffer);
  static final _$deregisterOrganizationDelegatedAdmin = $grpc.ClientMethod<
          $0.DeregisterOrganizationDelegatedAdminRequest,
          $0.DeregisterOrganizationDelegatedAdminResponse>(
      '/cloudtrail.CloudTrailService/DeregisterOrganizationDelegatedAdmin',
      ($0.DeregisterOrganizationDelegatedAdminRequest value) =>
          value.writeToBuffer(),
      $0.DeregisterOrganizationDelegatedAdminResponse.fromBuffer);
  static final _$describeQuery =
      $grpc.ClientMethod<$0.DescribeQueryRequest, $0.DescribeQueryResponse>(
          '/cloudtrail.CloudTrailService/DescribeQuery',
          ($0.DescribeQueryRequest value) => value.writeToBuffer(),
          $0.DescribeQueryResponse.fromBuffer);
  static final _$describeTrails =
      $grpc.ClientMethod<$0.DescribeTrailsRequest, $0.DescribeTrailsResponse>(
          '/cloudtrail.CloudTrailService/DescribeTrails',
          ($0.DescribeTrailsRequest value) => value.writeToBuffer(),
          $0.DescribeTrailsResponse.fromBuffer);
  static final _$disableFederation = $grpc.ClientMethod<
          $0.DisableFederationRequest, $0.DisableFederationResponse>(
      '/cloudtrail.CloudTrailService/DisableFederation',
      ($0.DisableFederationRequest value) => value.writeToBuffer(),
      $0.DisableFederationResponse.fromBuffer);
  static final _$enableFederation = $grpc.ClientMethod<
          $0.EnableFederationRequest, $0.EnableFederationResponse>(
      '/cloudtrail.CloudTrailService/EnableFederation',
      ($0.EnableFederationRequest value) => value.writeToBuffer(),
      $0.EnableFederationResponse.fromBuffer);
  static final _$generateQuery =
      $grpc.ClientMethod<$0.GenerateQueryRequest, $0.GenerateQueryResponse>(
          '/cloudtrail.CloudTrailService/GenerateQuery',
          ($0.GenerateQueryRequest value) => value.writeToBuffer(),
          $0.GenerateQueryResponse.fromBuffer);
  static final _$getChannel =
      $grpc.ClientMethod<$0.GetChannelRequest, $0.GetChannelResponse>(
          '/cloudtrail.CloudTrailService/GetChannel',
          ($0.GetChannelRequest value) => value.writeToBuffer(),
          $0.GetChannelResponse.fromBuffer);
  static final _$getDashboard =
      $grpc.ClientMethod<$0.GetDashboardRequest, $0.GetDashboardResponse>(
          '/cloudtrail.CloudTrailService/GetDashboard',
          ($0.GetDashboardRequest value) => value.writeToBuffer(),
          $0.GetDashboardResponse.fromBuffer);
  static final _$getEventConfiguration = $grpc.ClientMethod<
          $0.GetEventConfigurationRequest, $0.GetEventConfigurationResponse>(
      '/cloudtrail.CloudTrailService/GetEventConfiguration',
      ($0.GetEventConfigurationRequest value) => value.writeToBuffer(),
      $0.GetEventConfigurationResponse.fromBuffer);
  static final _$getEventDataStore = $grpc.ClientMethod<
          $0.GetEventDataStoreRequest, $0.GetEventDataStoreResponse>(
      '/cloudtrail.CloudTrailService/GetEventDataStore',
      ($0.GetEventDataStoreRequest value) => value.writeToBuffer(),
      $0.GetEventDataStoreResponse.fromBuffer);
  static final _$getEventSelectors = $grpc.ClientMethod<
          $0.GetEventSelectorsRequest, $0.GetEventSelectorsResponse>(
      '/cloudtrail.CloudTrailService/GetEventSelectors',
      ($0.GetEventSelectorsRequest value) => value.writeToBuffer(),
      $0.GetEventSelectorsResponse.fromBuffer);
  static final _$getImport =
      $grpc.ClientMethod<$0.GetImportRequest, $0.GetImportResponse>(
          '/cloudtrail.CloudTrailService/GetImport',
          ($0.GetImportRequest value) => value.writeToBuffer(),
          $0.GetImportResponse.fromBuffer);
  static final _$getInsightSelectors = $grpc.ClientMethod<
          $0.GetInsightSelectorsRequest, $0.GetInsightSelectorsResponse>(
      '/cloudtrail.CloudTrailService/GetInsightSelectors',
      ($0.GetInsightSelectorsRequest value) => value.writeToBuffer(),
      $0.GetInsightSelectorsResponse.fromBuffer);
  static final _$getQueryResults =
      $grpc.ClientMethod<$0.GetQueryResultsRequest, $0.GetQueryResultsResponse>(
          '/cloudtrail.CloudTrailService/GetQueryResults',
          ($0.GetQueryResultsRequest value) => value.writeToBuffer(),
          $0.GetQueryResultsResponse.fromBuffer);
  static final _$getResourcePolicy = $grpc.ClientMethod<
          $0.GetResourcePolicyRequest, $0.GetResourcePolicyResponse>(
      '/cloudtrail.CloudTrailService/GetResourcePolicy',
      ($0.GetResourcePolicyRequest value) => value.writeToBuffer(),
      $0.GetResourcePolicyResponse.fromBuffer);
  static final _$getTrail =
      $grpc.ClientMethod<$0.GetTrailRequest, $0.GetTrailResponse>(
          '/cloudtrail.CloudTrailService/GetTrail',
          ($0.GetTrailRequest value) => value.writeToBuffer(),
          $0.GetTrailResponse.fromBuffer);
  static final _$getTrailStatus =
      $grpc.ClientMethod<$0.GetTrailStatusRequest, $0.GetTrailStatusResponse>(
          '/cloudtrail.CloudTrailService/GetTrailStatus',
          ($0.GetTrailStatusRequest value) => value.writeToBuffer(),
          $0.GetTrailStatusResponse.fromBuffer);
  static final _$listChannels =
      $grpc.ClientMethod<$0.ListChannelsRequest, $0.ListChannelsResponse>(
          '/cloudtrail.CloudTrailService/ListChannels',
          ($0.ListChannelsRequest value) => value.writeToBuffer(),
          $0.ListChannelsResponse.fromBuffer);
  static final _$listDashboards =
      $grpc.ClientMethod<$0.ListDashboardsRequest, $0.ListDashboardsResponse>(
          '/cloudtrail.CloudTrailService/ListDashboards',
          ($0.ListDashboardsRequest value) => value.writeToBuffer(),
          $0.ListDashboardsResponse.fromBuffer);
  static final _$listEventDataStores = $grpc.ClientMethod<
          $0.ListEventDataStoresRequest, $0.ListEventDataStoresResponse>(
      '/cloudtrail.CloudTrailService/ListEventDataStores',
      ($0.ListEventDataStoresRequest value) => value.writeToBuffer(),
      $0.ListEventDataStoresResponse.fromBuffer);
  static final _$listImportFailures = $grpc.ClientMethod<
          $0.ListImportFailuresRequest, $0.ListImportFailuresResponse>(
      '/cloudtrail.CloudTrailService/ListImportFailures',
      ($0.ListImportFailuresRequest value) => value.writeToBuffer(),
      $0.ListImportFailuresResponse.fromBuffer);
  static final _$listImports =
      $grpc.ClientMethod<$0.ListImportsRequest, $0.ListImportsResponse>(
          '/cloudtrail.CloudTrailService/ListImports',
          ($0.ListImportsRequest value) => value.writeToBuffer(),
          $0.ListImportsResponse.fromBuffer);
  static final _$listInsightsData = $grpc.ClientMethod<
          $0.ListInsightsDataRequest, $0.ListInsightsDataResponse>(
      '/cloudtrail.CloudTrailService/ListInsightsData',
      ($0.ListInsightsDataRequest value) => value.writeToBuffer(),
      $0.ListInsightsDataResponse.fromBuffer);
  static final _$listInsightsMetricData = $grpc.ClientMethod<
          $0.ListInsightsMetricDataRequest, $0.ListInsightsMetricDataResponse>(
      '/cloudtrail.CloudTrailService/ListInsightsMetricData',
      ($0.ListInsightsMetricDataRequest value) => value.writeToBuffer(),
      $0.ListInsightsMetricDataResponse.fromBuffer);
  static final _$listPublicKeys =
      $grpc.ClientMethod<$0.ListPublicKeysRequest, $0.ListPublicKeysResponse>(
          '/cloudtrail.CloudTrailService/ListPublicKeys',
          ($0.ListPublicKeysRequest value) => value.writeToBuffer(),
          $0.ListPublicKeysResponse.fromBuffer);
  static final _$listQueries =
      $grpc.ClientMethod<$0.ListQueriesRequest, $0.ListQueriesResponse>(
          '/cloudtrail.CloudTrailService/ListQueries',
          ($0.ListQueriesRequest value) => value.writeToBuffer(),
          $0.ListQueriesResponse.fromBuffer);
  static final _$listTags =
      $grpc.ClientMethod<$0.ListTagsRequest, $0.ListTagsResponse>(
          '/cloudtrail.CloudTrailService/ListTags',
          ($0.ListTagsRequest value) => value.writeToBuffer(),
          $0.ListTagsResponse.fromBuffer);
  static final _$listTrails =
      $grpc.ClientMethod<$0.ListTrailsRequest, $0.ListTrailsResponse>(
          '/cloudtrail.CloudTrailService/ListTrails',
          ($0.ListTrailsRequest value) => value.writeToBuffer(),
          $0.ListTrailsResponse.fromBuffer);
  static final _$lookupEvents =
      $grpc.ClientMethod<$0.LookupEventsRequest, $0.LookupEventsResponse>(
          '/cloudtrail.CloudTrailService/LookupEvents',
          ($0.LookupEventsRequest value) => value.writeToBuffer(),
          $0.LookupEventsResponse.fromBuffer);
  static final _$putEventConfiguration = $grpc.ClientMethod<
          $0.PutEventConfigurationRequest, $0.PutEventConfigurationResponse>(
      '/cloudtrail.CloudTrailService/PutEventConfiguration',
      ($0.PutEventConfigurationRequest value) => value.writeToBuffer(),
      $0.PutEventConfigurationResponse.fromBuffer);
  static final _$putEventSelectors = $grpc.ClientMethod<
          $0.PutEventSelectorsRequest, $0.PutEventSelectorsResponse>(
      '/cloudtrail.CloudTrailService/PutEventSelectors',
      ($0.PutEventSelectorsRequest value) => value.writeToBuffer(),
      $0.PutEventSelectorsResponse.fromBuffer);
  static final _$putInsightSelectors = $grpc.ClientMethod<
          $0.PutInsightSelectorsRequest, $0.PutInsightSelectorsResponse>(
      '/cloudtrail.CloudTrailService/PutInsightSelectors',
      ($0.PutInsightSelectorsRequest value) => value.writeToBuffer(),
      $0.PutInsightSelectorsResponse.fromBuffer);
  static final _$putResourcePolicy = $grpc.ClientMethod<
          $0.PutResourcePolicyRequest, $0.PutResourcePolicyResponse>(
      '/cloudtrail.CloudTrailService/PutResourcePolicy',
      ($0.PutResourcePolicyRequest value) => value.writeToBuffer(),
      $0.PutResourcePolicyResponse.fromBuffer);
  static final _$registerOrganizationDelegatedAdmin = $grpc.ClientMethod<
          $0.RegisterOrganizationDelegatedAdminRequest,
          $0.RegisterOrganizationDelegatedAdminResponse>(
      '/cloudtrail.CloudTrailService/RegisterOrganizationDelegatedAdmin',
      ($0.RegisterOrganizationDelegatedAdminRequest value) =>
          value.writeToBuffer(),
      $0.RegisterOrganizationDelegatedAdminResponse.fromBuffer);
  static final _$removeTags =
      $grpc.ClientMethod<$0.RemoveTagsRequest, $0.RemoveTagsResponse>(
          '/cloudtrail.CloudTrailService/RemoveTags',
          ($0.RemoveTagsRequest value) => value.writeToBuffer(),
          $0.RemoveTagsResponse.fromBuffer);
  static final _$restoreEventDataStore = $grpc.ClientMethod<
          $0.RestoreEventDataStoreRequest, $0.RestoreEventDataStoreResponse>(
      '/cloudtrail.CloudTrailService/RestoreEventDataStore',
      ($0.RestoreEventDataStoreRequest value) => value.writeToBuffer(),
      $0.RestoreEventDataStoreResponse.fromBuffer);
  static final _$searchSampleQueries = $grpc.ClientMethod<
          $0.SearchSampleQueriesRequest, $0.SearchSampleQueriesResponse>(
      '/cloudtrail.CloudTrailService/SearchSampleQueries',
      ($0.SearchSampleQueriesRequest value) => value.writeToBuffer(),
      $0.SearchSampleQueriesResponse.fromBuffer);
  static final _$startDashboardRefresh = $grpc.ClientMethod<
          $0.StartDashboardRefreshRequest, $0.StartDashboardRefreshResponse>(
      '/cloudtrail.CloudTrailService/StartDashboardRefresh',
      ($0.StartDashboardRefreshRequest value) => value.writeToBuffer(),
      $0.StartDashboardRefreshResponse.fromBuffer);
  static final _$startEventDataStoreIngestion = $grpc.ClientMethod<
          $0.StartEventDataStoreIngestionRequest,
          $0.StartEventDataStoreIngestionResponse>(
      '/cloudtrail.CloudTrailService/StartEventDataStoreIngestion',
      ($0.StartEventDataStoreIngestionRequest value) => value.writeToBuffer(),
      $0.StartEventDataStoreIngestionResponse.fromBuffer);
  static final _$startImport =
      $grpc.ClientMethod<$0.StartImportRequest, $0.StartImportResponse>(
          '/cloudtrail.CloudTrailService/StartImport',
          ($0.StartImportRequest value) => value.writeToBuffer(),
          $0.StartImportResponse.fromBuffer);
  static final _$startLogging =
      $grpc.ClientMethod<$0.StartLoggingRequest, $0.StartLoggingResponse>(
          '/cloudtrail.CloudTrailService/StartLogging',
          ($0.StartLoggingRequest value) => value.writeToBuffer(),
          $0.StartLoggingResponse.fromBuffer);
  static final _$startQuery =
      $grpc.ClientMethod<$0.StartQueryRequest, $0.StartQueryResponse>(
          '/cloudtrail.CloudTrailService/StartQuery',
          ($0.StartQueryRequest value) => value.writeToBuffer(),
          $0.StartQueryResponse.fromBuffer);
  static final _$stopEventDataStoreIngestion = $grpc.ClientMethod<
          $0.StopEventDataStoreIngestionRequest,
          $0.StopEventDataStoreIngestionResponse>(
      '/cloudtrail.CloudTrailService/StopEventDataStoreIngestion',
      ($0.StopEventDataStoreIngestionRequest value) => value.writeToBuffer(),
      $0.StopEventDataStoreIngestionResponse.fromBuffer);
  static final _$stopImport =
      $grpc.ClientMethod<$0.StopImportRequest, $0.StopImportResponse>(
          '/cloudtrail.CloudTrailService/StopImport',
          ($0.StopImportRequest value) => value.writeToBuffer(),
          $0.StopImportResponse.fromBuffer);
  static final _$stopLogging =
      $grpc.ClientMethod<$0.StopLoggingRequest, $0.StopLoggingResponse>(
          '/cloudtrail.CloudTrailService/StopLogging',
          ($0.StopLoggingRequest value) => value.writeToBuffer(),
          $0.StopLoggingResponse.fromBuffer);
  static final _$updateChannel =
      $grpc.ClientMethod<$0.UpdateChannelRequest, $0.UpdateChannelResponse>(
          '/cloudtrail.CloudTrailService/UpdateChannel',
          ($0.UpdateChannelRequest value) => value.writeToBuffer(),
          $0.UpdateChannelResponse.fromBuffer);
  static final _$updateDashboard =
      $grpc.ClientMethod<$0.UpdateDashboardRequest, $0.UpdateDashboardResponse>(
          '/cloudtrail.CloudTrailService/UpdateDashboard',
          ($0.UpdateDashboardRequest value) => value.writeToBuffer(),
          $0.UpdateDashboardResponse.fromBuffer);
  static final _$updateEventDataStore = $grpc.ClientMethod<
          $0.UpdateEventDataStoreRequest, $0.UpdateEventDataStoreResponse>(
      '/cloudtrail.CloudTrailService/UpdateEventDataStore',
      ($0.UpdateEventDataStoreRequest value) => value.writeToBuffer(),
      $0.UpdateEventDataStoreResponse.fromBuffer);
  static final _$updateTrail =
      $grpc.ClientMethod<$0.UpdateTrailRequest, $0.UpdateTrailResponse>(
          '/cloudtrail.CloudTrailService/UpdateTrail',
          ($0.UpdateTrailRequest value) => value.writeToBuffer(),
          $0.UpdateTrailResponse.fromBuffer);
}

@$pb.GrpcServiceName('cloudtrail.CloudTrailService')
abstract class CloudTrailServiceBase extends $grpc.Service {
  $core.String get $name => 'cloudtrail.CloudTrailService';

  CloudTrailServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AddTagsRequest, $0.AddTagsResponse>(
        'AddTags',
        addTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.AddTagsRequest.fromBuffer(value),
        ($0.AddTagsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CancelQueryRequest, $0.CancelQueryResponse>(
            'CancelQuery',
            cancelQuery_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CancelQueryRequest.fromBuffer(value),
            ($0.CancelQueryResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateChannelRequest, $0.CreateChannelResponse>(
            'CreateChannel',
            createChannel_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateChannelRequest.fromBuffer(value),
            ($0.CreateChannelResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDashboardRequest,
            $0.CreateDashboardResponse>(
        'CreateDashboard',
        createDashboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDashboardRequest.fromBuffer(value),
        ($0.CreateDashboardResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateEventDataStoreRequest,
            $0.CreateEventDataStoreResponse>(
        'CreateEventDataStore',
        createEventDataStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateEventDataStoreRequest.fromBuffer(value),
        ($0.CreateEventDataStoreResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateTrailRequest, $0.CreateTrailResponse>(
            'CreateTrail',
            createTrail_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateTrailRequest.fromBuffer(value),
            ($0.CreateTrailResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteChannelRequest, $0.DeleteChannelResponse>(
            'DeleteChannel',
            deleteChannel_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteChannelRequest.fromBuffer(value),
            ($0.DeleteChannelResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDashboardRequest,
            $0.DeleteDashboardResponse>(
        'DeleteDashboard',
        deleteDashboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDashboardRequest.fromBuffer(value),
        ($0.DeleteDashboardResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteEventDataStoreRequest,
            $0.DeleteEventDataStoreResponse>(
        'DeleteEventDataStore',
        deleteEventDataStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteEventDataStoreRequest.fromBuffer(value),
        ($0.DeleteEventDataStoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourcePolicyRequest,
            $0.DeleteResourcePolicyResponse>(
        'DeleteResourcePolicy',
        deleteResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourcePolicyRequest.fromBuffer(value),
        ($0.DeleteResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteTrailRequest, $0.DeleteTrailResponse>(
            'DeleteTrail',
            deleteTrail_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteTrailRequest.fromBuffer(value),
            ($0.DeleteTrailResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeregisterOrganizationDelegatedAdminRequest,
            $0.DeregisterOrganizationDelegatedAdminResponse>(
        'DeregisterOrganizationDelegatedAdmin',
        deregisterOrganizationDelegatedAdmin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeregisterOrganizationDelegatedAdminRequest.fromBuffer(value),
        ($0.DeregisterOrganizationDelegatedAdminResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeQueryRequest, $0.DescribeQueryResponse>(
            'DescribeQuery',
            describeQuery_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeQueryRequest.fromBuffer(value),
            ($0.DescribeQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeTrailsRequest,
            $0.DescribeTrailsResponse>(
        'DescribeTrails',
        describeTrails_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeTrailsRequest.fromBuffer(value),
        ($0.DescribeTrailsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableFederationRequest,
            $0.DisableFederationResponse>(
        'DisableFederation',
        disableFederation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableFederationRequest.fromBuffer(value),
        ($0.DisableFederationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableFederationRequest,
            $0.EnableFederationResponse>(
        'EnableFederation',
        enableFederation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableFederationRequest.fromBuffer(value),
        ($0.EnableFederationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GenerateQueryRequest, $0.GenerateQueryResponse>(
            'GenerateQuery',
            generateQuery_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GenerateQueryRequest.fromBuffer(value),
            ($0.GenerateQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetChannelRequest, $0.GetChannelResponse>(
        'GetChannel',
        getChannel_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetChannelRequest.fromBuffer(value),
        ($0.GetChannelResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetDashboardRequest, $0.GetDashboardResponse>(
            'GetDashboard',
            getDashboard_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetDashboardRequest.fromBuffer(value),
            ($0.GetDashboardResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEventConfigurationRequest,
            $0.GetEventConfigurationResponse>(
        'GetEventConfiguration',
        getEventConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEventConfigurationRequest.fromBuffer(value),
        ($0.GetEventConfigurationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEventDataStoreRequest,
            $0.GetEventDataStoreResponse>(
        'GetEventDataStore',
        getEventDataStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEventDataStoreRequest.fromBuffer(value),
        ($0.GetEventDataStoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEventSelectorsRequest,
            $0.GetEventSelectorsResponse>(
        'GetEventSelectors',
        getEventSelectors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEventSelectorsRequest.fromBuffer(value),
        ($0.GetEventSelectorsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetImportRequest, $0.GetImportResponse>(
        'GetImport',
        getImport_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetImportRequest.fromBuffer(value),
        ($0.GetImportResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetInsightSelectorsRequest,
            $0.GetInsightSelectorsResponse>(
        'GetInsightSelectors',
        getInsightSelectors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetInsightSelectorsRequest.fromBuffer(value),
        ($0.GetInsightSelectorsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetQueryResultsRequest,
            $0.GetQueryResultsResponse>(
        'GetQueryResults',
        getQueryResults_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetQueryResultsRequest.fromBuffer(value),
        ($0.GetQueryResultsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcePolicyRequest,
            $0.GetResourcePolicyResponse>(
        'GetResourcePolicy',
        getResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcePolicyRequest.fromBuffer(value),
        ($0.GetResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTrailRequest, $0.GetTrailResponse>(
        'GetTrail',
        getTrail_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetTrailRequest.fromBuffer(value),
        ($0.GetTrailResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTrailStatusRequest,
            $0.GetTrailStatusResponse>(
        'GetTrailStatus',
        getTrailStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTrailStatusRequest.fromBuffer(value),
        ($0.GetTrailStatusResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListChannelsRequest, $0.ListChannelsResponse>(
            'ListChannels',
            listChannels_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListChannelsRequest.fromBuffer(value),
            ($0.ListChannelsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDashboardsRequest,
            $0.ListDashboardsResponse>(
        'ListDashboards',
        listDashboards_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDashboardsRequest.fromBuffer(value),
        ($0.ListDashboardsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEventDataStoresRequest,
            $0.ListEventDataStoresResponse>(
        'ListEventDataStores',
        listEventDataStores_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEventDataStoresRequest.fromBuffer(value),
        ($0.ListEventDataStoresResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListImportFailuresRequest,
            $0.ListImportFailuresResponse>(
        'ListImportFailures',
        listImportFailures_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListImportFailuresRequest.fromBuffer(value),
        ($0.ListImportFailuresResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListImportsRequest, $0.ListImportsResponse>(
            'ListImports',
            listImports_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListImportsRequest.fromBuffer(value),
            ($0.ListImportsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListInsightsDataRequest,
            $0.ListInsightsDataResponse>(
        'ListInsightsData',
        listInsightsData_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInsightsDataRequest.fromBuffer(value),
        ($0.ListInsightsDataResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListInsightsMetricDataRequest,
            $0.ListInsightsMetricDataResponse>(
        'ListInsightsMetricData',
        listInsightsMetricData_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInsightsMetricDataRequest.fromBuffer(value),
        ($0.ListInsightsMetricDataResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPublicKeysRequest,
            $0.ListPublicKeysResponse>(
        'ListPublicKeys',
        listPublicKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPublicKeysRequest.fromBuffer(value),
        ($0.ListPublicKeysResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListQueriesRequest, $0.ListQueriesResponse>(
            'ListQueries',
            listQueries_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListQueriesRequest.fromBuffer(value),
            ($0.ListQueriesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsRequest, $0.ListTagsResponse>(
        'ListTags',
        listTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListTagsRequest.fromBuffer(value),
        ($0.ListTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTrailsRequest, $0.ListTrailsResponse>(
        'ListTrails',
        listTrails_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListTrailsRequest.fromBuffer(value),
        ($0.ListTrailsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.LookupEventsRequest, $0.LookupEventsResponse>(
            'LookupEvents',
            lookupEvents_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.LookupEventsRequest.fromBuffer(value),
            ($0.LookupEventsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutEventConfigurationRequest,
            $0.PutEventConfigurationResponse>(
        'PutEventConfiguration',
        putEventConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutEventConfigurationRequest.fromBuffer(value),
        ($0.PutEventConfigurationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutEventSelectorsRequest,
            $0.PutEventSelectorsResponse>(
        'PutEventSelectors',
        putEventSelectors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutEventSelectorsRequest.fromBuffer(value),
        ($0.PutEventSelectorsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutInsightSelectorsRequest,
            $0.PutInsightSelectorsResponse>(
        'PutInsightSelectors',
        putInsightSelectors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutInsightSelectorsRequest.fromBuffer(value),
        ($0.PutInsightSelectorsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutResourcePolicyRequest,
            $0.PutResourcePolicyResponse>(
        'PutResourcePolicy',
        putResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutResourcePolicyRequest.fromBuffer(value),
        ($0.PutResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterOrganizationDelegatedAdminRequest,
            $0.RegisterOrganizationDelegatedAdminResponse>(
        'RegisterOrganizationDelegatedAdmin',
        registerOrganizationDelegatedAdmin_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterOrganizationDelegatedAdminRequest.fromBuffer(value),
        ($0.RegisterOrganizationDelegatedAdminResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveTagsRequest, $0.RemoveTagsResponse>(
        'RemoveTags',
        removeTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.RemoveTagsRequest.fromBuffer(value),
        ($0.RemoveTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RestoreEventDataStoreRequest,
            $0.RestoreEventDataStoreResponse>(
        'RestoreEventDataStore',
        restoreEventDataStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RestoreEventDataStoreRequest.fromBuffer(value),
        ($0.RestoreEventDataStoreResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SearchSampleQueriesRequest,
            $0.SearchSampleQueriesResponse>(
        'SearchSampleQueries',
        searchSampleQueries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SearchSampleQueriesRequest.fromBuffer(value),
        ($0.SearchSampleQueriesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartDashboardRefreshRequest,
            $0.StartDashboardRefreshResponse>(
        'StartDashboardRefresh',
        startDashboardRefresh_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartDashboardRefreshRequest.fromBuffer(value),
        ($0.StartDashboardRefreshResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartEventDataStoreIngestionRequest,
            $0.StartEventDataStoreIngestionResponse>(
        'StartEventDataStoreIngestion',
        startEventDataStoreIngestion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartEventDataStoreIngestionRequest.fromBuffer(value),
        ($0.StartEventDataStoreIngestionResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartImportRequest, $0.StartImportResponse>(
            'StartImport',
            startImport_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartImportRequest.fromBuffer(value),
            ($0.StartImportResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartLoggingRequest, $0.StartLoggingResponse>(
            'StartLogging',
            startLogging_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartLoggingRequest.fromBuffer(value),
            ($0.StartLoggingResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartQueryRequest, $0.StartQueryResponse>(
        'StartQuery',
        startQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.StartQueryRequest.fromBuffer(value),
        ($0.StartQueryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopEventDataStoreIngestionRequest,
            $0.StopEventDataStoreIngestionResponse>(
        'StopEventDataStoreIngestion',
        stopEventDataStoreIngestion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopEventDataStoreIngestionRequest.fromBuffer(value),
        ($0.StopEventDataStoreIngestionResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopImportRequest, $0.StopImportResponse>(
        'StopImport',
        stopImport_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.StopImportRequest.fromBuffer(value),
        ($0.StopImportResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StopLoggingRequest, $0.StopLoggingResponse>(
            'StopLogging',
            stopLogging_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StopLoggingRequest.fromBuffer(value),
            ($0.StopLoggingResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateChannelRequest, $0.UpdateChannelResponse>(
            'UpdateChannel',
            updateChannel_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateChannelRequest.fromBuffer(value),
            ($0.UpdateChannelResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDashboardRequest,
            $0.UpdateDashboardResponse>(
        'UpdateDashboard',
        updateDashboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDashboardRequest.fromBuffer(value),
        ($0.UpdateDashboardResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateEventDataStoreRequest,
            $0.UpdateEventDataStoreResponse>(
        'UpdateEventDataStore',
        updateEventDataStore_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateEventDataStoreRequest.fromBuffer(value),
        ($0.UpdateEventDataStoreResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateTrailRequest, $0.UpdateTrailResponse>(
            'UpdateTrail',
            updateTrail_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateTrailRequest.fromBuffer(value),
            ($0.UpdateTrailResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.AddTagsResponse> addTags_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddTagsRequest> $request) async {
    return addTags($call, await $request);
  }

  $async.Future<$0.AddTagsResponse> addTags(
      $grpc.ServiceCall call, $0.AddTagsRequest request);

  $async.Future<$0.CancelQueryResponse> cancelQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CancelQueryRequest> $request) async {
    return cancelQuery($call, await $request);
  }

  $async.Future<$0.CancelQueryResponse> cancelQuery(
      $grpc.ServiceCall call, $0.CancelQueryRequest request);

  $async.Future<$0.CreateChannelResponse> createChannel_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateChannelRequest> $request) async {
    return createChannel($call, await $request);
  }

  $async.Future<$0.CreateChannelResponse> createChannel(
      $grpc.ServiceCall call, $0.CreateChannelRequest request);

  $async.Future<$0.CreateDashboardResponse> createDashboard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDashboardRequest> $request) async {
    return createDashboard($call, await $request);
  }

  $async.Future<$0.CreateDashboardResponse> createDashboard(
      $grpc.ServiceCall call, $0.CreateDashboardRequest request);

  $async.Future<$0.CreateEventDataStoreResponse> createEventDataStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateEventDataStoreRequest> $request) async {
    return createEventDataStore($call, await $request);
  }

  $async.Future<$0.CreateEventDataStoreResponse> createEventDataStore(
      $grpc.ServiceCall call, $0.CreateEventDataStoreRequest request);

  $async.Future<$0.CreateTrailResponse> createTrail_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateTrailRequest> $request) async {
    return createTrail($call, await $request);
  }

  $async.Future<$0.CreateTrailResponse> createTrail(
      $grpc.ServiceCall call, $0.CreateTrailRequest request);

  $async.Future<$0.DeleteChannelResponse> deleteChannel_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteChannelRequest> $request) async {
    return deleteChannel($call, await $request);
  }

  $async.Future<$0.DeleteChannelResponse> deleteChannel(
      $grpc.ServiceCall call, $0.DeleteChannelRequest request);

  $async.Future<$0.DeleteDashboardResponse> deleteDashboard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDashboardRequest> $request) async {
    return deleteDashboard($call, await $request);
  }

  $async.Future<$0.DeleteDashboardResponse> deleteDashboard(
      $grpc.ServiceCall call, $0.DeleteDashboardRequest request);

  $async.Future<$0.DeleteEventDataStoreResponse> deleteEventDataStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteEventDataStoreRequest> $request) async {
    return deleteEventDataStore($call, await $request);
  }

  $async.Future<$0.DeleteEventDataStoreResponse> deleteEventDataStore(
      $grpc.ServiceCall call, $0.DeleteEventDataStoreRequest request);

  $async.Future<$0.DeleteResourcePolicyResponse> deleteResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourcePolicyRequest> $request) async {
    return deleteResourcePolicy($call, await $request);
  }

  $async.Future<$0.DeleteResourcePolicyResponse> deleteResourcePolicy(
      $grpc.ServiceCall call, $0.DeleteResourcePolicyRequest request);

  $async.Future<$0.DeleteTrailResponse> deleteTrail_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteTrailRequest> $request) async {
    return deleteTrail($call, await $request);
  }

  $async.Future<$0.DeleteTrailResponse> deleteTrail(
      $grpc.ServiceCall call, $0.DeleteTrailRequest request);

  $async.Future<$0.DeregisterOrganizationDelegatedAdminResponse>
      deregisterOrganizationDelegatedAdmin_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeregisterOrganizationDelegatedAdminRequest>
              $request) async {
    return deregisterOrganizationDelegatedAdmin($call, await $request);
  }

  $async.Future<$0.DeregisterOrganizationDelegatedAdminResponse>
      deregisterOrganizationDelegatedAdmin($grpc.ServiceCall call,
          $0.DeregisterOrganizationDelegatedAdminRequest request);

  $async.Future<$0.DescribeQueryResponse> describeQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeQueryRequest> $request) async {
    return describeQuery($call, await $request);
  }

  $async.Future<$0.DescribeQueryResponse> describeQuery(
      $grpc.ServiceCall call, $0.DescribeQueryRequest request);

  $async.Future<$0.DescribeTrailsResponse> describeTrails_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeTrailsRequest> $request) async {
    return describeTrails($call, await $request);
  }

  $async.Future<$0.DescribeTrailsResponse> describeTrails(
      $grpc.ServiceCall call, $0.DescribeTrailsRequest request);

  $async.Future<$0.DisableFederationResponse> disableFederation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DisableFederationRequest> $request) async {
    return disableFederation($call, await $request);
  }

  $async.Future<$0.DisableFederationResponse> disableFederation(
      $grpc.ServiceCall call, $0.DisableFederationRequest request);

  $async.Future<$0.EnableFederationResponse> enableFederation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.EnableFederationRequest> $request) async {
    return enableFederation($call, await $request);
  }

  $async.Future<$0.EnableFederationResponse> enableFederation(
      $grpc.ServiceCall call, $0.EnableFederationRequest request);

  $async.Future<$0.GenerateQueryResponse> generateQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GenerateQueryRequest> $request) async {
    return generateQuery($call, await $request);
  }

  $async.Future<$0.GenerateQueryResponse> generateQuery(
      $grpc.ServiceCall call, $0.GenerateQueryRequest request);

  $async.Future<$0.GetChannelResponse> getChannel_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetChannelRequest> $request) async {
    return getChannel($call, await $request);
  }

  $async.Future<$0.GetChannelResponse> getChannel(
      $grpc.ServiceCall call, $0.GetChannelRequest request);

  $async.Future<$0.GetDashboardResponse> getDashboard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDashboardRequest> $request) async {
    return getDashboard($call, await $request);
  }

  $async.Future<$0.GetDashboardResponse> getDashboard(
      $grpc.ServiceCall call, $0.GetDashboardRequest request);

  $async.Future<$0.GetEventConfigurationResponse> getEventConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEventConfigurationRequest> $request) async {
    return getEventConfiguration($call, await $request);
  }

  $async.Future<$0.GetEventConfigurationResponse> getEventConfiguration(
      $grpc.ServiceCall call, $0.GetEventConfigurationRequest request);

  $async.Future<$0.GetEventDataStoreResponse> getEventDataStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEventDataStoreRequest> $request) async {
    return getEventDataStore($call, await $request);
  }

  $async.Future<$0.GetEventDataStoreResponse> getEventDataStore(
      $grpc.ServiceCall call, $0.GetEventDataStoreRequest request);

  $async.Future<$0.GetEventSelectorsResponse> getEventSelectors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEventSelectorsRequest> $request) async {
    return getEventSelectors($call, await $request);
  }

  $async.Future<$0.GetEventSelectorsResponse> getEventSelectors(
      $grpc.ServiceCall call, $0.GetEventSelectorsRequest request);

  $async.Future<$0.GetImportResponse> getImport_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetImportRequest> $request) async {
    return getImport($call, await $request);
  }

  $async.Future<$0.GetImportResponse> getImport(
      $grpc.ServiceCall call, $0.GetImportRequest request);

  $async.Future<$0.GetInsightSelectorsResponse> getInsightSelectors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetInsightSelectorsRequest> $request) async {
    return getInsightSelectors($call, await $request);
  }

  $async.Future<$0.GetInsightSelectorsResponse> getInsightSelectors(
      $grpc.ServiceCall call, $0.GetInsightSelectorsRequest request);

  $async.Future<$0.GetQueryResultsResponse> getQueryResults_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetQueryResultsRequest> $request) async {
    return getQueryResults($call, await $request);
  }

  $async.Future<$0.GetQueryResultsResponse> getQueryResults(
      $grpc.ServiceCall call, $0.GetQueryResultsRequest request);

  $async.Future<$0.GetResourcePolicyResponse> getResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourcePolicyRequest> $request) async {
    return getResourcePolicy($call, await $request);
  }

  $async.Future<$0.GetResourcePolicyResponse> getResourcePolicy(
      $grpc.ServiceCall call, $0.GetResourcePolicyRequest request);

  $async.Future<$0.GetTrailResponse> getTrail_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetTrailRequest> $request) async {
    return getTrail($call, await $request);
  }

  $async.Future<$0.GetTrailResponse> getTrail(
      $grpc.ServiceCall call, $0.GetTrailRequest request);

  $async.Future<$0.GetTrailStatusResponse> getTrailStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTrailStatusRequest> $request) async {
    return getTrailStatus($call, await $request);
  }

  $async.Future<$0.GetTrailStatusResponse> getTrailStatus(
      $grpc.ServiceCall call, $0.GetTrailStatusRequest request);

  $async.Future<$0.ListChannelsResponse> listChannels_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListChannelsRequest> $request) async {
    return listChannels($call, await $request);
  }

  $async.Future<$0.ListChannelsResponse> listChannels(
      $grpc.ServiceCall call, $0.ListChannelsRequest request);

  $async.Future<$0.ListDashboardsResponse> listDashboards_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDashboardsRequest> $request) async {
    return listDashboards($call, await $request);
  }

  $async.Future<$0.ListDashboardsResponse> listDashboards(
      $grpc.ServiceCall call, $0.ListDashboardsRequest request);

  $async.Future<$0.ListEventDataStoresResponse> listEventDataStores_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEventDataStoresRequest> $request) async {
    return listEventDataStores($call, await $request);
  }

  $async.Future<$0.ListEventDataStoresResponse> listEventDataStores(
      $grpc.ServiceCall call, $0.ListEventDataStoresRequest request);

  $async.Future<$0.ListImportFailuresResponse> listImportFailures_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListImportFailuresRequest> $request) async {
    return listImportFailures($call, await $request);
  }

  $async.Future<$0.ListImportFailuresResponse> listImportFailures(
      $grpc.ServiceCall call, $0.ListImportFailuresRequest request);

  $async.Future<$0.ListImportsResponse> listImports_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListImportsRequest> $request) async {
    return listImports($call, await $request);
  }

  $async.Future<$0.ListImportsResponse> listImports(
      $grpc.ServiceCall call, $0.ListImportsRequest request);

  $async.Future<$0.ListInsightsDataResponse> listInsightsData_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListInsightsDataRequest> $request) async {
    return listInsightsData($call, await $request);
  }

  $async.Future<$0.ListInsightsDataResponse> listInsightsData(
      $grpc.ServiceCall call, $0.ListInsightsDataRequest request);

  $async.Future<$0.ListInsightsMetricDataResponse> listInsightsMetricData_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListInsightsMetricDataRequest> $request) async {
    return listInsightsMetricData($call, await $request);
  }

  $async.Future<$0.ListInsightsMetricDataResponse> listInsightsMetricData(
      $grpc.ServiceCall call, $0.ListInsightsMetricDataRequest request);

  $async.Future<$0.ListPublicKeysResponse> listPublicKeys_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPublicKeysRequest> $request) async {
    return listPublicKeys($call, await $request);
  }

  $async.Future<$0.ListPublicKeysResponse> listPublicKeys(
      $grpc.ServiceCall call, $0.ListPublicKeysRequest request);

  $async.Future<$0.ListQueriesResponse> listQueries_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListQueriesRequest> $request) async {
    return listQueries($call, await $request);
  }

  $async.Future<$0.ListQueriesResponse> listQueries(
      $grpc.ServiceCall call, $0.ListQueriesRequest request);

  $async.Future<$0.ListTagsResponse> listTags_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTagsRequest> $request) async {
    return listTags($call, await $request);
  }

  $async.Future<$0.ListTagsResponse> listTags(
      $grpc.ServiceCall call, $0.ListTagsRequest request);

  $async.Future<$0.ListTrailsResponse> listTrails_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTrailsRequest> $request) async {
    return listTrails($call, await $request);
  }

  $async.Future<$0.ListTrailsResponse> listTrails(
      $grpc.ServiceCall call, $0.ListTrailsRequest request);

  $async.Future<$0.LookupEventsResponse> lookupEvents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.LookupEventsRequest> $request) async {
    return lookupEvents($call, await $request);
  }

  $async.Future<$0.LookupEventsResponse> lookupEvents(
      $grpc.ServiceCall call, $0.LookupEventsRequest request);

  $async.Future<$0.PutEventConfigurationResponse> putEventConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutEventConfigurationRequest> $request) async {
    return putEventConfiguration($call, await $request);
  }

  $async.Future<$0.PutEventConfigurationResponse> putEventConfiguration(
      $grpc.ServiceCall call, $0.PutEventConfigurationRequest request);

  $async.Future<$0.PutEventSelectorsResponse> putEventSelectors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutEventSelectorsRequest> $request) async {
    return putEventSelectors($call, await $request);
  }

  $async.Future<$0.PutEventSelectorsResponse> putEventSelectors(
      $grpc.ServiceCall call, $0.PutEventSelectorsRequest request);

  $async.Future<$0.PutInsightSelectorsResponse> putInsightSelectors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutInsightSelectorsRequest> $request) async {
    return putInsightSelectors($call, await $request);
  }

  $async.Future<$0.PutInsightSelectorsResponse> putInsightSelectors(
      $grpc.ServiceCall call, $0.PutInsightSelectorsRequest request);

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutResourcePolicyRequest> $request) async {
    return putResourcePolicy($call, await $request);
  }

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy(
      $grpc.ServiceCall call, $0.PutResourcePolicyRequest request);

  $async.Future<$0.RegisterOrganizationDelegatedAdminResponse>
      registerOrganizationDelegatedAdmin_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.RegisterOrganizationDelegatedAdminRequest>
              $request) async {
    return registerOrganizationDelegatedAdmin($call, await $request);
  }

  $async.Future<$0.RegisterOrganizationDelegatedAdminResponse>
      registerOrganizationDelegatedAdmin($grpc.ServiceCall call,
          $0.RegisterOrganizationDelegatedAdminRequest request);

  $async.Future<$0.RemoveTagsResponse> removeTags_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemoveTagsRequest> $request) async {
    return removeTags($call, await $request);
  }

  $async.Future<$0.RemoveTagsResponse> removeTags(
      $grpc.ServiceCall call, $0.RemoveTagsRequest request);

  $async.Future<$0.RestoreEventDataStoreResponse> restoreEventDataStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RestoreEventDataStoreRequest> $request) async {
    return restoreEventDataStore($call, await $request);
  }

  $async.Future<$0.RestoreEventDataStoreResponse> restoreEventDataStore(
      $grpc.ServiceCall call, $0.RestoreEventDataStoreRequest request);

  $async.Future<$0.SearchSampleQueriesResponse> searchSampleQueries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SearchSampleQueriesRequest> $request) async {
    return searchSampleQueries($call, await $request);
  }

  $async.Future<$0.SearchSampleQueriesResponse> searchSampleQueries(
      $grpc.ServiceCall call, $0.SearchSampleQueriesRequest request);

  $async.Future<$0.StartDashboardRefreshResponse> startDashboardRefresh_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartDashboardRefreshRequest> $request) async {
    return startDashboardRefresh($call, await $request);
  }

  $async.Future<$0.StartDashboardRefreshResponse> startDashboardRefresh(
      $grpc.ServiceCall call, $0.StartDashboardRefreshRequest request);

  $async.Future<$0.StartEventDataStoreIngestionResponse>
      startEventDataStoreIngestion_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.StartEventDataStoreIngestionRequest>
              $request) async {
    return startEventDataStoreIngestion($call, await $request);
  }

  $async.Future<$0.StartEventDataStoreIngestionResponse>
      startEventDataStoreIngestion($grpc.ServiceCall call,
          $0.StartEventDataStoreIngestionRequest request);

  $async.Future<$0.StartImportResponse> startImport_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StartImportRequest> $request) async {
    return startImport($call, await $request);
  }

  $async.Future<$0.StartImportResponse> startImport(
      $grpc.ServiceCall call, $0.StartImportRequest request);

  $async.Future<$0.StartLoggingResponse> startLogging_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartLoggingRequest> $request) async {
    return startLogging($call, await $request);
  }

  $async.Future<$0.StartLoggingResponse> startLogging(
      $grpc.ServiceCall call, $0.StartLoggingRequest request);

  $async.Future<$0.StartQueryResponse> startQuery_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StartQueryRequest> $request) async {
    return startQuery($call, await $request);
  }

  $async.Future<$0.StartQueryResponse> startQuery(
      $grpc.ServiceCall call, $0.StartQueryRequest request);

  $async.Future<$0.StopEventDataStoreIngestionResponse>
      stopEventDataStoreIngestion_Pre($grpc.ServiceCall $call,
          $async.Future<$0.StopEventDataStoreIngestionRequest> $request) async {
    return stopEventDataStoreIngestion($call, await $request);
  }

  $async.Future<$0.StopEventDataStoreIngestionResponse>
      stopEventDataStoreIngestion($grpc.ServiceCall call,
          $0.StopEventDataStoreIngestionRequest request);

  $async.Future<$0.StopImportResponse> stopImport_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StopImportRequest> $request) async {
    return stopImport($call, await $request);
  }

  $async.Future<$0.StopImportResponse> stopImport(
      $grpc.ServiceCall call, $0.StopImportRequest request);

  $async.Future<$0.StopLoggingResponse> stopLogging_Pre($grpc.ServiceCall $call,
      $async.Future<$0.StopLoggingRequest> $request) async {
    return stopLogging($call, await $request);
  }

  $async.Future<$0.StopLoggingResponse> stopLogging(
      $grpc.ServiceCall call, $0.StopLoggingRequest request);

  $async.Future<$0.UpdateChannelResponse> updateChannel_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateChannelRequest> $request) async {
    return updateChannel($call, await $request);
  }

  $async.Future<$0.UpdateChannelResponse> updateChannel(
      $grpc.ServiceCall call, $0.UpdateChannelRequest request);

  $async.Future<$0.UpdateDashboardResponse> updateDashboard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDashboardRequest> $request) async {
    return updateDashboard($call, await $request);
  }

  $async.Future<$0.UpdateDashboardResponse> updateDashboard(
      $grpc.ServiceCall call, $0.UpdateDashboardRequest request);

  $async.Future<$0.UpdateEventDataStoreResponse> updateEventDataStore_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateEventDataStoreRequest> $request) async {
    return updateEventDataStore($call, await $request);
  }

  $async.Future<$0.UpdateEventDataStoreResponse> updateEventDataStore(
      $grpc.ServiceCall call, $0.UpdateEventDataStoreRequest request);

  $async.Future<$0.UpdateTrailResponse> updateTrail_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateTrailRequest> $request) async {
    return updateTrail($call, await $request);
  }

  $async.Future<$0.UpdateTrailResponse> updateTrail(
      $grpc.ServiceCall call, $0.UpdateTrailRequest request);
}
