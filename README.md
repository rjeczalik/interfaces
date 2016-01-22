# interfaces [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces?status.png)](https://godoc.org/github.com/rjeczalik/interfaces) [![Build Status](https://img.shields.io/travis/rjeczalik/interfaces/master.svg)](https://travis-ci.org/rjeczalik/interfaces "linux_amd64") [![Build status](https://img.shields.io/appveyor/ci/rjeczalik/interfaces.svg)](https://ci.appveyor.com/project/rjeczalik/interfaces "windows_amd64")
Code generation tools for Go's interfaces.

### cmd/structer [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces/cmd/structer?status.png)](https://godoc.org/github.com/rjeczalik/interfaces/cmd/structer)

Generates a struct for a formatted file. Currently supported formats are:

- CSV

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
"InvoiceID","PayerAccountId","LinkedAccountId","RecordType","RecordId","ProductName","RateId","SubscriptionId","PricingPlanId","UsageType","Operation","AvailabilityZone","ReservedInstance","ItemDescription","UsageStartDate","UsageEndDate","UsageQuantity","BlendedRate","BlendedCost","UnBlendedRate","UnBlendedCost","ResourceId"
"Estimated","54321","54321","LineItem","543212345","AWS CloudTrail","12345","12345","12345","USE1-FreeEventsRecorded","None","","N","0.0 per free event recorded in US East (N.Virginia) region","2016-01-01 00:00:00","2016-01-01 01:00:00","4105.00000000","0.0000000000","0.00000000","0.0000000000","0.00000000",""
```
```bash
~ $ structer -f aws-billing.csv -tag json -as billing.Record
```
```go
// Created by structer; DO NOT EDIT

package billing

// Record is an struct generated from "aws-billing.csv" file.
type Record struct {
        InvoiceID           string  `json:"invoiceID"`
        PayerAccountID      int64   `json:"payerAccountID"`
        LinkedAccountID     int64   `json:"linkedAccountID"`
        RecordType          string  `json:"recordType"`
        RecordID            float64 `json:"recordID"`
        ProductName         string  `json:"productName"`
        RateID              int64   `json:"rateID"`
        SubscriptionID      int64   `json:"subscriptionID"`
        PricingPlanID       int64   `json:"pricingPlanID"`
        UsageType           string  `json:"usageType"`
        Operation           string  `json:"operation"`
        AvailabilityZone    string  `json:"availabilityZone"`
        ReservedInstance    string  `json:"reservedInstance"`
        ItemDescription     string  `json:"itemDescription"`
        UsageStartDate      string  `json:"usageStartDate"`
        UsageEndDate        string  `json:"usageEndDate"`
        UsageQuantity       float64 `json:"usageQuantity"`
        BlendedRate         float64 `json:"blendedRate"`
        BlendedCost         float64 `json:"blendedCost"`
        UnBlendedRate       float64 `json:"unBlendedRate"`
        UnBlendedCost       float64 `json:"unBlendedCost"`
        ResourceID          string  `json:"resourceID"`
}
```

### cmd/interfacer [![GoDoc](https://godoc.org/github.com/rjeczalik/interfaces/cmd/interfacer?status.png)](https://godoc.org/github.com/rjeczalik/interfaces/cmd/interfacer)

Generates an interface for a named type.

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

```bash
~ $ interfacer -for \"os\".File -as mock.File
```
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
