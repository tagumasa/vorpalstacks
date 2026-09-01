package aws

// The handler Core guard enforces the thin-handler architecture rule on the
// AWS API HTTP plane: every registered operation handler is a thin transport
// adapter whose closure may only parse the wire request, build a DTO, call a
// Core function and serialise the result. Handlers — and the private helpers
// they reach — must not call into internal/store/aws packages; validation
// and persistence live in *_core.go Core functions alone.
//
// Detection is purely syntactic (go/parser) and hermetic: it reads only
// git-tracked source under internal/services/aws/, never the Smithy models.
// Handler entry points are methods matching the dispatcher Handler signature
// plus S3's BucketOperations/ObjectOperations methods, so engine files
// (executors, evaluators, protocol servers, workers) are excluded by
// construction. The same-package call closure is name-based and therefore
// over-approximate: a helper pulled in by a name collision can only add
// findings, never hide one.
//
// The sweep is complete: the guard is zero-tolerance. Any store call in a
// handler closure fails the test outright; there is no allowlist and
// reintroducing one is prohibited — new or edited handlers must be
// Core-routed from the start.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

const storeImportPrefix = "vorpalstacks/internal/store/aws/"

// storeTypeNames maps each store package import path to its declared type
// names. A selector expression through a store alias whose name is a type is
// a conversion (tolerated, like every type reference); only function calls
// through the alias are violations. Go forbids a type and a function sharing
// a name in one package, so the lookup is unambiguous.
var storeTypeNames map[string]map[string]bool

// storeInstanceNames are the conventional identifiers and struct fields that
// hold a store instance (or a per-region store bundle) in handler code.
var storeInstanceNames = map[string]bool{"store": true, "stores": true, "st": true}

// s3OpsReceivers lists the S3 operation-layer receiver types; S3 serves its
// own REST plane outside the dispatcher, so its handlers are selected by
// receiver instead of by the Handler signature.
var s3OpsReceivers = map[string]bool{"s3/BucketOperations": true, "s3/ObjectOperations": true}

type funcInfo struct {
	name string
	file string // path relative to internal/services/aws
	decl *ast.FuncDecl

	coreFile    bool // declared in a *_core.go file
	acquisition bool // sanctioned per-region store acquisition helper
}

type pkg struct {
	rel     string
	fset    *token.FileSet
	aliases map[string]map[string]string // file -> identifier -> store import path
	byName  map[string][]*funcInfo
}

type violation struct {
	key    string // "<relfile>:<FuncName>"
	detail string
}

func TestHandlerCoreGuard(t *testing.T) {
	var err error
	storeTypeNames, err = loadStoreTypeNames("../../store/aws")
	if err != nil {
		t.Fatalf("loadStoreTypeNames: %v", err)
	}
	pkgs, err := loadServicePackages(".")
	if err != nil {
		t.Fatalf("loadServicePackages: %v", err)
	}

	var found []violation
	for _, p := range pkgs {
		found = append(found, p.scan()...)
	}

	actual := map[string]violation{}
	for _, v := range found {
		if _, dup := actual[v.key]; dup {
			continue
		}
		actual[v.key] = v
	}

	if len(actual) == 0 {
		handlers := 0
		for _, p := range pkgs {
			for _, fis := range p.byName {
				for _, fi := range fis {
					if p.isHandler(fi) {
						handlers++
					}
				}
			}
		}
		t.Logf("handler Core guard: %d packages, %d handler entries, 0 violations", len(pkgs), handlers)
		return
	}

	keys := make([]string, 0, len(actual))
	for key := range actual {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("GUARD-VIOLATION %s :: %s\n", key, actual[key].detail)
	}
	t.Fatalf("handler Core guard: %d handler violations; see GUARD-VIOLATION lines above", len(actual))
}

// loadStoreTypeNames parses the store tree and indexes declared type names
// per package import path.
func loadStoreTypeNames(root string) (map[string]map[string]bool, error) {
	out := map[string]map[string]bool{}
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && (d.Name() == "testdata" || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		rel, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		importPath := storeImportPrefix + filepath.ToSlash(rel)
		if rel == "." {
			importPath = storeImportPrefix[:len(storeImportPrefix)-1]
		}
		if out[importPath] == nil {
			out[importPath] = map[string]bool{}
		}
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					out[importPath][ts.Name.Name] = true
				}
			}
		}
		return nil
	})
	return out, err
}

