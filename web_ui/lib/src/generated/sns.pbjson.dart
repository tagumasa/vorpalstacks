// This is a generated file - do not edit.
//
// Generated from sns.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports
// ignore_for_file: unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use languageCodeStringDescriptor instead')
const LanguageCodeString$json = {
  '1': 'LanguageCodeString',
  '2': [
    {'1': 'LANGUAGE_CODE_STRING_ZH_CN', '2': 0},
    {'1': 'LANGUAGE_CODE_STRING_EN_GB', '2': 1},
    {'1': 'LANGUAGE_CODE_STRING_IT_IT', '2': 2},
    {'1': 'LANGUAGE_CODE_STRING_JP_JP', '2': 3},
    {'1': 'LANGUAGE_CODE_STRING_ES_ES', '2': 4},
    {'1': 'LANGUAGE_CODE_STRING_DE_DE', '2': 5},
    {'1': 'LANGUAGE_CODE_STRING_KR_KR', '2': 6},
    {'1': 'LANGUAGE_CODE_STRING_FR_FR', '2': 7},
    {'1': 'LANGUAGE_CODE_STRING_EN_US', '2': 8},
    {'1': 'LANGUAGE_CODE_STRING_PT_BR', '2': 9},
    {'1': 'LANGUAGE_CODE_STRING_ZH_TW', '2': 10},
    {'1': 'LANGUAGE_CODE_STRING_ES_419', '2': 11},
    {'1': 'LANGUAGE_CODE_STRING_FR_CA', '2': 12},
  ],
};

/// Descriptor for `LanguageCodeString`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List languageCodeStringDescriptor = $convert.base64Decode(
    'ChJMYW5ndWFnZUNvZGVTdHJpbmcSHgoaTEFOR1VBR0VfQ09ERV9TVFJJTkdfWkhfQ04QABIeCh'
    'pMQU5HVUFHRV9DT0RFX1NUUklOR19FTl9HQhABEh4KGkxBTkdVQUdFX0NPREVfU1RSSU5HX0lU'
    'X0lUEAISHgoaTEFOR1VBR0VfQ09ERV9TVFJJTkdfSlBfSlAQAxIeChpMQU5HVUFHRV9DT0RFX1'
    'NUUklOR19FU19FUxAEEh4KGkxBTkdVQUdFX0NPREVfU1RSSU5HX0RFX0RFEAUSHgoaTEFOR1VB'
    'R0VfQ09ERV9TVFJJTkdfS1JfS1IQBhIeChpMQU5HVUFHRV9DT0RFX1NUUklOR19GUl9GUhAHEh'
    '4KGkxBTkdVQUdFX0NPREVfU1RSSU5HX0VOX1VTEAgSHgoaTEFOR1VBR0VfQ09ERV9TVFJJTkdf'
    'UFRfQlIQCRIeChpMQU5HVUFHRV9DT0RFX1NUUklOR19aSF9UVxAKEh8KG0xBTkdVQUdFX0NPRE'
    'VfU1RSSU5HX0VTXzQxORALEh4KGkxBTkdVQUdFX0NPREVfU1RSSU5HX0ZSX0NBEAw=');

@$core.Deprecated('Use numberCapabilityDescriptor instead')
const NumberCapability$json = {
  '1': 'NumberCapability',
  '2': [
    {'1': 'NUMBER_CAPABILITY_SMS', '2': 0},
    {'1': 'NUMBER_CAPABILITY_VOICE', '2': 1},
    {'1': 'NUMBER_CAPABILITY_MMS', '2': 2},
  ],
};

/// Descriptor for `NumberCapability`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List numberCapabilityDescriptor = $convert.base64Decode(
    'ChBOdW1iZXJDYXBhYmlsaXR5EhkKFU5VTUJFUl9DQVBBQklMSVRZX1NNUxAAEhsKF05VTUJFUl'
    '9DQVBBQklMSVRZX1ZPSUNFEAESGQoVTlVNQkVSX0NBUEFCSUxJVFlfTU1TEAI=');

@$core.Deprecated('Use routeTypeDescriptor instead')
const RouteType$json = {
  '1': 'RouteType',
  '2': [
    {'1': 'ROUTE_TYPE_TRANSACTIONAL', '2': 0},
    {'1': 'ROUTE_TYPE_PREMIUM', '2': 1},
    {'1': 'ROUTE_TYPE_PROMOTIONAL', '2': 2},
  ],
};

/// Descriptor for `RouteType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List routeTypeDescriptor = $convert.base64Decode(
    'CglSb3V0ZVR5cGUSHAoYUk9VVEVfVFlQRV9UUkFOU0FDVElPTkFMEAASFgoSUk9VVEVfVFlQRV'
    '9QUkVNSVVNEAESGgoWUk9VVEVfVFlQRV9QUk9NT1RJT05BTBAC');

@$core
    .Deprecated('Use sMSSandboxPhoneNumberVerificationStatusDescriptor instead')
const SMSSandboxPhoneNumberVerificationStatus$json = {
  '1': 'SMSSandboxPhoneNumberVerificationStatus',
  '2': [
    {'1': 'S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_VERIFIED', '2': 0},
    {'1': 'S_M_S_SANDBOX_PHONE_NUMBER_VERIFICATION_STATUS_PENDING', '2': 1},
  ],
};

/// Descriptor for `SMSSandboxPhoneNumberVerificationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sMSSandboxPhoneNumberVerificationStatusDescriptor =
    $convert.base64Decode(
        'CidTTVNTYW5kYm94UGhvbmVOdW1iZXJWZXJpZmljYXRpb25TdGF0dXMSOwo3U19NX1NfU0FORE'
        'JPWF9QSE9ORV9OVU1CRVJfVkVSSUZJQ0FUSU9OX1NUQVRVU19WRVJJRklFRBAAEjoKNlNfTV9T'
        'X1NBTkRCT1hfUEhPTkVfTlVNQkVSX1ZFUklGSUNBVElPTl9TVEFUVVNfUEVORElORxAB');

