package main

import (
	"fmt"
	"github.com/DannyMassa/dead-link-finder/services"
	flag "github.com/spf13/pflag"
)

var (
	DefaultDirectories                    = []string{"./"}
	directories, fileEndings, files, urls []string
	pass                                  = true
	DefaultFileEndings                    = []string{".markdown", ".mdown", ".mkdn", ".md", ".mkd", ".mdwn", ".mdtxt",
		".mdtext", "text", ".rmd", "rst"}

	urlService       services.URLService       = &services.URLServiceImpl{}
	directoryService services.DirectoryService = &services.DirectoryServiceImpl{}
)

func main() {
	flag.StringArrayVar(&fileEndings, "file_endings", DefaultFileEndings, "File Extensions to look for")
	flag.StringArrayVar(&directories, "directories", DefaultDirectories, "Directories to search for files")
	flag.Parse()

	for _, directory := range directories {
		fmt.Printf("%s\n", directory)
		files = directoryService.FindFiles(directory, fileEndings)
		for _, file := range files {
			urls = urlService.URLScraper(file)
			fmt.Printf("    %s\n", file)
			for _, url := range urls {
				if urlService.LinkLivenessChecker(url) {
					fmt.Printf("        PASS: %s\n", url)
				} else {
					fmt.Printf("        FAIL: %s\n", url)
					pass = false
				}
			}
		}
	}

	if !pass {
		panic("Some URLs are not accessible")
	}
}
