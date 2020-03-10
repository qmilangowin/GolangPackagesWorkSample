package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//GetFiles ... gets files from local directory
func (s *Server) GetFiles(sourceFolder string) ([]string, error) {

	var fileList []string
	s.RwMutex.Lock()
	files, err := ioutil.ReadDir(sourceFolder)
	if err != nil {
		fmt.Println("File doesn't exist")
	}
	for _, file := range files {

		fileList = append(fileList, file.Name())
	}

	s.RwMutex.Unlock()
	return fileList, err

}

//RenameFiles ... renames files in local directory
func (s *Server) RenameFiles(files FileNameInfo, c chan FileNameInfo) {

	dirPath := files.SourceFolder
	var originalFileName = filepath.Join(dirPath, files.FileNamesOriginal)
	var newFileName = filepath.Join(dirPath, files.FileNamesNew)
	s.Mutex.Lock()

	err := os.Rename(originalFileName, newFileName)
	if err != nil {
		files.Error = errors.New("Bad Request")
	}

	s.Mutex.Unlock()

	c <- files

}
