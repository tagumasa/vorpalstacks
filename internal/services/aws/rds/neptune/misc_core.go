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
	"vorpalstacks/internal/common/tags"
	"vorpalstacks/internal/core/logs"
	rdssvc "vorpalstacks/internal/services/aws/rds"
	neptunestore "vorpalstacks/internal/store/aws/rds/neptune"
	arnutil "vorpalstacks/internal/utils/aws/arn"
)

// CreateDBClusterEndpointInput carries the wire-parsed
// CreateDBClusterEndpoint request.
type CreateDBClusterEndpointInput struct {
	DBClusterEndpointIdentifier string
	DBClusterIdentifier         string
	EndpointType                string
	ExcludedMembers             []string
	StaticMembers               []string
	AccountID                   string
	Region                      string
}

// DescribeDBClusterEndpointsInput carries the wire-parsed
// DescribeDBClusterEndpoints request; an empty DBClusterEndpointIdentifier
// lists every endpoint of the cluster.
type DescribeDBClusterEndpointsInput struct {
	DBClusterIdentifier         string
	DBClusterEndpointIdentifier string
}

// DescribeDBClusterEndpointsResult holds either the single resolved endpoint
// or the full endpoint list, mirroring the two response shapes of the
// operation.
type DescribeDBClusterEndpointsResult struct {
	Endpoint  *neptunestore.DBClusterEndpoint
	Endpoints []*neptunestore.DBClusterEndpoint
}

// ModifyDBClusterEndpointInput carries the wire-parsed
// ModifyDBClusterEndpoint request. Empty member lists leave the stored lists
// unchanged — query-protocol lists cannot carry an explicit clear.
type ModifyDBClusterEndpointInput struct {
	DBClusterEndpointIdentifier string
	EndpointType                string
	ExcludedMembers             []string
	StaticMembers               []string
}

// DeleteDBClusterEndpointInput carries the wire-parsed
// DeleteDBClusterEndpoint request.
type DeleteDBClusterEndpointInput struct {
	DBClusterEndpointIdentifier string
}

// DescribeEventsInput carries the wire-parsed DescribeEvents request. The
// time members stay raw so the Core owns their RFC3339 validation.
type DescribeEventsInput struct {
	SourceType       string
	SourceIdentifier string
	StartTime        string
	EndTime          string
	Duration         int
	Marker           string
	MaxRecords       int
}

// ApplyPendingMaintenanceActionInput carries the wire-parsed
// ApplyPendingMaintenanceAction request.
type ApplyPendingMaintenanceActionInput struct {
	ResourceIdentifier string
	ApplyAction        string
	OptInType          string
}

// DescribeValidDBInstanceModificationsInput carries the wire-parsed
// DescribeValidDBInstanceModifications request.
type DescribeValidDBInstanceModificationsInput struct {
	DBInstanceIdentifier string
}

// recordEvent records a lifecycle event on the store, logging (not failing)
// when the write is rejected.
func recordEvent(store neptunestore.NeptuneStoreInterface, sourceType, sourceID, sourceArn, message string, categories []string) {
	evt := &neptunestore.Event{
		Date:             time.Now().UTC(),
		EventCategories:  categories,
		Message:          message,
		SourceArn:        sourceArn,
		SourceIdentifier: sourceID,
		SourceType:       sourceType,
	}
	if err := store.RecordEvent(evt); err != nil {
		logs.Warn("failed to record event", logs.Err(err))
	}
}

// removeTagsForResource removes all tags associated with the given resource ARN.
func removeTagsForResource(store neptunestore.NeptuneStoreInterface, resourceArn string) {
	tags, err := store.GetTags(resourceArn)
	if err != nil || len(tags) == 0 {
		return
	}
	keys := make([]string, len(tags))
	for i, t := range tags {
		keys[i] = t.Key
	}
	if err := store.RemoveTags(resourceArn, keys); err != nil {
		logs.Warn("failed to remove tags on delete", logs.String("arn", resourceArn), logs.Err(err))
	}
}

// validateNeptuneTagTarget resolves the resource behind the RDS-style ARN
// the tag operations take as ResourceName, so a tag against a nonexistent
// resource fails with the resource-specific NotFoundFault the model
// attaches to these operations instead of silently persisting tags under
// an unowned key. ARN kinds the platform does not host cannot address an
// existing resource and are rejected as invalid input.
func validateNeptuneTagTarget(store neptunestore.NeptuneStoreInterface, resourceArn string) error {
	_, _, _, _, resource := arnutil.SplitARN(resourceArn)
	kind, name, ok := strings.Cut(resource, ":")
	if !ok || name == "" {
		return awserrors.NewInvalidParameterValueException("ResourceName must be a valid Neptune resource ARN")
	}
	var err error
	switch kind {
	case "cluster":
		_, err = store.GetCluster(name)
	case "db":
		_, err = store.GetInstance(name)
	case "cluster-snapshot":
		_, err = store.GetSnapshot(name)
	case "snapshot":
		_, err = store.GetInstanceSnapshot(name)
	case "cluster-endpoint":
		_, err = store.GetClusterEndpoint(name)
	case "global-cluster":
		_, err = store.GetGlobalCluster(name)
	case "cluster-pg":
		_, err = store.GetClusterParameterGroup(name)
	case "pg":
		_, err = store.GetParameterGroup(name)
	case "subgrp":
		_, err = store.GetSubnetGroup(name)
	case "es":
		_, err = store.GetEventSubscription(name)
	default:
		return awserrors.NewInvalidParameterValueException("ResourceName must be a valid Neptune resource ARN")
	}
	if err != nil {
		return translateStoreError(err)
	}
	return nil
}

