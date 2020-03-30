package main

import (
	"bdi-test-service/app"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "bdi-test-service", log.LstdFlags)
	app := app.Server{}
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

	//start the server
	go func() {
		l.Println("Starting server on port: 8081")
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(tc)

}