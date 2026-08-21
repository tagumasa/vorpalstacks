package kms

import (
	"context"
	"fmt"
	"net/http"
	"vorpalstacks/internal/common/defaults"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
	svcerrors "vorpalstacks/internal/common/errors"
	types "vorpalstacks/internal/common/tags"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"

	pb "vorpalstacks/internal/pb/aws/kms"
	"vorpalstacks/internal/pb/aws/kms/kmsconnect"
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
	region := defaults.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ListKeys retrieves all KMS keys with pagination.
func (h *AdminHandler) ListKeys(ctx context.Context, req *connect.Request[pb.ListKeysRequest]) (*connect.Response[pb.ListKeysResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	result, err := h.service.listKeysCore(stores, req.Msg.Marker, int(req.Msg.GetLimit()))
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
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
		return nil, svcerrors.AWSErrorToGRPC(err)
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

	// Convert proto tags.
	var tags []types.Tag
	for _, t := range req.Msg.GetTags() {
		tags = append(tags, types.Tag{Key: t.GetTagkey(), Value: t.GetTagvalue()})
	}

	meta, err := h.service.createKeyCore(stores, CreateKeyInput{
		Description:        req.Msg.GetDescription(),
		KeyUsage:           protoKeyUsageToString(keyUsageProto),
		KeySpec:            protoKeySpecToString(keySpecProto),
		Origin:             protoOriginToString(originProto),
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
		Keymetadata: buildProtoKeyMetadata(meta, keyUsageProto, keySpecProto, originProto),
	}), nil
}

// ScheduleKeyDeletion schedules a KMS key for deletion via the admin
// console. It resolves the key through the service layer (supporting
// alias and ARN forms) and delegates to scheduleKeyDeletionCore.
func (h *AdminHandler) ScheduleKeyDeletion(ctx context.Context, req *connect.Request[pb.ScheduleKeyDeletionRequest]) (*connect.Response[pb.ScheduleKeyDeletionResponse], error) {
	stores, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	if req.Msg.GetKeyid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("KeyId is required"))
	}

	// Resolve the key through the service layer to support alias and ARN.
	resolvedKeyID, err := h.service.resolveKeyID(stores, req.Msg.GetKeyid())
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

	meta, days, err := h.service.scheduleKeyDeletionCore(stores, resolvedKeyID, pendingDays)
	if err != nil {
		return nil, svcerrors.AWSErrorToGRPC(err)
	}

	resp := &pb.ScheduleKeyDeletionResponse{
		Keyid:               meta.KeyID,
		Keystate:            keyStateToProto(meta.KeyState),
		Pendingwindowindays: proto.Int32(int32(days)),
	}
	if meta.DeletionDate != nil {
		resp.Deletiondate = meta.DeletionDate.UTC().Format(timeutils.ISO8601UTCFormat)
	}
	return connect.NewResponse(resp), nil
}

// ---------------------------------------------------------------------------
// Proto conversion helpers (no store imports)
// ---------------------------------------------------------------------------

func protoKeyUsageToString(v pb.KeyUsageType) string {
	switch v {
	case pb.KeyUsageType_KEY_USAGE_TYPE_SIGN_VERIFY:
		return "SIGN_VERIFY"
	case pb.KeyUsageType_KEY_USAGE_TYPE_GENERATE_VERIFY_MAC:
		return "GENERATE_VERIFY_MAC"
	default:
		return "ENCRYPT_DECRYPT"
	}
}

