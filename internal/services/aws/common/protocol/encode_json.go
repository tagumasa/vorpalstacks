package protocol

import (
	"encoding/json"
	"net/http"
	"strconv"

	"vorpalstacks/internal/utils/buffer"
)

// EncodeJSONResponse writes the response as JSON to the HTTP response with appropriate headers.
func EncodeJSONResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(response); err != nil {
		return err
	}

	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(buf.Bytes())
	return err
}

// EncodeJSONResponseWithStatus writes the response as JSON with the specified HTTP status code.
func EncodeJSONResponseWithStatus(w http.ResponseWriter, statusCode int, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	buf := buffer.GlobalPool.Get(4096)
	defer buffer.GlobalPool.Put(buf)

	if err := json.NewEncoder(buf).Encode(response); err != nil {
		return err
	}

	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.WriteHeader(statusCode)
	_, err := w.Write(buf.Bytes())
	return err
}
