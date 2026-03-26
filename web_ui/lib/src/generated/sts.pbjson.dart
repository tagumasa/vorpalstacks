// This is a generated file - do not edit.
//
// Generated from sts.proto.

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

@$core.Deprecated('Use assumeRoleRequestDescriptor instead')
const AssumeRoleRequest$json = {
  '1': 'AssumeRoleRequest',
  '2': [
    {
      '1': 'durationseconds',
      '3': 451873635,
      '4': 1,
      '5': 5,
      '10': 'durationseconds'
    },
    {'1': 'externalid', '3': 271401992, '4': 1, '5': 9, '10': 'externalid'},
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'policyarns',
      '3': 183785508,
      '4': 3,
      '5': 11,
      '6': '.sts.PolicyDescriptorType',
      '10': 'policyarns'
    },
    {
      '1': 'providedcontexts',
      '3': 228510151,
      '4': 3,
      '5': 11,
      '6': '.sts.ProvidedContext',
      '10': 'providedcontexts'
    },
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'rolesessionname',
      '3': 315098849,
      '4': 1,
      '5': 9,
      '10': 'rolesessionname'
    },
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {
      '1': 'sourceidentity',
      '3': 466635355,
      '4': 1,
      '5': 9,
      '10': 'sourceidentity'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sts.Tag',
      '10': 'tags'
    },
    {'1': 'tokencode', '3': 300671456, '4': 1, '5': 9, '10': 'tokencode'},
    {
      '1': 'transitivetagkeys',
      '3': 452608727,
      '4': 3,
      '5': 9,
      '10': 'transitivetagkeys'
    },
  ],
};

/// Descriptor for `AssumeRoleRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRoleRequestDescriptor = $convert.base64Decode(
    'ChFBc3N1bWVSb2xlUmVxdWVzdBIsCg9kdXJhdGlvbnNlY29uZHMY45a81wEgASgFUg9kdXJhdG'
    'lvbnNlY29uZHMSIgoKZXh0ZXJuYWxpZBiIiLWBASABKAlSCmV4dGVybmFsaWQSGgoGcG9saWN5'
    'GKDv8OABIAEoCVIGcG9saWN5EjwKCnBvbGljeWFybnMYpLDRVyADKAsyGS5zdHMuUG9saWN5RG'
    'VzY3JpcHRvclR5cGVSCnBvbGljeWFybnMSQwoQcHJvdmlkZWRjb250ZXh0cxjHk/tsIAMoCzIU'
    'LnN0cy5Qcm92aWRlZENvbnRleHRSEHByb3ZpZGVkY29udGV4dHMSHAoHcm9sZWFybhiB+OeZAS'
    'ABKAlSB3JvbGVhcm4SLAoPcm9sZXNlc3Npb25uYW1lGOGNoJYBIAEoCVIPcm9sZXNlc3Npb25u'
    'YW1lEiYKDHNlcmlhbG51bWJlchjlurnHASABKAlSDHNlcmlhbG51bWJlchIqCg5zb3VyY2VpZG'
    'VudGl0eRjblMHeASABKAlSDnNvdXJjZWlkZW50aXR5EiAKBHRhZ3MYwcH2tQEgAygLMgguc3Rz'
    'LlRhZ1IEdGFncxIgCgl0b2tlbmNvZGUY4MOvjwEgASgJUgl0b2tlbmNvZGUSMAoRdHJhbnNpdG'
    'l2ZXRhZ2tleXMY14Xp1wEgAygJUhF0cmFuc2l0aXZldGFna2V5cw==');

@$core.Deprecated('Use assumeRoleResponseDescriptor instead')
const AssumeRoleResponse$json = {
  '1': 'AssumeRoleResponse',
  '2': [
    {
      '1': 'assumedroleuser',
      '3': 314673579,
      '4': 1,
      '5': 11,
      '6': '.sts.AssumedRoleUser',
      '10': 'assumedroleuser'
    },
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.sts.Credentials',
      '10': 'credentials'
    },
    {
      '1': 'packedpolicysize',
      '3': 511234267,
      '4': 1,
      '5': 5,
      '10': 'packedpolicysize'
    },
    {
      '1': 'sourceidentity',
      '3': 466635355,
      '4': 1,
      '5': 9,
      '10': 'sourceidentity'
    },
  ],
};

