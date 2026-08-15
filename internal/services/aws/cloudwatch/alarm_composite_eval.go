package cloudwatch

import (
	"fmt"

	awserrors "vorpalstacks/internal/common/errors"
	cwstore "vorpalstacks/internal/store/aws/cloudwatch"
)

// evaluateCompositeAlarm evaluates a composite alarm by parsing its
// AlarmRule expression, looking up the current state of all referenced
// child alarms, and evaluating the boolean expression. Returns nil when
// no state transition is needed.
func evaluateCompositeAlarm(alarm *cwstore.Alarm, rule alarmRuleNode, alarmStateMap map[string]string) *alarmEvalResult {
	childNames := rule.childAlarmNames()

	// Build a map of alarm name -> is-in-ALARM-state using the
	// pre-fetched state map.
	alarmStates := make(map[string]bool, len(childNames))
	for _, name := range childNames {
		state := alarmStateMap[name]
		alarmStates[name] = state == "ALARM"
	}

	isBreaching := rule.evaluate(alarmStates)

	oldState := alarm.State
	var newState string
	var reason string

	allInsufficientData := true
	for _, name := range childNames {
		if alarmStateMap[name] != "INSUFFICIENT_DATA" {
			allInsufficientData = false
			break
		}
	}

	if isBreaching {
		newState = "ALARM"
		reason = fmt.Sprintf("Composite alarm rule evaluated to ALARM (was %s).", oldState)
	} else if allInsufficientData && len(childNames) > 0 {
		newState = "INSUFFICIENT_DATA"
		reason = fmt.Sprintf("All referenced alarms are in INSUFFICIENT_DATA state (was %s).", oldState)
	} else {
		newState = "OK"
		reason = fmt.Sprintf("Composite alarm rule evaluated to OK (was %s).", oldState)
	}

	if newState == oldState {
		return nil
	}

	var actionsToFire []string
	switch newState {
	case "ALARM":
		actionsToFire = alarm.AlarmActions
	case "OK":
		actionsToFire = alarm.OKActions
	case "INSUFFICIENT_DATA":
		actionsToFire = alarm.InsufficientDataActions
	}

	return &alarmEvalResult{
		alarm:         alarm,
		oldState:      oldState,
		newState:      newState,
		breachCount:   0,
		reason:        reason,
		actionsToFire: actionsToFire,
	}
}

// topologicalSortLevels performs a level-by-level topological sort using
// Kahn's algorithm. The input is a map of node → nodes-it-depends-on
// (i.e. nodes that must be evaluated before it). The output is a slice of
// levels, where each level contains nodes that can be evaluated in
// parallel. Nodes that remain unsorted are part of (or depend on) a
// dependency cycle and are returned separately so callers can skip them
// without abandoning the rest of the graph.
func topologicalSortLevels(dependencies map[string][]string) (levels [][]string, cyclic []string) {
	inDegree := make(map[string]int)
	dependents := make(map[string][]string)

	for node := range dependencies {
		if _, exists := inDegree[node]; !exists {
			inDegree[node] = 0
		}
		for _, dep := range dependencies[node] {
			if _, exists := inDegree[dep]; !exists {
				inDegree[dep] = 0
			}
			inDegree[node]++
			dependents[dep] = append(dependents[dep], node)
		}
	}

	for {
		var level []string
		for node, deg := range inDegree {
			if deg == 0 {
				level = append(level, node)
			}
		}
		if len(level) == 0 {
			break
		}
		for _, node := range level {
			delete(inDegree, node)
		}
		for _, node := range level {
			for _, dependent := range dependents[node] {
				if _, exists := inDegree[dependent]; exists {
					inDegree[dependent]--
				}
			}
		}
		levels = append(levels, level)
	}

	if len(inDegree) > 0 {
		for node := range inDegree {
			cyclic = append(cyclic, node)
		}
	}

	return levels, cyclic
}

// validateCompositeAcyclic rejects an AlarmRule that would create a
// circular reference — directly (ALARM("itself")) or transitively through
// other composite alarms. AWS rejects such rules when the composite alarm
// is created; a cycle left in the store would keep its alarms (and
// everything downstream of them) unevaluable on every tick.
func validateCompositeAcyclic(alarmStore *cwstore.AlarmStore, alarmName string, rule alarmRuleNode) error {
	for _, child := range rule.childAlarmNames() {
		if child == alarmName {
			return awserrors.NewInvalidParameterValueException(
				fmt.Sprintf("AlarmRule for %s must not reference itself", alarmName))
		}
	}

	// Build the composite-to-composite edges with the new rule in place of
	// any previously stored rule for this alarm.
	edges := map[string][]string{alarmName: rule.childAlarmNames()}

	composites, err := alarmStore.ListAlarms("")
	if err != nil {
		return err
	}
	for _, a := range composites {
		if a.AlarmType != cwstore.AlarmTypeCompositeAlarm || a.Name == alarmName {
			continue
		}
		node, err := parseAlarmRule(a.AlarmRule)
		if err != nil || node == nil {
			continue
		}
		edges[a.Name] = node.childAlarmNames()
	}

	// Only cycles the new rule closes are rejected: a cycle that already
	// exists between unrelated composites keeps those alarms from being
	// evaluated but must not block creation of acyclic alarms. The new
	// rule closes a cycle exactly when some path of composite references
	// leads from the new alarm back to itself.
	visited := map[string]bool{alarmName: true}
	var path []string
	var cyclePath []string
	var dfs func(node string)
	dfs = func(node string) {
		for _, child := range edges[node] {
			// References to metric alarms never form composite cycles.
			if _, isComposite := edges[child]; !isComposite {
				continue
			}
			if child == alarmName {
				cp := make([]string, 0, len(path)+2)
				cp = append(cp, alarmName)
				cp = append(cp, path...)
				cp = append(cp, alarmName)
				cyclePath = cp
				return
			}
			if visited[child] {
				continue
			}
			visited[child] = true
			path = append(path, child)
			dfs(child)
			path = path[:len(path)-1]
			if cyclePath != nil {
				return
			}
		}
	}
	dfs(alarmName)
	if cyclePath != nil {
		return awserrors.NewInvalidParameterValueException(
			fmt.Sprintf("AlarmRule for %s creates a circular reference between composite alarms: %v", alarmName, cyclePath))
	}
	return nil
}
