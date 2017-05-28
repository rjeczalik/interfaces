package main

import (
	"encoding/csv"
	"io"

	"github.com/rjeczalik/interfaces"
)

func init() {
	formats["csv"] = csvFmt{}
	formats["txt"] = csvFmt{}
}

type csvFmt struct{}

func (csvFmt) deps() []string {
	return []string{"strconv"}
}

func (csvFmt) parse(r io.Reader) (*interfaces.Options, error) {
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
}

var csvTmpl = mustTemplate(`{{with $v := .}}{{with $r := (receiver $v.StructName)}}
// MarshalCSV encodes {{$r}} as a single CSV record.
func ({{$r}} *{{$v.StructName}}) MarshalCSV() ([]string, error) {
	records := []string{ {{range $_, $s := $v.Struct}}{{if (eq $s.Type.Name "string")}}
		{{$r}}.{{$s.Name}},{{else if (eq $s.Type.Name "bool")}}
		strconv.FormatBool({{$r}}.{{$s.Name}}),{{else if (eq $s.Type.Name "int64")}}
		strconv.FormatInt({{$r}}.{{$s.Name}}, 10),{{else if (eq $s.Type.Name "float64")}}
		strconv.FormatFloat({{$r}}.{{$s.Name}}, 'f', -1, 64),{{else if (eq $s.Type.Name "time.Time")}}
		{{$r}}.{{$s.Name}}.Format("{{$v.TimeFormat}}"),{{else}}
		fmt.Sprintf("%v", {{$r}}.{{$s.Name}}),{{end}}{{end}}
	}
	return records, nil
}

// UnmarshalCSV decodes a single CSV record into {{$r}}.
func ({{$r}} *{{$v.StructName}}) UnmarshalCSV(record []string) error {
	if len(record) != {{$v.Struct | len}} {
		return fmt.Errorf("invalud number fields: want {{$v.Struct | len}}, got %d", len(record))
	}
{{range $i, $s := $v.Struct}}{{if (eq $s.Type.Name "string")}}	{{$r}}.{{$s.Name}} = record[{{$i}}]
{{else if (eq $s.Type.Name "bool")}}	if record[{{$i}}] != "" {
		if val, err := strconv.FormatBool(record[{{$i}}); err == nil {
			{{$r}}.{{$s.Name}} = val
		} else {
			return err
		}
	}
{{else if (eq $s.Type.Name "int64")}}	if record[{{$i}}] != "" {
		if val, err := strconv.ParseInt(record[{{$i}}], 10, 64); err == nil {
			{{$r}}.{{$s.Name}} = val
		} else {
			return err
		}
	}
{{else if (eq $s.Type.Name "float64")}}	if record[{{$i}}] != "" {
		if val, err := strconv.ParseFloat(record[{{$i}}], 64); err == nil {
			{{$r}}.{{$s.Name}} = val
		} else {
			return err
		}
	}
{{else if (eq $s.Type.Name "time.Time")}}	if record[{{$i}}] != "" {
		if val, err := time.Parse("{{$v.TimeFormat}}", record[{{$i}}]); err == nil {
			{{$r}}.{{$s.Name}} = val
		} else {
			return err
		}
	}
{{end}}{{end}}	return nil
}
{{end}}{{end}}`)

func (csvFmt) appendTemplate(v *vars, w io.Writer) error {
	return csvTmpl.Execute(w, v)
}
