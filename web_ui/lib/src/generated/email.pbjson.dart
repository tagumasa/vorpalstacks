// This is a generated file - do not edit.
//
// Generated from email.proto.

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

@$core.Deprecated('Use attachmentContentDispositionDescriptor instead')
const AttachmentContentDisposition$json = {
  '1': 'AttachmentContentDisposition',
  '2': [
    {'1': 'ATTACHMENT_CONTENT_DISPOSITION_INLINE', '2': 0},
    {'1': 'ATTACHMENT_CONTENT_DISPOSITION_ATTACHMENT', '2': 1},
  ],
};

/// Descriptor for `AttachmentContentDisposition`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List attachmentContentDispositionDescriptor =
    $convert.base64Decode(
        'ChxBdHRhY2htZW50Q29udGVudERpc3Bvc2l0aW9uEikKJUFUVEFDSE1FTlRfQ09OVEVOVF9ESV'
        'NQT1NJVElPTl9JTkxJTkUQABItCilBVFRBQ0hNRU5UX0NPTlRFTlRfRElTUE9TSVRJT05fQVRU'
        'QUNITUVOVBAB');

@$core.Deprecated('Use attachmentContentTransferEncodingDescriptor instead')
const AttachmentContentTransferEncoding$json = {
  '1': 'AttachmentContentTransferEncoding',
  '2': [
    {'1': 'ATTACHMENT_CONTENT_TRANSFER_ENCODING_QUOTED_PRINTABLE', '2': 0},
    {'1': 'ATTACHMENT_CONTENT_TRANSFER_ENCODING_BASE64', '2': 1},
    {'1': 'ATTACHMENT_CONTENT_TRANSFER_ENCODING_SEVEN_BIT', '2': 2},
  ],
};

/// Descriptor for `AttachmentContentTransferEncoding`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List attachmentContentTransferEncodingDescriptor =
    $convert.base64Decode(
        'CiFBdHRhY2htZW50Q29udGVudFRyYW5zZmVyRW5jb2RpbmcSOQo1QVRUQUNITUVOVF9DT05URU'
        '5UX1RSQU5TRkVSX0VOQ09ESU5HX1FVT1RFRF9QUklOVEFCTEUQABIvCitBVFRBQ0hNRU5UX0NP'
        'TlRFTlRfVFJBTlNGRVJfRU5DT0RJTkdfQkFTRTY0EAESMgouQVRUQUNITUVOVF9DT05URU5UX1'
        'RSQU5TRkVSX0VOQ09ESU5HX1NFVkVOX0JJVBAC');

@$core.Deprecated('Use behaviorOnMxFailureDescriptor instead')
const BehaviorOnMxFailure$json = {
  '1': 'BehaviorOnMxFailure',
  '2': [
    {'1': 'BEHAVIOR_ON_MX_FAILURE_USE_DEFAULT_VALUE', '2': 0},
    {'1': 'BEHAVIOR_ON_MX_FAILURE_REJECT_MESSAGE', '2': 1},
  ],
};

/// Descriptor for `BehaviorOnMxFailure`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List behaviorOnMxFailureDescriptor = $convert.base64Decode(
    'ChNCZWhhdmlvck9uTXhGYWlsdXJlEiwKKEJFSEFWSU9SX09OX01YX0ZBSUxVUkVfVVNFX0RFRk'
    'FVTFRfVkFMVUUQABIpCiVCRUhBVklPUl9PTl9NWF9GQUlMVVJFX1JFSkVDVF9NRVNTQUdFEAE=');

@$core.Deprecated('Use bounceTypeDescriptor instead')
const BounceType$json = {
  '1': 'BounceType',
  '2': [
    {'1': 'BOUNCE_TYPE_UNDETERMINED', '2': 0},
    {'1': 'BOUNCE_TYPE_TRANSIENT', '2': 1},
    {'1': 'BOUNCE_TYPE_PERMANENT', '2': 2},
  ],
};

/// Descriptor for `BounceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bounceTypeDescriptor = $convert.base64Decode(
    'CgpCb3VuY2VUeXBlEhwKGEJPVU5DRV9UWVBFX1VOREVURVJNSU5FRBAAEhkKFUJPVU5DRV9UWV'
    'BFX1RSQU5TSUVOVBABEhkKFUJPVU5DRV9UWVBFX1BFUk1BTkVOVBAC');

@$core.Deprecated('Use bulkEmailStatusDescriptor instead')
const BulkEmailStatus$json = {
  '1': 'BulkEmailStatus',
  '2': [
    {'1': 'BULK_EMAIL_STATUS_MESSAGE_REJECTED', '2': 0},
    {'1': 'BULK_EMAIL_STATUS_CONFIGURATION_SET_NOT_FOUND', '2': 1},
    {'1': 'BULK_EMAIL_STATUS_MAIL_FROM_DOMAIN_NOT_VERIFIED', '2': 2},
    {'1': 'BULK_EMAIL_STATUS_SUCCESS', '2': 3},
    {'1': 'BULK_EMAIL_STATUS_INVALID_PARAMETER', '2': 4},
    {'1': 'BULK_EMAIL_STATUS_TEMPLATE_NOT_FOUND', '2': 5},
    {'1': 'BULK_EMAIL_STATUS_ACCOUNT_SUSPENDED', '2': 6},
    {'1': 'BULK_EMAIL_STATUS_ACCOUNT_DAILY_QUOTA_EXCEEDED', '2': 7},
    {'1': 'BULK_EMAIL_STATUS_TRANSIENT_FAILURE', '2': 8},
    {'1': 'BULK_EMAIL_STATUS_ACCOUNT_THROTTLED', '2': 9},
    {'1': 'BULK_EMAIL_STATUS_INVALID_SENDING_POOL_NAME', '2': 10},
    {'1': 'BULK_EMAIL_STATUS_CONFIGURATION_SET_SENDING_PAUSED', '2': 11},
    {'1': 'BULK_EMAIL_STATUS_FAILED', '2': 12},
    {'1': 'BULK_EMAIL_STATUS_ACCOUNT_SENDING_PAUSED', '2': 13},
  ],
};

/// Descriptor for `BulkEmailStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List bulkEmailStatusDescriptor = $convert.base64Decode(
    'Cg9CdWxrRW1haWxTdGF0dXMSJgoiQlVMS19FTUFJTF9TVEFUVVNfTUVTU0FHRV9SRUpFQ1RFRB'
    'AAEjEKLUJVTEtfRU1BSUxfU1RBVFVTX0NPTkZJR1VSQVRJT05fU0VUX05PVF9GT1VORBABEjMK'
    'L0JVTEtfRU1BSUxfU1RBVFVTX01BSUxfRlJPTV9ET01BSU5fTk9UX1ZFUklGSUVEEAISHQoZQl'
    'VMS19FTUFJTF9TVEFUVVNfU1VDQ0VTUxADEicKI0JVTEtfRU1BSUxfU1RBVFVTX0lOVkFMSURf'
    'UEFSQU1FVEVSEAQSKAokQlVMS19FTUFJTF9TVEFUVVNfVEVNUExBVEVfTk9UX0ZPVU5EEAUSJw'
    'ojQlVMS19FTUFJTF9TVEFUVVNfQUNDT1VOVF9TVVNQRU5ERUQQBhIyCi5CVUxLX0VNQUlMX1NU'
    'QVRVU19BQ0NPVU5UX0RBSUxZX1FVT1RBX0VYQ0VFREVEEAcSJwojQlVMS19FTUFJTF9TVEFUVV'
    'NfVFJBTlNJRU5UX0ZBSUxVUkUQCBInCiNCVUxLX0VNQUlMX1NUQVRVU19BQ0NPVU5UX1RIUk9U'
    'VExFRBAJEi8KK0JVTEtfRU1BSUxfU1RBVFVTX0lOVkFMSURfU0VORElOR19QT09MX05BTUUQCh'
    'I2CjJCVUxLX0VNQUlMX1NUQVRVU19DT05GSUdVUkFUSU9OX1NFVF9TRU5ESU5HX1BBVVNFRBAL'
    'EhwKGEJVTEtfRU1BSUxfU1RBVFVTX0ZBSUxFRBAMEiwKKEJVTEtfRU1BSUxfU1RBVFVTX0FDQ0'
    '9VTlRfU0VORElOR19QQVVTRUQQDQ==');

@$core.Deprecated('Use contactLanguageDescriptor instead')
const ContactLanguage$json = {
  '1': 'ContactLanguage',
  '2': [
    {'1': 'CONTACT_LANGUAGE_EN', '2': 0},
    {'1': 'CONTACT_LANGUAGE_JA', '2': 1},
  ],
};

/// Descriptor for `ContactLanguage`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contactLanguageDescriptor = $convert.base64Decode(
    'Cg9Db250YWN0TGFuZ3VhZ2USFwoTQ09OVEFDVF9MQU5HVUFHRV9FThAAEhcKE0NPTlRBQ1RfTE'
    'FOR1VBR0VfSkEQAQ==');

@$core.Deprecated('Use contactListImportActionDescriptor instead')
const ContactListImportAction$json = {
  '1': 'ContactListImportAction',
  '2': [
    {'1': 'CONTACT_LIST_IMPORT_ACTION_DELETE', '2': 0},
    {'1': 'CONTACT_LIST_IMPORT_ACTION_PUT', '2': 1},
  ],
};

/// Descriptor for `ContactListImportAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List contactListImportActionDescriptor =
    $convert.base64Decode(
        'ChdDb250YWN0TGlzdEltcG9ydEFjdGlvbhIlCiFDT05UQUNUX0xJU1RfSU1QT1JUX0FDVElPTl'
        '9ERUxFVEUQABIiCh5DT05UQUNUX0xJU1RfSU1QT1JUX0FDVElPTl9QVVQQAQ==');

@$core.Deprecated('Use dataFormatDescriptor instead')
const DataFormat$json = {
  '1': 'DataFormat',
  '2': [
    {'1': 'DATA_FORMAT_JSON', '2': 0},
    {'1': 'DATA_FORMAT_CSV', '2': 1},
  ],
};

/// Descriptor for `DataFormat`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dataFormatDescriptor = $convert.base64Decode(
    'CgpEYXRhRm9ybWF0EhQKEERBVEFfRk9STUFUX0pTT04QABITCg9EQVRBX0ZPUk1BVF9DU1YQAQ'
    '==');

@$core.Deprecated('Use deliverabilityDashboardAccountStatusDescriptor instead')
const DeliverabilityDashboardAccountStatus$json = {
  '1': 'DeliverabilityDashboardAccountStatus',
  '2': [
    {'1': 'DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_DISABLED', '2': 0},
    {'1': 'DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_PENDING_EXPIRATION', '2': 1},
    {'1': 'DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_ACTIVE', '2': 2},
  ],
};

/// Descriptor for `DeliverabilityDashboardAccountStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deliverabilityDashboardAccountStatusDescriptor =
    $convert.base64Decode(
        'CiREZWxpdmVyYWJpbGl0eURhc2hib2FyZEFjY291bnRTdGF0dXMSNAowREVMSVZFUkFCSUxJVF'
        'lfREFTSEJPQVJEX0FDQ09VTlRfU1RBVFVTX0RJU0FCTEVEEAASPgo6REVMSVZFUkFCSUxJVFlf'
        'REFTSEJPQVJEX0FDQ09VTlRfU1RBVFVTX1BFTkRJTkdfRVhQSVJBVElPThABEjIKLkRFTElWRV'
        'JBQklMSVRZX0RBU0hCT0FSRF9BQ0NPVU5UX1NUQVRVU19BQ1RJVkUQAg==');

@$core.Deprecated('Use deliverabilityTestStatusDescriptor instead')
const DeliverabilityTestStatus$json = {
  '1': 'DeliverabilityTestStatus',
  '2': [
    {'1': 'DELIVERABILITY_TEST_STATUS_IN_PROGRESS', '2': 0},
    {'1': 'DELIVERABILITY_TEST_STATUS_COMPLETED', '2': 1},
  ],
};

/// Descriptor for `DeliverabilityTestStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deliverabilityTestStatusDescriptor = $convert.base64Decode(
    'ChhEZWxpdmVyYWJpbGl0eVRlc3RTdGF0dXMSKgomREVMSVZFUkFCSUxJVFlfVEVTVF9TVEFUVV'
    'NfSU5fUFJPR1JFU1MQABIoCiRERUxJVkVSQUJJTElUWV9URVNUX1NUQVRVU19DT01QTEVURUQQ'
    'AQ==');

@$core.Deprecated('Use deliveryEventTypeDescriptor instead')
const DeliveryEventType$json = {
  '1': 'DeliveryEventType',
  '2': [
    {'1': 'DELIVERY_EVENT_TYPE_TRANSIENT_BOUNCE', '2': 0},
    {'1': 'DELIVERY_EVENT_TYPE_SEND', '2': 1},
    {'1': 'DELIVERY_EVENT_TYPE_PERMANENT_BOUNCE', '2': 2},
    {'1': 'DELIVERY_EVENT_TYPE_COMPLAINT', '2': 3},
    {'1': 'DELIVERY_EVENT_TYPE_DELIVERY', '2': 4},
    {'1': 'DELIVERY_EVENT_TYPE_UNDETERMINED_BOUNCE', '2': 5},
  ],
};

/// Descriptor for `DeliveryEventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List deliveryEventTypeDescriptor = $convert.base64Decode(
    'ChFEZWxpdmVyeUV2ZW50VHlwZRIoCiRERUxJVkVSWV9FVkVOVF9UWVBFX1RSQU5TSUVOVF9CT1'
    'VOQ0UQABIcChhERUxJVkVSWV9FVkVOVF9UWVBFX1NFTkQQARIoCiRERUxJVkVSWV9FVkVOVF9U'
    'WVBFX1BFUk1BTkVOVF9CT1VOQ0UQAhIhCh1ERUxJVkVSWV9FVkVOVF9UWVBFX0NPTVBMQUlOVB'
    'ADEiAKHERFTElWRVJZX0VWRU5UX1RZUEVfREVMSVZFUlkQBBIrCidERUxJVkVSWV9FVkVOVF9U'
    'WVBFX1VOREVURVJNSU5FRF9CT1VOQ0UQBQ==');

@$core.Deprecated('Use dimensionValueSourceDescriptor instead')
const DimensionValueSource$json = {
  '1': 'DimensionValueSource',
  '2': [
    {'1': 'DIMENSION_VALUE_SOURCE_MESSAGE_TAG', '2': 0},
    {'1': 'DIMENSION_VALUE_SOURCE_EMAIL_HEADER', '2': 1},
    {'1': 'DIMENSION_VALUE_SOURCE_LINK_TAG', '2': 2},
  ],
};

/// Descriptor for `DimensionValueSource`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dimensionValueSourceDescriptor = $convert.base64Decode(
    'ChREaW1lbnNpb25WYWx1ZVNvdXJjZRImCiJESU1FTlNJT05fVkFMVUVfU09VUkNFX01FU1NBR0'
    'VfVEFHEAASJwojRElNRU5TSU9OX1ZBTFVFX1NPVVJDRV9FTUFJTF9IRUFERVIQARIjCh9ESU1F'
    'TlNJT05fVkFMVUVfU09VUkNFX0xJTktfVEFHEAI=');

@$core.Deprecated('Use dkimSigningAttributesOriginDescriptor instead')
const DkimSigningAttributesOrigin$json = {
  '1': 'DkimSigningAttributesOrigin',
  '2': [
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_3', '2': 0},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_CENTRAL_1', '2': 1},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_CENTRAL_1', '2': 2},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_2', '2': 3},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_EXTERNAL', '2': 4},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_1', '2': 5},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_1', '2': 6},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_1', '2': 7},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_5', '2': 8},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES', '2': 9},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_2', '2': 10},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_WEST_1', '2': 11},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_2', '2': 12},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_SA_EAST_1', '2': 13},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_3', '2': 14},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_SOUTH_1', '2': 15},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_2', '2': 16},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_2', '2': 17},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_NORTH_1', '2': 18},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_1', '2': 19},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_1', '2': 20},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AF_SOUTH_1', '2': 21},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_2', '2': 22},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_2', '2': 23},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_3', '2': 24},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_SOUTH_1', '2': 25},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_1', '2': 26},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_1', '2': 27},
    {'1': 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_IL_CENTRAL_1', '2': 28},
  ],
};

/// Descriptor for `DkimSigningAttributesOrigin`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dkimSigningAttributesOriginDescriptor = $convert.base64Decode(
    'ChtEa2ltU2lnbmluZ0F0dHJpYnV0ZXNPcmlnaW4SOQo1REtJTV9TSUdOSU5HX0FUVFJJQlVURV'
    'NfT1JJR0lOX0FXU19TRVNfQVBfTk9SVEhFQVNUXzMQABI3CjNES0lNX1NJR05JTkdfQVRUUklC'
    'VVRFU19PUklHSU5fQVdTX1NFU19DQV9DRU5UUkFMXzEQARI3CjNES0lNX1NJR05JTkdfQVRUUk'
    'lCVVRFU19PUklHSU5fQVdTX1NFU19NRV9DRU5UUkFMXzEQAhI3CjNES0lNX1NJR05JTkdfQVRU'
    'UklCVVRFU19PUklHSU5fQVdTX1NFU19FVV9DRU5UUkFMXzIQAxIrCidES0lNX1NJR05JTkdfQV'
    'RUUklCVVRFU19PUklHSU5fRVhURVJOQUwQBBI0CjBES0lNX1NJR05JTkdfQVRUUklCVVRFU19P'
    'UklHSU5fQVdTX1NFU19FVV9XRVNUXzEQBRI1CjFES0lNX1NJR05JTkdfQVRUUklCVVRFU19PUk'
    'lHSU5fQVdTX1NFU19BUF9TT1VUSF8xEAYSOQo1REtJTV9TSUdOSU5HX0FUVFJJQlVURVNfT1JJ'
    'R0lOX0FXU19TRVNfQVBfTk9SVEhFQVNUXzEQBxI5CjVES0lNX1NJR05JTkdfQVRUUklCVVRFU1'
    '9PUklHSU5fQVdTX1NFU19BUF9TT1VUSEVBU1RfNRAIEioKJkRLSU1fU0lHTklOR19BVFRSSUJV'
    'VEVTX09SSUdJTl9BV1NfU0VTEAkSOQo1REtJTV9TSUdOSU5HX0FUVFJJQlVURVNfT1JJR0lOX0'
    'FXU19TRVNfQVBfU09VVEhFQVNUXzIQChI0CjBES0lNX1NJR05JTkdfQVRUUklCVVRFU19PUklH'
    'SU5fQVdTX1NFU19DQV9XRVNUXzEQCxI5CjVES0lNX1NJR05JTkdfQVRUUklCVVRFU19PUklHSU'
    '5fQVdTX1NFU19BUF9OT1JUSEVBU1RfMhAMEjQKMERLSU1fU0lHTklOR19BVFRSSUJVVEVTX09S'
    'SUdJTl9BV1NfU0VTX1NBX0VBU1RfMRANEjQKMERLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSU'
    'dJTl9BV1NfU0VTX0VVX1dFU1RfMxAOEjUKMURLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJ'
    'Tl9BV1NfU0VTX0VVX1NPVVRIXzEQDxI1CjFES0lNX1NJR05JTkdfQVRUUklCVVRFU19PUklHSU'
    '5fQVdTX1NFU19BUF9TT1VUSF8yEBASNAowREtJTV9TSUdOSU5HX0FUVFJJQlVURVNfT1JJR0lO'
    'X0FXU19TRVNfVVNfRUFTVF8yEBESNQoxREtJTV9TSUdOSU5HX0FUVFJJQlVURVNfT1JJR0lOX0'
    'FXU19TRVNfRVVfTk9SVEhfMRASEjQKMERLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl9B'
    'V1NfU0VTX1VTX1dFU1RfMRATEjkKNURLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl9BV1'
    'NfU0VTX0FQX1NPVVRIRUFTVF8xEBQSNQoxREtJTV9TSUdOSU5HX0FUVFJJQlVURVNfT1JJR0lO'
    'X0FXU19TRVNfQUZfU09VVEhfMRAVEjQKMERLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl'
    '9BV1NfU0VTX0VVX1dFU1RfMhAWEjQKMERLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl9B'
    'V1NfU0VTX1VTX1dFU1RfMhAXEjkKNURLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl9BV1'
    'NfU0VTX0FQX1NPVVRIRUFTVF8zEBgSNQoxREtJTV9TSUdOSU5HX0FUVFJJQlVURVNfT1JJR0lO'
    'X0FXU19TRVNfTUVfU09VVEhfMRAZEjQKMERLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl'
    '9BV1NfU0VTX1VTX0VBU1RfMRAaEjcKM0RLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl9B'
    'V1NfU0VTX0VVX0NFTlRSQUxfMRAbEjcKM0RLSU1fU0lHTklOR19BVFRSSUJVVEVTX09SSUdJTl'
    '9BV1NfU0VTX0lMX0NFTlRSQUxfMRAc');

@$core.Deprecated('Use dkimSigningKeyLengthDescriptor instead')
const DkimSigningKeyLength$json = {
  '1': 'DkimSigningKeyLength',
  '2': [
    {'1': 'DKIM_SIGNING_KEY_LENGTH_RSA_1024_BIT', '2': 0},
    {'1': 'DKIM_SIGNING_KEY_LENGTH_RSA_2048_BIT', '2': 1},
  ],
};

/// Descriptor for `DkimSigningKeyLength`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dkimSigningKeyLengthDescriptor = $convert.base64Decode(
    'ChREa2ltU2lnbmluZ0tleUxlbmd0aBIoCiRES0lNX1NJR05JTkdfS0VZX0xFTkdUSF9SU0FfMT'
    'AyNF9CSVQQABIoCiRES0lNX1NJR05JTkdfS0VZX0xFTkdUSF9SU0FfMjA0OF9CSVQQAQ==');

@$core.Deprecated('Use dkimStatusDescriptor instead')
const DkimStatus$json = {
  '1': 'DkimStatus',
  '2': [
    {'1': 'DKIM_STATUS_PENDING', '2': 0},
    {'1': 'DKIM_STATUS_SUCCESS', '2': 1},
    {'1': 'DKIM_STATUS_TEMPORARY_FAILURE', '2': 2},
    {'1': 'DKIM_STATUS_FAILED', '2': 3},
    {'1': 'DKIM_STATUS_NOT_STARTED', '2': 4},
  ],
};

/// Descriptor for `DkimStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dkimStatusDescriptor = $convert.base64Decode(
    'CgpEa2ltU3RhdHVzEhcKE0RLSU1fU1RBVFVTX1BFTkRJTkcQABIXChNES0lNX1NUQVRVU19TVU'
    'NDRVNTEAESIQodREtJTV9TVEFUVVNfVEVNUE9SQVJZX0ZBSUxVUkUQAhIWChJES0lNX1NUQVRV'
    'U19GQUlMRUQQAxIbChdES0lNX1NUQVRVU19OT1RfU1RBUlRFRBAE');

@$core.Deprecated('Use emailAddressInsightsConfidenceVerdictDescriptor instead')
const EmailAddressInsightsConfidenceVerdict$json = {
  '1': 'EmailAddressInsightsConfidenceVerdict',
  '2': [
    {'1': 'EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_MEDIUM', '2': 0},
    {'1': 'EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_LOW', '2': 1},
    {'1': 'EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_HIGH', '2': 2},
  ],
};

/// Descriptor for `EmailAddressInsightsConfidenceVerdict`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List emailAddressInsightsConfidenceVerdictDescriptor =
    $convert.base64Decode(
        'CiVFbWFpbEFkZHJlc3NJbnNpZ2h0c0NvbmZpZGVuY2VWZXJkaWN0EjQKMEVNQUlMX0FERFJFU1'
        'NfSU5TSUdIVFNfQ09ORklERU5DRV9WRVJESUNUX01FRElVTRAAEjEKLUVNQUlMX0FERFJFU1Nf'
        'SU5TSUdIVFNfQ09ORklERU5DRV9WRVJESUNUX0xPVxABEjIKLkVNQUlMX0FERFJFU1NfSU5TSU'
        'dIVFNfQ09ORklERU5DRV9WRVJESUNUX0hJR0gQAg==');

@$core.Deprecated('Use engagementEventTypeDescriptor instead')
const EngagementEventType$json = {
  '1': 'EngagementEventType',
  '2': [
    {'1': 'ENGAGEMENT_EVENT_TYPE_OPEN', '2': 0},
    {'1': 'ENGAGEMENT_EVENT_TYPE_CLICK', '2': 1},
  ],
};

/// Descriptor for `EngagementEventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List engagementEventTypeDescriptor = $convert.base64Decode(
    'ChNFbmdhZ2VtZW50RXZlbnRUeXBlEh4KGkVOR0FHRU1FTlRfRVZFTlRfVFlQRV9PUEVOEAASHw'
    'obRU5HQUdFTUVOVF9FVkVOVF9UWVBFX0NMSUNLEAE=');

@$core.Deprecated('Use eventTypeDescriptor instead')
const EventType$json = {
  '1': 'EventType',
  '2': [
    {'1': 'EVENT_TYPE_DELIVERY_DELAY', '2': 0},
    {'1': 'EVENT_TYPE_RENDERING_FAILURE', '2': 1},
    {'1': 'EVENT_TYPE_SEND', '2': 2},
    {'1': 'EVENT_TYPE_BOUNCE', '2': 3},
    {'1': 'EVENT_TYPE_REJECT', '2': 4},
    {'1': 'EVENT_TYPE_OPEN', '2': 5},
    {'1': 'EVENT_TYPE_SUBSCRIPTION', '2': 6},
    {'1': 'EVENT_TYPE_COMPLAINT', '2': 7},
    {'1': 'EVENT_TYPE_DELIVERY', '2': 8},
    {'1': 'EVENT_TYPE_CLICK', '2': 9},
  ],
};

/// Descriptor for `EventType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List eventTypeDescriptor = $convert.base64Decode(
    'CglFdmVudFR5cGUSHQoZRVZFTlRfVFlQRV9ERUxJVkVSWV9ERUxBWRAAEiAKHEVWRU5UX1RZUE'
    'VfUkVOREVSSU5HX0ZBSUxVUkUQARITCg9FVkVOVF9UWVBFX1NFTkQQAhIVChFFVkVOVF9UWVBF'
    'X0JPVU5DRRADEhUKEUVWRU5UX1RZUEVfUkVKRUNUEAQSEwoPRVZFTlRfVFlQRV9PUEVOEAUSGw'
    'oXRVZFTlRfVFlQRV9TVUJTQ1JJUFRJT04QBhIYChRFVkVOVF9UWVBFX0NPTVBMQUlOVBAHEhcK'
    'E0VWRU5UX1RZUEVfREVMSVZFUlkQCBIUChBFVkVOVF9UWVBFX0NMSUNLEAk=');

@$core.Deprecated('Use exportSourceTypeDescriptor instead')
const ExportSourceType$json = {
  '1': 'ExportSourceType',
  '2': [
    {'1': 'EXPORT_SOURCE_TYPE_METRICS_DATA', '2': 0},
    {'1': 'EXPORT_SOURCE_TYPE_MESSAGE_INSIGHTS', '2': 1},
  ],
};

/// Descriptor for `ExportSourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List exportSourceTypeDescriptor = $convert.base64Decode(
    'ChBFeHBvcnRTb3VyY2VUeXBlEiMKH0VYUE9SVF9TT1VSQ0VfVFlQRV9NRVRSSUNTX0RBVEEQAB'
    'InCiNFWFBPUlRfU09VUkNFX1RZUEVfTUVTU0FHRV9JTlNJR0hUUxAB');

@$core.Deprecated('Use featureStatusDescriptor instead')
const FeatureStatus$json = {
  '1': 'FeatureStatus',
  '2': [
    {'1': 'FEATURE_STATUS_DISABLED', '2': 0},
    {'1': 'FEATURE_STATUS_ENABLED', '2': 1},
  ],
};

/// Descriptor for `FeatureStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List featureStatusDescriptor = $convert.base64Decode(
    'Cg1GZWF0dXJlU3RhdHVzEhsKF0ZFQVRVUkVfU1RBVFVTX0RJU0FCTEVEEAASGgoWRkVBVFVSRV'
    '9TVEFUVVNfRU5BQkxFRBAB');

@$core.Deprecated('Use httpsPolicyDescriptor instead')
const HttpsPolicy$json = {
  '1': 'HttpsPolicy',
  '2': [
    {'1': 'HTTPS_POLICY_OPTIONAL', '2': 0},
    {'1': 'HTTPS_POLICY_REQUIRE', '2': 1},
    {'1': 'HTTPS_POLICY_REQUIRE_OPEN_ONLY', '2': 2},
  ],
};

/// Descriptor for `HttpsPolicy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List httpsPolicyDescriptor = $convert.base64Decode(
    'CgtIdHRwc1BvbGljeRIZChVIVFRQU19QT0xJQ1lfT1BUSU9OQUwQABIYChRIVFRQU19QT0xJQ1'
    'lfUkVRVUlSRRABEiIKHkhUVFBTX1BPTElDWV9SRVFVSVJFX09QRU5fT05MWRAC');

@$core.Deprecated('Use identityTypeDescriptor instead')
const IdentityType$json = {
  '1': 'IdentityType',
  '2': [
    {'1': 'IDENTITY_TYPE_MANAGED_DOMAIN', '2': 0},
    {'1': 'IDENTITY_TYPE_DOMAIN', '2': 1},
    {'1': 'IDENTITY_TYPE_EMAIL_ADDRESS', '2': 2},
  ],
};

/// Descriptor for `IdentityType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List identityTypeDescriptor = $convert.base64Decode(
    'CgxJZGVudGl0eVR5cGUSIAocSURFTlRJVFlfVFlQRV9NQU5BR0VEX0RPTUFJThAAEhgKFElERU'
    '5USVRZX1RZUEVfRE9NQUlOEAESHwobSURFTlRJVFlfVFlQRV9FTUFJTF9BRERSRVNTEAI=');

@$core.Deprecated('Use importDestinationTypeDescriptor instead')
const ImportDestinationType$json = {
  '1': 'ImportDestinationType',
  '2': [
    {'1': 'IMPORT_DESTINATION_TYPE_SUPPRESSION_LIST', '2': 0},
    {'1': 'IMPORT_DESTINATION_TYPE_CONTACT_LIST', '2': 1},
  ],
};

/// Descriptor for `ImportDestinationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importDestinationTypeDescriptor = $convert.base64Decode(
    'ChVJbXBvcnREZXN0aW5hdGlvblR5cGUSLAooSU1QT1JUX0RFU1RJTkFUSU9OX1RZUEVfU1VQUF'
    'JFU1NJT05fTElTVBAAEigKJElNUE9SVF9ERVNUSU5BVElPTl9UWVBFX0NPTlRBQ1RfTElTVBAB');

@$core.Deprecated('Use jobStatusDescriptor instead')
const JobStatus$json = {
  '1': 'JobStatus',
  '2': [
    {'1': 'JOB_STATUS_PROCESSING', '2': 0},
    {'1': 'JOB_STATUS_CANCELLED', '2': 1},
    {'1': 'JOB_STATUS_COMPLETED', '2': 2},
    {'1': 'JOB_STATUS_CREATED', '2': 3},
    {'1': 'JOB_STATUS_FAILED', '2': 4},
  ],
};

/// Descriptor for `JobStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List jobStatusDescriptor = $convert.base64Decode(
    'CglKb2JTdGF0dXMSGQoVSk9CX1NUQVRVU19QUk9DRVNTSU5HEAASGAoUSk9CX1NUQVRVU19DQU'
    '5DRUxMRUQQARIYChRKT0JfU1RBVFVTX0NPTVBMRVRFRBACEhYKEkpPQl9TVEFUVVNfQ1JFQVRF'
    'RBADEhUKEUpPQl9TVEFUVVNfRkFJTEVEEAQ=');

@$core.Deprecated('Use listRecommendationsFilterKeyDescriptor instead')
const ListRecommendationsFilterKey$json = {
  '1': 'ListRecommendationsFilterKey',
  '2': [
    {'1': 'LIST_RECOMMENDATIONS_FILTER_KEY_STATUS', '2': 0},
    {'1': 'LIST_RECOMMENDATIONS_FILTER_KEY_IMPACT', '2': 1},
    {'1': 'LIST_RECOMMENDATIONS_FILTER_KEY_TYPE', '2': 2},
    {'1': 'LIST_RECOMMENDATIONS_FILTER_KEY_RESOURCE_ARN', '2': 3},
  ],
};

/// Descriptor for `ListRecommendationsFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List listRecommendationsFilterKeyDescriptor = $convert.base64Decode(
    'ChxMaXN0UmVjb21tZW5kYXRpb25zRmlsdGVyS2V5EioKJkxJU1RfUkVDT01NRU5EQVRJT05TX0'
    'ZJTFRFUl9LRVlfU1RBVFVTEAASKgomTElTVF9SRUNPTU1FTkRBVElPTlNfRklMVEVSX0tFWV9J'
    'TVBBQ1QQARIoCiRMSVNUX1JFQ09NTUVOREFUSU9OU19GSUxURVJfS0VZX1RZUEUQAhIwCixMSV'
    'NUX1JFQ09NTUVOREFUSU9OU19GSUxURVJfS0VZX1JFU09VUkNFX0FSThAD');

@$core.Deprecated('Use listTenantResourcesFilterKeyDescriptor instead')
const ListTenantResourcesFilterKey$json = {
  '1': 'ListTenantResourcesFilterKey',
  '2': [
    {'1': 'LIST_TENANT_RESOURCES_FILTER_KEY_RESOURCE_TYPE', '2': 0},
  ],
};

/// Descriptor for `ListTenantResourcesFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List listTenantResourcesFilterKeyDescriptor =
    $convert.base64Decode(
        'ChxMaXN0VGVuYW50UmVzb3VyY2VzRmlsdGVyS2V5EjIKLkxJU1RfVEVOQU5UX1JFU09VUkNFU1'
        '9GSUxURVJfS0VZX1JFU09VUkNFX1RZUEUQAA==');

@$core.Deprecated('Use mailFromDomainStatusDescriptor instead')
const MailFromDomainStatus$json = {
  '1': 'MailFromDomainStatus',
  '2': [
    {'1': 'MAIL_FROM_DOMAIN_STATUS_PENDING', '2': 0},
    {'1': 'MAIL_FROM_DOMAIN_STATUS_SUCCESS', '2': 1},
    {'1': 'MAIL_FROM_DOMAIN_STATUS_TEMPORARY_FAILURE', '2': 2},
    {'1': 'MAIL_FROM_DOMAIN_STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `MailFromDomainStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mailFromDomainStatusDescriptor = $convert.base64Decode(
    'ChRNYWlsRnJvbURvbWFpblN0YXR1cxIjCh9NQUlMX0ZST01fRE9NQUlOX1NUQVRVU19QRU5ESU'
    '5HEAASIwofTUFJTF9GUk9NX0RPTUFJTl9TVEFUVVNfU1VDQ0VTUxABEi0KKU1BSUxfRlJPTV9E'
    'T01BSU5fU1RBVFVTX1RFTVBPUkFSWV9GQUlMVVJFEAISIgoeTUFJTF9GUk9NX0RPTUFJTl9TVE'
    'FUVVNfRkFJTEVEEAM=');

@$core.Deprecated('Use mailTypeDescriptor instead')
const MailType$json = {
  '1': 'MailType',
  '2': [
    {'1': 'MAIL_TYPE_MARKETING', '2': 0},
    {'1': 'MAIL_TYPE_TRANSACTIONAL', '2': 1},
  ],
};

/// Descriptor for `MailType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List mailTypeDescriptor = $convert.base64Decode(
    'CghNYWlsVHlwZRIXChNNQUlMX1RZUEVfTUFSS0VUSU5HEAASGwoXTUFJTF9UWVBFX1RSQU5TQU'
    'NUSU9OQUwQAQ==');

@$core.Deprecated('Use metricDescriptor instead')
const Metric$json = {
  '1': 'Metric',
  '2': [
    {'1': 'METRIC_TRANSIENT_BOUNCE', '2': 0},
    {'1': 'METRIC_DELIVERY_CLICK', '2': 1},
    {'1': 'METRIC_SEND', '2': 2},
    {'1': 'METRIC_OPEN', '2': 3},
    {'1': 'METRIC_DELIVERY_COMPLAINT', '2': 4},
    {'1': 'METRIC_PERMANENT_BOUNCE', '2': 5},
    {'1': 'METRIC_DELIVERY_OPEN', '2': 6},
    {'1': 'METRIC_COMPLAINT', '2': 7},
    {'1': 'METRIC_DELIVERY', '2': 8},
    {'1': 'METRIC_CLICK', '2': 9},
  ],
};

/// Descriptor for `Metric`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricDescriptor = $convert.base64Decode(
    'CgZNZXRyaWMSGwoXTUVUUklDX1RSQU5TSUVOVF9CT1VOQ0UQABIZChVNRVRSSUNfREVMSVZFUl'
    'lfQ0xJQ0sQARIPCgtNRVRSSUNfU0VORBACEg8KC01FVFJJQ19PUEVOEAMSHQoZTUVUUklDX0RF'
    'TElWRVJZX0NPTVBMQUlOVBAEEhsKF01FVFJJQ19QRVJNQU5FTlRfQk9VTkNFEAUSGAoUTUVUUk'
    'lDX0RFTElWRVJZX09QRU4QBhIUChBNRVRSSUNfQ09NUExBSU5UEAcSEwoPTUVUUklDX0RFTElW'
    'RVJZEAgSEAoMTUVUUklDX0NMSUNLEAk=');

@$core.Deprecated('Use metricAggregationDescriptor instead')
const MetricAggregation$json = {
  '1': 'MetricAggregation',
  '2': [
    {'1': 'METRIC_AGGREGATION_VOLUME', '2': 0},
    {'1': 'METRIC_AGGREGATION_RATE', '2': 1},
  ],
};

/// Descriptor for `MetricAggregation`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricAggregationDescriptor = $convert.base64Decode(
    'ChFNZXRyaWNBZ2dyZWdhdGlvbhIdChlNRVRSSUNfQUdHUkVHQVRJT05fVk9MVU1FEAASGwoXTU'
    'VUUklDX0FHR1JFR0FUSU9OX1JBVEUQAQ==');

@$core.Deprecated('Use metricDimensionNameDescriptor instead')
const MetricDimensionName$json = {
  '1': 'MetricDimensionName',
  '2': [
    {'1': 'METRIC_DIMENSION_NAME_EMAIL_IDENTITY', '2': 0},
    {'1': 'METRIC_DIMENSION_NAME_ISP', '2': 1},
    {'1': 'METRIC_DIMENSION_NAME_CONFIGURATION_SET', '2': 2},
  ],
};

/// Descriptor for `MetricDimensionName`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricDimensionNameDescriptor = $convert.base64Decode(
    'ChNNZXRyaWNEaW1lbnNpb25OYW1lEigKJE1FVFJJQ19ESU1FTlNJT05fTkFNRV9FTUFJTF9JRE'
    'VOVElUWRAAEh0KGU1FVFJJQ19ESU1FTlNJT05fTkFNRV9JU1AQARIrCidNRVRSSUNfRElNRU5T'
    'SU9OX05BTUVfQ09ORklHVVJBVElPTl9TRVQQAg==');

@$core.Deprecated('Use metricNamespaceDescriptor instead')
const MetricNamespace$json = {
  '1': 'MetricNamespace',
  '2': [
    {'1': 'METRIC_NAMESPACE_VDM', '2': 0},
  ],
};

/// Descriptor for `MetricNamespace`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List metricNamespaceDescriptor = $convert.base64Decode(
    'Cg9NZXRyaWNOYW1lc3BhY2USGAoUTUVUUklDX05BTUVTUEFDRV9WRE0QAA==');

@$core.Deprecated('Use queryErrorCodeDescriptor instead')
const QueryErrorCode$json = {
  '1': 'QueryErrorCode',
  '2': [
    {'1': 'QUERY_ERROR_CODE_ACCESS_DENIED', '2': 0},
    {'1': 'QUERY_ERROR_CODE_INTERNAL_FAILURE', '2': 1},
  ],
};

/// Descriptor for `QueryErrorCode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List queryErrorCodeDescriptor = $convert.base64Decode(
    'Cg5RdWVyeUVycm9yQ29kZRIiCh5RVUVSWV9FUlJPUl9DT0RFX0FDQ0VTU19ERU5JRUQQABIlCi'
    'FRVUVSWV9FUlJPUl9DT0RFX0lOVEVSTkFMX0ZBSUxVUkUQAQ==');

@$core.Deprecated('Use recommendationImpactDescriptor instead')
const RecommendationImpact$json = {
  '1': 'RecommendationImpact',
  '2': [
    {'1': 'RECOMMENDATION_IMPACT_LOW', '2': 0},
    {'1': 'RECOMMENDATION_IMPACT_HIGH', '2': 1},
  ],
};

/// Descriptor for `RecommendationImpact`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recommendationImpactDescriptor = $convert.base64Decode(
    'ChRSZWNvbW1lbmRhdGlvbkltcGFjdBIdChlSRUNPTU1FTkRBVElPTl9JTVBBQ1RfTE9XEAASHg'
    'oaUkVDT01NRU5EQVRJT05fSU1QQUNUX0hJR0gQAQ==');

@$core.Deprecated('Use recommendationStatusDescriptor instead')
const RecommendationStatus$json = {
  '1': 'RecommendationStatus',
  '2': [
    {'1': 'RECOMMENDATION_STATUS_FIXED', '2': 0},
    {'1': 'RECOMMENDATION_STATUS_OPEN', '2': 1},
  ],
};

/// Descriptor for `RecommendationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recommendationStatusDescriptor = $convert.base64Decode(
    'ChRSZWNvbW1lbmRhdGlvblN0YXR1cxIfChtSRUNPTU1FTkRBVElPTl9TVEFUVVNfRklYRUQQAB'
    'IeChpSRUNPTU1FTkRBVElPTl9TVEFUVVNfT1BFThAB');

@$core.Deprecated('Use recommendationTypeDescriptor instead')
const RecommendationType$json = {
  '1': 'RecommendationType',
  '2': [
    {'1': 'RECOMMENDATION_TYPE_DMARC', '2': 0},
    {'1': 'RECOMMENDATION_TYPE_BOUNCE', '2': 1},
    {'1': 'RECOMMENDATION_TYPE_FEEDBACK_3P', '2': 2},
    {'1': 'RECOMMENDATION_TYPE_IP_LISTING', '2': 3},
    {'1': 'RECOMMENDATION_TYPE_COMPLAINT', '2': 4},
    {'1': 'RECOMMENDATION_TYPE_DKIM', '2': 5},
    {'1': 'RECOMMENDATION_TYPE_BIMI', '2': 6},
    {'1': 'RECOMMENDATION_TYPE_SPF', '2': 7},
  ],
};

/// Descriptor for `RecommendationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List recommendationTypeDescriptor = $convert.base64Decode(
    'ChJSZWNvbW1lbmRhdGlvblR5cGUSHQoZUkVDT01NRU5EQVRJT05fVFlQRV9ETUFSQxAAEh4KGl'
    'JFQ09NTUVOREFUSU9OX1RZUEVfQk9VTkNFEAESIwofUkVDT01NRU5EQVRJT05fVFlQRV9GRUVE'
    'QkFDS18zUBACEiIKHlJFQ09NTUVOREFUSU9OX1RZUEVfSVBfTElTVElORxADEiEKHVJFQ09NTU'
    'VOREFUSU9OX1RZUEVfQ09NUExBSU5UEAQSHAoYUkVDT01NRU5EQVRJT05fVFlQRV9ES0lNEAUS'
    'HAoYUkVDT01NRU5EQVRJT05fVFlQRV9CSU1JEAYSGwoXUkVDT01NRU5EQVRJT05fVFlQRV9TUE'
    'YQBw==');

@$core.Deprecated('Use reputationEntityFilterKeyDescriptor instead')
const ReputationEntityFilterKey$json = {
  '1': 'ReputationEntityFilterKey',
  '2': [
    {'1': 'REPUTATION_ENTITY_FILTER_KEY_STATUS', '2': 0},
    {'1': 'REPUTATION_ENTITY_FILTER_KEY_ENTITY_TYPE', '2': 1},
    {'1': 'REPUTATION_ENTITY_FILTER_KEY_ENTITY_REFERENCE_PREFIX', '2': 2},
    {'1': 'REPUTATION_ENTITY_FILTER_KEY_REPUTATION_IMPACT', '2': 3},
  ],
};

/// Descriptor for `ReputationEntityFilterKey`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reputationEntityFilterKeyDescriptor = $convert.base64Decode(
    'ChlSZXB1dGF0aW9uRW50aXR5RmlsdGVyS2V5EicKI1JFUFVUQVRJT05fRU5USVRZX0ZJTFRFUl'
    '9LRVlfU1RBVFVTEAASLAooUkVQVVRBVElPTl9FTlRJVFlfRklMVEVSX0tFWV9FTlRJVFlfVFlQ'
    'RRABEjgKNFJFUFVUQVRJT05fRU5USVRZX0ZJTFRFUl9LRVlfRU5USVRZX1JFRkVSRU5DRV9QUk'
    'VGSVgQAhIyCi5SRVBVVEFUSU9OX0VOVElUWV9GSUxURVJfS0VZX1JFUFVUQVRJT05fSU1QQUNU'
    'EAM=');

@$core.Deprecated('Use reputationEntityTypeDescriptor instead')
const ReputationEntityType$json = {
  '1': 'ReputationEntityType',
  '2': [
    {'1': 'REPUTATION_ENTITY_TYPE_RESOURCE', '2': 0},
  ],
};

/// Descriptor for `ReputationEntityType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reputationEntityTypeDescriptor = $convert.base64Decode(
    'ChRSZXB1dGF0aW9uRW50aXR5VHlwZRIjCh9SRVBVVEFUSU9OX0VOVElUWV9UWVBFX1JFU09VUk'
    'NFEAA=');

@$core.Deprecated('Use resourceTypeDescriptor instead')
const ResourceType$json = {
  '1': 'ResourceType',
  '2': [
    {'1': 'RESOURCE_TYPE_EMAIL_IDENTITY', '2': 0},
    {'1': 'RESOURCE_TYPE_CONFIGURATION_SET', '2': 1},
    {'1': 'RESOURCE_TYPE_EMAIL_TEMPLATE', '2': 2},
  ],
};

/// Descriptor for `ResourceType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List resourceTypeDescriptor = $convert.base64Decode(
    'CgxSZXNvdXJjZVR5cGUSIAocUkVTT1VSQ0VfVFlQRV9FTUFJTF9JREVOVElUWRAAEiMKH1JFU0'
    '9VUkNFX1RZUEVfQ09ORklHVVJBVElPTl9TRVQQARIgChxSRVNPVVJDRV9UWVBFX0VNQUlMX1RF'
    'TVBMQVRFEAI=');

@$core.Deprecated('Use reviewStatusDescriptor instead')
const ReviewStatus$json = {
  '1': 'ReviewStatus',
  '2': [
    {'1': 'REVIEW_STATUS_PENDING', '2': 0},
    {'1': 'REVIEW_STATUS_GRANTED', '2': 1},
    {'1': 'REVIEW_STATUS_FAILED', '2': 2},
    {'1': 'REVIEW_STATUS_DENIED', '2': 3},
  ],
};

/// Descriptor for `ReviewStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List reviewStatusDescriptor = $convert.base64Decode(
    'CgxSZXZpZXdTdGF0dXMSGQoVUkVWSUVXX1NUQVRVU19QRU5ESU5HEAASGQoVUkVWSUVXX1NUQV'
    'RVU19HUkFOVEVEEAESGAoUUkVWSUVXX1NUQVRVU19GQUlMRUQQAhIYChRSRVZJRVdfU1RBVFVT'
    'X0RFTklFRBAD');

@$core.Deprecated('Use scalingModeDescriptor instead')
const ScalingMode$json = {
  '1': 'ScalingMode',
  '2': [
    {'1': 'SCALING_MODE_STANDARD', '2': 0},
    {'1': 'SCALING_MODE_MANAGED', '2': 1},
  ],
};

/// Descriptor for `ScalingMode`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List scalingModeDescriptor = $convert.base64Decode(
    'CgtTY2FsaW5nTW9kZRIZChVTQ0FMSU5HX01PREVfU1RBTkRBUkQQABIYChRTQ0FMSU5HX01PRE'
    'VfTUFOQUdFRBAB');

@$core.Deprecated('Use sendingStatusDescriptor instead')
const SendingStatus$json = {
  '1': 'SendingStatus',
  '2': [
    {'1': 'SENDING_STATUS_DISABLED', '2': 0},
    {'1': 'SENDING_STATUS_REINSTATED', '2': 1},
    {'1': 'SENDING_STATUS_ENABLED', '2': 2},
  ],
};

/// Descriptor for `SendingStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List sendingStatusDescriptor = $convert.base64Decode(
    'Cg1TZW5kaW5nU3RhdHVzEhsKF1NFTkRJTkdfU1RBVFVTX0RJU0FCTEVEEAASHQoZU0VORElOR1'
    '9TVEFUVVNfUkVJTlNUQVRFRBABEhoKFlNFTkRJTkdfU1RBVFVTX0VOQUJMRUQQAg==');

@$core.Deprecated('Use statusDescriptor instead')
const Status$json = {
  '1': 'Status',
  '2': [
    {'1': 'STATUS_READY', '2': 0},
    {'1': 'STATUS_DELETING', '2': 1},
    {'1': 'STATUS_CREATING', '2': 2},
    {'1': 'STATUS_FAILED', '2': 3},
  ],
};

/// Descriptor for `Status`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List statusDescriptor = $convert.base64Decode(
    'CgZTdGF0dXMSEAoMU1RBVFVTX1JFQURZEAASEwoPU1RBVFVTX0RFTEVUSU5HEAESEwoPU1RBVF'
    'VTX0NSRUFUSU5HEAISEQoNU1RBVFVTX0ZBSUxFRBAD');

@$core.Deprecated('Use subscriptionStatusDescriptor instead')
const SubscriptionStatus$json = {
  '1': 'SubscriptionStatus',
  '2': [
    {'1': 'SUBSCRIPTION_STATUS_OPT_OUT', '2': 0},
    {'1': 'SUBSCRIPTION_STATUS_OPT_IN', '2': 1},
  ],
};

/// Descriptor for `SubscriptionStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List subscriptionStatusDescriptor = $convert.base64Decode(
    'ChJTdWJzY3JpcHRpb25TdGF0dXMSHwobU1VCU0NSSVBUSU9OX1NUQVRVU19PUFRfT1VUEAASHg'
    'oaU1VCU0NSSVBUSU9OX1NUQVRVU19PUFRfSU4QAQ==');

@$core.Deprecated('Use suppressionConfidenceVerdictThresholdDescriptor instead')
const SuppressionConfidenceVerdictThreshold$json = {
  '1': 'SuppressionConfidenceVerdictThreshold',
  '2': [
    {'1': 'SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MEDIUM', '2': 0},
    {'1': 'SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MANAGED', '2': 1},
    {'1': 'SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_HIGH', '2': 2},
  ],
};

/// Descriptor for `SuppressionConfidenceVerdictThreshold`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List suppressionConfidenceVerdictThresholdDescriptor =
    $convert.base64Decode(
        'CiVTdXBwcmVzc2lvbkNvbmZpZGVuY2VWZXJkaWN0VGhyZXNob2xkEjMKL1NVUFBSRVNTSU9OX0'
        'NPTkZJREVOQ0VfVkVSRElDVF9USFJFU0hPTERfTUVESVVNEAASNAowU1VQUFJFU1NJT05fQ09O'
        'RklERU5DRV9WRVJESUNUX1RIUkVTSE9MRF9NQU5BR0VEEAESMQotU1VQUFJFU1NJT05fQ09ORk'
        'lERU5DRV9WRVJESUNUX1RIUkVTSE9MRF9ISUdIEAI=');

@$core.Deprecated('Use suppressionListImportActionDescriptor instead')
const SuppressionListImportAction$json = {
  '1': 'SuppressionListImportAction',
  '2': [
    {'1': 'SUPPRESSION_LIST_IMPORT_ACTION_DELETE', '2': 0},
    {'1': 'SUPPRESSION_LIST_IMPORT_ACTION_PUT', '2': 1},
  ],
};

/// Descriptor for `SuppressionListImportAction`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List suppressionListImportActionDescriptor =
    $convert.base64Decode(
        'ChtTdXBwcmVzc2lvbkxpc3RJbXBvcnRBY3Rpb24SKQolU1VQUFJFU1NJT05fTElTVF9JTVBPUl'
        'RfQUNUSU9OX0RFTEVURRAAEiYKIlNVUFBSRVNTSU9OX0xJU1RfSU1QT1JUX0FDVElPTl9QVVQQ'
        'AQ==');

@$core.Deprecated('Use suppressionListReasonDescriptor instead')
const SuppressionListReason$json = {
  '1': 'SuppressionListReason',
  '2': [
    {'1': 'SUPPRESSION_LIST_REASON_BOUNCE', '2': 0},
    {'1': 'SUPPRESSION_LIST_REASON_COMPLAINT', '2': 1},
  ],
};

/// Descriptor for `SuppressionListReason`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List suppressionListReasonDescriptor = $convert.base64Decode(
    'ChVTdXBwcmVzc2lvbkxpc3RSZWFzb24SIgoeU1VQUFJFU1NJT05fTElTVF9SRUFTT05fQk9VTk'
    'NFEAASJQohU1VQUFJFU1NJT05fTElTVF9SRUFTT05fQ09NUExBSU5UEAE=');

@$core.Deprecated('Use tlsPolicyDescriptor instead')
const TlsPolicy$json = {
  '1': 'TlsPolicy',
  '2': [
    {'1': 'TLS_POLICY_OPTIONAL', '2': 0},
    {'1': 'TLS_POLICY_REQUIRE', '2': 1},
  ],
};

/// Descriptor for `TlsPolicy`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List tlsPolicyDescriptor = $convert.base64Decode(
    'CglUbHNQb2xpY3kSFwoTVExTX1BPTElDWV9PUFRJT05BTBAAEhYKElRMU19QT0xJQ1lfUkVRVU'
    'lSRRAB');

@$core.Deprecated('Use verificationErrorDescriptor instead')
const VerificationError$json = {
  '1': 'VerificationError',
  '2': [
    {
      '1': 'VERIFICATION_ERROR_REPLICATION_PRIMARY_BYO_DKIM_NOT_SUPPORTED',
      '2': 0
    },
    {'1': 'VERIFICATION_ERROR_REPLICATION_ACCESS_DENIED', '2': 1},
    {'1': 'VERIFICATION_ERROR_TYPE_NOT_FOUND', '2': 2},
    {'1': 'VERIFICATION_ERROR_DNS_SERVER_ERROR', '2': 3},
    {'1': 'VERIFICATION_ERROR_REPLICATION_PRIMARY_INVALID_REGION', '2': 4},
    {'1': 'VERIFICATION_ERROR_INVALID_VALUE', '2': 5},
    {'1': 'VERIFICATION_ERROR_REPLICATION_PRIMARY_NOT_FOUND', '2': 6},
    {'1': 'VERIFICATION_ERROR_HOST_NOT_FOUND', '2': 7},
    {
      '1': 'VERIFICATION_ERROR_REPLICATION_REPLICA_AS_PRIMARY_NOT_SUPPORTED',
      '2': 8
    },
    {'1': 'VERIFICATION_ERROR_SERVICE_ERROR', '2': 9},
  ],
};

/// Descriptor for `VerificationError`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List verificationErrorDescriptor = $convert.base64Decode(
    'ChFWZXJpZmljYXRpb25FcnJvchJBCj1WRVJJRklDQVRJT05fRVJST1JfUkVQTElDQVRJT05fUF'
    'JJTUFSWV9CWU9fREtJTV9OT1RfU1VQUE9SVEVEEAASMAosVkVSSUZJQ0FUSU9OX0VSUk9SX1JF'
    'UExJQ0FUSU9OX0FDQ0VTU19ERU5JRUQQARIlCiFWRVJJRklDQVRJT05fRVJST1JfVFlQRV9OT1'
    'RfRk9VTkQQAhInCiNWRVJJRklDQVRJT05fRVJST1JfRE5TX1NFUlZFUl9FUlJPUhADEjkKNVZF'
    'UklGSUNBVElPTl9FUlJPUl9SRVBMSUNBVElPTl9QUklNQVJZX0lOVkFMSURfUkVHSU9OEAQSJA'
    'ogVkVSSUZJQ0FUSU9OX0VSUk9SX0lOVkFMSURfVkFMVUUQBRI0CjBWRVJJRklDQVRJT05fRVJS'
    'T1JfUkVQTElDQVRJT05fUFJJTUFSWV9OT1RfRk9VTkQQBhIlCiFWRVJJRklDQVRJT05fRVJST1'
    'JfSE9TVF9OT1RfRk9VTkQQBxJDCj9WRVJJRklDQVRJT05fRVJST1JfUkVQTElDQVRJT05fUkVQ'
    'TElDQV9BU19QUklNQVJZX05PVF9TVVBQT1JURUQQCBIkCiBWRVJJRklDQVRJT05fRVJST1JfU0'
    'VSVklDRV9FUlJPUhAJ');

@$core.Deprecated('Use verificationStatusDescriptor instead')
const VerificationStatus$json = {
  '1': 'VerificationStatus',
  '2': [
    {'1': 'VERIFICATION_STATUS_PENDING', '2': 0},
    {'1': 'VERIFICATION_STATUS_SUCCESS', '2': 1},
    {'1': 'VERIFICATION_STATUS_TEMPORARY_FAILURE', '2': 2},
    {'1': 'VERIFICATION_STATUS_FAILED', '2': 3},
    {'1': 'VERIFICATION_STATUS_NOT_STARTED', '2': 4},
  ],
};

/// Descriptor for `VerificationStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List verificationStatusDescriptor = $convert.base64Decode(
    'ChJWZXJpZmljYXRpb25TdGF0dXMSHwobVkVSSUZJQ0FUSU9OX1NUQVRVU19QRU5ESU5HEAASHw'
    'obVkVSSUZJQ0FUSU9OX1NUQVRVU19TVUNDRVNTEAESKQolVkVSSUZJQ0FUSU9OX1NUQVRVU19U'
    'RU1QT1JBUllfRkFJTFVSRRACEh4KGlZFUklGSUNBVElPTl9TVEFUVVNfRkFJTEVEEAMSIwofVk'
    'VSSUZJQ0FUSU9OX1NUQVRVU19OT1RfU1RBUlRFRBAE');

@$core.Deprecated('Use warmupStatusDescriptor instead')
const WarmupStatus$json = {
  '1': 'WarmupStatus',
  '2': [
    {'1': 'WARMUP_STATUS_DONE', '2': 0},
    {'1': 'WARMUP_STATUS_NOT_APPLICABLE', '2': 1},
    {'1': 'WARMUP_STATUS_IN_PROGRESS', '2': 2},
  ],
};

/// Descriptor for `WarmupStatus`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List warmupStatusDescriptor = $convert.base64Decode(
    'CgxXYXJtdXBTdGF0dXMSFgoSV0FSTVVQX1NUQVRVU19ET05FEAASIAocV0FSTVVQX1NUQVRVU1'
    '9OT1RfQVBQTElDQUJMRRABEh0KGVdBUk1VUF9TVEFUVVNfSU5fUFJPR1JFU1MQAg==');

@$core.Deprecated('Use accountDetailsDescriptor instead')
const AccountDetails$json = {
  '1': 'AccountDetails',
  '2': [
    {
      '1': 'additionalcontactemailaddresses',
      '3': 322089693,
      '4': 3,
      '5': 9,
      '10': 'additionalcontactemailaddresses'
    },
    {
      '1': 'contactlanguage',
      '3': 114240022,
      '4': 1,
      '5': 14,
      '6': '.email.ContactLanguage',
      '10': 'contactlanguage'
    },
    {
      '1': 'mailtype',
      '3': 138144527,
      '4': 1,
      '5': 14,
      '6': '.email.MailType',
      '10': 'mailtype'
    },
    {
      '1': 'reviewdetails',
      '3': 378909498,
      '4': 1,
      '5': 11,
      '6': '.email.ReviewDetails',
      '10': 'reviewdetails'
    },
    {
      '1': 'usecasedescription',
      '3': 141053987,
      '4': 1,
      '5': 9,
      '10': 'usecasedescription'
    },
    {'1': 'websiteurl', '3': 201971828, '4': 1, '5': 9, '10': 'websiteurl'},
  ],
};

/// Descriptor for `AccountDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountDetailsDescriptor = $convert.base64Decode(
    'Cg5BY2NvdW50RGV0YWlscxJMCh9hZGRpdGlvbmFsY29udGFjdGVtYWlsYWRkcmVzc2VzGN3lyp'
    'kBIAMoCVIfYWRkaXRpb25hbGNvbnRhY3RlbWFpbGFkZHJlc3NlcxJDCg9jb250YWN0bGFuZ3Vh'
    'Z2UYltS8NiABKA4yFi5lbWFpbC5Db250YWN0TGFuZ3VhZ2VSD2NvbnRhY3RsYW5ndWFnZRIuCg'
    'htYWlsdHlwZRiP1u9BIAEoDjIPLmVtYWlsLk1haWxUeXBlUghtYWlsdHlwZRI+Cg1yZXZpZXdk'
    'ZXRhaWxzGLrm1rQBIAEoCzIULmVtYWlsLlJldmlld0RldGFpbHNSDXJldmlld2RldGFpbHMSMQ'
    'oSdXNlY2FzZWRlc2NyaXB0aW9uGKOgoUMgASgJUhJ1c2VjYXNlZGVzY3JpcHRpb24SIQoKd2Vi'
    'c2l0ZXVybBj0sKdgIAEoCVIKd2Vic2l0ZXVybA==');

@$core.Deprecated('Use accountSuspendedExceptionDescriptor instead')
const AccountSuspendedException$json = {
  '1': 'AccountSuspendedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AccountSuspendedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List accountSuspendedExceptionDescriptor =
    $convert.base64Decode(
        'ChlBY2NvdW50U3VzcGVuZGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use alreadyExistsExceptionDescriptor instead')
const AlreadyExistsException$json = {
  '1': 'AlreadyExistsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `AlreadyExistsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List alreadyExistsExceptionDescriptor =
    $convert.base64Decode(
        'ChZBbHJlYWR5RXhpc3RzRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use archivingOptionsDescriptor instead')
const ArchivingOptions$json = {
  '1': 'ArchivingOptions',
  '2': [
    {'1': 'archivearn', '3': 56866685, '4': 1, '5': 9, '10': 'archivearn'},
  ],
};

/// Descriptor for `ArchivingOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List archivingOptionsDescriptor = $convert.base64Decode(
    'ChBBcmNoaXZpbmdPcHRpb25zEiEKCmFyY2hpdmVhcm4Y/e6OGyABKAlSCmFyY2hpdmVhcm4=');

@$core.Deprecated('Use attachmentDescriptor instead')
const Attachment$json = {
  '1': 'Attachment',
  '2': [
    {
      '1': 'contentdescription',
      '3': 29266325,
      '4': 1,
      '5': 9,
      '10': 'contentdescription'
    },
    {
      '1': 'contentdisposition',
      '3': 120040130,
      '4': 1,
      '5': 14,
      '6': '.email.AttachmentContentDisposition',
      '10': 'contentdisposition'
    },
    {'1': 'contentid', '3': 431030096, '4': 1, '5': 9, '10': 'contentid'},
    {
      '1': 'contenttransferencoding',
      '3': 400458117,
      '4': 1,
      '5': 14,
      '6': '.email.AttachmentContentTransferEncoding',
      '10': 'contenttransferencoding'
    },
    {'1': 'contenttype', '3': 333064851, '4': 1, '5': 9, '10': 'contenttype'},
    {'1': 'filename', '3': 536729737, '4': 1, '5': 9, '10': 'filename'},
    {'1': 'rawcontent', '3': 230318007, '4': 1, '5': 12, '10': 'rawcontent'},
  ],
};

/// Descriptor for `Attachment`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List attachmentDescriptor = $convert.base64Decode(
    'CgpBdHRhY2htZW50EjEKEmNvbnRlbnRkZXNjcmlwdGlvbhiVo/oNIAEoCVISY29udGVudGRlc2'
    'NyaXB0aW9uElYKEmNvbnRlbnRkaXNwb3NpdGlvbhjC1Z45IAEoDjIjLmVtYWlsLkF0dGFjaG1l'
    'bnRDb250ZW50RGlzcG9zaXRpb25SEmNvbnRlbnRkaXNwb3NpdGlvbhIgCgljb250ZW50aWQY0P'
    '7DzQEgASgJUgljb250ZW50aWQSZgoXY29udGVudHRyYW5zZmVyZW5jb2RpbmcYhYP6vgEgASgO'
    'MiguZW1haWwuQXR0YWNobWVudENvbnRlbnRUcmFuc2ZlckVuY29kaW5nUhdjb250ZW50dHJhbn'
    'NmZXJlbmNvZGluZxIkCgtjb250ZW50dHlwZRiT1eieASABKAlSC2NvbnRlbnR0eXBlEh4KCGZp'
    'bGVuYW1lGImx9/8BIAEoCVIIZmlsZW5hbWUSIQoKcmF3Y29udGVudBi3v+ltIAEoDFIKcmF3Y2'
    '9udGVudA==');

@$core.Deprecated('Use badRequestExceptionDescriptor instead')
const BadRequestException$json = {
  '1': 'BadRequestException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `BadRequestException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List badRequestExceptionDescriptor =
    $convert.base64Decode(
        'ChNCYWRSZXF1ZXN0RXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use batchGetMetricDataQueryDescriptor instead')
const BatchGetMetricDataQuery$json = {
  '1': 'BatchGetMetricDataQuery',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.email.BatchGetMetricDataQuery.DimensionsEntry',
      '10': 'dimensions'
    },
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {
      '1': 'metric',
      '3': 386667298,
      '4': 1,
      '5': 14,
      '6': '.email.Metric',
      '10': 'metric'
    },
    {
      '1': 'namespace',
      '3': 355353153,
      '4': 1,
      '5': 14,
      '6': '.email.MetricNamespace',
      '10': 'namespace'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
  ],
  '3': [BatchGetMetricDataQuery_DimensionsEntry$json],
};

@$core.Deprecated('Use batchGetMetricDataQueryDescriptor instead')
const BatchGetMetricDataQuery_DimensionsEntry$json = {
  '1': 'DimensionsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `BatchGetMetricDataQuery`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetMetricDataQueryDescriptor = $convert.base64Decode(
    'ChdCYXRjaEdldE1ldHJpY0RhdGFRdWVyeRJSCgpkaW1lbnNpb25zGNGb39wBIAMoCzIuLmVtYW'
    'lsLkJhdGNoR2V0TWV0cmljRGF0YVF1ZXJ5LkRpbWVuc2lvbnNFbnRyeVIKZGltZW5zaW9ucxIb'
    'CgdlbmRkYXRlGM+z+SQgASgJUgdlbmRkYXRlEhIKAmlkGIHyorcBIAEoCVICaWQSKQoGbWV0cm'
    'ljGKKmsLgBIAEoDjINLmVtYWlsLk1ldHJpY1IGbWV0cmljEjgKCW5hbWVzcGFjZRjBhLmpASAB'
    'KA4yFi5lbWFpbC5NZXRyaWNOYW1lc3BhY2VSCW5hbWVzcGFjZRIgCglzdGFydGRhdGUY/Pig1A'
    'EgASgJUglzdGFydGRhdGUaPQoPRGltZW5zaW9uc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQK'
    'BXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use batchGetMetricDataRequestDescriptor instead')
const BatchGetMetricDataRequest$json = {
  '1': 'BatchGetMetricDataRequest',
  '2': [
    {
      '1': 'queries',
      '3': 305659620,
      '4': 3,
      '5': 11,
      '6': '.email.BatchGetMetricDataQuery',
      '10': 'queries'
    },
  ],
};

/// Descriptor for `BatchGetMetricDataRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetMetricDataRequestDescriptor =
    $convert.base64Decode(
        'ChlCYXRjaEdldE1ldHJpY0RhdGFSZXF1ZXN0EjwKB3F1ZXJpZXMY5P3fkQEgAygLMh4uZW1haW'
        'wuQmF0Y2hHZXRNZXRyaWNEYXRhUXVlcnlSB3F1ZXJpZXM=');

@$core.Deprecated('Use batchGetMetricDataResponseDescriptor instead')
const BatchGetMetricDataResponse$json = {
  '1': 'BatchGetMetricDataResponse',
  '2': [
    {
      '1': 'errors',
      '3': 166551719,
      '4': 3,
      '5': 11,
      '6': '.email.MetricDataError',
      '10': 'errors'
    },
    {
      '1': 'results',
      '3': 486024854,
      '4': 3,
      '5': 11,
      '6': '.email.MetricDataResult',
      '10': 'results'
    },
  ],
};

/// Descriptor for `BatchGetMetricDataResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List batchGetMetricDataResponseDescriptor =
    $convert.base64Decode(
        'ChpCYXRjaEdldE1ldHJpY0RhdGFSZXNwb25zZRIxCgZlcnJvcnMYp8G1TyADKAsyFi5lbWFpbC'
        '5NZXRyaWNEYXRhRXJyb3JSBmVycm9ycxI1CgdyZXN1bHRzGJbN4OcBIAMoCzIXLmVtYWlsLk1l'
        'dHJpY0RhdGFSZXN1bHRSB3Jlc3VsdHM=');

@$core.Deprecated('Use blacklistEntryDescriptor instead')
const BlacklistEntry$json = {
  '1': 'BlacklistEntry',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'listingtime', '3': 480893551, '4': 1, '5': 9, '10': 'listingtime'},
    {'1': 'rblname', '3': 317446625, '4': 1, '5': 9, '10': 'rblname'},
  ],
};

/// Descriptor for `BlacklistEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List blacklistEntryDescriptor = $convert.base64Decode(
    'Cg5CbGFja2xpc3RFbnRyeRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SJA'
    'oLbGlzdGluZ3RpbWUY77Sn5QEgASgJUgtsaXN0aW5ndGltZRIcCgdyYmxuYW1lGOGzr5cBIAEo'
    'CVIHcmJsbmFtZQ==');

@$core.Deprecated('Use bodyDescriptor instead')
const Body$json = {
  '1': 'Body',
  '2': [
    {
      '1': 'html',
      '3': 396489713,
      '4': 1,
      '5': 11,
      '6': '.email.Content',
      '10': 'html'
    },
    {
      '1': 'text',
      '3': 504638815,
      '4': 1,
      '5': 11,
      '6': '.email.Content',
      '10': 'text'
    },
  ],
};

/// Descriptor for `Body`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bodyDescriptor = $convert.base64Decode(
    'CgRCb2R5EiYKBGh0bWwY8eeHvQEgASgLMg4uZW1haWwuQ29udGVudFIEaHRtbBImCgR0ZXh0GN'
    '/a0PABIAEoCzIOLmVtYWlsLkNvbnRlbnRSBHRleHQ=');

@$core.Deprecated('Use bounceDescriptor instead')
const Bounce$json = {
  '1': 'Bounce',
  '2': [
    {
      '1': 'bouncesubtype',
      '3': 297220246,
      '4': 1,
      '5': 9,
      '10': 'bouncesubtype'
    },
    {
      '1': 'bouncetype',
      '3': 490222550,
      '4': 1,
      '5': 14,
      '6': '.email.BounceType',
      '10': 'bouncetype'
    },
    {
      '1': 'diagnosticcode',
      '3': 20043524,
      '4': 1,
      '5': 9,
      '10': 'diagnosticcode'
    },
  ],
};

/// Descriptor for `Bounce`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bounceDescriptor = $convert.base64Decode(
    'CgZCb3VuY2USKAoNYm91bmNlc3VidHlwZRiW8dyNASABKAlSDWJvdW5jZXN1YnR5cGUSNQoKYm'
    '91bmNldHlwZRjW5+DpASABKA4yES5lbWFpbC5Cb3VuY2VUeXBlUgpib3VuY2V0eXBlEikKDmRp'
    'YWdub3N0aWNjb2RlGISuxwkgASgJUg5kaWFnbm9zdGljY29kZQ==');

@$core.Deprecated('Use bulkEmailContentDescriptor instead')
const BulkEmailContent$json = {
  '1': 'BulkEmailContent',
  '2': [
    {
      '1': 'template',
      '3': 93743916,
      '4': 1,
      '5': 11,
      '6': '.email.Template',
      '10': 'template'
    },
  ],
};

/// Descriptor for `BulkEmailContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bulkEmailContentDescriptor = $convert.base64Decode(
    'ChBCdWxrRW1haWxDb250ZW50Ei4KCHRlbXBsYXRlGKzW2SwgASgLMg8uZW1haWwuVGVtcGxhdG'
    'VSCHRlbXBsYXRl');

@$core.Deprecated('Use bulkEmailEntryDescriptor instead')
const BulkEmailEntry$json = {
  '1': 'BulkEmailEntry',
  '2': [
    {
      '1': 'destination',
      '3': 457443680,
      '4': 1,
      '5': 11,
      '6': '.email.Destination',
      '10': 'destination'
    },
    {
      '1': 'replacementemailcontent',
      '3': 150255379,
      '4': 1,
      '5': 11,
      '6': '.email.ReplacementEmailContent',
      '10': 'replacementemailcontent'
    },
    {
      '1': 'replacementheaders',
      '3': 436072520,
      '4': 3,
      '5': 11,
      '6': '.email.MessageHeader',
      '10': 'replacementheaders'
    },
    {
      '1': 'replacementtags',
      '3': 379085467,
      '4': 3,
      '5': 11,
      '6': '.email.MessageTag',
      '10': 'replacementtags'
    },
  ],
};

/// Descriptor for `BulkEmailEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bulkEmailEntryDescriptor = $convert.base64Decode(
    'Cg5CdWxrRW1haWxFbnRyeRI4CgtkZXN0aW5hdGlvbhjgkpDaASABKAsyEi5lbWFpbC5EZXN0aW'
    '5hdGlvblILZGVzdGluYXRpb24SWwoXcmVwbGFjZW1lbnRlbWFpbGNvbnRlbnQYk+7SRyABKAsy'
    'Hi5lbWFpbC5SZXBsYWNlbWVudEVtYWlsQ29udGVudFIXcmVwbGFjZW1lbnRlbWFpbGNvbnRlbn'
    'QSSAoScmVwbGFjZW1lbnRoZWFkZXJzGMjg988BIAMoCzIULmVtYWlsLk1lc3NhZ2VIZWFkZXJS'
    'EnJlcGxhY2VtZW50aGVhZGVycxI/Cg9yZXBsYWNlbWVudHRhZ3MYm8XhtAEgAygLMhEuZW1haW'
    'wuTWVzc2FnZVRhZ1IPcmVwbGFjZW1lbnR0YWdz');

@$core.Deprecated('Use bulkEmailEntryResultDescriptor instead')
const BulkEmailEntryResult$json = {
  '1': 'BulkEmailEntryResult',
  '2': [
    {'1': 'error', '3': 328047858, '4': 1, '5': 9, '10': 'error'},
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.BulkEmailStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `BulkEmailEntryResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List bulkEmailEntryResultDescriptor = $convert.base64Decode(
    'ChRCdWxrRW1haWxFbnRyeVJlc3VsdBIYCgVlcnJvchjyubacASABKAlSBWVycm9yEiAKCW1lc3'
    'NhZ2VpZBiq5vSrASABKAlSCW1lc3NhZ2VpZBIxCgZzdGF0dXMYkOT7AiABKA4yFi5lbWFpbC5C'
    'dWxrRW1haWxTdGF0dXNSBnN0YXR1cw==');

@$core.Deprecated('Use cancelExportJobRequestDescriptor instead')
const CancelExportJobRequest$json = {
  '1': 'CancelExportJobRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
  ],
};

/// Descriptor for `CancelExportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelExportJobRequestDescriptor =
    $convert.base64Decode(
        'ChZDYW5jZWxFeHBvcnRKb2JSZXF1ZXN0EhcKBWpvYmlkGNLU3TMgASgJUgVqb2JpZA==');

@$core.Deprecated('Use cancelExportJobResponseDescriptor instead')
const CancelExportJobResponse$json = {
  '1': 'CancelExportJobResponse',
};

/// Descriptor for `CancelExportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelExportJobResponseDescriptor =
    $convert.base64Decode('ChdDYW5jZWxFeHBvcnRKb2JSZXNwb25zZQ==');

@$core.Deprecated('Use cloudWatchDestinationDescriptor instead')
const CloudWatchDestination$json = {
  '1': 'CloudWatchDestination',
  '2': [
    {
      '1': 'dimensionconfigurations',
      '3': 283186905,
      '4': 3,
      '5': 11,
      '6': '.email.CloudWatchDimensionConfiguration',
      '10': 'dimensionconfigurations'
    },
  ],
};

/// Descriptor for `CloudWatchDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchDestinationDescriptor = $convert.base64Decode(
    'ChVDbG91ZFdhdGNoRGVzdGluYXRpb24SZQoXZGltZW5zaW9uY29uZmlndXJhdGlvbnMY2a2Ehw'
    'EgAygLMicuZW1haWwuQ2xvdWRXYXRjaERpbWVuc2lvbkNvbmZpZ3VyYXRpb25SF2RpbWVuc2lv'
    'bmNvbmZpZ3VyYXRpb25z');

@$core.Deprecated('Use cloudWatchDimensionConfigurationDescriptor instead')
const CloudWatchDimensionConfiguration$json = {
  '1': 'CloudWatchDimensionConfiguration',
  '2': [
    {
      '1': 'defaultdimensionvalue',
      '3': 428650160,
      '4': 1,
      '5': 9,
      '10': 'defaultdimensionvalue'
    },
    {
      '1': 'dimensionname',
      '3': 324514281,
      '4': 1,
      '5': 9,
      '10': 'dimensionname'
    },
    {
      '1': 'dimensionvaluesource',
      '3': 406297926,
      '4': 1,
      '5': 14,
      '6': '.email.DimensionValueSource',
      '10': 'dimensionvaluesource'
    },
  ],
};

/// Descriptor for `CloudWatchDimensionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudWatchDimensionConfigurationDescriptor =
    $convert.base64Decode(
        'CiBDbG91ZFdhdGNoRGltZW5zaW9uQ29uZmlndXJhdGlvbhI4ChVkZWZhdWx0ZGltZW5zaW9udm'
        'FsdWUYsN2yzAEgASgJUhVkZWZhdWx0ZGltZW5zaW9udmFsdWUSKAoNZGltZW5zaW9ubmFtZRjp'
        '496aASABKAlSDWRpbWVuc2lvbm5hbWUSUwoUZGltZW5zaW9udmFsdWVzb3VyY2UYxrrewQEgAS'
        'gOMhsuZW1haWwuRGltZW5zaW9uVmFsdWVTb3VyY2VSFGRpbWVuc2lvbnZhbHVlc291cmNl');

@$core.Deprecated('Use complaintDescriptor instead')
const Complaint$json = {
  '1': 'Complaint',
  '2': [
    {
      '1': 'complaintfeedbacktype',
      '3': 251430634,
      '4': 1,
      '5': 9,
      '10': 'complaintfeedbacktype'
    },
    {
      '1': 'complaintsubtype',
      '3': 351511425,
      '4': 1,
      '5': 9,
      '10': 'complaintsubtype'
    },
  ],
};

/// Descriptor for `Complaint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List complaintDescriptor = $convert.base64Decode(
    'CglDb21wbGFpbnQSNwoVY29tcGxhaW50ZmVlZGJhY2t0eXBlGOqN8ncgASgJUhVjb21wbGFpbn'
    'RmZWVkYmFja3R5cGUSLgoQY29tcGxhaW50c3VidHlwZRiBx86nASABKAlSEGNvbXBsYWludHN1'
    'YnR5cGU=');

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

@$core.Deprecated('Use conflictExceptionDescriptor instead')
const ConflictException$json = {
  '1': 'ConflictException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictExceptionDescriptor = $convert.base64Decode(
    'ChFDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use contactDescriptor instead')
const Contact$json = {
  '1': 'Contact',
  '2': [
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {
      '1': 'topicdefaultpreferences',
      '3': 477710170,
      '4': 3,
      '5': 11,
      '6': '.email.TopicPreference',
      '10': 'topicdefaultpreferences'
    },
    {
      '1': 'topicpreferences',
      '3': 199822983,
      '4': 3,
      '5': 11,
      '6': '.email.TopicPreference',
      '10': 'topicpreferences'
    },
    {
      '1': 'unsubscribeall',
      '3': 49261174,
      '4': 1,
      '5': 8,
      '10': 'unsubscribeall'
    },
  ],
};

/// Descriptor for `Contact`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contactDescriptor = $convert.base64Decode(
    'CgdDb250YWN0EiYKDGVtYWlsYWRkcmVzcxjW8YrpASABKAlSDGVtYWlsYWRkcmVzcxI1ChRsYX'
    'N0dXBkYXRlZHRpbWVzdGFtcBiVy8g/IAEoCVIUbGFzdHVwZGF0ZWR0aW1lc3RhbXASVAoXdG9w'
    'aWNkZWZhdWx0cHJlZmVyZW5jZXMY2o7l4wEgAygLMhYuZW1haWwuVG9waWNQcmVmZXJlbmNlUh'
    'd0b3BpY2RlZmF1bHRwcmVmZXJlbmNlcxJFChB0b3BpY3ByZWZlcmVuY2VzGIedpF8gAygLMhYu'
    'ZW1haWwuVG9waWNQcmVmZXJlbmNlUhB0b3BpY3ByZWZlcmVuY2VzEikKDnVuc3Vic2NyaWJlYW'
    'xsGPbUvhcgASgIUg51bnN1YnNjcmliZWFsbA==');

@$core.Deprecated('Use contactListDescriptor instead')
const ContactList$json = {
  '1': 'ContactList',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
  ],
};

/// Descriptor for `ContactList`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contactListDescriptor = $convert.base64Decode(
    'CgtDb250YWN0TGlzdBIsCg9jb250YWN0bGlzdG5hbWUY8b2noQEgASgJUg9jb250YWN0bGlzdG'
    '5hbWUSNQoUbGFzdHVwZGF0ZWR0aW1lc3RhbXAYlcvIPyABKAlSFGxhc3R1cGRhdGVkdGltZXN0'
    'YW1w');

@$core.Deprecated('Use contactListDestinationDescriptor instead')
const ContactListDestination$json = {
  '1': 'ContactListDestination',
  '2': [
    {
      '1': 'contactlistimportaction',
      '3': 247602529,
      '4': 1,
      '5': 14,
      '6': '.email.ContactListImportAction',
      '10': 'contactlistimportaction'
    },
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
  ],
};

/// Descriptor for `ContactListDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contactListDestinationDescriptor = $convert.base64Decode(
    'ChZDb250YWN0TGlzdERlc3RpbmF0aW9uElsKF2NvbnRhY3RsaXN0aW1wb3J0YWN0aW9uGOG6iH'
    'YgASgOMh4uZW1haWwuQ29udGFjdExpc3RJbXBvcnRBY3Rpb25SF2NvbnRhY3RsaXN0aW1wb3J0'
    'YWN0aW9uEiwKD2NvbnRhY3RsaXN0bmFtZRjxvaehASABKAlSD2NvbnRhY3RsaXN0bmFtZQ==');

@$core.Deprecated('Use contentDescriptor instead')
const Content$json = {
  '1': 'Content',
  '2': [
    {'1': 'charset', '3': 348272428, '4': 1, '5': 9, '10': 'charset'},
    {'1': 'data', '3': 525498822, '4': 1, '5': 9, '10': 'data'},
  ],
};

/// Descriptor for `Content`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List contentDescriptor = $convert.base64Decode(
    'CgdDb250ZW50EhwKB2NoYXJzZXQYrO6IpgEgASgJUgdjaGFyc2V0EhYKBGRhdGEYxvPJ+gEgAS'
    'gJUgRkYXRh');

@$core.Deprecated(
    'Use createConfigurationSetEventDestinationRequestDescriptor instead')
const CreateConfigurationSetEventDestinationRequest$json = {
  '1': 'CreateConfigurationSetEventDestinationRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'eventdestination',
      '3': 469882902,
      '4': 1,
      '5': 11,
      '6': '.email.EventDestinationDefinition',
      '10': 'eventdestination'
    },
    {
      '1': 'eventdestinationname',
      '3': 477517655,
      '4': 1,
      '5': 9,
      '10': 'eventdestinationname'
    },
  ],
};

/// Descriptor for `CreateConfigurationSetEventDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createConfigurationSetEventDestinationRequestDescriptor =
    $convert.base64Decode(
        'Ci1DcmVhdGVDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvblJlcXVlc3QSNgoUY29uZm'
        'lndXJhdGlvbnNldG5hbWUYzYuxwAEgASgJUhRjb25maWd1cmF0aW9uc2V0bmFtZRJRChBldmVu'
        'dGRlc3RpbmF0aW9uGJawh+ABIAEoCzIhLmVtYWlsLkV2ZW50RGVzdGluYXRpb25EZWZpbml0aW'
        '9uUhBldmVudGRlc3RpbmF0aW9uEjYKFGV2ZW50ZGVzdGluYXRpb25uYW1lGNeu2eMBIAEoCVIU'
        'ZXZlbnRkZXN0aW5hdGlvbm5hbWU=');

@$core.Deprecated(
    'Use createConfigurationSetEventDestinationResponseDescriptor instead')
const CreateConfigurationSetEventDestinationResponse$json = {
  '1': 'CreateConfigurationSetEventDestinationResponse',
};

/// Descriptor for `CreateConfigurationSetEventDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createConfigurationSetEventDestinationResponseDescriptor =
    $convert.base64Decode(
        'Ci5DcmVhdGVDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvblJlc3BvbnNl');

@$core.Deprecated('Use createConfigurationSetRequestDescriptor instead')
const CreateConfigurationSetRequest$json = {
  '1': 'CreateConfigurationSetRequest',
  '2': [
    {
      '1': 'archivingoptions',
      '3': 54783281,
      '4': 1,
      '5': 11,
      '6': '.email.ArchivingOptions',
      '10': 'archivingoptions'
    },
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'deliveryoptions',
      '3': 332764278,
      '4': 1,
      '5': 11,
      '6': '.email.DeliveryOptions',
      '10': 'deliveryoptions'
    },
    {
      '1': 'reputationoptions',
      '3': 98302085,
      '4': 1,
      '5': 11,
      '6': '.email.ReputationOptions',
      '10': 'reputationoptions'
    },
    {
      '1': 'sendingoptions',
      '3': 319193106,
      '4': 1,
      '5': 11,
      '6': '.email.SendingOptions',
      '10': 'sendingoptions'
    },
    {
      '1': 'suppressionoptions',
      '3': 96075559,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionOptions',
      '10': 'suppressionoptions'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'trackingoptions',
      '3': 237447811,
      '4': 1,
      '5': 11,
      '6': '.email.TrackingOptions',
      '10': 'trackingoptions'
    },
    {
      '1': 'vdmoptions',
      '3': 344971073,
      '4': 1,
      '5': 11,
      '6': '.email.VdmOptions',
      '10': 'vdmoptions'
    },
  ],
};

/// Descriptor for `CreateConfigurationSetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConfigurationSetRequestDescriptor = $convert.base64Decode(
    'Ch1DcmVhdGVDb25maWd1cmF0aW9uU2V0UmVxdWVzdBJGChBhcmNoaXZpbmdvcHRpb25zGLHajx'
    'ogASgLMhcuZW1haWwuQXJjaGl2aW5nT3B0aW9uc1IQYXJjaGl2aW5nb3B0aW9ucxI2ChRjb25m'
    'aWd1cmF0aW9uc2V0bmFtZRjNi7HAASABKAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1lEkQKD2RlbG'
    'l2ZXJ5b3B0aW9ucxj2qNaeASABKAsyFi5lbWFpbC5EZWxpdmVyeU9wdGlvbnNSD2RlbGl2ZXJ5'
    'b3B0aW9ucxJJChFyZXB1dGF0aW9ub3B0aW9ucxiF8e8uIAEoCzIYLmVtYWlsLlJlcHV0YXRpb2'
    '5PcHRpb25zUhFyZXB1dGF0aW9ub3B0aW9ucxJBCg5zZW5kaW5nb3B0aW9ucxiSgJqYASABKAsy'
    'FS5lbWFpbC5TZW5kaW5nT3B0aW9uc1IOc2VuZGluZ29wdGlvbnMSTAoSc3VwcHJlc3Npb25vcH'
    'Rpb25zGKf+5y0gASgLMhkuZW1haWwuU3VwcHJlc3Npb25PcHRpb25zUhJzdXBwcmVzc2lvbm9w'
    'dGlvbnMSIgoEdGFncxjBwfa1ASADKAsyCi5lbWFpbC5UYWdSBHRhZ3MSQwoPdHJhY2tpbmdvcH'
    'Rpb25zGIPVnHEgASgLMhYuZW1haWwuVHJhY2tpbmdPcHRpb25zUg90cmFja2luZ29wdGlvbnMS'
    'NQoKdmRtb3B0aW9ucxjBrr+kASABKAsyES5lbWFpbC5WZG1PcHRpb25zUgp2ZG1vcHRpb25z');

@$core.Deprecated('Use createConfigurationSetResponseDescriptor instead')
const CreateConfigurationSetResponse$json = {
  '1': 'CreateConfigurationSetResponse',
};

/// Descriptor for `CreateConfigurationSetResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createConfigurationSetResponseDescriptor =
    $convert.base64Decode('Ch5DcmVhdGVDb25maWd1cmF0aW9uU2V0UmVzcG9uc2U=');

@$core.Deprecated('Use createContactListRequestDescriptor instead')
const CreateContactListRequest$json = {
  '1': 'CreateContactListRequest',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'topics',
      '3': 219850038,
      '4': 3,
      '5': 11,
      '6': '.email.Topic',
      '10': 'topics'
    },
  ],
};

/// Descriptor for `CreateContactListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactListRequestDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVDb250YWN0TGlzdFJlcXVlc3QSLAoPY29udGFjdGxpc3RuYW1lGPG9p6EBIAEoCV'
    'IPY29udGFjdGxpc3RuYW1lEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIi'
    'CgR0YWdzGMHB9rUBIAMoCzIKLmVtYWlsLlRhZ1IEdGFncxInCgZ0b3BpY3MYtsrqaCADKAsyDC'
    '5lbWFpbC5Ub3BpY1IGdG9waWNz');

@$core.Deprecated('Use createContactListResponseDescriptor instead')
const CreateContactListResponse$json = {
  '1': 'CreateContactListResponse',
};

/// Descriptor for `CreateContactListResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactListResponseDescriptor =
    $convert.base64Decode('ChlDcmVhdGVDb250YWN0TGlzdFJlc3BvbnNl');

@$core.Deprecated('Use createContactRequestDescriptor instead')
const CreateContactRequest$json = {
  '1': 'CreateContactRequest',
  '2': [
    {
      '1': 'attributesdata',
      '3': 497271421,
      '4': 1,
      '5': 9,
      '10': 'attributesdata'
    },
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {
      '1': 'topicpreferences',
      '3': 199822983,
      '4': 3,
      '5': 11,
      '6': '.email.TopicPreference',
      '10': 'topicpreferences'
    },
    {
      '1': 'unsubscribeall',
      '3': 49261174,
      '4': 1,
      '5': 8,
      '10': 'unsubscribeall'
    },
  ],
};

/// Descriptor for `CreateContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactRequestDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVDb250YWN0UmVxdWVzdBIqCg5hdHRyaWJ1dGVzZGF0YRj9hI/tASABKAlSDmF0dH'
    'JpYnV0ZXNkYXRhEiwKD2NvbnRhY3RsaXN0bmFtZRjxvaehASABKAlSD2NvbnRhY3RsaXN0bmFt'
    'ZRImCgxlbWFpbGFkZHJlc3MY1vGK6QEgASgJUgxlbWFpbGFkZHJlc3MSRQoQdG9waWNwcmVmZX'
    'JlbmNlcxiHnaRfIAMoCzIWLmVtYWlsLlRvcGljUHJlZmVyZW5jZVIQdG9waWNwcmVmZXJlbmNl'
    'cxIpCg51bnN1YnNjcmliZWFsbBj21L4XIAEoCFIOdW5zdWJzY3JpYmVhbGw=');

@$core.Deprecated('Use createContactResponseDescriptor instead')
const CreateContactResponse$json = {
  '1': 'CreateContactResponse',
};

/// Descriptor for `CreateContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createContactResponseDescriptor =
    $convert.base64Decode('ChVDcmVhdGVDb250YWN0UmVzcG9uc2U=');

@$core.Deprecated(
    'Use createCustomVerificationEmailTemplateRequestDescriptor instead')
const CreateCustomVerificationEmailTemplateRequest$json = {
  '1': 'CreateCustomVerificationEmailTemplateRequest',
  '2': [
    {
      '1': 'failureredirectionurl',
      '3': 376073007,
      '4': 1,
      '5': 9,
      '10': 'failureredirectionurl'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {
      '1': 'successredirectionurl',
      '3': 513750768,
      '4': 1,
      '5': 9,
      '10': 'successredirectionurl'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'templatecontent',
      '3': 528973453,
      '4': 1,
      '5': 9,
      '10': 'templatecontent'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
    {
      '1': 'templatesubject',
      '3': 144653738,
      '4': 1,
      '5': 9,
      '10': 'templatesubject'
    },
  ],
};

/// Descriptor for `CreateCustomVerificationEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createCustomVerificationEmailTemplateRequestDescriptor =
    $convert.base64Decode(
        'CixDcmVhdGVDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVxdWVzdBI4ChVmYWlsdX'
        'JlcmVkaXJlY3Rpb251cmwYr9apswEgASgJUhVmYWlsdXJlcmVkaXJlY3Rpb251cmwSLQoQZnJv'
        'bWVtYWlsYWRkcmVzcxiGmsssIAEoCVIQZnJvbWVtYWlsYWRkcmVzcxI4ChVzdWNjZXNzcmVkaX'
        'JlY3Rpb251cmwY8O389AEgASgJUhVzdWNjZXNzcmVkaXJlY3Rpb251cmwSIgoEdGFncxjBwfa1'
        'ASADKAsyCi5lbWFpbC5UYWdSBHRhZ3MSLAoPdGVtcGxhdGVjb250ZW50GI39nfwBIAEoCVIPdG'
        'VtcGxhdGVjb250ZW50EiYKDHRlbXBsYXRlbmFtZRixmJHlASABKAlSDHRlbXBsYXRlbmFtZRIr'
        'Cg90ZW1wbGF0ZXN1YmplY3QYqvv8RCABKAlSD3RlbXBsYXRlc3ViamVjdA==');

@$core.Deprecated(
    'Use createCustomVerificationEmailTemplateResponseDescriptor instead')
const CreateCustomVerificationEmailTemplateResponse$json = {
  '1': 'CreateCustomVerificationEmailTemplateResponse',
};

/// Descriptor for `CreateCustomVerificationEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    createCustomVerificationEmailTemplateResponseDescriptor =
    $convert.base64Decode(
        'Ci1DcmVhdGVDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVzcG9uc2U=');

@$core.Deprecated('Use createDedicatedIpPoolRequestDescriptor instead')
const CreateDedicatedIpPoolRequest$json = {
  '1': 'CreateDedicatedIpPoolRequest',
  '2': [
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
    {
      '1': 'scalingmode',
      '3': 210356138,
      '4': 1,
      '5': 14,
      '6': '.email.ScalingMode',
      '10': 'scalingmode'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateDedicatedIpPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDedicatedIpPoolRequestDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVEZWRpY2F0ZWRJcFBvb2xSZXF1ZXN0Eh0KCHBvb2xuYW1lGMmNhScgASgJUghwb2'
        '9sbmFtZRI3CgtzY2FsaW5nbW9kZRiqj6dkIAEoDjISLmVtYWlsLlNjYWxpbmdNb2RlUgtzY2Fs'
        'aW5nbW9kZRIiCgR0YWdzGMHB9rUBIAMoCzIKLmVtYWlsLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createDedicatedIpPoolResponseDescriptor instead')
const CreateDedicatedIpPoolResponse$json = {
  '1': 'CreateDedicatedIpPoolResponse',
};

/// Descriptor for `CreateDedicatedIpPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDedicatedIpPoolResponseDescriptor =
    $convert.base64Decode('Ch1DcmVhdGVEZWRpY2F0ZWRJcFBvb2xSZXNwb25zZQ==');

@$core.Deprecated('Use createDeliverabilityTestReportRequestDescriptor instead')
const CreateDeliverabilityTestReportRequest$json = {
  '1': 'CreateDeliverabilityTestReportRequest',
  '2': [
    {
      '1': 'content',
      '3': 23568227,
      '4': 1,
      '5': 11,
      '6': '.email.EmailContent',
      '10': 'content'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {'1': 'reportname', '3': 526054737, '4': 1, '5': 9, '10': 'reportname'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateDeliverabilityTestReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDeliverabilityTestReportRequestDescriptor =
    $convert.base64Decode(
        'CiVDcmVhdGVEZWxpdmVyYWJpbGl0eVRlc3RSZXBvcnRSZXF1ZXN0EjAKB2NvbnRlbnQY476eCy'
        'ABKAsyEy5lbWFpbC5FbWFpbENvbnRlbnRSB2NvbnRlbnQSLQoQZnJvbWVtYWlsYWRkcmVzcxiG'
        'msssIAEoCVIQZnJvbWVtYWlsYWRkcmVzcxIiCgpyZXBvcnRuYW1lGNHq6/oBIAEoCVIKcmVwb3'
        'J0bmFtZRIiCgR0YWdzGMHB9rUBIAMoCzIKLmVtYWlsLlRhZ1IEdGFncw==');

@$core
    .Deprecated('Use createDeliverabilityTestReportResponseDescriptor instead')
const CreateDeliverabilityTestReportResponse$json = {
  '1': 'CreateDeliverabilityTestReportResponse',
  '2': [
    {
      '1': 'deliverabilityteststatus',
      '3': 71311387,
      '4': 1,
      '5': 14,
      '6': '.email.DeliverabilityTestStatus',
      '10': 'deliverabilityteststatus'
    },
    {'1': 'reportid', '3': 420903847, '4': 1, '5': 9, '10': 'reportid'},
  ],
};

/// Descriptor for `CreateDeliverabilityTestReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createDeliverabilityTestReportResponseDescriptor =
    $convert.base64Decode(
        'CiZDcmVhdGVEZWxpdmVyYWJpbGl0eVRlc3RSZXBvcnRSZXNwb25zZRJeChhkZWxpdmVyYWJpbG'
        'l0eXRlc3RzdGF0dXMYm8CAIiABKA4yHy5lbWFpbC5EZWxpdmVyYWJpbGl0eVRlc3RTdGF0dXNS'
        'GGRlbGl2ZXJhYmlsaXR5dGVzdHN0YXR1cxIeCghyZXBvcnRpZBin99nIASABKAlSCHJlcG9ydG'
        'lk');

@$core.Deprecated('Use createEmailIdentityPolicyRequestDescriptor instead')
const CreateEmailIdentityPolicyRequest$json = {
  '1': 'CreateEmailIdentityPolicyRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `CreateEmailIdentityPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEmailIdentityPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVFbWFpbElkZW50aXR5UG9saWN5UmVxdWVzdBInCg1lbWFpbGlkZW50aXR5GJqU8k'
        'AgASgJUg1lbWFpbGlkZW50aXR5EhoKBnBvbGljeRig7/DgASABKAlSBnBvbGljeRIhCgpwb2xp'
        'Y3luYW1lGL31h38gASgJUgpwb2xpY3luYW1l');

@$core.Deprecated('Use createEmailIdentityPolicyResponseDescriptor instead')
const CreateEmailIdentityPolicyResponse$json = {
  '1': 'CreateEmailIdentityPolicyResponse',
};

/// Descriptor for `CreateEmailIdentityPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEmailIdentityPolicyResponseDescriptor =
    $convert.base64Decode('CiFDcmVhdGVFbWFpbElkZW50aXR5UG9saWN5UmVzcG9uc2U=');

@$core.Deprecated('Use createEmailIdentityRequestDescriptor instead')
const CreateEmailIdentityRequest$json = {
  '1': 'CreateEmailIdentityRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'dkimsigningattributes',
      '3': 81684873,
      '4': 1,
      '5': 11,
      '6': '.email.DkimSigningAttributes',
      '10': 'dkimsigningattributes'
    },
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateEmailIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEmailIdentityRequestDescriptor = $convert.base64Decode(
    'ChpDcmVhdGVFbWFpbElkZW50aXR5UmVxdWVzdBI2ChRjb25maWd1cmF0aW9uc2V0bmFtZRjNi7'
    'HAASABKAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1lElUKFWRraW1zaWduaW5nYXR0cmlidXRlcxiJ'
    '0/kmIAEoCzIcLmVtYWlsLkRraW1TaWduaW5nQXR0cmlidXRlc1IVZGtpbXNpZ25pbmdhdHRyaW'
    'J1dGVzEicKDWVtYWlsaWRlbnRpdHkYmpTyQCABKAlSDWVtYWlsaWRlbnRpdHkSIgoEdGFncxjB'
    'wfa1ASADKAsyCi5lbWFpbC5UYWdSBHRhZ3M=');

@$core.Deprecated('Use createEmailIdentityResponseDescriptor instead')
const CreateEmailIdentityResponse$json = {
  '1': 'CreateEmailIdentityResponse',
  '2': [
    {
      '1': 'dkimattributes',
      '3': 256039632,
      '4': 1,
      '5': 11,
      '6': '.email.DkimAttributes',
      '10': 'dkimattributes'
    },
    {
      '1': 'identitytype',
      '3': 499274628,
      '4': 1,
      '5': 14,
      '6': '.email.IdentityType',
      '10': 'identitytype'
    },
    {
      '1': 'verifiedforsendingstatus',
      '3': 163100765,
      '4': 1,
      '5': 8,
      '10': 'verifiedforsendingstatus'
    },
  ],
};

/// Descriptor for `CreateEmailIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEmailIdentityResponseDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVFbWFpbElkZW50aXR5UmVzcG9uc2USQAoOZGtpbWF0dHJpYnV0ZXMY0LWLeiABKA'
    'syFS5lbWFpbC5Ea2ltQXR0cmlidXRlc1IOZGtpbWF0dHJpYnV0ZXMSOwoMaWRlbnRpdHl0eXBl'
    'GISnie4BIAEoDjITLmVtYWlsLklkZW50aXR5VHlwZVIMaWRlbnRpdHl0eXBlEj0KGHZlcmlmaW'
    'VkZm9yc2VuZGluZ3N0YXR1cxjd8OJNIAEoCFIYdmVyaWZpZWRmb3JzZW5kaW5nc3RhdHVz');

@$core.Deprecated('Use createEmailTemplateRequestDescriptor instead')
const CreateEmailTemplateRequest$json = {
  '1': 'CreateEmailTemplateRequest',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'templatecontent',
      '3': 528973453,
      '4': 1,
      '5': 11,
      '6': '.email.EmailTemplateContent',
      '10': 'templatecontent'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `CreateEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEmailTemplateRequestDescriptor = $convert.base64Decode(
    'ChpDcmVhdGVFbWFpbFRlbXBsYXRlUmVxdWVzdBIiCgR0YWdzGMHB9rUBIAMoCzIKLmVtYWlsLl'
    'RhZ1IEdGFncxJJCg90ZW1wbGF0ZWNvbnRlbnQYjf2d/AEgASgLMhsuZW1haWwuRW1haWxUZW1w'
    'bGF0ZUNvbnRlbnRSD3RlbXBsYXRlY29udGVudBImCgx0ZW1wbGF0ZW5hbWUYsZiR5QEgASgJUg'
    'x0ZW1wbGF0ZW5hbWU=');

@$core.Deprecated('Use createEmailTemplateResponseDescriptor instead')
const CreateEmailTemplateResponse$json = {
  '1': 'CreateEmailTemplateResponse',
};

/// Descriptor for `CreateEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createEmailTemplateResponseDescriptor =
    $convert.base64Decode('ChtDcmVhdGVFbWFpbFRlbXBsYXRlUmVzcG9uc2U=');

@$core.Deprecated('Use createExportJobRequestDescriptor instead')
const CreateExportJobRequest$json = {
  '1': 'CreateExportJobRequest',
  '2': [
    {
      '1': 'exportdatasource',
      '3': 308051423,
      '4': 1,
      '5': 11,
      '6': '.email.ExportDataSource',
      '10': 'exportdatasource'
    },
    {
      '1': 'exportdestination',
      '3': 523408618,
      '4': 1,
      '5': 11,
      '6': '.email.ExportDestination',
      '10': 'exportdestination'
    },
  ],
};

/// Descriptor for `CreateExportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createExportJobRequestDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVFeHBvcnRKb2JSZXF1ZXN0EkcKEGV4cG9ydGRhdGFzb3VyY2UY3/vxkgEgASgLMh'
    'cuZW1haWwuRXhwb3J0RGF0YVNvdXJjZVIQZXhwb3J0ZGF0YXNvdXJjZRJKChFleHBvcnRkZXN0'
    'aW5hdGlvbhjqqcr5ASABKAsyGC5lbWFpbC5FeHBvcnREZXN0aW5hdGlvblIRZXhwb3J0ZGVzdG'
    'luYXRpb24=');

@$core.Deprecated('Use createExportJobResponseDescriptor instead')
const CreateExportJobResponse$json = {
  '1': 'CreateExportJobResponse',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
  ],
};

/// Descriptor for `CreateExportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createExportJobResponseDescriptor =
    $convert.base64Decode(
        'ChdDcmVhdGVFeHBvcnRKb2JSZXNwb25zZRIXCgVqb2JpZBjS1N0zIAEoCVIFam9iaWQ=');

@$core.Deprecated('Use createImportJobRequestDescriptor instead')
const CreateImportJobRequest$json = {
  '1': 'CreateImportJobRequest',
  '2': [
    {
      '1': 'importdatasource',
      '3': 486006026,
      '4': 1,
      '5': 11,
      '6': '.email.ImportDataSource',
      '10': 'importdatasource'
    },
    {
      '1': 'importdestination',
      '3': 146287461,
      '4': 1,
      '5': 11,
      '6': '.email.ImportDestination',
      '10': 'importdestination'
    },
  ],
};

/// Descriptor for `CreateImportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createImportJobRequestDescriptor = $convert.base64Decode(
    'ChZDcmVhdGVJbXBvcnRKb2JSZXF1ZXN0EkcKEGltcG9ydGRhdGFzb3VyY2UYirrf5wEgASgLMh'
    'cuZW1haWwuSW1wb3J0RGF0YVNvdXJjZVIQaW1wb3J0ZGF0YXNvdXJjZRJJChFpbXBvcnRkZXN0'
    'aW5hdGlvbhjl1uBFIAEoCzIYLmVtYWlsLkltcG9ydERlc3RpbmF0aW9uUhFpbXBvcnRkZXN0aW'
    '5hdGlvbg==');

@$core.Deprecated('Use createImportJobResponseDescriptor instead')
const CreateImportJobResponse$json = {
  '1': 'CreateImportJobResponse',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
  ],
};

/// Descriptor for `CreateImportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createImportJobResponseDescriptor =
    $convert.base64Decode(
        'ChdDcmVhdGVJbXBvcnRKb2JSZXNwb25zZRIXCgVqb2JpZBjS1N0zIAEoCVIFam9iaWQ=');

@$core.Deprecated('Use createMultiRegionEndpointRequestDescriptor instead')
const CreateMultiRegionEndpointRequest$json = {
  '1': 'CreateMultiRegionEndpointRequest',
  '2': [
    {
      '1': 'details',
      '3': 247611974,
      '4': 1,
      '5': 11,
      '6': '.email.Details',
      '10': 'details'
    },
    {'1': 'endpointname', '3': 209534392, '4': 1, '5': 9, '10': 'endpointname'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateMultiRegionEndpointRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMultiRegionEndpointRequestDescriptor =
    $convert.base64Decode(
        'CiBDcmVhdGVNdWx0aVJlZ2lvbkVuZHBvaW50UmVxdWVzdBIrCgdkZXRhaWxzGMaEiXYgASgLMg'
        '4uZW1haWwuRGV0YWlsc1IHZGV0YWlscxIlCgxlbmRwb2ludG5hbWUYuPv0YyABKAlSDGVuZHBv'
        'aW50bmFtZRIiCgR0YWdzGMHB9rUBIAMoCzIKLmVtYWlsLlRhZ1IEdGFncw==');

@$core.Deprecated('Use createMultiRegionEndpointResponseDescriptor instead')
const CreateMultiRegionEndpointResponse$json = {
  '1': 'CreateMultiRegionEndpointResponse',
  '2': [
    {'1': 'endpointid', '3': 35808946, '4': 1, '5': 9, '10': 'endpointid'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.Status',
      '10': 'status'
    },
  ],
};

/// Descriptor for `CreateMultiRegionEndpointResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createMultiRegionEndpointResponseDescriptor =
    $convert.base64Decode(
        'CiFDcmVhdGVNdWx0aVJlZ2lvbkVuZHBvaW50UmVzcG9uc2USIQoKZW5kcG9pbnRpZBiyzYkRIA'
        'EoCVIKZW5kcG9pbnRpZBIoCgZzdGF0dXMYkOT7AiABKA4yDS5lbWFpbC5TdGF0dXNSBnN0YXR1'
        'cw==');

@$core.Deprecated('Use createTenantRequestDescriptor instead')
const CreateTenantRequest$json = {
  '1': 'CreateTenantRequest',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `CreateTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTenantRequestDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVUZW5hbnRSZXF1ZXN0EiIKBHRhZ3MYwcH2tQEgAygLMgouZW1haWwuVGFnUgR0YW'
    'dzEiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbmFudG5hbWU=');

@$core
    .Deprecated('Use createTenantResourceAssociationRequestDescriptor instead')
const CreateTenantResourceAssociationRequest$json = {
  '1': 'CreateTenantResourceAssociationRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `CreateTenantResourceAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTenantResourceAssociationRequestDescriptor =
    $convert.base64Decode(
        'CiZDcmVhdGVUZW5hbnRSZXNvdXJjZUFzc29jaWF0aW9uUmVxdWVzdBIkCgtyZXNvdXJjZWFybh'
        'it+NmtASABKAlSC3Jlc291cmNlYXJuEiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbmFudG5h'
        'bWU=');

@$core
    .Deprecated('Use createTenantResourceAssociationResponseDescriptor instead')
const CreateTenantResourceAssociationResponse$json = {
  '1': 'CreateTenantResourceAssociationResponse',
};

/// Descriptor for `CreateTenantResourceAssociationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTenantResourceAssociationResponseDescriptor =
    $convert.base64Decode(
        'CidDcmVhdGVUZW5hbnRSZXNvdXJjZUFzc29jaWF0aW9uUmVzcG9uc2U=');

@$core.Deprecated('Use createTenantResponseDescriptor instead')
const CreateTenantResponse$json = {
  '1': 'CreateTenantResponse',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'sendingstatus',
      '3': 420634540,
      '4': 1,
      '5': 14,
      '6': '.email.SendingStatus',
      '10': 'sendingstatus'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {'1': 'tenantarn', '3': 19777181, '4': 1, '5': 9, '10': 'tenantarn'},
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `CreateTenantResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createTenantResponseDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVUZW5hbnRSZXNwb25zZRIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCVIQY3'
    'JlYXRlZHRpbWVzdGFtcBI+Cg1zZW5kaW5nc3RhdHVzGKy/ycgBIAEoDjIULmVtYWlsLlNlbmRp'
    'bmdTdGF0dXNSDXNlbmRpbmdzdGF0dXMSIgoEdGFncxjBwfa1ASADKAsyCi5lbWFpbC5UYWdSBH'
    'RhZ3MSHwoJdGVuYW50YXJuGJ2NtwkgASgJUgl0ZW5hbnRhcm4SHQoIdGVuYW50aWQYgfjvZCAB'
    'KAlSCHRlbmFudGlkEiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbmFudG5hbWU=');

@$core
    .Deprecated('Use customVerificationEmailTemplateMetadataDescriptor instead')
const CustomVerificationEmailTemplateMetadata$json = {
  '1': 'CustomVerificationEmailTemplateMetadata',
  '2': [
    {
      '1': 'failureredirectionurl',
      '3': 376073007,
      '4': 1,
      '5': 9,
      '10': 'failureredirectionurl'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {
      '1': 'successredirectionurl',
      '3': 513750768,
      '4': 1,
      '5': 9,
      '10': 'successredirectionurl'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
    {
      '1': 'templatesubject',
      '3': 144653738,
      '4': 1,
      '5': 9,
      '10': 'templatesubject'
    },
  ],
};

/// Descriptor for `CustomVerificationEmailTemplateMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customVerificationEmailTemplateMetadataDescriptor =
    $convert.base64Decode(
        'CidDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlTWV0YWRhdGESOAoVZmFpbHVyZXJlZG'
        'lyZWN0aW9udXJsGK/WqbMBIAEoCVIVZmFpbHVyZXJlZGlyZWN0aW9udXJsEi0KEGZyb21lbWFp'
        'bGFkZHJlc3MYhprLLCABKAlSEGZyb21lbWFpbGFkZHJlc3MSOAoVc3VjY2Vzc3JlZGlyZWN0aW'
        '9udXJsGPDt/PQBIAEoCVIVc3VjY2Vzc3JlZGlyZWN0aW9udXJsEiYKDHRlbXBsYXRlbmFtZRix'
        'mJHlASABKAlSDHRlbXBsYXRlbmFtZRIrCg90ZW1wbGF0ZXN1YmplY3QYqvv8RCABKAlSD3RlbX'
        'BsYXRlc3ViamVjdA==');

@$core.Deprecated('Use dailyVolumeDescriptor instead')
const DailyVolume$json = {
  '1': 'DailyVolume',
  '2': [
    {
      '1': 'domainispplacements',
      '3': 63848508,
      '4': 3,
      '5': 11,
      '6': '.email.DomainIspPlacement',
      '10': 'domainispplacements'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'volumestatistics',
      '3': 331663309,
      '4': 1,
      '5': 11,
      '6': '.email.VolumeStatistics',
      '10': 'volumestatistics'
    },
  ],
};

/// Descriptor for `DailyVolume`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dailyVolumeDescriptor = $convert.base64Decode(
    'CgtEYWlseVZvbHVtZRJOChNkb21haW5pc3BwbGFjZW1lbnRzGLyAuR4gAygLMhkuZW1haWwuRG'
    '9tYWluSXNwUGxhY2VtZW50UhNkb21haW5pc3BwbGFjZW1lbnRzEiAKCXN0YXJ0ZGF0ZRj8+KDU'
    'ASABKAlSCXN0YXJ0ZGF0ZRJHChB2b2x1bWVzdGF0aXN0aWNzGM2Pk54BIAEoCzIXLmVtYWlsLl'
    'ZvbHVtZVN0YXRpc3RpY3NSEHZvbHVtZXN0YXRpc3RpY3M=');

@$core.Deprecated('Use dashboardAttributesDescriptor instead')
const DashboardAttributes$json = {
  '1': 'DashboardAttributes',
  '2': [
    {
      '1': 'engagementmetrics',
      '3': 102186588,
      '4': 1,
      '5': 14,
      '6': '.email.FeatureStatus',
      '10': 'engagementmetrics'
    },
  ],
};

/// Descriptor for `DashboardAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dashboardAttributesDescriptor = $convert.base64Decode(
    'ChNEYXNoYm9hcmRBdHRyaWJ1dGVzEkUKEWVuZ2FnZW1lbnRtZXRyaWNzGNz83DAgASgOMhQuZW'
    '1haWwuRmVhdHVyZVN0YXR1c1IRZW5nYWdlbWVudG1ldHJpY3M=');

@$core.Deprecated('Use dashboardOptionsDescriptor instead')
const DashboardOptions$json = {
  '1': 'DashboardOptions',
  '2': [
    {
      '1': 'engagementmetrics',
      '3': 102186588,
      '4': 1,
      '5': 14,
      '6': '.email.FeatureStatus',
      '10': 'engagementmetrics'
    },
  ],
};

/// Descriptor for `DashboardOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dashboardOptionsDescriptor = $convert.base64Decode(
    'ChBEYXNoYm9hcmRPcHRpb25zEkUKEWVuZ2FnZW1lbnRtZXRyaWNzGNz83DAgASgOMhQuZW1haW'
    'wuRmVhdHVyZVN0YXR1c1IRZW5nYWdlbWVudG1ldHJpY3M=');

@$core.Deprecated('Use dedicatedIpDescriptor instead')
const DedicatedIp$json = {
  '1': 'DedicatedIp',
  '2': [
    {'1': 'ip', '3': 183031933, '4': 1, '5': 9, '10': 'ip'},
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
    {
      '1': 'warmuppercentage',
      '3': 389585846,
      '4': 1,
      '5': 5,
      '10': 'warmuppercentage'
    },
    {
      '1': 'warmupstatus',
      '3': 117093872,
      '4': 1,
      '5': 14,
      '6': '.email.WarmupStatus',
      '10': 'warmupstatus'
    },
  ],
};

/// Descriptor for `DedicatedIp`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dedicatedIpDescriptor = $convert.base64Decode(
    'CgtEZWRpY2F0ZWRJcBIRCgJpcBj9sKNXIAEoCVICaXASHQoIcG9vbG5hbWUYyY2FJyABKAlSCH'
    'Bvb2xuYW1lEi4KEHdhcm11cHBlcmNlbnRhZ2UYtrfiuQEgASgFUhB3YXJtdXBwZXJjZW50YWdl'
    'EjoKDHdhcm11cHN0YXR1cxjw6+o3IAEoDjITLmVtYWlsLldhcm11cFN0YXR1c1IMd2FybXVwc3'
    'RhdHVz');

@$core.Deprecated('Use dedicatedIpPoolDescriptor instead')
const DedicatedIpPool$json = {
  '1': 'DedicatedIpPool',
  '2': [
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
    {
      '1': 'scalingmode',
      '3': 210356138,
      '4': 1,
      '5': 14,
      '6': '.email.ScalingMode',
      '10': 'scalingmode'
    },
  ],
};

/// Descriptor for `DedicatedIpPool`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dedicatedIpPoolDescriptor = $convert.base64Decode(
    'Cg9EZWRpY2F0ZWRJcFBvb2wSHQoIcG9vbG5hbWUYyY2FJyABKAlSCHBvb2xuYW1lEjcKC3NjYW'
    'xpbmdtb2RlGKqPp2QgASgOMhIuZW1haWwuU2NhbGluZ01vZGVSC3NjYWxpbmdtb2Rl');

@$core.Deprecated(
    'Use deleteConfigurationSetEventDestinationRequestDescriptor instead')
const DeleteConfigurationSetEventDestinationRequest$json = {
  '1': 'DeleteConfigurationSetEventDestinationRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'eventdestinationname',
      '3': 477517655,
      '4': 1,
      '5': 9,
      '10': 'eventdestinationname'
    },
  ],
};

/// Descriptor for `DeleteConfigurationSetEventDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteConfigurationSetEventDestinationRequestDescriptor =
    $convert.base64Decode(
        'Ci1EZWxldGVDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvblJlcXVlc3QSNgoUY29uZm'
        'lndXJhdGlvbnNldG5hbWUYzYuxwAEgASgJUhRjb25maWd1cmF0aW9uc2V0bmFtZRI2ChRldmVu'
        'dGRlc3RpbmF0aW9ubmFtZRjXrtnjASABKAlSFGV2ZW50ZGVzdGluYXRpb25uYW1l');

@$core.Deprecated(
    'Use deleteConfigurationSetEventDestinationResponseDescriptor instead')
const DeleteConfigurationSetEventDestinationResponse$json = {
  '1': 'DeleteConfigurationSetEventDestinationResponse',
};

/// Descriptor for `DeleteConfigurationSetEventDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteConfigurationSetEventDestinationResponseDescriptor =
    $convert.base64Decode(
        'Ci5EZWxldGVDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvblJlc3BvbnNl');

@$core.Deprecated('Use deleteConfigurationSetRequestDescriptor instead')
const DeleteConfigurationSetRequest$json = {
  '1': 'DeleteConfigurationSetRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
  ],
};

/// Descriptor for `DeleteConfigurationSetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteConfigurationSetRequestDescriptor =
    $convert.base64Decode(
        'Ch1EZWxldGVDb25maWd1cmF0aW9uU2V0UmVxdWVzdBI2ChRjb25maWd1cmF0aW9uc2V0bmFtZR'
        'jNi7HAASABKAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1l');

@$core.Deprecated('Use deleteConfigurationSetResponseDescriptor instead')
const DeleteConfigurationSetResponse$json = {
  '1': 'DeleteConfigurationSetResponse',
};

/// Descriptor for `DeleteConfigurationSetResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteConfigurationSetResponseDescriptor =
    $convert.base64Decode('Ch5EZWxldGVDb25maWd1cmF0aW9uU2V0UmVzcG9uc2U=');

@$core.Deprecated('Use deleteContactListRequestDescriptor instead')
const DeleteContactListRequest$json = {
  '1': 'DeleteContactListRequest',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
  ],
};

/// Descriptor for `DeleteContactListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteContactListRequestDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVDb250YWN0TGlzdFJlcXVlc3QSLAoPY29udGFjdGxpc3RuYW1lGPG9p6EBIAEoCV'
        'IPY29udGFjdGxpc3RuYW1l');

@$core.Deprecated('Use deleteContactListResponseDescriptor instead')
const DeleteContactListResponse$json = {
  '1': 'DeleteContactListResponse',
};

/// Descriptor for `DeleteContactListResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteContactListResponseDescriptor =
    $convert.base64Decode('ChlEZWxldGVDb250YWN0TGlzdFJlc3BvbnNl');

@$core.Deprecated('Use deleteContactRequestDescriptor instead')
const DeleteContactRequest$json = {
  '1': 'DeleteContactRequest',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
  ],
};

/// Descriptor for `DeleteContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteContactRequestDescriptor = $convert.base64Decode(
    'ChREZWxldGVDb250YWN0UmVxdWVzdBIsCg9jb250YWN0bGlzdG5hbWUY8b2noQEgASgJUg9jb2'
    '50YWN0bGlzdG5hbWUSJgoMZW1haWxhZGRyZXNzGNbxiukBIAEoCVIMZW1haWxhZGRyZXNz');

@$core.Deprecated('Use deleteContactResponseDescriptor instead')
const DeleteContactResponse$json = {
  '1': 'DeleteContactResponse',
};

/// Descriptor for `DeleteContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteContactResponseDescriptor =
    $convert.base64Decode('ChVEZWxldGVDb250YWN0UmVzcG9uc2U=');

@$core.Deprecated(
    'Use deleteCustomVerificationEmailTemplateRequestDescriptor instead')
const DeleteCustomVerificationEmailTemplateRequest$json = {
  '1': 'DeleteCustomVerificationEmailTemplateRequest',
  '2': [
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `DeleteCustomVerificationEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteCustomVerificationEmailTemplateRequestDescriptor =
    $convert.base64Decode(
        'CixEZWxldGVDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVxdWVzdBImCgx0ZW1wbG'
        'F0ZW5hbWUYsZiR5QEgASgJUgx0ZW1wbGF0ZW5hbWU=');

@$core.Deprecated(
    'Use deleteCustomVerificationEmailTemplateResponseDescriptor instead')
const DeleteCustomVerificationEmailTemplateResponse$json = {
  '1': 'DeleteCustomVerificationEmailTemplateResponse',
};

/// Descriptor for `DeleteCustomVerificationEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    deleteCustomVerificationEmailTemplateResponseDescriptor =
    $convert.base64Decode(
        'Ci1EZWxldGVDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVzcG9uc2U=');

@$core.Deprecated('Use deleteDedicatedIpPoolRequestDescriptor instead')
const DeleteDedicatedIpPoolRequest$json = {
  '1': 'DeleteDedicatedIpPoolRequest',
  '2': [
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
  ],
};

/// Descriptor for `DeleteDedicatedIpPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDedicatedIpPoolRequestDescriptor =
    $convert.base64Decode(
        'ChxEZWxldGVEZWRpY2F0ZWRJcFBvb2xSZXF1ZXN0Eh0KCHBvb2xuYW1lGMmNhScgASgJUghwb2'
        '9sbmFtZQ==');

@$core.Deprecated('Use deleteDedicatedIpPoolResponseDescriptor instead')
const DeleteDedicatedIpPoolResponse$json = {
  '1': 'DeleteDedicatedIpPoolResponse',
};

/// Descriptor for `DeleteDedicatedIpPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteDedicatedIpPoolResponseDescriptor =
    $convert.base64Decode('Ch1EZWxldGVEZWRpY2F0ZWRJcFBvb2xSZXNwb25zZQ==');

@$core.Deprecated('Use deleteEmailIdentityPolicyRequestDescriptor instead')
const DeleteEmailIdentityPolicyRequest$json = {
  '1': 'DeleteEmailIdentityPolicyRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `DeleteEmailIdentityPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEmailIdentityPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVFbWFpbElkZW50aXR5UG9saWN5UmVxdWVzdBInCg1lbWFpbGlkZW50aXR5GJqU8k'
        'AgASgJUg1lbWFpbGlkZW50aXR5EiEKCnBvbGljeW5hbWUYvfWHfyABKAlSCnBvbGljeW5hbWU=');

@$core.Deprecated('Use deleteEmailIdentityPolicyResponseDescriptor instead')
const DeleteEmailIdentityPolicyResponse$json = {
  '1': 'DeleteEmailIdentityPolicyResponse',
};

/// Descriptor for `DeleteEmailIdentityPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEmailIdentityPolicyResponseDescriptor =
    $convert.base64Decode('CiFEZWxldGVFbWFpbElkZW50aXR5UG9saWN5UmVzcG9uc2U=');

@$core.Deprecated('Use deleteEmailIdentityRequestDescriptor instead')
const DeleteEmailIdentityRequest$json = {
  '1': 'DeleteEmailIdentityRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
  ],
};

/// Descriptor for `DeleteEmailIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEmailIdentityRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVFbWFpbElkZW50aXR5UmVxdWVzdBInCg1lbWFpbGlkZW50aXR5GJqU8kAgASgJUg'
        '1lbWFpbGlkZW50aXR5');

@$core.Deprecated('Use deleteEmailIdentityResponseDescriptor instead')
const DeleteEmailIdentityResponse$json = {
  '1': 'DeleteEmailIdentityResponse',
};

/// Descriptor for `DeleteEmailIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEmailIdentityResponseDescriptor =
    $convert.base64Decode('ChtEZWxldGVFbWFpbElkZW50aXR5UmVzcG9uc2U=');

@$core.Deprecated('Use deleteEmailTemplateRequestDescriptor instead')
const DeleteEmailTemplateRequest$json = {
  '1': 'DeleteEmailTemplateRequest',
  '2': [
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `DeleteEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEmailTemplateRequestDescriptor =
    $convert.base64Decode(
        'ChpEZWxldGVFbWFpbFRlbXBsYXRlUmVxdWVzdBImCgx0ZW1wbGF0ZW5hbWUYsZiR5QEgASgJUg'
        'x0ZW1wbGF0ZW5hbWU=');

@$core.Deprecated('Use deleteEmailTemplateResponseDescriptor instead')
const DeleteEmailTemplateResponse$json = {
  '1': 'DeleteEmailTemplateResponse',
};

/// Descriptor for `DeleteEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteEmailTemplateResponseDescriptor =
    $convert.base64Decode('ChtEZWxldGVFbWFpbFRlbXBsYXRlUmVzcG9uc2U=');

@$core.Deprecated('Use deleteMultiRegionEndpointRequestDescriptor instead')
const DeleteMultiRegionEndpointRequest$json = {
  '1': 'DeleteMultiRegionEndpointRequest',
  '2': [
    {'1': 'endpointname', '3': 209534392, '4': 1, '5': 9, '10': 'endpointname'},
  ],
};

/// Descriptor for `DeleteMultiRegionEndpointRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMultiRegionEndpointRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVNdWx0aVJlZ2lvbkVuZHBvaW50UmVxdWVzdBIlCgxlbmRwb2ludG5hbWUYuPv0Yy'
        'ABKAlSDGVuZHBvaW50bmFtZQ==');

@$core.Deprecated('Use deleteMultiRegionEndpointResponseDescriptor instead')
const DeleteMultiRegionEndpointResponse$json = {
  '1': 'DeleteMultiRegionEndpointResponse',
  '2': [
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.Status',
      '10': 'status'
    },
  ],
};

/// Descriptor for `DeleteMultiRegionEndpointResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteMultiRegionEndpointResponseDescriptor =
    $convert.base64Decode(
        'CiFEZWxldGVNdWx0aVJlZ2lvbkVuZHBvaW50UmVzcG9uc2USKAoGc3RhdHVzGJDk+wIgASgOMg'
        '0uZW1haWwuU3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use deleteSuppressedDestinationRequestDescriptor instead')
const DeleteSuppressedDestinationRequest$json = {
  '1': 'DeleteSuppressedDestinationRequest',
  '2': [
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
  ],
};

/// Descriptor for `DeleteSuppressedDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSuppressedDestinationRequestDescriptor =
    $convert.base64Decode(
        'CiJEZWxldGVTdXBwcmVzc2VkRGVzdGluYXRpb25SZXF1ZXN0EiYKDGVtYWlsYWRkcmVzcxjW8Y'
        'rpASABKAlSDGVtYWlsYWRkcmVzcw==');

@$core.Deprecated('Use deleteSuppressedDestinationResponseDescriptor instead')
const DeleteSuppressedDestinationResponse$json = {
  '1': 'DeleteSuppressedDestinationResponse',
};

/// Descriptor for `DeleteSuppressedDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteSuppressedDestinationResponseDescriptor =
    $convert
        .base64Decode('CiNEZWxldGVTdXBwcmVzc2VkRGVzdGluYXRpb25SZXNwb25zZQ==');

@$core.Deprecated('Use deleteTenantRequestDescriptor instead')
const DeleteTenantRequest$json = {
  '1': 'DeleteTenantRequest',
  '2': [
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `DeleteTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTenantRequestDescriptor = $convert.base64Decode(
    'ChNEZWxldGVUZW5hbnRSZXF1ZXN0EiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbmFudG5hbW'
    'U=');

@$core
    .Deprecated('Use deleteTenantResourceAssociationRequestDescriptor instead')
const DeleteTenantResourceAssociationRequest$json = {
  '1': 'DeleteTenantResourceAssociationRequest',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `DeleteTenantResourceAssociationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTenantResourceAssociationRequestDescriptor =
    $convert.base64Decode(
        'CiZEZWxldGVUZW5hbnRSZXNvdXJjZUFzc29jaWF0aW9uUmVxdWVzdBIkCgtyZXNvdXJjZWFybh'
        'it+NmtASABKAlSC3Jlc291cmNlYXJuEiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbmFudG5h'
        'bWU=');

@$core
    .Deprecated('Use deleteTenantResourceAssociationResponseDescriptor instead')
const DeleteTenantResourceAssociationResponse$json = {
  '1': 'DeleteTenantResourceAssociationResponse',
};

/// Descriptor for `DeleteTenantResourceAssociationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTenantResourceAssociationResponseDescriptor =
    $convert.base64Decode(
        'CidEZWxldGVUZW5hbnRSZXNvdXJjZUFzc29jaWF0aW9uUmVzcG9uc2U=');

@$core.Deprecated('Use deleteTenantResponseDescriptor instead')
const DeleteTenantResponse$json = {
  '1': 'DeleteTenantResponse',
};

/// Descriptor for `DeleteTenantResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteTenantResponseDescriptor =
    $convert.base64Decode('ChREZWxldGVUZW5hbnRSZXNwb25zZQ==');

@$core.Deprecated('Use deliverabilityTestReportDescriptor instead')
const DeliverabilityTestReport$json = {
  '1': 'DeliverabilityTestReport',
  '2': [
    {'1': 'createdate', '3': 37690514, '4': 1, '5': 9, '10': 'createdate'},
    {
      '1': 'deliverabilityteststatus',
      '3': 71311387,
      '4': 1,
      '5': 14,
      '6': '.email.DeliverabilityTestStatus',
      '10': 'deliverabilityteststatus'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {'1': 'reportid', '3': 420903847, '4': 1, '5': 9, '10': 'reportid'},
    {'1': 'reportname', '3': 526054737, '4': 1, '5': 9, '10': 'reportname'},
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
  ],
};

/// Descriptor for `DeliverabilityTestReport`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deliverabilityTestReportDescriptor = $convert.base64Decode(
    'ChhEZWxpdmVyYWJpbGl0eVRlc3RSZXBvcnQSIQoKY3JlYXRlZGF0ZRiSufwRIAEoCVIKY3JlYX'
    'RlZGF0ZRJeChhkZWxpdmVyYWJpbGl0eXRlc3RzdGF0dXMYm8CAIiABKA4yHy5lbWFpbC5EZWxp'
    'dmVyYWJpbGl0eVRlc3RTdGF0dXNSGGRlbGl2ZXJhYmlsaXR5dGVzdHN0YXR1cxItChBmcm9tZW'
    '1haWxhZGRyZXNzGIaayywgASgJUhBmcm9tZW1haWxhZGRyZXNzEh4KCHJlcG9ydGlkGKf32cgB'
    'IAEoCVIIcmVwb3J0aWQSIgoKcmVwb3J0bmFtZRjR6uv6ASABKAlSCnJlcG9ydG5hbWUSGwoHc3'
    'ViamVjdBjwyeQDIAEoCVIHc3ViamVjdA==');

@$core.Deprecated('Use deliveryOptionsDescriptor instead')
const DeliveryOptions$json = {
  '1': 'DeliveryOptions',
  '2': [
    {
      '1': 'maxdeliveryseconds',
      '3': 16229983,
      '4': 1,
      '5': 3,
      '10': 'maxdeliveryseconds'
    },
    {
      '1': 'sendingpoolname',
      '3': 89398333,
      '4': 1,
      '5': 9,
      '10': 'sendingpoolname'
    },
    {
      '1': 'tlspolicy',
      '3': 127629,
      '4': 1,
      '5': 14,
      '6': '.email.TlsPolicy',
      '10': 'tlspolicy'
    },
  ],
};

/// Descriptor for `DeliveryOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deliveryOptionsDescriptor = $convert.base64Decode(
    'Cg9EZWxpdmVyeU9wdGlvbnMSMQoSbWF4ZGVsaXZlcnlzZWNvbmRzGN/M3gcgASgDUhJtYXhkZW'
    'xpdmVyeXNlY29uZHMSKwoPc2VuZGluZ3Bvb2xuYW1lGL240CogASgJUg9zZW5kaW5ncG9vbG5h'
    'bWUSMAoJdGxzcG9saWN5GI3lByABKA4yEC5lbWFpbC5UbHNQb2xpY3lSCXRsc3BvbGljeQ==');

@$core.Deprecated('Use destinationDescriptor instead')
const Destination$json = {
  '1': 'Destination',
  '2': [
    {'1': 'bccaddresses', '3': 98755614, '4': 3, '5': 9, '10': 'bccaddresses'},
    {'1': 'ccaddresses', '3': 196613182, '4': 3, '5': 9, '10': 'ccaddresses'},
    {'1': 'toaddresses', '3': 227850817, '4': 3, '5': 9, '10': 'toaddresses'},
  ],
};

/// Descriptor for `Destination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List destinationDescriptor = $convert.base64Decode(
    'CgtEZXN0aW5hdGlvbhIlCgxiY2NhZGRyZXNzZXMYnsiLLyADKAlSDGJjY2FkZHJlc3NlcxIjCg'
    'tjY2FkZHJlc3Nlcxi+qOBdIAMoCVILY2NhZGRyZXNzZXMSIwoLdG9hZGRyZXNzZXMYwfTSbCAD'
    'KAlSC3RvYWRkcmVzc2Vz');

@$core.Deprecated('Use detailsDescriptor instead')
const Details$json = {
  '1': 'Details',
  '2': [
    {
      '1': 'routesdetails',
      '3': 318115776,
      '4': 3,
      '5': 11,
      '6': '.email.RouteDetails',
      '10': 'routesdetails'
    },
  ],
};

/// Descriptor for `Details`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List detailsDescriptor = $convert.base64Decode(
    'CgdEZXRhaWxzEj0KDXJvdXRlc2RldGFpbHMYwJ/YlwEgAygLMhMuZW1haWwuUm91dGVEZXRhaW'
    'xzUg1yb3V0ZXNkZXRhaWxz');

@$core.Deprecated('Use dkimAttributesDescriptor instead')
const DkimAttributes$json = {
  '1': 'DkimAttributes',
  '2': [
    {
      '1': 'currentsigningkeylength',
      '3': 260519333,
      '4': 1,
      '5': 14,
      '6': '.email.DkimSigningKeyLength',
      '10': 'currentsigningkeylength'
    },
    {
      '1': 'lastkeygenerationtimestamp',
      '3': 455425631,
      '4': 1,
      '5': 9,
      '10': 'lastkeygenerationtimestamp'
    },
    {
      '1': 'nextsigningkeylength',
      '3': 307141083,
      '4': 1,
      '5': 14,
      '6': '.email.DkimSigningKeyLength',
      '10': 'nextsigningkeylength'
    },
    {
      '1': 'signingattributesorigin',
      '3': 255219696,
      '4': 1,
      '5': 14,
      '6': '.email.DkimSigningAttributesOrigin',
      '10': 'signingattributesorigin'
    },
    {
      '1': 'signingenabled',
      '3': 289566486,
      '4': 1,
      '5': 8,
      '10': 'signingenabled'
    },
    {
      '1': 'signinghostedzone',
      '3': 442955686,
      '4': 1,
      '5': 9,
      '10': 'signinghostedzone'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.DkimStatus',
      '10': 'status'
    },
    {'1': 'tokens', '3': 50282100, '4': 3, '5': 9, '10': 'tokens'},
  ],
};

/// Descriptor for `DkimAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dkimAttributesDescriptor = $convert.base64Decode(
    'Cg5Ea2ltQXR0cmlidXRlcxJYChdjdXJyZW50c2lnbmluZ2tleWxlbmd0aBil65x8IAEoDjIbLm'
    'VtYWlsLkRraW1TaWduaW5nS2V5TGVuZ3RoUhdjdXJyZW50c2lnbmluZ2tleWxlbmd0aBJCChps'
    'YXN0a2V5Z2VuZXJhdGlvbnRpbWVzdGFtcBjf/JTZASABKAlSGmxhc3RrZXlnZW5lcmF0aW9udG'
    'ltZXN0YW1wElMKFG5leHRzaWduaW5na2V5bGVuZ3RoGNuzupIBIAEoDjIbLmVtYWlsLkRraW1T'
    'aWduaW5nS2V5TGVuZ3RoUhRuZXh0c2lnbmluZ2tleWxlbmd0aBJfChdzaWduaW5nYXR0cmlidX'
    'Rlc29yaWdpbhjwr9l5IAEoDjIiLmVtYWlsLkRraW1TaWduaW5nQXR0cmlidXRlc09yaWdpblIX'
    'c2lnbmluZ2F0dHJpYnV0ZXNvcmlnaW4SKgoOc2lnbmluZ2VuYWJsZWQYlt6JigEgASgIUg5zaW'
    'duaW5nZW5hYmxlZBIwChFzaWduaW5naG9zdGVkem9uZRim75vTASABKAlSEXNpZ25pbmdob3N0'
    'ZWR6b25lEiwKBnN0YXR1cxiQ5PsCIAEoDjIRLmVtYWlsLkRraW1TdGF0dXNSBnN0YXR1cxIZCg'
    'Z0b2tlbnMY9Pz8FyADKAlSBnRva2Vucw==');

@$core.Deprecated('Use dkimSigningAttributesDescriptor instead')
const DkimSigningAttributes$json = {
  '1': 'DkimSigningAttributes',
  '2': [
    {
      '1': 'domainsigningattributesorigin',
      '3': 248837556,
      '4': 1,
      '5': 14,
      '6': '.email.DkimSigningAttributesOrigin',
      '10': 'domainsigningattributesorigin'
    },
    {
      '1': 'domainsigningprivatekey',
      '3': 502133867,
      '4': 1,
      '5': 9,
      '10': 'domainsigningprivatekey'
    },
    {
      '1': 'domainsigningselector',
      '3': 454592926,
      '4': 1,
      '5': 9,
      '10': 'domainsigningselector'
    },
    {
      '1': 'nextsigningkeylength',
      '3': 307141083,
      '4': 1,
      '5': 14,
      '6': '.email.DkimSigningKeyLength',
      '10': 'nextsigningkeylength'
    },
  ],
};

/// Descriptor for `DkimSigningAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dkimSigningAttributesDescriptor = $convert.base64Decode(
    'ChVEa2ltU2lnbmluZ0F0dHJpYnV0ZXMSawodZG9tYWluc2lnbmluZ2F0dHJpYnV0ZXNvcmlnaW'
    '4YtOvTdiABKA4yIi5lbWFpbC5Ea2ltU2lnbmluZ0F0dHJpYnV0ZXNPcmlnaW5SHWRvbWFpbnNp'
    'Z25pbmdhdHRyaWJ1dGVzb3JpZ2luEjwKF2RvbWFpbnNpZ25pbmdwcml2YXRla2V5GOvot+8BIA'
    'EoCVIXZG9tYWluc2lnbmluZ3ByaXZhdGVrZXkSOAoVZG9tYWluc2lnbmluZ3NlbGVjdG9yGJ6T'
    '4tgBIAEoCVIVZG9tYWluc2lnbmluZ3NlbGVjdG9yElMKFG5leHRzaWduaW5na2V5bGVuZ3RoGN'
    'uzupIBIAEoDjIbLmVtYWlsLkRraW1TaWduaW5nS2V5TGVuZ3RoUhRuZXh0c2lnbmluZ2tleWxl'
    'bmd0aA==');

@$core.Deprecated('Use domainDeliverabilityCampaignDescriptor instead')
const DomainDeliverabilityCampaign$json = {
  '1': 'DomainDeliverabilityCampaign',
  '2': [
    {'1': 'campaignid', '3': 299715775, '4': 1, '5': 9, '10': 'campaignid'},
    {'1': 'deleterate', '3': 56287481, '4': 1, '5': 1, '10': 'deleterate'},
    {'1': 'esps', '3': 526407961, '4': 3, '5': 9, '10': 'esps'},
    {
      '1': 'firstseendatetime',
      '3': 77536972,
      '4': 1,
      '5': 9,
      '10': 'firstseendatetime'
    },
    {'1': 'fromaddress', '3': 84397700, '4': 1, '5': 9, '10': 'fromaddress'},
    {'1': 'imageurl', '3': 361905604, '4': 1, '5': 9, '10': 'imageurl'},
    {'1': 'inboxcount', '3': 96442763, '4': 1, '5': 3, '10': 'inboxcount'},
    {
      '1': 'lastseendatetime',
      '3': 2010772,
      '4': 1,
      '5': 9,
      '10': 'lastseendatetime'
    },
    {
      '1': 'projectedvolume',
      '3': 71581402,
      '4': 1,
      '5': 3,
      '10': 'projectedvolume'
    },
    {
      '1': 'readdeleterate',
      '3': 249983513,
      '4': 1,
      '5': 1,
      '10': 'readdeleterate'
    },
    {'1': 'readrate', '3': 412487272, '4': 1, '5': 1, '10': 'readrate'},
    {'1': 'sendingips', '3': 344517626, '4': 3, '5': 9, '10': 'sendingips'},
    {'1': 'spamcount', '3': 224064008, '4': 1, '5': 3, '10': 'spamcount'},
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
  ],
};

/// Descriptor for `DomainDeliverabilityCampaign`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainDeliverabilityCampaignDescriptor = $convert.base64Decode(
    'ChxEb21haW5EZWxpdmVyYWJpbGl0eUNhbXBhaWduEiIKCmNhbXBhaWduaWQYv5n1jgEgASgJUg'
    'pjYW1wYWlnbmlkEiEKCmRlbGV0ZXJhdGUY+cHrGiABKAFSCmRlbGV0ZXJhdGUSFgoEZXNwcxiZ'
    'soH7ASADKAlSBGVzcHMSLwoRZmlyc3RzZWVuZGF0ZXRpbWUYzL38JCABKAlSEWZpcnN0c2Vlbm'
    'RhdGV0aW1lEiMKC2Zyb21hZGRyZXNzGISdnyggASgJUgtmcm9tYWRkcmVzcxIeCghpbWFnZXVy'
    'bBjE+8isASABKAlSCGltYWdldXJsEiEKCmluYm94Y291bnQYi7P+LSABKANSCmluYm94Y291bn'
    'QSLAoQbGFzdHNlZW5kYXRldGltZRiU3XogASgJUhBsYXN0c2VlbmRhdGV0aW1lEisKD3Byb2pl'
    'Y3RlZHZvbHVtZRja/ZAiIAEoA1IPcHJvamVjdGVkdm9sdW1lEikKDnJlYWRkZWxldGVyYXRlGJ'
    'nkmXcgASgBUg5yZWFkZGVsZXRlcmF0ZRIeCghyZWFkcmF0ZRjonNjEASABKAFSCHJlYWRyYXRl'
    'EiIKCnNlbmRpbmdpcHMY+tejpAEgAygJUgpzZW5kaW5naXBzEh8KCXNwYW1jb3VudBiI5OtqIA'
    'EoA1IJc3BhbWNvdW50EhsKB3N1YmplY3QY8MnkAyABKAlSB3N1YmplY3Q=');

@$core.Deprecated('Use domainDeliverabilityTrackingOptionDescriptor instead')
const DomainDeliverabilityTrackingOption$json = {
  '1': 'DomainDeliverabilityTrackingOption',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {
      '1': 'inboxplacementtrackingoption',
      '3': 155827259,
      '4': 1,
      '5': 11,
      '6': '.email.InboxPlacementTrackingOption',
      '10': 'inboxplacementtrackingoption'
    },
    {
      '1': 'subscriptionstartdate',
      '3': 514104205,
      '4': 1,
      '5': 9,
      '10': 'subscriptionstartdate'
    },
  ],
};

/// Descriptor for `DomainDeliverabilityTrackingOption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainDeliverabilityTrackingOptionDescriptor =
    $convert.base64Decode(
        'CiJEb21haW5EZWxpdmVyYWJpbGl0eVRyYWNraW5nT3B0aW9uEhoKBmRvbWFpbhiSkvLwASABKA'
        'lSBmRvbWFpbhJqChxpbmJveHBsYWNlbWVudHRyYWNraW5nb3B0aW9uGLv4pkogASgLMiMuZW1h'
        'aWwuSW5ib3hQbGFjZW1lbnRUcmFja2luZ09wdGlvblIcaW5ib3hwbGFjZW1lbnR0cmFja2luZ2'
        '9wdGlvbhI4ChVzdWJzY3JpcHRpb25zdGFydGRhdGUYjbeS9QEgASgJUhVzdWJzY3JpcHRpb25z'
        'dGFydGRhdGU=');

@$core.Deprecated('Use domainIspPlacementDescriptor instead')
const DomainIspPlacement$json = {
  '1': 'DomainIspPlacement',
  '2': [
    {
      '1': 'inboxpercentage',
      '3': 81442932,
      '4': 1,
      '5': 1,
      '10': 'inboxpercentage'
    },
    {
      '1': 'inboxrawcount',
      '3': 215061115,
      '4': 1,
      '5': 3,
      '10': 'inboxrawcount'
    },
    {'1': 'ispname', '3': 79681213, '4': 1, '5': 9, '10': 'ispname'},
    {
      '1': 'spampercentage',
      '3': 364813289,
      '4': 1,
      '5': 1,
      '10': 'spampercentage'
    },
    {'1': 'spamrawcount', '3': 440722366, '4': 1, '5': 3, '10': 'spamrawcount'},
  ],
};

/// Descriptor for `DomainIspPlacement`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List domainIspPlacementDescriptor = $convert.base64Decode(
    'ChJEb21haW5Jc3BQbGFjZW1lbnQSKwoPaW5ib3hwZXJjZW50YWdlGPTw6iYgASgBUg9pbmJveH'
    'BlcmNlbnRhZ2USJwoNaW5ib3hyYXdjb3VudBj7pMZmIAEoA1INaW5ib3hyYXdjb3VudBIbCgdp'
    'c3BuYW1lGL2t/yUgASgJUgdpc3BuYW1lEioKDnNwYW1wZXJjZW50YWdlGOm3+q0BIAEoAVIOc3'
    'BhbXBlcmNlbnRhZ2USJgoMc3BhbXJhd2NvdW50GL7Hk9IBIAEoA1IMc3BhbXJhd2NvdW50');

@$core
    .Deprecated('Use emailAddressInsightsMailboxEvaluationsDescriptor instead')
const EmailAddressInsightsMailboxEvaluations$json = {
  '1': 'EmailAddressInsightsMailboxEvaluations',
  '2': [
    {
      '1': 'hasvaliddnsrecords',
      '3': 379952101,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsVerdict',
      '10': 'hasvaliddnsrecords'
    },
    {
      '1': 'hasvalidsyntax',
      '3': 368261973,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsVerdict',
      '10': 'hasvalidsyntax'
    },
    {
      '1': 'isdisposable',
      '3': 181632148,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsVerdict',
      '10': 'isdisposable'
    },
    {
      '1': 'israndominput',
      '3': 254287319,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsVerdict',
      '10': 'israndominput'
    },
    {
      '1': 'isroleaddress',
      '3': 61297080,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsVerdict',
      '10': 'isroleaddress'
    },
    {
      '1': 'mailboxexists',
      '3': 82590460,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsVerdict',
      '10': 'mailboxexists'
    },
  ],
};

/// Descriptor for `EmailAddressInsightsMailboxEvaluations`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailAddressInsightsMailboxEvaluationsDescriptor = $convert.base64Decode(
    'CiZFbWFpbEFkZHJlc3NJbnNpZ2h0c01haWxib3hFdmFsdWF0aW9ucxJWChJoYXN2YWxpZGRuc3'
    'JlY29yZHMY5beWtQEgASgLMiIuZW1haWwuRW1haWxBZGRyZXNzSW5zaWdodHNWZXJkaWN0UhJo'
    'YXN2YWxpZGRuc3JlY29yZHMSTgoOaGFzdmFsaWRzeW50YXgY1fbMrwEgASgLMiIuZW1haWwuRW'
    '1haWxBZGRyZXNzSW5zaWdodHNWZXJkaWN0Ug5oYXN2YWxpZHN5bnRheBJJCgxpc2Rpc3Bvc2Fi'
    'bGUYlPnNViABKAsyIi5lbWFpbC5FbWFpbEFkZHJlc3NJbnNpZ2h0c1ZlcmRpY3RSDGlzZGlzcG'
    '9zYWJsZRJLCg1pc3JhbmRvbWlucHV0GNe7oHkgASgLMiIuZW1haWwuRW1haWxBZGRyZXNzSW5z'
    'aWdodHNWZXJkaWN0Ug1pc3JhbmRvbWlucHV0EksKDWlzcm9sZWFkZHJlc3MYuKOdHSABKAsyIi'
    '5lbWFpbC5FbWFpbEFkZHJlc3NJbnNpZ2h0c1ZlcmRpY3RSDWlzcm9sZWFkZHJlc3MSSwoNbWFp'
    'bGJveGV4aXN0cxj89bAnIAEoCzIiLmVtYWlsLkVtYWlsQWRkcmVzc0luc2lnaHRzVmVyZGljdF'
    'INbWFpbGJveGV4aXN0cw==');

@$core.Deprecated('Use emailAddressInsightsVerdictDescriptor instead')
const EmailAddressInsightsVerdict$json = {
  '1': 'EmailAddressInsightsVerdict',
  '2': [
    {
      '1': 'confidenceverdict',
      '3': 145545115,
      '4': 1,
      '5': 14,
      '6': '.email.EmailAddressInsightsConfidenceVerdict',
      '10': 'confidenceverdict'
    },
  ],
};

/// Descriptor for `EmailAddressInsightsVerdict`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailAddressInsightsVerdictDescriptor =
    $convert.base64Decode(
        'ChtFbWFpbEFkZHJlc3NJbnNpZ2h0c1ZlcmRpY3QSXQoRY29uZmlkZW5jZXZlcmRpY3QYm6+zRS'
        'ABKA4yLC5lbWFpbC5FbWFpbEFkZHJlc3NJbnNpZ2h0c0NvbmZpZGVuY2VWZXJkaWN0UhFjb25m'
        'aWRlbmNldmVyZGljdA==');

@$core.Deprecated('Use emailContentDescriptor instead')
const EmailContent$json = {
  '1': 'EmailContent',
  '2': [
    {
      '1': 'raw',
      '3': 109685394,
      '4': 1,
      '5': 11,
      '6': '.email.RawMessage',
      '10': 'raw'
    },
    {
      '1': 'simple',
      '3': 187921568,
      '4': 1,
      '5': 11,
      '6': '.email.Message',
      '10': 'simple'
    },
    {
      '1': 'template',
      '3': 93743916,
      '4': 1,
      '5': 11,
      '6': '.email.Template',
      '10': 'template'
    },
  ],
};

/// Descriptor for `EmailContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailContentDescriptor = $convert.base64Decode(
    'CgxFbWFpbENvbnRlbnQSJgoDcmF3GJLVpjQgASgLMhEuZW1haWwuUmF3TWVzc2FnZVIDcmF3Ei'
    'kKBnNpbXBsZRig6c1ZIAEoCzIOLmVtYWlsLk1lc3NhZ2VSBnNpbXBsZRIuCgh0ZW1wbGF0ZRis'
    '1tksIAEoCzIPLmVtYWlsLlRlbXBsYXRlUgh0ZW1wbGF0ZQ==');

@$core.Deprecated('Use emailInsightsDescriptor instead')
const EmailInsights$json = {
  '1': 'EmailInsights',
  '2': [
    {'1': 'destination', '3': 457443680, '4': 1, '5': 9, '10': 'destination'},
    {
      '1': 'events',
      '3': 3416229,
      '4': 3,
      '5': 11,
      '6': '.email.InsightsEvent',
      '10': 'events'
    },
    {'1': 'isp', '3': 177492896, '4': 1, '5': 9, '10': 'isp'},
  ],
};

/// Descriptor for `EmailInsights`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailInsightsDescriptor = $convert.base64Decode(
    'Cg1FbWFpbEluc2lnaHRzEiQKC2Rlc3RpbmF0aW9uGOCSkNoBIAEoCVILZGVzdGluYXRpb24SLw'
    'oGZXZlbnRzGKXB0AEgAygLMhQuZW1haWwuSW5zaWdodHNFdmVudFIGZXZlbnRzEhMKA2lzcBig'
    'p9FUIAEoCVIDaXNw');

@$core.Deprecated('Use emailTemplateContentDescriptor instead')
const EmailTemplateContent$json = {
  '1': 'EmailTemplateContent',
  '2': [
    {'1': 'html', '3': 396489713, '4': 1, '5': 9, '10': 'html'},
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
    {'1': 'text', '3': 504638815, '4': 1, '5': 9, '10': 'text'},
  ],
};

/// Descriptor for `EmailTemplateContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailTemplateContentDescriptor = $convert.base64Decode(
    'ChRFbWFpbFRlbXBsYXRlQ29udGVudBIWCgRodG1sGPHnh70BIAEoCVIEaHRtbBIbCgdzdWJqZW'
    'N0GPDJ5AMgASgJUgdzdWJqZWN0EhYKBHRleHQY39rQ8AEgASgJUgR0ZXh0');

@$core.Deprecated('Use emailTemplateMetadataDescriptor instead')
const EmailTemplateMetadata$json = {
  '1': 'EmailTemplateMetadata',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `EmailTemplateMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List emailTemplateMetadataDescriptor = $convert.base64Decode(
    'ChVFbWFpbFRlbXBsYXRlTWV0YWRhdGESLgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEG'
    'NyZWF0ZWR0aW1lc3RhbXASJgoMdGVtcGxhdGVuYW1lGLGYkeUBIAEoCVIMdGVtcGxhdGVuYW1l');

@$core.Deprecated('Use eventBridgeDestinationDescriptor instead')
const EventBridgeDestination$json = {
  '1': 'EventBridgeDestination',
  '2': [
    {'1': 'eventbusarn', '3': 515755843, '4': 1, '5': 9, '10': 'eventbusarn'},
  ],
};

/// Descriptor for `EventBridgeDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventBridgeDestinationDescriptor =
    $convert.base64Decode(
        'ChZFdmVudEJyaWRnZURlc3RpbmF0aW9uEiQKC2V2ZW50YnVzYXJuGMOe9/UBIAEoCVILZXZlbn'
        'RidXNhcm4=');

@$core.Deprecated('Use eventDestinationDescriptor instead')
const EventDestination$json = {
  '1': 'EventDestination',
  '2': [
    {
      '1': 'cloudwatchdestination',
      '3': 124124536,
      '4': 1,
      '5': 11,
      '6': '.email.CloudWatchDestination',
      '10': 'cloudwatchdestination'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'eventbridgedestination',
      '3': 304968019,
      '4': 1,
      '5': 11,
      '6': '.email.EventBridgeDestination',
      '10': 'eventbridgedestination'
    },
    {
      '1': 'kinesisfirehosedestination',
      '3': 393614149,
      '4': 1,
      '5': 11,
      '6': '.email.KinesisFirehoseDestination',
      '10': 'kinesisfirehosedestination'
    },
    {
      '1': 'matchingeventtypes',
      '3': 888528,
      '4': 3,
      '5': 14,
      '6': '.email.EventType',
      '10': 'matchingeventtypes'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'pinpointdestination',
      '3': 448365711,
      '4': 1,
      '5': 11,
      '6': '.email.PinpointDestination',
      '10': 'pinpointdestination'
    },
    {
      '1': 'snsdestination',
      '3': 99189968,
      '4': 1,
      '5': 11,
      '6': '.email.SnsDestination',
      '10': 'snsdestination'
    },
  ],
};

/// Descriptor for `EventDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDestinationDescriptor = $convert.base64Decode(
    'ChBFdmVudERlc3RpbmF0aW9uElUKFWNsb3Vkd2F0Y2hkZXN0aW5hdGlvbhj4+pc7IAEoCzIcLm'
    'VtYWlsLkNsb3VkV2F0Y2hEZXN0aW5hdGlvblIVY2xvdWR3YXRjaGRlc3RpbmF0aW9uEhwKB2Vu'
    'YWJsZWQYv8ib5AEgASgIUgdlbmFibGVkElkKFmV2ZW50YnJpZGdlZGVzdGluYXRpb24Y0+K1kQ'
    'EgASgLMh0uZW1haWwuRXZlbnRCcmlkZ2VEZXN0aW5hdGlvblIWZXZlbnRicmlkZ2VkZXN0aW5h'
    'dGlvbhJlChpraW5lc2lzZmlyZWhvc2VkZXN0aW5hdGlvbhjFpti7ASABKAsyIS5lbWFpbC5LaW'
    '5lc2lzRmlyZWhvc2VEZXN0aW5hdGlvblIaa2luZXNpc2ZpcmVob3NlZGVzdGluYXRpb24SQgoS'
    'bWF0Y2hpbmdldmVudHR5cGVzGNCdNiADKA4yEC5lbWFpbC5FdmVudFR5cGVSEm1hdGNoaW5nZX'
    'ZlbnR0eXBlcxIVCgRuYW1lGIfmgX8gASgJUgRuYW1lElAKE3BpbnBvaW50ZGVzdGluYXRpb24Y'
    'j4nm1QEgASgLMhouZW1haWwuUGlucG9pbnREZXN0aW5hdGlvblITcGlucG9pbnRkZXN0aW5hdG'
    'lvbhJACg5zbnNkZXN0aW5hdGlvbhjQiaYvIAEoCzIVLmVtYWlsLlNuc0Rlc3RpbmF0aW9uUg5z'
    'bnNkZXN0aW5hdGlvbg==');

@$core.Deprecated('Use eventDestinationDefinitionDescriptor instead')
const EventDestinationDefinition$json = {
  '1': 'EventDestinationDefinition',
  '2': [
    {
      '1': 'cloudwatchdestination',
      '3': 124124536,
      '4': 1,
      '5': 11,
      '6': '.email.CloudWatchDestination',
      '10': 'cloudwatchdestination'
    },
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'eventbridgedestination',
      '3': 304968019,
      '4': 1,
      '5': 11,
      '6': '.email.EventBridgeDestination',
      '10': 'eventbridgedestination'
    },
    {
      '1': 'kinesisfirehosedestination',
      '3': 393614149,
      '4': 1,
      '5': 11,
      '6': '.email.KinesisFirehoseDestination',
      '10': 'kinesisfirehosedestination'
    },
    {
      '1': 'matchingeventtypes',
      '3': 888528,
      '4': 3,
      '5': 14,
      '6': '.email.EventType',
      '10': 'matchingeventtypes'
    },
    {
      '1': 'pinpointdestination',
      '3': 448365711,
      '4': 1,
      '5': 11,
      '6': '.email.PinpointDestination',
      '10': 'pinpointdestination'
    },
    {
      '1': 'snsdestination',
      '3': 99189968,
      '4': 1,
      '5': 11,
      '6': '.email.SnsDestination',
      '10': 'snsdestination'
    },
  ],
};

/// Descriptor for `EventDestinationDefinition`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDestinationDefinitionDescriptor = $convert.base64Decode(
    'ChpFdmVudERlc3RpbmF0aW9uRGVmaW5pdGlvbhJVChVjbG91ZHdhdGNoZGVzdGluYXRpb24Y+P'
    'qXOyABKAsyHC5lbWFpbC5DbG91ZFdhdGNoRGVzdGluYXRpb25SFWNsb3Vkd2F0Y2hkZXN0aW5h'
    'dGlvbhIcCgdlbmFibGVkGL/Im+QBIAEoCFIHZW5hYmxlZBJZChZldmVudGJyaWRnZWRlc3Rpbm'
    'F0aW9uGNPitZEBIAEoCzIdLmVtYWlsLkV2ZW50QnJpZGdlRGVzdGluYXRpb25SFmV2ZW50YnJp'
    'ZGdlZGVzdGluYXRpb24SZQoaa2luZXNpc2ZpcmVob3NlZGVzdGluYXRpb24YxabYuwEgASgLMi'
    'EuZW1haWwuS2luZXNpc0ZpcmVob3NlRGVzdGluYXRpb25SGmtpbmVzaXNmaXJlaG9zZWRlc3Rp'
    'bmF0aW9uEkIKEm1hdGNoaW5nZXZlbnR0eXBlcxjQnTYgAygOMhAuZW1haWwuRXZlbnRUeXBlUh'
    'JtYXRjaGluZ2V2ZW50dHlwZXMSUAoTcGlucG9pbnRkZXN0aW5hdGlvbhiPiebVASABKAsyGi5l'
    'bWFpbC5QaW5wb2ludERlc3RpbmF0aW9uUhNwaW5wb2ludGRlc3RpbmF0aW9uEkAKDnNuc2Rlc3'
    'RpbmF0aW9uGNCJpi8gASgLMhUuZW1haWwuU25zRGVzdGluYXRpb25SDnNuc2Rlc3RpbmF0aW9u');

@$core.Deprecated('Use eventDetailsDescriptor instead')
const EventDetails$json = {
  '1': 'EventDetails',
  '2': [
    {
      '1': 'bounce',
      '3': 54070366,
      '4': 1,
      '5': 11,
      '6': '.email.Bounce',
      '10': 'bounce'
    },
    {
      '1': 'complaint',
      '3': 108213147,
      '4': 1,
      '5': 11,
      '6': '.email.Complaint',
      '10': 'complaint'
    },
  ],
};

/// Descriptor for `EventDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventDetailsDescriptor = $convert.base64Decode(
    'CgxFdmVudERldGFpbHMSKAoGYm91bmNlGN6Y5BkgASgLMg0uZW1haWwuQm91bmNlUgZib3VuY2'
    'USMQoJY29tcGxhaW50GJvnzDMgASgLMhAuZW1haWwuQ29tcGxhaW50Ugljb21wbGFpbnQ=');

@$core.Deprecated('Use exportDataSourceDescriptor instead')
const ExportDataSource$json = {
  '1': 'ExportDataSource',
  '2': [
    {
      '1': 'messageinsightsdatasource',
      '3': 128258701,
      '4': 1,
      '5': 11,
      '6': '.email.MessageInsightsDataSource',
      '10': 'messageinsightsdatasource'
    },
    {
      '1': 'metricsdatasource',
      '3': 502846556,
      '4': 1,
      '5': 11,
      '6': '.email.MetricsDataSource',
      '10': 'metricsdatasource'
    },
  ],
};

/// Descriptor for `ExportDataSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportDataSourceDescriptor = $convert.base64Decode(
    'ChBFeHBvcnREYXRhU291cmNlEmEKGW1lc3NhZ2VpbnNpZ2h0c2RhdGFzb3VyY2UYjaWUPSABKA'
    'syIC5lbWFpbC5NZXNzYWdlSW5zaWdodHNEYXRhU291cmNlUhltZXNzYWdlaW5zaWdodHNkYXRh'
    'c291cmNlEkoKEW1ldHJpY3NkYXRhc291cmNlGNyo4+8BIAEoCzIYLmVtYWlsLk1ldHJpY3NEYX'
    'RhU291cmNlUhFtZXRyaWNzZGF0YXNvdXJjZQ==');

@$core.Deprecated('Use exportDestinationDescriptor instead')
const ExportDestination$json = {
  '1': 'ExportDestination',
  '2': [
    {
      '1': 'dataformat',
      '3': 89652083,
      '4': 1,
      '5': 14,
      '6': '.email.DataFormat',
      '10': 'dataformat'
    },
    {'1': 's3url', '3': 407792021, '4': 1, '5': 9, '10': 's3url'},
  ],
};

/// Descriptor for `ExportDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportDestinationDescriptor = $convert.base64Decode(
    'ChFFeHBvcnREZXN0aW5hdGlvbhI0CgpkYXRhZm9ybWF0GPP23yogASgOMhEuZW1haWwuRGF0YU'
    'Zvcm1hdFIKZGF0YWZvcm1hdBIYCgVzM3VybBiV07nCASABKAlSBXMzdXJs');

@$core.Deprecated('Use exportJobSummaryDescriptor instead')
const ExportJobSummary$json = {
  '1': 'ExportJobSummary',
  '2': [
    {
      '1': 'completedtimestamp',
      '3': 333312139,
      '4': 1,
      '5': 9,
      '10': 'completedtimestamp'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'exportsourcetype',
      '3': 248243607,
      '4': 1,
      '5': 14,
      '6': '.email.ExportSourceType',
      '10': 'exportsourcetype'
    },
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.email.JobStatus',
      '10': 'jobstatus'
    },
  ],
};

/// Descriptor for `ExportJobSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportJobSummaryDescriptor = $convert.base64Decode(
    'ChBFeHBvcnRKb2JTdW1tYXJ5EjIKEmNvbXBsZXRlZHRpbWVzdGFtcBiL4feeASABKAlSEmNvbX'
    'BsZXRlZHRpbWVzdGFtcBIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCVIQY3JlYXRlZHRp'
    'bWVzdGFtcBJGChBleHBvcnRzb3VyY2V0eXBlGJfLr3YgASgOMhcuZW1haWwuRXhwb3J0U291cm'
    'NlVHlwZVIQZXhwb3J0c291cmNldHlwZRIXCgVqb2JpZBjS1N0zIAEoCVIFam9iaWQSMQoJam9i'
    'c3RhdHVzGMec+zMgASgOMhAuZW1haWwuSm9iU3RhdHVzUglqb2JzdGF0dXM=');

@$core.Deprecated('Use exportMetricDescriptor instead')
const ExportMetric$json = {
  '1': 'ExportMetric',
  '2': [
    {
      '1': 'aggregation',
      '3': 132460038,
      '4': 1,
      '5': 14,
      '6': '.email.MetricAggregation',
      '10': 'aggregation'
    },
    {
      '1': 'name',
      '3': 266367751,
      '4': 1,
      '5': 14,
      '6': '.email.Metric',
      '10': 'name'
    },
  ],
};

/// Descriptor for `ExportMetric`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportMetricDescriptor = $convert.base64Decode(
    'CgxFeHBvcnRNZXRyaWMSPQoLYWdncmVnYXRpb24YhtyUPyABKA4yGC5lbWFpbC5NZXRyaWNBZ2'
    'dyZWdhdGlvblILYWdncmVnYXRpb24SJAoEbmFtZRiH5oF/IAEoDjINLmVtYWlsLk1ldHJpY1IE'
    'bmFtZQ==');

@$core.Deprecated('Use exportStatisticsDescriptor instead')
const ExportStatistics$json = {
  '1': 'ExportStatistics',
  '2': [
    {
      '1': 'exportedrecordscount',
      '3': 401601812,
      '4': 1,
      '5': 5,
      '10': 'exportedrecordscount'
    },
    {
      '1': 'processedrecordscount',
      '3': 507944491,
      '4': 1,
      '5': 5,
      '10': 'processedrecordscount'
    },
  ],
};

/// Descriptor for `ExportStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List exportStatisticsDescriptor = $convert.base64Decode(
    'ChBFeHBvcnRTdGF0aXN0aWNzEjYKFGV4cG9ydGVkcmVjb3Jkc2NvdW50GJTqv78BIAEoBVIUZX'
    'hwb3J0ZWRyZWNvcmRzY291bnQSOAoVcHJvY2Vzc2VkcmVjb3Jkc2NvdW50GKu8mvIBIAEoBVIV'
    'cHJvY2Vzc2VkcmVjb3Jkc2NvdW50');

@$core.Deprecated('Use failureInfoDescriptor instead')
const FailureInfo$json = {
  '1': 'FailureInfo',
  '2': [
    {'1': 'errormessage', '3': 518702377, '4': 1, '5': 9, '10': 'errormessage'},
    {
      '1': 'failedrecordss3url',
      '3': 275984286,
      '4': 1,
      '5': 9,
      '10': 'failedrecordss3url'
    },
  ],
};

/// Descriptor for `FailureInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List failureInfoDescriptor = $convert.base64Decode(
    'CgtGYWlsdXJlSW5mbxImCgxlcnJvcm1lc3NhZ2UYqYqr9wEgASgJUgxlcnJvcm1lc3NhZ2USMg'
    'oSZmFpbGVkcmVjb3Jkc3MzdXJsGJ7fzIMBIAEoCVISZmFpbGVkcmVjb3Jkc3MzdXJs');

@$core.Deprecated('Use getAccountRequestDescriptor instead')
const GetAccountRequest$json = {
  '1': 'GetAccountRequest',
};

/// Descriptor for `GetAccountRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountRequestDescriptor =
    $convert.base64Decode('ChFHZXRBY2NvdW50UmVxdWVzdA==');

@$core.Deprecated('Use getAccountResponseDescriptor instead')
const GetAccountResponse$json = {
  '1': 'GetAccountResponse',
  '2': [
    {
      '1': 'dedicatedipautowarmupenabled',
      '3': 96331178,
      '4': 1,
      '5': 8,
      '10': 'dedicatedipautowarmupenabled'
    },
    {
      '1': 'details',
      '3': 247611974,
      '4': 1,
      '5': 11,
      '6': '.email.AccountDetails',
      '10': 'details'
    },
    {
      '1': 'enforcementstatus',
      '3': 367265820,
      '4': 1,
      '5': 9,
      '10': 'enforcementstatus'
    },
    {
      '1': 'productionaccessenabled',
      '3': 471167534,
      '4': 1,
      '5': 8,
      '10': 'productionaccessenabled'
    },
    {
      '1': 'sendquota',
      '3': 486653062,
      '4': 1,
      '5': 11,
      '6': '.email.SendQuota',
      '10': 'sendquota'
    },
    {
      '1': 'sendingenabled',
      '3': 194846115,
      '4': 1,
      '5': 8,
      '10': 'sendingenabled'
    },
    {
      '1': 'suppressionattributes',
      '3': 17626546,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionAttributes',
      '10': 'suppressionattributes'
    },
    {
      '1': 'vdmattributes',
      '3': 506490540,
      '4': 1,
      '5': 11,
      '6': '.email.VdmAttributes',
      '10': 'vdmattributes'
    },
  ],
};

/// Descriptor for `GetAccountResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getAccountResponseDescriptor = $convert.base64Decode(
    'ChJHZXRBY2NvdW50UmVzcG9uc2USRQocZGVkaWNhdGVkaXBhdXRvd2FybXVwZW5hYmxlZBiqy/'
    'ctIAEoCFIcZGVkaWNhdGVkaXBhdXRvd2FybXVwZW5hYmxlZBIyCgdkZXRhaWxzGMaEiXYgASgL'
    'MhUuZW1haWwuQWNjb3VudERldGFpbHNSB2RldGFpbHMSMAoRZW5mb3JjZW1lbnRzdGF0dXMYnJ'
    'CQrwEgASgJUhFlbmZvcmNlbWVudHN0YXR1cxI8Chdwcm9kdWN0aW9uYWNjZXNzZW5hYmxlZBiu'
    '5NXgASABKAhSF3Byb2R1Y3Rpb25hY2Nlc3NlbmFibGVkEjIKCXNlbmRxdW90YRiG+YboASABKA'
    'syEC5lbWFpbC5TZW5kUXVvdGFSCXNlbmRxdW90YRIpCg5zZW5kaW5nZW5hYmxlZBiju/RcIAEo'
    'CFIOc2VuZGluZ2VuYWJsZWQSVQoVc3VwcHJlc3Npb25hdHRyaWJ1dGVzGLLrswggASgLMhwuZW'
    '1haWwuU3VwcHJlc3Npb25BdHRyaWJ1dGVzUhVzdXBwcmVzc2lvbmF0dHJpYnV0ZXMSPgoNdmRt'
    'YXR0cmlidXRlcxis3cHxASABKAsyFC5lbWFpbC5WZG1BdHRyaWJ1dGVzUg12ZG1hdHRyaWJ1dG'
    'Vz');

@$core.Deprecated('Use getBlacklistReportsRequestDescriptor instead')
const GetBlacklistReportsRequest$json = {
  '1': 'GetBlacklistReportsRequest',
  '2': [
    {
      '1': 'blacklistitemnames',
      '3': 434139554,
      '4': 3,
      '5': 9,
      '10': 'blacklistitemnames'
    },
  ],
};

/// Descriptor for `GetBlacklistReportsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlacklistReportsRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRCbGFja2xpc3RSZXBvcnRzUmVxdWVzdBIyChJibGFja2xpc3RpdGVtbmFtZXMYouOBzw'
        'EgAygJUhJibGFja2xpc3RpdGVtbmFtZXM=');

@$core.Deprecated('Use getBlacklistReportsResponseDescriptor instead')
const GetBlacklistReportsResponse$json = {
  '1': 'GetBlacklistReportsResponse',
  '2': [
    {
      '1': 'blacklistreport',
      '3': 64005855,
      '4': 3,
      '5': 11,
      '6': '.email.GetBlacklistReportsResponse.BlacklistreportEntry',
      '10': 'blacklistreport'
    },
  ],
  '3': [GetBlacklistReportsResponse_BlacklistreportEntry$json],
};

@$core.Deprecated('Use getBlacklistReportsResponseDescriptor instead')
const GetBlacklistReportsResponse_BlacklistreportEntry$json = {
  '1': 'BlacklistreportEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetBlacklistReportsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getBlacklistReportsResponseDescriptor = $convert.base64Decode(
    'ChtHZXRCbGFja2xpc3RSZXBvcnRzUmVzcG9uc2USZAoPYmxhY2tsaXN0cmVwb3J0GN/Nwh4gAy'
    'gLMjcuZW1haWwuR2V0QmxhY2tsaXN0UmVwb3J0c1Jlc3BvbnNlLkJsYWNrbGlzdHJlcG9ydEVu'
    'dHJ5Ug9ibGFja2xpc3RyZXBvcnQaQgoUQmxhY2tsaXN0cmVwb3J0RW50cnkSEAoDa2V5GAEgAS'
    'gJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated(
    'Use getConfigurationSetEventDestinationsRequestDescriptor instead')
const GetConfigurationSetEventDestinationsRequest$json = {
  '1': 'GetConfigurationSetEventDestinationsRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
  ],
};

/// Descriptor for `GetConfigurationSetEventDestinationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getConfigurationSetEventDestinationsRequestDescriptor =
    $convert.base64Decode(
        'CitHZXRDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvbnNSZXF1ZXN0EjYKFGNvbmZpZ3'
        'VyYXRpb25zZXRuYW1lGM2LscABIAEoCVIUY29uZmlndXJhdGlvbnNldG5hbWU=');

@$core.Deprecated(
    'Use getConfigurationSetEventDestinationsResponseDescriptor instead')
const GetConfigurationSetEventDestinationsResponse$json = {
  '1': 'GetConfigurationSetEventDestinationsResponse',
  '2': [
    {
      '1': 'eventdestinations',
      '3': 484944019,
      '4': 3,
      '5': 11,
      '6': '.email.EventDestination',
      '10': 'eventdestinations'
    },
  ],
};

/// Descriptor for `GetConfigurationSetEventDestinationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getConfigurationSetEventDestinationsResponseDescriptor =
    $convert.base64Decode(
        'CixHZXRDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvbnNSZXNwb25zZRJJChFldmVudG'
        'Rlc3RpbmF0aW9ucxiT0Z7nASADKAsyFy5lbWFpbC5FdmVudERlc3RpbmF0aW9uUhFldmVudGRl'
        'c3RpbmF0aW9ucw==');

@$core.Deprecated('Use getConfigurationSetRequestDescriptor instead')
const GetConfigurationSetRequest$json = {
  '1': 'GetConfigurationSetRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
  ],
};

/// Descriptor for `GetConfigurationSetRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConfigurationSetRequestDescriptor =
    $convert.base64Decode(
        'ChpHZXRDb25maWd1cmF0aW9uU2V0UmVxdWVzdBI2ChRjb25maWd1cmF0aW9uc2V0bmFtZRjNi7'
        'HAASABKAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1l');

@$core.Deprecated('Use getConfigurationSetResponseDescriptor instead')
const GetConfigurationSetResponse$json = {
  '1': 'GetConfigurationSetResponse',
  '2': [
    {
      '1': 'archivingoptions',
      '3': 54783281,
      '4': 1,
      '5': 11,
      '6': '.email.ArchivingOptions',
      '10': 'archivingoptions'
    },
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'deliveryoptions',
      '3': 332764278,
      '4': 1,
      '5': 11,
      '6': '.email.DeliveryOptions',
      '10': 'deliveryoptions'
    },
    {
      '1': 'reputationoptions',
      '3': 98302085,
      '4': 1,
      '5': 11,
      '6': '.email.ReputationOptions',
      '10': 'reputationoptions'
    },
    {
      '1': 'sendingoptions',
      '3': 319193106,
      '4': 1,
      '5': 11,
      '6': '.email.SendingOptions',
      '10': 'sendingoptions'
    },
    {
      '1': 'suppressionoptions',
      '3': 96075559,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionOptions',
      '10': 'suppressionoptions'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'trackingoptions',
      '3': 237447811,
      '4': 1,
      '5': 11,
      '6': '.email.TrackingOptions',
      '10': 'trackingoptions'
    },
    {
      '1': 'vdmoptions',
      '3': 344971073,
      '4': 1,
      '5': 11,
      '6': '.email.VdmOptions',
      '10': 'vdmoptions'
    },
  ],
};

/// Descriptor for `GetConfigurationSetResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getConfigurationSetResponseDescriptor = $convert.base64Decode(
    'ChtHZXRDb25maWd1cmF0aW9uU2V0UmVzcG9uc2USRgoQYXJjaGl2aW5nb3B0aW9ucxix2o8aIA'
    'EoCzIXLmVtYWlsLkFyY2hpdmluZ09wdGlvbnNSEGFyY2hpdmluZ29wdGlvbnMSNgoUY29uZmln'
    'dXJhdGlvbnNldG5hbWUYzYuxwAEgASgJUhRjb25maWd1cmF0aW9uc2V0bmFtZRJECg9kZWxpdm'
    'VyeW9wdGlvbnMY9qjWngEgASgLMhYuZW1haWwuRGVsaXZlcnlPcHRpb25zUg9kZWxpdmVyeW9w'
    'dGlvbnMSSQoRcmVwdXRhdGlvbm9wdGlvbnMYhfHvLiABKAsyGC5lbWFpbC5SZXB1dGF0aW9uT3'
    'B0aW9uc1IRcmVwdXRhdGlvbm9wdGlvbnMSQQoOc2VuZGluZ29wdGlvbnMYkoCamAEgASgLMhUu'
    'ZW1haWwuU2VuZGluZ09wdGlvbnNSDnNlbmRpbmdvcHRpb25zEkwKEnN1cHByZXNzaW9ub3B0aW'
    '9ucxin/uctIAEoCzIZLmVtYWlsLlN1cHByZXNzaW9uT3B0aW9uc1ISc3VwcHJlc3Npb25vcHRp'
    'b25zEiIKBHRhZ3MYwcH2tQEgAygLMgouZW1haWwuVGFnUgR0YWdzEkMKD3RyYWNraW5nb3B0aW'
    '9ucxiD1ZxxIAEoCzIWLmVtYWlsLlRyYWNraW5nT3B0aW9uc1IPdHJhY2tpbmdvcHRpb25zEjUK'
    'CnZkbW9wdGlvbnMYwa6/pAEgASgLMhEuZW1haWwuVmRtT3B0aW9uc1IKdmRtb3B0aW9ucw==');

@$core.Deprecated('Use getContactListRequestDescriptor instead')
const GetContactListRequest$json = {
  '1': 'GetContactListRequest',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
  ],
};

/// Descriptor for `GetContactListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContactListRequestDescriptor = $convert.base64Decode(
    'ChVHZXRDb250YWN0TGlzdFJlcXVlc3QSLAoPY29udGFjdGxpc3RuYW1lGPG9p6EBIAEoCVIPY2'
    '9udGFjdGxpc3RuYW1l');

@$core.Deprecated('Use getContactListResponseDescriptor instead')
const GetContactListResponse$json = {
  '1': 'GetContactListResponse',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'topics',
      '3': 219850038,
      '4': 3,
      '5': 11,
      '6': '.email.Topic',
      '10': 'topics'
    },
  ],
};

/// Descriptor for `GetContactListResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContactListResponseDescriptor = $convert.base64Decode(
    'ChZHZXRDb250YWN0TGlzdFJlc3BvbnNlEiwKD2NvbnRhY3RsaXN0bmFtZRjxvaehASABKAlSD2'
    'NvbnRhY3RsaXN0bmFtZRIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCVIQY3JlYXRlZHRp'
    'bWVzdGFtcBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SNQoUbGFzdHVwZG'
    'F0ZWR0aW1lc3RhbXAYlcvIPyABKAlSFGxhc3R1cGRhdGVkdGltZXN0YW1wEiIKBHRhZ3MYwcH2'
    'tQEgAygLMgouZW1haWwuVGFnUgR0YWdzEicKBnRvcGljcxi2yupoIAMoCzIMLmVtYWlsLlRvcG'
    'ljUgZ0b3BpY3M=');

@$core.Deprecated('Use getContactRequestDescriptor instead')
const GetContactRequest$json = {
  '1': 'GetContactRequest',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
  ],
};

/// Descriptor for `GetContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContactRequestDescriptor = $convert.base64Decode(
    'ChFHZXRDb250YWN0UmVxdWVzdBIsCg9jb250YWN0bGlzdG5hbWUY8b2noQEgASgJUg9jb250YW'
    'N0bGlzdG5hbWUSJgoMZW1haWxhZGRyZXNzGNbxiukBIAEoCVIMZW1haWxhZGRyZXNz');

@$core.Deprecated('Use getContactResponseDescriptor instead')
const GetContactResponse$json = {
  '1': 'GetContactResponse',
  '2': [
    {
      '1': 'attributesdata',
      '3': 497271421,
      '4': 1,
      '5': 9,
      '10': 'attributesdata'
    },
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {
      '1': 'topicdefaultpreferences',
      '3': 477710170,
      '4': 3,
      '5': 11,
      '6': '.email.TopicPreference',
      '10': 'topicdefaultpreferences'
    },
    {
      '1': 'topicpreferences',
      '3': 199822983,
      '4': 3,
      '5': 11,
      '6': '.email.TopicPreference',
      '10': 'topicpreferences'
    },
    {
      '1': 'unsubscribeall',
      '3': 49261174,
      '4': 1,
      '5': 8,
      '10': 'unsubscribeall'
    },
  ],
};

/// Descriptor for `GetContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getContactResponseDescriptor = $convert.base64Decode(
    'ChJHZXRDb250YWN0UmVzcG9uc2USKgoOYXR0cmlidXRlc2RhdGEY/YSP7QEgASgJUg5hdHRyaW'
    'J1dGVzZGF0YRIsCg9jb250YWN0bGlzdG5hbWUY8b2noQEgASgJUg9jb250YWN0bGlzdG5hbWUS'
    'LgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEGNyZWF0ZWR0aW1lc3RhbXASJgoMZW1haW'
    'xhZGRyZXNzGNbxiukBIAEoCVIMZW1haWxhZGRyZXNzEjUKFGxhc3R1cGRhdGVkdGltZXN0YW1w'
    'GJXLyD8gASgJUhRsYXN0dXBkYXRlZHRpbWVzdGFtcBJUChd0b3BpY2RlZmF1bHRwcmVmZXJlbm'
    'NlcxjajuXjASADKAsyFi5lbWFpbC5Ub3BpY1ByZWZlcmVuY2VSF3RvcGljZGVmYXVsdHByZWZl'
    'cmVuY2VzEkUKEHRvcGljcHJlZmVyZW5jZXMYh52kXyADKAsyFi5lbWFpbC5Ub3BpY1ByZWZlcm'
    'VuY2VSEHRvcGljcHJlZmVyZW5jZXMSKQoOdW5zdWJzY3JpYmVhbGwY9tS+FyABKAhSDnVuc3Vi'
    'c2NyaWJlYWxs');

@$core.Deprecated(
    'Use getCustomVerificationEmailTemplateRequestDescriptor instead')
const GetCustomVerificationEmailTemplateRequest$json = {
  '1': 'GetCustomVerificationEmailTemplateRequest',
  '2': [
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `GetCustomVerificationEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getCustomVerificationEmailTemplateRequestDescriptor = $convert.base64Decode(
        'CilHZXRDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVxdWVzdBImCgx0ZW1wbGF0ZW'
        '5hbWUYsZiR5QEgASgJUgx0ZW1wbGF0ZW5hbWU=');

@$core.Deprecated(
    'Use getCustomVerificationEmailTemplateResponseDescriptor instead')
const GetCustomVerificationEmailTemplateResponse$json = {
  '1': 'GetCustomVerificationEmailTemplateResponse',
  '2': [
    {
      '1': 'failureredirectionurl',
      '3': 376073007,
      '4': 1,
      '5': 9,
      '10': 'failureredirectionurl'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {
      '1': 'successredirectionurl',
      '3': 513750768,
      '4': 1,
      '5': 9,
      '10': 'successredirectionurl'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'templatecontent',
      '3': 528973453,
      '4': 1,
      '5': 9,
      '10': 'templatecontent'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
    {
      '1': 'templatesubject',
      '3': 144653738,
      '4': 1,
      '5': 9,
      '10': 'templatesubject'
    },
  ],
};

/// Descriptor for `GetCustomVerificationEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    getCustomVerificationEmailTemplateResponseDescriptor =
    $convert.base64Decode(
        'CipHZXRDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVzcG9uc2USOAoVZmFpbHVyZX'
        'JlZGlyZWN0aW9udXJsGK/WqbMBIAEoCVIVZmFpbHVyZXJlZGlyZWN0aW9udXJsEi0KEGZyb21l'
        'bWFpbGFkZHJlc3MYhprLLCABKAlSEGZyb21lbWFpbGFkZHJlc3MSOAoVc3VjY2Vzc3JlZGlyZW'
        'N0aW9udXJsGPDt/PQBIAEoCVIVc3VjY2Vzc3JlZGlyZWN0aW9udXJsEiIKBHRhZ3MYwcH2tQEg'
        'AygLMgouZW1haWwuVGFnUgR0YWdzEiwKD3RlbXBsYXRlY29udGVudBiN/Z38ASABKAlSD3RlbX'
        'BsYXRlY29udGVudBImCgx0ZW1wbGF0ZW5hbWUYsZiR5QEgASgJUgx0ZW1wbGF0ZW5hbWUSKwoP'
        'dGVtcGxhdGVzdWJqZWN0GKr7/EQgASgJUg90ZW1wbGF0ZXN1YmplY3Q=');

@$core.Deprecated('Use getDedicatedIpPoolRequestDescriptor instead')
const GetDedicatedIpPoolRequest$json = {
  '1': 'GetDedicatedIpPoolRequest',
  '2': [
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
  ],
};

/// Descriptor for `GetDedicatedIpPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDedicatedIpPoolRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXREZWRpY2F0ZWRJcFBvb2xSZXF1ZXN0Eh0KCHBvb2xuYW1lGMmNhScgASgJUghwb29sbm'
        'FtZQ==');

@$core.Deprecated('Use getDedicatedIpPoolResponseDescriptor instead')
const GetDedicatedIpPoolResponse$json = {
  '1': 'GetDedicatedIpPoolResponse',
  '2': [
    {
      '1': 'dedicatedippool',
      '3': 272906882,
      '4': 1,
      '5': 11,
      '6': '.email.DedicatedIpPool',
      '10': 'dedicatedippool'
    },
  ],
};

/// Descriptor for `GetDedicatedIpPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDedicatedIpPoolResponseDescriptor =
    $convert.base64Decode(
        'ChpHZXREZWRpY2F0ZWRJcFBvb2xSZXNwb25zZRJECg9kZWRpY2F0ZWRpcHBvb2wYgvWQggEgAS'
        'gLMhYuZW1haWwuRGVkaWNhdGVkSXBQb29sUg9kZWRpY2F0ZWRpcHBvb2w=');

@$core.Deprecated('Use getDedicatedIpRequestDescriptor instead')
const GetDedicatedIpRequest$json = {
  '1': 'GetDedicatedIpRequest',
  '2': [
    {'1': 'ip', '3': 183031933, '4': 1, '5': 9, '10': 'ip'},
  ],
};

/// Descriptor for `GetDedicatedIpRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDedicatedIpRequestDescriptor = $convert
    .base64Decode('ChVHZXREZWRpY2F0ZWRJcFJlcXVlc3QSEQoCaXAY/bCjVyABKAlSAmlw');

@$core.Deprecated('Use getDedicatedIpResponseDescriptor instead')
const GetDedicatedIpResponse$json = {
  '1': 'GetDedicatedIpResponse',
  '2': [
    {
      '1': 'dedicatedip',
      '3': 472460244,
      '4': 1,
      '5': 11,
      '6': '.email.DedicatedIp',
      '10': 'dedicatedip'
    },
  ],
};

/// Descriptor for `GetDedicatedIpResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDedicatedIpResponseDescriptor =
    $convert.base64Decode(
        'ChZHZXREZWRpY2F0ZWRJcFJlc3BvbnNlEjgKC2RlZGljYXRlZGlwGNTXpOEBIAEoCzISLmVtYW'
        'lsLkRlZGljYXRlZElwUgtkZWRpY2F0ZWRpcA==');

@$core.Deprecated('Use getDedicatedIpsRequestDescriptor instead')
const GetDedicatedIpsRequest$json = {
  '1': 'GetDedicatedIpsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
  ],
};

/// Descriptor for `GetDedicatedIpsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDedicatedIpsRequestDescriptor = $convert.base64Decode(
    'ChZHZXREZWRpY2F0ZWRJcHNSZXF1ZXN0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2'
    'VuEh4KCHBhZ2VzaXplGLiTgtEBIAEoBVIIcGFnZXNpemUSHQoIcG9vbG5hbWUYyY2FJyABKAlS'
    'CHBvb2xuYW1l');

@$core.Deprecated('Use getDedicatedIpsResponseDescriptor instead')
const GetDedicatedIpsResponse$json = {
  '1': 'GetDedicatedIpsResponse',
  '2': [
    {
      '1': 'dedicatedips',
      '3': 349154529,
      '4': 3,
      '5': 11,
      '6': '.email.DedicatedIp',
      '10': 'dedicatedips'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `GetDedicatedIpsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDedicatedIpsResponseDescriptor = $convert.base64Decode(
    'ChdHZXREZWRpY2F0ZWRJcHNSZXNwb25zZRI6CgxkZWRpY2F0ZWRpcHMY4dm+pgEgAygLMhIuZW'
    '1haWwuRGVkaWNhdGVkSXBSDGRlZGljYXRlZGlwcxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5l'
    'eHR0b2tlbg==');

@$core.Deprecated(
    'Use getDeliverabilityDashboardOptionsRequestDescriptor instead')
const GetDeliverabilityDashboardOptionsRequest$json = {
  '1': 'GetDeliverabilityDashboardOptionsRequest',
};

/// Descriptor for `GetDeliverabilityDashboardOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliverabilityDashboardOptionsRequestDescriptor =
    $convert.base64Decode(
        'CihHZXREZWxpdmVyYWJpbGl0eURhc2hib2FyZE9wdGlvbnNSZXF1ZXN0');

@$core.Deprecated(
    'Use getDeliverabilityDashboardOptionsResponseDescriptor instead')
const GetDeliverabilityDashboardOptionsResponse$json = {
  '1': 'GetDeliverabilityDashboardOptionsResponse',
  '2': [
    {
      '1': 'accountstatus',
      '3': 961735,
      '4': 1,
      '5': 14,
      '6': '.email.DeliverabilityDashboardAccountStatus',
      '10': 'accountstatus'
    },
    {
      '1': 'activesubscribeddomains',
      '3': 441853281,
      '4': 3,
      '5': 11,
      '6': '.email.DomainDeliverabilityTrackingOption',
      '10': 'activesubscribeddomains'
    },
    {
      '1': 'dashboardenabled',
      '3': 90846529,
      '4': 1,
      '5': 8,
      '10': 'dashboardenabled'
    },
    {
      '1': 'pendingexpirationsubscribeddomains',
      '3': 188697179,
      '4': 3,
      '5': 11,
      '6': '.email.DomainDeliverabilityTrackingOption',
      '10': 'pendingexpirationsubscribeddomains'
    },
    {
      '1': 'subscriptionexpirydate',
      '3': 381693284,
      '4': 1,
      '5': 9,
      '10': 'subscriptionexpirydate'
    },
  ],
};

/// Descriptor for `GetDeliverabilityDashboardOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliverabilityDashboardOptionsResponseDescriptor = $convert.base64Decode(
    'CilHZXREZWxpdmVyYWJpbGl0eURhc2hib2FyZE9wdGlvbnNSZXNwb25zZRJTCg1hY2NvdW50c3'
    'RhdHVzGMfZOiABKA4yKy5lbWFpbC5EZWxpdmVyYWJpbGl0eURhc2hib2FyZEFjY291bnRTdGF0'
    'dXNSDWFjY291bnRzdGF0dXMSZwoXYWN0aXZlc3Vic2NyaWJlZGRvbWFpbnMY4crY0gEgAygLMi'
    'kuZW1haWwuRG9tYWluRGVsaXZlcmFiaWxpdHlUcmFja2luZ09wdGlvblIXYWN0aXZlc3Vic2Ny'
    'aWJlZGRvbWFpbnMSLQoQZGFzaGJvYXJkZW5hYmxlZBjB6qgrIAEoCFIQZGFzaGJvYXJkZW5hYm'
    'xlZBJ8CiJwZW5kaW5nZXhwaXJhdGlvbnN1YnNjcmliZWRkb21haW5zGNuU/VkgAygLMikuZW1h'
    'aWwuRG9tYWluRGVsaXZlcmFiaWxpdHlUcmFja2luZ09wdGlvblIicGVuZGluZ2V4cGlyYXRpb2'
    '5zdWJzY3JpYmVkZG9tYWlucxI6ChZzdWJzY3JpcHRpb25leHBpcnlkYXRlGOTagLYBIAEoCVIW'
    'c3Vic2NyaXB0aW9uZXhwaXJ5ZGF0ZQ==');

@$core.Deprecated('Use getDeliverabilityTestReportRequestDescriptor instead')
const GetDeliverabilityTestReportRequest$json = {
  '1': 'GetDeliverabilityTestReportRequest',
  '2': [
    {'1': 'reportid', '3': 420903847, '4': 1, '5': 9, '10': 'reportid'},
  ],
};

/// Descriptor for `GetDeliverabilityTestReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliverabilityTestReportRequestDescriptor =
    $convert.base64Decode(
        'CiJHZXREZWxpdmVyYWJpbGl0eVRlc3RSZXBvcnRSZXF1ZXN0Eh4KCHJlcG9ydGlkGKf32cgBIA'
        'EoCVIIcmVwb3J0aWQ=');

@$core.Deprecated('Use getDeliverabilityTestReportResponseDescriptor instead')
const GetDeliverabilityTestReportResponse$json = {
  '1': 'GetDeliverabilityTestReportResponse',
  '2': [
    {
      '1': 'deliverabilitytestreport',
      '3': 77121923,
      '4': 1,
      '5': 11,
      '6': '.email.DeliverabilityTestReport',
      '10': 'deliverabilitytestreport'
    },
    {
      '1': 'ispplacements',
      '3': 334564328,
      '4': 3,
      '5': 11,
      '6': '.email.IspPlacement',
      '10': 'ispplacements'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
    {
      '1': 'overallplacement',
      '3': 259723246,
      '4': 1,
      '5': 11,
      '6': '.email.PlacementStatistics',
      '10': 'overallplacement'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `GetDeliverabilityTestReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDeliverabilityTestReportResponseDescriptor = $convert.base64Decode(
    'CiNHZXREZWxpdmVyYWJpbGl0eVRlc3RSZXBvcnRSZXNwb25zZRJeChhkZWxpdmVyYWJpbGl0eX'
    'Rlc3RyZXBvcnQYg5PjJCABKAsyHy5lbWFpbC5EZWxpdmVyYWJpbGl0eVRlc3RSZXBvcnRSGGRl'
    'bGl2ZXJhYmlsaXR5dGVzdHJlcG9ydBI9Cg1pc3BwbGFjZW1lbnRzGOiXxJ8BIAMoCzITLmVtYW'
    'lsLklzcFBsYWNlbWVudFINaXNwcGxhY2VtZW50cxIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNz'
    'YWdlEkkKEG92ZXJhbGxwbGFjZW1lbnQY7p/seyABKAsyGi5lbWFpbC5QbGFjZW1lbnRTdGF0aX'
    'N0aWNzUhBvdmVyYWxscGxhY2VtZW50EiIKBHRhZ3MYwcH2tQEgAygLMgouZW1haWwuVGFnUgR0'
    'YWdz');

@$core
    .Deprecated('Use getDomainDeliverabilityCampaignRequestDescriptor instead')
const GetDomainDeliverabilityCampaignRequest$json = {
  '1': 'GetDomainDeliverabilityCampaignRequest',
  '2': [
    {'1': 'campaignid', '3': 299715775, '4': 1, '5': 9, '10': 'campaignid'},
  ],
};

/// Descriptor for `GetDomainDeliverabilityCampaignRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDomainDeliverabilityCampaignRequestDescriptor =
    $convert.base64Decode(
        'CiZHZXREb21haW5EZWxpdmVyYWJpbGl0eUNhbXBhaWduUmVxdWVzdBIiCgpjYW1wYWlnbmlkGL'
        '+Z9Y4BIAEoCVIKY2FtcGFpZ25pZA==');

@$core
    .Deprecated('Use getDomainDeliverabilityCampaignResponseDescriptor instead')
const GetDomainDeliverabilityCampaignResponse$json = {
  '1': 'GetDomainDeliverabilityCampaignResponse',
  '2': [
    {
      '1': 'domaindeliverabilitycampaign',
      '3': 446153883,
      '4': 1,
      '5': 11,
      '6': '.email.DomainDeliverabilityCampaign',
      '10': 'domaindeliverabilitycampaign'
    },
  ],
};

/// Descriptor for `GetDomainDeliverabilityCampaignResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDomainDeliverabilityCampaignResponseDescriptor =
    $convert.base64Decode(
        'CidHZXREb21haW5EZWxpdmVyYWJpbGl0eUNhbXBhaWduUmVzcG9uc2USawocZG9tYWluZGVsaX'
        'ZlcmFiaWxpdHljYW1wYWlnbhibid/UASABKAsyIy5lbWFpbC5Eb21haW5EZWxpdmVyYWJpbGl0'
        'eUNhbXBhaWduUhxkb21haW5kZWxpdmVyYWJpbGl0eWNhbXBhaWdu');

@$core.Deprecated('Use getDomainStatisticsReportRequestDescriptor instead')
const GetDomainStatisticsReportRequest$json = {
  '1': 'GetDomainStatisticsReportRequest',
  '2': [
    {'1': 'domain', '3': 505186578, '4': 1, '5': 9, '10': 'domain'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
  ],
};

/// Descriptor for `GetDomainStatisticsReportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDomainStatisticsReportRequestDescriptor =
    $convert.base64Decode(
        'CiBHZXREb21haW5TdGF0aXN0aWNzUmVwb3J0UmVxdWVzdBIaCgZkb21haW4YkpLy8AEgASgJUg'
        'Zkb21haW4SGwoHZW5kZGF0ZRjPs/kkIAEoCVIHZW5kZGF0ZRIgCglzdGFydGRhdGUY/Pig1AEg'
        'ASgJUglzdGFydGRhdGU=');

@$core.Deprecated('Use getDomainStatisticsReportResponseDescriptor instead')
const GetDomainStatisticsReportResponse$json = {
  '1': 'GetDomainStatisticsReportResponse',
  '2': [
    {
      '1': 'dailyvolumes',
      '3': 473250602,
      '4': 3,
      '5': 11,
      '6': '.email.DailyVolume',
      '10': 'dailyvolumes'
    },
    {
      '1': 'overallvolume',
      '3': 138964909,
      '4': 1,
      '5': 11,
      '6': '.email.OverallVolume',
      '10': 'overallvolume'
    },
  ],
};

/// Descriptor for `GetDomainStatisticsReportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getDomainStatisticsReportResponseDescriptor =
    $convert.base64Decode(
        'CiFHZXREb21haW5TdGF0aXN0aWNzUmVwb3J0UmVzcG9uc2USOgoMZGFpbHl2b2x1bWVzGKr21O'
        'EBIAMoCzISLmVtYWlsLkRhaWx5Vm9sdW1lUgxkYWlseXZvbHVtZXMSPQoNb3ZlcmFsbHZvbHVt'
        'ZRit36FCIAEoCzIULmVtYWlsLk92ZXJhbGxWb2x1bWVSDW92ZXJhbGx2b2x1bWU=');

@$core.Deprecated('Use getEmailAddressInsightsRequestDescriptor instead')
const GetEmailAddressInsightsRequest$json = {
  '1': 'GetEmailAddressInsightsRequest',
  '2': [
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
  ],
};

/// Descriptor for `GetEmailAddressInsightsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailAddressInsightsRequestDescriptor =
    $convert.base64Decode(
        'Ch5HZXRFbWFpbEFkZHJlc3NJbnNpZ2h0c1JlcXVlc3QSJgoMZW1haWxhZGRyZXNzGNbxiukBIA'
        'EoCVIMZW1haWxhZGRyZXNz');

@$core.Deprecated('Use getEmailAddressInsightsResponseDescriptor instead')
const GetEmailAddressInsightsResponse$json = {
  '1': 'GetEmailAddressInsightsResponse',
  '2': [
    {
      '1': 'mailboxvalidation',
      '3': 38796029,
      '4': 1,
      '5': 11,
      '6': '.email.MailboxValidation',
      '10': 'mailboxvalidation'
    },
  ],
};

/// Descriptor for `GetEmailAddressInsightsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailAddressInsightsResponseDescriptor =
    $convert.base64Decode(
        'Ch9HZXRFbWFpbEFkZHJlc3NJbnNpZ2h0c1Jlc3BvbnNlEkkKEW1haWxib3h2YWxpZGF0aW9uGP'
        '31vxIgASgLMhguZW1haWwuTWFpbGJveFZhbGlkYXRpb25SEW1haWxib3h2YWxpZGF0aW9u');

@$core.Deprecated('Use getEmailIdentityPoliciesRequestDescriptor instead')
const GetEmailIdentityPoliciesRequest$json = {
  '1': 'GetEmailIdentityPoliciesRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
  ],
};

/// Descriptor for `GetEmailIdentityPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailIdentityPoliciesRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRFbWFpbElkZW50aXR5UG9saWNpZXNSZXF1ZXN0EicKDWVtYWlsaWRlbnRpdHkYmpTyQC'
        'ABKAlSDWVtYWlsaWRlbnRpdHk=');

@$core.Deprecated('Use getEmailIdentityPoliciesResponseDescriptor instead')
const GetEmailIdentityPoliciesResponse$json = {
  '1': 'GetEmailIdentityPoliciesResponse',
  '2': [
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.email.GetEmailIdentityPoliciesResponse.PoliciesEntry',
      '10': 'policies'
    },
  ],
  '3': [GetEmailIdentityPoliciesResponse_PoliciesEntry$json],
};

@$core.Deprecated('Use getEmailIdentityPoliciesResponseDescriptor instead')
const GetEmailIdentityPoliciesResponse_PoliciesEntry$json = {
  '1': 'PoliciesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetEmailIdentityPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailIdentityPoliciesResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXRFbWFpbElkZW50aXR5UG9saWNpZXNSZXNwb25zZRJUCghwb2xpY2llcxiYrIoTIAMoCz'
        'I1LmVtYWlsLkdldEVtYWlsSWRlbnRpdHlQb2xpY2llc1Jlc3BvbnNlLlBvbGljaWVzRW50cnlS'
        'CHBvbGljaWVzGjsKDVBvbGljaWVzRW50cnkSEAoDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAi'
        'ABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use getEmailIdentityRequestDescriptor instead')
const GetEmailIdentityRequest$json = {
  '1': 'GetEmailIdentityRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
  ],
};

/// Descriptor for `GetEmailIdentityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailIdentityRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRFbWFpbElkZW50aXR5UmVxdWVzdBInCg1lbWFpbGlkZW50aXR5GJqU8kAgASgJUg1lbW'
        'FpbGlkZW50aXR5');

@$core.Deprecated('Use getEmailIdentityResponseDescriptor instead')
const GetEmailIdentityResponse$json = {
  '1': 'GetEmailIdentityResponse',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'dkimattributes',
      '3': 256039632,
      '4': 1,
      '5': 11,
      '6': '.email.DkimAttributes',
      '10': 'dkimattributes'
    },
    {
      '1': 'feedbackforwardingstatus',
      '3': 304617468,
      '4': 1,
      '5': 8,
      '10': 'feedbackforwardingstatus'
    },
    {
      '1': 'identitytype',
      '3': 499274628,
      '4': 1,
      '5': 14,
      '6': '.email.IdentityType',
      '10': 'identitytype'
    },
    {
      '1': 'mailfromattributes',
      '3': 460277248,
      '4': 1,
      '5': 11,
      '6': '.email.MailFromAttributes',
      '10': 'mailfromattributes'
    },
    {
      '1': 'policies',
      '3': 40015384,
      '4': 3,
      '5': 11,
      '6': '.email.GetEmailIdentityResponse.PoliciesEntry',
      '10': 'policies'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'verificationinfo',
      '3': 66315551,
      '4': 1,
      '5': 11,
      '6': '.email.VerificationInfo',
      '10': 'verificationinfo'
    },
    {
      '1': 'verificationstatus',
      '3': 132712897,
      '4': 1,
      '5': 14,
      '6': '.email.VerificationStatus',
      '10': 'verificationstatus'
    },
    {
      '1': 'verifiedforsendingstatus',
      '3': 163100765,
      '4': 1,
      '5': 8,
      '10': 'verifiedforsendingstatus'
    },
  ],
  '3': [GetEmailIdentityResponse_PoliciesEntry$json],
};

@$core.Deprecated('Use getEmailIdentityResponseDescriptor instead')
const GetEmailIdentityResponse_PoliciesEntry$json = {
  '1': 'PoliciesEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GetEmailIdentityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailIdentityResponseDescriptor = $convert.base64Decode(
    'ChhHZXRFbWFpbElkZW50aXR5UmVzcG9uc2USNgoUY29uZmlndXJhdGlvbnNldG5hbWUYzYuxwA'
    'EgASgJUhRjb25maWd1cmF0aW9uc2V0bmFtZRJACg5ka2ltYXR0cmlidXRlcxjQtYt6IAEoCzIV'
    'LmVtYWlsLkRraW1BdHRyaWJ1dGVzUg5ka2ltYXR0cmlidXRlcxI+ChhmZWVkYmFja2Zvcndhcm'
    'RpbmdzdGF0dXMY/K+gkQEgASgIUhhmZWVkYmFja2ZvcndhcmRpbmdzdGF0dXMSOwoMaWRlbnRp'
    'dHl0eXBlGISnie4BIAEoDjITLmVtYWlsLklkZW50aXR5VHlwZVIMaWRlbnRpdHl0eXBlEk0KEm'
    '1haWxmcm9tYXR0cmlidXRlcxiAjL3bASABKAsyGS5lbWFpbC5NYWlsRnJvbUF0dHJpYnV0ZXNS'
    'Em1haWxmcm9tYXR0cmlidXRlcxJMCghwb2xpY2llcxiYrIoTIAMoCzItLmVtYWlsLkdldEVtYW'
    'lsSWRlbnRpdHlSZXNwb25zZS5Qb2xpY2llc0VudHJ5Ughwb2xpY2llcxIiCgR0YWdzGMHB9rUB'
    'IAMoCzIKLmVtYWlsLlRhZ1IEdGFncxJGChB2ZXJpZmljYXRpb25pbmZvGJ/Kzx8gASgLMhcuZW'
    '1haWwuVmVyaWZpY2F0aW9uSW5mb1IQdmVyaWZpY2F0aW9uaW5mbxJMChJ2ZXJpZmljYXRpb25z'
    'dGF0dXMYwZOkPyABKA4yGS5lbWFpbC5WZXJpZmljYXRpb25TdGF0dXNSEnZlcmlmaWNhdGlvbn'
    'N0YXR1cxI9Chh2ZXJpZmllZGZvcnNlbmRpbmdzdGF0dXMY3fDiTSABKAhSGHZlcmlmaWVkZm9y'
    'c2VuZGluZ3N0YXR1cxo7Cg1Qb2xpY2llc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbH'
    'VlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use getEmailTemplateRequestDescriptor instead')
const GetEmailTemplateRequest$json = {
  '1': 'GetEmailTemplateRequest',
  '2': [
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `GetEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailTemplateRequestDescriptor =
    $convert.base64Decode(
        'ChdHZXRFbWFpbFRlbXBsYXRlUmVxdWVzdBImCgx0ZW1wbGF0ZW5hbWUYsZiR5QEgASgJUgx0ZW'
        '1wbGF0ZW5hbWU=');

@$core.Deprecated('Use getEmailTemplateResponseDescriptor instead')
const GetEmailTemplateResponse$json = {
  '1': 'GetEmailTemplateResponse',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {
      '1': 'templatecontent',
      '3': 528973453,
      '4': 1,
      '5': 11,
      '6': '.email.EmailTemplateContent',
      '10': 'templatecontent'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `GetEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getEmailTemplateResponseDescriptor = $convert.base64Decode(
    'ChhHZXRFbWFpbFRlbXBsYXRlUmVzcG9uc2USIgoEdGFncxjBwfa1ASADKAsyCi5lbWFpbC5UYW'
    'dSBHRhZ3MSSQoPdGVtcGxhdGVjb250ZW50GI39nfwBIAEoCzIbLmVtYWlsLkVtYWlsVGVtcGxh'
    'dGVDb250ZW50Ug90ZW1wbGF0ZWNvbnRlbnQSJgoMdGVtcGxhdGVuYW1lGLGYkeUBIAEoCVIMdG'
    'VtcGxhdGVuYW1l');

@$core.Deprecated('Use getExportJobRequestDescriptor instead')
const GetExportJobRequest$json = {
  '1': 'GetExportJobRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
  ],
};

/// Descriptor for `GetExportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getExportJobRequestDescriptor =
    $convert.base64Decode(
        'ChNHZXRFeHBvcnRKb2JSZXF1ZXN0EhcKBWpvYmlkGNLU3TMgASgJUgVqb2JpZA==');

@$core.Deprecated('Use getExportJobResponseDescriptor instead')
const GetExportJobResponse$json = {
  '1': 'GetExportJobResponse',
  '2': [
    {
      '1': 'completedtimestamp',
      '3': 333312139,
      '4': 1,
      '5': 9,
      '10': 'completedtimestamp'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'exportdatasource',
      '3': 308051423,
      '4': 1,
      '5': 11,
      '6': '.email.ExportDataSource',
      '10': 'exportdatasource'
    },
    {
      '1': 'exportdestination',
      '3': 523408618,
      '4': 1,
      '5': 11,
      '6': '.email.ExportDestination',
      '10': 'exportdestination'
    },
    {
      '1': 'exportsourcetype',
      '3': 248243607,
      '4': 1,
      '5': 14,
      '6': '.email.ExportSourceType',
      '10': 'exportsourcetype'
    },
    {
      '1': 'failureinfo',
      '3': 451945802,
      '4': 1,
      '5': 11,
      '6': '.email.FailureInfo',
      '10': 'failureinfo'
    },
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.email.JobStatus',
      '10': 'jobstatus'
    },
    {
      '1': 'statistics',
      '3': 510636075,
      '4': 1,
      '5': 11,
      '6': '.email.ExportStatistics',
      '10': 'statistics'
    },
  ],
};

/// Descriptor for `GetExportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getExportJobResponseDescriptor = $convert.base64Decode(
    'ChRHZXRFeHBvcnRKb2JSZXNwb25zZRIyChJjb21wbGV0ZWR0aW1lc3RhbXAYi+H3ngEgASgJUh'
    'Jjb21wbGV0ZWR0aW1lc3RhbXASLgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEGNyZWF0'
    'ZWR0aW1lc3RhbXASRwoQZXhwb3J0ZGF0YXNvdXJjZRjf+/GSASABKAsyFy5lbWFpbC5FeHBvcn'
    'REYXRhU291cmNlUhBleHBvcnRkYXRhc291cmNlEkoKEWV4cG9ydGRlc3RpbmF0aW9uGOqpyvkB'
    'IAEoCzIYLmVtYWlsLkV4cG9ydERlc3RpbmF0aW9uUhFleHBvcnRkZXN0aW5hdGlvbhJGChBleH'
    'BvcnRzb3VyY2V0eXBlGJfLr3YgASgOMhcuZW1haWwuRXhwb3J0U291cmNlVHlwZVIQZXhwb3J0'
    'c291cmNldHlwZRI4CgtmYWlsdXJlaW5mbxjKysDXASABKAsyEi5lbWFpbC5GYWlsdXJlSW5mb1'
    'ILZmFpbHVyZWluZm8SFwoFam9iaWQY0tTdMyABKAlSBWpvYmlkEjEKCWpvYnN0YXR1cxjHnPsz'
    'IAEoDjIQLmVtYWlsLkpvYlN0YXR1c1IJam9ic3RhdHVzEjsKCnN0YXRpc3RpY3MYq+C+8wEgAS'
    'gLMhcuZW1haWwuRXhwb3J0U3RhdGlzdGljc1IKc3RhdGlzdGljcw==');

@$core.Deprecated('Use getImportJobRequestDescriptor instead')
const GetImportJobRequest$json = {
  '1': 'GetImportJobRequest',
  '2': [
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
  ],
};

/// Descriptor for `GetImportJobRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImportJobRequestDescriptor =
    $convert.base64Decode(
        'ChNHZXRJbXBvcnRKb2JSZXF1ZXN0EhcKBWpvYmlkGNLU3TMgASgJUgVqb2JpZA==');

@$core.Deprecated('Use getImportJobResponseDescriptor instead')
const GetImportJobResponse$json = {
  '1': 'GetImportJobResponse',
  '2': [
    {
      '1': 'completedtimestamp',
      '3': 333312139,
      '4': 1,
      '5': 9,
      '10': 'completedtimestamp'
    },
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'failedrecordscount',
      '3': 528801670,
      '4': 1,
      '5': 5,
      '10': 'failedrecordscount'
    },
    {
      '1': 'failureinfo',
      '3': 451945802,
      '4': 1,
      '5': 11,
      '6': '.email.FailureInfo',
      '10': 'failureinfo'
    },
    {
      '1': 'importdatasource',
      '3': 486006026,
      '4': 1,
      '5': 11,
      '6': '.email.ImportDataSource',
      '10': 'importdatasource'
    },
    {
      '1': 'importdestination',
      '3': 146287461,
      '4': 1,
      '5': 11,
      '6': '.email.ImportDestination',
      '10': 'importdestination'
    },
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.email.JobStatus',
      '10': 'jobstatus'
    },
    {
      '1': 'processedrecordscount',
      '3': 507944491,
      '4': 1,
      '5': 5,
      '10': 'processedrecordscount'
    },
  ],
};

/// Descriptor for `GetImportJobResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getImportJobResponseDescriptor = $convert.base64Decode(
    'ChRHZXRJbXBvcnRKb2JSZXNwb25zZRIyChJjb21wbGV0ZWR0aW1lc3RhbXAYi+H3ngEgASgJUh'
    'Jjb21wbGV0ZWR0aW1lc3RhbXASLgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEGNyZWF0'
    'ZWR0aW1lc3RhbXASMgoSZmFpbGVkcmVjb3Jkc2NvdW50GIa/k/wBIAEoBVISZmFpbGVkcmVjb3'
    'Jkc2NvdW50EjgKC2ZhaWx1cmVpbmZvGMrKwNcBIAEoCzISLmVtYWlsLkZhaWx1cmVJbmZvUgtm'
    'YWlsdXJlaW5mbxJHChBpbXBvcnRkYXRhc291cmNlGIq63+cBIAEoCzIXLmVtYWlsLkltcG9ydE'
    'RhdGFTb3VyY2VSEGltcG9ydGRhdGFzb3VyY2USSQoRaW1wb3J0ZGVzdGluYXRpb24Y5dbgRSAB'
    'KAsyGC5lbWFpbC5JbXBvcnREZXN0aW5hdGlvblIRaW1wb3J0ZGVzdGluYXRpb24SFwoFam9iaW'
    'QY0tTdMyABKAlSBWpvYmlkEjEKCWpvYnN0YXR1cxjHnPszIAEoDjIQLmVtYWlsLkpvYlN0YXR1'
    'c1IJam9ic3RhdHVzEjgKFXByb2Nlc3NlZHJlY29yZHNjb3VudBirvJryASABKAVSFXByb2Nlc3'
    'NlZHJlY29yZHNjb3VudA==');

@$core.Deprecated('Use getMessageInsightsRequestDescriptor instead')
const GetMessageInsightsRequest$json = {
  '1': 'GetMessageInsightsRequest',
  '2': [
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
  ],
};

/// Descriptor for `GetMessageInsightsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMessageInsightsRequestDescriptor =
    $convert.base64Decode(
        'ChlHZXRNZXNzYWdlSW5zaWdodHNSZXF1ZXN0EiAKCW1lc3NhZ2VpZBiq5vSrASABKAlSCW1lc3'
        'NhZ2VpZA==');

@$core.Deprecated('Use getMessageInsightsResponseDescriptor instead')
const GetMessageInsightsResponse$json = {
  '1': 'GetMessageInsightsResponse',
  '2': [
    {
      '1': 'emailtags',
      '3': 215809691,
      '4': 3,
      '5': 11,
      '6': '.email.MessageTag',
      '10': 'emailtags'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {
      '1': 'insights',
      '3': 3076585,
      '4': 3,
      '5': 11,
      '6': '.email.EmailInsights',
      '10': 'insights'
    },
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
    {'1': 'subject', '3': 7939312, '4': 1, '5': 9, '10': 'subject'},
  ],
};

/// Descriptor for `GetMessageInsightsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMessageInsightsResponseDescriptor = $convert.base64Decode(
    'ChpHZXRNZXNzYWdlSW5zaWdodHNSZXNwb25zZRIyCgllbWFpbHRhZ3MYm/3zZiADKAsyES5lbW'
    'FpbC5NZXNzYWdlVGFnUgllbWFpbHRhZ3MSLQoQZnJvbWVtYWlsYWRkcmVzcxiGmsssIAEoCVIQ'
    'ZnJvbWVtYWlsYWRkcmVzcxIzCghpbnNpZ2h0cxjp47sBIAMoCzIULmVtYWlsLkVtYWlsSW5zaW'
    'dodHNSCGluc2lnaHRzEiAKCW1lc3NhZ2VpZBiq5vSrASABKAlSCW1lc3NhZ2VpZBIbCgdzdWJq'
    'ZWN0GPDJ5AMgASgJUgdzdWJqZWN0');

@$core.Deprecated('Use getMultiRegionEndpointRequestDescriptor instead')
const GetMultiRegionEndpointRequest$json = {
  '1': 'GetMultiRegionEndpointRequest',
  '2': [
    {'1': 'endpointname', '3': 209534392, '4': 1, '5': 9, '10': 'endpointname'},
  ],
};

/// Descriptor for `GetMultiRegionEndpointRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMultiRegionEndpointRequestDescriptor =
    $convert.base64Decode(
        'Ch1HZXRNdWx0aVJlZ2lvbkVuZHBvaW50UmVxdWVzdBIlCgxlbmRwb2ludG5hbWUYuPv0YyABKA'
        'lSDGVuZHBvaW50bmFtZQ==');

@$core.Deprecated('Use getMultiRegionEndpointResponseDescriptor instead')
const GetMultiRegionEndpointResponse$json = {
  '1': 'GetMultiRegionEndpointResponse',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'endpointid', '3': 35808946, '4': 1, '5': 9, '10': 'endpointid'},
    {'1': 'endpointname', '3': 209534392, '4': 1, '5': 9, '10': 'endpointname'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {
      '1': 'routes',
      '3': 321835704,
      '4': 3,
      '5': 11,
      '6': '.email.Route',
      '10': 'routes'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.Status',
      '10': 'status'
    },
  ],
};

/// Descriptor for `GetMultiRegionEndpointResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getMultiRegionEndpointResponseDescriptor = $convert.base64Decode(
    'Ch5HZXRNdWx0aVJlZ2lvbkVuZHBvaW50UmVzcG9uc2USLgoQY3JlYXRlZHRpbWVzdGFtcBj628'
    '+fASABKAlSEGNyZWF0ZWR0aW1lc3RhbXASIQoKZW5kcG9pbnRpZBiyzYkRIAEoCVIKZW5kcG9p'
    'bnRpZBIlCgxlbmRwb2ludG5hbWUYuPv0YyABKAlSDGVuZHBvaW50bmFtZRI1ChRsYXN0dXBkYX'
    'RlZHRpbWVzdGFtcBiVy8g/IAEoCVIUbGFzdHVwZGF0ZWR0aW1lc3RhbXASKAoGcm91dGVzGLil'
    'u5kBIAMoCzIMLmVtYWlsLlJvdXRlUgZyb3V0ZXMSKAoGc3RhdHVzGJDk+wIgASgOMg0uZW1haW'
    'wuU3RhdHVzUgZzdGF0dXM=');

@$core.Deprecated('Use getReputationEntityRequestDescriptor instead')
const GetReputationEntityRequest$json = {
  '1': 'GetReputationEntityRequest',
  '2': [
    {
      '1': 'reputationentityreference',
      '3': 414929111,
      '4': 1,
      '5': 9,
      '10': 'reputationentityreference'
    },
    {
      '1': 'reputationentitytype',
      '3': 98287826,
      '4': 1,
      '5': 14,
      '6': '.email.ReputationEntityType',
      '10': 'reputationentitytype'
    },
  ],
};

/// Descriptor for `GetReputationEntityRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReputationEntityRequestDescriptor = $convert.base64Decode(
    'ChpHZXRSZXB1dGF0aW9uRW50aXR5UmVxdWVzdBJAChlyZXB1dGF0aW9uZW50aXR5cmVmZXJlbm'
    'NlGNeh7cUBIAEoCVIZcmVwdXRhdGlvbmVudGl0eXJlZmVyZW5jZRJSChRyZXB1dGF0aW9uZW50'
    'aXR5dHlwZRjSge8uIAEoDjIbLmVtYWlsLlJlcHV0YXRpb25FbnRpdHlUeXBlUhRyZXB1dGF0aW'
    '9uZW50aXR5dHlwZQ==');

@$core.Deprecated('Use getReputationEntityResponseDescriptor instead')
const GetReputationEntityResponse$json = {
  '1': 'GetReputationEntityResponse',
  '2': [
    {
      '1': 'reputationentity',
      '3': 126834426,
      '4': 1,
      '5': 11,
      '6': '.email.ReputationEntity',
      '10': 'reputationentity'
    },
  ],
};

/// Descriptor for `GetReputationEntityResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getReputationEntityResponseDescriptor =
    $convert.base64Decode(
        'ChtHZXRSZXB1dGF0aW9uRW50aXR5UmVzcG9uc2USRgoQcmVwdXRhdGlvbmVudGl0eRj6rb08IA'
        'EoCzIXLmVtYWlsLlJlcHV0YXRpb25FbnRpdHlSEHJlcHV0YXRpb25lbnRpdHk=');

@$core.Deprecated('Use getSuppressedDestinationRequestDescriptor instead')
const GetSuppressedDestinationRequest$json = {
  '1': 'GetSuppressedDestinationRequest',
  '2': [
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
  ],
};

/// Descriptor for `GetSuppressedDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSuppressedDestinationRequestDescriptor =
    $convert.base64Decode(
        'Ch9HZXRTdXBwcmVzc2VkRGVzdGluYXRpb25SZXF1ZXN0EiYKDGVtYWlsYWRkcmVzcxjW8YrpAS'
        'ABKAlSDGVtYWlsYWRkcmVzcw==');

@$core.Deprecated('Use getSuppressedDestinationResponseDescriptor instead')
const GetSuppressedDestinationResponse$json = {
  '1': 'GetSuppressedDestinationResponse',
  '2': [
    {
      '1': 'suppresseddestination',
      '3': 514454478,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressedDestination',
      '10': 'suppresseddestination'
    },
  ],
};

/// Descriptor for `GetSuppressedDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getSuppressedDestinationResponseDescriptor =
    $convert.base64Decode(
        'CiBHZXRTdXBwcmVzc2VkRGVzdGluYXRpb25SZXNwb25zZRJWChVzdXBwcmVzc2VkZGVzdGluYX'
        'Rpb24Yzuen9QEgASgLMhwuZW1haWwuU3VwcHJlc3NlZERlc3RpbmF0aW9uUhVzdXBwcmVzc2Vk'
        'ZGVzdGluYXRpb24=');

@$core.Deprecated('Use getTenantRequestDescriptor instead')
const GetTenantRequest$json = {
  '1': 'GetTenantRequest',
  '2': [
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `GetTenantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTenantRequestDescriptor = $convert.base64Decode(
    'ChBHZXRUZW5hbnRSZXF1ZXN0EiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbmFudG5hbWU=');

@$core.Deprecated('Use getTenantResponseDescriptor instead')
const GetTenantResponse$json = {
  '1': 'GetTenantResponse',
  '2': [
    {
      '1': 'tenant',
      '3': 40855750,
      '4': 1,
      '5': 11,
      '6': '.email.Tenant',
      '10': 'tenant'
    },
  ],
};

/// Descriptor for `GetTenantResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getTenantResponseDescriptor = $convert.base64Decode(
    'ChFHZXRUZW5hbnRSZXNwb25zZRIoCgZ0ZW5hbnQYxtG9EyABKAsyDS5lbWFpbC5UZW5hbnRSBn'
    'RlbmFudA==');

@$core.Deprecated('Use guardianAttributesDescriptor instead')
const GuardianAttributes$json = {
  '1': 'GuardianAttributes',
  '2': [
    {
      '1': 'optimizedshareddelivery',
      '3': 305138524,
      '4': 1,
      '5': 14,
      '6': '.email.FeatureStatus',
      '10': 'optimizedshareddelivery'
    },
  ],
};

/// Descriptor for `GuardianAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List guardianAttributesDescriptor = $convert.base64Decode(
    'ChJHdWFyZGlhbkF0dHJpYnV0ZXMSUgoXb3B0aW1pemVkc2hhcmVkZGVsaXZlcnkY3JbAkQEgAS'
    'gOMhQuZW1haWwuRmVhdHVyZVN0YXR1c1IXb3B0aW1pemVkc2hhcmVkZGVsaXZlcnk=');

@$core.Deprecated('Use guardianOptionsDescriptor instead')
const GuardianOptions$json = {
  '1': 'GuardianOptions',
  '2': [
    {
      '1': 'optimizedshareddelivery',
      '3': 305138524,
      '4': 1,
      '5': 14,
      '6': '.email.FeatureStatus',
      '10': 'optimizedshareddelivery'
    },
  ],
};

/// Descriptor for `GuardianOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List guardianOptionsDescriptor = $convert.base64Decode(
    'Cg9HdWFyZGlhbk9wdGlvbnMSUgoXb3B0aW1pemVkc2hhcmVkZGVsaXZlcnkY3JbAkQEgASgOMh'
    'QuZW1haWwuRmVhdHVyZVN0YXR1c1IXb3B0aW1pemVkc2hhcmVkZGVsaXZlcnk=');

@$core.Deprecated('Use identityInfoDescriptor instead')
const IdentityInfo$json = {
  '1': 'IdentityInfo',
  '2': [
    {'1': 'identityname', '3': 251627009, '4': 1, '5': 9, '10': 'identityname'},
    {
      '1': 'identitytype',
      '3': 499274628,
      '4': 1,
      '5': 14,
      '6': '.email.IdentityType',
      '10': 'identitytype'
    },
    {
      '1': 'sendingenabled',
      '3': 194846115,
      '4': 1,
      '5': 8,
      '10': 'sendingenabled'
    },
    {
      '1': 'verificationstatus',
      '3': 132712897,
      '4': 1,
      '5': 14,
      '6': '.email.VerificationStatus',
      '10': 'verificationstatus'
    },
  ],
};

/// Descriptor for `IdentityInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List identityInfoDescriptor = $convert.base64Decode(
    'CgxJZGVudGl0eUluZm8SJQoMaWRlbnRpdHluYW1lGIGM/ncgASgJUgxpZGVudGl0eW5hbWUSOw'
    'oMaWRlbnRpdHl0eXBlGISnie4BIAEoDjITLmVtYWlsLklkZW50aXR5VHlwZVIMaWRlbnRpdHl0'
    'eXBlEikKDnNlbmRpbmdlbmFibGVkGKO79FwgASgIUg5zZW5kaW5nZW5hYmxlZBJMChJ2ZXJpZm'
    'ljYXRpb25zdGF0dXMYwZOkPyABKA4yGS5lbWFpbC5WZXJpZmljYXRpb25TdGF0dXNSEnZlcmlm'
    'aWNhdGlvbnN0YXR1cw==');

@$core.Deprecated('Use importDataSourceDescriptor instead')
const ImportDataSource$json = {
  '1': 'ImportDataSource',
  '2': [
    {
      '1': 'dataformat',
      '3': 89652083,
      '4': 1,
      '5': 14,
      '6': '.email.DataFormat',
      '10': 'dataformat'
    },
    {'1': 's3url', '3': 407792021, '4': 1, '5': 9, '10': 's3url'},
  ],
};

/// Descriptor for `ImportDataSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importDataSourceDescriptor = $convert.base64Decode(
    'ChBJbXBvcnREYXRhU291cmNlEjQKCmRhdGFmb3JtYXQY8/bfKiABKA4yES5lbWFpbC5EYXRhRm'
    '9ybWF0UgpkYXRhZm9ybWF0EhgKBXMzdXJsGJXTucIBIAEoCVIFczN1cmw=');

@$core.Deprecated('Use importDestinationDescriptor instead')
const ImportDestination$json = {
  '1': 'ImportDestination',
  '2': [
    {
      '1': 'contactlistdestination',
      '3': 259529442,
      '4': 1,
      '5': 11,
      '6': '.email.ContactListDestination',
      '10': 'contactlistdestination'
    },
    {
      '1': 'suppressionlistdestination',
      '3': 142232257,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionListDestination',
      '10': 'suppressionlistdestination'
    },
  ],
};

/// Descriptor for `ImportDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importDestinationDescriptor = $convert.base64Decode(
    'ChFJbXBvcnREZXN0aW5hdGlvbhJYChZjb250YWN0bGlzdGRlc3RpbmF0aW9uGOK14HsgASgLMh'
    '0uZW1haWwuQ29udGFjdExpc3REZXN0aW5hdGlvblIWY29udGFjdGxpc3RkZXN0aW5hdGlvbhJk'
    'ChpzdXBwcmVzc2lvbmxpc3RkZXN0aW5hdGlvbhjBlelDIAEoCzIhLmVtYWlsLlN1cHByZXNzaW'
    '9uTGlzdERlc3RpbmF0aW9uUhpzdXBwcmVzc2lvbmxpc3RkZXN0aW5hdGlvbg==');

@$core.Deprecated('Use importJobSummaryDescriptor instead')
const ImportJobSummary$json = {
  '1': 'ImportJobSummary',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'failedrecordscount',
      '3': 528801670,
      '4': 1,
      '5': 5,
      '10': 'failedrecordscount'
    },
    {
      '1': 'importdestination',
      '3': 146287461,
      '4': 1,
      '5': 11,
      '6': '.email.ImportDestination',
      '10': 'importdestination'
    },
    {'1': 'jobid', '3': 108489298, '4': 1, '5': 9, '10': 'jobid'},
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.email.JobStatus',
      '10': 'jobstatus'
    },
    {
      '1': 'processedrecordscount',
      '3': 507944491,
      '4': 1,
      '5': 5,
      '10': 'processedrecordscount'
    },
  ],
};

/// Descriptor for `ImportJobSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importJobSummaryDescriptor = $convert.base64Decode(
    'ChBJbXBvcnRKb2JTdW1tYXJ5Ei4KEGNyZWF0ZWR0aW1lc3RhbXAY+tvPnwEgASgJUhBjcmVhdG'
    'VkdGltZXN0YW1wEjIKEmZhaWxlZHJlY29yZHNjb3VudBiGv5P8ASABKAVSEmZhaWxlZHJlY29y'
    'ZHNjb3VudBJJChFpbXBvcnRkZXN0aW5hdGlvbhjl1uBFIAEoCzIYLmVtYWlsLkltcG9ydERlc3'
    'RpbmF0aW9uUhFpbXBvcnRkZXN0aW5hdGlvbhIXCgVqb2JpZBjS1N0zIAEoCVIFam9iaWQSMQoJ'
    'am9ic3RhdHVzGMec+zMgASgOMhAuZW1haWwuSm9iU3RhdHVzUglqb2JzdGF0dXMSOAoVcHJvY2'
    'Vzc2VkcmVjb3Jkc2NvdW50GKu8mvIBIAEoBVIVcHJvY2Vzc2VkcmVjb3Jkc2NvdW50');

@$core.Deprecated('Use inboxPlacementTrackingOptionDescriptor instead')
const InboxPlacementTrackingOption$json = {
  '1': 'InboxPlacementTrackingOption',
  '2': [
    {'1': 'global', '3': 68321487, '4': 1, '5': 8, '10': 'global'},
    {'1': 'trackedisps', '3': 21912765, '4': 3, '5': 9, '10': 'trackedisps'},
  ],
};

/// Descriptor for `InboxPlacementTrackingOption`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List inboxPlacementTrackingOptionDescriptor =
    $convert.base64Decode(
        'ChxJbmJveFBsYWNlbWVudFRyYWNraW5nT3B0aW9uEhkKBmdsb2JhbBjPgcogIAEoCFIGZ2xvYm'
        'FsEiMKC3RyYWNrZWRpc3BzGL25uQogAygJUgt0cmFja2VkaXNwcw==');

@$core.Deprecated('Use insightsEventDescriptor instead')
const InsightsEvent$json = {
  '1': 'InsightsEvent',
  '2': [
    {
      '1': 'details',
      '3': 247611974,
      '4': 1,
      '5': 11,
      '6': '.email.EventDetails',
      '10': 'details'
    },
    {'1': 'timestamp', '3': 162390468, '4': 1, '5': 9, '10': 'timestamp'},
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.email.EventType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `InsightsEvent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List insightsEventDescriptor = $convert.base64Decode(
    'Cg1JbnNpZ2h0c0V2ZW50EjAKB2RldGFpbHMYxoSJdiABKAsyEy5lbWFpbC5FdmVudERldGFpbH'
    'NSB2RldGFpbHMSHwoJdGltZXN0YW1wGMTDt00gASgJUgl0aW1lc3RhbXASKAoEdHlwZRjuoNeK'
    'ASABKA4yEC5lbWFpbC5FdmVudFR5cGVSBHR5cGU=');

@$core.Deprecated('Use internalServiceErrorExceptionDescriptor instead')
const InternalServiceErrorException$json = {
  '1': 'InternalServiceErrorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServiceErrorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServiceErrorExceptionDescriptor =
    $convert.base64Decode(
        'Ch1JbnRlcm5hbFNlcnZpY2VFcnJvckV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use invalidNextTokenExceptionDescriptor instead')
const InvalidNextTokenException$json = {
  '1': 'InvalidNextTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidNextTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidNextTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkTmV4dFRva2VuRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use ispPlacementDescriptor instead')
const IspPlacement$json = {
  '1': 'IspPlacement',
  '2': [
    {'1': 'ispname', '3': 79681213, '4': 1, '5': 9, '10': 'ispname'},
    {
      '1': 'placementstatistics',
      '3': 449703852,
      '4': 1,
      '5': 11,
      '6': '.email.PlacementStatistics',
      '10': 'placementstatistics'
    },
  ],
};

/// Descriptor for `IspPlacement`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ispPlacementDescriptor = $convert.base64Decode(
    'CgxJc3BQbGFjZW1lbnQSGwoHaXNwbmFtZRi9rf8lIAEoCVIHaXNwbmFtZRJQChNwbGFjZW1lbn'
    'RzdGF0aXN0aWNzGKzft9YBIAEoCzIaLmVtYWlsLlBsYWNlbWVudFN0YXRpc3RpY3NSE3BsYWNl'
    'bWVudHN0YXRpc3RpY3M=');

@$core.Deprecated('Use kinesisFirehoseDestinationDescriptor instead')
const KinesisFirehoseDestination$json = {
  '1': 'KinesisFirehoseDestination',
  '2': [
    {
      '1': 'deliverystreamarn',
      '3': 446884461,
      '4': 1,
      '5': 9,
      '10': 'deliverystreamarn'
    },
    {'1': 'iamrolearn', '3': 228393658, '4': 1, '5': 9, '10': 'iamrolearn'},
  ],
};

/// Descriptor for `KinesisFirehoseDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kinesisFirehoseDestinationDescriptor =
    $convert.base64Decode(
        'ChpLaW5lc2lzRmlyZWhvc2VEZXN0aW5hdGlvbhIwChFkZWxpdmVyeXN0cmVhbWFybhjt1IvVAS'
        'ABKAlSEWRlbGl2ZXJ5c3RyZWFtYXJuEiEKCmlhbXJvbGVhcm4YuoX0bCABKAlSCmlhbXJvbGVh'
        'cm4=');

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

@$core.Deprecated('Use listConfigurationSetsRequestDescriptor instead')
const ListConfigurationSetsRequest$json = {
  '1': 'ListConfigurationSetsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListConfigurationSetsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConfigurationSetsRequestDescriptor =
    $convert.base64Decode(
        'ChxMaXN0Q29uZmlndXJhdGlvblNldHNSZXF1ZXN0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbm'
        'V4dHRva2VuEh4KCHBhZ2VzaXplGLiTgtEBIAEoBVIIcGFnZXNpemU=');

@$core.Deprecated('Use listConfigurationSetsResponseDescriptor instead')
const ListConfigurationSetsResponse$json = {
  '1': 'ListConfigurationSetsResponse',
  '2': [
    {
      '1': 'configurationsets',
      '3': 287812245,
      '4': 3,
      '5': 9,
      '10': 'configurationsets'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListConfigurationSetsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listConfigurationSetsResponseDescriptor =
    $convert.base64Decode(
        'Ch1MaXN0Q29uZmlndXJhdGlvblNldHNSZXNwb25zZRIwChFjb25maWd1cmF0aW9uc2V0cxiV1Z'
        '6JASADKAlSEWNvbmZpZ3VyYXRpb25zZXRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRv'
        'a2Vu');

@$core.Deprecated('Use listContactListsRequestDescriptor instead')
const ListContactListsRequest$json = {
  '1': 'ListContactListsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListContactListsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContactListsRequestDescriptor =
    $convert.base64Decode(
        'ChdMaXN0Q29udGFjdExpc3RzUmVxdWVzdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2'
        'tlbhIeCghwYWdlc2l6ZRi4k4LRASABKAVSCHBhZ2VzaXpl');

@$core.Deprecated('Use listContactListsResponseDescriptor instead')
const ListContactListsResponse$json = {
  '1': 'ListContactListsResponse',
  '2': [
    {
      '1': 'contactlists',
      '3': 367303753,
      '4': 3,
      '5': 11,
      '6': '.email.ContactList',
      '10': 'contactlists'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListContactListsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContactListsResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0Q29udGFjdExpc3RzUmVzcG9uc2USOgoMY29udGFjdGxpc3RzGMm4kq8BIAMoCzISLm'
    'VtYWlsLkNvbnRhY3RMaXN0Ugxjb250YWN0bGlzdHMSHwoJbmV4dHRva2VuGP6EumcgASgJUglu'
    'ZXh0dG9rZW4=');

@$core.Deprecated('Use listContactsFilterDescriptor instead')
const ListContactsFilter$json = {
  '1': 'ListContactsFilter',
  '2': [
    {
      '1': 'filteredstatus',
      '3': 310398937,
      '4': 1,
      '5': 14,
      '6': '.email.SubscriptionStatus',
      '10': 'filteredstatus'
    },
    {
      '1': 'topicfilter',
      '3': 211917895,
      '4': 1,
      '5': 11,
      '6': '.email.TopicFilter',
      '10': 'topicfilter'
    },
  ],
};

/// Descriptor for `ListContactsFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContactsFilterDescriptor = $convert.base64Decode(
    'ChJMaXN0Q29udGFjdHNGaWx0ZXISRQoOZmlsdGVyZWRzdGF0dXMY2Z+BlAEgASgOMhkuZW1haW'
    'wuU3Vic2NyaXB0aW9uU3RhdHVzUg5maWx0ZXJlZHN0YXR1cxI3Cgt0b3BpY2ZpbHRlchjHuIZl'
    'IAEoCzISLmVtYWlsLlRvcGljRmlsdGVyUgt0b3BpY2ZpbHRlcg==');

@$core.Deprecated('Use listContactsRequestDescriptor instead')
const ListContactsRequest$json = {
  '1': 'ListContactsRequest',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {
      '1': 'filter',
      '3': 346669208,
      '4': 1,
      '5': 11,
      '6': '.email.ListContactsFilter',
      '10': 'filter'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListContactsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContactsRequestDescriptor = $convert.base64Decode(
    'ChNMaXN0Q29udGFjdHNSZXF1ZXN0EiwKD2NvbnRhY3RsaXN0bmFtZRjxvaehASABKAlSD2Nvbn'
    'RhY3RsaXN0bmFtZRI1CgZmaWx0ZXIYmIGnpQEgASgLMhkuZW1haWwuTGlzdENvbnRhY3RzRmls'
    'dGVyUgZmaWx0ZXISHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SHgoIcGFnZXNpem'
    'UYuJOC0QEgASgFUghwYWdlc2l6ZQ==');

@$core.Deprecated('Use listContactsResponseDescriptor instead')
const ListContactsResponse$json = {
  '1': 'ListContactsResponse',
  '2': [
    {
      '1': 'contacts',
      '3': 145625361,
      '4': 3,
      '5': 11,
      '6': '.email.Contact',
      '10': 'contacts'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListContactsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listContactsResponseDescriptor = $convert.base64Decode(
    'ChRMaXN0Q29udGFjdHNSZXNwb25zZRItCghjb250YWN0cxiRorhFIAMoCzIOLmVtYWlsLkNvbn'
    'RhY3RSCGNvbnRhY3RzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated(
    'Use listCustomVerificationEmailTemplatesRequestDescriptor instead')
const ListCustomVerificationEmailTemplatesRequest$json = {
  '1': 'ListCustomVerificationEmailTemplatesRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListCustomVerificationEmailTemplatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listCustomVerificationEmailTemplatesRequestDescriptor =
    $convert.base64Decode(
        'CitMaXN0Q3VzdG9tVmVyaWZpY2F0aW9uRW1haWxUZW1wbGF0ZXNSZXF1ZXN0Eh8KCW5leHR0b2'
        'tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEh4KCHBhZ2VzaXplGLiTgtEBIAEoBVIIcGFnZXNpemU=');

@$core.Deprecated(
    'Use listCustomVerificationEmailTemplatesResponseDescriptor instead')
const ListCustomVerificationEmailTemplatesResponse$json = {
  '1': 'ListCustomVerificationEmailTemplatesResponse',
  '2': [
    {
      '1': 'customverificationemailtemplates',
      '3': 449684503,
      '4': 3,
      '5': 11,
      '6': '.email.CustomVerificationEmailTemplateMetadata',
      '10': 'customverificationemailtemplates'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListCustomVerificationEmailTemplatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listCustomVerificationEmailTemplatesResponseDescriptor =
    $convert.base64Decode(
        'CixMaXN0Q3VzdG9tVmVyaWZpY2F0aW9uRW1haWxUZW1wbGF0ZXNSZXNwb25zZRJ+CiBjdXN0b2'
        '12ZXJpZmljYXRpb25lbWFpbHRlbXBsYXRlcxiXyLbWASADKAsyLi5lbWFpbC5DdXN0b21WZXJp'
        'ZmljYXRpb25FbWFpbFRlbXBsYXRlTWV0YWRhdGFSIGN1c3RvbXZlcmlmaWNhdGlvbmVtYWlsdG'
        'VtcGxhdGVzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listDedicatedIpPoolsRequestDescriptor instead')
const ListDedicatedIpPoolsRequest$json = {
  '1': 'ListDedicatedIpPoolsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListDedicatedIpPoolsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDedicatedIpPoolsRequestDescriptor =
    $convert.base64Decode(
        'ChtMaXN0RGVkaWNhdGVkSXBQb29sc1JlcXVlc3QSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SHgoIcGFnZXNpemUYuJOC0QEgASgFUghwYWdlc2l6ZQ==');

@$core.Deprecated('Use listDedicatedIpPoolsResponseDescriptor instead')
const ListDedicatedIpPoolsResponse$json = {
  '1': 'ListDedicatedIpPoolsResponse',
  '2': [
    {
      '1': 'dedicatedippools',
      '3': 224971511,
      '4': 3,
      '5': 9,
      '10': 'dedicatedippools'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDedicatedIpPoolsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDedicatedIpPoolsResponseDescriptor =
    $convert.base64Decode(
        'ChxMaXN0RGVkaWNhdGVkSXBQb29sc1Jlc3BvbnNlEi0KEGRlZGljYXRlZGlwcG9vbHMY95Wjay'
        'ADKAlSEGRlZGljYXRlZGlwcG9vbHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listDeliverabilityTestReportsRequestDescriptor instead')
const ListDeliverabilityTestReportsRequest$json = {
  '1': 'ListDeliverabilityTestReportsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListDeliverabilityTestReportsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDeliverabilityTestReportsRequestDescriptor =
    $convert.base64Decode(
        'CiRMaXN0RGVsaXZlcmFiaWxpdHlUZXN0UmVwb3J0c1JlcXVlc3QSHwoJbmV4dHRva2VuGP6Eum'
        'cgASgJUgluZXh0dG9rZW4SHgoIcGFnZXNpemUYuJOC0QEgASgFUghwYWdlc2l6ZQ==');

@$core.Deprecated('Use listDeliverabilityTestReportsResponseDescriptor instead')
const ListDeliverabilityTestReportsResponse$json = {
  '1': 'ListDeliverabilityTestReportsResponse',
  '2': [
    {
      '1': 'deliverabilitytestreports',
      '3': 226879076,
      '4': 3,
      '5': 11,
      '6': '.email.DeliverabilityTestReport',
      '10': 'deliverabilitytestreports'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDeliverabilityTestReportsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDeliverabilityTestReportsResponseDescriptor =
    $convert.base64Decode(
        'CiVMaXN0RGVsaXZlcmFiaWxpdHlUZXN0UmVwb3J0c1Jlc3BvbnNlEmAKGWRlbGl2ZXJhYmlsaX'
        'R5dGVzdHJlcG9ydHMY5MyXbCADKAsyHy5lbWFpbC5EZWxpdmVyYWJpbGl0eVRlc3RSZXBvcnRS'
        'GWRlbGl2ZXJhYmlsaXR5dGVzdHJlcG9ydHMSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
        '9rZW4=');

@$core.Deprecated(
    'Use listDomainDeliverabilityCampaignsRequestDescriptor instead')
const ListDomainDeliverabilityCampaignsRequest$json = {
  '1': 'ListDomainDeliverabilityCampaignsRequest',
  '2': [
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {
      '1': 'subscribeddomain',
      '3': 484104338,
      '4': 1,
      '5': 9,
      '10': 'subscribeddomain'
    },
  ],
};

/// Descriptor for `ListDomainDeliverabilityCampaignsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listDomainDeliverabilityCampaignsRequestDescriptor =
    $convert.base64Decode(
        'CihMaXN0RG9tYWluRGVsaXZlcmFiaWxpdHlDYW1wYWlnbnNSZXF1ZXN0EhsKB2VuZGRhdGUYz7'
        'P5JCABKAlSB2VuZGRhdGUSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SHgoIcGFn'
        'ZXNpemUYuJOC0QEgASgFUghwYWdlc2l6ZRIgCglzdGFydGRhdGUY/Pig1AEgASgJUglzdGFydG'
        'RhdGUSLgoQc3Vic2NyaWJlZGRvbWFpbhiSsevmASABKAlSEHN1YnNjcmliZWRkb21haW4=');

@$core.Deprecated(
    'Use listDomainDeliverabilityCampaignsResponseDescriptor instead')
const ListDomainDeliverabilityCampaignsResponse$json = {
  '1': 'ListDomainDeliverabilityCampaignsResponse',
  '2': [
    {
      '1': 'domaindeliverabilitycampaigns',
      '3': 99285708,
      '4': 3,
      '5': 11,
      '6': '.email.DomainDeliverabilityCampaign',
      '10': 'domaindeliverabilitycampaigns'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListDomainDeliverabilityCampaignsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    listDomainDeliverabilityCampaignsResponseDescriptor = $convert.base64Decode(
        'CilMaXN0RG9tYWluRGVsaXZlcmFiaWxpdHlDYW1wYWlnbnNSZXNwb25zZRJsCh1kb21haW5kZW'
        'xpdmVyYWJpbGl0eWNhbXBhaWducxjM9asvIAMoCzIjLmVtYWlsLkRvbWFpbkRlbGl2ZXJhYmls'
        'aXR5Q2FtcGFpZ25SHWRvbWFpbmRlbGl2ZXJhYmlsaXR5Y2FtcGFpZ25zEh8KCW5leHR0b2tlbh'
        'j+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listEmailIdentitiesRequestDescriptor instead')
const ListEmailIdentitiesRequest$json = {
  '1': 'ListEmailIdentitiesRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListEmailIdentitiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEmailIdentitiesRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0RW1haWxJZGVudGl0aWVzUmVxdWVzdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leH'
        'R0b2tlbhIeCghwYWdlc2l6ZRi4k4LRASABKAVSCHBhZ2VzaXpl');

@$core.Deprecated('Use listEmailIdentitiesResponseDescriptor instead')
const ListEmailIdentitiesResponse$json = {
  '1': 'ListEmailIdentitiesResponse',
  '2': [
    {
      '1': 'emailidentities',
      '3': 182046870,
      '4': 3,
      '5': 11,
      '6': '.email.IdentityInfo',
      '10': 'emailidentities'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListEmailIdentitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEmailIdentitiesResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0RW1haWxJZGVudGl0aWVzUmVzcG9uc2USQAoPZW1haWxpZGVudGl0aWVzGJah51YgAy'
        'gLMhMuZW1haWwuSWRlbnRpdHlJbmZvUg9lbWFpbGlkZW50aXRpZXMSHwoJbmV4dHRva2VuGP6E'
        'umcgASgJUgluZXh0dG9rZW4=');

@$core.Deprecated('Use listEmailTemplatesRequestDescriptor instead')
const ListEmailTemplatesRequest$json = {
  '1': 'ListEmailTemplatesRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListEmailTemplatesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEmailTemplatesRequestDescriptor =
    $convert.base64Decode(
        'ChlMaXN0RW1haWxUZW1wbGF0ZXNSZXF1ZXN0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dH'
        'Rva2VuEh4KCHBhZ2VzaXplGLiTgtEBIAEoBVIIcGFnZXNpemU=');

@$core.Deprecated('Use listEmailTemplatesResponseDescriptor instead')
const ListEmailTemplatesResponse$json = {
  '1': 'ListEmailTemplatesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'templatesmetadata',
      '3': 29103678,
      '4': 3,
      '5': 11,
      '6': '.email.EmailTemplateMetadata',
      '10': 'templatesmetadata'
    },
  ],
};

/// Descriptor for `ListEmailTemplatesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listEmailTemplatesResponseDescriptor =
    $convert.base64Decode(
        'ChpMaXN0RW1haWxUZW1wbGF0ZXNSZXNwb25zZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leH'
        'R0b2tlbhJNChF0ZW1wbGF0ZXNtZXRhZGF0YRi+rPANIAMoCzIcLmVtYWlsLkVtYWlsVGVtcGxh'
        'dGVNZXRhZGF0YVIRdGVtcGxhdGVzbWV0YWRhdGE=');

@$core.Deprecated('Use listExportJobsRequestDescriptor instead')
const ListExportJobsRequest$json = {
  '1': 'ListExportJobsRequest',
  '2': [
    {
      '1': 'exportsourcetype',
      '3': 248243607,
      '4': 1,
      '5': 14,
      '6': '.email.ExportSourceType',
      '10': 'exportsourcetype'
    },
    {
      '1': 'jobstatus',
      '3': 108973639,
      '4': 1,
      '5': 14,
      '6': '.email.JobStatus',
      '10': 'jobstatus'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListExportJobsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExportJobsRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0RXhwb3J0Sm9ic1JlcXVlc3QSRgoQZXhwb3J0c291cmNldHlwZRiXy692IAEoDjIXLm'
    'VtYWlsLkV4cG9ydFNvdXJjZVR5cGVSEGV4cG9ydHNvdXJjZXR5cGUSMQoJam9ic3RhdHVzGMec'
    '+zMgASgOMhAuZW1haWwuSm9iU3RhdHVzUglqb2JzdGF0dXMSHwoJbmV4dHRva2VuGP6EumcgAS'
    'gJUgluZXh0dG9rZW4SHgoIcGFnZXNpemUYuJOC0QEgASgFUghwYWdlc2l6ZQ==');

@$core.Deprecated('Use listExportJobsResponseDescriptor instead')
const ListExportJobsResponse$json = {
  '1': 'ListExportJobsResponse',
  '2': [
    {
      '1': 'exportjobs',
      '3': 428957132,
      '4': 3,
      '5': 11,
      '6': '.email.ExportJobSummary',
      '10': 'exportjobs'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListExportJobsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listExportJobsResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0RXhwb3J0Sm9ic1Jlc3BvbnNlEjsKCmV4cG9ydGpvYnMYzLvFzAEgAygLMhcuZW1haW'
    'wuRXhwb3J0Sm9iU3VtbWFyeVIKZXhwb3J0am9icxIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5l'
    'eHR0b2tlbg==');

@$core.Deprecated('Use listImportJobsRequestDescriptor instead')
const ListImportJobsRequest$json = {
  '1': 'ListImportJobsRequest',
  '2': [
    {
      '1': 'importdestinationtype',
      '3': 338152013,
      '4': 1,
      '5': 14,
      '6': '.email.ImportDestinationType',
      '10': 'importdestinationtype'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListImportJobsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportJobsRequestDescriptor = $convert.base64Decode(
    'ChVMaXN0SW1wb3J0Sm9ic1JlcXVlc3QSVgoVaW1wb3J0ZGVzdGluYXRpb250eXBlGM2Un6EBIA'
    'EoDjIcLmVtYWlsLkltcG9ydERlc3RpbmF0aW9uVHlwZVIVaW1wb3J0ZGVzdGluYXRpb250eXBl'
    'Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEh4KCHBhZ2VzaXplGLiTgtEBIAEoBV'
    'IIcGFnZXNpemU=');

@$core.Deprecated('Use listImportJobsResponseDescriptor instead')
const ListImportJobsResponse$json = {
  '1': 'ListImportJobsResponse',
  '2': [
    {
      '1': 'importjobs',
      '3': 95920397,
      '4': 3,
      '5': 11,
      '6': '.email.ImportJobSummary',
      '10': 'importjobs'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListImportJobsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listImportJobsResponseDescriptor = $convert.base64Decode(
    'ChZMaXN0SW1wb3J0Sm9ic1Jlc3BvbnNlEjoKCmltcG9ydGpvYnMYjcLeLSADKAsyFy5lbWFpbC'
    '5JbXBvcnRKb2JTdW1tYXJ5UgppbXBvcnRqb2JzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4'
    'dHRva2Vu');

@$core.Deprecated('Use listManagementOptionsDescriptor instead')
const ListManagementOptions$json = {
  '1': 'ListManagementOptions',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {'1': 'topicname', '3': 446010240, '4': 1, '5': 9, '10': 'topicname'},
  ],
};

/// Descriptor for `ListManagementOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listManagementOptionsDescriptor = $convert.base64Decode(
    'ChVMaXN0TWFuYWdlbWVudE9wdGlvbnMSLAoPY29udGFjdGxpc3RuYW1lGPG9p6EBIAEoCVIPY2'
    '9udGFjdGxpc3RuYW1lEiAKCXRvcGljbmFtZRiAp9bUASABKAlSCXRvcGljbmFtZQ==');

@$core.Deprecated('Use listMultiRegionEndpointsRequestDescriptor instead')
const ListMultiRegionEndpointsRequest$json = {
  '1': 'ListMultiRegionEndpointsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListMultiRegionEndpointsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMultiRegionEndpointsRequestDescriptor =
    $convert.base64Decode(
        'Ch9MaXN0TXVsdGlSZWdpb25FbmRwb2ludHNSZXF1ZXN0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCV'
        'IJbmV4dHRva2VuEh4KCHBhZ2VzaXplGLiTgtEBIAEoBVIIcGFnZXNpemU=');

@$core.Deprecated('Use listMultiRegionEndpointsResponseDescriptor instead')
const ListMultiRegionEndpointsResponse$json = {
  '1': 'ListMultiRegionEndpointsResponse',
  '2': [
    {
      '1': 'multiregionendpoints',
      '3': 190926753,
      '4': 3,
      '5': 11,
      '6': '.email.MultiRegionEndpoint',
      '10': 'multiregionendpoints'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListMultiRegionEndpointsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listMultiRegionEndpointsResponseDescriptor =
    $convert.base64Decode(
        'CiBMaXN0TXVsdGlSZWdpb25FbmRwb2ludHNSZXNwb25zZRJRChRtdWx0aXJlZ2lvbmVuZHBvaW'
        '50cxihn4VbIAMoCzIaLmVtYWlsLk11bHRpUmVnaW9uRW5kcG9pbnRSFG11bHRpcmVnaW9uZW5k'
        'cG9pbnRzEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2Vu');

@$core.Deprecated('Use listRecommendationsRequestDescriptor instead')
const ListRecommendationsRequest$json = {
  '1': 'ListRecommendationsRequest',
  '2': [
    {
      '1': 'filter',
      '3': 346669208,
      '4': 3,
      '5': 11,
      '6': '.email.ListRecommendationsRequest.FilterEntry',
      '10': 'filter'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
  '3': [ListRecommendationsRequest_FilterEntry$json],
};

@$core.Deprecated('Use listRecommendationsRequestDescriptor instead')
const ListRecommendationsRequest_FilterEntry$json = {
  '1': 'FilterEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListRecommendationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRecommendationsRequestDescriptor = $convert.base64Decode(
    'ChpMaXN0UmVjb21tZW5kYXRpb25zUmVxdWVzdBJJCgZmaWx0ZXIYmIGnpQEgAygLMi0uZW1haW'
    'wuTGlzdFJlY29tbWVuZGF0aW9uc1JlcXVlc3QuRmlsdGVyRW50cnlSBmZpbHRlchIfCgluZXh0'
    'dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIeCghwYWdlc2l6ZRi4k4LRASABKAVSCHBhZ2VzaX'
    'plGjkKC0ZpbHRlckVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1'
    'ZToCOAE=');

@$core.Deprecated('Use listRecommendationsResponseDescriptor instead')
const ListRecommendationsResponse$json = {
  '1': 'ListRecommendationsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'recommendations',
      '3': 513799432,
      '4': 3,
      '5': 11,
      '6': '.email.Recommendation',
      '10': 'recommendations'
    },
  ],
};

/// Descriptor for `ListRecommendationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRecommendationsResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0UmVjb21tZW5kYXRpb25zUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SQwoPcmVjb21tZW5kYXRpb25zGIjq//QBIAMoCzIVLmVtYWlsLlJlY29tbWVuZGF0'
        'aW9uUg9yZWNvbW1lbmRhdGlvbnM=');

@$core.Deprecated('Use listReputationEntitiesRequestDescriptor instead')
const ListReputationEntitiesRequest$json = {
  '1': 'ListReputationEntitiesRequest',
  '2': [
    {
      '1': 'filter',
      '3': 346669208,
      '4': 3,
      '5': 11,
      '6': '.email.ListReputationEntitiesRequest.FilterEntry',
      '10': 'filter'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
  '3': [ListReputationEntitiesRequest_FilterEntry$json],
};

@$core.Deprecated('Use listReputationEntitiesRequestDescriptor instead')
const ListReputationEntitiesRequest_FilterEntry$json = {
  '1': 'FilterEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListReputationEntitiesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReputationEntitiesRequestDescriptor = $convert.base64Decode(
    'Ch1MaXN0UmVwdXRhdGlvbkVudGl0aWVzUmVxdWVzdBJMCgZmaWx0ZXIYmIGnpQEgAygLMjAuZW'
    '1haWwuTGlzdFJlcHV0YXRpb25FbnRpdGllc1JlcXVlc3QuRmlsdGVyRW50cnlSBmZpbHRlchIf'
    'CgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIeCghwYWdlc2l6ZRi4k4LRASABKAVSCH'
    'BhZ2VzaXplGjkKC0ZpbHRlckVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJ'
    'UgV2YWx1ZToCOAE=');

@$core.Deprecated('Use listReputationEntitiesResponseDescriptor instead')
const ListReputationEntitiesResponse$json = {
  '1': 'ListReputationEntitiesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'reputationentities',
      '3': 15832950,
      '4': 3,
      '5': 11,
      '6': '.email.ReputationEntity',
      '10': 'reputationentities'
    },
  ],
};

/// Descriptor for `ListReputationEntitiesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listReputationEntitiesResponseDescriptor =
    $convert.base64Decode(
        'Ch5MaXN0UmVwdXRhdGlvbkVudGl0aWVzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUg'
        'luZXh0dG9rZW4SSgoScmVwdXRhdGlvbmVudGl0aWVzGPauxgcgAygLMhcuZW1haWwuUmVwdXRh'
        'dGlvbkVudGl0eVIScmVwdXRhdGlvbmVudGl0aWVz');

@$core.Deprecated('Use listResourceTenantsRequestDescriptor instead')
const ListResourceTenantsRequest$json = {
  '1': 'ListResourceTenantsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
  ],
};

/// Descriptor for `ListResourceTenantsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceTenantsRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0UmVzb3VyY2VUZW5hbnRzUmVxdWVzdBIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leH'
        'R0b2tlbhIeCghwYWdlc2l6ZRi4k4LRASABKAVSCHBhZ2VzaXplEiQKC3Jlc291cmNlYXJuGK34'
        '2a0BIAEoCVILcmVzb3VyY2Vhcm4=');

@$core.Deprecated('Use listResourceTenantsResponseDescriptor instead')
const ListResourceTenantsResponse$json = {
  '1': 'ListResourceTenantsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'resourcetenants',
      '3': 465346835,
      '4': 3,
      '5': 11,
      '6': '.email.ResourceTenantMetadata',
      '10': 'resourcetenants'
    },
  ],
};

/// Descriptor for `ListResourceTenantsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceTenantsResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0UmVzb3VyY2VUZW5hbnRzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SSwoPcmVzb3VyY2V0ZW5hbnRzGJPC8t0BIAMoCzIdLmVtYWlsLlJlc291cmNlVGVu'
        'YW50TWV0YWRhdGFSD3Jlc291cmNldGVuYW50cw==');

@$core.Deprecated('Use listSuppressedDestinationsRequestDescriptor instead')
const ListSuppressedDestinationsRequest$json = {
  '1': 'ListSuppressedDestinationsRequest',
  '2': [
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
    {
      '1': 'reasons',
      '3': 176801663,
      '4': 3,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'reasons'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
  ],
};

/// Descriptor for `ListSuppressedDestinationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSuppressedDestinationsRequestDescriptor = $convert.base64Decode(
    'CiFMaXN0U3VwcHJlc3NlZERlc3RpbmF0aW9uc1JlcXVlc3QSGwoHZW5kZGF0ZRjPs/kkIAEoCV'
    'IHZW5kZGF0ZRIfCgluZXh0dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIeCghwYWdlc2l6ZRi4'
    'k4LRASABKAVSCHBhZ2VzaXplEjkKB3JlYXNvbnMY/46nVCADKA4yHC5lbWFpbC5TdXBwcmVzc2'
    'lvbkxpc3RSZWFzb25SB3JlYXNvbnMSIAoJc3RhcnRkYXRlGPz4oNQBIAEoCVIJc3RhcnRkYXRl');

@$core.Deprecated('Use listSuppressedDestinationsResponseDescriptor instead')
const ListSuppressedDestinationsResponse$json = {
  '1': 'ListSuppressedDestinationsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'suppresseddestinationsummaries',
      '3': 399108924,
      '4': 3,
      '5': 11,
      '6': '.email.SuppressedDestinationSummary',
      '10': 'suppresseddestinationsummaries'
    },
  ],
};

/// Descriptor for `ListSuppressedDestinationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSuppressedDestinationsResponseDescriptor =
    $convert.base64Decode(
        'CiJMaXN0U3VwcHJlc3NlZERlc3RpbmF0aW9uc1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIA'
        'EoCVIJbmV4dHRva2VuEm8KHnN1cHByZXNzZWRkZXN0aW5hdGlvbnN1bW1hcmllcxi81qe+ASAD'
        'KAsyIy5lbWFpbC5TdXBwcmVzc2VkRGVzdGluYXRpb25TdW1tYXJ5Uh5zdXBwcmVzc2VkZGVzdG'
        'luYXRpb25zdW1tYXJpZXM=');

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
      '6': '.email.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGFnc0ZvclJlc291cmNlUmVzcG9uc2USIgoEdGFncxjBwfa1ASADKAsyCi5lbWFpbC'
        '5UYWdSBHRhZ3M=');

@$core.Deprecated('Use listTenantResourcesRequestDescriptor instead')
const ListTenantResourcesRequest$json = {
  '1': 'ListTenantResourcesRequest',
  '2': [
    {
      '1': 'filter',
      '3': 346669208,
      '4': 3,
      '5': 11,
      '6': '.email.ListTenantResourcesRequest.FilterEntry',
      '10': 'filter'
    },
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
  '3': [ListTenantResourcesRequest_FilterEntry$json],
};

@$core.Deprecated('Use listTenantResourcesRequestDescriptor instead')
const ListTenantResourcesRequest_FilterEntry$json = {
  '1': 'FilterEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ListTenantResourcesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTenantResourcesRequestDescriptor = $convert.base64Decode(
    'ChpMaXN0VGVuYW50UmVzb3VyY2VzUmVxdWVzdBJJCgZmaWx0ZXIYmIGnpQEgAygLMi0uZW1haW'
    'wuTGlzdFRlbmFudFJlc291cmNlc1JlcXVlc3QuRmlsdGVyRW50cnlSBmZpbHRlchIfCgluZXh0'
    'dG9rZW4Y/oS6ZyABKAlSCW5leHR0b2tlbhIeCghwYWdlc2l6ZRi4k4LRASABKAVSCHBhZ2VzaX'
    'plEiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbmFudG5hbWUaOQoLRmlsdGVyRW50cnkSEAoD'
    'a2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use listTenantResourcesResponseDescriptor instead')
const ListTenantResourcesResponse$json = {
  '1': 'ListTenantResourcesResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tenantresources',
      '3': 121330067,
      '4': 3,
      '5': 11,
      '6': '.email.TenantResource',
      '10': 'tenantresources'
    },
  ],
};

/// Descriptor for `ListTenantResourcesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTenantResourcesResponseDescriptor =
    $convert.base64Decode(
        'ChtMaXN0VGVuYW50UmVzb3VyY2VzUmVzcG9uc2USHwoJbmV4dHRva2VuGP6EumcgASgJUgluZX'
        'h0dG9rZW4SQgoPdGVuYW50cmVzb3VyY2VzGJOz7TkgAygLMhUuZW1haWwuVGVuYW50UmVzb3Vy'
        'Y2VSD3RlbmFudHJlc291cmNlcw==');

@$core.Deprecated('Use listTenantsRequestDescriptor instead')
const ListTenantsRequest$json = {
  '1': 'ListTenantsRequest',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'pagesize', '3': 438340024, '4': 1, '5': 5, '10': 'pagesize'},
  ],
};

/// Descriptor for `ListTenantsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTenantsRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0VGVuYW50c1JlcXVlc3QSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG9rZW4SHg'
    'oIcGFnZXNpemUYuJOC0QEgASgFUghwYWdlc2l6ZQ==');

@$core.Deprecated('Use listTenantsResponseDescriptor instead')
const ListTenantsResponse$json = {
  '1': 'ListTenantsResponse',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'tenants',
      '3': 190961283,
      '4': 3,
      '5': 11,
      '6': '.email.TenantInfo',
      '10': 'tenants'
    },
  ],
};

/// Descriptor for `ListTenantsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTenantsResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0VGVuYW50c1Jlc3BvbnNlEh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEi'
    '4KB3RlbmFudHMYg62HWyADKAsyES5lbWFpbC5UZW5hbnRJbmZvUgd0ZW5hbnRz');

@$core.Deprecated('Use mailFromAttributesDescriptor instead')
const MailFromAttributes$json = {
  '1': 'MailFromAttributes',
  '2': [
    {
      '1': 'behavioronmxfailure',
      '3': 494873128,
      '4': 1,
      '5': 14,
      '6': '.email.BehaviorOnMxFailure',
      '10': 'behavioronmxfailure'
    },
    {
      '1': 'mailfromdomain',
      '3': 512250671,
      '4': 1,
      '5': 9,
      '10': 'mailfromdomain'
    },
    {
      '1': 'mailfromdomainstatus',
      '3': 364945809,
      '4': 1,
      '5': 14,
      '6': '.email.MailFromDomainStatus',
      '10': 'mailfromdomainstatus'
    },
  ],
};

/// Descriptor for `MailFromAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mailFromAttributesDescriptor = $convert.base64Decode(
    'ChJNYWlsRnJvbUF0dHJpYnV0ZXMSUAoTYmVoYXZpb3Jvbm14ZmFpbHVyZRio1PzrASABKA4yGi'
    '5lbWFpbC5CZWhhdmlvck9uTXhGYWlsdXJlUhNiZWhhdmlvcm9ubXhmYWlsdXJlEioKDm1haWxm'
    'cm9tZG9tYWluGK+mofQBIAEoCVIObWFpbGZyb21kb21haW4SUwoUbWFpbGZyb21kb21haW5zdG'
    'F0dXMYkcOCrgEgASgOMhsuZW1haWwuTWFpbEZyb21Eb21haW5TdGF0dXNSFG1haWxmcm9tZG9t'
    'YWluc3RhdHVz');

@$core.Deprecated('Use mailFromDomainNotVerifiedExceptionDescriptor instead')
const MailFromDomainNotVerifiedException$json = {
  '1': 'MailFromDomainNotVerifiedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MailFromDomainNotVerifiedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mailFromDomainNotVerifiedExceptionDescriptor =
    $convert.base64Decode(
        'CiJNYWlsRnJvbURvbWFpbk5vdFZlcmlmaWVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

@$core.Deprecated('Use mailboxValidationDescriptor instead')
const MailboxValidation$json = {
  '1': 'MailboxValidation',
  '2': [
    {
      '1': 'evaluations',
      '3': 207212261,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsMailboxEvaluations',
      '10': 'evaluations'
    },
    {
      '1': 'isvalid',
      '3': 163545986,
      '4': 1,
      '5': 11,
      '6': '.email.EmailAddressInsightsVerdict',
      '10': 'isvalid'
    },
  ],
};

/// Descriptor for `MailboxValidation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List mailboxValidationDescriptor = $convert.base64Decode(
    'ChFNYWlsYm94VmFsaWRhdGlvbhJSCgtldmFsdWF0aW9ucxjlnediIAEoCzItLmVtYWlsLkVtYW'
    'lsQWRkcmVzc0luc2lnaHRzTWFpbGJveEV2YWx1YXRpb25zUgtldmFsdWF0aW9ucxI/Cgdpc3Zh'
    'bGlkGIKH/k0gASgLMiIuZW1haWwuRW1haWxBZGRyZXNzSW5zaWdodHNWZXJkaWN0Ugdpc3ZhbG'
    'lk');

@$core.Deprecated('Use messageDescriptor instead')
const Message$json = {
  '1': 'Message',
  '2': [
    {
      '1': 'attachments',
      '3': 498946338,
      '4': 3,
      '5': 11,
      '6': '.email.Attachment',
      '10': 'attachments'
    },
    {
      '1': 'body',
      '3': 42602646,
      '4': 1,
      '5': 11,
      '6': '.email.Body',
      '10': 'body'
    },
    {
      '1': 'headers',
      '3': 323967370,
      '4': 3,
      '5': 11,
      '6': '.email.MessageHeader',
      '10': 'headers'
    },
    {
      '1': 'subject',
      '3': 7939312,
      '4': 1,
      '5': 11,
      '6': '.email.Content',
      '10': 'subject'
    },
  ],
};

/// Descriptor for `Message`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageDescriptor = $convert.base64Decode(
    'CgdNZXNzYWdlEjcKC2F0dGFjaG1lbnRzGKKi9e0BIAMoCzIRLmVtYWlsLkF0dGFjaG1lbnRSC2'
    'F0dGFjaG1lbnRzEiIKBGJvZHkYlqGoFCABKAsyCy5lbWFpbC5Cb2R5UgRib2R5EjIKB2hlYWRl'
    'cnMYirO9mgEgAygLMhQuZW1haWwuTWVzc2FnZUhlYWRlclIHaGVhZGVycxIrCgdzdWJqZWN0GP'
    'DJ5AMgASgLMg4uZW1haWwuQ29udGVudFIHc3ViamVjdA==');

@$core.Deprecated('Use messageHeaderDescriptor instead')
const MessageHeader$json = {
  '1': 'MessageHeader',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `MessageHeader`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageHeaderDescriptor = $convert.base64Decode(
    'Cg1NZXNzYWdlSGVhZGVyEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGAoFdmFsdWUY6/KfigEgAS'
    'gJUgV2YWx1ZQ==');

@$core.Deprecated('Use messageInsightsDataSourceDescriptor instead')
const MessageInsightsDataSource$json = {
  '1': 'MessageInsightsDataSource',
  '2': [
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'exclude',
      '3': 330273354,
      '4': 1,
      '5': 11,
      '6': '.email.MessageInsightsFilters',
      '10': 'exclude'
    },
    {
      '1': 'include',
      '3': 100816392,
      '4': 1,
      '5': 11,
      '6': '.email.MessageInsightsFilters',
      '10': 'include'
    },
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
  ],
};

/// Descriptor for `MessageInsightsDataSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageInsightsDataSourceDescriptor = $convert.base64Decode(
    'ChlNZXNzYWdlSW5zaWdodHNEYXRhU291cmNlEhsKB2VuZGRhdGUYz7P5JCABKAlSB2VuZGRhdG'
    'USOwoHZXhjbHVkZRjKpL6dASABKAsyHS5lbWFpbC5NZXNzYWdlSW5zaWdodHNGaWx0ZXJzUgdl'
    'eGNsdWRlEjoKB2luY2x1ZGUYiKyJMCABKAsyHS5lbWFpbC5NZXNzYWdlSW5zaWdodHNGaWx0ZX'
    'JzUgdpbmNsdWRlEiIKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEiAKCXN0YXJ0'
    'ZGF0ZRj8+KDUASABKAlSCXN0YXJ0ZGF0ZQ==');

@$core.Deprecated('Use messageInsightsFiltersDescriptor instead')
const MessageInsightsFilters$json = {
  '1': 'MessageInsightsFilters',
  '2': [
    {'1': 'destination', '3': 457443680, '4': 3, '5': 9, '10': 'destination'},
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 3,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {'1': 'isp', '3': 177492896, '4': 3, '5': 9, '10': 'isp'},
    {
      '1': 'lastdeliveryevent',
      '3': 528148396,
      '4': 3,
      '5': 14,
      '6': '.email.DeliveryEventType',
      '10': 'lastdeliveryevent'
    },
    {
      '1': 'lastengagementevent',
      '3': 426898599,
      '4': 3,
      '5': 14,
      '6': '.email.EngagementEventType',
      '10': 'lastengagementevent'
    },
    {'1': 'subject', '3': 7939312, '4': 3, '5': 9, '10': 'subject'},
  ],
};

/// Descriptor for `MessageInsightsFilters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageInsightsFiltersDescriptor = $convert.base64Decode(
    'ChZNZXNzYWdlSW5zaWdodHNGaWx0ZXJzEiQKC2Rlc3RpbmF0aW9uGOCSkNoBIAMoCVILZGVzdG'
    'luYXRpb24SLQoQZnJvbWVtYWlsYWRkcmVzcxiGmsssIAMoCVIQZnJvbWVtYWlsYWRkcmVzcxIT'
    'CgNpc3AYoKfRVCADKAlSA2lzcBJKChFsYXN0ZGVsaXZlcnlldmVudBisz+v7ASADKA4yGC5lbW'
    'FpbC5EZWxpdmVyeUV2ZW50VHlwZVIRbGFzdGRlbGl2ZXJ5ZXZlbnQSUAoTbGFzdGVuZ2FnZW1l'
    'bnRldmVudBin6cfLASADKA4yGi5lbWFpbC5FbmdhZ2VtZW50RXZlbnRUeXBlUhNsYXN0ZW5nYW'
    'dlbWVudGV2ZW50EhsKB3N1YmplY3QY8MnkAyADKAlSB3N1YmplY3Q=');

@$core.Deprecated('Use messageRejectedDescriptor instead')
const MessageRejected$json = {
  '1': 'MessageRejected',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MessageRejected`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageRejectedDescriptor = $convert.base64Decode(
    'Cg9NZXNzYWdlUmVqZWN0ZWQSGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use messageTagDescriptor instead')
const MessageTag$json = {
  '1': 'MessageTag',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `MessageTag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List messageTagDescriptor = $convert.base64Decode(
    'CgpNZXNzYWdlVGFnEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGAoFdmFsdWUY6/KfigEgASgJUg'
    'V2YWx1ZQ==');

@$core.Deprecated('Use metricDataErrorDescriptor instead')
const MetricDataError$json = {
  '1': 'MetricDataError',
  '2': [
    {
      '1': 'code',
      '3': 425572629,
      '4': 1,
      '5': 14,
      '6': '.email.QueryErrorCode',
      '10': 'code'
    },
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `MetricDataError`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricDataErrorDescriptor = $convert.base64Decode(
    'Cg9NZXRyaWNEYXRhRXJyb3ISLQoEY29kZRiV8vbKASABKA4yFS5lbWFpbC5RdWVyeUVycm9yQ2'
    '9kZVIEY29kZRISCgJpZBiB8qK3ASABKAlSAmlkEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3Nh'
    'Z2U=');

@$core.Deprecated('Use metricDataResultDescriptor instead')
const MetricDataResult$json = {
  '1': 'MetricDataResult',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
    {'1': 'timestamps', '3': 213534737, '4': 3, '5': 9, '10': 'timestamps'},
    {'1': 'values', '3': 223158876, '4': 3, '5': 3, '10': 'values'},
  ],
};

/// Descriptor for `MetricDataResult`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricDataResultDescriptor = $convert.base64Decode(
    'ChBNZXRyaWNEYXRhUmVzdWx0EhIKAmlkGIHyorcBIAEoCVICaWQSIQoKdGltZXN0YW1wcxiRkO'
    'llIAMoCVIKdGltZXN0YW1wcxIZCgZ2YWx1ZXMY3MS0aiADKANSBnZhbHVlcw==');

@$core.Deprecated('Use metricsDataSourceDescriptor instead')
const MetricsDataSource$json = {
  '1': 'MetricsDataSource',
  '2': [
    {
      '1': 'dimensions',
      '3': 462933457,
      '4': 3,
      '5': 11,
      '6': '.email.MetricsDataSource.DimensionsEntry',
      '10': 'dimensions'
    },
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'metrics',
      '3': 436365847,
      '4': 3,
      '5': 11,
      '6': '.email.ExportMetric',
      '10': 'metrics'
    },
    {
      '1': 'namespace',
      '3': 355353153,
      '4': 1,
      '5': 14,
      '6': '.email.MetricNamespace',
      '10': 'namespace'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
  ],
  '3': [MetricsDataSource_DimensionsEntry$json],
};

@$core.Deprecated('Use metricsDataSourceDescriptor instead')
const MetricsDataSource_DimensionsEntry$json = {
  '1': 'DimensionsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `MetricsDataSource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List metricsDataSourceDescriptor = $convert.base64Decode(
    'ChFNZXRyaWNzRGF0YVNvdXJjZRJMCgpkaW1lbnNpb25zGNGb39wBIAMoCzIoLmVtYWlsLk1ldH'
    'JpY3NEYXRhU291cmNlLkRpbWVuc2lvbnNFbnRyeVIKZGltZW5zaW9ucxIbCgdlbmRkYXRlGM+z'
    '+SQgASgJUgdlbmRkYXRlEjEKB21ldHJpY3MYl9SJ0AEgAygLMhMuZW1haWwuRXhwb3J0TWV0cm'
    'ljUgdtZXRyaWNzEjgKCW5hbWVzcGFjZRjBhLmpASABKA4yFi5lbWFpbC5NZXRyaWNOYW1lc3Bh'
    'Y2VSCW5hbWVzcGFjZRIgCglzdGFydGRhdGUY/Pig1AEgASgJUglzdGFydGRhdGUaPQoPRGltZW'
    '5zaW9uc0VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use multiRegionEndpointDescriptor instead')
const MultiRegionEndpoint$json = {
  '1': 'MultiRegionEndpoint',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'endpointid', '3': 35808946, '4': 1, '5': 9, '10': 'endpointid'},
    {'1': 'endpointname', '3': 209534392, '4': 1, '5': 9, '10': 'endpointname'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {'1': 'regions', '3': 36200107, '4': 3, '5': 9, '10': 'regions'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.Status',
      '10': 'status'
    },
  ],
};

/// Descriptor for `MultiRegionEndpoint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List multiRegionEndpointDescriptor = $convert.base64Decode(
    'ChNNdWx0aVJlZ2lvbkVuZHBvaW50Ei4KEGNyZWF0ZWR0aW1lc3RhbXAY+tvPnwEgASgJUhBjcm'
    'VhdGVkdGltZXN0YW1wEiEKCmVuZHBvaW50aWQYss2JESABKAlSCmVuZHBvaW50aWQSJQoMZW5k'
    'cG9pbnRuYW1lGLj79GMgASgJUgxlbmRwb2ludG5hbWUSNQoUbGFzdHVwZGF0ZWR0aW1lc3RhbX'
    'AYlcvIPyABKAlSFGxhc3R1cGRhdGVkdGltZXN0YW1wEhsKB3JlZ2lvbnMYq72hESADKAlSB3Jl'
    'Z2lvbnMSKAoGc3RhdHVzGJDk+wIgASgOMg0uZW1haWwuU3RhdHVzUgZzdGF0dXM=');

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

@$core.Deprecated('Use overallVolumeDescriptor instead')
const OverallVolume$json = {
  '1': 'OverallVolume',
  '2': [
    {
      '1': 'domainispplacements',
      '3': 63848508,
      '4': 3,
      '5': 11,
      '6': '.email.DomainIspPlacement',
      '10': 'domainispplacements'
    },
    {
      '1': 'readratepercent',
      '3': 458603677,
      '4': 1,
      '5': 1,
      '10': 'readratepercent'
    },
    {
      '1': 'volumestatistics',
      '3': 331663309,
      '4': 1,
      '5': 11,
      '6': '.email.VolumeStatistics',
      '10': 'volumestatistics'
    },
  ],
};

/// Descriptor for `OverallVolume`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List overallVolumeDescriptor = $convert.base64Decode(
    'Cg1PdmVyYWxsVm9sdW1lEk4KE2RvbWFpbmlzcHBsYWNlbWVudHMYvIC5HiADKAsyGS5lbWFpbC'
    '5Eb21haW5Jc3BQbGFjZW1lbnRSE2RvbWFpbmlzcHBsYWNlbWVudHMSLAoPcmVhZHJhdGVwZXJj'
    'ZW50GJ351toBIAEoAVIPcmVhZHJhdGVwZXJjZW50EkcKEHZvbHVtZXN0YXRpc3RpY3MYzY+Tng'
    'EgASgLMhcuZW1haWwuVm9sdW1lU3RhdGlzdGljc1IQdm9sdW1lc3RhdGlzdGljcw==');

@$core.Deprecated('Use pinpointDestinationDescriptor instead')
const PinpointDestination$json = {
  '1': 'PinpointDestination',
  '2': [
    {
      '1': 'applicationarn',
      '3': 524185341,
      '4': 1,
      '5': 9,
      '10': 'applicationarn'
    },
  ],
};

/// Descriptor for `PinpointDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List pinpointDestinationDescriptor = $convert.base64Decode(
    'ChNQaW5wb2ludERlc3RpbmF0aW9uEioKDmFwcGxpY2F0aW9uYXJuGP3d+fkBIAEoCVIOYXBwbG'
    'ljYXRpb25hcm4=');

@$core.Deprecated('Use placementStatisticsDescriptor instead')
const PlacementStatistics$json = {
  '1': 'PlacementStatistics',
  '2': [
    {
      '1': 'dkimpercentage',
      '3': 434155899,
      '4': 1,
      '5': 1,
      '10': 'dkimpercentage'
    },
    {
      '1': 'inboxpercentage',
      '3': 81442932,
      '4': 1,
      '5': 1,
      '10': 'inboxpercentage'
    },
    {
      '1': 'missingpercentage',
      '3': 466437664,
      '4': 1,
      '5': 1,
      '10': 'missingpercentage'
    },
    {
      '1': 'spampercentage',
      '3': 364813289,
      '4': 1,
      '5': 1,
      '10': 'spampercentage'
    },
    {
      '1': 'spfpercentage',
      '3': 52331233,
      '4': 1,
      '5': 1,
      '10': 'spfpercentage'
    },
  ],
};

/// Descriptor for `PlacementStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List placementStatisticsDescriptor = $convert.base64Decode(
    'ChNQbGFjZW1lbnRTdGF0aXN0aWNzEioKDmRraW1wZXJjZW50YWdlGPvigs8BIAEoAVIOZGtpbX'
    'BlcmNlbnRhZ2USKwoPaW5ib3hwZXJjZW50YWdlGPTw6iYgASgBUg9pbmJveHBlcmNlbnRhZ2US'
    'MAoRbWlzc2luZ3BlcmNlbnRhZ2UYoIy13gEgASgBUhFtaXNzaW5ncGVyY2VudGFnZRIqCg5zcG'
    'FtcGVyY2VudGFnZRjpt/qtASABKAFSDnNwYW1wZXJjZW50YWdlEicKDXNwZnBlcmNlbnRhZ2UY'
    '4YX6GCABKAFSDXNwZnBlcmNlbnRhZ2U=');

@$core.Deprecated(
    'Use putAccountDedicatedIpWarmupAttributesRequestDescriptor instead')
const PutAccountDedicatedIpWarmupAttributesRequest$json = {
  '1': 'PutAccountDedicatedIpWarmupAttributesRequest',
  '2': [
    {
      '1': 'autowarmupenabled',
      '3': 270184628,
      '4': 1,
      '5': 8,
      '10': 'autowarmupenabled'
    },
  ],
};

/// Descriptor for `PutAccountDedicatedIpWarmupAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putAccountDedicatedIpWarmupAttributesRequestDescriptor =
    $convert.base64Decode(
        'CixQdXRBY2NvdW50RGVkaWNhdGVkSXBXYXJtdXBBdHRyaWJ1dGVzUmVxdWVzdBIwChFhdXRvd2'
        'FybXVwZW5hYmxlZBi04eqAASABKAhSEWF1dG93YXJtdXBlbmFibGVk');

@$core.Deprecated(
    'Use putAccountDedicatedIpWarmupAttributesResponseDescriptor instead')
const PutAccountDedicatedIpWarmupAttributesResponse$json = {
  '1': 'PutAccountDedicatedIpWarmupAttributesResponse',
};

/// Descriptor for `PutAccountDedicatedIpWarmupAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putAccountDedicatedIpWarmupAttributesResponseDescriptor =
    $convert.base64Decode(
        'Ci1QdXRBY2NvdW50RGVkaWNhdGVkSXBXYXJtdXBBdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated('Use putAccountDetailsRequestDescriptor instead')
const PutAccountDetailsRequest$json = {
  '1': 'PutAccountDetailsRequest',
  '2': [
    {
      '1': 'additionalcontactemailaddresses',
      '3': 322089693,
      '4': 3,
      '5': 9,
      '10': 'additionalcontactemailaddresses'
    },
    {
      '1': 'contactlanguage',
      '3': 114240022,
      '4': 1,
      '5': 14,
      '6': '.email.ContactLanguage',
      '10': 'contactlanguage'
    },
    {
      '1': 'mailtype',
      '3': 138144527,
      '4': 1,
      '5': 14,
      '6': '.email.MailType',
      '10': 'mailtype'
    },
    {
      '1': 'productionaccessenabled',
      '3': 471167534,
      '4': 1,
      '5': 8,
      '10': 'productionaccessenabled'
    },
    {
      '1': 'usecasedescription',
      '3': 141053987,
      '4': 1,
      '5': 9,
      '10': 'usecasedescription'
    },
    {'1': 'websiteurl', '3': 201971828, '4': 1, '5': 9, '10': 'websiteurl'},
  ],
};

/// Descriptor for `PutAccountDetailsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountDetailsRequestDescriptor = $convert.base64Decode(
    'ChhQdXRBY2NvdW50RGV0YWlsc1JlcXVlc3QSTAofYWRkaXRpb25hbGNvbnRhY3RlbWFpbGFkZH'
    'Jlc3Nlcxjd5cqZASADKAlSH2FkZGl0aW9uYWxjb250YWN0ZW1haWxhZGRyZXNzZXMSQwoPY29u'
    'dGFjdGxhbmd1YWdlGJbUvDYgASgOMhYuZW1haWwuQ29udGFjdExhbmd1YWdlUg9jb250YWN0bG'
    'FuZ3VhZ2USLgoIbWFpbHR5cGUYj9bvQSABKA4yDy5lbWFpbC5NYWlsVHlwZVIIbWFpbHR5cGUS'
    'PAoXcHJvZHVjdGlvbmFjY2Vzc2VuYWJsZWQYruTV4AEgASgIUhdwcm9kdWN0aW9uYWNjZXNzZW'
    '5hYmxlZBIxChJ1c2VjYXNlZGVzY3JpcHRpb24Yo6ChQyABKAlSEnVzZWNhc2VkZXNjcmlwdGlv'
    'bhIhCgp3ZWJzaXRldXJsGPSwp2AgASgJUgp3ZWJzaXRldXJs');

@$core.Deprecated('Use putAccountDetailsResponseDescriptor instead')
const PutAccountDetailsResponse$json = {
  '1': 'PutAccountDetailsResponse',
};

/// Descriptor for `PutAccountDetailsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountDetailsResponseDescriptor =
    $convert.base64Decode('ChlQdXRBY2NvdW50RGV0YWlsc1Jlc3BvbnNl');

@$core.Deprecated('Use putAccountSendingAttributesRequestDescriptor instead')
const PutAccountSendingAttributesRequest$json = {
  '1': 'PutAccountSendingAttributesRequest',
  '2': [
    {
      '1': 'sendingenabled',
      '3': 194846115,
      '4': 1,
      '5': 8,
      '10': 'sendingenabled'
    },
  ],
};

/// Descriptor for `PutAccountSendingAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountSendingAttributesRequestDescriptor =
    $convert.base64Decode(
        'CiJQdXRBY2NvdW50U2VuZGluZ0F0dHJpYnV0ZXNSZXF1ZXN0EikKDnNlbmRpbmdlbmFibGVkGK'
        'O79FwgASgIUg5zZW5kaW5nZW5hYmxlZA==');

@$core.Deprecated('Use putAccountSendingAttributesResponseDescriptor instead')
const PutAccountSendingAttributesResponse$json = {
  '1': 'PutAccountSendingAttributesResponse',
};

/// Descriptor for `PutAccountSendingAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountSendingAttributesResponseDescriptor =
    $convert
        .base64Decode('CiNQdXRBY2NvdW50U2VuZGluZ0F0dHJpYnV0ZXNSZXNwb25zZQ==');

@$core
    .Deprecated('Use putAccountSuppressionAttributesRequestDescriptor instead')
const PutAccountSuppressionAttributesRequest$json = {
  '1': 'PutAccountSuppressionAttributesRequest',
  '2': [
    {
      '1': 'suppressedreasons',
      '3': 465922417,
      '4': 3,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'suppressedreasons'
    },
    {
      '1': 'validationattributes',
      '3': 171836944,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionValidationAttributes',
      '10': 'validationattributes'
    },
  ],
};

/// Descriptor for `PutAccountSuppressionAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountSuppressionAttributesRequestDescriptor =
    $convert.base64Decode(
        'CiZQdXRBY2NvdW50U3VwcHJlc3Npb25BdHRyaWJ1dGVzUmVxdWVzdBJOChFzdXBwcmVzc2Vkcm'
        'Vhc29ucxjx0pXeASADKA4yHC5lbWFpbC5TdXBwcmVzc2lvbkxpc3RSZWFzb25SEXN1cHByZXNz'
        'ZWRyZWFzb25zEl0KFHZhbGlkYXRpb25hdHRyaWJ1dGVzGJCM+FEgASgLMiYuZW1haWwuU3VwcH'
        'Jlc3Npb25WYWxpZGF0aW9uQXR0cmlidXRlc1IUdmFsaWRhdGlvbmF0dHJpYnV0ZXM=');

@$core
    .Deprecated('Use putAccountSuppressionAttributesResponseDescriptor instead')
const PutAccountSuppressionAttributesResponse$json = {
  '1': 'PutAccountSuppressionAttributesResponse',
};

/// Descriptor for `PutAccountSuppressionAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountSuppressionAttributesResponseDescriptor =
    $convert.base64Decode(
        'CidQdXRBY2NvdW50U3VwcHJlc3Npb25BdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated('Use putAccountVdmAttributesRequestDescriptor instead')
const PutAccountVdmAttributesRequest$json = {
  '1': 'PutAccountVdmAttributesRequest',
  '2': [
    {
      '1': 'vdmattributes',
      '3': 506490540,
      '4': 1,
      '5': 11,
      '6': '.email.VdmAttributes',
      '10': 'vdmattributes'
    },
  ],
};

/// Descriptor for `PutAccountVdmAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountVdmAttributesRequestDescriptor =
    $convert.base64Decode(
        'Ch5QdXRBY2NvdW50VmRtQXR0cmlidXRlc1JlcXVlc3QSPgoNdmRtYXR0cmlidXRlcxis3cHxAS'
        'ABKAsyFC5lbWFpbC5WZG1BdHRyaWJ1dGVzUg12ZG1hdHRyaWJ1dGVz');

@$core.Deprecated('Use putAccountVdmAttributesResponseDescriptor instead')
const PutAccountVdmAttributesResponse$json = {
  '1': 'PutAccountVdmAttributesResponse',
};

/// Descriptor for `PutAccountVdmAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putAccountVdmAttributesResponseDescriptor =
    $convert.base64Decode('Ch9QdXRBY2NvdW50VmRtQXR0cmlidXRlc1Jlc3BvbnNl');

@$core.Deprecated(
    'Use putConfigurationSetArchivingOptionsRequestDescriptor instead')
const PutConfigurationSetArchivingOptionsRequest$json = {
  '1': 'PutConfigurationSetArchivingOptionsRequest',
  '2': [
    {'1': 'archivearn', '3': 56866685, '4': 1, '5': 9, '10': 'archivearn'},
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
  ],
};

/// Descriptor for `PutConfigurationSetArchivingOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetArchivingOptionsRequestDescriptor =
    $convert.base64Decode(
        'CipQdXRDb25maWd1cmF0aW9uU2V0QXJjaGl2aW5nT3B0aW9uc1JlcXVlc3QSIQoKYXJjaGl2ZW'
        'Fybhj97o4bIAEoCVIKYXJjaGl2ZWFybhI2ChRjb25maWd1cmF0aW9uc2V0bmFtZRjNi7HAASAB'
        'KAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1l');

@$core.Deprecated(
    'Use putConfigurationSetArchivingOptionsResponseDescriptor instead')
const PutConfigurationSetArchivingOptionsResponse$json = {
  '1': 'PutConfigurationSetArchivingOptionsResponse',
};

/// Descriptor for `PutConfigurationSetArchivingOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetArchivingOptionsResponseDescriptor =
    $convert.base64Decode(
        'CitQdXRDb25maWd1cmF0aW9uU2V0QXJjaGl2aW5nT3B0aW9uc1Jlc3BvbnNl');

@$core.Deprecated(
    'Use putConfigurationSetDeliveryOptionsRequestDescriptor instead')
const PutConfigurationSetDeliveryOptionsRequest$json = {
  '1': 'PutConfigurationSetDeliveryOptionsRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'maxdeliveryseconds',
      '3': 16229983,
      '4': 1,
      '5': 3,
      '10': 'maxdeliveryseconds'
    },
    {
      '1': 'sendingpoolname',
      '3': 89398333,
      '4': 1,
      '5': 9,
      '10': 'sendingpoolname'
    },
    {
      '1': 'tlspolicy',
      '3': 127629,
      '4': 1,
      '5': 14,
      '6': '.email.TlsPolicy',
      '10': 'tlspolicy'
    },
  ],
};

/// Descriptor for `PutConfigurationSetDeliveryOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetDeliveryOptionsRequestDescriptor = $convert.base64Decode(
        'CilQdXRDb25maWd1cmF0aW9uU2V0RGVsaXZlcnlPcHRpb25zUmVxdWVzdBI2ChRjb25maWd1cm'
        'F0aW9uc2V0bmFtZRjNi7HAASABKAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1lEjEKEm1heGRlbGl2'
        'ZXJ5c2Vjb25kcxjfzN4HIAEoA1ISbWF4ZGVsaXZlcnlzZWNvbmRzEisKD3NlbmRpbmdwb29sbm'
        'FtZRi9uNAqIAEoCVIPc2VuZGluZ3Bvb2xuYW1lEjAKCXRsc3BvbGljeRiN5QcgASgOMhAuZW1h'
        'aWwuVGxzUG9saWN5Ugl0bHNwb2xpY3k=');

@$core.Deprecated(
    'Use putConfigurationSetDeliveryOptionsResponseDescriptor instead')
const PutConfigurationSetDeliveryOptionsResponse$json = {
  '1': 'PutConfigurationSetDeliveryOptionsResponse',
};

/// Descriptor for `PutConfigurationSetDeliveryOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetDeliveryOptionsResponseDescriptor =
    $convert.base64Decode(
        'CipQdXRDb25maWd1cmF0aW9uU2V0RGVsaXZlcnlPcHRpb25zUmVzcG9uc2U=');

@$core.Deprecated(
    'Use putConfigurationSetReputationOptionsRequestDescriptor instead')
const PutConfigurationSetReputationOptionsRequest$json = {
  '1': 'PutConfigurationSetReputationOptionsRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'reputationmetricsenabled',
      '3': 440357917,
      '4': 1,
      '5': 8,
      '10': 'reputationmetricsenabled'
    },
  ],
};

/// Descriptor for `PutConfigurationSetReputationOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetReputationOptionsRequestDescriptor =
    $convert.base64Decode(
        'CitQdXRDb25maWd1cmF0aW9uU2V0UmVwdXRhdGlvbk9wdGlvbnNSZXF1ZXN0EjYKFGNvbmZpZ3'
        'VyYXRpb25zZXRuYW1lGM2LscABIAEoCVIUY29uZmlndXJhdGlvbnNldG5hbWUSPgoYcmVwdXRh'
        'dGlvbm1ldHJpY3NlbmFibGVkGJ2o/dEBIAEoCFIYcmVwdXRhdGlvbm1ldHJpY3NlbmFibGVk');

@$core.Deprecated(
    'Use putConfigurationSetReputationOptionsResponseDescriptor instead')
const PutConfigurationSetReputationOptionsResponse$json = {
  '1': 'PutConfigurationSetReputationOptionsResponse',
};

/// Descriptor for `PutConfigurationSetReputationOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetReputationOptionsResponseDescriptor =
    $convert.base64Decode(
        'CixQdXRDb25maWd1cmF0aW9uU2V0UmVwdXRhdGlvbk9wdGlvbnNSZXNwb25zZQ==');

@$core.Deprecated(
    'Use putConfigurationSetSendingOptionsRequestDescriptor instead')
const PutConfigurationSetSendingOptionsRequest$json = {
  '1': 'PutConfigurationSetSendingOptionsRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'sendingenabled',
      '3': 194846115,
      '4': 1,
      '5': 8,
      '10': 'sendingenabled'
    },
  ],
};

/// Descriptor for `PutConfigurationSetSendingOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putConfigurationSetSendingOptionsRequestDescriptor =
    $convert.base64Decode(
        'CihQdXRDb25maWd1cmF0aW9uU2V0U2VuZGluZ09wdGlvbnNSZXF1ZXN0EjYKFGNvbmZpZ3VyYX'
        'Rpb25zZXRuYW1lGM2LscABIAEoCVIUY29uZmlndXJhdGlvbnNldG5hbWUSKQoOc2VuZGluZ2Vu'
        'YWJsZWQYo7v0XCABKAhSDnNlbmRpbmdlbmFibGVk');

@$core.Deprecated(
    'Use putConfigurationSetSendingOptionsResponseDescriptor instead')
const PutConfigurationSetSendingOptionsResponse$json = {
  '1': 'PutConfigurationSetSendingOptionsResponse',
};

/// Descriptor for `PutConfigurationSetSendingOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetSendingOptionsResponseDescriptor = $convert.base64Decode(
        'CilQdXRDb25maWd1cmF0aW9uU2V0U2VuZGluZ09wdGlvbnNSZXNwb25zZQ==');

@$core.Deprecated(
    'Use putConfigurationSetSuppressionOptionsRequestDescriptor instead')
const PutConfigurationSetSuppressionOptionsRequest$json = {
  '1': 'PutConfigurationSetSuppressionOptionsRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'suppressedreasons',
      '3': 465922417,
      '4': 3,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'suppressedreasons'
    },
    {
      '1': 'validationoptions',
      '3': 216495637,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionValidationOptions',
      '10': 'validationoptions'
    },
  ],
};

/// Descriptor for `PutConfigurationSetSuppressionOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetSuppressionOptionsRequestDescriptor =
    $convert.base64Decode(
        'CixQdXRDb25maWd1cmF0aW9uU2V0U3VwcHJlc3Npb25PcHRpb25zUmVxdWVzdBI2ChRjb25maW'
        'd1cmF0aW9uc2V0bmFtZRjNi7HAASABKAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1lEk4KEXN1cHBy'
        'ZXNzZWRyZWFzb25zGPHSld4BIAMoDjIcLmVtYWlsLlN1cHByZXNzaW9uTGlzdFJlYXNvblIRc3'
        'VwcHJlc3NlZHJlYXNvbnMSVAoRdmFsaWRhdGlvbm9wdGlvbnMYleydZyABKAsyIy5lbWFpbC5T'
        'dXBwcmVzc2lvblZhbGlkYXRpb25PcHRpb25zUhF2YWxpZGF0aW9ub3B0aW9ucw==');

@$core.Deprecated(
    'Use putConfigurationSetSuppressionOptionsResponseDescriptor instead')
const PutConfigurationSetSuppressionOptionsResponse$json = {
  '1': 'PutConfigurationSetSuppressionOptionsResponse',
};

/// Descriptor for `PutConfigurationSetSuppressionOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetSuppressionOptionsResponseDescriptor =
    $convert.base64Decode(
        'Ci1QdXRDb25maWd1cmF0aW9uU2V0U3VwcHJlc3Npb25PcHRpb25zUmVzcG9uc2U=');

@$core.Deprecated(
    'Use putConfigurationSetTrackingOptionsRequestDescriptor instead')
const PutConfigurationSetTrackingOptionsRequest$json = {
  '1': 'PutConfigurationSetTrackingOptionsRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'customredirectdomain',
      '3': 72478039,
      '4': 1,
      '5': 9,
      '10': 'customredirectdomain'
    },
    {
      '1': 'httpspolicy',
      '3': 147115397,
      '4': 1,
      '5': 14,
      '6': '.email.HttpsPolicy',
      '10': 'httpspolicy'
    },
  ],
};

/// Descriptor for `PutConfigurationSetTrackingOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetTrackingOptionsRequestDescriptor = $convert.base64Decode(
        'CilQdXRDb25maWd1cmF0aW9uU2V0VHJhY2tpbmdPcHRpb25zUmVxdWVzdBI2ChRjb25maWd1cm'
        'F0aW9uc2V0bmFtZRjNi7HAASABKAlSFGNvbmZpZ3VyYXRpb25zZXRuYW1lEjUKFGN1c3RvbXJl'
        'ZGlyZWN0ZG9tYWluGNfaxyIgASgJUhRjdXN0b21yZWRpcmVjdGRvbWFpbhI3CgtodHRwc3BvbG'
        'ljeRiFm5NGIAEoDjISLmVtYWlsLkh0dHBzUG9saWN5UgtodHRwc3BvbGljeQ==');

@$core.Deprecated(
    'Use putConfigurationSetTrackingOptionsResponseDescriptor instead')
const PutConfigurationSetTrackingOptionsResponse$json = {
  '1': 'PutConfigurationSetTrackingOptionsResponse',
};

/// Descriptor for `PutConfigurationSetTrackingOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putConfigurationSetTrackingOptionsResponseDescriptor =
    $convert.base64Decode(
        'CipQdXRDb25maWd1cmF0aW9uU2V0VHJhY2tpbmdPcHRpb25zUmVzcG9uc2U=');

@$core.Deprecated('Use putConfigurationSetVdmOptionsRequestDescriptor instead')
const PutConfigurationSetVdmOptionsRequest$json = {
  '1': 'PutConfigurationSetVdmOptionsRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'vdmoptions',
      '3': 344971073,
      '4': 1,
      '5': 11,
      '6': '.email.VdmOptions',
      '10': 'vdmoptions'
    },
  ],
};

/// Descriptor for `PutConfigurationSetVdmOptionsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putConfigurationSetVdmOptionsRequestDescriptor =
    $convert.base64Decode(
        'CiRQdXRDb25maWd1cmF0aW9uU2V0VmRtT3B0aW9uc1JlcXVlc3QSNgoUY29uZmlndXJhdGlvbn'
        'NldG5hbWUYzYuxwAEgASgJUhRjb25maWd1cmF0aW9uc2V0bmFtZRI1Cgp2ZG1vcHRpb25zGMGu'
        'v6QBIAEoCzIRLmVtYWlsLlZkbU9wdGlvbnNSCnZkbW9wdGlvbnM=');

@$core.Deprecated('Use putConfigurationSetVdmOptionsResponseDescriptor instead')
const PutConfigurationSetVdmOptionsResponse$json = {
  '1': 'PutConfigurationSetVdmOptionsResponse',
};

/// Descriptor for `PutConfigurationSetVdmOptionsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putConfigurationSetVdmOptionsResponseDescriptor =
    $convert
        .base64Decode('CiVQdXRDb25maWd1cmF0aW9uU2V0VmRtT3B0aW9uc1Jlc3BvbnNl');

@$core.Deprecated('Use putDedicatedIpInPoolRequestDescriptor instead')
const PutDedicatedIpInPoolRequest$json = {
  '1': 'PutDedicatedIpInPoolRequest',
  '2': [
    {
      '1': 'destinationpoolname',
      '3': 49627479,
      '4': 1,
      '5': 9,
      '10': 'destinationpoolname'
    },
    {'1': 'ip', '3': 183031933, '4': 1, '5': 9, '10': 'ip'},
  ],
};

/// Descriptor for `PutDedicatedIpInPoolRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDedicatedIpInPoolRequestDescriptor =
    $convert.base64Decode(
        'ChtQdXREZWRpY2F0ZWRJcEluUG9vbFJlcXVlc3QSMwoTZGVzdGluYXRpb25wb29sbmFtZRjXgt'
        'UXIAEoCVITZGVzdGluYXRpb25wb29sbmFtZRIRCgJpcBj9sKNXIAEoCVICaXA=');

@$core.Deprecated('Use putDedicatedIpInPoolResponseDescriptor instead')
const PutDedicatedIpInPoolResponse$json = {
  '1': 'PutDedicatedIpInPoolResponse',
};

/// Descriptor for `PutDedicatedIpInPoolResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDedicatedIpInPoolResponseDescriptor =
    $convert.base64Decode('ChxQdXREZWRpY2F0ZWRJcEluUG9vbFJlc3BvbnNl');

@$core.Deprecated(
    'Use putDedicatedIpPoolScalingAttributesRequestDescriptor instead')
const PutDedicatedIpPoolScalingAttributesRequest$json = {
  '1': 'PutDedicatedIpPoolScalingAttributesRequest',
  '2': [
    {'1': 'poolname', '3': 81872585, '4': 1, '5': 9, '10': 'poolname'},
    {
      '1': 'scalingmode',
      '3': 210356138,
      '4': 1,
      '5': 14,
      '6': '.email.ScalingMode',
      '10': 'scalingmode'
    },
  ],
};

/// Descriptor for `PutDedicatedIpPoolScalingAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putDedicatedIpPoolScalingAttributesRequestDescriptor =
    $convert.base64Decode(
        'CipQdXREZWRpY2F0ZWRJcFBvb2xTY2FsaW5nQXR0cmlidXRlc1JlcXVlc3QSHQoIcG9vbG5hbW'
        'UYyY2FJyABKAlSCHBvb2xuYW1lEjcKC3NjYWxpbmdtb2RlGKqPp2QgASgOMhIuZW1haWwuU2Nh'
        'bGluZ01vZGVSC3NjYWxpbmdtb2Rl');

@$core.Deprecated(
    'Use putDedicatedIpPoolScalingAttributesResponseDescriptor instead')
const PutDedicatedIpPoolScalingAttributesResponse$json = {
  '1': 'PutDedicatedIpPoolScalingAttributesResponse',
};

/// Descriptor for `PutDedicatedIpPoolScalingAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putDedicatedIpPoolScalingAttributesResponseDescriptor =
    $convert.base64Decode(
        'CitQdXREZWRpY2F0ZWRJcFBvb2xTY2FsaW5nQXR0cmlidXRlc1Jlc3BvbnNl');

@$core.Deprecated('Use putDedicatedIpWarmupAttributesRequestDescriptor instead')
const PutDedicatedIpWarmupAttributesRequest$json = {
  '1': 'PutDedicatedIpWarmupAttributesRequest',
  '2': [
    {'1': 'ip', '3': 183031933, '4': 1, '5': 9, '10': 'ip'},
    {
      '1': 'warmuppercentage',
      '3': 389585846,
      '4': 1,
      '5': 5,
      '10': 'warmuppercentage'
    },
  ],
};

/// Descriptor for `PutDedicatedIpWarmupAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDedicatedIpWarmupAttributesRequestDescriptor =
    $convert.base64Decode(
        'CiVQdXREZWRpY2F0ZWRJcFdhcm11cEF0dHJpYnV0ZXNSZXF1ZXN0EhEKAmlwGP2wo1cgASgJUg'
        'JpcBIuChB3YXJtdXBwZXJjZW50YWdlGLa34rkBIAEoBVIQd2FybXVwcGVyY2VudGFnZQ==');

@$core
    .Deprecated('Use putDedicatedIpWarmupAttributesResponseDescriptor instead')
const PutDedicatedIpWarmupAttributesResponse$json = {
  '1': 'PutDedicatedIpWarmupAttributesResponse',
};

/// Descriptor for `PutDedicatedIpWarmupAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDedicatedIpWarmupAttributesResponseDescriptor =
    $convert.base64Decode(
        'CiZQdXREZWRpY2F0ZWRJcFdhcm11cEF0dHJpYnV0ZXNSZXNwb25zZQ==');

@$core
    .Deprecated('Use putDeliverabilityDashboardOptionRequestDescriptor instead')
const PutDeliverabilityDashboardOptionRequest$json = {
  '1': 'PutDeliverabilityDashboardOptionRequest',
  '2': [
    {
      '1': 'dashboardenabled',
      '3': 90846529,
      '4': 1,
      '5': 8,
      '10': 'dashboardenabled'
    },
    {
      '1': 'subscribeddomains',
      '3': 243493831,
      '4': 3,
      '5': 11,
      '6': '.email.DomainDeliverabilityTrackingOption',
      '10': 'subscribeddomains'
    },
  ],
};

/// Descriptor for `PutDeliverabilityDashboardOptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliverabilityDashboardOptionRequestDescriptor =
    $convert.base64Decode(
        'CidQdXREZWxpdmVyYWJpbGl0eURhc2hib2FyZE9wdGlvblJlcXVlc3QSLQoQZGFzaGJvYXJkZW'
        '5hYmxlZBjB6qgrIAEoCFIQZGFzaGJvYXJkZW5hYmxlZBJaChFzdWJzY3JpYmVkZG9tYWlucxjH'
        '1410IAMoCzIpLmVtYWlsLkRvbWFpbkRlbGl2ZXJhYmlsaXR5VHJhY2tpbmdPcHRpb25SEXN1Yn'
        'NjcmliZWRkb21haW5z');

@$core.Deprecated(
    'Use putDeliverabilityDashboardOptionResponseDescriptor instead')
const PutDeliverabilityDashboardOptionResponse$json = {
  '1': 'PutDeliverabilityDashboardOptionResponse',
};

/// Descriptor for `PutDeliverabilityDashboardOptionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putDeliverabilityDashboardOptionResponseDescriptor =
    $convert.base64Decode(
        'CihQdXREZWxpdmVyYWJpbGl0eURhc2hib2FyZE9wdGlvblJlc3BvbnNl');

@$core.Deprecated(
    'Use putEmailIdentityConfigurationSetAttributesRequestDescriptor instead')
const PutEmailIdentityConfigurationSetAttributesRequest$json = {
  '1': 'PutEmailIdentityConfigurationSetAttributesRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
  ],
};

/// Descriptor for `PutEmailIdentityConfigurationSetAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityConfigurationSetAttributesRequestDescriptor =
    $convert.base64Decode(
        'CjFQdXRFbWFpbElkZW50aXR5Q29uZmlndXJhdGlvblNldEF0dHJpYnV0ZXNSZXF1ZXN0EjYKFG'
        'NvbmZpZ3VyYXRpb25zZXRuYW1lGM2LscABIAEoCVIUY29uZmlndXJhdGlvbnNldG5hbWUSJwoN'
        'ZW1haWxpZGVudGl0eRialPJAIAEoCVINZW1haWxpZGVudGl0eQ==');

@$core.Deprecated(
    'Use putEmailIdentityConfigurationSetAttributesResponseDescriptor instead')
const PutEmailIdentityConfigurationSetAttributesResponse$json = {
  '1': 'PutEmailIdentityConfigurationSetAttributesResponse',
};

/// Descriptor for `PutEmailIdentityConfigurationSetAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityConfigurationSetAttributesResponseDescriptor =
    $convert.base64Decode(
        'CjJQdXRFbWFpbElkZW50aXR5Q29uZmlndXJhdGlvblNldEF0dHJpYnV0ZXNSZXNwb25zZQ==');

@$core.Deprecated('Use putEmailIdentityDkimAttributesRequestDescriptor instead')
const PutEmailIdentityDkimAttributesRequest$json = {
  '1': 'PutEmailIdentityDkimAttributesRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
    {
      '1': 'signingenabled',
      '3': 289566486,
      '4': 1,
      '5': 8,
      '10': 'signingenabled'
    },
  ],
};

/// Descriptor for `PutEmailIdentityDkimAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEmailIdentityDkimAttributesRequestDescriptor =
    $convert.base64Decode(
        'CiVQdXRFbWFpbElkZW50aXR5RGtpbUF0dHJpYnV0ZXNSZXF1ZXN0EicKDWVtYWlsaWRlbnRpdH'
        'kYmpTyQCABKAlSDWVtYWlsaWRlbnRpdHkSKgoOc2lnbmluZ2VuYWJsZWQYlt6JigEgASgIUg5z'
        'aWduaW5nZW5hYmxlZA==');

@$core
    .Deprecated('Use putEmailIdentityDkimAttributesResponseDescriptor instead')
const PutEmailIdentityDkimAttributesResponse$json = {
  '1': 'PutEmailIdentityDkimAttributesResponse',
};

/// Descriptor for `PutEmailIdentityDkimAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putEmailIdentityDkimAttributesResponseDescriptor =
    $convert.base64Decode(
        'CiZQdXRFbWFpbElkZW50aXR5RGtpbUF0dHJpYnV0ZXNSZXNwb25zZQ==');

@$core.Deprecated(
    'Use putEmailIdentityDkimSigningAttributesRequestDescriptor instead')
const PutEmailIdentityDkimSigningAttributesRequest$json = {
  '1': 'PutEmailIdentityDkimSigningAttributesRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
    {
      '1': 'signingattributes',
      '3': 84270406,
      '4': 1,
      '5': 11,
      '6': '.email.DkimSigningAttributes',
      '10': 'signingattributes'
    },
    {
      '1': 'signingattributesorigin',
      '3': 255219696,
      '4': 1,
      '5': 14,
      '6': '.email.DkimSigningAttributesOrigin',
      '10': 'signingattributesorigin'
    },
  ],
};

/// Descriptor for `PutEmailIdentityDkimSigningAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityDkimSigningAttributesRequestDescriptor =
    $convert.base64Decode(
        'CixQdXRFbWFpbElkZW50aXR5RGtpbVNpZ25pbmdBdHRyaWJ1dGVzUmVxdWVzdBInCg1lbWFpbG'
        'lkZW50aXR5GJqU8kAgASgJUg1lbWFpbGlkZW50aXR5Ek0KEXNpZ25pbmdhdHRyaWJ1dGVzGMa6'
        'lyggASgLMhwuZW1haWwuRGtpbVNpZ25pbmdBdHRyaWJ1dGVzUhFzaWduaW5nYXR0cmlidXRlcx'
        'JfChdzaWduaW5nYXR0cmlidXRlc29yaWdpbhjwr9l5IAEoDjIiLmVtYWlsLkRraW1TaWduaW5n'
        'QXR0cmlidXRlc09yaWdpblIXc2lnbmluZ2F0dHJpYnV0ZXNvcmlnaW4=');

@$core.Deprecated(
    'Use putEmailIdentityDkimSigningAttributesResponseDescriptor instead')
const PutEmailIdentityDkimSigningAttributesResponse$json = {
  '1': 'PutEmailIdentityDkimSigningAttributesResponse',
  '2': [
    {
      '1': 'dkimstatus',
      '3': 42822693,
      '4': 1,
      '5': 14,
      '6': '.email.DkimStatus',
      '10': 'dkimstatus'
    },
    {'1': 'dkimtokens', '3': 307118741, '4': 3, '5': 9, '10': 'dkimtokens'},
    {
      '1': 'signinghostedzone',
      '3': 442955686,
      '4': 1,
      '5': 9,
      '10': 'signinghostedzone'
    },
  ],
};

/// Descriptor for `PutEmailIdentityDkimSigningAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityDkimSigningAttributesResponseDescriptor =
    $convert.base64Decode(
        'Ci1QdXRFbWFpbElkZW50aXR5RGtpbVNpZ25pbmdBdHRyaWJ1dGVzUmVzcG9uc2USNAoKZGtpbX'
        'N0YXR1cxil2LUUIAEoDjIRLmVtYWlsLkRraW1TdGF0dXNSCmRraW1zdGF0dXMSIgoKZGtpbXRv'
        'a2VucxiVhbmSASADKAlSCmRraW10b2tlbnMSMAoRc2lnbmluZ2hvc3RlZHpvbmUYpu+b0wEgAS'
        'gJUhFzaWduaW5naG9zdGVkem9uZQ==');

@$core.Deprecated(
    'Use putEmailIdentityFeedbackAttributesRequestDescriptor instead')
const PutEmailIdentityFeedbackAttributesRequest$json = {
  '1': 'PutEmailIdentityFeedbackAttributesRequest',
  '2': [
    {
      '1': 'emailforwardingenabled',
      '3': 378369798,
      '4': 1,
      '5': 8,
      '10': 'emailforwardingenabled'
    },
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
  ],
};

/// Descriptor for `PutEmailIdentityFeedbackAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityFeedbackAttributesRequestDescriptor = $convert.base64Decode(
        'CilQdXRFbWFpbElkZW50aXR5RmVlZGJhY2tBdHRyaWJ1dGVzUmVxdWVzdBI6ChZlbWFpbGZvcn'
        'dhcmRpbmdlbmFibGVkGIbutbQBIAEoCFIWZW1haWxmb3J3YXJkaW5nZW5hYmxlZBInCg1lbWFp'
        'bGlkZW50aXR5GJqU8kAgASgJUg1lbWFpbGlkZW50aXR5');

@$core.Deprecated(
    'Use putEmailIdentityFeedbackAttributesResponseDescriptor instead')
const PutEmailIdentityFeedbackAttributesResponse$json = {
  '1': 'PutEmailIdentityFeedbackAttributesResponse',
};

/// Descriptor for `PutEmailIdentityFeedbackAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityFeedbackAttributesResponseDescriptor =
    $convert.base64Decode(
        'CipQdXRFbWFpbElkZW50aXR5RmVlZGJhY2tBdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated(
    'Use putEmailIdentityMailFromAttributesRequestDescriptor instead')
const PutEmailIdentityMailFromAttributesRequest$json = {
  '1': 'PutEmailIdentityMailFromAttributesRequest',
  '2': [
    {
      '1': 'behavioronmxfailure',
      '3': 494873128,
      '4': 1,
      '5': 14,
      '6': '.email.BehaviorOnMxFailure',
      '10': 'behavioronmxfailure'
    },
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
    {
      '1': 'mailfromdomain',
      '3': 512250671,
      '4': 1,
      '5': 9,
      '10': 'mailfromdomain'
    },
  ],
};

/// Descriptor for `PutEmailIdentityMailFromAttributesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityMailFromAttributesRequestDescriptor = $convert.base64Decode(
        'CilQdXRFbWFpbElkZW50aXR5TWFpbEZyb21BdHRyaWJ1dGVzUmVxdWVzdBJQChNiZWhhdmlvcm'
        '9ubXhmYWlsdXJlGKjU/OsBIAEoDjIaLmVtYWlsLkJlaGF2aW9yT25NeEZhaWx1cmVSE2JlaGF2'
        'aW9yb25teGZhaWx1cmUSJwoNZW1haWxpZGVudGl0eRialPJAIAEoCVINZW1haWxpZGVudGl0eR'
        'IqCg5tYWlsZnJvbWRvbWFpbhivpqH0ASABKAlSDm1haWxmcm9tZG9tYWlu');

@$core.Deprecated(
    'Use putEmailIdentityMailFromAttributesResponseDescriptor instead')
const PutEmailIdentityMailFromAttributesResponse$json = {
  '1': 'PutEmailIdentityMailFromAttributesResponse',
};

/// Descriptor for `PutEmailIdentityMailFromAttributesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    putEmailIdentityMailFromAttributesResponseDescriptor =
    $convert.base64Decode(
        'CipQdXRFbWFpbElkZW50aXR5TWFpbEZyb21BdHRyaWJ1dGVzUmVzcG9uc2U=');

@$core.Deprecated('Use putSuppressedDestinationRequestDescriptor instead')
const PutSuppressedDestinationRequest$json = {
  '1': 'PutSuppressedDestinationRequest',
  '2': [
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {
      '1': 'reason',
      '3': 20005178,
      '4': 1,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'reason'
    },
  ],
};

/// Descriptor for `PutSuppressedDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putSuppressedDestinationRequestDescriptor =
    $convert.base64Decode(
        'Ch9QdXRTdXBwcmVzc2VkRGVzdGluYXRpb25SZXF1ZXN0EiYKDGVtYWlsYWRkcmVzcxjW8YrpAS'
        'ABKAlSDGVtYWlsYWRkcmVzcxI3CgZyZWFzb24YuoLFCSABKA4yHC5lbWFpbC5TdXBwcmVzc2lv'
        'bkxpc3RSZWFzb25SBnJlYXNvbg==');

@$core.Deprecated('Use putSuppressedDestinationResponseDescriptor instead')
const PutSuppressedDestinationResponse$json = {
  '1': 'PutSuppressedDestinationResponse',
};

/// Descriptor for `PutSuppressedDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putSuppressedDestinationResponseDescriptor =
    $convert.base64Decode('CiBQdXRTdXBwcmVzc2VkRGVzdGluYXRpb25SZXNwb25zZQ==');

@$core.Deprecated('Use rawMessageDescriptor instead')
const RawMessage$json = {
  '1': 'RawMessage',
  '2': [
    {'1': 'data', '3': 525498822, '4': 1, '5': 12, '10': 'data'},
  ],
};

/// Descriptor for `RawMessage`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rawMessageDescriptor =
    $convert.base64Decode('CgpSYXdNZXNzYWdlEhYKBGRhdGEYxvPJ+gEgASgMUgRkYXRh');

@$core.Deprecated('Use recommendationDescriptor instead')
const Recommendation$json = {
  '1': 'Recommendation',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'impact',
      '3': 75836600,
      '4': 1,
      '5': 14,
      '6': '.email.RecommendationImpact',
      '10': 'impact'
    },
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.RecommendationStatus',
      '10': 'status'
    },
    {
      '1': 'type',
      '3': 290836590,
      '4': 1,
      '5': 14,
      '6': '.email.RecommendationType',
      '10': 'type'
    },
  ],
};

/// Descriptor for `Recommendation`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recommendationDescriptor = $convert.base64Decode(
    'Cg5SZWNvbW1lbmRhdGlvbhIuChBjcmVhdGVkdGltZXN0YW1wGPrbz58BIAEoCVIQY3JlYXRlZH'
    'RpbWVzdGFtcBIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SNgoGaW1wYWN0'
    'GLjZlCQgASgOMhsuZW1haWwuUmVjb21tZW5kYXRpb25JbXBhY3RSBmltcGFjdBI1ChRsYXN0dX'
    'BkYXRlZHRpbWVzdGFtcBiVy8g/IAEoCVIUbGFzdHVwZGF0ZWR0aW1lc3RhbXASJAoLcmVzb3Vy'
    'Y2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJjZWFybhI2CgZzdGF0dXMYkOT7AiABKA4yGy5lbWFpbC'
    '5SZWNvbW1lbmRhdGlvblN0YXR1c1IGc3RhdHVzEjEKBHR5cGUY7qDXigEgASgOMhkuZW1haWwu'
    'UmVjb21tZW5kYXRpb25UeXBlUgR0eXBl');

@$core.Deprecated('Use replacementEmailContentDescriptor instead')
const ReplacementEmailContent$json = {
  '1': 'ReplacementEmailContent',
  '2': [
    {
      '1': 'replacementtemplate',
      '3': 265092678,
      '4': 1,
      '5': 11,
      '6': '.email.ReplacementTemplate',
      '10': 'replacementtemplate'
    },
  ],
};

/// Descriptor for `ReplacementEmailContent`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replacementEmailContentDescriptor = $convert.base64Decode(
    'ChdSZXBsYWNlbWVudEVtYWlsQ29udGVudBJPChNyZXBsYWNlbWVudHRlbXBsYXRlGMb8s34gAS'
    'gLMhouZW1haWwuUmVwbGFjZW1lbnRUZW1wbGF0ZVITcmVwbGFjZW1lbnR0ZW1wbGF0ZQ==');

@$core.Deprecated('Use replacementTemplateDescriptor instead')
const ReplacementTemplate$json = {
  '1': 'ReplacementTemplate',
  '2': [
    {
      '1': 'replacementtemplatedata',
      '3': 382231110,
      '4': 1,
      '5': 9,
      '10': 'replacementtemplatedata'
    },
  ],
};

/// Descriptor for `ReplacementTemplate`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replacementTemplateDescriptor = $convert.base64Decode(
    'ChNSZXBsYWNlbWVudFRlbXBsYXRlEjwKF3JlcGxhY2VtZW50dGVtcGxhdGVkYXRhGMbEobYBIA'
    'EoCVIXcmVwbGFjZW1lbnR0ZW1wbGF0ZWRhdGE=');

@$core.Deprecated('Use reputationEntityDescriptor instead')
const ReputationEntity$json = {
  '1': 'ReputationEntity',
  '2': [
    {
      '1': 'awssesmanagedstatus',
      '3': 20222799,
      '4': 1,
      '5': 11,
      '6': '.email.StatusRecord',
      '10': 'awssesmanagedstatus'
    },
    {
      '1': 'customermanagedstatus',
      '3': 21835973,
      '4': 1,
      '5': 11,
      '6': '.email.StatusRecord',
      '10': 'customermanagedstatus'
    },
    {
      '1': 'reputationentityreference',
      '3': 414929111,
      '4': 1,
      '5': 9,
      '10': 'reputationentityreference'
    },
    {
      '1': 'reputationentitytype',
      '3': 98287826,
      '4': 1,
      '5': 14,
      '6': '.email.ReputationEntityType',
      '10': 'reputationentitytype'
    },
    {
      '1': 'reputationimpact',
      '3': 345570753,
      '4': 1,
      '5': 14,
      '6': '.email.RecommendationImpact',
      '10': 'reputationimpact'
    },
    {
      '1': 'reputationmanagementpolicy',
      '3': 202871484,
      '4': 1,
      '5': 9,
      '10': 'reputationmanagementpolicy'
    },
    {
      '1': 'sendingstatusaggregate',
      '3': 10196739,
      '4': 1,
      '5': 14,
      '6': '.email.SendingStatus',
      '10': 'sendingstatusaggregate'
    },
  ],
};

/// Descriptor for `ReputationEntity`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reputationEntityDescriptor = $convert.base64Decode(
    'ChBSZXB1dGF0aW9uRW50aXR5EkgKE2F3c3Nlc21hbmFnZWRzdGF0dXMYz6bSCSABKAsyEy5lbW'
    'FpbC5TdGF0dXNSZWNvcmRSE2F3c3Nlc21hbmFnZWRzdGF0dXMSTAoVY3VzdG9tZXJtYW5hZ2Vk'
    'c3RhdHVzGMXhtAogASgLMhMuZW1haWwuU3RhdHVzUmVjb3JkUhVjdXN0b21lcm1hbmFnZWRzdG'
    'F0dXMSQAoZcmVwdXRhdGlvbmVudGl0eXJlZmVyZW5jZRjXoe3FASABKAlSGXJlcHV0YXRpb25l'
    'bnRpdHlyZWZlcmVuY2USUgoUcmVwdXRhdGlvbmVudGl0eXR5cGUY0oHvLiABKA4yGy5lbWFpbC'
    '5SZXB1dGF0aW9uRW50aXR5VHlwZVIUcmVwdXRhdGlvbmVudGl0eXR5cGUSSwoQcmVwdXRhdGlv'
    'bmltcGFjdBjB++OkASABKA4yGy5lbWFpbC5SZWNvbW1lbmRhdGlvbkltcGFjdFIQcmVwdXRhdG'
    'lvbmltcGFjdBJBChpyZXB1dGF0aW9ubWFuYWdlbWVudHBvbGljeRi8pd5gIAEoCVIacmVwdXRh'
    'dGlvbm1hbmFnZW1lbnRwb2xpY3kSTwoWc2VuZGluZ3N0YXR1c2FnZ3JlZ2F0ZRiDru4EIAEoDj'
    'IULmVtYWlsLlNlbmRpbmdTdGF0dXNSFnNlbmRpbmdzdGF0dXNhZ2dyZWdhdGU=');

@$core.Deprecated('Use reputationOptionsDescriptor instead')
const ReputationOptions$json = {
  '1': 'ReputationOptions',
  '2': [
    {
      '1': 'lastfreshstart',
      '3': 156685530,
      '4': 1,
      '5': 9,
      '10': 'lastfreshstart'
    },
    {
      '1': 'reputationmetricsenabled',
      '3': 440357917,
      '4': 1,
      '5': 8,
      '10': 'reputationmetricsenabled'
    },
  ],
};

/// Descriptor for `ReputationOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reputationOptionsDescriptor = $convert.base64Decode(
    'ChFSZXB1dGF0aW9uT3B0aW9ucxIpCg5sYXN0ZnJlc2hzdGFydBjaqdtKIAEoCVIObGFzdGZyZX'
    'Noc3RhcnQSPgoYcmVwdXRhdGlvbm1ldHJpY3NlbmFibGVkGJ2o/dEBIAEoCFIYcmVwdXRhdGlv'
    'bm1ldHJpY3NlbmFibGVk');

@$core.Deprecated('Use resourceTenantMetadataDescriptor instead')
const ResourceTenantMetadata$json = {
  '1': 'ResourceTenantMetadata',
  '2': [
    {
      '1': 'associatedtimestamp',
      '3': 169517010,
      '4': 1,
      '5': 9,
      '10': 'associatedtimestamp'
    },
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `ResourceTenantMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceTenantMetadataDescriptor = $convert.base64Decode(
    'ChZSZXNvdXJjZVRlbmFudE1ldGFkYXRhEjMKE2Fzc29jaWF0ZWR0aW1lc3RhbXAY0r/qUCABKA'
    'lSE2Fzc29jaWF0ZWR0aW1lc3RhbXASJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJj'
    'ZWFybhIdCgh0ZW5hbnRpZBiB+O9kIAEoCVIIdGVuYW50aWQSIQoKdGVuYW50bmFtZRiH3NNSIA'
    'EoCVIKdGVuYW50bmFtZQ==');

@$core.Deprecated('Use reviewDetailsDescriptor instead')
const ReviewDetails$json = {
  '1': 'ReviewDetails',
  '2': [
    {'1': 'caseid', '3': 181380677, '4': 1, '5': 9, '10': 'caseid'},
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.ReviewStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `ReviewDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reviewDetailsDescriptor = $convert.base64Decode(
    'Cg1SZXZpZXdEZXRhaWxzEhkKBmNhc2VpZBjFzL5WIAEoCVIGY2FzZWlkEi4KBnN0YXR1cxiQ5P'
    'sCIAEoDjITLmVtYWlsLlJldmlld1N0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use routeDescriptor instead')
const Route$json = {
  '1': 'Route',
  '2': [
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
  ],
};

/// Descriptor for `Route`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List routeDescriptor =
    $convert.base64Decode('CgVSb3V0ZRIZCgZyZWdpb24YnvG5SSABKAlSBnJlZ2lvbg==');

@$core.Deprecated('Use routeDetailsDescriptor instead')
const RouteDetails$json = {
  '1': 'RouteDetails',
  '2': [
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
  ],
};

/// Descriptor for `RouteDetails`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List routeDetailsDescriptor = $convert
    .base64Decode('CgxSb3V0ZURldGFpbHMSGQoGcmVnaW9uGJ7xuUkgASgJUgZyZWdpb24=');

@$core.Deprecated('Use sOARecordDescriptor instead')
const SOARecord$json = {
  '1': 'SOARecord',
  '2': [
    {'1': 'adminemail', '3': 289730349, '4': 1, '5': 9, '10': 'adminemail'},
    {
      '1': 'primarynameserver',
      '3': 522614982,
      '4': 1,
      '5': 9,
      '10': 'primarynameserver'
    },
    {'1': 'serialnumber', '3': 418274661, '4': 1, '5': 3, '10': 'serialnumber'},
  ],
};

/// Descriptor for `SOARecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sOARecordDescriptor = $convert.base64Decode(
    'CglTT0FSZWNvcmQSIgoKYWRtaW5lbWFpbBit3pOKASABKAlSCmFkbWluZW1haWwSMAoRcHJpbW'
    'FyeW5hbWVzZXJ2ZXIYxvGZ+QEgASgJUhFwcmltYXJ5bmFtZXNlcnZlchImCgxzZXJpYWxudW1i'
    'ZXIY5bq5xwEgASgDUgxzZXJpYWxudW1iZXI=');

@$core.Deprecated('Use sendBulkEmailRequestDescriptor instead')
const SendBulkEmailRequest$json = {
  '1': 'SendBulkEmailRequest',
  '2': [
    {
      '1': 'bulkemailentries',
      '3': 509093366,
      '4': 3,
      '5': 11,
      '6': '.email.BulkEmailEntry',
      '10': 'bulkemailentries'
    },
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'defaultcontent',
      '3': 178026400,
      '4': 1,
      '5': 11,
      '6': '.email.BulkEmailContent',
      '10': 'defaultcontent'
    },
    {
      '1': 'defaultemailtags',
      '3': 272909972,
      '4': 3,
      '5': 11,
      '6': '.email.MessageTag',
      '10': 'defaultemailtags'
    },
    {'1': 'endpointid', '3': 35808946, '4': 1, '5': 9, '10': 'endpointid'},
    {
      '1': 'feedbackforwardingemailaddress',
      '3': 237438826,
      '4': 1,
      '5': 9,
      '10': 'feedbackforwardingemailaddress'
    },
    {
      '1': 'feedbackforwardingemailaddressidentityarn',
      '3': 5181499,
      '4': 1,
      '5': 9,
      '10': 'feedbackforwardingemailaddressidentityarn'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {
      '1': 'fromemailaddressidentityarn',
      '3': 171860335,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddressidentityarn'
    },
    {
      '1': 'replytoaddresses',
      '3': 498368923,
      '4': 3,
      '5': 9,
      '10': 'replytoaddresses'
    },
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `SendBulkEmailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendBulkEmailRequestDescriptor = $convert.base64Decode(
    'ChRTZW5kQnVsa0VtYWlsUmVxdWVzdBJFChBidWxrZW1haWxlbnRyaWVzGPbL4PIBIAMoCzIVLm'
    'VtYWlsLkJ1bGtFbWFpbEVudHJ5UhBidWxrZW1haWxlbnRyaWVzEjYKFGNvbmZpZ3VyYXRpb25z'
    'ZXRuYW1lGM2LscABIAEoCVIUY29uZmlndXJhdGlvbnNldG5hbWUSQgoOZGVmYXVsdGNvbnRlbn'
    'QYoO/xVCABKAsyFy5lbWFpbC5CdWxrRW1haWxDb250ZW50Ug5kZWZhdWx0Y29udGVudBJBChBk'
    'ZWZhdWx0ZW1haWx0YWdzGJSNkYIBIAMoCzIRLmVtYWlsLk1lc3NhZ2VUYWdSEGRlZmF1bHRlbW'
    'FpbHRhZ3MSIQoKZW5kcG9pbnRpZBiyzYkRIAEoCVIKZW5kcG9pbnRpZBJJCh5mZWVkYmFja2Zv'
    'cndhcmRpbmdlbWFpbGFkZHJlc3MY6o6ccSABKAlSHmZlZWRiYWNrZm9yd2FyZGluZ2VtYWlsYW'
    'RkcmVzcxJfCilmZWVkYmFja2ZvcndhcmRpbmdlbWFpbGFkZHJlc3NpZGVudGl0eWFybhi7oLwC'
    'IAEoCVIpZmVlZGJhY2tmb3J3YXJkaW5nZW1haWxhZGRyZXNzaWRlbnRpdHlhcm4SLQoQZnJvbW'
    'VtYWlsYWRkcmVzcxiGmsssIAEoCVIQZnJvbWVtYWlsYWRkcmVzcxJDChtmcm9tZW1haWxhZGRy'
    'ZXNzaWRlbnRpdHlhcm4Y78L5USABKAlSG2Zyb21lbWFpbGFkZHJlc3NpZGVudGl0eWFybhIuCh'
    'ByZXBseXRvYWRkcmVzc2VzGJuD0u0BIAMoCVIQcmVwbHl0b2FkZHJlc3NlcxIhCgp0ZW5hbnRu'
    'YW1lGIfc01IgASgJUgp0ZW5hbnRuYW1l');

@$core.Deprecated('Use sendBulkEmailResponseDescriptor instead')
const SendBulkEmailResponse$json = {
  '1': 'SendBulkEmailResponse',
  '2': [
    {
      '1': 'bulkemailentryresults',
      '3': 140486386,
      '4': 3,
      '5': 11,
      '6': '.email.BulkEmailEntryResult',
      '10': 'bulkemailentryresults'
    },
  ],
};

/// Descriptor for `SendBulkEmailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendBulkEmailResponseDescriptor = $convert.base64Decode(
    'ChVTZW5kQnVsa0VtYWlsUmVzcG9uc2USVAoVYnVsa2VtYWlsZW50cnlyZXN1bHRzGPLN/kIgAy'
    'gLMhsuZW1haWwuQnVsa0VtYWlsRW50cnlSZXN1bHRSFWJ1bGtlbWFpbGVudHJ5cmVzdWx0cw==');

@$core.Deprecated('Use sendCustomVerificationEmailRequestDescriptor instead')
const SendCustomVerificationEmailRequest$json = {
  '1': 'SendCustomVerificationEmailRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `SendCustomVerificationEmailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendCustomVerificationEmailRequestDescriptor =
    $convert.base64Decode(
        'CiJTZW5kQ3VzdG9tVmVyaWZpY2F0aW9uRW1haWxSZXF1ZXN0EjYKFGNvbmZpZ3VyYXRpb25zZX'
        'RuYW1lGM2LscABIAEoCVIUY29uZmlndXJhdGlvbnNldG5hbWUSJgoMZW1haWxhZGRyZXNzGNbx'
        'iukBIAEoCVIMZW1haWxhZGRyZXNzEiYKDHRlbXBsYXRlbmFtZRixmJHlASABKAlSDHRlbXBsYX'
        'RlbmFtZQ==');

@$core.Deprecated('Use sendCustomVerificationEmailResponseDescriptor instead')
const SendCustomVerificationEmailResponse$json = {
  '1': 'SendCustomVerificationEmailResponse',
  '2': [
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
  ],
};

/// Descriptor for `SendCustomVerificationEmailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendCustomVerificationEmailResponseDescriptor =
    $convert.base64Decode(
        'CiNTZW5kQ3VzdG9tVmVyaWZpY2F0aW9uRW1haWxSZXNwb25zZRIgCgltZXNzYWdlaWQYqub0qw'
        'EgASgJUgltZXNzYWdlaWQ=');

@$core.Deprecated('Use sendEmailRequestDescriptor instead')
const SendEmailRequest$json = {
  '1': 'SendEmailRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'content',
      '3': 23568227,
      '4': 1,
      '5': 11,
      '6': '.email.EmailContent',
      '10': 'content'
    },
    {
      '1': 'destination',
      '3': 457443680,
      '4': 1,
      '5': 11,
      '6': '.email.Destination',
      '10': 'destination'
    },
    {
      '1': 'emailtags',
      '3': 215809691,
      '4': 3,
      '5': 11,
      '6': '.email.MessageTag',
      '10': 'emailtags'
    },
    {'1': 'endpointid', '3': 35808946, '4': 1, '5': 9, '10': 'endpointid'},
    {
      '1': 'feedbackforwardingemailaddress',
      '3': 237438826,
      '4': 1,
      '5': 9,
      '10': 'feedbackforwardingemailaddress'
    },
    {
      '1': 'feedbackforwardingemailaddressidentityarn',
      '3': 5181499,
      '4': 1,
      '5': 9,
      '10': 'feedbackforwardingemailaddressidentityarn'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {
      '1': 'fromemailaddressidentityarn',
      '3': 171860335,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddressidentityarn'
    },
    {
      '1': 'listmanagementoptions',
      '3': 419995741,
      '4': 1,
      '5': 11,
      '6': '.email.ListManagementOptions',
      '10': 'listmanagementoptions'
    },
    {
      '1': 'replytoaddresses',
      '3': 498368923,
      '4': 3,
      '5': 9,
      '10': 'replytoaddresses'
    },
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `SendEmailRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendEmailRequestDescriptor = $convert.base64Decode(
    'ChBTZW5kRW1haWxSZXF1ZXN0EjYKFGNvbmZpZ3VyYXRpb25zZXRuYW1lGM2LscABIAEoCVIUY2'
    '9uZmlndXJhdGlvbnNldG5hbWUSMAoHY29udGVudBjjvp4LIAEoCzITLmVtYWlsLkVtYWlsQ29u'
    'dGVudFIHY29udGVudBI4CgtkZXN0aW5hdGlvbhjgkpDaASABKAsyEi5lbWFpbC5EZXN0aW5hdG'
    'lvblILZGVzdGluYXRpb24SMgoJZW1haWx0YWdzGJv982YgAygLMhEuZW1haWwuTWVzc2FnZVRh'
    'Z1IJZW1haWx0YWdzEiEKCmVuZHBvaW50aWQYss2JESABKAlSCmVuZHBvaW50aWQSSQoeZmVlZG'
    'JhY2tmb3J3YXJkaW5nZW1haWxhZGRyZXNzGOqOnHEgASgJUh5mZWVkYmFja2ZvcndhcmRpbmdl'
    'bWFpbGFkZHJlc3MSXwopZmVlZGJhY2tmb3J3YXJkaW5nZW1haWxhZGRyZXNzaWRlbnRpdHlhcm'
    '4Yu6C8AiABKAlSKWZlZWRiYWNrZm9yd2FyZGluZ2VtYWlsYWRkcmVzc2lkZW50aXR5YXJuEi0K'
    'EGZyb21lbWFpbGFkZHJlc3MYhprLLCABKAlSEGZyb21lbWFpbGFkZHJlc3MSQwobZnJvbWVtYW'
    'lsYWRkcmVzc2lkZW50aXR5YXJuGO/C+VEgASgJUhtmcm9tZW1haWxhZGRyZXNzaWRlbnRpdHlh'
    'cm4SVgoVbGlzdG1hbmFnZW1lbnRvcHRpb25zGN3AosgBIAEoCzIcLmVtYWlsLkxpc3RNYW5hZ2'
    'VtZW50T3B0aW9uc1IVbGlzdG1hbmFnZW1lbnRvcHRpb25zEi4KEHJlcGx5dG9hZGRyZXNzZXMY'
    'm4PS7QEgAygJUhByZXBseXRvYWRkcmVzc2VzEiEKCnRlbmFudG5hbWUYh9zTUiABKAlSCnRlbm'
    'FudG5hbWU=');

@$core.Deprecated('Use sendEmailResponseDescriptor instead')
const SendEmailResponse$json = {
  '1': 'SendEmailResponse',
  '2': [
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
  ],
};

/// Descriptor for `SendEmailResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendEmailResponseDescriptor = $convert.base64Decode(
    'ChFTZW5kRW1haWxSZXNwb25zZRIgCgltZXNzYWdlaWQYqub0qwEgASgJUgltZXNzYWdlaWQ=');

@$core.Deprecated('Use sendQuotaDescriptor instead')
const SendQuota$json = {
  '1': 'SendQuota',
  '2': [
    {
      '1': 'max24hoursend',
      '3': 390803852,
      '4': 1,
      '5': 1,
      '10': 'max24hoursend'
    },
    {'1': 'maxsendrate', '3': 331016098, '4': 1, '5': 1, '10': 'maxsendrate'},
    {
      '1': 'sentlast24hours',
      '3': 272279377,
      '4': 1,
      '5': 1,
      '10': 'sentlast24hours'
    },
  ],
};

/// Descriptor for `SendQuota`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendQuotaDescriptor = $convert.base64Decode(
    'CglTZW5kUXVvdGESKAoNbWF4MjRob3Vyc2VuZBiM46y6ASABKAFSDW1heDI0aG91cnNlbmQSJA'
    'oLbWF4c2VuZHJhdGUYos/rnQEgASgBUgttYXhzZW5kcmF0ZRIsCg9zZW50bGFzdDI0aG91cnMY'
    '0c7qgQEgASgBUg9zZW50bGFzdDI0aG91cnM=');

@$core.Deprecated('Use sendingOptionsDescriptor instead')
const SendingOptions$json = {
  '1': 'SendingOptions',
  '2': [
    {
      '1': 'sendingenabled',
      '3': 194846115,
      '4': 1,
      '5': 8,
      '10': 'sendingenabled'
    },
  ],
};

/// Descriptor for `SendingOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendingOptionsDescriptor = $convert.base64Decode(
    'Cg5TZW5kaW5nT3B0aW9ucxIpCg5zZW5kaW5nZW5hYmxlZBiju/RcIAEoCFIOc2VuZGluZ2VuYW'
    'JsZWQ=');

@$core.Deprecated('Use sendingPausedExceptionDescriptor instead')
const SendingPausedException$json = {
  '1': 'SendingPausedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `SendingPausedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sendingPausedExceptionDescriptor =
    $convert.base64Decode(
        'ChZTZW5kaW5nUGF1c2VkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use snsDestinationDescriptor instead')
const SnsDestination$json = {
  '1': 'SnsDestination',
  '2': [
    {'1': 'topicarn', '3': 30652956, '4': 1, '5': 9, '10': 'topicarn'},
  ],
};

/// Descriptor for `SnsDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List snsDestinationDescriptor = $convert.base64Decode(
    'Cg5TbnNEZXN0aW5hdGlvbhIdCgh0b3BpY2Fybhic9M4OIAEoCVIIdG9waWNhcm4=');

@$core.Deprecated('Use statusRecordDescriptor instead')
const StatusRecord$json = {
  '1': 'StatusRecord',
  '2': [
    {'1': 'cause', '3': 283421889, '4': 1, '5': 9, '10': 'cause'},
    {
      '1': 'lastupdatedtimestamp',
      '3': 133309845,
      '4': 1,
      '5': 9,
      '10': 'lastupdatedtimestamp'
    },
    {
      '1': 'status',
      '3': 6222352,
      '4': 1,
      '5': 14,
      '6': '.email.SendingStatus',
      '10': 'status'
    },
  ],
};

/// Descriptor for `StatusRecord`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List statusRecordDescriptor = $convert.base64Decode(
    'CgxTdGF0dXNSZWNvcmQSGAoFY2F1c2UYwdmShwEgASgJUgVjYXVzZRI1ChRsYXN0dXBkYXRlZH'
    'RpbWVzdGFtcBiVy8g/IAEoCVIUbGFzdHVwZGF0ZWR0aW1lc3RhbXASLwoGc3RhdHVzGJDk+wIg'
    'ASgOMhQuZW1haWwuU2VuZGluZ1N0YXR1c1IGc3RhdHVz');

@$core.Deprecated('Use suppressedDestinationDescriptor instead')
const SuppressedDestination$json = {
  '1': 'SuppressedDestination',
  '2': [
    {
      '1': 'attributes',
      '3': 209638581,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressedDestinationAttributes',
      '10': 'attributes'
    },
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {
      '1': 'lastupdatetime',
      '3': 389145074,
      '4': 1,
      '5': 9,
      '10': 'lastupdatetime'
    },
    {
      '1': 'reason',
      '3': 20005178,
      '4': 1,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'reason'
    },
  ],
};

/// Descriptor for `SuppressedDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressedDestinationDescriptor = $convert.base64Decode(
    'ChVTdXBwcmVzc2VkRGVzdGluYXRpb24SSQoKYXR0cmlidXRlcxi1qftjIAEoCzImLmVtYWlsLl'
    'N1cHByZXNzZWREZXN0aW5hdGlvbkF0dHJpYnV0ZXNSCmF0dHJpYnV0ZXMSJgoMZW1haWxhZGRy'
    'ZXNzGNbxiukBIAEoCVIMZW1haWxhZGRyZXNzEioKDmxhc3R1cGRhdGV0aW1lGPLDx7kBIAEoCV'
    'IObGFzdHVwZGF0ZXRpbWUSNwoGcmVhc29uGLqCxQkgASgOMhwuZW1haWwuU3VwcHJlc3Npb25M'
    'aXN0UmVhc29uUgZyZWFzb24=');

@$core.Deprecated('Use suppressedDestinationAttributesDescriptor instead')
const SuppressedDestinationAttributes$json = {
  '1': 'SuppressedDestinationAttributes',
  '2': [
    {'1': 'feedbackid', '3': 291096336, '4': 1, '5': 9, '10': 'feedbackid'},
    {'1': 'messageid', '3': 360526634, '4': 1, '5': 9, '10': 'messageid'},
  ],
};

/// Descriptor for `SuppressedDestinationAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressedDestinationAttributesDescriptor =
    $convert.base64Decode(
        'Ch9TdXBwcmVzc2VkRGVzdGluYXRpb25BdHRyaWJ1dGVzEiIKCmZlZWRiYWNraWQYkI7nigEgAS'
        'gJUgpmZWVkYmFja2lkEiAKCW1lc3NhZ2VpZBiq5vSrASABKAlSCW1lc3NhZ2VpZA==');

@$core.Deprecated('Use suppressedDestinationSummaryDescriptor instead')
const SuppressedDestinationSummary$json = {
  '1': 'SuppressedDestinationSummary',
  '2': [
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {
      '1': 'lastupdatetime',
      '3': 389145074,
      '4': 1,
      '5': 9,
      '10': 'lastupdatetime'
    },
    {
      '1': 'reason',
      '3': 20005178,
      '4': 1,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'reason'
    },
  ],
};

/// Descriptor for `SuppressedDestinationSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressedDestinationSummaryDescriptor = $convert.base64Decode(
    'ChxTdXBwcmVzc2VkRGVzdGluYXRpb25TdW1tYXJ5EiYKDGVtYWlsYWRkcmVzcxjW8YrpASABKA'
    'lSDGVtYWlsYWRkcmVzcxIqCg5sYXN0dXBkYXRldGltZRjyw8e5ASABKAlSDmxhc3R1cGRhdGV0'
    'aW1lEjcKBnJlYXNvbhi6gsUJIAEoDjIcLmVtYWlsLlN1cHByZXNzaW9uTGlzdFJlYXNvblIGcm'
    'Vhc29u');

@$core.Deprecated('Use suppressionAttributesDescriptor instead')
const SuppressionAttributes$json = {
  '1': 'SuppressionAttributes',
  '2': [
    {
      '1': 'suppressedreasons',
      '3': 465922417,
      '4': 3,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'suppressedreasons'
    },
    {
      '1': 'validationattributes',
      '3': 171836944,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionValidationAttributes',
      '10': 'validationattributes'
    },
  ],
};

/// Descriptor for `SuppressionAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionAttributesDescriptor = $convert.base64Decode(
    'ChVTdXBwcmVzc2lvbkF0dHJpYnV0ZXMSTgoRc3VwcHJlc3NlZHJlYXNvbnMY8dKV3gEgAygOMh'
    'wuZW1haWwuU3VwcHJlc3Npb25MaXN0UmVhc29uUhFzdXBwcmVzc2VkcmVhc29ucxJdChR2YWxp'
    'ZGF0aW9uYXR0cmlidXRlcxiQjPhRIAEoCzImLmVtYWlsLlN1cHByZXNzaW9uVmFsaWRhdGlvbk'
    'F0dHJpYnV0ZXNSFHZhbGlkYXRpb25hdHRyaWJ1dGVz');

@$core.Deprecated('Use suppressionConditionThresholdDescriptor instead')
const SuppressionConditionThreshold$json = {
  '1': 'SuppressionConditionThreshold',
  '2': [
    {
      '1': 'conditionthresholdenabled',
      '3': 527217635,
      '4': 1,
      '5': 14,
      '6': '.email.FeatureStatus',
      '10': 'conditionthresholdenabled'
    },
    {
      '1': 'overallconfidencethreshold',
      '3': 440641390,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionConfidenceThreshold',
      '10': 'overallconfidencethreshold'
    },
  ],
};

/// Descriptor for `SuppressionConditionThreshold`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionConditionThresholdDescriptor = $convert.base64Decode(
    'Ch1TdXBwcmVzc2lvbkNvbmRpdGlvblRocmVzaG9sZBJWChljb25kaXRpb250aHJlc2hvbGRlbm'
    'FibGVkGOPnsvsBIAEoDjIULmVtYWlsLkZlYXR1cmVTdGF0dXNSGWNvbmRpdGlvbnRocmVzaG9s'
    'ZGVuYWJsZWQSaQoab3ZlcmFsbGNvbmZpZGVuY2V0aHJlc2hvbGQY7s6O0gEgASgLMiUuZW1haW'
    'wuU3VwcHJlc3Npb25Db25maWRlbmNlVGhyZXNob2xkUhpvdmVyYWxsY29uZmlkZW5jZXRocmVz'
    'aG9sZA==');

@$core.Deprecated('Use suppressionConfidenceThresholdDescriptor instead')
const SuppressionConfidenceThreshold$json = {
  '1': 'SuppressionConfidenceThreshold',
  '2': [
    {
      '1': 'confidenceverdictthreshold',
      '3': 259879446,
      '4': 1,
      '5': 14,
      '6': '.email.SuppressionConfidenceVerdictThreshold',
      '10': 'confidenceverdictthreshold'
    },
  ],
};

/// Descriptor for `SuppressionConfidenceThreshold`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionConfidenceThresholdDescriptor =
    $convert.base64Decode(
        'Ch5TdXBwcmVzc2lvbkNvbmZpZGVuY2VUaHJlc2hvbGQSbwoaY29uZmlkZW5jZXZlcmRpY3R0aH'
        'Jlc2hvbGQYluT1eyABKA4yLC5lbWFpbC5TdXBwcmVzc2lvbkNvbmZpZGVuY2VWZXJkaWN0VGhy'
        'ZXNob2xkUhpjb25maWRlbmNldmVyZGljdHRocmVzaG9sZA==');

@$core.Deprecated('Use suppressionListDestinationDescriptor instead')
const SuppressionListDestination$json = {
  '1': 'SuppressionListDestination',
  '2': [
    {
      '1': 'suppressionlistimportaction',
      '3': 402316168,
      '4': 1,
      '5': 14,
      '6': '.email.SuppressionListImportAction',
      '10': 'suppressionlistimportaction'
    },
  ],
};

/// Descriptor for `SuppressionListDestination`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionListDestinationDescriptor =
    $convert.base64Decode(
        'ChpTdXBwcmVzc2lvbkxpc3REZXN0aW5hdGlvbhJoChtzdXBwcmVzc2lvbmxpc3RpbXBvcnRhY3'
        'Rpb24YiLfrvwEgASgOMiIuZW1haWwuU3VwcHJlc3Npb25MaXN0SW1wb3J0QWN0aW9uUhtzdXBw'
        'cmVzc2lvbmxpc3RpbXBvcnRhY3Rpb24=');

@$core.Deprecated('Use suppressionOptionsDescriptor instead')
const SuppressionOptions$json = {
  '1': 'SuppressionOptions',
  '2': [
    {
      '1': 'suppressedreasons',
      '3': 465922417,
      '4': 3,
      '5': 14,
      '6': '.email.SuppressionListReason',
      '10': 'suppressedreasons'
    },
    {
      '1': 'validationoptions',
      '3': 216495637,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionValidationOptions',
      '10': 'validationoptions'
    },
  ],
};

/// Descriptor for `SuppressionOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionOptionsDescriptor = $convert.base64Decode(
    'ChJTdXBwcmVzc2lvbk9wdGlvbnMSTgoRc3VwcHJlc3NlZHJlYXNvbnMY8dKV3gEgAygOMhwuZW'
    '1haWwuU3VwcHJlc3Npb25MaXN0UmVhc29uUhFzdXBwcmVzc2VkcmVhc29ucxJUChF2YWxpZGF0'
    'aW9ub3B0aW9ucxiV7J1nIAEoCzIjLmVtYWlsLlN1cHByZXNzaW9uVmFsaWRhdGlvbk9wdGlvbn'
    'NSEXZhbGlkYXRpb25vcHRpb25z');

@$core.Deprecated('Use suppressionValidationAttributesDescriptor instead')
const SuppressionValidationAttributes$json = {
  '1': 'SuppressionValidationAttributes',
  '2': [
    {
      '1': 'conditionthreshold',
      '3': 528759594,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionConditionThreshold',
      '10': 'conditionthreshold'
    },
  ],
};

/// Descriptor for `SuppressionValidationAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionValidationAttributesDescriptor =
    $convert.base64Decode(
        'Ch9TdXBwcmVzc2lvblZhbGlkYXRpb25BdHRyaWJ1dGVzElgKEmNvbmRpdGlvbnRocmVzaG9sZB'
        'iq9pD8ASABKAsyJC5lbWFpbC5TdXBwcmVzc2lvbkNvbmRpdGlvblRocmVzaG9sZFISY29uZGl0'
        'aW9udGhyZXNob2xk');

@$core.Deprecated('Use suppressionValidationOptionsDescriptor instead')
const SuppressionValidationOptions$json = {
  '1': 'SuppressionValidationOptions',
  '2': [
    {
      '1': 'conditionthreshold',
      '3': 528759594,
      '4': 1,
      '5': 11,
      '6': '.email.SuppressionConditionThreshold',
      '10': 'conditionthreshold'
    },
  ],
};

/// Descriptor for `SuppressionValidationOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List suppressionValidationOptionsDescriptor =
    $convert.base64Decode(
        'ChxTdXBwcmVzc2lvblZhbGlkYXRpb25PcHRpb25zElgKEmNvbmRpdGlvbnRocmVzaG9sZBiq9p'
        'D8ASABKAsyJC5lbWFpbC5TdXBwcmVzc2lvbkNvbmRpdGlvblRocmVzaG9sZFISY29uZGl0aW9u'
        'dGhyZXNob2xk');

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
      '6': '.email.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSJAoLcmVzb3VyY2Vhcm4YrfjZrQEgASgJUgtyZXNvdXJjZW'
    'FybhIiCgR0YWdzGMHB9rUBIAMoCzIKLmVtYWlsLlRhZ1IEdGFncw==');

@$core.Deprecated('Use tagResourceResponseDescriptor instead')
const TagResourceResponse$json = {
  '1': 'TagResourceResponse',
};

/// Descriptor for `TagResourceResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceResponseDescriptor =
    $convert.base64Decode('ChNUYWdSZXNvdXJjZVJlc3BvbnNl');

@$core.Deprecated('Use templateDescriptor instead')
const Template$json = {
  '1': 'Template',
  '2': [
    {
      '1': 'attachments',
      '3': 498946338,
      '4': 3,
      '5': 11,
      '6': '.email.Attachment',
      '10': 'attachments'
    },
    {
      '1': 'headers',
      '3': 323967370,
      '4': 3,
      '5': 11,
      '6': '.email.MessageHeader',
      '10': 'headers'
    },
    {'1': 'templatearn', '3': 384599199, '4': 1, '5': 9, '10': 'templatearn'},
    {
      '1': 'templatecontent',
      '3': 528973453,
      '4': 1,
      '5': 11,
      '6': '.email.EmailTemplateContent',
      '10': 'templatecontent'
    },
    {'1': 'templatedata', '3': 284145300, '4': 1, '5': 9, '10': 'templatedata'},
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `Template`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List templateDescriptor = $convert.base64Decode(
    'CghUZW1wbGF0ZRI3CgthdHRhY2htZW50cxiiovXtASADKAsyES5lbWFpbC5BdHRhY2htZW50Ug'
    'thdHRhY2htZW50cxIyCgdoZWFkZXJzGIqzvZoBIAMoCzIULmVtYWlsLk1lc3NhZ2VIZWFkZXJS'
    'B2hlYWRlcnMSJAoLdGVtcGxhdGVhcm4Yn4mytwEgASgJUgt0ZW1wbGF0ZWFybhJJCg90ZW1wbG'
    'F0ZWNvbnRlbnQYjf2d/AEgASgLMhsuZW1haWwuRW1haWxUZW1wbGF0ZUNvbnRlbnRSD3RlbXBs'
    'YXRlY29udGVudBImCgx0ZW1wbGF0ZWRhdGEYlO2+hwEgASgJUgx0ZW1wbGF0ZWRhdGESJgoMdG'
    'VtcGxhdGVuYW1lGLGYkeUBIAEoCVIMdGVtcGxhdGVuYW1l');

@$core.Deprecated('Use tenantDescriptor instead')
const Tenant$json = {
  '1': 'Tenant',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {
      '1': 'sendingstatus',
      '3': 420634540,
      '4': 1,
      '5': 14,
      '6': '.email.SendingStatus',
      '10': 'sendingstatus'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.email.Tag',
      '10': 'tags'
    },
    {'1': 'tenantarn', '3': 19777181, '4': 1, '5': 9, '10': 'tenantarn'},
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `Tenant`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tenantDescriptor = $convert.base64Decode(
    'CgZUZW5hbnQSLgoQY3JlYXRlZHRpbWVzdGFtcBj628+fASABKAlSEGNyZWF0ZWR0aW1lc3RhbX'
    'ASPgoNc2VuZGluZ3N0YXR1cxisv8nIASABKA4yFC5lbWFpbC5TZW5kaW5nU3RhdHVzUg1zZW5k'
    'aW5nc3RhdHVzEiIKBHRhZ3MYwcH2tQEgAygLMgouZW1haWwuVGFnUgR0YWdzEh8KCXRlbmFudG'
    'FybhidjbcJIAEoCVIJdGVuYW50YXJuEh0KCHRlbmFudGlkGIH472QgASgJUgh0ZW5hbnRpZBIh'
    'Cgp0ZW5hbnRuYW1lGIfc01IgASgJUgp0ZW5hbnRuYW1l');

@$core.Deprecated('Use tenantInfoDescriptor instead')
const TenantInfo$json = {
  '1': 'TenantInfo',
  '2': [
    {
      '1': 'createdtimestamp',
      '3': 334753274,
      '4': 1,
      '5': 9,
      '10': 'createdtimestamp'
    },
    {'1': 'tenantarn', '3': 19777181, '4': 1, '5': 9, '10': 'tenantarn'},
    {'1': 'tenantid', '3': 211549185, '4': 1, '5': 9, '10': 'tenantid'},
    {'1': 'tenantname', '3': 173338119, '4': 1, '5': 9, '10': 'tenantname'},
  ],
};

/// Descriptor for `TenantInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tenantInfoDescriptor = $convert.base64Decode(
    'CgpUZW5hbnRJbmZvEi4KEGNyZWF0ZWR0aW1lc3RhbXAY+tvPnwEgASgJUhBjcmVhdGVkdGltZX'
    'N0YW1wEh8KCXRlbmFudGFybhidjbcJIAEoCVIJdGVuYW50YXJuEh0KCHRlbmFudGlkGIH472Qg'
    'ASgJUgh0ZW5hbnRpZBIhCgp0ZW5hbnRuYW1lGIfc01IgASgJUgp0ZW5hbnRuYW1l');

@$core.Deprecated('Use tenantResourceDescriptor instead')
const TenantResource$json = {
  '1': 'TenantResource',
  '2': [
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'resourcetype',
      '3': 301342558,
      '4': 1,
      '5': 14,
      '6': '.email.ResourceType',
      '10': 'resourcetype'
    },
  ],
};

/// Descriptor for `TenantResource`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tenantResourceDescriptor = $convert.base64Decode(
    'Cg5UZW5hbnRSZXNvdXJjZRIkCgtyZXNvdXJjZWFybhit+NmtASABKAlSC3Jlc291cmNlYXJuEj'
    'sKDHJlc291cmNldHlwZRjevtiPASABKA4yEy5lbWFpbC5SZXNvdXJjZVR5cGVSDHJlc291cmNl'
    'dHlwZQ==');

@$core.Deprecated('Use testRenderEmailTemplateRequestDescriptor instead')
const TestRenderEmailTemplateRequest$json = {
  '1': 'TestRenderEmailTemplateRequest',
  '2': [
    {'1': 'templatedata', '3': 284145300, '4': 1, '5': 9, '10': 'templatedata'},
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `TestRenderEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testRenderEmailTemplateRequestDescriptor =
    $convert.base64Decode(
        'Ch5UZXN0UmVuZGVyRW1haWxUZW1wbGF0ZVJlcXVlc3QSJgoMdGVtcGxhdGVkYXRhGJTtvocBIA'
        'EoCVIMdGVtcGxhdGVkYXRhEiYKDHRlbXBsYXRlbmFtZRixmJHlASABKAlSDHRlbXBsYXRlbmFt'
        'ZQ==');

@$core.Deprecated('Use testRenderEmailTemplateResponseDescriptor instead')
const TestRenderEmailTemplateResponse$json = {
  '1': 'TestRenderEmailTemplateResponse',
  '2': [
    {
      '1': 'renderedtemplate',
      '3': 74031487,
      '4': 1,
      '5': 9,
      '10': 'renderedtemplate'
    },
  ],
};

/// Descriptor for `TestRenderEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List testRenderEmailTemplateResponseDescriptor =
    $convert.base64Decode(
        'Ch9UZXN0UmVuZGVyRW1haWxUZW1wbGF0ZVJlc3BvbnNlEi0KEHJlbmRlcmVkdGVtcGxhdGUY/8'
        'KmIyABKAlSEHJlbmRlcmVkdGVtcGxhdGU=');

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

@$core.Deprecated('Use topicDescriptor instead')
const Topic$json = {
  '1': 'Topic',
  '2': [
    {
      '1': 'defaultsubscriptionstatus',
      '3': 212540838,
      '4': 1,
      '5': 14,
      '6': '.email.SubscriptionStatus',
      '10': 'defaultsubscriptionstatus'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'displayname', '3': 418161847, '4': 1, '5': 9, '10': 'displayname'},
    {'1': 'topicname', '3': 446010240, '4': 1, '5': 9, '10': 'topicname'},
  ],
};

/// Descriptor for `Topic`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List topicDescriptor = $convert.base64Decode(
    'CgVUb3BpYxJaChlkZWZhdWx0c3Vic2NyaXB0aW9uc3RhdHVzGKa7rGUgASgOMhkuZW1haWwuU3'
    'Vic2NyaXB0aW9uU3RhdHVzUhlkZWZhdWx0c3Vic2NyaXB0aW9uc3RhdHVzEiMKC2Rlc2NyaXB0'
    'aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIkCgtkaXNwbGF5bmFtZRi3ybLHASABKAlSC2Rpc3'
    'BsYXluYW1lEiAKCXRvcGljbmFtZRiAp9bUASABKAlSCXRvcGljbmFtZQ==');

@$core.Deprecated('Use topicFilterDescriptor instead')
const TopicFilter$json = {
  '1': 'TopicFilter',
  '2': [
    {'1': 'topicname', '3': 446010240, '4': 1, '5': 9, '10': 'topicname'},
    {
      '1': 'usedefaultifpreferenceunavailable',
      '3': 534808888,
      '4': 1,
      '5': 8,
      '10': 'usedefaultifpreferenceunavailable'
    },
  ],
};

/// Descriptor for `TopicFilter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List topicFilterDescriptor = $convert.base64Decode(
    'CgtUb3BpY0ZpbHRlchIgCgl0b3BpY25hbWUYgKfW1AEgASgJUgl0b3BpY25hbWUSUAohdXNlZG'
    'VmYXVsdGlmcHJlZmVyZW5jZXVuYXZhaWxhYmxlGLiSgv8BIAEoCFIhdXNlZGVmYXVsdGlmcHJl'
    'ZmVyZW5jZXVuYXZhaWxhYmxl');

@$core.Deprecated('Use topicPreferenceDescriptor instead')
const TopicPreference$json = {
  '1': 'TopicPreference',
  '2': [
    {
      '1': 'subscriptionstatus',
      '3': 444964671,
      '4': 1,
      '5': 14,
      '6': '.email.SubscriptionStatus',
      '10': 'subscriptionstatus'
    },
    {'1': 'topicname', '3': 446010240, '4': 1, '5': 9, '10': 'topicname'},
  ],
};

/// Descriptor for `TopicPreference`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List topicPreferenceDescriptor = $convert.base64Decode(
    'Cg9Ub3BpY1ByZWZlcmVuY2USTQoSc3Vic2NyaXB0aW9uc3RhdHVzGL++ltQBIAEoDjIZLmVtYW'
    'lsLlN1YnNjcmlwdGlvblN0YXR1c1ISc3Vic2NyaXB0aW9uc3RhdHVzEiAKCXRvcGljbmFtZRiA'
    'p9bUASABKAlSCXRvcGljbmFtZQ==');

@$core.Deprecated('Use trackingOptionsDescriptor instead')
const TrackingOptions$json = {
  '1': 'TrackingOptions',
  '2': [
    {
      '1': 'customredirectdomain',
      '3': 72478039,
      '4': 1,
      '5': 9,
      '10': 'customredirectdomain'
    },
    {
      '1': 'httpspolicy',
      '3': 147115397,
      '4': 1,
      '5': 14,
      '6': '.email.HttpsPolicy',
      '10': 'httpspolicy'
    },
  ],
};

/// Descriptor for `TrackingOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List trackingOptionsDescriptor = $convert.base64Decode(
    'Cg9UcmFja2luZ09wdGlvbnMSNQoUY3VzdG9tcmVkaXJlY3Rkb21haW4Y19rHIiABKAlSFGN1c3'
    'RvbXJlZGlyZWN0ZG9tYWluEjcKC2h0dHBzcG9saWN5GIWbk0YgASgOMhIuZW1haWwuSHR0cHNQ'
    'b2xpY3lSC2h0dHBzcG9saWN5');

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

@$core.Deprecated(
    'Use updateConfigurationSetEventDestinationRequestDescriptor instead')
const UpdateConfigurationSetEventDestinationRequest$json = {
  '1': 'UpdateConfigurationSetEventDestinationRequest',
  '2': [
    {
      '1': 'configurationsetname',
      '3': 403457485,
      '4': 1,
      '5': 9,
      '10': 'configurationsetname'
    },
    {
      '1': 'eventdestination',
      '3': 469882902,
      '4': 1,
      '5': 11,
      '6': '.email.EventDestinationDefinition',
      '10': 'eventdestination'
    },
    {
      '1': 'eventdestinationname',
      '3': 477517655,
      '4': 1,
      '5': 9,
      '10': 'eventdestinationname'
    },
  ],
};

/// Descriptor for `UpdateConfigurationSetEventDestinationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateConfigurationSetEventDestinationRequestDescriptor =
    $convert.base64Decode(
        'Ci1VcGRhdGVDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvblJlcXVlc3QSNgoUY29uZm'
        'lndXJhdGlvbnNldG5hbWUYzYuxwAEgASgJUhRjb25maWd1cmF0aW9uc2V0bmFtZRJRChBldmVu'
        'dGRlc3RpbmF0aW9uGJawh+ABIAEoCzIhLmVtYWlsLkV2ZW50RGVzdGluYXRpb25EZWZpbml0aW'
        '9uUhBldmVudGRlc3RpbmF0aW9uEjYKFGV2ZW50ZGVzdGluYXRpb25uYW1lGNeu2eMBIAEoCVIU'
        'ZXZlbnRkZXN0aW5hdGlvbm5hbWU=');

@$core.Deprecated(
    'Use updateConfigurationSetEventDestinationResponseDescriptor instead')
const UpdateConfigurationSetEventDestinationResponse$json = {
  '1': 'UpdateConfigurationSetEventDestinationResponse',
};

/// Descriptor for `UpdateConfigurationSetEventDestinationResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateConfigurationSetEventDestinationResponseDescriptor =
    $convert.base64Decode(
        'Ci5VcGRhdGVDb25maWd1cmF0aW9uU2V0RXZlbnREZXN0aW5hdGlvblJlc3BvbnNl');

@$core.Deprecated('Use updateContactListRequestDescriptor instead')
const UpdateContactListRequest$json = {
  '1': 'UpdateContactListRequest',
  '2': [
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'topics',
      '3': 219850038,
      '4': 3,
      '5': 11,
      '6': '.email.Topic',
      '10': 'topics'
    },
  ],
};

/// Descriptor for `UpdateContactListRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContactListRequestDescriptor = $convert.base64Decode(
    'ChhVcGRhdGVDb250YWN0TGlzdFJlcXVlc3QSLAoPY29udGFjdGxpc3RuYW1lGPG9p6EBIAEoCV'
    'IPY29udGFjdGxpc3RuYW1lEiMKC2Rlc2NyaXB0aW9uGIr0+TYgASgJUgtkZXNjcmlwdGlvbhIn'
    'CgZ0b3BpY3MYtsrqaCADKAsyDC5lbWFpbC5Ub3BpY1IGdG9waWNz');

@$core.Deprecated('Use updateContactListResponseDescriptor instead')
const UpdateContactListResponse$json = {
  '1': 'UpdateContactListResponse',
};

/// Descriptor for `UpdateContactListResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContactListResponseDescriptor =
    $convert.base64Decode('ChlVcGRhdGVDb250YWN0TGlzdFJlc3BvbnNl');

@$core.Deprecated('Use updateContactRequestDescriptor instead')
const UpdateContactRequest$json = {
  '1': 'UpdateContactRequest',
  '2': [
    {
      '1': 'attributesdata',
      '3': 497271421,
      '4': 1,
      '5': 9,
      '10': 'attributesdata'
    },
    {
      '1': 'contactlistname',
      '3': 338288369,
      '4': 1,
      '5': 9,
      '10': 'contactlistname'
    },
    {'1': 'emailaddress', '3': 488814806, '4': 1, '5': 9, '10': 'emailaddress'},
    {
      '1': 'topicpreferences',
      '3': 199822983,
      '4': 3,
      '5': 11,
      '6': '.email.TopicPreference',
      '10': 'topicpreferences'
    },
    {
      '1': 'unsubscribeall',
      '3': 49261174,
      '4': 1,
      '5': 8,
      '10': 'unsubscribeall'
    },
  ],
};

/// Descriptor for `UpdateContactRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContactRequestDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVDb250YWN0UmVxdWVzdBIqCg5hdHRyaWJ1dGVzZGF0YRj9hI/tASABKAlSDmF0dH'
    'JpYnV0ZXNkYXRhEiwKD2NvbnRhY3RsaXN0bmFtZRjxvaehASABKAlSD2NvbnRhY3RsaXN0bmFt'
    'ZRImCgxlbWFpbGFkZHJlc3MY1vGK6QEgASgJUgxlbWFpbGFkZHJlc3MSRQoQdG9waWNwcmVmZX'
    'JlbmNlcxiHnaRfIAMoCzIWLmVtYWlsLlRvcGljUHJlZmVyZW5jZVIQdG9waWNwcmVmZXJlbmNl'
    'cxIpCg51bnN1YnNjcmliZWFsbBj21L4XIAEoCFIOdW5zdWJzY3JpYmVhbGw=');

@$core.Deprecated('Use updateContactResponseDescriptor instead')
const UpdateContactResponse$json = {
  '1': 'UpdateContactResponse',
};

/// Descriptor for `UpdateContactResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateContactResponseDescriptor =
    $convert.base64Decode('ChVVcGRhdGVDb250YWN0UmVzcG9uc2U=');

@$core.Deprecated(
    'Use updateCustomVerificationEmailTemplateRequestDescriptor instead')
const UpdateCustomVerificationEmailTemplateRequest$json = {
  '1': 'UpdateCustomVerificationEmailTemplateRequest',
  '2': [
    {
      '1': 'failureredirectionurl',
      '3': 376073007,
      '4': 1,
      '5': 9,
      '10': 'failureredirectionurl'
    },
    {
      '1': 'fromemailaddress',
      '3': 93506822,
      '4': 1,
      '5': 9,
      '10': 'fromemailaddress'
    },
    {
      '1': 'successredirectionurl',
      '3': 513750768,
      '4': 1,
      '5': 9,
      '10': 'successredirectionurl'
    },
    {
      '1': 'templatecontent',
      '3': 528973453,
      '4': 1,
      '5': 9,
      '10': 'templatecontent'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
    {
      '1': 'templatesubject',
      '3': 144653738,
      '4': 1,
      '5': 9,
      '10': 'templatesubject'
    },
  ],
};

/// Descriptor for `UpdateCustomVerificationEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateCustomVerificationEmailTemplateRequestDescriptor =
    $convert.base64Decode(
        'CixVcGRhdGVDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVxdWVzdBI4ChVmYWlsdX'
        'JlcmVkaXJlY3Rpb251cmwYr9apswEgASgJUhVmYWlsdXJlcmVkaXJlY3Rpb251cmwSLQoQZnJv'
        'bWVtYWlsYWRkcmVzcxiGmsssIAEoCVIQZnJvbWVtYWlsYWRkcmVzcxI4ChVzdWNjZXNzcmVkaX'
        'JlY3Rpb251cmwY8O389AEgASgJUhVzdWNjZXNzcmVkaXJlY3Rpb251cmwSLAoPdGVtcGxhdGVj'
        'b250ZW50GI39nfwBIAEoCVIPdGVtcGxhdGVjb250ZW50EiYKDHRlbXBsYXRlbmFtZRixmJHlAS'
        'ABKAlSDHRlbXBsYXRlbmFtZRIrCg90ZW1wbGF0ZXN1YmplY3QYqvv8RCABKAlSD3RlbXBsYXRl'
        'c3ViamVjdA==');

@$core.Deprecated(
    'Use updateCustomVerificationEmailTemplateResponseDescriptor instead')
const UpdateCustomVerificationEmailTemplateResponse$json = {
  '1': 'UpdateCustomVerificationEmailTemplateResponse',
};

/// Descriptor for `UpdateCustomVerificationEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateCustomVerificationEmailTemplateResponseDescriptor =
    $convert.base64Decode(
        'Ci1VcGRhdGVDdXN0b21WZXJpZmljYXRpb25FbWFpbFRlbXBsYXRlUmVzcG9uc2U=');

@$core.Deprecated('Use updateEmailIdentityPolicyRequestDescriptor instead')
const UpdateEmailIdentityPolicyRequest$json = {
  '1': 'UpdateEmailIdentityPolicyRequest',
  '2': [
    {
      '1': 'emailidentity',
      '3': 136088090,
      '4': 1,
      '5': 9,
      '10': 'emailidentity'
    },
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `UpdateEmailIdentityPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateEmailIdentityPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiBVcGRhdGVFbWFpbElkZW50aXR5UG9saWN5UmVxdWVzdBInCg1lbWFpbGlkZW50aXR5GJqU8k'
        'AgASgJUg1lbWFpbGlkZW50aXR5EhoKBnBvbGljeRig7/DgASABKAlSBnBvbGljeRIhCgpwb2xp'
        'Y3luYW1lGL31h38gASgJUgpwb2xpY3luYW1l');

@$core.Deprecated('Use updateEmailIdentityPolicyResponseDescriptor instead')
const UpdateEmailIdentityPolicyResponse$json = {
  '1': 'UpdateEmailIdentityPolicyResponse',
};

/// Descriptor for `UpdateEmailIdentityPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateEmailIdentityPolicyResponseDescriptor =
    $convert.base64Decode('CiFVcGRhdGVFbWFpbElkZW50aXR5UG9saWN5UmVzcG9uc2U=');

@$core.Deprecated('Use updateEmailTemplateRequestDescriptor instead')
const UpdateEmailTemplateRequest$json = {
  '1': 'UpdateEmailTemplateRequest',
  '2': [
    {
      '1': 'templatecontent',
      '3': 528973453,
      '4': 1,
      '5': 11,
      '6': '.email.EmailTemplateContent',
      '10': 'templatecontent'
    },
    {'1': 'templatename', '3': 480529457, '4': 1, '5': 9, '10': 'templatename'},
  ],
};

/// Descriptor for `UpdateEmailTemplateRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateEmailTemplateRequestDescriptor =
    $convert.base64Decode(
        'ChpVcGRhdGVFbWFpbFRlbXBsYXRlUmVxdWVzdBJJCg90ZW1wbGF0ZWNvbnRlbnQYjf2d/AEgAS'
        'gLMhsuZW1haWwuRW1haWxUZW1wbGF0ZUNvbnRlbnRSD3RlbXBsYXRlY29udGVudBImCgx0ZW1w'
        'bGF0ZW5hbWUYsZiR5QEgASgJUgx0ZW1wbGF0ZW5hbWU=');

@$core.Deprecated('Use updateEmailTemplateResponseDescriptor instead')
const UpdateEmailTemplateResponse$json = {
  '1': 'UpdateEmailTemplateResponse',
};

/// Descriptor for `UpdateEmailTemplateResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateEmailTemplateResponseDescriptor =
    $convert.base64Decode('ChtVcGRhdGVFbWFpbFRlbXBsYXRlUmVzcG9uc2U=');

@$core.Deprecated(
    'Use updateReputationEntityCustomerManagedStatusRequestDescriptor instead')
const UpdateReputationEntityCustomerManagedStatusRequest$json = {
  '1': 'UpdateReputationEntityCustomerManagedStatusRequest',
  '2': [
    {
      '1': 'reputationentityreference',
      '3': 414929111,
      '4': 1,
      '5': 9,
      '10': 'reputationentityreference'
    },
    {
      '1': 'reputationentitytype',
      '3': 98287826,
      '4': 1,
      '5': 14,
      '6': '.email.ReputationEntityType',
      '10': 'reputationentitytype'
    },
    {
      '1': 'sendingstatus',
      '3': 420634540,
      '4': 1,
      '5': 14,
      '6': '.email.SendingStatus',
      '10': 'sendingstatus'
    },
  ],
};

/// Descriptor for `UpdateReputationEntityCustomerManagedStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateReputationEntityCustomerManagedStatusRequestDescriptor =
    $convert.base64Decode(
        'CjJVcGRhdGVSZXB1dGF0aW9uRW50aXR5Q3VzdG9tZXJNYW5hZ2VkU3RhdHVzUmVxdWVzdBJACh'
        'lyZXB1dGF0aW9uZW50aXR5cmVmZXJlbmNlGNeh7cUBIAEoCVIZcmVwdXRhdGlvbmVudGl0eXJl'
        'ZmVyZW5jZRJSChRyZXB1dGF0aW9uZW50aXR5dHlwZRjSge8uIAEoDjIbLmVtYWlsLlJlcHV0YX'
        'Rpb25FbnRpdHlUeXBlUhRyZXB1dGF0aW9uZW50aXR5dHlwZRI+Cg1zZW5kaW5nc3RhdHVzGKy/'
        'ycgBIAEoDjIULmVtYWlsLlNlbmRpbmdTdGF0dXNSDXNlbmRpbmdzdGF0dXM=');

@$core.Deprecated(
    'Use updateReputationEntityCustomerManagedStatusResponseDescriptor instead')
const UpdateReputationEntityCustomerManagedStatusResponse$json = {
  '1': 'UpdateReputationEntityCustomerManagedStatusResponse',
};

/// Descriptor for `UpdateReputationEntityCustomerManagedStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    updateReputationEntityCustomerManagedStatusResponseDescriptor =
    $convert.base64Decode(
        'CjNVcGRhdGVSZXB1dGF0aW9uRW50aXR5Q3VzdG9tZXJNYW5hZ2VkU3RhdHVzUmVzcG9uc2U=');

@$core.Deprecated('Use updateReputationEntityPolicyRequestDescriptor instead')
const UpdateReputationEntityPolicyRequest$json = {
  '1': 'UpdateReputationEntityPolicyRequest',
  '2': [
    {
      '1': 'reputationentitypolicy',
      '3': 166380308,
      '4': 1,
      '5': 9,
      '10': 'reputationentitypolicy'
    },
    {
      '1': 'reputationentityreference',
      '3': 414929111,
      '4': 1,
      '5': 9,
      '10': 'reputationentityreference'
    },
    {
      '1': 'reputationentitytype',
      '3': 98287826,
      '4': 1,
      '5': 14,
      '6': '.email.ReputationEntityType',
      '10': 'reputationentitytype'
    },
  ],
};

/// Descriptor for `UpdateReputationEntityPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateReputationEntityPolicyRequestDescriptor =
    $convert.base64Decode(
        'CiNVcGRhdGVSZXB1dGF0aW9uRW50aXR5UG9saWN5UmVxdWVzdBI5ChZyZXB1dGF0aW9uZW50aX'
        'R5cG9saWN5GJSGq08gASgJUhZyZXB1dGF0aW9uZW50aXR5cG9saWN5EkAKGXJlcHV0YXRpb25l'
        'bnRpdHlyZWZlcmVuY2UY16HtxQEgASgJUhlyZXB1dGF0aW9uZW50aXR5cmVmZXJlbmNlElIKFH'
        'JlcHV0YXRpb25lbnRpdHl0eXBlGNKB7y4gASgOMhsuZW1haWwuUmVwdXRhdGlvbkVudGl0eVR5'
        'cGVSFHJlcHV0YXRpb25lbnRpdHl0eXBl');

@$core.Deprecated('Use updateReputationEntityPolicyResponseDescriptor instead')
const UpdateReputationEntityPolicyResponse$json = {
  '1': 'UpdateReputationEntityPolicyResponse',
};

/// Descriptor for `UpdateReputationEntityPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateReputationEntityPolicyResponseDescriptor =
    $convert
        .base64Decode('CiRVcGRhdGVSZXB1dGF0aW9uRW50aXR5UG9saWN5UmVzcG9uc2U=');

@$core.Deprecated('Use vdmAttributesDescriptor instead')
const VdmAttributes$json = {
  '1': 'VdmAttributes',
  '2': [
    {
      '1': 'dashboardattributes',
      '3': 289648407,
      '4': 1,
      '5': 11,
      '6': '.email.DashboardAttributes',
      '10': 'dashboardattributes'
    },
    {
      '1': 'guardianattributes',
      '3': 163350388,
      '4': 1,
      '5': 11,
      '6': '.email.GuardianAttributes',
      '10': 'guardianattributes'
    },
    {
      '1': 'vdmenabled',
      '3': 74673388,
      '4': 1,
      '5': 14,
      '6': '.email.FeatureStatus',
      '10': 'vdmenabled'
    },
  ],
};

/// Descriptor for `VdmAttributes`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vdmAttributesDescriptor = $convert.base64Decode(
    'Cg1WZG1BdHRyaWJ1dGVzElAKE2Rhc2hib2FyZGF0dHJpYnV0ZXMYl96OigEgASgLMhouZW1haW'
    'wuRGFzaGJvYXJkQXR0cmlidXRlc1ITZGFzaGJvYXJkYXR0cmlidXRlcxJMChJndWFyZGlhbmF0'
    'dHJpYnV0ZXMY9I7yTSABKAsyGS5lbWFpbC5HdWFyZGlhbkF0dHJpYnV0ZXNSEmd1YXJkaWFuYX'
    'R0cmlidXRlcxI3Cgp2ZG1lbmFibGVkGOzZzSMgASgOMhQuZW1haWwuRmVhdHVyZVN0YXR1c1IK'
    'dmRtZW5hYmxlZA==');

@$core.Deprecated('Use vdmOptionsDescriptor instead')
const VdmOptions$json = {
  '1': 'VdmOptions',
  '2': [
    {
      '1': 'dashboardoptions',
      '3': 471420736,
      '4': 1,
      '5': 11,
      '6': '.email.DashboardOptions',
      '10': 'dashboardoptions'
    },
    {
      '1': 'guardianoptions',
      '3': 9430217,
      '4': 1,
      '5': 11,
      '6': '.email.GuardianOptions',
      '10': 'guardianoptions'
    },
  ],
};

/// Descriptor for `VdmOptions`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List vdmOptionsDescriptor = $convert.base64Decode(
    'CgpWZG1PcHRpb25zEkcKEGRhc2hib2FyZG9wdGlvbnMYwJ7l4AEgASgLMhcuZW1haWwuRGFzaG'
    'JvYXJkT3B0aW9uc1IQZGFzaGJvYXJkb3B0aW9ucxJDCg9ndWFyZGlhbm9wdGlvbnMYycm/BCAB'
    'KAsyFi5lbWFpbC5HdWFyZGlhbk9wdGlvbnNSD2d1YXJkaWFub3B0aW9ucw==');

@$core.Deprecated('Use verificationInfoDescriptor instead')
const VerificationInfo$json = {
  '1': 'VerificationInfo',
  '2': [
    {
      '1': 'errortype',
      '3': 398848954,
      '4': 1,
      '5': 14,
      '6': '.email.VerificationError',
      '10': 'errortype'
    },
    {
      '1': 'lastcheckedtimestamp',
      '3': 337895587,
      '4': 1,
      '5': 9,
      '10': 'lastcheckedtimestamp'
    },
    {
      '1': 'lastsuccesstimestamp',
      '3': 446188237,
      '4': 1,
      '5': 9,
      '10': 'lastsuccesstimestamp'
    },
    {
      '1': 'soarecord',
      '3': 456865288,
      '4': 1,
      '5': 11,
      '6': '.email.SOARecord',
      '10': 'soarecord'
    },
  ],
};

/// Descriptor for `VerificationInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verificationInfoDescriptor = $convert.base64Decode(
    'ChBWZXJpZmljYXRpb25JbmZvEjoKCWVycm9ydHlwZRi655e+ASABKA4yGC5lbWFpbC5WZXJpZm'
    'ljYXRpb25FcnJvclIJZXJyb3J0eXBlEjYKFGxhc3RjaGVja2VkdGltZXN0YW1wGKPBj6EBIAEo'
    'CVIUbGFzdGNoZWNrZWR0aW1lc3RhbXASNgoUbGFzdHN1Y2Nlc3N0aW1lc3RhbXAYzZXh1AEgAS'
    'gJUhRsYXN0c3VjY2Vzc3RpbWVzdGFtcBIyCglzb2FyZWNvcmQYiOzs2QEgASgLMhAuZW1haWwu'
    'U09BUmVjb3JkUglzb2FyZWNvcmQ=');

@$core.Deprecated('Use volumeStatisticsDescriptor instead')
const VolumeStatistics$json = {
  '1': 'VolumeStatistics',
  '2': [
    {
      '1': 'inboxrawcount',
      '3': 215061115,
      '4': 1,
      '5': 3,
      '10': 'inboxrawcount'
    },
    {
      '1': 'projectedinbox',
      '3': 411471566,
      '4': 1,
      '5': 3,
      '10': 'projectedinbox'
    },
    {
      '1': 'projectedspam',
      '3': 40322687,
      '4': 1,
      '5': 3,
      '10': 'projectedspam'
    },
    {'1': 'spamrawcount', '3': 440722366, '4': 1, '5': 3, '10': 'spamrawcount'},
  ],
};

/// Descriptor for `VolumeStatistics`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List volumeStatisticsDescriptor = $convert.base64Decode(
    'ChBWb2x1bWVTdGF0aXN0aWNzEicKDWluYm94cmF3Y291bnQY+6TGZiABKANSDWluYm94cmF3Y2'
    '91bnQSKgoOcHJvamVjdGVkaW5ib3gYzp2axAEgASgDUg5wcm9qZWN0ZWRpbmJveBInCg1wcm9q'
    'ZWN0ZWRzcGFtGP+MnRMgASgDUg1wcm9qZWN0ZWRzcGFtEiYKDHNwYW1yYXdjb3VudBi+x5PSAS'
    'ABKANSDHNwYW1yYXdjb3VudA==');
