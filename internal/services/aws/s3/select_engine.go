// Package s3 provides S3 service operations for vorpalstacks.
package s3

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"vorpalstacks/internal/core/logs"
	"vorpalstacks/pkg/sqlparser"

	"github.com/dgraph-io/ristretto/v2"
)

var (
	likePatternCache     *ristretto.Cache[string, *regexp.Regexp]
	likePatternCacheOnce sync.Once
)

func getLikePatternCache() *ristretto.Cache[string, *regexp.Regexp] {
	likePatternCacheOnce.Do(func() {
		cache, err := ristretto.NewCache(&ristretto.Config[string, *regexp.Regexp]{
			NumCounters: 1000,
			MaxCost:     5 << 20,
			BufferItems: 64,
		})
		if err != nil {
			logs.Error("failed to create LIKE pattern cache, regex caching disabled", logs.Err(err))
			return
		}
		likePatternCache = cache
	})
	return likePatternCache
}

// SelectEngine executes S3 Select queries on object content.
type SelectEngine struct {
	input *SelectObjectContentInput
	stmt  *sqlparser.Select
}

// NewSelectEngine creates a new SelectEngine for executing S3 Select queries.
func NewSelectEngine(input *SelectObjectContentInput) (*SelectEngine, error) {
	opts := sqlparser.ParserOptions{Dialect: sqlparser.DialectS3Select}
	stmt, err := sqlparser.ParseWithOptions(input.Expression, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL expression: %w", err)
	}

	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, fmt.Errorf("only SELECT statements are supported")
	}

	return &SelectEngine{
		input: input,
		stmt:  selectStmt,
	}, nil
}

// Execute runs the S3 Select query on the provided data.
func (e *SelectEngine) Execute(ctx context.Context, reader io.Reader, writer EventStreamWriter) error {
	var r io.Reader = reader

	switch e.input.InputSerialization.CompressionType {
	case CompressionTypeGZIP:
		gzReader, err := gzip.NewReader(reader)
		if err != nil {
			return fmt.Errorf("failed to create gzip reader: %w", err)
		}
		defer gzReader.Close()
		r = gzReader
	case CompressionTypeBZIP2:
		r = bzip2.NewReader(reader)
	}

	if e.input.InputSerialization.CSV != nil {
		return e.executeCSV(ctx, r, writer)
	}
	if e.input.InputSerialization.JSON != nil {
		return e.executeJSON(ctx, r, writer)
	}

	return fmt.Errorf("unsupported input serialization format")
}

func (e *SelectEngine) executeCSV(ctx context.Context, reader io.Reader, writer EventStreamWriter) error {
	csvInput := e.input.InputSerialization.CSV

	fieldDelimiter := ","
	if csvInput.FieldDelimiter != "" {
		fieldDelimiter = csvInput.FieldDelimiter
	}

	recordDelimiter := "\n"
	if csvInput.RecordDelimiter != "" {
		recordDelimiter = csvInput.RecordDelimiter
	}

	csvReader := csv.NewReader(reader)
	csvReader.Comma = rune(fieldDelimiter[0])
	csvReader.LazyQuotes = true

	var outputDelimiter string = "\n"
	var outputFieldDelimiter string = ","
	if e.input.OutputSerialization.CSV != nil {
		if e.input.OutputSerialization.CSV.RecordDelimiter != "" {
			outputDelimiter = e.input.OutputSerialization.CSV.RecordDelimiter
		}
		if e.input.OutputSerialization.CSV.FieldDelimiter != "" {
			outputFieldDelimiter = e.input.OutputSerialization.CSV.FieldDelimiter
		}
	} else if e.input.OutputSerialization.JSON != nil {
		if e.input.OutputSerialization.JSON.RecordDelimiter != "" {
			outputDelimiter = e.input.OutputSerialization.JSON.RecordDelimiter
		}
	}

	var headers []string
	rowNum := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		rowNum++
		rawBytes := strings.Join(record, fieldDelimiter) + recordDelimiter
		writer.AddBytesScanned(int64(len(rawBytes)))
		writer.AddBytesProcessed(int64(len(rawBytes)))

		if rowNum == 1 && csvInput.FileHeaderInfo == FileHeaderInfoUse {
			headers = record
			continue
		}
		if rowNum == 1 && csvInput.FileHeaderInfo == FileHeaderInfoIgnore {
			continue
		}

		row := e.csvRowToMap(record, headers)

		if !e.evaluateWhere(row) {
			continue
		}

		outputRecord := e.selectFields(record, headers)

		var outputData []byte
		if e.input.OutputSerialization.JSON != nil {
			jsonRecord := make(map[string]interface{})
			for i, val := range outputRecord {
				key := fmt.Sprintf("_%d", i+1)
				if i < len(headers) && headers[i] != "" {
					key = headers[i]
				}
				jsonRecord[key] = e.parseValue(val)
			}
			outputData, err = marshalJSONRecord(jsonRecord, outputDelimiter)
		} else {
			outputData, err = marshalCSVRecord(outputRecord, outputDelimiter, outputFieldDelimiter)
		}

		if err != nil {
			continue
		}

		if err := writer.WriteRecords(outputData); err != nil {
			return err
		}

		if err := writer.MaybeWriteProgress(); err != nil {
			return err
		}
	}

	return nil
}

