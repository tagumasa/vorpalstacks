package iam

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"vorpalstacks/internal/common/invokers"
	"vorpalstacks/internal/common/pagination"
	"vorpalstacks/internal/core/logs"
	iamstore "vorpalstacks/internal/store/aws/iam"
	awsarn "vorpalstacks/internal/utils/aws/arn"
)

// Granularity values for the service last accessed report.  AWS defines
// Granularity as the level of detail: SERVICE_LEVEL generates service
// data only, ACTION_LEVEL generates service and action data.
const (
	granularityServiceLevel = "SERVICE_LEVEL"
	granularityActionLevel  = "ACTION_LEVEL"
)

// serviceLastAccessedWindow is the report lookback window.  AWS documents
// that IAM reports activity for at least the last 400 days; the
// platform's CloudTrail retention bounds the events actually available.
const serviceLastAccessedWindow = 400 * 24 * time.Hour

// GenerateServiceLastAccessedDetailsInput carries the parsed
// GenerateServiceLastAccessedDetails request.
type GenerateServiceLastAccessedDetailsInput struct {
	Arn         string
	Granularity string
}

// ServiceLastAccessedEntity is one EntityDetails entry of the
// GetServiceLastAccessedDetailsWithEntities result: the resolved
// EntityInfo members plus the entity-level last-accessed timestamp.
type ServiceLastAccessedEntity struct {
	Arn               string
	Name              string
	Type              string
	Id                string
	Path              string
	LastAuthenticated *time.Time
}

// ServiceLastAccessedEntitiesResult is the Core result for
// GetServiceLastAccessedDetailsWithEntities.
type ServiceLastAccessedEntitiesResult struct {
	Job         *iamstore.ServiceLastAccessedJob
	Entities    []ServiceLastAccessedEntity
	IsTruncated bool
	Marker      string
}

// generateServiceLastAccessedDetailsCore validates input, persists the
// IN_PROGRESS report job and spawns the background report generation,
// returning the persisted pending job.
func (s *IAMService) generateServiceLastAccessedDetailsCore(reqCtxRegion string, store *iamstore.IAMStore, input *GenerateServiceLastAccessedDetailsInput) (*iamstore.ServiceLastAccessedJob, error) {
	if input.Arn == "" {
		return nil, NewValidationError("Arn")
	}
	if err := validateARNParameter("Arn", input.Arn); err != nil {
		return nil, err
	}
	granularity := input.Granularity
	if granularity == "" {
		granularity = granularityServiceLevel
	}
	if granularity != granularityServiceLevel && granularity != granularityActionLevel {
		return nil, NewInvalidInputError("Granularity", "must be SERVICE_LEVEL or ACTION_LEVEL")
	}

	now := time.Now().UTC()
	jobID := generateJobID()

	pendingJob := &iamstore.ServiceLastAccessedJob{
		JobID:           jobID,
		Arn:             input.Arn,
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
		completedJob := s.generateLastAccessedReport(store, input.Arn, granularity, "SERVICE_LAST_ACCESSED", reqCtxRegion)
		completedJob.JobID = jobID
		_ = store.ServiceLastAccessed().Put(completedJob)
	}()

	return pendingJob, nil
}

// getServiceLastAccessedDetailsCore validates input and retrieves the
// stored report job.  A job id that matches no stored report fails with
// ErrNoSuchJob.
func (s *IAMService) getServiceLastAccessedDetailsCore(store *iamstore.IAMStore, jobID string) (*iamstore.ServiceLastAccessedJob, error) {
	// jobIDType is a fixed-length shape: any other length is malformed input
	// rather than a missing report.
	if len(jobID) != ServiceLastAccessedJobIDLength {
		return nil, NewInvalidInputError("JobId", fmt.Sprintf("must be exactly %d characters", ServiceLastAccessedJobIDLength))
	}
	job, err := store.ServiceLastAccessed().Get(jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, ErrNoSuchJob
	}
	return job, nil
}

