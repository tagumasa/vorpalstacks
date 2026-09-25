// Package dynamodb provides DynamoDB service operations for vorpalstacks.
package dynamodb

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vorpalstacks/internal/common/request"
	dbstore "vorpalstacks/internal/store/aws/dynamodb"
)

func TestUpdateExpressionStrictness(t *testing.T) {
	values := map[string]*dbstore.AttributeValue{
		":v": dbstore.NumberValue("2"),
		":w": dbstore.NumberValue("3"),
	}

	attrs := map[string]*dbstore.AttributeValue{
		"src": dbstore.StringValue("copied"),
	}

	// The multiply operator is not part of SET arithmetic.
	_, err := applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :v * :w", nil, values)
	assert.Error(t, err)

	// An undefined value placeholder must not silently skip the action.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :missing", nil, values)
	assert.Error(t, err)

	// An undefined name placeholder must not pass through as a literal.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET #missing = :v", nil, values)
	assert.Error(t, err)

	// The same rule on the operand side: an undefined name placeholder
	// inside a SET operand is a validation error, never a literal path
	// that silently addresses nothing and drops the operand to a fallback.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = if_not_exists(#nope, :v)", nil, values)
	assert.Error(t, err)

	// SET a = src copies the existing attribute value.
	updated, err := applyUpdateExpressionWithTracking(attrs, "SET a = src", nil, nil)
	require.NoError(t, err)
	assert.Contains(t, updated, "a")
	assert.Equal(t, "copied", *attrs["a"].S)

	// DELETE with a mismatched operand type is rejected.
	ssAttrs := map[string]*dbstore.AttributeValue{
		"tags": dbstore.StringSet([]string{"a", "b"}),
	}
	_, err = applyUpdateExpressionWithTracking(ssAttrs, "DELETE tags :nums", nil, map[string]*dbstore.AttributeValue{
		":nums": dbstore.NumberSet([]string{"1"}),
	})
	assert.Error(t, err)

	// DELETE of matching elements still works.
	updated, err = applyUpdateExpressionWithTracking(ssAttrs, "DELETE tags :del", nil, map[string]*dbstore.AttributeValue{
		":del": dbstore.StringSet([]string{"a"}),
	})
	require.NoError(t, err)
	assert.Contains(t, updated, "tags")
	assert.Equal(t, []string{"b"}, ssAttrs["tags"].SS)

	// An undefined name placeholder in REMOVE must not pass through as a
	// literal attribute name.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "REMOVE #missing", nil, nil)
	assert.Error(t, err)

	// REMOVE with a defined name placeholder removes the attribute.
	rmAttrs := map[string]*dbstore.AttributeValue{
		"gsik": dbstore.StringValue("g"),
	}
	updated, err = applyUpdateExpressionWithTracking(rmAttrs, "REMOVE #name", map[string]string{"#name": "gsik"}, nil)
	require.NoError(t, err)
	assert.Contains(t, updated, "gsik")
	assert.NotContains(t, rmAttrs, "gsik")

	// The comma between same-clause actions is grammar, not formatting: a
	// comma-less concatenation is a syntax error on every clause, never a
	// second silently-applied action.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :v b = :w", nil, values)
	assert.Error(t, err)
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "REMOVE a b", nil, nil)
	assert.Error(t, err)
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "ADD a :v b :w", nil, values)
	assert.Error(t, err)
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{"ss": dbstore.StringSet([]string{"x"})}, "DELETE ss :v ss2 :w", nil, map[string]*dbstore.AttributeValue{
		":v": dbstore.StringSet([]string{"x"}),
		":w": dbstore.StringSet([]string{"y"}),
	})
	assert.Error(t, err)

	// A clause keyword may appear only once: the second occurrence is a
	// syntax error, never a path to an attribute with the keyword's
	// spelling.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{"b": dbstore.StringValue("keep")}, "REMOVE a REMOVE b", nil, nil)
	assert.Error(t, err)

	// A separator sits BETWEEN two actions of one clause: a trailing
	// comma, or one directly before a clause keyword, leaves the grammar
	// unsatisfied and is a syntax error on every clause.
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "SET a = :v ,", nil, values)
	assert.Error(t, err)
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{"b": dbstore.StringValue("keep")}, "SET a = :v, REMOVE b", nil, values)
	assert.Error(t, err)
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "REMOVE a ,", nil, nil)
	assert.Error(t, err)
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{}, "ADD a :v ,", nil, values)
	assert.Error(t, err)
	_, err = applyUpdateExpressionWithTracking(map[string]*dbstore.AttributeValue{"ss": dbstore.StringSet([]string{"x"})}, "DELETE ss :v ,", nil, map[string]*dbstore.AttributeValue{
		":v": dbstore.StringSet([]string{"x"}),
	})
	assert.Error(t, err)

	// The comma-separated and clause-boundary forms keep working.
	multiAttrs := map[string]*dbstore.AttributeValue{"gone": dbstore.StringValue("g")}
	updated, err = applyUpdateExpressionWithTracking(multiAttrs, "SET a = :v, b = :w REMOVE gone", nil, values)
	require.NoError(t, err)
	assert.Contains(t, updated, "a")
	assert.Contains(t, updated, "b")
	assert.Contains(t, updated, "gone")
	assert.NotContains(t, multiAttrs, "gone")
}

