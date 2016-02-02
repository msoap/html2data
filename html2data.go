/*
Package html2data - extract data from HTML via CSS selectors

Install package and command line utility:

	go get -u github.com/msoap/html2data/cmd/html2data

Install package only:

	go get -u github.com/msoap/html2data

Allowed pseudo-selectors:

 * `:attr(attr_name)` - for getting attributes instead text
 * `:html` - for getting HTML instead text

Command line utility:

    html2data URL "css selector"
    html2data file.html "css selector"
    cat file.html | html2data "css selector"

*/
package html2data

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Doc - html document for parse
type Doc struct {
	doc *goquery.Document
	Err error
}

// GetData - extract data by CSS-selectors
//  texts, err := doc.GetData(map[string]string{"h1": "h1"})
func (doc Doc) GetData(selectors map[string]string) (result map[string][]string, err error) {
	if doc.Err != nil {
		return result, fmt.Errorf("parse document error: %s", doc.Err)
	}

	result = map[string][]string{}
	for name, selector := range selectors {
		selector, attrName, getHTML, err := parseSelector(selector)
		if err != nil {
			return result, err
		}

		texts := []string{}
		doc.doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
			if attrName != "" {
				texts = append(texts, selection.AttrOr(attrName, ""))
			} else if getHTML {
				HTML, err := selection.Html()
				if err != nil {
					return
				}
				texts = append(texts, HTML)
			} else {
				texts = append(texts, selection.Text())
			}
		})
		result[name] = texts
	}

	return result, err
}

// parseSelector - parse pseudo-selectors:
// :attr(href) - for getting attribute instead text node
func parseSelector(inputSelector string) (outSelector string, attrName string, getHTML bool, err error) {
	htmlAttrRe := regexp.MustCompile(`^\s*(\w+)\s*(?:\(\s*(\w+)\s*\))?\s*$`)

	parts := strings.Split(inputSelector, ":")
	outSelector, parts = parts[0], parts[1:]
	for _, part := range parts {
		reParts := htmlAttrRe.FindStringSubmatch(part)
		switch {
		case len(reParts) == 3 && reParts[1] == "attr":
			attrName = reParts[2]
		case len(reParts) == 3 && reParts[1] == "html":
			getHTML = true
		default:
			return outSelector, attrName, getHTML, fmt.Errorf("pseudo-selector is invalid: %s", part)
		}
	}

	return outSelector, attrName, getHTML, nil
}

// GetDataSingle - extract data by one CSS-selector
//  title, err := doc.GetDataSingle("title")
func (doc Doc) GetDataSingle(selector string) (result string, err error) {
	if doc.Err != nil {
		return result, fmt.Errorf("parse document error: %s", doc.Err)
	}

	texts, err := doc.GetData(map[string]string{"single": selector})
	if err != nil {
		return result, err
	}

	if textOne, ok := texts["single"]; ok && len(textOne) > 0 {
		result = textOne[0]
	}

	return result, err
}

// FromReader - get doc from io.Reader
func FromReader(reader io.Reader) Doc {
	doc, err := goquery.NewDocumentFromReader(reader)
	return Doc{doc, err}
}

// FromFile - get doc from file
func FromFile(fileName string) Doc {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return Doc{Err: err}
	}
	defer fileReader.Close()

	return FromReader(fileReader)
}

// FromURL - get doc from URL
func FromURL(URL string) Doc {
	httpResponse, err := getHTMLPage(URL)
	if err != nil {
		return Doc{Err: err}
	}

	return FromReader(httpResponse.Body)
}

// getHTMLPage - get html by http(s) as http.Response
func getHTMLPage(url string) (response *http.Response, err error) {
	cookie, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookie,
	}

	response, err = client.Get(url)
	return response, err
}
