// This is a generated file - do not edit.
//
// Generated from cognito_idp.proto.

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

import 'cognito_idp.pb.dart' as $0;
import 'common.pb.dart' as $1;

export 'cognito_idp.pb.dart';

/// CognitoIdentityProviderService provides cognito_idp API operations.
@$pb.GrpcServiceName('cognito_idp.CognitoIdentityProviderService')
class CognitoIdentityProviderServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  CognitoIdentityProviderServiceClient(super.channel,
      {super.options, super.interceptors});

  /// Adds additional user attributes to the user pool schema. Custom attributes can be mutable or immutable and have a custom: or dev: prefix. For more information, see Custom attributes. Amazon Cognito...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AddCustomAttributesResponse> addCustomAttributes(
    $0.AddCustomAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addCustomAttributes, request, options: options);
  }

  /// Creates a new client secret for an existing confidential user pool app client. Supports up to 2 active secrets per app client for zero-downtime credential rotation workflows.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AddUserPoolClientSecretResponse>
      addUserPoolClientSecret(
    $0.AddUserPoolClientSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addUserPoolClientSecret, request,
        options: options);
  }

  /// Adds a user to a group. A user who is in a group can present a preferred-role claim to an identity pool, and populates a cognito:groups claim to their access and identity tokens. Amazon Cognito eva...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> adminAddUserToGroup(
    $0.AdminAddUserToGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminAddUserToGroup, request, options: options);
  }

  /// Confirms user sign-up as an administrator. This request sets a user account active in a user pool that requires confirmation of new user accounts before they can sign in. You can configure your use...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminConfirmSignUpResponse> adminConfirmSignUp(
    $0.AdminConfirmSignUpRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminConfirmSignUp, request, options: options);
  }

  /// Creates a new user in the specified user pool. If MessageAction isn't set, the default is to send a welcome message via email or phone (SMS). This message is based on a template that you configured...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminCreateUserResponse> adminCreateUser(
    $0.AdminCreateUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminCreateUser, request, options: options);
  }

  /// Deletes a user profile in your user pool. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API operation. For this operation, you must use IAM credentials...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> adminDeleteUser(
    $0.AdminDeleteUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminDeleteUser, request, options: options);
  }

  /// Deletes attribute values from a user. This operation doesn't affect tokens for existing user sessions. The next ID token that the user receives will no longer have the deleted attributes. Amazon Co...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminDeleteUserAttributesResponse>
      adminDeleteUserAttributes(
    $0.AdminDeleteUserAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminDeleteUserAttributes, request,
        options: options);
  }

  /// Prevents the user from signing in with the specified external (SAML or social) identity provider (IdP). If the user that you want to deactivate is a Amazon Cognito user pools native username + pass...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminDisableProviderForUserResponse>
      adminDisableProviderForUser(
    $0.AdminDisableProviderForUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminDisableProviderForUser, request,
        options: options);
  }

  /// Deactivates a user profile and revokes all access tokens for the user. A deactivated user can't sign in, but still appears in the responses to ListUsers API requests. Amazon Cognito evaluates Ident...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminDisableUserResponse> adminDisableUser(
    $0.AdminDisableUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminDisableUser, request, options: options);
  }

  /// Activates sign-in for a user profile that previously had sign-in access disabled. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API operation. For this...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminEnableUserResponse> adminEnableUser(
    $0.AdminEnableUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminEnableUser, request, options: options);
  }

  /// Forgets, or deletes, a remembered device from a user's profile. After you forget the device, the user can no longer complete device authentication with that device and when applicable, must submit ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> adminForgetDevice(
    $0.AdminForgetDeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminForgetDevice, request, options: options);
  }

  /// Given the device key, returns details for a user's device. For more information, see Working with devices. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for thi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminGetDeviceResponse> adminGetDevice(
    $0.AdminGetDeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminGetDevice, request, options: options);
  }

  /// Given a username, returns details about a user profile in a user pool. You can specify alias attributes in the Username request parameter. This operation contributes to your monthly active user (MA...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminGetUserResponse> adminGetUser(
    $0.AdminGetUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminGetUser, request, options: options);
  }

  /// Starts sign-in for applications with a server-side component, for example a traditional web application. This operation specifies the authentication flow that you'd like to begin. The authenticatio...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminInitiateAuthResponse> adminInitiateAuth(
    $0.AdminInitiateAuthRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminInitiateAuth, request, options: options);
  }

  /// Links an existing user account in a user pool, or DestinationUser, to an identity from an external IdP, or SourceUser, based on a specified attribute name and value from the external IdP. This oper...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminLinkProviderForUserResponse>
      adminLinkProviderForUser(
    $0.AdminLinkProviderForUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminLinkProviderForUser, request,
        options: options);
  }

  /// Lists a user's registered devices. Remembered devices are used in authentication services where you offer a "Remember me" option for users who you want to permit to sign in without MFA from a trust...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminListDevicesResponse> adminListDevices(
    $0.AdminListDevicesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminListDevices, request, options: options);
  }

  /// Lists the groups that a user belongs to. User pool groups are identifiers that you can reference from the contents of ID and access tokens, and set preferred IAM roles for identity-pool authenticat...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminListGroupsForUserResponse>
      adminListGroupsForUser(
    $0.AdminListGroupsForUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminListGroupsForUser, request,
        options: options);
  }

  /// Requests a history of user activity and any risks detected as part of Amazon Cognito threat protection. For more information, see Viewing user event history. Amazon Cognito evaluates Identity and A...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminListUserAuthEventsResponse>
      adminListUserAuthEvents(
    $0.AdminListUserAuthEventsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminListUserAuthEvents, request,
        options: options);
  }

  /// Given a username and a group name, removes them from the group. User pool groups are identifiers that you can reference from the contents of ID and access tokens, and set preferred IAM roles for id...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> adminRemoveUserFromGroup(
    $0.AdminRemoveUserFromGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminRemoveUserFromGroup, request,
        options: options);
  }

  /// Begins the password reset process. Sets the requested user’s account into a RESET_REQUIRED status, and sends them a password-reset code. Your user pool also sends the user a notification with a r...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminResetUserPasswordResponse>
      adminResetUserPassword(
    $0.AdminResetUserPasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminResetUserPassword, request,
        options: options);
  }

  /// Some API operations in a user pool generate a challenge, like a prompt for an MFA code, for device authentication that bypasses MFA, or for a custom authentication challenge. An AdminRespondToAuthC...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminRespondToAuthChallengeResponse>
      adminRespondToAuthChallenge(
    $0.AdminRespondToAuthChallengeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminRespondToAuthChallenge, request,
        options: options);
  }

  /// Sets the user's multi-factor authentication (MFA) preference, including which MFA options are activated, and if any are preferred. Only one factor can be set as preferred. The preferred MFA factor ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminSetUserMFAPreferenceResponse>
      adminSetUserMFAPreference(
    $0.AdminSetUserMFAPreferenceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminSetUserMFAPreference, request,
        options: options);
  }

  /// Sets the specified user's password in a user pool. This operation administratively sets a temporary or permanent password for a user. With this operation, you can bypass self-service password chang...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminSetUserPasswordResponse> adminSetUserPassword(
    $0.AdminSetUserPasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminSetUserPassword, request, options: options);
  }

  /// This action is no longer supported. You can use it to configure only SMS MFA. You can't use it to configure time-based one-time password (TOTP) software token MFA. Amazon Cognito evaluates Identity...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminSetUserSettingsResponse> adminSetUserSettings(
    $0.AdminSetUserSettingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminSetUserSettings, request, options: options);
  }

  /// Provides the feedback for an authentication event generated by threat protection features. Your response indicates that you think that the event either was from a valid user or was an unwanted auth...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminUpdateAuthEventFeedbackResponse>
      adminUpdateAuthEventFeedback(
    $0.AdminUpdateAuthEventFeedbackRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminUpdateAuthEventFeedback, request,
        options: options);
  }

  /// Updates the status of a user's device so that it is marked as remembered or not remembered for the purpose of device authentication. Device authentication is a "remember me" mechanism that silently...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminUpdateDeviceStatusResponse>
      adminUpdateDeviceStatus(
    $0.AdminUpdateDeviceStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminUpdateDeviceStatus, request,
        options: options);
  }

  /// Updates the specified user's attributes. To delete an attribute from your user, submit the attribute in your API request with a blank value. For custom attributes, you must add a custom: prefix to ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminUpdateUserAttributesResponse>
      adminUpdateUserAttributes(
    $0.AdminUpdateUserAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminUpdateUserAttributes, request,
        options: options);
  }

  /// Invalidates the identity, access, and refresh tokens that Amazon Cognito issued to a user. Call this operation with your administrative credentials when your user signs out of your app. This result...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AdminUserGlobalSignOutResponse>
      adminUserGlobalSignOut(
    $0.AdminUserGlobalSignOutRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$adminUserGlobalSignOut, request,
        options: options);
  }

  /// Begins setup of time-based one-time password (TOTP) multi-factor authentication (MFA) for a user, with a unique private key that Amazon Cognito generates and returns in the API response. You can au...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.AssociateSoftwareTokenResponse>
      associateSoftwareToken(
    $0.AssociateSoftwareTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateSoftwareToken, request,
        options: options);
  }

  /// Changes the password for the currently signed-in user. Authorize this action with a signed-in user's access token. It must include the scope aws.cognito.signin.user.admin. Amazon Cognito doesn't ev...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ChangePasswordResponse> changePassword(
    $0.ChangePasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changePassword, request, options: options);
  }

  /// Completes registration of a passkey authenticator for the currently signed-in user. Authorize this action with a signed-in user's access token. It must include the scope aws.cognito.signin.user.admin.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CompleteWebAuthnRegistrationResponse>
      completeWebAuthnRegistration(
    $0.CompleteWebAuthnRegistrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$completeWebAuthnRegistration, request,
        options: options);
  }

  /// Confirms a device that a user wants to remember. A remembered device is a "Remember me on this device" option for user pools that perform authentication with the device key of a trusted device in t...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ConfirmDeviceResponse> confirmDevice(
    $0.ConfirmDeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$confirmDevice, request, options: options);
  }

  /// This public API operation accepts a confirmation code that Amazon Cognito sent to a user and accepts a new password for that user. Amazon Cognito doesn't evaluate Identity and Access Management (IA...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ConfirmForgotPasswordResponse> confirmForgotPassword(
    $0.ConfirmForgotPasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$confirmForgotPassword, request, options: options);
  }

  /// Confirms the account of a new user. This public API operation submits a code that Amazon Cognito sent to your user when they signed up in your user pool. After your user enters their code, they con...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ConfirmSignUpResponse> confirmSignUp(
    $0.ConfirmSignUpRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$confirmSignUp, request, options: options);
  }

  /// Creates a new group in the specified user pool. For more information about user pool groups, see Adding groups to a user pool. Amazon Cognito evaluates Identity and Access Management (IAM) policies...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateGroupResponse> createGroup(
    $0.CreateGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createGroup, request, options: options);
  }

  /// Adds a configuration and trust relationship between a third-party identity provider (IdP) and a user pool. Amazon Cognito accepts sign-in with third-party identity providers through managed login a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateIdentityProviderResponse>
      createIdentityProvider(
    $0.CreateIdentityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createIdentityProvider, request,
        options: options);
  }

  /// Creates a new set of branding settings for a user pool style and associates it with an app client. This operation is the programmatic option for the creation of a new style in the branding editor. ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateManagedLoginBrandingResponse>
      createManagedLoginBranding(
    $0.CreateManagedLoginBrandingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createManagedLoginBranding, request,
        options: options);
  }

  /// Creates a new OAuth2.0 resource server and defines custom scopes within it. Resource servers are associated with custom scopes and machine-to-machine (M2M) authorization. For more information, see ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateResourceServerResponse> createResourceServer(
    $0.CreateResourceServerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createResourceServer, request, options: options);
  }

  /// Creates terms documents for the requested app client. When Terms and conditions and Privacy policy documents are configured, the app client displays links to them in the sign-up page of managed log...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateTermsResponse> createTerms(
    $0.CreateTermsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createTerms, request, options: options);
  }

  /// Creates a user import job. You can import users into user pools from a comma-separated values (CSV) file without adding Amazon Cognito MAU costs to your Amazon Web Services bill. Amazon Cognito eva...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateUserImportJobResponse> createUserImportJob(
    $0.CreateUserImportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createUserImportJob, request, options: options);
  }

  /// Creates a new Amazon Cognito user pool. This operation sets basic and advanced configuration options. If you don't provide a value for an attribute, Amazon Cognito sets it to its default value. Thi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateUserPoolResponse> createUserPool(
    $0.CreateUserPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createUserPool, request, options: options);
  }

  /// Creates an app client in a user pool. This operation sets basic and advanced configuration options. Unlike app clients created in the console, Amazon Cognito doesn't automatically assign a branding...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateUserPoolClientResponse> createUserPoolClient(
    $0.CreateUserPoolClientRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createUserPoolClient, request, options: options);
  }

  /// A user pool domain hosts managed login, an authorization server and web server for authentication in your application. This operation creates a new user pool prefix domain or custom domain and sets...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.CreateUserPoolDomainResponse> createUserPoolDomain(
    $0.CreateUserPoolDomainRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createUserPoolDomain, request, options: options);
  }

  /// Deletes a group from the specified user pool. When you delete a group, that group no longer contributes to users' cognito:preferred_group or cognito:groups claims, and no longer influence access-co...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteGroup(
    $0.DeleteGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteGroup, request, options: options);
  }

  /// Deletes a user pool identity provider (IdP). After you delete an IdP, users can no longer sign in to your user pool through that IdP. For more information about user pool IdPs, see Third-party IdP ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteIdentityProvider(
    $0.DeleteIdentityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteIdentityProvider, request,
        options: options);
  }

  /// Deletes a managed login branding style. When you delete a style, you delete the branding association for an app client. When an app client doesn't have a style assigned, your managed login pages fo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteManagedLoginBranding(
    $0.DeleteManagedLoginBrandingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteManagedLoginBranding, request,
        options: options);
  }

  /// Deletes a resource server. After you delete a resource server, users can no longer generate access tokens with scopes that are associate with that resource server. Resource servers are associated w...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteResourceServer(
    $0.DeleteResourceServerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteResourceServer, request, options: options);
  }

  /// Deletes the terms documents with the requested ID from your app client. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API operation. For this operation...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteTerms(
    $0.DeleteTermsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteTerms, request, options: options);
  }

  /// Deletes the profile of the currently signed-in user. A deleted user profile can no longer be used to sign in and can't be restored. Authorize this action with a signed-in user's access token. It mu...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteUser(
    $0.DeleteUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUser, request, options: options);
  }

  /// Deletes attributes from the currently signed-in user. For example, your application can submit a request to this operation when a user wants to remove their birthdate attribute value. Authorize thi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteUserAttributesResponse> deleteUserAttributes(
    $0.DeleteUserAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUserAttributes, request, options: options);
  }

  /// Deletes a user pool. After you delete a user pool, users can no longer sign in to any associated applications. When you delete a user pool, it's no longer visible or operational in your Amazon Web ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteUserPool(
    $0.DeleteUserPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUserPool, request, options: options);
  }

  /// Deletes a user pool app client. After you delete an app client, users can no longer sign in to the associated application.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> deleteUserPoolClient(
    $0.DeleteUserPoolClientRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUserPoolClient, request, options: options);
  }

  /// Deletes a specific client secret from a user pool app client. You cannot delete the last remaining secret for an app client.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteUserPoolClientSecretResponse>
      deleteUserPoolClientSecret(
    $0.DeleteUserPoolClientSecretRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUserPoolClientSecret, request,
        options: options);
  }

  /// Given a user pool ID and domain identifier, deletes a user pool domain. After you delete a user pool domain, your managed login pages and authorization server are no longer available.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteUserPoolDomainResponse> deleteUserPoolDomain(
    $0.DeleteUserPoolDomainRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUserPoolDomain, request, options: options);
  }

  /// Deletes a registered passkey, or WebAuthn, authenticator for the currently signed-in user. Authorize this action with a signed-in user's access token. It must include the scope aws.cognito.signin.u...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DeleteWebAuthnCredentialResponse>
      deleteWebAuthnCredential(
    $0.DeleteWebAuthnCredentialRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteWebAuthnCredential, request,
        options: options);
  }

  /// Given a user pool ID and identity provider (IdP) name, returns details about the IdP.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeIdentityProviderResponse>
      describeIdentityProvider(
    $0.DescribeIdentityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeIdentityProvider, request,
        options: options);
  }

  /// Given the ID of a managed login branding style, returns detailed information about the style.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeManagedLoginBrandingResponse>
      describeManagedLoginBranding(
    $0.DescribeManagedLoginBrandingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeManagedLoginBranding, request,
        options: options);
  }

  /// Given the ID of a user pool app client, returns detailed information about the style assigned to the app client.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeManagedLoginBrandingByClientResponse>
      describeManagedLoginBrandingByClient(
    $0.DescribeManagedLoginBrandingByClientRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeManagedLoginBrandingByClient, request,
        options: options);
  }

  /// Describes a resource server. For more information about resource servers, see Access control with resource servers.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeResourceServerResponse>
      describeResourceServer(
    $0.DescribeResourceServerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeResourceServer, request,
        options: options);
  }

  /// Given an app client or user pool ID where threat protection is configured, describes the risk configuration. This operation returns details about adaptive authentication, compromised credentials, a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeRiskConfigurationResponse>
      describeRiskConfiguration(
    $0.DescribeRiskConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeRiskConfiguration, request,
        options: options);
  }

  /// Returns details for the requested terms documents ID. For more information, see Terms documents. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API oper...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeTermsResponse> describeTerms(
    $0.DescribeTermsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeTerms, request, options: options);
  }

  /// Describes a user import job. For more information about user CSV import, see Importing users from a CSV file.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeUserImportJobResponse> describeUserImportJob(
    $0.DescribeUserImportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeUserImportJob, request, options: options);
  }

  /// Given a user pool ID, returns configuration information. This operation is useful when you want to inspect an existing user pool and programmatically replicate the configuration to another user poo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeUserPoolResponse> describeUserPool(
    $0.DescribeUserPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeUserPool, request, options: options);
  }

  /// Given an app client ID, returns configuration information. This operation is useful when you want to inspect an existing app client and programmatically replicate the configuration to another app c...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeUserPoolClientResponse>
      describeUserPoolClient(
    $0.DescribeUserPoolClientRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeUserPoolClient, request,
        options: options);
  }

  /// Given a user pool domain name, returns information about the domain configuration. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API operation. For thi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.DescribeUserPoolDomainResponse>
      describeUserPoolDomain(
    $0.DescribeUserPoolDomainRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$describeUserPoolDomain, request,
        options: options);
  }

  /// Given a device key, deletes a remembered device as the currently signed-in user. For more information about device authentication, see Working with user devices in your user pool. Authorize this ac...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$1.Empty> forgetDevice(
    $0.ForgetDeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$forgetDevice, request, options: options);
  }

  /// Sends a password-reset confirmation code to the email address or phone number of the requested username. The message delivery method is determined by the user's available attributes and the Account...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ForgotPasswordResponse> forgotPassword(
    $0.ForgotPasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$forgotPassword, request, options: options);
  }

  /// Given a user pool ID, generates a comma-separated value (CSV) list populated with available user attributes in the user pool. This list is the header for the CSV file that determines the users in a...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetCSVHeaderResponse> getCSVHeader(
    $0.GetCSVHeaderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCSVHeader, request, options: options);
  }

  /// Given a device key, returns information about a remembered device for the current user. For more information about device authentication, see Working with user devices in your user pool. Authorize ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetDeviceResponse> getDevice(
    $0.GetDeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDevice, request, options: options);
  }

  /// Given a user pool ID and a group name, returns information about the user group. For more information about user pool groups, see Adding groups to a user pool. Amazon Cognito evaluates Identity and...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetGroupResponse> getGroup(
    $0.GetGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGroup, request, options: options);
  }

  /// Given the identifier of an identity provider (IdP), for example examplecorp, returns information about the user pool configuration for that IdP. For more information about IdPs, see Third-party IdP...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetIdentityProviderByIdentifierResponse>
      getIdentityProviderByIdentifier(
    $0.GetIdentityProviderByIdentifierRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getIdentityProviderByIdentifier, request,
        options: options);
  }

  /// Given a user pool ID, returns the logging configuration. User pools can export message-delivery error and threat-protection activity logs to external Amazon Web Services services. For more informat...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetLogDeliveryConfigurationResponse>
      getLogDeliveryConfiguration(
    $0.GetLogDeliveryConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLogDeliveryConfiguration, request,
        options: options);
  }

  /// Given a user pool ID, returns the signing certificate for SAML 2.0 federation. Issued certificates are valid for 10 years from the date of issue. Amazon Cognito issues and assigns a new signing cer...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetSigningCertificateResponse> getSigningCertificate(
    $0.GetSigningCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSigningCertificate, request, options: options);
  }

  /// Given a refresh token, issues new ID, access, and optionally refresh tokens for the user who owns the submitted token. This operation issues a new refresh token and invalidates the original refresh...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetTokensFromRefreshTokenResponse>
      getTokensFromRefreshToken(
    $0.GetTokensFromRefreshTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getTokensFromRefreshToken, request,
        options: options);
  }

  /// Given a user pool ID or app client, returns information about classic hosted UI branding that you applied, if any. Returns user-pool level branding information if no app client branding is applied,...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetUICustomizationResponse> getUICustomization(
    $0.GetUICustomizationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUICustomization, request, options: options);
  }

  /// Gets user attributes and and MFA settings for the currently signed-in user. Authorize this action with a signed-in user's access token. It must include the scope aws.cognito.signin.user.admin. Amaz...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetUserResponse> getUser(
    $0.GetUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUser, request, options: options);
  }

  /// Given an attribute name, sends a user attribute verification code for the specified attribute name to the currently signed-in user. Authorize this action with a signed-in user's access token. It mu...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetUserAttributeVerificationCodeResponse>
      getUserAttributeVerificationCode(
    $0.GetUserAttributeVerificationCodeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUserAttributeVerificationCode, request,
        options: options);
  }

  /// Lists the authentication options for the currently signed-in user. Returns the following: The user's multi-factor authentication (MFA) preferences. The user's options for choice-based authenticatio...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetUserAuthFactorsResponse> getUserAuthFactors(
    $0.GetUserAuthFactorsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUserAuthFactors, request, options: options);
  }

  /// Given a user pool ID, returns configuration for sign-in with WebAuthn authenticators and for multi-factor authentication (MFA). This operation describes the following: The WebAuthn relying party (R...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GetUserPoolMfaConfigResponse> getUserPoolMfaConfig(
    $0.GetUserPoolMfaConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUserPoolMfaConfig, request, options: options);
  }

  /// Invalidates the identity, access, and refresh tokens that Amazon Cognito issued to a user. Call this operation when your user signs out of your app. This results in the following behavior. Amazon C...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.GlobalSignOutResponse> globalSignOut(
    $0.GlobalSignOutRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$globalSignOut, request, options: options);
  }

  /// Declares an authentication flow and initiates sign-in for a user in the Amazon Cognito user directory. Amazon Cognito might respond with an additional challenge or an AuthenticationResult that cont...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.InitiateAuthResponse> initiateAuth(
    $0.InitiateAuthRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$initiateAuth, request, options: options);
  }

  /// Lists the devices that Amazon Cognito has registered to the currently signed-in user. For more information about device authentication, see Working with user devices in your user pool. Authorize th...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListDevicesResponse> listDevices(
    $0.ListDevicesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDevices, request, options: options);
  }

  /// Given a user pool ID, returns user pool groups and their details. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API operation. For this operation, you ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListGroupsResponse> listGroups(
    $0.ListGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGroups, request, options: options);
  }

  /// Given a user pool ID, returns information about configured identity providers (IdPs). For more information about IdPs, see Third-party IdP sign-in. Amazon Cognito evaluates Identity and Access Mana...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListIdentityProvidersResponse> listIdentityProviders(
    $0.ListIdentityProvidersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listIdentityProviders, request, options: options);
  }

  /// Given a user pool ID, returns all resource servers and their details. For more information about resource servers, see Access control with resource servers. Amazon Cognito evaluates Identity and Ac...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListResourceServersResponse> listResourceServers(
    $0.ListResourceServersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listResourceServers, request, options: options);
  }

  /// Lists the tags that are assigned to an Amazon Cognito user pool. For more information, see Tagging resources.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTagsForResourceResponse> listTagsForResource(
    $0.ListTagsForResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTagsForResource, request, options: options);
  }

  /// Returns details about all terms documents for the requested user pool. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API operation. For this operation,...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListTermsResponse> listTerms(
    $0.ListTermsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listTerms, request, options: options);
  }

  /// Given a user pool ID, returns user import jobs and their details. Import jobs are retained in user pool configuration so that you can stage, stop, start, review, and delete them. For more informati...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListUserImportJobsResponse> listUserImportJobs(
    $0.ListUserImportJobsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUserImportJobs, request, options: options);
  }

  /// Given a user pool ID, lists app clients. App clients are sets of rules for the access that you want a user pool to grant to one application. For more information, see App clients. Amazon Cognito ev...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListUserPoolClientsResponse> listUserPoolClients(
    $0.ListUserPoolClientsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUserPoolClients, request, options: options);
  }

  /// Lists all client secrets associated with a user pool app client. Returns metadata about the secrets. The response does not include pagination tokens as there are only 2 secrets at any given time an...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListUserPoolClientSecretsResponse>
      listUserPoolClientSecrets(
    $0.ListUserPoolClientSecretsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUserPoolClientSecrets, request,
        options: options);
  }

  /// Lists user pools and their details in the current Amazon Web Services account. Amazon Cognito evaluates Identity and Access Management (IAM) policies in requests for this API operation. For this op...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListUserPoolsResponse> listUserPools(
    $0.ListUserPoolsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUserPools, request, options: options);
  }

  /// Given a user pool ID, returns a list of users and their basic details in a user pool. This operation is eventually consistent. You might experience a delay before results are up-to-date. To validat...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListUsersResponse> listUsers(
    $0.ListUsersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUsers, request, options: options);
  }

  /// Given a user pool ID and a group name, returns a list of users in the group. For more information about user pool groups, see Adding groups to a user pool. Amazon Cognito evaluates Identity and Acc...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListUsersInGroupResponse> listUsersInGroup(
    $0.ListUsersInGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUsersInGroup, request, options: options);
  }

  /// Generates a list of the currently signed-in user's registered passkey, or WebAuthn, credentials. Authorize this action with a signed-in user's access token. It must include the scope aws.cognito.si...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ListWebAuthnCredentialsResponse>
      listWebAuthnCredentials(
    $0.ListWebAuthnCredentialsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listWebAuthnCredentials, request,
        options: options);
  }

  /// Resends the code that confirms a new account for a user who has signed up in your user pool. Amazon Cognito sends confirmation codes to the user attribute in the AutoVerifiedAttributes property of ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.ResendConfirmationCodeResponse>
      resendConfirmationCode(
    $0.ResendConfirmationCodeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resendConfirmationCode, request,
        options: options);
  }

  /// Some API operations in a user pool generate a challenge, like a prompt for an MFA code, for device authentication that bypasses MFA, or for a custom authentication challenge. A RespondToAuthChallen...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RespondToAuthChallengeResponse>
      respondToAuthChallenge(
    $0.RespondToAuthChallengeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$respondToAuthChallenge, request,
        options: options);
  }

  /// Revokes all of the access tokens generated by, and at the same time as, the specified refresh token. After a token is revoked, you can't use the revoked token to access Amazon Cognito user APIs, or...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.RevokeTokenResponse> revokeToken(
    $0.RevokeTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$revokeToken, request, options: options);
  }

  /// Sets up or modifies the logging configuration of a user pool. User pools can export user notification logs and, when threat protection is active, user-activity logs. For more information, see Expor...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SetLogDeliveryConfigurationResponse>
      setLogDeliveryConfiguration(
    $0.SetLogDeliveryConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setLogDeliveryConfiguration, request,
        options: options);
  }

  /// Configures threat protection for a user pool or app client. Sets configuration for the following. Responses to risks with adaptive authentication Responses to vulnerable passwords with compromised-...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SetRiskConfigurationResponse> setRiskConfiguration(
    $0.SetRiskConfigurationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setRiskConfiguration, request, options: options);
  }

  /// Configures UI branding settings for domains with the hosted UI (classic) branding version. Your user pool must have a domain. Configure a domain with . Set the default configuration for all clients...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SetUICustomizationResponse> setUICustomization(
    $0.SetUICustomizationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setUICustomization, request, options: options);
  }

  /// Set the user's multi-factor authentication (MFA) method preference, including which MFA factors are activated and if any are preferred. Only one factor can be set as preferred. The preferred MFA fa...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SetUserMFAPreferenceResponse> setUserMFAPreference(
    $0.SetUserMFAPreferenceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setUserMFAPreference, request, options: options);
  }

  /// Sets user pool multi-factor authentication (MFA) and passkey configuration. For more information about user pool MFA, see Adding MFA. For more information about WebAuthn passkeys see Authentication...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SetUserPoolMfaConfigResponse> setUserPoolMfaConfig(
    $0.SetUserPoolMfaConfigRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setUserPoolMfaConfig, request, options: options);
  }

  /// This action is no longer supported. You can use it to configure only SMS MFA. You can't use it to configure time-based one-time password (TOTP) software token or email MFA. Authorize this action wi...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SetUserSettingsResponse> setUserSettings(
    $0.SetUserSettingsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setUserSettings, request, options: options);
  }

  /// Registers a user with an app client and requests a user name, password, and user attributes in the user pool. Amazon Cognito doesn't evaluate Identity and Access Management (IAM) policies in reques...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.SignUpResponse> signUp(
    $0.SignUpRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$signUp, request, options: options);
  }

  /// Instructs your user pool to start importing users from a CSV file that contains their usernames and attributes. For more information about importing users from a CSV file, see Importing users from ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartUserImportJobResponse> startUserImportJob(
    $0.StartUserImportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startUserImportJob, request, options: options);
  }

  /// Requests credential creation options from your user pool for the currently signed-in user. Returns information about the user pool, the user profile, and authentication requirements. Users must pro...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StartWebAuthnRegistrationResponse>
      startWebAuthnRegistration(
    $0.StartWebAuthnRegistrationRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$startWebAuthnRegistration, request,
        options: options);
  }

  /// Instructs your user pool to stop a running job that's importing users from a CSV file that contains their usernames and attributes. For more information about importing users from a CSV file, see I...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.StopUserImportJobResponse> stopUserImportJob(
    $0.StopUserImportJobRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$stopUserImportJob, request, options: options);
  }

  /// Assigns a set of tags to an Amazon Cognito user pool. A tag is a label that you can use to categorize and manage user pools in different ways, such as by purpose, owner, environment, or other crite...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.TagResourceResponse> tagResource(
    $0.TagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagResource, request, options: options);
  }

  /// Given tag IDs that you previously assigned to a user pool, removes them.
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UntagResourceResponse> untagResource(
    $0.UntagResourceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagResource, request, options: options);
  }

  /// Provides the feedback for an authentication event generated by threat protection features. The user's response indicates that you think that the event either was from a valid user or was an unwante...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateAuthEventFeedbackResponse>
      updateAuthEventFeedback(
    $0.UpdateAuthEventFeedbackRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAuthEventFeedback, request,
        options: options);
  }

  /// Updates the status of a the currently signed-in user's device so that it is marked as remembered or not remembered for the purpose of device authentication. Device authentication is a "remember me"...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateDeviceStatusResponse> updateDeviceStatus(
    $0.UpdateDeviceStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDeviceStatus, request, options: options);
  }

  /// Given the name of a user pool group, updates any of the properties for precedence, IAM role, or description. For more information about user pool groups, see Adding groups to a user pool. Amazon Co...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateGroupResponse> updateGroup(
    $0.UpdateGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateGroup, request, options: options);
  }

  /// Modifies the configuration and trust relationship between a third-party identity provider (IdP) and a user pool. Amazon Cognito accepts sign-in with third-party identity providers through managed l...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateIdentityProviderResponse>
      updateIdentityProvider(
    $0.UpdateIdentityProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateIdentityProvider, request,
        options: options);
  }

  /// Configures the branding settings for a user pool style. This operation is the programmatic option for the configuration of a style in the branding editor. Provides values for UI customization in a ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateManagedLoginBrandingResponse>
      updateManagedLoginBranding(
    $0.UpdateManagedLoginBrandingRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateManagedLoginBranding, request,
        options: options);
  }

  /// Updates the name and scopes of a resource server. All other fields are read-only. For more information about resource servers, see Access control with resource servers. If you don't provide a value...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateResourceServerResponse> updateResourceServer(
    $0.UpdateResourceServerRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateResourceServer, request, options: options);
  }

  /// Modifies existing terms documents for the requested app client. When Terms and conditions and Privacy policy documents are configured, the app client displays links to them in the sign-up page of m...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateTermsResponse> updateTerms(
    $0.UpdateTermsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateTerms, request, options: options);
  }

  /// Updates the currently signed-in user's attributes. To delete an attribute from the user, submit the attribute in your API request with a blank value. For custom attributes, you must add a custom: p...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateUserAttributesResponse> updateUserAttributes(
    $0.UpdateUserAttributesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateUserAttributes, request, options: options);
  }

  /// Updates the configuration of a user pool. To avoid setting parameters to Amazon Cognito defaults, construct this API request to pass the existing configuration of your user pool, modified to includ...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateUserPoolResponse> updateUserPool(
    $0.UpdateUserPoolRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateUserPool, request, options: options);
  }

  /// Given a user pool app client ID, updates the configuration. To avoid setting parameters to Amazon Cognito defaults, construct this API request to pass the existing configuration of your app client,...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateUserPoolClientResponse> updateUserPoolClient(
    $0.UpdateUserPoolClientRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateUserPoolClient, request, options: options);
  }

  /// A user pool domain hosts managed login, an authorization server and web server for authentication in your application. This operation updates the branding version for user pool domains between 1 fo...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.UpdateUserPoolDomainResponse> updateUserPoolDomain(
    $0.UpdateUserPoolDomainRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateUserPoolDomain, request, options: options);
  }

  /// Registers the current user's time-based one-time password (TOTP) authenticator with a code generated in their authenticator app from a private key that's supplied by your user pool. Marks the user'...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.VerifySoftwareTokenResponse> verifySoftwareToken(
    $0.VerifySoftwareTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$verifySoftwareToken, request, options: options);
  }

  /// Submits a verification code for a signed-in user who has added or changed a value of an auto-verified attribute. When successful, the user's attribute becomes verified and the attribute email_verif...
  /// HTTP:
  /// Protocol: awsJson1_1
  $grpc.ResponseFuture<$0.VerifyUserAttributeResponse> verifyUserAttribute(
    $0.VerifyUserAttributeRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$verifyUserAttribute, request, options: options);
  }

  // method descriptors

  static final _$addCustomAttributes = $grpc.ClientMethod<
          $0.AddCustomAttributesRequest, $0.AddCustomAttributesResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AddCustomAttributes',
      ($0.AddCustomAttributesRequest value) => value.writeToBuffer(),
      $0.AddCustomAttributesResponse.fromBuffer);
  static final _$addUserPoolClientSecret = $grpc.ClientMethod<
          $0.AddUserPoolClientSecretRequest,
          $0.AddUserPoolClientSecretResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AddUserPoolClientSecret',
      ($0.AddUserPoolClientSecretRequest value) => value.writeToBuffer(),
      $0.AddUserPoolClientSecretResponse.fromBuffer);
  static final _$adminAddUserToGroup =
      $grpc.ClientMethod<$0.AdminAddUserToGroupRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/AdminAddUserToGroup',
          ($0.AdminAddUserToGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$adminConfirmSignUp = $grpc.ClientMethod<
          $0.AdminConfirmSignUpRequest, $0.AdminConfirmSignUpResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminConfirmSignUp',
      ($0.AdminConfirmSignUpRequest value) => value.writeToBuffer(),
      $0.AdminConfirmSignUpResponse.fromBuffer);
  static final _$adminCreateUser =
      $grpc.ClientMethod<$0.AdminCreateUserRequest, $0.AdminCreateUserResponse>(
          '/cognito_idp.CognitoIdentityProviderService/AdminCreateUser',
          ($0.AdminCreateUserRequest value) => value.writeToBuffer(),
          $0.AdminCreateUserResponse.fromBuffer);
  static final _$adminDeleteUser =
      $grpc.ClientMethod<$0.AdminDeleteUserRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/AdminDeleteUser',
          ($0.AdminDeleteUserRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$adminDeleteUserAttributes = $grpc.ClientMethod<
          $0.AdminDeleteUserAttributesRequest,
          $0.AdminDeleteUserAttributesResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminDeleteUserAttributes',
      ($0.AdminDeleteUserAttributesRequest value) => value.writeToBuffer(),
      $0.AdminDeleteUserAttributesResponse.fromBuffer);
  static final _$adminDisableProviderForUser = $grpc.ClientMethod<
          $0.AdminDisableProviderForUserRequest,
          $0.AdminDisableProviderForUserResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminDisableProviderForUser',
      ($0.AdminDisableProviderForUserRequest value) => value.writeToBuffer(),
      $0.AdminDisableProviderForUserResponse.fromBuffer);
  static final _$adminDisableUser = $grpc.ClientMethod<
          $0.AdminDisableUserRequest, $0.AdminDisableUserResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminDisableUser',
      ($0.AdminDisableUserRequest value) => value.writeToBuffer(),
      $0.AdminDisableUserResponse.fromBuffer);
  static final _$adminEnableUser =
      $grpc.ClientMethod<$0.AdminEnableUserRequest, $0.AdminEnableUserResponse>(
          '/cognito_idp.CognitoIdentityProviderService/AdminEnableUser',
          ($0.AdminEnableUserRequest value) => value.writeToBuffer(),
          $0.AdminEnableUserResponse.fromBuffer);
  static final _$adminForgetDevice =
      $grpc.ClientMethod<$0.AdminForgetDeviceRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/AdminForgetDevice',
          ($0.AdminForgetDeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$adminGetDevice =
      $grpc.ClientMethod<$0.AdminGetDeviceRequest, $0.AdminGetDeviceResponse>(
          '/cognito_idp.CognitoIdentityProviderService/AdminGetDevice',
          ($0.AdminGetDeviceRequest value) => value.writeToBuffer(),
          $0.AdminGetDeviceResponse.fromBuffer);
  static final _$adminGetUser =
      $grpc.ClientMethod<$0.AdminGetUserRequest, $0.AdminGetUserResponse>(
          '/cognito_idp.CognitoIdentityProviderService/AdminGetUser',
          ($0.AdminGetUserRequest value) => value.writeToBuffer(),
          $0.AdminGetUserResponse.fromBuffer);
  static final _$adminInitiateAuth = $grpc.ClientMethod<
          $0.AdminInitiateAuthRequest, $0.AdminInitiateAuthResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminInitiateAuth',
      ($0.AdminInitiateAuthRequest value) => value.writeToBuffer(),
      $0.AdminInitiateAuthResponse.fromBuffer);
  static final _$adminLinkProviderForUser = $grpc.ClientMethod<
          $0.AdminLinkProviderForUserRequest,
          $0.AdminLinkProviderForUserResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminLinkProviderForUser',
      ($0.AdminLinkProviderForUserRequest value) => value.writeToBuffer(),
      $0.AdminLinkProviderForUserResponse.fromBuffer);
  static final _$adminListDevices = $grpc.ClientMethod<
          $0.AdminListDevicesRequest, $0.AdminListDevicesResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminListDevices',
      ($0.AdminListDevicesRequest value) => value.writeToBuffer(),
      $0.AdminListDevicesResponse.fromBuffer);
  static final _$adminListGroupsForUser = $grpc.ClientMethod<
          $0.AdminListGroupsForUserRequest, $0.AdminListGroupsForUserResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminListGroupsForUser',
      ($0.AdminListGroupsForUserRequest value) => value.writeToBuffer(),
      $0.AdminListGroupsForUserResponse.fromBuffer);
  static final _$adminListUserAuthEvents = $grpc.ClientMethod<
          $0.AdminListUserAuthEventsRequest,
          $0.AdminListUserAuthEventsResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminListUserAuthEvents',
      ($0.AdminListUserAuthEventsRequest value) => value.writeToBuffer(),
      $0.AdminListUserAuthEventsResponse.fromBuffer);
  static final _$adminRemoveUserFromGroup = $grpc.ClientMethod<
          $0.AdminRemoveUserFromGroupRequest, $1.Empty>(
      '/cognito_idp.CognitoIdentityProviderService/AdminRemoveUserFromGroup',
      ($0.AdminRemoveUserFromGroupRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$adminResetUserPassword = $grpc.ClientMethod<
          $0.AdminResetUserPasswordRequest, $0.AdminResetUserPasswordResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminResetUserPassword',
      ($0.AdminResetUserPasswordRequest value) => value.writeToBuffer(),
      $0.AdminResetUserPasswordResponse.fromBuffer);
  static final _$adminRespondToAuthChallenge = $grpc.ClientMethod<
          $0.AdminRespondToAuthChallengeRequest,
          $0.AdminRespondToAuthChallengeResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminRespondToAuthChallenge',
      ($0.AdminRespondToAuthChallengeRequest value) => value.writeToBuffer(),
      $0.AdminRespondToAuthChallengeResponse.fromBuffer);
  static final _$adminSetUserMFAPreference = $grpc.ClientMethod<
          $0.AdminSetUserMFAPreferenceRequest,
          $0.AdminSetUserMFAPreferenceResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminSetUserMFAPreference',
      ($0.AdminSetUserMFAPreferenceRequest value) => value.writeToBuffer(),
      $0.AdminSetUserMFAPreferenceResponse.fromBuffer);
  static final _$adminSetUserPassword = $grpc.ClientMethod<
          $0.AdminSetUserPasswordRequest, $0.AdminSetUserPasswordResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminSetUserPassword',
      ($0.AdminSetUserPasswordRequest value) => value.writeToBuffer(),
      $0.AdminSetUserPasswordResponse.fromBuffer);
  static final _$adminSetUserSettings = $grpc.ClientMethod<
          $0.AdminSetUserSettingsRequest, $0.AdminSetUserSettingsResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminSetUserSettings',
      ($0.AdminSetUserSettingsRequest value) => value.writeToBuffer(),
      $0.AdminSetUserSettingsResponse.fromBuffer);
  static final _$adminUpdateAuthEventFeedback = $grpc.ClientMethod<
          $0.AdminUpdateAuthEventFeedbackRequest,
          $0.AdminUpdateAuthEventFeedbackResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminUpdateAuthEventFeedback',
      ($0.AdminUpdateAuthEventFeedbackRequest value) => value.writeToBuffer(),
      $0.AdminUpdateAuthEventFeedbackResponse.fromBuffer);
  static final _$adminUpdateDeviceStatus = $grpc.ClientMethod<
          $0.AdminUpdateDeviceStatusRequest,
          $0.AdminUpdateDeviceStatusResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminUpdateDeviceStatus',
      ($0.AdminUpdateDeviceStatusRequest value) => value.writeToBuffer(),
      $0.AdminUpdateDeviceStatusResponse.fromBuffer);
  static final _$adminUpdateUserAttributes = $grpc.ClientMethod<
          $0.AdminUpdateUserAttributesRequest,
          $0.AdminUpdateUserAttributesResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminUpdateUserAttributes',
      ($0.AdminUpdateUserAttributesRequest value) => value.writeToBuffer(),
      $0.AdminUpdateUserAttributesResponse.fromBuffer);
  static final _$adminUserGlobalSignOut = $grpc.ClientMethod<
          $0.AdminUserGlobalSignOutRequest, $0.AdminUserGlobalSignOutResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AdminUserGlobalSignOut',
      ($0.AdminUserGlobalSignOutRequest value) => value.writeToBuffer(),
      $0.AdminUserGlobalSignOutResponse.fromBuffer);
  static final _$associateSoftwareToken = $grpc.ClientMethod<
          $0.AssociateSoftwareTokenRequest, $0.AssociateSoftwareTokenResponse>(
      '/cognito_idp.CognitoIdentityProviderService/AssociateSoftwareToken',
      ($0.AssociateSoftwareTokenRequest value) => value.writeToBuffer(),
      $0.AssociateSoftwareTokenResponse.fromBuffer);
  static final _$changePassword =
      $grpc.ClientMethod<$0.ChangePasswordRequest, $0.ChangePasswordResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ChangePassword',
          ($0.ChangePasswordRequest value) => value.writeToBuffer(),
          $0.ChangePasswordResponse.fromBuffer);
  static final _$completeWebAuthnRegistration = $grpc.ClientMethod<
          $0.CompleteWebAuthnRegistrationRequest,
          $0.CompleteWebAuthnRegistrationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/CompleteWebAuthnRegistration',
      ($0.CompleteWebAuthnRegistrationRequest value) => value.writeToBuffer(),
      $0.CompleteWebAuthnRegistrationResponse.fromBuffer);
  static final _$confirmDevice =
      $grpc.ClientMethod<$0.ConfirmDeviceRequest, $0.ConfirmDeviceResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ConfirmDevice',
          ($0.ConfirmDeviceRequest value) => value.writeToBuffer(),
          $0.ConfirmDeviceResponse.fromBuffer);
  static final _$confirmForgotPassword = $grpc.ClientMethod<
          $0.ConfirmForgotPasswordRequest, $0.ConfirmForgotPasswordResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ConfirmForgotPassword',
      ($0.ConfirmForgotPasswordRequest value) => value.writeToBuffer(),
      $0.ConfirmForgotPasswordResponse.fromBuffer);
  static final _$confirmSignUp =
      $grpc.ClientMethod<$0.ConfirmSignUpRequest, $0.ConfirmSignUpResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ConfirmSignUp',
          ($0.ConfirmSignUpRequest value) => value.writeToBuffer(),
          $0.ConfirmSignUpResponse.fromBuffer);
  static final _$createGroup =
      $grpc.ClientMethod<$0.CreateGroupRequest, $0.CreateGroupResponse>(
          '/cognito_idp.CognitoIdentityProviderService/CreateGroup',
          ($0.CreateGroupRequest value) => value.writeToBuffer(),
          $0.CreateGroupResponse.fromBuffer);
  static final _$createIdentityProvider = $grpc.ClientMethod<
          $0.CreateIdentityProviderRequest, $0.CreateIdentityProviderResponse>(
      '/cognito_idp.CognitoIdentityProviderService/CreateIdentityProvider',
      ($0.CreateIdentityProviderRequest value) => value.writeToBuffer(),
      $0.CreateIdentityProviderResponse.fromBuffer);
  static final _$createManagedLoginBranding = $grpc.ClientMethod<
          $0.CreateManagedLoginBrandingRequest,
          $0.CreateManagedLoginBrandingResponse>(
      '/cognito_idp.CognitoIdentityProviderService/CreateManagedLoginBranding',
      ($0.CreateManagedLoginBrandingRequest value) => value.writeToBuffer(),
      $0.CreateManagedLoginBrandingResponse.fromBuffer);
  static final _$createResourceServer = $grpc.ClientMethod<
          $0.CreateResourceServerRequest, $0.CreateResourceServerResponse>(
      '/cognito_idp.CognitoIdentityProviderService/CreateResourceServer',
      ($0.CreateResourceServerRequest value) => value.writeToBuffer(),
      $0.CreateResourceServerResponse.fromBuffer);
  static final _$createTerms =
      $grpc.ClientMethod<$0.CreateTermsRequest, $0.CreateTermsResponse>(
          '/cognito_idp.CognitoIdentityProviderService/CreateTerms',
          ($0.CreateTermsRequest value) => value.writeToBuffer(),
          $0.CreateTermsResponse.fromBuffer);
  static final _$createUserImportJob = $grpc.ClientMethod<
          $0.CreateUserImportJobRequest, $0.CreateUserImportJobResponse>(
      '/cognito_idp.CognitoIdentityProviderService/CreateUserImportJob',
      ($0.CreateUserImportJobRequest value) => value.writeToBuffer(),
      $0.CreateUserImportJobResponse.fromBuffer);
  static final _$createUserPool =
      $grpc.ClientMethod<$0.CreateUserPoolRequest, $0.CreateUserPoolResponse>(
          '/cognito_idp.CognitoIdentityProviderService/CreateUserPool',
          ($0.CreateUserPoolRequest value) => value.writeToBuffer(),
          $0.CreateUserPoolResponse.fromBuffer);
  static final _$createUserPoolClient = $grpc.ClientMethod<
          $0.CreateUserPoolClientRequest, $0.CreateUserPoolClientResponse>(
      '/cognito_idp.CognitoIdentityProviderService/CreateUserPoolClient',
      ($0.CreateUserPoolClientRequest value) => value.writeToBuffer(),
      $0.CreateUserPoolClientResponse.fromBuffer);
  static final _$createUserPoolDomain = $grpc.ClientMethod<
          $0.CreateUserPoolDomainRequest, $0.CreateUserPoolDomainResponse>(
      '/cognito_idp.CognitoIdentityProviderService/CreateUserPoolDomain',
      ($0.CreateUserPoolDomainRequest value) => value.writeToBuffer(),
      $0.CreateUserPoolDomainResponse.fromBuffer);
  static final _$deleteGroup =
      $grpc.ClientMethod<$0.DeleteGroupRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/DeleteGroup',
          ($0.DeleteGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteIdentityProvider =
      $grpc.ClientMethod<$0.DeleteIdentityProviderRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/DeleteIdentityProvider',
          ($0.DeleteIdentityProviderRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteManagedLoginBranding = $grpc.ClientMethod<
          $0.DeleteManagedLoginBrandingRequest, $1.Empty>(
      '/cognito_idp.CognitoIdentityProviderService/DeleteManagedLoginBranding',
      ($0.DeleteManagedLoginBrandingRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$deleteResourceServer =
      $grpc.ClientMethod<$0.DeleteResourceServerRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/DeleteResourceServer',
          ($0.DeleteResourceServerRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteTerms =
      $grpc.ClientMethod<$0.DeleteTermsRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/DeleteTerms',
          ($0.DeleteTermsRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUser =
      $grpc.ClientMethod<$0.DeleteUserRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/DeleteUser',
          ($0.DeleteUserRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUserAttributes = $grpc.ClientMethod<
          $0.DeleteUserAttributesRequest, $0.DeleteUserAttributesResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DeleteUserAttributes',
      ($0.DeleteUserAttributesRequest value) => value.writeToBuffer(),
      $0.DeleteUserAttributesResponse.fromBuffer);
  static final _$deleteUserPool =
      $grpc.ClientMethod<$0.DeleteUserPoolRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/DeleteUserPool',
          ($0.DeleteUserPoolRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUserPoolClient =
      $grpc.ClientMethod<$0.DeleteUserPoolClientRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/DeleteUserPoolClient',
          ($0.DeleteUserPoolClientRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUserPoolClientSecret = $grpc.ClientMethod<
          $0.DeleteUserPoolClientSecretRequest,
          $0.DeleteUserPoolClientSecretResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DeleteUserPoolClientSecret',
      ($0.DeleteUserPoolClientSecretRequest value) => value.writeToBuffer(),
      $0.DeleteUserPoolClientSecretResponse.fromBuffer);
  static final _$deleteUserPoolDomain = $grpc.ClientMethod<
          $0.DeleteUserPoolDomainRequest, $0.DeleteUserPoolDomainResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DeleteUserPoolDomain',
      ($0.DeleteUserPoolDomainRequest value) => value.writeToBuffer(),
      $0.DeleteUserPoolDomainResponse.fromBuffer);
  static final _$deleteWebAuthnCredential = $grpc.ClientMethod<
          $0.DeleteWebAuthnCredentialRequest,
          $0.DeleteWebAuthnCredentialResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DeleteWebAuthnCredential',
      ($0.DeleteWebAuthnCredentialRequest value) => value.writeToBuffer(),
      $0.DeleteWebAuthnCredentialResponse.fromBuffer);
  static final _$describeIdentityProvider = $grpc.ClientMethod<
          $0.DescribeIdentityProviderRequest,
          $0.DescribeIdentityProviderResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeIdentityProvider',
      ($0.DescribeIdentityProviderRequest value) => value.writeToBuffer(),
      $0.DescribeIdentityProviderResponse.fromBuffer);
  static final _$describeManagedLoginBranding = $grpc.ClientMethod<
          $0.DescribeManagedLoginBrandingRequest,
          $0.DescribeManagedLoginBrandingResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeManagedLoginBranding',
      ($0.DescribeManagedLoginBrandingRequest value) => value.writeToBuffer(),
      $0.DescribeManagedLoginBrandingResponse.fromBuffer);
  static final _$describeManagedLoginBrandingByClient = $grpc.ClientMethod<
          $0.DescribeManagedLoginBrandingByClientRequest,
          $0.DescribeManagedLoginBrandingByClientResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeManagedLoginBrandingByClient',
      ($0.DescribeManagedLoginBrandingByClientRequest value) =>
          value.writeToBuffer(),
      $0.DescribeManagedLoginBrandingByClientResponse.fromBuffer);
  static final _$describeResourceServer = $grpc.ClientMethod<
          $0.DescribeResourceServerRequest, $0.DescribeResourceServerResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeResourceServer',
      ($0.DescribeResourceServerRequest value) => value.writeToBuffer(),
      $0.DescribeResourceServerResponse.fromBuffer);
  static final _$describeRiskConfiguration = $grpc.ClientMethod<
          $0.DescribeRiskConfigurationRequest,
          $0.DescribeRiskConfigurationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeRiskConfiguration',
      ($0.DescribeRiskConfigurationRequest value) => value.writeToBuffer(),
      $0.DescribeRiskConfigurationResponse.fromBuffer);
  static final _$describeTerms =
      $grpc.ClientMethod<$0.DescribeTermsRequest, $0.DescribeTermsResponse>(
          '/cognito_idp.CognitoIdentityProviderService/DescribeTerms',
          ($0.DescribeTermsRequest value) => value.writeToBuffer(),
          $0.DescribeTermsResponse.fromBuffer);
  static final _$describeUserImportJob = $grpc.ClientMethod<
          $0.DescribeUserImportJobRequest, $0.DescribeUserImportJobResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeUserImportJob',
      ($0.DescribeUserImportJobRequest value) => value.writeToBuffer(),
      $0.DescribeUserImportJobResponse.fromBuffer);
  static final _$describeUserPool = $grpc.ClientMethod<
          $0.DescribeUserPoolRequest, $0.DescribeUserPoolResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeUserPool',
      ($0.DescribeUserPoolRequest value) => value.writeToBuffer(),
      $0.DescribeUserPoolResponse.fromBuffer);
  static final _$describeUserPoolClient = $grpc.ClientMethod<
          $0.DescribeUserPoolClientRequest, $0.DescribeUserPoolClientResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeUserPoolClient',
      ($0.DescribeUserPoolClientRequest value) => value.writeToBuffer(),
      $0.DescribeUserPoolClientResponse.fromBuffer);
  static final _$describeUserPoolDomain = $grpc.ClientMethod<
          $0.DescribeUserPoolDomainRequest, $0.DescribeUserPoolDomainResponse>(
      '/cognito_idp.CognitoIdentityProviderService/DescribeUserPoolDomain',
      ($0.DescribeUserPoolDomainRequest value) => value.writeToBuffer(),
      $0.DescribeUserPoolDomainResponse.fromBuffer);
  static final _$forgetDevice =
      $grpc.ClientMethod<$0.ForgetDeviceRequest, $1.Empty>(
          '/cognito_idp.CognitoIdentityProviderService/ForgetDevice',
          ($0.ForgetDeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$forgotPassword =
      $grpc.ClientMethod<$0.ForgotPasswordRequest, $0.ForgotPasswordResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ForgotPassword',
          ($0.ForgotPasswordRequest value) => value.writeToBuffer(),
          $0.ForgotPasswordResponse.fromBuffer);
  static final _$getCSVHeader =
      $grpc.ClientMethod<$0.GetCSVHeaderRequest, $0.GetCSVHeaderResponse>(
          '/cognito_idp.CognitoIdentityProviderService/GetCSVHeader',
          ($0.GetCSVHeaderRequest value) => value.writeToBuffer(),
          $0.GetCSVHeaderResponse.fromBuffer);
  static final _$getDevice =
      $grpc.ClientMethod<$0.GetDeviceRequest, $0.GetDeviceResponse>(
          '/cognito_idp.CognitoIdentityProviderService/GetDevice',
          ($0.GetDeviceRequest value) => value.writeToBuffer(),
          $0.GetDeviceResponse.fromBuffer);
  static final _$getGroup =
      $grpc.ClientMethod<$0.GetGroupRequest, $0.GetGroupResponse>(
          '/cognito_idp.CognitoIdentityProviderService/GetGroup',
          ($0.GetGroupRequest value) => value.writeToBuffer(),
          $0.GetGroupResponse.fromBuffer);
  static final _$getIdentityProviderByIdentifier = $grpc.ClientMethod<
          $0.GetIdentityProviderByIdentifierRequest,
          $0.GetIdentityProviderByIdentifierResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetIdentityProviderByIdentifier',
      ($0.GetIdentityProviderByIdentifierRequest value) =>
          value.writeToBuffer(),
      $0.GetIdentityProviderByIdentifierResponse.fromBuffer);
  static final _$getLogDeliveryConfiguration = $grpc.ClientMethod<
          $0.GetLogDeliveryConfigurationRequest,
          $0.GetLogDeliveryConfigurationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetLogDeliveryConfiguration',
      ($0.GetLogDeliveryConfigurationRequest value) => value.writeToBuffer(),
      $0.GetLogDeliveryConfigurationResponse.fromBuffer);
  static final _$getSigningCertificate = $grpc.ClientMethod<
          $0.GetSigningCertificateRequest, $0.GetSigningCertificateResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetSigningCertificate',
      ($0.GetSigningCertificateRequest value) => value.writeToBuffer(),
      $0.GetSigningCertificateResponse.fromBuffer);
  static final _$getTokensFromRefreshToken = $grpc.ClientMethod<
          $0.GetTokensFromRefreshTokenRequest,
          $0.GetTokensFromRefreshTokenResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetTokensFromRefreshToken',
      ($0.GetTokensFromRefreshTokenRequest value) => value.writeToBuffer(),
      $0.GetTokensFromRefreshTokenResponse.fromBuffer);
  static final _$getUICustomization = $grpc.ClientMethod<
          $0.GetUICustomizationRequest, $0.GetUICustomizationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetUICustomization',
      ($0.GetUICustomizationRequest value) => value.writeToBuffer(),
      $0.GetUICustomizationResponse.fromBuffer);
  static final _$getUser =
      $grpc.ClientMethod<$0.GetUserRequest, $0.GetUserResponse>(
          '/cognito_idp.CognitoIdentityProviderService/GetUser',
          ($0.GetUserRequest value) => value.writeToBuffer(),
          $0.GetUserResponse.fromBuffer);
  static final _$getUserAttributeVerificationCode = $grpc.ClientMethod<
          $0.GetUserAttributeVerificationCodeRequest,
          $0.GetUserAttributeVerificationCodeResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetUserAttributeVerificationCode',
      ($0.GetUserAttributeVerificationCodeRequest value) =>
          value.writeToBuffer(),
      $0.GetUserAttributeVerificationCodeResponse.fromBuffer);
  static final _$getUserAuthFactors = $grpc.ClientMethod<
          $0.GetUserAuthFactorsRequest, $0.GetUserAuthFactorsResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetUserAuthFactors',
      ($0.GetUserAuthFactorsRequest value) => value.writeToBuffer(),
      $0.GetUserAuthFactorsResponse.fromBuffer);
  static final _$getUserPoolMfaConfig = $grpc.ClientMethod<
          $0.GetUserPoolMfaConfigRequest, $0.GetUserPoolMfaConfigResponse>(
      '/cognito_idp.CognitoIdentityProviderService/GetUserPoolMfaConfig',
      ($0.GetUserPoolMfaConfigRequest value) => value.writeToBuffer(),
      $0.GetUserPoolMfaConfigResponse.fromBuffer);
  static final _$globalSignOut =
      $grpc.ClientMethod<$0.GlobalSignOutRequest, $0.GlobalSignOutResponse>(
          '/cognito_idp.CognitoIdentityProviderService/GlobalSignOut',
          ($0.GlobalSignOutRequest value) => value.writeToBuffer(),
          $0.GlobalSignOutResponse.fromBuffer);
  static final _$initiateAuth =
      $grpc.ClientMethod<$0.InitiateAuthRequest, $0.InitiateAuthResponse>(
          '/cognito_idp.CognitoIdentityProviderService/InitiateAuth',
          ($0.InitiateAuthRequest value) => value.writeToBuffer(),
          $0.InitiateAuthResponse.fromBuffer);
  static final _$listDevices =
      $grpc.ClientMethod<$0.ListDevicesRequest, $0.ListDevicesResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ListDevices',
          ($0.ListDevicesRequest value) => value.writeToBuffer(),
          $0.ListDevicesResponse.fromBuffer);
  static final _$listGroups =
      $grpc.ClientMethod<$0.ListGroupsRequest, $0.ListGroupsResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ListGroups',
          ($0.ListGroupsRequest value) => value.writeToBuffer(),
          $0.ListGroupsResponse.fromBuffer);
  static final _$listIdentityProviders = $grpc.ClientMethod<
          $0.ListIdentityProvidersRequest, $0.ListIdentityProvidersResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListIdentityProviders',
      ($0.ListIdentityProvidersRequest value) => value.writeToBuffer(),
      $0.ListIdentityProvidersResponse.fromBuffer);
  static final _$listResourceServers = $grpc.ClientMethod<
          $0.ListResourceServersRequest, $0.ListResourceServersResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListResourceServers',
      ($0.ListResourceServersRequest value) => value.writeToBuffer(),
      $0.ListResourceServersResponse.fromBuffer);
  static final _$listTagsForResource = $grpc.ClientMethod<
          $0.ListTagsForResourceRequest, $0.ListTagsForResourceResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListTagsForResource',
      ($0.ListTagsForResourceRequest value) => value.writeToBuffer(),
      $0.ListTagsForResourceResponse.fromBuffer);
  static final _$listTerms =
      $grpc.ClientMethod<$0.ListTermsRequest, $0.ListTermsResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ListTerms',
          ($0.ListTermsRequest value) => value.writeToBuffer(),
          $0.ListTermsResponse.fromBuffer);
  static final _$listUserImportJobs = $grpc.ClientMethod<
          $0.ListUserImportJobsRequest, $0.ListUserImportJobsResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListUserImportJobs',
      ($0.ListUserImportJobsRequest value) => value.writeToBuffer(),
      $0.ListUserImportJobsResponse.fromBuffer);
  static final _$listUserPoolClients = $grpc.ClientMethod<
          $0.ListUserPoolClientsRequest, $0.ListUserPoolClientsResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListUserPoolClients',
      ($0.ListUserPoolClientsRequest value) => value.writeToBuffer(),
      $0.ListUserPoolClientsResponse.fromBuffer);
  static final _$listUserPoolClientSecrets = $grpc.ClientMethod<
          $0.ListUserPoolClientSecretsRequest,
          $0.ListUserPoolClientSecretsResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListUserPoolClientSecrets',
      ($0.ListUserPoolClientSecretsRequest value) => value.writeToBuffer(),
      $0.ListUserPoolClientSecretsResponse.fromBuffer);
  static final _$listUserPools =
      $grpc.ClientMethod<$0.ListUserPoolsRequest, $0.ListUserPoolsResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ListUserPools',
          ($0.ListUserPoolsRequest value) => value.writeToBuffer(),
          $0.ListUserPoolsResponse.fromBuffer);
  static final _$listUsers =
      $grpc.ClientMethod<$0.ListUsersRequest, $0.ListUsersResponse>(
          '/cognito_idp.CognitoIdentityProviderService/ListUsers',
          ($0.ListUsersRequest value) => value.writeToBuffer(),
          $0.ListUsersResponse.fromBuffer);
  static final _$listUsersInGroup = $grpc.ClientMethod<
          $0.ListUsersInGroupRequest, $0.ListUsersInGroupResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListUsersInGroup',
      ($0.ListUsersInGroupRequest value) => value.writeToBuffer(),
      $0.ListUsersInGroupResponse.fromBuffer);
  static final _$listWebAuthnCredentials = $grpc.ClientMethod<
          $0.ListWebAuthnCredentialsRequest,
          $0.ListWebAuthnCredentialsResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ListWebAuthnCredentials',
      ($0.ListWebAuthnCredentialsRequest value) => value.writeToBuffer(),
      $0.ListWebAuthnCredentialsResponse.fromBuffer);
  static final _$resendConfirmationCode = $grpc.ClientMethod<
          $0.ResendConfirmationCodeRequest, $0.ResendConfirmationCodeResponse>(
      '/cognito_idp.CognitoIdentityProviderService/ResendConfirmationCode',
      ($0.ResendConfirmationCodeRequest value) => value.writeToBuffer(),
      $0.ResendConfirmationCodeResponse.fromBuffer);
  static final _$respondToAuthChallenge = $grpc.ClientMethod<
          $0.RespondToAuthChallengeRequest, $0.RespondToAuthChallengeResponse>(
      '/cognito_idp.CognitoIdentityProviderService/RespondToAuthChallenge',
      ($0.RespondToAuthChallengeRequest value) => value.writeToBuffer(),
      $0.RespondToAuthChallengeResponse.fromBuffer);
  static final _$revokeToken =
      $grpc.ClientMethod<$0.RevokeTokenRequest, $0.RevokeTokenResponse>(
          '/cognito_idp.CognitoIdentityProviderService/RevokeToken',
          ($0.RevokeTokenRequest value) => value.writeToBuffer(),
          $0.RevokeTokenResponse.fromBuffer);
  static final _$setLogDeliveryConfiguration = $grpc.ClientMethod<
          $0.SetLogDeliveryConfigurationRequest,
          $0.SetLogDeliveryConfigurationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/SetLogDeliveryConfiguration',
      ($0.SetLogDeliveryConfigurationRequest value) => value.writeToBuffer(),
      $0.SetLogDeliveryConfigurationResponse.fromBuffer);
  static final _$setRiskConfiguration = $grpc.ClientMethod<
          $0.SetRiskConfigurationRequest, $0.SetRiskConfigurationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/SetRiskConfiguration',
      ($0.SetRiskConfigurationRequest value) => value.writeToBuffer(),
      $0.SetRiskConfigurationResponse.fromBuffer);
  static final _$setUICustomization = $grpc.ClientMethod<
          $0.SetUICustomizationRequest, $0.SetUICustomizationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/SetUICustomization',
      ($0.SetUICustomizationRequest value) => value.writeToBuffer(),
      $0.SetUICustomizationResponse.fromBuffer);
  static final _$setUserMFAPreference = $grpc.ClientMethod<
          $0.SetUserMFAPreferenceRequest, $0.SetUserMFAPreferenceResponse>(
      '/cognito_idp.CognitoIdentityProviderService/SetUserMFAPreference',
      ($0.SetUserMFAPreferenceRequest value) => value.writeToBuffer(),
      $0.SetUserMFAPreferenceResponse.fromBuffer);
  static final _$setUserPoolMfaConfig = $grpc.ClientMethod<
          $0.SetUserPoolMfaConfigRequest, $0.SetUserPoolMfaConfigResponse>(
      '/cognito_idp.CognitoIdentityProviderService/SetUserPoolMfaConfig',
      ($0.SetUserPoolMfaConfigRequest value) => value.writeToBuffer(),
      $0.SetUserPoolMfaConfigResponse.fromBuffer);
  static final _$setUserSettings =
      $grpc.ClientMethod<$0.SetUserSettingsRequest, $0.SetUserSettingsResponse>(
          '/cognito_idp.CognitoIdentityProviderService/SetUserSettings',
          ($0.SetUserSettingsRequest value) => value.writeToBuffer(),
          $0.SetUserSettingsResponse.fromBuffer);
  static final _$signUp =
      $grpc.ClientMethod<$0.SignUpRequest, $0.SignUpResponse>(
          '/cognito_idp.CognitoIdentityProviderService/SignUp',
          ($0.SignUpRequest value) => value.writeToBuffer(),
          $0.SignUpResponse.fromBuffer);
  static final _$startUserImportJob = $grpc.ClientMethod<
          $0.StartUserImportJobRequest, $0.StartUserImportJobResponse>(
      '/cognito_idp.CognitoIdentityProviderService/StartUserImportJob',
      ($0.StartUserImportJobRequest value) => value.writeToBuffer(),
      $0.StartUserImportJobResponse.fromBuffer);
  static final _$startWebAuthnRegistration = $grpc.ClientMethod<
          $0.StartWebAuthnRegistrationRequest,
          $0.StartWebAuthnRegistrationResponse>(
      '/cognito_idp.CognitoIdentityProviderService/StartWebAuthnRegistration',
      ($0.StartWebAuthnRegistrationRequest value) => value.writeToBuffer(),
      $0.StartWebAuthnRegistrationResponse.fromBuffer);
  static final _$stopUserImportJob = $grpc.ClientMethod<
          $0.StopUserImportJobRequest, $0.StopUserImportJobResponse>(
      '/cognito_idp.CognitoIdentityProviderService/StopUserImportJob',
      ($0.StopUserImportJobRequest value) => value.writeToBuffer(),
      $0.StopUserImportJobResponse.fromBuffer);
  static final _$tagResource =
      $grpc.ClientMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
          '/cognito_idp.CognitoIdentityProviderService/TagResource',
          ($0.TagResourceRequest value) => value.writeToBuffer(),
          $0.TagResourceResponse.fromBuffer);
  static final _$untagResource =
      $grpc.ClientMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
          '/cognito_idp.CognitoIdentityProviderService/UntagResource',
          ($0.UntagResourceRequest value) => value.writeToBuffer(),
          $0.UntagResourceResponse.fromBuffer);
  static final _$updateAuthEventFeedback = $grpc.ClientMethod<
          $0.UpdateAuthEventFeedbackRequest,
          $0.UpdateAuthEventFeedbackResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateAuthEventFeedback',
      ($0.UpdateAuthEventFeedbackRequest value) => value.writeToBuffer(),
      $0.UpdateAuthEventFeedbackResponse.fromBuffer);
  static final _$updateDeviceStatus = $grpc.ClientMethod<
          $0.UpdateDeviceStatusRequest, $0.UpdateDeviceStatusResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateDeviceStatus',
      ($0.UpdateDeviceStatusRequest value) => value.writeToBuffer(),
      $0.UpdateDeviceStatusResponse.fromBuffer);
  static final _$updateGroup =
      $grpc.ClientMethod<$0.UpdateGroupRequest, $0.UpdateGroupResponse>(
          '/cognito_idp.CognitoIdentityProviderService/UpdateGroup',
          ($0.UpdateGroupRequest value) => value.writeToBuffer(),
          $0.UpdateGroupResponse.fromBuffer);
  static final _$updateIdentityProvider = $grpc.ClientMethod<
          $0.UpdateIdentityProviderRequest, $0.UpdateIdentityProviderResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateIdentityProvider',
      ($0.UpdateIdentityProviderRequest value) => value.writeToBuffer(),
      $0.UpdateIdentityProviderResponse.fromBuffer);
  static final _$updateManagedLoginBranding = $grpc.ClientMethod<
          $0.UpdateManagedLoginBrandingRequest,
          $0.UpdateManagedLoginBrandingResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateManagedLoginBranding',
      ($0.UpdateManagedLoginBrandingRequest value) => value.writeToBuffer(),
      $0.UpdateManagedLoginBrandingResponse.fromBuffer);
  static final _$updateResourceServer = $grpc.ClientMethod<
          $0.UpdateResourceServerRequest, $0.UpdateResourceServerResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateResourceServer',
      ($0.UpdateResourceServerRequest value) => value.writeToBuffer(),
      $0.UpdateResourceServerResponse.fromBuffer);
  static final _$updateTerms =
      $grpc.ClientMethod<$0.UpdateTermsRequest, $0.UpdateTermsResponse>(
          '/cognito_idp.CognitoIdentityProviderService/UpdateTerms',
          ($0.UpdateTermsRequest value) => value.writeToBuffer(),
          $0.UpdateTermsResponse.fromBuffer);
  static final _$updateUserAttributes = $grpc.ClientMethod<
          $0.UpdateUserAttributesRequest, $0.UpdateUserAttributesResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateUserAttributes',
      ($0.UpdateUserAttributesRequest value) => value.writeToBuffer(),
      $0.UpdateUserAttributesResponse.fromBuffer);
  static final _$updateUserPool =
      $grpc.ClientMethod<$0.UpdateUserPoolRequest, $0.UpdateUserPoolResponse>(
          '/cognito_idp.CognitoIdentityProviderService/UpdateUserPool',
          ($0.UpdateUserPoolRequest value) => value.writeToBuffer(),
          $0.UpdateUserPoolResponse.fromBuffer);
  static final _$updateUserPoolClient = $grpc.ClientMethod<
          $0.UpdateUserPoolClientRequest, $0.UpdateUserPoolClientResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateUserPoolClient',
      ($0.UpdateUserPoolClientRequest value) => value.writeToBuffer(),
      $0.UpdateUserPoolClientResponse.fromBuffer);
  static final _$updateUserPoolDomain = $grpc.ClientMethod<
          $0.UpdateUserPoolDomainRequest, $0.UpdateUserPoolDomainResponse>(
      '/cognito_idp.CognitoIdentityProviderService/UpdateUserPoolDomain',
      ($0.UpdateUserPoolDomainRequest value) => value.writeToBuffer(),
      $0.UpdateUserPoolDomainResponse.fromBuffer);
  static final _$verifySoftwareToken = $grpc.ClientMethod<
          $0.VerifySoftwareTokenRequest, $0.VerifySoftwareTokenResponse>(
      '/cognito_idp.CognitoIdentityProviderService/VerifySoftwareToken',
      ($0.VerifySoftwareTokenRequest value) => value.writeToBuffer(),
      $0.VerifySoftwareTokenResponse.fromBuffer);
  static final _$verifyUserAttribute = $grpc.ClientMethod<
          $0.VerifyUserAttributeRequest, $0.VerifyUserAttributeResponse>(
      '/cognito_idp.CognitoIdentityProviderService/VerifyUserAttribute',
      ($0.VerifyUserAttributeRequest value) => value.writeToBuffer(),
      $0.VerifyUserAttributeResponse.fromBuffer);
}

