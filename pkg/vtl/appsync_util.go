package vtl

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	appSyncUtilIsNullRegex              = regexp.MustCompile(`\$util\.isNull\(([^)]+)\)`)
	appSyncUtilIsNullOrEmptyRegex       = regexp.MustCompile(`\$util\.isNullOrEmpty\(([^)]+)\)`)
	appSyncUtilDefaultIfNullRegex       = regexp.MustCompile(`\$util\.defaultIfNull\(([^,]+),\s*([^)]+)\)`)
	appSyncUtilAutoIdRegex              = regexp.MustCompile(`\$util\.autoId\(\)`)
	appSyncUtilErrorRegex               = regexp.MustCompile(`\$util\.error\(([^)]*)\)`)
	appSyncUtilAppendErrorRegex         = regexp.MustCompile(`\$util\.appendError\(([^)]*)\)`)
	appSyncUtilUnauthorizedRegex        = regexp.MustCompile(`\$util\.unauthorized\(\)`)
	appSyncUtilToJsonRegex              = regexp.MustCompile(`\$util\.toJson\(([^)]+)\)`)
	appSyncUtilQuietIfRegex             = regexp.MustCompile(`\$util\.quietIf\(([^)]+)\)`)
	appSyncUtilQrRegex                  = regexp.MustCompile(`\$util\.qr\(([^)]+)\)`)
	appSyncUtilTimeNowISO8601Regex      = regexp.MustCompile(`\$util\.time\.nowISO8601\(\)`)
	appSyncUtilTimeNowEpochSecondsRegex = regexp.MustCompile(`\$util\.time\.nowEpochSeconds\(\)`)
	appSyncUtilTimeNowEpochMilliRegex   = regexp.MustCompile(`\$util\.time\.nowEpochMilliSeconds\(\)`)
	appSyncUtilTimeNowFormattedRegex    = regexp.MustCompile(`\$util\.time\.nowFormatted\(([^)]+)\)`)
	appSyncUtilLogInfoRegex             = regexp.MustCompile(`\$util\.log\.info\(([^)]+)\)`)
	appSyncUtilLogErrorRegex            = regexp.MustCompile(`\$util\.log\.error\(([^)]+)\)`)
	appSyncUtilLogWarnRegex             = regexp.MustCompile(`\$util\.log\.warn\(([^)]+)\)`)
	appSyncUtilDynamoDBToMapValuesRegex = regexp.MustCompile(`\$util\.(?:dynamodb\.)?toMapValues\(([^)]+)\)`)
	appSyncUtilValidateRegex            = regexp.MustCompile(`\$util\.validate\(([^,]+),\s*([^,)]+)(?:,\s*([^)]+))?\)`)
	appSyncUtilIsStringRegex            = regexp.MustCompile(`\$util\.isString\(([^)]+)\)`)
	appSyncUtilIsListRegex              = regexp.MustCompile(`\$util\.isList\(([^)]+)\)`)
	appSyncUtilMathEqualsRegex          = regexp.MustCompile(`\$util\.math\.equals\(([^,]+),\s*([^)]+)\)`)
)

func (e *Engine) processAppSyncUtil(templateStr string) string {
	if e.AppSyncCtx == nil {
		return templateStr
	}

	result := templateStr

	result = e.processUtilQuietIf(result)
	result = e.processUtilQr(result)

	result = e.processUtilAutoId(result)
	result = e.processUtilTimeNowISO8601(result)
	result = e.processUtilTimeNowEpochSeconds(result)
	result = e.processUtilTimeNowEpochMilliSeconds(result)
	result = e.processUtilTimeNowFormatted(result)

	result = e.processUtilIsNullOrEmpty(result)
	result = e.processUtilIsNull(result)
	result = e.processUtilDefaultIfNull(result)
	result = e.processUtilIsString(result)
	result = e.processUtilIsList(result)
	result = e.processUtilMathEquals(result)
	result = e.processUtilToJson(result)
	result = e.processUtilValidate(result)
	result = e.processUtilDynamoDBToMapValues(result)

	result = e.processUtilLogInfo(result)
	result = e.processUtilLogError(result)
	result = e.processUtilLogWarn(result)

	result = e.processUtilError(result)
	result = e.processUtilAppendError(result)
	result = e.processUtilUnauthorized(result)

	return result
}

