package neptune

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// clusterIDFromARN extracts the cluster identifier from an RDS cluster ARN.
// Uses SplitN to correctly handle identifiers that may contain colons.
func clusterIDFromARN(arn string) string {
	idx := strings.Index(arn, "cluster:")
	if idx < 0 {
		return arn
	}
	return arn[idx+len("cluster:"):]
}

// isValidSnsTopicArn validates that the given string is a well-formed SNS
// topic ARN (arn:aws:sns:<region>:<account>:<topic-name>).
func isValidSnsTopicArn(arn string) bool {
	partition, service, _, _, _ := arnutil.SplitARN(arn)
	return (partition == "aws" || partition == "aws-cn" || partition == "aws-us-gov") && service == "sns"
}

// isValidEventCategory checks whether the given category is one of the
// known Neptune event categories.  The list mirrors the output
// of DescribeEventCategories.
var validEventCategories = map[string]bool{
	"creation": true, "deletion": true, "failover": true,
	"failure": true, "maintenance": true, "notification": true,
	"read replica": true, "recovery": true, "restoration": true,
	"backup": true,
}

func isValidEventCategory(cat string) bool {
	return validEventCategories[cat]
}

// CreateGlobalCluster creates a new Neptune global cluster.
func (s *NeptuneService) CreateGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CreateGlobalClusterInput{
		GlobalClusterIdentifier: request.GetStringParam(params, "GlobalClusterIdentifier"),
		Engine:                  request.GetStringParam(params, "Engine"),
		EngineVersion:           request.GetStringParam(params, "EngineVersion"),
		StorageEncrypted:        request.GetBoolParam(params, "StorageEncrypted"),
		DeletionProtection:      request.GetBoolParam(params, "DeletionProtection"),
		AccountID:               reqCtx.GetAccountID(),
		Region:                  reqCtx.GetRegion(),
	}
	return s.createGlobalClusterCore(ctx, store, in)
}

// DeleteGlobalCluster deletes the specified Neptune global cluster.
func (s *NeptuneService) DeleteGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DeleteGlobalClusterInput{
		GlobalClusterIdentifier: request.GetStringParam(req.Parameters, "GlobalClusterIdentifier"),
	}
	return s.deleteGlobalClusterCore(ctx, store, in)
}

// DescribeGlobalClusters returns information about the specified global cluster
// or lists all global clusters when no identifier is provided.
func (s *NeptuneService) DescribeGlobalClusters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	clusters, nextMarker, err := rdssvc.QueryGlobalClusters(store, rdssvc.DescribeGlobalClustersInput{
		GlobalClusterIdentifier: request.GetStringParam(params, "GlobalClusterIdentifier"),
		Filters:                 nil,
		Marker:                  request.GetStringParam(params, "Marker"),
		MaxRecords:              int32(request.GetIntParam(params, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(clusters))
	for _, gc := range clusters {
		items = append(items, gc)
	}

	result := map[string]interface{}{
		"GlobalClusters": protocol.XMLElements{ElementName: "GlobalClusterMember", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// ModifyGlobalCluster updates the configuration of the specified global cluster.
func (s *NeptuneService) ModifyGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ModifyGlobalClusterInput{
		GlobalClusterIdentifier:    request.GetStringParam(params, "GlobalClusterIdentifier"),
		EngineVersion:              request.GetStringParam(params, "EngineVersion"),
		HasDeletionProtection:      request.HasParam(params, "DeletionProtection"),
		DeletionProtection:         request.GetBoolParam(params, "DeletionProtection"),
		NewGlobalClusterIdentifier: request.GetStringParam(params, "NewGlobalClusterIdentifier"),
	}
	return s.modifyGlobalClusterCore(ctx, store, in)
}

// FailoverGlobalCluster promotes the specified secondary DB cluster to be the
// primary DB cluster in the global database cluster.
func (s *NeptuneService) FailoverGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &FailoverGlobalClusterInput{
		GlobalClusterIdentifier:   request.GetStringParam(params, "GlobalClusterIdentifier"),
		TargetDbClusterIdentifier: request.GetStringParam(params, "TargetDbClusterIdentifier"),
	}
	return s.failoverGlobalClusterCore(ctx, store, in)
}

// SwitchoverGlobalCluster switches the primary cluster for the specified
// global cluster.
func (s *NeptuneService) SwitchoverGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &SwitchoverGlobalClusterInput{
		GlobalClusterIdentifier:   request.GetStringParam(params, "GlobalClusterIdentifier"),
		TargetDbClusterIdentifier: request.GetStringParam(params, "TargetDbClusterIdentifier"),
	}
	return s.switchoverGlobalClusterCore(ctx, store, in)
}

// RemoveFromGlobalCluster detaches a cluster from the specified global cluster.
func (s *NeptuneService) RemoveFromGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &RemoveFromGlobalClusterInput{
		GlobalClusterIdentifier: request.GetStringParam(params, "GlobalClusterIdentifier"),
		DbClusterIdentifier:     request.GetStringParam(params, "DbClusterIdentifier"),
	}
	return s.removeFromGlobalClusterCore(ctx, store, in)
}

// CreateEventSubscription creates a new Neptune event subscription that
// publishes events to an SNS topic.
func (s *NeptuneService) CreateEventSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &CreateEventSubscriptionInput{
		SubscriptionName: request.GetStringParam(params, "SubscriptionName"),
		SnsTopicArn:      request.GetStringParam(params, "SnsTopicArn"),
		SourceType:       request.GetStringParam(params, "SourceType"),
		Enabled:          request.GetBoolParam(params, "Enabled"),
		SourceIds:        getNeptuneStringList(params, "SourceIds", "SourceId", "member"),
		EventCategories:  getNeptuneStringList(params, "EventCategories", "EventCategory", "member"),
		AccountID:        reqCtx.GetAccountID(),
		Region:           reqCtx.GetRegion(),
	}
	return s.createEventSubscriptionCore(ctx, store, in)
}

