package iam

import (
	"context"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/core/logs"
	"vorpalstacks/internal/eventbus"
	iamstore "vorpalstacks/internal/store/aws/iam"
	awsarn "vorpalstacks/internal/utils/aws/arn"
	"vorpalstacks/internal/utils/timeutils"

	"github.com/google/uuid"
)

var (
	// ErrNoSuchJob is returned when a job ID does not match any stored report.
	ErrNoSuchJob = errors.NewAWSError("NoSuchEntity", "The job ID specified does not exist in this account.", http.StatusNotFound)
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

// generateJobID produces a unique UUID for a report job.
func generateJobID() string {
	return uuid.New().String()
}

// lookupAllEventsPageSize is the per-page MaxResults for CloudTrail
// LookupEvents. AWS spec (Smithy range trait) bounds MaxResults to 1-50.
const lookupAllEventsPageSize int32 = 50

// lookupAllEventsMaxPages caps the total number of pages fetched to avoid
// unbounded work in busy accounts. 20 pages × 50 events = 1000 events,
// which is ample for the ServiceLastAccessed report's 90-day window.
const lookupAllEventsMaxPages = 20

// lookupAllEvents paginates through CloudTrail LookupEvents results until
// either nextToken is empty or the page cap is reached. The caller passes
// a username filter ("" for unfiltered). The AWS LookupEvents spec bounds
// MaxResults to 1-50, so the per-page size is fixed at 50 by the helper.
func lookupAllEvents(ctx context.Context, invoker eventbus.CloudTrailInvoker, region, username string, startTime, endTime time.Time) ([]eventbus.CloudTrailEventInfo, error) {
	var all []eventbus.CloudTrailEventInfo
	nextToken := ""
	for page := 0; page < lookupAllEventsMaxPages; page++ {
		events, nt, err := invoker.LookupEvents(ctx, region, username, nextToken, startTime, endTime, lookupAllEventsPageSize)
		if err != nil {
			return nil, err
		}
		all = append(all, events...)
		if nt == "" {
			break
		}
		nextToken = nt
	}
	return all, nil
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

// parseIAMARNResource splits an IAM ARN into entity type and name.
// Uses arn.SplitARN for robust colon-delimited parsing instead of a
// local strings.Split.  Covers all IAM resource types defined in the
// AWS ARN reference; unknown types return empty strings so callers can
// fail-closed rather than silently defaulting to "User".
func parseIAMARNResource(arn string) (entityType, entityName string) {
	_, _, _, _, resource := awsarn.SplitARN(arn)
	if resource == "" {
		return "", ""
	}
	resourceParts := strings.SplitN(resource, "/", 2)
	switch resourceParts[0] {
	case "user":
		entityType = "User"
	case "role":
		entityType = "Role"
	case "group":
		entityType = "Group"
	case "instance-profile":
		entityType = "InstanceProfile"
	case "policy":
		entityType = "Policy"
	case "policy-version":
		entityType = "PolicyVersion"
	case "mfa":
		entityType = "MFADevice"
	case "server-certificate":
		entityType = "ServerCertificate"
	case "saml-provider":
		entityType = "SAMLProvider"
	case "oidc-provider":
		entityType = "OpenIDConnectProvider"
	case "access-key":
		entityType = "AccessKey"
	case "signing-certificate":
		entityType = "SigningCertificate"
	case "ssh-public-key":
		entityType = "SSHPublicKey"
	case "service-specific-credential":
		entityType = "ServiceSpecificCredential"
	default:
		return "", ""
	}
	if len(resourceParts) > 1 {
		entityName = resourceParts[1]
	} else {
		entityName = resource
	}
	return entityType, entityName
}

// resolveEntityName extracts the entity name (user, role, or group) from an IAM ARN.
func resolveEntityName(arn string) string {
	_, name := parseIAMARNResource(arn)
	return name
}

func resolveEntityType(arn string) string {
	entityType, _ := parseIAMARNResource(arn)
	return entityType
}

// generateLastAccessedReport queries CloudTrail events for the given entity within the
// specified time granularity and produces a ServiceLastAccessedJob with aggregated results.
func (s *IAMService) generateLastAccessedReport(arn, granularity, jobType, region string) *iamstore.ServiceLastAccessedJob {
	now := time.Now().UTC()
	duration := parseGranularity(granularity)
	startTime := now.Add(-duration)
	entityName := resolveEntityName(arn)

	var filteredEvents []eventbus.CloudTrailEventInfo
	if s.cloudTrailInvoker != nil {
		events, err := lookupAllEvents(context.Background(), s.cloudTrailInvoker, region, entityName, startTime, now)
		if err == nil {
			filteredEvents = events
		}
	}

	serviceMap := make(map[string]*iamstore.ServiceLastAccessed)
	actionMap := make(map[string]*iamstore.TrackedActionLastAccessed)

	// Build a per-service unique-entity count from an unfiltered event
	// query so that TotalAuthenticatedEntities reflects the actual number
	// of distinct entities that accessed each service (not just the
	// requested principal).
	entityCountMap := make(map[string]map[string]bool)
	if s.cloudTrailInvoker != nil {
		allEvents, err := lookupAllEvents(context.Background(), s.cloudTrailInvoker, region, "", startTime, now)
		if err == nil {
			for _, event := range allEvents {
				ns := eventSourceToServiceNamespace[event.EventSource]
				if ns == "" && event.EventSource != "" {
					parts := strings.SplitN(event.EventSource, ".", 2)
					ns = parts[0]
				}
				if entityCountMap[ns] == nil {
					entityCountMap[ns] = make(map[string]bool)
				}
				if event.Username != "" {
					entityCountMap[ns][event.Username] = true
				}
			}
		}
	}

	for _, event := range filteredEvents {
		eventRegion := region

		serviceNamespace := eventSourceToServiceNamespace[event.EventSource]
		if serviceNamespace == "" && event.EventSource != "" {
			parts := strings.SplitN(event.EventSource, ".", 2)
			serviceNamespace = parts[0]
		}

		serviceName := namespaceToDisplayName[serviceNamespace]
		if serviceName == "" {
			serviceName = serviceNamespace
		}

		svc, ok := serviceMap[serviceNamespace]
		if !ok {
			totalEntities := 1
			if users, ok := entityCountMap[serviceNamespace]; ok && len(users) > 0 {
				totalEntities = len(users)
			}
			svc = &iamstore.ServiceLastAccessed{
				ServiceName:                serviceName,
				ServiceNamespace:           serviceNamespace,
				TotalAuthenticatedEntities: totalEntities,
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

	now := time.Now().UTC()
	jobID := generateJobID()

	pendingJob := &iamstore.ServiceLastAccessedJob{
		JobID:           jobID,
		Arn:             arn,
		JobType:         "SERVICE_LAST_ACCESSED",
		JobStatus:       "IN_PROGRESS",
		JobCreationTime: now,
		Granularity:     granularity,
	}
	if err := store.ServiceLastAccessed().Put(pendingJob); err != nil {
		return nil, err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error("PANIC in ServiceLastAccessed report generation", logs.Any("panic", r))
			}
		}()
		completedJob := s.generateLastAccessedReport(arn, granularity, "SERVICE_LAST_ACCESSED", reqCtx.GetRegion())
		completedJob.JobID = jobID
		_ = store.ServiceLastAccessed().Put(completedJob)
	}()

	return map[string]interface{}{
		"JobId":     jobID,
		"JobType":   pendingJob.JobType,
		"JobStatus": "IN_PROGRESS",
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
		entityEntry := map[string]interface{}{
			"EntityName": resolveEntityName(entityPath),
			"EntityType": resolveEntityType(entityPath),
		}
		// Compute the entity-level LastAuthenticated as the most recent
		// action-level LastAccessedDate for this entity.
		var entityLastAccessed *time.Time
		for _, a := range actions {
			if a.LastAccessedDate != nil {
				if entityLastAccessed == nil || a.LastAccessedDate.After(*entityLastAccessed) {
					entityLastAccessed = a.LastAccessedDate
				}
			}
		}
		if entityLastAccessed != nil {
			entityEntry["LastAuthenticated"] = entityLastAccessed.Format(timeutils.ISO8601SimpleFormat)
		}
		entityPolicyNames = append(entityPolicyNames, entityEntry)
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
