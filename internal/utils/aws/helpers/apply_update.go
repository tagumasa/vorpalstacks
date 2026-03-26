// Package helpers provides AWS-specific utility functions for vorpalstacks.
package helpers

import (
	"fmt"
	"reflect"
	"strings"
)

// ApplyUpdate applies an update expression to an item.
// This implementation parses and applies update expression based on AWS DynamoDB update expression rules.
//
// Parameters:
//   - item: The item to update (must be a pointer to a struct)
//   - updateExpression: The update expression to apply (e.g., "SET #a = :val, #b = :val2")
//   - values: The values to bind to the expression (e.g., map[string]interface{}{":val": "value1", ":val2": "value2"})
//
// Returns:
//   - error: An error if update fails, nil otherwise
func ApplyUpdate(item interface{}, updateExpression string, values map[string]interface{}) error {
	if item == nil {
		return fmt.Errorf("item cannot be nil")
	}

	v := reflect.ValueOf(item)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("item must be a pointer to a struct")
	}

	v = v.Elem()
	if !v.IsValid() {
		return fmt.Errorf("item pointer is invalid")
	}
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("item must point to a struct")
	}

	parts := strings.Split(updateExpression, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		action := ""
		attributePath := ""
		if strings.HasPrefix(part, "SET ") {
			action = "SET"
			attributePath = strings.TrimSpace(strings.TrimPrefix(part, "SET "))
		} else if strings.HasPrefix(part, "REMOVE ") {
			action = "REMOVE"
			attributePath = strings.TrimSpace(strings.TrimPrefix(part, "REMOVE "))
		} else if strings.HasPrefix(part, "ADD ") {
			action = "ADD"
			attributePath = strings.TrimSpace(strings.TrimPrefix(part, "ADD "))
		} else if strings.HasPrefix(part, "DELETE ") {
			action = "DELETE"
			attributePath = strings.TrimSpace(strings.TrimPrefix(part, "DELETE "))
		} else {
			return fmt.Errorf("invalid update expression: %s", part)
		}

		if !strings.HasPrefix(attributePath, "#") {
			return fmt.Errorf("attribute path must start with #: %s", attributePath)
		}
		attrName := strings.TrimPrefix(strings.SplitN(attributePath, " ", 2)[0], "#")

		field := v.FieldByName(attrName)
		if !field.IsValid() {
			return fmt.Errorf("field %s not found in struct", attrName)
		}

		switch action {
		case "SET":
			setParts := strings.Split(attributePath, " = ")
			if len(setParts) != 2 {
				return fmt.Errorf("invalid SET expression: %s", part)
			}
			valueKey := setParts[1]
			if !strings.HasPrefix(valueKey, ":") {
				return fmt.Errorf("value must start with : %s", valueKey)
			}
			value, ok := values[valueKey]
			if !ok {
				return fmt.Errorf("value %s not provided", valueKey)
			}
			if !field.CanSet() {
				return fmt.Errorf("field %s cannot be set", attrName)
			}
			if err := setField(field, value); err != nil {
				return fmt.Errorf("failed to set field %s: %w", attrName, err)
			}

		case "REMOVE":
			if !field.CanSet() {
				return fmt.Errorf("field %s cannot be removed", attrName)
			}
			field.Set(reflect.Zero(field.Type()))

		case "ADD":
			addParts := strings.Split(attributePath, ":")
			if len(addParts) != 2 {
				return fmt.Errorf("invalid ADD expression: %s", part)
			}
			valueKey := ":" + addParts[1]
			value, ok := values[valueKey]
			if !ok {
				return fmt.Errorf("value %s not provided", valueKey)
			}
			if value == nil {
				return fmt.Errorf("value %s cannot be nil", valueKey)
			}
			if !field.CanSet() {
				return fmt.Errorf("field %s cannot be updated", attrName)
			}
			existingValue := field.Interface()
			if existingValue == nil {
				return fmt.Errorf("existing value for %s is nil", attrName)
			}
			switch ev := existingValue.(type) {
			case int:
				intVal, ok := value.(int)
				if !ok {
					return fmt.Errorf("type mismatch for %s: expected int, got %T", attrName, value)
				}
				field.Set(reflect.ValueOf(ev + intVal))
			case int64:
				intVal, ok := value.(int64)
				if !ok {
					return fmt.Errorf("type mismatch for %s: expected int64, got %T", attrName, value)
				}
				field.Set(reflect.ValueOf(ev + intVal))
			case float64:
				floatVal, ok := value.(float64)
				if !ok {
					return fmt.Errorf("type mismatch for %s: expected float64, got %T", attrName, value)
				}
				field.Set(reflect.ValueOf(ev + floatVal))
			case string:
				strVal, ok := value.(string)
				if !ok {
					return fmt.Errorf("type mismatch for %s: expected string, got %T", attrName, value)
				}
				field.Set(reflect.ValueOf(ev + strVal))
			default:
				return fmt.Errorf("unsupported type for ADD operation: %T", existingValue)
			}

		case "DELETE":
			if !field.CanSet() {
				return fmt.Errorf("field %s cannot be deleted", attrName)
			}
			field.Set(reflect.Zero(field.Type()))

		default:
			return fmt.Errorf("unsupported action: %s", action)
		}
	}

	return nil
}

func setField(field reflect.Value, value interface{}) error {
	fieldType := field.Type()
	valueType := reflect.TypeOf(value)

	if fieldType == valueType {
		field.Set(reflect.ValueOf(value))
		return nil
	}

	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		rv := reflect.ValueOf(value)
		switch rv.Type().Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetInt(rv.Int())
			return nil
		case reflect.Float32, reflect.Float64:
			field.SetInt(int64(rv.Float()))
			return nil
		}
	case reflect.Float32, reflect.Float64:
		rv := reflect.ValueOf(value)
		switch rv.Type().Kind() {
		case reflect.Float32, reflect.Float64:
			field.SetFloat(rv.Float())
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			field.SetFloat(float64(rv.Int()))
			return nil
		}
	case reflect.String:
		rv := reflect.ValueOf(value)
		if rv.Type().Kind() == reflect.String {
			field.SetString(rv.String())
			return nil
		}
	}

	return fmt.Errorf("cannot convert %T to %s", value, fieldType)
}

// ParseUpdateExpression parses an update expression into its components.
//
// Parameters:
//   - expression: The update expression to parse
//
// Returns:
//   - []string: The parsed components
func ParseUpdateExpression(expression string) ([]string, error) {
	// Split by comma to get individual update actions
	parts := strings.Split(expression, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Validate action
		if !strings.HasPrefix(part, "SET ") && !strings.HasPrefix(part, "REMOVE ") &&
			!strings.HasPrefix(part, "ADD ") && !strings.HasPrefix(part, "DELETE ") {
			return nil, fmt.Errorf("invalid update expression: %s", part)
		}

		result = append(result, part)
	}

	return result, nil
}

// ValidateUpdateExpression validates an update expression.
//
// Parameters:
//   - expression: The update expression to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func ValidateUpdateExpression(expression string) error {
	if expression == "" {
		return fmt.Errorf("update expression cannot be empty")
	}

	// Parse and validate
	_, err := ParseUpdateExpression(expression)
	return err
}
