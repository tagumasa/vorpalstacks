package cloudwatchlogs

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// queryResultRow represents a single row in the query results.
type queryResultRow struct {
	fields map[string]string
}

// queryResultField represents a field-value pair in the output.
type queryResultField struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// queryStats holds statistics about the query execution.
type queryStats struct {
	recordsScanned int64
	recordsMatched int64
	bytesScanned   int64
}

// executeQuery runs a CloudWatch Logs Insights query against the given events.
// It supports the CWLI commands: fields, filter, stats, sort, limit, parse, display.
func executeQuery(queryString string, events []logEventWithContext) ([]queryResultRow, queryStats) {
	stats := queryStats{
		recordsScanned: int64(len(events)),
	}
	for _, e := range events {
		stats.bytesScanned += int64(len(e.message))
	}

	commands := parseQueryPipeline(queryString)

	rows := make([]queryResultRow, 0, len(events))
	for _, evt := range events {
		row := buildRow(evt)
		rows = append(rows, row)
	}

	for _, cmd := range commands {
		rows = applyCommand(cmd, rows)
	}

	stats.recordsMatched = int64(len(rows))
	return rows, stats
}

// logEventWithContext carries a log event with its group/stream context.
type logEventWithContext struct {
	timestamp     int64
	message       string
	ingestionTime int64
	logGroup      string
	logStream     string
}

func buildRow(evt logEventWithContext) queryResultRow {
	fields := map[string]string{
		"@timestamp":     strconv.FormatInt(evt.timestamp, 10),
		"@message":       evt.message,
		"@logStream":     evt.logStream,
		"@logGroup":      evt.logGroup,
		"@ingestionTime": strconv.FormatInt(evt.ingestionTime, 10),
	}
	return queryResultRow{fields: fields}
}

type queryCommand struct {
	name string
	args string
}

func parseQueryPipeline(queryString string) []queryCommand {
	queryString = strings.TrimSpace(queryString)
	parts := strings.Split(queryString, "|")

	var commands []queryCommand
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		spaceIdx := strings.Index(part, " ")
		if spaceIdx < 0 {
			commands = append(commands, queryCommand{name: strings.ToLower(part), args: ""})
		} else {
			name := strings.ToLower(part[:spaceIdx])
			args := strings.TrimSpace(part[spaceIdx+1:])
			commands = append(commands, queryCommand{name: name, args: args})
		}
	}
	return commands
}

func applyCommand(cmd queryCommand, rows []queryResultRow) []queryResultRow {
	switch cmd.name {
	case "fields":
		return applyFields(cmd.args, rows)
	case "filter":
		return applyFilter(cmd.args, rows)
	case "stats":
		return applyStats(cmd.args, rows)
	case "sort":
		return applySort(cmd.args, rows)
	case "limit":
		return applyLimit(cmd.args, rows)
	case "parse":
		return applyParse(cmd.args, rows)
	case "display":
		return applyDisplay(cmd.args, rows)
	default:
		return rows
	}
}

func applyFields(args string, rows []queryResultRow) []queryResultRow {
	fieldNames := splitAndTrim(args, ",")
	if len(fieldNames) == 0 {
		return rows
	}
	for i := range rows {
		filtered := make(map[string]string)
		for _, fn := range fieldNames {
			fn = strings.TrimSpace(fn)
			if v, ok := rows[i].fields[fn]; ok {
				filtered[fn] = v
			}
		}
		rows[i].fields = filtered
	}
	return rows
}

func applyFilter(args string, rows []queryResultRow) []queryResultRow {
	var result []queryResultRow
	for _, row := range rows {
		if evalFilterExpr(args, row) {
			result = append(result, row)
		}
	}
	return result
}

func evalFilterExpr(expr string, row queryResultRow) bool {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return true
	}

	if orParts := splitLogicalExpr(expr, " or "); len(orParts) > 1 {
		for _, p := range orParts {
			if evalFilterExpr(p, row) {
				return true
			}
		}
		return false
	}

	if andParts := splitLogicalExpr(expr, " and "); len(andParts) > 1 {
		for _, p := range andParts {
			if !evalFilterExpr(p, row) {
				return false
			}
		}
		return true
	}

	return evalComparison(expr, row)
}

func splitLogicalExpr(expr, sep string) []string {
	if !strings.Contains(expr, sep) {
		return nil
	}
	return strings.Split(expr, sep)
}

