package main

import (
	"bufio"
	"encoding/json"
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

func getConfig() (userAgent, outerCSS, url string, getJSON bool, CSSSelectors map[string]string) {
	flag.StringVar(&userAgent, "user-agent", "", "set custom user-agent")
	flag.StringVar(&outerCSS, "find-in", "", "search in the specified elements instead document")
	flag.BoolVar(&getJSON, "json", false, "JSON output")
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
	return userAgent, outerCSS, url, getJSON, CSSSelectors
}

// printAsText - print result as text
func printAsText(texts map[string][]string, doPrintName bool) {
	for name, value := range texts {
		if doPrintName {
			fmt.Print(name + ":\t")
		}
		for _, text := range value {
			fmt.Println(text)
		}
	}
}

func main() {
	userAgent, outerCSS, url, getJSON, CSSSelectors := getConfig()
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

		if getJSON {
			jsonObject := []map[string][]string{}
			for _, texts := range textsOuter {
				jsonObject = append(jsonObject, texts)
			}
			json, _ := json.Marshal(jsonObject)
			fmt.Println(string(json))
		} else {
			for i, texts := range textsOuter {
				fmt.Printf("%d:\n", i)
				printAsText(texts, len(CSSSelectors) > 1)
			}
		}
	} else {
		texts, err := doc.GetData(CSSSelectors)
		if err != nil {
			log.Fatal(err)
		}

		if getJSON {
			json, _ := json.Marshal(texts)
			fmt.Println(string(json))
		} else {
			printAsText(texts, len(CSSSelectors) > 1)
		}
	}
}