// TestUpdateExpressionSeparatorRequiredAcrossFaces pins the grammar on the
// wire faces an UpdateExpression reaches: the comma between same-clause
// actions is required, so a comma-less concatenation answers the
// invalid-parameter rejection on UpdateItem and on a TransactWriteItems
// Update — never a second silently-applied action. The path-analysis walk
// the key-attribute guard runs shares the applier's grammar, so the
// rejection surfaces before any guard verdict.
func TestUpdateExpressionSeparatorRequiredAcrossFaces(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "SepTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "SepTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "n": map[string]interface{}{"N": "1"}},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}
	values := map[string]interface{}{":v": map[string]interface{}{"N": "2"}, ":w": map[string]interface{}{"N": "3"}}

	_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "SepTable", "Key": map[string]interface{}{"id": sVal("a")},
		"UpdateExpression":          "SET b = :v c = :w",
		"ExpressionAttributeValues": values,
	}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("UpdateItem face: err = %v, want ErrInvalidParameter", err)
	}

	_, err = svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"TableName": "SepTable", "Key": map[string]interface{}{"id": sVal("a")},
			"UpdateExpression":          "SET b = :v c = :w",
			"ExpressionAttributeValues": values,
		}}},
	}})
	if !errors.Is(err, ErrInvalidParameter) {
		t.Fatalf("TransactWriteItems face: err = %v, want ErrInvalidParameter", err)
	}

	// Nothing was written through either face.
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "SepTable", "Key": map[string]interface{}{"id": sVal("a")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	item := resp.(map[string]interface{})["Item"].(map[string]interface{})
	if _, present := item["b"]; present {
		t.Fatal("comma-less SET: b must not exist after the rejection")
	}
	if _, present := item["c"]; present {
		t.Fatal("comma-less SET: c must not exist after the rejection")
	}
}

// TestSameListSetAddsApplyInElementNumberOrder pins the documented ordering
// of multiple list-element adds in one SET clause: "If you add multiple
// elements in a single SET operation, the elements are sorted in order by
// element number." Composed with the append-at-end rule for out-of-range
// indices, the lower element number lands first whatever the expression
// order — the append positions depend on the list's length at application
// time, so expression order permutes the elements.
func TestSameListSetAddsApplyInElementNumberOrder(t *testing.T) {
	newItem := func(elements ...string) map[string]*dbstore.AttributeValue {
		l := make([]*dbstore.AttributeValue, 0, len(elements))
		for _, e := range elements {
			l = append(l, dbstore.StringValue(e))
		}
		return map[string]*dbstore.AttributeValue{"li": dbstore.ListValue(l)}
	}
	values := map[string]*dbstore.AttributeValue{
		":a": dbstore.StringValue("a"),
		":b": dbstore.StringValue("b"),
	}
	readList := func(attrs map[string]*dbstore.AttributeValue) []string {
		out := make([]string, 0, len(attrs["li"].L))
		for _, v := range attrs["li"].L {
			out = append(out, *v.S)
		}
		return out
	}

	// Two out-of-range adds in descending element order: element number 3
	// appends before element number 5.
	attrs := newItem("x", "y")
	_, err := applyUpdateExpressionWithTracking(attrs, "SET li[5] = :a, li[3] = :b", nil, values)
	require.NoError(t, err)
	assert.Equal(t, []string{"x", "y", "b", "a"}, readList(attrs))

	// An append-becomes-append cascade: element number 2 appends first, so
	// element number 3 finds the grown list and appends after it.
	attrs = newItem("x", "y")
	_, err = applyUpdateExpressionWithTracking(attrs, "SET li[3] = :a, li[2] = :b", nil, values)
	require.NoError(t, err)
	assert.Equal(t, []string{"x", "y", "b", "a"}, readList(attrs))

	// The converging case — one in-range write, one append — lands the same
	// under both orders; pinned as the regression guard for the composition.
	attrs = newItem("x", "y")
	_, err = applyUpdateExpressionWithTracking(attrs, "SET li[2] = :a, li[1] = :b", nil, values)
	require.NoError(t, err)
	assert.Equal(t, []string{"x", "b", "a"}, readList(attrs))

	// A single add and in-range-only writes keep their unconditional shape.
	attrs = newItem("x", "y")
	_, err = applyUpdateExpressionWithTracking(attrs, "SET li[2] = :a", nil, values)
	require.NoError(t, err)
	assert.Equal(t, []string{"x", "y", "a"}, readList(attrs))
	attrs = newItem("x", "y")
	_, err = applyUpdateExpressionWithTracking(attrs, "SET li[0] = :b, li[1] = :a", nil, values)
	require.NoError(t, err)
	assert.Equal(t, []string{"b", "a"}, readList(attrs))
}

