package html2data

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
			out:  map[string][]string{"title": {}},
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
			[]map[string][]string{{"urls": {}}},
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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintln(w, "<div>data</div>")
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
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		time.Sleep(1200 * time.Millisecond)
		_, _ = fmt.Fprintln(w, "<div>data</div>")
	}))

	doc = FromURL(ts.URL, URLCfg{TimeOut: 1})
	if doc.Err == nil {
		t.Errorf("Load url without timeout error")
	}
	ts.Close()

	// test parse
	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = fmt.Fprintln(w, "<div><a>data</a></div>")
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
		_, _ = fmt.Fprintln(w, "<div>"+r.UserAgent()+"</div><span id=2>Тест</span>")
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

	doc = FromURL("fake://url")
	_, err = doc.GetDataNested("", map[string]string{})
	if err == nil {
		t.Error("GetDataNested not got error on fake URL")
	}
	_, err = doc.GetDataNestedFirst("", map[string]string{})
	if err == nil {
		t.Error("GetDataNestedFirst not got error on fake URL")
	}
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
