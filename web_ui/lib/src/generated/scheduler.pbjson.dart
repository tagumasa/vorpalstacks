// This is a generated file - do not edit.
//
// Generated from scheduler.proto.

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

@$core.Deprecated('Use awsVpcConfigurationDescriptor instead')
const AwsVpcConfiguration$json = {
  '1': 'AwsVpcConfiguration',
  '2': [
    {
      '1': 'assignpublicip',
      '3': 461653589,
      '4': 1,
      '5': 9,
      '10': 'assignpublicip'
    },
    {
      '1': 'securitygroups',
      '3': 515282516,
      '4': 3,
      '5': 9,
      '10': 'securitygroups'
    },
    {'1': 'subnets', '3': 414921506, '4': 3, '5': 9, '10': 'subnets'},
  ],
};

/// Descriptor for `AwsVpcConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List awsVpcConfigurationDescriptor = $convert.base64Decode(
    'ChNBd3NWcGNDb25maWd1cmF0aW9uEioKDmFzc2lnbnB1YmxpY2lwGNWMkdwBIAEoCVIOYXNzaW'
    'ducHVibGljaXASKgoOc2VjdXJpdHlncm91cHMY1Kza9QEgAygJUg5zZWN1cml0eWdyb3VwcxIc'
    'CgdzdWJuZXRzGKLm7MUBIAMoCVIHc3VibmV0cw==');

@$core.Deprecated('Use capacityProviderStrategyItemDescriptor instead')
const CapacityProviderStrategyItem$json = {
  '1': 'CapacityProviderStrategyItem',
  '2': [
    {'1': 'base', '3': 500995289, '4': 1, '5': 5, '10': 'base'},
    {
      '1': 'capacityprovider',
      '3': 109086449,
      '4': 1,
      '5': 9,
      '10': 'capacityprovider'
    },
    {'1': 'weight', '3': 278961850, '4': 1, '5': 5, '10': 'weight'},
  ],
};

/// Descriptor for `CapacityProviderStrategyItem`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List capacityProviderStrategyItemDescriptor =
    $convert.base64Decode(
        'ChxDYXBhY2l0eVByb3ZpZGVyU3RyYXRlZ3lJdGVtEhYKBGJhc2UY2any7gEgASgFUgRiYXNlEi'
        '0KEGNhcGFjaXR5cHJvdmlkZXIY8Y2CNCABKAlSEGNhcGFjaXR5cHJvdmlkZXISGgoGd2VpZ2h0'
        'GLq9goUBIAEoBVIGd2VpZ2h0');

@$core.Deprecated('Use conflictExceptionDescriptor instead')
const ConflictException$json = {
  '1': 'ConflictException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ConflictException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List conflictExceptionDescriptor = $convert.base64Decode(
    'ChFDb25mbGljdEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use createScheduleGroupInputDescriptor instead')
const CreateScheduleGroupInput$json = {
  '1': 'CreateScheduleGroupInput',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.scheduler.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `CreateScheduleGroupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduleGroupInputDescriptor = $convert.base64Decode(
    'ChhDcmVhdGVTY2hlZHVsZUdyb3VwSW5wdXQSIwoLY2xpZW50dG9rZW4YzPu7QSABKAlSC2NsaW'
    'VudHRva2VuEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSJgoEdGFncxjBwfa1ASADKAsyDi5zY2hl'
    'ZHVsZXIuVGFnUgR0YWdz');

@$core.Deprecated('Use createScheduleGroupOutputDescriptor instead')
const CreateScheduleGroupOutput$json = {
  '1': 'CreateScheduleGroupOutput',
  '2': [
    {
      '1': 'schedulegrouparn',
      '3': 148398979,
      '4': 1,
      '5': 9,
      '10': 'schedulegrouparn'
    },
  ],
};

/// Descriptor for `CreateScheduleGroupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduleGroupOutputDescriptor =
    $convert.base64Decode(
        'ChlDcmVhdGVTY2hlZHVsZUdyb3VwT3V0cHV0Ei0KEHNjaGVkdWxlZ3JvdXBhcm4Yg8fhRiABKA'
        'lSEHNjaGVkdWxlZ3JvdXBhcm4=');

@$core.Deprecated('Use createScheduleInputDescriptor instead')
const CreateScheduleInput$json = {
  '1': 'CreateScheduleInput',
  '2': [
    {
      '1': 'actionaftercompletion',
      '3': 282350906,
      '4': 1,
      '5': 9,
      '10': 'actionaftercompletion'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'flexibletimewindow',
      '3': 2518952,
      '4': 1,
      '5': 11,
      '6': '.scheduler.FlexibleTimeWindow',
      '10': 'flexibletimewindow'
    },
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'kmskeyarn', '3': 110041649, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleexpressiontimezone',
      '3': 400730326,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpressiontimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {
      '1': 'target',
      '3': 191361385,
      '4': 1,
      '5': 11,
      '6': '.scheduler.Target',
      '10': 'target'
    },
  ],
};

/// Descriptor for `CreateScheduleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduleInputDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVTY2hlZHVsZUlucHV0EjgKFWFjdGlvbmFmdGVyY29tcGxldGlvbhi6qtGGASABKA'
    'lSFWFjdGlvbmFmdGVyY29tcGxldGlvbhIjCgtjbGllbnR0b2tlbhjM+7tBIAEoCVILY2xpZW50'
    'dG9rZW4SIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEhsKB2VuZGRhdGUYz7'
    'P5JCABKAlSB2VuZGRhdGUSUAoSZmxleGlibGV0aW1ld2luZG93GKjfmQEgASgLMh0uc2NoZWR1'
    'bGVyLkZsZXhpYmxlVGltZVdpbmRvd1ISZmxleGlibGV0aW1ld2luZG93EiAKCWdyb3VwbmFtZR'
    'jIyqCqASABKAlSCWdyb3VwbmFtZRIfCglrbXNrZXlhcm4YsbS8NCABKAlSCWttc2tleWFybhIV'
    'CgRuYW1lGIfmgX8gASgJUgRuYW1lEjIKEnNjaGVkdWxlZXhwcmVzc2lvbhj/kdvUASABKAlSEn'
    'NjaGVkdWxlZXhwcmVzc2lvbhJCChpzY2hlZHVsZWV4cHJlc3Npb250aW1lem9uZRjW0Yq/ASAB'
    'KAlSGnNjaGVkdWxlZXhwcmVzc2lvbnRpbWV6b25lEiAKCXN0YXJ0ZGF0ZRj8+KDUASABKAlSCX'
    'N0YXJ0ZGF0ZRIYCgVzdGF0ZRiXybLvASABKAlSBXN0YXRlEiwKBnRhcmdldBjp4p9bIAEoCzIR'
    'LnNjaGVkdWxlci5UYXJnZXRSBnRhcmdldA==');

