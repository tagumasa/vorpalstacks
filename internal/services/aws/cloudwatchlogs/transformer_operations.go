package cloudwatchlogs

import (
	"context"
	"fmt"
	"strings"

	"vorpalstacks/internal/common/request"
	"vorpalstacks/internal/common/response"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The log transformer family: per-group transformer CRUD and the
// TestTransformer dry run. The account-level TRANSFORMER_POLICY family
// is the same recipe engine's account form (the ingestion seam resolves
// the group's effective transformer: the group-level record, else the
// longest-prefix account policy).

// transformerMaxProcessors is the documented pipeline bound ("You can
// have as many as 20 processors in a transformer") — the model's
// Processors length trait 1-20.
const transformerMaxProcessors = 20

// transformerParserProcessors are the parser-type processors ("Each
// transformer must have at least one parser, and the first processor in
// a transformer must be a parser... as many as five parser-type
// processors").
var transformerParserProcessors = map[string]bool{
	"parseJSON": true, "grok": true, "csv": true, "parseKeyValue": true,
	"parseCloudfront": true, "parseRoute53": true, "parseVPC": true,
	"parseWAF": true, "parsePostgres": true,
}

// transformerUnsupportedProcessors are the built-in processors that
// reject at Put/Test with the recorded reason, never approximated.
// parseToOCSF's input contract is documented but the AWS-to-OCSF field
// mapping it applies is published nowhere (neither the parseToOCSF page
// nor the Smithy model carries an output example), so implementing it
// would invent AWS-internal specification. The unsupported rejection
// fires before any parser-set accounting, so a processor here never
// reaches the parser vocabulary below.
var transformerUnsupportedProcessors = map[string]string{
	"parseToOCSF": "The parseToOCSF processor's AWS-to-OCSF field mapping is not published in the AWS documentation or the service model",
}

// transformerSingletonProcessors may appear once in a pipeline ("You
// can include only one grok / addKeys / copyValue processor").
var transformerSingletonProcessors = map[string]bool{
	"grok": true, "addKeys": true, "copyValue": true,
}

// transformerVendedProcessors are the pre-configured parsers for AWS
// vended logs: "If you include a pre-configured parser for a type of
// AWS vended logs, it must be the first processor listed in the
// transformer. You can include only one such processor in a
// transformer."
var transformerVendedProcessors = map[string]bool{
	"parseCloudfront": true, "parseRoute53": true, "parseVPC": true,
	"parseWAF": true, "parsePostgres": true,
}

// validateTransformerConfig enforces the documented structural rules
// shared by PutTransformer, the account policy and TestTransformer.
func validateTransformerConfig(config []map[string]interface{}) error {
	if len(config) == 0 || len(config) > transformerMaxProcessors {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("A transformer must carry 1 to %d processors", transformerMaxProcessors), 400)
	}
	parsers := 0
	vended := 0
	kinds := map[string]int{}
	for i, processor := range config {
		kind, cfg, ok := processorKind(processor)
		if !ok {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Processor %d does not carry exactly one processor configuration", i+1), 400)
		}
		kinds[kind]++
		if reason, unsupported := transformerUnsupportedProcessors[kind]; unsupported {
			return NewLogsError("InvalidParameterException", reason, 400)
		}
		if !knownTransformerProcessor(kind) {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("Unknown processor %q", kind), 400)
		}
		if transformerParserProcessors[kind] {
			parsers++
		}
		if transformerVendedProcessors[kind] {
			vended++
			if i > 0 {
				return NewLogsError("InvalidParameterException",
					"A pre-configured parser for AWS vended logs must be the first processor in the transformer", 400)
			}
		}
		if i == 0 && !transformerParserProcessors[kind] {
			return NewLogsError("InvalidParameterException",
				"The first processor in a transformer must be a parser", 400)
		}
		if err := validateTransformerProcessorMembers(kind, cfg); err != nil {
			return err
		}
	}
	if parsers == 0 {
		return NewLogsError("InvalidParameterException",
			"A transformer must include at least one parser-type processor", 400)
	}
	if parsers > 5 {
		return NewLogsError("InvalidParameterException",
			"A transformer can include at most five parser-type processors", 400)
	}
	if vended > 1 {
		return NewLogsError("InvalidParameterException",
			"A transformer can include only one pre-configured parser for AWS vended logs", 400)
	}
	for kind := range transformerSingletonProcessors {
		if kinds[kind] > 1 {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("A transformer can include only one %s processor", kind), 400)
		}
	}
	return nil
}

