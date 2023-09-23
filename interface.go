package interfaces

import (
	"errors"
	"go/types"
	"sort"
	"unicode"

	"golang.org/x/tools/go/packages"
)

// Interface represents a typed interface.
type Interface []Func

// New builds an interface definition for a type specified by the query.
// Supported query format is "package".Type (similar to what gorename
// tool accepts).
//
// The function expects sources for the requested type to be present
// in current GOPATH.
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

// NewWithOptions builds an interface definition for a type specified by
// the given Options.
//
// The Options may be used to specify e.g. different GOPATH if sources
// for requested type are not available in the current one.
func NewWithOptions(opts *Options) (Interface, error) {
	if opts == nil || opts.Query == nil {
		panic("interfacer: called NewWithOptions with nil Options or nil Query")
	}
	if err := opts.Query.valid(); err != nil {
		return nil, errors.New("invalid query: " + err.Error())
	}
	return buildInterface(opts)
}

// Deps gives a list of packages the interface depends on.
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
		if pkg != "" {
			deps = append(deps, pkg)
		}
	}
	sort.Strings(deps)
	return deps
}

func buildInterface(opts *Options) (Interface, error) {
	var err error

	// If a requested type is defined in an external test package try to
	// build the interface using it before returning an error.  Therefore,
	// set the `Tests` flag to true so that any test packages are also
	// loaded and checked below.
	cfg := &packages.Config{
		Mode:  packages.NeedTypes | packages.NeedTypesInfo | packages.NeedImports | packages.NeedDeps,
		Tests: true,
	}

	pkgs, err := packages.Load(cfg, opts.Query.Package)
	if err != nil {
		return nil, err
	}

	for _, pkg := range pkgs {
		opts.Query.Package = pkg.Types.Path()
		i, err := buildInterfaceForPkg(pkg, opts)
		if err == nil {
			return i, nil
		}
	}
	return nil, err
}

func buildInterfaceForPkg(pkg *packages.Package, opts *Options) (Interface, error) {
	var typ *types.Named
	for _, obj := range pkg.TypesInfo.Defs {
		if obj == nil {
			continue
		}
		if obj.Name() != opts.Query.TypeName || obj.Pkg().Path() != opts.Query.Package {
			continue
		}
		tmp, ok := obj.Type().(*types.Named)
		if !ok {
			continue
		}
		if tmp.Obj() == obj {
			typ = tmp
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
			Name:       method.Name(),
			Ins:        make([]Type, ins.Len()),
			Outs:       make([]Type, outs.Len()),
			IsVariadic: sig.Variadic(),
		}
		for i := range fn.Ins {
			fn.Ins[i] = newType(ins.At(i))
			fixup(&fn.Ins[i], opts)
		}
		for i := range fn.Outs {
			fn.Outs[i] = newType(outs.At(i))
			fixup(&fn.Outs[i], opts)
		}
		inter = append(inter, fn)
	}
	if len(inter) == 0 {
		return nil, notFoundErr(opts)
	}
	sort.Sort(funcs(inter))
	return inter, nil
}

func collectMethods(methods map[string]*types.Func, typ *types.Named, depth int, orig types.Type) {
	if orig == nil {
		orig = typ
	}
	// TODO(rjeczalik): recursive types support
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
			field := typ.Field(i)
			if !field.Anonymous() {
				continue
			}
			typ := field.Type()
			if p, ok := typ.(*types.Pointer); ok {
				typ = p.Elem()
			}
			if named, ok := typ.(*types.Named); ok {
				collectMethods(methods, named, depth+1, orig)
			}
		}
	}
}
