package neptune

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/protocol"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// clusterIDFromARN extracts the cluster identifier from an RDS cluster ARN.
// Uses SplitN to correctly handle identifiers that may contain colons (L7 fix).
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
	parts := strings.Split(arn, ":")
	if len(parts) < 6 {
		return false
	}
	return parts[0] == "arn" && (parts[1] == "aws" || parts[1] == "aws-cn" || parts[1] == "aws-us-gov") && parts[2] == "sns"
}

// isValidEventCategory checks whether the given category is one of the
// known Neptune event categories (M15 fix).  The list mirrors the output
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
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
	}
	engine := request.GetStringParam(params, "Engine")
	if engine == "" {
		engine = "neptune"
	}
	if err := rdssvc.ValidateEngine(engine); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	engineVersion := request.GetStringParam(params, "EngineVersion")
	if engineVersion == "" {
		engineVersion = rdssvc.DefaultEngineVersion(engine)
	}
	if err := rdssvc.ValidateEngineVersion(engine, engineVersion); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc := &neptunestore.GlobalCluster{
		GlobalClusterIdentifier: id,
		GlobalClusterResourceId: fmt.Sprintf("cluster-%s", id),
		GlobalClusterArn:        arnutil.NewARNBuilder(reqCtx.GetAccountID(), reqCtx.GetRegion()).RDS().GlobalCluster(id),
		Engine:                  engine,
		EngineVersion:           engineVersion,
		Status:                  "creating",
		StorageEncrypted:        request.GetBoolParam(params, "StorageEncrypted"),
		DeletionProtection:      request.GetBoolParam(params, "DeletionProtection"),
		AccountID:               reqCtx.GetAccountID(),
		Region:                  reqCtx.GetRegion(),
	}

	if err := store.CreateGlobalCluster(gc); err != nil {
		return nil, translateStoreError(err)
	}

	// State machine: synchronous transition from 'creating' to 'available'
	// with safety-net goroutine.
	gc.Status = "available"
	if err := store.UpdateGlobalCluster(gc); err != nil {
		logs.Warn("failed to transition global cluster to available", logs.String("gc", id), logs.Err(err))
	}
	s.scheduleTransition(reqCtx.GetRegion(), 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
		g, err := st.GetGlobalCluster(id)
		if err != nil || g.Status != "creating" {
			return nil
		}
		g.Status = "available"
		return st.UpdateGlobalCluster(g)
	})

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// DeleteGlobalCluster deletes the specified Neptune global cluster.
func (s *NeptuneService) DeleteGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
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
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := request.GetStringParam(params, "EngineVersion"); v != "" {
		if err := rdssvc.ValidateEngineVersion(gc.Engine, v); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
		gc.EngineVersion = v
	}
	if request.HasParam(params, "DeletionProtection") {
		gc.DeletionProtection = request.GetBoolParam(params, "DeletionProtection")
	}

	if err := store.UpdateGlobalCluster(gc); err != nil {
		return nil, translateStoreError(err)
	}

	newID := request.GetStringParam(params, "NewGlobalClusterIdentifier")
	if newID != "" && newID != id {
		oldID := gc.GlobalClusterIdentifier
		gc.GlobalClusterIdentifier = newID
		gc.GlobalClusterArn = arnutil.NewARNBuilder(gc.AccountID, gc.Region).RDS().GlobalCluster(newID)
		if err := store.CreateGlobalCluster(gc); err != nil {
			return nil, translateStoreError(err)
		}

		// M10: Delete old entry before updating member clusters.  If
		// DeleteGlobalCluster fails, clean up the new entry and return
		// error; member clusters still reference oldID which remains in
		// the store, preserving referential integrity.
		if err := store.DeleteGlobalCluster(oldID); err != nil {
			_ = store.DeleteGlobalCluster(newID)
			return nil, translateStoreError(err)
		}

		// Update GlobalClusterIdentifier on all member clusters so they
		// reference the renamed global cluster. Without this, member
		// clusters retain a stale GC identifier that no longer resolves.
		for _, m := range gc.GlobalClusterMembers {
			clusterID := clusterIDFromARN(m.DBClusterArn)
			if member, err := store.GetCluster(clusterID); err == nil {
				member.GlobalClusterIdentifier = newID
				if err := store.UpdateCluster(member); err != nil {
					logs.Warn("rename GC: failed to update member cluster GC ref",
						logs.String("cluster", clusterID), logs.Err(err))
				}
			}
		}
	}

	return map[string]interface{}{
		"GlobalCluster": gc,
	}, nil
}

