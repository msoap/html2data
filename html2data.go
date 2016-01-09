package main

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

func get_html_page(weather_url string) *http.Response {
	cookie, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookie,
	}

	response, err := client.Get(weather_url)

	if err != nil {
		log.Fatal(err)
	}

	return response
}

func getData(url string, selectors map[string]string) (result map[string][]string) {
	var err error
	var doc *goquery.Document
	result = map[string][]string{}

	stat, _ := os.Stdin.Stat()
	if url == "-" || (stat.Mode()&os.ModeCharDevice) == 0 {
		stdinReader := bufio.NewReader(os.Stdin)
		doc, err = goquery.NewDocumentFromReader(stdinReader)
	} else {
		http_response := get_html_page(url)
		b, _ := ioutil.ReadAll(http_response.Body)
		fmt.Println(string(b))
		doc, err = goquery.NewDocumentFromReader(http_response.Body)
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

// go run html2data.go url "tr.quote__day:nth-child(2) td.quote__value"
func main() {
	texts := getData(os.Args[1], map[string]string{"one": os.Args[2]})
	for _, text := range texts["one"] {
		fmt.Println(text)
	}
}
