package kms

import (
	"context"

	"connectrpc.com/connect"

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
