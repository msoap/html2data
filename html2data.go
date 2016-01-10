/*
Package html2data - extract data from HTML via CSS selectors

Install package and command line utility:

	go get -u github.com/msoap/html2data/cmd/html2data

Install package only:

	go get -u github.com/msoap/html2data

*/
package html2data

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// Doc - html document for parse
type Doc struct {
	doc *goquery.Document
	Err error
}

// GetData - extract data by CSS selectors
//  texts, err := doc.GetData(map[string]string{"h1": "h1"})
func (doc Doc) GetData(selectors map[string]string) (result map[string][]string, err error) {
	if doc.Err != nil {
		return result, fmt.Errorf("parse document error: %s", doc.Err)
	}

	result = map[string][]string{}
	for name, selector := range selectors {
		texts := []string{}
		doc.doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
			texts = append(texts, selection.Text())
		})
		result[name] = texts
	}

	return result, err
}

// GetDataSingle - extract data by CSS selector
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
