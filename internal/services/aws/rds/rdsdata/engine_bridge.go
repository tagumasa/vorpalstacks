package rdsdata

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strings"
	"sync"
	"time"

	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/types"

	"vorpalstacks/internal/core/logs"
)

// AWS Data API operational constants. The binary response size cap (1 MiB)
// and the default statement timeout are documented in the AWS Data API
// reference and the Aurora User Guide ("If the binary response data from
// the database is more than 1 MB, the call is terminated."). The
// formattedRecords size limit (10 MB) is documented in the
// ExecuteStatement API reference ("The size limit for this field is
// currently 10 MB").
const (
	maxBinaryResponseBytes   = 1 << 20  // 1 MiB — binary Records + ColumnMetadata
	maxFormattedRecordsBytes = 10 << 20 // 10 MB — FormattedRecords (JSON)
	defaultStatementTimeout  = 45 * time.Second
	// maxBgStatementTime bounds how long a ContinueAfterTimeout statement
	// may run in the background after the client has received
	// StatementTimeoutException. AWS does not publish a hard ceiling for
	// Aurora Serverless v1; one hour covers all but the most extreme
	// schema migrations.
	maxBgStatementTime = time.Hour
)

// newSQLContext returns a fresh sql.Context bound to the given database.
// A new context per call mirrors the AWS Data API contract: each
// non-transactional call is independent and carries no session state.
func newSQLContext(database string) *sql.Context {
	ctx := sql.NewContext(context.Background())
	if database != "" {
		ctx.SetCurrentDatabase(database)
	}
	return ctx
}

// executeSQL runs a SQL string on the given engine and returns a formatted
// response. Behaviour follows the AWS Data API contract:
//
//   - Result-producing statements (SELECT / SHOW / DESCRIBE / EXPLAIN / WITH
//     CTEs / TABLE) populate Records and ColumnMetadata. Classification is
//     driven by the engine's returned schema (types.IsOkResultSchema) rather
//     than SQL prefix matching, so CTEs and parenthesised SELECTs are
//     detected correctly.
//   - DML / DDL statements produce an OkResult row carrying RowsAffected
//     and (for INSERT on auto-increment tables) the InsertID. The InsertID
//     is exposed via generatedFields; RowsAffected populates
//     NumberOfRecordsUpdated.
//   - The whole call is bounded by a context timeout (defaultStatementTimeout).
//     A timeout surfaces as StatementTimeoutException.
//   - The marshalled response is bounded to maxResponseBytes; an oversized
//     response surfaces as StatementTimeoutException, matching AWS which
//     terminates the call when binary response data exceeds 1 MiB.
func executeSQL(engine *sqle.Engine, sqlCtx *sql.Context, sqlStr string, includeMetadata bool, formatRecordsAs string, instanceID string) (*ExecuteStatementResponse, error) {
	return executeSQLOpts(engine, sqlCtx, sqlStr, includeMetadata, formatRecordsAs, nil, instanceID)
}