func protoKeySpecToString(v pb.KeySpec) string {
	switch v {
	case pb.KeySpec_KEY_SPEC_SYMMETRIC_DEFAULT:
		return "SYMMETRIC_DEFAULT"
	case pb.KeySpec_KEY_SPEC_RSA_2048:
		return "RSA_2048"
	case pb.KeySpec_KEY_SPEC_RSA_3072:
		return "RSA_3072"
	case pb.KeySpec_KEY_SPEC_RSA_4096:
		return "RSA_4096"
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P256:
		return "ECC_NIST_P256"
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P384:
		return "ECC_NIST_P384"
	case pb.KeySpec_KEY_SPEC_ECC_NIST_P521:
		return "ECC_NIST_P521"
	case pb.KeySpec_KEY_SPEC_ECC_SECG_P256K1:
		return "ECC_SECG_P256K1"
	case pb.KeySpec_KEY_SPEC_SM2:
		return "SM2"
	case pb.KeySpec_KEY_SPEC_HMAC_224:
		return "HMAC_224"
	case pb.KeySpec_KEY_SPEC_HMAC_256:
		return "HMAC_256"
	case pb.KeySpec_KEY_SPEC_HMAC_384:
		return "HMAC_384"
	case pb.KeySpec_KEY_SPEC_HMAC_512:
		return "HMAC_512"
	default:
		return "SYMMETRIC_DEFAULT"
	}
}

func protoOriginToString(v pb.OriginType) string {
	switch v {
	case pb.OriginType_ORIGIN_TYPE_EXTERNAL:
		return "EXTERNAL"
	default:
		return "AWS_KMS"
	}
}

// keyStateToProto maps a key-state string to the proto enum.
func keyStateToProto(state string) pb.KeyState {
	switch state {
	case "Enabled":
		return pb.KeyState_KEY_STATE_ENABLED
	case "Disabled":
		return pb.KeyState_KEY_STATE_DISABLED
	case "PendingDeletion":
		return pb.KeyState_KEY_STATE_PENDINGDELETION
	case "PendingImport":
		return pb.KeyState_KEY_STATE_PENDINGIMPORT
	case "Unavailable":
		return pb.KeyState_KEY_STATE_UNAVAILABLE
	default:
		return pb.KeyState_KEY_STATE_ENABLED
	}
}

// buildProtoKeyMetadata constructs a full pb.KeyMetadata from a
// KeyMetadataResult, mapping all fields that the HTTP API's
// buildKeyMetadata also emits.
func buildProtoKeyMetadata(meta *KeyMetadataResult, keyUsage pb.KeyUsageType, keySpec pb.KeySpec, origin pb.OriginType) *pb.KeyMetadata {
	_, _, _, accountID, _ := arnutil.SplitARN(meta.Arn)

	md := &pb.KeyMetadata{
		Awsaccountid:          accountID,
		Keyid:                 meta.KeyID,
		Arn:                   meta.Arn,
		Keystate:              keyStateToProto(meta.KeyState),
		Keyusage:              keyUsage,
		Keyspec:               keySpec,
		Customermasterkeyspec: pb.CustomerMasterKeySpec(keySpec),
		Description:           meta.Description,
		Enabled:               proto.Bool(meta.Enabled),
		Origin:                origin,
		Keymanager:            pb.KeyManagerType_KEY_MANAGER_TYPE_CUSTOMER,
		Multiregion:           proto.Bool(meta.MultiRegion),
		Creationdate:          meta.CreationDate.Format(timeutils.ISO8601UTCFormat),
	}

	if meta.DeletionDate != nil {
		md.Deletiondate = meta.DeletionDate.Format(timeutils.ISO8601UTCFormat)
	}
	if meta.KeyState == "PendingDeletion" && meta.PendingWindowInDays > 0 {
		md.Pendingdeletionwindowindays = proto.Int32(int32(meta.PendingWindowInDays))
	}
	if meta.ValidTo != nil {
		md.Validto = meta.ValidTo.Format(timeutils.ISO8601UTCFormat)
	}
	if meta.ExpirationModel == "KEY_MATERIAL_EXPIRES" {
		md.Expirationmodel = pb.ExpirationModelType_EXPIRATION_MODEL_TYPE_KEY_MATERIAL_EXPIRES
	} else if meta.ExpirationModel == "KEY_MATERIAL_DOES_NOT_EXPIRE" {
		md.Expirationmodel = pb.ExpirationModelType_EXPIRATION_MODEL_TYPE_KEY_MATERIAL_DOES_NOT_EXPIRE
	}

	return md
}

// NewConnectHandler creates a gRPC-Web connect handler for the Kms admin console.
func NewConnectHandler(svc *KMSService) (string, http.Handler) {
	return kmsconnect.NewKMSServiceHandler(NewAdminHandler(svc))
}
