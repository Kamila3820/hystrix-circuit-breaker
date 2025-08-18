package lib

import (
	"fmt"
	"log"

	"github.com/kapmahc/epub"
)

func Barsanuphe(file string) {
	book, err := epub.Open(file)
	if err != nil {
		log.Panic(err)
	}
	defer book.Close()
	title := book.Ncx.Points
	var count int
	for i, t := range title {
		// if t.Text == "" {
		// 	continue
		// }
		fmt.Printf("Chapter %d: %s\n", i+1, t.Text)
		count++
	}
}