// neptuneTagConfig builds the shared tag handler configuration binding the
// store's tag operations for the Neptune resource-tagging operations.
func (s *NeptuneService) neptuneTagConfig(store neptunestore.NeptuneStoreInterface) tags.TagHandlerConfig {
	return tags.TagHandlerConfig{
		Param: tags.TagOperationConfig{
			ResourceParam:   "ResourceName",
			TagsParam:       "Tags",
			TagKeysParam:    "TagKeys",
			TagKeyName:      "Key",
			TagValueName:    "Value",
			RequireTags:     false,
			RequireTagKeys:  false,
			RequireResource: true,
		},
		ParseTags: func(params map[string]interface{}) []tags.Tag {
			rawTags := getNeptuneTagList(params)
			result := make([]tags.Tag, 0, len(rawTags))
			for _, t := range rawTags {
				key, _ := t["Key"].(string)
				value, _ := t["Value"].(string)
				if key != "" {
					result = append(result, tags.Tag{Key: key, Value: value})
				}
			}
			return result
		},
		ValidateResource: func(ctx context.Context, resourceKey string) error {
			return validateNeptuneTagTarget(store, resourceKey)
		},
		ParseTagKeys: func(params map[string]interface{}) []string {
			return request.GetStringList(params, "TagKeys")
		},
		TagFunc: func(ctx context.Context, resourceKey string, tagList []tags.Tag) error {
			return store.AddTags(resourceKey, tagList)
		},
		UntagFunc: func(ctx context.Context, resourceKey string, tagKeys []string) error {
			return store.RemoveTags(resourceKey, tagKeys)
		},
		ListFunc: func(ctx context.Context, resourceKey string) ([]tags.Tag, error) {
			return store.GetTags(resourceKey)
		},
		FormatResponse: func(tagList []tags.Tag, _ string) (interface{}, error) {
			items := make([]interface{}, 0, len(tagList))
			for _, t := range tagList {
				items = append(items, map[string]interface{}{"Key": t.Key, "Value": t.Value})
			}
			return map[string]interface{}{
				"TagList": protocol.XMLElements{ElementName: "Tag", Items: items},
			}, nil
		},
		EmptyResponse: func() (interface{}, error) {
			return map[string]interface{}{}, nil
		},
		MapError: func(err error) error {
			return err
		},
	}
}

// createDBClusterEndpointCore validates cluster existence and persists a new
// custom cluster endpoint, returning the stored record.
func (s *NeptuneService) createDBClusterEndpointCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *CreateDBClusterEndpointInput) (*neptunestore.DBClusterEndpoint, error) {
	epID := in.DBClusterEndpointIdentifier
	if epID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterEndpointIdentifier is required")
	}
	clusterID := in.DBClusterIdentifier
	if clusterID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterIdentifier is required")
	}
	epType := in.EndpointType
	if epType == "" {
		return nil, awserrors.NewMissingParameter("EndpointType is required")
	}

	if _, err := store.GetCluster(clusterID); err != nil {
		return nil, awserrors.NewAWSError("DBClusterNotFoundFault", fmt.Sprintf("DBCluster %s not found", clusterID), http.StatusNotFound)
	}

	ep := &neptunestore.DBClusterEndpoint{
		DBClusterEndpointIdentifier: epID,
		DBClusterIdentifier:         clusterID,
		Endpoint:                    fmt.Sprintf("%s.cluster-%s.%s.amazonaws.com", epID, clusterID, in.Region),
		Status:                      "available",
		EndpointType:                epType,
		ExcludedMembers:             in.ExcludedMembers,
		StaticMembers:               in.StaticMembers,
		DBClusterEndpointArn:        arnutil.NewARNBuilder(in.AccountID, in.Region).RDS().ClusterEndpoint(epID),
	}

	if err := store.CreateClusterEndpoint(ep); err != nil {
		return nil, err
	}

	return ep, nil
}

// describeDBClusterEndpointsCore resolves one endpoint by identifier or all
// endpoints of a cluster.
func (s *NeptuneService) describeDBClusterEndpointsCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DescribeDBClusterEndpointsInput) (*DescribeDBClusterEndpointsResult, error) {
	if in.DBClusterEndpointIdentifier != "" {
		ep, err := store.GetClusterEndpoint(in.DBClusterEndpointIdentifier)
		if err != nil {
			return nil, translateStoreError(err)
		}
		return &DescribeDBClusterEndpointsResult{Endpoint: ep}, nil
	}

	endpoints, err := store.ListClusterEndpoints(in.DBClusterIdentifier)
	if err != nil {
		return nil, err
	}
	return &DescribeDBClusterEndpointsResult{Endpoints: endpoints}, nil
}

