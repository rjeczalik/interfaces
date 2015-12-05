package interfaces

import (
	"bytes"
	"fmt"
	"sort"
)

// Func
type Func struct {
	Name string `json:"name,omitempty"`
	Ins  []Type `json:"ins,omitempty"`
	Outs []Type `json:"outs,omitempty"`
}

// String
func (f Func) String() string {
	var buf bytes.Buffer
	if len(f.Ins) == 0 {
		fmt.Fprintf(&buf, "%s()", f.Name)
	} else {
		fmt.Fprintf(&buf, "%s(%s", f.Name, f.Ins[0])
		for _, typ := range f.Ins[1:] {
			fmt.Fprintf(&buf, ", %s", typ)
		}
		buf.WriteString(")")
	}
	if len(f.Outs) == 1 {
		fmt.Fprintf(&buf, " %s", f.Outs[0])
	} else if len(f.Outs) > 1 {
		fmt.Fprintf(&buf, " (%s", f.Outs[0])
		for _, typ := range f.Outs[1:] {
			fmt.Fprintf(&buf, ", %s", typ)
		}
		buf.WriteString(")")
	}
	return buf.String()
}

// Deps
func (f Func) Deps() []string {
	pkgs := make(map[string]struct{}, 0)
	for _, in := range f.Ins {
		pkgs[in.ImportPath] = struct{}{}
	}
	for _, out := range f.Outs {
		pkgs[out.ImportPath] = struct{}{}
	}
	delete(pkgs, "")
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

type byName []Func

func (p byName) Len() int           { return len(p) }
func (p byName) Less(i, j int) bool { return p[i].Name < p[j].Name }
func (p byName) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