// executeSQLOpts is the option-bearing form of executeSQL. resultSetOpts is
// optional; nil leaves DECIMAL / BIGINT formatting at the engine default.
// instanceID is included in StatementTimeoutException as dbConnectionId.
func executeSQLOpts(engine *sqle.Engine, sqlCtx *sql.Context, sqlStr string, includeMetadata bool, formatRecordsAs string, resultSetOpts *ResultSetOptions, instanceID string) (*ExecuteStatementResponse, error) {
	// Bind a deadline to the underlying context so engine.Query and any
	// downstream iterators observe StatementTimeout. AWS Data API's default
	// is 45 seconds for ExecuteStatement.
	parentCtx, cancel := context.WithTimeout(sqlCtx, defaultStatementTimeout)
	defer cancel()
	sqlCtx = sqlCtx.WithContext(parentCtx)

	// Wrap engine.Query in a goroutine + select so that even if the
	// engine's analyser or planner blocks without checking the context,
	// we still return StatementTimeoutException at the deadline instead
	// of hanging indefinitely.
	type queryResult struct {
		schema sql.Schema
		iter   sql.RowIter
		err    error
	}
	qCh := make(chan queryResult, 1)
	go func() {
		defer func() {
			if re := recover(); re != nil {
				qCh <- queryResult{err: fmt.Errorf("internal panic during query planning: %v", re)}
			}
		}()
		schema, iter, _, err := engine.Query(sqlCtx, sqlStr)
		qCh <- queryResult{schema, iter, err}
	}()

	var schema sql.Schema
	var rowIter sql.RowIter
	select {
	case qr := <-qCh:
		if qr.err != nil {
			if parentCtx.Err() == context.DeadlineExceeded {
				return nil, statementTimeout(fmt.Sprintf("statement exceeded %v timeout", defaultStatementTimeout), instanceID)
			}
			return nil, qr.err
		}
		schema = qr.schema
		rowIter = qr.iter
	case <-parentCtx.Done():
		return nil, statementTimeout(fmt.Sprintf("statement exceeded %v timeout", defaultStatementTimeout), instanceID)
	}

	// Wrap RowIterToRows in the same goroutine + select pattern so that a
	// blocking iterator (e.g. a full-table scan on a very large Pebble
	// store) cannot hang past the deadline.
	type rowsResult struct {
		rows []sql.Row
		err  error
	}
	rCh := make(chan rowsResult, 1)
	go func() {
		defer func() {
			if re := recover(); re != nil {
				rCh <- rowsResult{err: fmt.Errorf("internal panic during row iteration: %v", re)}
			}
		}()
		var rows []sql.Row
		var err error
		if rowIter != nil {
			rows, err = sql.RowIterToRows(sqlCtx, rowIter)
		}
		rCh <- rowsResult{rows, err}
	}()

	var rows []sql.Row
	select {
	case rr := <-rCh:
		if rr.err != nil {
			if parentCtx.Err() == context.DeadlineExceeded {
				return nil, statementTimeout(fmt.Sprintf("statement exceeded %v timeout", defaultStatementTimeout), instanceID)
			}
			return nil, fmt.Errorf("failed to read rows: %w", rr.err)
		}
		rows = rr.rows
	case <-parentCtx.Done():
		// Force-close the iterator to unblock the goroutine. The
		// goroutine may still leak if Close itself blocks, but many
		// iterators check for closure between rows.
		if rowIter != nil {
			rowIter.Close(sqlCtx)
		}
		return nil, statementTimeout(fmt.Sprintf("statement exceeded %v timeout", defaultStatementTimeout), instanceID)
	}

	resp := &ExecuteStatementResponse{}

	// AWS spec: formatRecordsAs only applies to SELECT-style statements.
	// DDL/DML do not produce a result set to format.
	isResultProducing := schema != nil && len(schema) > 0 && !types.IsOkResultSchema(schema)

	if isResultProducing {
		if strings.EqualFold(formatRecordsAs, "JSON") {
			// AWS docs: Records and ColumnMetadata are blank when
			// formatRecordsAs is set to JSON; only FormattedRecords
			// is populated.
			formatted, jerr := formatAsJSON(rows, schema)
			if jerr != nil {
				logs.Warn("rdsdata: failed to format records as JSON",
					logs.Err(jerr),
					logs.Int("rows", len(rows)))
			} else {
				resp.FormattedRecords = formatted
			}
		} else {
			resp.Records = convertRows(rows, schema, resultSetOpts)
			if includeMetadata {
				resp.ColumnMetadata = convertSchema(schema)
			}
		}
	} else {
		// DML / DDL: the engine yields a single OkResult row carrying
		// RowsAffected and the last InsertID (auto-increment). Fall back to
		// 'len(rows)' only if the engine neglected to emit an OkResult.
		var insertID uint64
		resp.NumberOfRecordsUpdated = int64(len(rows))
		if len(rows) > 0 && types.IsOkResult(rows[0]) {
			ok := types.GetOkResult(rows[0])
			resp.NumberOfRecordsUpdated = int64(ok.RowsAffected)
			insertID = ok.InsertID
		}
		if insertID > 0 {
			id := int64(insertID)
			resp.GeneratedFields = []Field{{LongValue: &id}}
		}
	}

	// AWS spec: "If the binary response data from the database is more
	// than 1 MB, the call is terminated." This applies to the binary
	// Records + ColumnMetadata payload. FormattedRecords (JSON) has a
	// separate 10 MB limit per the ExecuteStatement API reference.
	if approxBinarySize(resp) > maxBinaryResponseBytes {
		return nil, statementTimeout(fmt.Sprintf("binary response size exceeds %d bytes", maxBinaryResponseBytes), instanceID)
	}
	if len(resp.FormattedRecords) > maxFormattedRecordsBytes {
		return nil, statementTimeout(fmt.Sprintf("formattedRecords size exceeds %d bytes", maxFormattedRecordsBytes), instanceID)
	}

	return resp, nil
}

