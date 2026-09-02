package inspection

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	wafstore "vorpalstacks/internal/store/aws/waf"
	"vorpalstacks/internal/utils/aws/arn"
)

// Evaluator evaluates a WebACL's rules against a request.
type Evaluator struct {
	resolvers Resolvers
}

// NewEvaluator creates an evaluator with the given referenced-entity
// resolvers.
func NewEvaluator(resolvers Resolvers) *Evaluator {
	return &Evaluator{resolvers: resolvers}
}

// evalCtx carries the state of one WebACL evaluation: the request, the
// labels earlier matched rules added, and the result being built.
type evalCtx struct {
	e      *Evaluator
	req    *Request
	acl    *wafstore.WebACL
	result *Result
	labels map[string]bool
}

// statementOutcome classifies the evaluation of a statement tree. The
// forwarded-IP configurations document that a request whose configured
// header is absent does not have the rule applied at all — a state
// distinct from a non-match, because enclosing NOT, AND and OR
// statements must not turn it into a match.
type statementOutcome int

const (
	outcomeNoMatch statementOutcome = iota
	outcomeMatch
	outcomeNotApplied
)

// boolOutcome lifts a plain match boolean into an always-applied
// statement outcome.
func boolOutcome(matched bool) statementOutcome {
	if matched {
		return outcomeMatch
	}
	return outcomeNoMatch
}

// Evaluate runs the web ACL's rules against the request in priority
// order (lowest priority number first, per the AWS WAF rule processing
// order) and returns the terminating action plus the match record.
func (e *Evaluator) Evaluate(acl *wafstore.WebACL, req *Request) *Result {
	ctx := &evalCtx{
		e:      e,
		req:    req,
		acl:    acl,
		result: &Result{},
		labels: map[string]bool{},
	}

	rules := make([]*wafstore.Rule, len(acl.Rules))
	copy(rules, acl.Rules)
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].Priority < rules[j].Priority
	})

	for _, rule := range rules {
		if rule == nil || rule.Statement == nil {
			continue
		}
		if rule.Statement.RuleGroupReferenceStatement != nil {
			terminated := ctx.evaluateRuleGroupRule(rule)
			if terminated {
				return ctx.result
			}
			continue
		}
		if rule.Statement.ManagedRuleGroupStatement != nil {
			terminated := ctx.evaluateManagedRuleGroupRule(rule)
			if terminated {
				return ctx.result
			}
			continue
		}
		if ctx.evaluateStatement(rule.Statement, rule.Name, labelNamespaceOf(acl.LabelNamespace, acl.ARN, acl.Name)) != outcomeMatch {
			// A non-matching statement and one that leaves the rule
			// unapplied (forwarded-IP header absent) both skip the
			// rule entirely — no labels, no action — and evaluation
			// moves to the next rule.
			continue
		}
		ctx.recordLabels(rule, labelNamespaceOf(acl.LabelNamespace, acl.ARN, acl.Name))
		kind := actionKind(ruleAction(rule))
		switch kind {
		case ActionAllow, ActionBlock:
			ctx.recordMatch(rule, "", kind)
			ctx.result.Action = kind
			ctx.applyTerminatingAction(rule, ruleAction(rule))
			return ctx.result
		case ActionCount:
			ctx.recordMatch(rule, "", ActionCount)
			applyInsertHeaders(ctx.result, ruleAction(rule))
		case "":
			ctx.recordMatch(rule, "", ActionCount)
		case ActionCaptcha, ActionChallenge, ActionMonetize:
			if _, terminated := ctx.applyChallengeLikeAction(rule, "", kind, ruleAction(rule)); terminated {
				return ctx.result
			}
		}
	}

	defaultAction := webACLDefaultAction(acl)
	ctx.result.Action = ActionAllow
	if defaultAction != nil && defaultAction.Block != nil {
		ctx.result.Action = ActionBlock
		ctx.applyTerminatingAction(nil, defaultAction)
	}
	return ctx.result
}

