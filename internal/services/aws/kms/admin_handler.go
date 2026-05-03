package kms

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"

	"vorpalstacks/internal/core/storage"
	pb "vorpalstacks/internal/pb/aws/kms"
	"vorpalstacks/internal/pb/aws/kms/kmsconnect"
	kmsstore "vorpalstacks/internal/store/aws/kms"
)

// AdminHandler implements the KMS admin console gRPC-Web handler.
type AdminHandler struct {
	kmsconnect.UnimplementedKMSServiceHandler
	store *kmsstore.KeyStore
}

var _ kmsconnect.KMSServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new KMS admin handler with the given key store.
func NewAdminHandler(store *kmsstore.KeyStore) *AdminHandler {
	return &AdminHandler{store: store}
}

// ListKeys retrieves all KMS keys from the key store with pagination.
func (h *AdminHandler) ListKeys(ctx context.Context, req *connect.Request[pb.ListKeysRequest]) (*connect.Response[pb.ListKeysResponse], error) {
	result, err := h.store.List(req.Msg.Marker, int(req.Msg.Limit))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
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
		Truncated:  result.IsTruncated,
	}), nil
}

// CreateKey creates a new KMS key via the admin console.
func (h *AdminHandler) CreateKey(ctx context.Context, req *connect.Request[pb.CreateKeyRequest]) (*connect.Response[pb.CreateKeyResponse], error) {
	keyUsage := kmsstore.KeyUsageEncryptDecrypt
	switch req.Msg.GetKeyusage() {
	case pb.KeyUsageType_KEY_USAGE_TYPE_SIGN_VERIFY:
		keyUsage = kmsstore.KeyUsageSignVerify
	case pb.KeyUsageType_KEY_USAGE_TYPE_GENERATE_VERIFY_MAC:
		keyUsage = kmsstore.KeyUsageGenerateVerifyMAC
	}

	keySpec := kmsstore.KeySpecSymmetricDefault
	switch req.Msg.GetKeyspec() {
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
	case pb.KeySpec_KEY_SPEC_HMAC_256:
		keySpec = kmsstore.KeySpecHMAC256
	}

	origin := kmsstore.OriginTypeAWSKMS
	switch req.Msg.GetOrigin() {
	case pb.OriginType_ORIGIN_TYPE_EXTERNAL:
		origin = kmsstore.OriginTypeExternal
	}

	keyID := ""
	key, err := h.store.Create(keyID, keyUsage, keySpec, req.Msg.GetDescription(), origin, req.Msg.GetMultiregion())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateKeyResponse{
		Keymetadata: &pb.KeyMetadata{
			Keyid:       key.KeyID,
			Arn:         key.Arn,
			Keystate:    pb.KeyState_KEY_STATE_ENABLED,
			Keyusage:    req.Msg.GetKeyusage(),
			Keyspec:     req.Msg.GetKeyspec(),
			Description: key.Description,
			Enabled:     key.Enabled,
			Origin:      req.Msg.GetOrigin(),
			Keymanager:  pb.KeyManagerType_KEY_MANAGER_TYPE_CUSTOMER,
			Multiregion: key.MultiRegion,
			Creationdate: fmt.Sprintf("%d", key.CreationDate.Unix()),
		},
	}), nil
}

// ScheduleKeyDeletion schedules a KMS key for deletion via the admin console.
func (h *AdminHandler) ScheduleKeyDeletion(ctx context.Context, req *connect.Request[pb.ScheduleKeyDeletionRequest]) (*connect.Response[pb.ScheduleKeyDeletionResponse], error) {
	if req.Msg.GetKeyid() == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("KeyId is required"))
	}

	pendingDays := int(req.Msg.GetPendingwindowindays())
	if pendingDays <= 0 {
		pendingDays = 30
	}

	if err := h.store.ScheduleDeletion(req.Msg.GetKeyid(), pendingDays); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	key, err := h.store.Get(req.Msg.GetKeyid())
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &pb.ScheduleKeyDeletionResponse{
		Keyid: key.KeyID,
	}
	if key.DeletionDate != nil {
		resp.Deletiondate = fmt.Sprintf("%d", key.DeletionDate.Unix())
	}
	return connect.NewResponse(resp), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Kms admin console.
func NewConnectHandler(st storage.BasicStorage, accountID, region string) (string, http.Handler) {
	return kmsconnect.NewKMSServiceHandler(NewAdminHandler(kmsstore.NewKeyStore(st, accountID, region)))
}