// approxBinarySize estimates the marshalled size of the binary response
// payload (Records + ColumnMetadata + GeneratedFields) without performing
// the full JSON encode. FormattedRecords is excluded — it has a separate
// 10 MB cap enforced by the caller. Each Field is sized by its actual
// value type: StringValue/BlobValue contribute their byte length,
// LongValue up to 20 bytes (int64 max digits + sign), DoubleValue up to
// 24 bytes (float64 repr), BooleanValue 5 bytes. Column-metadata entries
// are sized at 256 bytes to account for type/name/label strings. The
// estimate is conservative (slightly over) so genuinely oversized
// payloads are caught.
func approxBinarySize(resp *ExecuteStatementResponse) int {
	total := 0
	for _, row := range resp.Records {
		for _, f := range row {
			if f.StringValue != nil {
				total += len(*f.StringValue) + 2 // quotes
			} else if f.BlobValue != nil {
				total += len(f.BlobValue)*2 + 2 // hex-encoded + quotes
			} else if f.LongValue != nil {
				total += 20 // max int64 digits + sign
			} else if f.DoubleValue != nil {
				total += 24 // float64 repr
			} else if f.BooleanValue != nil {
				total += 5 // "true"/"false"
			} else {
				total += 4 // "null"
			}
		}
	}
	total += len(resp.ColumnMetadata) * 256
	for _, f := range resp.GeneratedFields {
		if f.StringValue != nil {
			total += len(*f.StringValue) + 2
		} else {
			total += 20
		}
	}
	return total
}

// convertRows converts sql.Rows to RDS Data API Field format. When
// resultSetOpts is supplied, DECIMAL columns are stringified and BIGINT
// columns can be returned as strings per the AWS-spec LongReturnType
// / DecimalReturnType directives.
func convertRows(rows []sql.Row, schema sql.Schema, opts *ResultSetOptions) [][]Field {
	result := make([][]Field, len(rows))
	for i, row := range rows {
		fields := make([]Field, len(row))
		for j, val := range row {
			fields[j] = convertValue(val, schema, j, opts)
		}
		result[i] = fields
	}
	return result
}

// convertSchema converts sql.Schema to ColumnMetadata, populating all
// 14 Smithy fields. Precision and Scale are extracted from DECIMAL types
// via the sql.DecimalType interface. ArrayBaseColumnType is always 0 (MySQL
// has no ARRAY type). IsCurrency is always false (MySQL has no native
// currency type).
func convertSchema(schema sql.Schema) []ColumnMetadata {
	if schema == nil {
		return nil
	}
	metadata := make([]ColumnMetadata, len(schema))
	for i, col := range schema {
		cm := ColumnMetadata{
			Name:            col.Name,
			Label:           col.Name,
			IsSigned:        isSignedType(col.Type),
			IsCaseSensitive: isCaseSensitiveType(col.Type),
			Nullable:        nullableToInt(col.Nullable),
			Type:            sqlTypeToRDS(col.Type),
			TypeName:        col.Type.String(),
			IsAutoIncrement: col.AutoIncrement,
			TableName:       col.Source,
			SchemaName:      col.DatabaseSource,
		}
		if dt, ok := col.Type.(sql.DecimalType); ok {
			cm.Precision = int32(dt.Precision())
			cm.Scale = int32(dt.Scale())
		}
		metadata[i] = cm
	}
	return metadata
}