// FailoverGlobalCluster promotes the specified secondary DB cluster to be the
// primary DB cluster in the global database cluster.
func (s *NeptuneService) FailoverGlobalCluster(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	id := request.GetStringParam(params, "GlobalClusterIdentifier")
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
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
		return nil, awserrors.NewMissingParameter("TargetDbClusterIdentifier is required")
	}
	found := false
	for i := range gc.GlobalClusterMembers {
		if clusterIDFromARN(gc.GlobalClusterMembers[i].DBClusterArn) == targetID || gc.GlobalClusterMembers[i].DBClusterArn == targetID {
			if gc.GlobalClusterMembers[i].IsWriter {
				return nil, awserrors.NewAWSError("InvalidGlobalClusterStateFault", fmt.Sprintf("cluster %s is already the writer", targetID), http.StatusBadRequest)
			}
			for j := range gc.GlobalClusterMembers {
				gc.GlobalClusterMembers[j].IsWriter = j == i
			}
			found = true
			break
		}
	}
	if !found {
		return nil, awserrors.NewAWSError("GlobalClusterMemberNotFoundFault", fmt.Sprintf("cluster %s not found in global cluster %s", targetID, id), http.StatusNotFound)
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
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
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
		return nil, awserrors.NewMissingParameter("TargetDbClusterIdentifier is required")
	}
	found := false
	for i := range gc.GlobalClusterMembers {
		if gc.GlobalClusterMembers[i].IsWriter {
			gc.GlobalClusterMembers[i].IsWriter = false
		}
		if clusterIDFromARN(gc.GlobalClusterMembers[i].DBClusterArn) == targetID || gc.GlobalClusterMembers[i].DBClusterArn == targetID || gc.GlobalClusterMembers[i].GlobalClusterIdentifier == targetID {
			gc.GlobalClusterMembers[i].IsWriter = true
			found = true
		}
	}
	if !found {
		return nil, awserrors.NewAWSError("GlobalClusterNotFoundFault", fmt.Sprintf("TargetDbClusterIdentifier %s not found in global cluster members", targetID), http.StatusNotFound)
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
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
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
		return nil, awserrors.NewMissingParameter("DbClusterIdentifier is required")
	}
	found := false
	for _, m := range gc.GlobalClusterMembers {
		if clusterIDFromARN(m.DBClusterArn) == clusterID || m.DBClusterArn == clusterID || m.GlobalClusterIdentifier == clusterID {
			if m.IsWriter {
				return nil, awserrors.NewAWSError("InvalidGlobalClusterStateFault", "Cannot remove the writer member from global cluster", http.StatusBadRequest)
			}
			found = true
			break
		}
	}
	if !found {
		return nil, awserrors.NewAWSError("GlobalClusterNotFoundFault", fmt.Sprintf("Member %s not found in global cluster %s", clusterID, id), http.StatusNotFound)
	}
	filtered := make([]neptunestore.GlobalClusterMember, 0, len(gc.GlobalClusterMembers))
	for _, m := range gc.GlobalClusterMembers {
		if clusterIDFromARN(m.DBClusterArn) != clusterID && m.DBClusterArn != clusterID && m.GlobalClusterIdentifier != clusterID {
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
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
	}
	topicArn := request.GetStringParam(params, "SnsTopicArn")
	if topicArn == "" {
		return nil, awserrors.NewMissingParameter("SnsTopicArn is required")
	}
	// H3: Validate SNS topic ARN format (arn:aws:sns:region:account:topic).
	if !isValidSnsTopicArn(topicArn) {
		return nil, awserrors.NewAWSError("SNSInvalidTopicFault", fmt.Sprintf("Invalid SNS topic ARN: %s", topicArn), http.StatusBadRequest)
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	sub := &neptunestore.EventSubscription{
		CustSubscriptionId:       name,
		SnsTopicArn:              topicArn,
		Status:                   "creating",
		SubscriptionCreationTime: &now,
		SourceType:               request.GetStringParam(params, "SourceType"),
		Enabled:                  request.GetBoolParam(params, "Enabled"),
		CustSubscriptionArn:      neptunestore.EventSubscriptionARN(reqCtx.GetAccountID(), reqCtx.GetRegion(), name),
	}

	if sourceIds := request.GetStringList(params, "SourceIds"); len(sourceIds) > 0 {
		sub.SourceIdsList = sourceIds
	}
	if categories := request.GetStringList(params, "EventCategories"); len(categories) > 0 {
		// M15: Validate EventCategories against the known Neptune categories.
		for _, cat := range categories {
			if !isValidEventCategory(cat) {
				return nil, awserrors.NewAWSError("SubscriptionCategoryNotFoundFault",
					fmt.Sprintf("Invalid event category: %s", cat), http.StatusBadRequest)
			}
		}
		sub.EventCategoriesList = categories
	}

	if err := store.CreateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

	// M8: State machine — synchronous transition from 'creating' to 'active'
	// with safety-net goroutine.
	sub.Status = "active"
	if err := store.UpdateEventSubscription(sub); err != nil {
		logs.Warn("failed to transition event subscription to active", logs.String("sub", name), logs.Err(err))
	}
	s.scheduleTransition(reqCtx.GetRegion(), 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
		es, err := st.GetEventSubscription(name)
		if err != nil || es.Status != "creating" {
			return nil
		}
		es.Status = "active"
		return st.UpdateEventSubscription(es)
	})

	return map[string]interface{}{
		"EventSubscription": enrichEventSubscription(sub),
	}, nil
}

