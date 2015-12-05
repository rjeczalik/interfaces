package interfaces

import (
	"fmt"
	"go/types"
)

// Type
type Type struct {
	Name        string `json:"name,omitempty"`
	Package     string `json:"package,omitempty"`
	ImportPath  string `json:"importPath,omitempty"`
	IsPointer   bool   `json:"isPointer,omitempty"`
	IsComposite bool   `json:"isComposite,omitempty"`
}

// String
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

func (typ *Type) setFromNamed(t *types.Named) {
	if typ.Name == "" {
		typ.Name = t.Obj().Name()
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