func evalComparison(expr string, row queryResultRow) bool {
	expr = strings.TrimSpace(expr)

	ops := []string{">=", "<=", "!=", ">", "<", "=", " like ", " not like "}
	for _, op := range ops {
		if idx := strings.Index(strings.ToLower(expr), op); idx >= 0 {
			left := strings.TrimSpace(expr[:idx])
			right := strings.TrimSpace(expr[idx+len(op):])

			leftVal := getFieldValue(left, row)
			rightVal := strings.Trim(right, "\"'")

			if op == " like " || op == " not like " {
				pattern := strings.Trim(right, "\"'")
				matched := wildcardMatch(pattern, leftVal)
				if op == " not like " {
					return !matched
				}
				return matched
			}

			if op == "=" || op == "==" {
				return leftVal == rightVal
			}
			if op == "!=" {
				return leftVal != rightVal
			}

			leftNum, err1 := strconv.ParseFloat(leftVal, 64)
			rightNum, err2 := strconv.ParseFloat(rightVal, 64)
			if err1 != nil || err2 != nil {
				return false
			}
			switch op {
			case ">":
				return leftNum > rightNum
			case "<":
				return leftNum < rightNum
			case ">=":
				return leftNum >= rightNum
			case "<=":
				return leftNum <= rightNum
			}
		}
	}
	return true
}

func getFieldValue(field string, row queryResultRow) string {
	if v, ok := row.fields[field]; ok {
		return v
	}
	if v, ok := row.fields["@"+field]; ok {
		return v
	}
	return ""
}

func wildcardMatch(pattern, text string) bool {
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = "^" + pattern + "$"
	matched, _ := regexp.MatchString(pattern, text)
	return matched
}

func applyStats(args string, rows []queryResultRow) []queryResultRow {
	args = strings.TrimSpace(args)

	byIdx := strings.Index(strings.ToLower(args), " by ")
	var aggExpr string
	var groupFields []string
	if byIdx >= 0 {
		aggExpr = strings.TrimSpace(args[:byIdx])
		groupFields = splitAndTrim(args[byIdx+4:], ",")
	} else {
		aggExpr = args
	}

	if len(groupFields) == 0 {
		result := computeAggregations(aggExpr, rows)
		return []queryResultRow{{fields: result}}
	}

	groups := make(map[string][]queryResultRow)
	var groupOrder []string
	for _, row := range rows {
		keyParts := make([]string, len(groupFields))
		for i, gf := range groupFields {
			keyParts[i] = getFieldValue(gf, row)
		}
		key := strings.Join(keyParts, "|")
		if _, exists := groups[key]; !exists {
			groupOrder = append(groupOrder, key)
		}
		groups[key] = append(groups[key], row)
	}

	var result []queryResultRow
	for _, key := range groupOrder {
		groupRows := groups[key]
		resultFields := computeAggregations(aggExpr, groupRows)
		for i, gf := range groupFields {
			val := strings.Split(key, "|")
			if i < len(val) {
				resultFields[gf] = val[i]
			}
		}
		result = append(result, queryResultRow{fields: resultFields})
	}
	return result
}

func computeAggregations(expr string, rows []queryResultRow) map[string]string {
	result := make(map[string]string)
	funcs := splitAndTrim(expr, ",")
	for _, f := range funcs {
		f = strings.TrimSpace(f)
		openIdx := strings.Index(f, "(")
		closeIdx := strings.LastIndex(f, ")")
		if openIdx < 0 || closeIdx < 0 {
			continue
		}
		funcName := strings.ToLower(strings.TrimSpace(f[:openIdx]))
		arg := strings.TrimSpace(f[openIdx+1 : closeIdx])
		alias := funcName + "(" + arg + ")"
		if asIdx := strings.Index(strings.ToLower(f), " as "); asIdx >= 0 {
			alias = strings.TrimSpace(f[asIdx+4:])
		}

		switch funcName {
		case "count":
			if arg == "*" || arg == "" {
				result[alias] = strconv.Itoa(len(rows))
			} else {
				count := 0
				for _, row := range rows {
					if v := getFieldValue(arg, row); v != "" {
						count++
					}
				}
				result[alias] = strconv.Itoa(count)
			}
		case "sum":
			result[alias] = computeNumericAgg(rows, arg, "sum")
		case "avg":
			result[alias] = computeNumericAgg(rows, arg, "avg")
		case "min":
			result[alias] = computeNumericAgg(rows, arg, "min")
		case "max":
			result[alias] = computeNumericAgg(rows, arg, "max")
		case "count_distinct":
			seen := make(map[string]bool)
			for _, row := range rows {
				v := getFieldValue(arg, row)
				if v != "" {
					seen[v] = true
				}
			}
			result[alias] = strconv.Itoa(len(seen))
		}
	}
	return result
}

