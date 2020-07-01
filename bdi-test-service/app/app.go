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

var (
	configId     int
	sourcefolder string
	dataset      string
)

//ConfigurationInfo ... FileName as JSON description
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

//LogWriter ... to printout logs to http.ResponseWriter
type LogWriter struct {
	http.ResponseWriter
}

var configurations = make(map[string]ConfigurationInfo)
var doOnce sync.Once

//LogWriter ... to log out http.ResponseWriter errors
func (w LogWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	return
}

//Initialize .... initializes the server when first run to create default config
func (s *Server) Initialize() {

	doOnce.Do(func() {
		sourcefolder = "/home/data/stories.table/stories.parquet"
		dataset = "hacker"
		configuration := ConfigurationInfo{SourceFolder: sourcefolder, DatasetName: dataset}
		configurations["default"] = configuration

	})
}

//Run ... run's server for unit testing.
func (s *Server) Run(addr string) {}

//---------------------routes

//ShowFilesRoute ... route
func (s *Server) ShowFilesRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	fileList := make(map[string][]string)
	configId := mux.Vars(r)["configId"]
	if value, ok := configurations[configId]; ok {
		files, err := s.GetFiles(value.SourceFolder)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			w = LogWriter{w}
			w.Write([]byte(err.Error()))
			return
		}
		fileList["files"] = files
		json.NewEncoder(w).Encode(fileList)

	} else {
		http.Error(w, "Bad Request", http.StatusBadRequest)

	}

}

//ShowConfigurationByIdRoute ... route
func (s *Server) ShowConfigurationByIdRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configId := mux.Vars(r)["configId"]
	if value, ok := configurations[configId]; ok {
		json.NewEncoder(w).Encode(value)

	} else {
		http.Error(w, "Bad Request - route does not exist", http.StatusBadRequest)

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
		w = LogWriter{w}
		w.Write([]byte(err.Error()))
		return
	}

	json.Unmarshal(reqBody, &configuration)
	configurations[configIdStr] = configuration
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(configurations)

}

//DeleteConfigurationRoute ... route
func (s *Server) DeleteConfigurationRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configId := mux.Vars(r)["configId"]

	if configId == "default" {
		http.Error(w, "Forbidden", http.StatusForbidden)

	} else {
		delete(configurations, configId)
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(configurations)

	}

}

//SetFileNamesRoute ... route
func (s *Server) SetFileNamesRoute(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	configId := mux.Vars(r)["configId"]
	var sourceFolder string

	if value, ok := configurations[configId]; ok {
		sourceFolder = value.SourceFolder
	} else {
		fmt.Println("error")
	}

	var fileList []FileNameInfo
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal([]byte(body), &fileList)

	if err != nil {
		log.Printf("Cannot unmarshal")
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
