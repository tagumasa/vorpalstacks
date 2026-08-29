package dynamodb

import (
	"context"
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

type statementType int

const (
	statementTypeUnknown statementType = iota
	statementTypeRead
	statementTypeWrite
)

// executeTransactionInput carries the already-typed TransactStatements member
// plus the raw wire parameters (consumed for the ClientRequestToken
// idempotency member and the ReturnConsumedCapacity reporting).
type executeTransactionInput struct {
	TransactStatements []interface{}
	Parameters         map[string]interface{}
}

// executeTransactionCore is the single validation and persistence path of the
// PartiQL ExecuteTransaction operation: statement classification with the
// read/write exclusivity rule, then either one snapshot view or one
// transactional update across every statement.
func (s *DynamoDBService) executeTransactionCore(ctx context.Context, reqCtx *request.RequestContext, in executeTransactionInput) (map[string]interface{}, error) {
	statements := in.TransactStatements

	if len(statements) == 0 {
		return nil, ErrInvalidParameter
	}

	if len(statements) > transactMaxItems {
		return nil, ErrInvalidParameter
	}

	clientRequestToken := request.GetStringParam(in.Parameters, "ClientRequestToken")
	if !validateClientRequestToken(clientRequestToken) {
		return nil, ErrInvalidParameter
	}

	parsedStatements := make([]struct {
		statement string
		params    *partiQLParams
		stmtType  statementType
	}, 0, len(statements))

	var hasRead, hasWrite bool

	for _, stmt := range statements {
		stmtMap, ok := stmt.(map[string]interface{})
		if !ok {
			return nil, ErrInvalidParameter
		}

		statement, _ := stmtMap["Statement"].(string)
		if !validatePartiQLStatement(statement) {
			return nil, ErrInvalidParameter
		}

		params := parsePartiQLParams(stmtMap)

		upperStmt := strings.ToUpper(strings.TrimSpace(statement))
		var stmtType statementType
		switch {
		case strings.HasPrefix(upperStmt, "SELECT"):
			stmtType = statementTypeRead
			hasRead = true
		case strings.HasPrefix(upperStmt, "INSERT"),
			strings.HasPrefix(upperStmt, "UPDATE"),
			strings.HasPrefix(upperStmt, "DELETE"):
			stmtType = statementTypeWrite
			hasWrite = true
		default:
			return nil, ErrInvalidParameter
		}

		parsedStatements = append(parsedStatements, struct {
			statement string
			params    *partiQLParams
			stmtType  statementType
		}{statement, params, stmtType})
	}

	if hasRead && hasWrite {
		return nil, ErrTransactionConflict
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	responses := make([]map[string]interface{}, len(parsedStatements))

	if hasRead {
		err = store.View(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for i, ps := range parsedStatements {
				result, err := s.executePartiQLSelectInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				if err != nil {
					return err
				}

				var item interface{} = nil
				if resultMap, ok := result.(map[string]interface{}); ok {
					if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
						item = items[0]
					}
				}
				responses[i] = map[string]interface{}{"Item": item}
			}
			return nil
		})
	} else {
		err = store.Update(ctx, func(txn *dbstore.DynamoDBTxn) error {
			for i, ps := range parsedStatements {
				upperStmt := strings.ToUpper(strings.TrimSpace(ps.statement))
				var result interface{}
				var execErr error

				switch {
				case strings.HasPrefix(upperStmt, "INSERT"):
					result, execErr = s.executePartiQLInsertInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				case strings.HasPrefix(upperStmt, "UPDATE"):
					result, execErr = s.executePartiQLUpdateInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				case strings.HasPrefix(upperStmt, "DELETE"):
					result, execErr = s.executePartiQLDeleteInTxn(ctx, reqCtx, txn, ps.statement, ps.params)
				}

				if execErr != nil {
					return execErr
				}

				var item interface{} = nil
				if resultMap, ok := result.(map[string]interface{}); ok {
					if items, ok := resultMap["Items"].([]map[string]interface{}); ok && len(items) > 0 {
						item = items[0]
					}
				}
				responses[i] = map[string]interface{}{"Item": item}
			}
			return nil
		})
	}

	if err != nil {
		return nil, err
	}

	resp := map[string]interface{}{
		"Responses": responses,
	}

	returnConsumedCapacity := getReturnConsumedCapacity(in.Parameters)
	if returnConsumedCapacity == "TOTAL" || returnConsumedCapacity == "INDEXES" {
		tableNames := make(map[string]bool)
		for _, ps := range parsedStatements {
			if tn := extractTableNameFromStatement(ps.statement); tn != "" {
				tableNames[tn] = true
			}
		}
		var consumedCapacities []map[string]interface{}
		for tableName := range tableNames {
			consumedCapacities = append(consumedCapacities, buildConsumedCapacityResponse(tableName, 2.0))
		}
		if len(consumedCapacities) > 0 {
			resp["ConsumedCapacity"] = consumedCapacities
		}
	}

	return resp, nil
}
