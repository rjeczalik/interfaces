package interfaces

import (
	"bytes"
	"errors"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Tag
type Tag struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// String
func (t Tag) String() string {
	return t.Name + `:"` + t.Value + `"`
}

// Tags
type Tags []Tag

// String
func (t Tags) String() string {
	if len(t) == 0 {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteString("`")
	for _, tag := range t {
		buf.WriteString(tag.String())
		buf.WriteByte(' ')
	}
	p := buf.Bytes()
	p[len(p)-1] = '`' // replace danling whitespace
	return string(p)
}

// Field
type Field struct {
	Name string `json:"name,omitempty"`
	Type Type   `json:"type,omitempty"`
	Tags Tags   `json:"tags,omitempty"`
}

// String
func (f *Field) String() string {
	return f.Name + " " + f.Type.String() + " " + f.Tags.String()
}

// Struct
type Struct []Field

// Deps
func (s Struct) Deps() []string {
	pkgs := make(map[string]struct{}, 0)
	for i := range s {
		pkgs[s[i].Type.ImportPath] = struct{}{}
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

// String
func (s Struct) String() string {
	if len(s) == 0 {
		return ""
	}
	types := make([]string, 0, len(s))
	var maxName, maxType int
	var buf bytes.Buffer
	for i := range s {
		typ := s[i].Type.String()
		if n := len(typ); n > maxType {
			maxType = len(typ)
		}
		if n := len(s[i].Name); n > maxName {
			maxName = n
		}
		types = append(types, typ)
	}
	for i := range s {
		buf.WriteByte('\t')
		buf.WriteString(s[i].Name)
		for n := maxName - len(s[i].Name) + 1; n > 0; n-- {
			buf.WriteByte(' ')
		}

		buf.WriteString(types[i])
		for n := maxType - len(types[i]) + 1; n > 0; n-- {
			buf.WriteByte(' ')
		}
		buf.WriteString(s[i].Tags.String())
		buf.WriteByte('\n')
	}
	return buf.String()
}

// NewStruct
func NewStruct(opts *Options) (Struct, error) {
	if opts == nil {
		return nil, errors.New("interfaces: called NewStruct with nil Options")
	}
	if len(opts.CSVHeader) == 0 {
		return nil, errors.New("interfaces: empty CSV header")
	}
	if len(opts.CSVRecord) > len(opts.CSVHeader) {
		return nil, errors.New("interfaces: more CSV records than header fields")
	}
	var s Struct
	n := min(len(opts.CSVHeader), len(opts.CSVRecord))
	for i, h := range opts.CSVHeader[:n] {
		f := Field{
			Name: toFieldName(h),
		}
		v := opts.CSVRecord[i]
		if v == "" {
			f.Type.Name = "string"
		} else if _, err := strconv.ParseBool(v); err == nil {
			f.Type.Name = "bool"
		} else if _, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.Type.Name = "int64"
		} else if _, err := strconv.ParseFloat(v, 64); err == nil {
			f.Type.Name = "float64"
		} else if _, err := time.Parse(opts.TimeFormat, v); err == nil {
			f.Type.Name = "time.Time"
			f.Type.ImportPath = "time"
		} else {
			f.Type.Name = "string"
		}
		s = append(s, f)
	}
	for _, h := range opts.CSVHeader[n:] {
		f := Field{
			Name: toFieldName(h),
			Type: Type{
				Name: "string",
			},
		}
		s = append(s, f)
	}
	return s, nil
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

var (
	nameEscapeChars = "-:_"
	idiomatic       = strings.NewReplacer("Id", "ID")
)

func toFieldName(s string) string {
	r := []rune(idiomatic.Replace(s))
	upper := true
	i := 0
	for _, c := range r {
		if strings.IndexRune(nameEscapeChars, c) != -1 {
			upper = true
			continue
		}
		if upper {
			c = unicode.ToUpper(c)
			upper = false
		}
		r[i] = c
		i++
	}
	return string(r[:i])
}
