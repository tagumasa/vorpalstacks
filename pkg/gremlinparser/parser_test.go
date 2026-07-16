package gremlinparser

import (
	"reflect"
	"testing"
)

func TestParser_BasicTraversal(t *testing.T) {
	q, err := Parse("g.V()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(q.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(q.Statements))
	}
	s := q.Statements[0]
	if s.Traversal.Source != "" {
		t.Errorf("expected empty source, got %q", s.Traversal.Source)
	}
	if len(s.Traversal.Steps) != 1 || s.Traversal.Steps[0].Name != "V" {
		t.Errorf("expected step V, got %v", s.Traversal.Steps)
	}
}

func TestParser_WithTerminal(t *testing.T) {
	q, err := Parse("g.V().toList()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Terminal == nil || s.Traversal.Terminal.Name != "toList" {
		t.Errorf("expected terminal toList, got %v", s.Traversal.Terminal)
	}
}

func TestParser_ChainedSteps(t *testing.T) {
	q, err := Parse("g.V().has('person','name','marko').out('knows').values('name').toList()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	expectedSteps := []string{"V", "has", "out", "values"}
	if len(s.Traversal.Steps) != len(expectedSteps) {
		t.Fatalf("expected %d steps, got %d", len(expectedSteps), len(s.Traversal.Steps))
	}
	for i, name := range expectedSteps {
		if s.Traversal.Steps[i].Name != name {
			t.Errorf("step %d: expected %q, got %q", i, name, s.Traversal.Steps[i].Name)
		}
	}
	if s.Traversal.Terminal == nil || s.Traversal.Terminal.Name != "toList" {
		t.Errorf("expected terminal toList")
	}
}

func TestParser_HasWithPredicate(t *testing.T) {
	q, err := Parse("g.V().has('age',gt(30))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	hasStep := s.Traversal.Steps[1]
	if hasStep.Name != "has" {
		t.Fatalf("expected step has, got %q", hasStep.Name)
	}
	if len(hasStep.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(hasStep.Args))
	}
	if hasStep.Args[1].Kind != ArgPredicate {
		t.Fatalf("expected predicate arg, got %d", hasStep.Args[1].Kind)
	}
	pred := hasStep.Args[1].Pred
	if pred.Method != "gt" {
		t.Errorf("expected predicate gt, got %q", pred.Method)
	}
	if len(pred.Args) != 1 || pred.Args[0].Int != 30 {
		t.Errorf("expected pred arg 30, got %v", pred.Args)
	}
}

func TestParser_Predicate(t *testing.T) {
	q, err := Parse("g.V().has('age',P.gt(30))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	pred := s.Traversal.Steps[1].Args[1].Pred
	if pred.Type != "P" || pred.Method != "gt" {
		t.Errorf("expected P.gt, got %s.%s", pred.Type, pred.Method)
	}
}

func TestParser_TextPredicate(t *testing.T) {
	q, err := Parse("g.V().has('name',TextP.startingWith('mar'))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	pred := s.Traversal.Steps[1].Args[1].Pred
	if pred.Type != "TextP" || pred.Method != "startingWith" {
		t.Errorf("expected TextP.startingWith, got %s.%s", pred.Type, pred.Method)
	}
}

func TestParser_BetweenPredicate(t *testing.T) {
	q, err := Parse("g.V().has('age',between(20,30))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	pred := s.Traversal.Steps[1].Args[1].Pred
	if pred.Method != "between" {
		t.Errorf("expected between, got %q", pred.Method)
	}
	if len(pred.Args) != 2 || pred.Args[0].Int != 20 || pred.Args[1].Int != 30 {
		t.Errorf("expected args [20, 30], got %v", pred.Args)
	}
}

func TestParser_WithinPredicate(t *testing.T) {
	q, err := Parse("g.V().has('name',within('marko','vadas','josh'))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	pred := s.Traversal.Steps[1].Args[1].Pred
	if pred.Method != "within" {
		t.Errorf("expected within, got %q", pred.Method)
	}
	if len(pred.Args) != 3 {
		t.Errorf("expected 3 args, got %d", len(pred.Args))
	}
}

func TestParser_NestedTraversal(t *testing.T) {
	q, err := Parse("g.V().where(__.out('knows').has('age',gt(30)))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	whereStep := s.Traversal.Steps[1]
	if whereStep.Name != "where" {
		t.Fatalf("expected where step, got %q", whereStep.Name)
	}
	if len(whereStep.Args) != 1 || whereStep.Args[0].Kind != ArgNestedTraversal {
		t.Fatalf("expected nested traversal arg, got %d", whereStep.Args[0].Kind)
	}
	nested := whereStep.Args[0].Trav
	if nested.Source != "__" {
		t.Errorf("expected __ source, got %q", nested.Source)
	}
	if len(nested.Steps) != 2 || nested.Steps[0].Name != "out" || nested.Steps[1].Name != "has" {
		t.Errorf("expected nested steps [out, has], got %v", nested.Steps)
	}
}

func TestParser_Union(t *testing.T) {
	q, err := Parse("g.V().union(__.out('knows'),__.out('created'))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	unionStep := s.Traversal.Steps[1]
	if unionStep.Name != "union" {
		t.Fatalf("expected union step, got %q", unionStep.Name)
	}
	if len(unionStep.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(unionStep.Args))
	}
	if unionStep.Args[0].Trav.Steps[0].Name != "out" || unionStep.Args[0].Trav.Steps[0].Args[0].Str != "knows" {
		t.Errorf("expected first union branch out('knows'), got %v", unionStep.Args[0])
	}
}

func TestParser_RepeatUntil(t *testing.T) {
	q, err := Parse("g.V().repeat(__.out()).until(__.has('name','lop')).path()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if len(s.Traversal.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d: %v", len(s.Traversal.Steps), s.Traversal.Steps)
	}
	if s.Traversal.Steps[0].Name != "V" {
		t.Errorf("expected V, got %q", s.Traversal.Steps[0].Name)
	}
	if s.Traversal.Steps[1].Name != "repeat" {
		t.Errorf("expected repeat, got %q", s.Traversal.Steps[1].Name)
	}
	untilMod := getMod(s.Traversal.Steps[1], "until")
	if len(untilMod) != 1 {
		t.Errorf("expected until modulator on repeat, got %v", s.Traversal.Steps[1].Modulators)
	}
	if s.Traversal.Steps[2].Name != "path" {
		t.Errorf("expected path, got %q", s.Traversal.Steps[2].Name)
	}
}

func getMod(step Step, name string) []Step {
	var result []Step
	for _, m := range step.Modulators {
		if m.Name == name {
			result = append(result, m)
		}
	}
	return result
}

func TestParser_Enum(t *testing.T) {
	q, err := Parse("g.V().order().by('age',Order.desc)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	orderStep := s.Traversal.Steps[1]
	if orderStep.Name != "order" {
		t.Fatalf("expected order step, got %q", orderStep.Name)
	}
	byMods := getMod(orderStep, "by")
	if len(byMods) != 1 {
		t.Fatalf("expected 1 by() modulator on order, got %d", len(byMods))
	}
}

func TestParser_BareEnum(t *testing.T) {
	q, err := Parse("g.V().order().by(asc)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	orderStep := s.Traversal.Steps[1]
	if orderStep.Name != "order" {
		t.Fatalf("expected order step, got %q", orderStep.Name)
	}
}

func TestParser_TEnum(t *testing.T) {
	q, err := Parse("g.V().has(T.id,1)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	hasStep := s.Traversal.Steps[1]
	if hasStep.Args[0].Kind != ArgEnum || hasStep.Args[0].Enum.Value != "id" {
		t.Errorf("expected T.id, got %v", hasStep.Args[0])
	}
}

func TestParser_MapLiteral(t *testing.T) {
	q, err := Parse("g.mergeV([name:'marko',age:29])")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	mergeStep := s.Traversal.Steps[0]
	if mergeStep.Name != "mergeV" {
		t.Fatalf("expected mergeV step, got %q", mergeStep.Name)
	}
	if len(mergeStep.Args) != 1 || mergeStep.Args[0].Kind != ArgMap {
		t.Fatalf("expected map arg, got %d", mergeStep.Args[0].Kind)
	}
	m := mergeStep.Args[0].Map
	if len(m) != 2 {
		t.Fatalf("expected 2 map entries, got %d", len(m))
	}
	if m[0].Key.Str != "name" || m[0].Value.Str != "marko" {
		t.Errorf("expected name:marko, got %v", m[0])
	}
	if m[1].Key.Str != "age" || m[1].Value.Int != 29 {
		t.Errorf("expected age:29, got %v", m[1])
	}
}

func TestParser_AddV(t *testing.T) {
	q, err := Parse("g.addV('person').property('name','marko').property('age',29)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	expectedSteps := []string{"addV", "property", "property"}
	if len(s.Traversal.Steps) != len(expectedSteps) {
		t.Fatalf("expected %d steps, got %d", len(expectedSteps), len(s.Traversal.Steps))
	}
	for i, name := range expectedSteps {
		if s.Traversal.Steps[i].Name != name {
			t.Errorf("step %d: expected %q, got %q", i, name, s.Traversal.Steps[i].Name)
		}
	}
}

func TestParser_Group(t *testing.T) {
	q, err := Parse("g.V().group().by('name').by(__.values('age').fold())")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[1].Name != "group" {
		t.Errorf("expected group, got %q", s.Traversal.Steps[1].Name)
	}
	byMods := getMod(s.Traversal.Steps[1], "by")
	if len(byMods) != 2 {
		t.Errorf("expected 2 by() modulators on group, got %d", len(byMods))
	}
}

func TestParser_Choose(t *testing.T) {
	q, err := Parse("g.V().choose(__.hasLabel('person'),__.values('name'),__.values('title'))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	chooseStep := s.Traversal.Steps[1]
	if chooseStep.Name != "choose" {
		t.Fatalf("expected choose, got %q", chooseStep.Name)
	}
	if len(chooseStep.Args) != 3 {
		t.Fatalf("expected 3 args, got %d", len(chooseStep.Args))
	}
	for i, arg := range chooseStep.Args {
		if arg.Kind != ArgNestedTraversal {
			t.Errorf("arg %d: expected nested traversal, got %d", i, arg.Kind)
		}
	}
}

func TestParser_Coalesce(t *testing.T) {
	q, err := Parse("g.V().coalesce(__.out('knows'),__.out('created'))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	coalesceStep := s.Traversal.Steps[1]
	if coalesceStep.Name != "coalesce" {
		t.Fatalf("expected coalesce, got %q", coalesceStep.Name)
	}
}

func TestParser_AsSelect(t *testing.T) {
	q, err := Parse("g.V().as('a').out('knows').as('b').select('a','b')")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[0].Name != "V" {
		t.Errorf("expected V, got %q", s.Traversal.Steps[0].Name)
	}
	asModA := getMod(s.Traversal.Steps[0], "as")
	if len(asModA) != 1 || asModA[0].Args[0].Str != "a" {
		t.Errorf("expected as('a') modulator on V, got %v", s.Traversal.Steps[0].Modulators)
	}
	if s.Traversal.Steps[1].Name != "out" {
		t.Errorf("expected out, got %q", s.Traversal.Steps[1].Name)
	}
	asModB := getMod(s.Traversal.Steps[1], "as")
	if len(asModB) != 1 || asModB[0].Args[0].Str != "b" {
		t.Errorf("expected as('b') modulator on out, got %v", s.Traversal.Steps[1].Modulators)
	}
	selectStep := s.Traversal.Steps[2]
	if selectStep.Name != "select" || len(selectStep.Args) != 2 {
		t.Errorf("expected select('a','b'), got %v", selectStep)
	}
}

func TestParser_PropertyCardinality(t *testing.T) {
	q, err := Parse("g.addV('person').property(Cardinality.single,'name','marko')")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	propStep := s.Traversal.Steps[1]
	if propStep.Args[0].Kind != ArgEnum || propStep.Args[0].Enum.Value != "single" {
		t.Errorf("expected Cardinality.single, got %v", propStep.Args[0])
	}
}

func TestParser_NextWithArg(t *testing.T) {
	q, err := Parse("g.V(1).next(3)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Terminal == nil || s.Traversal.Terminal.Name != "next" {
		t.Fatalf("expected terminal next, got %v", s.Traversal.Terminal)
	}
	if len(s.Traversal.Terminal.Args) != 1 || s.Traversal.Terminal.Args[0].Int != 3 {
		t.Errorf("expected next(3), got %v", s.Traversal.Terminal)
	}
}

func TestParser_EdgeTraversal(t *testing.T) {
	q, err := Parse("g.E().hasLabel('knows').inV()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[0].Name != "E" {
		t.Errorf("expected E, got %q", s.Traversal.Steps[0].Name)
	}
	if s.Traversal.Steps[2].Name != "inV" {
		t.Errorf("expected inV, got %q", s.Traversal.Steps[2].Name)
	}
}

func TestParser_OptionMergeV(t *testing.T) {
	q, err := Parse("g.mergeV([name:'marko']).option(onCreate,property('age',29))")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if len(s.Traversal.Steps) != 1 {
		t.Fatalf("expected 1 step, got %d: %v", len(s.Traversal.Steps), s.Traversal.Steps)
	}
	optionMod := getMod(s.Traversal.Steps[0], "option")
	if len(optionMod) != 1 {
		t.Errorf("expected option modulator on mergeV, got %v", s.Traversal.Steps[0].Modulators)
	}
}

func TestParser_LimitSkip(t *testing.T) {
	q, err := Parse("g.V().hasLabel('person').limit(10).skip(5)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[2].Name != "limit" || s.Traversal.Steps[2].Args[0].Int != 10 {
		t.Errorf("expected limit(10), got %v", s.Traversal.Steps[2])
	}
	if s.Traversal.Steps[3].Name != "skip" || s.Traversal.Steps[3].Args[0].Int != 5 {
		t.Errorf("expected skip(5), got %v", s.Traversal.Steps[3])
	}
}

func TestParser_MultipleStatements(t *testing.T) {
	q, err := Parse("g.addV('person').property('name','marko');g.V().has('name','marko').values('name')")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(q.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(q.Statements))
	}
	if q.Statements[0].Traversal.Steps[0].Name != "addV" {
		t.Errorf("stmt 0: expected addV, got %q", q.Statements[0].Traversal.Steps[0].Name)
	}
	if q.Statements[1].Traversal.Steps[0].Name != "V" {
		t.Errorf("stmt 1: expected V, got %q", q.Statements[1].Traversal.Steps[0].Name)
	}
}

func TestParser_ValuesWithMultipleKeys(t *testing.T) {
	q, err := Parse("g.V().values('name','age')")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	valStep := s.Traversal.Steps[1]
	if len(valStep.Args) != 2 {
		t.Fatalf("expected 2 args, got %d", len(valStep.Args))
	}
}

func TestParser_RepeatTimes(t *testing.T) {
	q, err := Parse("g.V(1).repeat(__.out()).times(2).emit().path()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	names := []string{}
	for _, step := range s.Traversal.Steps {
		names = append(names, step.Name)
	}
	expected := []string{"V", "repeat", "path"}
	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected %v, got %v", expected, names)
	}
	repeatStep := s.Traversal.Steps[1]
	timesMod := getMod(repeatStep, "times")
	if len(timesMod) != 1 || timesMod[0].Args[0].Int != 2 {
		t.Errorf("expected times(2) modulator on repeat, got %v", repeatStep.Modulators)
	}
	emitMod := getMod(repeatStep, "emit")
	if len(emitMod) != 1 {
		t.Errorf("expected emit modulator on repeat, got %v", repeatStep.Modulators)
	}
}

func TestParser_Drop(t *testing.T) {
	q, err := Parse("g.V().has('name','marko').drop()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[2].Name != "drop" {
		t.Errorf("expected drop, got %q", s.Traversal.Steps[2].Name)
	}
}

func TestParser_Dedup(t *testing.T) {
	q, err := Parse("g.V().hasLabel('person').dedup().by('name')")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[2].Name != "dedup" {
		t.Errorf("expected dedup, got %q", s.Traversal.Steps[2].Name)
	}
	byMod := getMod(s.Traversal.Steps[2], "by")
	if len(byMod) != 1 || byMod[0].Args[0].Str != "name" {
		t.Errorf("expected by('name') modulator on dedup, got %v", s.Traversal.Steps[2].Modulators)
	}
}

func TestParser_Count(t *testing.T) {
	q, err := Parse("g.V().hasLabel('person').count()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[2].Name != "count" {
		t.Errorf("expected count, got %q", s.Traversal.Steps[2].Name)
	}
}

func TestParser_PropertyMap(t *testing.T) {
	q, err := Parse("g.V().propertyMap('name','age')")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[1].Name != "propertyMap" {
		t.Errorf("expected propertyMap, got %q", s.Traversal.Steps[1].Name)
	}
}

func TestParser_BothInOut(t *testing.T) {
	q, err := Parse("g.V(1).both('knows').in('created').outE().inV().otherV()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	names := []string{}
	for _, step := range s.Traversal.Steps {
		names = append(names, step.Name)
	}
	expected := []string{"V", "both", "in", "outE", "inV", "otherV"}
	if !reflect.DeepEqual(names, expected) {
		t.Errorf("expected %v, got %v", expected, names)
	}
}

func TestParser_LabelStep(t *testing.T) {
	q, err := Parse("g.V().label()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[1].Name != "label" {
		t.Errorf("expected label, got %q", s.Traversal.Steps[1].Name)
	}
}

func TestParser_IdStep(t *testing.T) {
	q, err := Parse("g.V().id()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[1].Name != "id" {
		t.Errorf("expected id, got %q", s.Traversal.Steps[1].Name)
	}
}

func TestParser_Fold(t *testing.T) {
	q, err := Parse("g.V().values('age').fold()")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	if s.Traversal.Steps[2].Name != "fold" {
		t.Errorf("expected fold, got %q", s.Traversal.Steps[2].Name)
	}
}

func TestParser_RangeStep(t *testing.T) {
	q, err := Parse("g.V().range(0,10)")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	s := q.Statements[0]
	rangeStep := s.Traversal.Steps[1]
	if rangeStep.Name != "range" || rangeStep.Args[0].Int != 0 || rangeStep.Args[1].Int != 10 {
		t.Errorf("expected range(0,10), got %v", rangeStep)
	}
}

func TestParser_ErrorNoG(t *testing.T) {
	_, err := Parse("V().out()")
	if err == nil {
		t.Fatal("expected error for query not starting with g")
	}
}

func TestParser_EmptyQuery(t *testing.T) {
	q, err := Parse("")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(q.Statements) != 0 {
		t.Errorf("expected 0 statements, got %d", len(q.Statements))
	}
}
