package iot

import "testing"

// "field:> 20" tokenises as "field:>" plus a bare "20"; the lone operator
// must be reported as an error instead of silently matching nothing.
func TestParseQueryRejectsOperatorWithoutValue(t *testing.T) {
	if _, err := parseQuery("attributes.temp:> 20"); err == nil {
		t.Fatal("expected an error for a comparison operator without a value")
	}
	if _, err := parseQuery("attributes.temp:>= 5"); err == nil {
		t.Fatal("expected an error for a comparison operator without a value")
	}

	node, err := parseQuery("attributes.temp:>20")
	if err != nil {
		t.Fatal(err)
	}
	fn, ok := node.(*fieldNode)
	if !ok {
		t.Fatalf("expected a field node, got %T", node)
	}
	if fn.field != "attributes.temp" || fn.op != ">" || fn.value != "20" || !fn.isNum {
		t.Fatalf("unexpected node: %+v", fn)
	}
}
