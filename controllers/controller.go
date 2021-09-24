package controllers

import (
	"errors"
	"fmt"
	"github.com/DannyMassa/dead-link-finder/services"
	"github.com/DannyMassa/dead-link-finder/types"
	"log"
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
		resultsChan: make(chan *types.URL, 1000),
		results:     []*types.URL{},
	}
)

type linkController interface {
	PrintResults(config *types.Config) error
	Run(c *types.Config)
}

type LinkControllerImpl struct {
	resultsChan chan *types.URL
	results     []*types.URL
}

func (l *LinkControllerImpl) PrintResults(config *types.Config) error {
	tmpUrls := l.results
	sort.SliceStable(tmpUrls, func(i, j int) bool {
		less := false
		// Sort by Directory
		if strings.ToLower(tmpUrls[i].Directory) == strings.ToLower(tmpUrls[j].Directory) {
			// If Directory == Directory, sort by File
			if strings.ToLower(tmpUrls[i].File) == strings.ToLower(tmpUrls[j].File) {
				// If Directory == Directory && File == File, sort by Line Number
				if tmpUrls[i].LineNumber < tmpUrls[j].LineNumber {
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

	if config.SuccessLogs == false {
		for i := len(tmpUrls) - 1; i >= 0; i-- {
			if tmpUrls[i].Result == "SUCCESS" {
				tmpUrls = append(tmpUrls[:i], tmpUrls[i+1:]...)
			}
		}
	}

	for i := 0; i < len(tmpUrls) - 1; i++ {
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

func (l *LinkControllerImpl) Run(config *types.Config) {
	go l.manageResults()
	waitGroup.Add(len(config.Directories))
	for _, directory := range config.Directories {
		go l.directorySearch(directory, config)
	}

	waitGroup.Wait()
	close(l.resultsChan)
}

// manageResults runs the serve loop, dispatching for checks that need it.
func (l *LinkControllerImpl) manageResults() {
	for p := range l.resultsChan {
		l.results = append(l.results, p)
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
	defer waitGroup.Done()
	var tmp = types.URL{
		Link:      url,
		Directory: directory,
		File:      file,
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
		log.Println("MAYDAY")
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