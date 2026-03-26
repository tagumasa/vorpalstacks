// This is a generated file - do not edit.
//
// Generated from cognito_identity.proto.

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

@$core.Deprecated('Use ambiguousRoleResolutionTypeDescriptor instead')
const AmbiguousRoleResolutionType$json = {
  '1': 'AmbiguousRoleResolutionType',
  '2': [
    {'1': 'AMBIGUOUS_ROLE_RESOLUTION_TYPE_AUTHENTICATED_ROLE', '2': 0},
    {'1': 'AMBIGUOUS_ROLE_RESOLUTION_TYPE_DENY', '2': 1},
  ],
};

/// Descriptor for `AmbiguousRoleResolutionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List ambiguousRoleResolutionTypeDescriptor =
    $convert.base64Decode(
        'ChtBbWJpZ3VvdXNSb2xlUmVzb2x1dGlvblR5cGUSNQoxQU1CSUdVT1VTX1JPTEVfUkVTT0xVVE'
        'lPTl9UWVBFX0FVVEhFTlRJQ0FURURfUk9MRRAAEicKI0FNQklHVU9VU19ST0xFX1JFU09MVVRJ'
        'T05fVFlQRV9ERU5ZEAE=');

@$core.Deprecated('Use errorCodeDescriptor instead')
const ErrorCode$json = {
  '1': 'ErrorCode',
  '2': [
    {'1': 'ERROR_CODE_ACCESS_DENIED', '2': 0},
    {'1': 'ERROR_CODE_INTERNAL_SERVER_ERROR', '2': 1},
  ],
};

/// Descriptor for `ErrorCode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List errorCodeDescriptor = $convert.base64Decode(
    'CglFcnJvckNvZGUSHAoYRVJST1JfQ09ERV9BQ0NFU1NfREVOSUVEEAASJAogRVJST1JfQ09ERV'
    '9JTlRFUk5BTF9TRVJWRVJfRVJST1IQAQ==');

@$core.Deprecated('Use mappingRuleMatchTypeDescriptor instead')
const MappingRuleMatchType$json = {
  '1': 'MappingRuleMatchType',
  '2': [
    {'1': 'MAPPING_RULE_MATCH_TYPE_CONTAINS', '2': 0},
    {'1': 'MAPPING_RULE_MATCH_TYPE_EQUALS', '2': 1},
    {'1': 'MAPPING_RULE_MATCH_TYPE_STARTS_WITH', '2': 2},
    {'1': 'MAPPING_RULE_MATCH_TYPE_NOT_EQUAL', '2': 3},
  ],
};

/// Descriptor for `MappingRuleMatchType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mappingRuleMatchTypeDescriptor = $convert.base64Decode(
    'ChRNYXBwaW5nUnVsZU1hdGNoVHlwZRIkCiBNQVBQSU5HX1JVTEVfTUFUQ0hfVFlQRV9DT05UQU'
    'lOUxAAEiIKHk1BUFBJTkdfUlVMRV9NQVRDSF9UWVBFX0VRVUFMUxABEicKI01BUFBJTkdfUlVM'
    'RV9NQVRDSF9UWVBFX1NUQVJUU19XSVRIEAISJQohTUFQUElOR19SVUxFX01BVENIX1RZUEVfTk'
    '9UX0VRVUFMEAM=');

@$core.Deprecated('Use roleMappingTypeDescriptor instead')
const RoleMappingType$json = {
  '1': 'RoleMappingType',
  '2': [
    {'1': 'ROLE_MAPPING_TYPE_RULES', '2': 0},
    {'1': 'ROLE_MAPPING_TYPE_TOKEN', '2': 1},
  ],
};

/// Descriptor for `RoleMappingType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List roleMappingTypeDescriptor = $convert.base64Decode(
    'Cg9Sb2xlTWFwcGluZ1R5cGUSGwoXUk9MRV9NQVBQSU5HX1RZUEVfUlVMRVMQABIbChdST0xFX0'
    '1BUFBJTkdfVFlQRV9UT0tFThAB');

@$core.Deprecated('Use cognitoIdentityProviderDescriptor instead')
const CognitoIdentityProvider$json = {
  '1': 'CognitoIdentityProvider',
  '2': [
    {'1': 'clientid', '3': 448902180, '4': 1, '5': 9, '10': 'clientid'},
    {'1': 'providername', '3': 485101816, '4': 1, '5': 9, '10': 'providername'},
    {
      '1': 'serversidetokencheck',
      '3': 291427543,
      '4': 1,
      '5': 8,
      '10': 'serversidetokencheck'
    },
  ],
};

/// Descriptor for `CognitoIdentityProvider`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cognitoIdentityProviderDescriptor = $convert.base64Decode(
    'ChdDb2duaXRvSWRlbnRpdHlQcm92aWRlchIeCghjbGllbnRpZBik6IbWASABKAlSCGNsaWVudG'
    'lkEiYKDHByb3ZpZGVybmFtZRj4oajnASABKAlSDHByb3ZpZGVybmFtZRI2ChRzZXJ2ZXJzaWRl'
    'dG9rZW5jaGVjaxjXqfuKASABKAhSFHNlcnZlcnNpZGV0b2tlbmNoZWNr');

@$core.Deprecated('Use concurrentModificationExceptionDescriptor instead')
const ConcurrentModificationException$json = {
  '1': 'ConcurrentModificationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConcurrentModificationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List concurrentModificationExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Db25jdXJyZW50TW9kaWZpY2F0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use createIdentityPoolInputDescriptor instead')
const CreateIdentityPoolInput$json = {
  '1': 'CreateIdentityPoolInput',
  '2': [
    {
      '1': 'allowclassicflow',
      '3': 101299523,
      '4': 1,
      '5': 8,
      '10': 'allowclassicflow'
    },
    {
      '1': 'allowunauthenticatedidentities',
      '3': 290546287,
      '4': 1,
      '5': 8,
      '10': 'allowunauthenticatedidentities'
    },
    {
      '1': 'cognitoidentityproviders',
      '3': 109610365,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.CognitoIdentityProvider',
      '10': 'cognitoidentityproviders'
    },
    {
      '1': 'developerprovidername',
      '3': 517589792,
      '4': 1,
      '5': 9,
      '10': 'developerprovidername'
    },
    {
      '1': 'identitypoolname',
      '3': 41275091,
      '4': 1,
      '5': 9,
      '10': 'identitypoolname'
    },
    {
      '1': 'identitypooltags',
      '3': 305630405,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.CreateIdentityPoolInput.IdentitypooltagsEntry',
      '10': 'identitypooltags'
    },
    {
      '1': 'openidconnectproviderarns',
      '3': 137395556,
      '4': 3,
      '5': 9,
      '10': 'openidconnectproviderarns'
    },
    {
      '1': 'samlproviderarns',
      '3': 24907670,
      '4': 3,
      '5': 9,
      '10': 'samlproviderarns'
    },
    {
      '1': 'supportedloginproviders',
      '3': 125752617,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_identity.CreateIdentityPoolInput.SupportedloginprovidersEntry',
      '10': 'supportedloginproviders'
    },
  ],
  '3': [
    CreateIdentityPoolInput_IdentitypooltagsEntry$json,
    CreateIdentityPoolInput_SupportedloginprovidersEntry$json
  ],
};

@$core.Deprecated('Use createIdentityPoolInputDescriptor instead')
const CreateIdentityPoolInput_IdentitypooltagsEntry$json = {
  '1': 'IdentitypooltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use createIdentityPoolInputDescriptor instead')
