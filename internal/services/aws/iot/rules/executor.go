package rules

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"
)

// ActionDispatcher dispatches a rule action with the extracted payload fields.
type ActionDispatcher func(ruleName string, topic string, payload map[string]interface{}) error

// Executor evaluates incoming MQTT messages against active rules and dispatches
// matching messages to the action dispatcher.
type Executor struct {
	mu         sync.RWMutex
	rules      map[string]*ActiveRule
	dispatcher ActionDispatcher
	logger     *slog.Logger
}

// ActiveRule represents a compiled rule bound to a topic pattern.
type ActiveRule struct {
	RuleName     string
	TopicPattern string
	SQL          string
	Parsed       *SelectExpr
}

// NewExecutor creates a rules executor backed by the given action dispatcher.
func NewExecutor(dispatcher ActionDispatcher, logger *slog.Logger) *Executor {
	return &Executor{
		rules:      make(map[string]*ActiveRule),
		dispatcher: dispatcher,
		logger:     logger,
	}
}

// AddRule compiles the SQL and registers a new active rule.
func (e *Executor) AddRule(ruleName, topicPattern, sql string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	p := NewParser(sql)
	parsed, err := p.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse rule SQL: %w", err)
	}

	e.rules[ruleName] = &ActiveRule{
		RuleName:     ruleName,
		TopicPattern: topicPattern,
		SQL:          sql,
		Parsed:       parsed,
	}

	return nil
}

// RemoveRule removes a rule by name.
func (e *Executor) RemoveRule(ruleName string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.rules, ruleName)
}

// OnMessage evaluates all active rules against an incoming MQTT message.
func (e *Executor) OnMessage(topic string, payload []byte) {
	e.mu.RLock()
	rules := make([]*ActiveRule, 0, len(e.rules))
	for _, r := range e.rules {
		if MatchTopicFilter(r.TopicPattern, topic) {
			rules = append(rules, r)
		}
	}
	e.mu.RUnlock()

	if len(rules) == 0 {
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		data = map[string]interface{}{"raw": string(payload)}
	}

	for _, rule := range rules {
		e.evaluateAndDispatch(rule, topic, data)
	}
}

func (e *Executor) evaluateAndDispatch(rule *ActiveRule, topic string, data map[string]interface{}) {
	eval := NewEvaluator(data, topic, "")

	if rule.Parsed.Where != nil {
		result, err := eval.Eval(rule.Parsed.Where)
		if err != nil {
			if e.logger != nil {
				e.logger.Error("rule WHERE evaluation failed", "rule", rule.RuleName, "error", err)
			}
			return
		}
		if shouldFilter(result) {
			return
		}
	}

	_, output, err := EvaluateSQL(rule.SQL, data, topic, "")
	if err != nil {
		if e.logger != nil {
			e.logger.Error("rule SQL evaluation failed", "rule", rule.RuleName, "error", err)
		}
		return
	}

	if e.dispatcher != nil {
		go func() {
			if err := e.dispatcher(rule.RuleName, topic, output); err != nil {
				if e.logger != nil {
					e.logger.Error("action dispatch failed", "rule", rule.RuleName, "error", err)
				}
			}
		}()
	}
}

func shouldFilter(result interface{}) bool {
	switch v := result.(type) {
	case bool:
		return !v
	case float64:
		return v == 0
	case nil:
		return true
	default:
		return false
	}
}

// TopicMatches checks if a concrete MQTT topic matches a pattern with + and # wildcards.
func TopicMatches(topic, pattern string) bool {
	return MatchTopicFilter(pattern, topic)
}

func extractSelectedFields(stmt *SelectExpr, data map[string]interface{}) map[string]interface{} {
	if len(stmt.Fields) == 0 {
		return data
	}

	for _, f := range stmt.Fields {
		if _, ok := f.Expression.(*StarExpr); ok {
			return data
		}
	}

	output := make(map[string]interface{})
	eval := NewEvaluator(data, "", "")
	for _, f := range stmt.Fields {
		name := f.Alias
		if name == "" {
			name = exprToName(f.Expression)
		}

		val, err := eval.Eval(f.Expression)
		if err == nil {
			output[name] = val
		}
	}

	return output
}

func exprToName(expr Expr) string {
	switch e := expr.(type) {
	case *Identifier:
		return e.Name
	case *StringLiteral:
		return e.Value
	case *JsonPath:
		return strings.Join(e.Path, ".")
	default:
		return fmt.Sprintf("%v", expr)
	}
}
