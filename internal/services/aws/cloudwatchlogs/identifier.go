package cloudwatchlogs

import (
	"strings"

	"vorpalstacks/internal/common/request"
	svcarn "vorpalstacks/internal/utils/aws/arn"
)

// The Smithy model documents several identifier members as name-or-ARN
// (GetLogEvents, DescribeLogStreams, FilterLogEvents, GetLogGroupFields,
// the data-protection-policy operations and PutLogGroupDeletionProtection
// all read "the name or ARN of the log group"; StartQuery and the
// scheduled-query operations accept the same form per element). Every such
// entry point routes through the resolvers below, so an ARN resolves to
// the store's name key through the platform ARN parser instead of failing
// the name lookup.

// resolveLogGroupIdentifier maps a name-or-ARN log group identifier to the
// store's name key. A bare name passes through unchanged; ARN input
// resolves through the ARN parser, tolerating the ":*" log-stream
// namespace suffix DescribeLogGroups appends to the object ARN. An ARN
// that does not resolve to a log-group resource passes through unchanged,
// so the subsequent store lookup reports it as not found.
func resolveLogGroupIdentifier(identifier string) string {
	if !strings.HasPrefix(identifier, "arn:") {
		return identifier
	}
	if name := svcarn.ExtractLogGroupNameFromARN(strings.TrimSuffix(identifier, ":*")); name != "" {
		return name
	}
	return identifier
}

// resolveLogGroupIdentifiers resolves a whole identifier list; the
// StartQuery and scheduled-query members document the same name-or-ARN
// form per element.
func resolveLogGroupIdentifiers(identifiers []string) []string {
	resolved := make([]string, len(identifiers))
	for i, id := range identifiers {
		resolved[i] = resolveLogGroupIdentifier(id)
	}
	return resolved
}

// resolveScheduledQueryIdentifier maps the ARN form of the scheduled
// query identifier to the store's id key; the bare name form resolves
// through resolveScheduledQueryRef's record scan (the members document
// "The ARN or name of the scheduled query").
func resolveScheduledQueryIdentifier(identifier string) string {
	if !strings.HasPrefix(identifier, "arn:") {
		return identifier
	}
	if id := svcarn.ExtractScheduledQueryIdFromARN(identifier); id != "" {
		return id
	}
	return identifier
}

// resolveLookupTableIdentifier maps a lookup-table ARN — the members
// document lookupTableArn as the ARN — to the table name; a bare name
// stays accepted as the platform's historical form.
func resolveLookupTableIdentifier(identifier string) string {
	if !strings.HasPrefix(identifier, "arn:") {
		return identifier
	}
	if name := svcarn.ExtractLookupTableNameFromARN(identifier); name != "" {
		return name
	}
	return identifier
}

// logGroupNameOrIdentifier reads the name-or-identifier member pair the
// read operations address log groups by. The pair is exclusive and
// exhaustive — "You must include either logGroupIdentifier or
// logGroupName, but not both" (member documentation) — so supplying both
// rejects and supplying neither rejects.
func logGroupNameOrIdentifier(params map[string]interface{}) (string, error) {
	name := request.GetParamLowerFirst(params, "LogGroupName")
	identifier := request.GetParamLowerFirst(params, "LogGroupIdentifier")
	if name != "" && identifier != "" {
		return "", NewLogsError("InvalidParameterException",
			"You must include either logGroupIdentifier or logGroupName, but not both", 400)
	}
	if name != "" {
		return name, nil
	}
	if identifier != "" {
		return resolveLogGroupIdentifier(identifier), nil
	}
	return "", errRequiredMember("logGroupIdentifier or logGroupName")
}

// logGroupTarget selects the log group a configuration core addresses: an
// explicit name wins, else the name-or-ARN identifier member resolves.
func logGroupTarget(name, identifier string) string {
	if name != "" {
		return name
	}
	return resolveLogGroupIdentifier(identifier)
}
