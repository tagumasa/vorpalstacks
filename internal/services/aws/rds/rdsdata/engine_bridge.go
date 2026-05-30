package rdsdata

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"
)

func newSQLContext(database string) *sql.Context {
	ctx := sql.NewContext(context.Background())
	if database != "" {
		ctx.SetCurrentDatabase(database)
	}
	return ctx
}

// executeSQL runs a SQL string on the given engine and returns a formatted response.
func executeSQL(engine *sqle.Engine, sqlCtx *sql.Context, sqlStr string, includeMetadata bool, formatRecordsAs string) (*ExecuteStatementResponse, error) {
	schema, rowIter, _, err := engine.Query(sqlCtx, sqlStr)
	if err != nil {
		return nil, err
	}

	var rows []sql.Row
	if rowIter != nil {
		rows, err = sql.RowIterToRows(sqlCtx, rowIter)
		if err != nil {
			return nil, fmt.Errorf("failed to read rows: %w", err)
		}
	}

	upperSQL := strings.ToUpper(strings.TrimSpace(sqlStr))
	isSelect := strings.HasPrefix(upperSQL, "SELECT") || strings.HasPrefix(upperSQL, "SHOW") || strings.HasPrefix(upperSQL, "DESCRIBE") || strings.HasPrefix(upperSQL, "EXPLAIN")

	resp := &ExecuteStatementResponse{}

	if isSelect && len(rows) > 0 {
		resp.Records = convertRows(rows, schema)
	} else if isSelect && len(rows) == 0 {
		resp.Records = [][]Field{}
	}

	if isSelect && includeMetadata && schema != nil {
		resp.ColumnMetadata = convertSchema(schema)
	}

	if !isSelect {
		resp.NumberOfRecordsUpdated = int64(len(rows))
	}

	if formatRecordsAs == "JSON" && isSelect {
		formatted, err := formatAsJSON(rows, schema)
		if err == nil {
			resp.FormattedRecords = formatted
		}
	}

	return resp, nil
}

// convertRows converts sql.Rows to RDS Data API Field format.
func convertRows(rows []sql.Row, schema sql.Schema) [][]Field {
	result := make([][]Field, len(rows))
	for i, row := range rows {
		fields := make([]Field, len(row))
		for j, val := range row {
			fields[j] = convertValue(val)
		}
		result[i] = fields
	}
	return result
}

// convertSchema converts sql.Schema to ColumnMetadata.
func convertSchema(schema sql.Schema) []ColumnMetadata {
	if schema == nil {
		return nil
	}
	metadata := make([]ColumnMetadata, len(schema))
	for i, col := range schema {
		metadata[i] = ColumnMetadata{
			Name:            col.Name,
			Label:           col.Name,
			IsSigned:        isSignedType(col.Type),
			IsCaseSensitive: isCaseSensitiveType(col.Type),
			Nullable:        nullableToInt(col.Nullable),
			Type:            sqlTypeToRDS(col.Type),
			TypeName:        col.Type.String(),
		}
	}
	return metadata
}

// convertValue converts a Go value to an RDS Data API Field.
func convertValue(val interface{}) Field {
	if val == nil {
		t := true
		return Field{IsNull: &t}
	}

	switch v := val.(type) {
	case int:
		lv := int64(v)
		return Field{LongValue: &lv}
	case int8:
		lv := int64(v)
		return Field{LongValue: &lv}
	case int16:
		lv := int64(v)
		return Field{LongValue: &lv}
	case int32:
		lv := int64(v)
		return Field{LongValue: &lv}
	case int64:
		return Field{LongValue: &v}
	case uint:
		lv := int64(v)
		return Field{LongValue: &lv}
	case uint8:
		lv := int64(v)
		return Field{LongValue: &lv}
	case uint16:
		lv := int64(v)
		return Field{LongValue: &lv}
	case uint32:
		lv := int64(v)
		return Field{LongValue: &lv}
	case float32:
		dv := float64(v)
		return Field{DoubleValue: &dv}
	case float64:
		return Field{DoubleValue: &v}
	case bool:
		return Field{BooleanValue: &v}
	case string:
		return Field{StringValue: &v}
	case []byte:
		return Field{BlobValue: v}
	case time.Time:
		sv := v.Format("2006-01-02 15:04:05.999")
		return Field{StringValue: &sv}
	case fmt.Stringer:
		sv := v.String()
		return Field{StringValue: &sv}
	default:
		sv := fmt.Sprintf("%v", v)
		return Field{StringValue: &sv}
	}
}