const CreateIdentityPoolInput_SupportedloginprovidersEntry$json = {
  '1': 'SupportedloginprovidersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateIdentityPoolInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createIdentityPoolInputDescriptor = $convert.base64Decode(
    'ChdDcmVhdGVJZGVudGl0eVBvb2xJbnB1dBItChBhbGxvd2NsYXNzaWNmbG93GMPqpjAgASgIUh'
    'BhbGxvd2NsYXNzaWNmbG93EkoKHmFsbG93dW5hdXRoZW50aWNhdGVkaWRlbnRpdGllcxjvxMWK'
    'ASABKAhSHmFsbG93dW5hdXRoZW50aWNhdGVkaWRlbnRpdGllcxJoChhjb2duaXRvaWRlbnRpdH'
    'lwcm92aWRlcnMY/YqiNCADKAsyKS5jb2duaXRvX2lkZW50aXR5LkNvZ25pdG9JZGVudGl0eVBy'
    'b3ZpZGVyUhhjb2duaXRvaWRlbnRpdHlwcm92aWRlcnMSOAoVZGV2ZWxvcGVycHJvdmlkZXJuYW'
    '1lGKCW5/YBIAEoCVIVZGV2ZWxvcGVycHJvdmlkZXJuYW1lEi0KEGlkZW50aXR5cG9vbG5hbWUY'
    '053XEyABKAlSEGlkZW50aXR5cG9vbG5hbWUSbwoQaWRlbnRpdHlwb29sdGFncxjFmd6RASADKA'
    'syPy5jb2duaXRvX2lkZW50aXR5LkNyZWF0ZUlkZW50aXR5UG9vbElucHV0LklkZW50aXR5cG9v'
    'bHRhZ3NFbnRyeVIQaWRlbnRpdHlwb29sdGFncxI/ChlvcGVuaWRjb25uZWN0cHJvdmlkZXJhcm'
    '5zGOT6wUEgAygJUhlvcGVuaWRjb25uZWN0cHJvdmlkZXJhcm5zEi0KEHNhbWxwcm92aWRlcmFy'
    'bnMYlp/wCyADKAlSEHNhbWxwcm92aWRlcmFybnMSgwEKF3N1cHBvcnRlZGxvZ2lucHJvdmlkZX'
    'JzGKmq+zsgAygLMkYuY29nbml0b19pZGVudGl0eS5DcmVhdGVJZGVudGl0eVBvb2xJbnB1dC5T'
    'dXBwb3J0ZWRsb2dpbnByb3ZpZGVyc0VudHJ5UhdzdXBwb3J0ZWRsb2dpbnByb3ZpZGVycxpDCh'
    'VJZGVudGl0eXBvb2x0YWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlS'
    'BXZhbHVlOgI4ARpKChxTdXBwb3J0ZWRsb2dpbnByb3ZpZGVyc0VudHJ5EhAKA2tleRgBIAEoCV'
    'IDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use credentialsDescriptor instead')
const Credentials$json = {
  '1': 'Credentials',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {'1': 'secretkey', '3': 465028505, '4': 1, '5': 9, '10': 'secretkey'},
    {'1': 'sessiontoken', '3': 211161069, '4': 1, '5': 9, '10': 'sessiontoken'},
  ],
};

/// Descriptor for `Credentials`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List credentialsDescriptor = $convert.base64Decode(
    'CgtDcmVkZW50aWFscxIkCgthY2Nlc3NrZXlpZBigt7fYASABKAlSC2FjY2Vzc2tleWlkEiEKCm'
    'V4cGlyYXRpb24YiamfdSABKAlSCmV4cGlyYXRpb24SIAoJc2VjcmV0a2V5GJmL390BIAEoCVIJ'
    'c2VjcmV0a2V5EiUKDHNlc3Npb250b2tlbhjtn9hkIAEoCVIMc2Vzc2lvbnRva2Vu');

@$core.Deprecated('Use deleteIdentitiesInputDescriptor instead')
const DeleteIdentitiesInput$json = {
  '1': 'DeleteIdentitiesInput',
  '2': [
    {
      '1': 'identityidstodelete',
      '3': 340357990,
      '4': 3,
      '5': 9,
      '10': 'identityidstodelete'
    },
  ],
};

/// Descriptor for `DeleteIdentitiesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIdentitiesInputDescriptor = $convert.base64Decode(
    'ChVEZWxldGVJZGVudGl0aWVzSW5wdXQSNAoTaWRlbnRpdHlpZHN0b2RlbGV0ZRjm5qWiASADKA'
    'lSE2lkZW50aXR5aWRzdG9kZWxldGU=');

@$core.Deprecated('Use deleteIdentitiesResponseDescriptor instead')
const DeleteIdentitiesResponse$json = {
  '1': 'DeleteIdentitiesResponse',
  '2': [
    {
      '1': 'unprocessedidentityids',
      '3': 532881071,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.UnprocessedIdentityId',
      '10': 'unprocessedidentityids'
    },
  ],
};

/// Descriptor for `DeleteIdentitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIdentitiesResponseDescriptor = $convert.base64Decode(
    'ChhEZWxldGVJZGVudGl0aWVzUmVzcG9uc2USYwoWdW5wcm9jZXNzZWRpZGVudGl0eWlkcxivvY'
    'z+ASADKAsyJy5jb2duaXRvX2lkZW50aXR5LlVucHJvY2Vzc2VkSWRlbnRpdHlJZFIWdW5wcm9j'
    'ZXNzZWRpZGVudGl0eWlkcw==');

@$core.Deprecated('Use deleteIdentityPoolInputDescriptor instead')
const DeleteIdentityPoolInput$json = {
  '1': 'DeleteIdentityPoolInput',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
  ],
};

/// Descriptor for `DeleteIdentityPoolInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteIdentityPoolInputDescriptor =
    $convert.base64Decode(
        'ChdEZWxldGVJZGVudGl0eVBvb2xJbnB1dBIpCg5pZGVudGl0eXBvb2xpZBj9/bQLIAEoCVIOaW'
        'RlbnRpdHlwb29saWQ=');

@$core.Deprecated('Use describeIdentityInputDescriptor instead')
const DescribeIdentityInput$json = {
  '1': 'DescribeIdentityInput',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
  ],
};

/// Descriptor for `DescribeIdentityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeIdentityInputDescriptor = $convert.base64Decode(
    'ChVEZXNjcmliZUlkZW50aXR5SW5wdXQSIQoKaWRlbnRpdHlpZBjX09VvIAEoCVIKaWRlbnRpdH'
    'lpZA==');

@$core.Deprecated('Use describeIdentityPoolInputDescriptor instead')
const DescribeIdentityPoolInput$json = {
  '1': 'DescribeIdentityPoolInput',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
  ],
};

/// Descriptor for `DescribeIdentityPoolInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeIdentityPoolInputDescriptor =
    $convert.base64Decode(
        'ChlEZXNjcmliZUlkZW50aXR5UG9vbElucHV0EikKDmlkZW50aXR5cG9vbGlkGP39tAsgASgJUg'
        '5pZGVudGl0eXBvb2xpZA==');

@$core
    .Deprecated('Use developerUserAlreadyRegisteredExceptionDescriptor instead')
