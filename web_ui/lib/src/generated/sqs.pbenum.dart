// This is a generated file - do not edit.
//
// Generated from sqs.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class MessageSystemAttributeName extends $pb.ProtobufEnum {
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_SEQUENCENUMBER =
      MessageSystemAttributeName._(0,
          _omitEnumNames ? '' : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_SEQUENCENUMBER');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEDEDUPLICATIONID =
      MessageSystemAttributeName._(
          1,
          _omitEnumNames
              ? ''
              : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEDEDUPLICATIONID');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_DEADLETTERQUEUESOURCEARN =
      MessageSystemAttributeName._(
          2,
          _omitEnumNames
              ? ''
              : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_DEADLETTERQUEUESOURCEARN');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENTTIMESTAMP =
      MessageSystemAttributeName._(3,
          _omitEnumNames ? '' : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENTTIMESTAMP');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENDERID = MessageSystemAttributeName._(
          4, _omitEnumNames ? '' : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENDERID');
  static const MessageSystemAttributeName MESSAGE_SYSTEM_ATTRIBUTE_NAME_ALL =
      MessageSystemAttributeName._(
          5, _omitEnumNames ? '' : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_ALL');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEGROUPID =
      MessageSystemAttributeName._(6,
          _omitEnumNames ? '' : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEGROUPID');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATERECEIVECOUNT =
      MessageSystemAttributeName._(
          7,
          _omitEnumNames
              ? ''
              : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATERECEIVECOUNT');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_AWSTRACEHEADER =
      MessageSystemAttributeName._(8,
          _omitEnumNames ? '' : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_AWSTRACEHEADER');
  static const MessageSystemAttributeName
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATEFIRSTRECEIVETIMESTAMP =
      MessageSystemAttributeName._(
          9,
          _omitEnumNames
              ? ''
              : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATEFIRSTRECEIVETIMESTAMP');

  static const $core.List<MessageSystemAttributeName> values =
      <MessageSystemAttributeName>[
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_SEQUENCENUMBER,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEDEDUPLICATIONID,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_DEADLETTERQUEUESOURCEARN,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENTTIMESTAMP,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_SENDERID,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_ALL,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_MESSAGEGROUPID,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATERECEIVECOUNT,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_AWSTRACEHEADER,
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_APPROXIMATEFIRSTRECEIVETIMESTAMP,
  ];

  static final $core.List<MessageSystemAttributeName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 9);
  static MessageSystemAttributeName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MessageSystemAttributeName._(super.value, super.name);
}

class MessageSystemAttributeNameForSends extends $pb.ProtobufEnum {
  static const MessageSystemAttributeNameForSends
      MESSAGE_SYSTEM_ATTRIBUTE_NAME_FOR_SENDS_AWSTRACEHEADER =
      MessageSystemAttributeNameForSends._(
          0,
          _omitEnumNames
              ? ''
              : 'MESSAGE_SYSTEM_ATTRIBUTE_NAME_FOR_SENDS_AWSTRACEHEADER');

  static const $core.List<MessageSystemAttributeNameForSends> values =
      <MessageSystemAttributeNameForSends>[
    MESSAGE_SYSTEM_ATTRIBUTE_NAME_FOR_SENDS_AWSTRACEHEADER,
  ];

  static final $core.List<MessageSystemAttributeNameForSends?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static MessageSystemAttributeNameForSends? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MessageSystemAttributeNameForSends._(super.value, super.name);
}

class QueueAttributeName extends $pb.ProtobufEnum {
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_MAXIMUMMESSAGESIZE =
      QueueAttributeName._(
          0, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_MAXIMUMMESSAGESIZE');
  static const QueueAttributeName
      QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESDELAYED =
      QueueAttributeName._(
          1,
          _omitEnumNames
              ? ''
              : 'QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESDELAYED');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_DEDUPLICATIONSCOPE =
      QueueAttributeName._(
          2, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_DEDUPLICATIONSCOPE');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_MESSAGERETENTIONPERIOD =
      QueueAttributeName._(3,
          _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_MESSAGERETENTIONPERIOD');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_DELAYSECONDS =
      QueueAttributeName._(
          4, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_DELAYSECONDS');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_QUEUEARN =
      QueueAttributeName._(
          5, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_QUEUEARN');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_FIFOTHROUGHPUTLIMIT =
      QueueAttributeName._(
          6, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_FIFOTHROUGHPUTLIMIT');
  static const QueueAttributeName
      QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGES = QueueAttributeName._(
          7,
          _omitEnumNames
              ? ''
              : 'QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGES');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_FIFOQUEUE =
      QueueAttributeName._(
          8, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_FIFOQUEUE');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_CREATEDTIMESTAMP =
      QueueAttributeName._(
          9, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_CREATEDTIMESTAMP');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_ALL =
      QueueAttributeName._(
          10, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_ALL');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_KMSMASTERKEYID =
      QueueAttributeName._(
          11, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_KMSMASTERKEYID');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_LASTMODIFIEDTIMESTAMP =
      QueueAttributeName._(12,
          _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_LASTMODIFIEDTIMESTAMP');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_REDRIVEALLOWPOLICY =
      QueueAttributeName._(
          13, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_REDRIVEALLOWPOLICY');
  static const QueueAttributeName
      QUEUE_ATTRIBUTE_NAME_RECEIVEMESSAGEWAITTIMESECONDS = QueueAttributeName._(
          14,
          _omitEnumNames
              ? ''
              : 'QUEUE_ATTRIBUTE_NAME_RECEIVEMESSAGEWAITTIMESECONDS');
  static const QueueAttributeName
      QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESNOTVISIBLE =
      QueueAttributeName._(
          15,
          _omitEnumNames
              ? ''
              : 'QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESNOTVISIBLE');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_SQSMANAGEDSSEENABLED =
      QueueAttributeName._(16,
          _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_SQSMANAGEDSSEENABLED');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_REDRIVEPOLICY =
      QueueAttributeName._(
          17, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_REDRIVEPOLICY');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_POLICY =
      QueueAttributeName._(
          18, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_POLICY');
  static const QueueAttributeName
      QUEUE_ATTRIBUTE_NAME_KMSDATAKEYREUSEPERIODSECONDS = QueueAttributeName._(
          19,
          _omitEnumNames
              ? ''
              : 'QUEUE_ATTRIBUTE_NAME_KMSDATAKEYREUSEPERIODSECONDS');
  static const QueueAttributeName QUEUE_ATTRIBUTE_NAME_VISIBILITYTIMEOUT =
      QueueAttributeName._(
          20, _omitEnumNames ? '' : 'QUEUE_ATTRIBUTE_NAME_VISIBILITYTIMEOUT');
  static const QueueAttributeName
      QUEUE_ATTRIBUTE_NAME_CONTENTBASEDDEDUPLICATION = QueueAttributeName._(
          21,
          _omitEnumNames
              ? ''
              : 'QUEUE_ATTRIBUTE_NAME_CONTENTBASEDDEDUPLICATION');

  static const $core.List<QueueAttributeName> values = <QueueAttributeName>[
    QUEUE_ATTRIBUTE_NAME_MAXIMUMMESSAGESIZE,
    QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESDELAYED,
    QUEUE_ATTRIBUTE_NAME_DEDUPLICATIONSCOPE,
    QUEUE_ATTRIBUTE_NAME_MESSAGERETENTIONPERIOD,
    QUEUE_ATTRIBUTE_NAME_DELAYSECONDS,
    QUEUE_ATTRIBUTE_NAME_QUEUEARN,
    QUEUE_ATTRIBUTE_NAME_FIFOTHROUGHPUTLIMIT,
    QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGES,
    QUEUE_ATTRIBUTE_NAME_FIFOQUEUE,
    QUEUE_ATTRIBUTE_NAME_CREATEDTIMESTAMP,
    QUEUE_ATTRIBUTE_NAME_ALL,
    QUEUE_ATTRIBUTE_NAME_KMSMASTERKEYID,
    QUEUE_ATTRIBUTE_NAME_LASTMODIFIEDTIMESTAMP,
    QUEUE_ATTRIBUTE_NAME_REDRIVEALLOWPOLICY,
    QUEUE_ATTRIBUTE_NAME_RECEIVEMESSAGEWAITTIMESECONDS,
    QUEUE_ATTRIBUTE_NAME_APPROXIMATENUMBEROFMESSAGESNOTVISIBLE,
    QUEUE_ATTRIBUTE_NAME_SQSMANAGEDSSEENABLED,
    QUEUE_ATTRIBUTE_NAME_REDRIVEPOLICY,
    QUEUE_ATTRIBUTE_NAME_POLICY,
    QUEUE_ATTRIBUTE_NAME_KMSDATAKEYREUSEPERIODSECONDS,
    QUEUE_ATTRIBUTE_NAME_VISIBILITYTIMEOUT,
    QUEUE_ATTRIBUTE_NAME_CONTENTBASEDDEDUPLICATION,
  ];

  static final $core.List<QueueAttributeName?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 21);
  static QueueAttributeName? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const QueueAttributeName._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