/// Descriptor for `AssumeRoleResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRoleResponseDescriptor = $convert.base64Decode(
    'ChJBc3N1bWVSb2xlUmVzcG9uc2USQgoPYXNzdW1lZHJvbGV1c2VyGKuThpYBIAEoCzIULnN0cy'
    '5Bc3N1bWVkUm9sZVVzZXJSD2Fzc3VtZWRyb2xldXNlchI2CgtjcmVkZW50aWFscxjymo62ASAB'
    'KAsyEC5zdHMuQ3JlZGVudGlhbHNSC2NyZWRlbnRpYWxzEi4KEHBhY2tlZHBvbGljeXNpemUY26'
    'Hj8wEgASgFUhBwYWNrZWRwb2xpY3lzaXplEioKDnNvdXJjZWlkZW50aXR5GNuUwd4BIAEoCVIO'
    'c291cmNlaWRlbnRpdHk=');

@$core.Deprecated('Use assumeRoleWithSAMLRequestDescriptor instead')
const AssumeRoleWithSAMLRequest$json = {
  '1': 'AssumeRoleWithSAMLRequest',
  '2': [
    {
      '1': 'durationseconds',
      '3': 451873635,
      '4': 1,
      '5': 5,
      '10': 'durationseconds'
    },
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'policyarns',
      '3': 183785508,
      '4': 3,
      '5': 11,
      '6': '.sts.PolicyDescriptorType',
      '10': 'policyarns'
    },
    {'1': 'principalarn', '3': 93469969, '4': 1, '5': 9, '10': 'principalarn'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'samlassertion',
      '3': 202921933,
      '4': 1,
      '5': 9,
      '10': 'samlassertion'
    },
  ],
};

/// Descriptor for `AssumeRoleWithSAMLRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRoleWithSAMLRequestDescriptor = $convert.base64Decode(
    'ChlBc3N1bWVSb2xlV2l0aFNBTUxSZXF1ZXN0EiwKD2R1cmF0aW9uc2Vjb25kcxjjlrzXASABKA'
    'VSD2R1cmF0aW9uc2Vjb25kcxIaCgZwb2xpY3kYoO/w4AEgASgJUgZwb2xpY3kSPAoKcG9saWN5'
    'YXJucxiksNFXIAMoCzIZLnN0cy5Qb2xpY3lEZXNjcmlwdG9yVHlwZVIKcG9saWN5YXJucxIlCg'
    'xwcmluY2lwYWxhcm4YkfrILCABKAlSDHByaW5jaXBhbGFybhIcCgdyb2xlYXJuGIH455kBIAEo'
    'CVIHcm9sZWFybhInCg1zYW1sYXNzZXJ0aW9uGM2v4WAgASgJUg1zYW1sYXNzZXJ0aW9u');

@$core.Deprecated('Use assumeRoleWithSAMLResponseDescriptor instead')
const AssumeRoleWithSAMLResponse$json = {
  '1': 'AssumeRoleWithSAMLResponse',
  '2': [
    {
      '1': 'assumedroleuser',
      '3': 314673579,
      '4': 1,
      '5': 11,
      '6': '.sts.AssumedRoleUser',
      '10': 'assumedroleuser'
    },
    {'1': 'audience', '3': 284892548, '4': 1, '5': 9, '10': 'audience'},
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.sts.Credentials',
      '10': 'credentials'
    },
    {'1': 'issuer', '3': 528708823, '4': 1, '5': 9, '10': 'issuer'},
    {
      '1': 'namequalifier',
      '3': 521907559,
      '4': 1,
      '5': 9,
      '10': 'namequalifier'
    },
    {
      '1': 'packedpolicysize',
      '3': 511234267,
      '4': 1,
      '5': 5,
      '10': 'packedpolicysize'
    },
    {
      '1': 'sourceidentity',
      '3': 466635355,
      '4': 1,
      '5': 9,
      '10': 'sourceidentity'
    },
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
    {'1': 'subjecttype', '3': 222881976, '4': 1, '5': 9, '10': 'subjecttype'},
  ],
};

/// Descriptor for `AssumeRoleWithSAMLResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRoleWithSAMLResponseDescriptor = $convert.base64Decode(
    'ChpBc3N1bWVSb2xlV2l0aFNBTUxSZXNwb25zZRJCCg9hc3N1bWVkcm9sZXVzZXIYq5OGlgEgAS'
    'gLMhQuc3RzLkFzc3VtZWRSb2xlVXNlclIPYXNzdW1lZHJvbGV1c2VyEh4KCGF1ZGllbmNlGIS7'
    '7IcBIAEoCVIIYXVkaWVuY2USNgoLY3JlZGVudGlhbHMY8pqOtgEgASgLMhAuc3RzLkNyZWRlbn'
    'RpYWxzUgtjcmVkZW50aWFscxIaCgZpc3N1ZXIY1+mN/AEgASgJUgZpc3N1ZXISKAoNbmFtZXF1'
    'YWxpZmllchjn2u74ASABKAlSDW5hbWVxdWFsaWZpZXISLgoQcGFja2VkcG9saWN5c2l6ZRjboe'
    'PzASABKAVSEHBhY2tlZHBvbGljeXNpemUSKgoOc291cmNlaWRlbnRpdHkY25TB3gEgASgJUg5z'
    'b3VyY2VpZGVudGl0eRIbCgdzdWJqZWN0GPDJ5AMgASgJUgdzdWJqZWN0EiMKC3N1YmplY3R0eX'
    'BlGLjRo2ogASgJUgtzdWJqZWN0dHlwZQ==');

