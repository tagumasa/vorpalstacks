package scheduler

import (
	"context"
	"strings"
	"time"
	"unicode/utf8"

	awserrors "vorpalstacks/internal/common/errors"
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

// DeleteScheduleGroupInput carries the DeleteScheduleGroup request.
// ClientToken preserves the wire member's idempotency semantics
// (@length(1,64), @pattern): an invalid token is rejected, and a replayed
// token reports the first deletion's outcome instead of surfacing not-found
// for a purged group — the same protocol DeleteSchedule carries, layered
// over the DELETING-state idempotency the operation already has.
type DeleteScheduleGroupInput struct {
	Name        string
	ClientToken string
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

// ListScheduleGroupsInput carries the ListScheduleGroups request. MaxResults
// is nil when the member was absent; the Core applies the default page size.
type ListScheduleGroupsInput struct {
	NamePrefix string
	// MaxResults is nil when the member was absent from the request.
	MaxResults *int32
	NextToken  string
}

// ListScheduleGroupsResult carries the ListScheduleGroups response. The list
// path reuses the store's summary type, the same convention the schedule
// listing follows — the summary is a projection of the store record, with no
// service-side transformation.
type ListScheduleGroupsResult struct {
	ScheduleGroups []schedulerstore.ScheduleGroupSummary
	NextToken      string
}

// createScheduleGroupCore validates and creates a schedule group, applies
// its tags, and honours ClientToken idempotency. A tagging failure rolls
// the group back so no orphan resource remains.
func (s *SchedulerService) createScheduleGroupCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *CreateScheduleGroupInput) (*CreateScheduleGroupResult, error) {
	if in.Name == "" || !namePattern.MatchString(in.Name) {
		return nil, ErrValidation
	}
	// Validate tags before creating the group so invalid tag sets are
	// rejected without creating an orphan resource.
	if err := ValidateScheduleGroupTags(in.Tags); err != nil {
		return nil, err
	}

	group := &schedulerstore.ScheduleGroup{Name: in.Name}

	replayArn, releaseToken, err := claimClientToken(store, in.ClientToken, store.BuildScheduleGroupARN(in.Name), "schedule-group")
	if err != nil {
		return nil, err
	}
	if replayArn != "" {
		return &CreateScheduleGroupResult{ScheduleGroupArn: replayArn}, nil
	}

	if err := store.CreateScheduleGroup(ctx, group); err != nil {
		releaseToken()
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
			releaseToken()
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
func (s *SchedulerService) deleteScheduleGroupCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *DeleteScheduleGroupInput) error {
	// The ScheduleGroupName shape (pattern + length, which also rejects an
	// empty name) is validated before any resource lookup.
	if err := validateScheduleGroupName(in.Name); err != nil {
		return err
	}
	// The default group cannot be deleted (User Guide: "You can't delete,
	// or edit, the default group").
	if in.Name == schedulerstore.DefaultGroupName {
		return awserrors.NewValidationException("cannot delete the default schedule group")
	}
	// A replayed idempotency token reports the first deletion's outcome;
	// an unrecoverable failure releases the claim so a retry re-executes.
	replayArn, releaseToken, err := claimClientToken(store, in.ClientToken, store.BuildScheduleGroupARN(in.Name), "schedule-group-delete")
	if err != nil {
		return err
	}
	if replayArn != "" {
		return nil
	}

	group, err := store.GetScheduleGroup(ctx, in.Name)
	if err != nil {
		releaseToken()
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return scheduleGroupNotFound(in.Name)
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
		releaseToken()
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return scheduleGroupNotFound(in.Name)
		}
		logs.Debug("Failed to mark schedule group deleting", logs.String("name", in.Name), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// getScheduleGroupCore validates and retrieves a schedule group.
func (s *SchedulerService) getScheduleGroupCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *GetScheduleGroupInput) (*GetScheduleGroupResult, error) {
	// The ScheduleGroupName shape (pattern + length, which also rejects an
	// empty name) is validated before any resource lookup.
	if err := validateScheduleGroupName(in.Name); err != nil {
		return nil, err
	}

	group, err := store.GetScheduleGroup(ctx, in.Name)
	if err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return nil, scheduleGroupNotFound(in.Name)
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

// listScheduleGroupsCore validates the filter and paging parameters and
// lists schedule groups. An absent MaxResults is defaulted to the model's
// page default; an explicitly invalid value is rejected.
func (s *SchedulerService) listScheduleGroupsCore(ctx context.Context, store *schedulerstore.SchedulerStore, in *ListScheduleGroupsInput) (*ListScheduleGroupsResult, error) {
	maxResults, err := resolveListMaxResults(in.MaxResults)
	if err != nil {
		return nil, err
	}
	if err := validateListNamePrefix(in.NamePrefix); err != nil {
		return nil, err
	}
	if err := validateNextToken(in.NextToken); err != nil {
		return nil, err
	}

	result, err := store.ListScheduleGroups(ctx, in.NamePrefix, maxResults, in.NextToken)
	if err != nil {
		logs.Debug("Failed to list schedule groups", logs.String("error", err.Error()))
		return nil, ErrInternalServer
	}
	return &ListScheduleGroupsResult{
		ScheduleGroups: result.ScheduleGroups,
		NextToken:      result.NextToken,
	}, nil
}

// scheduleGroupArnToName validates that the ARN addresses a schedule group
// (the only taggable scheduler resource) and returns the group name.
func scheduleGroupArnToName(resourceArn string) (string, error) {
	// TagResourceArn @length(1, 1011) in characters: an over-length ARN is
	// a validation failure, not a missing resource. The shape's pattern
	// charset is ASCII, so conforming ARNs never rely on the counting
	// basis; the check still counts characters like every @length bound.
	if utf8.RuneCountInString(resourceArn) > maxTagResourceArnLength {
		return "", ErrValidation
	}
	_, service, _, _, resource := svcarn.SplitARN(resourceArn)
	if service != "scheduler" {
		return "", ErrValidation
	}
	groupName, ok := strings.CutPrefix(resource, "schedule-group/")
	if !ok {
		return "", ErrValidation
	}
	// The ARN-embedded group name carries the ScheduleGroupName shape
	// constraints: TagResourceArn's pattern supplies the charset, and the
	// 1-64 bound holds because a tag target can only name a group whose
	// creation passed the same bound — namePattern is the single
	// definition of that grammar. A malformed ARN is a validation
	// failure, not a missing resource.
	if !namePattern.MatchString(groupName) {
		return "", ErrValidation
	}
	return groupName, nil
}

// validateScheduleGroupTagTargetCore verifies the tag target ARN addresses
// an existing schedule group. Only schedule groups carry tags: the
// TagResourceArn Smithy pattern accepts schedule-group ARNs only, and the
// TagResource documentation states "You can only assign tags to schedule
// groups."
func (s *SchedulerService) validateScheduleGroupTagTargetCore(ctx context.Context, store *schedulerstore.SchedulerStore, resourceArn string) error {
	groupName, err := scheduleGroupArnToName(resourceArn)
	if err != nil {
		return err
	}
	if _, err := store.GetScheduleGroup(ctx, groupName); err != nil {
		if err == schedulerstore.ErrScheduleGroupNotFound {
			return scheduleGroupNotFound(groupName)
		}
		logs.Debug("Failed to get schedule group", logs.String("name", groupName), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// tagScheduleGroupCore applies tags to the schedule group addressed by the
// ARN.
func (s *SchedulerService) tagScheduleGroupCore(ctx context.Context, store *schedulerstore.SchedulerStore, resourceArn string, tags []tagutil.Tag) error {
	if err := store.TagFromSlice(resourceArn, tags); err != nil {
		logs.Debug("Failed to tag resource", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// untagScheduleGroupCore removes tag keys from the schedule group addressed
// by the ARN.
func (s *SchedulerService) untagScheduleGroupCore(ctx context.Context, store *schedulerstore.SchedulerStore, resourceArn string, tagKeys []string) error {
	if err := store.Untag(resourceArn, tagKeys); err != nil {
		logs.Debug("Failed to untag resource", logs.String("arn", resourceArn), logs.String("error", err.Error()))
		return ErrInternalServer
	}
	return nil
}

// listScheduleGroupTagsCore lists the tags of the schedule group addressed
// by the ARN.
func (s *SchedulerService) listScheduleGroupTagsCore(ctx context.Context, store *schedulerstore.SchedulerStore, resourceArn string) ([]tagutil.Tag, error) {
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
func (s *SchedulerService) scheduleGroupTagConfig(store *schedulerstore.SchedulerStore) tagutil.TagHandlerConfig {
	return tagutil.TagHandlerConfig{
		Param: tagutil.StandardConfig,
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			return s.validateScheduleGroupTagTargetCore(ctx, store, resourceKey)
		},
		TagFunc: func(ctx context.Context, resourceKey string, tags []tagutil.Tag) error {
			return s.tagScheduleGroupCore(ctx, store, resourceKey, tags)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return s.untagScheduleGroupCore(ctx, store, resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tagutil.Tag, error) {
			return s.listScheduleGroupTagsCore(ctx, store, resourceKey)
		},
		ValidateTagsFunc: ValidateScheduleGroupTags,
		EmptyResponse:    func() (interface{}, error) { return response.EmptyResponse(), nil },
		MapError:         schedulerMapError,
	}
}
