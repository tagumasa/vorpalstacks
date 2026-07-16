package s3

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"vorpalstacks/pkg/sqlparser"
)

type testOutputWriter struct {
	records [][]byte
}

func (w *testOutputWriter) WriteRecords(data []byte) error {
	w.records = append(w.records, data)
	return nil
}

func (w *testOutputWriter) AddBytesScanned(n int64)   {}
func (w *testOutputWriter) AddBytesProcessed(n int64) {}
func (w *testOutputWriter) MaybeWriteProgress() error { return nil }

func (w *testOutputWriter) WriteStats() error {
	return nil
}

func (w *testOutputWriter) WriteEnd() error {
	return nil
}

func (w *testOutputWriter) Flush() error {
	return nil
}

func (w *testOutputWriter) String() string {
	return string(bytes.Join(w.records, nil))
}

func TestSelectEngineCSV(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		csvData   string
		wantMatch []string
		dontWant  []string
	}{
		{
			name:      "select all",
			sql:       "SELECT * FROM s3object",
			csvData:   "a,b,c\n1,2,3\n4,5,6",
			wantMatch: []string{"1,2,3", "4,5,6"},
		},
		{
			name:      "select with where",
			sql:       "SELECT _1 FROM s3object WHERE _2 = '5'",
			csvData:   "a,b,c\n1,2,3\n4,5,6",
			wantMatch: []string{"4"},
			dontWant:  []string{"1"},
		},
		{
			name:      "select with comparison",
			sql:       "SELECT _1, _2 FROM s3object WHERE _3 > '5'",
			csvData:   "col1,col2,col3\n1,2,3\n7,8,9",
			wantMatch: []string{"7,8"},
		},
		{
			name:      "select with LIKE",
			sql:       "SELECT * FROM s3object WHERE _1 LIKE '%test%'",
			csvData:   "value\ntest123\nnotmatch\ntesting",
			wantMatch: []string{"test123", "testing"},
			dontWant:  []string{"notmatch"},
		},
		{
			name:      "select with IN",
			sql:       "SELECT _1 FROM s3object WHERE _1 IN ('a', 'c')",
			csvData:   "col\na\nb\nc\nd",
			wantMatch: []string{"a", "c"},
			dontWant:  []string{"b", "d"},
		},
		{
			name:      "select with AND",
			sql:       "SELECT * FROM s3object WHERE _1 = '1' AND _2 = '2'",
			csvData:   "a,b\n1,2\n1,3\n2,2",
			wantMatch: []string{"1,2"},
			dontWant:  []string{"1,3", "2,2"},
		},
		{
			name:      "select with OR",
			sql:       "SELECT _1 FROM s3object WHERE _1 = '1' OR _1 = '3'",
			csvData:   "col\n1\n2\n3\n4",
			wantMatch: []string{"1", "3"},
			dontWant:  []string{"2", "4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &SelectObjectContentInput{
				Expression:     tt.sql,
				ExpressionType: ExpressionTypeSQL,
				InputSerialization: &InputSerialization{
					CSV: &CSVInput{
						FileHeaderInfo: FileHeaderInfoNone,
					},
				},
				OutputSerialization: &OutputSerialization{
					CSV: &CSVOutput{},
				},
			}

			engine, err := NewSelectEngine(input)
			if err != nil {
				t.Fatalf("NewSelectEngine error: %v", err)
			}

			writer := &testOutputWriter{}

			err = engine.Execute(context.Background(), strings.NewReader(tt.csvData), writer)
			if err != nil {
				t.Fatalf("Execute error: %v", err)
			}

			result := writer.String()

			for _, want := range tt.wantMatch {
				if !strings.Contains(result, want) {
					t.Errorf("expected result to contain %q, got %q", want, result)
				}
			}

			for _, dontWant := range tt.dontWant {
				if strings.Contains(result, dontWant) {
					t.Errorf("expected result NOT to contain %q, got %q", dontWant, result)
				}
			}
		})
	}
}