// evaluateRuleGroupRule evaluates one web ACL rule whose statement is
// a rule group reference. The inner rules run in priority order; a
// terminating inner action (Allow or Block) ends the whole evaluation
// with that action unless the reference's OverrideAction counts it.
// ExcludedRules converts an inner action to Count, and a rule action
// override replaces the inner rule's action entirely. The return value
// reports whether the whole evaluation terminated.
func (ctx *evalCtx) evaluateRuleGroupRule(rule *wafstore.Rule) bool {
	stmt := rule.Statement.RuleGroupReferenceStatement
	group, err := ctx.e.resolvers.RuleGroup(stmt.ARN)
	if err != nil || group == nil {
		return false
	}
	groupNamespace := labelNamespaceOf(group.LabelNamespace, group.ARN, group.Name)
	override := ruleOverrideAction(rule)
	overrideCounts := override != nil && override.Count != nil
	excluded := make(map[string]bool, len(stmt.ExcludedRules))
	for _, excludedRule := range stmt.ExcludedRules {
		excluded[excludedRule.Name] = true
	}

	innerRules := make([]*wafstore.Rule, len(group.Rules))
	copy(innerRules, group.Rules)
	sort.SliceStable(innerRules, func(i, j int) bool {
		return innerRules[i].Priority < innerRules[j].Priority
	})

	matched := false
	for _, inner := range innerRules {
		if inner == nil || inner.Statement == nil {
			continue
		}
		if ctx.evaluateStatement(inner.Statement, inner.Name, groupNamespace) != outcomeMatch {
			continue
		}
		matched = true
		ctx.recordLabels(inner, groupNamespace)
		if inner.Statement.RuleGroupReferenceStatement != nil {
			if ctx.evaluateRuleGroupRule(inner) {
				return true
			}
			continue
		}
		action := ruleAction(inner)
		if excluded[inner.Name] {
			// Excluding a rule converts its action to Count.
			action = &wafstore.Action{Count: &wafstore.CountAction{}}
		}
		configured := actionKind(ruleAction(inner))
		overridden := ""
		if replacement := ruleActionOverrideFor(stmt.RuleActionOverrides, inner.Name); replacement != nil {
			action = replacement
			overridden = configured
		}
		kind := actionKind(action)
		switch kind {
		case ActionAllow, ActionBlock:
			if overrideCounts {
				ctx.recordMatch(inner, group.Name, ActionCount)
				continue
			}
			entry := ctx.recordMatch(inner, group.Name, kind)
			entry.OverriddenAction = overridden
			ctx.result.Action = kind
			// With OverrideAction None the inner action's own custom
			// response and header insertion apply directly.
			ctx.applyTerminatingAction(inner, action)
			return true
		case ActionCount:
			entry := ctx.recordMatch(inner, group.Name, ActionCount)
			entry.OverriddenAction = overridden
		case "":
			entry := ctx.recordMatch(inner, group.Name, ActionCount)
			entry.OverriddenAction = overridden
		case ActionCaptcha, ActionChallenge, ActionMonetize:
			if overrideCounts {
				ctx.recordMatch(inner, group.Name, ActionCount)
				continue
			}
			entry, terminated := ctx.applyChallengeLikeAction(inner, group.Name, kind, action)
			entry.OverriddenAction = overridden
			if terminated {
				return true
			}
		}
	}
	if matched && overrideCounts {
		// The reference itself counted at least one inner match.
		ctx.recordMatch(rule, group.Name, ActionCount)
	}
	return false
}

// evaluateManagedRuleGroupRule evaluates one web ACL rule whose statement
// is a managed rule group reference. The group's rules come from the
// managed rules catalog, in the catalog's priority order. The statement's
// scope-down statement must match before any group rule runs; an excluded
// rule runs as Count; a rule action override replaces the catalog action,
// with override names that name no catalog rule silently ignored. The
// return value reports whether the whole evaluation terminated.
func (ctx *evalCtx) evaluateManagedRuleGroupRule(rule *wafstore.Rule) bool {
	stmt := rule.Statement.ManagedRuleGroupStatement
	group, ok := LookupManagedRuleGroup(stmt.VendorName, stmt.Name)
	if !ok {
		ctx.result.Unsupported = append(ctx.result.Unsupported, rule.Name)
		return false
	}
	if stmt.ScopeDownStatement != nil {
		if ctx.evaluateStatement(stmt.ScopeDownStatement, rule.Name, group.Namespace) != outcomeMatch {
			return false
		}
	}
	override := ruleOverrideAction(rule)
	overrideCounts := override != nil && override.Count != nil
	excluded := make(map[string]bool, len(stmt.ExcludedRules))
	for _, excludedRule := range stmt.ExcludedRules {
		excluded[excludedRule.Name] = true
	}

	for _, managed := range group.Rules {
		if managed.Statement == nil {
			// The rule's inspection input exists only inside AWS, so it
			// never matches locally; surface it for visibility.
			ctx.result.Unsupported = append(ctx.result.Unsupported, managed.Name)
			continue
		}
		if ctx.evaluateStatement(managed.Statement, managed.Name, group.Namespace) != outcomeMatch {
			continue
		}
		ctx.addManagedLabel(managed.Label)

		configured := actionKind(managedCatalogAction(managed.Action))
		action := managedCatalogAction(managed.Action)
		if excluded[managed.Name] {
			// Excluding a rule converts its action to Count.
			action = &wafstore.Action{Count: &wafstore.CountAction{}}
		}
		overridden := ""
		if replacement := ruleActionOverrideFor(stmt.RuleActionOverrides, managed.Name); replacement != nil {
			action = replacement
			overridden = configured
		}

		inner := &wafstore.Rule{Name: managed.Name}
		kind := actionKind(action)
		switch kind {
		case ActionAllow, ActionBlock:
			if overrideCounts {
				ctx.recordManagedMatch(managed, group.Name, ActionCount, "")
				continue
			}
			ctx.recordManagedMatch(managed, group.Name, kind, overridden)
			ctx.result.Action = kind
			ctx.applyTerminatingAction(inner, action)
			return true
		case ActionCount, "":
			ctx.recordManagedMatch(managed, group.Name, ActionCount, overridden)
		case ActionCaptcha, ActionChallenge, ActionMonetize:
			if overrideCounts {
				ctx.recordManagedMatch(managed, group.Name, ActionCount, "")
				continue
			}
			entry, terminated := ctx.applyChallengeLikeAction(inner, group.Name, kind, action)
			entry.RuleNameWithinRuleGroup = managedRuleWithinGroupName(group.Name, managed.Name)
			entry.OverriddenAction = overridden
			if terminated {
				return true
			}
		}
	}
	return false
}

// ruleActionOverrideFor returns the configured override action for one
// named rule, or nil when no override names it — the managed-statement
// documentation has overrides naming no rule silently ignored, while
// web ACL updates reject them for customer-owned rule groups before
// the statement ever reaches evaluation.
func ruleActionOverrideFor(overrides []wafstore.RuleActionOverride, ruleName string) *wafstore.Action {
	for _, override := range overrides {
		if override.Name != ruleName {
			continue
		}
		if typed, ok := normaliseThroughJSON[*wafstore.Action](override.ActionToUse); ok {
			return typed
		}
	}
	return nil
}

