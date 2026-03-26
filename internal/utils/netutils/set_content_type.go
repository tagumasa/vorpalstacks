// Package netutils provides network and HTTP utility functions for vorpalstacks.
package netutils

import "net/http"

// SetContentType sets the Content-Type header on the response writer.
//
// Parameters:
//   - w: The HTTP response writer
//   - contentType: The content type to set
//
// Example:
//
//	SetContentType(w, "application/json")
func SetContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}
