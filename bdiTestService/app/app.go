package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var configId int
var path string
var sourcefolder string
var dataset string

//ConfigurationInfo ... FileName as JSON description
type ConfigurationInfo struct {
	SourceFolder string `json:"sourcefolder"`
	DatasetName  string `json:"datasetname"`
}

//FileNameInfo ... struct defining the JSON schema
type FileNameInfo struct {
	FileNamesOriginal string `json:"oldFileName"`
	FileNamesNew      string `json:"newFileName"`
	Error             error
}

//Server ... server struct
type Server struct {
	Router *mux.Router
}

var configurations = make(map[string]ConfigurationInfo)
var mutex = &sync.Mutex{}
var doOnce sync.Once

//Initialize .... initializes the server when first run to create default config
func (s *Server) Initialize() {

	doOnce.Do(func() {
		sourcefolder = "/home/data"
		dataset = "hacker"
		configuration := ConfigurationInfo{SourceFolder: sourcefolder, DatasetName: dataset}
		configurations["default"] = configuration

	})
}

//---------------------routes

//ShowFilesRoute ... route
func (s *Server) ShowFilesRoute(w http.ResponseWriter, r *http.Request) {

	fileList := make(map[string][]string)
	w.Header().Set("Content-Type", "application/json")
	configId := mux.Vars(r)["configId"]
	if value, ok := configurations[configId]; ok {
		files, err := GetFiles(value.SourceFolder)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		fileList["files"] = files
		json.NewEncoder(w).Encode(fileList)
		return
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

}

//ShowConfigurationByIdRoute ... route
func (s *Server) ShowConfigurationByIdRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configId := mux.Vars(r)["configId"]
	if value, ok := configurations[configId]; ok {
		json.NewEncoder(w).Encode(value)
		return
	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
}

//ShowAllConfigurationsRoute ... route
func (s *Server) ShowAllConfigurationsRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(configurations)
}

//CreateNewConfiguration ... route
func (s *Server) CreateNewConfigurationRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var configuration ConfigurationInfo
	configId++
	configIdStr := strconv.Itoa(configId)
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Cannot Create Configuration")
		return
	}

	json.Unmarshal(reqBody, &configuration)
	if configId > 1 {
		configurations[configIdStr] = configurations["latest"]
		configurations["latest"] = configuration
	} else {
		configurations["latest"] = configuration
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(configurations)

}

//DeleteConfigurationRoute ... route
func (s *Server) DeleteConfigurationRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configId := mux.Vars(r)["configId"]

	if configId == "default" {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	} else {
		delete(configurations, configId)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(configurations)
		return
	}

}

//SetFileNamesRoute ... route
func (s *Server) SetFileNamesRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var fileList []FileNameInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err1 := json.Unmarshal([]byte(body), &fileList)

	if err1 != nil {
		log.Printf("Cannot unmarshal")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	c := make(chan FileNameInfo)

	for _, file := range fileList {

		go RenameFiles(file, c)

	}

	renamedFiles := <-c
	if renamedFiles.Error != nil {
		http.Error(w, "Couldn't rename - Check path, file names", http.StatusBadRequest)
	} else {

		http.Error(w, "Filenames changed", http.StatusOK)
		return
	}

}
