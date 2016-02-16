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

const usageString = "Usage:\n" +
	"  html2data [options] [url|file|-] 'css selector'\n" +
	"  html2data [options] [url|file|-] :name 'css1' :name2 'css2' ...\n\n" +
	"options:"

func getConfig() (userAgent, outerCSS, url string, CSSSelectors map[string]string) {
	flag.StringVar(&userAgent, "user-agent", "", "set custom user-agent")
	flag.StringVar(&outerCSS, "find-in", "", "search in the specified elements instead document")
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	url, CSSSelectors, err := parseArgs(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	return userAgent, outerCSS, url, CSSSelectors
}

func main() {
	userAgent, outerCSS, url, CSSSelectors := getConfig()
	var doc html2data.Doc
	stat, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

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

	if outerCSS != "" {
		textsOuter, err := doc.GetDataNested(outerCSS, CSSSelectors)
		if err != nil {
			log.Fatal(err)
		}

		for i, texts := range textsOuter {
			fmt.Printf("%d:\n", i)
			for name, value := range texts {
				if len(CSSSelectors) > 1 {
					fmt.Print(name + ":\t")
				}
				for _, text := range value {
					fmt.Println(text)
				}
			}
		}
	} else {
		texts, err := doc.GetData(CSSSelectors)
		if err != nil {
			log.Fatal(err)
		}

		for name, value := range texts {
			if len(CSSSelectors) > 1 {
				fmt.Print(name + ":\t")
			}
			for _, text := range value {
				fmt.Println(text)
			}
		}
	}
}
