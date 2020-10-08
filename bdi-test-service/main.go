package main

import (
	"bdi-test-service/handlers"
	logger "bdi-test-service/logging"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	//Setup flags to pass options
	flag.StringVar(&handlers.Sourcepath, "sourcepath", "/home/data/stories.table/stories.parquet", "Location of Data Files")
	flag.StringVar(&handlers.Outputpath, "outputpath", "/home/output", "Location of Indexed Files")
	flag.StringVar(&handlers.Dataset, "dataset", "hacker", "dataset name")
	addr := flag.String("addr", ":8081", "HTTP Service Port Number")
	flag.Parse()

	//Setup router and server
	server := handlers.Server{}
	server.Initialize()
	router := mux.NewRouter()
	router.HandleFunc("/sta/v1/bdi_test_service/configurations", server.ShowAllConfigurationsRoute).Methods("GET")
	router.HandleFunc("/sta/v1/bdi_test_service/configurations/{configID}", server.ShowConfigurationByIDRoute).Methods("GET")
	router.HandleFunc("/sta/v1/bdi_test_service/configurations/{configID}/files", server.ShowFilesRoute).Methods("GET")
	router.HandleFunc("/sta/v1/bdi_test_service/configurations", server.CreateNewConfigurationRoute).Methods("PATCH")
	router.HandleFunc("/sta/v1/bdi_test_service/configurations/{configID}", server.DeleteConfigurationRoute).Methods("DELETE")
	router.HandleFunc("/sta/v1/bdi_test_service/configurations/{configID}/output", server.RemoveIndexedFiles).Methods("DELETE")
	router.HandleFunc("/sta/v1/bdi_test_service/configurations/{configID}/files", server.SetFileNamesRoute).Methods("PATCH")
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bdi-test-service: Bad Request - Check your path", http.StatusBadRequest)
	})
	// TODO: Show Output folder

	srv := &http.Server{
		Handler:      router,
		Addr:         *addr,
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
	}

	//start the server
	go func() {
		logger.Info("Starting server on: %s", *addr)
		err := srv.ListenAndServe()
		if err != nil {
			logger.Errorfatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	logger.Info("Received terminate, graceful shutdown %s", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	srv.Shutdown(tc)
}