@$core.Deprecated('Use addPermissionInputDescriptor instead')
const AddPermissionInput$json = {
  '1': 'AddPermissionInput',
  '2': [
    {'1': 'awsaccountid', '3': 370093421, '4': 3, '5': 9, '10': 'awsaccountid'},
    {'1': 'actionname', '3': 115541053, '4': 3, '5': 9, '10': 'actionname'},
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `AddPermissionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addPermissionInputDescriptor = $convert.base64Decode(
    'ChJBZGRQZXJtaXNzaW9uSW5wdXQSJgoMYXdzYWNjb3VudGlkGO3avLABIAMoCVIMYXdzYWNjb3'
    'VudGlkEiEKCmFjdGlvbm5hbWUYvYiMNyADKAlSCmFjdGlvbm5hbWUSGAoFbGFiZWwYnuWz9gEg'
    'ASgJUgVsYWJlbBIdCgh0b3BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use authorizationErrorExceptionDescriptor instead')
const AuthorizationErrorException$json = {
  '1': 'AuthorizationErrorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AuthorizationErrorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List authorizationErrorExceptionDescriptor =
    $convert.base64Decode(
        'ChtBdXRob3JpemF0aW9uRXJyb3JFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use batchEntryIdsNotDistinctExceptionDescriptor instead')
const BatchEntryIdsNotDistinctException$json = {
  '1': 'BatchEntryIdsNotDistinctException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BatchEntryIdsNotDistinctException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchEntryIdsNotDistinctExceptionDescriptor =
    $convert.base64Decode(
        'CiFCYXRjaEVudHJ5SWRzTm90RGlzdGluY3RFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use batchRequestTooLongExceptionDescriptor instead')
const BatchRequestTooLongException$json = {
  '1': 'BatchRequestTooLongException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BatchRequestTooLongException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchRequestTooLongExceptionDescriptor =
    $convert.base64Decode(
        'ChxCYXRjaFJlcXVlc3RUb29Mb25nRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use batchResultErrorEntryDescriptor instead')
const BatchResultErrorEntry$json = {
  '1': 'BatchResultErrorEntry',
  '2': [
    {'1': 'code', '3': 425572629, '4': 1, '5': 9, '10': 'code'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'senderfault', '3': 28412929, '4': 1, '5': 8, '10': 'senderfault'},
  ],
};

/// Descriptor for `BatchResultErrorEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchResultErrorEntryDescriptor = $convert.base64Decode(
    'ChVCYXRjaFJlc3VsdEVycm9yRW50cnkSFgoEY29kZRiV8vbKASABKAlSBGNvZGUSEgoCaWQYgf'
    'KitwEgASgJUgJpZBIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdlEiMKC3NlbmRlcmZhdWx0'
    'GIGYxg0gASgIUgtzZW5kZXJmYXVsdA==');

@$core.Deprecated('Use checkIfPhoneNumberIsOptedOutInputDescriptor instead')
const CheckIfPhoneNumberIsOptedOutInput$json = {
  '1': 'CheckIfPhoneNumberIsOptedOutInput',
  '2': [
    {'1': 'phonenumber', '3': 156377999, '4': 1, '5': 9, '10': 'phonenumber'},
  ],
};

/// Descriptor for `CheckIfPhoneNumberIsOptedOutInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkIfPhoneNumberIsOptedOutInputDescriptor =
    $convert.base64Decode(
        'CiFDaGVja0lmUGhvbmVOdW1iZXJJc09wdGVkT3V0SW5wdXQSIwoLcGhvbmVudW1iZXIYj8fISi'
        'ABKAlSC3Bob25lbnVtYmVy');

@$core.Deprecated('Use checkIfPhoneNumberIsOptedOutResponseDescriptor instead')
const CheckIfPhoneNumberIsOptedOutResponse$json = {
  '1': 'CheckIfPhoneNumberIsOptedOutResponse',
  '2': [
    {'1': 'isoptedout', '3': 430412750, '4': 1, '5': 8, '10': 'isoptedout'},
  ],
};

/// Descriptor for `CheckIfPhoneNumberIsOptedOutResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List checkIfPhoneNumberIsOptedOutResponseDescriptor =
    $convert.base64Decode(
        'CiRDaGVja0lmUGhvbmVOdW1iZXJJc09wdGVkT3V0UmVzcG9uc2USIgoKaXNvcHRlZG91dBjOp5'
        '7NASABKAhSCmlzb3B0ZWRvdXQ=');

@$core.Deprecated('Use concurrentAccessExceptionDescriptor instead')
const ConcurrentAccessException$json = {
  '1': 'ConcurrentAccessException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConcurrentAccessException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrentAccessExceptionDescriptor =
    $convert.base64Decode(
        'ChlDb25jdXJyZW50QWNjZXNzRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use confirmSubscriptionInputDescriptor instead')
const ConfirmSubscriptionInput$json = {
  '1': 'ConfirmSubscriptionInput',
  '2': [
    {
      '1': 'authenticateonunsubscribe',
      '3': 529251911,
      '4': 1,
      '5': 9,
      '10': 'authenticateonunsubscribe'
    },
    {'1': 'token', '3': 439704531, '4': 1, '5': 9, '10': 'token'},
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `ConfirmSubscriptionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmSubscriptionInputDescriptor = $convert.base64Decode(
    'ChhDb25maXJtU3Vic2NyaXB0aW9uSW5wdXQSQAoZYXV0aGVudGljYXRlb251bnN1YnNjcmliZR'
    'jH/K78ASABKAlSGWF1dGhlbnRpY2F0ZW9udW5zdWJzY3JpYmUSGAoFdG9rZW4Y07fV0QEgASgJ'
    'UgV0b2tlbhIdCgh0b3BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use confirmSubscriptionResponseDescriptor instead')
const ConfirmSubscriptionResponse$json = {
  '1': 'ConfirmSubscriptionResponse',
  '2': [
    {
      '1': 'subscriptionarn',
      '3': 279547820,
      '4': 1,
      '5': 9,
      '10': 'subscriptionarn'
    },
  ],
};

/// Descriptor for `ConfirmSubscriptionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List confirmSubscriptionResponseDescriptor =
    $convert.base64Decode(
        'ChtDb25maXJtU3Vic2NyaXB0aW9uUmVzcG9uc2USLAoPc3Vic2NyaXB0aW9uYXJuGKyfpoUBIA'
        'EoCVIPc3Vic2NyaXB0aW9uYXJu');

@$core.Deprecated('Use createEndpointResponseDescriptor instead')
const CreateEndpointResponse$json = {
  '1': 'CreateEndpointResponse',
  '2': [
    {'1': 'endpointarn', '3': 32228660, '4': 1, '5': 9, '10': 'endpointarn'},
  ],
};

/// Descriptor for `CreateEndpointResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEndpointResponseDescriptor =
    $convert.base64Decode(
        'ChZDcmVhdGVFbmRwb2ludFJlc3BvbnNlEiMKC2VuZHBvaW50YXJuGLSKrw8gASgJUgtlbmRwb2'
        'ludGFybg==');

@$core.Deprecated('Use createPlatformApplicationInputDescriptor instead')
const CreatePlatformApplicationInput$json = {
  '1': 'CreatePlatformApplicationInput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.CreatePlatformApplicationInput.AttributesEntry',
      '10': 'attributes'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'platform', '3': 468905683, '4': 1, '5': 9, '10': 'platform'},
  ],
  '3': [CreatePlatformApplicationInput_AttributesEntry$json],
};

@$core.Deprecated('Use createPlatformApplicationInputDescriptor instead')
const CreatePlatformApplicationInput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreatePlatformApplicationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPlatformApplicationInputDescriptor = $convert.base64Decode(
    'Ch5DcmVhdGVQbGF0Zm9ybUFwcGxpY2F0aW9uSW5wdXQSVgoKYXR0cmlidXRlcxi1qftjIAMoCz'
    'IzLnNucy5DcmVhdGVQbGF0Zm9ybUFwcGxpY2F0aW9uSW5wdXQuQXR0cmlidXRlc0VudHJ5Ugph'
    'dHRyaWJ1dGVzEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSHgoIcGxhdGZvcm0Y093L3wEgASgJUg'
    'hwbGF0Zm9ybRo9Cg9BdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUY'
    'AiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createPlatformApplicationResponseDescriptor instead')
const CreatePlatformApplicationResponse$json = {
  '1': 'CreatePlatformApplicationResponse',
  '2': [
    {
      '1': 'platformapplicationarn',
      '3': 241250568,
      '4': 1,
      '5': 9,
      '10': 'platformapplicationarn'
    },
  ],
};

/// Descriptor for `CreatePlatformApplicationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPlatformApplicationResponseDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVQbGF0Zm9ybUFwcGxpY2F0aW9uUmVzcG9uc2USOQoWcGxhdGZvcm1hcHBsaWNhdG'
        'lvbmFybhiI4oRzIAEoCVIWcGxhdGZvcm1hcHBsaWNhdGlvbmFybg==');

@$core.Deprecated('Use createPlatformEndpointInputDescriptor instead')
const CreatePlatformEndpointInput$json = {
  '1': 'CreatePlatformEndpointInput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.CreatePlatformEndpointInput.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'customuserdata',
      '3': 388962962,
      '4': 1,
      '5': 9,
      '10': 'customuserdata'
    },
    {
      '1': 'platformapplicationarn',
      '3': 241250568,
      '4': 1,
      '5': 9,
      '10': 'platformapplicationarn'
    },
    {'1': 'token', '3': 439704531, '4': 1, '5': 9, '10': 'token'},
  ],
  '3': [CreatePlatformEndpointInput_AttributesEntry$json],
};

@$core.Deprecated('Use createPlatformEndpointInputDescriptor instead')
const CreatePlatformEndpointInput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreatePlatformEndpointInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createPlatformEndpointInputDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVQbGF0Zm9ybUVuZHBvaW50SW5wdXQSUwoKYXR0cmlidXRlcxi1qftjIAMoCzIwLn'
    'Nucy5DcmVhdGVQbGF0Zm9ybUVuZHBvaW50SW5wdXQuQXR0cmlidXRlc0VudHJ5UgphdHRyaWJ1'
    'dGVzEioKDmN1c3RvbXVzZXJkYXRhGJK1vLkBIAEoCVIOY3VzdG9tdXNlcmRhdGESOQoWcGxhdG'
    'Zvcm1hcHBsaWNhdGlvbmFybhiI4oRzIAEoCVIWcGxhdGZvcm1hcHBsaWNhdGlvbmFybhIYCgV0'
    'b2tlbhjTt9XRASABKAlSBXRva2VuGj0KD0F0dHJpYnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2'
    'tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use createSMSSandboxPhoneNumberInputDescriptor instead')
const CreateSMSSandboxPhoneNumberInput$json = {
  '1': 'CreateSMSSandboxPhoneNumberInput',
  '2': [
    {
      '1': 'languagecode',
      '3': 281903107,
      '4': 1,
      '5': 14,
      '6': '.sns.LanguageCodeString',
      '10': 'languagecode'
    },
    {'1': 'phonenumber', '3': 379600239, '4': 1, '5': 9, '10': 'phonenumber'},
  ],
};

/// Descriptor for `CreateSMSSandboxPhoneNumberInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSMSSandboxPhoneNumberInputDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVTTVNTYW5kYm94UGhvbmVOdW1iZXJJbnB1dBI/CgxsYW5ndWFnZWNvZGUYg4C2hg'
        'EgASgOMhcuc25zLkxhbmd1YWdlQ29kZVN0cmluZ1IMbGFuZ3VhZ2Vjb2RlEiQKC3Bob25lbnVt'
        'YmVyGO/6gLUBIAEoCVILcGhvbmVudW1iZXI=');

@$core.Deprecated('Use createSMSSandboxPhoneNumberResultDescriptor instead')
const CreateSMSSandboxPhoneNumberResult$json = {
  '1': 'CreateSMSSandboxPhoneNumberResult',
};

/// Descriptor for `CreateSMSSandboxPhoneNumberResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createSMSSandboxPhoneNumberResultDescriptor =
    $convert.base64Decode('CiFDcmVhdGVTTVNTYW5kYm94UGhvbmVOdW1iZXJSZXN1bHQ=');

@$core.Deprecated('Use createTopicInputDescriptor instead')
const CreateTopicInput$json = {
  '1': 'CreateTopicInput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.CreateTopicInput.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'dataprotectionpolicy',
      '3': 519444573,
      '4': 1,
      '5': 9,
      '10': 'dataprotectionpolicy'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sns.Tag',
      '10': 'tags'
    },
  ],
  '3': [CreateTopicInput_AttributesEntry$json],
};

@$core.Deprecated('Use createTopicInputDescriptor instead')
const CreateTopicInput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateTopicInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTopicInputDescriptor = $convert.base64Decode(
    'ChBDcmVhdGVUb3BpY0lucHV0EkgKCmF0dHJpYnV0ZXMYtan7YyADKAsyJS5zbnMuQ3JlYXRlVG'
    '9waWNJbnB1dC5BdHRyaWJ1dGVzRW50cnlSCmF0dHJpYnV0ZXMSNgoUZGF0YXByb3RlY3Rpb25w'
    'b2xpY3kY3bDY9wEgASgJUhRkYXRhcHJvdGVjdGlvbnBvbGljeRIVCgRuYW1lGIfmgX8gASgJUg'
    'RuYW1lEiAKBHRhZ3MYwcH2tQEgAygLMgguc25zLlRhZ1IEdGFncxo9Cg9BdHRyaWJ1dGVzRW50'
    'cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use createTopicResponseDescriptor instead')
