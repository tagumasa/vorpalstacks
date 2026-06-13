package rules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var templatePattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// ResolveTemplate replaces substitution tokens in the input string with
// runtime values from the MQTT message context.
func ResolveTemplate(input string, topic string, clientID string, payload map[string]interface{}) string {
	return templatePattern.ReplaceAllStringFunc(input, func(match string) string {
		inner := match[2 : len(match)-1]
		return resolveToken(inner, topic, clientID, payload)
	})
}

func resolveToken(token, topic, clientID string, payload map[string]interface{}) string {
	if strings.HasPrefix(token, "topic(") {
		return resolveTopicFunction(token, topic)
	}

	switch token {
	case "topic()":
		return topic
	case "timestamp()":
		return fmt.Sprintf("%d", time.Now().Unix())
	case "timestamp":
		return fmt.Sprintf("%d", time.Now().Unix())
	case "clientid()":
		return clientID
	case "clientid":
		return clientID
	case "payload()":
		return fmt.Sprintf("%v", payload)
	default:
		if strings.HasPrefix(token, "payload.") {
			key := strings.TrimPrefix(token, "payload.")
			if val, ok := payload[key]; ok {
				return fmt.Sprintf("%v", val)
			}
		}

		if val, ok := payload[token]; ok {
			return fmt.Sprintf("%v", val)
		}
		return "${" + token + "}"
	}
}

func resolveTopicFunction(token, topic string) string {
	argStr := token[6 : len(token)-1]
	if argStr == "" {
		return topic
	}

	parts := strings.Split(topic, "/")

	if strings.HasPrefix(argStr, "'") && strings.HasSuffix(argStr, "'") {
		return argStr[1 : len(argStr)-1]
	}

	idx, err := strconv.Atoi(argStr)
	if err != nil {
		return topic
	}
	if idx < 0 {
		idx = len(parts) + idx
	}
	if idx >= 0 && idx < len(parts) {
		return parts[idx]
	}
	return ""
}
