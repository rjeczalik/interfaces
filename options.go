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
	idx := strings.LastIndex(query, ".")
	if idx == -1 || query[:idx] == "" || query[idx+1:] == "" {
		return nil, errors.New("generating source should be path/to/package.type")
	}

	return &Query{
		Package:  query[:idx],
		TypeName: query[idx+1:],
	}, nil
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
	Query       *Query         // a named type
	PackageName string         // name of package to generate interface for
	Context     *build.Context // build context; see go/build godoc for details
	Unexported  bool           // whether to include unexported methods

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
