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
		return nil, capacityReservationNotFound(name)
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
		return nil, capacityReservationNotFound(name)
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
		return nil, capacityReservationNotFound(name)
	}

	cr.Status = athenastore.CapacityReservationStatusCancelled
	cr.LastModifiedTime = time.Now().UTC()
	if err := stores.capacityReservationStore.UpdateCapacityReservation(cr); err != nil {
		return nil, err
	}

	return formatCapacityReservation(cr), nil
}

// DeleteCapacityReservation deletes the specified capacity reservation.
func (s *AthenaService) DeleteCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewValidationException("Name is required")
	}

	if _, err := stores.capacityReservationStore.GetCapacityReservation(name); err != nil {
		return nil, capacityReservationNotFound(name)
	}

	if err := stores.capacityReservationStore.DeleteCapacityReservation(name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// GetCapacityAssignmentConfiguration returns the capacity assignment
// configuration for a capacity reservation, if one exists.
func (s *AthenaService) GetCapacityAssignmentConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "CapacityReservationName")
	if name == "" {
		return nil, awserrors.NewValidationException("CapacityReservationName is required")
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	assignments, err := stores.capacityReservationStore.GetCapacityAssignments(name)
	if err != nil {
		return nil, capacityReservationNotFound(name)
	}

	var assignmentList []map[string]interface{}
	for _, wgNames := range assignments {
		assignmentList = append(assignmentList, map[string]interface{}{
			"WorkGroupNames": wgNames,
		})
	}

	return map[string]interface{}{
		"CapacityAssignmentConfiguration": map[string]interface{}{
			"CapacityReservationName": name,
			"CapacityAssignments":     assignmentList,
		},
	}, nil
}

// PutCapacityAssignmentConfiguration puts a new capacity assignment
// configuration for a specified capacity reservation. If a configuration
// already exists, it is replaced.
func (s *AthenaService) PutCapacityAssignmentConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "CapacityReservationName")
	if name == "" {
		return nil, awserrors.NewValidationException("CapacityReservationName is required")
	}

	assignmentsRaw := request.GetArrayParam(req.Parameters, "CapacityAssignments")
	var assignments [][]string
	for _, a := range assignmentsRaw {
		if m, ok := a.(map[string]interface{}); ok {
			wgNamesRaw, _ := m["WorkGroupNames"].([]interface{})
			var wgNames []string
			for _, w := range wgNamesRaw {
				if ws, ok := w.(string); ok {
					wgNames = append(wgNames, ws)
				}
			}
			assignments = append(assignments, wgNames)
		}
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := stores.capacityReservationStore.PutCapacityAssignments(name, assignments); err != nil {
		if err == athenastore.ErrCapacityReservationNotFound {
			return nil, capacityReservationNotFound(name)
		}
		return nil, err
	}

	// Smithy defines PutCapacityAssignmentConfigurationOutput as an empty structure.
	return map[string]interface{}{}, nil
}

// ListApplicationDPUSizes returns the supported DPU sizes for the
// supported Athena application runtimes.
func (s *AthenaService) ListApplicationDPUSizes(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	return map[string]interface{}{
		"ApplicationDPUSizes": []map[string]interface{}{
			{
				"ApplicationRuntimeId": "Athena notebook version 1",
				"SupportedDPUSizes":    []int{24, 48, 72},
			},
		},
	}, nil
}
