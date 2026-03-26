// This is a generated file - do not edit.
//
// Generated from kinesis.proto.

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

@$core.Deprecated('Use consumerStatusDescriptor instead')
const ConsumerStatus$json = {
  '1': 'ConsumerStatus',
  '2': [
    {'1': 'CONSUMER_STATUS_ACTIVE', '2': 0},
    {'1': 'CONSUMER_STATUS_DELETING', '2': 1},
    {'1': 'CONSUMER_STATUS_CREATING', '2': 2},
  ],
};

/// Descriptor for `ConsumerStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List consumerStatusDescriptor = $convert.base64Decode(
    'Cg5Db25zdW1lclN0YXR1cxIaChZDT05TVU1FUl9TVEFUVVNfQUNUSVZFEAASHAoYQ09OU1VNRV'
    'JfU1RBVFVTX0RFTEVUSU5HEAESHAoYQ09OU1VNRVJfU1RBVFVTX0NSRUFUSU5HEAI=');

@$core.Deprecated('Use encryptionTypeDescriptor instead')
const EncryptionType$json = {
  '1': 'EncryptionType',
  '2': [
    {'1': 'ENCRYPTION_TYPE_NONE', '2': 0},
    {'1': 'ENCRYPTION_TYPE_KMS', '2': 1},
  ],
};

/// Descriptor for `EncryptionType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List encryptionTypeDescriptor = $convert.base64Decode(
    'Cg5FbmNyeXB0aW9uVHlwZRIYChRFTkNSWVBUSU9OX1RZUEVfTk9ORRAAEhcKE0VOQ1JZUFRJT0'
    '5fVFlQRV9LTVMQAQ==');

@$core.Deprecated('Use metricsNameDescriptor instead')
const MetricsName$json = {
  '1': 'MetricsName',
  '2': [
    {'1': 'METRICS_NAME_INCOMING_RECORDS', '2': 0},
    {'1': 'METRICS_NAME_INCOMING_BYTES', '2': 1},
    {'1': 'METRICS_NAME_WRITE_PROVISIONED_THROUGHPUT_EXCEEDED', '2': 2},
    {'1': 'METRICS_NAME_READ_PROVISIONED_THROUGHPUT_EXCEEDED', '2': 3},
    {'1': 'METRICS_NAME_OUTGOING_BYTES', '2': 4},
    {'1': 'METRICS_NAME_OUTGOING_RECORDS', '2': 5},
    {'1': 'METRICS_NAME_ITERATOR_AGE_MILLISECONDS', '2': 6},
    {'1': 'METRICS_NAME_ALL', '2': 7},
  ],
};

/// Descriptor for `MetricsName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricsNameDescriptor = $convert.base64Decode(
    'CgtNZXRyaWNzTmFtZRIhCh1NRVRSSUNTX05BTUVfSU5DT01JTkdfUkVDT1JEUxAAEh8KG01FVF'
    'JJQ1NfTkFNRV9JTkNPTUlOR19CWVRFUxABEjYKMk1FVFJJQ1NfTkFNRV9XUklURV9QUk9WSVNJ'
    'T05FRF9USFJPVUdIUFVUX0VYQ0VFREVEEAISNQoxTUVUUklDU19OQU1FX1JFQURfUFJPVklTSU'
    '9ORURfVEhST1VHSFBVVF9FWENFRURFRBADEh8KG01FVFJJQ1NfTkFNRV9PVVRHT0lOR19CWVRF'
    'UxAEEiEKHU1FVFJJQ1NfTkFNRV9PVVRHT0lOR19SRUNPUkRTEAUSKgomTUVUUklDU19OQU1FX0'
    'lURVJBVE9SX0FHRV9NSUxMSVNFQ09ORFMQBhIUChBNRVRSSUNTX05BTUVfQUxMEAc=');

@$core.Deprecated(
    'Use minimumThroughputBillingCommitmentInputStatusDescriptor instead')
const MinimumThroughputBillingCommitmentInputStatus$json = {
  '1': 'MinimumThroughputBillingCommitmentInputStatus',
  '2': [
    {
      '1': 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_DISABLED',
      '2': 0
    },
    {'1': 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_INPUT_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `MinimumThroughputBillingCommitmentInputStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List
    minimumThroughputBillingCommitmentInputStatusDescriptor =
    $convert.base64Decode(
        'Ci1NaW5pbXVtVGhyb3VnaHB1dEJpbGxpbmdDb21taXRtZW50SW5wdXRTdGF0dXMSPwo7TUlOSU'
        '1VTV9USFJPVUdIUFVUX0JJTExJTkdfQ09NTUlUTUVOVF9JTlBVVF9TVEFUVVNfRElTQUJMRUQQ'
        'ABI+CjpNSU5JTVVNX1RIUk9VR0hQVVRfQklMTElOR19DT01NSVRNRU5UX0lOUFVUX1NUQVRVU1'
        '9FTkFCTEVEEAE=');

@$core.Deprecated(
    'Use minimumThroughputBillingCommitmentOutputStatusDescriptor instead')
const MinimumThroughputBillingCommitmentOutputStatus$json = {
  '1': 'MinimumThroughputBillingCommitmentOutputStatus',
  '2': [
    {
      '1': 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_DISABLED',
      '2': 0
    },
    {
      '1': 'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED',
      '2': 1
    },
    {
      '1':
          'MINIMUM_THROUGHPUT_BILLING_COMMITMENT_OUTPUT_STATUS_ENABLED_UNTIL_EARLIEST_ALLOWED_END',
      '2': 2
    },
  ],
};

/// Descriptor for `MinimumThroughputBillingCommitmentOutputStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List
    minimumThroughputBillingCommitmentOutputStatusDescriptor =
    $convert.base64Decode(
        'Ci5NaW5pbXVtVGhyb3VnaHB1dEJpbGxpbmdDb21taXRtZW50T3V0cHV0U3RhdHVzEkAKPE1JTk'
        'lNVU1fVEhST1VHSFBVVF9CSUxMSU5HX0NPTU1JVE1FTlRfT1VUUFVUX1NUQVRVU19ESVNBQkxF'
        'RBAAEj8KO01JTklNVU1fVEhST1VHSFBVVF9CSUxMSU5HX0NPTU1JVE1FTlRfT1VUUFVUX1NUQV'
        'RVU19FTkFCTEVEEAESWgpWTUlOSU1VTV9USFJPVUdIUFVUX0JJTExJTkdfQ09NTUlUTUVOVF9P'
        'VVRQVVRfU1RBVFVTX0VOQUJMRURfVU5USUxfRUFSTElFU1RfQUxMT1dFRF9FTkQQAg==');

@$core.Deprecated('Use scalingTypeDescriptor instead')
const ScalingType$json = {
  '1': 'ScalingType',
  '2': [
    {'1': 'SCALING_TYPE_UNIFORM_SCALING', '2': 0},
  ],
};

/// Descriptor for `ScalingType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scalingTypeDescriptor = $convert.base64Decode(
    'CgtTY2FsaW5nVHlwZRIgChxTQ0FMSU5HX1RZUEVfVU5JRk9STV9TQ0FMSU5HEAA=');

@$core.Deprecated('Use shardFilterTypeDescriptor instead')
const ShardFilterType$json = {
  '1': 'ShardFilterType',
  '2': [
    {'1': 'SHARD_FILTER_TYPE_AT_TIMESTAMP', '2': 0},
    {'1': 'SHARD_FILTER_TYPE_FROM_TRIM_HORIZON', '2': 1},
    {'1': 'SHARD_FILTER_TYPE_AT_LATEST', '2': 2},
    {'1': 'SHARD_FILTER_TYPE_AT_TRIM_HORIZON', '2': 3},
    {'1': 'SHARD_FILTER_TYPE_AFTER_SHARD_ID', '2': 4},
    {'1': 'SHARD_FILTER_TYPE_FROM_TIMESTAMP', '2': 5},
  ],
};

/// Descriptor for `ShardFilterType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List shardFilterTypeDescriptor = $convert.base64Decode(
    'Cg9TaGFyZEZpbHRlclR5cGUSIgoeU0hBUkRfRklMVEVSX1RZUEVfQVRfVElNRVNUQU1QEAASJw'
    'ojU0hBUkRfRklMVEVSX1RZUEVfRlJPTV9UUklNX0hPUklaT04QARIfChtTSEFSRF9GSUxURVJf'
    'VFlQRV9BVF9MQVRFU1QQAhIlCiFTSEFSRF9GSUxURVJfVFlQRV9BVF9UUklNX0hPUklaT04QAx'
    'IkCiBTSEFSRF9GSUxURVJfVFlQRV9BRlRFUl9TSEFSRF9JRBAEEiQKIFNIQVJEX0ZJTFRFUl9U'
    'WVBFX0ZST01fVElNRVNUQU1QEAU=');

@$core.Deprecated('Use shardIteratorTypeDescriptor instead')
const ShardIteratorType$json = {
  '1': 'ShardIteratorType',
  '2': [
    {'1': 'SHARD_ITERATOR_TYPE_AT_TIMESTAMP', '2': 0},
    {'1': 'SHARD_ITERATOR_TYPE_TRIM_HORIZON', '2': 1},
    {'1': 'SHARD_ITERATOR_TYPE_AFTER_SEQUENCE_NUMBER', '2': 2},
    {'1': 'SHARD_ITERATOR_TYPE_AT_SEQUENCE_NUMBER', '2': 3},
    {'1': 'SHARD_ITERATOR_TYPE_LATEST', '2': 4},
  ],
};

/// Descriptor for `ShardIteratorType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List shardIteratorTypeDescriptor = $convert.base64Decode(
    'ChFTaGFyZEl0ZXJhdG9yVHlwZRIkCiBTSEFSRF9JVEVSQVRPUl9UWVBFX0FUX1RJTUVTVEFNUB'
    'AAEiQKIFNIQVJEX0lURVJBVE9SX1RZUEVfVFJJTV9IT1JJWk9OEAESLQopU0hBUkRfSVRFUkFU'
    'T1JfVFlQRV9BRlRFUl9TRVFVRU5DRV9OVU1CRVIQAhIqCiZTSEFSRF9JVEVSQVRPUl9UWVBFX0'
    'FUX1NFUVVFTkNFX05VTUJFUhADEh4KGlNIQVJEX0lURVJBVE9SX1RZUEVfTEFURVNUEAQ=');

@$core.Deprecated('Use streamModeDescriptor instead')
const StreamMode$json = {
  '1': 'StreamMode',
  '2': [
    {'1': 'STREAM_MODE_PROVISIONED', '2': 0},
    {'1': 'STREAM_MODE_ON_DEMAND', '2': 1},
  ],
};

/// Descriptor for `StreamMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List streamModeDescriptor = $convert.base64Decode(
    'CgpTdHJlYW1Nb2RlEhsKF1NUUkVBTV9NT0RFX1BST1ZJU0lPTkVEEAASGQoVU1RSRUFNX01PRE'
    'VfT05fREVNQU5EEAE=');

@$core.Deprecated('Use streamStatusDescriptor instead')
const StreamStatus$json = {
  '1': 'StreamStatus',
  '2': [
    {'1': 'STREAM_STATUS_UPDATING', '2': 0},
    {'1': 'STREAM_STATUS_ACTIVE', '2': 1},
    {'1': 'STREAM_STATUS_DELETING', '2': 2},
    {'1': 'STREAM_STATUS_CREATING', '2': 3},
  ],
};

/// Descriptor for `StreamStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List streamStatusDescriptor = $convert.base64Decode(
    'CgxTdHJlYW1TdGF0dXMSGgoWU1RSRUFNX1NUQVRVU19VUERBVElORxAAEhgKFFNUUkVBTV9TVE'
    'FUVVNfQUNUSVZFEAESGgoWU1RSRUFNX1NUQVRVU19ERUxFVElORxACEhoKFlNUUkVBTV9TVEFU'
    'VVNfQ1JFQVRJTkcQAw==');

@$core.Deprecated('Use accessDeniedExceptionDescriptor instead')
const AccessDeniedException$json = {
  '1': 'AccessDeniedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccessDeniedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accessDeniedExceptionDescriptor = $convert.base64Decode(
    'ChVBY2Nlc3NEZW5pZWRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use addTagsToStreamInputDescriptor instead')
