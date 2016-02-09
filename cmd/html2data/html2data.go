package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/msoap/html2data"
)

const usageString = "Usage: html2data [options] [url|file|-] 'css selector'"

func main() {
	userAgent := ""
	flag.StringVar(&userAgent, "user-agent", "", "set custom user-agent")
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	url, CSSSelector := "", ""
	args := flag.Args()
	if len(args) == 2 {
		// url and css selector
		url, CSSSelector = args[0], args[1]
	} else if len(args) == 1 {
		// css selector
		CSSSelector = args[0]
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
		doc = html2data.FromURL(url, html2data.Cfg{UA: userAgent})
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
