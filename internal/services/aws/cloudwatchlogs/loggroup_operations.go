package cloudwatchlogs

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	tagutil "vorpalstacks/internal/common/tags"
)

// CreateLogGroup creates a new CloudWatch Logs log group.
func (s *LogsService) CreateLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")

	// The Tags member's ceiling, entry traits and reserved-prefix rules are
	// validated in the core both planes share.
	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	input := CreateLogGroupInput{
		LogGroupName:              logGroupName,
		KmsKeyId:                  request.GetParamLowerFirst(req.Parameters, "KmsKeyId"),
		LogGroupClass:             request.GetParamLowerFirst(req.Parameters, "LogGroupClass"),
		Tags:                      tags,
		DeletionProtectionEnabled: request.GetBoolParam(req.Parameters, "DeletionProtectionEnabled"),
		Region:                    reqCtx.GetRegion(),
	}

	if _, err := s.createLogGroupCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteLogGroup deletes a CloudWatch Logs log group.
func (s *LogsService) DeleteLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := DeleteLogGroupInput{
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		Region:       reqCtx.GetRegion(),
	}

	if err := s.deleteLogGroupCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DescribeLogGroups returns a list of CloudWatch Logs log groups.
func (s *LogsService) DescribeLogGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := ListLogGroupsInput{
		LogGroupNamePrefix:    request.GetParamLowerFirst(req.Parameters, "LogGroupNamePrefix"),
		LogGroupNamePattern:   request.GetParamLowerFirst(req.Parameters, "LogGroupNamePattern"),
		LogGroupIdentifiers:   request.GetStringList(req.Parameters, "LogGroupIdentifiers"),
		AccountIdentifiers:    request.GetStringList(req.Parameters, "AccountIdentifiers"),
		IncludeLinkedAccounts: request.GetBoolParam(req.Parameters, "IncludeLinkedAccounts"),
		LogGroupClass:         request.GetParamLowerFirst(req.Parameters, "LogGroupClass"),
		NextToken:             request.GetParamLowerFirst(req.Parameters, "NextToken"),
		Limit:                 int32(request.GetIntParam(req.Parameters, "Limit")),
		Region:                reqCtx.GetRegion(),
	}

	result, err := s.describeLogGroupsCore(input)
	if err != nil {
		return nil, err
	}

	logGroups := make([]map[string]interface{}, 0)
	for _, lg := range result.LogGroups {
		// "If you specify logGroupNamePattern in your request, then only
		// arn, creationTime, and logGroupName are included in the
		// response" (member documentation).
		if input.LogGroupNamePattern != "" {
			logGroups = append(logGroups, map[string]interface{}{
				"logGroupName": lg.Name,
				"arn":          lg.ARN + ":*",
				"creationTime": lg.CreatedAt.UnixMilli(),
			})
			continue
		}
		entry := map[string]interface{}{
			"logGroupName":      lg.Name,
			"arn":               lg.ARN + ":*",
			"creationTime":      lg.CreatedAt.UnixMilli(),
			"metricFilterCount": lg.MetricFilterCount,
			"storedBytes":       lg.StoredBytes,
			"logGroupArn":       lg.ARN,
			"logGroupClass":     lg.LogGroupClass,
		}
		if lg.RetentionInDays > 0 {
			entry["retentionInDays"] = lg.RetentionInDays
		}
		if lg.KmsKeyId != "" {
			entry["kmsKeyId"] = lg.KmsKeyId
		}
		if lg.DeletionProtectionEnabled {
			entry["deletionProtectionEnabled"] = lg.DeletionProtectionEnabled
		}
		// dataProtectionStatus reports the policy lifecycle ("Displays
		// whether this log group has a protection policy, or whether it
		// had one in the past") — an untouched group carries no member.
		if lg.DataProtectionStatus != "" {
			entry["dataProtectionStatus"] = lg.DataProtectionStatus
		}
		// The bearer switch follows the file's conditional-emission
		// pattern for configuration booleans: it appears once enabled,
		// matching the deletion-protection member beside it.
		if lg.BearerTokenAuthenticationEnabled {
			entry["bearerTokenAuthenticationEnabled"] = lg.BearerTokenAuthenticationEnabled
		}
		logGroups = append(logGroups, entry)
	}

	resp := map[string]interface{}{
		"logGroups": logGroups,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// ListLogGroups returns a list of CloudWatch Logs log groups.
func (s *LogsService) ListLogGroups(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	pattern := request.GetParamLowerFirst(req.Parameters, "logGroupNamePattern")
	if pattern == "" {
		pattern = request.GetParamLowerFirst(req.Parameters, "LogGroupNamePattern")
	}
	legacyPrefix := request.GetParamLowerFirst(req.Parameters, "LogGroupNamePrefix")

	nextToken := request.GetParamLowerFirst(req.Parameters, "NextToken")
	limit := int32(request.GetIntParam(req.Parameters, "Limit"))

	// The structured filter members parse through the shared list-entry
	// helpers (DataSourceFilterInput is the type the aggregate-summaries
	// operation's own dataSources member already rides).
	var tagFilters []TagFilterInput
	for _, entry := range request.GetListParamLowerFirst(req.Parameters, "LogGroupTags") {
		tagFilters = append(tagFilters, TagFilterInput{
			Key:    request.GetStringParam(entry, "key"),
			Values: request.GetStringList(entry, "values"),
		})
	}
	var dataSources []DataSourceFilterInput
	for _, entry := range request.GetListParamLowerFirst(req.Parameters, "DataSources") {
		dataSources = append(dataSources, DataSourceFilterInput{
			Name: request.GetStringParam(entry, "name"),
			Type: request.GetStringParam(entry, "type"),
		})
	}

	input := ListLogGroupsInput{
		LogGroupClass:         request.GetParamLowerFirst(req.Parameters, "LogGroupClass"),
		AccountIdentifiers:    request.GetStringList(req.Parameters, "AccountIdentifiers"),
		IncludeLinkedAccounts: request.GetBoolParam(req.Parameters, "IncludeLinkedAccounts"),
		LogGroupTags:          tagFilters,
		DataSources:           dataSources,
		FieldIndexNames:       request.GetStringList(req.Parameters, "FieldIndexNames"),
		NextToken:             nextToken,
		Limit:                 limit,
		Region:                reqCtx.GetRegion(),
	}

	// The pattern member is up to five pipe-separated patterns of 3-24
	// restricted characters, each optionally anchored with ^: anchored
	// alternatives match by prefix, bare ones by substring ("For a
	// substring match, specify the string to match", ListLogGroups
	// logGroupNamePattern) — a single bare alternative included, so no
	// form of the member falls back to the store's prefix scan. A pattern
	// outside that form rejects.
	if pattern != "" {
		matcher, err := parseLogGroupNameMatcher(pattern)
		if err != nil {
			return nil, err
		}
		input.NameFilter = matcher.matches
		input.NameFilterPattern = pattern
	} else if legacyPrefix != "" {
		// LogGroupNamePrefix is not a member of the current model's
		// ListLogGroupsRequest; the fallback serves the legacy member's
		// own semantics alone — a plain prefix, never the pattern
		// member's substring rule.
		input.LogGroupNamePrefix = legacyPrefix
	}

	result, err := s.listLogGroupsCore(input)
	if err != nil {
		return nil, err
	}

	logGroups := make([]map[string]interface{}, 0)
	for _, lg := range result.LogGroups {
		// LogGroupSummary's members are exactly the three below
		// (logGroupArn, logGroupClass, logGroupName).
		entry := map[string]interface{}{
			"logGroupName":  lg.Name,
			"logGroupArn":   lg.ARN,
			"logGroupClass": lg.LogGroupClass,
		}
		logGroups = append(logGroups, entry)
	}

	resp := map[string]interface{}{
		"logGroups": logGroups,
	}
	if result.NextToken != "" {
		resp["nextToken"] = result.NextToken
	}

	return resp, nil
}

// logGroupNameMatcher is the parsed form of the ListLogGroups
// logGroupNamePattern member: a list of anchored-or-substring alternatives
// combined with OR. Every form filters through matches — the store's
// prefix scan serves the plain logGroupNamePrefix member alone, whose
// semantics a substring pattern must never inherit.
type logGroupNameMatcher struct {
	alternatives []logGroupNameAlternative
}

type logGroupNameAlternative struct {
	anchored bool
	text     string
}

// logGroupNameRegexPattern is the pattern trait on the
// LogGroupNameRegexPattern shape, verbatim: one to five pipe-separated
// patterns of 3-24 characters from the restricted set, each optionally
// preceded by the anchoring caret.
var logGroupNameRegexPattern = regexp.MustCompile(`^(\^?[\.\-_/#A-Za-z0-9]{3,24})(\|\^?[\.\-_/#A-Za-z0-9]{3,24}){0,4}$`)

// parseLogGroupNameMatcher validates the pattern against the shape's
// pattern trait and splits it into its alternatives.
func parseLogGroupNameMatcher(pattern string) (logGroupNameMatcher, error) {
	if !logGroupNameRegexPattern.MatchString(pattern) {
		return logGroupNameMatcher{}, NewLogsError("InvalidParameterException",
			fmt.Sprintf("Invalid logGroupNamePattern: %s. Up to five pipe-separated patterns of 3-24 characters from . - _ / # A-Z a-z 0-9, each optionally prefixed with ^", pattern), 400)
	}
	var m logGroupNameMatcher
	for _, alt := range strings.Split(pattern, "|") {
		a := logGroupNameAlternative{text: alt}
		if strings.HasPrefix(a.text, "^") {
			a.anchored = true
			a.text = a.text[1:]
		}
		m.alternatives = append(m.alternatives, a)
	}
	return m, nil
}

// matches reports whether the name matches any alternative: anchored
// alternatives compare by prefix, bare ones by substring.
func (m logGroupNameMatcher) matches(name string) bool {
	for _, a := range m.alternatives {
		if a.anchored {
			if strings.HasPrefix(name, a.text) {
				return true
			}
		} else if strings.Contains(name, a.text) {
			return true
		}
	}
	return false
}

// PutRetentionPolicy sets the retention policy for a CloudWatch Logs log group.
func (s *LogsService) PutRetentionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	logGroupName := request.GetParamLowerFirst(req.Parameters, "LogGroupName")

	retentionInDays := int32(request.GetIntParam(req.Parameters, "retentionInDays"))
	if retentionInDays == 0 {
		retentionInDays = int32(request.GetIntParam(req.Parameters, "RetentionInDays"))
	}

	input := PutRetentionPolicyInput{
		LogGroupName:    logGroupName,
		RetentionInDays: retentionInDays,
		Region:          reqCtx.GetRegion(),
	}

	if err := s.putRetentionPolicyCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// DeleteRetentionPolicy deletes the retention policy for a CloudWatch Logs log group.
func (s *LogsService) DeleteRetentionPolicy(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	input := DeleteRetentionPolicyInput{
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		Region:       reqCtx.GetRegion(),
	}

	if err := s.deleteRetentionPolicyCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// TagResource adds tags to a CloudWatch Logs log group.
func (s *LogsService) TagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleTag(ctx, req, s.tagHandlerConfig(store))
}

// UntagResource removes tags from a CloudWatch Logs resource.
func (s *LogsService) UntagResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleUntag(ctx, req, s.tagHandlerConfig(store))
}

