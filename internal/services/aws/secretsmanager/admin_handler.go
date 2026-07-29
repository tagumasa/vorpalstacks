package secretsmanager

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/protobuf/proto"
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
		maxResults = 100
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
		entry.Rotationenabled = proto.Bool(s.RotationEnabled)
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
	// Align with HTTP API validation (Batch 7).
	if len(req.Msg.Name) > 512 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("secret name must be between 1 and 512 characters long"))
	}
	if len(req.Msg.Description) > 2048 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("description must be between 0 and 2048 characters long"))
	}
	// SecretString and SecretBinary are mutually exclusive.
	if req.Msg.Secretstring != "" && len(req.Msg.Secretbinary) > 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("you can't specify both SecretString and SecretBinary in the same request"))
	}
	// ClientRequestToken length 32-64 when provided.
	if req.Msg.Clientrequesttoken != "" {
		if len(req.Msg.Clientrequesttoken) < 32 || len(req.Msg.Clientrequesttoken) > 64 {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("ClientRequestToken must be 32 to 64 characters long"))
		}
	}
	// Validate tag quotas (count, key length, value length).
	if len(req.Msg.Tags) > maxTagsPerSecret {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("you can't have more than %d tags on a secret", maxTagsPerSecret))
	}
	for _, tag := range req.Msg.Tags {
		if tag.Key == "" || len(tag.Key) > maxTagKeyLength {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("tag key length must be between 1 and %d characters", maxTagKeyLength))
		}
		if len(tag.Value) > maxTagValueLength {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("tag value length must be between 0 and %d characters", maxTagValueLength))
		}
	}

	store, err := h.getStoreFromHeaders(req.Header())
	if err != nil {
		return nil, svcerrors.StoreErrorToGRPC(err)
	}

	secret := &secretsmanagerstore.Secret{
		Name:             req.Msg.Name,
		Description:      req.Msg.Description,
		SecretString:     req.Msg.Secretstring,
		SecretBinary:     req.Msg.Secretbinary,
		KmsKeyId:         req.Msg.Kmskeyid,
		Type:             req.Msg.Type,
		InitialVersionId: req.Msg.Clientrequesttoken,
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

	forceDelete := false
	if req.Msg.Forcedeletewithoutrecovery != nil {
		forceDelete = *req.Msg.Forcedeletewithoutrecovery
	}
	hasRecoveryWindow := req.Msg.Recoverywindowindays > 0

	// You can't use both ForceDeleteWithoutRecovery and RecoveryWindowInDays.
	if forceDelete && hasRecoveryWindow {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("you can't use ForceDeleteWithoutRecovery in conjunction with RecoveryWindowInDays"))
	}
	// RecoveryWindowInDays must be between 7 and 30 (inclusive).
	if hasRecoveryWindow && !forceDelete {
		if req.Msg.Recoverywindowindays < 7 || req.Msg.Recoverywindowindays > 30 {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				fmt.Errorf("RecoveryWindowInDays must be between 7 and 30 days"))
		}
	}

	// You can't delete a primary secret that is replicated to other Regions.
	if len(secret.ReplicationStatus) > 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument,
			fmt.Errorf("You can't delete a primary secret that is replicated to other Regions. Remove the replicas first."))
	}

	var deletionDate time.Time
	if forceDelete {
		if err := store.DeleteSecret(req.Msg.Secretid); err != nil {
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
		deletionDate = time.Now().UTC()
	} else {
		recoveryWindow := int(req.Msg.Recoverywindowindays)
		if recoveryWindow == 0 {
			recoveryWindow = 30
		}
		deletionDate = time.Now().UTC().AddDate(0, 0, recoveryWindow)
		if err := store.ScheduleDeletion(req.Msg.Secretid, deletionDate); err != nil {
			return nil, svcerrors.StoreErrorToGRPC(err)
		}
	}

	resp := &pb.DeleteSecretResponse{
		Arn:          secret.ARN,
		Name:         secret.Name,
		Deletiondate: deletionDate.Format(time.RFC3339),
	}

	return connect.NewResponse(resp), nil
}

// NewConnectHandler creates a gRPC-Web connect handler for the Secrets Manager admin console.
func NewConnectHandler(svc *SecretsManagerService) (string, http.Handler) {
	return secretsmanagerconnect.NewSecretsManagerServiceHandler(NewAdminHandler(svc))
}
