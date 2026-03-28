package s3

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// RangeSpec represents a byte range specification for GetObject requests.
type RangeSpec struct {
	Start  int64
	End    int64
	Length int64
}

func parseRangeHeader(rangeHeader string) (ranges []RangeSpec, err error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, fmt.Errorf("invalid range header: missing bytes= prefix")
	}

	rangeSpec := strings.TrimPrefix(rangeHeader, "bytes=")

	for _, spec := range strings.Split(rangeSpec, ",") {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}

		parts := strings.SplitN(spec, "-", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range spec: %s", spec)
		}

		var r RangeSpec
		if parts[0] == "" {
			if parts[1] == "" {
				return nil, fmt.Errorf("invalid range spec: %s", spec)
			}
			suffixLength, parseErr := strconv.ParseInt(parts[1], 10, 64)
			if parseErr != nil {
				return nil, fmt.Errorf("invalid suffix length: %s", parts[1])
			}
			if suffixLength <= 0 {
				return nil, fmt.Errorf("suffix length must be positive: %d", suffixLength)
			}
			r.Start = -1
			r.Length = suffixLength
		} else if parts[1] == "" {
			start, parseErr := strconv.ParseInt(parts[0], 10, 64)
			if parseErr != nil {
				return nil, fmt.Errorf("invalid start offset: %s", parts[0])
			}
			if start < 0 {
				return nil, fmt.Errorf("start offset cannot be negative: %d", start)
			}
			r.Start = start
			r.End = -1
			r.Length = -1
		} else {
			start, parseErr := strconv.ParseInt(parts[0], 10, 64)
			if parseErr != nil {
				return nil, fmt.Errorf("invalid start offset: %s", parts[0])
			}
			end, parseErr := strconv.ParseInt(parts[1], 10, 64)
			if parseErr != nil {
				return nil, fmt.Errorf("invalid end offset: %s", parts[1])
			}
			if start < 0 || end < 0 {
				return nil, fmt.Errorf("offsets cannot be negative: start=%d, end=%d", start, end)
			}
			if start > end {
				return nil, fmt.Errorf("start offset cannot be greater than end: start=%d, end=%d", start, end)
			}
			r.Start = start
			r.End = end
			r.Length = end - start + 1
		}

		ranges = append(ranges, r)
	}

	if len(ranges) == 0 {
		return nil, fmt.Errorf("no valid range specified")
	}

	return ranges, nil
}

func parseCopySource(copySource string) (bucket, key, versionId string, err error) {
	copySource = strings.TrimPrefix(copySource, "/")

	if decoded, decodeErr := url.PathUnescape(copySource); decodeErr == nil {
		copySource = decoded
	}

	if idx := strings.Index(copySource, "?"); idx != -1 {
		queryPart := copySource[idx+1:]
		copySource = copySource[:idx]

		params := strings.Split(queryPart, "&")
		for _, param := range params {
			if strings.HasPrefix(param, "versionId=") {
				versionId = strings.TrimPrefix(param, "versionId=")
				break
			}
		}
	}

	parts := strings.SplitN(copySource, "/", 2)
	if len(parts) < 2 {
		return "", "", "", fmt.Errorf("invalid copy source")
	}
	key = parts[1]
	if err := validateObjectKey(key); err != nil {
		return "", "", "", err
	}
	return parts[0], key, versionId, nil
}
