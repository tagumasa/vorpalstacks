// This is a generated file - do not edit.
//
// Generated from sqs.proto.

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

@$core.Deprecated('Use messageSystemAttributeNameDescriptor instead')
const MessageSystemAttributeName$json = {
  '1': 'MessageSystemAttributeName',
  '2': [
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_SEQUENCENUMBER', '2': 0},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEDEDUPLICATIONID', '2': 1},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_DEADLETTERQUEUESOURCEARN', '2': 2},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENTTIMESTAMP', '2': 3},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENDERID', '2': 4},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_ALL', '2': 5},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEGROUPID', '2': 6},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATERECEIVECOUNT', '2': 7},
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_AWSTRACEHEADER', '2': 8},
    {
      '1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATEFIRSTRECEIVETIMESTAMP',
      '2': 9
    },
  ],
};

/// Descriptor for `MessageSystemAttributeName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List messageSystemAttributeNameDescriptor = $convert.base64Decode(
    'ChpNZXNzYWdlU3lzdGVtQXR0cmlidXRlTmFtZRIwCixNRVNTQUdFX1NZU1RFTV9BVFRSSUJVVE'
    'VfTkFNRV9TRVFVRU5DRU5VTUJFUhAAEjgKNE1FU1NBR0VfU1lTVEVNX0FUVFJJQlVURV9OQU1F'
    'X01FU1NBR0VERURVUExJQ0FUSU9OSUQQARI6CjZNRVNTQUdFX1NZU1RFTV9BVFRSSUJVVEVfTk'
    'FNRV9ERUFETEVUVEVSUVVFVUVTT1VSQ0VBUk4QAhIvCitNRVNTQUdFX1NZU1RFTV9BVFRSSUJV'
    'VEVfTkFNRV9TRU5UVElNRVNUQU1QEAMSKgomTUVTU0FHRV9TWVNURU1fQVRUUklCVVRFX05BTU'
    'VfU0VOREVSSUQQBBIlCiFNRVNTQUdFX1NZU1RFTV9BVFRSSUJVVEVfTkFNRV9BTEwQBRIwCixN'
    'RVNTQUdFX1NZU1RFTV9BVFRSSUJVVEVfTkFNRV9NRVNTQUdFR1JPVVBJRBAGEjkKNU1FU1NBR0'
    'VfU1lTVEVNX0FUVFJJQlVURV9OQU1FX0FQUFJPWElNQVRFUkVDRUlWRUNPVU5UEAcSMAosTUVT'
    'U0FHRV9TWVNURU1fQVRUUklCVVRFX05BTUVfQVdTVFJBQ0VIRUFERVIQCBJCCj5NRVNTQUdFX1'
    'NZU1RFTV9BVFRSSUJVVEVfTkFNRV9BUFBST1hJTUFURUZJUlNUUkVDRUlWRVRJTUVTVEFNUBAJ');

@$core.Deprecated('Use messageSystemAttributeNameForSendsDescriptor instead')
const MessageSystemAttributeNameForSends$json = {
  '1': 'MessageSystemAttributeNameForSends',
  '2': [
    {'1': 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_FOR_SENDS_AWSTRACEHEADER', '2': 0},
  ],
};

/// Descriptor for `MessageSystemAttributeNameForSends`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List messageSystemAttributeNameForSendsDescriptor =
    $convert.base64Decode(
        'CiJNZXNzYWdlU3lzdGVtQXR0cmlidXRlTmFtZUZvclNlbmRzEjoKNk1FU1NBR0VfU1lTVEVNX0'
        'FUVFJJQlVURV9OQU1FX0ZPUl9TRU5EU19BV1NUUkFDRUhFQURFUhAA');

@$core.Deprecated('Use queueAttributeNameDescriptor instead')
const QueueAttributeName$json = {
  '1': 'QueueAttributeName',
  '2': [
    {'1': 'QUEUE_ATTRIBUTE_NAME_MAXIMUMMESSAGESIZE', '2': 0},
    {'1': 'QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESDELAYED', '2': 1},
    {'1': 'QUEUE_ATTRIBUTE_NAME_DEDUPLICATIONSCOPE', '2': 2},
    {'1': 'QUEUE_ATTRIBUTE_NAME_MESSAGERETENTIONPERIOD', '2': 3},
    {'1': 'QUEUE_ATTRIBUTE_NAME_DELAYSECONDS', '2': 4},
    {'1': 'QUEUE_ATTRIBUTE_NAME_QUEUEARN', '2': 5},
    {'1': 'QUEUE_ATTRIBUTE_NAME_FIFOTHROUGHPUTLIMIT', '2': 6},
    {'1': 'QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGES', '2': 7},
    {'1': 'QUEUE_ATTRIBUTE_NAME_FIFOQUEUE', '2': 8},
    {'1': 'QUEUE_ATTRIBUTE_NAME_CREATEDTIMESTAMP', '2': 9},
    {'1': 'QUEUE_ATTRIBUTE_NAME_ALL', '2': 10},
    {'1': 'QUEUE_ATTRIBUTE_NAME_KMSMASTERKEYID', '2': 11},
    {'1': 'QUEUE_ATTRIBUTE_NAME_LASTMODIFIEDTIMESTAMP', '2': 12},
    {'1': 'QUEUE_ATTRIBUTE_NAME_REDRIVEALLOWPOLICY', '2': 13},
    {'1': 'QUEUE_ATTRIBUTE_NAME_RECEIVEMESSAGEWAITTIMESECONDS', '2': 14},
    {
      '1': 'QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESNOTVISIBLE',
      '2': 15
    },
    {'1': 'QUEUE_ATTRIBUTE_NAME_SQSMANAGEDSSEENABLED', '2': 16},
    {'1': 'QUEUE_ATTRIBUTE_NAME_REDRIVEPOLICY', '2': 17},
    {'1': 'QUEUE_ATTRIBUTE_NAME_POLICY', '2': 18},
    {'1': 'QUEUE_ATTRIBUTE_NAME_KMSDATAKEYREUSEPERIODSECONDS', '2': 19},
    {'1': 'QUEUE_ATTRIBUTE_NAME_VISIBILITYTIMEOUT', '2': 20},
    {'1': 'QUEUE_ATTRIBUTE_NAME_CONTENTBASEDDEDUPLICATION', '2': 21},
  ],
};

/// Descriptor for `QueueAttributeName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queueAttributeNameDescriptor = $convert.base64Decode(
    'ChJRdWV1ZUF0dHJpYnV0ZU5hbWUSKwonUVVFVUVfQVRUUklCVVRFX05BTUVfTUFYSU1VTU1FU1'
    'NBR0VTSVpFEAASOwo3UVVFVUVfQVRUUklCVVRFX05BTUVfQVBQUk9YSU1BVEVOVU1CRVJPRk1F'
    'U1NBR0VTREVMQVlFRBABEisKJ1FVRVVFX0FUVFJJQlVURV9OQU1FX0RFRFVQTElDQVRJT05TQ0'
    '9QRRACEi8KK1FVRVVFX0FUVFJJQlVURV9OQU1FX01FU1NBR0VSRVRFTlRJT05QRVJJT0QQAxIl'
    'CiFRVUVVRV9BVFRSSUJVVEVfTkFNRV9ERUxBWVNFQ09ORFMQBBIhCh1RVUVVRV9BVFRSSUJVVE'
    'VfTkFNRV9RVUVVRUFSThAFEiwKKFFVRVVFX0FUVFJJQlVURV9OQU1FX0ZJRk9USFJPVUdIUFVU'
    'TElNSVQQBhI0CjBRVUVVRV9BVFRSSUJVVEVfTkFNRV9BUFBST1hJTUFURU5VTUJFUk9GTUVTU0'
    'FHRVMQBxIiCh5RVUVVRV9BVFRSSUJVVEVfTkFNRV9GSUZPUVVFVUUQCBIpCiVRVUVVRV9BVFRS'
    'SUJVVEVfTkFNRV9DUkVBVEVEVElNRVNUQU1QEAkSHAoYUVVFVUVfQVRUUklCVVRFX05BTUVfQU'
    'xMEAoSJwojUVVFVUVfQVRUUklCVVRFX05BTUVfS01TTUFTVEVSS0VZSUQQCxIuCipRVUVVRV9B'
    'VFRSSUJVVEVfTkFNRV9MQVNUTU9ESUZJRURUSU1FU1RBTVAQDBIrCidRVUVVRV9BVFRSSUJVVE'
    'VfTkFNRV9SRURSSVZFQUxMT1dQT0xJQ1kQDRI2CjJRVUVVRV9BVFRSSUJVVEVfTkFNRV9SRUNF'
    'SVZFTUVTU0FHRVdBSVRUSU1FU0VDT05EUxAOEj4KOlFVRVVFX0FUVFJJQlVURV9OQU1FX0FQUF'
    'JPWElNQVRFTlVNQkVST0ZNRVNTQUdFU05PVFZJU0lCTEUQDxItCilRVUVVRV9BVFRSSUJVVEVf'
    'TkFNRV9TUVNNQU5BR0VEU1NFRU5BQkxFRBAQEiYKIlFVRVVFX0FUVFJJQlVURV9OQU1FX1JFRF'
    'JJVkVQT0xJQ1kQERIfChtRVUVVRV9BVFRSSUJVVEVfTkFNRV9QT0xJQ1kQEhI1CjFRVUVVRV9B'
    'VFRSSUJVVEVfTkFNRV9LTVNEQVRBS0VZUkVVU0VQRVJJT0RTRUNPTkRTEBMSKgomUVVFVUVfQV'
    'RUUklCVVRFX05BTUVfVklTSUJJTElUWVRJTUVPVVQQFBIyCi5RVUVVRV9BVFRSSUJVVEVfTkFN'
    'RV9DT05URU5UQkFTRURERURVUExJQ0FUSU9OEBU=');