const CreateTopicResponse$json = {
  '1': 'CreateTopicResponse',
  '2': [
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `CreateTopicResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTopicResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVUb3BpY1Jlc3BvbnNlEh0KCHRvcGljYXJuGJz0zg4gASgJUgh0b3BpY2Fybg==');

@$core.Deprecated('Use deleteEndpointInputDescriptor instead')
const DeleteEndpointInput$json = {
  '1': 'DeleteEndpointInput',
  '2': [
    {'1': 'endpointarn', '3': 32228660, '4': 1, '5': 9, '10': 'endpointarn'},
  ],
};

/// Descriptor for `DeleteEndpointInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEndpointInputDescriptor = $convert.base64Decode(
    'ChNEZWxldGVFbmRwb2ludElucHV0EiMKC2VuZHBvaW50YXJuGLSKrw8gASgJUgtlbmRwb2ludG'
    'Fybg==');

@$core.Deprecated('Use deletePlatformApplicationInputDescriptor instead')
const DeletePlatformApplicationInput$json = {
  '1': 'DeletePlatformApplicationInput',
  '2': [
    {
      '1': 'platformapplicationarn',
      '3': 241250568,
      '4': 1,
      '5': 9,
      '10': 'platformapplicationarn'
    },
  ],
};

/// Descriptor for `DeletePlatformApplicationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deletePlatformApplicationInputDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVQbGF0Zm9ybUFwcGxpY2F0aW9uSW5wdXQSOQoWcGxhdGZvcm1hcHBsaWNhdGlvbm'
        'FybhiI4oRzIAEoCVIWcGxhdGZvcm1hcHBsaWNhdGlvbmFybg==');

@$core.Deprecated('Use deleteSMSSandboxPhoneNumberInputDescriptor instead')
const DeleteSMSSandboxPhoneNumberInput$json = {
  '1': 'DeleteSMSSandboxPhoneNumberInput',
  '2': [
    {'1': 'phonenumber', '3': 379600239, '4': 1, '5': 9, '10': 'phonenumber'},
  ],
};

/// Descriptor for `DeleteSMSSandboxPhoneNumberInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSMSSandboxPhoneNumberInputDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVTTVNTYW5kYm94UGhvbmVOdW1iZXJJbnB1dBIkCgtwaG9uZW51bWJlchjv+oC1AS'
        'ABKAlSC3Bob25lbnVtYmVy');

@$core.Deprecated('Use deleteSMSSandboxPhoneNumberResultDescriptor instead')
const DeleteSMSSandboxPhoneNumberResult$json = {
  '1': 'DeleteSMSSandboxPhoneNumberResult',
};

/// Descriptor for `DeleteSMSSandboxPhoneNumberResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSMSSandboxPhoneNumberResultDescriptor =
    $convert.base64Decode('CiFEZWxldGVTTVNTYW5kYm94UGhvbmVOdW1iZXJSZXN1bHQ=');

@$core.Deprecated('Use deleteTopicInputDescriptor instead')
const DeleteTopicInput$json = {
  '1': 'DeleteTopicInput',
  '2': [
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `DeleteTopicInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTopicInputDescriptor = $convert.base64Decode(
    'ChBEZWxldGVUb3BpY0lucHV0Eh0KCHRvcGljYXJuGJz0zg4gASgJUgh0b3BpY2Fybg==');

@$core.Deprecated('Use emptyBatchRequestExceptionDescriptor instead')
const EmptyBatchRequestException$json = {
  '1': 'EmptyBatchRequestException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EmptyBatchRequestException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emptyBatchRequestExceptionDescriptor =
    $convert.base64Decode(
        'ChpFbXB0eUJhdGNoUmVxdWVzdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use endpointDescriptor instead')
const Endpoint$json = {
  '1': 'Endpoint',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.Endpoint.AttributesEntry',
      '10': 'attributes'
    },
    {'1': 'endpointarn', '3': 32228660, '4': 1, '5': 9, '10': 'endpointarn'},
  ],
  '3': [Endpoint_AttributesEntry$json],
};

@$core.Deprecated('Use endpointDescriptor instead')
const Endpoint_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `Endpoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endpointDescriptor = $convert.base64Decode(
    'CghFbmRwb2ludBJACgphdHRyaWJ1dGVzGLWp+2MgAygLMh0uc25zLkVuZHBvaW50LkF0dHJpYn'
    'V0ZXNFbnRyeVIKYXR0cmlidXRlcxIjCgtlbmRwb2ludGFybhi0iq8PIAEoCVILZW5kcG9pbnRh'
    'cm4aPQoPQXR0cmlidXRlc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUg'
    'V2YWx1ZToCOAE=');

@$core.Deprecated('Use endpointDisabledExceptionDescriptor instead')
const EndpointDisabledException$json = {
  '1': 'EndpointDisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EndpointDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List endpointDisabledExceptionDescriptor =
    $convert.base64Decode(
        'ChlFbmRwb2ludERpc2FibGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use filterPolicyLimitExceededExceptionDescriptor instead')
const FilterPolicyLimitExceededException$json = {
  '1': 'FilterPolicyLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `FilterPolicyLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List filterPolicyLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiJGaWx0ZXJQb2xpY3lMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use getDataProtectionPolicyInputDescriptor instead')
const GetDataProtectionPolicyInput$json = {
  '1': 'GetDataProtectionPolicyInput',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `GetDataProtectionPolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataProtectionPolicyInputDescriptor =
    $convert.base64Decode(
        'ChxHZXREYXRhUHJvdGVjdGlvblBvbGljeUlucHV0EiQKC3Jlc291cmNlYXJuGK342a0BIAEoCV'
        'ILcmVzb3VyY2Vhcm4=');

@$core.Deprecated('Use getDataProtectionPolicyResponseDescriptor instead')
const GetDataProtectionPolicyResponse$json = {
  '1': 'GetDataProtectionPolicyResponse',
  '2': [
    {
      '1': 'dataprotectionpolicy',
      '3': 519444573,
      '4': 1,
      '5': 9,
      '10': 'dataprotectionpolicy'
    },
  ],
};

/// Descriptor for `GetDataProtectionPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDataProtectionPolicyResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXREYXRhUHJvdGVjdGlvblBvbGljeVJlc3BvbnNlEjYKFGRhdGFwcm90ZWN0aW9ucG9saW'
        'N5GN2w2PcBIAEoCVIUZGF0YXByb3RlY3Rpb25wb2xpY3k=');

@$core.Deprecated('Use getEndpointAttributesInputDescriptor instead')
const GetEndpointAttributesInput$json = {
  '1': 'GetEndpointAttributesInput',
  '2': [
    {'1': 'endpointarn', '3': 32228660, '4': 1, '5': 9, '10': 'endpointarn'},
  ],
};

/// Descriptor for `GetEndpointAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEndpointAttributesInputDescriptor =
    $convert.base64Decode(
        'ChpHZXRFbmRwb2ludEF0dHJpYnV0ZXNJbnB1dBIjCgtlbmRwb2ludGFybhi0iq8PIAEoCVILZW'
        '5kcG9pbnRhcm4=');

@$core.Deprecated('Use getEndpointAttributesResponseDescriptor instead')
const GetEndpointAttributesResponse$json = {
  '1': 'GetEndpointAttributesResponse',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.GetEndpointAttributesResponse.AttributesEntry',
      '10': 'attributes'
    },
  ],
  '3': [GetEndpointAttributesResponse_AttributesEntry$json],
};

@$core.Deprecated('Use getEndpointAttributesResponseDescriptor instead')
const GetEndpointAttributesResponse_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetEndpointAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEndpointAttributesResponseDescriptor = $convert.base64Decode(
    'Ch1HZXRFbmRwb2ludEF0dHJpYnV0ZXNSZXNwb25zZRJVCgphdHRyaWJ1dGVzGLWp+2MgAygLMj'
    'Iuc25zLkdldEVuZHBvaW50QXR0cmlidXRlc1Jlc3BvbnNlLkF0dHJpYnV0ZXNFbnRyeVIKYXR0'
    'cmlidXRlcxo9Cg9BdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAi'
    'ABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getPlatformApplicationAttributesInputDescriptor instead')
const GetPlatformApplicationAttributesInput$json = {
  '1': 'GetPlatformApplicationAttributesInput',
  '2': [
    {
      '1': 'platformapplicationarn',
      '3': 241250568,
      '4': 1,
      '5': 9,
      '10': 'platformapplicationarn'
    },
  ],
};

/// Descriptor for `GetPlatformApplicationAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPlatformApplicationAttributesInputDescriptor =
    $convert.base64Decode(
        'CiVHZXRQbGF0Zm9ybUFwcGxpY2F0aW9uQXR0cmlidXRlc0lucHV0EjkKFnBsYXRmb3JtYXBwbG'
        'ljYXRpb25hcm4YiOKEcyABKAlSFnBsYXRmb3JtYXBwbGljYXRpb25hcm4=');

@$core.Deprecated(
    'Use getPlatformApplicationAttributesResponseDescriptor instead')
const GetPlatformApplicationAttributesResponse$json = {
  '1': 'GetPlatformApplicationAttributesResponse',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.GetPlatformApplicationAttributesResponse.AttributesEntry',
      '10': 'attributes'
    },
  ],
  '3': [GetPlatformApplicationAttributesResponse_AttributesEntry$json],
};

@$core.Deprecated(
    'Use getPlatformApplicationAttributesResponseDescriptor instead')
const GetPlatformApplicationAttributesResponse_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetPlatformApplicationAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPlatformApplicationAttributesResponseDescriptor =
    $convert.base64Decode(
        'CihHZXRQbGF0Zm9ybUFwcGxpY2F0aW9uQXR0cmlidXRlc1Jlc3BvbnNlEmAKCmF0dHJpYnV0ZX'
        'MYtan7YyADKAsyPS5zbnMuR2V0UGxhdGZvcm1BcHBsaWNhdGlvbkF0dHJpYnV0ZXNSZXNwb25z'
        'ZS5BdHRyaWJ1dGVzRW50cnlSCmF0dHJpYnV0ZXMaPQoPQXR0cmlidXRlc0VudHJ5EhAKA2tleR'
        'gBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getSMSAttributesInputDescriptor instead')
const GetSMSAttributesInput$json = {
  '1': 'GetSMSAttributesInput',
  '2': [
    {'1': 'attributes', '3': 33545109, '4': 3, '5': 9, '10': 'attributes'},
  ],
};

/// Descriptor for `GetSMSAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSMSAttributesInputDescriptor = $convert.base64Decode(
    'ChVHZXRTTVNBdHRyaWJ1dGVzSW5wdXQSIQoKYXR0cmlidXRlcxiVt/8PIAMoCVIKYXR0cmlidX'
    'Rlcw==');

@$core.Deprecated('Use getSMSAttributesResponseDescriptor instead')
const GetSMSAttributesResponse$json = {
  '1': 'GetSMSAttributesResponse',
  '2': [
    {
      '1': 'attributes',
      '3': 33545109,
      '4': 3,
      '5': 11,
      '6': '.sns.GetSMSAttributesResponse.AttributesEntry',
      '10': 'attributes'
    },
  ],
  '3': [GetSMSAttributesResponse_AttributesEntry$json],
};

