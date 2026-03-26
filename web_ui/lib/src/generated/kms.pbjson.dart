// This is a generated file - do not edit.
//
// Generated from kms.proto.

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

@$core.Deprecated('Use algorithmSpecDescriptor instead')
const AlgorithmSpec$json = {
  '1': 'AlgorithmSpec',
  '2': [
    {'1': 'ALGORITHM_SPEC_SM2PKE', '2': 0},
    {'1': 'ALGORITHM_SPEC_RSAES_PKCS1_V1_5', '2': 1},
    {'1': 'ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_256', '2': 2},
    {'1': 'ALGORITHM_SPEC_RSA_AES_KEY_WRAP_SHA_1', '2': 3},
    {'1': 'ALGORITHM_SPEC_RSAES_OAEP_SHA_1', '2': 4},
    {'1': 'ALGORITHM_SPEC_RSAES_OAEP_SHA_256', '2': 5},
  ],
};

/// Descriptor for `AlgorithmSpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List algorithmSpecDescriptor = $convert.base64Decode(
    'Cg1BbGdvcml0aG1TcGVjEhkKFUFMR09SSVRITV9TUEVDX1NNMlBLRRAAEiMKH0FMR09SSVRITV'
    '9TUEVDX1JTQUVTX1BLQ1MxX1YxXzUQARIrCidBTEdPUklUSE1fU1BFQ19SU0FfQUVTX0tFWV9X'
    'UkFQX1NIQV8yNTYQAhIpCiVBTEdPUklUSE1fU1BFQ19SU0FfQUVTX0tFWV9XUkFQX1NIQV8xEA'
    'MSIwofQUxHT1JJVEhNX1NQRUNfUlNBRVNfT0FFUF9TSEFfMRAEEiUKIUFMR09SSVRITV9TUEVD'
    'X1JTQUVTX09BRVBfU0hBXzI1NhAF');

@$core.Deprecated('Use connectionErrorCodeTypeDescriptor instead')
const ConnectionErrorCodeType$json = {
  '1': 'ConnectionErrorCodeType',
  '2': [
    {'1': 'CONNECTION_ERROR_CODE_TYPE_SUBNET_NOT_FOUND', '2': 0},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_TIMED_OUT', '2': 1},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_NETWORK_ERRORS', '2': 2},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_INVALID_CREDENTIALS', '2': 3},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_USER_LOGGED_IN', '2': 4},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_NOT_REACHABLE', '2': 5},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_CLUSTER_NOT_FOUND', '2': 6},
    {
      '1':
          'CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_INVALID_CONFIGURATION',
      '2': 7
    },
    {'1': 'CONNECTION_ERROR_CODE_TYPE_USER_NOT_FOUND', '2': 8},
    {
      '1': 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_TLS_CONFIGURATION',
      '2': 9
    },
    {'1': 'CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_CLOUDHSM_HSMS', '2': 10},
    {
      '1': 'CONNECTION_ERROR_CODE_TYPE_XKS_VPC_ENDPOINT_SERVICE_NOT_FOUND',
      '2': 11
    },
    {'1': 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_RESPONSE', '2': 12},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_USER_LOCKED_OUT', '2': 13},
    {
      '1': 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_INVALID_CONFIGURATION',
      '2': 14
    },
    {
      '1': 'CONNECTION_ERROR_CODE_TYPE_INSUFFICIENT_FREE_ADDRESSES_IN_SUBNET',
      '2': 15
    },
    {'1': 'CONNECTION_ERROR_CODE_TYPE_INTERNAL_ERROR', '2': 16},
    {'1': 'CONNECTION_ERROR_CODE_TYPE_XKS_PROXY_ACCESS_DENIED', '2': 17},
  ],
};

/// Descriptor for `ConnectionErrorCodeType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionErrorCodeTypeDescriptor = $convert.base64Decode(
    'ChdDb25uZWN0aW9uRXJyb3JDb2RlVHlwZRIvCitDT05ORUNUSU9OX0VSUk9SX0NPREVfVFlQRV'
    '9TVUJORVRfTk9UX0ZPVU5EEAASMgouQ09OTkVDVElPTl9FUlJPUl9DT0RFX1RZUEVfWEtTX1BS'
    'T1hZX1RJTUVEX09VVBABEi0KKUNPTk5FQ1RJT05fRVJST1JfQ09ERV9UWVBFX05FVFdPUktfRV'
    'JST1JTEAISMgouQ09OTkVDVElPTl9FUlJPUl9DT0RFX1RZUEVfSU5WQUxJRF9DUkVERU5USUFM'
    'UxADEi0KKUNPTk5FQ1RJT05fRVJST1JfQ09ERV9UWVBFX1VTRVJfTE9HR0VEX0lOEAQSNgoyQ0'
    '9OTkVDVElPTl9FUlJPUl9DT0RFX1RZUEVfWEtTX1BST1hZX05PVF9SRUFDSEFCTEUQBRIwCixD'
    'T05ORUNUSU9OX0VSUk9SX0NPREVfVFlQRV9DTFVTVEVSX05PVF9GT1VORBAGEk0KSUNPTk5FQ1'
    'RJT05fRVJST1JfQ09ERV9UWVBFX1hLU19WUENfRU5EUE9JTlRfU0VSVklDRV9JTlZBTElEX0NP'
    'TkZJR1VSQVRJT04QBxItCilDT05ORUNUSU9OX0VSUk9SX0NPREVfVFlQRV9VU0VSX05PVF9GT1'
    'VORBAIEkIKPkNPTk5FQ1RJT05fRVJST1JfQ09ERV9UWVBFX1hLU19QUk9YWV9JTlZBTElEX1RM'
    'U19DT05GSUdVUkFUSU9OEAkSOQo1Q09OTkVDVElPTl9FUlJPUl9DT0RFX1RZUEVfSU5TVUZGSU'
    'NJRU5UX0NMT1VESFNNX0hTTVMQChJBCj1DT05ORUNUSU9OX0VSUk9SX0NPREVfVFlQRV9YS1Nf'
    'VlBDX0VORFBPSU5UX1NFUlZJQ0VfTk9UX0ZPVU5EEAsSOQo1Q09OTkVDVElPTl9FUlJPUl9DT0'
    'RFX1RZUEVfWEtTX1BST1hZX0lOVkFMSURfUkVTUE9OU0UQDBIuCipDT05ORUNUSU9OX0VSUk9S'
    'X0NPREVfVFlQRV9VU0VSX0xPQ0tFRF9PVVQQDRI+CjpDT05ORUNUSU9OX0VSUk9SX0NPREVfVF'
    'lQRV9YS1NfUFJPWFlfSU5WQUxJRF9DT05GSUdVUkFUSU9OEA4SRApAQ09OTkVDVElPTl9FUlJP'
    'Ul9DT0RFX1RZUEVfSU5TVUZGSUNJRU5UX0ZSRUVfQUREUkVTU0VTX0lOX1NVQk5FVBAPEi0KKU'
    'NPTk5FQ1RJT05fRVJST1JfQ09ERV9UWVBFX0lOVEVSTkFMX0VSUk9SEBASNgoyQ09OTkVDVElP'
    'Tl9FUlJPUl9DT0RFX1RZUEVfWEtTX1BST1hZX0FDQ0VTU19ERU5JRUQQEQ==');

@$core.Deprecated('Use connectionStateTypeDescriptor instead')
const ConnectionStateType$json = {
  '1': 'ConnectionStateType',
  '2': [
    {'1': 'CONNECTION_STATE_TYPE_CONNECTING', '2': 0},
    {'1': 'CONNECTION_STATE_TYPE_DISCONNECTING', '2': 1},
    {'1': 'CONNECTION_STATE_TYPE_DISCONNECTED', '2': 2},
    {'1': 'CONNECTION_STATE_TYPE_CONNECTED', '2': 3},
    {'1': 'CONNECTION_STATE_TYPE_FAILED', '2': 4},
  ],
};

/// Descriptor for `ConnectionStateType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List connectionStateTypeDescriptor = $convert.base64Decode(
    'ChNDb25uZWN0aW9uU3RhdGVUeXBlEiQKIENPTk5FQ1RJT05fU1RBVEVfVFlQRV9DT05ORUNUSU'
    '5HEAASJwojQ09OTkVDVElPTl9TVEFURV9UWVBFX0RJU0NPTk5FQ1RJTkcQARImCiJDT05ORUNU'
    'SU9OX1NUQVRFX1RZUEVfRElTQ09OTkVDVEVEEAISIwofQ09OTkVDVElPTl9TVEFURV9UWVBFX0'
    'NPTk5FQ1RFRBADEiAKHENPTk5FQ1RJT05fU1RBVEVfVFlQRV9GQUlMRUQQBA==');

@$core.Deprecated('Use customKeyStoreTypeDescriptor instead')
const CustomKeyStoreType$json = {
  '1': 'CustomKeyStoreType',
  '2': [
    {'1': 'CUSTOM_KEY_STORE_TYPE_EXTERNAL_KEY_STORE', '2': 0},
    {'1': 'CUSTOM_KEY_STORE_TYPE_AWS_CLOUDHSM', '2': 1},
  ],
};

/// Descriptor for `CustomKeyStoreType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List customKeyStoreTypeDescriptor = $convert.base64Decode(
    'ChJDdXN0b21LZXlTdG9yZVR5cGUSLAooQ1VTVE9NX0tFWV9TVE9SRV9UWVBFX0VYVEVSTkFMX0'
    'tFWV9TVE9SRRAAEiYKIkNVU1RPTV9LRVlfU1RPUkVfVFlQRV9BV1NfQ0xPVURIU00QAQ==');

@$core.Deprecated('Use customerMasterKeySpecDescriptor instead')
const CustomerMasterKeySpec$json = {
  '1': 'CustomerMasterKeySpec',
  '2': [
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_HMAC_256', '2': 0},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P384', '2': 1},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_HMAC_384', '2': 2},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_SM2', '2': 3},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_ECC_SECG_P256K1', '2': 4},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P521', '2': 5},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_HMAC_512', '2': 6},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_HMAC_224', '2': 7},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_RSA_2048', '2': 8},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_RSA_3072', '2': 9},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_RSA_4096', '2': 10},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_SYMMETRIC_DEFAULT', '2': 11},
    {'1': 'CUSTOMER_MASTER_KEY_SPEC_ECC_NIST_P256', '2': 12},
  ],
};

/// Descriptor for `CustomerMasterKeySpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List customerMasterKeySpecDescriptor = $convert.base64Decode(
    'ChVDdXN0b21lck1hc3RlcktleVNwZWMSJQohQ1VTVE9NRVJfTUFTVEVSX0tFWV9TUEVDX0hNQU'
    'NfMjU2EAASKgomQ1VTVE9NRVJfTUFTVEVSX0tFWV9TUEVDX0VDQ19OSVNUX1AzODQQARIlCiFD'
    'VVNUT01FUl9NQVNURVJfS0VZX1NQRUNfSE1BQ18zODQQAhIgChxDVVNUT01FUl9NQVNURVJfS0'
    'VZX1NQRUNfU00yEAMSLAooQ1VTVE9NRVJfTUFTVEVSX0tFWV9TUEVDX0VDQ19TRUNHX1AyNTZL'
    'MRAEEioKJkNVU1RPTUVSX01BU1RFUl9LRVlfU1BFQ19FQ0NfTklTVF9QNTIxEAUSJQohQ1VTVE'
    '9NRVJfTUFTVEVSX0tFWV9TUEVDX0hNQUNfNTEyEAYSJQohQ1VTVE9NRVJfTUFTVEVSX0tFWV9T'
    'UEVDX0hNQUNfMjI0EAcSJQohQ1VTVE9NRVJfTUFTVEVSX0tFWV9TUEVDX1JTQV8yMDQ4EAgSJQ'
    'ohQ1VTVE9NRVJfTUFTVEVSX0tFWV9TUEVDX1JTQV8zMDcyEAkSJQohQ1VTVE9NRVJfTUFTVEVS'
    'X0tFWV9TUEVDX1JTQV80MDk2EAoSLgoqQ1VTVE9NRVJfTUFTVEVSX0tFWV9TUEVDX1NZTU1FVF'
    'JJQ19ERUZBVUxUEAsSKgomQ1VTVE9NRVJfTUFTVEVSX0tFWV9TUEVDX0VDQ19OSVNUX1AyNTYQ'
    'DA==');

@$core.Deprecated('Use dataKeyPairSpecDescriptor instead')
const DataKeyPairSpec$json = {
  '1': 'DataKeyPairSpec',
  '2': [
    {'1': 'DATA_KEY_PAIR_SPEC_ECC_NIST_EDWARDS25519', '2': 0},
    {'1': 'DATA_KEY_PAIR_SPEC_ECC_NIST_P384', '2': 1},
    {'1': 'DATA_KEY_PAIR_SPEC_SM2', '2': 2},
    {'1': 'DATA_KEY_PAIR_SPEC_ECC_SECG_P256K1', '2': 3},
    {'1': 'DATA_KEY_PAIR_SPEC_ECC_NIST_P521', '2': 4},
    {'1': 'DATA_KEY_PAIR_SPEC_RSA_2048', '2': 5},
    {'1': 'DATA_KEY_PAIR_SPEC_RSA_3072', '2': 6},
    {'1': 'DATA_KEY_PAIR_SPEC_RSA_4096', '2': 7},
    {'1': 'DATA_KEY_PAIR_SPEC_ECC_NIST_P256', '2': 8},
  ],
};

/// Descriptor for `DataKeyPairSpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dataKeyPairSpecDescriptor = $convert.base64Decode(
    'Cg9EYXRhS2V5UGFpclNwZWMSLAooREFUQV9LRVlfUEFJUl9TUEVDX0VDQ19OSVNUX0VEV0FSRF'
    'MyNTUxORAAEiQKIERBVEFfS0VZX1BBSVJfU1BFQ19FQ0NfTklTVF9QMzg0EAESGgoWREFUQV9L'
    'RVlfUEFJUl9TUEVDX1NNMhACEiYKIkRBVEFfS0VZX1BBSVJfU1BFQ19FQ0NfU0VDR19QMjU2Sz'
    'EQAxIkCiBEQVRBX0tFWV9QQUlSX1NQRUNfRUNDX05JU1RfUDUyMRAEEh8KG0RBVEFfS0VZX1BB'
    'SVJfU1BFQ19SU0FfMjA0OBAFEh8KG0RBVEFfS0VZX1BBSVJfU1BFQ19SU0FfMzA3MhAGEh8KG0'
    'RBVEFfS0VZX1BBSVJfU1BFQ19SU0FfNDA5NhAHEiQKIERBVEFfS0VZX1BBSVJfU1BFQ19FQ0Nf'
    'TklTVF9QMjU2EAg=');

@$core.Deprecated('Use dataKeySpecDescriptor instead')
const DataKeySpec$json = {
  '1': 'DataKeySpec',
  '2': [
    {'1': 'DATA_KEY_SPEC_AES_128', '2': 0},
    {'1': 'DATA_KEY_SPEC_AES_256', '2': 1},
  ],
};

/// Descriptor for `DataKeySpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dataKeySpecDescriptor = $convert.base64Decode(
    'CgtEYXRhS2V5U3BlYxIZChVEQVRBX0tFWV9TUEVDX0FFU18xMjgQABIZChVEQVRBX0tFWV9TUE'
    'VDX0FFU18yNTYQAQ==');

@$core.Deprecated('Use dryRunModifierTypeDescriptor instead')
const DryRunModifierType$json = {
  '1': 'DryRunModifierType',
  '2': [
    {'1': 'DRY_RUN_MODIFIER_TYPE_IGNORE_CIPHERTEXT', '2': 0},
  ],
};

/// Descriptor for `DryRunModifierType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List dryRunModifierTypeDescriptor = $convert.base64Decode(
    'ChJEcnlSdW5Nb2RpZmllclR5cGUSKwonRFJZX1JVTl9NT0RJRklFUl9UWVBFX0lHTk9SRV9DSV'
    'BIRVJURVhUEAA=');

@$core.Deprecated('Use encryptionAlgorithmSpecDescriptor instead')
const EncryptionAlgorithmSpec$json = {
  '1': 'EncryptionAlgorithmSpec',
  '2': [
    {'1': 'ENCRYPTION_ALGORITHM_SPEC_SM2PKE', '2': 0},
    {'1': 'ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_1', '2': 1},
    {'1': 'ENCRYPTION_ALGORITHM_SPEC_SYMMETRIC_DEFAULT', '2': 2},
    {'1': 'ENCRYPTION_ALGORITHM_SPEC_RSAES_OAEP_SHA_256', '2': 3},
  ],
};

/// Descriptor for `EncryptionAlgorithmSpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List encryptionAlgorithmSpecDescriptor = $convert.base64Decode(
    'ChdFbmNyeXB0aW9uQWxnb3JpdGhtU3BlYxIkCiBFTkNSWVBUSU9OX0FMR09SSVRITV9TUEVDX1'
    'NNMlBLRRAAEi4KKkVOQ1JZUFRJT05fQUxHT1JJVEhNX1NQRUNfUlNBRVNfT0FFUF9TSEFfMRAB'
    'Ei8KK0VOQ1JZUFRJT05fQUxHT1JJVEhNX1NQRUNfU1lNTUVUUklDX0RFRkFVTFQQAhIwCixFTk'
    'NSWVBUSU9OX0FMR09SSVRITV9TUEVDX1JTQUVTX09BRVBfU0hBXzI1NhAD');

@$core.Deprecated('Use expirationModelTypeDescriptor instead')
const ExpirationModelType$json = {
  '1': 'ExpirationModelType',
  '2': [
    {'1': 'EXPIRATION_MODEL_TYPE_KEY_MATERIAL_DOES_NOT_EXPIRE', '2': 0},
    {'1': 'EXPIRATION_MODEL_TYPE_KEY_MATERIAL_EXPIRES', '2': 1},
  ],
};

/// Descriptor for `ExpirationModelType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List expirationModelTypeDescriptor = $convert.base64Decode(
    'ChNFeHBpcmF0aW9uTW9kZWxUeXBlEjYKMkVYUElSQVRJT05fTU9ERUxfVFlQRV9LRVlfTUFURV'
    'JJQUxfRE9FU19OT1RfRVhQSVJFEAASLgoqRVhQSVJBVElPTl9NT0RFTF9UWVBFX0tFWV9NQVRF'
    'UklBTF9FWFBJUkVTEAE=');

@$core.Deprecated('Use grantOperationDescriptor instead')
const GrantOperation$json = {
  '1': 'GrantOperation',
  '2': [
    {'1': 'GRANT_OPERATION_GENERATEDATAKEYPAIR', '2': 0},
    {'1': 'GRANT_OPERATION_RETIREGRANT', '2': 1},
    {'1': 'GRANT_OPERATION_DESCRIBEKEY', '2': 2},
    {'1': 'GRANT_OPERATION_VERIFY', '2': 3},
    {'1': 'GRANT_OPERATION_ENCRYPT', '2': 4},
    {'1': 'GRANT_OPERATION_GENERATEDATAKEYWITHOUTPLAINTEXT', '2': 5},
    {'1': 'GRANT_OPERATION_GENERATEDATAKEYPAIRWITHOUTPLAINTEXT', '2': 6},
    {'1': 'GRANT_OPERATION_SIGN', '2': 7},
    {'1': 'GRANT_OPERATION_DERIVESHAREDSECRET', '2': 8},
    {'1': 'GRANT_OPERATION_VERIFYMAC', '2': 9},
    {'1': 'GRANT_OPERATION_GETPUBLICKEY', '2': 10},
    {'1': 'GRANT_OPERATION_CREATEGRANT', '2': 11},
    {'1': 'GRANT_OPERATION_REENCRYPTTO', '2': 12},
    {'1': 'GRANT_OPERATION_REENCRYPTFROM', '2': 13},
    {'1': 'GRANT_OPERATION_GENERATEDATAKEY', '2': 14},
    {'1': 'GRANT_OPERATION_DECRYPT', '2': 15},
    {'1': 'GRANT_OPERATION_GENERATEMAC', '2': 16},
  ],
};

/// Descriptor for `GrantOperation`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List grantOperationDescriptor = $convert.base64Decode(
    'Cg5HcmFudE9wZXJhdGlvbhInCiNHUkFOVF9PUEVSQVRJT05fR0VORVJBVEVEQVRBS0VZUEFJUh'
    'AAEh8KG0dSQU5UX09QRVJBVElPTl9SRVRJUkVHUkFOVBABEh8KG0dSQU5UX09QRVJBVElPTl9E'
    'RVNDUklCRUtFWRACEhoKFkdSQU5UX09QRVJBVElPTl9WRVJJRlkQAxIbChdHUkFOVF9PUEVSQV'
    'RJT05fRU5DUllQVBAEEjMKL0dSQU5UX09QRVJBVElPTl9HRU5FUkFURURBVEFLRVlXSVRIT1VU'
    'UExBSU5URVhUEAUSNwozR1JBTlRfT1BFUkFUSU9OX0dFTkVSQVRFREFUQUtFWVBBSVJXSVRIT1'
    'VUUExBSU5URVhUEAYSGAoUR1JBTlRfT1BFUkFUSU9OX1NJR04QBxImCiJHUkFOVF9PUEVSQVRJ'
    'T05fREVSSVZFU0hBUkVEU0VDUkVUEAgSHQoZR1JBTlRfT1BFUkFUSU9OX1ZFUklGWU1BQxAJEi'
    'AKHEdSQU5UX09QRVJBVElPTl9HRVRQVUJMSUNLRVkQChIfChtHUkFOVF9PUEVSQVRJT05fQ1JF'
    'QVRFR1JBTlQQCxIfChtHUkFOVF9PUEVSQVRJT05fUkVFTkNSWVBUVE8QDBIhCh1HUkFOVF9PUE'
    'VSQVRJT05fUkVFTkNSWVBURlJPTRANEiMKH0dSQU5UX09QRVJBVElPTl9HRU5FUkFURURBVEFL'
    'RVkQDhIbChdHUkFOVF9PUEVSQVRJT05fREVDUllQVBAPEh8KG0dSQU5UX09QRVJBVElPTl9HRU'
    '5FUkFURU1BQxAQ');

@$core.Deprecated('Use importStateDescriptor instead')
const ImportState$json = {
  '1': 'ImportState',
  '2': [
    {'1': 'IMPORT_STATE_IMPORTED', '2': 0},
    {'1': 'IMPORT_STATE_PENDING_IMPORT', '2': 1},
  ],
};

/// Descriptor for `ImportState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importStateDescriptor = $convert.base64Decode(
    'CgtJbXBvcnRTdGF0ZRIZChVJTVBPUlRfU1RBVEVfSU1QT1JURUQQABIfChtJTVBPUlRfU1RBVE'
    'VfUEVORElOR19JTVBPUlQQAQ==');

@$core.Deprecated('Use importTypeDescriptor instead')
const ImportType$json = {
  '1': 'ImportType',
  '2': [
    {'1': 'IMPORT_TYPE_EXISTING_KEY_MATERIAL', '2': 0},
    {'1': 'IMPORT_TYPE_NEW_KEY_MATERIAL', '2': 1},
  ],
};

/// Descriptor for `ImportType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List importTypeDescriptor = $convert.base64Decode(
    'CgpJbXBvcnRUeXBlEiUKIUlNUE9SVF9UWVBFX0VYSVNUSU5HX0tFWV9NQVRFUklBTBAAEiAKHE'
    'lNUE9SVF9UWVBFX05FV19LRVlfTUFURVJJQUwQAQ==');

@$core.Deprecated('Use includeKeyMaterialDescriptor instead')
const IncludeKeyMaterial$json = {
  '1': 'IncludeKeyMaterial',
  '2': [
    {'1': 'INCLUDE_KEY_MATERIAL_ALL_KEY_MATERIAL', '2': 0},
    {'1': 'INCLUDE_KEY_MATERIAL_ROTATIONS_ONLY', '2': 1},
  ],
};

/// Descriptor for `IncludeKeyMaterial`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List includeKeyMaterialDescriptor = $convert.base64Decode(
    'ChJJbmNsdWRlS2V5TWF0ZXJpYWwSKQolSU5DTFVERV9LRVlfTUFURVJJQUxfQUxMX0tFWV9NQV'
    'RFUklBTBAAEicKI0lOQ0xVREVfS0VZX01BVEVSSUFMX1JPVEFUSU9OU19PTkxZEAE=');

@$core.Deprecated('Use keyAgreementAlgorithmSpecDescriptor instead')
const KeyAgreementAlgorithmSpec$json = {
  '1': 'KeyAgreementAlgorithmSpec',
  '2': [
    {'1': 'KEY_AGREEMENT_ALGORITHM_SPEC_ECDH', '2': 0},
  ],
};