@$core.Deprecated('Use assumeRoleWithWebIdentityRequestDescriptor instead')
const AssumeRoleWithWebIdentityRequest$json = {
  '1': 'AssumeRoleWithWebIdentityRequest',
  '2': [
    {
      '1': 'durationseconds',
      '3': 451873635,
      '4': 1,
      '5': 5,
      '10': 'durationseconds'
    },
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'policyarns',
      '3': 183785508,
      '4': 3,
      '5': 11,
      '6': '.sts.PolicyDescriptorType',
      '10': 'policyarns'
    },
    {'1': 'providerid', '3': 509712370, '4': 1, '5': 9, '10': 'providerid'},
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'rolesessionname',
      '3': 315098849,
      '4': 1,
      '5': 9,
      '10': 'rolesessionname'
    },
    {
      '1': 'webidentitytoken',
      '3': 234014869,
      '4': 1,
      '5': 9,
      '10': 'webidentitytoken'
    },
  ],
};

/// Descriptor for `AssumeRoleWithWebIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRoleWithWebIdentityRequestDescriptor = $convert.base64Decode(
    'CiBBc3N1bWVSb2xlV2l0aFdlYklkZW50aXR5UmVxdWVzdBIsCg9kdXJhdGlvbnNlY29uZHMY45'
    'a81wEgASgFUg9kdXJhdGlvbnNlY29uZHMSGgoGcG9saWN5GKDv8OABIAEoCVIGcG9saWN5EjwK'
    'CnBvbGljeWFybnMYpLDRVyADKAsyGS5zdHMuUG9saWN5RGVzY3JpcHRvclR5cGVSCnBvbGljeW'
    'FybnMSIgoKcHJvdmlkZXJpZBjyr4bzASABKAlSCnByb3ZpZGVyaWQSHAoHcm9sZWFybhiB+OeZ'
    'ASABKAlSB3JvbGVhcm4SLAoPcm9sZXNlc3Npb25uYW1lGOGNoJYBIAEoCVIPcm9sZXNlc3Npb2'
    '5uYW1lEi0KEHdlYmlkZW50aXR5dG9rZW4YlZHLbyABKAlSEHdlYmlkZW50aXR5dG9rZW4=');

@$core.Deprecated('Use assumeRoleWithWebIdentityResponseDescriptor instead')
const AssumeRoleWithWebIdentityResponse$json = {
  '1': 'AssumeRoleWithWebIdentityResponse',
  '2': [
    {
      '1': 'assumedroleuser',
      '3': 314673579,
      '4': 1,
      '5': 11,
      '6': '.sts.AssumedRoleUser',
      '10': 'assumedroleuser'
    },
    {'1': 'audience', '3': 284892548, '4': 1, '5': 9, '10': 'audience'},
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.sts.Credentials',
      '10': 'credentials'
    },
    {
      '1': 'packedpolicysize',
      '3': 511234267,
      '4': 1,
      '5': 5,
      '10': 'packedpolicysize'
    },
    {'1': 'provider', '3': 363366621, '4': 1, '5': 9, '10': 'provider'},
    {
      '1': 'sourceidentity',
      '3': 466635355,
      '4': 1,
      '5': 9,
      '10': 'sourceidentity'
    },
    {
      '1': 'subjectfromwebidentitytoken',
      '3': 96354739,
      '4': 1,
      '5': 9,
      '10': 'subjectfromwebidentitytoken'
    },
  ],
};

/// Descriptor for `AssumeRoleWithWebIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRoleWithWebIdentityResponseDescriptor = $convert.base64Decode(
    'CiFBc3N1bWVSb2xlV2l0aFdlYklkZW50aXR5UmVzcG9uc2USQgoPYXNzdW1lZHJvbGV1c2VyGK'
    'uThpYBIAEoCzIULnN0cy5Bc3N1bWVkUm9sZVVzZXJSD2Fzc3VtZWRyb2xldXNlchIeCghhdWRp'
    'ZW5jZRiEu+yHASABKAlSCGF1ZGllbmNlEjYKC2NyZWRlbnRpYWxzGPKajrYBIAEoCzIQLnN0cy'
    '5DcmVkZW50aWFsc1ILY3JlZGVudGlhbHMSLgoQcGFja2VkcG9saWN5c2l6ZRjboePzASABKAVS'
    'EHBhY2tlZHBvbGljeXNpemUSHgoIcHJvdmlkZXIY3ZGirQEgASgJUghwcm92aWRlchIqCg5zb3'
    'VyY2VpZGVudGl0eRjblMHeASABKAlSDnNvdXJjZWlkZW50aXR5EkMKG3N1YmplY3Rmcm9td2Vi'
    'aWRlbnRpdHl0b2tlbhizg/ktIAEoCVIbc3ViamVjdGZyb213ZWJpZGVudGl0eXRva2Vu');

