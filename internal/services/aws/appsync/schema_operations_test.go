package appsync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractTypeName(t *testing.T) {
	tests := []struct {
		name string
		def  string
		want string
	}{
		{"type definition", "type Query {", "Query"},
		{"type definition with spaces", "type  MyType  {", "MyType"},
		{"input definition", "input CreateUserInput {", "CreateUserInput"},
		{"enum definition", "enum Status {", "Status"},
		{"interface definition", "interface Node {", "Node"},
		{"union definition", "union SearchResult = User | Post", "SearchResult"},
		{"scalar definition", "scalar DateTime ", "DateTime"},
		{"no matching prefix", "foo bar baz", ""},
		{"type on new line", "\ntype Query {\n", "Query"},
		{"substring should not match", "mytype Query {", "Query"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTypeName(tt.def)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTypeDefInSDL(t *testing.T) {
	sdl := "type Query {\n  getUser: User\n}\ntype User {\n  id: ID\n  name: String\n}\ntype UserProfile {\n  bio: String\n}"

	tests := []struct {
		name string
		def  string
		want bool
	}{
		{"type present", "type User {", true},
		{"type absent", "type Post {", false},
		{"substring should not false-match (Bug A-4)", "type User {", true},
		{"UserProfile distinct from User (Bug A-4 regression)", "type UserProfile {", true},
		{"User does not match UserProfile", "type User {\n  id: ID\n}", true},
		{"non-existent type", "type Order {", false},
		{"enum present in SDL", "enum Status {", false},
		{"input present", "input CreateUserInput {", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := typeDefInSDL(sdl, tt.def)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTypeDefInSDL_SubstringNoFalsePositive(t *testing.T) {
	sdl := "type UserProfile {\n  bio: String\n}\n"
	assert.False(t, typeDefInSDL(sdl, "type User {\n  id: ID\n}"),
		"'type User' should not match 'type UserProfile' in SDL")

	sdl2 := "type User {\n  id: ID\n}\ntype UserProfile {\n  bio: String\n}\n"
	assert.True(t, typeDefInSDL(sdl2, "type User {\n  id: ID\n}"),
		"'type User' should match when both User and UserProfile exist")
}