const DeveloperUserAlreadyRegisteredException$json = {
  '1': 'DeveloperUserAlreadyRegisteredException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DeveloperUserAlreadyRegisteredException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List developerUserAlreadyRegisteredExceptionDescriptor =
    $convert.base64Decode(
        'CidEZXZlbG9wZXJVc2VyQWxyZWFkeVJlZ2lzdGVyZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkc'
        'gnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use externalServiceExceptionDescriptor instead')
const ExternalServiceException$json = {
  '1': 'ExternalServiceException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExternalServiceException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List externalServiceExceptionDescriptor =
    $convert.base64Decode(
        'ChhFeHRlcm5hbFNlcnZpY2VFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use getCredentialsForIdentityInputDescriptor instead')
const GetCredentialsForIdentityInput$json = {
  '1': 'GetCredentialsForIdentityInput',
  '2': [
    {
      '1': 'customrolearn',
      '3': 146997938,
      '4': 1,
      '5': 9,
      '10': 'customrolearn'
    },
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {
      '1': 'logins',
      '3': 109702772,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.GetCredentialsForIdentityInput.LoginsEntry',
      '10': 'logins'
    },
  ],
  '3': [GetCredentialsForIdentityInput_LoginsEntry$json],
};

@$core.Deprecated('Use getCredentialsForIdentityInputDescriptor instead')
const GetCredentialsForIdentityInput_LoginsEntry$json = {
  '1': 'LoginsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetCredentialsForIdentityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCredentialsForIdentityInputDescriptor = $convert.base64Decode(
    'Ch5HZXRDcmVkZW50aWFsc0ZvcklkZW50aXR5SW5wdXQSJwoNY3VzdG9tcm9sZWFybhiyhYxGIA'
    'EoCVINY3VzdG9tcm9sZWFybhIhCgppZGVudGl0eWlkGNfT1W8gASgJUgppZGVudGl0eWlkElcK'
    'BmxvZ2lucxj03Kc0IAMoCzI8LmNvZ25pdG9faWRlbnRpdHkuR2V0Q3JlZGVudGlhbHNGb3JJZG'
    'VudGl0eUlucHV0LkxvZ2luc0VudHJ5UgZsb2dpbnMaOQoLTG9naW5zRW50cnkSEAoDa2V5GAEg'
    'ASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getCredentialsForIdentityResponseDescriptor instead')
const GetCredentialsForIdentityResponse$json = {
  '1': 'GetCredentialsForIdentityResponse',
  '2': [
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.cognito_identity.Credentials',
      '10': 'credentials'
    },
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
  ],
};

/// Descriptor for `GetCredentialsForIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCredentialsForIdentityResponseDescriptor =
    $convert.base64Decode(
        'CiFHZXRDcmVkZW50aWFsc0ZvcklkZW50aXR5UmVzcG9uc2USQwoLY3JlZGVudGlhbHMY8pqOtg'
        'EgASgLMh0uY29nbml0b19pZGVudGl0eS5DcmVkZW50aWFsc1ILY3JlZGVudGlhbHMSIQoKaWRl'
        'bnRpdHlpZBjX09VvIAEoCVIKaWRlbnRpdHlpZA==');

@$core.Deprecated('Use getIdInputDescriptor instead')
const GetIdInput$json = {
  '1': 'GetIdInput',
  '2': [
    {'1': 'accountid', '3': 65954002, '4': 1, '5': 9, '10': 'accountid'},
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'logins',
      '3': 109702772,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.GetIdInput.LoginsEntry',
      '10': 'logins'
    },
  ],
  '3': [GetIdInput_LoginsEntry$json],
};

@$core.Deprecated('Use getIdInputDescriptor instead')
const GetIdInput_LoginsEntry$json = {
  '1': 'LoginsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetIdInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIdInputDescriptor = $convert.base64Decode(
    'CgpHZXRJZElucHV0Eh8KCWFjY291bnRpZBjSwbkfIAEoCVIJYWNjb3VudGlkEikKDmlkZW50aX'
    'R5cG9vbGlkGP39tAsgASgJUg5pZGVudGl0eXBvb2xpZBJDCgZsb2dpbnMY9NynNCADKAsyKC5j'
    'b2duaXRvX2lkZW50aXR5LkdldElkSW5wdXQuTG9naW5zRW50cnlSBmxvZ2lucxo5CgtMb2dpbn'
    'NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use getIdResponseDescriptor instead')
const GetIdResponse$json = {
  '1': 'GetIdResponse',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
  ],
};

/// Descriptor for `GetIdResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIdResponseDescriptor = $convert.base64Decode(
    'Cg1HZXRJZFJlc3BvbnNlEiEKCmlkZW50aXR5aWQY19PVbyABKAlSCmlkZW50aXR5aWQ=');

@$core.Deprecated('Use getIdentityPoolRolesInputDescriptor instead')
const GetIdentityPoolRolesInput$json = {
  '1': 'GetIdentityPoolRolesInput',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
  ],
};

/// Descriptor for `GetIdentityPoolRolesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIdentityPoolRolesInputDescriptor =
    $convert.base64Decode(
        'ChlHZXRJZGVudGl0eVBvb2xSb2xlc0lucHV0EikKDmlkZW50aXR5cG9vbGlkGP39tAsgASgJUg'
        '5pZGVudGl0eXBvb2xpZA==');

@$core.Deprecated('Use getIdentityPoolRolesResponseDescriptor instead')
const GetIdentityPoolRolesResponse$json = {
  '1': 'GetIdentityPoolRolesResponse',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'rolemappings',
      '3': 96154047,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.GetIdentityPoolRolesResponse.RolemappingsEntry',
      '10': 'rolemappings'
    },
    {
      '1': 'roles',
      '3': 511168127,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.GetIdentityPoolRolesResponse.RolesEntry',
      '10': 'roles'
    },
  ],
  '3': [
    GetIdentityPoolRolesResponse_RolemappingsEntry$json,
    GetIdentityPoolRolesResponse_RolesEntry$json
  ],
};

@$core.Deprecated('Use getIdentityPoolRolesResponseDescriptor instead')
const GetIdentityPoolRolesResponse_RolemappingsEntry$json = {
  '1': 'RolemappingsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.cognito_identity.RoleMapping',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use getIdentityPoolRolesResponseDescriptor instead')
const GetIdentityPoolRolesResponse_RolesEntry$json = {
  '1': 'RolesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetIdentityPoolRolesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getIdentityPoolRolesResponseDescriptor = $convert.base64Decode(
    'ChxHZXRJZGVudGl0eVBvb2xSb2xlc1Jlc3BvbnNlEikKDmlkZW50aXR5cG9vbGlkGP39tAsgAS'
    'gJUg5pZGVudGl0eXBvb2xpZBJnCgxyb2xlbWFwcGluZ3MYv+PsLSADKAsyQC5jb2duaXRvX2lk'
    'ZW50aXR5LkdldElkZW50aXR5UG9vbFJvbGVzUmVzcG9uc2UuUm9sZW1hcHBpbmdzRW50cnlSDH'
    'JvbGVtYXBwaW5ncxJTCgVyb2xlcxj/nN/zASADKAsyOS5jb2duaXRvX2lkZW50aXR5LkdldElk'
    'ZW50aXR5UG9vbFJvbGVzUmVzcG9uc2UuUm9sZXNFbnRyeVIFcm9sZXMaXgoRUm9sZW1hcHBpbm'
    'dzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSMwoFdmFsdWUYAiABKAsyHS5jb2duaXRvX2lkZW50'
    'aXR5LlJvbGVNYXBwaW5nUgV2YWx1ZToCOAEaOAoKUm9sZXNFbnRyeRIQCgNrZXkYASABKAlSA2'
    'tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core
    .Deprecated('Use getOpenIdTokenForDeveloperIdentityInputDescriptor instead')
const GetOpenIdTokenForDeveloperIdentityInput$json = {
  '1': 'GetOpenIdTokenForDeveloperIdentityInput',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'logins',
      '3': 109702772,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_identity.GetOpenIdTokenForDeveloperIdentityInput.LoginsEntry',
      '10': 'logins'
    },
    {
      '1': 'principaltags',
      '3': 346698229,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_identity.GetOpenIdTokenForDeveloperIdentityInput.PrincipaltagsEntry',
      '10': 'principaltags'
    },
    {
      '1': 'tokenduration',
      '3': 402772367,
      '4': 1,
      '5': 3,
      '10': 'tokenduration'
    },
  ],
  '3': [
    GetOpenIdTokenForDeveloperIdentityInput_LoginsEntry$json,
    GetOpenIdTokenForDeveloperIdentityInput_PrincipaltagsEntry$json
  ],
};

@$core
    .Deprecated('Use getOpenIdTokenForDeveloperIdentityInputDescriptor instead')
