package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mbndr/figlet4go"
)

func main() {
	ascii := figlet4go.NewAsciiRender()
	options := figlet4go.NewRenderOptions()
	options.FontColor = []figlet4go.Color{
		figlet4go.ColorGreen,
	}
	options.FontName = "larry3d"
	renderStr, _ := ascii.RenderOpts("FINDER TOOL", options)
	fmt.Println(renderStr)

	// flags
	zipFilePath := flag.String("file", "", "Specify the compressed file path - required")
	searchText := flag.String("text", "", "Specify the text to search for or specify multiple texts separated by (,) - required")
	outputFile := flag.String("output", "", "Specify the path to the output file to save the results - optional")
	caseSensitive := flag.Bool("case-sensitive", false, "Specify whether the search should be case-sensitive (default: false) - optional")
	helpFlag := flag.Bool("help", false, "Help for using the finder tool")
	flag.Parse()

	// control
	if *helpFlag {
		fmt.Println("Flags")
		flag.PrintDefaults()
		return
	}

	if *zipFilePath == "" || *searchText == "" {
		fmt.Println("Please run the -help command to use the Finder tool.")
		return
	}

	// separate words
	searchKeywords := strings.Split(*searchText, ",")

	// open
	zipFile, err := zip.OpenReader(*zipFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	// results
	var searchResults []string

	var totalSize int64
	fileCount := 0

	for _, file := range zipFile.File {
		if file.FileInfo().IsDir() {
			continue
		}

		fileCount++

		// increment total size
		totalSize += file.FileInfo().Size()

		// read
		f, err := file.Open()
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(f)
		lineNumber := 1

		// scan
		for scanner.Scan() {
			line := scanner.Text()
			compareLine := line
			if !*caseSensitive {
				compareLine = strings.ToLower(line)
			}

			for _, keyword := range searchKeywords {
				compareKeyword := keyword
				if !*caseSensitive {
					compareKeyword = strings.ToLower(keyword)
				}

				if strings.Contains(compareLine, compareKeyword) {
					result := fmt.Sprintf("File: %s\n", file.Name)
					result += fmt.Sprintf("Line Number: %d\n", lineNumber)
					result += fmt.Sprintf("Line: %s\n", line)
					result += "----------\n"
					searchResults = append(searchResults, result)
					break
				}
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		f.Close()
	}

	// total size and file count
	fmt.Printf("Compressed File Properties\n")
	fmt.Printf("Total Size: %d bytes\n", totalSize)
	fmt.Printf("File Count: %d\n", fileCount)
	fmt.Printf("----------\n\n")

	// output
	if len(searchResults) > 0 {
		if *outputFile != "" {
			output, err := os.Create(*outputFile)
			if err != nil {
				log.Fatal(err)
			}
			defer output.Close()

			for _, result := range searchResults {
				_, err := output.WriteString(result)
				if err != nil {
					log.Fatal(err)
				}
			}

			fmt.Println("Search results were saved to", *outputFile)
		} else {
			for _, result := range searchResults {
				fmt.Println(result)
			}
		}
	} else {
		fmt.Println("The searched text was not found.")
	}
}
