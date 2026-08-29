package iot

import (
	"sort"
	"time"

	iotstore "vorpalstacks/internal/store/aws/iot"
)

// ---------------------------------------------------------------------------
// Command Core (remote command execution metadata management).
// Command definitions store a payload template, mandatory parameters, and
// an optional preprocessor Lambda ARN. Command executions are created by
// the MQTT data-plane when a device triggers a command; the control-plane
// API provides CRUD for definitions and read/delete for execution records.
// Records live under "iot-command/<commandId>" and
// "iot-command-execution/<executionId>" keys.
// ---------------------------------------------------------------------------

// CreateCommandInput carries the fields for CreateCommand. Payload,
// Preprocessor and MandatoryParameters keep the raw wire values (nested
// structures).
type CreateCommandInput struct {
	CommandID           string
	Namespace           string
	DisplayName         string
	Description         string
	Payload             interface{}
	PayloadTemplate     string
	Preprocessor        interface{}
	MandatoryParameters interface{}
	RoleArn             string
}

// CreateCommandResult is the transport-agnostic result of CreateCommand.
type CreateCommandResult struct {
	CommandID  string
	CommandArn string
}

// UpdateCommandInput carries the fields for UpdateCommand. DeprecatedProvided
// distinguishes an explicitly supplied flag from an omitted one.
type UpdateCommandInput struct {
	CommandID          string
	DisplayName        string
	Description        string
	Deprecated         bool
	DeprecatedProvided bool
}

// UpdateCommandResult is the transport-agnostic result of UpdateCommand.
type UpdateCommandResult struct {
	CommandID     string
	DisplayName   interface{}
	Description   interface{}
	Deprecated    interface{}
	LastUpdatedAt interface{}
}

// ListCommandsInput holds the parameters for listing commands.
// ListCommandsInput holds the parameters for listing commands: the
// namespace filter, the command-parameter-name filter (matching commands
// that declare the parameter), and the sort order.
type ListCommandsInput struct {
	Namespace            string
	CommandParameterName string
	SortOrder            string
}

// ListCommandExecutionsInput holds the parameters for listing command
// executions: the targetArn / commandArn / status body filters plus the
// namespace, sort-order, and started/completed time filters.
type ListCommandExecutionsInput struct {
	TargetArn           string
	CommandArn          string
	Status              string
	Namespace           string
	SortOrder           string
	StartedTimeFilter   interface{}
	CompletedTimeFilter interface{}
}

// createCommandCore validates and persists a command definition, rejecting
// duplicate command IDs.
func (s *IoTService) createCommandCore(store iotstore.IotStoreInterface, in CreateCommandInput) (*CreateCommandResult, error) {
	if in.CommandID == "" {
		return nil, iotstore.ErrMissingParam
	}
	// Reject duplicate command IDs.
	cmdExists, err := store.GetGenericExists("iot-command/"+in.CommandID, &map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if cmdExists {
		return nil, iotstore.ErrCommandAlreadyExists
	}
	now := time.Now().UTC().Unix()
	commandArn := iotstore.BuildCommandARN(store.GetAccountID(), store.GetRegion(), in.CommandID)
	rec := map[string]interface{}{
		"commandId":           in.CommandID,
		"commandArn":          commandArn,
		"namespace":           in.Namespace,
		"displayName":         in.DisplayName,
		"description":         in.Description,
		"payload":             in.Payload,
		"payloadTemplate":     in.PayloadTemplate,
		"preprocessor":        in.Preprocessor,
		"mandatoryParameters": in.MandatoryParameters,
		"roleArn":             in.RoleArn,
		"createdAt":           now,
		"lastUpdatedAt":       now,
		"deprecated":          false,
		"pendingDeletion":     false,
	}
	if err := store.PutGeneric("iot-command/"+in.CommandID, rec); err != nil {
		return nil, err
	}
	return &CreateCommandResult{
		CommandID:  in.CommandID,
		CommandArn: commandArn,
	}, nil
}

// getCommandCore retrieves a command definition record.
func (s *IoTService) getCommandCore(store iotstore.IotStoreInterface, cmdID string) (map[string]interface{}, error) {
	if cmdID == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-command/"+cmdID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandNotFound
	}
	return rec, nil
}