@$pb.GrpcServiceName('cognito_idp.CognitoIdentityProviderService')
abstract class CognitoIdentityProviderServiceBase extends $grpc.Service {
  $core.String get $name => 'cognito_idp.CognitoIdentityProviderService';

  CognitoIdentityProviderServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AddCustomAttributesRequest,
            $0.AddCustomAttributesResponse>(
        'AddCustomAttributes',
        addCustomAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddCustomAttributesRequest.fromBuffer(value),
        ($0.AddCustomAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AddUserPoolClientSecretRequest,
            $0.AddUserPoolClientSecretResponse>(
        'AddUserPoolClientSecret',
        addUserPoolClientSecret_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddUserPoolClientSecretRequest.fromBuffer(value),
        ($0.AddUserPoolClientSecretResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminAddUserToGroupRequest, $1.Empty>(
        'AdminAddUserToGroup',
        adminAddUserToGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminAddUserToGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminConfirmSignUpRequest,
            $0.AdminConfirmSignUpResponse>(
        'AdminConfirmSignUp',
        adminConfirmSignUp_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminConfirmSignUpRequest.fromBuffer(value),
        ($0.AdminConfirmSignUpResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminCreateUserRequest,
            $0.AdminCreateUserResponse>(
        'AdminCreateUser',
        adminCreateUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminCreateUserRequest.fromBuffer(value),
        ($0.AdminCreateUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminDeleteUserRequest, $1.Empty>(
        'AdminDeleteUser',
        adminDeleteUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminDeleteUserRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminDeleteUserAttributesRequest,
            $0.AdminDeleteUserAttributesResponse>(
        'AdminDeleteUserAttributes',
        adminDeleteUserAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminDeleteUserAttributesRequest.fromBuffer(value),
        ($0.AdminDeleteUserAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminDisableProviderForUserRequest,
            $0.AdminDisableProviderForUserResponse>(
        'AdminDisableProviderForUser',
        adminDisableProviderForUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminDisableProviderForUserRequest.fromBuffer(value),
        ($0.AdminDisableProviderForUserResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminDisableUserRequest,
            $0.AdminDisableUserResponse>(
        'AdminDisableUser',
        adminDisableUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminDisableUserRequest.fromBuffer(value),
        ($0.AdminDisableUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminEnableUserRequest,
            $0.AdminEnableUserResponse>(
        'AdminEnableUser',
        adminEnableUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminEnableUserRequest.fromBuffer(value),
        ($0.AdminEnableUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminForgetDeviceRequest, $1.Empty>(
        'AdminForgetDevice',
        adminForgetDevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminForgetDeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminGetDeviceRequest,
            $0.AdminGetDeviceResponse>(
        'AdminGetDevice',
        adminGetDevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminGetDeviceRequest.fromBuffer(value),
        ($0.AdminGetDeviceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.AdminGetUserRequest, $0.AdminGetUserResponse>(
            'AdminGetUser',
            adminGetUser_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.AdminGetUserRequest.fromBuffer(value),
            ($0.AdminGetUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminInitiateAuthRequest,
            $0.AdminInitiateAuthResponse>(
        'AdminInitiateAuth',
        adminInitiateAuth_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminInitiateAuthRequest.fromBuffer(value),
        ($0.AdminInitiateAuthResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminLinkProviderForUserRequest,
            $0.AdminLinkProviderForUserResponse>(
        'AdminLinkProviderForUser',
        adminLinkProviderForUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminLinkProviderForUserRequest.fromBuffer(value),
        ($0.AdminLinkProviderForUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminListDevicesRequest,
            $0.AdminListDevicesResponse>(
        'AdminListDevices',
        adminListDevices_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminListDevicesRequest.fromBuffer(value),
        ($0.AdminListDevicesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminListGroupsForUserRequest,
            $0.AdminListGroupsForUserResponse>(
        'AdminListGroupsForUser',
        adminListGroupsForUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminListGroupsForUserRequest.fromBuffer(value),
        ($0.AdminListGroupsForUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminListUserAuthEventsRequest,
            $0.AdminListUserAuthEventsResponse>(
        'AdminListUserAuthEvents',
        adminListUserAuthEvents_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminListUserAuthEventsRequest.fromBuffer(value),
        ($0.AdminListUserAuthEventsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.AdminRemoveUserFromGroupRequest, $1.Empty>(
            'AdminRemoveUserFromGroup',
            adminRemoveUserFromGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.AdminRemoveUserFromGroupRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminResetUserPasswordRequest,
            $0.AdminResetUserPasswordResponse>(
        'AdminResetUserPassword',
        adminResetUserPassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminResetUserPasswordRequest.fromBuffer(value),
        ($0.AdminResetUserPasswordResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminRespondToAuthChallengeRequest,
            $0.AdminRespondToAuthChallengeResponse>(
        'AdminRespondToAuthChallenge',
        adminRespondToAuthChallenge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminRespondToAuthChallengeRequest.fromBuffer(value),
        ($0.AdminRespondToAuthChallengeResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminSetUserMFAPreferenceRequest,
            $0.AdminSetUserMFAPreferenceResponse>(
        'AdminSetUserMFAPreference',
        adminSetUserMFAPreference_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminSetUserMFAPreferenceRequest.fromBuffer(value),
        ($0.AdminSetUserMFAPreferenceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminSetUserPasswordRequest,
            $0.AdminSetUserPasswordResponse>(
        'AdminSetUserPassword',
        adminSetUserPassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminSetUserPasswordRequest.fromBuffer(value),
        ($0.AdminSetUserPasswordResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminSetUserSettingsRequest,
            $0.AdminSetUserSettingsResponse>(
        'AdminSetUserSettings',
        adminSetUserSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminSetUserSettingsRequest.fromBuffer(value),
        ($0.AdminSetUserSettingsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminUpdateAuthEventFeedbackRequest,
            $0.AdminUpdateAuthEventFeedbackResponse>(
        'AdminUpdateAuthEventFeedback',
        adminUpdateAuthEventFeedback_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminUpdateAuthEventFeedbackRequest.fromBuffer(value),
        ($0.AdminUpdateAuthEventFeedbackResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminUpdateDeviceStatusRequest,
            $0.AdminUpdateDeviceStatusResponse>(
        'AdminUpdateDeviceStatus',
        adminUpdateDeviceStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminUpdateDeviceStatusRequest.fromBuffer(value),
        ($0.AdminUpdateDeviceStatusResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminUpdateUserAttributesRequest,
            $0.AdminUpdateUserAttributesResponse>(
        'AdminUpdateUserAttributes',
        adminUpdateUserAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminUpdateUserAttributesRequest.fromBuffer(value),
        ($0.AdminUpdateUserAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AdminUserGlobalSignOutRequest,
            $0.AdminUserGlobalSignOutResponse>(
        'AdminUserGlobalSignOut',
        adminUserGlobalSignOut_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AdminUserGlobalSignOutRequest.fromBuffer(value),
        ($0.AdminUserGlobalSignOutResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AssociateSoftwareTokenRequest,
            $0.AssociateSoftwareTokenResponse>(
        'AssociateSoftwareToken',
        associateSoftwareToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AssociateSoftwareTokenRequest.fromBuffer(value),
        ($0.AssociateSoftwareTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangePasswordRequest,
            $0.ChangePasswordResponse>(
        'ChangePassword',
        changePassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangePasswordRequest.fromBuffer(value),
        ($0.ChangePasswordResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CompleteWebAuthnRegistrationRequest,
            $0.CompleteWebAuthnRegistrationResponse>(
        'CompleteWebAuthnRegistration',
        completeWebAuthnRegistration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CompleteWebAuthnRegistrationRequest.fromBuffer(value),
        ($0.CompleteWebAuthnRegistrationResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ConfirmDeviceRequest, $0.ConfirmDeviceResponse>(
            'ConfirmDevice',
            confirmDevice_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ConfirmDeviceRequest.fromBuffer(value),
            ($0.ConfirmDeviceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ConfirmForgotPasswordRequest,
            $0.ConfirmForgotPasswordResponse>(
        'ConfirmForgotPassword',
        confirmForgotPassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ConfirmForgotPasswordRequest.fromBuffer(value),
        ($0.ConfirmForgotPasswordResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ConfirmSignUpRequest, $0.ConfirmSignUpResponse>(
            'ConfirmSignUp',
            confirmSignUp_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ConfirmSignUpRequest.fromBuffer(value),
            ($0.ConfirmSignUpResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateGroupRequest, $0.CreateGroupResponse>(
            'CreateGroup',
            createGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateGroupRequest.fromBuffer(value),
            ($0.CreateGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateIdentityProviderRequest,
            $0.CreateIdentityProviderResponse>(
        'CreateIdentityProvider',
        createIdentityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateIdentityProviderRequest.fromBuffer(value),
        ($0.CreateIdentityProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateManagedLoginBrandingRequest,
            $0.CreateManagedLoginBrandingResponse>(
        'CreateManagedLoginBranding',
        createManagedLoginBranding_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateManagedLoginBrandingRequest.fromBuffer(value),
        ($0.CreateManagedLoginBrandingResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateResourceServerRequest,
            $0.CreateResourceServerResponse>(
        'CreateResourceServer',
        createResourceServer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateResourceServerRequest.fromBuffer(value),
        ($0.CreateResourceServerResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateTermsRequest, $0.CreateTermsResponse>(
            'CreateTerms',
            createTerms_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateTermsRequest.fromBuffer(value),
            ($0.CreateTermsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateUserImportJobRequest,
            $0.CreateUserImportJobResponse>(
        'CreateUserImportJob',
        createUserImportJob_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateUserImportJobRequest.fromBuffer(value),
        ($0.CreateUserImportJobResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateUserPoolRequest,
            $0.CreateUserPoolResponse>(
        'CreateUserPool',
        createUserPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateUserPoolRequest.fromBuffer(value),
        ($0.CreateUserPoolResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateUserPoolClientRequest,
            $0.CreateUserPoolClientResponse>(
        'CreateUserPoolClient',
        createUserPoolClient_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateUserPoolClientRequest.fromBuffer(value),
        ($0.CreateUserPoolClientResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateUserPoolDomainRequest,
            $0.CreateUserPoolDomainResponse>(
        'CreateUserPoolDomain',
        createUserPoolDomain_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateUserPoolDomainRequest.fromBuffer(value),
        ($0.CreateUserPoolDomainResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteGroupRequest, $1.Empty>(
        'DeleteGroup',
        deleteGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteIdentityProviderRequest, $1.Empty>(
        'DeleteIdentityProvider',
        deleteIdentityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteIdentityProviderRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteManagedLoginBrandingRequest, $1.Empty>(
            'DeleteManagedLoginBranding',
            deleteManagedLoginBranding_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteManagedLoginBrandingRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteResourceServerRequest, $1.Empty>(
        'DeleteResourceServer',
        deleteResourceServer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteResourceServerRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteTermsRequest, $1.Empty>(
        'DeleteTerms',
        deleteTerms_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteTermsRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserRequest, $1.Empty>(
        'DeleteUser',
        deleteUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteUserRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserAttributesRequest,
            $0.DeleteUserAttributesResponse>(
        'DeleteUserAttributes',
        deleteUserAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUserAttributesRequest.fromBuffer(value),
        ($0.DeleteUserAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserPoolRequest, $1.Empty>(
        'DeleteUserPool',
        deleteUserPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUserPoolRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserPoolClientRequest, $1.Empty>(
        'DeleteUserPoolClient',
        deleteUserPoolClient_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUserPoolClientRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserPoolClientSecretRequest,
            $0.DeleteUserPoolClientSecretResponse>(
        'DeleteUserPoolClientSecret',
        deleteUserPoolClientSecret_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUserPoolClientSecretRequest.fromBuffer(value),
        ($0.DeleteUserPoolClientSecretResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserPoolDomainRequest,
            $0.DeleteUserPoolDomainResponse>(
        'DeleteUserPoolDomain',
        deleteUserPoolDomain_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUserPoolDomainRequest.fromBuffer(value),
        ($0.DeleteUserPoolDomainResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteWebAuthnCredentialRequest,
            $0.DeleteWebAuthnCredentialResponse>(
        'DeleteWebAuthnCredential',
        deleteWebAuthnCredential_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteWebAuthnCredentialRequest.fromBuffer(value),
        ($0.DeleteWebAuthnCredentialResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeIdentityProviderRequest,
            $0.DescribeIdentityProviderResponse>(
        'DescribeIdentityProvider',
        describeIdentityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeIdentityProviderRequest.fromBuffer(value),
        ($0.DescribeIdentityProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeManagedLoginBrandingRequest,
            $0.DescribeManagedLoginBrandingResponse>(
        'DescribeManagedLoginBranding',
        describeManagedLoginBranding_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeManagedLoginBrandingRequest.fromBuffer(value),
        ($0.DescribeManagedLoginBrandingResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DescribeManagedLoginBrandingByClientRequest,
            $0.DescribeManagedLoginBrandingByClientResponse>(
        'DescribeManagedLoginBrandingByClient',
        describeManagedLoginBrandingByClient_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeManagedLoginBrandingByClientRequest.fromBuffer(value),
        ($0.DescribeManagedLoginBrandingByClientResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeResourceServerRequest,
            $0.DescribeResourceServerResponse>(
        'DescribeResourceServer',
        describeResourceServer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeResourceServerRequest.fromBuffer(value),
        ($0.DescribeResourceServerResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeRiskConfigurationRequest,
            $0.DescribeRiskConfigurationResponse>(
        'DescribeRiskConfiguration',
        describeRiskConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeRiskConfigurationRequest.fromBuffer(value),
        ($0.DescribeRiskConfigurationResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DescribeTermsRequest, $0.DescribeTermsResponse>(
            'DescribeTerms',
            describeTerms_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DescribeTermsRequest.fromBuffer(value),
            ($0.DescribeTermsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeUserImportJobRequest,
            $0.DescribeUserImportJobResponse>(
        'DescribeUserImportJob',
        describeUserImportJob_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeUserImportJobRequest.fromBuffer(value),
        ($0.DescribeUserImportJobResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeUserPoolRequest,
            $0.DescribeUserPoolResponse>(
        'DescribeUserPool',
        describeUserPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeUserPoolRequest.fromBuffer(value),
        ($0.DescribeUserPoolResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeUserPoolClientRequest,
            $0.DescribeUserPoolClientResponse>(
        'DescribeUserPoolClient',
        describeUserPoolClient_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeUserPoolClientRequest.fromBuffer(value),
        ($0.DescribeUserPoolClientResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DescribeUserPoolDomainRequest,
            $0.DescribeUserPoolDomainResponse>(
        'DescribeUserPoolDomain',
        describeUserPoolDomain_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DescribeUserPoolDomainRequest.fromBuffer(value),
        ($0.DescribeUserPoolDomainResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ForgetDeviceRequest, $1.Empty>(
        'ForgetDevice',
        forgetDevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ForgetDeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ForgotPasswordRequest,
            $0.ForgotPasswordResponse>(
        'ForgotPassword',
        forgotPassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ForgotPasswordRequest.fromBuffer(value),
        ($0.ForgotPasswordResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetCSVHeaderRequest, $0.GetCSVHeaderResponse>(
            'GetCSVHeader',
            getCSVHeader_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetCSVHeaderRequest.fromBuffer(value),
            ($0.GetCSVHeaderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDeviceRequest, $0.GetDeviceResponse>(
        'GetDevice',
        getDevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetDeviceRequest.fromBuffer(value),
        ($0.GetDeviceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetGroupRequest, $0.GetGroupResponse>(
        'GetGroup',
        getGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetGroupRequest.fromBuffer(value),
        ($0.GetGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetIdentityProviderByIdentifierRequest,
            $0.GetIdentityProviderByIdentifierResponse>(
        'GetIdentityProviderByIdentifier',
        getIdentityProviderByIdentifier_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetIdentityProviderByIdentifierRequest.fromBuffer(value),
        ($0.GetIdentityProviderByIdentifierResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetLogDeliveryConfigurationRequest,
            $0.GetLogDeliveryConfigurationResponse>(
        'GetLogDeliveryConfiguration',
        getLogDeliveryConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetLogDeliveryConfigurationRequest.fromBuffer(value),
        ($0.GetLogDeliveryConfigurationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSigningCertificateRequest,
            $0.GetSigningCertificateResponse>(
        'GetSigningCertificate',
        getSigningCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSigningCertificateRequest.fromBuffer(value),
        ($0.GetSigningCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetTokensFromRefreshTokenRequest,
            $0.GetTokensFromRefreshTokenResponse>(
        'GetTokensFromRefreshToken',
        getTokensFromRefreshToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetTokensFromRefreshTokenRequest.fromBuffer(value),
        ($0.GetTokensFromRefreshTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUICustomizationRequest,
            $0.GetUICustomizationResponse>(
        'GetUICustomization',
        getUICustomization_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUICustomizationRequest.fromBuffer(value),
        ($0.GetUICustomizationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUserRequest, $0.GetUserResponse>(
        'GetUser',
        getUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetUserRequest.fromBuffer(value),
        ($0.GetUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUserAttributeVerificationCodeRequest,
            $0.GetUserAttributeVerificationCodeResponse>(
        'GetUserAttributeVerificationCode',
        getUserAttributeVerificationCode_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUserAttributeVerificationCodeRequest.fromBuffer(value),
        ($0.GetUserAttributeVerificationCodeResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUserAuthFactorsRequest,
            $0.GetUserAuthFactorsResponse>(
        'GetUserAuthFactors',
        getUserAuthFactors_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUserAuthFactorsRequest.fromBuffer(value),
        ($0.GetUserAuthFactorsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUserPoolMfaConfigRequest,
            $0.GetUserPoolMfaConfigResponse>(
        'GetUserPoolMfaConfig',
        getUserPoolMfaConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetUserPoolMfaConfigRequest.fromBuffer(value),
        ($0.GetUserPoolMfaConfigResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GlobalSignOutRequest, $0.GlobalSignOutResponse>(
            'GlobalSignOut',
            globalSignOut_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GlobalSignOutRequest.fromBuffer(value),
            ($0.GlobalSignOutResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.InitiateAuthRequest, $0.InitiateAuthResponse>(
            'InitiateAuth',
            initiateAuth_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.InitiateAuthRequest.fromBuffer(value),
            ($0.InitiateAuthResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListDevicesRequest, $0.ListDevicesResponse>(
            'ListDevices',
            listDevices_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListDevicesRequest.fromBuffer(value),
            ($0.ListDevicesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGroupsRequest, $0.ListGroupsResponse>(
        'ListGroups',
        listGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListGroupsRequest.fromBuffer(value),
        ($0.ListGroupsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListIdentityProvidersRequest,
            $0.ListIdentityProvidersResponse>(
        'ListIdentityProviders',
        listIdentityProviders_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListIdentityProvidersRequest.fromBuffer(value),
        ($0.ListIdentityProvidersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListResourceServersRequest,
            $0.ListResourceServersResponse>(
        'ListResourceServers',
        listResourceServers_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListResourceServersRequest.fromBuffer(value),
        ($0.ListResourceServersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTagsForResourceRequest,
            $0.ListTagsForResourceResponse>(
        'ListTagsForResource',
        listTagsForResource_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListTagsForResourceRequest.fromBuffer(value),
        ($0.ListTagsForResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListTermsRequest, $0.ListTermsResponse>(
        'ListTerms',
        listTerms_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListTermsRequest.fromBuffer(value),
        ($0.ListTermsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListUserImportJobsRequest,
            $0.ListUserImportJobsResponse>(
        'ListUserImportJobs',
        listUserImportJobs_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListUserImportJobsRequest.fromBuffer(value),
        ($0.ListUserImportJobsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListUserPoolClientsRequest,
            $0.ListUserPoolClientsResponse>(
        'ListUserPoolClients',
        listUserPoolClients_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListUserPoolClientsRequest.fromBuffer(value),
        ($0.ListUserPoolClientsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListUserPoolClientSecretsRequest,
            $0.ListUserPoolClientSecretsResponse>(
        'ListUserPoolClientSecrets',
        listUserPoolClientSecrets_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListUserPoolClientSecretsRequest.fromBuffer(value),
        ($0.ListUserPoolClientSecretsResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListUserPoolsRequest, $0.ListUserPoolsResponse>(
            'ListUserPools',
            listUserPools_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListUserPoolsRequest.fromBuffer(value),
            ($0.ListUserPoolsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListUsersRequest, $0.ListUsersResponse>(
        'ListUsers',
        listUsers_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListUsersRequest.fromBuffer(value),
        ($0.ListUsersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListUsersInGroupRequest,
            $0.ListUsersInGroupResponse>(
        'ListUsersInGroup',
        listUsersInGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListUsersInGroupRequest.fromBuffer(value),
        ($0.ListUsersInGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListWebAuthnCredentialsRequest,
            $0.ListWebAuthnCredentialsResponse>(
        'ListWebAuthnCredentials',
        listWebAuthnCredentials_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListWebAuthnCredentialsRequest.fromBuffer(value),
        ($0.ListWebAuthnCredentialsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResendConfirmationCodeRequest,
            $0.ResendConfirmationCodeResponse>(
        'ResendConfirmationCode',
        resendConfirmationCode_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResendConfirmationCodeRequest.fromBuffer(value),
        ($0.ResendConfirmationCodeResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RespondToAuthChallengeRequest,
            $0.RespondToAuthChallengeResponse>(
        'RespondToAuthChallenge',
        respondToAuthChallenge_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RespondToAuthChallengeRequest.fromBuffer(value),
        ($0.RespondToAuthChallengeResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RevokeTokenRequest, $0.RevokeTokenResponse>(
            'RevokeToken',
            revokeToken_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RevokeTokenRequest.fromBuffer(value),
            ($0.RevokeTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetLogDeliveryConfigurationRequest,
            $0.SetLogDeliveryConfigurationResponse>(
        'SetLogDeliveryConfiguration',
        setLogDeliveryConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetLogDeliveryConfigurationRequest.fromBuffer(value),
        ($0.SetLogDeliveryConfigurationResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetRiskConfigurationRequest,
            $0.SetRiskConfigurationResponse>(
        'SetRiskConfiguration',
        setRiskConfiguration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetRiskConfigurationRequest.fromBuffer(value),
        ($0.SetRiskConfigurationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetUICustomizationRequest,
            $0.SetUICustomizationResponse>(
        'SetUICustomization',
        setUICustomization_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetUICustomizationRequest.fromBuffer(value),
        ($0.SetUICustomizationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetUserMFAPreferenceRequest,
            $0.SetUserMFAPreferenceResponse>(
        'SetUserMFAPreference',
        setUserMFAPreference_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetUserMFAPreferenceRequest.fromBuffer(value),
        ($0.SetUserMFAPreferenceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetUserPoolMfaConfigRequest,
            $0.SetUserPoolMfaConfigResponse>(
        'SetUserPoolMfaConfig',
        setUserPoolMfaConfig_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetUserPoolMfaConfigRequest.fromBuffer(value),
        ($0.SetUserPoolMfaConfigResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetUserSettingsRequest,
            $0.SetUserSettingsResponse>(
        'SetUserSettings',
        setUserSettings_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetUserSettingsRequest.fromBuffer(value),
        ($0.SetUserSettingsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SignUpRequest, $0.SignUpResponse>(
        'SignUp',
        signUp_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.SignUpRequest.fromBuffer(value),
        ($0.SignUpResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartUserImportJobRequest,
            $0.StartUserImportJobResponse>(
        'StartUserImportJob',
        startUserImportJob_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartUserImportJobRequest.fromBuffer(value),
        ($0.StartUserImportJobResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StartWebAuthnRegistrationRequest,
            $0.StartWebAuthnRegistrationResponse>(
        'StartWebAuthnRegistration',
        startWebAuthnRegistration_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StartWebAuthnRegistrationRequest.fromBuffer(value),
        ($0.StartWebAuthnRegistrationResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.StopUserImportJobRequest,
            $0.StopUserImportJobResponse>(
        'StopUserImportJob',
        stopUserImportJob_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.StopUserImportJobRequest.fromBuffer(value),
        ($0.StopUserImportJobResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TagResourceRequest, $0.TagResourceResponse>(
            'TagResource',
            tagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TagResourceRequest.fromBuffer(value),
            ($0.TagResourceResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagResourceRequest, $0.UntagResourceResponse>(
            'UntagResource',
            untagResource_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagResourceRequest.fromBuffer(value),
            ($0.UntagResourceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAuthEventFeedbackRequest,
            $0.UpdateAuthEventFeedbackResponse>(
        'UpdateAuthEventFeedback',
        updateAuthEventFeedback_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAuthEventFeedbackRequest.fromBuffer(value),
        ($0.UpdateAuthEventFeedbackResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDeviceStatusRequest,
            $0.UpdateDeviceStatusResponse>(
        'UpdateDeviceStatus',
        updateDeviceStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDeviceStatusRequest.fromBuffer(value),
        ($0.UpdateDeviceStatusResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateGroupRequest, $0.UpdateGroupResponse>(
            'UpdateGroup',
            updateGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateGroupRequest.fromBuffer(value),
            ($0.UpdateGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateIdentityProviderRequest,
            $0.UpdateIdentityProviderResponse>(
        'UpdateIdentityProvider',
        updateIdentityProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateIdentityProviderRequest.fromBuffer(value),
        ($0.UpdateIdentityProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateManagedLoginBrandingRequest,
            $0.UpdateManagedLoginBrandingResponse>(
        'UpdateManagedLoginBranding',
        updateManagedLoginBranding_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateManagedLoginBrandingRequest.fromBuffer(value),
        ($0.UpdateManagedLoginBrandingResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateResourceServerRequest,
            $0.UpdateResourceServerResponse>(
        'UpdateResourceServer',
        updateResourceServer_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateResourceServerRequest.fromBuffer(value),
        ($0.UpdateResourceServerResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateTermsRequest, $0.UpdateTermsResponse>(
            'UpdateTerms',
            updateTerms_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateTermsRequest.fromBuffer(value),
            ($0.UpdateTermsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateUserAttributesRequest,
            $0.UpdateUserAttributesResponse>(
        'UpdateUserAttributes',
        updateUserAttributes_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateUserAttributesRequest.fromBuffer(value),
        ($0.UpdateUserAttributesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateUserPoolRequest,
            $0.UpdateUserPoolResponse>(
        'UpdateUserPool',
        updateUserPool_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateUserPoolRequest.fromBuffer(value),
        ($0.UpdateUserPoolResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateUserPoolClientRequest,
            $0.UpdateUserPoolClientResponse>(
        'UpdateUserPoolClient',
        updateUserPoolClient_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateUserPoolClientRequest.fromBuffer(value),
        ($0.UpdateUserPoolClientResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateUserPoolDomainRequest,
            $0.UpdateUserPoolDomainResponse>(
        'UpdateUserPoolDomain',
        updateUserPoolDomain_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateUserPoolDomainRequest.fromBuffer(value),
        ($0.UpdateUserPoolDomainResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.VerifySoftwareTokenRequest,
            $0.VerifySoftwareTokenResponse>(
        'VerifySoftwareToken',
        verifySoftwareToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.VerifySoftwareTokenRequest.fromBuffer(value),
        ($0.VerifySoftwareTokenResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.VerifyUserAttributeRequest,
            $0.VerifyUserAttributeResponse>(
        'VerifyUserAttribute',
        verifyUserAttribute_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.VerifyUserAttributeRequest.fromBuffer(value),
        ($0.VerifyUserAttributeResponse value) => value.writeToBuffer()));
  }

  $async.Future<$0.AddCustomAttributesResponse> addCustomAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AddCustomAttributesRequest> $request) async {
    return addCustomAttributes($call, await $request);
  }

  $async.Future<$0.AddCustomAttributesResponse> addCustomAttributes(
      $grpc.ServiceCall call, $0.AddCustomAttributesRequest request);

  $async.Future<$0.AddUserPoolClientSecretResponse> addUserPoolClientSecret_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AddUserPoolClientSecretRequest> $request) async {
    return addUserPoolClientSecret($call, await $request);
  }

  $async.Future<$0.AddUserPoolClientSecretResponse> addUserPoolClientSecret(
      $grpc.ServiceCall call, $0.AddUserPoolClientSecretRequest request);

  $async.Future<$1.Empty> adminAddUserToGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AdminAddUserToGroupRequest> $request) async {
    return adminAddUserToGroup($call, await $request);
  }

  $async.Future<$1.Empty> adminAddUserToGroup(
      $grpc.ServiceCall call, $0.AdminAddUserToGroupRequest request);

  $async.Future<$0.AdminConfirmSignUpResponse> adminConfirmSignUp_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminConfirmSignUpRequest> $request) async {
    return adminConfirmSignUp($call, await $request);
  }

  $async.Future<$0.AdminConfirmSignUpResponse> adminConfirmSignUp(
      $grpc.ServiceCall call, $0.AdminConfirmSignUpRequest request);

  $async.Future<$0.AdminCreateUserResponse> adminCreateUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminCreateUserRequest> $request) async {
    return adminCreateUser($call, await $request);
  }

  $async.Future<$0.AdminCreateUserResponse> adminCreateUser(
      $grpc.ServiceCall call, $0.AdminCreateUserRequest request);

  $async.Future<$1.Empty> adminDeleteUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AdminDeleteUserRequest> $request) async {
    return adminDeleteUser($call, await $request);
  }

  $async.Future<$1.Empty> adminDeleteUser(
      $grpc.ServiceCall call, $0.AdminDeleteUserRequest request);

  $async.Future<$0.AdminDeleteUserAttributesResponse>
      adminDeleteUserAttributes_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AdminDeleteUserAttributesRequest> $request) async {
    return adminDeleteUserAttributes($call, await $request);
  }

  $async.Future<$0.AdminDeleteUserAttributesResponse> adminDeleteUserAttributes(
      $grpc.ServiceCall call, $0.AdminDeleteUserAttributesRequest request);

  $async.Future<$0.AdminDisableProviderForUserResponse>
      adminDisableProviderForUser_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AdminDisableProviderForUserRequest> $request) async {
    return adminDisableProviderForUser($call, await $request);
  }

  $async.Future<$0.AdminDisableProviderForUserResponse>
      adminDisableProviderForUser($grpc.ServiceCall call,
          $0.AdminDisableProviderForUserRequest request);

  $async.Future<$0.AdminDisableUserResponse> adminDisableUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminDisableUserRequest> $request) async {
    return adminDisableUser($call, await $request);
  }

  $async.Future<$0.AdminDisableUserResponse> adminDisableUser(
      $grpc.ServiceCall call, $0.AdminDisableUserRequest request);

  $async.Future<$0.AdminEnableUserResponse> adminEnableUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminEnableUserRequest> $request) async {
    return adminEnableUser($call, await $request);
  }

  $async.Future<$0.AdminEnableUserResponse> adminEnableUser(
      $grpc.ServiceCall call, $0.AdminEnableUserRequest request);

  $async.Future<$1.Empty> adminForgetDevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AdminForgetDeviceRequest> $request) async {
    return adminForgetDevice($call, await $request);
  }

  $async.Future<$1.Empty> adminForgetDevice(
      $grpc.ServiceCall call, $0.AdminForgetDeviceRequest request);

  $async.Future<$0.AdminGetDeviceResponse> adminGetDevice_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminGetDeviceRequest> $request) async {
    return adminGetDevice($call, await $request);
  }

  $async.Future<$0.AdminGetDeviceResponse> adminGetDevice(
      $grpc.ServiceCall call, $0.AdminGetDeviceRequest request);

  $async.Future<$0.AdminGetUserResponse> adminGetUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminGetUserRequest> $request) async {
    return adminGetUser($call, await $request);
  }

  $async.Future<$0.AdminGetUserResponse> adminGetUser(
      $grpc.ServiceCall call, $0.AdminGetUserRequest request);

  $async.Future<$0.AdminInitiateAuthResponse> adminInitiateAuth_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminInitiateAuthRequest> $request) async {
    return adminInitiateAuth($call, await $request);
  }

  $async.Future<$0.AdminInitiateAuthResponse> adminInitiateAuth(
      $grpc.ServiceCall call, $0.AdminInitiateAuthRequest request);

  $async.Future<$0.AdminLinkProviderForUserResponse>
      adminLinkProviderForUser_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AdminLinkProviderForUserRequest> $request) async {
    return adminLinkProviderForUser($call, await $request);
  }

  $async.Future<$0.AdminLinkProviderForUserResponse> adminLinkProviderForUser(
      $grpc.ServiceCall call, $0.AdminLinkProviderForUserRequest request);

  $async.Future<$0.AdminListDevicesResponse> adminListDevices_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminListDevicesRequest> $request) async {
    return adminListDevices($call, await $request);
  }

  $async.Future<$0.AdminListDevicesResponse> adminListDevices(
      $grpc.ServiceCall call, $0.AdminListDevicesRequest request);

  $async.Future<$0.AdminListGroupsForUserResponse> adminListGroupsForUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminListGroupsForUserRequest> $request) async {
    return adminListGroupsForUser($call, await $request);
  }

  $async.Future<$0.AdminListGroupsForUserResponse> adminListGroupsForUser(
      $grpc.ServiceCall call, $0.AdminListGroupsForUserRequest request);

  $async.Future<$0.AdminListUserAuthEventsResponse> adminListUserAuthEvents_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminListUserAuthEventsRequest> $request) async {
    return adminListUserAuthEvents($call, await $request);
  }

  $async.Future<$0.AdminListUserAuthEventsResponse> adminListUserAuthEvents(
      $grpc.ServiceCall call, $0.AdminListUserAuthEventsRequest request);

  $async.Future<$1.Empty> adminRemoveUserFromGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AdminRemoveUserFromGroupRequest> $request) async {
    return adminRemoveUserFromGroup($call, await $request);
  }

  $async.Future<$1.Empty> adminRemoveUserFromGroup(
      $grpc.ServiceCall call, $0.AdminRemoveUserFromGroupRequest request);

  $async.Future<$0.AdminResetUserPasswordResponse> adminResetUserPassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminResetUserPasswordRequest> $request) async {
    return adminResetUserPassword($call, await $request);
  }

  $async.Future<$0.AdminResetUserPasswordResponse> adminResetUserPassword(
      $grpc.ServiceCall call, $0.AdminResetUserPasswordRequest request);

  $async.Future<$0.AdminRespondToAuthChallengeResponse>
      adminRespondToAuthChallenge_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AdminRespondToAuthChallengeRequest> $request) async {
    return adminRespondToAuthChallenge($call, await $request);
  }

  $async.Future<$0.AdminRespondToAuthChallengeResponse>
      adminRespondToAuthChallenge($grpc.ServiceCall call,
          $0.AdminRespondToAuthChallengeRequest request);

  $async.Future<$0.AdminSetUserMFAPreferenceResponse>
      adminSetUserMFAPreference_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AdminSetUserMFAPreferenceRequest> $request) async {
    return adminSetUserMFAPreference($call, await $request);
  }

  $async.Future<$0.AdminSetUserMFAPreferenceResponse> adminSetUserMFAPreference(
      $grpc.ServiceCall call, $0.AdminSetUserMFAPreferenceRequest request);

  $async.Future<$0.AdminSetUserPasswordResponse> adminSetUserPassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminSetUserPasswordRequest> $request) async {
    return adminSetUserPassword($call, await $request);
  }

  $async.Future<$0.AdminSetUserPasswordResponse> adminSetUserPassword(
      $grpc.ServiceCall call, $0.AdminSetUserPasswordRequest request);

  $async.Future<$0.AdminSetUserSettingsResponse> adminSetUserSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminSetUserSettingsRequest> $request) async {
    return adminSetUserSettings($call, await $request);
  }

  $async.Future<$0.AdminSetUserSettingsResponse> adminSetUserSettings(
      $grpc.ServiceCall call, $0.AdminSetUserSettingsRequest request);

  $async.Future<$0.AdminUpdateAuthEventFeedbackResponse>
      adminUpdateAuthEventFeedback_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.AdminUpdateAuthEventFeedbackRequest>
              $request) async {
    return adminUpdateAuthEventFeedback($call, await $request);
  }

  $async.Future<$0.AdminUpdateAuthEventFeedbackResponse>
      adminUpdateAuthEventFeedback($grpc.ServiceCall call,
          $0.AdminUpdateAuthEventFeedbackRequest request);

  $async.Future<$0.AdminUpdateDeviceStatusResponse> adminUpdateDeviceStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminUpdateDeviceStatusRequest> $request) async {
    return adminUpdateDeviceStatus($call, await $request);
  }

  $async.Future<$0.AdminUpdateDeviceStatusResponse> adminUpdateDeviceStatus(
      $grpc.ServiceCall call, $0.AdminUpdateDeviceStatusRequest request);

  $async.Future<$0.AdminUpdateUserAttributesResponse>
      adminUpdateUserAttributes_Pre($grpc.ServiceCall $call,
          $async.Future<$0.AdminUpdateUserAttributesRequest> $request) async {
    return adminUpdateUserAttributes($call, await $request);
  }

  $async.Future<$0.AdminUpdateUserAttributesResponse> adminUpdateUserAttributes(
      $grpc.ServiceCall call, $0.AdminUpdateUserAttributesRequest request);

  $async.Future<$0.AdminUserGlobalSignOutResponse> adminUserGlobalSignOut_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AdminUserGlobalSignOutRequest> $request) async {
    return adminUserGlobalSignOut($call, await $request);
  }

  $async.Future<$0.AdminUserGlobalSignOutResponse> adminUserGlobalSignOut(
      $grpc.ServiceCall call, $0.AdminUserGlobalSignOutRequest request);

  $async.Future<$0.AssociateSoftwareTokenResponse> associateSoftwareToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AssociateSoftwareTokenRequest> $request) async {
    return associateSoftwareToken($call, await $request);
  }

  $async.Future<$0.AssociateSoftwareTokenResponse> associateSoftwareToken(
      $grpc.ServiceCall call, $0.AssociateSoftwareTokenRequest request);

  $async.Future<$0.ChangePasswordResponse> changePassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ChangePasswordRequest> $request) async {
    return changePassword($call, await $request);
  }

  $async.Future<$0.ChangePasswordResponse> changePassword(
      $grpc.ServiceCall call, $0.ChangePasswordRequest request);

  $async.Future<$0.CompleteWebAuthnRegistrationResponse>
      completeWebAuthnRegistration_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CompleteWebAuthnRegistrationRequest>
              $request) async {
    return completeWebAuthnRegistration($call, await $request);
  }

  $async.Future<$0.CompleteWebAuthnRegistrationResponse>
      completeWebAuthnRegistration($grpc.ServiceCall call,
          $0.CompleteWebAuthnRegistrationRequest request);

  $async.Future<$0.ConfirmDeviceResponse> confirmDevice_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ConfirmDeviceRequest> $request) async {
    return confirmDevice($call, await $request);
  }

  $async.Future<$0.ConfirmDeviceResponse> confirmDevice(
      $grpc.ServiceCall call, $0.ConfirmDeviceRequest request);

  $async.Future<$0.ConfirmForgotPasswordResponse> confirmForgotPassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ConfirmForgotPasswordRequest> $request) async {
    return confirmForgotPassword($call, await $request);
  }

  $async.Future<$0.ConfirmForgotPasswordResponse> confirmForgotPassword(
      $grpc.ServiceCall call, $0.ConfirmForgotPasswordRequest request);

  $async.Future<$0.ConfirmSignUpResponse> confirmSignUp_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ConfirmSignUpRequest> $request) async {
    return confirmSignUp($call, await $request);
  }

  $async.Future<$0.ConfirmSignUpResponse> confirmSignUp(
      $grpc.ServiceCall call, $0.ConfirmSignUpRequest request);

  $async.Future<$0.CreateGroupResponse> createGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateGroupRequest> $request) async {
    return createGroup($call, await $request);
  }

  $async.Future<$0.CreateGroupResponse> createGroup(
      $grpc.ServiceCall call, $0.CreateGroupRequest request);

  $async.Future<$0.CreateIdentityProviderResponse> createIdentityProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateIdentityProviderRequest> $request) async {
    return createIdentityProvider($call, await $request);
  }

  $async.Future<$0.CreateIdentityProviderResponse> createIdentityProvider(
      $grpc.ServiceCall call, $0.CreateIdentityProviderRequest request);

  $async.Future<$0.CreateManagedLoginBrandingResponse>
      createManagedLoginBranding_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateManagedLoginBrandingRequest> $request) async {
    return createManagedLoginBranding($call, await $request);
  }

  $async.Future<$0.CreateManagedLoginBrandingResponse>
      createManagedLoginBranding(
          $grpc.ServiceCall call, $0.CreateManagedLoginBrandingRequest request);

  $async.Future<$0.CreateResourceServerResponse> createResourceServer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateResourceServerRequest> $request) async {
    return createResourceServer($call, await $request);
  }

  $async.Future<$0.CreateResourceServerResponse> createResourceServer(
      $grpc.ServiceCall call, $0.CreateResourceServerRequest request);

  $async.Future<$0.CreateTermsResponse> createTerms_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateTermsRequest> $request) async {
    return createTerms($call, await $request);
  }

  $async.Future<$0.CreateTermsResponse> createTerms(
      $grpc.ServiceCall call, $0.CreateTermsRequest request);

  $async.Future<$0.CreateUserImportJobResponse> createUserImportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateUserImportJobRequest> $request) async {
    return createUserImportJob($call, await $request);
  }

  $async.Future<$0.CreateUserImportJobResponse> createUserImportJob(
      $grpc.ServiceCall call, $0.CreateUserImportJobRequest request);

  $async.Future<$0.CreateUserPoolResponse> createUserPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateUserPoolRequest> $request) async {
    return createUserPool($call, await $request);
  }

  $async.Future<$0.CreateUserPoolResponse> createUserPool(
      $grpc.ServiceCall call, $0.CreateUserPoolRequest request);

  $async.Future<$0.CreateUserPoolClientResponse> createUserPoolClient_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateUserPoolClientRequest> $request) async {
    return createUserPoolClient($call, await $request);
  }

  $async.Future<$0.CreateUserPoolClientResponse> createUserPoolClient(
      $grpc.ServiceCall call, $0.CreateUserPoolClientRequest request);

  $async.Future<$0.CreateUserPoolDomainResponse> createUserPoolDomain_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateUserPoolDomainRequest> $request) async {
    return createUserPoolDomain($call, await $request);
  }

  $async.Future<$0.CreateUserPoolDomainResponse> createUserPoolDomain(
      $grpc.ServiceCall call, $0.CreateUserPoolDomainRequest request);

  $async.Future<$1.Empty> deleteGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteGroupRequest> $request) async {
    return deleteGroup($call, await $request);
  }

  $async.Future<$1.Empty> deleteGroup(
      $grpc.ServiceCall call, $0.DeleteGroupRequest request);

  $async.Future<$1.Empty> deleteIdentityProvider_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteIdentityProviderRequest> $request) async {
    return deleteIdentityProvider($call, await $request);
  }

  $async.Future<$1.Empty> deleteIdentityProvider(
      $grpc.ServiceCall call, $0.DeleteIdentityProviderRequest request);

  $async.Future<$1.Empty> deleteManagedLoginBranding_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteManagedLoginBrandingRequest> $request) async {
    return deleteManagedLoginBranding($call, await $request);
  }

  $async.Future<$1.Empty> deleteManagedLoginBranding(
      $grpc.ServiceCall call, $0.DeleteManagedLoginBrandingRequest request);

  $async.Future<$1.Empty> deleteResourceServer_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteResourceServerRequest> $request) async {
    return deleteResourceServer($call, await $request);
  }

  $async.Future<$1.Empty> deleteResourceServer(
      $grpc.ServiceCall call, $0.DeleteResourceServerRequest request);

  $async.Future<$1.Empty> deleteTerms_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteTermsRequest> $request) async {
    return deleteTerms($call, await $request);
  }

  $async.Future<$1.Empty> deleteTerms(
      $grpc.ServiceCall call, $0.DeleteTermsRequest request);

  $async.Future<$1.Empty> deleteUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserRequest> $request) async {
    return deleteUser($call, await $request);
  }

  $async.Future<$1.Empty> deleteUser(
      $grpc.ServiceCall call, $0.DeleteUserRequest request);

  $async.Future<$0.DeleteUserAttributesResponse> deleteUserAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserAttributesRequest> $request) async {
    return deleteUserAttributes($call, await $request);
  }

  $async.Future<$0.DeleteUserAttributesResponse> deleteUserAttributes(
      $grpc.ServiceCall call, $0.DeleteUserAttributesRequest request);

  $async.Future<$1.Empty> deleteUserPool_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserPoolRequest> $request) async {
    return deleteUserPool($call, await $request);
  }

  $async.Future<$1.Empty> deleteUserPool(
      $grpc.ServiceCall call, $0.DeleteUserPoolRequest request);

  $async.Future<$1.Empty> deleteUserPoolClient_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserPoolClientRequest> $request) async {
    return deleteUserPoolClient($call, await $request);
  }

  $async.Future<$1.Empty> deleteUserPoolClient(
      $grpc.ServiceCall call, $0.DeleteUserPoolClientRequest request);

  $async.Future<$0.DeleteUserPoolClientSecretResponse>
      deleteUserPoolClientSecret_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteUserPoolClientSecretRequest> $request) async {
    return deleteUserPoolClientSecret($call, await $request);
  }

  $async.Future<$0.DeleteUserPoolClientSecretResponse>
      deleteUserPoolClientSecret(
          $grpc.ServiceCall call, $0.DeleteUserPoolClientSecretRequest request);

  $async.Future<$0.DeleteUserPoolDomainResponse> deleteUserPoolDomain_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserPoolDomainRequest> $request) async {
    return deleteUserPoolDomain($call, await $request);
  }

  $async.Future<$0.DeleteUserPoolDomainResponse> deleteUserPoolDomain(
      $grpc.ServiceCall call, $0.DeleteUserPoolDomainRequest request);

  $async.Future<$0.DeleteWebAuthnCredentialResponse>
      deleteWebAuthnCredential_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DeleteWebAuthnCredentialRequest> $request) async {
    return deleteWebAuthnCredential($call, await $request);
  }

  $async.Future<$0.DeleteWebAuthnCredentialResponse> deleteWebAuthnCredential(
      $grpc.ServiceCall call, $0.DeleteWebAuthnCredentialRequest request);

  $async.Future<$0.DescribeIdentityProviderResponse>
      describeIdentityProvider_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeIdentityProviderRequest> $request) async {
    return describeIdentityProvider($call, await $request);
  }

  $async.Future<$0.DescribeIdentityProviderResponse> describeIdentityProvider(
      $grpc.ServiceCall call, $0.DescribeIdentityProviderRequest request);

  $async.Future<$0.DescribeManagedLoginBrandingResponse>
      describeManagedLoginBranding_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeManagedLoginBrandingRequest>
              $request) async {
    return describeManagedLoginBranding($call, await $request);
  }

  $async.Future<$0.DescribeManagedLoginBrandingResponse>
      describeManagedLoginBranding($grpc.ServiceCall call,
          $0.DescribeManagedLoginBrandingRequest request);

  $async.Future<$0.DescribeManagedLoginBrandingByClientResponse>
      describeManagedLoginBrandingByClient_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DescribeManagedLoginBrandingByClientRequest>
              $request) async {
    return describeManagedLoginBrandingByClient($call, await $request);
  }

  $async.Future<$0.DescribeManagedLoginBrandingByClientResponse>
      describeManagedLoginBrandingByClient($grpc.ServiceCall call,
          $0.DescribeManagedLoginBrandingByClientRequest request);

  $async.Future<$0.DescribeResourceServerResponse> describeResourceServer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeResourceServerRequest> $request) async {
    return describeResourceServer($call, await $request);
  }

  $async.Future<$0.DescribeResourceServerResponse> describeResourceServer(
      $grpc.ServiceCall call, $0.DescribeResourceServerRequest request);

  $async.Future<$0.DescribeRiskConfigurationResponse>
      describeRiskConfiguration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.DescribeRiskConfigurationRequest> $request) async {
    return describeRiskConfiguration($call, await $request);
  }

  $async.Future<$0.DescribeRiskConfigurationResponse> describeRiskConfiguration(
      $grpc.ServiceCall call, $0.DescribeRiskConfigurationRequest request);

  $async.Future<$0.DescribeTermsResponse> describeTerms_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeTermsRequest> $request) async {
    return describeTerms($call, await $request);
  }

  $async.Future<$0.DescribeTermsResponse> describeTerms(
      $grpc.ServiceCall call, $0.DescribeTermsRequest request);

  $async.Future<$0.DescribeUserImportJobResponse> describeUserImportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeUserImportJobRequest> $request) async {
    return describeUserImportJob($call, await $request);
  }

  $async.Future<$0.DescribeUserImportJobResponse> describeUserImportJob(
      $grpc.ServiceCall call, $0.DescribeUserImportJobRequest request);

  $async.Future<$0.DescribeUserPoolResponse> describeUserPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeUserPoolRequest> $request) async {
    return describeUserPool($call, await $request);
  }

  $async.Future<$0.DescribeUserPoolResponse> describeUserPool(
      $grpc.ServiceCall call, $0.DescribeUserPoolRequest request);

  $async.Future<$0.DescribeUserPoolClientResponse> describeUserPoolClient_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeUserPoolClientRequest> $request) async {
    return describeUserPoolClient($call, await $request);
  }

  $async.Future<$0.DescribeUserPoolClientResponse> describeUserPoolClient(
      $grpc.ServiceCall call, $0.DescribeUserPoolClientRequest request);

  $async.Future<$0.DescribeUserPoolDomainResponse> describeUserPoolDomain_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DescribeUserPoolDomainRequest> $request) async {
    return describeUserPoolDomain($call, await $request);
  }

  $async.Future<$0.DescribeUserPoolDomainResponse> describeUserPoolDomain(
      $grpc.ServiceCall call, $0.DescribeUserPoolDomainRequest request);

  $async.Future<$1.Empty> forgetDevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ForgetDeviceRequest> $request) async {
    return forgetDevice($call, await $request);
  }

  $async.Future<$1.Empty> forgetDevice(
      $grpc.ServiceCall call, $0.ForgetDeviceRequest request);

  $async.Future<$0.ForgotPasswordResponse> forgotPassword_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ForgotPasswordRequest> $request) async {
    return forgotPassword($call, await $request);
  }

  $async.Future<$0.ForgotPasswordResponse> forgotPassword(
      $grpc.ServiceCall call, $0.ForgotPasswordRequest request);

  $async.Future<$0.GetCSVHeaderResponse> getCSVHeader_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetCSVHeaderRequest> $request) async {
    return getCSVHeader($call, await $request);
  }

  $async.Future<$0.GetCSVHeaderResponse> getCSVHeader(
      $grpc.ServiceCall call, $0.GetCSVHeaderRequest request);

  $async.Future<$0.GetDeviceResponse> getDevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetDeviceRequest> $request) async {
    return getDevice($call, await $request);
  }

  $async.Future<$0.GetDeviceResponse> getDevice(
      $grpc.ServiceCall call, $0.GetDeviceRequest request);

  $async.Future<$0.GetGroupResponse> getGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetGroupRequest> $request) async {
    return getGroup($call, await $request);
  }

  $async.Future<$0.GetGroupResponse> getGroup(
      $grpc.ServiceCall call, $0.GetGroupRequest request);

  $async.Future<$0.GetIdentityProviderByIdentifierResponse>
      getIdentityProviderByIdentifier_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetIdentityProviderByIdentifierRequest>
              $request) async {
    return getIdentityProviderByIdentifier($call, await $request);
  }

  $async.Future<$0.GetIdentityProviderByIdentifierResponse>
      getIdentityProviderByIdentifier($grpc.ServiceCall call,
          $0.GetIdentityProviderByIdentifierRequest request);

  $async.Future<$0.GetLogDeliveryConfigurationResponse>
      getLogDeliveryConfiguration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetLogDeliveryConfigurationRequest> $request) async {
    return getLogDeliveryConfiguration($call, await $request);
  }

  $async.Future<$0.GetLogDeliveryConfigurationResponse>
      getLogDeliveryConfiguration($grpc.ServiceCall call,
          $0.GetLogDeliveryConfigurationRequest request);

  $async.Future<$0.GetSigningCertificateResponse> getSigningCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSigningCertificateRequest> $request) async {
    return getSigningCertificate($call, await $request);
  }

  $async.Future<$0.GetSigningCertificateResponse> getSigningCertificate(
      $grpc.ServiceCall call, $0.GetSigningCertificateRequest request);

  $async.Future<$0.GetTokensFromRefreshTokenResponse>
      getTokensFromRefreshToken_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetTokensFromRefreshTokenRequest> $request) async {
    return getTokensFromRefreshToken($call, await $request);
  }

  $async.Future<$0.GetTokensFromRefreshTokenResponse> getTokensFromRefreshToken(
      $grpc.ServiceCall call, $0.GetTokensFromRefreshTokenRequest request);

  $async.Future<$0.GetUICustomizationResponse> getUICustomization_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetUICustomizationRequest> $request) async {
    return getUICustomization($call, await $request);
  }

  $async.Future<$0.GetUICustomizationResponse> getUICustomization(
      $grpc.ServiceCall call, $0.GetUICustomizationRequest request);

  $async.Future<$0.GetUserResponse> getUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetUserRequest> $request) async {
    return getUser($call, await $request);
  }

  $async.Future<$0.GetUserResponse> getUser(
      $grpc.ServiceCall call, $0.GetUserRequest request);

  $async.Future<$0.GetUserAttributeVerificationCodeResponse>
      getUserAttributeVerificationCode_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetUserAttributeVerificationCodeRequest>
              $request) async {
    return getUserAttributeVerificationCode($call, await $request);
  }

  $async.Future<$0.GetUserAttributeVerificationCodeResponse>
      getUserAttributeVerificationCode($grpc.ServiceCall call,
          $0.GetUserAttributeVerificationCodeRequest request);

  $async.Future<$0.GetUserAuthFactorsResponse> getUserAuthFactors_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetUserAuthFactorsRequest> $request) async {
    return getUserAuthFactors($call, await $request);
  }

  $async.Future<$0.GetUserAuthFactorsResponse> getUserAuthFactors(
      $grpc.ServiceCall call, $0.GetUserAuthFactorsRequest request);

  $async.Future<$0.GetUserPoolMfaConfigResponse> getUserPoolMfaConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetUserPoolMfaConfigRequest> $request) async {
    return getUserPoolMfaConfig($call, await $request);
  }

  $async.Future<$0.GetUserPoolMfaConfigResponse> getUserPoolMfaConfig(
      $grpc.ServiceCall call, $0.GetUserPoolMfaConfigRequest request);

  $async.Future<$0.GlobalSignOutResponse> globalSignOut_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GlobalSignOutRequest> $request) async {
    return globalSignOut($call, await $request);
  }

  $async.Future<$0.GlobalSignOutResponse> globalSignOut(
      $grpc.ServiceCall call, $0.GlobalSignOutRequest request);

  $async.Future<$0.InitiateAuthResponse> initiateAuth_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.InitiateAuthRequest> $request) async {
    return initiateAuth($call, await $request);
  }

  $async.Future<$0.InitiateAuthResponse> initiateAuth(
      $grpc.ServiceCall call, $0.InitiateAuthRequest request);

  $async.Future<$0.ListDevicesResponse> listDevices_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListDevicesRequest> $request) async {
    return listDevices($call, await $request);
  }

  $async.Future<$0.ListDevicesResponse> listDevices(
      $grpc.ServiceCall call, $0.ListDevicesRequest request);

  $async.Future<$0.ListGroupsResponse> listGroups_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListGroupsRequest> $request) async {
    return listGroups($call, await $request);
  }

  $async.Future<$0.ListGroupsResponse> listGroups(
      $grpc.ServiceCall call, $0.ListGroupsRequest request);

  $async.Future<$0.ListIdentityProvidersResponse> listIdentityProviders_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListIdentityProvidersRequest> $request) async {
    return listIdentityProviders($call, await $request);
  }

  $async.Future<$0.ListIdentityProvidersResponse> listIdentityProviders(
      $grpc.ServiceCall call, $0.ListIdentityProvidersRequest request);

  $async.Future<$0.ListResourceServersResponse> listResourceServers_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListResourceServersRequest> $request) async {
    return listResourceServers($call, await $request);
  }

  $async.Future<$0.ListResourceServersResponse> listResourceServers(
      $grpc.ServiceCall call, $0.ListResourceServersRequest request);

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListTagsForResourceRequest> $request) async {
    return listTagsForResource($call, await $request);
  }

  $async.Future<$0.ListTagsForResourceResponse> listTagsForResource(
      $grpc.ServiceCall call, $0.ListTagsForResourceRequest request);

  $async.Future<$0.ListTermsResponse> listTerms_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListTermsRequest> $request) async {
    return listTerms($call, await $request);
  }

  $async.Future<$0.ListTermsResponse> listTerms(
      $grpc.ServiceCall call, $0.ListTermsRequest request);

  $async.Future<$0.ListUserImportJobsResponse> listUserImportJobs_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListUserImportJobsRequest> $request) async {
    return listUserImportJobs($call, await $request);
  }

  $async.Future<$0.ListUserImportJobsResponse> listUserImportJobs(
      $grpc.ServiceCall call, $0.ListUserImportJobsRequest request);

  $async.Future<$0.ListUserPoolClientsResponse> listUserPoolClients_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListUserPoolClientsRequest> $request) async {
    return listUserPoolClients($call, await $request);
  }

  $async.Future<$0.ListUserPoolClientsResponse> listUserPoolClients(
      $grpc.ServiceCall call, $0.ListUserPoolClientsRequest request);

  $async.Future<$0.ListUserPoolClientSecretsResponse>
      listUserPoolClientSecrets_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListUserPoolClientSecretsRequest> $request) async {
    return listUserPoolClientSecrets($call, await $request);
  }

  $async.Future<$0.ListUserPoolClientSecretsResponse> listUserPoolClientSecrets(
      $grpc.ServiceCall call, $0.ListUserPoolClientSecretsRequest request);

  $async.Future<$0.ListUserPoolsResponse> listUserPools_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListUserPoolsRequest> $request) async {
    return listUserPools($call, await $request);
  }

  $async.Future<$0.ListUserPoolsResponse> listUserPools(
      $grpc.ServiceCall call, $0.ListUserPoolsRequest request);

  $async.Future<$0.ListUsersResponse> listUsers_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListUsersRequest> $request) async {
    return listUsers($call, await $request);
  }

  $async.Future<$0.ListUsersResponse> listUsers(
      $grpc.ServiceCall call, $0.ListUsersRequest request);

  $async.Future<$0.ListUsersInGroupResponse> listUsersInGroup_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListUsersInGroupRequest> $request) async {
    return listUsersInGroup($call, await $request);
  }

  $async.Future<$0.ListUsersInGroupResponse> listUsersInGroup(
      $grpc.ServiceCall call, $0.ListUsersInGroupRequest request);

  $async.Future<$0.ListWebAuthnCredentialsResponse> listWebAuthnCredentials_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListWebAuthnCredentialsRequest> $request) async {
    return listWebAuthnCredentials($call, await $request);
  }

  $async.Future<$0.ListWebAuthnCredentialsResponse> listWebAuthnCredentials(
      $grpc.ServiceCall call, $0.ListWebAuthnCredentialsRequest request);

  $async.Future<$0.ResendConfirmationCodeResponse> resendConfirmationCode_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ResendConfirmationCodeRequest> $request) async {
    return resendConfirmationCode($call, await $request);
  }

  $async.Future<$0.ResendConfirmationCodeResponse> resendConfirmationCode(
      $grpc.ServiceCall call, $0.ResendConfirmationCodeRequest request);

  $async.Future<$0.RespondToAuthChallengeResponse> respondToAuthChallenge_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RespondToAuthChallengeRequest> $request) async {
    return respondToAuthChallenge($call, await $request);
  }

  $async.Future<$0.RespondToAuthChallengeResponse> respondToAuthChallenge(
      $grpc.ServiceCall call, $0.RespondToAuthChallengeRequest request);

  $async.Future<$0.RevokeTokenResponse> revokeToken_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RevokeTokenRequest> $request) async {
    return revokeToken($call, await $request);
  }

  $async.Future<$0.RevokeTokenResponse> revokeToken(
      $grpc.ServiceCall call, $0.RevokeTokenRequest request);

  $async.Future<$0.SetLogDeliveryConfigurationResponse>
      setLogDeliveryConfiguration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.SetLogDeliveryConfigurationRequest> $request) async {
    return setLogDeliveryConfiguration($call, await $request);
  }

  $async.Future<$0.SetLogDeliveryConfigurationResponse>
      setLogDeliveryConfiguration($grpc.ServiceCall call,
          $0.SetLogDeliveryConfigurationRequest request);

  $async.Future<$0.SetRiskConfigurationResponse> setRiskConfiguration_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetRiskConfigurationRequest> $request) async {
    return setRiskConfiguration($call, await $request);
  }

  $async.Future<$0.SetRiskConfigurationResponse> setRiskConfiguration(
      $grpc.ServiceCall call, $0.SetRiskConfigurationRequest request);

  $async.Future<$0.SetUICustomizationResponse> setUICustomization_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetUICustomizationRequest> $request) async {
    return setUICustomization($call, await $request);
  }

  $async.Future<$0.SetUICustomizationResponse> setUICustomization(
      $grpc.ServiceCall call, $0.SetUICustomizationRequest request);

  $async.Future<$0.SetUserMFAPreferenceResponse> setUserMFAPreference_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetUserMFAPreferenceRequest> $request) async {
    return setUserMFAPreference($call, await $request);
  }

  $async.Future<$0.SetUserMFAPreferenceResponse> setUserMFAPreference(
      $grpc.ServiceCall call, $0.SetUserMFAPreferenceRequest request);

  $async.Future<$0.SetUserPoolMfaConfigResponse> setUserPoolMfaConfig_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetUserPoolMfaConfigRequest> $request) async {
    return setUserPoolMfaConfig($call, await $request);
  }

  $async.Future<$0.SetUserPoolMfaConfigResponse> setUserPoolMfaConfig(
      $grpc.ServiceCall call, $0.SetUserPoolMfaConfigRequest request);

  $async.Future<$0.SetUserSettingsResponse> setUserSettings_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetUserSettingsRequest> $request) async {
    return setUserSettings($call, await $request);
  }

  $async.Future<$0.SetUserSettingsResponse> setUserSettings(
      $grpc.ServiceCall call, $0.SetUserSettingsRequest request);

  $async.Future<$0.SignUpResponse> signUp_Pre(
      $grpc.ServiceCall $call, $async.Future<$0.SignUpRequest> $request) async {
    return signUp($call, await $request);
  }

  $async.Future<$0.SignUpResponse> signUp(
      $grpc.ServiceCall call, $0.SignUpRequest request);

  $async.Future<$0.StartUserImportJobResponse> startUserImportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StartUserImportJobRequest> $request) async {
    return startUserImportJob($call, await $request);
  }

  $async.Future<$0.StartUserImportJobResponse> startUserImportJob(
      $grpc.ServiceCall call, $0.StartUserImportJobRequest request);

  $async.Future<$0.StartWebAuthnRegistrationResponse>
      startWebAuthnRegistration_Pre($grpc.ServiceCall $call,
          $async.Future<$0.StartWebAuthnRegistrationRequest> $request) async {
    return startWebAuthnRegistration($call, await $request);
  }

  $async.Future<$0.StartWebAuthnRegistrationResponse> startWebAuthnRegistration(
      $grpc.ServiceCall call, $0.StartWebAuthnRegistrationRequest request);

  $async.Future<$0.StopUserImportJobResponse> stopUserImportJob_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.StopUserImportJobRequest> $request) async {
    return stopUserImportJob($call, await $request);
  }

  $async.Future<$0.StopUserImportJobResponse> stopUserImportJob(
      $grpc.ServiceCall call, $0.StopUserImportJobRequest request);

  $async.Future<$0.TagResourceResponse> tagResource_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagResourceRequest> $request) async {
    return tagResource($call, await $request);
  }

  $async.Future<$0.TagResourceResponse> tagResource(
      $grpc.ServiceCall call, $0.TagResourceRequest request);

  $async.Future<$0.UntagResourceResponse> untagResource_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagResourceRequest> $request) async {
    return untagResource($call, await $request);
  }

  $async.Future<$0.UntagResourceResponse> untagResource(
      $grpc.ServiceCall call, $0.UntagResourceRequest request);

  $async.Future<$0.UpdateAuthEventFeedbackResponse> updateAuthEventFeedback_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateAuthEventFeedbackRequest> $request) async {
    return updateAuthEventFeedback($call, await $request);
  }

  $async.Future<$0.UpdateAuthEventFeedbackResponse> updateAuthEventFeedback(
      $grpc.ServiceCall call, $0.UpdateAuthEventFeedbackRequest request);

  $async.Future<$0.UpdateDeviceStatusResponse> updateDeviceStatus_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateDeviceStatusRequest> $request) async {
    return updateDeviceStatus($call, await $request);
  }

  $async.Future<$0.UpdateDeviceStatusResponse> updateDeviceStatus(
      $grpc.ServiceCall call, $0.UpdateDeviceStatusRequest request);

  $async.Future<$0.UpdateGroupResponse> updateGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateGroupRequest> $request) async {
    return updateGroup($call, await $request);
  }

  $async.Future<$0.UpdateGroupResponse> updateGroup(
      $grpc.ServiceCall call, $0.UpdateGroupRequest request);

  $async.Future<$0.UpdateIdentityProviderResponse> updateIdentityProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateIdentityProviderRequest> $request) async {
    return updateIdentityProvider($call, await $request);
  }

  $async.Future<$0.UpdateIdentityProviderResponse> updateIdentityProvider(
      $grpc.ServiceCall call, $0.UpdateIdentityProviderRequest request);

  $async.Future<$0.UpdateManagedLoginBrandingResponse>
      updateManagedLoginBranding_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UpdateManagedLoginBrandingRequest> $request) async {
    return updateManagedLoginBranding($call, await $request);
  }

  $async.Future<$0.UpdateManagedLoginBrandingResponse>
      updateManagedLoginBranding(
          $grpc.ServiceCall call, $0.UpdateManagedLoginBrandingRequest request);

  $async.Future<$0.UpdateResourceServerResponse> updateResourceServer_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateResourceServerRequest> $request) async {
    return updateResourceServer($call, await $request);
  }

  $async.Future<$0.UpdateResourceServerResponse> updateResourceServer(
      $grpc.ServiceCall call, $0.UpdateResourceServerRequest request);

  $async.Future<$0.UpdateTermsResponse> updateTerms_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateTermsRequest> $request) async {
    return updateTerms($call, await $request);
  }

  $async.Future<$0.UpdateTermsResponse> updateTerms(
      $grpc.ServiceCall call, $0.UpdateTermsRequest request);

  $async.Future<$0.UpdateUserAttributesResponse> updateUserAttributes_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateUserAttributesRequest> $request) async {
    return updateUserAttributes($call, await $request);
  }

  $async.Future<$0.UpdateUserAttributesResponse> updateUserAttributes(
      $grpc.ServiceCall call, $0.UpdateUserAttributesRequest request);

  $async.Future<$0.UpdateUserPoolResponse> updateUserPool_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateUserPoolRequest> $request) async {
    return updateUserPool($call, await $request);
  }

  $async.Future<$0.UpdateUserPoolResponse> updateUserPool(
      $grpc.ServiceCall call, $0.UpdateUserPoolRequest request);

  $async.Future<$0.UpdateUserPoolClientResponse> updateUserPoolClient_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateUserPoolClientRequest> $request) async {
    return updateUserPoolClient($call, await $request);
  }

  $async.Future<$0.UpdateUserPoolClientResponse> updateUserPoolClient(
      $grpc.ServiceCall call, $0.UpdateUserPoolClientRequest request);

  $async.Future<$0.UpdateUserPoolDomainResponse> updateUserPoolDomain_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateUserPoolDomainRequest> $request) async {
    return updateUserPoolDomain($call, await $request);
  }

  $async.Future<$0.UpdateUserPoolDomainResponse> updateUserPoolDomain(
      $grpc.ServiceCall call, $0.UpdateUserPoolDomainRequest request);

  $async.Future<$0.VerifySoftwareTokenResponse> verifySoftwareToken_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.VerifySoftwareTokenRequest> $request) async {
    return verifySoftwareToken($call, await $request);
  }

  $async.Future<$0.VerifySoftwareTokenResponse> verifySoftwareToken(
      $grpc.ServiceCall call, $0.VerifySoftwareTokenRequest request);

  $async.Future<$0.VerifyUserAttributeResponse> verifyUserAttribute_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.VerifyUserAttributeRequest> $request) async {
    return verifyUserAttribute($call, await $request);
  }

  $async.Future<$0.VerifyUserAttributeResponse> verifyUserAttribute(
      $grpc.ServiceCall call, $0.VerifyUserAttributeRequest request);
}
