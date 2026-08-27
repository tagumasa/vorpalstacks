package iot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Audit Suppression Core. AWS identifies suppressions by the
// (checkName, resourceIdentifier) tuple; resourceIdentifier is a structure
// with up to ten optional member fields, so a canonical-JSON SHA-256 digest
// gives a stable, collision-free key suffix without coupling the key layout
// to AWS-internal representation. Records live under
// "auditSuppression/<checkName>/<digest>".
// ---------------------------------------------------------------------------

// auditSuppressionKey builds the generic-KV key for an AuditSuppression
// record.
func auditSuppressionKey(checkName string, resourceIdentifier map[string]interface{}) string {
	canonical, _ := json.Marshal(resourceIdentifier)
	digest := sha256.Sum256(canonical)
	return fmt.Sprintf("auditSuppression/%s/%s", checkName, hex.EncodeToString(digest[:])[:16])
}

// CreateAuditSuppressionInput carries the parsed CreateAuditSuppression
// request. ExpirationDate is Unix epoch seconds (0 = absent); the
// *Provided flags distinguish explicitly supplied members from omitted ones
// because the exclusivity rule is evaluated on presence, not value.
type CreateAuditSuppressionInput struct {
	CheckName                    string
	ResourceIdentifier           map[string]interface{}
	ExpirationDate               int64
	ExpirationProvided           bool
	SuppressIndefinitely         bool
	SuppressIndefinitelyProvided bool
	Description                  string
	ClientRequestToken           string
}

// createAuditSuppressionCore validates and persists an audit suppression.
// The expirationDate and suppressIndefinitely members are mutually
// exclusive: exactly one must be specified. The clientRequestToken member
// is required, and recreating an existing (checkName, resourceIdentifier)
// tuple raises ResourceAlreadyExistsException instead of overwriting.
func (s *IoTService) createAuditSuppressionCore(store iotstore.IotStoreInterface, in CreateAuditSuppressionInput) error {
	if in.CheckName == "" {
		return iotstore.ErrMissingParam
	}
	if len(in.ResourceIdentifier) == 0 {
		return iotstore.ErrMissingParam
	}
	if in.ClientRequestToken == "" {
		return iotstore.ErrMissingParam
	}
	if in.ExpirationProvided == in.SuppressIndefinitelyProvided {
		return iotstore.ErrInvalidRequest
	}
	key := auditSuppressionKey(in.CheckName, in.ResourceIdentifier)
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if exists {
		return iotstore.ErrResourceAlreadyExists
	}
	rec := map[string]interface{}{
		"checkName":            in.CheckName,
		"resourceIdentifier":   in.ResourceIdentifier,
		"expirationDate":       in.ExpirationDate,
		"suppressIndefinitely": in.SuppressIndefinitely,
		"description":          in.Description,
		"clientRequestToken":   in.ClientRequestToken,
	}
	return store.PutGeneric(key, rec)
}

// deleteAuditSuppressionCore removes an audit suppression. An unknown
// (checkName, resourceIdentifier) tuple yields
// ErrAuditSuppressionNotFound.
func (s *IoTService) deleteAuditSuppressionCore(store iotstore.IotStoreInterface, checkName string, resourceIdentifier map[string]interface{}) error {
	if checkName == "" {
		return iotstore.ErrMissingParam
	}
	if len(resourceIdentifier) == 0 {
		return iotstore.ErrMissingParam
	}
	key := auditSuppressionKey(checkName, resourceIdentifier)
	exists, err := store.GetGenericExists(key, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrAuditSuppressionNotFound
	}
	return store.DeleteGeneric(key)
}

// AuditSuppressionRecord is the projected DescribeAuditSuppression response
// payload.
type AuditSuppressionRecord struct {
	CheckName            interface{}
	ResourceIdentifier   interface{}
	ExpirationDate       interface{}
	SuppressIndefinitely interface{}
	Description          interface{}
}

// describeAuditSuppressionCore loads an audit suppression. An unknown
// (checkName, resourceIdentifier) tuple yields
// ErrAuditSuppressionNotFound.
func (s *IoTService) describeAuditSuppressionCore(store iotstore.IotStoreInterface, checkName string, resourceIdentifier map[string]interface{}) (*AuditSuppressionRecord, error) {
	if checkName == "" {
		return nil, iotstore.ErrMissingParam
	}
	if len(resourceIdentifier) == 0 {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(auditSuppressionKey(checkName, resourceIdentifier), &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrAuditSuppressionNotFound
	}
	return auditSuppressionFromRecord(rec), nil
}

func auditSuppressionFromRecord(rec map[string]interface{}) *AuditSuppressionRecord {
	return &AuditSuppressionRecord{
		CheckName:            rec["checkName"],
		ResourceIdentifier:   rec["resourceIdentifier"],
		ExpirationDate:       rec["expirationDate"],
		SuppressIndefinitely: rec["suppressIndefinitely"],
		Description:          rec["description"],
	}
}

// listAuditSuppressionsCore lists every audit suppression projected onto the
// AuditSuppression member set.
func (s *IoTService) listAuditSuppressionsCore(store iotstore.IotStoreInterface) ([]*AuditSuppressionRecord, error) {
	items, err := store.ListGeneric("auditSuppression/")
	if err != nil {
		return nil, err
	}
	suppressions := make([]*AuditSuppressionRecord, 0, len(items))
	for _, rec := range items {
		suppressions = append(suppressions, auditSuppressionFromRecord(rec))
	}
	return suppressions, nil
}

// UpdateAuditSuppressionInput carries the parsed UpdateAuditSuppression
// request. The *Provided flags distinguish an explicitly supplied member
// from an omitted one; only supplied members overwrite the record.
type UpdateAuditSuppressionInput struct {
	CheckName                    string
	ResourceIdentifier           map[string]interface{}
	Description                  string
	ExpirationDate               int64
	ExpirationProvided           bool
	SuppressIndefinitely         bool
	SuppressIndefinitelyProvided bool
}

// updateAuditSuppressionCore applies a partial update to an existing audit
// suppression. An unknown (checkName, resourceIdentifier) tuple yields
// ErrAuditSuppressionNotFound.
func (s *IoTService) updateAuditSuppressionCore(store iotstore.IotStoreInterface, in UpdateAuditSuppressionInput) error {
	if in.CheckName == "" {
		return iotstore.ErrMissingParam
	}
	if len(in.ResourceIdentifier) == 0 {
		return iotstore.ErrMissingParam
	}
	key := auditSuppressionKey(in.CheckName, in.ResourceIdentifier)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrAuditSuppressionNotFound
	}
	// Partial update: only overwrite fields that are explicitly supplied.
	if in.Description != "" {
		rec["description"] = in.Description
	}
	if in.ExpirationProvided {
		// Normalise like the create path so the stored value keeps its
		// numeric epoch form regardless of how the client serialised it.
		rec["expirationDate"] = in.ExpirationDate
	}
	if in.SuppressIndefinitelyProvided {
		rec["suppressIndefinitely"] = in.SuppressIndefinitely
	}
	return store.PutGeneric(key, rec)
}