@$core.Deprecated('Use addPermissionRequestDescriptor instead')
const AddPermissionRequest$json = {
  '1': 'AddPermissionRequest',
  '2': [
    {
      '1': 'awsaccountids',
      '3': 417597390,
      '4': 3,
      '5': 9,
      '10': 'awsaccountids'
    },
    {'1': 'actions', '3': 106935557, '4': 3, '5': 9, '10': 'actions'},
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `AddPermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List addPermissionRequestDescriptor = $convert.base64Decode(
    'ChRBZGRQZXJtaXNzaW9uUmVxdWVzdBIoCg1hd3NhY2NvdW50aWRzGM6PkMcBIAMoCVINYXdzYW'
    'Njb3VudGlkcxIbCgdhY3Rpb25zGIXq/jIgAygJUgdhY3Rpb25zEhgKBWxhYmVsGJ7ls/YBIAEo'
    'CVIFbGFiZWwSHgoIcXVldWV1cmwYysG+8wEgASgJUghxdWV1ZXVybA==');

@$core.Deprecated('Use batchEntryIdsNotDistinctDescriptor instead')
const BatchEntryIdsNotDistinct$json = {
  '1': 'BatchEntryIdsNotDistinct',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BatchEntryIdsNotDistinct`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchEntryIdsNotDistinctDescriptor =
    $convert.base64Decode(
        'ChhCYXRjaEVudHJ5SWRzTm90RGlzdGluY3QSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use batchRequestTooLongDescriptor instead')
const BatchRequestTooLong$json = {
  '1': 'BatchRequestTooLong',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BatchRequestTooLong`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchRequestTooLongDescriptor =
    $convert.base64Decode(
        'ChNCYXRjaFJlcXVlc3RUb29Mb25nEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

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

@$core.Deprecated('Use cancelMessageMoveTaskRequestDescriptor instead')
const CancelMessageMoveTaskRequest$json = {
  '1': 'CancelMessageMoveTaskRequest',
  '2': [
    {'1': 'taskhandle', '3': 190544291, '4': 1, '5': 9, '10': 'taskhandle'},
  ],
};

/// Descriptor for `CancelMessageMoveTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelMessageMoveTaskRequestDescriptor =
    $convert.base64Decode(
        'ChxDYW5jZWxNZXNzYWdlTW92ZVRhc2tSZXF1ZXN0EiEKCnRhc2toYW5kbGUYo/PtWiABKAlSCn'
        'Rhc2toYW5kbGU=');

@$core.Deprecated('Use cancelMessageMoveTaskResultDescriptor instead')
const CancelMessageMoveTaskResult$json = {
  '1': 'CancelMessageMoveTaskResult',
  '2': [
    {
      '1': 'approximatenumberofmessagesmoved',
      '3': 510181657,
      '4': 1,
      '5': 3,
      '10': 'approximatenumberofmessagesmoved'
    },
  ],
};

/// Descriptor for `CancelMessageMoveTaskResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelMessageMoveTaskResultDescriptor =
    $convert.base64Decode(
        'ChtDYW5jZWxNZXNzYWdlTW92ZVRhc2tSZXN1bHQSTgogYXBwcm94aW1hdGVudW1iZXJvZm1lc3'
        'NhZ2VzbW92ZWQYmYKj8wEgASgDUiBhcHByb3hpbWF0ZW51bWJlcm9mbWVzc2FnZXNtb3ZlZA==');

@$core.Deprecated('Use changeMessageVisibilityBatchRequestDescriptor instead')
const ChangeMessageVisibilityBatchRequest$json = {
  '1': 'ChangeMessageVisibilityBatchRequest',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.sqs.ChangeMessageVisibilityBatchRequestEntry',
      '10': 'entries'
    },
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `ChangeMessageVisibilityBatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeMessageVisibilityBatchRequestDescriptor =
    $convert.base64Decode(
        'CiNDaGFuZ2VNZXNzYWdlVmlzaWJpbGl0eUJhdGNoUmVxdWVzdBJLCgdlbnRyaWVzGJTFsuUBIA'
        'MoCzItLnNxcy5DaGFuZ2VNZXNzYWdlVmlzaWJpbGl0eUJhdGNoUmVxdWVzdEVudHJ5UgdlbnRy'
        'aWVzEh4KCHF1ZXVldXJsGMrBvvMBIAEoCVIIcXVldWV1cmw=');

@$core.Deprecated(
    'Use changeMessageVisibilityBatchRequestEntryDescriptor instead')
const ChangeMessageVisibilityBatchRequestEntry$json = {
  '1': 'ChangeMessageVisibilityBatchRequestEntry',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'receipthandle',
      '3': 134471750,
      '4': 1,
      '5': 9,
      '10': 'receipthandle'
    },
    {
      '1': 'visibilitytimeout',
      '3': 460820073,
      '4': 1,
      '5': 5,
      '10': 'visibilitytimeout'
    },
  ],
};

/// Descriptor for `ChangeMessageVisibilityBatchRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeMessageVisibilityBatchRequestEntryDescriptor =
    $convert.base64Decode(
        'CihDaGFuZ2VNZXNzYWdlVmlzaWJpbGl0eUJhdGNoUmVxdWVzdEVudHJ5EhIKAmlkGIHyorcBIA'
        'EoCVICaWQSJwoNcmVjZWlwdGhhbmRsZRjGwI9AIAEoCVINcmVjZWlwdGhhbmRsZRIwChF2aXNp'
        'YmlsaXR5dGltZW91dBjpnN7bASABKAVSEXZpc2liaWxpdHl0aW1lb3V0');

@$core.Deprecated('Use changeMessageVisibilityBatchResultDescriptor instead')
const ChangeMessageVisibilityBatchResult$json = {
  '1': 'ChangeMessageVisibilityBatchResult',
  '2': [
    {
      '1': 'failed',
      '3': 360301525,
      '4': 3,
      '5': 11,
      '6': '.sqs.BatchResultErrorEntry',
      '10': 'failed'
    },
    {
      '1': 'successful',
      '3': 412818844,
      '4': 3,
      '5': 11,
      '6': '.sqs.ChangeMessageVisibilityBatchResultEntry',
      '10': 'successful'
    },
  ],
};

/// Descriptor for `ChangeMessageVisibilityBatchResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeMessageVisibilityBatchResultDescriptor =
    $convert.base64Decode(
        'CiJDaGFuZ2VNZXNzYWdlVmlzaWJpbGl0eUJhdGNoUmVzdWx0EjYKBmZhaWxlZBjVh+erASADKA'
        'syGi5zcXMuQmF0Y2hSZXN1bHRFcnJvckVudHJ5UgZmYWlsZWQSUAoKc3VjY2Vzc2Z1bBicu+zE'
        'ASADKAsyLC5zcXMuQ2hhbmdlTWVzc2FnZVZpc2liaWxpdHlCYXRjaFJlc3VsdEVudHJ5UgpzdW'
        'NjZXNzZnVs');

@$core
    .Deprecated('Use changeMessageVisibilityBatchResultEntryDescriptor instead')
const ChangeMessageVisibilityBatchResultEntry$json = {
  '1': 'ChangeMessageVisibilityBatchResultEntry',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `ChangeMessageVisibilityBatchResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeMessageVisibilityBatchResultEntryDescriptor =
    $convert.base64Decode(
        'CidDaGFuZ2VNZXNzYWdlVmlzaWJpbGl0eUJhdGNoUmVzdWx0RW50cnkSEgoCaWQYgfKitwEgAS'
        'gJUgJpZA==');

@$core.Deprecated('Use changeMessageVisibilityRequestDescriptor instead')
const ChangeMessageVisibilityRequest$json = {
  '1': 'ChangeMessageVisibilityRequest',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
    {
      '1': 'receipthandle',
      '3': 134471750,
      '4': 1,
      '5': 9,
      '10': 'receipthandle'
    },
    {
      '1': 'visibilitytimeout',
      '3': 460820073,
      '4': 1,
      '5': 5,
      '10': 'visibilitytimeout'
    },
  ],
};

/// Descriptor for `ChangeMessageVisibilityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List changeMessageVisibilityRequestDescriptor =
    $convert.base64Decode(
        'Ch5DaGFuZ2VNZXNzYWdlVmlzaWJpbGl0eVJlcXVlc3QSHgoIcXVldWV1cmwYysG+8wEgASgJUg'
        'hxdWV1ZXVybBInCg1yZWNlaXB0aGFuZGxlGMbAj0AgASgJUg1yZWNlaXB0aGFuZGxlEjAKEXZp'
        'c2liaWxpdHl0aW1lb3V0GOmc3tsBIAEoBVIRdmlzaWJpbGl0eXRpbWVvdXQ=');

@$core.Deprecated('Use createQueueRequestDescriptor instead')
const CreateQueueRequest$json = {
  '1': 'CreateQueueRequest',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sqs.CreateQueueRequest.AttributesEntry',
      '10': 'attributes'
    },
    {'1': 'queuename', '3': 154364636, '4': 1, '5': 9, '10': 'queuename'},
    {
      '1': 'tags',
      '3': 337046433,
      '4': 3,
      '5': 11,
      '6': '.sqs.CreateQueueRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [
    CreateQueueRequest_AttributesEntry$json,
    CreateQueueRequest_TagsEntry$json
  ],
};

@$core.Deprecated('Use createQueueRequestDescriptor instead')
const CreateQueueRequest_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use createQueueRequestDescriptor instead')
const CreateQueueRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `CreateQueueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createQueueRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVRdWV1ZVJlcXVlc3QSSgoKYXR0cmlidXRlcxi1qftjIAMoCzInLnNxcy5DcmVhdG'
    'VRdWV1ZVJlcXVlc3QuQXR0cmlidXRlc0VudHJ5UgphdHRyaWJ1dGVzEh8KCXF1ZXVlbmFtZRjc'
    '1c1JIAEoCVIJcXVldWVuYW1lEjkKBHRhZ3MYodfboAEgAygLMiEuc3FzLkNyZWF0ZVF1ZXVlUm'
    'VxdWVzdC5UYWdzRW50cnlSBHRhZ3MaPQoPQXR0cmlidXRlc0VudHJ5EhAKA2tleRgBIAEoCVID'
    'a2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaNwoJVGFnc0VudHJ5EhAKA2tleRgBIAEoCV'
    'IDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use createQueueResultDescriptor instead')
const CreateQueueResult$json = {
  '1': 'CreateQueueResult',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `CreateQueueResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createQueueResultDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVRdWV1ZVJlc3VsdBIeCghxdWV1ZXVybBjKwb7zASABKAlSCHF1ZXVldXJs');

@$core.Deprecated('Use deleteMessageBatchRequestDescriptor instead')
const DeleteMessageBatchRequest$json = {
  '1': 'DeleteMessageBatchRequest',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.sqs.DeleteMessageBatchRequestEntry',
      '10': 'entries'
    },
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `DeleteMessageBatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMessageBatchRequestDescriptor = $convert.base64Decode(
    'ChlEZWxldGVNZXNzYWdlQmF0Y2hSZXF1ZXN0EkEKB2VudHJpZXMYlMWy5QEgAygLMiMuc3FzLk'
    'RlbGV0ZU1lc3NhZ2VCYXRjaFJlcXVlc3RFbnRyeVIHZW50cmllcxIeCghxdWV1ZXVybBjKwb7z'
    'ASABKAlSCHF1ZXVldXJs');

@$core.Deprecated('Use deleteMessageBatchRequestEntryDescriptor instead')
const DeleteMessageBatchRequestEntry$json = {
  '1': 'DeleteMessageBatchRequestEntry',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'receipthandle',
      '3': 134471750,
      '4': 1,
      '5': 9,
      '10': 'receipthandle'
    },
  ],
};

/// Descriptor for `DeleteMessageBatchRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMessageBatchRequestEntryDescriptor =
    $convert.base64Decode(
        'Ch5EZWxldGVNZXNzYWdlQmF0Y2hSZXF1ZXN0RW50cnkSEgoCaWQYgfKitwEgASgJUgJpZBInCg'
        '1yZWNlaXB0aGFuZGxlGMbAj0AgASgJUg1yZWNlaXB0aGFuZGxl');

@$core.Deprecated('Use deleteMessageBatchResultDescriptor instead')
const DeleteMessageBatchResult$json = {
  '1': 'DeleteMessageBatchResult',
  '2': [
    {
      '1': 'failed',
      '3': 360301525,
      '4': 3,
      '5': 11,
      '6': '.sqs.BatchResultErrorEntry',
      '10': 'failed'
    },
    {
      '1': 'successful',
      '3': 412818844,
      '4': 3,
      '5': 11,
      '6': '.sqs.DeleteMessageBatchResultEntry',
      '10': 'successful'
    },
  ],
};

/// Descriptor for `DeleteMessageBatchResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMessageBatchResultDescriptor = $convert.base64Decode(
    'ChhEZWxldGVNZXNzYWdlQmF0Y2hSZXN1bHQSNgoGZmFpbGVkGNWH56sBIAMoCzIaLnNxcy5CYX'
    'RjaFJlc3VsdEVycm9yRW50cnlSBmZhaWxlZBJGCgpzdWNjZXNzZnVsGJy77MQBIAMoCzIiLnNx'
    'cy5EZWxldGVNZXNzYWdlQmF0Y2hSZXN1bHRFbnRyeVIKc3VjY2Vzc2Z1bA==');

@$core.Deprecated('Use deleteMessageBatchResultEntryDescriptor instead')
const DeleteMessageBatchResultEntry$json = {
  '1': 'DeleteMessageBatchResultEntry',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `DeleteMessageBatchResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMessageBatchResultEntryDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVNZXNzYWdlQmF0Y2hSZXN1bHRFbnRyeRISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use deleteMessageRequestDescriptor instead')
const DeleteMessageRequest$json = {
  '1': 'DeleteMessageRequest',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
    {
      '1': 'receipthandle',
      '3': 134471750,
      '4': 1,
      '5': 9,
      '10': 'receipthandle'
    },
  ],
};

/// Descriptor for `DeleteMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMessageRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVNZXNzYWdlUmVxdWVzdBIeCghxdWV1ZXVybBjKwb7zASABKAlSCHF1ZXVldXJsEi'
    'cKDXJlY2VpcHRoYW5kbGUYxsCPQCABKAlSDXJlY2VpcHRoYW5kbGU=');

@$core.Deprecated('Use deleteQueueRequestDescriptor instead')
const DeleteQueueRequest$json = {
  '1': 'DeleteQueueRequest',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `DeleteQueueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteQueueRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVRdWV1ZVJlcXVlc3QSHgoIcXVldWV1cmwYysG+8wEgASgJUghxdWV1ZXVybA==');

@$core.Deprecated('Use emptyBatchRequestDescriptor instead')
const EmptyBatchRequest$json = {
  '1': 'EmptyBatchRequest',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `EmptyBatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emptyBatchRequestDescriptor = $convert.base64Decode(
    'ChFFbXB0eUJhdGNoUmVxdWVzdBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use getQueueAttributesRequestDescriptor instead')
const GetQueueAttributesRequest$json = {
  '1': 'GetQueueAttributesRequest',
  '2': [
    {
      '1': 'attributenames',
      '3': 394468622,
      '4': 3,
      '5': 14,
      '6': '.sqs.QueueAttributeName',
      '10': 'attributenames'
    },
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `GetQueueAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueueAttributesRequestDescriptor = $convert.base64Decode(
    'ChlHZXRRdWV1ZUF0dHJpYnV0ZXNSZXF1ZXN0EkMKDmF0dHJpYnV0ZW5hbWVzGI66jLwBIAMoDj'
    'IXLnNxcy5RdWV1ZUF0dHJpYnV0ZU5hbWVSDmF0dHJpYnV0ZW5hbWVzEh4KCHF1ZXVldXJsGMrB'
    'vvMBIAEoCVIIcXVldWV1cmw=');

@$core.Deprecated('Use getQueueAttributesResultDescriptor instead')
const GetQueueAttributesResult$json = {
  '1': 'GetQueueAttributesResult',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sqs.GetQueueAttributesResult.AttributesEntry',
      '10': 'attributes'
    },
  ],
  '3': [GetQueueAttributesResult_AttributesEntry$json],
};

@$core.Deprecated('Use getQueueAttributesResultDescriptor instead')
const GetQueueAttributesResult_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetQueueAttributesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueueAttributesResultDescriptor = $convert.base64Decode(
    'ChhHZXRRdWV1ZUF0dHJpYnV0ZXNSZXN1bHQSUAoKYXR0cmlidXRlcxi1qftjIAMoCzItLnNxcy'
    '5HZXRRdWV1ZUF0dHJpYnV0ZXNSZXN1bHQuQXR0cmlidXRlc0VudHJ5UgphdHRyaWJ1dGVzGj0K'
    'D0F0dHJpYnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdW'
    'U6AjgB');

@$core.Deprecated('Use getQueueUrlRequestDescriptor instead')
const GetQueueUrlRequest$json = {
  '1': 'GetQueueUrlRequest',
  '2': [
    {'1': 'queuename', '3': 154364636, '4': 1, '5': 9, '10': 'queuename'},
    {
      '1': 'queueownerawsaccountid',
      '3': 248813375,
      '4': 1,
      '5': 9,
      '10': 'queueownerawsaccountid'
    },
  ],
};

/// Descriptor for `GetQueueUrlRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueueUrlRequestDescriptor = $convert.base64Decode(
    'ChJHZXRRdWV1ZVVybFJlcXVlc3QSHwoJcXVldWVuYW1lGNzVzUkgASgJUglxdWV1ZW5hbWUSOQ'
    'oWcXVldWVvd25lcmF3c2FjY291bnRpZBi/rtJ2IAEoCVIWcXVldWVvd25lcmF3c2FjY291bnRp'
    'ZA==');

@$core.Deprecated('Use getQueueUrlResultDescriptor instead')
const GetQueueUrlResult$json = {
  '1': 'GetQueueUrlResult',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `GetQueueUrlResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getQueueUrlResultDescriptor = $convert.base64Decode(
    'ChFHZXRRdWV1ZVVybFJlc3VsdBIeCghxdWV1ZXVybBjKwb7zASABKAlSCHF1ZXVldXJs');

@$core.Deprecated('Use invalidAddressDescriptor instead')
const InvalidAddress$json = {
  '1': 'InvalidAddress',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAddress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAddressDescriptor = $convert.base64Decode(
    'Cg5JbnZhbGlkQWRkcmVzcxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidAttributeNameDescriptor instead')
const InvalidAttributeName$json = {
  '1': 'InvalidAttributeName',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAttributeName`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAttributeNameDescriptor =
    $convert.base64Decode(
        'ChRJbnZhbGlkQXR0cmlidXRlTmFtZRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidAttributeValueDescriptor instead')
const InvalidAttributeValue$json = {
  '1': 'InvalidAttributeValue',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAttributeValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAttributeValueDescriptor = $convert.base64Decode(
    'ChVJbnZhbGlkQXR0cmlidXRlVmFsdWUSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use invalidBatchEntryIdDescriptor instead')
const InvalidBatchEntryId$json = {
  '1': 'InvalidBatchEntryId',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidBatchEntryId`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidBatchEntryIdDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkQmF0Y2hFbnRyeUlkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidIdFormatDescriptor instead')
const InvalidIdFormat$json = {
  '1': 'InvalidIdFormat',
};

/// Descriptor for `InvalidIdFormat`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidIdFormatDescriptor =
    $convert.base64Decode('Cg9JbnZhbGlkSWRGb3JtYXQ=');

@$core.Deprecated('Use invalidMessageContentsDescriptor instead')
const InvalidMessageContents$json = {
  '1': 'InvalidMessageContents',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidMessageContents`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidMessageContentsDescriptor =
    $convert.base64Decode(
        'ChZJbnZhbGlkTWVzc2FnZUNvbnRlbnRzEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidSecurityDescriptor instead')
const InvalidSecurity$json = {
  '1': 'InvalidSecurity',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidSecurity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidSecurityDescriptor = $convert.base64Decode(
    'Cg9JbnZhbGlkU2VjdXJpdHkSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use kmsAccessDeniedDescriptor instead')
const KmsAccessDenied$json = {
  '1': 'KmsAccessDenied',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsAccessDenied`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsAccessDeniedDescriptor = $convert.base64Decode(
    'Cg9LbXNBY2Nlc3NEZW5pZWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use kmsDisabledDescriptor instead')
const KmsDisabled$json = {
  '1': 'KmsDisabled',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsDisabled`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsDisabledDescriptor = $convert
    .base64Decode('CgtLbXNEaXNhYmxlZBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use kmsInvalidKeyUsageDescriptor instead')
const KmsInvalidKeyUsage$json = {
  '1': 'KmsInvalidKeyUsage',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsInvalidKeyUsage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsInvalidKeyUsageDescriptor =
    $convert.base64Decode(
        'ChJLbXNJbnZhbGlkS2V5VXNhZ2USGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use kmsInvalidStateDescriptor instead')
const KmsInvalidState$json = {
  '1': 'KmsInvalidState',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsInvalidState`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsInvalidStateDescriptor = $convert.base64Decode(
    'Cg9LbXNJbnZhbGlkU3RhdGUSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use kmsNotFoundDescriptor instead')
const KmsNotFound$json = {
  '1': 'KmsNotFound',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsNotFound`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsNotFoundDescriptor = $convert
    .base64Decode('CgtLbXNOb3RGb3VuZBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use kmsOptInRequiredDescriptor instead')
const KmsOptInRequired$json = {
  '1': 'KmsOptInRequired',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsOptInRequired`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsOptInRequiredDescriptor = $convert.base64Decode(
    'ChBLbXNPcHRJblJlcXVpcmVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use kmsThrottledDescriptor instead')
const KmsThrottled$json = {
  '1': 'KmsThrottled',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KmsThrottled`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kmsThrottledDescriptor = $convert.base64Decode(
    'CgxLbXNUaHJvdHRsZWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use listDeadLetterSourceQueuesRequestDescriptor instead')
const ListDeadLetterSourceQueuesRequest$json = {
  '1': 'ListDeadLetterSourceQueuesRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `ListDeadLetterSourceQueuesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDeadLetterSourceQueuesRequestDescriptor =
    $convert.base64Decode(
        'CiFMaXN0RGVhZExldHRlclNvdXJjZVF1ZXVlc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDAS'
        'ABKAVSCm1heHJlc3VsdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SHgoIcXVl'
        'dWV1cmwYysG+8wEgASgJUghxdWV1ZXVybA==');

@$core.Deprecated('Use listDeadLetterSourceQueuesResultDescriptor instead')
const ListDeadLetterSourceQueuesResult$json = {
  '1': 'ListDeadLetterSourceQueuesResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'queueurls', '3': 475536111, '4': 3, '5': 9, '10': 'queueurls'},
  ],
};

/// Descriptor for `ListDeadLetterSourceQueuesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDeadLetterSourceQueuesResultDescriptor =
    $convert.base64Decode(
        'CiBMaXN0RGVhZExldHRlclNvdXJjZVF1ZXVlc1Jlc3VsdBIfCgluZXh0dG9rZW4Y/oS6ZyABKA'
        'lSCW5leHR0b2tlbhIgCglxdWV1ZXVybHMY77Xg4gEgAygJUglxdWV1ZXVybHM=');

@$core.Deprecated('Use listMessageMoveTasksRequestDescriptor instead')
const ListMessageMoveTasksRequest$json = {
  '1': 'ListMessageMoveTasksRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'sourcearn', '3': 439903072, '4': 1, '5': 9, '10': 'sourcearn'},
  ],
};

/// Descriptor for `ListMessageMoveTasksRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMessageMoveTasksRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0TWVzc2FnZU1vdmVUYXNrc1JlcXVlc3QSIgoKbWF4cmVzdWx0cxiyqJuDASABKAVSCm'
        '1heHJlc3VsdHMSIAoJc291cmNlYXJuGODG4dEBIAEoCVIJc291cmNlYXJu');

@$core.Deprecated('Use listMessageMoveTasksResultDescriptor instead')
const ListMessageMoveTasksResult$json = {
  '1': 'ListMessageMoveTasksResult',
  '2': [
    {
      '1': 'results',
      '3': 486024854,
      '4': 3,
      '5': 11,
      '6': '.sqs.ListMessageMoveTasksResultEntry',
      '10': 'results'
    },
  ],
};

/// Descriptor for `ListMessageMoveTasksResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMessageMoveTasksResultDescriptor =
    $convert.base64Decode(
        'ChpMaXN0TWVzc2FnZU1vdmVUYXNrc1Jlc3VsdBJCCgdyZXN1bHRzGJbN4OcBIAMoCzIkLnNxcy'
        '5MaXN0TWVzc2FnZU1vdmVUYXNrc1Jlc3VsdEVudHJ5UgdyZXN1bHRz');

@$core.Deprecated('Use listMessageMoveTasksResultEntryDescriptor instead')
const ListMessageMoveTasksResultEntry$json = {
  '1': 'ListMessageMoveTasksResultEntry',
  '2': [
    {
      '1': 'approximatenumberofmessagesmoved',
      '3': 510181657,
      '4': 1,
      '5': 3,
      '10': 'approximatenumberofmessagesmoved'
    },
    {
      '1': 'approximatenumberofmessagestomove',
      '3': 108945690,
      '4': 1,
      '5': 3,
      '10': 'approximatenumberofmessagestomove'
    },
    {
      '1': 'destinationarn',
      '3': 375726595,
      '4': 1,
      '5': 9,
      '10': 'destinationarn'
    },
    {
      '1': 'failurereason',
      '3': 232322142,
      '4': 1,
      '5': 9,
      '10': 'failurereason'
    },
    {
      '1': 'maxnumberofmessagespersecond',
      '3': 335779921,
      '4': 1,
      '5': 5,
      '10': 'maxnumberofmessagespersecond'
    },
    {'1': 'sourcearn', '3': 439903072, '4': 1, '5': 9, '10': 'sourcearn'},
    {
      '1': 'startedtimestamp',
      '3': 397447975,
      '4': 1,
      '5': 3,
      '10': 'startedtimestamp'
    },
    {'1': 'status', '3': 6222352, '4': 1, '5': 9, '10': 'status'},
    {'1': 'taskhandle', '3': 190544291, '4': 1, '5': 9, '10': 'taskhandle'},
  ],
};

/// Descriptor for `ListMessageMoveTasksResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMessageMoveTasksResultEntryDescriptor = $convert.base64Decode(
    'Ch9MaXN0TWVzc2FnZU1vdmVUYXNrc1Jlc3VsdEVudHJ5Ek4KIGFwcHJveGltYXRlbnVtYmVyb2'
    'ZtZXNzYWdlc21vdmVkGJmCo/MBIAEoA1IgYXBwcm94aW1hdGVudW1iZXJvZm1lc3NhZ2VzbW92'
    'ZWQSTwohYXBwcm94aW1hdGVudW1iZXJvZm1lc3NhZ2VzdG9tb3ZlGJrC+TMgASgDUiFhcHByb3'
    'hpbWF0ZW51bWJlcm9mbWVzc2FnZXN0b21vdmUSKgoOZGVzdGluYXRpb25hcm4Yg8SUswEgASgJ'
    'Ug5kZXN0aW5hdGlvbmFybhInCg1mYWlsdXJlcmVhc29uGN7o424gASgJUg1mYWlsdXJlcmVhc2'
    '9uEkYKHG1heG51bWJlcm9mbWVzc2FnZXNwZXJzZWNvbmQY0bCOoAEgASgFUhxtYXhudW1iZXJv'
    'Zm1lc3NhZ2VzcGVyc2Vjb25kEiAKCXNvdXJjZWFybhjgxuHRASABKAlSCXNvdXJjZWFybhIuCh'
    'BzdGFydGVkdGltZXN0YW1wGKemwr0BIAEoA1IQc3RhcnRlZHRpbWVzdGFtcBIZCgZzdGF0dXMY'
    'kOT7AiABKAlSBnN0YXR1cxIhCgp0YXNraGFuZGxlGKPz7VogASgJUgp0YXNraGFuZGxl');

@$core.Deprecated('Use listQueueTagsRequestDescriptor instead')
const ListQueueTagsRequest$json = {
  '1': 'ListQueueTagsRequest',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `ListQueueTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueueTagsRequestDescriptor = $convert.base64Decode(
    'ChRMaXN0UXVldWVUYWdzUmVxdWVzdBIeCghxdWV1ZXVybBjKwb7zASABKAlSCHF1ZXVldXJs');

@$core.Deprecated('Use listQueueTagsResultDescriptor instead')
const ListQueueTagsResult$json = {
  '1': 'ListQueueTagsResult',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sqs.ListQueueTagsResult.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [ListQueueTagsResult_TagsEntry$json],
};

@$core.Deprecated('Use listQueueTagsResultDescriptor instead')
const ListQueueTagsResult_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListQueueTagsResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueueTagsResultDescriptor = $convert.base64Decode(
    'ChNMaXN0UXVldWVUYWdzUmVzdWx0EjoKBHRhZ3MYwcH2tQEgAygLMiIuc3FzLkxpc3RRdWV1ZV'
    'RhZ3NSZXN1bHQuVGFnc0VudHJ5UgR0YWdzGjcKCVRhZ3NFbnRyeRIQCgNrZXkYASABKAlSA2tl'
    'eRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use listQueuesRequestDescriptor instead')
const ListQueuesRequest$json = {
  '1': 'ListQueuesRequest',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'queuenameprefix',
      '3': 269478416,
      '4': 1,
      '5': 9,
      '10': 'queuenameprefix'
    },
  ],
};

/// Descriptor for `ListQueuesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueuesRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0UXVldWVzUmVxdWVzdBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cmVzdWx0cx'
    'IfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIsCg9xdWV1ZW5hbWVwcmVmaXgYkNS/'
    'gAEgASgJUg9xdWV1ZW5hbWVwcmVmaXg=');

@$core.Deprecated('Use listQueuesResultDescriptor instead')
const ListQueuesResult$json = {
  '1': 'ListQueuesResult',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'queueurls', '3': 62522575, '4': 3, '5': 9, '10': 'queueurls'},
  ],
};

/// Descriptor for `ListQueuesResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listQueuesResultDescriptor = $convert.base64Decode(
    'ChBMaXN0UXVldWVzUmVzdWx0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEh8KCX'
    'F1ZXVldXJscxjPiegdIAMoCVIJcXVldWV1cmxz');

@$core.Deprecated('Use messageDescriptor instead')
const Message$json = {
  '1': 'Message',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sqs.Message.AttributesEntry',
      '10': 'attributes'
    },
    {'1': 'body', '3': 42602646, '4': 1, '5': 9, '10': 'body'},
    {'1': 'md5ofbody', '3': 229335371, '4': 1, '5': 9, '10': 'md5ofbody'},
    {
      '1': 'md5ofmessageattributes',
      '3': 499647077,
      '4': 1,
      '5': 9,
      '10': 'md5ofmessageattributes'
    },
    {
      '1': 'messageattributes',
      '3': 56443766,
      '4': 3,
      '5': 11,
      '6': '.sqs.Message.MessageattributesEntry',
      '10': 'messageattributes'
    },
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
    {
      '1': 'receipthandle',
      '3': 134471750,
      '4': 1,
      '5': 9,
      '10': 'receipthandle'
    },
  ],
  '3': [Message_AttributesEntry$json, Message_MessageattributesEntry$json],
};

@$core.Deprecated('Use messageDescriptor instead')
const Message_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use messageDescriptor instead')
const Message_MessageattributesEntry$json = {
  '1': 'MessageattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.sqs.MessageAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `Message`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageDescriptor = $convert.base64Decode(
    'CgdNZXNzYWdlEj8KCmF0dHJpYnV0ZXMYtan7YyADKAsyHC5zcXMuTWVzc2FnZS5BdHRyaWJ1dG'
    'VzRW50cnlSCmF0dHJpYnV0ZXMSFQoEYm9keRiWoagUIAEoCVIEYm9keRIfCgltZDVvZmJvZHkY'
    'y8KtbSABKAlSCW1kNW9mYm9keRI6ChZtZDVvZm1lc3NhZ2VhdHRyaWJ1dGVzGOWEoO4BIAEoCV'
    'IWbWQ1b2ZtZXNzYWdlYXR0cmlidXRlcxJUChFtZXNzYWdlYXR0cmlidXRlcxj2hvUaIAMoCzIj'
    'LnNxcy5NZXNzYWdlLk1lc3NhZ2VhdHRyaWJ1dGVzRW50cnlSEW1lc3NhZ2VhdHRyaWJ1dGVzEi'
    'AKCW1lc3NhZ2VpZBiq5vSrASABKAlSCW1lc3NhZ2VpZBInCg1yZWNlaXB0aGFuZGxlGMbAj0Ag'
    'ASgJUg1yZWNlaXB0aGFuZGxlGj0KD0F0dHJpYnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleR'
    'IUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgBGmAKFk1lc3NhZ2VhdHRyaWJ1dGVzRW50cnkSEAoD'
    'a2V5GAEgASgJUgNrZXkSMAoFdmFsdWUYAiABKAsyGi5zcXMuTWVzc2FnZUF0dHJpYnV0ZVZhbH'
    'VlUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use messageAttributeValueDescriptor instead')
const MessageAttributeValue$json = {
  '1': 'MessageAttributeValue',
  '2': [
    {
      '1': 'binarylistvalues',
      '3': 158345259,
      '4': 3,
      '5': 12,
      '10': 'binarylistvalues'
    },
    {'1': 'binaryvalue', '3': 255476278, '4': 1, '5': 12, '10': 'binaryvalue'},
    {'1': 'datatype', '3': 67988590, '4': 1, '5': 9, '10': 'datatype'},
    {
      '1': 'stringlistvalues',
      '3': 86527527,
      '4': 3,
      '5': 9,
      '10': 'stringlistvalues'
    },
    {'1': 'stringvalue', '3': 184416138, '4': 1, '5': 9, '10': 'stringvalue'},
  ],
};

/// Descriptor for `MessageAttributeValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageAttributeValueDescriptor = $convert.base64Decode(
    'ChVNZXNzYWdlQXR0cmlidXRlVmFsdWUSLQoQYmluYXJ5bGlzdHZhbHVlcxir0MBLIAMoDFIQYm'
    'luYXJ5bGlzdHZhbHVlcxIjCgtiaW5hcnl2YWx1ZRi2hOl5IAEoDFILYmluYXJ5dmFsdWUSHQoI'
    'ZGF0YXR5cGUY7ti1ICABKAlSCGRhdGF0eXBlEi0KEHN0cmluZ2xpc3R2YWx1ZXMYp5yhKSADKA'
    'lSEHN0cmluZ2xpc3R2YWx1ZXMSIwoLc3RyaW5ndmFsdWUYiu/3VyABKAlSC3N0cmluZ3ZhbHVl');

@$core.Deprecated('Use messageNotInflightDescriptor instead')
const MessageNotInflight$json = {
  '1': 'MessageNotInflight',
};

/// Descriptor for `MessageNotInflight`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageNotInflightDescriptor =
    $convert.base64Decode('ChJNZXNzYWdlTm90SW5mbGlnaHQ=');

@$core.Deprecated('Use messageSystemAttributeValueDescriptor instead')
const MessageSystemAttributeValue$json = {
  '1': 'MessageSystemAttributeValue',
  '2': [
    {
      '1': 'binarylistvalues',
      '3': 158345259,
      '4': 3,
      '5': 12,
      '10': 'binarylistvalues'
    },
    {'1': 'binaryvalue', '3': 255476278, '4': 1, '5': 12, '10': 'binaryvalue'},
    {'1': 'datatype', '3': 67988590, '4': 1, '5': 9, '10': 'datatype'},
    {
      '1': 'stringlistvalues',
      '3': 86527527,
      '4': 3,
      '5': 9,
      '10': 'stringlistvalues'
    },
    {'1': 'stringvalue', '3': 184416138, '4': 1, '5': 9, '10': 'stringvalue'},
  ],
};

/// Descriptor for `MessageSystemAttributeValue`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageSystemAttributeValueDescriptor = $convert.base64Decode(
    'ChtNZXNzYWdlU3lzdGVtQXR0cmlidXRlVmFsdWUSLQoQYmluYXJ5bGlzdHZhbHVlcxir0MBLIA'
    'MoDFIQYmluYXJ5bGlzdHZhbHVlcxIjCgtiaW5hcnl2YWx1ZRi2hOl5IAEoDFILYmluYXJ5dmFs'
    'dWUSHQoIZGF0YXR5cGUY7ti1ICABKAlSCGRhdGF0eXBlEi0KEHN0cmluZ2xpc3R2YWx1ZXMYp5'
    'yhKSADKAlSEHN0cmluZ2xpc3R2YWx1ZXMSIwoLc3RyaW5ndmFsdWUYiu/3VyABKAlSC3N0cmlu'
    'Z3ZhbHVl');

@$core.Deprecated('Use overLimitDescriptor instead')
const OverLimit$json = {
  '1': 'OverLimit',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `OverLimit`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List overLimitDescriptor = $convert
    .base64Decode('CglPdmVyTGltaXQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use purgeQueueInProgressDescriptor instead')
const PurgeQueueInProgress$json = {
  '1': 'PurgeQueueInProgress',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `PurgeQueueInProgress`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List purgeQueueInProgressDescriptor =
    $convert.base64Decode(
        'ChRQdXJnZVF1ZXVlSW5Qcm9ncmVzcxIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use purgeQueueRequestDescriptor instead')
const PurgeQueueRequest$json = {
  '1': 'PurgeQueueRequest',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `PurgeQueueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List purgeQueueRequestDescriptor = $convert.base64Decode(
    'ChFQdXJnZVF1ZXVlUmVxdWVzdBIeCghxdWV1ZXVybBjKwb7zASABKAlSCHF1ZXVldXJs');

@$core.Deprecated('Use queueDeletedRecentlyDescriptor instead')
const QueueDeletedRecently$json = {
  '1': 'QueueDeletedRecently',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueueDeletedRecently`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queueDeletedRecentlyDescriptor =
    $convert.base64Decode(
        'ChRRdWV1ZURlbGV0ZWRSZWNlbnRseRIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use queueDoesNotExistDescriptor instead')
const QueueDoesNotExist$json = {
  '1': 'QueueDoesNotExist',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueueDoesNotExist`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queueDoesNotExistDescriptor = $convert.base64Decode(
    'ChFRdWV1ZURvZXNOb3RFeGlzdBIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use queueNameExistsDescriptor instead')
const QueueNameExists$json = {
  '1': 'QueueNameExists',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `QueueNameExists`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List queueNameExistsDescriptor = $convert.base64Decode(
    'Cg9RdWV1ZU5hbWVFeGlzdHMSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use receiptHandleIsInvalidDescriptor instead')
const ReceiptHandleIsInvalid$json = {
  '1': 'ReceiptHandleIsInvalid',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ReceiptHandleIsInvalid`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List receiptHandleIsInvalidDescriptor =
    $convert.base64Decode(
        'ChZSZWNlaXB0SGFuZGxlSXNJbnZhbGlkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use receiveMessageRequestDescriptor instead')
const ReceiveMessageRequest$json = {
  '1': 'ReceiveMessageRequest',
  '2': [
    {
      '1': 'attributenames',
      '3': 394468622,
      '4': 3,
      '5': 14,
      '6': '.sqs.QueueAttributeName',
      '10': 'attributenames'
    },
    {
      '1': 'maxnumberofmessages',
      '3': 139866822,
      '4': 1,
      '5': 5,
      '10': 'maxnumberofmessages'
    },
    {
      '1': 'messageattributenames',
      '3': 332558373,
      '4': 3,
      '5': 9,
      '10': 'messageattributenames'
    },
    {
      '1': 'messagesystemattributenames',
      '3': 42109014,
      '4': 3,
      '5': 14,
      '6': '.sqs.MessageSystemAttributeName',
      '10': 'messagesystemattributenames'
    },
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
    {
      '1': 'receiverequestattemptid',
      '3': 455135954,
      '4': 1,
      '5': 9,
      '10': 'receiverequestattemptid'
    },
    {
      '1': 'visibilitytimeout',
      '3': 460820073,
      '4': 1,
      '5': 5,
      '10': 'visibilitytimeout'
    },
    {
      '1': 'waittimeseconds',
      '3': 398991863,
      '4': 1,
      '5': 5,
      '10': 'waittimeseconds'
    },
  ],
};

/// Descriptor for `ReceiveMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List receiveMessageRequestDescriptor = $convert.base64Decode(
    'ChVSZWNlaXZlTWVzc2FnZVJlcXVlc3QSQwoOYXR0cmlidXRlbmFtZXMYjrqMvAEgAygOMhcuc3'
    'FzLlF1ZXVlQXR0cmlidXRlTmFtZVIOYXR0cmlidXRlbmFtZXMSMwoTbWF4bnVtYmVyb2ZtZXNz'
    'YWdlcxjG5dhCIAEoBVITbWF4bnVtYmVyb2ZtZXNzYWdlcxI4ChVtZXNzYWdlYXR0cmlidXRlbm'
    'FtZXMYpeDJngEgAygJUhVtZXNzYWdlYXR0cmlidXRlbmFtZXMSZAobbWVzc2FnZXN5c3RlbWF0'
    'dHJpYnV0ZW5hbWVzGNaQihQgAygOMh8uc3FzLk1lc3NhZ2VTeXN0ZW1BdHRyaWJ1dGVOYW1lUh'
    'ttZXNzYWdlc3lzdGVtYXR0cmlidXRlbmFtZXMSHgoIcXVldWV1cmwYysG+8wEgASgJUghxdWV1'
    'ZXVybBI8ChdyZWNlaXZlcmVxdWVzdGF0dGVtcHRpZBjSpYPZASABKAlSF3JlY2VpdmVyZXF1ZX'
    'N0YXR0ZW1wdGlkEjAKEXZpc2liaWxpdHl0aW1lb3V0GOmc3tsBIAEoBVIRdmlzaWJpbGl0eXRp'
    'bWVvdXQSLAoPd2FpdHRpbWVzZWNvbmRzGPfDoL4BIAEoBVIPd2FpdHRpbWVzZWNvbmRz');

@$core.Deprecated('Use receiveMessageResultDescriptor instead')
const ReceiveMessageResult$json = {
  '1': 'ReceiveMessageResult',
  '2': [
    {
      '1': 'messages',
      '3': 409018326,
      '4': 3,
      '5': 11,
      '6': '.sqs.Message',
      '10': 'messages'
    },
  ],
};

/// Descriptor for `ReceiveMessageResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List receiveMessageResultDescriptor = $convert.base64Decode(
    'ChRSZWNlaXZlTWVzc2FnZVJlc3VsdBIsCghtZXNzYWdlcxjWv4TDASADKAsyDC5zcXMuTWVzc2'
    'FnZVIIbWVzc2FnZXM=');

@$core.Deprecated('Use removePermissionRequestDescriptor instead')
const RemovePermissionRequest$json = {
  '1': 'RemovePermissionRequest',
  '2': [
    {'1': 'label', '3': 516747934, '4': 1, '5': 9, '10': 'label'},
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `RemovePermissionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List removePermissionRequestDescriptor =
    $convert.base64Decode(
        'ChdSZW1vdmVQZXJtaXNzaW9uUmVxdWVzdBIYCgVsYWJlbBie5bP2ASABKAlSBWxhYmVsEh4KCH'
        'F1ZXVldXJsGMrBvvMBIAEoCVIIcXVldWV1cmw=');

@$core.Deprecated('Use requestThrottledDescriptor instead')
const RequestThrottled$json = {
  '1': 'RequestThrottled',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `RequestThrottled`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List requestThrottledDescriptor = $convert.base64Decode(
    'ChBSZXF1ZXN0VGhyb3R0bGVkEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

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

@$core.Deprecated('Use sendMessageBatchRequestDescriptor instead')
const SendMessageBatchRequest$json = {
  '1': 'SendMessageBatchRequest',
  '2': [
    {
      '1': 'entries',
      '3': 481075860,
      '4': 3,
      '5': 11,
      '6': '.sqs.SendMessageBatchRequestEntry',
      '10': 'entries'
    },
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
};

/// Descriptor for `SendMessageBatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageBatchRequestDescriptor = $convert.base64Decode(
    'ChdTZW5kTWVzc2FnZUJhdGNoUmVxdWVzdBI/CgdlbnRyaWVzGJTFsuUBIAMoCzIhLnNxcy5TZW'
    '5kTWVzc2FnZUJhdGNoUmVxdWVzdEVudHJ5UgdlbnRyaWVzEh4KCHF1ZXVldXJsGMrBvvMBIAEo'
    'CVIIcXVldWV1cmw=');

@$core.Deprecated('Use sendMessageBatchRequestEntryDescriptor instead')
const SendMessageBatchRequestEntry$json = {
  '1': 'SendMessageBatchRequestEntry',
  '2': [
    {'1': 'delayseconds', '3': 48268198, '4': 1, '5': 5, '10': 'delayseconds'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'messageattributes',
      '3': 56443766,
      '4': 3,
      '5': 11,
      '6': '.sqs.SendMessageBatchRequestEntry.MessageattributesEntry',
      '10': 'messageattributes'
    },
    {'1': 'messagebody', '3': 56920001, '4': 1, '5': 9, '10': 'messagebody'},
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
      '1': 'messagesystemattributes',
      '3': 194951997,
      '4': 3,
      '5': 11,
      '6': '.sqs.SendMessageBatchRequestEntry.MessagesystemattributesEntry',
      '10': 'messagesystemattributes'
    },
  ],
  '3': [
    SendMessageBatchRequestEntry_MessageattributesEntry$json,
    SendMessageBatchRequestEntry_MessagesystemattributesEntry$json
  ],
};

@$core.Deprecated('Use sendMessageBatchRequestEntryDescriptor instead')
const SendMessageBatchRequestEntry_MessageattributesEntry$json = {
  '1': 'MessageattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.sqs.MessageAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use sendMessageBatchRequestEntryDescriptor instead')
const SendMessageBatchRequestEntry_MessagesystemattributesEntry$json = {
  '1': 'MessagesystemattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.sqs.MessageSystemAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `SendMessageBatchRequestEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageBatchRequestEntryDescriptor = $convert.base64Decode(
    'ChxTZW5kTWVzc2FnZUJhdGNoUmVxdWVzdEVudHJ5EiUKDGRlbGF5c2Vjb25kcximh4IXIAEoBV'
    'IMZGVsYXlzZWNvbmRzEhIKAmlkGIHyorcBIAEoCVICaWQSaQoRbWVzc2FnZWF0dHJpYnV0ZXMY'
    '9ob1GiADKAsyOC5zcXMuU2VuZE1lc3NhZ2VCYXRjaFJlcXVlc3RFbnRyeS5NZXNzYWdlYXR0cm'
    'lidXRlc0VudHJ5UhFtZXNzYWdlYXR0cmlidXRlcxIjCgttZXNzYWdlYm9keRjBj5IbIAEoCVIL'
    'bWVzc2FnZWJvZHkSOgoWbWVzc2FnZWRlZHVwbGljYXRpb25pZBjZxf60ASABKAlSFm1lc3NhZ2'
    'VkZWR1cGxpY2F0aW9uaWQSKgoObWVzc2FnZWdyb3VwaWQYm8SGyAEgASgJUg5tZXNzYWdlZ3Jv'
    'dXBpZBJ7ChdtZXNzYWdlc3lzdGVtYXR0cmlidXRlcxi99vpcIAMoCzI+LnNxcy5TZW5kTWVzc2'
    'FnZUJhdGNoUmVxdWVzdEVudHJ5Lk1lc3NhZ2VzeXN0ZW1hdHRyaWJ1dGVzRW50cnlSF21lc3Nh'
    'Z2VzeXN0ZW1hdHRyaWJ1dGVzGmAKFk1lc3NhZ2VhdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgAS'
    'gJUgNrZXkSMAoFdmFsdWUYAiABKAsyGi5zcXMuTWVzc2FnZUF0dHJpYnV0ZVZhbHVlUgV2YWx1'
    'ZToCOAEabAocTWVzc2FnZXN5c3RlbWF0dHJpYnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleR'
    'I2CgV2YWx1ZRgCIAEoCzIgLnNxcy5NZXNzYWdlU3lzdGVtQXR0cmlidXRlVmFsdWVSBXZhbHVl'
    'OgI4AQ==');

@$core.Deprecated('Use sendMessageBatchResultDescriptor instead')
const SendMessageBatchResult$json = {
  '1': 'SendMessageBatchResult',
  '2': [
    {
      '1': 'failed',
      '3': 360301525,
      '4': 3,
      '5': 11,
      '6': '.sqs.BatchResultErrorEntry',
      '10': 'failed'
    },
    {
      '1': 'successful',
      '3': 412818844,
      '4': 3,
      '5': 11,
      '6': '.sqs.SendMessageBatchResultEntry',
      '10': 'successful'
    },
  ],
};

/// Descriptor for `SendMessageBatchResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageBatchResultDescriptor = $convert.base64Decode(
    'ChZTZW5kTWVzc2FnZUJhdGNoUmVzdWx0EjYKBmZhaWxlZBjVh+erASADKAsyGi5zcXMuQmF0Y2'
    'hSZXN1bHRFcnJvckVudHJ5UgZmYWlsZWQSRAoKc3VjY2Vzc2Z1bBicu+zEASADKAsyIC5zcXMu'
    'U2VuZE1lc3NhZ2VCYXRjaFJlc3VsdEVudHJ5UgpzdWNjZXNzZnVs');

@$core.Deprecated('Use sendMessageBatchResultEntryDescriptor instead')
const SendMessageBatchResultEntry$json = {
  '1': 'SendMessageBatchResultEntry',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'md5ofmessageattributes',
      '3': 499647077,
      '4': 1,
      '5': 9,
      '10': 'md5ofmessageattributes'
    },
    {
      '1': 'md5ofmessagebody',
      '3': 28462758,
      '4': 1,
      '5': 9,
      '10': 'md5ofmessagebody'
    },
    {
      '1': 'md5ofmessagesystemattributes',
      '3': 304512206,
      '4': 1,
      '5': 9,
      '10': 'md5ofmessagesystemattributes'
    },
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

/// Descriptor for `SendMessageBatchResultEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageBatchResultEntryDescriptor = $convert.base64Decode(
    'ChtTZW5kTWVzc2FnZUJhdGNoUmVzdWx0RW50cnkSEgoCaWQYgfKitwEgASgJUgJpZBI6ChZtZD'
    'VvZm1lc3NhZ2VhdHRyaWJ1dGVzGOWEoO4BIAEoCVIWbWQ1b2ZtZXNzYWdlYXR0cmlidXRlcxIt'
    'ChBtZDVvZm1lc3NhZ2Vib2R5GKadyQ0gASgJUhBtZDVvZm1lc3NhZ2Vib2R5EkYKHG1kNW9mbW'
    'Vzc2FnZXN5c3RlbWF0dHJpYnV0ZXMYzvmZkQEgASgJUhxtZDVvZm1lc3NhZ2VzeXN0ZW1hdHRy'
    'aWJ1dGVzEiAKCW1lc3NhZ2VpZBiq5vSrASABKAlSCW1lc3NhZ2VpZBIpCg5zZXF1ZW5jZW51bW'
    'JlchiamuMuIAEoCVIOc2VxdWVuY2VudW1iZXI=');

@$core.Deprecated('Use sendMessageRequestDescriptor instead')
const SendMessageRequest$json = {
  '1': 'SendMessageRequest',
  '2': [
    {'1': 'delayseconds', '3': 48268198, '4': 1, '5': 5, '10': 'delayseconds'},
    {
      '1': 'messageattributes',
      '3': 56443766,
      '4': 3,
      '5': 11,
      '6': '.sqs.SendMessageRequest.MessageattributesEntry',
      '10': 'messageattributes'
    },
    {'1': 'messagebody', '3': 56920001, '4': 1, '5': 9, '10': 'messagebody'},
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
      '1': 'messagesystemattributes',
      '3': 194951997,
      '4': 3,
      '5': 11,
      '6': '.sqs.SendMessageRequest.MessagesystemattributesEntry',
      '10': 'messagesystemattributes'
    },
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
  '3': [
    SendMessageRequest_MessageattributesEntry$json,
    SendMessageRequest_MessagesystemattributesEntry$json
  ],
};

@$core.Deprecated('Use sendMessageRequestDescriptor instead')
const SendMessageRequest_MessageattributesEntry$json = {
  '1': 'MessageattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.sqs.MessageAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use sendMessageRequestDescriptor instead')
const SendMessageRequest_MessagesystemattributesEntry$json = {
  '1': 'MessagesystemattributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {
      '1': 'value',
      '3': 2,
      '4': 1,
      '5': 11,
      '6': '.sqs.MessageSystemAttributeValue',
      '10': 'value'
    },
  ],
  '7': {'7': true},
};

/// Descriptor for `SendMessageRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageRequestDescriptor = $convert.base64Decode(
    'ChJTZW5kTWVzc2FnZVJlcXVlc3QSJQoMZGVsYXlzZWNvbmRzGKaHghcgASgFUgxkZWxheXNlY2'
    '9uZHMSXwoRbWVzc2FnZWF0dHJpYnV0ZXMY9ob1GiADKAsyLi5zcXMuU2VuZE1lc3NhZ2VSZXF1'
    'ZXN0Lk1lc3NhZ2VhdHRyaWJ1dGVzRW50cnlSEW1lc3NhZ2VhdHRyaWJ1dGVzEiMKC21lc3NhZ2'
    'Vib2R5GMGPkhsgASgJUgttZXNzYWdlYm9keRI6ChZtZXNzYWdlZGVkdXBsaWNhdGlvbmlkGNnF'
    '/rQBIAEoCVIWbWVzc2FnZWRlZHVwbGljYXRpb25pZBIqCg5tZXNzYWdlZ3JvdXBpZBibxIbIAS'
    'ABKAlSDm1lc3NhZ2Vncm91cGlkEnEKF21lc3NhZ2VzeXN0ZW1hdHRyaWJ1dGVzGL32+lwgAygL'
    'MjQuc3FzLlNlbmRNZXNzYWdlUmVxdWVzdC5NZXNzYWdlc3lzdGVtYXR0cmlidXRlc0VudHJ5Uh'
    'dtZXNzYWdlc3lzdGVtYXR0cmlidXRlcxIeCghxdWV1ZXVybBjKwb7zASABKAlSCHF1ZXVldXJs'
    'GmAKFk1lc3NhZ2VhdHRyaWJ1dGVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSMAoFdmFsdWUYAi'
    'ABKAsyGi5zcXMuTWVzc2FnZUF0dHJpYnV0ZVZhbHVlUgV2YWx1ZToCOAEabAocTWVzc2FnZXN5'
    'c3RlbWF0dHJpYnV0ZXNFbnRyeRIQCgNrZXkYASABKAlSA2tleRI2CgV2YWx1ZRgCIAEoCzIgLn'
    'Nxcy5NZXNzYWdlU3lzdGVtQXR0cmlidXRlVmFsdWVSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use sendMessageResultDescriptor instead')
const SendMessageResult$json = {
  '1': 'SendMessageResult',
  '2': [
    {
      '1': 'md5ofmessageattributes',
      '3': 499647077,
      '4': 1,
      '5': 9,
      '10': 'md5ofmessageattributes'
    },
    {
      '1': 'md5ofmessagebody',
      '3': 28462758,
      '4': 1,
      '5': 9,
      '10': 'md5ofmessagebody'
    },
    {
      '1': 'md5ofmessagesystemattributes',
      '3': 304512206,
      '4': 1,
      '5': 9,
      '10': 'md5ofmessagesystemattributes'
    },
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

/// Descriptor for `SendMessageResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendMessageResultDescriptor = $convert.base64Decode(
    'ChFTZW5kTWVzc2FnZVJlc3VsdBI6ChZtZDVvZm1lc3NhZ2VhdHRyaWJ1dGVzGOWEoO4BIAEoCV'
    'IWbWQ1b2ZtZXNzYWdlYXR0cmlidXRlcxItChBtZDVvZm1lc3NhZ2Vib2R5GKadyQ0gASgJUhBt'
    'ZDVvZm1lc3NhZ2Vib2R5EkYKHG1kNW9mbWVzc2FnZXN5c3RlbWF0dHJpYnV0ZXMYzvmZkQEgAS'
    'gJUhxtZDVvZm1lc3NhZ2VzeXN0ZW1hdHRyaWJ1dGVzEiAKCW1lc3NhZ2VpZBiq5vSrASABKAlS'
    'CW1lc3NhZ2VpZBIpCg5zZXF1ZW5jZW51bWJlchiamuMuIAEoCVIOc2VxdWVuY2VudW1iZXI=');

@$core.Deprecated('Use setQueueAttributesRequestDescriptor instead')
const SetQueueAttributesRequest$json = {
  '1': 'SetQueueAttributesRequest',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 3,
      '5': 11,
      '6': '.sqs.SetQueueAttributesRequest.AttributesEntry',
      '10': 'attributes'
    },
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
  ],
  '3': [SetQueueAttributesRequest_AttributesEntry$json],
};

