package html2data_test

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/msoap/html2data"
)

func ExampleFromURL() {
	doc := html2data.FromURL("http://example.com")
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}

	// or with config
	doc = html2data.FromURL("http://example.com", html2data.URLCfg{UA: "userAgent", TimeOut: 10, DontDetectCharset: false})
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}
}

func ExampleFromFile() {
	doc := html2data.FromFile("file_name.html")
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}
}

func ExampleFromReader() {
	doc := html2data.FromReader(bufio.NewReader(os.Stdin))
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}
}

func ExampleDoc_GetDataSingle() {
	// get title
	title, err := html2data.FromFile("cmd/html2data/test.html").GetDataSingle("title")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title is:", title)
	// Output: Title is: Title
}

func ExampleDoc_GetData() {
	texts, _ := html2data.FromURL("http://example.com").GetData(map[string]string{"headers": "h1", "links": "a:attr(href)"})
	// get all H1 headers:
	if textOne, ok := texts["headers"]; ok {
		for _, text := range textOne {
			fmt.Println(text)
		}
	}
	// get all urls from links
	if links, ok := texts["links"]; ok {
		for _, text := range links {
			fmt.Println(text)
		}
	}
}

func ExampleDoc_GetDataFirst() {
	texts, err := html2data.FromURL("http://example.com").GetDataFirst(map[string]string{"header": "h1", "first_link": "a:attr(href)"})
	if err != nil {
		log.Fatal(err)
	}

	// get H1 header:
	fmt.Println("header: ", texts["header"])
	// get URL in first link:
	fmt.Println("first link: ", texts["first_link"])
}

func ExampleDoc_GetDataNested() {
	texts, _ := html2data.FromFile("test.html").GetDataNested("div.article", map[string]string{"headers": "h1", "links": "a:attr(href)"})
	for _, article := range texts {
		// get all H1 headers inside each <div class="article">:
		if textOne, ok := article["headers"]; ok {
			for _, text := range textOne {
				fmt.Println(text)
			}
		}
		// get all urls from links inside each <div class="article">
		if links, ok := article["links"]; ok {
			for _, text := range links {
				fmt.Println(text)
			}
		}
	}
}

func ExampleDoc_GetDataNestedFirst() {
	texts, err := html2data.FromFile("cmd/html2data/test.html").GetDataNestedFirst("div.block", map[string]string{"header": "h1", "link": "a:attr(href)", "sp": "span"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("")
	for _, block := range texts {
		// get first H1 header
		fmt.Printf("header - %s\n", block["header"])

		// get first link
		fmt.Printf("first URL - %s\n", block["link"])

		// get not exists span
		fmt.Printf("span - '%s'\n", block["span"])
	}

	// Output:
	// header - Head1.1
	// first URL - http://url1
	// span - ''
	// header - Head2.1
	// first URL - http://url2
	// span - ''
}

func Example() {
	doc := html2data.FromURL("http://example.com")
	// or with config
	// doc := FromURL("http://example.com", URLCfg{UA: "userAgent", TimeOut: 10, DontDetectCharset: true})
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}

	// get title
	title, _ := doc.GetDataSingle("title")
	fmt.Println("Title is:", title)

	title, _ = doc.GetDataSingle("title", html2data.Cfg{DontTrimSpaces: true})
	fmt.Println("Title as is, with spaces:", title)

	texts, _ := doc.GetData(map[string]string{"h1": "h1", "links": "a:attr(href)"})
	// get all H1 headers:
	if textOne, ok := texts["h1"]; ok {
		for _, text := range textOne {
			fmt.Println(text)
		}
	}
	// get all urls from links
	if links, ok := texts["links"]; ok {
		for _, text := range links {
			fmt.Println(text)
		}
	}
}