// TestUpdateActionsEvaluateAgainstPreUpdateItem pins the documented
// evaluation basis of a multi-action update expression: "DynamoDB
// evaluates every action against the item's attribute values as they were
// before the update. The actions aren't applied one after another from
// left to right, so an action's right-hand operand always refers to the
// pre-update value of an attribute—even if another action in the same
// expression also modifies that attribute." The documentation's own
// worked example is REMOVE a SET b = a, c = b over {a:1, b:2, c:3} →
// {b:1, c:2}.
func TestUpdateActionsEvaluateAgainstPreUpdateItem(t *testing.T) {
	// The documentation's worked example: both operands read the original
	// a and b, and a is removed.
	attrs := map[string]*dbstore.AttributeValue{
		"a": dbstore.NumberValue("1"),
		"b": dbstore.NumberValue("2"),
		"c": dbstore.NumberValue("3"),
	}
	_, err := applyUpdateExpressionWithTracking(attrs, "REMOVE a SET b = a, c = b", nil, nil)
	require.NoError(t, err)
	assert.NotContains(t, attrs, "a")
	assert.Equal(t, "1", *attrs["b"].N)
	assert.Equal(t, "2", *attrs["c"].N)

	// A SET operand reads the pre-update value of the attribute an earlier
	// action of the same expression replaces.
	attrs = map[string]*dbstore.AttributeValue{"a": dbstore.NumberValue("1")}
	_, err = applyUpdateExpressionWithTracking(attrs, "SET a = :one, b = a", nil, map[string]*dbstore.AttributeValue{
		":one": dbstore.NumberValue("2"),
	})
	require.NoError(t, err)
	assert.Equal(t, "2", *attrs["a"].N)
	assert.Equal(t, "1", *attrs["b"].N)

	// An ADD's stored base is equally a pre-update read: the later operand
	// reference sees the original number, not the incremented one.
	attrs = map[string]*dbstore.AttributeValue{"n": dbstore.NumberValue("5")}
	_, err = applyUpdateExpressionWithTracking(attrs, "ADD n :inc SET m = n", nil, map[string]*dbstore.AttributeValue{
		":inc": dbstore.NumberValue("2"),
	})
	require.NoError(t, err)
	assert.Equal(t, "7", *attrs["n"].N)
	assert.Equal(t, "5", *attrs["m"].N)

	// The evaluation basis is a deep copy: a nested write inside the same
	// expression must not leak into the later operand's read of the
	// enclosing map.
	attrs = map[string]*dbstore.AttributeValue{
		"m": dbstore.MapValue(map[string]*dbstore.AttributeValue{"k": dbstore.StringValue("old")}),
	}
	_, err = applyUpdateExpressionWithTracking(attrs, "SET m.k = :new, x = m.k", nil, map[string]*dbstore.AttributeValue{
		":new": dbstore.StringValue("new"),
	})
	require.NoError(t, err)
	assert.Equal(t, "new", *attrs["m"].M["k"].S)
	assert.Equal(t, "old", *attrs["x"].S)
}

