package kms

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"
	svccommon "vorpalstacks/internal/common"
	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/services/aws/kms/hsm"
	"vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/aws/types"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	pb "vorpalstacks/internal/pb/aws/kms"
	"vorpalstacks/internal/pb/aws/kms/kmsconnect"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// AdminHandler implements the KMS admin console gRPC-Web handler.
// It delegates to the shared KMSService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	kmsconnect.UnimplementedKMSServiceHandler
	service *KMSService
}

var _ kmsconnect.KMSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new KMS admin handler with the given key store.
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

// CreateKey creates a new KMS key via the admin console.
func (h *AdminHandler) CreateKey(ctx context.Context, req *connect.Request[pb.CreateKeyRequest]) (*connect.Response[pb.CreateKeyResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	keyUsage := kmsstore.KeyUsageEncryptDecrypt
	switch req.Msg.GetKeyusage() {
	case pb.KeyUsageType_KEY_USAGE_TYPE_SIGN_VERIFY:
		keyUsage = kmsstore.KeyUsageSignVerify
	case pb.KeyUsageType_KEY_USAGE_TYPE_GENERATE_VERIFY_MAC:
		keyUsage = kmsstore.KeyUsageGenerateVerifyMAC
	}

	keySpec := kmsstore.KeySpecSymmetricDefault
	switch req.Msg.GetKeyspec() {
	case pb.KeySpec_KEY_SPEC_SYMMETRIC_DEFAULT:
		keySpec = kmsstore.KeySpecSymmetricDefault
	case pb.KeySpec_KEY_SPEC_RSA_2048:
		keySpec = kmsstore.KeySpecRSA2048
	case pb.KeySpec_KEY_SPEC_RSA_3072:
		keySpec = kmsstore.KeySpecRSA3072
	case pb.KeySpec_KEY_SPEC_RSA_4096:
		keySpec = kmsstore.KeySpecRSA4096
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P256:
		keySpec = kmsstore.KeySpecECCNISTP256
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P384:
		keySpec = kmsstore.KeySpecECCNISTP384
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P521:
		keySpec = kmsstore.KeySpecECCNISTP521
	case pb.KeySpec_KEY_SPEC_ECC_SECG_P256K1:
		keySpec = kmsstore.KeySpecECCSECGP256K1
	case pb.KeySpec_KEY_SPEC_SM2:
		keySpec = kmsstore.KeySpecSM2
	case pb.KeySpec_KEY_SPEC_HMAC_224:
		keySpec = kmsstore.KeySpecHMAC224
	case pb.KeySpec_KEY_SPEC_HMAC_256:
		keySpec = kmsstore.KeySpecHMAC256
	case pb.KeySpec_KEY_SPEC_HMAC_384:
		keySpec = kmsstore.KeySpecHMAC384
	case pb.KeySpec_KEY_SPEC_HMAC_512:
		keySpec = kmsstore.KeySpecHMAC512
	}

	origin := kmsstore.OriginTypeAWSKMS
	switch req.Msg.GetOrigin() {
	case pb.OriginType_ORIGIN_TYPE_EXTERNAL:
		origin = kmsstore.OriginTypeExternal
	}

	keyID, err := kmsstore.GenerateKeyID()
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	key, err := stores.keys.Create(keyID, keyUsage, keySpec, req.Msg.GetDescription(), origin, req.Msg.GetMultiregion())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if origin == kmsstore.OriginTypeExternal {
		if err := stores.keys.SetPendingImport(keyID); err != nil {
			if delErr := stores.CascadeDeleteKey(h.service.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after SetPendingImport failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
		key.KeyState = kmsstore.KeyStatePendingImport
		key.Enabled = false
	} else {
		if err := h.service.hsmBackend.GenerateKey(keyID, hsm.KeySpec(keySpec)); err != nil {
			if delErr := stores.CascadeDeleteKey(h.service.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after HSM GenerateKey failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
	}

	// Apply caller-supplied policy. When omitted, inherit the primary's
	// policy or fall back to the default.
	keyPolicy := ""
	if req.Msg.GetPolicy() != "" {
		keyPolicy = req.Msg.GetPolicy()
	} else {
		keyPolicy = kmsstore.DefaultKeyPolicy
	}
	bypassLockoutCheck := req.Msg.GetBypasspolicylockoutsafetycheck()
	if !bypassLockoutCheck {
		if err := validatePolicyDoesNotLockOutRoot(keyPolicy); err != nil {
			if delErr := stores.CascadeDeleteKey(h.service.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after policy lockout check failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
	}
	key.BypassPolicyLockoutSafetyCheck = bypassLockoutCheck
	if err := stores.keys.Update(key); err != nil {
		if delErr := stores.CascadeDeleteKey(h.service.hsmBackend, keyID); delErr != nil {
			logs.Error("Failed to cascade-delete key after Update", logs.Err(delErr), logs.String("keyId", keyID))
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}
	if err := stores.keyPolicies.PutDefault(keyID, keyPolicy); err != nil {
		if delErr := stores.CascadeDeleteKey(h.service.hsmBackend, keyID); delErr != nil {
			logs.Error("Failed to cascade-delete key after PutDefault policy failure", logs.Err(delErr), logs.String("keyId", keyID))
		}
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	// Apply tags from the request.
	if len(req.Msg.GetTags()) > 0 {
		var tagList []types.Tag
		for _, t := range req.Msg.GetTags() {
			tagList = append(tagList, types.Tag{Key: t.GetTagkey(), Value: t.GetTagvalue()})
		}
		if err := validateKMSTags(tagList); err != nil {
			if delErr := stores.CascadeDeleteKey(h.service.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after tag validation failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, connect.NewError(connect.CodeInvalidArgument, err)
		}
		tagMap := make(map[string]string, len(tagList))
		for _, t := range tagList {
			tagMap[t.Key] = t.Value
		}
		if err := stores.keys.TagStore.Tag(keyID, tagMap); err != nil {
			if delErr := stores.CascadeDeleteKey(h.service.hsmBackend, keyID); delErr != nil {
				logs.Error("Failed to cascade-delete key after Tag failure", logs.Err(delErr), logs.String("keyId", keyID))
			}
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
	}

	return connect.NewResponse(&pb.CreateKeyResponse{
		Keymetadata: h.buildProtoKeyMetadata(key, req.Msg.GetKeyusage(), req.Msg.GetKeyspec(), req.Msg.GetOrigin()),
	}), nil
}

// ScheduleKeyDeletion schedules a KMS key for deletion via the admin console.
func (h *AdminHandler) ScheduleKeyDeletion(ctx context.Context, req *connect.Request[pb.ScheduleKeyDeletionRequest]) (*connect.Response[pb.ScheduleKeyDeletionResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	if req.Msg.GetKeyid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("KeyId is required"))
	}

	// AWS: PendingWindowInDays is optional (defaults to 30 when omitted)
	// but if explicitly supplied must be 7-30. Negative or zero values
	// are ValidationException. The previous code silently coerced
	// pendingDays <= 0 to 30, hiding invalid input from callers.
	pendingDays := int(req.Msg.GetPendingwindowindays())
	if pendingDays == 0 {
		pendingDays = 30
	}
	if pendingDays < 7 || pendingDays > 30 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("PendingWindowInDays must be between 7 and 30"))
	}

	if err := stores.keys.ScheduleDeletion(req.Msg.GetKeyid(), pendingDays); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	key, err := stores.keys.Get(req.Msg.GetKeyid())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	resp := &pb.ScheduleKeyDeletionResponse{
		Keyid:               key.KeyID,
		Pendingwindowindays: proto.Int32(int32(pendingDays)),
	}
	switch key.KeyState {
	case kmsstore.KeyStateEnabled:
		resp.Keystate = pb.KeyState_KEY_STATE_ENABLED
	case kmsstore.KeyStateDisabled:
		resp.Keystate = pb.KeyState_KEY_STATE_DISABLED
	case kmsstore.KeyStatePendingDeletion:
		resp.Keystate = pb.KeyState_KEY_STATE_PENDINGDELETION
	case kmsstore.KeyStatePendingImport:
		resp.Keystate = pb.KeyState_KEY_STATE_PENDINGIMPORT
	case kmsstore.KeyStateUnavailable:
		resp.Keystate = pb.KeyState_KEY_STATE_UNAVAILABLE
	}
	if key.DeletionDate != nil {
		resp.Deletiondate = key.DeletionDate.UTC().Format(timeutils.ISO8601UTCFormat)
	}
	return connect.NewResponse(resp), nil
}

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