// convertValue converts a Go value to an RDS Data API Field.
//
// When schema / colIdx are available and opts requests a particular return
// type, DECIMAL and BIGINT columns are stringified per the AWS spec:
//
//   - DecimalReturnType = "STRING"  -> DECIMAL values returned as StringValue
//   - LongReturnType    = "STRING"  -> BIGINT values returned as StringValue
//
// The default (no opts, or "DOUBLE_VALUE"/"LONG_VALUE") keeps the native
// numeric encoding.
//
// uint64: values exceeding int64's positive range still fit in LongValue
// as a negative int64, mirroring what JDBC clients do when reading
// BIGINT UNSIGNED via DatabaseMetaData. Callers that need the unsigned
// value should request LongReturnType=STRING.
func convertValue(val interface{}, schema sql.Schema, colIdx int, opts *ResultSetOptions) Field {
	if val == nil {
		t := true
		return Field{IsNull: &t}
	}

	// Apply DecimalReturnType / LongReturnType directives when we know the
	// originating column. Without schema information (e.g. for values that
	// bypassed convertRows) the directives are no-ops.
	//
	// AWS defaults (ResultSetOptions docs):
	//   - DecimalReturnType default = "STRING"
	//   - LongReturnType default    = "LONG"
	if schema != nil && colIdx < len(schema) {
		typeStr := ""
		if schema[colIdx].Type != nil {
			typeStr = schema[colIdx].Type.String()
		}

		decimalReturn := "STRING"
		if opts != nil && opts.DecimalReturnType != "" {
			decimalReturn = strings.ToUpper(opts.DecimalReturnType)
		}
		longReturn := "LONG"
		if opts != nil && opts.LongReturnType != "" {
			longReturn = strings.ToUpper(opts.LongReturnType)
		}

		switch {
		case strings.Contains(typeStr, "DECIMAL") && decimalReturn == "STRING":
			sv := fmtDecimalString(val)
			return Field{StringValue: &sv}
		case strings.Contains(typeStr, "BIGINT") && longReturn == "STRING":
			sv := fmt.Sprintf("%d", toInt64(val))
			return Field{StringValue: &sv}
		}
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
	case uint64:
		// uint64 is encoded as int64 (with wrap-around for values
		// exceeding MaxInt64) so the response shape matches what
		// AWS-spec clients expect for integer-typed columns. Callers
		// that need the unsigned magnitude should request
		// LongReturnType=STRING.
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

// toInt64 coerces integer-like Go values to int64. Used for LongReturnType
// stringification where the engine may return any of the int*/uint* types
// for a BIGINT column depending on the operation path.
func toInt64(val interface{}) int64 {
	switch v := val.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	}
	return 0
}

// fmtDecimalString renders a DECIMAL value to its canonical string form
// without using %v (which would lose precision through float coercion).
// strings.TrimRight trims trailing zeros after the decimal point while
// keeping at least one fractional digit for fractional inputs.
func fmtDecimalString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// fieldToValue converts a Field to a Value (deprecated ExecuteSql format).
// The deprecated Value type has finer-grained numeric variants (IntValue
// for int32, BigIntValue for int64, RealValue for float32) that do not
// exist in the Field union. We populate them from the Field's LongValue
// and DoubleValue, casting down where the value fits.
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
		bv := *f.LongValue
		v.BigIntValue = &bv
		if *f.LongValue >= math.MinInt32 && *f.LongValue <= math.MaxInt32 {
			iv := int32(*f.LongValue)
			v.IntValue = &iv
		}
	}
	if f.DoubleValue != nil {
		v.DoubleValue = f.DoubleValue
		rv := float32(*f.DoubleValue)
		v.RealValue = &rv
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
// The ARN must use the canonical colon-separated form
// "arn:aws:rds:<region>:<account>:<resource-type>:<identifier>". Any other
// prefix (or fewer than seven colon-separated parts) yields ("", "") so
// resolveEngine surfaces InvalidParameterException to the caller.
func parseArn(arn string) (string, string) {
	if !strings.HasPrefix(arn, "arn:aws:rds:") && !strings.HasPrefix(arn, "arn:aws-cn:rds:") && !strings.HasPrefix(arn, "arn:aws-us-gov:rds:") {
		return "", ""
	}
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

// splitSQL splits a multi-statement SQL string by semicolons, respecting
// single-quoted string literals, double-quoted identifiers, MySQL
// backtick-quoted identifiers, line comments (--), block comments (/* */),
// and backslash escapes inside single-quoted strings. Semicolons inside
// any of these contexts are not treated as statement separators.
//
// Dollar-quoting ($$...$$) is intentionally NOT handled: it is
// PostgreSQL-only syntax and does not exist in MySQL. MySQL identifiers
// may legally contain '$' (per the manual, "Schema Object Names"), and
// the sequence '$$' can appear inside a valid identifier such as
// `price$$amount`. Detecting '$$' as a delimiter would corrupt
// multi-statement batches and could shift statement boundaries in ways
// that enable SQL injection.
func splitSQL(sqlStr string) []string {
	var statements []string
	var buf strings.Builder
	inSingle := false
	inDouble := false
	inBacktick := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(sqlStr); i++ {
		ch := sqlStr[i]
		// Honour backslash escapes inside single-quoted strings so that
		// '\;' does not prematurely terminate the statement. MySQL's
		// default mode interprets '\' as an escape character.
		if inSingle && ch == '\\' && i+1 < len(sqlStr) {
			buf.WriteByte(ch)
			buf.WriteByte(sqlStr[i+1])
			i++
			continue
		}

		switch {
		case inLineComment:
			buf.WriteByte(ch)
			if ch == '\n' {
				inLineComment = false
			}

		case inBlockComment:
			buf.WriteByte(ch)
			if ch == '*' && i+1 < len(sqlStr) && sqlStr[i+1] == '/' {
				buf.WriteByte(sqlStr[i+1])
				i++
				inBlockComment = false
			}

		case inSingle:
			buf.WriteByte(ch)
			if ch == '\'' {
				if i+1 < len(sqlStr) && sqlStr[i+1] == '\'' {
					buf.WriteByte(sqlStr[i+1])
					i++
				} else {
					inSingle = false
				}
			}

		case inDouble:
			buf.WriteByte(ch)
			if ch == '"' {
				if i+1 < len(sqlStr) && sqlStr[i+1] == '"' {
					buf.WriteByte(sqlStr[i+1])
					i++
				} else {
					inDouble = false
				}
			}

		case inBacktick:
			buf.WriteByte(ch)
			if ch == '`' {
				if i+1 < len(sqlStr) && sqlStr[i+1] == '`' {
					buf.WriteByte(sqlStr[i+1])
					i++
				} else {
					inBacktick = false
				}
			}

		case ch == '-' && i+1 < len(sqlStr) && sqlStr[i+1] == '-':
			inLineComment = true
			buf.WriteByte(ch)
			buf.WriteByte(sqlStr[i+1])
			i++

		case ch == '/' && i+1 < len(sqlStr) && sqlStr[i+1] == '*':
			inBlockComment = true
			buf.WriteByte(ch)
			buf.WriteByte(sqlStr[i+1])
			i++

		case ch == '\'':
			inSingle = true
			buf.WriteByte(ch)

		case ch == '"':
			inDouble = true
			buf.WriteByte(ch)

		case ch == '`':
			inBacktick = true
			buf.WriteByte(ch)

		case ch == ';':
			stmt := strings.TrimSpace(buf.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			buf.Reset()

		default:
			buf.WriteByte(ch)
		}
	}

	stmt := strings.TrimSpace(buf.String())
	if stmt != "" {
		statements = append(statements, stmt)
	}
	return statements
}

// validateParameters ensures every SqlParameter is well-formed and that no
// array-valued parameters are present (AWS spec: 'Array parameters are not
// supported'). It also validates that TypeHint=DECIMAL StringValues match
// a strict numeric pattern and TypeHint=JSON StringValues are valid JSON,
// preventing SQL injection via raw string concatenation in
// fieldToSQLString. The AWS Data API Field is a union — at most one value
// discriminator may be set. Returns InvalidParameterException on the first
// violation.
func validateParameters(params []SqlParameter) error {
	for _, p := range params {
		if p.Name == "" {
			return invalidParam("parameter name is required")
		}
		if p.Value == nil {
			return invalidParam(fmt.Sprintf("parameter %q has no value", p.Name))
		}
		// Enforce the Field union constraint: at most one discriminator
		// may be set. AWS treats multiple simultaneous values as a
		// malformed request.
		setCount := 0
		if p.Value.IsNull != nil {
			setCount++
		}
		if p.Value.StringValue != nil {
			setCount++
		}
		if p.Value.LongValue != nil {
			setCount++
		}
		if p.Value.DoubleValue != nil {
			setCount++
		}
		if p.Value.BooleanValue != nil {
			setCount++
		}
		if len(p.Value.BlobValue) > 0 {
			setCount++
		}
		if p.Value.ArrayValue != nil {
			setCount++
		}
		if setCount > 1 {
			return invalidParam(fmt.Sprintf("parameter %q has multiple value discriminators set; Field is a union — set at most one", p.Name))
		}
		if p.Value.ArrayValue != nil {
			return invalidParam(fmt.Sprintf("array parameters are not supported (parameter %q)", p.Name))
		}
		// TypeHint-specific validation: DECIMAL and JSON TypeHints accept
		// StringValue but render it into SQL without single-quote
		// escaping (DECIMAL) or with only string-literal escaping (JSON).
		// DECIMAL must be a strict numeric literal to prevent injection;
		// JSON must be syntactically valid JSON.
		hint := strings.ToUpper(p.TypeHint)
		if hint == "DECIMAL" && p.Value.StringValue != nil {
			if !decimalLiteralPattern.MatchString(*p.Value.StringValue) {
				return invalidParam(fmt.Sprintf(
					"parameter %q with typeHint DECIMAL is not a valid numeric literal", p.Name))
			}
		}
		if hint == "JSON" && p.Value.StringValue != nil {
			if !json.Valid([]byte(*p.Value.StringValue)) {
				return invalidParam(fmt.Sprintf(
					"parameter %q with typeHint JSON is not valid JSON", p.Name))
			}
		}
	}
	return nil
}

// decimalLiteralPattern matches the textual representation that the AWS
// Data API accepts for TypeHint=DECIMAL: an optional sign, followed by
// digits with an optional decimal point. This is the same format that
// MySQL's DECIMAL parser accepts and prevents any non-numeric character
// from reaching the SQL stream.
var decimalLiteralPattern = regexp.MustCompile(`^[+-]?(\d+(\.\d*)?|\.\d+)$`)

// paramReCache caches compiled substituteParameters regexes by sorted
// parameter-name key. Building the alternation regex is the dominant cost
// of substituteParameters; caching eliminates re-compilation on repeated
// calls with the same parameter-name set (the common case in batch and
// loop callers).
var paramReCache sync.Map // key: string → value: *regexp.Regexp

func getParamRegex(names []string) *regexp.Regexp {
	key := strings.Join(names, "\x00")
	if cached, ok := paramReCache.Load(key); ok {
		return cached.(*regexp.Regexp)
	}
	re := regexp.MustCompile(`(^|\W):(` + strings.Join(names, "|") + `)($|\W)`)
	actual, _ := paramReCache.LoadOrStore(key, re)
	return actual.(*regexp.Regexp)
}

// substituteParameters replaces named parameters (:name) in SQL with values.
// It builds a single alternation regex covering every named parameter and
// applies it once, which avoids the per-parameter regexp.MustCompile cost
// of compiling one regex per parameter. Single-quoted string literals are
// skipped so that parameter-like substrings inside literals are not
// substituted. The compiled regex is cached across calls with the same
// parameter-name set so callers that reuse the same parameter list do not
// pay the regex-compile cost on every invocation.
//
// The replacement string is produced via fieldToSQLString, which:
//   - backslash-escapes '\' and single-quotes inside string values so the
//     substituted literal remains syntactically valid SQL
//   - honours TypeHint for DATE / DECIMAL / JSON / TIME / TIMESTAMP values
//     by rendering StringValue as the canonical literal form for that type
func substituteParameters(sqlStr string, params []SqlParameter) string {
	if len(params) == 0 {
		return sqlStr
	}
	names := make([]string, 0, len(params))
	values := make(map[string]*Field, len(params))
	hints := make(map[string]string, len(params))
	for _, p := range params {
		if p.Name == "" || p.Value == nil {
			continue
		}
		if _, exists := values[p.Name]; exists {
			continue
		}
		names = append(names, regexp.QuoteMeta(p.Name))
		values[p.Name] = p.Value
		hints[p.Name] = p.TypeHint
	}
	if len(names) == 0 {
		return sqlStr
	}
	pattern := getParamRegex(names)
	return replaceOutsideStringsWithCaptureMulti(sqlStr, pattern, func(name string) string {
		return fieldToSQLString(values[name], hints[name])
	})
}

// replaceOutsideStringsWithCaptureMulti is the multi-parameter generalisation
// of the previous replaceOutsideStringsWithCapture. The callback receives
// the captured parameter name (without the leading colon) and returns the
// substitution.
func replaceOutsideStringsWithCaptureMulti(s string, re *regexp.Regexp, replace func(name string) string) string {
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
		// Honour backslash escapes inside string literals so the next
		// character is consumed verbatim.
		if inString && s[i] == '\\' && i+1 < len(s) {
			b.WriteByte(s[i])
			b.WriteByte(s[i+1])
			i += 2
			continue
		}
		if !inString {
			loc := re.FindStringSubmatchIndex(s[i:])
			if loc != nil && loc[0] == 0 {
				// loc[2:4]  = group 1 (leading boundary)
				// loc[4:6]  = group 2 (parameter name)
				// loc[6:8]  = group 3 (trailing boundary)
				if loc[2] >= 0 && loc[3] > loc[2] {
					b.WriteString(s[i+loc[2] : i+loc[3]])
				}
				name := s[i+loc[4] : i+loc[5]]
				b.WriteString(replace(name))
				if loc[6] >= 0 && loc[7] > loc[6] {
					b.WriteString(s[i+loc[6] : i+loc[7]])
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

// fieldToSQLString renders a Field as a SQL literal suitable for inline
// substitution. The TypeHint influences the rendering for DATE / TIME /
// TIMESTAMP / DECIMAL / JSON values, matching the AWS Data API contract
// (SqlParameter.typeHint documentation).
//
// String values are escaped for the MySQL default mode:
//   - backslash is doubled ('\' -> '\\')
//   - single quote is doubled ("'" -> "”")
//   - NUL is replaced with '\0'
//
// Without the backslash escape, a hostile parameter such as "\'" could
// prematurely close the surrounding string literal and inject SQL.
func fieldToSQLString(f *Field, typeHint string) string {
	if f == nil {
		return "NULL"
	}
	if f.IsNull != nil && *f.IsNull {
		return "NULL"
	}
	// TypeHint-aware rendering (M-3): when the caller marks a value as
	// DATE / TIME / TIMESTAMP / DECIMAL / JSON, treat StringValue as the
	// canonical literal form and render accordingly. The default case
	// (no TypeHint) falls through to the value-type switch below.
	switch strings.ToUpper(typeHint) {
	case "DATE", "TIME", "TIMESTAMP":
		if f.StringValue != nil {
			return "'" + escapeMySQLLiteral(*f.StringValue) + "'"
		}
	case "DECIMAL":
		// DECIMAL literals must not be quoted; the engine parses the
		// textual representation directly without precision loss.
		if f.StringValue != nil {
			return *f.StringValue
		}
		if f.LongValue != nil {
			return fmt.Sprintf("%d", *f.LongValue)
		}
		if f.DoubleValue != nil {
			return fmt.Sprintf("%g", *f.DoubleValue)
		}
	case "JSON":
		if f.StringValue != nil {
			return "'" + escapeMySQLLiteral(*f.StringValue) + "'"
		}
	}

	if f.StringValue != nil {
		return "'" + escapeMySQLLiteral(*f.StringValue) + "'"
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

// escapeMySQLLiteral escapes a string for safe inclusion in a MySQL
// single-quoted literal under the default sql_mode (where '\' is an escape
// character). Single quotes, backslashes, and NUL bytes are escaped so
// that hostile parameter values cannot prematurely close the surrounding
// string literal and inject SQL.
func escapeMySQLLiteral(s string) string {
	// Pre-size for the worst case (every char doubled).
	out := make([]byte, 0, len(s)*2)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\'':
			out = append(out, '\'', '\'')
		case '\\':
			out = append(out, '\\', '\\')
		case 0:
			out = append(out, '\\', '0')
		default:
			out = append(out, s[i])
		}
	}
	return string(out)
}
