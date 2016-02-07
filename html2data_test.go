package html2data

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_GetDataSingle(t *testing.T) {
	testData := []struct {
		html string
		css  string
		out  string
	}{
		{
			"one<h1>head</h1>two",
			"h1",
			"head",
		}, {
			"one<h1>head</h1>two<h1>head2</h1>",
			"h1",
			"head",
		}, {
			"one<h1>head</h1>two<h1 id=2>head2</h1>",
			"h1#2",
			"head2",
		}, {
			"one<div><h1>head</h1>two</div><h1 id=2>head2</h1>",
			"div:html",
			"<h1>head</h1>two",
		}, {
			"one<h1>head</h1>two<a href='http://url'>link</a><h1>head2</h1>",
			"a:attr(href)",
			"http://url",
		}, {
			"one<h1>head1</h1>two<a href='http://url'>link</a><h1>head2</h1>",
			"h1:get(2)",
			"head2",
		},
	}

	for _, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetDataSingle(item.css)

		if err != nil {
			t.Errorf("Got error: %s", err)
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
		out  map[string][]string
	}{
		{
			"one<h1>head</h1>two",
			map[string]string{"h1": "h1"},
			map[string][]string{"h1": {"head"}},
		}, {
			"<title>Title</title>one<h1>head</h1>two<H1>Head 2</H1>",
			map[string]string{"title": "title", "h1": "h1"},
			map[string][]string{"title": {"Title"}, "h1": {"head", "Head 2"}},
		},
	}

	for _, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetData(item.css)

		if err != nil {
			t.Errorf("Got error: %s", err)
		}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("expected: %#v, real: %#v", item.out, out)
		}
	}
}

func Test_GetDataNested(t *testing.T) {
	testData := []struct {
		html     string
		outerCss string
		css      map[string]string
		out      []map[string][]string
	}{
		{
			"<div>one<h1>head</h1>two</div> <h1>head two</h1>",
			"div",
			map[string]string{"h1": "h1"},
			[]map[string][]string{{"h1": {"head"}}},
		},
		{
			"<div>one<h1>head</h1>two</div> <div><h1>head two</h1><div>",
			"div:get(1)",
			map[string]string{"h1": "h1"},
			[]map[string][]string{{"h1": {"head"}}},
		},
		{
			"<div>one<a href=url1>head</a>two</div> <div><a href=url2>head two</h1><div> <a href=url3>l3</a>",
			"div:get(1)",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string][]string{{"urls": {"url1"}}},
		},
		{
			"<div>one<a href=url1>head</a>two</div> <div><a href=url2>head two</h1><div>",
			"div:get(2)",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string][]string{{"urls": {"url2"}}},
		},
		{
			"<div class=cl>one<a href=url1>head</a>two<a href=url1.1>h1.1</a></div> <div><a href=url2>head two</a></div> <div class=cl><a href=url3>l3</a> </div>",
			"div.cl",
			map[string]string{"urls": "a:attr(href)"},
			[]map[string][]string{{"urls": {"url1", "url1.1"}}, {"urls": {"url3"}}},
		},
	}

	for _, item := range testData {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetDataNested(item.outerCss, item.css)

		if err != nil {
			t.Errorf("Got error: %s", err)
		}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("\nhtml: %s\ncss: %s\nexpected: %#v\nreal    : %#v", item.html, item.css, item.out, out)
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
