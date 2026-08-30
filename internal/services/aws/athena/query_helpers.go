package athena

import (
	"fmt"
	"strings"
	"time"

	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

func (s *AthenaService) detectStatementType(query string) athenastore.StatementType {
	upperQuery := strings.ToUpper(strings.TrimSpace(query))
	if strings.HasPrefix(upperQuery, "SELECT") || strings.HasPrefix(upperQuery, "WITH") {
		return athenastore.StatementTypeDML
	}
	if strings.HasPrefix(upperQuery, "CREATE") || strings.HasPrefix(upperQuery, "DROP") ||
		strings.HasPrefix(upperQuery, "ALTER") || strings.HasPrefix(upperQuery, "SHOW") ||
		strings.HasPrefix(upperQuery, "DESCRIBE") || strings.HasPrefix(upperQuery, "MSCK") {
		return athenastore.StatementTypeDDL
	}
	if strings.HasPrefix(upperQuery, "INSERT") || strings.HasPrefix(upperQuery, "UPDATE") ||
		strings.HasPrefix(upperQuery, "DELETE") {
		return athenastore.StatementTypeDML
	}
	return athenastore.StatementTypeUtility
}

func (s *AthenaService) queryExecutionToResponse(qe *athenastore.QueryExecution) map[string]interface{} {
	response := map[string]interface{}{
		"QueryExecutionId": qe.QueryExecutionId,
		"Query":            qe.Query,
		"StatementType":    qe.StatementType,
	}

	if qe.WorkGroup != "" {
		response["WorkGroup"] = qe.WorkGroup
	}

	if qe.QueryExecutionContext != nil {
		response["QueryExecutionContext"] = map[string]interface{}{
			"Database": qe.QueryExecutionContext.Database,
			"Catalog":  qe.QueryExecutionContext.Catalog,
		}
	}

	if qe.ResultConfiguration != nil {
		resultConfig := map[string]interface{}{}
		if qe.ResultConfiguration.OutputLocation != "" {
			resultConfig["OutputLocation"] = qe.ResultConfiguration.OutputLocation
		}
		if qe.ResultConfiguration.EncryptionConfiguration != nil {
			resultConfig["EncryptionConfiguration"] = map[string]interface{}{
				"EncryptionOption": qe.ResultConfiguration.EncryptionConfiguration.EncryptionOption,
				"KmsKey":           qe.ResultConfiguration.EncryptionConfiguration.KmsKey,
			}
		}
		if qe.ResultConfiguration.ExpectedBucketOwner != "" {
			resultConfig["ExpectedBucketOwner"] = qe.ResultConfiguration.ExpectedBucketOwner
		}
		if qe.ResultConfiguration.ACLConfiguration != nil {
			resultConfig["AclConfiguration"] = map[string]interface{}{
				"S3AclOption": qe.ResultConfiguration.ACLConfiguration.S3ACLOption,
			}
		}
		response["ResultConfiguration"] = resultConfig
	}

	if qe.Status != nil {
		status := map[string]interface{}{
			"State":              qe.Status.State,
			"SubmissionDateTime": float64(qe.Status.SubmissionDateTime.UnixNano()) / 1e9,
		}
		if qe.Status.StateChangeReason != "" {
			status["StateChangeReason"] = qe.Status.StateChangeReason
		}
		if !qe.Status.CompletionDateTime.IsZero() {
			status["CompletionDateTime"] = float64(qe.Status.CompletionDateTime.UnixNano()) / 1e9
		}
		if qe.Status.AthenaError != nil {
			// AWS SDK expects ErrorType as int32 (not string). Our store/proto
			// uses string, so we omit it if non-numeric to avoid deserialisation
			// errors. ErrorCategory is already int32-compatible.
			athenaErrMap := map[string]interface{}{
				"ErrorCategory": qe.Status.AthenaError.ErrorCategory,
				"Retryable":     qe.Status.AthenaError.Retryable,
				"ErrorMessage":  qe.Status.AthenaError.ErrorMessage,
			}
			if qe.Status.AthenaError.SyntaxErrorRow != 0 {
				athenaErrMap["SyntaxErrorRow"] = qe.Status.AthenaError.SyntaxErrorRow
			}
			if qe.Status.AthenaError.SyntaxErrorColumn != 0 {
				athenaErrMap["SyntaxErrorColumn"] = qe.Status.AthenaError.SyntaxErrorColumn
			}
			status["AthenaError"] = athenaErrMap
		}
		response["Status"] = status
	}

	if qe.Statistics != nil {
		statistics := map[string]interface{}{
			"EngineExecutionTimeInMillis":   qe.Statistics.EngineExecutionTimeInMillis,
			"DataScannedInBytes":            qe.Statistics.DataScannedInBytes,
			"TotalExecutionTimeInMillis":    qe.Statistics.TotalExecutionTimeInMillis,
			"QueryQueueTimeInMillis":        qe.Statistics.QueryQueueTimeInMillis,
			"QueryPlanningTimeInMillis":     qe.Statistics.QueryPlanningTimeInMillis,
			"ServiceProcessingTimeInMillis": qe.Statistics.ServiceProcessingTimeInMillis,
		}
		if qe.Statistics.ResultReuseInformation != nil {
			statistics["ResultReuseInformation"] = map[string]interface{}{
				"ReusedPreviousResult": qe.Statistics.ResultReuseInformation.ReusedPreviousResult,
			}
		}
		response["Statistics"] = statistics
	}

	return response
}

func (s *AthenaService) parseResultConfiguration(resultConfigMap map[string]interface{}) (*athenastore.ResultConfiguration, error) {
	resultConfiguration := &athenastore.ResultConfiguration{}

	if outputLocation, ok := resultConfigMap["OutputLocation"].(string); ok {
		resultConfiguration.OutputLocation = outputLocation
	}

	if encConfigMap, ok := resultConfigMap["EncryptionConfiguration"].(map[string]interface{}); ok {
		resultConfiguration.EncryptionConfiguration = &athenastore.EncryptionConfiguration{}
		if encOption, ok := encConfigMap["EncryptionOption"].(string); ok {
			resultConfiguration.EncryptionConfiguration.EncryptionOption = encOption
		}
		if kmsKey, ok := encConfigMap["KmsKey"].(string); ok {
			resultConfiguration.EncryptionConfiguration.KmsKey = kmsKey
		}
	}

	if expectedBucketOwner, ok := resultConfigMap["ExpectedBucketOwner"].(string); ok {
		resultConfiguration.ExpectedBucketOwner = expectedBucketOwner
	}

	if aclConfigMap, ok := resultConfigMap["AclConfiguration"].(map[string]interface{}); ok {
		aclOption, _ := aclConfigMap["S3AclOption"].(string)
		if aclOption != "BUCKET_OWNER_FULL_CONTROL" {
			return nil, invalidRequestParameter("AclConfiguration.S3AclOption must be BUCKET_OWNER_FULL_CONTROL")
		}
		resultConfiguration.ACLConfiguration = &athenastore.ACLConfiguration{S3ACLOption: aclOption}
	}

	return resultConfiguration, nil
}

func (s *AthenaService) hasJoin(selectStmt *sqlparser.Select) bool {
	for _, tableExpr := range selectStmt.From {
		if _, ok := tableExpr.(*sqlparser.JoinTableExpr); ok {
			return true
		}
	}
	return false
}

func (s *AthenaService) buildResultSetFromStoredTable(tableData *athenastore.StoredTable, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	var columnInfo []athenastore.ColumnInfo
	for _, col := range tableData.Columns {
		columnInfo = append(columnInfo, athenastore.ColumnInfo{
			Label: col.Name,
			Name:  col.Name,
			Type:  col.Type,
		})
	}

	var resultRows []athenastore.Row
	dataScanned := int64(0)
	for _, row := range tableData.Rows {
		var data []athenastore.Datum
		for _, col := range tableData.Columns {
			val := ""
			if v, ok := row.Values[col.Name]; ok {
				val = fmt.Sprintf("%v", v)
			}
			dataScanned += int64(len(val))
			data = append(data, athenastore.Datum{VarCharValue: val})
		}
		resultRows = append(resultRows, athenastore.Row{Data: data})
	}

	headerRow := athenastore.Row{}
	for _, col := range columnInfo {
		headerRow.Data = append(headerRow.Data, athenastore.Datum{VarCharValue: col.Name})
	}
	resultRows = append([]athenastore.Row{headerRow}, resultRows...)

	stats := &athenastore.QueryExecutionStatistics{
		QueryPlanningTimeInMillis: time.Since(startTime).Milliseconds(),
		DataScannedInBytes:        dataScanned,
	}

	return &athenastore.ResultSet{
		Rows:              resultRows,
		ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: columnInfo},
	}, stats, nil
}