/// Descriptor for `KeyAgreementAlgorithmSpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyAgreementAlgorithmSpecDescriptor =
    $convert.base64Decode(
        'ChlLZXlBZ3JlZW1lbnRBbGdvcml0aG1TcGVjEiUKIUtFWV9BR1JFRU1FTlRfQUxHT1JJVEhNX1'
        'NQRUNfRUNESBAA');

@$core.Deprecated('Use keyEncryptionMechanismDescriptor instead')
const KeyEncryptionMechanism$json = {
  '1': 'KeyEncryptionMechanism',
  '2': [
    {'1': 'KEY_ENCRYPTION_MECHANISM_RSAES_OAEP_SHA_256', '2': 0},
  ],
};

/// Descriptor for `KeyEncryptionMechanism`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyEncryptionMechanismDescriptor =
    $convert.base64Decode(
        'ChZLZXlFbmNyeXB0aW9uTWVjaGFuaXNtEi8KK0tFWV9FTkNSWVBUSU9OX01FQ0hBTklTTV9SU0'
        'FFU19PQUVQX1NIQV8yNTYQAA==');

@$core.Deprecated('Use keyManagerTypeDescriptor instead')
const KeyManagerType$json = {
  '1': 'KeyManagerType',
  '2': [
    {'1': 'KEY_MANAGER_TYPE_CUSTOMER', '2': 0},
    {'1': 'KEY_MANAGER_TYPE_AWS', '2': 1},
  ],
};

/// Descriptor for `KeyManagerType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyManagerTypeDescriptor = $convert.base64Decode(
    'Cg5LZXlNYW5hZ2VyVHlwZRIdChlLRVlfTUFOQUdFUl9UWVBFX0NVU1RPTUVSEAASGAoUS0VZX0'
    '1BTkFHRVJfVFlQRV9BV1MQAQ==');

@$core.Deprecated('Use keyMaterialStateDescriptor instead')
const KeyMaterialState$json = {
  '1': 'KeyMaterialState',
  '2': [
    {'1': 'KEY_MATERIAL_STATE_PENDING_ROTATION', '2': 0},
    {
      '1': 'KEY_MATERIAL_STATE_PENDING_MULTI_REGION_IMPORT_AND_ROTATION',
      '2': 1
    },
    {'1': 'KEY_MATERIAL_STATE_NON_CURRENT', '2': 2},
    {'1': 'KEY_MATERIAL_STATE_CURRENT', '2': 3},
  ],
};

/// Descriptor for `KeyMaterialState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyMaterialStateDescriptor = $convert.base64Decode(
    'ChBLZXlNYXRlcmlhbFN0YXRlEicKI0tFWV9NQVRFUklBTF9TVEFURV9QRU5ESU5HX1JPVEFUSU'
    '9OEAASPwo7S0VZX01BVEVSSUFMX1NUQVRFX1BFTkRJTkdfTVVMVElfUkVHSU9OX0lNUE9SVF9B'
    'TkRfUk9UQVRJT04QARIiCh5LRVlfTUFURVJJQUxfU1RBVEVfTk9OX0NVUlJFTlQQAhIeChpLRV'
    'lfTUFURVJJQUxfU1RBVEVfQ1VSUkVOVBAD');

@$core.Deprecated('Use keySpecDescriptor instead')
const KeySpec$json = {
  '1': 'KeySpec',
  '2': [
    {'1': 'KEY_SPEC_HMAC_256', '2': 0},
    {'1': 'KEY_SPEC_ML_DSA_87', '2': 1},
    {'1': 'KEY_SPEC_ECC_NIST_EDWARDS25519', '2': 2},
    {'1': 'KEY_SPEC_ECC_NIST_P384', '2': 3},
    {'1': 'KEY_SPEC_HMAC_384', '2': 4},
    {'1': 'KEY_SPEC_SM2', '2': 5},
    {'1': 'KEY_SPEC_ECC_SECG_P256K1', '2': 6},
    {'1': 'KEY_SPEC_ECC_NIST_P521', '2': 7},
    {'1': 'KEY_SPEC_HMAC_512', '2': 8},
    {'1': 'KEY_SPEC_HMAC_224', '2': 9},
    {'1': 'KEY_SPEC_ML_DSA_44', '2': 10},
    {'1': 'KEY_SPEC_RSA_2048', '2': 11},
    {'1': 'KEY_SPEC_RSA_3072', '2': 12},
    {'1': 'KEY_SPEC_ML_DSA_65', '2': 13},
    {'1': 'KEY_SPEC_RSA_4096', '2': 14},
    {'1': 'KEY_SPEC_SYMMETRIC_DEFAULT', '2': 15},
    {'1': 'KEY_SPEC_ECC_NIST_P256', '2': 16},
  ],
};

/// Descriptor for `KeySpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keySpecDescriptor = $convert.base64Decode(
    'CgdLZXlTcGVjEhUKEUtFWV9TUEVDX0hNQUNfMjU2EAASFgoSS0VZX1NQRUNfTUxfRFNBXzg3EA'
    'ESIgoeS0VZX1NQRUNfRUNDX05JU1RfRURXQVJEUzI1NTE5EAISGgoWS0VZX1NQRUNfRUNDX05J'
    'U1RfUDM4NBADEhUKEUtFWV9TUEVDX0hNQUNfMzg0EAQSEAoMS0VZX1NQRUNfU00yEAUSHAoYS0'
    'VZX1NQRUNfRUNDX1NFQ0dfUDI1NksxEAYSGgoWS0VZX1NQRUNfRUNDX05JU1RfUDUyMRAHEhUK'
    'EUtFWV9TUEVDX0hNQUNfNTEyEAgSFQoRS0VZX1NQRUNfSE1BQ18yMjQQCRIWChJLRVlfU1BFQ1'
    '9NTF9EU0FfNDQQChIVChFLRVlfU1BFQ19SU0FfMjA0OBALEhUKEUtFWV9TUEVDX1JTQV8zMDcy'
    'EAwSFgoSS0VZX1NQRUNfTUxfRFNBXzY1EA0SFQoRS0VZX1NQRUNfUlNBXzQwOTYQDhIeChpLRV'
    'lfU1BFQ19TWU1NRVRSSUNfREVGQVVMVBAPEhoKFktFWV9TUEVDX0VDQ19OSVNUX1AyNTYQEA==');

@$core.Deprecated('Use keyStateDescriptor instead')
const KeyState$json = {
  '1': 'KeyState',
  '2': [
    {'1': 'KEY_STATE_PENDINGREPLICADELETION', '2': 0},
    {'1': 'KEY_STATE_UPDATING', '2': 1},
    {'1': 'KEY_STATE_PENDINGIMPORT', '2': 2},
    {'1': 'KEY_STATE_PENDINGDELETION', '2': 3},
    {'1': 'KEY_STATE_CREATING', '2': 4},
    {'1': 'KEY_STATE_DISABLED', '2': 5},
    {'1': 'KEY_STATE_ENABLED', '2': 6},
    {'1': 'KEY_STATE_UNAVAILABLE', '2': 7},
  ],
};

/// Descriptor for `KeyState`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyStateDescriptor = $convert.base64Decode(
    'CghLZXlTdGF0ZRIkCiBLRVlfU1RBVEVfUEVORElOR1JFUExJQ0FERUxFVElPThAAEhYKEktFWV'
    '9TVEFURV9VUERBVElORxABEhsKF0tFWV9TVEFURV9QRU5ESU5HSU1QT1JUEAISHQoZS0VZX1NU'
    'QVRFX1BFTkRJTkdERUxFVElPThADEhYKEktFWV9TVEFURV9DUkVBVElORxAEEhYKEktFWV9TVE'
    'FURV9ESVNBQkxFRBAFEhUKEUtFWV9TVEFURV9FTkFCTEVEEAYSGQoVS0VZX1NUQVRFX1VOQVZB'
    'SUxBQkxFEAc=');

@$core.Deprecated('Use keyUsageTypeDescriptor instead')
const KeyUsageType$json = {
  '1': 'KeyUsageType',
  '2': [
    {'1': 'KEY_USAGE_TYPE_ENCRYPT_DECRYPT', '2': 0},
    {'1': 'KEY_USAGE_TYPE_KEY_AGREEMENT', '2': 1},
    {'1': 'KEY_USAGE_TYPE_GENERATE_VERIFY_MAC', '2': 2},
    {'1': 'KEY_USAGE_TYPE_SIGN_VERIFY', '2': 3},
  ],
};

/// Descriptor for `KeyUsageType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List keyUsageTypeDescriptor = $convert.base64Decode(
    'CgxLZXlVc2FnZVR5cGUSIgoeS0VZX1VTQUdFX1RZUEVfRU5DUllQVF9ERUNSWVBUEAASIAocS0'
    'VZX1VTQUdFX1RZUEVfS0VZX0FHUkVFTUVOVBABEiYKIktFWV9VU0FHRV9UWVBFX0dFTkVSQVRF'
    'X1ZFUklGWV9NQUMQAhIeChpLRVlfVVNBR0VfVFlQRV9TSUdOX1ZFUklGWRAD');

@$core.Deprecated('Use macAlgorithmSpecDescriptor instead')
const MacAlgorithmSpec$json = {
  '1': 'MacAlgorithmSpec',
  '2': [
    {'1': 'MAC_ALGORITHM_SPEC_HMAC_SHA_512', '2': 0},
    {'1': 'MAC_ALGORITHM_SPEC_HMAC_SHA_224', '2': 1},
    {'1': 'MAC_ALGORITHM_SPEC_HMAC_SHA_384', '2': 2},
    {'1': 'MAC_ALGORITHM_SPEC_HMAC_SHA_256', '2': 3},
  ],
};

/// Descriptor for `MacAlgorithmSpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List macAlgorithmSpecDescriptor = $convert.base64Decode(
    'ChBNYWNBbGdvcml0aG1TcGVjEiMKH01BQ19BTEdPUklUSE1fU1BFQ19ITUFDX1NIQV81MTIQAB'
    'IjCh9NQUNfQUxHT1JJVEhNX1NQRUNfSE1BQ19TSEFfMjI0EAESIwofTUFDX0FMR09SSVRITV9T'
    'UEVDX0hNQUNfU0hBXzM4NBACEiMKH01BQ19BTEdPUklUSE1fU1BFQ19ITUFDX1NIQV8yNTYQAw'
    '==');

@$core.Deprecated('Use messageTypeDescriptor instead')
const MessageType$json = {
  '1': 'MessageType',
  '2': [
    {'1': 'MESSAGE_TYPE_RAW', '2': 0},
    {'1': 'MESSAGE_TYPE_EXTERNAL_MU', '2': 1},
    {'1': 'MESSAGE_TYPE_DIGEST', '2': 2},
  ],
};

/// Descriptor for `MessageType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List messageTypeDescriptor = $convert.base64Decode(
    'CgtNZXNzYWdlVHlwZRIUChBNRVNTQUdFX1RZUEVfUkFXEAASHAoYTUVTU0FHRV9UWVBFX0VYVE'
    'VSTkFMX01VEAESFwoTTUVTU0FHRV9UWVBFX0RJR0VTVBAC');

@$core.Deprecated('Use multiRegionKeyTypeDescriptor instead')
const MultiRegionKeyType$json = {
  '1': 'MultiRegionKeyType',
  '2': [
    {'1': 'MULTI_REGION_KEY_TYPE_REPLICA', '2': 0},
    {'1': 'MULTI_REGION_KEY_TYPE_PRIMARY', '2': 1},
  ],
};

/// Descriptor for `MultiRegionKeyType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List multiRegionKeyTypeDescriptor = $convert.base64Decode(
    'ChJNdWx0aVJlZ2lvbktleVR5cGUSIQodTVVMVElfUkVHSU9OX0tFWV9UWVBFX1JFUExJQ0EQAB'
    'IhCh1NVUxUSV9SRUdJT05fS0VZX1RZUEVfUFJJTUFSWRAB');

@$core.Deprecated('Use originTypeDescriptor instead')
const OriginType$json = {
  '1': 'OriginType',
  '2': [
    {'1': 'ORIGIN_TYPE_EXTERNAL', '2': 0},
    {'1': 'ORIGIN_TYPE_AWS_KMS', '2': 1},
    {'1': 'ORIGIN_TYPE_EXTERNAL_KEY_STORE', '2': 2},
    {'1': 'ORIGIN_TYPE_AWS_CLOUDHSM', '2': 3},
  ],
};

/// Descriptor for `OriginType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List originTypeDescriptor = $convert.base64Decode(
    'CgpPcmlnaW5UeXBlEhgKFE9SSUdJTl9UWVBFX0VYVEVSTkFMEAASFwoTT1JJR0lOX1RZUEVfQV'
    'dTX0tNUxABEiIKHk9SSUdJTl9UWVBFX0VYVEVSTkFMX0tFWV9TVE9SRRACEhwKGE9SSUdJTl9U'
    'WVBFX0FXU19DTE9VREhTTRAD');

@$core.Deprecated('Use rotationTypeDescriptor instead')
const RotationType$json = {
  '1': 'RotationType',
  '2': [
    {'1': 'ROTATION_TYPE_AUTOMATIC', '2': 0},
    {'1': 'ROTATION_TYPE_ON_DEMAND', '2': 1},
  ],
};

/// Descriptor for `RotationType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List rotationTypeDescriptor = $convert.base64Decode(
    'CgxSb3RhdGlvblR5cGUSGwoXUk9UQVRJT05fVFlQRV9BVVRPTUFUSUMQABIbChdST1RBVElPTl'
    '9UWVBFX09OX0RFTUFORBAB');

@$core.Deprecated('Use signingAlgorithmSpecDescriptor instead')
const SigningAlgorithmSpec$json = {
  '1': 'SigningAlgorithmSpec',
  '2': [
    {'1': 'SIGNING_ALGORITHM_SPEC_ECDSA_SHA_384', '2': 0},
    {'1': 'SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_384', '2': 1},
    {'1': 'SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_256', '2': 2},
    {'1': 'SIGNING_ALGORITHM_SPEC_ML_DSA_SHAKE_256', '2': 3},
    {'1': 'SIGNING_ALGORITHM_SPEC_ED25519_SHA_512', '2': 4},
    {'1': 'SIGNING_ALGORITHM_SPEC_ECDSA_SHA_512', '2': 5},
    {'1': 'SIGNING_ALGORITHM_SPEC_RSASSA_PKCS1_V1_5_SHA_512', '2': 6},
    {'1': 'SIGNING_ALGORITHM_SPEC_ED25519_PH_SHA_512', '2': 7},
    {'1': 'SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_256', '2': 8},
    {'1': 'SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_384', '2': 9},
    {'1': 'SIGNING_ALGORITHM_SPEC_RSASSA_PSS_SHA_512', '2': 10},
    {'1': 'SIGNING_ALGORITHM_SPEC_ECDSA_SHA_256', '2': 11},
    {'1': 'SIGNING_ALGORITHM_SPEC_SM2DSA', '2': 12},
  ],
};

/// Descriptor for `SigningAlgorithmSpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List signingAlgorithmSpecDescriptor = $convert.base64Decode(
    'ChRTaWduaW5nQWxnb3JpdGhtU3BlYxIoCiRTSUdOSU5HX0FMR09SSVRITV9TUEVDX0VDRFNBX1'
    'NIQV8zODQQABI0CjBTSUdOSU5HX0FMR09SSVRITV9TUEVDX1JTQVNTQV9QS0NTMV9WMV81X1NI'
    'QV8zODQQARI0CjBTSUdOSU5HX0FMR09SSVRITV9TUEVDX1JTQVNTQV9QS0NTMV9WMV81X1NIQV'
    '8yNTYQAhIrCidTSUdOSU5HX0FMR09SSVRITV9TUEVDX01MX0RTQV9TSEFLRV8yNTYQAxIqCiZT'
    'SUdOSU5HX0FMR09SSVRITV9TUEVDX0VEMjU1MTlfU0hBXzUxMhAEEigKJFNJR05JTkdfQUxHT1'
    'JJVEhNX1NQRUNfRUNEU0FfU0hBXzUxMhAFEjQKMFNJR05JTkdfQUxHT1JJVEhNX1NQRUNfUlNB'
    'U1NBX1BLQ1MxX1YxXzVfU0hBXzUxMhAGEi0KKVNJR05JTkdfQUxHT1JJVEhNX1NQRUNfRUQyNT'
    'UxOV9QSF9TSEFfNTEyEAcSLQopU0lHTklOR19BTEdPUklUSE1fU1BFQ19SU0FTU0FfUFNTX1NI'
    'QV8yNTYQCBItCilTSUdOSU5HX0FMR09SSVRITV9TUEVDX1JTQVNTQV9QU1NfU0hBXzM4NBAJEi'
    '0KKVNJR05JTkdfQUxHT1JJVEhNX1NQRUNfUlNBU1NBX1BTU19TSEFfNTEyEAoSKAokU0lHTklO'
    'R19BTEdPUklUSE1fU1BFQ19FQ0RTQV9TSEFfMjU2EAsSIQodU0lHTklOR19BTEdPUklUSE1fU1'
    'BFQ19TTTJEU0EQDA==');

@$core.Deprecated('Use wrappingKeySpecDescriptor instead')
const WrappingKeySpec$json = {
  '1': 'WrappingKeySpec',
  '2': [
    {'1': 'WRAPPING_KEY_SPEC_SM2', '2': 0},
    {'1': 'WRAPPING_KEY_SPEC_RSA_2048', '2': 1},
    {'1': 'WRAPPING_KEY_SPEC_RSA_3072', '2': 2},
    {'1': 'WRAPPING_KEY_SPEC_RSA_4096', '2': 3},
  ],
};

/// Descriptor for `WrappingKeySpec`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List wrappingKeySpecDescriptor = $convert.base64Decode(
    'Cg9XcmFwcGluZ0tleVNwZWMSGQoVV1JBUFBJTkdfS0VZX1NQRUNfU00yEAASHgoaV1JBUFBJTk'
    'dfS0VZX1NQRUNfUlNBXzIwNDgQARIeChpXUkFQUElOR19LRVlfU1BFQ19SU0FfMzA3MhACEh4K'
    'GldSQVBQSU5HX0tFWV9TUEVDX1JTQV80MDk2EAM=');

@$core.Deprecated('Use xksProxyConnectivityTypeDescriptor instead')
const XksProxyConnectivityType$json = {
  '1': 'XksProxyConnectivityType',
  '2': [
    {'1': 'XKS_PROXY_CONNECTIVITY_TYPE_VPC_ENDPOINT_SERVICE', '2': 0},
    {'1': 'XKS_PROXY_CONNECTIVITY_TYPE_PUBLIC_ENDPOINT', '2': 1},
  ],
};

/// Descriptor for `XksProxyConnectivityType`. Decode as a `google.protobuf.EnumDescriptorProto`.
final $typed_data.Uint8List xksProxyConnectivityTypeDescriptor = $convert.base64Decode(
    'ChhYa3NQcm94eUNvbm5lY3Rpdml0eVR5cGUSNAowWEtTX1BST1hZX0NPTk5FQ1RJVklUWV9UWV'
    'BFX1ZQQ19FTkRQT0lOVF9TRVJWSUNFEAASLworWEtTX1BST1hZX0NPTk5FQ1RJVklUWV9UWVBF'
    'X1BVQkxJQ19FTkRQT0lOVBAB');

@$core.Deprecated('Use aliasListEntryDescriptor instead')
const AliasListEntry$json = {
  '1': 'AliasListEntry',
  '2': [
    {'1': 'aliasarn', '3': 461101595, '4': 1, '5': 9, '10': 'aliasarn'},
    {'1': 'aliasname', '3': 313250709, '4': 1, '5': 9, '10': 'aliasname'},
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'lastupdateddate',
      '3': 166338449,
      '4': 1,
      '5': 9,
      '10': 'lastupdateddate'
    },
    {'1': 'targetkeyid', '3': 406196123, '4': 1, '5': 9, '10': 'targetkeyid'},
  ],
};

/// Descriptor for `AliasListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List aliasListEntryDescriptor = $convert.base64Decode(
    'Cg5BbGlhc0xpc3RFbnRyeRIeCghhbGlhc2FybhibtO/bASABKAlSCGFsaWFzYXJuEiAKCWFsaW'
    'FzbmFtZRiVp6+VASABKAlSCWFsaWFzbmFtZRImCgxjcmVhdGlvbmRhdGUY4di3iQEgASgJUgxj'
    'cmVhdGlvbmRhdGUSKwoPbGFzdHVwZGF0ZWRkYXRlGJG/qE8gASgJUg9sYXN0dXBkYXRlZGRhdG'
    'USJAoLdGFyZ2V0a2V5aWQYm5/YwQEgASgJUgt0YXJnZXRrZXlpZA==');

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

