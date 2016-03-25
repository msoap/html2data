package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func mainWrapper(t *testing.T, args []string) (out string, err error) {
	oldStdout := os.Stdout // backup of the real stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Errorf("os.Pipe() failed")
	}
	os.Stdout = writer

	outCh := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, reader)
		outCh <- buf.String()
	}()

	os.Args = args
	config = cmdConfig{}
	err = runApp()
	writer.Close()

	out = <-outCh
	os.Stdout = oldStdout

	return strings.TrimSpace(out), err
}

func Test_main(t *testing.T) {
	// empty
	out, err := mainWrapper(t, []string{"html2data", "", "div"})
	if err != nil || !strings.HasPrefix(out, "Usage:") {
		t.Errorf("1. main() failed: got: '%s'", out)
	}

	// plain text
	out, err = mainWrapper(t, []string{"html2data", "test.html", "div.article h1"})
	if err != nil || out != "Head1\nHead2" {
		t.Errorf("2. main() failed: got: '%s'", out)
	}

	// plain text nested
	out, err = mainWrapper(t, []string{"html2data", "-find-in=div.article", "test.html", "h1"})
	if err != nil || out != "0:\nHead1\nHead2" {
		t.Errorf("3. main() failed: got: '%s'", out)
	}

	// plain text named selectors
	out, err = mainWrapper(t, []string{"html2data", "-find-in=div.article", "test.html", ":heads", "h1:get(1)", ":links", "a:attr(href)"})
	if err != nil || !(out == "0:\nheads:\tHead1\nlinks:\turl" || out == "0:\nlinks:\turl\nheads:\tHead1") {
		t.Errorf("4. main() failed: got: '%s'", out)
	}

	// json
	out, err = mainWrapper(t, []string{"html2data", "-json", "test.html", "div.article h1"})
	if err != nil || out != `{"one":["Head1","Head2"]}` {
		t.Errorf("5. main() failed: got: '%s'", out)
	}

	// json nested
	out, err = mainWrapper(t, []string{"html2data", "-json", "-find-in=div.article", "test.html", "h1"})
	if err != nil || out != `[{"one":["Head1","Head2"]}]` {
		t.Errorf("6. main() failed: got: '%s'", out)
	}

	// from URL
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<div>data</div>")
	}))
	out, err = mainWrapper(t, []string{"html2data", ts.URL, "div"})
	if err != nil || out != "data" {
		t.Errorf("7. main() failed: got: '%s'", out)
	}

	// error in args
	out, err = mainWrapper(t, []string{"html2data", "-find-in=div.article", "test.html", "h1", "h1"})
	if err == nil {
		t.Errorf("8. main() failed: got: '%s'", out)
	}

	// error in selectors
	out, err = mainWrapper(t, []string{"html2data", "test.html", "h1<<"})
	if err == nil {
		t.Errorf("9. main() failed: got: '%s'", out)
	}

	// error in selectors nested
	out, err = mainWrapper(t, []string{"html2data", "-find-in=div.article", "test.html", "h1<<"})
	if err == nil {
		t.Errorf("10. main() failed: got: '%s'", out)
	}
}
