package html2data

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_GetDataSingle(t *testing.T) {
	testData := []struct {
		html string
		css  string
		out  string
		err  error
	}{
		{
			"one<h1>head</h1>two",
			"h1",
			"head",
			nil,
		}, {
			"one<h1>  head  </h1>two",
			"h1",
			"head",
			nil,
		}, {
			"one<h1>head</h1>two<h1>head2</h1>",
			"h1",
			"head",
			nil,
		}, {
			"one<h1>head</h1>two<h1 id=2>head2</h1>",
			"h1#2",
			"head2",
			nil,
		}, {
			"one<div><h1>head</h1>two</div><h1 id=2>head2</h1>",
			"div:html",
			"<h1>head</h1>two",
			nil,
		}, {
			"one<h1>head</h1>two<a href='http://url'>link</a><h1>head2</h1>",
			"a:attr(href)",
			"http://url",
			nil,
		}, {
			"one<h1>head1</h1>two<a href='http://url'>link</a><h1>head2</h1>",
			"h1:get(2)",
			"head2",
			nil,
		}, {
			"<div>",
			"div<<<",
			"",
			nil,
		},
	}

	for i, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetDataSingle(item.css)

		if err != nil && item.err == nil {
			t.Errorf("Got error: %s", err)
		}
		if err == nil && item.err != nil {
			t.Errorf("Not got error, item: %d", i)
		}

		if out != item.out {
			t.Errorf("expected: %#v, real: %#v", item.out, out)
		}
	}
}

func Test_GetData(t *testing.T) {
	testData := []struct {
		html string
		css  map[string]string
		cfg  []Cfg
		out  map[string][]string
		err  error
	}{
		{
			html: "one<h1>head</h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{},
			out:  map[string][]string{"h1": {"head"}},
			err:  nil,
		}, {
			html: "one<h1> head </h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{},
			out:  map[string][]string{"h1": {"head"}},
			err:  nil,
		}, {
			html: "one<h1>head</h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{{DontTrimSpaces: true}},
			out:  map[string][]string{"h1": {"head"}},
			err:  nil,
		}, {
			html: "one<h1> head </h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{{DontTrimSpaces: true}},
			out:  map[string][]string{"h1": {" head "}},
			err:  nil,
		}, {
			html: "<title>Title</title>one<h1>head</h1>two<H1>Head 2</H1>",
			css:  map[string]string{"title": "title", "h1": "h1"},
			cfg:  []Cfg{},
			out:  map[string][]string{"title": {"Title"}, "h1": {"head", "Head 2"}},
			err:  nil,
		}, {
			html: "<title>Title</title>one<h1>head</h1>two<H1>Head 2</H1>",
			css:  map[string]string{"title": "title<<"},
			cfg:  []Cfg{},
			out:  map[string][]string{"title": []string{}},
			err:  nil,
		},
	}

	for i, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetData(item.css, item.cfg...)

		if err != nil && item.err == nil {
			t.Errorf("Got error: %s", err)
		}
		if err == nil && item.err != nil {
			t.Errorf("Not got error, item: %d", i)
		}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("%d. expected: %#v, real: %#v", i, item.out, out)
		}
	}
}

