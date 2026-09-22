package cloudwatchlogs

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// The documented Logs Insights QL function library's entry and its
// general family: the families the operations-functions page enumerates
// (general, numeric, string, IP address, JSON, datetime) each live in
// their own file, and the entry evaluates one call's arguments and
// dispatches by name.

// structureNullFunction reports whether a function belongs to the string,
// number, or datetime families — "Maps and lists are treated as null for
// string, number, and datetime functions." General, IP and JSON functions
// take structures by their own contracts (ispresent, jsonstringify,
// jsonarraysize, ...).
func structureNullFunction(name string) bool {
	switch name {
	// Number functions.
	case "abs", "ceil", "floor", "greatest", "least", "log", "round",
		"sqrt", "haversine", "tonumber", "toint", "tolong", "todouble":
		return true
	// String functions.
	case "isempty", "isblank", "concat", "ltrim", "rtrim", "trim",
		"strlen", "toupper", "tolower", "substr", "replace",
		"regexreplace", "strcontains", "startswith", "endswith",
		"urlencode", "urldecode", "base64encode", "base64decode", "split",
		"hextoascii", "hextodec", "dectohex":
		return true
	// Datetime functions.
	case "datefloor", "dateceil", "frommillis", "tomillis", "parsedate",
		"formatdate":
		return true
	}
	return false
}

// funcArgs carries one call's evaluated arguments and the execution
// context; the accessor set is the interface every family dispatcher
// reads its arguments through.
type funcArgs struct {
	vals []interface{}
	ctx  *execContext
}

func (f funcArgs) arg(i int) interface{} {
	if i < len(f.vals) {
		return f.vals[i]
	}
	return nil
}

func (f funcArgs) num(i int) (float64, bool) { return asNumber(f.arg(i)) }

func (f funcArgs) str(i int) string { return asString(f.arg(i)) }

// callQueryFunction evaluates one function call: the arguments evaluate
// against the row, the structure-null coercion the string, number and
// datetime families document runs here once, and the family dispatchers
// decide the name in turn.
func callQueryFunction(name string, args []exprNode, row *queryResultRow, ctx *execContext) interface{} {
	structureAsNull := structureNullFunction(name)
	vals := make([]interface{}, len(args))
	for i, a := range args {
		v := a.eval(row, ctx)
		if structureAsNull && isStructure(v) {
			// "Maps and lists are treated as null for string, number,
			// and datetime functions" — the documented example shows
			// toupper over a parsed map yielding empty results.
			v = nil
		}
		vals[i] = v
	}
	f := funcArgs{vals: vals, ctx: ctx}

	if v, ok := evalGeneralFunctions(name, f); ok {
		return v
	}
	if v, ok := evalNumericFunctions(name, f); ok {
		return v
	}
	if v, ok := evalStringFunctions(name, f); ok {
		return v
	}
	if v, ok := evalIPFunctions(name, f); ok {
		return v
	}
	if v, ok := evalJSONFunctions(name, f); ok {
		return v
	}
	if v, ok := evalDatetimeFunctions(name, f); ok {
		return v
	}
	return nil
}

// evalGeneralFunctions holds the general family's dispatch — the query
// metadata trio reads the execution context — plus the two documented
// hashing functions: false means the name is not this family's.
func evalGeneralFunctions(name string, f funcArgs) (interface{}, bool) {
	switch name {
	case "ispresent":
		return f.arg(0) != nil, true
	case "coalesce":
		// "Returns the first non-null value from the list" — an empty
		// string is a non-null value and wins over later arguments.
		for _, v := range f.vals {
			if v != nil {
				return v, true
			}
		}
		return nil, true
	case "case":
		// case(cond1, val1, ..., [default]): pairs then optional default.
		i := 0
		for i+1 < len(f.vals) {
			if truthy(f.vals[i]) {
				return f.vals[i+1], true
			}
			i += 2
		}
		if len(f.vals)%2 == 1 {
			return f.vals[len(f.vals)-1], true
		}
		return nil, true
	case "if":
		// The condition is validated as a three-argument call at parse time;
		// arg bounds defensively so a nil context cannot panic here.
		if truthy(f.arg(0)) {
			return f.arg(1), true
		}
		return f.arg(2), true
	case "isnumeric":
		_, ok := f.num(0)
		return ok, true
	case "messagesize":
		return float64(len(f.str(0))), true
	case "querystarttime":
		if f.ctx != nil {
			return float64(f.ctx.startTime), true
		}
		return nil, true
	case "queryendtime":
		if f.ctx != nil {
			return float64(f.ctx.endTime), true
		}
		return nil, true
	case "querytimerange":
		if f.ctx != nil {
			return float64(f.ctx.endTime - f.ctx.startTime), true
		}
		return nil, true

	// Hashing functions.
	case "md5":
		sum := md5.Sum([]byte(f.str(0)))
		return hex.EncodeToString(sum[:]), true
	case "sha256":
		sum := sha256.Sum256([]byte(f.str(0)))
		return hex.EncodeToString(sum[:]), true
	}
	return nil, false
}
