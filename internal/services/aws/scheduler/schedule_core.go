package scheduler

import (
	"context"
	"time"

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
	// ClientToken preserves the wire member's idempotency semantics
	// (@length(1,64), @pattern): an invalid token is rejected, and a
	// replayed token returns the first application's ARN without
	// re-applying the full-override mutation.
	ClientToken string
}

// UpdateScheduleResult holds the output of a successful schedule update.
type UpdateScheduleResult struct {
	ScheduleArn string
}

// DeleteScheduleInput is the transport-agnostic input for deleting a schedule.
type DeleteScheduleInput struct {
	Name      string
	GroupName string
	// ClientToken preserves the wire member's idempotency semantics
	// (@length(1,64), @pattern): an invalid token is rejected, and a
	// replayed token reports the deletion as already applied instead of
	// surfacing not-found for the removed schedule.
	ClientToken string
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
	// MaxResults is nil when the member was absent from the request; the
	// Core applies the default page size in that case.
	MaxResults *int32
	NextToken  string
}

// ListSchedulesResult holds the output of a successful schedule listing.
type ListSchedulesResult struct {
	Schedules []schedulerstore.ScheduleSummary
	NextToken string
}

// resolveScheduleGroup applies the documented default-group semantics:
// "If you omit this value, EventBridge Scheduler assumes the group is
// associated to the default group." (GroupName, UpdateSchedule API
// reference), and validates the ScheduleGroupName shape.
func resolveScheduleGroup(groupName string) (string, error) {
	if groupName == "" {
		return "default", nil
	}
	if err := validateScheduleGroupName(groupName); err != nil {
		return "", err
	}
	return groupName, nil
}

