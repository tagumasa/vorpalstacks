// This is a generated file - do not edit.
//
// Generated from dynamodb.proto.

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

import 'common.pb.dart' as $1;
import 'dynamodb.pb.dart' as $0;

export 'dynamodb.pb.dart';

/// DynamoDBService provides dynamodb API operations.
@$pb.GrpcServiceName('dynamodb.DynamoDBService')
class DynamoDBServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  DynamoDBServiceClient(super.channel, {super.options, super.interceptors});

  /// This operation allows you to perform batch reads or writes on data stored in DynamoDB, using PartiQL. Each read statement in a BatchExecuteStatement must specify an equality condition on all key at...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.BatchExecuteStatementOutput> batchExecuteStatement(
    $0.BatchExecuteStatementInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchExecuteStatement, request, options: options);
  }

  /// The BatchGetItem operation returns the attributes of one or more items from one or more tables. You identify requested items by primary key. A single operation can retrieve up to 16 MB of data, whi...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.BatchGetItemOutput> batchGetItem(
    $0.BatchGetItemInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchGetItem, request, options: options);
  }

  /// The BatchWriteItem operation puts or deletes multiple items in one or more tables. A single call to BatchWriteItem can transmit up to 16MB of data over the network, consisting of up to 25 item put ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.BatchWriteItemOutput> batchWriteItem(
    $0.BatchWriteItemInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$batchWriteItem, request, options: options);
  }

  /// Creates a backup for an existing table. Each time you create an on-demand backup, the entire table data is backed up. There is no limit to the number of on-demand backups that can be taken. When yo...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateBackupOutput> createBackup(
    $0.CreateBackupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createBackup, request, options: options);
  }

  /// Creates a global table from an existing table. A global table creates a replication relationship between two or more DynamoDB tables with the same table name in the provided Regions. This documenta...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateGlobalTableOutput> createGlobalTable(
    $0.CreateGlobalTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createGlobalTable, request, options: options);
  }

  /// The CreateTable operation adds a new table to your account. In an Amazon Web Services account, table names must be unique within each Region. That is, you can have two tables with same name if you ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.CreateTableOutput> createTable(
    $0.CreateTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTable, request, options: options);
  }

  /// Deletes an existing backup of a table. You can call DeleteBackup at a maximum rate of 10 times per second.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteBackupOutput> deleteBackup(
    $0.DeleteBackupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteBackup, request, options: options);
  }

  /// Deletes a single item in a table by primary key. You can perform a conditional delete operation that deletes the item if it exists, or if it has an expected attribute value. In addition to deleting...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteItemOutput> deleteItem(
    $0.DeleteItemInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteItem, request, options: options);
  }

  /// Deletes the resource-based policy attached to the resource, which can be a table or stream. DeleteResourcePolicy is an idempotent operation; running it multiple times on the same resource doesn't r...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteResourcePolicyOutput> deleteResourcePolicy(
    $0.DeleteResourcePolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourcePolicy, request, options: options);
  }

  /// The DeleteTable operation deletes a table and all of its items. After a DeleteTable request, the specified table is in the DELETING state until DynamoDB completes the deletion. If the table is in t...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DeleteTableOutput> deleteTable(
    $0.DeleteTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTable, request, options: options);
  }

  /// Describes an existing backup of a table. You can call DescribeBackup at a maximum rate of 10 times per second.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeBackupOutput> describeBackup(
    $0.DescribeBackupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeBackup, request, options: options);
  }

  /// Checks the status of continuous backups and point in time recovery on the specified table. Continuous backups are ENABLED on all tables at table creation. If point in time recovery is enabled, Poin...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeContinuousBackupsOutput>
      describeContinuousBackups(
    $0.DescribeContinuousBackupsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeContinuousBackups, request,
        options: options);
  }

  /// Returns information about contributor insights for a given table or global secondary index.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeContributorInsightsOutput>
      describeContributorInsights(
    $0.DescribeContributorInsightsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeContributorInsights, request,
        options: options);
  }

  /// Returns the regional endpoint information. For more information on policy permissions, please see Internetwork traffic privacy.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeEndpointsResponse> describeEndpoints(
    $0.DescribeEndpointsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeEndpoints, request, options: options);
  }

  /// Describes an existing table export.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeExportOutput> describeExport(
    $0.DescribeExportInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeExport, request, options: options);
  }

  /// Returns information about the specified global table. This documentation is for version 2017.11.29 (Legacy) of global tables, which should be avoided for new global tables. Customers should use Glo...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeGlobalTableOutput> describeGlobalTable(
    $0.DescribeGlobalTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeGlobalTable, request, options: options);
  }

  /// Describes Region-specific settings for a global table. This documentation is for version 2017.11.29 (Legacy) of global tables, which should be avoided for new global tables. Customers should use Gl...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeGlobalTableSettingsOutput>
      describeGlobalTableSettings(
    $0.DescribeGlobalTableSettingsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeGlobalTableSettings, request,
        options: options);
  }

  /// Represents the properties of the import.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeImportOutput> describeImport(
    $0.DescribeImportInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeImport, request, options: options);
  }

  /// Returns information about the status of Kinesis streaming.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeKinesisStreamingDestinationOutput>
      describeKinesisStreamingDestination(
    $0.DescribeKinesisStreamingDestinationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeKinesisStreamingDestination, request,
        options: options);
  }

  /// Returns the current provisioned-capacity quotas for your Amazon Web Services account in a Region, both for the Region as a whole and for any one DynamoDB table that you create there. When you estab...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeLimitsOutput> describeLimits(
    $0.DescribeLimitsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeLimits, request, options: options);
  }

  /// Returns information about the table, including the current status of the table, when it was created, the primary key schema, and any indexes on the table. If you issue a DescribeTable request immed...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeTableOutput> describeTable(
    $0.DescribeTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeTable, request, options: options);
  }

  /// Describes auto scaling settings across replicas of the global table at once.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeTableReplicaAutoScalingOutput>
      describeTableReplicaAutoScaling(
    $0.DescribeTableReplicaAutoScalingInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeTableReplicaAutoScaling, request,
        options: options);
  }

  /// Gives a description of the Time to Live (TTL) status on the specified table.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.DescribeTimeToLiveOutput> describeTimeToLive(
    $0.DescribeTimeToLiveInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeTimeToLive, request, options: options);
  }

  /// Stops replication from the DynamoDB table to the Kinesis data stream. This is done without deleting either of the resources.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.KinesisStreamingDestinationOutput>
      disableKinesisStreamingDestination(
    $0.KinesisStreamingDestinationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableKinesisStreamingDestination, request,
        options: options);
  }

  /// Starts table data replication to the specified Kinesis data stream at a timestamp chosen during the enable workflow. If this operation doesn't return results immediately, use DescribeKinesisStreami...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.KinesisStreamingDestinationOutput>
      enableKinesisStreamingDestination(
    $0.KinesisStreamingDestinationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableKinesisStreamingDestination, request,
        options: options);
  }

  /// This operation allows you to perform reads and singleton writes on data stored in DynamoDB, using PartiQL. For PartiQL reads (SELECT statement), if the total number of processed items exceeds the m...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ExecuteStatementOutput> executeStatement(
    $0.ExecuteStatementInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$executeStatement, request, options: options);
  }

  /// This operation allows you to perform transactional reads or writes on data stored in DynamoDB, using PartiQL. The entire transaction must consist of either read statements or write statements, you ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ExecuteTransactionOutput> executeTransaction(
    $0.ExecuteTransactionInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$executeTransaction, request, options: options);
  }

  /// Exports table data to an S3 bucket. The table must have point in time recovery enabled, and you can export data from any time within the point in time recovery window.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ExportTableToPointInTimeOutput>
      exportTableToPointInTime(
    $0.ExportTableToPointInTimeInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$exportTableToPointInTime, request,
        options: options);
  }

  /// The GetItem operation returns a set of attributes for the item with the given primary key. If there is no matching item, GetItem does not return any data and there will be no Item element in the re...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetItemOutput> getItem(
    $0.GetItemInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getItem, request, options: options);
  }

  /// Returns the resource-based policy document attached to the resource, which can be a table or stream, in JSON format. GetResourcePolicy follows an eventually consistent model. The following list des...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.GetResourcePolicyOutput> getResourcePolicy(
    $0.GetResourcePolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourcePolicy, request, options: options);
  }

  /// Imports table data from an S3 bucket.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ImportTableOutput> importTable(
    $0.ImportTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importTable, request, options: options);
  }

  /// List DynamoDB backups that are associated with an Amazon Web Services account and weren't made with Amazon Web Services Backup. To list these backups for a given table, specify TableName. ListBacku...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListBackupsOutput> listBackups(
    $0.ListBackupsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listBackups, request, options: options);
  }

  /// Returns a list of ContributorInsightsSummary for a table and all its global secondary indexes.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListContributorInsightsOutput>
      listContributorInsights(
    $0.ListContributorInsightsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listContributorInsights, request,
        options: options);
  }

  /// Lists completed exports within the past 90 days, in reverse alphanumeric order of ExportArn.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListExportsOutput> listExports(
    $0.ListExportsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listExports, request, options: options);
  }

  /// Lists all global tables that have a replica in the specified Region. This documentation is for version 2017.11.29 (Legacy) of global tables, which should be avoided for new global tables. Customers...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListGlobalTablesOutput> listGlobalTables(
    $0.ListGlobalTablesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGlobalTables, request, options: options);
  }

  /// Lists completed imports within the past 90 days.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListImportsOutput> listImports(
    $0.ListImportsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listImports, request, options: options);
  }

  /// Returns an array of table names associated with the current account and endpoint. The output from ListTables is paginated, with each page returning a maximum of 100 table names.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListTablesOutput> listTables(
    $0.ListTablesInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTables, request, options: options);
  }

  /// List all tags on an Amazon DynamoDB resource. You can call ListTagsOfResource up to 10 times per second, per account. For an overview on tagging DynamoDB resources, see Tagging for DynamoDB in the ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ListTagsOfResourceOutput> listTagsOfResource(
    $0.ListTagsOfResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsOfResource, request, options: options);
  }

  /// Creates a new item, or replaces an old item with a new item. If an item that has the same primary key as the new item already exists in the specified table, the new item completely replaces the exi...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PutItemOutput> putItem(
    $0.PutItemInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putItem, request, options: options);
  }

  /// Attaches a resource-based policy document to the resource, which can be a table or stream. When you attach a resource-based policy using this API, the policy application is eventually consistent . ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.PutResourcePolicyOutput> putResourcePolicy(
    $0.PutResourcePolicyInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putResourcePolicy, request, options: options);
  }

  /// You must provide the name of the partition key attribute and a single value for that attribute. Query returns all items with that partition key value. Optionally, you can provide a sort key attribu...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.QueryOutput> query(
    $0.QueryInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$query, request, options: options);
  }

  /// Creates a new table from an existing backup. Any number of users can execute up to 50 concurrent restores (any type of restore) in a given account. You can call RestoreTableFromBackup at a maximum ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.RestoreTableFromBackupOutput> restoreTableFromBackup(
    $0.RestoreTableFromBackupInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$restoreTableFromBackup, request,
        options: options);
  }

  /// Restores the specified table to the specified point in time within EarliestRestorableDateTime and LatestRestorableDateTime. You can restore your table to any point in time in the last 35 days. You ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.RestoreTableToPointInTimeOutput>
      restoreTableToPointInTime(
    $0.RestoreTableToPointInTimeInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$restoreTableToPointInTime, request,
        options: options);
  }

  /// The Scan operation returns one or more items and item attributes by accessing every item in a table or a secondary index. To have DynamoDB return fewer items, you can provide a FilterExpression ope...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.ScanOutput> scan(
    $0.ScanInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$scan, request, options: options);
  }

  /// Associate a set of tags with an Amazon DynamoDB resource. You can then activate these user-defined tags so that they appear on the Billing and Cost Management console for cost allocation tracking. ...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// TransactGetItems is a synchronous operation that atomically retrieves multiple items from one or more tables (but not from indexes) in a single account and Region. A TransactGetItems call can conta...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.TransactGetItemsOutput> transactGetItems(
    $0.TransactGetItemsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$transactGetItems, request, options: options);
  }

  /// TransactWriteItems is a synchronous write operation that groups up to 100 action requests. These actions can target items in different tables, but not in different Amazon Web Services accounts or R...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.TransactWriteItemsOutput> transactWriteItems(
    $0.TransactWriteItemsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$transactWriteItems, request, options: options);
  }

  /// Removes the association of tags from an Amazon DynamoDB resource. You can call UntagResource up to five times per second, per account. UntagResource is an asynchronous operation. If you issue a Lis...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// UpdateContinuousBackups enables or disables point in time recovery for the specified table. A successful UpdateContinuousBackups call returns the current ContinuousBackupsDescription. Continuous ba...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateContinuousBackupsOutput>
      updateContinuousBackups(
    $0.UpdateContinuousBackupsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateContinuousBackups, request,
        options: options);
  }

  /// Updates the status for contributor insights for a specific table or index. CloudWatch Contributor Insights for DynamoDB graphs display the partition key and (if applicable) sort key of frequently a...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateContributorInsightsOutput>
      updateContributorInsights(
    $0.UpdateContributorInsightsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateContributorInsights, request,
        options: options);
  }

  /// Adds or removes replicas in the specified global table. The global table must already exist to be able to use this operation. Any replica to be added must be empty, have the same name as the global...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateGlobalTableOutput> updateGlobalTable(
    $0.UpdateGlobalTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateGlobalTable, request, options: options);
  }

  /// Updates settings for a global table. This documentation is for version 2017.11.29 (Legacy) of global tables, which should be avoided for new global tables. Customers should use Global Tables versio...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateGlobalTableSettingsOutput>
      updateGlobalTableSettings(
    $0.UpdateGlobalTableSettingsInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateGlobalTableSettings, request,
        options: options);
  }

  /// Edits an existing item's attributes, or adds a new item to the table if it does not already exist. You can put, delete, or add attribute values. You can also perform a conditional update on an exis...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateItemOutput> updateItem(
    $0.UpdateItemInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateItem, request, options: options);
  }

  /// The command to update the Kinesis stream destination.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateKinesisStreamingDestinationOutput>
      updateKinesisStreamingDestination(
    $0.UpdateKinesisStreamingDestinationInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateKinesisStreamingDestination, request,
        options: options);
  }

  /// Modifies the provisioned throughput settings, global secondary indexes, or DynamoDB Streams settings for a given table. You can only perform one of the following operations at once: Modify the prov...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateTableOutput> updateTable(
    $0.UpdateTableInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTable, request, options: options);
  }

  /// Updates auto scaling settings on your global tables at once.
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateTableReplicaAutoScalingOutput>
      updateTableReplicaAutoScaling(
    $0.UpdateTableReplicaAutoScalingInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTableReplicaAutoScaling, request,
        options: options);
  }

  /// The UpdateTimeToLive method enables or disables Time to Live (TTL) for the specified table. A successful UpdateTimeToLive call returns the current TimeToLiveSpecification. It can take up to one hou...
  /// HTTP:
  /// Protocol: awsJson1_0
  $grpc.ResponseFuture<$0.UpdateTimeToLiveOutput> updateTimeToLive(
    $0.UpdateTimeToLiveInput request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTimeToLive, request, options: options);
  }

  // method descriptors

  static final _$batchExecuteStatement = $grpc.ClientMethod<
          $0.BatchExecuteStatementInput, $0.BatchExecuteStatementOutput>(
      '/dynamodb.DynamoDBService/BatchExecuteStatement',
      ($0.BatchExecuteStatementInput value) => value.writeToBuffer(),
      $0.BatchExecuteStatementOutput.fromBuffer);
  static final _$batchGetItem =
      $grpc.ClientMethod<$0.BatchGetItemInput, $0.BatchGetItemOutput>(
          '/dynamodb.DynamoDBService/BatchGetItem',
          ($0.BatchGetItemInput value) => value.writeToBuffer(),
          $0.BatchGetItemOutput.fromBuffer);
  static final _$batchWriteItem =
      $grpc.ClientMethod<$0.BatchWriteItemInput, $0.BatchWriteItemOutput>(
          '/dynamodb.DynamoDBService/BatchWriteItem',
          ($0.BatchWriteItemInput value) => value.writeToBuffer(),
          $0.BatchWriteItemOutput.fromBuffer);
  static final _$createBackup =
      $grpc.ClientMethod<$0.CreateBackupInput, $0.CreateBackupOutput>(
          '/dynamodb.DynamoDBService/CreateBackup',
          ($0.CreateBackupInput value) => value.writeToBuffer(),
          $0.CreateBackupOutput.fromBuffer);
  static final _$createGlobalTable =
      $grpc.ClientMethod<$0.CreateGlobalTableInput, $0.CreateGlobalTableOutput>(
          '/dynamodb.DynamoDBService/CreateGlobalTable',
          ($0.CreateGlobalTableInput value) => value.writeToBuffer(),
          $0.CreateGlobalTableOutput.fromBuffer);
  static final _$createTable =
      $grpc.ClientMethod<$0.CreateTableInput, $0.CreateTableOutput>(
          '/dynamodb.DynamoDBService/CreateTable',
          ($0.CreateTableInput value) => value.writeToBuffer(),
          $0.CreateTableOutput.fromBuffer);
  static final _$deleteBackup =
      $grpc.ClientMethod<$0.DeleteBackupInput, $0.DeleteBackupOutput>(
          '/dynamodb.DynamoDBService/DeleteBackup',
          ($0.DeleteBackupInput value) => value.writeToBuffer(),
          $0.DeleteBackupOutput.fromBuffer);
  static final _$deleteItem =
      $grpc.ClientMethod<$0.DeleteItemInput, $0.DeleteItemOutput>(
          '/dynamodb.DynamoDBService/DeleteItem',
          ($0.DeleteItemInput value) => value.writeToBuffer(),
          $0.DeleteItemOutput.fromBuffer);
  static final _$deleteResourcePolicy = $grpc.ClientMethod<
          $0.DeleteResourcePolicyInput, $0.DeleteResourcePolicyOutput>(
      '/dynamodb.DynamoDBService/DeleteResourcePolicy',
      ($0.DeleteResourcePolicyInput value) => value.writeToBuffer(),
      $0.DeleteResourcePolicyOutput.fromBuffer);
  static final _$deleteTable =
      $grpc.ClientMethod<$0.DeleteTableInput, $0.DeleteTableOutput>(
          '/dynamodb.DynamoDBService/DeleteTable',
          ($0.DeleteTableInput value) => value.writeToBuffer(),
          $0.DeleteTableOutput.fromBuffer);
  static final _$describeBackup =
      $grpc.ClientMethod<$0.DescribeBackupInput, $0.DescribeBackupOutput>(
          '/dynamodb.DynamoDBService/DescribeBackup',
          ($0.DescribeBackupInput value) => value.writeToBuffer(),
          $0.DescribeBackupOutput.fromBuffer);
  static final _$describeContinuousBackups = $grpc.ClientMethod<
          $0.DescribeContinuousBackupsInput,
          $0.DescribeContinuousBackupsOutput>(
      '/dynamodb.DynamoDBService/DescribeContinuousBackups',
      ($0.DescribeContinuousBackupsInput value) => value.writeToBuffer(),
      $0.DescribeContinuousBackupsOutput.fromBuffer);
  static final _$describeContributorInsights = $grpc.ClientMethod<
          $0.DescribeContributorInsightsInput,
          $0.DescribeContributorInsightsOutput>(
      '/dynamodb.DynamoDBService/DescribeContributorInsights',
      ($0.DescribeContributorInsightsInput value) => value.writeToBuffer(),
      $0.DescribeContributorInsightsOutput.fromBuffer);
  static final _$describeEndpoints = $grpc.ClientMethod<
          $0.DescribeEndpointsRequest, $0.DescribeEndpointsResponse>(
      '/dynamodb.DynamoDBService/DescribeEndpoints',
      ($0.DescribeEndpointsRequest value) => value.writeToBuffer(),
      $0.DescribeEndpointsResponse.fromBuffer);
  static final _$describeExport =
      $grpc.ClientMethod<$0.DescribeExportInput, $0.DescribeExportOutput>(
          '/dynamodb.DynamoDBService/DescribeExport',
          ($0.DescribeExportInput value) => value.writeToBuffer(),
          $0.DescribeExportOutput.fromBuffer);
  static final _$describeGlobalTable = $grpc.ClientMethod<
          $0.DescribeGlobalTableInput, $0.DescribeGlobalTableOutput>(
      '/dynamodb.DynamoDBService/DescribeGlobalTable',
      ($0.DescribeGlobalTableInput value) => value.writeToBuffer(),
      $0.DescribeGlobalTableOutput.fromBuffer);
  static final _$describeGlobalTableSettings = $grpc.ClientMethod<
          $0.DescribeGlobalTableSettingsInput,
          $0.DescribeGlobalTableSettingsOutput>(
      '/dynamodb.DynamoDBService/DescribeGlobalTableSettings',
      ($0.DescribeGlobalTableSettingsInput value) => value.writeToBuffer(),
      $0.DescribeGlobalTableSettingsOutput.fromBuffer);
  static final _$describeImport =
      $grpc.ClientMethod<$0.DescribeImportInput, $0.DescribeImportOutput>(
          '/dynamodb.DynamoDBService/DescribeImport',
          ($0.DescribeImportInput value) => value.writeToBuffer(),
          $0.DescribeImportOutput.fromBuffer);
  static final _$describeKinesisStreamingDestination = $grpc.ClientMethod<
          $0.DescribeKinesisStreamingDestinationInput,
          $0.DescribeKinesisStreamingDestinationOutput>(
      '/dynamodb.DynamoDBService/DescribeKinesisStreamingDestination',
      ($0.DescribeKinesisStreamingDestinationInput value) =>
          value.writeToBuffer(),
      $0.DescribeKinesisStreamingDestinationOutput.fromBuffer);
  static final _$describeLimits =
      $grpc.ClientMethod<$0.DescribeLimitsInput, $0.DescribeLimitsOutput>(
          '/dynamodb.DynamoDBService/DescribeLimits',
          ($0.DescribeLimitsInput value) => value.writeToBuffer(),
          $0.DescribeLimitsOutput.fromBuffer);
  static final _$describeTable =
      $grpc.ClientMethod<$0.DescribeTableInput, $0.DescribeTableOutput>(
          '/dynamodb.DynamoDBService/DescribeTable',
          ($0.DescribeTableInput value) => value.writeToBuffer(),
          $0.DescribeTableOutput.fromBuffer);
  static final _$describeTableReplicaAutoScaling = $grpc.ClientMethod<
          $0.DescribeTableReplicaAutoScalingInput,
          $0.DescribeTableReplicaAutoScalingOutput>(
      '/dynamodb.DynamoDBService/DescribeTableReplicaAutoScaling',
      ($0.DescribeTableReplicaAutoScalingInput value) => value.writeToBuffer(),
      $0.DescribeTableReplicaAutoScalingOutput.fromBuffer);
  static final _$describeTimeToLive = $grpc.ClientMethod<
          $0.DescribeTimeToLiveInput, $0.DescribeTimeToLiveOutput>(
      '/dynamodb.DynamoDBService/DescribeTimeToLive',
      ($0.DescribeTimeToLiveInput value) => value.writeToBuffer(),
      $0.DescribeTimeToLiveOutput.fromBuffer);
  static final _$disableKinesisStreamingDestination = $grpc.ClientMethod<
          $0.KinesisStreamingDestinationInput,
          $0.KinesisStreamingDestinationOutput>(
      '/dynamodb.DynamoDBService/DisableKinesisStreamingDestination',
      ($0.KinesisStreamingDestinationInput value) => value.writeToBuffer(),
      $0.KinesisStreamingDestinationOutput.fromBuffer);
  static final _$enableKinesisStreamingDestination = $grpc.ClientMethod<
          $0.KinesisStreamingDestinationInput,
          $0.KinesisStreamingDestinationOutput>(
      '/dynamodb.DynamoDBService/EnableKinesisStreamingDestination',
      ($0.KinesisStreamingDestinationInput value) => value.writeToBuffer(),
      $0.KinesisStreamingDestinationOutput.fromBuffer);
  static final _$executeStatement =
      $grpc.ClientMethod<$0.ExecuteStatementInput, $0.ExecuteStatementOutput>(
          '/dynamodb.DynamoDBService/ExecuteStatement',
          ($0.ExecuteStatementInput value) => value.writeToBuffer(),
          $0.ExecuteStatementOutput.fromBuffer);
  static final _$executeTransaction = $grpc.ClientMethod<
          $0.ExecuteTransactionInput, $0.ExecuteTransactionOutput>(
      '/dynamodb.DynamoDBService/ExecuteTransaction',
      ($0.ExecuteTransactionInput value) => value.writeToBuffer(),
      $0.ExecuteTransactionOutput.fromBuffer);
  static final _$exportTableToPointInTime = $grpc.ClientMethod<
          $0.ExportTableToPointInTimeInput, $0.ExportTableToPointInTimeOutput>(
      '/dynamodb.DynamoDBService/ExportTableToPointInTime',
      ($0.ExportTableToPointInTimeInput value) => value.writeToBuffer(),
      $0.ExportTableToPointInTimeOutput.fromBuffer);
  static final _$getItem =
      $grpc.ClientMethod<$0.GetItemInput, $0.GetItemOutput>(
          '/dynamodb.DynamoDBService/GetItem',
          ($0.GetItemInput value) => value.writeToBuffer(),
          $0.GetItemOutput.fromBuffer);
  static final _$getResourcePolicy =
      $grpc.ClientMethod<$0.GetResourcePolicyInput, $0.GetResourcePolicyOutput>(
          '/dynamodb.DynamoDBService/GetResourcePolicy',
          ($0.GetResourcePolicyInput value) => value.writeToBuffer(),
          $0.GetResourcePolicyOutput.fromBuffer);
  static final _$importTable =
      $grpc.ClientMethod<$0.ImportTableInput, $0.ImportTableOutput>(
          '/dynamodb.DynamoDBService/ImportTable',
          ($0.ImportTableInput value) => value.writeToBuffer(),
          $0.ImportTableOutput.fromBuffer);
  static final _$listBackups =
      $grpc.ClientMethod<$0.ListBackupsInput, $0.ListBackupsOutput>(
          '/dynamodb.DynamoDBService/ListBackups',
          ($0.ListBackupsInput value) => value.writeToBuffer(),
          $0.ListBackupsOutput.fromBuffer);
  static final _$listContributorInsights = $grpc.ClientMethod<
          $0.ListContributorInsightsInput, $0.ListContributorInsightsOutput>(
      '/dynamodb.DynamoDBService/ListContributorInsights',
      ($0.ListContributorInsightsInput value) => value.writeToBuffer(),
      $0.ListContributorInsightsOutput.fromBuffer);
  static final _$listExports =
      $grpc.ClientMethod<$0.ListExportsInput, $0.ListExportsOutput>(
          '/dynamodb.DynamoDBService/ListExports',
          ($0.ListExportsInput value) => value.writeToBuffer(),
          $0.ListExportsOutput.fromBuffer);
  static final _$listGlobalTables =
      $grpc.ClientMethod<$0.ListGlobalTablesInput, $0.ListGlobalTablesOutput>(
          '/dynamodb.DynamoDBService/ListGlobalTables',
          ($0.ListGlobalTablesInput value) => value.writeToBuffer(),
          $0.ListGlobalTablesOutput.fromBuffer);
  static final _$listImports =
      $grpc.ClientMethod<$0.ListImportsInput, $0.ListImportsOutput>(
          '/dynamodb.DynamoDBService/ListImports',
          ($0.ListImportsInput value) => value.writeToBuffer(),
          $0.ListImportsOutput.fromBuffer);
  static final _$listTables =
      $grpc.ClientMethod<$0.ListTablesInput, $0.ListTablesOutput>(
          '/dynamodb.DynamoDBService/ListTables',
          ($0.ListTablesInput value) => value.writeToBuffer(),
          $0.ListTablesOutput.fromBuffer);
  static final _$listTagsOfResource = $grpc.ClientMethod<
          $0.ListTagsOfResourceInput, $0.ListTagsOfResourceOutput>(
      '/dynamodb.DynamoDBService/ListTagsOfResource',
      ($0.ListTagsOfResourceInput value) => value.writeToBuffer(),
      $0.ListTagsOfResourceOutput.fromBuffer);
  static final _$putItem =
      $grpc.ClientMethod<$0.PutItemInput, $0.PutItemOutput>(
          '/dynamodb.DynamoDBService/PutItem',
          ($0.PutItemInput value) => value.writeToBuffer(),
          $0.PutItemOutput.fromBuffer);
  static final _$putResourcePolicy =
      $grpc.ClientMethod<$0.PutResourcePolicyInput, $0.PutResourcePolicyOutput>(
          '/dynamodb.DynamoDBService/PutResourcePolicy',
          ($0.PutResourcePolicyInput value) => value.writeToBuffer(),
          $0.PutResourcePolicyOutput.fromBuffer);
  static final _$query = $grpc.ClientMethod<$0.QueryInput, $0.QueryOutput>(
      '/dynamodb.DynamoDBService/Query',
      ($0.QueryInput value) => value.writeToBuffer(),
      $0.QueryOutput.fromBuffer);
  static final _$restoreTableFromBackup = $grpc.ClientMethod<
          $0.RestoreTableFromBackupInput, $0.RestoreTableFromBackupOutput>(
      '/dynamodb.DynamoDBService/RestoreTableFromBackup',
      ($0.RestoreTableFromBackupInput value) => value.writeToBuffer(),
      $0.RestoreTableFromBackupOutput.fromBuffer);
  static final _$restoreTableToPointInTime = $grpc.ClientMethod<
          $0.RestoreTableToPointInTimeInput,
          $0.RestoreTableToPointInTimeOutput>(
      '/dynamodb.DynamoDBService/RestoreTableToPointInTime',
      ($0.RestoreTableToPointInTimeInput value) => value.writeToBuffer(),
      $0.RestoreTableToPointInTimeOutput.fromBuffer);
  static final _$scan = $grpc.ClientMethod<$0.ScanInput, $0.ScanOutput>(
      '/dynamodb.DynamoDBService/Scan',
      ($0.ScanInput value) => value.writeToBuffer(),
      $0.ScanOutput.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceInput, $1.Empty>(
          '/dynamodb.DynamoDBService/TagResource',
          ($0.TagResourceInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$transactGetItems =
      $grpc.ClientMethod<$0.TransactGetItemsInput, $0.TransactGetItemsOutput>(
          '/dynamodb.DynamoDBService/TransactGetItems',
          ($0.TransactGetItemsInput value) => value.writeToBuffer(),
          $0.TransactGetItemsOutput.fromBuffer);
  static final _$transactWriteItems = $grpc.ClientMethod<
          $0.TransactWriteItemsInput, $0.TransactWriteItemsOutput>(
      '/dynamodb.DynamoDBService/TransactWriteItems',
      ($0.TransactWriteItemsInput value) => value.writeToBuffer(),
      $0.TransactWriteItemsOutput.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceInput, $1.Empty>(
          '/dynamodb.DynamoDBService/UntagResource',
          ($0.UntagResourceInput value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateContinuousBackups = $grpc.ClientMethod<
          $0.UpdateContinuousBackupsInput, $0.UpdateContinuousBackupsOutput>(
      '/dynamodb.DynamoDBService/UpdateContinuousBackups',
      ($0.UpdateContinuousBackupsInput value) => value.writeToBuffer(),
      $0.UpdateContinuousBackupsOutput.fromBuffer);
  static final _$updateContributorInsights = $grpc.ClientMethod<
          $0.UpdateContributorInsightsInput,
          $0.UpdateContributorInsightsOutput>(
      '/dynamodb.DynamoDBService/UpdateContributorInsights',
      ($0.UpdateContributorInsightsInput value) => value.writeToBuffer(),
      $0.UpdateContributorInsightsOutput.fromBuffer);
  static final _$updateGlobalTable =
      $grpc.ClientMethod<$0.UpdateGlobalTableInput, $0.UpdateGlobalTableOutput>(
          '/dynamodb.DynamoDBService/UpdateGlobalTable',
          ($0.UpdateGlobalTableInput value) => value.writeToBuffer(),
          $0.UpdateGlobalTableOutput.fromBuffer);
  static final _$updateGlobalTableSettings = $grpc.ClientMethod<
          $0.UpdateGlobalTableSettingsInput,
          $0.UpdateGlobalTableSettingsOutput>(
      '/dynamodb.DynamoDBService/UpdateGlobalTableSettings',
      ($0.UpdateGlobalTableSettingsInput value) => value.writeToBuffer(),
      $0.UpdateGlobalTableSettingsOutput.fromBuffer);
  static final _$updateItem =
      $grpc.ClientMethod<$0.UpdateItemInput, $0.UpdateItemOutput>(
          '/dynamodb.DynamoDBService/UpdateItem',
          ($0.UpdateItemInput value) => value.writeToBuffer(),
          $0.UpdateItemOutput.fromBuffer);
  static final _$updateKinesisStreamingDestination = $grpc.ClientMethod<
          $0.UpdateKinesisStreamingDestinationInput,
          $0.UpdateKinesisStreamingDestinationOutput>(
      '/dynamodb.DynamoDBService/UpdateKinesisStreamingDestination',
      ($0.UpdateKinesisStreamingDestinationInput value) =>
          value.writeToBuffer(),
      $0.UpdateKinesisStreamingDestinationOutput.fromBuffer);
  static final _$updateTable =
      $grpc.ClientMethod<$0.UpdateTableInput, $0.UpdateTableOutput>(
          '/dynamodb.DynamoDBService/UpdateTable',
          ($0.UpdateTableInput value) => value.writeToBuffer(),
          $0.UpdateTableOutput.fromBuffer);
  static final _$updateTableReplicaAutoScaling = $grpc.ClientMethod<
          $0.UpdateTableReplicaAutoScalingInput,
          $0.UpdateTableReplicaAutoScalingOutput>(
      '/dynamodb.DynamoDBService/UpdateTableReplicaAutoScaling',
      ($0.UpdateTableReplicaAutoScalingInput value) => value.writeToBuffer(),
      $0.UpdateTableReplicaAutoScalingOutput.fromBuffer);
  static final _$updateTimeToLive =
      $grpc.ClientMethod<$0.UpdateTimeToLiveInput, $0.UpdateTimeToLiveOutput>(
          '/dynamodb.DynamoDBService/UpdateTimeToLive',
          ($0.UpdateTimeToLiveInput value) => value.writeToBuffer(),
          $0.UpdateTimeToLiveOutput.fromBuffer);
}

@$pb.GrpcServiceName('dynamodb.DynamoDBService')
abstract class DynamoDBServiceBase extends $grpc.Service {
  $core.String get $name => 'dynamodb.DynamoDBService';

  DynamoDBServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.BatchExecuteStatementInput,
            $0.BatchExecuteStatementOutput>(
        'BatchExecuteStatement',
        batchExecuteStatement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.BatchExecuteStatementInput.fromBuffer(value),
        ($0.BatchExecuteStatementOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.BatchGetItemInput, $0.BatchGetItemOutput>(
        'BatchGetItem',
        batchGetItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.BatchGetItemInput.fromBuffer(value),
        ($0.BatchGetItemOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.BatchWriteItemInput, $0.BatchWriteItemOutput>(
            'BatchWriteItem',
            batchWriteItem_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.BatchWriteItemInput.fromBuffer(value),
            ($0.BatchWriteItemOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateBackupInput, $0.CreateBackupOutput>(
        'CreateBackup',
        createBackup_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateBackupInput.fromBuffer(value),
        ($0.CreateBackupOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateGlobalTableInput,
            $0.CreateGlobalTableOutput>(
        'CreateGlobalTable',
        createGlobalTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateGlobalTableInput.fromBuffer(value),
        ($0.CreateGlobalTableOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateTableInput, $0.CreateTableOutput>(
        'CreateTable',
        createTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateTableInput.fromBuffer(value),
        ($0.CreateTableOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteBackupInput, $0.DeleteBackupOutput>(
        'DeleteBackup',
        deleteBackup_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteBackupInput.fromBuffer(value),
        ($0.DeleteBackupOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteItemInput, $0.DeleteItemOutput>(
        'DeleteItem',
        deleteItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteItemInput.fromBuffer(value),
        ($0.DeleteItemOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourcePolicyInput,
            $0.DeleteResourcePolicyOutput>(
        'DeleteResourcePolicy',
        deleteResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourcePolicyInput.fromBuffer(value),
        ($0.DeleteResourcePolicyOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTableInput, $0.DeleteTableOutput>(
        'DeleteTable',
        deleteTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteTableInput.fromBuffer(value),
        ($0.DeleteTableOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeBackupInput, $0.DescribeBackupOutput>(
            'DescribeBackup',
            describeBackup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeBackupInput.fromBuffer(value),
            ($0.DescribeBackupOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeContinuousBackupsInput,
            $0.DescribeContinuousBackupsOutput>(
        'DescribeContinuousBackups',
        describeContinuousBackups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeContinuousBackupsInput.fromBuffer(value),
        ($0.DescribeContinuousBackupsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeContributorInsightsInput,
            $0.DescribeContributorInsightsOutput>(
        'DescribeContributorInsights',
        describeContributorInsights_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeContributorInsightsInput.fromBuffer(value),
        ($0.DescribeContributorInsightsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeEndpointsRequest,
            $0.DescribeEndpointsResponse>(
        'DescribeEndpoints',
        describeEndpoints_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeEndpointsRequest.fromBuffer(value),
        ($0.DescribeEndpointsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeExportInput, $0.DescribeExportOutput>(
            'DescribeExport',
            describeExport_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeExportInput.fromBuffer(value),
            ($0.DescribeExportOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeGlobalTableInput,
            $0.DescribeGlobalTableOutput>(
        'DescribeGlobalTable',
        describeGlobalTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeGlobalTableInput.fromBuffer(value),
        ($0.DescribeGlobalTableOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeGlobalTableSettingsInput,
            $0.DescribeGlobalTableSettingsOutput>(
        'DescribeGlobalTableSettings',
        describeGlobalTableSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeGlobalTableSettingsInput.fromBuffer(value),
        ($0.DescribeGlobalTableSettingsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeImportInput, $0.DescribeImportOutput>(
            'DescribeImport',
            describeImport_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeImportInput.fromBuffer(value),
            ($0.DescribeImportOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeKinesisStreamingDestinationInput,
            $0.DescribeKinesisStreamingDestinationOutput>(
        'DescribeKinesisStreamingDestination',
        describeKinesisStreamingDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeKinesisStreamingDestinationInput.fromBuffer(value),
        ($0.DescribeKinesisStreamingDestinationOutput value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeLimitsInput, $0.DescribeLimitsOutput>(
            'DescribeLimits',
            describeLimits_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeLimitsInput.fromBuffer(value),
            ($0.DescribeLimitsOutput value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeTableInput, $0.DescribeTableOutput>(
            'DescribeTable',
            describeTable_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeTableInput.fromBuffer(value),
            ($0.DescribeTableOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeTableReplicaAutoScalingInput,
            $0.DescribeTableReplicaAutoScalingOutput>(
        'DescribeTableReplicaAutoScaling',
        describeTableReplicaAutoScaling_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeTableReplicaAutoScalingInput.fromBuffer(value),
        ($0.DescribeTableReplicaAutoScalingOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeTimeToLiveInput,
            $0.DescribeTimeToLiveOutput>(
        'DescribeTimeToLive',
        describeTimeToLive_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeTimeToLiveInput.fromBuffer(value),
        ($0.DescribeTimeToLiveOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.KinesisStreamingDestinationInput,
            $0.KinesisStreamingDestinationOutput>(
        'DisableKinesisStreamingDestination',
        disableKinesisStreamingDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.KinesisStreamingDestinationInput.fromBuffer(value),
        ($0.KinesisStreamingDestinationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.KinesisStreamingDestinationInput,
            $0.KinesisStreamingDestinationOutput>(
        'EnableKinesisStreamingDestination',
        enableKinesisStreamingDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.KinesisStreamingDestinationInput.fromBuffer(value),
        ($0.KinesisStreamingDestinationOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ExecuteStatementInput,
            $0.ExecuteStatementOutput>(
        'ExecuteStatement',
        executeStatement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ExecuteStatementInput.fromBuffer(value),
        ($0.ExecuteStatementOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ExecuteTransactionInput,
            $0.ExecuteTransactionOutput>(
        'ExecuteTransaction',
        executeTransaction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ExecuteTransactionInput.fromBuffer(value),
        ($0.ExecuteTransactionOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ExportTableToPointInTimeInput,
            $0.ExportTableToPointInTimeOutput>(
        'ExportTableToPointInTime',
        exportTableToPointInTime_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ExportTableToPointInTimeInput.fromBuffer(value),
        ($0.ExportTableToPointInTimeOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetItemInput, $0.GetItemOutput>(
        'GetItem',
        getItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetItemInput.fromBuffer(value),
        ($0.GetItemOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcePolicyInput,
            $0.GetResourcePolicyOutput>(
        'GetResourcePolicy',
        getResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcePolicyInput.fromBuffer(value),
        ($0.GetResourcePolicyOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ImportTableInput, $0.ImportTableOutput>(
        'ImportTable',
        importTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ImportTableInput.fromBuffer(value),
        ($0.ImportTableOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListBackupsInput, $0.ListBackupsOutput>(
        'ListBackups',
        listBackups_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListBackupsInput.fromBuffer(value),
        ($0.ListBackupsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListContributorInsightsInput,
            $0.ListContributorInsightsOutput>(
        'ListContributorInsights',
        listContributorInsights_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListContributorInsightsInput.fromBuffer(value),
        ($0.ListContributorInsightsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListExportsInput, $0.ListExportsOutput>(
        'ListExports',
        listExports_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListExportsInput.fromBuffer(value),
        ($0.ListExportsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGlobalTablesInput,
            $0.ListGlobalTablesOutput>(
        'ListGlobalTables',
        listGlobalTables_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListGlobalTablesInput.fromBuffer(value),
        ($0.ListGlobalTablesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListImportsInput, $0.ListImportsOutput>(
        'ListImports',
        listImports_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListImportsInput.fromBuffer(value),
        ($0.ListImportsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTablesInput, $0.ListTablesOutput>(
        'ListTables',
        listTables_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListTablesInput.fromBuffer(value),
        ($0.ListTablesOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsOfResourceInput,
            $0.ListTagsOfResourceOutput>(
        'ListTagsOfResource',
        listTagsOfResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsOfResourceInput.fromBuffer(value),
        ($0.ListTagsOfResourceOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutItemInput, $0.PutItemOutput>(
        'PutItem',
        putItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.PutItemInput.fromBuffer(value),
        ($0.PutItemOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutResourcePolicyInput,
            $0.PutResourcePolicyOutput>(
        'PutResourcePolicy',
        putResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutResourcePolicyInput.fromBuffer(value),
        ($0.PutResourcePolicyOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.QueryInput, $0.QueryOutput>(
        'Query',
        query_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.QueryInput.fromBuffer(value),
        ($0.QueryOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RestoreTableFromBackupInput,
            $0.RestoreTableFromBackupOutput>(
        'RestoreTableFromBackup',
        restoreTableFromBackup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RestoreTableFromBackupInput.fromBuffer(value),
        ($0.RestoreTableFromBackupOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RestoreTableToPointInTimeInput,
            $0.RestoreTableToPointInTimeOutput>(
        'RestoreTableToPointInTime',
        restoreTableToPointInTime_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RestoreTableToPointInTimeInput.fromBuffer(value),
        ($0.RestoreTableToPointInTimeOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ScanInput, $0.ScanOutput>(
        'Scan',
        scan_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ScanInput.fromBuffer(value),
        ($0.ScanOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceInput, $1.Empty>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagResourceInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TransactGetItemsInput,
            $0.TransactGetItemsOutput>(
        'TransactGetItems',
        transactGetItems_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TransactGetItemsInput.fromBuffer(value),
        ($0.TransactGetItemsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TransactWriteItemsInput,
            $0.TransactWriteItemsOutput>(
        'TransactWriteItems',
        transactWriteItems_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TransactWriteItemsInput.fromBuffer(value),
        ($0.TransactWriteItemsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagResourceInput, $1.Empty>(
        'UntagResource',
        untagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagResourceInput.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateContinuousBackupsInput,
            $0.UpdateContinuousBackupsOutput>(
        'UpdateContinuousBackups',
        updateContinuousBackups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateContinuousBackupsInput.fromBuffer(value),
        ($0.UpdateContinuousBackupsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateContributorInsightsInput,
            $0.UpdateContributorInsightsOutput>(
        'UpdateContributorInsights',
        updateContributorInsights_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateContributorInsightsInput.fromBuffer(value),
        ($0.UpdateContributorInsightsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateGlobalTableInput,
            $0.UpdateGlobalTableOutput>(
        'UpdateGlobalTable',
        updateGlobalTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateGlobalTableInput.fromBuffer(value),
        ($0.UpdateGlobalTableOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateGlobalTableSettingsInput,
            $0.UpdateGlobalTableSettingsOutput>(
        'UpdateGlobalTableSettings',
        updateGlobalTableSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateGlobalTableSettingsInput.fromBuffer(value),
        ($0.UpdateGlobalTableSettingsOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateItemInput, $0.UpdateItemOutput>(
        'UpdateItem',
        updateItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateItemInput.fromBuffer(value),
        ($0.UpdateItemOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateKinesisStreamingDestinationInput,
            $0.UpdateKinesisStreamingDestinationOutput>(
        'UpdateKinesisStreamingDestination',
        updateKinesisStreamingDestination_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateKinesisStreamingDestinationInput.fromBuffer(value),
        ($0.UpdateKinesisStreamingDestinationOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateTableInput, $0.UpdateTableOutput>(
        'UpdateTable',
        updateTable_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateTableInput.fromBuffer(value),
        ($0.UpdateTableOutput value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateTableReplicaAutoScalingInput,
            $0.UpdateTableReplicaAutoScalingOutput>(
        'UpdateTableReplicaAutoScaling',
        updateTableReplicaAutoScaling_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateTableReplicaAutoScalingInput.fromBuffer(value),
        ($0.UpdateTableReplicaAutoScalingOutput value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateTimeToLiveInput,
            $0.UpdateTimeToLiveOutput>(
        'UpdateTimeToLive',
        updateTimeToLive_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateTimeToLiveInput.fromBuffer(value),
        ($0.UpdateTimeToLiveOutput value) => value.writeToBuffer()));
  }

  $async.Future<$0.BatchExecuteStatementOutput> batchExecuteStatement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.BatchExecuteStatementInput> $request) async {
    return batchExecuteStatement($call, await $request);
  }

  $async.Future<$0.BatchExecuteStatementOutput> batchExecuteStatement(
      $grpc.ServiceCall call, $0.BatchExecuteStatementInput request);

  $async.Future<$0.BatchGetItemOutput> batchGetItem_Pre($grpc.ServiceCall $call,
      $async.Future<$0.BatchGetItemInput> $request) async {
    return batchGetItem($call, await $request);
  }

  $async.Future<$0.BatchGetItemOutput> batchGetItem(
      $grpc.ServiceCall call, $0.BatchGetItemInput request);

  $async.Future<$0.BatchWriteItemOutput> batchWriteItem_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.BatchWriteItemInput> $request) async {
    return batchWriteItem($call, await $request);
  }

  $async.Future<$0.BatchWriteItemOutput> batchWriteItem(
      $grpc.ServiceCall call, $0.BatchWriteItemInput request);

  $async.Future<$0.CreateBackupOutput> createBackup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateBackupInput> $request) async {
    return createBackup($call, await $request);
  }

  $async.Future<$0.CreateBackupOutput> createBackup(
      $grpc.ServiceCall call, $0.CreateBackupInput request);

  $async.Future<$0.CreateGlobalTableOutput> createGlobalTable_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateGlobalTableInput> $request) async {
    return createGlobalTable($call, await $request);
  }

  $async.Future<$0.CreateGlobalTableOutput> createGlobalTable(
      $grpc.ServiceCall call, $0.CreateGlobalTableInput request);

  $async.Future<$0.CreateTableOutput> createTable_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateTableInput> $request) async {
    return createTable($call, await $request);
  }

  $async.Future<$0.CreateTableOutput> createTable(
      $grpc.ServiceCall call, $0.CreateTableInput request);

  $async.Future<$0.DeleteBackupOutput> deleteBackup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteBackupInput> $request) async {
    return deleteBackup($call, await $request);
  }

  $async.Future<$0.DeleteBackupOutput> deleteBackup(
      $grpc.ServiceCall call, $0.DeleteBackupInput request);

  $async.Future<$0.DeleteItemOutput> deleteItem_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteItemInput> $request) async {
    return deleteItem($call, await $request);
  }

  $async.Future<$0.DeleteItemOutput> deleteItem(
      $grpc.ServiceCall call, $0.DeleteItemInput request);

  $async.Future<$0.DeleteResourcePolicyOutput> deleteResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourcePolicyInput> $request) async {
    return deleteResourcePolicy($call, await $request);
  }

  $async.Future<$0.DeleteResourcePolicyOutput> deleteResourcePolicy(
      $grpc.ServiceCall call, $0.DeleteResourcePolicyInput request);

  $async.Future<$0.DeleteTableOutput> deleteTable_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteTableInput> $request) async {
    return deleteTable($call, await $request);
  }

  $async.Future<$0.DeleteTableOutput> deleteTable(
      $grpc.ServiceCall call, $0.DeleteTableInput request);

  $async.Future<$0.DescribeBackupOutput> describeBackup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeBackupInput> $request) async {
    return describeBackup($call, await $request);
  }

  $async.Future<$0.DescribeBackupOutput> describeBackup(
      $grpc.ServiceCall call, $0.DescribeBackupInput request);

  $async.Future<$0.DescribeContinuousBackupsOutput>
      describeContinuousBackups_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeContinuousBackupsInput> $request) async {
    return describeContinuousBackups($call, await $request);
  }

  $async.Future<$0.DescribeContinuousBackupsOutput> describeContinuousBackups(
      $grpc.ServiceCall call, $0.DescribeContinuousBackupsInput request);

  $async.Future<$0.DescribeContributorInsightsOutput>
      describeContributorInsights_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeContributorInsightsInput> $request) async {
    return describeContributorInsights($call, await $request);
  }

  $async.Future<$0.DescribeContributorInsightsOutput>
      describeContributorInsights(
          $grpc.ServiceCall call, $0.DescribeContributorInsightsInput request);

  $async.Future<$0.DescribeEndpointsResponse> describeEndpoints_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeEndpointsRequest> $request) async {
    return describeEndpoints($call, await $request);
  }

  $async.Future<$0.DescribeEndpointsResponse> describeEndpoints(
      $grpc.ServiceCall call, $0.DescribeEndpointsRequest request);

  $async.Future<$0.DescribeExportOutput> describeExport_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeExportInput> $request) async {
    return describeExport($call, await $request);
  }

  $async.Future<$0.DescribeExportOutput> describeExport(
      $grpc.ServiceCall call, $0.DescribeExportInput request);

  $async.Future<$0.DescribeGlobalTableOutput> describeGlobalTable_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeGlobalTableInput> $request) async {
    return describeGlobalTable($call, await $request);
  }

  $async.Future<$0.DescribeGlobalTableOutput> describeGlobalTable(
      $grpc.ServiceCall call, $0.DescribeGlobalTableInput request);

  $async.Future<$0.DescribeGlobalTableSettingsOutput>
      describeGlobalTableSettings_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeGlobalTableSettingsInput> $request) async {
    return describeGlobalTableSettings($call, await $request);
  }

  $async.Future<$0.DescribeGlobalTableSettingsOutput>
      describeGlobalTableSettings(
          $grpc.ServiceCall call, $0.DescribeGlobalTableSettingsInput request);

  $async.Future<$0.DescribeImportOutput> describeImport_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeImportInput> $request) async {
    return describeImport($call, await $request);
  }

  $async.Future<$0.DescribeImportOutput> describeImport(
      $grpc.ServiceCall call, $0.DescribeImportInput request);

  $async.Future<$0.DescribeKinesisStreamingDestinationOutput>
      describeKinesisStreamingDestination_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeKinesisStreamingDestinationInput>
              $request) async {
    return describeKinesisStreamingDestination($call, await $request);
  }

  $async.Future<$0.DescribeKinesisStreamingDestinationOutput>
      describeKinesisStreamingDestination($grpc.ServiceCall call,
          $0.DescribeKinesisStreamingDestinationInput request);

  $async.Future<$0.DescribeLimitsOutput> describeLimits_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeLimitsInput> $request) async {
    return describeLimits($call, await $request);
  }

  $async.Future<$0.DescribeLimitsOutput> describeLimits(
      $grpc.ServiceCall call, $0.DescribeLimitsInput request);

  $async.Future<$0.DescribeTableOutput> describeTable_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeTableInput> $request) async {
    return describeTable($call, await $request);
  }

  $async.Future<$0.DescribeTableOutput> describeTable(
      $grpc.ServiceCall call, $0.DescribeTableInput request);

  $async.Future<$0.DescribeTableReplicaAutoScalingOutput>
      describeTableReplicaAutoScaling_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeTableReplicaAutoScalingInput>
              $request) async {
    return describeTableReplicaAutoScaling($call, await $request);
  }

  $async.Future<$0.DescribeTableReplicaAutoScalingOutput>
      describeTableReplicaAutoScaling($grpc.ServiceCall call,
          $0.DescribeTableReplicaAutoScalingInput request);

  $async.Future<$0.DescribeTimeToLiveOutput> describeTimeToLive_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeTimeToLiveInput> $request) async {
    return describeTimeToLive($call, await $request);
  }

  $async.Future<$0.DescribeTimeToLiveOutput> describeTimeToLive(
      $grpc.ServiceCall call, $0.DescribeTimeToLiveInput request);

  $async.Future<$0.KinesisStreamingDestinationOutput>
      disableKinesisStreamingDestination_Pre($grpc.ServiceCall $call,
          $async.Future<$0.KinesisStreamingDestinationInput> $request) async {
    return disableKinesisStreamingDestination($call, await $request);
  }

  $async.Future<$0.KinesisStreamingDestinationOutput>
      disableKinesisStreamingDestination(
          $grpc.ServiceCall call, $0.KinesisStreamingDestinationInput request);

  $async.Future<$0.KinesisStreamingDestinationOutput>
      enableKinesisStreamingDestination_Pre($grpc.ServiceCall $call,
          $async.Future<$0.KinesisStreamingDestinationInput> $request) async {
    return enableKinesisStreamingDestination($call, await $request);
  }

  $async.Future<$0.KinesisStreamingDestinationOutput>
      enableKinesisStreamingDestination(
          $grpc.ServiceCall call, $0.KinesisStreamingDestinationInput request);

  $async.Future<$0.ExecuteStatementOutput> executeStatement_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExecuteStatementInput> $request) async {
    return executeStatement($call, await $request);
  }

  $async.Future<$0.ExecuteStatementOutput> executeStatement(
      $grpc.ServiceCall call, $0.ExecuteStatementInput request);

  $async.Future<$0.ExecuteTransactionOutput> executeTransaction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExecuteTransactionInput> $request) async {
    return executeTransaction($call, await $request);
  }

  $async.Future<$0.ExecuteTransactionOutput> executeTransaction(
      $grpc.ServiceCall call, $0.ExecuteTransactionInput request);

  $async.Future<$0.ExportTableToPointInTimeOutput> exportTableToPointInTime_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExportTableToPointInTimeInput> $request) async {
    return exportTableToPointInTime($call, await $request);
  }

  $async.Future<$0.ExportTableToPointInTimeOutput> exportTableToPointInTime(
      $grpc.ServiceCall call, $0.ExportTableToPointInTimeInput request);

  $async.Future<$0.GetItemOutput> getItem_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.GetItemInput> $request) async {
    return getItem($call, await $request);
  }

  $async.Future<$0.GetItemOutput> getItem(
      $grpc.ServiceCall call, $0.GetItemInput request);

  $async.Future<$0.GetResourcePolicyOutput> getResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourcePolicyInput> $request) async {
    return getResourcePolicy($call, await $request);
  }

  $async.Future<$0.GetResourcePolicyOutput> getResourcePolicy(
      $grpc.ServiceCall call, $0.GetResourcePolicyInput request);

  $async.Future<$0.ImportTableOutput> importTable_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ImportTableInput> $request) async {
    return importTable($call, await $request);
  }

  $async.Future<$0.ImportTableOutput> importTable(
      $grpc.ServiceCall call, $0.ImportTableInput request);

  $async.Future<$0.ListBackupsOutput> listBackups_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListBackupsInput> $request) async {
    return listBackups($call, await $request);
  }

  $async.Future<$0.ListBackupsOutput> listBackups(
      $grpc.ServiceCall call, $0.ListBackupsInput request);

  $async.Future<$0.ListContributorInsightsOutput> listContributorInsights_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListContributorInsightsInput> $request) async {
    return listContributorInsights($call, await $request);
  }

  $async.Future<$0.ListContributorInsightsOutput> listContributorInsights(
      $grpc.ServiceCall call, $0.ListContributorInsightsInput request);

  $async.Future<$0.ListExportsOutput> listExports_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListExportsInput> $request) async {
    return listExports($call, await $request);
  }

  $async.Future<$0.ListExportsOutput> listExports(
      $grpc.ServiceCall call, $0.ListExportsInput request);

  $async.Future<$0.ListGlobalTablesOutput> listGlobalTables_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListGlobalTablesInput> $request) async {
    return listGlobalTables($call, await $request);
  }

  $async.Future<$0.ListGlobalTablesOutput> listGlobalTables(
      $grpc.ServiceCall call, $0.ListGlobalTablesInput request);

  $async.Future<$0.ListImportsOutput> listImports_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListImportsInput> $request) async {
    return listImports($call, await $request);
  }

  $async.Future<$0.ListImportsOutput> listImports(
      $grpc.ServiceCall call, $0.ListImportsInput request);

  $async.Future<$0.ListTablesOutput> listTables_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTablesInput> $request) async {
    return listTables($call, await $request);
  }

  $async.Future<$0.ListTablesOutput> listTables(
      $grpc.ServiceCall call, $0.ListTablesInput request);

  $async.Future<$0.ListTagsOfResourceOutput> listTagsOfResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsOfResourceInput> $request) async {
    return listTagsOfResource($call, await $request);
  }

  $async.Future<$0.ListTagsOfResourceOutput> listTagsOfResource(
      $grpc.ServiceCall call, $0.ListTagsOfResourceInput request);

  $async.Future<$0.PutItemOutput> putItem_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.PutItemInput> $request) async {
    return putItem($call, await $request);
  }

  $async.Future<$0.PutItemOutput> putItem(
      $grpc.ServiceCall call, $0.PutItemInput request);

  $async.Future<$0.PutResourcePolicyOutput> putResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutResourcePolicyInput> $request) async {
    return putResourcePolicy($call, await $request);
  }

  $async.Future<$0.PutResourcePolicyOutput> putResourcePolicy(
      $grpc.ServiceCall call, $0.PutResourcePolicyInput request);

  $async.Future<$0.QueryOutput> query_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.QueryInput> $request) async {
    return query($call, await $request);
  }

  $async.Future<$0.QueryOutput> query(
      $grpc.ServiceCall call, $0.QueryInput request);

  $async.Future<$0.RestoreTableFromBackupOutput> restoreTableFromBackup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RestoreTableFromBackupInput> $request) async {
    return restoreTableFromBackup($call, await $request);
  }

  $async.Future<$0.RestoreTableFromBackupOutput> restoreTableFromBackup(
      $grpc.ServiceCall call, $0.RestoreTableFromBackupInput request);

  $async.Future<$0.RestoreTableToPointInTimeOutput>
      restoreTableToPointInTime_Pre($grpc.ServiceCall $call,
          $async.Future<$0.RestoreTableToPointInTimeInput> $request) async {
    return restoreTableToPointInTime($call, await $request);
  }

  $async.Future<$0.RestoreTableToPointInTimeOutput> restoreTableToPointInTime(
      $grpc.ServiceCall call, $0.RestoreTableToPointInTimeInput request);

  $async.Future<$0.ScanOutput> scan_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.ScanInput> $request) async {
    return scan($call, await $request);
  }

  $async.Future<$0.ScanOutput> scan(
      $grpc.ServiceCall call, $0.ScanInput request);

  $async.Future<$1.Empty> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceInput> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$1.Empty> tagResource(
      $grpc.ServiceCall call, $0.TagResourceInput request);

  $async.Future<$0.TransactGetItemsOutput> transactGetItems_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TransactGetItemsInput> $request) async {
    return transactGetItems($call, await $request);
  }

  $async.Future<$0.TransactGetItemsOutput> transactGetItems(
      $grpc.ServiceCall call, $0.TransactGetItemsInput request);

  $async.Future<$0.TransactWriteItemsOutput> transactWriteItems_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TransactWriteItemsInput> $request) async {
    return transactWriteItems($call, await $request);
  }

  $async.Future<$0.TransactWriteItemsOutput> transactWriteItems(
      $grpc.ServiceCall call, $0.TransactWriteItemsInput request);

  $async.Future<$1.Empty> untagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceInput> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$1.Empty> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceInput request);

  $async.Future<$0.UpdateContinuousBackupsOutput> updateContinuousBackups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateContinuousBackupsInput> $request) async {
    return updateContinuousBackups($call, await $request);
  }

  $async.Future<$0.UpdateContinuousBackupsOutput> updateContinuousBackups(
      $grpc.ServiceCall call, $0.UpdateContinuousBackupsInput request);

  $async.Future<$0.UpdateContributorInsightsOutput>
      updateContributorInsights_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateContributorInsightsInput> $request) async {
    return updateContributorInsights($call, await $request);
  }

  $async.Future<$0.UpdateContributorInsightsOutput> updateContributorInsights(
      $grpc.ServiceCall call, $0.UpdateContributorInsightsInput request);

  $async.Future<$0.UpdateGlobalTableOutput> updateGlobalTable_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateGlobalTableInput> $request) async {
    return updateGlobalTable($call, await $request);
  }

  $async.Future<$0.UpdateGlobalTableOutput> updateGlobalTable(
      $grpc.ServiceCall call, $0.UpdateGlobalTableInput request);

  $async.Future<$0.UpdateGlobalTableSettingsOutput>
      updateGlobalTableSettings_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateGlobalTableSettingsInput> $request) async {
    return updateGlobalTableSettings($call, await $request);
  }

  $async.Future<$0.UpdateGlobalTableSettingsOutput> updateGlobalTableSettings(
      $grpc.ServiceCall call, $0.UpdateGlobalTableSettingsInput request);

  $async.Future<$0.UpdateItemOutput> updateItem_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateItemInput> $request) async {
    return updateItem($call, await $request);
  }

  $async.Future<$0.UpdateItemOutput> updateItem(
      $grpc.ServiceCall call, $0.UpdateItemInput request);

  $async.Future<$0.UpdateKinesisStreamingDestinationOutput>
      updateKinesisStreamingDestination_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateKinesisStreamingDestinationInput>
              $request) async {
    return updateKinesisStreamingDestination($call, await $request);
  }

  $async.Future<$0.UpdateKinesisStreamingDestinationOutput>
      updateKinesisStreamingDestination($grpc.ServiceCall call,
          $0.UpdateKinesisStreamingDestinationInput request);

  $async.Future<$0.UpdateTableOutput> updateTable_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateTableInput> $request) async {
    return updateTable($call, await $request);
  }

  $async.Future<$0.UpdateTableOutput> updateTable(
      $grpc.ServiceCall call, $0.UpdateTableInput request);

  $async.Future<$0.UpdateTableReplicaAutoScalingOutput>
      updateTableReplicaAutoScaling_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateTableReplicaAutoScalingInput> $request) async {
    return updateTableReplicaAutoScaling($call, await $request);
  }

  $async.Future<$0.UpdateTableReplicaAutoScalingOutput>
      updateTableReplicaAutoScaling($grpc.ServiceCall call,
          $0.UpdateTableReplicaAutoScalingInput request);

  $async.Future<$0.UpdateTimeToLiveOutput> updateTimeToLive_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateTimeToLiveInput> $request) async {
    return updateTimeToLive($call, await $request);
  }

  $async.Future<$0.UpdateTimeToLiveOutput> updateTimeToLive(
      $grpc.ServiceCall call, $0.UpdateTimeToLiveInput request);
}