// DeleteEventSubscription deletes the specified event subscription.
func (s *NeptuneService) DeleteEventSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
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
		"EventSubscription": enrichEventSubscription(sub),
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
			"EventSubscriptionsList": protocol.XMLElements{ElementName: "EventSubscription", Items: []interface{}{enrichEventSubscription(sub)}},
		}, nil
	}

	subs, err := store.ListEventSubscriptions()
	if err != nil {
		return nil, translateStoreError(err)
	}

	items := make([]interface{}, 0, len(subs))
	for _, sub := range subs {
		items = append(items, enrichEventSubscription(sub))
	}

	marker := request.GetStringParam(params, "Marker")
	maxRecords := request.GetIntParam(params, "MaxRecords")
	resultItems, nextMarker, isTruncated := paginateItems(items, marker, maxRecords, func(item interface{}) string {
		m := item.(map[string]interface{})
		if v, ok := m["CustSubscriptionId"]; ok {
			return v.(string)
		}
		return ""
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
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
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
		categories := request.GetStringList(params, "EventCategories")
		// M15: Validate EventCategories on modify.
		for _, cat := range categories {
			if !isValidEventCategory(cat) {
				return nil, awserrors.NewAWSError("SubscriptionCategoryNotFoundFault",
					fmt.Sprintf("Invalid event category: %s", cat), http.StatusBadRequest)
			}
		}
		sub.EventCategoriesList = categories
	}
	if request.HasParam(params, "Enabled") {
		sub.Enabled = request.GetBoolParam(params, "Enabled")
	}

	if err := store.UpdateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"EventSubscription": enrichEventSubscription(sub),
	}, nil
}

// AddSourceIdentifierToSubscription adds a source identifier to the specified
// event subscription.
func (s *NeptuneService) AddSourceIdentifierToSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
	}
	sourceID := request.GetStringParam(params, "SourceIdentifier")
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceIdentifier is required")
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
			return nil, awserrors.NewAWSError("SourceNotFoundFault", fmt.Sprintf("Source identifier %s already exists in subscription %s", sourceID, name), http.StatusConflict)
		}
	}
	sub.SourceIdsList = append(sub.SourceIdsList, sourceID)
	if err := store.UpdateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"EventSubscription": enrichEventSubscription(sub),
	}, nil
}

// RemoveSourceIdentifierFromSubscription removes a source identifier from the
// specified event subscription.
func (s *NeptuneService) RemoveSourceIdentifierFromSubscription(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	params := req.Parameters
	name := request.GetStringParam(params, "SubscriptionName")
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
	}
	sourceID := request.GetStringParam(params, "SourceIdentifier")
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceIdentifier is required")
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
		"EventSubscription": enrichEventSubscription(sub),
	}, nil
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
