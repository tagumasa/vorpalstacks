package cloudwatchlogs

import (
	"fmt"
	"sort"

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

// ListLogGroupsInput is the transport-agnostic input for listing log groups.
type ListLogGroupsInput struct {
	LogGroupNamePrefix string
	LogGroupClass      string
	NextToken          string
	Limit              int32
	Region             string
}

// ListLogGroupsResult holds the outcome of a successful log group listing.
type ListLogGroupsResult struct {
	LogGroups []*logsstore.LogGroup
	NextToken string
}

// DescribeLogStreamsInput is the transport-agnostic input for listing log streams.
type DescribeLogStreamsInput struct {
	LogGroupName        string
	LogStreamNamePrefix string
	OrderBy             string
	Descending          bool
	NextToken           string
	Limit               int32
	Region              string
}

// DescribeLogStreamsResult holds the outcome of a successful log stream listing.
type DescribeLogStreamsResult struct {
	LogStreams []*logsstore.LogStream
	NextToken  string
}

// createLogGroupCore creates a new CloudWatch Logs log group after performing
// all validation. Both the HTTP API handler and the admin handler delegate here
// so that validation is shared in a single location.
func (s *LogsService) createLogGroupCore(input CreateLogGroupInput) (*CreateLogGroupResult, error) {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
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
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	lg := logsstore.NewLogGroup(input.LogGroupName, input.Region, s.accountID)
	lg.KmsKeyId = input.KmsKeyId
	lg.LogGroupClass = input.LogGroupClass
	lg.DeletionProtectionEnabled = input.DeletionProtectionEnabled
	lg.Tags = input.Tags

	if err := store.CreateLogGroup(lg); err != nil {
		return nil, mapStoreError(err)
	}

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

// deleteLogGroupCore deletes a log group after checking deletion protection.
// Both the HTTP API handler and the admin handler delegate here.
func (s *LogsService) deleteLogGroupCore(input DeleteLogGroupInput) error {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
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

	if lg.DeletionProtectionEnabled {
		return ErrOperationAborted
	}

	if err := store.DeleteLogGroup(input.LogGroupName); err != nil {
		return mapStoreError(err)
	}

	return nil
}

// listLogGroupsCore lists log groups with optional class filtering.
// The caller is responsible for enforcing the operation-specific Smithy limit
// (DescribeLimit max=50 for DescribeLogGroups, ListLimit max=1000 for
// ListLogGroups) before calling this method.
func (s *LogsService) listLogGroupsCore(input ListLogGroupsInput) (*ListLogGroupsResult, error) {
	if input.Limit <= 0 {
		input.Limit = 50
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	groups, nextToken, err := store.ListLogGroups(input.LogGroupNamePrefix, input.NextToken, int(input.Limit))
	if err != nil {
		return nil, mapStoreError(err)
	}

	if input.LogGroupClass != "" {
		filtered := make([]*logsstore.LogGroup, 0, len(groups))
		for _, lg := range groups {
			if lg.LogGroupClass == input.LogGroupClass {
				filtered = append(filtered, lg)
			}
		}
		groups = filtered
	}

	return &ListLogGroupsResult{
		LogGroups: groups,
		NextToken: nextToken,
	}, nil
}

// describeLogStreamsCore lists log streams with ordering and pagination.
func (s *LogsService) describeLogStreamsCore(input DescribeLogStreamsInput) (*DescribeLogStreamsResult, error) {
	if err := validateLogGroupName(input.LogGroupName); err != nil {
		return nil, err
	}

	limit, err := validateListLimit(input.Limit, 50, 50)
	if err != nil {
		return nil, err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return nil, err
	}

	if input.OrderBy == "LastEventTime" {
		if input.LogStreamNamePrefix != "" {
			return nil, NewLogsError("InvalidParameterException",
				"Cannot specify logStreamNamePrefix when orderBy is LastEventTime", 400)
		}

		allStreams, err := fetchAllLogStreams(store, input.LogGroupName, "")
		if err != nil {
			return nil, mapStoreError(err)
		}

		sort.Slice(allStreams, func(i, j int) bool {
			if input.Descending {
				return allStreams[i].LastEventTs > allStreams[j].LastEventTs
			}
			return allStreams[i].LastEventTs < allStreams[j].LastEventTs
		})

		_, offset, err := logsstore.ParsePaginationToken(input.NextToken)
		if err != nil {
			return nil, ErrInvalidParameter
		}
		if offset > len(allStreams) {
			offset = len(allStreams)
		}

		endIdx := offset + int(limit)
		if endIdx > len(allStreams) {
			endIdx = len(allStreams)
		}

		var pageItems []*logsstore.LogStream
		if offset < len(allStreams) {
			pageItems = allStreams[offset:endIdx]
		}

		result := &DescribeLogStreamsResult{LogStreams: pageItems}
		if endIdx < len(allStreams) {
			result.NextToken = logsstore.EncodePaginationToken(logsstore.PaginationForward, endIdx)
		}
		return result, nil
	}

	if input.Descending {
		allStreams, err := fetchAllLogStreams(store, input.LogGroupName, input.LogStreamNamePrefix)
		if err != nil {
			return nil, mapStoreError(err)
		}

		sort.Slice(allStreams, func(i, j int) bool {
			return allStreams[i].Name > allStreams[j].Name
		})

		result := &DescribeLogStreamsResult{}
		if len(allStreams) > 0 {
			pageStart := 0
			if input.NextToken != "" {
				for i, ls := range allStreams {
					if ls.Name == input.NextToken {
						pageStart = i
						break
					}
				}
			}
			endIdx := pageStart + int(limit)
			if endIdx > len(allStreams) {
				endIdx = len(allStreams)
			}
			result.LogStreams = allStreams[pageStart:endIdx]
			if endIdx < len(allStreams) {
				result.NextToken = allStreams[endIdx].Name
			}
		}
		return result, nil
	}

	streams, nextToken, err := store.ListLogStreams(input.LogGroupName, input.LogStreamNamePrefix, input.NextToken, int(limit))
	if err != nil {
		return nil, mapStoreError(err)
	}

	return &DescribeLogStreamsResult{
		LogStreams: streams,
		NextToken:  nextToken,
	}, nil
}
