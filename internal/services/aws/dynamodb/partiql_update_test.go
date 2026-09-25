package dynamodb

import (
	"errors"
	"strings"
	"testing"

	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func guardTestTable() *dbstore.Table {
	return &dbstore.Table{
		Name: "t",
		KeySchema: []*dbstore.KeySchemaElement{
			{AttributeName: "pk", KeyType: dbstore.KeyTypeHash},
			{AttributeName: "sk", KeyType: dbstore.KeyTypeRange},
		},
	}
}

func TestUpdateClauseTargetNames(t *testing.T) {
	_, clauses, _, err := parseUpdateStatement(`UPDATE t SET a = 1, b = 'x' REMOVE extra ADD n 5 DELETE tags 'v' WHERE pk = 'k1'`)
	if err != nil {
		t.Fatalf("multi-clause statement rejected: %v", err)
	}
	got := updateClauseTargetNames(clauses)
	want := []string{"a", "b", "extra", "n", "tags"}
	if len(got) != len(want) {
		t.Fatalf("collected targets %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("target %d = %q, want %q (all: %v)", i, got[i], want[i], got)
		}
	}

	// A nested target roots at its top-level attribute: the key-attribute
	// guard and the returning images report the attribute the write lands
	// in, never the literal dotted key.
	_, nested, _, err := parseUpdateStatement(`UPDATE t SET detail.name = 'v' REMOVE detail.tags[0] WHERE pk = 'k1'`)
	if err != nil {
		t.Fatalf("nested-target statement rejected: %v", err)
	}
	got = updateClauseTargetNames(nested)
	want = []string{"detail", "detail"}
	if len(got) != len(want) || got[0] != "detail" || got[1] != "detail" {
		t.Fatalf("nested targets %v, want %v", got, want)
	}
}

// A repeated action keyword keeps every clause in statement order: the
// reference's own UPDATE example carries two SET clauses, and a later
// clause must never overwrite an earlier one's assignments.
func TestManualUpdateKeepsEveryRepeatedClauseKeyword(t *testing.T) {
	_, clauses, _, err := parseUpdateStatement(`UPDATE t SET a = 1 SET b = 'x' WHERE pk = 'k1'`)
	if err != nil {
		t.Fatalf("two-SET statement rejected: %v", err)
	}
	if len(clauses.setAssignments) != 2 {
		t.Fatalf("expected both SET clauses parsed, got %d assignments (%v)",
			len(clauses.setAssignments), clauses.setAssignments)
	}
	if clauses.setAssignments[0].attrName != "a" || clauses.setAssignments[1].attrName != "b" {
		t.Fatalf("assignment order lost: first=%q second=%q",
			clauses.setAssignments[0].attrName, clauses.setAssignments[1].attrName)
	}
}

func TestPartiQLUpdateRejectsKeyAttributeTargets(t *testing.T) {
	table := guardTestTable()

	cases := []struct {
		name      string
		statement string
	}{
		{"SET partition key", `UPDATE t SET pk = 'other' WHERE pk = 'k1' AND sk = 's1'`},
		{"REMOVE sort key", `UPDATE t SET a = 1 REMOVE sk WHERE pk = 'k1' AND sk = 's1'`},
		{"ADD partition key", `UPDATE t SET a = 1 ADD pk 5 WHERE pk = 'k1' AND sk = 's1'`},
	}
	for _, tc := range cases {
		_, clauses, _, _ := parseUpdateStatement(tc.statement)
		err := validateNotKeyAttributes(table, updateClauseTargetNames(clauses))
		if err == nil {
			t.Fatalf("%s: expected rejection", tc.name)
		}
		apiErr, ok := err.(*APIError)
		if !ok || apiErr.Code != "com.amazon.coral.validate#ValidationException" {
			t.Fatalf("%s: expected ValidationException, got %v", tc.name, err)
		}
	}

	// A statement whose clauses touch only non-key attributes passes the guard.
	_, clauses, _, _ := parseUpdateStatement(`UPDATE t SET a = 1 REMOVE extra WHERE pk = 'k1' AND sk = 's1'`)
	if err := validateNotKeyAttributes(table, updateClauseTargetNames(clauses)); err != nil {
		t.Fatalf("non-key statement rejected: %v", err)
	}
}

// The RETURNING clause is part of the UPDATE grammar: the value validates
// against the documented four-value set (case-insensitively) and reaches
// the clause inventory in canonical form; anything else fails the
// statement, never executes with the requested return dropped.
func TestUpdateReturningClauseParsesAndValidates(t *testing.T) {
	_, clauses, _, err := parseUpdateStatement(`UPDATE t SET a = 1 WHERE pk = 'k1' RETURNING all old *`)
	if err != nil {
		t.Fatalf("RETURNING statement rejected: %v", err)
	}
	if clauses.returning != "ALL OLD *" {
		t.Fatalf("canonical returning = %q, want %q", clauses.returning, "ALL OLD *")
	}
	for _, bad := range []string{
		`UPDATE t SET a = 1 WHERE pk = 'k1' RETURNING ALL OLD`,    // missing star
		`UPDATE t SET a = 1 WHERE pk = 'k1' RETURNING EVERYTHING`, // outside the value set
	} {
		if _, _, _, err := parseUpdateStatement(bad); err == nil {
			t.Fatalf("invalid RETURNING value accepted: %s", bad)
		}
	}
}

// The DELETE grammar carries one optional RETURNING clause whose only
// documented value is ALL OLD *; a different value fails the statement,
// and a quoted literal never reads as the clause keyword.
func TestDeleteReturningClauseAllOldOnly(t *testing.T) {
	tableName, _, returning := parseDeleteStatement(`DELETE FROM t WHERE pk = 'k1' RETURNING ALL OLD *`)
	if tableName != "t" || returning != "ALL OLD *" {
		t.Fatalf("DELETE with RETURNING ALL OLD * = (%q, %q)", tableName, returning)
	}
	for _, bad := range []string{
		`DELETE FROM t WHERE pk = 'k1' RETURNING ALL NEW *`,
		`DELETE FROM t WHERE pk = 'k1' RETURNING MODIFIED OLD *`,
	} {
		if tableName, _, _ := parseDeleteStatement(bad); tableName != "" {
			t.Fatalf("DELETE returning value outside ALL OLD * accepted: %s", bad)
		}
	}
	if tableName, _, returning := parseDeleteStatement(`DELETE FROM t WHERE note = 'RETURNING ALL OLD *'`); tableName != "t" || returning != "" {
		t.Fatalf("quoted RETURNING misread as the clause: (%q, %q)", tableName, returning)
	}
}

// Every parse failure inside the manual UPDATE path fails the statement
// instead of executing it with partial intent: a dropped clause value, a
// value-less ADD target, or a WHERE that fails its reparse (which would
// otherwise run the update with no filter at all) are all validation
// errors, and the parse returns an empty table name so lock-key
// derivation locks nothing for a statement the engine is about to reject.
func TestManualUpdateParseFailuresFailTheStatement(t *testing.T) {
	cases := []struct {
		name      string
		statement string
	}{
		{"ADD value that does not parse", `UPDATE t SET a = 1 ADD n ) WHERE pk = 'k1'`},
		{"DELETE value that does not parse", `UPDATE t SET a = 1 DELETE tags ) WHERE pk = 'k1'`},
		{"ADD target without a value", `UPDATE t SET a = 1 ADD n WHERE pk = 'k1'`},
		{"SET value that does not parse", `UPDATE t SET a = ) WHERE pk = 'k1'`},
		{"WHERE that does not reparse", `UPDATE t SET a = 1 WHERE pk = )`},
	}
	for _, tc := range cases {
		tableName, _, _, err := parseUpdateStatement(tc.statement)
		if err == nil {
			t.Fatalf("%s: expected the statement to fail", tc.name)
		}
		if !errors.Is(err, ErrInvalidParameter) {
			t.Fatalf("%s: expected ErrInvalidParameter, got %v", tc.name, err)
		}
		if tableName != "" {
			t.Fatalf("%s: rejected statement must derive no table name, got %q", tc.name, tableName)
		}
	}
}

func matchSAttr(v string) *dbstore.AttributeValue { return &dbstore.AttributeValue{S: &v} }

// The ADD/DELETE action appliers are one implementation serving both the
// expression plane and the PartiQL clause plane. The pins below fix the two
// divergences the consolidation registered: every branch assigns a fresh
// attribute value (a snapshot taken before the application cannot alias
// the updated item); a type-incompatible pair answers ErrTypeMismatch,
// the shared contract of every calling plane.
func TestAddDeleteActionAppliersAssignFreshValues(t *testing.T) {
	attrs := map[string]*dbstore.AttributeValue{
		"tags": dbstore.StringSet([]string{"a", "b"}),
	}
	snapshot := attrs["tags"]

	changed, err := applyAddAction(attrs, attrs, "tags", dbstore.StringSet([]string{"c"}))
	if err != nil || !changed {
		t.Fatalf("ADD: changed=%v err=%v", changed, err)
	}
	if attrs["tags"] == snapshot {
		t.Fatal("ADD must assign a fresh attribute value, not mutate the stored one")
	}
	if len(snapshot.SS) != 2 {
		t.Fatalf("snapshot aliased the ADD: members = %v", snapshot.SS)
	}
	if len(attrs["tags"].SS) != 3 {
		t.Fatalf("updated set members = %v, want three", attrs["tags"].SS)
	}

	snapshot = attrs["tags"]
	changed, err = applyDeleteAction(attrs, attrs, "tags", dbstore.StringSet([]string{"a"}))
	if err != nil || !changed {
		t.Fatalf("DELETE: changed=%v err=%v", changed, err)
	}
	if attrs["tags"] == snapshot {
		t.Fatal("DELETE must assign a fresh attribute value, not mutate the stored one")
	}
	if len(snapshot.SS) != 3 || len(attrs["tags"].SS) != 2 {
		t.Fatalf("DELETE: snapshot=%v updated=%v", snapshot.SS, attrs["tags"].SS)
	}

	// Emptying the set removes the attribute entirely.
	if changed, err = applyDeleteAction(attrs, attrs, "tags", dbstore.StringSet([]string{"b", "c"})); err != nil || !changed {
		t.Fatalf("emptying DELETE: changed=%v err=%v", changed, err)
	}
	if _, exists := attrs["tags"]; exists {
		t.Fatal("a set emptied by DELETE must be removed from the item")
	}
}

func TestAddDeleteMismatchAnswersTheUnifiedError(t *testing.T) {
	attrs := map[string]*dbstore.AttributeValue{
		"n":  dbstore.NumberValue("1"),
		"ss": dbstore.StringSet([]string{"a"}),
	}

	// A type-incompatible pair answers ErrTypeMismatch on both appliers —
	// the documented ValidationException wording, identical on every
	// update face.
	if _, err := applyAddAction(attrs, attrs, "n", dbstore.StringSet([]string{"x"})); err != ErrTypeMismatch {
		t.Fatalf("ADD mismatch: %v", err)
	}
	if _, err := applyDeleteAction(attrs, attrs, "ss", dbstore.NumberValue("1")); err != ErrTypeMismatch {
		t.Fatalf("DELETE mismatch: %v", err)
	}

	// An ADD onto an absent attribute creates it from the operand; a
	// DELETE onto an absent attribute is a no-op.
	empty := map[string]*dbstore.AttributeValue{}
	if changed, err := applyAddAction(empty, empty, "n", dbstore.NumberValue("3")); err != nil || !changed {
		t.Fatalf("creating ADD: changed=%v err=%v", changed, err)
	}
	if changed, err := applyDeleteAction(empty, empty, "ss", dbstore.NumberSet([]string{"3"})); err != nil || changed {
		t.Fatalf("absent DELETE: changed=%v err=%v", changed, err)
	}
	if _, exists := empty["ss"]; exists {
		t.Fatal("absent DELETE must not create the attribute")
	}
}

// TestPartiQLColumnReferenceRendersEveryPathSegment pins the read-side
// rendering of qualified column references: the grammar nests one
// qualifier level, so a three-segment reference is a single ColName whose
// leading segment the renderer must keep — the write side addresses
// three-segment document paths, and the read side must name the same
// paths, not their two-segment shadows.
func TestPartiQLColumnReferenceRendersEveryPathSegment(t *testing.T) {
	ss := func(members ...string) *dbstore.AttributeValue { return dbstore.StringSet(members) }
	deepMap := func(inner *dbstore.AttributeValue) *dbstore.AttributeValue {
		return &dbstore.AttributeValue{M: map[string]*dbstore.AttributeValue{
			"b": &dbstore.AttributeValue{M: map[string]*dbstore.AttributeValue{"c": inner}},
		}}
	}
	// Each case starts from a fresh item: the cases write the same paths,
	// and a shared map would carry one case's fold into the next.
	freshAttrs := func() map[string]*dbstore.AttributeValue {
		return map[string]*dbstore.AttributeValue{
			"a": deepMap(ss("x")),
			// The head-dropping render reads a.b.c's two-segment shadow:
			// seed one so the wrong read is observable, not just absent.
			"b": &dbstore.AttributeValue{M: map[string]*dbstore.AttributeValue{"c": ss("wrong")}},
		}
	}
	apply := func(attrs map[string]*dbstore.AttributeValue, stmt string) error {
		_, clauses, _, _ := parseUpdateStatement(stmt)
		return applySetAssignments(attrs, copyAttributes(attrs), clauses.setAssignments, nil)
	}

	// The documented same-path set function form at depth three folds its
	// own target — the full path on both sides of the correspondence.
	attrs := freshAttrs()
	if err := apply(attrs, `UPDATE t SET "a"."b"."c" = set_add("a"."b"."c", <<'y'>>) WHERE pk = '1'`); err != nil {
		t.Fatalf("three-segment same-path set_add: %v", err)
	}
	if got := attrs["a"].M["b"].M["c"].SS; len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("three-segment set_add: a.b.c = %v, want [x y]", got)
	}

	// A bare-path operand at depth three reads the attribute it names, not
	// the two-segment shadow under the same tail.
	attrs = freshAttrs()
	if err := apply(attrs, `UPDATE t SET y = a.b.c WHERE pk = '1'`); err != nil {
		t.Fatalf("three-segment copy: %v", err)
	}
	if got := attrs["y"].SS; len(got) != 1 || got[0] != "x" {
		t.Fatalf("three-segment copy: y = %v, want [x] (a.b.c, not the b.c shadow)", got)
	}

	// The two-segment form keeps rendering as it always has.
	if err := apply(attrs, `UPDATE t SET y2 = b.c WHERE pk = '1'`); err != nil {
		t.Fatalf("two-segment copy: %v", err)
	}
	if got := attrs["y2"].SS; len(got) != 1 || got[0] != "wrong" {
		t.Fatalf("two-segment copy: y2 = %v, want [wrong] (b.c itself)", got)
	}
}

