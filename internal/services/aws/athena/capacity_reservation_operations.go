package athena

import (
	"context"
	"fmt"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	athenastore "vorpalstacks/internal/store/aws/athena"
)

// formatCapacityReservation builds the response map for a capacity reservation.
func formatCapacityReservation(cr *athenastore.CapacityReservation) map[string]interface{} {
	resp := map[string]interface{}{
		"Name":            cr.Name,
		"TargetWorkGroup": cr.TargetWorkGroup,
		"Capacity":        cr.Capacity,
		"Status":          cr.Status,
		"CreationTime":    cr.CreationTime,
	}
	if !cr.ExpectedDispatchTime.IsZero() {
		resp["ExpectedDispatchTime"] = cr.ExpectedDispatchTime
	}
	if !cr.LastModifiedTime.IsZero() {
		resp["LastModifiedTime"] = cr.LastModifiedTime
	}
	return resp
}

// CreateCapacityReservation creates a new Athena capacity reservation.
func (s *AthenaService) CreateCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Name is required")
	}
	targetWG := request.GetStringParam(req.Parameters, "TargetWorkGroup")
	if targetWG == "" {
		targetWG = "primary"
	}
	capacity := int32(request.GetIntParam(req.Parameters, "TargetDpus"))
	if capacity <= 0 {
		return nil, awserrors.NewValidationException("TargetDpus must be positive")
	}

	cr := &athenastore.CapacityReservation{
		Name:            name,
		TargetWorkGroup: targetWG,
		Capacity:        capacity,
	}

	created, err := stores.capacityReservationStore.CreateCapacityReservation(cr)
	if err != nil {
		if err == athenastore.ErrCapacityReservationAlreadyExists {
			return nil, awserrors.NewResourceAlreadyExistsException(fmt.Sprintf("capacity reservation %s", name))
		}
		return nil, err
	}

	return formatCapacityReservation(created), nil
}

// GetCapacityReservation retrieves details about the specified capacity reservation.
func (s *AthenaService) GetCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Name is required")
	}

	cr, err := stores.capacityReservationStore.GetCapacityReservation(name)
	if err != nil {
		return nil, awserrors.NewResourceNotFoundException("CapacityReservation", name)
	}

	return formatCapacityReservation(cr), nil
}

// ListCapacityReservations returns capacity reservations, optionally filtered by work group.
func (s *AthenaService) ListCapacityReservations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	workGroup := request.GetStringParam(req.Parameters, "WorkGroup")

	reservations, err := stores.capacityReservationStore.ListCapacityReservations(workGroup)
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(reservations))
	for _, cr := range reservations {
		items = append(items, formatCapacityReservation(cr))
	}

	return map[string]interface{}{
		"CapacityReservations": items,
	}, nil
}

// UpdateCapacityReservation updates the specified capacity reservation.
func (s *AthenaService) UpdateCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Name is required")
	}

	cr, err := stores.capacityReservationStore.GetCapacityReservation(name)
	if err != nil {
		return nil, awserrors.NewResourceNotFoundException("CapacityReservation", name)
	}

	if dpus := request.GetIntParam(req.Parameters, "TargetDpus"); dpus > 0 {
		cr.Capacity = int32(dpus)
	}
	cr.LastModifiedTime = time.Now().UTC()

	if err := stores.capacityReservationStore.UpdateCapacityReservation(cr); err != nil {
		return nil, err
	}

	return formatCapacityReservation(cr), nil
}

// CancelCapacityReservation cancels the specified capacity reservation.
func (s *AthenaService) CancelCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Name is required")
	}

	cr, err := stores.capacityReservationStore.GetCapacityReservation(name)
	if err != nil {
		return nil, awserrors.NewResourceNotFoundException("CapacityReservation", name)
	}

	cr.Status = athenastore.CapacityReservationStatusCancelled
	cr.LastModifiedTime = time.Now().UTC()
	if err := stores.capacityReservationStore.UpdateCapacityReservation(cr); err != nil {
		return nil, err
	}

	return formatCapacityReservation(cr), nil
}
