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

	// Set the author
	e.SetAuthor("Hingle McCringleberry")

	// Add a section
	section1Body := `<h1>Section 1</h1>
	<p>This is a paragraph.</p>`
	parent1, err := e.AddSection(section1Body, "Section 1", "", "")
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Add Section\n")

	// Add a section
	section11Body := `<h2>Section 1-1</h2>
	<p>This is a paragraph of h2.</p>`
	parent11, err := e.AddSubSection(parent1, section11Body, "Section 1-1", "SubSectionName", "")
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Add subSection\n")
	fmt.Printf("SubSection name: %v\n", parent11)

	// Add a section
	section111Body := `<h3>Section 1-1-1</h3>
	<p>This is a paragraph of h3.</p>`
	_, err = e.AddSubSection(parent11, section111Body, "Section 1-1-1", "", "")
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