const GetOpenIdTokenForDeveloperIdentityInput_LoginsEntry$json = {
  '1': 'LoginsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core
    .Deprecated('Use getOpenIdTokenForDeveloperIdentityInputDescriptor instead')
const GetOpenIdTokenForDeveloperIdentityInput_PrincipaltagsEntry$json = {
  '1': 'PrincipaltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetOpenIdTokenForDeveloperIdentityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpenIdTokenForDeveloperIdentityInputDescriptor = $convert.base64Decode(
    'CidHZXRPcGVuSWRUb2tlbkZvckRldmVsb3BlcklkZW50aXR5SW5wdXQSIQoKaWRlbnRpdHlpZB'
    'jX09VvIAEoCVIKaWRlbnRpdHlpZBIpCg5pZGVudGl0eXBvb2xpZBj9/bQLIAEoCVIOaWRlbnRp'
    'dHlwb29saWQSYAoGbG9naW5zGPTcpzQgAygLMkUuY29nbml0b19pZGVudGl0eS5HZXRPcGVuSW'
    'RUb2tlbkZvckRldmVsb3BlcklkZW50aXR5SW5wdXQuTG9naW5zRW50cnlSBmxvZ2lucxJ2Cg1w'
    'cmluY2lwYWx0YWdzGPXjqKUBIAMoCzJMLmNvZ25pdG9faWRlbnRpdHkuR2V0T3BlbklkVG9rZW'
    '5Gb3JEZXZlbG9wZXJJZGVudGl0eUlucHV0LlByaW5jaXBhbHRhZ3NFbnRyeVINcHJpbmNpcGFs'
    'dGFncxIoCg10b2tlbmR1cmF0aW9uGI+jh8ABIAEoA1INdG9rZW5kdXJhdGlvbho5CgtMb2dpbn'
    'NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGkAKElBy'
    'aW5jaXBhbHRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdW'
    'U6AjgB');

@$core.Deprecated(
    'Use getOpenIdTokenForDeveloperIdentityResponseDescriptor instead')
const GetOpenIdTokenForDeveloperIdentityResponse$json = {
  '1': 'GetOpenIdTokenForDeveloperIdentityResponse',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {'1': 'token', '3': 439704531, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `GetOpenIdTokenForDeveloperIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getOpenIdTokenForDeveloperIdentityResponseDescriptor =
    $convert.base64Decode(
        'CipHZXRPcGVuSWRUb2tlbkZvckRldmVsb3BlcklkZW50aXR5UmVzcG9uc2USIQoKaWRlbnRpdH'
        'lpZBjX09VvIAEoCVIKaWRlbnRpdHlpZBIYCgV0b2tlbhjTt9XRASABKAlSBXRva2Vu');

@$core.Deprecated('Use getOpenIdTokenInputDescriptor instead')
const GetOpenIdTokenInput$json = {
  '1': 'GetOpenIdTokenInput',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {
      '1': 'logins',
      '3': 109702772,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.GetOpenIdTokenInput.LoginsEntry',
      '10': 'logins'
    },
  ],
  '3': [GetOpenIdTokenInput_LoginsEntry$json],
};

@$core.Deprecated('Use getOpenIdTokenInputDescriptor instead')
const GetOpenIdTokenInput_LoginsEntry$json = {
  '1': 'LoginsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetOpenIdTokenInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpenIdTokenInputDescriptor = $convert.base64Decode(
    'ChNHZXRPcGVuSWRUb2tlbklucHV0EiEKCmlkZW50aXR5aWQY19PVbyABKAlSCmlkZW50aXR5aW'
    'QSTAoGbG9naW5zGPTcpzQgAygLMjEuY29nbml0b19pZGVudGl0eS5HZXRPcGVuSWRUb2tlbklu'
    'cHV0LkxvZ2luc0VudHJ5UgZsb2dpbnMaOQoLTG9naW5zRW50cnkSEAoDa2V5GAEgASgJUgNrZX'
    'kSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getOpenIdTokenResponseDescriptor instead')
const GetOpenIdTokenResponse$json = {
  '1': 'GetOpenIdTokenResponse',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {'1': 'token', '3': 439704531, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `GetOpenIdTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getOpenIdTokenResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXRPcGVuSWRUb2tlblJlc3BvbnNlEiEKCmlkZW50aXR5aWQY19PVbyABKAlSCmlkZW50aX'
        'R5aWQSGAoFdG9rZW4Y07fV0QEgASgJUgV0b2tlbg==');

@$core.Deprecated('Use getPrincipalTagAttributeMapInputDescriptor instead')
const GetPrincipalTagAttributeMapInput$json = {
  '1': 'GetPrincipalTagAttributeMapInput',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'identityprovidername',
      '3': 419175410,
      '4': 1,
      '5': 9,
      '10': 'identityprovidername'
    },
  ],
};

/// Descriptor for `GetPrincipalTagAttributeMapInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPrincipalTagAttributeMapInputDescriptor =
    $convert.base64Decode(
        'CiBHZXRQcmluY2lwYWxUYWdBdHRyaWJ1dGVNYXBJbnB1dBIpCg5pZGVudGl0eXBvb2xpZBj9/b'
        'QLIAEoCVIOaWRlbnRpdHlwb29saWQSNgoUaWRlbnRpdHlwcm92aWRlcm5hbWUY8rfwxwEgASgJ'
        'UhRpZGVudGl0eXByb3ZpZGVybmFtZQ==');

@$core.Deprecated('Use getPrincipalTagAttributeMapResponseDescriptor instead')
const GetPrincipalTagAttributeMapResponse$json = {
  '1': 'GetPrincipalTagAttributeMapResponse',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'identityprovidername',
      '3': 419175410,
      '4': 1,
      '5': 9,
      '10': 'identityprovidername'
    },
    {
      '1': 'principaltags',
      '3': 346698229,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_identity.GetPrincipalTagAttributeMapResponse.PrincipaltagsEntry',
      '10': 'principaltags'
    },
    {'1': 'usedefaults', '3': 306413999, '4': 1, '5': 8, '10': 'usedefaults'},
  ],
  '3': [GetPrincipalTagAttributeMapResponse_PrincipaltagsEntry$json],
};

@$core.Deprecated('Use getPrincipalTagAttributeMapResponseDescriptor instead')
const GetPrincipalTagAttributeMapResponse_PrincipaltagsEntry$json = {
  '1': 'PrincipaltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetPrincipalTagAttributeMapResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPrincipalTagAttributeMapResponseDescriptor = $convert.base64Decode(
    'CiNHZXRQcmluY2lwYWxUYWdBdHRyaWJ1dGVNYXBSZXNwb25zZRIpCg5pZGVudGl0eXBvb2xpZB'
    'j9/bQLIAEoCVIOaWRlbnRpdHlwb29saWQSNgoUaWRlbnRpdHlwcm92aWRlcm5hbWUY8rfwxwEg'
    'ASgJUhRpZGVudGl0eXByb3ZpZGVybmFtZRJyCg1wcmluY2lwYWx0YWdzGPXjqKUBIAMoCzJILm'
    'NvZ25pdG9faWRlbnRpdHkuR2V0UHJpbmNpcGFsVGFnQXR0cmlidXRlTWFwUmVzcG9uc2UuUHJp'
    'bmNpcGFsdGFnc0VudHJ5Ug1wcmluY2lwYWx0YWdzEiQKC3VzZWRlZmF1bHRzGK+DjpIBIAEoCF'
    'ILdXNlZGVmYXVsdHMaQAoSUHJpbmNpcGFsdGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQK'
    'BXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use identityDescriptionDescriptor instead')
const IdentityDescription$json = {
  '1': 'IdentityDescription',
  '2': [
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {
      '1': 'lastmodifieddate',
      '3': 24249427,
      '4': 1,
      '5': 9,
      '10': 'lastmodifieddate'
    },
    {'1': 'logins', '3': 109702772, '4': 3, '5': 9, '10': 'logins'},
  ],
};

