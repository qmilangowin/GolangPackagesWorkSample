//Package handlers ... workers
package handlers

import (
	logger "bdi-test-service/logging"
	"context"
	"errors"
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

		logger.Errorln(err.Error())
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
		logger.Error(err.Error())
	}

	s.Mutex.Unlock()

	c <- files

}

//RemoveFiles ... removes indexed files from default output directory
func (s *Server) RemoveFiles(ctx context.Context, dir string) error {

	done := make(chan bool, 1)

	d, err := os.Open(dir)
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}

	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		logger.Errorln(err.Error())
		return err
	}
	go func() {
		for _, name := range names {
			err := os.RemoveAll(filepath.Join(dir, name))
			if err != nil {
				logger.Errorln(err.Error())
				done <- false
			}
		}
		done <- true
	}()

	select {
	case result := <-done:
		if result == true {
			return nil
		}
		return err

	case <-ctx.Done():
		logger.Errorln(ctx.Err())
		return ctx.Err()
	}
}