// TestPartiQLClauseTargetsKeepRawSegmentNames pins the write-side rendering
// of clause targets: a segment the identifier printer would escape — an SQL
// keyword or a non-ASCII name — keeps its raw name on every clause face.
// The read side renders raw values; a target carrying an embedded escape
// would write a ghost member instead of the attribute it names, and the set
// functions' same-path correspondence would falsely reject the documented
// form.
func TestPartiQLClauseTargetsKeepRawSegmentNames(t *testing.T) {
	ss := func(members ...string) *dbstore.AttributeValue { return dbstore.StringSet(members) }
	freshAttrs := func() map[string]*dbstore.AttributeValue {
		return map[string]*dbstore.AttributeValue{
			"pk": dbstore.StringValue("1"),
			"m":  {M: map[string]*dbstore.AttributeValue{"add": ss("x")}},
		}
	}
	// applyAll mirrors the engine's clause application order: one statement,
	// one pre-statement image every read shares.
	applyAll := func(attrs map[string]*dbstore.AttributeValue, stmt string) error {
		_, clauses, _, _ := parseUpdateStatement(stmt)
		before := copyAttributes(attrs)
		if err := applySetAssignments(attrs, before, clauses.setAssignments, nil); err != nil {
			return err
		}
		if err := applyRemoveAttrs(attrs, clauses.removeAttrs); err != nil {
			return err
		}
		if err := applyAddAssignments(attrs, before, clauses.addAssignments, nil); err != nil {
			return err
		}
		return applyDeleteAssignments(attrs, before, clauses.deleteAssignments, nil)
	}

	// A qualified target whose member name is an SQL keyword writes the
	// member it names — never a ghost under the printer-escaped name.
	attrs := freshAttrs()
	if err := applyAll(attrs, `UPDATE t SET "m"."add" = <<'y'>> WHERE pk = '1'`); err != nil {
		t.Fatalf("keyword member overwrite: %v", err)
	}
	if inner := attrs["m"].M; len(inner) != 1 {
		t.Fatalf("keyword member overwrite: m holds %d members, want 1 (no ghost)", len(inner))
	} else if got := inner["add"].SS; len(got) != 1 || got[0] != "y" {
		t.Fatalf("keyword member overwrite: m.add = %v, want [y]", got)
	}

	// The same-path set-function form on a keyword member folds its target:
	// both sides of the correspondence render the same raw path.
	attrs = freshAttrs()
	if err := applyAll(attrs, `UPDATE t SET "m"."add" = set_add("m"."add", <<'y'>>) WHERE pk = '1'`); err != nil {
		t.Fatalf("keyword same-path set_add: %v", err)
	}
	if got := attrs["m"].M["add"].SS; len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("keyword same-path set_add: m.add = %v, want [x y]", got)
	}

	// A non-ASCII member name is escaped by the same printer rule and must
	// likewise keep its raw name.
	attrs = map[string]*dbstore.AttributeValue{
		"pk": dbstore.StringValue("1"),
		"m":  {M: map[string]*dbstore.AttributeValue{"年": ss("x")}},
	}
	if err := applyAll(attrs, `UPDATE t SET "m"."年" = <<'y'>> WHERE pk = '1'`); err != nil {
		t.Fatalf("non-ASCII member overwrite: %v", err)
	}
	if inner := attrs["m"].M; len(inner) != 1 {
		t.Fatalf("non-ASCII member overwrite: m holds %d members, want 1 (no ghost)", len(inner))
	} else if got := inner["年"].SS; len(got) != 1 || got[0] != "y" {
		t.Fatalf("non-ASCII member overwrite: m.年 = %v, want [y]", got)
	}

	// The ADD and REMOVE clause faces share the rendering: a qualified
	// keyword target folds into and removes the member it names.
	attrs = freshAttrs()
	if err := applyAll(attrs, `UPDATE t ADD "m"."add" <<'y'>> WHERE pk = '1'`); err != nil {
		t.Fatalf("keyword ADD: %v", err)
	}
	if got := attrs["m"].M["add"].SS; len(got) != 2 || got[0] != "x" || got[1] != "y" {
		t.Fatalf("keyword ADD: m.add = %v, want [x y]", got)
	}
	attrs = freshAttrs()
	if err := applyAll(attrs, `UPDATE t REMOVE "m"."add" WHERE pk = '1'`); err != nil {
		t.Fatalf("keyword REMOVE: %v", err)
	}
	if _, exists := attrs["m"].M["add"]; exists {
		t.Fatal("keyword REMOVE: m.add must be removed")
	}

	// A qualified target without escapable segments keeps its rendering.
	attrs = map[string]*dbstore.AttributeValue{
		"pk": dbstore.StringValue("1"),
		"m":  {M: map[string]*dbstore.AttributeValue{"n": ss("old")}},
	}
	if err := applyAll(attrs, `UPDATE t SET "m"."n" = <<'new'>> WHERE pk = '1'`); err != nil {
		t.Fatalf("plain qualified overwrite: %v", err)
	}
	if got := attrs["m"].M["n"].SS; len(got) != 1 || got[0] != "new" {
		t.Fatalf("plain qualified overwrite: m.n = %v, want [new]", got)
	}
}

