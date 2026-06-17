package rules

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolLifecycle(t *testing.T) {
	dispatched := atomic.Int32{}
	dispatcher := func(ruleName, topic string, _ []map[string]interface{}, payload map[string]interface{}) error {
		dispatched.Add(1)
		return nil
	}

	e := NewExecutor(dispatcher, nil, WithPoolSize(2))
	e.Start()

	e.mu.Lock()
	e.rules["r1"] = &ActiveRule{
		RuleName:     "r1",
		TopicPattern: "test/#",
		SQL:          "SELECT * FROM 'test/#'",
		Parsed:       &SelectExpr{},
	}
	e.mu.Unlock()

	for i := 0; i < 10; i++ {
		e.OnMessage("test/foo", []byte(`{"x":1}`))
	}

	e.Stop()

	if got := dispatched.Load(); got != 10 {
		t.Errorf("dispatched count = %d, want 10", got)
	}
}

func TestWorkerPoolPanicIsolation(t *testing.T) {
	shouldPanic := atomic.Bool{}
	dispatched := atomic.Int32{}

	dispatcher := func(ruleName, topic string, _ []map[string]interface{}, payload map[string]interface{}) error {
		if shouldPanic.Load() {
			panic("test panic")
		}
		dispatched.Add(1)
		return nil
	}

	e := NewExecutor(dispatcher, nil, WithPoolSize(2))
	e.Start()

	shouldPanic.Store(true)
	e.mu.Lock()
	e.rules["r1"] = &ActiveRule{
		RuleName:     "r1",
		TopicPattern: "test/#",
		SQL:          "SELECT * FROM 'test/#'",
		Parsed:       &SelectExpr{},
	}
	e.mu.Unlock()
	e.OnMessage("test/foo", []byte(`{}`))

	time.Sleep(100 * time.Millisecond)

	shouldPanic.Store(false)
	e.OnMessage("test/foo", []byte(`{}`))

	e.Stop()

	if got := dispatched.Load(); got != 1 {
		t.Errorf("dispatched count = %d, want 1 (panic task should be skipped, next should succeed)", got)
	}
}

func TestWorkerPoolNoGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()

	dispatcher := func(ruleName, topic string, _ []map[string]interface{}, payload map[string]interface{}) error {
		return nil
	}

	e := NewExecutor(dispatcher, nil, WithPoolSize(4))
	e.Start()

	e.mu.Lock()
	e.rules["r1"] = &ActiveRule{
		RuleName:     "r1",
		TopicPattern: "test/#",
		SQL:          "SELECT * FROM 'test/#'",
		Parsed:       &SelectExpr{},
	}
	e.mu.Unlock()

	for i := 0; i < 50; i++ {
		e.OnMessage("test/foo", []byte(`{"n":1}`))
	}

	e.Stop()

	runtime.GC()
	time.Sleep(50 * time.Millisecond)

	after := runtime.NumGoroutine()
	if delta := after - before; delta > 6 {
		t.Errorf("goroutine leak detected: before=%d, after=%d, delta=%d", before, after, delta)
	}
}

func TestWorkerPoolStopDrains(t *testing.T) {
	var mu sync.Mutex
	dispatched := []string{}
	dispatcher := func(ruleName, topic string, _ []map[string]interface{}, payload map[string]interface{}) error {
		mu.Lock()
		dispatched = append(dispatched, ruleName)
		mu.Unlock()
		return nil
	}

	e := NewExecutor(dispatcher, nil, WithPoolSize(1))
	e.Start()

	e.mu.Lock()
	e.rules["r1"] = &ActiveRule{
		RuleName:     "r1",
		TopicPattern: "test/#",
		SQL:          "SELECT * FROM 'test/#'",
		Parsed:       &SelectExpr{},
	}
	e.mu.Unlock()

	for i := 0; i < 5; i++ {
		e.OnMessage("test/foo", []byte(fmt.Sprintf(`{"i":%d}`, i)))
	}

	time.Sleep(100 * time.Millisecond)

	e.Stop()

	mu.Lock()
	if len(dispatched) != 5 {
		t.Errorf("dispatched count = %d, want 5", len(dispatched))
	}
	mu.Unlock()
}

func TestExecutorRulesCount(t *testing.T) {
	e := NewExecutor(nil, nil)
	e.AddRule("r1", "test/#", "SELECT * FROM 'test/#'", nil)
	e.AddRule("r2", "foo/#", "SELECT * FROM 'foo/#'", nil)

	if got := e.RulesCount(); got != 2 {
		t.Errorf("RulesCount = %d, want 2", got)
	}

	e.RemoveRule("r1")
	if got := e.RulesCount(); got != 1 {
		t.Errorf("RulesCount = %d, want 1", got)
	}
}

func TestExecutorEndToEnd_SQLParseTopicMatchDispatch(t *testing.T) {
	var mu sync.Mutex
	fired := []struct {
		rule   string
		topic  string
		fields map[string]interface{}
	}{}

	dispatcher := func(ruleName, topic string, _ []map[string]interface{}, payload map[string]interface{}) error {
		mu.Lock()
		fired = append(fired, struct {
			rule   string
			topic  string
			fields map[string]interface{}
		}{ruleName, topic, payload})
		mu.Unlock()
		return nil
	}

	e := NewExecutor(dispatcher, nil, WithPoolSize(2))
	e.Start()

	err := e.AddRule("temp_high", "sensors/+/temperature", "SELECT temperature, location FROM 'sensors/+/temperature' WHERE temperature > 30", nil)
	if err != nil {
		t.Fatalf("AddRule failed: %v", err)
	}

	e.OnMessage("sensors/device1/temperature", []byte(`{"temperature": 35, "location": "warehouse"}`))
	e.OnMessage("sensors/device2/temperature", []byte(`{"temperature": 20, "location": "office"}`))

	e.Stop()

	mu.Lock()
	if len(fired) != 1 {
		t.Fatalf("fired count = %d, want 1 (WHERE should filter cold reading)", len(fired))
	}
	if fired[0].rule != "temp_high" {
		t.Errorf("fired rule = %q, want %q", fired[0].rule, "temp_high")
	}
	if fired[0].topic != "sensors/device1/temperature" {
		t.Errorf("fired topic = %q, want %q", fired[0].topic, "sensors/device1/temperature")
	}
	mu.Unlock()
}

