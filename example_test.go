package interfaces_test

import (
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/rjeczalik/interfaces"
)

type Generic1[A any] struct {
	Value A
}

type Generic2[A any, B any] struct {
	Value  A
	Value2 B
}

type ExampleFoo int

type ExampleBar struct{}

type ExampleBaz struct {
	*ExampleBar
}

func (ExampleBar) A(int) int {
	return 0
}

func (*ExampleBar) B(*string, io.Writer, ExampleFoo) (*ExampleFoo, int) {
	return nil, 0
}

func (ExampleBar) C(map[string]int, *interfaces.Options, *http.Client) (chan []string, error) {
	return nil, nil
}

func (ExampleBaz) D(*map[interface{}]struct{}, interface{}) (chan struct{}, []interface{}) {
	return nil, nil
}

func (*ExampleBaz) E(*[]map[*flag.FlagSet]struct{}, [3]string) {}

func (*ExampleBaz) F(v Generic1[io.Writer]) (Generic2[io.Writer, ExampleFoo], int) {
	return Generic2[io.Writer, ExampleFoo]{}, 0
}

func (*ExampleBaz) G(v Generic2[string, io.Writer]) (Generic1[int], int) {
	return Generic1[int]{}, 0
}

func ExampleNew() {
	i, err := interfaces.New(`github.com/rjeczalik/interfaces.ExampleBaz`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Interface:")
	for _, fn := range i {
		fmt.Println(fn)
	}
	fmt.Println("Dependencies:")
	for _, dep := range i.Deps() {
		fmt.Println(dep)
	}
	// Output: Interface:
	// A(int) int
	// B(*string, io.Writer, interfaces_test.ExampleFoo) (*interfaces_test.ExampleFoo, int)
	// C(map[string]int, *interfaces.Options, *http.Client) (chan []string, error)
	// D(*map[interface{}]struct{}, interface{}) (chan struct{}, []interface{})
	// E(*[]map[*flag.FlagSet]struct{}, [3]string)
	// F(interfaces_test.Generic1[io.Writer]) (interfaces_test.Generic2[io.Writer, interfaces_test.ExampleFoo], int)
	// G(interfaces_test.Generic2[string, io.Writer]) (interfaces_test.Generic1[int], int)
	// Dependencies:
	// flag
	// github.com/rjeczalik/interfaces
	// github.com/rjeczalik/interfaces_test
	// io
	// net/http
}

func ExampleNewWithOptions() {
	opts := &interfaces.Options{
		Query: &interfaces.Query{
			Package:  "net",
			TypeName: "Interface",
		},
	}
	i, err := interfaces.NewWithOptions(opts)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Interface:")
	for _, fn := range i {
		fmt.Println(fn)
	}
	fmt.Println("Dependencies:")
	for _, dep := range i.Deps() {
		fmt.Println(dep)
	}
	// Output: Interface:
	// Addrs() ([]net.Addr, error)
	// MulticastAddrs() ([]net.Addr, error)
	// Dependencies:
	// net
}

func ExampleFunc_String() {
	f := interfaces.Func{
		Name: "Close",
		Outs: []interfaces.Type{{Name: "error"}},
	}
	fmt.Println(f)
	// Output: Close() error
}
