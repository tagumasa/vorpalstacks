// This is a generated file - do not edit.
//
// Generated from email.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AttachmentContentDisposition extends $pb.ProtobufEnum {
  static const AttachmentContentDisposition
      ATTACHMENT_CONTENT_DISPOSITION_INLINE = AttachmentContentDisposition._(
          0, _omitEnumNames ? '' : 'ATTACHMENT_CONTENT_DISPOSITION_INLINE');
  static const AttachmentContentDisposition
      ATTACHMENT_CONTENT_DISPOSITION_ATTACHMENT =
      AttachmentContentDisposition._(
          1, _omitEnumNames ? '' : 'ATTACHMENT_CONTENT_DISPOSITION_ATTACHMENT');

  static const $core.List<AttachmentContentDisposition> values =
      <AttachmentContentDisposition>[
    ATTACHMENT_CONTENT_DISPOSITION_INLINE,
    ATTACHMENT_CONTENT_DISPOSITION_ATTACHMENT,
  ];

  static final $core.List<AttachmentContentDisposition?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static AttachmentContentDisposition? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AttachmentContentDisposition._(super.value, super.name);
}

class AttachmentContentTransferEncoding extends $pb.ProtobufEnum {
  static const AttachmentContentTransferEncoding
      ATTACHMENT_CONTENT_TRANSFER_ENCODING_QUOTED_PRINTABLE =
      AttachmentContentTransferEncoding._(
          0,
          _omitEnumNames
              ? ''
              : 'ATTACHMENT_CONTENT_TRANSFER_ENCODING_QUOTED_PRINTABLE');
  static const AttachmentContentTransferEncoding
      ATTACHMENT_CONTENT_TRANSFER_ENCODING_BASE64 =
      AttachmentContentTransferEncoding._(1,
          _omitEnumNames ? '' : 'ATTACHMENT_CONTENT_TRANSFER_ENCODING_BASE64');
  static const AttachmentContentTransferEncoding
      ATTACHMENT_CONTENT_TRANSFER_ENCODING_SEVEN_BIT =
      AttachmentContentTransferEncoding._(
          2,
          _omitEnumNames
              ? ''
              : 'ATTACHMENT_CONTENT_TRANSFER_ENCODING_SEVEN_BIT');

  static const $core.List<AttachmentContentTransferEncoding> values =
      <AttachmentContentTransferEncoding>[
    ATTACHMENT_CONTENT_TRANSFER_ENCODING_QUOTED_PRINTABLE,
    ATTACHMENT_CONTENT_TRANSFER_ENCODING_BASE64,
    ATTACHMENT_CONTENT_TRANSFER_ENCODING_SEVEN_BIT,
  ];

  static final $core.List<AttachmentContentTransferEncoding?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static AttachmentContentTransferEncoding? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AttachmentContentTransferEncoding._(super.value, super.name);
}

class BehaviorOnMxFailure extends $pb.ProtobufEnum {
  static const BehaviorOnMxFailure BEHAVIOR_ON_MX_FAILURE_USE_DEFAULT_VALUE =
      BehaviorOnMxFailure._(
          0, _omitEnumNames ? '' : 'BEHAVIOR_ON_MX_FAILURE_USE_DEFAULT_VALUE');
  static const BehaviorOnMxFailure BEHAVIOR_ON_MX_FAILURE_REJECT_MESSAGE =
      BehaviorOnMxFailure._(
          1, _omitEnumNames ? '' : 'BEHAVIOR_ON_MX_FAILURE_REJECT_MESSAGE');

  static const $core.List<BehaviorOnMxFailure> values = <BehaviorOnMxFailure>[
    BEHAVIOR_ON_MX_FAILURE_USE_DEFAULT_VALUE,
    BEHAVIOR_ON_MX_FAILURE_REJECT_MESSAGE,
  ];

  static final $core.List<BehaviorOnMxFailure?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static BehaviorOnMxFailure? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BehaviorOnMxFailure._(super.value, super.name);
}

class BounceType extends $pb.ProtobufEnum {
  static const BounceType BOUNCE_TYPE_UNDETERMINED =
      BounceType._(0, _omitEnumNames ? '' : 'BOUNCE_TYPE_UNDETERMINED');
  static const BounceType BOUNCE_TYPE_TRANSIENT =
      BounceType._(1, _omitEnumNames ? '' : 'BOUNCE_TYPE_TRANSIENT');
  static const BounceType BOUNCE_TYPE_PERMANENT =
      BounceType._(2, _omitEnumNames ? '' : 'BOUNCE_TYPE_PERMANENT');

  static const $core.List<BounceType> values = <BounceType>[
    BOUNCE_TYPE_UNDETERMINED,
    BOUNCE_TYPE_TRANSIENT,
    BOUNCE_TYPE_PERMANENT,
  ];

  static final $core.List<BounceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static BounceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BounceType._(super.value, super.name);
}

