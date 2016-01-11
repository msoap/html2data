package html2data

import (
	"reflect"
	"strings"
	"testing"
)

func Test_GetDataSingle(t *testing.T) {
	test_data := []struct {
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
		},
	}

	for _, item := range test_data {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetDataSingle(item.css)

		if err != nil {
			t.Errorf("Got error:", err)
		}

		if out != item.out {
			t.Errorf("expected: %#v, real: %#v", item.out, out)
		}
	}
}

func Test_GetData(t *testing.T) {
	test_data := []struct {
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

	for _, item := range test_data {
		reader := strings.NewReader(item.html)
		out, err := FromReader(reader).GetData(item.css)

		if err != nil {
			t.Errorf("Got error:", err)
		}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("expected: %#v, real: %#v", item.out, out)
		}
	}
}
