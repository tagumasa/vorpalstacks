package classifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/fxamacker/cbor/v2"

	"vorpalstacks/internal/server/actionregistry"
	"vorpalstacks/internal/server/http/router"
	"vorpalstacks/internal/services/aws/common/request"
	"vorpalstacks/internal/utils/aws/authutil"
)

// Protocol represents the AWS request protocol type.
type Protocol int

const (
	ProtocolUnknown Protocol = iota
	ProtocolJSONRPC
	ProtocolQuery
	ProtocolRESTXML
	ProtocolRESTJSON
	ProtocolCBOR
	ProtocolNeptune
)

// ClassifiedRequest holds the result of classifying an incoming AWS HTTP request.
type ClassifiedRequest struct {
	Protocol    Protocol
	ServiceName string
	Operation   string
	Parameters  map[string]interface{}
	Headers     http.Header
	QueryParams url.Values
	PathParams  map[string]string
	Body        []byte
	Region      string
	AccessKeyID string
}

// Classifier determines the AWS service, operation, protocol, and parameters
// from an incoming HTTP request using a unified classification pipeline.
type Classifier struct {
	index *router.ServiceIndex
}

// NewClassifier creates a new Classifier backed by the given service index.
func NewClassifier(index *router.ServiceIndex) *Classifier {
	return &Classifier{index: index}
}

// Classify analyses an incoming HTTP request and returns a ClassifiedRequest
// containing the service name, operation, protocol, parameters, and
// authentication details. The request body is read, parsed, and then
// restored so that downstream handlers may re-read it.
func (c *Classifier) Classify(r *http.Request) (*ClassifiedRequest, error) {
	cr := &ClassifiedRequest{
		Headers:     r.Header,
		QueryParams: r.URL.Query(),
		PathParams:  make(map[string]string),
		Parameters:  make(map[string]interface{}),
	}

	bodyBytes, err := c.readBody(r)
	if err != nil {
		return nil, err
	}
	cr.Body = bodyBytes
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	c.extractAuth(r, cr)

	cr.Protocol = c.detectProtocol(r, bodyBytes)

	cr.ServiceName = c.determineServiceName(r, bodyBytes, cr)

	cr.Operation = c.determineOperation(r, bodyBytes, cr)

	c.parseBody(r, bodyBytes, cr)

	c.mergeQueryParams(cr)

	request.ExtractRESTPathParams(r.URL.Path, r.Method, cr.Parameters)

	if cr.Region == "" {
		cr.Region = request.DefaultRegion
	}

	return cr, nil
}

func (c *Classifier) readBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 64*1024*1024))
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func (c *Classifier) extractAuth(r *http.Request, cr *ClassifiedRequest) {
	cred := authutil.ParseAWSCredential(r.Header.Get("Authorization"))
	if cred == nil {
		return
	}
	cr.AccessKeyID = cred.AccessKeyID
	cr.Region = cred.Region
}

func (c *Classifier) detectProtocol(r *http.Request, bodyBytes []byte) Protocol {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/cbor") {
		return ProtocolCBOR
	}
	if r.Header.Get("X-Amz-Target") != "" {
		return ProtocolJSONRPC
	}
	if r.URL.Query().Get("Action") != "" {
		return ProtocolQuery
	}
	if r.Method == "POST" && strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return ProtocolQuery
	}
	if isNeptunedataPath(r.URL.Path) {
		return ProtocolNeptune
	}
	if strings.HasPrefix(r.URL.Path, "/2013-04-01/") || strings.HasPrefix(r.URL.Path, "/2020-05-31/") {
		return ProtocolRESTXML
	}
	if isLambdaRestPath(r.URL.Path) ||
		isApiGatewayPath(r.URL.Path) ||
		strings.HasPrefix(r.URL.Path, "/schedule-groups") ||
		strings.HasPrefix(r.URL.Path, "/schedules") ||
		strings.HasPrefix(r.URL.Path, "/v2/email/") ||
		strings.HasPrefix(r.URL.Path, "/service/GraniteServiceVersion20100801/") ||
		isAppSyncPath(r.URL.Path) {
		return ProtocolRESTJSON
	}
	return ProtocolUnknown
}

