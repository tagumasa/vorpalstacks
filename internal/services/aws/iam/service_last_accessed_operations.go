package iam

import (
	"context"
	"net/http"
	"strings"
	"time"

	"vorpalstacks/internal/common/errors"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/common/request"
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
// which is ample for the ServiceLastAccessed report window.
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

// GenerateServiceLastAccessedDetails generates a report of the last time
// the specified IAM entity accessed each AWS service.  The report is
// generated asynchronously and stored for retrieval via
// GetServiceLastAccessedDetails; the documented response carries the
// JobId only.
func (s *IAMService) GenerateServiceLastAccessedDetails(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	input := &GenerateServiceLastAccessedDetailsInput{
		Arn:         request.GetStringParam(req.Parameters, "Arn"),
		Granularity: request.GetStringParam(req.Parameters, "Granularity"),
	}
	pendingJob, err := s.generateServiceLastAccessedDetailsCore(reqCtx.GetRegion(), store, input)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"JobId": pendingJob.JobID,
	}, nil
}

// GetServiceLastAccessedDetails retrieves a previously generated report containing the last time
// each AWS service was accessed by the specified entity.  Supports
// pagination via Marker and MaxItems over the service list.
func (s *IAMService) GetServiceLastAccessedDetails(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetStringParam(req.Parameters, "JobId")
	if jobID == "" {
		jobID = request.GetStringParam(req.Parameters, "jobId")
	}
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	job, err := s.getServiceLastAccessedDetailsCore(store, jobID)
	if err != nil {
		return nil, err
	}

	paged := pagination.PaginateSlice(job.ServicesLastAccessed, marker, maxItems, func(svc iamstore.ServiceLastAccessed) string {
		return svc.ServiceNamespace
	})

	result := map[string]interface{}{
		"JobCreationDate":      job.JobCreationTime.Format(timeutils.ISO8601SimpleFormat),
		"JobStatus":            job.JobStatus,
		"JobType":              job.Granularity,
		"ServicesLastAccessed": serialiseServicesLastAccessed(paged.Items),
		"IsTruncated":          paged.IsTruncated,
	}
	if paged.NextMarker != "" {
		result["Marker"] = paged.NextMarker
	}

	if job.JobCompletionTime != nil {
		result["JobCompletionDate"] = job.JobCompletionTime.Format(timeutils.ISO8601SimpleFormat)
	}
	if job.Error != "" {
		result["Error"] = map[string]interface{}{
			"Message": job.Error,
			"Code":    "InternalError",
		}
	}

	return result, nil
}

// serialiseServicesLastAccessed projects the ServiceLastAccessed entries
// onto the documented response members, including the ARN of the entity
// that last attempted access — the report entity for user/role reports and
// the accessing member for group/policy reports.
func serialiseServicesLastAccessed(services []iamstore.ServiceLastAccessed) []map[string]interface{} {
	entries := make([]map[string]interface{}, 0, len(services))
	for _, svc := range services {
		entry := map[string]interface{}{
			"ServiceName":                svc.ServiceName,
			"ServiceNamespace":           svc.ServiceNamespace,
			"TotalAuthenticatedEntities": svc.TotalAuthenticatedEntities,
		}
		if svc.LastAuthenticated != nil {
			entry["LastAuthenticated"] = svc.LastAuthenticated.Format(timeutils.ISO8601SimpleFormat)
			if svc.LastAuthenticatedEntity != "" {
				entry["LastAuthenticatedEntity"] = svc.LastAuthenticatedEntity
			}
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
		entries = append(entries, entry)
	}
	return entries
}

// GetServiceLastAccessedDetailsWithEntities retrieves a previously
// generated report including the entity-level detail for the requested
// service namespace.  The entity list is scoped to that service and
// sorted by most recent access first; supports pagination via Marker and
// MaxItems.
func (s *IAMService) GetServiceLastAccessedDetailsWithEntities(_ context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	jobID := request.GetStringParam(req.Parameters, "JobId")
	if jobID == "" {
		jobID = request.GetStringParam(req.Parameters, "jobId")
	}
	serviceNamespace := request.GetStringParam(req.Parameters, "ServiceNamespace")
	marker := request.GetStringParam(req.Parameters, "Marker")
	maxItems := pagination.GetMaxItems(req.Parameters, pagination.DefaultMaxItems)

	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}

	result, err := s.getServiceLastAccessedDetailsWithEntitiesCore(store, jobID, serviceNamespace, marker, maxItems)
	if err != nil {
		return nil, err
	}

	entityList := make([]map[string]interface{}, 0, len(result.Entities))
	for _, entity := range result.Entities {
		entityInfo := map[string]interface{}{
			"Arn":  entity.Arn,
			"Id":   entity.Id,
			"Name": entity.Name,
			"Type": entity.Type,
		}
		if entity.Path != "" {
			entityInfo["Path"] = entity.Path
		}
		entry := map[string]interface{}{
			"EntityInfo": entityInfo,
		}
		if entity.LastAuthenticated != nil {
			entry["LastAuthenticated"] = entity.LastAuthenticated.Format(timeutils.ISO8601SimpleFormat)
		}
		entityList = append(entityList, entry)
	}

	resp := map[string]interface{}{
		"JobCreationDate":   result.Job.JobCreationTime.Format(timeutils.ISO8601SimpleFormat),
		"JobStatus":         result.Job.JobStatus,
		"EntityDetailsList": entityList,
		"IsTruncated":       result.IsTruncated,
	}
	if result.Marker != "" {
		resp["Marker"] = result.Marker
	}
	if result.Job.JobCompletionTime != nil {
		resp["JobCompletionDate"] = result.Job.JobCompletionTime.Format(timeutils.ISO8601SimpleFormat)
	}
	if result.Job.Error != "" {
		resp["Error"] = map[string]interface{}{
			"Message": result.Job.Error,
			"Code":    "InternalError",
		}
	}

	return resp, nil
}