func TestSelectEngineJSON(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		jsonData  string
		wantMatch []string
		dontWant  []string
	}{
		{
			name:      "select all",
			sql:       "SELECT * FROM s3object",
			jsonData:  `{"name":"test1","value":1}` + "\n" + `{"name":"test2","value":2}`,
			wantMatch: []string{"test1", "test2"},
		},
		{
			name:      "select with where",
			sql:       "SELECT name FROM s3object WHERE value = '2'",
			jsonData:  `{"name":"test1","value":1}` + "\n" + `{"name":"test2","value":2}`,
			wantMatch: []string{"test2"},
			dontWant:  []string{"test1"},
		},
		{
			name:      "select nested field",
			sql:       "SELECT name FROM s3object WHERE address.city = 'Tokyo'",
			jsonData:  `{"name":"user1","address":{"city":"Tokyo"}}` + "\n" + `{"name":"user2","address":{"city":"Osaka"}}`,
			wantMatch: []string{"user1"},
			dontWant:  []string{"user2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &SelectObjectContentInput{
				Expression:     tt.sql,
				ExpressionType: ExpressionTypeSQL,
				InputSerialization: &InputSerialization{
					JSON: &JSONInput{
						Type: "LINES",
					},
				},
				OutputSerialization: &OutputSerialization{
					JSON: &JSONOutput{},
				},
			}

			engine, err := NewSelectEngine(input)
			if err != nil {
				t.Fatalf("NewSelectEngine error: %v", err)
			}

			writer := &testOutputWriter{}

			err = engine.Execute(context.Background(), strings.NewReader(tt.jsonData), writer)
			if err != nil {
				t.Fatalf("Execute error: %v", err)
			}

			result := writer.String()

			for _, want := range tt.wantMatch {
				if !strings.Contains(result, want) {
					t.Errorf("expected result to contain %q, got %q", want, result)
				}
			}

			for _, dontWant := range tt.dontWant {
				if strings.Contains(result, dontWant) {
					t.Errorf("expected result NOT to contain %q, got %q", dontWant, result)
				}
			}
		})
	}
}

func TestSQLParsing(t *testing.T) {
	tests := []struct {
		sql  string
		desc string
	}{
		{"SELECT * FROM s3object", "select all"},
		{"SELECT _1, _2 FROM s3object", "select columns"},
		{"SELECT * FROM s3object WHERE _1 = 'value'", "where clause"},
		{"SELECT * FROM s3object WHERE _1 LIKE '%test%'", "like clause"},
		{"SELECT * FROM s3object WHERE _1 IN ('a', 'b')", "in clause"},
		{"SELECT * FROM s3object WHERE _1 IS NULL", "is null"},
		{"SELECT * FROM s3object s WHERE s._1 = 'value'", "alias"},
		{"SELECT _1 AS col1 FROM s3object", "column alias"},
		{"SELECT * FROM s3object WHERE _1 = 'a' AND _2 = 'b'", "and"},
		{"SELECT * FROM s3object WHERE _1 = 'a' OR _2 = 'b'", "or"},
	}

	opts := sqlparser.ParserOptions{Dialect: sqlparser.DialectS3Select}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := sqlparser.ParseWithOptions(tt.sql, opts)
			if err != nil {
				t.Errorf("Parse error for %q: %v", tt.sql, err)
			}
		})
	}
}

func TestEventStreamEncoding(t *testing.T) {
	var buf bytes.Buffer
	writer := NewSelectEventStreamWriter(&buf, nil)

	if err := writer.WriteRecords([]byte("test data")); err != nil {
		t.Fatalf("WriteRecords error: %v", err)
	}
	if err := writer.WriteStats(); err != nil {
		t.Fatalf("WriteStats error: %v", err)
	}
	if err := writer.WriteEnd(); err != nil {
		t.Fatalf("WriteEnd error: %v", err)
	}

	data := buf.Bytes()
	if len(data) == 0 {
		t.Error("expected non-empty event stream output")
	}
}
