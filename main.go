package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/DannyMassa/dead-link-finder/services"
	"gopkg.in/yaml.v2"
)

var (
	successCount, failCount, skipCount = 0, 0, 0
	files, urls                        []string
	pass                                                         = true
	urlService                         services.URLService       = &services.URLServiceImpl{}
	directoryService                   services.DirectoryService = &services.DirectoryServiceImpl{}
)

type Config struct {
	FileExtensions []string `yaml:"fileExtensions"`
	Directories    []string `yaml:"directories"`
	Ignored        []string `yaml:"ignored"`
}

func main() {
	c := &Config{}
	buf, err := ioutil.ReadFile(".deadlink")
	if err != nil {
		fmt.Printf("Could not find .deadlink file, using defaults\n")
		c.FileExtensions = []string{".markdown", ".mdown", ".mkdn", ".md", ".mkd", ".mdwn", ".mdtxt", ".mdtext",
			".text", ".txt", ".rmd", ".rst"}
		c.Directories = []string{"./"}
		c.Ignored = []string{}
	}

	err = yaml.Unmarshal(buf, c)
	if err != nil {
		panic("Could not parse .deadlink file")
	}

	for _, directory := range c.Directories {
		fmt.Printf("%s\n", directory)
		files = directoryService.FindFiles(directory, c.FileExtensions)
		for _, file := range files {
			urls = urlService.URLScraper(file)
			fmt.Printf("    %s\n", file)
			for _, url := range urls {
				if contains(url, c.Ignored) { //nolint
					fmt.Printf("        SKIP: %s\n", url)
					skipCount++
				} else if urlService.LinkLivenessChecker(url) {
					fmt.Printf("        SUCCESS: %s\n", url)
					successCount++
				} else {
					fmt.Printf("        FAILURE: %s\n", url)
					failCount++
					pass = false
				}
			}
		}
	}

	fmt.Printf("\n\nSUCCESS: %d     FAILURE: %d     SKIPPED: %d\n", successCount, failCount, skipCount)

	if !pass {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func contains(string1 string, strings []string) bool {
	for _, myString := range strings {
		if string1 == myString {
			return true
		}
	}

	return false
}
