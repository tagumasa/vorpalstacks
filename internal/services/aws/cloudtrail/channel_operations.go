package cloudtrail

import (
	"context"

	"vorpalstacks/internal/common/request"
	tags "vorpalstacks/internal/common/tags"
)

// CreateChannel creates a new CloudTrail channel.
func (s *CloudTrailService) CreateChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	in := CreateChannelInput{
		Name:   request.GetStringParam(req.Parameters, "Name"),
		Source: request.GetStringParam(req.Parameters, "Source"),
	}
	if destsRaw, ok := req.Parameters["Destinations"]; ok {
		in.DestinationsRaw = destsRaw
		in.DestinationsSet = true
	}
	if tagsRaw, ok := req.Parameters["Tags"]; ok {
		in.TagsRaw = tagsRaw
		in.TagsSet = true
		in.TagList = tags.ParseTags(req.Parameters, "Tags")
	}

	return s.createChannelCore(store, in)
}

// DeleteChannel deletes the specified CloudTrail channel.
func (s *CloudTrailService) DeleteChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.deleteChannelCore(store, ChannelInput{
		Channel: request.GetStringParam(req.Parameters, "Channel"),
	})
}

// GetChannel retrieves the specified CloudTrail channel.
func (s *CloudTrailService) GetChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.getChannelCore(store, ChannelInput{
		Channel: request.GetStringParam(req.Parameters, "Channel"),
	})
}

// ListChannels lists CloudTrail channels with pagination.
func (s *CloudTrailService) ListChannels(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	return s.listChannelsCore(store, ListChannelsInput{
		NextToken:  req.GetParam("NextToken"),
		MaxResults: request.GetIntParam(req.Parameters, "MaxResults"),
	})
}

// UpdateChannel updates the specified CloudTrail channel.
func (s *CloudTrailService) UpdateChannel(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, s.mapStoreError(err)
	}

	in := UpdateChannelInput{
		Channel: request.GetStringParam(req.Parameters, "Channel"),
		Name:    request.GetStringParam(req.Parameters, "Name"),
	}
	if destsRaw, ok := req.Parameters["Destinations"]; ok {
		in.DestinationsRaw = destsRaw
		in.DestinationsSet = true
	}

	return s.updateChannelCore(store, in)
}