func Test_GetDataFirst(t *testing.T) {
	testData := []struct {
		html string
		css  map[string]string
		cfg  []Cfg
		out  map[string]string
		err  error
	}{
		{
			html: "one<h1>head</h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{},
			out:  map[string]string{"h1": "head"},
			err:  nil,
		}, {
			html: "one<h1> head </h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{},
			out:  map[string]string{"h1": "head"},
			err:  nil,
		}, {
			html: "one<h1>head</h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{{DontTrimSpaces: true}},
			out:  map[string]string{"h1": "head"},
			err:  nil,
		}, {
			html: "one<h1> head </h1>two",
			css:  map[string]string{"h1": "h1"},
			cfg:  []Cfg{{DontTrimSpaces: true}},
			out:  map[string]string{"h1": " head "},
			err:  nil,
		}, {
			html: "<title>Title</title>one<h1>head</h1>two<H1>Head 2</H1>",
			css:  map[string]string{"title": "title", "h1": "h1"},
			cfg:  []Cfg{},
			out:  map[string]string{"title": "Title", "h1": "head"},
			err:  nil,
		}, {
			html: "<title>Title</title>one<h1>head</h1>two<H1>Head 2</H1>",
			css:  map[string]string{"title": "title<<"},
			cfg:  []Cfg{},
			out:  map[string]string{"title": ""},
			err:  nil,
		}, {
			html: "<title>Title</title>one<h1>head</h1>two<H1>Head 2</H1>",
			css:  map[string]string{"title": "title", "h3": "h3"},
			cfg:  []Cfg{},
			out:  map[string]string{"title": "Title", "h3": ""},
			err:  nil,
		}, {
			html: "<title>Title</title>one<h1>head</h1>two<H1>Head 2</H1>",
			css:  map[string]string{"h2": "h2:html", "h3": "h3"},
			cfg:  []Cfg{},
			out:  map[string]string{"h2": "", "h3": ""},
			err:  nil,
		},
	}

	for i, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetDataFirst(item.css, item.cfg...)

		if err != nil && item.err == nil {
			t.Errorf("Got error: %s", err)
		}
		if err == nil && item.err != nil {
			t.Errorf("Not got error, item: %d", i)
		}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("%d. expected: %#v, real: %#v", i, item.out, out)
		}
	}
}

func Test_GetDataNested(t *testing.T) {
	testData := []struct {
		html     string
		outerCSS string
		css      map[string]string
		out      []map[string][]string
		err      error
	}{
		{
			"<div>one<h1>head</h1>two</div> <h1>head two</h1>",
			"div",
			map[string]string{"h1": "h1"},
			[]map[string][]string{{"h1": {"head"}}},
			nil,
		},
		{
			"<div>one<h1>head</h1>two</div> <div><h1>head two</h1><div>",
			"div:get(1)",
			map[string]string{"h1": "h1"},
			[]map[string][]string{{"h1": {"head"}}},
			nil,
		},
		{
			"<div>one<a href=url1>head</a>two</div> <div><a href=url2>head two</h1><div> <a href=url3>l3</a>",
			"div:get(1)",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string][]string{{"urls": {"url1"}}},
			nil,
		},
		{
			"<div>one<a href=url1>head</a>two</div> <div><a href=url2>head two</h1><div>",
			"div:get(2)",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string][]string{{"urls": {"url2"}}},
			nil,
		},
		{
			"<div class=cl>one<a href=url1>head</a>two<a href=url1.1>h1.1</a></div> <div><a href=url2>head two</a></div> <div class=cl><a href=url3>l3</a> </div>",
			"div.cl",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string][]string{{"urls": {"url1", "url1.1"}}, {"urls": {"url3"}}},
			nil,
		},
		{
			"<div class=cl>one</div>",
			"div.cl<<",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string][]string{},
			nil,
		},
		{
			"<div class=cl>one</div>",
			"div.cl",
			map[string]string{"urls": "div<<"},
			[]map[string][]string{map[string][]string{"urls": []string{}}},
			nil,
		},
	}

	for i, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetDataNested(item.outerCSS, item.css)

		if err != nil && item.err == nil {
			t.Errorf("Got error: %s", err)
		}
		if err == nil && item.err != nil {
			t.Errorf("Not got error, item: %d", i)
		}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("\n%d. html: %s\ncss: %s\nexpected: %#v\nreal    : %#v", i, item.html, item.css, item.out, out)
		}
	}
}

func Test_GetDataNestedFirst(t *testing.T) {
	testData := []struct {
		html     string
		outerCSS string
		css      map[string]string
		out      []map[string]string
		err      error
	}{
		{
			"<div>one<h1>head</h1>two</div> <h1>head two</h1>",
			"div",
			map[string]string{"h1": "h1"},
			[]map[string]string{{"h1": "head"}},
			nil,
		},
		{
			"<div class=cl>one</div>",
			"div.cl<<",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string]string{},
			nil,
		},
		{
			"<div class=cl>one</div>",
			"div.cl",
			map[string]string{"urls": "div<<"},
			[]map[string]string{{"urls": ""}},
			nil,
		},
	}

	for i, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetDataNestedFirst(item.outerCSS, item.css)

		if err != nil && item.err == nil {
			t.Errorf("Got error: %s", err)
		}
		if err == nil && item.err != nil {
			t.Errorf("Not got error, item: %d", i)
		}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("\n%d. html: %s\ncss: %s\nexpected: %#v\nreal    : %#v", i, item.html, item.css, item.out, out)
		}
	}
}

