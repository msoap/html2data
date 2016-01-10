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

	var texts map[string][]string
	var err error
	stat, _ := os.Stdin.Stat()

	if url == "-" || (stat.Mode()&os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		texts, err = html2data.GetDataFromReader(reader, map[string]string{"one": CSSSelector})
	} else if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		texts, err = html2data.GetDataFromURL(url, map[string]string{"one": CSSSelector})
	} else if len(url) > 0 {
		texts, err = html2data.GetDataFromFile(url, map[string]string{"one": CSSSelector})
	} else {
		fmt.Println(usageString)
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	if textOne, ok := texts["one"]; ok {
		for _, text := range textOne {
			fmt.Println(text)
		}
	}
}