/// Descriptor for `IdentityDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List identityDescriptionDescriptor = $convert.base64Decode(
    'ChNJZGVudGl0eURlc2NyaXB0aW9uEiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW'
    '9uZGF0ZRIhCgppZGVudGl0eWlkGNfT1W8gASgJUgppZGVudGl0eWlkEi0KEGxhc3Rtb2RpZmll'
    'ZGRhdGUY04jICyABKAlSEGxhc3Rtb2RpZmllZGRhdGUSGQoGbG9naW5zGPTcpzQgAygJUgZsb2'
    'dpbnM=');

@$core.Deprecated('Use identityPoolDescriptor instead')
const IdentityPool$json = {
  '1': 'IdentityPool',
  '2': [
    {
      '1': 'allowclassicflow',
      '3': 101299523,
      '4': 1,
      '5': 8,
      '10': 'allowclassicflow'
    },
    {
      '1': 'allowunauthenticatedidentities',
      '3': 290546287,
      '4': 1,
      '5': 8,
      '10': 'allowunauthenticatedidentities'
    },
    {
      '1': 'cognitoidentityproviders',
      '3': 109610365,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.CognitoIdentityProvider',
      '10': 'cognitoidentityproviders'
    },
    {
      '1': 'developerprovidername',
      '3': 517589792,
      '4': 1,
      '5': 9,
      '10': 'developerprovidername'
    },
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'identitypoolname',
      '3': 41275091,
      '4': 1,
      '5': 9,
      '10': 'identitypoolname'
    },
    {
      '1': 'identitypooltags',
      '3': 305630405,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.IdentityPool.IdentitypooltagsEntry',
      '10': 'identitypooltags'
    },
    {
      '1': 'openidconnectproviderarns',
      '3': 137395556,
      '4': 3,
      '5': 9,
      '10': 'openidconnectproviderarns'
    },
    {
      '1': 'samlproviderarns',
      '3': 24907670,
      '4': 3,
      '5': 9,
      '10': 'samlproviderarns'
    },
    {
      '1': 'supportedloginproviders',
      '3': 125752617,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.IdentityPool.SupportedloginprovidersEntry',
      '10': 'supportedloginproviders'
    },
  ],
  '3': [
    IdentityPool_IdentitypooltagsEntry$json,
    IdentityPool_SupportedloginprovidersEntry$json
  ],
};

@$core.Deprecated('Use identityPoolDescriptor instead')
const IdentityPool_IdentitypooltagsEntry$json = {
  '1': 'IdentitypooltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use identityPoolDescriptor instead')
const IdentityPool_SupportedloginprovidersEntry$json = {
  '1': 'SupportedloginprovidersEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `IdentityPool`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List identityPoolDescriptor = $convert.base64Decode(
    'CgxJZGVudGl0eVBvb2wSLQoQYWxsb3djbGFzc2ljZmxvdxjD6qYwIAEoCFIQYWxsb3djbGFzc2'
    'ljZmxvdxJKCh5hbGxvd3VuYXV0aGVudGljYXRlZGlkZW50aXRpZXMY78TFigEgASgIUh5hbGxv'
    'd3VuYXV0aGVudGljYXRlZGlkZW50aXRpZXMSaAoYY29nbml0b2lkZW50aXR5cHJvdmlkZXJzGP'
    '2KojQgAygLMikuY29nbml0b19pZGVudGl0eS5Db2duaXRvSWRlbnRpdHlQcm92aWRlclIYY29n'
    'bml0b2lkZW50aXR5cHJvdmlkZXJzEjgKFWRldmVsb3BlcnByb3ZpZGVybmFtZRigluf2ASABKA'
    'lSFWRldmVsb3BlcnByb3ZpZGVybmFtZRIpCg5pZGVudGl0eXBvb2xpZBj9/bQLIAEoCVIOaWRl'
    'bnRpdHlwb29saWQSLQoQaWRlbnRpdHlwb29sbmFtZRjTndcTIAEoCVIQaWRlbnRpdHlwb29sbm'
    'FtZRJkChBpZGVudGl0eXBvb2x0YWdzGMWZ3pEBIAMoCzI0LmNvZ25pdG9faWRlbnRpdHkuSWRl'
    'bnRpdHlQb29sLklkZW50aXR5cG9vbHRhZ3NFbnRyeVIQaWRlbnRpdHlwb29sdGFncxI/ChlvcG'
    'VuaWRjb25uZWN0cHJvdmlkZXJhcm5zGOT6wUEgAygJUhlvcGVuaWRjb25uZWN0cHJvdmlkZXJh'
    'cm5zEi0KEHNhbWxwcm92aWRlcmFybnMYlp/wCyADKAlSEHNhbWxwcm92aWRlcmFybnMSeAoXc3'
    'VwcG9ydGVkbG9naW5wcm92aWRlcnMYqar7OyADKAsyOy5jb2duaXRvX2lkZW50aXR5LklkZW50'
    'aXR5UG9vbC5TdXBwb3J0ZWRsb2dpbnByb3ZpZGVyc0VudHJ5UhdzdXBwb3J0ZWRsb2dpbnByb3'
    'ZpZGVycxpDChVJZGVudGl0eXBvb2x0YWdzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFs'
    'dWUYAiABKAlSBXZhbHVlOgI4ARpKChxTdXBwb3J0ZWRsb2dpbnByb3ZpZGVyc0VudHJ5EhAKA2'
    'tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use identityPoolShortDescriptionDescriptor instead')
const IdentityPoolShortDescription$json = {
  '1': 'IdentityPoolShortDescription',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'identitypoolname',
      '3': 41275091,
      '4': 1,
      '5': 9,
      '10': 'identitypoolname'
    },
  ],
};

/// Descriptor for `IdentityPoolShortDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List identityPoolShortDescriptionDescriptor =
    $convert.base64Decode(
        'ChxJZGVudGl0eVBvb2xTaG9ydERlc2NyaXB0aW9uEikKDmlkZW50aXR5cG9vbGlkGP39tAsgAS'
        'gJUg5pZGVudGl0eXBvb2xpZBItChBpZGVudGl0eXBvb2xuYW1lGNOd1xMgASgJUhBpZGVudGl0'
        'eXBvb2xuYW1l');

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

@$core.Deprecated(
    'Use invalidIdentityPoolConfigurationExceptionDescriptor instead')
const InvalidIdentityPoolConfigurationException$json = {
  '1': 'InvalidIdentityPoolConfigurationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidIdentityPoolConfigurationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    invalidIdentityPoolConfigurationExceptionDescriptor = $convert.base64Decode(
        'CilJbnZhbGlkSWRlbnRpdHlQb29sQ29uZmlndXJhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGO'
        'WRyCcgASgJUgdtZXNzYWdl');

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

