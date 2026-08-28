package iam

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
)

// CreateServiceLinkedRole creates a service-linked role for a specified AWS service.
func (s *IAMService) CreateServiceLinkedRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	input := &CreateServiceLinkedRoleInput{
		AWSServiceName: request.GetStringParam(req.Parameters, "AWSServiceName"),
		CustomSuffix:   request.GetStringParam(req.Parameters, "CustomSuffix"),
		Description:    request.GetStringParam(req.Parameters, "Description"),
	}
	role, err := s.createServiceLinkedRoleCore(store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"Role": roleToResponse(role),
	}, nil
}

// DeleteServiceLinkedRole submits an asynchronous service-linked role
// deletion request and returns the deletion task identifier immediately.
// The actual cleanup (inline policies, attached policies, instance
// profiles, role deletion) runs in a background goroutine.  The task
// status is persisted so that GetServiceLinkedRoleDeletionStatus can
// report results even after a server restart.
func (s *IAMService) DeleteServiceLinkedRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	taskID, err := s.deleteServiceLinkedRoleCore(store, request.GetStringParam(req.Parameters, "RoleName"))
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"DeletionTaskId": taskID,
	}, nil
}

// GetServiceLinkedRoleDeletionStatus retrieves the status of a
// service-linked role deletion task from persistent storage.
func (s *IAMService) GetServiceLinkedRoleDeletionStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	task, err := s.getServiceLinkedRoleDeletionStatusCore(store, request.GetStringParam(req.Parameters, "DeletionTaskId"))
	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"Status": task.Status,
	}
	if task.DeletionFailed {
		resp["Reason"] = map[string]interface{}{
			"ReasonCode":    "DeleteConflict",
			"ReasonMessage": task.ErrorReason,
		}
	}

	return resp, nil
}

// RecoverOrphanedSLRoleDeletionTasks marks any IN_PROGRESS tasks as
// FAILED at startup.  Called during IAM service initialisation to handle
// tasks that were interrupted by a server restart.
//
// AWS does not automatically resume SL role deletion after a server-side
// interruption; the operator is expected to retry DeleteServiceLinkedRole
// manually after resolving any conflicts. Marking the task FAILED (rather
// than silently leaving it IN_PROGRESS) makes the failure observable via
// GetServiceLinkedRoleDeletionStatus and matches the AWS behaviour where
// an interrupted deletion surfaces a "Deletion in progress exceeded
// maximum time" style error.
func (s *IAMService) RecoverOrphanedSLRoleDeletionTasks() {
	store, err := s.GetStoreForRegion("global")
	if err != nil {
		logs.Error("RecoverOrphanedSLRoleDeletionTasks: failed to get global store",
			logs.Any("error", err))
		return
	}
	tasks, err := store.SLRoleDeletionTasks().List()
	if err != nil {
		logs.Error("RecoverOrphanedSLRoleDeletionTasks: failed to list tasks",
			logs.Any("error", err))
		return
	}
	for _, task := range tasks {
		if task.Status == "IN_PROGRESS" {
			task.Status = "FAILED"
			task.DeletionFailed = true
			task.ErrorReason = "Server restarted during deletion"
			if err := store.SLRoleDeletionTasks().Put(task); err != nil {
				logs.Error("RecoverOrphanedSLRoleDeletionTasks: failed to mark task FAILED",
					logs.Any("taskID", task.TaskID),
					logs.Any("error", err))
			}
		}
	}
}