func (e *SelectEngine) csvRowToMap(record []string, headers []string) map[string]interface{} {
	row := make(map[string]interface{})
	for i, val := range record {
		key := fmt.Sprintf("_%d", i+1)
		if i < len(headers) && headers[i] != "" {
			row[headers[i]] = e.parseValue(val)
		}
		row[key] = e.parseValue(val)
	}
	return row
}

func (e *SelectEngine) parseValue(s string) interface{} {
	if s == "" {
		return nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	if b, err := strconv.ParseBool(s); err == nil {
		return b
	}
	return s
}

func (e *SelectEngine) executeJSON(ctx context.Context, reader io.Reader, writer EventStreamWriter) error {
	var outputDelimiter string = "\n"
	if e.input.OutputSerialization.JSON != nil {
		if e.input.OutputSerialization.JSON.RecordDelimiter != "" {
			outputDelimiter = e.input.OutputSerialization.JSON.RecordDelimiter
		}
	} else if e.input.OutputSerialization.CSV != nil {
		if e.input.OutputSerialization.CSV.RecordDelimiter != "" {
			outputDelimiter = e.input.OutputSerialization.CSV.RecordDelimiter
		}
	}

	scanner := bufio.NewScanner(reader)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)

	rowNum := 0

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Bytes()
		rowNum++

		rawBytes := len(line) + 1
		writer.AddBytesScanned(int64(rawBytes))
		writer.AddBytesProcessed(int64(rawBytes))

		var row map[string]interface{}
		if err := json.Unmarshal(line, &row); err != nil {
			continue
		}

		flattened := e.flattenJSON(row, "")

		if !e.evaluateWhere(flattened) {
			continue
		}

		outputRecord := e.selectJSONFields(row, flattened)

		var outputData []byte
		var err error
		if e.input.OutputSerialization.CSV != nil {
			fields := e.jsonToCSVFields(outputRecord)
			fieldDelimiter := e.input.OutputSerialization.CSV.FieldDelimiter
			if fieldDelimiter == "" {
				fieldDelimiter = ","
			}
			outputData, err = marshalCSVRecord(fields, outputDelimiter, fieldDelimiter)
		} else {
			outputData, err = marshalJSONRecord(outputRecord, outputDelimiter)
		}

		if err != nil {
			continue
		}

		if err := writer.WriteRecords(outputData); err != nil {
			return err
		}

		if err := writer.MaybeWriteProgress(); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func (e *SelectEngine) flattenJSON(obj map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range obj {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case map[string]interface{}:
			nested := e.flattenJSON(val, key)
			for nk, nv := range nested {
				result[nk] = nv
			}
		default:
			result[key] = v
		}
		result[k] = v
	}
	return result
}

func (e *SelectEngine) evaluateWhere(row map[string]interface{}) bool {
	if e.stmt.Where == nil {
		return true
	}
	return e.evaluateExpr(e.stmt.Where.Expr, row)
}

