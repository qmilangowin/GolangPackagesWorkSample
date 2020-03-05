package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//GetFiles ... gets files from local directory
func GetFiles(sourceFolder string) ([]string, error) {
	var fileList []string
	files, err := ioutil.ReadDir(sourceFolder)
	if err != nil {
		fmt.Println("File doesn't exist")
	}
	for _, file := range files {

		fileList = append(fileList, file.Name())
	}

	return fileList, err

}

//RenameFiles ... renames files in local directory
func RenameFiles(files FileNameInfo, c chan FileNameInfo) {

	dirPath := sourcefolder
	var originalFileName = filepath.Join(dirPath, files.FileNamesOriginal)
	var newFileName = filepath.Join(dirPath, files.FileNamesNew)
	mutex.Lock()

	err := os.Rename(originalFileName, newFileName)
	if err != nil {
		files.Error = errors.New("Bad Request")
	}

	mutex.Unlock()

	c <- files

}