// getServiceLastAccessedDetailsWithEntitiesCore validates input,
// retrieves the stored report job and projects the entity-level detail
// for the requested service namespace: the tracked actions of that
// service are grouped by entity, each entity's EntityInfo is resolved
// from the store, and the list is sorted by most recent access first as
// documented.
func (s *IAMService) getServiceLastAccessedDetailsWithEntitiesCore(store *iamstore.IAMStore, jobID, serviceNamespace, marker string, maxItems int) (*ServiceLastAccessedEntitiesResult, error) {
	if serviceNamespace == "" {
		return nil, NewValidationError("ServiceNamespace")
	}
	if len(serviceNamespace) > 64 {
		return nil, NewInvalidInputError("ServiceNamespace", "must be 1-64 characters matching [\\w-]")
	}
	for _, r := range serviceNamespace {
		if !isServiceNamespaceRune(r) {
			return nil, NewInvalidInputError("ServiceNamespace", "must be 1-64 characters matching [\\w-]")
		}
	}

	job, err := s.getServiceLastAccessedDetailsCore(store, jobID)
	if err != nil {
		return nil, err
	}

	entitiesByPath := make(map[string]*ServiceLastAccessedEntity)
	for _, action := range job.TrackedActionsLastAccessed {
		if action.ServiceNamespace != serviceNamespace {
			continue
		}
		entity, ok := entitiesByPath[action.EntityPath]
		if !ok {
			entityType, entityName := parseIAMARNResource(action.EntityPath)
			id, path := resolveEntityInfo(store, action.EntityPath)
			entity = &ServiceLastAccessedEntity{
				Arn:  action.EntityPath,
				Name: entityName,
				Type: policyOwnerEntityTypeFromARN(entityType),
				Id:   id,
				Path: path,
			}
			entitiesByPath[action.EntityPath] = entity
		}
		if action.LastAccessedDate != nil {
			if entity.LastAuthenticated == nil || action.LastAccessedDate.After(*entity.LastAuthenticated) {
				entity.LastAuthenticated = action.LastAccessedDate
			}
		}
	}

	entities := make([]ServiceLastAccessedEntity, 0, len(entitiesByPath))
	for _, entity := range entitiesByPath {
		entities = append(entities, *entity)
	}
	// The documented default order is by date with the most recent access
	// listed first; entities without an access timestamp sort last.
	sort.SliceStable(entities, func(i, j int) bool {
		a, b := entities[i].LastAuthenticated, entities[j].LastAuthenticated
		if a == nil || b == nil {
			return b == nil && a != nil
		}
		return a.After(*b)
	})

	paged := pagination.PaginateSlice(entities, marker, maxItems, func(entity ServiceLastAccessedEntity) string {
		return entity.Arn
	})
	return &ServiceLastAccessedEntitiesResult{
		Job:         job,
		Entities:    paged.Items,
		IsTruncated: paged.IsTruncated,
		Marker:      paged.NextMarker,
	}, nil
}

