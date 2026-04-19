package iam

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	cloudtrailstore "vorpalstacks/internal/store/aws/cloudtrail"
	iamstore "vorpalstacks/internal/store/aws/iam"
	arnutil "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"
)

var (
	// ErrNoSuchJob is returned when a job ID does not match any stored report.
	ErrNoSuchJob = errors.NewAWSError("NoSuchEntity", "The request was rejected because no job was found with the provided job ID.", http.StatusNotFound)
)

// eventSourceToServiceNamespace maps CloudTrail event sources to AWS service namespace identifiers.
var eventSourceToServiceNamespace = map[string]string{
	"s3.amazonaws.com":             "s3",
	"lambda.amazonaws.com":         "lambda",
	"dynamodb.amazonaws.com":       "dynamodb",
	"sqs.amazonaws.com":            "sqs",
	"sns.amazonaws.com":            "sns",
	"kms.amazonaws.com":            "kms",
	"cloudtrail.amazonaws.com":     "cloudtrail",
	"monitoring.amazonaws.com":     "cloudwatch",
	"logs.amazonaws.com":           "logs",
	"iam.amazonaws.com":            "iam",
	"sts.amazonaws.com":            "sts",
	"events.amazonaws.com":         "events",
	"states.amazonaws.com":         "states",
	"cognito-idp.amazonaws.com":    "cognito-idp",
	"apigateway.amazonaws.com":     "apigateway",
	"cloudfront.amazonaws.com":     "cloudfront",
	"route53.amazonaws.com":        "route53",
	"acm.amazonaws.com":            "acm",
	"secretsmanager.amazonaws.com": "secretsmanager",
	"athena.amazonaws.com":         "athena",
	"kinesis.amazonaws.com":        "kinesis",
	"ssm.amazonaws.com":            "ssm",
	"email.amazonaws.com":          "ses",
	"timestream.amazonaws.com":     "timestream",
	"waf.amazonaws.com":            "waf",
}

// namespaceToDisplayName maps service namespaces to their AWS display names.
var namespaceToDisplayName = map[string]string{
	"lambda":         "AWS Lambda",
	"dynamodb":       "Amazon DynamoDB",
	"s3":             "Amazon S3",
	"sqs":            "Amazon SQS",
	"sns":            "Amazon SNS",
	"kms":            "AWS KMS",
	"iam":            "AWS Identity and Access Management",
	"cloudtrail":     "AWS CloudTrail",
	"cloudwatch":     "Amazon CloudWatch",
	"monitoring":     "Amazon CloudWatch",
	"logs":           "Amazon CloudWatch Logs",
	"sts":            "AWS Security Token Service",
	"events":         "Amazon EventBridge",
	"states":         "AWS Step Functions",
	"secretsmanager": "AWS Secrets Manager",
	"apigateway":     "Amazon API Gateway",
	"cloudfront":     "Amazon CloudFront",
	"route53":        "Amazon Route 53",
	"athena":         "Amazon Athena",
	"kinesis":        "Amazon Kinesis",
	"ssm":            "AWS Systems Manager",
	"ses":            "Amazon SES",
	"acm":            "AWS Certificate Manager",
	"cognito-idp":    "Amazon Cognito Identity Provider",
	"timestream":     "Amazon Timestream",
	"waf":            "AWS WAF",
}

// generateJobID produces a unique 32-character lowercase hex identifier for a report job.
func generateJobID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%032x", time.Now().UnixNano())
	}
	return strings.ToLower(hex.EncodeToString(b))
}

// parseGranularity converts an ISO 8601 duration string to a time.Duration.
// Supported values are P7D, P30D, and P90D; defaults to P30D.
func parseGranularity(granularity string) time.Duration {
	switch granularity {
	case "P7D":
		return 7 * 24 * time.Hour
	case "P30D":
		return 30 * 24 * time.Hour
	case "P90D":
		return 90 * 24 * time.Hour
	default:
		return 30 * 24 * time.Hour
	}
}

