package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"hystrix-circuit-breaker/lib"
	"io"
	"log"
	"strings"

	"github.com/kapmahc/epub"
	"golang.org/x/net/html"
)

// EPUB NCX structure
type NCX struct {
	NavMap NavMap `xml:"navMap"`
}

type NavMap struct {
	NavPoints []NavPoint `xml:"navPoint"`
}

type NavPoint struct {
	NavLabel NavLabel   `xml:"navLabel"`
	Children []NavPoint `xml:"navPoint"` // nested chapters
}

type NavLabel struct {
	Text string `xml:"text"`
}

func main() {
	// Epub2Contents("read.gz")
	// Epub2And3Contents("emoji.epub")

	// lib.KapMahc("world.epub")
	// lib.Pamphlet("epub2.epub")
	lib.Fitz("31999_release_admin.epub")

	// ReadContents("read.gz")
	// CompressFile("read.txt", "read.gz")
}

func Epub2Contents(file string) {
	epubPath := file

	r, err := zip.OpenReader(epubPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var tocFile *zip.File
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "toc.ncx") {
			tocFile = f
			break
		}
	}
	if tocFile == nil {
		log.Fatal("No toc.ncx found in EPUB")
	}

	rc, err := tocFile.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}

	var ncx NCX
	if err := xml.Unmarshal(data, &ncx); err != nil {
		log.Fatal(err)
	}

	// Flatten into []string
	var contents []string
	var walk func([]NavPoint)
	walk = func(points []NavPoint) {
		for _, np := range points {
			contents = append(contents, np.NavLabel.Text)
			if len(np.Children) > 0 {
				walk(np.Children)
			}
		}
	}
	walk(ncx.NavMap.NavPoints)

	fmt.Println("Contents:")
	var count int
	for _, c := range contents {
		count++
		fmt.Println("-", c)
	}
	fmt.Println("Total chapters:", count)
}

func Epub2And3Contents(file string) {
	epubPath := file
	fmt.Println(file)

	r, err := zip.OpenReader(epubPath)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var tocFile *zip.File
	var tocType string
	for _, f := range r.File {
		fmt.Println(f.Name)
		switch {
		case strings.HasSuffix(f.Name, "toc.ncx"):
			tocFile = f
			tocType = "ncx"
		case strings.HasSuffix(f.Name, "nav.xhtml"),
			strings.HasSuffix(f.Name, "toc.xhtml"),
			strings.HasSuffix(f.Name, "nav.html"):
			tocFile = f
			tocType = "xhtml"
		}
	}
	fmt.Println("TOC FROM ZIP -> ", tocFile.Name)
	if tocFile == nil {
		log.Fatal("No toc.ncx or nav.xhtml found in EPUB")
	}

	rc, err := tocFile.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		log.Fatal(err)
	}

	// Parser based on tocType
	var contents []string

	if tocType == "ncx" {
		// EPUB 2
		var ncx NCX
		if err := xml.Unmarshal(data, &ncx); err != nil {
			log.Fatal(err)
		}
		var walk func([]NavPoint)
		walk = func(points []NavPoint) {
			for _, np := range points {
				contents = append(contents, np.NavLabel.Text)
				if len(np.Children) > 0 {
					walk(np.Children)
				}
			}
		}
		walk(ncx.NavMap.NavPoints)

	} else if tocType == "xhtml" {
		// EPUB 3
		doc, err := html.Parse(strings.NewReader(string(data)))
		if err != nil {
			log.Fatal(err)
		}

		var f func(*html.Node)
		f = func(n *html.Node) {
			if n.Type == html.ElementNode && n.Data == "nav" {
				for _, a := range n.Attr {
					if a.Key == "epub:type" && a.Val == "toc" {
						extractLinks(n, &contents)
					}
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(doc)
	}

	fmt.Println("Contents:")
	for _, c := range contents {
		fmt.Println("-", c)
	}
}

// Helper to extract links from nav list
func extractLinks(n *html.Node, contents *[]string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.TextNode {
				*contents = append(*contents, c.Data)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractLinks(c, contents)
	}
}

func KapMahc(file string) {
	book, err := epub.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer book.Close()

	// Access table of contents
	toc := book.Ncx.Points
	for _, t := range toc {
		fmt.Println(t.Content.Src)
	}
}

// func CreateZipAndFilesInside() {
// 	fmt.Println("creating zip archive")

// 	//Create a new zip archive and named archive.zip
// 	archive, err := os.Create("archive.zip")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer archive.Close()
// 	fmt.Println("archive file created successfully....")

// 	// Create a new zip writer
// 	zipWriter := zip.NewWriter(archive)
// 	// Add files to the zip archive
// 	f1, err := os.Open("test.csv")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer f1.Close()

// 	w1, err := zipWriter.Create("test.csv")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if _, err := io.Copy(w1, f1); err != nil {
// 		panic(err)
// 	}

// 	defer zipWriter.Close()
// }

// func ReadContents(zipName string) {
// 	zipList, err := zip.OpenReader(zipName)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer zipList.Close()

// 	for _, file := range zipList.File {
// 		fmt.Println(file.Name)
// 	}
// }

// func CompressFile(input, gz string) {
// 	// Open the input file
// 	inputFile, err := os.Open(input)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer inputFile.Close()

// 	// Create a new gzip writer
// 	gzipWriter, err := os.Create(gz)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer gzipWriter.Close()

// 	zipWriter := gzip.NewWriter(gzipWriter)
// 	defer zipWriter.Close()

// 	// Copy the input file into gzip writer
// 	if _, err := io.Copy(zipWriter, inputFile); err != nil {
// 		log.Fatal(err)
// 	}
// }
