package lib

import (
	"fmt"
	"log"

	"github.com/gen2brain/go-fitz"
)

func Fitz(file string) {
	doc, err := fitz.New(file)
	if err != nil {
		log.Fatalf("Error opening document: %v", err)
	}
	defer doc.Close()

	toc, err := doc.ToC()
	if err != nil {
		log.Fatalf("Error extracting ToC: %v", err)
	}

	if len(toc) == 0 {
		fmt.Println("No Table of Contents found in the EPUB file.")
		return
	}

	var count int
	fmt.Println("Table of Contents:")
	for i, entry := range toc {
		fmt.Printf("Chapter %d: %s\n", i+1, entry.Title)
		count++
	}
	fmt.Printf("Total Contents: %d\n", count)
}

// func FitzZip(file string) {
// 	doc, err := fitz.New(file)
// 	if err != nil {
// 		log.Fatalf("Error opening document: %v", err)
// 	}
// 	defer doc.Close()

// 	toc, err := doc.ToC()
// 	if err != nil {
// 		log.Fatalf("Error extracting ToC: %v", err)
// 	}

// 	if len(toc) == 0 {
// 		fmt.Println("No Table of Contents found in the EPUB file.")
// 		return
// 	}

// 	fmt.Println("Table of Contents:")
// 	for _, entry := range toc {
// 		fmt.Printf("Title: %s\n", entry.Title)
// 	}
// }
