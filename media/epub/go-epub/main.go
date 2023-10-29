package main

import (
	"fmt"
	"log"

	epub "github.com/go-shiori/go-epub"
)

func main() {
	// Create a new EPUB
	e, err := epub.NewEpub("My title")
	if err != nil {
		log.Println(err)
	}

	// Add a section
	sectionBody := `<h1>Section 1</h1>
	<p>This is a paragraph.</p>`
	parent1, err := e.AddSection(sectionBody, "Section 1", "", "")
	if err != nil {
		log.Println(err)
	}

	// Add a subsection
	subsectionBody := `<h2>Section 1-1</h2>
	<p>This is a paragraph of h2.</p>`
	_, err = e.AddSubSection(parent1, subsectionBody, "Section 1-1", "SubSectionName", "")
	if err != nil {
		log.Println(err)
	}

	// Write the EPUB
	err = e.Write("My EPUB.epub")
	if err != nil {
		// handle error
		fmt.Printf("Falied to create epub, err: %v\n", err)
	}
}
