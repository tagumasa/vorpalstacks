package dynamodb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMatchingCloseParen(t *testing.T) {
	t.Run("simple match", func(t *testing.T) {
		assert.Equal(t, 4, findMatchingCloseParen("(abc)", 0))
	})

	t.Run("nested parentheses (Bug D-8)", func(t *testing.T) {
		assert.Equal(t, 10, findMatchingCloseParen("(a,(b,c),d)", 0))
	})

	t.Run("deeply nested", func(t *testing.T) {
		assert.Equal(t, 8, findMatchingCloseParen("(a(b(c)))", 0))
	})

	t.Run("list_append nested (Bug D-8 regression)", func(t *testing.T) {
		s := "list_append(:a, list_append(:b, :c))"
		openPos := strings.Index(s, "(")
		closePos := findMatchingCloseParen(s, openPos)
		assert.Equal(t, len(s)-1, closePos,
			"should match outermost close paren")
		assert.Equal(t, s[openPos:closePos+1], "(:a, list_append(:b, :c))")
	})

	t.Run("no matching close returns -1", func(t *testing.T) {
		assert.Equal(t, -1, findMatchingCloseParen("(abc", 0))
	})

	t.Run("empty parens", func(t *testing.T) {
		assert.Equal(t, 1, findMatchingCloseParen("()", 0))
	})

	t.Run("multiple groups", func(t *testing.T) {
		assert.Equal(t, 2, findMatchingCloseParen("(a)(b)", 0))
		assert.Equal(t, 5, findMatchingCloseParen("(a)(b)", 3))
	})
}
