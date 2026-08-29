package iot

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"

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
// A suppression carries exactly one expiration: an expiration date, or an
// indefinite suppression (suppressIndefinitely true). An explicit false
// alongside a date is the pairing the documented CLI example sends and is
// accepted. The clientRequestToken member is required, and recreating an
// existing (checkName, resourceIdentifier) tuple raises
// ResourceAlreadyExistsException instead of overwriting.
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
	indefinite := in.SuppressIndefinitelyProvided && in.SuppressIndefinitely
	if (in.ExpirationProvided && indefinite) || (!in.ExpirationProvided && !indefinite) {
		return iotstore.ErrInvalidRequest
	}
	key := auditSuppressionKey(in.CheckName, in.ResourceIdentifier)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if exists {
		// The clientRequestToken member carries the idempotencyToken
		// trait: replaying the token the suppression was created with
		// returns success instead of the duplicate conflict.
		if in.ClientRequestToken != "" && rec["clientRequestToken"] == in.ClientRequestToken {
			return nil
		}
		return iotstore.ErrResourceAlreadyExists
	}
	// Each suppression must have a unique client request token, so the
	// token index rejects reusing an existing suppression's token for a
	// different tuple.
	if err := s.claimClientToken(store, "auditSuppression", in.ClientRequestToken, key); err != nil {
		return err
	}
	rec = map[string]interface{}{
		"checkName":            in.CheckName,
		"resourceIdentifier":   in.ResourceIdentifier,
		"expirationDate":       in.ExpirationDate,
		"suppressIndefinitely": in.SuppressIndefinitely,
		"description":          in.Description,
		"clientRequestToken":   in.ClientRequestToken,
	}
	return store.PutGeneric(key, rec)
}

// deleteAuditSuppressionCore removes an audit suppression and releases the
// create-time clientRequestToken index so the token may back a later
// create. An unknown (checkName, resourceIdentifier) tuple yields
// ErrAuditSuppressionNotFound.
func (s *IoTService) deleteAuditSuppressionCore(store iotstore.IotStoreInterface, checkName string, resourceIdentifier map[string]interface{}) error {
	if checkName == "" {
		return iotstore.ErrMissingParam
	}
	if len(resourceIdentifier) == 0 {
		return iotstore.ErrMissingParam
	}
	key := auditSuppressionKey(checkName, resourceIdentifier)
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrAuditSuppressionNotFound
	}
	if token, _ := rec["clientRequestToken"].(string); token != "" {
		if err := s.releaseClientToken(store, "auditSuppression", token); err != nil {
			return err
		}
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

// ListAuditSuppressionsInput carries the parsed ListAuditSuppressions
// request. AscendingOrder is tri-state: nil keeps the documented ascending
// default, an explicit false sorts by descending expiration date.
type ListAuditSuppressionsInput struct {
	CheckName          string
	ResourceIdentifier map[string]interface{}
	AscendingOrder     *bool
}

// listAuditSuppressionsCore lists audit suppressions filtered by the optional
// checkName/resourceIdentifier members, ordered by expiration date (ascending
// unless descending order is requested explicitly).
func (s *IoTService) listAuditSuppressionsCore(store iotstore.IotStoreInterface, in ListAuditSuppressionsInput) ([]*AuditSuppressionRecord, error) {
	items, err := store.ListGeneric("auditSuppression/")
	if err != nil {
		return nil, err
	}
	suppressions := make([]*AuditSuppressionRecord, 0, len(items))
	for _, rec := range items {
		if !suppressionMatchesFilter(rec, in) {
			continue
		}
		suppressions = append(suppressions, auditSuppressionFromRecord(rec))
	}
	sortAuditSuppressions(suppressions, in.AscendingOrder)
	return suppressions, nil
}

// suppressionMatchesFilter reports whether a stored suppression record
// matches the optional checkName and resourceIdentifier filters. The
// resourceIdentifier filter carries the AWS ResourceIdentifier shape; every
// non-empty filter member must equal the stored member.
func suppressionMatchesFilter(rec map[string]interface{}, in ListAuditSuppressionsInput) bool {
	if in.CheckName != "" && rec["checkName"] != in.CheckName {
		return false
	}
	if len(in.ResourceIdentifier) > 0 {
		stored, _ := rec["resourceIdentifier"].(map[string]interface{})
		for k, v := range in.ResourceIdentifier {
			if stored[k] != v {
				return false
			}
		}
	}
	return true
}

// sortAuditSuppressions orders entries by expiration date (epoch seconds)
// with the checkName digest as the deterministic tie-breaker; a nil or true
// ascendingOrder keeps the documented ascending default. An indefinite
// suppression never expires, so it sorts as the farthest date.
func sortAuditSuppressions(suppressions []*AuditSuppressionRecord, ascending *bool) {
	descending := ascending != nil && !*ascending
	sort.SliceStable(suppressions, func(i, j int) bool {
		iExp, jExp := suppressionSortEpoch(suppressions[i]), suppressionSortEpoch(suppressions[j])
		if iExp != jExp {
			if descending {
				return iExp > jExp
			}
			return iExp < jExp
		}
		if descending {
			return fmt.Sprintf("%v", suppressions[i].CheckName) > fmt.Sprintf("%v", suppressions[j].CheckName)
		}
		return fmt.Sprintf("%v", suppressions[i].CheckName) < fmt.Sprintf("%v", suppressions[j].CheckName)
	})
}

// suppressionSortEpoch is the ordering key for one suppression: the stored
// expiration epoch, or the maximum value for an indefinite suppression.
func suppressionSortEpoch(rec *AuditSuppressionRecord) int64 {
	if suppressIndefinite, ok := rec.SuppressIndefinitely.(bool); ok && suppressIndefinite {
		return math.MaxInt64
	}
	return recordEpoch(rec.ExpirationDate)
}

// recordEpoch coerces a stored epoch value (JSON number or string) to int64.
func recordEpoch(v interface{}) int64 {
	switch n := v.(type) {
	case int64:
		return n
	case float64:
		return int64(n)
	case int:
		return int64(n)
	case string:
		var parsed int64
		fmt.Sscanf(n, "%d", &parsed)
		return parsed
	}
	return 0
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
	// The exclusivity rule binds on update as on create: an expiration
	// date together with an indefinite suppression (suppressIndefinitely
	// true) is rejected, while an explicit false alongside a date is the
	// documented pairing and accepted; supplying neither keeps the stored
	// state (a description-only update).
	indefinite := in.SuppressIndefinitelyProvided && in.SuppressIndefinitely
	if in.ExpirationProvided && indefinite {
		return iotstore.ErrInvalidRequest
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
		// numeric epoch form regardless of how the client serialised it;
		// a dated expiration replaces any stored indefinite flag.
		rec["expirationDate"] = in.ExpirationDate
		if !indefinite {
			rec["suppressIndefinitely"] = false
		}
	} else if indefinite {
		// An indefinite suppression replaces any stored expiration date.
		rec["suppressIndefinitely"] = true
		rec["expirationDate"] = int64(0)
	}
	return store.PutGeneric(key, rec)
}
