package athena

import (
	"fmt"
	"strings"
	"time"

	athenastore "vorpalstacks/internal/store/aws/athena"
	"vorpalstacks/pkg/sqlparser"
)

func (s *Service) convertCastOperator(sql string) string {
	var result strings.Builder
	inString := false
	i := 0

	for i < len(sql) {
		ch := sql[i]

		if ch == '\'' {
			inString = !inString
			result.WriteByte(ch)
			i++
			continue
		}

		if inString {
			result.WriteByte(ch)
			i++
			continue
		}

		if i+1 < len(sql) && sql[i:i+2] == "::" {
			j := i + 2
			for j < len(sql) && (isAlphaNum(sql[j]) || sql[j] == '_') {
				j++
			}
			i = j
			continue
		}

		result.WriteByte(ch)
		i++
	}

	return result.String()
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}

func (s *Service) detectStatementType(query string) athenastore.StatementType {
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

func (s *Service) queryExecutionToResponse(qe *athenastore.QueryExecution) map[string]interface{} {
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
			status["AthenaError"] = map[string]interface{}{
				"ErrorCategory":     qe.Status.AthenaError.ErrorCategory,
				"ErrorType":         qe.Status.AthenaError.ErrorType,
				"Retryable":         qe.Status.AthenaError.Retryable,
				"ErrorMessage":      qe.Status.AthenaError.ErrorMessage,
				"SyntaxErrorRow":    qe.Status.AthenaError.SyntaxErrorRow,
				"SyntaxErrorColumn": qe.Status.AthenaError.SyntaxErrorColumn,
			}
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

func (s *Service) parseResultConfiguration(resultConfigMap map[string]interface{}) *athenastore.ResultConfiguration {
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

	return resultConfiguration
}

func (s *Service) hasJoin(selectStmt *sqlparser.Select) bool {
	for _, tableExpr := range selectStmt.From {
		if _, ok := tableExpr.(*sqlparser.JoinTableExpr); ok {
			return true
		}
	}
	return false
}

func (s *Service) buildResultSetFromStoredTable(tableData *athenastore.StoredTable, startTime time.Time) (*athenastore.ResultSet, *athenastore.QueryExecutionStatistics, error) {
	var columnInfo []athenastore.ColumnInfo
	for _, col := range tableData.Columns {
		columnInfo = append(columnInfo, athenastore.ColumnInfo{
			Label: col.Name,
			Name:  col.Name,
			Type:  col.Type,
		})
	}

	var resultRows []athenastore.Row
	for _, row := range tableData.Rows {
		var data []athenastore.Datum
		for _, col := range tableData.Columns {
			val := ""
			if v, ok := row.Values[col.Name]; ok {
				val = fmt.Sprintf("%v", v)
			}
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
		DataScannedInBytes:        int64(len(tableData.Rows) * 100),
	}

	return &athenastore.ResultSet{
		Rows:              resultRows,
		ResultSetMetadata: &athenastore.ResultSetMetadata{ColumnInfo: columnInfo},
	}, stats, nil
}
