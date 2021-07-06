package main

import (
	"fmt"
	"github.com/DannyMassa/dead-link-finder/services"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var (
	passCount, failCount, skipCount = 0, 0, 0
	files, urls                     []string
	pass                                                      = true
	urlService                      services.URLService       = &services.URLServiceImpl{}
	directoryService                services.DirectoryService = &services.DirectoryServiceImpl{}
)

type Config struct {
	FileExtensions []string `yaml:"fileExtensions"`
	Directories    []string `yaml:"directories"`
	Ignored        []string `yaml:"ignored"`
}

func main() {
	buf, err := ioutil.ReadFile(".deadlink")
	if err != nil {
		panic("")
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		panic("")
	}

	for _, directory := range c.Directories {
		fmt.Printf("%s\n", directory)
		files = directoryService.FindFiles(directory, c.FileExtensions)
		for _, file := range files {
			urls = urlService.URLScraper(file)
			fmt.Printf("    %s\n", file)
			for _, url := range urls {
				if contains(url, c.Ignored) {
					fmt.Printf("        SKIP: %s\n", url)
					skipCount++
				} else if urlService.LinkLivenessChecker(url) {
					fmt.Printf("        PASS: %s\n", url)
					passCount++
				} else {
					fmt.Printf("        FAIL: %s\n", url)
					failCount++
					pass = false
				}
			}
		}
	}

	fmt.Printf("\n\nPASSED: %d     FAILED: %d     SKIPPED: %d\n", passCount, failCount, skipCount)

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