func (c *Classifier) determineServiceName(r *http.Request, bodyBytes []byte, cr *ClassifiedRequest) string {
	if svc := c.serviceFromSigningService(r, bodyBytes); svc != "" {
		return svc
	}
	if svc := c.serviceFromTarget(r); svc != "" {
		return svc
	}
	if svc := c.serviceFromAction(r, bodyBytes); svc != "" {
		return svc
	}
	if svc := c.serviceFromHost(r); svc != "" {
		return svc
	}
	if svc := lookupServiceByPath(r.URL.Path); svc != "" {
		return svc
	}
	if c.index != nil {
		if opRef := c.index.MatchOperationByPath(r.Method, r.URL.Path); opRef != nil {
			return opRef.ServiceName
		}
	}
	return ""
}

func (c *Classifier) serviceFromSigningService(r *http.Request, bodyBytes []byte) string {
	cred := authutil.ParseAWSCredential(r.Header.Get("Authorization"))
	if cred == nil {
		return ""
	}
	signingService := cred.SigningService

	mapping := map[string]string{
		"timestream":        "timestream-write",
		"email":             "email",
		"ses":               "email",
		"logs":              "logs",
		"states":            "states",
		"runtime.sagemaker": "lambda",
		"wafv2":             "wafv2",
		"events":            "eventbridge",
		"rds":               "neptune",
		"neptune-db":        "neptunedata",
	}
	if mapped, ok := mapping[signingService]; ok {
		if signingService == "timestream" {
			return c.timestreamSplit(r, bodyBytes, "timestream-write")
		}
		return mapped
	}
	return signingService
}

func (c *Classifier) timestreamSplit(r *http.Request, bodyBytes []byte, defaultSvc string) string {
	queryOps := map[string]bool{
		"Query": true, "CancelQuery": true, "PrepareQuery": true,
		"CreateScheduledQuery": true, "DeleteScheduledQuery": true,
		"DescribeScheduledQuery": true, "ListScheduledQueries": true,
		"UpdateScheduledQuery": true, "ExecuteScheduledQuery": true,
		"DescribeAccountSettings": true, "UpdateAccountSettings": true,
		"DescribeEndpoints": true,
		"TagResource":       true, "UntagResource": true, "ListTagsForResource": true,
	}

	xAmzTarget := r.Header.Get("X-Amz-Target")
	if xAmzTarget != "" {
		parts := strings.SplitN(xAmzTarget, ".", 2)
		if len(parts) == 2 && queryOps[parts[1]] {
			return "timestream-query"
		}
	}

	action := r.URL.Query().Get("Action")
	if action == "" && len(bodyBytes) > 0 {
		if values, err := url.ParseQuery(string(bodyBytes)); err == nil {
			action = values.Get("Action")
		}
	}
	if queryOps[action] {
		return "timestream-query"
	}

	return defaultSvc
}

func (c *Classifier) serviceFromTarget(r *http.Request) string {
	target := r.Header.Get("X-Amz-Target")
	if target == "" {
		return ""
	}
	parts := strings.SplitN(target, ".", 2)
	if len(parts) == 0 {
		return ""
	}
	svc := lookupServiceByTarget(target)
	if svc == "timestream-write" && len(parts) > 1 {
		return c.timestreamSplit(r, nil, "timestream-write")
	}
	return svc
}

func (c *Classifier) serviceFromAction(r *http.Request, bodyBytes []byte) string {
	action := r.URL.Query().Get("Action")
	if action == "" && len(bodyBytes) > 0 {
		if values, err := url.ParseQuery(string(bodyBytes)); err == nil {
			action = values.Get("Action")
		}
	}
	if action == "" {
		return ""
	}
	if svc := actionregistry.LookupServiceByAction(action); svc != "" {
		return svc
	}
	if c.index != nil {
		return c.index.FindServiceByActionPriority(action)
	}
	return ""
}

func (c *Classifier) serviceFromHost(r *http.Request) string {
	if c.index == nil {
		return ""
	}
	host := r.Host
	if host == "" {
		host = r.Header.Get("Host")
	}
	if i := strings.LastIndex(host, ":"); i >= 0 && !strings.Contains(host[i:], "]") {
		host = host[:i]
	}
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}
	parts := strings.Split(host, ".")
	if len(parts) >= 1 && parts[0] != "localhost" && !isIPAddress(parts[0]) {
		return c.index.GetServiceByEndpointPrefix(parts[0])
	}
	return ""
}