func (e *SelectEngine) evaluateExpr(expr sqlparser.Expr, row map[string]interface{}) bool {
	switch ex := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return e.evaluateComparison(ex, row)
	case *sqlparser.AndExpr:
		return e.evaluateExpr(ex.Left, row) && e.evaluateExpr(ex.Right, row)
	case *sqlparser.OrExpr:
		return e.evaluateExpr(ex.Left, row) || e.evaluateExpr(ex.Right, row)
	case *sqlparser.NotExpr:
		return !e.evaluateExpr(ex.Expr, row)
	case *sqlparser.ParenExpr:
		return e.evaluateExpr(ex.Expr, row)
	case *sqlparser.IsExpr:
		return e.evaluateIsExpr(ex, row)
	case *sqlparser.BoolVal:
		return bool(*ex)
	default:
		return true
	}
}

func (e *SelectEngine) evaluateComparison(expr *sqlparser.ComparisonExpr, row map[string]interface{}) bool {
	left := e.getExprValue(expr.Left, row)
	right := e.getExprValue(expr.Right, row)

	switch expr.Operator {
	case sqlparser.EqualStr:
		return e.compareEqual(left, right)
	case sqlparser.NotEqualStr:
		return !e.compareEqual(left, right)
	case sqlparser.LessThanStr:
		return e.compareLess(left, right)
	case sqlparser.LessEqualStr:
		return e.compareLessEqual(left, right)
	case sqlparser.GreaterThanStr:
		return e.compareGreater(left, right)
	case sqlparser.GreaterEqualStr:
		return e.compareGreaterEqual(left, right)
	case sqlparser.LikeStr:
		return e.matchLike(fmt.Sprintf("%v", left), fmt.Sprintf("%v", right))
	case sqlparser.NotLikeStr:
		return !e.matchLike(fmt.Sprintf("%v", left), fmt.Sprintf("%v", right))
	case sqlparser.InStr:
		return e.evaluateIn(left, expr.Right, row)
	case sqlparser.NotInStr:
		return !e.evaluateIn(left, expr.Right, row)
	}
	return false
}

func (e *SelectEngine) getExprValue(expr sqlparser.Expr, row map[string]interface{}) interface{} {
	switch ex := expr.(type) {
	case *sqlparser.ColName:
		name := ex.Name.String()
		if val, ok := row[name]; ok {
			return val
		}
		return nil
	case *sqlparser.SQLVal:
		return e.sqlValToInterface(ex)
	case *sqlparser.NullVal:
		return nil
	case sqlparser.BoolVal:
		return bool(ex)
	default:
		return nil
	}
}

func (e *SelectEngine) sqlValToInterface(v *sqlparser.SQLVal) interface{} {
	s := string(v.Val)
	switch v.Type {
	case sqlparser.StrVal:
		return s
	case sqlparser.IntVal:
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return i
		}
		return s
	case sqlparser.FloatVal:
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
		return s
	default:
		return s
	}
}

func (e *SelectEngine) compareEqual(left, right interface{}) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}

	switch l := left.(type) {
	case string:
		if r, ok := right.(string); ok {
			return l == r
		}
	case int64:
		if r, ok := right.(int64); ok {
			return l == r
		}
		if r, ok := right.(float64); ok {
			return float64(l) == r
		}
	case float64:
		if r, ok := right.(float64); ok {
			return l == r
		}
		if r, ok := right.(int64); ok {
			return l == float64(r)
		}
	case bool:
		if r, ok := right.(bool); ok {
			return l == r
		}
	}

	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

func (e *SelectEngine) compareLess(left, right interface{}) bool {
	l, lok := toFloat(left)
	r, rok := toFloat(right)
	if lok && rok {
		return l < r
	}
	lStr, lokStr := left.(string)
	rStr, rokStr := right.(string)
	if lokStr && rokStr {
		return lStr < rStr
	}
	return false
}

func (e *SelectEngine) compareLessEqual(left, right interface{}) bool {
	l, lok := toFloat(left)
	r, rok := toFloat(right)
	if lok && rok {
		return l <= r
	}
	lStr, lokStr := left.(string)
	rStr, rokStr := right.(string)
	if lokStr && rokStr {
		return lStr <= rStr
	}
	return false
}

func (e *SelectEngine) compareGreater(left, right interface{}) bool {
	l, lok := toFloat(left)
	r, rok := toFloat(right)
	if lok && rok {
		return l > r
	}
	lStr, lokStr := left.(string)
	rStr, rokStr := right.(string)
	if lokStr && rokStr {
		return lStr > rStr
	}
	return false
}

