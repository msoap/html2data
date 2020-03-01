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
	"  html2data [options] [url|file|-] :name1 'css1' :name2 'css2' ...\n\n" +
	"options:"

type cmdConfig struct {
	userAgent, outerCSS, url string
	timeOut                  int
	getJSON                  bool
	dontTrimSpaces           bool
	dontDetectCharset        bool
}

var (
	config cmdConfig
)

func init() {
	flag.StringVar(&config.userAgent, "user-agent", "", "set custom user-agent")
	flag.StringVar(&config.outerCSS, "find-in", "", "search in the specified elements instead document")
	flag.BoolVar(&config.getJSON, "json", false, "JSON output")
	flag.BoolVar(&config.dontTrimSpaces, "dont-trim-spaces", false, "don't trim spaces, get text as is")
	flag.BoolVar(&config.dontDetectCharset, "dont-detect-charset", false, "don't detect charset and convert text")
	flag.IntVar(&config.timeOut, "timeout", 0, "timeout in seconds")
}

func getConfig() (CSSSelectors map[string]string, err error) {
	flag.Usage = func() {
		fmt.Println(usageString)
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()

	config.url, CSSSelectors, err = parseArgs(flag.Args())
	return CSSSelectors, err
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

func runApp() error {
	CSSSelectors, err := getConfig()
	if err != nil {
		return err
	}
	var doc html2data.Doc
	stat, err := os.Stdin.Stat()
	if err != nil {
		return err
	}

	if config.url == "-" || (stat.Mode()&os.ModeCharDevice) == 0 {
		reader := bufio.NewReader(os.Stdin)
		doc = html2data.FromReader(reader)
	} else if strings.HasPrefix(config.url, "http://") || strings.HasPrefix(config.url, "https://") {
		doc = html2data.FromURL(config.url, html2data.URLCfg{UA: config.userAgent, TimeOut: config.timeOut, DontDetectCharset: config.dontDetectCharset})
	} else if len(config.url) > 0 {
		doc = html2data.FromFile(config.url)
	} else {
		fmt.Println(usageString)
		return nil
	}

	GetDocCfg := html2data.Cfg{DontTrimSpaces: config.dontTrimSpaces}
	if config.outerCSS != "" {
		textsOuter, err := doc.GetDataNested(config.outerCSS, CSSSelectors, GetDocCfg)
		if err != nil {
			return err
		}

		if config.getJSON {
			jsonBytes, err := json.Marshal(textsOuter)
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
		} else {
			for i, texts := range textsOuter {
				fmt.Printf("%d:\n", i)
				printAsText(texts, len(CSSSelectors) > 1)
			}
		}
	} else {
		texts, err := doc.GetData(CSSSelectors, GetDocCfg)
		if err != nil {
			return err
		}

		if config.getJSON {
			jsonBytes, err := json.Marshal(texts)
			if err != nil {
				return err
			}
			fmt.Println(string(jsonBytes))
		} else {
			printAsText(texts, len(CSSSelectors) > 1)
		}
	}

	return nil
}

func main() {
	err := runApp()
	if err != nil {
		log.Fatal(err)
	}
}