const AddTagsToStreamInput$json = {
  '1': 'AddTagsToStreamInput',
  '2': [
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kinesis.AddTagsToStreamInput.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [AddTagsToStreamInput_TagsEntry$json],
};

@$core.Deprecated('Use addTagsToStreamInputDescriptor instead')
const AddTagsToStreamInput_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `AddTagsToStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addTagsToStreamInputDescriptor = $convert.base64Decode(
    'ChRBZGRUYWdzVG9TdHJlYW1JbnB1dBIgCglzdHJlYW1hcm4Y3fOq8gEgASgJUglzdHJlYW1hcm'
    '4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZBIiCgpzdHJlYW1uYW1lGMe3ueABIAEo'
    'CVIKc3RyZWFtbmFtZRI/CgR0YWdzGMHB9rUBIAMoCzInLmtpbmVzaXMuQWRkVGFnc1RvU3RyZW'
    'FtSW5wdXQuVGFnc0VudHJ5UgR0YWdzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIU'
    'CgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use childShardDescriptor instead')
const ChildShard$json = {
  '1': 'ChildShard',
  '2': [
    {
      '1': 'hashkeyrange',
      '3': 981486,
      '4': 1,
      '5': 11,
      '6': '.kinesis.HashKeyRange',
      '10': 'hashkeyrange'
    },
    {'1': 'parentshards', '3': 152621201, '4': 3, '5': 9, '10': 'parentshards'},
    {'1': 'shardid', '3': 66410951, '4': 1, '5': 9, '10': 'shardid'},
  ],
};

/// Descriptor for `ChildShard`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List childShardDescriptor = $convert.base64Decode(
    'CgpDaGlsZFNoYXJkEjsKDGhhc2hrZXlyYW5nZRju8zsgASgLMhUua2luZXNpcy5IYXNoS2V5Um'
    'FuZ2VSDGhhc2hrZXlyYW5nZRIlCgxwYXJlbnRzaGFyZHMYkaHjSCADKAlSDHBhcmVudHNoYXJk'
    'cxIbCgdzaGFyZGlkGMez1R8gASgJUgdzaGFyZGlk');

@$core.Deprecated('Use consumerDescriptor instead')
const Consumer$json = {
  '1': 'Consumer',
  '2': [
    {'1': 'consumerarn', '3': 41107441, '4': 1, '5': 9, '10': 'consumerarn'},
    {
      '1': 'consumercreationtimestamp',
      '3': 356226265,
      '4': 1,
      '5': 9,
      '10': 'consumercreationtimestamp'
    },
    {'1': 'consumername', '3': 70979235, '4': 1, '5': 9, '10': 'consumername'},
    {
      '1': 'consumerstatus',
      '3': 33471404,
      '4': 1,
      '5': 14,
      '6': '.kinesis.ConsumerStatus',
      '10': 'consumerstatus'
    },
  ],
};

/// Descriptor for `Consumer`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List consumerDescriptor = $convert.base64Decode(
    'CghDb25zdW1lchIjCgtjb25zdW1lcmFybhjx/8wTIAEoCVILY29uc3VtZXJhcm4SQAoZY29uc3'
    'VtZXJjcmVhdGlvbnRpbWVzdGFtcBjZqe6pASABKAlSGWNvbnN1bWVyY3JlYXRpb250aW1lc3Rh'
    'bXASJQoMY29uc3VtZXJuYW1lGKOd7CEgASgJUgxjb25zdW1lcm5hbWUSQgoOY29uc3VtZXJzdG'
    'F0dXMYrPf6DyABKA4yFy5raW5lc2lzLkNvbnN1bWVyU3RhdHVzUg5jb25zdW1lcnN0YXR1cw==');

@$core.Deprecated('Use consumerDescriptionDescriptor instead')
const ConsumerDescription$json = {
  '1': 'ConsumerDescription',
  '2': [
    {'1': 'consumerarn', '3': 41107441, '4': 1, '5': 9, '10': 'consumerarn'},
    {
      '1': 'consumercreationtimestamp',
      '3': 356226265,
      '4': 1,
      '5': 9,
      '10': 'consumercreationtimestamp'
    },
    {'1': 'consumername', '3': 70979235, '4': 1, '5': 9, '10': 'consumername'},
    {
      '1': 'consumerstatus',
      '3': 33471404,
      '4': 1,
      '5': 14,
      '6': '.kinesis.ConsumerStatus',
      '10': 'consumerstatus'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
  ],
};

/// Descriptor for `ConsumerDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List consumerDescriptionDescriptor = $convert.base64Decode(
    'ChNDb25zdW1lckRlc2NyaXB0aW9uEiMKC2NvbnN1bWVyYXJuGPH/zBMgASgJUgtjb25zdW1lcm'
    'FybhJAChljb25zdW1lcmNyZWF0aW9udGltZXN0YW1wGNmp7qkBIAEoCVIZY29uc3VtZXJjcmVh'
    'dGlvbnRpbWVzdGFtcBIlCgxjb25zdW1lcm5hbWUYo53sISABKAlSDGNvbnN1bWVybmFtZRJCCg'
    '5jb25zdW1lcnN0YXR1cxis9/oPIAEoDjIXLmtpbmVzaXMuQ29uc3VtZXJTdGF0dXNSDmNvbnN1'
    'bWVyc3RhdHVzEiAKCXN0cmVhbWFybhjd86ryASABKAlSCXN0cmVhbWFybg==');

@$core.Deprecated('Use createStreamInputDescriptor instead')
const CreateStreamInput$json = {
  '1': 'CreateStreamInput',
  '2': [
    {
      '1': 'maxrecordsizeinkib',
      '3': 197267253,
      '4': 1,
      '5': 5,
      '10': 'maxrecordsizeinkib'
    },
    {'1': 'shardcount', '3': 92831547, '4': 1, '5': 5, '10': 'shardcount'},
    {
      '1': 'streammodedetails',
      '3': 12139665,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StreamModeDetails',
      '10': 'streammodedetails'
    },
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kinesis.CreateStreamInput.TagsEntry',
      '10': 'tags'
    },
    {
      '1': 'warmthroughputmibps',
      '3': 259219704,
      '4': 1,
      '5': 5,
      '10': 'warmthroughputmibps'
    },
  ],
  '3': [CreateStreamInput_TagsEntry$json],
};

@$core.Deprecated('Use createStreamInputDescriptor instead')
const CreateStreamInput_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createStreamInputDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVTdHJlYW1JbnB1dBIxChJtYXhyZWNvcmRzaXplaW5raWIYtZ6IXiABKAVSEm1heH'
    'JlY29yZHNpemVpbmtpYhIhCgpzaGFyZGNvdW50GLv+oSwgASgFUgpzaGFyZGNvdW50EksKEXN0'
    'cmVhbW1vZGVkZXRhaWxzGJH55AUgASgLMhoua2luZXNpcy5TdHJlYW1Nb2RlRGV0YWlsc1IRc3'
    'RyZWFtbW9kZWRldGFpbHMSIgoKc3RyZWFtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWUSPAoE'
    'dGFncxjBwfa1ASADKAsyJC5raW5lc2lzLkNyZWF0ZVN0cmVhbUlucHV0LlRhZ3NFbnRyeVIEdG'
    'FncxIzChN3YXJtdGhyb3VnaHB1dG1pYnBzGPjBzXsgASgFUhN3YXJtdGhyb3VnaHB1dG1pYnBz'
    'GjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6Aj'
    'gB');

