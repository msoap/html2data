html2data
=========

Simple wrapper for [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery/)

Install
-------

Install package and command line utility:

    go get -u github.com/msoap/html2data/cmd/html2data

Install package only:

    go get -u github.com/msoap/html2data

Example
-------

    // get title for url
    package main
    
    import (
    	"fmt"
    	"log"
    	"os"
        
    	"github.com/msoap/html2data"
    )
    
    func main() {
    	texts, err := html2data.GetData("url", map[string]string{"one": "title"})
    	if err != nil {
    		log.Fatal(err)
    	}
        
    	if textOne, ok := texts["one"]; ok {
    		for _, text := range textOne {
    			fmt.Println(text)
    		}
    	}
    }

Command line utility
--------------------

    html2data URL "css selector"
    html2data file.html "css selector"
    cat file.html | html2data "css selector"

### install

    brew tap msoap/tools
    brew install html2data
    # update:
    brew update; brew upgrade html2data

### examples

Get title of page:

    html2data https://golang.org/ title

Last blog posts:

    html2data https://blog.golang.org h3

See also
--------

  * [Python package with same name and functionality](https://pypi.python.org/pypi/html2data)
  * [Node.js module](https://www.npmjs.com/package/html2data)