// TestPartiQLUpdateEvaluatesAgainstThePreUpdateItem pins the evaluation
// basis of the PartiQL UPDATE plane: one statement is one update of one
// item, and every assignment's operand reads the item as it was before
// the statement — the basis the update-expression plane documents
// ("DynamoDB evaluates every action against the item's attribute values
// as they were before the update. The actions aren't applied one after
// another from left to right") — never the partially-updated map the
// statement's own earlier assignments have written.
func TestPartiQLUpdateEvaluatesAgainstThePreUpdateItem(t *testing.T) {
	num := func(v string) *dbstore.AttributeValue { return dbstore.NumberValue(v) }
	// applyOneStatement mirrors the engine's clause application for SET and
	// REMOVE: one statement, one pre-statement image every read shares.
	applyOneStatement := func(attrs map[string]*dbstore.AttributeValue, stmt string, params *partiQLParams) error {
		_, clauses, _, _ := parseUpdateStatement(stmt)
		before := copyAttributes(attrs)
		if err := applySetAssignments(attrs, before, clauses.setAssignments, params); err != nil {
			return err
		}
		return applyRemoveAttrs(attrs, clauses.removeAttrs)
	}

	// A chained copy reads both sources from the pre-statement item: over
	// {a:1,b:2,c:3}, b = a and c = b yield {b:1,c:2}.
	item := map[string]*dbstore.AttributeValue{"a": num("1"), "b": num("2"), "c": num("3")}
	if err := applyOneStatement(item, `UPDATE t SET b = a, c = b WHERE pk = 'k1'`, nil); err != nil {
		t.Fatalf("chained copy: %v", err)
	}
	if got := *item["b"].N; got != "1" {
		t.Fatalf("chained copy: b = %q, want 1", got)
	}
	if got := *item["c"].N; got != "2" {
		t.Fatalf("chained copy: c = %q, want 2 (the pre-statement b)", got)
	}

	// The two-clause form of the same statement — the documented example
	// shape: a swap reads both operands before either write lands.
	swap := map[string]*dbstore.AttributeValue{"a": num("1"), "b": num("2")}
	if err := applyOneStatement(swap, `UPDATE t SET b = a SET a = b WHERE pk = 'k1'`, nil); err != nil {
		t.Fatalf("swap: %v", err)
	}
	if got := *swap["a"].N; got != "2" || *swap["b"].N != "1" {
		t.Fatalf("swap: got a=%s b=%s, want a=2 b=1", *swap["a"].N, *swap["b"].N)
	}

	// The update-expression documentation's worked example in the PartiQL
	// form: the removal erases nothing the operands read.
	worked := map[string]*dbstore.AttributeValue{"a": num("1"), "b": num("2"), "c": num("3")}
	if err := applyOneStatement(worked, `UPDATE t REMOVE a SET b = a, c = b WHERE pk = 'k1'`, nil); err != nil {
		t.Fatalf("worked example: %v", err)
	}
	if _, exists := worked["a"]; exists {
		t.Fatal("worked example: a must be removed")
	}
	if got := *worked["b"].N; got != "1" || *worked["c"].N != "2" {
		t.Fatalf("worked example: got b=%s c=%s, want b=1 c=2", *worked["b"].N, *worked["c"].N)
	}

	// A SET function's fold shares the basis: n gains the operand while
	// the later copy still reads the pre-statement n.
	fold := map[string]*dbstore.AttributeValue{"n": num("1")}
	params := wireParams(map[string]interface{}{"N": "5"})
	if err := applyOneStatement(fold, `UPDATE t SET n = SET_ADD(n, :v1) SET mirror = n WHERE pk = 'k1'`, params); err != nil {
		t.Fatalf("set function basis: %v", err)
	}
	if got := *fold["n"].N; got != "6" {
		t.Fatalf("set function basis: n = %q, want 6", got)
	}
	if got := *fold["mirror"].N; got != "1" {
		t.Fatalf("set function basis: mirror = %q, want 1 (the pre-statement n)", got)
	}
}