@$core.Deprecated('Use getSMSAttributesResponseDescriptor instead')
const GetSMSAttributesResponse_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetSMSAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSMSAttributesResponseDescriptor = $convert.base64Decode(
    'ChhHZXRTTVNBdHRyaWJ1dGVzUmVzcG9uc2USUAoKYXR0cmlidXRlcxiVt/8PIAMoCzItLnNucy'
    '5HZXRTTVNBdHRyaWJ1dGVzUmVzcG9uc2UuQXR0cmlidXRlc0VudHJ5UgphdHRyaWJ1dGVzGj0K'
    'D0F0dHJpYnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdW'
    'U6AjgB');

@$core.Deprecated('Use getSMSSandboxAccountStatusInputDescriptor instead')
const GetSMSSandboxAccountStatusInput$json = {
  '1': 'GetSMSSandboxAccountStatusInput',
};

/// Descriptor for `GetSMSSandboxAccountStatusInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSMSSandboxAccountStatusInputDescriptor =
    $convert.base64Decode('Ch9HZXRTTVNTYW5kYm94QWNjb3VudFN0YXR1c0lucHV0');

@$core.Deprecated('Use getSMSSandboxAccountStatusResultDescriptor instead')
const GetSMSSandboxAccountStatusResult$json = {
  '1': 'GetSMSSandboxAccountStatusResult',
  '2': [
    {'1': 'isinsandbox', '3': 246264378, '4': 1, '5': 8, '10': 'isinsandbox'},
  ],
};

/// Descriptor for `GetSMSSandboxAccountStatusResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSMSSandboxAccountStatusResultDescriptor =
    $convert.base64Decode(
        'CiBHZXRTTVNTYW5kYm94QWNjb3VudFN0YXR1c1Jlc3VsdBIjCgtpc2luc2FuZGJveBi65LZ1IA'
        'EoCFILaXNpbnNhbmRib3g=');

@$core.Deprecated('Use getSubscriptionAttributesInputDescriptor instead')
const GetSubscriptionAttributesInput$json = {
  '1': 'GetSubscriptionAttributesInput',
  '2': [
    {
      '1': 'subscriptionarn',
      '3': 279547820,
      '4': 1,
      '5': 9,
      '10': 'subscriptionarn'
    },
  ],
};

/// Descriptor for `GetSubscriptionAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubscriptionAttributesInputDescriptor =
    $convert.base64Decode(
        'Ch5HZXRTdWJzY3JpcHRpb25BdHRyaWJ1dGVzSW5wdXQSLAoPc3Vic2NyaXB0aW9uYXJuGKyfpo'
        'UBIAEoCVIPc3Vic2NyaXB0aW9uYXJu');

@$core.Deprecated('Use getSubscriptionAttributesResponseDescriptor instead')
const GetSubscriptionAttributesResponse$json = {
  '1': 'GetSubscriptionAttributesResponse',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.GetSubscriptionAttributesResponse.AttributesEntry',
      '10': 'attributes'
    },
  ],
  '3': [GetSubscriptionAttributesResponse_AttributesEntry$json],
};

@$core.Deprecated('Use getSubscriptionAttributesResponseDescriptor instead')
const GetSubscriptionAttributesResponse_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetSubscriptionAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSubscriptionAttributesResponseDescriptor =
    $convert.base64Decode(
        'CiFHZXRTdWJzY3JpcHRpb25BdHRyaWJ1dGVzUmVzcG9uc2USWQoKYXR0cmlidXRlcxi1qftjIA'
        'MoCzI2LnNucy5HZXRTdWJzY3JpcHRpb25BdHRyaWJ1dGVzUmVzcG9uc2UuQXR0cmlidXRlc0Vu'
        'dHJ5UgphdHRyaWJ1dGVzGj0KD0F0dHJpYnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCg'
        'V2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use getTopicAttributesInputDescriptor instead')
const GetTopicAttributesInput$json = {
  '1': 'GetTopicAttributesInput',
  '2': [
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `GetTopicAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTopicAttributesInputDescriptor =
    $convert.base64Decode(
        'ChdHZXRUb3BpY0F0dHJpYnV0ZXNJbnB1dBIdCgh0b3BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm'
        '4=');

@$core.Deprecated('Use getTopicAttributesResponseDescriptor instead')
const GetTopicAttributesResponse$json = {
  '1': 'GetTopicAttributesResponse',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.GetTopicAttributesResponse.AttributesEntry',
      '10': 'attributes'
    },
  ],
  '3': [GetTopicAttributesResponse_AttributesEntry$json],
};

@$core.Deprecated('Use getTopicAttributesResponseDescriptor instead')
const GetTopicAttributesResponse_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetTopicAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTopicAttributesResponseDescriptor = $convert.base64Decode(
    'ChpHZXRUb3BpY0F0dHJpYnV0ZXNSZXNwb25zZRJSCgphdHRyaWJ1dGVzGLWp+2MgAygLMi8uc2'
    '5zLkdldFRvcGljQXR0cmlidXRlc1Jlc3BvbnNlLkF0dHJpYnV0ZXNFbnRyeVIKYXR0cmlidXRl'
    'cxo9Cg9BdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBX'
    'ZhbHVlOgI4AQ==');

@$core.Deprecated('Use internalErrorExceptionDescriptor instead')
const InternalErrorException$json = {
  '1': 'InternalErrorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalErrorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalErrorExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbnRlcm5hbEVycm9yRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidBatchEntryIdExceptionDescriptor instead')
const InvalidBatchEntryIdException$json = {
  '1': 'InvalidBatchEntryIdException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidBatchEntryIdException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidBatchEntryIdExceptionDescriptor =
    $convert.base64Decode(
        'ChxJbnZhbGlkQmF0Y2hFbnRyeUlkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use invalidParameterExceptionDescriptor instead')
const InvalidParameterException$json = {
  '1': 'InvalidParameterException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkUGFyYW1ldGVyRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidParameterValueExceptionDescriptor instead')
const InvalidParameterValueException$json = {
  '1': 'InvalidParameterValueException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidParameterValueException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidParameterValueExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JbnZhbGlkUGFyYW1ldGVyVmFsdWVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use invalidSecurityExceptionDescriptor instead')
const InvalidSecurityException$json = {
  '1': 'InvalidSecurityException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSecurityException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSecurityExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkU2VjdXJpdHlFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidStateExceptionDescriptor instead')
const InvalidStateException$json = {
  '1': 'InvalidStateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidStateExceptionDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkU3RhdGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use kMSAccessDeniedExceptionDescriptor instead')
const KMSAccessDeniedException$json = {
  '1': 'KMSAccessDeniedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSAccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSAccessDeniedExceptionDescriptor =
    $convert.base64Decode(
        'ChhLTVNBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use kMSDisabledExceptionDescriptor instead')
const KMSDisabledException$json = {
  '1': 'KMSDisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSDisabledExceptionDescriptor =
    $convert.base64Decode(
        'ChRLTVNEaXNhYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use kMSInvalidStateExceptionDescriptor instead')
const KMSInvalidStateException$json = {
  '1': 'KMSInvalidStateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSInvalidStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSInvalidStateExceptionDescriptor =
    $convert.base64Decode(
        'ChhLTVNJbnZhbGlkU3RhdGVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use kMSNotFoundExceptionDescriptor instead')
const KMSNotFoundException$json = {
  '1': 'KMSNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChRLTVNOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use kMSOptInRequiredDescriptor instead')
const KMSOptInRequired$json = {
  '1': 'KMSOptInRequired',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSOptInRequired`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSOptInRequiredDescriptor = $convert.base64Decode(
    'ChBLTVNPcHRJblJlcXVpcmVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use kMSThrottlingExceptionDescriptor instead')
const KMSThrottlingException$json = {
  '1': 'KMSThrottlingException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSThrottlingExceptionDescriptor =
    $convert.base64Decode(
        'ChZLTVNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core
    .Deprecated('Use listEndpointsByPlatformApplicationInputDescriptor instead')
const ListEndpointsByPlatformApplicationInput$json = {
  '1': 'ListEndpointsByPlatformApplicationInput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'platformapplicationarn',
      '3': 241250568,
      '4': 1,
      '5': 9,
      '10': 'platformapplicationarn'
    },
  ],
};

/// Descriptor for `ListEndpointsByPlatformApplicationInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEndpointsByPlatformApplicationInputDescriptor =
    $convert.base64Decode(
        'CidMaXN0RW5kcG9pbnRzQnlQbGF0Zm9ybUFwcGxpY2F0aW9uSW5wdXQSHwoJbmV4dHRva2VuGP'
        '6EumcgASgJUgluZXh0dG9rZW4SOQoWcGxhdGZvcm1hcHBsaWNhdGlvbmFybhiI4oRzIAEoCVIW'
        'cGxhdGZvcm1hcHBsaWNhdGlvbmFybg==');

@$core.Deprecated(
    'Use listEndpointsByPlatformApplicationResponseDescriptor instead')
const ListEndpointsByPlatformApplicationResponse$json = {
  '1': 'ListEndpointsByPlatformApplicationResponse',
  '2': [
    {
      '1': 'endpoints',
      '3': 16210494,
      '4': 3,
      '5': 11,
      '6': '.sns.Endpoint',
      '10': 'endpoints'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEndpointsByPlatformApplicationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listEndpointsByPlatformApplicationResponseDescriptor =
    $convert.base64Decode(
        'CipMaXN0RW5kcG9pbnRzQnlQbGF0Zm9ybUFwcGxpY2F0aW9uUmVzcG9uc2USLgoJZW5kcG9pbn'
        'RzGL603QcgAygLMg0uc25zLkVuZHBvaW50UgllbmRwb2ludHMSHwoJbmV4dHRva2VuGP6Eumcg'
        'ASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listOriginationNumbersRequestDescriptor instead')
const ListOriginationNumbersRequest$json = {
  '1': 'ListOriginationNumbersRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListOriginationNumbersRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOriginationNumbersRequestDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0T3JpZ2luYXRpb25OdW1iZXJzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBV'
        'IKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listOriginationNumbersResultDescriptor instead')
const ListOriginationNumbersResult$json = {
  '1': 'ListOriginationNumbersResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'phonenumbers',
      '3': 457192616,
      '4': 3,
      '5': 11,
      '6': '.sns.PhoneNumberInformation',
      '10': 'phonenumbers'
    },
  ],
};

/// Descriptor for `ListOriginationNumbersResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listOriginationNumbersResultDescriptor =
    $convert.base64Decode(
        'ChxMaXN0T3JpZ2luYXRpb25OdW1iZXJzUmVzdWx0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbm'
        'V4dHRva2VuEkMKDHBob25lbnVtYmVycxio6YDaASADKAsyGy5zbnMuUGhvbmVOdW1iZXJJbmZv'
        'cm1hdGlvblIMcGhvbmVudW1iZXJz');

@$core.Deprecated('Use listPhoneNumbersOptedOutInputDescriptor instead')
const ListPhoneNumbersOptedOutInput$json = {
  '1': 'ListPhoneNumbersOptedOutInput',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListPhoneNumbersOptedOutInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPhoneNumbersOptedOutInputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0UGhvbmVOdW1iZXJzT3B0ZWRPdXRJbnB1dBIfCgluZXh0dG9rZW4YnvOdNyABKAlSCW'
        '5leHR0b2tlbg==');

@$core.Deprecated('Use listPhoneNumbersOptedOutResponseDescriptor instead')
const ListPhoneNumbersOptedOutResponse$json = {
  '1': 'ListPhoneNumbersOptedOutResponse',
  '2': [
    {'1': 'nexttoken', '3': 115833246, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'phonenumbers', '3': 156149576, '4': 3, '5': 9, '10': 'phonenumbers'},
  ],
};

/// Descriptor for `ListPhoneNumbersOptedOutResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPhoneNumbersOptedOutResponseDescriptor =
    $convert.base64Decode(
        'CiBMaXN0UGhvbmVOdW1iZXJzT3B0ZWRPdXRSZXNwb25zZRIfCgluZXh0dG9rZW4YnvOdNyABKA'
        'lSCW5leHR0b2tlbhIlCgxwaG9uZW51bWJlcnMYyM66SiADKAlSDHBob25lbnVtYmVycw==');

@$core.Deprecated('Use listPlatformApplicationsInputDescriptor instead')
const ListPlatformApplicationsInput$json = {
  '1': 'ListPlatformApplicationsInput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListPlatformApplicationsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPlatformApplicationsInputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0UGxhdGZvcm1BcHBsaWNhdGlvbnNJbnB1dBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW'
        '5leHR0b2tlbg==');

@$core.Deprecated('Use listPlatformApplicationsResponseDescriptor instead')
const ListPlatformApplicationsResponse$json = {
  '1': 'ListPlatformApplicationsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'platformapplications',
      '3': 209735978,
      '4': 3,
      '5': 11,
      '6': '.sns.PlatformApplication',
      '10': 'platformapplications'
    },
  ],
};

/// Descriptor for `ListPlatformApplicationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listPlatformApplicationsResponseDescriptor =
    $convert.base64Decode(
        'CiBMaXN0UGxhdGZvcm1BcHBsaWNhdGlvbnNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
        'lSCW5leHR0b2tlbhJPChRwbGF0Zm9ybWFwcGxpY2F0aW9ucxiqooFkIAMoCzIYLnNucy5QbGF0'
        'Zm9ybUFwcGxpY2F0aW9uUhRwbGF0Zm9ybWFwcGxpY2F0aW9ucw==');

@$core.Deprecated('Use listSMSSandboxPhoneNumbersInputDescriptor instead')
const ListSMSSandboxPhoneNumbersInput$json = {
  '1': 'ListSMSSandboxPhoneNumbersInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListSMSSandboxPhoneNumbersInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSMSSandboxPhoneNumbersInputDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0U01TU2FuZGJveFBob25lTnVtYmVyc0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgAS'
        'gFUgptYXhyZXN1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listSMSSandboxPhoneNumbersResultDescriptor instead')
const ListSMSSandboxPhoneNumbersResult$json = {
  '1': 'ListSMSSandboxPhoneNumbersResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'phonenumbers',
      '3': 457192616,
      '4': 3,
      '5': 11,
      '6': '.sns.SMSSandboxPhoneNumber',
      '10': 'phonenumbers'
    },
  ],
};

