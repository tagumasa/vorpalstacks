package cloudtrail

import (
	"encoding/json"

	awserrors "vorpalstacks/internal/common/errors"
	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
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
		return nil, awserrors.NewAWSError("InvalidParameter", "Name is required", 400)
	}

	if in.Source == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Source is required", 400)
	}

	ch := cloudtrailstore.NewChannel(in.Name, in.Source, store.GetAccountID(), store.GetRegion())

	if in.DestinationsSet {
		ch.Destinations = parseDestinations(in.DestinationsRaw)
	}

	// Parse and validate tags before creation. CreateChannel uses the member
	// name "Tags" (not "TagsList" as Trail and EDS do), per Smithy model.
	if in.TagsSet {
		if err := validateCloudTrailTags(in.TagList); err != nil {
			return nil, err
		}
		applyChannelTags(ch, in.TagsRaw)
	}

	created, err := store.CreateChannel(ch)
	if err != nil {
		return nil, err
	}

	return formatChannel(created), nil
}

// deleteChannelCore is the single entry point for DeleteChannel. Per Smithy,
// DeleteChannel returns OperationNotPermittedException when the channel is
// referenced by an active EDS.
func (s *CloudTrailService) deleteChannelCore(store cloudtrailstore.CloudTrailStoreInterface, in ChannelInput) (map[string]interface{}, error) {
	if in.Channel == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Channel is required", 400)
	}

	ch, err := store.GetChannel(in.Channel)
	if err != nil {
		return nil, awserrors.NewAWSError("ChannelNotFoundException", "Channel not found", 404)
	}

	// Check if any event data store depends on this channel.
	for _, dest := range ch.Destinations {
		if dest.EDSARN != "" {
			if eds, err := store.GetEventDataStore(dest.EDSARN); err == nil {
				if eds.Status == "ENABLED" {
					return nil, awserrors.NewAWSError("OperationNotPermittedException",
						"Cannot delete channel because it is associated with an active event data store", 400)
				}
			}
		}
	}

	if err := store.DeleteChannel(in.Channel); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// getChannelCore is the single entry point for GetChannel.
func (s *CloudTrailService) getChannelCore(store cloudtrailstore.CloudTrailStoreInterface, in ChannelInput) (map[string]interface{}, error) {
	if in.Channel == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Channel is required", 400)
	}

	ch, err := store.GetChannel(in.Channel)
	if err != nil {
		return nil, awserrors.NewAWSError("ChannelNotFoundException", "Channel not found", 404)
	}

	return formatChannel(ch), nil
}

// listChannelsCore is the single entry point for ListChannels.
func (s *CloudTrailService) listChannelsCore(store cloudtrailstore.CloudTrailStoreInterface, in ListChannelsInput) (map[string]interface{}, error) {
	opts := storecommon.ListOptions{MaxItems: 100}
	if in.NextToken != "" {
		opts.Marker = in.NextToken
	}
	if in.MaxResults > 0 {
		opts.MaxItems = in.MaxResults
	}

	result, err := store.ListChannels(opts)
	if err != nil {
		return nil, err
	}

	chList := make([]map[string]interface{}, 0, len(result.Items))
	for _, ch := range result.Items {
		chList = append(chList, formatChannel(ch))
	}

	resp := map[string]interface{}{
		"Channels": chList,
	}
	if result.NextMarker != "" {
		resp["NextToken"] = result.NextMarker
	}

	return resp, nil
}

// updateChannelCore is the single entry point for UpdateChannel.
func (s *CloudTrailService) updateChannelCore(store cloudtrailstore.CloudTrailStoreInterface, in UpdateChannelInput) (map[string]interface{}, error) {
	if in.Channel == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Channel is required", 400)
	}

	ch, err := store.GetChannel(in.Channel)
	if err != nil {
		return nil, awserrors.NewAWSError("ChannelNotFoundException", "Channel not found", 404)
	}

	if in.Name != "" {
		ch.Name = in.Name
	}
	if in.DestinationsSet {
		ch.Destinations = parseDestinations(in.DestinationsRaw)
	}

	if err := store.UpdateChannel(ch); err != nil {
		return nil, err
	}

	return formatChannel(ch), nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// parseDestinations parses destinations from their raw wire value.
func parseDestinations(raw interface{}) []cloudtrailstore.Destination {
	var result []cloudtrailstore.Destination
	arr, ok := raw.([]interface{})
	if !ok {
		return result
	}
	for _, item := range arr {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		dest := cloudtrailstore.Destination{}
		if t, ok := m["Type"].(string); ok {
			dest.Type = t
		}
		if l, ok := m["Location"].(string); ok {
			dest.Location = l
		}
		result = append(result, dest)
	}
	return result
}

// formatChannel converts a channel to the API response format.
func formatChannel(ch *cloudtrailstore.Channel) map[string]interface{} {
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
	if ch.IngestionStatus != "" {
		resp["IngestionStatus"] = map[string]interface{}{
			"State": ch.IngestionStatus,
		}
	}
	if len(ch.Tags) > 0 {
		tagsList := make([]interface{}, 0, len(ch.Tags))
		for k, v := range ch.Tags {
			tagsList = append(tagsList, map[string]interface{}{
				"Key":   k,
				"Value": v,
			})
		}
		resp["Tags"] = tagsList
	}
	return resp
}

// applyChannelTags parses tags from the raw interface and applies them to the
// channel.
func applyChannelTags(ch *cloudtrailstore.Channel, raw interface{}) {
	if ch.Tags == nil {
		ch.Tags = make(map[string]string)
	}
	var tagsList []interface{}
	switch v := raw.(type) {
	case []interface{}:
		tagsList = v
	case string:
		if err := json.Unmarshal([]byte(v), &tagsList); err != nil {
			return
		}
	default:
		return
	}
	for _, item := range tagsList {
		if m, ok := item.(map[string]interface{}); ok {
			key, _ := m["Key"].(string)
			val, _ := m["Value"].(string)
			if key != "" {
				ch.Tags[key] = val
			}
		}
	}
}
