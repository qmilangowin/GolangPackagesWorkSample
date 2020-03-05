package main

import (
	"Repos/BDI_Service/bdiTestService/app"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	app.Initialize()
	router := mux.NewRouter()
	router.HandleFunc("/v1/bdi_test_service/configurations/{configId}/files", app.ShowFilesRoute).Methods("GET")
	router.HandleFunc("/v1/bdi_test_service/configurations/{configId}", app.ShowConfigurationByIdRoute).Methods("GET")
	router.HandleFunc("/v1/bdi_test_service/configurations", app.ShowAllConfigurationsRoute).Methods("GET")
	router.HandleFunc("/v1/bdi_test_service/configurations", app.CreateNewConfigurationRoute).Methods("PATCH")
	router.HandleFunc("/v1/bdi_test_service/configurations/{configId}", app.DeleteConfigurationRoute).Methods("DELETE")
	router.HandleFunc("/v1/bdi_test_service/configurations/{configId}/files", app.SetFileNamesRoute).Methods("PATCH")

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8081",
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
