// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"time"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

// DescribeContinuousBackups returns the continuous backup settings for a table.
func (s *DynamoDBService) DescribeContinuousBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy declares TableNotFoundException for this op.
	table, err := s.validateAndGetTableWithErr(reqCtx, req.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	pitr, err := store.Tables().GetPointInTimeRecovery(tableName)
	if err != nil {
		return nil, err
	}

	pitrStatus := "DISABLED"
	var earliest, latest int64
	recoveryPeriod := pitrDefaultRecoveryPeriodDays
	if pitr != nil && pitr.Status == dbstore.PITRStatusEnabled {
		pitrStatus = "ENABLED"
		earliest = pitr.EarliestRestorableDateTime.Unix()
		latest = pitr.LatestRestorableDateTime.Unix()
		if pitr.RecoveryPeriodInDays > 0 {
			recoveryPeriod = pitr.RecoveryPeriodInDays
		}
	}

	return map[string]interface{}{
		"ContinuousBackupsDescription": map[string]interface{}{
			"ContinuousBackupsStatus": pitrStatus,
			"PointInTimeRecoveryDescription": map[string]interface{}{
				"PointInTimeRecoveryStatus":  pitrStatus,
				"EarliestRestorableDateTime": earliest,
				"LatestRestorableDateTime":   latest,
				"RecoveryPeriodInDays":       recoveryPeriod,
			},
		},
	}, nil
}

// UpdateContinuousBackups enables or disables continuous backup for a table.
func (s *DynamoDBService) UpdateContinuousBackups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	// Smithy declares TableNotFoundException for this op.
	table, err := s.validateAndGetTableWithErr(reqCtx, req.Parameters, ErrTableNotFoundException)
	if err != nil {
		return nil, err
	}
	tableName := table.Name

	pitrSpec, ok := req.Parameters["PointInTimeRecoverySpecification"].(map[string]interface{})
	if !ok {
		return nil, ErrInvalidParameter
	}

	enabled, err := validateBoolParam(pitrSpec, "PointInTimeRecoveryEnabled", false)
	if err != nil {
		return nil, err
	}

	recoveryPeriod := pitrDefaultRecoveryPeriodDays
	if _, ok := pitrSpec["RecoveryPeriodInDays"]; ok {
		rp := request.GetIntParam(pitrSpec, "RecoveryPeriodInDays")
		if !validateRecoveryPeriodInDays(rp) {
			return nil, ErrInvalidParameter
		}
		recoveryPeriod = rp
	}

	pitr := &dbstore.PointInTimeRecoveryDescription{
		Status: dbstore.PITRStatusDisabled,
	}
	if enabled {
		pitr.Status = dbstore.PITRStatusEnabled
		pitr.EarliestRestorableDateTime = time.Now().Add(-time.Duration(recoveryPeriod) * 24 * time.Hour)
		pitr.LatestRestorableDateTime = time.Now()
		pitr.RecoveryPeriodInDays = recoveryPeriod
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	if err := store.Tables().SetPointInTimeRecovery(tableName, pitr); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"ContinuousBackupsDescription": map[string]interface{}{
			"ContinuousBackupsStatus": string(pitr.Status),
			"PointInTimeRecoveryDescription": map[string]interface{}{
				"PointInTimeRecoveryStatus": string(pitr.Status),
			},
		},
	}, nil
}