// modifyDBClusterEndpointCore applies member-list and type updates to a
// custom cluster endpoint, returning the updated record.
func (s *NeptuneService) modifyDBClusterEndpointCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ModifyDBClusterEndpointInput) (*neptunestore.DBClusterEndpoint, error) {
	epID := in.DBClusterEndpointIdentifier
	if epID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterEndpointIdentifier is required")
	}

	ep, err := store.GetClusterEndpoint(epID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if len(in.ExcludedMembers) > 0 {
		ep.ExcludedMembers = in.ExcludedMembers
	}
	if len(in.StaticMembers) > 0 {
		ep.StaticMembers = in.StaticMembers
	}
	if newType := in.EndpointType; newType != "" {
		ep.EndpointType = newType
	}

	ep.Status = "available"
	if err := store.UpdateClusterEndpoint(ep); err != nil {
		return nil, err
	}

	return ep, nil
}

// deleteDBClusterEndpointCore deletes a custom cluster endpoint, returning
// the removed record.
func (s *NeptuneService) deleteDBClusterEndpointCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DeleteDBClusterEndpointInput) (*neptunestore.DBClusterEndpoint, error) {
	epID := in.DBClusterEndpointIdentifier
	if epID == "" {
		return nil, awserrors.NewMissingParameter("DBClusterEndpointIdentifier is required")
	}

	ep, err := store.GetClusterEndpoint(epID)
	if err != nil {
		return nil, translateStoreError(err)
	}

	if err := store.DeleteClusterEndpoint(epID); err != nil {
		return nil, err
	}

	return ep, nil
}

// describeEventsCore validates the time filters, applies the duration window
// and loads the matching events from the store. MaxRecords follows the
// documented page-size window (default 100, bounds 20-100) through the
// shared resolver.
func (s *NeptuneService) describeEventsCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DescribeEventsInput) (*neptunestore.EventListResult, error) {
	// The documented MaxRecords window (default 100, bounds 20-100) is
	// resolved by the shared helper both engine planes use.
	maxRecords := rdssvc.ResolveDescribeMaxRecords(in.MaxRecords)

	var startTime time.Time
	if stStr := in.StartTime; stStr != "" {
		var err error
		startTime, err = time.Parse(time.RFC3339, stStr)
		if err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "Invalid StartTime format: use RFC3339", http.StatusBadRequest)
		}
	}

	var endTime time.Time
	if in.Duration > 0 && !startTime.IsZero() {
		endTime = startTime.Add(time.Duration(in.Duration) * time.Minute)
	} else if etStr := in.EndTime; etStr != "" {
		var err error
		endTime, err = time.Parse(time.RFC3339, etStr)
		if err != nil {
			return nil, awserrors.NewAWSError("InvalidParameterValue", "Invalid EndTime format: use RFC3339", http.StatusBadRequest)
		}
	}

	opts := neptunestore.EventListOptions{
		SourceType:       in.SourceType,
		SourceIdentifier: in.SourceIdentifier,
		StartTime:        startTime,
		EndTime:          endTime,
		Marker:           in.Marker,
		MaxRecords:       maxRecords,
	}

	return store.ListEvents(opts)
}

// applyPendingMaintenanceActionCore validates the action/opt-in values and
// the referenced resource's existence.
func (s *NeptuneService) applyPendingMaintenanceActionCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *ApplyPendingMaintenanceActionInput) error {
	resourceID := in.ResourceIdentifier
	if resourceID == "" {
		return awserrors.NewMissingParameter("ResourceIdentifier is required")
	}
	action := in.ApplyAction
	if action == "" {
		return awserrors.NewMissingParameter("ApplyAction is required")
	}
	if err := validateApplyAction(action); err != nil {
		return awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}
	optIn := in.OptInType
	if optIn == "" {
		return awserrors.NewMissingParameter("OptInType is required")
	}
	if err := validateOptInType(optIn); err != nil {
		return awserrors.NewAWSError("InvalidParameterValue", err.Error(), http.StatusBadRequest)
	}

	// Validate that the resource exists.
	if _, err := store.GetCluster(resourceID); err != nil {
		if _, err2 := store.GetInstance(resourceID); err2 != nil {
			return awserrors.NewAWSError("ResourceNotFoundFault",
				fmt.Sprintf("Resource %s not found", resourceID), http.StatusNotFound)
		}
	}
	return nil
}

// describeValidDBInstanceModificationsCore validates the instance identifier
// and returns the (currently empty) valid modification set.
func (s *NeptuneService) describeValidDBInstanceModificationsCore(ctx context.Context, store neptunestore.NeptuneStoreInterface, in *DescribeValidDBInstanceModificationsInput) (interface{}, error) {
	if in.DBInstanceIdentifier == "" {
		return nil, awserrors.NewMissingParameter("DBInstanceIdentifier is required")
	}

	return map[string]interface{}{
		"ValidDBInstanceModificationsMessage": map[string]interface{}{
			"Storage": protocol.XMLElements{ElementName: "ValidStorageOptions", Items: []interface{}{}},
		},
	}, nil
}