func computeNumericAgg(rows []queryResultRow, field, aggType string) string {
	var values []float64
	for _, row := range rows {
		v := getFieldValue(field, row)
		if n, err := strconv.ParseFloat(v, 64); err == nil {
			values = append(values, n)
		}
	}
	if len(values) == 0 {
		return "0"
	}
	switch aggType {
	case "sum":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return strconv.FormatFloat(sum, 'f', -1, 64)
	case "avg":
		sum := 0.0
		for _, v := range values {
			sum += v
		}
		return strconv.FormatFloat(sum/float64(len(values)), 'f', -1, 64)
	case "min":
		min := values[0]
		for _, v := range values {
			if v < min {
				min = v
			}
		}
		return strconv.FormatFloat(min, 'f', -1, 64)
	case "max":
		max := values[0]
		for _, v := range values {
			if v > max {
				max = v
			}
		}
		return strconv.FormatFloat(max, 'f', -1, 64)
	}
	return "0"
}

func applySort(args string, rows []queryResultRow) []queryResultRow {
	args = strings.TrimSpace(args)
	descending := false
	if strings.HasSuffix(strings.ToLower(args), " desc") {
		descending = true
		args = strings.TrimSpace(args[:len(args)-5])
	} else if strings.HasSuffix(strings.ToLower(args), " asc") {
		args = strings.TrimSpace(args[:len(args)-4])
	}

	field := strings.TrimSpace(args)
	sort.SliceStable(rows, func(i, j int) bool {
		vi := getFieldValue(field, rows[i])
		vj := getFieldValue(field, rows[j])
		if ni, err := strconv.ParseFloat(vi, 64); err == nil {
			if nj, err := strconv.ParseFloat(vj, 64); err == nil {
				if descending {
					return ni > nj
				}
				return ni < nj
			}
		}
		if descending {
			return vi > vj
		}
		return vi < vj
	})
	return rows
}

func applyLimit(args string, rows []queryResultRow) []queryResultRow {
	n, err := strconv.Atoi(strings.TrimSpace(args))
	if err != nil || n <= 0 {
		return rows
	}
	if n > len(rows) {
		n = len(rows)
	}
	return rows[:n]
}

func applyParse(args string, rows []queryResultRow) []queryResultRow {
	args = strings.TrimSpace(args)
	asIdx := strings.Index(strings.ToLower(args), " as ")
	if asIdx < 0 {
		return rows
	}

	pattern := strings.TrimSpace(args[:asIdx])
	rest := strings.TrimSpace(args[asIdx+4:])

	fieldNames := splitAndTrim(rest, ",")
	if len(fieldNames) == 0 {
		return rows
	}

	wildcardParts := strings.Split(pattern, "*")
	regexParts := make([]string, len(wildcardParts))
	for i, p := range wildcardParts {
		regexParts[i] = regexp.QuoteMeta(p)
	}
	regexStr := strings.Join(regexParts, "(.*?)")
	if len(wildcardParts)-1 != len(fieldNames) {
		return rows
	}

	re, err := regexp.Compile(regexStr)
	if err != nil {
		return rows
	}

	for i := range rows {
		msg := rows[i].fields["@message"]
		matches := re.FindStringSubmatch(msg)
		if matches != nil {
			for j, fn := range fieldNames {
				fn = strings.TrimSpace(fn)
				if j+1 < len(matches) {
					rows[i].fields[fn] = matches[j+1]
				}
			}
		}
	}
	return rows
}

func applyDisplay(args string, rows []queryResultRow) []queryResultRow {
	fieldNames := splitAndTrim(args, ",")
	for i := range rows {
		filtered := make(map[string]string)
		for _, fn := range fieldNames {
			fn = strings.TrimSpace(fn)
			if v, ok := rows[i].fields[fn]; ok {
				filtered[fn] = v
			}
		}
		rows[i].fields = filtered
	}
	return rows
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
