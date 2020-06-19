package interfaces

import "testing"

func Test_fixup(t *testing.T) {
	cases := map[string]struct {
		typ   *Type
		q     *Query
		fixed string
	}{
		"func (s *Sample) String() string": {
			typ: &Type{
				Name:        "string",
				Package:     "",
				ImportPath:  "",
				IsPointer:   false,
				IsComposite: false,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "string",
		},
		"func (s *Sample) Strings() []string": {
			typ: &Type{
				Name:        "[]string",
				Package:     "",
				ImportPath:  "",
				IsPointer:   false,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "[]string",
		},
		"func (s *Sample) Bar() *Bar": {
			typ: &Type{
				Name:        "Bar",
				Package:     "sample",
				ImportPath:  "github.com/rjeczalik/interfaces/test/sample",
				IsPointer:   true,
				IsComposite: false,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "Bar",
		},
		"func (s *Sample) Bars() []*Bar": {
			typ: &Type{
				Name:        "[]*github.com/rjeczalik/interfaces/test/sample.Bar",
				Package:     "sample",
				ImportPath:  "github.com/rjeczalik/interfaces/test/sample",
				IsPointer:   false,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "[]*sample.Bar",
		},
		"func (s *Sample) WithBar(b *Bar) error": {
			typ: &Type{
				Name:        "Bar",
				Package:     "sample",
				ImportPath:  "github.com/rjeczalik/interfaces/test/sample",
				IsPointer:   true,
				IsComposite: false,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "Bar",
		},
		"func (s *Sample) WithBars(b []*Bar) error": {
			typ: &Type{
				Name:        "[]*github.com/rjeczalik/interfaces/test/sample.Bar",
				Package:     "sample",
				ImportPath:  "github.com/rjeczalik/interfaces/test/sample",
				IsPointer:   false,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "[]*sample.Bar",
		},
		"func (s *Sample) Hoge() *hoge.Hoge": {
			typ: &Type{
				Name:        "Hoge",
				Package:     "hoge",
				ImportPath:  "github.com/rjeczalik/interfaces/test/hoge",
				IsPointer:   true,
				IsComposite: false,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "Hoge",
		},
		"func (s *Sample) Hoges() []*hoge.Hoge": {
			typ: &Type{
				Name:        "[]*github.com/rjeczalik/interfaces/test/hoge.Hoge",
				Package:     "hoge",
				ImportPath:  "github.com/rjeczalik/interfaces/test/hoge",
				IsPointer:   false,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "[]*hoge.Hoge",
		},
		"func (s *Sample) WithHoges(h []*hoge.Hoge) error": {
			typ: &Type{
				Name:        "[]*github.com/rjeczalik/interfaces/test/hoge.Hoge",
				Package:     "hoge",
				ImportPath:  "github.com/rjeczalik/interfaces/test/hoge",
				IsPointer:   false,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "[]*hoge.Hoge",
		},
		"func (s *Sample) Func() func(x int) error": {
			typ: &Type{
				Name:        "Func(x int) error",
				Package:     "",
				ImportPath:  "github.com/rjeczalik/interfaces/test/sample",
				IsPointer:   false,
				IsComposite: false,
				IsFunc:      true,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "Func(x int) error",
		},
		"func (s *Sample) WithFunc(func(x int)) error": {
			typ: &Type{
				Name:        "WithFunc(func(x int)) error",
				Package:     "",
				ImportPath:  "github.com/rjeczalik/interfaces/test/sample",
				IsPointer:   false,
				IsComposite: false,
				IsFunc:      true,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "WithFunc(func(x int)) error",
		},
		"func (s *Sample) WithMap(*[]map[*flag.FlagSet]struct{}, [3]string)": {
			typ: &Type{
				Name:        "WithMap(*[]map[*flag.FlagSet]struct{}, [3]string)",
				Package:     "flag",
				ImportPath:  "flag",
				IsPointer:   true,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "WithMap(*[]map[*flag.FlagSet]struct{}, [3]string)",
		},
		"func (s *Sample) WithMap2(*[]map[string]*hoge.Hoge, [3]string)": {
			typ: &Type{
				Name:        "WithMap2(*[]map[string]*github.com/rjeczalik/interfaces/test/hoge.Hoge)",
				Package:     "hoge",
				ImportPath:  "github.com/rjeczalik/interfaces/test/hoge",
				IsPointer:   true,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Sample",
				Package:  "github.com/rjeczalik/interfaces/test/sample",
			},
			fixed: "WithMap2(*[]map[string]*hoge.Hoge)",
		},
		"func (s *Sample) ListBuckets() ([]minio.BucketInfo)": {
			typ: &Type{
				Name:        "[]github.com/minio/minio-go/v6.BucketInfo",
				Package:     "minio",
				ImportPath:  "github.com/minio/minio-go/v6",
				IsPointer:   false,
				IsComposite: true,
				IsFunc:      false,
			},
			q: &Query{
				TypeName: "Client",
				Package:  "github.com/minio/minio-go/v6",
			},
			fixed: "[]minio.BucketInfo",
		},
	}

	for k, c := range cases {
		t.Run(k, func(t *testing.T) {
			fixup(c.typ, c.q)

			if c.typ.Name != c.fixed {
				t.Errorf("invalid fixedup result. got:%s want:%s", c.typ.Name, c.fixed)
			}
		})
	}
}