@$core.Deprecated('Use cancelKeyDeletionRequestDescriptor instead')
const CancelKeyDeletionRequest$json = {
  '1': 'CancelKeyDeletionRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `CancelKeyDeletionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelKeyDeletionRequestDescriptor =
    $convert.base64Decode(
        'ChhDYW5jZWxLZXlEZWxldGlvblJlcXVlc3QSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZA==');

@$core.Deprecated('Use cancelKeyDeletionResponseDescriptor instead')
const CancelKeyDeletionResponse$json = {
  '1': 'CancelKeyDeletionResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `CancelKeyDeletionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cancelKeyDeletionResponseDescriptor =
    $convert.base64Decode(
        'ChlDYW5jZWxLZXlEZWxldGlvblJlc3BvbnNlEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQ=');

@$core.Deprecated('Use cloudHsmClusterInUseExceptionDescriptor instead')
const CloudHsmClusterInUseException$json = {
  '1': 'CloudHsmClusterInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudHsmClusterInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudHsmClusterInUseExceptionDescriptor =
    $convert.base64Decode(
        'Ch1DbG91ZEhzbUNsdXN0ZXJJblVzZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated(
    'Use cloudHsmClusterInvalidConfigurationExceptionDescriptor instead')
const CloudHsmClusterInvalidConfigurationException$json = {
  '1': 'CloudHsmClusterInvalidConfigurationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudHsmClusterInvalidConfigurationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    cloudHsmClusterInvalidConfigurationExceptionDescriptor =
    $convert.base64Decode(
        'CixDbG91ZEhzbUNsdXN0ZXJJbnZhbGlkQ29uZmlndXJhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYW'
        'dlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use cloudHsmClusterNotActiveExceptionDescriptor instead')
const CloudHsmClusterNotActiveException$json = {
  '1': 'CloudHsmClusterNotActiveException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudHsmClusterNotActiveException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudHsmClusterNotActiveExceptionDescriptor =
    $convert.base64Decode(
        'CiFDbG91ZEhzbUNsdXN0ZXJOb3RBY3RpdmVFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use cloudHsmClusterNotFoundExceptionDescriptor instead')
const CloudHsmClusterNotFoundException$json = {
  '1': 'CloudHsmClusterNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudHsmClusterNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudHsmClusterNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CiBDbG91ZEhzbUNsdXN0ZXJOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use cloudHsmClusterNotRelatedExceptionDescriptor instead')
const CloudHsmClusterNotRelatedException$json = {
  '1': 'CloudHsmClusterNotRelatedException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CloudHsmClusterNotRelatedException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List cloudHsmClusterNotRelatedExceptionDescriptor =
    $convert.base64Decode(
        'CiJDbG91ZEhzbUNsdXN0ZXJOb3RSZWxhdGVkRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKA'
        'lSB21lc3NhZ2U=');

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

@$core.Deprecated('Use connectCustomKeyStoreRequestDescriptor instead')
const ConnectCustomKeyStoreRequest$json = {
  '1': 'ConnectCustomKeyStoreRequest',
  '2': [
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
  ],
};

/// Descriptor for `ConnectCustomKeyStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectCustomKeyStoreRequestDescriptor =
    $convert.base64Decode(
        'ChxDb25uZWN0Q3VzdG9tS2V5U3RvcmVSZXF1ZXN0Ei0KEGN1c3RvbWtleXN0b3JlaWQYxKyQKi'
        'ABKAlSEGN1c3RvbWtleXN0b3JlaWQ=');

@$core.Deprecated('Use connectCustomKeyStoreResponseDescriptor instead')
const ConnectCustomKeyStoreResponse$json = {
  '1': 'ConnectCustomKeyStoreResponse',
};

/// Descriptor for `ConnectCustomKeyStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List connectCustomKeyStoreResponseDescriptor =
    $convert.base64Decode('Ch1Db25uZWN0Q3VzdG9tS2V5U3RvcmVSZXNwb25zZQ==');

@$core.Deprecated('Use createAliasRequestDescriptor instead')
const CreateAliasRequest$json = {
  '1': 'CreateAliasRequest',
  '2': [
    {'1': 'aliasname', '3': 313250709, '4': 1, '5': 9, '10': 'aliasname'},
    {'1': 'targetkeyid', '3': 406196123, '4': 1, '5': 9, '10': 'targetkeyid'},
  ],
};

/// Descriptor for `CreateAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createAliasRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVBbGlhc1JlcXVlc3QSIAoJYWxpYXNuYW1lGJWnr5UBIAEoCVIJYWxpYXNuYW1lEi'
    'QKC3RhcmdldGtleWlkGJuf2MEBIAEoCVILdGFyZ2V0a2V5aWQ=');

@$core.Deprecated('Use createCustomKeyStoreRequestDescriptor instead')
const CreateCustomKeyStoreRequest$json = {
  '1': 'CreateCustomKeyStoreRequest',
  '2': [
    {
      '1': 'cloudhsmclusterid',
      '3': 56498754,
      '4': 1,
      '5': 9,
      '10': 'cloudhsmclusterid'
    },
    {
      '1': 'customkeystorename',
      '3': 170278046,
      '4': 1,
      '5': 9,
      '10': 'customkeystorename'
    },
    {
      '1': 'customkeystoretype',
      '3': 415647103,
      '4': 1,
      '5': 14,
      '6': '.kms.CustomKeyStoreType',
      '10': 'customkeystoretype'
    },
    {
      '1': 'keystorepassword',
      '3': 403136353,
      '4': 1,
      '5': 9,
      '10': 'keystorepassword'
    },
    {
      '1': 'trustanchorcertificate',
      '3': 48354588,
      '4': 1,
      '5': 9,
      '10': 'trustanchorcertificate'
    },
    {
      '1': 'xksproxyauthenticationcredential',
      '3': 350418199,
      '4': 1,
      '5': 11,
      '6': '.kms.XksProxyAuthenticationCredentialType',
      '10': 'xksproxyauthenticationcredential'
    },
    {
      '1': 'xksproxyconnectivity',
      '3': 298569161,
      '4': 1,
      '5': 14,
      '6': '.kms.XksProxyConnectivityType',
      '10': 'xksproxyconnectivity'
    },
    {
      '1': 'xksproxyuriendpoint',
      '3': 273255559,
      '4': 1,
      '5': 9,
      '10': 'xksproxyuriendpoint'
    },
    {
      '1': 'xksproxyuripath',
      '3': 436753509,
      '4': 1,
      '5': 9,
      '10': 'xksproxyuripath'
    },
    {
      '1': 'xksproxyvpcendpointservicename',
      '3': 372786130,
      '4': 1,
      '5': 9,
      '10': 'xksproxyvpcendpointservicename'
    },
    {
      '1': 'xksproxyvpcendpointserviceowner',
      '3': 55249590,
      '4': 1,
      '5': 9,
      '10': 'xksproxyvpcendpointserviceowner'
    },
  ],
};

/// Descriptor for `CreateCustomKeyStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCustomKeyStoreRequestDescriptor = $convert.base64Decode(
    'ChtDcmVhdGVDdXN0b21LZXlTdG9yZVJlcXVlc3QSLwoRY2xvdWRoc21jbHVzdGVyaWQYwrT4Gi'
    'ABKAlSEWNsb3VkaHNtY2x1c3RlcmlkEjEKEmN1c3RvbWtleXN0b3JlbmFtZRie+ZhRIAEoCVIS'
    'Y3VzdG9ta2V5c3RvcmVuYW1lEksKEmN1c3RvbWtleXN0b3JldHlwZRj/ipnGASABKA4yFy5rbX'
    'MuQ3VzdG9tS2V5U3RvcmVUeXBlUhJjdXN0b21rZXlzdG9yZXR5cGUSLgoQa2V5c3RvcmVwYXNz'
    'd29yZBjhvp3AASABKAlSEGtleXN0b3JlcGFzc3dvcmQSOQoWdHJ1c3RhbmNob3JjZXJ0aWZpY2'
    'F0ZRicqocXIAEoCVIWdHJ1c3RhbmNob3JjZXJ0aWZpY2F0ZRJ5CiB4a3Nwcm94eWF1dGhlbnRp'
    'Y2F0aW9uY3JlZGVudGlhbBiX6ounASABKAsyKS5rbXMuWGtzUHJveHlBdXRoZW50aWNhdGlvbk'
    'NyZWRlbnRpYWxUeXBlUiB4a3Nwcm94eWF1dGhlbnRpY2F0aW9uY3JlZGVudGlhbBJVChR4a3Nw'
    'cm94eWNvbm5lY3Rpdml0eRjJm6+OASABKA4yHS5rbXMuWGtzUHJveHlDb25uZWN0aXZpdHlUeX'
    'BlUhR4a3Nwcm94eWNvbm5lY3Rpdml0eRI0ChN4a3Nwcm94eXVyaWVuZHBvaW50GIeZpoIBIAEo'
    'CVITeGtzcHJveHl1cmllbmRwb2ludBIsCg94a3Nwcm94eXVyaXBhdGgY5aih0AEgASgJUg94a3'
    'Nwcm94eXVyaXBhdGgSSgoeeGtzcHJveHl2cGNlbmRwb2ludHNlcnZpY2VuYW1lGNKH4bEBIAEo'
    'CVIeeGtzcHJveHl2cGNlbmRwb2ludHNlcnZpY2VuYW1lEksKH3hrc3Byb3h5dnBjZW5kcG9pbn'
    'RzZXJ2aWNlb3duZXIYtpWsGiABKAlSH3hrc3Byb3h5dnBjZW5kcG9pbnRzZXJ2aWNlb3duZXI=');

@$core.Deprecated('Use createCustomKeyStoreResponseDescriptor instead')
const CreateCustomKeyStoreResponse$json = {
  '1': 'CreateCustomKeyStoreResponse',
  '2': [
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
  ],
};

/// Descriptor for `CreateCustomKeyStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createCustomKeyStoreResponseDescriptor =
    $convert.base64Decode(
        'ChxDcmVhdGVDdXN0b21LZXlTdG9yZVJlc3BvbnNlEi0KEGN1c3RvbWtleXN0b3JlaWQYxKyQKi'
        'ABKAlSEGN1c3RvbWtleXN0b3JlaWQ=');

@$core.Deprecated('Use createGrantRequestDescriptor instead')
const CreateGrantRequest$json = {
  '1': 'CreateGrantRequest',
  '2': [
    {
      '1': 'constraints',
      '3': 302297388,
      '4': 1,
      '5': 11,
      '6': '.kms.GrantConstraints',
      '10': 'constraints'
    },
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {
      '1': 'granteeprincipal',
      '3': 234727364,
      '4': 1,
      '5': 9,
      '10': 'granteeprincipal'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'operations',
      '3': 126776656,
      '4': 3,
      '5': 14,
      '6': '.kms.GrantOperation',
      '10': 'operations'
    },
    {
      '1': 'retiringprincipal',
      '3': 49541086,
      '4': 1,
      '5': 9,
      '10': 'retiringprincipal'
    },
  ],
};

/// Descriptor for `CreateGrantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGrantRequestDescriptor = $convert.base64Decode(
    'ChJDcmVhdGVHcmFudFJlcXVlc3QSOwoLY29uc3RyYWludHMYrOKSkAEgASgLMhUua21zLkdyYW'
    '50Q29uc3RyYWludHNSC2NvbnN0cmFpbnRzEhkKBmRyeXJ1bhi43/srIAEoCFIGZHJ5cnVuEiQK'
    'C2dyYW50dG9rZW5zGIyNgKIBIAMoCVILZ3JhbnR0b2tlbnMSLQoQZ3JhbnRlZXByaW5jaXBhbB'
    'jEz/ZvIAEoCVIQZ3JhbnRlZXByaW5jaXBhbBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEhUK'
    'BG5hbWUYh+aBfyABKAlSBG5hbWUSNgoKb3BlcmF0aW9ucxjQ6rk8IAMoDjITLmttcy5HcmFudE'
    '9wZXJhdGlvblIKb3BlcmF0aW9ucxIvChFyZXRpcmluZ3ByaW5jaXBhbBje388XIAEoCVIRcmV0'
    'aXJpbmdwcmluY2lwYWw=');

@$core.Deprecated('Use createGrantResponseDescriptor instead')
const CreateGrantResponse$json = {
  '1': 'CreateGrantResponse',
  '2': [
    {'1': 'grantid', '3': 66852281, '4': 1, '5': 9, '10': 'grantid'},
    {'1': 'granttoken', '3': 137683547, '4': 1, '5': 9, '10': 'granttoken'},
  ],
};

/// Descriptor for `CreateGrantResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createGrantResponseDescriptor = $convert.base64Decode(
    'ChNDcmVhdGVHcmFudFJlc3BvbnNlEhsKB2dyYW50aWQYuavwHyABKAlSB2dyYW50aWQSIQoKZ3'
    'JhbnR0b2tlbhjbxNNBIAEoCVIKZ3JhbnR0b2tlbg==');

@$core.Deprecated('Use createKeyRequestDescriptor instead')
const CreateKeyRequest$json = {
  '1': 'CreateKeyRequest',
  '2': [
    {
      '1': 'bypasspolicylockoutsafetycheck',
      '3': 177450851,
      '4': 1,
      '5': 8,
      '10': 'bypasspolicylockoutsafetycheck'
    },
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
    {
      '1': 'customermasterkeyspec',
      '3': 472470930,
      '4': 1,
      '5': 14,
      '6': '.kms.CustomerMasterKeySpec',
      '10': 'customermasterkeyspec'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {
      '1': 'keyspec',
      '3': 138220928,
      '4': 1,
      '5': 14,
      '6': '.kms.KeySpec',
      '10': 'keyspec'
    },
    {
      '1': 'keyusage',
      '3': 357216772,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyUsageType',
      '10': 'keyusage'
    },
    {'1': 'multiregion', '3': 405769103, '4': 1, '5': 8, '10': 'multiregion'},
    {
      '1': 'origin',
      '3': 529732720,
      '4': 1,
      '5': 14,
      '6': '.kms.OriginType',
      '10': 'origin'
    },
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kms.Tag',
      '10': 'tags'
    },
    {'1': 'xkskeyid', '3': 12647506, '4': 1, '5': 9, '10': 'xkskeyid'},
  ],
};

/// Descriptor for `CreateKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeyRequestDescriptor = $convert.base64Decode(
    'ChBDcmVhdGVLZXlSZXF1ZXN0EkkKHmJ5cGFzc3BvbGljeWxvY2tvdXRzYWZldHljaGVjaxjj3s'
    '5UIAEoCFIeYnlwYXNzcG9saWN5bG9ja291dHNhZmV0eWNoZWNrEi0KEGN1c3RvbWtleXN0b3Jl'
    'aWQYxKyQKiABKAlSEGN1c3RvbWtleXN0b3JlaWQSVAoVY3VzdG9tZXJtYXN0ZXJrZXlzcGVjGJ'
    'KrpeEBIAEoDjIaLmttcy5DdXN0b21lck1hc3RlcktleVNwZWNSFWN1c3RvbWVybWFzdGVya2V5'
    'c3BlYxIjCgtkZXNjcmlwdGlvbhiK9Pk2IAEoCVILZGVzY3JpcHRpb24SKQoHa2V5c3BlYxiAq/'
    'RBIAEoDjIMLmttcy5LZXlTcGVjUgdrZXlzcGVjEjEKCGtleXVzYWdlGITkqqoBIAEoDjIRLmtt'
    'cy5LZXlVc2FnZVR5cGVSCGtleXVzYWdlEiQKC211bHRpcmVnaW9uGI+XvsEBIAEoCFILbXVsdG'
    'lyZWdpb24SKwoGb3JpZ2luGPCozPwBIAEoDjIPLmttcy5PcmlnaW5UeXBlUgZvcmlnaW4SGgoG'
    'cG9saWN5GKDv8OABIAEoCVIGcG9saWN5EiAKBHRhZ3MYwcH2tQEgAygLMggua21zLlRhZ1IEdG'
    'FncxIdCgh4a3NrZXlpZBjS+IMGIAEoCVIIeGtza2V5aWQ=');

@$core.Deprecated('Use createKeyResponseDescriptor instead')
const CreateKeyResponse$json = {
  '1': 'CreateKeyResponse',
  '2': [
    {
      '1': 'keymetadata',
      '3': 149520202,
      '4': 1,
      '5': 11,
      '6': '.kms.KeyMetadata',
      '10': 'keymetadata'
    },
  ],
};

/// Descriptor for `CreateKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List createKeyResponseDescriptor = $convert.base64Decode(
    'ChFDcmVhdGVLZXlSZXNwb25zZRI1CgtrZXltZXRhZGF0YRjK/qVHIAEoCzIQLmttcy5LZXlNZX'
    'RhZGF0YVILa2V5bWV0YWRhdGE=');

@$core.Deprecated('Use customKeyStoreHasCMKsExceptionDescriptor instead')
const CustomKeyStoreHasCMKsException$json = {
  '1': 'CustomKeyStoreHasCMKsException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CustomKeyStoreHasCMKsException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customKeyStoreHasCMKsExceptionDescriptor =
    $convert.base64Decode(
        'Ch5DdXN0b21LZXlTdG9yZUhhc0NNS3NFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbW'
        'Vzc2FnZQ==');

@$core.Deprecated('Use customKeyStoreInvalidStateExceptionDescriptor instead')
const CustomKeyStoreInvalidStateException$json = {
  '1': 'CustomKeyStoreInvalidStateException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CustomKeyStoreInvalidStateException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customKeyStoreInvalidStateExceptionDescriptor =
    $convert.base64Decode(
        'CiNDdXN0b21LZXlTdG9yZUludmFsaWRTdGF0ZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use customKeyStoreNameInUseExceptionDescriptor instead')
const CustomKeyStoreNameInUseException$json = {
  '1': 'CustomKeyStoreNameInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CustomKeyStoreNameInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customKeyStoreNameInUseExceptionDescriptor =
    $convert.base64Decode(
        'CiBDdXN0b21LZXlTdG9yZU5hbWVJblVzZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use customKeyStoreNotFoundExceptionDescriptor instead')
const CustomKeyStoreNotFoundException$json = {
  '1': 'CustomKeyStoreNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `CustomKeyStoreNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customKeyStoreNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'Ch9DdXN0b21LZXlTdG9yZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated('Use customKeyStoresListEntryDescriptor instead')
const CustomKeyStoresListEntry$json = {
  '1': 'CustomKeyStoresListEntry',
  '2': [
    {
      '1': 'cloudhsmclusterid',
      '3': 56498754,
      '4': 1,
      '5': 9,
      '10': 'cloudhsmclusterid'
    },
    {
      '1': 'connectionerrorcode',
      '3': 324951101,
      '4': 1,
      '5': 14,
      '6': '.kms.ConnectionErrorCodeType',
      '10': 'connectionerrorcode'
    },
    {
      '1': 'connectionstate',
      '3': 404323675,
      '4': 1,
      '5': 14,
      '6': '.kms.ConnectionStateType',
      '10': 'connectionstate'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
    {
      '1': 'customkeystorename',
      '3': 170278046,
      '4': 1,
      '5': 9,
      '10': 'customkeystorename'
    },
    {
      '1': 'customkeystoretype',
      '3': 415647103,
      '4': 1,
      '5': 14,
      '6': '.kms.CustomKeyStoreType',
      '10': 'customkeystoretype'
    },
    {
      '1': 'trustanchorcertificate',
      '3': 48354588,
      '4': 1,
      '5': 9,
      '10': 'trustanchorcertificate'
    },
    {
      '1': 'xksproxyconfiguration',
      '3': 349047828,
      '4': 1,
      '5': 11,
      '6': '.kms.XksProxyConfigurationType',
      '10': 'xksproxyconfiguration'
    },
  ],
};

/// Descriptor for `CustomKeyStoresListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List customKeyStoresListEntryDescriptor = $convert.base64Decode(
    'ChhDdXN0b21LZXlTdG9yZXNMaXN0RW50cnkSLwoRY2xvdWRoc21jbHVzdGVyaWQYwrT4GiABKA'
    'lSEWNsb3VkaHNtY2x1c3RlcmlkElIKE2Nvbm5lY3Rpb25lcnJvcmNvZGUYvbj5mgEgASgOMhwu'
    'a21zLkNvbm5lY3Rpb25FcnJvckNvZGVUeXBlUhNjb25uZWN0aW9uZXJyb3Jjb2RlEkYKD2Nvbm'
    '5lY3Rpb25zdGF0ZRjb+uXAASABKA4yGC5rbXMuQ29ubmVjdGlvblN0YXRlVHlwZVIPY29ubmVj'
    'dGlvbnN0YXRlEiYKDGNyZWF0aW9uZGF0ZRjh2LeJASABKAlSDGNyZWF0aW9uZGF0ZRItChBjdX'
    'N0b21rZXlzdG9yZWlkGMSskCogASgJUhBjdXN0b21rZXlzdG9yZWlkEjEKEmN1c3RvbWtleXN0'
    'b3JlbmFtZRie+ZhRIAEoCVISY3VzdG9ta2V5c3RvcmVuYW1lEksKEmN1c3RvbWtleXN0b3JldH'
    'lwZRj/ipnGASABKA4yFy5rbXMuQ3VzdG9tS2V5U3RvcmVUeXBlUhJjdXN0b21rZXlzdG9yZXR5'
    'cGUSOQoWdHJ1c3RhbmNob3JjZXJ0aWZpY2F0ZRicqocXIAEoCVIWdHJ1c3RhbmNob3JjZXJ0aW'
    'ZpY2F0ZRJYChV4a3Nwcm94eWNvbmZpZ3VyYXRpb24YlJi4pgEgASgLMh4ua21zLlhrc1Byb3h5'
    'Q29uZmlndXJhdGlvblR5cGVSFXhrc3Byb3h5Y29uZmlndXJhdGlvbg==');

@$core.Deprecated('Use decryptRequestDescriptor instead')
const DecryptRequest$json = {
  '1': 'DecryptRequest',
  '2': [
    {
      '1': 'ciphertextblob',
      '3': 338198183,
      '4': 1,
      '5': 12,
      '10': 'ciphertextblob'
    },
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'dryrunmodifiers',
      '3': 24346424,
      '4': 3,
      '5': 14,
      '6': '.kms.DryRunModifierType',
      '10': 'dryrunmodifiers'
    },
    {
      '1': 'encryptionalgorithm',
      '3': 203224586,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'encryptionalgorithm'
    },
    {
      '1': 'encryptioncontext',
      '3': 286249536,
      '4': 3,
      '5': 11,
      '6': '.kms.DecryptRequest.EncryptioncontextEntry',
      '10': 'encryptioncontext'
    },
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'recipient',
      '3': 445981721,
      '4': 1,
      '5': 11,
      '6': '.kms.RecipientInfo',
      '10': 'recipient'
    },
  ],
  '3': [DecryptRequest_EncryptioncontextEntry$json],
};

@$core.Deprecated('Use decryptRequestDescriptor instead')
const DecryptRequest_EncryptioncontextEntry$json = {
  '1': 'EncryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `DecryptRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List decryptRequestDescriptor = $convert.base64Decode(
    'Cg5EZWNyeXB0UmVxdWVzdBIqCg5jaXBoZXJ0ZXh0YmxvYhin/aGhASABKAxSDmNpcGhlcnRleH'
    'RibG9iEhkKBmRyeXJ1bhi43/srIAEoCFIGZHJ5cnVuEkQKD2RyeXJ1bm1vZGlmaWVycxi4/s0L'
    'IAMoDjIXLmttcy5EcnlSdW5Nb2RpZmllclR5cGVSD2RyeXJ1bm1vZGlmaWVycxJRChNlbmNyeX'
    'B0aW9uYWxnb3JpdGhtGIrs82AgASgOMhwua21zLkVuY3J5cHRpb25BbGdvcml0aG1TcGVjUhNl'
    'bmNyeXB0aW9uYWxnb3JpdGhtElwKEWVuY3J5cHRpb25jb250ZXh0GMCkv4gBIAMoCzIqLmttcy'
    '5EZWNyeXB0UmVxdWVzdC5FbmNyeXB0aW9uY29udGV4dEVudHJ5UhFlbmNyeXB0aW9uY29udGV4'
    'dBIkCgtncmFudHRva2VucxiMjYCiASADKAlSC2dyYW50dG9rZW5zEhgKBWtleWlkGKKAyIMBIA'
    'EoCVIFa2V5aWQSNAoJcmVjaXBpZW50GJnI1NQBIAEoCzISLmttcy5SZWNpcGllbnRJbmZvUgly'
    'ZWNpcGllbnQaRAoWRW5jcnlwdGlvbmNvbnRleHRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCg'
    'V2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use decryptResponseDescriptor instead')
const DecryptResponse$json = {
  '1': 'DecryptResponse',
  '2': [
    {
      '1': 'ciphertextforrecipient',
      '3': 75299212,
      '4': 1,
      '5': 12,
      '10': 'ciphertextforrecipient'
    },
    {
      '1': 'encryptionalgorithm',
      '3': 203224586,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'encryptionalgorithm'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
    {'1': 'plaintext', '3': 88342721, '4': 1, '5': 12, '10': 'plaintext'},
  ],
};

/// Descriptor for `DecryptResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List decryptResponseDescriptor = $convert.base64Decode(
    'Cg9EZWNyeXB0UmVzcG9uc2USOQoWY2lwaGVydGV4dGZvcnJlY2lwaWVudBiM8/MjIAEoDFIWY2'
    'lwaGVydGV4dGZvcnJlY2lwaWVudBJRChNlbmNyeXB0aW9uYWxnb3JpdGhtGIrs82AgASgOMhwu'
    'a21zLkVuY3J5cHRpb25BbGdvcml0aG1TcGVjUhNlbmNyeXB0aW9uYWxnb3JpdGhtEhgKBWtleW'
    'lkGKKAyIMBIAEoCVIFa2V5aWQSJwoNa2V5bWF0ZXJpYWxpZBiB8IxGIAEoCVINa2V5bWF0ZXJp'
    'YWxpZBIfCglwbGFpbnRleHQYwYGQKiABKAxSCXBsYWludGV4dA==');

@$core.Deprecated('Use deleteAliasRequestDescriptor instead')
const DeleteAliasRequest$json = {
  '1': 'DeleteAliasRequest',
  '2': [
    {'1': 'aliasname', '3': 313250709, '4': 1, '5': 9, '10': 'aliasname'},
  ],
};

/// Descriptor for `DeleteAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteAliasRequestDescriptor = $convert.base64Decode(
    'ChJEZWxldGVBbGlhc1JlcXVlc3QSIAoJYWxpYXNuYW1lGJWnr5UBIAEoCVIJYWxpYXNuYW1l');

@$core.Deprecated('Use deleteCustomKeyStoreRequestDescriptor instead')
const DeleteCustomKeyStoreRequest$json = {
  '1': 'DeleteCustomKeyStoreRequest',
  '2': [
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
  ],
};

/// Descriptor for `DeleteCustomKeyStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCustomKeyStoreRequestDescriptor =
    $convert.base64Decode(
        'ChtEZWxldGVDdXN0b21LZXlTdG9yZVJlcXVlc3QSLQoQY3VzdG9ta2V5c3RvcmVpZBjErJAqIA'
        'EoCVIQY3VzdG9ta2V5c3RvcmVpZA==');

@$core.Deprecated('Use deleteCustomKeyStoreResponseDescriptor instead')
const DeleteCustomKeyStoreResponse$json = {
  '1': 'DeleteCustomKeyStoreResponse',
};

/// Descriptor for `DeleteCustomKeyStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteCustomKeyStoreResponseDescriptor =
    $convert.base64Decode('ChxEZWxldGVDdXN0b21LZXlTdG9yZVJlc3BvbnNl');

@$core.Deprecated('Use deleteImportedKeyMaterialRequestDescriptor instead')
const DeleteImportedKeyMaterialRequest$json = {
  '1': 'DeleteImportedKeyMaterialRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
  ],
};

/// Descriptor for `DeleteImportedKeyMaterialRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteImportedKeyMaterialRequestDescriptor =
    $convert.base64Decode(
        'CiBEZWxldGVJbXBvcnRlZEtleU1hdGVyaWFsUmVxdWVzdBIYCgVrZXlpZBiigMiDASABKAlSBW'
        'tleWlkEicKDWtleW1hdGVyaWFsaWQYgfCMRiABKAlSDWtleW1hdGVyaWFsaWQ=');

@$core.Deprecated('Use deleteImportedKeyMaterialResponseDescriptor instead')
const DeleteImportedKeyMaterialResponse$json = {
  '1': 'DeleteImportedKeyMaterialResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
  ],
};

/// Descriptor for `DeleteImportedKeyMaterialResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deleteImportedKeyMaterialResponseDescriptor =
    $convert.base64Decode(
        'CiFEZWxldGVJbXBvcnRlZEtleU1hdGVyaWFsUmVzcG9uc2USGAoFa2V5aWQYooDIgwEgASgJUg'
        'VrZXlpZBInCg1rZXltYXRlcmlhbGlkGIHwjEYgASgJUg1rZXltYXRlcmlhbGlk');

