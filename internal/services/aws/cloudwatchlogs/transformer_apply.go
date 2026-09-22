package cloudwatchlogs

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"vorpalstacks/internal/core/logs"
	logsstore "vorpalstacks/internal/store/aws/cloudwatchlogs"
)

// The transformer application seam. The effective transformer of an
// ingested batch is the group-level transformer, else the account-level
// TRANSFORMER_POLICY whose LogGroupNamePrefix is the group's longest
// matching prefix ("If a log group has a log group-level transformer,
// that transformer overrides any account-level transformer that would
// otherwise apply to that log group"). The transformed output is
// persisted per event (the query plane's surface — GetLogEvents and
// FilterLogEvents keep serving originals) and handed to the metric and
// subscription filters whose applyOnTransformedLogs member selects the
// transformed form.

// transformerCacheTTL bounds how long an ingested group's resolved
// recipe is memoised: policy edits take effect on the next batch past
// the TTL instead of consulting the policy family per PutLogEvents.
const transformerCacheTTL = 15 * time.Second

// The transformer levels the metrics dimension vocabulary carries:
// LogGroupname "is used only for log-group-level transformers" and
// PolicyLevel "is used only for account-level transformers. Currently
// the only valid value for this dimension is AccountPolicy".
const (
	transformerLevelGroup   = "group"
	transformerLevelAccount = "account"
)

// transformerRecipe is a group's resolved effective transformer: the
// processor list and which form (group-level record or account-level
// policy) it came from.
type transformerRecipe struct {
	config []map[string]interface{}
	level  string
}

type transformerCacheEntry struct {
	recipe transformerRecipe
	at     time.Time
}

// transformerMu guards transformerCache.
var transformerMu sync.Mutex

// transformerCache memoises (group → effective recipe) per service.
var transformerCache = map[string]transformerCacheEntry{}

// effectiveTransformer resolves a group's effective processor list:
// the group-level record, else the account policy's. A nil config
// means no transformation applies.
func (s *LogsService) effectiveTransformer(store *logsstore.Store, region, groupName string) transformerRecipe {
	transformerMu.Lock()
	cached, ok := transformerCache[region+"/"+groupName]
	transformerMu.Unlock()
	if ok && time.Since(cached.at) < transformerCacheTTL {
		return cached.recipe
	}

	recipe := transformerRecipe{}
	if t, err := store.GetTransformer(groupName); err == nil {
		recipe = transformerRecipe{config: t.Config, level: transformerLevelGroup}
	} else if policy := s.matchingTransformerPolicy(store, groupName); policy != nil {
		recipe = transformerRecipe{config: policy, level: transformerLevelAccount}
	}
	transformerMu.Lock()
	transformerCache[region+"/"+groupName] = transformerCacheEntry{recipe: recipe, at: time.Now()}
	transformerMu.Unlock()
	return recipe
}

// matchingTransformerPolicy resolves the account-level transformer for
// a group: the TRANSFORMER_POLICY whose prefix is the group's longest
// matching prefix (the documented no-overlapping-prefixes rule makes
// longest-prefix unambiguous in the valid configuration space). "Log
// transformation and enrichment is supported only for log groups in the
// Standard log class" — a non-Standard group matches no account policy,
// mirroring the group-level form's Put-time class check.
func (s *LogsService) matchingTransformerPolicy(store *logsstore.Store, groupName string) []map[string]interface{} {
	if group, err := store.GetLogGroup(groupName); err == nil {
		if group.LogGroupClass != "" && group.LogGroupClass != "STANDARD" {
			return nil
		}
	}
	policies, err := store.ListAccountPolicies("TRANSFORMER_POLICY", "")
	if err != nil {
		return nil
	}
	// The best-prefix length sentinel starts below every real prefix
	// length: an empty selectionCriteria is the catch-all (length 0) and
	// one-character prefixes are legal, so a sentinel that itself has
	// length 1 would bar both from ever being selected.
	bestLen := -1
	var best []map[string]interface{}
	for _, policy := range policies {
		prefix, ok := transformerPolicyPrefix(policy.SelectionCriteria)
		if !ok {
			continue
		}
		if prefix != "" && !strings.HasPrefix(groupName, prefix) {
			continue
		}
		if len(prefix) > bestLen {
			config := transformerPolicyConfig(policy.PolicyDocument)
			if config == nil {
				continue
			}
			if err := validateTransformerConfig(config); err != nil {
				logs.Warn("Account transformer policy carries an invalid recipe; skipping",
					logs.String("policyName", policy.PolicyName), logs.Err(err))
				continue
			}
			bestLen = len(prefix)
			best = config
		}
	}
	return best
}

// invalidateTransformerCache drops the memoised recipes (the transformer
// CRUD paths call it so an edit takes effect on the next batch rather
// than after the TTL).
func invalidateTransformerCache(region, groupName string) {
	transformerMu.Lock()
	delete(transformerCache, region+"/"+groupName)
	transformerMu.Unlock()
}

