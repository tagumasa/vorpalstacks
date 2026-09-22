package cloudwatchlogs

import (
	"errors"
	"fmt"
	"strings"

	awserrors "vorpalstacks/internal/common/errors"
	tagutil "vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// CreateLogGroupInput is the transport-agnostic input for creating a log group.
type CreateLogGroupInput struct {
	LogGroupName              string
	KmsKeyId                  string
	LogGroupClass             string
	Tags                      map[string]string
	DeletionProtectionEnabled bool
	Region                    string
}

// CreateLogGroupResult holds the outcome of a successful log group creation.
type CreateLogGroupResult struct {
	ARN string
}

// DeleteLogGroupInput is the transport-agnostic input for deleting a log group.
type DeleteLogGroupInput struct {
	LogGroupName string
	Region       string
}

// createLogGroupCore creates a new CloudWatch Logs log group after performing
// all validation. Both the HTTP API handler and the admin handler delegate here
// so that validation is shared in a single location.
func (s *LogsService) createLogGroupCore(input CreateLogGroupInput) (*CreateLogGroupResult, error) {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
	}
	// "Log group names can't start with the string aws/" (CreateLogGroup
	// naming guidelines) — the reservation is creation-specific, so it
	// lives here rather than in the shared name validator: a read or
	// delete of a missing aws/-prefixed name keeps the family's
	// ResourceNotFound identity, and the platform's own groups carry the
	// leading-slash /aws/ form that does not match this prefix.
	if strings.HasPrefix(input.LogGroupName, "aws/") {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Log group names can't start with the string aws/: %s", input.LogGroupName), 400)
	}

	if input.LogGroupClass == "" {
		input.LogGroupClass = "STANDARD"
	}
	if !validateLogGroupClass(input.LogGroupClass) {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid log group class: %s. Valid values: STANDARD, INFREQUENT_ACCESS, DELIVERY", input.LogGroupClass), 400)
	}

	if input.KmsKeyId != "" {
		if err := validateKmsKeyId(input.KmsKeyId); err != nil {
			return nil, err
		}
		if err := s.validateLogGroupKmsKey(input.KmsKeyId); err != nil {
			return nil, err
		}
	}

	// The same Tags target the tag operations validate against: 0-50
	// entries (the member is optional on this operation), per-entry
	// TagKey/TagValue traits and the reserved aws: prefix. TooManyTagsException
	// is undeclared here, so the ceiling is InvalidParameterException. Both
	// planes ride this core, so the validation lives here alone.
	if len(input.Tags) > tagutil.MaxTagsPerResource {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("A resource can have a maximum of %d tags", tagutil.MaxTagsPerResource), 400)
	}
	if err := validateTagEntries(input.Tags); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	lg := logsstore.NewLogGroup(input.LogGroupName, input.Region, s.accountID)
	lg.KmsKeyId = input.KmsKeyId
	lg.LogGroupClass = input.LogGroupClass
	lg.DeletionProtectionEnabled = input.DeletionProtectionEnabled

	// "You can create up to 1,000,000 log groups per Region per account"
	// (CreateLogGroup) — LimitExceededException is the operation's
	// declared quota identity.
	count, err := store.CountLogGroups()
	if err != nil {
		return nil, mapStoreError(err)
	}
	if count >= logsstore.MaxLogGroupsPerRegion {
		return nil, NewLogsError("LimitExceededException",
			fmt.Sprintf("You can create up to %d log groups per Region per account", logsstore.MaxLogGroupsPerRegion), 400)
	}

	if err := store.CreateLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

	// Tags live in the tag store alone: every tag-reading surface (the tag
	// operations, the SOURCE query filter) lists the tag store, so the log
	// group record carries no tag copy.
	if len(input.Tags) > 0 {
		if err := store.Tags().Tag(lg.ARN, input.Tags); err != nil {
			if delErr := store.DeleteLogGroup(lg.Name); delErr != nil {
				logs.Error("Failed to rollback log group after tag failure",
					logs.String("logGroup", lg.Name), logs.Err(delErr))
			}
			return nil, mapStoreError(err)
		}
	}

	return &CreateLogGroupResult{ARN: lg.ARN}, nil
}

