package protocol

import (
	"net/http"

	"github.com/fxamacker/cbor/v2"
)

var cborEncMode cbor.EncMode

func init() {
	opts := cbor.EncOptions{
		Time:    cbor.TimeUnixDynamic,
		TimeTag: cbor.EncTagRequired,
	}
	var err error
	cborEncMode, err = opts.EncMode()
	if err != nil {
		panic("failed to initialize CBOR encoder: " + err.Error())
	}
}

// EncodeCBORResponse writes the response as CBOR to the HTTP response with appropriate headers.
func EncodeCBORResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/cbor")
	w.Header().Set("smithy-protocol", "rpc-v2-cbor")
	data, err := cborEncMode.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
