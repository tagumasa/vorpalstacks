package iam

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateServiceLinkedRole creates a service-linked role for a specified AWS service.
func (s *IAMService) CreateServiceLinkedRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	awsServiceName := request.GetStringParam(req.Parameters, "AWSServiceName")
	if err := validateAWSServiceName(awsServiceName); err != nil {
		return nil, err
	}

	customSuffix := request.GetStringParam(req.Parameters, "CustomSuffix")
	description := request.GetStringParam(req.Parameters, "Description")

	// Smithy customSuffixType: length 1-64, pattern ^[\w+=,.@-]+$.
	// An empty suffix is allowed (no suffix appended to the role name).
	if customSuffix != "" && !validateCustomSuffix(customSuffix) {
		return nil, NewInvalidInputError("CustomSuffix", "must match pattern ^[\\w+=,.@-]+$ and be 1-64 characters")
	}

	// Derive the role name from the service's short name (the prefix
	// before the first dot) plus an optional custom suffix. validateAWSServiceName
	// guarantees the dotted <service>.amazonaws.com form so the prefix is never
	// empty. AWS enforces a 64-character maximum on IAM role names, so reject
	// any combination that would exceed the limit before attempting to create
	// the role.
	dotIdx := strings.Index(awsServiceName, ".")
	roleName := awsServiceName[:dotIdx]
	if customSuffix != "" {
		roleName = roleName + "-" + customSuffix
	}
	if len(roleName) > 64 {
		return nil, NewInvalidInputError("RoleName", "derived role name exceeds the 64-character limit (service prefix + custom suffix)")
	}
	path := "/aws-service-role/" + awsServiceName + "/"

	// Build the trust policy document that grants the service permission to
	// assume this role. validateAWSServiceName guarantees the dotted
	// <service>.amazonaws.com form required for a valid Service principal.
	trustPolicy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": awsServiceName,
				},
				"Action": "sts:AssumeRole",
			},
		},
	}

	trustPolicyJSON, err := json.Marshal(trustPolicy)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	role, err := store.Roles().Create(roleName, path, s.accountID, string(trustPolicyJSON), description, 3600, nil)
	if err != nil {
		if errors.Is(err, iamstore.ErrRoleAlreadyExists) {
			return nil, NewRoleAlreadyExistsError(roleName)
		}
		return nil, err
	}

	return map[string]interface{}{
		"Role": s.roleToResponse(reqCtx, role),
	}, nil
}

// DeleteServiceLinkedRole submits an asynchronous service-linked role
// deletion request and returns the deletion task identifier immediately.
// The actual cleanup (inline policies, attached policies, instance
// profiles, role deletion) runs in a background goroutine.  The task
// status is persisted so that GetServiceLinkedRoleDeletionStatus can
// report results even after a server restart.
func (s *IAMService) DeleteServiceLinkedRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, NewValidationError("RoleName")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	// Verify the role is actually a service-linked role. Service-linked roles
	// are created with a path under /aws-service-role/. Non-service-linked
	// roles must be deleted via DeleteRole instead.
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return nil, NewNoSuchRoleError(roleName)
	}
	if !strings.HasPrefix(role.Path, "/aws-service-role/") {
		return nil, NewDeleteConflictError("Cannot delete role " + roleName + " with DeleteServiceLinkedRole. Use DeleteRole instead.")
	}

	taskID := "task-" + uuid.New().String()

	// Persist the task with IN_PROGRESS status before returning so the
	// caller can poll GetServiceLinkedRoleDeletionStatus immediately.
	task := &iamstore.SLRoleDeletionTask{
		TaskID:    taskID,
		RoleName:  roleName,
		Status:    "IN_PROGRESS",
		CreatedAt: time.Now().UTC(),
	}
	if err := store.SLRoleDeletionTasks().Put(task); err != nil {
		return nil, err
	}

	// Spawn a goroutine to perform the cleanup asynchronously. The
	// WaitGroup is tracked so that shutdown can wait for in-flight
	// deletions to finish rather than aborting them mid-way.
	s.slRoleDeletionWg.Add(1)
	go s.executeSLRoleDeletion(store, taskID, roleName)

	return map[string]interface{}{
		"DeletionTaskId": taskID,
	}, nil
}

// executeSLRoleDeletion performs the actual cleanup of a service-linked
// role in the background. It delegates to cascadeDeleteRole (the same
// helper used by the admin DeleteRole handler) so that inline policies,
// attached managed policies, instance profile associations, and the role
// record itself are removed consistently and with proper error
// propagation (DecrementAttachmentCount errors are not silently dropped).
//
// The persisted task is updated via Get→mutate→Put so that the original
// CreatedAt timestamp is preserved.
func (s *IAMService) executeSLRoleDeletion(store *iamstore.IAMStore, taskID, roleName string) {
	defer s.slRoleDeletionWg.Done()
	defer func() {
		if r := recover(); r != nil {
			logs.Error("PANIC in service-linked role deletion goroutine",
				logs.Any("taskID", taskID),
				logs.Any("roleName", roleName),
				logs.Any("panic", r))
			s.updateSLRoleDeletionTask(store, taskID, "FAILED", true, "internal error during deletion")
		}
	}()

	// cascadeDeleteRole performs the full multi-step cleanup with proper
	// error propagation: inline policies → attached policies (with
	// DecrementAttachmentCount) → instance profile associations → role
	// record deletion. It does not take a lock because AWS itself does not
	// guarantee atomicity for async SL role deletion; concurrent
	// AttachRolePolicy calls during the cleanup window are a known
	// limitation and result in an error returned to the caller via the
	// task status.
	err := cascadeDeleteRole(store, roleName)
	if err != nil {
		s.updateSLRoleDeletionTask(store, taskID, "FAILED", true, err.Error())
		return
	}

	s.updateSLRoleDeletionTask(store, taskID, "SUCCEEDED", false, "")
}

// updateSLRoleDeletionTask fetches the existing persisted task, applies
// the new status fields, and saves it back. Reading the existing record
// preserves the original CreatedAt timestamp and any other metadata set
// at task creation time. Errors are logged but not propagated: the
// goroutine cannot return them to the caller, and the next
// GetServiceLinkedRoleDeletionStatus poll will simply observe the
// previously persisted status.
func (s *IAMService) updateSLRoleDeletionTask(store *iamstore.IAMStore, taskID, status string, failed bool, reason string) {
	existing, err := store.SLRoleDeletionTasks().Get(taskID)
	if err != nil {
		logs.Error("Failed to fetch SL role deletion task for update",
			logs.Any("taskID", taskID),
			logs.Any("error", err))
		return
	}
	existing.Status = status
	existing.DeletionFailed = failed
	existing.ErrorReason = reason
	if err := store.SLRoleDeletionTasks().Put(existing); err != nil {
		logs.Error("Failed to persist SL role deletion task update",
			logs.Any("taskID", taskID),
			logs.Any("error", err))
	}
}

// GetServiceLinkedRoleDeletionStatus retrieves the status of a
// service-linked role deletion task from persistent storage.
func (s *IAMService) GetServiceLinkedRoleDeletionStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetStringParam(req.Parameters, "DeletionTaskId")
	if taskID == "" {
		return nil, NewValidationError("DeletionTaskId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	task, err := store.SLRoleDeletionTasks().Get(taskID)
	if err != nil {
		return nil, NewNoSuchEntityError("deletion task", taskID)
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