// managedCatalogAction builds the action object of a catalog rule's
// declared action.
func managedCatalogAction(action string) *wafstore.Action {
	switch action {
	case ActionAllow:
		return &wafstore.Action{Allow: &wafstore.AllowAction{}}
	case ActionBlock:
		return &wafstore.Action{Block: &wafstore.BlockAction{}}
	case ActionCaptcha:
		return &wafstore.Action{Captcha: &wafstore.CaptchaAction{}}
	case ActionChallenge:
		return &wafstore.Action{Challenge: &wafstore.ChallengeAction{}}
	case ActionMonetize:
		return &wafstore.Action{Monetize: &wafstore.MonetizeAction{}}
	default:
		return &wafstore.Action{Count: &wafstore.CountAction{}}
	}
}

// addManagedLabel records a matching managed rule's label. Managed rules
// declare their labels as fully qualified names in the catalog, unlike
// custom rules whose labels combine a namespace with a relative name.
func (ctx *evalCtx) addManagedLabel(label string) {
	if label == "" || ctx.labels[label] {
		return
	}
	ctx.labels[label] = true
	ctx.result.Labels = append(ctx.result.Labels, label)
}

// recordManagedMatch records the match of a managed rule. The sampled
// request's rule name uses the managed rule name format
// vendor#ruleGroup#rule.
func (ctx *evalCtx) recordManagedMatch(managed ManagedRule, groupName, action, overriddenAction string) {
	entry := ctx.recordMatch(&wafstore.Rule{Name: managed.Name}, groupName, action)
	entry.RuleNameWithinRuleGroup = managedRuleWithinGroupName(groupName, managed.Name)
	entry.OverriddenAction = overriddenAction
}

// managedRuleWithinGroupName composes the rule name a managed rule's
// match carries in sampling records.
func managedRuleWithinGroupName(groupName, ruleName string) string {
	return "AWS#" + groupName + "#" + ruleName
}

// evaluateStatement reports how a statement evaluates against the
// request. ownerName keys rate aggregation and labels for the rule that
// owns the statement; ownerNamespace is that rule's label namespace.
// A branch whose forwarded-IP header is absent (outcomeNotApplied)
// propagates through the enclosing logical statements unchanged — the
// API documents it as leaving the whole rule unapplied, so NOT must
// not invert it into a match.
func (ctx *evalCtx) evaluateStatement(stmt *wafstore.Statement, ownerName, ownerNamespace string) statementOutcome {
	switch {
	case stmt.AndStatement != nil:
		result := outcomeMatch
		for _, sub := range stmt.AndStatement.Statements {
			if sub == nil {
				return outcomeNoMatch
			}
			switch outcome := ctx.evaluateStatement(sub, ownerName, ownerNamespace); outcome {
			case outcomeNotApplied:
				return outcomeNotApplied
			case outcomeNoMatch:
				result = outcomeNoMatch
			}
		}
		return result
	case stmt.OrStatement != nil:
		result := outcomeNoMatch
		for _, sub := range stmt.OrStatement.Statements {
			if sub == nil {
				continue
			}
			switch outcome := ctx.evaluateStatement(sub, ownerName, ownerNamespace); outcome {
			case outcomeNotApplied:
				return outcomeNotApplied
			case outcomeMatch:
				result = outcomeMatch
			}
		}
		return result
	case stmt.NotStatement != nil:
		if stmt.NotStatement.Statement == nil {
			return outcomeNoMatch
		}
		switch ctx.evaluateStatement(stmt.NotStatement.Statement, ownerName, ownerNamespace) {
		case outcomeNotApplied:
			return outcomeNotApplied
		case outcomeMatch:
			return outcomeNoMatch
		}
		return outcomeMatch
	case stmt.ByteMatchStatement != nil:
		return boolOutcome(ctx.byteMatch(stmt.ByteMatchStatement))
	case stmt.SqliMatchStatement != nil:
		return boolOutcome(ctx.sqliMatch(stmt.SqliMatchStatement))
	case stmt.XssMatchStatement != nil:
		return boolOutcome(ctx.xssMatch(stmt.XssMatchStatement))
	case stmt.SizeConstraintStatement != nil:
		return boolOutcome(ctx.sizeConstraint(stmt.SizeConstraintStatement))
	case stmt.IPSetReferenceStatement != nil:
		return ctx.e.ipSetMatches(ctx.req, stmt.IPSetReferenceStatement)
	case stmt.RegexMatchStatement != nil:
		return boolOutcome(ctx.regexMatch(stmt.RegexMatchStatement))
	case stmt.RegexPatternSetRefStatement != nil:
		return boolOutcome(ctx.regexPatternSetMatch(stmt.RegexPatternSetRefStatement))
	case stmt.LabelMatchStatement != nil:
		return boolOutcome(ctx.labelMatch(stmt.LabelMatchStatement))
	case stmt.RateBasedStatement != nil:
		return ctx.rateBasedMatch(stmt.RateBasedStatement, ownerName)
	case stmt.GeoMatchStatement != nil:
		return ctx.geoMatch(stmt.GeoMatchStatement)
	case stmt.AsnMatchStatement != nil:
		return ctx.asnMatch(stmt.AsnMatchStatement)
	case stmt.ManagedRuleGroupStatement != nil:
		// Managed rule groups evaluate only as a web ACL rule's
		// top-level statement; statement placement validation rejects
		// nesting, so this branch only guards malformed stored data.
		ctx.result.Unsupported = append(ctx.result.Unsupported, ownerName)
		return outcomeNoMatch
	case stmt.RuleGroupReferenceStatement != nil:
		// Nested references inside logical statements: the inner
		// terminating action cannot surface through the enclosing
		// statement, so the reference matches when any inner rule
		// matches, mirroring the boolean view AWS gives logical
		// statements.
		return boolOutcome(ctx.nestedRuleGroupMatches(stmt.RuleGroupReferenceStatement, ownerName, ownerNamespace))
	}
	return outcomeNoMatch
}

