package cloudtrail

import (
	"context"
	"encoding/json"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	tags "vorpalstacks/internal/common/tags"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	storecommon "vorpalstacks/internal/store/aws/common"
)

// parseDestinations parses destinations from request parameters.
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

// CreateChannel creates a new CloudTrail channel.
func (s *CloudTrailService) CreateChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	name := request.GetStringParam(req.Parameters, "Name")
	if name == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Name is required", 400)
	}

	source := request.GetStringParam(req.Parameters, "Source")
	if source == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Source is required", 400)
	}

	ch := cloudtrailstore.NewChannel(name, source, store.GetAccountID(), store.GetRegion())

	if destsRaw, ok := req.Parameters["Destinations"]; ok {
		ch.Destinations = parseDestinations(destsRaw)
	}

	// Parse and validate tags before creation. CreateChannel uses the member
	// name "Tags" (not "TagsList" as Trail and EDS do), per Smithy model.
	if tagsRaw, ok := req.Parameters["Tags"]; ok {
		tagList := tags.ParseTags(req.Parameters, "Tags")
		if err := validateCloudTrailTags(tagList); err != nil {
			return nil, err
		}
		applyChannelTags(ch, tagsRaw)
	}

	created, err := store.CreateChannel(ch)
	if err != nil {
		return nil, err
	}

	return formatChannel(created), nil
}

// DeleteChannel deletes the specified CloudTrail channel.
func (s *CloudTrailService) DeleteChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	arn := request.GetStringParam(req.Parameters, "Channel")
	if arn == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Channel is required", 400)
	}

	if _, err := store.GetChannel(arn); err != nil {
		return nil, awserrors.NewAWSError("ChannelNotFoundException", "Channel not found", 404)
	}

	if err := store.DeleteChannel(arn); err != nil {
		return nil, err
	}

	return map[string]interface{}{}, nil
}

// GetChannel retrieves the specified CloudTrail channel.
func (s *CloudTrailService) GetChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	arn := request.GetStringParam(req.Parameters, "Channel")
	if arn == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Channel is required", 400)
	}

	ch, err := store.GetChannel(arn)
	if err != nil {
		return nil, awserrors.NewAWSError("ChannelNotFoundException", "Channel not found", 404)
	}

	return formatChannel(ch), nil
}

// ListChannels lists CloudTrail channels with pagination.
func (s *CloudTrailService) ListChannels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	opts := storecommon.ListOptions{MaxItems: 100}
	if nextToken := req.GetParam("NextToken"); nextToken != "" {
		opts.Marker = nextToken
	}
	if maxResults := request.GetIntParam(req.Parameters, "MaxResults"); maxResults > 0 {
		opts.MaxItems = maxResults
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

// UpdateChannel updates the specified CloudTrail channel.
func (s *CloudTrailService) UpdateChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	arn := request.GetStringParam(req.Parameters, "Channel")
	if arn == "" {
		return nil, awserrors.NewAWSError("InvalidParameter", "Channel is required", 400)
	}

	ch, err := store.GetChannel(arn)
	if err != nil {
		return nil, awserrors.NewAWSError("ChannelNotFoundException", "Channel not found", 404)
	}

	if name := request.GetStringParam(req.Parameters, "Name"); name != "" {
		ch.Name = name
	}
	if destsRaw, ok := req.Parameters["Destinations"]; ok {
		ch.Destinations = parseDestinations(destsRaw)
	}

	if err := store.UpdateChannel(ch); err != nil {
		return nil, err
	}

	return formatChannel(ch), nil
}