@$core.Deprecated('Use createScheduleOutputDescriptor instead')
const CreateScheduleOutput$json = {
  '1': 'CreateScheduleOutput',
  '2': [
    {'1': 'schedulearn', '3': 178445188, '4': 1, '5': 9, '10': 'schedulearn'},
  ],
};

/// Descriptor for `CreateScheduleOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createScheduleOutputDescriptor = $convert.base64Decode(
    'ChRDcmVhdGVTY2hlZHVsZU91dHB1dBIjCgtzY2hlZHVsZWFybhiEt4tVIAEoCVILc2NoZWR1bG'
    'Vhcm4=');

@$core.Deprecated('Use deadLetterConfigDescriptor instead')
const DeadLetterConfig$json = {
  '1': 'DeadLetterConfig',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `DeadLetterConfig`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deadLetterConfigDescriptor = $convert
    .base64Decode('ChBEZWFkTGV0dGVyQ29uZmlnEhQKA2Fybhidm+2/ASABKAlSA2Fybg==');

@$core.Deprecated('Use deleteScheduleGroupInputDescriptor instead')
const DeleteScheduleGroupInput$json = {
  '1': 'DeleteScheduleGroupInput',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteScheduleGroupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteScheduleGroupInputDescriptor =
    $convert.base64Decode(
        'ChhEZWxldGVTY2hlZHVsZUdyb3VwSW5wdXQSIwoLY2xpZW50dG9rZW4YzPu7QSABKAlSC2NsaW'
        'VudHRva2VuEhUKBG5hbWUYh+aBfyABKAlSBG5hbWU=');

@$core.Deprecated('Use deleteScheduleGroupOutputDescriptor instead')
const DeleteScheduleGroupOutput$json = {
  '1': 'DeleteScheduleGroupOutput',
};

/// Descriptor for `DeleteScheduleGroupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteScheduleGroupOutputDescriptor =
    $convert.base64Decode('ChlEZWxldGVTY2hlZHVsZUdyb3VwT3V0cHV0');

@$core.Deprecated('Use deleteScheduleInputDescriptor instead')
const DeleteScheduleInput$json = {
  '1': 'DeleteScheduleInput',
  '2': [
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `DeleteScheduleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteScheduleInputDescriptor = $convert.base64Decode(
    'ChNEZWxldGVTY2hlZHVsZUlucHV0EiMKC2NsaWVudHRva2VuGMz7u0EgASgJUgtjbGllbnR0b2'
    'tlbhIgCglncm91cG5hbWUYyMqgqgEgASgJUglncm91cG5hbWUSFQoEbmFtZRiH5oF/IAEoCVIE'
    'bmFtZQ==');

@$core.Deprecated('Use deleteScheduleOutputDescriptor instead')
const DeleteScheduleOutput$json = {
  '1': 'DeleteScheduleOutput',
};

/// Descriptor for `DeleteScheduleOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteScheduleOutputDescriptor =
    $convert.base64Decode('ChREZWxldGVTY2hlZHVsZU91dHB1dA==');

@$core.Deprecated('Use ecsParametersDescriptor instead')
const EcsParameters$json = {
  '1': 'EcsParameters',
  '2': [
    {
      '1': 'capacityproviderstrategy',
      '3': 273957206,
      '4': 3,
      '5': 11,
      '6': '.scheduler.CapacityProviderStrategyItem',
      '10': 'capacityproviderstrategy'
    },
    {
      '1': 'enableecsmanagedtags',
      '3': 146161174,
      '4': 1,
      '5': 8,
      '10': 'enableecsmanagedtags'
    },
    {
      '1': 'enableexecutecommand',
      '3': 451374779,
      '4': 1,
      '5': 8,
      '10': 'enableexecutecommand'
    },
    {'1': 'group', '3': 91525165, '4': 1, '5': 9, '10': 'group'},
    {'1': 'launchtype', '3': 184333335, '4': 1, '5': 9, '10': 'launchtype'},
    {
      '1': 'networkconfiguration',
      '3': 240088634,
      '4': 1,
      '5': 11,
      '6': '.scheduler.NetworkConfiguration',
      '10': 'networkconfiguration'
    },
    {
      '1': 'placementconstraints',
      '3': 248464365,
      '4': 3,
      '5': 11,
      '6': '.scheduler.PlacementConstraint',
      '10': 'placementconstraints'
    },
    {
      '1': 'placementstrategy',
      '3': 25036678,
      '4': 3,
      '5': 11,
      '6': '.scheduler.PlacementStrategy',
      '10': 'placementstrategy'
    },
    {
      '1': 'platformversion',
      '3': 139924287,
      '4': 1,
      '5': 9,
      '10': 'platformversion'
    },
    {
      '1': 'propagatetags',
      '3': 405557622,
      '4': 1,
      '5': 9,
      '10': 'propagatetags'
    },
    {'1': 'referenceid', '3': 291739032, '4': 1, '5': 9, '10': 'referenceid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.scheduler.EcsParameters.TagsEntry',
      '10': 'tags'
    },
    {'1': 'taskcount', '3': 398407508, '4': 1, '5': 5, '10': 'taskcount'},
    {
      '1': 'taskdefinitionarn',
      '3': 82234477,
      '4': 1,
      '5': 9,
      '10': 'taskdefinitionarn'
    },
  ],
  '3': [EcsParameters_TagsEntry$json],
};

@$core.Deprecated('Use ecsParametersDescriptor instead')
const EcsParameters_TagsEntry$json = {
  '1': 'TagsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `EcsParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List ecsParametersDescriptor = $convert.base64Decode(
    'Cg1FY3NQYXJhbWV0ZXJzEmcKGGNhcGFjaXR5cHJvdmlkZXJzdHJhdGVneRjWgtGCASADKAsyJy'
    '5zY2hlZHVsZXIuQ2FwYWNpdHlQcm92aWRlclN0cmF0ZWd5SXRlbVIYY2FwYWNpdHlwcm92aWRl'
    'cnN0cmF0ZWd5EjUKFGVuYWJsZWVjc21hbmFnZWR0YWdzGJb82EUgASgIUhRlbmFibGVlY3NtYW'
    '5hZ2VkdGFncxI2ChRlbmFibGVleGVjdXRlY29tbWFuZBi73Z3XASABKAhSFGVuYWJsZWV4ZWN1'
    'dGVjb21tYW5kEhcKBWdyb3VwGK2g0isgASgJUgVncm91cBIhCgpsYXVuY2h0eXBlGJfo8lcgAS'
    'gJUgpsYXVuY2h0eXBlElYKFG5ldHdvcmtjb25maWd1cmF0aW9uGLrsvXIgASgLMh8uc2NoZWR1'
    'bGVyLk5ldHdvcmtDb25maWd1cmF0aW9uUhRuZXR3b3JrY29uZmlndXJhdGlvbhJVChRwbGFjZW'
    '1lbnRjb25zdHJhaW50cxjth712IAMoCzIeLnNjaGVkdWxlci5QbGFjZW1lbnRDb25zdHJhaW50'
    'UhRwbGFjZW1lbnRjb25zdHJhaW50cxJNChFwbGFjZW1lbnRzdHJhdGVneRiGj/gLIAMoCzIcLn'
    'NjaGVkdWxlci5QbGFjZW1lbnRTdHJhdGVneVIRcGxhY2VtZW50c3RyYXRlZ3kSKwoPcGxhdGZv'
    'cm12ZXJzaW9uGL+m3EIgASgJUg9wbGF0Zm9ybXZlcnNpb24SKAoNcHJvcGFnYXRldGFncxj2or'
    'HBASABKAlSDXByb3BhZ2F0ZXRhZ3MSJAoLcmVmZXJlbmNlaWQYmKuOiwEgASgJUgtyZWZlcmVu'
    'Y2VpZBI6CgR0YWdzGMHB9rUBIAMoCzIiLnNjaGVkdWxlci5FY3NQYXJhbWV0ZXJzLlRhZ3NFbn'
    'RyeVIEdGFncxIgCgl0YXNrY291bnQY1O78vQEgASgFUgl0YXNrY291bnQSLwoRdGFza2RlZmlu'
    'aXRpb25hcm4Y7ZibJyABKAlSEXRhc2tkZWZpbml0aW9uYXJuGjcKCVRhZ3NFbnRyeRIQCgNrZX'
    'kYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use eventBridgeParametersDescriptor instead')
const EventBridgeParameters$json = {
  '1': 'EventBridgeParameters',
  '2': [
    {'1': 'detailtype', '3': 11758589, '4': 1, '5': 9, '10': 'detailtype'},
    {'1': 'source', '3': 31630329, '4': 1, '5': 9, '10': 'source'},
  ],
};

/// Descriptor for `EventBridgeParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List eventBridgeParametersDescriptor = $convert.base64Decode(
    'ChVFdmVudEJyaWRnZVBhcmFtZXRlcnMSIQoKZGV0YWlsdHlwZRj9180FIAEoCVIKZGV0YWlsdH'
    'lwZRIZCgZzb3VyY2UY+ceKDyABKAlSBnNvdXJjZQ==');

@$core.Deprecated('Use flexibleTimeWindowDescriptor instead')
const FlexibleTimeWindow$json = {
  '1': 'FlexibleTimeWindow',
  '2': [
    {
      '1': 'maximumwindowinminutes',
      '3': 482931584,
      '4': 1,
      '5': 5,
      '10': 'maximumwindowinminutes'
    },
    {'1': 'mode', '3': 323909427, '4': 1, '5': 9, '10': 'mode'},
  ],
};

/// Descriptor for `FlexibleTimeWindow`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List flexibleTimeWindowDescriptor = $convert.base64Decode(
    'ChJGbGV4aWJsZVRpbWVXaW5kb3cSOgoWbWF4aW11bXdpbmRvd2lubWludXRlcxiA56PmASABKA'
    'VSFm1heGltdW13aW5kb3dpbm1pbnV0ZXMSFgoEbW9kZRiz7rmaASABKAlSBG1vZGU=');

@$core.Deprecated('Use getScheduleGroupInputDescriptor instead')
const GetScheduleGroupInput$json = {
  '1': 'GetScheduleGroupInput',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetScheduleGroupInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduleGroupInputDescriptor =
    $convert.base64Decode(
        'ChVHZXRTY2hlZHVsZUdyb3VwSW5wdXQSFQoEbmFtZRiH5oF/IAEoCVIEbmFtZQ==');

@$core.Deprecated('Use getScheduleGroupOutputDescriptor instead')
const GetScheduleGroupOutput$json = {
  '1': 'GetScheduleGroupOutput',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'lastmodificationdate',
      '3': 4750358,
      '4': 1,
      '5': 9,
      '10': 'lastmodificationdate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `GetScheduleGroupOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduleGroupOutputDescriptor = $convert.base64Decode(
    'ChZHZXRTY2hlZHVsZUdyb3VwT3V0cHV0EhQKA2Fybhidm+2/ASABKAlSA2FybhImCgxjcmVhdG'
    'lvbmRhdGUY4di3iQEgASgJUgxjcmVhdGlvbmRhdGUSNQoUbGFzdG1vZGlmaWNhdGlvbmRhdGUY'
    'lvihAiABKAlSFGxhc3Rtb2RpZmljYXRpb25kYXRlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGA'
    'oFc3RhdGUYl8my7wEgASgJUgVzdGF0ZQ==');

@$core.Deprecated('Use getScheduleInputDescriptor instead')
const GetScheduleInput$json = {
  '1': 'GetScheduleInput',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
  ],
};

/// Descriptor for `GetScheduleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduleInputDescriptor = $convert.base64Decode(
    'ChBHZXRTY2hlZHVsZUlucHV0EiAKCWdyb3VwbmFtZRjIyqCqASABKAlSCWdyb3VwbmFtZRIVCg'
    'RuYW1lGIfmgX8gASgJUgRuYW1l');

@$core.Deprecated('Use getScheduleOutputDescriptor instead')
const GetScheduleOutput$json = {
  '1': 'GetScheduleOutput',
  '2': [
    {
      '1': 'actionaftercompletion',
      '3': 282350906,
      '4': 1,
      '5': 9,
      '10': 'actionaftercompletion'
    },
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'flexibletimewindow',
      '3': 2518952,
      '4': 1,
      '5': 11,
      '6': '.scheduler.FlexibleTimeWindow',
      '10': 'flexibletimewindow'
    },
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'kmskeyarn', '3': 110041649, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {
      '1': 'lastmodificationdate',
      '3': 4750358,
      '4': 1,
      '5': 9,
      '10': 'lastmodificationdate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleexpressiontimezone',
      '3': 400730326,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpressiontimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {
      '1': 'target',
      '3': 191361385,
      '4': 1,
      '5': 11,
      '6': '.scheduler.Target',
      '10': 'target'
    },
  ],
};

/// Descriptor for `GetScheduleOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getScheduleOutputDescriptor = $convert.base64Decode(
    'ChFHZXRTY2hlZHVsZU91dHB1dBI4ChVhY3Rpb25hZnRlcmNvbXBsZXRpb24YuqrRhgEgASgJUh'
    'VhY3Rpb25hZnRlcmNvbXBsZXRpb24SFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEiYKDGNyZWF0aW9u'
    'ZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9uZGF0ZRIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZG'
    'VzY3JpcHRpb24SGwoHZW5kZGF0ZRjPs/kkIAEoCVIHZW5kZGF0ZRJQChJmbGV4aWJsZXRpbWV3'
    'aW5kb3cYqN+ZASABKAsyHS5zY2hlZHVsZXIuRmxleGlibGVUaW1lV2luZG93UhJmbGV4aWJsZX'
    'RpbWV3aW5kb3cSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1lEh8KCWttc2tleWFy'
    'bhixtLw0IAEoCVIJa21za2V5YXJuEjUKFGxhc3Rtb2RpZmljYXRpb25kYXRlGJb4oQIgASgJUh'
    'RsYXN0bW9kaWZpY2F0aW9uZGF0ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEjIKEnNjaGVkdWxl'
    'ZXhwcmVzc2lvbhj/kdvUASABKAlSEnNjaGVkdWxlZXhwcmVzc2lvbhJCChpzY2hlZHVsZWV4cH'
    'Jlc3Npb250aW1lem9uZRjW0Yq/ASABKAlSGnNjaGVkdWxlZXhwcmVzc2lvbnRpbWV6b25lEiAK'
    'CXN0YXJ0ZGF0ZRj8+KDUASABKAlSCXN0YXJ0ZGF0ZRIYCgVzdGF0ZRiXybLvASABKAlSBXN0YX'
    'RlEiwKBnRhcmdldBjp4p9bIAEoCzIRLnNjaGVkdWxlci5UYXJnZXRSBnRhcmdldA==');

