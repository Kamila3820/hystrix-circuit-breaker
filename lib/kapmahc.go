package lib

import (
	"fmt"
	"log"

	"github.com/kapmahc/epub"
)

func KapMahc(file string) {
	book, err := epub.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer book.Close()

	// List table of contents
	toc := book.Ncx.Points
	var count int
	for i, t := range toc {
		// if t.Text == "" {
		// 	continue
		// }
		fmt.Printf("Chapter %d: %s\n", i+1, t.Text)
		count++
	}
	fmt.Printf("Total Chapters: %d\n", count)
}
