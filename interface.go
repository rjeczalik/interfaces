package interfaces

import (
	"errors"
	"fmt"
	"sort"
	"unicode"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
)

// Interface
type Interface []Func

// New
func New(query string) (Interface, error) {
	q, err := ParseQuery(query)
	if err != nil {
		return nil, errors.New("invalid query: " + err.Error())
	}
	opts := &Options{
		Query: q,
	}
	return NewWithOptions(opts)
}

// NewWithOptions
func NewWithOptions(opts *Options) (Interface, error) {
	if opts == nil || opts.Query == nil {
		panic("interfacer: called NewWithOptions with nil Options or nil Query")
	}
	if err := opts.Query.valid(); err != nil {
		return nil, errors.New("invalid query: " + err.Error())
	}
	return buildInterface(opts)
}

// Deps
func (i Interface) Deps() []string {
	pkgs := make(map[string]struct{})
	for _, fn := range i {
		for _, pkg := range fn.Deps() {
			pkgs[pkg] = struct{}{}
		}
	}
	if len(pkgs) == 0 {
		return nil
	}
	deps := make([]string, 0, len(pkgs))
	for pkg := range pkgs {
		deps = append(deps, pkg)
	}
	sort.Strings(deps)
	return deps
}

func buildInterface(opts *Options) (Interface, error) {
	cfg := &loader.Config{
		AllowErrors:         true,
		Build:               opts.context(),
		ImportPkgs:          map[string]bool{opts.Query.Package: true},
		TypeCheckFuncBodies: func(string) bool { return false },
	}
	cfg.ImportWithTests(opts.Query.Package)
	prog, err := cfg.Load()
	if err != nil {
		return nil, err
	}
	pkg, ok := prog.Imported[opts.Query.Package]
	if !ok {
		return nil, fmt.Errorf("parsing successful, but package %q not found",
			opts.Query.Package)
	}
	i, err := buildInterfaceForPkg(pkg, opts)
	if err == nil {
		return i, nil
	}
	queryCopy := *opts.Query
	queryCopy.Package += "_test"
	optsCopy := *opts
	optsCopy.Query = &queryCopy
	for _, pkg := range prog.Created {
		if pkg.Pkg.Path() == optsCopy.Query.Package {
			return buildInterfaceForPkg(pkg, &optsCopy)
		}
	}
	return nil, err
}

func buildInterfaceForPkg(pkg *loader.PackageInfo, opts *Options) (Interface, error) {
	var typ *types.Named
	for _, obj := range pkg.Defs {
		if obj == nil {
			continue
		}
		if obj.Name() != opts.Query.TypeName || obj.Pkg().Path() != opts.Query.Package {
			continue
		}
		var ok bool
		typ, ok = obj.Type().(*types.Named)
		if ok {
			break
		}
	}
	if typ == nil {
		return nil, notFoundErr(opts)
	}
	var inter Interface
	var methods = make(map[string]*types.Func)
	collectMethods(methods, typ, 0, nil)
	for _, method := range methods {
		// TODO(rjeczalik): read rune
		if unicode.IsLower(rune(method.Name()[0])) && !opts.Unexported {
			continue
		}
		sig, ok := method.Type().(*types.Signature)
		if !ok {
			continue
		}
		ins := sig.Params()
		outs := sig.Results()
		fn := Func{
			Name: method.Name(),
			Ins:  make([]Type, ins.Len()),
			Outs: make([]Type, outs.Len()),
		}
		for i := range fn.Ins {
			fn.Ins[i] = newType(ins.At(i))
		}
		for i := range fn.Outs {
			fn.Outs[i] = newType(outs.At(i))
		}
		inter = append(inter, fn)
	}
	if len(inter) == 0 {
		return nil, notFoundErr(opts)
	}
	sort.Sort(byName(inter))
	return inter, nil
}

func collectMethods(methods map[string]*types.Func, typ *types.Named, depth int, orig types.Type) {
	if orig == nil {
		orig = typ
	}
	if depth > 128 {
		panic("recursive types not supported: " + orig.String())
	}
	for i := 0; i < typ.NumMethods(); i++ {
		method := typ.Method(i)
		if _, ok := methods[method.Name()]; ok {
			continue
		}
		methods[method.Name()] = method
	}
	if typ, ok := typ.Underlying().(*types.Struct); ok {
		for i := 0; i < typ.NumFields(); i++ {
			typ := typ.Field(i).Type()
			if p, ok := typ.(*types.Pointer); ok {
				typ = p.Elem()
			}
			if named, ok := typ.(*types.Named); ok {
				collectMethods(methods, named, depth+1, orig)
			}
		}
	}
}
