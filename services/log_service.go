package services

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/DannyMassa/dead-link-linter/types"
)

type LogService interface {
	PrintResults(results []*types.URL, config *types.Config) error
}

type LogServiceImpl struct{}

func (l LogServiceImpl) PrintResults(results []*types.URL, config *types.Config) error {
	successCount, failCount, skipCount := 0, 0, 0
	tmpResults := results

	// Get results summary before manipulating results for output
	for _, result := range tmpResults {
		if result.Result == "SUCCESS" {
			successCount++
		}
		if result.Result == "FAILURE" {
			failCount++
		}
		if result.Result == "SKIPPED" {
			skipCount++
		}
	}

	// Remove objects not going to be used in output
	if config.LogVerbosity < 2 {
		if config.LogVerbosity < 1 {
			for i := len(tmpResults) - 1; i >= 0; i-- {
				if tmpResults[i].Result == "SUCCESS" || tmpResults[i].Result == "SKIPPED" {
					tmpResults = append(tmpResults[:i], tmpResults[i+1:]...)
				}
			}
		} else {
			for i := len(tmpResults) - 1; i >= 0; i-- {
				if tmpResults[i].Result == "SUCCESS" {
					tmpResults = append(tmpResults[:i], tmpResults[i+1:]...)
				}
			}
		}
	}

	// Sort Output
	sort.SliceStable(tmpResults, func(i, j int) bool {
		less := false
		// Sort by Directory
		if strings.EqualFold(tmpResults[i].Directory, tmpResults[j].Directory) { //nolint
			// If Directory == Directory, sort by File
			if strings.EqualFold(tmpResults[i].File, tmpResults[j].File) { //nolint
				// If Directory == Directory && File == File, sort alphabetically
				if strings.ToLower(tmpResults[i].Link) < strings.ToLower(tmpResults[j].Link) {
					less = true
				} else {
					less = false
				}
			} else if strings.ToLower(tmpResults[i].File) < strings.ToLower(tmpResults[j].File) {
				less = true
			} else {
				less = false
			}
		} else if strings.ToLower(tmpResults[i].Directory) < strings.ToLower(tmpResults[j].Directory) {
			less = true
		} else {
			less = false
		}

		return less
	})

	// Print Output
	for i := 0; i < len(tmpResults); i++ {
		if i == 0 {
			fmt.Printf("%s\n", tmpResults[i].Directory)
			fmt.Printf("    %s\n", tmpResults[i].File)
			fmt.Printf("        %s - %s\n", tmpResults[i].Result, tmpResults[i].Link)
		} else {
			if tmpResults[i].Directory != tmpResults[i-1].Directory {
				fmt.Printf("%s\n", tmpResults[i].Directory)
			}
			if tmpResults[i].File != tmpResults[i-1].File {
				fmt.Printf("    %s\n", tmpResults[i].File)
			}
			fmt.Printf("        %s - %s\n", tmpResults[i].Result, tmpResults[i].Link)
		}
	}

	fmt.Printf("\n\nSUCCESS: %d     FAILED: %d     SKIPPED: %d\n", successCount, failCount, skipCount)

	if failCount > config.MaxFailures {
		return errors.New("maximum number of allowed failures was surpassed")
	}

	return nil
}
