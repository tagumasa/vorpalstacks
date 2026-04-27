package athena

import (
	"strconv"
	"strings"

	athenastore "vorpalstacks/internal/store/aws/athena"
)

func (s *AthenaService) extractDatabaseNameFromDDL(query, prefix1, prefix2 string) string {
	upperQuery := strings.ToUpper(query)

	var dbName string
	if strings.HasPrefix(upperQuery, prefix1) {
		rest := strings.TrimSpace(query[len(prefix1):])
		if strings.HasPrefix(strings.ToUpper(rest), "IF NOT EXISTS") {
			rest = strings.TrimSpace(rest[len("IF NOT EXISTS"):])
		}
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			dbName = strings.Trim(parts[0], "`\"';")
		}
	} else if strings.HasPrefix(upperQuery, prefix2) {
		rest := strings.TrimSpace(query[len(prefix2):])
		if strings.HasPrefix(strings.ToUpper(rest), "IF EXISTS") {
			rest = strings.TrimSpace(rest[len("IF EXISTS"):])
		}
		parts := strings.Fields(rest)
		if len(parts) > 0 {
			dbName = strings.Trim(parts[0], "`\"';")
		}
	}

	return dbName
}

func (s *AthenaService) extractTableNameFromDrop(query string) string {
	upperQuery := strings.ToUpper(query)
	if !strings.HasPrefix(upperQuery, "DROP TABLE") {
		return ""
	}

	rest := strings.TrimSpace(query[len("DROP TABLE"):])
	if strings.HasPrefix(strings.ToUpper(rest), "IF EXISTS") {
		rest = strings.TrimSpace(rest[len("IF EXISTS"):])
	}

	parts := strings.Fields(rest)
	if len(parts) > 0 {
		tableName := parts[0]
		if strings.Contains(tableName, ".") {
			parts := strings.Split(tableName, ".")
			tableName = parts[len(parts)-1]
		}
		return strings.Trim(tableName, "`\"';")
	}

	return ""
}

func (s *AthenaService) parseCreateTableStatement(query string) (string, string, []athenastore.Column) {
	upperQuery := strings.ToUpper(query)

	createIdx := strings.Index(upperQuery, "CREATE EXTERNAL TABLE")
	if createIdx != -1 {
		rest := query[createIdx+len("CREATE EXTERNAL TABLE"):]
		return s.parseCreateTableRest(rest)
	}

	createIdx = strings.Index(upperQuery, "CREATE TABLE")
	if createIdx == -1 {
		return "", "", nil
	}

	rest := query[createIdx+len("CREATE TABLE"):]
	return s.parseCreateTableRest(rest)
}

func (s *AthenaService) parseCreateTableRest(rest string) (string, string, []athenastore.Column) {
	rest = strings.TrimSpace(rest)

	if strings.HasPrefix(strings.ToUpper(rest), "IF NOT EXISTS") {
		rest = strings.TrimSpace(rest[len("IF NOT EXISTS"):])
	}

	var tableName string
	if strings.HasPrefix(rest, "(") {
		return "", "", nil
	}

	parenIdx := strings.Index(rest, "(")
	if parenIdx == -1 {
		return "", "", nil
	}

	tablePart := strings.TrimSpace(rest[:parenIdx])
	var dbName string
	if strings.Contains(tablePart, ".") {
		parts := strings.Split(tablePart, ".")
		if len(parts) == 2 {
			dbName = strings.Trim(parts[0], "`\"';")
		}
		tablePart = parts[len(parts)-1]
	}
	tableName = strings.Trim(tablePart, "`\"';")

	columnsPart := rest[parenIdx+1:]

	endParen := strings.LastIndex(columnsPart, ")")
	if endParen == -1 {
		return tableName, dbName, nil
	}

	columnsPart = columnsPart[:endParen]

	var columns []athenastore.Column
	for _, colDef := range s.splitColumnDefinitions(columnsPart) {
		colDef = strings.TrimSpace(colDef)
		if colDef == "" {
			continue
		}

		upperColDef := strings.ToUpper(colDef)
		if strings.HasPrefix(upperColDef, "PRIMARY KEY") ||
			strings.HasPrefix(upperColDef, "FOREIGN KEY") ||
			strings.HasPrefix(upperColDef, "UNIQUE") ||
			strings.HasPrefix(upperColDef, "CHECK") ||
			strings.HasPrefix(upperColDef, "INDEX") ||
			strings.HasPrefix(upperColDef, "KEY ") ||
			strings.HasPrefix(upperColDef, "CONSTRAINT") {
			continue
		}

		parts := strings.Fields(colDef)
		if len(parts) >= 2 {
			colName := strings.Trim(parts[0], "`\"';")
			colType := parts[1]
			columns = append(columns, athenastore.Column{
				Name: colName,
				Type: colType,
			})
		} else if len(parts) == 1 {
			colName := strings.Trim(parts[0], "`\"';")
			columns = append(columns, athenastore.Column{
				Name: colName,
				Type: "string",
			})
		}
	}

	return tableName, dbName, columns
}

