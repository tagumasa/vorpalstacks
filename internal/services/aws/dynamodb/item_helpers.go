package dynamodb

import (
	"strings"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

const maxItemSizeBytes = 400 * 1024

func getReturnConsumedCapacity(params map[string]interface{}) string {
	return request.GetStringParam(params, "ReturnConsumedCapacity")
}

func buildConsumedCapacityResponse(tableName string, capacityUnits float64) map[string]interface{} {
	return map[string]interface{}{
		"TableName":     tableName,
		"CapacityUnits": capacityUnits,
	}
}

func buildConsumedCapacityResponseWithIndex(tableName string, indexName string, capacityUnits float64, isLSI bool) map[string]interface{} {
	resp := map[string]interface{}{
		"TableName":     tableName,
		"CapacityUnits": capacityUnits,
	}
	if indexName != "" {
		resp["Table"] = map[string]interface{}{"CapacityUnits": capacityUnits}
		if isLSI {
			resp["LocalSecondaryIndexes"] = map[string]interface{}{
				indexName: map[string]interface{}{"CapacityUnits": capacityUnits},
			}
		} else {
			resp["GlobalSecondaryIndexes"] = map[string]interface{}{
				indexName: map[string]interface{}{"CapacityUnits": capacityUnits},
			}
		}
	}
	return resp
}



func (s *DynamoDBService) extractKeyFromItem(table *dbstore.Table, item map[string]*dbstore.AttributeValue) map[string]*dbstore.AttributeValue {
	key := make(map[string]*dbstore.AttributeValue)

	for _, ks := range table.KeySchema {
		attr, ok := item[ks.AttributeName]
		if !ok {
			return nil
		}
		key[ks.AttributeName] = attr
	}

	return key
}

func calculateItemSize(item map[string]*dbstore.AttributeValue) int64 {
	var size int64
	for attrName, av := range item {
		size += int64(len(attrName))
		size += calculateAttributeValueSize(av)
	}
	return size
}

func calculateAttributeValueSize(av *dbstore.AttributeValue) int64 {
	if av == nil {
		return 0
	}

	if av.S != nil {
		return int64(len(*av.S))
	}
	if av.N != nil {
		return int64(len(*av.N))
	}
	if av.B != nil {
		return int64(len(av.B))
	}
	if av.BOOL != nil {
		return 1
	}
	if av.NULL != nil {
		return 1
	}
	if av.SS != nil {
		var size int64
		for _, s := range av.SS {
			size += int64(len(s))
		}
		return size
	}
	if av.NS != nil {
		var size int64
		for _, n := range av.NS {
			size += int64(len(n))
		}
		return size
	}
	if av.BS != nil {
		var size int64
		for _, b := range av.BS {
			size += int64(len(b))
		}
		return size
	}
	if av.M != nil {
		var size int64
		for k, v := range av.M {
			size += int64(len(k))
			size += calculateAttributeValueSize(v)
		}
		return size
	}
	if av.L != nil {
		var size int64
		for _, v := range av.L {
			size += calculateAttributeValueSize(v)
		}
		return size
	}
	return 0
}

func validateKeyValueNotEmpty(key map[string]*dbstore.AttributeValue) bool {
	for _, av := range key {
		if av.S != nil && *av.S == "" {
			return false
		}
		if av.B != nil && len(av.B) == 0 {
			return false
		}
	}
	return true
}

func getBoolParamWithDefault(params map[string]interface{}, key string, defaultValue bool) bool {
	if v, ok := params[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return defaultValue
}

func tokenizeExpression(expr string) []string {
	var tokens []string
	var current strings.Builder
	i := 0

	for i < len(expr) {
		ch := expr[i]

		if ch == ' ' || ch == '\t' || ch == '\n' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			i++
			continue
		}

		if ch == '\'' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			j := i + 1
			for j < len(expr) && expr[j] != '\'' {
				j++
			}
			if j < len(expr) {
				tokens = append(tokens, expr[i:j+1])
				i = j + 1
			} else {
				tokens = append(tokens, expr[i:])
				i = j
			}
			continue
		}

		if ch == '=' || ch == '<' || ch == '>' || ch == ',' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}

			if ch == '<' && i+1 < len(expr) {
				next := expr[i+1]
				if next == '=' || next == '>' {
					tokens = append(tokens, string(ch)+string(next))
					i += 2
					continue
				}
			}

			if ch == '>' && i+1 < len(expr) && expr[i+1] == '=' {
				tokens = append(tokens, ">=")
				i += 2
				continue
			}

			tokens = append(tokens, string(ch))
			i++
			continue
		}

		if ch == '(' {
			if current.Len() > 0 && isIdentRune(current.String()) {
				current.WriteByte(ch)
				depth := 1
				i++
				for i < len(expr) && depth > 0 {
					c := expr[i]
					current.WriteByte(c)
					if c == '(' {
						depth++
					} else if c == ')' {
						depth--
					}
					i++
				}
				tokens = append(tokens, current.String())
				current.Reset()
				continue
			}
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, "(")
			i++
			continue
		}

		if ch == ')' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, ")")
			i++
			continue
		}

		current.WriteByte(ch)
		i++
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func splitAndTrim(s, sep string) []string {
	var result []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || string(s[i]) == sep {
			part := trim(s[start:i])
			if part != "" {
				result = append(result, part)
			}
			start = i + 1
		}
	}
	return result
}

func trim(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}

func isIdentRune(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}
	return true
}
