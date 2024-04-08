package main

import (
	"bufio"
	"flag"
	// "fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

func main() {
	fileName := flag.String("file", "first-post.txt", "The name of the text file to read in")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // This reads the first line, which is the title
	title := scanner.Text()

	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title   string
		Content string
	}{
		Title:   title,
		Content: content,
	}

	outputFileName := strings.TrimSuffix(*fileName, ".txt") + ".html"
	outFile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	err = tmpl.Execute(outFile, data)
	if err != nil {
		log.Fatal(err)
	}
}

