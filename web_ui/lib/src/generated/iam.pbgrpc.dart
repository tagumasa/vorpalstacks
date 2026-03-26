// This is a generated file - do not edit.
//
// Generated from iam.proto.

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
import 'iam.pb.dart' as $0;

export 'iam.pb.dart';

/// IAMService provides iam API operations.
@$pb.GrpcServiceName('iam.IAMService')
class IAMServiceClient extends $grpc.Client {
  /// The hostname for this service.
  static const $core.String defaultHost = '';

  /// OAuth scopes needed for the client.
  static const $core.List<$core.String> oauthScopes = [
    '',
  ];

  IAMServiceClient(super.channel, {super.options, super.interceptors});

  /// Accepts a delegation request, granting the requested temporary access. Once the delegation request is accepted, it is eligible to send the exchange token to the partner. The SendDelegationToken API...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> acceptDelegationRequest(
    $0.AcceptDelegationRequestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$acceptDelegationRequest, request,
        options: options);
  }

  /// Adds a new client ID (also known as audience) to the list of client IDs already registered for the specified IAM OpenID Connect (OIDC) provider resource. This operation is idempotent; it does not f...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> addClientIDToOpenIDConnectProvider(
    $0.AddClientIDToOpenIDConnectProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addClientIDToOpenIDConnectProvider, request,
        options: options);
  }

  /// Adds the specified IAM role to the specified instance profile. An instance profile can contain only one role, and this quota cannot be increased. You can remove the existing role and then add a dif...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> addRoleToInstanceProfile(
    $0.AddRoleToInstanceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addRoleToInstanceProfile, request,
        options: options);
  }

  /// Adds the specified user to the specified group.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> addUserToGroup(
    $0.AddUserToGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$addUserToGroup, request, options: options);
  }

  /// Associates a delegation request with the current identity. If the partner that created the delegation request has specified the owner account during creation, only an identity from that owner accou...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> associateDelegationRequest(
    $0.AssociateDelegationRequestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$associateDelegationRequest, request,
        options: options);
  }

  /// Attaches the specified managed policy to the specified IAM group. You use this operation to attach a managed policy to a group. To embed an inline policy in a group, use PutGroupPolicy . As a best ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> attachGroupPolicy(
    $0.AttachGroupPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$attachGroupPolicy, request, options: options);
  }

  /// Attaches the specified managed policy to the specified IAM role. When you attach a managed policy to a role, the managed policy becomes part of the role's permission (access) policy. You cannot use...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> attachRolePolicy(
    $0.AttachRolePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$attachRolePolicy, request, options: options);
  }

  /// Attaches the specified managed policy to the specified user. You use this operation to attach a managed policy to a user. To embed an inline policy in a user, use PutUserPolicy . As a best practice...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> attachUserPolicy(
    $0.AttachUserPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$attachUserPolicy, request, options: options);
  }

  /// Changes the password of the IAM user who is calling this operation. This operation can be performed using the CLI, the Amazon Web Services API, or the My Security Credentials page in the Amazon Web...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> changePassword(
    $0.ChangePasswordRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$changePassword, request, options: options);
  }

  /// Creates a new Amazon Web Services secret access key and corresponding Amazon Web Services access key ID for the specified user. The default status for new keys is Active. If you do not specify a us...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateAccessKeyResponse> createAccessKey(
    $0.CreateAccessKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAccessKey, request, options: options);
  }

  /// Creates an alias for your Amazon Web Services account. For information about using an Amazon Web Services account alias, see Creating, deleting, and listing an Amazon Web Services account alias in ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> createAccountAlias(
    $0.CreateAccountAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createAccountAlias, request, options: options);
  }

  /// Creates an IAM delegation request for temporary access delegation. This API is not available for general use. In order to use this API, a caller first need to go through an onboarding process descr...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateDelegationRequestResponse>
      createDelegationRequest(
    $0.CreateDelegationRequestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createDelegationRequest, request,
        options: options);
  }

  /// Creates a new group. For information about the number of groups you can create, see IAM and STS quotas in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateGroupResponse> createGroup(
    $0.CreateGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createGroup, request, options: options);
  }

  /// Creates a new instance profile. For information about instance profiles, see Using roles for applications on Amazon EC2 in the IAM User Guide, and Instance profiles in the Amazon EC2 User Guide. Fo...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateInstanceProfileResponse> createInstanceProfile(
    $0.CreateInstanceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createInstanceProfile, request, options: options);
  }

  /// Creates a password for the specified IAM user. A password allows an IAM user to access Amazon Web Services services through the Amazon Web Services Management Console. You can use the CLI, the Amaz...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateLoginProfileResponse> createLoginProfile(
    $0.CreateLoginProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createLoginProfile, request, options: options);
  }

  /// Creates an IAM entity to describe an identity provider (IdP) that supports OpenID Connect (OIDC). The OIDC provider that you create with this operation can be used as a principal in a role's trust ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateOpenIDConnectProviderResponse>
      createOpenIDConnectProvider(
    $0.CreateOpenIDConnectProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createOpenIDConnectProvider, request,
        options: options);
  }

  /// Creates a new managed policy for your Amazon Web Services account. This operation creates a policy version with a version identifier of v1 and sets v1 as the policy's default version. For more info...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreatePolicyResponse> createPolicy(
    $0.CreatePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPolicy, request, options: options);
  }

  /// Creates a new version of the specified managed policy. To update a managed policy, you create a new policy version. A managed policy can have up to five versions. If the policy has five versions, y...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreatePolicyVersionResponse> createPolicyVersion(
    $0.CreatePolicyVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createPolicyVersion, request, options: options);
  }

  /// Creates a new role for your Amazon Web Services account. For more information about roles, see IAM roles in the IAM User Guide. For information about quotas for role names and the number of roles y...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateRoleResponse> createRole(
    $0.CreateRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createRole, request, options: options);
  }

  /// Creates an IAM resource that describes an identity provider (IdP) that supports SAML 2.0. The SAML provider resource that you create with this operation can be used as a principal in an IAM role's ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateSAMLProviderResponse> createSAMLProvider(
    $0.CreateSAMLProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createSAMLProvider, request, options: options);
  }

  /// Creates an IAM role that is linked to a specific Amazon Web Services service. The service controls the attached policies and when the role can be deleted. This helps ensure that the service is not ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateServiceLinkedRoleResponse>
      createServiceLinkedRole(
    $0.CreateServiceLinkedRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createServiceLinkedRole, request,
        options: options);
  }

  /// Generates a set of credentials consisting of a user name and password that can be used to access the service specified in the request. These credentials are generated by IAM, and can be used only f...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateServiceSpecificCredentialResponse>
      createServiceSpecificCredential(
    $0.CreateServiceSpecificCredentialRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createServiceSpecificCredential, request,
        options: options);
  }

  /// Creates a new IAM user for your Amazon Web Services account. For information about quotas for the number of IAM users you can create, see IAM and STS quotas in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateUserResponse> createUser(
    $0.CreateUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createUser, request, options: options);
  }

  /// Creates a new virtual MFA device for the Amazon Web Services account. After creating the virtual MFA, use EnableMFADevice to attach the MFA device to an IAM user. For more information about creatin...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.CreateVirtualMFADeviceResponse>
      createVirtualMFADevice(
    $0.CreateVirtualMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$createVirtualMFADevice, request,
        options: options);
  }

  /// Deactivates the specified MFA device and removes it from association with the user name for which it was originally enabled. For more information about creating and working with virtual MFA devices...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deactivateMFADevice(
    $0.DeactivateMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deactivateMFADevice, request, options: options);
  }

  /// Deletes the access key pair associated with the specified IAM user. If you do not specify a user name, IAM determines the user name implicitly based on the Amazon Web Services access key ID signing...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteAccessKey(
    $0.DeleteAccessKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAccessKey, request, options: options);
  }

  /// Deletes the specified Amazon Web Services account alias. For information about using an Amazon Web Services account alias, see Creating, deleting, and listing an Amazon Web Services account alias i...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteAccountAlias(
    $0.DeleteAccountAliasRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAccountAlias, request, options: options);
  }

  /// Deletes the password policy for the Amazon Web Services account. There are no parameters.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteAccountPasswordPolicy(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteAccountPasswordPolicy, request,
        options: options);
  }

  /// Deletes the specified IAM group. The group must not contain any users or have any attached policies.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteGroup(
    $0.DeleteGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteGroup, request, options: options);
  }

  /// Deletes the specified inline policy that is embedded in the specified IAM group. A group can also have managed policies attached to it. To detach a managed policy from a group, use DetachGroupPolic...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteGroupPolicy(
    $0.DeleteGroupPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteGroupPolicy, request, options: options);
  }

  /// Deletes the specified instance profile. The instance profile must not have an associated role. Make sure that you do not have any Amazon EC2 instances running with the instance profile you are abou...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteInstanceProfile(
    $0.DeleteInstanceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteInstanceProfile, request, options: options);
  }

  /// Deletes the password for the specified IAM user or root user, For more information, see Managing passwords for IAM users. You can use the CLI, the Amazon Web Services API, or the Users page in the ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteLoginProfile(
    $0.DeleteLoginProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteLoginProfile, request, options: options);
  }

  /// Deletes an OpenID Connect identity provider (IdP) resource object in IAM. Deleting an IAM OIDC provider resource does not update any roles that reference the provider as a principal in their trust ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteOpenIDConnectProvider(
    $0.DeleteOpenIDConnectProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteOpenIDConnectProvider, request,
        options: options);
  }

  /// Deletes the specified managed policy. Before you can delete a managed policy, you must first detach the policy from all users, groups, and roles that it is attached to. In addition, you must delete...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deletePolicy(
    $0.DeletePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePolicy, request, options: options);
  }

  /// Deletes the specified version from the specified managed policy. You cannot delete the default version from a policy using this operation. To delete the default version from a policy, use DeletePol...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deletePolicyVersion(
    $0.DeletePolicyVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deletePolicyVersion, request, options: options);
  }

  /// Deletes the specified role. Unlike the Amazon Web Services Management Console, when you delete a role programmatically, you must delete the items attached to the role manually, or the deletion fail...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteRole(
    $0.DeleteRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRole, request, options: options);
  }

  /// Deletes the permissions boundary for the specified IAM role. You cannot set the boundary for a service-linked role. Deleting the permissions boundary for a role might increase its permissions. For ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteRolePermissionsBoundary(
    $0.DeleteRolePermissionsBoundaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRolePermissionsBoundary, request,
        options: options);
  }

  /// Deletes the specified inline policy that is embedded in the specified IAM role. A role can also have managed policies attached to it. To detach a managed policy from a role, use DetachRolePolicy. F...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteRolePolicy(
    $0.DeleteRolePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteRolePolicy, request, options: options);
  }

  /// Deletes a SAML provider resource in IAM. Deleting the provider resource from IAM does not update any roles that reference the SAML provider resource's ARN as a principal in their trust policies. An...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteSAMLProvider(
    $0.DeleteSAMLProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSAMLProvider, request, options: options);
  }

  /// Deletes the specified server certificate. For more information about working with server certificates, see Working with server certificates in the IAM User Guide. This topic also includes a list of...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteServerCertificate(
    $0.DeleteServerCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteServerCertificate, request,
        options: options);
  }

  /// Submits a service-linked role deletion request and returns a DeletionTaskId, which you can use to check the status of the deletion. Before you call this operation, confirm that the role has no acti...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.DeleteServiceLinkedRoleResponse>
      deleteServiceLinkedRole(
    $0.DeleteServiceLinkedRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteServiceLinkedRole, request,
        options: options);
  }

  /// Deletes the specified service-specific credential.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteServiceSpecificCredential(
    $0.DeleteServiceSpecificCredentialRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteServiceSpecificCredential, request,
        options: options);
  }

  /// Deletes a signing certificate associated with the specified IAM user. If you do not specify a user name, IAM determines the user name implicitly based on the Amazon Web Services access key ID signi...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteSigningCertificate(
    $0.DeleteSigningCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSigningCertificate, request,
        options: options);
  }

  /// Deletes the specified SSH public key. The SSH public key deleted by this operation is used only for authenticating the associated IAM user to an CodeCommit repository. For more information about us...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteSSHPublicKey(
    $0.DeleteSSHPublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteSSHPublicKey, request, options: options);
  }

  /// Deletes the specified IAM user. Unlike the Amazon Web Services Management Console, when you delete a user programmatically, you must delete the items attached to the user manually, or the deletion ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteUser(
    $0.DeleteUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUser, request, options: options);
  }

  /// Deletes the permissions boundary for the specified IAM user. Deleting the permissions boundary for a user might increase its permissions by allowing the user to perform all the actions granted in i...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteUserPermissionsBoundary(
    $0.DeleteUserPermissionsBoundaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUserPermissionsBoundary, request,
        options: options);
  }

  /// Deletes the specified inline policy that is embedded in the specified IAM user. A user can also have managed policies attached to it. To detach a managed policy from a user, use DetachUserPolicy. F...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteUserPolicy(
    $0.DeleteUserPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteUserPolicy, request, options: options);
  }

  /// Deletes a virtual MFA device. You must deactivate a user's virtual MFA device before you can delete it. For information about deactivating MFA devices, see DeactivateMFADevice.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> deleteVirtualMFADevice(
    $0.DeleteVirtualMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$deleteVirtualMFADevice, request,
        options: options);
  }

  /// Removes the specified managed policy from the specified IAM group. A group can also have inline policies embedded with it. To delete an inline policy, use DeleteGroupPolicy. For information about p...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> detachGroupPolicy(
    $0.DetachGroupPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$detachGroupPolicy, request, options: options);
  }

  /// Removes the specified managed policy from the specified role. A role can also have inline policies embedded with it. To delete an inline policy, use DeleteRolePolicy. For information about policies...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> detachRolePolicy(
    $0.DetachRolePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$detachRolePolicy, request, options: options);
  }

  /// Removes the specified managed policy from the specified user. A user can also have inline policies embedded with it. To delete an inline policy, use DeleteUserPolicy. For information about policies...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> detachUserPolicy(
    $0.DetachUserPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$detachUserPolicy, request, options: options);
  }

  /// Disables the management of privileged root user credentials across member accounts in your organization. When you disable this feature, the management account and the delegated administrator for IA...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.DisableOrganizationsRootCredentialsManagementResponse>
      disableOrganizationsRootCredentialsManagement(
    $0.DisableOrganizationsRootCredentialsManagementRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$disableOrganizationsRootCredentialsManagement, request,
        options: options);
  }

  /// Disables root user sessions for privileged tasks across member accounts in your organization. When you disable this feature, the management account and the delegated administrator for IAM can no lo...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.DisableOrganizationsRootSessionsResponse>
      disableOrganizationsRootSessions(
    $0.DisableOrganizationsRootSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableOrganizationsRootSessions, request,
        options: options);
  }

  /// Disables the outbound identity federation feature for your Amazon Web Services account. When disabled, IAM principals in the account cannot use the GetWebIdentityToken API to obtain JSON Web Tokens...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> disableOutboundWebIdentityFederation(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$disableOutboundWebIdentityFederation, request,
        options: options);
  }

  /// Enables the specified MFA device and associates it with the specified IAM user. When enabled, the MFA device is required for every subsequent login by the IAM user associated with the device.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> enableMFADevice(
    $0.EnableMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableMFADevice, request, options: options);
  }

  /// Enables the management of privileged root user credentials across member accounts in your organization. When you enable root credentials management for centralized root access, the management accou...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.EnableOrganizationsRootCredentialsManagementResponse>
      enableOrganizationsRootCredentialsManagement(
    $0.EnableOrganizationsRootCredentialsManagementRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$enableOrganizationsRootCredentialsManagement, request,
        options: options);
  }

  /// Allows the management account or delegated administrator to perform privileged tasks on member accounts in your organization. For more information, see Centrally manage root access for member accou...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.EnableOrganizationsRootSessionsResponse>
      enableOrganizationsRootSessions(
    $0.EnableOrganizationsRootSessionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableOrganizationsRootSessions, request,
        options: options);
  }

  /// Enables the outbound identity federation feature for your Amazon Web Services account. When enabled, IAM principals in your account can use the GetWebIdentityToken API to obtain JSON Web Tokens (JW...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.EnableOutboundWebIdentityFederationResponse>
      enableOutboundWebIdentityFederation(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$enableOutboundWebIdentityFederation, request,
        options: options);
  }

  /// Generates a credential report for the Amazon Web Services account. For more information about the credential report, see Getting credential reports in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GenerateCredentialReportResponse>
      generateCredentialReport(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateCredentialReport, request,
        options: options);
  }

  /// Generates a report for service last accessed data for Organizations. You can generate a report for any entities (organization root, organizational unit, or account) or policies in your organization...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GenerateOrganizationsAccessReportResponse>
      generateOrganizationsAccessReport(
    $0.GenerateOrganizationsAccessReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateOrganizationsAccessReport, request,
        options: options);
  }

  /// Generates a report that includes details about when an IAM resource (user, group, role, or policy) was last used in an attempt to access Amazon Web Services services. Recent activity usually appear...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GenerateServiceLastAccessedDetailsResponse>
      generateServiceLastAccessedDetails(
    $0.GenerateServiceLastAccessedDetailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$generateServiceLastAccessedDetails, request,
        options: options);
  }

  /// Retrieves information about when the specified access key was last used. The information includes the date and time of last use, along with the Amazon Web Services service and Region that were spec...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetAccessKeyLastUsedResponse> getAccessKeyLastUsed(
    $0.GetAccessKeyLastUsedRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccessKeyLastUsed, request, options: options);
  }

  /// Retrieves information about all IAM users, groups, roles, and policies in your Amazon Web Services account, including their relationships to one another. Use this operation to obtain a snapshot of ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetAccountAuthorizationDetailsResponse>
      getAccountAuthorizationDetails(
    $0.GetAccountAuthorizationDetailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccountAuthorizationDetails, request,
        options: options);
  }

  /// Retrieves the password policy for the Amazon Web Services account. This tells you the complexity requirements and mandatory rotation periods for the IAM user passwords in your account. For more inf...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetAccountPasswordPolicyResponse>
      getAccountPasswordPolicy(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccountPasswordPolicy, request,
        options: options);
  }

  /// Retrieves information about IAM entity usage and IAM quotas in the Amazon Web Services account. For information about IAM quotas, see IAM and STS quotas in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetAccountSummaryResponse> getAccountSummary(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getAccountSummary, request, options: options);
  }

  /// Gets a list of all of the context keys referenced in the input policies. The policies are supplied as a list of one or more strings. To get the context keys from policies associated with an IAM use...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetContextKeysForPolicyResponse>
      getContextKeysForCustomPolicy(
    $0.GetContextKeysForCustomPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContextKeysForCustomPolicy, request,
        options: options);
  }

  /// Gets a list of all of the context keys referenced in all the IAM policies that are attached to the specified IAM entity. The entity can be an IAM user, group, or role. If you specify a user, then t...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetContextKeysForPolicyResponse>
      getContextKeysForPrincipalPolicy(
    $0.GetContextKeysForPrincipalPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getContextKeysForPrincipalPolicy, request,
        options: options);
  }

  /// Retrieves a credential report for the Amazon Web Services account. For more information about the credential report, see Getting credential reports in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetCredentialReportResponse> getCredentialReport(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getCredentialReport, request, options: options);
  }

  /// Retrieves information about a specific delegation request. If a delegation request has no owner or owner account, GetDelegationRequest for that delegation request can be called by any account. If t...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetDelegationRequestResponse> getDelegationRequest(
    $0.GetDelegationRequestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getDelegationRequest, request, options: options);
  }

  /// Returns a list of IAM users that are in the specified IAM group. You can paginate the results using the MaxItems and Marker parameters.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetGroupResponse> getGroup(
    $0.GetGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGroup, request, options: options);
  }

  /// Retrieves the specified inline policy document that is embedded in the specified IAM group. Policies returned by this operation are URL-encoded compliant with RFC 3986. You can use a URL decoding m...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetGroupPolicyResponse> getGroupPolicy(
    $0.GetGroupPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getGroupPolicy, request, options: options);
  }

  /// Retrieves a human readable summary for a given entity. At this time, the only supported entity type is delegation-request This method uses a Large Language Model (LLM) to generate the summary. If a...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetHumanReadableSummaryResponse>
      getHumanReadableSummary(
    $0.GetHumanReadableSummaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getHumanReadableSummary, request,
        options: options);
  }

  /// Retrieves information about the specified instance profile, including the instance profile's path, GUID, ARN, and role. For more information about instance profiles, see Using instance profiles in ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetInstanceProfileResponse> getInstanceProfile(
    $0.GetInstanceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getInstanceProfile, request, options: options);
  }

  /// Retrieves the user name for the specified IAM user. A login profile is created when you create a password for the user to access the Amazon Web Services Management Console. If the user does not exi...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetLoginProfileResponse> getLoginProfile(
    $0.GetLoginProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getLoginProfile, request, options: options);
  }

  /// Retrieves information about an MFA device for a specified user.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetMFADeviceResponse> getMFADevice(
    $0.GetMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getMFADevice, request, options: options);
  }

  /// Returns information about the specified OpenID Connect (OIDC) provider resource object in IAM.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetOpenIDConnectProviderResponse>
      getOpenIDConnectProvider(
    $0.GetOpenIDConnectProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOpenIDConnectProvider, request,
        options: options);
  }

  /// Retrieves the service last accessed data report for Organizations that was previously generated using the GenerateOrganizationsAccessReport operation. This operation retrieves the status of your re...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetOrganizationsAccessReportResponse>
      getOrganizationsAccessReport(
    $0.GetOrganizationsAccessReportRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOrganizationsAccessReport, request,
        options: options);
  }

  /// Retrieves the configuration information for the outbound identity federation feature in your Amazon Web Services account. The response includes the unique issuer URL for your Amazon Web Services ac...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetOutboundWebIdentityFederationInfoResponse>
      getOutboundWebIdentityFederationInfo(
    $1.Empty request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getOutboundWebIdentityFederationInfo, request,
        options: options);
  }

  /// Retrieves information about the specified managed policy, including the policy's default version and the total number of IAM users, groups, and roles to which the policy is attached. To retrieve th...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetPolicyResponse> getPolicy(
    $0.GetPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPolicy, request, options: options);
  }

  /// Retrieves information about the specified version of the specified managed policy, including the policy document. Policies returned by this operation are URL-encoded compliant with RFC 3986. You ca...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetPolicyVersionResponse> getPolicyVersion(
    $0.GetPolicyVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getPolicyVersion, request, options: options);
  }

  /// Retrieves information about the specified role, including the role's path, GUID, ARN, and the role's trust policy that grants permission to assume the role. For more information about roles, see IA...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetRoleResponse> getRole(
    $0.GetRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRole, request, options: options);
  }

  /// Retrieves the specified inline policy document that is embedded with the specified IAM role. Policies returned by this operation are URL-encoded compliant with RFC 3986. You can use a URL decoding ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetRolePolicyResponse> getRolePolicy(
    $0.GetRolePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getRolePolicy, request, options: options);
  }

  /// Returns the SAML provider metadocument that was uploaded when the IAM SAML provider resource object was created or updated. This operation requires Signature Version 4.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetSAMLProviderResponse> getSAMLProvider(
    $0.GetSAMLProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSAMLProvider, request, options: options);
  }

  /// Retrieves information about the specified server certificate stored in IAM. For more information about working with server certificates, see Working with server certificates in the IAM User Guide. ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetServerCertificateResponse> getServerCertificate(
    $0.GetServerCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getServerCertificate, request, options: options);
  }

  /// Retrieves a service last accessed report that was created using the GenerateServiceLastAccessedDetails operation. You can use the JobId parameter in GetServiceLastAccessedDetails to retrieve the st...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetServiceLastAccessedDetailsResponse>
      getServiceLastAccessedDetails(
    $0.GetServiceLastAccessedDetailsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getServiceLastAccessedDetails, request,
        options: options);
  }

  /// After you generate a group or policy report using the GenerateServiceLastAccessedDetails operation, you can use the JobId parameter in GetServiceLastAccessedDetailsWithEntities. This operation retr...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetServiceLastAccessedDetailsWithEntitiesResponse>
      getServiceLastAccessedDetailsWithEntities(
    $0.GetServiceLastAccessedDetailsWithEntitiesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(
        _$getServiceLastAccessedDetailsWithEntities, request,
        options: options);
  }

  /// Retrieves the status of your service-linked role deletion. After you use DeleteServiceLinkedRole to submit a service-linked role for deletion, you can use the DeletionTaskId parameter in GetService...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetServiceLinkedRoleDeletionStatusResponse>
      getServiceLinkedRoleDeletionStatus(
    $0.GetServiceLinkedRoleDeletionStatusRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getServiceLinkedRoleDeletionStatus, request,
        options: options);
  }

  /// Retrieves the specified SSH public key, including metadata about the key. The SSH public key retrieved by this operation is used only for authenticating the associated IAM user to an CodeCommit rep...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetSSHPublicKeyResponse> getSSHPublicKey(
    $0.GetSSHPublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getSSHPublicKey, request, options: options);
  }

  /// Retrieves information about the specified IAM user, including the user's creation date, path, unique ID, and ARN. If you do not specify a user name, IAM determines the user name implicitly based on...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetUserResponse> getUser(
    $0.GetUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUser, request, options: options);
  }

  /// Retrieves the specified inline policy document that is embedded in the specified IAM user. Policies returned by this operation are URL-encoded compliant with RFC 3986. You can use a URL decoding me...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.GetUserPolicyResponse> getUserPolicy(
    $0.GetUserPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$getUserPolicy, request, options: options);
  }

  /// Returns information about the access key IDs associated with the specified IAM user. If there is none, the operation returns an empty list. Although each user is limited to a small number of keys, ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListAccessKeysResponse> listAccessKeys(
    $0.ListAccessKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAccessKeys, request, options: options);
  }

  /// Lists the account alias associated with the Amazon Web Services account (Note: you can have only one). For information about using an Amazon Web Services account alias, see Creating, deleting, and ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListAccountAliasesResponse> listAccountAliases(
    $0.ListAccountAliasesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAccountAliases, request, options: options);
  }

  /// Lists all managed policies that are attached to the specified IAM group. An IAM group can also have inline policies embedded with it. To list the inline policies for a group, use ListGroupPolicies....
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListAttachedGroupPoliciesResponse>
      listAttachedGroupPolicies(
    $0.ListAttachedGroupPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAttachedGroupPolicies, request,
        options: options);
  }

  /// Lists all managed policies that are attached to the specified IAM role. An IAM role can also have inline policies embedded with it. To list the inline policies for a role, use ListRolePolicies. For...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListAttachedRolePoliciesResponse>
      listAttachedRolePolicies(
    $0.ListAttachedRolePoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAttachedRolePolicies, request,
        options: options);
  }

  /// Lists all managed policies that are attached to the specified IAM user. An IAM user can also have inline policies embedded with it. To list the inline policies for a user, use ListUserPolicies. For...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListAttachedUserPoliciesResponse>
      listAttachedUserPolicies(
    $0.ListAttachedUserPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listAttachedUserPolicies, request,
        options: options);
  }

  /// Lists delegation requests based on the specified criteria. If a delegation request has no owner, even if it is assigned to a specific account, it will not be part of the ListDelegationRequests outp...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListDelegationRequestsResponse>
      listDelegationRequests(
    $0.ListDelegationRequestsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listDelegationRequests, request,
        options: options);
  }

  /// Lists all IAM users, groups, and roles that the specified managed policy is attached to. You can use the optional EntityFilter parameter to limit the results to a particular type of entity (users, ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListEntitiesForPolicyResponse> listEntitiesForPolicy(
    $0.ListEntitiesForPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listEntitiesForPolicy, request, options: options);
  }

  /// Lists the names of the inline policies that are embedded in the specified IAM group. An IAM group can also have managed policies attached to it. To list the managed policies that are attached to a ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListGroupPoliciesResponse> listGroupPolicies(
    $0.ListGroupPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGroupPolicies, request, options: options);
  }

  /// Lists the IAM groups that have the specified path prefix. You can paginate the results using the MaxItems and Marker parameters.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListGroupsResponse> listGroups(
    $0.ListGroupsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGroups, request, options: options);
  }

  /// Lists the IAM groups that the specified IAM user belongs to. You can paginate the results using the MaxItems and Marker parameters.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListGroupsForUserResponse> listGroupsForUser(
    $0.ListGroupsForUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listGroupsForUser, request, options: options);
  }

  /// Lists the instance profiles that have the specified path prefix. If there are none, the operation returns an empty list. For more information about instance profiles, see Using instance profiles in...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListInstanceProfilesResponse> listInstanceProfiles(
    $0.ListInstanceProfilesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInstanceProfiles, request, options: options);
  }

  /// Lists the instance profiles that have the specified associated IAM role. If there are none, the operation returns an empty list. For more information about instance profiles, go to Using instance p...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListInstanceProfilesForRoleResponse>
      listInstanceProfilesForRole(
    $0.ListInstanceProfilesForRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInstanceProfilesForRole, request,
        options: options);
  }

  /// Lists the tags that are attached to the specified IAM instance profile. The returned list of tags is sorted by tag key. For more information about tagging, see Tagging IAM resources in the IAM User...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListInstanceProfileTagsResponse>
      listInstanceProfileTags(
    $0.ListInstanceProfileTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listInstanceProfileTags, request,
        options: options);
  }

  /// Lists the MFA devices for an IAM user. If the request includes a IAM user name, then this operation lists all the MFA devices associated with the specified user. If you do not specify a user name, ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListMFADevicesResponse> listMFADevices(
    $0.ListMFADevicesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMFADevices, request, options: options);
  }

  /// Lists the tags that are attached to the specified IAM virtual multi-factor authentication (MFA) device. The returned list of tags is sorted by tag key. For more information about tagging, see Taggi...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListMFADeviceTagsResponse> listMFADeviceTags(
    $0.ListMFADeviceTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listMFADeviceTags, request, options: options);
  }

  /// Lists information about the IAM OpenID Connect (OIDC) provider resource objects defined in the Amazon Web Services account. IAM resource-listing operations return a subset of the available attribut...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListOpenIDConnectProvidersResponse>
      listOpenIDConnectProviders(
    $0.ListOpenIDConnectProvidersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOpenIDConnectProviders, request,
        options: options);
  }

  /// Lists the tags that are attached to the specified OpenID Connect (OIDC)-compatible identity provider. The returned list of tags is sorted by tag key. For more information, see About web identity fe...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListOpenIDConnectProviderTagsResponse>
      listOpenIDConnectProviderTags(
    $0.ListOpenIDConnectProviderTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOpenIDConnectProviderTags, request,
        options: options);
  }

  /// Lists the centralized root access features enabled for your organization. For more information, see Centrally manage root access for member accounts.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListOrganizationsFeaturesResponse>
      listOrganizationsFeatures(
    $0.ListOrganizationsFeaturesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listOrganizationsFeatures, request,
        options: options);
  }

  /// Lists all the managed policies that are available in your Amazon Web Services account, including your own customer-defined managed policies and all Amazon Web Services managed policies. You can fil...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListPoliciesResponse> listPolicies(
    $0.ListPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPolicies, request, options: options);
  }

  /// Retrieves a list of policies that the IAM identity (user, group, or role) can use to access each specified service. This operation does not use other policy types when determining whether a resourc...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListPoliciesGrantingServiceAccessResponse>
      listPoliciesGrantingServiceAccess(
    $0.ListPoliciesGrantingServiceAccessRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPoliciesGrantingServiceAccess, request,
        options: options);
  }

  /// Lists the tags that are attached to the specified IAM customer managed policy. The returned list of tags is sorted by tag key. For more information about tagging, see Tagging IAM resources in the I...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListPolicyTagsResponse> listPolicyTags(
    $0.ListPolicyTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPolicyTags, request, options: options);
  }

  /// Lists information about the versions of the specified managed policy, including the version that is currently set as the policy's default version. For more information about managed policies, see M...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListPolicyVersionsResponse> listPolicyVersions(
    $0.ListPolicyVersionsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listPolicyVersions, request, options: options);
  }

  /// Lists the names of the inline policies that are embedded in the specified IAM role. An IAM role can also have managed policies attached to it. To list the managed policies that are attached to a ro...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListRolePoliciesResponse> listRolePolicies(
    $0.ListRolePoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRolePolicies, request, options: options);
  }

  /// Lists the IAM roles that have the specified path prefix. If there are none, the operation returns an empty list. For more information about roles, see IAM roles in the IAM User Guide. IAM resource-...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListRolesResponse> listRoles(
    $0.ListRolesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRoles, request, options: options);
  }

  /// Lists the tags that are attached to the specified role. The returned list of tags is sorted by tag key. For more information about tagging, see Tagging IAM resources in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListRoleTagsResponse> listRoleTags(
    $0.ListRoleTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listRoleTags, request, options: options);
  }

  /// Lists the SAML provider resource objects defined in IAM in the account. IAM resource-listing operations return a subset of the available attributes for the resource. For example, this operation doe...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListSAMLProvidersResponse> listSAMLProviders(
    $0.ListSAMLProvidersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSAMLProviders, request, options: options);
  }

  /// Lists the tags that are attached to the specified Security Assertion Markup Language (SAML) identity provider. The returned list of tags is sorted by tag key. For more information, see About SAML 2...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListSAMLProviderTagsResponse> listSAMLProviderTags(
    $0.ListSAMLProviderTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSAMLProviderTags, request, options: options);
  }

  /// Lists the server certificates stored in IAM that have the specified path prefix. If none exist, the operation returns an empty list. You can paginate the results using the MaxItems and Marker param...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListServerCertificatesResponse>
      listServerCertificates(
    $0.ListServerCertificatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listServerCertificates, request,
        options: options);
  }

  /// Lists the tags that are attached to the specified IAM server certificate. The returned list of tags is sorted by tag key. For more information about tagging, see Tagging IAM resources in the IAM Us...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListServerCertificateTagsResponse>
      listServerCertificateTags(
    $0.ListServerCertificateTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listServerCertificateTags, request,
        options: options);
  }

  /// Returns information about the service-specific credentials associated with the specified IAM user. If none exists, the operation returns an empty list. The service-specific credentials returned by ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListServiceSpecificCredentialsResponse>
      listServiceSpecificCredentials(
    $0.ListServiceSpecificCredentialsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listServiceSpecificCredentials, request,
        options: options);
  }

  /// Returns information about the signing certificates associated with the specified IAM user. If none exists, the operation returns an empty list. Although each user is limited to a small number of si...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListSigningCertificatesResponse>
      listSigningCertificates(
    $0.ListSigningCertificatesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSigningCertificates, request,
        options: options);
  }

  /// Returns information about the SSH public keys associated with the specified IAM user. If none exists, the operation returns an empty list. The SSH public keys returned by this operation are used on...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListSSHPublicKeysResponse> listSSHPublicKeys(
    $0.ListSSHPublicKeysRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listSSHPublicKeys, request, options: options);
  }

  /// Lists the names of the inline policies embedded in the specified IAM user. An IAM user can also have managed policies attached to it. To list the managed policies that are attached to a user, use L...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListUserPoliciesResponse> listUserPolicies(
    $0.ListUserPoliciesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUserPolicies, request, options: options);
  }

  /// Lists the IAM users that have the specified path prefix. If no path prefix is specified, the operation returns all users in the Amazon Web Services account. If there are none, the operation returns...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListUsersResponse> listUsers(
    $0.ListUsersRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUsers, request, options: options);
  }

  /// Lists the tags that are attached to the specified IAM user. The returned list of tags is sorted by tag key. For more information about tagging, see Tagging IAM resources in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListUserTagsResponse> listUserTags(
    $0.ListUserTagsRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listUserTags, request, options: options);
  }

  /// Lists the virtual MFA devices defined in the Amazon Web Services account by assignment status. If you do not specify an assignment status, the operation returns a list of all virtual MFA devices. A...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ListVirtualMFADevicesResponse> listVirtualMFADevices(
    $0.ListVirtualMFADevicesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$listVirtualMFADevices, request, options: options);
  }

  /// Adds or updates an inline policy document that is embedded in the specified IAM group. A user can also have managed policies attached to it. To attach a managed policy to a group, use AttachGroupPo...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> putGroupPolicy(
    $0.PutGroupPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putGroupPolicy, request, options: options);
  }

  /// Adds or updates the policy that is specified as the IAM role's permissions boundary. You can use an Amazon Web Services managed policy or a customer managed policy to set the boundary for a role. U...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> putRolePermissionsBoundary(
    $0.PutRolePermissionsBoundaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putRolePermissionsBoundary, request,
        options: options);
  }

  /// Adds or updates an inline policy document that is embedded in the specified IAM role. When you embed an inline policy in a role, the inline policy is used as part of the role's access (permissions)...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> putRolePolicy(
    $0.PutRolePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putRolePolicy, request, options: options);
  }

  /// Adds or updates the policy that is specified as the IAM user's permissions boundary. You can use an Amazon Web Services managed policy or a customer managed policy to set the boundary for a user. U...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> putUserPermissionsBoundary(
    $0.PutUserPermissionsBoundaryRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putUserPermissionsBoundary, request,
        options: options);
  }

  /// Adds or updates an inline policy document that is embedded in the specified IAM user. An IAM user can also have a managed policy attached to it. To attach a managed policy to a user, use AttachUser...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> putUserPolicy(
    $0.PutUserPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$putUserPolicy, request, options: options);
  }

  /// Rejects a delegation request, denying the requested temporary access. Once a request is rejected, it cannot be accepted or updated later. Rejected requests expire after 7 days. When rejecting a req...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> rejectDelegationRequest(
    $0.RejectDelegationRequestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$rejectDelegationRequest, request,
        options: options);
  }

  /// Removes the specified client ID (also known as audience) from the list of client IDs registered for the specified IAM OpenID Connect (OIDC) provider resource object. This operation is idempotent; i...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> removeClientIDFromOpenIDConnectProvider(
    $0.RemoveClientIDFromOpenIDConnectProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeClientIDFromOpenIDConnectProvider, request,
        options: options);
  }

  /// Removes the specified IAM role from the specified Amazon EC2 instance profile. Make sure that you do not have any Amazon EC2 instances running with the role you are about to remove from the instanc...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> removeRoleFromInstanceProfile(
    $0.RemoveRoleFromInstanceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeRoleFromInstanceProfile, request,
        options: options);
  }

  /// Removes the specified user from the specified group.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> removeUserFromGroup(
    $0.RemoveUserFromGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$removeUserFromGroup, request, options: options);
  }

  /// Resets the password for a service-specific credential. The new password is Amazon Web Services generated and cryptographically strong. It cannot be configured by the user. Resetting the password im...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.ResetServiceSpecificCredentialResponse>
      resetServiceSpecificCredential(
    $0.ResetServiceSpecificCredentialRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resetServiceSpecificCredential, request,
        options: options);
  }

  /// Synchronizes the specified MFA device with its IAM resource object on the Amazon Web Services servers. For more information about creating and working with virtual MFA devices, see Using a virtual ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> resyncMFADevice(
    $0.ResyncMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$resyncMFADevice, request, options: options);
  }

  /// Sends the exchange token for an accepted delegation request. The exchange token is sent to the partner via an asynchronous notification channel, established by the partner. The delegation request m...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> sendDelegationToken(
    $0.SendDelegationTokenRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$sendDelegationToken, request, options: options);
  }

  /// Sets the specified version of the specified policy as the policy's default (operative) version. This operation affects all users, groups, and roles that the policy is attached to. To list the users...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> setDefaultPolicyVersion(
    $0.SetDefaultPolicyVersionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setDefaultPolicyVersion, request,
        options: options);
  }

  /// Sets the specified version of the global endpoint token as the token version used for the Amazon Web Services account. By default, Security Token Service (STS) is available as a global service, and...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> setSecurityTokenServicePreferences(
    $0.SetSecurityTokenServicePreferencesRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$setSecurityTokenServicePreferences, request,
        options: options);
  }

  /// Simulate how a set of IAM policies and optionally a resource-based policy works with a list of API operations and Amazon Web Services resources to determine the policies' effective permissions. The...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.SimulatePolicyResponse> simulateCustomPolicy(
    $0.SimulateCustomPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$simulateCustomPolicy, request, options: options);
  }

  /// Simulate how a set of IAM policies attached to an IAM entity works with a list of API operations and Amazon Web Services resources to determine the policies' effective permissions. The entity can b...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.SimulatePolicyResponse> simulatePrincipalPolicy(
    $0.SimulatePrincipalPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$simulatePrincipalPolicy, request,
        options: options);
  }

  /// Adds one or more tags to an IAM instance profile. If a tag with the same key name already exists, then that tag is overwritten with the new value. Each tag consists of a key name and an associated ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagInstanceProfile(
    $0.TagInstanceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagInstanceProfile, request, options: options);
  }

  /// Adds one or more tags to an IAM virtual multi-factor authentication (MFA) device. If a tag with the same key name already exists, then that tag is overwritten with the new value. A tag consists of ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagMFADevice(
    $0.TagMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagMFADevice, request, options: options);
  }

  /// Adds one or more tags to an OpenID Connect (OIDC)-compatible identity provider. For more information about these providers, see About web identity federation. If a tag with the same key name alread...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagOpenIDConnectProvider(
    $0.TagOpenIDConnectProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagOpenIDConnectProvider, request,
        options: options);
  }

  /// Adds one or more tags to an IAM customer managed policy. If a tag with the same key name already exists, then that tag is overwritten with the new value. A tag consists of a key name and an associa...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagPolicy(
    $0.TagPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagPolicy, request, options: options);
  }

  /// Adds one or more tags to an IAM role. The role can be a regular role or a service-linked role. If a tag with the same key name already exists, then that tag is overwritten with the new value. A tag...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagRole(
    $0.TagRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagRole, request, options: options);
  }

  /// Adds one or more tags to a Security Assertion Markup Language (SAML) identity provider. For more information about these providers, see About SAML 2.0-based federation . If a tag with the same key ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagSAMLProvider(
    $0.TagSAMLProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagSAMLProvider, request, options: options);
  }

  /// Adds one or more tags to an IAM server certificate. If a tag with the same key name already exists, then that tag is overwritten with the new value. For certificates in a Region supported by Certif...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagServerCertificate(
    $0.TagServerCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagServerCertificate, request, options: options);
  }

  /// Adds one or more tags to an IAM user. If a tag with the same key name already exists, then that tag is overwritten with the new value. A tag consists of a key name and an associated value. By assig...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> tagUser(
    $0.TagUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$tagUser, request, options: options);
  }

  /// Removes the specified tags from the IAM instance profile. For more information about tagging, see Tagging IAM resources in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagInstanceProfile(
    $0.UntagInstanceProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagInstanceProfile, request, options: options);
  }

  /// Removes the specified tags from the IAM virtual multi-factor authentication (MFA) device. For more information about tagging, see Tagging IAM resources in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagMFADevice(
    $0.UntagMFADeviceRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagMFADevice, request, options: options);
  }

  /// Removes the specified tags from the specified OpenID Connect (OIDC)-compatible identity provider in IAM. For more information about OIDC providers, see About web identity federation. For more infor...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagOpenIDConnectProvider(
    $0.UntagOpenIDConnectProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagOpenIDConnectProvider, request,
        options: options);
  }

  /// Removes the specified tags from the customer managed policy. For more information about tagging, see Tagging IAM resources in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagPolicy(
    $0.UntagPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagPolicy, request, options: options);
  }

  /// Removes the specified tags from the role. For more information about tagging, see Tagging IAM resources in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagRole(
    $0.UntagRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagRole, request, options: options);
  }

  /// Removes the specified tags from the specified Security Assertion Markup Language (SAML) identity provider in IAM. For more information about these providers, see About web identity federation. For ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagSAMLProvider(
    $0.UntagSAMLProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagSAMLProvider, request, options: options);
  }

  /// Removes the specified tags from the IAM server certificate. For more information about tagging, see Tagging IAM resources in the IAM User Guide. For certificates in a Region supported by Certificat...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagServerCertificate(
    $0.UntagServerCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagServerCertificate, request,
        options: options);
  }

  /// Removes the specified tags from the user. For more information about tagging, see Tagging IAM resources in the IAM User Guide.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> untagUser(
    $0.UntagUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$untagUser, request, options: options);
  }

  /// Changes the status of the specified access key from Active to Inactive, or vice versa. This operation can be used to disable a user's key as part of a key rotation workflow. If the UserName is not ...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateAccessKey(
    $0.UpdateAccessKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAccessKey, request, options: options);
  }

  /// Updates the password policy settings for the Amazon Web Services account. This operation does not support partial updates. No parameters are required, but if you do not specify a parameter, that pa...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateAccountPasswordPolicy(
    $0.UpdateAccountPasswordPolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAccountPasswordPolicy, request,
        options: options);
  }

  /// Updates the policy that grants an IAM entity permission to assume a role. This is typically referred to as the "role trust policy". For more information about roles, see Using roles to delegate per...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateAssumeRolePolicy(
    $0.UpdateAssumeRolePolicyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateAssumeRolePolicy, request,
        options: options);
  }

  /// Updates an existing delegation request with additional information. When the delegation request is updated, it reaches the PENDING_APPROVAL state. Once a delegation request has an owner, that owner...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateDelegationRequest(
    $0.UpdateDelegationRequestRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateDelegationRequest, request,
        options: options);
  }

  /// Updates the name and/or the path of the specified IAM group. You should understand the implications of changing a group's path or name. For more information, see Renaming users and groups in the IA...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateGroup(
    $0.UpdateGroupRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateGroup, request, options: options);
  }

  /// Changes the password for the specified IAM user. You can use the CLI, the Amazon Web Services API, or the Users page in the IAM console to change the password for any IAM user. Use ChangePassword t...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateLoginProfile(
    $0.UpdateLoginProfileRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateLoginProfile, request, options: options);
  }

  /// Replaces the existing list of server certificate thumbprints associated with an OpenID Connect (OIDC) provider resource object with a new list of thumbprints. The list that you pass with this opera...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateOpenIDConnectProviderThumbprint(
    $0.UpdateOpenIDConnectProviderThumbprintRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateOpenIDConnectProviderThumbprint, request,
        options: options);
  }

  /// Updates the description or maximum session duration setting of a role.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.UpdateRoleResponse> updateRole(
    $0.UpdateRoleRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRole, request, options: options);
  }

  /// Use UpdateRole instead. Modifies only the description of a role. This operation performs the same function as the Description parameter in the UpdateRole operation.
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.UpdateRoleDescriptionResponse> updateRoleDescription(
    $0.UpdateRoleDescriptionRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateRoleDescription, request, options: options);
  }

  /// Updates the metadata document, SAML encryption settings, and private keys for an existing SAML provider. To rotate private keys, add your new private key and then remove the old key in a separate r...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.UpdateSAMLProviderResponse> updateSAMLProvider(
    $0.UpdateSAMLProviderRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSAMLProvider, request, options: options);
  }

  /// Updates the name and/or the path of the specified server certificate stored in IAM. For more information about working with server certificates, see Working with server certificates in the IAM User...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateServerCertificate(
    $0.UpdateServerCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateServerCertificate, request,
        options: options);
  }

  /// Sets the status of a service-specific credential to Active or Inactive. Service-specific credentials that are inactive cannot be used for authentication to the service. This operation can be used t...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateServiceSpecificCredential(
    $0.UpdateServiceSpecificCredentialRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateServiceSpecificCredential, request,
        options: options);
  }

  /// Changes the status of the specified user signing certificate from active to disabled, or vice versa. This operation can be used to disable an IAM user's signing certificate as part of a certificate...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateSigningCertificate(
    $0.UpdateSigningCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSigningCertificate, request,
        options: options);
  }

  /// Sets the status of an IAM user's SSH public key to active or inactive. SSH public keys that are inactive cannot be used for authentication. This operation can be used to disable a user's SSH public...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateSSHPublicKey(
    $0.UpdateSSHPublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateSSHPublicKey, request, options: options);
  }

  /// Updates the name and/or the path of the specified IAM user. You should understand the implications of changing an IAM user's path or name. For more information, see Renaming an IAM user and Renamin...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$1.Empty> updateUser(
    $0.UpdateUserRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$updateUser, request, options: options);
  }

  /// Uploads a server certificate entity for the Amazon Web Services account. The server certificate entity includes a public key certificate, a private key, and an optional certificate chain, which sho...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.UploadServerCertificateResponse>
      uploadServerCertificate(
    $0.UploadServerCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$uploadServerCertificate, request,
        options: options);
  }

  /// Uploads an X.509 signing certificate and associates it with the specified IAM user. Some Amazon Web Services services require you to use certificates to validate requests that are signed with a cor...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.UploadSigningCertificateResponse>
      uploadSigningCertificate(
    $0.UploadSigningCertificateRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$uploadSigningCertificate, request,
        options: options);
  }

  /// Uploads an SSH public key and associates it with the specified IAM user. The SSH public key uploaded by this operation can be used only for authenticating the associated IAM user to an CodeCommit r...
  /// HTTP: POST /
  /// Protocol: awsQuery
  $grpc.ResponseFuture<$0.UploadSSHPublicKeyResponse> uploadSSHPublicKey(
    $0.UploadSSHPublicKeyRequest request, {
    $grpc.CallOptions? options,
  }) {
    return $createUnaryCall(_$uploadSSHPublicKey, request, options: options);
  }

  // method descriptors

  static final _$acceptDelegationRequest =
      $grpc.ClientMethod<$0.AcceptDelegationRequestRequest, $1.Empty>(
          '/iam.IAMService/AcceptDelegationRequest',
          ($0.AcceptDelegationRequestRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$addClientIDToOpenIDConnectProvider = $grpc.ClientMethod<
          $0.AddClientIDToOpenIDConnectProviderRequest, $1.Empty>(
      '/iam.IAMService/AddClientIDToOpenIDConnectProvider',
      ($0.AddClientIDToOpenIDConnectProviderRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$addRoleToInstanceProfile =
      $grpc.ClientMethod<$0.AddRoleToInstanceProfileRequest, $1.Empty>(
          '/iam.IAMService/AddRoleToInstanceProfile',
          ($0.AddRoleToInstanceProfileRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$addUserToGroup =
      $grpc.ClientMethod<$0.AddUserToGroupRequest, $1.Empty>(
          '/iam.IAMService/AddUserToGroup',
          ($0.AddUserToGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$associateDelegationRequest =
      $grpc.ClientMethod<$0.AssociateDelegationRequestRequest, $1.Empty>(
          '/iam.IAMService/AssociateDelegationRequest',
          ($0.AssociateDelegationRequestRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$attachGroupPolicy =
      $grpc.ClientMethod<$0.AttachGroupPolicyRequest, $1.Empty>(
          '/iam.IAMService/AttachGroupPolicy',
          ($0.AttachGroupPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$attachRolePolicy =
      $grpc.ClientMethod<$0.AttachRolePolicyRequest, $1.Empty>(
          '/iam.IAMService/AttachRolePolicy',
          ($0.AttachRolePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$attachUserPolicy =
      $grpc.ClientMethod<$0.AttachUserPolicyRequest, $1.Empty>(
          '/iam.IAMService/AttachUserPolicy',
          ($0.AttachUserPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$changePassword =
      $grpc.ClientMethod<$0.ChangePasswordRequest, $1.Empty>(
          '/iam.IAMService/ChangePassword',
          ($0.ChangePasswordRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$createAccessKey =
      $grpc.ClientMethod<$0.CreateAccessKeyRequest, $0.CreateAccessKeyResponse>(
          '/iam.IAMService/CreateAccessKey',
          ($0.CreateAccessKeyRequest value) => value.writeToBuffer(),
          $0.CreateAccessKeyResponse.fromBuffer);
  static final _$createAccountAlias =
      $grpc.ClientMethod<$0.CreateAccountAliasRequest, $1.Empty>(
          '/iam.IAMService/CreateAccountAlias',
          ($0.CreateAccountAliasRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$createDelegationRequest = $grpc.ClientMethod<
          $0.CreateDelegationRequestRequest,
          $0.CreateDelegationRequestResponse>(
      '/iam.IAMService/CreateDelegationRequest',
      ($0.CreateDelegationRequestRequest value) => value.writeToBuffer(),
      $0.CreateDelegationRequestResponse.fromBuffer);
  static final _$createGroup =
      $grpc.ClientMethod<$0.CreateGroupRequest, $0.CreateGroupResponse>(
          '/iam.IAMService/CreateGroup',
          ($0.CreateGroupRequest value) => value.writeToBuffer(),
          $0.CreateGroupResponse.fromBuffer);
  static final _$createInstanceProfile = $grpc.ClientMethod<
          $0.CreateInstanceProfileRequest, $0.CreateInstanceProfileResponse>(
      '/iam.IAMService/CreateInstanceProfile',
      ($0.CreateInstanceProfileRequest value) => value.writeToBuffer(),
      $0.CreateInstanceProfileResponse.fromBuffer);
  static final _$createLoginProfile = $grpc.ClientMethod<
          $0.CreateLoginProfileRequest, $0.CreateLoginProfileResponse>(
      '/iam.IAMService/CreateLoginProfile',
      ($0.CreateLoginProfileRequest value) => value.writeToBuffer(),
      $0.CreateLoginProfileResponse.fromBuffer);
  static final _$createOpenIDConnectProvider = $grpc.ClientMethod<
          $0.CreateOpenIDConnectProviderRequest,
          $0.CreateOpenIDConnectProviderResponse>(
      '/iam.IAMService/CreateOpenIDConnectProvider',
      ($0.CreateOpenIDConnectProviderRequest value) => value.writeToBuffer(),
      $0.CreateOpenIDConnectProviderResponse.fromBuffer);
  static final _$createPolicy =
      $grpc.ClientMethod<$0.CreatePolicyRequest, $0.CreatePolicyResponse>(
          '/iam.IAMService/CreatePolicy',
          ($0.CreatePolicyRequest value) => value.writeToBuffer(),
          $0.CreatePolicyResponse.fromBuffer);
  static final _$createPolicyVersion = $grpc.ClientMethod<
          $0.CreatePolicyVersionRequest, $0.CreatePolicyVersionResponse>(
      '/iam.IAMService/CreatePolicyVersion',
      ($0.CreatePolicyVersionRequest value) => value.writeToBuffer(),
      $0.CreatePolicyVersionResponse.fromBuffer);
  static final _$createRole =
      $grpc.ClientMethod<$0.CreateRoleRequest, $0.CreateRoleResponse>(
          '/iam.IAMService/CreateRole',
          ($0.CreateRoleRequest value) => value.writeToBuffer(),
          $0.CreateRoleResponse.fromBuffer);
  static final _$createSAMLProvider = $grpc.ClientMethod<
          $0.CreateSAMLProviderRequest, $0.CreateSAMLProviderResponse>(
      '/iam.IAMService/CreateSAMLProvider',
      ($0.CreateSAMLProviderRequest value) => value.writeToBuffer(),
      $0.CreateSAMLProviderResponse.fromBuffer);
  static final _$createServiceLinkedRole = $grpc.ClientMethod<
          $0.CreateServiceLinkedRoleRequest,
          $0.CreateServiceLinkedRoleResponse>(
      '/iam.IAMService/CreateServiceLinkedRole',
      ($0.CreateServiceLinkedRoleRequest value) => value.writeToBuffer(),
      $0.CreateServiceLinkedRoleResponse.fromBuffer);
  static final _$createServiceSpecificCredential = $grpc.ClientMethod<
          $0.CreateServiceSpecificCredentialRequest,
          $0.CreateServiceSpecificCredentialResponse>(
      '/iam.IAMService/CreateServiceSpecificCredential',
      ($0.CreateServiceSpecificCredentialRequest value) =>
          value.writeToBuffer(),
      $0.CreateServiceSpecificCredentialResponse.fromBuffer);
  static final _$createUser =
      $grpc.ClientMethod<$0.CreateUserRequest, $0.CreateUserResponse>(
          '/iam.IAMService/CreateUser',
          ($0.CreateUserRequest value) => value.writeToBuffer(),
          $0.CreateUserResponse.fromBuffer);
  static final _$createVirtualMFADevice = $grpc.ClientMethod<
          $0.CreateVirtualMFADeviceRequest, $0.CreateVirtualMFADeviceResponse>(
      '/iam.IAMService/CreateVirtualMFADevice',
      ($0.CreateVirtualMFADeviceRequest value) => value.writeToBuffer(),
      $0.CreateVirtualMFADeviceResponse.fromBuffer);
  static final _$deactivateMFADevice =
      $grpc.ClientMethod<$0.DeactivateMFADeviceRequest, $1.Empty>(
          '/iam.IAMService/DeactivateMFADevice',
          ($0.DeactivateMFADeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteAccessKey =
      $grpc.ClientMethod<$0.DeleteAccessKeyRequest, $1.Empty>(
          '/iam.IAMService/DeleteAccessKey',
          ($0.DeleteAccessKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteAccountAlias =
      $grpc.ClientMethod<$0.DeleteAccountAliasRequest, $1.Empty>(
          '/iam.IAMService/DeleteAccountAlias',
          ($0.DeleteAccountAliasRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteAccountPasswordPolicy =
      $grpc.ClientMethod<$1.Empty, $1.Empty>(
          '/iam.IAMService/DeleteAccountPasswordPolicy',
          ($1.Empty value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteGroup =
      $grpc.ClientMethod<$0.DeleteGroupRequest, $1.Empty>(
          '/iam.IAMService/DeleteGroup',
          ($0.DeleteGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteGroupPolicy =
      $grpc.ClientMethod<$0.DeleteGroupPolicyRequest, $1.Empty>(
          '/iam.IAMService/DeleteGroupPolicy',
          ($0.DeleteGroupPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteInstanceProfile =
      $grpc.ClientMethod<$0.DeleteInstanceProfileRequest, $1.Empty>(
          '/iam.IAMService/DeleteInstanceProfile',
          ($0.DeleteInstanceProfileRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteLoginProfile =
      $grpc.ClientMethod<$0.DeleteLoginProfileRequest, $1.Empty>(
          '/iam.IAMService/DeleteLoginProfile',
          ($0.DeleteLoginProfileRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteOpenIDConnectProvider =
      $grpc.ClientMethod<$0.DeleteOpenIDConnectProviderRequest, $1.Empty>(
          '/iam.IAMService/DeleteOpenIDConnectProvider',
          ($0.DeleteOpenIDConnectProviderRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deletePolicy =
      $grpc.ClientMethod<$0.DeletePolicyRequest, $1.Empty>(
          '/iam.IAMService/DeletePolicy',
          ($0.DeletePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deletePolicyVersion =
      $grpc.ClientMethod<$0.DeletePolicyVersionRequest, $1.Empty>(
          '/iam.IAMService/DeletePolicyVersion',
          ($0.DeletePolicyVersionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRole =
      $grpc.ClientMethod<$0.DeleteRoleRequest, $1.Empty>(
          '/iam.IAMService/DeleteRole',
          ($0.DeleteRoleRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRolePermissionsBoundary =
      $grpc.ClientMethod<$0.DeleteRolePermissionsBoundaryRequest, $1.Empty>(
          '/iam.IAMService/DeleteRolePermissionsBoundary',
          ($0.DeleteRolePermissionsBoundaryRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteRolePolicy =
      $grpc.ClientMethod<$0.DeleteRolePolicyRequest, $1.Empty>(
          '/iam.IAMService/DeleteRolePolicy',
          ($0.DeleteRolePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteSAMLProvider =
      $grpc.ClientMethod<$0.DeleteSAMLProviderRequest, $1.Empty>(
          '/iam.IAMService/DeleteSAMLProvider',
          ($0.DeleteSAMLProviderRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteServerCertificate =
      $grpc.ClientMethod<$0.DeleteServerCertificateRequest, $1.Empty>(
          '/iam.IAMService/DeleteServerCertificate',
          ($0.DeleteServerCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteServiceLinkedRole = $grpc.ClientMethod<
          $0.DeleteServiceLinkedRoleRequest,
          $0.DeleteServiceLinkedRoleResponse>(
      '/iam.IAMService/DeleteServiceLinkedRole',
      ($0.DeleteServiceLinkedRoleRequest value) => value.writeToBuffer(),
      $0.DeleteServiceLinkedRoleResponse.fromBuffer);
  static final _$deleteServiceSpecificCredential =
      $grpc.ClientMethod<$0.DeleteServiceSpecificCredentialRequest, $1.Empty>(
          '/iam.IAMService/DeleteServiceSpecificCredential',
          ($0.DeleteServiceSpecificCredentialRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteSigningCertificate =
      $grpc.ClientMethod<$0.DeleteSigningCertificateRequest, $1.Empty>(
          '/iam.IAMService/DeleteSigningCertificate',
          ($0.DeleteSigningCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteSSHPublicKey =
      $grpc.ClientMethod<$0.DeleteSSHPublicKeyRequest, $1.Empty>(
          '/iam.IAMService/DeleteSSHPublicKey',
          ($0.DeleteSSHPublicKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUser =
      $grpc.ClientMethod<$0.DeleteUserRequest, $1.Empty>(
          '/iam.IAMService/DeleteUser',
          ($0.DeleteUserRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUserPermissionsBoundary =
      $grpc.ClientMethod<$0.DeleteUserPermissionsBoundaryRequest, $1.Empty>(
          '/iam.IAMService/DeleteUserPermissionsBoundary',
          ($0.DeleteUserPermissionsBoundaryRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteUserPolicy =
      $grpc.ClientMethod<$0.DeleteUserPolicyRequest, $1.Empty>(
          '/iam.IAMService/DeleteUserPolicy',
          ($0.DeleteUserPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$deleteVirtualMFADevice =
      $grpc.ClientMethod<$0.DeleteVirtualMFADeviceRequest, $1.Empty>(
          '/iam.IAMService/DeleteVirtualMFADevice',
          ($0.DeleteVirtualMFADeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$detachGroupPolicy =
      $grpc.ClientMethod<$0.DetachGroupPolicyRequest, $1.Empty>(
          '/iam.IAMService/DetachGroupPolicy',
          ($0.DetachGroupPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$detachRolePolicy =
      $grpc.ClientMethod<$0.DetachRolePolicyRequest, $1.Empty>(
          '/iam.IAMService/DetachRolePolicy',
          ($0.DetachRolePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$detachUserPolicy =
      $grpc.ClientMethod<$0.DetachUserPolicyRequest, $1.Empty>(
          '/iam.IAMService/DetachUserPolicy',
          ($0.DetachUserPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$disableOrganizationsRootCredentialsManagement =
      $grpc.ClientMethod<
              $0.DisableOrganizationsRootCredentialsManagementRequest,
              $0.DisableOrganizationsRootCredentialsManagementResponse>(
          '/iam.IAMService/DisableOrganizationsRootCredentialsManagement',
          ($0.DisableOrganizationsRootCredentialsManagementRequest value) =>
              value.writeToBuffer(),
          $0.DisableOrganizationsRootCredentialsManagementResponse.fromBuffer);
  static final _$disableOrganizationsRootSessions = $grpc.ClientMethod<
          $0.DisableOrganizationsRootSessionsRequest,
          $0.DisableOrganizationsRootSessionsResponse>(
      '/iam.IAMService/DisableOrganizationsRootSessions',
      ($0.DisableOrganizationsRootSessionsRequest value) =>
          value.writeToBuffer(),
      $0.DisableOrganizationsRootSessionsResponse.fromBuffer);
  static final _$disableOutboundWebIdentityFederation =
      $grpc.ClientMethod<$1.Empty, $1.Empty>(
          '/iam.IAMService/DisableOutboundWebIdentityFederation',
          ($1.Empty value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$enableMFADevice =
      $grpc.ClientMethod<$0.EnableMFADeviceRequest, $1.Empty>(
          '/iam.IAMService/EnableMFADevice',
          ($0.EnableMFADeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$enableOrganizationsRootCredentialsManagement =
      $grpc.ClientMethod<$0.EnableOrganizationsRootCredentialsManagementRequest,
              $0.EnableOrganizationsRootCredentialsManagementResponse>(
          '/iam.IAMService/EnableOrganizationsRootCredentialsManagement',
          ($0.EnableOrganizationsRootCredentialsManagementRequest value) =>
              value.writeToBuffer(),
          $0.EnableOrganizationsRootCredentialsManagementResponse.fromBuffer);
  static final _$enableOrganizationsRootSessions = $grpc.ClientMethod<
          $0.EnableOrganizationsRootSessionsRequest,
          $0.EnableOrganizationsRootSessionsResponse>(
      '/iam.IAMService/EnableOrganizationsRootSessions',
      ($0.EnableOrganizationsRootSessionsRequest value) =>
          value.writeToBuffer(),
      $0.EnableOrganizationsRootSessionsResponse.fromBuffer);
  static final _$enableOutboundWebIdentityFederation = $grpc.ClientMethod<
          $1.Empty, $0.EnableOutboundWebIdentityFederationResponse>(
      '/iam.IAMService/EnableOutboundWebIdentityFederation',
      ($1.Empty value) => value.writeToBuffer(),
      $0.EnableOutboundWebIdentityFederationResponse.fromBuffer);
  static final _$generateCredentialReport =
      $grpc.ClientMethod<$1.Empty, $0.GenerateCredentialReportResponse>(
          '/iam.IAMService/GenerateCredentialReport',
          ($1.Empty value) => value.writeToBuffer(),
          $0.GenerateCredentialReportResponse.fromBuffer);
  static final _$generateOrganizationsAccessReport = $grpc.ClientMethod<
          $0.GenerateOrganizationsAccessReportRequest,
          $0.GenerateOrganizationsAccessReportResponse>(
      '/iam.IAMService/GenerateOrganizationsAccessReport',
      ($0.GenerateOrganizationsAccessReportRequest value) =>
          value.writeToBuffer(),
      $0.GenerateOrganizationsAccessReportResponse.fromBuffer);
  static final _$generateServiceLastAccessedDetails = $grpc.ClientMethod<
          $0.GenerateServiceLastAccessedDetailsRequest,
          $0.GenerateServiceLastAccessedDetailsResponse>(
      '/iam.IAMService/GenerateServiceLastAccessedDetails',
      ($0.GenerateServiceLastAccessedDetailsRequest value) =>
          value.writeToBuffer(),
      $0.GenerateServiceLastAccessedDetailsResponse.fromBuffer);
  static final _$getAccessKeyLastUsed = $grpc.ClientMethod<
          $0.GetAccessKeyLastUsedRequest, $0.GetAccessKeyLastUsedResponse>(
      '/iam.IAMService/GetAccessKeyLastUsed',
      ($0.GetAccessKeyLastUsedRequest value) => value.writeToBuffer(),
      $0.GetAccessKeyLastUsedResponse.fromBuffer);
  static final _$getAccountAuthorizationDetails = $grpc.ClientMethod<
          $0.GetAccountAuthorizationDetailsRequest,
          $0.GetAccountAuthorizationDetailsResponse>(
      '/iam.IAMService/GetAccountAuthorizationDetails',
      ($0.GetAccountAuthorizationDetailsRequest value) => value.writeToBuffer(),
      $0.GetAccountAuthorizationDetailsResponse.fromBuffer);
  static final _$getAccountPasswordPolicy =
      $grpc.ClientMethod<$1.Empty, $0.GetAccountPasswordPolicyResponse>(
          '/iam.IAMService/GetAccountPasswordPolicy',
          ($1.Empty value) => value.writeToBuffer(),
          $0.GetAccountPasswordPolicyResponse.fromBuffer);
  static final _$getAccountSummary =
      $grpc.ClientMethod<$1.Empty, $0.GetAccountSummaryResponse>(
          '/iam.IAMService/GetAccountSummary',
          ($1.Empty value) => value.writeToBuffer(),
          $0.GetAccountSummaryResponse.fromBuffer);
  static final _$getContextKeysForCustomPolicy = $grpc.ClientMethod<
          $0.GetContextKeysForCustomPolicyRequest,
          $0.GetContextKeysForPolicyResponse>(
      '/iam.IAMService/GetContextKeysForCustomPolicy',
      ($0.GetContextKeysForCustomPolicyRequest value) => value.writeToBuffer(),
      $0.GetContextKeysForPolicyResponse.fromBuffer);
  static final _$getContextKeysForPrincipalPolicy = $grpc.ClientMethod<
          $0.GetContextKeysForPrincipalPolicyRequest,
          $0.GetContextKeysForPolicyResponse>(
      '/iam.IAMService/GetContextKeysForPrincipalPolicy',
      ($0.GetContextKeysForPrincipalPolicyRequest value) =>
          value.writeToBuffer(),
      $0.GetContextKeysForPolicyResponse.fromBuffer);
  static final _$getCredentialReport =
      $grpc.ClientMethod<$1.Empty, $0.GetCredentialReportResponse>(
          '/iam.IAMService/GetCredentialReport',
          ($1.Empty value) => value.writeToBuffer(),
          $0.GetCredentialReportResponse.fromBuffer);
  static final _$getDelegationRequest = $grpc.ClientMethod<
          $0.GetDelegationRequestRequest, $0.GetDelegationRequestResponse>(
      '/iam.IAMService/GetDelegationRequest',
      ($0.GetDelegationRequestRequest value) => value.writeToBuffer(),
      $0.GetDelegationRequestResponse.fromBuffer);
  static final _$getGroup =
      $grpc.ClientMethod<$0.GetGroupRequest, $0.GetGroupResponse>(
          '/iam.IAMService/GetGroup',
          ($0.GetGroupRequest value) => value.writeToBuffer(),
          $0.GetGroupResponse.fromBuffer);
  static final _$getGroupPolicy =
      $grpc.ClientMethod<$0.GetGroupPolicyRequest, $0.GetGroupPolicyResponse>(
          '/iam.IAMService/GetGroupPolicy',
          ($0.GetGroupPolicyRequest value) => value.writeToBuffer(),
          $0.GetGroupPolicyResponse.fromBuffer);
  static final _$getHumanReadableSummary = $grpc.ClientMethod<
          $0.GetHumanReadableSummaryRequest,
          $0.GetHumanReadableSummaryResponse>(
      '/iam.IAMService/GetHumanReadableSummary',
      ($0.GetHumanReadableSummaryRequest value) => value.writeToBuffer(),
      $0.GetHumanReadableSummaryResponse.fromBuffer);
  static final _$getInstanceProfile = $grpc.ClientMethod<
          $0.GetInstanceProfileRequest, $0.GetInstanceProfileResponse>(
      '/iam.IAMService/GetInstanceProfile',
      ($0.GetInstanceProfileRequest value) => value.writeToBuffer(),
      $0.GetInstanceProfileResponse.fromBuffer);
  static final _$getLoginProfile =
      $grpc.ClientMethod<$0.GetLoginProfileRequest, $0.GetLoginProfileResponse>(
          '/iam.IAMService/GetLoginProfile',
          ($0.GetLoginProfileRequest value) => value.writeToBuffer(),
          $0.GetLoginProfileResponse.fromBuffer);
  static final _$getMFADevice =
      $grpc.ClientMethod<$0.GetMFADeviceRequest, $0.GetMFADeviceResponse>(
          '/iam.IAMService/GetMFADevice',
          ($0.GetMFADeviceRequest value) => value.writeToBuffer(),
          $0.GetMFADeviceResponse.fromBuffer);
  static final _$getOpenIDConnectProvider = $grpc.ClientMethod<
          $0.GetOpenIDConnectProviderRequest,
          $0.GetOpenIDConnectProviderResponse>(
      '/iam.IAMService/GetOpenIDConnectProvider',
      ($0.GetOpenIDConnectProviderRequest value) => value.writeToBuffer(),
      $0.GetOpenIDConnectProviderResponse.fromBuffer);
  static final _$getOrganizationsAccessReport = $grpc.ClientMethod<
          $0.GetOrganizationsAccessReportRequest,
          $0.GetOrganizationsAccessReportResponse>(
      '/iam.IAMService/GetOrganizationsAccessReport',
      ($0.GetOrganizationsAccessReportRequest value) => value.writeToBuffer(),
      $0.GetOrganizationsAccessReportResponse.fromBuffer);
  static final _$getOutboundWebIdentityFederationInfo = $grpc.ClientMethod<
          $1.Empty, $0.GetOutboundWebIdentityFederationInfoResponse>(
      '/iam.IAMService/GetOutboundWebIdentityFederationInfo',
      ($1.Empty value) => value.writeToBuffer(),
      $0.GetOutboundWebIdentityFederationInfoResponse.fromBuffer);
  static final _$getPolicy =
      $grpc.ClientMethod<$0.GetPolicyRequest, $0.GetPolicyResponse>(
          '/iam.IAMService/GetPolicy',
          ($0.GetPolicyRequest value) => value.writeToBuffer(),
          $0.GetPolicyResponse.fromBuffer);
  static final _$getPolicyVersion = $grpc.ClientMethod<
          $0.GetPolicyVersionRequest, $0.GetPolicyVersionResponse>(
      '/iam.IAMService/GetPolicyVersion',
      ($0.GetPolicyVersionRequest value) => value.writeToBuffer(),
      $0.GetPolicyVersionResponse.fromBuffer);
  static final _$getRole =
      $grpc.ClientMethod<$0.GetRoleRequest, $0.GetRoleResponse>(
          '/iam.IAMService/GetRole',
          ($0.GetRoleRequest value) => value.writeToBuffer(),
          $0.GetRoleResponse.fromBuffer);
  static final _$getRolePolicy =
      $grpc.ClientMethod<$0.GetRolePolicyRequest, $0.GetRolePolicyResponse>(
          '/iam.IAMService/GetRolePolicy',
          ($0.GetRolePolicyRequest value) => value.writeToBuffer(),
          $0.GetRolePolicyResponse.fromBuffer);
  static final _$getSAMLProvider =
      $grpc.ClientMethod<$0.GetSAMLProviderRequest, $0.GetSAMLProviderResponse>(
          '/iam.IAMService/GetSAMLProvider',
          ($0.GetSAMLProviderRequest value) => value.writeToBuffer(),
          $0.GetSAMLProviderResponse.fromBuffer);
  static final _$getServerCertificate = $grpc.ClientMethod<
          $0.GetServerCertificateRequest, $0.GetServerCertificateResponse>(
      '/iam.IAMService/GetServerCertificate',
      ($0.GetServerCertificateRequest value) => value.writeToBuffer(),
      $0.GetServerCertificateResponse.fromBuffer);
  static final _$getServiceLastAccessedDetails = $grpc.ClientMethod<
          $0.GetServiceLastAccessedDetailsRequest,
          $0.GetServiceLastAccessedDetailsResponse>(
      '/iam.IAMService/GetServiceLastAccessedDetails',
      ($0.GetServiceLastAccessedDetailsRequest value) => value.writeToBuffer(),
      $0.GetServiceLastAccessedDetailsResponse.fromBuffer);
  static final _$getServiceLastAccessedDetailsWithEntities = $grpc.ClientMethod<
          $0.GetServiceLastAccessedDetailsWithEntitiesRequest,
          $0.GetServiceLastAccessedDetailsWithEntitiesResponse>(
      '/iam.IAMService/GetServiceLastAccessedDetailsWithEntities',
      ($0.GetServiceLastAccessedDetailsWithEntitiesRequest value) =>
          value.writeToBuffer(),
      $0.GetServiceLastAccessedDetailsWithEntitiesResponse.fromBuffer);
  static final _$getServiceLinkedRoleDeletionStatus = $grpc.ClientMethod<
          $0.GetServiceLinkedRoleDeletionStatusRequest,
          $0.GetServiceLinkedRoleDeletionStatusResponse>(
      '/iam.IAMService/GetServiceLinkedRoleDeletionStatus',
      ($0.GetServiceLinkedRoleDeletionStatusRequest value) =>
          value.writeToBuffer(),
      $0.GetServiceLinkedRoleDeletionStatusResponse.fromBuffer);
  static final _$getSSHPublicKey =
      $grpc.ClientMethod<$0.GetSSHPublicKeyRequest, $0.GetSSHPublicKeyResponse>(
          '/iam.IAMService/GetSSHPublicKey',
          ($0.GetSSHPublicKeyRequest value) => value.writeToBuffer(),
          $0.GetSSHPublicKeyResponse.fromBuffer);
  static final _$getUser =
      $grpc.ClientMethod<$0.GetUserRequest, $0.GetUserResponse>(
          '/iam.IAMService/GetUser',
          ($0.GetUserRequest value) => value.writeToBuffer(),
          $0.GetUserResponse.fromBuffer);
  static final _$getUserPolicy =
      $grpc.ClientMethod<$0.GetUserPolicyRequest, $0.GetUserPolicyResponse>(
          '/iam.IAMService/GetUserPolicy',
          ($0.GetUserPolicyRequest value) => value.writeToBuffer(),
          $0.GetUserPolicyResponse.fromBuffer);
  static final _$listAccessKeys =
      $grpc.ClientMethod<$0.ListAccessKeysRequest, $0.ListAccessKeysResponse>(
          '/iam.IAMService/ListAccessKeys',
          ($0.ListAccessKeysRequest value) => value.writeToBuffer(),
          $0.ListAccessKeysResponse.fromBuffer);
  static final _$listAccountAliases = $grpc.ClientMethod<
          $0.ListAccountAliasesRequest, $0.ListAccountAliasesResponse>(
      '/iam.IAMService/ListAccountAliases',
      ($0.ListAccountAliasesRequest value) => value.writeToBuffer(),
      $0.ListAccountAliasesResponse.fromBuffer);
  static final _$listAttachedGroupPolicies = $grpc.ClientMethod<
          $0.ListAttachedGroupPoliciesRequest,
          $0.ListAttachedGroupPoliciesResponse>(
      '/iam.IAMService/ListAttachedGroupPolicies',
      ($0.ListAttachedGroupPoliciesRequest value) => value.writeToBuffer(),
      $0.ListAttachedGroupPoliciesResponse.fromBuffer);
  static final _$listAttachedRolePolicies = $grpc.ClientMethod<
          $0.ListAttachedRolePoliciesRequest,
          $0.ListAttachedRolePoliciesResponse>(
      '/iam.IAMService/ListAttachedRolePolicies',
      ($0.ListAttachedRolePoliciesRequest value) => value.writeToBuffer(),
      $0.ListAttachedRolePoliciesResponse.fromBuffer);
  static final _$listAttachedUserPolicies = $grpc.ClientMethod<
          $0.ListAttachedUserPoliciesRequest,
          $0.ListAttachedUserPoliciesResponse>(
      '/iam.IAMService/ListAttachedUserPolicies',
      ($0.ListAttachedUserPoliciesRequest value) => value.writeToBuffer(),
      $0.ListAttachedUserPoliciesResponse.fromBuffer);
  static final _$listDelegationRequests = $grpc.ClientMethod<
          $0.ListDelegationRequestsRequest, $0.ListDelegationRequestsResponse>(
      '/iam.IAMService/ListDelegationRequests',
      ($0.ListDelegationRequestsRequest value) => value.writeToBuffer(),
      $0.ListDelegationRequestsResponse.fromBuffer);
  static final _$listEntitiesForPolicy = $grpc.ClientMethod<
          $0.ListEntitiesForPolicyRequest, $0.ListEntitiesForPolicyResponse>(
      '/iam.IAMService/ListEntitiesForPolicy',
      ($0.ListEntitiesForPolicyRequest value) => value.writeToBuffer(),
      $0.ListEntitiesForPolicyResponse.fromBuffer);
  static final _$listGroupPolicies = $grpc.ClientMethod<
          $0.ListGroupPoliciesRequest, $0.ListGroupPoliciesResponse>(
      '/iam.IAMService/ListGroupPolicies',
      ($0.ListGroupPoliciesRequest value) => value.writeToBuffer(),
      $0.ListGroupPoliciesResponse.fromBuffer);
  static final _$listGroups =
      $grpc.ClientMethod<$0.ListGroupsRequest, $0.ListGroupsResponse>(
          '/iam.IAMService/ListGroups',
          ($0.ListGroupsRequest value) => value.writeToBuffer(),
          $0.ListGroupsResponse.fromBuffer);
  static final _$listGroupsForUser = $grpc.ClientMethod<
          $0.ListGroupsForUserRequest, $0.ListGroupsForUserResponse>(
      '/iam.IAMService/ListGroupsForUser',
      ($0.ListGroupsForUserRequest value) => value.writeToBuffer(),
      $0.ListGroupsForUserResponse.fromBuffer);
  static final _$listInstanceProfiles = $grpc.ClientMethod<
          $0.ListInstanceProfilesRequest, $0.ListInstanceProfilesResponse>(
      '/iam.IAMService/ListInstanceProfiles',
      ($0.ListInstanceProfilesRequest value) => value.writeToBuffer(),
      $0.ListInstanceProfilesResponse.fromBuffer);
  static final _$listInstanceProfilesForRole = $grpc.ClientMethod<
          $0.ListInstanceProfilesForRoleRequest,
          $0.ListInstanceProfilesForRoleResponse>(
      '/iam.IAMService/ListInstanceProfilesForRole',
      ($0.ListInstanceProfilesForRoleRequest value) => value.writeToBuffer(),
      $0.ListInstanceProfilesForRoleResponse.fromBuffer);
  static final _$listInstanceProfileTags = $grpc.ClientMethod<
          $0.ListInstanceProfileTagsRequest,
          $0.ListInstanceProfileTagsResponse>(
      '/iam.IAMService/ListInstanceProfileTags',
      ($0.ListInstanceProfileTagsRequest value) => value.writeToBuffer(),
      $0.ListInstanceProfileTagsResponse.fromBuffer);
  static final _$listMFADevices =
      $grpc.ClientMethod<$0.ListMFADevicesRequest, $0.ListMFADevicesResponse>(
          '/iam.IAMService/ListMFADevices',
          ($0.ListMFADevicesRequest value) => value.writeToBuffer(),
          $0.ListMFADevicesResponse.fromBuffer);
  static final _$listMFADeviceTags = $grpc.ClientMethod<
          $0.ListMFADeviceTagsRequest, $0.ListMFADeviceTagsResponse>(
      '/iam.IAMService/ListMFADeviceTags',
      ($0.ListMFADeviceTagsRequest value) => value.writeToBuffer(),
      $0.ListMFADeviceTagsResponse.fromBuffer);
  static final _$listOpenIDConnectProviders = $grpc.ClientMethod<
          $0.ListOpenIDConnectProvidersRequest,
          $0.ListOpenIDConnectProvidersResponse>(
      '/iam.IAMService/ListOpenIDConnectProviders',
      ($0.ListOpenIDConnectProvidersRequest value) => value.writeToBuffer(),
      $0.ListOpenIDConnectProvidersResponse.fromBuffer);
  static final _$listOpenIDConnectProviderTags = $grpc.ClientMethod<
          $0.ListOpenIDConnectProviderTagsRequest,
          $0.ListOpenIDConnectProviderTagsResponse>(
      '/iam.IAMService/ListOpenIDConnectProviderTags',
      ($0.ListOpenIDConnectProviderTagsRequest value) => value.writeToBuffer(),
      $0.ListOpenIDConnectProviderTagsResponse.fromBuffer);
  static final _$listOrganizationsFeatures = $grpc.ClientMethod<
          $0.ListOrganizationsFeaturesRequest,
          $0.ListOrganizationsFeaturesResponse>(
      '/iam.IAMService/ListOrganizationsFeatures',
      ($0.ListOrganizationsFeaturesRequest value) => value.writeToBuffer(),
      $0.ListOrganizationsFeaturesResponse.fromBuffer);
  static final _$listPolicies =
      $grpc.ClientMethod<$0.ListPoliciesRequest, $0.ListPoliciesResponse>(
          '/iam.IAMService/ListPolicies',
          ($0.ListPoliciesRequest value) => value.writeToBuffer(),
          $0.ListPoliciesResponse.fromBuffer);
  static final _$listPoliciesGrantingServiceAccess = $grpc.ClientMethod<
          $0.ListPoliciesGrantingServiceAccessRequest,
          $0.ListPoliciesGrantingServiceAccessResponse>(
      '/iam.IAMService/ListPoliciesGrantingServiceAccess',
      ($0.ListPoliciesGrantingServiceAccessRequest value) =>
          value.writeToBuffer(),
      $0.ListPoliciesGrantingServiceAccessResponse.fromBuffer);
  static final _$listPolicyTags =
      $grpc.ClientMethod<$0.ListPolicyTagsRequest, $0.ListPolicyTagsResponse>(
          '/iam.IAMService/ListPolicyTags',
          ($0.ListPolicyTagsRequest value) => value.writeToBuffer(),
          $0.ListPolicyTagsResponse.fromBuffer);
  static final _$listPolicyVersions = $grpc.ClientMethod<
          $0.ListPolicyVersionsRequest, $0.ListPolicyVersionsResponse>(
      '/iam.IAMService/ListPolicyVersions',
      ($0.ListPolicyVersionsRequest value) => value.writeToBuffer(),
      $0.ListPolicyVersionsResponse.fromBuffer);
  static final _$listRolePolicies = $grpc.ClientMethod<
          $0.ListRolePoliciesRequest, $0.ListRolePoliciesResponse>(
      '/iam.IAMService/ListRolePolicies',
      ($0.ListRolePoliciesRequest value) => value.writeToBuffer(),
      $0.ListRolePoliciesResponse.fromBuffer);
  static final _$listRoles =
      $grpc.ClientMethod<$0.ListRolesRequest, $0.ListRolesResponse>(
          '/iam.IAMService/ListRoles',
          ($0.ListRolesRequest value) => value.writeToBuffer(),
          $0.ListRolesResponse.fromBuffer);
  static final _$listRoleTags =
      $grpc.ClientMethod<$0.ListRoleTagsRequest, $0.ListRoleTagsResponse>(
          '/iam.IAMService/ListRoleTags',
          ($0.ListRoleTagsRequest value) => value.writeToBuffer(),
          $0.ListRoleTagsResponse.fromBuffer);
  static final _$listSAMLProviders = $grpc.ClientMethod<
          $0.ListSAMLProvidersRequest, $0.ListSAMLProvidersResponse>(
      '/iam.IAMService/ListSAMLProviders',
      ($0.ListSAMLProvidersRequest value) => value.writeToBuffer(),
      $0.ListSAMLProvidersResponse.fromBuffer);
  static final _$listSAMLProviderTags = $grpc.ClientMethod<
          $0.ListSAMLProviderTagsRequest, $0.ListSAMLProviderTagsResponse>(
      '/iam.IAMService/ListSAMLProviderTags',
      ($0.ListSAMLProviderTagsRequest value) => value.writeToBuffer(),
      $0.ListSAMLProviderTagsResponse.fromBuffer);
  static final _$listServerCertificates = $grpc.ClientMethod<
          $0.ListServerCertificatesRequest, $0.ListServerCertificatesResponse>(
      '/iam.IAMService/ListServerCertificates',
      ($0.ListServerCertificatesRequest value) => value.writeToBuffer(),
      $0.ListServerCertificatesResponse.fromBuffer);
  static final _$listServerCertificateTags = $grpc.ClientMethod<
          $0.ListServerCertificateTagsRequest,
          $0.ListServerCertificateTagsResponse>(
      '/iam.IAMService/ListServerCertificateTags',
      ($0.ListServerCertificateTagsRequest value) => value.writeToBuffer(),
      $0.ListServerCertificateTagsResponse.fromBuffer);
  static final _$listServiceSpecificCredentials = $grpc.ClientMethod<
          $0.ListServiceSpecificCredentialsRequest,
          $0.ListServiceSpecificCredentialsResponse>(
      '/iam.IAMService/ListServiceSpecificCredentials',
      ($0.ListServiceSpecificCredentialsRequest value) => value.writeToBuffer(),
      $0.ListServiceSpecificCredentialsResponse.fromBuffer);
  static final _$listSigningCertificates = $grpc.ClientMethod<
          $0.ListSigningCertificatesRequest,
          $0.ListSigningCertificatesResponse>(
      '/iam.IAMService/ListSigningCertificates',
      ($0.ListSigningCertificatesRequest value) => value.writeToBuffer(),
      $0.ListSigningCertificatesResponse.fromBuffer);
  static final _$listSSHPublicKeys = $grpc.ClientMethod<
          $0.ListSSHPublicKeysRequest, $0.ListSSHPublicKeysResponse>(
      '/iam.IAMService/ListSSHPublicKeys',
      ($0.ListSSHPublicKeysRequest value) => value.writeToBuffer(),
      $0.ListSSHPublicKeysResponse.fromBuffer);
  static final _$listUserPolicies = $grpc.ClientMethod<
          $0.ListUserPoliciesRequest, $0.ListUserPoliciesResponse>(
      '/iam.IAMService/ListUserPolicies',
      ($0.ListUserPoliciesRequest value) => value.writeToBuffer(),
      $0.ListUserPoliciesResponse.fromBuffer);
  static final _$listUsers =
      $grpc.ClientMethod<$0.ListUsersRequest, $0.ListUsersResponse>(
          '/iam.IAMService/ListUsers',
          ($0.ListUsersRequest value) => value.writeToBuffer(),
          $0.ListUsersResponse.fromBuffer);
  static final _$listUserTags =
      $grpc.ClientMethod<$0.ListUserTagsRequest, $0.ListUserTagsResponse>(
          '/iam.IAMService/ListUserTags',
          ($0.ListUserTagsRequest value) => value.writeToBuffer(),
          $0.ListUserTagsResponse.fromBuffer);
  static final _$listVirtualMFADevices = $grpc.ClientMethod<
          $0.ListVirtualMFADevicesRequest, $0.ListVirtualMFADevicesResponse>(
      '/iam.IAMService/ListVirtualMFADevices',
      ($0.ListVirtualMFADevicesRequest value) => value.writeToBuffer(),
      $0.ListVirtualMFADevicesResponse.fromBuffer);
  static final _$putGroupPolicy =
      $grpc.ClientMethod<$0.PutGroupPolicyRequest, $1.Empty>(
          '/iam.IAMService/PutGroupPolicy',
          ($0.PutGroupPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putRolePermissionsBoundary =
      $grpc.ClientMethod<$0.PutRolePermissionsBoundaryRequest, $1.Empty>(
          '/iam.IAMService/PutRolePermissionsBoundary',
          ($0.PutRolePermissionsBoundaryRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putRolePolicy =
      $grpc.ClientMethod<$0.PutRolePolicyRequest, $1.Empty>(
          '/iam.IAMService/PutRolePolicy',
          ($0.PutRolePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putUserPermissionsBoundary =
      $grpc.ClientMethod<$0.PutUserPermissionsBoundaryRequest, $1.Empty>(
          '/iam.IAMService/PutUserPermissionsBoundary',
          ($0.PutUserPermissionsBoundaryRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$putUserPolicy =
      $grpc.ClientMethod<$0.PutUserPolicyRequest, $1.Empty>(
          '/iam.IAMService/PutUserPolicy',
          ($0.PutUserPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$rejectDelegationRequest =
      $grpc.ClientMethod<$0.RejectDelegationRequestRequest, $1.Empty>(
          '/iam.IAMService/RejectDelegationRequest',
          ($0.RejectDelegationRequestRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$removeClientIDFromOpenIDConnectProvider = $grpc.ClientMethod<
          $0.RemoveClientIDFromOpenIDConnectProviderRequest, $1.Empty>(
      '/iam.IAMService/RemoveClientIDFromOpenIDConnectProvider',
      ($0.RemoveClientIDFromOpenIDConnectProviderRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$removeRoleFromInstanceProfile =
      $grpc.ClientMethod<$0.RemoveRoleFromInstanceProfileRequest, $1.Empty>(
          '/iam.IAMService/RemoveRoleFromInstanceProfile',
          ($0.RemoveRoleFromInstanceProfileRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$removeUserFromGroup =
      $grpc.ClientMethod<$0.RemoveUserFromGroupRequest, $1.Empty>(
          '/iam.IAMService/RemoveUserFromGroup',
          ($0.RemoveUserFromGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$resetServiceSpecificCredential = $grpc.ClientMethod<
          $0.ResetServiceSpecificCredentialRequest,
          $0.ResetServiceSpecificCredentialResponse>(
      '/iam.IAMService/ResetServiceSpecificCredential',
      ($0.ResetServiceSpecificCredentialRequest value) => value.writeToBuffer(),
      $0.ResetServiceSpecificCredentialResponse.fromBuffer);
  static final _$resyncMFADevice =
      $grpc.ClientMethod<$0.ResyncMFADeviceRequest, $1.Empty>(
          '/iam.IAMService/ResyncMFADevice',
          ($0.ResyncMFADeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$sendDelegationToken =
      $grpc.ClientMethod<$0.SendDelegationTokenRequest, $1.Empty>(
          '/iam.IAMService/SendDelegationToken',
          ($0.SendDelegationTokenRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$setDefaultPolicyVersion =
      $grpc.ClientMethod<$0.SetDefaultPolicyVersionRequest, $1.Empty>(
          '/iam.IAMService/SetDefaultPolicyVersion',
          ($0.SetDefaultPolicyVersionRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$setSecurityTokenServicePreferences = $grpc.ClientMethod<
          $0.SetSecurityTokenServicePreferencesRequest, $1.Empty>(
      '/iam.IAMService/SetSecurityTokenServicePreferences',
      ($0.SetSecurityTokenServicePreferencesRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$simulateCustomPolicy = $grpc.ClientMethod<
          $0.SimulateCustomPolicyRequest, $0.SimulatePolicyResponse>(
      '/iam.IAMService/SimulateCustomPolicy',
      ($0.SimulateCustomPolicyRequest value) => value.writeToBuffer(),
      $0.SimulatePolicyResponse.fromBuffer);
  static final _$simulatePrincipalPolicy = $grpc.ClientMethod<
          $0.SimulatePrincipalPolicyRequest, $0.SimulatePolicyResponse>(
      '/iam.IAMService/SimulatePrincipalPolicy',
      ($0.SimulatePrincipalPolicyRequest value) => value.writeToBuffer(),
      $0.SimulatePolicyResponse.fromBuffer);
  static final _$tagInstanceProfile =
      $grpc.ClientMethod<$0.TagInstanceProfileRequest, $1.Empty>(
          '/iam.IAMService/TagInstanceProfile',
          ($0.TagInstanceProfileRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$tagMFADevice =
      $grpc.ClientMethod<$0.TagMFADeviceRequest, $1.Empty>(
          '/iam.IAMService/TagMFADevice',
          ($0.TagMFADeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$tagOpenIDConnectProvider =
      $grpc.ClientMethod<$0.TagOpenIDConnectProviderRequest, $1.Empty>(
          '/iam.IAMService/TagOpenIDConnectProvider',
          ($0.TagOpenIDConnectProviderRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$tagPolicy = $grpc.ClientMethod<$0.TagPolicyRequest, $1.Empty>(
      '/iam.IAMService/TagPolicy',
      ($0.TagPolicyRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$tagRole = $grpc.ClientMethod<$0.TagRoleRequest, $1.Empty>(
      '/iam.IAMService/TagRole',
      ($0.TagRoleRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$tagSAMLProvider =
      $grpc.ClientMethod<$0.TagSAMLProviderRequest, $1.Empty>(
          '/iam.IAMService/TagSAMLProvider',
          ($0.TagSAMLProviderRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$tagServerCertificate =
      $grpc.ClientMethod<$0.TagServerCertificateRequest, $1.Empty>(
          '/iam.IAMService/TagServerCertificate',
          ($0.TagServerCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$tagUser = $grpc.ClientMethod<$0.TagUserRequest, $1.Empty>(
      '/iam.IAMService/TagUser',
      ($0.TagUserRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$untagInstanceProfile =
      $grpc.ClientMethod<$0.UntagInstanceProfileRequest, $1.Empty>(
          '/iam.IAMService/UntagInstanceProfile',
          ($0.UntagInstanceProfileRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagMFADevice =
      $grpc.ClientMethod<$0.UntagMFADeviceRequest, $1.Empty>(
          '/iam.IAMService/UntagMFADevice',
          ($0.UntagMFADeviceRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagOpenIDConnectProvider =
      $grpc.ClientMethod<$0.UntagOpenIDConnectProviderRequest, $1.Empty>(
          '/iam.IAMService/UntagOpenIDConnectProvider',
          ($0.UntagOpenIDConnectProviderRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagPolicy =
      $grpc.ClientMethod<$0.UntagPolicyRequest, $1.Empty>(
          '/iam.IAMService/UntagPolicy',
          ($0.UntagPolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagRole = $grpc.ClientMethod<$0.UntagRoleRequest, $1.Empty>(
      '/iam.IAMService/UntagRole',
      ($0.UntagRoleRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$untagSAMLProvider =
      $grpc.ClientMethod<$0.UntagSAMLProviderRequest, $1.Empty>(
          '/iam.IAMService/UntagSAMLProvider',
          ($0.UntagSAMLProviderRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagServerCertificate =
      $grpc.ClientMethod<$0.UntagServerCertificateRequest, $1.Empty>(
          '/iam.IAMService/UntagServerCertificate',
          ($0.UntagServerCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$untagUser = $grpc.ClientMethod<$0.UntagUserRequest, $1.Empty>(
      '/iam.IAMService/UntagUser',
      ($0.UntagUserRequest value) => value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$updateAccessKey =
      $grpc.ClientMethod<$0.UpdateAccessKeyRequest, $1.Empty>(
          '/iam.IAMService/UpdateAccessKey',
          ($0.UpdateAccessKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateAccountPasswordPolicy =
      $grpc.ClientMethod<$0.UpdateAccountPasswordPolicyRequest, $1.Empty>(
          '/iam.IAMService/UpdateAccountPasswordPolicy',
          ($0.UpdateAccountPasswordPolicyRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateAssumeRolePolicy =
      $grpc.ClientMethod<$0.UpdateAssumeRolePolicyRequest, $1.Empty>(
          '/iam.IAMService/UpdateAssumeRolePolicy',
          ($0.UpdateAssumeRolePolicyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateDelegationRequest =
      $grpc.ClientMethod<$0.UpdateDelegationRequestRequest, $1.Empty>(
          '/iam.IAMService/UpdateDelegationRequest',
          ($0.UpdateDelegationRequestRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateGroup =
      $grpc.ClientMethod<$0.UpdateGroupRequest, $1.Empty>(
          '/iam.IAMService/UpdateGroup',
          ($0.UpdateGroupRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateLoginProfile =
      $grpc.ClientMethod<$0.UpdateLoginProfileRequest, $1.Empty>(
          '/iam.IAMService/UpdateLoginProfile',
          ($0.UpdateLoginProfileRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateOpenIDConnectProviderThumbprint = $grpc.ClientMethod<
          $0.UpdateOpenIDConnectProviderThumbprintRequest, $1.Empty>(
      '/iam.IAMService/UpdateOpenIDConnectProviderThumbprint',
      ($0.UpdateOpenIDConnectProviderThumbprintRequest value) =>
          value.writeToBuffer(),
      $1.Empty.fromBuffer);
  static final _$updateRole =
      $grpc.ClientMethod<$0.UpdateRoleRequest, $0.UpdateRoleResponse>(
          '/iam.IAMService/UpdateRole',
          ($0.UpdateRoleRequest value) => value.writeToBuffer(),
          $0.UpdateRoleResponse.fromBuffer);
  static final _$updateRoleDescription = $grpc.ClientMethod<
          $0.UpdateRoleDescriptionRequest, $0.UpdateRoleDescriptionResponse>(
      '/iam.IAMService/UpdateRoleDescription',
      ($0.UpdateRoleDescriptionRequest value) => value.writeToBuffer(),
      $0.UpdateRoleDescriptionResponse.fromBuffer);
  static final _$updateSAMLProvider = $grpc.ClientMethod<
          $0.UpdateSAMLProviderRequest, $0.UpdateSAMLProviderResponse>(
      '/iam.IAMService/UpdateSAMLProvider',
      ($0.UpdateSAMLProviderRequest value) => value.writeToBuffer(),
      $0.UpdateSAMLProviderResponse.fromBuffer);
  static final _$updateServerCertificate =
      $grpc.ClientMethod<$0.UpdateServerCertificateRequest, $1.Empty>(
          '/iam.IAMService/UpdateServerCertificate',
          ($0.UpdateServerCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateServiceSpecificCredential =
      $grpc.ClientMethod<$0.UpdateServiceSpecificCredentialRequest, $1.Empty>(
          '/iam.IAMService/UpdateServiceSpecificCredential',
          ($0.UpdateServiceSpecificCredentialRequest value) =>
              value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateSigningCertificate =
      $grpc.ClientMethod<$0.UpdateSigningCertificateRequest, $1.Empty>(
          '/iam.IAMService/UpdateSigningCertificate',
          ($0.UpdateSigningCertificateRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateSSHPublicKey =
      $grpc.ClientMethod<$0.UpdateSSHPublicKeyRequest, $1.Empty>(
          '/iam.IAMService/UpdateSSHPublicKey',
          ($0.UpdateSSHPublicKeyRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$updateUser =
      $grpc.ClientMethod<$0.UpdateUserRequest, $1.Empty>(
          '/iam.IAMService/UpdateUser',
          ($0.UpdateUserRequest value) => value.writeToBuffer(),
          $1.Empty.fromBuffer);
  static final _$uploadServerCertificate = $grpc.ClientMethod<
          $0.UploadServerCertificateRequest,
          $0.UploadServerCertificateResponse>(
      '/iam.IAMService/UploadServerCertificate',
      ($0.UploadServerCertificateRequest value) => value.writeToBuffer(),
      $0.UploadServerCertificateResponse.fromBuffer);
  static final _$uploadSigningCertificate = $grpc.ClientMethod<
          $0.UploadSigningCertificateRequest,
          $0.UploadSigningCertificateResponse>(
      '/iam.IAMService/UploadSigningCertificate',
      ($0.UploadSigningCertificateRequest value) => value.writeToBuffer(),
      $0.UploadSigningCertificateResponse.fromBuffer);
  static final _$uploadSSHPublicKey = $grpc.ClientMethod<
          $0.UploadSSHPublicKeyRequest, $0.UploadSSHPublicKeyResponse>(
      '/iam.IAMService/UploadSSHPublicKey',
      ($0.UploadSSHPublicKeyRequest value) => value.writeToBuffer(),
      $0.UploadSSHPublicKeyResponse.fromBuffer);
}

@$pb.GrpcServiceName('iam.IAMService')
abstract class IAMServiceBase extends $grpc.Service {
  $core.String get $name => 'iam.IAMService';

  IAMServiceBase() {
    $addMethod($grpc.ServiceMethod<$0.AcceptDelegationRequestRequest, $1.Empty>(
        'AcceptDelegationRequest',
        acceptDelegationRequest_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AcceptDelegationRequestRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AddClientIDToOpenIDConnectProviderRequest,
            $1.Empty>(
        'AddClientIDToOpenIDConnectProvider',
        addClientIDToOpenIDConnectProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddClientIDToOpenIDConnectProviderRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.AddRoleToInstanceProfileRequest, $1.Empty>(
            'AddRoleToInstanceProfile',
            addRoleToInstanceProfile_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.AddRoleToInstanceProfileRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AddUserToGroupRequest, $1.Empty>(
        'AddUserToGroup',
        addUserToGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AddUserToGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.AssociateDelegationRequestRequest, $1.Empty>(
            'AssociateDelegationRequest',
            associateDelegationRequest_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.AssociateDelegationRequestRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AttachGroupPolicyRequest, $1.Empty>(
        'AttachGroupPolicy',
        attachGroupPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AttachGroupPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AttachRolePolicyRequest, $1.Empty>(
        'AttachRolePolicy',
        attachRolePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AttachRolePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.AttachUserPolicyRequest, $1.Empty>(
        'AttachUserPolicy',
        attachUserPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.AttachUserPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ChangePasswordRequest, $1.Empty>(
        'ChangePassword',
        changePassword_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ChangePasswordRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateAccessKeyRequest,
            $0.CreateAccessKeyResponse>(
        'CreateAccessKey',
        createAccessKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAccessKeyRequest.fromBuffer(value),
        ($0.CreateAccessKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateAccountAliasRequest, $1.Empty>(
        'CreateAccountAlias',
        createAccountAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateAccountAliasRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateDelegationRequestRequest,
            $0.CreateDelegationRequestResponse>(
        'CreateDelegationRequest',
        createDelegationRequest_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateDelegationRequestRequest.fromBuffer(value),
        ($0.CreateDelegationRequestResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreateGroupRequest, $0.CreateGroupResponse>(
            'CreateGroup',
            createGroup_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreateGroupRequest.fromBuffer(value),
            ($0.CreateGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateInstanceProfileRequest,
            $0.CreateInstanceProfileResponse>(
        'CreateInstanceProfile',
        createInstanceProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateInstanceProfileRequest.fromBuffer(value),
        ($0.CreateInstanceProfileResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateLoginProfileRequest,
            $0.CreateLoginProfileResponse>(
        'CreateLoginProfile',
        createLoginProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateLoginProfileRequest.fromBuffer(value),
        ($0.CreateLoginProfileResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateOpenIDConnectProviderRequest,
            $0.CreateOpenIDConnectProviderResponse>(
        'CreateOpenIDConnectProvider',
        createOpenIDConnectProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateOpenIDConnectProviderRequest.fromBuffer(value),
        ($0.CreateOpenIDConnectProviderResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.CreatePolicyRequest, $0.CreatePolicyResponse>(
            'CreatePolicy',
            createPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.CreatePolicyRequest.fromBuffer(value),
            ($0.CreatePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreatePolicyVersionRequest,
            $0.CreatePolicyVersionResponse>(
        'CreatePolicyVersion',
        createPolicyVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreatePolicyVersionRequest.fromBuffer(value),
        ($0.CreatePolicyVersionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateRoleRequest, $0.CreateRoleResponse>(
        'CreateRole',
        createRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateRoleRequest.fromBuffer(value),
        ($0.CreateRoleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateSAMLProviderRequest,
            $0.CreateSAMLProviderResponse>(
        'CreateSAMLProvider',
        createSAMLProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateSAMLProviderRequest.fromBuffer(value),
        ($0.CreateSAMLProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateServiceLinkedRoleRequest,
            $0.CreateServiceLinkedRoleResponse>(
        'CreateServiceLinkedRole',
        createServiceLinkedRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateServiceLinkedRoleRequest.fromBuffer(value),
        ($0.CreateServiceLinkedRoleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateServiceSpecificCredentialRequest,
            $0.CreateServiceSpecificCredentialResponse>(
        'CreateServiceSpecificCredential',
        createServiceSpecificCredential_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateServiceSpecificCredentialRequest.fromBuffer(value),
        ($0.CreateServiceSpecificCredentialResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateUserRequest, $0.CreateUserResponse>(
        'CreateUser',
        createUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.CreateUserRequest.fromBuffer(value),
        ($0.CreateUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.CreateVirtualMFADeviceRequest,
            $0.CreateVirtualMFADeviceResponse>(
        'CreateVirtualMFADevice',
        createVirtualMFADevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.CreateVirtualMFADeviceRequest.fromBuffer(value),
        ($0.CreateVirtualMFADeviceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeactivateMFADeviceRequest, $1.Empty>(
        'DeactivateMFADevice',
        deactivateMFADevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeactivateMFADeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAccessKeyRequest, $1.Empty>(
        'DeleteAccessKey',
        deleteAccessKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAccessKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteAccountAliasRequest, $1.Empty>(
        'DeleteAccountAlias',
        deleteAccountAlias_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteAccountAliasRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$1.Empty, $1.Empty>(
        'DeleteAccountPasswordPolicy',
        deleteAccountPasswordPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteGroupRequest, $1.Empty>(
        'DeleteGroup',
        deleteGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteGroupPolicyRequest, $1.Empty>(
        'DeleteGroupPolicy',
        deleteGroupPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteGroupPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteInstanceProfileRequest, $1.Empty>(
        'DeleteInstanceProfile',
        deleteInstanceProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteInstanceProfileRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteLoginProfileRequest, $1.Empty>(
        'DeleteLoginProfile',
        deleteLoginProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteLoginProfileRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteOpenIDConnectProviderRequest, $1.Empty>(
            'DeleteOpenIDConnectProvider',
            deleteOpenIDConnectProvider_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteOpenIDConnectProviderRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePolicyRequest, $1.Empty>(
        'DeletePolicy',
        deletePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeletePolicyVersionRequest, $1.Empty>(
        'DeletePolicyVersion',
        deletePolicyVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeletePolicyVersionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRoleRequest, $1.Empty>(
        'DeleteRole',
        deleteRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteRoleRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteRolePermissionsBoundaryRequest, $1.Empty>(
            'DeleteRolePermissionsBoundary',
            deleteRolePermissionsBoundary_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteRolePermissionsBoundaryRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteRolePolicyRequest, $1.Empty>(
        'DeleteRolePolicy',
        deleteRolePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteRolePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteSAMLProviderRequest, $1.Empty>(
        'DeleteSAMLProvider',
        deleteSAMLProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteSAMLProviderRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteServerCertificateRequest, $1.Empty>(
        'DeleteServerCertificate',
        deleteServerCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteServerCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteServiceLinkedRoleRequest,
            $0.DeleteServiceLinkedRoleResponse>(
        'DeleteServiceLinkedRole',
        deleteServiceLinkedRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteServiceLinkedRoleRequest.fromBuffer(value),
        ($0.DeleteServiceLinkedRoleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteServiceSpecificCredentialRequest,
            $1.Empty>(
        'DeleteServiceSpecificCredential',
        deleteServiceSpecificCredential_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteServiceSpecificCredentialRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteSigningCertificateRequest, $1.Empty>(
            'DeleteSigningCertificate',
            deleteSigningCertificate_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteSigningCertificateRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteSSHPublicKeyRequest, $1.Empty>(
        'DeleteSSHPublicKey',
        deleteSSHPublicKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteSSHPublicKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserRequest, $1.Empty>(
        'DeleteUser',
        deleteUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.DeleteUserRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.DeleteUserPermissionsBoundaryRequest, $1.Empty>(
            'DeleteUserPermissionsBoundary',
            deleteUserPermissionsBoundary_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.DeleteUserPermissionsBoundaryRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteUserPolicyRequest, $1.Empty>(
        'DeleteUserPolicy',
        deleteUserPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteUserPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DeleteVirtualMFADeviceRequest, $1.Empty>(
        'DeleteVirtualMFADevice',
        deleteVirtualMFADevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DeleteVirtualMFADeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DetachGroupPolicyRequest, $1.Empty>(
        'DetachGroupPolicy',
        detachGroupPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DetachGroupPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DetachRolePolicyRequest, $1.Empty>(
        'DetachRolePolicy',
        detachRolePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DetachRolePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DetachUserPolicyRequest, $1.Empty>(
        'DetachUserPolicy',
        detachUserPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DetachUserPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.DisableOrganizationsRootCredentialsManagementRequest,
            $0.DisableOrganizationsRootCredentialsManagementResponse>(
        'DisableOrganizationsRootCredentialsManagement',
        disableOrganizationsRootCredentialsManagement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableOrganizationsRootCredentialsManagementRequest.fromBuffer(
                value),
        ($0.DisableOrganizationsRootCredentialsManagementResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.DisableOrganizationsRootSessionsRequest,
            $0.DisableOrganizationsRootSessionsResponse>(
        'DisableOrganizationsRootSessions',
        disableOrganizationsRootSessions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.DisableOrganizationsRootSessionsRequest.fromBuffer(value),
        ($0.DisableOrganizationsRootSessionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$1.Empty, $1.Empty>(
        'DisableOutboundWebIdentityFederation',
        disableOutboundWebIdentityFederation_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableMFADeviceRequest, $1.Empty>(
        'EnableMFADevice',
        enableMFADevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableMFADeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.EnableOrganizationsRootCredentialsManagementRequest,
            $0.EnableOrganizationsRootCredentialsManagementResponse>(
        'EnableOrganizationsRootCredentialsManagement',
        enableOrganizationsRootCredentialsManagement_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableOrganizationsRootCredentialsManagementRequest.fromBuffer(
                value),
        ($0.EnableOrganizationsRootCredentialsManagementResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.EnableOrganizationsRootSessionsRequest,
            $0.EnableOrganizationsRootSessionsResponse>(
        'EnableOrganizationsRootSessions',
        enableOrganizationsRootSessions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.EnableOrganizationsRootSessionsRequest.fromBuffer(value),
        ($0.EnableOrganizationsRootSessionsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$1.Empty,
            $0.EnableOutboundWebIdentityFederationResponse>(
        'EnableOutboundWebIdentityFederation',
        enableOutboundWebIdentityFederation_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
        ($0.EnableOutboundWebIdentityFederationResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$1.Empty, $0.GenerateCredentialReportResponse>(
            'GenerateCredentialReport',
            generateCredentialReport_Pre,
            false,
            false,
            ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
            ($0.GenerateCredentialReportResponse value) =>
                value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateOrganizationsAccessReportRequest,
            $0.GenerateOrganizationsAccessReportResponse>(
        'GenerateOrganizationsAccessReport',
        generateOrganizationsAccessReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateOrganizationsAccessReportRequest.fromBuffer(value),
        ($0.GenerateOrganizationsAccessReportResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GenerateServiceLastAccessedDetailsRequest,
            $0.GenerateServiceLastAccessedDetailsResponse>(
        'GenerateServiceLastAccessedDetails',
        generateServiceLastAccessedDetails_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GenerateServiceLastAccessedDetailsRequest.fromBuffer(value),
        ($0.GenerateServiceLastAccessedDetailsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccessKeyLastUsedRequest,
            $0.GetAccessKeyLastUsedResponse>(
        'GetAccessKeyLastUsed',
        getAccessKeyLastUsed_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAccessKeyLastUsedRequest.fromBuffer(value),
        ($0.GetAccessKeyLastUsedResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetAccountAuthorizationDetailsRequest,
            $0.GetAccountAuthorizationDetailsResponse>(
        'GetAccountAuthorizationDetails',
        getAccountAuthorizationDetails_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetAccountAuthorizationDetailsRequest.fromBuffer(value),
        ($0.GetAccountAuthorizationDetailsResponse value) =>
            value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$1.Empty, $0.GetAccountPasswordPolicyResponse>(
            'GetAccountPasswordPolicy',
            getAccountPasswordPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
            ($0.GetAccountPasswordPolicyResponse value) =>
                value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$1.Empty, $0.GetAccountSummaryResponse>(
        'GetAccountSummary',
        getAccountSummary_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
        ($0.GetAccountSummaryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetContextKeysForCustomPolicyRequest,
            $0.GetContextKeysForPolicyResponse>(
        'GetContextKeysForCustomPolicy',
        getContextKeysForCustomPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetContextKeysForCustomPolicyRequest.fromBuffer(value),
        ($0.GetContextKeysForPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetContextKeysForPrincipalPolicyRequest,
            $0.GetContextKeysForPolicyResponse>(
        'GetContextKeysForPrincipalPolicy',
        getContextKeysForPrincipalPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetContextKeysForPrincipalPolicyRequest.fromBuffer(value),
        ($0.GetContextKeysForPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$1.Empty, $0.GetCredentialReportResponse>(
        'GetCredentialReport',
        getCredentialReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
        ($0.GetCredentialReportResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetDelegationRequestRequest,
            $0.GetDelegationRequestResponse>(
        'GetDelegationRequest',
        getDelegationRequest_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetDelegationRequestRequest.fromBuffer(value),
        ($0.GetDelegationRequestResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetGroupRequest, $0.GetGroupResponse>(
        'GetGroup',
        getGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetGroupRequest.fromBuffer(value),
        ($0.GetGroupResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetGroupPolicyRequest,
            $0.GetGroupPolicyResponse>(
        'GetGroupPolicy',
        getGroupPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetGroupPolicyRequest.fromBuffer(value),
        ($0.GetGroupPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetHumanReadableSummaryRequest,
            $0.GetHumanReadableSummaryResponse>(
        'GetHumanReadableSummary',
        getHumanReadableSummary_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetHumanReadableSummaryRequest.fromBuffer(value),
        ($0.GetHumanReadableSummaryResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetInstanceProfileRequest,
            $0.GetInstanceProfileResponse>(
        'GetInstanceProfile',
        getInstanceProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetInstanceProfileRequest.fromBuffer(value),
        ($0.GetInstanceProfileResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetLoginProfileRequest,
            $0.GetLoginProfileResponse>(
        'GetLoginProfile',
        getLoginProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetLoginProfileRequest.fromBuffer(value),
        ($0.GetLoginProfileResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetMFADeviceRequest, $0.GetMFADeviceResponse>(
            'GetMFADevice',
            getMFADevice_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetMFADeviceRequest.fromBuffer(value),
            ($0.GetMFADeviceResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOpenIDConnectProviderRequest,
            $0.GetOpenIDConnectProviderResponse>(
        'GetOpenIDConnectProvider',
        getOpenIDConnectProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOpenIDConnectProviderRequest.fromBuffer(value),
        ($0.GetOpenIDConnectProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetOrganizationsAccessReportRequest,
            $0.GetOrganizationsAccessReportResponse>(
        'GetOrganizationsAccessReport',
        getOrganizationsAccessReport_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetOrganizationsAccessReportRequest.fromBuffer(value),
        ($0.GetOrganizationsAccessReportResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$1.Empty,
            $0.GetOutboundWebIdentityFederationInfoResponse>(
        'GetOutboundWebIdentityFederationInfo',
        getOutboundWebIdentityFederationInfo_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $1.Empty.fromBuffer(value),
        ($0.GetOutboundWebIdentityFederationInfoResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPolicyRequest, $0.GetPolicyResponse>(
        'GetPolicy',
        getPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetPolicyRequest.fromBuffer(value),
        ($0.GetPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetPolicyVersionRequest,
            $0.GetPolicyVersionResponse>(
        'GetPolicyVersion',
        getPolicyVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetPolicyVersionRequest.fromBuffer(value),
        ($0.GetPolicyVersionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetRoleRequest, $0.GetRoleResponse>(
        'GetRole',
        getRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetRoleRequest.fromBuffer(value),
        ($0.GetRoleResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetRolePolicyRequest, $0.GetRolePolicyResponse>(
            'GetRolePolicy',
            getRolePolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetRolePolicyRequest.fromBuffer(value),
            ($0.GetRolePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSAMLProviderRequest,
            $0.GetSAMLProviderResponse>(
        'GetSAMLProvider',
        getSAMLProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSAMLProviderRequest.fromBuffer(value),
        ($0.GetSAMLProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetServerCertificateRequest,
            $0.GetServerCertificateResponse>(
        'GetServerCertificate',
        getServerCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetServerCertificateRequest.fromBuffer(value),
        ($0.GetServerCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetServiceLastAccessedDetailsRequest,
            $0.GetServiceLastAccessedDetailsResponse>(
        'GetServiceLastAccessedDetails',
        getServiceLastAccessedDetails_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetServiceLastAccessedDetailsRequest.fromBuffer(value),
        ($0.GetServiceLastAccessedDetailsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.GetServiceLastAccessedDetailsWithEntitiesRequest,
            $0.GetServiceLastAccessedDetailsWithEntitiesResponse>(
        'GetServiceLastAccessedDetailsWithEntities',
        getServiceLastAccessedDetailsWithEntities_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetServiceLastAccessedDetailsWithEntitiesRequest.fromBuffer(
                value),
        ($0.GetServiceLastAccessedDetailsWithEntitiesResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetServiceLinkedRoleDeletionStatusRequest,
            $0.GetServiceLinkedRoleDeletionStatusResponse>(
        'GetServiceLinkedRoleDeletionStatus',
        getServiceLinkedRoleDeletionStatus_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetServiceLinkedRoleDeletionStatusRequest.fromBuffer(value),
        ($0.GetServiceLinkedRoleDeletionStatusResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetSSHPublicKeyRequest,
            $0.GetSSHPublicKeyResponse>(
        'GetSSHPublicKey',
        getSSHPublicKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.GetSSHPublicKeyRequest.fromBuffer(value),
        ($0.GetSSHPublicKeyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.GetUserRequest, $0.GetUserResponse>(
        'GetUser',
        getUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.GetUserRequest.fromBuffer(value),
        ($0.GetUserResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.GetUserPolicyRequest, $0.GetUserPolicyResponse>(
            'GetUserPolicy',
            getUserPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.GetUserPolicyRequest.fromBuffer(value),
            ($0.GetUserPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAccessKeysRequest,
            $0.ListAccessKeysResponse>(
        'ListAccessKeys',
        listAccessKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAccessKeysRequest.fromBuffer(value),
        ($0.ListAccessKeysResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAccountAliasesRequest,
            $0.ListAccountAliasesResponse>(
        'ListAccountAliases',
        listAccountAliases_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAccountAliasesRequest.fromBuffer(value),
        ($0.ListAccountAliasesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAttachedGroupPoliciesRequest,
            $0.ListAttachedGroupPoliciesResponse>(
        'ListAttachedGroupPolicies',
        listAttachedGroupPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAttachedGroupPoliciesRequest.fromBuffer(value),
        ($0.ListAttachedGroupPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAttachedRolePoliciesRequest,
            $0.ListAttachedRolePoliciesResponse>(
        'ListAttachedRolePolicies',
        listAttachedRolePolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAttachedRolePoliciesRequest.fromBuffer(value),
        ($0.ListAttachedRolePoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListAttachedUserPoliciesRequest,
            $0.ListAttachedUserPoliciesResponse>(
        'ListAttachedUserPolicies',
        listAttachedUserPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListAttachedUserPoliciesRequest.fromBuffer(value),
        ($0.ListAttachedUserPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListDelegationRequestsRequest,
            $0.ListDelegationRequestsResponse>(
        'ListDelegationRequests',
        listDelegationRequests_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListDelegationRequestsRequest.fromBuffer(value),
        ($0.ListDelegationRequestsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListEntitiesForPolicyRequest,
            $0.ListEntitiesForPolicyResponse>(
        'ListEntitiesForPolicy',
        listEntitiesForPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListEntitiesForPolicyRequest.fromBuffer(value),
        ($0.ListEntitiesForPolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGroupPoliciesRequest,
            $0.ListGroupPoliciesResponse>(
        'ListGroupPolicies',
        listGroupPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListGroupPoliciesRequest.fromBuffer(value),
        ($0.ListGroupPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGroupsRequest, $0.ListGroupsResponse>(
        'ListGroups',
        listGroups_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListGroupsRequest.fromBuffer(value),
        ($0.ListGroupsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListGroupsForUserRequest,
            $0.ListGroupsForUserResponse>(
        'ListGroupsForUser',
        listGroupsForUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListGroupsForUserRequest.fromBuffer(value),
        ($0.ListGroupsForUserResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListInstanceProfilesRequest,
            $0.ListInstanceProfilesResponse>(
        'ListInstanceProfiles',
        listInstanceProfiles_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInstanceProfilesRequest.fromBuffer(value),
        ($0.ListInstanceProfilesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListInstanceProfilesForRoleRequest,
            $0.ListInstanceProfilesForRoleResponse>(
        'ListInstanceProfilesForRole',
        listInstanceProfilesForRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInstanceProfilesForRoleRequest.fromBuffer(value),
        ($0.ListInstanceProfilesForRoleResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListInstanceProfileTagsRequest,
            $0.ListInstanceProfileTagsResponse>(
        'ListInstanceProfileTags',
        listInstanceProfileTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListInstanceProfileTagsRequest.fromBuffer(value),
        ($0.ListInstanceProfileTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMFADevicesRequest,
            $0.ListMFADevicesResponse>(
        'ListMFADevices',
        listMFADevices_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListMFADevicesRequest.fromBuffer(value),
        ($0.ListMFADevicesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListMFADeviceTagsRequest,
            $0.ListMFADeviceTagsResponse>(
        'ListMFADeviceTags',
        listMFADeviceTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListMFADeviceTagsRequest.fromBuffer(value),
        ($0.ListMFADeviceTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOpenIDConnectProvidersRequest,
            $0.ListOpenIDConnectProvidersResponse>(
        'ListOpenIDConnectProviders',
        listOpenIDConnectProviders_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOpenIDConnectProvidersRequest.fromBuffer(value),
        ($0.ListOpenIDConnectProvidersResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOpenIDConnectProviderTagsRequest,
            $0.ListOpenIDConnectProviderTagsResponse>(
        'ListOpenIDConnectProviderTags',
        listOpenIDConnectProviderTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOpenIDConnectProviderTagsRequest.fromBuffer(value),
        ($0.ListOpenIDConnectProviderTagsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListOrganizationsFeaturesRequest,
            $0.ListOrganizationsFeaturesResponse>(
        'ListOrganizationsFeatures',
        listOrganizationsFeatures_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListOrganizationsFeaturesRequest.fromBuffer(value),
        ($0.ListOrganizationsFeaturesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListPoliciesRequest, $0.ListPoliciesResponse>(
            'ListPolicies',
            listPolicies_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListPoliciesRequest.fromBuffer(value),
            ($0.ListPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPoliciesGrantingServiceAccessRequest,
            $0.ListPoliciesGrantingServiceAccessResponse>(
        'ListPoliciesGrantingServiceAccess',
        listPoliciesGrantingServiceAccess_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPoliciesGrantingServiceAccessRequest.fromBuffer(value),
        ($0.ListPoliciesGrantingServiceAccessResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPolicyTagsRequest,
            $0.ListPolicyTagsResponse>(
        'ListPolicyTags',
        listPolicyTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPolicyTagsRequest.fromBuffer(value),
        ($0.ListPolicyTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListPolicyVersionsRequest,
            $0.ListPolicyVersionsResponse>(
        'ListPolicyVersions',
        listPolicyVersions_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListPolicyVersionsRequest.fromBuffer(value),
        ($0.ListPolicyVersionsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRolePoliciesRequest,
            $0.ListRolePoliciesResponse>(
        'ListRolePolicies',
        listRolePolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListRolePoliciesRequest.fromBuffer(value),
        ($0.ListRolePoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListRolesRequest, $0.ListRolesResponse>(
        'ListRoles',
        listRoles_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListRolesRequest.fromBuffer(value),
        ($0.ListRolesResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListRoleTagsRequest, $0.ListRoleTagsResponse>(
            'ListRoleTags',
            listRoleTags_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListRoleTagsRequest.fromBuffer(value),
            ($0.ListRoleTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSAMLProvidersRequest,
            $0.ListSAMLProvidersResponse>(
        'ListSAMLProviders',
        listSAMLProviders_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSAMLProvidersRequest.fromBuffer(value),
        ($0.ListSAMLProvidersResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSAMLProviderTagsRequest,
            $0.ListSAMLProviderTagsResponse>(
        'ListSAMLProviderTags',
        listSAMLProviderTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSAMLProviderTagsRequest.fromBuffer(value),
        ($0.ListSAMLProviderTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListServerCertificatesRequest,
            $0.ListServerCertificatesResponse>(
        'ListServerCertificates',
        listServerCertificates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListServerCertificatesRequest.fromBuffer(value),
        ($0.ListServerCertificatesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListServerCertificateTagsRequest,
            $0.ListServerCertificateTagsResponse>(
        'ListServerCertificateTags',
        listServerCertificateTags_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListServerCertificateTagsRequest.fromBuffer(value),
        ($0.ListServerCertificateTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListServiceSpecificCredentialsRequest,
            $0.ListServiceSpecificCredentialsResponse>(
        'ListServiceSpecificCredentials',
        listServiceSpecificCredentials_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListServiceSpecificCredentialsRequest.fromBuffer(value),
        ($0.ListServiceSpecificCredentialsResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSigningCertificatesRequest,
            $0.ListSigningCertificatesResponse>(
        'ListSigningCertificates',
        listSigningCertificates_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSigningCertificatesRequest.fromBuffer(value),
        ($0.ListSigningCertificatesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListSSHPublicKeysRequest,
            $0.ListSSHPublicKeysResponse>(
        'ListSSHPublicKeys',
        listSSHPublicKeys_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListSSHPublicKeysRequest.fromBuffer(value),
        ($0.ListSSHPublicKeysResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListUserPoliciesRequest,
            $0.ListUserPoliciesResponse>(
        'ListUserPolicies',
        listUserPolicies_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListUserPoliciesRequest.fromBuffer(value),
        ($0.ListUserPoliciesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListUsersRequest, $0.ListUsersResponse>(
        'ListUsers',
        listUsers_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.ListUsersRequest.fromBuffer(value),
        ($0.ListUsersResponse value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.ListUserTagsRequest, $0.ListUserTagsResponse>(
            'ListUserTags',
            listUserTags_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.ListUserTagsRequest.fromBuffer(value),
            ($0.ListUserTagsResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ListVirtualMFADevicesRequest,
            $0.ListVirtualMFADevicesResponse>(
        'ListVirtualMFADevices',
        listVirtualMFADevices_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ListVirtualMFADevicesRequest.fromBuffer(value),
        ($0.ListVirtualMFADevicesResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutGroupPolicyRequest, $1.Empty>(
        'PutGroupPolicy',
        putGroupPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutGroupPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutRolePermissionsBoundaryRequest, $1.Empty>(
            'PutRolePermissionsBoundary',
            putRolePermissionsBoundary_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutRolePermissionsBoundaryRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutRolePolicyRequest, $1.Empty>(
        'PutRolePolicy',
        putRolePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutRolePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.PutUserPermissionsBoundaryRequest, $1.Empty>(
            'PutUserPermissionsBoundary',
            putUserPermissionsBoundary_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.PutUserPermissionsBoundaryRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.PutUserPolicyRequest, $1.Empty>(
        'PutUserPolicy',
        putUserPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.PutUserPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RejectDelegationRequestRequest, $1.Empty>(
        'RejectDelegationRequest',
        rejectDelegationRequest_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RejectDelegationRequestRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.RemoveClientIDFromOpenIDConnectProviderRequest, $1.Empty>(
        'RemoveClientIDFromOpenIDConnectProvider',
        removeClientIDFromOpenIDConnectProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveClientIDFromOpenIDConnectProviderRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.RemoveRoleFromInstanceProfileRequest, $1.Empty>(
            'RemoveRoleFromInstanceProfile',
            removeRoleFromInstanceProfile_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.RemoveRoleFromInstanceProfileRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.RemoveUserFromGroupRequest, $1.Empty>(
        'RemoveUserFromGroup',
        removeUserFromGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.RemoveUserFromGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResetServiceSpecificCredentialRequest,
            $0.ResetServiceSpecificCredentialResponse>(
        'ResetServiceSpecificCredential',
        resetServiceSpecificCredential_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResetServiceSpecificCredentialRequest.fromBuffer(value),
        ($0.ResetServiceSpecificCredentialResponse value) =>
            value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.ResyncMFADeviceRequest, $1.Empty>(
        'ResyncMFADevice',
        resyncMFADevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.ResyncMFADeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SendDelegationTokenRequest, $1.Empty>(
        'SendDelegationToken',
        sendDelegationToken_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SendDelegationTokenRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetDefaultPolicyVersionRequest, $1.Empty>(
        'SetDefaultPolicyVersion',
        setDefaultPolicyVersion_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetDefaultPolicyVersionRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SetSecurityTokenServicePreferencesRequest,
            $1.Empty>(
        'SetSecurityTokenServicePreferences',
        setSecurityTokenServicePreferences_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SetSecurityTokenServicePreferencesRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SimulateCustomPolicyRequest,
            $0.SimulatePolicyResponse>(
        'SimulateCustomPolicy',
        simulateCustomPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SimulateCustomPolicyRequest.fromBuffer(value),
        ($0.SimulatePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.SimulatePrincipalPolicyRequest,
            $0.SimulatePolicyResponse>(
        'SimulatePrincipalPolicy',
        simulatePrincipalPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.SimulatePrincipalPolicyRequest.fromBuffer(value),
        ($0.SimulatePolicyResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagInstanceProfileRequest, $1.Empty>(
        'TagInstanceProfile',
        tagInstanceProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagInstanceProfileRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagMFADeviceRequest, $1.Empty>(
        'TagMFADevice',
        tagMFADevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagMFADeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.TagOpenIDConnectProviderRequest, $1.Empty>(
            'TagOpenIDConnectProvider',
            tagOpenIDConnectProvider_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.TagOpenIDConnectProviderRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagPolicyRequest, $1.Empty>(
        'TagPolicy',
        tagPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagRoleRequest, $1.Empty>(
        'TagRole',
        tagRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagRoleRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagSAMLProviderRequest, $1.Empty>(
        'TagSAMLProvider',
        tagSAMLProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagSAMLProviderRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagServerCertificateRequest, $1.Empty>(
        'TagServerCertificate',
        tagServerCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.TagServerCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.TagUserRequest, $1.Empty>(
        'TagUser',
        tagUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.TagUserRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagInstanceProfileRequest, $1.Empty>(
        'UntagInstanceProfile',
        untagInstanceProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagInstanceProfileRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagMFADeviceRequest, $1.Empty>(
        'UntagMFADevice',
        untagMFADevice_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagMFADeviceRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UntagOpenIDConnectProviderRequest, $1.Empty>(
            'UntagOpenIDConnectProvider',
            untagOpenIDConnectProvider_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UntagOpenIDConnectProviderRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagPolicyRequest, $1.Empty>(
        'UntagPolicy',
        untagPolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagPolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagRoleRequest, $1.Empty>(
        'UntagRole',
        untagRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UntagRoleRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagSAMLProviderRequest, $1.Empty>(
        'UntagSAMLProvider',
        untagSAMLProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagSAMLProviderRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagServerCertificateRequest, $1.Empty>(
        'UntagServerCertificate',
        untagServerCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UntagServerCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UntagUserRequest, $1.Empty>(
        'UntagUser',
        untagUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UntagUserRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAccessKeyRequest, $1.Empty>(
        'UpdateAccessKey',
        updateAccessKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAccessKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateAccountPasswordPolicyRequest, $1.Empty>(
            'UpdateAccountPasswordPolicy',
            updateAccountPasswordPolicy_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateAccountPasswordPolicyRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateAssumeRolePolicyRequest, $1.Empty>(
        'UpdateAssumeRolePolicy',
        updateAssumeRolePolicy_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateAssumeRolePolicyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateDelegationRequestRequest, $1.Empty>(
        'UpdateDelegationRequest',
        updateDelegationRequest_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateDelegationRequestRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateGroupRequest, $1.Empty>(
        'UpdateGroup',
        updateGroup_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateGroupRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateLoginProfileRequest, $1.Empty>(
        'UpdateLoginProfile',
        updateLoginProfile_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateLoginProfileRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<
            $0.UpdateOpenIDConnectProviderThumbprintRequest, $1.Empty>(
        'UpdateOpenIDConnectProviderThumbprint',
        updateOpenIDConnectProviderThumbprint_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateOpenIDConnectProviderThumbprintRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRoleRequest, $0.UpdateRoleResponse>(
        'UpdateRole',
        updateRole_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateRoleRequest.fromBuffer(value),
        ($0.UpdateRoleResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateRoleDescriptionRequest,
            $0.UpdateRoleDescriptionResponse>(
        'UpdateRoleDescription',
        updateRoleDescription_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateRoleDescriptionRequest.fromBuffer(value),
        ($0.UpdateRoleDescriptionResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateSAMLProviderRequest,
            $0.UpdateSAMLProviderResponse>(
        'UpdateSAMLProvider',
        updateSAMLProvider_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateSAMLProviderRequest.fromBuffer(value),
        ($0.UpdateSAMLProviderResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateServerCertificateRequest, $1.Empty>(
        'UpdateServerCertificate',
        updateServerCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateServerCertificateRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateServiceSpecificCredentialRequest,
            $1.Empty>(
        'UpdateServiceSpecificCredential',
        updateServiceSpecificCredential_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateServiceSpecificCredentialRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod(
        $grpc.ServiceMethod<$0.UpdateSigningCertificateRequest, $1.Empty>(
            'UpdateSigningCertificate',
            updateSigningCertificate_Pre,
            false,
            false,
            ($core.List<$core.int> value) =>
                $0.UpdateSigningCertificateRequest.fromBuffer(value),
            ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateSSHPublicKeyRequest, $1.Empty>(
        'UpdateSSHPublicKey',
        updateSSHPublicKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UpdateSSHPublicKeyRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UpdateUserRequest, $1.Empty>(
        'UpdateUser',
        updateUser_Pre,
        false,
        false,
        ($core.List<$core.int> value) => $0.UpdateUserRequest.fromBuffer(value),
        ($1.Empty value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UploadServerCertificateRequest,
            $0.UploadServerCertificateResponse>(
        'UploadServerCertificate',
        uploadServerCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UploadServerCertificateRequest.fromBuffer(value),
        ($0.UploadServerCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UploadSigningCertificateRequest,
            $0.UploadSigningCertificateResponse>(
        'UploadSigningCertificate',
        uploadSigningCertificate_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UploadSigningCertificateRequest.fromBuffer(value),
        ($0.UploadSigningCertificateResponse value) => value.writeToBuffer()));
    $addMethod($grpc.ServiceMethod<$0.UploadSSHPublicKeyRequest,
            $0.UploadSSHPublicKeyResponse>(
        'UploadSSHPublicKey',
        uploadSSHPublicKey_Pre,
        false,
        false,
        ($core.List<$core.int> value) =>
            $0.UploadSSHPublicKeyRequest.fromBuffer(value),
        ($0.UploadSSHPublicKeyResponse value) => value.writeToBuffer()));
  }

  $async.Future<$1.Empty> acceptDelegationRequest_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AcceptDelegationRequestRequest> $request) async {
    return acceptDelegationRequest($call, await $request);
  }

  $async.Future<$1.Empty> acceptDelegationRequest(
      $grpc.ServiceCall call, $0.AcceptDelegationRequestRequest request);

  $async.Future<$1.Empty> addClientIDToOpenIDConnectProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AddClientIDToOpenIDConnectProviderRequest>
          $request) async {
    return addClientIDToOpenIDConnectProvider($call, await $request);
  }

  $async.Future<$1.Empty> addClientIDToOpenIDConnectProvider(
      $grpc.ServiceCall call,
      $0.AddClientIDToOpenIDConnectProviderRequest request);

  $async.Future<$1.Empty> addRoleToInstanceProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddRoleToInstanceProfileRequest> $request) async {
    return addRoleToInstanceProfile($call, await $request);
  }

  $async.Future<$1.Empty> addRoleToInstanceProfile(
      $grpc.ServiceCall call, $0.AddRoleToInstanceProfileRequest request);

  $async.Future<$1.Empty> addUserToGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AddUserToGroupRequest> $request) async {
    return addUserToGroup($call, await $request);
  }

  $async.Future<$1.Empty> addUserToGroup(
      $grpc.ServiceCall call, $0.AddUserToGroupRequest request);

  $async.Future<$1.Empty> associateDelegationRequest_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.AssociateDelegationRequestRequest> $request) async {
    return associateDelegationRequest($call, await $request);
  }

  $async.Future<$1.Empty> associateDelegationRequest(
      $grpc.ServiceCall call, $0.AssociateDelegationRequestRequest request);

  $async.Future<$1.Empty> attachGroupPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AttachGroupPolicyRequest> $request) async {
    return attachGroupPolicy($call, await $request);
  }

  $async.Future<$1.Empty> attachGroupPolicy(
      $grpc.ServiceCall call, $0.AttachGroupPolicyRequest request);

  $async.Future<$1.Empty> attachRolePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AttachRolePolicyRequest> $request) async {
    return attachRolePolicy($call, await $request);
  }

  $async.Future<$1.Empty> attachRolePolicy(
      $grpc.ServiceCall call, $0.AttachRolePolicyRequest request);

  $async.Future<$1.Empty> attachUserPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.AttachUserPolicyRequest> $request) async {
    return attachUserPolicy($call, await $request);
  }

  $async.Future<$1.Empty> attachUserPolicy(
      $grpc.ServiceCall call, $0.AttachUserPolicyRequest request);

  $async.Future<$1.Empty> changePassword_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ChangePasswordRequest> $request) async {
    return changePassword($call, await $request);
  }

  $async.Future<$1.Empty> changePassword(
      $grpc.ServiceCall call, $0.ChangePasswordRequest request);

  $async.Future<$0.CreateAccessKeyResponse> createAccessKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateAccessKeyRequest> $request) async {
    return createAccessKey($call, await $request);
  }

  $async.Future<$0.CreateAccessKeyResponse> createAccessKey(
      $grpc.ServiceCall call, $0.CreateAccessKeyRequest request);

  $async.Future<$1.Empty> createAccountAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateAccountAliasRequest> $request) async {
    return createAccountAlias($call, await $request);
  }

  $async.Future<$1.Empty> createAccountAlias(
      $grpc.ServiceCall call, $0.CreateAccountAliasRequest request);

  $async.Future<$0.CreateDelegationRequestResponse> createDelegationRequest_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateDelegationRequestRequest> $request) async {
    return createDelegationRequest($call, await $request);
  }

  $async.Future<$0.CreateDelegationRequestResponse> createDelegationRequest(
      $grpc.ServiceCall call, $0.CreateDelegationRequestRequest request);

  $async.Future<$0.CreateGroupResponse> createGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateGroupRequest> $request) async {
    return createGroup($call, await $request);
  }

  $async.Future<$0.CreateGroupResponse> createGroup(
      $grpc.ServiceCall call, $0.CreateGroupRequest request);

  $async.Future<$0.CreateInstanceProfileResponse> createInstanceProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateInstanceProfileRequest> $request) async {
    return createInstanceProfile($call, await $request);
  }

  $async.Future<$0.CreateInstanceProfileResponse> createInstanceProfile(
      $grpc.ServiceCall call, $0.CreateInstanceProfileRequest request);

  $async.Future<$0.CreateLoginProfileResponse> createLoginProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateLoginProfileRequest> $request) async {
    return createLoginProfile($call, await $request);
  }

  $async.Future<$0.CreateLoginProfileResponse> createLoginProfile(
      $grpc.ServiceCall call, $0.CreateLoginProfileRequest request);

  $async.Future<$0.CreateOpenIDConnectProviderResponse>
      createOpenIDConnectProvider_Pre($grpc.ServiceCall $call,
          $async.Future<$0.CreateOpenIDConnectProviderRequest> $request) async {
    return createOpenIDConnectProvider($call, await $request);
  }

  $async.Future<$0.CreateOpenIDConnectProviderResponse>
      createOpenIDConnectProvider($grpc.ServiceCall call,
          $0.CreateOpenIDConnectProviderRequest request);

  $async.Future<$0.CreatePolicyResponse> createPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreatePolicyRequest> $request) async {
    return createPolicy($call, await $request);
  }

  $async.Future<$0.CreatePolicyResponse> createPolicy(
      $grpc.ServiceCall call, $0.CreatePolicyRequest request);

  $async.Future<$0.CreatePolicyVersionResponse> createPolicyVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreatePolicyVersionRequest> $request) async {
    return createPolicyVersion($call, await $request);
  }

  $async.Future<$0.CreatePolicyVersionResponse> createPolicyVersion(
      $grpc.ServiceCall call, $0.CreatePolicyVersionRequest request);

  $async.Future<$0.CreateRoleResponse> createRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateRoleRequest> $request) async {
    return createRole($call, await $request);
  }

  $async.Future<$0.CreateRoleResponse> createRole(
      $grpc.ServiceCall call, $0.CreateRoleRequest request);

  $async.Future<$0.CreateSAMLProviderResponse> createSAMLProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateSAMLProviderRequest> $request) async {
    return createSAMLProvider($call, await $request);
  }

  $async.Future<$0.CreateSAMLProviderResponse> createSAMLProvider(
      $grpc.ServiceCall call, $0.CreateSAMLProviderRequest request);

  $async.Future<$0.CreateServiceLinkedRoleResponse> createServiceLinkedRole_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateServiceLinkedRoleRequest> $request) async {
    return createServiceLinkedRole($call, await $request);
  }

  $async.Future<$0.CreateServiceLinkedRoleResponse> createServiceLinkedRole(
      $grpc.ServiceCall call, $0.CreateServiceLinkedRoleRequest request);

  $async.Future<$0.CreateServiceSpecificCredentialResponse>
      createServiceSpecificCredential_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.CreateServiceSpecificCredentialRequest>
              $request) async {
    return createServiceSpecificCredential($call, await $request);
  }

  $async.Future<$0.CreateServiceSpecificCredentialResponse>
      createServiceSpecificCredential($grpc.ServiceCall call,
          $0.CreateServiceSpecificCredentialRequest request);

  $async.Future<$0.CreateUserResponse> createUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.CreateUserRequest> $request) async {
    return createUser($call, await $request);
  }

  $async.Future<$0.CreateUserResponse> createUser(
      $grpc.ServiceCall call, $0.CreateUserRequest request);

  $async.Future<$0.CreateVirtualMFADeviceResponse> createVirtualMFADevice_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.CreateVirtualMFADeviceRequest> $request) async {
    return createVirtualMFADevice($call, await $request);
  }

  $async.Future<$0.CreateVirtualMFADeviceResponse> createVirtualMFADevice(
      $grpc.ServiceCall call, $0.CreateVirtualMFADeviceRequest request);

  $async.Future<$1.Empty> deactivateMFADevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeactivateMFADeviceRequest> $request) async {
    return deactivateMFADevice($call, await $request);
  }

  $async.Future<$1.Empty> deactivateMFADevice(
      $grpc.ServiceCall call, $0.DeactivateMFADeviceRequest request);

  $async.Future<$1.Empty> deleteAccessKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAccessKeyRequest> $request) async {
    return deleteAccessKey($call, await $request);
  }

  $async.Future<$1.Empty> deleteAccessKey(
      $grpc.ServiceCall call, $0.DeleteAccessKeyRequest request);

  $async.Future<$1.Empty> deleteAccountAlias_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteAccountAliasRequest> $request) async {
    return deleteAccountAlias($call, await $request);
  }

  $async.Future<$1.Empty> deleteAccountAlias(
      $grpc.ServiceCall call, $0.DeleteAccountAliasRequest request);

  $async.Future<$1.Empty> deleteAccountPasswordPolicy_Pre(
      $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return deleteAccountPasswordPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteAccountPasswordPolicy(
      $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$1.Empty> deleteGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteGroupRequest> $request) async {
    return deleteGroup($call, await $request);
  }

  $async.Future<$1.Empty> deleteGroup(
      $grpc.ServiceCall call, $0.DeleteGroupRequest request);

  $async.Future<$1.Empty> deleteGroupPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteGroupPolicyRequest> $request) async {
    return deleteGroupPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteGroupPolicy(
      $grpc.ServiceCall call, $0.DeleteGroupPolicyRequest request);

  $async.Future<$1.Empty> deleteInstanceProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteInstanceProfileRequest> $request) async {
    return deleteInstanceProfile($call, await $request);
  }

  $async.Future<$1.Empty> deleteInstanceProfile(
      $grpc.ServiceCall call, $0.DeleteInstanceProfileRequest request);

  $async.Future<$1.Empty> deleteLoginProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteLoginProfileRequest> $request) async {
    return deleteLoginProfile($call, await $request);
  }

  $async.Future<$1.Empty> deleteLoginProfile(
      $grpc.ServiceCall call, $0.DeleteLoginProfileRequest request);

  $async.Future<$1.Empty> deleteOpenIDConnectProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteOpenIDConnectProviderRequest> $request) async {
    return deleteOpenIDConnectProvider($call, await $request);
  }

  $async.Future<$1.Empty> deleteOpenIDConnectProvider(
      $grpc.ServiceCall call, $0.DeleteOpenIDConnectProviderRequest request);

  $async.Future<$1.Empty> deletePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeletePolicyRequest> $request) async {
    return deletePolicy($call, await $request);
  }

  $async.Future<$1.Empty> deletePolicy(
      $grpc.ServiceCall call, $0.DeletePolicyRequest request);

  $async.Future<$1.Empty> deletePolicyVersion_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeletePolicyVersionRequest> $request) async {
    return deletePolicyVersion($call, await $request);
  }

  $async.Future<$1.Empty> deletePolicyVersion(
      $grpc.ServiceCall call, $0.DeletePolicyVersionRequest request);

  $async.Future<$1.Empty> deleteRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRoleRequest> $request) async {
    return deleteRole($call, await $request);
  }

  $async.Future<$1.Empty> deleteRole(
      $grpc.ServiceCall call, $0.DeleteRoleRequest request);

  $async.Future<$1.Empty> deleteRolePermissionsBoundary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteRolePermissionsBoundaryRequest> $request) async {
    return deleteRolePermissionsBoundary($call, await $request);
  }

  $async.Future<$1.Empty> deleteRolePermissionsBoundary(
      $grpc.ServiceCall call, $0.DeleteRolePermissionsBoundaryRequest request);

  $async.Future<$1.Empty> deleteRolePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteRolePolicyRequest> $request) async {
    return deleteRolePolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteRolePolicy(
      $grpc.ServiceCall call, $0.DeleteRolePolicyRequest request);

  $async.Future<$1.Empty> deleteSAMLProvider_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteSAMLProviderRequest> $request) async {
    return deleteSAMLProvider($call, await $request);
  }

  $async.Future<$1.Empty> deleteSAMLProvider(
      $grpc.ServiceCall call, $0.DeleteSAMLProviderRequest request);

  $async.Future<$1.Empty> deleteServerCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteServerCertificateRequest> $request) async {
    return deleteServerCertificate($call, await $request);
  }

  $async.Future<$1.Empty> deleteServerCertificate(
      $grpc.ServiceCall call, $0.DeleteServerCertificateRequest request);

  $async.Future<$0.DeleteServiceLinkedRoleResponse> deleteServiceLinkedRole_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteServiceLinkedRoleRequest> $request) async {
    return deleteServiceLinkedRole($call, await $request);
  }

  $async.Future<$0.DeleteServiceLinkedRoleResponse> deleteServiceLinkedRole(
      $grpc.ServiceCall call, $0.DeleteServiceLinkedRoleRequest request);

  $async.Future<$1.Empty> deleteServiceSpecificCredential_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteServiceSpecificCredentialRequest> $request) async {
    return deleteServiceSpecificCredential($call, await $request);
  }

  $async.Future<$1.Empty> deleteServiceSpecificCredential(
      $grpc.ServiceCall call,
      $0.DeleteServiceSpecificCredentialRequest request);

  $async.Future<$1.Empty> deleteSigningCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteSigningCertificateRequest> $request) async {
    return deleteSigningCertificate($call, await $request);
  }

  $async.Future<$1.Empty> deleteSigningCertificate(
      $grpc.ServiceCall call, $0.DeleteSigningCertificateRequest request);

  $async.Future<$1.Empty> deleteSSHPublicKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteSSHPublicKeyRequest> $request) async {
    return deleteSSHPublicKey($call, await $request);
  }

  $async.Future<$1.Empty> deleteSSHPublicKey(
      $grpc.ServiceCall call, $0.DeleteSSHPublicKeyRequest request);

  $async.Future<$1.Empty> deleteUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserRequest> $request) async {
    return deleteUser($call, await $request);
  }

  $async.Future<$1.Empty> deleteUser(
      $grpc.ServiceCall call, $0.DeleteUserRequest request);

  $async.Future<$1.Empty> deleteUserPermissionsBoundary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserPermissionsBoundaryRequest> $request) async {
    return deleteUserPermissionsBoundary($call, await $request);
  }

  $async.Future<$1.Empty> deleteUserPermissionsBoundary(
      $grpc.ServiceCall call, $0.DeleteUserPermissionsBoundaryRequest request);

  $async.Future<$1.Empty> deleteUserPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteUserPolicyRequest> $request) async {
    return deleteUserPolicy($call, await $request);
  }

  $async.Future<$1.Empty> deleteUserPolicy(
      $grpc.ServiceCall call, $0.DeleteUserPolicyRequest request);

  $async.Future<$1.Empty> deleteVirtualMFADevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DeleteVirtualMFADeviceRequest> $request) async {
    return deleteVirtualMFADevice($call, await $request);
  }

  $async.Future<$1.Empty> deleteVirtualMFADevice(
      $grpc.ServiceCall call, $0.DeleteVirtualMFADeviceRequest request);

  $async.Future<$1.Empty> detachGroupPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DetachGroupPolicyRequest> $request) async {
    return detachGroupPolicy($call, await $request);
  }

  $async.Future<$1.Empty> detachGroupPolicy(
      $grpc.ServiceCall call, $0.DetachGroupPolicyRequest request);

  $async.Future<$1.Empty> detachRolePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DetachRolePolicyRequest> $request) async {
    return detachRolePolicy($call, await $request);
  }

  $async.Future<$1.Empty> detachRolePolicy(
      $grpc.ServiceCall call, $0.DetachRolePolicyRequest request);

  $async.Future<$1.Empty> detachUserPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.DetachUserPolicyRequest> $request) async {
    return detachUserPolicy($call, await $request);
  }

  $async.Future<$1.Empty> detachUserPolicy(
      $grpc.ServiceCall call, $0.DetachUserPolicyRequest request);

  $async.Future<$0.DisableOrganizationsRootCredentialsManagementResponse>
      disableOrganizationsRootCredentialsManagement_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DisableOrganizationsRootCredentialsManagementRequest>
              $request) async {
    return disableOrganizationsRootCredentialsManagement($call, await $request);
  }

  $async.Future<$0.DisableOrganizationsRootCredentialsManagementResponse>
      disableOrganizationsRootCredentialsManagement($grpc.ServiceCall call,
          $0.DisableOrganizationsRootCredentialsManagementRequest request);

  $async.Future<$0.DisableOrganizationsRootSessionsResponse>
      disableOrganizationsRootSessions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.DisableOrganizationsRootSessionsRequest>
              $request) async {
    return disableOrganizationsRootSessions($call, await $request);
  }

  $async.Future<$0.DisableOrganizationsRootSessionsResponse>
      disableOrganizationsRootSessions($grpc.ServiceCall call,
          $0.DisableOrganizationsRootSessionsRequest request);

  $async.Future<$1.Empty> disableOutboundWebIdentityFederation_Pre(
      $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return disableOutboundWebIdentityFederation($call, await $request);
  }

  $async.Future<$1.Empty> disableOutboundWebIdentityFederation(
      $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$1.Empty> enableMFADevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.EnableMFADeviceRequest> $request) async {
    return enableMFADevice($call, await $request);
  }

  $async.Future<$1.Empty> enableMFADevice(
      $grpc.ServiceCall call, $0.EnableMFADeviceRequest request);

  $async.Future<$0.EnableOrganizationsRootCredentialsManagementResponse>
      enableOrganizationsRootCredentialsManagement_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.EnableOrganizationsRootCredentialsManagementRequest>
              $request) async {
    return enableOrganizationsRootCredentialsManagement($call, await $request);
  }

  $async.Future<$0.EnableOrganizationsRootCredentialsManagementResponse>
      enableOrganizationsRootCredentialsManagement($grpc.ServiceCall call,
          $0.EnableOrganizationsRootCredentialsManagementRequest request);

  $async.Future<$0.EnableOrganizationsRootSessionsResponse>
      enableOrganizationsRootSessions_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.EnableOrganizationsRootSessionsRequest>
              $request) async {
    return enableOrganizationsRootSessions($call, await $request);
  }

  $async.Future<$0.EnableOrganizationsRootSessionsResponse>
      enableOrganizationsRootSessions($grpc.ServiceCall call,
          $0.EnableOrganizationsRootSessionsRequest request);

  $async.Future<$0.EnableOutboundWebIdentityFederationResponse>
      enableOutboundWebIdentityFederation_Pre(
          $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return enableOutboundWebIdentityFederation($call, await $request);
  }

  $async.Future<$0.EnableOutboundWebIdentityFederationResponse>
      enableOutboundWebIdentityFederation(
          $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$0.GenerateCredentialReportResponse>
      generateCredentialReport_Pre(
          $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return generateCredentialReport($call, await $request);
  }

  $async.Future<$0.GenerateCredentialReportResponse> generateCredentialReport(
      $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$0.GenerateOrganizationsAccessReportResponse>
      generateOrganizationsAccessReport_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GenerateOrganizationsAccessReportRequest>
              $request) async {
    return generateOrganizationsAccessReport($call, await $request);
  }

  $async.Future<$0.GenerateOrganizationsAccessReportResponse>
      generateOrganizationsAccessReport($grpc.ServiceCall call,
          $0.GenerateOrganizationsAccessReportRequest request);

  $async.Future<$0.GenerateServiceLastAccessedDetailsResponse>
      generateServiceLastAccessedDetails_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GenerateServiceLastAccessedDetailsRequest>
              $request) async {
    return generateServiceLastAccessedDetails($call, await $request);
  }

  $async.Future<$0.GenerateServiceLastAccessedDetailsResponse>
      generateServiceLastAccessedDetails($grpc.ServiceCall call,
          $0.GenerateServiceLastAccessedDetailsRequest request);

  $async.Future<$0.GetAccessKeyLastUsedResponse> getAccessKeyLastUsed_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetAccessKeyLastUsedRequest> $request) async {
    return getAccessKeyLastUsed($call, await $request);
  }

  $async.Future<$0.GetAccessKeyLastUsedResponse> getAccessKeyLastUsed(
      $grpc.ServiceCall call, $0.GetAccessKeyLastUsedRequest request);

  $async.Future<$0.GetAccountAuthorizationDetailsResponse>
      getAccountAuthorizationDetails_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetAccountAuthorizationDetailsRequest>
              $request) async {
    return getAccountAuthorizationDetails($call, await $request);
  }

  $async.Future<$0.GetAccountAuthorizationDetailsResponse>
      getAccountAuthorizationDetails($grpc.ServiceCall call,
          $0.GetAccountAuthorizationDetailsRequest request);

  $async.Future<$0.GetAccountPasswordPolicyResponse>
      getAccountPasswordPolicy_Pre(
          $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return getAccountPasswordPolicy($call, await $request);
  }

  $async.Future<$0.GetAccountPasswordPolicyResponse> getAccountPasswordPolicy(
      $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$0.GetAccountSummaryResponse> getAccountSummary_Pre(
      $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return getAccountSummary($call, await $request);
  }

  $async.Future<$0.GetAccountSummaryResponse> getAccountSummary(
      $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$0.GetContextKeysForPolicyResponse>
      getContextKeysForCustomPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetContextKeysForCustomPolicyRequest>
              $request) async {
    return getContextKeysForCustomPolicy($call, await $request);
  }

  $async.Future<$0.GetContextKeysForPolicyResponse>
      getContextKeysForCustomPolicy($grpc.ServiceCall call,
          $0.GetContextKeysForCustomPolicyRequest request);

  $async.Future<$0.GetContextKeysForPolicyResponse>
      getContextKeysForPrincipalPolicy_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetContextKeysForPrincipalPolicyRequest>
              $request) async {
    return getContextKeysForPrincipalPolicy($call, await $request);
  }

  $async.Future<$0.GetContextKeysForPolicyResponse>
      getContextKeysForPrincipalPolicy($grpc.ServiceCall call,
          $0.GetContextKeysForPrincipalPolicyRequest request);

  $async.Future<$0.GetCredentialReportResponse> getCredentialReport_Pre(
      $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return getCredentialReport($call, await $request);
  }

  $async.Future<$0.GetCredentialReportResponse> getCredentialReport(
      $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$0.GetDelegationRequestResponse> getDelegationRequest_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetDelegationRequestRequest> $request) async {
    return getDelegationRequest($call, await $request);
  }

  $async.Future<$0.GetDelegationRequestResponse> getDelegationRequest(
      $grpc.ServiceCall call, $0.GetDelegationRequestRequest request);

  $async.Future<$0.GetGroupResponse> getGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetGroupRequest> $request) async {
    return getGroup($call, await $request);
  }

  $async.Future<$0.GetGroupResponse> getGroup(
      $grpc.ServiceCall call, $0.GetGroupRequest request);

  $async.Future<$0.GetGroupPolicyResponse> getGroupPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetGroupPolicyRequest> $request) async {
    return getGroupPolicy($call, await $request);
  }

  $async.Future<$0.GetGroupPolicyResponse> getGroupPolicy(
      $grpc.ServiceCall call, $0.GetGroupPolicyRequest request);

  $async.Future<$0.GetHumanReadableSummaryResponse> getHumanReadableSummary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetHumanReadableSummaryRequest> $request) async {
    return getHumanReadableSummary($call, await $request);
  }

  $async.Future<$0.GetHumanReadableSummaryResponse> getHumanReadableSummary(
      $grpc.ServiceCall call, $0.GetHumanReadableSummaryRequest request);

  $async.Future<$0.GetInstanceProfileResponse> getInstanceProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetInstanceProfileRequest> $request) async {
    return getInstanceProfile($call, await $request);
  }

  $async.Future<$0.GetInstanceProfileResponse> getInstanceProfile(
      $grpc.ServiceCall call, $0.GetInstanceProfileRequest request);

  $async.Future<$0.GetLoginProfileResponse> getLoginProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetLoginProfileRequest> $request) async {
    return getLoginProfile($call, await $request);
  }

  $async.Future<$0.GetLoginProfileResponse> getLoginProfile(
      $grpc.ServiceCall call, $0.GetLoginProfileRequest request);

  $async.Future<$0.GetMFADeviceResponse> getMFADevice_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetMFADeviceRequest> $request) async {
    return getMFADevice($call, await $request);
  }

  $async.Future<$0.GetMFADeviceResponse> getMFADevice(
      $grpc.ServiceCall call, $0.GetMFADeviceRequest request);

  $async.Future<$0.GetOpenIDConnectProviderResponse>
      getOpenIDConnectProvider_Pre($grpc.ServiceCall $call,
          $async.Future<$0.GetOpenIDConnectProviderRequest> $request) async {
    return getOpenIDConnectProvider($call, await $request);
  }

  $async.Future<$0.GetOpenIDConnectProviderResponse> getOpenIDConnectProvider(
      $grpc.ServiceCall call, $0.GetOpenIDConnectProviderRequest request);

  $async.Future<$0.GetOrganizationsAccessReportResponse>
      getOrganizationsAccessReport_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetOrganizationsAccessReportRequest>
              $request) async {
    return getOrganizationsAccessReport($call, await $request);
  }

  $async.Future<$0.GetOrganizationsAccessReportResponse>
      getOrganizationsAccessReport($grpc.ServiceCall call,
          $0.GetOrganizationsAccessReportRequest request);

  $async.Future<$0.GetOutboundWebIdentityFederationInfoResponse>
      getOutboundWebIdentityFederationInfo_Pre(
          $grpc.ServiceCall $call, $async.Future<$1.Empty> $request) async {
    return getOutboundWebIdentityFederationInfo($call, await $request);
  }

  $async.Future<$0.GetOutboundWebIdentityFederationInfoResponse>
      getOutboundWebIdentityFederationInfo(
          $grpc.ServiceCall call, $1.Empty request);

  $async.Future<$0.GetPolicyResponse> getPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetPolicyRequest> $request) async {
    return getPolicy($call, await $request);
  }

  $async.Future<$0.GetPolicyResponse> getPolicy(
      $grpc.ServiceCall call, $0.GetPolicyRequest request);

  $async.Future<$0.GetPolicyVersionResponse> getPolicyVersion_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetPolicyVersionRequest> $request) async {
    return getPolicyVersion($call, await $request);
  }

  $async.Future<$0.GetPolicyVersionResponse> getPolicyVersion(
      $grpc.ServiceCall call, $0.GetPolicyVersionRequest request);

  $async.Future<$0.GetRoleResponse> getRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetRoleRequest> $request) async {
    return getRole($call, await $request);
  }

  $async.Future<$0.GetRoleResponse> getRole(
      $grpc.ServiceCall call, $0.GetRoleRequest request);

  $async.Future<$0.GetRolePolicyResponse> getRolePolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetRolePolicyRequest> $request) async {
    return getRolePolicy($call, await $request);
  }

  $async.Future<$0.GetRolePolicyResponse> getRolePolicy(
      $grpc.ServiceCall call, $0.GetRolePolicyRequest request);

  $async.Future<$0.GetSAMLProviderResponse> getSAMLProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSAMLProviderRequest> $request) async {
    return getSAMLProvider($call, await $request);
  }

  $async.Future<$0.GetSAMLProviderResponse> getSAMLProvider(
      $grpc.ServiceCall call, $0.GetSAMLProviderRequest request);

  $async.Future<$0.GetServerCertificateResponse> getServerCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetServerCertificateRequest> $request) async {
    return getServerCertificate($call, await $request);
  }

  $async.Future<$0.GetServerCertificateResponse> getServerCertificate(
      $grpc.ServiceCall call, $0.GetServerCertificateRequest request);

  $async.Future<$0.GetServiceLastAccessedDetailsResponse>
      getServiceLastAccessedDetails_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetServiceLastAccessedDetailsRequest>
              $request) async {
    return getServiceLastAccessedDetails($call, await $request);
  }

  $async.Future<$0.GetServiceLastAccessedDetailsResponse>
      getServiceLastAccessedDetails($grpc.ServiceCall call,
          $0.GetServiceLastAccessedDetailsRequest request);

  $async.Future<$0.GetServiceLastAccessedDetailsWithEntitiesResponse>
      getServiceLastAccessedDetailsWithEntities_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetServiceLastAccessedDetailsWithEntitiesRequest>
              $request) async {
    return getServiceLastAccessedDetailsWithEntities($call, await $request);
  }

  $async.Future<$0.GetServiceLastAccessedDetailsWithEntitiesResponse>
      getServiceLastAccessedDetailsWithEntities($grpc.ServiceCall call,
          $0.GetServiceLastAccessedDetailsWithEntitiesRequest request);

  $async.Future<$0.GetServiceLinkedRoleDeletionStatusResponse>
      getServiceLinkedRoleDeletionStatus_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.GetServiceLinkedRoleDeletionStatusRequest>
              $request) async {
    return getServiceLinkedRoleDeletionStatus($call, await $request);
  }

  $async.Future<$0.GetServiceLinkedRoleDeletionStatusResponse>
      getServiceLinkedRoleDeletionStatus($grpc.ServiceCall call,
          $0.GetServiceLinkedRoleDeletionStatusRequest request);

  $async.Future<$0.GetSSHPublicKeyResponse> getSSHPublicKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetSSHPublicKeyRequest> $request) async {
    return getSSHPublicKey($call, await $request);
  }

  $async.Future<$0.GetSSHPublicKeyResponse> getSSHPublicKey(
      $grpc.ServiceCall call, $0.GetSSHPublicKeyRequest request);

  $async.Future<$0.GetUserResponse> getUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.GetUserRequest> $request) async {
    return getUser($call, await $request);
  }

  $async.Future<$0.GetUserResponse> getUser(
      $grpc.ServiceCall call, $0.GetUserRequest request);

  $async.Future<$0.GetUserPolicyResponse> getUserPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.GetUserPolicyRequest> $request) async {
    return getUserPolicy($call, await $request);
  }

  $async.Future<$0.GetUserPolicyResponse> getUserPolicy(
      $grpc.ServiceCall call, $0.GetUserPolicyRequest request);

  $async.Future<$0.ListAccessKeysResponse> listAccessKeys_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAccessKeysRequest> $request) async {
    return listAccessKeys($call, await $request);
  }

  $async.Future<$0.ListAccessKeysResponse> listAccessKeys(
      $grpc.ServiceCall call, $0.ListAccessKeysRequest request);

  $async.Future<$0.ListAccountAliasesResponse> listAccountAliases_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListAccountAliasesRequest> $request) async {
    return listAccountAliases($call, await $request);
  }

  $async.Future<$0.ListAccountAliasesResponse> listAccountAliases(
      $grpc.ServiceCall call, $0.ListAccountAliasesRequest request);

  $async.Future<$0.ListAttachedGroupPoliciesResponse>
      listAttachedGroupPolicies_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListAttachedGroupPoliciesRequest> $request) async {
    return listAttachedGroupPolicies($call, await $request);
  }

  $async.Future<$0.ListAttachedGroupPoliciesResponse> listAttachedGroupPolicies(
      $grpc.ServiceCall call, $0.ListAttachedGroupPoliciesRequest request);

  $async.Future<$0.ListAttachedRolePoliciesResponse>
      listAttachedRolePolicies_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListAttachedRolePoliciesRequest> $request) async {
    return listAttachedRolePolicies($call, await $request);
  }

  $async.Future<$0.ListAttachedRolePoliciesResponse> listAttachedRolePolicies(
      $grpc.ServiceCall call, $0.ListAttachedRolePoliciesRequest request);

  $async.Future<$0.ListAttachedUserPoliciesResponse>
      listAttachedUserPolicies_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListAttachedUserPoliciesRequest> $request) async {
    return listAttachedUserPolicies($call, await $request);
  }

  $async.Future<$0.ListAttachedUserPoliciesResponse> listAttachedUserPolicies(
      $grpc.ServiceCall call, $0.ListAttachedUserPoliciesRequest request);

  $async.Future<$0.ListDelegationRequestsResponse> listDelegationRequests_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListDelegationRequestsRequest> $request) async {
    return listDelegationRequests($call, await $request);
  }

  $async.Future<$0.ListDelegationRequestsResponse> listDelegationRequests(
      $grpc.ServiceCall call, $0.ListDelegationRequestsRequest request);

  $async.Future<$0.ListEntitiesForPolicyResponse> listEntitiesForPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListEntitiesForPolicyRequest> $request) async {
    return listEntitiesForPolicy($call, await $request);
  }

  $async.Future<$0.ListEntitiesForPolicyResponse> listEntitiesForPolicy(
      $grpc.ServiceCall call, $0.ListEntitiesForPolicyRequest request);

  $async.Future<$0.ListGroupPoliciesResponse> listGroupPolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListGroupPoliciesRequest> $request) async {
    return listGroupPolicies($call, await $request);
  }

  $async.Future<$0.ListGroupPoliciesResponse> listGroupPolicies(
      $grpc.ServiceCall call, $0.ListGroupPoliciesRequest request);

  $async.Future<$0.ListGroupsResponse> listGroups_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListGroupsRequest> $request) async {
    return listGroups($call, await $request);
  }

  $async.Future<$0.ListGroupsResponse> listGroups(
      $grpc.ServiceCall call, $0.ListGroupsRequest request);

  $async.Future<$0.ListGroupsForUserResponse> listGroupsForUser_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListGroupsForUserRequest> $request) async {
    return listGroupsForUser($call, await $request);
  }

  $async.Future<$0.ListGroupsForUserResponse> listGroupsForUser(
      $grpc.ServiceCall call, $0.ListGroupsForUserRequest request);

  $async.Future<$0.ListInstanceProfilesResponse> listInstanceProfiles_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListInstanceProfilesRequest> $request) async {
    return listInstanceProfiles($call, await $request);
  }

  $async.Future<$0.ListInstanceProfilesResponse> listInstanceProfiles(
      $grpc.ServiceCall call, $0.ListInstanceProfilesRequest request);

  $async.Future<$0.ListInstanceProfilesForRoleResponse>
      listInstanceProfilesForRole_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListInstanceProfilesForRoleRequest> $request) async {
    return listInstanceProfilesForRole($call, await $request);
  }

  $async.Future<$0.ListInstanceProfilesForRoleResponse>
      listInstanceProfilesForRole($grpc.ServiceCall call,
          $0.ListInstanceProfilesForRoleRequest request);

  $async.Future<$0.ListInstanceProfileTagsResponse> listInstanceProfileTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListInstanceProfileTagsRequest> $request) async {
    return listInstanceProfileTags($call, await $request);
  }

  $async.Future<$0.ListInstanceProfileTagsResponse> listInstanceProfileTags(
      $grpc.ServiceCall call, $0.ListInstanceProfileTagsRequest request);

  $async.Future<$0.ListMFADevicesResponse> listMFADevices_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListMFADevicesRequest> $request) async {
    return listMFADevices($call, await $request);
  }

  $async.Future<$0.ListMFADevicesResponse> listMFADevices(
      $grpc.ServiceCall call, $0.ListMFADevicesRequest request);

  $async.Future<$0.ListMFADeviceTagsResponse> listMFADeviceTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListMFADeviceTagsRequest> $request) async {
    return listMFADeviceTags($call, await $request);
  }

  $async.Future<$0.ListMFADeviceTagsResponse> listMFADeviceTags(
      $grpc.ServiceCall call, $0.ListMFADeviceTagsRequest request);

  $async.Future<$0.ListOpenIDConnectProvidersResponse>
      listOpenIDConnectProviders_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListOpenIDConnectProvidersRequest> $request) async {
    return listOpenIDConnectProviders($call, await $request);
  }

  $async.Future<$0.ListOpenIDConnectProvidersResponse>
      listOpenIDConnectProviders(
          $grpc.ServiceCall call, $0.ListOpenIDConnectProvidersRequest request);

  $async.Future<$0.ListOpenIDConnectProviderTagsResponse>
      listOpenIDConnectProviderTags_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListOpenIDConnectProviderTagsRequest>
              $request) async {
    return listOpenIDConnectProviderTags($call, await $request);
  }

  $async.Future<$0.ListOpenIDConnectProviderTagsResponse>
      listOpenIDConnectProviderTags($grpc.ServiceCall call,
          $0.ListOpenIDConnectProviderTagsRequest request);

  $async.Future<$0.ListOrganizationsFeaturesResponse>
      listOrganizationsFeatures_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListOrganizationsFeaturesRequest> $request) async {
    return listOrganizationsFeatures($call, await $request);
  }

  $async.Future<$0.ListOrganizationsFeaturesResponse> listOrganizationsFeatures(
      $grpc.ServiceCall call, $0.ListOrganizationsFeaturesRequest request);

  $async.Future<$0.ListPoliciesResponse> listPolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPoliciesRequest> $request) async {
    return listPolicies($call, await $request);
  }

  $async.Future<$0.ListPoliciesResponse> listPolicies(
      $grpc.ServiceCall call, $0.ListPoliciesRequest request);

  $async.Future<$0.ListPoliciesGrantingServiceAccessResponse>
      listPoliciesGrantingServiceAccess_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListPoliciesGrantingServiceAccessRequest>
              $request) async {
    return listPoliciesGrantingServiceAccess($call, await $request);
  }

  $async.Future<$0.ListPoliciesGrantingServiceAccessResponse>
      listPoliciesGrantingServiceAccess($grpc.ServiceCall call,
          $0.ListPoliciesGrantingServiceAccessRequest request);

  $async.Future<$0.ListPolicyTagsResponse> listPolicyTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPolicyTagsRequest> $request) async {
    return listPolicyTags($call, await $request);
  }

  $async.Future<$0.ListPolicyTagsResponse> listPolicyTags(
      $grpc.ServiceCall call, $0.ListPolicyTagsRequest request);

  $async.Future<$0.ListPolicyVersionsResponse> listPolicyVersions_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListPolicyVersionsRequest> $request) async {
    return listPolicyVersions($call, await $request);
  }

  $async.Future<$0.ListPolicyVersionsResponse> listPolicyVersions(
      $grpc.ServiceCall call, $0.ListPolicyVersionsRequest request);

  $async.Future<$0.ListRolePoliciesResponse> listRolePolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRolePoliciesRequest> $request) async {
    return listRolePolicies($call, await $request);
  }

  $async.Future<$0.ListRolePoliciesResponse> listRolePolicies(
      $grpc.ServiceCall call, $0.ListRolePoliciesRequest request);

  $async.Future<$0.ListRolesResponse> listRoles_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListRolesRequest> $request) async {
    return listRoles($call, await $request);
  }

  $async.Future<$0.ListRolesResponse> listRoles(
      $grpc.ServiceCall call, $0.ListRolesRequest request);

  $async.Future<$0.ListRoleTagsResponse> listRoleTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListRoleTagsRequest> $request) async {
    return listRoleTags($call, await $request);
  }

  $async.Future<$0.ListRoleTagsResponse> listRoleTags(
      $grpc.ServiceCall call, $0.ListRoleTagsRequest request);

  $async.Future<$0.ListSAMLProvidersResponse> listSAMLProviders_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSAMLProvidersRequest> $request) async {
    return listSAMLProviders($call, await $request);
  }

  $async.Future<$0.ListSAMLProvidersResponse> listSAMLProviders(
      $grpc.ServiceCall call, $0.ListSAMLProvidersRequest request);

  $async.Future<$0.ListSAMLProviderTagsResponse> listSAMLProviderTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSAMLProviderTagsRequest> $request) async {
    return listSAMLProviderTags($call, await $request);
  }

  $async.Future<$0.ListSAMLProviderTagsResponse> listSAMLProviderTags(
      $grpc.ServiceCall call, $0.ListSAMLProviderTagsRequest request);

  $async.Future<$0.ListServerCertificatesResponse> listServerCertificates_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListServerCertificatesRequest> $request) async {
    return listServerCertificates($call, await $request);
  }

  $async.Future<$0.ListServerCertificatesResponse> listServerCertificates(
      $grpc.ServiceCall call, $0.ListServerCertificatesRequest request);

  $async.Future<$0.ListServerCertificateTagsResponse>
      listServerCertificateTags_Pre($grpc.ServiceCall $call,
          $async.Future<$0.ListServerCertificateTagsRequest> $request) async {
    return listServerCertificateTags($call, await $request);
  }

  $async.Future<$0.ListServerCertificateTagsResponse> listServerCertificateTags(
      $grpc.ServiceCall call, $0.ListServerCertificateTagsRequest request);

  $async.Future<$0.ListServiceSpecificCredentialsResponse>
      listServiceSpecificCredentials_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ListServiceSpecificCredentialsRequest>
              $request) async {
    return listServiceSpecificCredentials($call, await $request);
  }

  $async.Future<$0.ListServiceSpecificCredentialsResponse>
      listServiceSpecificCredentials($grpc.ServiceCall call,
          $0.ListServiceSpecificCredentialsRequest request);

  $async.Future<$0.ListSigningCertificatesResponse> listSigningCertificates_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSigningCertificatesRequest> $request) async {
    return listSigningCertificates($call, await $request);
  }

  $async.Future<$0.ListSigningCertificatesResponse> listSigningCertificates(
      $grpc.ServiceCall call, $0.ListSigningCertificatesRequest request);

  $async.Future<$0.ListSSHPublicKeysResponse> listSSHPublicKeys_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListSSHPublicKeysRequest> $request) async {
    return listSSHPublicKeys($call, await $request);
  }

  $async.Future<$0.ListSSHPublicKeysResponse> listSSHPublicKeys(
      $grpc.ServiceCall call, $0.ListSSHPublicKeysRequest request);

  $async.Future<$0.ListUserPoliciesResponse> listUserPolicies_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListUserPoliciesRequest> $request) async {
    return listUserPolicies($call, await $request);
  }

  $async.Future<$0.ListUserPoliciesResponse> listUserPolicies(
      $grpc.ServiceCall call, $0.ListUserPoliciesRequest request);

  $async.Future<$0.ListUsersResponse> listUsers_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ListUsersRequest> $request) async {
    return listUsers($call, await $request);
  }

  $async.Future<$0.ListUsersResponse> listUsers(
      $grpc.ServiceCall call, $0.ListUsersRequest request);

  $async.Future<$0.ListUserTagsResponse> listUserTags_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListUserTagsRequest> $request) async {
    return listUserTags($call, await $request);
  }

  $async.Future<$0.ListUserTagsResponse> listUserTags(
      $grpc.ServiceCall call, $0.ListUserTagsRequest request);

  $async.Future<$0.ListVirtualMFADevicesResponse> listVirtualMFADevices_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.ListVirtualMFADevicesRequest> $request) async {
    return listVirtualMFADevices($call, await $request);
  }

  $async.Future<$0.ListVirtualMFADevicesResponse> listVirtualMFADevices(
      $grpc.ServiceCall call, $0.ListVirtualMFADevicesRequest request);

  $async.Future<$1.Empty> putGroupPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutGroupPolicyRequest> $request) async {
    return putGroupPolicy($call, await $request);
  }

  $async.Future<$1.Empty> putGroupPolicy(
      $grpc.ServiceCall call, $0.PutGroupPolicyRequest request);

  $async.Future<$1.Empty> putRolePermissionsBoundary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutRolePermissionsBoundaryRequest> $request) async {
    return putRolePermissionsBoundary($call, await $request);
  }

  $async.Future<$1.Empty> putRolePermissionsBoundary(
      $grpc.ServiceCall call, $0.PutRolePermissionsBoundaryRequest request);

  $async.Future<$1.Empty> putRolePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutRolePolicyRequest> $request) async {
    return putRolePolicy($call, await $request);
  }

  $async.Future<$1.Empty> putRolePolicy(
      $grpc.ServiceCall call, $0.PutRolePolicyRequest request);

  $async.Future<$1.Empty> putUserPermissionsBoundary_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.PutUserPermissionsBoundaryRequest> $request) async {
    return putUserPermissionsBoundary($call, await $request);
  }

  $async.Future<$1.Empty> putUserPermissionsBoundary(
      $grpc.ServiceCall call, $0.PutUserPermissionsBoundaryRequest request);

  $async.Future<$1.Empty> putUserPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.PutUserPolicyRequest> $request) async {
    return putUserPolicy($call, await $request);
  }

  $async.Future<$1.Empty> putUserPolicy(
      $grpc.ServiceCall call, $0.PutUserPolicyRequest request);

  $async.Future<$1.Empty> rejectDelegationRequest_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RejectDelegationRequestRequest> $request) async {
    return rejectDelegationRequest($call, await $request);
  }

  $async.Future<$1.Empty> rejectDelegationRequest(
      $grpc.ServiceCall call, $0.RejectDelegationRequestRequest request);

  $async.Future<$1.Empty> removeClientIDFromOpenIDConnectProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RemoveClientIDFromOpenIDConnectProviderRequest>
          $request) async {
    return removeClientIDFromOpenIDConnectProvider($call, await $request);
  }

  $async.Future<$1.Empty> removeClientIDFromOpenIDConnectProvider(
      $grpc.ServiceCall call,
      $0.RemoveClientIDFromOpenIDConnectProviderRequest request);

  $async.Future<$1.Empty> removeRoleFromInstanceProfile_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.RemoveRoleFromInstanceProfileRequest> $request) async {
    return removeRoleFromInstanceProfile($call, await $request);
  }

  $async.Future<$1.Empty> removeRoleFromInstanceProfile(
      $grpc.ServiceCall call, $0.RemoveRoleFromInstanceProfileRequest request);

  $async.Future<$1.Empty> removeUserFromGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.RemoveUserFromGroupRequest> $request) async {
    return removeUserFromGroup($call, await $request);
  }

  $async.Future<$1.Empty> removeUserFromGroup(
      $grpc.ServiceCall call, $0.RemoveUserFromGroupRequest request);

  $async.Future<$0.ResetServiceSpecificCredentialResponse>
      resetServiceSpecificCredential_Pre(
          $grpc.ServiceCall $call,
          $async.Future<$0.ResetServiceSpecificCredentialRequest>
              $request) async {
    return resetServiceSpecificCredential($call, await $request);
  }

  $async.Future<$0.ResetServiceSpecificCredentialResponse>
      resetServiceSpecificCredential($grpc.ServiceCall call,
          $0.ResetServiceSpecificCredentialRequest request);

  $async.Future<$1.Empty> resyncMFADevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.ResyncMFADeviceRequest> $request) async {
    return resyncMFADevice($call, await $request);
  }

  $async.Future<$1.Empty> resyncMFADevice(
      $grpc.ServiceCall call, $0.ResyncMFADeviceRequest request);

  $async.Future<$1.Empty> sendDelegationToken_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SendDelegationTokenRequest> $request) async {
    return sendDelegationToken($call, await $request);
  }

  $async.Future<$1.Empty> sendDelegationToken(
      $grpc.ServiceCall call, $0.SendDelegationTokenRequest request);

  $async.Future<$1.Empty> setDefaultPolicyVersion_Pre($grpc.ServiceCall $call,
      $async.Future<$0.SetDefaultPolicyVersionRequest> $request) async {
    return setDefaultPolicyVersion($call, await $request);
  }

  $async.Future<$1.Empty> setDefaultPolicyVersion(
      $grpc.ServiceCall call, $0.SetDefaultPolicyVersionRequest request);

  $async.Future<$1.Empty> setSecurityTokenServicePreferences_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SetSecurityTokenServicePreferencesRequest>
          $request) async {
    return setSecurityTokenServicePreferences($call, await $request);
  }

  $async.Future<$1.Empty> setSecurityTokenServicePreferences(
      $grpc.ServiceCall call,
      $0.SetSecurityTokenServicePreferencesRequest request);

  $async.Future<$0.SimulatePolicyResponse> simulateCustomPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SimulateCustomPolicyRequest> $request) async {
    return simulateCustomPolicy($call, await $request);
  }

  $async.Future<$0.SimulatePolicyResponse> simulateCustomPolicy(
      $grpc.ServiceCall call, $0.SimulateCustomPolicyRequest request);

  $async.Future<$0.SimulatePolicyResponse> simulatePrincipalPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.SimulatePrincipalPolicyRequest> $request) async {
    return simulatePrincipalPolicy($call, await $request);
  }

  $async.Future<$0.SimulatePolicyResponse> simulatePrincipalPolicy(
      $grpc.ServiceCall call, $0.SimulatePrincipalPolicyRequest request);

  $async.Future<$1.Empty> tagInstanceProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagInstanceProfileRequest> $request) async {
    return tagInstanceProfile($call, await $request);
  }

  $async.Future<$1.Empty> tagInstanceProfile(
      $grpc.ServiceCall call, $0.TagInstanceProfileRequest request);

  $async.Future<$1.Empty> tagMFADevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagMFADeviceRequest> $request) async {
    return tagMFADevice($call, await $request);
  }

  $async.Future<$1.Empty> tagMFADevice(
      $grpc.ServiceCall call, $0.TagMFADeviceRequest request);

  $async.Future<$1.Empty> tagOpenIDConnectProvider_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagOpenIDConnectProviderRequest> $request) async {
    return tagOpenIDConnectProvider($call, await $request);
  }

  $async.Future<$1.Empty> tagOpenIDConnectProvider(
      $grpc.ServiceCall call, $0.TagOpenIDConnectProviderRequest request);

  $async.Future<$1.Empty> tagPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagPolicyRequest> $request) async {
    return tagPolicy($call, await $request);
  }

  $async.Future<$1.Empty> tagPolicy(
      $grpc.ServiceCall call, $0.TagPolicyRequest request);

  $async.Future<$1.Empty> tagRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagRoleRequest> $request) async {
    return tagRole($call, await $request);
  }

  $async.Future<$1.Empty> tagRole(
      $grpc.ServiceCall call, $0.TagRoleRequest request);

  $async.Future<$1.Empty> tagSAMLProvider_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagSAMLProviderRequest> $request) async {
    return tagSAMLProvider($call, await $request);
  }

  $async.Future<$1.Empty> tagSAMLProvider(
      $grpc.ServiceCall call, $0.TagSAMLProviderRequest request);

  $async.Future<$1.Empty> tagServerCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagServerCertificateRequest> $request) async {
    return tagServerCertificate($call, await $request);
  }

  $async.Future<$1.Empty> tagServerCertificate(
      $grpc.ServiceCall call, $0.TagServerCertificateRequest request);

  $async.Future<$1.Empty> tagUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.TagUserRequest> $request) async {
    return tagUser($call, await $request);
  }

  $async.Future<$1.Empty> tagUser(
      $grpc.ServiceCall call, $0.TagUserRequest request);

  $async.Future<$1.Empty> untagInstanceProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagInstanceProfileRequest> $request) async {
    return untagInstanceProfile($call, await $request);
  }

  $async.Future<$1.Empty> untagInstanceProfile(
      $grpc.ServiceCall call, $0.UntagInstanceProfileRequest request);

  $async.Future<$1.Empty> untagMFADevice_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagMFADeviceRequest> $request) async {
    return untagMFADevice($call, await $request);
  }

  $async.Future<$1.Empty> untagMFADevice(
      $grpc.ServiceCall call, $0.UntagMFADeviceRequest request);

  $async.Future<$1.Empty> untagOpenIDConnectProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UntagOpenIDConnectProviderRequest> $request) async {
    return untagOpenIDConnectProvider($call, await $request);
  }

  $async.Future<$1.Empty> untagOpenIDConnectProvider(
      $grpc.ServiceCall call, $0.UntagOpenIDConnectProviderRequest request);

  $async.Future<$1.Empty> untagPolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagPolicyRequest> $request) async {
    return untagPolicy($call, await $request);
  }

  $async.Future<$1.Empty> untagPolicy(
      $grpc.ServiceCall call, $0.UntagPolicyRequest request);

  $async.Future<$1.Empty> untagRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagRoleRequest> $request) async {
    return untagRole($call, await $request);
  }

  $async.Future<$1.Empty> untagRole(
      $grpc.ServiceCall call, $0.UntagRoleRequest request);

  $async.Future<$1.Empty> untagSAMLProvider_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagSAMLProviderRequest> $request) async {
    return untagSAMLProvider($call, await $request);
  }

  $async.Future<$1.Empty> untagSAMLProvider(
      $grpc.ServiceCall call, $0.UntagSAMLProviderRequest request);

  $async.Future<$1.Empty> untagServerCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagServerCertificateRequest> $request) async {
    return untagServerCertificate($call, await $request);
  }

  $async.Future<$1.Empty> untagServerCertificate(
      $grpc.ServiceCall call, $0.UntagServerCertificateRequest request);

  $async.Future<$1.Empty> untagUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UntagUserRequest> $request) async {
    return untagUser($call, await $request);
  }

  $async.Future<$1.Empty> untagUser(
      $grpc.ServiceCall call, $0.UntagUserRequest request);

  $async.Future<$1.Empty> updateAccessKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAccessKeyRequest> $request) async {
    return updateAccessKey($call, await $request);
  }

  $async.Future<$1.Empty> updateAccessKey(
      $grpc.ServiceCall call, $0.UpdateAccessKeyRequest request);

  $async.Future<$1.Empty> updateAccountPasswordPolicy_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateAccountPasswordPolicyRequest> $request) async {
    return updateAccountPasswordPolicy($call, await $request);
  }

  $async.Future<$1.Empty> updateAccountPasswordPolicy(
      $grpc.ServiceCall call, $0.UpdateAccountPasswordPolicyRequest request);

  $async.Future<$1.Empty> updateAssumeRolePolicy_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateAssumeRolePolicyRequest> $request) async {
    return updateAssumeRolePolicy($call, await $request);
  }

  $async.Future<$1.Empty> updateAssumeRolePolicy(
      $grpc.ServiceCall call, $0.UpdateAssumeRolePolicyRequest request);

  $async.Future<$1.Empty> updateDelegationRequest_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateDelegationRequestRequest> $request) async {
    return updateDelegationRequest($call, await $request);
  }

  $async.Future<$1.Empty> updateDelegationRequest(
      $grpc.ServiceCall call, $0.UpdateDelegationRequestRequest request);

  $async.Future<$1.Empty> updateGroup_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateGroupRequest> $request) async {
    return updateGroup($call, await $request);
  }

  $async.Future<$1.Empty> updateGroup(
      $grpc.ServiceCall call, $0.UpdateGroupRequest request);

  $async.Future<$1.Empty> updateLoginProfile_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateLoginProfileRequest> $request) async {
    return updateLoginProfile($call, await $request);
  }

  $async.Future<$1.Empty> updateLoginProfile(
      $grpc.ServiceCall call, $0.UpdateLoginProfileRequest request);

  $async.Future<$1.Empty> updateOpenIDConnectProviderThumbprint_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateOpenIDConnectProviderThumbprintRequest>
          $request) async {
    return updateOpenIDConnectProviderThumbprint($call, await $request);
  }

  $async.Future<$1.Empty> updateOpenIDConnectProviderThumbprint(
      $grpc.ServiceCall call,
      $0.UpdateOpenIDConnectProviderThumbprintRequest request);

  $async.Future<$0.UpdateRoleResponse> updateRole_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateRoleRequest> $request) async {
    return updateRole($call, await $request);
  }

  $async.Future<$0.UpdateRoleResponse> updateRole(
      $grpc.ServiceCall call, $0.UpdateRoleRequest request);

  $async.Future<$0.UpdateRoleDescriptionResponse> updateRoleDescription_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateRoleDescriptionRequest> $request) async {
    return updateRoleDescription($call, await $request);
  }

  $async.Future<$0.UpdateRoleDescriptionResponse> updateRoleDescription(
      $grpc.ServiceCall call, $0.UpdateRoleDescriptionRequest request);

  $async.Future<$0.UpdateSAMLProviderResponse> updateSAMLProvider_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateSAMLProviderRequest> $request) async {
    return updateSAMLProvider($call, await $request);
  }

  $async.Future<$0.UpdateSAMLProviderResponse> updateSAMLProvider(
      $grpc.ServiceCall call, $0.UpdateSAMLProviderRequest request);

  $async.Future<$1.Empty> updateServerCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateServerCertificateRequest> $request) async {
    return updateServerCertificate($call, await $request);
  }

  $async.Future<$1.Empty> updateServerCertificate(
      $grpc.ServiceCall call, $0.UpdateServerCertificateRequest request);

  $async.Future<$1.Empty> updateServiceSpecificCredential_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UpdateServiceSpecificCredentialRequest> $request) async {
    return updateServiceSpecificCredential($call, await $request);
  }

  $async.Future<$1.Empty> updateServiceSpecificCredential(
      $grpc.ServiceCall call,
      $0.UpdateServiceSpecificCredentialRequest request);

  $async.Future<$1.Empty> updateSigningCertificate_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateSigningCertificateRequest> $request) async {
    return updateSigningCertificate($call, await $request);
  }

  $async.Future<$1.Empty> updateSigningCertificate(
      $grpc.ServiceCall call, $0.UpdateSigningCertificateRequest request);

  $async.Future<$1.Empty> updateSSHPublicKey_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateSSHPublicKeyRequest> $request) async {
    return updateSSHPublicKey($call, await $request);
  }

  $async.Future<$1.Empty> updateSSHPublicKey(
      $grpc.ServiceCall call, $0.UpdateSSHPublicKeyRequest request);

  $async.Future<$1.Empty> updateUser_Pre($grpc.ServiceCall $call,
      $async.Future<$0.UpdateUserRequest> $request) async {
    return updateUser($call, await $request);
  }

  $async.Future<$1.Empty> updateUser(
      $grpc.ServiceCall call, $0.UpdateUserRequest request);

  $async.Future<$0.UploadServerCertificateResponse> uploadServerCertificate_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UploadServerCertificateRequest> $request) async {
    return uploadServerCertificate($call, await $request);
  }

  $async.Future<$0.UploadServerCertificateResponse> uploadServerCertificate(
      $grpc.ServiceCall call, $0.UploadServerCertificateRequest request);

  $async.Future<$0.UploadSigningCertificateResponse>
      uploadSigningCertificate_Pre($grpc.ServiceCall $call,
          $async.Future<$0.UploadSigningCertificateRequest> $request) async {
    return uploadSigningCertificate($call, await $request);
  }

  $async.Future<$0.UploadSigningCertificateResponse> uploadSigningCertificate(
      $grpc.ServiceCall call, $0.UploadSigningCertificateRequest request);

  $async.Future<$0.UploadSSHPublicKeyResponse> uploadSSHPublicKey_Pre(
      $grpc.ServiceCall $call,
      $async.Future<$0.UploadSSHPublicKeyRequest> $request) async {
    return uploadSSHPublicKey($call, await $request);
  }

  $async.Future<$0.UploadSSHPublicKeyResponse> uploadSSHPublicKey(
      $grpc.ServiceCall call, $0.UploadSSHPublicKeyRequest request);
}
