package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
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

	if gc.DeletionProtection {
		return nil, awserrors.NewAWSError("InvalidGlobalClusterStateFault", "Cannot delete global cluster when DeletionProtection is enabled", http.StatusBadRequest)
	}

	gc.Status = "deleting"

	if err := store.DeleteGlobalCluster(id); err != nil {
		return nil, translateStoreError(err)
	}

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

	items := make([]interface{}, 0, len(clusters))
	for _, c := range clusters {
		items = append(items, c)
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(*neptunestore.GlobalCluster).GlobalClusterIdentifier
	})

	result := map[string]interface{}{
		"GlobalClusters": protocol.XMLElements{ElementName: "GlobalClusterMember", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
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
		oldID := gc.GlobalClusterIdentifier
		gc.GlobalClusterIdentifier = v
		if err := store.CreateGlobalCluster(gc); err != nil {
			return nil, translateStoreError(err)
		}
		if err := store.DeleteGlobalCluster(oldID); err != nil {
			store.DeleteGlobalCluster(v)
			return nil, translateStoreError(err)
		}
		return map[string]interface{}{
			"GlobalCluster": gc,
		}, nil
	}
	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		gc.EngineVersion = v
	}
	if request.HasParam(params, "DeletionProtection") {
		gc.DeletionProtection = request.GetBoolParam(params, "DeletionProtection")
	}

	if err := store.UpdateGlobalCluster(gc); err != nil {
		return nil, translateStoreError(err)
	}

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
	if err := store.UpdateGlobalCluster(gc); err != nil {
		return nil, translateStoreError(err)
	}

	gc.Status = "available"
	if err := store.UpdateGlobalCluster(gc); err != nil {
		return nil, translateStoreError(err)
	}

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

	targetID := request.GetStringParam(params, "TargetDbClusterIdentifier")
	if targetID == "" {
		return nil, fmt.Errorf("neptune: TargetDbClusterIdentifier is required")
	}
	found := false
	for i := range gc.GlobalClusterMembers {
		if gc.GlobalClusterMembers[i].IsWriter {
			gc.GlobalClusterMembers[i].IsWriter = false
		}
		if gc.GlobalClusterMembers[i].DBClusterArn == targetID || gc.GlobalClusterMembers[i].GlobalClusterIdentifier == targetID {
			gc.GlobalClusterMembers[i].IsWriter = true
			found = true
		}
	}
	if !found {
		return nil, fmt.Errorf("neptune: TargetDbClusterIdentifier %s not found in global cluster members", targetID)
	}
	gc.Status = "available"
	if err := store.UpdateGlobalCluster(gc); err != nil {
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

	clusterID := request.GetStringParam(params, "DbClusterIdentifier")
	if clusterID == "" {
		return nil, fmt.Errorf("neptune: DbClusterIdentifier is required")
	}
	found := false
	for _, m := range gc.GlobalClusterMembers {
		if m.DBClusterArn == clusterID || m.GlobalClusterIdentifier == clusterID {
			if m.IsWriter {
				return nil, fmt.Errorf("neptune: cannot remove the writer member from global cluster")
			}
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("neptune: member %s not found in global cluster %s", clusterID, id)
	}
	filtered := make([]neptunestore.GlobalClusterMember, 0, len(gc.GlobalClusterMembers))
	for _, m := range gc.GlobalClusterMembers {
		if m.DBClusterArn != clusterID && m.GlobalClusterIdentifier != clusterID {
			filtered = append(filtered, m)
		}
	}
	gc.GlobalClusterMembers = filtered
	if err := store.UpdateGlobalCluster(gc); err != nil {
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

	sub.Status = "deleted"

	if err := store.DeleteEventSubscription(name); err != nil {
		return nil, translateStoreError(err)
	}

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

	items := make([]interface{}, 0, len(subs))
	for _, sub := range subs {
		items = append(items, sub)
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		return item.(*neptunestore.EventSubscription).CustSubscriptionId
	})

	result := map[string]interface{}{
		"EventSubscriptionsList": protocol.XMLElements{ElementName: "EventSubscription", Items: resultItems},
	}
	if isTruncated {
		result["Marker"] = nextMarker
	}
	return result, nil
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
	if request.HasParam(params, "EventCategories") {
		sub.EventCategoriesList = request.GetStringList(params, "EventCategories")
	}
	if request.HasParam(params, "Enabled") {
		sub.Enabled = request.GetBoolParam(params, "Enabled")
	}

	if err := store.UpdateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

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

	for _, existing := range sub.SourceIdsList {
		if existing == sourceID {
			return nil, fmt.Errorf("neptune: source identifier %s already exists in subscription %s", sourceID, name)
		}
	}
	sub.SourceIdsList = append(sub.SourceIdsList, sourceID)
	if err := store.UpdateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

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

	filtered := make([]string, 0, len(sub.SourceIdsList))
	for _, id := range sub.SourceIdsList {
		if id != sourceID {
			filtered = append(filtered, id)
		}
	}
	sub.SourceIdsList = filtered
	if err := store.UpdateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"EventSubscription": sub,
	}, nil
}
