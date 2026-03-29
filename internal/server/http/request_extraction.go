package http

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"vorpalstacks/internal/server/actionregistry"
	"vorpalstacks/internal/server/http/router"
)

func (s *Server) extractServiceFromRequest(r *http.Request) string {
	if s.serviceRouter != nil {
		if svcName := s.serviceRouter.DetermineService(r); svcName != "" {
			return svcName
		}
	}
	return extractServiceFromRequestFallback(r)
}

func extractServiceFromRequestFallback(r *http.Request) string {
	if router.IsLambdaRestPath(r.URL.Path) {
		return "lambda"
	}

	if isApiGatewayPath(r.URL.Path) {
		return "apigateway"
	}

	if isSchedulerPath(r.URL.Path) {
		return "scheduler"
	}

	if isSESv2Path(r.URL.Path) {
		return "email"
	}

	if isRoute53Path(r.URL.Path) {
		return "route53"
	}

	if isCloudFrontPath(r.URL.Path) {
		return "cloudfront"
	}

	xAmzTarget := r.Header.Get("X-Amz-Target")
	if xAmzTarget != "" {
		parts := splitByDot(xAmzTarget)
		if len(parts) > 0 {
			svc := extractServiceFromTarget(parts[0])
			if svc == "timestream-write" && len(parts) > 1 {
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
			return svc
		}
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		if strings.Contains(authHeader, "Credential=") {
			parts := strings.Split(authHeader, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if idx := strings.Index(part, "Credential="); idx >= 0 {
					cred := part[idx+11:]
					credParts := strings.Split(cred, "/")
					if len(credParts) >= 5 {
						return credParts[3]
					}
				}
			}
		}
	}

	action := r.URL.Query().Get("Action")
	if action == "" && r.Method == "POST" && strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			r.Body = io.NopCloser(bytes.NewReader([]byte{}))
			return ""
		}
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		if err := r.ParseForm(); err == nil {
			action = r.Form.Get("Action")
		}
		r.Form = nil
		r.PostForm = nil
		r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	if action != "" {
		if svc := extractServiceFromAction(action); svc != "" {
			return svc
		}
	}

	host := r.Host
	if i := strings.LastIndex(host, ":"); i >= 0 && !strings.Contains(host[i:], "]") {
		host = host[:i]
	}
	parts := splitByDot(host)
	if len(parts) >= 1 && parts[0] != "localhost" && !isIPAddress(parts[0]) {
		return parts[0]
	}

	return ""
}

func extractServiceFromAction(action string) string {
	return actionregistry.LookupServiceByAction(action)
}

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

func extractServiceFromTarget(target string) string {
	switch target {
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
		return "timestream-write"
	case "AmazonAthena":
		return "athena"
	case "AWSSTS":
		return "sts"
	case "AWSWAF_20150824":
		return "waf"
	case "AWSWAF_20190729":
		return "wafv2"
	}
	if len(target) > 0 {
		switch strings.ToLower(target) {
		case "amazonsns":
			return "sns"
		case "amazonsqs":
			return "sqs"
		case "amazonssm":
			return "ssm"
		case "amazonlambda":
			return "lambda"
		case "amazoncloudwatch":
			return "cloudwatch"
		case "amazondynamodb":
			return "dynamodb"
		case "amazonkinesis":
			return "kinesis"
		case "amazonevents":
			return "eventbridge"
		case "amazonathena":
			return "athena"
		case "awscognitoidentityproviderservice":
			return "cognito-idp"
		case "awscognitoidentityservice":
			return "cognito-identity"
		case "awsstepfunctions":
			return "states"
		case "amazoncloudwatchlogs":
			return "logs"
		case "timestream_20181101":
			return "timestream-write"
		case "timestreamquery_20181101":
			return "timestream-query"
		}
		return strings.ToLower(target)
	}
	return ""
}

func splitByDot(s string) []string {
	var result []string
	var current strings.Builder
	for _, c := range s {
		if c == '.' {
			if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(c)
		}
	}
	if current.Len() > 0 {
		result = append(result, current.String())
	}
	return result
}