// loadServicePackages parses every non-test Go file under root (the
// internal/services/aws directory when the test runs) and groups the
// declarations by package directory.
func loadServicePackages(root string) ([]*pkg, error) {
	byDir := map[string]*pkg{}
	fset := token.NewFileSet()

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && (d.Name() == "testdata" || strings.HasPrefix(d.Name(), ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
		rel, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		p, ok := byDir[rel]
		if !ok {
			p = &pkg{rel: rel, fset: fset, aliases: map[string]map[string]string{}, byName: map[string][]*funcInfo{}}
			byDir[rel] = p
		}
		storeAliases := map[string]string{}
		for _, imp := range file.Imports {
			pathValue := strings.Trim(imp.Path.Value, `"`)
			if !strings.HasPrefix(pathValue, storeImportPrefix) {
				continue
			}
			name := filepath.Base(pathValue)
			if imp.Name != nil {
				name = imp.Name.Name
			}
			storeAliases[name] = pathValue
		}
		// Import aliases are file-scoped; the per-dir map is their union so
		// that an unusual alias in one file is still detected in its own file.
		if p.aliases[rel] == nil {
			p.aliases[rel] = map[string]string{}
		}
		for name, pathValue := range storeAliases {
			p.aliases[rel][name] = pathValue
		}
		base := filepath.Base(path)
		for _, decl := range file.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			fi := &funcInfo{name: fd.Name.Name, file: filepath.Join(rel, base), decl: fd}
			fi.coreFile = strings.HasSuffix(base, "_core.go")
			fi.acquisition = isAcquisitionHelper(fi, base)
			p.byName[fi.name] = append(p.byName[fi.name], fi)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	var out []*pkg
	for _, p := range byDir {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].rel < out[j].rel })
	return out, nil
}

// neverFollowCallees are call names the closure must not descend into: Close
// is overwhelmingly io/lifecycle cleanup (deferring a body or connection
// close), and descending would drag service shutdown helpers into operation
// closures by pure name collision.
var neverFollowCallees = map[string]bool{"Close": true}

// isAcquisitionHelper reports whether the function is the sanctioned
// per-region store acquisition layer (thin helpers on service.go wrapping
// GetOrCreateStoreE, or the conventional *StoreForRegion/store spellings).
// Acquisition helpers may construct stores; they must not be abused for
// anything else, which the wave audits confirm file by file.
func isAcquisitionHelper(fi *funcInfo, base string) bool {
	if fi.name == "store" || strings.Contains(fi.name, "StoreForRegion") || strings.HasPrefix(fi.name, "getStore") {
		return true
	}
	if base != "service.go" || fi.decl.Body == nil {
		return false
	}
	acquires := false
	ast.Inspect(fi.decl.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		if sel, ok := call.Fun.(*ast.SelectorExpr); ok {
			if sel.Sel.Name == "GetOrCreateStoreE" || sel.Sel.Name == "GetOrCreateStore" {
				acquires = true
			}
		}
		return !acquires
	})
	return acquires
}

// scan walks every handler entry point in the package and reports store
// calls reachable from it through same-package, non-Core callees.
func (p *pkg) scan() []violation {
	var out []violation
	for _, fis := range p.byName {
		for _, fi := range fis {
			if !p.isHandler(fi) {
				continue
			}
			for _, v := range p.closureViolations(fi) {
				out = append(out, v)
			}
		}
	}
	return out
}

// isHandler selects HTTP-plane operation handlers: exported methods matching
// the dispatcher Handler signature, or exported methods on the S3 operation
// structs (S3 serves its own REST plane outside the dispatcher).
func (p *pkg) isHandler(fi *funcInfo) bool {
	if fi.decl.Recv == nil || fi.decl.Recv.NumFields() == 0 {
		return false
	}
	if s3OpsReceivers[p.rel+"/"+recvTypeName(fi.decl.Recv)] {
		return true
	}
	if !token.IsExported(fi.name) {
		return false
	}
	ft := fi.decl.Type
	if ft.Params == nil || ft.Params.NumFields() != 3 || ft.Results == nil || ft.Results.NumFields() != 2 {
		return false
	}
	if !isPackageType(ft.Params.List[0].Type, "Context") ||
		!isStarPackageType(ft.Params.List[1].Type, "RequestContext") ||
		!isStarPackageType(ft.Params.List[2].Type, "ParsedRequest") {
		return false
	}
	return isEmptyInterface(ft.Results.List[0].Type) && isError(ft.Results.List[1].Type)
}

