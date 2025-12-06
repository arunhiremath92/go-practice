package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	categoryUncategorized = "uncategorized"
	categoryHidden        = "hidden"
)

func ReadDirectory(directoryName string) ([]os.DirEntry, error) {
	dirContents, err := os.ReadDir(directoryName)
	if err != nil {
		return nil, err
	}
	return dirContents, nil

}

func SortFiles(files []os.DirEntry) map[string][]string {
	fileCategory := make(map[string][]string)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()
		var fileExtension string
		// Handle hidden files (starts with . but isn't .. or .)
		if strings.HasPrefix(fileName, ".") {
			fileExtension = categoryHidden
		} else {
			ext := filepath.Ext(fileName)
			if ext == "" {
				fileExtension = categoryUncategorized
			} else {
				fileExtension = strings.ToLower(strings.TrimPrefix(ext, "."))
			}
		}
		if fileCategory[fileExtension] == nil {
			fileCategory[fileExtension] = []string{}
		}

		fileCategory[fileExtension] = append(fileCategory[fileExtension], fileName)
	}
	return fileCategory

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("atleast 1 argument when running this program")
		os.Exit(1)
	}

	folderToProcess := os.Args[1]
	dirContents, err := ReadDirectory(folderToProcess)

	if err != nil {
		fmt.Println("failed to read the directory contents", err)
		os.Exit(1)
	}

	filesCategorized := SortFiles(dirContents)

	for extension, fileList := range filesCategorized {
		fmt.Println("files under the extension ", extension)
		for _, file := range fileList {
			fmt.Println("\t", file)
		}
	}

}
