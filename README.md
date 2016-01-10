html2data
=========

[![GoDoc](https://godoc.org/github.com/msoap/html2data?status.svg)](https://godoc.org/github.com/msoap/html2data)

Extract data from HTML via CSS selectors

Install
-------

Install package and command line utility:

    go get -u github.com/msoap/html2data/cmd/html2data

Install package only:

    go get -u github.com/msoap/html2data

Methods
-------

  * `FromReader(io.Reader)` - create document for parse
  * `FromURL(URL)` - create document from http(s) URL
  * `FromFile(file)` - create document from local file
  * `doc.GetData(map[string]string)` - get texts by CSS selectors
  * `doc.GetDataSingle(string)` - get text by one CSS selector

Example
-------

```go
package main

import (
    "fmt"
    "log"

    "github.com/msoap/html2data"
)

func main() {
    doc := html2data.FromURL("http://example.com")
    if doc.Err != nil {
        log.Fatal(doc.Err)
    }

    // get title
    title, _ := doc.GetDataSingle("title")
    fmt.Println("Title is:", title)

    // get all H1 headers:
    texts, _ := doc.GetData(map[string]string{"h1": "h1"})
    if textOne, ok := texts["h1"]; ok {
        for _, text := range textOne {
            fmt.Println(text)
        }
    }
}
```

TODO
----

  * get tag attributes
  * html2data: get by several selectors
  * html2data: get JSON

Command line utility
--------------------

    html2data URL "css selector"
    html2data file.html "css selector"
    cat file.html | html2data "css selector"

### TODO: install from homebrew

    brew tap msoap/tools
    brew install html2data
    # update:
    brew update; brew upgrade html2data

### examples

Get title of page:

    html2data https://golang.org/ title

Last blog posts:

    html2data https://blog.golang.org/ h3

See also
--------

  * [Python package with same name and functionality](https://pypi.python.org/pypi/html2data)
  * [Node.js module](https://www.npmjs.com/package/html2data)
  * [Go package for CSS selectors](https://github.com/PuerkitoBio/goquery/)