func Test_parseSelector(t *testing.T) {
	testData := []struct {
		inSelector  string
		outSelector CSSSelector
	}{
		{
			"div",
			CSSSelector{
				"div",
				"",
				false,
				0,
			},
		}, {
			"div:attr(href)",
			CSSSelector{
				"div",
				"href",
				false,
				0,
			},
		}, {
			"div: attr ( href ) ",
			CSSSelector{
				"div",
				"href",
				false,
				0,
			},
		}, {
			"div#1: attr ( href ) ",
			CSSSelector{
				"div#1",
				"href",
				false,
				0,
			},
		}, {
			"div#1:html",
			CSSSelector{
				"div#1",
				"",
				true,
				0,
			},
		}, {
			"div#1",
			CSSSelector{
				"div#1",
				"",
				false,
				0,
			},
		}, {
			"div:nth-child(1):attr(href)",
			CSSSelector{
				"div:nth-child(1)",
				"href",
				false,
				0,
			},
		}, {
			"div:nth-child(1):get(3)",
			CSSSelector{
				"div:nth-child(1)",
				"",
				false,
				3,
			},
		},
	}

	for _, item := range testData {
		outSelector := parseSelector(item.inSelector)
		inString := fmt.Sprintf("%#v", item.outSelector)
		outString := fmt.Sprintf("%#v", outSelector)

		if inString != outString {
			t.Errorf("For: %s\nexpected: %s\nreal: %s",
				item.inSelector,
				inString,
				outString,
			)
		}
	}
}

func assertDontPanic(t *testing.T, fn func(), name string) {
	defer func() {
		if recoverInfo := recover(); recoverInfo != nil {
			t.Errorf("The code panic: %s\npanic: %s", name, recoverInfo)
		}
	}()
	fn()
}

func assertPanic(t *testing.T, fn func(), name string) {
	defer func() {
		if recover() == nil {
			t.Errorf("The code did not panic: %s", name)
		}
	}()
	fn()
}

func Test_FromURL(t *testing.T) {
	assertDontPanic(t, func() { FromURL("url") }, "FromURL() with 0 arguments")
	assertDontPanic(t, func() { FromURL("url", URLCfg{}) }, "FromURL() with 1 arguments")
	assertPanic(t, func() { FromURL("url", URLCfg{}, URLCfg{}) }, "FromURL() with 2 arguments")

	// test get Url
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<div>data</div>")
	}))

	doc := FromURL(ts.URL)
	if doc.Err != nil {
		t.Errorf("Dont load url (%s): %s", ts.URL, doc.Err)
	}
	ts.Close()

	// test dont get Url
	doc = FromURL("fake://invalid/url")
	if doc.Err == nil {
		t.Errorf("Load fake url without error")
	}
	doc = FromURL("fake://%%%%/")
	if doc.Err == nil {
		t.Errorf("Load invalid url without error")
	}
	doc = FromURL("")
	if doc.Err == nil {
		t.Errorf("Load empty url without error")
	}

	// test timeout
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1200 * time.Millisecond)
		fmt.Fprintln(w, "<div>data</div>")
	}))

	doc = FromURL(ts.URL, URLCfg{TimeOut: 1})
	if doc.Err == nil {
		t.Errorf("Load url without timeout error")
	}
	ts.Close()

	// test parse
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<div><a>data</a></div>")
	}))

	doc = FromURL(ts.URL)
	if doc.Err != nil {
		t.Errorf("Dont load url, error: %s", doc.Err)
	}
	div, err := doc.GetDataSingle("div")
	if err != nil || div != "data" {
		t.Errorf("Dont load url, div: '%s', error: %s", div, doc.Err)
	}
	div, err = doc.GetDataSingle("div:html")
	if err != nil || div != "<a>data</a>" {
		t.Errorf("Dont load url, div: '%s', error: %s", div, doc.Err)
	}
	ts.Close()

	// UA test
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<div>"+r.UserAgent()+"</div><span id=2>Тест</span>")
	}))

	customUA := "CustomUA/1.0"
	doc = FromURL(ts.URL, URLCfg{UA: customUA})
	if doc.Err != nil {
		t.Errorf("Dont load url, error: %s", doc.Err)
	}
	div, err = doc.GetDataSingle("div")
	if err != nil || div != customUA {
		t.Errorf("User-agent test failed, div: '%s'", div)
	}

	doc = FromURL(ts.URL, URLCfg{DontDetectCharset: true})
	if doc.Err != nil {
		t.Errorf("Dont load url, error: %s", doc.Err)
	}
	span, err := doc.GetDataSingle("span#2")
	if err != nil || span != "Тест" {
		t.Errorf("DontDetectCharset failed, span: '%s'", div)
	}
	ts.Close()
}

