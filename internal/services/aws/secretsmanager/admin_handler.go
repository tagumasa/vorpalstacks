package secretsmanager

import (
	"context"
	"fmt"
	"net/http"

	svcerrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/utils/timeutils"

	"connectrpc.com/connect"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/secretsmanager"
	secretsmanagerconnect "vorpalstacks/internal/pb/aws/secretsmanager/secretsmanagerconnect"
	"vorpalstacks/internal/store/aws/common"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
)

// AdminHandler implements the Secrets Manager admin console gRPC-Web handler.
// It delegates to the shared SecretsManagerService store cache so that the same
// per-region store instances are used by both the HTTP API handlers and the
// admin console gRPC-Web handlers.
type AdminHandler struct {
	secretsmanagerconnect.UnimplementedSecretsManagerServiceHandler
	service *SecretsManagerService
}

var _ secretsmanagerconnect.SecretsManagerServiceHandler = (*AdminHandler)(nil)

// NewAdminHandler creates a new Secrets Manager admin handler backed by the
// given service instance.
func NewAdminHandler(svc *SecretsManagerService) *AdminHandler {
	return &AdminHandler{
		service: svc,
	}
}

func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (*secretsmanagerstore.SecretStore, error) {
	region := svccommon.GetRegionFromHeader(headers)
	store, err := h.service.GetStoreForRegion(region)
	if err != nil {
		return nil, err
	}
	return store.(*secretsmanagerstore.SecretStore), nil
}

// ListSecrets returns all Secrets Manager secrets visible to the admin console.
func (h *AdminHandler) ListSecrets(ctx context.Context, req *connect.Request[pb.ListSecretsRequest]) (*connect.Response[pb.ListSecretsResponse], error) {
	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	maxResults := req.Msg.Maxresults
	if maxResults <= 0 {
		maxResults = 50
	}

	opts := common.ListOptions{
		MaxItems: int(maxResults),
		Marker:   req.Msg.Nexttoken,
	}
	result, err := store.ListSecrets(opts, nil)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	var entries []*pb.SecretListEntry
	for _, s := range result.Items {
		entry := &pb.SecretListEntry{
			Arn:         s.ARN,
			Name:        s.Name,
			Description: s.Description,
			Kmskeyid:    s.KmsKeyId,
			Type:        s.Type,
		}
		if !s.CreatedDate.IsZero() {
			entry.Createddate = s.CreatedDate.Format(timeutils.ISO8601UTCFormat)
		}
		if !s.LastChangedDate.IsZero() {
			entry.Lastchangeddate = s.LastChangedDate.Format(timeutils.ISO8601UTCFormat)
		}
		if !s.LastAccessedDate.IsZero() {
			entry.Lastaccesseddate = s.LastAccessedDate.Format(timeutils.ISO8601UTCFormat)
		}
		if !s.LastRotatedDate.IsZero() {
			entry.Lastrotateddate = s.LastRotatedDate.Format(timeutils.ISO8601UTCFormat)
		}
		if !s.NextRotationDate.IsZero() {
			entry.Nextrotationdate = s.NextRotationDate.Format(timeutils.ISO8601UTCFormat)
		}
		entry.Rotationenabled = s.RotationEnabled
		entry.Rotationlambdaarn = s.RotationLambdaARN
		if s.RotationRules != nil {
			entry.Rotationrules = &pb.RotationRulesType{
				Automaticallyafterdays: int64(s.RotationRules.AutomaticallyAfterDays),
			}
		}
		if s.DeletedDate != nil {
			entry.Deleteddate = s.DeletedDate.Format(timeutils.ISO8601UTCFormat)
		}
		entry.Owningservice = s.OwningService
		entry.Primaryregion = s.PrimaryRegion
		entries = append(entries, entry)
	}

	return connect.NewResponse(&pb.ListSecretsResponse{
		Secretlist: entries,
		Nexttoken:  result.NextMarker,
	}), nil
}

// CreateSecret creates a new secret via the admin console.
func (h *AdminHandler) CreateSecret(ctx context.Context, req *connect.Request[pb.CreateSecretRequest]) (*connect.Response[pb.CreateSecretResponse], error) {
	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	secret := &secretsmanagerstore.Secret{
		Name:         req.Msg.Name,
		Description:  req.Msg.Description,
		SecretString: req.Msg.Secretstring,
		SecretBinary: req.Msg.Secretbinary,
		KmsKeyId:     req.Msg.Kmskeyid,
		Type:         req.Msg.Type,
	}

	if len(req.Msg.Tags) > 0 {
		secret.Tags = make(map[string]string)
		for _, tag := range req.Msg.Tags {
			secret.Tags[tag.Key] = tag.Value
		}
	}

	created, err := store.CreateSecret(secret)
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	resp := &pb.CreateSecretResponse{
		Arn:  created.ARN,
		Name: created.Name,
	}
	if created.CurrentVersion != "" {
		resp.Versionid = created.CurrentVersion
	}

	return connect.NewResponse(resp), nil
}

// DeleteSecret deletes a secret via the admin console.
func (h *AdminHandler) DeleteSecret(ctx context.Context, req *connect.Request[pb.DeleteSecretRequest]) (*connect.Response[pb.DeleteSecretResponse], error) {
	if req.Msg.Secretid == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("SecretId is required"))
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	secret, getErr := store.GetSecretForMetadata(req.Msg.Secretid)
	if getErr != nil {
		return nil, svcerrors.StoreErrorToGRPC(getErr)
	}

	if err := store.DeleteSecret(req.Msg.Secretid); err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	resp := &pb.DeleteSecretResponse{
		Arn:  secret.ARN,
		Name: secret.Name,
	}

	return connect.NewResponse(resp), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Secrets Manager admin console.
func NewConnectHandler(svc *SecretsManagerService) (string, http.Handler) {
	return secretsmanagerconnect.NewSecretsManagerServiceHandler(NewAdminHandler(svc))
}