// knownTransformerProcessor is the platform's implemented processor
// vocabulary: the configurable processors plus the five built-in
// vended-format decoders.
func knownTransformerProcessor(kind string) bool {
	switch kind {
	case "parseJSON", "grok", "csv", "parseKeyValue",
		"parseCloudfront", "parseRoute53", "parseVPC", "parseWAF", "parsePostgres",
		"addKeys", "deleteKeys", "moveKeys", "renameKeys", "copyValue", "listToMap",
		"lowerCaseString", "upperCaseString", "trimString", "splitString", "substituteString",
		"typeConverter", "dateTimeConverter":
		return true
	}
	return false
}

// The per-member bounds the model and the processor pages carry: path
// members are 1-128 characters with at most three nested key levels,
// plain string members carry their own lengths.
const (
	transformerMaxPathLength = 128
	transformerMaxPathDepth  = 3
	transformerMaxValueLen   = 256
	transformerMaxTargetFmt  = 64
	transformerMaxColumns    = 100
	transformerMaxPatterns   = 5
	// transformerMaxGrokMatchLen is the grok match member's documented
	// ceiling ("Maximum length: 512").
	transformerMaxGrokMatchLen = 512
	// transformerMaxDryRunMessages is the TestTransformer logEventMessages
	// member's documented ceiling (1 to 50 messages).
	transformerMaxDryRunMessages = 50
)

// validTransformerPath reports whether a dotted path member obeys its
// documented length and nested-key-depth bounds.
func validTransformerPath(path string) bool {
	if len(path) == 0 || len(path) > transformerMaxPathLength {
		return false
	}
	return len(strings.Split(path, ".")) <= transformerMaxPathDepth
}

// validTransformerMember reports whether a plain string member obeys
// its documented 1-128 length bound.
func validTransformerMember(value string) bool {
	return len(value) > 0 && len(value) <= transformerMaxPathLength
}

// validateTransformerProcessorMembers enforces the per-processor member
// rows the documented configuration tables carry (required members,
// entry ceilings, member lengths and depths); the per-kind validators
// live in transformer_validators.go.
func validateTransformerProcessorMembers(kind string, cfg map[string]interface{}) error {
	// The built-in vended-format parsers "accept only @message as the
	// input".
	switch kind {
	case "parseCloudfront", "parseRoute53", "parseVPC", "parseWAF", "parsePostgres":
		if source, has := cfg["source"].(string); has && source != "@message" {
			return NewLogsError("InvalidParameterException",
				fmt.Sprintf("The %s processor accepts only @message as the input", kind), 400)
		}
	}
	switch kind {
	case "deleteKeys":
		return validateWithKeysMembers(kind, cfg, 5)
	case "lowerCaseString", "upperCaseString", "trimString":
		return validateWithKeysMembers(kind, cfg, 10)
	case "grok":
		return validateGrokMembers(cfg)
	case "parseJSON":
		return validateParseJSONMembers(kind, cfg)
	case "csv":
		return validateCSVMembers(kind, cfg)
	case "parseKeyValue":
		return validateParseKeyValueMembers(kind, cfg)
	case "addKeys":
		return validateAddKeysMembers(cfg)
	case "moveKeys", "copyValue":
		return validateMoveKeysMembers(kind, cfg)
	case "renameKeys":
		return validateRenameKeysMembers(cfg)
	case "splitString":
		return validateSplitStringMembers(kind, cfg)
	case "substituteString":
		return validateSubstituteStringMembers(kind, cfg)
	case "typeConverter":
		return validateTypeConverterMembers(kind, cfg)
	case "dateTimeConverter":
		return validateDateTimeConverterMembers(kind, cfg)
	case "listToMap":
		return validateListToMapMembers(kind, cfg)
	}
	return nil
}