// The SET functions' first argument must address the assignment's own
// target: the only documented form is set_add(path, value) with the path
// naming the target the fold writes, and a first argument naming any
// other path fits no documented form — the statement fails rather than
// silently folding the value into the target while the named path is
// never read.
func TestPartiQLSetFunctionRequiresItsTargetPath(t *testing.T) {
	attrs := map[string]*dbstore.AttributeValue{
		"sizes": dbstore.StringSet([]string{"M"}),
	}
	apply := func(stmt string) error {
		_, clauses, _, _ := parseUpdateStatement(stmt)
		return applySetAssignments(attrs, copyAttributes(attrs), clauses.setAssignments, nil)
	}
	if err := apply(`UPDATE t SET colors = set_add(sizes, <<'L'>>) WHERE pk = 'k1'`); err != ErrInvalidParameter {
		t.Fatalf("set_add onto another path: %v, want ErrInvalidParameter", err)
	}
	if err := apply(`UPDATE t SET colors = set_delete(sizes, <<'M'>>) WHERE pk = 'k1'`); err != ErrInvalidParameter {
		t.Fatalf("set_delete onto another path: %v, want ErrInvalidParameter", err)
	}
	if _, exists := attrs["colors"]; exists {
		t.Fatal("a mismatched set function must not create or write its target")
	}

	// The same-path documented form keeps folding into its target.
	if err := apply(`UPDATE t SET sizes = set_add(sizes, <<'L'>>) WHERE pk = 'k1'`); err != nil {
		t.Fatalf("same-path set_add: %v", err)
	}
	if got := attrs["sizes"].SS; len(got) != 2 {
		t.Fatalf("same-path set_add: sizes = %v, want both members", got)
	}
}

