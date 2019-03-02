package main

import (
	"fmt"
	"reflect"
	"testing"
)

// func parseArgs(args []string) (url string, selectors map[string]string) {
type parseArgsResult struct {
	url       string
	selectors map[string]string
	err       string
}

func Test_parseArgs(t *testing.T) {
	testData := []struct {
		in  []string
		out parseArgsResult
	}{
		{
			in: []string{},
			out: parseArgsResult{
				url:       "",
				selectors: nil,
				err:       "arguments is empty",
			},
		},
		{
			in: []string{"div"},
			out: parseArgsResult{
				url:       "-",
				selectors: map[string]string{"one": "div"},
				err:       "",
			},
		},
		{
			in: []string{":name", "div"},
			out: parseArgsResult{
				url:       "-",
				selectors: map[string]string{"name": "div"},
				err:       "",
			},
		},
		{
			in: []string{"http://url", ":name", "div"},
			out: parseArgsResult{
				url:       "http://url",
				selectors: map[string]string{"name": "div"},
				err:       "",
			},
		},
		{
			in: []string{"http://url", "div"},
			out: parseArgsResult{
				url:       "http://url",
				selectors: map[string]string{"one": "div"},
				err:       "",
			},
		},
		{
			in: []string{":name1", "div1", ":name2", "div2"},
			out: parseArgsResult{
				url:       "-",
				selectors: map[string]string{"name1": "div1", "name2": "div2"},
				err:       "",
			},
		},
		{
			in: []string{"file", ":name1", "div1", ":name2", "div2"},
			out: parseArgsResult{
				url:       "file",
				selectors: map[string]string{"name1": "div1", "name2": "div2"},
				err:       "",
			},
		},
		{
			in: []string{"file", ":name1", "div1", "name2", "div2"},
			out: parseArgsResult{
				url:       "",
				selectors: nil,
				err:       fmt.Sprintf("name '%s' is not valid, must begin from ':'", "name2"),
			},
		},
	}

	for i, item := range testData {
		var selectors map[string]string
		url, selectors, err := parseArgs(item.in)
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		out := parseArgsResult{url, selectors, errMsg}

		if !reflect.DeepEqual(item.out, out) {
			t.Errorf("\n%d. expected: %#v\n       real: %#v", i, item.out, out)
		}
	}
}
