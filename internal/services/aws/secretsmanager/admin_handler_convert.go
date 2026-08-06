package secretsmanager

import (
	"net/http"
	"time"

	"google.golang.org/protobuf/proto"

	svccommon "vorpalstacks/internal/common"
	pb "vorpalstacks/internal/pb/aws/secretsmanager"
	secretsmanagerstore "vorpalstacks/internal/store/aws/secretsmanager"
	"vorpalstacks/internal/utils/timeutils"
)

// getStoreFromHeaders resolves the region from request headers and returns
// the per-region SecretStoreInterface. This is the sole function in the
// admin handler layer that touches the store package directly.
func (h *AdminHandler) getStoreFromHeaders(headers http.Header) (secretsmanagerstore.SecretStoreInterface, error) {
	region := svccommon.GetRegionFromHeader(headers)
	return h.service.GetStoreForRegion(region)
}

// ---------------------------------------------------------------------------
// Proto → Input converters
// ---------------------------------------------------------------------------

// pbToCreateSecretInput converts a proto CreateSecretRequest to the
// transport-agnostic CreateSecretInput.
func pbToCreateSecretInput(req *pb.CreateSecretRequest) CreateSecretInput {
	in := CreateSecretInput{
		Name:               req.Name,
		SecretString:       req.Secretstring,
		SecretBinary:       req.Secretbinary,
		Description:        req.Description,
		KmsKeyId:           req.Kmskeyid,
		Type:               req.Type,
		ClientRequestToken: req.Clientrequesttoken,
	}
	if len(req.Tags) > 0 {
		in.Tags = make(map[string]string, len(req.Tags))
		for _, tag := range req.Tags {
			in.Tags[tag.Key] = tag.Value
		}
	}
	return in
}

// pbToDeleteSecretInput converts a proto DeleteSecretRequest to the
// transport-agnostic DeleteSecretInput.
func pbToDeleteSecretInput(req *pb.DeleteSecretRequest) DeleteSecretInput {
	in := DeleteSecretInput{
		SecretId: req.Secretid,
	}
	if req.Forcedeletewithoutrecovery != nil {
		in.ForceDeleteWithoutRecovery = *req.Forcedeletewithoutrecovery
	}
	if req.Recoverywindowindays != nil {
		in.RecoveryWindowInDays = int(*req.Recoverywindowindays)
		in.HasRecoveryWindow = true
	}
	return in
}

// pbToListSecretsInput converts a proto ListSecretsRequest to the
// transport-agnostic ListSecretsInput.
func pbToListSecretsInput(req *pb.ListSecretsRequest) ListSecretsInput {
	in := ListSecretsInput{
		NextToken: req.Nexttoken,
	}
	if req.Maxresults != nil {
		in.MaxResults = int(*req.Maxresults)
	}
	return in
}

// ---------------------------------------------------------------------------
// Result → Proto converters
// ---------------------------------------------------------------------------

// createSecretResultToPb converts a CreateSecretResult to a proto response.
func createSecretResultToPb(r *CreateSecretResult) *pb.CreateSecretResponse {
	resp := &pb.CreateSecretResponse{
		Arn:  r.ARN,
		Name: r.Name,
	}
	if r.VersionID != "" {
		resp.Versionid = r.VersionID
	}
	return resp
}

// deleteSecretResultToPb converts a DeleteSecretResult to a proto response.
func deleteSecretResultToPb(r *DeleteSecretResult) *pb.DeleteSecretResponse {
	return &pb.DeleteSecretResponse{
		Arn:          r.ARN,
		Name:         r.Name,
		Deletiondate: r.DeletionDate.Format(time.RFC3339),
	}
}

// secretToPbEntry converts a store Secret to a proto SecretListEntry.
func secretToPbEntry(s *secretsmanagerstore.Secret) *pb.SecretListEntry {
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
			Automaticallyafterdays: proto.Int64(int64(s.RotationRules.AutomaticallyAfterDays)),
		}
	}
	if s.DeletedDate != nil {
		entry.Deleteddate = s.DeletedDate.Format(timeutils.ISO8601UTCFormat)
	}
	entry.Owningservice = s.OwningService
	entry.Primaryregion = s.PrimaryRegion
	return entry
}