// nestedRuleGroupMatches evaluates a rule group reference nested
// inside a logical statement and reports whether any inner rule
// matched, without propagating terminating actions. Inner rules the
// request leaves unapplied are skipped like non-matching ones.
func (ctx *evalCtx) nestedRuleGroupMatches(stmt *wafstore.RuleGroupReferenceStatement, ownerName, ownerNamespace string) bool {
	group, err := ctx.e.resolvers.RuleGroup(stmt.ARN)
	if err != nil || group == nil {
		return false
	}
	groupNamespace := labelNamespaceOf(group.LabelNamespace, group.ARN, group.Name)
	for _, inner := range group.Rules {
		if inner == nil || inner.Statement == nil {
			continue
		}
		if ctx.evaluateStatement(inner.Statement, inner.Name, groupNamespace) == outcomeMatch {
			ctx.recordLabels(inner, groupNamespace)
			return true
		}
	}
	return false
}

func (ctx *evalCtx) byteMatch(stmt *wafstore.ByteMatchStatement) bool {
	for _, candidate := range extractFieldToMatch(ctx.req, stmt.FieldToMatch) {
		if isMatchSentinel(candidate) {
			return true
		}
		transformed := applyTextTransformations(candidate.value, stmt.TextTransformations)
		if byteMatchConstraint(transformed, stmt.SearchString, stmt.PositionalConstraint) {
			return true
		}
	}
	return false
}

func (ctx *evalCtx) sqliMatch(stmt *wafstore.SqliMatchStatement) bool {
	for _, candidate := range extractFieldToMatch(ctx.req, stmt.FieldToMatch) {
		if isMatchSentinel(candidate) {
			return true
		}
		transformed := applyTextTransformations(candidate.value, stmt.TextTransformations)
		if sqliDetected(transformed, stmt.SensitivityLevel) {
			return true
		}
	}
	return false
}

func (ctx *evalCtx) xssMatch(stmt *wafstore.XssMatchStatement) bool {
	for _, candidate := range extractFieldToMatch(ctx.req, stmt.FieldToMatch) {
		if isMatchSentinel(candidate) {
			return true
		}
		transformed := applyTextTransformations(candidate.value, stmt.TextTransformations)
		if xssDetected(transformed) {
			return true
		}
	}
	return false
}

func (ctx *evalCtx) sizeConstraint(stmt *wafstore.SizeConstraintStatement) bool {
	for _, candidate := range extractFieldToMatch(ctx.req, stmt.FieldToMatch) {
		if isMatchSentinel(candidate) {
			return true
		}
		transformed := applyTextTransformations(candidate.value, stmt.TextTransformations)
		if sizeConstraintHolds(int64(len(transformed)), stmt.ComparisonOperator, stmt.Size) {
			return true
		}
	}
	return false
}

func (ctx *evalCtx) regexMatch(stmt *wafstore.RegexMatchStatement) bool {
	re, err := regexp.Compile(stmt.RegexString)
	if err != nil {
		return false
	}
	for _, candidate := range extractFieldToMatch(ctx.req, stmt.FieldToMatch) {
		if isMatchSentinel(candidate) {
			return true
		}
		transformed := applyTextTransformations(candidate.value, stmt.TextTransformations)
		if re.Match(transformed) {
			return true
		}
	}
	return false
}

func (ctx *evalCtx) regexPatternSetMatch(stmt *wafstore.RegexPatternSetRefStatement) bool {
	set, err := ctx.e.resolvers.RegexSet(stmt.ARN)
	if err != nil || set == nil {
		return false
	}
	patterns := make([]*regexp.Regexp, 0, len(set.RegularPatterns))
	for _, pattern := range set.RegularPatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			patterns = append(patterns, re)
		}
	}
	for _, candidate := range extractFieldToMatch(ctx.req, stmt.FieldToMatch) {
		if isMatchSentinel(candidate) {
			return true
		}
		transformed := applyTextTransformations(candidate.value, stmt.TextTransformations)
		for _, re := range patterns {
			if re.Match(transformed) {
				return true
			}
		}
	}
	return false
}

// labelMatch implements LabelMatchStatement. Labels are case
// sensitive. Scope LABEL matches when the key equals the fully
// qualified label or names a trailing run of whole components of it
// (the key must include the label name and may include any number of
// preceding namespaces and the prefix). Scope NAMESPACE matches when
// the key names a contiguous run of namespace components adjacent to
// the label name, optionally extending to the entire label namespace
// prefix.
func (ctx *evalCtx) labelMatch(stmt *wafstore.LabelMatchStatement) bool {
	for label := range ctx.labels {
		if stmt.Scope == "NAMESPACE" {
			if labelNamespacePrefixMatches(label, stmt.Key) {
				return true
			}
			continue
		}
		if label == stmt.Key || strings.HasSuffix(label, ":"+stmt.Key) {
			return true
		}
	}
	return false
}

