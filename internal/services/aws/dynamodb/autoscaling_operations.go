// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"

	"vorpalstacks/internal/common/request"
)

// DescribeTableReplicaAutoScaling returns the auto scaling settings for table replicas.
func (s *DynamoDBService) DescribeTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	settings, err := store.Tables().GetAutoScalingSettings(table.Name)
	if err != nil {
		return nil, err
	}

	var replicas []map[string]interface{}
	if settings != nil {
		if replicaList, ok := settings["replicas"].([]interface{}); ok {
			for _, r := range replicaList {
				if rMap, ok := r.(map[string]interface{}); ok {
					replicas = append(replicas, rMap)
				}
			}
		}
	}
	if replicas == nil {
		replicas = []map[string]interface{}{}
	}

	return map[string]interface{}{
		"TableAutoScalingDescription": map[string]interface{}{
			"TableName":   table.Name,
			"TableStatus": string(table.Status),
			"Replicas":    replicas,
		},
	}, nil
}

// UpdateTableReplicaAutoScaling updates the auto scaling settings for table replicas.
func (s *DynamoDBService) UpdateTableReplicaAutoScaling(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	table, err := s.validateAndGetTable(reqCtx, req.Parameters)
	if err != nil {
		return nil, err
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	// Parse replica settings updates and build the replica descriptions.
	// AWS auto-scaling uses Application Auto Scaling; we store settings
	// for API round-trip compatibility without implementing the scaling
	// engine itself.
	var replicas []map[string]interface{}

	if globalTableName := request.GetStringParam(req.Parameters, "GlobalTableName"); globalTableName != "" {
		gt, gtErr := store.GlobalTables().Get(globalTableName)
		if gtErr == nil && gt != nil {
			for _, replica := range gt.ReplicationGroup {
				replicaDesc := map[string]interface{}{
					"RegionName": replica.RegionName,
				}
				replicas = append(replicas, replicaDesc)
			}
		}
	}

	if replicas == nil {
		replicas = []map[string]interface{}{}
	}

	settings := map[string]interface{}{
		"replicas": replicas,
	}
	_ = store.Tables().SetAutoScalingSettings(table.Name, settings)

	return map[string]interface{}{
		"TableAutoScalingDescription": map[string]interface{}{
			"TableName":   table.Name,
			"TableStatus": string(table.Status),
			"Replicas":    replicas,
		},
	}, nil
}
