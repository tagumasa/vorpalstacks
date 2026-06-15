package policy

type ConditionOperator func(actual, expected string) bool

var conditionOperators = map[string]ConditionOperator{
	"StringEquals":    func(a, b string) bool { return a == b },
	"StringNotEquals": func(a, b string) bool { return a != b },
	"StringLike":      func(a, b string) bool { return wildcardMatch(a, b) },
	"StringNotLike":   func(a, b string) bool { return !wildcardMatch(a, b) },
	"Bool":            func(a, b string) bool { return a == b },
}