@$core.Deprecated('Use internalServerExceptionDescriptor instead')
const InternalServerException$json = {
  '1': 'InternalServerException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InternalServerException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List internalServerExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnRlcm5hbFNlcnZlckV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use kinesisParametersDescriptor instead')
const KinesisParameters$json = {
  '1': 'KinesisParameters',
  '2': [
    {'1': 'partitionkey', '3': 379379617, '4': 1, '5': 9, '10': 'partitionkey'},
  ],
};

/// Descriptor for `KinesisParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kinesisParametersDescriptor = $convert.base64Decode(
    'ChFLaW5lc2lzUGFyYW1ldGVycxImCgxwYXJ0aXRpb25rZXkYob/ztAEgASgJUgxwYXJ0aXRpb2'
    '5rZXk=');

@$core.Deprecated('Use listScheduleGroupsInputDescriptor instead')
const ListScheduleGroupsInput$json = {
  '1': 'ListScheduleGroupsInput',
  '2': [
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
  ],
};

/// Descriptor for `ListScheduleGroupsInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listScheduleGroupsInputDescriptor = $convert.base64Decode(
    'ChdMaXN0U2NoZWR1bGVHcm91cHNJbnB1dBIiCgptYXhyZXN1bHRzGLKom4MBIAEoBVIKbWF4cm'
    'VzdWx0cxIiCgpuYW1lcHJlZml4GJvzvKwBIAEoCVIKbmFtZXByZWZpeBIfCgluZXh0dG9rZW4Y'
    '/oS6ZyABKAlSCW5leHR0b2tlbg==');

