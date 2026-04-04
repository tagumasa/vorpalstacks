package neptune

import (
	"context"
	"fmt"
	"time"

	"vorpalstacks/internal/services/aws/common/protocol"
	"vorpalstacks/internal/services/aws/common/request"
	neptunestore "vorpalstacks/internal/store/aws/neptune"
)

// CreateGlobalCluster creates a new Neptune global cluster.
func (s *NeptuneService) CreateGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: GlobalClusterIdentifier is required")
	}
	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc := &neptunestore.GlobalCluster{
		GlobalClusterIdentifier: id,
		GlobalClusterResourceId: fmt.Sprintf("cluster-%s", id),
		GlobalClusterArn:        fmt.Sprintf("arn:aws:rds:%s:%s:global-cluster:%s", reqCtx.GetRegion(), reqCtx.GetAccountID(), id),
		Engine:                  engine,
		EngineVersion:           request.GetStringParam(params, "EngineVersion"),
		Status:                  "available",
		StorageEncrypted:        request.GetBoolParam(params, "StorageEncrypted"),
		DeletionProtection:      request.GetBoolParam(params, "DeletionProtection"),
		AccountID:               reqCtx.GetAccountID(),
		Region:                  reqCtx.GetRegion(),
	}

	if err := store.CreateGlobalCluster(gc); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// DeleteGlobalCluster deletes the specified Neptune global cluster.
func (s *NeptuneService) DeleteGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: GlobalClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	_ = store.DeleteGlobalCluster(id)
	gc.Status = "deleting"

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// DescribeGlobalClusters returns information about the specified global cluster
// or lists all global clusters when no identifier is provided.
func (s *NeptuneService) DescribeGlobalClusters(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id != "" {
		gc, err := store.GetGlobalCluster(id)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"GlobalClusters": protocol.XMLElements{ElementName: "GlobalClusterMember", Items: []interface{}{gc}},
		}, nil
	}

	clusters, err := store.ListGlobalClusters()
	if err != nil {
		return nil, translateStoreError(err)
	}

	result := make([]interface{}, 0, len(clusters))
	for _, c := range clusters {
		result = append(result, c)
	}

	return map[string]interface{}{
		"GlobalClusters": protocol.XMLElements{ElementName: "GlobalClusterMember", Items: result},
	}, nil
}

// ModifyGlobalCluster updates the configuration of the specified global cluster.
func (s *NeptuneService) ModifyGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: GlobalClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := request.GetStringParam(params, "NewGlobalClusterIdentifier"); v != "" {
		gc.GlobalClusterIdentifier = v
	}
	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		gc.EngineVersion = v
	}
	if request.GetBoolParam(params, "DeletionProtection") {
		gc.DeletionProtection = true
	}

	_ = store.UpdateGlobalCluster(gc)

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// FailoverGlobalCluster forces a failover for the specified global cluster.
func (s *NeptuneService) FailoverGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: GlobalClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	gc.Status = "failing-over"
	_ = store.UpdateGlobalCluster(gc)
	gc.Status = "available"
	_ = store.UpdateGlobalCluster(gc)

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// SwitchoverGlobalCluster switches the primary cluster for the specified
// global cluster.
func (s *NeptuneService) SwitchoverGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: GlobalClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// RemoveFromGlobalCluster detaches a cluster from the specified global cluster.
func (s *NeptuneService) RemoveFromGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, fmt.Errorf("neptune: GlobalClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// CreateEventSubscription creates a new Neptune event subscription that
// publishes events to an SNS topic.
func (s *NeptuneService) CreateEventSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, fmt.Errorf("neptune: SubscriptionName is required")
	}
	topicArn := request.GetStringParam(params, "SnsTopicArn")
	if topicArn == "" {
		return nil, fmt.Errorf("neptune: SnsTopicArn is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sub := &neptunestore.EventSubscription{
		CustSubscriptionId:       name,
		SnsTopicArn:              topicArn,
		Status:                   "active",
		SubscriptionCreationTime: &now,
		SourceType:               request.GetStringParam(params, "SourceType"),
		Enabled:                  request.GetBoolParam(params, "Enabled"),
		CustSubscriptionArn:      neptunestore.EventSubscriptionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}

	if sourceIds := request.GetStringList(params, "SourceIds"); len(sourceIds) > 0 {
		sub.SourceIdsList = sourceIds
	}
	if categories := request.GetStringList(params, "EventCategories"); len(categories) > 0 {
		sub.EventCategoriesList = categories
	}

	if err := store.CreateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"EventSubscription": sub,
	}, nil
}

