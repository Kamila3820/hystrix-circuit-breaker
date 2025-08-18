package lib

import (
	"archive/zip"
	"fmt"
	"log"

	"github.com/timsims/pamphlet"
)

const (
	HTMLMimetype = "application/xhtml+xml"
	SVGMimetype  = "image/svg+xml"
)

type ZipFile struct {
	file *zip.File
}

type Book struct {
	Title       string
	Author      string
	Language    string
	Identifier  string
	Publisher   string
	Description string
	Subject     string
	Date        string
	Chapters    []Chapter
}

type Chapter struct {
	ZipFile
	ID        string
	Title     string
	Href      string
	MediaType string
	HasToc    bool
	Order     int
}

func Pamphlet(file string) {
	parser, err := pamphlet.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	book := parser.GetBook()
	fmt.Printf("Title: %s\n", book.Title)
	fmt.Printf("Author: %s\n", book.Author)

	chapters := book.Chapters
	var count int
	for i, chapter := range chapters {
		// if chapter.Title == "" {
		// 	continue
		// }
		fmt.Printf("Chapter %d: %s, Order: %s\n", i+1, chapter.Title, chapter.ID)
		count++
	}
	fmt.Printf("Total chapters: %d\n", count)
}
