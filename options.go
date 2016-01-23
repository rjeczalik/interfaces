package interfaces

import (
	"errors"
	"fmt"
	"go/build"
	"strings"
)

var errSyntax = errors.New("query string syntax error")

// Query represents a named type request.
type Query struct {
	TypeName string `json:"name,omitempty"`
	Package  string `json:"package,omitempty"`
}

// ParseQuery gives new Query for the given query text.
func ParseQuery(query string) (*Query, error) {
	if query == "" {
		return nil, errors.New("query string is empty")
	}
	var q Query
	if query[0] != '"' {
		return nil, errSyntax
	}
	query = query[1:]
	i := strings.LastIndex(query, `"`)
	if i == -1 || i+1 == len(query) || query[i+1] != '.' {
		return nil, errSyntax
	}
	q.Package = query[:i]
	q.TypeName = query[i+2:]
	if err := q.valid(); err != nil {
		return nil, err
	}
	return &q, nil
}

func (q *Query) valid() error {
	if q == nil {
		return errors.New("query is nil")
	}
	if q.Package == "" {
		return errors.New("package is empty")
	}
	if q.TypeName == "" {
		return errors.New("type name is empty")
	}
	return nil
}

// Options is used for altering behavior of New() function.
type Options struct {
	Query      *Query         // a named type
	Context    *build.Context // build context; see go/build godoc for details
	Unexported bool           // whether to include unexported methods

	CSVHeader  []string
	CSVRecord  []string
	TimeFormat string
}

func (opts *Options) context() *build.Context {
	if opts.Context != nil {
		return opts.Context
	}
	return &build.Default
}

func notFoundErr(opts *Options) error {
	return fmt.Errorf("no exported methods found for %q (package %q)",
		opts.Query.TypeName, opts.Query.Package)
}