@$core.Deprecated('Use listScheduleGroupsOutputDescriptor instead')
const ListScheduleGroupsOutput$json = {
  '1': 'ListScheduleGroupsOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'schedulegroups',
      '3': 136082885,
      '4': 3,
      '5': 11,
      '6': '.scheduler.ScheduleGroupSummary',
      '10': 'schedulegroups'
    },
  ],
};

/// Descriptor for `ListScheduleGroupsOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listScheduleGroupsOutputDescriptor = $convert.base64Decode(
    'ChhMaXN0U2NoZWR1bGVHcm91cHNPdXRwdXQSHwoJbmV4dHRva2VuGP6EumcgASgJUgluZXh0dG'
    '9rZW4SSgoOc2NoZWR1bGVncm91cHMYxevxQCADKAsyHy5zY2hlZHVsZXIuU2NoZWR1bGVHcm91'
    'cFN1bW1hcnlSDnNjaGVkdWxlZ3JvdXBz');

@$core.Deprecated('Use listSchedulesInputDescriptor instead')
const ListSchedulesInput$json = {
  '1': 'ListSchedulesInput',
  '2': [
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'maxresults', '3': 275174450, '4': 1, '5': 5, '10': 'maxresults'},
    {'1': 'nameprefix', '3': 361707931, '4': 1, '5': 9, '10': 'nameprefix'},
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `ListSchedulesInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSchedulesInputDescriptor = $convert.base64Decode(
    'ChJMaXN0U2NoZWR1bGVzSW5wdXQSIAoJZ3JvdXBuYW1lGMjKoKoBIAEoCVIJZ3JvdXBuYW1lEi'
    'IKCm1heHJlc3VsdHMYsqibgwEgASgFUgptYXhyZXN1bHRzEiIKCm5hbWVwcmVmaXgYm/O8rAEg'
    'ASgJUgpuYW1lcHJlZml4Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEhgKBXN0YX'
    'RlGJfJsu8BIAEoCVIFc3RhdGU=');

@$core.Deprecated('Use listSchedulesOutputDescriptor instead')
const ListSchedulesOutput$json = {
  '1': 'ListSchedulesOutput',
  '2': [
    {'1': 'nexttoken', '3': 216957566, '4': 1, '5': 9, '10': 'nexttoken'},
    {
      '1': 'schedules',
      '3': 18925646,
      '4': 3,
      '5': 11,
      '6': '.scheduler.ScheduleSummary',
      '10': 'schedules'
    },
  ],
};

/// Descriptor for `ListSchedulesOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listSchedulesOutputDescriptor = $convert.base64Decode(
    'ChNMaXN0U2NoZWR1bGVzT3V0cHV0Eh8KCW5leHR0b2tlbhj+hLpnIAEoCVIJbmV4dHRva2VuEj'
    'sKCXNjaGVkdWxlcxjOkIMJIAMoCzIaLnNjaGVkdWxlci5TY2hlZHVsZVN1bW1hcnlSCXNjaGVk'
    'dWxlcw==');

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

@$core.Deprecated('Use listTagsForResourceOutputDescriptor instead')
const ListTagsForResourceOutput$json = {
  '1': 'ListTagsForResourceOutput',
  '2': [
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.scheduler.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ListTagsForResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listTagsForResourceOutputDescriptor =
    $convert.base64Decode(
        'ChlMaXN0VGFnc0ZvclJlc291cmNlT3V0cHV0EiYKBHRhZ3MYwcH2tQEgAygLMg4uc2NoZWR1bG'
        'VyLlRhZ1IEdGFncw==');

@$core.Deprecated('Use networkConfigurationDescriptor instead')
const NetworkConfiguration$json = {
  '1': 'NetworkConfiguration',
  '2': [
    {
      '1': 'awsvpcconfiguration',
      '3': 223852630,
      '4': 1,
      '5': 11,
      '6': '.scheduler.AwsVpcConfiguration',
      '10': 'awsvpcconfiguration'
    },
  ],
};

/// Descriptor for `NetworkConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List networkConfigurationDescriptor = $convert.base64Decode(
    'ChROZXR3b3JrQ29uZmlndXJhdGlvbhJTChNhd3N2cGNjb25maWd1cmF0aW9uGNbw3mogASgLMh'
    '4uc2NoZWR1bGVyLkF3c1ZwY0NvbmZpZ3VyYXRpb25SE2F3c3ZwY2NvbmZpZ3VyYXRpb24=');

@$core.Deprecated('Use placementConstraintDescriptor instead')
const PlacementConstraint$json = {
  '1': 'PlacementConstraint',
  '2': [
    {'1': 'expression', '3': 253079532, '4': 1, '5': 9, '10': 'expression'},
    {'1': 'type', '3': 287830350, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `PlacementConstraint`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List placementConstraintDescriptor = $convert.base64Decode(
    'ChNQbGFjZW1lbnRDb25zdHJhaW50EiEKCmV4cHJlc3Npb24Y7N/WeCABKAlSCmV4cHJlc3Npb2'
    '4SFgoEdHlwZRjO4p+JASABKAlSBHR5cGU=');

@$core.Deprecated('Use placementStrategyDescriptor instead')
const PlacementStrategy$json = {
  '1': 'PlacementStrategy',
  '2': [
    {'1': 'field', '3': 125985384, '4': 1, '5': 9, '10': 'field'},
    {'1': 'type', '3': 287830350, '4': 1, '5': 9, '10': 'type'},
  ],
};

/// Descriptor for `PlacementStrategy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List placementStrategyDescriptor = $convert.base64Decode(
    'ChFQbGFjZW1lbnRTdHJhdGVneRIXCgVmaWVsZBjoxIk8IAEoCVIFZmllbGQSFgoEdHlwZRjO4p'
    '+JASABKAlSBHR5cGU=');

@$core.Deprecated('Use resourceNotFoundExceptionDescriptor instead')
const ResourceNotFoundException$json = {
  '1': 'ResourceNotFoundException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ResourceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List resourceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChlSZXNvdXJjZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use retryPolicyDescriptor instead')
const RetryPolicy$json = {
  '1': 'RetryPolicy',
  '2': [
    {
      '1': 'maximumeventageinseconds',
      '3': 393041563,
      '4': 1,
      '5': 5,
      '10': 'maximumeventageinseconds'
    },
    {
      '1': 'maximumretryattempts',
      '3': 112088128,
      '4': 1,
      '5': 5,
      '10': 'maximumretryattempts'
    },
  ],
};

/// Descriptor for `RetryPolicy`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List retryPolicyDescriptor = $convert.base64Decode(
    'CgtSZXRyeVBvbGljeRI+ChhtYXhpbXVtZXZlbnRhZ2VpbnNlY29uZHMYm621uwEgASgFUhhtYX'
    'hpbXVtZXZlbnRhZ2VpbnNlY29uZHMSNQoUbWF4aW11bXJldHJ5YXR0ZW1wdHMYwKi5NSABKAVS'
    'FG1heGltdW1yZXRyeWF0dGVtcHRz');

@$core.Deprecated('Use sageMakerPipelineParameterDescriptor instead')
const SageMakerPipelineParameter$json = {
  '1': 'SageMakerPipelineParameter',
  '2': [
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'value', '3': 289929579, '4': 1, '5': 9, '10': 'value'},
  ],
};

/// Descriptor for `SageMakerPipelineParameter`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sageMakerPipelineParameterDescriptor =
    $convert.base64Decode(
        'ChpTYWdlTWFrZXJQaXBlbGluZVBhcmFtZXRlchIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhgKBX'
        'ZhbHVlGOvyn4oBIAEoCVIFdmFsdWU=');

@$core.Deprecated('Use sageMakerPipelineParametersDescriptor instead')
const SageMakerPipelineParameters$json = {
  '1': 'SageMakerPipelineParameters',
  '2': [
    {
      '1': 'pipelineparameterlist',
      '3': 198270119,
      '4': 3,
      '5': 11,
      '6': '.scheduler.SageMakerPipelineParameter',
      '10': 'pipelineparameterlist'
    },
  ],
};

/// Descriptor for `SageMakerPipelineParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sageMakerPipelineParametersDescriptor =
    $convert.base64Decode(
        'ChtTYWdlTWFrZXJQaXBlbGluZVBhcmFtZXRlcnMSXgoVcGlwZWxpbmVwYXJhbWV0ZXJsaXN0GK'
        'e5xV4gAygLMiUuc2NoZWR1bGVyLlNhZ2VNYWtlclBpcGVsaW5lUGFyYW1ldGVyUhVwaXBlbGlu'
        'ZXBhcmFtZXRlcmxpc3Q=');

@$core.Deprecated('Use scheduleGroupSummaryDescriptor instead')
const ScheduleGroupSummary$json = {
  '1': 'ScheduleGroupSummary',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'lastmodificationdate',
      '3': 4750358,
      '4': 1,
      '5': 9,
      '10': 'lastmodificationdate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
  ],
};

/// Descriptor for `ScheduleGroupSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduleGroupSummaryDescriptor = $convert.base64Decode(
    'ChRTY2hlZHVsZUdyb3VwU3VtbWFyeRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SJgoMY3JlYXRpb2'
    '5kYXRlGOHYt4kBIAEoCVIMY3JlYXRpb25kYXRlEjUKFGxhc3Rtb2RpZmljYXRpb25kYXRlGJb4'
    'oQIgASgJUhRsYXN0bW9kaWZpY2F0aW9uZGF0ZRIVCgRuYW1lGIfmgX8gASgJUgRuYW1lEhgKBX'
    'N0YXRlGJfJsu8BIAEoCVIFc3RhdGU=');

