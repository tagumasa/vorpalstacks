package scheduler

import (
	"context"

	"vorpalstacks/internal/common/iam"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
)

// CreateScheduleInput is the transport-agnostic input for creating a
// schedule. Both the HTTP API and the admin gRPC handler build this from
// their respective request formats and delegate to createScheduleCore.
type CreateScheduleInput struct {
	Spec         *ScheduleSpec
	ClientToken  string
	Region       string
	IAMValidator *iam.IAMValidator
}

// CreateScheduleResult holds the output of a successful schedule creation.
type CreateScheduleResult struct {
	ScheduleArn string
}

// UpdateScheduleInput is the transport-agnostic input for updating a
// schedule. UpdateSchedule is a PUT operation (full replacement).
type UpdateScheduleInput struct {
	Spec         *ScheduleSpec
	Region       string
	IAMValidator *iam.IAMValidator
}

// UpdateScheduleResult holds the output of a successful schedule update.
type UpdateScheduleResult struct {
	ScheduleArn string
}

// DeleteScheduleInput is the transport-agnostic input for deleting a schedule.
type DeleteScheduleInput struct {
	Name      string
	GroupName string
}

// GetScheduleInput is the transport-agnostic input for retrieving a schedule.
type GetScheduleInput struct {
	Name      string
	GroupName string
}

// ListSchedulesInput is the transport-agnostic input for listing schedules.
type ListSchedulesInput struct {
	GroupName  string
	NamePrefix string
	State      string
	MaxResults int32
	NextToken  string
}

// ListSchedulesResult holds the output of a successful schedule listing.
type ListSchedulesResult struct {
	Schedules []schedulerstore.ScheduleSummary
	NextToken string
}

// createScheduleCore is the single entry point for schedule creation shared
// by the HTTP API and the admin gRPC handler. It performs validation, IAM
// role validation, VPC validation, group existence check, ClientToken
// idempotency, and store creation.
func (s *SchedulerService) createScheduleCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *CreateScheduleInput) (*CreateScheduleResult, error) {
	validated, err := validateScheduleFields(in.Spec)
	if err != nil {
		return nil, err
	}

	target := in.Spec.Target

	if in.IAMValidator != nil && target != nil && target.RoleArn != "" {
		if err := in.IAMValidator.ValidateRoleForService(ctx, target.RoleArn, iam.ServicePrincipalScheduler); err != nil {
			return nil, err
		}
	}

	if err := s.validateVpcConfig(ctx, in.Region, target); err != nil {
		return nil, err
	}

	groupName := in.Spec.GroupName
	if groupName == "" {
		groupName = "default"
	}

	clientToken := in.ClientToken
	tokenClaimed := false
	if clientToken != "" {
		if err := validateClientToken(clientToken); err != nil {
			return nil, err
		}
		expectedArn := store.BuildScheduleARN(groupName, in.Spec.Name)
		if entry, created := store.ClientTokens().LookupOrClaim(clientToken, expectedArn, "schedule"); !created {
			return &CreateScheduleResult{ScheduleArn: entry.ResourceArn}, nil
		}
		tokenClaimed = true
	}

	if groupName != "default" {
		if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
			if tokenClaimed {
				store.ClientTokens().Release(clientToken)
			}
			if err == schedulerstore.ErrScheduleGroupNotFound {
				return nil, ErrScheduleGroupNotFound
			}
			return nil, ErrInternalServer
		}
	}

	schedule := &schedulerstore.Schedule{
		Name:                       in.Spec.Name,
		GroupName:                  groupName,
		ScheduleExpression:         in.Spec.ScheduleExpression,
		Target:                     target,
		FlexibleTimeWindow:         in.Spec.FlexibleTimeWindow,
		State:                      validated.State,
		ScheduleExpressionTimezone: in.Spec.ScheduleExpressionTimezone,
		Description:                in.Spec.Description,
		KmsKeyArn:                  in.Spec.KmsKeyArn,
		StartDate:                  validated.StartDate,
		EndDate:                    validated.EndDate,
		ActionAfterCompletion:      validated.ActionAfterCompletion,
	}

	if err := store.CreateSchedule(ctx, schedule); err != nil {
		if tokenClaimed {
			store.ClientTokens().Release(clientToken)
		}
		if err == schedulerstore.ErrScheduleAlreadyExists {
			return nil, ErrScheduleAlreadyExists
		}
		return nil, ErrInternalServer
	}

	return &CreateScheduleResult{ScheduleArn: schedule.ARN}, nil
}

// updateScheduleCore is the single entry point for schedule updates shared
// by the HTTP API and the admin gRPC handler. UpdateSchedule is a PUT
// operation: all fields from the request replace the existing values.
func (s *SchedulerService) updateScheduleCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *UpdateScheduleInput) (*UpdateScheduleResult, error) {
	existing, err := store.GetSchedule(ctx, in.Spec.GroupName, in.Spec.Name)
	if err != nil {
		if err == schedulerstore.ErrScheduleNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, ErrInternalServer
	}

	validated, err := validateScheduleFields(in.Spec)
	if err != nil {
		return nil, err
	}

	target := in.Spec.Target

	if in.IAMValidator != nil && target != nil && target.RoleArn != "" {
		if err := in.IAMValidator.ValidateRoleForService(ctx, target.RoleArn, iam.ServicePrincipalScheduler); err != nil {
			return nil, err
		}
	}

	if err := s.validateVpcConfig(ctx, in.Region, target); err != nil {
		return nil, err
	}

	existing.ScheduleExpression = in.Spec.ScheduleExpression
	existing.Target = target
	existing.FlexibleTimeWindow = in.Spec.FlexibleTimeWindow
	existing.Description = in.Spec.Description
	existing.ScheduleExpressionTimezone = in.Spec.ScheduleExpressionTimezone
	existing.KmsKeyArn = in.Spec.KmsKeyArn
	existing.State = validated.State
	existing.ActionAfterCompletion = validated.ActionAfterCompletion
	existing.StartDate = validated.StartDate
	existing.EndDate = validated.EndDate

	if err := store.UpdateSchedule(ctx, existing); err != nil {
		return nil, ErrInternalServer
	}

	return &UpdateScheduleResult{ScheduleArn: existing.ARN}, nil
}

// deleteScheduleCore is the single entry point for schedule deletion.
func (s *SchedulerService) deleteScheduleCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *DeleteScheduleInput) error {
	if err := store.DeleteSchedule(ctx, in.GroupName, in.Name); err != nil {
		if err == schedulerstore.ErrScheduleNotFound {
			return ErrScheduleNotFound
		}
		return ErrInternalServer
	}
	return nil
}

// getScheduleCore is the single entry point for retrieving a schedule.
func (s *SchedulerService) getScheduleCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *GetScheduleInput) (*schedulerstore.Schedule, error) {
	schedule, err := store.GetSchedule(ctx, in.GroupName, in.Name)
	if err != nil {
		if err == schedulerstore.ErrScheduleNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, ErrInternalServer
	}
	return schedule, nil
}

// listSchedulesCore is the single entry point for listing schedules.
func (s *SchedulerService) listSchedulesCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *ListSchedulesInput) (*ListSchedulesResult, error) {
	result, err := store.ListSchedules(ctx, in.GroupName, in.NamePrefix, schedulerstore.ScheduleState(in.State), in.MaxResults, in.NextToken)
	if err != nil {
		return nil, ErrInternalServer
	}
	return &ListSchedulesResult{
		Schedules: result.Schedules,
		NextToken: result.NextToken,
	}, nil
}
