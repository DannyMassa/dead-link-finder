package services

import (
	"errors"
	"fmt"
	"github.com/DannyMassa/dead-link-linter/types"
	"sort"
	"strings"
)

type LogService interface {
	PrintResults(directory string, fileEndings []string) []string
}

type LogServiceImpl struct{}

func (l LogServiceImpl) PrintResults(results []*types.URL, config *types.Config) error {
	successCount, failCount, skipCount := 0, 0, 0
	tmpUrls := results
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
		return errors.New("maximum number of allowed failures was surpassed")
	}

	return nil
}
