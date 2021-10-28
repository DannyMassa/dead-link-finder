package controllers

import (
	"log"
	"os"
	"sync"

	"github.com/DannyMassa/dead-link-linter/services"
	"github.com/DannyMassa/dead-link-linter/types"
)

var (
	waitGroup                          sync.WaitGroup
	successCount, failCount, skipCount                           = 0, 0, 0
	urlService                         services.URLService       = &services.URLServiceImpl{}
	directoryService                   services.DirectoryService = &services.DirectoryServiceImpl{}
	logService                         services.LogService       = &services.LogServiceImpl{}
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

func (l *LinkControllerImpl) Run(config *types.Config) error {
	go l.manageResults()
	waitGroup.Add(len(config.Directories))
	for _, directory := range config.Directories {
		go l.directorySearch(directory, config)
	}

	waitGroup.Wait()
	close(l.resultsChan)

	return logService.PrintResults(l.results, config)
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