@$core.Deprecated('Use scheduleSummaryDescriptor instead')
const ScheduleSummary$json = {
  '1': 'ScheduleSummary',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {
      '1': 'lastmodificationdate',
      '3': 4750358,
      '4': 1,
      '5': 9,
      '10': 'lastmodificationdate'
    },
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {
      '1': 'target',
      '3': 191361385,
      '4': 1,
      '5': 11,
      '6': '.scheduler.TargetSummary',
      '10': 'target'
    },
  ],
};

/// Descriptor for `ScheduleSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduleSummaryDescriptor = $convert.base64Decode(
    'Cg9TY2hlZHVsZVN1bW1hcnkSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEiYKDGNyZWF0aW9uZGF0ZR'
    'jh2LeJASABKAlSDGNyZWF0aW9uZGF0ZRIgCglncm91cG5hbWUYyMqgqgEgASgJUglncm91cG5h'
    'bWUSNQoUbGFzdG1vZGlmaWNhdGlvbmRhdGUYlvihAiABKAlSFGxhc3Rtb2RpZmljYXRpb25kYX'
    'RlEhUKBG5hbWUYh+aBfyABKAlSBG5hbWUSGAoFc3RhdGUYl8my7wEgASgJUgVzdGF0ZRIzCgZ0'
    'YXJnZXQY6eKfWyABKAsyGC5zY2hlZHVsZXIuVGFyZ2V0U3VtbWFyeVIGdGFyZ2V0');