@$core.Deprecated('Use decreaseStreamRetentionPeriodInputDescriptor instead')
const DecreaseStreamRetentionPeriodInput$json = {
  '1': 'DecreaseStreamRetentionPeriodInput',
  '2': [
    {
      '1': 'retentionperiodhours',
      '3': 396381944,
      '4': 1,
      '5': 5,
      '10': 'retentionperiodhours'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `DecreaseStreamRetentionPeriodInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List decreaseStreamRetentionPeriodInputDescriptor =
    $convert.base64Decode(
        'CiJEZWNyZWFzZVN0cmVhbVJldGVudGlvblBlcmlvZElucHV0EjYKFHJldGVudGlvbnBlcmlvZG'
        'hvdXJzGPidgb0BIAEoBVIUcmV0ZW50aW9ucGVyaW9kaG91cnMSIAoJc3RyZWFtYXJuGN3zqvIB'
        'IAEoCVIJc3RyZWFtYXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSIgoKc3RyZW'
        'FtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWU=');

@$core.Deprecated('Use deleteResourcePolicyInputDescriptor instead')
const DeleteResourcePolicyInput$json = {
  '1': 'DeleteResourcePolicyInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `DeleteResourcePolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteResourcePolicyInputDescriptor =
    $convert.base64Decode(
        'ChlEZWxldGVSZXNvdXJjZVBvbGljeUlucHV0EiQKC3Jlc291cmNlYXJuGO3AmbABIAEoCVILcm'
        'Vzb3VyY2Vhcm4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZA==');

@$core.Deprecated('Use deleteStreamInputDescriptor instead')
const DeleteStreamInput$json = {
  '1': 'DeleteStreamInput',
  '2': [
    {
      '1': 'enforceconsumerdeletion',
      '3': 513184424,
      '4': 1,
      '5': 8,
      '10': 'enforceconsumerdeletion'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `DeleteStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteStreamInputDescriptor = $convert.base64Decode(
    'ChFEZWxldGVTdHJlYW1JbnB1dBI8ChdlbmZvcmNlY29uc3VtZXJkZWxldGlvbhiopdr0ASABKA'
    'hSF2VuZm9yY2Vjb25zdW1lcmRlbGV0aW9uEiAKCXN0cmVhbWFybhjd86ryASABKAlSCXN0cmVh'
    'bWFybhIeCghzdHJlYW1pZBjB7YHGASABKAlSCHN0cmVhbWlkEiIKCnN0cmVhbW5hbWUYx7e54A'
    'EgASgJUgpzdHJlYW1uYW1l');

@$core.Deprecated('Use deregisterStreamConsumerInputDescriptor instead')
const DeregisterStreamConsumerInput$json = {
  '1': 'DeregisterStreamConsumerInput',
  '2': [
    {'1': 'consumerarn', '3': 41107441, '4': 1, '5': 9, '10': 'consumerarn'},
    {'1': 'consumername', '3': 70979235, '4': 1, '5': 9, '10': 'consumername'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `DeregisterStreamConsumerInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deregisterStreamConsumerInputDescriptor = $convert.base64Decode(
    'Ch1EZXJlZ2lzdGVyU3RyZWFtQ29uc3VtZXJJbnB1dBIjCgtjb25zdW1lcmFybhjx/8wTIAEoCV'
    'ILY29uc3VtZXJhcm4SJQoMY29uc3VtZXJuYW1lGKOd7CEgASgJUgxjb25zdW1lcm5hbWUSIAoJ'
    'c3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3'
    'RyZWFtaWQ=');

@$core.Deprecated('Use describeAccountSettingsInputDescriptor instead')
const DescribeAccountSettingsInput$json = {
  '1': 'DescribeAccountSettingsInput',
};

/// Descriptor for `DescribeAccountSettingsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAccountSettingsInputDescriptor =
    $convert.base64Decode('ChxEZXNjcmliZUFjY291bnRTZXR0aW5nc0lucHV0');

@$core.Deprecated('Use describeAccountSettingsOutputDescriptor instead')
const DescribeAccountSettingsOutput$json = {
  '1': 'DescribeAccountSettingsOutput',
  '2': [
    {
      '1': 'minimumthroughputbillingcommitment',
      '3': 335326962,
      '4': 1,
      '5': 11,
      '6': '.kinesis.MinimumThroughputBillingCommitmentOutput',
      '10': 'minimumthroughputbillingcommitment'
    },
  ],
};

/// Descriptor for `DescribeAccountSettingsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeAccountSettingsOutputDescriptor = $convert.base64Decode(
    'Ch1EZXNjcmliZUFjY291bnRTZXR0aW5nc091dHB1dBKFAQoibWluaW11bXRocm91Z2hwdXRiaW'
    'xsaW5nY29tbWl0bWVudBjy3fKfASABKAsyMS5raW5lc2lzLk1pbmltdW1UaHJvdWdocHV0Qmls'
    'bGluZ0NvbW1pdG1lbnRPdXRwdXRSIm1pbmltdW10aHJvdWdocHV0YmlsbGluZ2NvbW1pdG1lbn'
    'Q=');

@$core.Deprecated('Use describeLimitsInputDescriptor instead')
const DescribeLimitsInput$json = {
  '1': 'DescribeLimitsInput',
};

/// Descriptor for `DescribeLimitsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLimitsInputDescriptor =
    $convert.base64Decode('ChNEZXNjcmliZUxpbWl0c0lucHV0');

@$core.Deprecated('Use describeLimitsOutputDescriptor instead')
const DescribeLimitsOutput$json = {
  '1': 'DescribeLimitsOutput',
  '2': [
    {
      '1': 'ondemandstreamcount',
      '3': 386875707,
      '4': 1,
      '5': 5,
      '10': 'ondemandstreamcount'
    },
    {
      '1': 'ondemandstreamcountlimit',
      '3': 38712458,
      '4': 1,
      '5': 5,
      '10': 'ondemandstreamcountlimit'
    },
    {
      '1': 'openshardcount',
      '3': 476409287,
      '4': 1,
      '5': 5,
      '10': 'openshardcount'
    },
    {'1': 'shardlimit', '3': 515722839, '4': 1, '5': 5, '10': 'shardlimit'},
  ],
};

/// Descriptor for `DescribeLimitsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeLimitsOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZUxpbWl0c091dHB1dBI0ChNvbmRlbWFuZHN0cmVhbWNvdW50GLuCvbgBIAEoBV'
    'ITb25kZW1hbmRzdHJlYW1jb3VudBI9ChhvbmRlbWFuZHN0cmVhbWNvdW50bGltaXQYium6EiAB'
    'KAVSGG9uZGVtYW5kc3RyZWFtY291bnRsaW1pdBIqCg5vcGVuc2hhcmRjb3VudBjH25XjASABKA'
    'VSDm9wZW5zaGFyZGNvdW50EiIKCnNoYXJkbGltaXQY15z19QEgASgFUgpzaGFyZGxpbWl0');

@$core.Deprecated('Use describeStreamConsumerInputDescriptor instead')
const DescribeStreamConsumerInput$json = {
  '1': 'DescribeStreamConsumerInput',
  '2': [
    {'1': 'consumerarn', '3': 41107441, '4': 1, '5': 9, '10': 'consumerarn'},
    {'1': 'consumername', '3': 70979235, '4': 1, '5': 9, '10': 'consumername'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `DescribeStreamConsumerInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStreamConsumerInputDescriptor = $convert.base64Decode(
    'ChtEZXNjcmliZVN0cmVhbUNvbnN1bWVySW5wdXQSIwoLY29uc3VtZXJhcm4Y8f/MEyABKAlSC2'
    'NvbnN1bWVyYXJuEiUKDGNvbnN1bWVybmFtZRijnewhIAEoCVIMY29uc3VtZXJuYW1lEiAKCXN0'
    'cmVhbWFybhjd86ryASABKAlSCXN0cmVhbWFybhIeCghzdHJlYW1pZBjB7YHGASABKAlSCHN0cm'
    'VhbWlk');

@$core.Deprecated('Use describeStreamConsumerOutputDescriptor instead')
const DescribeStreamConsumerOutput$json = {
  '1': 'DescribeStreamConsumerOutput',
  '2': [
    {
      '1': 'consumerdescription',
      '3': 524123998,
      '4': 1,
      '5': 11,
      '6': '.kinesis.ConsumerDescription',
      '10': 'consumerdescription'
    },
  ],
};

/// Descriptor for `DescribeStreamConsumerOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStreamConsumerOutputDescriptor =
    $convert.base64Decode(
        'ChxEZXNjcmliZVN0cmVhbUNvbnN1bWVyT3V0cHV0ElIKE2NvbnN1bWVyZGVzY3JpcHRpb24Y3v'
        '71+QEgASgLMhwua2luZXNpcy5Db25zdW1lckRlc2NyaXB0aW9uUhNjb25zdW1lcmRlc2NyaXB0'
        'aW9u');

@$core.Deprecated('Use describeStreamInputDescriptor instead')
const DescribeStreamInput$json = {
  '1': 'DescribeStreamInput',
  '2': [
    {
      '1': 'exclusivestartshardid',
      '3': 44771587,
      '4': 1,
      '5': 9,
      '10': 'exclusivestartshardid'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `DescribeStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStreamInputDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZVN0cmVhbUlucHV0EjcKFWV4Y2x1c2l2ZXN0YXJ0c2hhcmRpZBiD0qwVIAEoCV'
    'IVZXhjbHVzaXZlc3RhcnRzaGFyZGlkEhgKBWxpbWl0GNWV2cQBIAEoBVIFbGltaXQSIAoJc3Ry'
    'ZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZW'
    'FtaWQSIgoKc3RyZWFtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWU=');

@$core.Deprecated('Use describeStreamOutputDescriptor instead')
const DescribeStreamOutput$json = {
  '1': 'DescribeStreamOutput',
  '2': [
    {
      '1': 'streamdescription',
      '3': 363737034,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StreamDescription',
      '10': 'streamdescription'
    },
  ],
};

/// Descriptor for `DescribeStreamOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStreamOutputDescriptor = $convert.base64Decode(
    'ChREZXNjcmliZVN0cmVhbU91dHB1dBJMChFzdHJlYW1kZXNjcmlwdGlvbhjK37itASABKAsyGi'
    '5raW5lc2lzLlN0cmVhbURlc2NyaXB0aW9uUhFzdHJlYW1kZXNjcmlwdGlvbg==');

@$core.Deprecated('Use describeStreamSummaryInputDescriptor instead')
const DescribeStreamSummaryInput$json = {
  '1': 'DescribeStreamSummaryInput',
  '2': [
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `DescribeStreamSummaryInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStreamSummaryInputDescriptor =
    $convert.base64Decode(
        'ChpEZXNjcmliZVN0cmVhbVN1bW1hcnlJbnB1dBIgCglzdHJlYW1hcm4Y3fOq8gEgASgJUglzdH'
        'JlYW1hcm4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZBIiCgpzdHJlYW1uYW1lGMe3'
        'ueABIAEoCVIKc3RyZWFtbmFtZQ==');

@$core.Deprecated('Use describeStreamSummaryOutputDescriptor instead')
const DescribeStreamSummaryOutput$json = {
  '1': 'DescribeStreamSummaryOutput',
  '2': [
    {
      '1': 'streamdescriptionsummary',
      '3': 478899272,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StreamDescriptionSummary',
      '10': 'streamdescriptionsummary'
    },
  ],
};

/// Descriptor for `DescribeStreamSummaryOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeStreamSummaryOutputDescriptor =
    $convert.base64Decode(
        'ChtEZXNjcmliZVN0cmVhbVN1bW1hcnlPdXRwdXQSYQoYc3RyZWFtZGVzY3JpcHRpb25zdW1tYX'
        'J5GMjYreQBIAEoCzIhLmtpbmVzaXMuU3RyZWFtRGVzY3JpcHRpb25TdW1tYXJ5UhhzdHJlYW1k'
        'ZXNjcmlwdGlvbnN1bW1hcnk=');

@$core.Deprecated('Use disableEnhancedMonitoringInputDescriptor instead')
const DisableEnhancedMonitoringInput$json = {
  '1': 'DisableEnhancedMonitoringInput',
  '2': [
    {
      '1': 'shardlevelmetrics',
      '3': 406711021,
      '4': 3,
      '5': 14,
      '6': '.kinesis.MetricsName',
      '10': 'shardlevelmetrics'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `DisableEnhancedMonitoringInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableEnhancedMonitoringInputDescriptor =
    $convert.base64Decode(
        'Ch5EaXNhYmxlRW5oYW5jZWRNb25pdG9yaW5nSW5wdXQSRgoRc2hhcmRsZXZlbG1ldHJpY3MY7d'
        'X3wQEgAygOMhQua2luZXNpcy5NZXRyaWNzTmFtZVIRc2hhcmRsZXZlbG1ldHJpY3MSIAoJc3Ry'
        'ZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZW'
        'FtaWQSIgoKc3RyZWFtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWU=');

@$core.Deprecated('Use enableEnhancedMonitoringInputDescriptor instead')
const EnableEnhancedMonitoringInput$json = {
  '1': 'EnableEnhancedMonitoringInput',
  '2': [
    {
      '1': 'shardlevelmetrics',
      '3': 406711021,
      '4': 3,
      '5': 14,
      '6': '.kinesis.MetricsName',
      '10': 'shardlevelmetrics'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `EnableEnhancedMonitoringInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableEnhancedMonitoringInputDescriptor = $convert.base64Decode(
    'Ch1FbmFibGVFbmhhbmNlZE1vbml0b3JpbmdJbnB1dBJGChFzaGFyZGxldmVsbWV0cmljcxjt1f'
    'fBASADKA4yFC5raW5lc2lzLk1ldHJpY3NOYW1lUhFzaGFyZGxldmVsbWV0cmljcxIgCglzdHJl'
    'YW1hcm4Y3fOq8gEgASgJUglzdHJlYW1hcm4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW'
    '1pZBIiCgpzdHJlYW1uYW1lGMe3ueABIAEoCVIKc3RyZWFtbmFtZQ==');

@$core.Deprecated('Use enhancedMetricsDescriptor instead')
const EnhancedMetrics$json = {
  '1': 'EnhancedMetrics',
  '2': [
    {
      '1': 'shardlevelmetrics',
      '3': 406711021,
      '4': 3,
      '5': 14,
      '6': '.kinesis.MetricsName',
      '10': 'shardlevelmetrics'
    },
  ],
};

/// Descriptor for `EnhancedMetrics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enhancedMetricsDescriptor = $convert.base64Decode(
    'Cg9FbmhhbmNlZE1ldHJpY3MSRgoRc2hhcmRsZXZlbG1ldHJpY3MY7dX3wQEgAygOMhQua2luZX'
    'Npcy5NZXRyaWNzTmFtZVIRc2hhcmRsZXZlbG1ldHJpY3M=');

@$core.Deprecated('Use enhancedMonitoringOutputDescriptor instead')
const EnhancedMonitoringOutput$json = {
  '1': 'EnhancedMonitoringOutput',
  '2': [
    {
      '1': 'currentshardlevelmetrics',
      '3': 453144018,
      '4': 3,
      '5': 14,
      '6': '.kinesis.MetricsName',
      '10': 'currentshardlevelmetrics'
    },
    {
      '1': 'desiredshardlevelmetrics',
      '3': 231120285,
      '4': 3,
      '5': 14,
      '6': '.kinesis.MetricsName',
      '10': 'desiredshardlevelmetrics'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `EnhancedMonitoringOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enhancedMonitoringOutputDescriptor = $convert.base64Decode(
    'ChhFbmhhbmNlZE1vbml0b3JpbmdPdXRwdXQSVAoYY3VycmVudHNoYXJkbGV2ZWxtZXRyaWNzGN'
    'LbidgBIAMoDjIULmtpbmVzaXMuTWV0cmljc05hbWVSGGN1cnJlbnRzaGFyZGxldmVsbWV0cmlj'
    'cxJTChhkZXNpcmVkc2hhcmRsZXZlbG1ldHJpY3MYnbuabiADKA4yFC5raW5lc2lzLk1ldHJpY3'
    'NOYW1lUhhkZXNpcmVkc2hhcmRsZXZlbG1ldHJpY3MSIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJ'
    'c3RyZWFtYXJuEiIKCnN0cmVhbW5hbWUYx7e54AEgASgJUgpzdHJlYW1uYW1l');

@$core.Deprecated('Use expiredIteratorExceptionDescriptor instead')
const ExpiredIteratorException$json = {
  '1': 'ExpiredIteratorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExpiredIteratorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expiredIteratorExceptionDescriptor =
    $convert.base64Decode(
        'ChhFeHBpcmVkSXRlcmF0b3JFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use expiredNextTokenExceptionDescriptor instead')
const ExpiredNextTokenException$json = {
  '1': 'ExpiredNextTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExpiredNextTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expiredNextTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChlFeHBpcmVkTmV4dFRva2VuRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use getRecordsInputDescriptor instead')
const GetRecordsInput$json = {
  '1': 'GetRecordsInput',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {
      '1': 'sharditerator',
      '3': 379619650,
      '4': 1,
      '5': 9,
      '10': 'sharditerator'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `GetRecordsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRecordsInputDescriptor = $convert.base64Decode(
    'Cg9HZXRSZWNvcmRzSW5wdXQSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBIoCg1zaGFyZGl0ZX'
    'JhdG9yGMKSgrUBIAEoCVINc2hhcmRpdGVyYXRvchIgCglzdHJlYW1hcm4Y3fOq8gEgASgJUglz'
    'dHJlYW1hcm4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZA==');

@$core.Deprecated('Use getRecordsOutputDescriptor instead')
const GetRecordsOutput$json = {
  '1': 'GetRecordsOutput',
  '2': [
    {
      '1': 'childshards',
      '3': 348740657,
      '4': 3,
      '5': 11,
      '6': '.kinesis.ChildShard',
      '10': 'childshards'
    },
    {
      '1': 'millisbehindlatest',
      '3': 456422057,
      '4': 1,
      '5': 3,
      '10': 'millisbehindlatest'
    },
    {
      '1': 'nextsharditerator',
      '3': 442470571,
      '4': 1,
      '5': 9,
      '10': 'nextsharditerator'
    },
    {
      '1': 'records',
      '3': 423557454,
      '4': 3,
      '5': 11,
      '6': '.kinesis.Record',
      '10': 'records'
    },
  ],
};

/// Descriptor for `GetRecordsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getRecordsOutputDescriptor = $convert.base64Decode(
    'ChBHZXRSZWNvcmRzT3V0cHV0EjkKC2NoaWxkc2hhcmRzGLG4paYBIAMoCzITLmtpbmVzaXMuQ2'
    'hpbGRTaGFyZFILY2hpbGRzaGFyZHMSMgoSbWlsbGlzYmVoaW5kbGF0ZXN0GKnl0dkBIAEoA1IS'
    'bWlsbGlzYmVoaW5kbGF0ZXN0EjAKEW5leHRzaGFyZGl0ZXJhdG9yGKuh/tIBIAEoCVIRbmV4dH'
    'NoYXJkaXRlcmF0b3ISLQoHcmVjb3JkcxjO8vvJASADKAsyDy5raW5lc2lzLlJlY29yZFIHcmVj'
    'b3Jkcw==');

@$core.Deprecated('Use getResourcePolicyInputDescriptor instead')
const GetResourcePolicyInput$json = {
  '1': 'GetResourcePolicyInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `GetResourcePolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyInputDescriptor =
    $convert.base64Decode(
        'ChZHZXRSZXNvdXJjZVBvbGljeUlucHV0EiQKC3Jlc291cmNlYXJuGO3AmbABIAEoCVILcmVzb3'
        'VyY2Vhcm4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZA==');

@$core.Deprecated('Use getResourcePolicyOutputDescriptor instead')
const GetResourcePolicyOutput$json = {
  '1': 'GetResourcePolicyOutput',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
  ],
};

/// Descriptor for `GetResourcePolicyOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getResourcePolicyOutputDescriptor =
    $convert.base64Decode(
        'ChdHZXRSZXNvdXJjZVBvbGljeU91dHB1dBIaCgZwb2xpY3kYoO/w4AEgASgJUgZwb2xpY3k=');

@$core.Deprecated('Use getShardIteratorInputDescriptor instead')
const GetShardIteratorInput$json = {
  '1': 'GetShardIteratorInput',
  '2': [
    {'1': 'shardid', '3': 66410951, '4': 1, '5': 9, '10': 'shardid'},
    {
      '1': 'sharditeratortype',
      '3': 229371818,
      '4': 1,
      '5': 14,
      '6': '.kinesis.ShardIteratorType',
      '10': 'sharditeratortype'
    },
    {
      '1': 'startingsequencenumber',
      '3': 88770150,
      '4': 1,
      '5': 9,
      '10': 'startingsequencenumber'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
  ],
};

/// Descriptor for `GetShardIteratorInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getShardIteratorInputDescriptor = $convert.base64Decode(
    'ChVHZXRTaGFyZEl0ZXJhdG9ySW5wdXQSGwoHc2hhcmRpZBjHs9UfIAEoCVIHc2hhcmRpZBJLCh'
    'FzaGFyZGl0ZXJhdG9ydHlwZRiq369tIAEoDjIaLmtpbmVzaXMuU2hhcmRJdGVyYXRvclR5cGVS'
    'EXNoYXJkaXRlcmF0b3J0eXBlEjkKFnN0YXJ0aW5nc2VxdWVuY2VudW1iZXIY5oyqKiABKAlSFn'
    'N0YXJ0aW5nc2VxdWVuY2VudW1iZXISIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJu'
    'Eh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSIgoKc3RyZWFtbmFtZRjHt7ngASABKA'
    'lSCnN0cmVhbW5hbWUSHwoJdGltZXN0YW1wGMTDt00gASgJUgl0aW1lc3RhbXA=');

@$core.Deprecated('Use getShardIteratorOutputDescriptor instead')
const GetShardIteratorOutput$json = {
  '1': 'GetShardIteratorOutput',
  '2': [
    {
      '1': 'sharditerator',
      '3': 379619650,
      '4': 1,
      '5': 9,
      '10': 'sharditerator'
    },
  ],
};

/// Descriptor for `GetShardIteratorOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getShardIteratorOutputDescriptor =
    $convert.base64Decode(
        'ChZHZXRTaGFyZEl0ZXJhdG9yT3V0cHV0EigKDXNoYXJkaXRlcmF0b3IYwpKCtQEgASgJUg1zaG'
        'FyZGl0ZXJhdG9y');

@$core.Deprecated('Use hashKeyRangeDescriptor instead')
const HashKeyRange$json = {
  '1': 'HashKeyRange',
  '2': [
    {
      '1': 'endinghashkey',
      '3': 399409906,
      '4': 1,
      '5': 9,
      '10': 'endinghashkey'
    },
    {
      '1': 'startinghashkey',
      '3': 202551085,
      '4': 1,
      '5': 9,
      '10': 'startinghashkey'
    },
  ],
};

/// Descriptor for `HashKeyRange`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List hashKeyRangeDescriptor = $convert.base64Decode(
    'CgxIYXNoS2V5UmFuZ2USKAoNZW5kaW5naGFzaGtleRjyhbq+ASABKAlSDWVuZGluZ2hhc2hrZX'
    'kSKwoPc3RhcnRpbmdoYXNoa2V5GK3eymAgASgJUg9zdGFydGluZ2hhc2hrZXk=');

@$core.Deprecated('Use increaseStreamRetentionPeriodInputDescriptor instead')
const IncreaseStreamRetentionPeriodInput$json = {
  '1': 'IncreaseStreamRetentionPeriodInput',
  '2': [
    {
      '1': 'retentionperiodhours',
      '3': 396381944,
      '4': 1,
      '5': 5,
      '10': 'retentionperiodhours'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `IncreaseStreamRetentionPeriodInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List increaseStreamRetentionPeriodInputDescriptor =
    $convert.base64Decode(
        'CiJJbmNyZWFzZVN0cmVhbVJldGVudGlvblBlcmlvZElucHV0EjYKFHJldGVudGlvbnBlcmlvZG'
        'hvdXJzGPidgb0BIAEoBVIUcmV0ZW50aW9ucGVyaW9kaG91cnMSIAoJc3RyZWFtYXJuGN3zqvIB'
        'IAEoCVIJc3RyZWFtYXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSIgoKc3RyZW'
        'FtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWU=');

@$core.Deprecated('Use internalFailureExceptionDescriptor instead')
const InternalFailureException$json = {
  '1': 'InternalFailureException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalFailureException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalFailureExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnRlcm5hbEZhaWx1cmVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidArgumentExceptionDescriptor instead')
const InvalidArgumentException$json = {
  '1': 'InvalidArgumentException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidArgumentException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidArgumentExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkQXJndW1lbnRFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

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

@$core.Deprecated('Use listShardsInputDescriptor instead')
const ListShardsInput$json = {
  '1': 'ListShardsInput',
  '2': [
    {
      '1': 'exclusivestartshardid',
      '3': 44771587,
      '4': 1,
      '5': 9,
      '10': 'exclusivestartshardid'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'shardfilter',
      '3': 230254710,
      '4': 1,
      '5': 11,
      '6': '.kinesis.ShardFilter',
      '10': 'shardfilter'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {
      '1': 'streamcreationtimestamp',
      '3': 224951013,
      '4': 1,
      '5': 9,
      '10': 'streamcreationtimestamp'
    },
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `ListShardsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listShardsInputDescriptor = $convert.base64Decode(
    'Cg9MaXN0U2hhcmRzSW5wdXQSNwoVZXhjbHVzaXZlc3RhcnRzaGFyZGlkGIPSrBUgASgJUhVleG'
    'NsdXNpdmVzdGFydHNoYXJkaWQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heHJlc3VsdHMS'
    'HwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SOQoLc2hhcmRmaWx0ZXIY9tDlbSABKA'
    'syFC5raW5lc2lzLlNoYXJkRmlsdGVyUgtzaGFyZGZpbHRlchIgCglzdHJlYW1hcm4Y3fOq8gEg'
    'ASgJUglzdHJlYW1hcm4SOwoXc3RyZWFtY3JlYXRpb250aW1lc3RhbXAY5fWhayABKAlSF3N0cm'
    'VhbWNyZWF0aW9udGltZXN0YW1wEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSIgoK'
    'c3RyZWFtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWU=');

@$core.Deprecated('Use listShardsOutputDescriptor instead')
const ListShardsOutput$json = {
  '1': 'ListShardsOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'shards',
      '3': 437117641,
      '4': 3,
      '5': 11,
      '6': '.kinesis.Shard',
      '10': 'shards'
    },
  ],
};

/// Descriptor for `ListShardsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listShardsOutputDescriptor = $convert.base64Decode(
    'ChBMaXN0U2hhcmRzT3V0cHV0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEioKBn'
    'NoYXJkcxjJxbfQASADKAsyDi5raW5lc2lzLlNoYXJkUgZzaGFyZHM=');

@$core.Deprecated('Use listStreamConsumersInputDescriptor instead')
const ListStreamConsumersInput$json = {
  '1': 'ListStreamConsumersInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {
      '1': 'streamcreationtimestamp',
      '3': 224951013,
      '4': 1,
      '5': 9,
      '10': 'streamcreationtimestamp'
    },
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `ListStreamConsumersInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStreamConsumersInputDescriptor = $convert.base64Decode(
    'ChhMaXN0U3RyZWFtQ29uc3VtZXJzSW5wdXQSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm1heH'
    'Jlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SIAoJc3RyZWFtYXJuGN3z'
    'qvIBIAEoCVIJc3RyZWFtYXJuEjsKF3N0cmVhbWNyZWF0aW9udGltZXN0YW1wGOX1oWsgASgJUh'
    'dzdHJlYW1jcmVhdGlvbnRpbWVzdGFtcBIeCghzdHJlYW1pZBjB7YHGASABKAlSCHN0cmVhbWlk');

@$core.Deprecated('Use listStreamConsumersOutputDescriptor instead')
const ListStreamConsumersOutput$json = {
  '1': 'ListStreamConsumersOutput',
  '2': [
    {
      '1': 'consumers',
      '3': 348153455,
      '4': 3,
      '5': 11,
      '6': '.kinesis.Consumer',
      '10': 'consumers'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListStreamConsumersOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStreamConsumersOutputDescriptor = $convert.base64Decode(
    'ChlMaXN0U3RyZWFtQ29uc3VtZXJzT3V0cHV0EjMKCWNvbnN1bWVycxjvzIGmASADKAsyES5raW'
    '5lc2lzLkNvbnN1bWVyUgljb25zdW1lcnMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9r'
    'ZW4=');

@$core.Deprecated('Use listStreamsInputDescriptor instead')
const ListStreamsInput$json = {
  '1': 'ListStreamsInput',
  '2': [
    {
      '1': 'exclusivestartstreamname',
      '3': 431108955,
      '4': 1,
      '5': 9,
      '10': 'exclusivestartstreamname'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListStreamsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStreamsInputDescriptor = $convert.base64Decode(
    'ChBMaXN0U3RyZWFtc0lucHV0Ej4KGGV4Y2x1c2l2ZXN0YXJ0c3RyZWFtbmFtZRjb5sjNASABKA'
    'lSGGV4Y2x1c2l2ZXN0YXJ0c3RyZWFtbmFtZRIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0Eh8K'
    'CW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listStreamsOutputDescriptor instead')
const ListStreamsOutput$json = {
  '1': 'ListStreamsOutput',
  '2': [
    {
      '1': 'hasmorestreams',
      '3': 181206864,
      '4': 1,
      '5': 8,
      '10': 'hasmorestreams'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'streamnames', '3': 530210288, '4': 3, '5': 9, '10': 'streamnames'},
    {
      '1': 'streamsummaries',
      '3': 334316244,
      '4': 3,
      '5': 11,
      '6': '.kinesis.StreamSummary',
      '10': 'streamsummaries'
    },
  ],
};

/// Descriptor for `ListStreamsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listStreamsOutputDescriptor = $convert.base64Decode(
    'ChFMaXN0U3RyZWFtc091dHB1dBIpCg5oYXNtb3Jlc3RyZWFtcxjQ/rNWIAEoCFIOaGFzbW9yZX'
    'N0cmVhbXMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SJAoLc3RyZWFtbmFtZXMY'
    '8Lvp/AEgAygJUgtzdHJlYW1uYW1lcxJECg9zdHJlYW1zdW1tYXJpZXMY1IW1nwEgAygLMhYua2'
    'luZXNpcy5TdHJlYW1TdW1tYXJ5Ug9zdHJlYW1zdW1tYXJpZXM=');

@$core.Deprecated('Use listTagsForResourceInputDescriptor instead')
const ListTagsForResourceInput$json = {
  '1': 'ListTagsForResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `ListTagsForResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceInputDescriptor =
    $convert.base64Decode(
        'ChhMaXN0VGFnc0ZvclJlc291cmNlSW5wdXQSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUgtyZX'
        'NvdXJjZWFybhIeCghzdHJlYW1pZBjB7YHGASABKAlSCHN0cmVhbWlk');

@$core.Deprecated('Use listTagsForResourceOutputDescriptor instead')
const ListTagsForResourceOutput$json = {
  '1': 'ListTagsForResourceOutput',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kinesis.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceOutputDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFnc0ZvclJlc291cmNlT3V0cHV0EiQKBHRhZ3MYwcH2tQEgAygLMgwua2luZXNpcy'
        '5UYWdSBHRhZ3M=');

@$core.Deprecated('Use listTagsForStreamInputDescriptor instead')
const ListTagsForStreamInput$json = {
  '1': 'ListTagsForStreamInput',
  '2': [
    {
      '1': 'exclusivestarttagkey',
      '3': 470465287,
      '4': 1,
      '5': 9,
      '10': 'exclusivestarttagkey'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `ListTagsForStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForStreamInputDescriptor = $convert.base64Decode(
    'ChZMaXN0VGFnc0ZvclN0cmVhbUlucHV0EjYKFGV4Y2x1c2l2ZXN0YXJ0dGFna2V5GIf2quABIA'
    'EoCVIUZXhjbHVzaXZlc3RhcnR0YWdrZXkSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBIgCglz'
    'dHJlYW1hcm4Y3fOq8gEgASgJUglzdHJlYW1hcm4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdH'
    'JlYW1pZBIiCgpzdHJlYW1uYW1lGMe3ueABIAEoCVIKc3RyZWFtbmFtZQ==');

@$core.Deprecated('Use listTagsForStreamOutputDescriptor instead')
const ListTagsForStreamOutput$json = {
  '1': 'ListTagsForStreamOutput',
  '2': [
    {'1': 'hasmoretags', '3': 397963248, '4': 1, '5': 8, '10': 'hasmoretags'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kinesis.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForStreamOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForStreamOutputDescriptor =
    $convert.base64Decode(
        'ChdMaXN0VGFnc0ZvclN0cmVhbU91dHB1dBIkCgtoYXNtb3JldGFncxjw3+G9ASABKAhSC2hhc2'
        '1vcmV0YWdzEiQKBHRhZ3MYwcH2tQEgAygLMgwua2luZXNpcy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use mergeShardsInputDescriptor instead')
const MergeShardsInput$json = {
  '1': 'MergeShardsInput',
  '2': [
    {
      '1': 'adjacentshardtomerge',
      '3': 435883605,
      '4': 1,
      '5': 9,
      '10': 'adjacentshardtomerge'
    },
    {'1': 'shardtomerge', '3': 473343671, '4': 1, '5': 9, '10': 'shardtomerge'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `MergeShardsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mergeShardsInputDescriptor = $convert.base64Decode(
    'ChBNZXJnZVNoYXJkc0lucHV0EjYKFGFkamFjZW50c2hhcmR0b21lcmdlGNWc7M8BIAEoCVIUYW'
    'RqYWNlbnRzaGFyZHRvbWVyZ2USJgoMc2hhcmR0b21lcmdlGLfN2uEBIAEoCVIMc2hhcmR0b21l'
    'cmdlEiAKCXN0cmVhbWFybhjd86ryASABKAlSCXN0cmVhbWFybhIeCghzdHJlYW1pZBjB7YHGAS'
    'ABKAlSCHN0cmVhbWlkEiIKCnN0cmVhbW5hbWUYx7e54AEgASgJUgpzdHJlYW1uYW1l');

@$core
    .Deprecated('Use minimumThroughputBillingCommitmentInputDescriptor instead')
const MinimumThroughputBillingCommitmentInput$json = {
  '1': 'MinimumThroughputBillingCommitmentInput',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.kinesis.MinimumThroughputBillingCommitmentInputStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `MinimumThroughputBillingCommitmentInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List minimumThroughputBillingCommitmentInputDescriptor =
    $convert.base64Decode(
        'CidNaW5pbXVtVGhyb3VnaHB1dEJpbGxpbmdDb21taXRtZW50SW5wdXQSUQoGc3RhdHVzGJDk+w'
        'IgASgOMjYua2luZXNpcy5NaW5pbXVtVGhyb3VnaHB1dEJpbGxpbmdDb21taXRtZW50SW5wdXRT'
        'dGF0dXNSBnN0YXR1cw==');

@$core.Deprecated(
    'Use minimumThroughputBillingCommitmentOutputDescriptor instead')
const MinimumThroughputBillingCommitmentOutput$json = {
  '1': 'MinimumThroughputBillingCommitmentOutput',
  '2': [
    {
      '1': 'earliestallowedendat',
      '3': 77927885,
      '4': 1,
      '5': 9,
      '10': 'earliestallowedendat'
    },
    {'1': 'endedat', '3': 104122351, '4': 1, '5': 9, '10': 'endedat'},
    {'1': 'startedat', '3': 77629404, '4': 1, '5': 9, '10': 'startedat'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.kinesis.MinimumThroughputBillingCommitmentOutputStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `MinimumThroughputBillingCommitmentOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List minimumThroughputBillingCommitmentOutputDescriptor =
    $convert.base64Decode(
        'CihNaW5pbXVtVGhyb3VnaHB1dEJpbGxpbmdDb21taXRtZW50T3V0cHV0EjUKFGVhcmxpZXN0YW'
        'xsb3dlZGVuZGF0GM2rlCUgASgJUhRlYXJsaWVzdGFsbG93ZWRlbmRhdBIbCgdlbmRlZGF0GO+P'
        '0zEgASgJUgdlbmRlZGF0Eh8KCXN0YXJ0ZWRhdBjcj4IlIAEoCVIJc3RhcnRlZGF0ElIKBnN0YX'
        'R1cxiQ5PsCIAEoDjI3LmtpbmVzaXMuTWluaW11bVRocm91Z2hwdXRCaWxsaW5nQ29tbWl0bWVu'
        'dE91dHB1dFN0YXR1c1IGc3RhdHVz');

@$core
    .Deprecated('Use provisionedThroughputExceededExceptionDescriptor instead')
const ProvisionedThroughputExceededException$json = {
  '1': 'ProvisionedThroughputExceededException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ProvisionedThroughputExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List provisionedThroughputExceededExceptionDescriptor =
    $convert.base64Decode(
        'CiZQcm92aXNpb25lZFRocm91Z2hwdXRFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyC'
        'cgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use putRecordInputDescriptor instead')
const PutRecordInput$json = {
  '1': 'PutRecordInput',
  '2': [
    {'1': 'data', '3': 525498822, '4': 1, '5': 12, '10': 'data'},
    {
      '1': 'explicithashkey',
      '3': 438141821,
      '4': 1,
      '5': 9,
      '10': 'explicithashkey'
    },
    {'1': 'partitionkey', '3': 379379617, '4': 1, '5': 9, '10': 'partitionkey'},
    {
      '1': 'sequencenumberforordering',
      '3': 73727749,
      '4': 1,
      '5': 9,
      '10': 'sequencenumberforordering'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `PutRecordInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRecordInputDescriptor = $convert.base64Decode(
    'Cg5QdXRSZWNvcmRJbnB1dBIWCgRkYXRhGMbzyfoBIAEoDFIEZGF0YRIsCg9leHBsaWNpdGhhc2'
    'hrZXkY/Yb20AEgASgJUg9leHBsaWNpdGhhc2hrZXkSJgoMcGFydGl0aW9ua2V5GKG/87QBIAEo'
    'CVIMcGFydGl0aW9ua2V5Ej8KGXNlcXVlbmNlbnVtYmVyZm9yb3JkZXJpbmcYhf6TIyABKAlSGX'
    'NlcXVlbmNlbnVtYmVyZm9yb3JkZXJpbmcSIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFt'
    'YXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSIgoKc3RyZWFtbmFtZRjHt7ngAS'
    'ABKAlSCnN0cmVhbW5hbWU=');

@$core.Deprecated('Use putRecordOutputDescriptor instead')
const PutRecordOutput$json = {
  '1': 'PutRecordOutput',
  '2': [
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.kinesis.EncryptionType',
      '10': 'encryptiontype'
    },
    {
      '1': 'sequencenumber',
      '3': 98094362,
      '4': 1,
      '5': 9,
      '10': 'sequencenumber'
    },
    {'1': 'shardid', '3': 66410951, '4': 1, '5': 9, '10': 'shardid'},
  ],
};

/// Descriptor for `PutRecordOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRecordOutputDescriptor = $convert.base64Decode(
    'Cg9QdXRSZWNvcmRPdXRwdXQSQgoOZW5jcnlwdGlvbnR5cGUYtd/xfSABKA4yFy5raW5lc2lzLk'
    'VuY3J5cHRpb25UeXBlUg5lbmNyeXB0aW9udHlwZRIpCg5zZXF1ZW5jZW51bWJlchiamuMuIAEo'
    'CVIOc2VxdWVuY2VudW1iZXISGwoHc2hhcmRpZBjHs9UfIAEoCVIHc2hhcmRpZA==');

@$core.Deprecated('Use putRecordsInputDescriptor instead')
const PutRecordsInput$json = {
  '1': 'PutRecordsInput',
  '2': [
    {
      '1': 'records',
      '3': 423557454,
      '4': 3,
      '5': 11,
      '6': '.kinesis.PutRecordsRequestEntry',
      '10': 'records'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `PutRecordsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRecordsInputDescriptor = $convert.base64Decode(
    'Cg9QdXRSZWNvcmRzSW5wdXQSPQoHcmVjb3JkcxjO8vvJASADKAsyHy5raW5lc2lzLlB1dFJlY2'
    '9yZHNSZXF1ZXN0RW50cnlSB3JlY29yZHMSIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFt'
    'YXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSIgoKc3RyZWFtbmFtZRjHt7ngAS'
    'ABKAlSCnN0cmVhbW5hbWU=');

@$core.Deprecated('Use putRecordsOutputDescriptor instead')
const PutRecordsOutput$json = {
  '1': 'PutRecordsOutput',
  '2': [
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.kinesis.EncryptionType',
      '10': 'encryptiontype'
    },
    {
      '1': 'failedrecordcount',
      '3': 89270457,
      '4': 1,
      '5': 5,
      '10': 'failedrecordcount'
    },
    {
      '1': 'records',
      '3': 423557454,
      '4': 3,
      '5': 11,
      '6': '.kinesis.PutRecordsResultEntry',
      '10': 'records'
    },
  ],
};

/// Descriptor for `PutRecordsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRecordsOutputDescriptor = $convert.base64Decode(
    'ChBQdXRSZWNvcmRzT3V0cHV0EkIKDmVuY3J5cHRpb250eXBlGLXf8X0gASgOMhcua2luZXNpcy'
    '5FbmNyeXB0aW9uVHlwZVIOZW5jcnlwdGlvbnR5cGUSLwoRZmFpbGVkcmVjb3JkY291bnQYudHI'
    'KiABKAVSEWZhaWxlZHJlY29yZGNvdW50EjwKB3JlY29yZHMYzvL7yQEgAygLMh4ua2luZXNpcy'
    '5QdXRSZWNvcmRzUmVzdWx0RW50cnlSB3JlY29yZHM=');

@$core.Deprecated('Use putRecordsRequestEntryDescriptor instead')
const PutRecordsRequestEntry$json = {
  '1': 'PutRecordsRequestEntry',
  '2': [
    {'1': 'data', '3': 525498822, '4': 1, '5': 12, '10': 'data'},
    {
      '1': 'explicithashkey',
      '3': 438141821,
      '4': 1,
      '5': 9,
      '10': 'explicithashkey'
    },
    {'1': 'partitionkey', '3': 379379617, '4': 1, '5': 9, '10': 'partitionkey'},
  ],
};

/// Descriptor for `PutRecordsRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRecordsRequestEntryDescriptor = $convert.base64Decode(
    'ChZQdXRSZWNvcmRzUmVxdWVzdEVudHJ5EhYKBGRhdGEYxvPJ+gEgASgMUgRkYXRhEiwKD2V4cG'
    'xpY2l0aGFzaGtleRj9hvbQASABKAlSD2V4cGxpY2l0aGFzaGtleRImCgxwYXJ0aXRpb25rZXkY'
    'ob/ztAEgASgJUgxwYXJ0aXRpb25rZXk=');

@$core.Deprecated('Use putRecordsResultEntryDescriptor instead')
const PutRecordsResultEntry$json = {
  '1': 'PutRecordsResultEntry',
  '2': [
    {'1': 'errorcode', '3': 34663193, '4': 1, '5': 9, '10': 'errorcode'},
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'sequencenumber',
      '3': 98094362,
      '4': 1,
      '5': 9,
      '10': 'sequencenumber'
    },
    {'1': 'shardid', '3': 66410951, '4': 1, '5': 9, '10': 'shardid'},
  ],
};

/// Descriptor for `PutRecordsResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putRecordsResultEntryDescriptor = $convert.base64Decode(
    'ChVQdXRSZWNvcmRzUmVzdWx0RW50cnkSHwoJZXJyb3Jjb2RlGJnWwxAgASgJUgllcnJvcmNvZG'
    'USJgoMZXJyb3JtZXNzYWdlGKmKq/cBIAEoCVIMZXJyb3JtZXNzYWdlEikKDnNlcXVlbmNlbnVt'
    'YmVyGJqa4y4gASgJUg5zZXF1ZW5jZW51bWJlchIbCgdzaGFyZGlkGMez1R8gASgJUgdzaGFyZG'
    'lk');

@$core.Deprecated('Use putResourcePolicyInputDescriptor instead')
const PutResourcePolicyInput$json = {
  '1': 'PutResourcePolicyInput',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `PutResourcePolicyInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putResourcePolicyInputDescriptor = $convert.base64Decode(
    'ChZQdXRSZXNvdXJjZVBvbGljeUlucHV0EhoKBnBvbGljeRig7/DgASABKAlSBnBvbGljeRIkCg'
    'tyZXNvdXJjZWFybhjtwJmwASABKAlSC3Jlc291cmNlYXJuEh4KCHN0cmVhbWlkGMHtgcYBIAEo'
    'CVIIc3RyZWFtaWQ=');

@$core.Deprecated('Use recordDescriptor instead')
const Record$json = {
  '1': 'Record',
  '2': [
    {
      '1': 'approximatearrivaltimestamp',
      '3': 95039887,
      '4': 1,
      '5': 9,
      '10': 'approximatearrivaltimestamp'
    },
    {'1': 'data', '3': 525498822, '4': 1, '5': 12, '10': 'data'},
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.kinesis.EncryptionType',
      '10': 'encryptiontype'
    },
    {'1': 'partitionkey', '3': 379379617, '4': 1, '5': 9, '10': 'partitionkey'},
    {
      '1': 'sequencenumber',
      '3': 98094362,
      '4': 1,
      '5': 9,
      '10': 'sequencenumber'
    },
  ],
};

/// Descriptor for `Record`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recordDescriptor = $convert.base64Decode(
    'CgZSZWNvcmQSQwobYXBwcm94aW1hdGVhcnJpdmFsdGltZXN0YW1wGI/jqC0gASgJUhthcHByb3'
    'hpbWF0ZWFycml2YWx0aW1lc3RhbXASFgoEZGF0YRjG88n6ASABKAxSBGRhdGESQgoOZW5jcnlw'
    'dGlvbnR5cGUYtd/xfSABKA4yFy5raW5lc2lzLkVuY3J5cHRpb25UeXBlUg5lbmNyeXB0aW9udH'
    'lwZRImCgxwYXJ0aXRpb25rZXkYob/ztAEgASgJUgxwYXJ0aXRpb25rZXkSKQoOc2VxdWVuY2Vu'
    'dW1iZXIYmprjLiABKAlSDnNlcXVlbmNlbnVtYmVy');

@$core.Deprecated('Use registerStreamConsumerInputDescriptor instead')
const RegisterStreamConsumerInput$json = {
  '1': 'RegisterStreamConsumerInput',
  '2': [
    {'1': 'consumername', '3': 70979235, '4': 1, '5': 9, '10': 'consumername'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kinesis.RegisterStreamConsumerInput.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [RegisterStreamConsumerInput_TagsEntry$json],
};

@$core.Deprecated('Use registerStreamConsumerInputDescriptor instead')
const RegisterStreamConsumerInput_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `RegisterStreamConsumerInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerStreamConsumerInputDescriptor = $convert.base64Decode(
    'ChtSZWdpc3RlclN0cmVhbUNvbnN1bWVySW5wdXQSJQoMY29uc3VtZXJuYW1lGKOd7CEgASgJUg'
    'xjb25zdW1lcm5hbWUSIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJuEh4KCHN0cmVh'
    'bWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSRgoEdGFncxjBwfa1ASADKAsyLi5raW5lc2lzLlJlZ2'
    'lzdGVyU3RyZWFtQ29uc3VtZXJJbnB1dC5UYWdzRW50cnlSBHRhZ3MaNwoJVGFnc0VudHJ5EhAK'
    'A2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use registerStreamConsumerOutputDescriptor instead')
const RegisterStreamConsumerOutput$json = {
  '1': 'RegisterStreamConsumerOutput',
  '2': [
    {
      '1': 'consumer',
      '3': 482032874,
      '4': 1,
      '5': 11,
      '6': '.kinesis.Consumer',
      '10': 'consumer'
    },
  ],
};

/// Descriptor for `RegisterStreamConsumerOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List registerStreamConsumerOutputDescriptor =
    $convert.base64Decode(
        'ChxSZWdpc3RlclN0cmVhbUNvbnN1bWVyT3V0cHV0EjEKCGNvbnN1bWVyGOr57OUBIAEoCzIRLm'
        'tpbmVzaXMuQ29uc3VtZXJSCGNvbnN1bWVy');

@$core.Deprecated('Use removeTagsFromStreamInputDescriptor instead')
const RemoveTagsFromStreamInput$json = {
  '1': 'RemoveTagsFromStreamInput',
  '2': [
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `RemoveTagsFromStreamInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removeTagsFromStreamInputDescriptor = $convert.base64Decode(
    'ChlSZW1vdmVUYWdzRnJvbVN0cmVhbUlucHV0EiAKCXN0cmVhbWFybhjd86ryASABKAlSCXN0cm'
    'VhbWFybhIeCghzdHJlYW1pZBjB7YHGASABKAlSCHN0cmVhbWlkEiIKCnN0cmVhbW5hbWUYx7e5'
    '4AEgASgJUgpzdHJlYW1uYW1lEhwKB3RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use resourceInUseExceptionDescriptor instead')
const ResourceInUseException$json = {
  '1': 'ResourceInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChZSZXNvdXJjZUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

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

@$core.Deprecated('Use sequenceNumberRangeDescriptor instead')
const SequenceNumberRange$json = {
  '1': 'SequenceNumberRange',
  '2': [
    {
      '1': 'endingsequencenumber',
      '3': 328882583,
      '4': 1,
      '5': 9,
      '10': 'endingsequencenumber'
    },
    {
      '1': 'startingsequencenumber',
      '3': 88770150,
      '4': 1,
      '5': 9,
      '10': 'startingsequencenumber'
    },
  ],
};

/// Descriptor for `SequenceNumberRange`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sequenceNumberRangeDescriptor = $convert.base64Decode(
    'ChNTZXF1ZW5jZU51bWJlclJhbmdlEjYKFGVuZGluZ3NlcXVlbmNlbnVtYmVyGJez6ZwBIAEoCV'
    'IUZW5kaW5nc2VxdWVuY2VudW1iZXISOQoWc3RhcnRpbmdzZXF1ZW5jZW51bWJlchjmjKoqIAEo'
    'CVIWc3RhcnRpbmdzZXF1ZW5jZW51bWJlcg==');

@$core.Deprecated('Use shardDescriptor instead')
const Shard$json = {
  '1': 'Shard',
  '2': [
    {
      '1': 'adjacentparentshardid',
      '3': 310461245,
      '4': 1,
      '5': 9,
      '10': 'adjacentparentshardid'
    },
    {
      '1': 'hashkeyrange',
      '3': 981486,
      '4': 1,
      '5': 11,
      '6': '.kinesis.HashKeyRange',
      '10': 'hashkeyrange'
    },
    {
      '1': 'parentshardid',
      '3': 431117103,
      '4': 1,
      '5': 9,
      '10': 'parentshardid'
    },
    {
      '1': 'sequencenumberrange',
      '3': 86578567,
      '4': 1,
      '5': 11,
      '6': '.kinesis.SequenceNumberRange',
      '10': 'sequencenumberrange'
    },
    {'1': 'shardid', '3': 66410951, '4': 1, '5': 9, '10': 'shardid'},
  ],
};

/// Descriptor for `Shard`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List shardDescriptor = $convert.base64Decode(
    'CgVTaGFyZBI4ChVhZGphY2VudHBhcmVudHNoYXJkaWQYvYaFlAEgASgJUhVhZGphY2VudHBhcm'
    'VudHNoYXJkaWQSOwoMaGFzaGtleXJhbmdlGO7zOyABKAsyFS5raW5lc2lzLkhhc2hLZXlSYW5n'
    'ZVIMaGFzaGtleXJhbmdlEigKDXBhcmVudHNoYXJkaWQYr6bJzQEgASgJUg1wYXJlbnRzaGFyZG'
    'lkElEKE3NlcXVlbmNlbnVtYmVycmFuZ2UYh6ukKSABKAsyHC5raW5lc2lzLlNlcXVlbmNlTnVt'
    'YmVyUmFuZ2VSE3NlcXVlbmNlbnVtYmVycmFuZ2USGwoHc2hhcmRpZBjHs9UfIAEoCVIHc2hhcm'
    'RpZA==');

@$core.Deprecated('Use shardFilterDescriptor instead')
const ShardFilter$json = {
  '1': 'ShardFilter',
  '2': [
    {'1': 'shardid', '3': 66410951, '4': 1, '5': 9, '10': 'shardid'},
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.kinesis.ShardFilterType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `ShardFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List shardFilterDescriptor = $convert.base64Decode(
    'CgtTaGFyZEZpbHRlchIbCgdzaGFyZGlkGMez1R8gASgJUgdzaGFyZGlkEh8KCXRpbWVzdGFtcB'
    'jEw7dNIAEoCVIJdGltZXN0YW1wEjAKBHR5cGUY7qDXigEgASgOMhgua2luZXNpcy5TaGFyZEZp'
    'bHRlclR5cGVSBHR5cGU=');

@$core.Deprecated('Use splitShardInputDescriptor instead')
const SplitShardInput$json = {
  '1': 'SplitShardInput',
  '2': [
    {
      '1': 'newstartinghashkey',
      '3': 351052649,
      '4': 1,
      '5': 9,
      '10': 'newstartinghashkey'
    },
    {'1': 'shardtosplit', '3': 124434671, '4': 1, '5': 9, '10': 'shardtosplit'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `SplitShardInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List splitShardInputDescriptor = $convert.base64Decode(
    'Cg9TcGxpdFNoYXJkSW5wdXQSMgoSbmV3c3RhcnRpbmdoYXNoa2V5GOnGsqcBIAEoCVISbmV3c3'
    'RhcnRpbmdoYXNoa2V5EiUKDHNoYXJkdG9zcGxpdBjv8ao7IAEoCVIMc2hhcmR0b3NwbGl0EiAK'
    'CXN0cmVhbWFybhjd86ryASABKAlSCXN0cmVhbWFybhIeCghzdHJlYW1pZBjB7YHGASABKAlSCH'
    'N0cmVhbWlkEiIKCnN0cmVhbW5hbWUYx7e54AEgASgJUgpzdHJlYW1uYW1l');

@$core.Deprecated('Use startStreamEncryptionInputDescriptor instead')
const StartStreamEncryptionInput$json = {
  '1': 'StartStreamEncryptionInput',
  '2': [
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.kinesis.EncryptionType',
      '10': 'encryptiontype'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `StartStreamEncryptionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startStreamEncryptionInputDescriptor = $convert.base64Decode(
    'ChpTdGFydFN0cmVhbUVuY3J5cHRpb25JbnB1dBJCCg5lbmNyeXB0aW9udHlwZRi13/F9IAEoDj'
    'IXLmtpbmVzaXMuRW5jcnlwdGlvblR5cGVSDmVuY3J5cHRpb250eXBlEhgKBWtleWlkGKKAyIMB'
    'IAEoCVIFa2V5aWQSIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJuEh4KCHN0cmVhbW'
    'lkGMHtgcYBIAEoCVIIc3RyZWFtaWQSIgoKc3RyZWFtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5h'
    'bWU=');

@$core.Deprecated('Use startingPositionDescriptor instead')
const StartingPosition$json = {
  '1': 'StartingPosition',
  '2': [
    {
      '1': 'sequencenumber',
      '3': 98094362,
      '4': 1,
      '5': 9,
      '10': 'sequencenumber'
    },
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.kinesis.ShardIteratorType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `StartingPosition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startingPositionDescriptor = $convert.base64Decode(
    'ChBTdGFydGluZ1Bvc2l0aW9uEikKDnNlcXVlbmNlbnVtYmVyGJqa4y4gASgJUg5zZXF1ZW5jZW'
    '51bWJlchIfCgl0aW1lc3RhbXAYxMO3TSABKAlSCXRpbWVzdGFtcBIyCgR0eXBlGO6g14oBIAEo'
    'DjIaLmtpbmVzaXMuU2hhcmRJdGVyYXRvclR5cGVSBHR5cGU=');

@$core.Deprecated('Use stopStreamEncryptionInputDescriptor instead')
const StopStreamEncryptionInput$json = {
  '1': 'StopStreamEncryptionInput',
  '2': [
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.kinesis.EncryptionType',
      '10': 'encryptiontype'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
  ],
};

/// Descriptor for `StopStreamEncryptionInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List stopStreamEncryptionInputDescriptor = $convert.base64Decode(
    'ChlTdG9wU3RyZWFtRW5jcnlwdGlvbklucHV0EkIKDmVuY3J5cHRpb250eXBlGLXf8X0gASgOMh'
    'cua2luZXNpcy5FbmNyeXB0aW9uVHlwZVIOZW5jcnlwdGlvbnR5cGUSGAoFa2V5aWQYooDIgwEg'
    'ASgJUgVrZXlpZBIgCglzdHJlYW1hcm4Y3fOq8gEgASgJUglzdHJlYW1hcm4SHgoIc3RyZWFtaW'
    'QYwe2BxgEgASgJUghzdHJlYW1pZBIiCgpzdHJlYW1uYW1lGMe3ueABIAEoCVIKc3RyZWFtbmFt'
    'ZQ==');

@$core.Deprecated('Use streamDescriptionDescriptor instead')
const StreamDescription$json = {
  '1': 'StreamDescription',
  '2': [
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.kinesis.EncryptionType',
      '10': 'encryptiontype'
    },
    {
      '1': 'enhancedmonitoring',
      '3': 452259826,
      '4': 3,
      '5': 11,
      '6': '.kinesis.EnhancedMetrics',
      '10': 'enhancedmonitoring'
    },
    {
      '1': 'hasmoreshards',
      '3': 10836604,
      '4': 1,
      '5': 8,
      '10': 'hasmoreshards'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'retentionperiodhours',
      '3': 396381944,
      '4': 1,
      '5': 5,
      '10': 'retentionperiodhours'
    },
    {
      '1': 'shards',
      '3': 437117641,
      '4': 3,
      '5': 11,
      '6': '.kinesis.Shard',
      '10': 'shards'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {
      '1': 'streamcreationtimestamp',
      '3': 224951013,
      '4': 1,
      '5': 9,
      '10': 'streamcreationtimestamp'
    },
    {
      '1': 'streammodedetails',
      '3': 12139665,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StreamModeDetails',
      '10': 'streammodedetails'
    },
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'streamstatus',
      '3': 245792976,
      '4': 1,
      '5': 14,
      '6': '.kinesis.StreamStatus',
      '10': 'streamstatus'
    },
  ],
};

/// Descriptor for `StreamDescription`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamDescriptionDescriptor = $convert.base64Decode(
    'ChFTdHJlYW1EZXNjcmlwdGlvbhJCCg5lbmNyeXB0aW9udHlwZRi13/F9IAEoDjIXLmtpbmVzaX'
    'MuRW5jcnlwdGlvblR5cGVSDmVuY3J5cHRpb250eXBlEkwKEmVuaGFuY2VkbW9uaXRvcmluZxjy'
    '39PXASADKAsyGC5raW5lc2lzLkVuaGFuY2VkTWV0cmljc1ISZW5oYW5jZWRtb25pdG9yaW5nEi'
    'cKDWhhc21vcmVzaGFyZHMY/LSVBSABKAhSDWhhc21vcmVzaGFyZHMSGAoFa2V5aWQYooDIgwEg'
    'ASgJUgVrZXlpZBI2ChRyZXRlbnRpb25wZXJpb2Rob3Vycxj4nYG9ASABKAVSFHJldGVudGlvbn'
    'BlcmlvZGhvdXJzEioKBnNoYXJkcxjJxbfQASADKAsyDi5raW5lc2lzLlNoYXJkUgZzaGFyZHMS'
    'IAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJuEjsKF3N0cmVhbWNyZWF0aW9udGltZX'
    'N0YW1wGOX1oWsgASgJUhdzdHJlYW1jcmVhdGlvbnRpbWVzdGFtcBJLChFzdHJlYW1tb2RlZGV0'
    'YWlscxiR+eQFIAEoCzIaLmtpbmVzaXMuU3RyZWFtTW9kZURldGFpbHNSEXN0cmVhbW1vZGVkZX'
    'RhaWxzEiIKCnN0cmVhbW5hbWUYx7e54AEgASgJUgpzdHJlYW1uYW1lEjwKDHN0cmVhbXN0YXR1'
    'cxjQgZp1IAEoDjIVLmtpbmVzaXMuU3RyZWFtU3RhdHVzUgxzdHJlYW1zdGF0dXM=');

@$core.Deprecated('Use streamDescriptionSummaryDescriptor instead')
const StreamDescriptionSummary$json = {
  '1': 'StreamDescriptionSummary',
  '2': [
    {
      '1': 'consumercount',
      '3': 448084721,
      '4': 1,
      '5': 5,
      '10': 'consumercount'
    },
    {
      '1': 'encryptiontype',
      '3': 264007605,
      '4': 1,
      '5': 14,
      '6': '.kinesis.EncryptionType',
      '10': 'encryptiontype'
    },
    {
      '1': 'enhancedmonitoring',
      '3': 452259826,
      '4': 3,
      '5': 11,
      '6': '.kinesis.EnhancedMetrics',
      '10': 'enhancedmonitoring'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'maxrecordsizeinkib',
      '3': 197267253,
      '4': 1,
      '5': 5,
      '10': 'maxrecordsizeinkib'
    },
    {
      '1': 'openshardcount',
      '3': 476409287,
      '4': 1,
      '5': 5,
      '10': 'openshardcount'
    },
    {
      '1': 'retentionperiodhours',
      '3': 396381944,
      '4': 1,
      '5': 5,
      '10': 'retentionperiodhours'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {
      '1': 'streamcreationtimestamp',
      '3': 224951013,
      '4': 1,
      '5': 9,
      '10': 'streamcreationtimestamp'
    },
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {
      '1': 'streammodedetails',
      '3': 12139665,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StreamModeDetails',
      '10': 'streammodedetails'
    },
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'streamstatus',
      '3': 245792976,
      '4': 1,
      '5': 14,
      '6': '.kinesis.StreamStatus',
      '10': 'streamstatus'
    },
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.kinesis.WarmThroughputObject',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `StreamDescriptionSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamDescriptionSummaryDescriptor = $convert.base64Decode(
    'ChhTdHJlYW1EZXNjcmlwdGlvblN1bW1hcnkSKAoNY29uc3VtZXJjb3VudBjx9dTVASABKAVSDW'
    'NvbnN1bWVyY291bnQSQgoOZW5jcnlwdGlvbnR5cGUYtd/xfSABKA4yFy5raW5lc2lzLkVuY3J5'
    'cHRpb25UeXBlUg5lbmNyeXB0aW9udHlwZRJMChJlbmhhbmNlZG1vbml0b3JpbmcY8t/T1wEgAy'
    'gLMhgua2luZXNpcy5FbmhhbmNlZE1ldHJpY3NSEmVuaGFuY2VkbW9uaXRvcmluZxIYCgVrZXlp'
    'ZBiigMiDASABKAlSBWtleWlkEjEKEm1heHJlY29yZHNpemVpbmtpYhi1noheIAEoBVISbWF4cm'
    'Vjb3Jkc2l6ZWlua2liEioKDm9wZW5zaGFyZGNvdW50GMfbleMBIAEoBVIOb3BlbnNoYXJkY291'
    'bnQSNgoUcmV0ZW50aW9ucGVyaW9kaG91cnMY+J2BvQEgASgFUhRyZXRlbnRpb25wZXJpb2Rob3'
    'VycxIgCglzdHJlYW1hcm4Y3fOq8gEgASgJUglzdHJlYW1hcm4SOwoXc3RyZWFtY3JlYXRpb250'
    'aW1lc3RhbXAY5fWhayABKAlSF3N0cmVhbWNyZWF0aW9udGltZXN0YW1wEh4KCHN0cmVhbWlkGM'
    'HtgcYBIAEoCVIIc3RyZWFtaWQSSwoRc3RyZWFtbW9kZWRldGFpbHMYkfnkBSABKAsyGi5raW5l'
    'c2lzLlN0cmVhbU1vZGVEZXRhaWxzUhFzdHJlYW1tb2RlZGV0YWlscxIiCgpzdHJlYW1uYW1lGM'
    'e3ueABIAEoCVIKc3RyZWFtbmFtZRI8CgxzdHJlYW1zdGF0dXMY0IGadSABKA4yFS5raW5lc2lz'
    'LlN0cmVhbVN0YXR1c1IMc3RyZWFtc3RhdHVzEkkKDndhcm10aHJvdWdocHV0GIPeyIoBIAEoCz'
    'IdLmtpbmVzaXMuV2FybVRocm91Z2hwdXRPYmplY3RSDndhcm10aHJvdWdocHV0');

@$core.Deprecated('Use streamModeDetailsDescriptor instead')
const StreamModeDetails$json = {
  '1': 'StreamModeDetails',
  '2': [
    {
      '1': 'streammode',
      '3': 457304819,
      '4': 1,
      '5': 14,
      '6': '.kinesis.StreamMode',
      '10': 'streammode'
    },
  ],
};

/// Descriptor for `StreamModeDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamModeDetailsDescriptor = $convert.base64Decode(
    'ChFTdHJlYW1Nb2RlRGV0YWlscxI3CgpzdHJlYW1tb2RlGPPVh9oBIAEoDjITLmtpbmVzaXMuU3'
    'RyZWFtTW9kZVIKc3RyZWFtbW9kZQ==');

@$core.Deprecated('Use streamSummaryDescriptor instead')
const StreamSummary$json = {
  '1': 'StreamSummary',
  '2': [
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {
      '1': 'streamcreationtimestamp',
      '3': 224951013,
      '4': 1,
      '5': 9,
      '10': 'streamcreationtimestamp'
    },
    {
      '1': 'streammodedetails',
      '3': 12139665,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StreamModeDetails',
      '10': 'streammodedetails'
    },
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'streamstatus',
      '3': 245792976,
      '4': 1,
      '5': 14,
      '6': '.kinesis.StreamStatus',
      '10': 'streamstatus'
    },
  ],
};

/// Descriptor for `StreamSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List streamSummaryDescriptor = $convert.base64Decode(
    'Cg1TdHJlYW1TdW1tYXJ5EiAKCXN0cmVhbWFybhjd86ryASABKAlSCXN0cmVhbWFybhI7ChdzdH'
    'JlYW1jcmVhdGlvbnRpbWVzdGFtcBjl9aFrIAEoCVIXc3RyZWFtY3JlYXRpb250aW1lc3RhbXAS'
    'SwoRc3RyZWFtbW9kZWRldGFpbHMYkfnkBSABKAsyGi5raW5lc2lzLlN0cmVhbU1vZGVEZXRhaW'
    'xzUhFzdHJlYW1tb2RlZGV0YWlscxIiCgpzdHJlYW1uYW1lGMe3ueABIAEoCVIKc3RyZWFtbmFt'
    'ZRI8CgxzdHJlYW1zdGF0dXMY0IGadSABKA4yFS5raW5lc2lzLlN0cmVhbVN0YXR1c1IMc3RyZW'
    'Ftc3RhdHVz');

@$core.Deprecated('Use subscribeToShardEventDescriptor instead')
const SubscribeToShardEvent$json = {
  '1': 'SubscribeToShardEvent',
  '2': [
    {
      '1': 'childshards',
      '3': 348740657,
      '4': 3,
      '5': 11,
      '6': '.kinesis.ChildShard',
      '10': 'childshards'
    },
    {
      '1': 'continuationsequencenumber',
      '3': 143046235,
      '4': 1,
      '5': 9,
      '10': 'continuationsequencenumber'
    },
    {
      '1': 'millisbehindlatest',
      '3': 456422057,
      '4': 1,
      '5': 3,
      '10': 'millisbehindlatest'
    },
    {
      '1': 'records',
      '3': 423557454,
      '4': 3,
      '5': 11,
      '6': '.kinesis.Record',
      '10': 'records'
    },
  ],
};

/// Descriptor for `SubscribeToShardEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeToShardEventDescriptor = $convert.base64Decode(
    'ChVTdWJzY3JpYmVUb1NoYXJkRXZlbnQSOQoLY2hpbGRzaGFyZHMYsbilpgEgAygLMhMua2luZX'
    'Npcy5DaGlsZFNoYXJkUgtjaGlsZHNoYXJkcxJBChpjb250aW51YXRpb25zZXF1ZW5jZW51bWJl'
    'chjb7JpEIAEoCVIaY29udGludWF0aW9uc2VxdWVuY2VudW1iZXISMgoSbWlsbGlzYmVoaW5kbG'
    'F0ZXN0GKnl0dkBIAEoA1ISbWlsbGlzYmVoaW5kbGF0ZXN0Ei0KB3JlY29yZHMYzvL7yQEgAygL'
    'Mg8ua2luZXNpcy5SZWNvcmRSB3JlY29yZHM=');

@$core.Deprecated('Use subscribeToShardEventStreamDescriptor instead')
const SubscribeToShardEventStream$json = {
  '1': 'SubscribeToShardEventStream',
  '2': [
    {
      '1': 'internalfailureexception',
      '3': 445206748,
      '4': 1,
      '5': 11,
      '6': '.kinesis.InternalFailureException',
      '10': 'internalfailureexception'
    },
    {
      '1': 'kmsaccessdeniedexception',
      '3': 110824201,
      '4': 1,
      '5': 11,
      '6': '.kinesis.KMSAccessDeniedException',
      '10': 'kmsaccessdeniedexception'
    },
    {
      '1': 'kmsdisabledexception',
      '3': 352104756,
      '4': 1,
      '5': 11,
      '6': '.kinesis.KMSDisabledException',
      '10': 'kmsdisabledexception'
    },
    {
      '1': 'kmsinvalidstateexception',
      '3': 183514916,
      '4': 1,
      '5': 11,
      '6': '.kinesis.KMSInvalidStateException',
      '10': 'kmsinvalidstateexception'
    },
    {
      '1': 'kmsnotfoundexception',
      '3': 251973779,
      '4': 1,
      '5': 11,
      '6': '.kinesis.KMSNotFoundException',
      '10': 'kmsnotfoundexception'
    },
    {
      '1': 'kmsoptinrequired',
      '3': 72448154,
      '4': 1,
      '5': 11,
      '6': '.kinesis.KMSOptInRequired',
      '10': 'kmsoptinrequired'
    },
    {
      '1': 'kmsthrottlingexception',
      '3': 109345891,
      '4': 1,
      '5': 11,
      '6': '.kinesis.KMSThrottlingException',
      '10': 'kmsthrottlingexception'
    },
    {
      '1': 'resourceinuseexception',
      '3': 360549483,
      '4': 1,
      '5': 11,
      '6': '.kinesis.ResourceInUseException',
      '10': 'resourceinuseexception'
    },
    {
      '1': 'resourcenotfoundexception',
      '3': 396489184,
      '4': 1,
      '5': 11,
      '6': '.kinesis.ResourceNotFoundException',
      '10': 'resourcenotfoundexception'
    },
    {
      '1': 'subscribetoshardevent',
      '3': 309751871,
      '4': 1,
      '5': 11,
      '6': '.kinesis.SubscribeToShardEvent',
      '10': 'subscribetoshardevent'
    },
  ],
};

/// Descriptor for `SubscribeToShardEventStream`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeToShardEventStreamDescriptor = $convert.base64Decode(
    'ChtTdWJzY3JpYmVUb1NoYXJkRXZlbnRTdHJlYW0SYQoYaW50ZXJuYWxmYWlsdXJlZXhjZXB0aW'
    '9uGNyhpdQBIAEoCzIhLmtpbmVzaXMuSW50ZXJuYWxGYWlsdXJlRXhjZXB0aW9uUhhpbnRlcm5h'
    'bGZhaWx1cmVleGNlcHRpb24SYAoYa21zYWNjZXNzZGVuaWVkZXhjZXB0aW9uGImW7DQgASgLMi'
    'Eua2luZXNpcy5LTVNBY2Nlc3NEZW5pZWRFeGNlcHRpb25SGGttc2FjY2Vzc2RlbmllZGV4Y2Vw'
    'dGlvbhJVChRrbXNkaXNhYmxlZGV4Y2VwdGlvbhi04vKnASABKAsyHS5raW5lc2lzLktNU0Rpc2'
    'FibGVkRXhjZXB0aW9uUhRrbXNkaXNhYmxlZGV4Y2VwdGlvbhJgChhrbXNpbnZhbGlkc3RhdGVl'
    'eGNlcHRpb24YpO7AVyABKAsyIS5raW5lc2lzLktNU0ludmFsaWRTdGF0ZUV4Y2VwdGlvblIYa2'
    '1zaW52YWxpZHN0YXRlZXhjZXB0aW9uElQKFGttc25vdGZvdW5kZXhjZXB0aW9uGJOhk3ggASgL'
    'Mh0ua2luZXNpcy5LTVNOb3RGb3VuZEV4Y2VwdGlvblIUa21zbm90Zm91bmRleGNlcHRpb24SSA'
    'oQa21zb3B0aW5yZXF1aXJlZBia8cUiIAEoCzIZLmtpbmVzaXMuS01TT3B0SW5SZXF1aXJlZFIQ'
    'a21zb3B0aW5yZXF1aXJlZBJaChZrbXN0aHJvdHRsaW5nZXhjZXB0aW9uGOP4kTQgASgLMh8ua2'
    'luZXNpcy5LTVNUaHJvdHRsaW5nRXhjZXB0aW9uUhZrbXN0aHJvdHRsaW5nZXhjZXB0aW9uElsK'
    'FnJlc291cmNlaW51c2VleGNlcHRpb24Y65j2qwEgASgLMh8ua2luZXNpcy5SZXNvdXJjZUluVX'
    'NlRXhjZXB0aW9uUhZyZXNvdXJjZWludXNlZXhjZXB0aW9uEmQKGXJlc291cmNlbm90Zm91bmRl'
    'eGNlcHRpb24Y4OOHvQEgASgLMiIua2luZXNpcy5SZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uUh'
    'lyZXNvdXJjZW5vdGZvdW5kZXhjZXB0aW9uElgKFXN1YnNjcmliZXRvc2hhcmRldmVudBi/4NmT'
    'ASABKAsyHi5raW5lc2lzLlN1YnNjcmliZVRvU2hhcmRFdmVudFIVc3Vic2NyaWJldG9zaGFyZG'
    'V2ZW50');

@$core.Deprecated('Use subscribeToShardInputDescriptor instead')
const SubscribeToShardInput$json = {
  '1': 'SubscribeToShardInput',
  '2': [
    {'1': 'consumerarn', '3': 41107441, '4': 1, '5': 9, '10': 'consumerarn'},
    {'1': 'shardid', '3': 66410951, '4': 1, '5': 9, '10': 'shardid'},
    {
      '1': 'startingposition',
      '3': 428771919,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StartingPosition',
      '10': 'startingposition'
    },
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `SubscribeToShardInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeToShardInputDescriptor = $convert.base64Decode(
    'ChVTdWJzY3JpYmVUb1NoYXJkSW5wdXQSIwoLY29uc3VtZXJhcm4Y8f/MEyABKAlSC2NvbnN1bW'
    'VyYXJuEhsKB3NoYXJkaWQYx7PVHyABKAlSB3NoYXJkaWQSSQoQc3RhcnRpbmdwb3NpdGlvbhjP'
    'lLrMASABKAsyGS5raW5lc2lzLlN0YXJ0aW5nUG9zaXRpb25SEHN0YXJ0aW5ncG9zaXRpb24SHg'
    'oIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZA==');

@$core.Deprecated('Use subscribeToShardOutputDescriptor instead')
const SubscribeToShardOutput$json = {
  '1': 'SubscribeToShardOutput',
  '2': [
    {
      '1': 'eventstream',
      '3': 26857888,
      '4': 1,
      '5': 11,
      '6': '.kinesis.SubscribeToShardEventStream',
      '10': 'eventstream'
    },
  ],
};

/// Descriptor for `SubscribeToShardOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscribeToShardOutputDescriptor =
    $convert.base64Decode(
        'ChZTdWJzY3JpYmVUb1NoYXJkT3V0cHV0EkkKC2V2ZW50c3RyZWFtGKCj5wwgASgLMiQua2luZX'
        'Npcy5TdWJzY3JpYmVUb1NoYXJkRXZlbnRTdHJlYW1SC2V2ZW50c3RyZWFt');

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

@$core.Deprecated('Use tagResourceInputDescriptor instead')
const TagResourceInput$json = {
  '1': 'TagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kinesis.TagResourceInput.TagsEntry',
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
    'ChBUYWdSZXNvdXJjZUlucHV0EiQKC3Jlc291cmNlYXJuGO3AmbABIAEoCVILcmVzb3VyY2Vhcm'
    '4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZBI7CgR0YWdzGMHB9rUBIAMoCzIjLmtp'
    'bmVzaXMuVGFnUmVzb3VyY2VJbnB1dC5UYWdzRW50cnlSBHRhZ3MaNwoJVGFnc0VudHJ5EhAKA2'
    'tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use untagResourceInputDescriptor instead')
const UntagResourceInput$json = {
  '1': 'UntagResourceInput',
  '2': [
    {'1': 'resourcearn', '3': 369516653, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceInputDescriptor = $convert.base64Decode(
    'ChJVbnRhZ1Jlc291cmNlSW5wdXQSJAoLcmVzb3VyY2Vhcm4Y7cCZsAEgASgJUgtyZXNvdXJjZW'
    'FybhIeCghzdHJlYW1pZBjB7YHGASABKAlSCHN0cmVhbWlkEhwKB3RhZ2tleXMY/MPzmAEgAygJ'
    'Ugd0YWdrZXlz');

@$core.Deprecated('Use updateAccountSettingsInputDescriptor instead')
const UpdateAccountSettingsInput$json = {
  '1': 'UpdateAccountSettingsInput',
  '2': [
    {
      '1': 'minimumthroughputbillingcommitment',
      '3': 335326962,
      '4': 1,
      '5': 11,
      '6': '.kinesis.MinimumThroughputBillingCommitmentInput',
      '10': 'minimumthroughputbillingcommitment'
    },
  ],
};

/// Descriptor for `UpdateAccountSettingsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAccountSettingsInputDescriptor = $convert.base64Decode(
    'ChpVcGRhdGVBY2NvdW50U2V0dGluZ3NJbnB1dBKEAQoibWluaW11bXRocm91Z2hwdXRiaWxsaW'
    '5nY29tbWl0bWVudBjy3fKfASABKAsyMC5raW5lc2lzLk1pbmltdW1UaHJvdWdocHV0QmlsbGlu'
    'Z0NvbW1pdG1lbnRJbnB1dFIibWluaW11bXRocm91Z2hwdXRiaWxsaW5nY29tbWl0bWVudA==');

@$core.Deprecated('Use updateAccountSettingsOutputDescriptor instead')
const UpdateAccountSettingsOutput$json = {
  '1': 'UpdateAccountSettingsOutput',
  '2': [
    {
      '1': 'minimumthroughputbillingcommitment',
      '3': 335326962,
      '4': 1,
      '5': 11,
      '6': '.kinesis.MinimumThroughputBillingCommitmentOutput',
      '10': 'minimumthroughputbillingcommitment'
    },
  ],
};

/// Descriptor for `UpdateAccountSettingsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAccountSettingsOutputDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVBY2NvdW50U2V0dGluZ3NPdXRwdXQShQEKIm1pbmltdW10aHJvdWdocHV0YmlsbG'
    'luZ2NvbW1pdG1lbnQY8t3ynwEgASgLMjEua2luZXNpcy5NaW5pbXVtVGhyb3VnaHB1dEJpbGxp'
    'bmdDb21taXRtZW50T3V0cHV0UiJtaW5pbXVtdGhyb3VnaHB1dGJpbGxpbmdjb21taXRtZW50');

@$core.Deprecated('Use updateMaxRecordSizeInputDescriptor instead')
const UpdateMaxRecordSizeInput$json = {
  '1': 'UpdateMaxRecordSizeInput',
  '2': [
    {
      '1': 'maxrecordsizeinkib',
      '3': 197267253,
      '4': 1,
      '5': 5,
      '10': 'maxrecordsizeinkib'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
  ],
};

/// Descriptor for `UpdateMaxRecordSizeInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateMaxRecordSizeInputDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVNYXhSZWNvcmRTaXplSW5wdXQSMQoSbWF4cmVjb3Jkc2l6ZWlua2liGLWeiF4gAS'
    'gFUhJtYXhyZWNvcmRzaXplaW5raWISIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYXJu'
    'Eh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQ=');

@$core.Deprecated('Use updateShardCountInputDescriptor instead')
const UpdateShardCountInput$json = {
  '1': 'UpdateShardCountInput',
  '2': [
    {
      '1': 'scalingtype',
      '3': 34064531,
      '4': 1,
      '5': 14,
      '6': '.kinesis.ScalingType',
      '10': 'scalingtype'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'targetshardcount',
      '3': 361168816,
      '4': 1,
      '5': 5,
      '10': 'targetshardcount'
    },
  ],
};

/// Descriptor for `UpdateShardCountInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateShardCountInputDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVTaGFyZENvdW50SW5wdXQSOQoLc2NhbGluZ3R5cGUYk5GfECABKA4yFC5raW5lc2'
    'lzLlNjYWxpbmdUeXBlUgtzY2FsaW5ndHlwZRIgCglzdHJlYW1hcm4Y3fOq8gEgASgJUglzdHJl'
    'YW1hcm4SHgoIc3RyZWFtaWQYwe2BxgEgASgJUghzdHJlYW1pZBIiCgpzdHJlYW1uYW1lGMe3ue'
    'ABIAEoCVIKc3RyZWFtbmFtZRIuChB0YXJnZXRzaGFyZGNvdW50GLD/m6wBIAEoBVIQdGFyZ2V0'
    'c2hhcmRjb3VudA==');

@$core.Deprecated('Use updateShardCountOutputDescriptor instead')
const UpdateShardCountOutput$json = {
  '1': 'UpdateShardCountOutput',
  '2': [
    {
      '1': 'currentshardcount',
      '3': 21258174,
      '4': 1,
      '5': 5,
      '10': 'currentshardcount'
    },
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'targetshardcount',
      '3': 361168816,
      '4': 1,
      '5': 5,
      '10': 'targetshardcount'
    },
  ],
};

/// Descriptor for `UpdateShardCountOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateShardCountOutputDescriptor = $convert.base64Decode(
    'ChZVcGRhdGVTaGFyZENvdW50T3V0cHV0Ei8KEWN1cnJlbnRzaGFyZGNvdW50GL6/kQogASgFUh'
    'FjdXJyZW50c2hhcmRjb3VudBIgCglzdHJlYW1hcm4Y3fOq8gEgASgJUglzdHJlYW1hcm4SIgoK'
    'c3RyZWFtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWUSLgoQdGFyZ2V0c2hhcmRjb3VudBiw/5'
    'usASABKAVSEHRhcmdldHNoYXJkY291bnQ=');

@$core.Deprecated('Use updateStreamModeInputDescriptor instead')
const UpdateStreamModeInput$json = {
  '1': 'UpdateStreamModeInput',
  '2': [
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {
      '1': 'streammodedetails',
      '3': 12139665,
      '4': 1,
      '5': 11,
      '6': '.kinesis.StreamModeDetails',
      '10': 'streammodedetails'
    },
    {
      '1': 'warmthroughputmibps',
      '3': 259219704,
      '4': 1,
      '5': 5,
      '10': 'warmthroughputmibps'
    },
  ],
};

/// Descriptor for `UpdateStreamModeInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStreamModeInputDescriptor = $convert.base64Decode(
    'ChVVcGRhdGVTdHJlYW1Nb2RlSW5wdXQSIAoJc3RyZWFtYXJuGN3zqvIBIAEoCVIJc3RyZWFtYX'
    'JuEh4KCHN0cmVhbWlkGMHtgcYBIAEoCVIIc3RyZWFtaWQSSwoRc3RyZWFtbW9kZWRldGFpbHMY'
    'kfnkBSABKAsyGi5raW5lc2lzLlN0cmVhbU1vZGVEZXRhaWxzUhFzdHJlYW1tb2RlZGV0YWlscx'
    'IzChN3YXJtdGhyb3VnaHB1dG1pYnBzGPjBzXsgASgFUhN3YXJtdGhyb3VnaHB1dG1pYnBz');

@$core.Deprecated('Use updateStreamWarmThroughputInputDescriptor instead')
const UpdateStreamWarmThroughputInput$json = {
  '1': 'UpdateStreamWarmThroughputInput',
  '2': [
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamid', '3': 415266497, '4': 1, '5': 9, '10': 'streamid'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'warmthroughputmibps',
      '3': 259219704,
      '4': 1,
      '5': 5,
      '10': 'warmthroughputmibps'
    },
  ],
};

/// Descriptor for `UpdateStreamWarmThroughputInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStreamWarmThroughputInputDescriptor =
    $convert.base64Decode(
        'Ch9VcGRhdGVTdHJlYW1XYXJtVGhyb3VnaHB1dElucHV0EiAKCXN0cmVhbWFybhjd86ryASABKA'
        'lSCXN0cmVhbWFybhIeCghzdHJlYW1pZBjB7YHGASABKAlSCHN0cmVhbWlkEiIKCnN0cmVhbW5h'
        'bWUYx7e54AEgASgJUgpzdHJlYW1uYW1lEjMKE3dhcm10aHJvdWdocHV0bWlicHMY+MHNeyABKA'
        'VSE3dhcm10aHJvdWdocHV0bWlicHM=');

@$core.Deprecated('Use updateStreamWarmThroughputOutputDescriptor instead')
const UpdateStreamWarmThroughputOutput$json = {
  '1': 'UpdateStreamWarmThroughputOutput',
  '2': [
    {'1': 'streamarn', '3': 508213725, '4': 1, '5': 9, '10': 'streamarn'},
    {'1': 'streamname', '3': 470703047, '4': 1, '5': 9, '10': 'streamname'},
    {
      '1': 'warmthroughput',
      '3': 290598659,
      '4': 1,
      '5': 11,
      '6': '.kinesis.WarmThroughputObject',
      '10': 'warmthroughput'
    },
  ],
};

/// Descriptor for `UpdateStreamWarmThroughputOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateStreamWarmThroughputOutputDescriptor =
    $convert.base64Decode(
        'CiBVcGRhdGVTdHJlYW1XYXJtVGhyb3VnaHB1dE91dHB1dBIgCglzdHJlYW1hcm4Y3fOq8gEgAS'
        'gJUglzdHJlYW1hcm4SIgoKc3RyZWFtbmFtZRjHt7ngASABKAlSCnN0cmVhbW5hbWUSSQoOd2Fy'
        'bXRocm91Z2hwdXQYg97IigEgASgLMh0ua2luZXNpcy5XYXJtVGhyb3VnaHB1dE9iamVjdFIOd2'
        'FybXRocm91Z2hwdXQ=');

@$core.Deprecated('Use validationExceptionDescriptor instead')
const ValidationException$json = {
  '1': 'ValidationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ValidationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List validationExceptionDescriptor =
    $convert.base64Decode(
        'ChNWYWxpZGF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use warmThroughputObjectDescriptor instead')
const WarmThroughputObject$json = {
  '1': 'WarmThroughputObject',
  '2': [
    {'1': 'currentmibps', '3': 153746976, '4': 1, '5': 5, '10': 'currentmibps'},
    {'1': 'targetmibps', '3': 535445626, '4': 1, '5': 5, '10': 'targetmibps'},
  ],
};

/// Descriptor for `WarmThroughputObject`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List warmThroughputObjectDescriptor = $convert.base64Decode(
    'ChRXYXJtVGhyb3VnaHB1dE9iamVjdBIlCgxjdXJyZW50bWlicHMYoPynSSABKAVSDGN1cnJlbn'
    'RtaWJwcxIkCgt0YXJnZXRtaWJwcxj6gKn/ASABKAVSC3RhcmdldG1pYnBz');
