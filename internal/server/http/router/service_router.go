// Package router provides HTTP request routing functionality for AWS services.
package router

import (
	"io"
	"net/http"
	"sort"
	"strings"

	"vorpalstacks/internal/server/actionregistry"
)

// Credential represents AWS signing credentials extracted from the Authorization header.
type Credential struct {
	AccessKey string
	Date      string
	Region    string
	Service   string
	Request   string
}

// ServiceRouter routes HTTP requests to the appropriate AWS service.
type ServiceRouter struct {
	index *ServiceIndex
}

// NewServiceRouter creates a new ServiceRouter with the given service index.
func NewServiceRouter(index *ServiceIndex) *ServiceRouter {
	return &ServiceRouter{
		index: index,
	}
}

// DetermineService determines the service name from the HTTP request.
func (r *ServiceRouter) DetermineService(req *http.Request) string {
	if r.index == nil {
		return ""
	}

	// Check Authorization header for signing name first (most reliable)
	if cred := r.extractCredential(req); cred != nil {
		if cred.Service != "" {
			if cred.Service == "timestream" {
				if target := req.Header.Get("X-Amz-Target"); target != "" {
					parts := strings.SplitN(target, ".", 2)
					if len(parts) == 2 {
						queryOps := map[string]bool{
							"Query": true, "CancelQuery": true, "PrepareQuery": true,
							"CreateScheduledQuery": true, "DeleteScheduledQuery": true,
							"DescribeScheduledQuery": true, "ListScheduledQueries": true,
							"UpdateScheduledQuery": true, "ExecuteScheduledQuery": true,
							"DescribeAccountSettings": true, "UpdateAccountSettings": true,
							"DescribeEndpoints": true,
						}
						if queryOps[parts[1]] {
							return "timestream-query"
						}
					}
				}
				return "timestream-write"
			}
			if mapped := mapAuthServiceToService(cred.Service); mapped != "" {
				return mapped
			}
			return cred.Service
		}
	}

	// 1. Check Action query parameter (awsQuery protocols)
	action := req.URL.Query().Get("Action")
	if action == "" && req.Method == "POST" && strings.Contains(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			req.Body = io.NopCloser(strings.NewReader(""))
		} else {
			req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			if err := req.ParseForm(); err == nil {
				action = req.Form.Get("Action")
			}
			req.Form = nil
			req.PostForm = nil
			req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		}
	}

	if action != "" {
		if svcName := mapActionToService(action); svcName != "" {
			return svcName
		}
		if svcName := r.findServiceByActionPriority(action); svcName != "" {
			return svcName
		}
	}

	// 2. Check X-Amz-Target header (JSON protocols)
	if target := req.Header.Get("X-Amz-Target"); target != "" {
		if opRef := r.index.GetOperationByTarget(target); opRef != nil {
			extractedSvc := extractServiceFromTargetHeader(target)
			if extractedSvc != "" && extractedSvc != opRef.ServiceName {
				return extractedSvc
			}
			return opRef.ServiceName
		}
		return extractServiceFromTargetHeader(target)
	}

	// 3. Check Host header for endpoint prefix
	if host := r.extractHost(req); host != "" {
		parts := strings.Split(host, ".")
		if len(parts) >= 1 {
			potentialPrefix := parts[0]
			if svcName := r.index.GetServiceByEndpointPrefix(potentialPrefix); svcName != "" {
				return svcName
			}
		}
	}

	// 4. Check path patterns for known REST API prefixes
	path := req.URL.Path
	if strings.HasPrefix(path, "/v2/email/") {
		return "email"
	}
	if strings.HasPrefix(path, "/2013-04-01/") {
		return "route53"
	}
	if strings.HasPrefix(path, "/2020-05-31/") {
		return "cloudfront"
	}
	if strings.HasPrefix(path, "/schedule-groups") || strings.HasPrefix(path, "/schedules") || strings.HasPrefix(path, "/tags/arn:aws:scheduler") {
		return "scheduler"
	}
	if isNeptunedataPath(path) {
		return "neptunedata"
	}
	if strings.HasPrefix(path, "/service/GraniteServiceVersion20100801/operation/") {
		return "monitoring"
	}

	// 5. Check path patterns (REST protocols)
	if opRef := r.index.MatchOperationByPath(req.Method, req.URL.Path); opRef != nil {
		return opRef.ServiceName
	}

	return ""
}

