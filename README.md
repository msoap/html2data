html2data
=========

[![GoDoc](https://godoc.org/github.com/msoap/html2data?status.svg)](https://godoc.org/github.com/msoap/html2data)
[![Build Status](https://travis-ci.org/msoap/html2data.svg?branch=master)](https://travis-ci.org/msoap/html2data)
[![Coverage Status](https://coveralls.io/repos/github/msoap/html2data/badge.svg?branch=master)](https://coveralls.io/github/msoap/html2data?branch=master)
[![Coverage](https://gocover.io/_badge/github.com/msoap/html2data?1)](https://gocover.io/github.com/msoap/html2data)
[![Report Card](https://goreportcard.com/badge/github.com/msoap/html2data)](https://goreportcard.com/report/github.com/msoap/html2data)

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
  * `doc.GetDataNested(outerCss string, map[string]string)` - extract nested data by CSS-selectors from another CSS-selector
  * `doc.GetDataSingle(string)` - get text by one CSS selector

Pseudo-selectors
----------------

  * `:attr(attr_name)` - getting attribute instead text, for example getting urls from links: `a:attr(href)`
  * `:html` - getting HTML instead text
  * `:get(N)` - get n-th element from list

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

    texts, _ := doc.GetData(map[string]string{"h1": "h1", "links": "a:attr(href)"})
    // get all H1 headers:
    if textOne, ok := texts["h1"]; ok {
        for _, text := range textOne {
            fmt.Println(text)
        }
    }
    // get all urls from links
    if links, ok := texts["links"]; ok {
        for _, text := range links {
            fmt.Println(text)
        }
    }
}
```

TODO
----

  * html2data: get JSON

Command line utility
--------------------

    html2data [options] URL "css selector"
    html2data [options] URL :name1 "css1" :name2 "css2"...
    html2data [options] file.html "css selector"
    cat file.html | html2data "css selector"

### Options

  * `-user-agent="Custom UA"` -- set custom user-agent
  * `-find-in="outer.css.selector"` -- search in the specified elements instead document

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

Getting RSS URL:

    html2data https://blog.golang.org/ 'link[type="application/atom+xml"]:attr(href)'

See also
--------

  * [Python package with same name and functionality](https://pypi.python.org/pypi/html2data)
  * [Node.js module](https://www.npmjs.com/package/html2data)
  * [Go package for CSS selectors](https://github.com/PuerkitoBio/goquery/)
