// This is a generated file - do not edit.
//
// Generated from ssm.proto.

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

import 'ssm.pb.dart' as $0;

export 'ssm.pb.dart';

/// SSMService provides ssm API operations.
@$pb.GrpcServiceName('ssm.SSMService')
class SSMServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  SSMServiceClient(super.channel, {super.options, super.interceptors});

  /// Adds or overwrites one or more tags for the specified resource. Tags are metadata that you can assign to your automations, documents, managed nodes, maintenance windows, Parameter Store parameters,...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AddTagsToResourceResult> addTagsToResource(
    $0.AddTagsToResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addTagsToResource, request, options: options);
  }

  /// Associates a related item to a Systems Manager OpsCenter OpsItem. For example, you can associate an Incident Manager incident or analysis with an OpsItem. Incident Manager and OpsCenter are tools i...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AssociateOpsItemRelatedItemResponse>
      associateOpsItemRelatedItem(
    $0.AssociateOpsItemRelatedItemRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateOpsItemRelatedItem, request,
        options: options);
  }

  /// Attempts to cancel the command specified by the Command ID. There is no guarantee that the command will be terminated and the underlying process stopped.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelCommandResult> cancelCommand(
    $0.CancelCommandRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelCommand, request, options: options);
  }

  /// Stops a maintenance window execution that is already in progress and cancels any tasks in the window that haven't already starting running. Tasks already in progress will continue to completion.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CancelMaintenanceWindowExecutionResult>
      cancelMaintenanceWindowExecution(
    $0.CancelMaintenanceWindowExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$cancelMaintenanceWindowExecution, request,
        options: options);
  }

  /// Generates an activation code and activation ID you can use to register your on-premises servers, edge devices, or virtual machine (VM) with Amazon Web Services Systems Manager. Registering these ma...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateActivationResult> createActivation(
    $0.CreateActivationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createActivation, request, options: options);
  }

  /// A State Manager association defines the state that you want to maintain on your managed nodes. For example, an association can specify that anti-virus software must be installed and running on your...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateAssociationResult> createAssociation(
    $0.CreateAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAssociation, request, options: options);
  }

  /// Associates the specified Amazon Web Services Systems Manager document (SSM document) with the specified managed nodes or targets. When you associate a document with one or more managed nodes using ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateAssociationBatchResult> createAssociationBatch(
    $0.CreateAssociationBatchRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAssociationBatch, request,
        options: options);
  }

  /// Creates a Amazon Web Services Systems Manager (SSM document). An SSM document defines the actions that Systems Manager performs on your managed nodes. For more information about SSM documents, incl...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateDocumentResult> createDocument(
    $0.CreateDocumentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDocument, request, options: options);
  }

  /// Creates a new maintenance window. The value you specify for Duration determines the specific end time for the maintenance window based on the time it begins. No maintenance window tasks are permitt...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateMaintenanceWindowResult>
      createMaintenanceWindow(
    $0.CreateMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createMaintenanceWindow, request,
        options: options);
  }

  /// Creates a new OpsItem. You must have permission in Identity and Access Management (IAM) to create a new OpsItem. For more information, see Set up OpsCenter in the Amazon Web Services Systems Manage...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateOpsItemResponse> createOpsItem(
    $0.CreateOpsItemRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createOpsItem, request, options: options);
  }

  /// If you create a new application in Application Manager, Amazon Web Services Systems Manager calls this API operation to specify information about the new application, including the application type.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateOpsMetadataResult> createOpsMetadata(
    $0.CreateOpsMetadataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createOpsMetadata, request, options: options);
  }

  /// Creates a patch baseline. For information about valid key-value pairs in PatchFilters for each supported operating system type, see PatchFilter.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreatePatchBaselineResult> createPatchBaseline(
    $0.CreatePatchBaselineRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPatchBaseline, request, options: options);
  }

  /// A resource data sync helps you view data from multiple sources in a single location. Amazon Web Services Systems Manager offers two types of resource data sync: SyncToDestination and SyncFromSource...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateResourceDataSyncResult> createResourceDataSync(
    $0.CreateResourceDataSyncRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createResourceDataSync, request,
        options: options);
  }

  /// Deletes an activation. You aren't required to delete an activation. If you delete an activation, you can no longer use it to register additional managed nodes. Deleting an activation doesn't de-reg...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteActivationResult> deleteActivation(
    $0.DeleteActivationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteActivation, request, options: options);
  }

  /// Disassociates the specified Amazon Web Services Systems Manager document (SSM document) from the specified managed node. If you created the association by using the Targets parameter, then you must...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteAssociationResult> deleteAssociation(
    $0.DeleteAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAssociation, request, options: options);
  }

  /// Deletes the Amazon Web Services Systems Manager document (SSM document) and all managed node associations to the document. Before you delete the document, we recommend that you use DeleteAssociatio...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteDocumentResult> deleteDocument(
    $0.DeleteDocumentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteDocument, request, options: options);
  }

  /// Delete a custom inventory type or the data associated with a custom Inventory type. Deleting a custom inventory type is also referred to as deleting a custom inventory schema.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteInventoryResult> deleteInventory(
    $0.DeleteInventoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteInventory, request, options: options);
  }

  /// Deletes a maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteMaintenanceWindowResult>
      deleteMaintenanceWindow(
    $0.DeleteMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteMaintenanceWindow, request,
        options: options);
  }

  /// Delete an OpsItem. You must have permission in Identity and Access Management (IAM) to delete an OpsItem. Note the following important information about this operation. Deleting an OpsItem is irrev...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteOpsItemResponse> deleteOpsItem(
    $0.DeleteOpsItemRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteOpsItem, request, options: options);
  }

  /// Delete OpsMetadata related to an application.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteOpsMetadataResult> deleteOpsMetadata(
    $0.DeleteOpsMetadataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteOpsMetadata, request, options: options);
  }

  /// Delete a parameter from the system. After deleting a parameter, wait for at least 30 seconds to create a parameter with the same name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteParameterResult> deleteParameter(
    $0.DeleteParameterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteParameter, request, options: options);
  }

  /// Delete a list of parameters. After deleting a parameter, wait for at least 30 seconds to create a parameter with the same name.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteParametersResult> deleteParameters(
    $0.DeleteParametersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteParameters, request, options: options);
  }

  /// Deletes a patch baseline.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeletePatchBaselineResult> deletePatchBaseline(
    $0.DeletePatchBaselineRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePatchBaseline, request, options: options);
  }

  /// Deletes a resource data sync configuration. After the configuration is deleted, changes to data on managed nodes are no longer synced to or from the target. Deleting a sync configuration doesn't de...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteResourceDataSyncResult> deleteResourceDataSync(
    $0.DeleteResourceDataSyncRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourceDataSync, request,
        options: options);
  }

  /// Deletes a Systems Manager resource policy. A resource policy helps you to define the IAM entity (for example, an Amazon Web Services account) that can manage your Systems Manager resources. The fol...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteResourcePolicyResponse> deleteResourcePolicy(
    $0.DeleteResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourcePolicy, request, options: options);
  }

  /// Removes the server or virtual machine from the list of registered servers. If you want to reregister an on-premises server, edge device, or VM, you must use a different Activation Code and Activati...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeregisterManagedInstanceResult>
      deregisterManagedInstance(
    $0.DeregisterManagedInstanceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deregisterManagedInstance, request,
        options: options);
  }

  /// Removes a patch group from a patch baseline.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeregisterPatchBaselineForPatchGroupResult>
      deregisterPatchBaselineForPatchGroup(
    $0.DeregisterPatchBaselineForPatchGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deregisterPatchBaselineForPatchGroup, request,
        options: options);
  }

  /// Removes a target from a maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeregisterTargetFromMaintenanceWindowResult>
      deregisterTargetFromMaintenanceWindow(
    $0.DeregisterTargetFromMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deregisterTargetFromMaintenanceWindow, request,
        options: options);
  }

  /// Removes a task from a maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeregisterTaskFromMaintenanceWindowResult>
      deregisterTaskFromMaintenanceWindow(
    $0.DeregisterTaskFromMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deregisterTaskFromMaintenanceWindow, request,
        options: options);
  }

  /// Describes details about the activation, such as the date and time the activation was created, its expiration date, the Identity and Access Management (IAM) role assigned to the managed nodes in the...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeActivationsResult> describeActivations(
    $0.DescribeActivationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeActivations, request, options: options);
  }

  /// Describes the association for the specified target or managed node. If you created the association by using the Targets parameter, then you must retrieve the association by using the association ID.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAssociationResult> describeAssociation(
    $0.DescribeAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAssociation, request, options: options);
  }

  /// Views all executions for a specific association ID.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAssociationExecutionsResult>
      describeAssociationExecutions(
    $0.DescribeAssociationExecutionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAssociationExecutions, request,
        options: options);
  }

  /// Views information about a specific execution of a specific association.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAssociationExecutionTargetsResult>
      describeAssociationExecutionTargets(
    $0.DescribeAssociationExecutionTargetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAssociationExecutionTargets, request,
        options: options);
  }

  /// Provides details about all active and terminated Automation executions.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAutomationExecutionsResult>
      describeAutomationExecutions(
    $0.DescribeAutomationExecutionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAutomationExecutions, request,
        options: options);
  }

  /// Information about all active and terminated step executions in an Automation workflow.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAutomationStepExecutionsResult>
      describeAutomationStepExecutions(
    $0.DescribeAutomationStepExecutionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAutomationStepExecutions, request,
        options: options);
  }

  /// Lists all patches eligible to be included in a patch baseline. Currently, DescribeAvailablePatches supports only the Amazon Linux 1, Amazon Linux 2, and Windows Server operating systems.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeAvailablePatchesResult>
      describeAvailablePatches(
    $0.DescribeAvailablePatchesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeAvailablePatches, request,
        options: options);
  }

  /// Describes the specified Amazon Web Services Systems Manager document (SSM document).
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeDocumentResult> describeDocument(
    $0.DescribeDocumentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeDocument, request, options: options);
  }

  /// Describes the permissions for a Amazon Web Services Systems Manager document (SSM document). If you created the document, you are the owner. If a document is shared, it can either be shared private...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeDocumentPermissionResponse>
      describeDocumentPermission(
    $0.DescribeDocumentPermissionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeDocumentPermission, request,
        options: options);
  }

  /// All associations for the managed nodes.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeEffectiveInstanceAssociationsResult>
      describeEffectiveInstanceAssociations(
    $0.DescribeEffectiveInstanceAssociationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeEffectiveInstanceAssociations, request,
        options: options);
  }

  /// Retrieves the current effective patches (the patch and the approval state) for the specified patch baseline. Applies to patch baselines for Windows only.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeEffectivePatchesForPatchBaselineResult>
      describeEffectivePatchesForPatchBaseline(
    $0.DescribeEffectivePatchesForPatchBaselineRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeEffectivePatchesForPatchBaseline, request,
        options: options);
  }

  /// The status of the associations for the managed nodes.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeInstanceAssociationsStatusResult>
      describeInstanceAssociationsStatus(
    $0.DescribeInstanceAssociationsStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInstanceAssociationsStatus, request,
        options: options);
  }

  /// Provides information about one or more of your managed nodes, including the operating system platform, SSM Agent version, association status, and IP address. This operation does not return informat...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeInstanceInformationResult>
      describeInstanceInformation(
    $0.DescribeInstanceInformationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInstanceInformation, request,
        options: options);
  }

  /// Retrieves information about the patches on the specified managed node and their state relative to the patch baseline being used for the node.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeInstancePatchesResult>
      describeInstancePatches(
    $0.DescribeInstancePatchesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInstancePatches, request,
        options: options);
  }

  /// Retrieves the high-level patch state of one or more managed nodes.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeInstancePatchStatesResult>
      describeInstancePatchStates(
    $0.DescribeInstancePatchStatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInstancePatchStates, request,
        options: options);
  }

  /// Retrieves the high-level patch state for the managed nodes in the specified patch group.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeInstancePatchStatesForPatchGroupResult>
      describeInstancePatchStatesForPatchGroup(
    $0.DescribeInstancePatchStatesForPatchGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInstancePatchStatesForPatchGroup, request,
        options: options);
  }

  /// An API operation used by the Systems Manager console to display information about Systems Manager managed nodes.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeInstancePropertiesResult>
      describeInstanceProperties(
    $0.DescribeInstancePropertiesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInstanceProperties, request,
        options: options);
  }

  /// Describes a specific delete inventory operation.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeInventoryDeletionsResult>
      describeInventoryDeletions(
    $0.DescribeInventoryDeletionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeInventoryDeletions, request,
        options: options);
  }

  /// Lists the executions of a maintenance window. This includes information about when the maintenance window was scheduled to be active, and information about tasks registered and run with the mainten...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMaintenanceWindowExecutionsResult>
      describeMaintenanceWindowExecutions(
    $0.DescribeMaintenanceWindowExecutionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMaintenanceWindowExecutions, request,
        options: options);
  }

  /// Retrieves the individual task executions (one per target) for a particular task run as part of a maintenance window execution.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<
          $0.DescribeMaintenanceWindowExecutionTaskInvocationsResult>
      describeMaintenanceWindowExecutionTaskInvocations(
    $0.DescribeMaintenanceWindowExecutionTaskInvocationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$describeMaintenanceWindowExecutionTaskInvocations, request,
        options: options);
  }

  /// For a given maintenance window execution, lists the tasks that were run.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMaintenanceWindowExecutionTasksResult>
      describeMaintenanceWindowExecutionTasks(
    $0.DescribeMaintenanceWindowExecutionTasksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMaintenanceWindowExecutionTasks, request,
        options: options);
  }

  /// Retrieves the maintenance windows in an Amazon Web Services account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMaintenanceWindowsResult>
      describeMaintenanceWindows(
    $0.DescribeMaintenanceWindowsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMaintenanceWindows, request,
        options: options);
  }

  /// Retrieves information about upcoming executions of a maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMaintenanceWindowScheduleResult>
      describeMaintenanceWindowSchedule(
    $0.DescribeMaintenanceWindowScheduleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMaintenanceWindowSchedule, request,
        options: options);
  }

  /// Retrieves information about the maintenance window targets or tasks that a managed node is associated with.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMaintenanceWindowsForTargetResult>
      describeMaintenanceWindowsForTarget(
    $0.DescribeMaintenanceWindowsForTargetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMaintenanceWindowsForTarget, request,
        options: options);
  }

  /// Lists the targets registered with the maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMaintenanceWindowTargetsResult>
      describeMaintenanceWindowTargets(
    $0.DescribeMaintenanceWindowTargetsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMaintenanceWindowTargets, request,
        options: options);
  }

  /// Lists the tasks in a maintenance window. For maintenance window tasks without a specified target, you can't supply values for --max-errors and --max-concurrency. Instead, the system inserts a place...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeMaintenanceWindowTasksResult>
      describeMaintenanceWindowTasks(
    $0.DescribeMaintenanceWindowTasksRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeMaintenanceWindowTasks, request,
        options: options);
  }

  /// Query a set of OpsItems. You must have permission in Identity and Access Management (IAM) to query a list of OpsItems. For more information, see Set up OpsCenter in the Amazon Web Services Systems ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeOpsItemsResponse> describeOpsItems(
    $0.DescribeOpsItemsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeOpsItems, request, options: options);
  }

  /// Lists the parameters in your Amazon Web Services account or the parameters shared with you when you enable the Shared option. Request results are returned on a best-effort basis. If you specify Max...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeParametersResult> describeParameters(
    $0.DescribeParametersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeParameters, request, options: options);
  }

  /// Lists the patch baselines in your Amazon Web Services account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribePatchBaselinesResult> describePatchBaselines(
    $0.DescribePatchBaselinesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describePatchBaselines, request,
        options: options);
  }

  /// Lists all patch groups that have been registered with patch baselines.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribePatchGroupsResult> describePatchGroups(
    $0.DescribePatchGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describePatchGroups, request, options: options);
  }

  /// Returns high-level aggregated patch compliance state information for a patch group.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribePatchGroupStateResult>
      describePatchGroupState(
    $0.DescribePatchGroupStateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describePatchGroupState, request,
        options: options);
  }

  /// Lists the properties of available patches organized by product, product family, classification, severity, and other properties of available patches. You can use the reported properties in the filte...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribePatchPropertiesResult>
      describePatchProperties(
    $0.DescribePatchPropertiesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describePatchProperties, request,
        options: options);
  }

  /// Retrieves a list of all active sessions (both connected and disconnected) or terminated sessions from the past 30 days.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeSessionsResponse> describeSessions(
    $0.DescribeSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeSessions, request, options: options);
  }

  /// Deletes the association between an OpsItem and a related item. For example, this API operation can delete an Incident Manager incident from an OpsItem. Incident Manager is a tool in Amazon Web Serv...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DisassociateOpsItemRelatedItemResponse>
      disassociateOpsItemRelatedItem(
    $0.DisassociateOpsItemRelatedItemRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disassociateOpsItemRelatedItem, request,
        options: options);
  }

  /// Returns a credentials set to be used with just-in-time node access.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetAccessTokenResponse> getAccessToken(
    $0.GetAccessTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccessToken, request, options: options);
  }

  /// Get detailed information about a particular Automation execution.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetAutomationExecutionResult> getAutomationExecution(
    $0.GetAutomationExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAutomationExecution, request,
        options: options);
  }

  /// Gets the state of a Amazon Web Services Systems Manager change calendar at the current time or a specified time. If you specify a time, GetCalendarState returns the state of the calendar at that sp...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCalendarStateResponse> getCalendarState(
    $0.GetCalendarStateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCalendarState, request, options: options);
  }

  /// Returns detailed information about command execution for an invocation or plugin. The Run Command API follows an eventual consistency model, due to the distributed nature of the system supporting t...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCommandInvocationResult> getCommandInvocation(
    $0.GetCommandInvocationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCommandInvocation, request, options: options);
  }

  /// Retrieves the Session Manager connection status for a managed node to determine whether it is running and ready to receive Session Manager connections.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetConnectionStatusResponse> getConnectionStatus(
    $0.GetConnectionStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getConnectionStatus, request, options: options);
  }

  /// Retrieves the default patch baseline. Amazon Web Services Systems Manager supports creating multiple default patch baselines. For example, you can create a default patch baseline for each operating...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDefaultPatchBaselineResult>
      getDefaultPatchBaseline(
    $0.GetDefaultPatchBaselineRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDefaultPatchBaseline, request,
        options: options);
  }

  /// Retrieves the current snapshot for the patch baseline the managed node uses. This API is primarily used by the AWS-RunPatchBaseline Systems Manager document (SSM document). If you run the command l...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDeployablePatchSnapshotForInstanceResult>
      getDeployablePatchSnapshotForInstance(
    $0.GetDeployablePatchSnapshotForInstanceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDeployablePatchSnapshotForInstance, request,
        options: options);
  }

  /// Gets the contents of the specified Amazon Web Services Systems Manager document (SSM document).
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDocumentResult> getDocument(
    $0.GetDocumentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDocument, request, options: options);
  }

  /// Initiates the process of retrieving an existing preview that shows the effects that running a specified Automation runbook would have on the targeted resources.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetExecutionPreviewResponse> getExecutionPreview(
    $0.GetExecutionPreviewRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getExecutionPreview, request, options: options);
  }

  /// Query inventory information. This includes managed node status, such as Stopped or Terminated.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetInventoryResult> getInventory(
    $0.GetInventoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInventory, request, options: options);
  }

  /// Return a list of inventory type names for the account, or return a list of attribute names for a specific Inventory item type.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetInventorySchemaResult> getInventorySchema(
    $0.GetInventorySchemaRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInventorySchema, request, options: options);
  }

  /// Retrieves a maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetMaintenanceWindowResult> getMaintenanceWindow(
    $0.GetMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMaintenanceWindow, request, options: options);
  }

  /// Retrieves details about a specific a maintenance window execution.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetMaintenanceWindowExecutionResult>
      getMaintenanceWindowExecution(
    $0.GetMaintenanceWindowExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMaintenanceWindowExecution, request,
        options: options);
  }

  /// Retrieves the details about a specific task run as part of a maintenance window execution.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetMaintenanceWindowExecutionTaskResult>
      getMaintenanceWindowExecutionTask(
    $0.GetMaintenanceWindowExecutionTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMaintenanceWindowExecutionTask, request,
        options: options);
  }

  /// Retrieves information about a specific task running on a specific target.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetMaintenanceWindowExecutionTaskInvocationResult>
      getMaintenanceWindowExecutionTaskInvocation(
    $0.GetMaintenanceWindowExecutionTaskInvocationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$getMaintenanceWindowExecutionTaskInvocation, request,
        options: options);
  }

  /// Retrieves the details of a maintenance window task. For maintenance window tasks without a specified target, you can't supply values for --max-errors and --max-concurrency. Instead, the system inse...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetMaintenanceWindowTaskResult>
      getMaintenanceWindowTask(
    $0.GetMaintenanceWindowTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMaintenanceWindowTask, request,
        options: options);
  }

  /// Get information about an OpsItem by using the ID. You must have permission in Identity and Access Management (IAM) to view information about an OpsItem. For more information, see Set up OpsCenter i...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetOpsItemResponse> getOpsItem(
    $0.GetOpsItemRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOpsItem, request, options: options);
  }

  /// View operational metadata related to an application in Application Manager.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetOpsMetadataResult> getOpsMetadata(
    $0.GetOpsMetadataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOpsMetadata, request, options: options);
  }

  /// View a summary of operations metadata (OpsData) based on specified filters and aggregators. OpsData can include information about Amazon Web Services Systems Manager OpsCenter operational workitems...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetOpsSummaryResult> getOpsSummary(
    $0.GetOpsSummaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOpsSummary, request, options: options);
  }

  /// Get information about a single parameter by specifying the parameter name. Parameter names can't contain spaces. The service removes any spaces specified for the beginning or end of a parameter nam...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetParameterResult> getParameter(
    $0.GetParameterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getParameter, request, options: options);
  }

  /// Retrieves the history of all changes to a parameter. Parameter names can't contain spaces. The service removes any spaces specified for the beginning or end of a parameter name. If the specified na...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetParameterHistoryResult> getParameterHistory(
    $0.GetParameterHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getParameterHistory, request, options: options);
  }

  /// Get information about one or more parameters by specifying multiple parameter names. To get information about a single parameter, you can use the GetParameter operation instead. Parameter names can...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetParametersResult> getParameters(
    $0.GetParametersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getParameters, request, options: options);
  }

  /// Retrieve information about one or more parameters under a specified level in a hierarchy. Request results are returned on a best-effort basis. If you specify MaxResults in the request, the response...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetParametersByPathResult> getParametersByPath(
    $0.GetParametersByPathRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getParametersByPath, request, options: options);
  }

  /// Retrieves information about a patch baseline.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetPatchBaselineResult> getPatchBaseline(
    $0.GetPatchBaselineRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPatchBaseline, request, options: options);
  }

  /// Retrieves the patch baseline that should be used for the specified patch group.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetPatchBaselineForPatchGroupResult>
      getPatchBaselineForPatchGroup(
    $0.GetPatchBaselineForPatchGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPatchBaselineForPatchGroup, request,
        options: options);
  }

  /// Returns an array of the Policy object.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetResourcePoliciesResponse> getResourcePolicies(
    $0.GetResourcePoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getResourcePolicies, request, options: options);
  }

  /// ServiceSetting is an account-level setting for an Amazon Web Services service. This setting defines how a user interacts with or uses a service or a feature of a service. For example, if an Amazon ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetServiceSettingResult> getServiceSetting(
    $0.GetServiceSettingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getServiceSetting, request, options: options);
  }

  /// A parameter label is a user-defined alias to help you manage different versions of a parameter. When you modify a parameter, Amazon Web Services Systems Manager automatically saves a new version an...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.LabelParameterVersionResult> labelParameterVersion(
    $0.LabelParameterVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$labelParameterVersion, request, options: options);
  }

  /// Returns all State Manager associations in the current Amazon Web Services account and Amazon Web Services Region. You can limit the results to a specific State Manager association document or manag...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAssociationsResult> listAssociations(
    $0.ListAssociationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAssociations, request, options: options);
  }

  /// Retrieves all versions of an association for a specific association ID.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListAssociationVersionsResult>
      listAssociationVersions(
    $0.ListAssociationVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAssociationVersions, request,
        options: options);
  }

  /// An invocation is copy of a command sent to a specific managed node. A command can apply to one or more managed nodes. A command invocation applies to one managed node. For example, if a user runs S...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListCommandInvocationsResult> listCommandInvocations(
    $0.ListCommandInvocationsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCommandInvocations, request,
        options: options);
  }

  /// Lists the commands requested by users of the Amazon Web Services account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListCommandsResult> listCommands(
    $0.ListCommandsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCommands, request, options: options);
  }

  /// For a specified resource ID, this API operation returns a list of compliance statuses for different resource types. Currently, you can only specify one resource ID per call. List results depend on ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListComplianceItemsResult> listComplianceItems(
    $0.ListComplianceItemsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listComplianceItems, request, options: options);
  }

  /// Returns a summary count of compliant and non-compliant resources for a compliance type. For example, this call can return State Manager associations, patches, or custom compliance types according t...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListComplianceSummariesResult>
      listComplianceSummaries(
    $0.ListComplianceSummariesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listComplianceSummaries, request,
        options: options);
  }

  /// Amazon Web Services Systems Manager Change Manager is no longer open to new customers. Existing customers can continue to use the service as normal. For more information, see Amazon Web Services Sy...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListDocumentMetadataHistoryResponse>
      listDocumentMetadataHistory(
    $0.ListDocumentMetadataHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDocumentMetadataHistory, request,
        options: options);
  }

  /// Returns all Systems Manager (SSM) documents in the current Amazon Web Services account and Amazon Web Services Region. You can limit the results of this request by using a filter.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListDocumentsResult> listDocuments(
    $0.ListDocumentsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDocuments, request, options: options);
  }

  /// List all versions for a document.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListDocumentVersionsResult> listDocumentVersions(
    $0.ListDocumentVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDocumentVersions, request, options: options);
  }

  /// A list of inventory items returned by the request.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListInventoryEntriesResult> listInventoryEntries(
    $0.ListInventoryEntriesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInventoryEntries, request, options: options);
  }

  /// Takes in filters and returns a list of managed nodes matching the filter criteria.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListNodesResult> listNodes(
    $0.ListNodesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listNodes, request, options: options);
  }

  /// Generates a summary of managed instance/node metadata based on the filters and aggregators you specify. Results are grouped by the input aggregator you specify.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListNodesSummaryResult> listNodesSummary(
    $0.ListNodesSummaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listNodesSummary, request, options: options);
  }

  /// Returns a list of all OpsItem events in the current Amazon Web Services Region and Amazon Web Services account. You can limit the results to events associated with specific OpsItems by specifying a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListOpsItemEventsResponse> listOpsItemEvents(
    $0.ListOpsItemEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOpsItemEvents, request, options: options);
  }

  /// Lists all related-item resources associated with a Systems Manager OpsCenter OpsItem. OpsCenter is a tool in Amazon Web Services Systems Manager.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListOpsItemRelatedItemsResponse>
      listOpsItemRelatedItems(
    $0.ListOpsItemRelatedItemsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOpsItemRelatedItems, request,
        options: options);
  }

  /// Amazon Web Services Systems Manager calls this API operation when displaying all Application Manager OpsMetadata objects or blobs.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListOpsMetadataResult> listOpsMetadata(
    $0.ListOpsMetadataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOpsMetadata, request, options: options);
  }

  /// Returns a resource-level summary count. The summary includes information about compliant and non-compliant statuses and detailed compliance-item severity counts, according to the filter criteria yo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListResourceComplianceSummariesResult>
      listResourceComplianceSummaries(
    $0.ListResourceComplianceSummariesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResourceComplianceSummaries, request,
        options: options);
  }

  /// Lists your resource data sync configurations. Includes information about the last time a sync attempted to start, the last sync status, and the last time a sync successfully completed. The number o...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListResourceDataSyncResult> listResourceDataSync(
    $0.ListResourceDataSyncRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResourceDataSync, request, options: options);
  }

  /// Returns a list of the tags assigned to the specified resource. For information about the ID format for each supported resource type, see AddTagsToResource.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceResult> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Shares a Amazon Web Services Systems Manager document (SSM document)publicly or privately. If you share a document privately, you must specify the Amazon Web Services user IDs for those people who ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ModifyDocumentPermissionResponse>
      modifyDocumentPermission(
    $0.ModifyDocumentPermissionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$modifyDocumentPermission, request,
        options: options);
  }

  /// Registers a compliance type and other compliance details on a designated resource. This operation lets you register custom compliance details with a resource. This call overwrites existing complian...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutComplianceItemsResult> putComplianceItems(
    $0.PutComplianceItemsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putComplianceItems, request, options: options);
  }

  /// Bulk update custom inventory items on one or more managed nodes. The request adds an inventory item, if it doesn't already exist, or updates an inventory item, if it does exist.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutInventoryResult> putInventory(
    $0.PutInventoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putInventory, request, options: options);
  }

  /// Create or update a parameter in Parameter Store.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutParameterResult> putParameter(
    $0.PutParameterRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putParameter, request, options: options);
  }

  /// Creates or updates a Systems Manager resource policy. A resource policy helps you to define the IAM entity (for example, an Amazon Web Services account) that can manage your Systems Manager resourc...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.PutResourcePolicyResponse> putResourcePolicy(
    $0.PutResourcePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putResourcePolicy, request, options: options);
  }

  /// Defines the default patch baseline for the relevant operating system. To reset the Amazon Web Services-predefined patch baseline as the default, specify the full patch baseline Amazon Resource Name...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RegisterDefaultPatchBaselineResult>
      registerDefaultPatchBaseline(
    $0.RegisterDefaultPatchBaselineRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerDefaultPatchBaseline, request,
        options: options);
  }

  /// Registers a patch baseline for a patch group.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RegisterPatchBaselineForPatchGroupResult>
      registerPatchBaselineForPatchGroup(
    $0.RegisterPatchBaselineForPatchGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerPatchBaselineForPatchGroup, request,
        options: options);
  }

  /// Registers a target with a maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RegisterTargetWithMaintenanceWindowResult>
      registerTargetWithMaintenanceWindow(
    $0.RegisterTargetWithMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerTargetWithMaintenanceWindow, request,
        options: options);
  }

  /// Adds a new task to a maintenance window.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RegisterTaskWithMaintenanceWindowResult>
      registerTaskWithMaintenanceWindow(
    $0.RegisterTaskWithMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$registerTaskWithMaintenanceWindow, request,
        options: options);
  }

  /// Removes tag keys from the specified resource.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RemoveTagsFromResourceResult> removeTagsFromResource(
    $0.RemoveTagsFromResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeTagsFromResource, request,
        options: options);
  }

  /// ServiceSetting is an account-level setting for an Amazon Web Services service. This setting defines how a user interacts with or uses a service or a feature of a service. For example, if an Amazon ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ResetServiceSettingResult> resetServiceSetting(
    $0.ResetServiceSettingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resetServiceSetting, request, options: options);
  }

  /// Reconnects a session to a managed node after it has been disconnected. Connections can be resumed for disconnected sessions, but not terminated sessions. This command is primarily for use by client...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ResumeSessionResponse> resumeSession(
    $0.ResumeSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resumeSession, request, options: options);
  }

  /// Sends a signal to an Automation execution to change the current behavior or status of the execution.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SendAutomationSignalResult> sendAutomationSignal(
    $0.SendAutomationSignalRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendAutomationSignal, request, options: options);
  }

  /// Runs commands on one or more managed nodes.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SendCommandResult> sendCommand(
    $0.SendCommandRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendCommand, request, options: options);
  }

  /// Starts the workflow for just-in-time node access sessions.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartAccessRequestResponse> startAccessRequest(
    $0.StartAccessRequestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startAccessRequest, request, options: options);
  }

  /// Runs an association immediately and only one time. This operation can be helpful when troubleshooting associations.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartAssociationsOnceResult> startAssociationsOnce(
    $0.StartAssociationsOnceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startAssociationsOnce, request, options: options);
  }

  /// Initiates execution of an Automation runbook.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartAutomationExecutionResult>
      startAutomationExecution(
    $0.StartAutomationExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startAutomationExecution, request,
        options: options);
  }

  /// Amazon Web Services Systems Manager Change Manager is no longer open to new customers. Existing customers can continue to use the service as normal. For more information, see Amazon Web Services Sy...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartChangeRequestExecutionResult>
      startChangeRequestExecution(
    $0.StartChangeRequestExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startChangeRequestExecution, request,
        options: options);
  }

  /// Initiates the process of creating a preview showing the effects that running a specified Automation runbook would have on the targeted resources.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartExecutionPreviewResponse> startExecutionPreview(
    $0.StartExecutionPreviewRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startExecutionPreview, request, options: options);
  }

  /// Initiates a connection to a target (for example, a managed node) for a Session Manager session. Returns a URL and token that can be used to open a WebSocket connection for sending input and receivi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartSessionResponse> startSession(
    $0.StartSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startSession, request, options: options);
  }

  /// Stop an Automation that is currently running.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopAutomationExecutionResult>
      stopAutomationExecution(
    $0.StopAutomationExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopAutomationExecution, request,
        options: options);
  }

  /// Permanently ends a session and closes the data connection between the Session Manager client and SSM Agent on the managed node. A terminated session can't be resumed.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TerminateSessionResponse> terminateSession(
    $0.TerminateSessionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$terminateSession, request, options: options);
  }

  /// Remove a label or labels from a parameter. Parameter names can't contain spaces. The service removes any spaces specified for the beginning or end of a parameter name. If the specified name for a p...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UnlabelParameterVersionResult>
      unlabelParameterVersion(
    $0.UnlabelParameterVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$unlabelParameterVersion, request,
        options: options);
  }

  /// Updates an association. You can update the association name and version, the document version, schedule, parameters, and Amazon Simple Storage Service (Amazon S3) output. When you call UpdateAssoci...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateAssociationResult> updateAssociation(
    $0.UpdateAssociationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAssociation, request, options: options);
  }

  /// Updates the status of the Amazon Web Services Systems Manager document (SSM document) associated with the specified managed node. UpdateAssociationStatus is primarily used by the Amazon Web Service...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateAssociationStatusResult>
      updateAssociationStatus(
    $0.UpdateAssociationStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAssociationStatus, request,
        options: options);
  }

  /// Updates one or more values for an SSM document.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateDocumentResult> updateDocument(
    $0.UpdateDocumentRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDocument, request, options: options);
  }

  /// Set the default version of a document. If you change a document version for a State Manager association, Systems Manager immediately runs the association unless you previously specifed the apply-on...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateDocumentDefaultVersionResult>
      updateDocumentDefaultVersion(
    $0.UpdateDocumentDefaultVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDocumentDefaultVersion, request,
        options: options);
  }

  /// Amazon Web Services Systems Manager Change Manager is no longer open to new customers. Existing customers can continue to use the service as normal. For more information, see Amazon Web Services Sy...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateDocumentMetadataResponse>
      updateDocumentMetadata(
    $0.UpdateDocumentMetadataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDocumentMetadata, request,
        options: options);
  }

  /// Updates an existing maintenance window. Only specified parameters are modified. The value you specify for Duration determines the specific end time for the maintenance window based on the time it b...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateMaintenanceWindowResult>
      updateMaintenanceWindow(
    $0.UpdateMaintenanceWindowRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMaintenanceWindow, request,
        options: options);
  }

  /// Modifies the target of an existing maintenance window. You can change the following: Name Description Owner IDs for an ID target Tags for a Tag target From any supported tag type to another. The th...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateMaintenanceWindowTargetResult>
      updateMaintenanceWindowTarget(
    $0.UpdateMaintenanceWindowTargetRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMaintenanceWindowTarget, request,
        options: options);
  }

  /// Modifies a task assigned to a maintenance window. You can't change the task type, but you can change the following values: TaskARN. For example, you can change a RUN_COMMAND task from AWS-RunPowerS...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateMaintenanceWindowTaskResult>
      updateMaintenanceWindowTask(
    $0.UpdateMaintenanceWindowTaskRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateMaintenanceWindowTask, request,
        options: options);
  }

  /// Changes the Identity and Access Management (IAM) role that is assigned to the on-premises server, edge device, or virtual machines (VM). IAM roles are first assigned to these hybrid nodes during th...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateManagedInstanceRoleResult>
      updateManagedInstanceRole(
    $0.UpdateManagedInstanceRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateManagedInstanceRole, request,
        options: options);
  }

  /// Edit or change an OpsItem. You must have permission in Identity and Access Management (IAM) to update an OpsItem. For more information, see Set up OpsCenter in the Amazon Web Services Systems Manag...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateOpsItemResponse> updateOpsItem(
    $0.UpdateOpsItemRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateOpsItem, request, options: options);
  }

  /// Amazon Web Services Systems Manager calls this API operation when you edit OpsMetadata in Application Manager.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateOpsMetadataResult> updateOpsMetadata(
    $0.UpdateOpsMetadataRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateOpsMetadata, request, options: options);
  }

  /// Modifies an existing patch baseline. Fields not specified in the request are left unchanged. For information about valid key-value pairs in PatchFilters for each supported operating system type, se...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdatePatchBaselineResult> updatePatchBaseline(
    $0.UpdatePatchBaselineRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updatePatchBaseline, request, options: options);
  }

  /// Update a resource data sync. After you create a resource data sync for a Region, you can't change the account options for that sync. For example, if you create a sync in the us-east-2 (Ohio) Region...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateResourceDataSyncResult> updateResourceDataSync(
    $0.UpdateResourceDataSyncRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateResourceDataSync, request,
        options: options);
  }

  /// ServiceSetting is an account-level setting for an Amazon Web Services service. This setting defines how a user interacts with or uses a service or a feature of a service. For example, if an Amazon ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateServiceSettingResult> updateServiceSetting(
    $0.UpdateServiceSettingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateServiceSetting, request, options: options);
  }

  // method descriptors

  static final _$addTagsToResource = $grpc.ClientMethod<
          $0.AddTagsToResourceRequest, $0.AddTagsToResourceResult>(
      '/ssm.SSMService/AddTagsToResource',
      ($0.AddTagsToResourceRequest value) => value.writeToBuffer(),
      $0.AddTagsToResourceResult.fromBuffer);
  static final _$associateOpsItemRelatedItem = $grpc.ClientMethod<
          $0.AssociateOpsItemRelatedItemRequest,
          $0.AssociateOpsItemRelatedItemResponse>(
      '/ssm.SSMService/AssociateOpsItemRelatedItem',
      ($0.AssociateOpsItemRelatedItemRequest value) => value.writeToBuffer(),
      $0.AssociateOpsItemRelatedItemResponse.fromBuffer);
  static final _$cancelCommand =
      $grpc.ClientMethod<$0.CancelCommandRequest, $0.CancelCommandResult>(
          '/ssm.SSMService/CancelCommand',
          ($0.CancelCommandRequest value) => value.writeToBuffer(),
          $0.CancelCommandResult.fromBuffer);
  static final _$cancelMaintenanceWindowExecution = $grpc.ClientMethod<
          $0.CancelMaintenanceWindowExecutionRequest,
          $0.CancelMaintenanceWindowExecutionResult>(
      '/ssm.SSMService/CancelMaintenanceWindowExecution',
      ($0.CancelMaintenanceWindowExecutionRequest value) =>
          value.writeToBuffer(),
      $0.CancelMaintenanceWindowExecutionResult.fromBuffer);
  static final _$createActivation =
      $grpc.ClientMethod<$0.CreateActivationRequest, $0.CreateActivationResult>(
          '/ssm.SSMService/CreateActivation',
          ($0.CreateActivationRequest value) => value.writeToBuffer(),
          $0.CreateActivationResult.fromBuffer);
  static final _$createAssociation = $grpc.ClientMethod<
          $0.CreateAssociationRequest, $0.CreateAssociationResult>(
      '/ssm.SSMService/CreateAssociation',
      ($0.CreateAssociationRequest value) => value.writeToBuffer(),
      $0.CreateAssociationResult.fromBuffer);
  static final _$createAssociationBatch = $grpc.ClientMethod<
          $0.CreateAssociationBatchRequest, $0.CreateAssociationBatchResult>(
      '/ssm.SSMService/CreateAssociationBatch',
      ($0.CreateAssociationBatchRequest value) => value.writeToBuffer(),
      $0.CreateAssociationBatchResult.fromBuffer);
  static final _$createDocument =
      $grpc.ClientMethod<$0.CreateDocumentRequest, $0.CreateDocumentResult>(
          '/ssm.SSMService/CreateDocument',
          ($0.CreateDocumentRequest value) => value.writeToBuffer(),
          $0.CreateDocumentResult.fromBuffer);
  static final _$createMaintenanceWindow = $grpc.ClientMethod<
          $0.CreateMaintenanceWindowRequest, $0.CreateMaintenanceWindowResult>(
      '/ssm.SSMService/CreateMaintenanceWindow',
      ($0.CreateMaintenanceWindowRequest value) => value.writeToBuffer(),
      $0.CreateMaintenanceWindowResult.fromBuffer);
  static final _$createOpsItem =
      $grpc.ClientMethod<$0.CreateOpsItemRequest, $0.CreateOpsItemResponse>(
          '/ssm.SSMService/CreateOpsItem',
          ($0.CreateOpsItemRequest value) => value.writeToBuffer(),
          $0.CreateOpsItemResponse.fromBuffer);
  static final _$createOpsMetadata = $grpc.ClientMethod<
          $0.CreateOpsMetadataRequest, $0.CreateOpsMetadataResult>(
      '/ssm.SSMService/CreateOpsMetadata',
      ($0.CreateOpsMetadataRequest value) => value.writeToBuffer(),
      $0.CreateOpsMetadataResult.fromBuffer);
  static final _$createPatchBaseline = $grpc.ClientMethod<
          $0.CreatePatchBaselineRequest, $0.CreatePatchBaselineResult>(
      '/ssm.SSMService/CreatePatchBaseline',
      ($0.CreatePatchBaselineRequest value) => value.writeToBuffer(),
      $0.CreatePatchBaselineResult.fromBuffer);
  static final _$createResourceDataSync = $grpc.ClientMethod<
          $0.CreateResourceDataSyncRequest, $0.CreateResourceDataSyncResult>(
      '/ssm.SSMService/CreateResourceDataSync',
      ($0.CreateResourceDataSyncRequest value) => value.writeToBuffer(),
      $0.CreateResourceDataSyncResult.fromBuffer);
  static final _$deleteActivation =
      $grpc.ClientMethod<$0.DeleteActivationRequest, $0.DeleteActivationResult>(
          '/ssm.SSMService/DeleteActivation',
          ($0.DeleteActivationRequest value) => value.writeToBuffer(),
          $0.DeleteActivationResult.fromBuffer);
  static final _$deleteAssociation = $grpc.ClientMethod<
          $0.DeleteAssociationRequest, $0.DeleteAssociationResult>(
      '/ssm.SSMService/DeleteAssociation',
      ($0.DeleteAssociationRequest value) => value.writeToBuffer(),
      $0.DeleteAssociationResult.fromBuffer);
  static final _$deleteDocument =
      $grpc.ClientMethod<$0.DeleteDocumentRequest, $0.DeleteDocumentResult>(
          '/ssm.SSMService/DeleteDocument',
          ($0.DeleteDocumentRequest value) => value.writeToBuffer(),
          $0.DeleteDocumentResult.fromBuffer);
  static final _$deleteInventory =
      $grpc.ClientMethod<$0.DeleteInventoryRequest, $0.DeleteInventoryResult>(
          '/ssm.SSMService/DeleteInventory',
          ($0.DeleteInventoryRequest value) => value.writeToBuffer(),
          $0.DeleteInventoryResult.fromBuffer);
  static final _$deleteMaintenanceWindow = $grpc.ClientMethod<
          $0.DeleteMaintenanceWindowRequest, $0.DeleteMaintenanceWindowResult>(
      '/ssm.SSMService/DeleteMaintenanceWindow',
      ($0.DeleteMaintenanceWindowRequest value) => value.writeToBuffer(),
      $0.DeleteMaintenanceWindowResult.fromBuffer);
  static final _$deleteOpsItem =
      $grpc.ClientMethod<$0.DeleteOpsItemRequest, $0.DeleteOpsItemResponse>(
          '/ssm.SSMService/DeleteOpsItem',
          ($0.DeleteOpsItemRequest value) => value.writeToBuffer(),
          $0.DeleteOpsItemResponse.fromBuffer);
  static final _$deleteOpsMetadata = $grpc.ClientMethod<
          $0.DeleteOpsMetadataRequest, $0.DeleteOpsMetadataResult>(
      '/ssm.SSMService/DeleteOpsMetadata',
      ($0.DeleteOpsMetadataRequest value) => value.writeToBuffer(),
      $0.DeleteOpsMetadataResult.fromBuffer);
  static final _$deleteParameter =
      $grpc.ClientMethod<$0.DeleteParameterRequest, $0.DeleteParameterResult>(
          '/ssm.SSMService/DeleteParameter',
          ($0.DeleteParameterRequest value) => value.writeToBuffer(),
          $0.DeleteParameterResult.fromBuffer);
  static final _$deleteParameters =
      $grpc.ClientMethod<$0.DeleteParametersRequest, $0.DeleteParametersResult>(
          '/ssm.SSMService/DeleteParameters',
          ($0.DeleteParametersRequest value) => value.writeToBuffer(),
          $0.DeleteParametersResult.fromBuffer);
  static final _$deletePatchBaseline = $grpc.ClientMethod<
          $0.DeletePatchBaselineRequest, $0.DeletePatchBaselineResult>(
      '/ssm.SSMService/DeletePatchBaseline',
      ($0.DeletePatchBaselineRequest value) => value.writeToBuffer(),
      $0.DeletePatchBaselineResult.fromBuffer);
  static final _$deleteResourceDataSync = $grpc.ClientMethod<
          $0.DeleteResourceDataSyncRequest, $0.DeleteResourceDataSyncResult>(
      '/ssm.SSMService/DeleteResourceDataSync',
      ($0.DeleteResourceDataSyncRequest value) => value.writeToBuffer(),
      $0.DeleteResourceDataSyncResult.fromBuffer);
  static final _$deleteResourcePolicy = $grpc.ClientMethod<
          $0.DeleteResourcePolicyRequest, $0.DeleteResourcePolicyResponse>(
      '/ssm.SSMService/DeleteResourcePolicy',
      ($0.DeleteResourcePolicyRequest value) => value.writeToBuffer(),
      $0.DeleteResourcePolicyResponse.fromBuffer);
  static final _$deregisterManagedInstance = $grpc.ClientMethod<
          $0.DeregisterManagedInstanceRequest,
          $0.DeregisterManagedInstanceResult>(
      '/ssm.SSMService/DeregisterManagedInstance',
      ($0.DeregisterManagedInstanceRequest value) => value.writeToBuffer(),
      $0.DeregisterManagedInstanceResult.fromBuffer);
  static final _$deregisterPatchBaselineForPatchGroup = $grpc.ClientMethod<
          $0.DeregisterPatchBaselineForPatchGroupRequest,
          $0.DeregisterPatchBaselineForPatchGroupResult>(
      '/ssm.SSMService/DeregisterPatchBaselineForPatchGroup',
      ($0.DeregisterPatchBaselineForPatchGroupRequest value) =>
          value.writeToBuffer(),
      $0.DeregisterPatchBaselineForPatchGroupResult.fromBuffer);
  static final _$deregisterTargetFromMaintenanceWindow = $grpc.ClientMethod<
          $0.DeregisterTargetFromMaintenanceWindowRequest,
          $0.DeregisterTargetFromMaintenanceWindowResult>(
      '/ssm.SSMService/DeregisterTargetFromMaintenanceWindow',
      ($0.DeregisterTargetFromMaintenanceWindowRequest value) =>
          value.writeToBuffer(),
      $0.DeregisterTargetFromMaintenanceWindowResult.fromBuffer);
  static final _$deregisterTaskFromMaintenanceWindow = $grpc.ClientMethod<
          $0.DeregisterTaskFromMaintenanceWindowRequest,
          $0.DeregisterTaskFromMaintenanceWindowResult>(
      '/ssm.SSMService/DeregisterTaskFromMaintenanceWindow',
      ($0.DeregisterTaskFromMaintenanceWindowRequest value) =>
          value.writeToBuffer(),
      $0.DeregisterTaskFromMaintenanceWindowResult.fromBuffer);
  static final _$describeActivations = $grpc.ClientMethod<
          $0.DescribeActivationsRequest, $0.DescribeActivationsResult>(
      '/ssm.SSMService/DescribeActivations',
      ($0.DescribeActivationsRequest value) => value.writeToBuffer(),
      $0.DescribeActivationsResult.fromBuffer);
  static final _$describeAssociation = $grpc.ClientMethod<
          $0.DescribeAssociationRequest, $0.DescribeAssociationResult>(
      '/ssm.SSMService/DescribeAssociation',
      ($0.DescribeAssociationRequest value) => value.writeToBuffer(),
      $0.DescribeAssociationResult.fromBuffer);
  static final _$describeAssociationExecutions = $grpc.ClientMethod<
          $0.DescribeAssociationExecutionsRequest,
          $0.DescribeAssociationExecutionsResult>(
      '/ssm.SSMService/DescribeAssociationExecutions',
      ($0.DescribeAssociationExecutionsRequest value) => value.writeToBuffer(),
      $0.DescribeAssociationExecutionsResult.fromBuffer);
  static final _$describeAssociationExecutionTargets = $grpc.ClientMethod<
          $0.DescribeAssociationExecutionTargetsRequest,
          $0.DescribeAssociationExecutionTargetsResult>(
      '/ssm.SSMService/DescribeAssociationExecutionTargets',
      ($0.DescribeAssociationExecutionTargetsRequest value) =>
          value.writeToBuffer(),
      $0.DescribeAssociationExecutionTargetsResult.fromBuffer);
  static final _$describeAutomationExecutions = $grpc.ClientMethod<
          $0.DescribeAutomationExecutionsRequest,
          $0.DescribeAutomationExecutionsResult>(
      '/ssm.SSMService/DescribeAutomationExecutions',
      ($0.DescribeAutomationExecutionsRequest value) => value.writeToBuffer(),
      $0.DescribeAutomationExecutionsResult.fromBuffer);
  static final _$describeAutomationStepExecutions = $grpc.ClientMethod<
          $0.DescribeAutomationStepExecutionsRequest,
          $0.DescribeAutomationStepExecutionsResult>(
      '/ssm.SSMService/DescribeAutomationStepExecutions',
      ($0.DescribeAutomationStepExecutionsRequest value) =>
          value.writeToBuffer(),
      $0.DescribeAutomationStepExecutionsResult.fromBuffer);
  static final _$describeAvailablePatches = $grpc.ClientMethod<
          $0.DescribeAvailablePatchesRequest,
          $0.DescribeAvailablePatchesResult>(
      '/ssm.SSMService/DescribeAvailablePatches',
      ($0.DescribeAvailablePatchesRequest value) => value.writeToBuffer(),
      $0.DescribeAvailablePatchesResult.fromBuffer);
  static final _$describeDocument =
      $grpc.ClientMethod<$0.DescribeDocumentRequest, $0.DescribeDocumentResult>(
          '/ssm.SSMService/DescribeDocument',
          ($0.DescribeDocumentRequest value) => value.writeToBuffer(),
          $0.DescribeDocumentResult.fromBuffer);
  static final _$describeDocumentPermission = $grpc.ClientMethod<
          $0.DescribeDocumentPermissionRequest,
          $0.DescribeDocumentPermissionResponse>(
      '/ssm.SSMService/DescribeDocumentPermission',
      ($0.DescribeDocumentPermissionRequest value) => value.writeToBuffer(),
      $0.DescribeDocumentPermissionResponse.fromBuffer);
  static final _$describeEffectiveInstanceAssociations = $grpc.ClientMethod<
          $0.DescribeEffectiveInstanceAssociationsRequest,
          $0.DescribeEffectiveInstanceAssociationsResult>(
      '/ssm.SSMService/DescribeEffectiveInstanceAssociations',
      ($0.DescribeEffectiveInstanceAssociationsRequest value) =>
          value.writeToBuffer(),
      $0.DescribeEffectiveInstanceAssociationsResult.fromBuffer);
  static final _$describeEffectivePatchesForPatchBaseline = $grpc.ClientMethod<
          $0.DescribeEffectivePatchesForPatchBaselineRequest,
          $0.DescribeEffectivePatchesForPatchBaselineResult>(
      '/ssm.SSMService/DescribeEffectivePatchesForPatchBaseline',
      ($0.DescribeEffectivePatchesForPatchBaselineRequest value) =>
          value.writeToBuffer(),
      $0.DescribeEffectivePatchesForPatchBaselineResult.fromBuffer);
  static final _$describeInstanceAssociationsStatus = $grpc.ClientMethod<
          $0.DescribeInstanceAssociationsStatusRequest,
          $0.DescribeInstanceAssociationsStatusResult>(
      '/ssm.SSMService/DescribeInstanceAssociationsStatus',
      ($0.DescribeInstanceAssociationsStatusRequest value) =>
          value.writeToBuffer(),
      $0.DescribeInstanceAssociationsStatusResult.fromBuffer);
  static final _$describeInstanceInformation = $grpc.ClientMethod<
          $0.DescribeInstanceInformationRequest,
          $0.DescribeInstanceInformationResult>(
      '/ssm.SSMService/DescribeInstanceInformation',
      ($0.DescribeInstanceInformationRequest value) => value.writeToBuffer(),
      $0.DescribeInstanceInformationResult.fromBuffer);
  static final _$describeInstancePatches = $grpc.ClientMethod<
          $0.DescribeInstancePatchesRequest, $0.DescribeInstancePatchesResult>(
      '/ssm.SSMService/DescribeInstancePatches',
      ($0.DescribeInstancePatchesRequest value) => value.writeToBuffer(),
      $0.DescribeInstancePatchesResult.fromBuffer);
  static final _$describeInstancePatchStates = $grpc.ClientMethod<
          $0.DescribeInstancePatchStatesRequest,
          $0.DescribeInstancePatchStatesResult>(
      '/ssm.SSMService/DescribeInstancePatchStates',
      ($0.DescribeInstancePatchStatesRequest value) => value.writeToBuffer(),
      $0.DescribeInstancePatchStatesResult.fromBuffer);
  static final _$describeInstancePatchStatesForPatchGroup = $grpc.ClientMethod<
          $0.DescribeInstancePatchStatesForPatchGroupRequest,
          $0.DescribeInstancePatchStatesForPatchGroupResult>(
      '/ssm.SSMService/DescribeInstancePatchStatesForPatchGroup',
      ($0.DescribeInstancePatchStatesForPatchGroupRequest value) =>
          value.writeToBuffer(),
      $0.DescribeInstancePatchStatesForPatchGroupResult.fromBuffer);
  static final _$describeInstanceProperties = $grpc.ClientMethod<
          $0.DescribeInstancePropertiesRequest,
          $0.DescribeInstancePropertiesResult>(
      '/ssm.SSMService/DescribeInstanceProperties',
      ($0.DescribeInstancePropertiesRequest value) => value.writeToBuffer(),
      $0.DescribeInstancePropertiesResult.fromBuffer);
  static final _$describeInventoryDeletions = $grpc.ClientMethod<
          $0.DescribeInventoryDeletionsRequest,
          $0.DescribeInventoryDeletionsResult>(
      '/ssm.SSMService/DescribeInventoryDeletions',
      ($0.DescribeInventoryDeletionsRequest value) => value.writeToBuffer(),
      $0.DescribeInventoryDeletionsResult.fromBuffer);
  static final _$describeMaintenanceWindowExecutions = $grpc.ClientMethod<
          $0.DescribeMaintenanceWindowExecutionsRequest,
          $0.DescribeMaintenanceWindowExecutionsResult>(
      '/ssm.SSMService/DescribeMaintenanceWindowExecutions',
      ($0.DescribeMaintenanceWindowExecutionsRequest value) =>
          value.writeToBuffer(),
      $0.DescribeMaintenanceWindowExecutionsResult.fromBuffer);
  static final _$describeMaintenanceWindowExecutionTaskInvocations =
      $grpc.ClientMethod<
              $0.DescribeMaintenanceWindowExecutionTaskInvocationsRequest,
              $0.DescribeMaintenanceWindowExecutionTaskInvocationsResult>(
          '/ssm.SSMService/DescribeMaintenanceWindowExecutionTaskInvocations',
          ($0.DescribeMaintenanceWindowExecutionTaskInvocationsRequest value) =>
              value.writeToBuffer(),
          $0.DescribeMaintenanceWindowExecutionTaskInvocationsResult
              .fromBuffer);
  static final _$describeMaintenanceWindowExecutionTasks = $grpc.ClientMethod<
          $0.DescribeMaintenanceWindowExecutionTasksRequest,
          $0.DescribeMaintenanceWindowExecutionTasksResult>(
      '/ssm.SSMService/DescribeMaintenanceWindowExecutionTasks',
      ($0.DescribeMaintenanceWindowExecutionTasksRequest value) =>
          value.writeToBuffer(),
      $0.DescribeMaintenanceWindowExecutionTasksResult.fromBuffer);
  static final _$describeMaintenanceWindows = $grpc.ClientMethod<
          $0.DescribeMaintenanceWindowsRequest,
          $0.DescribeMaintenanceWindowsResult>(
      '/ssm.SSMService/DescribeMaintenanceWindows',
      ($0.DescribeMaintenanceWindowsRequest value) => value.writeToBuffer(),
      $0.DescribeMaintenanceWindowsResult.fromBuffer);
  static final _$describeMaintenanceWindowSchedule = $grpc.ClientMethod<
          $0.DescribeMaintenanceWindowScheduleRequest,
          $0.DescribeMaintenanceWindowScheduleResult>(
      '/ssm.SSMService/DescribeMaintenanceWindowSchedule',
      ($0.DescribeMaintenanceWindowScheduleRequest value) =>
          value.writeToBuffer(),
      $0.DescribeMaintenanceWindowScheduleResult.fromBuffer);
  static final _$describeMaintenanceWindowsForTarget = $grpc.ClientMethod<
          $0.DescribeMaintenanceWindowsForTargetRequest,
          $0.DescribeMaintenanceWindowsForTargetResult>(
      '/ssm.SSMService/DescribeMaintenanceWindowsForTarget',
      ($0.DescribeMaintenanceWindowsForTargetRequest value) =>
          value.writeToBuffer(),
      $0.DescribeMaintenanceWindowsForTargetResult.fromBuffer);
  static final _$describeMaintenanceWindowTargets = $grpc.ClientMethod<
          $0.DescribeMaintenanceWindowTargetsRequest,
          $0.DescribeMaintenanceWindowTargetsResult>(
      '/ssm.SSMService/DescribeMaintenanceWindowTargets',
      ($0.DescribeMaintenanceWindowTargetsRequest value) =>
          value.writeToBuffer(),
      $0.DescribeMaintenanceWindowTargetsResult.fromBuffer);
  static final _$describeMaintenanceWindowTasks = $grpc.ClientMethod<
          $0.DescribeMaintenanceWindowTasksRequest,
          $0.DescribeMaintenanceWindowTasksResult>(
      '/ssm.SSMService/DescribeMaintenanceWindowTasks',
      ($0.DescribeMaintenanceWindowTasksRequest value) => value.writeToBuffer(),
      $0.DescribeMaintenanceWindowTasksResult.fromBuffer);
  static final _$describeOpsItems = $grpc.ClientMethod<
          $0.DescribeOpsItemsRequest, $0.DescribeOpsItemsResponse>(
      '/ssm.SSMService/DescribeOpsItems',
      ($0.DescribeOpsItemsRequest value) => value.writeToBuffer(),
      $0.DescribeOpsItemsResponse.fromBuffer);
  static final _$describeParameters = $grpc.ClientMethod<
          $0.DescribeParametersRequest, $0.DescribeParametersResult>(
      '/ssm.SSMService/DescribeParameters',
      ($0.DescribeParametersRequest value) => value.writeToBuffer(),
      $0.DescribeParametersResult.fromBuffer);
  static final _$describePatchBaselines = $grpc.ClientMethod<
          $0.DescribePatchBaselinesRequest, $0.DescribePatchBaselinesResult>(
      '/ssm.SSMService/DescribePatchBaselines',
      ($0.DescribePatchBaselinesRequest value) => value.writeToBuffer(),
      $0.DescribePatchBaselinesResult.fromBuffer);
  static final _$describePatchGroups = $grpc.ClientMethod<
          $0.DescribePatchGroupsRequest, $0.DescribePatchGroupsResult>(
      '/ssm.SSMService/DescribePatchGroups',
      ($0.DescribePatchGroupsRequest value) => value.writeToBuffer(),
      $0.DescribePatchGroupsResult.fromBuffer);
  static final _$describePatchGroupState = $grpc.ClientMethod<
          $0.DescribePatchGroupStateRequest, $0.DescribePatchGroupStateResult>(
      '/ssm.SSMService/DescribePatchGroupState',
      ($0.DescribePatchGroupStateRequest value) => value.writeToBuffer(),
      $0.DescribePatchGroupStateResult.fromBuffer);
  static final _$describePatchProperties = $grpc.ClientMethod<
          $0.DescribePatchPropertiesRequest, $0.DescribePatchPropertiesResult>(
      '/ssm.SSMService/DescribePatchProperties',
      ($0.DescribePatchPropertiesRequest value) => value.writeToBuffer(),
      $0.DescribePatchPropertiesResult.fromBuffer);
  static final _$describeSessions = $grpc.ClientMethod<
          $0.DescribeSessionsRequest, $0.DescribeSessionsResponse>(
      '/ssm.SSMService/DescribeSessions',
      ($0.DescribeSessionsRequest value) => value.writeToBuffer(),
      $0.DescribeSessionsResponse.fromBuffer);
  static final _$disassociateOpsItemRelatedItem = $grpc.ClientMethod<
          $0.DisassociateOpsItemRelatedItemRequest,
          $0.DisassociateOpsItemRelatedItemResponse>(
      '/ssm.SSMService/DisassociateOpsItemRelatedItem',
      ($0.DisassociateOpsItemRelatedItemRequest value) => value.writeToBuffer(),
      $0.DisassociateOpsItemRelatedItemResponse.fromBuffer);
  static final _$getAccessToken =
      $grpc.ClientMethod<$0.GetAccessTokenRequest, $0.GetAccessTokenResponse>(
          '/ssm.SSMService/GetAccessToken',
          ($0.GetAccessTokenRequest value) => value.writeToBuffer(),
          $0.GetAccessTokenResponse.fromBuffer);
  static final _$getAutomationExecution = $grpc.ClientMethod<
          $0.GetAutomationExecutionRequest, $0.GetAutomationExecutionResult>(
      '/ssm.SSMService/GetAutomationExecution',
      ($0.GetAutomationExecutionRequest value) => value.writeToBuffer(),
      $0.GetAutomationExecutionResult.fromBuffer);
  static final _$getCalendarState = $grpc.ClientMethod<
          $0.GetCalendarStateRequest, $0.GetCalendarStateResponse>(
      '/ssm.SSMService/GetCalendarState',
      ($0.GetCalendarStateRequest value) => value.writeToBuffer(),
      $0.GetCalendarStateResponse.fromBuffer);
  static final _$getCommandInvocation = $grpc.ClientMethod<
          $0.GetCommandInvocationRequest, $0.GetCommandInvocationResult>(
      '/ssm.SSMService/GetCommandInvocation',
      ($0.GetCommandInvocationRequest value) => value.writeToBuffer(),
      $0.GetCommandInvocationResult.fromBuffer);
  static final _$getConnectionStatus = $grpc.ClientMethod<
          $0.GetConnectionStatusRequest, $0.GetConnectionStatusResponse>(
      '/ssm.SSMService/GetConnectionStatus',
      ($0.GetConnectionStatusRequest value) => value.writeToBuffer(),
      $0.GetConnectionStatusResponse.fromBuffer);
  static final _$getDefaultPatchBaseline = $grpc.ClientMethod<
          $0.GetDefaultPatchBaselineRequest, $0.GetDefaultPatchBaselineResult>(
      '/ssm.SSMService/GetDefaultPatchBaseline',
      ($0.GetDefaultPatchBaselineRequest value) => value.writeToBuffer(),
      $0.GetDefaultPatchBaselineResult.fromBuffer);
  static final _$getDeployablePatchSnapshotForInstance = $grpc.ClientMethod<
          $0.GetDeployablePatchSnapshotForInstanceRequest,
          $0.GetDeployablePatchSnapshotForInstanceResult>(
      '/ssm.SSMService/GetDeployablePatchSnapshotForInstance',
      ($0.GetDeployablePatchSnapshotForInstanceRequest value) =>
          value.writeToBuffer(),
      $0.GetDeployablePatchSnapshotForInstanceResult.fromBuffer);
  static final _$getDocument =
      $grpc.ClientMethod<$0.GetDocumentRequest, $0.GetDocumentResult>(
          '/ssm.SSMService/GetDocument',
          ($0.GetDocumentRequest value) => value.writeToBuffer(),
          $0.GetDocumentResult.fromBuffer);
  static final _$getExecutionPreview = $grpc.ClientMethod<
          $0.GetExecutionPreviewRequest, $0.GetExecutionPreviewResponse>(
      '/ssm.SSMService/GetExecutionPreview',
      ($0.GetExecutionPreviewRequest value) => value.writeToBuffer(),
      $0.GetExecutionPreviewResponse.fromBuffer);
  static final _$getInventory =
      $grpc.ClientMethod<$0.GetInventoryRequest, $0.GetInventoryResult>(
          '/ssm.SSMService/GetInventory',
          ($0.GetInventoryRequest value) => value.writeToBuffer(),
          $0.GetInventoryResult.fromBuffer);
  static final _$getInventorySchema = $grpc.ClientMethod<
          $0.GetInventorySchemaRequest, $0.GetInventorySchemaResult>(
      '/ssm.SSMService/GetInventorySchema',
      ($0.GetInventorySchemaRequest value) => value.writeToBuffer(),
      $0.GetInventorySchemaResult.fromBuffer);
  static final _$getMaintenanceWindow = $grpc.ClientMethod<
          $0.GetMaintenanceWindowRequest, $0.GetMaintenanceWindowResult>(
      '/ssm.SSMService/GetMaintenanceWindow',
      ($0.GetMaintenanceWindowRequest value) => value.writeToBuffer(),
      $0.GetMaintenanceWindowResult.fromBuffer);
  static final _$getMaintenanceWindowExecution = $grpc.ClientMethod<
          $0.GetMaintenanceWindowExecutionRequest,
          $0.GetMaintenanceWindowExecutionResult>(
      '/ssm.SSMService/GetMaintenanceWindowExecution',
      ($0.GetMaintenanceWindowExecutionRequest value) => value.writeToBuffer(),
      $0.GetMaintenanceWindowExecutionResult.fromBuffer);
  static final _$getMaintenanceWindowExecutionTask = $grpc.ClientMethod<
          $0.GetMaintenanceWindowExecutionTaskRequest,
          $0.GetMaintenanceWindowExecutionTaskResult>(
      '/ssm.SSMService/GetMaintenanceWindowExecutionTask',
      ($0.GetMaintenanceWindowExecutionTaskRequest value) =>
          value.writeToBuffer(),
      $0.GetMaintenanceWindowExecutionTaskResult.fromBuffer);
  static final _$getMaintenanceWindowExecutionTaskInvocation =
      $grpc.ClientMethod<$0.GetMaintenanceWindowExecutionTaskInvocationRequest,
              $0.GetMaintenanceWindowExecutionTaskInvocationResult>(
          '/ssm.SSMService/GetMaintenanceWindowExecutionTaskInvocation',
          ($0.GetMaintenanceWindowExecutionTaskInvocationRequest value) =>
              value.writeToBuffer(),
          $0.GetMaintenanceWindowExecutionTaskInvocationResult.fromBuffer);
  static final _$getMaintenanceWindowTask = $grpc.ClientMethod<
          $0.GetMaintenanceWindowTaskRequest,
          $0.GetMaintenanceWindowTaskResult>(
      '/ssm.SSMService/GetMaintenanceWindowTask',
      ($0.GetMaintenanceWindowTaskRequest value) => value.writeToBuffer(),
      $0.GetMaintenanceWindowTaskResult.fromBuffer);
  static final _$getOpsItem =
      $grpc.ClientMethod<$0.GetOpsItemRequest, $0.GetOpsItemResponse>(
          '/ssm.SSMService/GetOpsItem',
          ($0.GetOpsItemRequest value) => value.writeToBuffer(),
          $0.GetOpsItemResponse.fromBuffer);
  static final _$getOpsMetadata =
      $grpc.ClientMethod<$0.GetOpsMetadataRequest, $0.GetOpsMetadataResult>(
          '/ssm.SSMService/GetOpsMetadata',
          ($0.GetOpsMetadataRequest value) => value.writeToBuffer(),
          $0.GetOpsMetadataResult.fromBuffer);
  static final _$getOpsSummary =
      $grpc.ClientMethod<$0.GetOpsSummaryRequest, $0.GetOpsSummaryResult>(
          '/ssm.SSMService/GetOpsSummary',
          ($0.GetOpsSummaryRequest value) => value.writeToBuffer(),
          $0.GetOpsSummaryResult.fromBuffer);
  static final _$getParameter =
      $grpc.ClientMethod<$0.GetParameterRequest, $0.GetParameterResult>(
          '/ssm.SSMService/GetParameter',
          ($0.GetParameterRequest value) => value.writeToBuffer(),
          $0.GetParameterResult.fromBuffer);
  static final _$getParameterHistory = $grpc.ClientMethod<
          $0.GetParameterHistoryRequest, $0.GetParameterHistoryResult>(
      '/ssm.SSMService/GetParameterHistory',
      ($0.GetParameterHistoryRequest value) => value.writeToBuffer(),
      $0.GetParameterHistoryResult.fromBuffer);
  static final _$getParameters =
      $grpc.ClientMethod<$0.GetParametersRequest, $0.GetParametersResult>(
          '/ssm.SSMService/GetParameters',
          ($0.GetParametersRequest value) => value.writeToBuffer(),
          $0.GetParametersResult.fromBuffer);
  static final _$getParametersByPath = $grpc.ClientMethod<
          $0.GetParametersByPathRequest, $0.GetParametersByPathResult>(
      '/ssm.SSMService/GetParametersByPath',
      ($0.GetParametersByPathRequest value) => value.writeToBuffer(),
      $0.GetParametersByPathResult.fromBuffer);
  static final _$getPatchBaseline =
      $grpc.ClientMethod<$0.GetPatchBaselineRequest, $0.GetPatchBaselineResult>(
          '/ssm.SSMService/GetPatchBaseline',
          ($0.GetPatchBaselineRequest value) => value.writeToBuffer(),
          $0.GetPatchBaselineResult.fromBuffer);
  static final _$getPatchBaselineForPatchGroup = $grpc.ClientMethod<
          $0.GetPatchBaselineForPatchGroupRequest,
          $0.GetPatchBaselineForPatchGroupResult>(
      '/ssm.SSMService/GetPatchBaselineForPatchGroup',
      ($0.GetPatchBaselineForPatchGroupRequest value) => value.writeToBuffer(),
      $0.GetPatchBaselineForPatchGroupResult.fromBuffer);
  static final _$getResourcePolicies = $grpc.ClientMethod<
          $0.GetResourcePoliciesRequest, $0.GetResourcePoliciesResponse>(
      '/ssm.SSMService/GetResourcePolicies',
      ($0.GetResourcePoliciesRequest value) => value.writeToBuffer(),
      $0.GetResourcePoliciesResponse.fromBuffer);
  static final _$getServiceSetting = $grpc.ClientMethod<
          $0.GetServiceSettingRequest, $0.GetServiceSettingResult>(
      '/ssm.SSMService/GetServiceSetting',
      ($0.GetServiceSettingRequest value) => value.writeToBuffer(),
      $0.GetServiceSettingResult.fromBuffer);
  static final _$labelParameterVersion = $grpc.ClientMethod<
          $0.LabelParameterVersionRequest, $0.LabelParameterVersionResult>(
      '/ssm.SSMService/LabelParameterVersion',
      ($0.LabelParameterVersionRequest value) => value.writeToBuffer(),
      $0.LabelParameterVersionResult.fromBuffer);
  static final _$listAssociations =
      $grpc.ClientMethod<$0.ListAssociationsRequest, $0.ListAssociationsResult>(
          '/ssm.SSMService/ListAssociations',
          ($0.ListAssociationsRequest value) => value.writeToBuffer(),
          $0.ListAssociationsResult.fromBuffer);
  static final _$listAssociationVersions = $grpc.ClientMethod<
          $0.ListAssociationVersionsRequest, $0.ListAssociationVersionsResult>(
      '/ssm.SSMService/ListAssociationVersions',
      ($0.ListAssociationVersionsRequest value) => value.writeToBuffer(),
      $0.ListAssociationVersionsResult.fromBuffer);
  static final _$listCommandInvocations = $grpc.ClientMethod<
          $0.ListCommandInvocationsRequest, $0.ListCommandInvocationsResult>(
      '/ssm.SSMService/ListCommandInvocations',
      ($0.ListCommandInvocationsRequest value) => value.writeToBuffer(),
      $0.ListCommandInvocationsResult.fromBuffer);
  static final _$listCommands =
      $grpc.ClientMethod<$0.ListCommandsRequest, $0.ListCommandsResult>(
          '/ssm.SSMService/ListCommands',
          ($0.ListCommandsRequest value) => value.writeToBuffer(),
          $0.ListCommandsResult.fromBuffer);
  static final _$listComplianceItems = $grpc.ClientMethod<
          $0.ListComplianceItemsRequest, $0.ListComplianceItemsResult>(
      '/ssm.SSMService/ListComplianceItems',
      ($0.ListComplianceItemsRequest value) => value.writeToBuffer(),
      $0.ListComplianceItemsResult.fromBuffer);
  static final _$listComplianceSummaries = $grpc.ClientMethod<
          $0.ListComplianceSummariesRequest, $0.ListComplianceSummariesResult>(
      '/ssm.SSMService/ListComplianceSummaries',
      ($0.ListComplianceSummariesRequest value) => value.writeToBuffer(),
      $0.ListComplianceSummariesResult.fromBuffer);
  static final _$listDocumentMetadataHistory = $grpc.ClientMethod<
          $0.ListDocumentMetadataHistoryRequest,
          $0.ListDocumentMetadataHistoryResponse>(
      '/ssm.SSMService/ListDocumentMetadataHistory',
      ($0.ListDocumentMetadataHistoryRequest value) => value.writeToBuffer(),
      $0.ListDocumentMetadataHistoryResponse.fromBuffer);
  static final _$listDocuments =
      $grpc.ClientMethod<$0.ListDocumentsRequest, $0.ListDocumentsResult>(
          '/ssm.SSMService/ListDocuments',
          ($0.ListDocumentsRequest value) => value.writeToBuffer(),
          $0.ListDocumentsResult.fromBuffer);
  static final _$listDocumentVersions = $grpc.ClientMethod<
          $0.ListDocumentVersionsRequest, $0.ListDocumentVersionsResult>(
      '/ssm.SSMService/ListDocumentVersions',
      ($0.ListDocumentVersionsRequest value) => value.writeToBuffer(),
      $0.ListDocumentVersionsResult.fromBuffer);
  static final _$listInventoryEntries = $grpc.ClientMethod<
          $0.ListInventoryEntriesRequest, $0.ListInventoryEntriesResult>(
      '/ssm.SSMService/ListInventoryEntries',
      ($0.ListInventoryEntriesRequest value) => value.writeToBuffer(),
      $0.ListInventoryEntriesResult.fromBuffer);
  static final _$listNodes =
      $grpc.ClientMethod<$0.ListNodesRequest, $0.ListNodesResult>(
          '/ssm.SSMService/ListNodes',
          ($0.ListNodesRequest value) => value.writeToBuffer(),
          $0.ListNodesResult.fromBuffer);
  static final _$listNodesSummary =
      $grpc.ClientMethod<$0.ListNodesSummaryRequest, $0.ListNodesSummaryResult>(
          '/ssm.SSMService/ListNodesSummary',
          ($0.ListNodesSummaryRequest value) => value.writeToBuffer(),
          $0.ListNodesSummaryResult.fromBuffer);
  static final _$listOpsItemEvents = $grpc.ClientMethod<
          $0.ListOpsItemEventsRequest, $0.ListOpsItemEventsResponse>(
      '/ssm.SSMService/ListOpsItemEvents',
      ($0.ListOpsItemEventsRequest value) => value.writeToBuffer(),
      $0.ListOpsItemEventsResponse.fromBuffer);
  static final _$listOpsItemRelatedItems = $grpc.ClientMethod<
          $0.ListOpsItemRelatedItemsRequest,
          $0.ListOpsItemRelatedItemsResponse>(
      '/ssm.SSMService/ListOpsItemRelatedItems',
      ($0.ListOpsItemRelatedItemsRequest value) => value.writeToBuffer(),
      $0.ListOpsItemRelatedItemsResponse.fromBuffer);
  static final _$listOpsMetadata =
      $grpc.ClientMethod<$0.ListOpsMetadataRequest, $0.ListOpsMetadataResult>(
          '/ssm.SSMService/ListOpsMetadata',
          ($0.ListOpsMetadataRequest value) => value.writeToBuffer(),
          $0.ListOpsMetadataResult.fromBuffer);
  static final _$listResourceComplianceSummaries = $grpc.ClientMethod<
          $0.ListResourceComplianceSummariesRequest,
          $0.ListResourceComplianceSummariesResult>(
      '/ssm.SSMService/ListResourceComplianceSummaries',
      ($0.ListResourceComplianceSummariesRequest value) =>
          value.writeToBuffer(),
      $0.ListResourceComplianceSummariesResult.fromBuffer);
  static final _$listResourceDataSync = $grpc.ClientMethod<
          $0.ListResourceDataSyncRequest, $0.ListResourceDataSyncResult>(
      '/ssm.SSMService/ListResourceDataSync',
      ($0.ListResourceDataSyncRequest value) => value.writeToBuffer(),
      $0.ListResourceDataSyncResult.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResult>(
      '/ssm.SSMService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResult.fromBuffer);
  static final _$modifyDocumentPermission = $grpc.ClientMethod<
          $0.ModifyDocumentPermissionRequest,
          $0.ModifyDocumentPermissionResponse>(
      '/ssm.SSMService/ModifyDocumentPermission',
      ($0.ModifyDocumentPermissionRequest value) => value.writeToBuffer(),
      $0.ModifyDocumentPermissionResponse.fromBuffer);
  static final _$putComplianceItems = $grpc.ClientMethod<
          $0.PutComplianceItemsRequest, $0.PutComplianceItemsResult>(
      '/ssm.SSMService/PutComplianceItems',
      ($0.PutComplianceItemsRequest value) => value.writeToBuffer(),
      $0.PutComplianceItemsResult.fromBuffer);
  static final _$putInventory =
      $grpc.ClientMethod<$0.PutInventoryRequest, $0.PutInventoryResult>(
          '/ssm.SSMService/PutInventory',
          ($0.PutInventoryRequest value) => value.writeToBuffer(),
          $0.PutInventoryResult.fromBuffer);
  static final _$putParameter =
      $grpc.ClientMethod<$0.PutParameterRequest, $0.PutParameterResult>(
          '/ssm.SSMService/PutParameter',
          ($0.PutParameterRequest value) => value.writeToBuffer(),
          $0.PutParameterResult.fromBuffer);
  static final _$putResourcePolicy = $grpc.ClientMethod<
          $0.PutResourcePolicyRequest, $0.PutResourcePolicyResponse>(
      '/ssm.SSMService/PutResourcePolicy',
      ($0.PutResourcePolicyRequest value) => value.writeToBuffer(),
      $0.PutResourcePolicyResponse.fromBuffer);
  static final _$registerDefaultPatchBaseline = $grpc.ClientMethod<
          $0.RegisterDefaultPatchBaselineRequest,
          $0.RegisterDefaultPatchBaselineResult>(
      '/ssm.SSMService/RegisterDefaultPatchBaseline',
      ($0.RegisterDefaultPatchBaselineRequest value) => value.writeToBuffer(),
      $0.RegisterDefaultPatchBaselineResult.fromBuffer);
  static final _$registerPatchBaselineForPatchGroup = $grpc.ClientMethod<
          $0.RegisterPatchBaselineForPatchGroupRequest,
          $0.RegisterPatchBaselineForPatchGroupResult>(
      '/ssm.SSMService/RegisterPatchBaselineForPatchGroup',
      ($0.RegisterPatchBaselineForPatchGroupRequest value) =>
          value.writeToBuffer(),
      $0.RegisterPatchBaselineForPatchGroupResult.fromBuffer);
  static final _$registerTargetWithMaintenanceWindow = $grpc.ClientMethod<
          $0.RegisterTargetWithMaintenanceWindowRequest,
          $0.RegisterTargetWithMaintenanceWindowResult>(
      '/ssm.SSMService/RegisterTargetWithMaintenanceWindow',
      ($0.RegisterTargetWithMaintenanceWindowRequest value) =>
          value.writeToBuffer(),
      $0.RegisterTargetWithMaintenanceWindowResult.fromBuffer);
  static final _$registerTaskWithMaintenanceWindow = $grpc.ClientMethod<
          $0.RegisterTaskWithMaintenanceWindowRequest,
          $0.RegisterTaskWithMaintenanceWindowResult>(
      '/ssm.SSMService/RegisterTaskWithMaintenanceWindow',
      ($0.RegisterTaskWithMaintenanceWindowRequest value) =>
          value.writeToBuffer(),
      $0.RegisterTaskWithMaintenanceWindowResult.fromBuffer);
  static final _$removeTagsFromResource = $grpc.ClientMethod<
          $0.RemoveTagsFromResourceRequest, $0.RemoveTagsFromResourceResult>(
      '/ssm.SSMService/RemoveTagsFromResource',
      ($0.RemoveTagsFromResourceRequest value) => value.writeToBuffer(),
      $0.RemoveTagsFromResourceResult.fromBuffer);
  static final _$resetServiceSetting = $grpc.ClientMethod<
          $0.ResetServiceSettingRequest, $0.ResetServiceSettingResult>(
      '/ssm.SSMService/ResetServiceSetting',
      ($0.ResetServiceSettingRequest value) => value.writeToBuffer(),
      $0.ResetServiceSettingResult.fromBuffer);
  static final _$resumeSession =
      $grpc.ClientMethod<$0.ResumeSessionRequest, $0.ResumeSessionResponse>(
          '/ssm.SSMService/ResumeSession',
          ($0.ResumeSessionRequest value) => value.writeToBuffer(),
          $0.ResumeSessionResponse.fromBuffer);
  static final _$sendAutomationSignal = $grpc.ClientMethod<
          $0.SendAutomationSignalRequest, $0.SendAutomationSignalResult>(
      '/ssm.SSMService/SendAutomationSignal',
      ($0.SendAutomationSignalRequest value) => value.writeToBuffer(),
      $0.SendAutomationSignalResult.fromBuffer);
  static final _$sendCommand =
      $grpc.ClientMethod<$0.SendCommandRequest, $0.SendCommandResult>(
          '/ssm.SSMService/SendCommand',
          ($0.SendCommandRequest value) => value.writeToBuffer(),
          $0.SendCommandResult.fromBuffer);
  static final _$startAccessRequest = $grpc.ClientMethod<
          $0.StartAccessRequestRequest, $0.StartAccessRequestResponse>(
      '/ssm.SSMService/StartAccessRequest',
      ($0.StartAccessRequestRequest value) => value.writeToBuffer(),
      $0.StartAccessRequestResponse.fromBuffer);
  static final _$startAssociationsOnce = $grpc.ClientMethod<
          $0.StartAssociationsOnceRequest, $0.StartAssociationsOnceResult>(
      '/ssm.SSMService/StartAssociationsOnce',
      ($0.StartAssociationsOnceRequest value) => value.writeToBuffer(),
      $0.StartAssociationsOnceResult.fromBuffer);
  static final _$startAutomationExecution = $grpc.ClientMethod<
          $0.StartAutomationExecutionRequest,
          $0.StartAutomationExecutionResult>(
      '/ssm.SSMService/StartAutomationExecution',
      ($0.StartAutomationExecutionRequest value) => value.writeToBuffer(),
      $0.StartAutomationExecutionResult.fromBuffer);
  static final _$startChangeRequestExecution = $grpc.ClientMethod<
          $0.StartChangeRequestExecutionRequest,
          $0.StartChangeRequestExecutionResult>(
      '/ssm.SSMService/StartChangeRequestExecution',
      ($0.StartChangeRequestExecutionRequest value) => value.writeToBuffer(),
      $0.StartChangeRequestExecutionResult.fromBuffer);
  static final _$startExecutionPreview = $grpc.ClientMethod<
          $0.StartExecutionPreviewRequest, $0.StartExecutionPreviewResponse>(
      '/ssm.SSMService/StartExecutionPreview',
      ($0.StartExecutionPreviewRequest value) => value.writeToBuffer(),
      $0.StartExecutionPreviewResponse.fromBuffer);
  static final _$startSession =
      $grpc.ClientMethod<$0.StartSessionRequest, $0.StartSessionResponse>(
          '/ssm.SSMService/StartSession',
          ($0.StartSessionRequest value) => value.writeToBuffer(),
          $0.StartSessionResponse.fromBuffer);
  static final _$stopAutomationExecution = $grpc.ClientMethod<
          $0.StopAutomationExecutionRequest, $0.StopAutomationExecutionResult>(
      '/ssm.SSMService/StopAutomationExecution',
      ($0.StopAutomationExecutionRequest value) => value.writeToBuffer(),
      $0.StopAutomationExecutionResult.fromBuffer);
  static final _$terminateSession = $grpc.ClientMethod<
          $0.TerminateSessionRequest, $0.TerminateSessionResponse>(
      '/ssm.SSMService/TerminateSession',
      ($0.TerminateSessionRequest value) => value.writeToBuffer(),
      $0.TerminateSessionResponse.fromBuffer);
  static final _$unlabelParameterVersion = $grpc.ClientMethod<
          $0.UnlabelParameterVersionRequest, $0.UnlabelParameterVersionResult>(
      '/ssm.SSMService/UnlabelParameterVersion',
      ($0.UnlabelParameterVersionRequest value) => value.writeToBuffer(),
      $0.UnlabelParameterVersionResult.fromBuffer);
  static final _$updateAssociation = $grpc.ClientMethod<
          $0.UpdateAssociationRequest, $0.UpdateAssociationResult>(
      '/ssm.SSMService/UpdateAssociation',
      ($0.UpdateAssociationRequest value) => value.writeToBuffer(),
      $0.UpdateAssociationResult.fromBuffer);
  static final _$updateAssociationStatus = $grpc.ClientMethod<
          $0.UpdateAssociationStatusRequest, $0.UpdateAssociationStatusResult>(
      '/ssm.SSMService/UpdateAssociationStatus',
      ($0.UpdateAssociationStatusRequest value) => value.writeToBuffer(),
      $0.UpdateAssociationStatusResult.fromBuffer);
  static final _$updateDocument =
      $grpc.ClientMethod<$0.UpdateDocumentRequest, $0.UpdateDocumentResult>(
          '/ssm.SSMService/UpdateDocument',
          ($0.UpdateDocumentRequest value) => value.writeToBuffer(),
          $0.UpdateDocumentResult.fromBuffer);
  static final _$updateDocumentDefaultVersion = $grpc.ClientMethod<
          $0.UpdateDocumentDefaultVersionRequest,
          $0.UpdateDocumentDefaultVersionResult>(
      '/ssm.SSMService/UpdateDocumentDefaultVersion',
      ($0.UpdateDocumentDefaultVersionRequest value) => value.writeToBuffer(),
      $0.UpdateDocumentDefaultVersionResult.fromBuffer);
  static final _$updateDocumentMetadata = $grpc.ClientMethod<
          $0.UpdateDocumentMetadataRequest, $0.UpdateDocumentMetadataResponse>(
      '/ssm.SSMService/UpdateDocumentMetadata',
      ($0.UpdateDocumentMetadataRequest value) => value.writeToBuffer(),
      $0.UpdateDocumentMetadataResponse.fromBuffer);
  static final _$updateMaintenanceWindow = $grpc.ClientMethod<
          $0.UpdateMaintenanceWindowRequest, $0.UpdateMaintenanceWindowResult>(
      '/ssm.SSMService/UpdateMaintenanceWindow',
      ($0.UpdateMaintenanceWindowRequest value) => value.writeToBuffer(),
      $0.UpdateMaintenanceWindowResult.fromBuffer);
  static final _$updateMaintenanceWindowTarget = $grpc.ClientMethod<
          $0.UpdateMaintenanceWindowTargetRequest,
          $0.UpdateMaintenanceWindowTargetResult>(
      '/ssm.SSMService/UpdateMaintenanceWindowTarget',
      ($0.UpdateMaintenanceWindowTargetRequest value) => value.writeToBuffer(),
      $0.UpdateMaintenanceWindowTargetResult.fromBuffer);
  static final _$updateMaintenanceWindowTask = $grpc.ClientMethod<
          $0.UpdateMaintenanceWindowTaskRequest,
          $0.UpdateMaintenanceWindowTaskResult>(
      '/ssm.SSMService/UpdateMaintenanceWindowTask',
      ($0.UpdateMaintenanceWindowTaskRequest value) => value.writeToBuffer(),
      $0.UpdateMaintenanceWindowTaskResult.fromBuffer);
  static final _$updateManagedInstanceRole = $grpc.ClientMethod<
          $0.UpdateManagedInstanceRoleRequest,
          $0.UpdateManagedInstanceRoleResult>(
      '/ssm.SSMService/UpdateManagedInstanceRole',
      ($0.UpdateManagedInstanceRoleRequest value) => value.writeToBuffer(),
      $0.UpdateManagedInstanceRoleResult.fromBuffer);
  static final _$updateOpsItem =
      $grpc.ClientMethod<$0.UpdateOpsItemRequest, $0.UpdateOpsItemResponse>(
          '/ssm.SSMService/UpdateOpsItem',
          ($0.UpdateOpsItemRequest value) => value.writeToBuffer(),
          $0.UpdateOpsItemResponse.fromBuffer);
  static final _$updateOpsMetadata = $grpc.ClientMethod<
          $0.UpdateOpsMetadataRequest, $0.UpdateOpsMetadataResult>(
      '/ssm.SSMService/UpdateOpsMetadata',
      ($0.UpdateOpsMetadataRequest value) => value.writeToBuffer(),
      $0.UpdateOpsMetadataResult.fromBuffer);
  static final _$updatePatchBaseline = $grpc.ClientMethod<
          $0.UpdatePatchBaselineRequest, $0.UpdatePatchBaselineResult>(
      '/ssm.SSMService/UpdatePatchBaseline',
      ($0.UpdatePatchBaselineRequest value) => value.writeToBuffer(),
      $0.UpdatePatchBaselineResult.fromBuffer);
  static final _$updateResourceDataSync = $grpc.ClientMethod<
          $0.UpdateResourceDataSyncRequest, $0.UpdateResourceDataSyncResult>(
      '/ssm.SSMService/UpdateResourceDataSync',
      ($0.UpdateResourceDataSyncRequest value) => value.writeToBuffer(),
      $0.UpdateResourceDataSyncResult.fromBuffer);
  static final _$updateServiceSetting = $grpc.ClientMethod<
          $0.UpdateServiceSettingRequest, $0.UpdateServiceSettingResult>(
      '/ssm.SSMService/UpdateServiceSetting',
      ($0.UpdateServiceSettingRequest value) => value.writeToBuffer(),
      $0.UpdateServiceSettingResult.fromBuffer);
}

@$pb.GrpcServiceName('ssm.SSMService')
abstract class SSMServiceBase extends $grpc.Service {
  $core.String get $name => 'ssm.SSMService';

  SSMServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AddTagsToResourceRequest,
            $0.AddTagsToResourceResult>(
        'AddTagsToResource',
        addTagsToResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddTagsToResourceRequest.fromBuffer(value),
        ($0.AddTagsToResourceResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssociateOpsItemRelatedItemRequest,
            $0.AssociateOpsItemRelatedItemResponse>(
        'AssociateOpsItemRelatedItem',
        associateOpsItemRelatedItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateOpsItemRelatedItemRequest.fromBuffer(value),
        ($0.AssociateOpsItemRelatedItemResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CancelCommandRequest, $0.CancelCommandResult>(
            'CancelCommand',
            cancelCommand_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CancelCommandRequest.fromBuffer(value),
            ($0.CancelCommandResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CancelMaintenanceWindowExecutionRequest,
            $0.CancelMaintenanceWindowExecutionResult>(
        'CancelMaintenanceWindowExecution',
        cancelMaintenanceWindowExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CancelMaintenanceWindowExecutionRequest.fromBuffer(value),
        ($0.CancelMaintenanceWindowExecutionResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateActivationRequest,
            $0.CreateActivationResult>(
        'CreateActivation',
        createActivation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateActivationRequest.fromBuffer(value),
        ($0.CreateActivationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateAssociationRequest,
            $0.CreateAssociationResult>(
        'CreateAssociation',
        createAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAssociationRequest.fromBuffer(value),
        ($0.CreateAssociationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateAssociationBatchRequest,
            $0.CreateAssociationBatchResult>(
        'CreateAssociationBatch',
        createAssociationBatch_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAssociationBatchRequest.fromBuffer(value),
        ($0.CreateAssociationBatchResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateDocumentRequest, $0.CreateDocumentResult>(
            'CreateDocument',
            createDocument_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateDocumentRequest.fromBuffer(value),
            ($0.CreateDocumentResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateMaintenanceWindowRequest,
            $0.CreateMaintenanceWindowResult>(
        'CreateMaintenanceWindow',
        createMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateMaintenanceWindowRequest.fromBuffer(value),
        ($0.CreateMaintenanceWindowResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateOpsItemRequest, $0.CreateOpsItemResponse>(
            'CreateOpsItem',
            createOpsItem_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateOpsItemRequest.fromBuffer(value),
            ($0.CreateOpsItemResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateOpsMetadataRequest,
            $0.CreateOpsMetadataResult>(
        'CreateOpsMetadata',
        createOpsMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateOpsMetadataRequest.fromBuffer(value),
        ($0.CreateOpsMetadataResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePatchBaselineRequest,
            $0.CreatePatchBaselineResult>(
        'CreatePatchBaseline',
        createPatchBaseline_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePatchBaselineRequest.fromBuffer(value),
        ($0.CreatePatchBaselineResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateResourceDataSyncRequest,
            $0.CreateResourceDataSyncResult>(
        'CreateResourceDataSync',
        createResourceDataSync_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateResourceDataSyncRequest.fromBuffer(value),
        ($0.CreateResourceDataSyncResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteActivationRequest,
            $0.DeleteActivationResult>(
        'DeleteActivation',
        deleteActivation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteActivationRequest.fromBuffer(value),
        ($0.DeleteActivationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAssociationRequest,
            $0.DeleteAssociationResult>(
        'DeleteAssociation',
        deleteAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAssociationRequest.fromBuffer(value),
        ($0.DeleteAssociationResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteDocumentRequest, $0.DeleteDocumentResult>(
            'DeleteDocument',
            deleteDocument_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteDocumentRequest.fromBuffer(value),
            ($0.DeleteDocumentResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteInventoryRequest,
            $0.DeleteInventoryResult>(
        'DeleteInventory',
        deleteInventory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteInventoryRequest.fromBuffer(value),
        ($0.DeleteInventoryResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteMaintenanceWindowRequest,
            $0.DeleteMaintenanceWindowResult>(
        'DeleteMaintenanceWindow',
        deleteMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteMaintenanceWindowRequest.fromBuffer(value),
        ($0.DeleteMaintenanceWindowResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteOpsItemRequest, $0.DeleteOpsItemResponse>(
            'DeleteOpsItem',
            deleteOpsItem_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteOpsItemRequest.fromBuffer(value),
            ($0.DeleteOpsItemResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteOpsMetadataRequest,
            $0.DeleteOpsMetadataResult>(
        'DeleteOpsMetadata',
        deleteOpsMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteOpsMetadataRequest.fromBuffer(value),
        ($0.DeleteOpsMetadataResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteParameterRequest,
            $0.DeleteParameterResult>(
        'DeleteParameter',
        deleteParameter_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteParameterRequest.fromBuffer(value),
        ($0.DeleteParameterResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteParametersRequest,
            $0.DeleteParametersResult>(
        'DeleteParameters',
        deleteParameters_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteParametersRequest.fromBuffer(value),
        ($0.DeleteParametersResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePatchBaselineRequest,
            $0.DeletePatchBaselineResult>(
        'DeletePatchBaseline',
        deletePatchBaseline_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePatchBaselineRequest.fromBuffer(value),
        ($0.DeletePatchBaselineResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourceDataSyncRequest,
            $0.DeleteResourceDataSyncResult>(
        'DeleteResourceDataSync',
        deleteResourceDataSync_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourceDataSyncRequest.fromBuffer(value),
        ($0.DeleteResourceDataSyncResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourcePolicyRequest,
            $0.DeleteResourcePolicyResponse>(
        'DeleteResourcePolicy',
        deleteResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourcePolicyRequest.fromBuffer(value),
        ($0.DeleteResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeregisterManagedInstanceRequest,
            $0.DeregisterManagedInstanceResult>(
        'DeregisterManagedInstance',
        deregisterManagedInstance_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeregisterManagedInstanceRequest.fromBuffer(value),
        ($0.DeregisterManagedInstanceResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeregisterPatchBaselineForPatchGroupRequest,
            $0.DeregisterPatchBaselineForPatchGroupResult>(
        'DeregisterPatchBaselineForPatchGroup',
        deregisterPatchBaselineForPatchGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeregisterPatchBaselineForPatchGroupRequest.fromBuffer(value),
        ($0.DeregisterPatchBaselineForPatchGroupResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeregisterTargetFromMaintenanceWindowRequest,
            $0.DeregisterTargetFromMaintenanceWindowResult>(
        'DeregisterTargetFromMaintenanceWindow',
        deregisterTargetFromMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeregisterTargetFromMaintenanceWindowRequest.fromBuffer(value),
        ($0.DeregisterTargetFromMaintenanceWindowResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DeregisterTaskFromMaintenanceWindowRequest,
            $0.DeregisterTaskFromMaintenanceWindowResult>(
        'DeregisterTaskFromMaintenanceWindow',
        deregisterTaskFromMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeregisterTaskFromMaintenanceWindowRequest.fromBuffer(value),
        ($0.DeregisterTaskFromMaintenanceWindowResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeActivationsRequest,
            $0.DescribeActivationsResult>(
        'DescribeActivations',
        describeActivations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeActivationsRequest.fromBuffer(value),
        ($0.DescribeActivationsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAssociationRequest,
            $0.DescribeAssociationResult>(
        'DescribeAssociation',
        describeAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAssociationRequest.fromBuffer(value),
        ($0.DescribeAssociationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAssociationExecutionsRequest,
            $0.DescribeAssociationExecutionsResult>(
        'DescribeAssociationExecutions',
        describeAssociationExecutions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAssociationExecutionsRequest.fromBuffer(value),
        ($0.DescribeAssociationExecutionsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeAssociationExecutionTargetsRequest,
            $0.DescribeAssociationExecutionTargetsResult>(
        'DescribeAssociationExecutionTargets',
        describeAssociationExecutionTargets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAssociationExecutionTargetsRequest.fromBuffer(value),
        ($0.DescribeAssociationExecutionTargetsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAutomationExecutionsRequest,
            $0.DescribeAutomationExecutionsResult>(
        'DescribeAutomationExecutions',
        describeAutomationExecutions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAutomationExecutionsRequest.fromBuffer(value),
        ($0.DescribeAutomationExecutionsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAutomationStepExecutionsRequest,
            $0.DescribeAutomationStepExecutionsResult>(
        'DescribeAutomationStepExecutions',
        describeAutomationStepExecutions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAutomationStepExecutionsRequest.fromBuffer(value),
        ($0.DescribeAutomationStepExecutionsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeAvailablePatchesRequest,
            $0.DescribeAvailablePatchesResult>(
        'DescribeAvailablePatches',
        describeAvailablePatches_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeAvailablePatchesRequest.fromBuffer(value),
        ($0.DescribeAvailablePatchesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeDocumentRequest,
            $0.DescribeDocumentResult>(
        'DescribeDocument',
        describeDocument_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeDocumentRequest.fromBuffer(value),
        ($0.DescribeDocumentResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeDocumentPermissionRequest,
            $0.DescribeDocumentPermissionResponse>(
        'DescribeDocumentPermission',
        describeDocumentPermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeDocumentPermissionRequest.fromBuffer(value),
        ($0.DescribeDocumentPermissionResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeEffectiveInstanceAssociationsRequest,
            $0.DescribeEffectiveInstanceAssociationsResult>(
        'DescribeEffectiveInstanceAssociations',
        describeEffectiveInstanceAssociations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeEffectiveInstanceAssociationsRequest.fromBuffer(value),
        ($0.DescribeEffectiveInstanceAssociationsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeEffectivePatchesForPatchBaselineRequest,
            $0.DescribeEffectivePatchesForPatchBaselineResult>(
        'DescribeEffectivePatchesForPatchBaseline',
        describeEffectivePatchesForPatchBaseline_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeEffectivePatchesForPatchBaselineRequest.fromBuffer(
                value),
        ($0.DescribeEffectivePatchesForPatchBaselineResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeInstanceAssociationsStatusRequest,
            $0.DescribeInstanceAssociationsStatusResult>(
        'DescribeInstanceAssociationsStatus',
        describeInstanceAssociationsStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInstanceAssociationsStatusRequest.fromBuffer(value),
        ($0.DescribeInstanceAssociationsStatusResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeInstanceInformationRequest,
            $0.DescribeInstanceInformationResult>(
        'DescribeInstanceInformation',
        describeInstanceInformation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInstanceInformationRequest.fromBuffer(value),
        ($0.DescribeInstanceInformationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeInstancePatchesRequest,
            $0.DescribeInstancePatchesResult>(
        'DescribeInstancePatches',
        describeInstancePatches_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInstancePatchesRequest.fromBuffer(value),
        ($0.DescribeInstancePatchesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeInstancePatchStatesRequest,
            $0.DescribeInstancePatchStatesResult>(
        'DescribeInstancePatchStates',
        describeInstancePatchStates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInstancePatchStatesRequest.fromBuffer(value),
        ($0.DescribeInstancePatchStatesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeInstancePatchStatesForPatchGroupRequest,
            $0.DescribeInstancePatchStatesForPatchGroupResult>(
        'DescribeInstancePatchStatesForPatchGroup',
        describeInstancePatchStatesForPatchGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInstancePatchStatesForPatchGroupRequest.fromBuffer(
                value),
        ($0.DescribeInstancePatchStatesForPatchGroupResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeInstancePropertiesRequest,
            $0.DescribeInstancePropertiesResult>(
        'DescribeInstanceProperties',
        describeInstanceProperties_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInstancePropertiesRequest.fromBuffer(value),
        ($0.DescribeInstancePropertiesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeInventoryDeletionsRequest,
            $0.DescribeInventoryDeletionsResult>(
        'DescribeInventoryDeletions',
        describeInventoryDeletions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeInventoryDeletionsRequest.fromBuffer(value),
        ($0.DescribeInventoryDeletionsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeMaintenanceWindowExecutionsRequest,
            $0.DescribeMaintenanceWindowExecutionsResult>(
        'DescribeMaintenanceWindowExecutions',
        describeMaintenanceWindowExecutions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowExecutionsRequest.fromBuffer(value),
        ($0.DescribeMaintenanceWindowExecutionsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeMaintenanceWindowExecutionTaskInvocationsRequest,
            $0.DescribeMaintenanceWindowExecutionTaskInvocationsResult>(
        'DescribeMaintenanceWindowExecutionTaskInvocations',
        describeMaintenanceWindowExecutionTaskInvocations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowExecutionTaskInvocationsRequest
                .fromBuffer(value),
        ($0.DescribeMaintenanceWindowExecutionTaskInvocationsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeMaintenanceWindowExecutionTasksRequest,
            $0.DescribeMaintenanceWindowExecutionTasksResult>(
        'DescribeMaintenanceWindowExecutionTasks',
        describeMaintenanceWindowExecutionTasks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowExecutionTasksRequest.fromBuffer(value),
        ($0.DescribeMaintenanceWindowExecutionTasksResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeMaintenanceWindowsRequest,
            $0.DescribeMaintenanceWindowsResult>(
        'DescribeMaintenanceWindows',
        describeMaintenanceWindows_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowsRequest.fromBuffer(value),
        ($0.DescribeMaintenanceWindowsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeMaintenanceWindowScheduleRequest,
            $0.DescribeMaintenanceWindowScheduleResult>(
        'DescribeMaintenanceWindowSchedule',
        describeMaintenanceWindowSchedule_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowScheduleRequest.fromBuffer(value),
        ($0.DescribeMaintenanceWindowScheduleResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeMaintenanceWindowsForTargetRequest,
            $0.DescribeMaintenanceWindowsForTargetResult>(
        'DescribeMaintenanceWindowsForTarget',
        describeMaintenanceWindowsForTarget_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowsForTargetRequest.fromBuffer(value),
        ($0.DescribeMaintenanceWindowsForTargetResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeMaintenanceWindowTargetsRequest,
            $0.DescribeMaintenanceWindowTargetsResult>(
        'DescribeMaintenanceWindowTargets',
        describeMaintenanceWindowTargets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowTargetsRequest.fromBuffer(value),
        ($0.DescribeMaintenanceWindowTargetsResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeMaintenanceWindowTasksRequest,
            $0.DescribeMaintenanceWindowTasksResult>(
        'DescribeMaintenanceWindowTasks',
        describeMaintenanceWindowTasks_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeMaintenanceWindowTasksRequest.fromBuffer(value),
        ($0.DescribeMaintenanceWindowTasksResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeOpsItemsRequest,
            $0.DescribeOpsItemsResponse>(
        'DescribeOpsItems',
        describeOpsItems_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeOpsItemsRequest.fromBuffer(value),
        ($0.DescribeOpsItemsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeParametersRequest,
            $0.DescribeParametersResult>(
        'DescribeParameters',
        describeParameters_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeParametersRequest.fromBuffer(value),
        ($0.DescribeParametersResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribePatchBaselinesRequest,
            $0.DescribePatchBaselinesResult>(
        'DescribePatchBaselines',
        describePatchBaselines_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribePatchBaselinesRequest.fromBuffer(value),
        ($0.DescribePatchBaselinesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribePatchGroupsRequest,
            $0.DescribePatchGroupsResult>(
        'DescribePatchGroups',
        describePatchGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribePatchGroupsRequest.fromBuffer(value),
        ($0.DescribePatchGroupsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribePatchGroupStateRequest,
            $0.DescribePatchGroupStateResult>(
        'DescribePatchGroupState',
        describePatchGroupState_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribePatchGroupStateRequest.fromBuffer(value),
        ($0.DescribePatchGroupStateResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribePatchPropertiesRequest,
            $0.DescribePatchPropertiesResult>(
        'DescribePatchProperties',
        describePatchProperties_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribePatchPropertiesRequest.fromBuffer(value),
        ($0.DescribePatchPropertiesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeSessionsRequest,
            $0.DescribeSessionsResponse>(
        'DescribeSessions',
        describeSessions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeSessionsRequest.fromBuffer(value),
        ($0.DescribeSessionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisassociateOpsItemRelatedItemRequest,
            $0.DisassociateOpsItemRelatedItemResponse>(
        'DisassociateOpsItemRelatedItem',
        disassociateOpsItemRelatedItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisassociateOpsItemRelatedItemRequest.fromBuffer(value),
        ($0.DisassociateOpsItemRelatedItemResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccessTokenRequest,
            $0.GetAccessTokenResponse>(
        'GetAccessToken',
        getAccessToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAccessTokenRequest.fromBuffer(value),
        ($0.GetAccessTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAutomationExecutionRequest,
            $0.GetAutomationExecutionResult>(
        'GetAutomationExecution',
        getAutomationExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAutomationExecutionRequest.fromBuffer(value),
        ($0.GetAutomationExecutionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCalendarStateRequest,
            $0.GetCalendarStateResponse>(
        'GetCalendarState',
        getCalendarState_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCalendarStateRequest.fromBuffer(value),
        ($0.GetCalendarStateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCommandInvocationRequest,
            $0.GetCommandInvocationResult>(
        'GetCommandInvocation',
        getCommandInvocation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCommandInvocationRequest.fromBuffer(value),
        ($0.GetCommandInvocationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetConnectionStatusRequest,
            $0.GetConnectionStatusResponse>(
        'GetConnectionStatus',
        getConnectionStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetConnectionStatusRequest.fromBuffer(value),
        ($0.GetConnectionStatusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDefaultPatchBaselineRequest,
            $0.GetDefaultPatchBaselineResult>(
        'GetDefaultPatchBaseline',
        getDefaultPatchBaseline_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDefaultPatchBaselineRequest.fromBuffer(value),
        ($0.GetDefaultPatchBaselineResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetDeployablePatchSnapshotForInstanceRequest,
            $0.GetDeployablePatchSnapshotForInstanceResult>(
        'GetDeployablePatchSnapshotForInstance',
        getDeployablePatchSnapshotForInstance_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDeployablePatchSnapshotForInstanceRequest.fromBuffer(value),
        ($0.GetDeployablePatchSnapshotForInstanceResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDocumentRequest, $0.GetDocumentResult>(
        'GetDocument',
        getDocument_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDocumentRequest.fromBuffer(value),
        ($0.GetDocumentResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetExecutionPreviewRequest,
            $0.GetExecutionPreviewResponse>(
        'GetExecutionPreview',
        getExecutionPreview_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetExecutionPreviewRequest.fromBuffer(value),
        ($0.GetExecutionPreviewResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetInventoryRequest, $0.GetInventoryResult>(
            'GetInventory',
            getInventory_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetInventoryRequest.fromBuffer(value),
            ($0.GetInventoryResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetInventorySchemaRequest,
            $0.GetInventorySchemaResult>(
        'GetInventorySchema',
        getInventorySchema_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetInventorySchemaRequest.fromBuffer(value),
        ($0.GetInventorySchemaResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMaintenanceWindowRequest,
            $0.GetMaintenanceWindowResult>(
        'GetMaintenanceWindow',
        getMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMaintenanceWindowRequest.fromBuffer(value),
        ($0.GetMaintenanceWindowResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMaintenanceWindowExecutionRequest,
            $0.GetMaintenanceWindowExecutionResult>(
        'GetMaintenanceWindowExecution',
        getMaintenanceWindowExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMaintenanceWindowExecutionRequest.fromBuffer(value),
        ($0.GetMaintenanceWindowExecutionResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMaintenanceWindowExecutionTaskRequest,
            $0.GetMaintenanceWindowExecutionTaskResult>(
        'GetMaintenanceWindowExecutionTask',
        getMaintenanceWindowExecutionTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMaintenanceWindowExecutionTaskRequest.fromBuffer(value),
        ($0.GetMaintenanceWindowExecutionTaskResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetMaintenanceWindowExecutionTaskInvocationRequest,
            $0.GetMaintenanceWindowExecutionTaskInvocationResult>(
        'GetMaintenanceWindowExecutionTaskInvocation',
        getMaintenanceWindowExecutionTaskInvocation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMaintenanceWindowExecutionTaskInvocationRequest.fromBuffer(
                value),
        ($0.GetMaintenanceWindowExecutionTaskInvocationResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetMaintenanceWindowTaskRequest,
            $0.GetMaintenanceWindowTaskResult>(
        'GetMaintenanceWindowTask',
        getMaintenanceWindowTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetMaintenanceWindowTaskRequest.fromBuffer(value),
        ($0.GetMaintenanceWindowTaskResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOpsItemRequest, $0.GetOpsItemResponse>(
        'GetOpsItem',
        getOpsItem_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetOpsItemRequest.fromBuffer(value),
        ($0.GetOpsItemResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetOpsMetadataRequest, $0.GetOpsMetadataResult>(
            'GetOpsMetadata',
            getOpsMetadata_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetOpsMetadataRequest.fromBuffer(value),
            ($0.GetOpsMetadataResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetOpsSummaryRequest, $0.GetOpsSummaryResult>(
            'GetOpsSummary',
            getOpsSummary_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetOpsSummaryRequest.fromBuffer(value),
            ($0.GetOpsSummaryResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetParameterRequest, $0.GetParameterResult>(
            'GetParameter',
            getParameter_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetParameterRequest.fromBuffer(value),
            ($0.GetParameterResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetParameterHistoryRequest,
            $0.GetParameterHistoryResult>(
        'GetParameterHistory',
        getParameterHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetParameterHistoryRequest.fromBuffer(value),
        ($0.GetParameterHistoryResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetParametersRequest, $0.GetParametersResult>(
            'GetParameters',
            getParameters_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetParametersRequest.fromBuffer(value),
            ($0.GetParametersResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetParametersByPathRequest,
            $0.GetParametersByPathResult>(
        'GetParametersByPath',
        getParametersByPath_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetParametersByPathRequest.fromBuffer(value),
        ($0.GetParametersByPathResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPatchBaselineRequest,
            $0.GetPatchBaselineResult>(
        'GetPatchBaseline',
        getPatchBaseline_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPatchBaselineRequest.fromBuffer(value),
        ($0.GetPatchBaselineResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPatchBaselineForPatchGroupRequest,
            $0.GetPatchBaselineForPatchGroupResult>(
        'GetPatchBaselineForPatchGroup',
        getPatchBaselineForPatchGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPatchBaselineForPatchGroupRequest.fromBuffer(value),
        ($0.GetPatchBaselineForPatchGroupResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetResourcePoliciesRequest,
            $0.GetResourcePoliciesResponse>(
        'GetResourcePolicies',
        getResourcePolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetResourcePoliciesRequest.fromBuffer(value),
        ($0.GetResourcePoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetServiceSettingRequest,
            $0.GetServiceSettingResult>(
        'GetServiceSetting',
        getServiceSetting_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetServiceSettingRequest.fromBuffer(value),
        ($0.GetServiceSettingResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.LabelParameterVersionRequest,
            $0.LabelParameterVersionResult>(
        'LabelParameterVersion',
        labelParameterVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.LabelParameterVersionRequest.fromBuffer(value),
        ($0.LabelParameterVersionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAssociationsRequest,
            $0.ListAssociationsResult>(
        'ListAssociations',
        listAssociations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAssociationsRequest.fromBuffer(value),
        ($0.ListAssociationsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAssociationVersionsRequest,
            $0.ListAssociationVersionsResult>(
        'ListAssociationVersions',
        listAssociationVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAssociationVersionsRequest.fromBuffer(value),
        ($0.ListAssociationVersionsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCommandInvocationsRequest,
            $0.ListCommandInvocationsResult>(
        'ListCommandInvocations',
        listCommandInvocations_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCommandInvocationsRequest.fromBuffer(value),
        ($0.ListCommandInvocationsResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListCommandsRequest, $0.ListCommandsResult>(
            'ListCommands',
            listCommands_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListCommandsRequest.fromBuffer(value),
            ($0.ListCommandsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListComplianceItemsRequest,
            $0.ListComplianceItemsResult>(
        'ListComplianceItems',
        listComplianceItems_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListComplianceItemsRequest.fromBuffer(value),
        ($0.ListComplianceItemsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListComplianceSummariesRequest,
            $0.ListComplianceSummariesResult>(
        'ListComplianceSummaries',
        listComplianceSummaries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListComplianceSummariesRequest.fromBuffer(value),
        ($0.ListComplianceSummariesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDocumentMetadataHistoryRequest,
            $0.ListDocumentMetadataHistoryResponse>(
        'ListDocumentMetadataHistory',
        listDocumentMetadataHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDocumentMetadataHistoryRequest.fromBuffer(value),
        ($0.ListDocumentMetadataHistoryResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListDocumentsRequest, $0.ListDocumentsResult>(
            'ListDocuments',
            listDocuments_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDocumentsRequest.fromBuffer(value),
            ($0.ListDocumentsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDocumentVersionsRequest,
            $0.ListDocumentVersionsResult>(
        'ListDocumentVersions',
        listDocumentVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDocumentVersionsRequest.fromBuffer(value),
        ($0.ListDocumentVersionsResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListInventoryEntriesRequest,
            $0.ListInventoryEntriesResult>(
        'ListInventoryEntries',
        listInventoryEntries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInventoryEntriesRequest.fromBuffer(value),
        ($0.ListInventoryEntriesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListNodesRequest, $0.ListNodesResult>(
        'ListNodes',
        listNodes_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListNodesRequest.fromBuffer(value),
        ($0.ListNodesResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListNodesSummaryRequest,
            $0.ListNodesSummaryResult>(
        'ListNodesSummary',
        listNodesSummary_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListNodesSummaryRequest.fromBuffer(value),
        ($0.ListNodesSummaryResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOpsItemEventsRequest,
            $0.ListOpsItemEventsResponse>(
        'ListOpsItemEvents',
        listOpsItemEvents_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOpsItemEventsRequest.fromBuffer(value),
        ($0.ListOpsItemEventsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOpsItemRelatedItemsRequest,
            $0.ListOpsItemRelatedItemsResponse>(
        'ListOpsItemRelatedItems',
        listOpsItemRelatedItems_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOpsItemRelatedItemsRequest.fromBuffer(value),
        ($0.ListOpsItemRelatedItemsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOpsMetadataRequest,
            $0.ListOpsMetadataResult>(
        'ListOpsMetadata',
        listOpsMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOpsMetadataRequest.fromBuffer(value),
        ($0.ListOpsMetadataResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResourceComplianceSummariesRequest,
            $0.ListResourceComplianceSummariesResult>(
        'ListResourceComplianceSummaries',
        listResourceComplianceSummaries_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResourceComplianceSummariesRequest.fromBuffer(value),
        ($0.ListResourceComplianceSummariesResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResourceDataSyncRequest,
            $0.ListResourceDataSyncResult>(
        'ListResourceDataSync',
        listResourceDataSync_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResourceDataSyncRequest.fromBuffer(value),
        ($0.ListResourceDataSyncResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResult>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ModifyDocumentPermissionRequest,
            $0.ModifyDocumentPermissionResponse>(
        'ModifyDocumentPermission',
        modifyDocumentPermission_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ModifyDocumentPermissionRequest.fromBuffer(value),
        ($0.ModifyDocumentPermissionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutComplianceItemsRequest,
            $0.PutComplianceItemsResult>(
        'PutComplianceItems',
        putComplianceItems_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutComplianceItemsRequest.fromBuffer(value),
        ($0.PutComplianceItemsResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutInventoryRequest, $0.PutInventoryResult>(
            'PutInventory',
            putInventory_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutInventoryRequest.fromBuffer(value),
            ($0.PutInventoryResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutParameterRequest, $0.PutParameterResult>(
            'PutParameter',
            putParameter_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutParameterRequest.fromBuffer(value),
            ($0.PutParameterResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutResourcePolicyRequest,
            $0.PutResourcePolicyResponse>(
        'PutResourcePolicy',
        putResourcePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutResourcePolicyRequest.fromBuffer(value),
        ($0.PutResourcePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterDefaultPatchBaselineRequest,
            $0.RegisterDefaultPatchBaselineResult>(
        'RegisterDefaultPatchBaseline',
        registerDefaultPatchBaseline_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterDefaultPatchBaselineRequest.fromBuffer(value),
        ($0.RegisterDefaultPatchBaselineResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterPatchBaselineForPatchGroupRequest,
            $0.RegisterPatchBaselineForPatchGroupResult>(
        'RegisterPatchBaselineForPatchGroup',
        registerPatchBaselineForPatchGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterPatchBaselineForPatchGroupRequest.fromBuffer(value),
        ($0.RegisterPatchBaselineForPatchGroupResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.RegisterTargetWithMaintenanceWindowRequest,
            $0.RegisterTargetWithMaintenanceWindowResult>(
        'RegisterTargetWithMaintenanceWindow',
        registerTargetWithMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterTargetWithMaintenanceWindowRequest.fromBuffer(value),
        ($0.RegisterTargetWithMaintenanceWindowResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RegisterTaskWithMaintenanceWindowRequest,
            $0.RegisterTaskWithMaintenanceWindowResult>(
        'RegisterTaskWithMaintenanceWindow',
        registerTaskWithMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RegisterTaskWithMaintenanceWindowRequest.fromBuffer(value),
        ($0.RegisterTaskWithMaintenanceWindowResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveTagsFromResourceRequest,
            $0.RemoveTagsFromResourceResult>(
        'RemoveTagsFromResource',
        removeTagsFromResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveTagsFromResourceRequest.fromBuffer(value),
        ($0.RemoveTagsFromResourceResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResetServiceSettingRequest,
            $0.ResetServiceSettingResult>(
        'ResetServiceSetting',
        resetServiceSetting_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResetServiceSettingRequest.fromBuffer(value),
        ($0.ResetServiceSettingResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ResumeSessionRequest, $0.ResumeSessionResponse>(
            'ResumeSession',
            resumeSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ResumeSessionRequest.fromBuffer(value),
            ($0.ResumeSessionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendAutomationSignalRequest,
            $0.SendAutomationSignalResult>(
        'SendAutomationSignal',
        sendAutomationSignal_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendAutomationSignalRequest.fromBuffer(value),
        ($0.SendAutomationSignalResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendCommandRequest, $0.SendCommandResult>(
        'SendCommand',
        sendCommand_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendCommandRequest.fromBuffer(value),
        ($0.SendCommandResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartAccessRequestRequest,
            $0.StartAccessRequestResponse>(
        'StartAccessRequest',
        startAccessRequest_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartAccessRequestRequest.fromBuffer(value),
        ($0.StartAccessRequestResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartAssociationsOnceRequest,
            $0.StartAssociationsOnceResult>(
        'StartAssociationsOnce',
        startAssociationsOnce_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartAssociationsOnceRequest.fromBuffer(value),
        ($0.StartAssociationsOnceResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartAutomationExecutionRequest,
            $0.StartAutomationExecutionResult>(
        'StartAutomationExecution',
        startAutomationExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartAutomationExecutionRequest.fromBuffer(value),
        ($0.StartAutomationExecutionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartChangeRequestExecutionRequest,
            $0.StartChangeRequestExecutionResult>(
        'StartChangeRequestExecution',
        startChangeRequestExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartChangeRequestExecutionRequest.fromBuffer(value),
        ($0.StartChangeRequestExecutionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartExecutionPreviewRequest,
            $0.StartExecutionPreviewResponse>(
        'StartExecutionPreview',
        startExecutionPreview_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartExecutionPreviewRequest.fromBuffer(value),
        ($0.StartExecutionPreviewResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.StartSessionRequest, $0.StartSessionResponse>(
            'StartSession',
            startSession_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.StartSessionRequest.fromBuffer(value),
            ($0.StartSessionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopAutomationExecutionRequest,
            $0.StopAutomationExecutionResult>(
        'StopAutomationExecution',
        stopAutomationExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopAutomationExecutionRequest.fromBuffer(value),
        ($0.StopAutomationExecutionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TerminateSessionRequest,
            $0.TerminateSessionResponse>(
        'TerminateSession',
        terminateSession_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TerminateSessionRequest.fromBuffer(value),
        ($0.TerminateSessionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UnlabelParameterVersionRequest,
            $0.UnlabelParameterVersionResult>(
        'UnlabelParameterVersion',
        unlabelParameterVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UnlabelParameterVersionRequest.fromBuffer(value),
        ($0.UnlabelParameterVersionResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAssociationRequest,
            $0.UpdateAssociationResult>(
        'UpdateAssociation',
        updateAssociation_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAssociationRequest.fromBuffer(value),
        ($0.UpdateAssociationResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAssociationStatusRequest,
            $0.UpdateAssociationStatusResult>(
        'UpdateAssociationStatus',
        updateAssociationStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAssociationStatusRequest.fromBuffer(value),
        ($0.UpdateAssociationStatusResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateDocumentRequest, $0.UpdateDocumentResult>(
            'UpdateDocument',
            updateDocument_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateDocumentRequest.fromBuffer(value),
            ($0.UpdateDocumentResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDocumentDefaultVersionRequest,
            $0.UpdateDocumentDefaultVersionResult>(
        'UpdateDocumentDefaultVersion',
        updateDocumentDefaultVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDocumentDefaultVersionRequest.fromBuffer(value),
        ($0.UpdateDocumentDefaultVersionResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDocumentMetadataRequest,
            $0.UpdateDocumentMetadataResponse>(
        'UpdateDocumentMetadata',
        updateDocumentMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDocumentMetadataRequest.fromBuffer(value),
        ($0.UpdateDocumentMetadataResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateMaintenanceWindowRequest,
            $0.UpdateMaintenanceWindowResult>(
        'UpdateMaintenanceWindow',
        updateMaintenanceWindow_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateMaintenanceWindowRequest.fromBuffer(value),
        ($0.UpdateMaintenanceWindowResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateMaintenanceWindowTargetRequest,
            $0.UpdateMaintenanceWindowTargetResult>(
        'UpdateMaintenanceWindowTarget',
        updateMaintenanceWindowTarget_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateMaintenanceWindowTargetRequest.fromBuffer(value),
        ($0.UpdateMaintenanceWindowTargetResult value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateMaintenanceWindowTaskRequest,
            $0.UpdateMaintenanceWindowTaskResult>(
        'UpdateMaintenanceWindowTask',
        updateMaintenanceWindowTask_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateMaintenanceWindowTaskRequest.fromBuffer(value),
        ($0.UpdateMaintenanceWindowTaskResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateManagedInstanceRoleRequest,
            $0.UpdateManagedInstanceRoleResult>(
        'UpdateManagedInstanceRole',
        updateManagedInstanceRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateManagedInstanceRoleRequest.fromBuffer(value),
        ($0.UpdateManagedInstanceRoleResult value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateOpsItemRequest, $0.UpdateOpsItemResponse>(
            'UpdateOpsItem',
            updateOpsItem_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateOpsItemRequest.fromBuffer(value),
            ($0.UpdateOpsItemResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateOpsMetadataRequest,
            $0.UpdateOpsMetadataResult>(
        'UpdateOpsMetadata',
        updateOpsMetadata_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateOpsMetadataRequest.fromBuffer(value),
        ($0.UpdateOpsMetadataResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdatePatchBaselineRequest,
            $0.UpdatePatchBaselineResult>(
        'UpdatePatchBaseline',
        updatePatchBaseline_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdatePatchBaselineRequest.fromBuffer(value),
        ($0.UpdatePatchBaselineResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateResourceDataSyncRequest,
            $0.UpdateResourceDataSyncResult>(
        'UpdateResourceDataSync',
        updateResourceDataSync_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateResourceDataSyncRequest.fromBuffer(value),
        ($0.UpdateResourceDataSyncResult value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateServiceSettingRequest,
            $0.UpdateServiceSettingResult>(
        'UpdateServiceSetting',
        updateServiceSetting_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateServiceSettingRequest.fromBuffer(value),
        ($0.UpdateServiceSettingResult value) => value.writeToBuffer()));
  }

  $async.Future<$0.AddTagsToResourceResult> addTagsToResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AddTagsToResourceRequest> $request) async {
    return addTagsToResource($call, await $request);
  }

  $async.Future<$0.AddTagsToResourceResult> addTagsToResource(
      $grpc.ServiceCall call, $0.AddTagsToResourceRequest request);

  $async.Future<$0.AssociateOpsItemRelatedItemResponse>
      associateOpsItemRelatedItem_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AssociateOpsItemRelatedItemRequest> $request) async {
    return associateOpsItemRelatedItem($call, await $request);
  }

  $async.Future<$0.AssociateOpsItemRelatedItemResponse>
      associateOpsItemRelatedItem($grpc.ServiceCall call,
          $0.AssociateOpsItemRelatedItemRequest request);

  $async.Future<$0.CancelCommandResult> cancelCommand_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CancelCommandRequest> $request) async {
    return cancelCommand($call, await $request);
  }

  $async.Future<$0.CancelCommandResult> cancelCommand(
      $grpc.ServiceCall call, $0.CancelCommandRequest request);

  $async.Future<$0.CancelMaintenanceWindowExecutionResult>
      cancelMaintenanceWindowExecution_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CancelMaintenanceWindowExecutionRequest>
              $request) async {
    return cancelMaintenanceWindowExecution($call, await $request);
  }

  $async.Future<$0.CancelMaintenanceWindowExecutionResult>
      cancelMaintenanceWindowExecution($grpc.ServiceCall call,
          $0.CancelMaintenanceWindowExecutionRequest request);

  $async.Future<$0.CreateActivationResult> createActivation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateActivationRequest> $request) async {
    return createActivation($call, await $request);
  }

  $async.Future<$0.CreateActivationResult> createActivation(
      $grpc.ServiceCall call, $0.CreateActivationRequest request);

  $async.Future<$0.CreateAssociationResult> createAssociation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateAssociationRequest> $request) async {
    return createAssociation($call, await $request);
  }

  $async.Future<$0.CreateAssociationResult> createAssociation(
      $grpc.ServiceCall call, $0.CreateAssociationRequest request);

  $async.Future<$0.CreateAssociationBatchResult> createAssociationBatch_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateAssociationBatchRequest> $request) async {
    return createAssociationBatch($call, await $request);
  }

  $async.Future<$0.CreateAssociationBatchResult> createAssociationBatch(
      $grpc.ServiceCall call, $0.CreateAssociationBatchRequest request);

  $async.Future<$0.CreateDocumentResult> createDocument_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDocumentRequest> $request) async {
    return createDocument($call, await $request);
  }

  $async.Future<$0.CreateDocumentResult> createDocument(
      $grpc.ServiceCall call, $0.CreateDocumentRequest request);

  $async.Future<$0.CreateMaintenanceWindowResult> createMaintenanceWindow_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateMaintenanceWindowRequest> $request) async {
    return createMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.CreateMaintenanceWindowResult> createMaintenanceWindow(
      $grpc.ServiceCall call, $0.CreateMaintenanceWindowRequest request);

  $async.Future<$0.CreateOpsItemResponse> createOpsItem_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateOpsItemRequest> $request) async {
    return createOpsItem($call, await $request);
  }

  $async.Future<$0.CreateOpsItemResponse> createOpsItem(
      $grpc.ServiceCall call, $0.CreateOpsItemRequest request);

  $async.Future<$0.CreateOpsMetadataResult> createOpsMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateOpsMetadataRequest> $request) async {
    return createOpsMetadata($call, await $request);
  }

  $async.Future<$0.CreateOpsMetadataResult> createOpsMetadata(
      $grpc.ServiceCall call, $0.CreateOpsMetadataRequest request);

  $async.Future<$0.CreatePatchBaselineResult> createPatchBaseline_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreatePatchBaselineRequest> $request) async {
    return createPatchBaseline($call, await $request);
  }

  $async.Future<$0.CreatePatchBaselineResult> createPatchBaseline(
      $grpc.ServiceCall call, $0.CreatePatchBaselineRequest request);

  $async.Future<$0.CreateResourceDataSyncResult> createResourceDataSync_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateResourceDataSyncRequest> $request) async {
    return createResourceDataSync($call, await $request);
  }

  $async.Future<$0.CreateResourceDataSyncResult> createResourceDataSync(
      $grpc.ServiceCall call, $0.CreateResourceDataSyncRequest request);

  $async.Future<$0.DeleteActivationResult> deleteActivation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteActivationRequest> $request) async {
    return deleteActivation($call, await $request);
  }

  $async.Future<$0.DeleteActivationResult> deleteActivation(
      $grpc.ServiceCall call, $0.DeleteActivationRequest request);

  $async.Future<$0.DeleteAssociationResult> deleteAssociation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteAssociationRequest> $request) async {
    return deleteAssociation($call, await $request);
  }

  $async.Future<$0.DeleteAssociationResult> deleteAssociation(
      $grpc.ServiceCall call, $0.DeleteAssociationRequest request);

  $async.Future<$0.DeleteDocumentResult> deleteDocument_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteDocumentRequest> $request) async {
    return deleteDocument($call, await $request);
  }

  $async.Future<$0.DeleteDocumentResult> deleteDocument(
      $grpc.ServiceCall call, $0.DeleteDocumentRequest request);

  $async.Future<$0.DeleteInventoryResult> deleteInventory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteInventoryRequest> $request) async {
    return deleteInventory($call, await $request);
  }

  $async.Future<$0.DeleteInventoryResult> deleteInventory(
      $grpc.ServiceCall call, $0.DeleteInventoryRequest request);

  $async.Future<$0.DeleteMaintenanceWindowResult> deleteMaintenanceWindow_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteMaintenanceWindowRequest> $request) async {
    return deleteMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.DeleteMaintenanceWindowResult> deleteMaintenanceWindow(
      $grpc.ServiceCall call, $0.DeleteMaintenanceWindowRequest request);

  $async.Future<$0.DeleteOpsItemResponse> deleteOpsItem_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteOpsItemRequest> $request) async {
    return deleteOpsItem($call, await $request);
  }

  $async.Future<$0.DeleteOpsItemResponse> deleteOpsItem(
      $grpc.ServiceCall call, $0.DeleteOpsItemRequest request);

  $async.Future<$0.DeleteOpsMetadataResult> deleteOpsMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteOpsMetadataRequest> $request) async {
    return deleteOpsMetadata($call, await $request);
  }

  $async.Future<$0.DeleteOpsMetadataResult> deleteOpsMetadata(
      $grpc.ServiceCall call, $0.DeleteOpsMetadataRequest request);

  $async.Future<$0.DeleteParameterResult> deleteParameter_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteParameterRequest> $request) async {
    return deleteParameter($call, await $request);
  }

  $async.Future<$0.DeleteParameterResult> deleteParameter(
      $grpc.ServiceCall call, $0.DeleteParameterRequest request);

  $async.Future<$0.DeleteParametersResult> deleteParameters_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteParametersRequest> $request) async {
    return deleteParameters($call, await $request);
  }

  $async.Future<$0.DeleteParametersResult> deleteParameters(
      $grpc.ServiceCall call, $0.DeleteParametersRequest request);

  $async.Future<$0.DeletePatchBaselineResult> deletePatchBaseline_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeletePatchBaselineRequest> $request) async {
    return deletePatchBaseline($call, await $request);
  }

  $async.Future<$0.DeletePatchBaselineResult> deletePatchBaseline(
      $grpc.ServiceCall call, $0.DeletePatchBaselineRequest request);

  $async.Future<$0.DeleteResourceDataSyncResult> deleteResourceDataSync_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourceDataSyncRequest> $request) async {
    return deleteResourceDataSync($call, await $request);
  }

  $async.Future<$0.DeleteResourceDataSyncResult> deleteResourceDataSync(
      $grpc.ServiceCall call, $0.DeleteResourceDataSyncRequest request);

  $async.Future<$0.DeleteResourcePolicyResponse> deleteResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourcePolicyRequest> $request) async {
    return deleteResourcePolicy($call, await $request);
  }

  $async.Future<$0.DeleteResourcePolicyResponse> deleteResourcePolicy(
      $grpc.ServiceCall call, $0.DeleteResourcePolicyRequest request);

  $async.Future<$0.DeregisterManagedInstanceResult>
      deregisterManagedInstance_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeregisterManagedInstanceRequest> $request) async {
    return deregisterManagedInstance($call, await $request);
  }

  $async.Future<$0.DeregisterManagedInstanceResult> deregisterManagedInstance(
      $grpc.ServiceCall call, $0.DeregisterManagedInstanceRequest request);

  $async.Future<$0.DeregisterPatchBaselineForPatchGroupResult>
      deregisterPatchBaselineForPatchGroup_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeregisterPatchBaselineForPatchGroupRequest>
              $request) async {
    return deregisterPatchBaselineForPatchGroup($call, await $request);
  }

  $async.Future<$0.DeregisterPatchBaselineForPatchGroupResult>
      deregisterPatchBaselineForPatchGroup($grpc.ServiceCall call,
          $0.DeregisterPatchBaselineForPatchGroupRequest request);

  $async.Future<$0.DeregisterTargetFromMaintenanceWindowResult>
      deregisterTargetFromMaintenanceWindow_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeregisterTargetFromMaintenanceWindowRequest>
              $request) async {
    return deregisterTargetFromMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.DeregisterTargetFromMaintenanceWindowResult>
      deregisterTargetFromMaintenanceWindow($grpc.ServiceCall call,
          $0.DeregisterTargetFromMaintenanceWindowRequest request);

  $async.Future<$0.DeregisterTaskFromMaintenanceWindowResult>
      deregisterTaskFromMaintenanceWindow_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DeregisterTaskFromMaintenanceWindowRequest>
              $request) async {
    return deregisterTaskFromMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.DeregisterTaskFromMaintenanceWindowResult>
      deregisterTaskFromMaintenanceWindow($grpc.ServiceCall call,
          $0.DeregisterTaskFromMaintenanceWindowRequest request);

  $async.Future<$0.DescribeActivationsResult> describeActivations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeActivationsRequest> $request) async {
    return describeActivations($call, await $request);
  }

  $async.Future<$0.DescribeActivationsResult> describeActivations(
      $grpc.ServiceCall call, $0.DescribeActivationsRequest request);

  $async.Future<$0.DescribeAssociationResult> describeAssociation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAssociationRequest> $request) async {
    return describeAssociation($call, await $request);
  }

  $async.Future<$0.DescribeAssociationResult> describeAssociation(
      $grpc.ServiceCall call, $0.DescribeAssociationRequest request);

  $async.Future<$0.DescribeAssociationExecutionsResult>
      describeAssociationExecutions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeAssociationExecutionsRequest>
              $request) async {
    return describeAssociationExecutions($call, await $request);
  }

  $async.Future<$0.DescribeAssociationExecutionsResult>
      describeAssociationExecutions($grpc.ServiceCall call,
          $0.DescribeAssociationExecutionsRequest request);

  $async.Future<$0.DescribeAssociationExecutionTargetsResult>
      describeAssociationExecutionTargets_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeAssociationExecutionTargetsRequest>
              $request) async {
    return describeAssociationExecutionTargets($call, await $request);
  }

  $async.Future<$0.DescribeAssociationExecutionTargetsResult>
      describeAssociationExecutionTargets($grpc.ServiceCall call,
          $0.DescribeAssociationExecutionTargetsRequest request);

  $async.Future<$0.DescribeAutomationExecutionsResult>
      describeAutomationExecutions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeAutomationExecutionsRequest>
              $request) async {
    return describeAutomationExecutions($call, await $request);
  }

  $async.Future<$0.DescribeAutomationExecutionsResult>
      describeAutomationExecutions($grpc.ServiceCall call,
          $0.DescribeAutomationExecutionsRequest request);

  $async.Future<$0.DescribeAutomationStepExecutionsResult>
      describeAutomationStepExecutions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeAutomationStepExecutionsRequest>
              $request) async {
    return describeAutomationStepExecutions($call, await $request);
  }

  $async.Future<$0.DescribeAutomationStepExecutionsResult>
      describeAutomationStepExecutions($grpc.ServiceCall call,
          $0.DescribeAutomationStepExecutionsRequest request);

  $async.Future<$0.DescribeAvailablePatchesResult> describeAvailablePatches_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeAvailablePatchesRequest> $request) async {
    return describeAvailablePatches($call, await $request);
  }

  $async.Future<$0.DescribeAvailablePatchesResult> describeAvailablePatches(
      $grpc.ServiceCall call, $0.DescribeAvailablePatchesRequest request);

  $async.Future<$0.DescribeDocumentResult> describeDocument_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeDocumentRequest> $request) async {
    return describeDocument($call, await $request);
  }

  $async.Future<$0.DescribeDocumentResult> describeDocument(
      $grpc.ServiceCall call, $0.DescribeDocumentRequest request);

  $async.Future<$0.DescribeDocumentPermissionResponse>
      describeDocumentPermission_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeDocumentPermissionRequest> $request) async {
    return describeDocumentPermission($call, await $request);
  }

  $async.Future<$0.DescribeDocumentPermissionResponse>
      describeDocumentPermission(
          $grpc.ServiceCall call, $0.DescribeDocumentPermissionRequest request);

  $async.Future<$0.DescribeEffectiveInstanceAssociationsResult>
      describeEffectiveInstanceAssociations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeEffectiveInstanceAssociationsRequest>
              $request) async {
    return describeEffectiveInstanceAssociations($call, await $request);
  }

  $async.Future<$0.DescribeEffectiveInstanceAssociationsResult>
      describeEffectiveInstanceAssociations($grpc.ServiceCall call,
          $0.DescribeEffectiveInstanceAssociationsRequest request);

  $async.Future<$0.DescribeEffectivePatchesForPatchBaselineResult>
      describeEffectivePatchesForPatchBaseline_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeEffectivePatchesForPatchBaselineRequest>
              $request) async {
    return describeEffectivePatchesForPatchBaseline($call, await $request);
  }

  $async.Future<$0.DescribeEffectivePatchesForPatchBaselineResult>
      describeEffectivePatchesForPatchBaseline($grpc.ServiceCall call,
          $0.DescribeEffectivePatchesForPatchBaselineRequest request);

  $async.Future<$0.DescribeInstanceAssociationsStatusResult>
      describeInstanceAssociationsStatus_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeInstanceAssociationsStatusRequest>
              $request) async {
    return describeInstanceAssociationsStatus($call, await $request);
  }

  $async.Future<$0.DescribeInstanceAssociationsStatusResult>
      describeInstanceAssociationsStatus($grpc.ServiceCall call,
          $0.DescribeInstanceAssociationsStatusRequest request);

  $async.Future<$0.DescribeInstanceInformationResult>
      describeInstanceInformation_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeInstanceInformationRequest> $request) async {
    return describeInstanceInformation($call, await $request);
  }

  $async.Future<$0.DescribeInstanceInformationResult>
      describeInstanceInformation($grpc.ServiceCall call,
          $0.DescribeInstanceInformationRequest request);

  $async.Future<$0.DescribeInstancePatchesResult> describeInstancePatches_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeInstancePatchesRequest> $request) async {
    return describeInstancePatches($call, await $request);
  }

  $async.Future<$0.DescribeInstancePatchesResult> describeInstancePatches(
      $grpc.ServiceCall call, $0.DescribeInstancePatchesRequest request);

  $async.Future<$0.DescribeInstancePatchStatesResult>
      describeInstancePatchStates_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeInstancePatchStatesRequest> $request) async {
    return describeInstancePatchStates($call, await $request);
  }

  $async.Future<$0.DescribeInstancePatchStatesResult>
      describeInstancePatchStates($grpc.ServiceCall call,
          $0.DescribeInstancePatchStatesRequest request);

  $async.Future<$0.DescribeInstancePatchStatesForPatchGroupResult>
      describeInstancePatchStatesForPatchGroup_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeInstancePatchStatesForPatchGroupRequest>
              $request) async {
    return describeInstancePatchStatesForPatchGroup($call, await $request);
  }

  $async.Future<$0.DescribeInstancePatchStatesForPatchGroupResult>
      describeInstancePatchStatesForPatchGroup($grpc.ServiceCall call,
          $0.DescribeInstancePatchStatesForPatchGroupRequest request);

  $async.Future<$0.DescribeInstancePropertiesResult>
      describeInstanceProperties_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeInstancePropertiesRequest> $request) async {
    return describeInstanceProperties($call, await $request);
  }

  $async.Future<$0.DescribeInstancePropertiesResult> describeInstanceProperties(
      $grpc.ServiceCall call, $0.DescribeInstancePropertiesRequest request);

  $async.Future<$0.DescribeInventoryDeletionsResult>
      describeInventoryDeletions_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeInventoryDeletionsRequest> $request) async {
    return describeInventoryDeletions($call, await $request);
  }

  $async.Future<$0.DescribeInventoryDeletionsResult> describeInventoryDeletions(
      $grpc.ServiceCall call, $0.DescribeInventoryDeletionsRequest request);

  $async.Future<$0.DescribeMaintenanceWindowExecutionsResult>
      describeMaintenanceWindowExecutions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeMaintenanceWindowExecutionsRequest>
              $request) async {
    return describeMaintenanceWindowExecutions($call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowExecutionsResult>
      describeMaintenanceWindowExecutions($grpc.ServiceCall call,
          $0.DescribeMaintenanceWindowExecutionsRequest request);

  $async.Future<$0.DescribeMaintenanceWindowExecutionTaskInvocationsResult>
      describeMaintenanceWindowExecutionTaskInvocations_Pre(
          $grpc.ServiceCall $call,
          $async.Future<
                  $0.DescribeMaintenanceWindowExecutionTaskInvocationsRequest>
              $request) async {
    return describeMaintenanceWindowExecutionTaskInvocations(
        $call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowExecutionTaskInvocationsResult>
      describeMaintenanceWindowExecutionTaskInvocations($grpc.ServiceCall call,
          $0.DescribeMaintenanceWindowExecutionTaskInvocationsRequest request);

  $async.Future<$0.DescribeMaintenanceWindowExecutionTasksResult>
      describeMaintenanceWindowExecutionTasks_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeMaintenanceWindowExecutionTasksRequest>
              $request) async {
    return describeMaintenanceWindowExecutionTasks($call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowExecutionTasksResult>
      describeMaintenanceWindowExecutionTasks($grpc.ServiceCall call,
          $0.DescribeMaintenanceWindowExecutionTasksRequest request);

  $async.Future<$0.DescribeMaintenanceWindowsResult>
      describeMaintenanceWindows_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeMaintenanceWindowsRequest> $request) async {
    return describeMaintenanceWindows($call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowsResult> describeMaintenanceWindows(
      $grpc.ServiceCall call, $0.DescribeMaintenanceWindowsRequest request);

  $async.Future<$0.DescribeMaintenanceWindowScheduleResult>
      describeMaintenanceWindowSchedule_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeMaintenanceWindowScheduleRequest>
              $request) async {
    return describeMaintenanceWindowSchedule($call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowScheduleResult>
      describeMaintenanceWindowSchedule($grpc.ServiceCall call,
          $0.DescribeMaintenanceWindowScheduleRequest request);

  $async.Future<$0.DescribeMaintenanceWindowsForTargetResult>
      describeMaintenanceWindowsForTarget_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeMaintenanceWindowsForTargetRequest>
              $request) async {
    return describeMaintenanceWindowsForTarget($call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowsForTargetResult>
      describeMaintenanceWindowsForTarget($grpc.ServiceCall call,
          $0.DescribeMaintenanceWindowsForTargetRequest request);

  $async.Future<$0.DescribeMaintenanceWindowTargetsResult>
      describeMaintenanceWindowTargets_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeMaintenanceWindowTargetsRequest>
              $request) async {
    return describeMaintenanceWindowTargets($call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowTargetsResult>
      describeMaintenanceWindowTargets($grpc.ServiceCall call,
          $0.DescribeMaintenanceWindowTargetsRequest request);

  $async.Future<$0.DescribeMaintenanceWindowTasksResult>
      describeMaintenanceWindowTasks_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeMaintenanceWindowTasksRequest>
              $request) async {
    return describeMaintenanceWindowTasks($call, await $request);
  }

  $async.Future<$0.DescribeMaintenanceWindowTasksResult>
      describeMaintenanceWindowTasks($grpc.ServiceCall call,
          $0.DescribeMaintenanceWindowTasksRequest request);

  $async.Future<$0.DescribeOpsItemsResponse> describeOpsItems_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeOpsItemsRequest> $request) async {
    return describeOpsItems($call, await $request);
  }

  $async.Future<$0.DescribeOpsItemsResponse> describeOpsItems(
      $grpc.ServiceCall call, $0.DescribeOpsItemsRequest request);

  $async.Future<$0.DescribeParametersResult> describeParameters_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeParametersRequest> $request) async {
    return describeParameters($call, await $request);
  }

  $async.Future<$0.DescribeParametersResult> describeParameters(
      $grpc.ServiceCall call, $0.DescribeParametersRequest request);

  $async.Future<$0.DescribePatchBaselinesResult> describePatchBaselines_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribePatchBaselinesRequest> $request) async {
    return describePatchBaselines($call, await $request);
  }

  $async.Future<$0.DescribePatchBaselinesResult> describePatchBaselines(
      $grpc.ServiceCall call, $0.DescribePatchBaselinesRequest request);

  $async.Future<$0.DescribePatchGroupsResult> describePatchGroups_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribePatchGroupsRequest> $request) async {
    return describePatchGroups($call, await $request);
  }

  $async.Future<$0.DescribePatchGroupsResult> describePatchGroups(
      $grpc.ServiceCall call, $0.DescribePatchGroupsRequest request);

  $async.Future<$0.DescribePatchGroupStateResult> describePatchGroupState_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribePatchGroupStateRequest> $request) async {
    return describePatchGroupState($call, await $request);
  }

  $async.Future<$0.DescribePatchGroupStateResult> describePatchGroupState(
      $grpc.ServiceCall call, $0.DescribePatchGroupStateRequest request);

  $async.Future<$0.DescribePatchPropertiesResult> describePatchProperties_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribePatchPropertiesRequest> $request) async {
    return describePatchProperties($call, await $request);
  }

  $async.Future<$0.DescribePatchPropertiesResult> describePatchProperties(
      $grpc.ServiceCall call, $0.DescribePatchPropertiesRequest request);

  $async.Future<$0.DescribeSessionsResponse> describeSessions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeSessionsRequest> $request) async {
    return describeSessions($call, await $request);
  }

  $async.Future<$0.DescribeSessionsResponse> describeSessions(
      $grpc.ServiceCall call, $0.DescribeSessionsRequest request);

  $async.Future<$0.DisassociateOpsItemRelatedItemResponse>
      disassociateOpsItemRelatedItem_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DisassociateOpsItemRelatedItemRequest>
              $request) async {
    return disassociateOpsItemRelatedItem($call, await $request);
  }

  $async.Future<$0.DisassociateOpsItemRelatedItemResponse>
      disassociateOpsItemRelatedItem($grpc.ServiceCall call,
          $0.DisassociateOpsItemRelatedItemRequest request);

  $async.Future<$0.GetAccessTokenResponse> getAccessToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAccessTokenRequest> $request) async {
    return getAccessToken($call, await $request);
  }

  $async.Future<$0.GetAccessTokenResponse> getAccessToken(
      $grpc.ServiceCall call, $0.GetAccessTokenRequest request);

  $async.Future<$0.GetAutomationExecutionResult> getAutomationExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAutomationExecutionRequest> $request) async {
    return getAutomationExecution($call, await $request);
  }

  $async.Future<$0.GetAutomationExecutionResult> getAutomationExecution(
      $grpc.ServiceCall call, $0.GetAutomationExecutionRequest request);

  $async.Future<$0.GetCalendarStateResponse> getCalendarState_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCalendarStateRequest> $request) async {
    return getCalendarState($call, await $request);
  }

  $async.Future<$0.GetCalendarStateResponse> getCalendarState(
      $grpc.ServiceCall call, $0.GetCalendarStateRequest request);

  $async.Future<$0.GetCommandInvocationResult> getCommandInvocation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCommandInvocationRequest> $request) async {
    return getCommandInvocation($call, await $request);
  }

  $async.Future<$0.GetCommandInvocationResult> getCommandInvocation(
      $grpc.ServiceCall call, $0.GetCommandInvocationRequest request);

  $async.Future<$0.GetConnectionStatusResponse> getConnectionStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetConnectionStatusRequest> $request) async {
    return getConnectionStatus($call, await $request);
  }

  $async.Future<$0.GetConnectionStatusResponse> getConnectionStatus(
      $grpc.ServiceCall call, $0.GetConnectionStatusRequest request);

  $async.Future<$0.GetDefaultPatchBaselineResult> getDefaultPatchBaseline_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDefaultPatchBaselineRequest> $request) async {
    return getDefaultPatchBaseline($call, await $request);
  }

  $async.Future<$0.GetDefaultPatchBaselineResult> getDefaultPatchBaseline(
      $grpc.ServiceCall call, $0.GetDefaultPatchBaselineRequest request);

  $async.Future<$0.GetDeployablePatchSnapshotForInstanceResult>
      getDeployablePatchSnapshotForInstance_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetDeployablePatchSnapshotForInstanceRequest>
              $request) async {
    return getDeployablePatchSnapshotForInstance($call, await $request);
  }

  $async.Future<$0.GetDeployablePatchSnapshotForInstanceResult>
      getDeployablePatchSnapshotForInstance($grpc.ServiceCall call,
          $0.GetDeployablePatchSnapshotForInstanceRequest request);

  $async.Future<$0.GetDocumentResult> getDocument_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDocumentRequest> $request) async {
    return getDocument($call, await $request);
  }

  $async.Future<$0.GetDocumentResult> getDocument(
      $grpc.ServiceCall call, $0.GetDocumentRequest request);

  $async.Future<$0.GetExecutionPreviewResponse> getExecutionPreview_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetExecutionPreviewRequest> $request) async {
    return getExecutionPreview($call, await $request);
  }

  $async.Future<$0.GetExecutionPreviewResponse> getExecutionPreview(
      $grpc.ServiceCall call, $0.GetExecutionPreviewRequest request);

  $async.Future<$0.GetInventoryResult> getInventory_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetInventoryRequest> $request) async {
    return getInventory($call, await $request);
  }

  $async.Future<$0.GetInventoryResult> getInventory(
      $grpc.ServiceCall call, $0.GetInventoryRequest request);

  $async.Future<$0.GetInventorySchemaResult> getInventorySchema_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetInventorySchemaRequest> $request) async {
    return getInventorySchema($call, await $request);
  }

  $async.Future<$0.GetInventorySchemaResult> getInventorySchema(
      $grpc.ServiceCall call, $0.GetInventorySchemaRequest request);

  $async.Future<$0.GetMaintenanceWindowResult> getMaintenanceWindow_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMaintenanceWindowRequest> $request) async {
    return getMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.GetMaintenanceWindowResult> getMaintenanceWindow(
      $grpc.ServiceCall call, $0.GetMaintenanceWindowRequest request);

  $async.Future<$0.GetMaintenanceWindowExecutionResult>
      getMaintenanceWindowExecution_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetMaintenanceWindowExecutionRequest>
              $request) async {
    return getMaintenanceWindowExecution($call, await $request);
  }

  $async.Future<$0.GetMaintenanceWindowExecutionResult>
      getMaintenanceWindowExecution($grpc.ServiceCall call,
          $0.GetMaintenanceWindowExecutionRequest request);

  $async.Future<$0.GetMaintenanceWindowExecutionTaskResult>
      getMaintenanceWindowExecutionTask_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetMaintenanceWindowExecutionTaskRequest>
              $request) async {
    return getMaintenanceWindowExecutionTask($call, await $request);
  }

  $async.Future<$0.GetMaintenanceWindowExecutionTaskResult>
      getMaintenanceWindowExecutionTask($grpc.ServiceCall call,
          $0.GetMaintenanceWindowExecutionTaskRequest request);

  $async.Future<$0.GetMaintenanceWindowExecutionTaskInvocationResult>
      getMaintenanceWindowExecutionTaskInvocation_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetMaintenanceWindowExecutionTaskInvocationRequest>
              $request) async {
    return getMaintenanceWindowExecutionTaskInvocation($call, await $request);
  }

  $async.Future<$0.GetMaintenanceWindowExecutionTaskInvocationResult>
      getMaintenanceWindowExecutionTaskInvocation($grpc.ServiceCall call,
          $0.GetMaintenanceWindowExecutionTaskInvocationRequest request);

  $async.Future<$0.GetMaintenanceWindowTaskResult> getMaintenanceWindowTask_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMaintenanceWindowTaskRequest> $request) async {
    return getMaintenanceWindowTask($call, await $request);
  }

  $async.Future<$0.GetMaintenanceWindowTaskResult> getMaintenanceWindowTask(
      $grpc.ServiceCall call, $0.GetMaintenanceWindowTaskRequest request);

  $async.Future<$0.GetOpsItemResponse> getOpsItem_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetOpsItemRequest> $request) async {
    return getOpsItem($call, await $request);
  }

  $async.Future<$0.GetOpsItemResponse> getOpsItem(
      $grpc.ServiceCall call, $0.GetOpsItemRequest request);

  $async.Future<$0.GetOpsMetadataResult> getOpsMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetOpsMetadataRequest> $request) async {
    return getOpsMetadata($call, await $request);
  }

  $async.Future<$0.GetOpsMetadataResult> getOpsMetadata(
      $grpc.ServiceCall call, $0.GetOpsMetadataRequest request);

  $async.Future<$0.GetOpsSummaryResult> getOpsSummary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetOpsSummaryRequest> $request) async {
    return getOpsSummary($call, await $request);
  }

  $async.Future<$0.GetOpsSummaryResult> getOpsSummary(
      $grpc.ServiceCall call, $0.GetOpsSummaryRequest request);

  $async.Future<$0.GetParameterResult> getParameter_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetParameterRequest> $request) async {
    return getParameter($call, await $request);
  }

  $async.Future<$0.GetParameterResult> getParameter(
      $grpc.ServiceCall call, $0.GetParameterRequest request);

  $async.Future<$0.GetParameterHistoryResult> getParameterHistory_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetParameterHistoryRequest> $request) async {
    return getParameterHistory($call, await $request);
  }

  $async.Future<$0.GetParameterHistoryResult> getParameterHistory(
      $grpc.ServiceCall call, $0.GetParameterHistoryRequest request);

  $async.Future<$0.GetParametersResult> getParameters_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetParametersRequest> $request) async {
    return getParameters($call, await $request);
  }

  $async.Future<$0.GetParametersResult> getParameters(
      $grpc.ServiceCall call, $0.GetParametersRequest request);

  $async.Future<$0.GetParametersByPathResult> getParametersByPath_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetParametersByPathRequest> $request) async {
    return getParametersByPath($call, await $request);
  }

  $async.Future<$0.GetParametersByPathResult> getParametersByPath(
      $grpc.ServiceCall call, $0.GetParametersByPathRequest request);

  $async.Future<$0.GetPatchBaselineResult> getPatchBaseline_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPatchBaselineRequest> $request) async {
    return getPatchBaseline($call, await $request);
  }

  $async.Future<$0.GetPatchBaselineResult> getPatchBaseline(
      $grpc.ServiceCall call, $0.GetPatchBaselineRequest request);

  $async.Future<$0.GetPatchBaselineForPatchGroupResult>
      getPatchBaselineForPatchGroup_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetPatchBaselineForPatchGroupRequest>
              $request) async {
    return getPatchBaselineForPatchGroup($call, await $request);
  }

  $async.Future<$0.GetPatchBaselineForPatchGroupResult>
      getPatchBaselineForPatchGroup($grpc.ServiceCall call,
          $0.GetPatchBaselineForPatchGroupRequest request);

  $async.Future<$0.GetResourcePoliciesResponse> getResourcePolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetResourcePoliciesRequest> $request) async {
    return getResourcePolicies($call, await $request);
  }

  $async.Future<$0.GetResourcePoliciesResponse> getResourcePolicies(
      $grpc.ServiceCall call, $0.GetResourcePoliciesRequest request);

  $async.Future<$0.GetServiceSettingResult> getServiceSetting_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetServiceSettingRequest> $request) async {
    return getServiceSetting($call, await $request);
  }

  $async.Future<$0.GetServiceSettingResult> getServiceSetting(
      $grpc.ServiceCall call, $0.GetServiceSettingRequest request);

  $async.Future<$0.LabelParameterVersionResult> labelParameterVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.LabelParameterVersionRequest> $request) async {
    return labelParameterVersion($call, await $request);
  }

  $async.Future<$0.LabelParameterVersionResult> labelParameterVersion(
      $grpc.ServiceCall call, $0.LabelParameterVersionRequest request);

  $async.Future<$0.ListAssociationsResult> listAssociations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAssociationsRequest> $request) async {
    return listAssociations($call, await $request);
  }

  $async.Future<$0.ListAssociationsResult> listAssociations(
      $grpc.ServiceCall call, $0.ListAssociationsRequest request);

  $async.Future<$0.ListAssociationVersionsResult> listAssociationVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAssociationVersionsRequest> $request) async {
    return listAssociationVersions($call, await $request);
  }

  $async.Future<$0.ListAssociationVersionsResult> listAssociationVersions(
      $grpc.ServiceCall call, $0.ListAssociationVersionsRequest request);

  $async.Future<$0.ListCommandInvocationsResult> listCommandInvocations_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCommandInvocationsRequest> $request) async {
    return listCommandInvocations($call, await $request);
  }

  $async.Future<$0.ListCommandInvocationsResult> listCommandInvocations(
      $grpc.ServiceCall call, $0.ListCommandInvocationsRequest request);

  $async.Future<$0.ListCommandsResult> listCommands_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListCommandsRequest> $request) async {
    return listCommands($call, await $request);
  }

  $async.Future<$0.ListCommandsResult> listCommands(
      $grpc.ServiceCall call, $0.ListCommandsRequest request);

  $async.Future<$0.ListComplianceItemsResult> listComplianceItems_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListComplianceItemsRequest> $request) async {
    return listComplianceItems($call, await $request);
  }

  $async.Future<$0.ListComplianceItemsResult> listComplianceItems(
      $grpc.ServiceCall call, $0.ListComplianceItemsRequest request);

  $async.Future<$0.ListComplianceSummariesResult> listComplianceSummaries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListComplianceSummariesRequest> $request) async {
    return listComplianceSummaries($call, await $request);
  }

  $async.Future<$0.ListComplianceSummariesResult> listComplianceSummaries(
      $grpc.ServiceCall call, $0.ListComplianceSummariesRequest request);

  $async.Future<$0.ListDocumentMetadataHistoryResponse>
      listDocumentMetadataHistory_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListDocumentMetadataHistoryRequest> $request) async {
    return listDocumentMetadataHistory($call, await $request);
  }

  $async.Future<$0.ListDocumentMetadataHistoryResponse>
      listDocumentMetadataHistory($grpc.ServiceCall call,
          $0.ListDocumentMetadataHistoryRequest request);

  $async.Future<$0.ListDocumentsResult> listDocuments_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDocumentsRequest> $request) async {
    return listDocuments($call, await $request);
  }

  $async.Future<$0.ListDocumentsResult> listDocuments(
      $grpc.ServiceCall call, $0.ListDocumentsRequest request);

  $async.Future<$0.ListDocumentVersionsResult> listDocumentVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDocumentVersionsRequest> $request) async {
    return listDocumentVersions($call, await $request);
  }

  $async.Future<$0.ListDocumentVersionsResult> listDocumentVersions(
      $grpc.ServiceCall call, $0.ListDocumentVersionsRequest request);

  $async.Future<$0.ListInventoryEntriesResult> listInventoryEntries_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListInventoryEntriesRequest> $request) async {
    return listInventoryEntries($call, await $request);
  }

  $async.Future<$0.ListInventoryEntriesResult> listInventoryEntries(
      $grpc.ServiceCall call, $0.ListInventoryEntriesRequest request);

  $async.Future<$0.ListNodesResult> listNodes_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListNodesRequest> $request) async {
    return listNodes($call, await $request);
  }

  $async.Future<$0.ListNodesResult> listNodes(
      $grpc.ServiceCall call, $0.ListNodesRequest request);

  $async.Future<$0.ListNodesSummaryResult> listNodesSummary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListNodesSummaryRequest> $request) async {
    return listNodesSummary($call, await $request);
  }

  $async.Future<$0.ListNodesSummaryResult> listNodesSummary(
      $grpc.ServiceCall call, $0.ListNodesSummaryRequest request);

  $async.Future<$0.ListOpsItemEventsResponse> listOpsItemEvents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListOpsItemEventsRequest> $request) async {
    return listOpsItemEvents($call, await $request);
  }

  $async.Future<$0.ListOpsItemEventsResponse> listOpsItemEvents(
      $grpc.ServiceCall call, $0.ListOpsItemEventsRequest request);

  $async.Future<$0.ListOpsItemRelatedItemsResponse> listOpsItemRelatedItems_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListOpsItemRelatedItemsRequest> $request) async {
    return listOpsItemRelatedItems($call, await $request);
  }

  $async.Future<$0.ListOpsItemRelatedItemsResponse> listOpsItemRelatedItems(
      $grpc.ServiceCall call, $0.ListOpsItemRelatedItemsRequest request);

  $async.Future<$0.ListOpsMetadataResult> listOpsMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListOpsMetadataRequest> $request) async {
    return listOpsMetadata($call, await $request);
  }

  $async.Future<$0.ListOpsMetadataResult> listOpsMetadata(
      $grpc.ServiceCall call, $0.ListOpsMetadataRequest request);

  $async.Future<$0.ListResourceComplianceSummariesResult>
      listResourceComplianceSummaries_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListResourceComplianceSummariesRequest>
              $request) async {
    return listResourceComplianceSummaries($call, await $request);
  }

  $async.Future<$0.ListResourceComplianceSummariesResult>
      listResourceComplianceSummaries($grpc.ServiceCall call,
          $0.ListResourceComplianceSummariesRequest request);

  $async.Future<$0.ListResourceDataSyncResult> listResourceDataSync_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListResourceDataSyncRequest> $request) async {
    return listResourceDataSync($call, await $request);
  }

  $async.Future<$0.ListResourceDataSyncResult> listResourceDataSync(
      $grpc.ServiceCall call, $0.ListResourceDataSyncRequest request);

  $async.Future<$0.ListTagsForResourceResult> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResult> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ModifyDocumentPermissionResponse>
      modifyDocumentPermission_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ModifyDocumentPermissionRequest> $request) async {
    return modifyDocumentPermission($call, await $request);
  }

  $async.Future<$0.ModifyDocumentPermissionResponse> modifyDocumentPermission(
      $grpc.ServiceCall call, $0.ModifyDocumentPermissionRequest request);

  $async.Future<$0.PutComplianceItemsResult> putComplianceItems_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutComplianceItemsRequest> $request) async {
    return putComplianceItems($call, await $request);
  }

  $async.Future<$0.PutComplianceItemsResult> putComplianceItems(
      $grpc.ServiceCall call, $0.PutComplianceItemsRequest request);

  $async.Future<$0.PutInventoryResult> putInventory_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutInventoryRequest> $request) async {
    return putInventory($call, await $request);
  }

  $async.Future<$0.PutInventoryResult> putInventory(
      $grpc.ServiceCall call, $0.PutInventoryRequest request);

  $async.Future<$0.PutParameterResult> putParameter_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutParameterRequest> $request) async {
    return putParameter($call, await $request);
  }

  $async.Future<$0.PutParameterResult> putParameter(
      $grpc.ServiceCall call, $0.PutParameterRequest request);

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutResourcePolicyRequest> $request) async {
    return putResourcePolicy($call, await $request);
  }

  $async.Future<$0.PutResourcePolicyResponse> putResourcePolicy(
      $grpc.ServiceCall call, $0.PutResourcePolicyRequest request);

  $async.Future<$0.RegisterDefaultPatchBaselineResult>
      registerDefaultPatchBaseline_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.RegisterDefaultPatchBaselineRequest>
              $request) async {
    return registerDefaultPatchBaseline($call, await $request);
  }

  $async.Future<$0.RegisterDefaultPatchBaselineResult>
      registerDefaultPatchBaseline($grpc.ServiceCall call,
          $0.RegisterDefaultPatchBaselineRequest request);

  $async.Future<$0.RegisterPatchBaselineForPatchGroupResult>
      registerPatchBaselineForPatchGroup_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.RegisterPatchBaselineForPatchGroupRequest>
              $request) async {
    return registerPatchBaselineForPatchGroup($call, await $request);
  }

  $async.Future<$0.RegisterPatchBaselineForPatchGroupResult>
      registerPatchBaselineForPatchGroup($grpc.ServiceCall call,
          $0.RegisterPatchBaselineForPatchGroupRequest request);

  $async.Future<$0.RegisterTargetWithMaintenanceWindowResult>
      registerTargetWithMaintenanceWindow_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.RegisterTargetWithMaintenanceWindowRequest>
              $request) async {
    return registerTargetWithMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.RegisterTargetWithMaintenanceWindowResult>
      registerTargetWithMaintenanceWindow($grpc.ServiceCall call,
          $0.RegisterTargetWithMaintenanceWindowRequest request);

  $async.Future<$0.RegisterTaskWithMaintenanceWindowResult>
      registerTaskWithMaintenanceWindow_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.RegisterTaskWithMaintenanceWindowRequest>
              $request) async {
    return registerTaskWithMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.RegisterTaskWithMaintenanceWindowResult>
      registerTaskWithMaintenanceWindow($grpc.ServiceCall call,
          $0.RegisterTaskWithMaintenanceWindowRequest request);

  $async.Future<$0.RemoveTagsFromResourceResult> removeTagsFromResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RemoveTagsFromResourceRequest> $request) async {
    return removeTagsFromResource($call, await $request);
  }

  $async.Future<$0.RemoveTagsFromResourceResult> removeTagsFromResource(
      $grpc.ServiceCall call, $0.RemoveTagsFromResourceRequest request);

  $async.Future<$0.ResetServiceSettingResult> resetServiceSetting_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ResetServiceSettingRequest> $request) async {
    return resetServiceSetting($call, await $request);
  }

  $async.Future<$0.ResetServiceSettingResult> resetServiceSetting(
      $grpc.ServiceCall call, $0.ResetServiceSettingRequest request);

  $async.Future<$0.ResumeSessionResponse> resumeSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ResumeSessionRequest> $request) async {
    return resumeSession($call, await $request);
  }

  $async.Future<$0.ResumeSessionResponse> resumeSession(
      $grpc.ServiceCall call, $0.ResumeSessionRequest request);

  $async.Future<$0.SendAutomationSignalResult> sendAutomationSignal_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SendAutomationSignalRequest> $request) async {
    return sendAutomationSignal($call, await $request);
  }

  $async.Future<$0.SendAutomationSignalResult> sendAutomationSignal(
      $grpc.ServiceCall call, $0.SendAutomationSignalRequest request);

  $async.Future<$0.SendCommandResult> sendCommand_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendCommandRequest> $request) async {
    return sendCommand($call, await $request);
  }

  $async.Future<$0.SendCommandResult> sendCommand(
      $grpc.ServiceCall call, $0.SendCommandRequest request);

  $async.Future<$0.StartAccessRequestResponse> startAccessRequest_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartAccessRequestRequest> $request) async {
    return startAccessRequest($call, await $request);
  }

  $async.Future<$0.StartAccessRequestResponse> startAccessRequest(
      $grpc.ServiceCall call, $0.StartAccessRequestRequest request);

  $async.Future<$0.StartAssociationsOnceResult> startAssociationsOnce_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartAssociationsOnceRequest> $request) async {
    return startAssociationsOnce($call, await $request);
  }

  $async.Future<$0.StartAssociationsOnceResult> startAssociationsOnce(
      $grpc.ServiceCall call, $0.StartAssociationsOnceRequest request);

  $async.Future<$0.StartAutomationExecutionResult> startAutomationExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartAutomationExecutionRequest> $request) async {
    return startAutomationExecution($call, await $request);
  }

  $async.Future<$0.StartAutomationExecutionResult> startAutomationExecution(
      $grpc.ServiceCall call, $0.StartAutomationExecutionRequest request);

  $async.Future<$0.StartChangeRequestExecutionResult>
      startChangeRequestExecution_Pre($grpc.ServiceCall $call,
          $async.Future<$0.StartChangeRequestExecutionRequest> $request) async {
    return startChangeRequestExecution($call, await $request);
  }

  $async.Future<$0.StartChangeRequestExecutionResult>
      startChangeRequestExecution($grpc.ServiceCall call,
          $0.StartChangeRequestExecutionRequest request);

  $async.Future<$0.StartExecutionPreviewResponse> startExecutionPreview_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartExecutionPreviewRequest> $request) async {
    return startExecutionPreview($call, await $request);
  }

  $async.Future<$0.StartExecutionPreviewResponse> startExecutionPreview(
      $grpc.ServiceCall call, $0.StartExecutionPreviewRequest request);

  $async.Future<$0.StartSessionResponse> startSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartSessionRequest> $request) async {
    return startSession($call, await $request);
  }

  $async.Future<$0.StartSessionResponse> startSession(
      $grpc.ServiceCall call, $0.StartSessionRequest request);

  $async.Future<$0.StopAutomationExecutionResult> stopAutomationExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StopAutomationExecutionRequest> $request) async {
    return stopAutomationExecution($call, await $request);
  }

  $async.Future<$0.StopAutomationExecutionResult> stopAutomationExecution(
      $grpc.ServiceCall call, $0.StopAutomationExecutionRequest request);

  $async.Future<$0.TerminateSessionResponse> terminateSession_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.TerminateSessionRequest> $request) async {
    return terminateSession($call, await $request);
  }

  $async.Future<$0.TerminateSessionResponse> terminateSession(
      $grpc.ServiceCall call, $0.TerminateSessionRequest request);

  $async.Future<$0.UnlabelParameterVersionResult> unlabelParameterVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UnlabelParameterVersionRequest> $request) async {
    return unlabelParameterVersion($call, await $request);
  }

  $async.Future<$0.UnlabelParameterVersionResult> unlabelParameterVersion(
      $grpc.ServiceCall call, $0.UnlabelParameterVersionRequest request);

  $async.Future<$0.UpdateAssociationResult> updateAssociation_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateAssociationRequest> $request) async {
    return updateAssociation($call, await $request);
  }

  $async.Future<$0.UpdateAssociationResult> updateAssociation(
      $grpc.ServiceCall call, $0.UpdateAssociationRequest request);

  $async.Future<$0.UpdateAssociationStatusResult> updateAssociationStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateAssociationStatusRequest> $request) async {
    return updateAssociationStatus($call, await $request);
  }

  $async.Future<$0.UpdateAssociationStatusResult> updateAssociationStatus(
      $grpc.ServiceCall call, $0.UpdateAssociationStatusRequest request);

  $async.Future<$0.UpdateDocumentResult> updateDocument_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDocumentRequest> $request) async {
    return updateDocument($call, await $request);
  }

  $async.Future<$0.UpdateDocumentResult> updateDocument(
      $grpc.ServiceCall call, $0.UpdateDocumentRequest request);

  $async.Future<$0.UpdateDocumentDefaultVersionResult>
      updateDocumentDefaultVersion_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateDocumentDefaultVersionRequest>
              $request) async {
    return updateDocumentDefaultVersion($call, await $request);
  }

  $async.Future<$0.UpdateDocumentDefaultVersionResult>
      updateDocumentDefaultVersion($grpc.ServiceCall call,
          $0.UpdateDocumentDefaultVersionRequest request);

  $async.Future<$0.UpdateDocumentMetadataResponse> updateDocumentMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDocumentMetadataRequest> $request) async {
    return updateDocumentMetadata($call, await $request);
  }

  $async.Future<$0.UpdateDocumentMetadataResponse> updateDocumentMetadata(
      $grpc.ServiceCall call, $0.UpdateDocumentMetadataRequest request);

  $async.Future<$0.UpdateMaintenanceWindowResult> updateMaintenanceWindow_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateMaintenanceWindowRequest> $request) async {
    return updateMaintenanceWindow($call, await $request);
  }

  $async.Future<$0.UpdateMaintenanceWindowResult> updateMaintenanceWindow(
      $grpc.ServiceCall call, $0.UpdateMaintenanceWindowRequest request);

  $async.Future<$0.UpdateMaintenanceWindowTargetResult>
      updateMaintenanceWindowTarget_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateMaintenanceWindowTargetRequest>
              $request) async {
    return updateMaintenanceWindowTarget($call, await $request);
  }

  $async.Future<$0.UpdateMaintenanceWindowTargetResult>
      updateMaintenanceWindowTarget($grpc.ServiceCall call,
          $0.UpdateMaintenanceWindowTargetRequest request);

  $async.Future<$0.UpdateMaintenanceWindowTaskResult>
      updateMaintenanceWindowTask_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateMaintenanceWindowTaskRequest> $request) async {
    return updateMaintenanceWindowTask($call, await $request);
  }

  $async.Future<$0.UpdateMaintenanceWindowTaskResult>
      updateMaintenanceWindowTask($grpc.ServiceCall call,
          $0.UpdateMaintenanceWindowTaskRequest request);

  $async.Future<$0.UpdateManagedInstanceRoleResult>
      updateManagedInstanceRole_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateManagedInstanceRoleRequest> $request) async {
    return updateManagedInstanceRole($call, await $request);
  }

  $async.Future<$0.UpdateManagedInstanceRoleResult> updateManagedInstanceRole(
      $grpc.ServiceCall call, $0.UpdateManagedInstanceRoleRequest request);

  $async.Future<$0.UpdateOpsItemResponse> updateOpsItem_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateOpsItemRequest> $request) async {
    return updateOpsItem($call, await $request);
  }

  $async.Future<$0.UpdateOpsItemResponse> updateOpsItem(
      $grpc.ServiceCall call, $0.UpdateOpsItemRequest request);

  $async.Future<$0.UpdateOpsMetadataResult> updateOpsMetadata_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateOpsMetadataRequest> $request) async {
    return updateOpsMetadata($call, await $request);
  }

  $async.Future<$0.UpdateOpsMetadataResult> updateOpsMetadata(
      $grpc.ServiceCall call, $0.UpdateOpsMetadataRequest request);

  $async.Future<$0.UpdatePatchBaselineResult> updatePatchBaseline_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdatePatchBaselineRequest> $request) async {
    return updatePatchBaseline($call, await $request);
  }

  $async.Future<$0.UpdatePatchBaselineResult> updatePatchBaseline(
      $grpc.ServiceCall call, $0.UpdatePatchBaselineRequest request);

  $async.Future<$0.UpdateResourceDataSyncResult> updateResourceDataSync_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateResourceDataSyncRequest> $request) async {
    return updateResourceDataSync($call, await $request);
  }

  $async.Future<$0.UpdateResourceDataSyncResult> updateResourceDataSync(
      $grpc.ServiceCall call, $0.UpdateResourceDataSyncRequest request);

  $async.Future<$0.UpdateServiceSettingResult> updateServiceSetting_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateServiceSettingRequest> $request) async {
    return updateServiceSetting($call, await $request);
  }

  $async.Future<$0.UpdateServiceSettingResult> updateServiceSetting(
      $grpc.ServiceCall call, $0.UpdateServiceSettingRequest request);
}