@$core.Deprecated('Use assumeRootRequestDescriptor instead')
const AssumeRootRequest$json = {
  '1': 'AssumeRootRequest',
  '2': [
    {
      '1': 'durationseconds',
      '3': 451873635,
      '4': 1,
      '5': 5,
      '10': 'durationseconds'
    },
    {
      '1': 'targetprincipal',
      '3': 215916667,
      '4': 1,
      '5': 9,
      '10': 'targetprincipal'
    },
    {
      '1': 'taskpolicyarn',
      '3': 332406370,
      '4': 1,
      '5': 11,
      '6': '.sts.PolicyDescriptorType',
      '10': 'taskpolicyarn'
    },
  ],
};

/// Descriptor for `AssumeRootRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRootRequestDescriptor = $convert.base64Decode(
    'ChFBc3N1bWVSb290UmVxdWVzdBIsCg9kdXJhdGlvbnNlY29uZHMY45a81wEgASgFUg9kdXJhdG'
    'lvbnNlY29uZHMSKwoPdGFyZ2V0cHJpbmNpcGFsGPvA+mYgASgJUg90YXJnZXRwcmluY2lwYWwS'
    'QwoNdGFza3BvbGljeWFybhjivMCeASABKAsyGS5zdHMuUG9saWN5RGVzY3JpcHRvclR5cGVSDX'
    'Rhc2twb2xpY3lhcm4=');

@$core.Deprecated('Use assumeRootResponseDescriptor instead')
const AssumeRootResponse$json = {
  '1': 'AssumeRootResponse',
  '2': [
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.sts.Credentials',
      '10': 'credentials'
    },
    {
      '1': 'sourceidentity',
      '3': 466635355,
      '4': 1,
      '5': 9,
      '10': 'sourceidentity'
    },
  ],
};

/// Descriptor for `AssumeRootResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumeRootResponseDescriptor = $convert.base64Decode(
    'ChJBc3N1bWVSb290UmVzcG9uc2USNgoLY3JlZGVudGlhbHMY8pqOtgEgASgLMhAuc3RzLkNyZW'
    'RlbnRpYWxzUgtjcmVkZW50aWFscxIqCg5zb3VyY2VpZGVudGl0eRjblMHeASABKAlSDnNvdXJj'
    'ZWlkZW50aXR5');

@$core.Deprecated('Use assumedRoleUserDescriptor instead')
const AssumedRoleUser$json = {
  '1': 'AssumedRoleUser',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'assumedroleid',
      '3': 244783337,
      '4': 1,
      '5': 9,
      '10': 'assumedroleid'
    },
  ],
};

/// Descriptor for `AssumedRoleUser`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List assumedRoleUserDescriptor = $convert.base64Decode(
    'Cg9Bc3N1bWVkUm9sZVVzZXISFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEicKDWFzc3VtZWRyb2xlaW'
    'QY6bHcdCABKAlSDWFzc3VtZWRyb2xlaWQ=');

@$core.Deprecated('Use credentialsDescriptor instead')
const Credentials$json = {
  '1': 'Credentials',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {
      '1': 'secretaccesskey',
      '3': 172288487,
      '4': 1,
      '5': 9,
      '10': 'secretaccesskey'
    },
    {'1': 'sessiontoken', '3': 211161069, '4': 1, '5': 9, '10': 'sessiontoken'},
  ],
};

/// Descriptor for `Credentials`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List credentialsDescriptor = $convert.base64Decode(
    'CgtDcmVkZW50aWFscxIkCgthY2Nlc3NrZXlpZBigt7fYASABKAlSC2FjY2Vzc2tleWlkEiEKCm'
    'V4cGlyYXRpb24YiamfdSABKAlSCmV4cGlyYXRpb24SKwoPc2VjcmV0YWNjZXNza2V5GOfTk1Ig'
    'ASgJUg9zZWNyZXRhY2Nlc3NrZXkSJQoMc2Vzc2lvbnRva2VuGO2f2GQgASgJUgxzZXNzaW9udG'
    '9rZW4=');

