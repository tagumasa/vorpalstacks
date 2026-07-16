package cloudtrail

import (
	"context"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
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
	return resp
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

// ListChannels lists all CloudTrail channels.
func (s *CloudTrailService) ListChannels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	channels, err := store.ListAllChannels()
	if err != nil {
		return nil, err
	}

	chList := make([]map[string]interface{}, 0, len(channels))
	for _, ch := range channels {
		chList = append(chList, formatChannel(ch))
	}

	result := map[string]interface{}{
		"Channels": chList,
	}

	return result, nil
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
