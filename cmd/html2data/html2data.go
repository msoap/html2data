package main

import (
	"fmt"
	"os"

	"github.com/msoap/html2data"
)

func main() {
	texts := html2data.GetData(os.Args[1], map[string]string{"one": os.Args[2]})
	for _, text := range texts["one"] {
		fmt.Println(text)
	}
}
