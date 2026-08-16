package request

import (
	"net/http"
	"net/url"
	"strings"
)

// cloudFrontPathValue decodes a percent-encoded path segment. Path labels
// such as the ListDistributionsByWebACLId WebACLId accept WAFv2 ARNs,
// whose slashes clients send escaped (%2F); the path must be split before
// decoding so those slashes stay inside a single segment.
func cloudFrontPathValue(seg string) string {
	if !strings.Contains(seg, "%") {
		return seg
	}
	if v, err := url.PathUnescape(seg); err == nil {
		return v
	}
	return seg
}

// cloudFrontRESTParser implements RESTServiceParser for Amazon CloudFront.
type cloudFrontRESTParser struct{}

// MatchPath returns true if the path belongs to CloudFront.
func (p *cloudFrontRESTParser) MatchPath(path string) bool {
	return strings.HasPrefix(path, "/2020-05-31/")
}

// ExtractOperation returns the CloudFront operation name, or empty if the path does not match.
func (p *cloudFrontRESTParser) ExtractOperation(r *http.Request) string {
	op, _ := matchCloudFrontRoute(r.Method, r.URL.EscapedPath(), r.URL.Query())
	return op
}

// ExtractPathParams extracts URI-bound parameters from the CloudFront request path.
func (p *cloudFrontRESTParser) ExtractPathParams(r *http.Request, params map[string]interface{}) {
	_, labels := matchCloudFrontRoute(r.Method, r.URL.EscapedPath(), r.URL.Query())
	for name, value := range labels {
		if _, ok := params[name]; !ok {
			params[name] = value
		}
	}
}