class BulkEmailStatus extends $pb.ProtobufEnum {
  static const BulkEmailStatus BULK_EMAIL_STATUS_MESSAGE_REJECTED =
      BulkEmailStatus._(
          0, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_MESSAGE_REJECTED');
  static const BulkEmailStatus BULK_EMAIL_STATUS_CONFIGURATION_SET_NOT_FOUND =
      BulkEmailStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'BULK_EMAIL_STATUS_CONFIGURATION_SET_NOT_FOUND');
  static const BulkEmailStatus BULK_EMAIL_STATUS_MAIL_FROM_DOMAIN_NOT_VERIFIED =
      BulkEmailStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'BULK_EMAIL_STATUS_MAIL_FROM_DOMAIN_NOT_VERIFIED');
  static const BulkEmailStatus BULK_EMAIL_STATUS_SUCCESS =
      BulkEmailStatus._(3, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_SUCCESS');
  static const BulkEmailStatus BULK_EMAIL_STATUS_INVALID_PARAMETER =
      BulkEmailStatus._(
          4, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_INVALID_PARAMETER');
  static const BulkEmailStatus BULK_EMAIL_STATUS_TEMPLATE_NOT_FOUND =
      BulkEmailStatus._(
          5, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_TEMPLATE_NOT_FOUND');
  static const BulkEmailStatus BULK_EMAIL_STATUS_ACCOUNT_SUSPENDED =
      BulkEmailStatus._(
          6, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_ACCOUNT_SUSPENDED');
  static const BulkEmailStatus BULK_EMAIL_STATUS_ACCOUNT_DAILY_QUOTA_EXCEEDED =
      BulkEmailStatus._(
          7,
          _omitEnumNames
              ? ''
              : 'BULK_EMAIL_STATUS_ACCOUNT_DAILY_QUOTA_EXCEEDED');
  static const BulkEmailStatus BULK_EMAIL_STATUS_TRANSIENT_FAILURE =
      BulkEmailStatus._(
          8, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_TRANSIENT_FAILURE');
  static const BulkEmailStatus BULK_EMAIL_STATUS_ACCOUNT_THROTTLED =
      BulkEmailStatus._(
          9, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_ACCOUNT_THROTTLED');
  static const BulkEmailStatus BULK_EMAIL_STATUS_INVALID_SENDING_POOL_NAME =
      BulkEmailStatus._(10,
          _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_INVALID_SENDING_POOL_NAME');
  static const BulkEmailStatus
      BULK_EMAIL_STATUS_CONFIGURATION_SET_SENDING_PAUSED = BulkEmailStatus._(
          11,
          _omitEnumNames
              ? ''
              : 'BULK_EMAIL_STATUS_CONFIGURATION_SET_SENDING_PAUSED');
  static const BulkEmailStatus BULK_EMAIL_STATUS_FAILED =
      BulkEmailStatus._(12, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_FAILED');
  static const BulkEmailStatus BULK_EMAIL_STATUS_ACCOUNT_SENDING_PAUSED =
      BulkEmailStatus._(
          13, _omitEnumNames ? '' : 'BULK_EMAIL_STATUS_ACCOUNT_SENDING_PAUSED');

  static const $core.List<BulkEmailStatus> values = <BulkEmailStatus>[
    BULK_EMAIL_STATUS_MESSAGE_REJECTED,
    BULK_EMAIL_STATUS_CONFIGURATION_SET_NOT_FOUND,
    BULK_EMAIL_STATUS_MAIL_FROM_DOMAIN_NOT_VERIFIED,
    BULK_EMAIL_STATUS_SUCCESS,
    BULK_EMAIL_STATUS_INVALID_PARAMETER,
    BULK_EMAIL_STATUS_TEMPLATE_NOT_FOUND,
    BULK_EMAIL_STATUS_ACCOUNT_SUSPENDED,
    BULK_EMAIL_STATUS_ACCOUNT_DAILY_QUOTA_EXCEEDED,
    BULK_EMAIL_STATUS_TRANSIENT_FAILURE,
    BULK_EMAIL_STATUS_ACCOUNT_THROTTLED,
    BULK_EMAIL_STATUS_INVALID_SENDING_POOL_NAME,
    BULK_EMAIL_STATUS_CONFIGURATION_SET_SENDING_PAUSED,
    BULK_EMAIL_STATUS_FAILED,
    BULK_EMAIL_STATUS_ACCOUNT_SENDING_PAUSED,
  ];

  static final $core.List<BulkEmailStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 13);
  static BulkEmailStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const BulkEmailStatus._(super.value, super.name);
}

class ContactLanguage extends $pb.ProtobufEnum {
  static const ContactLanguage CONTACT_LANGUAGE_EN =
      ContactLanguage._(0, _omitEnumNames ? '' : 'CONTACT_LANGUAGE_EN');
  static const ContactLanguage CONTACT_LANGUAGE_JA =
      ContactLanguage._(1, _omitEnumNames ? '' : 'CONTACT_LANGUAGE_JA');

  static const $core.List<ContactLanguage> values = <ContactLanguage>[
    CONTACT_LANGUAGE_EN,
    CONTACT_LANGUAGE_JA,
  ];

  static final $core.List<ContactLanguage?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ContactLanguage? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContactLanguage._(super.value, super.name);
}

class ContactListImportAction extends $pb.ProtobufEnum {
  static const ContactListImportAction CONTACT_LIST_IMPORT_ACTION_DELETE =
      ContactListImportAction._(
          0, _omitEnumNames ? '' : 'CONTACT_LIST_IMPORT_ACTION_DELETE');
  static const ContactListImportAction CONTACT_LIST_IMPORT_ACTION_PUT =
      ContactListImportAction._(
          1, _omitEnumNames ? '' : 'CONTACT_LIST_IMPORT_ACTION_PUT');

  static const $core.List<ContactListImportAction> values =
      <ContactListImportAction>[
    CONTACT_LIST_IMPORT_ACTION_DELETE,
    CONTACT_LIST_IMPORT_ACTION_PUT,
  ];

  static final $core.List<ContactListImportAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ContactListImportAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ContactListImportAction._(super.value, super.name);
}

class DataFormat extends $pb.ProtobufEnum {
  static const DataFormat DATA_FORMAT_JSON =
      DataFormat._(0, _omitEnumNames ? '' : 'DATA_FORMAT_JSON');
  static const DataFormat DATA_FORMAT_CSV =
      DataFormat._(1, _omitEnumNames ? '' : 'DATA_FORMAT_CSV');

  static const $core.List<DataFormat> values = <DataFormat>[
    DATA_FORMAT_JSON,
    DATA_FORMAT_CSV,
  ];

  static final $core.List<DataFormat?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DataFormat? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataFormat._(super.value, super.name);
}

class DeliverabilityDashboardAccountStatus extends $pb.ProtobufEnum {
  static const DeliverabilityDashboardAccountStatus
      DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_DISABLED =
      DeliverabilityDashboardAccountStatus._(
          0,
          _omitEnumNames
              ? ''
              : 'DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_DISABLED');
  static const DeliverabilityDashboardAccountStatus
      DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_PENDING_EXPIRATION =
      DeliverabilityDashboardAccountStatus._(
          1,
          _omitEnumNames
              ? ''
              : 'DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_PENDING_EXPIRATION');
  static const DeliverabilityDashboardAccountStatus
      DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_ACTIVE =
      DeliverabilityDashboardAccountStatus._(
          2,
          _omitEnumNames
              ? ''
              : 'DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_ACTIVE');

  static const $core.List<DeliverabilityDashboardAccountStatus> values =
      <DeliverabilityDashboardAccountStatus>[
    DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_DISABLED,
    DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_PENDING_EXPIRATION,
    DELIVERABILITY_DASHBOARD_ACCOUNT_STATUS_ACTIVE,
  ];

  static final $core.List<DeliverabilityDashboardAccountStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static DeliverabilityDashboardAccountStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeliverabilityDashboardAccountStatus._(super.value, super.name);
}

class DeliverabilityTestStatus extends $pb.ProtobufEnum {
  static const DeliverabilityTestStatus DELIVERABILITY_TEST_STATUS_IN_PROGRESS =
      DeliverabilityTestStatus._(
          0, _omitEnumNames ? '' : 'DELIVERABILITY_TEST_STATUS_IN_PROGRESS');
  static const DeliverabilityTestStatus DELIVERABILITY_TEST_STATUS_COMPLETED =
      DeliverabilityTestStatus._(
          1, _omitEnumNames ? '' : 'DELIVERABILITY_TEST_STATUS_COMPLETED');

  static const $core.List<DeliverabilityTestStatus> values =
      <DeliverabilityTestStatus>[
    DELIVERABILITY_TEST_STATUS_IN_PROGRESS,
    DELIVERABILITY_TEST_STATUS_COMPLETED,
  ];

  static final $core.List<DeliverabilityTestStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DeliverabilityTestStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeliverabilityTestStatus._(super.value, super.name);
}

class DeliveryEventType extends $pb.ProtobufEnum {
  static const DeliveryEventType DELIVERY_EVENT_TYPE_TRANSIENT_BOUNCE =
      DeliveryEventType._(
          0, _omitEnumNames ? '' : 'DELIVERY_EVENT_TYPE_TRANSIENT_BOUNCE');
  static const DeliveryEventType DELIVERY_EVENT_TYPE_SEND =
      DeliveryEventType._(1, _omitEnumNames ? '' : 'DELIVERY_EVENT_TYPE_SEND');
  static const DeliveryEventType DELIVERY_EVENT_TYPE_PERMANENT_BOUNCE =
      DeliveryEventType._(
          2, _omitEnumNames ? '' : 'DELIVERY_EVENT_TYPE_PERMANENT_BOUNCE');
  static const DeliveryEventType DELIVERY_EVENT_TYPE_COMPLAINT =
      DeliveryEventType._(
          3, _omitEnumNames ? '' : 'DELIVERY_EVENT_TYPE_COMPLAINT');
  static const DeliveryEventType DELIVERY_EVENT_TYPE_DELIVERY =
      DeliveryEventType._(
          4, _omitEnumNames ? '' : 'DELIVERY_EVENT_TYPE_DELIVERY');
  static const DeliveryEventType DELIVERY_EVENT_TYPE_UNDETERMINED_BOUNCE =
      DeliveryEventType._(
          5, _omitEnumNames ? '' : 'DELIVERY_EVENT_TYPE_UNDETERMINED_BOUNCE');

  static const $core.List<DeliveryEventType> values = <DeliveryEventType>[
    DELIVERY_EVENT_TYPE_TRANSIENT_BOUNCE,
    DELIVERY_EVENT_TYPE_SEND,
    DELIVERY_EVENT_TYPE_PERMANENT_BOUNCE,
    DELIVERY_EVENT_TYPE_COMPLAINT,
    DELIVERY_EVENT_TYPE_DELIVERY,
    DELIVERY_EVENT_TYPE_UNDETERMINED_BOUNCE,
  ];

  static final $core.List<DeliveryEventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static DeliveryEventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DeliveryEventType._(super.value, super.name);
}

class DimensionValueSource extends $pb.ProtobufEnum {
  static const DimensionValueSource DIMENSION_VALUE_SOURCE_MESSAGE_TAG =
      DimensionValueSource._(
          0, _omitEnumNames ? '' : 'DIMENSION_VALUE_SOURCE_MESSAGE_TAG');
  static const DimensionValueSource DIMENSION_VALUE_SOURCE_EMAIL_HEADER =
      DimensionValueSource._(
          1, _omitEnumNames ? '' : 'DIMENSION_VALUE_SOURCE_EMAIL_HEADER');
  static const DimensionValueSource DIMENSION_VALUE_SOURCE_LINK_TAG =
      DimensionValueSource._(
          2, _omitEnumNames ? '' : 'DIMENSION_VALUE_SOURCE_LINK_TAG');

  static const $core.List<DimensionValueSource> values = <DimensionValueSource>[
    DIMENSION_VALUE_SOURCE_MESSAGE_TAG,
    DIMENSION_VALUE_SOURCE_EMAIL_HEADER,
    DIMENSION_VALUE_SOURCE_LINK_TAG,
  ];

  static final $core.List<DimensionValueSource?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static DimensionValueSource? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DimensionValueSource._(super.value, super.name);
}

class DkimSigningAttributesOrigin extends $pb.ProtobufEnum {
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_3 =
      DkimSigningAttributesOrigin._(
          0,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_3');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_CENTRAL_1 =
      DkimSigningAttributesOrigin._(
          1,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_CENTRAL_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_CENTRAL_1 =
      DkimSigningAttributesOrigin._(
          2,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_CENTRAL_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_2 =
      DkimSigningAttributesOrigin._(
          3,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_2');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_EXTERNAL = DkimSigningAttributesOrigin._(
          4, _omitEnumNames ? '' : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_EXTERNAL');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_1 =
      DkimSigningAttributesOrigin._(
          5,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_1 =
      DkimSigningAttributesOrigin._(
          6,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_1 =
      DkimSigningAttributesOrigin._(
          7,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_5 =
      DkimSigningAttributesOrigin._(
          8,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_5');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES = DkimSigningAttributesOrigin._(
          9, _omitEnumNames ? '' : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_2 =
      DkimSigningAttributesOrigin._(
          10,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_2');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_WEST_1 =
      DkimSigningAttributesOrigin._(
          11,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_WEST_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_2 =
      DkimSigningAttributesOrigin._(
          12,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_2');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_SA_EAST_1 =
      DkimSigningAttributesOrigin._(
          13,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_SA_EAST_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_3 =
      DkimSigningAttributesOrigin._(
          14,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_3');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_SOUTH_1 =
      DkimSigningAttributesOrigin._(
          15,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_SOUTH_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_2 =
      DkimSigningAttributesOrigin._(
          16,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_2');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_2 =
      DkimSigningAttributesOrigin._(
          17,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_2');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_NORTH_1 =
      DkimSigningAttributesOrigin._(
          18,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_NORTH_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_1 =
      DkimSigningAttributesOrigin._(
          19,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_1 =
      DkimSigningAttributesOrigin._(
          20,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AF_SOUTH_1 =
      DkimSigningAttributesOrigin._(
          21,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AF_SOUTH_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_2 =
      DkimSigningAttributesOrigin._(
          22,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_2');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_2 =
      DkimSigningAttributesOrigin._(
          23,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_2');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_3 =
      DkimSigningAttributesOrigin._(
          24,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_3');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_SOUTH_1 =
      DkimSigningAttributesOrigin._(
          25,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_SOUTH_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_1 =
      DkimSigningAttributesOrigin._(
          26,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_1 =
      DkimSigningAttributesOrigin._(
          27,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_1');
  static const DkimSigningAttributesOrigin
      DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_IL_CENTRAL_1 =
      DkimSigningAttributesOrigin._(
          28,
          _omitEnumNames
              ? ''
              : 'DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_IL_CENTRAL_1');

  static const $core.List<DkimSigningAttributesOrigin> values =
      <DkimSigningAttributesOrigin>[
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_3,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_CENTRAL_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_CENTRAL_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_2,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_EXTERNAL,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_5,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_2,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_CA_WEST_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_NORTHEAST_2,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_SA_EAST_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_3,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_SOUTH_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTH_2,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_2,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_NORTH_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AF_SOUTH_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_WEST_2,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_WEST_2,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_AP_SOUTHEAST_3,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_ME_SOUTH_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_US_EAST_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_EU_CENTRAL_1,
    DKIM_SIGNING_ATTRIBUTES_ORIGIN_AWS_SES_IL_CENTRAL_1,
  ];

  static final $core.List<DkimSigningAttributesOrigin?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 28);
  static DkimSigningAttributesOrigin? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DkimSigningAttributesOrigin._(super.value, super.name);
}

class DkimSigningKeyLength extends $pb.ProtobufEnum {
  static const DkimSigningKeyLength DKIM_SIGNING_KEY_LENGTH_RSA_1024_BIT =
      DkimSigningKeyLength._(
          0, _omitEnumNames ? '' : 'DKIM_SIGNING_KEY_LENGTH_RSA_1024_BIT');
  static const DkimSigningKeyLength DKIM_SIGNING_KEY_LENGTH_RSA_2048_BIT =
      DkimSigningKeyLength._(
          1, _omitEnumNames ? '' : 'DKIM_SIGNING_KEY_LENGTH_RSA_2048_BIT');

  static const $core.List<DkimSigningKeyLength> values = <DkimSigningKeyLength>[
    DKIM_SIGNING_KEY_LENGTH_RSA_1024_BIT,
    DKIM_SIGNING_KEY_LENGTH_RSA_2048_BIT,
  ];

  static final $core.List<DkimSigningKeyLength?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DkimSigningKeyLength? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DkimSigningKeyLength._(super.value, super.name);
}

class DkimStatus extends $pb.ProtobufEnum {
  static const DkimStatus DKIM_STATUS_PENDING =
      DkimStatus._(0, _omitEnumNames ? '' : 'DKIM_STATUS_PENDING');
  static const DkimStatus DKIM_STATUS_SUCCESS =
      DkimStatus._(1, _omitEnumNames ? '' : 'DKIM_STATUS_SUCCESS');
  static const DkimStatus DKIM_STATUS_TEMPORARY_FAILURE =
      DkimStatus._(2, _omitEnumNames ? '' : 'DKIM_STATUS_TEMPORARY_FAILURE');
  static const DkimStatus DKIM_STATUS_FAILED =
      DkimStatus._(3, _omitEnumNames ? '' : 'DKIM_STATUS_FAILED');
  static const DkimStatus DKIM_STATUS_NOT_STARTED =
      DkimStatus._(4, _omitEnumNames ? '' : 'DKIM_STATUS_NOT_STARTED');

  static const $core.List<DkimStatus> values = <DkimStatus>[
    DKIM_STATUS_PENDING,
    DKIM_STATUS_SUCCESS,
    DKIM_STATUS_TEMPORARY_FAILURE,
    DKIM_STATUS_FAILED,
    DKIM_STATUS_NOT_STARTED,
  ];

  static final $core.List<DkimStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static DkimStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DkimStatus._(super.value, super.name);
}

class EmailAddressInsightsConfidenceVerdict extends $pb.ProtobufEnum {
  static const EmailAddressInsightsConfidenceVerdict
      EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_MEDIUM =
      EmailAddressInsightsConfidenceVerdict._(
          0,
          _omitEnumNames
              ? ''
              : 'EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_MEDIUM');
  static const EmailAddressInsightsConfidenceVerdict
      EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_LOW =
      EmailAddressInsightsConfidenceVerdict._(
          1,
          _omitEnumNames
              ? ''
              : 'EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_LOW');
  static const EmailAddressInsightsConfidenceVerdict
      EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_HIGH =
      EmailAddressInsightsConfidenceVerdict._(
          2,
          _omitEnumNames
              ? ''
              : 'EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_HIGH');

  static const $core.List<EmailAddressInsightsConfidenceVerdict> values =
      <EmailAddressInsightsConfidenceVerdict>[
    EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_MEDIUM,
    EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_LOW,
    EMAIL_ADDRESS_INSIGHTS_CONFIDENCE_VERDICT_HIGH,
  ];

  static final $core.List<EmailAddressInsightsConfidenceVerdict?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static EmailAddressInsightsConfidenceVerdict? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EmailAddressInsightsConfidenceVerdict._(super.value, super.name);
}

class EngagementEventType extends $pb.ProtobufEnum {
  static const EngagementEventType ENGAGEMENT_EVENT_TYPE_OPEN =
      EngagementEventType._(
          0, _omitEnumNames ? '' : 'ENGAGEMENT_EVENT_TYPE_OPEN');
  static const EngagementEventType ENGAGEMENT_EVENT_TYPE_CLICK =
      EngagementEventType._(
          1, _omitEnumNames ? '' : 'ENGAGEMENT_EVENT_TYPE_CLICK');

  static const $core.List<EngagementEventType> values = <EngagementEventType>[
    ENGAGEMENT_EVENT_TYPE_OPEN,
    ENGAGEMENT_EVENT_TYPE_CLICK,
  ];

  static final $core.List<EngagementEventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static EngagementEventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EngagementEventType._(super.value, super.name);
}

class EventType extends $pb.ProtobufEnum {
  static const EventType EVENT_TYPE_DELIVERY_DELAY =
      EventType._(0, _omitEnumNames ? '' : 'EVENT_TYPE_DELIVERY_DELAY');
  static const EventType EVENT_TYPE_RENDERING_FAILURE =
      EventType._(1, _omitEnumNames ? '' : 'EVENT_TYPE_RENDERING_FAILURE');
  static const EventType EVENT_TYPE_SEND =
      EventType._(2, _omitEnumNames ? '' : 'EVENT_TYPE_SEND');
  static const EventType EVENT_TYPE_BOUNCE =
      EventType._(3, _omitEnumNames ? '' : 'EVENT_TYPE_BOUNCE');
  static const EventType EVENT_TYPE_REJECT =
      EventType._(4, _omitEnumNames ? '' : 'EVENT_TYPE_REJECT');
  static const EventType EVENT_TYPE_OPEN =
      EventType._(5, _omitEnumNames ? '' : 'EVENT_TYPE_OPEN');
  static const EventType EVENT_TYPE_SUBSCRIPTION =
      EventType._(6, _omitEnumNames ? '' : 'EVENT_TYPE_SUBSCRIPTION');
  static const EventType EVENT_TYPE_COMPLAINT =
      EventType._(7, _omitEnumNames ? '' : 'EVENT_TYPE_COMPLAINT');
  static const EventType EVENT_TYPE_DELIVERY =
      EventType._(8, _omitEnumNames ? '' : 'EVENT_TYPE_DELIVERY');
  static const EventType EVENT_TYPE_CLICK =
      EventType._(9, _omitEnumNames ? '' : 'EVENT_TYPE_CLICK');

  static const $core.List<EventType> values = <EventType>[
    EVENT_TYPE_DELIVERY_DELAY,
    EVENT_TYPE_RENDERING_FAILURE,
    EVENT_TYPE_SEND,
    EVENT_TYPE_BOUNCE,
    EVENT_TYPE_REJECT,
    EVENT_TYPE_OPEN,
    EVENT_TYPE_SUBSCRIPTION,
    EVENT_TYPE_COMPLAINT,
    EVENT_TYPE_DELIVERY,
    EVENT_TYPE_CLICK,
  ];

  static final $core.List<EventType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 9);
  static EventType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EventType._(super.value, super.name);
}

class ExportSourceType extends $pb.ProtobufEnum {
  static const ExportSourceType EXPORT_SOURCE_TYPE_METRICS_DATA =
      ExportSourceType._(
          0, _omitEnumNames ? '' : 'EXPORT_SOURCE_TYPE_METRICS_DATA');
  static const ExportSourceType EXPORT_SOURCE_TYPE_MESSAGE_INSIGHTS =
      ExportSourceType._(
          1, _omitEnumNames ? '' : 'EXPORT_SOURCE_TYPE_MESSAGE_INSIGHTS');

  static const $core.List<ExportSourceType> values = <ExportSourceType>[
    EXPORT_SOURCE_TYPE_METRICS_DATA,
    EXPORT_SOURCE_TYPE_MESSAGE_INSIGHTS,
  ];

  static final $core.List<ExportSourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExportSourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExportSourceType._(super.value, super.name);
}

class FeatureStatus extends $pb.ProtobufEnum {
  static const FeatureStatus FEATURE_STATUS_DISABLED =
      FeatureStatus._(0, _omitEnumNames ? '' : 'FEATURE_STATUS_DISABLED');
  static const FeatureStatus FEATURE_STATUS_ENABLED =
      FeatureStatus._(1, _omitEnumNames ? '' : 'FEATURE_STATUS_ENABLED');

  static const $core.List<FeatureStatus> values = <FeatureStatus>[
    FEATURE_STATUS_DISABLED,
    FEATURE_STATUS_ENABLED,
  ];

  static final $core.List<FeatureStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static FeatureStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const FeatureStatus._(super.value, super.name);
}

class HttpsPolicy extends $pb.ProtobufEnum {
  static const HttpsPolicy HTTPS_POLICY_OPTIONAL =
      HttpsPolicy._(0, _omitEnumNames ? '' : 'HTTPS_POLICY_OPTIONAL');
  static const HttpsPolicy HTTPS_POLICY_REQUIRE =
      HttpsPolicy._(1, _omitEnumNames ? '' : 'HTTPS_POLICY_REQUIRE');
  static const HttpsPolicy HTTPS_POLICY_REQUIRE_OPEN_ONLY =
      HttpsPolicy._(2, _omitEnumNames ? '' : 'HTTPS_POLICY_REQUIRE_OPEN_ONLY');

  static const $core.List<HttpsPolicy> values = <HttpsPolicy>[
    HTTPS_POLICY_OPTIONAL,
    HTTPS_POLICY_REQUIRE,
    HTTPS_POLICY_REQUIRE_OPEN_ONLY,
  ];

  static final $core.List<HttpsPolicy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static HttpsPolicy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const HttpsPolicy._(super.value, super.name);
}

class IdentityType extends $pb.ProtobufEnum {
  static const IdentityType IDENTITY_TYPE_MANAGED_DOMAIN =
      IdentityType._(0, _omitEnumNames ? '' : 'IDENTITY_TYPE_MANAGED_DOMAIN');
  static const IdentityType IDENTITY_TYPE_DOMAIN =
      IdentityType._(1, _omitEnumNames ? '' : 'IDENTITY_TYPE_DOMAIN');
  static const IdentityType IDENTITY_TYPE_EMAIL_ADDRESS =
      IdentityType._(2, _omitEnumNames ? '' : 'IDENTITY_TYPE_EMAIL_ADDRESS');

  static const $core.List<IdentityType> values = <IdentityType>[
    IDENTITY_TYPE_MANAGED_DOMAIN,
    IDENTITY_TYPE_DOMAIN,
    IDENTITY_TYPE_EMAIL_ADDRESS,
  ];

  static final $core.List<IdentityType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static IdentityType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IdentityType._(super.value, super.name);
}

class ImportDestinationType extends $pb.ProtobufEnum {
  static const ImportDestinationType IMPORT_DESTINATION_TYPE_SUPPRESSION_LIST =
      ImportDestinationType._(
          0, _omitEnumNames ? '' : 'IMPORT_DESTINATION_TYPE_SUPPRESSION_LIST');
  static const ImportDestinationType IMPORT_DESTINATION_TYPE_CONTACT_LIST =
      ImportDestinationType._(
          1, _omitEnumNames ? '' : 'IMPORT_DESTINATION_TYPE_CONTACT_LIST');

  static const $core.List<ImportDestinationType> values =
      <ImportDestinationType>[
    IMPORT_DESTINATION_TYPE_SUPPRESSION_LIST,
    IMPORT_DESTINATION_TYPE_CONTACT_LIST,
  ];

  static final $core.List<ImportDestinationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ImportDestinationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportDestinationType._(super.value, super.name);
}

class JobStatus extends $pb.ProtobufEnum {
  static const JobStatus JOB_STATUS_PROCESSING =
      JobStatus._(0, _omitEnumNames ? '' : 'JOB_STATUS_PROCESSING');
  static const JobStatus JOB_STATUS_CANCELLED =
      JobStatus._(1, _omitEnumNames ? '' : 'JOB_STATUS_CANCELLED');
  static const JobStatus JOB_STATUS_COMPLETED =
      JobStatus._(2, _omitEnumNames ? '' : 'JOB_STATUS_COMPLETED');
  static const JobStatus JOB_STATUS_CREATED =
      JobStatus._(3, _omitEnumNames ? '' : 'JOB_STATUS_CREATED');
  static const JobStatus JOB_STATUS_FAILED =
      JobStatus._(4, _omitEnumNames ? '' : 'JOB_STATUS_FAILED');

  static const $core.List<JobStatus> values = <JobStatus>[
    JOB_STATUS_PROCESSING,
    JOB_STATUS_CANCELLED,
    JOB_STATUS_COMPLETED,
    JOB_STATUS_CREATED,
    JOB_STATUS_FAILED,
  ];

  static final $core.List<JobStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static JobStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const JobStatus._(super.value, super.name);
}

class ListRecommendationsFilterKey extends $pb.ProtobufEnum {
  static const ListRecommendationsFilterKey
      LIST_RECOMMENDATIONS_FILTER_KEY_STATUS = ListRecommendationsFilterKey._(
          0, _omitEnumNames ? '' : 'LIST_RECOMMENDATIONS_FILTER_KEY_STATUS');
  static const ListRecommendationsFilterKey
      LIST_RECOMMENDATIONS_FILTER_KEY_IMPACT = ListRecommendationsFilterKey._(
          1, _omitEnumNames ? '' : 'LIST_RECOMMENDATIONS_FILTER_KEY_IMPACT');
  static const ListRecommendationsFilterKey
      LIST_RECOMMENDATIONS_FILTER_KEY_TYPE = ListRecommendationsFilterKey._(
          2, _omitEnumNames ? '' : 'LIST_RECOMMENDATIONS_FILTER_KEY_TYPE');
  static const ListRecommendationsFilterKey
      LIST_RECOMMENDATIONS_FILTER_KEY_RESOURCE_ARN =
      ListRecommendationsFilterKey._(3,
          _omitEnumNames ? '' : 'LIST_RECOMMENDATIONS_FILTER_KEY_RESOURCE_ARN');

  static const $core.List<ListRecommendationsFilterKey> values =
      <ListRecommendationsFilterKey>[
    LIST_RECOMMENDATIONS_FILTER_KEY_STATUS,
    LIST_RECOMMENDATIONS_FILTER_KEY_IMPACT,
    LIST_RECOMMENDATIONS_FILTER_KEY_TYPE,
    LIST_RECOMMENDATIONS_FILTER_KEY_RESOURCE_ARN,
  ];

  static final $core.List<ListRecommendationsFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ListRecommendationsFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ListRecommendationsFilterKey._(super.value, super.name);
}

class ListTenantResourcesFilterKey extends $pb.ProtobufEnum {
  static const ListTenantResourcesFilterKey
      LIST_TENANT_RESOURCES_FILTER_KEY_RESOURCE_TYPE =
      ListTenantResourcesFilterKey._(
          0,
          _omitEnumNames
              ? ''
              : 'LIST_TENANT_RESOURCES_FILTER_KEY_RESOURCE_TYPE');

  static const $core.List<ListTenantResourcesFilterKey> values =
      <ListTenantResourcesFilterKey>[
    LIST_TENANT_RESOURCES_FILTER_KEY_RESOURCE_TYPE,
  ];

  static final $core.List<ListTenantResourcesFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ListTenantResourcesFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ListTenantResourcesFilterKey._(super.value, super.name);
}

class MailFromDomainStatus extends $pb.ProtobufEnum {
  static const MailFromDomainStatus MAIL_FROM_DOMAIN_STATUS_PENDING =
      MailFromDomainStatus._(
          0, _omitEnumNames ? '' : 'MAIL_FROM_DOMAIN_STATUS_PENDING');
  static const MailFromDomainStatus MAIL_FROM_DOMAIN_STATUS_SUCCESS =
      MailFromDomainStatus._(
          1, _omitEnumNames ? '' : 'MAIL_FROM_DOMAIN_STATUS_SUCCESS');
  static const MailFromDomainStatus MAIL_FROM_DOMAIN_STATUS_TEMPORARY_FAILURE =
      MailFromDomainStatus._(
          2, _omitEnumNames ? '' : 'MAIL_FROM_DOMAIN_STATUS_TEMPORARY_FAILURE');
  static const MailFromDomainStatus MAIL_FROM_DOMAIN_STATUS_FAILED =
      MailFromDomainStatus._(
          3, _omitEnumNames ? '' : 'MAIL_FROM_DOMAIN_STATUS_FAILED');

  static const $core.List<MailFromDomainStatus> values = <MailFromDomainStatus>[
    MAIL_FROM_DOMAIN_STATUS_PENDING,
    MAIL_FROM_DOMAIN_STATUS_SUCCESS,
    MAIL_FROM_DOMAIN_STATUS_TEMPORARY_FAILURE,
    MAIL_FROM_DOMAIN_STATUS_FAILED,
  ];

  static final $core.List<MailFromDomainStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static MailFromDomainStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MailFromDomainStatus._(super.value, super.name);
}

class MailType extends $pb.ProtobufEnum {
  static const MailType MAIL_TYPE_MARKETING =
      MailType._(0, _omitEnumNames ? '' : 'MAIL_TYPE_MARKETING');
  static const MailType MAIL_TYPE_TRANSACTIONAL =
      MailType._(1, _omitEnumNames ? '' : 'MAIL_TYPE_TRANSACTIONAL');

  static const $core.List<MailType> values = <MailType>[
    MAIL_TYPE_MARKETING,
    MAIL_TYPE_TRANSACTIONAL,
  ];

  static final $core.List<MailType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MailType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MailType._(super.value, super.name);
}

class Metric extends $pb.ProtobufEnum {
  static const Metric METRIC_TRANSIENT_BOUNCE =
      Metric._(0, _omitEnumNames ? '' : 'METRIC_TRANSIENT_BOUNCE');
  static const Metric METRIC_DELIVERY_CLICK =
      Metric._(1, _omitEnumNames ? '' : 'METRIC_DELIVERY_CLICK');
  static const Metric METRIC_SEND =
      Metric._(2, _omitEnumNames ? '' : 'METRIC_SEND');
  static const Metric METRIC_OPEN =
      Metric._(3, _omitEnumNames ? '' : 'METRIC_OPEN');
  static const Metric METRIC_DELIVERY_COMPLAINT =
      Metric._(4, _omitEnumNames ? '' : 'METRIC_DELIVERY_COMPLAINT');
  static const Metric METRIC_PERMANENT_BOUNCE =
      Metric._(5, _omitEnumNames ? '' : 'METRIC_PERMANENT_BOUNCE');
  static const Metric METRIC_DELIVERY_OPEN =
      Metric._(6, _omitEnumNames ? '' : 'METRIC_DELIVERY_OPEN');
  static const Metric METRIC_COMPLAINT =
      Metric._(7, _omitEnumNames ? '' : 'METRIC_COMPLAINT');
  static const Metric METRIC_DELIVERY =
      Metric._(8, _omitEnumNames ? '' : 'METRIC_DELIVERY');
  static const Metric METRIC_CLICK =
      Metric._(9, _omitEnumNames ? '' : 'METRIC_CLICK');

  static const $core.List<Metric> values = <Metric>[
    METRIC_TRANSIENT_BOUNCE,
    METRIC_DELIVERY_CLICK,
    METRIC_SEND,
    METRIC_OPEN,
    METRIC_DELIVERY_COMPLAINT,
    METRIC_PERMANENT_BOUNCE,
    METRIC_DELIVERY_OPEN,
    METRIC_COMPLAINT,
    METRIC_DELIVERY,
    METRIC_CLICK,
  ];

  static final $core.List<Metric?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 9);
  static Metric? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Metric._(super.value, super.name);
}

class MetricAggregation extends $pb.ProtobufEnum {
  static const MetricAggregation METRIC_AGGREGATION_VOLUME =
      MetricAggregation._(0, _omitEnumNames ? '' : 'METRIC_AGGREGATION_VOLUME');
  static const MetricAggregation METRIC_AGGREGATION_RATE =
      MetricAggregation._(1, _omitEnumNames ? '' : 'METRIC_AGGREGATION_RATE');

  static const $core.List<MetricAggregation> values = <MetricAggregation>[
    METRIC_AGGREGATION_VOLUME,
    METRIC_AGGREGATION_RATE,
  ];

  static final $core.List<MetricAggregation?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MetricAggregation? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricAggregation._(super.value, super.name);
}

class MetricDimensionName extends $pb.ProtobufEnum {
  static const MetricDimensionName METRIC_DIMENSION_NAME_EMAIL_IDENTITY =
      MetricDimensionName._(
          0, _omitEnumNames ? '' : 'METRIC_DIMENSION_NAME_EMAIL_IDENTITY');
  static const MetricDimensionName METRIC_DIMENSION_NAME_ISP =
      MetricDimensionName._(
          1, _omitEnumNames ? '' : 'METRIC_DIMENSION_NAME_ISP');
  static const MetricDimensionName METRIC_DIMENSION_NAME_CONFIGURATION_SET =
      MetricDimensionName._(
          2, _omitEnumNames ? '' : 'METRIC_DIMENSION_NAME_CONFIGURATION_SET');

  static const $core.List<MetricDimensionName> values = <MetricDimensionName>[
    METRIC_DIMENSION_NAME_EMAIL_IDENTITY,
    METRIC_DIMENSION_NAME_ISP,
    METRIC_DIMENSION_NAME_CONFIGURATION_SET,
  ];

  static final $core.List<MetricDimensionName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static MetricDimensionName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricDimensionName._(super.value, super.name);
}

class MetricNamespace extends $pb.ProtobufEnum {
  static const MetricNamespace METRIC_NAMESPACE_VDM =
      MetricNamespace._(0, _omitEnumNames ? '' : 'METRIC_NAMESPACE_VDM');

  static const $core.List<MetricNamespace> values = <MetricNamespace>[
    METRIC_NAMESPACE_VDM,
  ];

  static final $core.List<MetricNamespace?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static MetricNamespace? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MetricNamespace._(super.value, super.name);
}

class QueryErrorCode extends $pb.ProtobufEnum {
  static const QueryErrorCode QUERY_ERROR_CODE_ACCESS_DENIED = QueryErrorCode._(
      0, _omitEnumNames ? '' : 'QUERY_ERROR_CODE_ACCESS_DENIED');
  static const QueryErrorCode QUERY_ERROR_CODE_INTERNAL_FAILURE =
      QueryErrorCode._(
          1, _omitEnumNames ? '' : 'QUERY_ERROR_CODE_INTERNAL_FAILURE');

  static const $core.List<QueryErrorCode> values = <QueryErrorCode>[
    QUERY_ERROR_CODE_ACCESS_DENIED,
    QUERY_ERROR_CODE_INTERNAL_FAILURE,
  ];

  static final $core.List<QueryErrorCode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static QueryErrorCode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueryErrorCode._(super.value, super.name);
}

class RecommendationImpact extends $pb.ProtobufEnum {
  static const RecommendationImpact RECOMMENDATION_IMPACT_LOW =
      RecommendationImpact._(
          0, _omitEnumNames ? '' : 'RECOMMENDATION_IMPACT_LOW');
  static const RecommendationImpact RECOMMENDATION_IMPACT_HIGH =
      RecommendationImpact._(
          1, _omitEnumNames ? '' : 'RECOMMENDATION_IMPACT_HIGH');

  static const $core.List<RecommendationImpact> values = <RecommendationImpact>[
    RECOMMENDATION_IMPACT_LOW,
    RECOMMENDATION_IMPACT_HIGH,
  ];

  static final $core.List<RecommendationImpact?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RecommendationImpact? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecommendationImpact._(super.value, super.name);
}

class RecommendationStatus extends $pb.ProtobufEnum {
  static const RecommendationStatus RECOMMENDATION_STATUS_FIXED =
      RecommendationStatus._(
          0, _omitEnumNames ? '' : 'RECOMMENDATION_STATUS_FIXED');
  static const RecommendationStatus RECOMMENDATION_STATUS_OPEN =
      RecommendationStatus._(
          1, _omitEnumNames ? '' : 'RECOMMENDATION_STATUS_OPEN');

  static const $core.List<RecommendationStatus> values = <RecommendationStatus>[
    RECOMMENDATION_STATUS_FIXED,
    RECOMMENDATION_STATUS_OPEN,
  ];

  static final $core.List<RecommendationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RecommendationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecommendationStatus._(super.value, super.name);
}

class RecommendationType extends $pb.ProtobufEnum {
  static const RecommendationType RECOMMENDATION_TYPE_DMARC =
      RecommendationType._(
          0, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_DMARC');
  static const RecommendationType RECOMMENDATION_TYPE_BOUNCE =
      RecommendationType._(
          1, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_BOUNCE');
  static const RecommendationType RECOMMENDATION_TYPE_FEEDBACK_3P =
      RecommendationType._(
          2, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_FEEDBACK_3P');
  static const RecommendationType RECOMMENDATION_TYPE_IP_LISTING =
      RecommendationType._(
          3, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_IP_LISTING');
  static const RecommendationType RECOMMENDATION_TYPE_COMPLAINT =
      RecommendationType._(
          4, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_COMPLAINT');
  static const RecommendationType RECOMMENDATION_TYPE_DKIM =
      RecommendationType._(5, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_DKIM');
  static const RecommendationType RECOMMENDATION_TYPE_BIMI =
      RecommendationType._(6, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_BIMI');
  static const RecommendationType RECOMMENDATION_TYPE_SPF =
      RecommendationType._(7, _omitEnumNames ? '' : 'RECOMMENDATION_TYPE_SPF');

  static const $core.List<RecommendationType> values = <RecommendationType>[
    RECOMMENDATION_TYPE_DMARC,
    RECOMMENDATION_TYPE_BOUNCE,
    RECOMMENDATION_TYPE_FEEDBACK_3P,
    RECOMMENDATION_TYPE_IP_LISTING,
    RECOMMENDATION_TYPE_COMPLAINT,
    RECOMMENDATION_TYPE_DKIM,
    RECOMMENDATION_TYPE_BIMI,
    RECOMMENDATION_TYPE_SPF,
  ];

  static final $core.List<RecommendationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static RecommendationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RecommendationType._(super.value, super.name);
}

class ReputationEntityFilterKey extends $pb.ProtobufEnum {
  static const ReputationEntityFilterKey REPUTATION_ENTITY_FILTER_KEY_STATUS =
      ReputationEntityFilterKey._(
          0, _omitEnumNames ? '' : 'REPUTATION_ENTITY_FILTER_KEY_STATUS');
  static const ReputationEntityFilterKey
      REPUTATION_ENTITY_FILTER_KEY_ENTITY_TYPE = ReputationEntityFilterKey._(
          1, _omitEnumNames ? '' : 'REPUTATION_ENTITY_FILTER_KEY_ENTITY_TYPE');
  static const ReputationEntityFilterKey
      REPUTATION_ENTITY_FILTER_KEY_ENTITY_REFERENCE_PREFIX =
      ReputationEntityFilterKey._(
          2,
          _omitEnumNames
              ? ''
              : 'REPUTATION_ENTITY_FILTER_KEY_ENTITY_REFERENCE_PREFIX');
  static const ReputationEntityFilterKey
      REPUTATION_ENTITY_FILTER_KEY_REPUTATION_IMPACT =
      ReputationEntityFilterKey._(
          3,
          _omitEnumNames
              ? ''
              : 'REPUTATION_ENTITY_FILTER_KEY_REPUTATION_IMPACT');

  static const $core.List<ReputationEntityFilterKey> values =
      <ReputationEntityFilterKey>[
    REPUTATION_ENTITY_FILTER_KEY_STATUS,
    REPUTATION_ENTITY_FILTER_KEY_ENTITY_TYPE,
    REPUTATION_ENTITY_FILTER_KEY_ENTITY_REFERENCE_PREFIX,
    REPUTATION_ENTITY_FILTER_KEY_REPUTATION_IMPACT,
  ];

  static final $core.List<ReputationEntityFilterKey?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ReputationEntityFilterKey? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReputationEntityFilterKey._(super.value, super.name);
}

class ReputationEntityType extends $pb.ProtobufEnum {
  static const ReputationEntityType REPUTATION_ENTITY_TYPE_RESOURCE =
      ReputationEntityType._(
          0, _omitEnumNames ? '' : 'REPUTATION_ENTITY_TYPE_RESOURCE');

  static const $core.List<ReputationEntityType> values = <ReputationEntityType>[
    REPUTATION_ENTITY_TYPE_RESOURCE,
  ];

  static final $core.List<ReputationEntityType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static ReputationEntityType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReputationEntityType._(super.value, super.name);
}

class ResourceType extends $pb.ProtobufEnum {
  static const ResourceType RESOURCE_TYPE_EMAIL_IDENTITY =
      ResourceType._(0, _omitEnumNames ? '' : 'RESOURCE_TYPE_EMAIL_IDENTITY');
  static const ResourceType RESOURCE_TYPE_CONFIGURATION_SET = ResourceType._(
      1, _omitEnumNames ? '' : 'RESOURCE_TYPE_CONFIGURATION_SET');
  static const ResourceType RESOURCE_TYPE_EMAIL_TEMPLATE =
      ResourceType._(2, _omitEnumNames ? '' : 'RESOURCE_TYPE_EMAIL_TEMPLATE');

  static const $core.List<ResourceType> values = <ResourceType>[
    RESOURCE_TYPE_EMAIL_IDENTITY,
    RESOURCE_TYPE_CONFIGURATION_SET,
    RESOURCE_TYPE_EMAIL_TEMPLATE,
  ];

  static final $core.List<ResourceType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static ResourceType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ResourceType._(super.value, super.name);
}

class ReviewStatus extends $pb.ProtobufEnum {
  static const ReviewStatus REVIEW_STATUS_PENDING =
      ReviewStatus._(0, _omitEnumNames ? '' : 'REVIEW_STATUS_PENDING');
  static const ReviewStatus REVIEW_STATUS_GRANTED =
      ReviewStatus._(1, _omitEnumNames ? '' : 'REVIEW_STATUS_GRANTED');
  static const ReviewStatus REVIEW_STATUS_FAILED =
      ReviewStatus._(2, _omitEnumNames ? '' : 'REVIEW_STATUS_FAILED');
  static const ReviewStatus REVIEW_STATUS_DENIED =
      ReviewStatus._(3, _omitEnumNames ? '' : 'REVIEW_STATUS_DENIED');

  static const $core.List<ReviewStatus> values = <ReviewStatus>[
    REVIEW_STATUS_PENDING,
    REVIEW_STATUS_GRANTED,
    REVIEW_STATUS_FAILED,
    REVIEW_STATUS_DENIED,
  ];

  static final $core.List<ReviewStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static ReviewStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ReviewStatus._(super.value, super.name);
}

class ScalingMode extends $pb.ProtobufEnum {
  static const ScalingMode SCALING_MODE_STANDARD =
      ScalingMode._(0, _omitEnumNames ? '' : 'SCALING_MODE_STANDARD');
  static const ScalingMode SCALING_MODE_MANAGED =
      ScalingMode._(1, _omitEnumNames ? '' : 'SCALING_MODE_MANAGED');

  static const $core.List<ScalingMode> values = <ScalingMode>[
    SCALING_MODE_STANDARD,
    SCALING_MODE_MANAGED,
  ];

  static final $core.List<ScalingMode?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ScalingMode? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ScalingMode._(super.value, super.name);
}

class SendingStatus extends $pb.ProtobufEnum {
  static const SendingStatus SENDING_STATUS_DISABLED =
      SendingStatus._(0, _omitEnumNames ? '' : 'SENDING_STATUS_DISABLED');
  static const SendingStatus SENDING_STATUS_REINSTATED =
      SendingStatus._(1, _omitEnumNames ? '' : 'SENDING_STATUS_REINSTATED');
  static const SendingStatus SENDING_STATUS_ENABLED =
      SendingStatus._(2, _omitEnumNames ? '' : 'SENDING_STATUS_ENABLED');

  static const $core.List<SendingStatus> values = <SendingStatus>[
    SENDING_STATUS_DISABLED,
    SENDING_STATUS_REINSTATED,
    SENDING_STATUS_ENABLED,
  ];

  static final $core.List<SendingStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SendingStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SendingStatus._(super.value, super.name);
}

class Status extends $pb.ProtobufEnum {
  static const Status STATUS_READY =
      Status._(0, _omitEnumNames ? '' : 'STATUS_READY');
  static const Status STATUS_DELETING =
      Status._(1, _omitEnumNames ? '' : 'STATUS_DELETING');
  static const Status STATUS_CREATING =
      Status._(2, _omitEnumNames ? '' : 'STATUS_CREATING');
  static const Status STATUS_FAILED =
      Status._(3, _omitEnumNames ? '' : 'STATUS_FAILED');

  static const $core.List<Status> values = <Status>[
    STATUS_READY,
    STATUS_DELETING,
    STATUS_CREATING,
    STATUS_FAILED,
  ];

  static final $core.List<Status?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static Status? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const Status._(super.value, super.name);
}

class SubscriptionStatus extends $pb.ProtobufEnum {
  static const SubscriptionStatus SUBSCRIPTION_STATUS_OPT_OUT =
      SubscriptionStatus._(
          0, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_OPT_OUT');
  static const SubscriptionStatus SUBSCRIPTION_STATUS_OPT_IN =
      SubscriptionStatus._(
          1, _omitEnumNames ? '' : 'SUBSCRIPTION_STATUS_OPT_IN');

  static const $core.List<SubscriptionStatus> values = <SubscriptionStatus>[
    SUBSCRIPTION_STATUS_OPT_OUT,
    SUBSCRIPTION_STATUS_OPT_IN,
  ];

  static final $core.List<SubscriptionStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SubscriptionStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SubscriptionStatus._(super.value, super.name);
}

class SuppressionConfidenceVerdictThreshold extends $pb.ProtobufEnum {
  static const SuppressionConfidenceVerdictThreshold
      SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MEDIUM =
      SuppressionConfidenceVerdictThreshold._(
          0,
          _omitEnumNames
              ? ''
              : 'SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MEDIUM');
  static const SuppressionConfidenceVerdictThreshold
      SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MANAGED =
      SuppressionConfidenceVerdictThreshold._(
          1,
          _omitEnumNames
              ? ''
              : 'SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MANAGED');
  static const SuppressionConfidenceVerdictThreshold
      SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_HIGH =
      SuppressionConfidenceVerdictThreshold._(
          2,
          _omitEnumNames
              ? ''
              : 'SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_HIGH');

  static const $core.List<SuppressionConfidenceVerdictThreshold> values =
      <SuppressionConfidenceVerdictThreshold>[
    SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MEDIUM,
    SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_MANAGED,
    SUPPRESSION_CONFIDENCE_VERDICT_THRESHOLD_HIGH,
  ];

  static final $core.List<SuppressionConfidenceVerdictThreshold?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static SuppressionConfidenceVerdictThreshold? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SuppressionConfidenceVerdictThreshold._(super.value, super.name);
}

class SuppressionListImportAction extends $pb.ProtobufEnum {
  static const SuppressionListImportAction
      SUPPRESSION_LIST_IMPORT_ACTION_DELETE = SuppressionListImportAction._(
          0, _omitEnumNames ? '' : 'SUPPRESSION_LIST_IMPORT_ACTION_DELETE');
  static const SuppressionListImportAction SUPPRESSION_LIST_IMPORT_ACTION_PUT =
      SuppressionListImportAction._(
          1, _omitEnumNames ? '' : 'SUPPRESSION_LIST_IMPORT_ACTION_PUT');

  static const $core.List<SuppressionListImportAction> values =
      <SuppressionListImportAction>[
    SUPPRESSION_LIST_IMPORT_ACTION_DELETE,
    SUPPRESSION_LIST_IMPORT_ACTION_PUT,
  ];

  static final $core.List<SuppressionListImportAction?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SuppressionListImportAction? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SuppressionListImportAction._(super.value, super.name);
}

class SuppressionListReason extends $pb.ProtobufEnum {
  static const SuppressionListReason SUPPRESSION_LIST_REASON_BOUNCE =
      SuppressionListReason._(
          0, _omitEnumNames ? '' : 'SUPPRESSION_LIST_REASON_BOUNCE');
  static const SuppressionListReason SUPPRESSION_LIST_REASON_COMPLAINT =
      SuppressionListReason._(
          1, _omitEnumNames ? '' : 'SUPPRESSION_LIST_REASON_COMPLAINT');

  static const $core.List<SuppressionListReason> values =
      <SuppressionListReason>[
    SUPPRESSION_LIST_REASON_BOUNCE,
    SUPPRESSION_LIST_REASON_COMPLAINT,
  ];

  static final $core.List<SuppressionListReason?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static SuppressionListReason? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SuppressionListReason._(super.value, super.name);
}

class TlsPolicy extends $pb.ProtobufEnum {
  static const TlsPolicy TLS_POLICY_OPTIONAL =
      TlsPolicy._(0, _omitEnumNames ? '' : 'TLS_POLICY_OPTIONAL');
  static const TlsPolicy TLS_POLICY_REQUIRE =
      TlsPolicy._(1, _omitEnumNames ? '' : 'TLS_POLICY_REQUIRE');

  static const $core.List<TlsPolicy> values = <TlsPolicy>[
    TLS_POLICY_OPTIONAL,
    TLS_POLICY_REQUIRE,
  ];

  static final $core.List<TlsPolicy?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static TlsPolicy? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const TlsPolicy._(super.value, super.name);
}

class VerificationError extends $pb.ProtobufEnum {
  static const VerificationError
      VERIFICATION_ERROR_REPLICATION_PRIMARY_BYO_DKIM_NOT_SUPPORTED =
      VerificationError._(
          0,
          _omitEnumNames
              ? ''
              : 'VERIFICATION_ERROR_REPLICATION_PRIMARY_BYO_DKIM_NOT_SUPPORTED');
  static const VerificationError VERIFICATION_ERROR_REPLICATION_ACCESS_DENIED =
      VerificationError._(1,
          _omitEnumNames ? '' : 'VERIFICATION_ERROR_REPLICATION_ACCESS_DENIED');
  static const VerificationError VERIFICATION_ERROR_TYPE_NOT_FOUND =
      VerificationError._(
          2, _omitEnumNames ? '' : 'VERIFICATION_ERROR_TYPE_NOT_FOUND');
  static const VerificationError VERIFICATION_ERROR_DNS_SERVER_ERROR =
      VerificationError._(
          3, _omitEnumNames ? '' : 'VERIFICATION_ERROR_DNS_SERVER_ERROR');
  static const VerificationError
      VERIFICATION_ERROR_REPLICATION_PRIMARY_INVALID_REGION =
      VerificationError._(
          4,
          _omitEnumNames
              ? ''
              : 'VERIFICATION_ERROR_REPLICATION_PRIMARY_INVALID_REGION');
  static const VerificationError VERIFICATION_ERROR_INVALID_VALUE =
      VerificationError._(
          5, _omitEnumNames ? '' : 'VERIFICATION_ERROR_INVALID_VALUE');
  static const VerificationError
      VERIFICATION_ERROR_REPLICATION_PRIMARY_NOT_FOUND = VerificationError._(
          6,
          _omitEnumNames
              ? ''
              : 'VERIFICATION_ERROR_REPLICATION_PRIMARY_NOT_FOUND');
  static const VerificationError VERIFICATION_ERROR_HOST_NOT_FOUND =
      VerificationError._(
          7, _omitEnumNames ? '' : 'VERIFICATION_ERROR_HOST_NOT_FOUND');
  static const VerificationError
      VERIFICATION_ERROR_REPLICATION_REPLICA_AS_PRIMARY_NOT_SUPPORTED =
      VerificationError._(
          8,
          _omitEnumNames
              ? ''
              : 'VERIFICATION_ERROR_REPLICATION_REPLICA_AS_PRIMARY_NOT_SUPPORTED');
  static const VerificationError VERIFICATION_ERROR_SERVICE_ERROR =
      VerificationError._(
          9, _omitEnumNames ? '' : 'VERIFICATION_ERROR_SERVICE_ERROR');

  static const $core.List<VerificationError> values = <VerificationError>[
    VERIFICATION_ERROR_REPLICATION_PRIMARY_BYO_DKIM_NOT_SUPPORTED,
    VERIFICATION_ERROR_REPLICATION_ACCESS_DENIED,
    VERIFICATION_ERROR_TYPE_NOT_FOUND,
    VERIFICATION_ERROR_DNS_SERVER_ERROR,
    VERIFICATION_ERROR_REPLICATION_PRIMARY_INVALID_REGION,
    VERIFICATION_ERROR_INVALID_VALUE,
    VERIFICATION_ERROR_REPLICATION_PRIMARY_NOT_FOUND,
    VERIFICATION_ERROR_HOST_NOT_FOUND,
    VERIFICATION_ERROR_REPLICATION_REPLICA_AS_PRIMARY_NOT_SUPPORTED,
    VERIFICATION_ERROR_SERVICE_ERROR,
  ];

  static final $core.List<VerificationError?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 9);
  static VerificationError? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VerificationError._(super.value, super.name);
}

class VerificationStatus extends $pb.ProtobufEnum {
  static const VerificationStatus VERIFICATION_STATUS_PENDING =
      VerificationStatus._(
          0, _omitEnumNames ? '' : 'VERIFICATION_STATUS_PENDING');
  static const VerificationStatus VERIFICATION_STATUS_SUCCESS =
      VerificationStatus._(
          1, _omitEnumNames ? '' : 'VERIFICATION_STATUS_SUCCESS');
  static const VerificationStatus VERIFICATION_STATUS_TEMPORARY_FAILURE =
      VerificationStatus._(
          2, _omitEnumNames ? '' : 'VERIFICATION_STATUS_TEMPORARY_FAILURE');
  static const VerificationStatus VERIFICATION_STATUS_FAILED =
      VerificationStatus._(
          3, _omitEnumNames ? '' : 'VERIFICATION_STATUS_FAILED');
  static const VerificationStatus VERIFICATION_STATUS_NOT_STARTED =
      VerificationStatus._(
          4, _omitEnumNames ? '' : 'VERIFICATION_STATUS_NOT_STARTED');

  static const $core.List<VerificationStatus> values = <VerificationStatus>[
    VERIFICATION_STATUS_PENDING,
    VERIFICATION_STATUS_SUCCESS,
    VERIFICATION_STATUS_TEMPORARY_FAILURE,
    VERIFICATION_STATUS_FAILED,
    VERIFICATION_STATUS_NOT_STARTED,
  ];

  static final $core.List<VerificationStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static VerificationStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const VerificationStatus._(super.value, super.name);
}

class WarmupStatus extends $pb.ProtobufEnum {
  static const WarmupStatus WARMUP_STATUS_DONE =
      WarmupStatus._(0, _omitEnumNames ? '' : 'WARMUP_STATUS_DONE');
  static const WarmupStatus WARMUP_STATUS_NOT_APPLICABLE =
      WarmupStatus._(1, _omitEnumNames ? '' : 'WARMUP_STATUS_NOT_APPLICABLE');
  static const WarmupStatus WARMUP_STATUS_IN_PROGRESS =
      WarmupStatus._(2, _omitEnumNames ? '' : 'WARMUP_STATUS_IN_PROGRESS');

  static const $core.List<WarmupStatus> values = <WarmupStatus>[
    WARMUP_STATUS_DONE,
    WARMUP_STATUS_NOT_APPLICABLE,
    WARMUP_STATUS_IN_PROGRESS,
  ];

  static final $core.List<WarmupStatus?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static WarmupStatus? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WarmupStatus._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
