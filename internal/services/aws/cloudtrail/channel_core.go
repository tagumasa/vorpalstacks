package cloudtrail

import (
	"errors"
	"fmt"
	"strings"

	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// ---------------------------------------------------------------------------
// Transport-agnostic Input structs
// ---------------------------------------------------------------------------

// CreateChannelInput carries the create members for a channel. Destinations
// and Tags keep their raw wire values plus presence flags so the Core can
// distinguish "not provided" from an explicit empty value.
type CreateChannelInput struct {
	Name            string
	Source          string
	DestinationsRaw interface{}
	DestinationsSet bool
	TagList         []tags.Tag
	TagsRaw         interface{}
	TagsSet         bool
}

// ChannelInput carries the channel ARN.
type ChannelInput struct {
	Channel string
}

// ListChannelsInput carries pagination parameters for listing channels.
type ListChannelsInput struct {
	NextToken  string
	MaxResults int
}

// UpdateChannelInput carries the update members for a channel. Only provided
// members are applied.
type UpdateChannelInput struct {
	Channel         string
	Name            string
	DestinationsRaw interface{}
	DestinationsSet bool
}

// ---------------------------------------------------------------------------
// Core functions — single validation + persistence path
// ---------------------------------------------------------------------------

// createChannelCore is the single entry point for CreateChannel.
func (s *CloudTrailService) createChannelCore(store cloudtrailstore.CloudTrailStoreInterface, in CreateChannelInput) (map[string]interface{}, error) {
	if in.Name == "" {
		return nil, newInvalidParameterException("Name is required")
	}
	if err := validateChannelName(in.Name); err != nil {
		return nil, err
	}

	if in.Source == "" {
		return nil, newInvalidSourceException("Source is required")
	}
	if err := validateChannelSource(in.Source); err != nil {
		return nil, err
	}

	// Destinations is model-required with a minimum of one entry: a
	// destination-less channel is never created.
	if !in.DestinationsSet {
		return nil, newInvalidParameterException("Destinations is required")
	}
	destinations, err := parseDestinations(in.DestinationsRaw)
	if err != nil {
		return nil, err
	}
	if err := validateDestinations(destinations); err != nil {
		return nil, err
	}
	if _, err := resolveChannelDestinationEDSs(store, destinations); err != nil {
		return nil, err
	}

	ch := cloudtrailstore.NewChannel(in.Name, in.Source, store.GetAccountID(), store.GetRegion())
	ch.Destinations = destinations

	// Parse and validate tags before creation. CreateChannel uses the member
	// name "Tags" (not "TagsList" as Trail and EDS do), per Smithy model.
	if in.TagsSet {
		if err := validateCloudTrailTags(in.TagList); err != nil {
			return nil, err
		}
		applyTags(&ch.Tags, in.TagsRaw)
	}

	// The channel quota is enforced inside CreateChannel (count and write
	// under the store mutex): "Channels — 25" (Quotas in AWS CloudTrail);
	// the model declares ChannelMaxLimitExceededException on CreateChannel.

	created, err := store.CreateChannel(ch)
	if err != nil {
		if errors.Is(err, cloudtrailstore.ErrChannelQuotaExceeded) {
			return nil, newChannelMaxLimitExceededException(
				fmt.Sprintf("The maximum number of channels (%d) has been reached", cloudtrailstore.MaxChannelsPerRegion))
		}
		return nil, s.mapStoreError(err)
	}

	// CreateChannelResponse carries the channel's identity and source plus
	// the applied tags.
	resp := formatChannelBase(created)
	if len(created.Tags) > 0 {
		resp["Tags"] = formatTagsList(created.Tags)
	}
	return resp, nil
}

// deleteChannelCore is the single entry point for DeleteChannel: the
// reference documents no event-data-store precondition ("Deletes a
// channel"), so a channel is deletable regardless of its destinations. The
// association constraint lives on the other side: DeleteEventDataStore
// refuses while a channel points at the store
// (ChannelExistsForEDSException), so teardown deletes the channel first.
func (s *CloudTrailService) deleteChannelCore(store cloudtrailstore.CloudTrailStoreInterface, in ChannelInput) (map[string]interface{}, error) {
	if in.Channel == "" {
		return nil, newInvalidParameterException("Channel is required")
	}

	ch, err := s.resolveChannelCore(store, in.Channel)
	if err != nil {
		return nil, err
	}

	if err := store.DeleteChannelIf(ch.ChannelARN, func(_ *cloudtrailstore.Channel) error {
		return nil
	}); err != nil {
		return nil, s.mapStoreError(err)
	}

	return map[string]interface{}{}, nil
}

// getChannelCore is the single entry point for GetChannel.
func (s *CloudTrailService) getChannelCore(store cloudtrailstore.CloudTrailStoreInterface, in ChannelInput) (map[string]interface{}, error) {
	if in.Channel == "" {
		return nil, newInvalidParameterException("Channel is required")
	}

	ch, err := s.resolveChannelCore(store, in.Channel)
	if err != nil {
		return nil, err
	}

	// GetChannelResponse adds SourceConfig and IngestionStatus to the base;
	// both describe the channel's ingestion machinery, which has no data
	// until channel ingestion exists — they stay absent rather than
	// synthesised.
	return formatChannelBase(ch), nil
}

// listChannelsCore is the single entry point for ListChannels.
func (s *CloudTrailService) listChannelsCore(store cloudtrailstore.CloudTrailStoreInterface, in ListChannelsInput) (map[string]interface{}, error) {
	opts := storecommon.ListOptions{MaxItems: cloudtrailstore.DefaultListChannelsResults}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		// ListChannels declares no InvalidMaxResultsException; its generic
		// declared invalid-input error carries the bound rejection.
		if in.MaxResults > cloudtrailstore.MaxListChannelsResults {
			return nil, newInvalidParameterException(
				fmt.Sprintf("MaxResults exceeds the maximum of %d", cloudtrailstore.MaxListChannelsResults))
		}
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListChannels(opts)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// The list item is the two-member Channel shape: ARN and name alone.
	chList := make([]map[string]interface{}, 0, len(result.Items))
	for _, ch := range result.Items {
		chList = append(chList, map[string]interface{}{
			"ChannelArn": ch.ChannelARN,
			"Name":       ch.Name,
		})
	}

	resp := map[string]interface{}{
		"Channels": chList,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// updateChannelCore is the single entry point for UpdateChannel. The
// provided members are validated against the model constraints and applied
// to the stored record under the store mutex, so concurrent updates cannot
// lose writes.
func (s *CloudTrailService) updateChannelCore(store cloudtrailstore.CloudTrailStoreInterface, in UpdateChannelInput) (map[string]interface{}, error) {
	if in.Channel == "" {
		return nil, newInvalidParameterException("Channel is required")
	}

	if in.Name != "" {
		if err := validateChannelName(in.Name); err != nil {
			return nil, err
		}
	}
	var destinations []cloudtrailstore.Destination
	if in.DestinationsSet {
		parsed, parseErr := parseDestinations(in.DestinationsRaw)
		if parseErr != nil {
			return nil, parseErr
		}
		destinations = parsed
		if err := validateDestinations(destinations); err != nil {
			return nil, err
		}
		if _, err := resolveChannelDestinationEDSs(store, destinations); err != nil {
			return nil, err
		}
	}

	ch, err := s.resolveChannelCore(store, in.Channel)
	if err != nil {
		return nil, err
	}

	updated, err := store.MutateChannel(ch.ChannelARN, func(ch *cloudtrailstore.Channel) error {
		if in.Name != "" {
			ch.Name = in.Name
		}
		if in.DestinationsSet {
			ch.Destinations = destinations
		}
		return nil
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	// UpdateChannelResponse carries the base alone — no tags, no ingestion
	// members.
	return formatChannelBase(updated), nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// resolveChannelByARNOrUUID resolves a channel by full ARN or bare UUID
// ("The ARN or UUID of a channel", GetChannel). A value that carries the
// ARN prefix but is not a well-formed channel ARN answers invalidARN;
// channels are keyed by full ARN, so a bare value is matched against the
// UUID suffix of every stored channel and a miss returns the store's
// not-found sentinel for the caller's own mapping.
func resolveChannelByARNOrUUID(store cloudtrailstore.CloudTrailStoreInterface, value string, invalidARN func(string) error) (*cloudtrailstore.Channel, error) {
	if strings.HasPrefix(value, "arn:") {
		if !isCloudTrailResourceARN(value, "channel/") {
			return nil, invalidARN(value)
		}
		return store.GetChannel(value)
	}

	marker := ""
	for {
		result, err := store.ListChannels(storecommon.ListOptions{
			MaxItems: cloudtrailstore.MaxListChannelsResults,
			Marker:   marker,
		})
		if err != nil {
			return nil, err
		}
		for _, ch := range result.Items {
			_, _, _, _, resource := svcarn.SplitARN(ch.ChannelARN)
			if id, ok := strings.CutPrefix(resource, "channel/"); ok && id == value {
				return ch, nil
			}
		}
		if result.NextMarker == "" {
			return nil, cloudtrailstore.ErrChannelNotFound
		}
		marker = result.NextMarker
	}
}

// resolveChannelCore resolves a channel selector for the management
// plane: an invalid ARN answers the channel operations' declared
// ChannelARNInvalidException and an absent channel the declared
// ChannelNotFoundException, whatever identifier form named it.
func (s *CloudTrailService) resolveChannelCore(store cloudtrailstore.CloudTrailStoreInterface, value string) (*cloudtrailstore.Channel, error) {
	ch, err := resolveChannelByARNOrUUID(store, value, func(v string) error {
		return newChannelARNInvalidException(
			fmt.Sprintf("The specified channel ARN is not valid: %s", v))
	})
	if err != nil {
		return nil, s.mapStoreError(err)
	}
	return ch, nil
}

// parseDestinations parses destinations from their raw wire value.
// parseDestinations parses the Destinations wire member strictly: an entry
// that is not an object, or a Type or Location that is not a string, is
// rejected rather than silently dropped — a partially parsed destination
// list would create a channel that ignores part of the request.
func parseDestinations(raw interface{}) ([]cloudtrailstore.Destination, error) {
	arr, ok := raw.([]interface{})
	if !ok {
		return nil, newInvalidParameterException(
			"Destinations must be a list of destination objects")
	}
	result := make([]cloudtrailstore.Destination, 0, len(arr))
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, newInvalidParameterException(
				"Each destination must be an object")
		}
		destType, ok := m["Type"].(string)
		if !ok {
			return nil, newInvalidParameterException(
				"Destination Type must be a string")
		}
		location, ok := m["Location"].(string)
		if !ok {
			return nil, newInvalidParameterException(
				"Destination Location must be a string")
		}
		result = append(result, cloudtrailstore.Destination{
			Type:     destType,
			Location: location,
		})
	}
	return result, nil
}

// resolveChannelDestinationEDSs resolves every EVENT_DATA_STORE destination
// to its event data store. CreateChannel and UpdateChannel both declare
// EventDataStoreARNInvalidException and EventDataStoreNotFoundException —
// the Location member is "the ARN of an event data store that receives
// events from a channel" (Destination), so a destination that is not a
// well-formed event data store ARN, or that names no store this account
// owns, is refused. AWS_SERVICE destinations name a service, not a store,
// and are not resolved.
func resolveChannelDestinationEDSs(store cloudtrailstore.CloudTrailStoreInterface, destinations []cloudtrailstore.Destination) ([]*cloudtrailstore.EventDataStore, error) {
	var resolved []*cloudtrailstore.EventDataStore
	for _, dest := range destinations {
		if dest.Type != cloudtrailstore.DestinationTypeEventDataStore {
			continue
		}
		if !isCloudTrailResourceARN(dest.Location, "eventdatastore/") {
			return nil, newEventDataStoreARNInvalidException(
				fmt.Sprintf("The destination location is not a valid event data store ARN: %s", dest.Location))
		}
		eds, getErr := store.GetEventDataStore(dest.Location)
		if getErr != nil {
			if errors.Is(getErr, cloudtrailstore.ErrEventDataStoreNotFound) {
				return nil, newEventDataStoreNotFoundException(
					fmt.Sprintf("The destination event data store was not found: %s", dest.Location))
			}
			return nil, getErr
		}
		resolved = append(resolved, eds)
	}
	return resolved, nil
}

// formatChannelBase renders the members every channel response shape
// shares: identity, source and destinations. The per-operation shapes add
// their own members (tags on create, ingestion state on get) at their call
// sites; the IngestionStatus wire structure carries the Latest* ingestion
// attempt and success members, none of which exist before channel
// ingestion does.
func formatChannelBase(ch *cloudtrailstore.Channel) map[string]interface{} {
	resp := map[string]interface{}{
		"ChannelArn": ch.ChannelARN,
		"Name":       ch.Name,
		"Source":     ch.Source,
	}
	if len(ch.Destinations) > 0 {
		dests := make([]map[string]interface{}, 0, len(ch.Destinations))
		for _, d := range ch.Destinations {
			dm := map[string]interface{}{"Type": d.Type}
			if d.Location != "" {
				dm["Location"] = d.Location
			}
			dests = append(dests, dm)
		}
		resp["Destinations"] = dests
	}
	return resp
}