@$core.Deprecated('Use setQueueAttributesRequestDescriptor instead')
const SetQueueAttributesRequest_AttributesEntry$json = {
  '1': 'AttributesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `SetQueueAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List setQueueAttributesRequestDescriptor = $convert.base64Decode(
    'ChlTZXRRdWV1ZUF0dHJpYnV0ZXNSZXF1ZXN0ElEKCmF0dHJpYnV0ZXMYtan7YyADKAsyLi5zcX'
    'MuU2V0UXVldWVBdHRyaWJ1dGVzUmVxdWVzdC5BdHRyaWJ1dGVzRW50cnlSCmF0dHJpYnV0ZXMS'
    'HgoIcXVldWV1cmwYysG+8wEgASgJUghxdWV1ZXVybBo9Cg9BdHRyaWJ1dGVzRW50cnkSEAoDa2'
    'V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use startMessageMoveTaskRequestDescriptor instead')
const StartMessageMoveTaskRequest$json = {
  '1': 'StartMessageMoveTaskRequest',
  '2': [
    {
      '1': 'destinationarn',
      '3': 375726595,
      '4': 1,
      '5': 9,
      '10': 'destinationarn'
    },
    {
      '1': 'maxnumberofmessagespersecond',
      '3': 335779921,
      '4': 1,
      '5': 5,
      '10': 'maxnumberofmessagespersecond'
    },
    {'1': 'sourcearn', '3': 439903072, '4': 1, '5': 9, '10': 'sourcearn'},
  ],
};

