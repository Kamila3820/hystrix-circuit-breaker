package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	ReadContents("archive.zip")
}

func CreateZipAndFilesInside() {
	fmt.Println("creating zip archive")

	//Create a new zip archive and named archive.zip
	archive, err := os.Create("archive.zip")
	if err != nil {
		panic(err)
	}
	defer archive.Close()
	fmt.Println("archive file created successfully....")

	// Create a new zip writer
	zipWriter := zip.NewWriter(archive)
	// Add files to the zip archive
	f1, err := os.Open("test.csv")
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	w1, err := zipWriter.Create("test.csv")
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(w1, f1); err != nil {
		panic(err)
	}

	defer zipWriter.Close()
}

func ReadContents(zipName string) {
	zipList, err := zip.OpenReader(zipName)
	if err != nil {
		log.Fatal(err)
	}
	defer zipList.Close()

	for _, file := range zipList.File {
		fmt.Println(file.Name)
	}
}