/// Descriptor for `ListSMSSandboxPhoneNumbersResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSMSSandboxPhoneNumbersResultDescriptor =
    $convert.base64Decode(
        'CiBMaXN0U01TU2FuZGJveFBob25lTnVtYmVyc1Jlc3VsdBIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
        'lSCW5leHR0b2tlbhJCCgxwaG9uZW51bWJlcnMYqOmA2gEgAygLMhouc25zLlNNU1NhbmRib3hQ'
        'aG9uZU51bWJlclIMcGhvbmVudW1iZXJz');

@$core.Deprecated('Use listSubscriptionsByTopicInputDescriptor instead')
const ListSubscriptionsByTopicInput$json = {
  '1': 'ListSubscriptionsByTopicInput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `ListSubscriptionsByTopicInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubscriptionsByTopicInputDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0U3Vic2NyaXB0aW9uc0J5VG9waWNJbnB1dBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW'
        '5leHR0b2tlbhIdCgh0b3BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use listSubscriptionsByTopicResponseDescriptor instead')
const ListSubscriptionsByTopicResponse$json = {
  '1': 'ListSubscriptionsByTopicResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'subscriptions',
      '3': 169711430,
      '4': 3,
      '5': 11,
      '6': '.sns.Subscription',
      '10': 'subscriptions'
    },
  ],
};

/// Descriptor for `ListSubscriptionsByTopicResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubscriptionsByTopicResponseDescriptor =
    $convert.base64Decode(
        'CiBMaXN0U3Vic2NyaXB0aW9uc0J5VG9waWNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
        'lSCW5leHR0b2tlbhI6Cg1zdWJzY3JpcHRpb25zGMau9lAgAygLMhEuc25zLlN1YnNjcmlwdGlv'
        'blINc3Vic2NyaXB0aW9ucw==');

@$core.Deprecated('Use listSubscriptionsInputDescriptor instead')
const ListSubscriptionsInput$json = {
  '1': 'ListSubscriptionsInput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListSubscriptionsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubscriptionsInputDescriptor =
    $convert.base64Decode(
        'ChZMaXN0U3Vic2NyaXB0aW9uc0lucHV0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2'
        'Vu');

@$core.Deprecated('Use listSubscriptionsResponseDescriptor instead')
const ListSubscriptionsResponse$json = {
  '1': 'ListSubscriptionsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'subscriptions',
      '3': 169711430,
      '4': 3,
      '5': 11,
      '6': '.sns.Subscription',
      '10': 'subscriptions'
    },
  ],
};

/// Descriptor for `ListSubscriptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSubscriptionsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0U3Vic2NyaXB0aW9uc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
    'Rva2VuEjoKDXN1YnNjcmlwdGlvbnMYxq72UCADKAsyES5zbnMuU3Vic2NyaXB0aW9uUg1zdWJz'
    'Y3JpcHRpb25z');

@$core.Deprecated('Use listTagsForResourceRequestDescriptor instead')
const ListTagsForResourceRequest$json = {
  '1': 'ListTagsForResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0VGFnc0ZvclJlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3'
        'Jlc291cmNlYXJu');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sns.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USIAoEdGFncxjBwfa1ASADKAsyCC5zbnMuVG'
        'FnUgR0YWdz');

@$core.Deprecated('Use listTopicsInputDescriptor instead')
const ListTopicsInput$json = {
  '1': 'ListTopicsInput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListTopicsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTopicsInputDescriptor = $convert.base64Decode(
    'Cg9MaXN0VG9waWNzSW5wdXQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listTopicsResponseDescriptor instead')
const ListTopicsResponse$json = {
  '1': 'ListTopicsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'topics',
      '3': 219850038,
      '4': 3,
      '5': 11,
      '6': '.sns.Topic',
      '10': 'topics'
    },
  ],
};

/// Descriptor for `ListTopicsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTopicsResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0VG9waWNzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SJQ'
    'oGdG9waWNzGLbK6mggAygLMgouc25zLlRvcGljUgZ0b3BpY3M=');

@$core.Deprecated('Use messageAttributeValueDescriptor instead')
const MessageAttributeValue$json = {
  '1': 'MessageAttributeValue',
  '2': [
    {'1': 'binaryvalue', '3': 255476278, '4': 1, '5': 12, '10': 'binaryvalue'},
    {'1': 'datatype', '3': 67988590, '4': 1, '5': 9, '10': 'datatype'},
    {'1': 'stringvalue', '3': 184416138, '4': 1, '5': 9, '10': 'stringvalue'},
  ],
};

/// Descriptor for `MessageAttributeValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageAttributeValueDescriptor = $convert.base64Decode(
    'ChVNZXNzYWdlQXR0cmlidXRlVmFsdWUSIwoLYmluYXJ5dmFsdWUYtoTpeSABKAxSC2JpbmFyeX'
    'ZhbHVlEh0KCGRhdGF0eXBlGO7YtSAgASgJUghkYXRhdHlwZRIjCgtzdHJpbmd2YWx1ZRiK7/dX'
    'IAEoCVILc3RyaW5ndmFsdWU=');

@$core.Deprecated('Use notFoundExceptionDescriptor instead')
const NotFoundException$json = {
  '1': 'NotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notFoundExceptionDescriptor = $convert.base64Decode(
    'ChFOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use optInPhoneNumberInputDescriptor instead')
const OptInPhoneNumberInput$json = {
  '1': 'OptInPhoneNumberInput',
  '2': [
    {'1': 'phonenumber', '3': 156377999, '4': 1, '5': 9, '10': 'phonenumber'},
  ],
};

/// Descriptor for `OptInPhoneNumberInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List optInPhoneNumberInputDescriptor = $convert.base64Decode(
    'ChVPcHRJblBob25lTnVtYmVySW5wdXQSIwoLcGhvbmVudW1iZXIYj8fISiABKAlSC3Bob25lbn'
    'VtYmVy');

@$core.Deprecated('Use optInPhoneNumberResponseDescriptor instead')
const OptInPhoneNumberResponse$json = {
  '1': 'OptInPhoneNumberResponse',
};

/// Descriptor for `OptInPhoneNumberResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List optInPhoneNumberResponseDescriptor =
    $convert.base64Decode('ChhPcHRJblBob25lTnVtYmVyUmVzcG9uc2U=');

@$core.Deprecated('Use optedOutExceptionDescriptor instead')
const OptedOutException$json = {
  '1': 'OptedOutException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OptedOutException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List optedOutExceptionDescriptor = $convert.base64Decode(
    'ChFPcHRlZE91dEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use phoneNumberInformationDescriptor instead')
