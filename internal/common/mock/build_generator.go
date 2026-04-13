// Package mock provides mock service functionality for vorpalstacks.
package mock

import (
	"strconv"
	"time"

	"vorpalstacks/internal/common/response"
	"vorpalstacks/internal/store/api"
)

// Generator generates mock responses for API operations based on shape definitions.
type Generator struct {
	shapeStore  *api.ShapeStore
	memberStore *api.MemberStore
}

// NewGenerator creates a new Generator with the given shape and member stores.
func NewGenerator(shapeStore *api.ShapeStore, memberStore *api.MemberStore) *Generator {
	return &Generator{
		shapeStore:  shapeStore,
		memberStore: memberStore,
	}
}

// GenerateResponse generates a mock response for the given output shape ID.
// Returns an empty response if the shape ID is nil or not found.
func (g *Generator) GenerateResponse(outputShapeID *int64) (map[string]interface{}, error) {
	if outputShapeID == nil {
		return response.EmptyResponse(), nil
	}

	shape, err := g.shapeStore.GetByID(*outputShapeID)
	if err != nil || shape == nil {
		return response.EmptyResponse(), nil
	}

	return g.generateStructureValue(shape.ID)
}

func (g *Generator) generateValue(shapeID int64) (interface{}, error) {
	shape, err := g.shapeStore.GetByID(shapeID)
	if err != nil || shape == nil {
		return "example-value", nil
	}

	switch shape.Type {
	case "string":
		return "example-value", nil
	case "integer":
		return 123, nil
	case "long":
		return int64(1234567890), nil
	case "boolean":
		return true, nil
	case "double", "float":
		return 123.45, nil
	case "timestamp":
		return time.Now().Format(time.RFC3339), nil
	case "blob":
		return "dGVzdC1ibG9iLWRhdGE=", nil
	case "list":
		return g.generateListValue(shape)
	case "map":
		return g.generateMapValue(shape)
	case "structure":
		return g.generateStructureValue(shape.ID)
	case "enum":
		return "ENUM_VALUE", nil
	default:
		return "example-value", nil
	}
}

func (g *Generator) generateListValue(shape *api.Shape) (interface{}, error) {
	if shape.ListMemberShapeID == nil {
		return []interface{}{}, nil
	}

	item, err := g.generateValue(*shape.ListMemberShapeID)
	if err != nil {
		return []interface{}{}, nil
	}

	return []interface{}{item}, nil
}

func (g *Generator) generateMapValue(shape *api.Shape) (interface{}, error) {
	if shape.MapKeyShapeID == nil || shape.MapValueShapeID == nil {
		return response.EmptyResponse(), nil
	}

	key, err := g.generateValue(*shape.MapKeyShapeID)
	if err != nil {
		return response.EmptyResponse(), nil
	}

	value, err := g.generateValue(*shape.MapValueShapeID)
	if err != nil {
		return response.EmptyResponse(), nil
	}

	return map[string]interface{}{
		toString(key): value,
	}, nil
}

func (g *Generator) generateStructureValue(shapeID int64) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	err := g.memberStore.ForEachByShape(shapeID, func(member *api.Member) error {
		jsonName := member.Name
		if member.JSONName != nil && *member.JSONName != "" {
			jsonName = *member.JSONName
		}

		value, err := g.generateValue(member.TargetShapeID)
		if err != nil {
			return err
		}

		if m, ok := value.(map[string]interface{}); ok && len(m) == 0 {
			result[jsonName] = ""
		} else {
			result[jsonName] = value
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case int:
		return intToString(int64(val))
	case int64:
		return intToString(val)
	case float64:
		return intToString(int64(val))
	default:
		return "key"
	}
}

func intToString(n int64) string {
	return strconv.FormatInt(n, 10)
}