@$core.Deprecated('Use dependencyTimeoutExceptionDescriptor instead')
const DependencyTimeoutException$json = {
  '1': 'DependencyTimeoutException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DependencyTimeoutException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dependencyTimeoutExceptionDescriptor =
    $convert.base64Decode(
        'ChpEZXBlbmRlbmN5VGltZW91dEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use deriveSharedSecretRequestDescriptor instead')
const DeriveSharedSecretRequest$json = {
  '1': 'DeriveSharedSecretRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {
      '1': 'keyagreementalgorithm',
      '3': 99147702,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyAgreementAlgorithmSpec',
      '10': 'keyagreementalgorithm'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'publickey', '3': 167335776, '4': 1, '5': 12, '10': 'publickey'},
    {
      '1': 'recipient',
      '3': 445981721,
      '4': 1,
      '5': 11,
      '6': '.kms.RecipientInfo',
      '10': 'recipient'
    },
  ],
};

/// Descriptor for `DeriveSharedSecretRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deriveSharedSecretRequestDescriptor = $convert.base64Decode(
    'ChlEZXJpdmVTaGFyZWRTZWNyZXRSZXF1ZXN0EhkKBmRyeXJ1bhi43/srIAEoCFIGZHJ5cnVuEi'
    'QKC2dyYW50dG9rZW5zGIyNgKIBIAMoCVILZ3JhbnR0b2tlbnMSVwoVa2V5YWdyZWVtZW50YWxn'
    'b3JpdGhtGLa/oy8gASgOMh4ua21zLktleUFncmVlbWVudEFsZ29yaXRobVNwZWNSFWtleWFncm'
    'VlbWVudGFsZ29yaXRobRIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEh8KCXB1YmxpY2tleRjg'
    'ruVPIAEoDFIJcHVibGlja2V5EjQKCXJlY2lwaWVudBiZyNTUASABKAsyEi5rbXMuUmVjaXBpZW'
    '50SW5mb1IJcmVjaXBpZW50');

@$core.Deprecated('Use deriveSharedSecretResponseDescriptor instead')
const DeriveSharedSecretResponse$json = {
  '1': 'DeriveSharedSecretResponse',
  '2': [
    {
      '1': 'ciphertextforrecipient',
      '3': 75299212,
      '4': 1,
      '5': 12,
      '10': 'ciphertextforrecipient'
    },
    {
      '1': 'keyagreementalgorithm',
      '3': 99147702,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyAgreementAlgorithmSpec',
      '10': 'keyagreementalgorithm'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keyorigin',
      '3': 50912311,
      '4': 1,
      '5': 14,
      '6': '.kms.OriginType',
      '10': 'keyorigin'
    },
    {
      '1': 'sharedsecret',
      '3': 382876889,
      '4': 1,
      '5': 12,
      '10': 'sharedsecret'
    },
  ],
};

/// Descriptor for `DeriveSharedSecretResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List deriveSharedSecretResponseDescriptor = $convert.base64Decode(
    'ChpEZXJpdmVTaGFyZWRTZWNyZXRSZXNwb25zZRI5ChZjaXBoZXJ0ZXh0Zm9ycmVjaXBpZW50GI'
    'zz8yMgASgMUhZjaXBoZXJ0ZXh0Zm9ycmVjaXBpZW50ElcKFWtleWFncmVlbWVudGFsZ29yaXRo'
    'bRi2v6MvIAEoDjIeLmttcy5LZXlBZ3JlZW1lbnRBbGdvcml0aG1TcGVjUhVrZXlhZ3JlZW1lbn'
    'RhbGdvcml0aG0SGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBIwCglrZXlvcmlnaW4Yt7ijGCAB'
    'KA4yDy5rbXMuT3JpZ2luVHlwZVIJa2V5b3JpZ2luEiYKDHNoYXJlZHNlY3JldBjZ+ci2ASABKA'
    'xSDHNoYXJlZHNlY3JldA==');

@$core.Deprecated('Use describeCustomKeyStoresRequestDescriptor instead')
const DescribeCustomKeyStoresRequest$json = {
  '1': 'DescribeCustomKeyStoresRequest',
  '2': [
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
    {
      '1': 'customkeystorename',
      '3': 170278046,
      '4': 1,
      '5': 9,
      '10': 'customkeystorename'
    },
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `DescribeCustomKeyStoresRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeCustomKeyStoresRequestDescriptor =
    $convert.base64Decode(
        'Ch5EZXNjcmliZUN1c3RvbUtleVN0b3Jlc1JlcXVlc3QSLQoQY3VzdG9ta2V5c3RvcmVpZBjErJ'
        'AqIAEoCVIQY3VzdG9ta2V5c3RvcmVpZBIxChJjdXN0b21rZXlzdG9yZW5hbWUYnvmYUSABKAlS'
        'EmN1c3RvbWtleXN0b3JlbmFtZRIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0EhkKBm1hcmtlch'
        'i43c0qIAEoCVIGbWFya2Vy');

@$core.Deprecated('Use describeCustomKeyStoresResponseDescriptor instead')
const DescribeCustomKeyStoresResponse$json = {
  '1': 'DescribeCustomKeyStoresResponse',
  '2': [
    {
      '1': 'customkeystores',
      '3': 200763800,
      '4': 3,
      '5': 11,
      '6': '.kms.CustomKeyStoresListEntry',
      '10': 'customkeystores'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `DescribeCustomKeyStoresResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeCustomKeyStoresResponseDescriptor =
    $convert.base64Decode(
        'Ch9EZXNjcmliZUN1c3RvbUtleVN0b3Jlc1Jlc3BvbnNlEkoKD2N1c3RvbWtleXN0b3JlcxiY09'
        '1fIAMoCzIdLmttcy5DdXN0b21LZXlTdG9yZXNMaXN0RW50cnlSD2N1c3RvbWtleXN0b3JlcxIi'
        'CgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG1hcmtlchIfCgl0cnVuY2F0ZWQYyu/YSCABKA'
        'hSCXRydW5jYXRlZA==');

@$core.Deprecated('Use describeKeyRequestDescriptor instead')
const DescribeKeyRequest$json = {
  '1': 'DescribeKeyRequest',
  '2': [
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `DescribeKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeKeyRequestDescriptor = $convert.base64Decode(
    'ChJEZXNjcmliZUtleVJlcXVlc3QSJAoLZ3JhbnR0b2tlbnMYjI2AogEgAygJUgtncmFudHRva2'
    'VucxIYCgVrZXlpZBiigMiDASABKAlSBWtleWlk');

@$core.Deprecated('Use describeKeyResponseDescriptor instead')
const DescribeKeyResponse$json = {
  '1': 'DescribeKeyResponse',
  '2': [
    {
      '1': 'keymetadata',
      '3': 149520202,
      '4': 1,
      '5': 11,
      '6': '.kms.KeyMetadata',
      '10': 'keymetadata'
    },
  ],
};

/// Descriptor for `DescribeKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List describeKeyResponseDescriptor = $convert.base64Decode(
    'ChNEZXNjcmliZUtleVJlc3BvbnNlEjUKC2tleW1ldGFkYXRhGMr+pUcgASgLMhAua21zLktleU'
    '1ldGFkYXRhUgtrZXltZXRhZGF0YQ==');

@$core.Deprecated('Use disableKeyRequestDescriptor instead')
const DisableKeyRequest$json = {
  '1': 'DisableKeyRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `DisableKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableKeyRequestDescriptor = $convert.base64Decode(
    'ChFEaXNhYmxlS2V5UmVxdWVzdBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlk');

@$core.Deprecated('Use disableKeyRotationRequestDescriptor instead')
const DisableKeyRotationRequest$json = {
  '1': 'DisableKeyRotationRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `DisableKeyRotationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disableKeyRotationRequestDescriptor =
    $convert.base64Decode(
        'ChlEaXNhYmxlS2V5Um90YXRpb25SZXF1ZXN0EhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQ=');

@$core.Deprecated('Use disabledExceptionDescriptor instead')
const DisabledException$json = {
  '1': 'DisabledException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DisabledException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disabledExceptionDescriptor = $convert.base64Decode(
    'ChFEaXNhYmxlZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use disconnectCustomKeyStoreRequestDescriptor instead')
const DisconnectCustomKeyStoreRequest$json = {
  '1': 'DisconnectCustomKeyStoreRequest',
  '2': [
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
  ],
};

/// Descriptor for `DisconnectCustomKeyStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disconnectCustomKeyStoreRequestDescriptor =
    $convert.base64Decode(
        'Ch9EaXNjb25uZWN0Q3VzdG9tS2V5U3RvcmVSZXF1ZXN0Ei0KEGN1c3RvbWtleXN0b3JlaWQYxK'
        'yQKiABKAlSEGN1c3RvbWtleXN0b3JlaWQ=');

@$core.Deprecated('Use disconnectCustomKeyStoreResponseDescriptor instead')
const DisconnectCustomKeyStoreResponse$json = {
  '1': 'DisconnectCustomKeyStoreResponse',
};

/// Descriptor for `DisconnectCustomKeyStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List disconnectCustomKeyStoreResponseDescriptor =
    $convert.base64Decode('CiBEaXNjb25uZWN0Q3VzdG9tS2V5U3RvcmVSZXNwb25zZQ==');

@$core.Deprecated('Use dryRunOperationExceptionDescriptor instead')
const DryRunOperationException$json = {
  '1': 'DryRunOperationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `DryRunOperationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List dryRunOperationExceptionDescriptor =
    $convert.base64Decode(
        'ChhEcnlSdW5PcGVyYXRpb25FeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use enableKeyRequestDescriptor instead')
const EnableKeyRequest$json = {
  '1': 'EnableKeyRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `EnableKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableKeyRequestDescriptor = $convert.base64Decode(
    'ChBFbmFibGVLZXlSZXF1ZXN0EhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQ=');

@$core.Deprecated('Use enableKeyRotationRequestDescriptor instead')
const EnableKeyRotationRequest$json = {
  '1': 'EnableKeyRotationRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'rotationperiodindays',
      '3': 118357231,
      '4': 1,
      '5': 5,
      '10': 'rotationperiodindays'
    },
  ],
};

/// Descriptor for `EnableKeyRotationRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List enableKeyRotationRequestDescriptor =
    $convert.base64Decode(
        'ChhFbmFibGVLZXlSb3RhdGlvblJlcXVlc3QSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBI1Ch'
        'Ryb3RhdGlvbnBlcmlvZGluZGF5cxjv+bc4IAEoBVIUcm90YXRpb25wZXJpb2RpbmRheXM=');

@$core.Deprecated('Use encryptRequestDescriptor instead')
const EncryptRequest$json = {
  '1': 'EncryptRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'encryptionalgorithm',
      '3': 203224586,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'encryptionalgorithm'
    },
    {
      '1': 'encryptioncontext',
      '3': 286249536,
      '4': 3,
      '5': 11,
      '6': '.kms.EncryptRequest.EncryptioncontextEntry',
      '10': 'encryptioncontext'
    },
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'plaintext', '3': 88342721, '4': 1, '5': 12, '10': 'plaintext'},
  ],
  '3': [EncryptRequest_EncryptioncontextEntry$json],
};

@$core.Deprecated('Use encryptRequestDescriptor instead')
const EncryptRequest_EncryptioncontextEntry$json = {
  '1': 'EncryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `EncryptRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptRequestDescriptor = $convert.base64Decode(
    'Cg5FbmNyeXB0UmVxdWVzdBIZCgZkcnlydW4YuN/7KyABKAhSBmRyeXJ1bhJRChNlbmNyeXB0aW'
    '9uYWxnb3JpdGhtGIrs82AgASgOMhwua21zLkVuY3J5cHRpb25BbGdvcml0aG1TcGVjUhNlbmNy'
    'eXB0aW9uYWxnb3JpdGhtElwKEWVuY3J5cHRpb25jb250ZXh0GMCkv4gBIAMoCzIqLmttcy5Fbm'
    'NyeXB0UmVxdWVzdC5FbmNyeXB0aW9uY29udGV4dEVudHJ5UhFlbmNyeXB0aW9uY29udGV4dBIk'
    'CgtncmFudHRva2VucxiMjYCiASADKAlSC2dyYW50dG9rZW5zEhgKBWtleWlkGKKAyIMBIAEoCV'
    'IFa2V5aWQSHwoJcGxhaW50ZXh0GMGBkCogASgMUglwbGFpbnRleHQaRAoWRW5jcnlwdGlvbmNv'
    'bnRleHRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use encryptResponseDescriptor instead')
const EncryptResponse$json = {
  '1': 'EncryptResponse',
  '2': [
    {
      '1': 'ciphertextblob',
      '3': 338198183,
      '4': 1,
      '5': 12,
      '10': 'ciphertextblob'
    },
    {
      '1': 'encryptionalgorithm',
      '3': 203224586,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'encryptionalgorithm'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `EncryptResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List encryptResponseDescriptor = $convert.base64Decode(
    'Cg9FbmNyeXB0UmVzcG9uc2USKgoOY2lwaGVydGV4dGJsb2IYp/2hoQEgASgMUg5jaXBoZXJ0ZX'
    'h0YmxvYhJRChNlbmNyeXB0aW9uYWxnb3JpdGhtGIrs82AgASgOMhwua21zLkVuY3J5cHRpb25B'
    'bGdvcml0aG1TcGVjUhNlbmNyeXB0aW9uYWxnb3JpdGhtEhgKBWtleWlkGKKAyIMBIAEoCVIFa2'
    'V5aWQ=');

@$core.Deprecated('Use expiredImportTokenExceptionDescriptor instead')
const ExpiredImportTokenException$json = {
  '1': 'ExpiredImportTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `ExpiredImportTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List expiredImportTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChtFeHBpcmVkSW1wb3J0VG9rZW5FeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use generateDataKeyPairRequestDescriptor instead')
const GenerateDataKeyPairRequest$json = {
  '1': 'GenerateDataKeyPairRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'encryptioncontext',
      '3': 286249536,
      '4': 3,
      '5': 11,
      '6': '.kms.GenerateDataKeyPairRequest.EncryptioncontextEntry',
      '10': 'encryptioncontext'
    },
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keypairspec',
      '3': 142696380,
      '4': 1,
      '5': 14,
      '6': '.kms.DataKeyPairSpec',
      '10': 'keypairspec'
    },
    {
      '1': 'recipient',
      '3': 445981721,
      '4': 1,
      '5': 11,
      '6': '.kms.RecipientInfo',
      '10': 'recipient'
    },
  ],
  '3': [GenerateDataKeyPairRequest_EncryptioncontextEntry$json],
};

@$core.Deprecated('Use generateDataKeyPairRequestDescriptor instead')
const GenerateDataKeyPairRequest_EncryptioncontextEntry$json = {
  '1': 'EncryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GenerateDataKeyPairRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDataKeyPairRequestDescriptor = $convert.base64Decode(
    'ChpHZW5lcmF0ZURhdGFLZXlQYWlyUmVxdWVzdBIZCgZkcnlydW4YuN/7KyABKAhSBmRyeXJ1bh'
    'JoChFlbmNyeXB0aW9uY29udGV4dBjApL+IASADKAsyNi5rbXMuR2VuZXJhdGVEYXRhS2V5UGFp'
    'clJlcXVlc3QuRW5jcnlwdGlvbmNvbnRleHRFbnRyeVIRZW5jcnlwdGlvbmNvbnRleHQSJAoLZ3'
    'JhbnR0b2tlbnMYjI2AogEgAygJUgtncmFudHRva2VucxIYCgVrZXlpZBiigMiDASABKAlSBWtl'
    'eWlkEjkKC2tleXBhaXJzcGVjGLy/hUQgASgOMhQua21zLkRhdGFLZXlQYWlyU3BlY1ILa2V5cG'
    'FpcnNwZWMSNAoJcmVjaXBpZW50GJnI1NQBIAEoCzISLmttcy5SZWNpcGllbnRJbmZvUglyZWNp'
    'cGllbnQaRAoWRW5jcnlwdGlvbmNvbnRleHRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YW'
    'x1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core.Deprecated('Use generateDataKeyPairResponseDescriptor instead')
const GenerateDataKeyPairResponse$json = {
  '1': 'GenerateDataKeyPairResponse',
  '2': [
    {
      '1': 'ciphertextforrecipient',
      '3': 75299212,
      '4': 1,
      '5': 12,
      '10': 'ciphertextforrecipient'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
    {
      '1': 'keypairspec',
      '3': 142696380,
      '4': 1,
      '5': 14,
      '6': '.kms.DataKeyPairSpec',
      '10': 'keypairspec'
    },
    {
      '1': 'privatekeyciphertextblob',
      '3': 295238401,
      '4': 1,
      '5': 12,
      '10': 'privatekeyciphertextblob'
    },
    {
      '1': 'privatekeyplaintext',
      '3': 422534247,
      '4': 1,
      '5': 12,
      '10': 'privatekeyplaintext'
    },
    {'1': 'publickey', '3': 167335776, '4': 1, '5': 12, '10': 'publickey'},
  ],
};

/// Descriptor for `GenerateDataKeyPairResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDataKeyPairResponseDescriptor = $convert.base64Decode(
    'ChtHZW5lcmF0ZURhdGFLZXlQYWlyUmVzcG9uc2USOQoWY2lwaGVydGV4dGZvcnJlY2lwaWVudB'
    'iM8/MjIAEoDFIWY2lwaGVydGV4dGZvcnJlY2lwaWVudBIYCgVrZXlpZBiigMiDASABKAlSBWtl'
    'eWlkEicKDWtleW1hdGVyaWFsaWQYgfCMRiABKAlSDWtleW1hdGVyaWFsaWQSOQoLa2V5cGFpcn'
    'NwZWMYvL+FRCABKA4yFC5rbXMuRGF0YUtleVBhaXJTcGVjUgtrZXlwYWlyc3BlYxI+Chhwcml2'
    'YXRla2V5Y2lwaGVydGV4dGJsb2IYgfbjjAEgASgMUhhwcml2YXRla2V5Y2lwaGVydGV4dGJsb2'
    'ISNAoTcHJpdmF0ZWtleXBsYWludGV4dBjnuL3JASABKAxSE3ByaXZhdGVrZXlwbGFpbnRleHQS'
    'HwoJcHVibGlja2V5GOCu5U8gASgMUglwdWJsaWNrZXk=');

@$core.Deprecated(
    'Use generateDataKeyPairWithoutPlaintextRequestDescriptor instead')
const GenerateDataKeyPairWithoutPlaintextRequest$json = {
  '1': 'GenerateDataKeyPairWithoutPlaintextRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'encryptioncontext',
      '3': 286249536,
      '4': 3,
      '5': 11,
      '6':
          '.kms.GenerateDataKeyPairWithoutPlaintextRequest.EncryptioncontextEntry',
      '10': 'encryptioncontext'
    },
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keypairspec',
      '3': 142696380,
      '4': 1,
      '5': 14,
      '6': '.kms.DataKeyPairSpec',
      '10': 'keypairspec'
    },
  ],
  '3': [GenerateDataKeyPairWithoutPlaintextRequest_EncryptioncontextEntry$json],
};

@$core.Deprecated(
    'Use generateDataKeyPairWithoutPlaintextRequestDescriptor instead')
const GenerateDataKeyPairWithoutPlaintextRequest_EncryptioncontextEntry$json = {
  '1': 'EncryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GenerateDataKeyPairWithoutPlaintextRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDataKeyPairWithoutPlaintextRequestDescriptor = $convert.base64Decode(
    'CipHZW5lcmF0ZURhdGFLZXlQYWlyV2l0aG91dFBsYWludGV4dFJlcXVlc3QSGQoGZHJ5cnVuGL'
    'jf+ysgASgIUgZkcnlydW4SeAoRZW5jcnlwdGlvbmNvbnRleHQYwKS/iAEgAygLMkYua21zLkdl'
    'bmVyYXRlRGF0YUtleVBhaXJXaXRob3V0UGxhaW50ZXh0UmVxdWVzdC5FbmNyeXB0aW9uY29udG'
    'V4dEVudHJ5UhFlbmNyeXB0aW9uY29udGV4dBIkCgtncmFudHRva2VucxiMjYCiASADKAlSC2dy'
    'YW50dG9rZW5zEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSOQoLa2V5cGFpcnNwZWMYvL+FRC'
    'ABKA4yFC5rbXMuRGF0YUtleVBhaXJTcGVjUgtrZXlwYWlyc3BlYxpEChZFbmNyeXB0aW9uY29u'
    'dGV4dEVudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated(
    'Use generateDataKeyPairWithoutPlaintextResponseDescriptor instead')
const GenerateDataKeyPairWithoutPlaintextResponse$json = {
  '1': 'GenerateDataKeyPairWithoutPlaintextResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
    {
      '1': 'keypairspec',
      '3': 142696380,
      '4': 1,
      '5': 14,
      '6': '.kms.DataKeyPairSpec',
      '10': 'keypairspec'
    },
    {
      '1': 'privatekeyciphertextblob',
      '3': 295238401,
      '4': 1,
      '5': 12,
      '10': 'privatekeyciphertextblob'
    },
    {'1': 'publickey', '3': 167335776, '4': 1, '5': 12, '10': 'publickey'},
  ],
};

/// Descriptor for `GenerateDataKeyPairWithoutPlaintextResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    generateDataKeyPairWithoutPlaintextResponseDescriptor =
    $convert.base64Decode(
        'CitHZW5lcmF0ZURhdGFLZXlQYWlyV2l0aG91dFBsYWludGV4dFJlc3BvbnNlEhgKBWtleWlkGK'
        'KAyIMBIAEoCVIFa2V5aWQSJwoNa2V5bWF0ZXJpYWxpZBiB8IxGIAEoCVINa2V5bWF0ZXJpYWxp'
        'ZBI5CgtrZXlwYWlyc3BlYxi8v4VEIAEoDjIULmttcy5EYXRhS2V5UGFpclNwZWNSC2tleXBhaX'
        'JzcGVjEj4KGHByaXZhdGVrZXljaXBoZXJ0ZXh0YmxvYhiB9uOMASABKAxSGHByaXZhdGVrZXlj'
        'aXBoZXJ0ZXh0YmxvYhIfCglwdWJsaWNrZXkY4K7lTyABKAxSCXB1YmxpY2tleQ==');

@$core.Deprecated('Use generateDataKeyRequestDescriptor instead')
const GenerateDataKeyRequest$json = {
  '1': 'GenerateDataKeyRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'encryptioncontext',
      '3': 286249536,
      '4': 3,
      '5': 11,
      '6': '.kms.GenerateDataKeyRequest.EncryptioncontextEntry',
      '10': 'encryptioncontext'
    },
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keyspec',
      '3': 138220928,
      '4': 1,
      '5': 14,
      '6': '.kms.DataKeySpec',
      '10': 'keyspec'
    },
    {
      '1': 'numberofbytes',
      '3': 277086201,
      '4': 1,
      '5': 5,
      '10': 'numberofbytes'
    },
    {
      '1': 'recipient',
      '3': 445981721,
      '4': 1,
      '5': 11,
      '6': '.kms.RecipientInfo',
      '10': 'recipient'
    },
  ],
  '3': [GenerateDataKeyRequest_EncryptioncontextEntry$json],
};

@$core.Deprecated('Use generateDataKeyRequestDescriptor instead')
const GenerateDataKeyRequest_EncryptioncontextEntry$json = {
  '1': 'EncryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GenerateDataKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDataKeyRequestDescriptor = $convert.base64Decode(
    'ChZHZW5lcmF0ZURhdGFLZXlSZXF1ZXN0EhkKBmRyeXJ1bhi43/srIAEoCFIGZHJ5cnVuEmQKEW'
    'VuY3J5cHRpb25jb250ZXh0GMCkv4gBIAMoCzIyLmttcy5HZW5lcmF0ZURhdGFLZXlSZXF1ZXN0'
    'LkVuY3J5cHRpb25jb250ZXh0RW50cnlSEWVuY3J5cHRpb25jb250ZXh0EiQKC2dyYW50dG9rZW'
    '5zGIyNgKIBIAMoCVILZ3JhbnR0b2tlbnMSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBItCgdr'
    'ZXlzcGVjGICr9EEgASgOMhAua21zLkRhdGFLZXlTcGVjUgdrZXlzcGVjEigKDW51bWJlcm9mYn'
    'l0ZXMY+f+PhAEgASgFUg1udW1iZXJvZmJ5dGVzEjQKCXJlY2lwaWVudBiZyNTUASABKAsyEi5r'
    'bXMuUmVjaXBpZW50SW5mb1IJcmVjaXBpZW50GkQKFkVuY3J5cHRpb25jb250ZXh0RW50cnkSEA'
    'oDa2V5GAEgASgJUgNrZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4AQ==');