/// Descriptor for `StartMessageMoveTaskRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startMessageMoveTaskRequestDescriptor = $convert.base64Decode(
    'ChtTdGFydE1lc3NhZ2VNb3ZlVGFza1JlcXVlc3QSKgoOZGVzdGluYXRpb25hcm4Yg8SUswEgAS'
    'gJUg5kZXN0aW5hdGlvbmFybhJGChxtYXhudW1iZXJvZm1lc3NhZ2VzcGVyc2Vjb25kGNGwjqAB'
    'IAEoBVIcbWF4bnVtYmVyb2ZtZXNzYWdlc3BlcnNlY29uZBIgCglzb3VyY2Vhcm4Y4Mbh0QEgAS'
    'gJUglzb3VyY2Vhcm4=');

@$core.Deprecated('Use startMessageMoveTaskResultDescriptor instead')
const StartMessageMoveTaskResult$json = {
  '1': 'StartMessageMoveTaskResult',
  '2': [
    {'1': 'taskhandle', '3': 190544291, '4': 1, '5': 9, '10': 'taskhandle'},
  ],
};

/// Descriptor for `StartMessageMoveTaskResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List startMessageMoveTaskResultDescriptor =
    $convert.base64Decode(
        'ChpTdGFydE1lc3NhZ2VNb3ZlVGFza1Jlc3VsdBIhCgp0YXNraGFuZGxlGKPz7VogASgJUgp0YX'
        'NraGFuZGxl');

