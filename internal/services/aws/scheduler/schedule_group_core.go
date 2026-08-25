package scheduler

import (
	"context"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	schedulerstore "vorpalstacks/internal/store/aws/scheduler"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// Schedule-group operation DTOs and Core functions: the single
// validation/persistence path for the schedule-group CRUD and the tag trio.
// The HTTP handlers are thin transport adapters (parse → DTO → Core →
// serialise); they must not import the store package directly.

// CreateScheduleGroupInput carries the CreateScheduleGroup request.
type CreateScheduleGroupInput struct {
	Name        string
	Tags        []tagutil.Tag
	ClientToken string
}

// CreateScheduleGroupResult carries the CreateScheduleGroup response.
type CreateScheduleGroupResult struct {
	ScheduleGroupArn string
}

// DeleteScheduleGroupInput carries the DeleteScheduleGroup request. The
// wire member ClientToken is not carried here: the operation's
// idempotency is realised by the DELETING-state handling, so a repeated
// delete with the same token returns success without a second sweep.
type DeleteScheduleGroupInput struct {
	Name string
}

// GetScheduleGroupInput carries the GetScheduleGroup request.
type GetScheduleGroupInput struct {
	Name string
}

// GetScheduleGroupResult carries the GetScheduleGroup response.
type GetScheduleGroupResult struct {
	Arn                  string
	Name                 string
	State                string
	CreationDate         time.Time
	LastModificationDate time.Time
}

// ScheduleGroupSummary is one list entry of ListScheduleGroups.
type ScheduleGroupSummary struct {
	Arn                  string
	Name                 string
	State                string
	CreationDate         *time.Time
	LastModificationDate *time.Time
}

// ListScheduleGroupsInput carries the ListScheduleGroups request. A zero
// MaxResults means the parameter was absent and is defaulted by the Core.
type ListScheduleGroupsInput struct {
	NamePrefix string
	MaxResults int32
	NextToken  string
}

// ListScheduleGroupsResult carries the ListScheduleGroups response.
type ListScheduleGroupsResult struct {
	ScheduleGroups []ScheduleGroupSummary
	NextToken      string
}

// createScheduleGroupCore validates and creates a schedule group, applies
// its tags, and honours ClientToken idempotency. A tagging failure rolls
// the group back so no orphan resource remains.
func (s *SchedulerService) createScheduleGroupCore(ctx context.Context, reqCtx *request.RequestContext, in *CreateScheduleGroupInput) (*CreateScheduleGroupResult, error) {
	if in.Name == "" || !namePattern.MatchString(in.Name) {
		return nil, ErrValidation
	}
	// Validate tags before creating the group so invalid tag sets are
	// rejected without creating an orphan resource.
	if err := ValidateScheduleGroupTags(in.Tags); err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	group := &schedulerstore.ScheduleGroup{Name: in.Name}

	tokenClaimed := false
	if in.ClientToken != "" {
		if err := validateClientToken(in.ClientToken); err != nil {
			return nil, err
		}
		expectedArn := store.BuildScheduleGroupARN(in.Name)
		if entry, created := store.ClientTokens().LookupOrClaim(in.ClientToken, expectedArn, "schedule-group"); !created {
			return &CreateScheduleGroupResult{ScheduleGroupArn: entry.ResourceArn}, nil
		}
		tokenClaimed = true
	}

	if err := store.CreateScheduleGroup(ctx, group); err != nil {
		if tokenClaimed {
			store.ClientTokens().Release(in.ClientToken)
		}
		if err == schedulerstore.ErrScheduleGroupAlreadyExists {
			return nil, ErrScheduleGroupAlreadyExists
		}
		logs.Debug("Failed to create schedule group", logs.String("name", in.Name), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	// Apply tags atomically: if tagging fails after group creation, roll
	// the group back so we never leave an orphan resource.
	if len(in.Tags) > 0 {
		if err := store.TagFromSlice(group.ARN, in.Tags); err != nil {
			logs.Warn("Failed to tag schedule group, rolling back",
				logs.String("arn", group.ARN),
				logs.String("error", err.Error()))
			if tokenClaimed {
				store.ClientTokens().Release(in.ClientToken)
			}
			// Roll the group back: mark deleting and purge immediately
			// (the group was just created and has no member schedules).
			_ = store.MarkScheduleGroupDeleting(ctx, in.Name)
			_ = store.PurgeDeletedScheduleGroup(ctx, in.Name)
			return nil, ErrInternalServer
		}
	}

	return &CreateScheduleGroupResult{ScheduleGroupArn: group.ARN}, nil
}

// deleteScheduleGroupCore validates and marks a schedule group for
// deletion. Deleting a group cascades (the model documentation: the group
// remains in a DELETING state until all of its schedules are deleted); the
// engine's sweep deletes the member schedules and then purges the group.
func (s *SchedulerService) deleteScheduleGroupCore(ctx context.Context, reqCtx *request.RequestContext, in *DeleteScheduleGroupInput) error {
	if in.Name == "" {
		return ErrValidation
	}
	// The default group cannot be deleted (User Guide: "You can't delete,
	// or edit, the default group").
	if in.Name == "default" {
		return awserrors.NewValidationException("cannot delete the default schedule group")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}

	group, err := store.GetScheduleGroup(ctx, in.Name)
	if err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return ErrScheduleGroupNotFound
		}
		logs.Debug("Failed to get schedule group", logs.String("name", in.Name), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	// A repeated delete of a group already in DELETING is idempotent
	// (DeleteScheduleGroup carries the idempotent trait).
	if group.State == schedulerstore.ScheduleGroupStateDeleting {
		return nil
	}

	if err := store.MarkScheduleGroupDeleting(ctx, in.Name); err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return ErrScheduleGroupNotFound
		}
		logs.Debug("Failed to mark schedule group deleting", logs.String("name", in.Name), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// getScheduleGroupCore validates and retrieves a schedule group.
func (s *SchedulerService) getScheduleGroupCore(ctx context.Context, reqCtx *request.RequestContext, in *GetScheduleGroupInput) (*GetScheduleGroupResult, error) {
	if in.Name == "" {
		return nil, ErrValidation
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	group, err := store.GetScheduleGroup(ctx, in.Name)
	if err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return nil, ErrScheduleGroupNotFound
		}
		logs.Debug("Failed to get schedule group", logs.String("name", in.Name), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	return &GetScheduleGroupResult{
		Arn:                  group.ARN,
		Name:                 group.Name,
		State:                string(group.State),
		CreationDate:         group.CreationDate,
		LastModificationDate: group.LastModificationDate,
	}, nil
}

// listScheduleGroupsCore validates the paging parameters and lists schedule
// groups. A zero MaxResults is defaulted to the model's page default.
func (s *SchedulerService) listScheduleGroupsCore(ctx context.Context, reqCtx *request.RequestContext, in *ListScheduleGroupsInput) (*ListScheduleGroupsResult, error) {
	if in.MaxResults == 0 {
		in.MaxResults = DefaultListMaxResults
	}
	if in.MaxResults < 1 || in.MaxResults > MaxListMaxResults {
		return nil, ErrValidation
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := store.ListScheduleGroups(ctx, in.NamePrefix, in.MaxResults, in.NextToken)
	if err != nil {
		logs.Debug("Failed to list schedule groups", logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}

	groups := make([]ScheduleGroupSummary, len(result.ScheduleGroups))
	for i, g := range result.ScheduleGroups {
		groups[i] = ScheduleGroupSummary{
			Arn:                  g.Arn,
			Name:                 g.Name,
			State:                string(g.State),
			CreationDate:         g.CreationDate,
			LastModificationDate: g.LastModificationDate,
		}
	}
	return &ListScheduleGroupsResult{
		ScheduleGroups: groups,
		NextToken:      result.NextToken,
	}, nil
}

// scheduleGroupArnToName validates that the ARN addresses a schedule group
// (the only taggable scheduler resource) and returns the group name.
func scheduleGroupArnToName(resourceArn string) (string, error) {
	_, service, _, _, resource := svcarn.SplitARN(resourceArn)
	if service != "scheduler" {
		return "", ErrValidation
	}
	groupName, ok := strings.CutPrefix(resource, "schedule-group/")
	if !ok || groupName == "" {
		return "", ErrValidation
	}
	// The TagResourceArn pattern constrains the group-name portion to
	// [0-9a-zA-Z-_.]+ ; a malformed ARN is a validation failure, not a
	// missing resource.
	for _, r := range groupName {
		switch {
		case r >= '0' && r <= '9':
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r == '-' || r == '_' || r == '.':
		default:
			return "", ErrValidation
		}
	}
	return groupName, nil
}

// validateScheduleGroupTagTargetCore verifies the tag target ARN addresses
// an existing schedule group. Only schedule groups carry tags: the
// TagResourceArn Smithy pattern accepts schedule-group ARNs only, and the
// TagResource documentation states "You can only assign tags to schedule
// groups."
func (s *SchedulerService) validateScheduleGroupTagTargetCore(ctx context.Context, reqCtx *request.RequestContext, resourceArn string) error {
	groupName, err := scheduleGroupArnToName(resourceArn)
	if err != nil {
		return err
	}
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return ErrScheduleGroupNotFound
		}
		logs.Debug("Failed to get schedule group", logs.String("name", groupName), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// tagScheduleGroupCore applies tags to the schedule group addressed by the
// ARN.
func (s *SchedulerService) tagScheduleGroupCore(ctx context.Context, reqCtx *request.RequestContext, resourceArn string, tags []tagutil.Tag) error {
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	if err := store.TagFromSlice(resourceArn, tags); err != nil {
		logs.Debug("Failed to tag resource", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// untagScheduleGroupCore removes tag keys from the schedule group addressed
// by the ARN.
func (s *SchedulerService) untagScheduleGroupCore(ctx context.Context, reqCtx *request.RequestContext, resourceArn string, tagKeys []string) error {
	store, err := s.store(reqCtx)
	if err != nil {
		return err
	}
	if err := store.Untag(resourceArn, tagKeys); err != nil {
		logs.Debug("Failed to untag resource", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// listScheduleGroupTagsCore lists the tags of the schedule group addressed
// by the ARN.
func (s *SchedulerService) listScheduleGroupTagsCore(ctx context.Context, reqCtx *request.RequestContext, resourceArn string) ([]tagutil.Tag, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	tags, err := store.ListAsSlice(resourceArn)
	if err != nil {
		logs.Debug("Failed to list tags", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}
	return tags, nil
}

// scheduleGroupTagConfig builds the shared tag-handler configuration; every
// closure delegates to a Core function so the handler files carry no store
// access.
func (s *SchedulerService) scheduleGroupTagConfig(reqCtx *request.RequestContext) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			return s.validateScheduleGroupTagTargetCore(ctx, reqCtx, resourceKey)
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			return s.tagScheduleGroupCore(ctx, reqCtx, resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return s.untagScheduleGroupCore(ctx, reqCtx, resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return s.listScheduleGroupTagsCore(ctx, reqCtx, resourceKey)
		},
		ValidateTagsFunc: ValidateScheduleGroupTags,
		EmptyResponse:    func() (interface{}, error) { return response.EmptyResponse(), nil },
		MapError:         schedulerMapError,
	}
}
