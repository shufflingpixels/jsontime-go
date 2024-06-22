# jsontime

[![Test](https://github.com/shufflingpixels/jsontime-go/actions/workflows/test.yml/badge.svg)](https://github.com/shufflingpixels/jsontime-go/actions/workflows/test.yml)

A [json iterator](https://github.com/json-iterator/go) extension that support custom time format.

# Install

`go get github.com/shufflingpixels/jsontime`

or

`go mod edit -require=github.com/shufflingpixels/jsontime@latest`


## Usage
100% compatibility with standard lib

Replace
```go
import "encoding/json"

json.Marshal(&data)
json.Unmarshal(input, &data)
```

with
```go
import "github.com/shufflingpixels/jsontime-go"

var json = jsontime.ConfigWithCustomTimeFormat

json.Marshal(&data)
json.Unmarshal(input, &data)
```

## Example

```go
package main

import (
	"fmt"
	"time"

	"github.com/shufflingpixels/jsontime-go"
)

var json = jsontime.Default

func init() {
	jsontime.SetDefaultTimeFormat(time.RFC3339, time.Local)
}

type Book struct {
	Id          int        `json:"id"`
	PublishedAt time.Time  `json:"published_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

func main() {
	book := Book{
		Id:          1,
		PublishedAt: time.Now(),
		UpdatedAt:   nil,
		CreatedAt:   time.Now(),
	}

	bytes, _ := json.Marshal(book)
	fmt.Printf("%s", bytes)
}
```
