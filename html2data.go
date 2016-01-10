package html2data

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// GetDataFromReader - extract data by CSS selectors from io.Reader
func GetDataFromReader(reader io.Reader, selectors map[string]string) (result map[string][]string, err error) {
	result = map[string][]string{}
	doc, err := goquery.NewDocumentFromReader(reader)

	if err != nil {
		return result, fmt.Errorf("GetData error: ", err)
	}

	for name, selector := range selectors {
		texts := []string{}
		doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
			texts = append(texts, selection.Text())
		})
		result[name] = texts
	}

	return result, err
}

// GetDataFromFile - extract data by CSS selectors from file
func GetDataFromFile(fileName string, selectors map[string]string) (result map[string][]string, err error) {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return result, err
	}
	defer fileReader.Close()

	return GetDataFromReader(fileReader, selectors)
}

// GetDataFromURL - extract data by CSS selectors from URL
func GetDataFromURL(URL string, selectors map[string]string) (result map[string][]string, err error) {
	httpResponse, err := getHTMLPage(URL)
	if err != nil {
		return result, err
	}

	return GetDataFromReader(httpResponse.Body, selectors)
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
