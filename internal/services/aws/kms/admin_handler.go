package kms

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/kms"
	"vorpalstacks/internal/pb/aws/kms/kmsconnect"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// AdminHandler implements the KMS admin console gRPC-Web handler.
// It delegates to the shared KMSService *Core methods so that the same
// validation, key creation, and policy handling code path is used by
// both the HTTP API handlers and the admin console handlers.
type AdminHandler struct {
	kmsconnect.UnimplementedKMSServiceHandler
	service *KMSService
}

var _ kmsconnect.KMSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new KMS admin handler with the given service.
func NewAdminHandler(svc *KMSService) *AdminHandler {
	return &AdminHandler{service: svc}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*kmsStores, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListKeys retrieves all KMS keys from the key store with pagination.
func (h *AdminHandler) ListKeys(ctx context.Context, req *connect.Request[pb.ListKeysRequest]) (*connect.Response[pb.ListKeysResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	limit := int(req.Msg.GetLimit())
	if limit <= 0 {
		limit = 100
	}
	// Smithy LimitType: @range(1-1000).
	if limit > 1000 {
		limit = 1000
	}
	result, err := stores.keys.List(req.Msg.Marker, limit)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	keys := make([]*pb.KeyListEntry, len(result.Keys))
	for i, k := range result.Keys {
		keys[i] = &pb.KeyListEntry{
			Keyid:  k.KeyID,
			Keyarn: k.KeyArn,
		}
	}

	return connect.NewResponse(&pb.ListKeysResponse{
		Keys:       keys,
		Nextmarker: result.NextMarker,
		Truncated:  proto.Bool(result.IsTruncated),
	}), nil
}

// CreateKey creates a new KMS key via the admin console. It converts the
// proto request into a CreateKeyInput DTO and delegates all validation
// and persistence to KMSService.createKeyCore, ensuring the same code
// path as the HTTP API.
func (h *AdminHandler) CreateKey(ctx context.Context, req *connect.Request[pb.CreateKeyRequest]) (*connect.Response[pb.CreateKeyResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Apply proto-level defaults so the response carries the resolved
	// enum values, not the zero-value "unspecified".
	keyUsageProto := req.Msg.GetKeyusage()
	if keyUsageProto == 0 {
		keyUsageProto = pb.KeyUsageType_KEY_USAGE_TYPE_ENCRYPT_DECRYPT
	}
	keySpecProto := req.Msg.GetKeyspec()
	if keySpecProto == 0 {
		keySpecProto = pb.KeySpec_KEY_SPEC_SYMMETRIC_DEFAULT
	}
	originProto := req.Msg.GetOrigin()
	if originProto == 0 {
		originProto = pb.OriginType_ORIGIN_TYPE_AWS_KMS
	}

	// Convert proto enums to store types.
	keyUsage := protoKeyUsageToStore(keyUsageProto)
	keySpec := protoKeySpecToStore(keySpecProto)
	origin := protoOriginToStore(originProto)

	// Convert proto tags.
	var tags []types.Tag
	for _, t := range req.Msg.GetTags() {
		tags = append(tags, types.Tag{Key: t.GetTagkey(), Value: t.GetTagvalue()})
	}

	key, err := h.service.createKeyCore(stores, CreateKeyInput{
		Description:        req.Msg.GetDescription(),
		KeyUsage:           keyUsage,
		KeySpec:            keySpec,
		Origin:             origin,
		MultiRegion:        req.Msg.GetMultiregion(),
		CustomKeyStoreID:   req.Msg.GetCustomkeystoreid(),
		XksKeyID:           req.Msg.GetXkskeyid(),
		Policy:             req.Msg.GetPolicy(),
		BypassLockoutCheck: req.Msg.GetBypasspolicylockoutsafetycheck(),
		Tags:               tags,
		AccountID:          h.service.accountID,
	})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	return connect.NewResponse(&pb.CreateKeyResponse{
		Keymetadata: h.buildProtoKeyMetadata(key, keyUsageProto, keySpecProto, originProto),
	}), nil
}

// ScheduleKeyDeletion schedules a KMS key for deletion via the admin
// console. It resolves the key through the service layer (supporting
// alias and ARN forms) and delegates to scheduleKeyDeletionCore.
func (h *AdminHandler) ScheduleKeyDeletion(ctx context.Context, req *connect.Request[pb.ScheduleKeyDeletionRequest]) (*connect.Response[pb.ScheduleKeyDeletionResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.GetKeyid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("KeyId is required"))
	}

	// Resolve the key through the service layer to support alias and ARN.
	key, err := h.service.resolveKey(stores, map[string]interface{}{"KeyId": req.Msg.GetKeyid()})
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	// AWS: PendingWindowInDays defaults to 30 when omitted. Proto3 zero
	// value (0) cannot distinguish unset from explicit 0, so we treat 0
	// as "default to 30". Range validation (7-30) is in the Core method.
	pendingDays := int(req.Msg.GetPendingwindowindays())
	if pendingDays == 0 {
		pendingDays = 30
	}

	updatedKey, days, err := h.service.scheduleKeyDeletionCore(stores, key.KeyID, pendingDays)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	resp := &pb.ScheduleKeyDeletionResponse{
		Keyid:               updatedKey.KeyID,
		Keystate:            storeKeyStateToProto(updatedKey.KeyState),
		Pendingwindowindays: proto.Int32(int32(days)),
	}
	if updatedKey.DeletionDate != nil {
		resp.Deletiondate = updatedKey.DeletionDate.UTC().Format(timeutils.ISO8601UTCFormat)
	}
	return connect.NewResponse(resp), nil
}

// ---------------------------------------------------------------------------
// Proto ↔ store type conversion helpers
// ---------------------------------------------------------------------------

func protoKeyUsageToStore(v pb.KeyUsageType) kmsstore.KeyUsage {
	switch v {
	case pb.KeyUsageType_KEY_USAGE_TYPE_SIGN_VERIFY:
		return kmsstore.KeyUsageSignVerify
	case pb.KeyUsageType_KEY_USAGE_TYPE_GENERATE_VERIFY_MAC:
		return kmsstore.KeyUsageGenerateVerifyMAC
	default:
		return kmsstore.KeyUsageEncryptDecrypt
	}
}

func protoKeySpecToStore(v pb.KeySpec) kmsstore.KeySpec {
	switch v {
	case pb.KeySpec_KEY_SPEC_SYMMETRIC_DEFAULT:
		return kmsstore.KeySpecSymmetricDefault
	case pb.KeySpec_KEY_SPEC_RSA_2048:
		return kmsstore.KeySpecRSA2048
	case pb.KeySpec_KEY_SPEC_RSA_3072:
		return kmsstore.KeySpecRSA3072
	case pb.KeySpec_KEY_SPEC_RSA_4096:
		return kmsstore.KeySpecRSA4096
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P256:
		return kmsstore.KeySpecECCNISTP256
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P384:
		return kmsstore.KeySpecECCNISTP384
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P521:
		return kmsstore.KeySpecECCNISTP521
	case pb.KeySpec_KEY_SPEC_ECC_SECG_P256K1:
		return kmsstore.KeySpecECCSECGP256K1
	case pb.KeySpec_KEY_SPEC_SM2:
		return kmsstore.KeySpecSM2
	case pb.KeySpec_KEY_SPEC_HMAC_224:
		return kmsstore.KeySpecHMAC224
	case pb.KeySpec_KEY_SPEC_HMAC_256:
		return kmsstore.KeySpecHMAC256
	case pb.KeySpec_KEY_SPEC_HMAC_384:
		return kmsstore.KeySpecHMAC384
	case pb.KeySpec_KEY_SPEC_HMAC_512:
		return kmsstore.KeySpecHMAC512
	default:
		return kmsstore.KeySpecSymmetricDefault
	}
}

func protoOriginToStore(v pb.OriginType) kmsstore.OriginType {
	switch v {
	case pb.OriginType_ORIGIN_TYPE_EXTERNAL:
		return kmsstore.OriginTypeExternal
	default:
		return kmsstore.OriginTypeAWSKMS
	}
}

// ---------------------------------------------------------------------------

// buildProtoKeyMetadata constructs a full pb.KeyMetadata from a store Key,
// mapping all fields that the HTTP API's buildKeyMetadata also emits. This
// ensures admin console clients see the same metadata richness as SDK clients.
func (h *AdminHandler) buildProtoKeyMetadata(key *kmsstore.Key, keyUsage pb.KeyUsageType, keySpec pb.KeySpec, origin pb.OriginType) *pb.KeyMetadata {
	_, _, _, accountID, _ := arn.SplitARN(key.Arn)

	md := &pb.KeyMetadata{
		Awsaccountid:          accountID,
		Keyid:                 key.KeyID,
		Arn:                   key.Arn,
		Keystate:              storeKeyStateToProto(key.KeyState),
		Keyusage:              keyUsage,
		Keyspec:               keySpec,
		Customermasterkeyspec: pb.CustomerMasterKeySpec(keySpec),
		Description:           key.Description,
		Enabled:               proto.Bool(key.Enabled),
		Origin:                origin,
		Keymanager:            pb.KeyManagerType_KEY_MANAGER_TYPE_CUSTOMER,
		Multiregion:           proto.Bool(key.MultiRegion),
		Creationdate:          key.CreationDate.Format(timeutils.ISO8601UTCFormat),
	}

	if key.DeletionDate != nil {
		md.Deletiondate = key.DeletionDate.Format(timeutils.ISO8601UTCFormat)
	}
	if key.KeyState == kmsstore.KeyStatePendingDeletion && key.PendingWindowInDays > 0 {
		md.Pendingdeletionwindowindays = proto.Int32(int32(key.PendingWindowInDays))
	}
	if key.ValidTo != nil {
		md.Validto = key.ValidTo.Format(timeutils.ISO8601UTCFormat)
	}
	if key.ExpirationModel == "KEY_MATERIAL_EXPIRES" {
		md.Expirationmodel = pb.ExpirationModelType_EXPIRATION_MODEL_TYPE_KEY_MATERIAL_EXPIRES
	} else if key.ExpirationModel == "KEY_MATERIAL_DOES_NOT_EXPIRE" {
		md.Expirationmodel = pb.ExpirationModelType_EXPIRATION_MODEL_TYPE_KEY_MATERIAL_DOES_NOT_EXPIRE
	}

	return md
}

func storeKeyStateToProto(state kmsstore.KeyState) pb.KeyState {
	switch state {
	case kmsstore.KeyStateEnabled:
		return pb.KeyState_KEY_STATE_ENABLED
	case kmsstore.KeyStateDisabled:
		return pb.KeyState_KEY_STATE_DISABLED
	case kmsstore.KeyStatePendingDeletion:
		return pb.KeyState_KEY_STATE_PENDINGDELETION
	case kmsstore.KeyStatePendingImport:
		return pb.KeyState_KEY_STATE_PENDINGIMPORT
	case kmsstore.KeyStateUnavailable:
		return pb.KeyState_KEY_STATE_UNAVAILABLE
	default:
		return pb.KeyState_KEY_STATE_ENABLED
	}
}

// NewConnectHandler creates a gRPC-Web connect handler for the Kms admin console.
func NewConnectHandler(svc *KMSService) (string, http.Handler) {
	return kmsconnect.NewKMSServiceHandler(NewAdminHandler(svc))
}
