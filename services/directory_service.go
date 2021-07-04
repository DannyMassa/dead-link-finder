package services

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DirectoryService interface {
	FindFiles(directory string, fileEndings []string) []string
}

type DirectoryServiceImpl struct{}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		*files = append(*files, path)
		return nil
	}
}

func (d DirectoryServiceImpl) FindFiles(directory string, fileEndings []string) []string {
	var files, filesWithCorrectEndings []string
	err := filepath.Walk(directory, visit(&files))
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		for _, fileEnding := range fileEndings {
			if strings.HasSuffix(file, fileEnding) {
				filesWithCorrectEndings = append(filesWithCorrectEndings, file)
				break
			}
		}
	}

	return filesWithCorrectEndings
}
