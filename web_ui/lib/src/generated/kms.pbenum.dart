// This is a generated file - do not edit.
//
// Generated from kms.proto.

// @dart = 3.3

// ignore_for_file: annotate_overrides, camel_case_types, comment_references
// ignore_for_file: constant_identifier_names
// ignore_for_file: curly_braces_in_flow_control_structures
// ignore_for_file: deprecated_member_use_from_same_package, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_relative_imports

import 'dart:core' as $core;

import 'package:protobuf/protobuf.dart' as $pb;

class AlgorithmSpec extends $pb.ProtobufEnum {
  static const AlgorithmSpec ALGORITHM_SPEC_SM2PKE =
      AlgorithmSpec._(0, _omitEnumNames ? '' : 'ALGORITHM_SPEC_SM2PKE');
  static const AlgorithmSpec ALGORITHM_SPEC_RSAES_PKCS1_V1_5 = AlgorithmSpec._(
      1, _omitEnumNames ? '' : 'ALGORITHM_SPEC_RSAES_PKCS1_V1_5');
  static const AlgorithmSpec ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_256 =
      AlgorithmSpec._(
          2, _omitEnumNames ? '' : 'ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_256');
  static const AlgorithmSpec ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_1 =
      AlgorithmSpec._(
          3, _omitEnumNames ? '' : 'ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_1');
  static const AlgorithmSpec ALGORITHM_SPEC_RSAES_OAEP_SHA_1 = AlgorithmSpec._(
      4, _omitEnumNames ? '' : 'ALGORITHM_SPEC_RSAES_OAEP_SHA_1');
  static const AlgorithmSpec ALGORITHM_SPEC_RSAES_OAEP_SHA_256 =
      AlgorithmSpec._(
          5, _omitEnumNames ? '' : 'ALGORITHM_SPEC_RSAES_OAEP_SHA_256');

  static const $core.List<AlgorithmSpec> values = <AlgorithmSpec>[
    ALGORITHM_SPEC_SM2PKE,
    ALGORITHM_SPEC_RSAES_PKCS1_V1_5,
    ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_256,
    ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_1,
    ALGORITHM_SPEC_RSAES_OAEP_SHA_1,
    ALGORITHM_SPEC_RSAES_OAEP_SHA_256,
  ];

  static final $core.List<AlgorithmSpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 5);
  static AlgorithmSpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const AlgorithmSpec._(super.value, super.name);
}

