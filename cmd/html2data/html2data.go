package main

import (
	"fmt"
	"log"
	"os"

	"github.com/msoap/html2data"
)

func main() {
	texts, err := html2data.GetData(os.Args[1], map[string]string{"one": os.Args[2]})
	if err != nil {
		log.Fatal(err)
	}

	if textOne, ok := texts["one"]; ok {
		for _, text := range textOne {
			fmt.Println(text)
		}
	}
}