@$core.Deprecated('Use serviceQuotaExceededExceptionDescriptor instead')
const ServiceQuotaExceededException$json = {
  '1': 'ServiceQuotaExceededException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ServiceQuotaExceededException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List serviceQuotaExceededExceptionDescriptor =
    $convert.base64Decode(
        'Ch1TZXJ2aWNlUXVvdGFFeGNlZWRlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGIWzu3AgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use sqsParametersDescriptor instead')
const SqsParameters$json = {
  '1': 'SqsParameters',
  '2': [
    {
      '1': 'messagegroupid',
      '3': 419537435,
      '4': 1,
      '5': 9,
      '10': 'messagegroupid'
    },
  ],
};

/// Descriptor for `SqsParameters`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List sqsParametersDescriptor = $convert.base64Decode(
    'Cg1TcXNQYXJhbWV0ZXJzEioKDm1lc3NhZ2Vncm91cGlkGJvEhsgBIAEoCVIObWVzc2FnZWdyb3'
    'VwaWQ=');

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
    {'1': 'resourcearn', '3': 364280877, '4': 1, '5': 9, '10': 'resourcearn'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.scheduler.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceInputDescriptor = $convert.base64Decode(
    'ChBUYWdSZXNvdXJjZUlucHV0EiQKC3Jlc291cmNlYXJuGK342a0BIAEoCVILcmVzb3VyY2Vhcm'
    '4SJgoEdGFncxjBwfa1ASADKAsyDi5zY2hlZHVsZXIuVGFnUgR0YWdz');

@$core.Deprecated('Use tagResourceOutputDescriptor instead')
const TagResourceOutput$json = {
  '1': 'TagResourceOutput',
};

/// Descriptor for `TagResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceOutputDescriptor =
    $convert.base64Decode('ChFUYWdSZXNvdXJjZU91dHB1dA==');

@$core.Deprecated('Use targetDescriptor instead')
const Target$json = {
  '1': 'Target',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'deadletterconfig',
      '3': 79786642,
      '4': 1,
      '5': 11,
      '6': '.scheduler.DeadLetterConfig',
      '10': 'deadletterconfig'
    },
    {
      '1': 'ecsparameters',
      '3': 501521183,
      '4': 1,
      '5': 11,
      '6': '.scheduler.EcsParameters',
      '10': 'ecsparameters'
    },
    {
      '1': 'eventbridgeparameters',
      '3': 285439471,
      '4': 1,
      '5': 11,
      '6': '.scheduler.EventBridgeParameters',
      '10': 'eventbridgeparameters'
    },
    {'1': 'input', '3': 529785116, '4': 1, '5': 9, '10': 'input'},
    {
      '1': 'kinesisparameters',
      '3': 70111902,
      '4': 1,
      '5': 11,
      '6': '.scheduler.KinesisParameters',
      '10': 'kinesisparameters'
    },
    {
      '1': 'retrypolicy',
      '3': 266827188,
      '4': 1,
      '5': 11,
      '6': '.scheduler.RetryPolicy',
      '10': 'retrypolicy'
    },
    {'1': 'rolearn', '3': 322567169, '4': 1, '5': 9, '10': 'rolearn'},
    {
      '1': 'sagemakerpipelineparameters',
      '3': 513800606,
      '4': 1,
      '5': 11,
      '6': '.scheduler.SageMakerPipelineParameters',
      '10': 'sagemakerpipelineparameters'
    },
    {
      '1': 'sqsparameters',
      '3': 91345999,
      '4': 1,
      '5': 11,
      '6': '.scheduler.SqsParameters',
      '10': 'sqsparameters'
    },
  ],
};