// resolveScheduleIdentifier validates the schedule identifier pair: the
// Name shape (@length(1,64), @pattern) plus the group resolution above.
// Every read/update/delete path resolves through here so a malformed
// identifier is a ValidationException, never a resource lookup.
func resolveScheduleIdentifier(name, groupName string) (string, error) {
	if name == "" || !namePattern.MatchString(name) {
		return "", ErrValidation
	}
	return resolveScheduleGroup(groupName)
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

	if err := s.validateKmsKey(ctx, in.Region, in.Spec.KmsKeyArn); err != nil {
		return nil, err
	}

	groupName, err := resolveScheduleGroup(in.Spec.GroupName)
	if err != nil {
		return nil, err
	}

	clientToken := in.ClientToken
	tokenClaimed := false
	var expectedArn string
	if clientToken != "" {
		if err := validateClientToken(clientToken); err != nil {
			return nil, err
		}
		expectedArn = store.BuildScheduleARN(groupName, in.Spec.Name)
		if entry, created := store.ClientTokens().LookupOrClaim(clientToken, expectedArn, "schedule"); !created {
			return &CreateScheduleResult{ScheduleArn: entry.ResourceArn}, nil
		}
		tokenClaimed = true
	}

	if groupName != "default" {
		if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
			if tokenClaimed {
				store.ClientTokens().Release(clientToken, expectedArn, "schedule")
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
			store.ClientTokens().Release(clientToken, expectedArn, "schedule")
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
	// Malformed identifiers are ValidationException on AWS, so the pair is
	// resolved and validated before the existence probe; for a well-formed
	// identifier a missing schedule still reports not-found ahead of the
	// remaining input validation (historical error precedence).
	groupName, err := resolveScheduleIdentifier(in.Spec.Name, in.Spec.GroupName)
	if err != nil {
		return nil, err
	}
	if _, err := store.GetSchedule(ctx, groupName, in.Spec.Name); err != nil {
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

	if err := s.validateKmsKey(ctx, in.Region, in.Spec.KmsKeyArn); err != nil {
		return nil, err
	}

	// The idempotency token is claimed only after every validation passes;
	// a replay returns the first application's ARN without re-applying the
	// full-override mutation (so a replay does not re-stamp the record's
	// modification date).
	clientToken := in.ClientToken
	tokenClaimed := false
	var expectedArn string
	if clientToken != "" {
		if err := validateClientToken(clientToken); err != nil {
			return nil, err
		}
		expectedArn = store.BuildScheduleARN(groupName, in.Spec.Name)
		if entry, created := store.ClientTokens().LookupOrClaim(clientToken, expectedArn, "schedule-update"); !created {
			return &UpdateScheduleResult{ScheduleArn: entry.ResourceArn}, nil
		}
		tokenClaimed = true
	}

	// The user's fields are applied through the store-level atomic mutation
	// so a concurrent engine write (completion or firing markers) can never
	// be lost to this update's read-modify-write cycle.
	var scheduleARN string
	err = store.MutateSchedule(ctx, groupName, in.Spec.Name, func(existing *schedulerstore.Schedule) error {
		// Captured before the assignments: re-lifecycling a completed
		// schedule and changing the expression both start a new firing
		// lifecycle and must reset the delivered-boundary marker.
		completionWasSet := existing.CompletionDate != nil
		exprChanged := existing.ScheduleExpression != in.Spec.ScheduleExpression
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
		// Updating a schedule starts a new execution lifecycle: a one-time
		// schedule whose previous lifecycle had ended becomes eligible to
		// fire again for its new expression or schedule time.
		existing.CompletionDate = nil
		if completionWasSet || exprChanged {
			// AWS UpdateSchedule is a full override that resets execution
			// state ("uses all values, including empty values, specified in
			// the request and overrides the existing schedule"), so the
			// previous lifecycle's delivered boundary must not suppress
			// the new lifecycle — not before and not after a restart.
			existing.LastFiredAt = nil
		}
		existing.LastModificationDate = time.Now().UTC()
		scheduleARN = existing.ARN
		return nil
	})
	if err != nil {
		if tokenClaimed {
			store.ClientTokens().Release(clientToken, expectedArn, "schedule-update")
		}
		if err == schedulerstore.ErrScheduleNotFound {
			return nil, ErrScheduleNotFound
		}
		return nil, ErrInternalServer
	}

	return &UpdateScheduleResult{ScheduleArn: scheduleARN}, nil
}

// deleteScheduleCore is the single entry point for schedule deletion.
func (s *SchedulerService) deleteScheduleCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *DeleteScheduleInput) error {
	groupName, err := resolveScheduleIdentifier(in.Name, in.GroupName)
	if err != nil {
		return err
	}
	// A replayed idempotency token reports the first deletion's outcome;
	// an unrecoverable failure releases the claim so a retry re-executes.
	clientToken := in.ClientToken
	tokenClaimed := false
	var expectedArn string
	if clientToken != "" {
		if err := validateClientToken(clientToken); err != nil {
			return err
		}
		expectedArn = store.BuildScheduleARN(groupName, in.Name)
		if _, created := store.ClientTokens().LookupOrClaim(clientToken, expectedArn, "schedule-delete"); !created {
			return nil
		}
		tokenClaimed = true
	}
	if err := store.DeleteSchedule(ctx, groupName, in.Name); err != nil {
		if tokenClaimed {
			store.ClientTokens().Release(clientToken, expectedArn, "schedule-delete")
		}
		if err == schedulerstore.ErrScheduleNotFound {
			return ErrScheduleNotFound
		}
		return ErrInternalServer
	}
	return nil
}

// getScheduleCore is the single entry point for retrieving a schedule.
func (s *SchedulerService) getScheduleCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *GetScheduleInput) (*schedulerstore.Schedule, error) {
	groupName, err := resolveScheduleIdentifier(in.Name, in.GroupName)
	if err != nil {
		return nil, err
	}
	schedule, err := store.GetSchedule(ctx, groupName, in.Name)
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
	if err := validateListStateFilter(in.State); err != nil {
		return nil, err
	}
	if err := validateListNamePrefix(in.NamePrefix); err != nil {
		return nil, err
	}
	if err := validateNextToken(in.NextToken); err != nil {
		return nil, err
	}
	maxResults, err := resolveListMaxResults(in.MaxResults)
	if err != nil {
		return nil, err
	}
	// A scoped listing must reference an existing group: the model defines
	// ResourceNotFoundException on ListSchedules — the group-filtered
	// operation — while the unfiltered ListScheduleGroups defines no such
	// error, so the group filter is the resource this operation references.
	if in.GroupName != "" {
		if err := validateScheduleGroupName(in.GroupName); err != nil {
			return nil, err
		}
		if _, err := store.GetScheduleGroup(ctx, in.GroupName); err != nil {
			if err == schedulerstore.ErrScheduleGroupNotFound {
				return nil, ErrScheduleGroupNotFound
			}
			return nil, ErrInternalServer
		}
	}

	result, err := store.ListSchedules(ctx, in.GroupName, in.NamePrefix, schedulerstore.ScheduleState(in.State), maxResults, in.NextToken)
	if err != nil {
		return nil, ErrInternalServer
	}
	return &ListSchedulesResult{
		Schedules: result.Schedules,
		NextToken: result.NextToken,
	}, nil
}
