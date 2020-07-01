//Package handlers ... workers
package handlers

import (
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

		log.Errorlog.Println(err.Error())
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
		log.Errorlog.Println(err.Error())
	}

	s.Mutex.Unlock()

	c <- files

}

//RemoveFiles ... removes indexed files from default output directory
func (s *Server) RemoveFiles(dir string) error {

	d, err := os.Open(dir)
	if err != nil {
		log.Errorlog.Println(err.Error())
		return err
	}

	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		log.Errorlog.Println(err.Error())
	}
	s.Mutex.Lock()
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			log.Errorlog.Println(err.Error())
		}
	}
	s.Mutex.Unlock()

	return nil
}