func TestAddDeleteRejectDocumentPaths(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "TopLevelTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TopLevelTable",
		"Item":      map[string]interface{}{"id": sVal("a"), "n": map[string]interface{}{"N": "1"}},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	update := func(expr string, names, values map[string]interface{}) error {
		params := map[string]interface{}{
			"TableName": "TopLevelTable", "Key": map[string]interface{}{"id": sVal("a")},
			"UpdateExpression": expr,
		}
		if names != nil {
			params["ExpressionAttributeNames"] = names
		}
		if values != nil {
			params["ExpressionAttributeValues"] = values
		}
		_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}
	expectTopLevelRejection := func(op string, err error) {
		t.Helper()
		if err == nil || !strings.Contains(err.Error(), "requires a top-level attribute path") {
			t.Errorf("%s document path: expected the top-level rejection, got %v", op, err)
		}
	}

	expectTopLevelRejection("ADD dotted", update("ADD m.n :v", nil, map[string]interface{}{":v": map[string]interface{}{"N": "1"}}))
	expectTopLevelRejection("ADD bracket", update("ADD li[0] :v", nil, map[string]interface{}{":v": map[string]interface{}{"N": "1"}}))
	expectTopLevelRejection("DELETE dotted", update("DELETE m.s :v", nil, map[string]interface{}{":v": map[string]interface{}{"SS": []interface{}{"x"}}}))

	// The plain top-level forms keep working.
	if err := update("ADD n :one", nil, map[string]interface{}{":one": map[string]interface{}{"N": "1"}}); err != nil {
		t.Fatalf("top-level ADD: unexpected error %v", err)
	}

	// An attribute whose name contains a dot is addressed through an
	// expression attribute name: the alias token carries no separator, so
	// it is a top-level attribute, not a path.
	if err := update("ADD #d :v",
		map[string]interface{}{"#d": "dotted.name"},
		map[string]interface{}{":v": map[string]interface{}{"N": "7"}}); err != nil {
		t.Fatalf("aliased dotted-name ADD: unexpected error %v", err)
	}
	got, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "TopLevelTable", "Key": map[string]interface{}{"id": sVal("a")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	item := got.(map[string]interface{})["Item"].(map[string]interface{})
	dotted, ok := item["dotted.name"].(map[string]interface{})
	if !ok || dotted["N"] != "7" {
		t.Fatalf("aliased dotted-name ADD: expected the literal attribute dotted.name=7, got %v", item)
	}
}

func TestSetOperandsAreDocumentPaths(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "OperandTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "OperandTable",
		"Item": map[string]interface{}{
			"id": sVal("a"),
			"m":  map[string]interface{}{"M": map[string]interface{}{"c": sVal("deep")}},
			"li": map[string]interface{}{"L": []interface{}{sVal("i0"), sVal("i1")}},
		},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	update := func(expr string, values map[string]interface{}) error {
		params := map[string]interface{}{
			"TableName": "OperandTable", "Key": map[string]interface{}{"id": sVal("a")},
			"UpdateExpression": expr,
		}
		if values != nil {
			params["ExpressionAttributeValues"] = values
		}
		_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}
	readBack := func(name string) interface{} {
		resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "OperandTable", "Key": map[string]interface{}{"id": sVal("a")},
		}})
		if err != nil {
			t.Fatalf("read-back get: %v", err)
		}
		return resp.(map[string]interface{})["Item"].(map[string]interface{})[name]
	}

	// A document-path operand copies the nested value.
	if err := update("SET x = m.c", nil); err != nil {
		t.Fatalf("SET x = m.c: unexpected error %v", err)
	}
	if got := readBack("x").(map[string]interface{})["S"]; got != "deep" {
		t.Fatalf("SET x = m.c: expected x=deep, got %v", got)
	}
	if err := update("SET y = li[1]", nil); err != nil {
		t.Fatalf("SET y = li[1]: unexpected error %v", err)
	}
	if got := readBack("y").(map[string]interface{})["S"]; got != "i1" {
		t.Fatalf("SET y = li[1]: expected y=i1, got %v", got)
	}

	// if_not_exists honours an existing nested value at its path operand
	// and takes the default operand when the path addresses nothing.
	if err := update("SET keep = if_not_exists(m.c, :d)", map[string]interface{}{":d": sVal("default")}); err != nil {
		t.Fatalf("if_not_exists(m.c, :d): unexpected error %v", err)
	}
	if got := readBack("keep").(map[string]interface{})["S"]; got != "deep" {
		t.Fatalf("if_not_exists(m.c, :d): expected keep=deep, got %v", got)
	}
	if err := update("SET fresh = if_not_exists(m.zz, :d)", map[string]interface{}{":d": sVal("default")}); err != nil {
		t.Fatalf("if_not_exists(m.zz, :d): unexpected error %v", err)
	}
	if got := readBack("fresh").(map[string]interface{})["S"]; got != "default" {
		t.Fatalf("if_not_exists(m.zz, :d): expected fresh=default, got %v", got)
	}
	if err := update("SET top = if_not_exists(absent, :d)", map[string]interface{}{":d": sVal("default")}); err != nil {
		t.Fatalf("if_not_exists(absent, :d): unexpected error %v", err)
	}
	if got := readBack("top").(map[string]interface{})["S"]; got != "default" {
		t.Fatalf("if_not_exists(absent, :d): expected top=default, got %v", got)
	}
}