// The PartiQL SET plane evaluates through the shared operand model: the
// two SET functions and column-reference operands behave exactly as the
// UpdateExpression plane's grammar defines them, and the model's type
// mismatch reaches the statement error unchanged.
func TestPartiQLSetAssignmentsUseTheSharedOperandModel(t *testing.T) {
	str := func(v string) *dbstore.AttributeValue { return &dbstore.AttributeValue{S: &v} }
	attrs := map[string]*dbstore.AttributeValue{
		"keep": str("kept"),
		"lst":  dbstore.ListValue([]*dbstore.AttributeValue{str("a")}),
	}
	apply := func(stmt string) error {
		_, clauses, _, _ := parseUpdateStatement(stmt)
		return applySetAssignments(attrs, copyAttributes(attrs), clauses.setAssignments, nil)
	}

	// if_not_exists keeps the existing attribute the path addresses...
	if err := apply(`UPDATE t SET v = if_not_exists(keep, 'd') WHERE pk = 'k1'`); err != nil {
		t.Fatalf("if_not_exists keep: %v", err)
	}
	if got := *attrs["v"].S; got != "kept" {
		t.Fatalf("if_not_exists keep: v = %q, want kept", got)
	}

	// ...and takes the fallback when the path addresses nothing.
	if err := apply(`UPDATE t SET w = if_not_exists(absent, 'd') WHERE pk = 'k1'`); err != nil {
		t.Fatalf("if_not_exists fallback: %v", err)
	}
	if got := *attrs["w"].S; got != "d" {
		t.Fatalf("if_not_exists fallback: w = %q, want d", got)
	}

	// A bare column reference is a document-path operand: the assignment
	// copies the addressed value, the operand grammar the UpdateExpression
	// plane has always carried.
	if err := apply(`UPDATE t SET copy = keep WHERE pk = 'k1'`); err != nil {
		t.Fatalf("column copy: %v", err)
	}
	if got := *attrs["copy"].S; got != "kept" {
		t.Fatalf("column copy: copy = %q, want kept", got)
	}

	// list_append concatenates the item path's list with a literal list.
	if err := apply(`UPDATE t SET x = list_append(lst, ['b']) WHERE pk = 'k1'`); err != nil {
		t.Fatalf("list_append: %v", err)
	}
	if l := attrs["x"].L; len(l) != 2 || *l[0].S != "a" || *l[1].S != "b" {
		t.Fatalf("list_append: x = %+v, want [a b]", l)
	}

	// A non-list operand is the shared model's type mismatch — the same
	// error the UpdateExpression plane answers.
	if err := apply(`UPDATE t SET y = list_append(keep, ['b']) WHERE pk = 'k1'`); err != ErrTypeMismatch {
		t.Fatalf("list_append non-list operand: %v, want ErrTypeMismatch", err)
	}

	// The SET set functions fold their value operand into the assignment's
	// own target through the ADD and DELETE action appliers: set_add unions
	// a set member in (or adds numerically), set_delete removes one — and
	// deletes the attribute when the subtraction empties it — while a
	// type-incompatible pair answers the plane's mismatch sentinel.
	attrs["tags"] = dbstore.StringSet([]string{"a"})
	attrs["n"] = dbstore.NumberValue("10")
	setParams := wireParams(
		map[string]interface{}{"SS": []interface{}{"b"}},
		map[string]interface{}{"N": "5"},
		map[string]interface{}{"SS": []interface{}{"a"}},
	)
	applyParams := func(stmt string) error {
		_, clauses, _, pErr := parseUpdateStatement(stmt)
		if pErr != nil {
			return pErr
		}
		return applySetAssignments(attrs, copyAttributes(attrs), clauses.setAssignments, setParams)
	}
	if err := applyParams(`UPDATE t SET tags = SET_ADD(tags, :v1) WHERE pk = 'k1'`); err != nil {
		t.Fatalf("set_add: %v", err)
	}
	if got := attrs["tags"].SS; len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("set_add: tags = %v, want [a b]", got)
	}
	if err := applyParams(`UPDATE t SET n = SET_ADD(n, :v2) WHERE pk = 'k1'`); err != nil {
		t.Fatalf("set_add number: %v", err)
	}
	if got := *attrs["n"].N; got != "15" {
		t.Fatalf("set_add number: n = %q, want 15", got)
	}
	if err := applyParams(`UPDATE t SET tags = SET_ADD(tags, :v2) WHERE pk = 'k1'`); err != ErrTypeMismatch {
		t.Fatalf("set_add number into string set: %v, want ErrTypeMismatch", err)
	}
	if err := applyParams(`UPDATE t SET tags = SET_DELETE(tags, :v3) WHERE pk = 'k1'`); err != nil {
		t.Fatalf("set_delete: %v", err)
	}
	if got := attrs["tags"].SS; len(got) != 1 || got[0] != "b" {
		t.Fatalf("set_delete: tags = %v, want [b]", got)
	}
	if err := applyParams(`UPDATE t SET keep = SET_DELETE(keep, :v3) WHERE pk = 'k1'`); err != ErrTypeMismatch {
		t.Fatalf("set_delete from a non-set attribute: %v, want ErrTypeMismatch", err)
	}

	// A SET operand naming an unknown function is a validation error, never
	// a stored value.
	before := len(attrs)
	if err := apply(`UPDATE t SET z = no_such_function(1) WHERE pk = 'k1'`); err != ErrInvalidParameter {
		t.Fatalf("unknown SET function: %v, want ErrInvalidParameter", err)
	}
	if len(attrs) != before {
		t.Fatalf("unknown SET function must store nothing; attrs grew to %d", len(attrs))
	}

	// A nested left-hand side writes the member it names — the assignment
	// target parses as the document path it is, never a literal flat key —
	// and a missing parent is the documented invalid document path.
	attrs["detail"] = dbstore.MapValue(map[string]*dbstore.AttributeValue{
		"name": str("old"),
		"tags": dbstore.StringSet([]string{"a"}),
	})
	if err := apply(`UPDATE t SET detail.name = 'new' WHERE pk = 'k1'`); err != nil {
		t.Fatalf("nested SET: %v", err)
	}
	if got := *attrs["detail"].M["name"].S; got != "new" {
		t.Fatalf("nested SET: detail.name = %q, want new", got)
	}
	if _, flat := attrs["detail.name"]; flat {
		t.Fatal("nested SET must not store a literal flat key")
	}
	if err := apply(`UPDATE t SET absent.nested = 'v' WHERE pk = 'k1'`); err == nil || !strings.Contains(err.Error(), "invalid for update") {
		t.Fatalf("nested SET with a missing parent: %v, want the invalid-document-path rejection", err)
	}

	// The set functions fold at a nested target through the same walks, and
	// REMOVE resolves its target as a document path too — an indexed member
	// of a nested list.
	if err := applyParams(`UPDATE t SET detail.tags = SET_ADD(detail.tags, :v1) WHERE pk = 'k1'`); err != nil {
		t.Fatalf("nested set_add: %v", err)
	}
	if got := attrs["detail"].M["tags"].SS; len(got) != 2 || got[1] != "b" {
		t.Fatalf("nested set_add: detail.tags = %v, want [a b]", got)
	}
	attrs["detail"].M["nums"] = dbstore.ListValue([]*dbstore.AttributeValue{
		dbstore.NumberValue("10"), dbstore.NumberValue("20"), dbstore.NumberValue("30"),
	})
	_, clauses, _, pErr := parseUpdateStatement(`UPDATE t SET a = 1 REMOVE detail.nums[1] WHERE pk = 'k1'`)
	if pErr != nil {
		t.Fatalf("nested REMOVE parse: %v", pErr)
	}
	if err := applyRemoveAttrs(attrs, clauses.removeAttrs); err != nil {
		t.Fatalf("nested REMOVE apply: %v", err)
	}
	remaining := attrs["detail"].M["nums"].L
	if len(remaining) != 2 || *remaining[0].N != "10" || *remaining[1].N != "30" {
		t.Fatalf("nested REMOVE: detail.nums = %v, want [10 30]", remaining)
	}
}