// DeleteEventSubscription deletes the specified event subscription.
func (s *NeptuneService) DeleteEventSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, fmt.Errorf("neptune: SubscriptionName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sub, err := store.GetEventSubscription(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	_ = store.DeleteEventSubscription(name)
	sub.Status = "deleted"

	return map[string]interface{}{
		"EventSubscription": sub,
	}, nil
}

// DescribeEventSubscriptions returns information about the specified event
// subscription or lists all subscriptions when no name is provided.
func (s *NeptuneService) DescribeEventSubscriptions(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	name := request.GetStringParam(params, "SubscriptionName")
	if name != "" {
		sub, err := store.GetEventSubscription(name)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"EventSubscriptionsList": protocol.XMLElements{ElementName: "EventSubscription", Items: []interface{}{sub}},
		}, nil
	}

	subs, err := store.ListEventSubscriptions()
	if err != nil {
		return nil, translateStoreError(err)
	}

	result := make([]interface{}, 0, len(subs))
	for _, sub := range subs {
		result = append(result, sub)
	}

	return map[string]interface{}{
		"EventSubscriptionsList": protocol.XMLElements{ElementName: "EventSubscription", Items: result},
	}, nil
}

// ModifyEventSubscription updates the configuration of the specified event
// subscription.
func (s *NeptuneService) ModifyEventSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, fmt.Errorf("neptune: SubscriptionName is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sub, err := store.GetEventSubscription(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if topicArn := request.GetStringParam(params, "SnsTopicArn"); topicArn != "" {
		sub.SnsTopicArn = topicArn
	}
	if sourceType := request.GetStringParam(params, "SourceType"); sourceType != "" {
		sub.SourceType = sourceType
	}

	_ = store.UpdateEventSubscription(sub)

	return map[string]interface{}{
		"EventSubscription": sub,
	}, nil
}

// AddSourceIdentifierToSubscription adds a source identifier to the specified
// event subscription.
func (s *NeptuneService) AddSourceIdentifierToSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, fmt.Errorf("neptune: SubscriptionName is required")
	}
	sourceID := request.GetStringParam(params, "SourceIdentifier")
	if sourceID == "" {
		return nil, fmt.Errorf("neptune: SourceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sub, err := store.GetEventSubscription(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	sub.SourceIdsList = append(sub.SourceIdsList, sourceID)
	_ = store.UpdateEventSubscription(sub)

	return map[string]interface{}{
		"EventSubscription": sub,
	}, nil
}

// RemoveSourceIdentifierFromSubscription removes a source identifier from the
// specified event subscription.
func (s *NeptuneService) RemoveSourceIdentifierFromSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, fmt.Errorf("neptune: SubscriptionName is required")
	}
	sourceID := request.GetStringParam(params, "SourceIdentifier")
	if sourceID == "" {
		return nil, fmt.Errorf("neptune: SourceIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	sub, err := store.GetEventSubscription(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	filtered := sub.SourceIdsList[:0]
	for _, id := range sub.SourceIdsList {
		if id != sourceID {
			filtered = append(filtered, id)
		}
	}
	sub.SourceIdsList = filtered
	_ = store.UpdateEventSubscription(sub)

	return map[string]interface{}{
		"EventSubscription": sub,
	}, nil
}
