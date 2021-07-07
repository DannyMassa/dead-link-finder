package services

import (
	"io/ioutil"
	"net/http"
	"regexp"
)

type URLService interface {
	URLScraper(file string) []string
	LinkLivenessChecker(url string) bool
}

type URLServiceImpl struct{}

func (u URLServiceImpl) URLScraper(file string) []string {
	var urls [][]string
	fileByteSlice, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Could not parse file... crashing")
	}
	fileString := string(fileByteSlice)
	var urlList []string
	rex := regexp.MustCompile(`(http|ftp|https)://([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?`)
	urls = rex.FindAllStringSubmatch(fileString, -1)

	for i := 1; i < len(urls); i++ {
		urlList = append(urlList, urls[i][0])
	}

	return urlList
}

func (u URLServiceImpl) LinkLivenessChecker(url string) bool {
	resp, err := http.Get(url) //nolint
	if err != nil || resp.StatusCode/100 == 4 || resp.StatusCode/100 == 5 {
		return false
	}
	defer resp.Body.Close()
	return true
}