const PhoneNumberInformation$json = {
  '1': 'PhoneNumberInformation',
  '2': [
    {'1': 'createdat', '3': 258192751, '4': 1, '5': 9, '10': 'createdat'},
    {
      '1': 'iso2countrycode',
      '3': 283246908,
      '4': 1,
      '5': 9,
      '10': 'iso2countrycode'
    },
    {
      '1': 'numbercapabilities',
      '3': 54004711,
      '4': 3,
      '5': 14,
      '6': '.sns.NumberCapability',
      '10': 'numbercapabilities'
    },
    {'1': 'phonenumber', '3': 379600239, '4': 1, '5': 9, '10': 'phonenumber'},
    {
      '1': 'routetype',
      '3': 170172127,
      '4': 1,
      '5': 14,
      '6': '.sns.RouteType',
      '10': 'routetype'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `PhoneNumberInformation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List phoneNumberInformationDescriptor = $convert.base64Decode(
    'ChZQaG9uZU51bWJlckluZm9ybWF0aW9uEh8KCWNyZWF0ZWRhdBjv6o57IAEoCVIJY3JlYXRlZG'
    'F0EiwKD2lzbzJjb3VudHJ5Y29kZRi8goiHASABKAlSD2lzbzJjb3VudHJ5Y29kZRJIChJudW1i'
    'ZXJjYXBhYmlsaXRpZXMY55fgGSADKA4yFS5zbnMuTnVtYmVyQ2FwYWJpbGl0eVISbnVtYmVyY2'
    'FwYWJpbGl0aWVzEiQKC3Bob25lbnVtYmVyGO/6gLUBIAEoCVILcGhvbmVudW1iZXISLwoJcm91'
    'dGV0eXBlGN+9klEgASgOMg4uc25zLlJvdXRlVHlwZVIJcm91dGV0eXBlEhkKBnN0YXR1cxiQ5P'
    'sCIAEoCVIGc3RhdHVz');

@$core.Deprecated('Use platformApplicationDescriptor instead')
const PlatformApplication$json = {
  '1': 'PlatformApplication',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.PlatformApplication.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'platformapplicationarn',
      '3': 241250568,
      '4': 1,
      '5': 9,
      '10': 'platformapplicationarn'
    },
  ],
  '3': [PlatformApplication_AttributesEntry$json],
};

@$core.Deprecated('Use platformApplicationDescriptor instead')
const PlatformApplication_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `PlatformApplication`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List platformApplicationDescriptor = $convert.base64Decode(
    'ChNQbGF0Zm9ybUFwcGxpY2F0aW9uEksKCmF0dHJpYnV0ZXMYtan7YyADKAsyKC5zbnMuUGxhdG'
    'Zvcm1BcHBsaWNhdGlvbi5BdHRyaWJ1dGVzRW50cnlSCmF0dHJpYnV0ZXMSOQoWcGxhdGZvcm1h'
    'cHBsaWNhdGlvbmFybhiI4oRzIAEoCVIWcGxhdGZvcm1hcHBsaWNhdGlvbmFybho9Cg9BdHRyaW'
    'J1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use platformApplicationDisabledExceptionDescriptor instead')
const PlatformApplicationDisabledException$json = {
  '1': 'PlatformApplicationDisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PlatformApplicationDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List platformApplicationDisabledExceptionDescriptor =
    $convert.base64Decode(
        'CiRQbGF0Zm9ybUFwcGxpY2F0aW9uRGlzYWJsZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use publishBatchInputDescriptor instead')
const PublishBatchInput$json = {
  '1': 'PublishBatchInput',
  '2': [
    {
      '1': 'publishbatchrequestentries',
      '3': 327245576,
      '4': 3,
      '5': 11,
      '6': '.sns.PublishBatchRequestEntry',
      '10': 'publishbatchrequestentries'
    },
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `PublishBatchInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishBatchInputDescriptor = $convert.base64Decode(
    'ChFQdWJsaXNoQmF0Y2hJbnB1dBJhChpwdWJsaXNoYmF0Y2hyZXF1ZXN0ZW50cmllcxiIvoWcAS'
    'ADKAsyHS5zbnMuUHVibGlzaEJhdGNoUmVxdWVzdEVudHJ5UhpwdWJsaXNoYmF0Y2hyZXF1ZXN0'
    'ZW50cmllcxIdCgh0b3BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use publishBatchRequestEntryDescriptor instead')
const PublishBatchRequestEntry$json = {
  '1': 'PublishBatchRequestEntry',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'messageattributes',
      '3': 56443766,
      '4': 3,
      '5': 11,
      '6': '.sns.PublishBatchRequestEntry.MessageattributesEntry',
      '10': 'messageattributes'
    },
    {
      '1': 'messagededuplicationid',
      '3': 379560665,
      '4': 1,
      '5': 9,
      '10': 'messagededuplicationid'
    },
    {
      '1': 'messagegroupid',
      '3': 419537435,
      '4': 1,
      '5': 9,
      '10': 'messagegroupid'
    },
    {
      '1': 'messagestructure',
      '3': 402672330,
      '4': 1,
      '5': 9,
      '10': 'messagestructure'
    },
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
  ],
  '3': [PublishBatchRequestEntry_MessageattributesEntry$json],
};

@$core.Deprecated('Use publishBatchRequestEntryDescriptor instead')
const PublishBatchRequestEntry_MessageattributesEntry$json = {
  '1': 'MessageattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.sns.MessageAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `PublishBatchRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishBatchRequestEntryDescriptor = $convert.base64Decode(
    'ChhQdWJsaXNoQmF0Y2hSZXF1ZXN0RW50cnkSEgoCaWQYgfKitwEgASgJUgJpZBIbCgdtZXNzYW'
    'dlGIWzu3AgASgJUgdtZXNzYWdlEmUKEW1lc3NhZ2VhdHRyaWJ1dGVzGPaG9RogAygLMjQuc25z'
    'LlB1Ymxpc2hCYXRjaFJlcXVlc3RFbnRyeS5NZXNzYWdlYXR0cmlidXRlc0VudHJ5UhFtZXNzYW'
    'dlYXR0cmlidXRlcxI6ChZtZXNzYWdlZGVkdXBsaWNhdGlvbmlkGNnF/rQBIAEoCVIWbWVzc2Fn'
    'ZWRlZHVwbGljYXRpb25pZBIqCg5tZXNzYWdlZ3JvdXBpZBibxIbIASABKAlSDm1lc3NhZ2Vncm'
    '91cGlkEi4KEG1lc3NhZ2VzdHJ1Y3R1cmUYypWBwAEgASgJUhBtZXNzYWdlc3RydWN0dXJlEhsK'
    'B3N1YmplY3QY8MnkAyABKAlSB3N1YmplY3QaYAoWTWVzc2FnZWF0dHJpYnV0ZXNFbnRyeRIQCg'
    'NrZXkYASABKAlSA2tleRIwCgV2YWx1ZRgCIAEoCzIaLnNucy5NZXNzYWdlQXR0cmlidXRlVmFs'
    'dWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use publishBatchResponseDescriptor instead')
const PublishBatchResponse$json = {
  '1': 'PublishBatchResponse',
  '2': [
    {
      '1': 'failed',
      '3': 360301525,
      '4': 3,
      '5': 11,
      '6': '.sns.BatchResultErrorEntry',
      '10': 'failed'
    },
    {
      '1': 'successful',
      '3': 412818844,
      '4': 3,
      '5': 11,
      '6': '.sns.PublishBatchResultEntry',
      '10': 'successful'
    },
  ],
};

/// Descriptor for `PublishBatchResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishBatchResponseDescriptor = $convert.base64Decode(
    'ChRQdWJsaXNoQmF0Y2hSZXNwb25zZRI2CgZmYWlsZWQY1YfnqwEgAygLMhouc25zLkJhdGNoUm'
    'VzdWx0RXJyb3JFbnRyeVIGZmFpbGVkEkAKCnN1Y2Nlc3NmdWwYnLvsxAEgAygLMhwuc25zLlB1'
    'Ymxpc2hCYXRjaFJlc3VsdEVudHJ5UgpzdWNjZXNzZnVs');

@$core.Deprecated('Use publishBatchResultEntryDescriptor instead')
const PublishBatchResultEntry$json = {
  '1': 'PublishBatchResultEntry',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
    {
      '1': 'sequencenumber',
      '3': 98094362,
      '4': 1,
      '5': 9,
      '10': 'sequencenumber'
    },
  ],
};

/// Descriptor for `PublishBatchResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishBatchResultEntryDescriptor = $convert.base64Decode(
    'ChdQdWJsaXNoQmF0Y2hSZXN1bHRFbnRyeRISCgJpZBiB8qK3ASABKAlSAmlkEiAKCW1lc3NhZ2'
    'VpZBiq5vSrASABKAlSCW1lc3NhZ2VpZBIpCg5zZXF1ZW5jZW51bWJlchiamuMuIAEoCVIOc2Vx'
    'dWVuY2VudW1iZXI=');

@$core.Deprecated('Use publishInputDescriptor instead')
const PublishInput$json = {
  '1': 'PublishInput',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'messageattributes',
      '3': 56443766,
      '4': 3,
      '5': 11,
      '6': '.sns.PublishInput.MessageattributesEntry',
      '10': 'messageattributes'
    },
    {
      '1': 'messagededuplicationid',
      '3': 379560665,
      '4': 1,
      '5': 9,
      '10': 'messagededuplicationid'
    },
    {
      '1': 'messagegroupid',
      '3': 419537435,
      '4': 1,
      '5': 9,
      '10': 'messagegroupid'
    },
    {
      '1': 'messagestructure',
      '3': 402672330,
      '4': 1,
      '5': 9,
      '10': 'messagestructure'
    },
    {'1': 'phonenumber', '3': 379600239, '4': 1, '5': 9, '10': 'phonenumber'},
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
    {'1': 'targetarn', '3': 217664144, '4': 1, '5': 9, '10': 'targetarn'},
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
  '3': [PublishInput_MessageattributesEntry$json],
};