func extractServiceFromTargetHeader(target string) string {
	parts := strings.Split(target, ".")
	if len(parts) == 0 {
		return ""
	}
	svcPart := parts[0]
	switch svcPart {
	case "AmazonSNS":
		return "sns"
	case "AmazonSQS":
		return "sqs"
	case "AmazonSSM":
		return "ssm"
	case "AWSCognitoIdentityProviderService":
		return "cognito-idp"
	case "AWSCognitoIdentityService":
		return "cognito-identity"
	case "AWSEvents":
		return "eventbridge"
	case "Logs_20140328":
		return "logs"
	case "DynamoDB_20120810":
		return "dynamodb"
	case "Kinesis_20131202":
		return "kinesis"
	case "KMS_20141211", "TrentService":
		return "kms"
	case "secretsmanager":
		return "secretsmanager"
	case "AWSStepFunctions":
		return "states"
	case "Timestream_20181101":
		if len(parts) > 1 {
			queryOps := map[string]bool{
				"Query": true, "CancelQuery": true, "PrepareQuery": true,
				"CreateScheduledQuery": true, "DeleteScheduledQuery": true, "DescribeScheduledQuery": true,
				"ListScheduledQueries": true, "UpdateScheduledQuery": true, "ExecuteScheduledQuery": true,
				"DescribeAccountSettings": true, "UpdateAccountSettings": true,
			}
			if queryOps[parts[1]] {
				return "timestream-query"
			}
		}
		return "timestream-write"
	case "GraniteServiceVersion20100801":
		return "monitoring"
	case "AmazonAthena":
		return "athena"
	case "CertificateManager":
		return "acm"
	case "AWSWAF_20190729":
		return "wafv2"
	case "AWSWAF_20150824":
		return "waf"
	case "CloudTrail_20131101":
		return "cloudtrail"
	case "com":
		if len(parts) >= 3 {
			if parts[1] == "amazonaws" {
				switch parts[2] {
				case "cloudtrail":
					return "cloudtrail"
				}
			}
		}
	}
	return strings.ToLower(svcPart)
}

func (r *ServiceRouter) findServiceByActionPriority(action string) string {
	priorityServices := []string{"sns", "sqs", "iam", "sts", "kms", "dynamodb", "s3", "lambda", "events", "states"}

	for _, svcName := range priorityServices {
		key := svcName + ":" + action
		if _, exists := r.index.actionToOperation[key]; exists {
			return svcName
		}
	}

	var matches []string
	for key, opRef := range r.index.actionToOperation {
		if strings.HasSuffix(key, ":"+action) {
			matches = append(matches, opRef.ServiceName)
		}
	}
	if len(matches) > 0 {
		sort.Strings(matches)
		return matches[0]
	}
	return ""
}

func mapActionToService(action string) string {
	return actionregistry.LookupServiceByAction(action)
}

func mapAuthServiceToService(authService string) string {
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
	if mapped, ok := mapping[authService]; ok {
		return mapped
	}
	return ""
}

// DetermineOperation determines the operation name for the given service from the HTTP request.
func (r *ServiceRouter) DetermineOperation(req *http.Request, serviceName string) string {
	// 1. Check X-Amz-Target header
	if target := req.Header.Get("X-Amz-Target"); target != "" {
		if opRef := r.index.GetOperationByTarget(target); opRef != nil {
			if opRef.ServiceName == serviceName {
				return opRef.OpName
			}
		}
	}

	// 2. Check Action query parameter (in URL or POST body)
	action := req.URL.Query().Get("Action")
	if action == "" && req.Method == "POST" && strings.Contains(req.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		bodyBytes, err := io.ReadAll(req.Body)
		if err != nil {
			req.Body = io.NopCloser(strings.NewReader(""))
		} else {
			req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
			if err := req.ParseForm(); err == nil {
				action = req.Form.Get("Action")
			}
			req.Form = nil
			req.PostForm = nil
			req.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		}
	}
	if action != "" {
		if opRef := r.index.GetOperationByAction(serviceName, action); opRef != nil {
			return opRef.OpName
		}
	}

	// 3. Check path patterns
	if opRef := r.index.MatchOperationByPath(req.Method, req.URL.Path); opRef != nil {
		if opRef.ServiceName == serviceName {
			return opRef.OpName
		}
	}

	return ""
}

func (r *ServiceRouter) extractCredential(req *http.Request) *Credential {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return nil
	}

	if !strings.Contains(authHeader, "Credential=") {
		return nil
	}

	cred := &Credential{}
	parts := strings.Split(authHeader, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if idx := strings.Index(part, "Credential="); idx >= 0 {
			credStr := part[idx+11:]
			credParts := strings.Split(credStr, "/")
			if len(credParts) >= 5 {
				cred.AccessKey = credParts[0]
				cred.Date = credParts[1]
				cred.Region = credParts[2]
				cred.Service = credParts[3]
				cred.Request = credParts[4]
				return cred
			}
		}
	}

	return nil
}

func (r *ServiceRouter) extractHost(req *http.Request) string {
	host := req.Host
	if host == "" {
		host = req.Header.Get("Host")
	}

	if i := strings.LastIndex(host, ":"); i >= 0 && !strings.Contains(host[i:], "]") {
		host = host[:i]
	}

	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return host
}

// GetIndex returns the service index used by the router.
func (r *ServiceRouter) GetIndex() *ServiceIndex {
	return r.index
}

func isNeptunedataPath(path string) bool {
	return strings.HasPrefix(path, "/gremlin") ||
		strings.HasPrefix(path, "/opencypher") ||
		strings.HasPrefix(path, "/status") ||
		strings.HasPrefix(path, "/system") ||
		strings.HasPrefix(path, "/loader") ||
		strings.HasPrefix(path, "/propertygraph") ||
		strings.HasPrefix(path, "/sparql") ||
		strings.HasPrefix(path, "/rdf") ||
		strings.HasPrefix(path, "/ml")
}