func (s *AthenaService) parseCreateTableStatementWithLocation(query string) (string, string, []athenastore.Column, string, string) {
	tableName, dbName, columns := s.parseCreateTableStatement(query)

	var location string
	var format string

	upperQuery := strings.ToUpper(query)

	locIdx := strings.Index(upperQuery, "LOCATION")
	if locIdx != -1 {
		locPart := query[locIdx+len("LOCATION"):]
		locPart = strings.TrimSpace(locPart)

		locPart = strings.TrimPrefix(locPart, "'")
		endIdx := strings.Index(locPart, "'")
		if endIdx != -1 {
			location = locPart[:endIdx]
		}
	}

	fmtIdx := strings.Index(upperQuery, "STORED AS")
	if fmtIdx != -1 {
		fmtPart := query[fmtIdx+len("STORED AS"):]
		fmtPart = strings.TrimSpace(fmtPart)

		parts := strings.Fields(fmtPart)
		if len(parts) > 0 {
			format = strings.ToLower(strings.Trim(parts[0], "'\";"))
		}
	}

	if strings.Contains(upperQuery, "ROW FORMAT SERDE") && strings.Contains(upperQuery, "org.apache.hadoop.hive.serde2.lazy.LazySimpleSerDe") {
		if format == "" {
			format = "csv"
		}
	}

	if strings.Contains(upperQuery, "OPENCSVSERDE") || strings.Contains(upperQuery, "opencsv") {
		format = "csv"
	}

	if strings.Contains(upperQuery, "JSON") || strings.Contains(upperQuery, "HiveJsonSerDe") {
		if format == "" {
			format = "json"
		}
	}

	return tableName, dbName, columns, location, format
}