// deleteLogGroupCore deletes a log group after checking deletion
// protection. Both the HTTP API handler and the admin handler delegate
// here; the protection decision runs inside the store delete's critical
// section, so an enable that commits after this function starts is
// still honoured.
func (s *LogsService) deleteLogGroupCore(input DeleteLogGroupInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	if err := store.DeleteLogGroupIfUnprotected(input.LogGroupName); err != nil {
		return mapStoreError(err)
	}

	// The group's memoised transformer recipe dies with the group: a
	// same-named group recreated later starts with no recipe, not with
	// the deleted one until the cache TTL lapses.
	invalidateTransformerCache(input.Region, input.LogGroupName)
	return nil
}

// TagLogGroupInput is the transport-agnostic input for tagging a log group.
type TagLogGroupInput struct {
	LogGroupName string
	Tags         map[string]string
	Region       string
}

func (s *LogsService) tagLogGroupCore(input TagLogGroupInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}
	// The legacy op's tags member is required over the same Tags target
	// (1-50); TooManyTagsException is declared on TagResource alone, so the
	// ceiling rejects with InvalidParameterException here.
	if len(input.Tags) == 0 || len(input.Tags) > tagutil.MaxTagsPerResource {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The tags member is required and must contain between 1 and %d entries", tagutil.MaxTagsPerResource), 400)
	}
	if err := validateTagEntries(input.Tags); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	lg, err := store.GetLogGroup(input.LogGroupName)
	if err != nil {
		return mapStoreError(err)
	}

	if err := store.Tags().Tag(lg.ARN, input.Tags); err != nil {
		// The cumulative 50-tag ceiling the tag store's merge seam
		// enforces rejects here with the operation's own parameter
		// identity — TooManyTagsException is declared on TagResource
		// alone, and the store's construction budget carries it, so the
		// conversion keys on that wire identity.
		var tagCeiling *awserrors.AWSError
		if errors.As(err, &tagCeiling) && tagCeiling.Code == "TooManyTagsException" {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("A resource can have a maximum of %d tags", tagutil.MaxTagsPerResource), 400)
		}
		return mapStoreError(err)
	}
	return nil
}

// UntagLogGroupInput is the transport-agnostic input for removing tags from a log group.
type UntagLogGroupInput struct {
	LogGroupName string
	TagKeys      []string
	Region       string
}

func (s *LogsService) untagLogGroupCore(input UntagLogGroupInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return err
	}
	// The legacy op's tags member targets TagList, whose length trait is
	// min 1 — an empty list rejects. The typed SDK blocks this client-side
	// (the member is required), so the row is reachable through raw HTTP
	// alone; the operation declares no parameter error beyond
	// ResourceNotFoundException, and InvalidParameterException is the
	// service-wide parameter identity (the least-contradictory choice,
	// recorded in the plan register).
	if len(input.TagKeys) == 0 {
		return NewLogsError("InvalidParameterException",
			"The tags member is required and must contain at least one tag key", 400)
	}
	if err := validateTagKeyElements(input.TagKeys); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	lg, err := store.GetLogGroup(input.LogGroupName)
	if err != nil {
		return mapStoreError(err)
	}

	if err := store.Tags().Untag(lg.ARN, input.TagKeys); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// ListTagsLogGroupInput is the transport-agnostic input for listing tags on a log group.
type ListTagsLogGroupInput struct {
	LogGroupName string
	Region       string
}

func (s *LogsService) listTagsLogGroupCore(input ListTagsLogGroupInput) (map[string]string, error) {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	lg, err := store.GetLogGroup(input.LogGroupName)
	if err != nil {
		return nil, mapStoreError(err)
	}

	tags, err := store.Tags().List(lg.ARN)
	if err != nil {
		return nil, mapStoreError(err)
	}
	return tags, nil
}