@$core.Deprecated('Use publishInputDescriptor instead')
const PublishInput_MessageattributesEntry$json = {
  '1': 'MessageattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.sns.MessageAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `PublishInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishInputDescriptor = $convert.base64Decode(
    'CgxQdWJsaXNoSW5wdXQSGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZRJZChFtZXNzYWdlYX'
    'R0cmlidXRlcxj2hvUaIAMoCzIoLnNucy5QdWJsaXNoSW5wdXQuTWVzc2FnZWF0dHJpYnV0ZXNF'
    'bnRyeVIRbWVzc2FnZWF0dHJpYnV0ZXMSOgoWbWVzc2FnZWRlZHVwbGljYXRpb25pZBjZxf60AS'
    'ABKAlSFm1lc3NhZ2VkZWR1cGxpY2F0aW9uaWQSKgoObWVzc2FnZWdyb3VwaWQYm8SGyAEgASgJ'
    'Ug5tZXNzYWdlZ3JvdXBpZBIuChBtZXNzYWdlc3RydWN0dXJlGMqVgcABIAEoCVIQbWVzc2FnZX'
    'N0cnVjdHVyZRIkCgtwaG9uZW51bWJlchjv+oC1ASABKAlSC3Bob25lbnVtYmVyEhsKB3N1Ympl'
    'Y3QY8MnkAyABKAlSB3N1YmplY3QSHwoJdGFyZ2V0YXJuGJCV5WcgASgJUgl0YXJnZXRhcm4SHQ'
    'oIdG9waWNhcm4YnPTODiABKAlSCHRvcGljYXJuGmAKFk1lc3NhZ2VhdHRyaWJ1dGVzRW50cnkS'
    'EAoDa2V5GAEgASgJUgNrZXkSMAoFdmFsdWUYAiABKAsyGi5zbnMuTWVzc2FnZUF0dHJpYnV0ZV'
    'ZhbHVlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use publishResponseDescriptor instead')
const PublishResponse$json = {
  '1': 'PublishResponse',
  '2': [
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
    {
      '1': 'sequencenumber',
      '3': 98094362,
      '4': 1,
      '5': 9,
      '10': 'sequencenumber'
    },
  ],
};

/// Descriptor for `PublishResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List publishResponseDescriptor = $convert.base64Decode(
    'Cg9QdWJsaXNoUmVzcG9uc2USIAoJbWVzc2FnZWlkGKrm9KsBIAEoCVIJbWVzc2FnZWlkEikKDn'
    'NlcXVlbmNlbnVtYmVyGJqa4y4gASgJUg5zZXF1ZW5jZW51bWJlcg==');

