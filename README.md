# Sqwiggle Go Client

[![Build Status](https://travis-ci.org/hermanschaaf/sqwiggle.svg)](https://travis-ci.org/hermanschaaf/sqwiggle) [![Coverage Status](https://coveralls.io/repos/hermanschaaf/sqwiggle/badge.svg?branch=master)](https://coveralls.io/r/hermanschaaf/sqwiggle?branch=master) [![Go Report Card](http://goreportcard.com/badge/hermanschaaf/sqwiggle)](http:/goreportcard.com/report/hermanschaaf/sqwiggle)

A Golang client library for the [Sqwiggle API](https://www.sqwiggle.com/docs/overview/getting-started). 

#### Supports

- [x] Attachments
- [x] Conversations
- [x] Info
- [x] Invites
- [x] Messages
- [x] Organizations
- [x] Streams
- [x] Users

#### Installation

```shell
go get github.com/hermanschaaf/sqwiggle
```

#### Example usage

The following example lists the 50 most recent messages across all streams, and prints them out.

```go
package main

import (
	"fmt"
	"github.com/hermanschaaf/sqwiggle"
)

// The following code instantiates a client, then calls the
// ListMessages method to return a slice of Messages. If no error occurred, it
// iterates through the messages and prints them out one by one.
func main() {
	client := sqwiggle.NewClient("YOUR-API-KEY")

	page, limit := 0, 50
	msgs, err := client.ListMessages(page, limit)
	if err != nil {
		panic(err)
	}

	for _, m := range msgs {
		fmt.Printf("%s: %s\n", m.Author.Name, m.Text)
	}
}
```

When instantiating a new client, it is also possible to use your own HTTPClient:

```go
client := Client{
	APIKey:     "YOUR-API-KEY",
	RootURL:    "https://api.sqwiggle.com/", // customize the URL
	HTTPClient: &http.Client{}, // your own custom http.Client
}
``` 

This is useful in testing environments, for example (and is used in the tests for this package).

#### Full Docs

[https://godoc.org/github.com/hermanschaaf/sqwiggle](https://godoc.org/github.com/hermanschaaf/sqwiggle)

There are plenty more usage examples in the docs and also in [sqwiggle_test.go](sqwiggle_test.go).

***** 

MIT License