// requiredEntryString enforces a required entry member (the model marks
// delimiter, from and to @required): the member must be present, a string,
// and 1 to max characters.
func requiredEntryString(kind string, entry map[string]interface{}, member string, max int) error {
	value, has := entry[member].(string)
	if !has || len(value) == 0 || len(value) > max {
		return NewLogsError("InvalidParameterException",
			fmt.Sprintf("The %s processor's %s accepts 1 to %d characters", kind, member, max), 400)
	}
	return nil
}

// --- Cores ---

// putTransformerCore validates and stores one log group's transformer.
func (s *LogsService) putTransformerCore(identifier string, config []map[string]interface{}, region string) (*logsstore.Transformer, error) {
	if identifier == "" {
		return nil, errRequiredMember("logGroupIdentifier")
	}
	if err := validateLogGroupIdentifierElements([]string{identifier}); err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errRequiredMember("transformerConfig")
	}
	if err := validateTransformerConfig(config); err != nil {
		return nil, err
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	groupName := resolveLogGroupIdentifier(identifier)
	group, err := store.GetLogGroup(groupName)
	if err != nil {
		return nil, mapStoreError(err)
	}
	// "Log transformation and enrichment is supported only for log
	// groups in the Standard log class."
	if group.LogGroupClass != "" && group.LogGroupClass != "STANDARD" {
		return nil, NewLogsError("InvalidParameterException",
			"Log transformation is supported only for log groups in the Standard log class", 400)
	}
	transformer := &logsstore.Transformer{
		LogGroupName:       groupName,
		LogGroupIdentifier: store.ARNBuilder().CloudWatch().LogGroup(groupName),
		Config:             config,
	}
	if err := store.PutTransformer(transformer); err != nil {
		return nil, mapStoreError(err)
	}
	invalidateTransformerCache(region, groupName)
	return transformer, nil
}

func (s *LogsService) getTransformerCore(identifier, region string) (*logsstore.Transformer, error) {
	if identifier == "" {
		return nil, errRequiredMember("logGroupIdentifier")
	}
	if err := validateLogGroupIdentifierElements([]string{identifier}); err != nil {
		return nil, err
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return nil, err
	}
	transformer, err := store.GetTransformer(resolveLogGroupIdentifier(identifier))
	if err != nil {
		return nil, mapStoreError(err)
	}
	return transformer, nil
}

func (s *LogsService) deleteTransformerCore(identifier, region string) error {
	if identifier == "" {
		return errRequiredMember("logGroupIdentifier")
	}
	if err := validateLogGroupIdentifierElements([]string{identifier}); err != nil {
		return err
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	groupName := resolveLogGroupIdentifier(identifier)
	if err := store.DeleteTransformer(groupName); err != nil {
		return mapStoreError(err)
	}
	invalidateTransformerCache(region, groupName)
	return nil
}

// testTransformerCore runs the recipe over supplied messages (the
// TestTransformer dry run: 1 to 50 messages).
func (s *LogsService) testTransformerCore(config []map[string]interface{}, messages []string, tctx TransformContext) ([]map[string]interface{}, error) {
	if config == nil {
		return nil, errRequiredMember("transformerConfig")
	}
	if len(messages) == 0 || len(messages) > transformerMaxDryRunMessages {
		return nil, NewLogsError("InvalidParameterException",
			fmt.Sprintf("logEventMessages accepts 1 to %d messages", transformerMaxDryRunMessages), 400)
	}
	// Each logEventMessages member carries a minimum length of one —
	// an empty-string message is not a log event.
	for _, message := range messages {
		if message == "" {
			return nil, NewLogsError("InvalidParameterException",
				"Each entry in logEventMessages accepts at least one character", 400)
		}
	}
	if err := validateTransformerConfig(config); err != nil {
		return nil, err
	}
	// The response's TransformedLogRecord carries the event number, the
	// original message and the transformed message — a failed
	// transformation carries its @transformationError marker record.
	records := make([]map[string]interface{}, len(messages))
	for i, message := range messages {
		transformed, _ := transformEvent(config, message, tctx)
		if transformed == "" {
			transformed = message
		}
		records[i] = map[string]interface{}{
			"eventNumber":             int64(i + 1),
			"eventMessage":            message,
			"transformedEventMessage": transformed,
		}
	}
	return records, nil
}

// --- HTTP handlers ---

// transformerConfigFromParams reads the transformerConfig processor list.
// Every element must be the object form of one Processor structure — the
// member's type; a scalar or array element is a malformed request and
// rejects here rather than being dropped (a silently shrunken pipeline
// would persist a transformer the caller never sent).
func transformerConfigFromParams(params map[string]interface{}) ([]map[string]interface{}, error) {
	raw := request.GetArrayParamLowerFirst(params, "TransformerConfig")
	config := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]interface{})
		if !ok {
			return nil, NewLogsError("InvalidParameterException",
				"Every transformerConfig entry must be a processor object", 400)
		}
		config = append(config, m)
	}
	return config, nil
}

