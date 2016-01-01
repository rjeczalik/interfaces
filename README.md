# interfaces [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces?status.png)](https://godoc.org/github.com/rjeczalik/interfaces) [![Build Status](https://img.shields.io/travis/rjeczalik/interfaces/master.svg)](https://travis-ci.org/rjeczalik/cmd "linux_amd64") [![Build status](https://img.shields.io/appveyor/ci/rjeczalik/interfaces.svg)](https://ci.appveyor.com/project/rjeczalik/interfaces "windows_amd64") 
Code generation tools for Go's interfaces.

*Installation*

```
~ $ go get -u github.com/rjeczalik/interfaces/...
```

### cmd/interfacer [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces/cmd/interfacer?status.png)](https://godoc.org/github.com/rjeczalik/interfaces/cmd/interfacer)

Generates an interface for a named type.

*Usage*

```
~ $ interfacer -help
```
```
Usage of interfacer:
  -all
        Include also unexproted methods.
  -as string
        Generated interface name. (default "main.Interface")
  -for string
        Type to generate an interface for.
  -o string
        Output file. (default "-")
```

*Example*

```
~ $ $ interfacer -for \"os\".File -as mock.File
```
```go
// Created by interfacer; DO NOT EDIT

package mock

import (
        "os"
)

type File interface {
        Chdir() error
        Chmod(os.FileMode) error
        Chown(int, int) error
        Close() error
        Fd() uintptr
        Name() string
        Read([]byte) (int, error)
        ReadAt([]byte, int64) (int, error)
        Readdir(int) ([]os.FileInfo, error)
        Readdirnames(int) ([]string, error)
        Seek(int64, int) (int64, error)
        Stat() (os.FileInfo, error)
        Sync() error
        Truncate(int64) error
        Write([]byte) (int, error)
        WriteAt([]byte, int64) (int, error)
        WriteString(string) (int, error)
}
```