func TestExpressionAttributeNamesPerSegment(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "AliasTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	nVal := func(n string) map[string]interface{} { return map[string]interface{}{"N": n} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "AliasTable",
		"Item": map[string]interface{}{
			"id": sVal("a"),
			"m":  map[string]interface{}{"M": map[string]interface{}{"k": sVal("deep"), "li": map[string]interface{}{"L": []interface{}{sVal("i0"), sVal("i1")}}}},
		},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	segmentNames := map[string]interface{}{"#pr": "m", "#k": "k", "#li": "li"}
	update := func(expr string, names map[string]interface{}, values map[string]interface{}, condExpr string) error {
		params := map[string]interface{}{
			"TableName": "AliasTable", "Key": map[string]interface{}{"id": sVal("a")},
			"UpdateExpression": expr,
		}
		if names != nil {
			params["ExpressionAttributeNames"] = names
		}
		if values != nil {
			params["ExpressionAttributeValues"] = values
		}
		if condExpr != "" {
			params["ConditionExpression"] = condExpr
		}
		_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: params})
		return err
	}
	item := func() map[string]interface{} {
		resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "AliasTable", "Key": map[string]interface{}{"id": sVal("a")},
		}})
		if err != nil {
			t.Fatalf("read-back get: %v", err)
		}
		return resp.(map[string]interface{})["Item"].(map[string]interface{})
	}
	nested := func(item map[string]interface{}, segs ...string) interface{} {
		var cur interface{} = item[segs[0]].(map[string]interface{})["M"]
		for _, seg := range segs[1:] {
			cur = cur.(map[string]interface{})[seg]
		}
		return cur
	}

	// SET target: #pr.#k addresses the nested map member.
	if err := update("SET #pr.#k = :v", segmentNames, map[string]interface{}{":v": sVal("new")}, ""); err != nil {
		t.Fatalf("SET #pr.#k: unexpected error %v", err)
	}
	if got := nested(item(), "m", "k").(map[string]interface{})["S"]; got != "new" {
		t.Fatalf("SET #pr.#k: expected m.k=new, got %v", got)
	}

	// Bare SET operand and if_not_exists read their document paths with
	// the same per-segment resolution.
	if err := update("SET x = #pr.#k", segmentNames, nil, ""); err != nil {
		t.Fatalf("SET x = #pr.#k: unexpected error %v", err)
	}
	if got := item()["x"].(map[string]interface{})["S"]; got != "new" {
		t.Fatalf("SET x = #pr.#k: expected x=new, got %v", got)
	}
	if err := update("SET keep = if_not_exists(#pr.#k, :d)", segmentNames, map[string]interface{}{":d": sVal("default")}, ""); err != nil {
		t.Fatalf("if_not_exists(#pr.#k, :d): unexpected error %v", err)
	}
	if got := item()["keep"].(map[string]interface{})["S"]; got != "new" {
		t.Fatalf("if_not_exists(#pr.#k, :d): expected keep=new, got %v", got)
	}
	if err := update("SET fresh = if_not_exists(#pr.#absent, :d)",
		map[string]interface{}{"#pr": "m", "#absent": "absent"},
		map[string]interface{}{":d": sVal("default")}, ""); err != nil {
		t.Fatalf("if_not_exists(#pr.#absent, :d): unexpected error %v", err)
	}
	if got := item()["fresh"].(map[string]interface{})["S"]; got != "default" {
		t.Fatalf("if_not_exists(#pr.#absent, :d): expected fresh=default, got %v", got)
	}

	// A list element inside the nested map, as operand and as target.
	if err := update("SET y = #pr.#li[1]", segmentNames, nil, ""); err != nil {
		t.Fatalf("SET y = #pr.#li[1]: unexpected error %v", err)
	}
	if got := item()["y"].(map[string]interface{})["S"]; got != "i1" {
		t.Fatalf("SET y = #pr.#li[1]: expected y=i1, got %v", got)
	}
	if err := update("SET #pr.#li[0] = :v2", segmentNames, map[string]interface{}{":v2": sVal("y0")}, ""); err != nil {
		t.Fatalf("SET #pr.#li[0]: unexpected error %v", err)
	}
	if got := nested(item(), "m", "li").(map[string]interface{})["L"].([]interface{})[0].(map[string]interface{})["S"]; got != "y0" {
		t.Fatalf("SET #pr.#li[0]: expected m.li[0]=y0, got %v", got)
	}

	// Condition plane: a per-segment comparison resolves the nested value
	// — true when it matches, ConditionalCheckFailed when it does not.
	if err := update("SET c1 = :ok", segmentNames, map[string]interface{}{":ok": sVal("ok"), ":expect": sVal("new")}, "#pr.#k = :expect"); err != nil {
		t.Fatalf("condition #pr.#k = :expect: unexpected error %v", err)
	}
	if err := update("SET c1 = :ok", segmentNames, map[string]interface{}{":ok": sVal("ok"), ":expect": sVal("nope")}, "#pr.#k = :expect"); !errors.Is(err, ErrConditionalCheckFailed) {
		t.Fatalf("condition #pr.#k = :nope: expected ConditionalCheckFailed, got %v", err)
	}
	if err := update("SET c2 = :ok", segmentNames, map[string]interface{}{":ok": sVal("ok")}, "attribute_exists(#pr.#k)"); err != nil {
		t.Fatalf("attribute_exists(#pr.#k): unexpected error %v", err)
	}
	if err := update("SET c3 = :ok", segmentNames, map[string]interface{}{":ok": sVal("ok"), ":n": nVal("2")}, "size(#pr.#li) = :n"); err != nil {
		t.Fatalf("size(#pr.#li) = :n: unexpected error %v", err)
	}
	if err := update("SET c4 = :ok", segmentNames, map[string]interface{}{":ok": sVal("ok"), ":sub": sVal("ew")}, "contains(#pr.#k, :sub)"); err != nil {
		t.Fatalf("contains(#pr.#k, :sub): unexpected error %v", err)
	}

	// REMOVE addresses the nested member, and the existence functions
	// observe the removal through the same per-segment paths.
	if err := update("REMOVE #pr.#k", segmentNames, nil, ""); err != nil {
		t.Fatalf("REMOVE #pr.#k: unexpected error %v", err)
	}
	if _, still := nested(item(), "m").(map[string]interface{})["k"]; still {
		t.Fatal("REMOVE #pr.#k: expected m.k to be gone")
	}
	if err := update("SET c5 = :ok", segmentNames, map[string]interface{}{":ok": sVal("ok")}, "attribute_not_exists(#pr.#k)"); err != nil {
		t.Fatalf("attribute_not_exists(#pr.#k) after REMOVE: unexpected error %v", err)
	}
	if err := update("SET c5 = :ok", segmentNames, map[string]interface{}{":ok": sVal("ok")}, "attribute_exists(#pr.#k)"); !errors.Is(err, ErrConditionalCheckFailed) {
		t.Fatalf("attribute_exists(#pr.#k) after REMOVE: expected ConditionalCheckFailed, got %v", err)
	}

	// A whole-token alias standing for a name that itself contains a dot
	// addresses ONE top-level attribute: resolution happens on parsed
	// segments, never by re-serialising the resolved name into
	// dot-splitting.
	if err := update("SET #d = :v", map[string]interface{}{"#d": "dotted.name"}, map[string]interface{}{":v": sVal("dv")}, ""); err != nil {
		t.Fatalf("SET #d (dotted target name): unexpected error %v", err)
	}
	gotItem := item()
	if got := gotItem["dotted.name"].(map[string]interface{})["S"]; got != "dv" {
		t.Fatalf("SET #d: expected the top-level attribute %q=dv, got %v", "dotted.name", got)
	}
	if _, split := gotItem["dotted"]; split {
		t.Fatal("SET #d: the resolved name must stay one segment, not split into a nested path")
	}

	// An alias the names map does not define remains a validation error on
	// the update plane's strict resolution.
	if err := update("SET #pr.#undef = :v", segmentNames, map[string]interface{}{":v": sVal("x")}, ""); err == nil {
		t.Fatal("SET #pr.#undef: expected rejection of the undefined alias")
	}
}