func (e *Engine) processUtilIsNull(templateStr string) string {
	return appSyncUtilIsNullRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		arg := strings.TrimSpace(appSyncUtilIsNullRegex.FindStringSubmatch(match)[1])
		if isNullLiteral(arg) {
			return "true"
		}
		if arg == "" {
			return "false"
		}
		return "false"
	})
}

func (e *Engine) processUtilIsNullOrEmpty(templateStr string) string {
	return appSyncUtilIsNullOrEmptyRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		arg := strings.TrimSpace(appSyncUtilIsNullOrEmptyRegex.FindStringSubmatch(match)[1])
		if isNullLiteral(arg) {
			return "true"
		}
		if arg == `""` || arg == `''` {
			return "true"
		}
		if arg == "[]" || arg == "{}" {
			return "true"
		}
		return "false"
	})
}

func (e *Engine) processUtilDefaultIfNull(templateStr string) string {
	return appSyncUtilDefaultIfNullRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		val := strings.TrimSpace(appSyncUtilDefaultIfNullRegex.FindStringSubmatch(match)[1])
		def := strings.TrimSpace(appSyncUtilDefaultIfNullRegex.FindStringSubmatch(match)[2])
		if isNullLiteral(val) {
			return e.formatValue(e.resolveValue(def))
		}
		return e.formatValue(e.resolveValue(val))
	})
}

func (e *Engine) processUtilAutoId(templateStr string) string {
	return appSyncUtilAutoIdRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return uuid.New().String()
	})
}

func (e *Engine) processUtilError(templateStr string) string {
	return appSyncUtilErrorRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		submatches := appSyncUtilErrorRegex.FindStringSubmatch(match)
		args := strings.TrimSpace(submatches[1])
		parts := splitErrorArgs(args)
		if e.AppSyncCtx == nil {
			return ""
		}
		msg := strings.TrimSpace(strings.Trim(parts[0], `"'`))
		if msg == "" {
			msg = "Unknown error"
		}
		errType := ""
		if len(parts) > 1 {
			trimmed := strings.TrimSpace(parts[1])
			errType = strings.Trim(trimmed, `"'`)
		}
		data := parseErrorData(parts)
		e.AppSyncCtx.Errors = append(e.AppSyncCtx.Errors, AppSyncError{
			Message:   msg,
			ErrorType: errType,
			Data:      data,
		})
		return ""
	})
}

func (e *Engine) processUtilAppendError(templateStr string) string {
	return appSyncUtilAppendErrorRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		submatches := appSyncUtilAppendErrorRegex.FindStringSubmatch(match)
		args := strings.TrimSpace(submatches[1])
		if e.AppSyncCtx == nil {
			return ""
		}
		parts := splitErrorArgs(args)
		msg := strings.TrimSpace(strings.Trim(parts[0], `"'`))
		errType := ""
		if len(parts) > 1 {
			trimmed := strings.TrimSpace(parts[1])
			errType = strings.Trim(trimmed, `"'`)
		}
		data := parseErrorData(parts)
		e.AppSyncCtx.Errors = append(e.AppSyncCtx.Errors, AppSyncError{
			Message:   msg,
			ErrorType: errType,
			Data:      data,
		})
		return ""
	})
}

func parseErrorData(parts []string) interface{} {
	if len(parts) > 2 {
		dataStr := strings.TrimSpace(parts[2])
		var parsed interface{}
		if json.Unmarshal([]byte(dataStr), &parsed) == nil {
			return parsed
		}
		return strings.Trim(dataStr, `"'`)
	}
	return nil
}

