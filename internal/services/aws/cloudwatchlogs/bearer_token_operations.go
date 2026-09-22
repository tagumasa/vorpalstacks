package cloudwatchlogs

import (
	"context"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The modelled bearer token authentication switch: PutBearerToken
// Authentication enables or disables ACWL bearer tokens (IAM
// service-specific credentials of service logs.amazonaws.com) for one log
// group. While the switch is on, the group's HTTP ingestion endpoints
// accept bearer-token requests ("When enabled on a log group, bearer token
// authentication is enabled on operations until it is explicitly
// disabled", operation documentation); the per-group record also gates the
// endpoints' documented rule that bearer token authentication must be
// enabled on the log group before it can accept bearer-authenticated
// logs.

type PutBearerTokenAuthenticationInput struct {
	LogGroupIdentifier string
	// BearerTokenAuthenticationEnabled is the modelled setting; the wire
	// presence flag carries the member's required trait (an explicit false
	// is the documented disable operation, so absence must travel
	// separately).
	BearerTokenAuthenticationEnabled bool
	BearerTokenAuthenticationSet     bool
	Region                           string
}

// putBearerTokenAuthenticationCore validates and stores the group's
// bearer token authentication switch. Both members are required: the
// identifier (1..2048, the shared LogGroupIdentifier length) and the
// boolean, whose explicit false is a legitimate value.
func (s *LogsService) putBearerTokenAuthenticationCore(input PutBearerTokenAuthenticationInput) error {
	if !input.BearerTokenAuthenticationSet {
		return NewLogsError("InvalidParameterException",
			"The bearerTokenAuthenticationEnabled member is required", 400)
	}
	identifier := input.LogGroupIdentifier
	if identifier == "" {
		return errRequiredMember("logGroupIdentifier")
	}
	// The identifier member targets the shared LogGroupIdentifier shape:
	// the per-element length and alphabet traits, not the length alone.
	if err := validateLogGroupIdentifierElements([]string{identifier}); err != nil {
		return err
	}

	store, err := s.getLogsStoreByRegion(input.Region)
	if err != nil {
		return err
	}

	target := resolveLogGroupIdentifier(identifier)
	// The switch rides the group record itself (the read surfaces — the
	// DescribeLogGroups member and the ingestion gate — observe it with
	// the record they already hold). The mutate seam's group re-read under
	// the family lock keeps the teardown race closed: a delete interleaving
	// this put fails the mutate instead of leaving an orphaned switch.
	if err := store.MutateLogGroup(target, func(lg *logsstore.LogGroup) error {
		lg.BearerTokenAuthenticationEnabled = input.BearerTokenAuthenticationEnabled
		return nil
	}); err != nil {
		return mapStoreError(err)
	}
	return nil
}

// --- HTTP handler ---

func (s *LogsService) PutBearerTokenAuthentication(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	_, enabledSet := req.Parameters["bearerTokenAuthenticationEnabled"]
	if !enabledSet {
		_, enabledSet = req.Parameters["BearerTokenAuthenticationEnabled"]
	}
	input := PutBearerTokenAuthenticationInput{
		LogGroupIdentifier:               request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"),
		BearerTokenAuthenticationEnabled: request.GetBoolParam(req.Parameters, "BearerTokenAuthenticationEnabled"),
		BearerTokenAuthenticationSet:     enabledSet,
		Region:                           reqCtx.GetRegion(),
	}

	if err := s.putBearerTokenAuthenticationCore(input); err != nil {
		return nil, err
	}

	return response.EmptyResponse(), nil
}