// resolveEntityName extracts the entity name (user, role, or group) from an IAM ARN.
func resolveEntityName(arn string) string {
	arn = strings.TrimPrefix(arn, "arn:aws:iam::")

	for _, prefix := range []string{"user/", "role/", "group/"} {
		if idx := strings.Index(arn, prefix); idx >= 0 {
			return arn[idx+len(prefix):]
		}
	}
	return arn
}

// generateLastAccessedReport queries CloudTrail events for the given entity within the
// specified time granularity and produces a ServiceLastAccessedJob with aggregated results.
func (s *IAMService) generateLastAccessedReport(arn, granularity, jobType string, ctStore cloudtrailstore.CloudTrailStoreInterface) *iamstore.ServiceLastAccessedJob {
	now := time.Now().UTC()
	duration := parseGranularity(granularity)
	startTime := now.Add(-duration)
	entityName := resolveEntityName(arn)

	query := cloudtrailstore.EventQuery{
		MaxResults: 1000,
	}

	var allEvents []*cloudtrailstore.Event
	nextToken := ""
	for {
		query.NextToken = nextToken
		events, token, err := ctStore.LookupEvents(query)
		if err != nil {
			break
		}
		allEvents = append(allEvents, events...)
		if token == "" {
			break
		}
		nextToken = token
	}

	filteredEvents := make([]*cloudtrailstore.Event, 0, len(allEvents))
	for _, event := range allEvents {
		if event == nil || event.UserIdentity == nil {
			continue
		}
		if event.UserIdentity.UserName != entityName {
			continue
		}
		if event.EventTime.Before(startTime) || event.EventTime.After(now) {
			continue
		}
		filteredEvents = append(filteredEvents, event)
	}

	serviceMap := make(map[string]*iamstore.ServiceLastAccessed)
	actionMap := make(map[string]*iamstore.TrackedActionLastAccessed)

	for _, event := range filteredEvents {
		if event == nil {
			continue
		}

		eventRegion := request.DefaultRegion
		if event.UserIdentity != nil && event.UserIdentity.ARN != "" {
			if _, _, r, _, _ := arnutil.SplitARN(event.UserIdentity.ARN); r != "" {
				eventRegion = r
			}
		} else if len(event.Resources) > 0 {
			for _, res := range event.Resources {
				if res.ARN != "" {
					if _, _, r, _, _ := arnutil.SplitARN(res.ARN); r != "" {
						eventRegion = r
						break
					}
				}
			}
		}

		serviceNamespace := eventSourceToServiceNamespace[event.EventSource]
		if serviceNamespace == "" {
			parts := strings.SplitN(event.EventSource, ".", 2)
			if len(parts) > 0 {
				serviceNamespace = parts[0]
			} else {
				serviceNamespace = event.EventSource
			}
		}

		serviceName := namespaceToDisplayName[serviceNamespace]
		if serviceName == "" {
			serviceName = serviceNamespace
		}

		svc, ok := serviceMap[serviceNamespace]
		if !ok {
			svc = &iamstore.ServiceLastAccessed{
				ServiceName:                serviceName,
				ServiceNamespace:           serviceNamespace,
				TotalAuthenticatedEntities: 1,
			}
			serviceMap[serviceNamespace] = svc
		}

		if svc.LastAuthenticated == nil || event.EventTime.After(svc.LastAuthenticated.UTC()) {
			t := event.EventTime
			svc.LastAuthenticated = &t
			svc.LastAuthenticatedRegion = eventRegion
		}

		actionKey := serviceNamespace + ":" + event.EventName
		action, ok := actionMap[actionKey]
		if !ok {
			action = &iamstore.TrackedActionLastAccessed{
				ActionName:       event.EventName,
				ServiceNamespace: serviceNamespace,
				EntityPath:       arn,
			}
			actionMap[actionKey] = action
		}

		if action.LastAccessedDate == nil || event.EventTime.After(action.LastAccessedDate.UTC()) {
			t := event.EventTime
			action.LastAccessedDate = &t
			action.LastAccessedRegion = eventRegion
		}
	}

	services := make([]iamstore.ServiceLastAccessed, 0, len(serviceMap))
	for _, svc := range serviceMap {
		actions := make([]iamstore.TrackedActionLastAccessed, 0)
		for _, a := range actionMap {
			if a.ServiceNamespace == svc.ServiceNamespace {
				actions = append(actions, *a)
			}
		}
		svc.TrackedActionsLastAccessed = actions
		services = append(services, *svc)
	}

	actions := make([]iamstore.TrackedActionLastAccessed, 0, len(actionMap))
	for _, a := range actionMap {
		actions = append(actions, *a)
	}

	completionTime := now
	job := &iamstore.ServiceLastAccessedJob{
		JobID:                      generateJobID(),
		Arn:                        arn,
		JobType:                    jobType,
		JobStatus:                  "COMPLETED",
		JobCreationTime:            now,
		JobCompletionTime:          &completionTime,
		Granularity:                granularity,
		TrackedActionsLastAccessed: actions,
		ServicesLastAccessed:       services,
	}

	return job
}