// DeleteEventSubscription deletes the specified event subscription.
func (s *NeptuneService) DeleteEventSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &DeleteEventSubscriptionInput{
		SubscriptionName: request.GetStringParam(req.Parameters, "SubscriptionName"),
	}
	return s.deleteEventSubscriptionCore(ctx, store, in)
}

// DescribeEventSubscriptions returns information about the specified event
// subscription or lists all subscriptions when no name is provided.
func (s *NeptuneService) DescribeEventSubscriptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	subs, nextMarker, err := rdssvc.QueryEventSubscriptions(store, rdssvc.DescribeEventSubscriptionsInput{
		SubscriptionName: request.GetStringParam(params, "SubscriptionName"),
		Filters:          nil,
		Marker:           request.GetStringParam(params, "Marker"),
		MaxRecords:       int32(request.GetIntParam(params, "MaxRecords")),
	})
	if err != nil {
		return nil, err
	}

	items := make([]interface{}, 0, len(subs))
	for _, sub := range subs {
		items = append(items, enrichEventSubscription(sub))
	}

	result := map[string]interface{}{
		"EventSubscriptionsList": protocol.XMLElements{ElementName: "EventSubscription", Items: items},
	}
	if nextMarker != "" {
		result["Marker"] = nextMarker
	}
	return result, nil
}

// ModifyEventSubscription updates the configuration of the specified event
// subscription.
func (s *NeptuneService) ModifyEventSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &ModifyEventSubscriptionInput{
		SubscriptionName: request.GetStringParam(params, "SubscriptionName"),
		SnsTopicArn:      request.GetStringParam(params, "SnsTopicArn"),
		SourceType:       request.GetStringParam(params, "SourceType"),
		EventCategories:  getNeptuneStringList(params, "EventCategories", "EventCategory", "member"),
		HasEnabled:       request.HasParam(params, "Enabled"),
		Enabled:          request.GetBoolParam(params, "Enabled"),
	}
	return s.modifyEventSubscriptionCore(ctx, store, in)
}

// AddSourceIdentifierToSubscription adds a source identifier to the specified
// event subscription.
func (s *NeptuneService) AddSourceIdentifierToSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &AddSourceIdentifierToSubscriptionInput{
		SubscriptionName: request.GetStringParam(params, "SubscriptionName"),
		SourceIdentifier: request.GetStringParam(params, "SourceIdentifier"),
	}
	return s.addSourceIdentifierToSubscriptionCore(ctx, store, in)
}

// RemoveSourceIdentifierFromSubscription removes a source identifier from the
// specified event subscription.
func (s *NeptuneService) RemoveSourceIdentifierFromSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	in := &RemoveSourceIdentifierFromSubscriptionInput{
		SubscriptionName: request.GetStringParam(params, "SubscriptionName"),
		SourceIdentifier: request.GetStringParam(params, "SourceIdentifier"),
	}
	return s.removeSourceIdentifierFromSubscriptionCore(ctx, store, in)
}

func enrichEventSubscription(sub *neptunestore.EventSubscription) map[string]interface{} {
	m := map[string]interface{}{
		"CustSubscriptionId":   sub.CustSubscriptionId,
		"SnsTopicArn":          sub.SnsTopicArn,
		"Status":               sub.Status,
		"SourceType":           sub.SourceType,
		"Enabled":              sub.Enabled,
		"EventSubscriptionArn": sub.CustSubscriptionArn,
	}
	if sub.SubscriptionCreationTime != nil {
		m["SubscriptionCreationTime"] = *sub.SubscriptionCreationTime
	}
	if len(sub.SourceIdsList) > 0 {
		items := make([]interface{}, 0, len(sub.SourceIdsList))
		for _, id := range sub.SourceIdsList {
			items = append(items, id)
		}
		m["SourceIdsList"] = protocol.XMLElements{ElementName: "SourceId", Items: items}
	}
	if len(sub.EventCategoriesList) > 0 {
		items := make([]interface{}, 0, len(sub.EventCategoriesList))
		for _, cat := range sub.EventCategoriesList {
			items = append(items, cat)
		}
		m["EventCategoriesList"] = protocol.XMLElements{ElementName: "EventCategory", Items: items}
	}
	return m
}
