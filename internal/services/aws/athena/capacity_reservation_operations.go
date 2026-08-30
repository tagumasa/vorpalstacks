package athena

import (
	"context"

	"vorpalstacks/internal/common/pagination"
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
	input := CreateCapacityReservationInput{
		Name:            request.GetStringParam(req.Parameters, "Name"),
		TargetWorkGroup: request.GetStringParam(req.Parameters, "TargetWorkGroup"),
		TargetDpus:      int32(request.GetIntParam(req.Parameters, "TargetDpus")),
	}

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	created, err := createCapacityReservationCore(stores, input)
	if err != nil {
		return nil, err
	}

	return formatCapacityReservation(created), nil
}

// GetCapacityReservation retrieves details about the specified capacity reservation.
func (s *AthenaService) GetCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cr, err := getCapacityReservationCore(stores, name)
	if err != nil {
		return nil, err
	}

	return formatCapacityReservation(cr), nil
}

// ListCapacityReservations returns capacity reservations, optionally filtered by work group.
func (s *AthenaService) ListCapacityReservations(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	maxResults, hasMaxResults := request.GetIntParamCaseInsensitive(req.Parameters, "MaxResults")
	input := ListCapacityReservationsInput{
		WorkGroup:     request.GetStringParam(req.Parameters, "WorkGroup"),
		MaxResults:    maxResults,
		HasMaxResults: hasMaxResults,
		NextToken:     pagination.GetMarker(req.Parameters, "NextToken"),
	}

	reservations, nextMarker, err := s.listCapacityReservationsCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	summaries := make([]map[string]interface{}, 0, len(reservations))
	for _, cr := range reservations {
		summaries = append(summaries, formatCapacityReservation(cr))
	}

	return pagination.BuildListResponse("CapacityReservations", summaries, nextMarker), nil
}

// UpdateCapacityReservation updates the specified capacity reservation.
// Per the Smithy model, both Name and TargetDpus are REQUIRED on this operation.
func (s *AthenaService) UpdateCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := UpdateCapacityReservationInput{
		Name:       request.GetStringParam(req.Parameters, "Name"),
		TargetDpus: int32(request.GetIntParam(req.Parameters, "TargetDpus")),
	}

	cr, err := s.updateCapacityReservationCore(reqCtx, input)
	if err != nil {
		return nil, err
	}

	return formatCapacityReservation(cr), nil
}

// CancelCapacityReservation cancels the specified capacity reservation.
func (s *AthenaService) CancelCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	cr, err := cancelCapacityReservationCore(stores, name)
	if err != nil {
		return nil, err
	}

	return formatCapacityReservation(cr), nil
}

// DeleteCapacityReservation deletes the specified capacity reservation.
func (s *AthenaService) DeleteCapacityReservation(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "Name")

	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	if err := deleteCapacityReservationCore(stores, name); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// GetCapacityAssignmentConfiguration returns the capacity assignment
// configuration for a capacity reservation, if one exists.
func (s *AthenaService) GetCapacityAssignmentConfiguration(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	name := request.GetStringParam(req.Parameters, "CapacityReservationName")

	assignments, err := s.getCapacityAssignmentConfigurationCore(reqCtx, name)
	if err != nil {
		return nil, err
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

	input := PutCapacityAssignmentConfigurationInput{
		CapacityReservationName: request.GetStringParam(req.Parameters, "CapacityReservationName"),
		CapacityAssignments:     assignments,
	}

	if err := s.putCapacityAssignmentConfigurationCore(reqCtx, input); err != nil {
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
