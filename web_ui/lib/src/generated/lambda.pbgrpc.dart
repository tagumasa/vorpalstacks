// This is a generated file - do not edit.
//
// Generated from lambda.proto.

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
import 'lambda.pb.dart' as $0;

export 'lambda.pb.dart';

/// LambdaService provides lambda API operations.
@$pb.GrpcServiceName('lambda.LambdaService')
class LambdaServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  LambdaServiceClient(super.channel, {super.options, super.interceptors});

  /// Saves the progress of a durable function execution during runtime. This API is used by the Lambda durable functions SDK to checkpoint completed steps and schedule asynchronous operations. You typic...
  /// HTTP: POST /2025-12-01/durable-executions/{DurableExecutionArn}/checkpoint
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CheckpointDurableExecutionResponse>
      checkpointDurableExecution(
    $0.CheckpointDurableExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$checkpointDurableExecution, request,
        options: options);
  }

  /// Deletes a Lambda function. To delete a specific function version, use the Qualifier parameter. Otherwise, all versions and aliases are deleted. This doesn't require the user to have explicit permis...
  /// HTTP: DELETE /2015-03-31/functions/{FunctionName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteFunctionResponse> deleteFunction(
    $0.DeleteFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFunction, request, options: options);
  }

  /// Deletes the configuration for asynchronous invocation for a function, version, or alias. To configure options for asynchronous invocation, use PutFunctionEventInvokeConfig.
  /// HTTP: DELETE /2019-09-25/functions/{FunctionName}/event-invoke-config
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteFunctionEventInvokeConfig(
    $0.DeleteFunctionEventInvokeConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFunctionEventInvokeConfig, request,
        options: options);
  }

  /// Retrieves details about your account's limits and usage in an Amazon Web Services Region.
  /// HTTP: GET /2016-08-19/account-settings
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetAccountSettingsResponse> getAccountSettings(
    $0.GetAccountSettingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccountSettings, request, options: options);
  }

  /// Retrieves detailed information about a specific durable execution, including its current status, input payload, result or error information, and execution metadata such as start time and usage stat...
  /// HTTP: GET /2025-12-01/durable-executions/{DurableExecutionArn}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDurableExecutionResponse> getDurableExecution(
    $0.GetDurableExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDurableExecution, request, options: options);
  }

  /// Retrieves the execution history for a durable execution, showing all the steps, callbacks, and events that occurred during the execution. This provides a detailed audit trail of the execution's pro...
  /// HTTP: GET /2025-12-01/durable-executions/{DurableExecutionArn}/history
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDurableExecutionHistoryResponse>
      getDurableExecutionHistory(
    $0.GetDurableExecutionHistoryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDurableExecutionHistory, request,
        options: options);
  }

  /// Retrieves the current execution state required for the replay process during durable function execution. This API is used by the Lambda durable functions SDK to get state information needed for rep...
  /// HTTP: GET /2025-12-01/durable-executions/{DurableExecutionArn}/state
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetDurableExecutionStateResponse>
      getDurableExecutionState(
    $0.GetDurableExecutionStateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDurableExecutionState, request,
        options: options);
  }

  /// Retrieves the configuration for asynchronous invocation for a function, version, or alias. To configure options for asynchronous invocation, use PutFunctionEventInvokeConfig.
  /// HTTP: GET /2019-09-25/functions/{FunctionName}/event-invoke-config
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.FunctionEventInvokeConfig>
      getFunctionEventInvokeConfig(
    $0.GetFunctionEventInvokeConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFunctionEventInvokeConfig, request,
        options: options);
  }

  /// Returns a list of durable executions for a specified Lambda function. You can filter the results by execution name, status, and start time range. This API supports pagination for large result sets.
  /// HTTP: GET /2025-12-01/functions/{FunctionName}/durable-executions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListDurableExecutionsByFunctionResponse>
      listDurableExecutionsByFunction(
    $0.ListDurableExecutionsByFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDurableExecutionsByFunction, request,
        options: options);
  }

  /// Retrieves a list of configurations for asynchronous invocation for a function. To configure options for asynchronous invocation, use PutFunctionEventInvokeConfig.
  /// HTTP: GET /2019-09-25/functions/{FunctionName}/event-invoke-config/list
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListFunctionEventInvokeConfigsResponse>
      listFunctionEventInvokeConfigs(
    $0.ListFunctionEventInvokeConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFunctionEventInvokeConfigs, request,
        options: options);
  }

  /// Returns a function, event source mapping, or code signing configuration's tags. You can also view function tags with GetFunction.
  /// HTTP: GET /2017-03-31/tags/{Resource}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListTagsResponse> listTags(
    $0.ListTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTags, request, options: options);
  }

  /// Configures options for asynchronous invocation on a function, version, or alias. If a configuration already exists for a function, version, or alias, this operation overwrites it. If you exclude an...
  /// HTTP: PUT /2019-09-25/functions/{FunctionName}/event-invoke-config
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.FunctionEventInvokeConfig>
      putFunctionEventInvokeConfig(
    $0.PutFunctionEventInvokeConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putFunctionEventInvokeConfig, request,
        options: options);
  }

  /// Sends a failure response for a callback operation in a durable execution. Use this API when an external system cannot complete a callback operation successfully.
  /// HTTP: POST /2025-12-01/durable-execution-callbacks/{CallbackId}/fail
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SendDurableExecutionCallbackFailureResponse>
      sendDurableExecutionCallbackFailure(
    $0.SendDurableExecutionCallbackFailureRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendDurableExecutionCallbackFailure, request,
        options: options);
  }

  /// Sends a heartbeat signal for a long-running callback operation to prevent timeout. Use this API to extend the callback timeout period while the external operation is still in progress.
  /// HTTP: POST /2025-12-01/durable-execution-callbacks/{CallbackId}/heartbeat
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SendDurableExecutionCallbackHeartbeatResponse>
      sendDurableExecutionCallbackHeartbeat(
    $0.SendDurableExecutionCallbackHeartbeatRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendDurableExecutionCallbackHeartbeat, request,
        options: options);
  }

  /// Sends a successful completion response for a callback operation in a durable execution. Use this API when an external system has successfully completed a callback operation.
  /// HTTP: POST /2025-12-01/durable-execution-callbacks/{CallbackId}/succeed
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.SendDurableExecutionCallbackSuccessResponse>
      sendDurableExecutionCallbackSuccess(
    $0.SendDurableExecutionCallbackSuccessRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendDurableExecutionCallbackSuccess, request,
        options: options);
  }

  /// Stops a running durable execution. The execution transitions to STOPPED status and cannot be resumed. Any in-progress operations are terminated.
  /// HTTP: POST /2025-12-01/durable-executions/{DurableExecutionArn}/stop
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.StopDurableExecutionResponse> stopDurableExecution(
    $0.StopDurableExecutionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopDurableExecution, request, options: options);
  }

  /// Adds tags to a function, event source mapping, or code signing configuration.
  /// HTTP: POST /2017-03-31/tags/{Resource}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Removes tags from a function, event source mapping, or code signing configuration.
  /// HTTP: DELETE /2017-03-31/tags/{Resource}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Updates the configuration for asynchronous invocation for a function, version, or alias. To configure options for asynchronous invocation, use PutFunctionEventInvokeConfig.
  /// HTTP: POST /2019-09-25/functions/{FunctionName}/event-invoke-config
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.FunctionEventInvokeConfig>
      updateFunctionEventInvokeConfig(
    $0.UpdateFunctionEventInvokeConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateFunctionEventInvokeConfig, request,
        options: options);
  }

  /// Returns a list of capacity providers in your account.
  /// HTTP: GET /2025-11-30/capacity-providers
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListCapacityProvidersResponse> listCapacityProviders(
    $0.ListCapacityProvidersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCapacityProviders, request, options: options);
  }

  /// Creates a capacity provider that manages compute resources for Lambda functions
  /// HTTP: POST /2025-11-30/capacity-providers
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateCapacityProviderResponse>
      createCapacityProvider(
    $0.CreateCapacityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCapacityProvider, request,
        options: options);
  }

  /// Retrieves information about a specific capacity provider, including its configuration, state, and associated resources.
  /// HTTP: GET /2025-11-30/capacity-providers/{CapacityProviderName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetCapacityProviderResponse> getCapacityProvider(
    $0.GetCapacityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCapacityProvider, request, options: options);
  }

  /// Updates the configuration of an existing capacity provider.
  /// HTTP: PUT /2025-11-30/capacity-providers/{CapacityProviderName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateCapacityProviderResponse>
      updateCapacityProvider(
    $0.UpdateCapacityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateCapacityProvider, request,
        options: options);
  }

  /// Deletes a capacity provider. You cannot delete a capacity provider that is currently being used by Lambda functions.
  /// HTTP: DELETE /2025-11-30/capacity-providers/{CapacityProviderName}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.DeleteCapacityProviderResponse>
      deleteCapacityProvider(
    $0.DeleteCapacityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCapacityProvider, request,
        options: options);
  }

  /// Returns a list of function versions that are configured to use a specific capacity provider.
  /// HTTP: GET /2025-11-30/capacity-providers/{CapacityProviderName}/function-versions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListFunctionVersionsByCapacityProviderResponse>
      listFunctionVersionsByCapacityProvider(
    $0.ListFunctionVersionsByCapacityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFunctionVersionsByCapacityProvider, request,
        options: options);
  }

  /// Returns a list of code signing configurations. A request returns up to 10,000 configurations per call. You can use the MaxItems parameter to return fewer configurations per call.
  /// HTTP: GET /2020-04-22/code-signing-configs
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListCodeSigningConfigsResponse>
      listCodeSigningConfigs(
    $0.ListCodeSigningConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCodeSigningConfigs, request,
        options: options);
  }

  /// Creates a code signing configuration. A code signing configuration defines a list of allowed signing profiles and defines the code-signing validation policy (action to be taken if deployment valida...
  /// HTTP: POST /2020-04-22/code-signing-configs
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateCodeSigningConfigResponse>
      createCodeSigningConfig(
    $0.CreateCodeSigningConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createCodeSigningConfig, request,
        options: options);
  }

  /// Lists event source mappings. Specify an EventSourceArn to show only event source mappings for a single event source.
  /// HTTP: GET /2015-03-31/event-source-mappings
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListEventSourceMappingsResponse>
      listEventSourceMappings(
    $0.ListEventSourceMappingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEventSourceMappings, request,
        options: options);
  }

  /// Creates a mapping between an event source and an Lambda function. Lambda reads items from the event source and invokes the function. For details about how to configure different event sources, see ...
  /// HTTP: POST /2015-03-31/event-source-mappings
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.EventSourceMappingConfiguration>
      createEventSourceMapping(
    $0.CreateEventSourceMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createEventSourceMapping, request,
        options: options);
  }

  /// Returns details about an event source mapping. You can get the identifier of a mapping from the output of ListEventSourceMappings.
  /// HTTP: GET /2015-03-31/event-source-mappings/{UUID}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.EventSourceMappingConfiguration>
      getEventSourceMapping(
    $0.GetEventSourceMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getEventSourceMapping, request, options: options);
  }

  /// Updates an event source mapping. You can change the function that Lambda invokes, or pause invocation and resume later from the same location. For details about how to configure different event sou...
  /// HTTP: PUT /2015-03-31/event-source-mappings/{UUID}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.EventSourceMappingConfiguration>
      updateEventSourceMapping(
    $0.UpdateEventSourceMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateEventSourceMapping, request,
        options: options);
  }

  /// Deletes an event source mapping. You can get the identifier of a mapping from the output of ListEventSourceMappings. When you delete an event source mapping, it enters a Deleting state and might no...
  /// HTTP: DELETE /2015-03-31/event-source-mappings/{UUID}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.EventSourceMappingConfiguration>
      deleteEventSourceMapping(
    $0.DeleteEventSourceMappingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteEventSourceMapping, request,
        options: options);
  }

  /// Returns a list of Lambda functions, with the version-specific configuration of each. Lambda returns up to 50 functions per call. Set FunctionVersion to ALL to include all published versions of each...
  /// HTTP: GET /2015-03-31/functions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListFunctionsResponse> listFunctions(
    $0.ListFunctionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFunctions, request, options: options);
  }

  /// Creates a Lambda function. To create a function, you need a deployment package and an execution role. The deployment package is a .zip file archive or container image that contains your function co...
  /// HTTP: POST /2015-03-31/functions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.FunctionConfiguration> createFunction(
    $0.CreateFunctionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createFunction, request, options: options);
  }

  /// Creates a Lambda function URL with the specified configuration parameters. A function URL is a dedicated HTTP(S) endpoint that you can use to invoke your function.
  /// HTTP: POST /2021-10-31/functions/{FunctionName}/url
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.CreateFunctionUrlConfigResponse>
      createFunctionUrlConfig(
    $0.CreateFunctionUrlConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createFunctionUrlConfig, request,
        options: options);
  }

  /// Removes a concurrent execution limit from a function.
  /// HTTP: DELETE /2017-10-31/functions/{FunctionName}/concurrency
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteFunctionConcurrency(
    $0.DeleteFunctionConcurrencyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFunctionConcurrency, request,
        options: options);
  }

  /// Deletes a Lambda function URL. When you delete a function URL, you can't recover it. Creating a new function URL results in a different URL address.
  /// HTTP: DELETE /2021-10-31/functions/{FunctionName}/url
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteFunctionUrlConfig(
    $0.DeleteFunctionUrlConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteFunctionUrlConfig, request,
        options: options);
  }

  /// Returns details about the reserved concurrency configuration for a function. To set a concurrency limit for a function, use PutFunctionConcurrency.
  /// HTTP: GET /2019-09-30/functions/{FunctionName}/concurrency
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetFunctionConcurrencyResponse>
      getFunctionConcurrency(
    $0.GetFunctionConcurrencyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFunctionConcurrency, request,
        options: options);
  }

  /// Returns details about a Lambda function URL.
  /// HTTP: GET /2021-10-31/functions/{FunctionName}/url
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetFunctionUrlConfigResponse> getFunctionUrlConfig(
    $0.GetFunctionUrlConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getFunctionUrlConfig, request, options: options);
  }

  /// Returns a list of Lambda function URLs for the specified function.
  /// HTTP: GET /2021-10-31/functions/{FunctionName}/urls
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListFunctionUrlConfigsResponse>
      listFunctionUrlConfigs(
    $0.ListFunctionUrlConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listFunctionUrlConfigs, request,
        options: options);
  }

  /// Retrieves a list of provisioned concurrency configurations for a function.
  /// HTTP: GET /2019-09-30/functions/{FunctionName}/provisioned-concurrency?List=ALL
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListProvisionedConcurrencyConfigsResponse>
      listProvisionedConcurrencyConfigs(
    $0.ListProvisionedConcurrencyConfigsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listProvisionedConcurrencyConfigs, request,
        options: options);
  }

  /// Sets the maximum number of simultaneous executions for a function, and reserves capacity for that concurrency level. Concurrency settings apply to the function as a whole, including all published v...
  /// HTTP: PUT /2017-10-31/functions/{FunctionName}/concurrency
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.Concurrency> putFunctionConcurrency(
    $0.PutFunctionConcurrencyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putFunctionConcurrency, request,
        options: options);
  }

  /// Updates a Lambda function's code. If code signing is enabled for the function, the code package must be signed by a trusted publisher. For more information, see Configuring code signing for Lambda....
  /// HTTP: PUT /2015-03-31/functions/{FunctionName}/code
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.FunctionConfiguration> updateFunctionCode(
    $0.UpdateFunctionCodeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateFunctionCode, request, options: options);
  }

  /// Modify the version-specific settings of a Lambda function. When you update a function, Lambda provisions an instance of the function and its supporting resources. If your function connects to a VPC...
  /// HTTP: PUT /2015-03-31/functions/{FunctionName}/configuration
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.FunctionConfiguration> updateFunctionConfiguration(
    $0.UpdateFunctionConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateFunctionConfiguration, request,
        options: options);
  }

  /// Updates the configuration for a Lambda function URL.
  /// HTTP: PUT /2021-10-31/functions/{FunctionName}/url
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.UpdateFunctionUrlConfigResponse>
      updateFunctionUrlConfig(
    $0.UpdateFunctionUrlConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateFunctionUrlConfig, request,
        options: options);
  }

  /// Returns a list of aliases for a Lambda function.
  /// HTTP: GET /2015-03-31/functions/{FunctionName}/aliases
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListAliasesResponse> listAliases(
    $0.ListAliasesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAliases, request, options: options);
  }

  /// Creates an alias for a Lambda function version. Use aliases to provide clients with a function identifier that you can update to invoke a different version. You can also map an alias to split invoc...
  /// HTTP: POST /2015-03-31/functions/{FunctionName}/aliases
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.AliasConfiguration> createAlias(
    $0.CreateAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAlias, request, options: options);
  }

  /// Returns details about a Lambda function alias.
  /// HTTP: GET /2015-03-31/functions/{FunctionName}/aliases/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.AliasConfiguration> getAlias(
    $0.GetAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAlias, request, options: options);
  }

  /// Updates the configuration of a Lambda function alias.
  /// HTTP: PUT /2015-03-31/functions/{FunctionName}/aliases/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.AliasConfiguration> updateAlias(
    $0.UpdateAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAlias, request, options: options);
  }

  /// Deletes a Lambda function alias.
  /// HTTP: DELETE /2015-03-31/functions/{FunctionName}/aliases/{Name}
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteAlias(
    $0.DeleteAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAlias, request, options: options);
  }

  /// Creates a version from the current code and configuration of a function. Use versions to create a snapshot of your function code and configuration that doesn't change. Lambda doesn't publish a vers...
  /// HTTP: POST /2015-03-31/functions/{FunctionName}/versions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.FunctionConfiguration> publishVersion(
    $0.PublishVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$publishVersion, request, options: options);
  }

  /// Lists Lambda layers and shows information about the latest version of each. Specify a runtime identifier to list only layers that indicate that they're compatible with that runtime. Specify a compa...
  /// HTTP: GET /2018-10-31/layers
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListLayersResponse> listLayers(
    $0.ListLayersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listLayers, request, options: options);
  }

  /// Lists the versions of an Lambda layer. Versions that have been deleted aren't listed. Specify a runtime identifier to list only versions that indicate that they're compatible with that runtime. Spe...
  /// HTTP: GET /2018-10-31/layers/{LayerName}/versions
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.ListLayerVersionsResponse> listLayerVersions(
    $0.ListLayerVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listLayerVersions, request, options: options);
  }

  /// Adds a provisioned concurrency configuration to a function's alias or version.
  /// HTTP: PUT /2019-09-30/functions/{FunctionName}/provisioned-concurrency
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.PutProvisionedConcurrencyConfigResponse>
      putProvisionedConcurrencyConfig(
    $0.PutProvisionedConcurrencyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putProvisionedConcurrencyConfig, request,
        options: options);
  }

  /// Retrieves the provisioned concurrency configuration for a function's alias or version.
  /// HTTP: GET /2019-09-30/functions/{FunctionName}/provisioned-concurrency
  /// Protocol: restJson1
  $grpc.ResponseFuture<$0.GetProvisionedConcurrencyConfigResponse>
      getProvisionedConcurrencyConfig(
    $0.GetProvisionedConcurrencyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getProvisionedConcurrencyConfig, request,
        options: options);
  }

  /// Deletes the provisioned concurrency configuration for a function.
  /// HTTP: DELETE /2019-09-30/functions/{FunctionName}/provisioned-concurrency
  /// Protocol: restJson1
  $grpc.ResponseFuture<$1.Empty> deleteProvisionedConcurrencyConfig(
    $0.DeleteProvisionedConcurrencyConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteProvisionedConcurrencyConfig, request,
        options: options);
  }

  // method descriptors

  static final _$checkpointDurableExecution = $grpc.ClientMethod<
          $0.CheckpointDurableExecutionRequest,
          $0.CheckpointDurableExecutionResponse>(
      '/lambda.LambdaService/CheckpointDurableExecution',
      ($0.CheckpointDurableExecutionRequest value) => value.writeToBuffer(),
      $0.CheckpointDurableExecutionResponse.fromBuffer);
  static final _$deleteFunction =
      $grpc.ClientMethod<$0.DeleteFunctionRequest, $0.DeleteFunctionResponse>(
          '/lambda.LambdaService/DeleteFunction',
          ($0.DeleteFunctionRequest value) => value.writeToBuffer(),
          $0.DeleteFunctionResponse.fromBuffer);
  static final _$deleteFunctionEventInvokeConfig =
      $grpc.ClientMethod<$0.DeleteFunctionEventInvokeConfigRequest, $1.Empty>(
          '/lambda.LambdaService/DeleteFunctionEventInvokeConfig',
          ($0.DeleteFunctionEventInvokeConfigRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$getAccountSettings = $grpc.ClientMethod<
          $0.GetAccountSettingsRequest, $0.GetAccountSettingsResponse>(
      '/lambda.LambdaService/GetAccountSettings',
      ($0.GetAccountSettingsRequest value) => value.writeToBuffer(),
      $0.GetAccountSettingsResponse.fromBuffer);
  static final _$getDurableExecution = $grpc.ClientMethod<
          $0.GetDurableExecutionRequest, $0.GetDurableExecutionResponse>(
      '/lambda.LambdaService/GetDurableExecution',
      ($0.GetDurableExecutionRequest value) => value.writeToBuffer(),
      $0.GetDurableExecutionResponse.fromBuffer);
  static final _$getDurableExecutionHistory = $grpc.ClientMethod<
          $0.GetDurableExecutionHistoryRequest,
          $0.GetDurableExecutionHistoryResponse>(
      '/lambda.LambdaService/GetDurableExecutionHistory',
      ($0.GetDurableExecutionHistoryRequest value) => value.writeToBuffer(),
      $0.GetDurableExecutionHistoryResponse.fromBuffer);
  static final _$getDurableExecutionState = $grpc.ClientMethod<
          $0.GetDurableExecutionStateRequest,
          $0.GetDurableExecutionStateResponse>(
      '/lambda.LambdaService/GetDurableExecutionState',
      ($0.GetDurableExecutionStateRequest value) => value.writeToBuffer(),
      $0.GetDurableExecutionStateResponse.fromBuffer);
  static final _$getFunctionEventInvokeConfig = $grpc.ClientMethod<
          $0.GetFunctionEventInvokeConfigRequest, $0.FunctionEventInvokeConfig>(
      '/lambda.LambdaService/GetFunctionEventInvokeConfig',
      ($0.GetFunctionEventInvokeConfigRequest value) => value.writeToBuffer(),
      $0.FunctionEventInvokeConfig.fromBuffer);
  static final _$listDurableExecutionsByFunction = $grpc.ClientMethod<
          $0.ListDurableExecutionsByFunctionRequest,
          $0.ListDurableExecutionsByFunctionResponse>(
      '/lambda.LambdaService/ListDurableExecutionsByFunction',
      ($0.ListDurableExecutionsByFunctionRequest value) =>
          value.writeToBuffer(),
      $0.ListDurableExecutionsByFunctionResponse.fromBuffer);
  static final _$listFunctionEventInvokeConfigs = $grpc.ClientMethod<
          $0.ListFunctionEventInvokeConfigsRequest,
          $0.ListFunctionEventInvokeConfigsResponse>(
      '/lambda.LambdaService/ListFunctionEventInvokeConfigs',
      ($0.ListFunctionEventInvokeConfigsRequest value) => value.writeToBuffer(),
      $0.ListFunctionEventInvokeConfigsResponse.fromBuffer);
  static final _$listTags =
      $grpc.ClientMethod<$0.ListTagsRequest, $0.ListTagsResponse>(
          '/lambda.LambdaService/ListTags',
          ($0.ListTagsRequest value) => value.writeToBuffer(),
          $0.ListTagsResponse.fromBuffer);
  static final _$putFunctionEventInvokeConfig = $grpc.ClientMethod<
          $0.PutFunctionEventInvokeConfigRequest, $0.FunctionEventInvokeConfig>(
      '/lambda.LambdaService/PutFunctionEventInvokeConfig',
      ($0.PutFunctionEventInvokeConfigRequest value) => value.writeToBuffer(),
      $0.FunctionEventInvokeConfig.fromBuffer);
  static final _$sendDurableExecutionCallbackFailure = $grpc.ClientMethod<
          $0.SendDurableExecutionCallbackFailureRequest,
          $0.SendDurableExecutionCallbackFailureResponse>(
      '/lambda.LambdaService/SendDurableExecutionCallbackFailure',
      ($0.SendDurableExecutionCallbackFailureRequest value) =>
          value.writeToBuffer(),
      $0.SendDurableExecutionCallbackFailureResponse.fromBuffer);
  static final _$sendDurableExecutionCallbackHeartbeat = $grpc.ClientMethod<
          $0.SendDurableExecutionCallbackHeartbeatRequest,
          $0.SendDurableExecutionCallbackHeartbeatResponse>(
      '/lambda.LambdaService/SendDurableExecutionCallbackHeartbeat',
      ($0.SendDurableExecutionCallbackHeartbeatRequest value) =>
          value.writeToBuffer(),
      $0.SendDurableExecutionCallbackHeartbeatResponse.fromBuffer);
  static final _$sendDurableExecutionCallbackSuccess = $grpc.ClientMethod<
          $0.SendDurableExecutionCallbackSuccessRequest,
          $0.SendDurableExecutionCallbackSuccessResponse>(
      '/lambda.LambdaService/SendDurableExecutionCallbackSuccess',
      ($0.SendDurableExecutionCallbackSuccessRequest value) =>
          value.writeToBuffer(),
      $0.SendDurableExecutionCallbackSuccessResponse.fromBuffer);
  static final _$stopDurableExecution = $grpc.ClientMethod<
          $0.StopDurableExecutionRequest, $0.StopDurableExecutionResponse>(
      '/lambda.LambdaService/StopDurableExecution',
      ($0.StopDurableExecutionRequest value) => value.writeToBuffer(),
      $0.StopDurableExecutionResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $1.Empty>(
          '/lambda.LambdaService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $1.Empty>(
          '/lambda.LambdaService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateFunctionEventInvokeConfig = $grpc.ClientMethod<
          $0.UpdateFunctionEventInvokeConfigRequest,
          $0.FunctionEventInvokeConfig>(
      '/lambda.LambdaService/UpdateFunctionEventInvokeConfig',
      ($0.UpdateFunctionEventInvokeConfigRequest value) =>
          value.writeToBuffer(),
      $0.FunctionEventInvokeConfig.fromBuffer);
  static final _$listCapacityProviders = $grpc.ClientMethod<
          $0.ListCapacityProvidersRequest, $0.ListCapacityProvidersResponse>(
      '/lambda.LambdaService/ListCapacityProviders',
      ($0.ListCapacityProvidersRequest value) => value.writeToBuffer(),
      $0.ListCapacityProvidersResponse.fromBuffer);
  static final _$createCapacityProvider = $grpc.ClientMethod<
          $0.CreateCapacityProviderRequest, $0.CreateCapacityProviderResponse>(
      '/lambda.LambdaService/CreateCapacityProvider',
      ($0.CreateCapacityProviderRequest value) => value.writeToBuffer(),
      $0.CreateCapacityProviderResponse.fromBuffer);
  static final _$getCapacityProvider = $grpc.ClientMethod<
          $0.GetCapacityProviderRequest, $0.GetCapacityProviderResponse>(
      '/lambda.LambdaService/GetCapacityProvider',
      ($0.GetCapacityProviderRequest value) => value.writeToBuffer(),
      $0.GetCapacityProviderResponse.fromBuffer);
  static final _$updateCapacityProvider = $grpc.ClientMethod<
          $0.UpdateCapacityProviderRequest, $0.UpdateCapacityProviderResponse>(
      '/lambda.LambdaService/UpdateCapacityProvider',
      ($0.UpdateCapacityProviderRequest value) => value.writeToBuffer(),
      $0.UpdateCapacityProviderResponse.fromBuffer);
  static final _$deleteCapacityProvider = $grpc.ClientMethod<
          $0.DeleteCapacityProviderRequest, $0.DeleteCapacityProviderResponse>(
      '/lambda.LambdaService/DeleteCapacityProvider',
      ($0.DeleteCapacityProviderRequest value) => value.writeToBuffer(),
      $0.DeleteCapacityProviderResponse.fromBuffer);
  static final _$listFunctionVersionsByCapacityProvider = $grpc.ClientMethod<
          $0.ListFunctionVersionsByCapacityProviderRequest,
          $0.ListFunctionVersionsByCapacityProviderResponse>(
      '/lambda.LambdaService/ListFunctionVersionsByCapacityProvider',
      ($0.ListFunctionVersionsByCapacityProviderRequest value) =>
          value.writeToBuffer(),
      $0.ListFunctionVersionsByCapacityProviderResponse.fromBuffer);
  static final _$listCodeSigningConfigs = $grpc.ClientMethod<
          $0.ListCodeSigningConfigsRequest, $0.ListCodeSigningConfigsResponse>(
      '/lambda.LambdaService/ListCodeSigningConfigs',
      ($0.ListCodeSigningConfigsRequest value) => value.writeToBuffer(),
      $0.ListCodeSigningConfigsResponse.fromBuffer);
  static final _$createCodeSigningConfig = $grpc.ClientMethod<
          $0.CreateCodeSigningConfigRequest,
          $0.CreateCodeSigningConfigResponse>(
      '/lambda.LambdaService/CreateCodeSigningConfig',
      ($0.CreateCodeSigningConfigRequest value) => value.writeToBuffer(),
      $0.CreateCodeSigningConfigResponse.fromBuffer);
  static final _$listEventSourceMappings = $grpc.ClientMethod<
          $0.ListEventSourceMappingsRequest,
          $0.ListEventSourceMappingsResponse>(
      '/lambda.LambdaService/ListEventSourceMappings',
      ($0.ListEventSourceMappingsRequest value) => value.writeToBuffer(),
      $0.ListEventSourceMappingsResponse.fromBuffer);
  static final _$createEventSourceMapping = $grpc.ClientMethod<
          $0.CreateEventSourceMappingRequest,
          $0.EventSourceMappingConfiguration>(
      '/lambda.LambdaService/CreateEventSourceMapping',
      ($0.CreateEventSourceMappingRequest value) => value.writeToBuffer(),
      $0.EventSourceMappingConfiguration.fromBuffer);
  static final _$getEventSourceMapping = $grpc.ClientMethod<
          $0.GetEventSourceMappingRequest, $0.EventSourceMappingConfiguration>(
      '/lambda.LambdaService/GetEventSourceMapping',
      ($0.GetEventSourceMappingRequest value) => value.writeToBuffer(),
      $0.EventSourceMappingConfiguration.fromBuffer);
  static final _$updateEventSourceMapping = $grpc.ClientMethod<
          $0.UpdateEventSourceMappingRequest,
          $0.EventSourceMappingConfiguration>(
      '/lambda.LambdaService/UpdateEventSourceMapping',
      ($0.UpdateEventSourceMappingRequest value) => value.writeToBuffer(),
      $0.EventSourceMappingConfiguration.fromBuffer);
  static final _$deleteEventSourceMapping = $grpc.ClientMethod<
          $0.DeleteEventSourceMappingRequest,
          $0.EventSourceMappingConfiguration>(
      '/lambda.LambdaService/DeleteEventSourceMapping',
      ($0.DeleteEventSourceMappingRequest value) => value.writeToBuffer(),
      $0.EventSourceMappingConfiguration.fromBuffer);
  static final _$listFunctions =
      $grpc.ClientMethod<$0.ListFunctionsRequest, $0.ListFunctionsResponse>(
          '/lambda.LambdaService/ListFunctions',
          ($0.ListFunctionsRequest value) => value.writeToBuffer(),
          $0.ListFunctionsResponse.fromBuffer);
  static final _$createFunction =
      $grpc.ClientMethod<$0.CreateFunctionRequest, $0.FunctionConfiguration>(
          '/lambda.LambdaService/CreateFunction',
          ($0.CreateFunctionRequest value) => value.writeToBuffer(),
          $0.FunctionConfiguration.fromBuffer);
  static final _$createFunctionUrlConfig = $grpc.ClientMethod<
          $0.CreateFunctionUrlConfigRequest,
          $0.CreateFunctionUrlConfigResponse>(
      '/lambda.LambdaService/CreateFunctionUrlConfig',
      ($0.CreateFunctionUrlConfigRequest value) => value.writeToBuffer(),
      $0.CreateFunctionUrlConfigResponse.fromBuffer);
  static final _$deleteFunctionConcurrency =
      $grpc.ClientMethod<$0.DeleteFunctionConcurrencyRequest, $1.Empty>(
          '/lambda.LambdaService/DeleteFunctionConcurrency',
          ($0.DeleteFunctionConcurrencyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteFunctionUrlConfig =
      $grpc.ClientMethod<$0.DeleteFunctionUrlConfigRequest, $1.Empty>(
          '/lambda.LambdaService/DeleteFunctionUrlConfig',
          ($0.DeleteFunctionUrlConfigRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$getFunctionConcurrency = $grpc.ClientMethod<
          $0.GetFunctionConcurrencyRequest, $0.GetFunctionConcurrencyResponse>(
      '/lambda.LambdaService/GetFunctionConcurrency',
      ($0.GetFunctionConcurrencyRequest value) => value.writeToBuffer(),
      $0.GetFunctionConcurrencyResponse.fromBuffer);
  static final _$getFunctionUrlConfig = $grpc.ClientMethod<
          $0.GetFunctionUrlConfigRequest, $0.GetFunctionUrlConfigResponse>(
      '/lambda.LambdaService/GetFunctionUrlConfig',
      ($0.GetFunctionUrlConfigRequest value) => value.writeToBuffer(),
      $0.GetFunctionUrlConfigResponse.fromBuffer);
  static final _$listFunctionUrlConfigs = $grpc.ClientMethod<
          $0.ListFunctionUrlConfigsRequest, $0.ListFunctionUrlConfigsResponse>(
      '/lambda.LambdaService/ListFunctionUrlConfigs',
      ($0.ListFunctionUrlConfigsRequest value) => value.writeToBuffer(),
      $0.ListFunctionUrlConfigsResponse.fromBuffer);
  static final _$listProvisionedConcurrencyConfigs = $grpc.ClientMethod<
          $0.ListProvisionedConcurrencyConfigsRequest,
          $0.ListProvisionedConcurrencyConfigsResponse>(
      '/lambda.LambdaService/ListProvisionedConcurrencyConfigs',
      ($0.ListProvisionedConcurrencyConfigsRequest value) =>
          value.writeToBuffer(),
      $0.ListProvisionedConcurrencyConfigsResponse.fromBuffer);
  static final _$putFunctionConcurrency =
      $grpc.ClientMethod<$0.PutFunctionConcurrencyRequest, $0.Concurrency>(
          '/lambda.LambdaService/PutFunctionConcurrency',
          ($0.PutFunctionConcurrencyRequest value) => value.writeToBuffer(),
          $0.Concurrency.fromBuffer);
  static final _$updateFunctionCode = $grpc.ClientMethod<
          $0.UpdateFunctionCodeRequest, $0.FunctionConfiguration>(
      '/lambda.LambdaService/UpdateFunctionCode',
      ($0.UpdateFunctionCodeRequest value) => value.writeToBuffer(),
      $0.FunctionConfiguration.fromBuffer);
  static final _$updateFunctionConfiguration = $grpc.ClientMethod<
          $0.UpdateFunctionConfigurationRequest, $0.FunctionConfiguration>(
      '/lambda.LambdaService/UpdateFunctionConfiguration',
      ($0.UpdateFunctionConfigurationRequest value) => value.writeToBuffer(),
      $0.FunctionConfiguration.fromBuffer);
  static final _$updateFunctionUrlConfig = $grpc.ClientMethod<
          $0.UpdateFunctionUrlConfigRequest,
          $0.UpdateFunctionUrlConfigResponse>(
      '/lambda.LambdaService/UpdateFunctionUrlConfig',
      ($0.UpdateFunctionUrlConfigRequest value) => value.writeToBuffer(),
      $0.UpdateFunctionUrlConfigResponse.fromBuffer);
  static final _$listAliases =
      $grpc.ClientMethod<$0.ListAliasesRequest, $0.ListAliasesResponse>(
          '/lambda.LambdaService/ListAliases',
          ($0.ListAliasesRequest value) => value.writeToBuffer(),
          $0.ListAliasesResponse.fromBuffer);
  static final _$createAlias =
      $grpc.ClientMethod<$0.CreateAliasRequest, $0.AliasConfiguration>(
          '/lambda.LambdaService/CreateAlias',
          ($0.CreateAliasRequest value) => value.writeToBuffer(),
          $0.AliasConfiguration.fromBuffer);
  static final _$getAlias =
      $grpc.ClientMethod<$0.GetAliasRequest, $0.AliasConfiguration>(
          '/lambda.LambdaService/GetAlias',
          ($0.GetAliasRequest value) => value.writeToBuffer(),
          $0.AliasConfiguration.fromBuffer);
  static final _$updateAlias =
      $grpc.ClientMethod<$0.UpdateAliasRequest, $0.AliasConfiguration>(
          '/lambda.LambdaService/UpdateAlias',
          ($0.UpdateAliasRequest value) => value.writeToBuffer(),
          $0.AliasConfiguration.fromBuffer);
  static final _$deleteAlias =
      $grpc.ClientMethod<$0.DeleteAliasRequest, $1.Empty>(
          '/lambda.LambdaService/DeleteAlias',
          ($0.DeleteAliasRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$publishVersion =
      $grpc.ClientMethod<$0.PublishVersionRequest, $0.FunctionConfiguration>(
          '/lambda.LambdaService/PublishVersion',
          ($0.PublishVersionRequest value) => value.writeToBuffer(),
          $0.FunctionConfiguration.fromBuffer);
  static final _$listLayers =
      $grpc.ClientMethod<$0.ListLayersRequest, $0.ListLayersResponse>(
          '/lambda.LambdaService/ListLayers',
          ($0.ListLayersRequest value) => value.writeToBuffer(),
          $0.ListLayersResponse.fromBuffer);
  static final _$listLayerVersions = $grpc.ClientMethod<
          $0.ListLayerVersionsRequest, $0.ListLayerVersionsResponse>(
      '/lambda.LambdaService/ListLayerVersions',
      ($0.ListLayerVersionsRequest value) => value.writeToBuffer(),
      $0.ListLayerVersionsResponse.fromBuffer);
  static final _$putProvisionedConcurrencyConfig = $grpc.ClientMethod<
          $0.PutProvisionedConcurrencyConfigRequest,
          $0.PutProvisionedConcurrencyConfigResponse>(
      '/lambda.LambdaService/PutProvisionedConcurrencyConfig',
      ($0.PutProvisionedConcurrencyConfigRequest value) =>
          value.writeToBuffer(),
      $0.PutProvisionedConcurrencyConfigResponse.fromBuffer);
  static final _$getProvisionedConcurrencyConfig = $grpc.ClientMethod<
          $0.GetProvisionedConcurrencyConfigRequest,
          $0.GetProvisionedConcurrencyConfigResponse>(
      '/lambda.LambdaService/GetProvisionedConcurrencyConfig',
      ($0.GetProvisionedConcurrencyConfigRequest value) =>
          value.writeToBuffer(),
      $0.GetProvisionedConcurrencyConfigResponse.fromBuffer);
  static final _$deleteProvisionedConcurrencyConfig = $grpc.ClientMethod<
          $0.DeleteProvisionedConcurrencyConfigRequest, $1.Empty>(
      '/lambda.LambdaService/DeleteProvisionedConcurrencyConfig',
      ($0.DeleteProvisionedConcurrencyConfigRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
}

@$pb.GrpcServiceName('lambda.LambdaService')
abstract class LambdaServiceBase extends $grpc.Service {
  $core.String get $name => 'lambda.LambdaService';

  LambdaServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.CheckpointDurableExecutionRequest,
            $0.CheckpointDurableExecutionResponse>(
        'CheckpointDurableExecution',
        checkpointDurableExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CheckpointDurableExecutionRequest.fromBuffer(value),
        ($0.CheckpointDurableExecutionResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteFunctionRequest,
            $0.DeleteFunctionResponse>(
        'DeleteFunction',
        deleteFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteFunctionRequest.fromBuffer(value),
        ($0.DeleteFunctionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteFunctionEventInvokeConfigRequest,
            $1.Empty>(
        'DeleteFunctionEventInvokeConfig',
        deleteFunctionEventInvokeConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteFunctionEventInvokeConfigRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccountSettingsRequest,
            $0.GetAccountSettingsResponse>(
        'GetAccountSettings',
        getAccountSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAccountSettingsRequest.fromBuffer(value),
        ($0.GetAccountSettingsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDurableExecutionRequest,
            $0.GetDurableExecutionResponse>(
        'GetDurableExecution',
        getDurableExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDurableExecutionRequest.fromBuffer(value),
        ($0.GetDurableExecutionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDurableExecutionHistoryRequest,
            $0.GetDurableExecutionHistoryResponse>(
        'GetDurableExecutionHistory',
        getDurableExecutionHistory_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDurableExecutionHistoryRequest.fromBuffer(value),
        ($0.GetDurableExecutionHistoryResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDurableExecutionStateRequest,
            $0.GetDurableExecutionStateResponse>(
        'GetDurableExecutionState',
        getDurableExecutionState_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDurableExecutionStateRequest.fromBuffer(value),
        ($0.GetDurableExecutionStateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFunctionEventInvokeConfigRequest,
            $0.FunctionEventInvokeConfig>(
        'GetFunctionEventInvokeConfig',
        getFunctionEventInvokeConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFunctionEventInvokeConfigRequest.fromBuffer(value),
        ($0.FunctionEventInvokeConfig value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDurableExecutionsByFunctionRequest,
            $0.ListDurableExecutionsByFunctionResponse>(
        'ListDurableExecutionsByFunction',
        listDurableExecutionsByFunction_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDurableExecutionsByFunctionRequest.fromBuffer(value),
        ($0.ListDurableExecutionsByFunctionResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListFunctionEventInvokeConfigsRequest,
            $0.ListFunctionEventInvokeConfigsResponse>(
        'ListFunctionEventInvokeConfigs',
        listFunctionEventInvokeConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListFunctionEventInvokeConfigsRequest.fromBuffer(value),
        ($0.ListFunctionEventInvokeConfigsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsRequest, $0.ListTagsResponse>(
        'ListTags',
        listTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListTagsRequest.fromBuffer(value),
        ($0.ListTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutFunctionEventInvokeConfigRequest,
            $0.FunctionEventInvokeConfig>(
        'PutFunctionEventInvokeConfig',
        putFunctionEventInvokeConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutFunctionEventInvokeConfigRequest.fromBuffer(value),
        ($0.FunctionEventInvokeConfig value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.SendDurableExecutionCallbackFailureRequest,
            $0.SendDurableExecutionCallbackFailureResponse>(
        'SendDurableExecutionCallbackFailure',
        sendDurableExecutionCallbackFailure_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendDurableExecutionCallbackFailureRequest.fromBuffer(value),
        ($0.SendDurableExecutionCallbackFailureResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.SendDurableExecutionCallbackHeartbeatRequest,
            $0.SendDurableExecutionCallbackHeartbeatResponse>(
        'SendDurableExecutionCallbackHeartbeat',
        sendDurableExecutionCallbackHeartbeat_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendDurableExecutionCallbackHeartbeatRequest.fromBuffer(value),
        ($0.SendDurableExecutionCallbackHeartbeatResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.SendDurableExecutionCallbackSuccessRequest,
            $0.SendDurableExecutionCallbackSuccessResponse>(
        'SendDurableExecutionCallbackSuccess',
        sendDurableExecutionCallbackSuccess_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendDurableExecutionCallbackSuccessRequest.fromBuffer(value),
        ($0.SendDurableExecutionCallbackSuccessResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopDurableExecutionRequest,
            $0.StopDurableExecutionResponse>(
        'StopDurableExecution',
        stopDurableExecution_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopDurableExecutionRequest.fromBuffer(value),
        ($0.StopDurableExecutionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagResourceRequest, $1.Empty>(
        'TagResource',
        tagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagResourceRequest, $1.Empty>(
        'UntagResource',
        untagResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagResourceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateFunctionEventInvokeConfigRequest,
            $0.FunctionEventInvokeConfig>(
        'UpdateFunctionEventInvokeConfig',
        updateFunctionEventInvokeConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateFunctionEventInvokeConfigRequest.fromBuffer(value),
        ($0.FunctionEventInvokeConfig value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCapacityProvidersRequest,
            $0.ListCapacityProvidersResponse>(
        'ListCapacityProviders',
        listCapacityProviders_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCapacityProvidersRequest.fromBuffer(value),
        ($0.ListCapacityProvidersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateCapacityProviderRequest,
            $0.CreateCapacityProviderResponse>(
        'CreateCapacityProvider',
        createCapacityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCapacityProviderRequest.fromBuffer(value),
        ($0.CreateCapacityProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCapacityProviderRequest,
            $0.GetCapacityProviderResponse>(
        'GetCapacityProvider',
        getCapacityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCapacityProviderRequest.fromBuffer(value),
        ($0.GetCapacityProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateCapacityProviderRequest,
            $0.UpdateCapacityProviderResponse>(
        'UpdateCapacityProvider',
        updateCapacityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateCapacityProviderRequest.fromBuffer(value),
        ($0.UpdateCapacityProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteCapacityProviderRequest,
            $0.DeleteCapacityProviderResponse>(
        'DeleteCapacityProvider',
        deleteCapacityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCapacityProviderRequest.fromBuffer(value),
        ($0.DeleteCapacityProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.ListFunctionVersionsByCapacityProviderRequest,
            $0.ListFunctionVersionsByCapacityProviderResponse>(
        'ListFunctionVersionsByCapacityProvider',
        listFunctionVersionsByCapacityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListFunctionVersionsByCapacityProviderRequest.fromBuffer(value),
        ($0.ListFunctionVersionsByCapacityProviderResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCodeSigningConfigsRequest,
            $0.ListCodeSigningConfigsResponse>(
        'ListCodeSigningConfigs',
        listCodeSigningConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCodeSigningConfigsRequest.fromBuffer(value),
        ($0.ListCodeSigningConfigsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateCodeSigningConfigRequest,
            $0.CreateCodeSigningConfigResponse>(
        'CreateCodeSigningConfig',
        createCodeSigningConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateCodeSigningConfigRequest.fromBuffer(value),
        ($0.CreateCodeSigningConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEventSourceMappingsRequest,
            $0.ListEventSourceMappingsResponse>(
        'ListEventSourceMappings',
        listEventSourceMappings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEventSourceMappingsRequest.fromBuffer(value),
        ($0.ListEventSourceMappingsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateEventSourceMappingRequest,
            $0.EventSourceMappingConfiguration>(
        'CreateEventSourceMapping',
        createEventSourceMapping_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateEventSourceMappingRequest.fromBuffer(value),
        ($0.EventSourceMappingConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetEventSourceMappingRequest,
            $0.EventSourceMappingConfiguration>(
        'GetEventSourceMapping',
        getEventSourceMapping_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetEventSourceMappingRequest.fromBuffer(value),
        ($0.EventSourceMappingConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateEventSourceMappingRequest,
            $0.EventSourceMappingConfiguration>(
        'UpdateEventSourceMapping',
        updateEventSourceMapping_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateEventSourceMappingRequest.fromBuffer(value),
        ($0.EventSourceMappingConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteEventSourceMappingRequest,
            $0.EventSourceMappingConfiguration>(
        'DeleteEventSourceMapping',
        deleteEventSourceMapping_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteEventSourceMappingRequest.fromBuffer(value),
        ($0.EventSourceMappingConfiguration value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListFunctionsRequest, $0.ListFunctionsResponse>(
            'ListFunctions',
            listFunctions_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListFunctionsRequest.fromBuffer(value),
            ($0.ListFunctionsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateFunctionRequest, $0.FunctionConfiguration>(
            'CreateFunction',
            createFunction_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateFunctionRequest.fromBuffer(value),
            ($0.FunctionConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateFunctionUrlConfigRequest,
            $0.CreateFunctionUrlConfigResponse>(
        'CreateFunctionUrlConfig',
        createFunctionUrlConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateFunctionUrlConfigRequest.fromBuffer(value),
        ($0.CreateFunctionUrlConfigResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteFunctionConcurrencyRequest, $1.Empty>(
            'DeleteFunctionConcurrency',
            deleteFunctionConcurrency_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteFunctionConcurrencyRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteFunctionUrlConfigRequest, $1.Empty>(
        'DeleteFunctionUrlConfig',
        deleteFunctionUrlConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteFunctionUrlConfigRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFunctionConcurrencyRequest,
            $0.GetFunctionConcurrencyResponse>(
        'GetFunctionConcurrency',
        getFunctionConcurrency_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFunctionConcurrencyRequest.fromBuffer(value),
        ($0.GetFunctionConcurrencyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetFunctionUrlConfigRequest,
            $0.GetFunctionUrlConfigResponse>(
        'GetFunctionUrlConfig',
        getFunctionUrlConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetFunctionUrlConfigRequest.fromBuffer(value),
        ($0.GetFunctionUrlConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListFunctionUrlConfigsRequest,
            $0.ListFunctionUrlConfigsResponse>(
        'ListFunctionUrlConfigs',
        listFunctionUrlConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListFunctionUrlConfigsRequest.fromBuffer(value),
        ($0.ListFunctionUrlConfigsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListProvisionedConcurrencyConfigsRequest,
            $0.ListProvisionedConcurrencyConfigsResponse>(
        'ListProvisionedConcurrencyConfigs',
        listProvisionedConcurrencyConfigs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListProvisionedConcurrencyConfigsRequest.fromBuffer(value),
        ($0.ListProvisionedConcurrencyConfigsResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutFunctionConcurrencyRequest, $0.Concurrency>(
            'PutFunctionConcurrency',
            putFunctionConcurrency_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutFunctionConcurrencyRequest.fromBuffer(value),
            ($0.Concurrency value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateFunctionCodeRequest,
            $0.FunctionConfiguration>(
        'UpdateFunctionCode',
        updateFunctionCode_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateFunctionCodeRequest.fromBuffer(value),
        ($0.FunctionConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateFunctionConfigurationRequest,
            $0.FunctionConfiguration>(
        'UpdateFunctionConfiguration',
        updateFunctionConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateFunctionConfigurationRequest.fromBuffer(value),
        ($0.FunctionConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateFunctionUrlConfigRequest,
            $0.UpdateFunctionUrlConfigResponse>(
        'UpdateFunctionUrlConfig',
        updateFunctionUrlConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateFunctionUrlConfigRequest.fromBuffer(value),
        ($0.UpdateFunctionUrlConfigResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListAliasesRequest, $0.ListAliasesResponse>(
            'ListAliases',
            listAliases_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListAliasesRequest.fromBuffer(value),
            ($0.ListAliasesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateAliasRequest, $0.AliasConfiguration>(
            'CreateAlias',
            createAlias_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateAliasRequest.fromBuffer(value),
            ($0.AliasConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAliasRequest, $0.AliasConfiguration>(
        'GetAlias',
        getAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetAliasRequest.fromBuffer(value),
        ($0.AliasConfiguration value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateAliasRequest, $0.AliasConfiguration>(
            'UpdateAlias',
            updateAlias_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateAliasRequest.fromBuffer(value),
            ($0.AliasConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAliasRequest, $1.Empty>(
        'DeleteAlias',
        deleteAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAliasRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PublishVersionRequest, $0.FunctionConfiguration>(
            'PublishVersion',
            publishVersion_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PublishVersionRequest.fromBuffer(value),
            ($0.FunctionConfiguration value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListLayersRequest, $0.ListLayersResponse>(
        'ListLayers',
        listLayers_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListLayersRequest.fromBuffer(value),
        ($0.ListLayersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListLayerVersionsRequest,
            $0.ListLayerVersionsResponse>(
        'ListLayerVersions',
        listLayerVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListLayerVersionsRequest.fromBuffer(value),
        ($0.ListLayerVersionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutProvisionedConcurrencyConfigRequest,
            $0.PutProvisionedConcurrencyConfigResponse>(
        'PutProvisionedConcurrencyConfig',
        putProvisionedConcurrencyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutProvisionedConcurrencyConfigRequest.fromBuffer(value),
        ($0.PutProvisionedConcurrencyConfigResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetProvisionedConcurrencyConfigRequest,
            $0.GetProvisionedConcurrencyConfigResponse>(
        'GetProvisionedConcurrencyConfig',
        getProvisionedConcurrencyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetProvisionedConcurrencyConfigRequest.fromBuffer(value),
        ($0.GetProvisionedConcurrencyConfigResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteProvisionedConcurrencyConfigRequest,
            $1.Empty>(
        'DeleteProvisionedConcurrencyConfig',
        deleteProvisionedConcurrencyConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteProvisionedConcurrencyConfigRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
  }

  $async.Future<$0.CheckpointDurableExecutionResponse>
      checkpointDurableExecution_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CheckpointDurableExecutionRequest> $request) async {
    return checkpointDurableExecution($call, await $request);
  }

  $async.Future<$0.CheckpointDurableExecutionResponse>
      checkpointDurableExecution(
          $grpc.ServiceCall call, $0.CheckpointDurableExecutionRequest request);

  $async.Future<$0.DeleteFunctionResponse> deleteFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteFunctionRequest> $request) async {
    return deleteFunction($call, await $request);
  }

  $async.Future<$0.DeleteFunctionResponse> deleteFunction(
      $grpc.ServiceCall call, $0.DeleteFunctionRequest request);

  $async.Future<$1.Empty> deleteFunctionEventInvokeConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteFunctionEventInvokeConfigRequest> $request) async {
    return deleteFunctionEventInvokeConfig($call, await $request);
  }

  $async.Future<$1.Empty> deleteFunctionEventInvokeConfig(
      $grpc.ServiceCall call,
      $0.DeleteFunctionEventInvokeConfigRequest request);

  $async.Future<$0.GetAccountSettingsResponse> getAccountSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAccountSettingsRequest> $request) async {
    return getAccountSettings($call, await $request);
  }

  $async.Future<$0.GetAccountSettingsResponse> getAccountSettings(
      $grpc.ServiceCall call, $0.GetAccountSettingsRequest request);

  $async.Future<$0.GetDurableExecutionResponse> getDurableExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDurableExecutionRequest> $request) async {
    return getDurableExecution($call, await $request);
  }

  $async.Future<$0.GetDurableExecutionResponse> getDurableExecution(
      $grpc.ServiceCall call, $0.GetDurableExecutionRequest request);

  $async.Future<$0.GetDurableExecutionHistoryResponse>
      getDurableExecutionHistory_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetDurableExecutionHistoryRequest> $request) async {
    return getDurableExecutionHistory($call, await $request);
  }

  $async.Future<$0.GetDurableExecutionHistoryResponse>
      getDurableExecutionHistory(
          $grpc.ServiceCall call, $0.GetDurableExecutionHistoryRequest request);

  $async.Future<$0.GetDurableExecutionStateResponse>
      getDurableExecutionState_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetDurableExecutionStateRequest> $request) async {
    return getDurableExecutionState($call, await $request);
  }

  $async.Future<$0.GetDurableExecutionStateResponse> getDurableExecutionState(
      $grpc.ServiceCall call, $0.GetDurableExecutionStateRequest request);

  $async.Future<$0.FunctionEventInvokeConfig> getFunctionEventInvokeConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetFunctionEventInvokeConfigRequest> $request) async {
    return getFunctionEventInvokeConfig($call, await $request);
  }

  $async.Future<$0.FunctionEventInvokeConfig> getFunctionEventInvokeConfig(
      $grpc.ServiceCall call, $0.GetFunctionEventInvokeConfigRequest request);

  $async.Future<$0.ListDurableExecutionsByFunctionResponse>
      listDurableExecutionsByFunction_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListDurableExecutionsByFunctionRequest>
              $request) async {
    return listDurableExecutionsByFunction($call, await $request);
  }

  $async.Future<$0.ListDurableExecutionsByFunctionResponse>
      listDurableExecutionsByFunction($grpc.ServiceCall call,
          $0.ListDurableExecutionsByFunctionRequest request);

  $async.Future<$0.ListFunctionEventInvokeConfigsResponse>
      listFunctionEventInvokeConfigs_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListFunctionEventInvokeConfigsRequest>
              $request) async {
    return listFunctionEventInvokeConfigs($call, await $request);
  }

  $async.Future<$0.ListFunctionEventInvokeConfigsResponse>
      listFunctionEventInvokeConfigs($grpc.ServiceCall call,
          $0.ListFunctionEventInvokeConfigsRequest request);

  $async.Future<$0.ListTagsResponse> listTags_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTagsRequest> $request) async {
    return listTags($call, await $request);
  }

  $async.Future<$0.ListTagsResponse> listTags(
      $grpc.ServiceCall call, $0.ListTagsRequest request);

  $async.Future<$0.FunctionEventInvokeConfig> putFunctionEventInvokeConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutFunctionEventInvokeConfigRequest> $request) async {
    return putFunctionEventInvokeConfig($call, await $request);
  }

  $async.Future<$0.FunctionEventInvokeConfig> putFunctionEventInvokeConfig(
      $grpc.ServiceCall call, $0.PutFunctionEventInvokeConfigRequest request);

  $async.Future<$0.SendDurableExecutionCallbackFailureResponse>
      sendDurableExecutionCallbackFailure_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.SendDurableExecutionCallbackFailureRequest>
              $request) async {
    return sendDurableExecutionCallbackFailure($call, await $request);
  }

  $async.Future<$0.SendDurableExecutionCallbackFailureResponse>
      sendDurableExecutionCallbackFailure($grpc.ServiceCall call,
          $0.SendDurableExecutionCallbackFailureRequest request);

  $async.Future<$0.SendDurableExecutionCallbackHeartbeatResponse>
      sendDurableExecutionCallbackHeartbeat_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.SendDurableExecutionCallbackHeartbeatRequest>
              $request) async {
    return sendDurableExecutionCallbackHeartbeat($call, await $request);
  }

  $async.Future<$0.SendDurableExecutionCallbackHeartbeatResponse>
      sendDurableExecutionCallbackHeartbeat($grpc.ServiceCall call,
          $0.SendDurableExecutionCallbackHeartbeatRequest request);

  $async.Future<$0.SendDurableExecutionCallbackSuccessResponse>
      sendDurableExecutionCallbackSuccess_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.SendDurableExecutionCallbackSuccessRequest>
              $request) async {
    return sendDurableExecutionCallbackSuccess($call, await $request);
  }

  $async.Future<$0.SendDurableExecutionCallbackSuccessResponse>
      sendDurableExecutionCallbackSuccess($grpc.ServiceCall call,
          $0.SendDurableExecutionCallbackSuccessRequest request);

  $async.Future<$0.StopDurableExecutionResponse> stopDurableExecution_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StopDurableExecutionRequest> $request) async {
    return stopDurableExecution($call, await $request);
  }

  $async.Future<$0.StopDurableExecutionResponse> stopDurableExecution(
      $grpc.ServiceCall call, $0.StopDurableExecutionRequest request);

  $async.Future<$1.Empty> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$1.Empty> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$1.Empty> untagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$1.Empty> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.FunctionEventInvokeConfig>
      updateFunctionEventInvokeConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.UpdateFunctionEventInvokeConfigRequest>
              $request) async {
    return updateFunctionEventInvokeConfig($call, await $request);
  }

  $async.Future<$0.FunctionEventInvokeConfig> updateFunctionEventInvokeConfig(
      $grpc.ServiceCall call,
      $0.UpdateFunctionEventInvokeConfigRequest request);

  $async.Future<$0.ListCapacityProvidersResponse> listCapacityProviders_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCapacityProvidersRequest> $request) async {
    return listCapacityProviders($call, await $request);
  }

  $async.Future<$0.ListCapacityProvidersResponse> listCapacityProviders(
      $grpc.ServiceCall call, $0.ListCapacityProvidersRequest request);

  $async.Future<$0.CreateCapacityProviderResponse> createCapacityProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateCapacityProviderRequest> $request) async {
    return createCapacityProvider($call, await $request);
  }

  $async.Future<$0.CreateCapacityProviderResponse> createCapacityProvider(
      $grpc.ServiceCall call, $0.CreateCapacityProviderRequest request);

  $async.Future<$0.GetCapacityProviderResponse> getCapacityProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCapacityProviderRequest> $request) async {
    return getCapacityProvider($call, await $request);
  }

  $async.Future<$0.GetCapacityProviderResponse> getCapacityProvider(
      $grpc.ServiceCall call, $0.GetCapacityProviderRequest request);

  $async.Future<$0.UpdateCapacityProviderResponse> updateCapacityProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateCapacityProviderRequest> $request) async {
    return updateCapacityProvider($call, await $request);
  }

  $async.Future<$0.UpdateCapacityProviderResponse> updateCapacityProvider(
      $grpc.ServiceCall call, $0.UpdateCapacityProviderRequest request);

  $async.Future<$0.DeleteCapacityProviderResponse> deleteCapacityProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteCapacityProviderRequest> $request) async {
    return deleteCapacityProvider($call, await $request);
  }

  $async.Future<$0.DeleteCapacityProviderResponse> deleteCapacityProvider(
      $grpc.ServiceCall call, $0.DeleteCapacityProviderRequest request);

  $async.Future<$0.ListFunctionVersionsByCapacityProviderResponse>
      listFunctionVersionsByCapacityProvider_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListFunctionVersionsByCapacityProviderRequest>
              $request) async {
    return listFunctionVersionsByCapacityProvider($call, await $request);
  }

  $async.Future<$0.ListFunctionVersionsByCapacityProviderResponse>
      listFunctionVersionsByCapacityProvider($grpc.ServiceCall call,
          $0.ListFunctionVersionsByCapacityProviderRequest request);

  $async.Future<$0.ListCodeSigningConfigsResponse> listCodeSigningConfigs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCodeSigningConfigsRequest> $request) async {
    return listCodeSigningConfigs($call, await $request);
  }

  $async.Future<$0.ListCodeSigningConfigsResponse> listCodeSigningConfigs(
      $grpc.ServiceCall call, $0.ListCodeSigningConfigsRequest request);

  $async.Future<$0.CreateCodeSigningConfigResponse> createCodeSigningConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateCodeSigningConfigRequest> $request) async {
    return createCodeSigningConfig($call, await $request);
  }

  $async.Future<$0.CreateCodeSigningConfigResponse> createCodeSigningConfig(
      $grpc.ServiceCall call, $0.CreateCodeSigningConfigRequest request);

  $async.Future<$0.ListEventSourceMappingsResponse> listEventSourceMappings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEventSourceMappingsRequest> $request) async {
    return listEventSourceMappings($call, await $request);
  }

  $async.Future<$0.ListEventSourceMappingsResponse> listEventSourceMappings(
      $grpc.ServiceCall call, $0.ListEventSourceMappingsRequest request);

  $async.Future<$0.EventSourceMappingConfiguration>
      createEventSourceMapping_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateEventSourceMappingRequest> $request) async {
    return createEventSourceMapping($call, await $request);
  }

  $async.Future<$0.EventSourceMappingConfiguration> createEventSourceMapping(
      $grpc.ServiceCall call, $0.CreateEventSourceMappingRequest request);

  $async.Future<$0.EventSourceMappingConfiguration> getEventSourceMapping_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetEventSourceMappingRequest> $request) async {
    return getEventSourceMapping($call, await $request);
  }

  $async.Future<$0.EventSourceMappingConfiguration> getEventSourceMapping(
      $grpc.ServiceCall call, $0.GetEventSourceMappingRequest request);

  $async.Future<$0.EventSourceMappingConfiguration>
      updateEventSourceMapping_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateEventSourceMappingRequest> $request) async {
    return updateEventSourceMapping($call, await $request);
  }

  $async.Future<$0.EventSourceMappingConfiguration> updateEventSourceMapping(
      $grpc.ServiceCall call, $0.UpdateEventSourceMappingRequest request);

  $async.Future<$0.EventSourceMappingConfiguration>
      deleteEventSourceMapping_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteEventSourceMappingRequest> $request) async {
    return deleteEventSourceMapping($call, await $request);
  }

  $async.Future<$0.EventSourceMappingConfiguration> deleteEventSourceMapping(
      $grpc.ServiceCall call, $0.DeleteEventSourceMappingRequest request);

  $async.Future<$0.ListFunctionsResponse> listFunctions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListFunctionsRequest> $request) async {
    return listFunctions($call, await $request);
  }

  $async.Future<$0.ListFunctionsResponse> listFunctions(
      $grpc.ServiceCall call, $0.ListFunctionsRequest request);

  $async.Future<$0.FunctionConfiguration> createFunction_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateFunctionRequest> $request) async {
    return createFunction($call, await $request);
  }

  $async.Future<$0.FunctionConfiguration> createFunction(
      $grpc.ServiceCall call, $0.CreateFunctionRequest request);

  $async.Future<$0.CreateFunctionUrlConfigResponse> createFunctionUrlConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateFunctionUrlConfigRequest> $request) async {
    return createFunctionUrlConfig($call, await $request);
  }

  $async.Future<$0.CreateFunctionUrlConfigResponse> createFunctionUrlConfig(
      $grpc.ServiceCall call, $0.CreateFunctionUrlConfigRequest request);

  $async.Future<$1.Empty> deleteFunctionConcurrency_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteFunctionConcurrencyRequest> $request) async {
    return deleteFunctionConcurrency($call, await $request);
  }

  $async.Future<$1.Empty> deleteFunctionConcurrency(
      $grpc.ServiceCall call, $0.DeleteFunctionConcurrencyRequest request);

  $async.Future<$1.Empty> deleteFunctionUrlConfig_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteFunctionUrlConfigRequest> $request) async {
    return deleteFunctionUrlConfig($call, await $request);
  }

  $async.Future<$1.Empty> deleteFunctionUrlConfig(
      $grpc.ServiceCall call, $0.DeleteFunctionUrlConfigRequest request);

  $async.Future<$0.GetFunctionConcurrencyResponse> getFunctionConcurrency_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetFunctionConcurrencyRequest> $request) async {
    return getFunctionConcurrency($call, await $request);
  }

  $async.Future<$0.GetFunctionConcurrencyResponse> getFunctionConcurrency(
      $grpc.ServiceCall call, $0.GetFunctionConcurrencyRequest request);

  $async.Future<$0.GetFunctionUrlConfigResponse> getFunctionUrlConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetFunctionUrlConfigRequest> $request) async {
    return getFunctionUrlConfig($call, await $request);
  }

  $async.Future<$0.GetFunctionUrlConfigResponse> getFunctionUrlConfig(
      $grpc.ServiceCall call, $0.GetFunctionUrlConfigRequest request);

  $async.Future<$0.ListFunctionUrlConfigsResponse> listFunctionUrlConfigs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListFunctionUrlConfigsRequest> $request) async {
    return listFunctionUrlConfigs($call, await $request);
  }

  $async.Future<$0.ListFunctionUrlConfigsResponse> listFunctionUrlConfigs(
      $grpc.ServiceCall call, $0.ListFunctionUrlConfigsRequest request);

  $async.Future<$0.ListProvisionedConcurrencyConfigsResponse>
      listProvisionedConcurrencyConfigs_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListProvisionedConcurrencyConfigsRequest>
              $request) async {
    return listProvisionedConcurrencyConfigs($call, await $request);
  }

  $async.Future<$0.ListProvisionedConcurrencyConfigsResponse>
      listProvisionedConcurrencyConfigs($grpc.ServiceCall call,
          $0.ListProvisionedConcurrencyConfigsRequest request);

  $async.Future<$0.Concurrency> putFunctionConcurrency_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutFunctionConcurrencyRequest> $request) async {
    return putFunctionConcurrency($call, await $request);
  }

  $async.Future<$0.Concurrency> putFunctionConcurrency(
      $grpc.ServiceCall call, $0.PutFunctionConcurrencyRequest request);

  $async.Future<$0.FunctionConfiguration> updateFunctionCode_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateFunctionCodeRequest> $request) async {
    return updateFunctionCode($call, await $request);
  }

  $async.Future<$0.FunctionConfiguration> updateFunctionCode(
      $grpc.ServiceCall call, $0.UpdateFunctionCodeRequest request);

  $async.Future<$0.FunctionConfiguration> updateFunctionConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateFunctionConfigurationRequest> $request) async {
    return updateFunctionConfiguration($call, await $request);
  }

  $async.Future<$0.FunctionConfiguration> updateFunctionConfiguration(
      $grpc.ServiceCall call, $0.UpdateFunctionConfigurationRequest request);

  $async.Future<$0.UpdateFunctionUrlConfigResponse> updateFunctionUrlConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateFunctionUrlConfigRequest> $request) async {
    return updateFunctionUrlConfig($call, await $request);
  }

  $async.Future<$0.UpdateFunctionUrlConfigResponse> updateFunctionUrlConfig(
      $grpc.ServiceCall call, $0.UpdateFunctionUrlConfigRequest request);

  $async.Future<$0.ListAliasesResponse> listAliases_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListAliasesRequest> $request) async {
    return listAliases($call, await $request);
  }

  $async.Future<$0.ListAliasesResponse> listAliases(
      $grpc.ServiceCall call, $0.ListAliasesRequest request);

  $async.Future<$0.AliasConfiguration> createAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateAliasRequest> $request) async {
    return createAlias($call, await $request);
  }

  $async.Future<$0.AliasConfiguration> createAlias(
      $grpc.ServiceCall call, $0.CreateAliasRequest request);

  $async.Future<$0.AliasConfiguration> getAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetAliasRequest> $request) async {
    return getAlias($call, await $request);
  }

  $async.Future<$0.AliasConfiguration> getAlias(
      $grpc.ServiceCall call, $0.GetAliasRequest request);

  $async.Future<$0.AliasConfiguration> updateAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAliasRequest> $request) async {
    return updateAlias($call, await $request);
  }

  $async.Future<$0.AliasConfiguration> updateAlias(
      $grpc.ServiceCall call, $0.UpdateAliasRequest request);

  $async.Future<$1.Empty> deleteAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAliasRequest> $request) async {
    return deleteAlias($call, await $request);
  }

  $async.Future<$1.Empty> deleteAlias(
      $grpc.ServiceCall call, $0.DeleteAliasRequest request);

  $async.Future<$0.FunctionConfiguration> publishVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PublishVersionRequest> $request) async {
    return publishVersion($call, await $request);
  }

  $async.Future<$0.FunctionConfiguration> publishVersion(
      $grpc.ServiceCall call, $0.PublishVersionRequest request);

  $async.Future<$0.ListLayersResponse> listLayers_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListLayersRequest> $request) async {
    return listLayers($call, await $request);
  }

  $async.Future<$0.ListLayersResponse> listLayers(
      $grpc.ServiceCall call, $0.ListLayersRequest request);

  $async.Future<$0.ListLayerVersionsResponse> listLayerVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListLayerVersionsRequest> $request) async {
    return listLayerVersions($call, await $request);
  }

  $async.Future<$0.ListLayerVersionsResponse> listLayerVersions(
      $grpc.ServiceCall call, $0.ListLayerVersionsRequest request);

  $async.Future<$0.PutProvisionedConcurrencyConfigResponse>
      putProvisionedConcurrencyConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.PutProvisionedConcurrencyConfigRequest>
              $request) async {
    return putProvisionedConcurrencyConfig($call, await $request);
  }

  $async.Future<$0.PutProvisionedConcurrencyConfigResponse>
      putProvisionedConcurrencyConfig($grpc.ServiceCall call,
          $0.PutProvisionedConcurrencyConfigRequest request);

  $async.Future<$0.GetProvisionedConcurrencyConfigResponse>
      getProvisionedConcurrencyConfig_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetProvisionedConcurrencyConfigRequest>
              $request) async {
    return getProvisionedConcurrencyConfig($call, await $request);
  }

  $async.Future<$0.GetProvisionedConcurrencyConfigResponse>
      getProvisionedConcurrencyConfig($grpc.ServiceCall call,
          $0.GetProvisionedConcurrencyConfigRequest request);

  $async.Future<$1.Empty> deleteProvisionedConcurrencyConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteProvisionedConcurrencyConfigRequest>
          $request) async {
    return deleteProvisionedConcurrencyConfig($call, await $request);
  }

  $async.Future<$1.Empty> deleteProvisionedConcurrencyConfig(
      $grpc.ServiceCall call,
      $0.DeleteProvisionedConcurrencyConfigRequest request);
}