@$core.Deprecated('Use tagQueueRequestDescriptor instead')
const TagQueueRequest$json = {
  '1': 'TagQueueRequest',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.sqs.TagQueueRequest.TagsEntry',
      '10': 'tags'
    },
  ],
  '3': [TagQueueRequest_TagsEntry$json],
};

@$core.Deprecated('Use tagQueueRequestDescriptor instead')
const TagQueueRequest_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `TagQueueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagQueueRequestDescriptor = $convert.base64Decode(
    'Cg9UYWdRdWV1ZVJlcXVlc3QSHgoIcXVldWV1cmwYysG+8wEgASgJUghxdWV1ZXVybBI2CgR0YW'
    'dzGMHB9rUBIAMoCzIeLnNxcy5UYWdRdWV1ZVJlcXVlc3QuVGFnc0VudHJ5UgR0YWdzGjcKCVRh'
    'Z3NFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use tooManyEntriesInBatchRequestDescriptor instead')
const TooManyEntriesInBatchRequest$json = {
  '1': 'TooManyEntriesInBatchRequest',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TooManyEntriesInBatchRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tooManyEntriesInBatchRequestDescriptor =
    $convert.base64Decode(
        'ChxUb29NYW55RW50cmllc0luQmF0Y2hSZXF1ZXN0EhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

@$core.Deprecated('Use unsupportedOperationDescriptor instead')
const UnsupportedOperation$json = {
  '1': 'UnsupportedOperation',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedOperation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedOperationDescriptor =
    $convert.base64Decode(
        'ChRVbnN1cHBvcnRlZE9wZXJhdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use untagQueueRequestDescriptor instead')
const UntagQueueRequest$json = {
  '1': 'UntagQueueRequest',
  '2': [
    {'1': 'queueurl', '3': 510632138, '4': 1, '5': 9, '10': 'queueurl'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagQueueRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagQueueRequestDescriptor = $convert.base64Decode(
    'ChFVbnRhZ1F1ZXVlUmVxdWVzdBIeCghxdWV1ZXVybBjKwb7zASABKAlSCHF1ZXVldXJsEhwKB3'
    'RhZ2tleXMY/MPzmAEgAygJUgd0YWdrZXlz');