// updateCommandCore applies the supplied fields to an existing command
// definition.
func (s *IoTService) updateCommandCore(store iotstore.IotStoreInterface, in UpdateCommandInput) (*UpdateCommandResult, error) {
	if in.CommandID == "" {
		return nil, iotstore.ErrMissingParam
	}
	key := "iot-command/" + in.CommandID
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists(key, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandNotFound
	}
	if in.DisplayName != "" {
		rec["displayName"] = in.DisplayName
	}
	if in.Description != "" {
		rec["description"] = in.Description
	}
	if in.DeprecatedProvided {
		rec["deprecated"] = in.Deprecated
	}
	rec["lastUpdatedAt"] = time.Now().UTC().Unix()
	if err := store.PutGeneric(key, rec); err != nil {
		return nil, err
	}
	return &UpdateCommandResult{
		CommandID:     in.CommandID,
		DisplayName:   rec["displayName"],
		Description:   rec["description"],
		Deprecated:    rec["deprecated"],
		LastUpdatedAt: rec["lastUpdatedAt"],
	}, nil
}

// deleteCommandCore removes a command definition record.
func (s *IoTService) deleteCommandCore(store iotstore.IotStoreInterface, cmdID string) error {
	if cmdID == "" {
		return iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("iot-command/"+cmdID, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrCommandNotFound
	}
	return store.DeleteGeneric("iot-command/" + cmdID)
}

// listCommandsCore lists command summaries with the namespace filter, the
// command-parameter-name filter, and the createdAt sort order.
func (s *IoTService) listCommandsCore(store iotstore.IotStoreInterface, in ListCommandsInput) ([]map[string]interface{}, error) {
	if err := validateCommandSortOrder(in.SortOrder); err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("iot-command/")
	if err != nil {
		return nil, err
	}
	commands := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if in.Namespace != "" {
			if ns, _ := rec["namespace"].(string); ns != in.Namespace {
				continue
			}
		}
		if in.CommandParameterName != "" && !commandDeclaresParameter(rec["mandatoryParameters"], in.CommandParameterName) {
			continue
		}
		commands = append(commands, map[string]interface{}{
			"commandId":       rec["commandId"],
			"commandArn":      rec["commandArn"],
			"displayName":     rec["displayName"],
			"deprecated":      rec["deprecated"],
			"createdAt":       rec["createdAt"],
			"lastUpdatedAt":   rec["lastUpdatedAt"],
			"pendingDeletion": rec["pendingDeletion"],
		})
	}
	// The API lists commands in descending creation order by default.
	sortOrder := in.SortOrder
	if sortOrder == "" {
		sortOrder = "DESCENDING"
	}
	sortCommandRecords(commands, sortOrder, "createdAt", "commandId")
	return commands, nil
}

// commandDeclaresParameter reports whether a stored mandatoryParameters
// list carries an entry with the given name.
func commandDeclaresParameter(raw interface{}, name string) bool {
	parameters, ok := raw.([]interface{})
	if !ok {
		return false
	}
	for _, p := range parameters {
		pMap, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		if entryName, _ := pMap["name"].(string); entryName == name {
			return true
		}
	}
	return false
}

// validateCommandSortOrder enforces the SortOrder enum wire values.
func validateCommandSortOrder(sortOrder string) error {
	switch sortOrder {
	case "", "ASCENDING", "DESCENDING":
		return nil
	default:
		return iotstore.ErrValidation
	}
}

// sortCommandRecords sorts command summaries or execution summaries by the
// given time key (with the record id key as the deterministic tie-breaker).
func sortCommandRecords(records []map[string]interface{}, sortOrder, timeKey, idKey string) {
	sort.SliceStable(records, func(i, j int) bool {
		iTime, iOK := recordInt64(records[i][timeKey])
		jTime, jOK := recordInt64(records[j][timeKey])
		if iOK && jOK && iTime != jTime {
			if sortOrder == "DESCENDING" {
				return iTime > jTime
			}
			return iTime < jTime
		}
		iID, _ := records[i][idKey].(string)
		jID, _ := records[j][idKey].(string)
		if sortOrder == "DESCENDING" {
			return iID > jID
		}
		return iID < jID
	})
}

// getCommandExecutionCore retrieves a command execution record. The model
// marks targetArn required alongside executionId; both must be present.
func (s *IoTService) getCommandExecutionCore(store iotstore.IotStoreInterface, execID, targetArn string) (map[string]interface{}, error) {
	if execID == "" || targetArn == "" {
		return nil, iotstore.ErrMissingParam
	}
	rec := map[string]interface{}{}
	exists, err := store.GetGenericExists("iot-command-execution/"+execID, &rec)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, iotstore.ErrCommandExecutionNotFound
	}
	return rec, nil
}