// fieldToValue converts a Field to a Value (deprecated ExecuteSql format).
func fieldToValue(f Field) Value {
	v := Value{}
	if f.IsNull != nil && *f.IsNull {
		t := true
		v.IsNull = &t
		return v
	}
	if f.StringValue != nil {
		v.StringValue = f.StringValue
	}
	if f.LongValue != nil {
		v.LongValue = f.LongValue
	}
	if f.DoubleValue != nil {
		v.DoubleValue = f.DoubleValue
	}
	if f.BooleanValue != nil {
		v.BitValue = f.BooleanValue
	}
	if f.BlobValue != nil {
		v.BlobValue = f.BlobValue
	}
	return v
}

// formatAsJSON converts rows to a JSON string for formatRecordsAs=JSON.
func formatAsJSON(rows []sql.Row, schema sql.Schema) (string, error) {
	if schema == nil || len(rows) == 0 {
		return "[]", nil
	}

	records := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		m := make(map[string]interface{}, len(schema))
		for j, col := range schema {
			if j < len(row) {
				m[col.Name] = valueToJSON(row[j])
			}
		}
		records[i] = m
	}

	data, err := json.Marshal(records)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func valueToJSON(val interface{}) interface{} {
	if val == nil {
		return nil
	}
	switch v := val.(type) {
	case time.Time:
		return v.Format("2006-01-02 15:04:05.999")
	case []byte:
		return string(v)
	default:
		return v
	}
}

// sqlTypeToRDS maps go-mysql-server types to JDBC type codes used by RDS Data API.
func sqlTypeToRDS(t sql.Type) int32 {
	if t == nil {
		return 12
	}
	typeStr := t.String()
	switch {
	case strings.Contains(typeStr, "BIGINT"):
		return -5
	case strings.Contains(typeStr, "INT"), strings.Contains(typeStr, "TINYINT"),
		strings.Contains(typeStr, "SMALLINT"), strings.Contains(typeStr, "MEDIUMINT"):
		return 4
	case strings.Contains(typeStr, "FLOAT"):
		return 7
	case strings.Contains(typeStr, "DOUBLE"):
		return 8
	case strings.Contains(typeStr, "DECIMAL"), strings.Contains(typeStr, "NUMERIC"):
		return 3
	case strings.Contains(typeStr, "BIT"), strings.Contains(typeStr, "BOOLEAN"), strings.Contains(typeStr, "BOOL"):
		return 16
	case strings.Contains(typeStr, "BLOB"), strings.Contains(typeStr, "BINARY"):
		return -2
	case strings.Contains(typeStr, "DATE"):
		return 91
	case strings.Contains(typeStr, "DATETIME"), strings.Contains(typeStr, "TIMESTAMP"):
		return 93
	default:
		return 12
	}
}

func isSignedType(t sql.Type) bool {
	if t == nil {
		return true
	}
	return !strings.Contains(t.String(), "UNSIGNED")
}

func isCaseSensitiveType(t sql.Type) bool {
	if t == nil {
		return true
	}
	s := t.String()
	return !strings.Contains(s, "TEXT") && !strings.Contains(s, "BLOB")
}

func nullableToInt(nullable bool) int32 {
	if nullable {
		return 1
	}
	return 0
}

