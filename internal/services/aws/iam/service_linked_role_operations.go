package iam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/services/aws/common/request"
	iamstore "vorpalstacks/internal/store/aws/iam"
)

var (
	slRoleDeletionTasks sync.Map
	slRoleCleanupOnce   sync.Once
)

func init() {
	slRoleCleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(30 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				now := time.Now()
				slRoleDeletionTasks.Range(func(key, value interface{}) bool {
					task := value.(*serviceLinkedRoleDeletionTask)
					if task.CreatedAt.Add(1 * time.Hour).Before(now) {
						slRoleDeletionTasks.Delete(key)
					}
					return true
				})
			}
		}()
	})
}

type serviceLinkedRoleDeletionTask struct {
	TaskID         string    `json:"taskId"`
	RoleName       string    `json:"roleName"`
	Status         string    `json:"status"`
	DeletionFailed bool      `json:"deletionFailed"`
	ErrorReason    string    `json:"errorReason,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
}

// CreateServiceLinkedRole creates a service-linked role for a specified AWS service.
func (s *IAMService) CreateServiceLinkedRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	awsServiceName := request.GetStringParam(req.Parameters, "AWSServiceName")
	if awsServiceName == "" {
		return nil, NewInvalidInputError("AWSServiceName", "cannot be empty")
	}

	description := request.GetStringParam(req.Parameters, "Description")
	customSuffix := request.GetStringParam(req.Parameters, "CustomSuffix")

	roleName := fmt.Sprintf("AWSServiceRoleFor%s", awsServiceName)
	if customSuffix != "" {
		roleName = fmt.Sprintf("%s_%s", roleName, customSuffix)
	}

	principalService := awsServiceName
	if !strings.HasSuffix(principalService, ".amazonaws.com") {
		principalService = awsServiceName + ".amazonaws.com"
	}

	assumeRolePolicy := map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect":    "Allow",
				"Principal": map[string]interface{}{"Service": []string{principalService}},
				"Action":    []string{"sts:AssumeRole"},
			},
		},
	}
	policyDoc, err := json.Marshal(assumeRolePolicy)
	if err != nil {
		return nil, err
	}

	path := "/aws-service-role/" + awsServiceName + "/"

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	role, err := store.Roles().Create(roleName, path, s.accountID, string(policyDoc), description, 3600, nil)
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

// DeleteServiceLinkedRole submits a service-linked role deletion request and returns the deletion task identifier.
func (s *IAMService) DeleteServiceLinkedRole(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	roleName := request.GetStringParam(req.Parameters, "RoleName")
	if roleName == "" {
		return nil, ErrNoSuchRole
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if !store.Roles().Exists(roleName) {
		return nil, NewNoSuchRoleError(roleName)
	}

	taskID := fmt.Sprintf("task-%d", time.Now().UnixNano())

	task := &serviceLinkedRoleDeletionTask{
		TaskID:    taskID,
		RoleName:  roleName,
		Status:    "SUCCEEDED",
		CreatedAt: time.Now(),
	}

	if err := store.InlinePolicies().DeleteAllForPrincipal(PrincipalTypeRole, roleName); err != nil {
		task.Status = "FAILED"
		task.DeletionFailed = true
		task.ErrorReason = "Failed to delete inline policies"
	}

	if task.Status == "SUCCEEDED" {
		attachedPolicies, err := store.AttachedPolicies().ListAttachedPolicies(PrincipalTypeRole, roleName)
		if err != nil {
			task.Status = "FAILED"
			task.DeletionFailed = true
			task.ErrorReason = "Failed to list attached policies"
		}
		for _, policyArn := range attachedPolicies {
			if err := store.AttachedPolicies().Detach(PrincipalTypeRole, roleName, policyArn); err != nil {
				task.Status = "FAILED"
				task.DeletionFailed = true
				task.ErrorReason = "Failed to detach managed policies"
				break
			}
			if err := store.Policies().DecrementAttachmentCount(policyArn); err != nil {
				task.Status = "FAILED"
				task.DeletionFailed = true
				task.ErrorReason = "Failed to update policy attachment count"
				break
			}
		}
	}

	if task.Status == "SUCCEEDED" {
		instanceProfiles, err := store.InstanceProfiles().ListForRole(roleName, "", 1)
		if err != nil {
			task.Status = "FAILED"
			task.DeletionFailed = true
			task.ErrorReason = "Failed to list instance profiles"
		}
		for _, ip := range instanceProfiles.InstanceProfiles {
			if err := store.InstanceProfiles().RemoveRole(ip.InstanceProfileName, roleName); err != nil {
				task.Status = "FAILED"
				task.DeletionFailed = true
				task.ErrorReason = "Failed to remove role from instance profiles"
				break
			}
		}
	}

	if task.Status == "SUCCEEDED" {
		if err := store.Roles().Delete(roleName); err != nil {
			task.Status = "FAILED"
			task.DeletionFailed = true
			task.ErrorReason = "Failed to delete role"
		}
	}

	slRoleDeletionTasks.Store(taskID, task)

	return map[string]interface{}{
		"DeletionTaskId": taskID,
	}, nil
}

// GetServiceLinkedRoleDeletionStatus retrieves the status of the service-linked role deletion task.
func (s *IAMService) GetServiceLinkedRoleDeletionStatus(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	taskID := request.GetStringParam(req.Parameters, "DeletionTaskId")
	if taskID == "" {
		return nil, NewValidationError("DeletionTaskId")
	}

	val, ok := slRoleDeletionTasks.Load(taskID)
	if !ok {
		return map[string]interface{}{
			"Status": "FAILED",
		}, nil
	}

	task := val.(*serviceLinkedRoleDeletionTask)

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