func (e *Engine) processUtilUnauthorized(templateStr string) string {
	return appSyncUtilUnauthorizedRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		if e.AppSyncCtx == nil {
			return match
		}
		e.AppSyncCtx.Errors = append(e.AppSyncCtx.Errors, AppSyncError{
			Message:   "Unauthorized",
			ErrorType: "Unauthorized",
		})
		return match
	})
}

func (e *Engine) processUtilToJson(templateStr string) string {
	return appSyncUtilToJsonRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		arg := strings.TrimSpace(appSyncUtilToJsonRegex.FindStringSubmatch(match)[1])
		if isNullLiteral(arg) {
			return "null"
		}
		if strings.HasPrefix(arg, `"{`) || strings.HasPrefix(arg, `"[`) {
			return arg
		}
		if strings.HasPrefix(arg, "{") || strings.HasPrefix(arg, "[") {
			return arg
		}
		if (strings.HasPrefix(arg, `"`) && strings.HasSuffix(arg, `"`)) ||
			(strings.HasPrefix(arg, `'`) && strings.HasSuffix(arg, `'`)) {
			inner := arg[1 : len(arg)-1]
			return inner
		}
		return arg
	})
}

func (e *Engine) processUtilQuietIf(templateStr string) string {
	return appSyncUtilQuietIfRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return ""
	})
}

func (e *Engine) processUtilQr(templateStr string) string {
	return appSyncUtilQrRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		inner := strings.TrimSpace(appSyncUtilQrRegex.FindStringSubmatch(match)[1])
		if strings.HasPrefix(inner, "$") && strings.Contains(inner, "=") {
			parts := strings.SplitN(inner, "=", 2)
			varName := strings.TrimSpace(parts[0])
			if strings.HasPrefix(varName, "$") {
				varName = varName[1:]
			}
			valueExpr := strings.TrimSpace(parts[1])
			value := e.resolveValue(valueExpr)
			if e.context != nil && e.context.Context != nil {
				e.context.Context[varName] = value
			}
		}
		return ""
	})
}

func (e *Engine) processUtilTimeNowISO8601(templateStr string) string {
	return appSyncUtilTimeNowISO8601Regex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	})
}

func (e *Engine) processUtilTimeNowEpochSeconds(templateStr string) string {
	return appSyncUtilTimeNowEpochSecondsRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return fmt.Sprintf("%d", time.Now().Unix())
	})
}

func (e *Engine) processUtilTimeNowEpochMilliSeconds(templateStr string) string {
	return appSyncUtilTimeNowEpochMilliRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return fmt.Sprintf("%d", time.Now().UnixMilli())
	})
}

func (e *Engine) processUtilTimeNowFormatted(templateStr string) string {
	return appSyncUtilTimeNowFormattedRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		format := strings.TrimSpace(appSyncUtilTimeNowFormattedRegex.FindStringSubmatch(match)[1])
		format = strings.Trim(format, `"'`)
		if format == "" {
			format = "yyyy-MM-dd'T'HH:mm:ss.SSS'Z'"
		}
		goFormat := javaToGoTimeFormat(format)
		return time.Now().UTC().Format(goFormat)
	})
}

func (e *Engine) processUtilLogInfo(templateStr string) string {
	return appSyncUtilLogInfoRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return ""
	})
}

func (e *Engine) processUtilLogError(templateStr string) string {
	return appSyncUtilLogErrorRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return ""
	})
}

func (e *Engine) processUtilLogWarn(templateStr string) string {
	return appSyncUtilLogWarnRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		return ""
	})
}

func (e *Engine) processUtilDynamoDBToMapValues(templateStr string) string {
	return appSyncUtilDynamoDBToMapValuesRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		arg := strings.TrimSpace(appSyncUtilDynamoDBToMapValuesRegex.FindStringSubmatch(match)[1])
		if isNullLiteral(arg) {
			return "null"
		}
		return arg
	})
}

