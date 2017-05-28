package interfaces_test

import (
	"bytes"
	"fmt"
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
					"-for", `"os".File`,
					"-as", "mock.File",
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
					"-as", "billing.Record",
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

	for name, cas := range cases {
		t.Run(name, func(t *testing.T) {
			genpkg := filepath.Join(src, name)

			if err := os.MkdirAll(genpkg, 0755); err != nil {
				t.Fatalf("MkdirAll()=%s", err)
			}

			if err := cas.run(genpkg); err != nil {
				t.Fatalf("run()=%s", err)
			}

			var buf bytes.Buffer

			gobuild := exec.Command("go", "build", name)
			gobuild.Stderr = &buf
			gobuild.Stdout = &buf
			gobuild.Env = []string{
				"GOROOT=" + os.Getenv("GOROOT"),
				"GOPATH=" + gopath,
			}

			if err := gobuild.Run(); err != nil {
				t.Fatalf("gobuild.Run()=%s:\n%s", err, &buf)
			}
		})
	}
}