@$core.Deprecated('Use decodeAuthorizationMessageRequestDescriptor instead')
const DecodeAuthorizationMessageRequest$json = {
  '1': 'DecodeAuthorizationMessageRequest',
  '2': [
    {
      '1': 'encodedmessage',
      '3': 176632049,
      '4': 1,
      '5': 9,
      '10': 'encodedmessage'
    },
  ],
};

/// Descriptor for `DecodeAuthorizationMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List decodeAuthorizationMessageRequestDescriptor =
    $convert.base64Decode(
        'CiFEZWNvZGVBdXRob3JpemF0aW9uTWVzc2FnZVJlcXVlc3QSKQoOZW5jb2RlZG1lc3NhZ2UY8e'
        'GcVCABKAlSDmVuY29kZWRtZXNzYWdl');

@$core.Deprecated('Use decodeAuthorizationMessageResponseDescriptor instead')
const DecodeAuthorizationMessageResponse$json = {
  '1': 'DecodeAuthorizationMessageResponse',
  '2': [
    {
      '1': 'decodedmessage',
      '3': 475373641,
      '4': 1,
      '5': 9,
      '10': 'decodedmessage'
    },
  ],
};

/// Descriptor for `DecodeAuthorizationMessageResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List decodeAuthorizationMessageResponseDescriptor =
    $convert.base64Decode(
        'CiJEZWNvZGVBdXRob3JpemF0aW9uTWVzc2FnZVJlc3BvbnNlEioKDmRlY29kZWRtZXNzYWdlGM'
        'nA1uIBIAEoCVIOZGVjb2RlZG1lc3NhZ2U=');

@$core.Deprecated('Use expiredTokenExceptionDescriptor instead')
const ExpiredTokenException$json = {
  '1': 'ExpiredTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExpiredTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expiredTokenExceptionDescriptor = $convert.base64Decode(
    'ChVFeHBpcmVkVG9rZW5FeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use expiredTradeInTokenExceptionDescriptor instead')
const ExpiredTradeInTokenException$json = {
  '1': 'ExpiredTradeInTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExpiredTradeInTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expiredTradeInTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChxFeHBpcmVkVHJhZGVJblRva2VuRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use federatedUserDescriptor instead')
const FederatedUser$json = {
  '1': 'FederatedUser',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'federateduserid',
      '3': 366276542,
      '4': 1,
      '5': 9,
      '10': 'federateduserid'
    },
  ],
};

/// Descriptor for `FederatedUser`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List federatedUserDescriptor = $convert.base64Decode(
    'Cg1GZWRlcmF0ZWRVc2VyEhQKA2Fybhidm+2/ASABKAlSA2FybhIsCg9mZWRlcmF0ZWR1c2VyaW'
    'QYvt/TrgEgASgJUg9mZWRlcmF0ZWR1c2VyaWQ=');

@$core.Deprecated('Use getAccessKeyInfoRequestDescriptor instead')
const GetAccessKeyInfoRequest$json = {
  '1': 'GetAccessKeyInfoRequest',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
  ],
};

/// Descriptor for `GetAccessKeyInfoRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccessKeyInfoRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRBY2Nlc3NLZXlJbmZvUmVxdWVzdBIkCgthY2Nlc3NrZXlpZBigt7fYASABKAlSC2FjY2'
        'Vzc2tleWlk');

@$core.Deprecated('Use getAccessKeyInfoResponseDescriptor instead')
const GetAccessKeyInfoResponse$json = {
  '1': 'GetAccessKeyInfoResponse',
  '2': [
    {'1': 'account', '3': 435725053, '4': 1, '5': 9, '10': 'account'},
  ],
};

/// Descriptor for `GetAccessKeyInfoResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccessKeyInfoResponseDescriptor =
    $convert.base64Decode(
        'ChhHZXRBY2Nlc3NLZXlJbmZvUmVzcG9uc2USHAoHYWNjb3VudBj9xeLPASABKAlSB2FjY291bn'
        'Q=');

@$core.Deprecated('Use getCallerIdentityRequestDescriptor instead')
const GetCallerIdentityRequest$json = {
  '1': 'GetCallerIdentityRequest',
};

/// Descriptor for `GetCallerIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCallerIdentityRequestDescriptor =
    $convert.base64Decode('ChhHZXRDYWxsZXJJZGVudGl0eVJlcXVlc3Q=');

@$core.Deprecated('Use getCallerIdentityResponseDescriptor instead')
const GetCallerIdentityResponse$json = {
  '1': 'GetCallerIdentityResponse',
  '2': [
    {'1': 'account', '3': 435725053, '4': 1, '5': 9, '10': 'account'},
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'userid', '3': 10274112, '4': 1, '5': 9, '10': 'userid'},
  ],
};