func recvTypeName(field *ast.FieldList) string {
	if field == nil || len(field.List) == 0 {
		return ""
	}
	var name string
	switch ty := field.List[0].Type.(type) {
	case *ast.StarExpr:
		if ident, ok := ty.X.(*ast.Ident); ok {
			name = ident.Name
		}
	case *ast.Ident:
		name = ty.Name
	}
	return name
}

func isPackageType(expr ast.Expr, sel string) bool {
	s, ok := expr.(*ast.SelectorExpr)
	return ok && s.Sel.Name == sel
}

func isStarPackageType(expr ast.Expr, sel string) bool {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}
	return isPackageType(star.X, sel)
}

func isEmptyInterface(expr ast.Expr) bool {
	it, ok := expr.(*ast.InterfaceType)
	return ok && (it.Methods == nil || it.Methods.List == nil)
}

func isError(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "error"
}

// closureViolations collects store calls in the transitive same-package
// closure of the handler, stopping at Core functions (*_core.go files or
// names ending in Core) and sanctioned store acquisition helpers.
func (p *pkg) closureViolations(entry *funcInfo) []violation {
	visited := map[*funcInfo]bool{}
	var work []*funcInfo
	var out []violation

	push := func(fi *funcInfo) {
		if fi == nil || visited[fi] || fi.decl.Body == nil {
			return
		}
		if fi.coreFile || strings.HasSuffix(fi.name, "Core") || fi.acquisition {
			return
		}
		visited[fi] = true
		work = append(work, fi)
	}
	push(entry)

	for len(work) > 0 {
		fi := work[len(work)-1]
		work = work[:len(work)-1]
		aliases := p.aliases[p.rel]
		for _, call := range callsIn(fi.decl.Body) {
			for _, callee := range calleeNames(call) {
				if neverFollowCallees[callee] {
					continue
				}
				for _, target := range p.byName[callee] {
					push(target)
				}
			}
			if detail, bad := storeCall(call, aliases); bad {
				line := p.fset.Position(call.Pos()).Line
				out = append(out, violation{
					key:    fmt.Sprintf("%s:%s", fi.file, fi.name),
					detail: fmt.Sprintf("%s at %s:%d", detail, fi.file, line),
				})
			}
		}
	}
	return out
}

// callsIn returns every call expression in the body, including nested
// function literals.
func callsIn(body ast.Node) []*ast.CallExpr {
	var calls []*ast.CallExpr
	ast.Inspect(body, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			calls = append(calls, call)
		}
		return true
	})
	return calls
}

// calleeNames returns the same-package callee candidate names for a call:
// the identifier for plain function calls, the selector for method calls.
func calleeNames(call *ast.CallExpr) []string {
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
		return []string{f.Name}
	case *ast.SelectorExpr:
		return []string{f.Sel.Name}
	}
	return nil
}

// storeCall reports whether the call is a store-package call (through an
// import alias of internal/store/aws/...) or a call on a conventional store
// instance identifier. Type references, constant references and struct
// literals through the alias are tolerated, matching the SFN exemplar.
func storeCall(call *ast.CallExpr, aliases map[string]string) (string, bool) {
	fun := call.Fun
	for {
		if idx, ok := fun.(*ast.IndexExpr); ok {
			fun = idx.X
			continue
		}
		break
	}
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}
	if ident, ok := sel.X.(*ast.Ident); ok && aliases != nil {
		if importPath, isStore := aliases[ident.Name]; isStore {
			if types := storeTypeNames[importPath]; types != nil && types[sel.Sel.Name] {
				return "", false // type conversion, tolerated like a type reference
			}
			return fmt.Sprintf("store package call %s.%s", ident.Name, sel.Sel.Name), true
		}
	}
	if !token.IsExported(sel.Sel.Name) {
		return "", false
	}
	root, fields := unwrapChain(sel.X)
	if root == "" && len(fields) == 0 {
		return "", false
	}
	hit := storeInstanceNames[root]
	for _, f := range fields {
		if storeInstanceNames[f] {
			hit = true
		}
	}
	if hit {
		return fmt.Sprintf("store instance call .%s", sel.Sel.Name), true
	}
	return "", false
}

// unwrapChain flattens a.b.c into the root identifier name and the
// intermediate selector names, stopping at anything that is not a plain
// selector chain (calls, type assertions, index expressions).
func unwrapChain(expr ast.Expr) (string, []string) {
	var fields []string
	for {
		switch e := expr.(type) {
		case *ast.SelectorExpr:
			fields = append([]string{e.Sel.Name}, fields...)
			expr = e.X
		case *ast.Ident:
			return e.Name, fields
		default:
			return "", fields
		}
	}
}