func (s *AthenaService) splitColumnDefinitions(columnsPart string) []string {
	var result []string
	var current strings.Builder
	parenDepth := 0

	for _, ch := range columnsPart {
		switch ch {
		case '(':
			parenDepth++
			current.WriteRune(ch)
		case ')':
			parenDepth--
			current.WriteRune(ch)
		case ',':
			if parenDepth == 0 {
				result = append(result, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

func (s *AthenaService) parseInsertStatement(query string) (string, []string, []interface{}) {
	upperQuery := strings.ToUpper(query)

	insertIdx := strings.Index(upperQuery, "INSERT INTO")
	if insertIdx == -1 {
		return "", nil, nil
	}

	rest := query[insertIdx+len("INSERT INTO"):]
	rest = strings.TrimSpace(rest)

	var tableName string
	var columnNames []string
	var values []interface{}

	valuesKeywordIdx := strings.Index(strings.ToUpper(rest), "VALUES")
	if valuesKeywordIdx == -1 {
		return "", nil, nil
	}

	tableAndColumnsPart := rest[:valuesKeywordIdx]
	tableAndColumnsPart = strings.TrimSpace(tableAndColumnsPart)

	openParenIdx := strings.Index(tableAndColumnsPart, "(")
	if openParenIdx == -1 {
		tablePart := strings.TrimSpace(tableAndColumnsPart)
		if strings.Contains(tablePart, ".") {
			parts := strings.Split(tablePart, ".")
			tablePart = parts[len(parts)-1]
		}
		tableName = strings.Trim(tablePart, "`\"';")
	} else {
		tablePart := strings.TrimSpace(tableAndColumnsPart[:openParenIdx])
		if strings.Contains(tablePart, ".") {
			parts := strings.Split(tablePart, ".")
			tablePart = parts[len(parts)-1]
		}
		tableName = strings.Trim(tablePart, "`\"';")

		closeParenIdx := strings.Index(tableAndColumnsPart, ")")
		if closeParenIdx == -1 {
			return tableName, nil, nil
		}
		columnsPart := tableAndColumnsPart[openParenIdx+1 : closeParenIdx]

		for _, col := range strings.Split(columnsPart, ",") {
			col = strings.TrimSpace(col)
			if col != "" {
				columnNames = append(columnNames, strings.Trim(col, "`\"';"))
			}
		}
	}

	rest = rest[valuesKeywordIdx+len("VALUES"):]
	rest = strings.TrimSpace(rest)

	rowSets := s.parseValueSets(rest)
	for _, rowSet := range rowSets {
		rowValues := s.parseValues(rowSet)
		values = append(values, rowValues)
	}

	return tableName, columnNames, values
}

func (s *AthenaService) parseValueSets(valuesPart string) []string {
	var sets []string
	var current strings.Builder
	parenDepth := 0
	inString := false
	stringDelim := byte(0)

	for i := 0; i < len(valuesPart); i++ {
		ch := valuesPart[i]

		if inString {
			if ch == '\\' && i+1 < len(valuesPart) {
				current.WriteByte(ch)
				i++
				current.WriteByte(valuesPart[i])
				continue
			}
			if ch == stringDelim {
				inString = false
			}
			current.WriteByte(ch)
			continue
		}

		if ch == '\'' || ch == '"' {
			inString = true
			stringDelim = ch
			current.WriteByte(ch)
			continue
		}

		if ch == '(' {
			if parenDepth == 0 {
				current.Reset()
			} else {
				current.WriteByte(ch)
			}
			parenDepth++
			continue
		}

		if ch == ')' {
			parenDepth--
			if parenDepth == 0 {
				sets = append(sets, current.String())
				current.Reset()
			} else {
				current.WriteByte(ch)
			}
			continue
		}

		if parenDepth > 0 {
			current.WriteByte(ch)
		}
	}

	return sets
}

func (s *AthenaService) parseValues(valuesPart string) []interface{} {
	var values []interface{}
	var current strings.Builder
	inString := false
	stringDelim := byte(0)

	for i := 0; i < len(valuesPart); i++ {
		ch := valuesPart[i]

		if inString {
			if ch == '\\' && i+1 < len(valuesPart) {
				current.WriteByte(ch)
				i++
				current.WriteByte(valuesPart[i])
				continue
			}
			if ch == stringDelim {
				inString = false
			}
			current.WriteByte(ch)
			continue
		}

		if ch == '\'' || ch == '"' {
			inString = true
			stringDelim = ch
			current.WriteByte(ch)
			continue
		}

		if ch == ',' {
			val := strings.TrimSpace(current.String())
			values = append(values, s.parseValue(val))
			current.Reset()
			continue
		}

		current.WriteByte(ch)
	}

	if current.Len() > 0 {
		val := strings.TrimSpace(current.String())
		values = append(values, s.parseValue(val))
	}

	return values
}

func (s *AthenaService) parseValue(val string) interface{} {
	val = strings.TrimSpace(val)

	if strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'") {
		return strings.Trim(val, "'")
	}
	if strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"") {
		return strings.Trim(val, "\"")
	}

	upperVal := strings.ToUpper(val)
	if upperVal == "NULL" {
		return nil
	}
	if upperVal == "TRUE" {
		return true
	}
	if upperVal == "FALSE" {
		return false
	}

	if intVal, err := strconv.ParseInt(val, 10, 64); err == nil {
		return intVal
	}
	if floatVal, err := strconv.ParseFloat(val, 64); err == nil {
		return floatVal
	}

	return val
}