func Test_FromFile(t *testing.T) {
	doc := FromFile("/dont exists file")
	_, err := doc.GetDataSingle("div")
	if err == nil {
		t.Errorf("FromFile(): open dont exists file")
	}
}

func Test_getConfig(t *testing.T) {
	cfg := getConfig([]Cfg{})
	if !reflect.DeepEqual(cfg, Cfg{}) {
		t.Errorf("1. getConfig(): empty list")
	}

	cfg = getConfig([]Cfg{{}})
	if !reflect.DeepEqual(cfg, Cfg{}) {
		t.Errorf("2. getConfig(): list with one element")
	}

	cfg = getConfig([]Cfg{{DontTrimSpaces: true}})
	if !reflect.DeepEqual(cfg, Cfg{DontTrimSpaces: true}) {
		t.Errorf("3. getConfig(): list with one element with true")
	}

	assertPanic(t, func() {
		getConfig([]Cfg{{}, {}})
	}, "3. getConfig() must panic")
}

func ExampleFromURL() {
	doc := FromURL("http://example.com")
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}

	// or with config
	doc = FromURL("http://example.com", URLCfg{UA: "userAgent", TimeOut: 10, DontDetectCharset: false})
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}
}

func ExampleFromFile() {
	doc := FromFile("file_name.html")
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}
}

func ExampleFromReader() {
	doc := FromReader(bufio.NewReader(os.Stdin))
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}
}

func ExampleDoc_GetDataSingle() {
	// get title
	title, err := FromFile("cmd/html2data/test.html").GetDataSingle("title")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Title is:", title)
	// Output: Title is: Title
}

func ExampleDoc_GetData() {
	texts, _ := FromURL("http://example.com").GetData(map[string]string{"headers": "h1", "links": "a:attr(href)"})
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
	texts, err := FromURL("http://example.com").GetDataFirst(map[string]string{"header": "h1", "first_link": "a:attr(href)"})
	if err != nil {
		log.Fatal(err)
	}

	// get H1 header:
	fmt.Println("header: ", texts["header"])
	// get URL in first link:
	fmt.Println("first link: ", texts["first_link"])
}

func ExampleDoc_GetDataNested() {
	texts, _ := FromFile("test.html").GetDataNested("div.article", map[string]string{"headers": "h1", "links": "a:attr(href)"})
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
	texts, err := FromFile("cmd/html2data/test.html").GetDataNestedFirst("div.block", map[string]string{"header": "h1", "link": "a:attr(href)", "sp": "span"})
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
	doc := FromURL("http://example.com")
	// or with config
	// doc := FromURL("http://example.com", URLCfg{UA: "userAgent", TimeOut: 10, DontDetectCharset: true})
	if doc.Err != nil {
		log.Fatal(doc.Err)
	}

	// get title
	title, _ := doc.GetDataSingle("title")
	fmt.Println("Title is:", title)

	title, _ = doc.GetDataSingle("title", Cfg{DontTrimSpaces: true})
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