// GenerateServiceLastAccessedDetails generates a report of the last time the specified IAM entity
// accessed each AWS service. The report is generated synchronously and stored for retrieval via
// GetServiceLastAccessedDetails.
func (s *IAMService) GenerateServiceLastAccessedDetails(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	arn := request.GetStringParam(req.Parameters, "Arn")
	if arn == "" {
		arn = request.GetStringParam(req.Parameters, "arn")
	}
	granularity := request.GetStringParam(req.Parameters, "Granularity")
	if granularity == "" {
		granularity = request.GetStringParam(req.Parameters, "granularity")
	}
	if granularity == "" {
		granularity = "P30D"
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	regionStorage, err := reqCtx.GetStorage()
	if err != nil {
		return nil, errors.NewAWSError("InternalFailure", "Failed to access CloudTrail storage", http.StatusInternalServerError)
	}
	ctStore := cloudtrailstore.NewCloudTrailStore(regionStorage, s.accountID, reqCtx.GetRegion())

	job := s.generateLastAccessedReport(arn, granularity, "SERVICE_LAST_ACCESSED", ctStore)

	if err := store.ServiceLastAccessed().Put(job); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"JobId":     job.JobID,
		"JobType":   job.JobType,
		"JobStatus": "COMPLETED",
	}, nil
}

