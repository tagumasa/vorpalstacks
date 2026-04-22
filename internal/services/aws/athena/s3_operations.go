package athena

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	athenastore "vorpalstacks/internal/store/aws/athena"
	s3store "vorpalstacks/internal/store/aws/s3"

	"github.com/parquet-go/parquet-go"
)

func (s *AthenaService) hasS3Support() bool {
	return s.s3ObjectStore != nil
}

func (s *AthenaService) parseS3Location(location string) (bucket, prefix string, err error) {
	if !strings.HasPrefix(location, "s3://") {
		return "", "", fmt.Errorf("invalid S3 location format: %s", location)
	}

	rest := location[5:]
	idx := strings.Index(rest, "/")
	if idx == -1 {
		return rest, "", nil
	}
	return rest[:idx], rest[idx+1:], nil
}

func (s *AthenaService) writeQueryResultsToS3(ctx context.Context, queryExecutionId string, result *athenastore.QueryResult, outputLocation string) error {
	if !s.hasS3Support() {
		return nil
	}

	bucket, prefix, err := s.parseS3Location(outputLocation)
	if err != nil {
		return err
	}

	csvData := s.resultSetToCSV(result.ResultSet)

	key := prefix
	if !strings.HasSuffix(key, "/") {
		key += "/"
	}
	key += queryExecutionId + ".csv"

	_, err = s.s3ObjectStore.Put(ctx, bucket, key, bytes.NewReader(csvData), "text/csv", nil)
	return err
}

func (s *AthenaService) resultSetToCSV(resultSet *athenastore.ResultSet) []byte {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	for _, row := range resultSet.Rows {
		var record []string
		for _, datum := range row.Data {
			record = append(record, datum.VarCharValue)
		}
		if err := writer.Write(record); err != nil {
			logs.Error("Failed to write CSV record", logs.Err(err))
		}
	}
	writer.Flush()

	return buf.Bytes()
}

func (s *AthenaService) readDataFromS3(ctx context.Context, location string, format string) ([]map[string]interface{}, error) {
	if !s.hasS3Support() {
		return nil, fmt.Errorf("S3 integration not available")
	}

	bucket, prefix, err := s.parseS3Location(location)
	if err != nil {
		return nil, err
	}

	listResult, err := s.s3ObjectStore.List(bucket, prefix, "", "", 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to list S3 objects: %w", err)
	}

	var allData []map[string]interface{}

	for _, obj := range listResult.Objects {
		reader, _, err := s.s3ObjectStore.Get(ctx, bucket, obj.Key)
		if err != nil {
			continue
		}

		data, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			continue
		}

		var parsedData []map[string]interface{}
		detectedFormat := format
		if detectedFormat == "" {
			detectedFormat = s.detectFileFormat(obj.Key, data)
		}

		switch detectedFormat {
		case "parquet":
			parsedData, err = s.parseParquetData(data)
		case "csv":
			parsedData, err = s.parseCSVData(data)
		case "json":
			parsedData, err = s.parseJSONData(data)
		default:
			parsedData, err = s.parseCSVData(data)
		}

		if err != nil {
			continue
		}

		allData = append(allData, parsedData...)
	}

	return allData, nil
}

func (s *AthenaService) detectFileFormat(key string, data []byte) string {
	lowerKey := strings.ToLower(key)
	if strings.HasSuffix(lowerKey, ".parquet") || strings.HasSuffix(lowerKey, ".par") {
		return "parquet"
	}
	if strings.HasSuffix(lowerKey, ".json") {
		return "json"
	}
	if strings.HasSuffix(lowerKey, ".csv") {
		return "csv"
	}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) > 4 {
		magic := string(trimmed[:4])
		if magic == "PAR1" {
			return "parquet"
		}
	}
	if len(trimmed) > 0 {
		if trimmed[0] == '{' || trimmed[0] == '[' {
			return "json"
		}
	}
	return "csv"
}

func (s *AthenaService) parseCSVData(data []byte) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 1 {
		return results, nil
	}

	headers := records[0]
	for i := 1; i < len(records); i++ {
		row := make(map[string]interface{})
		for j, val := range records[i] {
			if j < len(headers) {
				row[headers[j]] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

func (s *AthenaService) parseJSONData(data []byte) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return results, nil
	}

	if trimmed[0] == '[' {
		if err := json.Unmarshal(data, &results); err != nil {
			return nil, err
		}
	} else if trimmed[0] == '{' {
		var singleObj map[string]interface{}
		if err := json.Unmarshal(data, &singleObj); err != nil {
			return nil, err
		}
		results = append(results, singleObj)
	}

	return results, nil
}

func (s *AthenaService) parseParquetData(data []byte) ([]map[string]interface{}, error) {
	reader := parquet.NewReader(bytes.NewReader(data))
	defer reader.Close()

	var results []map[string]interface{}
	for {
		row := make(map[string]interface{})
		err := reader.Read(&row)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read parquet row: %w", err)
		}
		results = append(results, row)
	}

	return results, nil
}

func (s *AthenaService) loadExternalTableData(reqCtx *request.RequestContext, catalog, database, tableName string) ([]map[string]interface{}, error) {
	stores, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	table, err := stores.tableStore.GetTable(catalog, database, tableName)
	if err != nil {
		return nil, err
	}

	location := ""
	if table.Parameters != nil {
		location = table.Parameters["LOCATION"]
	}
	if location == "" {
		return nil, fmt.Errorf("table %s has no LOCATION parameter", tableName)
	}

	format := ""
	if table.Parameters != nil {
		format = table.Parameters["format"]
	}

	return s.readDataFromS3(reqCtx, location, format)
}

func (s *AthenaService) convertS3DataToStoredTable(data []map[string]interface{}, columns []athenastore.Column) *athenastore.StoredTable {
	storedTable := &athenastore.StoredTable{
		Columns: columns,
		Rows:    []*athenastore.StoredRow{},
	}

	for _, rowData := range data {
		row := &athenastore.StoredRow{Values: make(map[string]interface{})}
		for _, col := range columns {
			if val, ok := rowData[col.Name]; ok {
				row.Values[col.Name] = val
			} else {
				row.Values[col.Name] = nil
			}
		}
		storedTable.Rows = append(storedTable.Rows, row)
	}

	return storedTable
}

var _ S3ObjectStore = (*s3store.ObjectStore)(nil)