func TestExecutorEndToEnd_DisableDeleteRemovesRule(t *testing.T) {
	dispatched := atomic.Int32{}
	dispatcher := func(ruleName, topic string, _ []map[string]interface{}, payload map[string]interface{}) error {
		dispatched.Add(1)
		return nil
	}

	e := NewExecutor(dispatcher, nil, WithPoolSize(1))
	e.Start()

	e.AddRule("alert", "alerts/#", "SELECT * FROM 'alerts/#'", nil)

	e.OnMessage("alerts/fire", []byte(`{"msg":"fire"}`))
	time.Sleep(50 * time.Millisecond)
	if dispatched.Load() != 1 {
		t.Fatalf("before disable: dispatched = %d, want 1", dispatched.Load())
	}

	e.RemoveRule("alert")

	e.OnMessage("alerts/fire", []byte(`{"msg":"fire"}`))
	time.Sleep(50 * time.Millisecond)

	e.Stop()

	if dispatched.Load() != 1 {
		t.Errorf("after remove: dispatched = %d, want 1 (removed rule should not fire)", dispatched.Load())
	}
}

func TestExecutorEndToEnd_ReplaceRule(t *testing.T) {
	dispatched := atomic.Int32{}
	dispatcher := func(ruleName, topic string, _ []map[string]interface{}, payload map[string]interface{}) error {
		dispatched.Add(1)
		return nil
	}

	e := NewExecutor(dispatcher, nil, WithPoolSize(1))
	e.Start()

	e.AddRule("r1", "old/#", "SELECT * FROM 'old/#'", nil)

	e.OnMessage("old/topic", []byte(`{}`))
	time.Sleep(50 * time.Millisecond)
	if dispatched.Load() != 1 {
		t.Fatalf("before replace: dispatched = %d, want 1", dispatched.Load())
	}

	e.RemoveRule("r1")
	e.AddRule("r1", "new/#", "SELECT * FROM 'new/#'", nil)

	e.OnMessage("old/topic", []byte(`{}`))
	e.OnMessage("new/topic", []byte(`{}`))
	time.Sleep(50 * time.Millisecond)

	e.Stop()

	if dispatched.Load() != 2 {
		t.Errorf("after replace: dispatched = %d, want 2", dispatched.Load())
	}
}

func TestShouldFilter_UnknownValue(t *testing.T) {
	cases := []struct {
		input  interface{}
		want   bool
		reason string
	}{
		{unknownValue{}, true, "SQL three-valued logic: UNKNOWN treated as false"},
		{true, false, "true passes through"},
		{false, true, "false is filtered out"},
		{float64(0), true, "zero is filtered out"},
		{float64(1), false, "non-zero passes through"},
		{nil, true, "nil is filtered out"},
		{"", false, "non-empty string passes through"},
	}
	for _, c := range cases {
		got := shouldFilter(c.input)
		if got != c.want {
			t.Errorf("shouldFilter(%v) = %v, want %v (%s)", c.input, got, c.want, c.reason)
		}
	}
}

func TestExtractSelectedFields_TopicPropagation(t *testing.T) {
	data := map[string]interface{}{
		"temperature": 25.5,
		"humidity":    60.0,
	}

	stmt := &SelectExpr{
		Fields: []SelectField{
			{Expression: &Identifier{Name: "temperature"}, Alias: "temp"},
			{Expression: &Identifier{Name: "humidity"}, Alias: "hum"},
		},
	}

	result := extractSelectedFields(stmt, data, "sensor/data", "client-42")

	if result["temp"] != 25.5 {
		t.Errorf("expected temp=25.5, got %v", result["temp"])
	}
	if result["hum"] != 60.0 {
		t.Errorf("expected hum=60.0, got %v", result["hum"])
	}
	if _, exists := result["temperature"]; exists {
		t.Error("expected original field 'temperature' to NOT be present (aliased to 'temp')")
	}
	if _, exists := result["humidity"]; exists {
		t.Error("expected original field 'humidity' to NOT be present (aliased to 'hum')")
	}
}

func TestExtractSelectedFields_StarReturnsAll(t *testing.T) {
	data := map[string]interface{}{
		"temperature": 25.5,
		"humidity":    60.0,
	}

	stmt := &SelectExpr{
		Fields: []SelectField{
			{Expression: &StarExpr{}},
		},
	}

	result := extractSelectedFields(stmt, data, "sensor/data", "client-42")

	if result["temperature"] != 25.5 {
		t.Errorf("expected temperature=25.5, got %v", result["temperature"])
	}
	if result["humidity"] != 60.0 {
		t.Errorf("expected humidity=60.0, got %v", result["humidity"])
	}
}

func TestExtractSelectedFields_NoFieldsReturnsAll(t *testing.T) {
	data := map[string]interface{}{
		"temperature": 25.5,
	}

	stmt := &SelectExpr{}

	result := extractSelectedFields(stmt, data, "sensor/data", "client-42")

	if result["temperature"] != 25.5 {
		t.Errorf("expected temperature=25.5, got %v", result["temperature"])
	}
}