@$core.Deprecated('Use limitExceededExceptionDescriptor instead')
const LimitExceededException$json = {
  '1': 'LimitExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `LimitExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List limitExceededExceptionDescriptor =
    $convert.base64Decode(
        'ChZMaW1pdEV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use listIdentitiesInputDescriptor instead')
const ListIdentitiesInput$json = {
  '1': 'ListIdentitiesInput',
  '2': [
    {'1': 'hidedisabled', '3': 189789070, '4': 1, '5': 8, '10': 'hidedisabled'},
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListIdentitiesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIdentitiesInputDescriptor = $convert.base64Decode(
    'ChNMaXN0SWRlbnRpdGllc0lucHV0EiUKDGhpZGVkaXNhYmxlZBiO579aIAEoCFIMaGlkZWRpc2'
    'FibGVkEikKDmlkZW50aXR5cG9vbGlkGP39tAsgASgJUg5pZGVudGl0eXBvb2xpZBIiCgptYXhy'
    'ZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leH'
    'R0b2tlbg==');

@$core.Deprecated('Use listIdentitiesResponseDescriptor instead')
const ListIdentitiesResponse$json = {
  '1': 'ListIdentitiesResponse',
  '2': [
    {
      '1': 'identities',
      '3': 452470428,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.IdentityDescription',
      '10': 'identities'
    },
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListIdentitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIdentitiesResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0SWRlbnRpdGllc1Jlc3BvbnNlEkkKCmlkZW50aXRpZXMYnM3g1wEgAygLMiUuY29nbm'
    'l0b19pZGVudGl0eS5JZGVudGl0eURlc2NyaXB0aW9uUgppZGVudGl0aWVzEikKDmlkZW50aXR5'
    'cG9vbGlkGP39tAsgASgJUg5pZGVudGl0eXBvb2xpZBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW'
    '5leHR0b2tlbg==');

@$core.Deprecated('Use listIdentityPoolsInputDescriptor instead')
const ListIdentityPoolsInput$json = {
  '1': 'ListIdentityPoolsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListIdentityPoolsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIdentityPoolsInputDescriptor =
    $convert.base64Decode(
        'ChZMaXN0SWRlbnRpdHlQb29sc0lucHV0EiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZX'
        'N1bHRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listIdentityPoolsResponseDescriptor instead')
const ListIdentityPoolsResponse$json = {
  '1': 'ListIdentityPoolsResponse',
  '2': [
    {
      '1': 'identitypools',
      '3': 158350303,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.IdentityPoolShortDescription',
      '10': 'identitypools'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListIdentityPoolsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listIdentityPoolsResponseDescriptor = $convert.base64Decode(
    'ChlMaXN0SWRlbnRpdHlQb29sc1Jlc3BvbnNlElcKDWlkZW50aXR5cG9vbHMY3/fASyADKAsyLi'
    '5jb2duaXRvX2lkZW50aXR5LklkZW50aXR5UG9vbFNob3J0RGVzY3JpcHRpb25SDWlkZW50aXR5'
    'cG9vbHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listTagsForResourceInputDescriptor instead')
const ListTagsForResourceInput$json = {
  '1': 'ListTagsForResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListTagsForResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceInputDescriptor =
    $convert.base64Decode(
        'ChhMaXN0VGFnc0ZvclJlc291cmNlSW5wdXQSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZX'
        'NvdXJjZWFybg==');

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse$json = {
  '1': 'ListTagsForResourceResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.ListTagsForResourceResponse.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [ListTagsForResourceResponse_TagsEntry$json],
};

@$core.Deprecated('Use listTagsForResourceResponseDescriptor instead')
const ListTagsForResourceResponse_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor = $convert.base64Decode(
    'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USTwoEdGFncxjBwfa1ASADKAsyNy5jb2duaX'
    'RvX2lkZW50aXR5Lkxpc3RUYWdzRm9yUmVzb3VyY2VSZXNwb25zZS5UYWdzRW50cnlSBHRhZ3Ma'
    'NwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOA'
    'E=');

@$core.Deprecated('Use lookupDeveloperIdentityInputDescriptor instead')
const LookupDeveloperIdentityInput$json = {
  '1': 'LookupDeveloperIdentityInput',
  '2': [
    {
      '1': 'developeruseridentifier',
      '3': 349826618,
      '4': 1,
      '5': 9,
      '10': 'developeruseridentifier'
    },
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `LookupDeveloperIdentityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lookupDeveloperIdentityInputDescriptor = $convert.base64Decode(
    'ChxMb29rdXBEZXZlbG9wZXJJZGVudGl0eUlucHV0EjwKF2RldmVsb3BlcnVzZXJpZGVudGlmaW'
    'VyGLrc56YBIAEoCVIXZGV2ZWxvcGVydXNlcmlkZW50aWZpZXISIQoKaWRlbnRpdHlpZBjX09Vv'
    'IAEoCVIKaWRlbnRpdHlpZBIpCg5pZGVudGl0eXBvb2xpZBj9/bQLIAEoCVIOaWRlbnRpdHlwb2'
    '9saWQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6E'
    'umcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use lookupDeveloperIdentityResponseDescriptor instead')
const LookupDeveloperIdentityResponse$json = {
  '1': 'LookupDeveloperIdentityResponse',
  '2': [
    {
      '1': 'developeruseridentifierlist',
      '3': 86156134,
      '4': 3,
      '5': 9,
      '10': 'developeruseridentifierlist'
    },
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `LookupDeveloperIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List lookupDeveloperIdentityResponseDescriptor =
    $convert.base64Decode(
        'Ch9Mb29rdXBEZXZlbG9wZXJJZGVudGl0eVJlc3BvbnNlEkMKG2RldmVsb3BlcnVzZXJpZGVudG'
        'lmaWVybGlzdBjmxoopIAMoCVIbZGV2ZWxvcGVydXNlcmlkZW50aWZpZXJsaXN0EiEKCmlkZW50'
        'aXR5aWQY19PVbyABKAlSCmlkZW50aXR5aWQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
        '9rZW4=');

@$core.Deprecated('Use mappingRuleDescriptor instead')
const MappingRule$json = {
  '1': 'MappingRule',
  '2': [
    {'1': 'claim', '3': 48235024, '4': 1, '5': 9, '10': 'claim'},
    {
      '1': 'matchtype',
      '3': 94592295,
      '4': 1,
      '5': 14,
      '6': '.cognito_identity.MappingRuleMatchType',
      '10': 'matchtype'
    },
    {'1': 'rolearn', '3': 327777153, '4': 1, '5': 9, '10': 'rolearn'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `MappingRule`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mappingRuleDescriptor = $convert.base64Decode(
    'CgtNYXBwaW5nUnVsZRIXCgVjbGFpbRiQhIAXIAEoCVIFY2xhaW0SRwoJbWF0Y2h0eXBlGKe6jS'
    '0gASgOMiYuY29nbml0b19pZGVudGl0eS5NYXBwaW5nUnVsZU1hdGNoVHlwZVIJbWF0Y2h0eXBl'
    'EhwKB3JvbGVhcm4YgfelnAEgASgJUgdyb2xlYXJuEhgKBXZhbHVlGOvyn4oBIAEoCVIFdmFsdW'
    'U=');

@$core.Deprecated('Use mergeDeveloperIdentitiesInputDescriptor instead')
const MergeDeveloperIdentitiesInput$json = {
  '1': 'MergeDeveloperIdentitiesInput',
  '2': [
    {
      '1': 'destinationuseridentifier',
      '3': 333657056,
      '4': 1,
      '5': 9,
      '10': 'destinationuseridentifier'
    },
    {
      '1': 'developerprovidername',
      '3': 517589792,
      '4': 1,
      '5': 9,
      '10': 'developerprovidername'
    },
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'sourceuseridentifier',
      '3': 222471153,
      '4': 1,
      '5': 9,
      '10': 'sourceuseridentifier'
    },
  ],
};

/// Descriptor for `MergeDeveloperIdentitiesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mergeDeveloperIdentitiesInputDescriptor = $convert.base64Decode(
    'Ch1NZXJnZURldmVsb3BlcklkZW50aXRpZXNJbnB1dBJAChlkZXN0aW5hdGlvbnVzZXJpZGVudG'
    'lmaWVyGODnjJ8BIAEoCVIZZGVzdGluYXRpb251c2VyaWRlbnRpZmllchI4ChVkZXZlbG9wZXJw'
    'cm92aWRlcm5hbWUYoJbn9gEgASgJUhVkZXZlbG9wZXJwcm92aWRlcm5hbWUSKQoOaWRlbnRpdH'
    'lwb29saWQY/f20CyABKAlSDmlkZW50aXR5cG9vbGlkEjUKFHNvdXJjZXVzZXJpZGVudGlmaWVy'
    'GPHHimogASgJUhRzb3VyY2V1c2VyaWRlbnRpZmllcg==');

@$core.Deprecated('Use mergeDeveloperIdentitiesResponseDescriptor instead')
const MergeDeveloperIdentitiesResponse$json = {
  '1': 'MergeDeveloperIdentitiesResponse',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
  ],
};

/// Descriptor for `MergeDeveloperIdentitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mergeDeveloperIdentitiesResponseDescriptor =
    $convert.base64Decode(
        'CiBNZXJnZURldmVsb3BlcklkZW50aXRpZXNSZXNwb25zZRIhCgppZGVudGl0eWlkGNfT1W8gAS'
        'gJUgppZGVudGl0eWlk');

@$core.Deprecated('Use notAuthorizedExceptionDescriptor instead')
const NotAuthorizedException$json = {
  '1': 'NotAuthorizedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `NotAuthorizedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List notAuthorizedExceptionDescriptor =
    $convert.base64Decode(
        'ChZOb3RBdXRob3JpemVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use resourceConflictExceptionDescriptor instead')
const ResourceConflictException$json = {
  '1': 'ResourceConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceConflictExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZUNvbmZsaWN0RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

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

@$core.Deprecated('Use roleMappingDescriptor instead')
const RoleMapping$json = {
  '1': 'RoleMapping',
  '2': [
    {
      '1': 'ambiguousroleresolution',
      '3': 410304778,
      '4': 1,
      '5': 14,
      '6': '.cognito_identity.AmbiguousRoleResolutionType',
      '10': 'ambiguousroleresolution'
    },
    {
      '1': 'rulesconfiguration',
      '3': 123571251,
      '4': 1,
      '5': 11,
      '6': '.cognito_identity.RulesConfigurationType',
      '10': 'rulesconfiguration'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.cognito_identity.RoleMappingType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `RoleMapping`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List roleMappingDescriptor = $convert.base64Decode(
    'CgtSb2xlTWFwcGluZxJrChdhbWJpZ3VvdXNyb2xlcmVzb2x1dGlvbhiKgtPDASABKA4yLS5jb2'
    'duaXRvX2lkZW50aXR5LkFtYmlndW91c1JvbGVSZXNvbHV0aW9uVHlwZVIXYW1iaWd1b3Vzcm9s'
    'ZXJlc29sdXRpb24SWwoScnVsZXNjb25maWd1cmF0aW9uGLOY9jogASgLMiguY29nbml0b19pZG'
    'VudGl0eS5SdWxlc0NvbmZpZ3VyYXRpb25UeXBlUhJydWxlc2NvbmZpZ3VyYXRpb24SOQoEdHlw'
    'ZRjuoNeKASABKA4yIS5jb2duaXRvX2lkZW50aXR5LlJvbGVNYXBwaW5nVHlwZVIEdHlwZQ==');

@$core.Deprecated('Use rulesConfigurationTypeDescriptor instead')
const RulesConfigurationType$json = {
  '1': 'RulesConfigurationType',
  '2': [
    {
      '1': 'rules',
      '3': 42675585,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.MappingRule',
      '10': 'rules'
    },
  ],
};

/// Descriptor for `RulesConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rulesConfigurationTypeDescriptor =
    $convert.base64Decode(
        'ChZSdWxlc0NvbmZpZ3VyYXRpb25UeXBlEjYKBXJ1bGVzGIHbrBQgAygLMh0uY29nbml0b19pZG'
        'VudGl0eS5NYXBwaW5nUnVsZVIFcnVsZXM=');

@$core.Deprecated('Use setIdentityPoolRolesInputDescriptor instead')
const SetIdentityPoolRolesInput$json = {
  '1': 'SetIdentityPoolRolesInput',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'rolemappings',
      '3': 96154047,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.SetIdentityPoolRolesInput.RolemappingsEntry',
      '10': 'rolemappings'
    },
    {
      '1': 'roles',
      '3': 511168127,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.SetIdentityPoolRolesInput.RolesEntry',
      '10': 'roles'
    },
  ],
  '3': [
    SetIdentityPoolRolesInput_RolemappingsEntry$json,
    SetIdentityPoolRolesInput_RolesEntry$json
  ],
};

@$core.Deprecated('Use setIdentityPoolRolesInputDescriptor instead')
const SetIdentityPoolRolesInput_RolemappingsEntry$json = {
  '1': 'RolemappingsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.cognito_identity.RoleMapping',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use setIdentityPoolRolesInputDescriptor instead')
const SetIdentityPoolRolesInput_RolesEntry$json = {
  '1': 'RolesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SetIdentityPoolRolesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setIdentityPoolRolesInputDescriptor = $convert.base64Decode(
    'ChlTZXRJZGVudGl0eVBvb2xSb2xlc0lucHV0EikKDmlkZW50aXR5cG9vbGlkGP39tAsgASgJUg'
    '5pZGVudGl0eXBvb2xpZBJkCgxyb2xlbWFwcGluZ3MYv+PsLSADKAsyPS5jb2duaXRvX2lkZW50'
    'aXR5LlNldElkZW50aXR5UG9vbFJvbGVzSW5wdXQuUm9sZW1hcHBpbmdzRW50cnlSDHJvbGVtYX'
    'BwaW5ncxJQCgVyb2xlcxj/nN/zASADKAsyNi5jb2duaXRvX2lkZW50aXR5LlNldElkZW50aXR5'
    'UG9vbFJvbGVzSW5wdXQuUm9sZXNFbnRyeVIFcm9sZXMaXgoRUm9sZW1hcHBpbmdzRW50cnkSEA'
    'oDa2V5GAEgASgJUgNrZXkSMwoFdmFsdWUYAiABKAsyHS5jb2duaXRvX2lkZW50aXR5LlJvbGVN'
    'YXBwaW5nUgV2YWx1ZToCOAEaOAoKUm9sZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YW'
    'x1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use setPrincipalTagAttributeMapInputDescriptor instead')
const SetPrincipalTagAttributeMapInput$json = {
  '1': 'SetPrincipalTagAttributeMapInput',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'identityprovidername',
      '3': 419175410,
      '4': 1,
      '5': 9,
      '10': 'identityprovidername'
    },
    {
      '1': 'principaltags',
      '3': 346698229,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_identity.SetPrincipalTagAttributeMapInput.PrincipaltagsEntry',
      '10': 'principaltags'
    },
    {'1': 'usedefaults', '3': 306413999, '4': 1, '5': 8, '10': 'usedefaults'},
  ],
  '3': [SetPrincipalTagAttributeMapInput_PrincipaltagsEntry$json],
};

@$core.Deprecated('Use setPrincipalTagAttributeMapInputDescriptor instead')
const SetPrincipalTagAttributeMapInput_PrincipaltagsEntry$json = {
  '1': 'PrincipaltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SetPrincipalTagAttributeMapInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setPrincipalTagAttributeMapInputDescriptor = $convert.base64Decode(
    'CiBTZXRQcmluY2lwYWxUYWdBdHRyaWJ1dGVNYXBJbnB1dBIpCg5pZGVudGl0eXBvb2xpZBj9/b'
    'QLIAEoCVIOaWRlbnRpdHlwb29saWQSNgoUaWRlbnRpdHlwcm92aWRlcm5hbWUY8rfwxwEgASgJ'
    'UhRpZGVudGl0eXByb3ZpZGVybmFtZRJvCg1wcmluY2lwYWx0YWdzGPXjqKUBIAMoCzJFLmNvZ2'
    '5pdG9faWRlbnRpdHkuU2V0UHJpbmNpcGFsVGFnQXR0cmlidXRlTWFwSW5wdXQuUHJpbmNpcGFs'
    'dGFnc0VudHJ5Ug1wcmluY2lwYWx0YWdzEiQKC3VzZWRlZmF1bHRzGK+DjpIBIAEoCFILdXNlZG'
    'VmYXVsdHMaQAoSUHJpbmNpcGFsdGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVl'
    'GAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use setPrincipalTagAttributeMapResponseDescriptor instead')
const SetPrincipalTagAttributeMapResponse$json = {
  '1': 'SetPrincipalTagAttributeMapResponse',
  '2': [
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
    {
      '1': 'identityprovidername',
      '3': 419175410,
      '4': 1,
      '5': 9,
      '10': 'identityprovidername'
    },
    {
      '1': 'principaltags',
      '3': 346698229,
      '4': 3,
      '5': 11,
      '6':
          '.cognito_identity.SetPrincipalTagAttributeMapResponse.PrincipaltagsEntry',
      '10': 'principaltags'
    },
    {'1': 'usedefaults', '3': 306413999, '4': 1, '5': 8, '10': 'usedefaults'},
  ],
  '3': [SetPrincipalTagAttributeMapResponse_PrincipaltagsEntry$json],
};

@$core.Deprecated('Use setPrincipalTagAttributeMapResponseDescriptor instead')
const SetPrincipalTagAttributeMapResponse_PrincipaltagsEntry$json = {
  '1': 'PrincipaltagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SetPrincipalTagAttributeMapResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setPrincipalTagAttributeMapResponseDescriptor = $convert.base64Decode(
    'CiNTZXRQcmluY2lwYWxUYWdBdHRyaWJ1dGVNYXBSZXNwb25zZRIpCg5pZGVudGl0eXBvb2xpZB'
    'j9/bQLIAEoCVIOaWRlbnRpdHlwb29saWQSNgoUaWRlbnRpdHlwcm92aWRlcm5hbWUY8rfwxwEg'
    'ASgJUhRpZGVudGl0eXByb3ZpZGVybmFtZRJyCg1wcmluY2lwYWx0YWdzGPXjqKUBIAMoCzJILm'
    'NvZ25pdG9faWRlbnRpdHkuU2V0UHJpbmNpcGFsVGFnQXR0cmlidXRlTWFwUmVzcG9uc2UuUHJp'
    'bmNpcGFsdGFnc0VudHJ5Ug1wcmluY2lwYWx0YWdzEiQKC3VzZWRlZmF1bHRzGK+DjpIBIAEoCF'
    'ILdXNlZGVmYXVsdHMaQAoSUHJpbmNpcGFsdGFnc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQK'
    'BXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use tagResourceInputDescriptor instead')
const TagResourceInput$json = {
  '1': 'TagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.TagResourceInput.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [TagResourceInput_TagsEntry$json],
};

@$core.Deprecated('Use tagResourceInputDescriptor instead')
const TagResourceInput_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceInputDescriptor = $convert.base64Decode(
    'ChBUYWdSZXNvdXJjZUlucHV0EiQKC3Jlc291cmNlYXJuGK342a0BIAEoCVILcmVzb3VyY2Vhcm'
    '4SRAoEdGFncxjBwfa1ASADKAsyLC5jb2duaXRvX2lkZW50aXR5LlRhZ1Jlc291cmNlSW5wdXQu'
    'VGFnc0VudHJ5UgR0YWdzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZR'
    'gCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use tagResourceResponseDescriptor instead')
const TagResourceResponse$json = {
  '1': 'TagResourceResponse',
};

/// Descriptor for `TagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceResponseDescriptor =
    $convert.base64Decode('ChNUYWdSZXNvdXJjZVJlc3BvbnNl');

@$core.Deprecated('Use tooManyRequestsExceptionDescriptor instead')
const TooManyRequestsException$json = {
  '1': 'TooManyRequestsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyRequestsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyRequestsExceptionDescriptor =
    $convert.base64Decode(
        'ChhUb29NYW55UmVxdWVzdHNFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use unlinkDeveloperIdentityInputDescriptor instead')
const UnlinkDeveloperIdentityInput$json = {
  '1': 'UnlinkDeveloperIdentityInput',
  '2': [
    {
      '1': 'developerprovidername',
      '3': 517589792,
      '4': 1,
      '5': 9,
      '10': 'developerprovidername'
    },
    {
      '1': 'developeruseridentifier',
      '3': 349826618,
      '4': 1,
      '5': 9,
      '10': 'developeruseridentifier'
    },
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {
      '1': 'identitypoolid',
      '3': 23936765,
      '4': 1,
      '5': 9,
      '10': 'identitypoolid'
    },
  ],
};

/// Descriptor for `UnlinkDeveloperIdentityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unlinkDeveloperIdentityInputDescriptor = $convert.base64Decode(
    'ChxVbmxpbmtEZXZlbG9wZXJJZGVudGl0eUlucHV0EjgKFWRldmVsb3BlcnByb3ZpZGVybmFtZR'
    'igluf2ASABKAlSFWRldmVsb3BlcnByb3ZpZGVybmFtZRI8ChdkZXZlbG9wZXJ1c2VyaWRlbnRp'
    'Zmllchi63OemASABKAlSF2RldmVsb3BlcnVzZXJpZGVudGlmaWVyEiEKCmlkZW50aXR5aWQY19'
    'PVbyABKAlSCmlkZW50aXR5aWQSKQoOaWRlbnRpdHlwb29saWQY/f20CyABKAlSDmlkZW50aXR5'
    'cG9vbGlk');

@$core.Deprecated('Use unlinkIdentityInputDescriptor instead')
const UnlinkIdentityInput$json = {
  '1': 'UnlinkIdentityInput',
  '2': [
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
    {
      '1': 'logins',
      '3': 109702772,
      '4': 3,
      '5': 11,
      '6': '.cognito_identity.UnlinkIdentityInput.LoginsEntry',
      '10': 'logins'
    },
    {
      '1': 'loginstoremove',
      '3': 490649531,
      '4': 3,
      '5': 9,
      '10': 'loginstoremove'
    },
  ],
  '3': [UnlinkIdentityInput_LoginsEntry$json],
};

@$core.Deprecated('Use unlinkIdentityInputDescriptor instead')
const UnlinkIdentityInput_LoginsEntry$json = {
  '1': 'LoginsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `UnlinkIdentityInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unlinkIdentityInputDescriptor = $convert.base64Decode(
    'ChNVbmxpbmtJZGVudGl0eUlucHV0EiEKCmlkZW50aXR5aWQY19PVbyABKAlSCmlkZW50aXR5aW'
    'QSTAoGbG9naW5zGPTcpzQgAygLMjEuY29nbml0b19pZGVudGl0eS5VbmxpbmtJZGVudGl0eUlu'
    'cHV0LkxvZ2luc0VudHJ5UgZsb2dpbnMSKgoObG9naW5zdG9yZW1vdmUYu+/66QEgAygJUg5sb2'
    'dpbnN0b3JlbW92ZRo5CgtMb2dpbnNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgC'
    'IAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use unprocessedIdentityIdDescriptor instead')
const UnprocessedIdentityId$json = {
  '1': 'UnprocessedIdentityId',
  '2': [
    {
      '1': 'errorcode',
      '3': 34663193,
      '4': 1,
      '5': 14,
      '6': '.cognito_identity.ErrorCode',
      '10': 'errorcode'
    },
    {'1': 'identityid', '3': 234187223, '4': 1, '5': 9, '10': 'identityid'},
  ],
};

/// Descriptor for `UnprocessedIdentityId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unprocessedIdentityIdDescriptor = $convert.base64Decode(
    'ChVVbnByb2Nlc3NlZElkZW50aXR5SWQSPAoJZXJyb3Jjb2RlGJnWwxAgASgOMhsuY29nbml0b1'
    '9pZGVudGl0eS5FcnJvckNvZGVSCWVycm9yY29kZRIhCgppZGVudGl0eWlkGNfT1W8gASgJUgpp'
    'ZGVudGl0eWlk');

@$core.Deprecated('Use untagResourceInputDescriptor instead')
const UntagResourceInput$json = {
  '1': 'UntagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceInputDescriptor = $convert.base64Decode(
    'ChJVbnRhZ1Jlc291cmNlSW5wdXQSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJjZW'
    'FybhIcCgd0YWdrZXlzGPzD85gBIAMoCVIHdGFna2V5cw==');

@$core.Deprecated('Use untagResourceResponseDescriptor instead')
const UntagResourceResponse$json = {
  '1': 'UntagResourceResponse',
};

/// Descriptor for `UntagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceResponseDescriptor =
    $convert.base64Decode('ChVVbnRhZ1Jlc291cmNlUmVzcG9uc2U=');
