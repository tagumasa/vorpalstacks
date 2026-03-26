// This is a generated file - do not edit.
//
// Generated from athena.proto.

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

import 'athena.pb.dart' as $0;

export 'athena.pb.dart';

/// AthenaService provides athena API operations.
@$pb.GrpcServiceName('athena.AthenaService')
class AthenaServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  AthenaServiceClient(super.channel, {super.options, super.interceptors});

  /// Returns the details of a single named query or a list of up to 50 queries, which you provide as an array of query ID strings. Requires you to have access to the workgroup in which the queries were ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.BatchGetNamedQueryOutput> batchGetNamedQuery(
    $0.BatchGetNamedQueryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchGetNamedQuery, request, options: options);
  }

  /// Returns the details of a single prepared statement or a list of up to 256 prepared statements for the array of prepared statement names that you provide. Requires you to have access to the workgrou...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.BatchGetPreparedStatementOutput>
      batchGetPreparedStatement(
    $0.BatchGetPreparedStatementInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchGetPreparedStatement, request,
        options: options);
  }

  /// Returns the details of a single query execution or a list of up to 50 query executions, which you provide as an array of query execution ID strings. Requires you to have access to the workgroup in ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.BatchGetQueryExecutionOutput> batchGetQueryExecution(
    $0.BatchGetQueryExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchGetQueryExecution, request,
        options: options);
  }

  /// Cancels the capacity reservation with the specified name. Cancelled reservations remain in your account and will be deleted 45 days after cancellation. During the 45 days, you cannot re-purpose or ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelCapacityReservationOutput>
      cancelCapacityReservation(
    $0.CancelCapacityReservationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelCapacityReservation, request,
        options: options);
  }

  /// Creates a capacity reservation with the specified name and number of requested data processing units.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateCapacityReservationOutput>
      createCapacityReservation(
    $0.CreateCapacityReservationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCapacityReservation, request,
        options: options);
  }

  /// Creates (registers) a data catalog with the specified name and properties. Catalogs created are visible to all users of the same Amazon Web Services account. For a FEDERATED catalog, this API opera...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateDataCatalogOutput> createDataCatalog(
    $0.CreateDataCatalogInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDataCatalog, request, options: options);
  }

  /// Creates a named query in the specified workgroup. Requires that you have access to the workgroup.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateNamedQueryOutput> createNamedQuery(
    $0.CreateNamedQueryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createNamedQuery, request, options: options);
  }

  /// Creates an empty ipynb file in the specified Apache Spark enabled workgroup. Throws an error if a file in the workgroup with the same name already exists.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateNotebookOutput> createNotebook(
    $0.CreateNotebookInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createNotebook, request, options: options);
  }

  /// Creates a prepared statement for use with SQL queries in Athena.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreatePreparedStatementOutput>
      createPreparedStatement(
    $0.CreatePreparedStatementInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPreparedStatement, request,
        options: options);
  }

  /// Gets an authentication token and the URL at which the notebook can be accessed. During programmatic access, CreatePresignedNotebookUrl must be called every 10 minutes to refresh the authentication ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreatePresignedNotebookUrlResponse>
      createPresignedNotebookUrl(
    $0.CreatePresignedNotebookUrlRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPresignedNotebookUrl, request,
        options: options);
  }

  /// Creates a workgroup with the specified name. A workgroup can be an Apache Spark enabled workgroup or an Athena SQL workgroup.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateWorkGroupOutput> createWorkGroup(
    $0.CreateWorkGroupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createWorkGroup, request, options: options);
  }

  /// Deletes a cancelled capacity reservation. A reservation must be cancelled before it can be deleted. A deleted reservation is immediately removed from your account and can no longer be referenced, i...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteCapacityReservationOutput>
      deleteCapacityReservation(
    $0.DeleteCapacityReservationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCapacityReservation, request,
        options: options);
  }

  /// Deletes a data catalog.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteDataCatalogOutput> deleteDataCatalog(
    $0.DeleteDataCatalogInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDataCatalog, request, options: options);
  }

  /// Deletes the named query if you have access to the workgroup in which the query was saved.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteNamedQueryOutput> deleteNamedQuery(
    $0.DeleteNamedQueryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteNamedQuery, request, options: options);
  }

  /// Deletes the specified notebook.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteNotebookOutput> deleteNotebook(
    $0.DeleteNotebookInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteNotebook, request, options: options);
  }

  /// Deletes the prepared statement with the specified name from the specified workgroup.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeletePreparedStatementOutput>
      deletePreparedStatement(
    $0.DeletePreparedStatementInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePreparedStatement, request,
        options: options);
  }

  /// Deletes the workgroup with the specified name. The primary workgroup cannot be deleted.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteWorkGroupOutput> deleteWorkGroup(
    $0.DeleteWorkGroupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteWorkGroup, request, options: options);
  }

  /// Exports the specified notebook and its metadata.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ExportNotebookOutput> exportNotebook(
    $0.ExportNotebookInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$exportNotebook, request, options: options);
  }

  /// Describes a previously submitted calculation execution.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCalculationExecutionResponse>
      getCalculationExecution(
    $0.GetCalculationExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCalculationExecution, request,
        options: options);
  }

  /// Retrieves the unencrypted code that was executed for the calculation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCalculationExecutionCodeResponse>
      getCalculationExecutionCode(
    $0.GetCalculationExecutionCodeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCalculationExecutionCode, request,
        options: options);
  }

  /// Gets the status of a current calculation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCalculationExecutionStatusResponse>
      getCalculationExecutionStatus(
    $0.GetCalculationExecutionStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCalculationExecutionStatus, request,
        options: options);
  }

  /// Gets the capacity assignment configuration for a capacity reservation, if one exists.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCapacityAssignmentConfigurationOutput>
      getCapacityAssignmentConfiguration(
    $0.GetCapacityAssignmentConfigurationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCapacityAssignmentConfiguration, request,
        options: options);
  }

  /// Returns information about the capacity reservation with the specified name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCapacityReservationOutput> getCapacityReservation(
    $0.GetCapacityReservationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCapacityReservation, request,
        options: options);
  }

  /// Returns a database object for the specified database and data catalog.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDatabaseOutput> getDatabase(
    $0.GetDatabaseInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDatabase, request, options: options);
  }

  /// Returns the specified data catalog.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDataCatalogOutput> getDataCatalog(
    $0.GetDataCatalogInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDataCatalog, request, options: options);
  }

  /// Returns information about a single query. Requires that you have access to the workgroup in which the query was saved.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetNamedQueryOutput> getNamedQuery(
    $0.GetNamedQueryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getNamedQuery, request, options: options);
  }

  /// Retrieves notebook metadata for the specified notebook ID.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetNotebookMetadataOutput> getNotebookMetadata(
    $0.GetNotebookMetadataInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getNotebookMetadata, request, options: options);
  }

  /// Retrieves the prepared statement with the specified name from the specified workgroup.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetPreparedStatementOutput> getPreparedStatement(
    $0.GetPreparedStatementInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPreparedStatement, request, options: options);
  }

  /// Returns information about a single execution of a query if you have access to the workgroup in which the query ran. Each time a query executes, information about the query execution is saved with a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetQueryExecutionOutput> getQueryExecution(
    $0.GetQueryExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueryExecution, request, options: options);
  }

  /// Streams the results of a single query execution specified by QueryExecutionId from the Athena query results location in Amazon S3. For more information, see Working with query results, recent queri...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetQueryResultsOutput> getQueryResults(
    $0.GetQueryResultsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueryResults, request, options: options);
  }

  /// Returns query execution runtime statistics related to a single execution of a query if you have access to the workgroup in which the query ran. Statistics from the Timeline section of the response ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetQueryRuntimeStatisticsOutput>
      getQueryRuntimeStatistics(
    $0.GetQueryRuntimeStatisticsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getQueryRuntimeStatistics, request,
        options: options);
  }

  /// Gets the Live UI/Persistence UI for a session.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetResourceDashboardResponse> getResourceDashboard(
    $0.GetResourceDashboardRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourceDashboard, request, options: options);
  }

  /// Gets the full details of a previously created session, including the session status and configuration.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSessionResponse> getSession(
    $0.GetSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSession, request, options: options);
  }

  /// Gets a connection endpoint and authentication token for a given session Id.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSessionEndpointResponse> getSessionEndpoint(
    $0.GetSessionEndpointRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSessionEndpoint, request, options: options);
  }

  /// Gets the current status of a session.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSessionStatusResponse> getSessionStatus(
    $0.GetSessionStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSessionStatus, request, options: options);
  }

  /// Returns table metadata for the specified catalog, database, and table.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetTableMetadataOutput> getTableMetadata(
    $0.GetTableMetadataInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTableMetadata, request, options: options);
  }

  /// Returns information about the workgroup with the specified name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetWorkGroupOutput> getWorkGroup(
    $0.GetWorkGroupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getWorkGroup, request, options: options);
  }

  /// Imports a single ipynb file to a Spark enabled workgroup. To import the notebook, the request must specify a value for either Payload or NoteBookS3LocationUri. If neither is specified or both are s...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ImportNotebookOutput> importNotebook(
    $0.ImportNotebookInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importNotebook, request, options: options);
  }

  /// Returns the supported DPU sizes for the supported application runtimes (for example, Athena notebook version 1).
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListApplicationDPUSizesOutput>
      listApplicationDPUSizes(
    $0.ListApplicationDPUSizesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listApplicationDPUSizes, request,
        options: options);
  }

  /// Lists the calculations that have been submitted to a session in descending order. Newer calculations are listed first; older calculations are listed later.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListCalculationExecutionsResponse>
      listCalculationExecutions(
    $0.ListCalculationExecutionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCalculationExecutions, request,
        options: options);
  }

  /// Lists the capacity reservations for the current account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListCapacityReservationsOutput>
      listCapacityReservations(
    $0.ListCapacityReservationsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCapacityReservations, request,
        options: options);
  }

  /// Lists the databases in the specified data catalog.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListDatabasesOutput> listDatabases(
    $0.ListDatabasesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDatabases, request, options: options);
  }

  /// Lists the data catalogs in the current Amazon Web Services account. In the Athena console, data catalogs are listed as "data sources" on the Data sources page under the Data source name column.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListDataCatalogsOutput> listDataCatalogs(
    $0.ListDataCatalogsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDataCatalogs, request, options: options);
  }

  /// Returns a list of engine versions that are available to choose from, including the Auto option.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListEngineVersionsOutput> listEngineVersions(
    $0.ListEngineVersionsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEngineVersions, request, options: options);
  }

  /// Lists, in descending order, the executors that joined a session. Newer executors are listed first; older executors are listed later. The result can be optionally filtered by state.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListExecutorsResponse> listExecutors(
    $0.ListExecutorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listExecutors, request, options: options);
  }

  /// Provides a list of available query IDs only for queries saved in the specified workgroup. Requires that you have access to the specified workgroup. If a workgroup is not specified, lists the saved ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListNamedQueriesOutput> listNamedQueries(
    $0.ListNamedQueriesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listNamedQueries, request, options: options);
  }

  /// Displays the notebook files for the specified workgroup in paginated format.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListNotebookMetadataOutput> listNotebookMetadata(
    $0.ListNotebookMetadataInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listNotebookMetadata, request, options: options);
  }

  /// Lists, in descending order, the sessions that have been created in a notebook that are in an active state like CREATING, CREATED, IDLE or BUSY. Newer sessions are listed first; older sessions are l...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListNotebookSessionsResponse> listNotebookSessions(
    $0.ListNotebookSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listNotebookSessions, request, options: options);
  }

  /// Lists the prepared statements in the specified workgroup.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListPreparedStatementsOutput> listPreparedStatements(
    $0.ListPreparedStatementsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPreparedStatements, request,
        options: options);
  }

  /// Provides a list of available query execution IDs for the queries in the specified workgroup. Athena keeps a query history for 45 days. If a workgroup is not specified, returns a list of query execu...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListQueryExecutionsOutput> listQueryExecutions(
    $0.ListQueryExecutionsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listQueryExecutions, request, options: options);
  }

  /// Lists the sessions in a workgroup that are in an active state like CREATING, CREATED, IDLE, or BUSY. Newer sessions are listed first; older sessions are listed later.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListSessionsResponse> listSessions(
    $0.ListSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSessions, request, options: options);
  }

  /// Lists the metadata for the tables in the specified data catalog database.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTableMetadataOutput> listTableMetadata(
    $0.ListTableMetadataInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTableMetadata, request, options: options);
  }

  /// Lists the tags associated with an Athena resource.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceOutput> listTagsForResource(
    $0.ListTagsForResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Lists available workgroups for the account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListWorkGroupsOutput> listWorkGroups(
    $0.ListWorkGroupsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listWorkGroups, request, options: options);
  }

  /// Puts a new capacity assignment configuration for a specified capacity reservation. If a capacity assignment configuration already exists for the capacity reservation, replaces the existing capacity...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutCapacityAssignmentConfigurationOutput>
      putCapacityAssignmentConfiguration(
    $0.PutCapacityAssignmentConfigurationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putCapacityAssignmentConfiguration, request,
        options: options);
  }

  /// Submits calculations for execution within a session. You can supply the code to run as an inline code block within the request. The request syntax requires the StartCalculationExecutionRequest$Code...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartCalculationExecutionResponse>
      startCalculationExecution(
    $0.StartCalculationExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startCalculationExecution, request,
        options: options);
  }

  /// Runs the SQL query statements contained in the Query. Requires you to have access to the workgroup in which the query ran. Running queries against an external catalog requires GetDataCatalog permis...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartQueryExecutionOutput> startQueryExecution(
    $0.StartQueryExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startQueryExecution, request, options: options);
  }

  /// Creates a session for running calculations within a workgroup. The session is ready when it reaches an IDLE state.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartSessionResponse> startSession(
    $0.StartSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startSession, request, options: options);
  }

  /// Requests the cancellation of a calculation. A StopCalculationExecution call on a calculation that is already in a terminal state (for example, STOPPED, FAILED, or COMPLETED) succeeds but has no eff...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopCalculationExecutionResponse>
      stopCalculationExecution(
    $0.StopCalculationExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopCalculationExecution, request,
        options: options);
  }

  /// Stops a query execution. Requires you to have access to the workgroup in which the query ran.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopQueryExecutionOutput> stopQueryExecution(
    $0.StopQueryExecutionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopQueryExecution, request, options: options);
  }

  /// Adds one or more tags to an Athena resource. A tag is a label that you assign to a resource. Each tag consists of a key and an optional value, both of which you define. For example, you can use tag...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TagResourceOutput> tagResource(
    $0.TagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Terminates an active session. A TerminateSession call on a session that is already inactive (for example, in a FAILED, TERMINATED or TERMINATING state) succeeds but has no effect. Calculations runn...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TerminateSessionResponse> terminateSession(
    $0.TerminateSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$terminateSession, request, options: options);
  }

  /// Removes one or more tags from an Athena resource.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UntagResourceOutput> untagResource(
    $0.UntagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates the number of requested data processing units for the capacity reservation with the specified name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateCapacityReservationOutput>
      updateCapacityReservation(
    $0.UpdateCapacityReservationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateCapacityReservation, request,
        options: options);
  }

  /// Updates the data catalog that has the specified name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateDataCatalogOutput> updateDataCatalog(
    $0.UpdateDataCatalogInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDataCatalog, request, options: options);
  }

  /// Updates a NamedQuery object. The database or workgroup cannot be updated.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateNamedQueryOutput> updateNamedQuery(
    $0.UpdateNamedQueryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateNamedQuery, request, options: options);
  }

  /// Updates the contents of a Spark notebook.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateNotebookOutput> updateNotebook(
    $0.UpdateNotebookInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateNotebook, request, options: options);
  }

  /// Updates the metadata for a notebook.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateNotebookMetadataOutput> updateNotebookMetadata(
    $0.UpdateNotebookMetadataInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateNotebookMetadata, request,
        options: options);
  }

  /// Updates a prepared statement.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdatePreparedStatementOutput>
      updatePreparedStatement(
    $0.UpdatePreparedStatementInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updatePreparedStatement, request,
        options: options);
  }

  /// Updates the workgroup with the specified name. The workgroup's name cannot be changed. Only ConfigurationUpdates can be specified.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateWorkGroupOutput> updateWorkGroup(
    $0.UpdateWorkGroupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateWorkGroup, request, options: options);
  }

  // method descriptors

  static final _$batchGetNamedQuery = $grpc.ClientMethod<
          $0.BatchGetNamedQueryInput, $0.BatchGetNamedQueryOutput>(
      '/athena.AthenaService/BatchGetNamedQuery',
      ($0.BatchGetNamedQueryInput value) => value.writeToBuffer(),
      $0.BatchGetNamedQueryOutput.fromBuffer);
  static final _$batchGetPreparedStatement = $grpc.ClientMethod<
          $0.BatchGetPreparedStatementInput,
          $0.BatchGetPreparedStatementOutput>(
      '/athena.AthenaService/BatchGetPreparedStatement',
      ($0.BatchGetPreparedStatementInput value) => value.writeToBuffer(),
      $0.BatchGetPreparedStatementOutput.fromBuffer);
  static final _$batchGetQueryExecution = $grpc.ClientMethod<
          $0.BatchGetQueryExecutionInput, $0.BatchGetQueryExecutionOutput>(
      '/athena.AthenaService/BatchGetQueryExecution',
      ($0.BatchGetQueryExecutionInput value) => value.writeToBuffer(),
      $0.BatchGetQueryExecutionOutput.fromBuffer);
  static final _$cancelCapacityReservation = $grpc.ClientMethod<
          $0.CancelCapacityReservationInput,
          $0.CancelCapacityReservationOutput>(
      '/athena.AthenaService/CancelCapacityReservation',
      ($0.CancelCapacityReservationInput value) => value.writeToBuffer(),
      $0.CancelCapacityReservationOutput.fromBuffer);
  static final _$createCapacityReservation = $grpc.ClientMethod<
          $0.CreateCapacityReservationInput,
          $0.CreateCapacityReservationOutput>(
      '/athena.AthenaService/CreateCapacityReservation',
      ($0.CreateCapacityReservationInput value) => value.writeToBuffer(),
      $0.CreateCapacityReservationOutput.fromBuffer);
  static final _$createDataCatalog =
      $grpc.ClientMethod<$0.CreateDataCatalogInput, $0.CreateDataCatalogOutput>(
          '/athena.AthenaService/CreateDataCatalog',
          ($0.CreateDataCatalogInput value) => value.writeToBuffer(),
          $0.CreateDataCatalogOutput.fromBuffer);
  static final _$createNamedQuery =
      $grpc.ClientMethod<$0.CreateNamedQueryInput, $0.CreateNamedQueryOutput>(
          '/athena.AthenaService/CreateNamedQuery',
          ($0.CreateNamedQueryInput value) => value.writeToBuffer(),
          $0.CreateNamedQueryOutput.fromBuffer);
  static final _$createNotebook =
      $grpc.ClientMethod<$0.CreateNotebookInput, $0.CreateNotebookOutput>(
          '/athena.AthenaService/CreateNotebook',
          ($0.CreateNotebookInput value) => value.writeToBuffer(),
          $0.CreateNotebookOutput.fromBuffer);
  static final _$createPreparedStatement = $grpc.ClientMethod<
          $0.CreatePreparedStatementInput, $0.CreatePreparedStatementOutput>(
      '/athena.AthenaService/CreatePreparedStatement',
      ($0.CreatePreparedStatementInput value) => value.writeToBuffer(),
      $0.CreatePreparedStatementOutput.fromBuffer);
  static final _$createPresignedNotebookUrl = $grpc.ClientMethod<
          $0.CreatePresignedNotebookUrlRequest,
          $0.CreatePresignedNotebookUrlResponse>(
      '/athena.AthenaService/CreatePresignedNotebookUrl',
      ($0.CreatePresignedNotebookUrlRequest value) => value.writeToBuffer(),
      $0.CreatePresignedNotebookUrlResponse.fromBuffer);
  static final _$createWorkGroup =
      $grpc.ClientMethod<$0.CreateWorkGroupInput, $0.CreateWorkGroupOutput>(
          '/athena.AthenaService/CreateWorkGroup',
          ($0.CreateWorkGroupInput value) => value.writeToBuffer(),
          $0.CreateWorkGroupOutput.fromBuffer);
  static final _$deleteCapacityReservation = $grpc.ClientMethod<
          $0.DeleteCapacityReservationInput,
          $0.DeleteCapacityReservationOutput>(
      '/athena.AthenaService/DeleteCapacityReservation',
      ($0.DeleteCapacityReservationInput value) => value.writeToBuffer(),
      $0.DeleteCapacityReservationOutput.fromBuffer);
  static final _$deleteDataCatalog =
      $grpc.ClientMethod<$0.DeleteDataCatalogInput, $0.DeleteDataCatalogOutput>(
          '/athena.AthenaService/DeleteDataCatalog',
          ($0.DeleteDataCatalogInput value) => value.writeToBuffer(),
          $0.DeleteDataCatalogOutput.fromBuffer);
  static final _$deleteNamedQuery =
      $grpc.ClientMethod<$0.DeleteNamedQueryInput, $0.DeleteNamedQueryOutput>(
          '/athena.AthenaService/DeleteNamedQuery',
          ($0.DeleteNamedQueryInput value) => value.writeToBuffer(),
          $0.DeleteNamedQueryOutput.fromBuffer);
  static final _$deleteNotebook =
      $grpc.ClientMethod<$0.DeleteNotebookInput, $0.DeleteNotebookOutput>(
          '/athena.AthenaService/DeleteNotebook',
          ($0.DeleteNotebookInput value) => value.writeToBuffer(),
          $0.DeleteNotebookOutput.fromBuffer);
  static final _$deletePreparedStatement = $grpc.ClientMethod<
          $0.DeletePreparedStatementInput, $0.DeletePreparedStatementOutput>(
      '/athena.AthenaService/DeletePreparedStatement',
      ($0.DeletePreparedStatementInput value) => value.writeToBuffer(),
      $0.DeletePreparedStatementOutput.fromBuffer);
  static final _$deleteWorkGroup =
      $grpc.ClientMethod<$0.DeleteWorkGroupInput, $0.DeleteWorkGroupOutput>(
          '/athena.AthenaService/DeleteWorkGroup',
          ($0.DeleteWorkGroupInput value) => value.writeToBuffer(),
          $0.DeleteWorkGroupOutput.fromBuffer);
  static final _$exportNotebook =
      $grpc.ClientMethod<$0.ExportNotebookInput, $0.ExportNotebookOutput>(
          '/athena.AthenaService/ExportNotebook',
          ($0.ExportNotebookInput value) => value.writeToBuffer(),
          $0.ExportNotebookOutput.fromBuffer);
  static final _$getCalculationExecution = $grpc.ClientMethod<
          $0.GetCalculationExecutionRequest,
          $0.GetCalculationExecutionResponse>(
      '/athena.AthenaService/GetCalculationExecution',
      ($0.GetCalculationExecutionRequest value) => value.writeToBuffer(),
      $0.GetCalculationExecutionResponse.fromBuffer);
  static final _$getCalculationExecutionCode = $grpc.ClientMethod<
          $0.GetCalculationExecutionCodeRequest,
          $0.GetCalculationExecutionCodeResponse>(
      '/athena.AthenaService/GetCalculationExecutionCode',
      ($0.GetCalculationExecutionCodeRequest value) => value.writeToBuffer(),
      $0.GetCalculationExecutionCodeResponse.fromBuffer);
  static final _$getCalculationExecutionStatus = $grpc.ClientMethod<
          $0.GetCalculationExecutionStatusRequest,
          $0.GetCalculationExecutionStatusResponse>(
      '/athena.AthenaService/GetCalculationExecutionStatus',
      ($0.GetCalculationExecutionStatusRequest value) => value.writeToBuffer(),
      $0.GetCalculationExecutionStatusResponse.fromBuffer);
  static final _$getCapacityAssignmentConfiguration = $grpc.ClientMethod<
          $0.GetCapacityAssignmentConfigurationInput,
          $0.GetCapacityAssignmentConfigurationOutput>(
      '/athena.AthenaService/GetCapacityAssignmentConfiguration',
      ($0.GetCapacityAssignmentConfigurationInput value) =>
          value.writeToBuffer(),
      $0.GetCapacityAssignmentConfigurationOutput.fromBuffer);
  static final _$getCapacityReservation = $grpc.ClientMethod<
          $0.GetCapacityReservationInput, $0.GetCapacityReservationOutput>(
      '/athena.AthenaService/GetCapacityReservation',
      ($0.GetCapacityReservationInput value) => value.writeToBuffer(),
      $0.GetCapacityReservationOutput.fromBuffer);
  static final _$getDatabase =
      $grpc.ClientMethod<$0.GetDatabaseInput, $0.GetDatabaseOutput>(
          '/athena.AthenaService/GetDatabase',
          ($0.GetDatabaseInput value) => value.writeToBuffer(),
          $0.GetDatabaseOutput.fromBuffer);
  static final _$getDataCatalog =
      $grpc.ClientMethod<$0.GetDataCatalogInput, $0.GetDataCatalogOutput>(
          '/athena.AthenaService/GetDataCatalog',
          ($0.GetDataCatalogInput value) => value.writeToBuffer(),
          $0.GetDataCatalogOutput.fromBuffer);
  static final _$getNamedQuery =
      $grpc.ClientMethod<$0.GetNamedQueryInput, $0.GetNamedQueryOutput>(
          '/athena.AthenaService/GetNamedQuery',
          ($0.GetNamedQueryInput value) => value.writeToBuffer(),
          $0.GetNamedQueryOutput.fromBuffer);
  static final _$getNotebookMetadata = $grpc.ClientMethod<
          $0.GetNotebookMetadataInput, $0.GetNotebookMetadataOutput>(
      '/athena.AthenaService/GetNotebookMetadata',
      ($0.GetNotebookMetadataInput value) => value.writeToBuffer(),
      $0.GetNotebookMetadataOutput.fromBuffer);
  static final _$getPreparedStatement = $grpc.ClientMethod<
          $0.GetPreparedStatementInput, $0.GetPreparedStatementOutput>(
      '/athena.AthenaService/GetPreparedStatement',
      ($0.GetPreparedStatementInput value) => value.writeToBuffer(),
      $0.GetPreparedStatementOutput.fromBuffer);
  static final _$getQueryExecution =
      $grpc.ClientMethod<$0.GetQueryExecutionInput, $0.GetQueryExecutionOutput>(
          '/athena.AthenaService/GetQueryExecution',
          ($0.GetQueryExecutionInput value) => value.writeToBuffer(),
          $0.GetQueryExecutionOutput.fromBuffer);
  static final _$getQueryResults =
      $grpc.ClientMethod<$0.GetQueryResultsInput, $0.GetQueryResultsOutput>(
          '/athena.AthenaService/GetQueryResults',
          ($0.GetQueryResultsInput value) => value.writeToBuffer(),
          $0.GetQueryResultsOutput.fromBuffer);
  static final _$getQueryRuntimeStatistics = $grpc.ClientMethod<
          $0.GetQueryRuntimeStatisticsInput,
          $0.GetQueryRuntimeStatisticsOutput>(
      '/athena.AthenaService/GetQueryRuntimeStatistics',
      ($0.GetQueryRuntimeStatisticsInput value) => value.writeToBuffer(),
      $0.GetQueryRuntimeStatisticsOutput.fromBuffer);
  static final _$getResourceDashboard = $grpc.ClientMethod<
          $0.GetResourceDashboardRequest, $0.GetResourceDashboardResponse>(
      '/athena.AthenaService/GetResourceDashboard',
      ($0.GetResourceDashboardRequest value) => value.writeToBuffer(),
      $0.GetResourceDashboardResponse.fromBuffer);
  static final _$getSession =
      $grpc.ClientMethod<$0.GetSessionRequest, $0.GetSessionResponse>(
          '/athena.AthenaService/GetSession',
          ($0.GetSessionRequest value) => value.writeToBuffer(),
          $0.GetSessionResponse.fromBuffer);
  static final _$getSessionEndpoint = $grpc.ClientMethod<
          $0.GetSessionEndpointRequest, $0.GetSessionEndpointResponse>(
      '/athena.AthenaService/GetSessionEndpoint',
      ($0.GetSessionEndpointRequest value) => value.writeToBuffer(),
      $0.GetSessionEndpointResponse.fromBuffer);
  static final _$getSessionStatus = $grpc.ClientMethod<
          $0.GetSessionStatusRequest, $0.GetSessionStatusResponse>(
      '/athena.AthenaService/GetSessionStatus',
      ($0.GetSessionStatusRequest value) => value.writeToBuffer(),
      $0.GetSessionStatusResponse.fromBuffer);
  static final _$getTableMetadata =
      $grpc.ClientMethod<$0.GetTableMetadataInput, $0.GetTableMetadataOutput>(
          '/athena.AthenaService/GetTableMetadata',
          ($0.GetTableMetadataInput value) => value.writeToBuffer(),
          $0.GetTableMetadataOutput.fromBuffer);
  static final _$getWorkGroup =
      $grpc.ClientMethod<$0.GetWorkGroupInput, $0.GetWorkGroupOutput>(
          '/athena.AthenaService/GetWorkGroup',
          ($0.GetWorkGroupInput value) => value.writeToBuffer(),
          $0.GetWorkGroupOutput.fromBuffer);
  static final _$importNotebook =
      $grpc.ClientMethod<$0.ImportNotebookInput, $0.ImportNotebookOutput>(
          '/athena.AthenaService/ImportNotebook',
          ($0.ImportNotebookInput value) => value.writeToBuffer(),
          $0.ImportNotebookOutput.fromBuffer);
  static final _$listApplicationDPUSizes = $grpc.ClientMethod<
          $0.ListApplicationDPUSizesInput, $0.ListApplicationDPUSizesOutput>(
      '/athena.AthenaService/ListApplicationDPUSizes',
      ($0.ListApplicationDPUSizesInput value) => value.writeToBuffer(),
      $0.ListApplicationDPUSizesOutput.fromBuffer);
  static final _$listCalculationExecutions = $grpc.ClientMethod<
          $0.ListCalculationExecutionsRequest,
          $0.ListCalculationExecutionsResponse>(
      '/athena.AthenaService/ListCalculationExecutions',
      ($0.ListCalculationExecutionsRequest value) => value.writeToBuffer(),
      $0.ListCalculationExecutionsResponse.fromBuffer);
  static final _$listCapacityReservations = $grpc.ClientMethod<
          $0.ListCapacityReservationsInput, $0.ListCapacityReservationsOutput>(
      '/athena.AthenaService/ListCapacityReservations',
      ($0.ListCapacityReservationsInput value) => value.writeToBuffer(),
      $0.ListCapacityReservationsOutput.fromBuffer);
  static final _$listDatabases =
      $grpc.ClientMethod<$0.ListDatabasesInput, $0.ListDatabasesOutput>(
          '/athena.AthenaService/ListDatabases',
          ($0.ListDatabasesInput value) => value.writeToBuffer(),
          $0.ListDatabasesOutput.fromBuffer);
  static final _$listDataCatalogs =
      $grpc.ClientMethod<$0.ListDataCatalogsInput, $0.ListDataCatalogsOutput>(
          '/athena.AthenaService/ListDataCatalogs',
          ($0.ListDataCatalogsInput value) => value.writeToBuffer(),
          $0.ListDataCatalogsOutput.fromBuffer);
  static final _$listEngineVersions = $grpc.ClientMethod<
          $0.ListEngineVersionsInput, $0.ListEngineVersionsOutput>(
      '/athena.AthenaService/ListEngineVersions',
      ($0.ListEngineVersionsInput value) => value.writeToBuffer(),
      $0.ListEngineVersionsOutput.fromBuffer);
  static final _$listExecutors =
      $grpc.ClientMethod<$0.ListExecutorsRequest, $0.ListExecutorsResponse>(
          '/athena.AthenaService/ListExecutors',
          ($0.ListExecutorsRequest value) => value.writeToBuffer(),
          $0.ListExecutorsResponse.fromBuffer);
  static final _$listNamedQueries =
      $grpc.ClientMethod<$0.ListNamedQueriesInput, $0.ListNamedQueriesOutput>(
          '/athena.AthenaService/ListNamedQueries',
          ($0.ListNamedQueriesInput value) => value.writeToBuffer(),
          $0.ListNamedQueriesOutput.fromBuffer);
  static final _$listNotebookMetadata = $grpc.ClientMethod<
          $0.ListNotebookMetadataInput, $0.ListNotebookMetadataOutput>(
      '/athena.AthenaService/ListNotebookMetadata',
      ($0.ListNotebookMetadataInput value) => value.writeToBuffer(),
      $0.ListNotebookMetadataOutput.fromBuffer);
  static final _$listNotebookSessions = $grpc.ClientMethod<
          $0.ListNotebookSessionsRequest, $0.ListNotebookSessionsResponse>(
      '/athena.AthenaService/ListNotebookSessions',
      ($0.ListNotebookSessionsRequest value) => value.writeToBuffer(),
      $0.ListNotebookSessionsResponse.fromBuffer);
  static final _$listPreparedStatements = $grpc.ClientMethod<
          $0.ListPreparedStatementsInput, $0.ListPreparedStatementsOutput>(
      '/athena.AthenaService/ListPreparedStatements',
      ($0.ListPreparedStatementsInput value) => value.writeToBuffer(),
      $0.ListPreparedStatementsOutput.fromBuffer);
  static final _$listQueryExecutions = $grpc.ClientMethod<
          $0.ListQueryExecutionsInput, $0.ListQueryExecutionsOutput>(
      '/athena.AthenaService/ListQueryExecutions',
      ($0.ListQueryExecutionsInput value) => value.writeToBuffer(),
      $0.ListQueryExecutionsOutput.fromBuffer);
  static final _$listSessions =
      $grpc.ClientMethod<$0.ListSessionsRequest, $0.ListSessionsResponse>(
          '/athena.AthenaService/ListSessions',
          ($0.ListSessionsRequest value) => value.writeToBuffer(),
          $0.ListSessionsResponse.fromBuffer);
  static final _$listTableMetadata =
      $grpc.ClientMethod<$0.ListTableMetadataInput, $0.ListTableMetadataOutput>(
          '/athena.AthenaService/ListTableMetadata',
          ($0.ListTableMetadataInput value) => value.writeToBuffer(),
          $0.ListTableMetadataOutput.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceInput, $0.ListTagsForResourceOutput>(
      '/athena.AthenaService/ListTagsForResource',
      ($0.ListTagsForResourceInput value) => value.writeToBuffer(),
      $0.ListTagsForResourceOutput.fromBuffer);
  static final _$listWorkGroups =
      $grpc.ClientMethod<$0.ListWorkGroupsInput, $0.ListWorkGroupsOutput>(
          '/athena.AthenaService/ListWorkGroups',
          ($0.ListWorkGroupsInput value) => value.writeToBuffer(),
          $0.ListWorkGroupsOutput.fromBuffer);
  static final _$putCapacityAssignmentConfiguration = $grpc.ClientMethod<
          $0.PutCapacityAssignmentConfigurationInput,
          $0.PutCapacityAssignmentConfigurationOutput>(
      '/athena.AthenaService/PutCapacityAssignmentConfiguration',
      ($0.PutCapacityAssignmentConfigurationInput value) =>
          value.writeToBuffer(),
      $0.PutCapacityAssignmentConfigurationOutput.fromBuffer);
  static final _$startCalculationExecution = $grpc.ClientMethod<
          $0.StartCalculationExecutionRequest,
          $0.StartCalculationExecutionResponse>(
      '/athena.AthenaService/StartCalculationExecution',
      ($0.StartCalculationExecutionRequest value) => value.writeToBuffer(),
      $0.StartCalculationExecutionResponse.fromBuffer);
  static final _$startQueryExecution = $grpc.ClientMethod<
          $0.StartQueryExecutionInput, $0.StartQueryExecutionOutput>(
      '/athena.AthenaService/StartQueryExecution',
      ($0.StartQueryExecutionInput value) => value.writeToBuffer(),
      $0.StartQueryExecutionOutput.fromBuffer);
  static final _$startSession =
      $grpc.ClientMethod<$0.StartSessionRequest, $0.StartSessionResponse>(
          '/athena.AthenaService/StartSession',
          ($0.StartSessionRequest value) => value.writeToBuffer(),
          $0.StartSessionResponse.fromBuffer);
  static final _$stopCalculationExecution = $grpc.ClientMethod<
          $0.StopCalculationExecutionRequest,
          $0.StopCalculationExecutionResponse>(
      '/athena.AthenaService/StopCalculationExecution',
      ($0.StopCalculationExecutionRequest value) => value.writeToBuffer(),
      $0.StopCalculationExecutionResponse.fromBuffer);
  static final _$stopQueryExecution = $grpc.ClientMethod<
          $0.StopQueryExecutionInput, $0.StopQueryExecutionOutput>(
      '/athena.AthenaService/StopQueryExecution',
      ($0.StopQueryExecutionInput value) => value.writeToBuffer(),
      $0.StopQueryExecutionOutput.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceInput, $0.TagResourceOutput>(
          '/athena.AthenaService/TagResource',
          ($0.TagResourceInput value) => value.writeToBuffer(),
          $0.TagResourceOutput.fromBuffer);
  static final _$terminateSession = $grpc.ClientMethod<
          $0.TerminateSessionRequest, $0.TerminateSessionResponse>(
      '/athena.AthenaService/TerminateSession',
      ($0.TerminateSessionRequest value) => value.writeToBuffer(),
      $0.TerminateSessionResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
          '/athena.AthenaService/UntagResource',
          ($0.UntagResourceInput value) => value.writeToBuffer(),
          $0.UntagResourceOutput.fromBuffer);
  static final _$updateCapacityReservation = $grpc.ClientMethod<
          $0.UpdateCapacityReservationInput,
          $0.UpdateCapacityReservationOutput>(
      '/athena.AthenaService/UpdateCapacityReservation',
      ($0.UpdateCapacityReservationInput value) => value.writeToBuffer(),
      $0.UpdateCapacityReservationOutput.fromBuffer);
  static final _$updateDataCatalog =
      $grpc.ClientMethod<$0.UpdateDataCatalogInput, $0.UpdateDataCatalogOutput>(
          '/athena.AthenaService/UpdateDataCatalog',
          ($0.UpdateDataCatalogInput value) => value.writeToBuffer(),
          $0.UpdateDataCatalogOutput.fromBuffer);
  static final _$updateNamedQuery =
      $grpc.ClientMethod<$0.UpdateNamedQueryInput, $0.UpdateNamedQueryOutput>(
          '/athena.AthenaService/UpdateNamedQuery',
          ($0.UpdateNamedQueryInput value) => value.writeToBuffer(),
          $0.UpdateNamedQueryOutput.fromBuffer);
  static final _$updateNotebook =
      $grpc.ClientMethod<$0.UpdateNotebookInput, $0.UpdateNotebookOutput>(
          '/athena.AthenaService/UpdateNotebook',
          ($0.UpdateNotebookInput value) => value.writeToBuffer(),
          $0.UpdateNotebookOutput.fromBuffer);
  static final _$updateNotebookMetadata = $grpc.ClientMethod<
          $0.UpdateNotebookMetadataInput, $0.UpdateNotebookMetadataOutput>(
      '/athena.AthenaService/UpdateNotebookMetadata',
      ($0.UpdateNotebookMetadataInput value) => value.writeToBuffer(),
      $0.UpdateNotebookMetadataOutput.fromBuffer);
  static final _$updatePreparedStatement = $grpc.ClientMethod<
          $0.UpdatePreparedStatementInput, $0.UpdatePreparedStatementOutput>(
      '/athena.AthenaService/UpdatePreparedStatement',
      ($0.UpdatePreparedStatementInput value) => value.writeToBuffer(),
      $0.UpdatePreparedStatementOutput.fromBuffer);
  static final _$updateWorkGroup =
      $grpc.ClientMethod<$0.UpdateWorkGroupInput, $0.UpdateWorkGroupOutput>(
          '/athena.AthenaService/UpdateWorkGroup',
          ($0.UpdateWorkGroupInput value) => value.writeToBuffer(),
          $0.UpdateWorkGroupOutput.fromBuffer);
}

@$pb.GrpcServiceName('athena.AthenaService')
abstract class AthenaServiceBase extends $grpc.Service {
  $core.String get $name => 'athena.AthenaService';

  AthenaServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.BatchGetNamedQueryInput,
            $0.BatchGetNamedQueryOutput>(
        'BatchGetNamedQuery',
        batchGetNamedQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.BatchGetNamedQueryInput.fromBuffer(value),
        ($0.BatchGetNamedQueryOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.BatchGetPreparedStatementInput,
            $0.BatchGetPreparedStatementOutput>(
        'BatchGetPreparedStatement',
        batchGetPreparedStatement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.BatchGetPreparedStatementInput.fromBuffer(value),
        ($0.BatchGetPreparedStatementOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.BatchGetQueryExecutionInput,
            $0.BatchGetQueryExecutionOutput>(
        'BatchGetQueryExecution',
        batchGetQueryExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.BatchGetQueryExecutionInput.fromBuffer(value),
        ($0.BatchGetQueryExecutionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelCapacityReservationInput,
            $0.CancelCapacityReservationOutput>(
        'CancelCapacityReservation',
        cancelCapacityReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelCapacityReservationInput.fromBuffer(value),
        ($0.CancelCapacityReservationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateCapacityReservationInput,
            $0.CreateCapacityReservationOutput>(
        'CreateCapacityReservation',
        createCapacityReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCapacityReservationInput.fromBuffer(value),
        ($0.CreateCapacityReservationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDataCatalogInput,
            $0.CreateDataCatalogOutput>(
        'CreateDataCatalog',
        createDataCatalog_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDataCatalogInput.fromBuffer(value),
        ($0.CreateDataCatalogOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateNamedQueryInput,
            $0.CreateNamedQueryOutput>(
        'CreateNamedQuery',
        createNamedQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateNamedQueryInput.fromBuffer(value),
        ($0.CreateNamedQueryOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateNotebookInput, $0.CreateNotebookOutput>(
            'CreateNotebook',
            createNotebook_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateNotebookInput.fromBuffer(value),
            ($0.CreateNotebookOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePreparedStatementInput,
            $0.CreatePreparedStatementOutput>(
        'CreatePreparedStatement',
        createPreparedStatement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePreparedStatementInput.fromBuffer(value),
        ($0.CreatePreparedStatementOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePresignedNotebookUrlRequest,
            $0.CreatePresignedNotebookUrlResponse>(
        'CreatePresignedNotebookUrl',
        createPresignedNotebookUrl_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePresignedNotebookUrlRequest.fromBuffer(value),
        ($0.CreatePresignedNotebookUrlResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateWorkGroupInput, $0.CreateWorkGroupOutput>(
            'CreateWorkGroup',
            createWorkGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateWorkGroupInput.fromBuffer(value),
            ($0.CreateWorkGroupOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteCapacityReservationInput,
            $0.DeleteCapacityReservationOutput>(
        'DeleteCapacityReservation',
        deleteCapacityReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCapacityReservationInput.fromBuffer(value),
        ($0.DeleteCapacityReservationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteDataCatalogInput,
            $0.DeleteDataCatalogOutput>(
        'DeleteDataCatalog',
        deleteDataCatalog_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteDataCatalogInput.fromBuffer(value),
        ($0.DeleteDataCatalogOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteNamedQueryInput,
            $0.DeleteNamedQueryOutput>(
        'DeleteNamedQuery',
        deleteNamedQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteNamedQueryInput.fromBuffer(value),
        ($0.DeleteNamedQueryOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteNotebookInput, $0.DeleteNotebookOutput>(
            'DeleteNotebook',
            deleteNotebook_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteNotebookInput.fromBuffer(value),
            ($0.DeleteNotebookOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePreparedStatementInput,
            $0.DeletePreparedStatementOutput>(
        'DeletePreparedStatement',
        deletePreparedStatement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePreparedStatementInput.fromBuffer(value),
        ($0.DeletePreparedStatementOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteWorkGroupInput, $0.DeleteWorkGroupOutput>(
            'DeleteWorkGroup',
            deleteWorkGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteWorkGroupInput.fromBuffer(value),
            ($0.DeleteWorkGroupOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ExportNotebookInput, $0.ExportNotebookOutput>(
            'ExportNotebook',
            exportNotebook_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ExportNotebookInput.fromBuffer(value),
            ($0.ExportNotebookOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCalculationExecutionRequest,
            $0.GetCalculationExecutionResponse>(
        'GetCalculationExecution',
        getCalculationExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCalculationExecutionRequest.fromBuffer(value),
        ($0.GetCalculationExecutionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCalculationExecutionCodeRequest,
            $0.GetCalculationExecutionCodeResponse>(
        'GetCalculationExecutionCode',
        getCalculationExecutionCode_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCalculationExecutionCodeRequest.fromBuffer(value),
        ($0.GetCalculationExecutionCodeResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCalculationExecutionStatusRequest,
            $0.GetCalculationExecutionStatusResponse>(
        'GetCalculationExecutionStatus',
        getCalculationExecutionStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCalculationExecutionStatusRequest.fromBuffer(value),
        ($0.GetCalculationExecutionStatusResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCapacityAssignmentConfigurationInput,
            $0.GetCapacityAssignmentConfigurationOutput>(
        'GetCapacityAssignmentConfiguration',
        getCapacityAssignmentConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCapacityAssignmentConfigurationInput.fromBuffer(value),
        ($0.GetCapacityAssignmentConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCapacityReservationInput,
            $0.GetCapacityReservationOutput>(
        'GetCapacityReservation',
        getCapacityReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCapacityReservationInput.fromBuffer(value),
        ($0.GetCapacityReservationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDatabaseInput, $0.GetDatabaseOutput>(
        'GetDatabase',
        getDatabase_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetDatabaseInput.fromBuffer(value),
        ($0.GetDatabaseOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetDataCatalogInput, $0.GetDataCatalogOutput>(
            'GetDataCatalog',
            getDataCatalog_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetDataCatalogInput.fromBuffer(value),
            ($0.GetDataCatalogOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetNamedQueryInput, $0.GetNamedQueryOutput>(
            'GetNamedQuery',
            getNamedQuery_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetNamedQueryInput.fromBuffer(value),
            ($0.GetNamedQueryOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetNotebookMetadataInput,
            $0.GetNotebookMetadataOutput>(
        'GetNotebookMetadata',
        getNotebookMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetNotebookMetadataInput.fromBuffer(value),
        ($0.GetNotebookMetadataOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPreparedStatementInput,
            $0.GetPreparedStatementOutput>(
        'GetPreparedStatement',
        getPreparedStatement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPreparedStatementInput.fromBuffer(value),
        ($0.GetPreparedStatementOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetQueryExecutionInput,
            $0.GetQueryExecutionOutput>(
        'GetQueryExecution',
        getQueryExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetQueryExecutionInput.fromBuffer(value),
        ($0.GetQueryExecutionOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetQueryResultsInput, $0.GetQueryResultsOutput>(
            'GetQueryResults',
            getQueryResults_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetQueryResultsInput.fromBuffer(value),
            ($0.GetQueryResultsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetQueryRuntimeStatisticsInput,
            $0.GetQueryRuntimeStatisticsOutput>(
        'GetQueryRuntimeStatistics',
        getQueryRuntimeStatistics_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetQueryRuntimeStatisticsInput.fromBuffer(value),
        ($0.GetQueryRuntimeStatisticsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourceDashboardRequest,
            $0.GetResourceDashboardResponse>(
        'GetResourceDashboard',
        getResourceDashboard_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourceDashboardRequest.fromBuffer(value),
        ($0.GetResourceDashboardResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSessionRequest, $0.GetSessionResponse>(
        'GetSession',
        getSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetSessionRequest.fromBuffer(value),
        ($0.GetSessionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSessionEndpointRequest,
            $0.GetSessionEndpointResponse>(
        'GetSessionEndpoint',
        getSessionEndpoint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSessionEndpointRequest.fromBuffer(value),
        ($0.GetSessionEndpointResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSessionStatusRequest,
            $0.GetSessionStatusResponse>(
        'GetSessionStatus',
        getSessionStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSessionStatusRequest.fromBuffer(value),
        ($0.GetSessionStatusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTableMetadataInput,
            $0.GetTableMetadataOutput>(
        'GetTableMetadata',
        getTableMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTableMetadataInput.fromBuffer(value),
        ($0.GetTableMetadataOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetWorkGroupInput, $0.GetWorkGroupOutput>(
        'GetWorkGroup',
        getWorkGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetWorkGroupInput.fromBuffer(value),
        ($0.GetWorkGroupOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ImportNotebookInput, $0.ImportNotebookOutput>(
            'ImportNotebook',
            importNotebook_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ImportNotebookInput.fromBuffer(value),
            ($0.ImportNotebookOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListApplicationDPUSizesInput,
            $0.ListApplicationDPUSizesOutput>(
        'ListApplicationDPUSizes',
        listApplicationDPUSizes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListApplicationDPUSizesInput.fromBuffer(value),
        ($0.ListApplicationDPUSizesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCalculationExecutionsRequest,
            $0.ListCalculationExecutionsResponse>(
        'ListCalculationExecutions',
        listCalculationExecutions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCalculationExecutionsRequest.fromBuffer(value),
        ($0.ListCalculationExecutionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCapacityReservationsInput,
            $0.ListCapacityReservationsOutput>(
        'ListCapacityReservations',
        listCapacityReservations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCapacityReservationsInput.fromBuffer(value),
        ($0.ListCapacityReservationsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListDatabasesInput, $0.ListDatabasesOutput>(
            'ListDatabases',
            listDatabases_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDatabasesInput.fromBuffer(value),
            ($0.ListDatabasesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDataCatalogsInput,
            $0.ListDataCatalogsOutput>(
        'ListDataCatalogs',
        listDataCatalogs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDataCatalogsInput.fromBuffer(value),
        ($0.ListDataCatalogsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEngineVersionsInput,
            $0.ListEngineVersionsOutput>(
        'ListEngineVersions',
        listEngineVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEngineVersionsInput.fromBuffer(value),
        ($0.ListEngineVersionsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListExecutorsRequest, $0.ListExecutorsResponse>(
            'ListExecutors',
            listExecutors_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListExecutorsRequest.fromBuffer(value),
            ($0.ListExecutorsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListNamedQueriesInput,
            $0.ListNamedQueriesOutput>(
        'ListNamedQueries',
        listNamedQueries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListNamedQueriesInput.fromBuffer(value),
        ($0.ListNamedQueriesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListNotebookMetadataInput,
            $0.ListNotebookMetadataOutput>(
        'ListNotebookMetadata',
        listNotebookMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListNotebookMetadataInput.fromBuffer(value),
        ($0.ListNotebookMetadataOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListNotebookSessionsRequest,
            $0.ListNotebookSessionsResponse>(
        'ListNotebookSessions',
        listNotebookSessions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListNotebookSessionsRequest.fromBuffer(value),
        ($0.ListNotebookSessionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPreparedStatementsInput,
            $0.ListPreparedStatementsOutput>(
        'ListPreparedStatements',
        listPreparedStatements_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPreparedStatementsInput.fromBuffer(value),
        ($0.ListPreparedStatementsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListQueryExecutionsInput,
            $0.ListQueryExecutionsOutput>(
        'ListQueryExecutions',
        listQueryExecutions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListQueryExecutionsInput.fromBuffer(value),
        ($0.ListQueryExecutionsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListSessionsRequest, $0.ListSessionsResponse>(
            'ListSessions',
            listSessions_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListSessionsRequest.fromBuffer(value),
            ($0.ListSessionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTableMetadataInput,
            $0.ListTableMetadataOutput>(
        'ListTableMetadata',
        listTableMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTableMetadataInput.fromBuffer(value),
        ($0.ListTableMetadataOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceInput,
            $0.ListTagsForResourceOutput>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceInput.fromBuffer(value),
        ($0.ListTagsForResourceOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListWorkGroupsInput, $0.ListWorkGroupsOutput>(
            'ListWorkGroups',
            listWorkGroups_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListWorkGroupsInput.fromBuffer(value),
            ($0.ListWorkGroupsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutCapacityAssignmentConfigurationInput,
            $0.PutCapacityAssignmentConfigurationOutput>(
        'PutCapacityAssignmentConfiguration',
        putCapacityAssignmentConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutCapacityAssignmentConfigurationInput.fromBuffer(value),
        ($0.PutCapacityAssignmentConfigurationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartCalculationExecutionRequest,
            $0.StartCalculationExecutionResponse>(
        'StartCalculationExecution',
        startCalculationExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartCalculationExecutionRequest.fromBuffer(value),
        ($0.StartCalculationExecutionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartQueryExecutionInput,
            $0.StartQueryExecutionOutput>(
        'StartQueryExecution',
        startQueryExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartQueryExecutionInput.fromBuffer(value),
        ($0.StartQueryExecutionOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartSessionRequest, $0.StartSessionResponse>(
            'StartSession',
            startSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartSessionRequest.fromBuffer(value),
            ($0.StartSessionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopCalculationExecutionRequest,
            $0.StopCalculationExecutionResponse>(
        'StopCalculationExecution',
        stopCalculationExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopCalculationExecutionRequest.fromBuffer(value),
        ($0.StopCalculationExecutionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopQueryExecutionInput,
            $0.StopQueryExecutionOutput>(
        'StopQueryExecution',
        stopQueryExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopQueryExecutionInput.fromBuffer(value),
        ($0.StopQueryExecutionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceInput, $0.TagResourceOutput>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagResourceInput.fromBuffer(value),
        ($0.TagResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TerminateSessionRequest,
            $0.TerminateSessionResponse>(
        'TerminateSession',
        terminateSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TerminateSessionRequest.fromBuffer(value),
        ($0.TerminateSessionResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceInput, $0.UntagResourceOutput>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceInput.fromBuffer(value),
            ($0.UntagResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateCapacityReservationInput,
            $0.UpdateCapacityReservationOutput>(
        'UpdateCapacityReservation',
        updateCapacityReservation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateCapacityReservationInput.fromBuffer(value),
        ($0.UpdateCapacityReservationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDataCatalogInput,
            $0.UpdateDataCatalogOutput>(
        'UpdateDataCatalog',
        updateDataCatalog_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDataCatalogInput.fromBuffer(value),
        ($0.UpdateDataCatalogOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateNamedQueryInput,
            $0.UpdateNamedQueryOutput>(
        'UpdateNamedQuery',
        updateNamedQuery_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateNamedQueryInput.fromBuffer(value),
        ($0.UpdateNamedQueryOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateNotebookInput, $0.UpdateNotebookOutput>(
            'UpdateNotebook',
            updateNotebook_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateNotebookInput.fromBuffer(value),
            ($0.UpdateNotebookOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateNotebookMetadataInput,
            $0.UpdateNotebookMetadataOutput>(
        'UpdateNotebookMetadata',
        updateNotebookMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateNotebookMetadataInput.fromBuffer(value),
        ($0.UpdateNotebookMetadataOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdatePreparedStatementInput,
            $0.UpdatePreparedStatementOutput>(
        'UpdatePreparedStatement',
        updatePreparedStatement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdatePreparedStatementInput.fromBuffer(value),
        ($0.UpdatePreparedStatementOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateWorkGroupInput, $0.UpdateWorkGroupOutput>(
            'UpdateWorkGroup',
            updateWorkGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateWorkGroupInput.fromBuffer(value),
            ($0.UpdateWorkGroupOutput value) => value.writeToBuffer()));
  }

  $async.Future<$0.BatchGetNamedQueryOutput> batchGetNamedQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.BatchGetNamedQueryInput> $request) async {
    return batchGetNamedQuery($call, await $request);
  }

  $async.Future<$0.BatchGetNamedQueryOutput> batchGetNamedQuery(
      $grpc.ServiceCall call, $0.BatchGetNamedQueryInput request);

  $async.Future<$0.BatchGetPreparedStatementOutput>
      batchGetPreparedStatement_Pre($grpc.ServiceCall $call,
          $async.Future<$0.BatchGetPreparedStatementInput> $request) async {
    return batchGetPreparedStatement($call, await $request);
  }

  $async.Future<$0.BatchGetPreparedStatementOutput> batchGetPreparedStatement(
      $grpc.ServiceCall call, $0.BatchGetPreparedStatementInput request);

  $async.Future<$0.BatchGetQueryExecutionOutput> batchGetQueryExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.BatchGetQueryExecutionInput> $request) async {
    return batchGetQueryExecution($call, await $request);
  }

  $async.Future<$0.BatchGetQueryExecutionOutput> batchGetQueryExecution(
      $grpc.ServiceCall call, $0.BatchGetQueryExecutionInput request);

  $async.Future<$0.CancelCapacityReservationOutput>
      cancelCapacityReservation_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CancelCapacityReservationInput> $request) async {
    return cancelCapacityReservation($call, await $request);
  }

  $async.Future<$0.CancelCapacityReservationOutput> cancelCapacityReservation(
      $grpc.ServiceCall call, $0.CancelCapacityReservationInput request);

  $async.Future<$0.CreateCapacityReservationOutput>
      createCapacityReservation_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateCapacityReservationInput> $request) async {
    return createCapacityReservation($call, await $request);
  }

  $async.Future<$0.CreateCapacityReservationOutput> createCapacityReservation(
      $grpc.ServiceCall call, $0.CreateCapacityReservationInput request);

  $async.Future<$0.CreateDataCatalogOutput> createDataCatalog_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDataCatalogInput> $request) async {
    return createDataCatalog($call, await $request);
  }

  $async.Future<$0.CreateDataCatalogOutput> createDataCatalog(
      $grpc.ServiceCall call, $0.CreateDataCatalogInput request);

  $async.Future<$0.CreateNamedQueryOutput> createNamedQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateNamedQueryInput> $request) async {
    return createNamedQuery($call, await $request);
  }

  $async.Future<$0.CreateNamedQueryOutput> createNamedQuery(
      $grpc.ServiceCall call, $0.CreateNamedQueryInput request);

  $async.Future<$0.CreateNotebookOutput> createNotebook_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateNotebookInput> $request) async {
    return createNotebook($call, await $request);
  }

  $async.Future<$0.CreateNotebookOutput> createNotebook(
      $grpc.ServiceCall call, $0.CreateNotebookInput request);

  $async.Future<$0.CreatePreparedStatementOutput> createPreparedStatement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreatePreparedStatementInput> $request) async {
    return createPreparedStatement($call, await $request);
  }

  $async.Future<$0.CreatePreparedStatementOutput> createPreparedStatement(
      $grpc.ServiceCall call, $0.CreatePreparedStatementInput request);

  $async.Future<$0.CreatePresignedNotebookUrlResponse>
      createPresignedNotebookUrl_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreatePresignedNotebookUrlRequest> $request) async {
    return createPresignedNotebookUrl($call, await $request);
  }

  $async.Future<$0.CreatePresignedNotebookUrlResponse>
      createPresignedNotebookUrl(
          $grpc.ServiceCall call, $0.CreatePresignedNotebookUrlRequest request);

  $async.Future<$0.CreateWorkGroupOutput> createWorkGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateWorkGroupInput> $request) async {
    return createWorkGroup($call, await $request);
  }

  $async.Future<$0.CreateWorkGroupOutput> createWorkGroup(
      $grpc.ServiceCall call, $0.CreateWorkGroupInput request);

  $async.Future<$0.DeleteCapacityReservationOutput>
      deleteCapacityReservation_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteCapacityReservationInput> $request) async {
    return deleteCapacityReservation($call, await $request);
  }

  $async.Future<$0.DeleteCapacityReservationOutput> deleteCapacityReservation(
      $grpc.ServiceCall call, $0.DeleteCapacityReservationInput request);

  $async.Future<$0.DeleteDataCatalogOutput> deleteDataCatalog_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDataCatalogInput> $request) async {
    return deleteDataCatalog($call, await $request);
  }

  $async.Future<$0.DeleteDataCatalogOutput> deleteDataCatalog(
      $grpc.ServiceCall call, $0.DeleteDataCatalogInput request);

  $async.Future<$0.DeleteNamedQueryOutput> deleteNamedQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteNamedQueryInput> $request) async {
    return deleteNamedQuery($call, await $request);
  }

  $async.Future<$0.DeleteNamedQueryOutput> deleteNamedQuery(
      $grpc.ServiceCall call, $0.DeleteNamedQueryInput request);

  $async.Future<$0.DeleteNotebookOutput> deleteNotebook_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteNotebookInput> $request) async {
    return deleteNotebook($call, await $request);
  }

  $async.Future<$0.DeleteNotebookOutput> deleteNotebook(
      $grpc.ServiceCall call, $0.DeleteNotebookInput request);

  $async.Future<$0.DeletePreparedStatementOutput> deletePreparedStatement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeletePreparedStatementInput> $request) async {
    return deletePreparedStatement($call, await $request);
  }

  $async.Future<$0.DeletePreparedStatementOutput> deletePreparedStatement(
      $grpc.ServiceCall call, $0.DeletePreparedStatementInput request);

  $async.Future<$0.DeleteWorkGroupOutput> deleteWorkGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteWorkGroupInput> $request) async {
    return deleteWorkGroup($call, await $request);
  }

  $async.Future<$0.DeleteWorkGroupOutput> deleteWorkGroup(
      $grpc.ServiceCall call, $0.DeleteWorkGroupInput request);

  $async.Future<$0.ExportNotebookOutput> exportNotebook_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExportNotebookInput> $request) async {
    return exportNotebook($call, await $request);
  }

  $async.Future<$0.ExportNotebookOutput> exportNotebook(
      $grpc.ServiceCall call, $0.ExportNotebookInput request);

  $async.Future<$0.GetCalculationExecutionResponse> getCalculationExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCalculationExecutionRequest> $request) async {
    return getCalculationExecution($call, await $request);
  }

  $async.Future<$0.GetCalculationExecutionResponse> getCalculationExecution(
      $grpc.ServiceCall call, $0.GetCalculationExecutionRequest request);

  $async.Future<$0.GetCalculationExecutionCodeResponse>
      getCalculationExecutionCode_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetCalculationExecutionCodeRequest> $request) async {
    return getCalculationExecutionCode($call, await $request);
  }

  $async.Future<$0.GetCalculationExecutionCodeResponse>
      getCalculationExecutionCode($grpc.ServiceCall call,
          $0.GetCalculationExecutionCodeRequest request);

  $async.Future<$0.GetCalculationExecutionStatusResponse>
      getCalculationExecutionStatus_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetCalculationExecutionStatusRequest>
              $request) async {
    return getCalculationExecutionStatus($call, await $request);
  }

  $async.Future<$0.GetCalculationExecutionStatusResponse>
      getCalculationExecutionStatus($grpc.ServiceCall call,
          $0.GetCalculationExecutionStatusRequest request);

  $async.Future<$0.GetCapacityAssignmentConfigurationOutput>
      getCapacityAssignmentConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetCapacityAssignmentConfigurationInput>
              $request) async {
    return getCapacityAssignmentConfiguration($call, await $request);
  }

  $async.Future<$0.GetCapacityAssignmentConfigurationOutput>
      getCapacityAssignmentConfiguration($grpc.ServiceCall call,
          $0.GetCapacityAssignmentConfigurationInput request);

  $async.Future<$0.GetCapacityReservationOutput> getCapacityReservation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCapacityReservationInput> $request) async {
    return getCapacityReservation($call, await $request);
  }

  $async.Future<$0.GetCapacityReservationOutput> getCapacityReservation(
      $grpc.ServiceCall call, $0.GetCapacityReservationInput request);

  $async.Future<$0.GetDatabaseOutput> getDatabase_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDatabaseInput> $request) async {
    return getDatabase($call, await $request);
  }

  $async.Future<$0.GetDatabaseOutput> getDatabase(
      $grpc.ServiceCall call, $0.GetDatabaseInput request);

  $async.Future<$0.GetDataCatalogOutput> getDataCatalog_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDataCatalogInput> $request) async {
    return getDataCatalog($call, await $request);
  }

  $async.Future<$0.GetDataCatalogOutput> getDataCatalog(
      $grpc.ServiceCall call, $0.GetDataCatalogInput request);

  $async.Future<$0.GetNamedQueryOutput> getNamedQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetNamedQueryInput> $request) async {
    return getNamedQuery($call, await $request);
  }

  $async.Future<$0.GetNamedQueryOutput> getNamedQuery(
      $grpc.ServiceCall call, $0.GetNamedQueryInput request);

  $async.Future<$0.GetNotebookMetadataOutput> getNotebookMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetNotebookMetadataInput> $request) async {
    return getNotebookMetadata($call, await $request);
  }

  $async.Future<$0.GetNotebookMetadataOutput> getNotebookMetadata(
      $grpc.ServiceCall call, $0.GetNotebookMetadataInput request);

  $async.Future<$0.GetPreparedStatementOutput> getPreparedStatement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPreparedStatementInput> $request) async {
    return getPreparedStatement($call, await $request);
  }

  $async.Future<$0.GetPreparedStatementOutput> getPreparedStatement(
      $grpc.ServiceCall call, $0.GetPreparedStatementInput request);

  $async.Future<$0.GetQueryExecutionOutput> getQueryExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetQueryExecutionInput> $request) async {
    return getQueryExecution($call, await $request);
  }

  $async.Future<$0.GetQueryExecutionOutput> getQueryExecution(
      $grpc.ServiceCall call, $0.GetQueryExecutionInput request);

  $async.Future<$0.GetQueryResultsOutput> getQueryResults_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetQueryResultsInput> $request) async {
    return getQueryResults($call, await $request);
  }

  $async.Future<$0.GetQueryResultsOutput> getQueryResults(
      $grpc.ServiceCall call, $0.GetQueryResultsInput request);

  $async.Future<$0.GetQueryRuntimeStatisticsOutput>
      getQueryRuntimeStatistics_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetQueryRuntimeStatisticsInput> $request) async {
    return getQueryRuntimeStatistics($call, await $request);
  }

  $async.Future<$0.GetQueryRuntimeStatisticsOutput> getQueryRuntimeStatistics(
      $grpc.ServiceCall call, $0.GetQueryRuntimeStatisticsInput request);

  $async.Future<$0.GetResourceDashboardResponse> getResourceDashboard_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourceDashboardRequest> $request) async {
    return getResourceDashboard($call, await $request);
  }

  $async.Future<$0.GetResourceDashboardResponse> getResourceDashboard(
      $grpc.ServiceCall call, $0.GetResourceDashboardRequest request);

  $async.Future<$0.GetSessionResponse> getSession_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetSessionRequest> $request) async {
    return getSession($call, await $request);
  }

  $async.Future<$0.GetSessionResponse> getSession(
      $grpc.ServiceCall call, $0.GetSessionRequest request);

  $async.Future<$0.GetSessionEndpointResponse> getSessionEndpoint_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSessionEndpointRequest> $request) async {
    return getSessionEndpoint($call, await $request);
  }

  $async.Future<$0.GetSessionEndpointResponse> getSessionEndpoint(
      $grpc.ServiceCall call, $0.GetSessionEndpointRequest request);

  $async.Future<$0.GetSessionStatusResponse> getSessionStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSessionStatusRequest> $request) async {
    return getSessionStatus($call, await $request);
  }

  $async.Future<$0.GetSessionStatusResponse> getSessionStatus(
      $grpc.ServiceCall call, $0.GetSessionStatusRequest request);

  $async.Future<$0.GetTableMetadataOutput> getTableMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetTableMetadataInput> $request) async {
    return getTableMetadata($call, await $request);
  }

  $async.Future<$0.GetTableMetadataOutput> getTableMetadata(
      $grpc.ServiceCall call, $0.GetTableMetadataInput request);

  $async.Future<$0.GetWorkGroupOutput> getWorkGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetWorkGroupInput> $request) async {
    return getWorkGroup($call, await $request);
  }

  $async.Future<$0.GetWorkGroupOutput> getWorkGroup(
      $grpc.ServiceCall call, $0.GetWorkGroupInput request);

  $async.Future<$0.ImportNotebookOutput> importNotebook_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ImportNotebookInput> $request) async {
    return importNotebook($call, await $request);
  }

  $async.Future<$0.ImportNotebookOutput> importNotebook(
      $grpc.ServiceCall call, $0.ImportNotebookInput request);

  $async.Future<$0.ListApplicationDPUSizesOutput> listApplicationDPUSizes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListApplicationDPUSizesInput> $request) async {
    return listApplicationDPUSizes($call, await $request);
  }

  $async.Future<$0.ListApplicationDPUSizesOutput> listApplicationDPUSizes(
      $grpc.ServiceCall call, $0.ListApplicationDPUSizesInput request);

  $async.Future<$0.ListCalculationExecutionsResponse>
      listCalculationExecutions_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListCalculationExecutionsRequest> $request) async {
    return listCalculationExecutions($call, await $request);
  }

  $async.Future<$0.ListCalculationExecutionsResponse> listCalculationExecutions(
      $grpc.ServiceCall call, $0.ListCalculationExecutionsRequest request);

  $async.Future<$0.ListCapacityReservationsOutput> listCapacityReservations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCapacityReservationsInput> $request) async {
    return listCapacityReservations($call, await $request);
  }

  $async.Future<$0.ListCapacityReservationsOutput> listCapacityReservations(
      $grpc.ServiceCall call, $0.ListCapacityReservationsInput request);

  $async.Future<$0.ListDatabasesOutput> listDatabases_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDatabasesInput> $request) async {
    return listDatabases($call, await $request);
  }

  $async.Future<$0.ListDatabasesOutput> listDatabases(
      $grpc.ServiceCall call, $0.ListDatabasesInput request);

  $async.Future<$0.ListDataCatalogsOutput> listDataCatalogs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDataCatalogsInput> $request) async {
    return listDataCatalogs($call, await $request);
  }

  $async.Future<$0.ListDataCatalogsOutput> listDataCatalogs(
      $grpc.ServiceCall call, $0.ListDataCatalogsInput request);

  $async.Future<$0.ListEngineVersionsOutput> listEngineVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEngineVersionsInput> $request) async {
    return listEngineVersions($call, await $request);
  }

  $async.Future<$0.ListEngineVersionsOutput> listEngineVersions(
      $grpc.ServiceCall call, $0.ListEngineVersionsInput request);

  $async.Future<$0.ListExecutorsResponse> listExecutors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListExecutorsRequest> $request) async {
    return listExecutors($call, await $request);
  }

  $async.Future<$0.ListExecutorsResponse> listExecutors(
      $grpc.ServiceCall call, $0.ListExecutorsRequest request);

  $async.Future<$0.ListNamedQueriesOutput> listNamedQueries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListNamedQueriesInput> $request) async {
    return listNamedQueries($call, await $request);
  }

  $async.Future<$0.ListNamedQueriesOutput> listNamedQueries(
      $grpc.ServiceCall call, $0.ListNamedQueriesInput request);

  $async.Future<$0.ListNotebookMetadataOutput> listNotebookMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListNotebookMetadataInput> $request) async {
    return listNotebookMetadata($call, await $request);
  }

  $async.Future<$0.ListNotebookMetadataOutput> listNotebookMetadata(
      $grpc.ServiceCall call, $0.ListNotebookMetadataInput request);

  $async.Future<$0.ListNotebookSessionsResponse> listNotebookSessions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListNotebookSessionsRequest> $request) async {
    return listNotebookSessions($call, await $request);
  }

  $async.Future<$0.ListNotebookSessionsResponse> listNotebookSessions(
      $grpc.ServiceCall call, $0.ListNotebookSessionsRequest request);

  $async.Future<$0.ListPreparedStatementsOutput> listPreparedStatements_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPreparedStatementsInput> $request) async {
    return listPreparedStatements($call, await $request);
  }

  $async.Future<$0.ListPreparedStatementsOutput> listPreparedStatements(
      $grpc.ServiceCall call, $0.ListPreparedStatementsInput request);

  $async.Future<$0.ListQueryExecutionsOutput> listQueryExecutions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListQueryExecutionsInput> $request) async {
    return listQueryExecutions($call, await $request);
  }

  $async.Future<$0.ListQueryExecutionsOutput> listQueryExecutions(
      $grpc.ServiceCall call, $0.ListQueryExecutionsInput request);

  $async.Future<$0.ListSessionsResponse> listSessions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSessionsRequest> $request) async {
    return listSessions($call, await $request);
  }

  $async.Future<$0.ListSessionsResponse> listSessions(
      $grpc.ServiceCall call, $0.ListSessionsRequest request);

  $async.Future<$0.ListTableMetadataOutput> listTableMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTableMetadataInput> $request) async {
    return listTableMetadata($call, await $request);
  }

  $async.Future<$0.ListTableMetadataOutput> listTableMetadata(
      $grpc.ServiceCall call, $0.ListTableMetadataInput request);

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceInput> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceOutput> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceInput request);

  $async.Future<$0.ListWorkGroupsOutput> listWorkGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListWorkGroupsInput> $request) async {
    return listWorkGroups($call, await $request);
  }

  $async.Future<$0.ListWorkGroupsOutput> listWorkGroups(
      $grpc.ServiceCall call, $0.ListWorkGroupsInput request);

  $async.Future<$0.PutCapacityAssignmentConfigurationOutput>
      putCapacityAssignmentConfiguration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutCapacityAssignmentConfigurationInput>
              $request) async {
    return putCapacityAssignmentConfiguration($call, await $request);
  }

  $async.Future<$0.PutCapacityAssignmentConfigurationOutput>
      putCapacityAssignmentConfiguration($grpc.ServiceCall call,
          $0.PutCapacityAssignmentConfigurationInput request);

  $async.Future<$0.StartCalculationExecutionResponse>
      startCalculationExecution_Pre($grpc.ServiceCall $call,
          $async.Future<$0.StartCalculationExecutionRequest> $request) async {
    return startCalculationExecution($call, await $request);
  }

  $async.Future<$0.StartCalculationExecutionResponse> startCalculationExecution(
      $grpc.ServiceCall call, $0.StartCalculationExecutionRequest request);

  $async.Future<$0.StartQueryExecutionOutput> startQueryExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartQueryExecutionInput> $request) async {
    return startQueryExecution($call, await $request);
  }

  $async.Future<$0.StartQueryExecutionOutput> startQueryExecution(
      $grpc.ServiceCall call, $0.StartQueryExecutionInput request);

  $async.Future<$0.StartSessionResponse> startSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartSessionRequest> $request) async {
    return startSession($call, await $request);
  }

  $async.Future<$0.StartSessionResponse> startSession(
      $grpc.ServiceCall call, $0.StartSessionRequest request);

  $async.Future<$0.StopCalculationExecutionResponse>
      stopCalculationExecution_Pre($grpc.ServiceCall $call,
          $async.Future<$0.StopCalculationExecutionRequest> $request) async {
    return stopCalculationExecution($call, await $request);
  }

  $async.Future<$0.StopCalculationExecutionResponse> stopCalculationExecution(
      $grpc.ServiceCall call, $0.StopCalculationExecutionRequest request);

  $async.Future<$0.StopQueryExecutionOutput> stopQueryExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StopQueryExecutionInput> $request) async {
    return stopQueryExecution($call, await $request);
  }

  $async.Future<$0.StopQueryExecutionOutput> stopQueryExecution(
      $grpc.ServiceCall call, $0.StopQueryExecutionInput request);

  $async.Future<$0.TagResourceOutput> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceInput> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceOutput> tagResource(
      $grpc.ServiceCall call, $0.TagResourceInput request);

  $async.Future<$0.TerminateSessionResponse> terminateSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TerminateSessionRequest> $request) async {
    return terminateSession($call, await $request);
  }

  $async.Future<$0.TerminateSessionResponse> terminateSession(
      $grpc.ServiceCall call, $0.TerminateSessionRequest request);

  $async.Future<$0.UntagResourceOutput> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceInput> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceOutput> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceInput request);

  $async.Future<$0.UpdateCapacityReservationOutput>
      updateCapacityReservation_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateCapacityReservationInput> $request) async {
    return updateCapacityReservation($call, await $request);
  }

  $async.Future<$0.UpdateCapacityReservationOutput> updateCapacityReservation(
      $grpc.ServiceCall call, $0.UpdateCapacityReservationInput request);

  $async.Future<$0.UpdateDataCatalogOutput> updateDataCatalog_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDataCatalogInput> $request) async {
    return updateDataCatalog($call, await $request);
  }

  $async.Future<$0.UpdateDataCatalogOutput> updateDataCatalog(
      $grpc.ServiceCall call, $0.UpdateDataCatalogInput request);

  $async.Future<$0.UpdateNamedQueryOutput> updateNamedQuery_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateNamedQueryInput> $request) async {
    return updateNamedQuery($call, await $request);
  }

  $async.Future<$0.UpdateNamedQueryOutput> updateNamedQuery(
      $grpc.ServiceCall call, $0.UpdateNamedQueryInput request);

  $async.Future<$0.UpdateNotebookOutput> updateNotebook_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateNotebookInput> $request) async {
    return updateNotebook($call, await $request);
  }

  $async.Future<$0.UpdateNotebookOutput> updateNotebook(
      $grpc.ServiceCall call, $0.UpdateNotebookInput request);

  $async.Future<$0.UpdateNotebookMetadataOutput> updateNotebookMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateNotebookMetadataInput> $request) async {
    return updateNotebookMetadata($call, await $request);
  }

  $async.Future<$0.UpdateNotebookMetadataOutput> updateNotebookMetadata(
      $grpc.ServiceCall call, $0.UpdateNotebookMetadataInput request);

  $async.Future<$0.UpdatePreparedStatementOutput> updatePreparedStatement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdatePreparedStatementInput> $request) async {
    return updatePreparedStatement($call, await $request);
  }

  $async.Future<$0.UpdatePreparedStatementOutput> updatePreparedStatement(
      $grpc.ServiceCall call, $0.UpdatePreparedStatementInput request);

  $async.Future<$0.UpdateWorkGroupOutput> updateWorkGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateWorkGroupInput> $request) async {
    return updateWorkGroup($call, await $request);
  }

  $async.Future<$0.UpdateWorkGroupOutput> updateWorkGroup(
      $grpc.ServiceCall call, $0.UpdateWorkGroupInput request);
}