@$core.Deprecated('Use putDataProtectionPolicyInputDescriptor instead')
const PutDataProtectionPolicyInput$json = {
  '1': 'PutDataProtectionPolicyInput',
  '2': [
    {
      '1': 'dataprotectionpolicy',
      '3': 519444573,
      '4': 1,
      '5': 9,
      '10': 'dataprotectionpolicy'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `PutDataProtectionPolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDataProtectionPolicyInputDescriptor =
    $convert.base64Decode(
        'ChxQdXREYXRhUHJvdGVjdGlvblBvbGljeUlucHV0EjYKFGRhdGFwcm90ZWN0aW9ucG9saWN5GN'
        '2w2PcBIAEoCVIUZGF0YXByb3RlY3Rpb25wb2xpY3kSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJ'
        'UgtyZXNvdXJjZWFybg==');

@$core.Deprecated('Use removePermissionInputDescriptor instead')
const RemovePermissionInput$json = {
  '1': 'RemovePermissionInput',
  '2': [
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `RemovePermissionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removePermissionInputDescriptor = $convert.base64Decode(
    'ChVSZW1vdmVQZXJtaXNzaW9uSW5wdXQSGAoFbGFiZWwYnuWz9gEgASgJUgVsYWJlbBIdCgh0b3'
    'BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use replayLimitExceededExceptionDescriptor instead')
const ReplayLimitExceededException$json = {
  '1': 'ReplayLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ReplayLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replayLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChxSZXBsYXlMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use sMSSandboxPhoneNumberDescriptor instead')
const SMSSandboxPhoneNumber$json = {
  '1': 'SMSSandboxPhoneNumber',
  '2': [
    {'1': 'phonenumber', '3': 379600239, '4': 1, '5': 9, '10': 'phonenumber'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.sns.SMSSandboxPhoneNumberVerificationStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `SMSSandboxPhoneNumber`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sMSSandboxPhoneNumberDescriptor = $convert.base64Decode(
    'ChVTTVNTYW5kYm94UGhvbmVOdW1iZXISJAoLcGhvbmVudW1iZXIY7/qAtQEgASgJUgtwaG9uZW'
    '51bWJlchJHCgZzdGF0dXMYkOT7AiABKA4yLC5zbnMuU01TU2FuZGJveFBob25lTnVtYmVyVmVy'
    'aWZpY2F0aW9uU3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use setEndpointAttributesInputDescriptor instead')
const SetEndpointAttributesInput$json = {
  '1': 'SetEndpointAttributesInput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.SetEndpointAttributesInput.AttributesEntry',
      '10': 'attributes'
    },
    {'1': 'endpointarn', '3': 32228660, '4': 1, '5': 9, '10': 'endpointarn'},
  ],
  '3': [SetEndpointAttributesInput_AttributesEntry$json],
};

@$core.Deprecated('Use setEndpointAttributesInputDescriptor instead')
const SetEndpointAttributesInput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SetEndpointAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setEndpointAttributesInputDescriptor = $convert.base64Decode(
    'ChpTZXRFbmRwb2ludEF0dHJpYnV0ZXNJbnB1dBJSCgphdHRyaWJ1dGVzGLWp+2MgAygLMi8uc2'
    '5zLlNldEVuZHBvaW50QXR0cmlidXRlc0lucHV0LkF0dHJpYnV0ZXNFbnRyeVIKYXR0cmlidXRl'
    'cxIjCgtlbmRwb2ludGFybhi0iq8PIAEoCVILZW5kcG9pbnRhcm4aPQoPQXR0cmlidXRlc0VudH'
    'J5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use setPlatformApplicationAttributesInputDescriptor instead')
const SetPlatformApplicationAttributesInput$json = {
  '1': 'SetPlatformApplicationAttributesInput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.SetPlatformApplicationAttributesInput.AttributesEntry',
      '10': 'attributes'
    },
    {
      '1': 'platformapplicationarn',
      '3': 241250568,
      '4': 1,
      '5': 9,
      '10': 'platformapplicationarn'
    },
  ],
  '3': [SetPlatformApplicationAttributesInput_AttributesEntry$json],
};

@$core.Deprecated('Use setPlatformApplicationAttributesInputDescriptor instead')
const SetPlatformApplicationAttributesInput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SetPlatformApplicationAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setPlatformApplicationAttributesInputDescriptor =
    $convert.base64Decode(
        'CiVTZXRQbGF0Zm9ybUFwcGxpY2F0aW9uQXR0cmlidXRlc0lucHV0El0KCmF0dHJpYnV0ZXMYta'
        'n7YyADKAsyOi5zbnMuU2V0UGxhdGZvcm1BcHBsaWNhdGlvbkF0dHJpYnV0ZXNJbnB1dC5BdHRy'
        'aWJ1dGVzRW50cnlSCmF0dHJpYnV0ZXMSOQoWcGxhdGZvcm1hcHBsaWNhdGlvbmFybhiI4oRzIA'
        'EoCVIWcGxhdGZvcm1hcHBsaWNhdGlvbmFybho9Cg9BdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEg'
        'ASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use setSMSAttributesInputDescriptor instead')
const SetSMSAttributesInput$json = {
  '1': 'SetSMSAttributesInput',
  '2': [
    {
      '1': 'attributes',
      '3': 33545109,
      '4': 3,
      '5': 11,
      '6': '.sns.SetSMSAttributesInput.AttributesEntry',
      '10': 'attributes'
    },
  ],
  '3': [SetSMSAttributesInput_AttributesEntry$json],
};

@$core.Deprecated('Use setSMSAttributesInputDescriptor instead')
const SetSMSAttributesInput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SetSMSAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setSMSAttributesInputDescriptor = $convert.base64Decode(
    'ChVTZXRTTVNBdHRyaWJ1dGVzSW5wdXQSTQoKYXR0cmlidXRlcxiVt/8PIAMoCzIqLnNucy5TZX'
    'RTTVNBdHRyaWJ1dGVzSW5wdXQuQXR0cmlidXRlc0VudHJ5UgphdHRyaWJ1dGVzGj0KD0F0dHJp'
    'YnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use setSMSAttributesResponseDescriptor instead')
const SetSMSAttributesResponse$json = {
  '1': 'SetSMSAttributesResponse',
};

/// Descriptor for `SetSMSAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setSMSAttributesResponseDescriptor =
    $convert.base64Decode('ChhTZXRTTVNBdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated('Use setSubscriptionAttributesInputDescriptor instead')
const SetSubscriptionAttributesInput$json = {
  '1': 'SetSubscriptionAttributesInput',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'attributevalue',
      '3': 96769221,
      '4': 1,
      '5': 9,
      '10': 'attributevalue'
    },
    {
      '1': 'subscriptionarn',
      '3': 279547820,
      '4': 1,
      '5': 9,
      '10': 'subscriptionarn'
    },
  ],
};

/// Descriptor for `SetSubscriptionAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setSubscriptionAttributesInputDescriptor =
    $convert.base64Decode(
        'Ch5TZXRTdWJzY3JpcHRpb25BdHRyaWJ1dGVzSW5wdXQSKAoNYXR0cmlidXRlbmFtZRitlZioAS'
        'ABKAlSDWF0dHJpYnV0ZW5hbWUSKQoOYXR0cmlidXRldmFsdWUYxamSLiABKAlSDmF0dHJpYnV0'
        'ZXZhbHVlEiwKD3N1YnNjcmlwdGlvbmFybhisn6aFASABKAlSD3N1YnNjcmlwdGlvbmFybg==');

@$core.Deprecated('Use setTopicAttributesInputDescriptor instead')
const SetTopicAttributesInput$json = {
  '1': 'SetTopicAttributesInput',
  '2': [
    {
      '1': 'attributename',
      '3': 352717485,
      '4': 1,
      '5': 9,
      '10': 'attributename'
    },
    {
      '1': 'attributevalue',
      '3': 96769221,
      '4': 1,
      '5': 9,
      '10': 'attributevalue'
    },
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `SetTopicAttributesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setTopicAttributesInputDescriptor = $convert.base64Decode(
    'ChdTZXRUb3BpY0F0dHJpYnV0ZXNJbnB1dBIoCg1hdHRyaWJ1dGVuYW1lGK2VmKgBIAEoCVINYX'
    'R0cmlidXRlbmFtZRIpCg5hdHRyaWJ1dGV2YWx1ZRjFqZIuIAEoCVIOYXR0cmlidXRldmFsdWUS'
    'HQoIdG9waWNhcm4YnPTODiABKAlSCHRvcGljYXJu');

@$core.Deprecated('Use staleTagExceptionDescriptor instead')
const StaleTagException$json = {
  '1': 'StaleTagException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `StaleTagException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List staleTagExceptionDescriptor = $convert.base64Decode(
    'ChFTdGFsZVRhZ0V4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use subscribeInputDescriptor instead')
const SubscribeInput$json = {
  '1': 'SubscribeInput',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sns.SubscribeInput.AttributesEntry',
      '10': 'attributes'
    },
    {'1': 'endpoint', '3': 132634269, '4': 1, '5': 9, '10': 'endpoint'},
    {'1': 'protocol', '3': 173534166, '4': 1, '5': 9, '10': 'protocol'},
    {
      '1': 'returnsubscriptionarn',
      '3': 78754574,
      '4': 1,
      '5': 8,
      '10': 'returnsubscriptionarn'
    },
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
  '3': [SubscribeInput_AttributesEntry$json],
};

@$core.Deprecated('Use subscribeInputDescriptor instead')
const SubscribeInput_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SubscribeInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeInputDescriptor = $convert.base64Decode(
    'Cg5TdWJzY3JpYmVJbnB1dBJGCgphdHRyaWJ1dGVzGLWp+2MgAygLMiMuc25zLlN1YnNjcmliZU'
    'lucHV0LkF0dHJpYnV0ZXNFbnRyeVIKYXR0cmlidXRlcxIdCghlbmRwb2ludBidrZ8/IAEoCVII'
    'ZW5kcG9pbnQSHQoIcHJvdG9jb2wY1tffUiABKAlSCHByb3RvY29sEjcKFXJldHVybnN1YnNjcm'
    'lwdGlvbmFybhiO5sYlIAEoCFIVcmV0dXJuc3Vic2NyaXB0aW9uYXJuEh0KCHRvcGljYXJuGJz0'
    'zg4gASgJUgh0b3BpY2Fybho9Cg9BdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFA'
    'oFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use subscribeResponseDescriptor instead')
const SubscribeResponse$json = {
  '1': 'SubscribeResponse',
  '2': [
    {
      '1': 'subscriptionarn',
      '3': 279547820,
      '4': 1,
      '5': 9,
      '10': 'subscriptionarn'
    },
  ],
};

/// Descriptor for `SubscribeResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeResponseDescriptor = $convert.base64Decode(
    'ChFTdWJzY3JpYmVSZXNwb25zZRIsCg9zdWJzY3JpcHRpb25hcm4YrJ+mhQEgASgJUg9zdWJzY3'
    'JpcHRpb25hcm4=');

@$core.Deprecated('Use subscriptionDescriptor instead')
const Subscription$json = {
  '1': 'Subscription',
  '2': [
    {'1': 'endpoint', '3': 132634269, '4': 1, '5': 9, '10': 'endpoint'},
    {'1': 'owner', '3': 455261813, '4': 1, '5': 9, '10': 'owner'},
    {'1': 'protocol', '3': 173534166, '4': 1, '5': 9, '10': 'protocol'},
    {
      '1': 'subscriptionarn',
      '3': 279547820,
      '4': 1,
      '5': 9,
      '10': 'subscriptionarn'
    },
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `Subscription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscriptionDescriptor = $convert.base64Decode(
    'CgxTdWJzY3JpcHRpb24SHQoIZW5kcG9pbnQYna2fPyABKAlSCGVuZHBvaW50EhgKBW93bmVyGP'
    'X8itkBIAEoCVIFb3duZXISHQoIcHJvdG9jb2wY1tffUiABKAlSCHByb3RvY29sEiwKD3N1YnNj'
    'cmlwdGlvbmFybhisn6aFASABKAlSD3N1YnNjcmlwdGlvbmFybhIdCgh0b3BpY2Fybhic9M4OIA'
    'EoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use subscriptionLimitExceededExceptionDescriptor instead')
const SubscriptionLimitExceededException$json = {
  '1': 'SubscriptionLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SubscriptionLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscriptionLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiJTdWJzY3JpcHRpb25MaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use tagDescriptor instead')
const Tag$json = {
  '1': 'Tag',
  '2': [
    {'1': 'key', '3': 219859213, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `Tag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagDescriptor = $convert.base64Decode(
    'CgNUYWcSEwoDa2V5GI2S62ggASgJUgNrZXkSGAoFdmFsdWUY6/KfigEgASgJUgV2YWx1ZQ==');

@$core.Deprecated('Use tagLimitExceededExceptionDescriptor instead')
const TagLimitExceededException$json = {
  '1': 'TagLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TagLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChlUYWdMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use tagPolicyExceptionDescriptor instead')
const TagPolicyException$json = {
  '1': 'TagPolicyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TagPolicyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagPolicyExceptionDescriptor =
    $convert.base64Decode(
        'ChJUYWdQb2xpY3lFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sns.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJjZW'
    'FybhIgCgR0YWdzGMHB9rUBIAMoCzIILnNucy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use tagResourceResponseDescriptor instead')
const TagResourceResponse$json = {
  '1': 'TagResourceResponse',
};

/// Descriptor for `TagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceResponseDescriptor =
    $convert.base64Decode('ChNUYWdSZXNvdXJjZVJlc3BvbnNl');

@$core.Deprecated('Use throttledExceptionDescriptor instead')
const ThrottledException$json = {
  '1': 'ThrottledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ThrottledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttledExceptionDescriptor =
    $convert.base64Decode(
        'ChJUaHJvdHRsZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use tooManyEntriesInBatchRequestExceptionDescriptor instead')
const TooManyEntriesInBatchRequestException$json = {
  '1': 'TooManyEntriesInBatchRequestException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyEntriesInBatchRequestException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyEntriesInBatchRequestExceptionDescriptor =
    $convert.base64Decode(
        'CiVUb29NYW55RW50cmllc0luQmF0Y2hSZXF1ZXN0RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use topicDescriptor instead')
const Topic$json = {
  '1': 'Topic',
  '2': [
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `Topic`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List topicDescriptor = $convert
    .base64Decode('CgVUb3BpYxIdCgh0b3BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use topicLimitExceededExceptionDescriptor instead')
const TopicLimitExceededException$json = {
  '1': 'TopicLimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TopicLimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List topicLimitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChtUb3BpY0xpbWl0RXhjZWVkZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use unsubscribeInputDescriptor instead')
const UnsubscribeInput$json = {
  '1': 'UnsubscribeInput',
  '2': [
    {
      '1': 'subscriptionarn',
      '3': 279547820,
      '4': 1,
      '5': 9,
      '10': 'subscriptionarn'
    },
  ],
};

/// Descriptor for `UnsubscribeInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsubscribeInputDescriptor = $convert.base64Decode(
    'ChBVbnN1YnNjcmliZUlucHV0EiwKD3N1YnNjcmlwdGlvbmFybhisn6aFASABKAlSD3N1YnNjcm'
    'lwdGlvbmFybg==');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc291cm'
    'NlYXJuEhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use untagResourceResponseDescriptor instead')
const UntagResourceResponse$json = {
  '1': 'UntagResourceResponse',
};

/// Descriptor for `UntagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceResponseDescriptor =
    $convert.base64Decode('ChVVbnRhZ1Jlc291cmNlUmVzcG9uc2U=');

@$core.Deprecated('Use userErrorExceptionDescriptor instead')
const UserErrorException$json = {
  '1': 'UserErrorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UserErrorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List userErrorExceptionDescriptor =
    $convert.base64Decode(
        'ChJVc2VyRXJyb3JFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use validationExceptionDescriptor instead')
const ValidationException$json = {
  '1': 'ValidationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationExceptionDescriptor =
    $convert.base64Decode(
        'ChNWYWxpZGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use verificationExceptionDescriptor instead')
const VerificationException$json = {
  '1': 'VerificationException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
  ],
};

/// Descriptor for `VerificationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verificationExceptionDescriptor = $convert.base64Decode(
    'ChVWZXJpZmljYXRpb25FeGNlcHRpb24SGwoHbWVzc2FnZRiFs7twIAEoCVIHbWVzc2FnZRIZCg'
    'ZzdGF0dXMYkOT7AiABKAlSBnN0YXR1cw==');

@$core.Deprecated('Use verifySMSSandboxPhoneNumberInputDescriptor instead')
const VerifySMSSandboxPhoneNumberInput$json = {
  '1': 'VerifySMSSandboxPhoneNumberInput',
  '2': [
    {
      '1': 'onetimepassword',
      '3': 417020900,
      '4': 1,
      '5': 9,
      '10': 'onetimepassword'
    },
    {'1': 'phonenumber', '3': 379600239, '4': 1, '5': 9, '10': 'phonenumber'},
  ],
};

/// Descriptor for `VerifySMSSandboxPhoneNumberInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifySMSSandboxPhoneNumberInputDescriptor =
    $convert.base64Decode(
        'CiBWZXJpZnlTTVNTYW5kYm94UGhvbmVOdW1iZXJJbnB1dBIsCg9vbmV0aW1lcGFzc3dvcmQY5P'
        'fsxgEgASgJUg9vbmV0aW1lcGFzc3dvcmQSJAoLcGhvbmVudW1iZXIY7/qAtQEgASgJUgtwaG9u'
        'ZW51bWJlcg==');

@$core.Deprecated('Use verifySMSSandboxPhoneNumberResultDescriptor instead')
const VerifySMSSandboxPhoneNumberResult$json = {
  '1': 'VerifySMSSandboxPhoneNumberResult',
};

/// Descriptor for `VerifySMSSandboxPhoneNumberResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifySMSSandboxPhoneNumberResultDescriptor =
    $convert.base64Decode('CiFWZXJpZnlTTVNTYW5kYm94UGhvbmVOdW1iZXJSZXN1bHQ=');
