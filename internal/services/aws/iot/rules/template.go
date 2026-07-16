package rules

import (
	"vorpalstacks/internal/common/iotutil"
)

// ResolveTemplate is a thin wrapper around iotutil.ResolveTemplate for
// backward compatibility within the rules package.
func ResolveTemplate(input string, topic string, clientID string, payload map[string]interface{}) string {
	return iotutil.ResolveTemplate(input, topic, clientID, payload)
}