/// Descriptor for `Target`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetDescriptor = $convert.base64Decode(
    'CgZUYXJnZXQSFAoDYXJuGJ2b7b8BIAEoCVIDYXJuEkoKEGRlYWRsZXR0ZXJjb25maWcYkuWFJi'
    'ABKAsyGy5zY2hlZHVsZXIuRGVhZExldHRlckNvbmZpZ1IQZGVhZGxldHRlcmNvbmZpZxJCCg1l'
    'Y3NwYXJhbWV0ZXJzGJ+2ku8BIAEoCzIYLnNjaGVkdWxlci5FY3NQYXJhbWV0ZXJzUg1lY3NwYX'
    'JhbWV0ZXJzEloKFWV2ZW50YnJpZGdlcGFyYW1ldGVycxjv642IASABKAsyIC5zY2hlZHVsZXIu'
    'RXZlbnRCcmlkZ2VQYXJhbWV0ZXJzUhVldmVudGJyaWRnZXBhcmFtZXRlcnMSGAoFaW5wdXQYnM'
    'LP/AEgASgJUgVpbnB1dBJNChFraW5lc2lzcGFyYW1ldGVycxiepbchIAEoCzIcLnNjaGVkdWxl'
    'ci5LaW5lc2lzUGFyYW1ldGVyc1IRa2luZXNpc3BhcmFtZXRlcnMSOwoLcmV0cnlwb2xpY3kYtO'
    'udfyABKAsyFi5zY2hlZHVsZXIuUmV0cnlQb2xpY3lSC3JldHJ5cG9saWN5EhwKB3JvbGVhcm4Y'
    'gfjnmQEgASgJUgdyb2xlYXJuEmwKG3NhZ2VtYWtlcnBpcGVsaW5lcGFyYW1ldGVycxie8//0AS'
    'ABKAsyJi5zY2hlZHVsZXIuU2FnZU1ha2VyUGlwZWxpbmVQYXJhbWV0ZXJzUhtzYWdlbWFrZXJw'
    'aXBlbGluZXBhcmFtZXRlcnMSQQoNc3FzcGFyYW1ldGVycxjPqMcrIAEoCzIYLnNjaGVkdWxlci'
    '5TcXNQYXJhbWV0ZXJzUg1zcXNwYXJhbWV0ZXJz');

