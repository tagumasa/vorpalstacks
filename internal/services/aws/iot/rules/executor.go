package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"sync"
	"time"
)

const defaultDispatchTimeout = 30 * time.Second

// ActionDispatcher dispatches rule actions with the extracted payload fields.
// Each action in the actions list is dispatched by its type (lambda, sns, sqs, etc.).
type ActionDispatcher func(ruleName string, topic string, actions []map[string]interface{}, payload map[string]interface{}) error

// dispatchTask represents a single rule-fire action queued for the worker pool.
type dispatchTask struct {
	ruleName string
	topic    string
	actions  []map[string]interface{}
	payload  map[string]interface{}
}

// workerPool provides a bounded pool of goroutines for rule dispatch.
type workerPool struct {
	queue      chan dispatchTask
	dispatcher ActionDispatcher
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFn   context.CancelFunc
	logger     *slog.Logger
}

func newWorkerPool(size, queueDepth int, dispatcher ActionDispatcher, logger *slog.Logger) *workerPool {
	ctx, cancelFn := context.WithCancel(context.Background())
	return &workerPool{
		queue:      make(chan dispatchTask, queueDepth),
		dispatcher: dispatcher,
		ctx:        ctx,
		cancelFn:   cancelFn,
		logger:     logger,
	}
}

func (p *workerPool) Start(workers int) {
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker()
	}
}

func (p *workerPool) worker() {
	defer p.wg.Done()
	for task := range p.queue {
		p.execute(task)
	}
}

func (p *workerPool) execute(task dispatchTask) {
	defer func() {
		if r := recover(); r != nil {
			if p.logger != nil {
				p.logger.Error("rule dispatch panic recovered", "rule", task.ruleName, "panic", r)
			}
		}
	}()

	ctx, cancel := context.WithTimeout(p.ctx, defaultDispatchTimeout)
	defer cancel()

	if p.dispatcher == nil {
		return
	}
	if err := p.dispatcher(task.ruleName, task.topic, task.actions, task.payload); err != nil {
		if ctx.Err() == nil && p.logger != nil {
			p.logger.Error("action dispatch failed", "rule", task.ruleName, "error", err)
		}
	}
}

func (p *workerPool) Stop(timeout time.Duration) {
	close(p.queue)
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		if p.logger != nil {
			p.logger.Warn("worker pool stop timed out, some tasks may be lost")
		}
	}
	p.cancelFn()
}

func (p *workerPool) Enqueue(task dispatchTask) error {
	select {
	case p.queue <- task:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool shut down")
	}
}

// Executor evaluates incoming MQTT messages against active rules and dispatches
// matching messages to the action dispatcher.
type Executor struct {
	mu         sync.RWMutex
	rules      map[string]*ActiveRule
	dispatcher ActionDispatcher
	logger     *slog.Logger
	pool       *workerPool
	poolSize   int
}

// ActiveRule represents a compiled rule bound to a topic pattern.
type ActiveRule struct {
	RuleName     string
	TopicPattern string
	SQL          string
	Actions      []map[string]interface{}
	Parsed       *SelectExpr
}

// ExecutorOption configures the Executor.
type ExecutorOption func(*Executor)

// WithPoolSize sets the worker pool size. If 0, defaults to runtime.NumCPU() * 2.
func WithPoolSize(n int) ExecutorOption {
	return func(e *Executor) {
		if n > 0 {
			e.poolSize = n
		}
	}
}

// NewExecutor creates a rules executor backed by the given action dispatcher.
func NewExecutor(dispatcher ActionDispatcher, logger *slog.Logger, opts ...ExecutorOption) *Executor {
	e := &Executor{
		rules:      make(map[string]*ActiveRule),
		dispatcher: dispatcher,
		logger:     logger,
		poolSize:   runtime.NumCPU() * 2,
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

// Start initialises the worker pool and makes the executor ready to process messages.
func (e *Executor) Start() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.pool != nil {
		return
	}
	e.pool = newWorkerPool(e.poolSize, e.poolSize*4, e.dispatcher, e.logger)
	e.pool.Start(e.poolSize)
}

// Stop drains and shuts down the worker pool.
func (e *Executor) Stop() {
	e.mu.Lock()
	pool := e.pool
	e.pool = nil
	e.mu.Unlock()
	if pool != nil {
		pool.Stop(5 * time.Second)
	}
}

// RulesCount returns the number of active rules.
func (e *Executor) RulesCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.rules)
}

// AddRule compiles the SQL and registers a new active rule.
// The topic pattern is extracted from the SQL FROM clause; the topicPattern
// parameter is used only when the SQL has no FROM clause (e.g. for testing).
func (e *Executor) AddRule(ruleName, topicPattern, sql string, actions []map[string]interface{}) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	p := NewParser(sql)
	parsed, err := p.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse rule SQL: %w", err)
	}

	if parsed.FromTopic != "" {
		topicPattern = parsed.FromTopic
	}

	e.rules[ruleName] = &ActiveRule{
		RuleName:     ruleName,
		TopicPattern: topicPattern,
		SQL:          sql,
		Actions:      actions,
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

	output := extractSelectedFields(rule.Parsed, data, topic, "")

	if e.pool != nil {
		e.pool.Enqueue(dispatchTask{
			ruleName: rule.RuleName,
			topic:    topic,
			actions:  rule.Actions,
			payload:  output,
		})
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
	case unknownValue:
		// SQL three-valued logic: UNKNOWN in a WHERE clause is treated as
		// false, which means the row is filtered out.
		return true
	default:
		return false
	}
}

// TopicMatches checks if a concrete MQTT topic matches a pattern with + and # wildcards.
func TopicMatches(topic, pattern string) bool {
	return MatchTopicFilter(pattern, topic)
}

func extractSelectedFields(stmt *SelectExpr, data map[string]interface{}, topic, clientID string) map[string]interface{} {
	if len(stmt.Fields) == 0 {
		return data
	}

	for _, f := range stmt.Fields {
		if _, ok := f.Expression.(*StarExpr); ok {
			return data
		}
	}

	output := make(map[string]interface{})
	eval := NewEvaluator(data, topic, clientID)
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
