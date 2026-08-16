package dispatcher

import (
	"net/http"

	"vorpalstacks/internal/common/request"
)

// isCloudFrontPayloadOperation reports whether the operation is part of the
// CloudFront routing table and therefore encodes its response with the
// table's payload root rather than the generic operation-named envelope.
func isCloudFrontPayloadOperation(opName string) bool {
	_, ok := request.CloudFrontPayloadRoot(opName)
	return ok
}

// getCloudFrontPayloadRoot returns the XML payload root element configured
// for a CloudFront operation; empty for operations with no response body.
func getCloudFrontPayloadRoot(opName string) string {
	root, _ := request.CloudFrontPayloadRoot(opName)
	return root
}

// extractCloudFrontETag lifts the ETag and Location members of a CloudFront
// handler response into response headers, removing them from the body. The
// values may sit at the top level of the response map or inside the payload
// root member.
func extractCloudFrontETag(w http.ResponseWriter, response interface{}, payloadRoot string) {
	m, ok := response.(map[string]interface{})
	if !ok {
		return
	}

	var etag string
	var location string

	if inner, exists := m[payloadRoot]; exists {
		if innerMap, ok := inner.(map[string]interface{}); ok {
			if e, ok := innerMap["ETag"].(string); ok && e != "" {
				etag = e
			}
			if l, ok := innerMap["Location"].(string); ok && l != "" {
				location = l
				delete(innerMap, "Location")
			}
		}
	}

	if etag == "" {
		if e, ok := m["ETag"].(string); ok && e != "" {
			etag = e
		}
	}

	if location == "" {
		if l, ok := m["Location"].(string); ok && l != "" {
			location = l
			delete(m, "Location")
		}
	}

	if etag != "" {
		w.Header().Set("ETag", etag)
	}
	if location != "" {
		w.Header().Set("Location", location)
	}
}