/// Descriptor for `GetCallerIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getCallerIdentityResponseDescriptor =
    $convert.base64Decode(
        'ChlHZXRDYWxsZXJJZGVudGl0eVJlc3BvbnNlEhwKB2FjY291bnQY/cXizwEgASgJUgdhY2NvdW'
        '50EhQKA2Fybhidm+2/ASABKAlSA2FybhIZCgZ1c2VyaWQYwIrzBCABKAlSBnVzZXJpZA==');

@$core.Deprecated('Use getDelegatedAccessTokenRequestDescriptor instead')
const GetDelegatedAccessTokenRequest$json = {
  '1': 'GetDelegatedAccessTokenRequest',
  '2': [
    {'1': 'tradeintoken', '3': 10008816, '4': 1, '5': 9, '10': 'tradeintoken'},
  ],
};

/// Descriptor for `GetDelegatedAccessTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDelegatedAccessTokenRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXREZWxlZ2F0ZWRBY2Nlc3NUb2tlblJlcXVlc3QSJQoMdHJhZGVpbnRva2VuGPDx4gQgAS'
        'gJUgx0cmFkZWludG9rZW4=');

@$core.Deprecated('Use getDelegatedAccessTokenResponseDescriptor instead')
const GetDelegatedAccessTokenResponse$json = {
  '1': 'GetDelegatedAccessTokenResponse',
  '2': [
    {
      '1': 'assumedprincipal',
      '3': 359093742,
      '4': 1,
      '5': 9,
      '10': 'assumedprincipal'
    },
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.sts.Credentials',
      '10': 'credentials'
    },
    {
      '1': 'packedpolicysize',
      '3': 511234267,
      '4': 1,
      '5': 5,
      '10': 'packedpolicysize'
    },
  ],
};

/// Descriptor for `GetDelegatedAccessTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDelegatedAccessTokenResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXREZWxlZ2F0ZWRBY2Nlc3NUb2tlblJlc3BvbnNlEi4KEGFzc3VtZWRwcmluY2lwYWwY7q'
        'udqwEgASgJUhBhc3N1bWVkcHJpbmNpcGFsEjYKC2NyZWRlbnRpYWxzGPKajrYBIAEoCzIQLnN0'
        'cy5DcmVkZW50aWFsc1ILY3JlZGVudGlhbHMSLgoQcGFja2VkcG9saWN5c2l6ZRjboePzASABKA'
        'VSEHBhY2tlZHBvbGljeXNpemU=');