// isServiceNamespaceRune reports whether r is accepted by the documented
// ServiceNamespace pattern [\w-]* (word characters and hyphen).
func isServiceNamespaceRune(r rune) bool {
	return r == '-' || r == '_' || (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// policyOwnerEntityTypeFromARN maps the parsed IAM resource type onto the
// documented EntityInfo Type enum (USER | ROLE | GROUP).
func policyOwnerEntityTypeFromARN(entityType string) string {
	switch entityType {
	case "User":
		return "USER"
	case "Role":
		return "ROLE"
	case "Group":
		return "GROUP"
	}
	return entityType
}

// resolveEntityInfo best-effort resolves the entity's stable Id and Path
// from the store; empty values are returned when the entity no longer
// exists.
func resolveEntityInfo(store *iamstore.IAMStore, arn string) (id, path string) {
	entityType, name := parseIAMARNResource(arn)
	switch entityType {
	case "User":
		if user, err := store.Users().Get(name); err == nil {
			return user.ID, user.Path
		}
	case "Role":
		if role, err := store.Roles().Get(name); err == nil {
			return role.ID, role.Path
		}
	case "Group":
		if group, err := store.Groups().Get(name); err == nil {
			return group.ID, group.Path
		}
	}
	return "", ""
}

// reportPrincipal is one entity whose CloudTrail activity feeds a service
// last accessed report.
type reportPrincipal struct {
	arn  string
	name string
}

// reportPrincipals resolves the entities whose activity the report covers.
// The documented semantics depend on the report resource: a user or role
// report covers the entity itself; a group report covers its member users;
// a policy report covers the users and roles the policy is attached to.
func reportPrincipals(store *iamstore.IAMStore, arn string) []reportPrincipal {
	entityType, entityName := parseIAMARNResource(arn)
	partition, _, _, accountID, _ := awsarn.SplitARN(arn)
	userARN := func(name string) string {
		return fmt.Sprintf("arn:%s:iam::%s:user/%s", partition, accountID, name)
	}
	roleARN := func(name string) string {
		return fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, accountID, name)
	}

	switch entityType {
	case "Group":
		if store == nil {
			return nil
		}
		members, err := store.UserGroups().ListUsersInGroup(entityName)
		if err != nil {
			return nil
		}
		principals := make([]reportPrincipal, 0, len(members))
		for _, member := range members {
			principals = append(principals, reportPrincipal{arn: userARN(member), name: member})
		}
		return principals
	case "Policy":
		if store == nil {
			return nil
		}
		var principals []reportPrincipal
		refs, err := store.AttachedPolicies().ListPrincipalsForPolicy(arn)
		if err != nil {
			return nil
		}
		for _, ref := range refs {
			switch ref.PrincipalType {
			case PrincipalTypeUser:
				principals = append(principals, reportPrincipal{arn: userARN(ref.PrincipalName), name: ref.PrincipalName})
			case PrincipalTypeRole:
				principals = append(principals, reportPrincipal{arn: roleARN(ref.PrincipalName), name: ref.PrincipalName})
			}
		}
		return principals
	default:
		// User and role reports cover the entity itself.
		return []reportPrincipal{{arn: arn, name: entityName}}
	}
}

// generateLastAccessedReport queries CloudTrail events for every entity
// the report covers within the report tracking window and produces a
// ServiceLastAccessedJob with aggregated results.  Each service entry
// records the ARN of the entity that produced the most recent access
// attempt, so group and policy reports surface the actual member activity
// the documented response describes.  Action-level tracking is populated
// only for ACTION_LEVEL granularity, matching the documented granularity
// semantics.
func (s *IAMService) generateLastAccessedReport(store *iamstore.IAMStore, arn, granularity, jobType, region string) *iamstore.ServiceLastAccessedJob {
	now := time.Now().UTC()
	startTime := now.Add(-serviceLastAccessedWindow)
	trackActions := granularity == granularityActionLevel

	type principalEvent struct {
		principal reportPrincipal
		event     invokers.CloudTrailEventInfo
	}
	var filteredEvents []principalEvent
	if s.cloudTrailInvoker != nil {
		for _, principal := range reportPrincipals(store, arn) {
			events, err := lookupAllEvents(context.Background(), s.cloudTrailInvoker, region, principal.name, startTime, now)
			if err != nil {
				continue
			}
			for _, event := range events {
				filteredEvents = append(filteredEvents, principalEvent{principal: principal, event: event})
			}
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

	for _, pe := range filteredEvents {
		event := pe.event
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
			svc.LastAuthenticatedEntity = pe.principal.arn
		}

		if !trackActions {
			continue
		}

		actionKey := serviceNamespace + ":" + event.EventName
		action, ok := actionMap[actionKey]
		if !ok {
			action = &iamstore.TrackedActionLastAccessed{
				ActionName:       event.EventName,
				ServiceNamespace: serviceNamespace,
				EntityPath:       pe.principal.arn,
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
	// The documented default ordering of the service list is by service
	// namespace.
	sort.SliceStable(services, func(i, j int) bool {
		return services[i].ServiceNamespace < services[j].ServiceNamespace
	})

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