func (s *LogsService) PutTransformer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	config, err := transformerConfigFromParams(req.Parameters)
	if err != nil {
		return nil, err
	}
	if _, err := s.putTransformerCore(
		request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"),
		config,
		reqCtx.GetRegion()); err != nil {
		return nil, err
	}
	// The operation's response shape is Unit.
	return response.EmptyResponse(), nil
}

func (s *LogsService) GetTransformer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	transformer, err := s.getTransformerCore(
		request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"), reqCtx.GetRegion())
	if err != nil {
		return nil, err
	}
	// GetTransformerResponse is a flat shape: the identifier, the stamps
	// and the configuration are its members.
	return formatTransformer(transformer), nil
}

func (s *LogsService) DeleteTransformer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	if err := s.deleteTransformerCore(
		request.GetParamLowerFirst(req.Parameters, "LogGroupIdentifier"), reqCtx.GetRegion()); err != nil {
		return nil, err
	}
	return response.EmptyResponse(), nil
}

func (s *LogsService) TestTransformer(ctx context.Context, reqCtx *request.RequestContext, req *request.ParsedRequest) (interface{}, error) {
	messages := request.GetStringList(req.Parameters, "LogEventMessages")
	config, cfgErr := transformerConfigFromParams(req.Parameters)
	if cfgErr != nil {
		return nil, cfgErr
	}
	results, err := s.testTransformerCore(
		config,
		messages,
		TransformContext{AccountID: s.accountID, Region: reqCtx.GetRegion()})
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"transformedLogs": results,
	}, nil
}

// formatTransformer renders the model's Transformer shape (the response
// carries logGroupIdentifier, the stamps and the configuration).
func formatTransformer(t *logsstore.Transformer) map[string]interface{} {
	return map[string]interface{}{
		"logGroupIdentifier": t.LogGroupIdentifier,
		"creationTime":       t.CreationTime,
		"lastModifiedTime":   t.LastModifiedTime,
		"transformerConfig":  t.Config,
	}
}

// --- The account-policy form ---

// transformerPolicyConfig extracts the processor list from a
// TRANSFORMER_POLICY account policy document ("A transformer policy
// must include one JSON block with the array of processors and their
// configurations"): the document is the processor array itself. A nil
// return is a malformed document the Put path rejects.
func transformerPolicyConfig(document string) []map[string]interface{} {
	document = strings.TrimSpace(document)
	if document == "" {
		return nil
	}
	var asList []map[string]interface{}
	if err := jsonUnmarshalString(document, &asList); err == nil {
		return asList
	}
	return nil
}

// transformerPolicyPrefix parses a TRANSFORMER_POLICY's
// selectionCriteria ("the only supported selectionCriteria filter is
// LogGroupNamePrefix"). The criteria form the console writes is
// `LogGroupNamePrefix = <prefix>`; the parser tolerates the separator
// and quoting variants. An empty criteria selects every group.
func transformerPolicyPrefix(selectionCriteria string) (string, bool) {
	sc := strings.TrimSpace(selectionCriteria)
	if sc == "" {
		return "", true
	}
	key, value, found := strings.Cut(sc, "=")
	if !found {
		key, value, found = strings.Cut(sc, ":")
	}
	if !found {
		return "", false
	}
	if strings.TrimSpace(key) != "LogGroupNamePrefix" {
		return "", false
	}
	prefix := strings.TrimSpace(value)
	prefix = strings.Trim(prefix, `"'`)
	return prefix, true
}
