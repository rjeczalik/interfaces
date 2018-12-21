# interfaces [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces?status.png)](https://godoc.org/github.com/rjeczalik/interfaces) [![Build Status](https://img.shields.io/travis/rjeczalik/interfaces/master.svg)](https://travis-ci.org/rjeczalik/interfaces "linux_amd64") [![Build status](https://img.shields.io/appveyor/ci/rjeczalik/interfaces.svg)](https://ci.appveyor.com/project/rjeczalik/interfaces "windows_amd64")
Code generation tools for Go's interfaces.

Tools available in this repository:

- [cmd/interfacer](#cmdinterfacer-)
- [cmd/structer](#cmdstructer-)

### cmd/interfacer [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces/cmd/interfacer?status.png)](https://godoc.org/github.com/rjeczalik/interfaces/cmd/interfacer)

Generates an interface for a named type.

*Installation*
```bash
~ $ go get github.com/rjeczalik/interfaces/cmd/interfacer
```

*Usage*

```bash
~ $ interfacer -help
```
```
Usage of interfacer:
  -all
        Include also unexported methods.
  -as string
        Generated interface name. (default "main.Interface")
  -for string
        Type to generate an interface for.
  -o string
        Output file. (default "-")
```

*Example*
- generate by manually
```bash
~ $ interfacer -for os.File -as mock.File
```
- generate by go generate
```go
//go:generate interfacer -for os.File -as mock.File -o file_iface.go
```
```bash
~ $ go generate  ./...
```
- output
```go
// Created by interfacer; DO NOT EDIT

package mock

import (
        "os"
)

// File is an interface generated for "os".File.
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

### cmd/structer [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces/cmd/structer?status.png)](https://godoc.org/github.com/rjeczalik/interfaces/cmd/structer)

Generates a struct for a formatted file. Currently supported formats are:

- CSV

*Installation*
```bash
~ $ go get github.com/rjeczalik/interfaces/cmd/structer
```

*Usage*

```bash
~ $ structer -help
```
```
Usage of structer:
  -as string
        Generated struct name. (default "main.Struct")
  -f string
        Input file. (default "-")
  -o string
        Output file. (default "-")
  -tag string
        Name for a struct tag to add to each field.
  -type string
        Type of the input, overwrites inferred from file name.
```

*Example*

```bash
~ $ head -2 aws-billing.csv         # first line is a CSV header, second - first line of values
```
```
"InvoiceID","PayerAccountId","LinkedAccountId","RecordType","RecordID","BillingPeriodStartDate","BillingPeriodEndDate","InvoiceDate"
"Estimated","123456","","PayerLineItem","5433212345","2016/01/01 00:00:00","2016/01/31 23:59:59","2016/01/21 19:19:06"
```
```bash
~ $ structer -f aws-billing.csv -tag json -as billing.Record
```
```go
// Created by structer; DO NOT EDIT

package billing

import (
        "strconv"
        "time"
)

// Record is a struct generated from "aws-billing.csv" file.
type Record struct {
        InvoiceID              string    `json:"invoiceID"`
        PayerAccountID         int64     `json:"payerAccountID"`
        LinkedAccountID        string    `json:"linkedAccountID"`
        RecordType             string    `json:"recordType"`
        RecordID               int64     `json:"recordID"`
        BillingPeriodStartDate time.Time `json:"billingPeriodStartDate"`
        BillingPeriodEndDate   time.Time `json:"billingPeriodEndDate"`
        InvoiceDate            time.Time `json:"invoiceDate"`
}

// MarshalCSV encodes r as a single CSV record.
func (r *Record) MarshalCSV() ([]string, error) {
        records := []string{
                r.InvoiceID,
                strconv.FormatInt(r.PayerAccountID, 10),
                r.LinkedAccountID,
                r.RecordType,
                strconv.FormatInt(r.RecordID, 10),
                time.Parse("2006/01/02 15:04:05", r.BillingPeriodStartDate),
                time.Parse("2006/01/02 15:04:05", r.BillingPeriodEndDate),
                time.Parse("2006/01/02 15:04:05", r.InvoiceDate),
        }
        return records, nil
}

// UnmarshalCSV decodes a single CSV record into r.
func (r *Record) UnmarshalCSV(record []string) error {
        if len(record) != 8 {
                return fmt.Errorf("invalud number fields: want 8, got %d", len(record))
        }
        r.InvoiceID = record[0]
        if record[1] != "" {
                if val, err := strconv.ParseInt(record[1], 10, 64); err == nil {
                        r.PayerAccountID = val
                } else {
                        return err
                }
        }
        r.LinkedAccountID = record[2]
        r.RecordType = record[3]
        if record[4] != "" {
                if val, err := strconv.ParseInt(record[4], 10, 64); err == nil {
                        r.RecordID = val
                } else {
                        return err
                }
        }
        if record[5] != "" {
                if val, err := time.Parse("2006/01/02 15:04:05", record[5]); err == nil {
                        r.BillingPeriodStartDate = val
                } else {
                        return err
                }
        }
        if record[6] != "" {
                if val, err := time.Parse("2006/01/02 15:04:05", record[6]); err == nil {
                        r.BillingPeriodEndDate = val
                } else {
                        return err
                }
        }
        if record[7] != "" {
                if val, err := time.Parse("2006/01/02 15:04:05", record[7]); err == nil {
                        r.InvoiceDate = val
                } else {
                        return err
                }
        }
        return nil
}
```
