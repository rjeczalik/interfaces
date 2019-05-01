package interfaces_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	gopath, err := ioutil.TempDir("", "interfaces_test")
	if err != nil {
		t.Fatalf("TempDir()=%s", err)
	}
	defer os.RemoveAll(gopath)

	src := filepath.Join(gopath, "src")

	if err := os.MkdirAll(src, 0755); err != nil {
		t.Fatalf("MkdirAll()=%s", err)
	}

	cases := map[string]struct {
		run func(string) error
	}{
		"interfacer": {
			run: func(base string) error {
				args := []string{
					"-for", `os.File`,
					"-as", "interfacer.File",
					"-o", filepath.Join(base, "package.go"),
				}

				p, err := exec.Command("interfacer", args...).CombinedOutput()
				if err != nil {
					return fmt.Errorf("%s:\n%s", err, p)
				}

				return nil
			},
		},
		"structer": {
			run: func(base string) error {
				testdata, err := ioutil.ReadFile(filepath.FromSlash("testdata/aws-billing.csv"))
				if err != nil {
					return err
				}

				args := []string{
					"-tag", "json",
					"-as", "structer.Record",
					"-format", "csv",
					"-o", filepath.Join(base, "package.go"),
				}

				var buf bytes.Buffer

				cmd := exec.Command("structer", args...)
				cmd.Stdin = bytes.NewReader(testdata)
				cmd.Stdout = &buf
				cmd.Stderr = &buf

				if err := cmd.Run(); err != nil {
					return fmt.Errorf("%s:\n%s", err, &buf)
				}

				return nil
			},
		},
	}

	gocommand := func(out io.Writer, pkg string, args ...string) *exec.Cmd {
		c := exec.Command("go", args...)
		c.Stderr = out
		c.Stdout = out
		c.Dir = filepath.Join(gopath, "src", pkg)
		c.Env = []string{
			"PATH=" + os.Getenv("PATH"),
			"GOROOT=" + os.Getenv("GOROOT"),
			"GOPATH=" + gopath,
			"GOCACHE=" + os.Getenv("GOCACHE"),
			"GO111MODULE=on",
		}
		return c
	}

	for pkg, cas := range cases {
		t.Run(pkg, func(t *testing.T) {
			genpkg := filepath.Join(src, pkg)

			if err := os.MkdirAll(genpkg, 0755); err != nil {
				t.Fatalf("MkdirAll()=%s", err)
			}

			if err := cas.run(genpkg); err != nil {
				t.Fatalf("run()=%s", err)
			}

			var buf bytes.Buffer

			if err := gocommand(&buf, pkg, "mod", "init").Run(); err != nil {
				t.Fatalf("gomod.Run()=%s:\n%s", err, &buf)
			}

			buf.Reset()

			if err := gocommand(&buf, pkg, "build", ".").Run(); err != nil {
				t.Fatalf("gobuild.Run()=%s:\n%s", err, &buf)
			}
		})
	}
}
