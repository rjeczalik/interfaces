package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/rjeczalik/interfaces"
)

var (
	time   = flag.String("time", "2006/01/02 15:04:05", "Time format for use with date fields.")
	tag    = flag.String("tag", "", "Name for a struct tag to add to each field.")
	typ    = flag.String("format", "", "Type of the input, overwrites inferred from file name.")
	as     = flag.String("as", "main.Struct", "Generated struct name.")
	input  = flag.String("f", "-", "Input file.")
	output = flag.String("o", "-", "Output file.")
)

// formatter is a helper interface used to build struct definition and
// custom marshallers for a particular format type.
type formatter interface {
	deps() []string
	parse(io.Reader) (*interfaces.Options, error)
	appendTemplate(*vars, io.Writer) error
}

// formats map holds all registered formatters
var formats = make(map[string]formatter)

// deps gives list of import paths that the format depends on.
func deps(typ string) ([]string, error) {
	f, ok := formats[typ]
	if !ok {
		return nil, errors.New("unsupported format type: " + typ)
	}
	return f.deps(), nil
}

// parse reads user-provided file and returns options, which are used
// to create struct definition.
func parse(typ string, r io.Reader) (*interfaces.Options, error) {
	f, ok := formats[typ]
	if !ok {
		return nil, errors.New("unsupported format type: " + typ)
	}
	return f.parse(r)
}

// appendTemplate writes to w custom marshaller/unmarshaller methods for
// struct definition given by the v.
func appendTemplate(typ string, v *vars, w io.Writer) error {
	f, ok := formats[typ]
	if !ok {
		return errors.New("unsupported format type: " + typ)
	}
	return f.appendTemplate(v, w)
}

var tmpl = mustTemplate(`// Created by structer; DO NOT EDIT

package {{.PackageName}}
{{if (eq (.Deps | len) 1)}}{{println}}import "{{(index .Deps 0)}}"{{println}}{{else if (gt (.Deps | len) 1)}}{{println}}import ({{println}}{{range .Deps}}	"{{.}}"{{println}}{{end}}){{println}}{{end}}
// {{.StructName}} is a struct generated from "{{.FileName}}" file.
type {{.StructName}} struct {
{{.Struct}}}
`)

type vars struct {
	Deps        []string
	TimeFormat  string
	PackageName string
	StructName  string
	FileName    string
	Struct      interfaces.Struct
}

func nonil(err ...error) error {
	for _, e := range err {
		if e != nil {
			return e
		}
	}
	return nil
}

func die(v interface{}) {
	fmt.Fprintln(os.Stderr, v)
	os.Exit(1)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		die(err)
	}
}

func run() (err error) {
	var v vars
	var inferredType string

	r := os.Stdin
	if *input != "-" {
		r, err = os.Open(*input)
		if err != nil {
			return err
		}
		defer r.Close()

		v.FileName = filepath.Base(*input)
		if i := strings.LastIndex(v.FileName, "."); i != -1 {
			inferredType = v.FileName[i+1:]
		}
	}

	w := os.Stdout
	if *output != "-" {
		w, err = os.OpenFile(*output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}

	if *typ == "" {
		*typ = inferredType
	}

	opts, err := parse(*typ, r)
	if err != nil {
		return err
	}

	opts.TimeFormat = *time
	v.TimeFormat = *time

	v.Struct, err = interfaces.NewStruct(opts)
	if err != nil {
		return err
	}

	v.Deps, err = deps(*typ)
	if err != nil {
		return err
	}

	v.Deps = append(v.Deps, v.Struct.Deps()...)
	sort.Strings(v.Deps)

	if i := strings.IndexRune(*as, '.'); i != -1 {
		v.PackageName = (*as)[:i]
		v.StructName = (*as)[i+1:]
	} else {
		v.StructName = *as
	}

	if *tag != "" {
		for i := range v.Struct {
			t := interfaces.Tag{
				Name:  *tag,
				Value: camelcase(v.Struct[i].Name),
			}
			v.Struct[i].Tags = append(v.Struct[i].Tags, t)
		}
	}

	return nonil(tmpl.Execute(w, &v), appendTemplate(*typ, &v, w), w.Close())
}

var tmplFuncs = template.FuncMap{
	"receiver": func(typ string) string {
		return string(unicode.ToLower(rune(typ[0])))
	},
	"camelcase": camelcase,
}

func mustTemplate(content string) *template.Template {
	return template.Must(template.New("").Funcs(tmplFuncs).Parse(content))
}

func camelcase(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
