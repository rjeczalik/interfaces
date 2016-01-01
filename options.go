package interfaces

import (
	"errors"
	"fmt"
	"go/build"
	"strings"
)

var errSyntax = errors.New("query string syntax error")

// Query
type Query struct {
	TypeName string `json:"name,omitempty"`
	Package  string `json:"package,omitempty"`
}

// ParseQuery
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

// Options
type Options struct {
	Query      *Query
	Context    *build.Context
	Unexported bool
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
