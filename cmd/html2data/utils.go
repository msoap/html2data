package main

import (
	"fmt"
	"strings"
)

/*
Args variants:

	http://url css
	file css
	css

	http://url :name css
	file :name css
	:name css

	http://url :name css :name css
	file :name css :name css
	:name css :name css
*/

func parseArgs(args []string) (url string, selectors map[string]string, err error) {
	var tail []string
	selectors = map[string]string{}

	switch {
	case len(args) == 0:
		return "", nil, fmt.Errorf("arguments is empty")
	case len(args) == 1:
		selectors["one"] = args[0]
		url = "-"
		return url, selectors, err
	case len(args) == 2 && strings.HasPrefix(args[0], ":"):
		selectors[strings.TrimLeft(args[0], ":")] = args[1]
		url = "-"
	case len(args) == 2 && !strings.HasPrefix(args[0], ":"):
		selectors["one"] = args[1]
		url = args[0]
	case len(args)%2 == 0:
		// even arguments
		url = "-"
		tail = args[:]
	case len(args)%2 != 0:
		// not even arguments
		url = args[0]
		tail = args[1:]
	}

	for i := 0; i < len(tail); i += 2 {
		name := tail[i]
		if !strings.HasPrefix(name, ":") {
			return "", nil, fmt.Errorf("name '%s' is not valid, must begin from ':'", name)
		}
		name = strings.TrimLeft(name, ":")
		selectors[name] = tail[i+1]
	}

	return url, selectors, err
}