func TestNestedSetRequiresExistingParents(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "NestedSetTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "NestedSetTable",
		"Item": map[string]interface{}{
			"id": sVal("a"),
			"s":  sVal("scalar"),
			"m":  map[string]interface{}{"M": map[string]interface{}{"c": sVal("leaf")}},
		},
	}}); err != nil {
		t.Fatalf("seed put: %v", err)
	}

	update := func(expr string) error {
		_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "NestedSetTable", "Key": map[string]interface{}{"id": sVal("a")},
			"UpdateExpression": expr,
			"ExpressionAttributeValues": map[string]interface{}{
				":v": sVal("x"), ":v2": sVal("v2"),
			},
		}})
		return err
	}
	expectInvalidPath := func(shape string, err error) {
		t.Helper()
		if err == nil || !strings.Contains(err.Error(), "The document path provided in the update expression is invalid for update") {
			t.Errorf("%s: expected the invalid-document-path rejection, got %v", shape, err)
		}
	}

	// A scalar parent is never silently replaced by a fresh map — the
	// pre-fix behaviour at the top level — and matches the deeper
	// rejection.
	expectInvalidPath("scalar parent", update("SET s.n = :v"))
	// A missing parent at the top level and one level deeper are the same
	// documented rejection, not implicit map creation.
	expectInvalidPath("absent top-level parent", update("SET m2.q = :v"))
	expectInvalidPath("absent deeper parent", update("SET m.d.e = :v"))

	// The scalar attribute is intact after the rejected updates.
	resp, err := svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "NestedSetTable", "Key": map[string]interface{}{"id": sVal("a")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	item := resp.(map[string]interface{})["Item"].(map[string]interface{})
	if got := item["s"].(map[string]interface{})["S"]; got != "scalar" {
		t.Fatalf("scalar parent after rejection: expected s=scalar intact, got %v", got)
	}

	// A fully existing path keeps overwriting its leaf.
	if err := update("SET m.c = :v2"); err != nil {
		t.Fatalf("existing-path SET: unexpected error %v", err)
	}
	resp, err = svc.GetItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName": "NestedSetTable", "Key": map[string]interface{}{"id": sVal("a")},
	}})
	if err != nil {
		t.Fatalf("read-back get: %v", err)
	}
	item = resp.(map[string]interface{})["Item"].(map[string]interface{})
	if got := item["m"].(map[string]interface{})["M"].(map[string]interface{})["c"].(map[string]interface{})["S"]; got != "v2" {
		t.Fatalf("existing-path SET: expected m.c=v2, got %v", got)
	}
}