// transformIngestedBatch applies the group's effective recipe to a
// written batch and persists the transformed messages. The write has
// already committed when this runs: a transformation failure never
// fails ingestion ("Any log events greater than 512kb will fail in
// transformation and emit an error" — the original still ingests, and
// the failure persists its @transformationError marker copy).
func (s *LogsService) transformIngestedBatch(store *logsstore.Store, region, groupName, logStream string, events []logsstore.LogEntry) map[string]string {
	recipe := s.effectiveTransformer(store, region, groupName)
	if recipe.config == nil {
		return nil
	}
	tctx := TransformContext{
		LogGroup:  groupName,
		LogStream: logStream,
		AccountID: s.accountID,
		Region:    region,
	}
	entries := make(map[string]logsstore.TransformedMessage, len(events))
	transformed := make(map[string]string, len(events))
	transformedCount, transformedBytes, errorCount := 0, 0, 0
	for _, e := range events {
		out, failed := transformEvent(recipe.config, e.Message, tctx)
		if out == "" {
			continue
		}
		digest := logsstore.TransformedMessageDigest(e.Timestamp, logStream, e.Message)
		entries[digest] = logsstore.TransformedMessage{Timestamp: e.Timestamp, Message: out}
		transformed[digest] = out
		if failed {
			errorCount++
			continue
		}
		transformedCount++
		transformedBytes += len(out)
	}
	if len(entries) == 0 {
		return nil
	}
	if err := store.PutTransformedMessages(groupName, entries); err != nil {
		logs.Error("Failed to persist transformed messages",
			logs.String("logGroup", groupName), logs.Err(err))
	}
	s.emitTransformationMetrics(region, groupName, recipe.level, transformedCount, transformedBytes, errorCount)
	return transformed
}

// transformedForFilters picks the message form a filter evaluates: the
// transformed form when the filter's applyOnTransformedLogs member is
// set and a transformed copy exists, else the original.
func transformedForFilters(original string, digest string, transformed map[string]string, applyOnTransformed bool) string {
	if !applyOnTransformed {
		return original
	}
	if out, ok := transformed[digest]; ok {
		return out
	}
	return original
}

// invalidateAllTransformerCaches drops every memoised recipe (an
// account-policy edit can change any prefix's recipe).
func invalidateAllTransformerCaches() {
	transformerMu.Lock()
	transformerCache = map[string]transformerCacheEntry{}
	transformerMu.Unlock()
}

// validateTransformerPolicySelection enforces the TRANSFORMER_POLICY
// prefix rules and the account-level cap: the criteria is a
// LogGroupNamePrefix (or empty for all standard groups), "You can't
// create two transformer policies in the same Region that use the same
// prefix or have one prefix contained within another", and "you can
// create as many as 20 account-level transformer policies" per Region.
func (s *LogsService) validateTransformerPolicySelection(region, policyName, selectionCriteria string) error {
	prefix, ok := transformerPolicyPrefix(selectionCriteria)
	if !ok {
		return NewLogsError("InvalidParameterException",
			"The only supported selectionCriteria for a TRANSFORMER_POLICY is LogGroupNamePrefix", 400)
	}
	store, err := s.getLogsStoreByRegion(region)
	if err != nil {
		return err
	}
	policies, err := store.ListAccountPolicies("TRANSFORMER_POLICY", "")
	if err != nil {
		return mapStoreError(err)
	}
	others := 0
	for _, policy := range policies {
		if policy.PolicyName == policyName {
			continue
		}
		others++
		existingPrefix, ok := transformerPolicyPrefix(policy.SelectionCriteria)
		if !ok {
			continue
		}
		// "You can have one account-level transformer policy that
		// applies to all log groups in the account. Or you can create
		// as many as 20 account-level transformer policies that are
		// each scoped to a subset of log groups" — the account-wide
		// policy and prefix-scoped policies never coexist, in either
		// creation order.
		if prefix == "" {
			return NewLogsError("InvalidParameterException",
				"An account-wide transformer policy cannot be created while other transformer policies exist", 400)
		}
		if existingPrefix == "" {
			return NewLogsError("InvalidParameterException",
				"A prefix-scoped transformer policy cannot be created while an account-wide transformer policy exists", 400)
		}
		if prefix == existingPrefix ||
			strings.HasPrefix(existingPrefix, prefix) ||
			strings.HasPrefix(prefix, existingPrefix) {
			return NewLogsError("InvalidParameterException",
				"A transformer policy for prefix "+existingPrefix+" already exists with an overlapping prefix", 400)
		}
	}
	if others >= maxAccountTransformerPolicies {
		return NewLogsError("LimitExceededException",
			fmt.Sprintf("You can create as many as %d account-level transformer policies in a Region", maxAccountTransformerPolicies), 400)
	}
	return nil
}

// maxAccountTransformerPolicies is the documented per-Region ceiling on
// account-level transformer policies.
const maxAccountTransformerPolicies = 20