@$core.Deprecated('Use generateDataKeyResponseDescriptor instead')
const GenerateDataKeyResponse$json = {
  '1': 'GenerateDataKeyResponse',
  '2': [
    {
      '1': 'ciphertextblob',
      '3': 338198183,
      '4': 1,
      '5': 12,
      '10': 'ciphertextblob'
    },
    {
      '1': 'ciphertextforrecipient',
      '3': 75299212,
      '4': 1,
      '5': 12,
      '10': 'ciphertextforrecipient'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
    {'1': 'plaintext', '3': 88342721, '4': 1, '5': 12, '10': 'plaintext'},
  ],
};

/// Descriptor for `GenerateDataKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDataKeyResponseDescriptor = $convert.base64Decode(
    'ChdHZW5lcmF0ZURhdGFLZXlSZXNwb25zZRIqCg5jaXBoZXJ0ZXh0YmxvYhin/aGhASABKAxSDm'
    'NpcGhlcnRleHRibG9iEjkKFmNpcGhlcnRleHRmb3JyZWNpcGllbnQYjPPzIyABKAxSFmNpcGhl'
    'cnRleHRmb3JyZWNpcGllbnQSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBInCg1rZXltYXRlcm'
    'lhbGlkGIHwjEYgASgJUg1rZXltYXRlcmlhbGlkEh8KCXBsYWludGV4dBjBgZAqIAEoDFIJcGxh'
    'aW50ZXh0');

@$core
    .Deprecated('Use generateDataKeyWithoutPlaintextRequestDescriptor instead')
const GenerateDataKeyWithoutPlaintextRequest$json = {
  '1': 'GenerateDataKeyWithoutPlaintextRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'encryptioncontext',
      '3': 286249536,
      '4': 3,
      '5': 11,
      '6': '.kms.GenerateDataKeyWithoutPlaintextRequest.EncryptioncontextEntry',
      '10': 'encryptioncontext'
    },
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keyspec',
      '3': 138220928,
      '4': 1,
      '5': 14,
      '6': '.kms.DataKeySpec',
      '10': 'keyspec'
    },
    {
      '1': 'numberofbytes',
      '3': 277086201,
      '4': 1,
      '5': 5,
      '10': 'numberofbytes'
    },
  ],
  '3': [GenerateDataKeyWithoutPlaintextRequest_EncryptioncontextEntry$json],
};

@$core
    .Deprecated('Use generateDataKeyWithoutPlaintextRequestDescriptor instead')
const GenerateDataKeyWithoutPlaintextRequest_EncryptioncontextEntry$json = {
  '1': 'EncryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GenerateDataKeyWithoutPlaintextRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDataKeyWithoutPlaintextRequestDescriptor = $convert.base64Decode(
    'CiZHZW5lcmF0ZURhdGFLZXlXaXRob3V0UGxhaW50ZXh0UmVxdWVzdBIZCgZkcnlydW4YuN/7Ky'
    'ABKAhSBmRyeXJ1bhJ0ChFlbmNyeXB0aW9uY29udGV4dBjApL+IASADKAsyQi5rbXMuR2VuZXJh'
    'dGVEYXRhS2V5V2l0aG91dFBsYWludGV4dFJlcXVlc3QuRW5jcnlwdGlvbmNvbnRleHRFbnRyeV'
    'IRZW5jcnlwdGlvbmNvbnRleHQSJAoLZ3JhbnR0b2tlbnMYjI2AogEgAygJUgtncmFudHRva2Vu'
    'cxIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEi0KB2tleXNwZWMYgKv0QSABKA4yEC5rbXMuRG'
    'F0YUtleVNwZWNSB2tleXNwZWMSKAoNbnVtYmVyb2ZieXRlcxj5/4+EASABKAVSDW51bWJlcm9m'
    'Ynl0ZXMaRAoWRW5jcnlwdGlvbmNvbnRleHRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YW'
    'x1ZRgCIAEoCVIFdmFsdWU6AjgB');

@$core
    .Deprecated('Use generateDataKeyWithoutPlaintextResponseDescriptor instead')
const GenerateDataKeyWithoutPlaintextResponse$json = {
  '1': 'GenerateDataKeyWithoutPlaintextResponse',
  '2': [
    {
      '1': 'ciphertextblob',
      '3': 338198183,
      '4': 1,
      '5': 12,
      '10': 'ciphertextblob'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
  ],
};

/// Descriptor for `GenerateDataKeyWithoutPlaintextResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateDataKeyWithoutPlaintextResponseDescriptor =
    $convert.base64Decode(
        'CidHZW5lcmF0ZURhdGFLZXlXaXRob3V0UGxhaW50ZXh0UmVzcG9uc2USKgoOY2lwaGVydGV4dG'
        'Jsb2IYp/2hoQEgASgMUg5jaXBoZXJ0ZXh0YmxvYhIYCgVrZXlpZBiigMiDASABKAlSBWtleWlk'
        'EicKDWtleW1hdGVyaWFsaWQYgfCMRiABKAlSDWtleW1hdGVyaWFsaWQ=');

@$core.Deprecated('Use generateMacRequestDescriptor instead')
const GenerateMacRequest$json = {
  '1': 'GenerateMacRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'macalgorithm',
      '3': 253519878,
      '4': 1,
      '5': 14,
      '6': '.kms.MacAlgorithmSpec',
      '10': 'macalgorithm'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 12, '10': 'message'},
  ],
};

/// Descriptor for `GenerateMacRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateMacRequestDescriptor = $convert.base64Decode(
    'ChJHZW5lcmF0ZU1hY1JlcXVlc3QSGQoGZHJ5cnVuGLjf+ysgASgIUgZkcnlydW4SJAoLZ3Jhbn'
    'R0b2tlbnMYjI2AogEgAygJUgtncmFudHRva2VucxIYCgVrZXlpZBiigMiDASABKAlSBWtleWlk'
    'EjwKDG1hY2FsZ29yaXRobRiG0PF4IAEoDjIVLmttcy5NYWNBbGdvcml0aG1TcGVjUgxtYWNhbG'
    'dvcml0aG0SGwoHbWVzc2FnZRiFs7twIAEoDFIHbWVzc2FnZQ==');

@$core.Deprecated('Use generateMacResponseDescriptor instead')
const GenerateMacResponse$json = {
  '1': 'GenerateMacResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'mac', '3': 296223945, '4': 1, '5': 12, '10': 'mac'},
    {
      '1': 'macalgorithm',
      '3': 253519878,
      '4': 1,
      '5': 14,
      '6': '.kms.MacAlgorithmSpec',
      '10': 'macalgorithm'
    },
  ],
};

/// Descriptor for `GenerateMacResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateMacResponseDescriptor = $convert.base64Decode(
    'ChNHZW5lcmF0ZU1hY1Jlc3BvbnNlEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSFAoDbWFjGM'
    'mJoI0BIAEoDFIDbWFjEjwKDG1hY2FsZ29yaXRobRiG0PF4IAEoDjIVLmttcy5NYWNBbGdvcml0'
    'aG1TcGVjUgxtYWNhbGdvcml0aG0=');

@$core.Deprecated('Use generateRandomRequestDescriptor instead')
const GenerateRandomRequest$json = {
  '1': 'GenerateRandomRequest',
  '2': [
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
    {
      '1': 'numberofbytes',
      '3': 277086201,
      '4': 1,
      '5': 5,
      '10': 'numberofbytes'
    },
    {
      '1': 'recipient',
      '3': 445981721,
      '4': 1,
      '5': 11,
      '6': '.kms.RecipientInfo',
      '10': 'recipient'
    },
  ],
};

/// Descriptor for `GenerateRandomRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateRandomRequestDescriptor = $convert.base64Decode(
    'ChVHZW5lcmF0ZVJhbmRvbVJlcXVlc3QSLQoQY3VzdG9ta2V5c3RvcmVpZBjErJAqIAEoCVIQY3'
    'VzdG9ta2V5c3RvcmVpZBIoCg1udW1iZXJvZmJ5dGVzGPn/j4QBIAEoBVINbnVtYmVyb2ZieXRl'
    'cxI0CglyZWNpcGllbnQYmcjU1AEgASgLMhIua21zLlJlY2lwaWVudEluZm9SCXJlY2lwaWVudA'
    '==');

@$core.Deprecated('Use generateRandomResponseDescriptor instead')
const GenerateRandomResponse$json = {
  '1': 'GenerateRandomResponse',
  '2': [
    {
      '1': 'ciphertextforrecipient',
      '3': 75299212,
      '4': 1,
      '5': 12,
      '10': 'ciphertextforrecipient'
    },
    {'1': 'plaintext', '3': 88342721, '4': 1, '5': 12, '10': 'plaintext'},
  ],
};

/// Descriptor for `GenerateRandomResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List generateRandomResponseDescriptor = $convert.base64Decode(
    'ChZHZW5lcmF0ZVJhbmRvbVJlc3BvbnNlEjkKFmNpcGhlcnRleHRmb3JyZWNpcGllbnQYjPPzIy'
    'ABKAxSFmNpcGhlcnRleHRmb3JyZWNpcGllbnQSHwoJcGxhaW50ZXh0GMGBkCogASgMUglwbGFp'
    'bnRleHQ=');

@$core.Deprecated('Use getKeyPolicyRequestDescriptor instead')
const GetKeyPolicyRequest$json = {
  '1': 'GetKeyPolicyRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `GetKeyPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyPolicyRequestDescriptor = $convert.base64Decode(
    'ChNHZXRLZXlQb2xpY3lSZXF1ZXN0EhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSIQoKcG9saW'
    'N5bmFtZRi99Yd/IAEoCVIKcG9saWN5bmFtZQ==');

@$core.Deprecated('Use getKeyPolicyResponseDescriptor instead')
const GetKeyPolicyResponse$json = {
  '1': 'GetKeyPolicyResponse',
  '2': [
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `GetKeyPolicyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyPolicyResponseDescriptor = $convert.base64Decode(
    'ChRHZXRLZXlQb2xpY3lSZXNwb25zZRIaCgZwb2xpY3kYoO/w4AEgASgJUgZwb2xpY3kSIQoKcG'
    '9saWN5bmFtZRi99Yd/IAEoCVIKcG9saWN5bmFtZQ==');

@$core.Deprecated('Use getKeyRotationStatusRequestDescriptor instead')
const GetKeyRotationStatusRequest$json = {
  '1': 'GetKeyRotationStatusRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `GetKeyRotationStatusRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyRotationStatusRequestDescriptor =
    $convert.base64Decode(
        'ChtHZXRLZXlSb3RhdGlvblN0YXR1c1JlcXVlc3QSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZA'
        '==');

@$core.Deprecated('Use getKeyRotationStatusResponseDescriptor instead')
const GetKeyRotationStatusResponse$json = {
  '1': 'GetKeyRotationStatusResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keyrotationenabled',
      '3': 525956616,
      '4': 1,
      '5': 8,
      '10': 'keyrotationenabled'
    },
    {
      '1': 'nextrotationdate',
      '3': 192035355,
      '4': 1,
      '5': 9,
      '10': 'nextrotationdate'
    },
    {
      '1': 'ondemandrotationstartdate',
      '3': 360279652,
      '4': 1,
      '5': 9,
      '10': 'ondemandrotationstartdate'
    },
    {
      '1': 'rotationperiodindays',
      '3': 118357231,
      '4': 1,
      '5': 5,
      '10': 'rotationperiodindays'
    },
  ],
};

/// Descriptor for `GetKeyRotationStatusResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getKeyRotationStatusResponseDescriptor = $convert.base64Decode(
    'ChxHZXRLZXlSb3RhdGlvblN0YXR1c1Jlc3BvbnNlEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aW'
    'QSMgoSa2V5cm90YXRpb25lbmFibGVkGIjs5foBIAEoCFISa2V5cm90YXRpb25lbmFibGVkEi0K'
    'EG5leHRyb3RhdGlvbmRhdGUYm/TIWyABKAlSEG5leHRyb3RhdGlvbmRhdGUSQAoZb25kZW1hbm'
    'Ryb3RhdGlvbnN0YXJ0ZGF0ZRjk3OWrASABKAlSGW9uZGVtYW5kcm90YXRpb25zdGFydGRhdGUS'
    'NQoUcm90YXRpb25wZXJpb2RpbmRheXMY7/m3OCABKAVSFHJvdGF0aW9ucGVyaW9kaW5kYXlz');

@$core.Deprecated('Use getParametersForImportRequestDescriptor instead')
const GetParametersForImportRequest$json = {
  '1': 'GetParametersForImportRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'wrappingalgorithm',
      '3': 164083247,
      '4': 1,
      '5': 14,
      '6': '.kms.AlgorithmSpec',
      '10': 'wrappingalgorithm'
    },
    {
      '1': 'wrappingkeyspec',
      '3': 521228924,
      '4': 1,
      '5': 14,
      '6': '.kms.WrappingKeySpec',
      '10': 'wrappingkeyspec'
    },
  ],
};

/// Descriptor for `GetParametersForImportRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParametersForImportRequestDescriptor = $convert.base64Decode(
    'Ch1HZXRQYXJhbWV0ZXJzRm9ySW1wb3J0UmVxdWVzdBIYCgVrZXlpZBiigMiDASABKAlSBWtleW'
    'lkEkMKEXdyYXBwaW5nYWxnb3JpdGhtGK/snk4gASgOMhIua21zLkFsZ29yaXRobVNwZWNSEXdy'
    'YXBwaW5nYWxnb3JpdGhtEkIKD3dyYXBwaW5na2V5c3BlYxj8pMX4ASABKA4yFC5rbXMuV3JhcH'
    'BpbmdLZXlTcGVjUg93cmFwcGluZ2tleXNwZWM=');

@$core.Deprecated('Use getParametersForImportResponseDescriptor instead')
const GetParametersForImportResponse$json = {
  '1': 'GetParametersForImportResponse',
  '2': [
    {'1': 'importtoken', '3': 461726162, '4': 1, '5': 12, '10': 'importtoken'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'parametersvalidto',
      '3': 394913717,
      '4': 1,
      '5': 9,
      '10': 'parametersvalidto'
    },
    {'1': 'publickey', '3': 167335776, '4': 1, '5': 12, '10': 'publickey'},
  ],
};

/// Descriptor for `GetParametersForImportResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getParametersForImportResponseDescriptor =
    $convert.base64Decode(
        'Ch5HZXRQYXJhbWV0ZXJzRm9ySW1wb3J0UmVzcG9uc2USJAoLaW1wb3J0dG9rZW4Y0sOV3AEgAS'
        'gMUgtpbXBvcnR0b2tlbhIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEjAKEXBhcmFtZXRlcnN2'
        'YWxpZHRvGLXPp7wBIAEoCVIRcGFyYW1ldGVyc3ZhbGlkdG8SHwoJcHVibGlja2V5GOCu5U8gAS'
        'gMUglwdWJsaWNrZXk=');

@$core.Deprecated('Use getPublicKeyRequestDescriptor instead')
const GetPublicKeyRequest$json = {
  '1': 'GetPublicKeyRequest',
  '2': [
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `GetPublicKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyRequestDescriptor = $convert.base64Decode(
    'ChNHZXRQdWJsaWNLZXlSZXF1ZXN0EiQKC2dyYW50dG9rZW5zGIyNgKIBIAMoCVILZ3JhbnR0b2'
    'tlbnMSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZA==');

@$core.Deprecated('Use getPublicKeyResponseDescriptor instead')
const GetPublicKeyResponse$json = {
  '1': 'GetPublicKeyResponse',
  '2': [
    {
      '1': 'customermasterkeyspec',
      '3': 472470930,
      '4': 1,
      '5': 14,
      '6': '.kms.CustomerMasterKeySpec',
      '10': 'customermasterkeyspec'
    },
    {
      '1': 'encryptionalgorithms',
      '3': 194511375,
      '4': 3,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'encryptionalgorithms'
    },
    {
      '1': 'keyagreementalgorithms',
      '3': 328746163,
      '4': 3,
      '5': 14,
      '6': '.kms.KeyAgreementAlgorithmSpec',
      '10': 'keyagreementalgorithms'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keyspec',
      '3': 138220928,
      '4': 1,
      '5': 14,
      '6': '.kms.KeySpec',
      '10': 'keyspec'
    },
    {
      '1': 'keyusage',
      '3': 357216772,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyUsageType',
      '10': 'keyusage'
    },
    {'1': 'publickey', '3': 167335776, '4': 1, '5': 12, '10': 'publickey'},
    {
      '1': 'signingalgorithms',
      '3': 508241975,
      '4': 3,
      '5': 14,
      '6': '.kms.SigningAlgorithmSpec',
      '10': 'signingalgorithms'
    },
  ],
};

/// Descriptor for `GetPublicKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List getPublicKeyResponseDescriptor = $convert.base64Decode(
    'ChRHZXRQdWJsaWNLZXlSZXNwb25zZRJUChVjdXN0b21lcm1hc3RlcmtleXNwZWMYkqul4QEgAS'
    'gOMhoua21zLkN1c3RvbWVyTWFzdGVyS2V5U3BlY1IVY3VzdG9tZXJtYXN0ZXJrZXlzcGVjElMK'
    'FGVuY3J5cHRpb25hbGdvcml0aG1zGI+E4FwgAygOMhwua21zLkVuY3J5cHRpb25BbGdvcml0aG'
    '1TcGVjUhRlbmNyeXB0aW9uYWxnb3JpdGhtcxJaChZrZXlhZ3JlZW1lbnRhbGdvcml0aG1zGLOJ'
    '4ZwBIAMoDjIeLmttcy5LZXlBZ3JlZW1lbnRBbGdvcml0aG1TcGVjUhZrZXlhZ3JlZW1lbnRhbG'
    'dvcml0aG1zEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSKQoHa2V5c3BlYxiAq/RBIAEoDjIM'
    'Lmttcy5LZXlTcGVjUgdrZXlzcGVjEjEKCGtleXVzYWdlGITkqqoBIAEoDjIRLmttcy5LZXlVc2'
    'FnZVR5cGVSCGtleXVzYWdlEh8KCXB1YmxpY2tleRjgruVPIAEoDFIJcHVibGlja2V5EksKEXNp'
    'Z25pbmdhbGdvcml0aG1zGLfQrPIBIAMoDjIZLmttcy5TaWduaW5nQWxnb3JpdGhtU3BlY1IRc2'
    'lnbmluZ2FsZ29yaXRobXM=');

@$core.Deprecated('Use grantConstraintsDescriptor instead')
const GrantConstraints$json = {
  '1': 'GrantConstraints',
  '2': [
    {
      '1': 'encryptioncontextequals',
      '3': 68403171,
      '4': 3,
      '5': 11,
      '6': '.kms.GrantConstraints.EncryptioncontextequalsEntry',
      '10': 'encryptioncontextequals'
    },
    {
      '1': 'encryptioncontextsubset',
      '3': 72310514,
      '4': 3,
      '5': 11,
      '6': '.kms.GrantConstraints.EncryptioncontextsubsetEntry',
      '10': 'encryptioncontextsubset'
    },
  ],
  '3': [
    GrantConstraints_EncryptioncontextequalsEntry$json,
    GrantConstraints_EncryptioncontextsubsetEntry$json
  ],
};

@$core.Deprecated('Use grantConstraintsDescriptor instead')
const GrantConstraints_EncryptioncontextequalsEntry$json = {
  '1': 'EncryptioncontextequalsEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use grantConstraintsDescriptor instead')
const GrantConstraints_EncryptioncontextsubsetEntry$json = {
  '1': 'EncryptioncontextsubsetEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `GrantConstraints`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grantConstraintsDescriptor = $convert.base64Decode(
    'ChBHcmFudENvbnN0cmFpbnRzEm8KF2VuY3J5cHRpb25jb250ZXh0ZXF1YWxzGOP/ziAgAygLMj'
    'Iua21zLkdyYW50Q29uc3RyYWludHMuRW5jcnlwdGlvbmNvbnRleHRlcXVhbHNFbnRyeVIXZW5j'
    'cnlwdGlvbmNvbnRleHRlcXVhbHMSbwoXZW5jcnlwdGlvbmNvbnRleHRzdWJzZXQY8r29IiADKA'
    'syMi5rbXMuR3JhbnRDb25zdHJhaW50cy5FbmNyeXB0aW9uY29udGV4dHN1YnNldEVudHJ5Uhdl'
    'bmNyeXB0aW9uY29udGV4dHN1YnNldBpKChxFbmNyeXB0aW9uY29udGV4dGVxdWFsc0VudHJ5Eh'
    'AKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAEaSgocRW5jcnlwdGlv'
    'bmNvbnRleHRzdWJzZXRFbnRyeRIQCgNrZXkYASABKAlSA2tleRIUCgV2YWx1ZRgCIAEoCVIFdm'
    'FsdWU6AjgB');

@$core.Deprecated('Use grantListEntryDescriptor instead')
const GrantListEntry$json = {
  '1': 'GrantListEntry',
  '2': [
    {
      '1': 'constraints',
      '3': 302297388,
      '4': 1,
      '5': 11,
      '6': '.kms.GrantConstraints',
      '10': 'constraints'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {'1': 'grantid', '3': 66852281, '4': 1, '5': 9, '10': 'grantid'},
    {
      '1': 'granteeprincipal',
      '3': 234727364,
      '4': 1,
      '5': 9,
      '10': 'granteeprincipal'
    },
    {
      '1': 'issuingaccount',
      '3': 47662575,
      '4': 1,
      '5': 9,
      '10': 'issuingaccount'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'name', '3': 266367751, '4': 1, '5': 9, '10': 'name'},
    {
      '1': 'operations',
      '3': 126776656,
      '4': 3,
      '5': 14,
      '6': '.kms.GrantOperation',
      '10': 'operations'
    },
    {
      '1': 'retiringprincipal',
      '3': 49541086,
      '4': 1,
      '5': 9,
      '10': 'retiringprincipal'
    },
  ],
};

/// Descriptor for `GrantListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List grantListEntryDescriptor = $convert.base64Decode(
    'Cg5HcmFudExpc3RFbnRyeRI7Cgtjb25zdHJhaW50cxis4pKQASABKAsyFS5rbXMuR3JhbnRDb2'
    '5zdHJhaW50c1ILY29uc3RyYWludHMSJgoMY3JlYXRpb25kYXRlGOHYt4kBIAEoCVIMY3JlYXRp'
    'b25kYXRlEhsKB2dyYW50aWQYuavwHyABKAlSB2dyYW50aWQSLQoQZ3JhbnRlZXByaW5jaXBhbB'
    'jEz/ZvIAEoCVIQZ3JhbnRlZXByaW5jaXBhbBIpCg5pc3N1aW5nYWNjb3VudBjvi90WIAEoCVIO'
    'aXNzdWluZ2FjY291bnQSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBIVCgRuYW1lGIfmgX8gAS'
    'gJUgRuYW1lEjYKCm9wZXJhdGlvbnMY0Oq5PCADKA4yEy5rbXMuR3JhbnRPcGVyYXRpb25SCm9w'
    'ZXJhdGlvbnMSLwoRcmV0aXJpbmdwcmluY2lwYWwY3t/PFyABKAlSEXJldGlyaW5ncHJpbmNpcG'
    'Fs');

@$core.Deprecated('Use importKeyMaterialRequestDescriptor instead')
const ImportKeyMaterialRequest$json = {
  '1': 'ImportKeyMaterialRequest',
  '2': [
    {
      '1': 'encryptedkeymaterial',
      '3': 3185968,
      '4': 1,
      '5': 12,
      '10': 'encryptedkeymaterial'
    },
    {
      '1': 'expirationmodel',
      '3': 113933558,
      '4': 1,
      '5': 14,
      '6': '.kms.ExpirationModelType',
      '10': 'expirationmodel'
    },
    {'1': 'importtoken', '3': 461726162, '4': 1, '5': 12, '10': 'importtoken'},
    {
      '1': 'importtype',
      '3': 331980349,
      '4': 1,
      '5': 14,
      '6': '.kms.ImportType',
      '10': 'importtype'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialdescription',
      '3': 277153546,
      '4': 1,
      '5': 9,
      '10': 'keymaterialdescription'
    },
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
    {'1': 'validto', '3': 522718673, '4': 1, '5': 9, '10': 'validto'},
  ],
};

/// Descriptor for `ImportKeyMaterialRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importKeyMaterialRequestDescriptor = $convert.base64Decode(
    'ChhJbXBvcnRLZXlNYXRlcmlhbFJlcXVlc3QSNQoUZW5jcnlwdGVka2V5bWF0ZXJpYWwYsLrCAS'
    'ABKAxSFGVuY3J5cHRlZGtleW1hdGVyaWFsEkUKD2V4cGlyYXRpb25tb2RlbBj2+ak2IAEoDjIY'
    'Lmttcy5FeHBpcmF0aW9uTW9kZWxUeXBlUg9leHBpcmF0aW9ubW9kZWwSJAoLaW1wb3J0dG9rZW'
    '4Y0sOV3AEgASgMUgtpbXBvcnR0b2tlbhIzCgppbXBvcnR0eXBlGL28pp4BIAEoDjIPLmttcy5J'
    'bXBvcnRUeXBlUgppbXBvcnR0eXBlEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSOgoWa2V5bW'
    'F0ZXJpYWxkZXNjcmlwdGlvbhiKjpSEASABKAlSFmtleW1hdGVyaWFsZGVzY3JpcHRpb24SJwoN'
    'a2V5bWF0ZXJpYWxpZBiB8IxGIAEoCVINa2V5bWF0ZXJpYWxpZBIcCgd2YWxpZHRvGNGboPkBIA'
    'EoCVIHdmFsaWR0bw==');

@$core.Deprecated('Use importKeyMaterialResponseDescriptor instead')
const ImportKeyMaterialResponse$json = {
  '1': 'ImportKeyMaterialResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
  ],
};

