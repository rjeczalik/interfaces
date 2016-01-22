package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/rjeczalik/interfaces"
)

var (
	tag    = flag.String("tag", "", "Name for a struct tag to add to each field.")
	typ    = flag.String("type", "", "Type of the input, overwrites inferred from file name.")
	as     = flag.String("as", "main.Struct", "Generated struct name.")
	input  = flag.String("f", "-", "Input file.")
	output = flag.String("o", "-", "Output file.")
)

var tmpl = template.Must(template.New("").Parse(`// Created by structer; DO NOT EDIT

package {{.PackageName}}

// {{.StructName}} is a struct generated from "{{.FileName}}" file.
type {{.StructName}} struct {
{{.Struct}}}
`))

type vars struct {
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
		w, err = os.OpenFile(*output, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			return err
		}
	}

	var opts *interfaces.Options
	if *typ != "" {
		opts, err = parse(*typ, r)
	} else {
		opts, err = parse(inferredType, r)
	}
	if err != nil {
		return err
	}

	v.Struct, err = interfaces.NewStruct(opts)
	if err != nil {
		return err
	}

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

	return nonil(tmpl.Execute(w, v), w.Sync(), w.Close())
}

func parse(typ string, r io.Reader) (*interfaces.Options, error) {
	switch typ {
	case "", "txt", "csv":
		dec := csv.NewReader(r)

		header, err := dec.Read()
		if err != nil {
			return nil, err
		}

		record, err := dec.Read()
		if err != nil {
			return nil, err
		}

		opts := &interfaces.Options{
			CSVHeader: header,
			CSVRecord: record,
		}

		return opts, nil
	default:
		return nil, errors.New("unsupported input type: " + typ)
	}
}

func camelcase(s string) string {
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
