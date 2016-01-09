package html2data

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func getHTMLPage(url string) *http.Response {
	cookie, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookie,
	}

	response, err := client.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	return response
}

// GetData - extract data by CSS selectors from URL or HTML page
func GetData(url string, selectors map[string]string) (result map[string][]string) {
	var err error
	var doc *goquery.Document
	result = map[string][]string{}

	stat, _ := os.Stdin.Stat()
	if url == "-" || (stat.Mode()&os.ModeCharDevice) == 0 {
		stdinReader := bufio.NewReader(os.Stdin)
		doc, err = goquery.NewDocumentFromReader(stdinReader)
	} else {
		httpResponse := getHTMLPage(url)
		doc, err = goquery.NewDocumentFromReader(httpResponse.Body)
	}

	if err != nil {
		log.Fatal(err)
	}

	for name, selector := range selectors {
		texts := []string{}
		doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
			texts = append(texts, selection.Text())
		})
		result[name] = texts
	}

	return result
}