// parseArn extracts the resource identifier and type from an RDS ARN.
func parseArn(arn string) (string, string) {
	parts := strings.Split(arn, ":")
	if len(parts) < 7 {
		return "", ""
	}
	resourceType := parts[5]
	identifier := parts[6]

	switch resourceType {
	case "db":
		return identifier, "db"
	case "cluster":
		return identifier, "cluster"
	default:
		return "", ""
	}
}

// splitSQL splits a multi-statement SQL string by semicolons.
func splitSQL(sqlStr string) []string {
	var statements []string
	for _, stmt := range strings.Split(sqlStr, ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}
	return statements
}

// substituteParameters replaces named parameters (:name) in SQL with values.
// It uses word-boundary matching to avoid partial replacement (e.g. :id
// matching inside :id2) and skips occurrences inside string literals.
func substituteParameters(sqlStr string, params []SqlParameter) string {
	result := sqlStr
	for _, p := range params {
		if p.Name == "" || p.Value == nil {
			continue
		}
		replacement := fieldToSQLString(p.Value)
		// \b cannot anchor before : because : is a non-word character; use (^|\W) instead.
		pattern := regexp.MustCompile(`(^|\W)` + regexp.QuoteMeta(":"+p.Name) + `($|\W)`)
		result = replaceOutsideStringsWithCapture(result, pattern, replacement)
	}
	return result
}

// replaceOutsideStringsWithCapture is like replaceOutsideStrings but preserves
// the captured boundary characters around the match.
func replaceOutsideStringsWithCapture(s string, re *regexp.Regexp, replacement string) string {
	var b strings.Builder
	inString := false
	i := 0
	for i < len(s) {
		if s[i] == '\'' {
			if inString && i+1 < len(s) && s[i+1] == '\'' {
				b.WriteString("''")
				i += 2
				continue
			}
			b.WriteByte(s[i])
			inString = !inString
			i++
			continue
		}
		if !inString {
			loc := re.FindStringSubmatchIndex(s[i:])
			if loc != nil && loc[0] == 0 {
				// loc[2:4] = group 1 (leading boundary), loc[4:6] = group 2 (trailing boundary)
				if loc[2] >= 0 && loc[3] > loc[2] {
					b.WriteString(s[i+loc[2] : i+loc[3]])
				}
				b.WriteString(replacement)
				if loc[4] >= 0 && loc[5] > loc[4] {
					b.WriteString(s[i+loc[4] : i+loc[5]])
				}
				i += loc[1]
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// replaceOutsideStrings replaces all matches of re in s, but only in parts
// outside single-quoted string literals.
func replaceOutsideStrings(s string, re *regexp.Regexp, replacement string) string {
	var b strings.Builder
	inString := false
	i := 0
	for i < len(s) {
		if s[i] == '\'' {
			// Check for escaped quote ''
			if inString && i+1 < len(s) && s[i+1] == '\'' {
				b.WriteString("''")
				i += 2
				continue
			}
			b.WriteByte(s[i])
			inString = !inString
			i++
			continue
		}
		if !inString {
			// Try to match the pattern at this position
			loc := re.FindStringIndex(s[i:])
			if loc != nil && loc[0] == 0 {
				b.WriteString(replacement)
				i += loc[1]
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

func fieldToSQLString(f *Field) string {
	if f == nil {
		return "NULL"
	}
	if f.IsNull != nil && *f.IsNull {
		return "NULL"
	}
	if f.StringValue != nil {
		return "'" + strings.ReplaceAll(*f.StringValue, "'", "''") + "'"
	}
	if f.LongValue != nil {
		return fmt.Sprintf("%d", *f.LongValue)
	}
	if f.DoubleValue != nil {
		return fmt.Sprintf("%g", *f.DoubleValue)
	}
	if f.BooleanValue != nil {
		if *f.BooleanValue {
			return "TRUE"
		}
		return "FALSE"
	}
	if f.BlobValue != nil {
		return fmt.Sprintf("'%x'", f.BlobValue)
	}
	return "NULL"
}
