package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Test_main(t *testing.T) {
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

	os.Args = []string{"html2data", "-json", "-find-in=div.article", "test.html", "h1"}
	main()
	writer.Close()

	out := <-outCh
	os.Stdout = oldStdout
	if strings.TrimSpace(out) != `[{"one":["Head1","Head2"]}]` {
		t.Errorf("main() failed: got: '%s'", strings.TrimSpace(out))
	}
	// fmt.Printf("out: <%s>\n", out)
}
