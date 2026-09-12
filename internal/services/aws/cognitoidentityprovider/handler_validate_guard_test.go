package cognitoidentityprovider

// The validate-reachability guard extends the thin-handler architecture
// rule to the validation half: a handler closure (the handler itself plus
// every same-package private helper it reaches, stopping at *_core.go
// functions and *Core-named functions) must not call a validate* function.
// Validation lives in Core functions alone; a handler may only parse the
// wire request, build a DTO, call a Core function and serialise the result.
//
// Detection mirrors the handler Core guard (internal/services/aws): purely
// syntactic (go/parser), hermetic, reading only this package's non-test
// sources, and name-based over-approximate for callees — a helper pulled in
// by a name collision can only add findings, never hide one. The guard is
// zero-tolerance: any validate* or isValid* call reachable from a handler
// closure fails the test outright - validator detection is name-based over
// exactly those two spellings.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
	"testing"
)

type validateGuardFunc struct {
	name     string
	file     string // base name of the declaring file
	decl     *ast.FuncDecl
	coreFile bool // declared in a *_core.go file
}

// neverFollowCallees are call names the closure must not descend into; Close
// is io/lifecycle cleanup and descending would drag unrelated helpers in by
// name collision.
var neverFollowCallees = map[string]bool{"Close": true}

func TestHandlerValidateGuard(t *testing.T) {
	fset := token.NewFileSet()
	byName := map[string][]*validateGuardFunc{}
	handlers := 0
	var violations []validateGuardViolation

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read package dir: %v", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(fset, e.Name(), nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", e.Name(), err)
		}
		for _, decl := range file.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			fi := &validateGuardFunc{
				name:     fd.Name.Name,
				file:     e.Name(),
				decl:     fd,
				coreFile: strings.HasSuffix(e.Name(), "_core.go"),
			}
			byName[fi.name] = append(byName[fi.name], fi)
			if isValidateGuardHandler(fd) {
				handlers++
				violations = append(violations, validateGuardClosureViolations(fset, fi, byName)...)
			}
		}
	}
	if handlers == 0 {
		t.Fatal("validate guard found no handler entry points; the guard is broken")
	}
	if len(violations) == 0 {
		t.Logf("validate guard: %d handler entries, 0 violations", handlers)
		return
	}
	for _, v := range violations {
		t.Errorf("GUARD-VIOLATION %s :: %s", v.key, v.detail)
	}
}

type validateGuardViolation struct {
	key    string // "<file>:<HandlerName>"
	detail string
}

// isValidateGuardHandler selects HTTP-plane operation handlers: exported
// methods matching the dispatcher Handler signature.
func isValidateGuardHandler(fd *ast.FuncDecl) bool {
	if fd.Recv == nil || fd.Recv.NumFields() == 0 || !fd.Name.IsExported() {
		return false
	}
	ft := fd.Type
	if ft.Params == nil || ft.Params.NumFields() != 3 || ft.Results == nil || ft.Results.NumFields() != 2 {
		return false
	}
	if !validateGuardIsPkgType(ft.Params.List[0].Type, "Context") ||
		!validateGuardIsStarPkgType(ft.Params.List[1].Type, "RequestContext") ||
		!validateGuardIsStarPkgType(ft.Params.List[2].Type, "ParsedRequest") {
		return false
	}
	empty, isEmpty := ft.Results.List[0].Type.(*ast.InterfaceType)
	emptyOK := isEmpty && (empty.Methods == nil || empty.Methods.List == nil)
	errIdent, isErr := ft.Results.List[1].Type.(*ast.Ident)
	return emptyOK && isErr && errIdent.Name == "error"
}

func validateGuardIsPkgType(expr ast.Expr, sel string) bool {
	s, ok := expr.(*ast.SelectorExpr)
	return ok && s.Sel.Name == sel
}

func validateGuardIsStarPkgType(expr ast.Expr, sel string) bool {
	star, ok := expr.(*ast.StarExpr)
	return ok && validateGuardIsPkgType(star.X, sel)
}

// validateGuardClosureViolations collects validate* calls reachable from the
// handler through same-package, non-Core callees.
func validateGuardClosureViolations(fset *token.FileSet, entry *validateGuardFunc, byName map[string][]*validateGuardFunc) []validateGuardViolation {
	visited := map[*validateGuardFunc]bool{}
	var work []*validateGuardFunc
	var out []validateGuardViolation

	push := func(fi *validateGuardFunc) {
		if fi == nil || visited[fi] || fi.decl.Body == nil {
			return
		}
		if fi.coreFile || strings.HasSuffix(fi.name, "Core") {
			return
		}
		visited[fi] = true
		work = append(work, fi)
	}
	push(entry)

	for len(work) > 0 {
		fi := work[len(work)-1]
		work = work[:len(work)-1]
		ast.Inspect(fi.decl.Body, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			callee := validateGuardCalleeName(call)
			if callee == "" {
				return true
			}
			if validateGuardIsValidatorName(callee) {
				line := fset.Position(call.Pos()).Line
				out = append(out, validateGuardViolation{
					key:    fmt.Sprintf("%s:%s", fi.file, fi.name),
					detail: fmt.Sprintf("validate call %s at %s:%d", callee, fi.file, line),
				})
			}
			if neverFollowCallees[callee] {
				return true
			}
			for _, target := range byName[callee] {
				push(target)
			}
			return true
		})
	}

	sort.Slice(out, func(i, j int) bool { return out[i].detail < out[j].detail })
	return out
}

func validateGuardCalleeName(call *ast.CallExpr) string {
	fun := call.Fun
	for {
		if idx, ok := fun.(*ast.IndexExpr); ok {
			fun = idx.X
			continue
		}
		break
	}
	switch f := fun.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.SelectorExpr:
		return f.Sel.Name
	}
	return ""
}

// validateGuardIsValidatorName reports whether a callee name reads as a
// validator: validateXxx or isValidXxx. Detection is name-based; the
// isValid* arm covers validators spelled as predicates, so the same
// spelling cannot appear in a handler closure unnoticed.
func validateGuardIsValidatorName(callee string) bool {
	for _, prefix := range []string{"validate", "isValid"} {
		if len(callee) > len(prefix) && strings.HasPrefix(callee, prefix) {
			c := callee[len(prefix)]
			if c >= 'A' && c <= 'Z' {
				return true
			}
		}
	}
	return false
}
