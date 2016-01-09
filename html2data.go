package html2data

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func getHTMLPage(url string) (response *http.Response, err error) {
	cookie, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookie,
	}

	response, err = client.Get(url)

	if err != nil {
		return response, err
	}

	return response, err
}

// GetData - extract data by CSS selectors from URL or HTML page
func GetData(url string, selectors map[string]string) (result map[string][]string, err error) {
	var doc *goquery.Document
	result = map[string][]string{}

	stat, _ := os.Stdin.Stat()
	if url == "-" || (stat.Mode()&os.ModeCharDevice) == 0 {
		stdinReader := bufio.NewReader(os.Stdin)
		doc, err = goquery.NewDocumentFromReader(stdinReader)
	} else {
		httpResponse, err := getHTMLPage(url)
		if err != nil {
			return result, fmt.Errorf("GetData error: %s", err)
		}

		doc, err = goquery.NewDocumentFromReader(httpResponse.Body)
	}

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
