package interfaces

import (
	"fmt"
	"go/types"
	"path"
	"strings"
)

// Type is a simple representation of a single parameter type.
type Type struct {
	Name        string `json:"name,omitempty"`        // type name
	Package     string `json:"package,omitempty"`     // package name the type is defined in; empty for builtin
	ImportPath  string `json:"importPath,omitempty"`  // import path of the package
	IsPointer   bool   `json:"isPointer,omitempty"`   // whether the parameter is a pointer
	IsComposite bool   `json:"isComposite,omitempty"` // whether the type is map, slice, chan or array
	IsFunc      bool   `json:"isFunc,omitempty"`      // whether the type if function
}

// String gives Go code representation of the type.
func (typ Type) String() (s string) {
	if typ.IsPointer {
		s = "*"
	}
	if !typ.IsComposite && typ.Package != "" {
		s = s + typ.Package + "."
	}
	return s + typ.Name
}

func newType(v *types.Var) (typ Type) {
	typ.setFromType(v.Type(), 0, nil)
	return typ
}

type compositeType interface {
	types.Type
	Elem() types.Type
}

func (typ *Type) setFromType(t types.Type, depth int, orig types.Type) {
	if orig == nil {
		orig = t
	}
	if depth > 128 {
		panic("recursive types not supported: " + orig.String())
	}
	switch t := t.(type) {
	case *types.Basic:
		typ.setFromBasic(t)
	case *types.Interface:
		typ.setFromInterface(t)
	case *types.Struct:
		typ.setFromStruct(t)
	case *types.Named:
		typ.setFromNamed(t)
	case *types.Signature:
		typ.IsFunc = true
		typ.setFromSignature(t)
	case *types.Pointer:
		if depth == 0 {
			typ.IsPointer = true
		}
		typ.setFromType(t.Elem(), depth+1, orig)
	case *types.Map:
		typ.setFromComposite(t, depth, orig)
		typ.setFromType(t.Key(), depth+1, orig)
	case compositeType:
		typ.setFromComposite(t, depth, orig)
	default:
		panic(fmt.Sprintf("internal: t=%T, orig=%T", t, orig))
	}
}

func (typ *Type) setFromBasic(t *types.Basic) {
	if typ.Name == "" {
		typ.Name = t.Name()
	}
}

func (typ *Type) setFromInterface(t *types.Interface) {
	if typ.Name == "" {
		typ.Name = t.String()
	}
}

func (typ *Type) setFromStruct(t *types.Struct) {
	if typ.Name == "" {
		typ.Name = t.String()
	}
}

func (typ *Type) setFromSignature(t *types.Signature) {
	if typ.Name == "" {
		typ.Name = t.String()
	}
}

func (typ *Type) setFromNamed(t *types.Named) {
	if typ.Name == "" {
		typ.Name = t.Obj().Name()
		if typeArgs := t.TypeArgs(); typeArgs != nil && typeArgs.Len() > 0 {
			argValues := make([]string, typeArgs.Len())
			for i := 0; i < typeArgs.Len(); i++ {
				argValues[i] = typeArgs.At(i).String()
			}
			typ.Name = fmt.Sprintf("%s[%s]", typ.Name, strings.Join(argValues, ", "))
		}
	}
	if typ.Package != "" || typ.ImportPath != "" {
		return
	}
	if pkg := t.Obj().Pkg(); pkg != nil {
		typ.Package = pkg.Name()
		typ.ImportPath = pkg.Path()
	}
}

func (typ *Type) setFromComposite(t compositeType, depth int, orig types.Type) {
	typ.IsComposite = true
	if typ.Name == "" {
		typ.Name = t.String()
	}
	typ.setFromType(t.Elem(), depth+1, orig)
}

func fixup(typ *Type, opts *Options) {
	query := opts.Query
	packageName := opts.PackageName

	// Hacky fixup for renaming:
	//
	//   GeoAdd(string, []*github.com/go-redis/redis.GeoLocation) *redis.IntCmd
	//
	// to:
	//
	//   GeoAdd(string, []*redis.GeoLocation) *redis.IntCmd
	//
	// Should be fixed layer below, in type.go.

	// when include other package struct
	if typ.ImportPath != "" && typ.IsComposite {
		if typ.ImportPath == query.Package {
			typ.Name = strings.Replace(typ.Name, typ.ImportPath, typ.Package, -1)
		}

		if typ.ImportPath != query.Package {
			pkgIdx := strings.LastIndex(typ.ImportPath, typ.Package)
			if 0 < pkgIdx {
				typ.Name = strings.Replace(typ.Name, typ.ImportPath[:pkgIdx], "", -1)
			}
		}
	}

	typ.Name = strings.Replace(typ.Name, query.Package, path.Base(query.Package), -1)
	typ.ImportPath = trimVendorPath(typ.ImportPath)

	if typ.Package == packageName {
		typ.Package = ""
		typ.ImportPath = ""
	}
}

// trimVendorPath removes the vendor dir prefix from a package path.
// example: github.com/foo/bar/vendor/github.com/pkg/errors -> github.com/pkg/errors.
func trimVendorPath(p string) string {
	parts := strings.Split(p, "/vendor/")
	if len(parts) == 1 {
		return p
	}
	return strings.TrimLeft(path.Join(parts[1:]...), "/")
}