@$core.Deprecated('Use getFederationTokenRequestDescriptor instead')
const GetFederationTokenRequest$json = {
  '1': 'GetFederationTokenRequest',
  '2': [
    {
      '1': 'durationseconds',
      '3': 451873635,
      '4': 1,
      '5': 5,
      '10': 'durationseconds'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'policyarns',
      '3': 183785508,
      '4': 3,
      '5': 11,
      '6': '.sts.PolicyDescriptorType',
      '10': 'policyarns'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sts.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `GetFederationTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFederationTokenRequestDescriptor = $convert.base64Decode(
    'ChlHZXRGZWRlcmF0aW9uVG9rZW5SZXF1ZXN0EiwKD2R1cmF0aW9uc2Vjb25kcxjjlrzXASABKA'
    'VSD2R1cmF0aW9uc2Vjb25kcxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhoKBnBvbGljeRig7/Dg'
    'ASABKAlSBnBvbGljeRI8Cgpwb2xpY3lhcm5zGKSw0VcgAygLMhkuc3RzLlBvbGljeURlc2NyaX'
    'B0b3JUeXBlUgpwb2xpY3lhcm5zEiAKBHRhZ3MYwcH2tQEgAygLMgguc3RzLlRhZ1IEdGFncw==');

@$core.Deprecated('Use getFederationTokenResponseDescriptor instead')
const GetFederationTokenResponse$json = {
  '1': 'GetFederationTokenResponse',
  '2': [
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.sts.Credentials',
      '10': 'credentials'
    },
    {
      '1': 'federateduser',
      '3': 328666081,
      '4': 1,
      '5': 11,
      '6': '.sts.FederatedUser',
      '10': 'federateduser'
    },
    {
      '1': 'packedpolicysize',
      '3': 511234267,
      '4': 1,
      '5': 5,
      '10': 'packedpolicysize'
    },
  ],
};

/// Descriptor for `GetFederationTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getFederationTokenResponseDescriptor = $convert.base64Decode(
    'ChpHZXRGZWRlcmF0aW9uVG9rZW5SZXNwb25zZRI2CgtjcmVkZW50aWFscxjymo62ASABKAsyEC'
    '5zdHMuQ3JlZGVudGlhbHNSC2NyZWRlbnRpYWxzEjwKDWZlZGVyYXRlZHVzZXIY4ZfcnAEgASgL'
    'MhIuc3RzLkZlZGVyYXRlZFVzZXJSDWZlZGVyYXRlZHVzZXISLgoQcGFja2VkcG9saWN5c2l6ZR'
    'jboePzASABKAVSEHBhY2tlZHBvbGljeXNpemU=');

@$core.Deprecated('Use getSessionTokenRequestDescriptor instead')
const GetSessionTokenRequest$json = {
  '1': 'GetSessionTokenRequest',
  '2': [
    {
      '1': 'durationseconds',
      '3': 451873635,
      '4': 1,
      '5': 5,
      '10': 'durationseconds'
    },
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 9, '10': 'serialnumber'},
    {'1': 'tokencode', '3': 300671456, '4': 1, '5': 9, '10': 'tokencode'},
  ],
};

/// Descriptor for `GetSessionTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionTokenRequestDescriptor = $convert.base64Decode(
    'ChZHZXRTZXNzaW9uVG9rZW5SZXF1ZXN0EiwKD2R1cmF0aW9uc2Vjb25kcxjjlrzXASABKAVSD2'
    'R1cmF0aW9uc2Vjb25kcxImCgxzZXJpYWxudW1iZXIY5bq5xwEgASgJUgxzZXJpYWxudW1iZXIS'
    'IAoJdG9rZW5jb2RlGODDr48BIAEoCVIJdG9rZW5jb2Rl');

@$core.Deprecated('Use getSessionTokenResponseDescriptor instead')
const GetSessionTokenResponse$json = {
  '1': 'GetSessionTokenResponse',
  '2': [
    {
      '1': 'credentials',
      '3': 381914482,
      '4': 1,
      '5': 11,
      '6': '.sts.Credentials',
      '10': 'credentials'
    },
  ],
};

/// Descriptor for `GetSessionTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSessionTokenResponseDescriptor =
    $convert.base64Decode(
        'ChdHZXRTZXNzaW9uVG9rZW5SZXNwb25zZRI2CgtjcmVkZW50aWFscxjymo62ASABKAsyEC5zdH'
        'MuQ3JlZGVudGlhbHNSC2NyZWRlbnRpYWxz');

@$core.Deprecated('Use getWebIdentityTokenRequestDescriptor instead')
const GetWebIdentityTokenRequest$json = {
  '1': 'GetWebIdentityTokenRequest',
  '2': [
    {'1': 'audience', '3': 284892548, '4': 3, '5': 9, '10': 'audience'},
    {
      '1': 'durationseconds',
      '3': 451873635,
      '4': 1,
      '5': 5,
      '10': 'durationseconds'
    },
    {
      '1': 'signingalgorithm',
      '3': 488091842,
      '4': 1,
      '5': 9,
      '10': 'signingalgorithm'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sts.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `GetWebIdentityTokenRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWebIdentityTokenRequestDescriptor = $convert.base64Decode(
    'ChpHZXRXZWJJZGVudGl0eVRva2VuUmVxdWVzdBIeCghhdWRpZW5jZRiEu+yHASADKAlSCGF1ZG'
    'llbmNlEiwKD2R1cmF0aW9uc2Vjb25kcxjjlrzXASABKAVSD2R1cmF0aW9uc2Vjb25kcxIuChBz'
    'aWduaW5nYWxnb3JpdGhtGMLh3ugBIAEoCVIQc2lnbmluZ2FsZ29yaXRobRIgCgR0YWdzGMHB9r'
    'UBIAMoCzIILnN0cy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use getWebIdentityTokenResponseDescriptor instead')
const GetWebIdentityTokenResponse$json = {
  '1': 'GetWebIdentityTokenResponse',
  '2': [
    {'1': 'expiration', '3': 245879945, '4': 1, '5': 9, '10': 'expiration'},
    {
      '1': 'webidentitytoken',
      '3': 234014869,
      '4': 1,
      '5': 9,
      '10': 'webidentitytoken'
    },
  ],
};

/// Descriptor for `GetWebIdentityTokenResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getWebIdentityTokenResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRXZWJJZGVudGl0eVRva2VuUmVzcG9uc2USIQoKZXhwaXJhdGlvbhiJqZ91IAEoCVIKZX'
        'hwaXJhdGlvbhItChB3ZWJpZGVudGl0eXRva2VuGJWRy28gASgJUhB3ZWJpZGVudGl0eXRva2Vu');

@$core.Deprecated('Use iDPCommunicationErrorExceptionDescriptor instead')
const IDPCommunicationErrorException$json = {
  '1': 'IDPCommunicationErrorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IDPCommunicationErrorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List iDPCommunicationErrorExceptionDescriptor =
    $convert.base64Decode(
        'Ch5JRFBDb21tdW5pY2F0aW9uRXJyb3JFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use iDPRejectedClaimExceptionDescriptor instead')
const IDPRejectedClaimException$json = {
  '1': 'IDPRejectedClaimException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IDPRejectedClaimException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List iDPRejectedClaimExceptionDescriptor =
    $convert.base64Decode(
        'ChlJRFBSZWplY3RlZENsYWltRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidAuthorizationMessageExceptionDescriptor instead')
const InvalidAuthorizationMessageException$json = {
  '1': 'InvalidAuthorizationMessageException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAuthorizationMessageException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAuthorizationMessageExceptionDescriptor =
    $convert.base64Decode(
        'CiRJbnZhbGlkQXV0aG9yaXphdGlvbk1lc3NhZ2VFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIA'
        'EoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidIdentityTokenExceptionDescriptor instead')
const InvalidIdentityTokenException$json = {
  '1': 'InvalidIdentityTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidIdentityTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidIdentityTokenExceptionDescriptor =
    $convert.base64Decode(
        'Ch1JbnZhbGlkSWRlbnRpdHlUb2tlbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use jWTPayloadSizeExceededExceptionDescriptor instead')
const JWTPayloadSizeExceededException$json = {
  '1': 'JWTPayloadSizeExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `JWTPayloadSizeExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List jWTPayloadSizeExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch9KV1RQYXlsb2FkU2l6ZUV4Y2VlZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use malformedPolicyDocumentExceptionDescriptor instead')
const MalformedPolicyDocumentException$json = {
  '1': 'MalformedPolicyDocumentException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MalformedPolicyDocumentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List malformedPolicyDocumentExceptionDescriptor =
    $convert.base64Decode(
        'CiBNYWxmb3JtZWRQb2xpY3lEb2N1bWVudEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated(
    'Use outboundWebIdentityFederationDisabledExceptionDescriptor instead')
const OutboundWebIdentityFederationDisabledException$json = {
  '1': 'OutboundWebIdentityFederationDisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OutboundWebIdentityFederationDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    outboundWebIdentityFederationDisabledExceptionDescriptor =
    $convert.base64Decode(
        'Ci5PdXRib3VuZFdlYklkZW50aXR5RmVkZXJhdGlvbkRpc2FibGVkRXhjZXB0aW9uEhsKB21lc3'
        'NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use packedPolicyTooLargeExceptionDescriptor instead')
const PackedPolicyTooLargeException$json = {
  '1': 'PackedPolicyTooLargeException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PackedPolicyTooLargeException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List packedPolicyTooLargeExceptionDescriptor =
    $convert.base64Decode(
        'Ch1QYWNrZWRQb2xpY3lUb29MYXJnZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use policyDescriptorTypeDescriptor instead')
const PolicyDescriptorType$json = {
  '1': 'PolicyDescriptorType',
  '2': [
    {'1': 'arn', '3': 359604989, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `PolicyDescriptorType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List policyDescriptorTypeDescriptor =
    $convert.base64Decode(
        'ChRQb2xpY3lEZXNjcmlwdG9yVHlwZRIUCgNhcm4Y/cW8qwEgASgJUgNhcm4=');

@$core.Deprecated('Use providedContextDescriptor instead')
const ProvidedContext$json = {
  '1': 'ProvidedContext',
  '2': [
    {
      '1': 'contextassertion',
      '3': 351907089,
      '4': 1,
      '5': 9,
      '10': 'contextassertion'
    },
    {'1': 'providerarn', '3': 426083188, '4': 1, '5': 9, '10': 'providerarn'},
  ],
};

/// Descriptor for `ProvidedContext`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List providedContextDescriptor = $convert.base64Decode(
    'Cg9Qcm92aWRlZENvbnRleHQSLgoQY29udGV4dGFzc2VydGlvbhiR2uanASABKAlSEGNvbnRleH'
    'Rhc3NlcnRpb24SJAoLcHJvdmlkZXJhcm4Y9IaWywEgASgJUgtwcm92aWRlcmFybg==');

@$core.Deprecated('Use regionDisabledExceptionDescriptor instead')
const RegionDisabledException$json = {
  '1': 'RegionDisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RegionDisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List regionDisabledExceptionDescriptor =
    $convert.base64Decode(
        'ChdSZWdpb25EaXNhYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use sessionDurationEscalationExceptionDescriptor instead')
const SessionDurationEscalationException$json = {
  '1': 'SessionDurationEscalationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SessionDurationEscalationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sessionDurationEscalationExceptionDescriptor =
    $convert.base64Decode(
        'CiJTZXNzaW9uRHVyYXRpb25Fc2NhbGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
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
