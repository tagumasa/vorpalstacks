package helpers

import (
	"testing"
)

type TestStruct struct {
	Name     string
	Age      int
	Score    float64
	Active   bool
	Password string
}

func TestApplyUpdate(t *testing.T) {
	tests := []struct {
		name       string
		item       interface{}
		updateExpr string
		values     map[string]interface{}
		wantErr    bool
		checkField func(*TestStruct) bool
	}{
		{
			name:       "SET string field",
			item:       &TestStruct{Name: "old"},
			updateExpr: "SET #Name = :name",
			values:     map[string]interface{}{":name": "new"},
			wantErr:    false,
			checkField: func(s *TestStruct) bool { return s.Name == "new" },
		},
		{
			name:       "SET int field",
			item:       &TestStruct{Age: 20},
			updateExpr: "SET #Age = :age",
			values:     map[string]interface{}{":age": 30},
			wantErr:    false,
			checkField: func(s *TestStruct) bool { return s.Age == 30 },
		},
		{
			name:       "SET float field",
			item:       &TestStruct{Score: 1.5},
			updateExpr: "SET #Score = :score",
			values:     map[string]interface{}{":score": 2.5},
			wantErr:    false,
			checkField: func(s *TestStruct) bool { return s.Score == 2.5 },
		},
		{
			name:       "nil item",
			item:       nil,
			updateExpr: "SET #Name = :name",
			values:     map[string]interface{}{":name": "new"},
			wantErr:    true,
		},
		{
			name:       "non-pointer item",
			item:       TestStruct{Name: "test"},
			updateExpr: "SET #Name = :name",
			values:     map[string]interface{}{":name": "new"},
			wantErr:    true,
		},
		{
			name:       "invalid update expression",
			item:       &TestStruct{},
			updateExpr: "INVALID #Name = :name",
			values:     map[string]interface{}{":name": "new"},
			wantErr:    true,
		},
		{
			name:       "missing attribute prefix",
			item:       &TestStruct{},
			updateExpr: "SET Name = :name",
			values:     map[string]interface{}{":name": "new"},
			wantErr:    true,
		},
		{
			name:       "missing field in struct",
			item:       &TestStruct{},
			updateExpr: "SET #Nonexistent = :value",
			values:     map[string]interface{}{":value": "test"},
			wantErr:    true,
		},
		{
			name:       "value not provided",
			item:       &TestStruct{},
			updateExpr: "SET #Name = :name",
			values:     map[string]interface{}{},
			wantErr:    true,
		},
		{
			name:       "REMOVE operation",
			item:       &TestStruct{Name: "test"},
			updateExpr: "REMOVE #Name",
			values:     map[string]interface{}{},
			wantErr:    false,
			checkField: func(s *TestStruct) bool { return s.Name == "" },
		},
		{
			name:       "type conversion int to int64",
			item:       &TestStruct{Age: 10},
			updateExpr: "SET #Age = :age",
			values:     map[string]interface{}{":age": int64(20)},
			wantErr:    false,
			checkField: func(s *TestStruct) bool { return s.Age == 20 },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ApplyUpdate(tt.item, tt.updateExpr, tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ApplyUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.checkField != nil {
				if !tt.checkField(tt.item.(*TestStruct)) {
					t.Errorf("field check failed")
				}
			}
		})
	}
}

func TestParseUpdateExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantErr    bool
		wantLen    int
	}{
		{
			name:       "single SET expression",
			expression: "SET #Name = :name",
			wantErr:    false,
			wantLen:    1,
		},
		{
			name:       "multiple SET expressions",
			expression: "SET #Name = :name, SET #Age = :age",
			wantErr:    false,
			wantLen:    2,
		},
		{
			name:       "REMOVE expression",
			expression: "REMOVE #Name",
			wantErr:    false,
			wantLen:    1,
		},
		{
			name:       "ADD expression",
			expression: "ADD #Count :count",
			wantErr:    false,
			wantLen:    1,
		},
		{
			name:       "DELETE expression",
			expression: "DELETE #Tags :tag",
			wantErr:    false,
			wantLen:    1,
		},
		{
			name:       "invalid expression",
			expression: "INVALID #Name = :name",
			wantErr:    true,
			wantLen:    0,
		},
		{
			name:       "empty expression",
			expression: "",
			wantErr:    false,
			wantLen:    0,
		},
		{
			name:       "whitespace only",
			expression: "   ",
			wantErr:    false,
			wantLen:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseUpdateExpression(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUpdateExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(result) != tt.wantLen {
				t.Errorf("ParseUpdateExpression() len = %v, want %v", len(result), tt.wantLen)
			}
		})
	}
}

func TestValidateUpdateExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		wantErr    bool
	}{
		{
			name:       "valid SET expression",
			expression: "SET #Name = :name",
			wantErr:    false,
		},
		{
			name:       "empty expression",
			expression: "",
			wantErr:    true,
		},
		{
			name:       "invalid expression",
			expression: "INVALID #Name",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateExpression(tt.expression)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUpdateExpression() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
