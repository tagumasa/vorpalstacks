package athena

import (
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// --- DTOs ---

// CreateCapacityReservationInput carries the parsed wire members of a
// CreateCapacityReservation request.
type CreateCapacityReservationInput struct {
	Name            string
	TargetWorkGroup string
	TargetDpus      int32
}

// ListCapacityReservationsInput carries the workgroup filter plus the raw
// MaxResults window (presence-flagged) and pagination marker.
type ListCapacityReservationsInput struct {
	WorkGroup     string
	MaxResults    int
	HasMaxResults bool
	NextToken     string
}

// UpdateCapacityReservationInput carries the parsed wire members of an
// UpdateCapacityReservation request (per the Smithy model, Name and
// TargetDpus are both REQUIRED).
type UpdateCapacityReservationInput struct {
	Name       string
	TargetDpus int32
}

// PutCapacityAssignmentConfigurationInput carries the parsed wire members of
// a PutCapacityAssignmentConfiguration request (assignments already
// converted from the raw wire array).
type PutCapacityAssignmentConfigurationInput struct {
	CapacityReservationName string
	CapacityAssignments     [][]string
}

// --- Core functions ---

// createCapacityReservationCore validates the create request and persists
// the reservation.
func createCapacityReservationCore(stores *athenaStores, input CreateCapacityReservationInput) (*athenastore.CapacityReservation, error) {
	if input.Name == "" {
		return nil, awserrors.NewAWSError("InvalidRequestException", "Name is required", http.StatusBadRequest)
	}
	if err := validateCapacityReservationName(input.Name); err != nil {
		return nil, err
	}
	targetWG := input.TargetWorkGroup
	if targetWG == "" {
		targetWG = "primary"
	}
	if err := validateTargetDpus(input.TargetDpus); err != nil {
		return nil, err
	}

	cr := &athenastore.CapacityReservation{
		Name:            input.Name,
		TargetWorkGroup: targetWG,
		Capacity:        input.TargetDpus,
	}

	created, err := stores.capacityReservationStore.CreateCapacityReservation(cr)
	if err != nil {
		if err == athenastore.ErrCapacityReservationAlreadyExists {
			return nil, alreadyExistsInvalidRequest("CapacityReservation", input.Name)
		}
		return nil, err
	}

	return created, nil
}

// getCapacityReservationCore fetches a capacity reservation, mapping any
// store error onto the documented InvalidRequestException (the capacity
// family declares no ResourceNotFoundException).
func getCapacityReservationCore(stores *athenaStores, name string) (*athenastore.CapacityReservation, error) {
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidRequestException", "Name is required", http.StatusBadRequest)
	}

	cr, err := stores.capacityReservationStore.GetCapacityReservation(name)
	if err != nil {
		return nil, capacityReservationNotFound(name)
	}
	return cr, nil
}

// listCapacityReservationsCore lists reservations optionally filtered by
// workgroup and pages them by name with the documented window semantics
// (default 50, range 1-50) applied before the walk, matching the original
// validation position.
func (s *AthenaService) listCapacityReservationsCore(reqCtx *request.RequestContext, input ListCapacityReservationsInput) ([]*athenastore.CapacityReservation, string, error) {
	maxResults, err := resolveMaxResults(input.MaxResults, input.HasMaxResults, 50, 1, 50)
	if err != nil {
		return nil, "", err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, "", err
	}

	reservations, err := stores.capacityReservationStore.ListCapacityReservations(input.WorkGroup)
	if err != nil {
		return nil, "", err
	}

	pageResult := pagination.PaginateSlice(reservations, input.NextToken, maxResults, func(item *athenastore.CapacityReservation) string {
		return item.Name
	})

	return pageResult.Items, pageResult.NextMarker, nil
}

// updateCapacityReservationCore validates the update request, applies the
// new DPU capacity and persists it.
func (s *AthenaService) updateCapacityReservationCore(reqCtx *request.RequestContext, input UpdateCapacityReservationInput) (*athenastore.CapacityReservation, error) {
	if input.Name == "" {
		return nil, awserrors.NewAWSError("InvalidRequestException", "Name is required", http.StatusBadRequest)
	}

	if err := validateTargetDpus(input.TargetDpus); err != nil {
		return nil, err
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cr, err := stores.capacityReservationStore.GetCapacityReservation(input.Name)
	if err != nil {
		return nil, capacityReservationNotFound(input.Name)
	}

	cr.Capacity = input.TargetDpus
	cr.LastModifiedTime = time.Now().UTC()

	if err := stores.capacityReservationStore.UpdateCapacityReservation(cr); err != nil {
		return nil, err
	}

	return cr, nil
}

// cancelCapacityReservationCore cancels an ACTIVE reservation, refusing any
// other state.
func cancelCapacityReservationCore(stores *athenaStores, name string) (*athenastore.CapacityReservation, error) {
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidRequestException", "Name is required", http.StatusBadRequest)
	}

	cr, err := stores.capacityReservationStore.GetCapacityReservation(name)
	if err != nil {
		return nil, capacityReservationNotFound(name)
	}

	if err := validateCapacityReservationStatusForCancel(string(cr.Status)); err != nil {
		return nil, err
	}

	cr.Status = athenastore.CapacityReservationStatusCancelled
	cr.LastModifiedTime = time.Now().UTC()
	if err := stores.capacityReservationStore.UpdateCapacityReservation(cr); err != nil {
		return nil, err
	}

	return cr, nil
}

// deleteCapacityReservationCore deletes a CANCELLED reservation, refusing
// any other state.
func deleteCapacityReservationCore(stores *athenaStores, name string) error {
	if name == "" {
		return awserrors.NewAWSError("InvalidRequestException", "Name is required", http.StatusBadRequest)
	}

	cr, err := stores.capacityReservationStore.GetCapacityReservation(name)
	if err != nil {
		return capacityReservationNotFound(name)
	}

	if err := validateCapacityReservationStatusForDelete(string(cr.Status)); err != nil {
		return err
	}

	if err := stores.capacityReservationStore.DeleteCapacityReservation(name); err != nil {
		return err
	}

	return nil
}

// getCapacityAssignmentConfigurationCore fetches the capacity assignment
// configuration of a reservation, mapping any store error onto the
// documented InvalidRequestException.
func (s *AthenaService) getCapacityAssignmentConfigurationCore(reqCtx *request.RequestContext, name string) ([][]string, error) {
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidRequestException", "CapacityReservationName is required", http.StatusBadRequest)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	assignments, err := stores.capacityReservationStore.GetCapacityAssignments(name)
	if err != nil {
		return nil, capacityReservationNotFound(name)
	}

	return assignments, nil
}

// putCapacityAssignmentConfigurationCore replaces the capacity assignment
// configuration of a reservation.
func (s *AthenaService) putCapacityAssignmentConfigurationCore(reqCtx *request.RequestContext, input PutCapacityAssignmentConfigurationInput) error {
	if input.CapacityReservationName == "" {
		return awserrors.NewAWSError("InvalidRequestException", "CapacityReservationName is required", http.StatusBadRequest)
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	if err := stores.capacityReservationStore.PutCapacityAssignments(input.CapacityReservationName, input.CapacityAssignments); err != nil {
		if err == athenastore.ErrCapacityReservationNotFound {
			return capacityReservationNotFound(input.CapacityReservationName)
		}
		return err
	}

	return nil
}
