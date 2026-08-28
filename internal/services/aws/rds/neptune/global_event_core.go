package neptune

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awserrors "vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// CreateGlobalClusterInput carries the wire-parsed CreateGlobalCluster
// request. EngineVersion carries the raw wire value; the Core applies the
// engine default when it is empty.
type CreateGlobalClusterInput struct {
	GlobalClusterIdentifier string
	Engine                  string
	EngineVersion           string
	StorageEncrypted        bool
	DeletionProtection      bool
	AccountID               string
	Region                  string
}

// DeleteGlobalClusterInput carries the wire-parsed DeleteGlobalCluster
// request.
type DeleteGlobalClusterInput struct {
	GlobalClusterIdentifier string
}

// ModifyGlobalClusterInput carries the wire-parsed ModifyGlobalCluster
// request.
type ModifyGlobalClusterInput struct {
	GlobalClusterIdentifier    string
	EngineVersion              string
	HasDeletionProtection      bool
	DeletionProtection         bool
	NewGlobalClusterIdentifier string
}

// FailoverGlobalClusterInput carries the wire-parsed FailoverGlobalCluster
// request.
type FailoverGlobalClusterInput struct {
	GlobalClusterIdentifier   string
	TargetDbClusterIdentifier string
}

// SwitchoverGlobalClusterInput carries the wire-parsed
// SwitchoverGlobalCluster request.
type SwitchoverGlobalClusterInput struct {
	GlobalClusterIdentifier   string
	TargetDbClusterIdentifier string
}

// RemoveFromGlobalClusterInput carries the wire-parsed
// RemoveFromGlobalCluster request.
type RemoveFromGlobalClusterInput struct {
	GlobalClusterIdentifier string
	DbClusterIdentifier     string
}

// CreateEventSubscriptionInput carries the wire-parsed CreateEventSubscription
// request.
type CreateEventSubscriptionInput struct {
	SubscriptionName string
	SnsTopicArn      string
	SourceType       string
	Enabled          bool
	SourceIds        []string
	EventCategories  []string
	AccountID        string
	Region           string
}

// DeleteEventSubscriptionInput carries the wire-parsed DeleteEventSubscription
// request.
type DeleteEventSubscriptionInput struct {
	SubscriptionName string
}

// ModifyEventSubscriptionInput carries the wire-parsed ModifyEventSubscription
// request. The Has* members preserve the wire presence of optional SCALAR
// members; the EventCategories list is value-gated because query-protocol
// lists arrive as flattened member keys.
type ModifyEventSubscriptionInput struct {
	SubscriptionName string
	SnsTopicArn      string
	SourceType       string
	EventCategories  []string
	HasEnabled       bool
	Enabled          bool
}

// AddSourceIdentifierToSubscriptionInput carries the wire-parsed
// AddSourceIdentifierToSubscription request.
type AddSourceIdentifierToSubscriptionInput struct {
	SubscriptionName string
	SourceIdentifier string
}

// RemoveSourceIdentifierFromSubscriptionInput carries the wire-parsed
// RemoveSourceIdentifierFromSubscription request.
type RemoveSourceIdentifierFromSubscriptionInput struct {
	SubscriptionName string
	SourceIdentifier string
}