// GetServiceLastAccessedDetails retrieves a previously generated report containing the last time
// each AWS service was accessed by the specified entity.
func (s *IAMService) GetServiceLastAccessedDetails(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetStringParam(req.Parameters, "JobId")
	if jobID == "" {
		jobID = request.GetStringParam(req.Parameters, "jobId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.ServiceLastAccessed().Get(jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrNoSuchJob
	}

	result := map[string]interface{}{
		"JobCreationDate": job.JobCreationTime.Format(timeutils.ISO8601SimpleFormat),
		"JobStatus":       job.JobStatus,
		"JobType":         job.JobType,
		"JobId":           job.JobID,
	}

	if job.JobCompletionTime != nil {
		result["JobCompletionDate"] = job.JobCompletionTime.Format(timeutils.ISO8601SimpleFormat)
	}
	if job.Error != "" {
		result["Error"] = job.Error
	}

	services := make([]map[string]interface{}, 0, len(job.ServicesLastAccessed))
	for _, svc := range job.ServicesLastAccessed {
		entry := map[string]interface{}{
			"ServiceName":                svc.ServiceName,
			"ServiceNamespace":           svc.ServiceNamespace,
			"TotalAuthenticatedEntities": svc.TotalAuthenticatedEntities,
		}
		if svc.LastAuthenticated != nil {
			entry["LastAuthenticated"] = svc.LastAuthenticated.Format(timeutils.ISO8601SimpleFormat)
		}
		if svc.LastAuthenticatedRegion != "" {
			entry["LastAuthenticatedRegion"] = svc.LastAuthenticatedRegion
		}
		actions := make([]map[string]interface{}, 0, len(svc.TrackedActionsLastAccessed))
		for _, a := range svc.TrackedActionsLastAccessed {
			actionEntry := map[string]interface{}{
				"ActionName":         a.ActionName,
				"ServiceNamespace":   a.ServiceNamespace,
				"LastAccessedEntity": a.EntityPath,
			}
			if a.LastAccessedDate != nil {
				actionEntry["LastAccessedTime"] = a.LastAccessedDate.Format(timeutils.ISO8601SimpleFormat)
			}
			if a.LastAccessedRegion != "" {
				actionEntry["LastAccessedRegion"] = a.LastAccessedRegion
			}
			actions = append(actions, actionEntry)
		}
		if len(actions) > 0 {
			entry["TrackedActionsLastAccessed"] = actions
		}
		services = append(services, entry)
	}
	result["ServicesLastAccessed"] = services

	return result, nil
}

// GetServiceLastAccessedDetailsWithEntities retrieves a previously generated report including
// entity-level detail for each service access event.
func (s *IAMService) GetServiceLastAccessedDetailsWithEntities(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetStringParam(req.Parameters, "JobId")
	if jobID == "" {
		jobID = request.GetStringParam(req.Parameters, "jobId")
	}

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := store.ServiceLastAccessed().Get(jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrNoSuchJob
	}

	result := map[string]interface{}{
		"JobCreationDate": job.JobCreationTime.Format(timeutils.ISO8601SimpleFormat),
		"JobStatus":       job.JobStatus,
		"JobType":         job.JobType,
		"JobId":           job.JobID,
	}

	if job.JobCompletionTime != nil {
		result["JobCompletionDate"] = job.JobCompletionTime.Format(timeutils.ISO8601SimpleFormat)
	}
	if job.Error != "" {
		result["Error"] = job.Error
	}

	entityList := make([]map[string]interface{}, 0)
	entityServices := make(map[string][]iamstore.TrackedActionLastAccessed)

	for _, action := range job.TrackedActionsLastAccessed {
		entityServices[action.EntityPath] = append(entityServices[action.EntityPath], action)
	}

	for entityPath, actions := range entityServices {
		entity := map[string]interface{}{
			"EntityPath": entityPath,
		}
		entityPolicyNames := make([]map[string]interface{}, 0)
		entityPolicyNames = append(entityPolicyNames, map[string]interface{}{
			"EntityName":        resolveEntityName(entityPath),
			"EntityType":        "User",
			"LastAuthenticated": job.JobCreationTime.Format(timeutils.ISO8601SimpleFormat),
		})
		entity["EntityPolicyList"] = entityPolicyNames

		svcMap := make(map[string][]iamstore.TrackedActionLastAccessed)
		for _, a := range actions {
			svcMap[a.ServiceNamespace] = append(svcMap[a.ServiceNamespace], a)
		}

		serviceList := make([]map[string]interface{}, 0, len(svcMap))
		for ns, svcActions := range svcMap {
			svcEntry := map[string]interface{}{
				"ServiceNamespace": ns,
			}
			actionList := make([]map[string]interface{}, 0, len(svcActions))
			for _, a := range svcActions {
				aEntry := map[string]interface{}{
					"ActionName": a.ActionName,
				}
				if a.LastAccessedDate != nil {
					aEntry["LastAccessedTime"] = a.LastAccessedDate.Format(timeutils.ISO8601SimpleFormat)
				}
				if a.LastAccessedRegion != "" {
					aEntry["LastAccessedRegion"] = a.LastAccessedRegion
				}
				actionList = append(actionList, aEntry)
			}
			svcEntry["TrackedActionsLastAccessed"] = actionList
			serviceList = append(serviceList, svcEntry)
		}
		entity["ServicesLastAccessed"] = serviceList
		entityList = append(entityList, entity)
	}

	result["EntityDetailsList"] = entityList

	return result, nil
}