class ConnectionErrorCodeType extends $pb.ProtobufEnum {
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_SUBNET_NOT_FOUND = ConnectionErrorCodeType._(0,
          _omitEnumNames ? '' : 'CONNECTION_ERROR_CODE_TYPE_SUBNET_NOT_FOUND');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_TIMED_OUT =
      ConnectionErrorCodeType._(
          1,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_TIMED_OUT');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_NETWORK_ERRORS = ConnectionErrorCodeType._(
          2, _omitEnumNames ? '' : 'CONNECTION_ERROR_CODE_TYPE_NETWORK_ERRORS');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_INVALID_CREDENTIALS =
      ConnectionErrorCodeType._(
          3,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_INVALID_CREDENTIALS');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_USER_LOGGED_IN = ConnectionErrorCodeType._(
          4, _omitEnumNames ? '' : 'CONNECTION_ERROR_CODE_TYPE_USER_LOGGED_IN');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_NOT_REACHABLE =
      ConnectionErrorCodeType._(
          5,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_NOT_REACHABLE');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_CLUSTER_NOT_FOUND = ConnectionErrorCodeType._(
          6,
          _omitEnumNames ? '' : 'CONNECTION_ERROR_CODE_TYPE_CLUSTER_NOT_FOUND');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_INVALID_CONFIGURATION =
      ConnectionErrorCodeType._(
          7,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_INVALID_CONFIGURATION');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_USER_NOT_FOUND = ConnectionErrorCodeType._(
          8, _omitEnumNames ? '' : 'CONNECTION_ERROR_CODE_TYPE_USER_NOT_FOUND');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_TLS_CONFIGURATION =
      ConnectionErrorCodeType._(
          9,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_TLS_CONFIGURATION');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_CLOUDHSM_HSMS =
      ConnectionErrorCodeType._(
          10,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_CLOUDHSM_HSMS');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_NOT_FOUND =
      ConnectionErrorCodeType._(
          11,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_NOT_FOUND');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_RESPONSE =
      ConnectionErrorCodeType._(
          12,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_RESPONSE');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_USER_LOCKED_OUT = ConnectionErrorCodeType._(13,
          _omitEnumNames ? '' : 'CONNECTION_ERROR_CODE_TYPE_USER_LOCKED_OUT');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_CONFIGURATION =
      ConnectionErrorCodeType._(
          14,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_CONFIGURATION');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_FREE_ADDRESSES_IN_SUBNET =
      ConnectionErrorCodeType._(
          15,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_FREE_ADDRESSES_IN_SUBNET');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_INTERNAL_ERROR = ConnectionErrorCodeType._(16,
          _omitEnumNames ? '' : 'CONNECTION_ERROR_CODE_TYPE_INTERNAL_ERROR');
  static const ConnectionErrorCodeType
      CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_ACCESS_DENIED =
      ConnectionErrorCodeType._(
          17,
          _omitEnumNames
              ? ''
              : 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_ACCESS_DENIED');

  static const $core.List<ConnectionErrorCodeType> values =
      <ConnectionErrorCodeType>[
    CONNECTION_ERROR_CODE_TYPE_SUBNET_NOT_FOUND,
    CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_TIMED_OUT,
    CONNECTION_ERROR_CODE_TYPE_NETWORK_ERRORS,
    CONNECTION_ERROR_CODE_TYPE_INVALID_CREDENTIALS,
    CONNECTION_ERROR_CODE_TYPE_USER_LOGGED_IN,
    CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_NOT_REACHABLE,
    CONNECTION_ERROR_CODE_TYPE_CLUSTER_NOT_FOUND,
    CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_INVALID_CONFIGURATION,
    CONNECTION_ERROR_CODE_TYPE_USER_NOT_FOUND,
    CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_TLS_CONFIGURATION,
    CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_CLOUDHSM_HSMS,
    CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_NOT_FOUND,
    CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_RESPONSE,
    CONNECTION_ERROR_CODE_TYPE_USER_LOCKED_OUT,
    CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_CONFIGURATION,
    CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_FREE_ADDRESSES_IN_SUBNET,
    CONNECTION_ERROR_CODE_TYPE_INTERNAL_ERROR,
    CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_ACCESS_DENIED,
  ];

  static final $core.List<ConnectionErrorCodeType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 17);
  static ConnectionErrorCodeType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionErrorCodeType._(super.value, super.name);
}

class ConnectionStateType extends $pb.ProtobufEnum {
  static const ConnectionStateType CONNECTION_STATE_TYPE_CONNECTING =
      ConnectionStateType._(
          0, _omitEnumNames ? '' : 'CONNECTION_STATE_TYPE_CONNECTING');
  static const ConnectionStateType CONNECTION_STATE_TYPE_DISCONNECTING =
      ConnectionStateType._(
          1, _omitEnumNames ? '' : 'CONNECTION_STATE_TYPE_DISCONNECTING');
  static const ConnectionStateType CONNECTION_STATE_TYPE_DISCONNECTED =
      ConnectionStateType._(
          2, _omitEnumNames ? '' : 'CONNECTION_STATE_TYPE_DISCONNECTED');
  static const ConnectionStateType CONNECTION_STATE_TYPE_CONNECTED =
      ConnectionStateType._(
          3, _omitEnumNames ? '' : 'CONNECTION_STATE_TYPE_CONNECTED');
  static const ConnectionStateType CONNECTION_STATE_TYPE_FAILED =
      ConnectionStateType._(
          4, _omitEnumNames ? '' : 'CONNECTION_STATE_TYPE_FAILED');

  static const $core.List<ConnectionStateType> values = <ConnectionStateType>[
    CONNECTION_STATE_TYPE_CONNECTING,
    CONNECTION_STATE_TYPE_DISCONNECTING,
    CONNECTION_STATE_TYPE_DISCONNECTED,
    CONNECTION_STATE_TYPE_CONNECTED,
    CONNECTION_STATE_TYPE_FAILED,
  ];

  static final $core.List<ConnectionStateType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 4);
  static ConnectionStateType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ConnectionStateType._(super.value, super.name);
}

class CustomKeyStoreType extends $pb.ProtobufEnum {
  static const CustomKeyStoreType CUSTOM_KEY_STORE_TYPE_EXTERNAL_KEY_STORE =
      CustomKeyStoreType._(
          0, _omitEnumNames ? '' : 'CUSTOM_KEY_STORE_TYPE_EXTERNAL_KEY_STORE');
  static const CustomKeyStoreType CUSTOM_KEY_STORE_TYPE_AWS_CLOUDHSM =
      CustomKeyStoreType._(
          1, _omitEnumNames ? '' : 'CUSTOM_KEY_STORE_TYPE_AWS_CLOUDHSM');

  static const $core.List<CustomKeyStoreType> values = <CustomKeyStoreType>[
    CUSTOM_KEY_STORE_TYPE_EXTERNAL_KEY_STORE,
    CUSTOM_KEY_STORE_TYPE_AWS_CLOUDHSM,
  ];

  static final $core.List<CustomKeyStoreType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static CustomKeyStoreType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CustomKeyStoreType._(super.value, super.name);
}

class CustomerMasterKeySpec extends $pb.ProtobufEnum {
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_HMAC_256 =
      CustomerMasterKeySpec._(
          0, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_HMAC_256');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P384 =
      CustomerMasterKeySpec._(
          1, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P384');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_HMAC_384 =
      CustomerMasterKeySpec._(
          2, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_HMAC_384');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_SM2 =
      CustomerMasterKeySpec._(
          3, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_SM2');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_ECC_SECG_P256K1 =
      CustomerMasterKeySpec._(
          4, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_ECC_SECG_P256K1');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P521 =
      CustomerMasterKeySpec._(
          5, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P521');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_HMAC_512 =
      CustomerMasterKeySpec._(
          6, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_HMAC_512');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_HMAC_224 =
      CustomerMasterKeySpec._(
          7, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_HMAC_224');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_RSA_2048 =
      CustomerMasterKeySpec._(
          8, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_RSA_2048');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_RSA_3072 =
      CustomerMasterKeySpec._(
          9, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_RSA_3072');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_RSA_4096 =
      CustomerMasterKeySpec._(
          10, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_RSA_4096');
  static const CustomerMasterKeySpec
      CUSTOMER_MASTER_KEY_SPEC_SYMMETRIC_DEFAULT = CustomerMasterKeySpec._(11,
          _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_SYMMETRIC_DEFAULT');
  static const CustomerMasterKeySpec CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P256 =
      CustomerMasterKeySpec._(
          12, _omitEnumNames ? '' : 'CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P256');

  static const $core.List<CustomerMasterKeySpec> values =
      <CustomerMasterKeySpec>[
    CUSTOMER_MASTER_KEY_SPEC_HMAC_256,
    CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P384,
    CUSTOMER_MASTER_KEY_SPEC_HMAC_384,
    CUSTOMER_MASTER_KEY_SPEC_SM2,
    CUSTOMER_MASTER_KEY_SPEC_ECC_SECG_P256K1,
    CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P521,
    CUSTOMER_MASTER_KEY_SPEC_HMAC_512,
    CUSTOMER_MASTER_KEY_SPEC_HMAC_224,
    CUSTOMER_MASTER_KEY_SPEC_RSA_2048,
    CUSTOMER_MASTER_KEY_SPEC_RSA_3072,
    CUSTOMER_MASTER_KEY_SPEC_RSA_4096,
    CUSTOMER_MASTER_KEY_SPEC_SYMMETRIC_DEFAULT,
    CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P256,
  ];

  static final $core.List<CustomerMasterKeySpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 12);
  static CustomerMasterKeySpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const CustomerMasterKeySpec._(super.value, super.name);
}

class DataKeyPairSpec extends $pb.ProtobufEnum {
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_ECC_NIST_EDWARDS25519 =
      DataKeyPairSpec._(
          0, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_ECC_NIST_EDWARDS25519');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_ECC_NIST_P384 =
      DataKeyPairSpec._(
          1, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_ECC_NIST_P384');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_SM2 =
      DataKeyPairSpec._(2, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_SM2');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_ECC_SECG_P256K1 =
      DataKeyPairSpec._(
          3, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_ECC_SECG_P256K1');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_ECC_NIST_P521 =
      DataKeyPairSpec._(
          4, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_ECC_NIST_P521');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_RSA_2048 =
      DataKeyPairSpec._(5, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_RSA_2048');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_RSA_3072 =
      DataKeyPairSpec._(6, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_RSA_3072');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_RSA_4096 =
      DataKeyPairSpec._(7, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_RSA_4096');
  static const DataKeyPairSpec DATA_KEY_PAIR_SPEC_ECC_NIST_P256 =
      DataKeyPairSpec._(
          8, _omitEnumNames ? '' : 'DATA_KEY_PAIR_SPEC_ECC_NIST_P256');

  static const $core.List<DataKeyPairSpec> values = <DataKeyPairSpec>[
    DATA_KEY_PAIR_SPEC_ECC_NIST_EDWARDS25519,
    DATA_KEY_PAIR_SPEC_ECC_NIST_P384,
    DATA_KEY_PAIR_SPEC_SM2,
    DATA_KEY_PAIR_SPEC_ECC_SECG_P256K1,
    DATA_KEY_PAIR_SPEC_ECC_NIST_P521,
    DATA_KEY_PAIR_SPEC_RSA_2048,
    DATA_KEY_PAIR_SPEC_RSA_3072,
    DATA_KEY_PAIR_SPEC_RSA_4096,
    DATA_KEY_PAIR_SPEC_ECC_NIST_P256,
  ];

  static final $core.List<DataKeyPairSpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 8);
  static DataKeyPairSpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataKeyPairSpec._(super.value, super.name);
}

class DataKeySpec extends $pb.ProtobufEnum {
  static const DataKeySpec DATA_KEY_SPEC_AES_128 =
      DataKeySpec._(0, _omitEnumNames ? '' : 'DATA_KEY_SPEC_AES_128');
  static const DataKeySpec DATA_KEY_SPEC_AES_256 =
      DataKeySpec._(1, _omitEnumNames ? '' : 'DATA_KEY_SPEC_AES_256');

  static const $core.List<DataKeySpec> values = <DataKeySpec>[
    DATA_KEY_SPEC_AES_128,
    DATA_KEY_SPEC_AES_256,
  ];

  static final $core.List<DataKeySpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static DataKeySpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DataKeySpec._(super.value, super.name);
}

class DryRunModifierType extends $pb.ProtobufEnum {
  static const DryRunModifierType DRY_RUN_MODIFIER_TYPE_IGNORE_CIPHERTEXT =
      DryRunModifierType._(
          0, _omitEnumNames ? '' : 'DRY_RUN_MODIFIER_TYPE_IGNORE_CIPHERTEXT');

  static const $core.List<DryRunModifierType> values = <DryRunModifierType>[
    DRY_RUN_MODIFIER_TYPE_IGNORE_CIPHERTEXT,
  ];

  static final $core.List<DryRunModifierType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static DryRunModifierType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const DryRunModifierType._(super.value, super.name);
}

class EncryptionAlgorithmSpec extends $pb.ProtobufEnum {
  static const EncryptionAlgorithmSpec ENCRYPTION_ALGORITHM_SPEC_SM2PKE =
      EncryptionAlgorithmSpec._(
          0, _omitEnumNames ? '' : 'ENCRYPTION_ALGORITHM_SPEC_SM2PKE');
  static const EncryptionAlgorithmSpec
      ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_1 = EncryptionAlgorithmSpec._(1,
          _omitEnumNames ? '' : 'ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_1');
  static const EncryptionAlgorithmSpec
      ENCRYPTION_ALGORITHM_SPEC_SYMMETRIC_DEFAULT = EncryptionAlgorithmSpec._(2,
          _omitEnumNames ? '' : 'ENCRYPTION_ALGORITHM_SPEC_SYMMETRIC_DEFAULT');
  static const EncryptionAlgorithmSpec
      ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_256 = EncryptionAlgorithmSpec._(
          3,
          _omitEnumNames ? '' : 'ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_256');

  static const $core.List<EncryptionAlgorithmSpec> values =
      <EncryptionAlgorithmSpec>[
    ENCRYPTION_ALGORITHM_SPEC_SM2PKE,
    ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_1,
    ENCRYPTION_ALGORITHM_SPEC_SYMMETRIC_DEFAULT,
    ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_256,
  ];

  static final $core.List<EncryptionAlgorithmSpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static EncryptionAlgorithmSpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const EncryptionAlgorithmSpec._(super.value, super.name);
}

class ExpirationModelType extends $pb.ProtobufEnum {
  static const ExpirationModelType
      EXPIRATION_MODEL_TYPE_KEY_MATERIAL_DOES_NOT_EXPIRE =
      ExpirationModelType._(
          0,
          _omitEnumNames
              ? ''
              : 'EXPIRATION_MODEL_TYPE_KEY_MATERIAL_DOES_NOT_EXPIRE');
  static const ExpirationModelType EXPIRATION_MODEL_TYPE_KEY_MATERIAL_EXPIRES =
      ExpirationModelType._(1,
          _omitEnumNames ? '' : 'EXPIRATION_MODEL_TYPE_KEY_MATERIAL_EXPIRES');

  static const $core.List<ExpirationModelType> values = <ExpirationModelType>[
    EXPIRATION_MODEL_TYPE_KEY_MATERIAL_DOES_NOT_EXPIRE,
    EXPIRATION_MODEL_TYPE_KEY_MATERIAL_EXPIRES,
  ];

  static final $core.List<ExpirationModelType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ExpirationModelType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ExpirationModelType._(super.value, super.name);
}

class GrantOperation extends $pb.ProtobufEnum {
  static const GrantOperation GRANT_OPERATION_GENERATEDATAKEYPAIR =
      GrantOperation._(
          0, _omitEnumNames ? '' : 'GRANT_OPERATION_GENERATEDATAKEYPAIR');
  static const GrantOperation GRANT_OPERATION_RETIREGRANT =
      GrantOperation._(1, _omitEnumNames ? '' : 'GRANT_OPERATION_RETIREGRANT');
  static const GrantOperation GRANT_OPERATION_DESCRIBEKEY =
      GrantOperation._(2, _omitEnumNames ? '' : 'GRANT_OPERATION_DESCRIBEKEY');
  static const GrantOperation GRANT_OPERATION_VERIFY =
      GrantOperation._(3, _omitEnumNames ? '' : 'GRANT_OPERATION_VERIFY');
  static const GrantOperation GRANT_OPERATION_ENCRYPT =
      GrantOperation._(4, _omitEnumNames ? '' : 'GRANT_OPERATION_ENCRYPT');
  static const GrantOperation GRANT_OPERATION_GENERATEDATAKEYWITHOUTPLAINTEXT =
      GrantOperation._(
          5,
          _omitEnumNames
              ? ''
              : 'GRANT_OPERATION_GENERATEDATAKEYWITHOUTPLAINTEXT');
  static const GrantOperation
      GRANT_OPERATION_GENERATEDATAKEYPAIRWITHOUTPLAINTEXT = GrantOperation._(
          6,
          _omitEnumNames
              ? ''
              : 'GRANT_OPERATION_GENERATEDATAKEYPAIRWITHOUTPLAINTEXT');
  static const GrantOperation GRANT_OPERATION_SIGN =
      GrantOperation._(7, _omitEnumNames ? '' : 'GRANT_OPERATION_SIGN');
  static const GrantOperation GRANT_OPERATION_DERIVESHAREDSECRET =
      GrantOperation._(
          8, _omitEnumNames ? '' : 'GRANT_OPERATION_DERIVESHAREDSECRET');
  static const GrantOperation GRANT_OPERATION_VERIFYMAC =
      GrantOperation._(9, _omitEnumNames ? '' : 'GRANT_OPERATION_VERIFYMAC');
  static const GrantOperation GRANT_OPERATION_GETPUBLICKEY = GrantOperation._(
      10, _omitEnumNames ? '' : 'GRANT_OPERATION_GETPUBLICKEY');
  static const GrantOperation GRANT_OPERATION_CREATEGRANT =
      GrantOperation._(11, _omitEnumNames ? '' : 'GRANT_OPERATION_CREATEGRANT');
  static const GrantOperation GRANT_OPERATION_REENCRYPTTO =
      GrantOperation._(12, _omitEnumNames ? '' : 'GRANT_OPERATION_REENCRYPTTO');
  static const GrantOperation GRANT_OPERATION_REENCRYPTFROM = GrantOperation._(
      13, _omitEnumNames ? '' : 'GRANT_OPERATION_REENCRYPTFROM');
  static const GrantOperation GRANT_OPERATION_GENERATEDATAKEY =
      GrantOperation._(
          14, _omitEnumNames ? '' : 'GRANT_OPERATION_GENERATEDATAKEY');
  static const GrantOperation GRANT_OPERATION_DECRYPT =
      GrantOperation._(15, _omitEnumNames ? '' : 'GRANT_OPERATION_DECRYPT');
  static const GrantOperation GRANT_OPERATION_GENERATEMAC =
      GrantOperation._(16, _omitEnumNames ? '' : 'GRANT_OPERATION_GENERATEMAC');

  static const $core.List<GrantOperation> values = <GrantOperation>[
    GRANT_OPERATION_GENERATEDATAKEYPAIR,
    GRANT_OPERATION_RETIREGRANT,
    GRANT_OPERATION_DESCRIBEKEY,
    GRANT_OPERATION_VERIFY,
    GRANT_OPERATION_ENCRYPT,
    GRANT_OPERATION_GENERATEDATAKEYWITHOUTPLAINTEXT,
    GRANT_OPERATION_GENERATEDATAKEYPAIRWITHOUTPLAINTEXT,
    GRANT_OPERATION_SIGN,
    GRANT_OPERATION_DERIVESHAREDSECRET,
    GRANT_OPERATION_VERIFYMAC,
    GRANT_OPERATION_GETPUBLICKEY,
    GRANT_OPERATION_CREATEGRANT,
    GRANT_OPERATION_REENCRYPTTO,
    GRANT_OPERATION_REENCRYPTFROM,
    GRANT_OPERATION_GENERATEDATAKEY,
    GRANT_OPERATION_DECRYPT,
    GRANT_OPERATION_GENERATEMAC,
  ];

  static final $core.List<GrantOperation?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 16);
  static GrantOperation? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const GrantOperation._(super.value, super.name);
}

class ImportState extends $pb.ProtobufEnum {
  static const ImportState IMPORT_STATE_IMPORTED =
      ImportState._(0, _omitEnumNames ? '' : 'IMPORT_STATE_IMPORTED');
  static const ImportState IMPORT_STATE_PENDING_IMPORT =
      ImportState._(1, _omitEnumNames ? '' : 'IMPORT_STATE_PENDING_IMPORT');

  static const $core.List<ImportState> values = <ImportState>[
    IMPORT_STATE_IMPORTED,
    IMPORT_STATE_PENDING_IMPORT,
  ];

  static final $core.List<ImportState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ImportState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportState._(super.value, super.name);
}

class ImportType extends $pb.ProtobufEnum {
  static const ImportType IMPORT_TYPE_EXISTING_KEY_MATERIAL = ImportType._(
      0, _omitEnumNames ? '' : 'IMPORT_TYPE_EXISTING_KEY_MATERIAL');
  static const ImportType IMPORT_TYPE_NEW_KEY_MATERIAL =
      ImportType._(1, _omitEnumNames ? '' : 'IMPORT_TYPE_NEW_KEY_MATERIAL');

  static const $core.List<ImportType> values = <ImportType>[
    IMPORT_TYPE_EXISTING_KEY_MATERIAL,
    IMPORT_TYPE_NEW_KEY_MATERIAL,
  ];

  static final $core.List<ImportType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static ImportType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const ImportType._(super.value, super.name);
}

class IncludeKeyMaterial extends $pb.ProtobufEnum {
  static const IncludeKeyMaterial INCLUDE_KEY_MATERIAL_ALL_KEY_MATERIAL =
      IncludeKeyMaterial._(
          0, _omitEnumNames ? '' : 'INCLUDE_KEY_MATERIAL_ALL_KEY_MATERIAL');
  static const IncludeKeyMaterial INCLUDE_KEY_MATERIAL_ROTATIONS_ONLY =
      IncludeKeyMaterial._(
          1, _omitEnumNames ? '' : 'INCLUDE_KEY_MATERIAL_ROTATIONS_ONLY');

  static const $core.List<IncludeKeyMaterial> values = <IncludeKeyMaterial>[
    INCLUDE_KEY_MATERIAL_ALL_KEY_MATERIAL,
    INCLUDE_KEY_MATERIAL_ROTATIONS_ONLY,
  ];

  static final $core.List<IncludeKeyMaterial?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static IncludeKeyMaterial? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const IncludeKeyMaterial._(super.value, super.name);
}

class KeyAgreementAlgorithmSpec extends $pb.ProtobufEnum {
  static const KeyAgreementAlgorithmSpec KEY_AGREEMENT_ALGORITHM_SPEC_ECDH =
      KeyAgreementAlgorithmSpec._(
          0, _omitEnumNames ? '' : 'KEY_AGREEMENT_ALGORITHM_SPEC_ECDH');

  static const $core.List<KeyAgreementAlgorithmSpec> values =
      <KeyAgreementAlgorithmSpec>[
    KEY_AGREEMENT_ALGORITHM_SPEC_ECDH,
  ];

  static final $core.List<KeyAgreementAlgorithmSpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static KeyAgreementAlgorithmSpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyAgreementAlgorithmSpec._(super.value, super.name);
}

class KeyEncryptionMechanism extends $pb.ProtobufEnum {
  static const KeyEncryptionMechanism
      KEY_ENCRYPTION_MECHANISM_RSAES_OAEP_SHA_256 = KeyEncryptionMechanism._(0,
          _omitEnumNames ? '' : 'KEY_ENCRYPTION_MECHANISM_RSAES_OAEP_SHA_256');

  static const $core.List<KeyEncryptionMechanism> values =
      <KeyEncryptionMechanism>[
    KEY_ENCRYPTION_MECHANISM_RSAES_OAEP_SHA_256,
  ];

  static final $core.List<KeyEncryptionMechanism?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 0);
  static KeyEncryptionMechanism? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyEncryptionMechanism._(super.value, super.name);
}

class KeyManagerType extends $pb.ProtobufEnum {
  static const KeyManagerType KEY_MANAGER_TYPE_CUSTOMER =
      KeyManagerType._(0, _omitEnumNames ? '' : 'KEY_MANAGER_TYPE_CUSTOMER');
  static const KeyManagerType KEY_MANAGER_TYPE_AWS =
      KeyManagerType._(1, _omitEnumNames ? '' : 'KEY_MANAGER_TYPE_AWS');

  static const $core.List<KeyManagerType> values = <KeyManagerType>[
    KEY_MANAGER_TYPE_CUSTOMER,
    KEY_MANAGER_TYPE_AWS,
  ];

  static final $core.List<KeyManagerType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static KeyManagerType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyManagerType._(super.value, super.name);
}

class KeyMaterialState extends $pb.ProtobufEnum {
  static const KeyMaterialState KEY_MATERIAL_STATE_PENDING_ROTATION =
      KeyMaterialState._(
          0, _omitEnumNames ? '' : 'KEY_MATERIAL_STATE_PENDING_ROTATION');
  static const KeyMaterialState
      KEY_MATERIAL_STATE_PENDING_MULTI_REGION_IMPORT_AND_ROTATION =
      KeyMaterialState._(
          1,
          _omitEnumNames
              ? ''
              : 'KEY_MATERIAL_STATE_PENDING_MULTI_REGION_IMPORT_AND_ROTATION');
  static const KeyMaterialState KEY_MATERIAL_STATE_NON_CURRENT =
      KeyMaterialState._(
          2, _omitEnumNames ? '' : 'KEY_MATERIAL_STATE_NON_CURRENT');
  static const KeyMaterialState KEY_MATERIAL_STATE_CURRENT =
      KeyMaterialState._(3, _omitEnumNames ? '' : 'KEY_MATERIAL_STATE_CURRENT');

  static const $core.List<KeyMaterialState> values = <KeyMaterialState>[
    KEY_MATERIAL_STATE_PENDING_ROTATION,
    KEY_MATERIAL_STATE_PENDING_MULTI_REGION_IMPORT_AND_ROTATION,
    KEY_MATERIAL_STATE_NON_CURRENT,
    KEY_MATERIAL_STATE_CURRENT,
  ];

  static final $core.List<KeyMaterialState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static KeyMaterialState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyMaterialState._(super.value, super.name);
}

class KeySpec extends $pb.ProtobufEnum {
  static const KeySpec KEY_SPEC_HMAC_256 =
      KeySpec._(0, _omitEnumNames ? '' : 'KEY_SPEC_HMAC_256');
  static const KeySpec KEY_SPEC_ML_DSA_87 =
      KeySpec._(1, _omitEnumNames ? '' : 'KEY_SPEC_ML_DSA_87');
  static const KeySpec KEY_SPEC_ECC_NIST_EDWARDS25519 =
      KeySpec._(2, _omitEnumNames ? '' : 'KEY_SPEC_ECC_NIST_EDWARDS25519');
  static const KeySpec KEY_SPEC_ECC_NIST_P384 =
      KeySpec._(3, _omitEnumNames ? '' : 'KEY_SPEC_ECC_NIST_P384');
  static const KeySpec KEY_SPEC_HMAC_384 =
      KeySpec._(4, _omitEnumNames ? '' : 'KEY_SPEC_HMAC_384');
  static const KeySpec KEY_SPEC_SM2 =
      KeySpec._(5, _omitEnumNames ? '' : 'KEY_SPEC_SM2');
  static const KeySpec KEY_SPEC_ECC_SECG_P256K1 =
      KeySpec._(6, _omitEnumNames ? '' : 'KEY_SPEC_ECC_SECG_P256K1');
  static const KeySpec KEY_SPEC_ECC_NIST_P521 =
      KeySpec._(7, _omitEnumNames ? '' : 'KEY_SPEC_ECC_NIST_P521');
  static const KeySpec KEY_SPEC_HMAC_512 =
      KeySpec._(8, _omitEnumNames ? '' : 'KEY_SPEC_HMAC_512');
  static const KeySpec KEY_SPEC_HMAC_224 =
      KeySpec._(9, _omitEnumNames ? '' : 'KEY_SPEC_HMAC_224');
  static const KeySpec KEY_SPEC_ML_DSA_44 =
      KeySpec._(10, _omitEnumNames ? '' : 'KEY_SPEC_ML_DSA_44');
  static const KeySpec KEY_SPEC_RSA_2048 =
      KeySpec._(11, _omitEnumNames ? '' : 'KEY_SPEC_RSA_2048');
  static const KeySpec KEY_SPEC_RSA_3072 =
      KeySpec._(12, _omitEnumNames ? '' : 'KEY_SPEC_RSA_3072');
  static const KeySpec KEY_SPEC_ML_DSA_65 =
      KeySpec._(13, _omitEnumNames ? '' : 'KEY_SPEC_ML_DSA_65');
  static const KeySpec KEY_SPEC_RSA_4096 =
      KeySpec._(14, _omitEnumNames ? '' : 'KEY_SPEC_RSA_4096');
  static const KeySpec KEY_SPEC_SYMMETRIC_DEFAULT =
      KeySpec._(15, _omitEnumNames ? '' : 'KEY_SPEC_SYMMETRIC_DEFAULT');
  static const KeySpec KEY_SPEC_ECC_NIST_P256 =
      KeySpec._(16, _omitEnumNames ? '' : 'KEY_SPEC_ECC_NIST_P256');

  static const $core.List<KeySpec> values = <KeySpec>[
    KEY_SPEC_HMAC_256,
    KEY_SPEC_ML_DSA_87,
    KEY_SPEC_ECC_NIST_EDWARDS25519,
    KEY_SPEC_ECC_NIST_P384,
    KEY_SPEC_HMAC_384,
    KEY_SPEC_SM2,
    KEY_SPEC_ECC_SECG_P256K1,
    KEY_SPEC_ECC_NIST_P521,
    KEY_SPEC_HMAC_512,
    KEY_SPEC_HMAC_224,
    KEY_SPEC_ML_DSA_44,
    KEY_SPEC_RSA_2048,
    KEY_SPEC_RSA_3072,
    KEY_SPEC_ML_DSA_65,
    KEY_SPEC_RSA_4096,
    KEY_SPEC_SYMMETRIC_DEFAULT,
    KEY_SPEC_ECC_NIST_P256,
  ];

  static final $core.List<KeySpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 16);
  static KeySpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeySpec._(super.value, super.name);
}

class KeyState extends $pb.ProtobufEnum {
  static const KeyState KEY_STATE_PENDINGREPLICADELETION =
      KeyState._(0, _omitEnumNames ? '' : 'KEY_STATE_PENDINGREPLICADELETION');
  static const KeyState KEY_STATE_UPDATING =
      KeyState._(1, _omitEnumNames ? '' : 'KEY_STATE_UPDATING');
  static const KeyState KEY_STATE_PENDINGIMPORT =
      KeyState._(2, _omitEnumNames ? '' : 'KEY_STATE_PENDINGIMPORT');
  static const KeyState KEY_STATE_PENDINGDELETION =
      KeyState._(3, _omitEnumNames ? '' : 'KEY_STATE_PENDINGDELETION');
  static const KeyState KEY_STATE_CREATING =
      KeyState._(4, _omitEnumNames ? '' : 'KEY_STATE_CREATING');
  static const KeyState KEY_STATE_DISABLED =
      KeyState._(5, _omitEnumNames ? '' : 'KEY_STATE_DISABLED');
  static const KeyState KEY_STATE_ENABLED =
      KeyState._(6, _omitEnumNames ? '' : 'KEY_STATE_ENABLED');
  static const KeyState KEY_STATE_UNAVAILABLE =
      KeyState._(7, _omitEnumNames ? '' : 'KEY_STATE_UNAVAILABLE');

  static const $core.List<KeyState> values = <KeyState>[
    KEY_STATE_PENDINGREPLICADELETION,
    KEY_STATE_UPDATING,
    KEY_STATE_PENDINGIMPORT,
    KEY_STATE_PENDINGDELETION,
    KEY_STATE_CREATING,
    KEY_STATE_DISABLED,
    KEY_STATE_ENABLED,
    KEY_STATE_UNAVAILABLE,
  ];

  static final $core.List<KeyState?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 7);
  static KeyState? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyState._(super.value, super.name);
}

class KeyUsageType extends $pb.ProtobufEnum {
  static const KeyUsageType KEY_USAGE_TYPE_ENCRYPT_DECRYPT =
      KeyUsageType._(0, _omitEnumNames ? '' : 'KEY_USAGE_TYPE_ENCRYPT_DECRYPT');
  static const KeyUsageType KEY_USAGE_TYPE_KEY_AGREEMENT =
      KeyUsageType._(1, _omitEnumNames ? '' : 'KEY_USAGE_TYPE_KEY_AGREEMENT');
  static const KeyUsageType KEY_USAGE_TYPE_GENERATE_VERIFY_MAC = KeyUsageType._(
      2, _omitEnumNames ? '' : 'KEY_USAGE_TYPE_GENERATE_VERIFY_MAC');
  static const KeyUsageType KEY_USAGE_TYPE_SIGN_VERIFY =
      KeyUsageType._(3, _omitEnumNames ? '' : 'KEY_USAGE_TYPE_SIGN_VERIFY');

  static const $core.List<KeyUsageType> values = <KeyUsageType>[
    KEY_USAGE_TYPE_ENCRYPT_DECRYPT,
    KEY_USAGE_TYPE_KEY_AGREEMENT,
    KEY_USAGE_TYPE_GENERATE_VERIFY_MAC,
    KEY_USAGE_TYPE_SIGN_VERIFY,
  ];

  static final $core.List<KeyUsageType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static KeyUsageType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const KeyUsageType._(super.value, super.name);
}

class MacAlgorithmSpec extends $pb.ProtobufEnum {
  static const MacAlgorithmSpec MAC_ALGORITHM_SPEC_HMAC_SHA_512 =
      MacAlgorithmSpec._(
          0, _omitEnumNames ? '' : 'MAC_ALGORITHM_SPEC_HMAC_SHA_512');
  static const MacAlgorithmSpec MAC_ALGORITHM_SPEC_HMAC_SHA_224 =
      MacAlgorithmSpec._(
          1, _omitEnumNames ? '' : 'MAC_ALGORITHM_SPEC_HMAC_SHA_224');
  static const MacAlgorithmSpec MAC_ALGORITHM_SPEC_HMAC_SHA_384 =
      MacAlgorithmSpec._(
          2, _omitEnumNames ? '' : 'MAC_ALGORITHM_SPEC_HMAC_SHA_384');
  static const MacAlgorithmSpec MAC_ALGORITHM_SPEC_HMAC_SHA_256 =
      MacAlgorithmSpec._(
          3, _omitEnumNames ? '' : 'MAC_ALGORITHM_SPEC_HMAC_SHA_256');

  static const $core.List<MacAlgorithmSpec> values = <MacAlgorithmSpec>[
    MAC_ALGORITHM_SPEC_HMAC_SHA_512,
    MAC_ALGORITHM_SPEC_HMAC_SHA_224,
    MAC_ALGORITHM_SPEC_HMAC_SHA_384,
    MAC_ALGORITHM_SPEC_HMAC_SHA_256,
  ];

  static final $core.List<MacAlgorithmSpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static MacAlgorithmSpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MacAlgorithmSpec._(super.value, super.name);
}

class MessageType extends $pb.ProtobufEnum {
  static const MessageType MESSAGE_TYPE_RAW =
      MessageType._(0, _omitEnumNames ? '' : 'MESSAGE_TYPE_RAW');
  static const MessageType MESSAGE_TYPE_EXTERNAL_MU =
      MessageType._(1, _omitEnumNames ? '' : 'MESSAGE_TYPE_EXTERNAL_MU');
  static const MessageType MESSAGE_TYPE_DIGEST =
      MessageType._(2, _omitEnumNames ? '' : 'MESSAGE_TYPE_DIGEST');

  static const $core.List<MessageType> values = <MessageType>[
    MESSAGE_TYPE_RAW,
    MESSAGE_TYPE_EXTERNAL_MU,
    MESSAGE_TYPE_DIGEST,
  ];

  static final $core.List<MessageType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 2);
  static MessageType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MessageType._(super.value, super.name);
}

class MultiRegionKeyType extends $pb.ProtobufEnum {
  static const MultiRegionKeyType MULTI_REGION_KEY_TYPE_REPLICA =
      MultiRegionKeyType._(
          0, _omitEnumNames ? '' : 'MULTI_REGION_KEY_TYPE_REPLICA');
  static const MultiRegionKeyType MULTI_REGION_KEY_TYPE_PRIMARY =
      MultiRegionKeyType._(
          1, _omitEnumNames ? '' : 'MULTI_REGION_KEY_TYPE_PRIMARY');

  static const $core.List<MultiRegionKeyType> values = <MultiRegionKeyType>[
    MULTI_REGION_KEY_TYPE_REPLICA,
    MULTI_REGION_KEY_TYPE_PRIMARY,
  ];

  static final $core.List<MultiRegionKeyType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static MultiRegionKeyType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const MultiRegionKeyType._(super.value, super.name);
}

class OriginType extends $pb.ProtobufEnum {
  static const OriginType ORIGIN_TYPE_EXTERNAL =
      OriginType._(0, _omitEnumNames ? '' : 'ORIGIN_TYPE_EXTERNAL');
  static const OriginType ORIGIN_TYPE_AWS_KMS =
      OriginType._(1, _omitEnumNames ? '' : 'ORIGIN_TYPE_AWS_KMS');
  static const OriginType ORIGIN_TYPE_EXTERNAL_KEY_STORE =
      OriginType._(2, _omitEnumNames ? '' : 'ORIGIN_TYPE_EXTERNAL_KEY_STORE');
  static const OriginType ORIGIN_TYPE_AWS_CLOUDHSM =
      OriginType._(3, _omitEnumNames ? '' : 'ORIGIN_TYPE_AWS_CLOUDHSM');

  static const $core.List<OriginType> values = <OriginType>[
    ORIGIN_TYPE_EXTERNAL,
    ORIGIN_TYPE_AWS_KMS,
    ORIGIN_TYPE_EXTERNAL_KEY_STORE,
    ORIGIN_TYPE_AWS_CLOUDHSM,
  ];

  static final $core.List<OriginType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static OriginType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const OriginType._(super.value, super.name);
}

class RotationType extends $pb.ProtobufEnum {
  static const RotationType ROTATION_TYPE_AUTOMATIC =
      RotationType._(0, _omitEnumNames ? '' : 'ROTATION_TYPE_AUTOMATIC');
  static const RotationType ROTATION_TYPE_ON_DEMAND =
      RotationType._(1, _omitEnumNames ? '' : 'ROTATION_TYPE_ON_DEMAND');

  static const $core.List<RotationType> values = <RotationType>[
    ROTATION_TYPE_AUTOMATIC,
    ROTATION_TYPE_ON_DEMAND,
  ];

  static final $core.List<RotationType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static RotationType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const RotationType._(super.value, super.name);
}

class SigningAlgorithmSpec extends $pb.ProtobufEnum {
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_ECDSA_SHA_384 =
      SigningAlgorithmSpec._(
          0, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_ECDSA_SHA_384');
  static const SigningAlgorithmSpec
      SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_384 = SigningAlgorithmSpec._(
          1,
          _omitEnumNames
              ? ''
              : 'SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_384');
  static const SigningAlgorithmSpec
      SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_256 = SigningAlgorithmSpec._(
          2,
          _omitEnumNames
              ? ''
              : 'SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_256');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_ML_DSA_SHAKE_256 =
      SigningAlgorithmSpec._(
          3, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_ML_DSA_SHAKE_256');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_ED25519_SHA_512 =
      SigningAlgorithmSpec._(
          4, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_ED25519_SHA_512');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_ECDSA_SHA_512 =
      SigningAlgorithmSpec._(
          5, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_ECDSA_SHA_512');
  static const SigningAlgorithmSpec
      SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_512 = SigningAlgorithmSpec._(
          6,
          _omitEnumNames
              ? ''
              : 'SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_512');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_ED25519_PH_SHA_512 =
      SigningAlgorithmSpec._(
          7, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_ED25519_PH_SHA_512');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_256 =
      SigningAlgorithmSpec._(
          8, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_256');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_384 =
      SigningAlgorithmSpec._(
          9, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_384');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_512 =
      SigningAlgorithmSpec._(10,
          _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_512');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_ECDSA_SHA_256 =
      SigningAlgorithmSpec._(
          11, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_ECDSA_SHA_256');
  static const SigningAlgorithmSpec SIGNING_ALGORITHM_SPEC_SM2DSA =
      SigningAlgorithmSpec._(
          12, _omitEnumNames ? '' : 'SIGNING_ALGORITHM_SPEC_SM2DSA');

  static const $core.List<SigningAlgorithmSpec> values = <SigningAlgorithmSpec>[
    SIGNING_ALGORITHM_SPEC_ECDSA_SHA_384,
    SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_384,
    SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_256,
    SIGNING_ALGORITHM_SPEC_ML_DSA_SHAKE_256,
    SIGNING_ALGORITHM_SPEC_ED25519_SHA_512,
    SIGNING_ALGORITHM_SPEC_ECDSA_SHA_512,
    SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_512,
    SIGNING_ALGORITHM_SPEC_ED25519_PH_SHA_512,
    SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_256,
    SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_384,
    SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_512,
    SIGNING_ALGORITHM_SPEC_ECDSA_SHA_256,
    SIGNING_ALGORITHM_SPEC_SM2DSA,
  ];

  static final $core.List<SigningAlgorithmSpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 12);
  static SigningAlgorithmSpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const SigningAlgorithmSpec._(super.value, super.name);
}

class WrappingKeySpec extends $pb.ProtobufEnum {
  static const WrappingKeySpec WRAPPING_KEY_SPEC_SM2 =
      WrappingKeySpec._(0, _omitEnumNames ? '' : 'WRAPPING_KEY_SPEC_SM2');
  static const WrappingKeySpec WRAPPING_KEY_SPEC_RSA_2048 =
      WrappingKeySpec._(1, _omitEnumNames ? '' : 'WRAPPING_KEY_SPEC_RSA_2048');
  static const WrappingKeySpec WRAPPING_KEY_SPEC_RSA_3072 =
      WrappingKeySpec._(2, _omitEnumNames ? '' : 'WRAPPING_KEY_SPEC_RSA_3072');
  static const WrappingKeySpec WRAPPING_KEY_SPEC_RSA_4096 =
      WrappingKeySpec._(3, _omitEnumNames ? '' : 'WRAPPING_KEY_SPEC_RSA_4096');

  static const $core.List<WrappingKeySpec> values = <WrappingKeySpec>[
    WRAPPING_KEY_SPEC_SM2,
    WRAPPING_KEY_SPEC_RSA_2048,
    WRAPPING_KEY_SPEC_RSA_3072,
    WRAPPING_KEY_SPEC_RSA_4096,
  ];

  static final $core.List<WrappingKeySpec?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 3);
  static WrappingKeySpec? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const WrappingKeySpec._(super.value, super.name);
}

class XksProxyConnectivityType extends $pb.ProtobufEnum {
  static const XksProxyConnectivityType
      XKS_PROXY_CONNECTIVITY_TYPE_VPC_ENDPOINT_SERVICE =
      XksProxyConnectivityType._(
          0,
          _omitEnumNames
              ? ''
              : 'XKS_PROXY_CONNECTIVITY_TYPE_VPC_ENDPOINT_SERVICE');
  static const XksProxyConnectivityType
      XKS_PROXY_CONNECTIVITY_TYPE_PUBLIC_ENDPOINT = XksProxyConnectivityType._(
          1,
          _omitEnumNames ? '' : 'XKS_PROXY_CONNECTIVITY_TYPE_PUBLIC_ENDPOINT');

  static const $core.List<XksProxyConnectivityType> values =
      <XksProxyConnectivityType>[
    XKS_PROXY_CONNECTIVITY_TYPE_VPC_ENDPOINT_SERVICE,
    XKS_PROXY_CONNECTIVITY_TYPE_PUBLIC_ENDPOINT,
  ];

  static final $core.List<XksProxyConnectivityType?> _byValue =
      $pb.ProtobufEnum.$_initByValueList(values, 1);
  static XksProxyConnectivityType? valueOf($core.int value) =>
      value < 0 || value >= _byValue.length ? null : _byValue[value];

  const XksProxyConnectivityType._(super.value, super.name);
}

const $core.bool _omitEnumNames =
    $core.bool.fromEnvironment('protobuf.omit_enum_names');