// createGlobalClusterCore validates and persists a new global cluster, then
// transitions it to available with a safety-net goroutine.
func (s *NeptuneService) createGlobalClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateGlobalClusterInput) (interface{}, error) {
	id := in.GlobalClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
	}
	engine := in.Engine
	if engine == "" {
		return nil, awserrors.NewMissingParameter("Engine is required")
	}
	if err := rdssvc.ValidateEngine(engine); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	engineVersion := in.EngineVersion
	if engineVersion == "" {
		engineVersion = rdssvc.DefaultEngineVersion(engine)
	}
	if err := rdssvc.ValidateEngineVersion(engine, engineVersion); err != nil {
		return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	gc := &neptunestore.GlobalCluster{
		GlobalClusterIdentifier: id,
		GlobalClusterResourceId: fmt.Sprintf("cluster-%s", id),
		GlobalClusterArn:        arnutil.NewARNBuilder(in.AccountID, in.Region).RDS().GlobalCluster(id),
		Engine:                  engine,
		EngineVersion:           engineVersion,
		Status:                  "creating",
		StorageEncrypted:        in.StorageEncrypted,
		DeletionProtection:      in.DeletionProtection,
		AccountID:               in.AccountID,
		Region:                  in.Region,
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
	s.scheduleTransition(in.Region, 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
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

// deleteGlobalClusterCore deletes a global cluster after the
// DeletionProtection check.
func (s *NeptuneService) deleteGlobalClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteGlobalClusterInput) (interface{}, error) {
	id := in.GlobalClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
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

// modifyGlobalClusterCore applies member updates to a global cluster and
// handles the rename path, repointing member clusters at the new identifier.
func (s *NeptuneService) modifyGlobalClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyGlobalClusterInput) (interface{}, error) {
	id := in.GlobalClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if v := in.EngineVersion; v != "" {
		if err := rdssvc.ValidateEngineVersion(gc.Engine, v); err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
		}
		gc.EngineVersion = v
	}
	if in.HasDeletionProtection {
		gc.DeletionProtection = in.DeletionProtection
	}

	if err := store.UpdateGlobalCluster(gc); err != nil {
		return nil, translateStoreError(err)
	}

	newID := in.NewGlobalClusterIdentifier
	if newID != "" && newID != id {
		oldID := gc.GlobalClusterIdentifier
		gc.GlobalClusterIdentifier = newID
		gc.GlobalClusterArn = arnutil.NewARNBuilder(gc.AccountID, gc.Region).RDS().GlobalCluster(newID)
		if err := store.CreateGlobalCluster(gc); err != nil {
			return nil, translateStoreError(err)
		}

		// Delete old entry before updating member clusters.  If
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

// failoverGlobalClusterCore promotes the target secondary member to writer.
func (s *NeptuneService) failoverGlobalClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *FailoverGlobalClusterInput) (interface{}, error) {
	id := in.GlobalClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	targetID := in.TargetDbClusterIdentifier
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

// switchoverGlobalClusterCore switches the writer member of a global
// cluster to the specified target.
func (s *NeptuneService) switchoverGlobalClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *SwitchoverGlobalClusterInput) (interface{}, error) {
	id := in.GlobalClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	targetID := in.TargetDbClusterIdentifier
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

// removeFromGlobalClusterCore detaches a member cluster from its global
// cluster.
func (s *NeptuneService) removeFromGlobalClusterCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *RemoveFromGlobalClusterInput) (interface{}, error) {
	id := in.GlobalClusterIdentifier
	if id == "" {
		return nil, awserrors.NewMissingParameter("GlobalClusterIdentifier is required")
	}

	gc, err := store.GetGlobalCluster(id)
	if err != nil {
		return nil, translateStoreError(err)
	}

	clusterID := in.DbClusterIdentifier
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

// createEventSubscriptionCore validates and persists a new event
// subscription, then transitions it to active with a safety-net goroutine.
func (s *NeptuneService) createEventSubscriptionCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateEventSubscriptionInput) (interface{}, error) {
	name := in.SubscriptionName
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
	}
	topicArn := in.SnsTopicArn
	if topicArn == "" {
		return nil, awserrors.NewMissingParameter("SnsTopicArn is required")
	}
	// Validate SNS topic ARN format (arn:aws:sns:region:account:topic).
	if !isValidSnsTopicArn(topicArn) {
		return nil, awserrors.NewAWSError("SNSInvalidTopicFault", fmt.Sprintf("Invalid SNS topic ARN: %s", topicArn), http.StatusBadRequest)
	}

	now := time.Now()
	sub := &neptunestore.EventSubscription{
		CustSubscriptionId:       name,
		SnsTopicArn:              topicArn,
		Status:                   "creating",
		SubscriptionCreationTime: &now,
		SourceType:               in.SourceType,
		Enabled:                  in.Enabled,
		CustSubscriptionArn:      neptunestore.EventSubscriptionARN(in.AccountID, in.Region, name),
	}

	if len(in.SourceIds) > 0 {
		sub.SourceIdsList = in.SourceIds
	}
	if len(in.EventCategories) > 0 {
		// Validate EventCategories against the known Neptune categories.
		for _, cat := range in.EventCategories {
			if !isValidEventCategory(cat) {
				return nil, awserrors.NewAWSError("SubscriptionCategoryNotFoundFault",
					fmt.Sprintf("Invalid event category: %s", cat), http.StatusBadRequest)
			}
		}
		sub.EventCategoriesList = in.EventCategories
	}

	if err := store.CreateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

	// State machine — synchronous transition from 'creating' to 'active'
	// with safety-net goroutine.
	sub.Status = "active"
	if err := store.UpdateEventSubscription(sub); err != nil {
		logs.Warn("failed to transition event subscription to active", logs.String("sub", name), logs.Err(err))
	}
	s.scheduleTransition(in.Region, 500*time.Millisecond, func(st neptunestore.NeptuneStoreInterface) error {
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

// deleteEventSubscriptionCore deletes an event subscription.
func (s *NeptuneService) deleteEventSubscriptionCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteEventSubscriptionInput) (interface{}, error) {
	name := in.SubscriptionName
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
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

// modifyEventSubscriptionCore applies member updates to an event
// subscription.
func (s *NeptuneService) modifyEventSubscriptionCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyEventSubscriptionInput) (interface{}, error) {
	name := in.SubscriptionName
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
	}

	sub, err := store.GetEventSubscription(name)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if topicArn := in.SnsTopicArn; topicArn != "" {
		sub.SnsTopicArn = topicArn
	}
	if sourceType := in.SourceType; sourceType != "" {
		sub.SourceType = sourceType
	}
	if len(in.EventCategories) > 0 {
		// Validate EventCategories on modify.
		for _, cat := range in.EventCategories {
			if !isValidEventCategory(cat) {
				return nil, awserrors.NewAWSError("SubscriptionCategoryNotFoundFault",
					fmt.Sprintf("Invalid event category: %s", cat), http.StatusBadRequest)
			}
		}
		sub.EventCategoriesList = in.EventCategories
	}
	if in.HasEnabled {
		sub.Enabled = in.Enabled
	}

	if err := store.UpdateEventSubscription(sub); err != nil {
		return nil, translateStoreError(err)
	}

	return map[string]interface{}{
		"EventSubscription": enrichEventSubscription(sub),
	}, nil
}

// addSourceIdentifierToSubscriptionCore adds a source identifier to an
// event subscription.
func (s *NeptuneService) addSourceIdentifierToSubscriptionCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *AddSourceIdentifierToSubscriptionInput) (interface{}, error) {
	name := in.SubscriptionName
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
	}
	sourceID := in.SourceIdentifier
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceIdentifier is required")
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

// removeSourceIdentifierFromSubscriptionCore removes a source identifier
// from an event subscription.
func (s *NeptuneService) removeSourceIdentifierFromSubscriptionCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *RemoveSourceIdentifierFromSubscriptionInput) (interface{}, error) {
	name := in.SubscriptionName
	if name == "" {
		return nil, awserrors.NewMissingParameter("SubscriptionName is required")
	}
	sourceID := in.SourceIdentifier
	if sourceID == "" {
		return nil, awserrors.NewMissingParameter("SourceIdentifier is required")
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