func TestNumberSetIdentityAcrossBothPlanes(t *testing.T) {
	svc, reqCtx, store := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "NSTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	seed := func() {
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "NSTable",
			"Item":      map[string]interface{}{"id": sVal("a"), "ns": map[string]interface{}{"NS": []interface{}{"1"}}},
		}}); err != nil {
			t.Fatalf("seed put: %v", err)
		}
	}
	storedNS := func() []string {
		item, err := store.Items().Get("NSTable", map[string]*dbstore.AttributeValue{"id": dbstore.StringValue("a")})
		if err != nil {
			t.Fatalf("store read: %v", err)
		}
		attr := item.Attributes["ns"]
		if attr == nil {
			return nil
		}
		return attr.NS
	}

	// UpdateExpression plane.
	seed()
	update := func(expr string) error {
		_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName":                 "NSTable",
			"Key":                       map[string]interface{}{"id": sVal("a")},
			"UpdateExpression":          expr,
			"ExpressionAttributeValues": map[string]interface{}{":v": map[string]interface{}{"NS": []interface{}{"01"}}},
		}})
		return err
	}
	if err := update("ADD ns :v"); err != nil {
		t.Fatalf("UpdateExpression ADD: %v", err)
	}
	if ns := storedNS(); len(ns) != 1 {
		t.Fatalf("UpdateExpression ADD: expected the set to stay single-member, got %v", ns)
	}
	if err := update("DELETE ns :v"); err != nil {
		t.Fatalf("UpdateExpression DELETE: %v", err)
	}
	if ns := storedNS(); ns != nil {
		t.Fatalf("UpdateExpression DELETE: expected the attribute removed, got %v", ns)
	}

	// A fractional member stores in the canonical decimal form: the
	// appended spelling must itself be a valid DynamoDB number, never
	// the fraction form an exact rational's own string would produce.
	seed()
	if _, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "NSTable",
		"Key":                       map[string]interface{}{"id": sVal("a")},
		"UpdateExpression":          "ADD ns :f",
		"ExpressionAttributeValues": map[string]interface{}{":f": map[string]interface{}{"NS": []interface{}{"1.5"}}},
	}}); err != nil {
		t.Fatalf("UpdateExpression fraction ADD: %v", err)
	}
	if ns := storedNS(); len(ns) != 2 || ns[1] != "1.5" || !isValidDynamoDBNumber(ns[1]) {
		t.Fatalf("UpdateExpression fraction ADD: expected the appended member in decimal form, got %v", ns)
	}

	// PartiQL plane.
	seed()
	partiql := func(stmt string) error {
		_, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"Statement":  stmt,
			"Parameters": []interface{}{map[string]interface{}{"NS": []interface{}{"01"}}},
		}})
		return err
	}
	if err := partiql(`UPDATE "NSTable" ADD ns ? WHERE id = 'a'`); err != nil {
		t.Fatalf("PartiQL ADD: %v", err)
	}
	if ns := storedNS(); len(ns) != 1 {
		t.Fatalf("PartiQL ADD: expected the set to stay single-member, got %v", ns)
	}
	if err := partiql(`UPDATE "NSTable" DELETE ns ? WHERE id = 'a'`); err != nil {
		t.Fatalf("PartiQL DELETE: %v", err)
	}
	if ns := storedNS(); ns != nil {
		t.Fatalf("PartiQL DELETE: expected the attribute removed, got %v", ns)
	}

	// The same decimal-form rule on the PartiQL plane.
	seed()
	if _, err := svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement":  `UPDATE "NSTable" ADD ns ? WHERE id = 'a'`,
		"Parameters": []interface{}{map[string]interface{}{"NS": []interface{}{"1.5"}}},
	}}); err != nil {
		t.Fatalf("PartiQL fraction ADD: %v", err)
	}
	if ns := storedNS(); len(ns) != 2 || ns[1] != "1.5" || !isValidDynamoDBNumber(ns[1]) {
		t.Fatalf("PartiQL fraction ADD: expected the appended member in decimal form, got %v", ns)
	}
}

