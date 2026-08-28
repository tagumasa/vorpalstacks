package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

// CreateServiceLinkedRoleInput carries the parsed CreateServiceLinkedRole
// request.
type CreateServiceLinkedRoleInput struct {
	AWSServiceName string
	CustomSuffix   string
	Description    string
}

// createServiceLinkedRoleCore validates input and creates the
// service-linked role for the specified AWS service.
func (s *IAMService) createServiceLinkedRoleCore(store *iamstore.IAMStore, input *CreateServiceLinkedRoleInput) (*iamstore.Role, error) {
	if err := validateAWSServiceName(input.AWSServiceName); err != nil {
		return nil, err
	}

	// Smithy customSuffixType: length 1-64, pattern ^[\w+=,.@-]+$.
	// An empty suffix is allowed (no suffix appended to the role name).
	if input.CustomSuffix != "" && !validateCustomSuffix(input.CustomSuffix) {
		return nil, NewInvalidInputError("CustomSuffix", "must match pattern ^[\\w+=,.@-]+$ and be 1-64 characters")
	}

	// Service-linked role names combine a service-provided prefix with the
	// optional custom suffix; the documented prefix form is
	// AWSServiceRoleFor<Service> (for example ecs.amazonaws.com creates
	// AWSServiceRoleForECS). AWS defines the exact capitalisation per
	// service; this derivation uses the uppercase short name, which matches
	// the documented example. AWS enforces a 64-character maximum on IAM
	// role names, so reject any combination that would exceed the limit
	// before attempting to create the role.
	dotIdx := strings.Index(input.AWSServiceName, ".")
	roleName := "AWSServiceRoleFor" + strings.ToUpper(input.AWSServiceName[:dotIdx])
	if input.CustomSuffix != "" {
		roleName = roleName + "-" + input.CustomSuffix
	}
	if len(roleName) > 64 {
		return nil, NewInvalidInputError("RoleName", "derived role name exceeds the 64-character limit (service prefix + custom suffix)")
	}
	path := "/aws-service-role/" + input.AWSServiceName + "/"

	// Build the trust policy document that grants the service permission to
	// assume this role. validateAWSServiceName guarantees the dotted
	// <service>.amazonaws.com form required for a valid Service principal.
	trustPolicy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": input.AWSServiceName,
				},
				"Action": "sts:AssumeRole",
			},
		},
	}

	trustPolicyJSON, err := json.Marshal(trustPolicy)
	if err != nil {
		return nil, err
	}

	role, err := store.Roles().Create(roleName, path, s.accountID, string(trustPolicyJSON), input.Description, 3600, nil)
	if err != nil {
		if errors.Is(err, iamstore.ErrRoleAlreadyExists) {
			return nil, NewRoleAlreadyExistsError(roleName)
		}
		return nil, err
	}
	return role, nil
}

// deleteServiceLinkedRoleCore validates input and submits an asynchronous
// service-linked role deletion request, returning the deletion task
// identifier immediately.  The actual cleanup (inline policies, attached
// policies, instance profiles, role deletion) runs in a background
// goroutine.  The task status is persisted so that
// GetServiceLinkedRoleDeletionStatus can report results even after a
// server restart.
func (s *IAMService) deleteServiceLinkedRoleCore(store *iamstore.IAMStore, roleName string) (string, error) {
	if roleName == "" {
		return "", NewValidationError("RoleName")
	}

	if !store.Roles().Exists(roleName) {
		return "", NewNoSuchRoleError(roleName)
	}

	// Verify the role is actually a service-linked role. Service-linked roles
	// are created with a path under /aws-service-role/. Non-service-linked
	// roles must be deleted via DeleteRole instead.
	role, err := store.Roles().Get(roleName)
	if err != nil {
		return "", NewNoSuchRoleError(roleName)
	}
	if !strings.HasPrefix(role.Path, "/aws-service-role/") {
		return "", NewDeleteConflictError("Cannot delete role " + roleName + " with DeleteServiceLinkedRole. Use DeleteRole instead.")
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
		return "", err
	}

	// Spawn a goroutine to perform the cleanup asynchronously. The
	// WaitGroup is tracked so that shutdown can wait for in-flight
	// deletions to finish rather than aborting them mid-way.
	s.slRoleDeletionWg.Add(1)
	go s.executeSLRoleDeletion(store, taskID, roleName)

	return taskID, nil
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

// getServiceLinkedRoleDeletionStatusCore validates input and retrieves
// the persisted deletion task.
func (s *IAMService) getServiceLinkedRoleDeletionStatusCore(store *iamstore.IAMStore, taskID string) (*iamstore.SLRoleDeletionTask, error) {
	if taskID == "" {
		return nil, NewValidationError("DeletionTaskId")
	}
	if len(taskID) > MaxDeletionTaskIdLength {
		return nil, NewInvalidInputError("DeletionTaskId", fmt.Sprintf("must be 1 to %d characters", MaxDeletionTaskIdLength))
	}
	task, err := store.SLRoleDeletionTasks().Get(taskID)
	if err != nil {
		return nil, NewNoSuchEntityError("deletion task", taskID)
	}
	return task, nil
}
