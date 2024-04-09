package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

// findTxtFiles uses filepath.Walk to recursively find all .txt files in dirPath and its subdirectories.
func findTxtFiles(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".txt") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// generateHTML reads a .txt file, extracts the title and content, and uses a template to generate an HTML file.
func generateHTML(fileName string, tmpl *template.Template) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan() // Reads the first line, which is the title
	title := scanner.Text()

	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title   string
		Content string
	}{
		Title:   title,
		Content: content,
	}

	outputFileName := strings.TrimSuffix(fileName, ".txt") + ".html"
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

func main() {
	dirPtr := flag.String("dir", ".", "Directory to search for .txt files")
	flag.Parse()

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	files, err := findTxtFiles(*dirPtr)
	if err != nil {
		log.Fatal(err)
	}

	totalBytes := int64(0)

	for _, fileName := range files {
		fmt.Println("Generating HTML for:", fileName)
		generateHTML(fileName, tmpl)
		info, err := os.Stat(strings.TrimSuffix(fileName, ".txt") + ".html")
		if err != nil {
			log.Fatal(err)
		}
		totalBytes += info.Size()
	}

	duration := time.Since(start)
	fmt.Printf("\033[1;32mSuccess!\033[0m Generated \033[1m%d\033[0m pages (%.1fkB total) in %.2f seconds.\n", len(files), float64(totalBytes)/1024, duration.Seconds())
}
