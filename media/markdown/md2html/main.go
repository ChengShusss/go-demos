package main

import (
	"fmt"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

func main() {
	input, err := os.ReadFile("input.md")
	if err != nil {
		panic(err)
	}

	output := markdown.ToHTML(input, nil, &html.Renderer{})

	fmt.Println(string(output))
}