/// Descriptor for `ImportKeyMaterialResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List importKeyMaterialResponseDescriptor =
    $convert.base64Decode(
        'ChlJbXBvcnRLZXlNYXRlcmlhbFJlc3BvbnNlEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSJw'
        'oNa2V5bWF0ZXJpYWxpZBiB8IxGIAEoCVINa2V5bWF0ZXJpYWxpZA==');

@$core.Deprecated('Use incorrectKeyExceptionDescriptor instead')
const IncorrectKeyException$json = {
  '1': 'IncorrectKeyException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IncorrectKeyException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List incorrectKeyExceptionDescriptor = $convert.base64Decode(
    'ChVJbmNvcnJlY3RLZXlFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use incorrectKeyMaterialExceptionDescriptor instead')
const IncorrectKeyMaterialException$json = {
  '1': 'IncorrectKeyMaterialException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IncorrectKeyMaterialException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List incorrectKeyMaterialExceptionDescriptor =
    $convert.base64Decode(
        'Ch1JbmNvcnJlY3RLZXlNYXRlcmlhbEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use incorrectTrustAnchorExceptionDescriptor instead')
const IncorrectTrustAnchorException$json = {
  '1': 'IncorrectTrustAnchorException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `IncorrectTrustAnchorException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List incorrectTrustAnchorExceptionDescriptor =
    $convert.base64Decode(
        'Ch1JbmNvcnJlY3RUcnVzdEFuY2hvckV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use invalidAliasNameExceptionDescriptor instead')
const InvalidAliasNameException$json = {
  '1': 'InvalidAliasNameException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidAliasNameException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidAliasNameExceptionDescriptor =
    $convert.base64Decode(
        'ChlJbnZhbGlkQWxpYXNOYW1lRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use invalidArnExceptionDescriptor instead')
const InvalidArnException$json = {
  '1': 'InvalidArnException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidArnException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidArnExceptionDescriptor =
    $convert.base64Decode(
        'ChNJbnZhbGlkQXJuRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use invalidCiphertextExceptionDescriptor instead')
const InvalidCiphertextException$json = {
  '1': 'InvalidCiphertextException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidCiphertextException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidCiphertextExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkQ2lwaGVydGV4dEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidGrantIdExceptionDescriptor instead')
const InvalidGrantIdException$json = {
  '1': 'InvalidGrantIdException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidGrantIdException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidGrantIdExceptionDescriptor =
    $convert.base64Decode(
        'ChdJbnZhbGlkR3JhbnRJZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use invalidGrantTokenExceptionDescriptor instead')
const InvalidGrantTokenException$json = {
  '1': 'InvalidGrantTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidGrantTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidGrantTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChpJbnZhbGlkR3JhbnRUb2tlbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYW'
        'dl');

@$core.Deprecated('Use invalidImportTokenExceptionDescriptor instead')
const InvalidImportTokenException$json = {
  '1': 'InvalidImportTokenException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidImportTokenException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidImportTokenExceptionDescriptor =
    $convert.base64Decode(
        'ChtJbnZhbGlkSW1wb3J0VG9rZW5FeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use invalidKeyUsageExceptionDescriptor instead')
const InvalidKeyUsageException$json = {
  '1': 'InvalidKeyUsageException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidKeyUsageException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidKeyUsageExceptionDescriptor =
    $convert.base64Decode(
        'ChhJbnZhbGlkS2V5VXNhZ2VFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ'
        '==');

@$core.Deprecated('Use invalidMarkerExceptionDescriptor instead')
const InvalidMarkerException$json = {
  '1': 'InvalidMarkerException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `InvalidMarkerException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List invalidMarkerExceptionDescriptor =
    $convert.base64Decode(
        'ChZJbnZhbGlkTWFya2VyRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use kMSInternalExceptionDescriptor instead')
const KMSInternalException$json = {
  '1': 'KMSInternalException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSInternalException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSInternalExceptionDescriptor =
    $convert.base64Decode(
        'ChRLTVNJbnRlcm5hbEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use kMSInvalidMacExceptionDescriptor instead')
const KMSInvalidMacException$json = {
  '1': 'KMSInvalidMacException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSInvalidMacException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSInvalidMacExceptionDescriptor =
    $convert.base64Decode(
        'ChZLTVNJbnZhbGlkTWFjRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use kMSInvalidSignatureExceptionDescriptor instead')
const KMSInvalidSignatureException$json = {
  '1': 'KMSInvalidSignatureException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KMSInvalidSignatureException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List kMSInvalidSignatureExceptionDescriptor =
    $convert.base64Decode(
        'ChxLTVNJbnZhbGlkU2lnbmF0dXJlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3'
        'NhZ2U=');

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

@$core.Deprecated('Use keyListEntryDescriptor instead')
const KeyListEntry$json = {
  '1': 'KeyListEntry',
  '2': [
    {'1': 'keyarn', '3': 418055012, '4': 1, '5': 9, '10': 'keyarn'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `KeyListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyListEntryDescriptor = $convert.base64Decode(
    'CgxLZXlMaXN0RW50cnkSGgoGa2V5YXJuGOSGrMcBIAEoCVIGa2V5YXJuEhgKBWtleWlkGKKAyI'
    'MBIAEoCVIFa2V5aWQ=');

@$core.Deprecated('Use keyMetadataDescriptor instead')
const KeyMetadata$json = {
  '1': 'KeyMetadata',
  '2': [
    {'1': 'awsaccountid', '3': 370093421, '4': 1, '5': 9, '10': 'awsaccountid'},
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {
      '1': 'cloudhsmclusterid',
      '3': 56498754,
      '4': 1,
      '5': 9,
      '10': 'cloudhsmclusterid'
    },
    {'1': 'creationdate', '3': 288222305, '4': 1, '5': 9, '10': 'creationdate'},
    {
      '1': 'currentkeymaterialid',
      '3': 183721586,
      '4': 1,
      '5': 9,
      '10': 'currentkeymaterialid'
    },
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
    {
      '1': 'customermasterkeyspec',
      '3': 472470930,
      '4': 1,
      '5': 14,
      '6': '.kms.CustomerMasterKeySpec',
      '10': 'customermasterkeyspec'
    },
    {'1': 'deletiondate', '3': 347845564, '4': 1, '5': 9, '10': 'deletiondate'},
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'enabled', '3': 478602303, '4': 1, '5': 8, '10': 'enabled'},
    {
      '1': 'encryptionalgorithms',
      '3': 194511375,
      '4': 3,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'encryptionalgorithms'
    },
    {
      '1': 'expirationmodel',
      '3': 113933558,
      '4': 1,
      '5': 14,
      '6': '.kms.ExpirationModelType',
      '10': 'expirationmodel'
    },
    {
      '1': 'keyagreementalgorithms',
      '3': 328746163,
      '4': 3,
      '5': 14,
      '6': '.kms.KeyAgreementAlgorithmSpec',
      '10': 'keyagreementalgorithms'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymanager',
      '3': 432733972,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyManagerType',
      '10': 'keymanager'
    },
    {
      '1': 'keyspec',
      '3': 138220928,
      '4': 1,
      '5': 14,
      '6': '.kms.KeySpec',
      '10': 'keyspec'
    },
    {
      '1': 'keystate',
      '3': 27894226,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyState',
      '10': 'keystate'
    },
    {
      '1': 'keyusage',
      '3': 357216772,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyUsageType',
      '10': 'keyusage'
    },
    {
      '1': 'macalgorithms',
      '3': 532181443,
      '4': 3,
      '5': 14,
      '6': '.kms.MacAlgorithmSpec',
      '10': 'macalgorithms'
    },
    {'1': 'multiregion', '3': 405769103, '4': 1, '5': 8, '10': 'multiregion'},
    {
      '1': 'multiregionconfiguration',
      '3': 460103813,
      '4': 1,
      '5': 11,
      '6': '.kms.MultiRegionConfiguration',
      '10': 'multiregionconfiguration'
    },
    {
      '1': 'origin',
      '3': 529732720,
      '4': 1,
      '5': 14,
      '6': '.kms.OriginType',
      '10': 'origin'
    },
    {
      '1': 'pendingdeletionwindowindays',
      '3': 480447795,
      '4': 1,
      '5': 5,
      '10': 'pendingdeletionwindowindays'
    },
    {
      '1': 'signingalgorithms',
      '3': 508241975,
      '4': 3,
      '5': 14,
      '6': '.kms.SigningAlgorithmSpec',
      '10': 'signingalgorithms'
    },
    {'1': 'validto', '3': 522718673, '4': 1, '5': 9, '10': 'validto'},
    {
      '1': 'xkskeyconfiguration',
      '3': 359766455,
      '4': 1,
      '5': 11,
      '6': '.kms.XksKeyConfigurationType',
      '10': 'xkskeyconfiguration'
    },
  ],
};

/// Descriptor for `KeyMetadata`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyMetadataDescriptor = $convert.base64Decode(
    'CgtLZXlNZXRhZGF0YRImCgxhd3NhY2NvdW50aWQY7dq8sAEgASgJUgxhd3NhY2NvdW50aWQSFA'
    'oDYXJuGJ2b7b8BIAEoCVIDYXJuEi8KEWNsb3VkaHNtY2x1c3RlcmlkGMK0+BogASgJUhFjbG91'
    'ZGhzbWNsdXN0ZXJpZBImCgxjcmVhdGlvbmRhdGUY4di3iQEgASgJUgxjcmVhdGlvbmRhdGUSNQ'
    'oUY3VycmVudGtleW1hdGVyaWFsaWQY8rzNVyABKAlSFGN1cnJlbnRrZXltYXRlcmlhbGlkEi0K'
    'EGN1c3RvbWtleXN0b3JlaWQYxKyQKiABKAlSEGN1c3RvbWtleXN0b3JlaWQSVAoVY3VzdG9tZX'
    'JtYXN0ZXJrZXlzcGVjGJKrpeEBIAEoDjIaLmttcy5DdXN0b21lck1hc3RlcktleVNwZWNSFWN1'
    'c3RvbWVybWFzdGVya2V5c3BlYxImCgxkZWxldGlvbmRhdGUYvOfupQEgASgJUgxkZWxldGlvbm'
    'RhdGUSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2Rlc2NyaXB0aW9uEhwKB2VuYWJsZWQYv8ib'
    '5AEgASgIUgdlbmFibGVkElMKFGVuY3J5cHRpb25hbGdvcml0aG1zGI+E4FwgAygOMhwua21zLk'
    'VuY3J5cHRpb25BbGdvcml0aG1TcGVjUhRlbmNyeXB0aW9uYWxnb3JpdGhtcxJFCg9leHBpcmF0'
    'aW9ubW9kZWwY9vmpNiABKA4yGC5rbXMuRXhwaXJhdGlvbk1vZGVsVHlwZVIPZXhwaXJhdGlvbm'
    '1vZGVsEloKFmtleWFncmVlbWVudGFsZ29yaXRobXMYs4nhnAEgAygOMh4ua21zLktleUFncmVl'
    'bWVudEFsZ29yaXRobVNwZWNSFmtleWFncmVlbWVudGFsZ29yaXRobXMSGAoFa2V5aWQYooDIgw'
    'EgASgJUgVrZXlpZBI3CgprZXltYW5hZ2VyGJT+q84BIAEoDjITLmttcy5LZXlNYW5hZ2VyVHlw'
    'ZVIKa2V5bWFuYWdlchIpCgdrZXlzcGVjGICr9EEgASgOMgwua21zLktleVNwZWNSB2tleXNwZW'
    'MSLAoIa2V5c3RhdGUY0sOmDSABKA4yDS5rbXMuS2V5U3RhdGVSCGtleXN0YXRlEjEKCGtleXVz'
    'YWdlGITkqqoBIAEoDjIRLmttcy5LZXlVc2FnZVR5cGVSCGtleXVzYWdlEj8KDW1hY2FsZ29yaX'
    'RobXMYw+Ph/QEgAygOMhUua21zLk1hY0FsZ29yaXRobVNwZWNSDW1hY2FsZ29yaXRobXMSJAoL'
    'bXVsdGlyZWdpb24Yj5e+wQEgASgIUgttdWx0aXJlZ2lvbhJdChhtdWx0aXJlZ2lvbmNvbmZpZ3'
    'VyYXRpb24YhcGy2wEgASgLMh0ua21zLk11bHRpUmVnaW9uQ29uZmlndXJhdGlvblIYbXVsdGly'
    'ZWdpb25jb25maWd1cmF0aW9uEisKBm9yaWdpbhjwqMz8ASABKA4yDy5rbXMuT3JpZ2luVHlwZV'
    'IGb3JpZ2luEkQKG3BlbmRpbmdkZWxldGlvbndpbmRvd2luZGF5cxizmozlASABKAVSG3BlbmRp'
    'bmdkZWxldGlvbndpbmRvd2luZGF5cxJLChFzaWduaW5nYWxnb3JpdGhtcxi30KzyASADKA4yGS'
    '5rbXMuU2lnbmluZ0FsZ29yaXRobVNwZWNSEXNpZ25pbmdhbGdvcml0aG1zEhwKB3ZhbGlkdG8Y'
    '0Zug+QEgASgJUgd2YWxpZHRvElIKE3hrc2tleWNvbmZpZ3VyYXRpb24Yt7PGqwEgASgLMhwua2'
    '1zLlhrc0tleUNvbmZpZ3VyYXRpb25UeXBlUhN4a3NrZXljb25maWd1cmF0aW9u');

@$core.Deprecated('Use keyUnavailableExceptionDescriptor instead')
const KeyUnavailableException$json = {
  '1': 'KeyUnavailableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `KeyUnavailableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List keyUnavailableExceptionDescriptor =
    $convert.base64Decode(
        'ChdLZXlVbmF2YWlsYWJsZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

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

@$core.Deprecated('Use listAliasesRequestDescriptor instead')
const ListAliasesRequest$json = {
  '1': 'ListAliasesRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListAliasesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAliasesRequestDescriptor = $convert.base64Decode(
    'ChJMaXN0QWxpYXNlc1JlcXVlc3QSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBIYCgVsaW1pdB'
    'jVldnEASABKAVSBWxpbWl0EhkKBm1hcmtlchi43c0qIAEoCVIGbWFya2Vy');

@$core.Deprecated('Use listAliasesResponseDescriptor instead')
const ListAliasesResponse$json = {
  '1': 'ListAliasesResponse',
  '2': [
    {
      '1': 'aliases',
      '3': 476693696,
      '4': 3,
      '5': 11,
      '6': '.kms.AliasListEntry',
      '10': 'aliases'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ListAliasesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listAliasesResponseDescriptor = $convert.base64Decode(
    'ChNMaXN0QWxpYXNlc1Jlc3BvbnNlEjEKB2FsaWFzZXMYwImn4wEgAygLMhMua21zLkFsaWFzTG'
    'lzdEVudHJ5UgdhbGlhc2VzEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2VyEh8K'
    'CXRydW5jYXRlZBjK79hIIAEoCFIJdHJ1bmNhdGVk');

@$core.Deprecated('Use listGrantsRequestDescriptor instead')
const ListGrantsRequest$json = {
  '1': 'ListGrantsRequest',
  '2': [
    {'1': 'grantid', '3': 66852281, '4': 1, '5': 9, '10': 'grantid'},
    {
      '1': 'granteeprincipal',
      '3': 234727364,
      '4': 1,
      '5': 9,
      '10': 'granteeprincipal'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListGrantsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGrantsRequestDescriptor = $convert.base64Decode(
    'ChFMaXN0R3JhbnRzUmVxdWVzdBIbCgdncmFudGlkGLmr8B8gASgJUgdncmFudGlkEi0KEGdyYW'
    '50ZWVwcmluY2lwYWwYxM/2byABKAlSEGdyYW50ZWVwcmluY2lwYWwSGAoFa2V5aWQYooDIgwEg'
    'ASgJUgVrZXlpZBIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0EhkKBm1hcmtlchi43c0qIAEoCV'
    'IGbWFya2Vy');

@$core.Deprecated('Use listGrantsResponseDescriptor instead')
const ListGrantsResponse$json = {
  '1': 'ListGrantsResponse',
  '2': [
    {
      '1': 'grants',
      '3': 226910555,
      '4': 3,
      '5': 11,
      '6': '.kms.GrantListEntry',
      '10': 'grants'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ListGrantsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listGrantsResponseDescriptor = $convert.base64Decode(
    'ChJMaXN0R3JhbnRzUmVzcG9uc2USLgoGZ3JhbnRzGNvCmWwgAygLMhMua21zLkdyYW50TGlzdE'
    'VudHJ5UgZncmFudHMSIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leHRtYXJrZXISHwoJdHJ1'
    'bmNhdGVkGMrv2EggASgIUgl0cnVuY2F0ZWQ=');

@$core.Deprecated('Use listKeyPoliciesRequestDescriptor instead')
const ListKeyPoliciesRequest$json = {
  '1': 'ListKeyPoliciesRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListKeyPoliciesRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyPoliciesRequestDescriptor =
    $convert.base64Decode(
        'ChZMaXN0S2V5UG9saWNpZXNSZXF1ZXN0EhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSGAoFbG'
        'ltaXQY1ZXZxAEgASgFUgVsaW1pdBIZCgZtYXJrZXIYuN3NKiABKAlSBm1hcmtlcg==');

@$core.Deprecated('Use listKeyPoliciesResponseDescriptor instead')
const ListKeyPoliciesResponse$json = {
  '1': 'ListKeyPoliciesResponse',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'policynames', '3': 264098782, '4': 3, '5': 9, '10': 'policynames'},
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ListKeyPoliciesResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyPoliciesResponseDescriptor = $convert.base64Decode(
    'ChdMaXN0S2V5UG9saWNpZXNSZXNwb25zZRIiCgpuZXh0bWFya2VyGKOBrv0BIAEoCVIKbmV4dG'
    '1hcmtlchIjCgtwb2xpY3luYW1lcxjep/d9IAMoCVILcG9saWN5bmFtZXMSHwoJdHJ1bmNhdGVk'
    'GMrv2EggASgIUgl0cnVuY2F0ZWQ=');

@$core.Deprecated('Use listKeyRotationsRequestDescriptor instead')
const ListKeyRotationsRequest$json = {
  '1': 'ListKeyRotationsRequest',
  '2': [
    {
      '1': 'includekeymaterial',
      '3': 531559428,
      '4': 1,
      '5': 14,
      '6': '.kms.IncludeKeyMaterial',
      '10': 'includekeymaterial'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListKeyRotationsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyRotationsRequestDescriptor = $convert.base64Decode(
    'ChdMaXN0S2V5Um90YXRpb25zUmVxdWVzdBJLChJpbmNsdWRla2V5bWF0ZXJpYWwYhOi7/QEgAS'
    'gOMhcua21zLkluY2x1ZGVLZXlNYXRlcmlhbFISaW5jbHVkZWtleW1hdGVyaWFsEhgKBWtleWlk'
    'GKKAyIMBIAEoCVIFa2V5aWQSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBIZCgZtYXJrZXIYuN'
    '3NKiABKAlSBm1hcmtlcg==');

@$core.Deprecated('Use listKeyRotationsResponseDescriptor instead')
const ListKeyRotationsResponse$json = {
  '1': 'ListKeyRotationsResponse',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {
      '1': 'rotations',
      '3': 24209381,
      '4': 3,
      '5': 11,
      '6': '.kms.RotationsListEntry',
      '10': 'rotations'
    },
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ListKeyRotationsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeyRotationsResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0S2V5Um90YXRpb25zUmVzcG9uc2USIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leH'
    'RtYXJrZXISOAoJcm90YXRpb25zGOXPxQsgAygLMhcua21zLlJvdGF0aW9uc0xpc3RFbnRyeVIJ'
    'cm90YXRpb25zEh8KCXRydW5jYXRlZBjK79hIIAEoCFIJdHJ1bmNhdGVk');

@$core.Deprecated('Use listKeysRequestDescriptor instead')
const ListKeysRequest$json = {
  '1': 'ListKeysRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListKeysRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeysRequestDescriptor = $convert.base64Decode(
    'Cg9MaXN0S2V5c1JlcXVlc3QSGAoFbGltaXQY1ZXZxAEgASgFUgVsaW1pdBIZCgZtYXJrZXIYuN'
    '3NKiABKAlSBm1hcmtlcg==');

@$core.Deprecated('Use listKeysResponseDescriptor instead')
const ListKeysResponse$json = {
  '1': 'ListKeysResponse',
  '2': [
    {
      '1': 'keys',
      '3': 2831086,
      '4': 3,
      '5': 11,
      '6': '.kms.KeyListEntry',
      '10': 'keys'
    },
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ListKeysResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listKeysResponseDescriptor = $convert.base64Decode(
    'ChBMaXN0S2V5c1Jlc3BvbnNlEigKBGtleXMY7uWsASADKAsyES5rbXMuS2V5TGlzdEVudHJ5Ug'
    'RrZXlzEiIKCm5leHRtYXJrZXIYo4Gu/QEgASgJUgpuZXh0bWFya2VyEh8KCXRydW5jYXRlZBjK'
    '79hIIAEoCFIJdHJ1bmNhdGVk');

@$core.Deprecated('Use listResourceTagsRequestDescriptor instead')
const ListResourceTagsRequest$json = {
  '1': 'ListResourceTagsRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
  ],
};

/// Descriptor for `ListResourceTagsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceTagsRequestDescriptor =
    $convert.base64Decode(
        'ChdMaXN0UmVzb3VyY2VUYWdzUmVxdWVzdBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEhgKBW'
        'xpbWl0GNWV2cQBIAEoBVIFbGltaXQSGQoGbWFya2VyGLjdzSogASgJUgZtYXJrZXI=');

@$core.Deprecated('Use listResourceTagsResponseDescriptor instead')
const ListResourceTagsResponse$json = {
  '1': 'ListResourceTagsResponse',
  '2': [
    {'1': 'nextmarker', '3': 531333283, '4': 1, '5': 9, '10': 'nextmarker'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kms.Tag',
      '10': 'tags'
    },
    {'1': 'truncated', '3': 152451018, '4': 1, '5': 8, '10': 'truncated'},
  ],
};

/// Descriptor for `ListResourceTagsResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listResourceTagsResponseDescriptor = $convert.base64Decode(
    'ChhMaXN0UmVzb3VyY2VUYWdzUmVzcG9uc2USIgoKbmV4dG1hcmtlchijga79ASABKAlSCm5leH'
    'RtYXJrZXISIAoEdGFncxjBwfa1ASADKAsyCC5rbXMuVGFnUgR0YWdzEh8KCXRydW5jYXRlZBjK'
    '79hIIAEoCFIJdHJ1bmNhdGVk');

@$core.Deprecated('Use listRetirableGrantsRequestDescriptor instead')
const ListRetirableGrantsRequest$json = {
  '1': 'ListRetirableGrantsRequest',
  '2': [
    {'1': 'limit', '3': 412502741, '4': 1, '5': 5, '10': 'limit'},
    {'1': 'marker', '3': 89353912, '4': 1, '5': 9, '10': 'marker'},
    {
      '1': 'retiringprincipal',
      '3': 49541086,
      '4': 1,
      '5': 9,
      '10': 'retiringprincipal'
    },
  ],
};

/// Descriptor for `ListRetirableGrantsRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List listRetirableGrantsRequestDescriptor =
    $convert.base64Decode(
        'ChpMaXN0UmV0aXJhYmxlR3JhbnRzUmVxdWVzdBIYCgVsaW1pdBjVldnEASABKAVSBWxpbWl0Eh'
        'kKBm1hcmtlchi43c0qIAEoCVIGbWFya2VyEi8KEXJldGlyaW5ncHJpbmNpcGFsGN7fzxcgASgJ'
        'UhFyZXRpcmluZ3ByaW5jaXBhbA==');

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

@$core.Deprecated('Use multiRegionConfigurationDescriptor instead')
const MultiRegionConfiguration$json = {
  '1': 'MultiRegionConfiguration',
  '2': [
    {
      '1': 'multiregionkeytype',
      '3': 483927198,
      '4': 1,
      '5': 14,
      '6': '.kms.MultiRegionKeyType',
      '10': 'multiregionkeytype'
    },
    {
      '1': 'primarykey',
      '3': 174019825,
      '4': 1,
      '5': 11,
      '6': '.kms.MultiRegionKey',
      '10': 'primarykey'
    },
    {
      '1': 'replicakeys',
      '3': 266801624,
      '4': 3,
      '5': 11,
      '6': '.kms.MultiRegionKey',
      '10': 'replicakeys'
    },
  ],
};

/// Descriptor for `MultiRegionConfiguration`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List multiRegionConfigurationDescriptor = $convert.base64Decode(
    'ChhNdWx0aVJlZ2lvbkNvbmZpZ3VyYXRpb24SSwoSbXVsdGlyZWdpb25rZXl0eXBlGJ7J4OYBIA'
    'EoDjIXLmttcy5NdWx0aVJlZ2lvbktleVR5cGVSEm11bHRpcmVnaW9ua2V5dHlwZRI2Cgpwcmlt'
    'YXJ5a2V5GPGp/VIgASgLMhMua21zLk11bHRpUmVnaW9uS2V5UgpwcmltYXJ5a2V5EjgKC3JlcG'
    'xpY2FrZXlzGNijnH8gAygLMhMua21zLk11bHRpUmVnaW9uS2V5UgtyZXBsaWNha2V5cw==');

@$core.Deprecated('Use multiRegionKeyDescriptor instead')
const MultiRegionKey$json = {
  '1': 'MultiRegionKey',
  '2': [
    {'1': 'arn', '3': 402345373, '4': 1, '5': 9, '10': 'arn'},
    {'1': 'region', '3': 154040478, '4': 1, '5': 9, '10': 'region'},
  ],
};

/// Descriptor for `MultiRegionKey`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List multiRegionKeyDescriptor = $convert.base64Decode(
    'Cg5NdWx0aVJlZ2lvbktleRIUCgNhcm4YnZvtvwEgASgJUgNhcm4SGQoGcmVnaW9uGJ7xuUkgAS'
    'gJUgZyZWdpb24=');

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

@$core.Deprecated('Use putKeyPolicyRequestDescriptor instead')
const PutKeyPolicyRequest$json = {
  '1': 'PutKeyPolicyRequest',
  '2': [
    {
      '1': 'bypasspolicylockoutsafetycheck',
      '3': 177450851,
      '4': 1,
      '5': 8,
      '10': 'bypasspolicylockoutsafetycheck'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {'1': 'policyname', '3': 266468029, '4': 1, '5': 9, '10': 'policyname'},
  ],
};

/// Descriptor for `PutKeyPolicyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List putKeyPolicyRequestDescriptor = $convert.base64Decode(
    'ChNQdXRLZXlQb2xpY3lSZXF1ZXN0EkkKHmJ5cGFzc3BvbGljeWxvY2tvdXRzYWZldHljaGVjax'
    'jj3s5UIAEoCFIeYnlwYXNzcG9saWN5bG9ja291dHNhZmV0eWNoZWNrEhgKBWtleWlkGKKAyIMB'
    'IAEoCVIFa2V5aWQSGgoGcG9saWN5GKDv8OABIAEoCVIGcG9saWN5EiEKCnBvbGljeW5hbWUYvf'
    'WHfyABKAlSCnBvbGljeW5hbWU=');

@$core.Deprecated('Use reEncryptRequestDescriptor instead')
const ReEncryptRequest$json = {
  '1': 'ReEncryptRequest',
  '2': [
    {
      '1': 'ciphertextblob',
      '3': 338198183,
      '4': 1,
      '5': 12,
      '10': 'ciphertextblob'
    },
    {
      '1': 'destinationencryptionalgorithm',
      '3': 3500944,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'destinationencryptionalgorithm'
    },
    {
      '1': 'destinationencryptioncontext',
      '3': 116435710,
      '4': 3,
      '5': 11,
      '6': '.kms.ReEncryptRequest.DestinationencryptioncontextEntry',
      '10': 'destinationencryptioncontext'
    },
    {
      '1': 'destinationkeyid',
      '3': 396219964,
      '4': 1,
      '5': 9,
      '10': 'destinationkeyid'
    },
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {
      '1': 'dryrunmodifiers',
      '3': 24346424,
      '4': 3,
      '5': 14,
      '6': '.kms.DryRunModifierType',
      '10': 'dryrunmodifiers'
    },
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {
      '1': 'sourceencryptionalgorithm',
      '3': 283215847,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'sourceencryptionalgorithm'
    },
    {
      '1': 'sourceencryptioncontext',
      '3': 203785001,
      '4': 3,
      '5': 11,
      '6': '.kms.ReEncryptRequest.SourceencryptioncontextEntry',
      '10': 'sourceencryptioncontext'
    },
    {'1': 'sourcekeyid', '3': 137105771, '4': 1, '5': 9, '10': 'sourcekeyid'},
  ],
  '3': [
    ReEncryptRequest_DestinationencryptioncontextEntry$json,
    ReEncryptRequest_SourceencryptioncontextEntry$json
  ],
};

@$core.Deprecated('Use reEncryptRequestDescriptor instead')
const ReEncryptRequest_DestinationencryptioncontextEntry$json = {
  '1': 'DestinationencryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

@$core.Deprecated('Use reEncryptRequestDescriptor instead')
const ReEncryptRequest_SourceencryptioncontextEntry$json = {
  '1': 'SourceencryptioncontextEntry',
  '2': [
    {'1': 'key', '3': 1, '4': 1, '5': 9, '10': 'key'},
    {'1': 'value', '3': 2, '4': 1, '5': 9, '10': 'value'},
  ],
  '7': {'7': true},
};

/// Descriptor for `ReEncryptRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reEncryptRequestDescriptor = $convert.base64Decode(
    'ChBSZUVuY3J5cHRSZXF1ZXN0EioKDmNpcGhlcnRleHRibG9iGKf9oaEBIAEoDFIOY2lwaGVydG'
    'V4dGJsb2ISZwoeZGVzdGluYXRpb25lbmNyeXB0aW9uYWxnb3JpdGhtGJDX1QEgASgOMhwua21z'
    'LkVuY3J5cHRpb25BbGdvcml0aG1TcGVjUh5kZXN0aW5hdGlvbmVuY3J5cHRpb25hbGdvcml0aG'
    '0SfgocZGVzdGluYXRpb25lbmNyeXB0aW9uY29udGV4dBj+1cI3IAMoCzI3Lmttcy5SZUVuY3J5'
    'cHRSZXF1ZXN0LkRlc3RpbmF0aW9uZW5jcnlwdGlvbmNvbnRleHRFbnRyeVIcZGVzdGluYXRpb2'
    '5lbmNyeXB0aW9uY29udGV4dBIuChBkZXN0aW5hdGlvbmtleWlkGLys97wBIAEoCVIQZGVzdGlu'
    'YXRpb25rZXlpZBIZCgZkcnlydW4YuN/7KyABKAhSBmRyeXJ1bhJECg9kcnlydW5tb2RpZmllcn'
    'MYuP7NCyADKA4yFy5rbXMuRHJ5UnVuTW9kaWZpZXJUeXBlUg9kcnlydW5tb2RpZmllcnMSJAoL'
    'Z3JhbnR0b2tlbnMYjI2AogEgAygJUgtncmFudHRva2VucxJeChlzb3VyY2VlbmNyeXB0aW9uYW'
    'xnb3JpdGhtGOePhocBIAEoDjIcLmttcy5FbmNyeXB0aW9uQWxnb3JpdGhtU3BlY1IZc291cmNl'
    'ZW5jcnlwdGlvbmFsZ29yaXRobRJvChdzb3VyY2VlbmNyeXB0aW9uY29udGV4dBiphpZhIAMoCz'
    'IyLmttcy5SZUVuY3J5cHRSZXF1ZXN0LlNvdXJjZWVuY3J5cHRpb25jb250ZXh0RW50cnlSF3Nv'
    'dXJjZWVuY3J5cHRpb25jb250ZXh0EiMKC3NvdXJjZWtleWlkGOuisEEgASgJUgtzb3VyY2VrZX'
    'lpZBpPCiFEZXN0aW5hdGlvbmVuY3J5cHRpb25jb250ZXh0RW50cnkSEAoDa2V5GAEgASgJUgNr'
    'ZXkSFAoFdmFsdWUYAiABKAlSBXZhbHVlOgI4ARpKChxTb3VyY2VlbmNyeXB0aW9uY29udGV4dE'
    'VudHJ5EhAKA2tleRgBIAEoCVIDa2V5EhQKBXZhbHVlGAIgASgJUgV2YWx1ZToCOAE=');

@$core.Deprecated('Use reEncryptResponseDescriptor instead')
const ReEncryptResponse$json = {
  '1': 'ReEncryptResponse',
  '2': [
    {
      '1': 'ciphertextblob',
      '3': 338198183,
      '4': 1,
      '5': 12,
      '10': 'ciphertextblob'
    },
    {
      '1': 'destinationencryptionalgorithm',
      '3': 3500944,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'destinationencryptionalgorithm'
    },
    {
      '1': 'destinationkeymaterialid',
      '3': 469834039,
      '4': 1,
      '5': 9,
      '10': 'destinationkeymaterialid'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'sourceencryptionalgorithm',
      '3': 283215847,
      '4': 1,
      '5': 14,
      '6': '.kms.EncryptionAlgorithmSpec',
      '10': 'sourceencryptionalgorithm'
    },
    {'1': 'sourcekeyid', '3': 137105771, '4': 1, '5': 9, '10': 'sourcekeyid'},
    {
      '1': 'sourcekeymaterialid',
      '3': 34789220,
      '4': 1,
      '5': 9,
      '10': 'sourcekeymaterialid'
    },
  ],
};

/// Descriptor for `ReEncryptResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List reEncryptResponseDescriptor = $convert.base64Decode(
    'ChFSZUVuY3J5cHRSZXNwb25zZRIqCg5jaXBoZXJ0ZXh0YmxvYhin/aGhASABKAxSDmNpcGhlcn'
    'RleHRibG9iEmcKHmRlc3RpbmF0aW9uZW5jcnlwdGlvbmFsZ29yaXRobRiQ19UBIAEoDjIcLmtt'
    'cy5FbmNyeXB0aW9uQWxnb3JpdGhtU3BlY1IeZGVzdGluYXRpb25lbmNyeXB0aW9uYWxnb3JpdG'
    'htEj4KGGRlc3RpbmF0aW9ua2V5bWF0ZXJpYWxpZBi3soTgASABKAlSGGRlc3RpbmF0aW9ua2V5'
    'bWF0ZXJpYWxpZBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEl4KGXNvdXJjZWVuY3J5cHRpb2'
    '5hbGdvcml0aG0Y54+GhwEgASgOMhwua21zLkVuY3J5cHRpb25BbGdvcml0aG1TcGVjUhlzb3Vy'
    'Y2VlbmNyeXB0aW9uYWxnb3JpdGhtEiMKC3NvdXJjZWtleWlkGOuisEEgASgJUgtzb3VyY2VrZX'
    'lpZBIzChNzb3VyY2VrZXltYXRlcmlhbGlkGOSuyxAgASgJUhNzb3VyY2VrZXltYXRlcmlhbGlk');

@$core.Deprecated('Use recipientInfoDescriptor instead')
const RecipientInfo$json = {
  '1': 'RecipientInfo',
  '2': [
    {
      '1': 'attestationdocument',
      '3': 217786849,
      '4': 1,
      '5': 12,
      '10': 'attestationdocument'
    },
    {
      '1': 'keyencryptionalgorithm',
      '3': 478234803,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyEncryptionMechanism',
      '10': 'keyencryptionalgorithm'
    },
  ],
};

/// Descriptor for `RecipientInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List recipientInfoDescriptor = $convert.base64Decode(
    'Cg1SZWNpcGllbnRJbmZvEjMKE2F0dGVzdGF0aW9uZG9jdW1lbnQY4dPsZyABKAxSE2F0dGVzdG'
    'F0aW9uZG9jdW1lbnQSVwoWa2V5ZW5jcnlwdGlvbmFsZ29yaXRobRizkYXkASABKA4yGy5rbXMu'
    'S2V5RW5jcnlwdGlvbk1lY2hhbmlzbVIWa2V5ZW5jcnlwdGlvbmFsZ29yaXRobQ==');

@$core.Deprecated('Use replicateKeyRequestDescriptor instead')
const ReplicateKeyRequest$json = {
  '1': 'ReplicateKeyRequest',
  '2': [
    {
      '1': 'bypasspolicylockoutsafetycheck',
      '3': 177450851,
      '4': 1,
      '5': 8,
      '10': 'bypasspolicylockoutsafetycheck'
    },
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'policy', '3': 471611296, '4': 1, '5': 9, '10': 'policy'},
    {
      '1': 'replicaregion',
      '3': 160061584,
      '4': 1,
      '5': 9,
      '10': 'replicaregion'
    },
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kms.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `ReplicateKeyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicateKeyRequestDescriptor = $convert.base64Decode(
    'ChNSZXBsaWNhdGVLZXlSZXF1ZXN0EkkKHmJ5cGFzc3BvbGljeWxvY2tvdXRzYWZldHljaGVjax'
    'jj3s5UIAEoCFIeYnlwYXNzcG9saWN5bG9ja291dHNhZmV0eWNoZWNrEiMKC2Rlc2NyaXB0aW9u'
    'GIr0+TYgASgJUgtkZXNjcmlwdGlvbhIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEhoKBnBvbG'
    'ljeRig7/DgASABKAlSBnBvbGljeRInCg1yZXBsaWNhcmVnaW9uGJCxqUwgASgJUg1yZXBsaWNh'
    'cmVnaW9uEiAKBHRhZ3MYwcH2tQEgAygLMggua21zLlRhZ1IEdGFncw==');

@$core.Deprecated('Use replicateKeyResponseDescriptor instead')
const ReplicateKeyResponse$json = {
  '1': 'ReplicateKeyResponse',
  '2': [
    {
      '1': 'replicakeymetadata',
      '3': 68328236,
      '4': 1,
      '5': 11,
      '6': '.kms.KeyMetadata',
      '10': 'replicakeymetadata'
    },
    {
      '1': 'replicapolicy',
      '3': 279018266,
      '4': 1,
      '5': 9,
      '10': 'replicapolicy'
    },
    {
      '1': 'replicatags',
      '3': 17934651,
      '4': 3,
      '5': 11,
      '6': '.kms.Tag',
      '10': 'replicatags'
    },
  ],
};

/// Descriptor for `ReplicateKeyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List replicateKeyResponseDescriptor = $convert.base64Decode(
    'ChRSZXBsaWNhdGVLZXlSZXNwb25zZRJDChJyZXBsaWNha2V5bWV0YWRhdGEYrLbKICABKAsyEC'
    '5rbXMuS2V5TWV0YWRhdGFSEnJlcGxpY2FrZXltZXRhZGF0YRIoCg1yZXBsaWNhcG9saWN5GJr2'
    'hYUBIAEoCVINcmVwbGljYXBvbGljeRItCgtyZXBsaWNhdGFncxi70sYIIAMoCzIILmttcy5UYW'
    'dSC3JlcGxpY2F0YWdz');

@$core.Deprecated('Use retireGrantRequestDescriptor instead')
const RetireGrantRequest$json = {
  '1': 'RetireGrantRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'grantid', '3': 66852281, '4': 1, '5': 9, '10': 'grantid'},
    {'1': 'granttoken', '3': 137683547, '4': 1, '5': 9, '10': 'granttoken'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `RetireGrantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List retireGrantRequestDescriptor = $convert.base64Decode(
    'ChJSZXRpcmVHcmFudFJlcXVlc3QSGQoGZHJ5cnVuGLjf+ysgASgIUgZkcnlydW4SGwoHZ3Jhbn'
    'RpZBi5q/AfIAEoCVIHZ3JhbnRpZBIhCgpncmFudHRva2VuGNvE00EgASgJUgpncmFudHRva2Vu'
    'EhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQ=');

@$core.Deprecated('Use revokeGrantRequestDescriptor instead')
const RevokeGrantRequest$json = {
  '1': 'RevokeGrantRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'grantid', '3': 66852281, '4': 1, '5': 9, '10': 'grantid'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `RevokeGrantRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List revokeGrantRequestDescriptor = $convert.base64Decode(
    'ChJSZXZva2VHcmFudFJlcXVlc3QSGQoGZHJ5cnVuGLjf+ysgASgIUgZkcnlydW4SGwoHZ3Jhbn'
    'RpZBi5q/AfIAEoCVIHZ3JhbnRpZBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlk');

@$core.Deprecated('Use rotateKeyOnDemandRequestDescriptor instead')
const RotateKeyOnDemandRequest$json = {
  '1': 'RotateKeyOnDemandRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `RotateKeyOnDemandRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rotateKeyOnDemandRequestDescriptor =
    $convert.base64Decode(
        'ChhSb3RhdGVLZXlPbkRlbWFuZFJlcXVlc3QSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZA==');

@$core.Deprecated('Use rotateKeyOnDemandResponseDescriptor instead')
const RotateKeyOnDemandResponse$json = {
  '1': 'RotateKeyOnDemandResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `RotateKeyOnDemandResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rotateKeyOnDemandResponseDescriptor =
    $convert.base64Decode(
        'ChlSb3RhdGVLZXlPbkRlbWFuZFJlc3BvbnNlEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQ=');

@$core.Deprecated('Use rotationsListEntryDescriptor instead')
const RotationsListEntry$json = {
  '1': 'RotationsListEntry',
  '2': [
    {
      '1': 'expirationmodel',
      '3': 113933558,
      '4': 1,
      '5': 14,
      '6': '.kms.ExpirationModelType',
      '10': 'expirationmodel'
    },
    {
      '1': 'importstate',
      '3': 32548970,
      '4': 1,
      '5': 14,
      '6': '.kms.ImportState',
      '10': 'importstate'
    },
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keymaterialdescription',
      '3': 277153546,
      '4': 1,
      '5': 9,
      '10': 'keymaterialdescription'
    },
    {
      '1': 'keymaterialid',
      '3': 147011585,
      '4': 1,
      '5': 9,
      '10': 'keymaterialid'
    },
    {
      '1': 'keymaterialstate',
      '3': 431806871,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyMaterialState',
      '10': 'keymaterialstate'
    },
    {'1': 'rotationdate', '3': 529238652, '4': 1, '5': 9, '10': 'rotationdate'},
    {
      '1': 'rotationtype',
      '3': 122951592,
      '4': 1,
      '5': 14,
      '6': '.kms.RotationType',
      '10': 'rotationtype'
    },
    {'1': 'validto', '3': 522718673, '4': 1, '5': 9, '10': 'validto'},
  ],
};

/// Descriptor for `RotationsListEntry`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List rotationsListEntryDescriptor = $convert.base64Decode(
    'ChJSb3RhdGlvbnNMaXN0RW50cnkSRQoPZXhwaXJhdGlvbm1vZGVsGPb5qTYgASgOMhgua21zLk'
    'V4cGlyYXRpb25Nb2RlbFR5cGVSD2V4cGlyYXRpb25tb2RlbBI1CgtpbXBvcnRzdGF0ZRjq0MIP'
    'IAEoDjIQLmttcy5JbXBvcnRTdGF0ZVILaW1wb3J0c3RhdGUSGAoFa2V5aWQYooDIgwEgASgJUg'
    'VrZXlpZBI6ChZrZXltYXRlcmlhbGRlc2NyaXB0aW9uGIqOlIQBIAEoCVIWa2V5bWF0ZXJpYWxk'
    'ZXNjcmlwdGlvbhInCg1rZXltYXRlcmlhbGlkGIHwjEYgASgJUg1rZXltYXRlcmlhbGlkEkUKEG'
    'tleW1hdGVyaWFsc3RhdGUYl7PzzQEgASgOMhUua21zLktleU1hdGVyaWFsU3RhdGVSEGtleW1h'
    'dGVyaWFsc3RhdGUSJgoMcm90YXRpb25kYXRlGPyUrvwBIAEoCVIMcm90YXRpb25kYXRlEjgKDH'
    'JvdGF0aW9udHlwZRior9A6IAEoDjIRLmttcy5Sb3RhdGlvblR5cGVSDHJvdGF0aW9udHlwZRIc'
    'Cgd2YWxpZHRvGNGboPkBIAEoCVIHdmFsaWR0bw==');

@$core.Deprecated('Use scheduleKeyDeletionRequestDescriptor instead')
const ScheduleKeyDeletionRequest$json = {
  '1': 'ScheduleKeyDeletionRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'pendingwindowindays',
      '3': 532945081,
      '4': 1,
      '5': 5,
      '10': 'pendingwindowindays'
    },
  ],
};

/// Descriptor for `ScheduleKeyDeletionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduleKeyDeletionRequestDescriptor =
    $convert.base64Decode(
        'ChpTY2hlZHVsZUtleURlbGV0aW9uUmVxdWVzdBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEj'
        'QKE3BlbmRpbmd3aW5kb3dpbmRheXMYubGQ/gEgASgFUhNwZW5kaW5nd2luZG93aW5kYXlz');

@$core.Deprecated('Use scheduleKeyDeletionResponseDescriptor instead')
const ScheduleKeyDeletionResponse$json = {
  '1': 'ScheduleKeyDeletionResponse',
  '2': [
    {'1': 'deletiondate', '3': 347845564, '4': 1, '5': 9, '10': 'deletiondate'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'keystate',
      '3': 27894226,
      '4': 1,
      '5': 14,
      '6': '.kms.KeyState',
      '10': 'keystate'
    },
    {
      '1': 'pendingwindowindays',
      '3': 532945081,
      '4': 1,
      '5': 5,
      '10': 'pendingwindowindays'
    },
  ],
};

/// Descriptor for `ScheduleKeyDeletionResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List scheduleKeyDeletionResponseDescriptor = $convert.base64Decode(
    'ChtTY2hlZHVsZUtleURlbGV0aW9uUmVzcG9uc2USJgoMZGVsZXRpb25kYXRlGLzn7qUBIAEoCV'
    'IMZGVsZXRpb25kYXRlEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSLAoIa2V5c3RhdGUY0sOm'
    'DSABKA4yDS5rbXMuS2V5U3RhdGVSCGtleXN0YXRlEjQKE3BlbmRpbmd3aW5kb3dpbmRheXMYub'
    'GQ/gEgASgFUhNwZW5kaW5nd2luZG93aW5kYXlz');

@$core.Deprecated('Use signRequestDescriptor instead')
const SignRequest$json = {
  '1': 'SignRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 12, '10': 'message'},
    {
      '1': 'messagetype',
      '3': 452676013,
      '4': 1,
      '5': 14,
      '6': '.kms.MessageType',
      '10': 'messagetype'
    },
    {
      '1': 'signingalgorithm',
      '3': 488091842,
      '4': 1,
      '5': 14,
      '6': '.kms.SigningAlgorithmSpec',
      '10': 'signingalgorithm'
    },
  ],
};

/// Descriptor for `SignRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signRequestDescriptor = $convert.base64Decode(
    'CgtTaWduUmVxdWVzdBIZCgZkcnlydW4YuN/7KyABKAhSBmRyeXJ1bhIkCgtncmFudHRva2Vucx'
    'iMjYCiASADKAlSC2dyYW50dG9rZW5zEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQSGwoHbWVz'
    'c2FnZRiFs7twIAEoDFIHbWVzc2FnZRI2CgttZXNzYWdldHlwZRitk+3XASABKA4yEC5rbXMuTW'
    'Vzc2FnZVR5cGVSC21lc3NhZ2V0eXBlEkkKEHNpZ25pbmdhbGdvcml0aG0YwuHe6AEgASgOMhku'
    'a21zLlNpZ25pbmdBbGdvcml0aG1TcGVjUhBzaWduaW5nYWxnb3JpdGht');

@$core.Deprecated('Use signResponseDescriptor instead')
const SignResponse$json = {
  '1': 'SignResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'signature', '3': 4785422, '4': 1, '5': 12, '10': 'signature'},
    {
      '1': 'signingalgorithm',
      '3': 488091842,
      '4': 1,
      '5': 14,
      '6': '.kms.SigningAlgorithmSpec',
      '10': 'signingalgorithm'
    },
  ],
};

/// Descriptor for `SignResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List signResponseDescriptor = $convert.base64Decode(
    'CgxTaWduUmVzcG9uc2USGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBIfCglzaWduYXR1cmUYjo'
    'qkAiABKAxSCXNpZ25hdHVyZRJJChBzaWduaW5nYWxnb3JpdGhtGMLh3ugBIAEoDjIZLmttcy5T'
    'aWduaW5nQWxnb3JpdGhtU3BlY1IQc2lnbmluZ2FsZ29yaXRobQ==');

@$core.Deprecated('Use tagDescriptor instead')
const Tag$json = {
  '1': 'Tag',
  '2': [
    {'1': 'tagkey', '3': 175603339, '4': 1, '5': 9, '10': 'tagkey'},
    {'1': 'tagvalue', '3': 487739505, '4': 1, '5': 9, '10': 'tagvalue'},
  ],
};

/// Descriptor for `Tag`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagDescriptor = $convert.base64Decode(
    'CgNUYWcSGQoGdGFna2V5GIv93VMgASgJUgZ0YWdrZXkSHgoIdGFndmFsdWUY8aDJ6AEgASgJUg'
    'h0YWd2YWx1ZQ==');

@$core.Deprecated('Use tagExceptionDescriptor instead')
const TagException$json = {
  '1': 'TagException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `TagException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagExceptionDescriptor = $convert.base64Decode(
    'CgxUYWdFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2FnZQ==');

@$core.Deprecated('Use tagResourceRequestDescriptor instead')
const TagResourceRequest$json = {
  '1': 'TagResourceRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'tags',
      '3': 381526209,
      '4': 3,
      '5': 11,
      '6': '.kms.Tag',
      '10': 'tags'
    },
  ],
};

/// Descriptor for `TagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List tagResourceRequestDescriptor = $convert.base64Decode(
    'ChJUYWdSZXNvdXJjZVJlcXVlc3QSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBIgCgR0YWdzGM'
    'HB9rUBIAMoCzIILmttcy5UYWdSBHRhZ3M=');

@$core.Deprecated('Use unsupportedOperationExceptionDescriptor instead')
const UnsupportedOperationException$json = {
  '1': 'UnsupportedOperationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `UnsupportedOperationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List unsupportedOperationExceptionDescriptor =
    $convert.base64Decode(
        'Ch1VbnN1cHBvcnRlZE9wZXJhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZX'
        'NzYWdl');

@$core.Deprecated('Use untagResourceRequestDescriptor instead')
const UntagResourceRequest$json = {
  '1': 'UntagResourceRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'tagkeys', '3': 320659964, '4': 3, '5': 9, '10': 'tagkeys'},
  ],
};

/// Descriptor for `UntagResourceRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List untagResourceRequestDescriptor = $convert.base64Decode(
    'ChRVbnRhZ1Jlc291cmNlUmVxdWVzdBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEhwKB3RhZ2'
    'tleXMY/MPzmAEgAygJUgd0YWdrZXlz');

@$core.Deprecated('Use updateAliasRequestDescriptor instead')
const UpdateAliasRequest$json = {
  '1': 'UpdateAliasRequest',
  '2': [
    {'1': 'aliasname', '3': 313250709, '4': 1, '5': 9, '10': 'aliasname'},
    {'1': 'targetkeyid', '3': 406196123, '4': 1, '5': 9, '10': 'targetkeyid'},
  ],
};

/// Descriptor for `UpdateAliasRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateAliasRequestDescriptor = $convert.base64Decode(
    'ChJVcGRhdGVBbGlhc1JlcXVlc3QSIAoJYWxpYXNuYW1lGJWnr5UBIAEoCVIJYWxpYXNuYW1lEi'
    'QKC3RhcmdldGtleWlkGJuf2MEBIAEoCVILdGFyZ2V0a2V5aWQ=');

@$core.Deprecated('Use updateCustomKeyStoreRequestDescriptor instead')
const UpdateCustomKeyStoreRequest$json = {
  '1': 'UpdateCustomKeyStoreRequest',
  '2': [
    {
      '1': 'cloudhsmclusterid',
      '3': 56498754,
      '4': 1,
      '5': 9,
      '10': 'cloudhsmclusterid'
    },
    {
      '1': 'customkeystoreid',
      '3': 88348228,
      '4': 1,
      '5': 9,
      '10': 'customkeystoreid'
    },
    {
      '1': 'keystorepassword',
      '3': 403136353,
      '4': 1,
      '5': 9,
      '10': 'keystorepassword'
    },
    {
      '1': 'newcustomkeystorename',
      '3': 10936866,
      '4': 1,
      '5': 9,
      '10': 'newcustomkeystorename'
    },
    {
      '1': 'xksproxyauthenticationcredential',
      '3': 350418199,
      '4': 1,
      '5': 11,
      '6': '.kms.XksProxyAuthenticationCredentialType',
      '10': 'xksproxyauthenticationcredential'
    },
    {
      '1': 'xksproxyconnectivity',
      '3': 298569161,
      '4': 1,
      '5': 14,
      '6': '.kms.XksProxyConnectivityType',
      '10': 'xksproxyconnectivity'
    },
    {
      '1': 'xksproxyuriendpoint',
      '3': 273255559,
      '4': 1,
      '5': 9,
      '10': 'xksproxyuriendpoint'
    },
    {
      '1': 'xksproxyuripath',
      '3': 436753509,
      '4': 1,
      '5': 9,
      '10': 'xksproxyuripath'
    },
    {
      '1': 'xksproxyvpcendpointservicename',
      '3': 372786130,
      '4': 1,
      '5': 9,
      '10': 'xksproxyvpcendpointservicename'
    },
    {
      '1': 'xksproxyvpcendpointserviceowner',
      '3': 55249590,
      '4': 1,
      '5': 9,
      '10': 'xksproxyvpcendpointserviceowner'
    },
  ],
};

/// Descriptor for `UpdateCustomKeyStoreRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCustomKeyStoreRequestDescriptor = $convert.base64Decode(
    'ChtVcGRhdGVDdXN0b21LZXlTdG9yZVJlcXVlc3QSLwoRY2xvdWRoc21jbHVzdGVyaWQYwrT4Gi'
    'ABKAlSEWNsb3VkaHNtY2x1c3RlcmlkEi0KEGN1c3RvbWtleXN0b3JlaWQYxKyQKiABKAlSEGN1'
    'c3RvbWtleXN0b3JlaWQSLgoQa2V5c3RvcmVwYXNzd29yZBjhvp3AASABKAlSEGtleXN0b3JlcG'
    'Fzc3dvcmQSNwoVbmV3Y3VzdG9ta2V5c3RvcmVuYW1lGKLEmwUgASgJUhVuZXdjdXN0b21rZXlz'
    'dG9yZW5hbWUSeQogeGtzcHJveHlhdXRoZW50aWNhdGlvbmNyZWRlbnRpYWwYl+qLpwEgASgLMi'
    'kua21zLlhrc1Byb3h5QXV0aGVudGljYXRpb25DcmVkZW50aWFsVHlwZVIgeGtzcHJveHlhdXRo'
    'ZW50aWNhdGlvbmNyZWRlbnRpYWwSVQoUeGtzcHJveHljb25uZWN0aXZpdHkYyZuvjgEgASgOMh'
    '0ua21zLlhrc1Byb3h5Q29ubmVjdGl2aXR5VHlwZVIUeGtzcHJveHljb25uZWN0aXZpdHkSNAoT'
    'eGtzcHJveHl1cmllbmRwb2ludBiHmaaCASABKAlSE3hrc3Byb3h5dXJpZW5kcG9pbnQSLAoPeG'
    'tzcHJveHl1cmlwYXRoGOWoodABIAEoCVIPeGtzcHJveHl1cmlwYXRoEkoKHnhrc3Byb3h5dnBj'
    'ZW5kcG9pbnRzZXJ2aWNlbmFtZRjSh+GxASABKAlSHnhrc3Byb3h5dnBjZW5kcG9pbnRzZXJ2aW'
    'NlbmFtZRJLCh94a3Nwcm94eXZwY2VuZHBvaW50c2VydmljZW93bmVyGLaVrBogASgJUh94a3Nw'
    'cm94eXZwY2VuZHBvaW50c2VydmljZW93bmVy');

@$core.Deprecated('Use updateCustomKeyStoreResponseDescriptor instead')
const UpdateCustomKeyStoreResponse$json = {
  '1': 'UpdateCustomKeyStoreResponse',
};

/// Descriptor for `UpdateCustomKeyStoreResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateCustomKeyStoreResponseDescriptor =
    $convert.base64Decode('ChxVcGRhdGVDdXN0b21LZXlTdG9yZVJlc3BvbnNl');

@$core.Deprecated('Use updateKeyDescriptionRequestDescriptor instead')
const UpdateKeyDescriptionRequest$json = {
  '1': 'UpdateKeyDescriptionRequest',
  '2': [
    {'1': 'description', '3': 115243530, '4': 1, '5': 9, '10': 'description'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
  ],
};

/// Descriptor for `UpdateKeyDescriptionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updateKeyDescriptionRequestDescriptor =
    $convert.base64Decode(
        'ChtVcGRhdGVLZXlEZXNjcmlwdGlvblJlcXVlc3QSIwoLZGVzY3JpcHRpb24YivT5NiABKAlSC2'
        'Rlc2NyaXB0aW9uEhgKBWtleWlkGKKAyIMBIAEoCVIFa2V5aWQ=');

@$core.Deprecated('Use updatePrimaryRegionRequestDescriptor instead')
const UpdatePrimaryRegionRequest$json = {
  '1': 'UpdatePrimaryRegionRequest',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'primaryregion',
      '3': 480901186,
      '4': 1,
      '5': 9,
      '10': 'primaryregion'
    },
  ],
};

/// Descriptor for `UpdatePrimaryRegionRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List updatePrimaryRegionRequestDescriptor =
    $convert.base64Decode(
        'ChpVcGRhdGVQcmltYXJ5UmVnaW9uUmVxdWVzdBIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEi'
        'gKDXByaW1hcnlyZWdpb24YwvCn5QEgASgJUg1wcmltYXJ5cmVnaW9u');

@$core.Deprecated('Use verifyMacRequestDescriptor instead')
const VerifyMacRequest$json = {
  '1': 'VerifyMacRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'mac', '3': 296223945, '4': 1, '5': 12, '10': 'mac'},
    {
      '1': 'macalgorithm',
      '3': 253519878,
      '4': 1,
      '5': 14,
      '6': '.kms.MacAlgorithmSpec',
      '10': 'macalgorithm'
    },
    {'1': 'message', '3': 235854213, '4': 1, '5': 12, '10': 'message'},
  ],
};

/// Descriptor for `VerifyMacRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyMacRequestDescriptor = $convert.base64Decode(
    'ChBWZXJpZnlNYWNSZXF1ZXN0EhkKBmRyeXJ1bhi43/srIAEoCFIGZHJ5cnVuEiQKC2dyYW50dG'
    '9rZW5zGIyNgKIBIAMoCVILZ3JhbnR0b2tlbnMSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBIU'
    'CgNtYWMYyYmgjQEgASgMUgNtYWMSPAoMbWFjYWxnb3JpdGhtGIbQ8XggASgOMhUua21zLk1hY0'
    'FsZ29yaXRobVNwZWNSDG1hY2FsZ29yaXRobRIbCgdtZXNzYWdlGIWzu3AgASgMUgdtZXNzYWdl');

@$core.Deprecated('Use verifyMacResponseDescriptor instead')
const VerifyMacResponse$json = {
  '1': 'VerifyMacResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'macalgorithm',
      '3': 253519878,
      '4': 1,
      '5': 14,
      '6': '.kms.MacAlgorithmSpec',
      '10': 'macalgorithm'
    },
    {'1': 'macvalid', '3': 482746075, '4': 1, '5': 8, '10': 'macvalid'},
  ],
};

/// Descriptor for `VerifyMacResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyMacResponseDescriptor = $convert.base64Decode(
    'ChFWZXJpZnlNYWNSZXNwb25zZRIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEjwKDG1hY2FsZ2'
    '9yaXRobRiG0PF4IAEoDjIVLmttcy5NYWNBbGdvcml0aG1TcGVjUgxtYWNhbGdvcml0aG0SHgoI'
    'bWFjdmFsaWQY272Y5gEgASgIUghtYWN2YWxpZA==');

@$core.Deprecated('Use verifyRequestDescriptor instead')
const VerifyRequest$json = {
  '1': 'VerifyRequest',
  '2': [
    {'1': 'dryrun', '3': 92204984, '4': 1, '5': 8, '10': 'dryrun'},
    {'1': 'granttokens', '3': 339740300, '4': 3, '5': 9, '10': 'granttokens'},
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {'1': 'message', '3': 235854213, '4': 1, '5': 12, '10': 'message'},
    {
      '1': 'messagetype',
      '3': 452676013,
      '4': 1,
      '5': 14,
      '6': '.kms.MessageType',
      '10': 'messagetype'
    },
    {'1': 'signature', '3': 4785422, '4': 1, '5': 12, '10': 'signature'},
    {
      '1': 'signingalgorithm',
      '3': 488091842,
      '4': 1,
      '5': 14,
      '6': '.kms.SigningAlgorithmSpec',
      '10': 'signingalgorithm'
    },
  ],
};

/// Descriptor for `VerifyRequest`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyRequestDescriptor = $convert.base64Decode(
    'Cg1WZXJpZnlSZXF1ZXN0EhkKBmRyeXJ1bhi43/srIAEoCFIGZHJ5cnVuEiQKC2dyYW50dG9rZW'
    '5zGIyNgKIBIAMoCVILZ3JhbnR0b2tlbnMSGAoFa2V5aWQYooDIgwEgASgJUgVrZXlpZBIbCgdt'
    'ZXNzYWdlGIWzu3AgASgMUgdtZXNzYWdlEjYKC21lc3NhZ2V0eXBlGK2T7dcBIAEoDjIQLmttcy'
    '5NZXNzYWdlVHlwZVILbWVzc2FnZXR5cGUSHwoJc2lnbmF0dXJlGI6KpAIgASgMUglzaWduYXR1'
    'cmUSSQoQc2lnbmluZ2FsZ29yaXRobRjC4d7oASABKA4yGS5rbXMuU2lnbmluZ0FsZ29yaXRobV'
    'NwZWNSEHNpZ25pbmdhbGdvcml0aG0=');

@$core.Deprecated('Use verifyResponseDescriptor instead')
const VerifyResponse$json = {
  '1': 'VerifyResponse',
  '2': [
    {'1': 'keyid', '3': 275906594, '4': 1, '5': 9, '10': 'keyid'},
    {
      '1': 'signaturevalid',
      '3': 272180330,
      '4': 1,
      '5': 8,
      '10': 'signaturevalid'
    },
    {
      '1': 'signingalgorithm',
      '3': 488091842,
      '4': 1,
      '5': 14,
      '6': '.kms.SigningAlgorithmSpec',
      '10': 'signingalgorithm'
    },
  ],
};

/// Descriptor for `VerifyResponse`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List verifyResponseDescriptor = $convert.base64Decode(
    'Cg5WZXJpZnlSZXNwb25zZRIYCgVrZXlpZBiigMiDASABKAlSBWtleWlkEioKDnNpZ25hdHVyZX'
    'ZhbGlkGOrI5IEBIAEoCFIOc2lnbmF0dXJldmFsaWQSSQoQc2lnbmluZ2FsZ29yaXRobRjC4d7o'
    'ASABKA4yGS5rbXMuU2lnbmluZ0FsZ29yaXRobVNwZWNSEHNpZ25pbmdhbGdvcml0aG0=');

@$core.Deprecated('Use xksKeyAlreadyInUseExceptionDescriptor instead')
const XksKeyAlreadyInUseException$json = {
  '1': 'XksKeyAlreadyInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksKeyAlreadyInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksKeyAlreadyInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChtYa3NLZXlBbHJlYWR5SW5Vc2VFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCVIHbWVzc2'
        'FnZQ==');

@$core.Deprecated('Use xksKeyConfigurationTypeDescriptor instead')
const XksKeyConfigurationType$json = {
  '1': 'XksKeyConfigurationType',
  '2': [
    {'1': 'id', '3': 384350465, '4': 1, '5': 9, '10': 'id'},
  ],
};

/// Descriptor for `XksKeyConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksKeyConfigurationTypeDescriptor =
    $convert.base64Decode(
        'ChdYa3NLZXlDb25maWd1cmF0aW9uVHlwZRISCgJpZBiB8qK3ASABKAlSAmlk');

@$core.Deprecated('Use xksKeyInvalidConfigurationExceptionDescriptor instead')
const XksKeyInvalidConfigurationException$json = {
  '1': 'XksKeyInvalidConfigurationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksKeyInvalidConfigurationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksKeyInvalidConfigurationExceptionDescriptor =
    $convert.base64Decode(
        'CiNYa3NLZXlJbnZhbGlkQ29uZmlndXJhdGlvbkV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgAS'
        'gJUgdtZXNzYWdl');

@$core.Deprecated('Use xksKeyNotFoundExceptionDescriptor instead')
const XksKeyNotFoundException$json = {
  '1': 'XksKeyNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksKeyNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksKeyNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'ChdYa3NLZXlOb3RGb3VuZEV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use xksProxyAuthenticationCredentialTypeDescriptor instead')
const XksProxyAuthenticationCredentialType$json = {
  '1': 'XksProxyAuthenticationCredentialType',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {
      '1': 'rawsecretaccesskey',
      '3': 68137571,
      '4': 1,
      '5': 9,
      '10': 'rawsecretaccesskey'
    },
  ],
};

/// Descriptor for `XksProxyAuthenticationCredentialType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyAuthenticationCredentialTypeDescriptor =
    $convert.base64Decode(
        'CiRYa3NQcm94eUF1dGhlbnRpY2F0aW9uQ3JlZGVudGlhbFR5cGUSJAoLYWNjZXNza2V5aWQYoL'
        'e32AEgASgJUgthY2Nlc3NrZXlpZBIxChJyYXdzZWNyZXRhY2Nlc3NrZXkY4+S+ICABKAlSEnJh'
        'd3NlY3JldGFjY2Vzc2tleQ==');

@$core.Deprecated('Use xksProxyConfigurationTypeDescriptor instead')
const XksProxyConfigurationType$json = {
  '1': 'XksProxyConfigurationType',
  '2': [
    {'1': 'accesskeyid', '3': 453893024, '4': 1, '5': 9, '10': 'accesskeyid'},
    {
      '1': 'connectivity',
      '3': 210638147,
      '4': 1,
      '5': 14,
      '6': '.kms.XksProxyConnectivityType',
      '10': 'connectivity'
    },
    {'1': 'uriendpoint', '3': 79142005, '4': 1, '5': 9, '10': 'uriendpoint'},
    {'1': 'uripath', '3': 288340351, '4': 1, '5': 9, '10': 'uripath'},
    {
      '1': 'vpcendpointservicename',
      '3': 269882444,
      '4': 1,
      '5': 9,
      '10': 'vpcendpointservicename'
    },
    {
      '1': 'vpcendpointserviceowner',
      '3': 298819456,
      '4': 1,
      '5': 9,
      '10': 'vpcendpointserviceowner'
    },
  ],
};

/// Descriptor for `XksProxyConfigurationType`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyConfigurationTypeDescriptor = $convert.base64Decode(
    'ChlYa3NQcm94eUNvbmZpZ3VyYXRpb25UeXBlEiQKC2FjY2Vzc2tleWlkGKC3t9gBIAEoCVILYW'
    'NjZXNza2V5aWQSRAoMY29ubmVjdGl2aXR5GMOquGQgASgOMh0ua21zLlhrc1Byb3h5Q29ubmVj'
    'dGl2aXR5VHlwZVIMY29ubmVjdGl2aXR5EiMKC3VyaWVuZHBvaW50GPW43iUgASgJUgt1cmllbm'
    'Rwb2ludBIcCgd1cmlwYXRoGP/yvokBIAEoCVIHdXJpcGF0aBI6ChZ2cGNlbmRwb2ludHNlcnZp'
    'Y2VuYW1lGMyo2IABIAEoCVIWdnBjZW5kcG9pbnRzZXJ2aWNlbmFtZRI8Chd2cGNlbmRwb2ludH'
    'NlcnZpY2Vvd25lchiAv76OASABKAlSF3ZwY2VuZHBvaW50c2VydmljZW93bmVy');

@$core.Deprecated(
    'Use xksProxyIncorrectAuthenticationCredentialExceptionDescriptor instead')
const XksProxyIncorrectAuthenticationCredentialException$json = {
  '1': 'XksProxyIncorrectAuthenticationCredentialException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyIncorrectAuthenticationCredentialException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    xksProxyIncorrectAuthenticationCredentialExceptionDescriptor =
    $convert.base64Decode(
        'CjJYa3NQcm94eUluY29ycmVjdEF1dGhlbnRpY2F0aW9uQ3JlZGVudGlhbEV4Y2VwdGlvbhIbCg'
        'dtZXNzYWdlGOWRyCcgASgJUgdtZXNzYWdl');

@$core.Deprecated('Use xksProxyInvalidConfigurationExceptionDescriptor instead')
const XksProxyInvalidConfigurationException$json = {
  '1': 'XksProxyInvalidConfigurationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyInvalidConfigurationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyInvalidConfigurationExceptionDescriptor =
    $convert.base64Decode(
        'CiVYa3NQcm94eUludmFsaWRDb25maWd1cmF0aW9uRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJy'
        'ABKAlSB21lc3NhZ2U=');

@$core.Deprecated('Use xksProxyInvalidResponseExceptionDescriptor instead')
const XksProxyInvalidResponseException$json = {
  '1': 'XksProxyInvalidResponseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyInvalidResponseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyInvalidResponseExceptionDescriptor =
    $convert.base64Decode(
        'CiBYa3NQcm94eUludmFsaWRSZXNwb25zZUV4Y2VwdGlvbhIbCgdtZXNzYWdlGOWRyCcgASgJUg'
        'dtZXNzYWdl');

@$core.Deprecated('Use xksProxyUriEndpointInUseExceptionDescriptor instead')
const XksProxyUriEndpointInUseException$json = {
  '1': 'XksProxyUriEndpointInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyUriEndpointInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyUriEndpointInUseExceptionDescriptor =
    $convert.base64Decode(
        'CiFYa3NQcm94eVVyaUVuZHBvaW50SW5Vc2VFeGNlcHRpb24SGwoHbWVzc2FnZRjlkcgnIAEoCV'
        'IHbWVzc2FnZQ==');

@$core.Deprecated('Use xksProxyUriInUseExceptionDescriptor instead')
const XksProxyUriInUseException$json = {
  '1': 'XksProxyUriInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyUriInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyUriInUseExceptionDescriptor =
    $convert.base64Decode(
        'ChlYa3NQcm94eVVyaUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2'
        'U=');

@$core.Deprecated('Use xksProxyUriUnreachableExceptionDescriptor instead')
const XksProxyUriUnreachableException$json = {
  '1': 'XksProxyUriUnreachableException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyUriUnreachableException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyUriUnreachableExceptionDescriptor =
    $convert.base64Decode(
        'Ch9Ya3NQcm94eVVyaVVucmVhY2hhYmxlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB2'
        '1lc3NhZ2U=');

@$core.Deprecated(
    'Use xksProxyVpcEndpointServiceInUseExceptionDescriptor instead')
const XksProxyVpcEndpointServiceInUseException$json = {
  '1': 'XksProxyVpcEndpointServiceInUseException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyVpcEndpointServiceInUseException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List xksProxyVpcEndpointServiceInUseExceptionDescriptor =
    $convert.base64Decode(
        'CihYa3NQcm94eVZwY0VuZHBvaW50U2VydmljZUluVXNlRXhjZXB0aW9uEhsKB21lc3NhZ2UY5Z'
        'HIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use xksProxyVpcEndpointServiceInvalidConfigurationExceptionDescriptor instead')
const XksProxyVpcEndpointServiceInvalidConfigurationException$json = {
  '1': 'XksProxyVpcEndpointServiceInvalidConfigurationException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyVpcEndpointServiceInvalidConfigurationException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    xksProxyVpcEndpointServiceInvalidConfigurationExceptionDescriptor =
    $convert.base64Decode(
        'CjdYa3NQcm94eVZwY0VuZHBvaW50U2VydmljZUludmFsaWRDb25maWd1cmF0aW9uRXhjZXB0aW'
        '9uEhsKB21lc3NhZ2UY5ZHIJyABKAlSB21lc3NhZ2U=');

@$core.Deprecated(
    'Use xksProxyVpcEndpointServiceNotFoundExceptionDescriptor instead')
const XksProxyVpcEndpointServiceNotFoundException$json = {
  '1': 'XksProxyVpcEndpointServiceNotFoundException',
  '2': [
    {'1': 'message', '3': 82970853, '4': 1, '5': 9, '10': 'message'},
  ],
};

/// Descriptor for `XksProxyVpcEndpointServiceNotFoundException`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List
    xksProxyVpcEndpointServiceNotFoundExceptionDescriptor =
    $convert.base64Decode(
        'CitYa3NQcm94eVZwY0VuZHBvaW50U2VydmljZU5vdEZvdW5kRXhjZXB0aW9uEhsKB21lc3NhZ2'
        'UY5ZHIJyABKAlSB21lc3NhZ2U=');
