package controllers

import (
	"errors"
	"fmt"
	"github.com/DannyMassa/dead-link-linter/services"
	"github.com/DannyMassa/dead-link-linter/types"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

var (
	waitGroup                          sync.WaitGroup
	successCount, failCount, skipCount                           = 0, 0, 0
	urlService                         services.URLService       = &services.URLServiceImpl{}
	directoryService                   services.DirectoryService = &services.DirectoryServiceImpl{}
	Controller                         linkController            = &LinkControllerImpl{
		resultsChan: make(chan *types.URL, 512),
		results:     []*types.URL{},
	}
)

type linkController interface {
	Run(c *types.Config) error
}

type LinkControllerImpl struct {
	resultsChan chan *types.URL
	results     []*types.URL
}

func (l *LinkControllerImpl) printResults(config *types.Config) error {
	tmpUrls := l.results
	sort.SliceStable(tmpUrls, func(i, j int) bool {
		less := false
		// Sort by Directory
		if strings.EqualFold(tmpUrls[i].Directory, tmpUrls[j].Directory) { //nolint
			// If Directory == Directory, sort by File
			if strings.EqualFold(tmpUrls[i].File, tmpUrls[j].File) { //nolint
				// If Directory == Directory && File == File, sort alphabetically
				if strings.ToLower(tmpUrls[i].Link) < strings.ToLower(tmpUrls[j].Link) {
					less = true
				} else {
					less = false
				}
			} else if strings.ToLower(tmpUrls[i].File) < strings.ToLower(tmpUrls[j].File) {
				less = true
			} else {
				less = false
			}
		} else if strings.ToLower(tmpUrls[i].Directory) < strings.ToLower(tmpUrls[j].Directory) {
			less = true
		} else {
			less = false
		}

		return less
	})

	if config.LogVerbosity <= 1 {
		for i := len(tmpUrls) - 1; i >= 0; i-- {
			if tmpUrls[i].Result == "SUCCESS" {
				tmpUrls = append(tmpUrls[:i], tmpUrls[i+1:]...)
			}
		}
	}

	for i := 0; i < len(tmpUrls); i++ {
		if i == 0 {
			fmt.Printf("%s\n", tmpUrls[i].Directory)
			fmt.Printf("    %s\n", tmpUrls[i].File)
			fmt.Printf("        %s - %s\n", tmpUrls[i].Result, tmpUrls[i].Link)
		} else {
			if tmpUrls[i].Directory != tmpUrls[i-1].Directory {
				fmt.Printf("%s\n", tmpUrls[i].Directory)
			}
			if tmpUrls[i].File != tmpUrls[i-1].File {
				fmt.Printf("    %s\n", tmpUrls[i].File)
			}
			fmt.Printf("        %s - %s\n", tmpUrls[i].Result, tmpUrls[i].Link)
		}
	}

	fmt.Printf("\n\nSUCCESS: %d     FAILED: %d     SKIPPED: %d\n", successCount, failCount, skipCount)

	if failCount > config.MaxFailures {
		return errors.New("wow")
	}

	return nil
}

func (l *LinkControllerImpl) Run(config *types.Config) error {
	go l.manageResults()
	waitGroup.Add(len(config.Directories))
	for _, directory := range config.Directories {
		go l.directorySearch(directory, config)
	}

	waitGroup.Wait()
	close(l.resultsChan)

	return l.printResults(config)
}

// manageResults runs the serve loop, dispatching for checks that need it.
func (l *LinkControllerImpl) manageResults() {
	for p := range l.resultsChan {
		l.results = append(l.results, p)
		waitGroup.Done()
	}
}

func (l *LinkControllerImpl) directorySearch(directory string, config *types.Config) {
	defer waitGroup.Done()
	files := directoryService.FindFiles(directory, config.FileExtensions)
	for _, file := range files {
		l.fileSearch(directory, file, config)
	}
}

func (l *LinkControllerImpl) fileSearch(directory string, file string, config *types.Config) {
	urls := urlService.URLScraper(file)
	for _, url := range urls {
		waitGroup.Add(1)
		go l.urlCheck(directory, file, url, config)
	}
}

func (l *LinkControllerImpl) urlCheck(directory string, file string, url string, config *types.Config) {
	var tmp = types.URL{
		Link:      url,
		Directory: directory,
		File:      file,
		Result:    "",
	}
	if l.contains(url, config.Ignored) { //nolint
		tmp.Result = "SKIPPED"
		skipCount++
	} else if urlService.LinkLivenessChecker(url) {
		tmp.Result = "SUCCESS"
		successCount++
	} else {
		tmp.Result = "FAILURE"
		failCount++
	}

	select {
	case l.resultsChan <- &tmp:
	default:
		log.Println("Channel Overload, Dead Link Linter must exit")
		os.Exit(6)
	}
}

func (l *LinkControllerImpl) contains(string1 string, strings []string) bool {
	for _, myString := range strings {
		if string1 == myString {
			return true
		}
	}

	return false
}