// TestOperandTypeMismatchAnswersOneErrorAcrossFaces pins the unified
// operand-failure contract: the same type-incompatible pair — a string-set
// operand ADDed onto a stored number — answers ErrTypeMismatch verbatim on
// every update face (UpdateItem's expression, a TransactWriteItems Update,
// and a PartiQL UPDATE), the documented ValidationException wording, where
// the faces previously answered the generic Invalid parameter or whichever
// sentinel the calling plane passed into the shared appliers.
func TestOperandTypeMismatchAnswersOneErrorAcrossFaces(t *testing.T) {
	svc, reqCtx, _ := serviceTableFixture(t, dbstore.CreateTableParams{
		Name:                 "MismatchTable",
		KeySchema:            []*dbstore.KeySchemaElement{{AttributeName: "id", KeyType: dbstore.KeyTypeHash}},
		AttributeDefinitions: []*dbstore.AttributeDefinition{{AttributeName: "id", AttributeType: dbstore.ScalarAttributeTypeS}},
		BillingMode:          dbstore.BillingModePayPerRequest,
	})
	sVal := func(s string) map[string]interface{} { return map[string]interface{}{"S": s} }
	seed := func() {
		t.Helper()
		if _, err := svc.PutItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
			"TableName": "MismatchTable",
			"Item":      map[string]interface{}{"id": sVal("k1"), "n": map[string]interface{}{"N": "1"}},
		}}); err != nil {
			t.Fatalf("seed put: %v", err)
		}
	}
	mismatchOperand := map[string]interface{}{"SS": []interface{}{"x"}}

	seed()
	_, err := svc.UpdateItem(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TableName":                 "MismatchTable",
		"Key":                       map[string]interface{}{"id": sVal("k1")},
		"UpdateExpression":          "ADD n :ss",
		"ExpressionAttributeValues": map[string]interface{}{":ss": mismatchOperand},
	}})
	if !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("UpdateItem face: err = %v, want ErrTypeMismatch", err)
	}

	seed()
	_, err = svc.TransactWriteItems(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"TransactItems": []interface{}{map[string]interface{}{"Update": map[string]interface{}{
			"TableName":                 "MismatchTable",
			"Key":                       map[string]interface{}{"id": sVal("k1")},
			"UpdateExpression":          "ADD n :ss",
			"ExpressionAttributeValues": map[string]interface{}{":ss": mismatchOperand},
		}}},
	}})
	if !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("TransactWriteItems face: err = %v, want ErrTypeMismatch", err)
	}

	seed()
	_, err = svc.ExecuteStatement(context.Background(), reqCtx, &request.ParsedRequest{Parameters: map[string]interface{}{
		"Statement":  `UPDATE "MismatchTable" ADD n ? WHERE id = 'k1'`,
		"Parameters": []interface{}{mismatchOperand},
	}})
	if !errors.Is(err, ErrTypeMismatch) {
		t.Fatalf("PartiQL face: err = %v, want ErrTypeMismatch", err)
	}
}