func TestMatchKeyedTargetItemEnforcesSinglePrimaryKeyContract(t *testing.T) {
	table := guardTestTable()
	mkItem := func(sk, v string) *dbstore.Item {
		return &dbstore.Item{Attributes: map[string]*dbstore.AttributeValue{
			"pk": matchSAttr("a"),
			"sk": matchSAttr(sk),
			"v":  matchSAttr(v),
		}}
	}
	items := []*dbstore.Item{mkItem("s1", "x"), mkItem("s2", "y")}

	// The getItem adapter emulates a key lookup over the stored items.
	run := func(where string) (item *dbstore.Item, found, conditionHeld bool, err error) {
		_, whereExpr := parseSelectStatement("SELECT * FROM t WHERE " + where)
		getItem := func(key map[string]*dbstore.AttributeValue) (*dbstore.Item, error) {
			for _, it := range items {
				if *it.Attributes["pk"].S == *key["pk"].S && *it.Attributes["sk"].S == *key["sk"].S {
					return it, nil
				}
			}
			return nil, dbstore.ErrItemNotFound
		}
		return matchKeyedTargetItem(getItem, table, whereExpr, nil, "UPDATE")
	}

	// A nil WHERE clause cannot resolve a primary key.
	if _, _, _, err := run(""); err == nil {
		t.Fatal("empty WHERE text: expected rejection")
	}

	// A WHERE naming only the partition key of a composite-key table does
	// not resolve to a single primary-key value.
	_, _, _, err := run("pk = 'a'")
	apiErr, ok := err.(*APIError)
	if !ok || apiErr.Code != "com.amazon.coral.validate#ValidationException" {
		t.Fatalf("partition-only WHERE: expected ValidationException, got %v", err)
	}

	// A non-key WHERE is rejected the same way.
	if _, _, _, err = run("v = 'x'"); err == nil {
		t.Fatal("non-key WHERE: expected rejection")
	}

	// A full-key WHERE whose item exists and satisfies the remaining
	// predicates resolves the item.
	item, found, conditionHeld, err := run("pk = 'a' AND sk = 's2'")
	if err != nil || item == nil || !found || !conditionHeld {
		t.Fatalf("full-key match: expected found item holding the condition, got item=%v found=%v held=%v err=%v", item, found, conditionHeld, err)
	}
	if got := *item.Attributes["sk"].S; got != "s2" {
		t.Fatalf("full-key match: expected sk=s2, got %q", got)
	}

	// A full-key WHERE addressing no stored item reports the absent key.
	item, found, conditionHeld, err = run("pk = 'a' AND sk = 'absent'")
	if err != nil || item != nil || found || conditionHeld {
		t.Fatalf("absent key: expected no item and no error, got item=%v found=%v held=%v err=%v", item, found, conditionHeld, err)
	}

	// A full-key WHERE whose item exists but fails a further predicate
	// reports the item with the condition rejected.
	item, found, conditionHeld, err = run("pk = 'a' AND sk = 's1' AND v = 'nope'")
	if err != nil || item == nil || !found || conditionHeld {
		t.Fatalf("rejected condition: expected the item with the condition rejected, got item=%v found=%v held=%v err=%v", item, found, conditionHeld, err)
	}
	if got := *item.Attributes["v"].S; got != "x" {
		t.Fatalf("rejected condition: expected the stored item, got v=%q", got)
	}
}