// labelNamespacePrefixMatches reports whether key covers a contiguous
// run of the label's namespace components, adjacent to the label name
// and compared at component boundaries, so a key like header:encoding
// cannot match the component header:encoding2 and a key carrying the
// label name itself never matches.
func labelNamespacePrefixMatches(label, key string) bool {
	key = strings.TrimSuffix(key, ":")
	if key == "" {
		return true
	}
	lastColon := strings.LastIndexByte(label, ':')
	if lastColon < 0 {
		return false
	}
	namespacePart := label[:lastColon]
	return namespacePart == key || strings.HasSuffix(namespacePart, ":"+key)
}

// rateBasedMatch implements the rate-based statement: the optional
// ScopeDownStatement must match first, then the request count for the
// aggregated key within the evaluation window must exceed Limit. A
// forwarded-IP configuration whose header is absent leaves the rule
// unapplied rather than non-matching.
func (ctx *evalCtx) rateBasedMatch(stmt *wafstore.RateBasedStatement, ownerName string) statementOutcome {
	acl := ctx.acl
	if stmt.ScopeDownStatement != nil {
		switch ctx.evaluateStatement(stmt.ScopeDownStatement, ownerName, labelNamespaceOf(acl.LabelNamespace, acl.ARN, acl.Name)) {
		case outcomeMatch:
		case outcomeNotApplied:
			return outcomeNotApplied
		default:
			return outcomeNoMatch
		}
	}
	if ctx.e.resolvers.Rate == nil {
		return outcomeNoMatch
	}
	window := time.Duration(stmt.EvaluationWindowSec) * time.Second
	if window <= 0 {
		window = time.Duration(wafstore.RateBasedEvalWindowDefault) * time.Second
	}
	now := ctx.req.Now
	if now.IsZero() {
		now = time.Now()
	}
	if stmt.AggregateKeyType == "FORWARDED_IP" {
		// An absent forwarded header leaves the rule unapplied; a
		// present header without a parseable address follows the
		// fallback behaviour immediately instead of being counted
		// under a key.
		if forwardedHeaderAbsent(ctx.req, stmt.ForwardedIPConfig) {
			return outcomeNotApplied
		}
		ip, matched := ctx.req.forwardedIP(stmt.ForwardedIPConfig)
		if ip == "" {
			return boolOutcome(matched)
		}
		count := ctx.e.resolvers.Rate.Hit(RateKey{WebACLARN: acl.ARN, RuleName: ownerName, Value: ip}, window, now)
		return boolOutcome(count > stmt.Limit)
	}
	// Custom aggregation keys whose value derives from the forwarded
	// header (ForwardedIP and ASN) also leave the rule unapplied when
	// the header is absent.
	if stmt.AggregateKeyType == "CUSTOM_KEYS" && forwardedDerivedCustomKey(stmt) &&
		forwardedHeaderAbsent(ctx.req, stmt.ForwardedIPConfig) {
		return outcomeNotApplied
	}
	for _, key := range ctx.rateKeys(stmt) {
		count := ctx.e.resolvers.Rate.Hit(RateKey{WebACLARN: acl.ARN, RuleName: ownerName, Value: key}, window, now)
		if count > stmt.Limit {
			return outcomeMatch
		}
	}
	return outcomeNoMatch
}

// forwardedHeaderAbsent reports whether a forwarded-IP configuration
// is present but its header is not on the request — the condition the
// API documents as leaving the rule unapplied.
func forwardedHeaderAbsent(req *Request, cfg *wafstore.ForwardedIPConfig) bool {
	return cfg != nil && cfg.HeaderName != "" && !req.hasHeader(cfg.HeaderName)
}

// forwardedDerivedCustomKey reports whether a rate-based statement's
// custom keys include one whose value derives from the forwarded
// header.
func forwardedDerivedCustomKey(stmt *wafstore.RateBasedStatement) bool {
	for _, custom := range stmt.CustomKeys {
		if custom == nil {
			continue
		}
		if custom.ForwardedIP != nil || custom.ASN != nil {
			return true
		}
	}
	return false
}

// rateKeys computes the aggregation key values for a rate-based
// statement according to its AggregateKeyType (the FORWARDED_IP type
// is resolved directly by rateBasedMatch). IP uses the client
// address; CONSTANT aggregates every request under a single key;
// CUSTOM_KEYS builds one aggregation key per unique combination of
// the configured custom key fields.
func (ctx *evalCtx) rateKeys(stmt *wafstore.RateBasedStatement) []string {
	switch stmt.AggregateKeyType {
	case "CONSTANT":
		return []string{"constant"}
	case "CUSTOM_KEYS":
		return ctx.customAggregationKeys(stmt)
	default: // IP
		if !parseIP(ctx.req.SourceIP).IsValid() {
			return nil
		}
		return []string{ctx.req.SourceIP}
	}
}

// customAggregationKeys computes the composite aggregation keys of a
// CUSTOM_KEYS rate-based statement: each unique set of values for the
// aggregation keys is a separate aggregation instance, so the values
// of all keys combine into one composite (joined by a unit separator
// so no combination is ambiguous), and a key with several distinct
// values contributes one component value per combination. A request
// that is missing any key component is omitted from the rate-based
// rule evaluation and handling entirely, per the model's note on
// aggregation keys, and yields no keys.
func (ctx *evalCtx) customAggregationKeys(stmt *wafstore.RateBasedStatement) []string {
	composites := []string{""}
	for _, custom := range stmt.CustomKeys {
		if custom == nil {
			return nil
		}
		values := ctx.customKeyValues(custom, stmt)
		if len(values) == 0 {
			return nil
		}
		var next []string
		for _, prefix := range composites {
			for _, value := range values {
				next = append(next, prefix+value+"\x1f")
			}
		}
		composites = next
	}
	if len(composites) == 1 && len(stmt.CustomKeys) == 0 {
		return nil
	}
	return composites
}

