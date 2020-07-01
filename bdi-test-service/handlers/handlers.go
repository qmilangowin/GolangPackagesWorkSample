//Package handlers ... handlers
package handlers

import (
	"bdi-test-service/config"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var (
	Sourcepath     string
	Outputpath     string
	Dataset        string
	configID       int
	configurations = make(map[string]ConfigurationInfo)
	doOnce         sync.Once
	log            *config.Logging
)

//ConfigurationInfo ... Filename as JSON Description
type ConfigurationInfo struct {
	SourceFolder string `json:"sourcefolder"`
	DatasetName  string `json:"datasetname"`
}

//FileNameInfo ... struct defining the JSON schema
type FileNameInfo struct {
	FileNamesOriginal string `json:"oldFileName"`
	FileNamesNew      string `json:"newFileName"`
	SourceFolder      string
	Error             error
}

//Server ... defined as struct
type Server struct {
	Router  *mux.Router
	Mutex   sync.Mutex
	RwMutex sync.RWMutex
}

//Initialize .... initializes the server when first run to create default config
func (s *Server) Initialize() {

	doOnce.Do(func() {
		sourcefolder := Sourcepath
		dataset := Dataset
		configuration := ConfigurationInfo{SourceFolder: sourcefolder, DatasetName: dataset}
		configurations["default"] = configuration
		log = config.Logger()

	})
}

//ShowAllConfigurationsRoute ... route
func (s *Server) ShowAllConfigurationsRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configurations)

}

//ShowConfigurationByIDRoute ... route
func (s *Server) ShowConfigurationByIDRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configID := mux.Vars(r)["configID"]
	if value, ok := configurations[configID]; ok {
		json.NewEncoder(w).Encode(value)

	} else {
		http.Error(w, "Bad Request - route does not exist", http.StatusBadRequest)

	}
}

//ShowFilesRoute ... displays the files for a given configuration
func (s *Server) ShowFilesRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	fileList := make(map[string][]string)
	configID := mux.Vars(r)["configID"]
	if value, ok := configurations[configID]; ok {
		files, err := s.GetFiles(value.SourceFolder)
		if err != nil {

			http.Error(w, "Bad Request - no such file/directory", http.StatusBadRequest)
			return
		}
		fileList["files"] = files
		json.NewEncoder(w).Encode(fileList)

	} else {

		http.Error(w, "Bad Request", http.StatusBadRequest)

	}

}

//CreateNewConfigurationRoute ... route
func (s *Server) CreateNewConfigurationRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var configuration ConfigurationInfo
	configID++
	configIDStr := strconv.Itoa(configID)
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {

		log.Errorlog.Println(w, "Cannot Create Configuration")
		return
	}

	json.Unmarshal(reqBody, &configuration)
	configurations[configIDStr] = configuration
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(configurations)
	log.Infolog.Printf("Configuration %d created ", configID)

}

//DeleteConfigurationRoute ... route
func (s *Server) DeleteConfigurationRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configID := mux.Vars(r)["configID"]

	if configID == "default" {

		http.Error(w, "Forbidden", http.StatusForbidden)

	} else {

		if _, ok := configurations[configID]; ok {
			delete(configurations, configID)
			w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(configurations)
			log.Infolog.Printf("Configuration %s deleted ", configID)

		} else {
			http.Error(w, "Bad Request - Config not found", http.StatusBadRequest)
			log.Infolog.Printf("Configuration %s not found ", configID)
		}

	}

}

//RemoveIndexedFiles ... purges the output folder
func (s *Server) RemoveIndexedFiles(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	err := s.RemoveFiles(Outputpath)
	if err != nil {
		http.Error(w, "Bad Request - Could not delete output folder", http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusAccepted)
		log.Infolog.Println("Output Folder Purged")
	}
}

//SetFileNamesRoute ... route
func (s *Server) SetFileNamesRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configID := mux.Vars(r)["configID"]
	var sourceFolder string

	if value, ok := configurations[configID]; ok {
		sourceFolder = value.SourceFolder
	} else {
		log.Infolog.Println("Could not read configuration")
	}

	var fileList []FileNameInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorlog.Printf("Error reading body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(body), &fileList)

	if err != nil {
		log.Errorlog.Printf("Cannot unmarshal")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	for index := range fileList {
		fileList[index].SourceFolder = sourceFolder
	}

	c := make(chan FileNameInfo)

	for _, file := range fileList {

		go s.RenameFiles(file, c)

	}

	renamedFiles := <-c
	if renamedFiles.Error != nil {
		http.Error(w, "Couldn't rename - Check path, file names", http.StatusBadRequest)

	} else {

		http.Error(w, "Filenames changed", http.StatusOK)

	}

}