func (e *Engine) processUtilValidate(templateStr string) string {
	return appSyncUtilValidateRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		cond := strings.TrimSpace(appSyncUtilValidateRegex.FindStringSubmatch(match)[1])
		if e.evaluateCondition(cond) {
			return ""
		}
		submatches := appSyncUtilValidateRegex.FindStringSubmatch(match)
		msg := "Validation error"
		errType := "Validation"
		if len(submatches) > 2 && submatches[2] != "" {
			msg = strings.TrimSpace(strings.Trim(submatches[2], `"'`))
		}
		if len(submatches) > 3 && submatches[3] != "" {
			errType = strings.TrimSpace(strings.Trim(submatches[3], `"'`))
		}
		if e.AppSyncCtx != nil {
			e.AppSyncCtx.Errors = append(e.AppSyncCtx.Errors, AppSyncError{
				Message:   msg,
				ErrorType: errType,
			})
		}
		return match
	})
}

func (e *Engine) processUtilIsString(templateStr string) string {
	return appSyncUtilIsStringRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		arg := strings.TrimSpace(appSyncUtilIsStringRegex.FindStringSubmatch(match)[1])
		if isNullLiteral(arg) {
			return "false"
		}
		if (strings.HasPrefix(arg, `"`) && strings.HasSuffix(arg, `"`)) ||
			(strings.HasPrefix(arg, `'`) && strings.HasSuffix(arg, `'`)) {
			return "true"
		}
		return "false"
	})
}

func (e *Engine) processUtilIsList(templateStr string) string {
	return appSyncUtilIsListRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		arg := strings.TrimSpace(appSyncUtilIsListRegex.FindStringSubmatch(match)[1])
		if isNullLiteral(arg) {
			return "false"
		}
		if strings.HasPrefix(arg, "[") {
			return "true"
		}
		return "false"
	})
}

func (e *Engine) processUtilMathEquals(templateStr string) string {
	return appSyncUtilMathEqualsRegex.ReplaceAllStringFunc(templateStr, func(match string) string {
		submatches := appSyncUtilMathEqualsRegex.FindStringSubmatch(match)
		left := strings.TrimSpace(submatches[1])
		right := strings.TrimSpace(submatches[2])
		return fmt.Sprintf("%v", left == right)
	})
}

func isNullLiteral(arg string) bool {
	trimmed := strings.TrimSpace(arg)
	return trimmed == "null" || trimmed == "$ctx.result" || trimmed == "$ctx.error"
}

func splitErrorArgs(args string) []string {
	var parts []string
	depth := 0
	inString := false
	stringChar := byte(0)
	current := strings.Builder{}
	for i := 0; i < len(args); i++ {
		ch := args[i]
		if inString {
			current.WriteByte(ch)
			if ch == stringChar && (i == 0 || args[i-1] != '\\') {
				inString = false
			}
			continue
		}
		if ch == '"' || ch == '\'' {
			inString = true
			stringChar = ch
			current.WriteByte(ch)
			continue
		}
		if ch == '{' || ch == '[' || ch == '(' {
			depth++
			current.WriteByte(ch)
			continue
		}
		if ch == '}' || ch == ']' || ch == ')' {
			depth--
			current.WriteByte(ch)
			continue
		}
		if ch == ',' && depth == 0 {
			parts = append(parts, current.String())
			current.Reset()
			continue
		}
		current.WriteByte(ch)
	}
	parts = append(parts, current.String())
	return parts
}

func javaToGoTimeFormat(format string) string {
	r := strings.NewReplacer(
		"yyyy", "2006",
		"yy", "06",
		"MM", "01",
		"dd", "02",
		"HH", "15",
		"hh", "03",
		"mm", "04",
		"ss", "05",
		"SSS", "000",
	)
	return r.Replace(format)
}