// customKeyValues returns the distinct aggregation values one custom
// key contributes, with the key's own text transformations applied.
// An empty return means the component is missing on this request.
func (ctx *evalCtx) customKeyValues(custom *wafstore.RateBasedStatementCustomKey, stmt *wafstore.RateBasedStatement) []string {
	switch {
	case custom.Header != nil:
		var out []string
		for _, value := range ctx.req.headerValues(custom.Header.Name) {
			out = append(out, string(applyTextTransformations([]byte(value), custom.Header.TextTransformations)))
		}
		return distinctStrings(out)
	case custom.Cookie != nil:
		var out []string
		for _, c := range ctx.req.Cookies {
			if c.Name == custom.Cookie.Name {
				out = append(out, string(applyTextTransformations([]byte(c.Value), custom.Cookie.TextTransformations)))
			}
		}
		return distinctStrings(out)
	case custom.QueryArgument != nil:
		var out []string
		for _, candidate := range queryArgumentCandidates(ctx.req.RawQuery, custom.QueryArgument.Name) {
			out = append(out, string(applyTextTransformations(candidate.value, custom.QueryArgument.TextTransformations)))
		}
		return distinctStrings(out)
	case custom.QueryString != nil:
		return []string{string(applyTextTransformations([]byte(ctx.req.RawQuery), custom.QueryString.TextTransformations))}
	case custom.HTTPMethod != nil:
		return []string{ctx.req.Method}
	case custom.UriPath != nil:
		return []string{string(applyTextTransformations([]byte(ctx.req.URIPath), custom.UriPath.TextTransformations))}
	case custom.IP != nil:
		if !parseIP(ctx.req.SourceIP).IsValid() {
			return nil
		}
		return []string{ctx.req.SourceIP}
	case custom.ForwardedIP != nil:
		// The forwarded address comes from the statement-level
		// configuration; an absent header or unparseable value leaves
		// the component missing.
		ip, _ := ctx.req.forwardedIP(stmt.ForwardedIPConfig)
		if ip == "" {
			return nil
		}
		return []string{ip}
	case custom.LabelNamespace != nil:
		// Each distinct fully qualified label name under the specified
		// namespace contributes; only labels added by rules evaluated
		// before this one are visible.
		namespace := strings.TrimSuffix(custom.LabelNamespace.Namespace, ":") + ":"
		var out []string
		for label := range ctx.labels {
			if strings.HasPrefix(label, namespace) {
				out = append(out, label)
			}
		}
		sort.Strings(out)
		return out
	case custom.JA3Fingerprint != nil:
		return fingerprintKeyValues(custom.JA3Fingerprint.FallbackBehavior)
	case custom.JA4Fingerprint != nil:
		return fingerprintKeyValues(custom.JA4Fingerprint.FallbackBehavior)
	case custom.ASN != nil:
		// The ASN component resolves through the embedded routing-table
		// snapshot; an unresolvable address leaves the component
		// missing, which omits the request from the aggregation.
		addrString, state := ctx.statementAddress(stmt.ForwardedIPConfig)
		if state != addressResolved {
			return nil
		}
		asn, ok := ASNForIP(parseIP(addrString))
		if !ok {
			return nil
		}
		return []string{strconv.FormatUint(uint64(asn), 10)}
	}
	return nil
}

// fingerprintKeyValues reports the aggregation values of a TLS
// fingerprint key. The platform's HTTP listeners do not expose
// ClientHello fingerprints, so the configured fallback decides: MATCH
// counts the component as present with a constant value, any other
// setting leaves it missing.
func fingerprintKeyValues(fallbackBehavior string) []string {
	if fallbackBehavior == "MATCH" {
		return []string{"fingerprint-fallback"}
	}
	return nil
}

