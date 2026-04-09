package s3

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var errAwsChunkedInvalidFormat = errors.New("invalid aws-chunked encoding format")

// awsChunkedDecoder decodes an HTTP request body encoded with Transfer-Encoding: aws-chunked.
type awsChunkedDecoder struct {
	reader    *bufio.Reader
	remaining int
	eof       bool
}

// Read decodes the next chunk from the aws-chunked transfer encoding stream.
func (d *awsChunkedDecoder) Read(p []byte) (int, error) {
	if d.eof {
		return 0, io.EOF
	}

	if d.remaining > 0 {
		n, err := d.reader.Read(p)
		if n > d.remaining {
			n = d.remaining
		}
		d.remaining -= n
		if err != nil {
			return n, err
		}
		if d.remaining > 0 {
			return n, nil
		}
		if _, err := d.reader.ReadString('\n'); err != nil {
			d.eof = true
			return n, err
		}
		return n, nil
	}

	for {
		line, err := d.reader.ReadString('\n')
		if err != nil {
			d.eof = true
			return 0, err
		}

		line = strings.TrimRight(line, "\r\n")
		semicolonIdx := strings.IndexByte(line, ';')
		var sizeStr string
		if semicolonIdx >= 0 {
			sizeStr = line[:semicolonIdx]
		} else {
			sizeStr = line
		}

		size, err := strconv.ParseInt(sizeStr, 16, 64)
		if err != nil {
			d.eof = true
			return 0, fmt.Errorf("%w: failed to parse chunk size %q: %v", errAwsChunkedInvalidFormat, sizeStr, err)
		}

		if size == 0 {
			d.eof = true
			d.reader.ReadString('\n')
			return 0, io.EOF
		}

		d.remaining = int(size)
		n, err := d.reader.Read(p)
		if n > d.remaining {
			n = d.remaining
		}
		d.remaining -= n
		if d.remaining > 0 {
			return n, err
		}
		if _, discardErr := d.reader.ReadString('\n'); discardErr != nil {
			d.eof = true
		}
		return n, err
	}
}

func newAwsChunkedDecoder(r io.Reader) io.Reader {
	return &awsChunkedDecoder{reader: bufio.NewReader(r)}
}

func isAwsChunkedRequest(r *http.Request) bool {
	if strings.Contains(r.Header.Get("Content-Encoding"), "aws-chunked") {
		return true
	}
	contentSha := r.Header.Get("X-Amz-Content-Sha256")
	return strings.HasPrefix(contentSha, "STREAMING-")
}

func decodeAwsChunkedBody(r io.Reader) io.Reader {
	return newAwsChunkedDecoder(r)
}