// deleteCommandExecutionCore removes a command execution record. The model
// marks targetArn required alongside executionId; both must be present.
func (s *IoTService) deleteCommandExecutionCore(store iotstore.IotStoreInterface, execID, targetArn string) error {
	if execID == "" || targetArn == "" {
		return iotstore.ErrMissingParam
	}
	exists, err := store.GetGenericExists("iot-command-execution/"+execID, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if !exists {
		return iotstore.ErrCommandExecutionNotFound
	}
	return store.DeleteGeneric("iot-command-execution/" + execID)
}

// listCommandExecutionsCore lists command execution summaries with the
// targetArn / commandArn / status / namespace body filters, the
// started/completed time filters, and the start/completion-time sort order.
func (s *IoTService) listCommandExecutionsCore(store iotstore.IotStoreInterface, in ListCommandExecutionsInput) ([]map[string]interface{}, error) {
	if err := validateCommandSortOrder(in.SortOrder); err != nil {
		return nil, err
	}
	// The API documents that only the started or the completed time
	// filter may be provided, and only the command or the target ARN;
	// providing both members of either pair generates an error.
	if in.StartedTimeFilter != nil && in.CompletedTimeFilter != nil {
		return nil, iotstore.ErrValidation
	}
	if in.TargetArn != "" && in.CommandArn != "" {
		return nil, iotstore.ErrValidation
	}
	startedAfter, startedBefore, err := parseTimeFilter(in.StartedTimeFilter)
	if err != nil {
		return nil, err
	}
	completedAfter, completedBefore, err := parseTimeFilter(in.CompletedTimeFilter)
	if err != nil {
		return nil, err
	}
	items, err := store.ListGeneric("iot-command-execution/")
	if err != nil {
		return nil, err
	}
	executions := make([]map[string]interface{}, 0, len(items))
	for _, rec := range items {
		if in.TargetArn != "" {
			if ta, _ := rec["targetArn"].(string); ta != in.TargetArn {
				continue
			}
		}
		if in.CommandArn != "" {
			if ca, _ := rec["commandArn"].(string); ca != in.CommandArn {
				continue
			}
		}
		if in.Status != "" {
			if st, _ := rec["status"].(string); st != in.Status {
				continue
			}
		}
		if in.Namespace != "" {
			if ns, _ := rec["namespace"].(string); ns != in.Namespace {
				continue
			}
		}
		if !recordTimeWithin(rec["startedAt"], startedAfter, startedBefore) {
			continue
		}
		if !recordTimeWithin(rec["completedAt"], completedAfter, completedBefore) {
			continue
		}
		executions = append(executions, map[string]interface{}{
			"executionId": rec["executionId"],
			"commandArn":  rec["commandArn"],
			"targetArn":   rec["targetArn"],
			"status":      rec["status"],
			"createdAt":   rec["createdAt"],
			"startedAt":   rec["startedAt"],
			"completedAt": rec["completedAt"],
		})
	}
	// The API lists executions in descending order by default, keyed on
	// the start or completion time the provided filter selects.
	sortOrder := in.SortOrder
	if sortOrder == "" {
		sortOrder = "DESCENDING"
	}
	sortKey := "startedAt"
	if in.CompletedTimeFilter != nil {
		sortKey = "completedAt"
	}
	sortCommandRecords(executions, sortOrder, sortKey, "executionId")
	return executions, nil
}

// commandTimeFilterLayouts lists the accepted date-time layouts in
// preference order: the documented yyyy-MM-dd'T'HH:mm format, then the
// RFC 3339 form.
var commandTimeFilterLayouts = []string{"2006-01-02T15:04", time.RFC3339}

// parseCommandTimeFilterText parses one TimeFilter date-time string with the
// accepted layouts.
func parseCommandTimeFilterText(text string) (time.Time, bool) {
	for _, layout := range commandTimeFilterLayouts {
		if parsed, err := time.Parse(layout, text); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

// parseTimeFilter parses a TimeFilter wire structure ({after, before} as
// date-time strings) into unix-second bounds. Zero bounds mean absent.
func parseTimeFilter(raw interface{}) (after, before int64, err error) {
	filterMap, ok := raw.(map[string]interface{})
	if !ok || raw == nil {
		return 0, 0, nil
	}
	parse := func(key string) (int64, error) {
		value, present := filterMap[key]
		if !present || value == nil {
			return 0, nil
		}
		text, ok := value.(string)
		if !ok {
			return 0, iotstore.ErrValidation
		}
		parsed, ok := parseCommandTimeFilterText(text)
		if !ok {
			return 0, iotstore.ErrValidation
		}
		return parsed.Unix(), nil
	}
	after, err = parse("after")
	if err != nil {
		return 0, 0, err
	}
	before, err = parse("before")
	if err != nil {
		return 0, 0, err
	}
	return after, before, nil
}

// recordTimeWithin reports whether a stored unix-second timestamp falls
// within the after/before bounds (zero bounds are absent).
func recordTimeWithin(value interface{}, after, before int64) bool {
	if after == 0 && before == 0 {
		return true
	}
	stored, ok := recordInt64(value)
	if !ok {
		return false
	}
	if after != 0 && stored < after {
		return false
	}
	if before != 0 && stored > before {
		return false
	}
	return true
}
