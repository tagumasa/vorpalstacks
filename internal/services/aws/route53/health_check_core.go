package route53

import (
	"reflect"
	"strconv"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	route53store "vorpalstacks/internal/store/aws/route53"
)

// deletedHealthCheckRefRetention is how long a deleted health check keeps
// resolving its CallerReference. The API documents that a retry reusing the
// CallerReference of a deleted health check fails with
// HealthCheckAlreadyExists, but not for how long the reference is retained
// ("a number of days"); this is the platform's window.
const deletedHealthCheckRefRetention = 7 * 24 * time.Hour

// ---------------------------------------------------------------------------
// Transport-agnostic Input / Result structs
// ---------------------------------------------------------------------------

// CreateHealthCheckInput carries every field that CreateHealthCheck needs,
// independent of the wire protocol. The HTTP API handler parses the wire
// HealthCheckConfig into the store config shape and delegates validation and
// persistence to createHealthCheckCore.
type CreateHealthCheckInput struct {
	CallerReference string
	Config          *route53store.HealthCheckConfig
	Region          string
}

// ListHealthChecksInput carries the pagination parameters for ListHealthChecks.
type ListHealthChecksInput struct {
	Marker   string
	MaxItems int
}

// DeleteHealthCheckInput carries the parameters for DeleteHealthCheck.
type DeleteHealthCheckInput struct {
	Id string
}

// UpdateHealthCheckInput carries the parameters for UpdateHealthCheck.
// Updates is the raw HealthCheckConfig member map — or the full parameter
// map when the request carries the fields at the top level — whose keys
// follow the wire member names; applyHealthCheckConfigUpdates interprets
// them with presence semantics.
type UpdateHealthCheckInput struct {
	HealthCheckId      string
	HealthCheckVersion int
	Updates            map[string]interface{}
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// getHealthCheckCore looks up a health check by ID, mapping store misses to
// the NoSuchHealthCheck AWS error.
func getHealthCheckCore(stores *route53store.Route53Stores, id string) (*route53store.HealthCheck, error) {
	hc, err := stores.HealthChecks().Get(id)
	if err != nil {
		if route53store.IsNotFound(err) {
			return nil, NewNoSuchHealthCheckError(id)
		}
		return nil, mapStoreError(err)
	}
	return hc, nil
}

// createHealthCheckCore is the single entry point for creating a health
// check, shared by every protocol plane. It validates the parsed config,
// applies the CallerReference retry semantics, assigns the health check ID
// and persists it.
func (s *Route53Service) createHealthCheckCore(stores *route53store.Route53Stores, input CreateHealthCheckInput) (*route53store.HealthCheck, error) {
	if input.Config == nil {
		return nil, awserrors.NewAWSError("InvalidInput", "HealthCheckConfig is required", 400)
	}
	// CallerReference is a required member of CreateHealthCheckRequest; AWS
	// rejects an omitted member rather than synthesising an idempotency
	// token for the caller.
	if input.CallerReference == "" {
		return nil, awserrors.NewAWSError("InvalidInput", "CallerReference is required", 400)
	}
	// Validate Type enum and field ranges before persisting.
	if err := validateHealthCheckConfig(input.Config); err != nil {
		return nil, err
	}

	// CreateHealthCheck documents retry semantics for the CallerReference:
	// the same reference with identical settings returns the existing
	// health check, the same reference with different settings fails with
	// HealthCheckAlreadyExists, and so does a retry whose reference matches
	// a recently deleted health check (the reference is only retained for a
	// limited period after deletion).
	if existing, err := stores.HealthChecks().GetByCallerReference(input.CallerReference); err == nil && existing != nil {
		if existing.Deleted {
			if time.Since(existing.DeletedAt) < deletedHealthCheckRefRetention {
				return nil, NewHealthCheckAlreadyExistsError()
			}
		} else if reflect.DeepEqual(existing.HealthCheckConfig, input.Config) {
			return existing, nil
		} else {
			return nil, NewHealthCheckAlreadyExistsError()
		}
	}

	healthCheck := &route53store.HealthCheck{
		ID:                 generateHealthCheckId(),
		CallerReference:    input.CallerReference,
		HealthCheckConfig:  input.Config,
		HealthCheckVersion: "1",
		Region:             input.Region,
		AccountID:          s.accountID,
	}

	if err := stores.HealthChecks().Create(healthCheck); err != nil {
		if route53store.IsAlreadyExists(err) {
			return nil, NewHealthCheckAlreadyExistsError()
		}
		return nil, mapStoreError(err)
	}
	return healthCheck, nil
}

// listHealthChecksCore is the single entry point for listing health checks.
func listHealthChecksCore(stores *route53store.Route53Stores, input ListHealthChecksInput) (*route53store.HealthCheckListResult, error) {
	result, err := stores.HealthChecks().List(input.Marker, input.MaxItems)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return result, nil
}

// deleteHealthCheckCore is the single entry point for deleting a health
// check. It rejects deletion while any record set still references the
// health check and removes the health check's tags with it.
func deleteHealthCheckCore(stores *route53store.Route53Stores, input DeleteHealthCheckInput) error {
	// Reject deletion if any record set references this health check.
	if stores.RecordSets().IsHealthCheckReferenced(input.Id) {
		return NewHealthCheckInUseError(input.Id)
	}

	if err := stores.Tags().Raw().Delete("healthcheck/" + input.Id); err != nil {
		return awserrors.NewAWSError("InvalidInput", "Failed to delete tags: "+err.Error(), 500)
	}

	if err := stores.HealthChecks().Delete(input.Id); err != nil {
		if route53store.IsNotFound(err) {
			return NewNoSuchHealthCheckError(input.Id)
		}
		return mapStoreError(err)
	}
	return nil
}

// updateHealthCheckCore is the single entry point for updating a health
// check. It enforces the optional HealthCheckVersion match, applies the
// requested config updates, and validates the resulting configuration —
// the same constraints as createHealthCheckCore apply to updated values.
func updateHealthCheckCore(stores *route53store.Route53Stores, input UpdateHealthCheckInput) (*route53store.HealthCheck, error) {
	healthCheck, err := getHealthCheckCore(stores, input.HealthCheckId)
	if err != nil {
		return nil, err
	}

	// Verify HealthCheckVersion matches when provided.
	if input.HealthCheckVersion > 0 {
		existingVersion, _ := strconv.ParseInt(healthCheck.HealthCheckVersion, 10, 64)
		if int64(input.HealthCheckVersion) != existingVersion {
			return nil, NewHealthCheckVersionMismatchError(
				healthCheck.HealthCheckVersion, strconv.Itoa(input.HealthCheckVersion))
		}
	}

	if input.Updates != nil {
		applyHealthCheckConfigUpdates(healthCheck.HealthCheckConfig, input.Updates)
	}

	if err := validateHealthCheckConfig(healthCheck.HealthCheckConfig); err != nil {
		return nil, err
	}

	if err := stores.HealthChecks().Update(healthCheck); err != nil {
		return nil, mapStoreError(err)
	}
	return healthCheck, nil
}
