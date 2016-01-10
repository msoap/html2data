package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/msoap/html2data"
)

const usageString = "Usage: html2data [url|file|-] 'css selector'"

func main() {
	url, CSSSelector := "", ""
	if len(os.Args) == 3 {
		// url and css selector
		url, CSSSelector = os.Args[1], os.Args[2]
	} else if len(os.Args) == 2 {
		// css selector
		CSSSelector = os.Args[1]
	} else {
		fmt.Println(usageString)
		return
	}

	var err error
	var doc html2data.Doc
	stat, _ := os.Stdin.Stat()

	if url == "-" || (stat.Mode()&os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		doc = html2data.FromReader(reader)
	} else if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		doc = html2data.FromURL(url)
	} else if len(url) > 0 {
		doc = html2data.FromFile(url)
	} else {
		fmt.Println(usageString)
		return
	}

	texts, err := doc.GetData(map[string]string{"one": CSSSelector})
	if err != nil {
		log.Fatal(err)
	}

	if textOne, ok := texts["one"]; ok {
		for _, text := range textOne {
			fmt.Println(text)
		}
	}
}
