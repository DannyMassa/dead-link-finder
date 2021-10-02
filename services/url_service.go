package services

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

var (
	httpClient = http.Client{
		Timeout: 10 * time.Second,
	}
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

	rex := regexp.MustCompile(`(http|https|ftp|mailto|file|data|irc):(//([0-9a-zA-Z-._~?/#!@$%^&*:\[\]]+@)?[0-9a-zA-Z-._~?/#!@$%^&*:\[\]]+(:[0-9]+)?)?[0-9a-zA-Z-._~?/#!@$%^&*:\[\]]+(\?[0-9a-zA-Z-._~?/#!@$%^&*:\[\]])?(#[0-9a-zA-Z-._~?/#!@$%^&*:\[\]])?`) //nolint
	urls = rex.FindAllStringSubmatch(fileString, -1)

	for i := 1; i < len(urls); i++ {
		urlList = append(urlList, urls[i][0])
	}

	return urlList
}

func (u URLServiceImpl) LinkLivenessChecker(url string) bool {
	resp, err := httpClient.Get(url) //nolint
	if err != nil || resp.StatusCode/100 == 4 || resp.StatusCode/100 == 5 {
		return false
	}
	defer resp.Body.Close()
	return true
}