func (e *SelectEngine) compareGreaterEqual(left, right interface{}) bool {
	l, lok := toFloat(left)
	r, rok := toFloat(right)
	if lok && rok {
		return l >= r
	}
	lStr, lokStr := left.(string)
	rStr, rokStr := right.(string)
	if lokStr && rokStr {
		return lStr >= rStr
	}
	return false
}

func toFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		f, err := strconv.ParseFloat(val, 64)
		return f, err == nil
	default:
		return 0, false
	}
}

func (e *SelectEngine) matchLike(s, pattern string) bool {
	cache := getLikePatternCache()

	if re, ok := cache.Get(pattern); ok {
		return re.MatchString(s)
	}

	regexPattern := e.likeToRegex(pattern)
	re, err := regexp.Compile("^" + regexPattern + "$")
	if err != nil {
		return false
	}

	cache.Set(pattern, re, 1)
	return re.MatchString(s)
}

func (e *SelectEngine) likeToRegex(pattern string) string {
	var result strings.Builder
	for _, c := range pattern {
		switch c {
		case '%':
			result.WriteString(".*")
		case '_':
			result.WriteString(".")
		case '.', '^', '$', '*', '+', '?', '(', ')', '[', ']', '{', '}', '\\', '|':
			result.WriteString("\\")
			result.WriteRune(c)
		default:
			result.WriteRune(c)
		}
	}
	return result.String()
}

func (e *SelectEngine) evaluateIn(value interface{}, list sqlparser.Expr, row map[string]interface{}) bool {
	switch l := list.(type) {
	case sqlparser.ValTuple:
		for _, item := range l {
			if e.compareEqual(value, e.getExprValue(item, row)) {
				return true
			}
		}
	}
	return false
}

func (e *SelectEngine) evaluateIsExpr(expr *sqlparser.IsExpr, row map[string]interface{}) bool {
	val := e.getExprValue(expr.Expr, row)
	switch expr.Operator {
	case sqlparser.IsNullStr:
		return val == nil
	case sqlparser.IsNotNullStr:
		return val != nil
	}
	return false
}

func (e *SelectEngine) selectFields(record []string, headers []string) []string {
	if e.stmt.SelectExprs == nil {
		return record
	}

	var result []string
	for _, expr := range e.stmt.SelectExprs {
		switch se := expr.(type) {
		case *sqlparser.StarExpr:
			return record
		case *sqlparser.AliasedExpr:
			colName := e.getColumnName(se.Expr)
			if colName == "" {
				continue
			}

			idx := e.getColumnIndex(colName, headers)
			if idx >= 0 && idx < len(record) {
				result = append(result, record[idx])
			}
		}
	}

	return result
}

func (e *SelectEngine) selectJSONFields(original map[string]interface{}, flattened map[string]interface{}) map[string]interface{} {
	if e.stmt.SelectExprs == nil {
		return original
	}

	result := make(map[string]interface{})
	for _, expr := range e.stmt.SelectExprs {
		switch se := expr.(type) {
		case *sqlparser.StarExpr:
			return original
		case *sqlparser.AliasedExpr:
			colName := e.getColumnName(se.Expr)
			if colName == "" {
				continue
			}

			if val, ok := flattened[colName]; ok {
				alias := se.As.String()
				if alias != "" {
					result[alias] = val
				} else {
					result[colName] = val
				}
			}
		}
	}

	return result
}

func (e *SelectEngine) getColumnName(expr sqlparser.Expr) string {
	switch ex := expr.(type) {
	case *sqlparser.ColName:
		return ex.Name.String()
	}
	return ""
}

func (e *SelectEngine) getColumnIndex(name string, headers []string) int {
	if len(headers) == 0 {
		if strings.HasPrefix(name, "_") {
			if idx, err := strconv.Atoi(name[1:]); err == nil {
				return idx - 1
			}
		}
		return -1
	}

	for i, h := range headers {
		if h == name {
			return i
		}
	}

	if strings.HasPrefix(name, "_") {
		if idx, err := strconv.Atoi(name[1:]); err == nil {
			return idx - 1
		}
	}

	return -1
}

func (e *SelectEngine) jsonToCSVFields(obj map[string]interface{}) []string {
	var keys []string
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var fields []string
	for _, k := range keys {
		fields = append(fields, fmt.Sprintf("%v", obj[k]))
	}
	return fields
}
