// This is a generated file - do not edit.
//
// Generated from acm.proto.

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

import 'acm.pb.dart' as $0;
import 'common.pb.dart' as $1;

export 'acm.pb.dart';

/// ACMService provides acm API operations.
@$pb.GrpcServiceName('acm.ACMService')
class ACMServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  ACMServiceClient(super.channel, {super.options, super.interceptors});

  /// Adds one or more tags to an ACM certificate. Tags are labels that you can use to identify and organize your Amazon Web Services resources. Each tag consists of a key and an optional value. You spec...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> addTagsToCertificate(
    $0.AddTagsToCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addTagsToCertificate, request, options: options);
  }

  /// Deletes a certificate and its associated private key. If this action succeeds, the certificate no longer appears in the list that can be displayed by calling the ListCertificates action or be retri...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteCertificate(
    $0.DeleteCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteCertificate, request, options: options);
  }

  /// Returns detailed metadata about the specified ACM certificate. If you have just created a certificate using the RequestCertificate action, there is a delay of several seconds before you can retriev...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeCertificateResponse> describeCertificate(
    $0.DescribeCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeCertificate, request, options: options);
  }

  /// Exports a private certificate issued by a private certificate authority (CA) or public certificate for use anywhere. The exported file contains the certificate, the certificate chain, and the encry...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ExportCertificateResponse> exportCertificate(
    $0.ExportCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$exportCertificate, request, options: options);
  }

  /// Returns the account configuration options associated with an Amazon Web Services account.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetAccountConfigurationResponse>
      getAccountConfiguration(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccountConfiguration, request,
        options: options);
  }

  /// Retrieves a certificate and its certificate chain. The certificate may be either a public or private certificate issued using the ACM RequestCertificate action, or a certificate imported into ACM u...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCertificateResponse> getCertificate(
    $0.GetCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCertificate, request, options: options);
  }

  /// Imports a certificate into Certificate Manager (ACM) to use with services that are integrated with ACM. Note that integrated services allow only certificate types and keys they support to be associ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ImportCertificateResponse> importCertificate(
    $0.ImportCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$importCertificate, request, options: options);
  }

  /// Retrieves a list of certificate ARNs and domain names. You can request that only certificates that match a specific status be listed. You can also filter by specific attributes of the certificate. ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListCertificatesResponse> listCertificates(
    $0.ListCertificatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listCertificates, request, options: options);
  }

  /// Lists the tags that have been applied to the ACM certificate. Use the certificate's Amazon Resource Name (ARN) to specify the certificate. To add a tag to an ACM certificate, use the AddTagsToCerti...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForCertificateResponse>
      listTagsForCertificate(
    $0.ListTagsForCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForCertificate, request,
        options: options);
  }

  /// Adds or modifies account-level configurations in ACM. The supported configuration option is DaysBeforeExpiry. This option specifies the number of days prior to certificate expiration when ACM start...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> putAccountConfiguration(
    $0.PutAccountConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putAccountConfiguration, request,
        options: options);
  }

  /// Remove one or more tags from an ACM certificate. A tag consists of a key-value pair. If you do not specify the value portion of the tag when calling this function, the tag will be removed regardles...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> removeTagsFromCertificate(
    $0.RemoveTagsFromCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeTagsFromCertificate, request,
        options: options);
  }

  /// Renews an eligible ACM certificate. In order to renew your Amazon Web Services Private CA certificates with ACM, you must first grant the ACM service principal permission to do so. For more informa...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> renewCertificate(
    $0.RenewCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$renewCertificate, request, options: options);
  }

  /// Requests an ACM certificate for use with other Amazon Web Services services. To request an ACM certificate, you must specify a fully qualified domain name (FQDN) in the DomainName parameter. You ca...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RequestCertificateResponse> requestCertificate(
    $0.RequestCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$requestCertificate, request, options: options);
  }

  /// Resends the email that requests domain ownership validation. The domain owner or an authorized representative must approve the ACM certificate before it can be issued. The certificate can be approv...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> resendValidationEmail(
    $0.ResendValidationEmailRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resendValidationEmail, request, options: options);
  }

  /// Revokes a public ACM certificate. You can only revoke certificates that have been previously exported.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RevokeCertificateResponse> revokeCertificate(
    $0.RevokeCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$revokeCertificate, request, options: options);
  }

  /// Updates a certificate. You can use this function to specify whether to opt in to or out of recording your certificate in a certificate transparency log and exporting. For more information, see Opti...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> updateCertificateOptions(
    $0.UpdateCertificateOptionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateCertificateOptions, request,
        options: options);
  }

  // method descriptors

  static final _$addTagsToCertificate =
      $grpc.ClientMethod<$0.AddTagsToCertificateRequest, $1.Empty>(
          '/acm.ACMService/AddTagsToCertificate',
          ($0.AddTagsToCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteCertificate =
      $grpc.ClientMethod<$0.DeleteCertificateRequest, $1.Empty>(
          '/acm.ACMService/DeleteCertificate',
          ($0.DeleteCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$describeCertificate = $grpc.ClientMethod<
          $0.DescribeCertificateRequest, $0.DescribeCertificateResponse>(
      '/acm.ACMService/DescribeCertificate',
      ($0.DescribeCertificateRequest value) => value.writeToBuffer(),
      $0.DescribeCertificateResponse.fromBuffer);
  static final _$exportCertificate = $grpc.ClientMethod<
          $0.ExportCertificateRequest, $0.ExportCertificateResponse>(
      '/acm.ACMService/ExportCertificate',
      ($0.ExportCertificateRequest value) => value.writeToBuffer(),
      $0.ExportCertificateResponse.fromBuffer);
  static final _$getAccountConfiguration =
      $grpc.ClientMethod<$1.Empty, $0.GetAccountConfigurationResponse>(
          '/acm.ACMService/GetAccountConfiguration',
          ($1.Empty value) => value.writeToBuffer(),
          $0.GetAccountConfigurationResponse.fromBuffer);
  static final _$getCertificate =
      $grpc.ClientMethod<$0.GetCertificateRequest, $0.GetCertificateResponse>(
          '/acm.ACMService/GetCertificate',
          ($0.GetCertificateRequest value) => value.writeToBuffer(),
          $0.GetCertificateResponse.fromBuffer);
  static final _$importCertificate = $grpc.ClientMethod<
          $0.ImportCertificateRequest, $0.ImportCertificateResponse>(
      '/acm.ACMService/ImportCertificate',
      ($0.ImportCertificateRequest value) => value.writeToBuffer(),
      $0.ImportCertificateResponse.fromBuffer);
  static final _$listCertificates = $grpc.ClientMethod<
          $0.ListCertificatesRequest, $0.ListCertificatesResponse>(
      '/acm.ACMService/ListCertificates',
      ($0.ListCertificatesRequest value) => value.writeToBuffer(),
      $0.ListCertificatesResponse.fromBuffer);
  static final _$listTagsForCertificate = $grpc.ClientMethod<
          $0.ListTagsForCertificateRequest, $0.ListTagsForCertificateResponse>(
      '/acm.ACMService/ListTagsForCertificate',
      ($0.ListTagsForCertificateRequest value) => value.writeToBuffer(),
      $0.ListTagsForCertificateResponse.fromBuffer);
  static final _$putAccountConfiguration =
      $grpc.ClientMethod<$0.PutAccountConfigurationRequest, $1.Empty>(
          '/acm.ACMService/PutAccountConfiguration',
          ($0.PutAccountConfigurationRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$removeTagsFromCertificate =
      $grpc.ClientMethod<$0.RemoveTagsFromCertificateRequest, $1.Empty>(
          '/acm.ACMService/RemoveTagsFromCertificate',
          ($0.RemoveTagsFromCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$renewCertificate =
      $grpc.ClientMethod<$0.RenewCertificateRequest, $1.Empty>(
          '/acm.ACMService/RenewCertificate',
          ($0.RenewCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$requestCertificate = $grpc.ClientMethod<
          $0.RequestCertificateRequest, $0.RequestCertificateResponse>(
      '/acm.ACMService/RequestCertificate',
      ($0.RequestCertificateRequest value) => value.writeToBuffer(),
      $0.RequestCertificateResponse.fromBuffer);
  static final _$resendValidationEmail =
      $grpc.ClientMethod<$0.ResendValidationEmailRequest, $1.Empty>(
          '/acm.ACMService/ResendValidationEmail',
          ($0.ResendValidationEmailRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$revokeCertificate = $grpc.ClientMethod<
          $0.RevokeCertificateRequest, $0.RevokeCertificateResponse>(
      '/acm.ACMService/RevokeCertificate',
      ($0.RevokeCertificateRequest value) => value.writeToBuffer(),
      $0.RevokeCertificateResponse.fromBuffer);
  static final _$updateCertificateOptions =
      $grpc.ClientMethod<$0.UpdateCertificateOptionsRequest, $1.Empty>(
          '/acm.ACMService/UpdateCertificateOptions',
          ($0.UpdateCertificateOptionsRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
}

@$pb.GrpcServiceName('acm.ACMService')
abstract class ACMServiceBase extends $grpc.Service {
  $core.String get $name => 'acm.ACMService';

  ACMServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AddTagsToCertificateRequest, $1.Empty>(
        'AddTagsToCertificate',
        addTagsToCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddTagsToCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteCertificateRequest, $1.Empty>(
        'DeleteCertificate',
        deleteCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeCertificateRequest,
            $0.DescribeCertificateResponse>(
        'DescribeCertificate',
        describeCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeCertificateRequest.fromBuffer(value),
        ($0.DescribeCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ExportCertificateRequest,
            $0.ExportCertificateResponse>(
        'ExportCertificate',
        exportCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ExportCertificateRequest.fromBuffer(value),
        ($0.ExportCertificateResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$1.Empty, $0.GetAccountConfigurationResponse>(
            'GetAccountConfiguration',
            getAccountConfiguration_Pre,
            false,
            false,
            ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
            ($0.GetAccountConfigurationResponse value) =>
                value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetCertificateRequest,
            $0.GetCertificateResponse>(
        'GetCertificate',
        getCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetCertificateRequest.fromBuffer(value),
        ($0.GetCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ImportCertificateRequest,
            $0.ImportCertificateResponse>(
        'ImportCertificate',
        importCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ImportCertificateRequest.fromBuffer(value),
        ($0.ImportCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListCertificatesRequest,
            $0.ListCertificatesResponse>(
        'ListCertificates',
        listCertificates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListCertificatesRequest.fromBuffer(value),
        ($0.ListCertificatesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForCertificateRequest,
            $0.ListTagsForCertificateResponse>(
        'ListTagsForCertificate',
        listTagsForCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForCertificateRequest.fromBuffer(value),
        ($0.ListTagsForCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutAccountConfigurationRequest, $1.Empty>(
        'PutAccountConfiguration',
        putAccountConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutAccountConfigurationRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RemoveTagsFromCertificateRequest, $1.Empty>(
            'RemoveTagsFromCertificate',
            removeTagsFromCertificate_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RemoveTagsFromCertificateRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RenewCertificateRequest, $1.Empty>(
        'RenewCertificate',
        renewCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RenewCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RequestCertificateRequest,
            $0.RequestCertificateResponse>(
        'RequestCertificate',
        requestCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RequestCertificateRequest.fromBuffer(value),
        ($0.RequestCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResendValidationEmailRequest, $1.Empty>(
        'ResendValidationEmail',
        resendValidationEmail_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResendValidationEmailRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RevokeCertificateRequest,
            $0.RevokeCertificateResponse>(
        'RevokeCertificate',
        revokeCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RevokeCertificateRequest.fromBuffer(value),
        ($0.RevokeCertificateResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateCertificateOptionsRequest, $1.Empty>(
            'UpdateCertificateOptions',
            updateCertificateOptions_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateCertificateOptionsRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> addTagsToCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddTagsToCertificateRequest> $request) async {
    return addTagsToCertificate($call, await $request);
  }

  $async.Future<$1.Empty> addTagsToCertificate(
      $grpc.ServiceCall call, $0.AddTagsToCertificateRequest request);

  $async.Future<$1.Empty> deleteCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteCertificateRequest> $request) async {
    return deleteCertificate($call, await $request);
  }

  $async.Future<$1.Empty> deleteCertificate(
      $grpc.ServiceCall call, $0.DeleteCertificateRequest request);

  $async.Future<$0.DescribeCertificateResponse> describeCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeCertificateRequest> $request) async {
    return describeCertificate($call, await $request);
  }

  $async.Future<$0.DescribeCertificateResponse> describeCertificate(
      $grpc.ServiceCall call, $0.DescribeCertificateRequest request);

  $async.Future<$0.ExportCertificateResponse> exportCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ExportCertificateRequest> $request) async {
    return exportCertificate($call, await $request);
  }

  $async.Future<$0.ExportCertificateResponse> exportCertificate(
      $grpc.ServiceCall call, $0.ExportCertificateRequest request);

  $async.Future<$0.GetAccountConfigurationResponse> getAccountConfiguration_Pre(
      $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return getAccountConfiguration($call, await $request);
  }

  $async.Future<$0.GetAccountConfigurationResponse> getAccountConfiguration(
      $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$0.GetCertificateResponse> getCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCertificateRequest> $request) async {
    return getCertificate($call, await $request);
  }

  $async.Future<$0.GetCertificateResponse> getCertificate(
      $grpc.ServiceCall call, $0.GetCertificateRequest request);

  $async.Future<$0.ImportCertificateResponse> importCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ImportCertificateRequest> $request) async {
    return importCertificate($call, await $request);
  }

  $async.Future<$0.ImportCertificateResponse> importCertificate(
      $grpc.ServiceCall call, $0.ImportCertificateRequest request);

  $async.Future<$0.ListCertificatesResponse> listCertificates_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListCertificatesRequest> $request) async {
    return listCertificates($call, await $request);
  }

  $async.Future<$0.ListCertificatesResponse> listCertificates(
      $grpc.ServiceCall call, $0.ListCertificatesRequest request);

  $async.Future<$0.ListTagsForCertificateResponse> listTagsForCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForCertificateRequest> $request) async {
    return listTagsForCertificate($call, await $request);
  }

  $async.Future<$0.ListTagsForCertificateResponse> listTagsForCertificate(
      $grpc.ServiceCall call, $0.ListTagsForCertificateRequest request);

  $async.Future<$1.Empty> putAccountConfiguration_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutAccountConfigurationRequest> $request) async {
    return putAccountConfiguration($call, await $request);
  }

  $async.Future<$1.Empty> putAccountConfiguration(
      $grpc.ServiceCall call, $0.PutAccountConfigurationRequest request);

  $async.Future<$1.Empty> removeTagsFromCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemoveTagsFromCertificateRequest> $request) async {
    return removeTagsFromCertificate($call, await $request);
  }

  $async.Future<$1.Empty> removeTagsFromCertificate(
      $grpc.ServiceCall call, $0.RemoveTagsFromCertificateRequest request);

  $async.Future<$1.Empty> renewCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RenewCertificateRequest> $request) async {
    return renewCertificate($call, await $request);
  }

  $async.Future<$1.Empty> renewCertificate(
      $grpc.ServiceCall call, $0.RenewCertificateRequest request);

  $async.Future<$0.RequestCertificateResponse> requestCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RequestCertificateRequest> $request) async {
    return requestCertificate($call, await $request);
  }

  $async.Future<$0.RequestCertificateResponse> requestCertificate(
      $grpc.ServiceCall call, $0.RequestCertificateRequest request);

  $async.Future<$1.Empty> resendValidationEmail_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ResendValidationEmailRequest> $request) async {
    return resendValidationEmail($call, await $request);
  }

  $async.Future<$1.Empty> resendValidationEmail(
      $grpc.ServiceCall call, $0.ResendValidationEmailRequest request);

  $async.Future<$0.RevokeCertificateResponse> revokeCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RevokeCertificateRequest> $request) async {
    return revokeCertificate($call, await $request);
  }

  $async.Future<$0.RevokeCertificateResponse> revokeCertificate(
      $grpc.ServiceCall call, $0.RevokeCertificateRequest request);

  $async.Future<$1.Empty> updateCertificateOptions_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateCertificateOptionsRequest> $request) async {
    return updateCertificateOptions($call, await $request);
  }

  $async.Future<$1.Empty> updateCertificateOptions(
      $grpc.ServiceCall call, $0.UpdateCertificateOptionsRequest request);
}
