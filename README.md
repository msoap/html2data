html2data
=========

Extract data from HTML via CSS selectors

Install
-------

Install package and command line utility:

    go get -u github.com/msoap/html2data/cmd/html2data

Install package only:

    go get -u github.com/msoap/html2data

Example
-------

```go
// get title for url
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/msoap/html2data"
)

func main() {
	texts, err := html2data.FromURL("http://example.com").GetData(map[string]string{"title": "title"})
	if err != nil {
		log.Fatal(err)
	}

	if textOne, ok := texts["title"]; ok {
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