// ListTagsForResource lists the tags for a CloudWatch Logs resource.
func (s *LogsService) ListTagsForResource(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	store, err := s.store(reqCtx)
	if err != nil {
		return nil, err
	}
	return tagutil.HandleList(ctx, req, s.tagHandlerConfig(store))
}

// TagLogGroup adds tags to the specified CloudWatch Logs log group.
func (s *LogsService) TagLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tags := tagutil.ToMap(tagutil.ParseTagsWithQueryFallback(req.Parameters, "Tags"))

	if err := s.tagLogGroupCore(TagLogGroupInput{
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		Tags:         tags,
		Region:       reqCtx.GetRegion(),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// UntagLogGroup removes tags from the specified CloudWatch Logs log group.
// Deprecated: use UntagResource instead.
func (s *LogsService) UntagLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.untagLogGroupCore(UntagLogGroupInput{
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		TagKeys:      request.GetStringList(req.Parameters, "Tags"),
		Region:       reqCtx.GetRegion(),
	}); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}

// ListTagsLogGroup retrieves the tags for the specified CloudWatch Logs log group.
func (s *LogsService) ListTagsLogGroup(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	tags, err := s.listTagsLogGroupCore(ListTagsLogGroupInput{
		LogGroupName: request.GetParamLowerFirst(req.Parameters, "LogGroupName"),
		Region:       reqCtx.GetRegion(),
	})
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"tags": tags,
	}, nil
}