// distinctStrings keeps the first occurrence of every value, in
// order: each distinct value contributes to the aggregation instance
// definition, so repeats of one value are not separate instances.
func distinctStrings(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, v := range values {
		if !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	return out
}

// recordLabels adds the rule's labels to the visible label set and the
// result. A fully qualified label is the concatenation of the rule's
// web ACL or rule group label namespace and the label specification;
// the rule's own name is not part of the qualified form.
func (ctx *evalCtx) recordLabels(rule *wafstore.Rule, namespace string) {
	for _, label := range ruleLabels(rule) {
		if label.Name == "" {
			continue
		}
		full := namespace + label.Name
		if !ctx.labels[full] {
			ctx.labels[full] = true
			ctx.result.Labels = append(ctx.result.Labels, full)
		}
	}
}

func (ctx *evalCtx) recordMatch(rule *wafstore.Rule, ruleGroup string, action string) *MatchedRule {
	// The rule's metric name is the one configured in its visibility
	// configuration; the rule name is the fallback for rules created
	// without one.
	metricName := rule.MetricName
	if metricName == "" && rule.VisibilityConfig != nil {
		metricName = rule.VisibilityConfig.MetricName
	}
	var sampled *bool
	if rule.VisibilityConfig != nil {
		enabled := rule.VisibilityConfig.SampledRequestsEnabled
		sampled = &enabled
	}
	entry := MatchedRule{
		RuleName:               rule.Name,
		RuleGroup:              ruleGroup,
		MetricName:             metricName,
		Action:                 action,
		SampledRequestsEnabled: sampled,
	}
	if ruleGroup != "" {
		entry.RuleNameWithinRuleGroup = ruleGroup + "#" + rule.Name
	}
	ctx.result.MatchedRules = append(ctx.result.MatchedRules, entry)
	return &ctx.result.MatchedRules[len(ctx.result.MatchedRules)-1]
}

// applyTerminatingAction resolves the custom response and header
// insertion for an Allow or Block outcome. rule is nil when the action
// comes from the web ACL's default action. Every Block outcome carries
// a resolved response so the enforcement plane always has a status
// code (403 unless configured otherwise).
func (ctx *evalCtx) applyTerminatingAction(rule *wafstore.Rule, action *wafstore.Action) {
	if action == nil {
		return
	}
	if action.Block != nil {
		var custom *wafstore.CustomResponse
		if action.Block.CustomResponse != nil {
			custom = action.Block.CustomResponse
		}
		ctx.result.CustomResponse = ctx.resolveCustomResponse(custom)
	}
	applyInsertHeaders(ctx.result, action)
}

func (ctx *evalCtx) resolveCustomResponse(custom *wafstore.CustomResponse) *ResolvedResponse {
	resolved := &ResolvedResponse{StatusCode: 403}
	if custom != nil {
		resolved.Headers = custom.ResponseHeaders
		if custom.ResponseCode != 0 {
			resolved.StatusCode = custom.ResponseCode
		}
		if custom.CustomResponseBodyKey != "" {
			if body, ok := customResponseBodies(ctx.acl)[custom.CustomResponseBodyKey]; ok {
				resolved.Body = body.Content
			}
		}
	}
	return resolved
}

// challengeActionHeader names the response header the CaptchaAction and
// ChallengeAction documentation attaches to the interrupting response.
const challengeActionHeader = "x-amzn-waf-action"

// applyChallengeLikeAction runs the action of a matching Captcha,
// Challenge or Monetize rule and reports whether the whole evaluation
// terminated, returning the recorded match for the caller to complete.
// Captcha and Challenge first run the token check: with a
// valid, unexpired token the rule behaves like a CountAction — the
// configured custom request handling applies and evaluation continues
// with the next rule — while a missing, invalid or expired token
// interrupts the request with the action's documented response, HTTP
// 405 for Captcha and 202 for Challenge, both carrying
// x-amzn-waf-action. Monetize always interrupts with HTTP 402 carrying
// the price manifest derived from the web ACL's monetization
// configuration.
func (ctx *evalCtx) applyChallengeLikeAction(rule *wafstore.Rule, ruleGroup, kind string, action *wafstore.Action) (*MatchedRule, bool) {
	switch kind {
	case ActionCaptcha, ActionChallenge:
		inspection, satisfied := ctx.tokenInspection(kind, rule)
		entry := ctx.recordMatch(rule, ruleGroup, kind)
		if kind == ActionCaptcha {
			entry.Captcha = inspection
		} else {
			entry.Challenge = inspection
		}
		if satisfied {
			applyInsertHeaders(ctx.result, action)
			return entry, false
		}
		ctx.result.Action = kind
		ctx.result.CustomResponse = challengeInterruptResponse(kind)
		ctx.result.InterstitialRequested = requestAcceptsHTML(ctx.req)
		return entry, true
	case ActionMonetize:
		entry := ctx.recordMatch(rule, ruleGroup, kind)
		ctx.result.Action = ActionMonetize
		ctx.result.CustomResponse = ctx.monetizeResponse(action)
		return entry, true
	}
	return nil, false
}

// tokenInspection runs the aws-waf-token check of a Captcha or
// Challenge rule. The returned inspection carries the kind's solve
// timestamp and, when the check failed, the failure reason; the boolean
// reports whether the token grants immunity for this request.
func (ctx *evalCtx) tokenInspection(kind string, rule *wafstore.Rule) (*TokenInspection, bool) {
	inspection := &TokenInspection{}
	value := ctx.req.cookieValue(TokenCookieName)
	if value == "" {
		inspection.FailureReason = "TOKEN_MISSING"
		return inspection, false
	}
	if ctx.e.resolvers.Token == nil {
		inspection.FailureReason = "TOKEN_INVALID"
		return inspection, false
	}
	token, ok := ctx.e.resolvers.Token.ValidateToken(value)
	if !ok {
		inspection.FailureReason = "TOKEN_INVALID"
		return inspection, false
	}
	switch kind {
	case ActionCaptcha:
		inspection.SolveTimestamp = token.CaptchaSolvedAt
	case ActionChallenge:
		inspection.SolveTimestamp = token.ChallengeSolvedAt
	}
	if !token.CoversHost(ctx.requestHost()) {
		inspection.FailureReason = "TOKEN_DOMAIN_MISMATCH"
		return inspection, false
	}
	now := ctx.req.Now
	if now.IsZero() {
		now = time.Now()
	}
	if !token.SolvedWithin(kind, now, ctx.immunityFor(kind, rule)) {
		inspection.FailureReason = "TOKEN_EXPIRED"
		return inspection, false
	}
	return inspection, true
}

// requestHost returns the request's Host header value; the enforcement
// planes inject Host into the header set they build for inspection.
func (ctx *evalCtx) requestHost() string {
	values := ctx.req.headerValues("Host")
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// immunityFor resolves the immunity window of a Captcha or Challenge
// rule. The rule's own configuration overrides the web ACL's, and an
// unset level falls through to the documented default of 300 seconds.
func (ctx *evalCtx) immunityFor(kind string, rule *wafstore.Rule) time.Duration {
	if rule != nil {
		if duration, ok := immunityOf(ruleImmunityConfig(rule, kind)); ok {
			return duration
		}
	}
	if duration, ok := immunityOf(aclImmunityConfig(ctx.acl, kind)); ok {
		return duration
	}
	return time.Duration(wafstore.ImmunityTimeDefault) * time.Second
}

func ruleImmunityConfig(rule *wafstore.Rule, kind string) interface{} {
	if kind == ActionChallenge {
		return rule.ChallengeConfig
	}
	return rule.CaptchaConfig
}

func aclImmunityConfig(acl *wafstore.WebACL, kind string) interface{} {
	if kind == ActionChallenge {
		return acl.ChallengeConfig
	}
	return acl.CaptchaConfig
}

// immunityOf reads the immunity time from one CaptchaConfig- or
// ChallengeConfig-shaped value; the two configurations carry the same
// single member, so one decoding covers both kinds.
func immunityOf(raw interface{}) (time.Duration, bool) {
	config, ok := normaliseThroughJSON[*wafstore.CaptchaConfig](raw)
	if !ok || config == nil || config.ImmunityTimeProperty == nil {
		return 0, false
	}
	seconds := config.ImmunityTimeProperty.ImmunityTime
	if seconds <= 0 {
		return 0, false
	}
	return time.Duration(seconds) * time.Second, true
}

// challengeInterruptResponse builds the response an interrupting
// Captcha or Challenge evaluation sends: HTTP 405 with the header value
// captcha for Captcha, HTTP 202 with challenge for Challenge. The
// interstitial body is filled in by the serving plane when the request
// accepts text/html.
func challengeInterruptResponse(kind string) *ResolvedResponse {
	response := &ResolvedResponse{
		StatusCode: 405,
		Headers:    []wafstore.CustomHTTPHeader{{Name: challengeActionHeader, Value: strings.ToLower(ActionCaptcha)}},
	}
	if kind == ActionChallenge {
		response.StatusCode = 202
		response.Headers = []wafstore.CustomHTTPHeader{{Name: challengeActionHeader, Value: strings.ToLower(ActionChallenge)}}
	}
	return response
}

// insertHeaderPrefix is the prefix WAF applies to every custom request
// header it inserts, so the names cannot collide with headers already
// on the request (CustomHTTPHeader Name documentation: for the header
// name sample, WAF inserts the header x-amzn-waf-sample).
const insertHeaderPrefix = "x-amzn-waf-"

func applyInsertHeaders(result *Result, action *wafstore.Action) {
	if action == nil {
		return
	}
	var headers []wafstore.CustomHTTPHeader
	switch {
	case action.Allow != nil && action.Allow.CustomRequestHandling != nil:
		headers = action.Allow.CustomRequestHandling.InsertHeaders
	case action.Count != nil && action.Count.CustomRequestHandling != nil:
		headers = action.Count.CustomRequestHandling.InsertHeaders
	}
	for _, h := range headers {
		result.InsertHeaders = append(result.InsertHeaders, wafstore.CustomHTTPHeader{
			Name:  insertHeaderPrefix + h.Name,
			Value: h.Value,
		})
	}
}

// labelNamespaceOf returns the label namespace for a web ACL or rule
// group: the configured value when set, otherwise the default
// awswaf:<account>:<entity type>:<entity name>: that the Developer
// Guide's label syntax specifies. The returned namespace always ends
// with a single colon so qualified labels concatenate directly.
func labelNamespaceOf(configured, entityARN, name string) string {
	if configured == "" {
		configured = defaultLabelNamespace(entityARN, name)
	}
	return strings.TrimSuffix(configured, ":") + ":"
}

// defaultLabelNamespace derives awswaf:<account>:<entity type>:<entity
// name> from the entity's ARN, the fully qualified prefix the
// Developer Guide generates for labels created in web ACLs and rule
// groups. The ARN's resource path starts with the scope segment
// (regional or global), which is not part of the namespace. Entity
// names are unique per account and entity type, so a name without its
// ARN cannot disambiguate the namespace.
func defaultLabelNamespace(entityARN, name string) string {
	account, entityType := "", ""
	if parsed, err := arn.ParseARN(entityARN); err == nil {
		account = parsed.AccountID
		segments := strings.Split(parsed.Resource, "/")
		if len(segments) > 0 {
			index := 0
			if segments[0] == "regional" || segments[0] == "global" {
				index = 1
			}
			if index < len(segments) {
				entityType = segments[index]
			}
		}
	}
	return "awswaf:" + account + ":" + entityType + ":" + name
}

// isMatchSentinel reports whether a candidate is the MATCH fallback
// sentinel, which satisfies every condition.
func isMatchSentinel(candidate fieldCandidate) bool {
	return candidate.oversized && len(candidate.value) == 0
}

func sizeConstraintHolds(actual int64, operator string, size int64) bool {
	switch operator {
	case "EQ":
		return actual == size
	case "NE":
		return actual != size
	case "LE":
		return actual <= size
	case "LT":
		return actual < size
	case "GE":
		return actual >= size
	case "GT":
		return actual > size
	}
	return false
}
