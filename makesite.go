package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"html/template"
	"time"
	"github.com/gomarkdown/markdown"
)

func findFiles(dirPath string, includeMd bool) ([]string, error) {
    var files []string
    err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if !info.IsDir() && (strings.HasSuffix(info.Name(), ".txt") || (includeMd && strings.HasSuffix(info.Name(), ".md"))) {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

func generateHTML(fileName string, tmpl *template.Template, outputFileName string) {
    
    file, err := os.Open(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    
    scanner := bufio.NewScanner(file)
    scanner.Scan() 
    title := scanner.Text()

    var contentBuilder strings.Builder
    for scanner.Scan() {
        contentBuilder.WriteString(scanner.Text() + "\n")
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }

    
    var htmlContent []byte
    if strings.HasSuffix(fileName, ".md") {
        
        htmlContent = markdown.ToHTML([]byte(contentBuilder.String()), nil, nil)
    } else {
        
        
        content := contentBuilder.String()
        htmlContent = []byte("<p>" + strings.ReplaceAll(content, "\n", "</p><p>") + "</p>")
    }

    
    data := struct {
        Title   string
        Content template.HTML 
    }{
        Title:   title,
        Content: template.HTML(htmlContent),
    }

    
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
    dirPtr := flag.String("dir", ".", "Directory to search for text files")
    mdPtr := flag.Bool("md", false, "Include markdown (.md) files")
    flag.Parse()

    tmpl, err := template.ParseFiles("template.tmpl")
    if err != nil {
        log.Fatal(err)
    }

    start := time.Now()
    files, err := findFiles(*dirPtr, *mdPtr)
    if err != nil {
        log.Fatal(err)
    }

    totalBytes := int64(0)

    for _, fileName := range files {
        fmt.Println("Generating HTML for:", fileName)

        
        baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
        outputFileName := baseName + ".html"

        generateHTML(fileName, tmpl, outputFileName)

        
        info, err := os.Stat(outputFileName)
        if err != nil {
            log.Fatal(err)
        }
        totalBytes += info.Size()
    }

    duration := time.Since(start)
    fmt.Printf("\033[1;32mSuccess!\033[0m Generated \033[1m%d\033[0m pages (%.1fkB total) in %.2f seconds.\n", len(files), float64(totalBytes)/1024, duration.Seconds())
}