// isIPAddress reports whether s is an IPv4 or IPv6 address literal.
func isIPAddress(s string) bool {
	if s == "" {
		return false
	}
	if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
		s = s[1 : len(s)-1]
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == ':' {
			return true
		}
		if c != '.' && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

func (c *Classifier) determineOperation(r *http.Request, bodyBytes []byte, cr *ClassifiedRequest) string {
	xAmzTarget := r.Header.Get("X-Amz-Target")
	if xAmzTarget != "" {
		parts := strings.Split(xAmzTarget, ".")
		if len(parts) > 1 {
			return parts[len(parts)-1]
		}
		return xAmzTarget
	}

	if strings.HasPrefix(r.URL.Path, "/service/") {
		pathParts := strings.Split(r.URL.Path, "/")
		for i, p := range pathParts {
			if p == "operation" && i+1 < len(pathParts) {
				return pathParts[i+1]
			}
		}
	}

	action := r.URL.Query().Get("Action")
	if action != "" {
		return action
	}
	if len(bodyBytes) > 0 {
		if values, err := url.ParseQuery(string(bodyBytes)); err == nil {
			if a := values.Get("Action"); a != "" {
				return a
			}
		}
	}

	if op := request.ExtractRESTOperation(r, bodyBytes); op != "" {
		return op
	}

	return ""
}

func (c *Classifier) parseBody(r *http.Request, bodyBytes []byte, cr *ClassifiedRequest) {
	if len(bodyBytes) == 0 {
		return
	}

	contentType := r.Header.Get("Content-Type")

	switch {
	case strings.Contains(contentType, "application/cbor"):
		var rawParams interface{}
		if err := cbor.Unmarshal(bodyBytes, &rawParams); err == nil {
			if m, ok := convertCBORMapToStringMap(rawParams).(map[string]interface{}); ok {
				cr.Parameters = m
			}
		}
	case strings.Contains(contentType, "application/json") || strings.Contains(contentType, "application/x-amz-json"):
		if err := json.Unmarshal(bodyBytes, &cr.Parameters); err != nil {
			cr.Parameters = make(map[string]interface{})
		}
	case strings.Contains(contentType, "application/x-www-form-urlencoded"):
		if values, err := url.ParseQuery(string(bodyBytes)); err == nil {
			for k, v := range values {
				if len(v) == 1 {
					cr.Parameters[k] = v[0]
				} else {
					cr.Parameters[k] = v
				}
			}
		}
	case strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") ||
		strings.HasPrefix(r.URL.Path, "/2013-04-01/") || strings.HasPrefix(r.URL.Path, "/2020-05-31/"):
		request.ParseXMLBody(bodyBytes, cr.Parameters)
	default:
		if err := json.Unmarshal(bodyBytes, &cr.Parameters); err != nil {
			if values, err := url.ParseQuery(string(bodyBytes)); err == nil {
				for k, v := range values {
					if len(v) == 1 {
						cr.Parameters[k] = v[0]
					} else {
						cr.Parameters[k] = v
					}
				}
			}
		}
	}
}

func (c *Classifier) mergeQueryParams(cr *ClassifiedRequest) {
	for k, v := range cr.QueryParams {
		if len(v) == 1 {
			cr.Parameters[k] = v[0]
		} else {
			cr.Parameters[k] = v
		}
	}
}

// convertCBORMapToStringMap recursively converts CBOR-decoded maps
// (which use interface{} keys) into maps with string keys.
func convertCBORMapToStringMap(v interface{}) interface{} {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			keyStr, ok := k.(string)
			if !ok {
				switch kv := k.(type) {
				case string:
					keyStr = kv
				case []byte:
					keyStr = string(kv)
				default:
					keyStr = fmt.Sprintf("%v", kv)
				}
			}
			result[keyStr] = convertCBORMapToStringMap(v)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = convertCBORMapToStringMap(item)
		}
		return result
	default:
		return v
	}
}