@$core.Deprecated('Use targetSummaryDescriptor instead')
const TargetSummary$json = {
  '1': 'TargetSummary',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
  ],
};

/// Descriptor for `TargetSummary`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List targetSummaryDescriptor = $convert
    .base64Decode('Cg1UYXJnZXRTdW1tYXJ5EhQKA2Fybhidm+2/ASABKAlSA2Fybg==');

@$core.Deprecated('Use throttlingExceptionDescriptor instead')
const ThrottlingException$json = {
  '1': 'ThrottlingException',
  '2': [
    {'1': 'message', '3': 235854213, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ThrottlingException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List throttlingExceptionDescriptor =
    $convert.base64Decode(
        'ChNUaHJvdHRsaW5nRXhjZXB0aW9uEhsKB21lc3NhZ2UYhbO7cCABKAlSB21lc3NhZ2U=');

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

@$core.Deprecated('Use untagResourceOutputDescriptor instead')
const UntagResourceOutput$json = {
  '1': 'UntagResourceOutput',
};

/// Descriptor for `UntagResourceOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceOutputDescriptor =
    $convert.base64Decode('ChNVbnRhZ1Jlc291cmNlT3V0cHV0');

@$core.Deprecated('Use updateScheduleInputDescriptor instead')
const UpdateScheduleInput$json = {
  '1': 'UpdateScheduleInput',
  '2': [
    {
      '1': 'actionaftercompletion',
      '3': 282350906,
      '4': 1,
      '5': 9,
      '10': 'actionaftercompletion'
    },
    {'1': 'clienttoken', '3': 137297356, '4': 1, '5': 9, '10': 'clienttoken'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'enddate', '3': 77486543, '4': 1, '5': 9, '10': 'enddate'},
    {
      '1': 'flexibletimewindow',
      '3': 2518952,
      '4': 1,
      '5': 11,
      '6': '.scheduler.FlexibleTimeWindow',
      '10': 'flexibletimewindow'
    },
    {'1': 'groupname', '3': 357049672, '4': 1, '5': 9, '10': 'groupname'},
    {'1': 'kmskeyarn', '3': 110041649, '4': 1, '5': 9, '10': 'kmskeyarn'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'scheduleexpression',
      '3': 446089471,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpression'
    },
    {
      '1': 'scheduleexpressiontimezone',
      '3': 400730326,
      '4': 1,
      '5': 9,
      '10': 'scheduleexpressiontimezone'
    },
    {'1': 'startdate', '3': 445135996, '4': 1, '5': 9, '10': 'startdate'},
    {'1': 'state', '3': 502047895, '4': 1, '5': 9, '10': 'state'},
    {
      '1': 'target',
      '3': 191361385,
      '4': 1,
      '5': 11,
      '6': '.scheduler.Target',
      '10': 'target'
    },
  ],
};

/// Descriptor for `UpdateScheduleInput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateScheduleInputDescriptor = $convert.base64Decode(
    'ChNVcGRhdGVTY2hlZHVsZUlucHV0EjgKFWFjdGlvbmFmdGVyY29tcGxldGlvbhi6qtGGASABKA'
    'lSFWFjdGlvbmFmdGVyY29tcGxldGlvbhIjCgtjbGllbnR0b2tlbhjM+7tBIAEoCVILY2xpZW50'
    'dG9rZW4SIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEhsKB2VuZGRhdGUYz7'
    'P5JCABKAlSB2VuZGRhdGUSUAoSZmxleGlibGV0aW1ld2luZG93GKjfmQEgASgLMh0uc2NoZWR1'
    'bGVyLkZsZXhpYmxlVGltZVdpbmRvd1ISZmxleGlibGV0aW1ld2luZG93EiAKCWdyb3VwbmFtZR'
    'jIyqCqASABKAlSCWdyb3VwbmFtZRIfCglrbXNrZXlhcm4YsbS8NCABKAlSCWttc2tleWFybhIV'
    'CgRuYW1lGIfmgX8gASgJUgRuYW1lEjIKEnNjaGVkdWxlZXhwcmVzc2lvbhj/kdvUASABKAlSEn'
    'NjaGVkdWxlZXhwcmVzc2lvbhJCChpzY2hlZHVsZWV4cHJlc3Npb250aW1lem9uZRjW0Yq/ASAB'
    'KAlSGnNjaGVkdWxlZXhwcmVzc2lvbnRpbWV6b25lEiAKCXN0YXJ0ZGF0ZRj8+KDUASABKAlSCX'
    'N0YXJ0ZGF0ZRIYCgVzdGF0ZRiXybLvASABKAlSBXN0YXRlEiwKBnRhcmdldBjp4p9bIAEoCzIR'
    'LnNjaGVkdWxlci5UYXJnZXRSBnRhcmdldA==');

@$core.Deprecated('Use updateScheduleOutputDescriptor instead')
const UpdateScheduleOutput$json = {
  '1': 'UpdateScheduleOutput',
  '2': [
    {'1': 'schedulearn', '3': 178445188, '4': 1, '5': 9, '10': 'schedulearn'},
  ],
};

/// Descriptor for `UpdateScheduleOutput`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateScheduleOutputDescriptor = $convert.base64Decode(
    'ChRVcGRhdGVTY2hlZHVsZU91dHB1dBIjCgtzY2hlZHVsZWFybhiEt4tVIAEoCVILc2NoZWR1bG'
    'Vhcm4=');

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
